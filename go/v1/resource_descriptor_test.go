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

const wantFullRd = `{"name":"theName","uri":"https://example.com","digest":{"alg1":"abc123"},"content":"Ynl0ZXNjb250ZW50","downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType","annotations":{"a1":{"keyNum": 13,"keyStr":"value1"},"a2":{"keyObj":{"subKey":"subVal"}}}}`

const badRd = `{"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}`

func createTestResourceDescriptor() (*ResourceDescriptor, error) {
	// Create a ResourceDescriptor
	a1, err := structpb.NewValue(map[string]interface{}{
		"keyStr": "value1",
		"keyNum": 13})
	if err != nil {
		return nil, err
	}
	a2, err := structpb.NewValue(map[string]interface{}{
		"keyObj": map[string]interface{}{
			"subKey": "subVal"}})
	if err != nil {
		return nil, err
	}

	return &ResourceDescriptor{
		Name: "theName",
		Uri:  "https://example.com",
		Digest: map[string]string{
			"alg1": "abc123",
		},
		Content:          []byte("bytescontent"),
		DownloadLocation: "https://example.com/test.zip",
		MediaType:        "theMediaType",
		Annotations:      map[string]*structpb.Value{"a1": a1, "a2": a2},
	}, nil
}

func TestJsonUnmarshalResourceDescriptor(t *testing.T) {
	got := &ResourceDescriptor{}
	err := protojson.Unmarshal([]byte(wantFullRd), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	want, err := createTestResourceDescriptor()

	assert.NoError(t, err, "Error during test RD creation")
	assert.True(t, proto.Equal(got, want), "Protos do not match")
}

func TestBadResourceDescriptor(t *testing.T) {
	got := &ResourceDescriptor{}
	err := protojson.Unmarshal([]byte(badRd), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()

	assert.Error(t, err, "Error: created malformed ResourceDescriptor")
}
