# Predicate type: EIP-712 Payload Attestation

Type URI: https://payperbyte.io/attestation/eip712-payload-attestation/v1

Version: 1.0.0

Predicate Name: EIP-712 Payload Attestation (ForeSeal receipt)

Authors: BYTEDev Inc (PayPerByte / ForeSeal) — DRAFT 2026-08-20, not yet submitted

## Purpose

This predicate carries an **EIP-712 `PayloadAttestation`** — a secp256k1 typed-data
signature over the exact bytes of a served payload — inside an in-toto Statement, so
that a signed-bytes receipt produced in the Ethereum/x402 ecosystem can be consumed by
in-toto/DSSE-native policy engines without re-signing or translation.

The receipt proves **which key signed exactly these bytes, and until when the signer
stood behind that delivery**. It is an authenticity and tamper-evidence record. It
makes **no claim about the correctness or truth of the payload's contents** — that
distinction is load-bearing and is repeated in the Fields section.

The EIP-712 struct, domain, and encoding are those already deployed as the
PayPerByte / ForeSeal `X-BYTE-Attestation` HTTP receipt (domain name "BYTE Library",
struct `PayloadAttestation{publisher,payloadHash,payloadLength,deadline}`), but the
predicate is generic: any producer that signs `PayloadAttestation` typed data under
any EIP-712 domain can use it — the domain is part of the predicate, and verifiers pin
the domains they trust.

## Use Cases

1.  **Verify-before-act for AI agents.** An agent buys a data verdict (e.g. a
   sanctions screen, a package-install verdict) over HTTP and receives the bytes plus
   an `X-BYTE-Attestation` receipt. Before acting, it (or a policy gate in front of
   it) re-hashes the bytes, checks the signature, and checks the signer against an
   allowlist. Wrapping that receipt as an in-toto Statement lets the same evidence be
   stored, transported, and evaluated by in-toto/DSSE tooling alongside the agent's
   other supply-chain attestations.
2.  **Provenance handoff across organizations.** A data publisher signs its payload
   (provenance tier); a delivery intermediary signs the exact bytes it forwarded
   (delivery tier). Two Statements with the same subject digest and different
   `publisher` values describe the chain without either party re-keying.
