# Attesation v1

## Purpose

This is the "base class" for all attestations.

## Schema

```jsonc
{
  "attestation_type": "<URI>",
  "subject": {
    "<URI-OR-PATH>": {
      "<DIGEST_TYPE>": "<DIGEST_VALUE>"
    }
  },
  "materials": {
    "<URI-OR-PATH>": {
      "<DIGEST_TYPE>": "<DIGEST_VALUE>"
    }
  }
  // Other fields may follow, depending on attestation_type.
}
```

### Standard attestation fields

<a id="attestation_type"></a>
`attestation_type` _string ([TypeURI]), required_

> URI representing the meaning of this attestation and how to interpret the rest
> of the fields. Example:
>
> ```json
> "attestation_type": "https://in-toto.io/Provenance/v1"
> ```
>
> (Somewhat similar to `name` in [in-toto 0.9].)

<a id="subject"></a>
`subject` _object ([ArtifactCollection]), required_

> The collection of software artifacts this attestation is about. Example:
>
> ```json
> "subject": {
>   "curl-7.72.0.tar.bz2": { "sha256": "ad9197…" },
>   "curl-7.72.0.tar.gz": { "sha256": "d4d589…" }
> }
> ```
>
> When there is a single artifact whose name is not meaningful, is RECOMMENDED
> to use `"_"` as the name.
>
> (Roughly equivalent to `products` in [in-toto 0.9].)

<a id="materials"></a>
`materials` _object ([ArtifactCollection]), optional_

> The collection of software artifacts that influenced the attestation, aside
> from the `subject` itself. Example:
>
> ```json
> "materials": {
>   "git+https://github.com/curl/curl@curl-7_72_0": { "git_commit": "9d954e4…" },
>   "pkg:deb/debian/stunnel4@5.50-3?arch=amd64": { "sha256": "e1731ae…" },
>   "pkg:deb/debian/python-impacket@0.9.15-5?arch=all": { "sha256": "71fa2e6…" },
>   "pkg:deb/debian/libzstd-dev@1.3.8+dfsg-3?arch=amd64": { "sha256": "91442b0…" },
>   "pkg:deb/debian/libbrotli-dev@1.0.7-2+deb10u1?arch=amd64": { "sha256": "05b6e46…" }
> }
> ```
>
> (Unchanged from [in-toto 0.9].)

[ArtifactCollection]: field_types.md#ArtifactCollection
[TypeURI]: field_types.md#TypeURI
[in-toto 0.9]: https://github.com/in-toto/docs/blob/v0.9/in-toto-spec.md
