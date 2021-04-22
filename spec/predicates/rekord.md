# Predicate type: Rekord v1

Type URI: https://sigstore.dev/Rekord

## Purpose
Support for the sigstore's [rekord type](https://github.com/sigstore/rekor). Rekord is a default type for the rekor, which enables software maintainers and build systems to record signed metadata to an immutable record.

## Schema

```jsonc
{
  "subject": [{ ... }],
  "predicateType": "https://sigstore.dev/Rekord", // or https://rekord.sigstore.dev/Rekord - is there already a pre-defined URI?
  "predicate": {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "$id": "http://rekor.sigstore.dev/types/rekord/rekord_schema.json",
    "title": "Rekor Schema",
    "description": "Schema for Rekord objects",
    "type": "object",
    "oneOf": [
        {
            "$ref": "v0.0.1/rekord_v0_0_1_schema.json"
        }
    ]
  }
}