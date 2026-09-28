# ADR-0130: gRPC-Verwaltungs-API — neun RPCs im neuen `Administration`-Service (Supersedes ADR-0057, teilweise)

**Status:** Accepted — Supersedes [`ADR-0057`](0057-http-grpc-api.md) in genau
einer Stelle: Teilfrage 1 („Protokoll"), dort Option D („HTTP/JSON,
ausschließlich") wird durch Option C („beide parallel, HTTP + gRPC") ersetzt —
für die Verwaltungsfähigkeiten, nicht für `ADR-0060`s Live-Streaming-Fähigkeit,
die davon unberührt bleibt. Alles Übrige von `ADR-0057` bleibt in Kraft,
insbesondere Teilfrage 2 (Fähigkeits-Umfang: dieselben neun Port-gedeckten
Fähigkeiten, jetzt zusätzlich über gRPC statt nur über HTTP), Teilfrage 3
(Token-Klassen `CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`, unverändert
wiederverwendet) und Teilfrage 4 (HTTP-Adapter-Platzierung, unverändert). Der
in `ADR-0057`s Re-Evaluierungs-Trigger vorgezeichnete Fall („ein konkreter
gRPC-Consumer benennt sich") ist mit dieser ADR zum zweiten Mal eingetreten —
das erste Mal löste `ADR-0060`s disjunkte Fähigkeit (Live-Streaming) aus,
dieses Mal die Verwaltungsfähigkeiten selbst.

**Datum:** 2026-09-28

**Autor:** Architect-Agent (Modul 8), Architect-Zug nach explizitem
Auftraggeber-Auftrag (2026-09-28): gRPC soll dieselben Verwaltungsfähigkeiten
wie die HTTP-/JSON-API bekommen. Jede Aussage über eine Menge nennt die
Menge, an der sie geprüft ist (`AGENTS.md` §3.12); vom Architect selbst
gefahren wurde nichts — diese ADR entscheidet die Umsetzungsform vor der
Implementierung, ihre Fitness-Function-Zeilen sind Erwartungen an den
umsetzenden Slice, keine Messungen.

**Bezug:** [`LH-FA-SST-006`](../../../spec/lastenheft.md) (Haupt-Bezug —
Out-of-Scope nennt die Protokoll-Wahl ausdrücklich als Architektur-Frage),
[`LH-FA-CON-001`](../../../spec/lastenheft.md),
[`LH-FA-CON-003`](../../../spec/lastenheft.md),
[`LH-FA-CON-004`](../../../spec/lastenheft.md),
[`LH-FA-CON-005`](../../../spec/lastenheft.md),
[`LH-FA-CON-006`](../../../spec/lastenheft.md),
[`LH-FA-CFG-001`](../../../spec/lastenheft.md),
[`LH-FA-CFG-002`](../../../spec/lastenheft.md),
[`LH-FA-CFG-003`](../../../spec/lastenheft.md),
[`LH-FA-CFG-004`](../../../spec/lastenheft.md),
[`LH-FA-RET-002`](../../../spec/lastenheft.md)…[`LH-FA-RET-004`](../../../spec/lastenheft.md)
(die neun Fähigkeiten, unverändert dieselben Use Cases wie beim HTTP-Adapter),
[`ADR-0057`](0057-http-grpc-api.md) (teilweise superseded — Haupt-Bezug,
Token-Modell und Fähigkeits-Umfang bleiben in Kraft),
[`ADR-0060`](0060-grpc-streaming-mechanismus.md) (Bezug, unberührt — disjunkte
Fähigkeit, eigener Service, eigenes `.proto`-Paket),
[`ADR-0047`](0047-rollenspezifische-dsn-verdrahtung.md) (Rollen-Modell, das
die Token-Klassen hier unverändert spiegeln), [`ADR-0028`](0028-inbound-use-cases.md)
(Inbound Use Cases — dieselben, die HTTP bereits aufruft, kein zweiter
Domänenpfad), [`ADR-0034`](0034-ports-nach-faehigkeiten.md) (Ports nach
Fähigkeit — Vorbild für die Paket-Trennung `ChangeStream`/`Administration`)

**Schärft:** [`ARC-005`](../../../spec/architecture.md) (macht das dort
bereits generisch genannte „später HTTP-/gRPC" für die
Verwaltungsfähigkeiten ein zweites Mal konkret: gRPC zusätzlich zu HTTP/JSON,
unabhängig von `ADR-0060`s Konkretisierung für Live-Streaming)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0057`](0057-http-grpc-api.md) wählte für die Verwaltungs-API
ausschließlich HTTP/JSON und verwarf gRPC (Option B) explizit mit der
Begründung „kein bekannter gRPC-Consumer benennt diesen Bedarf" — mit einem
eigenen Re-Evaluierungs-Trigger für genau diesen Fall. Dieser Trigger ist
bereits einmal eingetreten: [`ADR-0060`](0060-grpc-streaming-mechanismus.md)
führte gRPC für eine **andere**, disjunkte Fähigkeit ein (Live-Streaming
vollständiger Change-Inhalte, `LH-FA-SST-008`) und band das ausdrücklich
nicht an `ADR-0057`s Verwaltungs-API zurück. Der Auftraggeber (pt9912) hat am
2026-09-28 den Trigger ein zweites Mal ausgelöst — diesmal für die
Verwaltungsfähigkeiten selbst: gRPC soll dieselben neun Fähigkeiten tragen,
die heute ausschließlich über HTTP/JSON erreichbar sind
(Consumer-Registrierung/-Bestätigung/-Position/-Entfernung,
Tabellen-Aktivierung/-Deaktivierung/-Status/-Liste, Retention-Lauf). Das ist
keine Korrektur von `ADR-0057`s damaliger Begründung — die Prämisse „kein
gRPC-Consumer" galt zum Entscheidungszeitpunkt real —, sondern eine neue
Tatsache, die eine neue Entscheidung über Teilfrage 1 verlangt (Baseline-
Regelwerk `modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz,
Verdikt 2: Folge-ADR statt stiller Korrektur).

