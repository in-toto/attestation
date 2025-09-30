# Envelope layer specification

Version: [DSSE v1.0]

The Envelope is the outermost layer of the attestation, handling
authentication and serialization.

## Schema

The format and protocol are defined per [DSSE v1.0].

## Fields

The in-toto Attestation Framework has the following requirements for the
standard DSSE fields.

-   `payloadType` MUST be set to `application/vnd.in-toto.<predicate>+json` or to
    `application/vnd.in-toto+json`. This indicates that the Envelope contains
    a JSON object with a `_type` field specifying its schema. If the
    predicate-specific media type is used, the following requirements apply:
    -   `<predicate>` MUST match the [predicate specification filename] without
        the file extension.
    -   Consumers SHOULD NOT rely upon the media type for individual attestations
        as faithful indicators of predicate type. Consumer SHOULD only rely on the
        `predicateType` field in the [Statement] layer.
    -   To obtain predicate information that is authenticated, consumers MUST
        parse the Envelope's `payload`, and verify it against its `signatures`.
    -   The predicate version is not specified in the media type; it is handled
        in the [Statement] layer.
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

[Bundle]: bundle.md
[DSSE v1.0]: https://github.com/secure-systems-lab/dsse/blob/v1.0.0/envelope.md
[KEYID]: https://github.com/in-toto/docs/blob/v1.0/in-toto-spec.md#421-key-formats
[Statement]: statement.md
[in-toto-verify]: https://github.com/in-toto/in-toto#verification
[functionaries]: https://github.com/in-toto/docs/blob/v1.0/in-toto-spec.md#212-functionaries
[predicate specification filename]: ../predicates
