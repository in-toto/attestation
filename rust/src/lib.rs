/// Bindings for generating unsigned in-toto attestations
/// compliant with the in-toto Attestation Framework spec.
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
use protobuf_json_mapping::{PrintError, parse_from_str, print_to_string};

pub mod predicates;
pub mod v1;

/// Utility function to convert any in-toto attestation data structure
/// into a protobuf-compatible Struct type.
pub fn to_struct(proto_msg: &dyn MessageDyn) -> Result<Struct, PrintError> {
    let msg_json = print_to_string(proto_msg)?;

    let msg_struct =
        parse_from_str::<Struct>(&msg_json).map_err(|_| PrintError::from(std::fmt::Error))?;

    Ok(msg_struct)
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
