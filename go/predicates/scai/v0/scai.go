/*
Wrapper APIs for SCAI AttributeAssertion and AttributeReport protos.
*/

package v0

import "fmt"

func (a *AttributeAssertion) Validate() error {
	// at least the attribute field is required
	if a.GetAttribute() == "" {
		return fmt.Errorf("The attribute field is required")
	}

	// check target and evidence are valid ResourceDescriptors
	if a.GetTarget() != nil {
		if err := a.GetTarget().Validate(); err != nil {
			return fmt.Errorf("Target validation failed with an error: %w", err)
		}
	}

	if a.GetEvidence() != nil {
		if err := a.GetEvidence().Validate(); err != nil {
			return fmt.Errorf("Evidence validation failed with an error: %w", err)
		}
	}

	return nil
}

func (r *AttributeReport) Validate() error {
	// at least the attributes field is required
	attrs := r.GetAttributes()
	if attrs == nil || len(attrs) == 0 {
		return fmt.Errorf("At least one AttributeAssertion is required")
	}

	// ensure all AttributeAssertions are valid
	for _, a := range attrs {
		if err := a.Validate(); err != nil {
			return fmt.Errorf("AttributeAssertion validation failed with an error: %w", err)
		}
	}

	// ensure the producer is a valid ResourceDescriptor
	if r.GetProducer() != nil {
		if err := r.GetProducer().Validate(); err != nil {
			return fmt.Errorf("Producer validation failed with an error: %w", err)
		}
	}

	return nil
}
