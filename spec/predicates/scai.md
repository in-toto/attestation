# Predicate type: Software Supply Chain Attribute Integrity (SCAI)

Type URI: scai/attribute-report/v0.1

Version: 0.1

Author: Marcela Melara <marcela.melara@intel.com>

## Purpose

The Software Supply Chain Attribute Integrity, or SCAI (pronounced "sky"), specification proposes a data
format for capturing functional attribute and integrity information about software artifacts and their supply
chain. SCAI data can be associated with executable binaries, statically- or dynamically-linked libraries,
software packages, container images, software toolchains, and compute environments.

Existing supply chain data formats do not capture any information about the security functionality or
behavior of the resulting software artifact, nor do they provide sufficient evidence to support any claims of
integrity of the supply chain processes they describe. The SCAI data format is designed to bridge this gap.

For more details, see the [SCAI v0.1 specification document](https://arxiv.org/pdf/2210.05813.pdf)

## Prerequisites

In addition to the in-toto attestation framework,
SCAI assumes that implementers have appropriate processes and tooling in place for
capturing other types of software supply chain metadata, which can be extended to add support for SCAI.

## Model

SCAI metadata, referred to as an Attribute Assertion, describes functional attributes of a software
artifact and its supply chain, capable of covering the full software stack of the
toolchain that produced the artifact down to the hardware platform. SCAI Attribute Assertions include
information about the conditions under which certain functional attributes arise, as well as (authenticated)
evidence for the asserted attributes. The set of Assertions about a subject artifact and its producer
is referred to as the **Attribute Report**.
Similarly, SCAI Attribute Reports about the producer of a subject artifact can be
generated separately, with the attestation subject indicating an artifact producer.

SCAI is intended to be implemented as part of an existing software supply chain attestation
framework by software development tools or services (e.g., builders, CI/CD pipelines, software analysis tools)
seeking to capture more granular information about the attributes and behavior of the software artifacts they
produce.
As such, we envision SCAI metadata being explictly bound to, or included within, other metadata
objects; we recommend [in-toto Attestation Bundles](https://github.com/in-toto/attestation/blob/main/spec/bundle.md) for this purpose.

## Schema

The core metadata in SCAI is the Attribute Assertion. A collection of Attribute Assertions
for a specific supply chain step or operation are issued together in a SCAI Attribute Report predicate.

```
{
    "predicateType": "https://in-toto.io/scai/attribute-report/v0.1",
    "predicate": {
        "attributes": [{
            "attribute": "<ATTRIBUTE>",
            "target": { // optional
                /* SCAI Object Reference */
                "name": "<NAME>",
                "digest": { "<ALGORITHM>": "<HEX VALUE>", ... },
                "locationURI": "<RESOURCE URI>",
                "objectType": "<OBJECT TYPE>" // optional
            },
            "conditions": { /* object */ }, // optional
            "evidence": { /* SCAI Object Reference */ } // optional
        }],
        "producer": {
            "type": "<TYPE URI>",
            "reference": { // optional
                "name": "<NAME>",
                "digest": { "<ALGORITHM>": "<HEX VALUE>", ... },
                "locationURI": "<RESOURCE URI>",
                "objectType": "<OBJECT TYPE>" // optional
            }
        }
    }
}
```

This predicate has been adapted from the [SCAI v0.1 specification](https://arxiv.org/pdf/2210.05813.pdf) for greater interoperability.

### Parsing Rules

At a high-level, Attribute Reports MUST allow humans and programs to easily parse the asserted
attributes. Additional fields MUST enable program-based consumers to automatically parse and evaluate
the given information.

The following parsing rules apply in addition:
* Consumers MUST ignore unrecognized fields.
* Producers SHOULD omit _optional_ fields when unused to avoid ambiguity.
* Acceptable formats of the `attribute` and `conditions` fields are up to the producer and consumer.
* Because consumers evaluate this predicate against a policy, the semantics SHOULD be consistent and monotonic between Attestations (see in-toto Attestation Spec [parsing rules](https://github.com/in-toto/attestation/tree/main/spec#parsing-rules)).

### Fields

`predicateType` _string ([TypeURI](https://github.com/in-toto/attestation/blob/main/spec/field_types.md#TypeURI)), required_

> Identifier for the schema of the Attribute Report. Always
> `scai/attribute-report/v0.1` for this version of the spec.

`predicate.attributes` _array of objects, required_

> An array of one or more SCAI Attribute Assertions about the subject.

`predicate.attributes[*].attribute` _string, required_

> A string describing a specific functional feature of the Attestation subject or producer.
>
> Attributes are expected to be domain- or application-specific.

`predicate.attributes[*].target` _object ([SCAI Object Reference](https://arxiv.org/pdf/2210.05813.pdf), optional_

> An object reference to a specific artifact or metadata object to which the `attribute` field applies. 
>
> The semantics of the optional `objectType` field are up to the producer and consumer.

`predicate.attributes[*].conditions` _object, optional_

> An object representing specific conditions under which the associated attribute arises.

`predicate.attributes[*].evidence` _object ([SCAI Object Reference](https://arxiv.org/pdf/2210.05813.pdf), optional_

> An object reference to (authenticated) evidence for the asserted `attribute`.
>
> If the evidence object is generated by the producer in parallel to the SCAI predicate
> the producer MAY include the attestation for the evidence object in an Attestation Bundle,
> and omit the `locationURI` field.
> The semantics of the optional `objectType` field are up to the producer and consumer.
>
> When `evidence` is omitted, a consumer MAY choose to evaluate the Attestation
> on the basis of the producer's identity.

`predicate.producer` _object, optional_

> An object identifying the producer of the Attestation subject.

`predicate.producer.type` _string ([TypeURI](https://github.com/in-toto/attestation/blob/main/spec/field_types.md#TypeURI)), required_

> A URI describing the type of producer of the Attestation subject.

`predicate.producer.reference` _object ([SCAI Object Reference](https://arxiv.org/pdf/2210.05813.pdf), optional_

> An object reference to the specific producer of the Attestation subject, if applicable. 
>
> The semantics of the optional `objectType` field are up to the producer and consumer.

## Examples

#### Attestation for gcc compiler attributes
```
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v0.1",
    "subject": [{
        "name": "gcc9.3.0",
        "digest": { "sha256": "78ab6a8..." }
    }],
        
    "predicateType": "scai/attribute-report/v0.1",
    "predicate": {
        "attributes": [{
            "attribute": "WITH_STACK_PROTECTION",
            "conditions": { "flags": "-fstack-protector*" }
        },
        {
            "attribute": "REPRODUCIBLE",
            "evidence": {
                "name": "rebuilderd-attestation",
                "digest": { "sha256": "abcdabcde..." },
                "locationURI": "http://example.com/rebuilderd-instance/gcc_9.3.0-1ubuntu2_amd64.att",
                "objectType": "application/vnd.in-toto+json"
            }
        }]
    }
}
```

#### Attestation for binary attributes
```
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v0.1",
    "subject": [{
        "name": "my-app",
        "digest": { "sha256": "78ab6a8..." }
    }],
        
    "predicateType": "scai/attribute-report/v0.1",
    "predicate": {
        "subjectAttributes": [{
            "attribute": "WITH_STACK_PROTECTION",
            "conditions": { "flags": "-fstack-protector*" },
        }],
        "producer": {
            "type": "file:/usr/bin/gcc", 
            "reference": {
                "name": "gcc9.3.0",
                "digest": {
                    "sha256": "78ab6a8..."
                },
                "locationURI": "http://us.archive.ubuntu.com/ubuntu/pool/main/g/gcc-defaults/gcc_9.3.0-1ubuntu2_amd64.deb"
            }
        }
    }
}
```

#### Attestation for binary attributes with evidence
```
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v0.1",
    "subject": [{
        "name": "my-app",
        "digest": { "sha256": "78ab6a8..." }
    }],
        
    "predicateType": "scai/attribute-report/v0.1",
    "predicate": {
        "subjectAttributes": [{
            "attribute": "WITH_STACK_PROTECTION",
            "conditions": { "flags": "-fstack-protector*" },
            "evidence": {
                "name": "gcc9.3.0-attribute-report",
                "digest": { "sha256": "abcdabcde..." },
                "locationURI": "http://example.com/rekor-instance",
                "objectType": "application/vnd.in-toto+json"
            }
        }],
        "producer": {
            "type": "file:/usr/bin/gcc", 
            "reference": {
                "name": "gcc9.3.0",
                "digest": {
                    "sha256": "78ab6a8..."
                },
                "locationURI": "http://us.archive.ubuntu.com/ubuntu/pool/main/g/gcc-defaults/gcc_9.3.0-1ubuntu2_amd64.deb"
            }
        }
    }
}
```

#### Attestation for binary with attested dependencies
```
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v0.1",
    "subject": [{
        "name": "my-app",
        "digest": { "sha256": "78ab6a8..." }
    }],
        
    "predicateType": "scai/attribute-report/v0.1",
    "predicate": {
        "attributes": [{
            "attribute": "ATTESTED_DEPENDENCIES",
            "target": {
                "name": "my-rsa-lib.so",
                "digest": { "sha256": "ebebebe..." },
                "locationURI": "http://example.com/libraries/my-rsa-lib.so"
            }
            "evidence": {
                "name": "rsa-lib-attribute-report",
                "digest": { "sha256": "0987654..." },
                "locationURI": "http://example.com/rekor-instance",
                "objectType": "application/vnd.in-toto+json"
            }
        }],
        "producer": {
            "type": "https://example.com/my-github-actions-runner", 
        }
    }
}
```

#### Attestation for build on Intel(R) SGX hardware
```
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v0.1",
    "subject": [{
        "name": "my-sgx-builder",
        "digest": { "sha256": "78ab6a8..." }
    }],
        
    "predicateType": "scai/attribute-report/v0.1"
    "predicate": {
        "attributes": [{
            "attribute": "HARDWARE_ENCLAVE",
            "target": {
                "name": "enclave.signed.so",
                "digest": { "sha256": "e3b0c44..." },
                "locationURI": "http://example.com/enclaves/enclave.signed.so",
            },
            "evidence": {
                "name": "my-sgx-builder.att",
                "digest": { "sha256": "0987654..." },
                "locationURI": "http://example.com/sgx-attestations/my-sgx-builder.att",
                "objectType": "https://download.01.org/intel-sgx/sgx-dcap/1.14/linux"
            }
       }]
    }
}
```