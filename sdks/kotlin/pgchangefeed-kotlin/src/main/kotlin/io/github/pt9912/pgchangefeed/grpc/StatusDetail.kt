package io.github.pt9912.pgchangefeed.grpc

import com.google.protobuf.CodedInputStream
import com.google.protobuf.WireFormat
import io.grpc.Metadata
import io.grpc.StatusException
import com.google.protobuf.Any as ProtoAny

/**
 * Reads the message code out of the `grpc-status-details-bin` trailer of a
 * failed call. The trailer carries a serialized `google.rpc.Status`; its
 * repeated field 3 holds `google.protobuf.Any` entries, and an entry whose
 * type is `google.rpc.ErrorInfo` carries the message code as `reason`
 * (field 1) next to its `domain` (field 2). The package depends on no
 * `google.rpc` classes, so this reads exactly these three levels of fields
 * with the Protobuf runtime that is already a dependency.
 */
internal object StatusDetail {
    private val TRAILER_KEY: Metadata.Key<ByteArray> =
        Metadata.Key.of("grpc-status-details-bin", Metadata.BINARY_BYTE_MARSHALLER)
    private const val ERROR_INFO_TYPE_SUFFIX = "google.rpc.ErrorInfo"
    private const val SERVER_DOMAIN = "pg-change-feed"

    private const val STATUS_DETAILS_FIELD = 3
    private const val ERROR_INFO_REASON_FIELD = 1
    private const val ERROR_INFO_DOMAIN_FIELD = 2

    /**
     * Returns the `reason` of the first `ErrorInfo` of domain `pg-change-feed`
     * in the status detail of [exception], or `null`. Never throws.
     */
    fun readMessageCode(exception: StatusException): String? =
        try {
            exception.trailers?.get(TRAILER_KEY)?.let(::readFromStatus)
        } catch (ex: Exception) {
            null
        }

    private fun readFromStatus(status: ByteArray): String? {
        val input = CodedInputStream.newInstance(status)
        while (true) {
            val tag = input.readTag()
            if (tag == 0) {
                return null
            }
            if (WireFormat.getTagFieldNumber(tag) != STATUS_DETAILS_FIELD ||
                WireFormat.getTagWireType(tag) != WireFormat.WIRETYPE_LENGTH_DELIMITED
            ) {
                input.skipField(tag)
                continue
            }
            val code = readFromAny(ProtoAny.parseFrom(input.readBytes()))
            if (code != null) {
                return code
            }
        }
    }

    private fun readFromAny(detail: ProtoAny): String? {
        if (!detail.typeUrl.endsWith(ERROR_INFO_TYPE_SUFFIX)) {
            return null
        }
        val input = detail.value.newCodedInput()
        var reason = ""
        var domain = ""
        while (true) {
            val tag = input.readTag()
            if (tag == 0) {
                break
            }
            val lengthDelimited = WireFormat.getTagWireType(tag) == WireFormat.WIRETYPE_LENGTH_DELIMITED
            when {
                lengthDelimited && WireFormat.getTagFieldNumber(tag) == ERROR_INFO_REASON_FIELD ->
                    reason = input.readStringRequireUtf8()
                lengthDelimited && WireFormat.getTagFieldNumber(tag) == ERROR_INFO_DOMAIN_FIELD ->
                    domain = input.readStringRequireUtf8()
                else -> input.skipField(tag)
            }
        }
        return if (domain == SERVER_DOMAIN && reason.isNotEmpty()) reason else null
    }
}
