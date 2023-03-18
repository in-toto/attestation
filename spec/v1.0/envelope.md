# Envelope layer specification

Version: v1.0

The Envelope is the outermost layer of the attestation, handling
authentication and serialization. The format and protocol are defined in
[DSSE] and adopted by in-toto in [ITE-5].

## Schema

```jsonc
{
  "payloadType": "application/vnd.in-toto+json",
  "payload": "<Base64(Statement)>",
  "signatures": [{"sig": "<Base64(Signature)>"}]
}
```

## Fields

`payloadType` _string, required_

> Identifier for the encoding of the payload. Always
> `application/vnd.in-toto+json`, which indicates that it is a JSON object
> with a `_type` field indicating its schema.

`payload` _string, required_

> Base64-encoded JSON [Statement].

`signatures` _array of objects, required_

> One or more signatures over `payloadType` and `payload`, as defined in
> [DSSE].

[DSSE]: https://github.com/secure-systems-lab/dsse
[ITE-5]: https://github.com/in-toto/ITE/blob/master/ITE/5/README.adoc
[Statement]: statement.md
