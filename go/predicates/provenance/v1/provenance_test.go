/*
Tests for SLSA Provenance v1 protos.
*/

package v1

import (
	"testing"

	ita1 "github.com/in-toto/attestation/go/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func createTestProvenance(t *testing.T) *Provenance {
	// Create a Provenance

	t.Helper()

	rd := &ita1.ResourceDescriptor{
		Name:   "theResource",
		Digest: map[string]string{"alg1": "abc123"},
	}

	builder := &Builder{
		Id:                  "theId",
		Version:             map[string]string{"theComponent": "v0.1"},
		BuilderDependencies: []*ita1.ResourceDescriptor{rd},
	}

	buildMeta := &BuildMetadata{
		InvocationId: "theInvocationId",
	}

	runDetails := &RunDetails{
		Builder:    builder,
		Metadata:   buildMeta,
		Byproducts: []*ita1.ResourceDescriptor{rd},
	}

	externalParams, err := structpb.NewStruct(map[string]interface{}{
		"param1": map[string]interface{}{
			"subKey": "subVal"}})
	if err != nil {
		t.Fatal(err)
	}

	buildDef := &BuildDefinition{
		BuildType:            "theBuildType",
		ExternalParameters:   externalParams,
		ResolvedDependencies: []*ita1.ResourceDescriptor{rd},
	}

	return &Provenance{
		BuildDefinition: buildDef,
		RunDetails:      runDetails,
	}
}

func TestJsonUnmarshalProvenance(t *testing.T) {
	var wantSt = `{"buildDefinition":{"buildType":"theBuildType","externalParameters":{"param1":{"subKey":"subVal"}},"resolvedDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"runDetails":{"builder":{"id":"theId","version":{"theComponent":"v0.1"},"builderDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"metadata":{"invocationId":"theInvocationId"},"byproducts":[{"name":"theResource","digest":{"alg1":"abc123"}}]}}`

	got := &Provenance{}
	err := protojson.Unmarshal([]byte(wantSt), got)
	assert.NoError(t, err, "error during JSON unmarshalling")

	want := createTestProvenance(t)
	assert.NoError(t, err, "unexpected error during test Statement creation")
	assert.True(t, proto.Equal(got, want), "protos do not match")
}
