/*
Wrapper APIs for in-toto attestation ResourceDescriptor protos.
*/

package v1

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrIncorrectDigestLength = errors.New("digest has incorrect length")
	ErrInvalidDigestEncoding = errors.New("digest is not valid hex-encoded string")
	ErrEmptyDigest           = errors.New("digest algorithm or value cannot be empty")
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

// HexLength returns the expected length of an algorithm's hash in bytes.
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

func (algo HashAlgorithm) isGitObject() bool {
	switch algo {
	case AlgorithmGitBlob, AlgorithmGitCommit, AlgorithmGitTag, AlgorithmGitTree:
		return true
	default:
		return false
	}
}

// digestLengths returns every digest size in bytes that an algorithm accepts.
//
// Most algorithms have exactly one. The git object algorithms have two: git
// names objects with either a SHA-1 or a SHA-256 hash, and digest_set.md
// accepts both, telling them apart by length.
func (algo HashAlgorithm) digestLengths() []int {
	if algo.isGitObject() {
		return []int{20, 32}
	}

	if size := algo.HexLength(); size > 0 {
		return []int{size}
	}

	return nil
}

func (algo HashAlgorithm) validByteLength(byteLen int) bool {
	if algo.isGitObject() {
		return byteLen == 20 || byteLen == 32
	}
	want := algo.HexLength()
	return want > 0 && byteLen == want
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

// isValidHex reports whether s is an even-length hexadecimal string without
// allocating a decoded byte slice (accepting upper- and lower-case hex to
// match encoding/hex.DecodeString behavior).
func isValidHex(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// ValidateDigestSet checks that digestSet contains at least one entry and that
// all entries using supported default hash algorithms contain valid hex-encoded
// digests of an accepted byte length per
// https://github.com/in-toto/attestation/blob/main/spec/v1/digest_set.md.
// Custom, unrecognized algorithms are permitted provided both the algorithm
// name and digest value are non-empty.
func ValidateDigestSet(digestSet map[string]string) error {
	if len(digestSet) == 0 {
		return ErrDigestRequired
	}

	for alg, digest := range digestSet {
		if alg == "" {
			return ErrEmptyDigest
		}

		algo := HashAlgorithm(alg)
		if algo.HexLength() == 0 {
			if digest == "" {
				return fmt.Errorf("%w (%s)", ErrEmptyDigest, alg)
			}
			continue
		}

		if !isValidHex(digest) {
			return fmt.Errorf("%w (%s: %s)", ErrInvalidDigestEncoding, alg, digest)
		}

		byteLen := len(digest) / 2
		if !algo.validByteLength(byteLen) {
			_, sizes := isSupportedAlgorithm(alg)
			return fmt.Errorf("%w: got %d bytes, want %s bytes (%s: %s)", ErrIncorrectDigestLength, byteLen, formatSizes(sizes), alg, digest)
		}
	}

	return nil
}

func (d *ResourceDescriptor) Validate() error {
	// at least one of name, URI or digest are required
	if d.GetName() == "" && d.GetUri() == "" && len(d.GetDigest()) == 0 {
		return ErrRDRequiredField
	}

	if len(d.GetDigest()) > 0 {
		if err := ValidateDigestSet(d.GetDigest()); err != nil {
			return err
		}
	}

	return nil
}
