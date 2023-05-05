# Predicate type: Link

Type URI: https://in-toto.io/attestation/link/v0.3

Deprecated Type URI: https://in-toto.io/Link/*

Version: 0.3

## Purpose

A generic attestation type with a schema isomorphic to [in-toto specification].
This allows existing in-toto users to make minimal changes to upgrade to the new
attestation format.

Depending on the context, a more specific predicate type such as [Provenance]
may be more appropriate.

## Prerequisites

Understanding of [in-toto specification] and the in-toto attestation framework.

## Model

Every link attestation corresponds to the execution of one step in the software
supply chain. The `subject` field corresponds to the products of the operation
while `materials` indicates the inputs to the step. The `name` field identifies
the step the attestation corresponds to and the `command` field optionally
describes the command that was run. Link attestations allow for recording
arbitrary but relevant information in the opaque `environment` and `byproducts`
fields.

For every step described in the layout of the software supply chain, one (or
more, depending on the threshold) signed link attestations must be presented.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/link/v0.3",
  "predicate": {
    "name": "...",
    "command": [ ... ],
    "materials": [<ResourceDescriptor>, ...],
    "byproducts": { ... },
    "environment": { ... }
  }
}
```

### Fields

The `predicate` has the same schema as the link's `signed` field in
[in-toto specification] except:

-   `predicate._type` is omitted. `predicateType` serves the same purpose.
-   `predicate.products` is omitted. `subject` serves the same purpose.
-   `predicate.materials` is updated to a list of `ResourceDescriptors`.
    Each `ResourceDescriptor` entry MUST include `name` and `digest`.

`name`, _string, required_

Name of the step. When used with an in-toto layout as described in the
[in-toto specification], this field is used when identifying link metadata to
verify for a step.

`command`, _list of strings, optional_

Contains the command and its arguments executed during the step. While the field
is required, it may be empty.

`materials`, _list of [ResourceDescriptor] objects, optional_

List of artifacts that make up the materials of the step. The `name` and
`digest` fields of each entry MUST be set. [ITE-4] artifact types may be
identified using the `uri` or `mediaType` fields instead of overloading the
`name` field.

`byproducts`, _object, optional_

An opaque dictionary that contains additional information about the step
performed. Consult the [in-toto specification] for how the verification workflow
may use it.

`environment`, _object, optional_

An opaque dictionary that contains information about the environment in which
the step was carried out. Consult the [in-toto specification] for how the
verification workflow may use it.

### Converting to old-style links

A Link predicate may be converted into an [in-toto specification] link as
follows:

-   Set `link._type` to `"link"`.
-   Set `link.name`, `link.environment`, and `link.byproducts` using
    `predicate.name`, `predicate.environment`, and `predicate.byproducts`
    respectively.
-   Set `link.products` to be a map from `subject[*].name` to
    `subject[*].digest`.
-   Set `link.materials` to be a map from `predicate.materials[*].name` to
    `predicate.materials[*].digest`.

In Python:

```python
def convert(statement):
    assert statement.predicateType == 'https://in-toto.io/attestation/link/v0.3'
    link = {}
    link['_type'] = 'link'
    link['name'] = statement.predicate['name']
    link['byproducts'] = statement.predicate['byproducts']
    link['environment'] = statement.predicate['environment']
    link['materials'] = {s['name'] : s['digest'] for s in statement.predicate['materials']}
    link['products'] = {s['name'] : s['digest'] for s in statement.subject}
    return link
```

## Version History

-   0.3: Updated `materials` to use a list of `ResourceDescriptor` objects.
    Reverted `command` to a list of strings to match the original link
    specification.
-   0.2: Removed `_type` and `products`. Defined conversion rules.
-   0.1: Initial version as described in [in-toto specification].

<!-- TODO: Fix link-->

[in-toto specification]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md
[ResourceDescriptor]: ../v1.0/resource_descriptor.md
[Provenance]: provenance.md
[ITE-4]: https://github.com/in-toto/ITE/blob/master/ITE/4/README.adoc
