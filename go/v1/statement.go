/*
Wrapper APIs for in-toto attestation Statement layer protos.
*/

package v1

import "errors"

const StatementTypeUri = "https://in-toto.io/Statement/v1"

func (s *Statement) Validate() error {
	if s.GetType() != StatementTypeUri {
		return errors.New("Wrong statement type")
	}

	if s.GetSubject() == nil || len(s.GetSubject()) == 0 {
		return errors.New("At least one subject required")
	}

	// check all resource descriptors in the subject
	subject := s.GetSubject()
	for _, rd := range subject {
		if err := rd.Validate(); err != nil {
			return err
		}

		// v1 statements require the digest to be set in the subject
		if len(rd.GetDigest()) == 0 {
			return errors.New("At least one digest required")
		}
	}

	if s.GetPredicateType() == "" {
		return errors.New("Predicate type required")
	}

	if s.GetPredicate() == nil {
		return errors.New("Predicate object required")
	}

	return nil
}
