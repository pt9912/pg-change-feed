package io.github.pt9912.pgchangefeed

import java.net.URI

/**
 * Shared connection configuration for PG Change Feed client surfaces. The
 * HTTP API (SPEC-018) and the gRPC change stream (SPEC-020) both
 * authenticate with a bearer token against a single server address — this
 * is the one unstrittige, shared configuration denominator identified
 * while writing this project skeleton (ADR-0109 Festlegung 1, referencing
 * ADR-0106 Festlegung 1: "HTTP und gRPC teilen Auth-Header-Form (Bearer
 * Token) und Grundkonfiguration (Adresse, Token)"). Surface-specific
 * behavior (which HTTP paths, which gRPC stub, retry/backoff policy) is
 * deliberately NOT part of this class — it is added by the follow-up
 * slices that build the actual client surfaces.
 *
 * @property address the base address of the PG Change Feed server (HTTP or
 *   gRPC endpoint, depending on the surface that consumes these options).
 * @property apiToken the bearer token sent as an authorization credential
 *   (SPEC-018).
 */
class PgChangeFeedClientOptions(
    val address: URI,
    val apiToken: String,
) {
    init {
        require(apiToken.isNotBlank()) { "API token must not be blank." }
    }
}
