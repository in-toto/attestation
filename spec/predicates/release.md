# Predicate type: Release

Type URI: https://in-toto.io/attestation/release

Version 1

## Purpose

To authoritatively link from a specific release name and version string in a
package registry, to the artifact names and hashes that make up that release.

## Use Cases

When receiving a new release version, package registries can publish a release
attestation covering the artifact names and hashes that make up that release.
This allows consumers of that release version to ensure the artifacts they are
consuming have not been modified, no matter how many network links or
intermediate caches were used to acquire the artifact.

If these release attestations are optionally published to a transparency log,
package authors (or other interested parties) can monitor when a new version of
a package is released.

The release attestation provides integrity for the release artifacts, but it
does not provide availability as it does not require the registry to serve an
attested release artifact that is later determined to be untrustworthy.

These use cases are not hypothetical; both of them are the case today for npm's
[build provenance feature], which includes a [publish attestation]. Note that
while npm calls this a publish attestation, calling it a release attestation
better reflects that it's coming from the package registry. Publishing often
refers to an author sending content to the registry, as in
[PyPI's trusted publishers feature].

Perhaps surprisingly, this predicate does not depend on [SLSA Provenance], but
they are better together. In an ecosystem where some (but not all) packages
have SLSA Provenance and you query by artifact hash, you can't tell the
difference between an artifact that has been tampered with and one that does
not yet have SLSA provenance. If the ecosystem has release attestations, you
can give authoritative answers to what artifacts make up a given release
version, and what hashes those artifacts should have.

## Prerequisites

This predicate depends on the [in-toto Attestation Framework], as well as the
[purl-spec] for identifying packaging ecosystems.

## Model

This predicate is for the final stages of the software supply chain, where
consumers are looking for attestations that the software has not been tampered
with during distribution.

If a registry supports immutable releases, there SHOULD be one release
attestation for a given `predicate.purl`.

Even if the registry does not support immutable releases, the attestation
subject SHOULD include all the artifacts associated with the release at that
time; otherwise it will be unclear if an artifact was later removed from a
release. For example, if a release version has an initial release attestation
with artifact A, and then later has a release attestation with only artifact B,
that SHOULD be interpreted as the release version now only containing artifact
B.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/release/v1",
  "predicate": {
    "purl": <ResourceURI>,
    "releaseId": "..."
  }
}
```

### Fields

- **`subject.name`** string
  - The filename of the artifact as it would appear on disk.

- **`predicate.purl`, required** string ([ResourceURI])
  - A purl uniquely identifying a specific release name and version from a
    package registry.

- **`predicate.releaseId`** string
  - Stable identifier for a release; this should remain unchanged between
    release versions (e.g. it's associated with urllib3, not urllib3 v2.1.0).
    This will allow users to confirm that a release has moved to a new name,
    and prevent confusion if the old name is re-used. This could be an
    automatically incrementing database key or a randomly generated UUID.

### Parsing Rules

The purl field MUST be parsed using the [purl-spec]. It MUST include a purl
`version` (which is OPTIONAL in the [purl-spec]). It SHOULD NOT include purl
`qualifiers` or `subpath`, unless the `type` requires them to uniquely identify
a release (as a counter-example `type:oci` would include `qualifiers`).

## Example

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "http-7.2.16.tgz",
      "digest": {
        "sha256": "4faeb1b2..."
      }
    }
  ],
  "predicateType": "https://in-toto.io/attestation/release/v1",
  "predicate": {
    "purl": "pkg:npm/@angular/http@7.2.16",
    "releaseId": 1234567890
  }
}
```

If a release has multiple artifacts that might be consumed separately, the
attestation SHOULD have a subject per artifact:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "urllib3-2.1.0.tar.gz",
      "digest": {
        "sha256": "df7aa8af..."
      }
    },
    {
      "name": "urllib3-2.1.0-py3-none-any.whl",
      "digest": {
        "sha256": "55901e91...",
      }
    }
  ],
  "predicateType": "https://in-toto.io/attestation/release/v1",
  "predicate": {
    "purl": "pkg:pypi/urllib3@2.1.0"
  }
}
```

## Changelog and Migrations

As this is the initial version, no changes or migrations to previous versions.
This proposal is a subset of the information in the existing npm
[publish attestation], so npm could easily migrate to this specification.

[build provenance feature]:
https://github.blog/2023-04-19-introducing-npm-package-provenance/
[publish attestation]:
https://github.com/npm/attestation/tree/main/specs/publish/v0.1
[PyPI's trusted publishers feature]: https://docs.pypi.org/trusted-publishers/
[SLSA Provenance]: https://slsa.dev/provenance
[in-toto Attestation Framework]: ../README.md
[purl-spec]: https://github.com/package-url/purl-spec
[ResourceURI]: ../v1/field_types.md#resourceuri
