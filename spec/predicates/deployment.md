# Predicate type: Deployment

Type URI: https://in-toto.io/attestation/deployment

Version 1.0

## Purpose

To authoritatively express which environment an artifact is allowed to be deployed to.

## Use Cases

When deploying an artifact (e.g., a container), we want to restrict which environment
the artifact is allowed to be deployed / run. The environment has access
to resources we want to protect, such as a service account, a Spiffe ID, a Kubernetes pod ID, etc.
The deployment attestation authoritatively binds an artifact to a deployment environment
where an artifact is allowed to be deployed.

The ability to bind an artifact to an environment is paramount to reduce the blast radius
if vulnerabilties are exploited or environments are compromised. Attackers who gain access
to an environment will pivot based on the privileges of this environment, so it is imperative to
follow the privilege of least principle and restrict _which_ code is allowed to run in _which_ environment.
For example, we would not want to deploy a container with remote shell capabilities on a pod that processes
user credentials, even if this container is integrity protected at the highest SLSA level.
Conceptually, this is similar to how we think about sandboxing and least-privilege principle on
operating systems. The same concepts apply to different types of environments, including cloud environments.

These use cases are not hypothetical. Binding an artifact to its expected deployment environment is
one of the principles used internally at Google; it is also a feature provided by [Google Cloud Binauthz](https://cloud.google.com/binary-authorization/).

The decision to allow or deny a deployment request may happen in "real-time", i.e. the control plane
may query an online authorization service at the time of the deployment. Such an authorization service
requires low-latency / high-availability SLOs to avoid deployment outage. This is exarcebated in systems
like Kubernetes where admission webhooks run for every pod deployed. Thus it is often desirable
to "shift-left" and perform an authorization evaluation ahead of time _before_ a deployment request
reaches the control plane. The deployment attestation _is_ the proof of authorization that the control plane may
use to make its final decision, instead of querying an online service itself. 
Verification of the deployment attestation is simple, fast and may be performed entirely offline.
Overall, this shift-left strategy provides the following advantages: less likely to cause production issues,
better debugging UX for devs, less auditing and production noise for SREs and security teams.

## Prerequisites

This predicate depends on the [in-toto Attestation Framework](https://github.com/laurentsimon/attestation/blob/feat/deploy-att/spec/README.md).

## Model

This predicate is for the deployment stage of the software supply chain, where
consumers want to bind an artifact to a deployment environment.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {

    // Required: creation time.
    "creationTime": "...",

    // Optional: decision details.
    "decisionDetails": {
       "evidence": []<ResourceDescriptor>,
       "policy": []<ResourceDescriptor>
     },

    // Optional: scopes.
    "scopes": map[string]string{
        "<scope-name1>/<version>": "value1",
        "<scope-name2>/<version>": "value2"
    }
  }
}

```

### Fields

**`creationTime`, required** string ([Timestamp](https://github.com/in-toto/attestation/blob/main/spec/v1/field_types.md#Timestamp))

The timestamp indicating what time the attestation was created.

**`decisionDetails.evidence`, optional** (list of [ResourceDescriptor](https://github.com/in-toto/attestation/blob/main/spec/v1/resource_descriptor.md))

List of evidence used to make a decision. Resources may include attestations or other relevant evidence.

**`decisionDetails.policy`, optional** (list of [ResourceDescriptor](https://github.com/in-toto/attestation/blob/main/spec/v1/resource_descriptor.md))

List of policies used to make the decision, if any.

**`scopes`, optional** map of string to string (scope type to scope value)

A set of protection scopes of different types. A protection scope identifies the deployment environment to be protected and binds the attestation subject (image, artifact) to it.
A protection scope SHOULD identify the resources to be protected explicitly via their value (e.g., a service account, a Spiffe ID, a Kubernetes pod ID). A protection scope MAY identify the resource implicitly via an authorization / policy URI that hides these details.
A scope has a type and a value. A type ends with its version encoded with `/version`, such as `/v1`. Examples of explicit scope types include Kubernetes's objects such as a pod's cluster ID or a GCP service account. Examples of implicit scope types incude [Google Cloud Binauthz](https://cloud.google.com/binary-authorization/) policy URIs.
Let's see some examples:

- If we want to "restrict an image to run only on GKE cluster X", the scope type is a "Kubernetes cluster" and "X" is the scope value.
- If we want to "restrict an image to run only under service account Y", the scope type is "service account" and the scope value is "Y".
- If we want to "run an image only if it has been vuln scanned", the scope is empty in the sense that the image is allowed to run in any target environment.

### Parsing Rules

This predicate follows the in-toto attestation [parsing rules](https://github.com/in-toto/attestation/blob/main/spec/v1/README.md#parsing-rules).
Summary:

- Consumers MUST ignore the `decisionDetails` field during verification. The field is purely informational and is intended only
for troubleshooting and logging.
- Consumers MUST reject attestations with scope types they do not recognize.
- The `predicateType` URI includes the major version number and will always change whenever there is a backwards incompatible change.
- Minor version changes are always backwards compatible and "monotonic". Such changes do not update the `predicateType`.
- Producers MAY add custom scope types to the `scopes` field. To avoid type name collisions, a scope type MUST be an URI. See [custom scopes](#custom-scopes).

### Verification

#### Configuration
Verification of a deployment attestation is typically performed by an admission controller prior to deploying an artifact / container.
The verification configuration MUST be done out-of-band and contain the following pieces of information:

1. Required: The "trusted roots". A trusted root defines an entity that is trusted to generate attestations. A trusted root MUST be configured with at least the following pieces of information:
    - Required: The unique identity of the attestation generator. The identity may be a cryptographic public key, an identity in an x509 certificate, etc.
    - Required: Which scope types the attestation generator is authoritative for.
2. Optional: Required scopes, which is a set of mandatory scope types that MUST be non-empty for verification to pass. Images MUST have attestation(s) over each scope type in the set in order to be admitted. Required scopes are necessary in an attestation, but not sufficient; other scopes present in the attestation MUST match the current environment in order for it to be considered valid.
3. Optional: Required URI for each scope type that identifies resources implicitely (see [Schema](#schema)).

#### Logic

Verification happens in two phases:

1. Attestation authenticity verification. It takes as input an artifact, an attestation, an attestation signature and the trusted roots. If verification passes, it outputs the attestation's intoto statement. If verification fails, the attestation is considered invalid and MUST be rejected. Using the trusted roots, this phase verifies:
    - The attestation signature.
    - Only the authoritative scopes are present in the attestation. If a scope type is non-empty and the generator is _not_ authoritative for the scope type, verification MUST fail. If a scope type is unrecognized or not supported by the verifier, verification MUST fail.
    - Required scopes are present in the attestation.
2. Scope match verification. It takes as input the intoto payload from the previous phase. For scope types that identify a resource explicitly (see [Schema](#schema)), the verifier matches each scope value against its corresponding environment value where the artifact is to be deployed
(e.g., a service account, a pod ID). For scope types that identify a resource implicitly via an authorization URI (see [Schema](#schema)), the verifier matches the value against the URI in the configuration (See [Configuration](#configuration)). Non-empty fields add constraints to the protection scope and are _always_ interpreted as a logical "AND". The verifier MUST compare each scope value to its expected value using an equality comparison. If the values are all equal, verification passes. Otherwise, it MUST fail. Unset scopes (either a scope type with an empty value or a non-present scope) are interpreted as "any value" and are ignored.

### Supported Scopes

#### Kubernetes pod scope

The specifications defines the Kubernetes's pod scope as follows:

```shell
kubernetes.io/pod/service_account/v1  string: A k8 service account
kubernetes.io/pod/cluster_id/v1       string: A cluster ID
kubernetes.io/pod/namespace/          string: A namespace
kubernetes.io/pod/cluster_name/v1     string: A cluster name
```

The scope match verification compares each non-empty scope value against the corresponding
environment value the artifact is to be deployed:

```shell
attestation's "kubernetes.io/pod/service_account" == environment's "k8's service_account" AND
attestation's "kubernetes.io/pod/cluster_id" == environment's "k8'scluster_id" AND 
...
```

#### GCP scope

The specifications defines the GCP scope as follows:

```shell
cloud.google.com/service_account/v1 string: A GCP service account
cloud.google.com/location/v1        string: A location
cloud.google.com/project_id/v1      string: A project id
```

The scope match verification compares each non-empty scope value against the corresponding
environment value the artifact is to be deployed:


```shell
attestation's "cloud.google.com/service_account" == environment's "GCP service_account" AND
attestation's "cloud.google.com/project_id" == environment's "GCP project ID" AND 
...
```

#### Spiffe

The specifications define the Spiffe scope as follows:

```shell
spiffe.io/id/v1 string: The Spiffe ID
```

The scope match verification compares each non-empty scope value against the corresponding
environment value the artifact is to be deployed:


```shell
attestation's "spiffe.io/id" == environment's "Spiffe ID" AND
...
```

#### Custom scopes

One can define their own scopes. To avoid scope type name collisions, the scope type name MUST b a unique URI, such as:

```shell
my.myproject.com/resource/v1   string: resource for environment my.myproject.com
```

If a scope type is unrecognized or not supported by the verifier, verification MUST fail.

## Examples

### Example 1: Single scope

**Trusted (single) root configuration**:
- A public key (ignored by the example)
- Authoritative scope types: "cloud.google.com/service_account"
- Required scope types: "cloud.google.com/service_account"

**Deployment request**:
- For a container running under service account "sa-name@project.iam.gserviceaccount.com"
- With the following attestation matching the container's digest in the request:

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {
    "creationTime": "...",
    "scopes": {
      "cloud.google.com/service_account/v1": "sa-name@project.iam.gserviceaccount.com"
    }
}
```

**Verification result**:
- The attestation authentication verification passes, because the trusted root
  is authoritative for scope type "cloud.google.com/service_account/v1".
- The scope match verification passes, because the value of the scope matches
  the value of the environment "sa-name@project.iam.gserviceaccount.com".

### Example 2: Non-authoritative scope

**Trusted (single) root configuration**:
- A public key (ignored by the example)
- Authoritative scope types: "cloud.google.com/project_id"
- Required scope types: "cloud.google.com/project_id"

**Deployment request**:
- For a container running under GCP service account "sa-name@project.iam.gserviceaccount.com"
- With the following attestation matching the container's digest in the request:

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {
    "creationTime": "...",
    "scopes": {
      "cloud.google.com/service_account/v1": "sa-name@project.iam.gserviceaccount.com"
    }
}
```

**Verification result**:
- The attestation authentication verification fails, because the trusted root
  is _not_ authoritative for scope type "cloud.google.com/service_account/v1".
- The scope match verification is _not_ performed, because the authentication
verification failed.

### Example 3: Two required scopes

**Trusted (single) root configuration**:
- A public key (ignored by the example)
- Authoritative scope types: "cloud.google.com/service_account" and "kubernetes.io/pod/cluster_id"
- Required scope types: "cloud.google.com/service_account" and "kubernetes.io/pod/cluster_id"

**Deployment request**:
- For a container running in cluster ID "unique-cluster-id" under GCP service account "sa-name@project.iam.gserviceaccount.com"
- With the following attestation matching the container's digest in the request:

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {
    "creationTime": "...",
    "scopes": {
      "cloud.google.com/service_account/v1": "sa-name@project.iam.gserviceaccount.com",
      "kubernetes.io/pod/cluster_id/v1": "unique-cluster-id"
    }
}
```

**Verification result**:
- The attestation authentication verification succeeds, because the trusted root
  is authoritative for scope types "cloud.google.com/service_account/v1" and "kubernetes.io/pod/cluster_id/v1"; and both scope types are present in the attestation.
- The scope match verification succeeds, because the  environment the container is to be deployed matches the scope values.

### Example 4: Single required scope

**Trusted (single) root configuration**:
- A public key (ignored by the example)
- Authoritative scope types: "cloud.google.com/service_account" and "kubernetes.io/pod/cluster_id"
- Required scope types: "cloud.google.com/service_account"

**Deployment request**:
- For a container running under GCP service account "sa-name@project.iam.gserviceaccount.com"
- With the following attestation matching the container's digest in the request:

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {
    "creationTime": "...",
    "scopes": {
      "cloud.google.com/service_account/v1": "sa-name@project.iam.gserviceaccount.com"
    }
}
```

**Verification result**:
- The attestation authentication verification succeeds, because the trusted root
  is authoritative for scope type "cloud.google.com/service_account/v1". The only required scope type is "cloud.google.com/service_account" and is present in the attestation.
- The scope match verification succeeds, because the  environment the container is to be deployed matches the scope values.

### Example 5: Multiple roots

**Trusted root configurations**:
- Root 1:
  - A public key (ignored by the example)
  - Authoritative scope types: "cloud.google.com/service_account"
  - Required scope types: "cloud.google.com/service_account"
- Root 2:
  - A public key (ignored by the example)
  - Authoritative scope types: "kubernetes.io/pod/cluster_id"
  - Required scope types: "kubernetes.io/pod/cluster_id"

**Deployment request**:
- For a container running in cluster ID "unique-cluster-id" under GCP service account "sa-name@project.iam.gserviceaccount.com"
- With the following attestations matching the container's digest in the request:

```jsonc
{
  // Assumption: Signed by Root 1
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {
    "creationTime": "...",
    "scopes": {
      "cloud.google.com/service_account/v1": "sa-name@project.iam.gserviceaccount.com",
    }
}
```

```jsonc
{
  // Assumption: Signed by Root 2
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {
    "creationTime": "...",
    "scopes": {
      "kubernetes.io/pod/cluster_id": "unique-cluster-id",
    }
}
```
**Verification result**:
- The attestation authentication verification succeeds, because:
  - Root 1 is authoritative for scope type "cloud.google.com/service_account/v1" and the scope is present in the (first) attestation signed by Root 1.
  - Root 2 is authoritative for scope type "kubernetes.io/pod/cluster_id" and the scope is present in the (second) attestation signed by Root 2.
- The scope match verification succeeds, because:
  - Root 1: The service account (in the first attesttation) matches the environment the container is to be deployed
  - Root 2: The cluster ID (in the second attesttation) matches the environment the container is to be deployed.

### Example 6: Implicit scope type

**Trusted (single) root configuration**:
- A public key (ignored by the example)
- Authoritative scope types: "binaryauthorization.googleapis.com/policy_uri" and "kubernetes.io/pod/namespace"
- Required scope types: "binaryauthorization.googleapis.com/policy_uri" with value "projects/foo/platforms/gke/policies/bar".

**Deployment request**:
- For a container running in a Kubernetes namespace "prod-namespace"
- With the following attestation matching the container's digest in the request:

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {
    "creationTime": "...",
    "scopes": {
      "kubernetes.io/pod/namespace/v1": "prod-namespace",
      "binaryauthorization.googleapis.com/policy_uri": "projects/foo/platforms/gke/policies/bar"
    }
}
```

**Verification result**:
- The attestation authentication verification succeeds, because the trusted root
  is authoritative for scope type "kubernetes.io/pod/namespace/v1" and "binaryauthorization.googleapis.com/policy_uri". The only required scope type is "binaryauthorization.googleapis.com/policy_uri" and is present in the attestation.
- The scope match verification succeeds, because the namespace matches the environment and
the "binaryauthorization.googleapis.com/policy_uri" value matches the configuration value "projects/foo/platforms/gke/policies/bar".

### Example 7: Unrecognized scope

**Trusted (single) root configuration**:
- A public key (ignored by the example)
- Authoritative scope types: "cloud.google.com/service_account"
- Required scope types: "cloud.google.com/service_account"

**Deployment request**:
- For a container running under service account "sa-name@project.iam.gserviceaccount.com"
- With the following attestation matching the container's digest in the request:

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {
    "creationTime": "...",
    "scopes": {
      "cloud.google.com/service_account/v1": "sa-name@project.iam.gserviceaccount.com",
      "my.custom-scope.com/some-field/v1": "some-value"
    }
}
```

**Verification result**:
- The attestation authentication verification fails, because the trusted root
  is _not_ authoritative for the unrecognized scope type "my.custom-scope.com/some-field/v1".
- The scope match verification is _not_ performed, because the authentication
verification failed.

### Example 8: No scope

**Trusted (single) root configuration**:
- A public key (ignored by the example)
- Authoritative scope types: "cloud.google.com/service_account"
- Required scope types: none

**Deployment request**:
- For a container running under service account "sa-name@project.iam.gserviceaccount.com"
- With the following attestation matching the container's digest in the request:

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  "predicateType": "https://in-toto.io/attestation/deployment/v1",
  "predicate": {
    "creationTime": "...",
}
```

**Verification result**:
- The attestation authentication verification passes, because no scopes are required.
- The scope match verification passes, because there is _no_ verification to perform.
