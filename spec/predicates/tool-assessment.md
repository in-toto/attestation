# Predicate type: Tool Assessment

Type URI: https://in-toto.io/attestation/tool-assessment/v0.1

Version: v0.1.0

## Purpose

**The tool assessment attestation references the results and metadata associated with tools used to assess software before or after its creation.** It's primary purpose is to provide an immutable attestation of tool assessment of a software so that it can be bundled with its provenance. This enables to mapping of build provenances to assessments on a per build basis and can reflect entire DevSecOps pipeline processes.

 However, this attestation can be used to describe the assessment of any target using any tool for any purpose and is not restricted to pipelines.

There are many existing predicates that describe the use of specific tool types. While these predicates are well defined, they are narrowly scoped. There should exist a predicate that is general enough to effectively attest the use of any tool that can be used if a tool type does not have a predicate type yet. The tool assessment attestation type aims to solve this.

-   The [cyclonedx](cyclonedx.md) and [spdx](spdx.md) predicate types describe SBOM standards.
-   The [test result](test-result.md) predicate type describes test running tests in the software supply chain.
-   The [vulnerabilities](vulns_02.md) predicate type describes the results of a vulnerability scan. This predicate closely resembles the type of information desired to be captured by the tool assessment attestation but is too narrowly scoped to producers of vulnerability information.
-   The [SCAI](scai.md) predicate type captures functional attribute and integrity information about software and its supply chain. It is the closest predicate for this use-case but fails to cleanly map a result to its tooling while providing appropriate metadata on the tooling or the policy requiring its execution. The tool assessment attestation would serve well as an attribute predicate in the SCAI framework.

Prior existing predicates still have their own important use-cases. This predicate type does not aim to replace them but to provide a specification flexible enough to use for any type of tool.

## Use Cases

### Control Gates

Control gates are an essential and increasingly prevalent requirement in many
DevSecOps pipelines where continuous integration is adopted as a standard of
assessment for a piece of software to be deployed and operated.
Tool assessment attestations enable security control assessors to audit compliance to policy in an immutable format.

### Policy as Code Enabling via Attachment to Build Artifacts

Tool assessment attestations bundled with container build provenance can enable policy-as-code enforcement of containers or software on IT systems.

## Prerequisites

-   [in-toto Attestation Framework](https://github.com/in-toto/attestation/blob/main/spec/README.md)
-   Appropriate knowledge in capturing tool metadata and processing results

## Model

This predicate type is based on three parts: describing the tool, its configuration, and the results of the tool. Due to the generalness of the predicate, some fields will be optional. The `summary` field is included as a non-prescribed object for producers to include more specific data to be attested to from the ran tool.

This also defines the `Profile` object type.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/tool-assessment/v1",
  "predicate": { 
      "tool": {
          "name": "<NAME>",
          "type": "<TOOL_DESCRIPTION>",
          "uri": "<URI>",
          "version": "<VERSION>"
      },
      "config": {
          "profiles": [ Profile ],
          "exclusions": ["<EXCLUSIONS>"],
          "files": [ ResourceDescriptor ], // .ignores, config files, etc
          "full_command": "<FULL_CMD_LINE>"
      },
      "result": "<RESULT>",
      "output": [ ResourceDescriptor ],
      "summary": { /* object */ },
   }
}

// Profile Object Type
Profile: {
    "profile": "<PROFILE_NAME_OR_DECRIPTION>",
    "uri": "<PROFILE_URI>",
    "version": "<VERISON>",
    "last_updated": "<TIME_STAMP>",
    "annotations": { /* object */ }
}

