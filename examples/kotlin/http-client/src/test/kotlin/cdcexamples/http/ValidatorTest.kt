package cdcexamples.http

import kotlin.test.Test
import kotlin.test.assertNull
import kotlin.test.assertTrue

/**
 * Prüft [Validator.validate] — netzlos, reine Funktion. Form-Vorbild:
 * `examples/http-client/main_test.go`.
 */
class ValidatorTest {
    private fun baseConfig(verb: String) = Config(
        addr = "", token = "", adminToken = "", verb = verb,
        source = "", publication = "", consumerId = "", name = "", offset = 0,
        schema = "", table = "", tableId = "", schemaVersionId = "", version = 1,
        from = "", to = "", limit = "", minAgeNanos = 0,
    )

    @Test
    fun validateRejectsUnknownVerb() {
        val err = Validator.validate(baseConfig("unbekannt").copy(addr = "feed:8080"))
        assertTrue(err != null && err.contains("unbekanntes --verb"))
    }

    @Test
    fun validateRequiresAddrBeforeVerbFields() {
        val err = Validator.validate(baseConfig("tables"))
        assertTrue(err != null && err.contains("CDC_HTTP_ADDR"))
    }

    @Test
    fun validateEnableTableAllowsDefaultTableId() {
        val cfg = baseConfig("enable-table").copy(
            addr = "feed:8080", adminToken = "admin-token",
            source = "quelle-1", schema = "public", table = "orders", publication = "pub_quelle_1",
        )
        assertNull(Validator.validate(cfg))
    }

    @Test
    fun validateRetentionRunAllowsZeroMinAge() {
        val cfg = baseConfig("retention-run").copy(addr = "feed:8080", adminToken = "admin-token", source = "quelle-1")
        assertNull(Validator.validate(cfg))
    }

    @Test
    fun validateAcknowledgeRequiresAdminTokenNotReaderToken() {
        val cfg = baseConfig("acknowledge").copy(
            addr = "feed:8080", token = "reader-token", consumerId = "consumer-1", source = "quelle-1",
        )
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("Admin-Token"))
    }

    @Test
    fun validateChangesRequiresOnlySource() {
        val cfg = baseConfig("changes").copy(addr = "feed:8080", token = "reader-token", source = "quelle-1")
        assertNull(Validator.validate(cfg))
    }

    @Test
    fun validateConsumerPositionRequiresReaderTokenNotAdminToken() {
        val cfg = baseConfig("consumer-position").copy(addr = "feed:8080", adminToken = "admin-token", consumerId = "c-1")
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("CDC_API_TOKEN_READER"))
    }
}