**Out-of-Scope-Klarstellung.** `LH-FA-REA-*`/`SPEC-022` (`GET /changes`,
Bestandslesen über einen Bereich) ist **nicht** Gegenstand dieser ADR: Für
gRPC existiert für neue Changes bereits ein Vollinhalts-Livepfad
(`StreamChanges`, `ADR-0060`) — ein zusätzlicher unärer Range-Read-RPC über
gRPC wäre eine eigene, hier nicht beauftragte Fähigkeit. Der Auftrag nennt
ausdrücklich die neun Verwaltungsfähigkeiten der HTTP-API, nicht das Lesen.

Drei Konstraints prägen den Lösungsraum, alle bereits von `ADR-0057`/`ADR-0060`
gesetzt:

- **Dieselben Inbound Use Cases, kein zweiter Domänenpfad.** `LH-FA-SST-006`s
  Boundary-Kriterium verlangt fachliche Gleichwertigkeit über alle
  Zugriffswege. Der neue gRPC-Adapter ruft exakt dieselben
  `internal/application/port/inbound/*`-Interfaces auf, die der HTTP-Adapter
  bereits verwendet (`consumer.go`, `verwaltung.go`, `retention.go`) — keine
  neue Anwendungsschicht-Entscheidung.
- **Bestehender gRPC-Adapter/-Port.** `internal/adapters/driving/grpc/`
  trägt bereits einen `grpc.Server` (`ADR-0060`) samt Auth-Infrastruktur
  (`role`, `classifyToken`, `credentialToken` in `interceptor.go`) und die
  additive Aktivierung über `CDC_GRPC_ADDR` — diese Bausteine sind
  wiederverwendbar, keine neue Secret-Klasse oder ENV-Variable nötig.
