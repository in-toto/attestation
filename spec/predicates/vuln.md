# Predicate type: Vulnerabilities

Type URI: https://in-toto.io/attestation/vulns

Version: 0.1

## Purpose

The definition of a vulnerability attestation type has been discussed in the past in [in-toto attestation](https://github.com/in-toto/attestation/issues/58) and [issue](https://github.com/sigstore/cosign/issues/442). However we need to identify two different purposes from these initial conversations:

-   The definition of a common format to represent the results of a vulnerability report.
-   The definition of a certain set of metadata fields that could help to consume these vulnerability attestations from the different scanning tools.

Obviously the first goal is quite challenging and requires a bigger community to agree upon a specific format. As a consequence, the following attestation type focuses on the definition of that common metadata which could enable the beginning of an exportable and manageable vulnerability attestation.

This document describes a vulnerability attestation type to represent vulnerability reports from the scanners in an "exportable" manner and independently of the format chosen to output the results.

## Prerequisites

The in-toto [attestation] framework and a Vulnerability scanner tools such as [Grype](https://github.com/anchore/grype), [Trivy](https://github.com/aquasecurity/trivy) and others.

## Use cases

When sharing the results of a vulnerability scan using an attestation, there is certain metadata that is crucial to trust and reuse this information.
Information about the scanner used during the scanning is relevant to trust these results. The state of the vulnerability database used to search for vulnerabilities defines the accuracy of the results. Other metadata information such as the timestamp when the scan finished could define the reusability of these results.

## Model

This is a predicate type that fits within the larger [Attestation] framework.
The following model aims to provide a well defined list of fields so that consumers know how to start exchanging their scanner results.

This predicate model is inspired by [cosign vulnerability attestation](https://github.com/sigstore/cosign/blob/main/specs/COSIGN_VULN_ATTESTATION_SPEC.md).

## Schema

The schema of this predicate type is documented below.

### Fields

The fields that make up this predicate type are:

The `subject` contains whatever software artifacts are to be associated with this vulnerability report document.
The `predicate` contains a JSON-encoded data with the following fields:

**scanner, required** object

> There are lots of container image scanners such as Trivy, Grype, Clair, etc.
> This field describes which scanner is used while performing a container image scan,
> as well as version information and which Vulnerability DB is used.

**scanner.uri, required** string (ResourceURI)

> > URI indicating the identity of the source of the scanner.

**scanner.version, optional** string (ResourceURI)

> > The version of the scanner.

**scanner.db.uri, optional** string (ResourceURI)

> > > URI indicating the identity of the source of the Vulnerability DB.

**scanner.db.version, optional** string

> > > The version of the Vulnerability DB.

**scanner.db.lastUpdate, optional** string (timestamp)

> > > The timestamp of when the vulnerability DB was updated last time.

**scanner.result, required** object list

> > The result contains a list of vulnerabilities. Note that an empty list means the **scanner** found no vulnerabilities.
> > This is the most important part of this field because it'll store the scan result as a whole. So, people might want
> > to use this field to take decisions based on them by making use of Policy Engines tooling whether allow or deny these images.
> > Each object defines information about each one of the vulnerabilities found by the scanner.


**scanner.result.[*].id, required** string

> > > This is the identifier of the vulnerability, e.g. [GHSA-fxph-q3j8-mv87](https://github.com/advisories/GHSA-fxph-q3j8-mv87) whose CVE id is [CVE-2017-5645](https://nvd.nist.gov/vuln/detail/CVE-2017-5645).

**scanner.result.[*].severity, required** object list

> > > The severity contains a list to describe the severity of a vulnerability using one or more quantitative scoring method.

**scanner.result.[*].severity.[*].method, required** string

> > > > The method describes the quantitative method used to calculate the associated severity score such as nvd, cvss and others.

**scanner.result.[*].severity.[*].score, required** string

> > > > This is a string representing the severity score based on the selected method.

**scanner.result.[*].annotations, optional** object list

> > > This is a list of key/value pairs where scanners can add additional custom information.

**metadata.scanStartedOn, optional** Timestamp

> > The timestamp of when the scan started.

**metadata.scanFinishedOn, required** Timestamp

> > The timestamp of when the scan completed.

## Example

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "foo.jar",
      "digest": {"sha256": "fe4fe40ac7250263c5dbe1cf3138912f3f416140aa248637a60d65fe22c47da4"}
    }
  ],
  // Predicate:
  "predicateType": "https://in-toto.io/attestation/vulns/v0.1",
  "predicate": {
    "scanner": {
      "uri": "pkg:github/aquasecurity/trivy@244fd47e07d1004f0aed9",
      "version": "0.19.2",
      "db": {
        "uri": "pkg:github/aquasecurity/trivy-db/commit/4c76bb580b2736d67751410fa4ab66d2b6b9b27d",
        "version": "v1-2021080612",
        "lastUpdate": "2021-08-06T17:45:50.52Z"
      },
      "result": [
        {
         "id": "CVE-123",
         "severity": [
          { "method": "nvd", "score": "medium"},
          { "method": "cvss_score", "score", "5.2" }
         ]
        },
        {...}
      ]
    },
    "metadata": {
      "scanStartedOn": "2021-08-06T17:45:50.52Z",
      "scanFinishedOn": "2021-08-06T17:50:50.52Z"
    }
  }
}
```

## Changelog and Migrations

Not applicable for this initial version.

2024-10-11:
- List types updated to disambiguate list types.
- Remove invocation from example (which is not in the spec).
- Removed the additional "vulnerability" object indirection
  from the spec since it is the only field within it (as reflected in
  the example).
- Severity changed to a list to contain multiple scores (as reflected in
  the example).
- ScanStartedOn changed from required to optional as it may not always be
  available.
- scanner.db.lastUpdate changed from required to optional since this information
  may not always be available when a DB uri/version is known.

[Attestation]: ../README.md
