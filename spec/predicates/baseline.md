# Predicate type: Baseline

Type URI: https://baseline.openssf.org/attestation/0.1

Version: 0.1

## Purpose

Provide a standard interchange format for measurement of controls on a software
project. In particular, supporting project reporting of [OSPS Baseline] status
through a combination of manual and automated processes.

## Use Cases

OSPS Baseline aims to provide a common set of abstract best-practice controls
which may be implemented by software projects in a number of ways. Some control
implementations can easily be measured by software (possibly using elevated
privileges, such as to read from the GitHub API surface), while other
implementations require human assessment (code collaborators are reviewed prior
to granting escalated permissions). Consumers of this information need a
stardardized format as well as the ability to distinguish the author of
particular measurements and the evidence relied upon.

The [SCAI predicate] includes some elements in common with the baseline
predicate in that each indicates a list of attributes about a set of in-toto
subjects, but Baseline includes a number of features which differ from SCAI:

- `control` values are defined within a `framework` of expected values, rather
  than domain-specific `attributes`.
- controls support explicit `result`s of the evaluation, rather than optional
  `conditions` under which the attribute arrises.
- Baseline `evidence` is a list of mostly human-readable assessment criteria
  (with possible extension via the monotonic principle) rather than a singular
  ResourceDescriptor object.

## Prerequisites

The in-toto Attestation Framework is a necessary pre-requisite, but Baseline may
be implemented as an automated way to capture human answers to control questions
as well as automated measurements of controls.

## Model

Baseline metadata is intended to capture the status of an open source software
project relative to certain expected maturity levels. As open source software
projects take many forms with many different types of assets, Baseline aims to
flexibly define the sets of resources which may encompass a particular project,
including:

- Version controlled source code repositories (generally recorded _without_
  reference to a specific commit)
- Build and release pipelines or environments
- Specific released software assets, such as binaries, libraries, or container
  images
- Accounts in distribution platforms such as NPM, DockerHub, or distributed
  through websites
- Documentation and governance materials for the project
- Subprojects which contain one or more of the above and are managed through the
  same group of collaborators.

As Baseline aims to establish standard levels of project controls which can be
applied by projects regardless of implementation, we expect that the most common
mechanism for reporting controls compliance will be through a combination of
automated and manual reporting. It is expected that the [level 1 criteria] will
generally be amenable to automation on common hosting platforms, but that some
criteria at the higher assessment levels may require manual attestation.

In order to support clear attribution of controls measurement to actors, we
expect that human and automated tooling will generate separate (possibly
overlapping) baseline attestations, each mapped to a specific `author` and set
of in-toto subjects. Consumers will consume an [in-toto attestation bundle]
which contains all the attestations related to the project, and will be
responsible for unifying the control claims across the attestations.

## Schema

