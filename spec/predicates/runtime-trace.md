# Predicate type: Runtime Trace

Type URI: https://in-toto.io/attestation/runtime-trace/v0.2

Version: 0.2.0

Authors: Parth Patel (@pxp928), Shripad Nadgowda (@nadgowdas),
  Aditya Sirish A Yelgundhalli (@adityasaky), Hiroki Suezawa (@rung)

## Purpose

This predicate is used to describe system events that were part of some software
supply chain step, for example, the build process of an artifact. Supply chain
operations today are complex and opaque with little to no insights into
underlying system operations. However, there are emerging introspection
techniques that now allow us to monitor and capture system event logs during the
monitored operation.

## Use Cases

Generally, this predicate can be used to express a runtime trace of any
operation. The predicate specification does not mandate any specific monitoring
tool or technology. The schema can be used to express runtime traces of anything
that can be traced or monitored, from an operation spawned by a user using a CLI
command to tasks performed in CI systems. The primary use case driving the
development of this predicate is its applicability to runtime traces of build
processes.

This attestation can be used in conjunction with provenance attestations to
leverage observability to attest to various SLSA requirements. The runtime trace
can prove the build was invoked via a script, that the build was executed in a
hermetic environment with no network access, and so on.

## Prerequisites

Understanding of monitoring tools and the larger in-toto attestation framework.

## Model

`monitor` identifies the specific instance of a tool that’s observing the target
process and details the configurations of the monitor, while `monitoredProcess`
contains information that when put together uniquely identifies the exact
instance of the job being observed. The `monitorLog` field contains the actual
runtime trace information.

## Schema

```json
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [{ ... }],
    "predicateType": "https://in-toto.io/attestation/runtime-trace/v0.2",
    "predicate": {
        "monitor": {
            "type": "<TypeURI>",
            "configSource": "<ResourceDescriptor>",
            "tracePolicy": { /* object */ }
        },
        "monitoredProcess": {
            "hostID": "<URI>",
            "type": "<URI>",
            "event": "<STRING>"
        },
        "monitorLog": {
            "process": [
                { /* object */ }
            ],
            "network": [
                {
                    "ip": "<STRING>",
                    "hostname": "<STRING>",
                    "port": /* integer */,
                    "protocol": "<STRING>",
                    "invokingProcess": { /* object */ }
                }
            ],
            "fileAccess": ["<ResourceDescriptor>", ...]
        },
        "metadata": {
            "buildStartedOn": "<TIMESTAMP>",
            "buildFinishedOn": "<TIMESTAMP>"
        }
    }
}
```

### Parsing Rules

This predicate follows the in-toto attestation parsing rules. Summary:

-   Consumers MUST ignore unrecognized fields.
-   The `predicateType` URI includes the major version number and will always
    change whenever there is a backwards incompatible change.
-   Minor version changes are always backwards compatible and “monotonic.” Such
    changes do not update the `predicateType`.
-   Producers MAY add extension fields using field names that are unlikely to
    collide with names used by other producers. Field names SHOULD avoid using
    characters like `.` and `$`.
-   Fields marked _optional_ MAY be unset or null, and should be treated
    equivalently. Both are equivalent to empty for object or array values.

### Fields

`monitor` _object_, _required_

Identifies the specific monitor instance used to trace the runtime.

`monitor.type` _TypeURI_, _required_

URI indicating the monitor’s type.

`monitor.configSource` ResourceDescriptor, _optional_

Effectively a pointer to the monitor's configuration.

`monitor.tracePolicy` _object_, _optional_

Indicates the trace policy used for the monitoring event that generated the
runtime trace. The trace policy's format is dependent on the monitor used, and
so it must be parsed after identifying the monitor using `monitor.type`.

FIXME: should this be optional?

`monitoredProcess` _object_, _required_

Identifies the process being monitored.

`monitoredProcess.hostID` _string (URI)_, _required_

URI indicating the process host’s identity. Ex: when monitoring a job on a
CI/CD platform, this field records the platform’s identity, such as
`https://github.com`.

`monitoredProcess.type` _string (TypeURI)_, _required_

URI indicating the type of process performed. Ex: when monitoring a build, this
field records the build type. For builds on CI/CD platforms, producers SHOULD
reuse the platform’s published SLSA `buildType` URI, such as
`https://slsa-framework.github.io/github-actions-buildtypes/workflow/v1` for
GitHub Actions.

`monitoredProcess.event` _string_, _required_

String identifying the specific job or task associated with the attestation.
The value SHOULD uniquely identify the monitored job or task, so that
consumers can tie the trace back to the exact run and cross-reference it with
other attestations for the same run, such as build provenance. Ex:
`https://github.com/<org>/<repo>/actions/runs/<run_id>/attempts/<n>` for a
GitHub Actions job.

`monitorLog` _object_, _required_

Record of events that were traced by the monitor. At least one of `process`,
`network`, and `fileAccess` must be present.

`monitorLog.process` _list of objects_, _optional_

