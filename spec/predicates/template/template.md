# Predicate type: <name>

Type URI:

Version:

Predicate Name:

## Purpose

This section must discuss at a high level the reason for creating the new
predicate type. Essentially, it describes the problems that can be solved using
the new predicate.

## Use Cases

This section must discuss at least one concrete use case that will be covered by
this predicate and how the current predicate types (if any) do not solve this.

## Prerequisites

New predicate types may be very specific to a particular technology or
specification, which must be detailed in the prerequisites section. The in-toto
Attestation Framework may be a necessary prerequisite for all predicate types,
if so it SHOULD be explicitly mentioned to err on the side of caution.

## Model

The in-toto Attestation Framework was designed with software supply chain
specific metadata in mind. However, supply chains are vast and composed of many
different components and phases. Therefore, it is important for the predicate
type to declare at a high level the steps and functionaries (in the generic) it
is relevant to.

## Schema

The schema of the predicate type is the core part of the document. It defines
the fields included in an instantiation of the predicate, as well as rules to
parse them.

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "<Type URI>",
  "predicate": { ... }
}
```

### Parsing Rules

Parsing rules are important to define explicitly to ensure implementations can
handle example attestations correctly. For example, this section can discuss how
the predicate type is versioned, how non-specified fields must be handled, and
so on. Attestations definitions MUST use this section to define how the parsing
rules differ from the framework’s
[standard parsing rules](/spec/v1/README.md#parsing-rules).

### Fields

This subsection defines the fields that make up the new predicate type. In
addition to the predicate's fields, this subsection can be used to specify
additional constraints at the statement layer.

The fields must be listed with their names, the type of information they can
hold, whether they are required or not, and a textual description of their
expected contents. If the legal values for a field are known, or belong in a
range, this must be specified in the description.

## Example

A new predicate type MUST include one or more examples. It is recommended that
this example corresponds to the Purpose defined at the start of the document.

The example included MUST be an entire in-toto attestation statement, and not
just the predicate type. This will allow adopters and implementers to work with
examples more easily.

## Changelog and Migrations

For predicate definitions, it is important to be able to clearly map the changes
made to the schema of the predicate to its different versions. These changes and
the steps to migrate from one to another must be detailed where appropriate.
This section MAY be skipped for the initial version of a predicate type.
