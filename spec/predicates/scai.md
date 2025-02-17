# Predicate type: Software Supply Chain Attribute Integrity (SCAI)

Type URI: https://in-toto.io/attestation/scai

Version: 0.4

Author: Marcela Melara ([@marcelamelara](https://github.com/marcelamelara))

## Purpose

The Software Supply Chain Attribute Integrity, or SCAI (pronounced "sky"),
specification proposes a data format for capturing functional attribute and
integrity information about software artifacts and their supply chain. SCAI
data can be associated with executable binaries, statically- or dynamically-
linked libraries, software packages, container images, software toolchains,
and compute environments.

Existing supply chain data formats do not capture any information about the
security functionality or behavior of the resulting software artifact, nor
do they provide sufficient evidence to support any claims of integrity of
the supply chain processes they describe. The SCAI data format is designed
to bridge this gap.

For more details, see the [SCAI specification] document.

## Use Cases

We highlight three core use cases of SCAI, although there may be others.

### Fine-grained software attributes

A producer seeks to make assertions about fine-grained code-level attributes
of the attestation `subject`. Examples of such code-level attributes include
include cases like code features instrumented by a compiler, or properties
tested via static analysis.

See an example for [attested binary attributes](#attestation-for-binary-attributes-with-evidence) below.

### Producer compute platform attributes

A producer seeks to assert certain properties about its underlying compute
platform, which produced the `subject`. Examples of compute platform
attributes include properties like secure boot enablement or trusted
execution enablement. Since these types of attributes are often accompanied
by an attestation document, these hardware attestations serve as evidence
for the asserted attributes.

See an example for a [build on attested hardware](#attestation-for-build-on-intelr-sgx-hardware) below.

### Evidence collection

A producer seeks to assert that a collection of evidence or attestations
applies to the `subject`. In essence, this use case serves to bind multiple
attestations to a single `subject`. For example, the asserted attributes
may correspond to specific attestation types.

See an example for an [evidence collection](#attestation-for-evidence-collection) below.

## Prerequisites

In addition to the in-toto Attestation Framework, SCAI assumes that
implementers have appropriate processes and tooling in place for capturing
other types of software supply chain metadata, which can be extended to add
support for SCAI.

## Model

SCAI metadata, referred to as an Attribute Assertion, describes functional
attributes of a software artifact and its supply chain, capable of covering
the full software stack of the toolchain that produced the artifact down to
the hardware platform. SCAI Attribute Assertions include information about
the conditions under which certain functional attributes arise, as well as
(authenticated) evidence for the asserted attributes.

The set of Assertions about a subject artifact and its producer is referred
to as the **Attribute Report**. Similarly, SCAI Attribute Reports about the
producer of a subject artifact can be generated separately, with the
attestation subject indicating an artifact producer.

We emphasize that this predicate does not dictate the types of attributes
that can be attested, nor does it prescribe any format to the evidence for an
assertion.

SCAI is intended to be implemented as part of an existing software supply
chain attestation framework by software development tools or services (e.g.,
builders, CI/CD pipelines, software analysis tools) seeking to capture more
granular information about the attributes and behavior of the software
artifacts they produce.

As such, we envision SCAI metadata being explictly bound to, or included
within, other metadata objects; we recommend an in-toto [attestation Bundle]
for this purpose.

## Data description

The core metadata in SCAI is the Attribute Assertion. A collection of
Attribute Assertions for a specific supply chain step or operation are issued
together in a SCAI Attribute Report predicate.

The following CDDL depends on definitions from the base v1.2 in-toto specification.

```cddl
scai-predicate = (
  predicateType-label => "https://in-toto.io/attestation/scai/v0.3",
  predicate-label => {
      attributes-label: [ * scai-attribute-map ],
      ? producer-label => ResourceDescriptor
    }
)
scai-attributes-label = JC<"attributes", 0>
scai-producer-label   = JC<"producer",   1>

scai-attribute-map = {
  attribute-label => text,
  ? target-label => ResourceDescriptor,
  ? conditions-label => object,
  ? evidence-label => ResourceDescriptor
}
scai-attribute-label  = JC<"attribute",  0>
scai-target-label     = JC<"target",     1>
scai-conditions-label = JC<"conditions", 2>
scai-evidence-label   = JC<"evidence",   3>
```

This predicate has been adapted from the [SCAI specification] for greater
interoperability.

### Parsing Rules

At a high-level, Attribute Reports MUST allow humans and programs to easily
parse the asserted attributes. Additional fields MUST enable program-based
consumers to automatically parse and evaluate the given information.

The following parsing rules apply in addition:

-   Consumers MUST ignore unrecognized fields.
-   Producers SHOULD omit _optional_ fields when unused to avoid ambiguity.
-   Acceptable formats of the `attribute` and `conditions` fields are up to
    the producer and consumer.
-   Because consumers evaluate this predicate against a policy, the semantics
    SHOULD be consistent and monotonic between attestations (see in-toto
    Attestation Framework [parsing rules]).

### Fields

`predicateType` _string ([TypeURI]), required_

> Identifier for the schema of the Attribute Report. Always
> `https://in-toto.io/attestation/scai/v0.3` for this version of the
> spec.

`predicate.attributes` _array of objects, required_

> An array of one or more SCAI Attribute Assertions about the subject.

`predicate.attributes[*].attribute` _string, required_

> A string describing a specific functional feature of the attestation
> subject or producer.
>
> Attributes are expected to be domain- or application-specific.

`predicate.attributes[*].target` _object ([ResourceDescriptor]), optional_

> An object reference to a specific artifact or metadata object to which the
> `attribute` field applies.
>
> The producer and consumer SHOULD agree on the ResourceDescriptor fields
> needed for identification and validation of the target.

`predicate.attributes[*].conditions` _object, optional_

> An object representing specific conditions under which the associated
> attribute arises.

`predicate.attributes[*].evidence` _object ([ResourceDescriptor]), optional_

> A description of (authenticated) evidence for the asserted `attribute`.
>
> If the evidence object is generated by the producer in conjunction with the
> SCAI predicate the producer MAY include the attestation for the evidence
> object in the same in-toto [attestation Bundle],
>
> The producer and consumer SHOULD agree on the ResourceDescriptor fields
> needed for identification and validation of the evidence.
>
> When `evidence` is omitted, a consumer MAY choose to evaluate the
> atestation on the basis of the producer's identity.

`predicate.producer` _object, ([ResourceDescriptor]) optional_

> A description of the producer of the attestation subject, if applicable.
>
> The producer and consumer SHOULD agree on the ResourceDescriptor fields
> needed for identification and validation of the producer.

## Examples

### Attestation for binary attributes

```jsonc
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "name": "my-app",
        "digest": { "sha256": "78ab6a8..." }
    }],

    "predicateType": "https://in-toto.io/attestation/scai/v0.3",
    "predicate": {
        "attributes": [{
            "attribute": "WITH_STACK_PROTECTION",
            "conditions": { "flags": "-fstack-protector*" },
        }],
        "producer": {
            "uri": "file:///usr/bin/gcc",
            "name": "gcc9.3.0",
            "digest": {
                "sha256": "78ab6a8..."
            },
            "downloadLocation": "http://us.archive.ubuntu.com/ubuntu/pool/main/g/gcc-defaults/gcc_9.3.0-1ubuntu2_amd64.deb"
        }
    }
}
```

### Attestation for gcc compiler attributes

```jsonc
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "name": "gcc9.3.0",
        "digest": { "sha256": "78ab6a8..." }
    }],

    "predicateType": "https://in-toto.io/attestation/scai/v0.3",
    "predicate": {
        "attributes": [{
            "attribute": "WITH_STACK_PROTECTION",
            "conditions": { "flags": "-fstack-protector*" }
        },
        {
            "attribute": "REPRODUCIBLE",
            "evidence": {
                "name": "gcc_9.3.0-1ubuntu2_amd64.json",
                "digest": { "sha256": "abcdabcde..." },
                "uri": "http://example.com/rebuilderd-instance/gcc_9.3.0-1ubuntu2_amd64.json",
                "mediaType": "application/x.dsse+json"
            }
        }]
    }
}
```

### Attestation for binary attributes with evidence

```jsonc
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "name": "my-app",
        "digest": { "sha256": "78ab6a8..." }
    }],

    "predicateType": "https://in-toto.io/attestation/scai/v0.3",
    "predicate": {
        "attributes": [{
            "attribute": "WITH_STACK_PROTECTION",
            "conditions": { "flags": "-fstack-protector*" },
            "evidence": {
                "name": "gcc9.3.0-attribute-report.json",
                "digest": { "sha256": "abcdabcde..." },
                "mediaType": "application/x.dsse+json"
            }
        }],
        "producer": {
            "uri": "file:///usr/bin/gcc",
            "name": "gcc9.3.0",
            "digest": {
                "sha256": "78ab6a8..."
            },
            "downloadLocation": "http://us.archive.ubuntu.com/ubuntu/pool/main/g/gcc-defaults/gcc_9.3.0-1ubuntu2_amd64.deb"
        }
    }
}
```

### Attestation for binary with attested dependencies

```jsonc
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "name": "my-app",
        "digest": { "sha256": "78ab6a8..." }
    }],

    "predicateType": "https://in-toto.io/attestation/scai/v0.3",
    "predicate": {
        "attributes": [{
            "attribute": "ATTESTED_DEPENDENCIES",
            "target": {
                "name": "my-rsa-lib.so",
                "digest": { "sha256": "ebebebe..." },
                "uri": "http://example.com/libraries/my-rsa-lib.so"
            }
            "evidence": {
                "name": "rsa-lib-attribute-report.json",
                "digest": { "sha256": "0987654..." },
                "mediaType": "application/x.dsse+json"
            }
        }],
        "producer": {
            "uri": "https://example.com/my-github-actions-runner",
        }
    }
}
```

### Attestation for build on Intel(R) SGX hardware

```jsonc
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "name": "my-sgx-builder",
        "digest": { "sha256": "78ab6a8..." }
    }],

    "predicateType": "https://in-toto.io/attestation/scai/v0.3"
    "predicate": {
        "attributes": [{
            "attribute": "VALID_ENCLAVE",
            "target": {
                "name": "enclave.signed.so",
                "digest": { "sha256": "e3b0c44..." },
                "uri": "http://example.com/enclaves/enclave.signed.so",
            },
            "evidence": {
                "name": "my-sgx-builder.json",
                "digest": { "sha256": "0987654..." },
                "downloadLocation": "http://example.com/sgx-attestations/my-sgx-builder.json",
                "mediaType": "application/x.sgx.dcap1.14+json"
            }
       }]
    }
}
```

### Attestation for evidence collection

```jsonc
{
    // Standard attestation fields
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{
        "name": "app-evidence-collection",
        "digest": { "sha256": "88888888..." }
    }],

    "predicateType": "https://in-toto.io/attestation/scai/v0.3",
    "predicate": {
        "attributes": [{
            "attribute": "attestation-1",
            "evidence": {
                "uri": "https://example.com/attestations/attestation-1"
                "digest": { "sha256": "abcdabcd..." },
                "mediaType": "application/x.dsse+json"
            }
        },
        {
            "attribute": "attestation-2",
            "evidence": {
                "uri": "https://example.com/attestations/attestation-2"
                "digest": { "sha256": "01234567..." },
                "mediaType": "application/x.dsse+json"
            }
        },
        {
            "attribute": "attestation-3",
            "evidence": {
                "uri": "https://example.com/attestations/attestation-3"
                "digest": { "sha256": "deadbeef..." },
                "mediaType": "application/x.dsse+json"
            }
        }],
        "producer": { "uri": "https://my-sw-attestor" }
    }
}
```

## Changelog and Migrations

### New in v0.4

-   Data description updated to use CDDL and provide more concise representation in CBOR.

### New in v0.3

-   Simplify the predicate `TypeURI` suffix from `/scai/attribute-report` to `/scai` per the latest [predicate naming convention].

### New in v0.2

-   Change the `target` and `evidence` field type of a SCAI Attribute
    Assertion to [ResourceDescriptor].
-   Simplify the `producer` field: remove the `type` sub-field, and change
    the field type to [ResourceDescriptor].

[ResourceDescriptor]: ../v1/resource_descriptor.md
[SCAI specification]: https://arxiv.org/pdf/2210.05813.pdf
[TypeURI]: ../v1/field_types.md#typeuri
[attestation Bundle]: ../v1/bundle.md
[parsing rules]: ../v1#parsing-rules
[predicate naming convention]: ../../docs/new_predicate_guidelines.md#vetting-process
