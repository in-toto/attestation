//! # in-toto Attestation Errors
//!
//! This module defines custom error types and a result type alias for the
//! in-toto attestation data structure library.
//!
//! The `Error` enum provides a comprehensive set of error variants to represent
//! different kinds of errors that can occur in the application.
//!
//! ## Example Usage
//!
//! ```
//! use in_toto_attestation::error::{Error, Result};
//!
//! fn example_function() -> Result<()> {
//!     Err(Error::FieldValidationError("This metadata field has an issue".to_string()))
//! }
//!
//! match example_function() {
//!     Ok(_) => println!("Operation succeeded"),
//!     Err(e) => eprintln!("Error occurred: {}", e),
//! }
//! ```
use thiserror::Error;

/// Represents the various errors that can occur in the application.
#[derive(Debug, Error)]
pub enum Error {
    /// Represents a metadata field validation error.
    ///
    /// This variant includes a string describing the validation error.
    #[error("Field validation error: {0}")]
    FieldValidationError(String),

    /// Represents an error that occurs during parsing of serialized data.
    ///
    /// This variant includes a string describing the parsing error.
    #[error("Parse error: {0}")]
    ParseError(String),

    /// Represents an error that occurs during data serialization.
    ///
    /// This variant includes a string describing the serialization error.
    #[error("Serialization error: {0}")]
    SerializationError(String),
}

/// A type alias for results that use the custom `Error` type.
///
/// This alias simplifies function signatures by using the `Error` enum as the
/// error type in `std::result::Result`.
pub type Result<T> = std::result::Result<T, Error>;
