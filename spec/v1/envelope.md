# Envelope layer specification

The Envelope is the outermost layer of the attestation, handling authentication and serialization.

A signing envelope may either be [DSSE v1.0] for JSON encoding or [RFC9052] for CBOR encoding.

## Version [DSSE v1.0]

### Schema

The format and protocol are defined per [DSSE v1.0].

### Fields

The in-toto Attestation Framework has the following requirements for the
standard DSSE fields.

-   `payloadType` MUST be set to `application/vnd.in-toto+json`, which
    indicates that the Envelope contains a JSON object with a `_type` field
    specifying its schema.
-   `payload` MUST be a base64-encoded JSON [Statement].

### File naming convention

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

### Storage convention

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

#### Examples

Example media types for single DSSE-signed attestation predicates include:

-   SLSA Provenance: `application/vnd.in-toto.provenance+dsse`
-   SPDX: `application/vnd.in-toto.spdx+dsse`
-   VSA: `application/vnd.in-toto.vsa+dsse`

## Version [RFC9052]

### Schema

The format is defined to be the single signer envelope `COSE_Sign1` from [RFC9052].

### File naming convention

COSE objects do not have a standard file extension.
If stored in a dedicated file by itself, and not as part of a [Bundle], an Envelope SHOULD use the suffix `.intoto.cose`.

### Headers

The in-toto Attestation Framework has the following requirements for the standard COSE headers.

-   content type (label 3) MUST be set to `"application/vnd.in-toto+cbor`, which indicates that the Envelop contains a CBOR object with a `_type` field specifying its schema.

No other headers are required.
The envelope's payload array entry MUST be a CBOR-encoded [Statement].

[Bundle]: bundle.md
[DSSE v1.0]: https://github.com/secure-systems-lab/dsse/blob/v1.0.0/envelope.md
[KEYID]: https://github.com/in-toto/docs/blob/v1.0/in-toto-spec.md#421-key-formats
[Statement]: statement.md
[in-toto-verify]: https://github.com/in-toto/in-toto#verification
[functionaries]: https://github.com/in-toto/docs/blob/v1.0/in-toto-spec.md#212-functionaries
[predicate specification filename]: ../predicates
[RFC9052]: https://www.rfc-editor.org/rfc/rfc9052.html