- **`.a-check.yml`s bestehende Globs sind bereits rekursiv** (`adapters:
  ["internal/adapters/**"]`, `contract: ["gen/**"]`, geprüft durch Lesen der
  Datei) — jeder neue Pfad unter diesen beiden Wurzeln braucht keine
  Erweiterung der Architekturregeln.

## Entscheidung

Wir erweitern die gRPC-Fläche um einen neuen, eigenständigen Service
`Administration` mit neun unären RPCs, im selben Driving-Adapter-Paket wie
der bestehende Stream-Server, auf demselben Port, mit denselben zwei
Token-Klassen. Sechs Teilfragen:

### Teilfrage 1 — Ein Service oder zwei? Ein `.proto`-Paket oder zwei?

| Option | Pro | Contra |
|---|---|---|
| A — die neun RPCs in den bestehenden `ChangeStream`-Service/das bestehende Paket `cdc.stream.v1` aufnehmen | ein Client-Stub für die gesamte gRPC-Fläche; keine neue Datei | vermischt zwei von `ADR-0060` bewusst getrennte Konzerne (reines Live-Streaming vs. administrative Verwaltung) unter einem Namen, der laut `ADR-0060`s eigenem Wortlaut exakt auf Live-Streaming beschränkt ist; jede künftige Verwaltungs-Erweiterung wüchse in dieselbe `.proto`-Datei wie der Stream-Vertrag; ein Client, der nur streamen will, generiert/importiert unnötig die volle Verwaltungsfläche |
| **B — neuer Service `Administration` in einem neuen Paket `cdc.administration.v1`, eigene `.proto`-Datei `proto/cdc/administration/v1/administration.proto` (gewählt)** | `changestream.proto` bleibt **byte-identisch unverändert** — null Diff-Risiko für bestehende `StreamChanges`-Clients, stärker als bloße Rückwärtskompatibilität; folgt `ADR-0034`s Prinzip „Ports/Verträge nach Fähigkeit" — Live-Streaming und Verwaltung sind zwei Konzerne mit unterschiedlicher Zustellsemantik (Stream vs. Request/Response); ein reiner Streaming-Client importiert weiterhin nur `gen/cdc/stream/v1` | zwei generierte Go-Paket-Wurzeln (`gen/cdc/stream/v1`, `gen/cdc/administration/v1`) statt einer; ein zweiter `protoc`-Lauf in derselben Build-Stufe (mechanisch, kein neuer Build-Schritt) |
| C — neuer Service `Administration`, eigene `.proto`-Datei, aber im **bestehenden** Paket `cdc.stream.v1` | ein Go-Import-Pfad (`gen/cdc/stream/v1`) deckt beide Services | bindet die Versionierung zweier fachlich unabhängiger Konzerne aneinander — eine künftige Breaking-Change-Version der Verwaltungsfläche zwänge entweder eine paketweite Versionsanhebung, die den unberührten Stream-Vertrag mitzöge, oder ohnehin ein zweites Paket; derselbe Einwand wie bei `ADR-0034`s Ports-nach-Fähigkeit-Prinzip |

`ChangeStream`/`changestream.proto` (`ADR-0060`) bleiben durch diese
Entscheidung **unverändert** — kein Zeilen-Diff an dieser bestehenden Datei
oder ihrem generierten Code ist Folge dieser ADR.

### Teilfrage 2 — RPC-Form je Fähigkeit

Alle neun Fähigkeiten sind unäre Request/Response-RPCs — dieselbe Begründung,
mit der `ADR-0057` Teilfrage 1 gRPC für die Verwaltungs-API ursprünglich als
unverhältnismäßig verwarf (kein Streaming-Vorteil bei Request/Response-
Aufrufen), jetzt als bewusste Formwahl statt als Totalablehnung des
Protokolls:

```proto
service Administration {
  rpc RegisterConsumer(RegisterConsumerRequest) returns (RegisterConsumerResponse);
  rpc AcknowledgeConsumer(AcknowledgeConsumerRequest) returns (AcknowledgeConsumerResponse);
  rpc GetConsumerPosition(GetConsumerPositionRequest) returns (GetConsumerPositionResponse);
  rpc RemoveConsumer(RemoveConsumerRequest) returns (RemoveConsumerResponse);
  rpc EnableTable(EnableTableRequest) returns (EnableTableResponse);
  rpc DisableTable(DisableTableRequest) returns (DisableTableResponse);
  rpc GetTableStatus(GetTableStatusRequest) returns (GetTableStatusResponse);
  rpc ListTables(ListTablesRequest) returns (ListTablesResponse);
  rpc RunRetention(RunRetentionRequest) returns (RunRetentionResponse);
}
```

### Teilfrage 3 — Nachrichtenschema je RPC

Feldnamen/Typen 1:1 aus dem bestehenden HTTP-Kontrakt (`SPEC-018`,
`internal/adapters/driving/http/{consumer,registerconsumer,verwaltung,retention}.go`)
übernommen; `snake_case` in proto3 wird zu `camelCase` im generierten
Go-Code — Formsache der Sprache, keine inhaltliche Abweichung.

| RPC | Request-Felder | Response-Felder |
|---|---|---|
| `RegisterConsumer` | `string consumer_id`, `string name` | `string consumer_id`, `string name`, `bool already_registered` |
| `AcknowledgeConsumer` | `string consumer_id`, `string source_id`, `uint64 offset` | `string consumer_id`, `string source_id`, `uint64 offset` |
| `GetConsumerPosition` | `string consumer_id` | `string consumer_id`, `string source_id`, `uint64 offset`, `bool acknowledged` |
| `RemoveConsumer` | `string consumer_id` | `string consumer_id`, `bool removed` |
| `EnableTable` | `string source`, `string schema`, `string table`, `string table_id`, `string schema_version_id`, `int64 version`, `string publication` | `string table_id`, `string source`, `string schema`, `string table`, `bool already_enabled` |
| `DisableTable` | `string source`, `string schema`, `string table`, `string publication` | `bool removed`, `bool retained` |
| `GetTableStatus` | `string source`, `string schema`, `string table`, `string publication` | `bool enabled`, `bool retained` |
| `ListTables` | `string source`, `string publication` | `repeated SourceTable tables`, `repeated SourceTable retained` |
| `RunRetention` | `string source`, `int64 min_age_nanos` | `int64 deleted` |

`SourceTable` (Hilfsnachricht für `ListTables`): `string table_id`,
`string source`, `string schema`, `string table` — dieselben vier Felder wie
`sourceTableResponse` im HTTP-Adapter.