The core metadata in the baseline predicate is the project assessment. A project
assessment includes both a set of control assessments (a list with unique
records for each control) and a reference to the entity which performed the
assessment.

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ "uri": "https://github.com/in-toto/attestation.git" }],

  // Predicate:
  "predicateType": "https://baseline.openssf.org/attestation/0.1",
  "predicate": {
    "author": { "uri": "https://github.com/evankanderson" },
    "framework": "https://baseline.openssf.org/versions/2025-10-10",
    "assessedAt": "2025-11-20T23:15:00Z",
    "controls": [
      {
        "control": "OSPS-BR-03.01",
        "result": "passed",
        "evidence": [
          {
            "name": "https://github.com/in-toto/attestation",
            "result": "passed",
            "description": "HTTPS URL"
          },
          {
            "name": "https://in-toto.io/",
            "result": "passed",
            "description": "HTTPS URL"
          }
        ]
      },
      {
        "control": "OSPS-BR-03.02",
        "result": "passed"
        // evidence is optional, not required.
      }
    ]
  }
}
```

This predicate has been somewhat adopted from the [Gemara layer 4] format.

### Parsing Rules

The baseline predicate follows the [in-toto parsing rules]. In particular:

- The baseline admits the possibility of extension fields as a means of
  developing and evolving the base specification. Specific software
  implementations MAY emit additional control assessment and evidence data as
  long as that data follows the monotonic principle: that is, if the field were
  omitted, it would not transform a DENY decision into an ALLOW.
- Consumers MUST accept and ignore unrecognized fields in the predicate.
- While the version of the baseline predicate is "0.1" to indicate that feedback
  and improvement is desired, it is expected that efforts will be made to ensure
  backwards compatibility between "0.1" and future versions.

The following fields are considered mandatory:

- `author` and `framework` in ProjectAssessment.
- `control` and `result` in ControlAssessment.

Additionally, the `result` field of a control assessment in a baseline predicate
is bounded by (no stronger than) the `result`s of the corresponding assessment
results (`evidence`):

- If there is at least one "failure" in the evidence results, the control's
  result is "failure".
- If there are no "failure" results but at least one "needs review" evidence
  result, the control's result may only be "needs review" or "failure".
- A "passed" result for the control requires that all evidence either has a
  "passed" result or does not include an explicit result.

Predicates which do not properly aggregate the evidence results into control
results MAY be discarded by consumers.

### Fields

`predicateType` _string ([TypeURI]), required_

> Identifier for the schema of the baseline predicate. Always
> `https://baseline.openssf.org/attestation/0.1` for this version of the spec.

`predicate.author` _object ([ResourceDescriptor]), required_

> Identifier for the entity (either a human or software program) which assessed
> the control implementations for the project.

`predicate.framework` _string, required_

> Identifier (such as a URL) for the control framework used to assess the
> project.

`predicate.assessedAt` _timestamp, optional_

> Timestamp at which the author completed the assessment.

`predicate.controls` _array of objects (ControlAssessment), required_

> An array of ControlAssessment objects, each of which represents an evaluation
> of the implementation of a particular control indicated by the `control`
> field. No two ControlAssessment objects in the same baseline predicate should
> have the same `control` value.

`predicate.controls[*].control` _string, required_

> A short string representing a specific control from the specified framework.
> Control names should be capitalized, but compared case-insensitively (e.g. for
> the uniqueness requirement).

`predicate.controls[*].result` _enum, required_

> A summarry assessment of whether project meets the specified control. Three
> values are accepted:
>
> - `"passed"`: the control could be determined to be present and enforced
> - `"needs review"`: indeterminate; the control status was unable to be
>   measured but might be present in some other way
> - `"failed"`: the control was not present

`predicate.controls[*].evidence` _array of objects (AssessmentResult), optional_

> An array of AssessmentResult objects, reporting the evidence which the author
> examined to make the `result` determination. Authors are not required to
> provide evidence, but doing so can provide stronger assurance to consumers. If
> evidence is provided, each item of evidence must have a unique `name`.

`predicate.controls[*].evidence[*].name` _string, required_

> A unique identifier for the evidence considered. This MAY be a URL or other
> identifying resource, but may also be a simple string.

`predicate.controls[*].evidence[*].result` _enum, optional_

> A detailed assessment of the control against the evidence considered. Values
> match those of `predicate.controls[*].result`, but the presence of the
> evidence's `result` field is _optional_.

`predicate.controls[*].evidence[*].description` _string, optional_

> A description of the assessment performed.

`predicate.controls[*].evidence[*].message` _string, optional_

> A human-readable explanation of why the evidence did not pass. A message MUST
> only be set if `result` is not `"passed"`.

## Example

### Manual attestation

