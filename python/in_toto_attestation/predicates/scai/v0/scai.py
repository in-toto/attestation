# Wrapper class for in-toto attestation SCAI predicate protos.

import in_toto_attestation.predicates.scai.v0.scai_pb2 as scaipb
from in_toto_attestation.v1.resource_descriptor import ResourceDescriptor

SCAI_PREDICATE_TYPE = "https://in-toto.io/attestation/scai/attribute-report/"
SCAI_PREDICATE_VERSION = "v0.2"


class AttributeAssertion:
    def __init__(self, attribute, target=None, conditions=None, evidence=None) -> None:
        self.pb = scaipb.AttributeAssertion()  # type: ignore[attr-defined]
        self.pb.attribute = attribute
        if target:
            self.pb.target.CopyFrom(target)
        if conditions:
            self.pb.conditions.update(conditions)
        if evidence:
            self.pb.evidence.CopyFrom(evidence)

    @staticmethod
    def copy_from_pb(proto: type[scaipb.AttributeAssertion]) -> "AttributeAssertion":  # type: ignore[name-defined]
        assertion = AttributeAssertion("tmp-attr")
        assertion.pb.CopyFrom(proto)
        return assertion

    def validate(self) -> None:
        if len(self.pb.attribute) == 0:
            raise ValueError("The attribute field is required")

        # check any resource descriptors are valid
        if self.pb.target.ByteSize() > 0:
            rd = ResourceDescriptor.copy_from_pb(self.pb.target)
            rd.validate()

        if self.pb.evidence.ByteSize() > 0:
            rd = ResourceDescriptor.copy_from_pb(self.pb.evidence)
            rd.validate()


class AttributeReport:
    def __init__(self, attributes: list, producer=None) -> None:
        self.pb = scaipb.AttributeReport()  # type: ignore[attr-defined]
        self.pb.attributes.extend(attributes)
        if producer:
            self.pb.producer.CopyFrom(producer)

    @staticmethod
    def copy_from_pb(proto: type[scaipb.AttributeReport]) -> "AttributeReport":  # type: ignore[name-defined]
        report = AttributeReport([None])
        report.pb.CopyFrom(proto)
        return report

    def validate(self) -> None:
        if len(self.pb.attributes) == 0:
            raise ValueError("The attributes field is required")

        # check any attribute assertions are valid
        attributes = self.pb.attributes
        for i, attrpb in enumerate(attributes):
            assertion = AttributeAssertion.copy_from_pb(attrpb)

            try:
                assertion.validate()
            except ValueError as e:
                # return index in the attributes list in case of failure:
                # can't assume any other fields in attribute assertion are set
                raise ValueError(
                    "Invalid attributes field (index: {0})".format(i)
                ) from e

        if self.pb.producer.ByteSize() > 0:
            rd = ResourceDescriptor.copy_from_pb(self.pb.producer)
            rd.validate()
