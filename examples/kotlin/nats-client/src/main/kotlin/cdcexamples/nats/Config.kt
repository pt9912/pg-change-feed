package cdcexamples.nats

/**
 * Config trägt die Laufzeit-Eingabe des Beispiels: die NATS-Server-URL des
 * Wecksignals, die Horch-Adresse der HTTP-API, das Token der lesenden
 * Rechtsklasse und die drei Bestandteile des Subjekts. URL, Adresse und
 * Token kommen aus denselben Umgebungsvariablen, die das Benutzerhandbuch
 * führt (`CDC_NATS_URL`, `CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`), und
 * lassen sich per Flag übersteuern (`ADR-0076` Festlegung 1). Quelle, Schema
 * und Tabelle sind Flags — das Handbuch führt für eine einzelne Tabelle
 * keine `CDC_*`-Variable. Form-Vorbild: `examples/nats-client` (Go),
 * `examples/csharp/nats-client/Config.cs` (C#).
 */
data class Config(
    val natsUrl: String,
    val addr: String,
    val token: String,
    val source: String,
    val schema: String,
    val table: String,
)