3.  **Audit trail for paid machine-to-machine calls.** Receipts are retained as
   Statements keyed by payload digest, so a later dispute ("what bytes did the seller
   actually serve for this call?") resolves by digest match plus signature recovery.

**Why existing predicates do not fit.** There is no predicate in the directory for a
secp256k1/EIP-712 signature artifact, and no predicate whose subject is "the exact
bytes served" with a signer-declared validity deadline. SLSA Provenance describes how
an artifact was built; SCAI and Simple Verification Result describe evaluations of an
artifact; none carries an independently verifiable third-party typed-data signature
whose verification key is recovered from the signature itself (no certificate, no
key identifier). Nesting the EIP-712 signature as predicate content (rather than as
the DSSE envelope signature) is deliberate: DSSE signs the Statement; the EIP-712
signature inside remains verifiable on its own terms, with its own domain separation.

## Prerequisites

-   The [in-toto Attestation Framework](https://github.com/in-toto/attestation)
  (Statement v1, DSSE envelope).
-   [EIP-712](https://eips.ethereum.org/EIPS/eip-712) typed structured data hashing and
  signing, and secp256k1 ECDSA public-key recovery (Ethereum `ecrecover`).
-   keccak-256 (the Ethereum Keccak, not SHA3-256) for `payloadHash`.

## Model

Functionaries and steps this predicate applies to:

-   **Signer / attester** (the `publisher` field): the party holding the secp256k1 key
  that signs the typed data. Two tiers are common: *provenance* (the originator of the
  bytes) and *delivery* (an intermediary vouching for the exact bytes it served).
-   **Consumer / verifier**: the party that received the bytes and the receipt and wants
  evidence-toward authenticity and integrity before acting, storing, or paying.
-   **Subject**: the exact payload bytes as served (the HTTP response body, byte for
  byte — producers MUST NOT re-serialize before hashing). Producers SHOULD include both
  a standard digest (`sha256`) and the Ethereum hash (`keccak256`) of the same bytes so
  that in-toto tooling can match the subject by `sha256` while the EIP-712 `payloadHash`
  can be matched directly to `keccak256`.

Verification procedure (normative for consumers):

1.  Obtain the exact payload bytes `B`.
2.  Check `len(B) == predicate.payloadLength`.
3.  Check `keccak256(B) == predicate.payloadHash` (and, if present, equals
   `subject[].digest.keccak256`; `sha256(B)` equals `subject[].digest.sha256`).
4.  Recover the signer: `recover(EIP712Hash(domain, PayloadAttestation{publisher,
   payloadHash, payloadLength, deadline}), signature)` and check it equals
   `predicate.publisher`.
5.  Policy decisions (outside this spec, but the predicate must make them answerable):
   is `publisher` an accepted key for this resource/tier; is `domain` one the verifier
   trusts; is `deadline` acceptable for the consumer's freshness requirement.

Any failure in steps 2–4 means the Statement's predicate is **not verified** — it is
not "partially verified". Step 5 is policy, not validity: an expired `deadline` does
not invalidate the signature as historical evidence, but a verify-before-act consumer
SHOULD treat it as stale.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "<resource URI the bytes were served from, optional>",
      "digest": {
        "sha256": "<hex>",
        "keccak256": "<hex>"     // Ethereum keccak-256 of the same bytes; equals predicate.payloadHash
      }
    }
  ],

  // Predicate:
  "predicateType": "https://payperbyte.io/attestation/eip712-payload-attestation/v1",
  "predicate": {
    "alg": "EIP712-PayloadAttestation",
    "domain": {
      "name": "<EIP-712 domain name>",
      "version": "<EIP-712 domain version>",
      "chainId": <integer>,
      "verifyingContract": "<0x-prefixed 20-byte hex>"
    },
    "publisher": "<0x-prefixed 20-byte hex address of the signer>",
    "payloadHash": "<0x-prefixed 32-byte hex, keccak256 of the exact bytes>",
    "payloadLength": <integer, byte length of the exact bytes>,
    "deadline": <integer, Unix seconds; the signed validity deadline>,
    "signature": "<0x-prefixed 65-byte hex r||s||v>",
    "tier": "provenance" | "delivery",               // optional
    "expiresAt": "<RFC 3339 UTC, derived from deadline>",   // optional, informative
    "attestedAt": "<RFC 3339 UTC>"                   // optional, informative
  }
}
```

### Parsing Rules

This predicate follows the in-toto Attestation Framework's
[standard parsing rules](../v1/README.md#parsing-rules), with these clarifications:

-   **Signed fields are authoritative.** `publisher`, `payloadHash`, `payloadLength`,
  `deadline`, `domain`, and `signature` reproduce the typed data exactly as signed.
  Verifiers MUST recompute the EIP-712 hash from these fields; they MUST NOT trust
  `expiresAt`, `attestedAt`, `tier`, or any subject metadata as evidence — those are
  convenience/annotation fields and are not covered by the EIP-712 signature (they are
  covered only by the DSSE envelope signature, if any).
-   `deadline` is kept as a Unix-seconds integer (not RFC 3339) because it is the
  `uint256` value inside the signed struct; altering its representation would break
  recovery. `expiresAt` is its RFC 3339 rendering for human/tool convenience.
-   The EIP-712 type definition is fixed for this predicate version:
  `PayloadAttestation(address publisher,bytes32 payloadHash,uint256 payloadLength,uint256 deadline)`.
  A future version that changes the struct MUST bump the predicate major version.
-   `alg` MUST be the literal `EIP712-PayloadAttestation` for this version; unknown
  `alg` values MUST cause the verifier to reject the predicate (fail-closed).
-   Unknown additional fields MUST be ignored (monotonic principle), but MUST NOT be
  treated as evidence.
-   Hex strings are case-insensitive on input; producers SHOULD emit lowercase hex for
  hashes/signatures and EIP-55 checksum casing for addresses.
-   **`keccak256` DigestSet encoding (required by the DigestSet rules for non-standard
  algorithm names):** in `subject[].digest`, `keccak256` is the Ethereum Keccak-256 of
  the exact bytes, encoded as 64 lowercase hex characters **without** a `0x` prefix
  (matching the other DigestSet entries). In the predicate, `payloadHash` is the same
  32 bytes **with** the `0x` prefix, as produced by EVM tooling. Consumers that do not
  recognize `keccak256` MUST ignore it (DigestSet rule) and can still verify via
  `payloadHash` + `sha256`.

### Fields

`alg` *string*, *required*

Literal `EIP712-PayloadAttestation`. Identifies the signing scheme and the struct.

`domain` *object*, *required*

The EIP-712 domain separator parameters used when signing:
`name` (string), `version` (string), `chainId` (integer), `verifyingContract`
(20-byte hex address). The domain is part of what the signature commits to; verifiers
MUST pin the domains they accept (a receipt under an unexpected domain is not evidence
for the expected one). `chainId` here identifies the domain anchor and does not imply
anything about where payment for the payload settled.

`publisher` *string (address)*, *required*

The signer's Ethereum address. This is the value recovered from `signature` over the
typed data; it is both a signed field and the expected recovery result. The name
`publisher` is historical (it is the struct's field name on the deployed verifier);
for the delivery tier it identifies the delivering attester rather than the data's
originator.

`payloadHash` *string (bytes32 hex)*, *required*

keccak-256 of the exact subject bytes. MUST equal `subject[].digest.keccak256` when
that digest is present.

`payloadLength` *integer*, *required*

Byte length of the exact subject bytes. Included in the signed struct as a cheap
second integrity check and to defeat length-extension-style confusion.

`deadline` *integer (Unix seconds)*, *required*

The signed validity deadline: the signer stands behind "I served exactly these bytes"
until this instant. Freshness policy is the consumer's; the signature remains
cryptographically valid as historical evidence after the deadline.

`signature` *string (65-byte hex)*, *required*

secp256k1 ECDSA signature (`r||s||v`) over the EIP-712 digest of the typed data.

`tier` *enum*, *optional*

`provenance` (signer is the originator of the bytes) or `delivery` (signer is an
intermediary vouching for the exact bytes it served). Informative; not signed.

`expiresAt` *Timestamp (RFC 3339 UTC)*, *optional*

`deadline` rendered as RFC 3339 (e.g. `2026-08-18T20:00:00Z`). Informative; not signed.

`attestedAt` *Timestamp (RFC 3339 UTC)*, *optional*

When the receipt was produced, if known. Informative; not signed.

**Scope statement (normative for producers' documentation):** this predicate is
evidence-toward *authenticity and tamper-evidence of the exact bytes* and the signer's
*validity window*. It MUST NOT be described as certifying the correctness, quality, or
truth of the payload's contents.

## Example

### Example 1 — minimal (ephemeral key)

A complete Statement follows. The signature is real and re-verifies with the
procedure in **Model** (length, keccak256, signer recovery). The subject is the exact
292-byte JSON body below. The signer is an **ephemeral example key** generated for this
document (it signs nothing else); production receipts are signed by the producer's own
attester key.

Exact subject bytes (UTF-8, no trailing newline):

```text
{"feed":"sanctions-screen","query":{"address":"0x1111111111111111111111111111111111111111"},"verdict":"ALLOW","list_state":{"source":"OFAC SDN + Consolidated","date":"2026-08-19","sha256":"3b7f0c5a9d2e4f6a8c1b3d5e7f9a0b2c4d6e8f0a1b3c5d7e9f1a3b5c7d9e0f2a"},"checked_at":"2026-08-20T20:00:00Z"}
```

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "https://x402.payperbyte.io/feeds/sanctions-screen",
      "digest": {
        "sha256": "e08fb5b71183bec5943322bd56d2e60c5c6d2e17ed6987f4f44ca96970d1e8b4",
        "keccak256": "a32d8a7c4a3486fd297e181841f7e12e9cf84c5b944fde3f2e4543ceae7b555b"
      }
    }
  ],
  "predicateType": "https://payperbyte.io/attestation/eip712-payload-attestation/v1",
  "predicate": {
    "alg": "EIP712-PayloadAttestation",
    "domain": {
      "name": "BYTE Library",
      "version": "1",
      "chainId": 421614,
      "verifyingContract": "0x44729bB148F46d8Db509E47b0453edc271e06e95"
    },
    "publisher": "0xB9700F44D3458681cbE0C7f95d281BD2cA29e532",
    "payloadHash": "0xa32d8a7c4a3486fd297e181841f7e12e9cf84c5b944fde3f2e4543ceae7b555b",
    "payloadLength": 292,
    "deadline": 1787342400,
    "signature": "0x34c2b618faae4e8207e172d7f003fd69725c0c9de4dedcd4cb4c7f5c91b637e42459592d02f60de69978b9f43e10a075a511e156886c9c9f2563d6641ffeae301b",
    "tier": "delivery",
    "expiresAt": "2026-08-21T20:00:00Z",
    "attestedAt": "2026-08-20T20:01:00Z"
  }
}
```

