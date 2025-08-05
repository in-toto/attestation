use protobuf::well_known_types::struct_::Struct;
/// # in-toto Attestation Framework Rust Bindings
///
/// Bindings for generating unsigned in-toto attestations
/// compliant with the [in-toto Attestation Framework spec](https://github.com/in-toto/attestation/tree/main/spec#in-toto-attestation-framework-spec).
/// Please note: The v1 and predicates APIs are auto-generated. DO NOT EDIT
///
/// # Example
///
/// ```
/// use in_toto_attestation::v1::statement::Statement;
/// use in_toto_attestation::v1::resource_descriptor::ResourceDescriptor;
/// use in_toto_attestation::predicates::link::v0::link::Link;
/// use in_toto_attestation::{generate_statement, to_struct};
/// use protobuf::well_known_types::struct_::Struct;
/// use protobuf_json_mapping::print_to_string;
/// use std::collections::HashMap;
///
/// let digest_m = HashMap::from([("alg".to_string(), "01dfab13".to_string())]);
/// let digest_s = HashMap::from([("alg".to_string(), "de4db33f".to_string())]);
///
/// let mut materials = ResourceDescriptor::new();
/// materials.name = "hello-world.c".to_string();
/// materials.digest = digest_m;
///
/// let mut subject = ResourceDescriptor::new();
/// subject.name = "hello-world".to_string();
/// subject.digest = digest_s;
///
/// let mut link = Link::new();
/// link.name = "hello world".to_string();
/// link.materials = [materials].to_vec();
///
/// let link_struct = to_struct(&link).unwrap();
///
/// let statement = generate_statement(&[subject], "https://in-toto.io/attestation/link/v0.3", &link_struct).unwrap();
///
/// let statement_json = print_to_string(&statement).unwrap();
/// println!("JSON statement: {}", statement_json.as_str());
/// ```
use protobuf::{MessageDyn, MessageField};
use protobuf_json_mapping::{parse_from_str, print_to_string};

pub mod error;
pub mod hash_algorithm;
pub mod predicates;
pub mod v1;
pub mod validator;

use error::{Error, Result};
use hash_algorithm::HashAlgorithm;
use v1::resource_descriptor::ResourceDescriptor;
use v1::statement::Statement;
use validator::MetadataValidator;

/// The in-toto Statement v1 type URI
const STATEMENT_TYPE_URI_V1: &str = "https://in-toto.io/Statement/v1";

/// Utility function to convert any in-toto attestation data structure
/// into a protobuf-compatible Struct type.
pub fn to_struct(proto_msg: &dyn MessageDyn) -> Result<Struct> {
    let msg_json =
        print_to_string(proto_msg).map_err(|e| Error::SerializationError(e.to_string()))?;

    let msg_struct =
        parse_from_str::<Struct>(&msg_json).map_err(|e| Error::ParseError(e.to_string()))?;

    Ok(msg_struct)
}

/// Utility function to generate an in-toto Statement
pub fn generate_statement(
    subject: &[ResourceDescriptor],
    predicate_type: &str,
    predicate: &Struct,
) -> Result<Statement> {
    let mut statement = Statement::new();
    statement.type_ = STATEMENT_TYPE_URI_V1.to_string();
    statement.subject = subject.to_vec();
    statement.predicate_type = predicate_type.to_string();

    statement.predicate = MessageField::some(predicate.clone());

    Ok(statement)
}

impl MetadataValidator for ResourceDescriptor {
    /// Field validator for in-toto [ResourceDescriptor](https://github.com/in-toto/attestation/blob/main/spec/v1/resource_descriptor.md) metadata.
    fn validate_fields(&self) -> Result<bool> {
        // at least one of name, URI or digest are required
        if self.name.is_empty() && self.uri.is_empty() && self.digest.is_empty() {
            return Err(Error::FieldValidationError(
                "at least one of name, URI, or digest are required".to_string(),
            ));
        }

        if !self.digest.is_empty() {
            for (alg, digest) in &self.digest {
                // Per https://github.com/in-toto/attestation/blob/main/spec/v1/digest_set.md
                // check encoding and length for supported algorithms;
                // use of custom, unsupported algorithms is allowed and does not not generate
                // validation errors.
                let supported = HashAlgorithm::from_string(alg);
                if supported != HashAlgorithm::Unsupported {
                    // the in-toto spec expects a hex-encoded string in DigestSets for supported
                    // algorithms
                    let hash_bytes = hex::decode(digest).map_err(|e| {
                        Error::FieldValidationError(
                            format_args!("{e}: ({alg}: {digest})").to_string(),
                        )
                    })?;

                    // check the length of the digest
                    if hash_bytes.len() != supported.hex_len() {
                        return Err(Error::FieldValidationError(
                            format_args!("digest has incorrect length: ({alg}:{digest})")
                                .to_string(),
                        ));
                    }
                }
            }
        }

        Ok(true)
    }
}

