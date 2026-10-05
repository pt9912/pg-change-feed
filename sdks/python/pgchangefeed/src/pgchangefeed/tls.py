"""Connections to the PG Change Feed server that honor ``ClientOptions.trust_anchor_file``.

``create_http_client`` builds the ``httpx.Client`` for the HTTP and SSE clients,
``create_grpc_channel`` the ``grpc.Channel`` for the gRPC stream client and the
administration client. Without a trust anchor the default trust of the runtime
applies (``httpx`` and ``grpc`` use their own bundles); with one, the
connection trusts exactly the certificates of the anchor file. Chain, validity
period and server name (the host of the address) are always checked, and a
failed check surfaces as the connection error of the transport (``httpx``
exception, ``grpc.RpcError`` with status ``UNAVAILABLE``), never as a retry in
plaintext. A client or channel you build yourself keeps its own TLS
configuration; ``trust_anchor_file`` does not apply to it.
"""

from __future__ import annotations

import ssl

import grpc
import httpx

from pgchangefeed.options import ClientOptions

_DEFAULT_HTTP_TIMEOUT = httpx.Timeout(30.0)


def create_http_client(
    options: ClientOptions, timeout: httpx.Timeout | float | None = _DEFAULT_HTTP_TIMEOUT
) -> httpx.Client:
    """Builds the ``httpx.Client`` for ``PgChangeFeedHttpClient`` and ``PgChangeFeedSseClient``.

    ``timeout`` is the ``httpx`` timeout of the client; the SSE stream needs the
    read timeout off (``httpx.Timeout(10.0, read=None)``). The caller closes the
    returned client. A ``trust_anchor_file`` with an address that is not
    ``https://`` raises ``ValueError``.
    """
    if options.trust_anchor_file is None:
        return httpx.Client(timeout=timeout)
    if not options.address.lower().startswith("https://"):
        raise ValueError(f"a trust anchor requires an https address, got '{options.address}'")
    context = ssl.create_default_context(cadata=options._trust_anchor_pem.decode("ascii"))
    return httpx.Client(verify=context, timeout=timeout)


def create_grpc_channel(options: ClientOptions) -> grpc.Channel:
    """Builds the ``grpc.Channel`` for ``PgChangeFeedGrpcClient`` and ``PgChangeFeedAdministrationClient``.

    The address is ``host:port`` or a URL: ``https://host:port`` connects over
    TLS, ``http://host:port`` in plaintext, and ``host:port`` over TLS when a
    ``trust_anchor_file`` is set and in plaintext otherwise. A trust anchor with
    an ``http://`` address raises ``ValueError``. The caller closes the returned
    channel.
    """
    address = options.address.strip()
    lowered = address.lower()
    if lowered.startswith("https://"):
        target, tls = address[len("https://") :], True
    elif lowered.startswith("http://"):
        target, tls = address[len("http://") :], False
    else:
        target, tls = address, options.trust_anchor_file is not None
    target = target.rstrip("/")
    if options.trust_anchor_file is not None and not tls:
        raise ValueError(f"a trust anchor requires a TLS address, got '{options.address}'")
    if not tls:
        return grpc.insecure_channel(target)
    return grpc.secure_channel(target, grpc.ssl_channel_credentials(root_certificates=options._trust_anchor_pem))
