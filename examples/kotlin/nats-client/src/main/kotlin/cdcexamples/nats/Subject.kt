package cdcexamples.nats

/**
 * Subject leitet das tabellen-granulare Wecksignal-Subjekt aus Quelle,
 * Schema und Tabelle ab (`SPEC-017`):
 * `cdc.changes.<source_id>.<schema>.<table>`. Das Subjekt wird abgeleitet,
 * nicht handgetippt — dieselbe Form, die der Feed-Container beim
 * Publizieren bildet (`ADR-0056`). Reine Funktion, Form-Vorbild:
 * `examples/nats-client/subject.go`, `examples/csharp/nats-client/Subject.cs`
 * (C#).
 */
object Subject {
    fun build(sourceId: String, schema: String, table: String): String =
        "cdc.changes.$sourceId.$schema.$table"
}
