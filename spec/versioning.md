# Versioning Rules

## Objective

This document explains how to update the version number when making changes to
the spec. For regular consumers of the spec, see
[README](README.md#versioning-and-extensions).

## Versioning rules

We adhere to [Semantic Versioning 2.0.0](https://semver.org), i.e.
MAJOR.MINOR.PATCH. (PATCH isn't really needed, but we do it just so we're SemVer
compliant.) This way implementations can state specifically which version they
support.

Any message MUST be semantically correct when parsed as any other version with
the same major number. That is:

-   Parsing an older message as a newer minor/patch version MUST be semantically
    equivalent to parsing it as the original version. (All fields present in
    both versions have the same meaning, and the absence of any new fields has
    no semantic meaning.)
-   Parsing a newer message as an older minor/patch version MUST be semantically
    equivalent to parsing it as the newer version with all new fields removed.
    (All fields present in both version have the same meaning, regardless of the
    presence or absence of the new fields.)

The type ID only includes the major version. This allows implementations to
support any major version without having to parse the URI to pull out the
version. (Parsing is possible since we don't mandate the format of the URI.)
Implementations just use the URI as an opaque string.

The advantage of having minor versions is that we can add new information
without requiring consumers to update.

## Examples

-   Provenance 1.1.0 (minor version) adds a new `buildFinished` timestamp. This
    is OK because the absence of the `buildFinished` has no semantic meaning and
    the presence or absence of `buildFinished` does not affect the rest of the
    message. The type ID is still `https://in-toto.io/Provenance/v1`.

-   Provenance 2.0.0 (major version) adds a new `recipe.extraArgs` field. This
    is required because the meaning of `recipe` changes if `extraArgs` is
    present. Parsing it as 1.x.y will result in an incorrect interpretation.

-   Provenance 3.0.0 (major version) modifies the meaning of `builder.id`. This
    is required because an existing field was modified.
