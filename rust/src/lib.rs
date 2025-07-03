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
/// use in_toto_attestation::to_struct;
/// use protobuf::MessageField;
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
/// let mut statement = Statement::new();
/// statement.type_ = "https://in-toto.io/Statement/v1".to_string();
/// statement.subject = [subject].to_vec();
/// statement.predicate_type = "https://in-toto.io/attestation/link/v0.3".to_string();
/// let predicate_struct = to_struct(&link).unwrap();
/// statement.predicate = MessageField::some(predicate_struct);
///
/// let statement_json = print_to_string(&statement).unwrap();
/// println!("JSON statement: {}", statement_json.as_str());
/// ```
use protobuf::MessageDyn;
use protobuf::well_known_types::struct_::Struct;
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

        if self.predicate.is_none() {
            return Err(Error::FieldValidationError(
                "predicate object required".to_string(),
            ));
        }

        Ok(true)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use protobuf::MessageField;
    use std::collections::HashMap;

    #[test]
    fn test_create_simple_statement() {
        let digest_m = HashMap::from([("alg".to_string(), "01dfab13".to_string())]);
        let digest_s = HashMap::from([("alg".to_string(), "de4db33f".to_string())]);

        let mut materials = v1::resource_descriptor::ResourceDescriptor::new();
        materials.name = "hello-world.c".to_string();
        materials.digest = digest_m;

        let mut subject = v1::resource_descriptor::ResourceDescriptor::new();
        subject.name = "hello-world".to_string();
        subject.digest = digest_s;

        let mut link = predicates::link::v0::link::Link::new();
        link.name = "hello world".to_string();
        link.materials = [materials].to_vec();

        let mut statement = v1::statement::Statement::new();
        statement.type_ = "https://in-toto.io/Statement/v1".to_string();
        statement.subject = [subject].to_vec();
        statement.predicate_type = "https://in-toto.io/attestation/link/v0.3".to_string();
        let predicate_struct = to_struct(&link).unwrap();
        statement.predicate = MessageField::some(predicate_struct);

        let statement_json = print_to_string(&statement).unwrap();
        println!("JSON statement: {}", statement_json.as_str());
    }
}
