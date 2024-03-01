/*
Wrapper APIs for in-toto attestation ResourceDescriptor protos.
*/

package v1

import (
	"encoding/hex"
	"errors"
)

var (
	ErrInvalidDigestEncoding = errors.New("digest is not valid hex-encoded string")
	ErrRDRequiredField       = errors.New("at least one of name, URI, or digest are required")
)

func (d *ResourceDescriptor) Validate() error {
	// at least one of name, URI or digest are required
	if d.GetName() == "" && d.GetUri() == "" && len(d.GetDigest()) == 0 {
		return ErrRDRequiredField
	}

	// the in-toto spec expects a hex-encoded string in DigestSets
	// https://github.com/in-toto/attestation/blob/main/spec/v1/digest_set.md
	if len(d.GetDigest()) > 0 {
		for _, digest := range d.GetDigest() {
			_, err := hex.DecodeString(digest)

			if err != nil {
				return ErrInvalidDigestEncoding
			}
		}
	}

	return nil
}
