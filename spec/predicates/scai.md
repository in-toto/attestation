# Predicate type: Software Supply Chain Attribute Integrity (SCAI)

Type URI: scai/attribute-assertion/v0.1

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

SCAI metadata, referred to as **Attribute Assertions**, comprises a set of statements that describe functional
attributes of a software artifact and its supply chain, capable of covering the full software stack of the
toolchain that produced the artifact down to the hardware platform. SCAI Attribute Assertions include
information about the conditions under which certain functional attributes arise, as well as (authenticated)
evidence for the asserted attributes.

SCAI is intended to be implemented as part of an existing software supply chain attestation
framework by software development tools or services (e.g., builders, CI/CD pipelines, software analysis tools)
seeking to capture more granular information about the attributes and behavior of the software artifacts they
produce. As such, we envision SCAI metadata being explicitly bound to, or embedded within, other metadata
objects.

## Schema

```
{
    "predicateType": "scai/attribute-assertion/v0.1",
    "predicate": {
        "attributes": [{
            "attribute": "<ATTRIBUTE>",
            "target": { // optional
                "name": "<NAME>",
                "digest": { "<ALGORITHM>": "<HEX VALUE>", ... },
                "locationURI": "<URI>",
                "objectType": "<OBJECT TYPE>" // optional
            },
            "conditions": { /*object */ }, // optional
            "evidence": { /* object */ } // optional
        }]
    }
}
```

