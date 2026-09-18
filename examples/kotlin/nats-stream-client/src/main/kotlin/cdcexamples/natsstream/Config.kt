package cdcexamples.natsstream

/**
 * Config trägt die Laufzeit-Eingabe des Beispiels: die NATS-Server-URL und
 * der Verbindungs-Token des dritten, vollinhaltstragenden Zustellwegs
 * (`ADR-0100`). Beide kommen aus denselben Umgebungsvariablen, die das
 * Benutzerhandbuch führt (`CDC_NATS_URL`, `CDC_NATS_STREAM_TOKEN`), und
 * lassen sich per Flag übersteuern (`ADR-0076` Festlegung 1). Form-Vorbild:
 * `examples/nats-stream-client/main.go` (`config`).
 */
data class Config(
    val natsUrl: String,
    val token: String,
)
