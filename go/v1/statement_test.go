/*
Tests for in-toto attestation ResourceDescriptor protos.
*/

package v1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func createTestStatement(t *testing.T) *Statement {
	// Create a Statement

	t.Helper()

	sub := &ResourceDescriptor{
		Name:   "theSub",
		Digest: map[string]string{"alg1": "abc123"},
	}

	pred, err := structpb.NewStruct(map[string]interface{}{
		"keyObj": map[string]interface{}{
			"subKey": "subVal"}})
	if err != nil {
		t.Fatal(err)
	}

	return &Statement{
		Type:          StatementTypeUri,
		Subject:       []*ResourceDescriptor{sub},
		PredicateType: "thePredicate",
		Predicate:     pred,
	}
}

func TestProtojsonUnmarshalStatement(t *testing.T) {
	var wantSt = `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

	got := &Statement{}
	err := protojson.Unmarshal([]byte(wantSt), got)
	assert.NoError(t, err, "error during JSON unmarshalling")

	want := createTestStatement(t)
	assert.True(t, proto.Equal(got, want), "protos do not match")
}

func TestJsonUnmarshalStatement(t *testing.T) {
	var wantSt = `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

	got := &Statement{}
	err := json.Unmarshal([]byte(wantSt), got)
	assert.NoError(t, err, "error during JSON unmarshalling")

	want := createTestStatement(t)
	assert.True(t, proto.Equal(got, want), "protos do not match")
}

func TestProtojsonMarshalStatement(t *testing.T) {
	var wantSt = `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`
	want := createTestStatement(t)

	gotSt, err := protojson.Marshal(want)
	assert.NoError(t, err, "error during JSON marshalling")
	assert.JSONEq(t, wantSt, string(gotSt), "JSON objects do not match")
}

func TestJsonMarshalStatement(t *testing.T) {
	var wantSt = `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`
	want := createTestStatement(t)

	gotSt, err := json.Marshal(want)
	assert.NoError(t, err, "error during JSON marshalling")
	assert.JSONEq(t, wantSt, string(gotSt), "JSON objects do not match")
}

func TestBadStatementType(t *testing.T) {
	var badStType = `{"_type":"https://in-toto.io/Statement/v0","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

	got := &Statement{}
	err := protojson.Unmarshal([]byte(badStType), got)
	assert.NoError(t, err, "error during JSON unmarshalling")

	err = got.Validate()
	assert.ErrorIs(t, err, ErrInvalidStatementType, "created malformed Statement (bad type)")
}

func TestBadStatementSubject(t *testing.T) {
	tests := map[string]struct {
		input        string
		err          error
		noErrMessage string
	}{
		"no subjects": {
			input:        `{"_type":"https://in-toto.io/Statement/v1","subject":[],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`,
			err:          ErrSubjectRequired,
			noErrMessage: "created malformed Statement (empty subject)",
		},
		"subject missing RD required field": {
			input:        `{"_type":"https://in-toto.io/Statement/v1","subject":[{"downloadLocation":"https://example.com/test.zip"}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`,
			err:          ErrRDRequiredField,
			noErrMessage: "created malformed Statement (subject missing one of RD required fields, 'name', 'uri', 'digest')",
		},
		"subject missing Statement requirement on digest": {
			input:        `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name": "theSub", "downloadLocation":"https://example.com/test.zip"}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`,
			err:          ErrDigestRequired,
			noErrMessage: "created malformed Statement (subject missing 'digest')",
		},
	}

	for name, test := range tests {
		got := &Statement{}
		err := protojson.Unmarshal([]byte(test.input), got)
		assert.NoError(t, err, fmt.Sprintf("error during JSON unmarshalling in test '%s'", name))

		err = got.Validate()
		assert.ErrorIs(t, err, test.err, fmt.Sprintf("%s in test '%s'", test.noErrMessage, name))
	}
}

func TestBadStatementPredicate(t *testing.T) {
	tests := map[string]struct {
		input        string
		err          error
		noErrMessage string
	}{
		"missing predicate type": {
			input:        `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"","predicate":{"keyObj":{"subKey":"subVal"}}}`,
			err:          ErrPredicateTypeRequired,
			noErrMessage: "created malformed Statement (missing predicate type)",
		},
		"missing predicate": {
			input:        `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate"}`,
			err:          ErrPredicateRequired,
			noErrMessage: "created malformed Statement (no predicate)",
		},
	}

	for name, test := range tests {
		got := &Statement{}
		err := protojson.Unmarshal([]byte(test.input), got)
		assert.NoError(t, err, fmt.Sprintf("error during JSON unmarshalling in test '%s'", name))

		err = got.Validate()
		assert.ErrorIs(t, err, test.err, fmt.Sprintf("%s in test '%s'", test.noErrMessage, name))
	}
}
