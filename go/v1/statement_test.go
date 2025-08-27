/*
Tests for in-toto attestation ResourceDescriptor protos.
*/

package v1

import (
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

func TestJsonUnmarshalStatement(t *testing.T) {
	var wantSt = `{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

	got := &Statement{}
	err := protojson.Unmarshal([]byte(wantSt), got)
	assert.NoError(t, err, "error during JSON unmarshalling")

	want := createTestStatement(t)
	assert.NoError(t, err, "unexpected error during test Statement creation")
	assert.True(t, proto.Equal(got, want), "protos do not match")
}

func TestBadStatementType(t *testing.T) {
	var badStType = `{"_type":"https://not-in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

	got := &Statement{}
	err := protojson.Unmarshal([]byte(badStType), got)
	assert.NoError(t, err, "error during JSON unmarshalling")

	err = got.Validate()
	assert.ErrorIs(t, err, ErrInvalidStatementType, "created malformed Statement (bad type)")
}

func TestStatementTypeVersions(t *testing.T) {
	var tests = []struct {
		version string
		err     error
	}{
		{
			version: "1",
			err:     nil,
		},
		{
			version: "0.1",
			err:     nil,
		},
		{
			version: "2",
			err:     ErrInvalidStatementType,
		},
		{
			version: "0.2",
			err:     ErrInvalidStatementType,
		},
	}

	stStringTemplate := `{"_type":"https://in-toto.io/Statement/v%s","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}`

	for _, tt := range tests {
		stString := fmt.Sprintf(stStringTemplate, tt.version)
		got := &Statement{}
		err := protojson.Unmarshal([]byte(stString), got)
		assert.NoError(t, err, "error during JSON unmarshalling")

		err = got.Validate()
		assert.Equal(t, tt.err, err)
	}
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
