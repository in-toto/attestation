# Predicate type: Safety Context

Type URI: https://in-toto.io/attestation/safety-context/v0.1 (proposed)

Version: 0.1 (draft)

Predicate Name: safety-context

## Purpose

Bind an artifact to a safety assessment in an explicit system and intended-use
context, including assessment outcome, optional classification, review status,
and evidence. The vocabulary and URI in this proposal are subject to review.

Software identity alone cannot express the safety significance of its use in a
particular system. The same communications library may have different safety
relevance in a medical control system and an entertainment application.

## Use Cases

A producer publishes a safety context sidecar alongside an SBOM. A downstream
integrator needs to determine whether the assessment applies to the artifact
and configuration being integrated, whether it has been reviewed, and which
evidence supports the assessment.

A second producer publishes equivalent assessment data within a native SPDX
or CycloneDX document. Consumers need the same policy questions to remain
answerable while the location of the underlying data changes.

### Existing predicates and formats

The gaps below concern shared semantics, not whether generic formats can
encode arbitrary JSON.

| Existing predicate or format | Scope and remaining gap |
| --- | --- |
| [Reference] | Binds external documents to an artifact, but does not define safety assessment outcomes, classifications, review states, or context semantics. A separately standardized sidecar referenced through Reference remains a viable alternative. |
| [SCAI] | Supports attributes, conditions, targets, and evidence, but does not prescribe common safety vocabulary. A safety profile of SCAI could address this gap and must be evaluated alongside a dedicated predicate. |
| [VEX], including OpenVEX | Communicates vulnerability status for products. Its vulnerability-specific vocabulary does not by itself define a general system safety assessment. |
| [SLSA Provenance] | Describes how an artifact was built. It does not define its safety relevance in a downstream operating context. |

A dedicated predicate is proposed to give consumers one explicit contract for
context, outcome, classification, review, and evidence. Acceptance depends on
whether that contract offers sufficient benefit over a standardized sidecar
with Reference, or a safety profile of SCAI.

### Policy questions

-   Does this assessment refer to the exact artifact under consideration?
-   Does its system revision, intended use, and assumptions match deployment?
-   Is the outcome relevant, not relevant, or undetermined in that context?
-   Does the classification use a policy-recognized scheme?
-   Has the assessment completed the required review?
-   Are the evidence and source documents available and digest-valid?

An assessment is not a declaration of certification or general safety.
A completed review is not an approval decision. Local policy determines
whether an assessment is sufficient for a specific use.

## Prerequisites

The in-toto Attestation Framework, artifact digest computation, and processes
for producing and reviewing safety assessments are required. Consumers need
an explicit trust policy for attestation signers and reviewers, and support
for the source document versions they interpret.

## Model

One statement describes one assessment of exactly one artifact in one context.
The statement subject identifies the artifact by digest. The predicate records
the contextual assessment. Evidence supports the assessment; sources identify
documents from which the contextual assessment was derived.

A producer can publish an assessment independently of a full SBOM. Native SPDX
or CycloneDX safety data can also be identified as sources when available.
The predicate does not redefine those standards' internal safety models.

## Schema

The [JSON Schema] validates the predicate payload. The statement itself must
also satisfy the in-toto Statement v1 requirements and the subject constraint
above. Schema validation does not establish evidence integrity, trust,
context applicability, or correct mapping from source documents.

### Parsing Rules

This proposal adopts the framework's [standard parsing rules], including the
monotonic principle.

-   Consumers must ignore unrecognized fields. Required policy decisions must
  depend on recognized fields with understood semantics.
-   A missing field or unsupported value must not be interpreted as a favorable
  outcome. The three outcome values are distinct from approval or acceptance.
-   Unsupported predicate versions cannot satisfy a policy requiring this schema.
-   Timestamps must use RFC 3339 UTC syntax ending in `Z`.
-   Consumers must verify the subject digest against the artifact being assessed.
-   Document digests cover the exact referenced bytes, before parsing. Consumers
  must support and verify the declared digest algorithms before relying on
  those documents. Multiple declared digests must all verify.
-   A URL or document selector alone is insufficient for an integrity decision.
-   Context matching must consider the system revision, intended use, and
  assumptions. An artifact digest match alone does not establish applicability.
-   Consumers must not infer review authorization from the `reviewer` field.
  Authorization must be established under the consumer's trust policy.
-   When native data is used, a producer must document its source mapping.
  Unmapped semantics or conflicting sources must not be silently treated as
  equivalent or as a successful policy result.