Kein Feld trägt einen JSON-Decode-Fehlerfall (`400` bei „kein gültiges
JSON" in `verwaltung.go`/`consumer.go`): Das Nachrichtenschema ist bereits
typisiert (Protobuf), dieser Fehlerpfad entfällt strukturell — es bleibt nur
die Klasse der Domänen-Invarianten-Verletzung (siehe Teilfrage 5).

### Teilfrage 4 — Authentifizierung/Autorisierung

`internal/adapters/driving/grpc/interceptor.go` trägt bereits
`role`/`classifyToken`/`credentialToken` für den bestehenden
`authStreamInterceptor` (`ADR-0060` Teilfrage 4) — ein
`grpc.StreamServerInterceptor`, der ausschließlich Streaming-RPCs schützt.
`Administration` trägt ausschließlich unäre RPCs; sie brauchen einen
`grpc.UnaryServerInterceptor`, der dieselben drei bereits vorhandenen
Bausteine wiederverwendet (keine dritte, wortgleiche Fassung nötig — beide
Interceptoren liegen im selben Paket `grpc`).

| Option | Pro | Contra |
|---|---|---|
| A — ein Interceptor pro RPC-Methode (neun separate Wrapper, analog zu HTTP `withToken(reader, admin, required, handler)` je Handler) | pro Methode explizit lesbar, welche Rolle sie verlangt | gRPC verdrahtet Interceptoren am `grpc.Server` als ganzes, nicht pro Methode — ein Wrapper je Methode bräuchte eine eigene, nicht-standardmäßige Dispatch-Konstruktion |
| **B — ein gemeinsamer `grpc.UnaryServerInterceptor` mit einer Rechtsklassen-Tabelle je Methodenname (gewählt)** | ein Interceptor, an `grpc.NewServer(grpc.UnaryInterceptor(...))` registriert, deckt alle neun RPCs; die Tabelle (`map[string]role`) ist so explizit wie neun einzelne Wrapper, aber an einer Stelle; folgt demselben Muster wie `withToken`s `required role`-Parameter, nur über einen Lookup statt neun Aufrufstellen | ein unbekannter/fehlgeschriebener Methodenname in der Tabelle wäre ein stiller Bug — dagegen: fail-closed-Default (siehe unten) macht ihn sichtbar als `PermissionDenied` statt als klaffende Lücke |
| C — keine Rechtsklassen-Differenzierung, jedes gültige Token (`reader` oder `admin`) erreicht jede RPC | einfachster Interceptor | widerspricht `ADR-0057` Teilfrage 3s Least-Privilege-Modell direkt — ein `reader`-Token dürfte dann `RegisterConsumer` aufrufen, das HTTP-Äquivalent verlangt `admin`; fachliche Ungleichheit zwischen den Zugriffswegen, verboten durch `LH-FA-SST-006`s Boundary-Kriterium |

Der Interceptor liest den bareren RPC-Namen aus `info.FullMethod`
(Form `/cdc.administration.v1.Administration/<Methode>`, letztes
Pfadsegment) und schlägt ihn in einer bei Konstruktion übergebenen Tabelle
nach. Rollen-Zuordnung — identisch zu `ADR-0057` Teilfrage 3s HTTP-Zuordnung:

- **`roleAdmin` (deckt `roleReader` implizit ab):** `RegisterConsumer`,
  `AcknowledgeConsumer`, `RemoveConsumer`, `EnableTable`, `DisableTable`,
  `RunRetention`.
- **`roleReader`:** `GetConsumerPosition`, `GetTableStatus`, `ListTables`.

Ein nicht in der Tabelle gefundener Methodenname (sollte am fest verdrahteten
RPC-Satz nicht vorkommen) fällt **fail-closed** auf `roleAdmin` — dieselbe
Vorsichtsrichtung wie `classifyToken`s Behandlung eines leer konfigurierten
Tokens (ein Lücken-Fall eröffnet keine implizite dritte, schwächere
Rechtsklasse).

`grpc.NewServer(grpc.StreamInterceptor(authStreamInterceptor(...)),
grpc.UnaryInterceptor(authUnaryInterceptor(...)))`: Beide Optionen koexistieren
konfliktfrei, weil `ChangeStream` ausschließlich Streaming-RPCs und
`Administration` ausschließlich unäre RPCs trägt — kein RPC durchläuft beide
Interceptoren, keine Lücke zwischen ihnen.

### Teilfrage 5 — Fehlerform

Ein `codes.*`-Status ersetzt den jeweiligen HTTP-Statuscode aus
`writeError`/`writeDomainError` (`internal/adapters/driving/http/{middleware,errors}.go`):

| HTTP-Status (Ursache) | gRPC-Code |
|---|---|
| `400` — Request-Body kein gültiges JSON | entfällt strukturell (Teilfrage 3) |
| `400` — Domänen-Invariante verletzt (`domainerrors.Err*`, `outbound.ErrNonPositiveLimit`, `outbound.ErrRangeInverted`) | `codes.InvalidArgument` |
| `401` — fehlender/unbekannter Bearer-Token | `codes.Unauthenticated` |
| `403` — Rechtsklasse unzureichend | `codes.PermissionDenied` |
| `404` — `inbound.ErrSourceTableMissing` | `codes.NotFound` |
| `500` — unerwarteter interner Fehler | `codes.Internal` |

Die Zuordnung ist eine reine Übersetzungstabelle: Der Handler ruft denselben
Inbound Use Case wie der HTTP-Adapter auf und übersetzt dessen Fehler-Sentinel
in den `codes.*`-Wert dieser Tabelle statt in einen HTTP-Statuscode — dieselbe
Fehlerklassifikation (`ADR-0023`), zwei Zieldarstellungen.

### Teilfrage 6 — Server-Instanz und Adapter-Platzierung

| Option | Pro | Contra |
|---|---|---|
| A — eigener `grpc.Server` auf neuem Port (z. B. `CDC_GRPC_ADMIN_ADDR`) | isoliert Administration- von Streaming-Traffic auf Netzwerkebene | neue ENV-Variable ohne fachlichen Grund — dasselbe Trust-Modell, dieselben zwei Tokens; zwei TCP-Listener für dieselbe Auth-Domäne, ohne dass ein Akzeptanzkriterium Trennung verlangt |
| **B — derselbe `grpc.Server`, dieselbe Adresse `CDC_GRPC_ADDR`, zusätzliche `RegisterAdministrationServer`-Registrierung im bestehenden Paket `internal/adapters/driving/grpc/` (gewählt)** | keine neue ENV-Variable — additiv im engsten Sinn: ein Deployment, das `CDC_GRPC_ADDR` bereits setzt, bekommt die neuen Fähigkeiten ohne weitere Konfiguration; ein Deployment ohne `CDC_GRPC_ADDR` bleibt beim bestehenden No-Op; gRPC registriert nativ mehrere Services auf einem Server (Standardmuster); ein Interceptor-Paar auf demselben `grpc.NewServer(...)`-Aufruf deckt beide Dienste | ein Ausfall des gemeinsamen Listeners betrifft künftig beide Fähigkeiten statt nur Streaming — akzeptiert (Konsequenzen) |
| C — neues Paket `internal/adapters/driving/grpcadmin/` mit eigenem `grpc.Server` | physische Datei-Trennung nach Fähigkeit | inkonsistent zum HTTP-Vorbild: `consumer.go`/`verwaltung.go`/`retention.go` trennen nach **Datei**, nicht nach **Paket** — alle liegen im selben `http`-Paket hinter demselben Server/Port; ein Paket-Split hier bräuchte ohnehin einen zweiten Port (Option A) oder eine erklärungsbedürftige geteilte-Listener-Konstruktion |

Der Adapter importiert für die neun Handler ausschließlich Inbound Ports
(`internal/application/port/inbound/`) und Domain-Typen zur
Request-/Response-Übersetzung — keine Driven-Adapter-Interna, keine
Application-Interna, dieselbe Grenze wie beim bestehenden `StreamChanges`-
Handler und beim HTTP-Adapter (`ARC-002`/`ARC-003`/`ARC-005`).

### Kompatibilität

Das Hinzufügen des neuen `Administration`-Service in einem **neuen** Paket
(`cdc.administration.v1`) ist stärker als „nur rückwärtskompatibel": Die
bestehende Datei `proto/cdc/stream/v1/changestream.proto` und ihr generierter
Code bleiben **byte-identisch unverändert** — kein Diff, kein Re-Compile-
Risiko für bestehende `StreamChanges`-Clients (Go/C#/Kotlin-SDKs, die drei
Beispiel-Clients). Selbst im verworfenen Fall (Option A, RPCs im bestehenden
Service) wäre das Hinzufügen neuer RPC-Methoden zu einem bestehenden
`.proto`-Service laut Protobuf-/gRPC-Semantik rückwärtskompatibel (ein alter
Client ruft nur die ihm bekannten Methoden auf) — diese Frage stellt sich
mit der gewählten Option B aber gar nicht erst, weil die betroffene Datei
unangetastet bleibt.

## Konsequenzen

- Positiv: Ein Consumer, der ausschließlich gRPC spricht, erreicht jetzt
  dieselben neun Verwaltungsfähigkeiten wie über HTTP — schließt genau die
  Lücke, die `ADR-0057`s Negativ-Konsequenz „ein Netzwerk-Consumer, der
  [Changes-Lesen/Diagnose] braucht, muss auf einen Folge-Slice warten"
  seinerzeit für die HTTP-exklusiven Fähigkeiten offenließ — hier
  spiegelbildlich für die gRPC-exklusive Lücke.
- Positiv: Keine neue Secret-Klasse, keine neue ENV-Variable — `CDC_GRPC_ADDR`
  und die beiden bestehenden Token-Klassen tragen die Erweiterung vollständig
  (Least-Privilege-Kontinuität aus `ADR-0047`/`ADR-0057`/`ADR-0060`).
- Positiv: `ChangeStream`/`changestream.proto` bleiben unverändert — keine
  Regressionsgefahr für die bestehende Live-Streaming-Fähigkeit.
- Positiv: `.a-check.yml` bleibt unverändert — beide neuen Pfade
  (`internal/adapters/driving/grpc/*.go`, `gen/cdc/administration/v1/**`)
  liegen bereits in den bestehenden, rekursiven Globs (geprüft durch Lesen
  von `.a-check.yml`).
- Negativ: Der bislang rein auf Live-Streaming zugeschnittene
  Driving-Adapter `internal/adapters/driving/grpc/` wächst um neun
  Handler-Implementierungen samt Unary-Interceptor — dieselbe Wachstums-
  Konsequenz, die `ADR-0057` für den HTTP-Adapter bereits einging.
- Negativ: Ein Ausfall des gemeinsamen `grpc.Server`-Listeners
  (Entscheidung 6, Option B) betrifft künftig sowohl Streaming als auch
  Administration — akzeptiert, weil beide dasselbe Trust-/Betriebsmodell
  teilen und kein Akzeptanzkriterium eine Netzwerktrennung verlangt.
- Negativ: Zwei generierte Go-Paket-Wurzeln (`gen/cdc/stream/v1`,
  `gen/cdc/administration/v1`) statt einer — ein zweiter `protoc`-Lauf in
  derselben Docker-Build-Stufe (`proto`/`proto-export`,
  `tools/harness/proto-generate.sh`), mechanisch, aber Folgepflicht des
  umsetzenden Slices.
- Folgepflicht: Der umsetzende Slice-Schnitt implementiert `administration.proto`,
  den generierten Code, die neun Handler in
  `internal/adapters/driving/grpc/` (jeder ruft denselben Inbound Use Case
  wie sein HTTP-Äquivalent auf — kein zweiter Domänenpfad, `LH-FA-SST-006`
  Boundary), `authUnaryInterceptor` samt Rechtsklassen-Tabelle, die
  Registrierung von `Administration` am bestehenden `grpc.Server` und die
  Erweiterung von `tools/harness/proto-generate.sh`/Dockerfile um die zweite
  `.proto`-Quelle.
- Folgepflicht: `spec/pflichtenheft.md` erhält einen neuen `SPEC-*`-Eintrag
  für das konkrete Nachrichtenschema, den Service-/RPC-Namen und die
  Fehlercode-Tabelle (analog zu `ADR-0057`/`ADR-0060`s gleichlautender
  Folgepflicht) — Gegenstand des umsetzenden Slices, nicht dieser ADR.
- Folgepflicht: Beispiel-Clients (Go/C#/Kotlin, volle Matrix laut
  `SPEC-018`§„Sprachen und Umfang") und die Erweiterung der drei SDK-Packages
  (`PgChangeFeed.Client`, `pgchangefeed`, `pgchangefeed-kotlin`) um die neun
  Administration-RPCs — je eigener Folge-Schritt/Slice, nicht Gegenstand
  dieser ADR.
- Folgepflicht: E2E-Beleg (`make test-integration`) über mindestens eine RPC
  je Token-Klasse, plus die beiden Negative-Pfade
  (`Unauthenticated`/`PermissionDenied`) — analog zum bestehenden
  gRPC-Stream-Rundlauf, Umsetzung Gegenstand des umsetzenden Slices.
- Folgepflicht: `spec/architecture.md` braucht **keine** Änderung —
  `ARC-005` nennt gRPC bereits generisch als Driving-Adapter-Technologie;
  diese ADR macht die Wahl für die Verwaltungsfähigkeiten ein zweites Mal
  konkret, ohne die Sicht selbst zu ändern.
- Nicht Gegenstand dieser ADR: `LH-FA-REA-*`/`SPEC-022` (Bestandslesen über
  einen Bereich) über gRPC — kein Auftrag benennt diese Fähigkeit, und der
  bestehende Live-Stream (`ADR-0060`) deckt neue Changes bereits ab (siehe
  Kontext, Out-of-Scope-Klarstellung).

## Fitness Function (falls maschinell prüfbar; Erwartung an den umsetzenden Slice, nicht durch diesen Architect-Lauf gefahren)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driving/grpc`, Whitebox, `authUnaryInterceptor`) | Ein Aufruf ohne oder mit unbekanntem `authorization`-Metadata-Wert endet mit `codes.Unauthenticated`; ein gültiges `reader`-Token gegen eine `roleAdmin`-RPC (z. B. `RegisterConsumer`) endet mit `codes.PermissionDenied`; ein gültiges `admin`-Token erreicht sowohl `roleReader`- als auch `roleAdmin`-RPCs | `make test` |
| Go-Unit-Test (`internal/adapters/driving/grpc`, Whitebox, je Handler) | Jeder der neun Handler ruft exakt den Inbound Use Case seines HTTP-Äquivalents auf (z. B. `RegisterConsumer`-RPC ruft `inbound.RegisterConsumerUseCase.Register` — Fake-Use-Case im Test) und bildet dessen Fehler-Sentinels auf die `codes.*`-Tabelle aus Teilfrage 5 ab | `make test` |
| `.a-check` | unverändert — beide neuen Pfade liegen vollständig in den bestehenden, rekursiven Globs `adapters: ["internal/adapters/**"]`/`contract: ["gen/**"]`, kein neuer Hexagon-Schichten-Edge nötig (geprüft durch Lesen von `.a-check.yml` bei Abfassung dieser ADR) | `make a-check` |
| `make test-integration` | realer gRPC-Rundlauf gegen den laufenden Feed-Container: mindestens eine `Administration`-RPC je Token-Klasse erreicht ihr Ziel und ist über `cdc.consumer`/`cdc.source_table` unabhängig bestätigt; ein Aufruf ohne Token endet mit `Unauthenticated`, ein `reader`-Token gegen eine `roleAdmin`-RPC mit `PermissionDenied` — Umsetzung Gegenstand des umsetzenden Slice-Schnitts | `make test-integration` |

## Re-Evaluierungs-Trigger

Permanent — die Wahl „`Administration`-Service im eigenen `.proto`-Paket,
neun unäre RPCs, gemeinsamer `grpc.Server`/Port, wiederverwendete
Token-Klassen" gilt unabhängig vom Umsetzungszeitpunkt. Ausnahmen, die eine
eigene Folge-ADR brauchen: (1) ein konkreter Bedarf an Netzwerktrennung
zwischen Streaming- und Administration-Traffic (Entscheidung 6, Option A
neu bewerten); (2) ein konkret benannter Bedarf an `LH-FA-REA-*`/`SPEC-022`
über gRPC (Out-of-Scope-Klarstellung im Kontext); (3) ein Bedarf an
Rechtsklassen-Verfeinerung über `reader`/`admin` hinaus.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-28 | Accepted — Architect-Entscheidung nach explizitem Auftraggeber-Auftrag (gRPC bekommt dieselben Verwaltungsfähigkeiten wie die HTTP-API); Supersedes `ADR-0057` Teilfrage 1, teilweise | [`ADR-0057`](0057-http-grpc-api.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0130` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
