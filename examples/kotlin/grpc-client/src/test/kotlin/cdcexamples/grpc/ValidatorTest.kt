package cdcexamples.grpc

import kotlin.test.Test
import kotlin.test.assertNull
import kotlin.test.assertTrue

/**
 * Prüft [Validator.validate] — netzlos, reine Funktion. Deckt für jedes der
 * zwölf Verben die Rechtsklassen-Bindung (das jeweils falsche Token allein
 * genügt nicht) und die dokumentierten Pflichtfeld-/Default-Ausnahmen.
 * Form-Vorbild: `examples/csharp/grpc-client/GrpcClient.Tests/ValidatorTests.cs`,
 * `examples/grpc-client/main_test.go`.
 */
class ValidatorTest {
    private fun baseConfig(verb: String) = Config(
        addr = "", token = "", adminToken = "", verb = verb,
        schema = "", table = "", consumerId = "", name = "", offset = 0,
        tableId = "", schemaVersionId = "", version = 1,
        source = "", publication = "", from = 0, to = 0, limit = 0, minAgeNanos = 0,
    )

    @Test
    fun validateRejectsUnknownVerb() {
        val err = Validator.validate(baseConfig("unbekannt").copy(addr = "feed:9090"))
        assertTrue(err != null && err.contains("unbekanntes --verb"))
    }

    @Test
    fun validateRequiresAddrBeforeVerbFields() {
        val err = Validator.validate(baseConfig("list-tables"))
        assertTrue(err != null && err.contains("CDC_GRPC_ADDR"))
    }

    @Test
    fun validateStreamDefaultAcceptsReaderTokenWithoutFilter() {
        val cfg = baseConfig("stream").copy(addr = "feed:9090", token = "reader-token")
        assertNull(Validator.validate(cfg))
    }

    @Test
    fun validateStreamRequiresReaderTokenNotAdminToken() {
        val cfg = baseConfig("stream").copy(addr = "feed:9090", adminToken = "admin-token")
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("CDC_API_TOKEN_READER"))
    }

    @Test
    fun validateEnableTableAllowsDefaultTableId() {
        val cfg = baseConfig("enable-table").copy(
            addr = "feed:9090", adminToken = "admin-token",
            source = "quelle-1", schema = "public", table = "orders", publication = "pub_quelle_1",
        )
        assertNull(Validator.validate(cfg))
    }

    @Test
    fun validateRunRetentionAllowsZeroMinAge() {
        val cfg = baseConfig("run-retention").copy(addr = "feed:9090", adminToken = "admin-token", source = "quelle-1")
        assertNull(Validator.validate(cfg))
    }

    @Test
    fun validateAcknowledgeConsumerRequiresAdminTokenNotReaderToken() {
        val cfg = baseConfig("acknowledge-consumer").copy(
            addr = "feed:9090", token = "reader-token", consumerId = "consumer-1", source = "quelle-1",
        )
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("Admin-Token"))
    }

    @Test
    fun validateRegisterConsumerRequiresAdminTokenNotReaderToken() {
        val cfg = baseConfig("register-consumer").copy(
            addr = "feed:9090", token = "reader-token", consumerId = "consumer-1", name = "Consumer",
        )
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("Admin-Token"))
    }

    @Test
    fun validateRemoveConsumerRequiresAdminTokenNotReaderToken() {
        val cfg = baseConfig("remove-consumer").copy(addr = "feed:9090", token = "reader-token", consumerId = "consumer-1")
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("Admin-Token"))
    }

    @Test
    fun validateGetConsumerPositionRequiresReaderTokenNotAdminToken() {
        val cfg = baseConfig("get-consumer-position").copy(addr = "feed:9090", adminToken = "admin-token", consumerId = "consumer-1")
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("CDC_API_TOKEN_READER"))
    }

    @Test
    fun validateEnableTableRequiresAdminTokenNotReaderToken() {
        val cfg = baseConfig("enable-table").copy(
            addr = "feed:9090", token = "reader-token",
            source = "quelle-1", schema = "public", table = "orders", publication = "pub_quelle_1",
        )
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("Admin-Token"))
    }

    @Test
    fun validateDisableTableRequiresAdminTokenNotReaderToken() {
        val cfg = baseConfig("disable-table").copy(
            addr = "feed:9090", token = "reader-token",
            source = "quelle-1", schema = "public", table = "orders", publication = "pub_quelle_1",
        )
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("Admin-Token"))
    }

    @Test
    fun validateGetTableStatusRequiresReaderTokenNotAdminToken() {
        val cfg = baseConfig("get-table-status").copy(
            addr = "feed:9090", adminToken = "admin-token",
            source = "quelle-1", schema = "public", table = "orders", publication = "pub_quelle_1",
        )
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("CDC_API_TOKEN_READER"))
    }

    @Test
    fun validateListTablesRequiresReaderTokenNotAdminToken() {
        val cfg = baseConfig("list-tables").copy(addr = "feed:9090", adminToken = "admin-token", source = "quelle-1", publication = "pub_quelle_1")
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("CDC_API_TOKEN_READER"))
    }

    @Test
    fun validateRunRetentionRequiresAdminTokenNotReaderToken() {
        val cfg = baseConfig("run-retention").copy(addr = "feed:9090", token = "reader-token", source = "quelle-1")
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("Admin-Token"))
    }

    @Test
    fun validateDiagnoseRequiresReaderTokenNotAdminToken() {
        val cfg = baseConfig("diagnose").copy(addr = "feed:9090", adminToken = "admin-token", source = "quelle-1")
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("CDC_API_TOKEN_READER"))
    }

    @Test
    fun validateReadChangesRequiresReaderTokenNotAdminToken() {
        val cfg = baseConfig("read-changes").copy(addr = "feed:9090", adminToken = "admin-token", source = "quelle-1")
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("CDC_API_TOKEN_READER"))
    }

    @Test
    fun validateDiagnoseRequiresSource() {
        val cfg = baseConfig("diagnose").copy(addr = "feed:9090", token = "reader-token")
        val err = Validator.validate(cfg)
        assertTrue(err != null && err.contains("-source ist Pflicht"))
    }
}
