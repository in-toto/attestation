# Predicate type: Dependency Reviews (crev)

Type URI: (tentative) https://in-toto.io/attestation/human-review/crev/v0.1

Version: Do we use Crev's versioning? Currently -1? What about their URI?

## Purpose

This attestation type is used to describe the results of human review of
dependency source code. The format is based on the
[crev project](https://github.com/crev-dev/crev).

## Use Cases

As noted above, crev enables social review of popular open source software
dependencies. A crev review includes information such as the thoroughness of
the review, understanding of the source code, and a final rating.

## Model

Most modern software have external dependencies. Dependency review is the
process of reviewing and verifying the source code of a particular
dependency, and can be performed by one or more of several actors in the supply
chain. The developer importing a new dependency can perform the review or a
dedicated security team can be tasked with it.

## Schema

```json
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{...}],
    "predicateType": "https://in-toto.io/attestation/human-review/crev/v0.1",
    "predicate": {
        "result": "positive|negative",
        "timestamp": "<TIMESTAMP>",
        "reviewer": {
            "idType": "crev",
            "id": "",
            "url": "<URL>"        
        },
        "thoroughness": "high|medium|low",
        "understanding": "high|medium|low",
        "comment": "<STRING>"
    }
}
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1/README.md#parsing-rules).

### Fields

The subject of this predicate type is a specific package and its version in some
ecosystem.

`result`, _enum_, _required_

Specifies if the overall rating of the dependency is `positive` or `negative`.

`reviewer` _object_, _required_

Identifies the reviewer. This has some meaning for crev's trust proliferation
aspects, but the identity of the reviewer can also be mapped based on in-toto's
functionary handling. `idType` is used to determine the contents of `reviewer`.
The `url` is a reference to the reviewer's crev-proofs repository.

`timestamp` _Timestamp_, _required_

Indicates time of review creation.

`thoroughness` _enum_, _required_

Describes how thorough the reviewer was. Must be set to one of `low`, `medium`,
or `high`.

`understanding` _enum_, _required_

Describes the reviewer's understanding of the dependency code. Must be set to
one of `low`, `medium`, or `high`.

`comment` _string_, _optional_

Optional field with any other comments a reviewer may have about the
dependency.

## Example

```json
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [
        {
            "name": "in-toto",
            "uri": "purl+pkg:pypi/in-toto@1.3.2",
            "digest": {
                "sha256": "aa12e63298425cfc4773ed03febd68a384c63b2690959dd788f8c4511ea97bbe"
            },
            "downloadLocation": "https://github.com/in-toto/in-toto/releases/download/v1.3.2/in_toto-1.3.2-py3-none-any.whl"
        },
    ],
    "predicateType": "https://in-toto.io/attestation/human-review/crev/v0.1",
    "predicate": {
        "result": "positive",
        "timestamp": "2023-03-16T00:09:27Z",
        "reviewer": {
            "idType": "github",
            "id": "adityasaky",
            "url": "https://github.com/adityasaky/crev-proofs"
        },
        "thoroughness": "high",
        "understanding": "high",
        "comment": "This dependency is well written and can be used safely."
    }
}
```
