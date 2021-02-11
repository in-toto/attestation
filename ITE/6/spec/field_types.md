# Field type definitions

<a id="TypeURI"></a>
_TypeURI (string)_

> Uniform Resource Identifier as specified in [RFC 3986], used as a
> collision-resistant type identifier. Case sensitive and MUST be case
> normalized as per section 6.2.2.1 of RFC 3986, meaning that the scheme and
> authority MUST be in lowercase. SHOULD resolve to a human-readable
> description, but MAY be unresolvable. SHOULD include a version number to allow
> for revisions.
>
> Example: `"https://in-toto.io/Attestation/v1"`.

<a id="ArtifactCollection"></a>
_ArtifactCollection (object)_

> A collection of software artifacts. Each key/value pair represents a single
> software artifact, such as a file, a container image, or a git commit.
>
> The key identifies the artifact relative to this attestation. It MUST be a URI
> or path-noscheme ([RFC 3986]). If it is a URI, it MUST be case normalized and
> SHOULD resolve to the artifact, but MAY be unresolvable. It is RECOMMENDED to
> use [Package URL][] (`pkg:`) or [SPDX Download Location][] (e.g.
> `git+https:`). If path-noscheme, it SHOULD represent a relative filesystem
> path.
>
> The value contains cryptographic digests of the artifact's content. It is a
> map from digest type to digest value (string), encoded as lowercase hex. The
> digest type unambiguously identifies the hash algorithm and how it is applied
> to the artifact. An artifact MUST be considered matching if *any* of its
> digests match. Verifiers MUST choose which digest types they accept and MUST
> ignore digest types they do not accept or recognize. The value MAY be empty or
> null, meaning that no digest is available.
>
> The following digests types are RECOMMENDED:
>
> *   Regular File: `sha256`
> *   Git repository: `git_commit`
> *   Mercurial repository: `hg_changeset`
> *   Container image: [`oci_image_id`][oci_image_id] and `oci_repo_digest`. It
>     is best to list both when available. The former is registry independent
>     (over the uncompressed manifest) while the latter depends on the registry
>     (computed over the compressed manifest).
>
> Additional pre-defined digest types: `sha224`, `sha384`, `sha512`,
> `sha512_224`, `sha512_256`, `sha3_224`, `sha3_256`, `sha3_384`, `sha3_512`,
> `shake128`, `shake256`, `blake2b`, `blake2s`, `md5` (DISCOURAGED), `sha1`
> (DISCOURAGED).
>
> Custom digest types MAY be used if none of the recommended or pre-defined
> digest types work.
>
> Example:
>
> ```jsonc
> {
>   "pkg:docker/alpine@3.13.1?arch=amd64": {
>     "oci_image_id": "sha256:e50c909a8df2b7c8b92a6e8730e210ebe98e5082871e66edd8ef4d90838cbd25",
>     "oci_repo_digest": "sha256:3747d4eb5e7f0825d54c8e80452f1e245e24bd715972c919d189a62da97af2ae"
>   },
>   "git+https://github.com/curl/curl@curl-7_72_0": {
>     "git_commit": "9d954e49bce3706a9a2efb119ecd05767f0f2a9e"
>   },
>   "pkg:deb/debian/stunnel4@5.50-3?arch=amd64":
>     "sha256": "e1731ae217fcbc64d4c00d707dcead45c828c5f762bcf8cc56d87de511e096fa"
>   },
>   "github_hosted_vm:ubuntu-18.04:20210123.1": null  // no digest available
> }
> ```

<a id="Timestamp"></a>
_Timestamp (string)_

> A point in time, represented as a string in [RFC 3339] format in the UTC time
> zone ("Z").
>
> Example: `"1985-04-12T23:20:50.52Z"`.

[Package URL]: https://github.com/package-url/purl-spec/
[RFC 3339]: https://tools.ietf.org/html/rfc3339
[RFC 3986]: https://tools.ietf.org/html/rfc3986
[SPDX Download Location]: https://spdx.github.io/spdx-spec/3-package-information/#37-package-download-location
[oci_image_id]: https://github.com/opencontainers/image-spec/blob/master/config.md#imageid
