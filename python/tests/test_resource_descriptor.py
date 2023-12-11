"""
Tests for in-toto attestation ResourceDescriptor protos.
"""

import unittest
import google.protobuf.json_format as pb_json

import in_toto_attestation.v1.resource_descriptor_pb2 as rdpb
from in_toto_attestation.v1.resource_descriptor import ResourceDescriptor


def create_test_desc():
    desc = ResourceDescriptor(
        "theName",
        "https://example.com",
        {"alg": "abc123"},
        b"bytescontent",
        "https://example.com/test.zip",
        "theMediaType",
        {"keyStr": "val1", "keyNum": 13, "keyObj": {"subKey": "subVal"}},
    )

    return desc


class TestResourceDescriptor(unittest.TestCase):
    def test_create_resource_descriptor(self):
        test_rd = create_test_desc()
        test_rd.validate()

    def test_json_parse_resource_descriptor(self):
        full_rd = '{"name":"theName","uri":"https://example.com","digest":{"alg":"abc123"},"content":"Ynl0ZXNjb250ZW50","downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType","annotations":{"keyNum": 13,"keyStr":"val1","keyObj":{"subKey":"subVal"}}}'
        got_pb = pb_json.Parse(full_rd, rdpb.ResourceDescriptor())
        got = got_pb.SerializeToString(deterministic=True)

        test_rd = create_test_desc()
        want = test_rd.pb.SerializeToString(deterministic=True)

        self.assertEqual(got, want, "Protos do not match")

    def test_bad_resource_descriptor(self):
        bad_rd = '{"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}'

        got_pb = pb_json.Parse(bad_rd, rdpb.ResourceDescriptor())
        got = ResourceDescriptor.copy_from_pb(got_pb)

        with self.assertRaises(
            ValueError,
            msg="Error: created malformed ResourceDescriptor (no required fields)",
        ):
            got.validate()

    def test_empty_name_only(self):
        bad_rd = '{"name":"","downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}'

        got_pb = pb_json.Parse(bad_rd, rdpb.ResourceDescriptor())
        got = ResourceDescriptor.copy_from_pb(got_pb)

        with self.assertRaises(
            ValueError,
            msg="Error: created malformed ResourceDescriptor (only empty required fields)",
        ):
            got.validate()

    def test_empty_digest(self):
        empty_digest = '{"name":"theName","digest":{},"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}'

        got_pb = pb_json.Parse(empty_digest, rdpb.ResourceDescriptor())
        got = ResourceDescriptor.copy_from_pb(got_pb)

        # this should not raise an error
        got.validate()


if __name__ == "__main__":
    unittest.main()
