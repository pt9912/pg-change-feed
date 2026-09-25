package io.github.pt9912.pgchangefeed.http.model

import com.google.gson.annotations.SerializedName

/*
 * Typed request and response data classes of the consumer calls. Every Kotlin
 * property maps to its JSON field name via `@SerializedName`, keeping the
 * property itself idiomatic camelCase.
 */

/** The consumer to register: a unique `consumerId` and a display `name`, both mandatory. */
data class RegisterConsumerRequest(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("name") val name: String,
)

/**
 * The registered consumer. `alreadyRegistered` is true when it existed before;
 * there is no separate status code for that.
 */
data class RegisterConsumerResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("name") val name: String,
    @SerializedName("already_registered") val alreadyRegistered: Boolean,
)

/**
 * The position to store: `offset` is the `commit_position` of the last change
 * the consumer has processed for `sourceId`; all three fields are mandatory. A
 * position older than the current one, or one from a different source, ends
 * `400`.
 *
 * `offset` is a `Long` (signed 64-bit), although the wire field is unsigned
 * 64-bit: Gson's reflection-based codec does not support Kotlin's unsigned
 * value classes (it would (de)serialize the wrapped `Long` field instead of
 * the value). `Long` matches the other 64-bit fields of this SDK (`sequence`,
 * `commit_position`, `version`, `min_age_nanos`); offsets in the top half of
 * the unsigned range are not representable.
 */
data class AcknowledgeConsumerRequest(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("source_id") val sourceId: String,
    @SerializedName("offset") val offset: Long,
)

/** The position now stored for the consumer after the acknowledgement. */
data class AcknowledgeConsumerResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("source_id") val sourceId: String,
    @SerializedName("offset") val offset: Long,
)

/**
 * The stored position of a consumer. `acknowledged = false` reads the defined
 * starting position of a consumer that never acknowledged.
 */
data class ConsumerPositionResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("source_id") val sourceId: String,
    @SerializedName("offset") val offset: Long,
    @SerializedName("acknowledged") val acknowledged: Boolean,
)

/**
 * A consumer that was never registered reports `removed = false`, never
 * `404`.
 */
data class RemoveConsumerResponse(
    @SerializedName("consumer_id") val consumerId: String,
    @SerializedName("removed") val removed: Boolean,
)
