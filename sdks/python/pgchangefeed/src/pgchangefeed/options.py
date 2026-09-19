"""Shared connection configuration for PG Change Feed client surfaces.

ADR-0107 Festlegung 1 scopes the first Python package release to the HTTP
API (SPEC-018) only -- no gRPC/SSE/NATS surface exists yet in this package.
This class still holds only the two values every wire surface needs
regardless of transport (base address, bearer token), so a later surface
can reuse it without a breaking change to this constructor -- no
anticipation of endpoint methods themselves (analogy:
sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs).
"""

from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True)
class ClientOptions:
    """Connection configuration for a PG Change Feed client surface.

    Attributes:
        address: Base address of the PG Change Feed server (the HTTP
            endpoint, for the surface currently covered by this package).
        api_token: Bearer token sent as an authorization credential
            (SPEC-018).
    """

    address: str
    api_token: str

    def __post_init__(self) -> None:
        if not self.address or not self.address.strip():
            raise ValueError("address must not be empty")
        if not self.api_token or not self.api_token.strip():
            raise ValueError("api_token must not be empty")
