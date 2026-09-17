namespace CdcExamples.Http;

/// <summary>
/// Config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse der
/// HTTP-API, das Token der lesenden Rechtsklasse und die zwei Pflichtfelder
/// des aufgerufenen Endpunkts. Adresse und Token kommen aus denselben
/// Umgebungsvariablen, die das Benutzerhandbuch führt (<c>CDC_HTTP_ADDR</c>,
/// <c>CDC_API_TOKEN_READER</c>), und lassen sich per Flag übersteuern
/// (<c>ADR-0076</c> Festlegung 1). Form-Vorbild: <c>examples/http-client</c>
/// (Go), <c>config</c> in <c>main.go</c>.
/// </summary>
public sealed record Config(string Addr, string Token, string Source, string Publication);
