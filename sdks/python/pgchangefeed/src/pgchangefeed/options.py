"""Shared connection configuration for the PG Change Feed clients.

``ClientOptions`` holds the two values every client needs regardless of the
transport: the server address and the bearer token. The HTTP, gRPC, SSE and
NATS clients all take the same class.
"""

from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True)
class ClientOptions:
    """Connection configuration for a PG Change Feed client.

    Attributes:
        address: Address of the PG Change Feed server, in the form the client
            needs: the HTTP base URL for the HTTP and SSE clients, ``host:port``
            for the gRPC client, a ``nats://`` URL for the NATS client.
        api_token: Bearer token sent as the authorization credential: as a
            header by the HTTP and SSE clients, as call metadata by the gRPC
            client, and once when the connection is opened by the NATS client.

    Raises:
        ValueError: ``address`` or ``api_token`` is empty or blank.
    """

    address: str
    api_token: str

    def __post_init__(self) -> None:
        if not self.address or not self.address.strip():
            raise ValueError("address must not be empty")
        if not self.api_token or not self.api_token.strip():
            raise ValueError("api_token must not be empty")