This predicate has been adapted from the [SCAI v0.1 specification](https://arxiv.org/pdf/2210.05813.pdf) for greater interoperability.

### Parsing Rules

At a high-level, Attribute Assertions MUST allow humans and programs to easily parse the asserted
attributes. Additional fields MUST enable program-based consumers to automatically parse and evaluate
the given information.

The following parsing rules apply in addition:
* Consumers MUST ignore unrecognized fields.
* Producers SHOULD omit _optional_ fields when unused to avoid ambiguity.
* Acceptable formats of the `attribute`, `conditions` and `evidence` fields are up to the producer and consumer.
* Because consumers evaluate this predicate against a policy, the semantics SHOULD be consistent and monotonic between Attestations (see in-toto Attestation Spec [parsing rules](https://github.com/in-toto/attestation/tree/main/spec#parsing-rules)).

### Fields

`predicateType` _string ([TypeURI](https://github.com/in-toto/attestation/blob/main/spec/field_types.md#TypeURI)), required_

> Identifier for the schema of the Assertion. Always
> `scai/attribute-assertion/v0.1` for this version of the spec.

`predicate.attributes` _array of objects, required_

> An array of one or more SCAI Attribute Assertions about the subject.

`predicate.attributes[*].attribute` _string, required_

> A string describing a specific functional feature of the Attestation subject.
>
> Attributes are expected to be domain- or application-specific.

`predicate.attributes[*].target` _object ([SCAI Object Reference](https://arxiv.org/pdf/2210.05813.pdf), optional_

> An object reference to a specific artifact or metadata object to which the `attribute` field applies. 
>
> The semantics of the optional `objectType` field are up to the producer and consumer.

`predicate.attributes[*].conditions` _object, optional_

> An object representing specific conditions under which the associated attribute arises.

`predicate.attributes[*].evidence` _object, optional_

> An object representing (authenticated) evidence for the asserted `attribute`.
>
> Regardless of format, the `evidence` object MUST be self-describing (i.e, have a [TypeURI](https://github.com/in-toto/attestation/blob/main/spec/field_types.md#TypeURI)) to
> enable the consumer to evaluate the object against a policy.
>
> When omitted, a consumer may choose to evaluate the Attestation
> on the basis of the producer's identity.

## Examples

#### Attestation for gcc compiler attributes
```
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v0.1",
    "subject": [{
        "name": "gcc9.3.0",
        "digest": { "sha256": "78ab6a8..." },
        "locationURI": "http://us.archive.ubuntu.com/ubuntu/pool/main/g/gcc-defaults/gcc_9.3.0-1ubuntu2_amd64.deb"
    }],
        
    "predicateType": "scai/attribute-assertion/v0.1",
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
                "objectType": "https://in-toto.io/link/v0.1"
            }
        }]
    }
}
```

#### Attestation for gcc compilation with *endorsed* compiler flags
```
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v0.1",
    "subject": [{
        "name": "hello-world",
        "digest": { "sha256": "78ab6a8..." },
        "locationURI":  "http://example.com/binaries/hello-world"
    }],
        
    "predicateType": "scai/attribute-assertion/v0.1",
    "predicate": {
        "attributes": [{
            "attribute": "WITH_STACK_PROTECTION",
            "conditions": { "flags": "-fstack-protector*" },
            "evidence": {
                "_type": "https://github.com/secure-systems-lab/dsse",
                "payloadType": "application/vnd.in-toto+json",
                "payload": "eyJfdHlwZSI6ICJodHRwczovL2l...",
                "signatures": [{ "sig": "MEQCIAZjdOJnQddF14Rpq..." }]
            }
        }]
    }
}
```

#### Attestation for build with SLSA Provenance (hypothetical)
```
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v0.1",
    "subject": [{
        "name": "hello-world",
        "digest": { "sha256": "78ab6a8..." },
        "locationURI":  "http://example.com/binaries/hello-world"
    }],
        
    "predicateType": "https://slsa.dev/provenance/vX"
    "predicate": {
        "buildType": "https://example.com/sources/Makefile",
        "builder": { "id": "mailto:hi@example.com" },
        "invocation": {
            "configSource": {
                "uri": "https://example.com/sources/",
                "digest": { "sha256": "cdcdcdc..." }
            },
            "parameters": { "CFLAGS": "-fstack-protector" }
        },
        "materials": [{
            "uri": "https://example.com/sources/hello-world.c",
            "digest": { "sha256": "e5564d4..." }
        }],

       // hypothetical field for SCAI
       "subjectAttributes": [{
           "attribute": "WITH_STACK_PROTECTION",
           "conditions": { "flags": "-fstack-protector*" },
           "evidence": {
               "_type": "https://github.com/secure-systems-lab/dsse",
               "payloadType": "application/vnd.in-toto+json",
               "payload": "eyJfdHlwZSI6ICJodHRwczovL2l...",
               "signatures": [{ "sig": "MEQCIAZjdOJnQddF14Rpq..." }]
           }
       }]
    }
}
```

#### Attestation for build on Intel(R) SGX hardware with SLSA Provenance (hypothetical)
```
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v0.1",
    "subject": [{
        "name": "some-security-critical-app",
        "digest": { "sha256": "abcdefa..." },
        "locationURI":  "http://example.com/binaries/some-security-critical-app"
    }],
        
    "predicateType": "https://slsa.dev/provenance/vX"
    "predicate": {
        "buildType": "https://example.com/sources/Makefile",
        "builder": { "id": "mailto:hi@example.com" },
        "invocation": {
            "configSource": {
                "uri": "https://example.com/sources/",
                "digest": { "sha256": "cdcdcdc..." }
            },
            "parameters": { "CFLAGS": "-fstack-protector" }
        },
        "materials": [{
            "uri": "https://example.com/sources/some-app.c",
            "digest": { "sha256": "e456f78..." }
        }],

       // hypothetical field for SCAI
       "metadata": {
           "attribute": "ATTESTED_HARDWARE",
           "target": {
               "name": "enclave.signed.so",
               "digest": { "sha256": "e3b0c44..." },
               "locationURI": "http://example.com/enclaves/enclave.signed.so",
           },
           "evidence": {
               "_type": "https://download.01.org/intel-sgx/sgx-dcap/1.14/linux",
               // Intel SGX DCAP attestation
               "quote_body": {
                   "version": 3,
                   "sign_type": 2,
                   "..."
               },
               "report_body": {
                   "cpu_svn": "0c0dfff...",
                   "mr_enclave": "0a1b2c3...",
                   "mr_signer": "a1b2c3d...",
                   "report_data":"9876543...",
                   "..."
               },
               "signature_size": 4163,
               "signature": "0123456..."
           }
       }
    }
}
```