Re-verification sketch (viem, TypeScript):

```ts
import { keccak256, recoverTypedDataAddress } from "viem";
const types = { PayloadAttestation: [
  { name: "publisher", type: "address" }, { name: "payloadHash", type: "bytes32" },
  { name: "payloadLength", type: "uint256" }, { name: "deadline", type: "uint256" } ] };
const bytes = new TextEncoder().encode(body);               // the exact subject bytes
const p = statement.predicate;
const ok =
  bytes.length === p.payloadLength &&
  keccak256(bytes) === p.payloadHash &&
  (await recoverTypedDataAddress({ domain: p.domain, types, primaryType: "PayloadAttestation",
     message: { publisher: p.publisher, payloadHash: p.payloadHash,
                payloadLength: BigInt(p.payloadLength), deadline: BigInt(p.deadline) },
     signature: p.signature })).toLowerCase() === p.publisher.toLowerCase();
// ok === true for the example above; then apply policy: accepted publisher? accepted domain? deadline fresh enough?
```

### Example 2 — production receipt (sanctions-screen)

A complete Statement built from a real production receipt captured from the live
gateway (2026-08-21T03:55:08.648Z, `https://x402.payperbyte.io/feeds/sanctions-screen`). The signature is real and re-verifies
with the procedure in **Model** (length, keccak256, signer recovery) — checked
automatically by `evidence/regen_from_capture.mjs` before this example was
generated. The signer, `0xB48CCc9e3ab67041e3b5D09700138E45cda6AeA8`, is the gateway's live delivery attester
(not an ephemeral example key). The subject is the exact 3312-byte JSON
body below.

