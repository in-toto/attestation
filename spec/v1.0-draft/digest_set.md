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

<dl>

<dt>

`sha256`, `sha224`, `sha384`, `sha512`, `sha512_224`, `sha512_256`,
`sha3_224`, `sha3_256`, `sha3_384`, `sha3_512`, `shake128`, `shake256`,
`blake2b`, `blake2s`, `ripemd160`, `sm3`, `gost`, `sha1`, `md5`

<dd>

Standard cryptographic hash algorithms using [the NIST names][] (converting to
lowercase and replacing `-` with `_`) as keys and lowercase hex-encoded values,
for cases when the method of serialization is obvious or well known:

<dt>

`dirHash`

<dd>

The [directory Hash1][] function, omitting the "h1:" prefix
and output in lowercase hexadecimal instead of base64. This algorithm was
designed for go modules but can be used to digest the _contents_ of an
arbitrary archive or file tree. Equivalent to extracting the archive to an
empty directory and running the following command in that directory:

```bash
find . -type f | cut -c3- | LC_ALL=C sort | xargs -r sha256sum | sha256sum | cut -f1 -d' '
```

For example, the module dirhash
`h1:Khu2En+0gcYPZ2kuIihfswbzxv/mIHXgzPZ018Oty48=` would be encoded as
`{"dirHash1": "2a1bb6127fb481c60f67692e22285fb306f3c6ffe62075e0ccf674d7c3adcb8f"}`.

<details>
<summary>Detailed example</summary>

The go module `github.com/marklodato/go-hello-world@v0.0.1` has module
dirhash `h1:Khu2En+0gcYPZ2kuIihfswbzxv/mIHXgzPZ018Oty48=`:

```bash
$ curl https://sum.golang.org/lookup/github.com/marklodato/go-hello-world@v0.0.1
...
github.com/marklodato/go-hello-world v0.0.1 h1:Khu2En+0gcYPZ2kuIihfswbzxv/mIHXgzPZ018Oty48=
...
```

To compute the dirhash by hand, first fetch the module archive and extract
it to an empty directory:

```bash
curl -O https://proxy.golang.org/github.com/marklodato/go-hello-world/@v/v0.0.1.zip
mkdir tmp
cd tmp
unzip ../v0.0.1.zip
```

We can see all of the files in the directory using the first part of the
command above:

```bash
$ find . -type f | cut -c3- | LC_ALL=C sort | xargs -r sha256sum
3a137eef6458bfb76bb2c63fc29ffc7166604d2d2e09ed9d8250a534122a8364  github.com/marklodato/go-hello-world@v0.0.1/README.md
28e7c942a036902d981759d0bf5704d2bfc7cb500caf68b84711b234af01c6a5  github.com/marklodato/go-hello-world@v0.0.1/go.mod
ddc4da627d9a9f45fb29641a1b185d6f53287ecfd921aacbf4fe54b7a86fe8d1  github.com/marklodato/go-hello-world@v0.0.1/main.go
```

The dirhash is the sha256 sum over the output of the previous command:

```bash
$ find . -type f | cut -c3- | LC_ALL=C sort | xargs -r sha256sum | sha256sum | cut -f1 -d' '
2a1bb6127fb481c60f67692e22285fb306f3c6ffe62075e0ccf674d7c3adcb8f
```

This is equivalent to the base64 encoded version:

```bash
$ echo '2a1bb6127fb481c60f67692e22285fb306f3c6ffe62075e0ccf674d7c3adcb8f' | xxd -r -p | {printf 'h1:'; base64}
h1:Khu2En+0gcYPZ2kuIihfswbzxv/mIHXgzPZ018Oty48=
```

</details>

</dl>

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
