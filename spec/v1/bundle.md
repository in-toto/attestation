# Bundle layer specification

An attestation Bundle is a collection of multiple attestations in a single
file. This allows attestations from multiple different points in the software
supply chain (e.g. Provenance, Code Review, Test Result, vuln scan, ...) to
be grouped together, allowing users to make decisions based on all available
information.

**NOTE**: The Bundle is not authenticated as a whole. Instead each individual
attestation is authenticated using signature formats like [DSSE]. As such,
an attacker might be able to delete valid attestations, replay obsolete
attestations, and/or inject invalid or irrelevant attestations in a Bundle
without being detected.

Our Bundle specification aims to reduce the risks of such attacks.
In addition, predicates that follow the [monotonic principle] should not
have any issues with deletion attacks.

## Data structure

Attestation Bundles use [JSON Lines] to store multiple attestations.

-   Each attestation within a Bundle MAY have different signing keys,
    `_type`, `subject`, and/or `predicateType`.
-   Each line within a Bundle SHOULD be an [Envelope]. Consumers MUST ignore
    unrecognized lines.
-   Consumers MUST ignore attestations with unrecognized keys, types,
    subjects, or predicates.
-   Processing of a Bundle MUST NOT depend on the order of the attestations.

## File naming convention

Bundles SHOULD use the suffix `.intoto.jsonl`.

-   A Bundle of attestations relevant for `<filename>` SHOULD be named
    `<filename>.intoto.jsonl`.
-   Attestations in the Bundle MAY have different subjects, but they SHOULD
    all be relevant to that file or its dependencies.

### Example

A package named `foo-1.2.3.tar.gz` with hash `abcd` that was built from git
commit `1234` could have a Bundle name `foo-1.2.3.tar.gz.intoto.jsonl` with
two attestations, one with subject `abcd` and one with subject `1234`.

## Storage convention

The media type `application/vnd.in-toto.bundle` SHOULD be used to denote
a Bundle in arbitrary storage systems.

-   The media type for a Bundles does not indicate an encoding since Bundles
    MUST be encoded as JSON lines, and the encoding of each attestation
    within the Bundle SHOULD be indicated at the [Envelope] layer.
-   The predicate type of individual attestations within the stored Bundle
    SHOULD NOT be indicated in the media type for the Bundle, as this
    information is not authenticated at the Bundle layer.
-   To obtain predicate information that is authenticated, consumers MUST
    download and parse each line in a Bundle separately.

## Example use case

The Fooly app has a CI/CD system which builds the application from source,
runs a NoVulnz vulnerability scan on the results, produces an SPDX SBOM, and
then deploys the app to an app store.

### Build

The Fooly builder builds the app (`fooly.apk` with hash `aaa...`) and
produces a generic [SLSA Provenance].  The Fooly builder also produces a more
detailed attestation that contains all the logs of the build as an in-toto
[Statement] with `predicateType=https://fooly.app/Builder/v1`. The builder
places _both_ of these signed attestations in a new file named
`fooly.apk.intoto.jsonl`.

```jsonl
{ "payloadType": "application/vnd.in-toto+json", "payload": "a...", "signatures": [w...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "b...", "signatures": [w...] }
```

### Vulnerability scanning

The CI/CD system then requests a third-party vulnerability scan on
`fooly.apk`. The vulnerability scanner doesn't use in-toto Statements but
instead uses their own custom `payloadType=application/vnd.novulz+cbor`,
which they put in a [DSSE] envelope. This attestation indicates `fooly.apk`
with hash `aaa...` has 0 critical vulnerabilities and 3 medium
vulnerabilities.

The TestResult is then appended to the contents of `fooly.apk.intoto.jsonl`

```jsonl
{ "payloadType": "application/vnd.in-toto+json", "payload": "a...", "signatures": [w...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "b...", "signatures": [w...] }
{ "payloadType": "application/vnd.novulz+cbor", "payload": "c...", "signatures": [x...] }
```

### SBOM generation

The CI/CD system then generates an SPDX SBOM attestation for `fooly.apk`
with hash `aaa...` using an in-toto Statement with
[`predicateType=https://spdx.dev/Document`](https://github.com/in-toto/attestation/blob/main/spec/predicates/spdx.md)
and appending that to the contents of `fooly.apk.intoto.jsonl`.

```jsonl
{ "payloadType": "application/vnd.in-toto+json", "payload": "a...", "signatures": [w...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "b...", "signatures": [w...] }
{ "payloadType": "application/vnd.novulz+cbor", "payload": "c...", "signatures": [x...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "d...", "signatures": [y...] }
```

### Deployment

Just prior to deployment the CI/CD system checks `fooly.apk` with a policy
engine (providing `fooly.apk.intoto.jsonl` as it does so) to ensure the app
is safe to publish. The policy engine used doesn't understand
`predicateType=https://spdx.dev/Document`, so it is ignored.

Satisfied with the result CI/CD system now deploys `fooly.apk` to the app
store.

### Attestation publishing

Fooly Inc. wants to publish all of the accumulated attestations for evey
published app _except for_ the internal build attestation. The CI/CD system
then iterates through all the attestations, removing the attestation with
`predicateType=https://fooly.app/Builder/v1` and publishes to their website:

```jsonl
{ "payloadType": "application/vnd.in-toto+json", "payload": "a...", "signatures": [w...] }
{ "payloadType": "application/vnd.novulz+cbor", "payload": "c...", "signatures": [x...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "d...", "signatures": [y...] }
```

## Viewing bundles

Large attestations might wind up being difficult for a person to read if
they're serialized to a single line each.  An easy way to make it more
legible for people is to use the `jq` command.

Example:

```shell
$ grep . ~/multiple.intoto.jsonl | jq
{
  "payloadType": "application/vnd.in-toto+json",
  "payload": "a...",
  "signatures": []
}
{
  "payloadType": "application/vnd.in-toto+json",
  "payload": "b...",
  "signatures": []
}
{
  "payloadType": "application/vnd.novulz+cbor",
  "payload": "c...",
  "signatures": []
}
{
  "payloadType": "application/vnd.in-toto+json",
  "payload": "d...",
  "signatures": []
}
```

[DSSE]: https://github.com/secure-systems-lab/dsse
[Envelope]: envelope.md
[JSON Lines]: https://jsonlines.org/
[SLSA Provenance]: https://slsa.dev/provenance
[monotonic principle]: https://github.com/in-toto/attestation/tree/main/spec/v1#parsing-rules
