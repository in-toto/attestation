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

// Package v1 provides validation for SLSA Verification Summary Attestation
// (VSA) v1 protos (https://slsa.dev/verification_summary/v1).
//
// Callers marshalling or unmarshalling JSON VSA predicates should use
// google.golang.org/protobuf/encoding/protojson (with
// protojson.UnmarshalOptions{DiscardUnknown: true} to honor the in-toto
// unrecognized-fields parsing rule) rather than encoding/json.
package v1

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	ita1 "github.com/in-toto/attestation/go/v1"
)

const (
	// PredicateTypeURI is the canonical in-toto predicateType URI for SLSA VSA v1.
	PredicateTypeURI = "https://slsa.dev/verification_summary/v1"

	// VerificationResultPassed indicates the artifact satisfied policy verification.
	VerificationResultPassed = "PASSED"
	// VerificationResultFailed indicates the artifact failed policy verification.
	VerificationResultFailed = "FAILED"

	slsaPrefix  = "SLSA_"
	levelMarker = "_LEVEL_"
	levelUneval = "UNEVALUATED"
)

var (
	ErrVerifierRequired               = errors.New("verifier required")
	ErrVerifierIDRequired             = errors.New("verifier.id required")
	ErrInvalidTimeVerified            = errors.New("timeVerified must be a valid timestamp")
	ErrResourceURIRequired            = errors.New("resourceUri required")
	ErrPolicyRequired                 = errors.New("policy required")
	ErrPolicyURIRequired              = errors.New("policy.uri required")
	ErrInputAttestationDigestRequired = errors.New("inputAttestation.digest required")
	ErrVerificationResultRequired     = errors.New("verificationResult required")
	ErrInvalidVerificationResult      = errors.New("verificationResult must be PASSED or FAILED")
	ErrVerifiedLevelsRequired         = errors.New("verifiedLevels required")
	ErrInvalidVerifiedLevels          = errors.New("invalid verifiedLevels")
	ErrDuplicateSLSATrack             = errors.New("verifiedLevels must not contain multiple levels for the same SLSA track")
	ErrInvalidDependencyLevels        = errors.New("invalid dependencyLevels")
)

// Validate checks that the Verifier sub-message is present and specifies a non-empty id.
func (v *VerificationSummary_Verifier) Validate() error {
	if v == nil {
		return ErrVerifierRequired
	}
	if v.GetId() == "" {
		return ErrVerifierIDRequired
	}
	return nil
}

// Validate checks that the Policy sub-message is present, specifies a non-empty uri,
// and contains valid DigestSet entries whenever digest is populated.
func (p *VerificationSummary_Policy) Validate() error {
	if p == nil {
		return ErrPolicyRequired
	}
	if p.GetUri() == "" {
		return ErrPolicyURIRequired
	}
	if len(p.GetDigest()) > 0 {
		if err := ita1.ValidateDigestSet(p.GetDigest()); err != nil {
			return fmt.Errorf("policy.digest: %w", err)
		}
	}
	return nil
}

// Validate checks that an InputAttestation sub-message contains at least one digest
// and that all supported digest algorithms use valid hex encodings and byte lengths.
func (ia *VerificationSummary_InputAttestation) Validate() error {
	if len(ia.GetDigest()) == 0 {
		return ErrInputAttestationDigestRequired
	}
	if err := ita1.ValidateDigestSet(ia.GetDigest()); err != nil {
		return fmt.Errorf("inputAttestation.digest: %w", err)
	}
	return nil
}

// Validate checks structural conformance of a SLSA VerificationSummary v1 predicate
// against https://slsa.dev/verification_summary/v1.
func (s *VerificationSummary) Validate() error {
	if err := s.GetVerifier().Validate(); err != nil {
		return err
	}

	if tv := s.GetTimeVerified(); tv != nil {
		if err := tv.CheckValid(); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidTimeVerified, err)
		}
	}

	if s.GetResourceUri() == "" {
		return ErrResourceURIRequired
	}

	if err := s.GetPolicy().Validate(); err != nil {
		return err
	}

	for i, ia := range s.GetInputAttestations() {
		if err := ia.Validate(); err != nil {
			return fmt.Errorf("inputAttestations[%d]: %w", i, err)
		}
	}

	res := s.GetVerificationResult()
	if res == "" {
		return ErrVerificationResultRequired
	}
	if res != VerificationResultPassed && res != VerificationResultFailed {
		return fmt.Errorf("%w: got %q", ErrInvalidVerificationResult, res)
	}

	if err := validateVerifiedLevels(res, s.GetVerifiedLevels()); err != nil {
		return err
	}

	for k := range s.GetDependencyLevels() {
		if k == "" {
			return fmt.Errorf("%w: empty key", ErrInvalidDependencyLevels)
		}
	}

	return nil
}

func validateVerifiedLevels(result string, levels []string) error {
	if len(levels) == 0 {
		if result == VerificationResultPassed {
			return ErrVerifiedLevelsRequired
		}
		// In proto3, an explicit empty JSON array ("verifiedLevels": []) is
		// indistinguishable from an unset field; allow empty on FAILED summaries.
		return nil
	}

	for i, lvl := range levels {
		if lvl == "" {
			return fmt.Errorf("%w: empty entry at index %d", ErrInvalidVerifiedLevels, i)
		}
		if result == VerificationResultPassed && lvl == VerificationResultFailed {
			return fmt.Errorf("%w: PASSED summary cannot contain %q", ErrInvalidVerifiedLevels, VerificationResultFailed)
		}
		if slices.Contains(levels[:i], lvl) {
			return fmt.Errorf("%w: duplicate entry %q", ErrInvalidVerifiedLevels, lvl)
		}

		track, ok := slsaTrackFromLevel(lvl)
		if !ok {
			continue
		}
		for _, prev := range levels[:i] {
			if prevTrack, prevOK := slsaTrackFromLevel(prev); prevOK && prevTrack == track {
				return fmt.Errorf("%w: %q", ErrDuplicateSLSATrack, track)
			}
		}
	}

	return nil
}

// slsaTrackFromLevel extracts the track name from a SlsaResult string when val
// conforms to the SLSA track-level pattern "SLSA_<TRACK>_LEVEL_<N|UNEVALUATED>".
// Non-matching strings (such as SLSA Verified Properties, custom domain
// properties, or future minor-version additions) return ("", false) to preserve
// monotonic forward compatibility.
func slsaTrackFromLevel(val string) (string, bool) {
	rest, ok := strings.CutPrefix(val, slsaPrefix)
	if !ok {
		return "", false
	}

	idx := strings.LastIndex(rest, levelMarker)
	if idx <= 0 {
		return "", false
	}

	track := rest[:idx]
	suffix := rest[idx+len(levelMarker):]
	if !isTrackIdentifier(track) || !isLevelSuffix(suffix) {
		return "", false
	}
	return track, true
}

func isTrackIdentifier(s string) bool {
	if len(s) == 0 || s[0] == '_' || s[len(s)-1] == '_' {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

func isLevelSuffix(s string) bool {
	if s == levelUneval || s == "0" {
		return true
	}
	if len(s) == 0 || s[0] < '1' || s[0] > '9' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
