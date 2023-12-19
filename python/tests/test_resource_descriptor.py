"""
Tests for in-toto attestation ResourceDescriptor protos.
"""

import unittest

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
        create_test_desc()

    def test_json_parse_resource_descriptor(self):
        full_rd = '{"name":"theName","uri":"https://example.com","digest":{"alg":"abc123"},"content":"Ynl0ZXNjb250ZW50","downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType","annotations":{"keyNum": 13,"keyStr":"val1","keyObj":{"subKey":"subVal"}}}'

        got_pb = ResourceDescriptor.from_json(full_rd)
        got = got_pb.to_json()

        test_rd = create_test_desc()
        want = test_rd.to_json()

        self.assertEqual(got, want, "Protos do not match")

    def test_bad_resource_descriptor(self):
        bad_rd = '{"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}'

        with self.assertRaises(
            ValueError,
            msg="Error: created malformed ResourceDescriptor (no required fields)",
        ):
            ResourceDescriptor.from_json(bad_rd)

    def test_empty_name_only(self):
        bad_rd = '{"name":"","downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}'

        with self.assertRaises(
            ValueError,
            msg="Error: created malformed ResourceDescriptor (only empty required fields)",
        ):
            ResourceDescriptor.from_json(bad_rd)

    def test_empty_digest(self):
        empty_digest = '{"name":"theName","digest":{},"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}'

        # this should not raise an error
        ResourceDescriptor.from_json(empty_digest)


if __name__ == "__main__":
    unittest.main()
