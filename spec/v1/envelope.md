# Envelope layer specification

The Envelope is the outermost layer of the attestation, handling  serialization
and authentication (via digital signatures).

## Schema

The RECOMMENDED format and protocol for Envelopes are defined per [DSSE v1.0].
Producers MAY use other signature methods and formats that meet the [ITE-5]
specification:

-   MUST support the inclusion of multiple signatures in a single envelope
-   SHOULD include an authenticated payload type
-   SHOULD avoid depending on canonicalization for security
-   SHOULD support a hint indicating what signing key was used, i.e., a KEYID
-   SHOULD NOT require the verifier to parse the payload before verifying
-   SHOULD NOT require the inclusion of signing key algorithms in the signature

### Alternative Envelope schemas

-   The [Sigstore Bundle], while [supporting DSSE], is not currently [ITE-5]
    compliant because it requires a _single signature_ in the envelope.[^1]
-   The [COSE_Sign] structure is [ITE-5] compliant, whereas the `COSE_Sign1`
    format that supports only one signer is NOT compliant.

## Fields

The in-toto Attestation Framework has the following general field requirements
for an Envelope:

-   `signature` (or equivalent) is REQUIRED and MUST be defined as an array.
-   A `keyid` (or equivalent) SHOULD be included for each signing key used.
-   `payload` (or equivalent) SHOULD be included and contain the attestation
    data that was signed.
-   `payloadType` (or equivalent) SHOULD be signed along with the `payload`.

In addition, the Envelope spec has the following specific requirements for the
standard [DSSE][DSSE v1.0] fields.

-   `payloadType` MUST be set to `application/vnd.in-toto+json`, which
    indicates that the Envelope contains a JSON object with a `_type` field
    specifying its schema.
-   `payload` MUST be a base64-encoded JSON [Statement].

## File naming convention

If stored in a dedicated file by itself, and not as part of a [Bundle], an
Envelope SHOULD use the suffix `.json`.

-   For attestations intended for consumption by [in-toto-verify], an
    Envelope containing an attestation about a particular software supply
    chain step `<step-name>` SHOULD be named `<step-name>.json`.
-   For other verifiers, or cases in which a step name cannot be easily
    determined, the attestation producer and consumer MAY agree on an
    arbitrary filename: `<env-name>.json`.
-   If multiple Envelopes are produced for the same step by different
    [functionaries] uniquely identified by a public key, an Envelope name
    SHOULD include the truncated [KEYID] of the public key `<keyid[0:8]>` of
    the signing functionary: `<step/env-name>.<keyid[0:8]>.json`.

## Storage convention

The media type `application/vnd.in-toto.<predicate>+dsse` SHOULD
be used to denote an individual attestation in arbitrary storage systems.

-   The `<predicate>` MUST match the [predicate specification filename]
    without the file extension. Predicate versioning is handled in the
    [Statement] layer.
-   Consumers SHOULD NOT rely upon the media type for individual attestations
    as faithful indicators of predicate type. Consumer SHOULD only rely on the
    `predicateType` field in the [Statement] layer.
-   To obtain predicate information that is authenticated, consumers MUST
    parse the Envelope's `payload`, and verify it against its `signatures`.

### Examples

Example media types for single DSSE-signed attestation predicates include:

-   SLSA Provenance: `application/vnd.in-toto.provenance+dsse`
-   SPDX: `application/vnd.in-toto.spdx+dsse`
-   VSA: `application/vnd.in-toto.vsa+dsse`

[^1]: There is an [ongoing discussion](https://github.com/sigstore/sig-clients/issues/9) about supporting [DSSE Signature Extensions](https://github.com/secure-systems-lab/dsse/blob/devel/envelope.md#signature-extensions-experimental) to extend the current features of Sigstore Bundles.

[Bundle]: bundle.md
[COSE_Sign]: https://datatracker.ietf.org/doc/html/rfc8152#section-4.1
[DSSE v1.0]: https://github.com/secure-systems-lab/dsse/blob/v1.0.2/envelope.md
[ITE-5]: https://github.com/in-toto/ITE/tree/master/ITE/5#specification
[KEYID]: https://github.com/in-toto/docs/blob/v1.0/in-toto-spec.md#421-key-formats
[Sigstore Bundle]: https://docs.sigstore.dev/about/bundle/
[Statement]: statement.md
[in-toto-verify]: https://github.com/in-toto/in-toto#verification
[functionaries]: https://github.com/in-toto/docs/blob/v1.0/in-toto-spec.md#212-functionaries
[predicate specification filename]: ../predicates
[supporting DSSE]: https://docs.sigstore.dev/about/bundle/#dsse
