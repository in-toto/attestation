/*
Tests for in-toto attestation ResourceDescriptor protos.
*/

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

const wantSt = `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

func createTestStatement() (*Statement, error) {
	// Create a Statement

	sub := &ResourceDescriptor{
		Name:   "theSub",
		Digest: map[string]string{"alg1": "abc123"},
	}

	pred, err := structpb.NewStruct(map[string]interface{}{
		"keyObj": map[string]interface{}{
			"subKey": "subVal"}})
	if err != nil {
		return nil, err
	}

	return &Statement{
		Type:          StatementTypeUri,
		Subject:       []*ResourceDescriptor{sub},
		PredicateType: "thePredicate",
		Predicate:     pred,
	}, nil
}

func TestJsonUnmarshalStatement(t *testing.T) {
	got := &Statement{}
	err := protojson.Unmarshal([]byte(wantSt), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	want, err := createTestStatement()

	assert.NoError(t, err, "Error during test Statement creation")
	assert.True(t, proto.Equal(got, want), "Protos do not match")
}

func TestBadStatementType(t *testing.T) {
	var badStType = `{"_type":"https://in-toto.io/Statement/v0","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

	got := &Statement{}
	err := protojson.Unmarshal([]byte(badStType), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()

	assert.Error(t, err, "Error: created malformed Statement (bad type)")
}

func TestBadStatementSubject(t *testing.T) {
	var badStNoSub = `{"_type":"https://in-toto.io/Statement/v1","subject":[],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

	got := &Statement{}
	err := protojson.Unmarshal([]byte(badStNoSub), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()

	assert.Error(t, err, "Error: created malformed Statement (empty subject)")

	var badStBadSub = `{"_type":"https://in-toto.io/Statement/v1","subject":[{"downloadLocation":"https://example.com/test.zip"}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

	got = &Statement{}
	err = protojson.Unmarshal([]byte(badStBadSub), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()

	assert.Error(t, err, "Error: created malformed Statement (bad subject)")
}

func TestBadStatementPredicate(t *testing.T) {
	var badStPredType = `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"","predicate":{"keyObj":{"subKey":"subVal"}}}`

	got := &Statement{}
	err := protojson.Unmarshal([]byte(badStPredType), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()

	assert.Error(t, err, "Error: created malformed Statement (bad predicate type)")

	var badStPred = `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate"}`

	got = &Statement{}
	err = protojson.Unmarshal([]byte(badStPred), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()

	assert.Error(t, err, "Error: created malformed Statement (no prdicate)")
}
