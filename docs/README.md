# Introduction

The **in-toto Attestation Framework**

1.  defines a standard [format](../spec/v1/) for attestations which
    bind subjects, the artifacts being described, to arbitrary authenticated
    metadata about the artifact
2.  provides a set of [pre-defined predicates](../spec/predicates/) for
    communicating authenticated metadata throughout and across software supply
    chains

The framework was designed to support use of the data by automated policy
engines, such as [in-toto-verify] and [Binary Authorization]. However, any
system that creates or consumes verifiable claims about how a piece of software
is produced can benefit from the in-toto Attestation Framework.

## Project goals

-   Standardize artifact metadata without being specific to the producer or
    consumer. This way CI/CD pipelines, vulnerability scanners, and other
    systems can generate a single set of attestations that can be consumed by
    anyone.
-   Make it possible to write automated policies that take advantage of
    structured information.
-   Fit within the [SLSA Framework][SLSA].
-   Support an ecosystem of verifiable metadata about software artifacts to
    improve software supply chain security.

## Example uses

Attestations can be used to provide data for a variety of policy decisions
including, but not limited to, the following high-level examples:

-   [Provenance][SLSA Provenance]: GitHub Actions attests to the fact that it
    built a container image with digest "sha256:87f7fe…" from git commit
    "f0c93d…" in the "main" branch of "https://github.com/example/foo".
-   Test result: GitHub Actions attests to the fact that the npm tests passed on
    git commit "f0c93d…".
-   Vulnerability scan: Google Container Analysis attests to the fact that no
    vulnerabilities were found in container image "sha256:87f7fe…" at a
    particular time.
-   Policy decision: Binary Authorization attests to the fact that container
    image "sha256:87f7fe…" is allowed to run under GKE project "example-project"
    within the next 4 hours, and that it used the four attestations above and as
    well as the policy with sha256 hash "79e572" to make its decision.

## Detailed documentation

-   [background on the framework and its requirements](background.md)
-   [motivating use case](motivating_use_case.md)
-   [guidelines for creating new attestation schemas](new_predicate_guidelines.md)
-   [validation model](validation.md)
-   [ideas for future schemas](schema_ideas.md)
-   [protobuf definitions](protos.md)
-   [testing the implementations](testing.md)

[Binary Authorization]: https://cloud.google.com/binary-authorization
[SLSA Provenance]: https://slsa.dev/provenance
[SLSA]: https://github.com/slsa-framework/slsa
[in-toto-verify]: https://github.com/in-toto/in-toto#verification
