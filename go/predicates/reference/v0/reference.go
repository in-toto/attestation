/*
Validator APIs for reference v0 protos.
*/
package v0

import (
	"errors"
	"fmt"
)

var (
	ErrAttesterRequired         = errors.New("attester required")
	ErrAttesterIdRequired       = errors.New("attester.id required")
	ErrReferencesRequired       = errors.New("at least one reference required")
	ErrReferenceRequired        = errors.New("reference must not be empty")
	ErrDownloadLocationRequired = errors.New("downloadLocation required")
	ErrMediaTypeRequired        = errors.New("mediaType required")
)

func (a *Reference_Attester) Validate() error {
	if a.GetId() == "" {
		return ErrAttesterIdRequired
	}
	return nil
}

// Validate checks the fields reference.md marks as required. Each entry in
// references has to be a valid ResourceDescriptor with a downloadLocation and
// a mediaType.
func (r *Reference) Validate() error {
	attester := r.GetAttester()
	if attester == nil {
		return ErrAttesterRequired
	}
	if err := attester.Validate(); err != nil {
		return err
	}

	refs := r.GetReferences()
	if len(refs) == 0 {
		return ErrReferencesRequired
	}
	for i, ref := range refs {
		if ref == nil {
			return fmt.Errorf("references[%d]: %w", i, ErrReferenceRequired)
		}
		if ref.GetDownloadLocation() == "" {
			return fmt.Errorf("references[%d]: %w", i, ErrDownloadLocationRequired)
		}
		if ref.GetMediaType() == "" {
			return fmt.Errorf("references[%d]: %w", i, ErrMediaTypeRequired)
		}
		if err := ref.Validate(); err != nil {
			return fmt.Errorf("references[%d]: %w", i, err)
		}
	}

	return nil
}
