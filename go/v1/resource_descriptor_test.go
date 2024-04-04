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

const supportedRdDigest = `{"digest":{"sha256":"a1234567b1234567c1234567d1234567e1234567f1234567a1234567b1234567","custom":"myCustomEnvoding","sha1":"a1234567b1234567c1234567d1234567e1234567"}}`

const badRd = `{"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}`

const badRdDigestEncoding = `{"digest":{"sha256":"badDigest"},"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}`

const badRdDigestLength = `{"digest":{"sha256":"abc123"},"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}`

func createTestResourceDescriptor() (*ResourceDescriptor, error) {
	// Create a ResourceDescriptor
	a, err := structpb.NewStruct(map[string]interface{}{
		"a1": map[string]interface{}{
			"keyStr": "value1",
			"keyNum": 13},
		"a2": map[string]interface{}{
			"keyObj": map[string]interface{}{
				"subKey": "subVal"}},
	})
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
		Annotations:      a,
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

func TestSupportedResourceDescriptorDigest(t *testing.T) {
	got := &ResourceDescriptor{}
	err := protojson.Unmarshal([]byte(supportedRdDigest), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()
	assert.NoError(t, err, "Error during validation of valid supported RD digests")
}

func TestBadResourceDescriptor(t *testing.T) {
	got := &ResourceDescriptor{}
	err := protojson.Unmarshal([]byte(badRd), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()
	assert.ErrorIs(t, err, ErrRDRequiredField, "created malformed ResourceDescriptor")
}

func TestBadResourceDescriptorDigestEncoding(t *testing.T) {
	got := &ResourceDescriptor{}
	err := protojson.Unmarshal([]byte(badRdDigestEncoding), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()
	assert.ErrorIs(t, err, ErrInvalidDigestEncoding, "created ResourceDescriptor with invalid digest encoding")
}

func TestBadResourceDescriptorDigestLength(t *testing.T) {
	got := &ResourceDescriptor{}
	err := protojson.Unmarshal([]byte(badRdDigestLength), got)

	assert.NoError(t, err, "Error during JSON unmarshalling")

	err = got.Validate()
	assert.ErrorIs(t, err, ErrIncorrectDigestLength, "created ResourceDescriptor with incorrect digest length")
}
