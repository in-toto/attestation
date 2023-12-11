#  Wrapper class for in-toto attestation Statement protos.

import in_toto_attestation.v1.statement_pb2 as spb
from in_toto_attestation.v1.resource_descriptor import ResourceDescriptor

STATEMENT_TYPE_URI = "https://in-toto.io/Statement/v1"


class Statement:
    def __init__(self, subjects: list, predicate_type: str, predicate: dict) -> None:
        self.pb = spb.Statement()  # type: ignore[attr-defined]
        self.pb.type = STATEMENT_TYPE_URI
        self.pb.subject.extend(subjects)
        self.pb.predicate_type = predicate_type
        self.pb.predicate.update(predicate)

    @staticmethod
    def copy_from_pb(proto: spb.Statement) -> "Statement":  # type: ignore[name-defined]
        stmt = Statement([], "", {})
        stmt.pb.CopyFrom(proto)
        return stmt

    def validate(self) -> None:
        if self.pb.type != STATEMENT_TYPE_URI:
            raise ValueError("Wrong statement type")

        if len(self.pb.subject) == 0:
            raise ValueError("At least one subject required")

        # check all resource descriptors in the subject
        subject = self.pb.subject
        for i, rdpb in enumerate(subject):
            rd = ResourceDescriptor.copy_from_pb(rdpb)
            rd.validate()

            # v1 statements require the digest to be set in the subject
            if len(rd.pb.digest) == 0:
                # return index in the subjects list in case of failure:
                # can't assume any other fields in subject are set
                raise ValueError("At least one digest required (subject {0})".format(i))

        if self.pb.predicate_type == "":
            raise ValueError("Predicate type required")

        if len(self.pb.predicate) == 0:
            raise ValueError("Predicate object required")
