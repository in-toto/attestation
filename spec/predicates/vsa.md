# Predicate type: Verification Summary

Type URI: https://in-toto.io/attestation/verification_summary/v1

Version: 1.1

## Purpose

Verification summary attestations communicate that an artifact has been verified
against some policy and details about that verification. Such details may
include, but are not limited to, what [SLSA](https://slsa.dev) level the
artifact and/or its dependencies were verified to meet.

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
"predicateType": "https://in-toto.io/attestation/verification_summary/v1",
"predicate": {
  "verifier": {
    "id": "<URI>",
    "version": {
      "<COMPONENT>": "<VERSION>",
      ...
    }
  },
  "timeVerified": <TIMESTAMP>,
  "verifiedUri": <artifact-URI-in-request>,
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
> its dependencies).

<a id="dependencyProperties"></a>
`dependencyProperties` _object, optional_

> A count of the dependencies verified to have each property.
>
> Map from _string_ to the number of the artifact's _transitive_ dependencies
> that were verified at the indicated level. Absence of a given value
>  MUST be interpreted as reporting _0_ dependencies with that value.

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
}
```

<div id="verificationresult">

## How to verify

VSA consumers use VSAs to accomplish goals based on delegated trust. We call the
process of establishing a VSA's authenticity and determining whether it meets
the consumer's goals 'verification'. Goals differ, as do levels of confidence
in VSA producers, so the verification procedure changes to suit its context.
However, there are certain steps that most verification procedures have in
common.

Verification MUST include the following steps:

1.  Verify the signature on the VSA envelope using the preconfigured roots of
    trust. This step ensures that the VSA was produced by a trusted producer
    and that it hasn't been tampered with.

2.  Verify the statement's `subject` matches the digest of the artifact in
    question. This step ensures that the VSA pertains to the intended artifact.

3.  Verify that the `predicateType` is
    `https://in-toto.io/attestation/verification_summary/v1`. This step ensures
    that the in-toto predicate is using this version of the VSA format.

4.  Verify that the `verifier` matches the public key (or equivalent) used to
    verify the signature in step 1. This step identifies the VSA producer in
    cases where their identity is not implicitly revealed in step 1.

5.  Verify that the value for `resourceUri` in the VSA matches the expected
    value. This step ensures that the consumer is using the VSA for the
    producer's intended purpose.

6.  Verify that the value for `verificationResult` is `PASSED`. This step
    ensures the artifact is suitable for the consumer's purposes.

7.  Verify that `verifiedLevels` contains the expected value. This step ensures
    that the artifact is suitable for the consumer's purposes.

Verification MAY additionally contain the following step:

1.  (Optional) Verify additional fields required to determine whether the VSA
    meets your goal.

Verification mitigates different threats depending on the VSA's contents and the
verification procudure.

IMPORTANT: A VSA does not protect against compromise of the verifier, such as by
a malicious insider. Instead, VSA consumers SHOULD carefully consider which
verifiers they add to their roots of trust.

### Examples

1.  Suppose consumer C wants to delegate to verifier V the decision for whether
    to accept artifact A as resource R. Consumer C verifies that:

    -   The signature on the VSA envelope using V's public signing key from their
      preconfigured root of trust.

    -   `subject` is A.

    -   `predicateType` is `https://in-toto.io/attestation/verification_summary/v1`.

    -   `verifier.id` is V.

    -   `resourceUri` is R.

    -   `verificationResult` is `PASSED`.

    -   `verifiedlevels` contains `SLSA_BUILD_LEVEL_UNEVALUATED`.

    Note: This example is analogous to traditional code signing. The expected
    value for `verifiedProperties` is arbitrary but prenegotiated by the
    producer and the consumer. The consumer does not need to check additional
    fields, as C fully delegates the decision to V.

2.  Suppose consumer C wants to enforce the rule "Artifact A at resource R must
    have a passing VSA from verifier V showing it meets SLSA Build Level 2+."
    Consumer C verifies that:

    -   The signature on the VSA envelope using V's public signing key from their
      preconfigured root of trust.

    -   `subject` is A.

    -   `predicateType` is `https://in-toto.io/attestation/verification_summary/v1`.

    -   `verifier.id` is V.

    -   `resourceUri` is R.

    -   `verificationResult` is `PASSED`.

    -   `verifiedProperties` is `SLSA_BUILD_LEVEL_2` or `SLSA_BUILD_LEVEL_3`.

    Note: In this example, verifying the VSA mitigates the same threats as
    verifying the artifact's SLSA provenance. See
    [Verifying artifacts](https://slsa.dev/spec/v1.0/verifying-artifacts) for
    details about which threats are addressed by verifying each SLSA level.

## Change history

-   1.1:
    -   Added optional `verifier.version` field to record verification tools.
    -   Added Verification section with examples.
    -   Made `timeVerified` optional.
    -   Moved from SLSA specification to in-toto attestation framework.
-   1.0:
    -   Replaced `materials` with `resolvedDependencies`.
    -   Relaxed `VerificationResult` to allow other values.
    -   Converted to lowerCamelCase for consistency with [SLSA Provenance].
    -   Added `slsaVersion` field.
-   0.2:
    -   Added `resource_uri` field.
    -   Added optional `input_attestations` field.
-   0.1: Initial version.

[SLSA Provenance]: https://slsa.dev/spec/v1.0/provenance
[DigestSet]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/digest_set.md
[ResourceURI]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/field_types.md#resourceuri
[ResourceDescriptor]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/resource_descriptor.md
[Statement]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/statement.md
[Timestamp]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/field_types.md#timestamp
[TypeURI]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/field_types.md#TypeURI
[in-toto attestation]: https://github.com/in-toto/attestation
[parsing rules]: https://github.com/in-toto/attestation/blob/7aefca35a0f74a6e0cb397a8c4a76558f54de571/spec/v1/README.md#parsing-rules
[in-toto specification]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md