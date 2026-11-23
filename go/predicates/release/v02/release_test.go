/*
Tests for release v0.2 protos.
*/

package v02

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestValidRelease(t *testing.T) {
	// The purls from the examples in the release predicate spec.
	for _, purl := range []string{
		"pkg:npm/@angular/http@7.2.16",
		"pkg:npm/%40angular/http@7.2.16",
		"pkg:pypi/urllib3@2.1.0",
		"pkg:oci/my-image@sha256%3Afedcba09?repository_url=registry.example.com/my-project/my-image&tag=v1.2.3",
		"pkg:npm/semver@7.5.4",
		"pkg:maven/org.apache.commons/io@1.3.4#some/subpath",
	} {
		r := &Release{Purl: purl}
		assert.NoError(t, r.Validate(), purl)
	}

	got := &Release{}
	err := protojson.Unmarshal([]byte(`{"purl": "pkg:npm/@angular/http@7.2.16", "packageId": "1234567890"}`), got)
	assert.NoError(t, err)
	assert.NoError(t, got.Validate())
}

func TestBadRelease(t *testing.T) {
	tests := map[string]struct {
		purl string
		err  error
	}{
		"no purl":                     {"", ErrPurlRequired},
		"not a purl":                  {"https://example.com/urllib3-2.1.0.tar.gz", ErrInvalidPurl},
		"no type":                     {"pkg:/urllib3@2.1.0", ErrInvalidPurl},
		"no name":                     {"pkg:pypi/@2.1.0", ErrInvalidPurl},
		"no version":                  {"pkg:pypi/urllib3", ErrPurlVersionRequired},
		"empty version":               {"pkg:pypi/urllib3@", ErrPurlVersionRequired},
		"version only in qualifiers":  {"pkg:oci/my-image?tag=v1.2.3", ErrPurlVersionRequired},
		"at sign only in namespace":   {"pkg:npm/@angular/http", ErrPurlVersionRequired},
		"version only in the subpath": {"pkg:pypi/urllib3#dist/urllib3@2.1.0", ErrPurlVersionRequired},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &Release{Purl: tt.purl}
			assert.ErrorIs(t, r.Validate(), tt.err)
		})
	}
}