### Fields

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `assessmentId` | URI string | Yes | Identity of this assessment instance. |
| `assessedAt` | UTC timestamp | Yes | Time this assessment was completed by its producer. |
| `context` | Object | Yes | Context in which the outcome applies. |
| `context.system` | Resource object | Yes | System configuration or revision, identified by an immutable or version-specific `uri`, with optional `name`. |
| `context.intendedUse` | String | Yes | Intended use assessed by the producer. |
| `context.assumptions` | Array of strings | Yes | Explicit assumptions; an empty list means none are stated. |
| `outcome` | Object | Yes | Assessment result and rationale. |
| `outcome.value` | String | Yes | `relevant`, `notRelevant`, or `undetermined`. These express safety relevance in the stated context. |
| `outcome.rationale` | String | Yes | Explanation supporting the result, including uncertainty. |
| `classification` | Object | No | An optional scheme-specific classification. Omission means no classification is asserted. |
| `classification.scheme` | URI string | If classification is present | Identifies the classification scheme and version. |
| `classification.value` | String | If classification is present | Value defined by that scheme; consumers must understand the scheme before interpreting it. |
| `review` | Object | Yes | Review state for the assessment. |
| `review.status` | String | Yes | `pending` or `completed`; neither is an approval or certification result. |
| `review.reviewedAt` | UTC timestamp | For completed review | Time the review was completed; must not precede `assessedAt`. |
| `review.reviewer` | Resource object | For completed review | Reviewer identity as a `uri`, with optional `name`. |
| `evidence` | Array of document objects | Yes | Supporting evidence; empty means no external evidence is supplied. |
| `sources` | Array of source objects | Yes | Derivation sources; empty means this predicate is the primary assessment record. |

Every document object contains a `uri`, `digest` map, and `mediaType`, plus
an optional `name`. Digest values use hexadecimal encoding. This draft uses
a constrained subset of the framework's [ResourceDescriptor] fields.

Each source object contains:

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `format` | String | Yes | `sidecar`, `spdx`, or `cyclonedx`. |
| `version` | String | Yes | Exact source format or schema version, including prerelease status. |
| `document` | Document object | Yes | Digest-bound source document. |
| `selector` | String | For native sources | Source-format-specific element identifier or selector; interpretation must be specified by the mapping. |

## Example

The following complete statement illustrates a sidecar-derived assessment.
Names, URLs, classifications, source versions, and digests are synthetic.
The classification does not correspond to an industry assurance level.
The example source version is not a claim of compatibility with a vendor's
sidecar format. No signatures or certification are implied.

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "device-controller.bin",
      "digest": {
        "sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
      }
    }
  ],
  "predicateType": "https://in-toto.io/attestation/safety-context/v0.1",
  "predicate": {
    "assessmentId": "https://example.org/assessments/device-controller-001",
    "assessedAt": "2026-11-06T18:00:00Z",
    "context": {
      "system": {
        "uri": "https://example.org/systems/infusion-device/revisions/1",
        "name": "Example infusion device, revision 1"
      },
      "intendedUse": "Control infusion delivery in the assessed system configuration.",
      "assumptions": [
        "The assessed configuration uses the documented communications service."
      ]
    },
    "outcome": {
      "value": "relevant",
      "rationale": "The communications service participates in a control path whose failure can affect infusion delivery."
    },
    "classification": {
      "scheme": "https://example.org/safety-classification/v1",
      "value": "requiresSafetyReview"
    },
    "review": {
      "status": "completed",
      "reviewedAt": "2026-11-06T18:30:00Z",
      "reviewer": {
        "uri": "https://example.org/teams/safety-review"
      }
    },
    "evidence": [
      {
        "uri": "https://example.org/evidence/control-path-analysis.json",
        "digest": {
          "sha256": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
        },
        "mediaType": "application/json"
      }
    ],
    "sources": [
      {
        "format": "sidecar",
        "version": "example-1",
        "document": {
          "uri": "https://example.org/metadata/safety-context.json",
          "digest": {
            "sha256": "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
          },
          "mediaType": "application/json"
        }
      }
    ]
  }
}
```

A machine-readable copy is available as [Example JSON].

## Changelog and Migrations

### Initial proposal

This is the first proposed version. It has not been vetted or assigned an
official predicate URI.

### From sidecars to native safety data

Current deployments can publish SPDX or CycloneDX documents alongside a
separate safety context sidecar. This predicate records the assessment and
references that sidecar using `sources`.

When a selected SPDX or CycloneDX version supports the necessary native
safety data, the producer can reference the native document instead. The
source document version and selected elements must be explicit; an adapter
must document which native fields map to the outcome, context,
classification, review, and evidence.

| Stage | Safety data location | Predicate source |
| --- | --- | --- |
| Sidecar | Safety context document accompanying an SBOM | Digest-bound sidecar reference |
| Native | Safety data embedded in a supported SPDX or CycloneDX document | Digest-bound native document reference and element selector |

Migration must preserve assessment meaning and traceability. Where the source
models are not semantically equivalent, producers must issue a new assessment
or retain the sidecar rather than inventing an equivalence. Old sidecar
references remain valid for their original artifact and context.

Changes to the artifact, system configuration, intended use, or assumptions
require an applicability check and may require reassessment. Existing
assessments must not silently carry over to a new artifact digest or context.
A reassessment uses a new assessment identifier. Assessment identifiers alone
do not establish freshness or resolve conflicting assessments.

[Reference]: reference.md
[SCAI]: scai.md
[VEX]: https://github.com/openvex/spec/blob/main/OPENVEX-SPEC.md
[SLSA Provenance]: https://slsa.dev/spec/v1.2/build-provenance
[JSON Schema]: safety-context.schema.json
[Example JSON]: safety-context.example.json
[ResourceDescriptor]: ../v1/resource_descriptor.md
[standard parsing rules]: ../v1/README.md#parsing-rules
