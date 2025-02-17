# ResourceDescriptor field type specification

A size-efficient description of any software artifact or resource (mutable
or immutable).

## Data description

```cddl
ResourceDescriptor = {
  ? name-label => text,
  at-least-one-of-uri-digest-content,
  ? download-location-label => uri-type,
  ? media-type-label => media-type,
  ? annotations-label => annotations-map
}

at-least-one-of-uri-digest-content = (
  uri-label => uri-type,
  ? digest-label => DigestSet,
  ? content-label => content
) // (
  ? uri-label => uri-type,
  digest-label => DigestSet,
  ? content-label => content
) // (
  ? uri-label => uri-type,
  ? digest-label => DigestSet,
  content-label => content
)

content = JC<text, bytes>

name-label              = JC<"name",              0>
uri-label               = JC<"uri",               1>
digest-label            = JC<"digest",            2>
content-label           = JC<"content",           3>
download-location-label = JC<"download-location", 4>
media-type-label        = JC<"mediaType",         5>
annotations-label       = JC<"annotations",       6>

media-type = text

annotations-map = { * annotation-type => any }
annotation-label = JC<text, text / int>
```

## Fields

Though all fields are optional, a ResourceDescriptor MUST specify one of `uri`,
`digest` or `content` at a minimum. Further, a context that uses the
ResourceDescriptor can require one or more fields. For example, a predicate may
require the `name` and `digest` fields. Note that those requirements cannot
override the minimum requirement of one of `uri`, `digest`, or `content`
specified here.

`name` _string, optional_

> Machine-readable identifier for distinguishing between descriptors.
>
> The semantics are up to the producer and consumer. The `name` SHOULD be
> stable, such as a filename, to allow consumers to reliably use the `name`
> as part of their policy.

`uri` _[ResourceURI], optional_

> A URI used to identify the resource or artifact globally.
> This field is REQUIRED unless either `digest` or `content` is set.

`digest` _[DigestSet], optional_

> A set of cryptographic digests of the contents of the resource or artifact.
> This field is REQUIRED unless either `uri` or `content` is set.
>
> When known, the producer SHOULD set this field to denote an immutable
> artifact or resource. The producer and consumer SHOULD agree on acceptable
> algorithms.

`content` _bytes, optional_

> The contents of the resource or artifact.
> This field is REQUIRED unless either `uri` or `digest` is set.
>
> The producer MAY use this field in scenarios where including the contents
> of the resource/artifact directly in the attestation is deemed more
> efficient for consumers than providing a pointer to another location. To
> maintain size efficiency, the size of `content` SHOULD be less than 1KB.
>
> The semantics are up to the producer and consumer. The `uri` or
> `mediaType` MAY be used by the producer as hints for how consumers should
> parse `content`.

`downloadLocation` _[ResourceURI], optional_

> The location of the described resource or artifact, if different from the
> `uri`.
>
> To enable automated downloads by consumers, the specified location SHOULD
> be resolvable.

`mediaType` _string, optional_

> The [MIME Type][] (i.e., media type) of the described resource or artifact.
>
> For resources or artifacts that do not have a standardized MIME type,
> producers SHOULD follow [RFC 6838][] (Sections 3.2-3.4) conventions of
> prefixing types with `x.`, `prs.`, or `vnd.` to avoid collisions with other
> producers.

`annotations` _map <string, value>, optional_

> This field MAY be used to provide additional information or metadata about
> the resource or artifact that may be useful to the consumer when evaluating
> the attestation against a policy.
>
> For maximum flexibility annotations may be any mapping from a field name to
> any JSON value (string, number, object, array, boolean or _null_).
>
> The producer and consumer SHOULD agree on the semantics, and acceptable
> fields and values in the `annotations` map. Producers SHOULD follow the
> same naming conventions for annotation fields as for [extension fields].
>
> In CBOR, the map key MAY be an integer for a more concise representation.
> Negative integers may be used for enumerations agreed upon between producer
> and consumer. Non-negative integers are reserved for in-toto specifications
> to standardize.

## Semantics

Though the ResourceDescriptor allows for a high degree of flexibility,
certain field combinations typically have specific semantics.
For consistency, we RECOMMEND the following:

-   A descriptor that specifies a `digest` is assumed to refer to an
immutable resource or artifact. The `digest` SHOULD match the resource or
artifact specified in one of the `name`, `uri`, `content` or
`downloadLocation` fields. The field that consumers are expected to match
the `digest` against is ultimately determined by the predicate type, and
SHOULD be documented by the predicate specification.
-   A descriptor without a `content` field, is assumed to serve as a
pointer to the resource/artifact.
-   When `uri` and `name` are specified, the scope of `name` is assumed to be
local to the attestation. The `name` SHOULD be compared to the `uri` to match
the descriptor against the resource or artifact referenced universally.

## Examples

Pointer to a file:

```jsonc
{
  "name": "Foo.apk",
  "uri": "git+https://android.googlesource.com/platform/vendor/foo/bar@16244f4e7524d44a8f3060905eaf9190e96e9fb0#prebuilts/Foo/Foo.apk",
  "digest": { "sha256": "7f4714fd..." }
}
```

Pointer to a git repo (with annotations):

```jsonc
{
  "uri": "git+https://github.com/actions/runner",
  "digest": { "sha1": "d61b27b8395512..." },
  "annotations": { "twoPartyReview": false }
}
```

Pointer to another in-toto attestation:

```jsonc
 {
   "name": "gcc_9.3.0-1ubuntu2_amd64.intoto.json",
   "digest": { "sha256": "abcdabcde..." },
   "downloadLocation": "http://example.com/rebuilderd-instance/gcc_9.3.0-1ubuntu2_amd64.intoto.json",
   "mediaType": "application/x.dsse+json"
 }
```

Pointer to build service:
o
```jsonc
{
  "uri": "https://cloudbuild.googleapis.com/GoogleHostedWorker@v1"
}
```

Descriptor for Tekton TaskRun:

```jsonc
{
  "uri": "https://tekton.dev/TaskRun/check/result/report",
  "digest": {
    "sha256": "ec87961dbfe8e7d8a73890c602ac7bd407b80b9e31c326beb9110bdd255f12e6"
  },
  "content": "eyAicmVzdWx0IjogIlNVQ0NFU1MiLCAidGltZXN0YW1wIjogIjE2Nzc4NzIyMzYiLCAic3VjY2Vzc2VzIjogMjIsICJmYWlsdXJlcyI6IDAsICJ3YXJuaW5ncyI6MCB9",
  "mediaType": "application/json"
}
```

[DigestSet]: digest_set.md
[MIME Type]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
[RFC 6838]: https://www.rfc-editor.org/rfc/rfc6838.html#section-3.2
[ResourceURI]: field_types.md#ResourceURI
[extension fields]: ./#parsing-rules
