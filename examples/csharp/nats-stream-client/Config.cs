namespace CdcExamples.NatsStream;

/// <summary>
/// Config trägt die Laufzeit-Eingabe des Beispiels: die NATS-Server-URL und
/// der Verbindungs-Token des dritten, vollinhaltstragenden Zustellwegs
/// (<c>ADR-0100</c>). Beide kommen aus denselben Umgebungsvariablen, die das
/// Benutzerhandbuch führt (<c>CDC_NATS_URL</c>, <c>CDC_NATS_STREAM_TOKEN</c>),
/// und lassen sich per Flag übersteuern (<c>ADR-0076</c> Festlegung 1).
/// Form-Vorbild: <c>examples/nats-stream-client/main.go</c> (<c>config</c>).
/// </summary>
public sealed record Config(string NatsUrl, string Token);
