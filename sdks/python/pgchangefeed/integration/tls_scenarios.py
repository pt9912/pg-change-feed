"""The three TLS refusals every TLS phase proves against the running feed container.

No trust anchor (the certificate of the server is unknown to the runtime), a
foreign trust anchor, and a server name that is not in the certificate (the
container name, which resolves in the Docker network but is not listed in the
certificate; the anchor is the right one). Each must end with a failed TLS
check -- a refusal for any other reason (an unreachable name, a refused port)
does not count.
"""

from __future__ import annotations

import os
import ssl
from typing import Callable

import grpc
import httpx

from pgchangefeed.exceptions import PgChangeFeedGrpcError
from pgchangefeed.options import ClientOptions

CA_FILE = os.environ["PGCHANGEFEED_TLS_CA_FILE"]
FOREIGN_CA_FILE = os.environ["PGCHANGEFEED_TLS_FOREIGN_CA_FILE"]
MISMATCH_HOST = os.environ["PGCHANGEFEED_TLS_MISMATCH_HOST"]

_GRPC_TLS_MARKERS = ("handshake", "hostname verification", "certificate", "ssl")


def with_mismatch_host(address: str) -> str:
    """The address with its host replaced by the name that is not in the certificate."""
    scheme, _, rest = address.partition("://") if "://" in address else ("", "", address)
    host, _, port = rest.partition(":")
    mismatch = f"{MISMATCH_HOST}:{port}" if port else MISMATCH_HOST
    return f"{scheme}://{mismatch}" if scheme else mismatch


def _causes(error: BaseException | None) -> list[BaseException]:
    chain: list[BaseException] = []
    while error is not None and error not in chain:
        chain.append(error)
        error = error.__cause__ or error.__context__
    return chain


def is_verification_failure(error: BaseException) -> bool:
    """True when the error is a failed TLS check of the transport.

    ``httpx`` reports it as a ``ConnectError`` caused by an ``ssl.SSLCertVerificationError``;
    ``grpc`` as an ``UNAVAILABLE`` call whose details name the handshake, the
    certificate or the hostname verification.
    """
    for cause in _causes(error):
        if isinstance(cause, ssl.SSLCertVerificationError):
            return True
        if isinstance(cause, grpc.RpcError) and hasattr(cause, "code") and cause.code() == grpc.StatusCode.UNAVAILABLE:
            details = (cause.details() or "").lower()
            if any(marker in details for marker in _GRPC_TLS_MARKERS):
                return True
    return False


def expect_refusals(address: str, token: str, call: Callable[[ClientOptions], None]) -> None:
    """Runs ``call`` with the three option sets, asserts a TLS check failure for each and prints the runner marker."""
    cases = {
        "no_anchor": ClientOptions(address=address, api_token=token),
        "foreign_anchor": ClientOptions(address=address, api_token=token, trust_anchor_file=FOREIGN_CA_FILE),
        "name_mismatch": ClientOptions(
            address=with_mismatch_host(address), api_token=token, trust_anchor_file=CA_FILE
        ),
    }
    for name, options in cases.items():
        try:
            call(options)
        except (httpx.HTTPError, grpc.RpcError, PgChangeFeedGrpcError) as error:
            assert is_verification_failure(error), f"{name}: not a TLS check failure: {error!r}"
        else:
            raise AssertionError(f"{name}: the call succeeded")
    print("REJECTED tls no_anchor=ok foreign_anchor=ok name_mismatch=ok", flush=True)
