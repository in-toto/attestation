# Predicate type: Agent Threat Scan

Type URI: https://in-toto.io/attestation/agent-threat-scan/v0.1

Version: 0.1

Authors: Adam Lin (@eeee2345)

## Purpose

The Agent Threat Scan predicate communicates that an AI agent artifact (such as
a Model Context Protocol server, a Claude Code skill, or another agent
configuration manifest) has been scanned against a named detection ruleset for
agent-specific threats including prompt injection, tool poisoning, MCP request
forgery, and skill compromise. It enables consumers of agent artifacts to
verify which ruleset and ruleset version were applied, what the scan outcome
was, and which rules (if any) matched, without needing to re-run the scan or
read the full set of rule definitions.

## Use Cases

Existing predicates such as Vulnerabilities (vulns), Test Result, and Simple
Verification Result do not capture the structure of agent-specific scanning.
Vulnerability scanners are oriented around package CVEs in conventional
software artifacts. Test Result is generic and does not carry the
ruleset-versioning, rule-identifier, and threat-class taxonomy that are
characteristic of agent threat detection. Simple Verification Result records
that policies passed but does not enumerate which rules matched on a fail.

Concrete use cases the Agent Threat Scan predicate enables include the
following. A registry of MCP servers can require an attestation that an
artifact was scanned with a named open ruleset at or above a minimum version
before the artifact is listed. A CI pipeline producing agent skills can attach
the predicate to released artifacts so downstream consumers can decide whether
to install. A policy engine can answer the question, was this agent scanned
with ruleset X version Y on date Z, and did any rules of severity high or
critical match.

## Prerequisites

This predicate depends on the in-toto Attestation Framework and on the
existence of a named detection ruleset that issues stable rule identifiers.
One example of such a ruleset is Agent Threat Rules (ATR), an open detection
standard for AI agent threats licensed Apache-2.0, available at
https://github.com/Agent-Threat-Rule/agent-threat-rules. The predicate is not
tied to ATR; any ruleset that issues stable identifiers and a version may be
referenced through the ruleset.uri and ruleset.version fields. ATR is shipped
in production at Cisco AI Defense and Microsoft agent-governance-toolkit.

## Model

An Agent Threat Scan attestation declares that a named scanner, applying a
named ruleset at a specific version, evaluated the subject artifact at
scannedAt and produced an outcome of pass, warn, or fail, optionally
enumerating the rules that matched. Multiple Agent Threat Scan attestations
can be produced for the same subject across rulesets or across time. The
predicate does not include the rule definitions themselves; it references
them by stable identifier and by ruleset version so that consumers can resolve
them through the ruleset.uri.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/agent-threat-scan/v0.1",
  "predicate": {
    "scanner": {
      "uri": "<URI>",
      "version": "<STRING>"
    },
    "ruleset": {
      "uri": "<URI>",
      "version": "<STRING>"
    },
    "scannedAt": "<TIMESTAMP>",
    "outcome": "pass|warn|fail",
    "matches": [
      {
        "ruleId": "<STRING>",
        "severity": "low|medium|high|critical",
        "threatClass": "<STRING>",
        "evidence": "<STRING>"
      }
    ]
  }
}
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1/README.md#parsing-rules).

### Fields

**`scanner`, required** object

> Identifies the tool that performed the scan.

**`scanner.uri`, required** string (ResourceURI)

> URI identifying the scanner. May reference a package, repository, or service.

**`scanner.version`, optional** string

> The version of the scanner.

**`ruleset`, required** object

> Identifies the detection ruleset applied during the scan.

**`ruleset.uri`, required** string (ResourceURI)

> URI identifying the ruleset. SHOULD resolve to a definition of the rules,
> or to documentation describing how to resolve a rule identifier to its
> rule definition.

**`ruleset.version`, required** string

> The version of the ruleset that was applied. The version SHOULD be stable
> and citable so that a consumer can re-resolve any ruleId in matches against
> the exact ruleset version used.

**`scannedAt`, required** string (Timestamp)

> RFC 3339 timestamp indicating when the scan completed.

**`outcome`, required** string

> One of pass, warn, or fail. pass means no rules matched. warn means at most
> rules below the producer's chosen blocking threshold matched. fail means at
> least one rule at or above the blocking threshold matched. Producers SHOULD
> document their threshold mapping where the predicate is consumed.

**`matches`, optional** array of object

> The list of rules that matched on the subject. MAY be empty when outcome
> is pass. SHOULD be present and non-empty when outcome is warn or fail.

**`matches[*].ruleId`, required** string

> The stable identifier of the rule that matched, as issued by the ruleset.

**`matches[*].severity`, optional** string

> One of low, medium, high, critical. The severity assigned by the ruleset
> to the matched rule.

**`matches[*].threatClass`, optional** string

> A short label describing the class of threat the rule covers, such as
> prompt_injection, tool_poisoning, mcp_request_forgery, or skill_compromise.

**`matches[*].evidence`, optional** string

> A free-form string containing scanner-specific evidence such as a snippet,
> a path within the subject, or a hash. Producers SHOULD avoid including
> sensitive content in this field.

## Example

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "example-mcp-server-1.2.3.tgz",
      "digest": {"sha256": "fe4fe40ac7250263c5dbe1cf3138912f3f416140aa248637a60d65fe22c47da4"}
    }
  ],
  "predicateType": "https://in-toto.io/attestation/agent-threat-scan/v0.1",
  "predicate": {
    "scanner": {
      "uri": "pkg:npm/%40panguard-ai/cli@1.4.13",
      "version": "1.4.13"
    },
    "ruleset": {
      "uri": "https://github.com/Agent-Threat-Rule/agent-threat-rules",
      "version": "v2.0.17"
    },
    "scannedAt": "2026-05-09T10:00:00Z",
    "outcome": "fail",
    "matches": [
      {
        "ruleId": "ATR-2026-00081",
        "severity": "high",
        "threatClass": "prompt_injection",
        "evidence": "tools[2].description contained a hidden instruction directive"
      }
    ]
  }
}
```

## Changelog and Migrations

Not applicable for this initial version.

[Attestation]: ../README.md
