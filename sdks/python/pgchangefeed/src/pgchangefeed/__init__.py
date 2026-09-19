"""PG Change Feed Python client library.

Package skeleton (slice-sdk-python-projektgeruest, ADR-0107 Festlegung
1/3/4). The HTTP API client surface itself (the SPEC-018 capabilities) is
added by the follow-up slice (slice-sdk-python-http-client-flaeche) -- this
release exposes only the shared connection configuration.
"""

from pgchangefeed.options import ClientOptions

__all__ = ["ClientOptions"]
