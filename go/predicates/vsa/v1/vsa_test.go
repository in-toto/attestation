// Copyright 2026 in-toto Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"testing"
	"time"

	ita1 "github.com/in-toto/attestation/go/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	validSHA256    = "a1234567b1234567c1234567d1234567e1234567f1234567a1234567b1234567"
	validSHA256Alt = "b1234567b1234567c1234567d1234567e1234567f1234567a1234567b1234567"
)

func createTestVerificationSummary(t testing.TB) *VerificationSummary {
	t.Helper()

	ts, err := time.Parse(time.RFC3339Nano, "2026-11-16T22:00:00.52Z")
	if err != nil {
		t.Fatal(err)
	}

	return &VerificationSummary{
		Verifier: &VerificationSummary_Verifier{
			Id: "https://example.com/verifier/v1",
		},
		TimeVerified: timestamppb.New(ts),
		ResourceUri:  "oci://registry.example.com/prod/service@sha256:" + validSHA256,
		Policy: &VerificationSummary_Policy{
			Uri:    "https://example.com/policies/prod-admission/v3",
			Digest: map[string]string{"sha256": validSHA256},
		},
		InputAttestations: []*VerificationSummary_InputAttestation{
			{
				Uri:    "https://example.com/attestations/provenance.intoto.jsonl",
				Digest: map[string]string{"sha256": validSHA256Alt},
			},
		},
		VerificationResult: VerificationResultPassed,
		VerifiedLevels: []string{
			"SLSA_BUILD_LEVEL_3",
			"SLSA_SOURCE_LEVEL_3",
			"SLSA_BUILD_REPRODUCED",
			"SLSA_SOURCE_TWO_PARTY_REVIEWED",
		},
		DependencyLevels: map[string]uint64{
			"SLSA_BUILD_LEVEL_3": 5,
			"SLSA_BUILD_LEVEL_2": 2,
		},
		SlsaVersion: "1.2",
	}
}

func TestJSONRoundTripVerificationSummary(t *testing.T) {
	rawJSON := `{
		"verifier": {"id": "https://example.com/verifier/v1"},
		"timeVerified": "2026-11-16T22:00:00.520Z",
		"resourceUri": "oci://registry.example.com/prod/service@sha256:` + validSHA256 + `",
		"policy": {
			"uri": "https://example.com/policies/prod-admission/v3",
			"digest": {"sha256": "` + validSHA256 + `"}
		},
		"inputAttestations": [
			{
				"uri": "https://example.com/attestations/provenance.intoto.jsonl",
				"digest": {"sha256": "` + validSHA256Alt + `"}
			}
		],
		"verificationResult": "PASSED",
		"verifiedLevels": [
			"SLSA_BUILD_LEVEL_3",
			"SLSA_SOURCE_LEVEL_3",
			"SLSA_BUILD_REPRODUCED",
			"SLSA_SOURCE_TWO_PARTY_REVIEWED"
		],
		"dependencyLevels": {
			"SLSA_BUILD_LEVEL_3": 5,
			"SLSA_BUILD_LEVEL_2": 2
		},
		"slsaVersion": "1.2"
	}`

	got := &VerificationSummary{}
	require.NoError(t, protojson.Unmarshal([]byte(rawJSON), got))
	require.NoError(t, got.Validate())

	want := createTestVerificationSummary(t)
	assert.True(t, proto.Equal(got, want), "unmarshalled proto did not match expected VerificationSummary")

	marshalled, err := protojson.Marshal(got)
	require.NoError(t, err)

	roundTripped := &VerificationSummary{}
	require.NoError(t, protojson.Unmarshal(marshalled, roundTripped))
	require.NoError(t, roundTripped.Validate())
	assert.True(t, proto.Equal(want, roundTripped))
}

