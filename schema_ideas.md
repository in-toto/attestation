# Alternative schema design ideas

This doc lists several alternative ideas that are under discussion for the
schema of attestations.

## Terminology and Concepts

*   **Artifact:** blob of data, identified solely by contents
*   **Resource:** uniquely identifiable thing that can be associated with an
    artifact

The following example shows various resources associated with the software
[curl](https://curl.se). Each is a way that a user may retrieve the software,
and each maps to a specific artifact at a point in time.

| Description    | Example Resource IDs                            | Example Artifact ID     |
| -------------- | ----------------------------------------------- | -------------------     |
| Git commit     | https://github.com/curl/curl, branch “master”   | Commit ID “9d954…”      |
|                | https://github.com/curl/curl, tag “curl-7_72_0” |                         |
| Source tarball | https://curl.se/download/curl-7.72.0.tar.gz     | SHA256 “d4d58…”         |
| Debian package | [curl][curl-debian] (latest version)            | SHA256 “aec36e…”        |
|                | [curl 7.72.0-1_amd64][curl-debian]              |                         |
| Docker image   | [docker.io/curlimages/curl:latest][curl-docker] | Digest “sha256:09a374…” |
|                | [docker.io/curlimages/curl:7.72.0][curl-docker] |                         |

[curl-debian]: https://packages.debian.org/search?keywords=curl
[curl-docker]: https://hub.docker.com/r/curlimages/curl

## Use cases

Attestation archetypes:

*   **Provenance:** how an artifact (or set of artifacts) came into being.
    *   Subject: artifact
*   **Artifact Analysis:** metadata about an artifact
*   **Artifact Equivalence:** two artifacts are equivalent
    *   Example: git commit == git tree == tar.gz == zip

**Open Question:** Are there any use cases that don't cleanly fit into one of
the archetypes above?

## Current proposal ([link](README.md))

Code review:

```json
{
  "attestation_type": "https://example.com/CodeReview/v1",
  "subject": { "git_commit_id": "859b387b985ea0f414e4e8099c9f874acb217b94" },
  "details": {
    "timestamp": "2020-04-12T13:50:00Z",
    "repo_type": "git",
    "repo_url": "https://github.com/my-company/my-product",
    "repo_branch": "master",
    "code_reviewed": true
  }
}
```

Provenance:

```json
{
  "attestation_type": "https://example.com/GitHubActionProduct/v1",
  "subject": { "container_image_digest": "sha256:c201c331d6142766c866..." },
  "relations": {
    "top_level_source": [{
      "artifact": { "git_commit_id": "859b387b985ea0f414e4e8099c9f874acb217b94" },
      "git_repo": "https://github.com/example/repo"
    }],
    "dependent_sources": [{
      "artifact": { "git_commit_id": "2f02c094e6a9afe8e889c3f1d3cb66b437797af4" },
      "git_repo": "https://github.com/example/submodule1"
      }, {
      "artifact": { "git_commit_id": "5215a97a7978d8ee0de859ccac1bbfd2475bfe92" },
      "git_repo": "https://github.com/example/submodule2"
    }],
    "tools": [{
      "artifact": { "sha256": "411c1dfb3c8f3bea29da934d61a884baad341af8..." },
      "name": "clang"
      }, {
      "artifact": { "sha256": "9f5068311eb98e6dd9bb554d4b7b9ee126b13693..." },
      "name": "bazel"
    }]
  },
  "details": {
    "workflow_name": "Build",
    "hermetic": true
  }
}
```

Artifact analysis:

```json
{
  "attestation_type": "https://example.com/VulnerabilityScan/v1",
  "subject": { "git_commit_id": "859b387b985ea0f414e4e8099c9f874acb217b94" },
  "details": {
    "timestamp": "2020-04-12T13:55:02Z",
    "vulnerability_counts": {
      "high": 0,
      "medium": 1,
      "low": 17
    }
  }
}
```

## NEW IDEA: Provenance should be generic way to reproduce build

Contains all information so that one can build it hermetically.

*   Environment variables
*   Working directory
*   Entry point
*   Architecture
*   How to get artifacts out
*   All dependencies to be fetched up front
    *   Where to fetch from. Try to minimize number of custom schemes.
        *   HTTPS get
        *   Git checkout
        *   Hg checkout
        *   Container registry
    *   Digest of artifact
    *   Where to place it

Would probably be verbose, but maybe templating (below) can help with that.

Look at Debian buildinfo for inspiration.

## NEW IDEA: Templating

For provenance (at least) allow templates and instancing. That way all common
stuff can be defined by the template without having to be standardized. Tempalte
must be a retrievable URL (though not necessarily publicly accessible) and must
resolve to a Canonical JSON representation of the template. Maybe also a ".sig"
contains a signature.

```json
{
  "attestation_type": "...Provenance...",
  "template": "https://github.com/Actions/Template/v1",
  "params": {
    "workflow": "build"
  }
}
```

Or for Debian:

```javascript
{
  "attestation_type": "...Provenance...",
  "template": "...DebianBuild/v1",
  "params": {
    "source_date_epoch": "2020-08-24T08:26:12Z"  # or just int
  }
}
```

which could get translated into


