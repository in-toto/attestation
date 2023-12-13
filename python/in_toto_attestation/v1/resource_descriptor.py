#  Wrapper class for in-toto attestation ResourceDescriptor protos.

from __future__ import annotations

import json

import google.protobuf.json_format as pb_json
from google.protobuf.struct_pb2 import Value

import in_toto_attestation.v1.resource_descriptor_pb2 as rdpb


class ResourceDescriptor:
    def __init__(
        self,
        name: str = "",
        uri: str = "",
        digest: dict | None = None,
        content: bytes = bytes(),
        download_location: str = "",
        media_type: str = "",
        annotations: dict | None = None,
    ) -> None:
        self.pb = rdpb.ResourceDescriptor()  # type: ignore[attr-defined]
        self.pb.name = name
        self.pb.uri = uri
        if digest:
            self.pb.digest.update(digest)
        self.pb.content = content
        self.pb.download_location = download_location
        self.pb.media_type = media_type
        if annotations:
            self.pb.annotations.update(annotations)

    @staticmethod
    def copy_from_pb(proto: rdpb.ResourceDescriptor) -> "ResourceDescriptor":  # type: ignore[name-defined]
        rd = ResourceDescriptor()
        rd.pb.CopyFrom(proto)
        return rd

    def validate(self) -> None:
        # at least one of name, URI or digest are required
        if self.pb.name == "" and self.pb.uri == "" and len(self.pb.digest) == 0:
            raise ValueError("At least one of name, URI, or digest need to be set")
