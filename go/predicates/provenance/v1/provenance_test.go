/*
Tests for SLSA Provenance v1 protos.
*/

package v1

import (
	"fmt"
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

func TestBadProvenanceBuildDefinition(t *testing.T) {
	tests := map[string]struct {
		input        string
		err          error
		noErrMessage string
	}{
		"no buildDefinition": {
			input:        `{"buildDefinition":{},"runDetails":{"builder":{"id":"theId","version":{"theComponent":"v0.1"},"builderDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"metadata":{"invocationId":"theInvocationId"},"byproducts":[{"name":"theResource","digest":{"alg1":"abc123"}}]}}`,
			err:          ErrBuildDefinitionRequired,
			noErrMessage: "created malformed Provenance (empty buildDefinition)",
		},
		"buildDefinition missing buildType required field": {
			input:        `{"buildDefinition":{"externalParameters":{"param1":{"subKey":"subVal"}},"resolvedDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"runDetails":{"builder":{"id":"theId","version":{"theComponent":"v0.1"},"builderDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"metadata":{"invocationId":"theInvocationId"},"byproducts":[{"name":"theResource","digest":{"alg1":"abc123"}}]}}`,
			err:          ErrBuildTypeRequired,
			noErrMessage: "created malformed Provenance (buildDefinition missing required buildType field)",
		},
		"buildDefinition missing externalParameters required field": {
			input:        `{"buildDefinition":{"buildType":"theBuildType","resolvedDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"runDetails":{"builder":{"id":"theId","version":{"theComponent":"v0.1"},"builderDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"metadata":{"invocationId":"theInvocationId"},"byproducts":[{"name":"theResource","digest":{"alg1":"abc123"}}]}}`,
			err:          ErrExternalParamsRequired,
			noErrMessage: "created malformed Provenance (buildDefinition missing requried externalParameters field)",
		},
	}

	for name, test := range tests {
		got := &Provenance{}
		err := protojson.Unmarshal([]byte(test.input), got)
		assert.NoError(t, err, fmt.Sprintf("error during JSON unmarshalling in test '%s'", name))

		err = got.Validate()
		assert.ErrorIs(t, err, test.err, fmt.Sprintf("%s in test '%s'", test.noErrMessage, name))
	}
}

func TestBadProvenanceRunDetails(t *testing.T) {
	tests := map[string]struct {
		input        string
		err          error
		noErrMessage string
	}{
		"no runDetails": {
			input:        `{"buildDefinition":{"buildType":"theBuildType","externalParameters":{"param1":{"subKey":"subVal"}},"resolvedDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"runDetails":{}}`,
			err:          ErrRunDetailsRequired,
			noErrMessage: "created malformed Provenance (empty runDetails)",
		},
		"runDetails missing builder required field": {
			input:        `{"buildDefinition":{"buildType":"theBuildType","externalParameters":{"param1":{"subKey":"subVal"}},"resolvedDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"runDetails":{"builder":{},"metadata":{"invocationId":"theInvocationId"},"byproducts":[{"name":"theResource","digest":{"alg1":"abc123"}}]}}`,
			err:          ErrBuilderRequired,
			noErrMessage: "created malformed Provenance (runDetails missing required builder field)",
		},
		"runDetails.builder missing id required field": {
			input:        `{"buildDefinition":{"buildType":"theBuildType","externalParameters":{"param1":{"subKey":"subVal"}},"resolvedDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"runDetails":{"builder":{"id":"","version":{"theComponent":"v0.1"},"builderDependencies":[{"name":"theResource","digest":{"alg1":"abc123"}}]},"metadata":{"invocationId":"theInvocationId"},"byproducts":[{"name":"theResource","digest":{"alg1":"abc123"}}]}}`,
			err:          ErrBuilderIdRequired,
			noErrMessage: "created malformed Provenance (runDetails.builder missing requried id field)",
		},
	}

	for name, test := range tests {
		got := &Provenance{}
		err := protojson.Unmarshal([]byte(test.input), got)
		assert.NoError(t, err, fmt.Sprintf("error during JSON unmarshalling in test '%s'", name))

		err = got.Validate()
		assert.ErrorIs(t, err, test.err, fmt.Sprintf("%s in test '%s'", test.noErrMessage, name))
	}
}