Record of processes observed by monitor. The exact format of this field is
currently dependent on the monitor, and can be determined using `monitor.type`.
In future, after consulting different monitors, this field may have a consistent
schema.

`monitorLog.network` _list of objects_, _optional_

Record of network activity observed by monitor. Entries SHOULD use the
following fields so that consumers can process network activity consistently
across monitors. Producers MAY include additional, monitor-specific fields in
each entry.

`monitorLog.network[*].ip` _string_, _optional_

IP address of the remote endpoint. Both IPv4 and IPv6 addresses use this
field. The canonical textual form is recommended (for IPv6, [RFC 5952]).

`monitorLog.network[*].hostname` _string_, _optional_

Observed or resolved hostname of the remote endpoint. An entry SHOULD carry at
least one of `ip` or `hostname`, and MAY carry just one of them when the
association between the two is unknown.

`monitorLog.network[*].port` _integer_, _optional_

Port number of the remote endpoint, between 1 and 65535.

`monitorLog.network[*].protocol` _string_, _optional_

Transport protocol, lowercase. Ex: `tcp`, `udp`. If `port` is present,
`protocol` SHOULD also be present.

`monitorLog.network[*].invokingProcess` _object_, _optional_

The process that generated the network activity. The exact format of this
field is dependent on the monitor, and can be determined using `monitor.type`.
Ex: `{ "pid": 1234, "comm": "curl" }`.

`monitorLog.fileAccess` _list of ResourceDescriptor objects_, _optional_

Record of files accessed during the monitored process. A complete list of
_materials_ can be derived from this information. Each entry in this list is
expected to record the path of the file and one or more digests of the file, but
as each entry is an instance of _ResourceDescriptor_, other supported
information can also be captured. This field is a list rather than a key-value
map because a single file may be used multiple times during the build process.
Further, some files that are accessed may _change_ during the build process, and
so, different entries for the same file may have different digests.

Note: While this predicate can be used to log file accesses, the actual
technique used to capture the file access event has some implications. If a
synchronous monitor, for example one that uses `ptrace` to trace the file access
system calls, is used, then the build process can be paused while the file's
digest is calculated and stored. However, asynchronous monitors such as those
using eBPF cannot pause the build process before the file is actually used.
Therefore, they cannot make as strong guarantees about the digests of the files
accessed. Verifiers using runtime trace attestations for file accesses must be
careful about what guarantees they are actually getting based on how the build
process was monitored.

`metadata` _object_, _optional_

Other properties of the monitoring event.

`metadata.buildStartedOn` _Timestamp_, _optional_

`metadata.buildFinishedOn` _Timestamp_, _optional_

