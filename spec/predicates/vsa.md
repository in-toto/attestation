# Predicate type: Verification Summary Attestation (VSA)

Type URI: https://in-toto.io/attestation/verification_summary/v1

Version: 1.1

## Purpose

Verification summary attestations communicate that an artifact has been verified
against some policy and details about that verification. Such details may
include, but are not limited to, what SLSA level the artifact and/or its
dependencies were verified to meet.

## Use Cases

Software consumers make a decision about the validity of an
artifact without needing to have access to all of the attestations about the
artifact or all of its transitive dependencies.  They delegate
complex policy decisions to some trusted party and then simply trust that
party's decision regarding the artifact. That decision is represented by a VSA.

VSAs also allow software producers to keep the details of their build pipeline
confidential while still communicating that some verification has taken place.
This might be necessary for legal reasons (keeping a software supplier
confidential) or for security reasons (not revealing that an embargoed patch has
been included).

## Prerequisites

Understanding of [in-toto specification] and the in-toto attestation framework.

## Model

A Verification Summary Attestation (VSA) is an attestation that some entity
(`verifier`) verified one or more software artifacts (the `subject` of an
in-toto attestation [Statement]) by evaluating the artifact and a `bundle`
of attestations against some `policy`.  Users who trust the `verifier` may
assume that the artifacts satisfied the indicated security polciy without
themselves needing to evaluate the artifact or to have access to the
attestations the `verifier` used to make its determination.

The VSA also allows consumers to determine the verified levels of
all of an artifact’s _transitive_ dependencies.  The verifier does this by
either a) verifying the provenance of each non-source dependency listed in
the [resolvedDependencies](/provenance/v1#resolvedDependencies) of the artifact
being verified (recursively) or b) matching the non-source dependency
listed in `resolvedDependencies` (`subject.digest` ==
`resolvedDependencies.digest` and, ideally, `vsa.resourceUri` ==
`resolvedDependencies.uri`) to a VSA _for that dependency_ and using
`vsa.verifiedProperties` and `vsa.dependencyProperties`.  Policy verifiers
wishing to establish minimum requirements on dependencies may use
`vsa.dependencyProperties` to do so.

## Schema

```jsonc
// Standard attestation fields:
"_type": "https://in-toto.io/Statement/v1",
"subject": [{
  "name": <NAME>,
  "digest": { <digest-in-request> }
}],

// Predicate
"predicateType": "https://slsa.dev/verification_summary/v1",
"predicate": {
  "verifier": {
    "id": "<URI>",
    "version": {
      "<COMPONENT>": "<VERSION>",
      ...
    }
  },
  "timeVerified": <TIMESTAMP>,
  "resourceUri": <artifact-URI-in-request>,
  "policy": {
    "uri": "<URI>",
    "digest": { /* DigestSet */ }
  }
  "inputAttestations": [
    {
      "uri": "<URI>",
      "digest": { <digest-of-attestation-data> }
    },
    ...
  ],
  "verificationResult": "<PASSED|FAILED>",
  "verifiedProperties": ["<String>"],
  "dependencyProperties": {
    "<String>": <Int>,
    "<String>": <Int>,
    ...
  },
  "slsaVersion": "<MAJOR>.<MINOR>",
}
```

### Parsing rules

This predicate follows the in-toto attestation [parsing rules]. Summary:

-   Consumers MUST ignore unrecognized fields.
-   The `predicateType` URI includes the major version number and will always
    change whenever there is a backwards incompatible change.
-   Minor version changes are always backwards compatible and "monotonic." Such
    changes do not update the `predicateType`.
-   Producers MAY add extension fields using field names that are URIs.

### Fields

_NOTE: This section describes the fields within `predicate`. For a description
of the other top-level fields, such as `subject`, see [Statement]._

<a id="verifier"></a>
`verifier` _object, required_

> Identifies the entity that performed the verification.
>
> The identity MUST reflect the trust base that consumers care about. How
> detailed to be is a judgment call.
>
> Consumers MUST accept only specific (signer, verifier) pairs. For example,
> "GitHub" can sign provenance for the "GitHub Actions" verifier, and "Google"
> can sign provenance for the "Google Cloud Deploy" verifier, but "GitHub" cannot
> sign for the "Google Cloud Deploy" verifier.
>
> The field is required, even if it is implicit from the signer, to aid readability and
> debugging. It is an object to allow additional fields in the future, in case one
> URI is not sufficient.

<a id="verifier.id"></a>
`verifier.id` _string ([TypeURI]), required_

> URI indicating the verifier’s identity.

<a id="verifier.version"></a>
`verifier.version` _map (string->string), optional_

> Map of names of components of the verification platform to their version.

<a id="timeVerified"></a>
`timeVerified` _string ([Timestamp]), optional_

> Timestamp indicating what time the verification occurred.

