# Predicate type: Human Reviews

Type URI: https://in-toto.io/attestation/human-review/v0.1

Version: 0.1.0

Authors:

## Purpose

This is a generic predicate type used to describe human reviews of software
artifacts. For example, this predicate (or a derivative) can be used to wrap
results of code reviews, software audits, and dependency reviews such as the
crev project.

## Use Cases

Software supply chains encompass many types of human reviews. Best practices
including standards like SLSA recommend two or multi-party review for source
code. Another example is the crev project where open source software
dependencies are socially reviewed.

## Model

This predicate type includes one compulsory field, `result`, that describes the
result of the review. Derivatives of this predicate may include other contextual
fields pertaining to different code review systems or use cases. The predicate
includes an optional `reviewLink` field that can be used to find the full
review. In the general case, the subject of a code review attestation must
identify the software artifacts the review applies to.

## Schema

```json
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{...}],
    "predicateType": "https://in-toto.io/attestation/human-review/v0.1",
    "predicate": {
        "result": "pass|fail",
        "reviewLink": "<URL>",
        "timestamp": "<TIMESTAMP>"
    }
}
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1/README.md#parsing-rules).

### Fields

`result` _enum_, _optional_

Indicates the result of the review. For example, it may have the values "pass"
or "fail".

`reviewLink` _URI_, _optional_

Contains a link to the full review. Useful to point to information that cannot
be captured in the attestation.

`timestamp` _Timestamp_, _optional_

Indicates time of review creation.

## Example

```json
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "digest": {
            "gitCommit": "b0bd8ecab5607e174fa403e002c74a666e7edd51"
        }
    }],
    "predicateType": "https://in-toto.io/attestation/human-review/v0.1",
    "predicate": {
        "result": "pass",
        "reviewLink": "https://github.com/in-toto/in-toto/pull/503#pullrequestreview-1341209941",
        "timestamp": "2023-03-15T11:05:00Z"
    }
}
```