## Example

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "predicateType": "https://in-toto.io/attestation/runtime-trace/v0.2",
  "subject": [
    {
      "name": "ttl.sh/testin123",
      "digest": {
        "sha256": "def2bdf8cee687d5889d51923d7907c441f1a61958f1e5dfb07f53041c83745f"
      }
    }
  ],
  "predicate": {
    "monitor": {
      "type": "https://github.com/cilium/tetragon/v0.8.4",
      "configSource": {},
      "tracePolicy": {
        "policies": [
          {
            "Name": "connect",
            "Config": ""
          },
          {
            "Name": "sys-read-follow-prefix",
            "Config": ""
          }
        ]
      }
    },
    "monitoredProcess": {
      "hostID": "https://tekton.dev/chains/v2",
      "type": "https://tekton.dev/attestations/chains@v2",
      "event": "run-image-pipelinerun-build-trusted"
    },
    "monitorLog": {
      "process": [
        {
          "eventType": "process",
          "processBinary": "/ko-app/entrypoint",
          "arguments": [
            "init /ko-app/entrypoint /tekton/bin/entrypoint step-prepare step-create step-results"
          ],
          "privileged": null
        },
        {
          "eventType": "process",
          "processBinary": "/bin/sh",
          "arguments": [
            "-c \"scriptfile=\"/tekton/scripts/script-0-vjklb\"\ntouch ${scriptfile} && chmod +x ${scriptfile}\ncat > ${s\" \"riptfile} << '_EOF_'\nIyEvdXNyL2Jpbi9lbnYgYmFzaApzZXQgLWUKCmlmIFtbICJ0cnVlIiA9PSAidHJ1ZSIgXV07IHRoZW\" KICBlY2hvICI+IFNldHRpbmcgcGVybWlzc2lvbnMgb24gJy93b3Jrc3BhY2UvY2FjaGUnLi4uIgogIGNob3duIC1SICIxMDAwOj wMDAiICIvd29ya3NwYWNlL2NhY2hlIgpmaQoKZm9yIHBhdGggaW4gIi90ZWt0b24vaG9tZSIgIi9sYXllcnMiICIvd29ya3NwYW lL3NvdXJjZSI7IGRvCiAgZWNobyAiPiBTZXR0aW5nIHBlcm1pc3Npb25zIG9uICckcGF0aCcuLi4iCiAgY2hvd24gLVIgIjEwMD 6MTAwMCIgIiRwYXRoIgoKICBpZiBbWyAiJHBhdGgiID09ICIvd29ya3NwYWNlL3NvdXJjZSIgXV07IHRoZW4KICAgICAgY2htb2 gNzc1ICIvd29ya3NwYWNlL3NvdXJjZSIKICBmaQpkb25lCgplY2hvICI+IFBhcnNpbmcgYWRkaXRpb25hbCBjb25maWd1cmF0aW uLi4uIgpwYXJzaW5nX2ZsYWc9IiIKZW52cz0oKQpmb3IgYXJnIGluICIkQCI7IGRvCiAgICBpZiBbWyAiJGFyZyIgPT0gIi0tZW 2LXZhcnMiIF1dOyB0aGVuCiAgICAgICAgZWNobyAiLT4gUGFyc2luZyBlbnYgdmFyaWFibGVzLi4uIgogICAgICAgIHBhcnNpbm fZmxhZz0iZW52LXZhcnMiCiAgICBlbGlmIFtbICIkcGFyc2luZ19mbGFnIiA9PSAiZW52LXZhcnMiIF1dOyB0aGVuCiAgICAgIC"
          ],
          "privileged": null
        },
        {
          "eventType": "process",
          "processBinary": "/bin/touch",
          "arguments": [
            "/tekton/scripts/script-0-vjklb"
          ],
          "privileged": null
        },
        ...,
        {
          "eventType": "exit",
          "processBinary": "/bin/sh",
          "arguments": [
            "",
            "-c \"scriptfile=\"/tekton/scripts/script-0-vjklb\"\ntouch ${scriptfile} && chmod +x ${scriptfile}\ncat > ${s\" \"riptfile} << '_EOF_'\nIyEvdXNyL2Jpbi9lbnYgYmFzaApzZXQgLWUKCmlmIFtbICJ0cnVlIiA9PSAidHJ1ZSIgXV07IHRoZW\" KICBlY2hvICI+IFNldHRpbmcgcGVybWlzc2lvbnMgb24gJy93b3Jrc3BhY2UvY2FjaGUnLi4uIgogIGNob3duIC1SICIxMDAwOj wMDAiICIvd29ya3NwYWNlL2NhY2hlIgpmaQoKZm9yIHBhdGggaW4gIi90ZWt0b24vaG9tZSIgIi9sYXllcnMiICIvd29ya3NwYW lL3NvdXJjZSI7IGRvCiAgZWNobyAiPiBTZXR0aW5nIHBlcm1pc3Npb25zIG9uICckcGF0aCcuLi4iCiAgY2hvd24gLVIgIjEwMD 6MTAwMCIgIiRwYXRoIgoKICBpZiBbWyAiJHBhdGgiID09ICIvd29ya3NwYWNlL3NvdXJjZSIgXV07IHRoZW4KICAgICAgY2htb2 gNzc1ICIvd29ya3NwYWNlL3NvdXJjZSIKICBmaQpkb25lCgplY2hvICI+IFBhcnNpbmcgYWRkaXRpb25hbCBjb25maWd1cmF0aW uLi4uIgpwYXJzaW5nX2ZsYWc9IiIKZW52cz0oKQpmb3IgYXJnIGluICIkQCI7IGRvCiAgICBpZiBbWyAiJGFyZyIgPT0gIi0tZW 2LXZhcnMiIF1dOyB0aGVuCiAgICAgICAgZWNobyAiLT4gUGFyc2luZyBlbnYgdmFyaWFibGVzLi4uIgogICAgICAgIHBhcnNpbm fZmxhZz0iZW52LXZhcnMiCiAgICBlbGlmIFtbICIkcGFyc2luZ19mbGFnIiA9PSAiZW52LXZhcnMiIF1dOyB0aGVuCiAgICAgIC"
          ],
          "privileged": null
        },
        {
          "eventType": "exit",
          "processBinary": "/tekton/bin/entrypoint",
          "arguments": [
            "",
            "-wait_file /tekton/downward/ready -wait_file_content -post_file /tekton/run/0/out -termination_path /tekton/termination -step_metadata_dir /tekton/run/0/status -results APP_IMAGE_DIGEST,APP_IMAGE_URL -entrypoint /tekton/scripts/script-0-vjklb -- --env-vars"
          ],
          "privileged": null
        }
      ],
      "network": [
        {
          "ip": "140.82.112.3",
          "hostname": "github.com",
          "port": 443,
          "protocol": "tcp"
        },
        {
          "ip": "8.8.8.8",
          "port": 53,
          "protocol": "udp"
        }
      ]
    },
    "metadata": {
      "buildStartedOn": "2022-10-13T22:55:50Z",
      "buildFinishedOn": "2022-10-13T22:56:17Z"
    }
  }
}
```

## Changelog and Migrations

### v0.2

-   Defined a recommended entry shape for `monitorLog.network` entries;
    the format was previously monitor-dependent, so producers SHOULD
    migrate to the recommended field names.
-   Clarified how `monitoredProcess` identifies CI/CD jobs.

### v0.1

Initial version.

[RFC 5952]: https://www.rfc-editor.org/rfc/rfc5952
