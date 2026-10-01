package cdcexamples.http

/**
 * Config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse der
 * HTTP-API, die zwei Token-Klassen (`CDC_API_TOKEN_READER`/
 * `CDC_API_TOKEN_ADMIN`) und die Felder aller zehn Fähigkeiten — je Verb
 * prüft [Validator.validate] nur die tatsächlich nötigen Felder.
 * Form-Vorbild: `examples/http-client` (Go), `config` in `main.go`.
 */
data class Config(
    val addr: String,
    val token: String,
    val adminToken: String,
    val verb: String,
    val source: String,
    val publication: String,
    val consumerId: String,
    val name: String,
    val offset: Long,
    val schema: String,
    val table: String,
    val tableId: String,
    val schemaVersionId: String,
    val version: Long,
    val from: String,
    val to: String,
    val limit: String,
    val minAgeNanos: Long,
    val target: String = "",
)
