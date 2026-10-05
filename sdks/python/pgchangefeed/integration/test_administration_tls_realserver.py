"""Real-server TLS phase for the gRPC administration client.

With the trust anchor of the options the SDK registers a disposable consumer
with the admin token and lists tables with the reader token over TLS (the
runner reads the registration back from `cdc.consumer`); without an anchor,
with a foreign anchor and with a server name outside the certificate the same
call fails the TLS check.
"""

from __future__ import annotations

import os
import uuid

from tls_scenarios import CA_FILE, expect_refusals

from pgchangefeed import ClientOptions, PgChangeFeedAdministrationClient, create_grpc_channel
from pgchangefeed.grpc_gen import administration_pb2

_ADDR = os.environ["PGCHANGEFEED_GRPC_ADDR"]
_ADMIN_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN_ADMIN"]
_READER_TOKEN = os.environ["PGCHANGEFEED_API_TOKEN_READER"]
_SOURCE = os.environ["PGCHANGEFEED_SOURCE_ID"]
_PUBLICATION = os.environ["PGCHANGEFEED_HTTP_PUBLICATION"]


def test_realserver_registers_a_consumer_and_lists_tables_over_tls() -> None:
    print("READY", flush=True)

    admin_options = ClientOptions(address=f"https://{_ADDR}", api_token=_ADMIN_TOKEN, trust_anchor_file=CA_FILE)
    consumer_id = f"python-sdk-tls-admin-{uuid.uuid4().hex[:12]}"
    with create_grpc_channel(admin_options) as channel:
        registered = PgChangeFeedAdministrationClient(channel, admin_options).register_consumer(
            administration_pb2.RegisterConsumerRequest(consumer_id=consumer_id, name=f"Python SDK TLS {consumer_id}")
        )
    assert registered.consumer_id == consumer_id

    reader_options = ClientOptions(address=f"https://{_ADDR}", api_token=_READER_TOKEN, trust_anchor_file=CA_FILE)
    with create_grpc_channel(reader_options) as channel:
        tables = PgChangeFeedAdministrationClient(channel, reader_options).list_tables(
            administration_pb2.ListTablesRequest(source=_SOURCE, publication=_PUBLICATION)
        )

    print(f"RECEIVED consumer_id={registered.consumer_id} tables={len(tables.tables)}", flush=True)


def test_realserver_without_anchor_foreign_anchor_and_wrong_name_fail_the_tls_check() -> None:
    def call(options: ClientOptions) -> None:
        with create_grpc_channel(options) as channel:
            PgChangeFeedAdministrationClient(channel, options).list_tables(
                administration_pb2.ListTablesRequest(source=_SOURCE, publication=_PUBLICATION)
            )

    expect_refusals(f"https://{_ADDR}", _READER_TOKEN, call)
