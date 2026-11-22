/*
Tests for test-result v0.1 protos.
*/

package v0

import (
	"testing"

	ita1 "github.com/in-toto/attestation/go/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

const validTestResult = `{
  "result": "PASSED",
  "configuration": [{
    "name": ".github/workflows/ci.yml",
    "downloadLocation": "https://github.com/in-toto/in-toto/blob/d20ace7968ba43c0219f62d71334c1095bab1602/.github/workflows/ci.yml",
    "digest": {
      "gitBlob": "ebe4add40f63c3c98bc9b32ff1e736f04120b023"
    }
  }],
  "url": "https://github.com/in-toto/in-toto/actions/runs/4425592351",
  "passedTests": ["build (3.10, ubuntu-latest, py)"],
  "warnedTests": [],
  "failedTests": []
}`

func TestValidTestResult(t *testing.T) {
	got := &TestResult{}
	err := protojson.Unmarshal([]byte(validTestResult), got)
	assert.NoError(t, err)
	assert.NoError(t, got.Validate())

	for _, r := range Results {
		got.Result = r
		assert.NoError(t, got.Validate(), r)
	}
}

func TestBadTestResult(t *testing.T) {
	tests := map[string]struct {
		edit func(*TestResult)
		err  error
	}{
		"no result": {
			edit: func(r *TestResult) { r.Result = "" },
			err:  ErrResultRequired,
		},
		"unknown result": {
			edit: func(r *TestResult) { r.Result = "SKIPPED" },
			err:  ErrInvalidResult,
		},
		"lowercase result": {
			edit: func(r *TestResult) { r.Result = "passed" },
			err:  ErrInvalidResult,
		},
		"no configuration": {
			edit: func(r *TestResult) { r.Configuration = nil },
			err:  ErrConfigurationRequired,
		},
		"empty configuration entry": {
			edit: func(r *TestResult) { r.Configuration = append(r.Configuration, &ita1.ResourceDescriptor{}) },
			err:  ita1.ErrRDRequiredField,
		},
		"nil configuration entry": {
			edit: func(r *TestResult) { r.Configuration = append(r.Configuration, nil) },
			err:  ita1.ErrRDRequiredField,
		},
		"bad configuration digest": {
			edit: func(r *TestResult) {
				r.Configuration[0].Digest = map[string]string{"sha256": "not-hex"}
			},
			err: ita1.ErrInvalidDigestEncoding,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := &TestResult{}
			if err := protojson.Unmarshal([]byte(validTestResult), got); err != nil {
				t.Fatal(err)
			}
			tt.edit(got)
			err := got.Validate()
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestBadTestResultPath(t *testing.T) {
	got := &TestResult{}
	if err := protojson.Unmarshal([]byte(validTestResult), got); err != nil {
		t.Fatal(err)
	}
	got.Configuration = append(got.Configuration, &ita1.ResourceDescriptor{})
	assert.ErrorContains(t, got.Validate(), "configuration[1]")
}
