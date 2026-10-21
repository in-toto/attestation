# Conformance vectors

A predicate specification decides what a conforming verifier accepts. Nothing in this
repository lets two verifiers of the same predicate find out whether they agree, and
nothing lets a vetting review run anything: `docs/validation.md` is pseudocode, and
no predicate ships an executable case.

This directory is that layer, in one shape for every predicate rather than one shape
per predicate.

## Layout

```text
conformance/
  <predicate-directory-name>/
    manifest.json
    statements/
      <vector-id>.json
```

`<predicate-directory-name>` matches the predicate's document under
`spec/predicates/`, so `spec/predicates/svr.md` is served by `conformance/svr/`.

## What a manifest declares

One entry per vector, and each entry states the verdict a conforming verifier must
reach. That is the whole interface: a verifier reads the manifest, runs every
statement it names, and compares its own verdict to the declared one.

| member | meaning |
| --- | --- |
| `predicateType` | the type URI these vectors are for |
| `specPath` | the document in this repository that decides them |
| `comparisonSurface` | which fields of a verdict are compared, so two verifiers are not held to each other's private detail |
| `vectors[].id` | a stable identifier for the case |
| `vectors[].kind` | `accept`, `reject`, or `indeterminate` |
| `vectors[].file` | the statement, relative to the manifest |
| `vectors[].expected` | the verdict, and for a refusal the codes the refusal must carry |
| `vectors[].conditions` | the specification rules the vector exercises, so a failure names a sentence rather than a file |

A refusal declares its codes because a verifier that refuses the right statement for
the wrong reason is not interoperable with one that refuses it for the right one, and
verdict-only comparison scores that as agreement.

`indeterminate` is a real third kind rather than a gap. A vector whose verdict the
specification genuinely leaves open belongs in the corpus, labelled, with the
readings a conforming verifier could take; dropping it would hide the one case where
two verifiers are both correct and disagree.

## What this directory is not

It is not a second normative surface. Every vector is decided by the document its
manifest names, and where a vector and that document disagree the document governs
and the vector is the defect.

## Adding a predicate

Add a directory, a manifest and the statements. No tool changes: a verifier is
invoked the same way for every predicate, which is what makes the cost of the tenth
predicate the same as the cost of the second.

## What is here now

`adversarial-execution-evidence/` holds four vectors, two accepts and two refusals
carrying different codes, as a worked example of the layout under review in
[#570](https://github.com/in-toto/attestation/pull/570). It is deliberately small:
the layout is the thing to settle first, because a layout adopted after two
predicates have shipped their own is a migration.

See [Verifying an attestation against conformance vectors](../docs/verifying.md) for
how to run one.
