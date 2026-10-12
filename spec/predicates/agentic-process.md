# Predicate type: Test Result

Type URI: https://in-toto.io/attestation/agentic-process/v0.1

Version: 0.1

Authors: Yuval Kainan (@KainanYuval), Carmit Hershman (@carmithersh) 

## Purpose

This predicate type defines a generic provenance of an agentic proccss. The schema may be extended to support different agentic processes within a software supply chain or other realms, while requiring a minimal schema required for any type of agentic process. 

## Use Cases

Software development processes that were ran by AI agents, either autonomously or with a human in the loop. This attestation
can be used to verify generic agnetic process questions, such as:

1. that the agentic process ran on approved model, harnes and agent 
2. that the agentic process used only approved agnetic tools
3. that an appropreate user is accountable for the process
4. there were humans in the loop (optional)
5. the process session logs are available  

Also, specific process types validations are supported:

6. is the process aligned?
7. is this code-development agnetic process referenced the requirements resources
8. were the organization review-guidelines resources used within the pr-approval agentic process 

Each attestation corresponds to a single agnetic process and may reference multiple sessions logs that were run as part of that process.

## Prerequisites

1. Understanding of the [in-toto attestation specification](https://github.com/in-toto/attestation).
2. A runtime tool monitoring agnetic processes on-device or through a proxy that has visibility to the agent sessions and the process events 

## Model

This predicate type includes few compulsory fields which should allow basic provenance validations of any agentic process (e.g. code development, PR approval, version promotion, code-scan, support-case), while also allowing specific process types validations through custom fields.

This predicate model is inspired by https://slsa.dev/provenance

## Schema
The schema of this predicate type is documented in [Agentic Process Evidence
](https://github.com/jfrog/agentic-process-evidence/blob/main/spec/agentic-process-evidence.md)

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [...],
  "predicateType": "https://in-toto.io/attestation/agentic-process/v0.1.0",
  "predicate": {
    "processType": "<STRING>",
    "providers": [
      {
        "harness": {
          "name": "<STRING>",
          "version": "<STRING>"
        },
        "agent": {
          "id": "<STRING>",
          "name": "<STRING>",
          "version": "<STRING>"
        },
        "languageModels": [
          {
            "inferenceProvider": "<STRING>",
            "resolved": "<STRING>"
          }
        ]
      }
    ],
    "traceId": "<STRING>",
    "sessionsLogs": [<ResourceDescriptor>, ...],
    "tools": [{"name": "<STRING>", "version": "<STRING>"}],
    "contextArtifacts": [<ResourceDescriptor>, ...],
    "custom": { /* object */ },
    "result": "<STRING>",
    "intents": ["<STRING>"],
    "processSummary": "<STRING>",
    "owner": "<STRING>",
    "reviewers": ["<STRING>"],
    "startTimestamp": "<TIMESTAMP>",
    "endTimestamp": "<TIMESTAMP>"
  },
  "createdAt": "<TIMESTAMP>",
  "createdBy": "<STRING>"
}
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1/README.md#parsing-rules).

### Fields

#### mandatory fields:
`processType` _string_, _required_
Type of process, e.g. code-development, pr-review, promotion, support-case  

`providers[]` list of _object_, _required_
Models, agents and harnesses running the process  
`providers[].harness`  _object_, _required_
`providers[].harness.name`  _string_, _required_
Harness product name
`providers[].harness.version`  _string_, _required_
Harness version as emitted by the producer (not necessarily SemVer)
Agentic provider harness
`providers[].agent`  _string_, _optional_
Agent running in the harness
`providers[].languageModels`  list of _object_, _optional_
Language models involved in the run
`providers[].languageModels[].inferenceProvider`  _string_, _optional_
LLM requested (name/version), if known
`providers[].languageModels[].resolved`  _string_, _optional_
LLM actually used (name/version), if known

`result` _string_, _required_
Short process completion state

`owner` _string_, _required_
Login or email of the accountable user

`startTimestamp` _ISO 8601 timestamp_, _required_
Process start
`endTimestamp`  _ISO 8601 timestamp_, _required_
Process end 

#### optional fields:
`traceId` _string_, _optional_ 
Unique identifier of the agentic process

`sessionsLogs[]` list of ResourceDescriptor, _optional_
Links to session logs (full agentic chat context)

`tools[]` list of _object_, _optional_
list of agentic tools used throughout the process
`tools[].name` _string__, _required_ 
Tool name as the harness reports it
`tools[].version` _string__, _optional_ 
Tool version as the producer reports it (not necessarily SemVer)

`contextArtifacts[]` list of _object_, _optional_
Artifacts used by the process (policies, instructions, prior logs, guidelines, …)
`contextArtifacts[].tags` list of _string_ , _optional_
values SHOULD be from Context artifact tags; additional tags allowed
Context artifact tags: policy, guideline, instructions, priorLog
`contextArtifacts[].uri`  _string (URI)_, exactly one of uri or data
Location of the artifact
`contextArtifacts[].data`  _string_, exactly one of uri or data
Inline content of the artifact
`contextArtifacts[].digest`  _DigestSet_, _optional_
Immutable identity of the artifact bytes

`custom` _object_, _optional_
Custom information for the specific agentic process

`intents` list of _string_, _optional_
Short descriptions of intents identified in the process

`processSummary` _string_, _optional_
Summarized description of the process

`reviewers[]` list of _string_, _optional_
Logins or emails of human overseers (e.g. developer chatting during development or PR review)


## Example

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [...],
  "predicateType": "https://in-toto.io/attestation/agentic-process/v0.1.0",
  "predicate": {
    "processType": "code-development",
    "providers": [
      {
        "harness": {
          "name": "cursor",
          "version": "3.5.33"
        },
        "agent": {
          "id": "stable-agent-id",
          "name": "cursor-agent",
          "version": "1.0"
        },
        "languageModels": [
          {
            "inferenceProvider": "anthropic/claude-opus-4-8",
            "resolved": "anthropic/claude-opus-4-8-20260701"
          }
        ]
      }
    ],
    "traceId": "process id | ci_job_run_uri",
    "sessionsLogs": [
      {
        "uri": "session log url",
        "digest": {
          "sha256": "session log sha..."
        }
      }
    ],
    "tools": [{"name": "tool-name", "version": "tool version"}],
    "contextArtifacts": [
      {
        "tags": ["policy"],
        "uri": "link-to-development-guideline",
        "data": "The complete artifact document inline",
        "digest": {
          "sha256": "policy-sha256"
        }
      },
      {
        "tags": ["guideline"],
        "uri": "link-to-architecture-guidelines",
        "digest": {
          "sha256": "architecture-sha256"
        }
      }
    ],
    "custom": {
      "baseCommit": {
        "uri": "https://github.com/myorg/hello-world/commit/bf02510c8ce0b804a099797510af...c19325acfb979538c5b521304c83cde63892",
        "digest": {
          "gitCommit": "bf02510c8ce0b804a099797510af...c19325acfb979538c5b521304c83cde63892"
        }
      },
      "requirements": [
        {
          "issue": "https://myorg.atlassian.net/browse/FIN-3678",
          "title": "my requirements subject"
        }
      ]
    },
    "result": "COMPLETED",
    "intents": ["short description"],
    "processSummary": "long text",
    "owner": "login|email",
    "reviewers": ["login|email"],
    "startTimestamp": "ISO 8601 timestamp",
    "endTimestamp": "ISO 8601 timestamp"
  },
  "createdAt": "2026-04-08T08:55:17.242Z",
  "createdBy": "dev-auto-reviewer"
}
```