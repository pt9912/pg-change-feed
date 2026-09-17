namespace CdcExamples.Nats;

/// <summary>
/// Config trägt die Laufzeit-Eingabe des Beispiels: die NATS-Server-URL des
/// Wecksignals, die Horch-Adresse der HTTP-API, das Token der lesenden
/// Rechtsklasse und die drei Bestandteile des Subjekts. URL, Adresse und
/// Token kommen aus denselben Umgebungsvariablen, die das Benutzerhandbuch
/// führt (<c>CDC_NATS_URL</c>, <c>CDC_HTTP_ADDR</c>,
/// <c>CDC_API_TOKEN_READER</c>), und lassen sich per Flag übersteuern
/// (<c>ADR-0076</c> Festlegung 1). Quelle, Schema und Tabelle sind Flags —
/// das Handbuch führt für eine einzelne Tabelle keine <c>CDC_*</c>-Variable.
/// Form-Vorbild: <c>examples/nats-client</c> (Go), <c>config</c> in
/// <c>main.go</c>.
/// </summary>
public sealed record Config(
    string NatsUrl,
    string Addr,
    string Token,
    string Source,
    string Schema,
    string Table);
