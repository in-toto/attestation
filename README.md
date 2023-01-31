# in-toto Attestation Framework

<!-- msm: This is real clunky -->
The in-toto Attestation Framework is a framework for generating verifiable
claims about any aspect of how a piece of software is produced, allowing
anyone to validate the origins of software and establish trust in its supply
chain.

## Learning about in-toto attestations

Visit [https://in-toto.io](https://in-toto.io) to learn about the larger
in-toto project.

For a deeper dive, we recommend reading through our
[documentation](https://github.com/in-toto/attestation/tree/main/docs) to
learn more about the goals of the in-toto Attestation Framework.

## Working with in-toto attestations

The core of the in-toto Attestation Framework is the
[specification](https://github.com/in-toto/attestation/tree/main/spec)
that defines the format for in-toto attestations and the metadata they
contain.

Before opening a request for a new metadata format, please review the
following guidelines.

### Want to share information about how your software is produced?

As a first step, peruse the existing attestation [predicate types]. These are
metadata formats vetted by our maintainers to cover a number of common use
cases.

### Is your use case not covered by existing predicate types?

Take a look at the open [issues] or [pull requests] to see if your usage has
already been reported. We can help with use cases, thinking through options,
and questions about existing predicates. Feel free to comment on an existing
issue or PR.

### Need to start a new discussion?

If you still can't find what you're looking for, open a new issue. Make sure
you cover the following questions, so that we can better understand your
request:

-   What’s you use case?
-   Why don’t existing predicates cover this use case?
-   What might a new predicate type for your use case look like?
(concrete examples in JSON or CUE preferred)
-   What policy questions do you want to be able to answer with the predicate?

### Want to propose a new predicate type?

We love to see our list of vetted [predicates] grow. If you would like your
new metadata format to go through our vetting process, submit a PR following
the [ITE-9] formatting guidelines.

**Note**: Your predicate type is yours. This means that in-toto Attestation
Framework maintainers can provide feedback, but will not write the
specification for you.

## Governance

The in-toto Attestation Framework is part of the [in-toto] project under the
[CNCF].

Use `@in-toto/attestation-maintainers` to tag the maintainers on GitHub.

## Disclaimer

The in-toto Attestation Framework is still under development. We are in the
process of developing tooling to enable better integration and adoption of
the framework. In the meantime, please visit any of the language-specific
[in-toto implementations] to become familiar with current tooling options.

[CNCF]: https://www.cncf.io/projects/in-toto/
[ITE-9]: https://github.com/in-toto/ITE/tree/master/ITE/9#document-format
[in-toto]: https://in-toto.io
[in-toto implementations]: https://github.com/in-toto
[issues]: https://github.com/in-toto/attestation/issues?q=is%3Aopen+is%3Aissue
[predicate types]: https://github.com/in-toto/attestation/tree/main/spec/predicates
[pull requests]: https://github.com/in-toto/attestation/pulls?q=is%3Aopen+is%3Apr
