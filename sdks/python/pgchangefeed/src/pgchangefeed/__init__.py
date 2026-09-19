"""PG Change Feed Python client library.

This release exposes only the shared connection configuration
(`ClientOptions`) that every wire surface needs regardless of transport --
address and bearer token (ADR-0107 Festlegung 1/3/4). The HTTP API client
surface itself (the SPEC-018 capabilities) is added by a follow-up release.
"""

from pgchangefeed.options import ClientOptions

__all__ = ["ClientOptions"]
