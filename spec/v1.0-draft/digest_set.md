# DigestSet field type specification

Version v1.0-draft

Set of one or more cryptographic digests for a single software artifact or
metadata object.

## Schema

```json
{
  "<ALGORITHM_1>": "<HEX/BASE64 VALUE>",
  "<ALGORITHM_2>": "<HEX/BASE64 VALUE>",
  ... 
}
```

## Fields

A DigestSet is represented as a _JSON object_ mapping algorithm name to
a string encoding of the digest using that algorithm. The named standard
algorithms below use lowercase hex encoding. Usually there is just a
single key/value pair, but multiple entries MAY be used for algorithm
agility.

Supported algorithms:

-   Standard cryptographic hash algorithms using [the NIST names][]
    (converting to lowercase and replacing `-` replaced with `_`) as keys
    and lowercase hex-encoded values, for cases when the method of
    serialization is obvious or well known:
    `sha256`, `sha224`, `sha384`, `sha512`, `sha512_224`, `sha512_256`,
    `sha3_224`, `sha3_256`, `sha3_384`, `sha3_512`, `shake128`, `shake256`,
    `blake2b`, `blake2s`, `ripemd160`, `sm3`, `gost`, `sha1`, `md5`

-   `goModuleH1`: The go module [directory Hash1][], omitting the "h1:"
    prefix and output in lowercase hexadecimal instead of base64. Can
    be computed over a directory named `name@version`, or the contents
    of zip file containing such a directory:

    ```bash
    find name@version -type f | LC_ALL=C sort | xargs -r sha256sum | sha256sum | cut -f1 -d' '
    ```

    For example, the module dirhash
    `h1:47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=` would be encoded as
    `{"goModuleH1": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}`.

> It is RECOMMENDED to use at least `sha256` for compatibility between
> producers and consumers.
>
> Consumers MUST only accept algorithms that they consider secure and MUST
> ignore unrecognized or unaccepted algorithms. For example, most
> applications SHOULD NOT accept "md5" because it lacks collision resistance.
>
> Two DigestSets SHOULD be considered matching if ANY acceptable field
> matches.
>
> New algorithms MUST document how the value is encoded, e.g. URL-safe base64,
> lowercase hex, etc...

## Examples

-   `{"sha256": "abcd", "sha512": "1234"}` matches `{"sha256": "abcd"}`
-   `{"sha256": "abcd"}` does not match `{"sha256": "fedb", "sha512": "abcd"}`
-   `{"somecoolhash": "abcd"}` uses a non-predefined algorithm

[the NIST names]: https://csrc.nist.gov/projects/hash-functions
[directory Hash1]: https://cs.opensource.google/go/x/mod/+/refs/tags/v0.5.0:sumdb/dirhash/hash.go
