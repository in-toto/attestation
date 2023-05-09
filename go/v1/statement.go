/*
Wrapper APIs for in-toto attestation Statement layer protos.
*/

package v1

const StatementTypeUri = "https://in-toto.io/Statement/v1"

func (s *Statement) Validate() bool {
	if s.GetType() != StatementTypeUri {
		return false
	}

	if s.GetSubject() == nil || len(s.GetSubject()) == 0 {
		return false
	}

	// check all resource descriptors in the subject
	subject := s.GetSubject()
	for _, rd := range subject {
		if !rd.Validate() {
			return false
		}

		// v1 statements require the digest to be set in the subject
		if len(rd.GetDigest()) == 0 {
			return false
		}
	}

	if s.GetPredicateType() == "" {
		return false
	}

	if s.GetPredicate() == nil {
		return false
	}

	return true
}
