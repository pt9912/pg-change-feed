package cdcexamples.sse

/**
 * Config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse der
 * HTTP-API und das Token der lesenden Rechtsklasse. Beide kommen aus
 * denselben Umgebungsvariablen, die das Benutzerhandbuch führt
 * (`CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`), und lassen sich per Flag
 * übersteuern (`ADR-0076` Festlegung 1). Form-Vorbild: `examples/sse-client`
 * (Go), `examples/csharp/sse-client/Config.cs` (C#).
 */
data class Config(
    val addr: String,
    val token: String,
)
