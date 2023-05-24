/*
Wrapper APIs for SCAI AttributeAssertion and AttributeReport protos.
*/

package v0

import "errors"

func (a *AttributeAssertion) Validate() error {
	// at least the attribute field is required
	if a.GetAttribute() == "" {
		return errors.New("The attribute field is required")
	}

	// check target and evidence are valid ResourceDescriptors
	if a.GetTarget() != nil {
		if err := a.GetTarget().Validate(); err != nil {
			return err
		}
	}

	if a.GetEvidence() != nil {
		if err := a.GetEvidence().Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (r *AttributeReport) Validate() error {
	// at least the attributes field is required
	attrs := r.GetAttributes()
	if attrs == nil || len(attrs) == 0 {
		return errors.New("At least one AttributeAssertion required")
	}

	// ensure all AttributeAssertions are valid
	for _, a := range attrs {
		if err := a.Validate(); err != nil {
			return err
		}
	}

	// ensure the producer is a valid ResourceDescriptor
	if r.GetProducer() != nil {
		if err := r.GetProducer().Validate(); err != nil {
			return err
		}
	}

	return nil
}
