package main

import (
	"fmt"
	"log"
	"strings"

	vpb "github.com/in-toto/attestation/go/spec/predicates/vsa"
	spb "github.com/in-toto/attestation/go/spec/v1.0"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createStatementPbFromJson(subName string, subSha256 string, predicateType string, predicateJson []byte) (*spb.Statement, error) {
	pred := &structpb.Struct{}
	err := protojson.Unmarshal(predicateJson, pred)
	if err != nil {
		fmt.Errorf("failed to unmarshal predicate: %w", err)
		return nil, err
	}
	return createStatementPb(subName, subSha256, predicateType, pred), nil
}

func createStatementPb(subName string, subSha256 string, predicateType string, predicate *structpb.Struct) *spb.Statement {
	sub := []*spb.Statement_Subject{{
		Name:   subName,
		Digest: map[string]string{"sha256": strings.ToLower(subSha256)},
	}}
	statement := &spb.Statement{
		Type:          "https://in-toto.io/Statement/v1.0",
		Subject:       sub,
		PredicateType: predicateType,
		Predicate:     predicate,
	}
	return statement
}

func createVsa(subName string, subSha256 string, vsaBody *vpb.VerificationSummaryV02) (*spb.Statement, error) {
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

func createTestResourceDescriptor() (*spb.ResourceDescriptor, error) {
	// Create a ResourceDescriptor
	a1, err := structpb.NewStruct(map[string]interface{}{
		"keyStr": "value1",
		"keyNum": 13})
	if err != nil {
		return nil, err
	}
	a2, err := structpb.NewStruct(map[string]interface{}{
		"keyObj": map[string]interface{}{
			"subKey": "subVal"}})
	if err != nil {
		return nil, err
	}
	r := &spb.ResourceDescriptor{
		Name:             "theName",
		Uri:              "http://example.com",
		Digest:           map[string]string{"sha256": "abc123"},
		Content:          []byte("bytescontent"),
		DownloadLocation: "http://example.com/test.zip",
		MediaType:        "theMediaType",
		Annotations:      map[string]*structpb.Struct{"a1": a1, "a2": a2}}
	return r, nil
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
	vsaPred := &vpb.VerificationSummaryV02{
		Verifier: &vpb.VerificationSummaryV02_Verifier{
			Id: "verifier-id"},
		TimeVerified: timestamppb.Now(),
		ResourceUri:  "http://example.com/the/protected/resource.tar",
		Policy: &vpb.VerificationSummaryV02_Policy{
			Uri: "http://example.com/policy/uri"},
		InputAttestations: []*vpb.VerificationSummaryV02_InputAttestation{{
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
            "_type": "https://in-toto.io/Statement/v1.0",
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

	// Test ResourceDescriptor
	r, err := createTestResourceDescriptor()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nResourceDescriptor as json:\n%v\n", protojson.Format(r))
}
