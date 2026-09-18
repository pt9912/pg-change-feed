# Beispiel-Clients

Index der Beispiel-Programme, die den Zugriff auf `pg-change-feed` in drei
Sprachen zeigen — fünf Zugriffsarten × drei Sprachen. Die volle Beschreibung
je Zugriffsart (Voraussetzungen, Umgebungsvariablen, erwartete Ausgabe) steht
in [`docs/user/benutzerhandbuch.md`](../docs/user/benutzerhandbuch.md) §4,
Abschnitte „Zugriff über die HTTP-/JSON-API", „Zugriff über den
gRPC-Change-Stream", „Zugriff über Server-Sent-Events", „Zugriff über das
NATS-Wecksignal" und „Zugriff über den NATS-Vollinhalts-Stream". Diese Datei
dupliziert deren Inhalt nicht, sondern zeigt nur, wo welches Programm liegt
und wie es gebaut/gestartet wird.

**Voraussetzung:** ein laufender Feed-Container gegen eine erreichbare
PostgreSQL-Quelle (Image bauen, Container starten) — siehe
[`docs/user/benutzerhandbuch.md`](../docs/user/benutzerhandbuch.md)
§3 „Erste Schritte", **oder** die Demo-Umgebung unten, die diese
Voraussetzung mit einem Befehl herstellt.

## Demo-Umgebung

`make example-demo-up` (`LH-QA-OPS-001`, `ADR-0098`) fährt PostgreSQL, NATS
und den Feed-Container real hoch und stellt ohne manuellen SQL-Schritt sofort
lesbare Demo-Daten bereit: Schema-Rollout über d-migrate
(`tools/schema/schema.yaml`), eine registrierte Beispiel-Quelle
(`demo-source`) und eine Beispiel-Tabelle (`public.orders`, eine Demo-Zeile),
über `CDC_TABLES` beim Feed-Start automatisch aktiviert. Danach liefert
`GET /changes?source=demo-source` (Token `demo-reader-token`, Adresse
`pg-change-feed:8090` — siehe `examples/.env`) die Demo-Zeile, und jedes der
fünfzehn Beispiel-Programme kann sofort gegen echte Daten laufen (z. B. `make
example-run-go SURFACE=http ARGS="-source demo-source -publication
pub_demo"`).

Eigenständig von [`docs/user/benutzerhandbuch.md`](../docs/user/benutzerhandbuch.md)
§3 „Erste Schritte" (Produktionsanleitung) und von der Wurzel-`compose.yaml`
(CI-/Testtier, `LH-QA-POR-003`) — eigenes Docker-Netzwerk `cdc-examples`,
eigene Container, eigene Daten; beide Compose-Umgebungen können gleichzeitig
laufen. Die Umgebungsdatei `examples/.env` (committet, **keine** Vorlage,
`ADR-0098` Festlegung 3) trägt einen Klartext-Kopfkommentar zur
Netzwerk-Grenze — die Werte gelten ausschließlich innerhalb von
`cdc-examples`, niemals gegen eine Produktionsinstanz.

```bash
make example-demo-up     # hochfahren + bootstrappen
make example-demo-down   # abräumen (Container + Netzwerk)
```

## Go

Teil des Root-Moduls (`go test ./...` bleibt der Anti-Verrottungs-Träger,
`go run ./examples/<name> ...` bleibt technisch funktionsfähig), gebaut/gestartet
über `examples/Dockerfile` (Bau-Kontext Repo-Wurzel, isoliert über
`examples/Dockerfile.dockerignore`) und `make example-run-go`:

