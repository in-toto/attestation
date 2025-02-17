# Field type specifications

The in-toto Attestation Framework defines several field types that may be of
common interest:

-   [DigestSet]
-   [ResourceURI]
-   [TypeURI]
-   [Timestamp]
-   [ResourceDescriptor]

## ResourceURI

Uniform Resource Identifier as specified in [RFC 3986][], used to identify
and locate any resource, service, or software artifact.

**Format:**

A ResourceURI is represented as a case sensitive _string_ and MUST be case
normalized as per section 6.2.2.1 of RFC 3986, meaning that the scheme and
authority MUST be in lowercase.

> SHOULD resolve to the artifact, but MAY be unresolvable. It is RECOMMENDED
> to use [Package URL][] (`pkg:`) or [SPDX Download Location][] (e.g.
> `git+https:`).

**Example:**

`"pkg:deb/debian/stunnel@5.50-3?arch=amd64"`.

## TypeURI

Uniform Resource Identifier as specified in [RFC 3986][], used as a
collision-resistant type identifier.

**Format:**

A TypeURI is represented as a case sensitive _string_ and MUST be case
normalized as per section 6.2.2.1 of RFC 3986, meaning that the scheme and
authority MUST be in lowercase.

> SHOULD resolve to a human-readable description, but MAY be unresolvable.
> SHOULD include a version number to allow for revisions.
>
> TypeURIs are not registered. The natural namespacing of URIs is sufficient
> to prevent collisions.

**Example:**

`"https://in-toto.io/Statement/v1"`.

## Timestamp

A point in time.

**Format:**

```cddl
Timestamp = JC<text, text / tdate>
```

A timestamp is represented as a _string_ and MUST be in [RFC 3339][] format
in the UTC timezone ("Z").

In CBOR, a timestamp SHOULD be use the 0 tag to indicate the formatting, see [Section 3.4.1 of RFC8949].
Note that tag 0 further refines RFC3339 to have a required `"T"` between date and time, following [Section 3.3 of RFC4287].

**Example:**

`"1985-04-12T23:20:50.52Z"`.

[DigestSet]: digest_set.md
[Package URL]: https://github.com/package-url/purl-spec/
[RFC 3339]: https://tools.ietf.org/html/rfc3339
[RFC 3986]: https://tools.ietf.org/html/rfc3986
[ResourceDescriptor]: resource_descriptor.md
[ResourceURI]: #resourceuri
[SPDX Download Location]: https://spdx.github.io/spdx-spec/v2.3/package-information/#77-package-download-location-field
[Timestamp]: #timestamp
[TypeURI]: #typeuri
[Section 3.4.1 of RFC8949]: https://www.rfc-editor.org/rfc/rfc8949.html#name-standard-date-time-string
[Section 3.3 of RFC4287]: https://www.rfc-editor.org/rfc/rfc4287#section-3.3