impl MetadataValidator for Statement {
    /// Field validator for in-toto [Statement](https://github.com/in-toto/attestation/blob/main/spec/v1/statement.md) metadata.
    fn validate_fields(&self) -> Result<bool> {
        if self.type_ != STATEMENT_TYPE_URI_V1 {
            return Err(Error::FieldValidationError(
                "wrong statement type".to_string(),
            ));
        }

        let subject = &self.subject;
        if subject.is_empty() {
            return Err(Error::FieldValidationError(
                "at least one subject required".to_string(),
            ));
        }

        // Check all resource descriptors in the subject
        for rd in subject.iter() {
            rd.validate_fields()?;

            // v1 statements require the digest to be set in the subject
            if rd.digest.is_empty() {
                return Err(Error::FieldValidationError(
                    "at least one digest required".to_string(),
                ));
            }
        }

        if self.predicate_type.is_empty() {
            return Err(Error::FieldValidationError(
                "predicate type required".to_string(),
            ));
        }

        // an unset predicate field is considered set-but-empty
        // see: https://github.com/in-toto/attestation/blob/main/spec/v1/statement.md#fields

        Ok(true)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use protobuf::MessageField;
    use std::collections::HashMap;

    fn create_test_resource_desc() -> v1::resource_descriptor::ResourceDescriptor {
        let digest = HashMap::from([("alg".to_string(), "01dfab13".to_string())]);

        let mut rd = v1::resource_descriptor::ResourceDescriptor::new();
        rd.name = "hello-world.c".to_string();
        rd.uri = "https://example.com".to_string();
        rd.digest = digest;
        rd.content = "bytecontent".as_bytes().to_vec();
        rd.download_location = "https://example.com/test.zip".to_string();
        rd.media_type = "theMediaType".to_string();

        rd
    }

    fn create_test_statement() -> v1::statement::Statement {
        let digest_s = HashMap::from([("alg".to_string(), "de4db33f".to_string())]);

        let materials = create_test_resource_desc();

        let mut subject = v1::resource_descriptor::ResourceDescriptor::new();
        subject.name = "hello-world".to_string();
        subject.digest = digest_s;

        let mut link = predicates::link::v0::link::Link::new();
        link.name = "hello world".to_string();
        link.materials = [materials].to_vec();

        let link_struct = to_struct(&link).unwrap();

        let statement = generate_statement(
            &[subject],
            "https://in-toto.io/attestation/link/v0.3",
            &link_struct,
        )
        .unwrap();

        statement
    }

    #[test]
    fn test_create_resource_descriptor() {
        let rd = create_test_resource_desc();

        assert!(rd.validate_fields().expect("This RD should be valid."));

        let rd_json = print_to_string(&rd).unwrap();
        println!("JSON RD: {}", rd_json.as_str());
    }

    #[test]
    fn test_supported_resource_descriptor_digest() {
        let digest = HashMap::from([
            (
                "sha256".to_string(),
                "a1234567b1234567c1234567d1234567e1234567f1234567a1234567b1234567".to_string(),
            ),
            ("alg".to_string(), "01dfab13".to_string()),
            (
                "gitCommit".to_string(),
                "a1234567b1234567c1234567d1234567e1234567".to_string(),
            ),
        ]);

        let mut rd = create_test_resource_desc();
        rd.digest = digest;

        assert!(rd.validate_fields().expect("This RD should be valid."));
    }

    #[test]
    fn test_bad_resource_descriptor() -> Result<()> {
        let mut bad_rd = v1::resource_descriptor::ResourceDescriptor::new();
        bad_rd.content = "bytecontent".as_bytes().to_vec();
        bad_rd.download_location = "https://example.com/test.zip".to_string();
        bad_rd.media_type = "theMediaType".to_string();

        match bad_rd.validate_fields() {
            Ok(true) => Err(Error::FieldValidationError(
                "RD with missing required fields should throw an error".to_string(),
            )),
            _ => Ok(()),
        }
    }

    #[test]
    fn test_bad_resource_descriptor_digest_encoding() -> Result<()> {
        let bad_digest = HashMap::from([("sha256".to_string(), "bad_digest".to_string())]);

        let mut bad_rd = create_test_resource_desc();
        bad_rd.digest = bad_digest;

        match bad_rd.validate_fields() {
            Ok(true) => Err(Error::FieldValidationError(
                "Bad RD digest encoding should throw an error".to_string(),
            )),
            _ => Ok(()),
        }
    }

    #[test]
    fn test_bad_resource_descriptor_digest_length() -> Result<()> {
        let bad_digest =
            HashMap::from([("sha256".to_string(), "a1234567b1234567c123".to_string())]);

        let mut bad_rd = create_test_resource_desc();
        bad_rd.digest = bad_digest;

        match bad_rd.validate_fields() {
            Ok(true) => Err(Error::FieldValidationError(
                "Bad RD digest length should throw an error".to_string(),
            )),
            _ => Ok(()),
        }
    }

    #[test]
    fn test_create_statement() {
        let statement = create_test_statement();

        assert!(
            statement
                .validate_fields()
                .expect("This Statement should be valid.")
        );

        let statement_json = print_to_string(&statement).unwrap();
        println!("JSON statement: {}", statement_json.as_str());
    }

    #[test]
    fn test_statement_empty_predicate() {
        let mut statement = create_test_statement();
        statement.predicate = MessageField::some(Struct::new());

        assert!(
            statement
                .validate_fields()
                .expect("A Statement with an empty predicate should be valid.")
        );
    }

    #[test]
    fn test_statement_unset_predicate() {
        let mut statement = create_test_statement();
        statement.predicate = MessageField::none();

        assert!(
            statement
                .validate_fields()
                .expect("A Statement with no predicate object should be valid.")
        );
    }

    #[test]
    fn test_bad_statement_type() -> Result<()> {
        let mut bad_statement = create_test_statement();
        bad_statement.type_ = "https://example.com/in-toto".to_string();

        match bad_statement.validate_fields() {
            Ok(true) => Err(Error::FieldValidationError(
                "Statement with unexpected type should throw an error".to_string(),
            )),
            _ => Ok(()),
        }
    }

    #[test]
    fn test_bad_statement_no_subject() -> Result<()> {
        let mut bad_statement = create_test_statement();
        bad_statement.subject = Vec::new();

        match bad_statement.validate_fields() {
            Ok(true) => Err(Error::FieldValidationError(
                "Statement without subject should throw an error".to_string(),
            )),
            _ => Ok(()),
        }
    }

    #[test]
    fn test_bad_statement_no_subject_digest() -> Result<()> {
        let mut bad_rd = v1::resource_descriptor::ResourceDescriptor::new();
        bad_rd.name = "badSubjectRd".to_string();
        bad_rd.download_location = "https://example.com/test.zip".to_string();
        bad_rd.media_type = "theMediaType".to_string();

        let mut bad_statement = create_test_statement();
        bad_statement.subject = [bad_rd].to_vec();

        match bad_statement.validate_fields() {
            Ok(true) => Err(Error::FieldValidationError(
                "Statement without subject digest should throw an error".to_string(),
            )),
            _ => Ok(()),
        }
    }

    #[test]
    fn test_bad_statement_predicate_type() -> Result<()> {
        let mut bad_statement = create_test_statement();
        bad_statement.predicate_type = "".to_string();

        match bad_statement.validate_fields() {
            Ok(true) => Err(Error::FieldValidationError(
                "Statement without predicate type should throw an error".to_string(),
            )),
            _ => Ok(()),
        }
    }
}
