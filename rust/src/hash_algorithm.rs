//! # in-toto Supported Hash Algorithms
//!
//! This module provides an enum for all cryptographic hash
//! algorithms currently supported for [in-toto DigestSets](https://github.com/in-toto/attestation/blob/main/spec/v1/digest_set.md) and a function to obtain their
//! string represtation.

/// Supported hash algorithms for in-toto DigestSets
#[derive(Debug, PartialEq)]
pub enum HashAlgorithm {
    MD5,
    SHA1,
    SHA224,
    SHA512_224,
    SHA256,
    SHA512_256,
    SHA384,
    SHA512,
    SHA3_224,
    SHA3_256,
    SHA3_384,
    SHA3_512,
    GitBlob,
    GitCommit,
    GitTag,
    GitTree,
    DirHash,
    Unsupported,
}

impl HashAlgorithm {
    /// Converts a given hash algorithm into its string representation.
    pub fn as_str(&self) -> &'static str {
        match self {
            HashAlgorithm::MD5 => "md5",
            HashAlgorithm::SHA1 => "sha1",
            HashAlgorithm::SHA224 => "sha224",
            HashAlgorithm::SHA512_224 => "sha512_224",
            HashAlgorithm::SHA256 => "sha256",
            HashAlgorithm::SHA512_256 => "sha512_256",
            HashAlgorithm::SHA384 => "sha384",
            HashAlgorithm::SHA512 => "sha512",
            HashAlgorithm::SHA3_224 => "sha3_224",
            HashAlgorithm::SHA3_256 => "sha3_256",
            HashAlgorithm::SHA3_384 => "sha3_384",
            HashAlgorithm::SHA3_512 => "sha3_512",
            HashAlgorithm::GitBlob => "gitBlob",
            HashAlgorithm::GitCommit => "gitCommit",
            HashAlgorithm::GitTag => "gitTag",
            HashAlgorithm::GitTree => "gitTree",
            HashAlgorithm::DirHash => "dirHash",
            HashAlgorithm::Unsupported => "unsupported",
        }
    }

    /// Returns the enum value of a hash algorithm given its string representation
    pub fn from_string(alg_str: &str) -> HashAlgorithm {
        match alg_str {
            "md5" => HashAlgorithm::MD5,
            "sha1" => HashAlgorithm::SHA1,
            "sha224" => HashAlgorithm::SHA224,
            "sha512_224" => HashAlgorithm::SHA512_224,
            "sha256" => HashAlgorithm::SHA256,
            "sha512_256" => HashAlgorithm::SHA512_256,
            "sha384" => HashAlgorithm::SHA384,
            "sha512" => HashAlgorithm::SHA512,
            "sha3_224" => HashAlgorithm::SHA3_224,
            "sha3_256" => HashAlgorithm::SHA3_256,
            "sha3_384" => HashAlgorithm::SHA3_384,
            "sha3_512" => HashAlgorithm::SHA3_512,
            "gitBlob" => HashAlgorithm::GitBlob,
            "gitCommit" => HashAlgorithm::GitCommit,
            "gitTag" => HashAlgorithm::GitTag,
            "gitTree" => HashAlgorithm::GitTree,
            "dirHash" => HashAlgorithm::DirHash,
            _ => HashAlgorithm::Unsupported,
        }
    }

    /// Returns the expected length of an algorithm's hash when hex-encoded.
    ///
    /// We assume gitCommit and dirHash are aliases for sha1 and sha256, respectively.
    ///
    /// SHA digest sizes from https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.202.pdf
    /// MD5 digest size from https://www.rfc-editor.org/rfc/rfc1321.html#section-1
    pub fn hex_len(&self) -> usize {
        match self {
            HashAlgorithm::MD5 => 16,
            HashAlgorithm::SHA1
            | HashAlgorithm::GitBlob
            | HashAlgorithm::GitCommit
            | HashAlgorithm::GitTag
            | HashAlgorithm::GitTree => 20,
            HashAlgorithm::SHA224 | HashAlgorithm::SHA512_224 | HashAlgorithm::SHA3_224 => 28,
            HashAlgorithm::SHA256
            | HashAlgorithm::SHA512_256
            | HashAlgorithm::SHA3_256
            | HashAlgorithm::DirHash => 32,
            HashAlgorithm::SHA384 | HashAlgorithm::SHA3_384 => 48,
            HashAlgorithm::SHA512 | HashAlgorithm::SHA3_512 => 64,
            HashAlgorithm::Unsupported => 0,
        }
    }
}
