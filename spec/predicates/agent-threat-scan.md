# Predicate type: Agent Threat Scan

Type URI: https://in-toto.io/attestation/agent-threat-scan/v0.1

Version: 0.1

Authors: Adam Lin (@eeee2345)

## Purpose

The Agent Threat Scan predicate communicates that an AI agent artifact (such as
a Model Context Protocol server, an agent skill package, or another agent
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
One example of such a ruleset is Agent Threat Rules (ATR), an MIT-licensed
detection ruleset for AI agent threats. The predicate is not tied to ATR; any ruleset that issues stable identifiers and can be identified
by a versioned URI and, preferably, a digest may be referenced through the
`ruleset` [ResourceDescriptor].

## Model

An Agent Threat Scan attestation declares that a named scanner, applying a
named ruleset identified by a versioned URI and, where available, a digest,
evaluated the subject artifact at scannedAt and produced an outcome of pass,
warn, fail, or incomplete, recording which rules it evaluated and, optionally,
enumerating the rules that matched. Multiple Agent Threat Scan attestations
can be produced for the same subject across rulesets
or across time. The predicate does not include the rule definitions
themselves; it references them by stable identifier, and identifies the
ruleset that issued them with a [ResourceDescriptor] so that consumers can
resolve them.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/agent-threat-scan/v0.1",
  "predicate": {
    "scanner": "<ResourceDescriptor>",
    "ruleset": "<ResourceDescriptor>",
    "scannedAt": "<TIMESTAMP>",
    "outcome": "pass|warn|fail|incomplete",
    "evaluatedRules": ["<STRING>"],
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

**`scanner`, required** [ResourceDescriptor]

> Identifies the tool that performed the scan. `uri` SHOULD identify the
> scanner and its version, for example a package URL such as
> `pkg:npm/example-agent-scanner@1.0.0`. `digest` MAY carry the digest of the
> scanner artifact that ran.

**`ruleset`, required** [ResourceDescriptor]

> Identifies the detection ruleset applied during the scan. `uri` MUST be set
> and SHOULD identify the exact version that was applied, for example a
> package URL with a version or a tag URL, so that a consumer can re-resolve
> any ruleId in matches against the ruleset that was used. The ruleset SHOULD
> resolve to a definition of the rules, or to documentation describing how to
> resolve a rule identifier to its rule definition. `digest` SHOULD carry the
> digest of the ruleset artifact that was evaluated, such as a package
> tarball, so that the ruleset is pinned by content and not only by a version
> string.

**`scannedAt`, required** string (Timestamp)

> RFC 3339 timestamp indicating when the scan completed.

**`outcome`, required** string

> One of pass, warn, fail, or incomplete. fail means at least one rule at or
> above the producer's chosen blocking threshold matched. incomplete means no
> rule at or above that threshold matched, but the scanner did not evaluate
> every rule in the ruleset that `ruleset` identifies, for example because a rule was
> unsupported, failed to run, or was excluded by configuration; a scan that
> did not evaluate every rule MUST NOT report pass or warn. warn means every
> rule was evaluated and only rules below the blocking threshold matched.
> pass means every rule was evaluated and no rule matched. Producers SHOULD
> document their threshold mapping where the predicate is consumed.

**`evaluatedRules`, optional** array of string

> The ruleIds of the rules the scanner evaluated against the subject.
> REQUIRED, and non-empty, when outcome is pass, warn, or incomplete; SHOULD
> be present when outcome is fail. When outcome is pass or warn it lists every
> rule in the identified ruleset, which lets a consumer check the coverage
> claim against the ruleset itself; when outcome is incomplete it lists the
> rules that did run. This follows `passedTests` in
> [Test Result](test-result.md), which records which tests ran.

**`matches`, optional** array of object

> The list of rules that matched on the subject. MAY be empty when outcome
> is pass. SHOULD be present and non-empty when outcome is warn or fail, and
> when outcome is incomplete and any rule matched.

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
      "uri": "pkg:npm/example-agent-scanner@1.0.0"
    },
    "ruleset": {
      "uri": "pkg:npm/agent-threat-rules@4.1.5",
      "digest": {"sha256": "a872606780acdd7b15482182b152c28f80ea5bc6d7055c8d16f560be6c08d4c4"}
    },
    "scannedAt": "2026-10-20T10:00:00Z",
    "outcome": "fail",
    "matches": [
      {
        "ruleId": "ATR-2026-00579",
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

[ResourceDescriptor]: ../v1/resource_descriptor.md
