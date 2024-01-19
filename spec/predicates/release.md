# Predicate type: Release

Type URI: https://in-toto.io/attestation/release

Version 0.1

## Purpose

To authoritatively link from a specific release name and version string in a registry, to the artifact names and hashes that make up that release.

## Use Cases

When receiving a new release version, package registries can publish a release attestation covering the artifact names and hashes that make up that release. This allows consumers of that release version to ensure the artifacts they are consuming have not been tampered with.

If these release attestations are optionally published to a transparency log, package authors (or other interested parties) can monitor when a new version of a package is released.

These use cases are not hypothetical; both of them are the case today for npm's [build provenance feature], which includes a [publish attestation]. Note that while npm calls this a publish attestation, calling it a release attestation better reflects that it's coming from the package registry. Publishing often refers to an author sending content to the registry, as in [PyPI's trusted publishers feature].

## Prerequisites

This predicate depends on the [in-toto attestation framework], as well as the [purl-spec] for identifying packaging ecosystems.

Perhaps surprisingly, this predicate does not depend on [SLSA Provenance], but they are better together. In an ecosystem where some (but not all) packages have SLSA Provenance and you query by artifact hash, you can't tell the difference between an artifact that has been tampered with and one that does not yet have SLSA provenance. If the ecosystem has release attestations, you can give authoritative answers to what artifacts make up a given release version, and what hashes those artifacts should have.

## Model

This predicate is for the final stages of the software supply chain, where consumers are looking for attestations that the software has not been tampered with during distribution.

## Schema

### Fields

- **name, required** string
  - The filename of the artifact as it would appear on disk.

- **purl, required** string (ResourceURI)
  - A purl uniquely identifying a specific release name and version from a package registry.

### Parsing Rules

The purl field MUST be parsed using the [purl-spec]. It MUST include a version (which is OPTIONAL in the [purl-spec]).

## Example

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "http-7.2.16.tgz
      "digest": {
        "sha256": "4faeb1b21ad0612c1553752dffe2ec006020ef3914b0e9ff7315ca77121b79a5"
      }
    }
  ],
  "predicateType": "https://in-toto.io/attestation/release/v0.1",
  "predicate": {
    "purl": "pkg:npm/@angular/http@7.2.16",
  }
}
```

If a release has multiple artifacts that might be consumed separately, the attestation SHOULD have a subject per artifact:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "urllib3-2.1.0.tar.gz",
      "digest": {
        "sha256": "df7aa8afb0148fa78488e7899b2c59b5f4ffcfa82e6c54ccb9dd37c1d7b52d54"
      }
    },
    {
      "name": "urllib3-2.1.0-py3-none-any.whl",
      "digest": {
        "sha256": "55901e917a5896a349ff771be919f8bd99aff50b79fe58fec595eb37bbc56bb3",
      }
    }
  ],
  "predicateType": "https://in-toto.io/attestation/release/v0.1",
  "predicate": {
    "purl": "pkg:pypi/urllib3@2.1.0"
  }
}
```

## Changelog and Migrations

As this is the initial version, no changes or migrations to previous versions. This proposal is a subset of the information in the existing npm [publish attestation], so npm could easily migrate to this specification.

[build provenance feature]: https://github.blog/2023-04-19-introducing-npm-package-provenance/
[publish attestation]: https://github.com/npm/attestation/tree/main/specs/publish/v0.1
[PyPI's trusted publishers feature]: https://github.com/npm/attestation/tree/main/specs/publish/v0.1
[in-toto attestation framework]: ../README.md
[purl-spec]: https://github.com/package-url/purl-spec
[SLSA Provenance]: https://slsa.dev/provenance
