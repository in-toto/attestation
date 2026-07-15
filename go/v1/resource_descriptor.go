/*
Wrapper APIs for in-toto attestation ResourceDescriptor protos.
*/

package v1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrIncorrectDigestLength = errors.New("digest has incorrect length")
	ErrInvalidDigestEncoding = errors.New("digest is not valid hex-encoded string")
	ErrRDRequiredField       = errors.New("at least one of name, URI, or digest are required")
)

type HashAlgorithm string

const (
	AlgorithmMD5        HashAlgorithm = "md5"
	AlgorithmSHA1       HashAlgorithm = "sha1"
	AlgorithmSHA224     HashAlgorithm = "sha224"
	AlgorithmSHA512_224 HashAlgorithm = "sha512_224"
	AlgorithmSHA256     HashAlgorithm = "sha256"
	AlgorithmSHA512_256 HashAlgorithm = "sha512_256"
	AlgorithmSHA384     HashAlgorithm = "sha384"
	AlgorithmSHA512     HashAlgorithm = "sha512"
	AlgorithmSHA3_224   HashAlgorithm = "sha3_224"
	AlgorithmSHA3_256   HashAlgorithm = "sha3_256"
	AlgorithmSHA3_384   HashAlgorithm = "sha3_384"
	AlgorithmSHA3_512   HashAlgorithm = "sha3_512"
	AlgorithmGitBlob    HashAlgorithm = "gitBlob"
	AlgorithmGitCommit  HashAlgorithm = "gitCommit"
	AlgorithmGitTag     HashAlgorithm = "gitTag"
	AlgorithmGitTree    HashAlgorithm = "gitTree"
	AlgorithmDirHash    HashAlgorithm = "dirHash"
)

// HashAlgorithms indexes the known algorithms in a dictionary
// by their string value
var HashAlgorithms = map[string]HashAlgorithm{
	"md5":        AlgorithmMD5,
	"sha1":       AlgorithmSHA1,
	"sha224":     AlgorithmSHA224,
	"sha512_224": AlgorithmSHA512_224,
	"sha256":     AlgorithmSHA256,
	"sha512_256": AlgorithmSHA512_256,
	"sha384":     AlgorithmSHA384,
	"sha512":     AlgorithmSHA512,
	"sha3_224":   AlgorithmSHA3_224,
	"sha3_256":   AlgorithmSHA3_256,
	"sha3_384":   AlgorithmSHA3_384,
	"sha3_512":   AlgorithmSHA3_512,
	"gitBlob":    AlgorithmGitBlob,
	"gitCommit":  AlgorithmGitCommit,
	"gitTag":     AlgorithmGitTag,
	"gitTree":    AlgorithmGitTree,
	"dirHash":    AlgorithmDirHash,
}

// HexLength returns the expected length of an algorithm's hash when hexencoded
func (algo HashAlgorithm) HexLength() int {
	switch algo {
	case AlgorithmMD5:
		return 16
	case AlgorithmSHA1, AlgorithmGitBlob, AlgorithmGitCommit, AlgorithmGitTag, AlgorithmGitTree:
		return 20
	case AlgorithmSHA224, AlgorithmSHA512_224, AlgorithmSHA3_224:
		return 28
	case AlgorithmSHA256, AlgorithmSHA512_256, AlgorithmSHA3_256, AlgorithmDirHash:
		return 32
	case AlgorithmSHA384, AlgorithmSHA3_384:
		return 48
	case AlgorithmSHA512, AlgorithmSHA3_512:
		return 64
	default:
		return 0
	}
}

// String returns the hash algorithm name as a string
func (algo HashAlgorithm) String() string {
	return string(algo)
}

// digestLengths returns every digest size in bytes that an algorithm accepts.
//
// Most algorithms have exactly one. The git object algorithms have two: git
// names objects with either a SHA-1 or a SHA-256 hash, and digest_set.md
// accepts both, telling them apart by length.
func (algo HashAlgorithm) digestLengths() []int {
	switch algo {
	case AlgorithmGitBlob, AlgorithmGitCommit, AlgorithmGitTag, AlgorithmGitTree:
		return []int{20, 32}
	}

	if size := algo.HexLength(); size > 0 {
		return []int{size}
	}

	return nil
}

// Indicates if a given hash algorithm is supported by default and returns the
// digest sizes in bytes that it accepts, if supported.
//
// SHA digest sizes from https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.202.pdf
// MD5 digest size from https://www.rfc-editor.org/rfc/rfc1321.html#section-1
func isSupportedAlgorithm(algString string) (bool, []int) {
	sizes := HashAlgorithm(algString).digestLengths()
	return len(sizes) > 0, sizes
}

// formatSizes renders the accepted digest sizes for an error message, e.g.
// "20" or "20 or 32".
func formatSizes(sizes []int) string {
	strs := make([]string, 0, len(sizes))
	for _, size := range sizes {
		strs = append(strs, strconv.Itoa(size))
	}

	return strings.Join(strs, " or ")
}

func (d *ResourceDescriptor) Validate() error {
	// at least one of name, URI or digest are required
	if d.GetName() == "" && d.GetUri() == "" && len(d.GetDigest()) == 0 {
		return ErrRDRequiredField
	}

	if len(d.GetDigest()) > 0 {
		for alg, digest := range d.GetDigest() {

			// Per https://github.com/in-toto/attestation/blob/main/spec/v1/digest_set.md
			// check encoding and length for supported algorithms;
			// use of custom, unsupported algorithms is allowed and does not not generate validation errors.
			supported, sizes := isSupportedAlgorithm(alg)
			if supported {
				// the in-toto spec expects a hex-encoded string in DigestSets for supported algorithms
				hashBytes, err := hex.DecodeString(digest)

				if err != nil {
					return fmt.Errorf("%w (%s: %s)", ErrInvalidDigestEncoding, alg, digest)
				}

				// check the length of the digest
				if !slices.Contains(sizes, len(hashBytes)) {
					return fmt.Errorf("%w: got %d bytes, want %s bytes (%s: %s)", ErrIncorrectDigestLength, len(hashBytes), formatSizes(sizes), alg, digest)
				}
			}
		}
	}

	return nil
}
