package io.github.pt9912.pgchangefeed.integration

/**
 * Shared environment plumbing for the real-server integration phases: every
 * phase names its test class explicitly through the
 * runner's `PGCHANGEFEED_TEST_NAME` `--tests` filter (no silent exclusion),
 * and every marker (`READY`/`RECEIVED`/`REJECTED`) is what
 * `tools/harness/run-sdk-kotlin-integration-tests.sh` asserts on via
 * `docker logs` while the test still runs — the Gradle task streams the
 * standard output live (`showStandardStreams`), so every print reaches the
 * runner in time.
 */
object PhaseEnvironment {
    private fun required(name: String): String =
        System.getenv(name)
            ?: throw IllegalStateException("Umgebungsvariable $name fehlt — die Runner-Phase trägt sie.")

    val table: String get() = required("PGCHANGEFEED_E2E_TABLE")
    val sentinel: String get() = required("PGCHANGEFEED_E2E_SENTINEL")
    val apiToken: String get() = required("PGCHANGEFEED_API_TOKEN")
    val adminToken: String get() = required("PGCHANGEFEED_API_TOKEN_ADMIN")
    val readerToken: String get() = required("PGCHANGEFEED_API_TOKEN_READER")
    val grpcAddr: String get() = required("PGCHANGEFEED_GRPC_ADDR")
    val httpAddr: String get() = required("PGCHANGEFEED_HTTP_ADDR")
    val natsUrl: String get() = required("PGCHANGEFEED_NATS_URL")
    val natsStreamToken: String get() = required("PGCHANGEFEED_NATS_STREAM_TOKEN")
    val sourceId: String get() = required("PGCHANGEFEED_SOURCE_ID")
    val httpPublication: String get() = required("PGCHANGEFEED_HTTP_PUBLICATION")
}