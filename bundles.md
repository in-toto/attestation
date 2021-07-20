# Attestation Bundle

An attestation bundle is a collection of multiple attestations in a single file.
This allows attestations from multiple different points in the software supply
chain (e.g. Provenance, Code Review, Test Result, vuln scan, ...) to be grouped
together, allowing users to make decisions based on all available information.

**NOTE**: Attestation Bundles themselves are **not authenticated** instead each
individual attestation is authenticated
([using DSSE](https://github.com/secure-systems-lab/dsse)). As such, an attacker
might be able to _delete_ an attestation without being detected.  Predicates that
follow [the monotonic principle](spec/README.md#parsing-rules) should not have any
issues with this behavior.

## Data structure

Attestation Bundles use [JSON Lines](https://jsonlines.org/) to store multiple
[DSSEs](https://github.com/secure-systems-lab/dsse).

*  Each line MUST contain a _single_ DSSE.
*  Each line MAY contain any DSSE `payloadType`
*  in-toto attestations (`payloadType` == `application/vnd.in-toto+json`) MAY use
   any ` _type`/`predicateType`
*  [in-toto Statements](spec/README.md#statement)
   in a Bundle MAY reference different Subjects
*  Consumers MUST ignore any attestations whose `payloadType`, `_type`, or `predicateType`
   they do not understand.
*  Attestations MAY be signed by different keys
*  New attestations MAY be added to existing bundles
*  Processing of a bundle MUST NOT depend on the order of the attestations.

## File naming convention

* Attestation Bundles SHOULD use the `.intoto.jsonl` extension.
* Bundles that concern a single artifact SHOULD name the bundle file
  `<artifact filename>.intoto.jsonl`.
* Bundles that concern multiple artifacts SHOULD name the bundle file
  `intoto.jsonl`.

## Example Use Case

The Fooly app has a CI/CD system which builds the application from source, runs a
NoVulnz vulnerability scan on the results, produces an SPDX SBOM, and then deploys the
app to an app store.

### Build

The Fooly builder builds the app (`fooly.apk` with hash `aaa...`) and produces a generic
[in-toto Provenance](spec/predicates/provenance.md).  The Fooly builder also
produces a more detailed attestation that contains all the logs of the build as an
in-toto Statement with `predicateType=https://fooly.app/Builder/v1`.  The builder places
_both_ of these signed attestations in a new file named `fooly.apk.intoto.jsonl`.

```jsonl
{ "payloadType": "application/vnd.in-toto+json", "payload": "a...", "signatures": [w...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "b...", "signatures": [w...] }
```

### Vulnerability Scanning

The CI/CD system then requests a third-party vulnerability scan on `fooly.apk`.  The
vulnerability scanner doesn't use in-toto Statements but instead uses their own custom
`payloadType=application/vnd.novulz+cbor`, which they put in a DSSE envelope. This
attestation indicates `fooly.apk` with hash `aaa...` has 0 critical vulnerabilities and
3 medium vulnerabilities.

The TestResult is then appended to the contents of `fooly.apk.intoto.jsonl`

```jsonl
{ "payloadType": "application/vnd.in-toto+json", "payload": "a...", "signatures": [w...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "b...", "signatures": [w...] }
{ "payloadType": "application/vnd.novulz+cbor", "payload": "c...", "signatures": [x...] }
```

### SBOM Generation

The CI/CD system then generates an SPDX SBOM attestation for `fooly.apk` with hash
`aaa...` using an in-toto Statement with
[`predicateType=https://spdx.dev/Document`](https://github.com/in-toto/attestation/blob/main/spec/predicates/spdx.md)
and appending that to the contents of `fooly.apk.intoto.jsonl`.

```jsonl
{ "payloadType": "application/vnd.in-toto+json", "payload": "a...", "signatures": [w...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "b...", "signatures": [w...] }
{ "payloadType": "application/vnd.novulz+cbor", "payload": "c...", "signatures": [x...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "d...", "signatures": [y...] }
```

### Deployment

Just prior to deployment the CI/CD system checks `fooly.apk` with a policy engine
(providing `fooly.apk.intoto.jsonl` as it does so) to ensure the app is safe to publish.
The policy engine used doesn't understand `predicateType=https://spdx.dev/Document`, so
it is ignored.

Satisfied with the result CI/CD system now deploys `fooly.apk` to the app store.

### Attestation Publishing

Fooly Inc. wants to publish all of the accumulated attestations for evey published app
_except for_ the internal build attestation. The CI/CD system then iterates through all
the attestations, removing the attestation with
`predicateType=https://fooly.app/Builder/v1` and publishes to their website:

```jsonl
{ "payloadType": "application/vnd.in-toto+json", "payload": "a...", "signatures": [w...] }
{ "payloadType": "application/vnd.novulz+cbor", "payload": "c...", "signatures": [x...] }
{ "payloadType": "application/vnd.in-toto+json", "payload": "d...", "signatures": [y...] }
```
