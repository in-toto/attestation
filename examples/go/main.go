package main

import (
	"fmt"
	"log"
	"strings"

	vpb "github.com/in-toto/attestation/go/predicates/vsa/v0"
	spb "github.com/in-toto/attestation/go/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createStatementPbFromJson(subName string, subSha256 string, predicateType string, predicateJson []byte) (*spb.Statement, error) {
	pred := &structpb.Struct{}
	err := protojson.Unmarshal(predicateJson, pred)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal predicate: %w", err)
	}
	return createStatementPb(subName, subSha256, predicateType, pred), nil
}

func createStatementPb(subName string, subSha256 string, predicateType string, predicate *structpb.Struct) *spb.Statement {
	sub := []*spb.ResourceDescriptor{{
		Name:   subName,
		Digest: map[string]string{"sha256": strings.ToLower(subSha256)},
	}}
	return &spb.Statement{
		Type:          spb.StatementTypeUri,
		Subject:       sub,
		PredicateType: predicateType,
		Predicate:     predicate,
	}
}

func createVsa(subName string, subSha256 string, vsaBody *vpb.VerificationSummary) (*spb.Statement, error) {
	vsaJson, err := protojson.Marshal(vsaBody)
	if err != nil {
		return nil, err
	}
	vsaStruct := &structpb.Struct{}
	err = protojson.Unmarshal(vsaJson, vsaStruct)
	if err != nil {
		return nil, err
	}
	return createStatementPb(subName, subSha256, "https://slsa.dev/verification_summary/v0.2", vsaStruct), nil
}

// Example of how to use protobuf to create in-toto statements.
// Users will still likely want to put the json output in a DSSE.
func main() {
	// Create a statement with an unknown predicate.
	s, err := createStatementPbFromJson(
		"sub-name",
		"ABC123",
		"https://example.com/unknownPred1",
		[]byte(`{
                  "foo": "bar",
                  "baz": [1,2,3]
                }`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Statement as json:\n%v\n", protojson.Format(s))

	// Create a statement of a VSA
	vsaPred := &vpb.VerificationSummary{
		Verifier: &vpb.VerificationSummary_Verifier{
			Id: "verifier-id"},
		TimeVerified: timestamppb.Now(),
		ResourceUri:  "http://example.com/the/protected/resource.tar",
		Policy: &vpb.VerificationSummary_Policy{
			Uri: "http://example.com/policy/uri"},
		InputAttestations: []*vpb.VerificationSummary_InputAttestation{{
			Uri:    "http://example.com/attestation/foo.intoto.jsonl",
			Digest: map[string]string{"sha256": "def456"}},
		},
		VerificationResult: "PASSED",
		PolicyLevel:        "SLSA_LEVEL_3",
		DependencyLevels:   map[string]uint64{"SLSA_LEVEL_0": 1},
	}
	v, err := createVsa("vsa-sub", "abc123", vsaPred)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nVSA as json:\n%v\n", protojson.Format(v))

	// Read JSON text into a Statement
	err = protojson.Unmarshal([]byte(`{
            "_type": "https://in-toto.io/Statement/v1",
            "subject": [{
              "name": "sub name",
              "digest": { "sha256": "abc123" }
            }],
            "predicateType": "https://example.com/unknownPred2",
            "predicate": {
              "foo": {
                "bar": "baz"
              }
            }
          }`), s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nRead statement with predicateType %v\n", s.PredicateType)
	fmt.Printf("Predicate %v\n", s.Predicate)
}
