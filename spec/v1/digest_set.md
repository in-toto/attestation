# DigestSet field type specification

Version: v1.1

Set of one or more immutable digests for a single software artifact or
metadata object.

## Schema

```json
{
  "<ALGORITHM_1>": "<VALUE>",
  "<ALGORITHM_2>": "<VALUE>",
  ... 
}
```

## Fields

A DigestSet is represented as a _JSON object_ mapping algorithm name to
a string encoding of the digest using that algorithm. The named standard
algorithms below use lowercase hex encoding. Usually there is just a
single key/value pair, but multiple entries MAY be used for algorithm
agility.

Users SHOULD use a _cryptographic_ digest, but MAY use another identifier
if the underlying implementation ensures immutability via other means.

### Supported algorithms

#### `sha256`, `sha224`, `sha384`, `sha512`, `sha512_224`, `sha512_256`, `sha3_224`, `sha3_256`, `sha3_384`, `sha3_512`, `shake128`, `shake256`, `blake2b`, `blake2s`, `ripemd160`, `sm3`, `gost`, `sha1`, `md5`

Standard cryptographic hash algorithms using [the NIST names][] (converting to
lowercase and replacing `-` with `_`) as keys and lowercase hex-encoded values,
for cases when the method of serialization is obvious or well known.

#### `dirHash`

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

#### `gitCommit`, `gitTree`, `gitBlob`, `gitTag`

The lowercase hex SHA-1 (40 character) or SHA-256 (64 character) of a git
commit, tree, blob, or tag object, respectively. The `gitTree` and `gitBlob` in
particular can be used for arbitrary trees or files, even outside git.

This hash is computed over `<type> SP <size> NUL <content>`, where:

-   `<type>` is one of `commit`, `tree`, `blob`, `tag`
-   `SP` is the ASCII space character, 0x20
-   `<size>` is the number of bytes in `<content>`, represented as a decimal
    ASCII number with no leading zeros
-   `NUL` is the ASCII NUL character, 0x00
-   `<content>` is git representation of the object:
    -   For `commit`, the raw commit object ([more info][so-commit][^git-docs])
    -   For `tree`,  the raw tree object, which is a series of
        `<unix-octal-mode> <name> NUL <binary-digest>` entries, sorted by
        `<name>` in the C locale ([more info][so-tree][^git-docs])
    -   For `blob`, the raw file contents
    -   For `tag`, the raw tag object

For more information, see [Git Objects] in the Git Book.[^git-docs]

Example of `gitBlob` for the file containing the 5 bytes `Hello` (no newline):

```bash
$ printf 'blob 5\0Hello' | sha1sum | cut -f1 -d' '
5ab2f8a4323abafb10abb68657d9d39f1a775057
$ printf 'Hello' | git hash-object -t blob --stdin
5ab2f8a4323abafb10abb68657d9d39f1a775057
```

### Guidelines

It is RECOMMENDED to use at least `sha256` for compatibility between
producers and consumers, unless a different hash algorithm is more
conventional (e.g. `gitCommit` for git).

Consumers MUST only accept algorithms that they consider secure and MUST
ignore unrecognized or unaccepted algorithms. For example, most
applications SHOULD NOT accept "md5" because it lacks collision resistance.

Two DigestSets SHOULD be considered matching if ANY acceptable field
matches.

New algorithms MUST document how the value is encoded, e.g. URL-safe base64,
lowercase hex, etc...

## Examples

-   `{"sha256": "abcd", "sha512": "1234"}` matches `{"sha256": "abcd"}`
-   `{"sha256": "abcd"}` does not match `{"sha256": "fedb", "sha512": "abcd"}`
-   `{"somecoolhash": "abcd"}` uses a non-predefined algorithm

<!-- Add a horizontal rule to separate footnotes -->

---

[^git-docs]: At the time of writing (2023-03), git has no official documentation
    of the internal object format used for hashing. The [Git Objects]
    chapter of the Git Book is the closest thing to official documentation, but
    it lacks many details, such as the raw tree object format. The best
    documentation we have found are the linked Stack Overflow articles. If you
    can find a better, more official reference, please open an issue.

[Git Objects]: https://git-scm.com/book/en/v2/Git-Internals-Git-Objects
[directory Hash1]: https://cs.opensource.google/go/x/mod/+/refs/tags/v0.5.0:sumdb/dirhash/hash.go
[so-commit]: https://stackoverflow.com/a/37438460
[so-tree]: https://stackoverflow.com/a/35902553
[the NIST names]: https://csrc.nist.gov/projects/hash-functions
