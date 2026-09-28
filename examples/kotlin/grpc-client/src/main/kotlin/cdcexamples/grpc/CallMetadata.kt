package cdcexamples.grpc

import io.grpc.Metadata
import java.util.concurrent.TimeUnit

/**
 * CallMetadata baut den `authorization`-Metadata-Eintrag, den der
 * Stream-Aufruf und alle elf Administration-RPCs teilen (`SPEC-020`) —
 * derselbe `Bearer`-Vorsprung wie der Auth-Interceptor des Adapters.
 * [REQUEST_TIMEOUT_SECONDS] begrenzt einen einzelnen unären RPC-Aufruf; der
 * Stream trägt keine eigene Frist. Form-Vorbild:
 * `examples/csharp/grpc-client/CallMetadata.cs`,
 * `authorizationMetadataKey`/`bearerPrefix`/`callCtx` in
 * `examples/grpc-client/main.go`.
 */
object CallMetadata {
    private const val AUTHORIZATION_METADATA_KEY = "authorization"
    private const val BEARER_PREFIX = "Bearer "
    private val AUTHORIZATION_METADATA_ENTRY: Metadata.Key<String> =
        Metadata.Key.of(AUTHORIZATION_METADATA_KEY, Metadata.ASCII_STRING_MARSHALLER)

    const val REQUEST_TIMEOUT_SECONDS = 10L

    fun headers(token: String): Metadata = Metadata().apply { put(AUTHORIZATION_METADATA_ENTRY, BEARER_PREFIX + token) }
}

/**
 * withCallTimeout setzt [CallMetadata.REQUEST_TIMEOUT_SECONDS] als Deadline
 * auf einem Coroutine-Stub — die generische Signatur trifft sowohl den
 * `Administration`- als auch (ungenutzt) den `ChangeStream`-Stub, ohne
 * Wiederholung an jeder Aufrufstelle.
 */
fun <T : io.grpc.kotlin.AbstractCoroutineStub<T>> withCallTimeout(stub: T): T =
    stub.withDeadlineAfter(CallMetadata.REQUEST_TIMEOUT_SECONDS, TimeUnit.SECONDS)
