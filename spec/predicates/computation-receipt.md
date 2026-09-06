# Predicate type: Computation Receipt

Type URI: https://in-toto.io/attestation/computation-receipt/v0.1

Version: v0.1

Predicate Name: Computation Receipt

Authors: Anomly, Inc. (https://github.com/anomly-labs)

## Purpose

To attest that a numeric computation — an HPC reduction, a training step, a model
inference — produced a specific output, in a form checkable by re-execution. The
predicate binds the output, inputs, model weights, computation identity, and the
arithmetic profile by digest into a canonical manifest whose SHA-256 is the claim's
certificate. Under an order-independent exact arithmetic profile, every honest
re-execution reproduces the output bit-identically on any hardware, so verification
is: re-execute, compare digests.

## Use Cases

- A client verifies on its own (different) hardware that a provider's inference or
  reduction result is exactly correct.
- An auditor re-derives a recorded model output months later from the pinned
  weights, inputs, and arithmetic profile.
- A mixed-hardware fleet proves two nodes computed the same answer to the bit.

## Prerequisites

The in-toto Attestation Framework and the Computation Receipt v0.1 specification
(https://github.com/anomly-labs/computation-receipts — spec, conformance vectors
including required refusals, and reference implementations).

## Model

The subject is the certified output artifact. `subject[0].digest.sha256` MUST equal
the hex digest in the predicate's `manifest.output.digest`. The predicate is a
complete CR receipt object, carried verbatim.

## Schema

