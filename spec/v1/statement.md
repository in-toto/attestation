# Statement layer specification

The Statement is the middle layer of the attestation, binding it to a
particular subject and unambiguously identifying the types of the
[Predicate].

## Data description

The following CDDL describes both the JSON and CBOR encoding of an in-toto Statement.

```cddl
statement =
{
  _type-label => "https://in-toto.io/Statement/v1",
  subject-label => [ + ResourceDescriptor ],
  predicate-group
}

JC<J,C> = J .feature "json" / C .feature "cbor"

_type-label         = JC<"_type",         0>
subject-label       = JC<"subject",       1>

uri-type        = JC<uri-text, uri-type-choice>
uri-type-choice = uri-text / uri
uri-text        = text
```

The format of the `uri-type` SHOULD be a valid URI according to [RFC3986].

## Fields

The Statement is represented as a [JSON] object with the following fields.
Additional [parsing rules] apply.

`_type` _string ([TypeURI]), required_

> Identifier for the schema of the Statement. Always
> `https://in-toto.io/Statement/v1` for this version of the spec.

`subject` _array of [ResourceDescriptor] objects, required_

> Set of software artifacts that the attestation applies to. Each element
> represents a single software artifact. Each element MUST have `digest` set.
>
> Subjects are assumed to be _immutable_, i.e. the artifacts identified by the
> subject SHOULD NOT change.
>
> The `name` field may be used as an identifier to distinguish this artifact
> from others within the `subject`. Similarly, other ResourceDescriptor fields
> may be used as required by the context. The semantics are up to the producer
> and consumer and they MAY use them when evaluating policy. If the name is not
> meaningful, leave the field unset or use "\_". For example, a
> [SLSA Provenance] attestation might use the name to specify output filename,
> expecting the consumer to only consider entries with a particular name.
> Alternatively, a vulnerability scan attestation might leave name unset because
> the results apply regardless of what the artifact is named.
>
> If set, `name` and `uri` SHOULD be unique within subject.
>
> IMPORTANT: Subject artifacts are matched purely by digest, regardless of
> content type. If this matters to you, please comment on
> [GitHub Issue #28](https://github.com/in-toto/attestation/issues/28)

`predicate-group`

See [Predicate].

[ResourceDescriptor]: resource_descriptor.md
[JSON]: https://www.json.org/json-en.html
[Predicate]: predicate.md
[SLSA Provenance]: https://slsa.dev/provenance
[TypeURI]: field_types.md#TypeURI
[parsing rules]: README.md#parsing-rules
[RFC3986]: https://datatracker.ietf.org/doc/html/rfc3986
