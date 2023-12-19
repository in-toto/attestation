#  Wrapper class for in-toto attestation Statement protos.

from in_toto_attestation._internal.in_toto_attestation.v1 import Statement as _Statement
from in_toto_attestation.v1.resource_descriptor import ResourceDescriptor

# import in_toto_attestation.v1.statement_pb2 as spb
# from in_toto_attestation.v1.resource_descriptor import ResourceDescriptor

STATEMENT_TYPE_URI = "https://in-toto.io/Statement/v1"


class Statement(_Statement):
    def validate(self) -> None:
        if self.type != STATEMENT_TYPE_URI:
            raise ValueError("Wrong statement type")

        if len(self.subject) == 0:
            raise ValueError("At least one subject required")

        # check all resource descriptors in the subject
        subject = self.subject
        for i, rdpb in enumerate(subject):
            rd = ResourceDescriptor().from_dict(rdpb.to_dict())
            rd.validate()

            # v1 statements require the digest to be set in the subject
            if len(rd.digest) == 0:
                # return index in the subjects list in case of failure:
                # can't assume any other fields in subject are set
                raise ValueError(f"At least one digest required (subject {i})")

        if self.predicate_type == "":
            raise ValueError("Predicate type required")

        if len(self.predicate) == 0:
            raise ValueError("Predicate object required")
