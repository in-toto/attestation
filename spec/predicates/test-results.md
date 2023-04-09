# Predicate type: Test Results

Type URI: https://in-toto.io/attestation/test-result/v0.1

Version: 0.1.0

Authors: Aditya Sirish A Yelgundhalli (@adityasaky)

## Purpose

This predicate type defines a generic schema to express the result of running
tests in software supply chains. The schema may be extended to support different
types of testing or specific test harnesses.

## Use Cases

Software development processes include several types of tests. This attestation
can be used to express the results of running those tests. It can be used to
verify:

1.  that all tests were in fact run, and
2.  that all required tests passed

Therefore, each attestation corresponds to one invocation of a test suite, and
may include the results of several individual tests.

### Asserting Test Configurations Used

The supply chain owner creates a policy that records the expected test
configurations. During verification, the policy checks that the test attestation
used the right configurations. If verified using in-toto layouts, a custom
inspection may optionally parse the `url` field to verify attestation matches
the test run.

### Asserting Test Results

In addition to the previous use case, the supply chain owner creates a policy
verifying that test results passed. In the simplest case, the policy applies to
all tests. Therefore, it asserts the contents of `result` and that
`failedTests` is empty. In more nuanced cases, a subset of tests may
matter. For example, the tested artifact may be an OS image that's to be
deployed to three types of devices, A, B, and C. As such, the test harness
validates the new image on an instance of each device. When verifying the
attestation prior to an image being installed on a device of type A, it only
matters that the tests passed on the corresponding test device of type A and not
necessarily the others. The test result attestation can be used successfully
even if the tests failed on C. As such, the verification policy ensures that
tests for A are all listed in `passedTests` or possibly `warnedTests`.

As before, if the attestation is verified using an in-toto layout, a custom
inspection may examine the `url` contents to verify the contents of the
attestation.

## Prerequisites

Understanding of the
[in-toto attestation specification](https://github.com/in-toto/attestation).

## Model

This predicate type includes two compulsory fields, `result` that describes the
result of the test run, and `configuration` that contains the configuration used
for the test invocation. The optional `url` field contains a link to the test
run. `passedTests`, `warnedTests`, and `failedTests` are lists that record the
names of tests that passed with no errors or warnings, passed with warnings, and
failed with errors respectively. The expected `subject` are the source artifacts
tested.

## Schema

```json
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{...}],
    "predicateType": "https://in-toto.io/attestation/test-result/v0.1",
    "predicate": {
        "result": "pass|fail",
        "configuration": ["<ResourceDescriptor>", ...],
        "url": "<URL>",
        "passedTests": ["<TEST_NAME>", ...],
        "warnedTests": ["<TEST_NAME>", ...],
        "failedTests": ["<TEST_NAME>", ...]
    }
}
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1.0/README.md#parsing-rules).

### Fields

`result` _boolean_ , _required_

Indicates the result of the test run. If true, it indicates _all_ tests passed
in the corresponding run. This means that `warnedTests` and `failedTests` must
both be empty when they are used.

`configuration` _list of ResourceDescriptor_, _required_

Reference to the configuration used for the test run.

`url` _ResourceURI_, _optional_

Contains a URL to the test run, if applicable. This may be used to find
information such as the logs for test execution, to confirm the test was
performed against the expected subject, and more. The predicate makes no
assumptions about the test harness or system used which this field would point
to. Instead, verifiers are expected to ascertain how to use this field's
response separately, perhaps via the URL itself. For example, this may point to
a GitHub Actions run which may be determined using the domain of the URL. On the
other hand, for custom Jenkins deployments and the like, the verifier should
determine out of band how to use this field.

`passedTests` _list of strings_, _optional_

Each entry corresponds to the name of a single test that passed. The semantics
of the name must be determined separately between the producer and consumer.

`warnedTests` _list of strings_, _optional_

Each entry corresponds to the name of a single test that expressed a warning.
The semantics of the name must be determined separately between the producer and
consumer.

`failedTests` _list of strings_, _optional_

Each entry corresponds to the name of a single test that failed. The semantics
of the name must be determined separately between the producer and consumer.

## Example

```json
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [
        {
            "digest": {
                "gitCommit": "d20ace7968ba43c0219f62d71334c1095bab1602"
            }
        }
    ],
    "predicateType": "https://in-toto.io/attestation/test-result/v0.1",
    "predicate": {
        "result": "pass",
        "configuration": [{
            "name": ".github/workflows/ci.yml",
            "downloadLocation": "https://github.com/in-toto/in-toto/blob/d20ace7968ba43c0219f62d71334c1095bab1602/.github/workflows/ci.yml",
            "digest": {
                "gitBlob": "ebe4add40f63c3c98bc9b32ff1e736f04120b023"
            }
        }],
        "url": "https://github.com/in-toto/in-toto/actions/runs/4425592351",
        "passedTests": [
            "build (3.7, ubuntu-latest, py)",
            "build (3.7, macos-latest, py)",
            "build (3.7, windows-latest, py)",
            "build (3.8, ubuntu-latest, py)",
            "build (3.8, macos-latest, py)",
            "build (3.8, windows-latest, py)",
            "build (3.9, ubuntu-latest, py)",
            "build (3.9, macos-latest, py)",
            "build (3.9, windows-latest, py)",
            "build (3.10, ubuntu-latest, py)",
            "build (3.10, macos-latest, py)",
            "build (3.10, windows-latest, py)",
            "build (3.x, ubuntu-latest, lint)"
        ],
        "warnedTests": [],
        "failedTests": []
    }
}
```