```jsonc
{
  "type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "Minder",
      "uri": "https://github.com/mindersec/minder.git"
    },
    {
      "name": "Minder Rules and Profiles",
      "uri": "https://github.com/mindersec/minder-rules-and-profiles.git"
    },
    {
      "name": "Minder Website",
      "uri": "https://mindersec.github.io/"
    }
  ],
  "predicateType": "https://baseline.openssf.org/attestation/0.1",
  "predicate": {
    "author": {
      "name": "Evan Anderson",
      "uri": "https://github.com/evankanderson"
    },
    "framework": "https://baseline.openssf.org/versions/2025-02-25",
    "assessedAt": "2025-08-12T18:52:00Z",
    "controls": [
      {
        "control": "OSPS-AC-01.01",
        "result": "passed",
        "evidence": [
          {
            "name": "manual",
            "description": "Checked org MFA controls"
          }
        ]
      },
      {
        "control": "OSPS-AC-02.01",
        "result": "passed"
      },
      {
        "control": "OSPS-AC-03.01",
        "result": "passed",
        "evidence": [
          {
            "name": "manual",
            "description": "Checked minder branch protection"
          }
        ]
      }
    ]
  }
}
```

### Automated attestation

The following is an example of an attestation with additional extension evidence
fields specific to the generating application:

```jsonc
{
  "type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "hello-linux-amd64",
      "uri": "_test/hello-linux-amd64",
      "digest": {
        "sha256": "ec5cbb4dfea31ebb0a69499dbdc77dc6655a9faac2c24183b0e4dfe378f98684",
        "sha512": "353658d0025b1b7ac02eda74f140687ec089a1f394178e9974b15d031b3cdb780fd324082cf098fb29c4c74c5c65df44a8e50cc4cd9eefc11b774d2e7ac6d988"
      }
    }
  ],
  "predicateType": "https://baseline.openssf.org/attestation/0.1",
  "predicate": {
    "author": {
      "name": "Ampel",
      "uri": "https://github.com/carabiner-dev/ampel",
      "digest": {
        "sha256": "09534776aedb9ff9adbc35931bfdc91bc45691737b8486fd8423ba2d020338fe"
      }
    },
    "framework": "https://baseline.openssf.org/versions/2025-02-25",
    "assessedAt": "2025-07-28T21:47:27.767Z",
    "controls": [
      {
        "control": "OSPS-AC-01.01",
        "result": "failed",
        "evidence": [
          {
            "name": "mfa-carabiner-dev",
            "result": "failed",
            "message": "Multifactor authentication is not enabled for some members",
            "subject": {
              "name": "github.com/carabiner-dev",
              "uri": "https://github.com/carabiner-dev",
              "digest": {
                "sha256": "2775bba8b2170bef2f91b79d4f179fd87724ffee32b4a20b8304856fd3bf4b8f"
              }
            },
            "statements": [
              {
                "type": "http://github.com/carabiner-dev/snappy/specs/mfa.yaml",
                "attestation": {
                  "name": "jsonl:_test/attestations.jsonl#2",
                  "uri": "jsonl:_test/attestations.jsonl#2",
                  "digest": {
                    "sha256": "ef7557732dc75f9d9ea85c16785ff0ca736bf8432fedef78d59740c74bc0edae",
                    "sha512": "ee0c5f43e831f137cf0d3ee40482809aeb035fe357d3cd4212464653a831f0dfefdc031ead03ef78dfdf980f106f072f9d7b8a32e6e1dc2f26b3ee28077c4d63"
                  }
                }
              }
            ],
            "guidance": "Enable MFA in the GitHub organization settings to force all members to turn it on"
          }
        ]
      }
    ]
  }
}
```

## Changelog and Migrations

### 0.1

This is the initial version of the predicate.

[OSPS Baseline]: https://baseline.openssf.org/
[SCAI predicate]: ./scai.md
[level 1 criteria]: https://baseline.openssf.org/versions/2025-10-10#level-1
[in-toto attestation bundle]: ../v1/bundle.md
[Gemara layer 4]:
  https://github.com/ossf/gemara?tab=readme-ov-file#layer-4-evaluation
[in-toto parsing rules]: ../v1/README.md#parsing-rules
[ResourceDescriptior]: ../v1/resource_descriptor.md
[TypeURI]: ../v1/field_types.md#typeuri
