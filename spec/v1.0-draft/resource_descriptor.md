# ResourceDescriptor field type specification

Version v1.0-draft

A size-efficient description of any mutable or immutable software artifact
or resource.

## Schema

```json
{
  "name": "<NAME>",
  "uri": "<RESOURCE URI>",
  "digest": { "<ALGORITHM>": "<HEX VALUE>", ... },
  "representation": { /* object */ },
  "downloadLocation": "<RESOURCE URI>"
  "mediaType": "<MIME TYPE>",
  "annotations": { /* object */ }
}
```

## Fields

Though all fields are optional, a ResourceDescriptor MUST specify at least
one of `uri` or `digest` at a minimum.

`name` _string, optional_

> Human-readable identifier to distinguish the resource or artifact locally.
> The semantics are up to the producer and consumer.

`uri` _[ResourceURI], optional_

> A URI used to identify the resource or artifact. When a `digest` cannot be
> computed for the resource or artifact, the producer SHOULD set this field.

`digest` _[DigestSet], optional_

> A set of crypographic digests of the contents of the resource or artifact.
> When known, the producer SHOULD set this field to denote an immutable
> artifact or resource. The producer and consumer must agree on acceptable
> algorithms.

`representation` _object, optional_

> A representation of the contents of the resource or artifact.
> The producer MAY use this field in scenarios where including the contents
> of the resource/artifact directly in the attestation is deemed more
> efficient for consumers than providing a reference to another location.
>
> The producer and consumer must agree on the semantics and acceptable
> formats for the `representation` object.

`downloadLocation` _[ResourceURI], optional_

> The location of the described resource or artifact, if different from the
> `uri`. To enable automated downloads by consumers, the specified location
> SHOULD be resolvable.

`mediaType` _string, optional_

> The [MIME Type][] (i.e., media type) of the described resource or artifact.

`annotations` _object, optional_

> An object describing annotations to resource or artifact description.
> The producer MAY use this field to provide additional information or
> metadata about the resource or artifact that may be useful to the consumer.
>
> The producer and consumer must agree on the semantics and acceptable
> formats for the `annotations` object.

## Semantics

Though the ResourceDescriptor allows for a high degree of flexibility,
certain field combinations typically have specific semantics.
For consistency, we RECOMMEND the following:

-   A descriptor that specifies a `digest` is assumed to refer to an
immutable resource or artifact.
-   When `uri` AND `digest` are specified, the descriptor is bound to the
`digest` field, and all other fields are considered informational.
-   A descriptor without a `representation`, is assumed to serve as a
pointer to the resource/artifact.

## Examples

Pointer to a local file:

```jsonc
{
  "name": "foo.c",
  "digest": { "sha256": "abc123def456..." }
}
```

Pointer to a remote file:

```jsonc
{
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
   "name": "gcc9.3.0-rebuilderd-attestation",
   "digest": { "sha256": "abcdabcde..." },
   "downloadLocation": "http://example.com/rebuilderd-instance/gcc_9.3.0-1ubuntu2_amd64.att",
   "mediaType": "application/vnd.in-toto+json"
 }
```

Pointer to build service:

```jsonc
{
  "uri": "https://cloudbuild.googleapis.com/GoogleHostedWorker@v1"
}
```

<!-- TODO: Representation of small file -->

[DigestSet]: digest_set.md
[MIME Type]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
[ResourceURI]: scalar_field_types.md#ResourceURI
