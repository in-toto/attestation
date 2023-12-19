#  Wrapper class for in-toto attestation ResourceDescriptor protos.

from __future__ import annotations

from in_toto_attestation._internal.in_toto_attestation.v1 import (
    ResourceDescriptor as _ResourceDescriptor,
)


class ResourceDescriptor:
    def __init__(
        self,
        name: str = "",
        uri: str = "",
        digest: dict | None = None,
        content: bytes = b"",
        download_location: str = "",
        media_type: str = "",
        annotations: dict | None = None,
    ) -> None:
        self.pb = _ResourceDescriptor()
        self.pb.name = name
        self.pb.uri = uri
        self.pb.digest = digest
        self.pb.content = content
        self.pb.download_location = download_location
        self.pb.media_type = media_type
        if annotations:
            self.pb.annotations.fields = annotations

        self._validate()

    @classmethod
    def from_pb(cls, raw: bytes) -> ResourceDescriptor:
        pb = _ResourceDescriptor.parse(raw)
        return cls(
            pb.name,
            pb.uri,
            pb.digest,
            pb.content,
            pb.download_location,
            pb.media_type,
            pb.annotations,
        )

    @classmethod
    def from_json(cls, raw: str) -> ResourceDescriptor:
        pb = _ResourceDescriptor().from_json(raw)
        return cls(
            pb.name,
            pb.uri,
            pb.digest,
            pb.content,
            pb.download_location,
            pb.media_type,
            pb.annotations,
        )

    def to_json(self) -> str:
        self.pb.to_json()

    def _validate(self) -> None:
        # at least one of name, URI or digest are required
        if self.pb.name == "" and self.pb.uri == "" and len(self.pb.digest) == 0:
            raise ValueError("At least one of name, URI, or digest need to be set")
