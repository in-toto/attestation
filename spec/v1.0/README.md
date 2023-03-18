# Specification for in-toto attestation layers

Version: v1.0

Index:

-   [Envelope]: Handles authentication and serialization.
-   [Statement]: Binds the attestation to a particular subject and
    unambiguously identifies the types of the predicate.
-   [Predicate]: Contains arbitrary metadata about a subject artifact, with a
    type-specific schema.
-   [Bundle]: Defines a method of grouping multiple attestations together.

## Common field types

As part of this specification, we also define a set of common [field
types] for any layer of an in-toto attestation.

## Parsing rules

The following rules apply to [Statement] and [Predicates] that opt-in to this
model.

-   **Unrecognized fields:** Consumers MUST ignore unrecognized fields unless
    otherwise noted in the predicate specification. This is to allow minor
    version upgrades and extension fields. Ignoring fields is safe due to the
    monotonic principle.

-   **Versioning:** Each type has a [SemVer2](https://semver.org) version
    number and the [TypeURI] reflects the major version number. A message is
    always semantically correct, but possibly incomplete, when parsed as any
    other version with the same major version number and thus the same
    [TypeURI]. Minor version changes always follow the monotonic principle.
    NOTE: 0.X versions are considered major versions.

-   **Extension fields:** Producers MAY add extension fields to any JSON
    object. Extension fields SHOULD use names that are unlikely to collide
    with names used by other orgs. Producers MAY file an issue/PR to document
    extension fields they're using. Any consumer MAY parse and use these
    extensions if desired.

    Field names SHOULD avoid special characters like `.` and `$` as these
    can make querying these fields in some databases more difficult.

    The presence or absence of the extension field MUST NOT influence the
    meaning of any other field, and the field MUST follow the monotonic
    principle.

-   **Monotonic:** A policy is considered monotonic if ignoring an
    attestation, or a field within an attestation, will never turn a DENY
    decision into an ALLOW. A predicate or field follows the monotonic
    principle if the expected policy that consumes it is monotonic.
    Consumers SHOULD design policies to be monotonic. Example: instead of
    "deny if a 'has vulnerabilities' attestation exists", prefer "deny
    unless a 'no vulnerabilities' attestation exists".

See [versioning rules](../versioning.md) for details and examples.

[Bundle]: bundle.md
[Envelope]: envelope.md
[Predicate]: predicate.md
[Statement]: statement.md
[TypeURI]: field_types.md#TypeURI
[field types]: field_types.md
