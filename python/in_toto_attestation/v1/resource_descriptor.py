#  Wrapper class for in-toto attestation ResourceDescriptor protos.

import in_toto_attestation.v1.resource_descriptor_pb2 as rdpb
from google.protobuf.struct_pb2 import Value
import google.protobuf.json_format as pb_json
import json

class ResourceDescriptor:
    def __init__(self, name: str='', uri: str='', digest: dict=None, content: bytes=bytes(), download_location: str='', media_type: str='', annotations: dict=None) -> None:
        self.pb = rdpb.ResourceDescriptor()
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
    def copy_from_pb(proto: type[rdpb.ResourceDescriptor]) -> 'ResourceDescriptor':
        rd = ResourceDescriptor()
        rd.pb.CopyFrom(proto)
        return rd

    def validate(self) -> None:
        # at least one of name, URI or digest are required
        if self.pb.name == '' and self.pb.uri == '' and len(self.pb.digest) == 0:
            raise ValueError("At least one of name, URI, or digest need to be set")
