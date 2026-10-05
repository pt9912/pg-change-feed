"""The four surfaces that connect over TLS, against servers with a certificate created at run time.

With the trust anchor of ``ClientOptions`` the HTTP client, the SSE client, the
gRPC stream client and the administration client connect. A foreign anchor, a
server name that is not in the certificate, an expired certificate and no
anchor at all (the certificate is unknown to the runtime) each end with the
connection error of the transport; the anchor option is checked when the
options are created; a plaintext address without an anchor works unchanged,
and an anchor with a plaintext address is refused. Certificates are created in
a temporary directory per test; nothing is committed.
"""

from __future__ import annotations

import datetime
import json
import ssl
import threading
from concurrent import futures
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from typing import Iterator

import grpc
import httpx
import pytest
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.x509.oid import NameOID

from pgchangefeed import (
    ClientOptions,
    PgChangeFeedAdministrationClient,
    PgChangeFeedGrpcClient,
    PgChangeFeedHttpClient,
    PgChangeFeedSseClient,
    create_grpc_channel,
    create_http_client,
)
from pgchangefeed.exceptions import PgChangeFeedGrpcUnexpectedStatusError
from pgchangefeed.grpc_gen import (
    administration_pb2,
    administration_pb2_grpc,
    changestream_pb2,
    changestream_pb2_grpc,
)

SENTINEL = "tls-test-change"
SURFACES = ["http", "sse", "grpc_stream", "grpc_admin"]


class Certificate:
    """A self-signed certificate that is its own trust anchor, with its PEM files."""

    def __init__(self, directory: Path, name: str, dns: str, expired: bool = False) -> None:
        key = ec.generate_private_key(ec.SECP256R1())
        subject = x509.Name([x509.NameAttribute(NameOID.COMMON_NAME, f"pgcf-test-{name}")])
        now = datetime.datetime.now(datetime.timezone.utc)
        start, end = (
            (now - datetime.timedelta(days=3), now - datetime.timedelta(days=2))
            if expired
            else (now - datetime.timedelta(minutes=5), now + datetime.timedelta(days=1))
        )
        certificate = (
            x509.CertificateBuilder()
            .subject_name(subject)
            .issuer_name(subject)
            .public_key(key.public_key())
            .serial_number(x509.random_serial_number())
            .not_valid_before(start)
            .not_valid_after(end)
            .add_extension(x509.SubjectAlternativeName([x509.DNSName(dns)]), critical=False)
            .add_extension(x509.BasicConstraints(ca=True, path_length=None), critical=True)
            .sign(key, hashes.SHA256())
        )
        self.cert_pem = certificate.public_bytes(serialization.Encoding.PEM)
        self.key_pem = key.private_bytes(
            serialization.Encoding.PEM,
            serialization.PrivateFormat.PKCS8,
            serialization.NoEncryption(),
        )
        self.cert_path = directory / f"{name}.pem"
        self.cert_path.write_bytes(self.cert_pem)
        self.server_cert_path = self.cert_path
        self.key_path = directory / f"{name}-key.pem"
        self.key_path.write_bytes(self.key_pem)


