package io.github.pt9912.pgchangefeed.http.model

import com.google.gson.annotations.SerializedName

/**
 * Typed request/response data classes mirroring the SPEC-018 JSON schemas
 * exactly — field names taken directly from `spec/pflichtenheft.md`
 * SPEC-018, not from the C#/Python sibling packages, which serve only as a
 * structural reference (`ADR-0109` §Kontext, same rule as `ADR-0107` for
 * Python). Every Kotlin property maps to its JSON field name via
 * `@SerializedName`, keeping the property itself idiomatic camelCase.
 */

/** `RegisterConsumer` request — `POST /consumers` (SPEC-018), both fields mandatory. */
data class RegisterConsumerRequest(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("name") val name: String,
)

/**
 * `RegisterConsumer` response (`201`) — `alreadyRegistered` carries
 * idempotency forward, there is no separate status code for it (SPEC-018).
 */
data class RegisterConsumerResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("name") val name: String,
    @SerializedName("already_registered") val alreadyRegistered: Boolean,
)

/**
 * `AcknowledgeConsumer` request — `POST /consumers/acknowledge` (SPEC-018),
 * all three fields mandatory. A position older than the current one, or one
 * from a different source, ends `400`.
 *
 * `offset` is modeled as `Long` (signed 64-bit), not Kotlin's `ULong`, even
 * though the wire field is `uint64` (same as C#'s `ulong`) — Gson's
 * reflection-based codec does not natively support Kotlin's unsigned
 * inline/value classes; it would (de)serialize the type's internal wrapped
 * `Long` field rather than the value itself, silently producing the wrong
 * JSON shape. `Long` matches how the other 64-bit wire fields in this SDK
 * (`sequence`, `commit_position`, `version`, `min_age_nanos`) are already
 * modeled and avoids that Gson footgun, at the cost of not representing
 * offsets in the top half of the `uint64` range — a narrow, named boundary
 * for this pre-1.0 release.
 */
data class AcknowledgeConsumerRequest(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("source_id") val sourceId: String,
    @SerializedName("offset") val offset: Long,
)

/** `AcknowledgeConsumer` response (`200`) — the position now in effect after the acknowledgement. */
data class AcknowledgeConsumerResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("source_id") val sourceId: String,
    @SerializedName("offset") val offset: Long,
)

/**
 * `GetConsumerPosition` response (`200`) — `acknowledged = false` reads the
 * defined starting position without any prior acknowledgement (SPEC-018).
 */
data class ConsumerPositionResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("source_id") val sourceId: String,
    @SerializedName("offset") val offset: Long,
    @SerializedName("acknowledged") val acknowledged: Boolean,
)

/**
 * `RemoveConsumer` response (`200`) — a never-registered consumer reports
 * `removed = false`, never `404` (idempotency instead of an error against
 * an unknown resource, SPEC-018).
 */
data class RemoveConsumerResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("removed") val removed: Boolean,
)
