# Attestation type: SBOM v1

## Purpose

A Software Bill of Materials type following the SPDX 3.0/3T SBOM std.

This allows to represent an "exportable" or "published" software artifact. It
can also be used as an entry point for other types of in-toto attestation when
performing policy decisions

## Schema

```jsonc
{
  "attestation_type": "https://in-toto.io/Provenance/v1",
  "subject": {
    "<URI-OR-PATH>": {
      "<DIGEST_TYPE>": "<DIGEST_VALUE>"
    }
  },
  "materials": {
    "<URI-OR-PATH>": {
      "<DIGEST_TYPE>": "<DIGEST_VALUE>"
    }
  },
  "_name": "<STRING>",
  "document_information": {
    "description": "<STRING>",
    "data_license": "<STRING>",
    "SPDXID": "<SPDXID>",
    "doc_name": "<STRING>",
    "doc_namespace": "<STRING>",
    "external_references": ["<SPDXID>", ...],
    "created": "<TIMESTAMP>",
    "creator": "<STRING>",
    "comment": "<STRING>",
  },

  /* optional sub-fields (maybe subtypes?) */
  "package_information": {...},
  "file_information": {...},
  "snippet_information": {...},
  "other_licensing": {...},
  "relationships": {...},
  "annotations": {...},

}
```
