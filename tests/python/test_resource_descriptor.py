'''
Tests for in-toto attestation ResourceDescriptor protos.
'''

import unittest
import google.protobuf.json_format as pb_json

import in_toto_attestation.v1.resource_descriptor_pb2 as rdpb
from in_toto_attestation.v1.resource_descriptor import ResourceDescriptor

def create_test_desc():
    desc = rdpb.ResourceDescriptor()
    desc.name = 'theName'
    desc.uri = 'https://example.com'
    desc.digest['alg'] = 'abc123'
    desc.content = b'bytescontent'
    desc.downloadLocation = 'https://example.com/test.zip'
    desc.mediaType = 'theMediaType'
    desc.annotations['a1'].update({'keyStr': 'val1', 'keyNum': 13})
    desc.annotations['a2'].update({'keyObj': {'subKey': 'subVal'}})
    return ResourceDescriptor(desc)

class TestResourceDescriptor(unittest.TestCase):
    def setUp(self):
        self.want_full_rd = '{"name":"theName","uri":"https://example.com","digest":{"alg":"abc123"},"content":"Ynl0ZXNjb250ZW50","downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType","annotations":{"a1":{"keyNum": 13,"keyStr":"val1"},"a2":{"keyObj":{"subKey":"subVal"}}}}'

        self.bad_rd = '{"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}'

        self.bad_digest = '{"name":"theName","digest":{},"downloadLocation":"https://example.com/test.zip","mediaType":"theMediaType"}'

        self.test_rd = create_test_desc()

    def test_create_resource_descriptor(self):
        self.test_rd.validate()
        
    def test_json_parse_resource_descriptor(self):
        got_pb = pb_json.Parse(self.want_full_rd, rdpb.ResourceDescriptor())
        got = got_pb.SerializeToString(deterministic=True)

        want = self.test_rd.pb.SerializeToString(deterministic=True)

        self.assertEqual(got, want, 'Protos do not match')

    def test_bad_resource_descriptor(self):
        got_pb = pb_json.Parse(self.bad_rd, rdpb.ResourceDescriptor())
        got = ResourceDescriptor(got_pb)

        with self.assertRaises(ValueError, msg='Error: created malformed ResourceDescriptor'):
            got.validate()

    def test_bad_digest(self):
        got_pb = pb_json.Parse(self.bad_digest, rdpb.ResourceDescriptor())
        got = ResourceDescriptor(got_pb)

        with self.assertRaises(ValueError, msg='Error: created ResourceDescriptor with malformed digest field'):
            got.validate()

if __name__ == '__main__':
    unittest.main()
