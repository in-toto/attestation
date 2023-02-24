# Examples

> **TODO: remove these, we have sufficient examples in `spec/predicates`.

## SLSA Provenance v0.1

**TODO: Update to v1.0 when draft is released.**

A [SLSA Provenance]-type attestation describing how the
[curl 7.72.0 source tarballs](https://curl.se/download.html) were built,
pretending they were built on
[GitHub Actions](https://github.com/features/actions).

```json
{
  "payloadType": "application/vnd.in-toto+json",
  "payload": "ewogICJzdWJqZWN0IjogWwogICAg...",
  "signatures": [{"sig": "MeQyap6MyFyc9Y..."}]
}
```

where `payload` base64-decodes as the following [Statement]:

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [
    { "name": "curl-7.72.0.tar.bz2",
      "digest": { "sha256": "ad91970864102a59765e20ce16216efc9d6ad381471f7accceceab7d905703ef" }},
    { "name": "curl-7.72.0.tar.gz",
      "digest": { "sha256": "d4d5899a3868fbb6ae1856c3e55a32ce35913de3956d1973caccd37bd0174fa2" }},
    { "name": "curl-7.72.0.tar.xz",
      "digest": { "sha256": "0ded0808c4d85f2ee0db86980ae610cc9d165e9ca9da466196cc73c346513713" }},
    { "name": "curl-7.72.0.zip",
      "digest": { "sha256": "e363cc5b4e500bfc727106434a2578b38440aa18e105d57576f3d8f2abebf888" }}
  ],
  "predicateType": "https://slsa.dev/provenance/v0.1",
  "predicate": {
    "builder": { "id": "https://github.com/Attestations/GitHubHostedActions@v1" },
    "recipe": {
      "type": "https://github.com/Attestations/GitHubActionsWorkflow@v1",
      "definedInMaterial": 0,
      "entryPoint": "build.yaml:maketgz"
    },
    "metadata": {
      "buildStartedOn": "2020-08-19T08:38:00Z"
    },
    "materials": [
      {
        "uri": "git+https://github.com/curl/curl-docker@master",
        "digest": { "sha1": "d6525c840a62b398424a78d792f457477135d0cf" }
      }, {
        "uri": "github_hosted_vm:ubuntu-18.04:20210123.1"
      }, {
        "uri": "git+https://github.com/actions/checkout@v2",
        "digest": {"sha1": "5a4ac9002d0be2fb38bd78e4b4dbde5606d7042f"}
      }, {
        "uri": "git+https://github.com/actions/upload-artifact@v2",
        "digest": { "sha1": "e448a9b857ee2131e752b06002bf0e093c65e571" }
      }, {
        "uri": "pkg:deb/debian/stunnel4@5.50-3?arch=amd64",
        "digest": { "sha256": "e1731ae217fcbc64d4c00d707dcead45c828c5f762bcf8cc56d87de511e096fa" }
      }, {
        "uri": "pkg:deb/debian/python-impacket@0.9.15-5?arch=all",
        "digest": { "sha256": "71fa2e67376c8bc03429e154628ddd7b196ccf9e79dec7319f9c3a312fd76469" }
      }, {
        "uri": "pkg:deb/debian/libzstd-dev@1.3.8+dfsg-3?arch=amd64",
        "digest": { "sha256": "91442b0ae04afc25ab96426761bbdf04b0e3eb286fdfbddb1e704444cb12a625" }
      }, {
        "uri": "pkg:deb/debian/libbrotli-dev@1.0.7-2+deb10u1?arch=amd64",
        "digest": { "sha256": "05b6e467173c451b6211945de47ac0eda2a3dccb3cc7203e800c633f74de8b4f" }
      }
    ]
  }
}
```

### Custom-type examples

In many cases, custom-type predicates would be a more natural fit, as shown
below. Such custom attestations are not yet supported by in-toto because the
layout format has no way to reference such attestations. Still, we show the
examples to explain the benefits for the new link format.

The initial step is often to write code. This has no materials and no real
command. The existing [Link] schema has little benefit. Instead, a custom
`predicateType` would avoid all of the meaningless boilerplate fields.

(Only the base64-decoded `payload` [Statement] is shown, since the outer
[Envelope] looks the same in all cases.)

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [
    { "digest": { "sha1":  "859b387b985ea0f414e4e8099c9f874acb217b94" } }
  ],
  "predicateType": "https://example.com/CodeReview/v1",
  "predicate": {
    "repo": {
      "type": "git",
      "uri": "https://github.com/example/my-project",
      "branch": "main"
    },
    "author": "mailto:alice@example.com",
    "reviewers": ["mailto:bob@example.com"]
  }
}
```

Test results are also an awkward fit for the Link schema, since the subject is
really the materials, not the products. Again, a custom `predicateType` is a
better fit:

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [
    { "digest": { "sha1": "859b387b985ea0f414e4e8099c9f874acb217b94" } }
  ],
  "predicateType": "https://example.com/TestResult/v1",
  "predicate": {
    "passed": true
  }
}
```

[Envelope]: spec/README.md#envelope
[Link]: spec/predicates/link.md
[SLSA Provenance]: https://slsa.dev/provenance
[Statement]: spec/README.md#statement
