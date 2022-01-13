# Field type definitions

<a id="DigestSet"></a>
_DigestSet (object)_

> Set of alternative cryptographic digests for a single software artifact,
> expressed as a JSON map from algorithm name to lowercase hex-encoded value.
> Usually there is just a single key/value pair, but multiple entires MAY be
> used for algorithm agility.
>
> It is RECOMMENDED to use at least `sha256` for compatibility between producers
> and consumers.
>
> Pre-defined algorithm names: `sha256`, `sha224`, `sha384`, `sha512`,
> `sha512_224`, `sha512_256`, `sha3_224`, `sha3_256`, `sha3_384`, `sha3_512`,
> `shake128`, `shake256`, `blake2b`, `blake2s`, `ripemd160`, `sm3`, `gost`,
> `sha1`, `md5`.
>
> Consumers MUST only accept algorithms that they consider secure and MUST
> ignore unrecognized or unaccepted algorithms. For example, most applications
> SHOULD NOT accept "md5" because it lacks collision resistance.
>
> Two DigestSets SHOULD be considered matching if ANY acceptable field matches.
>
> Examples:
>
> -   `{"sha256": "abcd", "sha512": "1234"}` matches `{"sha256": "abcd"}`
> -   `{"sha256": "abcd"}` does not match `{"sha256": "fedb", "sha512": "abcd"}`
> -   `{"somecoolhash": "abcd"}` uses a non-predefined algorithm

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
> TypeURIs are not registered. The natural namespacing of URIs is sufficient to
> prevent collisions.
>
> Example: `"https://in-toto.io/Statement/v1"`.

<a id="Timestamp"></a>
_Timestamp (string)_

> A point in time, represented as a string in [RFC 3339] format in the UTC time
> zone ("Z").
>
> Example: `"1985-04-12T23:20:50.52Z"`.

[Package URL]: https://github.com/package-url/purl-spec/
[RFC 3339]: https://tools.ietf.org/html/rfc3339
[RFC 3986]: https://tools.ietf.org/html/rfc3986
[SPDX Download Location]: https://spdx.github.io/spdx-spec/package-information/#77-package-download-location-field
[oci_image_id]: https://github.com/opencontainers/image-spec/blob/master/config.md#imageid
