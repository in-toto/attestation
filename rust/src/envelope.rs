//! # DSSE envelopes for in-toto Statements
//!
//! Available with the `dsse` feature. An in-toto attestation is a Statement
//! carried in a [DSSE](https://github.com/secure-systems-lab/dsse) envelope
//! with the payload type `application/vnd.in-toto+json`. This module signs a
//! Statement into an envelope and verifies an envelope back into a Statement.
//!
//! Verification returns the Statement parsed from the exact bytes the
//! signatures cover, so the caller never reads the envelope a second time
//! after verifying it, as the DSSE protocol requires.
//!
//! ## Example Usage
//!
//! ```
//! use in_toto_attestation::envelope::{sign_statement, verify_statement};
//! use in_toto_attestation::envelope::dsse::{Ed25519Signer, Verifier};
//! use in_toto_attestation::generate_statement_v1;
//! use in_toto_attestation::v1::resource_descriptor::ResourceDescriptor;
//! use protobuf::well_known_types::struct_::Struct;
//!
//! let mut subject = ResourceDescriptor::new();
//! subject.name = "app.tar.gz".to_string();
//! subject.digest.insert(
//!     "sha256".to_string(),
//!     "a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4e5f60718293a4b5c6d7e8f90".to_string(),
//! );
//! let statement =
//!     generate_statement_v1(&[subject], "https://example.com/predicate/v1", &Struct::new())?;
//!
//! let signer = Ed25519Signer::from_bytes(&[7u8; 32])?.with_key_id("builder");
//! let envelope = sign_statement(&statement, &signer)?;
//!
//! let verifier = signer.verifier();
//! let verified = verify_statement(&envelope, &[&verifier as &dyn Verifier], 1)?;
//! assert_eq!(verified, statement);
//! # Ok::<(), Box<dyn std::error::Error>>(())
//! ```

use protobuf_json_mapping::{parse_from_str, print_to_string};

pub use dsse;

use crate::error::{Error, Result};
use crate::v1::statement::Statement;

/// The DSSE payload type for an in-toto Statement.
pub const PAYLOAD_TYPE: &str = "application/vnd.in-toto+json";

/// Signs `statement` into a DSSE envelope and returns the envelope's JSON bytes.
pub fn sign_statement(statement: &Statement, signer: &dyn dsse::Signer) -> Result<Vec<u8>> {
    let payload =
        print_to_string(statement).map_err(|e| Error::SerializationError(e.to_string()))?;
    let envelope = dsse::sign(PAYLOAD_TYPE, payload.as_bytes(), signer)
        .map_err(|e| Error::EnvelopeError(e.to_string()))?;
    envelope
        .to_json()
        .map_err(|e| Error::SerializationError(e.to_string()))
}

/// Verifies a DSSE envelope and returns the Statement it carries.
///
/// The envelope must carry signatures from at least `threshold` distinct keys
/// in `verifiers`, and its authenticated payload type must be
/// [`PAYLOAD_TYPE`]. The Statement is parsed from the verified payload bytes.
/// Field validation is left to the caller, through
/// [`MetadataValidator`](crate::validator::MetadataValidator).
pub fn verify_statement(
    envelope_json: &[u8],
    verifiers: &[&dyn dsse::Verifier],
    threshold: usize,
) -> Result<Statement> {
    let envelope = dsse::Envelope::from_json(envelope_json)
        .map_err(|e| Error::EnvelopeError(e.to_string()))?;
    let verified = dsse::verify(&envelope, verifiers, threshold)
        .map_err(|e| Error::EnvelopeError(e.to_string()))?;
    if verified.payload_type != PAYLOAD_TYPE {
        return Err(Error::EnvelopeError(format!(
            "payload type {:?} is not {PAYLOAD_TYPE:?}",
            verified.payload_type
        )));
    }
    let payload =
        std::str::from_utf8(&verified.payload).map_err(|e| Error::ParseError(e.to_string()))?;
    parse_from_str::<Statement>(payload).map_err(|e| Error::ParseError(e.to_string()))
}
