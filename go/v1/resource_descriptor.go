/*
Wrapper APIs for in-toto attestation ResourceDescriptor protos.
*/

package v1

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"errors"
)

var (
	ErrIncorrectDigestLength = errors.New("digest is not correct length")
	ErrInvalidDigestEncoding = errors.New("digest is not valid hex-encoded string")
	ErrRDRequiredField       = errors.New("at least one of name, URI, or digest are required")
)

// Supported standard hash algorithms
func isSupportedAlgorithm(alg string) (bool, int) {
	algos := map[string]int{"md5": md5.Size, "sha1": sha1.Size, "shake128": md5.Size, "sha224": sha512.Size224, "sha3_224": sha512.Size224, "sha512_224": sha512.Size224, "sha256": sha512.Size256, "sha3_256": sha512.Size256, "sha512_256": sha512.Size256, "shake256": sha512.Size256, "sha384": sha512.Size384, "sha3_384": sha512.Size384, "sha512_384": sha512.Size384, "sha512": sha512.Size, "sha3_512": sha512.Size, "dirHash": sha512.Size256, "gitCommit": sha1.Size}

	size, ok := algos[alg]
	return ok, size
}

func (d *ResourceDescriptor) Validate() error {
	// at least one of name, URI or digest are required
	if d.GetName() == "" && d.GetUri() == "" && len(d.GetDigest()) == 0 {
		return ErrRDRequiredField
	}

	if len(d.GetDigest()) > 0 {
		for alg, digest := range d.GetDigest() {

			// check encoding and length for supported algorithms
			supported, size := isSupportedAlgorithm(alg)
			if supported {
				// the in-toto spec expects a hex-encoded string in DigestSets for supported algorithms
				// https://github.com/in-toto/attestation/blob/main/spec/v1/digest_set.md
				hashBytes, err := hex.DecodeString(digest)

				if err != nil {
					return ErrInvalidDigestEncoding
				}

				// check the length of the digest
				if len(hashBytes) != size {
					return ErrIncorrectDigestLength
				}
			}
		}
	}

	return nil
}
