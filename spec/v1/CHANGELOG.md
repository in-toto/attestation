# Changelog

## v1.2

-   Expanded the allowed values for the envelope's `payloadType` field to
`application/vnd.in-toto.<predicate>+json`. See [Envelope Fields].

## v1.1

-   Clarified that subjects are assumed to be immutable and that it is
acceptable to use a non-cryptographic digest (though cryptographic
digests are still strongly recommended).

## v1

Initial release.

[Envelope Fields]: envelope.md#Fields
