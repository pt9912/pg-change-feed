namespace CdcExamples.Grpc;

/// <summary>
/// Config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse des
/// gRPC-Streaming-Servers und das Token der lesenden Rechtsklasse. Beide
/// kommen aus denselben Umgebungsvariablen, die das Benutzerhandbuch führt
/// (<c>CDC_GRPC_ADDR</c>, <c>CDC_API_TOKEN_READER</c>), und lassen sich per
/// Flag übersteuern (<c>ADR-0076</c> Festlegung 1). Form-Vorbild:
/// <c>examples/grpc-client</c> (Go), <c>config</c> in <c>main.go</c>.
/// </summary>
public sealed record Config(string Addr, string Token);
