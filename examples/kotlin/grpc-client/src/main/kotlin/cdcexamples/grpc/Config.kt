package cdcexamples.grpc

/**
 * Config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse des
 * gRPC-Servers, die zwei Token-Klassen und die Felder aller zwölf
 * Fähigkeiten — je Verb prüft [Validator.validate] nur die tatsächlich
 * nötigen Felder. `schema`/`table` dienen doppelt: als optionaler
 * Stream-Filter (`ADR-0133`) und als Tabellen-Identität der
 * Verwaltungs-RPCs; `target` (`--target`) wählt das Zustellziel einer Change im
 * Stream und in `read-changes`, leer ist kein Filter. Form-Vorbild: `examples/csharp/grpc-client/Config.cs`,
 * `config` in `examples/grpc-client/main.go`.
 */
data class Config(
    val addr: String,
    val token: String,
    val adminToken: String,
    val verb: String,
    val schema: String,
    val table: String,
    val consumerId: String,
    val name: String,
    val offset: Long,
    val tableId: String,
    val schemaVersionId: String,
    val version: Long,
    val source: String,
    val publication: String,
    val from: Long,
    val to: Long,
    val limit: Long,
    val minAgeNanos: Long,
    val target: String = "",
)
