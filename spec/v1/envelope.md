# Envelope layer specification

Version: [DSSE v1.0]

The Envelope is the outermost layer of the attestation, handling
authentication and serialization.

## Schema

The RECOMMENDED format and protocol for Envelopes are defined per [DSSE v1.0].
Producers MAY use other signature methods and formats so long as they meet
the [Bundle] data structure requirements.

## Fields

The in-toto Attestation Framework has the following requirements for the
standard DSSE fields.

-   `payloadType` MUST be set to `application/vnd.in-toto+json`, which
    indicates that the Envelope contains a JSON object with a `_type` field
    specifying its schema.
-   `payload` MUST be a base64-encoded JSON [Statement].

## File naming convention

Envelopes SHOULD use the suffix `.json`.

-   An Envelope containing an attestation about a particular SW supply chain
    step `<step-name>` SHOULD be named `<step-name>.json`.
-   If multiple Envelopes are produced for the same step by different
    [functionaries] uniquely identified by a public key, an Envelope name
    SHOULD include the hash of the public key `<pubkey-hash>` of the signing
    functionary: `<step-name>.<pubkey-hash>.json`.

## Storage convention

The media type `application/vnd.in-toto.<predicate>+<sig>` SHOULD
be used to denote an individual attestation in arbitrary storage systems.

-   The `<predicate>` MUST match the [predicate specification name]
    without the file extension. Predicate versioning is handled in the
    [Statement] layer.
-   The `<sig>` MUST be a succint alias that unambiguously identifies
    the Envelope signature format.
-   Consumers SHOULD NOT rely upon the media type for individual attestations
    as faithful indicators of predicate type because this information is only
    authenticated at the [Statement] layer.
-   To obtain predicate information that is authenticated, consumers MUST
    parse the Envelope's `payload`.

### Example

The media type for a single DSSE-signed attestation containing an
[SPDX predicate] SHOULD be `application/vnd.in-toto.spdx+dsse`.

[Bundle]: bundle.md
[DSSE v1.0]: https://github.com/secure-systems-lab/dsse/blob/v1.0.0/envelope.md
[SPDX predicate]: ../predicates/spdx.md
[Statement]: statement.md
[functionaries]: https://github.com/in-toto/docs/blob/v1.0/in-toto-spec.md#212-functionaries
[predicate specification name]: ../predicates
