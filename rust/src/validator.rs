//! # in-toto Metadata Structure Validator
//!
//! This module provides the `MetadataValidator` trait, which consumers of
//! structured in-toto attestation metadata can use to implement schema and
//! field requirements-related validation functions.

use crate::error::Result;

pub trait MetadataValidator {
    fn validate_fields(&self) -> Result<bool>;
}