| Zugriffsart | Verzeichnis | Start |
|---|---|---|
| [HTTP-/JSON-API](../docs/user/benutzerhandbuch.md#zugriff-über-die-http-json-api) | [`http-client`](http-client) | `make example-run-go SURFACE=http ARGS="-source <quelle> -publication <publication>"` |
| [gRPC-Change-Stream](../docs/user/benutzerhandbuch.md#zugriff-über-den-grpc-change-stream) | [`grpc-client`](grpc-client) | `make example-run-go SURFACE=grpc` |
| [Server-Sent-Events](../docs/user/benutzerhandbuch.md#zugriff-über-server-sent-events) | [`sse-client`](sse-client) | `make example-run-go SURFACE=sse` |
| [NATS-Wecksignal](../docs/user/benutzerhandbuch.md#zugriff-über-das-nats-wecksignal) | [`nats-client`](nats-client) | `make example-run-go SURFACE=nats ARGS="-source <quelle> -schema <schema> -table <tabelle>"` |
| [NATS-Vollinhalts-Stream](../docs/user/benutzerhandbuch.md#zugriff-über-den-nats-vollinhalts-stream) | [`nats-stream-client`](nats-stream-client) | `make example-run-go SURFACE=nats-stream` |

`make example-run-go` baut bei Bedarf `pg-change-feed-examples:go[-<surface>]`
und startet den Container real gegen das Docker-Netzwerk `cdc-examples` mit
`--env-file examples/.env` (Umgebungsdatei-Kontrakt und Netzwerk legt
`slice-beispiele-compose-bootstrap` an).

## C#

Eigene Sprach-Wurzel [`csharp/`](csharp), gebaut über
`make examples-csharp` (Docker, digest-gepinntes Dockerfile, Bau-Kontext
`examples/csharp/`), gestartet über `make example-run-csharp`:

| Zugriffsart | Verzeichnis | Start |
|---|---|---|
| [HTTP-/JSON-API](../docs/user/benutzerhandbuch.md#zugriff-über-die-http-json-api) | [`csharp/http-client`](csharp/http-client) | `make example-run-csharp SURFACE=http ARGS="--source <quelle> --publication <publication>"` |
| [gRPC-Change-Stream](../docs/user/benutzerhandbuch.md#zugriff-über-den-grpc-change-stream) | [`csharp/grpc-client`](csharp/grpc-client) | `make example-run-csharp SURFACE=grpc` |
| [Server-Sent-Events](../docs/user/benutzerhandbuch.md#zugriff-über-server-sent-events) | [`csharp/sse-client`](csharp/sse-client) | `make example-run-csharp SURFACE=sse` |
| [NATS-Wecksignal](../docs/user/benutzerhandbuch.md#zugriff-über-das-nats-wecksignal) | [`csharp/nats-client`](csharp/nats-client) | `make example-run-csharp SURFACE=nats ARGS="--source <quelle> --schema <schema> --table <tabelle>"` |
| [NATS-Vollinhalts-Stream](../docs/user/benutzerhandbuch.md#zugriff-über-den-nats-vollinhalts-stream) | [`csharp/nats-stream-client`](csharp/nats-stream-client) | `make example-run-csharp SURFACE=nats-stream` |

`make example-run-csharp` baut **nicht** — es startet den bereits von
`make examples-csharp` gebauten Image-Tag
`pg-change-feed-examples:csharp[-<surface>]` real gegen das Docker-Netzwerk
`cdc-examples` mit `--env-file examples/.env` (Umgebungsdatei-Kontrakt und
Netzwerk legt `slice-beispiele-compose-bootstrap` an).

## Kotlin

Eigene Sprach-Wurzel [`kotlin/`](kotlin), gebaut über
`make examples-kotlin` (Docker, digest-gepinntes Dockerfile, Bau-Kontext
`examples/kotlin/`), gestartet über `make example-run-kotlin`:

| Zugriffsart | Verzeichnis | Start |
|---|---|---|
| [HTTP-/JSON-API](../docs/user/benutzerhandbuch.md#zugriff-über-die-http-json-api) | [`kotlin/http-client`](kotlin/http-client) | `make example-run-kotlin SURFACE=http ARGS="--source <quelle> --publication <publication>"` |
| [gRPC-Change-Stream](../docs/user/benutzerhandbuch.md#zugriff-über-den-grpc-change-stream) | [`kotlin/grpc-client`](kotlin/grpc-client) | `make example-run-kotlin SURFACE=grpc` |
| [Server-Sent-Events](../docs/user/benutzerhandbuch.md#zugriff-über-server-sent-events) | [`kotlin/sse-client`](kotlin/sse-client) | `make example-run-kotlin SURFACE=sse` |
| [NATS-Wecksignal](../docs/user/benutzerhandbuch.md#zugriff-über-das-nats-wecksignal) | [`kotlin/nats-client`](kotlin/nats-client) | `make example-run-kotlin SURFACE=nats ARGS="--source <quelle> --schema <schema> --table <tabelle>"` |
| [NATS-Vollinhalts-Stream](../docs/user/benutzerhandbuch.md#zugriff-über-den-nats-vollinhalts-stream) | [`kotlin/nats-stream-client`](kotlin/nats-stream-client) | `make example-run-kotlin SURFACE=nats-stream` |

`make example-run-kotlin` baut **nicht** — es startet den bereits von
`make examples-kotlin` gebauten Image-Tag
`pg-change-feed-examples:kotlin[-<surface>]` real gegen das Docker-Netzwerk
`cdc-examples` mit `--env-file examples/.env` (Umgebungsdatei-Kontrakt und
Netzwerk legt `slice-beispiele-compose-bootstrap` an).

## Abgrenzung

Die Wegwerf-Clients unter [`tools/harness/`](../tools/harness) (`httpclient`,
`sseclient`, `grpcclient`) sind **keine** Beispiel-Programme — sie dienen
ausschließlich `make test-integration` als E2E-Testwerkzeug und tauchen im
Benutzerhandbuch nicht auf.