func TestValidVerificationSummaryVariants(t *testing.T) {
	tests := map[string]struct {
		unmarshalOpts protojson.UnmarshalOptions
		input         string
	}{
		"minimal PASSED without optional fields (SLSA v1.1+ optional timeVerified)": {
			input: `{
				"verifier": {"id": "https://example.com/verifier"},
				"resourceUri": "pkg:oci/app@sha256:` + validSHA256 + `",
				"policy": {"uri": "https://example.com/policy"},
				"verificationResult": "PASSED",
				"verifiedLevels": ["SLSA_BUILD_LEVEL_0"]
			}`,
		},
		"PASSED with SLSA_BUILD_LEVEL_UNEVALUATED and multi-word custom track": {
			input: `{
				"verifier": {"id": "https://example.com/verifier"},
				"resourceUri": "pkg:oci/app@sha256:` + validSHA256 + `",
				"policy": {"uri": "https://example.com/policy"},
				"verificationResult": "PASSED",
				"verifiedLevels": ["SLSA_BUILD_LEVEL_UNEVALUATED", "SLSA_CUSTOM_TRACK_LEVEL_12"]
			}`,
		},
		"FAILED with single FAILED entry": {
			input: `{
				"verifier": {"id": "https://example.com/verifier"},
				"resourceUri": "pkg:oci/app@sha256:` + validSHA256 + `",
				"policy": {"uri": "https://example.com/policy"},
				"verificationResult": "FAILED",
				"verifiedLevels": ["FAILED"]
			}`,
		},
		"FAILED with artifact track levels alongside FAILED (e.g. transitive dependency policy failure)": {
			input: `{
				"verifier": {"id": "https://example.com/verifier"},
				"resourceUri": "pkg:oci/app@sha256:` + validSHA256 + `",
				"policy": {"uri": "https://example.com/policy"},
				"verificationResult": "FAILED",
				"verifiedLevels": ["FAILED", "BCID_L3", "SLSA_BUILD_LEVEL_3"]
			}`,
		},
		"FAILED with empty verifiedLevels array": {
			input: `{
				"verifier": {"id": "https://example.com/verifier"},
				"resourceUri": "pkg:oci/app@sha256:` + validSHA256 + `",
				"policy": {"uri": "https://example.com/policy"},
				"verificationResult": "FAILED",
				"verifiedLevels": []
			}`,
		},
		"monotonic forward compatibility with future SLSA minor-version properties and RC slsaVersion": {
			unmarshalOpts: protojson.UnmarshalOptions{DiscardUnknown: true},
			input: `{
				"verifier": {
					"id": "https://k8s.io/promo-tools/verifier/v1",
					"version": {"slsa-verifier": "v2.6.0"}
				},
				"resourceUri": "registry.k8s.io/kube-apiserver:v1.34.0",
				"policy": {
					"uri": "https://github.com/kubernetes/k8s.io/policy",
					"digest": {"gitCommit": "a1234567b1234567c1234567d1234567e1234567"}
				},
				"verificationResult": "PASSED",
				"verifiedLevels": [
					"SLSA_BUILD_LEVEL_3",
					"SLSA_DEPENDENCY_LEVEL_PINNED",
					"K8S_PROMOTION_MANIFEST_REVIEWED"
				],
				"slsaVersion": "1.2-rc1"
			}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := &VerificationSummary{}
			require.NoError(t, tc.unmarshalOpts.Unmarshal([]byte(tc.input), got))
			assert.NoError(t, got.Validate())
		})
	}
}

func TestInvalidVerificationSummaryFromJSON(t *testing.T) {
	tests := map[string]struct {
		input   string
		wantErr error
	}{
		"missing verifier": {
			input:   `{"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ErrVerifierRequired,
		},
		"empty verifier object": {
			input:   `{"verifier":{},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ErrVerifierIDRequired,
		},
		"missing resourceUri": {
			input:   `{"verifier":{"id":"https://example.com/v"},"policy":{"uri":"https://example.com/p"},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ErrResourceURIRequired,
		},
		"missing policy": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ErrPolicyRequired,
		},
		"empty policy object": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ErrPolicyURIRequired,
		},
		"policy with digest but missing uri": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"digest":{"sha256":"` + validSHA256 + `"}},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ErrPolicyURIRequired,
		},
		"policy with invalid sha256 digest length": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p","digest":{"sha256":"abcd"}},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ita1.ErrIncorrectDigestLength,
		},
		"policy with invalid sha256 hex encoding": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p","digest":{"sha256":"zzzz"}},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ita1.ErrInvalidDigestEncoding,
		},
		"inputAttestation missing digest": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"inputAttestations":[{"uri":"https://example.com/att.jsonl"}],"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ErrInputAttestationDigestRequired,
		},
		"inputAttestation with empty custom digest value": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"inputAttestations":[{"uri":"https://example.com/att.jsonl","digest":{"custom":""}}],"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ita1.ErrEmptyDigest,
		},
		"inputAttestation with bad digest length": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"inputAttestations":[{"uri":"https://example.com/att.jsonl","digest":{"sha256":"1234"}}],"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ita1.ErrIncorrectDigestLength,
		},
		"missing verificationResult": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ErrVerificationResultRequired,
		},
		"invalid verificationResult enum": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"UNKNOWN","verifiedLevels":["SLSA_BUILD_LEVEL_1"]}`,
			wantErr: ErrInvalidVerificationResult,
		},
		"PASSED missing verifiedLevels": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"PASSED"}`,
			wantErr: ErrVerifiedLevelsRequired,
		},
		"PASSED with empty entry in verifiedLevels": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"PASSED","verifiedLevels":[""]}`,
			wantErr: ErrInvalidVerifiedLevels,
		},
		"FAILED with empty entry in verifiedLevels": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"FAILED","verifiedLevels":[""]}`,
			wantErr: ErrInvalidVerifiedLevels,
		},
		"PASSED summary containing FAILED in verifiedLevels": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"PASSED","verifiedLevels":["FAILED"]}`,
			wantErr: ErrInvalidVerifiedLevels,
		},
		"duplicate verifiedLevels entry on PASSED": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_3","SLSA_BUILD_LEVEL_3"]}`,
			wantErr: ErrInvalidVerifiedLevels,
		},
		"duplicate FAILED entry on FAILED": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"FAILED","verifiedLevels":["FAILED","FAILED"]}`,
			wantErr: ErrInvalidVerifiedLevels,
		},
		"multiple levels for same SLSA track on PASSED": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_2","SLSA_BUILD_LEVEL_3"]}`,
			wantErr: ErrDuplicateSLSATrack,
		},
		"multiple levels for same SLSA track on FAILED": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"FAILED","verifiedLevels":["FAILED","SLSA_BUILD_LEVEL_1","SLSA_BUILD_LEVEL_3"]}`,
			wantErr: ErrDuplicateSLSATrack,
		},
		"empty key in dependencyLevels": {
			input:   `{"verifier":{"id":"https://example.com/v"},"resourceUri":"https://example.com/a","policy":{"uri":"https://example.com/p"},"verificationResult":"PASSED","verifiedLevels":["SLSA_BUILD_LEVEL_3"],"dependencyLevels":{"":1}}`,
			wantErr: ErrInvalidDependencyLevels,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := &VerificationSummary{}
			require.NoError(t, protojson.Unmarshal([]byte(tc.input), got))
			assert.ErrorIs(t, got.Validate(), tc.wantErr)
		})
	}
}

func TestSubMessageAndDirectStructEdgeCases(t *testing.T) {
	assert.ErrorIs(t, (*VerificationSummary)(nil).Validate(), ErrVerifierRequired)
	assert.ErrorIs(t, (*VerificationSummary_Verifier)(nil).Validate(), ErrVerifierRequired)
	assert.ErrorIs(t, (&VerificationSummary_Verifier{}).Validate(), ErrVerifierIDRequired)
	assert.ErrorIs(t, (*VerificationSummary_Policy)(nil).Validate(), ErrPolicyRequired)
	assert.ErrorIs(t, (&VerificationSummary_Policy{}).Validate(), ErrPolicyURIRequired)
	assert.ErrorIs(t, (*VerificationSummary_InputAttestation)(nil).Validate(), ErrInputAttestationDigestRequired)

	badTime := createTestVerificationSummary(t)
	badTime.TimeVerified = &timestamppb.Timestamp{Seconds: 100, Nanos: -5}
	assert.ErrorIs(t, badTime.Validate(), ErrInvalidTimeVerified)

	nilInput := createTestVerificationSummary(t)
	nilInput.InputAttestations = append(nilInput.InputAttestations, nil)
	assert.ErrorIs(t, nilInput.Validate(), ErrInputAttestationDigestRequired)
}

func TestSLSATrackFromLevel(t *testing.T) {
	tests := []struct {
		input     string
		wantTrack string
		wantOK    bool
	}{
		{"SLSA_BUILD_LEVEL_0", "BUILD", true},
		{"SLSA_BUILD_LEVEL_3", "BUILD", true},
		{"SLSA_SOURCE_LEVEL_UNEVALUATED", "SOURCE", true},
		{"SLSA_CUSTOM_TRACK_2_LEVEL_12", "CUSTOM_TRACK_2", true},
		{"SLSA_BUILD_REPRODUCED", "", false},
		{"SLSA_SOURCE_TWO_PARTY_REVIEWED", "", false},
		{"SLSA_DEPENDENCY_LEVEL_PINNED", "", false},
		{"SLSA_BUILD_LEVEL_003", "", false},
		{"SLSA_BUILD_LEVEL_1A", "", false},
		{"SLSA_BUILD_LEVEL_", "", false},
		{"SLSA__LEVEL_3", "", false},
		{"SLSA__BUILD_LEVEL_3", "", false},
		{"SLSA_BUILD__LEVEL_3", "", false},
		{"SLSA_build_LEVEL_3", "", false},
		{"SLSA_LEVEL_3", "", false},
		{"SLSA_", "", false},
		{"FAILED", "", false},
		{"BCID_L3", "", false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			gotTrack, gotOK := slsaTrackFromLevel(tc.input)
			assert.Equal(t, tc.wantOK, gotOK)
			assert.Equal(t, tc.wantTrack, gotTrack)
		})
	}
}

func BenchmarkVerificationSummaryValidate(b *testing.B) {
	vsa := createTestVerificationSummary(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := vsa.Validate(); err != nil {
			b.Fatal(err)
		}
	}
}
