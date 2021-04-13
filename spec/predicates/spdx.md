# Predicate type: SPDX

Type URI: (tentative) https://spdx.dev/Document

TODO: Ask SPDX project to choose a URI and to review this spec. Ideally the URI
would resolve to this file.

## Purpose

A Software Bill of Materials type following the
[SPDX standard](https://spdx.dev/specifications/).

This allows to represent an "exportable" or "published" software artifact. It
can also be used as an entry point for other types of in-toto attestation when
performing policy decisions

## Schema

```jsonc
{
  "subject": { ... }
  "predicateType": "https://spdx.dev/Document",
  "predicate": {
    "SPDXID" : "SPDXRef-DOCUMENT",
    "spdxVersion" : "SPDX-2.2",
    ...
  }
}
```

The `predicate` contains a JSON-encoded SPDX document. The SPDX format has a
mandatory `spdxVersion` field, so there is no need to version the predicate
type.

The `subject` contains whatever software artifacts are to be associated with
this SPDX document.
