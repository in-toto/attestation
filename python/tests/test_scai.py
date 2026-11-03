"""
Tests for in-toto attestation SCAI predicate protos.
"""

import unittest

import google.protobuf.json_format as pb_json

import in_toto_attestation.predicates.scai.v0.scai_pb2 as scaipb
from in_toto_attestation.predicates.scai.v0.scai import AttributeReport


class TestAttributeReport(unittest.TestCase):
    def test_copy_from_pb(self):
        report_json = '{"attributes":[{"attribute":"theAttr","target":{"name":"theTarget"}}],"producer":{"uri":"https://example.com/producer"}}'
        got_pb = pb_json.Parse(report_json, scaipb.AttributeReport())

        got = AttributeReport.copy_from_pb(got_pb)

        self.assertEqual(got.pb, got_pb, "Protos do not match")
        got.validate()

    def test_copy_from_pb_bad_attribute(self):
        bad_report = '{"attributes":[{"target":{"name":"theTarget"}}]}'
        got_pb = pb_json.Parse(bad_report, scaipb.AttributeReport())

        got = AttributeReport.copy_from_pb(got_pb)

        with self.assertRaises(
            ValueError, msg="Error: created malformed AttributeReport (no attribute)"
        ):
            got.validate()


if __name__ == "__main__":
    unittest.main()
