'''
Tests for in-toto attestation Statement protos.
'''

import unittest
import google.protobuf.json_format as pb_json

import in_toto_attestation.v1.statement_pb2 as stpb
import in_toto_attestation.v1.resource_descriptor_pb2 as rdpb

import in_toto_attestation.v1.statement as its

def create_test_statement():
    sub = rdpb.ResourceDescriptor()
    sub.name = 'theSub'
    sub.digest['alg1'] = 'abc123'

    st = stpb.Statement()
    st.type = its.STATEMENT_TYPE_URI
    st.subject.append(sub)
    st.predicateType = 'thePredicate'
    st.predicate.update({'keyObj': {'subKey': 'subVal'}})
    return its.Statement(st)

class TestStatement(unittest.TestCase):
    def setUp(self):
        self.want_st = '{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}'

        self.test_st = create_test_statement()

    def test_create_statement(self):
        self.test_st.validate()
        
    def test_json_parse_statement(self):
        got_pb = pb_json.Parse(self.want_st, stpb.Statement())
        got = got_pb.SerializeToString(deterministic=True)

        want = self.test_st.pb.SerializeToString(deterministic=True)

        self.assertEqual(got, want, 'Protos do not match')

    def test_bad_statement_type(self):
        bad_st = '{"_type":"https://in-toto.io/Statement/v0","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}'
        
        got_pb = pb_json.Parse(bad_st, stpb.Statement())
        got = its.Statement(got_pb)

        with self.assertRaises(ValueError, msg='Error: created malformed Statement (bad type)'):
            got.validate()

    def test_bad_statement_empty_subject(self):
        bad_st = '{"_type":"https://in-toto.io/Statement/v1","subject":[],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}'
        
        got_pb = pb_json.Parse(bad_st, stpb.Statement())
        got = its.Statement(got_pb)

        with self.assertRaises(ValueError, msg='Error: created malformed Statement (empty subject)'):
            got.validate()

    def test_bad_statement_bad_subject(self):
        bad_st = '{"_type":"https://in-toto.io/Statement/v1","subject":[{"digest":{}}],"predicateType":"thePredicate","predicate":{"keyObj":{"subKey":"subVal"}}}'
        
        got_pb = pb_json.Parse(bad_st, stpb.Statement())
        got = its.Statement(got_pb)

        with self.assertRaises(ValueError, msg='Error: created malformed Statement (bad subject)'):
            got.validate()

    def test_bad_predicate_type(self):
        bad_st = '{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"","predicate":{"keyObj":{"subKey":"subVal"}}}'
        
        got_pb = pb_json.Parse(bad_st, stpb.Statement())
        got = its.Statement(got_pb)

        with self.assertRaises(ValueError, msg='Error: created malformed Statement (bad predicate type)'):
            got.validate()

    def test_bad_predicate(self):
        bad_st = '{"_type":"https://in-toto.io/Statement/v1","subject":[{"name":"theSub","digest":{"alg1":"abc123"}}],"predicateType":"thePredicate"}'
        
        got_pb = pb_json.Parse(bad_st, stpb.Statement())
        got = its.Statement(got_pb)

        with self.assertRaises(ValueError, msg='Error: created malformed Statement (bad predicate)'):
            got.validate()

if __name__ == '__main__':
    unittest.main()
