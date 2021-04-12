# Field type definitions

<a id="DigestSet"></a>
_DigestSet (object)_

> Set of alternative cryptographic digests, expressed as a JSON map from
> algorithm name to lowercase hex-encoded value. See
> [digest.proto](digest.proto) for more details, including standard algorithm
> names.

<a id="ResourceURI"></a>
_ResourceURI (string)_

> Uniform Resource Identifier as specified in [RFC 3986], used to identify and
> locate a software artifact. Case sensitive and MUST be case normalized as per
> section 6.2.2.1 of RFC 3986, meaning that the scheme and authority MUST be in
> lowercase.
>
> SHOULD resolve to the artifact, but MAY be unresolvable. It is RECOMMENDED to
> use [Package URL][] (`pkg:`) or [SPDX Download Location][] (e.g.
> `git+https:`).
>
> Example: `"pkg:deb/debian/stunnel@5.50-3?arch=amd64"`.

<a id="TypeURI"></a>
_TypeURI (string)_

> Uniform Resource Identifier as specified in [RFC 3986], used as a
> collision-resistant type identifier. Case sensitive and MUST be case
> normalized as per section 6.2.2.1 of RFC 3986, meaning that the scheme and
> authority MUST be in lowercase.
>
> SHOULD resolve to a human-readable description, but MAY be unresolvable.
> SHOULD include a version number to allow for revisions.
>
> Example: `"https://in-toto.io/Attestation/v1"`.

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