class IssuedCertificate(Certificate):
    """A server certificate issued by a certificate authority of its own.

    ``cert_path`` is the authority certificate, which is the trust anchor; the
    server presents ``server_cert_pem`` and ``server_key_pem``.
    """

    def __init__(self, directory: Path, name: str, dns: str) -> None:
        authority_key = ec.generate_private_key(ec.SECP256R1())
        authority_name = x509.Name([x509.NameAttribute(NameOID.COMMON_NAME, f"pgcf-test-ca-{name}")])
        now = datetime.datetime.now(datetime.timezone.utc)
        authority = (
            x509.CertificateBuilder()
            .subject_name(authority_name)
            .issuer_name(authority_name)
            .public_key(authority_key.public_key())
            .serial_number(x509.random_serial_number())
            .not_valid_before(now - datetime.timedelta(minutes=5))
            .not_valid_after(now + datetime.timedelta(days=1))
            .add_extension(x509.BasicConstraints(ca=True, path_length=None), critical=True)
            .add_extension(x509.SubjectKeyIdentifier.from_public_key(authority_key.public_key()), critical=False)
            .add_extension(
                x509.KeyUsage(
                    digital_signature=False,
                    content_commitment=False,
                    key_encipherment=False,
                    data_encipherment=False,
                    key_agreement=False,
                    key_cert_sign=True,
                    crl_sign=False,
                    encipher_only=False,
                    decipher_only=False,
                ),
                critical=True,
            )
            .sign(authority_key, hashes.SHA256())
        )
        server_key = ec.generate_private_key(ec.SECP256R1())
        server = (
            x509.CertificateBuilder()
            .subject_name(x509.Name([x509.NameAttribute(NameOID.COMMON_NAME, f"pgcf-test-server-{name}")]))
            .issuer_name(authority_name)
            .public_key(server_key.public_key())
            .serial_number(x509.random_serial_number())
            .not_valid_before(now - datetime.timedelta(minutes=5))
            .not_valid_after(now + datetime.timedelta(days=1))
            .add_extension(x509.SubjectAlternativeName([x509.DNSName(dns)]), critical=False)
            .add_extension(x509.BasicConstraints(ca=False, path_length=None), critical=True)
            .add_extension(
                x509.AuthorityKeyIdentifier.from_issuer_public_key(authority_key.public_key()), critical=False
            )
            .sign(authority_key, hashes.SHA256())
        )
        self.cert_pem = server.public_bytes(serialization.Encoding.PEM)
        self.key_pem = server_key.private_bytes(
            serialization.Encoding.PEM, serialization.PrivateFormat.PKCS8, serialization.NoEncryption()
        )
        self.cert_path = directory / f"{name}-ca.pem"
        self.cert_path.write_bytes(authority.public_bytes(serialization.Encoding.PEM))
        self.server_cert_path = directory / f"{name}-server.pem"
        self.server_cert_path.write_bytes(self.cert_pem)
        self.key_path = directory / f"{name}-server-key.pem"
        self.key_path.write_bytes(self.key_pem)


class Servers:
    """Servers on loopback ports; every one is stopped by ``close``."""

    def __init__(self) -> None:
        self._http: list[ThreadingHTTPServer] = []
        self._grpc: list[grpc.Server] = []

    def start_http(self, certificate: Certificate | None) -> int:
        class Handler(BaseHTTPRequestHandler):
            def do_GET(self) -> None:
                if self.path.startswith("/tables"):
                    content_type, body = "application/json", b'{"tables":[],"retained":[]}'
                else:
                    data = json.dumps(
                        {
                            "change_id": SENTINEL,
                            "transaction_id": "tx-1",
                            "source_table_id": "t-1",
                            "sequence": 1,
                            "operation": "INSERT",
                            "old_image": None,
                            "new_image": {"id": 1},
                            "schema_version": "v1",
                            "schema": "public",
                            "table": "orders",
                        }
                    )
                    content_type, body = "text/event-stream", f"event: change\ndata: {data}\n\n".encode()
                self.send_response(200)
                self.send_header("Content-Type", content_type)
                self.send_header("Content-Length", str(len(body)))
                self.end_headers()
                self.wfile.write(body)

            def log_message(self, format: str, *args: object) -> None:
                pass

        server = ThreadingHTTPServer(("127.0.0.1", 0), Handler)
        if certificate is not None:
            context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
            context.load_cert_chain(certificate.server_cert_path, certificate.key_path)
            server.socket = context.wrap_socket(server.socket, server_side=True)
        threading.Thread(target=server.serve_forever, daemon=True).start()
        self._http.append(server)
        return server.server_address[1]

    def start_grpc(self, certificate: Certificate | None) -> int:
        class Stream(changestream_pb2_grpc.ChangeStreamServicer):
            def StreamChanges(self, request, context) -> Iterator[changestream_pb2.Change]:
                yield changestream_pb2.Change(change_id=SENTINEL, operation="INSERT", table="orders")

        class Administration(administration_pb2_grpc.AdministrationServicer):
            def ListTables(self, request, context):
                return administration_pb2.ListTablesResponse()

        server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
        changestream_pb2_grpc.add_ChangeStreamServicer_to_server(Stream(), server)
        administration_pb2_grpc.add_AdministrationServicer_to_server(Administration(), server)
        if certificate is None:
            port = server.add_insecure_port("127.0.0.1:0")
        else:
            credentials = grpc.ssl_server_credentials([(certificate.key_pem, certificate.cert_pem)])
            port = server.add_secure_port("127.0.0.1:0", credentials)
        server.start()
        self._grpc.append(server)
        return port

    def start(self, surface: str, certificate: Certificate | None) -> int:
        if surface in ("http", "sse"):
            return self.start_http(certificate)
        return self.start_grpc(certificate)

    def close(self) -> None:
        for http_server in self._http:
            http_server.shutdown()
            http_server.server_close()
        for grpc_server in self._grpc:
            grpc_server.stop(0)


