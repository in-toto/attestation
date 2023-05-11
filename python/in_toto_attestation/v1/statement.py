#  Wrapper class for in-toto attestation Statement protos.

import in_toto_attestation.v1.statement_pb2 as spb
from in_toto_attestation.v1.resource_descriptor import ResourceDescriptor

STATEMENT_TYPE_URI = 'https://in-toto.io/Statement/v1'

class Statement:
    def __init__(self, s=None):
        self.pb = spb.Statement()

        if s:
            self.pb.CopyFrom(s)

    def validate(self):    
        if not self.pb.type or self.pb.type != STATEMENT_TYPE_URI:
            raise ValueError('Wrong statement type')

        if not self.pb.subject or len(self.pb.subject) == 0:
            raise ValueError('At least one subject required')

        # check all resource descriptors in the subject
        subject = self.pb.subject
        for rdpb in subject:
            rd = ResourceDescriptor(rdpb)
            rd.validate()

            # v1 statements require the digest to be set in the subject
            if len(rd.pb.digest) == 0:
                raise ValueError('At least one digest required')

        if not self.pb.predicateType:
            raise ValueError('Predicate type required')

        if not self.pb.predicate:
            raise ValueError('Predicate object required')
