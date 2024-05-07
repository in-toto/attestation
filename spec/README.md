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

The source of this diagram can be found [here](../images/envelope_relationships.excalidraw).

The [validation model] provides pseudocode showing how these layers fit
together. See the [documentation] for more background and examples.

## Tagged Releases

The latest [tagged release] version matches the [SemVer](https://semver.org)
MAJ.MIN version of the Attestation Framework spec.

Backwards-compatible semantic updates to the spec (except predicates) are
indicated through new tagged minor version releases.
We use new tagged patch version releases to indicate updates to predicate
specifications and/or backwards-compatible changes to the language bindings.

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
