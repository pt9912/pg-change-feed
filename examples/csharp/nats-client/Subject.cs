namespace CdcExamples.Nats;

/// <summary>
/// Subject leitet das tabellen-granulare Wecksignal-Subjekt aus Quelle,
/// Schema und Tabelle ab (<c>SPEC-017</c>):
/// <c>cdc.changes.&lt;source_id&gt;.&lt;schema&gt;.&lt;table&gt;</c>. Das
/// Subjekt wird abgeleitet, nicht handgetippt — dieselbe Form, die der
/// Feed-Container beim Publizieren bildet (<c>ADR-0056</c>). Reine Funktion,
/// Form-Vorbild: <c>examples/nats-client/subject.go</c>.
/// </summary>
public static class Subject
{
    public static string Build(string sourceId, string schema, string table) =>
        $"cdc.changes.{sourceId}.{schema}.{table}";
}