Exact subject bytes (UTF-8 — the fence content below PLUS one trailing LF (0x0A); total 3312 bytes, sha256 `2fae9140f9b95c09bf7e711cd93086b24464373579710522e774d001d7fa9be3`):

```text
{"answer":{"v":"sanctions-screen/v1","ts":1787284507,"query":{"address":"0x833589fcd6edb6e08f4c7c32d4f71b54bda02913","name":null,"chain":null},"verdict":"ALLOW","score":100,"reasons":["no match on the OFAC SDN list (19249 entries; list published 2026-08-20, fetched 2026-08-20T23:42:38Z, sha256 50213298d936901a\u2026)","no match on the OFAC Consolidated (non-SDN) list (481 entries; list published 2026-08-20, fetched 2026-08-20T23:42:41Z, sha256 5a629469398539ac\u2026)"],"signals":{"sdn":{"list_available":true,"address_hit":false,"address_matches":[],"name_exact_hit":false,"name_exact_matches":[],"name_fuzzy_hit":false,"name_fuzzy_matches":[],"list_state":{"source":"OFAC SDN (Specially Designated Nationals and Blocked Persons)","source_url":"https://sanctionslistservice.ofac.treas.gov/api/PublicationPreview/exports/SDN.CSV","published_date":"2026-08-20","fetched_at":"2026-08-20T23:42:38Z","content_sha256":"50213298d936901a1aaad7bb19c968dab9e82fa07e8c808aacfae8fcea3d870e","entry_count":19249,"age_days":0,"stale":false},"error":null},"consolidated":{"list_available":true,"address_hit":false,"address_matches":[],"name_exact_hit":false,"name_exact_matches":[],"name_fuzzy_hit":false,"name_fuzzy_matches":[],"list_state":{"source":"OFAC Consolidated (non-SDN) Sanctions List","source_url":"https://sanctionslistservice.ofac.treas.gov/api/PublicationPreview/exports/CONS_PRIM.CSV","published_date":"2026-08-20","fetched_at":"2026-08-20T23:42:41Z","content_sha256":"5a629469398539aca2d180a086543e2161d1203fb2a3c9c737b1d682544df5b1","entry_count":481,"age_days":0,"stale":false},"error":null}},"list_state":{"sdn":{"source":"OFAC SDN (Specially Designated Nationals and Blocked Persons)","source_url":"https://sanctionslistservice.ofac.treas.gov/api/PublicationPreview/exports/SDN.CSV","published_date":"2026-08-20","fetched_at":"2026-08-20T23:42:38Z","content_sha256":"50213298d936901a1aaad7bb19c968dab9e82fa07e8c808aacfae8fcea3d870e","entry_count":19249,"age_days":0,"stale":false},"consolidated":{"source":"OFAC Consolidated (non-SDN) Sanctions List","source_url":"https://sanctionslistservice.ofac.treas.gov/api/PublicationPreview/exports/CONS_PRIM.CSV","published_date":"2026-08-20","fetched_at":"2026-08-20T23:42:41Z","content_sha256":"5a629469398539aca2d180a086543e2161d1203fb2a3c9c737b1d682544df5b1","entry_count":481,"age_days":0,"stale":false}},"retrieved_at":"2026-08-21T03:55:07Z","methodology":"ss-v1","input_hashes":{"sdn":"0x9bcbbaa69c4040ffc3513afab8080366074718c49a18071fc2c9fd865af6f8d4","consolidated":"0xeee09b9059da1fb85d29fdc01bd7e7c6d3b8d76d3712a5cd7937e2bb66d0469b"},"source":"OFAC SDN + OFAC Consolidated (non-SDN) via sanctionslistservice.ofac.treas.gov (official Treasury exports)","error":null},"broadcast":{"ok":false,"tx":null,"delivered":0,"note":"broadcast disabled (SANCTIONS_SCREEN_BROADCAST=0)"},"attestation":{"payloadHash":"0xbe58daa362cf94a4b4d6dc90c8415c306c06d69eedb5f599a69e14e62cc79464","payloadLength":2720,"deadline":2102644507,"signer":"0x344ECaCDe6566294c31397445c98b62a3EEEA456","signature":"0xb63cf806e4d74bc8323de684f502ceda8c04e2e3bbc049dc9d631bd276214dac02dfbb22d624b28deff1359b857e5110247d4c2f7e147232f2fddaca0a084ed21b","domain":{"name":"BYTE Library","version":"1","chainId":421614,"verifyingContract":"0x44729bB148F46d8Db509E47b0453edc271e06e95"}}}
```

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "https://x402.payperbyte.io/feeds/sanctions-screen",
      "digest": {
        "sha256": "2fae9140f9b95c09bf7e711cd93086b24464373579710522e774d001d7fa9be3",
        "keccak256": "b14ef4b30838a2964800ace5f02f592834e14c695be5862b54b6ff8d2e1647d3"
      }
    }
  ],
  "predicateType": "https://payperbyte.io/attestation/eip712-payload-attestation/v1",
  "predicate": {
    "alg": "EIP712-PayloadAttestation",
    "domain": {
      "name": "BYTE Library",
      "version": "1",
      "chainId": 421614,
      "verifyingContract": "0x44729bB148F46d8Db509E47b0453edc271e06e95"
    },
    "publisher": "0xB48CCc9e3ab67041e3b5D09700138E45cda6AeA8",
    "payloadHash": "0xb14ef4b30838a2964800ace5f02f592834e14c695be5862b54b6ff8d2e1647d3",
    "payloadLength": 3312,
    "deadline": 2102644507,
    "signature": "0x575399d1e3f8fdcfc5586c93be797951a83802b718621aa1e1d938dbf56f443434e1ec4cb18bead30bc4ea2582f587ee1797ea8580334bcd42d682ed5eea6cf11c",
    "tier": "delivery",
    "expiresAt": "2036-08-18T03:55:07Z",
    "attestedAt": "2026-08-21T03:55:08.648Z"
  }
}
```

Re-verification: `evidence/regen_from_capture.mjs` re-derived this exact receipt
(length, keccak256, EIP-712 signer recovery) from `./attestation-capture-sanctions-2026-08-21.json` and scanned its
body for degraded-response markers before emitting this example (zero matches).
The check output is recorded in `capture_verify.json` alongside this file.

## Changelog and Migrations

-   1.0.0 (draft 2026-08-20): initial version. Struct and domain encoding are identical
  to the deployed PayPerByte/ForeSeal `X-BYTE-Attestation` HTTP receipt, so existing
  receipts wrap into this predicate without re-signing.