@pytest.fixture
def servers() -> Iterator[Servers]:
    fixture = Servers()
    yield fixture
    fixture.close()


def options(host: str, port: int, scheme: str, anchor: Certificate | Path | None) -> ClientOptions:
    path = anchor.cert_path if isinstance(anchor, Certificate) else anchor
    return ClientOptions(
        address=f"{scheme}://{host}:{port}", api_token="tls-test-token", trust_anchor_file=str(path) if path else None
    )


def use(surface: str, opts: ClientOptions) -> None:
    """Makes one call on the surface and checks that the test server answered."""
    if surface in ("http", "sse"):
        with create_http_client(opts, timeout=httpx.Timeout(10.0, read=None)) as http:
            if surface == "http":
                assert PgChangeFeedHttpClient(http, opts).list_tables("src", "pub").tables == []
            else:
                first = next(iter(PgChangeFeedSseClient(http, opts).stream_changes()))
                assert first.change_id == SENTINEL
        return
    with create_grpc_channel(opts) as channel:
        if surface == "grpc_stream":
            first = next(iter(PgChangeFeedGrpcClient(channel, opts).stream_changes(timeout=10)))
            assert first.change_id == SENTINEL
        else:
            PgChangeFeedAdministrationClient(channel, opts).list_tables(
                administration_pb2.ListTablesRequest(source="src", publication="pub")
            )


def expect_connection_failure(surface: str, opts: ClientOptions) -> None:
    """The surface fails with the connection error of its transport; the server is never reached."""
    if surface in ("http", "sse"):
        with pytest.raises(httpx.ConnectError):
            use(surface, opts)
    elif surface == "grpc_stream":
        with pytest.raises(grpc.RpcError) as error:
            use(surface, opts)
        assert error.value.code() == grpc.StatusCode.UNAVAILABLE
    else:
        with pytest.raises(PgChangeFeedGrpcUnexpectedStatusError) as error:
            use(surface, opts)
        assert error.value.code == grpc.StatusCode.UNAVAILABLE
        assert error.value.message_code is None


@pytest.fixture
def directory(tmp_path: Path) -> Path:
    return tmp_path


@pytest.mark.parametrize("surface", SURFACES)
def test_trust_anchor_connects_over_tls(surface: str, servers: Servers, directory: Path) -> None:
    certificate = Certificate(directory, "server", "localhost")
    port = servers.start(surface, certificate)
    use(surface, options("localhost", port, "https", certificate))


@pytest.mark.parametrize("surface", SURFACES)
def test_issuer_as_trust_anchor_connects_over_tls(surface: str, servers: Servers, directory: Path) -> None:
    issued = IssuedCertificate(directory, "issued", "localhost")
    port = servers.start(surface, issued)
    use(surface, options("localhost", port, "https", issued))


@pytest.mark.parametrize("surface", SURFACES)
def test_issued_server_certificate_as_only_trust_anchor_connects_over_tls(
    surface: str, servers: Servers, directory: Path
) -> None:
    issued = IssuedCertificate(directory, "issued", "localhost")
    port = servers.start(surface, issued)
    use(surface, options("localhost", port, "https", issued.server_cert_path))


@pytest.mark.parametrize("surface", SURFACES)
def test_issued_server_certificate_of_another_server_as_trust_anchor_fails_the_connection(
    surface: str, servers: Servers, directory: Path
) -> None:
    port = servers.start(surface, IssuedCertificate(directory, "issued", "localhost"))
    other = IssuedCertificate(directory, "other", "localhost")
    expect_connection_failure(surface, options("localhost", port, "https", other.server_cert_path))


@pytest.mark.parametrize("surface", SURFACES)
def test_issuer_of_another_server_as_trust_anchor_fails_the_connection(
    surface: str, servers: Servers, directory: Path
) -> None:
    port = servers.start(surface, IssuedCertificate(directory, "issued", "localhost"))
    other = IssuedCertificate(directory, "other", "localhost")
    expect_connection_failure(surface, options("localhost", port, "https", other))


@pytest.mark.parametrize("surface", SURFACES)
def test_anchor_file_with_several_certificates_trusts_the_one_that_matches(
    surface: str, servers: Servers, directory: Path
) -> None:
    certificate = Certificate(directory, "server", "localhost")
    foreign = Certificate(directory, "foreign", "localhost")
    bundle = directory / "bundle.pem"
    bundle.write_bytes(foreign.cert_pem + certificate.cert_pem)
    port = servers.start(surface, certificate)
    use(surface, options("localhost", port, "https", bundle))


