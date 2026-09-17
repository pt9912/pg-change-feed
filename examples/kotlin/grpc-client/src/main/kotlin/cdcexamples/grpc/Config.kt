package cdcexamples.grpc

/**
 * Config trägt die Laufzeit-Eingabe des Beispiels: die Horch-Adresse des
 * gRPC-Streaming-Servers und das Token der lesenden Rechtsklasse. Beide
 * kommen aus denselben Umgebungsvariablen, die das Benutzerhandbuch führt
 * (`CDC_GRPC_ADDR`, `CDC_API_TOKEN_READER`), und lassen sich per Flag
 * übersteuern (`ADR-0076` Festlegung 1). Form-Vorbild:
 * `examples/csharp/grpc-client/Config.cs` (`slice-102`),
 * `examples/grpc-client` (Go).
 */
data class Config(
    val addr: String,
    val token: String,
)
