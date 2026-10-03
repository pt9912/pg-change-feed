package io.github.pt9912.pgchangefeed.grpc

import io.grpc.Status
import io.grpc.StatusException

/**
 * Typed errors for the non-`OK` gRPC statuses of [PgChangeFeedAdministrationClient]:
 * the gRPC status detail becomes a typed exception instead of a raw status
 * code mixed with the success path, the same design as
 * [io.github.pt9912.pgchangefeed.http.PgChangeFeedException] for the HTTP
 * surface. The original [StatusException] is always the `cause`.
 *
 * `PgChangeFeedGrpcException` is a Kotlin **sealed class**: a caller can
 * `when` over every subtype exhaustively (compiler-checked, no silent
 * fallthrough), while `catch (e: PgChangeFeedGrpcException)` still catches
 * them all. The subtypes group the status codes as follows: `INVALID_ARGUMENT`,
 * `UNAUTHENTICATED`, `PERMISSION_DENIED`, `NOT_FOUND`, `INTERNAL` and any
 * other status.
 *
 * `messageCode` is the message code of the server (`PCF-<E|W|I><4 digits>`),
 * read from the `reason` of the `google.rpc.ErrorInfo` status detail of domain
 * `pg-change-feed` and passed through unchanged, whatever the status; it is
 * `null` when the server sent none (authentication errors, a server without
 * message codes, an unreadable detail).
 */
sealed class PgChangeFeedGrpcException @JvmOverloads constructor(
    val statusCode: Status.Code,
    message: String,
    cause: StatusException,
    val messageCode: String? = null,
) : Exception(message, cause)

/**
 * `INVALID_ARGUMENT` — a violated rule of the request (e.g. an empty
 * identifier, an inverted range, a non-positive limit).
 */
class PgChangeFeedGrpcInvalidArgumentException @JvmOverloads constructor(
    message: String,
    cause: StatusException,
    messageCode: String? = null,
) : PgChangeFeedGrpcException(Status.Code.INVALID_ARGUMENT, message, cause, messageCode)

/**
 * `UNAUTHENTICATED` — a missing bearer token, or one that matches no
 * configured token class.
 */
class PgChangeFeedGrpcUnauthenticatedException @JvmOverloads constructor(
    message: String,
    cause: StatusException,
    messageCode: String? = null,
) : PgChangeFeedGrpcException(Status.Code.UNAUTHENTICATED, message, cause, messageCode)

/**
 * `PERMISSION_DENIED` — a known token whose class does not reach the called
 * RPC (e.g. a reader token against an admin RPC).
 */
class PgChangeFeedGrpcPermissionDeniedException @JvmOverloads constructor(
    message: String,
    cause: StatusException,
    messageCode: String? = null,
) : PgChangeFeedGrpcException(Status.Code.PERMISSION_DENIED, message, cause, messageCode)

/**
 * `NOT_FOUND` — the addressed table does not exist in the source database
 * (`enableTable`, `disableTable`, `getTableStatus`).
 */
class PgChangeFeedGrpcNotFoundException @JvmOverloads constructor(
    message: String,
    cause: StatusException,
    messageCode: String? = null,
) : PgChangeFeedGrpcException(Status.Code.NOT_FOUND, message, cause, messageCode)

/** `INTERNAL` — an unexpected internal error of the server. */
class PgChangeFeedGrpcInternalException @JvmOverloads constructor(
    message: String,
    cause: StatusException,
    messageCode: String? = null,
) : PgChangeFeedGrpcException(Status.Code.INTERNAL, message, cause, messageCode)

/**
 * Any other non-`OK` gRPC status than `INVALID_ARGUMENT`, `UNAUTHENTICATED`,
 * `PERMISSION_DENIED`, `NOT_FOUND` and `INTERNAL`. It is distinct from the
 * five status-specific exceptions above so a caller can still catch
 * [PgChangeFeedGrpcException] uniformly across every method.
 */
class PgChangeFeedGrpcUnexpectedStatusException @JvmOverloads constructor(
    statusCode: Status.Code,
    message: String,
    cause: StatusException,
    messageCode: String? = null,
) : PgChangeFeedGrpcException(statusCode, message, cause, messageCode)
