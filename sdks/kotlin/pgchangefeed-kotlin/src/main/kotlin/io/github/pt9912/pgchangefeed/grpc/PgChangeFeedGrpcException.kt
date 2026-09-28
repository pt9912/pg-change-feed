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
 */
sealed class PgChangeFeedGrpcException(
    val statusCode: Status.Code,
    message: String,
    cause: StatusException,
) : Exception(message, cause)

/**
 * `INVALID_ARGUMENT` — a violated rule of the request (e.g. an empty
 * identifier, an inverted range, a non-positive limit).
 */
class PgChangeFeedGrpcInvalidArgumentException(message: String, cause: StatusException) :
    PgChangeFeedGrpcException(Status.Code.INVALID_ARGUMENT, message, cause)

/**
 * `UNAUTHENTICATED` — a missing bearer token, or one that matches no
 * configured token class.
 */
class PgChangeFeedGrpcUnauthenticatedException(message: String, cause: StatusException) :
    PgChangeFeedGrpcException(Status.Code.UNAUTHENTICATED, message, cause)

/**
 * `PERMISSION_DENIED` — a known token whose class does not reach the called
 * RPC (e.g. a reader token against an admin RPC).
 */
class PgChangeFeedGrpcPermissionDeniedException(message: String, cause: StatusException) :
    PgChangeFeedGrpcException(Status.Code.PERMISSION_DENIED, message, cause)

/**
 * `NOT_FOUND` — the addressed table does not exist in the source database
 * (`enableTable`, `disableTable`, `getTableStatus`).
 */
class PgChangeFeedGrpcNotFoundException(message: String, cause: StatusException) :
    PgChangeFeedGrpcException(Status.Code.NOT_FOUND, message, cause)

/** `INTERNAL` — an unexpected internal error of the server. */
class PgChangeFeedGrpcInternalException(message: String, cause: StatusException) :
    PgChangeFeedGrpcException(Status.Code.INTERNAL, message, cause)

/**
 * Any other non-`OK` gRPC status than `INVALID_ARGUMENT`, `UNAUTHENTICATED`,
 * `PERMISSION_DENIED`, `NOT_FOUND` and `INTERNAL`. It is distinct from the
 * five status-specific exceptions above so a caller can still catch
 * [PgChangeFeedGrpcException] uniformly across every method.
 */
class PgChangeFeedGrpcUnexpectedStatusException(
    statusCode: Status.Code,
    message: String,
    cause: StatusException,
) : PgChangeFeedGrpcException(statusCode, message, cause)