```

### Parsing Rules

-   Consumers MUST ignore unrecognized fields unless otherwise noted.
-   Acceptable formats of the `summary` and `annotations` fields are up to the producer and consumer.

### Fields

**tool, required** object

> Object associated with identifying the tool.

**tool.name, required** string

>> Name of tool.

**tool.type, required** string

>> Description of the type of tool (SAST, DAST, SECRETS, etc).

**tool.uri, required** string (ResourceURI)

>> URI indicating the identity of the source of the tool.

**tool.version, optional** string

>> Version of the tool.

**config, required** object

> Object that describes the configuration of the tool.
> This object should be descriptive enough to reproduce the results of the tool based on its entries.

**config.profiles, optional** Profile list

>> Contains a list of profiles used in the tool. Profiles describe the set of data that the tool references in its execution that may modify the behavior of the tool or its results. This includes: rulesets, databases, policies, etc.

**config.exclusions, optional** string list

>> List of deviations from the profile, such as rule IDs, file names, ignores, etc.

**config.files, optional** ResourceDescriptor list

>> Reference to files used by the tool to modify configuration.

**config.full_command, required** string

>> Command used to run the tool.

**result, required** string

> Result of the tool execution. Usually `PASS` or `FAIL`.

**output, required** ResourceDescriptor list

> Artifacts associated with the result of the execution of the tool.

**summary, optional** object

> Object containing extra fields associated with the execution of the tool that contribute to the understanding of a tools results. Acceptable formats are up to the producer and consumer of the attestation.

---

**Profile.profile, required** string

>> Name or description of the profile.

**Profile.uri, required** string

>> URI identifying the source of the profile

**Profile.version, optional** string

>> Version of the profile

**Profile.last_updated, optional** string

>> Timestamp of the last update of the profile

**Profile.annoations, optional** object

>> Extraneous data associated with the tool assessment.

## Example

### Semgrep

```jsonc
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "name": "foo",
        "digest": { "sha256": "78ec328..." }
    }],
    "predicateType": "https://in-toto.io/attestation/tool-assessment/v0.1",
    "predicate": {
          "tool": {
              "name": "Semgrep",
              "type": "SAST",
              "uri": "pkg:github/semgrep/semgrep@984f760",
              "version": "1.139.0"
          },
          "config": {
              "profiles": [
                {
                  "profile": "Default Python",
                  "uri": "https://semgrep.dev/p/python"
                },
                {
                  "profile": "Community Python",
                  "uri": "https://github.com/semgrep/semgrep-rules/tree/d375208f04370b4e8d3ca7fe668db6f0465bb643/python",
                  "last_updated": "2025-06-04T19:25:00Z"
                }],
              "exclusions": ["bar.py"],
              "full_command": "semgrep scan --config p/python --config rules/python --exclude='bar.py'"
          },
          "result": "PASS",
          "output": ["<ResourceDescriptor(semgrep_output.txt)>"]
    }
}
```

### Trufflehog

```jsonc
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "name": "foo",
        "digest": { "sha256": "78ec328..." }
    }],
    "predicateType": "https://in-toto.io/attestation/tool-assessment/v0.1",
    "predicate": {
          "tool": {
              "name": "Trufflehog",
              "type": "Secrets Scanning",
              "uri": "pkg:github/trufflesecurity/trufflehog@466da4b",
              "version": "3.90.8"
          },
          "config": {
              "profiles": [
                {
                  "profile": "Custom",
                  "uri": "https://example.com/trufflehog_config.yml",
                  "last_updated": "2025-06-04T19:25:00Z"
                }],
              "exclusions": ["excluded_files.txt"],
              "files": [
                "<ResourceDescriptor(trufflehog_config.yml)>",
                "<ResourceDescriptor(excluded_files.txt)>"
              ],
              "full_command": "trufflehog --config=trugglehog_config.yml --no-update git file://. --exclude-paths='excluded_files.txt' --json > th.json"
          },
          "result": "PASS",
          "output": ["<ResourceDescriptor(th.json)>"]
    }
}
```

### OpenSCAP

```jsonc
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "name": "foo",
        "digest": { "sha256": "78ec328..." }
    }],
    "predicateType": "https://in-toto.io/attestation/tool-assessment/v0.1",
    "predicate": {
          "tool": {
              "name": "Openscap",
              "type": "STIG Compliance Scan",
              "uri": "pkg:github/OpenSCAP/openscap@e9b2a41",
              "version": "1.4.2"
          },
          "config": {
              "profiles": [
                {
                  "profile": "Ubuntu",
                  "uri": "https://example.com/1.3/xccdf_ubuntu_profile.xml",
                  "last_updated": "2025-06-04T19:25:00Z",
                  "version": "v1.3"
                }],
              "files": ["<ResourceDescriptor(xccdf_ubuntu_profile.xml)>"],
              "full_command": "oscap xccdf eval --profile Ubuntu --results xccdf-results.xml xccdf_ubuntu_profile.xml"
          },
          "result": "PASS",
          "output": ["<ResourceDescriptor(xccdf-results.xml)>"],
          "summary": {
            "score": 98,
            "total": 214,
            "pass": 21,
            "fail": 2,
            "not_checked": 3,
            "not_applicable": 188
          }
    }
}
```

## Changelog and Migrations