@pytest.mark.parametrize("surface", SURFACES)
def test_foreign_anchor_fails_the_connection(surface: str, servers: Servers, directory: Path) -> None:
    port = servers.start(surface, Certificate(directory, "server", "localhost"))
    foreign = Certificate(directory, "foreign", "localhost")
    expect_connection_failure(surface, options("localhost", port, "https", foreign))


@pytest.mark.parametrize("surface", SURFACES)
def test_server_name_not_in_certificate_fails_the_connection(surface: str, servers: Servers, directory: Path) -> None:
    certificate = Certificate(directory, "server", "other.example")
    port = servers.start(surface, certificate)
    # The anchor is the server certificate itself, so the chain is right; the address names `localhost`.
    expect_connection_failure(surface, options("localhost", port, "https", certificate))


@pytest.mark.parametrize("surface", SURFACES)
def test_expired_certificate_fails_the_connection(surface: str, servers: Servers, directory: Path) -> None:
    expired = Certificate(directory, "server", "localhost", expired=True)
    port = servers.start(surface, expired)
    expect_connection_failure(surface, options("localhost", port, "https", expired))


@pytest.mark.parametrize("surface", SURFACES)
def test_without_anchor_the_runtime_refuses_an_unknown_certificate(
    surface: str, servers: Servers, directory: Path
) -> None:
    port = servers.start(surface, Certificate(directory, "server", "localhost"))
    expect_connection_failure(surface, options("localhost", port, "https", None))


@pytest.mark.parametrize("surface", ["http", "sse"])
def test_with_anchor_the_system_trust_does_not_apply(
    surface: str, servers: Servers, directory: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    # The certificate of the server is trusted by the system store of this test
    # process (control: no anchor connects); a foreign anchor must not add it.
    certificate = Certificate(directory, "server", "localhost")
    foreign = Certificate(directory, "foreign", "localhost")
    monkeypatch.setenv("SSL_CERT_FILE", str(certificate.cert_path))
    monkeypatch.delenv("SSL_CERT_DIR", raising=False)
    port = servers.start(surface, certificate)

    use(surface, options("localhost", port, "https", None))
    expect_connection_failure(surface, options("localhost", port, "https", foreign))


@pytest.mark.parametrize("surface", SURFACES)
def test_without_anchor_a_plaintext_address_works_unchanged(surface: str, servers: Servers) -> None:
    port = servers.start(surface, None)
    use(surface, options("localhost", port, "http", None))


@pytest.mark.parametrize("surface", SURFACES)
def test_trust_anchor_with_plaintext_address_is_refused(surface: str, directory: Path) -> None:
    certificate = Certificate(directory, "server", "localhost")
    with pytest.raises(ValueError):
        use(surface, options("localhost", 9, "http", certificate))


def test_grpc_address_without_scheme_goes_over_tls_when_an_anchor_is_set(servers: Servers, directory: Path) -> None:
    certificate = Certificate(directory, "server", "localhost")
    port = servers.start_grpc(certificate)
    opts = ClientOptions(address=f"localhost:{port}", api_token="tls-test-token", trust_anchor_file=str(certificate.cert_path))
    use("grpc_admin", opts)


def test_anchor_option_is_checked_when_the_options_are_created(directory: Path) -> None:
    certificate = Certificate(directory, "server", "localhost")
    assert ClientOptions(address="https://x.invalid", api_token="t").trust_anchor_file is None
    assert ClientOptions(
        address="https://x.invalid", api_token="t", trust_anchor_file=str(certificate.cert_path)
    ).trust_anchor_file == str(certificate.cert_path)

    with pytest.raises(ValueError, match="missing.pem"):
        ClientOptions(address="https://x.invalid", api_token="t", trust_anchor_file=str(directory / "missing.pem"))
    text = directory / "text.txt"
    text.write_text("this is not a certificate\n")
    with pytest.raises(ValueError, match="no PEM certificate"):
        ClientOptions(address="https://x.invalid", api_token="t", trust_anchor_file=str(text))
    empty = directory / "empty.pem"
    empty.write_text("")
    with pytest.raises(ValueError, match="no PEM certificate"):
        ClientOptions(address="https://x.invalid", api_token="t", trust_anchor_file=str(empty))
