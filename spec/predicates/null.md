# Predicate type: Null

Type URI: https://in-toto.io/Null

Version: 1.0.0

## Purpose

Denotes the absence of a predicate. May be used as a drop-in replacement for
traditional code signing, where the notion of a predicate does not exist.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/Null",
  "predicate": null  // or unset
}
```

_(Note: This is a Predicate type that fits within the larger
[Attestation](../README.md) framework.)_

This predicate has no fields.
