#  Wrapper class for in-toto attestation ResourceDescriptor protos.

import in_toto_attestation.v1.resource_descriptor_pb2 as rdpb

class ResourceDescriptor:
    def __init__(self, rd=None):
        self.pb = rdpb.ResourceDescriptor()

        if rd:
            self.pb.CopyFrom(rd)

    def validate(self):
        # at least one of name, URI or digest are required
        if (not self.pb.name and not self.pb.uri and not self.pb.digest) or len(self.pb.digest) == 0:
            raise ValueError("At least one of name, URI, or digest are required")
