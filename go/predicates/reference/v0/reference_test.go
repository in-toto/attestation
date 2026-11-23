/*
Tests for reference v0 protos.
*/

package v0

import (
	"testing"

	ita1 "github.com/in-toto/attestation/go/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

const validReference = `{
  "attester": {
    "id": "https://my-builder.com/v1"
  },
  "references": [{
    "downloadLocation": "https://cloud.com/my-sboms/sbom.spdx.json",
    "digest": {
      "sha256": "3b0d19b348f1e46a571c6e9df32897637c20dd9803261c2fc9cbe38b0c8422c4"
    },
    "mediaType": "application/spdx+json"
  }]
}`

func parseReference(t *testing.T, s string) *Reference {
	t.Helper()
	r := &Reference{}
	if err := protojson.Unmarshal([]byte(s), r); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestValidReference(t *testing.T) {
	r := parseReference(t, validReference)
	assert.NoError(t, r.Validate())
}

func TestValidReferenceMultiple(t *testing.T) {
	r := parseReference(t, validReference)
	r.References = append(r.References, &ita1.ResourceDescriptor{
		Uri:              "https://cloud.com/my-sboms/sbom.cdx.json",
		DownloadLocation: "https://cloud.com/my-sboms/sbom.cdx.json",
		MediaType:        "application/octet-stream",
	})
	assert.NoError(t, r.Validate())
}

func TestReferenceJSONRoundTrip(t *testing.T) {
	r := parseReference(t, validReference)
	b, err := protojson.Marshal(r)
	assert.NoError(t, err)

	got := parseReference(t, string(b))
	assert.NoError(t, got.Validate())
	assert.Equal(t, r.GetAttester().GetId(), got.GetAttester().GetId())
	assert.Equal(t, r.GetReferences()[0].GetDigest(), got.GetReferences()[0].GetDigest())
}

func TestBadReference(t *testing.T) {
	tests := map[string]struct {
		mutate  func(*Reference)
		wantErr error
		wantMsg string
	}{
		"no attester": {
			mutate:  func(r *Reference) { r.Attester = nil },
			wantErr: ErrAttesterRequired,
		},
		"empty attester id": {
			mutate:  func(r *Reference) { r.Attester.Id = "" },
			wantErr: ErrAttesterIdRequired,
		},
		"no references": {
			mutate:  func(r *Reference) { r.References = nil },
			wantErr: ErrReferencesRequired,
		},
		"nil reference": {
			mutate: func(r *Reference) {
				r.References = append(r.References, nil)
			},
			wantErr: ErrReferenceRequired,
			wantMsg: "references[1]: reference must not be empty",
		},
		"no downloadLocation": {
			mutate:  func(r *Reference) { r.References[0].DownloadLocation = "" },
			wantErr: ErrDownloadLocationRequired,
			wantMsg: "references[0]: downloadLocation required",
		},
		"no mediaType": {
			mutate:  func(r *Reference) { r.References[0].MediaType = "" },
			wantErr: ErrMediaTypeRequired,
			wantMsg: "references[0]: mediaType required",
		},
		"no name, uri or digest": {
			mutate:  func(r *Reference) { r.References[0].Digest = nil },
			wantErr: ita1.ErrRDRequiredField,
		},
		"bad digest": {
			mutate:  func(r *Reference) { r.References[0].Digest["sha256"] = "abc" },
			wantErr: ita1.ErrInvalidDigestEncoding,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := parseReference(t, validReference)
			tt.mutate(r)
			err := r.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantMsg != "" {
				assert.EqualError(t, err, tt.wantMsg)
			}
		})
	}
}
