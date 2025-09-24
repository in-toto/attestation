# Validation model

The following pseudocode shows how to verify and extract metadata about a
single artifact from a single attestation. The expectation is that consumers
will feed the resulting metadata into a policy engine.

> **TODO**: Explain how to process multiple artifacts and/or multiple attestations.

Inputs:

-   `artifactToVerify`: blob of data
-   `attestation`: JSON-encoded [Envelope]
-   `recognizedAttesters`: collection of (`name`, `publicKey`) pairs
-   `acceptableDigestAlgorithms`: collection of acceptable cryptographic hash
    algorithms (usually just `sha256`)

Steps:

-   Envelope layer:
-   `envelope` := decode `attestation` as a JSON-encoded [Envelope];
        reject if decoding fails
    -   `attesterNames` := empty set of names
    -   For each `signature` in `envelope.signatures`:
        -   For each (`name`, `publicKey`) in `recognizedAttesters`:
        -   Optional: skip if `signature.keyid` does not match
            `publicKey`
            -   If `signature.sig` matches `publicKey`:
                -   Add `name` to `attesterNames`
    -   Reject if `attesterNames` is empty
-   Intermediate state: `envelope.payloadType`, `envelope.payload`,
    `attesterNames`
-   Statement layer:
    -   Reject if (`envelope.payloadType` != `application/vnd.in-toto+json` OR `envelope.payloadType` does not match `application/vnd.in-toto.<predicate>+json`)
    -   `statement` := decode `envelope.payload` as a JSON-encoded
        [Statement]; reject if decoding fails
    -   Reject if `statement.type` != `https://in-toto.io/Statement/v1`
    -   `matchedSubjects` := the subset of entries `s` in `statement.subject`
        where:
        -   there exists at least one `(alg, value)` in `s.digest` where:
            -   `alg` is in `acceptableDigestAlgorithms` AND
            -   `hash(alg, artifactToVerify)` == `hexDecode(value)`
    -   Reject if `matchedSubjects` is empty

Output (to be fed into policy engine):

-   `statement.predicateType`
-   `statement.predicate`
-   `matchedSubjects`
-   `attesterNames`
