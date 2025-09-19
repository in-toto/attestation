# Predicate type: Runtime Trace

Type URI: https://in-toto.io/attestation/runtime-trace/v0.1

Version: 0.1.0

Authors: Parth Patel (@pxp928), Shripad Nadgowda (@nadgowdas),
  Aditya Sirish A Yelgundhalli (@adityasaky)

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

## Data definition

```json
rt-trace-predicate = (
    predicateType-label => "https://in-toto.io/attestation/runtime-trace/v0.1",
    predicate-label => rt-trace-predicate-map
)
rt-trace-predicate-map = {
        rt-trace-monitor-label => rt-trace-monitor-map,
        rt-trace-monitoredProcess-label => rt-trace-monitoredProcess-map,
        rt-trace-monitorLog-label => rt-trace-monitorLog-map,
        ? rt-trace-metadata-label => rt-trace-metadata-map
}
rt-trace-monitor-label          = JC<"monitor",          0>
rt-trace-monitoredProcess-label = JC<"monitoredProcess", 1>
rt-trace-monitorLog-label       = JC<"monitorLog",       2>
rt-trace-metadata-label         = JC<"metadata",         3>

rt-trace-monitor-map = {
  rt-trace-monitor-type-label => uri-type,
  ? rt-trace-monitor-configSource-label => ResourceDescriptor,
  ? rt-trace-monitor-tracePolicy-label => object
}
rt-trace-monitor-type-label         = JC<"type",         0>
rt-trace-monitor-configSource-label = JC<"configSource", 1>
rt-trace-monitor-tracePolicy-label  = JC<"tracePolicy",  2>

rt-trace-monitoredProcess-map {
  rt-trace-monitoredProcess-hostID-label => uri-type,
  rt-trace-monitoredProcess-type-label => uri-type,
  rt-trace-monitoredProcess-event-label => text
}
rt-trace-monitoredProcess-hostID-label = JC<"hostID", 0>
rt-trace-monitoredProcess-type-label   = JC<"type",   1>
rt-trace-monitoredProcess-event-label  = JC<"event",  2>

rt-trace-monitorLog-map = nonempty<{
  ? rt-trace-monitorLog-process-label => [ * object ],
  ? rt-trace-monitorLog-network-label => [ * object ],
  ? rt-trace-monitorLog-fileAccess-label => [ * ResourceDescriptor ]
}>
rt-trace-monitorLog-process-label    = JC<"process",    0>
rt-trace-monitorLog-network-label    = JC<"network",    1>
rt-trace-monitorLog-fileAccess-label = JC<"fileAccess", 2>

rt-trace-metadata-map = {
  ? rt-trace-buildStartedOn-label => Timestamp,
  ? rt-trace-buildFinishedOn-label => Timestamp
}
rt-trace-buildStartedOn-label  = JC<"buildStartedOn",  0>
rt-trace-buildFinishedOn-label = JC<"buildFinishedOn", 1>
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

URI indicating the process host’s identity.

`monitoredProcess.type` _string (TypeURI)_, _required_

URI indicating the type of process performed. Ex: when monitoring a build, this
field records the build type.

`monitoredProcess.event` _string_, _required_

String identifying the specific job or task associated with the attestation.

`monitorLog` _object_, _required_

Record of events that were traced by the monitor. At least one of `process`,
`network`, and `fileAccess` must be present.

`monitorLog.process` _list of objects_, _optional_

Record of processes observed by monitor. The exact format of this field is
currently dependent on the monitor, and can be determined using `monitor.type`.
In future, after consulting different monitors, this field may have a consistent
schema.

`monitorLog.network` _list of objects_, _optional_

Record of network activity observed by monitor. The exact format of this field
is currently dependent on the monitor, and can be determined using
`monitor.type`. In future, after consulting different monitors, this field may
have a consistent schema.

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
  "predicateType": "https://in-toto.io/attestation/runtime-trace/v0.1",
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
      ]
    },
    "metadata": {
      "buildStartedOn": "2022-10-13T22:55:50Z",
      "buildFinishedOn": "2022-10-13T22:56:17Z"
    }
  }
}
```