<a id="resourceUri"></a>
`resourceUri` _string ([ResourceURI]), required_

> URI that identifies the resource associated with the artifact being verified.

<a id="policy"></a>
`policy` _object ([ResourceDescriptor]), required_

> Describes the policy that the `subject` was verified against.
>
> The entry MUST contain a `uri`.

<a id="inputAttestations"></a>
`inputAttestations` _array ([ResourceDescriptor]), optional_

> The collection of attestations that were used to perform verification.
> Conceptually similar to the `resolvedDependencies` field in [SLSA Provenance].
>
> This field MAY be absent if the verifier does not support this feature.
> If non-empty, this field MUST contain information on _all_ the attestations
> used to perform verification.
>
> Each entry MUST contain a `digest` of the attestation and SHOULD contains a
> `uri` that can be used to fetch the attestation.

<a id="verificationResult"></a>
`verificationResult` _string, required_

> Either “PASSED” or “FAILED” to indicate if the artifact passed or failed the policy verification.

<a id="verifiedProperties"></a>
`verifiedProperties` _array (_string), required_

> Indicates the properties verified for the artifact (and not
> its dependencies), or "FAILED" if policy verification failed.
>
> If including SLSA levels, users MUST NOT include more than one level per SLSA
> track. Note that each SLSA level implies all levels below it (e.g.
> SLSA_BUILD_LEVEL_3 implies SLSA_BUILD_LEVEL_2 and SLSA_BUILD_LEVEL_1), so
> there is no need to include more than one level per track.

<a id="dependencyProperties"></a>
`dependencyProperties` _object, optional_

> A count of the dependencies verified to have each property.
>
> Map from _string to the number of the artifact's _transitive_ dependencies
> that were verified at the indicated level. Absence of a given value
>  MUST be interpreted as reporting _0_ dependencies with that value.
>
> If including SLSA levels, users MUST count each dependency only once per SLSA
> track, at the highest level verified. For example, if a dependency meets
> SLSA_BUILD_LEVEL_2, you include it with the count for SLSA_BUILD_LEVEL_2 but
> not the count for SLSA_BUILD_LEVEL_1.

<a id="slsaVersion"></a>
`slsaVersion` _string, optional_

> Indicates the version of the SLSA specification that the verifier used, in the
> form `<MAJOR>.<MINOR>`. Example: `1.0`. If unset, the default is an
> unspecified minor version of `1.x`.

## Example

WARNING: This is just for demonstration purposes.

```jsonc
"_type": "https://in-toto.io/Statement/v1",
"subject": [{
  "name": "out/example-1.2.3.tar.gz",
  "digest": {"sha256": "5678..."}
}],

// Predicate
"predicateType": "https://in-toto.io/attestation/verification_summary/v1",
"predicate": {
  "verifier": {
    "id": "https://example.com/publication_verifier",
    "version": {
      "slsa-verifier-linux-amd64": "v2.3.0",
      "slsa-framework/slsa-verifier/actions/installer": "v2.3.0"
    }
  },
  "timeVerified": "1985-04-12T23:20:50.52Z",
  "resourceUri": "https://example.com/example-1.2.3.tar.gz",
  "policy": {
    "uri": "https://example.com/example_tarball.policy",
    "digest": {"sha256": "1234..."}
  },
  "inputAttestations": [
    {
      "uri": "https://example.com/provenances/example-1.2.3.tar.gz.intoto.jsonl",
      "digest": {"sha256": "abcd..."}
    }
  ],
  "verificationResult": "PASSED",
  "verifiedProperties": ["SLSA_BUILD_LEVEL_3"],
  "dependencyProperties": {
    "SLSA_BUILD_LEVEL_3": 5,
    "SLSA_BUILD_LEVEL_2": 7,
    "SLSA_BUILD_LEVEL_1": 1,
  },
  "slsaVersion": "1.0"
}
```

<div id="verificationresult">

## Change history

-   1.1:
    -   Added optional `verifier.version` field to record verification tools.
    -   Added Verification section with examples.
    -   Made `timeVerified` optional.
-   1.0:
    -   Replaced `materials` with `resolvedDependencies`.
    -   Relaxed `VerificationResult` to allow other values.
    -   Converted to lowerCamelCase for consistency with [SLSA Provenance].
    -   Added `slsaVersion` field.
-   0.2:
    -   Added `resource_uri` field.
    -   Added optional `input_attestations` field.
-   0.1: Initial version.

[SLSA Provenance]: /provenance
[DigestSet]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/digest_set.md
[ResourceURI]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/field_types.md#resourceuri
[ResourceDescriptor]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/resource_descriptor.md
[Statement]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/statement.md
[Timestamp]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/field_types.md#timestamp
[TypeURI]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/field_types.md#TypeURI
[in-toto attestation]: https://github.com/in-toto/attestation
[parsing rules]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/README.md#parsing-rules
[in-toto specification]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md