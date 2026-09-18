# Beispiel-Clients

Index der Beispiel-Programme, die den Zugriff auf `pg-change-feed` in drei
Sprachen zeigen — vier Zugriffsarten × drei Sprachen. Die volle Beschreibung
je Zugriffsart (Voraussetzungen, Umgebungsvariablen, erwartete Ausgabe) steht
in [`docs/user/benutzerhandbuch.md`](../docs/user/benutzerhandbuch.md) §4,
Abschnitte „Zugriff über die HTTP-/JSON-API", „Zugriff über den
gRPC-Change-Stream", „Zugriff über Server-Sent-Events" und „Zugriff über das
NATS-Wecksignal". Diese Datei dupliziert deren Inhalt nicht, sondern zeigt
nur, wo welches Programm liegt und wie es gebaut/gestartet wird.

## Go

Teil des Root-Moduls, kein eigener Bau-Schritt nötig:

| Zugriffsart | Verzeichnis | Start |
|---|---|---|
| [HTTP-/JSON-API](../docs/user/benutzerhandbuch.md#zugriff-über-die-http-json-api) | [`http-client`](http-client) | `go run ./examples/http-client ...` |
| [gRPC-Change-Stream](../docs/user/benutzerhandbuch.md#zugriff-über-den-grpc-change-stream) | [`grpc-client`](grpc-client) | `go run ./examples/grpc-client ...` |
| [Server-Sent-Events](../docs/user/benutzerhandbuch.md#zugriff-über-server-sent-events) | [`sse-client`](sse-client) | `go run ./examples/sse-client ...` |
| [NATS-Wecksignal](../docs/user/benutzerhandbuch.md#zugriff-über-das-nats-wecksignal) | [`nats-client`](nats-client) | `go run ./examples/nats-client ...` |

## C#

Eigene Sprach-Wurzel [`csharp/`](csharp), gebaut über
`make examples-csharp` (Docker, digest-gepinntes Dockerfile, Bau-Kontext
`examples/csharp/`):

| Zugriffsart | Verzeichnis | Image |
|---|---|---|
| [HTTP-/JSON-API](../docs/user/benutzerhandbuch.md#zugriff-über-die-http-json-api) | [`csharp/http-client`](csharp/http-client) | `pg-change-feed-examples:csharp` |
| [gRPC-Change-Stream](../docs/user/benutzerhandbuch.md#zugriff-über-den-grpc-change-stream) | [`csharp/grpc-client`](csharp/grpc-client) | `pg-change-feed-examples:csharp-grpc` |
| [Server-Sent-Events](../docs/user/benutzerhandbuch.md#zugriff-über-server-sent-events) | [`csharp/sse-client`](csharp/sse-client) | `pg-change-feed-examples:csharp-sse` |
| [NATS-Wecksignal](../docs/user/benutzerhandbuch.md#zugriff-über-das-nats-wecksignal) | [`csharp/nats-client`](csharp/nats-client) | `pg-change-feed-examples:csharp-nats` |

Start je Image per `docker run --rm <ENV-Variablen> <Image> <Argumente>` —
die konkrete ENV-/Argument-Form je Zugriffsart steht in der oben verlinkten
Handbuch-Sektion.

## Kotlin

Eigene Sprach-Wurzel [`kotlin/`](kotlin), gebaut über
`make examples-kotlin` (Docker, digest-gepinntes Dockerfile, Bau-Kontext
`examples/kotlin/`):

| Zugriffsart | Verzeichnis | Image |
|---|---|---|
| [HTTP-/JSON-API](../docs/user/benutzerhandbuch.md#zugriff-über-die-http-json-api) | [`kotlin/http-client`](kotlin/http-client) | `pg-change-feed-examples:kotlin` |
| [gRPC-Change-Stream](../docs/user/benutzerhandbuch.md#zugriff-über-den-grpc-change-stream) | [`kotlin/grpc-client`](kotlin/grpc-client) | `pg-change-feed-examples:kotlin-grpc` |
| [Server-Sent-Events](../docs/user/benutzerhandbuch.md#zugriff-über-server-sent-events) | [`kotlin/sse-client`](kotlin/sse-client) | `pg-change-feed-examples:kotlin-sse` |
| [NATS-Wecksignal](../docs/user/benutzerhandbuch.md#zugriff-über-das-nats-wecksignal) | [`kotlin/nats-client`](kotlin/nats-client) | `pg-change-feed-examples:kotlin-nats` |

Start je Image per `docker run --rm <ENV-Variablen> <Image> <Argumente>` —
die konkrete ENV-/Argument-Form je Zugriffsart steht in der oben verlinkten
Handbuch-Sektion.

## Abgrenzung

Die Wegwerf-Clients unter [`tools/harness/`](../tools/harness) (`httpclient`,
`sseclient`, `grpcclient`) sind **keine** Beispiel-Programme — sie dienen
ausschließlich `make test-integration` als E2E-Testwerkzeug und tauchen im
Benutzerhandbuch nicht auf.
