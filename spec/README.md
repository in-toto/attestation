# in-toto Attestation Framework Spec

Latest version: [v1.1]

An **in-toto attestation** is authenticated metadata about one or more
software artifacts[^1]. The intended consumers are automated policy engines,
such as [in-toto-verify] and [Binary Authorization].

It has four layers that are independent but designed to work together:

-   [Predicate]: Contains arbitrary metadata about a subject artifact, with a
    type-specific schema.
-   [Statement]: Binds the attestation to a particular subject and
    unambiguously identifies the types of the predicate.
-   [Envelope]: Handles authentication and serialization.
-   [Bundle]: Defines a method of grouping multiple attestations together.

The following diagram visualises the relationships between the envelope, statement and predicate layers.

<img src="../images/envelope_relationships.png" alt="Relationships between the envelope, statement and predicate layers" width="600">

For future edits, we provide the [source](../images/envelope_relationships.excalidraw) of this diagram.

The [validation model] provides pseudocode showing how these layers fit
together. See the [documentation] for more background and examples.

## Tagged Releases

The latest [tagged release] version matches the [SemVer](https://semver.org)
MAJOR.MINOR version of the Attestation Framework spec.

Backwards-compatible semantic updates to the spec (except predicates) are
indicated through new tagged MINOR version releases.
We use new tagged PATCH version releases to indicate updates to predicate
specifications and/or backwards-compatible changes to the language bindings.

### Examples

-   Attestation Framework tagged release v1.0.2 (PATCH version) incorporates
    refinements to the predicate specification process, a new predicate type,
    and a small patch to the Golang language bindings. None of these changes
    affects the semantics of the core spec. The `_type` of a `Statement` is
    still `https://in-toto.io/Statement/v1`.

-   Tagged release v1.1.0 (MINOR version) generalizes the semantics of the
    `DigestSet` field type to support any type of immutable identifier.
    This change is backwards comptabile because cryptographic digests are
    strongly recommended to achieve immutability, so any implementations that
    only support cryptographic `DigestSet` still meet the modified semantics.
    The `_type` of a `Statement` is still `https://in-toto.io/Statement/v1`
    but a new entry in the `v1` CHANGELOG is added.

-   Tagged release v2.0.0 (MAJOR version) changes the meaning of the
    `predicateType` field. A new `v2` directory is added to `/spec` and the
    `_type` of a `Statement` becomes `https://in-toto.io/Statement/v2`.

## Keywords

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in all documents under
this specification are to be interpreted as described in [RFC 2119].

[^1]: This is compatible with the [SLSA Attestation Model].

[Binary Authorization]: https://cloud.google.com/binary-authorization
[Bundle]: v1/bundle.md
[Envelope]: v1/envelope.md
[Predicate]: v1/predicate.md
[RFC 2119]: https://www.rfc-editor.org/rfc/rfc2119
[SLSA Attestation Model]: https://slsa.dev/attestation-model
[Statement]: v1/statement.md
[documentation]: ../docs
[in-toto-verify]: https://github.com/in-toto/in-toto#verification
[tagged release]: https://github.com/in-toto/attestation/releases
[v1.1]: v1/README.md
[validation model]: ../docs/validation.md
