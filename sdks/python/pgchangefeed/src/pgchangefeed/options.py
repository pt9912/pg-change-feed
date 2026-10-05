"""Shared connection configuration for the PG Change Feed clients.

``ClientOptions`` holds the values every client needs regardless of the
transport: the server address, the bearer token and, optionally, a trust
anchor for TLS. The HTTP, gRPC, SSE and NATS clients all take the same class.
"""

from __future__ import annotations

import ssl
from dataclasses import dataclass, field


@dataclass(frozen=True)
class ClientOptions:
    """Connection configuration for a PG Change Feed client.

    Attributes:
        address: Address of the PG Change Feed server, in the form the client
            needs: the HTTP base URL for the HTTP and SSE clients, ``host:port``
            for the gRPC client (``https://host:port`` asks for TLS), a
            ``nats://`` URL for the NATS client.
        api_token: Bearer token sent as the authorization credential: as a
            header by the HTTP and SSE clients, as call metadata by the gRPC
            client, and once when the connection is opened by the NATS client.
        trust_anchor_file: Path of a PEM file with the certificates a TLS
            connection trusts (the certificate of the issuer or of the server),
            or ``None`` when the default trust of the runtime applies. The
            functions ``create_http_client`` and ``create_grpc_channel`` use it;
            a client or channel you build yourself keeps its own TLS
            configuration. Chain, validity period and server name are always
            checked; there is no switch that turns the check off.

    Raises:
        ValueError: ``address`` or ``api_token`` is empty or blank, or
            ``trust_anchor_file`` cannot be read or holds no PEM certificate.
    """

    address: str
    api_token: str
    trust_anchor_file: str | None = None
    _trust_anchor_pem: bytes | None = field(default=None, init=False, repr=False, compare=False)

    def __post_init__(self) -> None:
        if not self.address or not self.address.strip():
            raise ValueError("address must not be empty")
        if not self.api_token or not self.api_token.strip():
            raise ValueError("api_token must not be empty")
        if self.trust_anchor_file is not None:
            object.__setattr__(self, "_trust_anchor_pem", _read_trust_anchor(self.trust_anchor_file))


def _read_trust_anchor(path: str) -> bytes:
    """Reads the PEM bytes of the anchor file and checks that they hold a certificate."""
    try:
        with open(path, "rb") as handle:
            pem = handle.read()
    except OSError as error:
        raise ValueError(f"trust anchor file '{path}' cannot be read: {error}") from error
    context = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    try:
        context.load_verify_locations(cadata=pem.decode("ascii"))
    except (ssl.SSLError, UnicodeDecodeError, ValueError) as error:
        raise ValueError(f"trust anchor file '{path}' holds no PEM certificate: {error}") from error
    return pem
