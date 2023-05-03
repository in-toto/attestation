# Envelope layer specification

Version: [DSSE v1.0]

The Envelope is the outermost layer of the attestation, handling
authentication and serialization.

## Schema

The format and protocol are defined per [DSSE v1.0].

## Fields

The in-toto Attestation Framework has the following requirements for the
standard DSSE fields.

-   `payloadType` MUST be set to `application/vnd.in-toto+json`, which
    indicates that the Envelope contains a JSON object with a `_type` field
    specifying its schema.
-   `payload` MUST be a base64-encoded JSON [Statement].

[DSSE v1.0]: https://github.com/secure-systems-lab/dsse/blob/v1.0.0/envelope.md
[Statement]: statement.md
