//! Tests for signing and verifying Statements in DSSE envelopes.
#![cfg(feature = "dsse")]

use in_toto_attestation::envelope::dsse::{self, Ed25519Signer, Verifier};
use in_toto_attestation::envelope::{PAYLOAD_TYPE, sign_statement, verify_statement};
use in_toto_attestation::error::Error;
use in_toto_attestation::generate_statement_v1;
use in_toto_attestation::v1::resource_descriptor::ResourceDescriptor;
use in_toto_attestation::v1::statement::Statement;
use protobuf::well_known_types::struct_::Struct;

fn statement() -> Statement {
    let mut subject = ResourceDescriptor::new();
    subject.name = "app.tar.gz".to_string();
    subject.digest.insert(
        "sha256".to_string(),
        "a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4e5f60718293a4b5c6d7e8f90".to_string(),
    );
    generate_statement_v1(
        &[subject],
        "https://example.com/predicate/v1",
        &Struct::new(),
    )
    .unwrap()
}

fn signer(seed: u8, key_id: &str) -> Ed25519Signer {
    Ed25519Signer::from_bytes(&[seed; 32])
        .unwrap()
        .with_key_id(key_id)
}

fn envelope_error<T: std::fmt::Debug>(result: Result<T, Error>) -> String {
    match result {
        Err(Error::EnvelopeError(msg)) => msg,
        other => panic!("expected an envelope error, got {other:?}"),
    }
}

#[test]
fn signed_statement_verifies_and_round_trips() {
    let a = signer(1, "a");
    let envelope = sign_statement(&statement(), &a).unwrap();
    let verified = verify_statement(&envelope, &[&a.verifier() as &dyn Verifier], 1).unwrap();
    assert_eq!(verified, statement());
}

#[test]
fn envelope_carries_the_in_toto_payload_type() {
    let envelope = sign_statement(&statement(), &signer(1, "a")).unwrap();
    let parsed = dsse::Envelope::from_json(&envelope).unwrap();
    assert_eq!(parsed.payload_type, PAYLOAD_TYPE);
}

#[test]
fn untrusted_key_is_refused() {
    let envelope = sign_statement(&statement(), &signer(1, "a")).unwrap();
    let other = signer(2, "b").verifier();
    envelope_error(verify_statement(&envelope, &[&other as &dyn Verifier], 1));
}

#[test]
fn tampered_payload_is_refused() {
    let a = signer(1, "a");
    let envelope = sign_statement(&statement(), &a).unwrap();
    let mut parsed = dsse::Envelope::from_json(&envelope).unwrap();
    let mut other = statement();
    other.predicate_type = "https://example.com/predicate/v2".to_string();
    let forged = sign_statement(&other, &a).unwrap();
    parsed.payload = dsse::Envelope::from_json(&forged).unwrap().payload;
    let tampered = parsed.to_json().unwrap();
    envelope_error(verify_statement(
        &tampered,
        &[&a.verifier() as &dyn Verifier],
        1,
    ));
}

#[test]
fn other_payload_type_is_refused_even_when_signed() {
    let a = signer(1, "a");
    let payload = protobuf_json_mapping::print_to_string(&statement()).unwrap();
    let envelope = dsse::sign("application/json", payload.as_bytes(), &a)
        .unwrap()
        .to_json()
        .unwrap();
    let msg = envelope_error(verify_statement(
        &envelope,
        &[&a.verifier() as &dyn Verifier],
        1,
    ));
    assert!(msg.contains("application/json"), "{msg}");
}

#[test]
fn one_key_signing_twice_does_not_meet_a_threshold_of_two() {
    let a = signer(1, "a");
    let b = signer(2, "b");
    let payload = protobuf_json_mapping::print_to_string(&statement()).unwrap();
    let envelope = dsse::sign_with(PAYLOAD_TYPE, payload.as_bytes(), &[&a, &a])
        .unwrap()
        .to_json()
        .unwrap();
    let keys = [&a.verifier() as &dyn Verifier, &b.verifier()];
    envelope_error(verify_statement(&envelope, &keys, 2));
}

#[test]
fn two_distinct_keys_meet_a_threshold_of_two() {
    let a = signer(1, "a");
    let b = signer(2, "b");
    let payload = protobuf_json_mapping::print_to_string(&statement()).unwrap();
    let envelope = dsse::sign_with(PAYLOAD_TYPE, payload.as_bytes(), &[&a, &b])
        .unwrap()
        .to_json()
        .unwrap();
    let keys = [&a.verifier() as &dyn Verifier, &b.verifier()];
    assert_eq!(verify_statement(&envelope, &keys, 2).unwrap(), statement());
}

#[test]
fn verified_payload_that_is_not_a_statement_is_a_parse_error() {
    let a = signer(1, "a");
    let envelope = dsse::sign(PAYLOAD_TYPE, b"not json", &a)
        .unwrap()
        .to_json()
        .unwrap();
    match verify_statement(&envelope, &[&a.verifier() as &dyn Verifier], 1) {
        Err(Error::ParseError(_)) => {}
        other => panic!("expected a parse error, got {other:?}"),
    }
}
