package cdcexamples.http

/**
 * Config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse der
 * HTTP-API, das Token der lesenden Rechtsklasse und die zwei Pflichtfelder
 * des aufgerufenen Endpunkts. Adresse und Token kommen aus denselben
 * Umgebungsvariablen, die das Benutzerhandbuch führt (`CDC_HTTP_ADDR`,
 * `CDC_API_TOKEN_READER`), und lassen sich per Flag übersteuern
 * (`ADR-0076` Festlegung 1). Form-Vorbild: `examples/http-client` (Go),
 * `examples/csharp/http-client/Config.cs` (C#).
 */
data class Config(
    val addr: String,
    val token: String,
    val source: String,
    val publication: String,
)
