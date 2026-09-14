# ADR-0060: gRPC-Server-Streaming für Live-Change-Zustellung

**Status:** Accepted

**Datum:** 2026-09-14

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-14)

**Bezug:** [`LH-FA-SST-008`](../../../spec/lastenheft.md) (Haupt-Bezug —
Live-Streaming vollständiger Change-Inhalte), [`LH-FA-SST-007`](../../../spec/lastenheft.md)
(Abgrenzung — NATS bleibt eigenständiges Wecksignal), [`LH-FA-REA-001`](../../../spec/lastenheft.md)
(bestehender Lesezugriffsweg als Nachhol-Pfad), [`LH-FA-CON-003`](../../../spec/lastenheft.md),
[`LH-FA-CON-005`](../../../spec/lastenheft.md) (bestätigte Consumer-Position
als Fortsetzungs-Referenz), [ADR-0057](0057-http-grpc-api.md) (Abgrenzung —
betrifft die administrative Verwaltungs-API, nicht diese Entscheidung; trägt
den hier eingetretenen Re-Evaluierungs-Trigger), [ADR-0055](0055-nats-change-notification-wecksignal.md),
[ADR-0056](0056-nats-tabellen-granulares-subjekt.md) (NATS-Präzedenzfall:
reines Wecksignal, Core-NATS-Fire-and-Forget-Muster als Vorbild für die
Zustellsemantik hier), [ADR-0030](0030-testpyramide.md) (Testpyramide),
[ADR-0028](0028-inbound-use-cases.md) (Inbound Use Cases — Abgrenzung: diese
ADR führt keinen neuen Inbound Use Case, sondern einen neuen Outbound Port),
[ADR-0034](0034-ports-nach-faehigkeiten.md) (Port-Zuschnitt nach Fähigkeit,
Vorbild für den neuen `ChangeStreamPort`)

**Schärft:** [`ARC-005`](../../../spec/architecture.md) (macht das dort
bereits generisch genannte „später HTTP-/gRPC" für **diese** Fähigkeit
konkret: gRPC Server-Streaming, für Live-Change-Zustellung — unabhängig von
und zusätzlich zu `ADR-0057`s HTTP/JSON-Wahl für die Verwaltungs-API)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0057`](0057-http-grpc-api.md) wählte HTTP/JSON für die administrative
Verwaltungs-API und verwarf gRPC (Option B) explizit mit zwei
Begründungen: kein konkreter gRPC-Consumer-Bedarf zum Entscheidungszeitpunkt,
und die dort betroffenen Fähigkeiten (Register/Enable/Disable/Retention …)
sind reine Request/Response-Aufrufe ohne nativen Streaming-Vorteil. Diese
ADR ist **keine Korrektur und keine Supersedes-ADR zu `ADR-0057`** — die
Verwaltungs-API bleibt unverändert HTTP/JSON mit Token-Authn. Gegenstand
hier ist eine **andere** Fähigkeit: `LH-FA-SST-008` (im Lastenheft seit
Version 0.8.0, 2026-09-14) verlangt Live-Streaming vollständiger
Change-Inhalte über gRPC, ohne dafür den bestehenden Lesezugriffsweg
pollen zu müssen. Der Auftraggeber (pt9912) hat diesen Bedarf am
2026-09-14 explizit benannt — exakt der in `ADR-0057`s
Re-Evaluierungs-Trigger vorgezeichnete Fall („ein konkreter
gRPC-Consumer benennt sich"), aber als **eigene Entscheidung für eine
eigene Anforderung** modelliert, nicht als Wiedereröffnung von
`ADR-0057`s Teilfrage 1 (Baseline-Regelwerk `modul-08-agentenrollen.md`
§Konflikt-Pfad als Rollen-Sequenz, Verdikt 2: Folge-ADR statt stiller
Korrektur — hier sogar strenger, weil die neue ADR nicht einmal
`ADR-0057`s Textkörper anfasst, sondern eine disjunkte Fähigkeit trifft).

Der fachliche Unterschied trägt die Protokoll-Wahl: Live-Change-Zustellung
ist ein Server-Streaming-Lehrbuchfall — ein Client öffnet eine Verbindung,
der Server sendet über die Zeit viele Nachrichten, ohne dass der Client
erneut anfragen muss. gRPC hat dafür nativen Support (HTTP/2-basiertes
Server-Streaming, ein RPC-Aufruf, ein lang laufender Antwortstrom).
HTTP/JSON bräuchte eine Ersatzkonstruktion (Server-Sent Events oder
Long-Polling), die dasselbe Ziel über Mittel erreicht, für die das
Protokoll nicht gebaut ist.

Drei Konstraints prägen den Lösungsraum:

- **`LH-FA-SST-008`s eigenes Boundary-Akzeptanzkriterium ist eng gefasst:**
  Bei Unterbrechung des Streams dürfen Changes „nicht verloren" gehen —
  sie müssen über den bestehenden Lesezugriffsweg (`LH-FA-REA-001` ff.)
  und über die bestätigte Consumer-Position (`LH-FA-CON-003`/`005`)
  **abrufbar/fortsetzbar bleiben**. Ein Replay **innerhalb** des Streams
  selbst wird nicht verlangt — der Store bleibt die Quelle der
  Nachvollziehbarkeit, nicht der Stream.
- **`LH-FA-SST-008`s Out-of-Scope grenzt `LH-FA-SST-007` (NATS) ausdrücklich
  ab:** „NATS bleibt ein eigenständiges Wecksignal, dieses Streaming
  ersetzt es nicht." Beide Mechanismen existieren nebeneinander mit
  unterschiedlichem Zweck (Wecksignal ohne Inhalt vs. vollständiger
  Inhalt über einen dedizierten Kanal).
- **Bestehender Capture-kritischer Pfad ist unantastbar
  (`ADR-0011`/`ADR-0027`):** `Receive → Decode → Persist → COMMIT Store →
  ACK Source` darf durch keinen neuen Schritt verzögert, blockiert oder
  rückgängig gemacht werden — dieselbe Grenze, die `ADR-0055` Punkt 4 für
  das NATS-Wecksignal bereits gezogen hat.

## Entscheidung

Wir wählen **gRPC Server-Streaming über ein neues Adapter-Paar** —
Driving-Adapter `internal/adapters/driving/grpc/` (Server, Auth-Interceptor)
und Driven-Adapter `internal/adapters/driven/grpcstream/` (In-Prozess-
Broadcaster hinter einem neuen Outbound Port `ChangeStreamPort`) — mit
sechs Teilfragen.

### Teilfrage 1 — Mechanismus für Live-Streaming

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (nur SQL-Poll/NATS-Wecksignal bleiben) | kein Aufwand | erfüllt `LH-FA-SST-008` nicht: Wecksignal trägt laut `ADR-0055` Punkt 3 bewusst keinen Inhalt, ein Consumer müsste nach jedem Signal zusätzlich pollen — genau das Pollen, das `LH-FA-SST-008` vermeiden soll |
| B — HTTP mit Server-Sent Events (SSE) | passt zum bestehenden `net/http`-Adapter aus `ADR-0057`, keine neue Fremd-Abhängigkeit | SSE ist reine Ersatzkonstruktion für Server-Streaming auf einem Protokoll ohne nativen Support (kein typisiertes Nachrichtenschema, keine eingebaute Flusskontrolle, textbasiertes Framing); `LH-FA-SST-008` verlangt „vollständigen Change-Inhalt" — ohne Protobuf-Schema bliebe die Nachrichtenform so lose wie bei JSON-über-SSE, während gRPC dafür die native, im Auftrag mitgenannte Form ist |
| **C — gRPC Server-Streaming (gewählt)** | nativer Server-Streaming-Support (ein RPC-Aufruf, viele Antworten über die Zeit) — genau der hier vorliegende Lehrbuchfall; typisiertes Nachrichtenschema (Protobuf) für „vollständigen Change-Inhalt"; deckt den explizit benannten Consumer-Bedarf direkt | neue Toolchain-Abhängigkeit (protoc/buf, neue Docker-Build-Stufe) — derselbe Kostenpunkt, den `ADR-0057` Teilfrage 1 Option B für die Verwaltungs-API noch als unverhältnismäßig verworfen hatte; hier aber gerechtfertigt, weil ein konkreter Bedarf **und** ein struktureller Streaming-Vorteil zusammenfallen, die bei der Verwaltungs-API beide fehlten |
| D — NATS selbst auf vollständige Payloads erweitern (Core NATS oder JetStream, `ChangeNotificationPort` trägt künftig den Change-Inhalt) | keine neue Protokoll-/Toolchain-Einführung, ein Mechanismus statt zwei | widerspricht `ADR-0055` Punkt 3 (leerer Payload als bewusste Entscheidung gegen jede Zweitrepräsentation) und `LH-FA-SST-008`s eigenem Out-of-Scope („NATS bleibt ein eigenständiges Wecksignal, dieses Streaming ersetzt es nicht") — eine Änderung an `ChangeNotificationPort`s Vertrag wäre eine stillschweigende inhaltliche Änderung einer `Accepted`-ADR (`AGENTS.md` §3.5, verboten) und bräuchte selbst eine Folge-ADR mit `Supersedes ADR-0055`, die diese ADR nicht mitentscheidet; verlöre zudem `ADR-0055`s zentralen Architektur-Hebel (Core NATS ohne Zustellgarantie ist für ein reines Wecksignal richtig, aber falsch für vollständige, nicht-verlustbehaftete Inhaltszustellung) |

Option D wird nicht nur als schwächer, sondern als **strukturell
unzulässig auf diesem Weg** verworfen: Sie würde eine bereits
`Accepted`-ADR inhaltlich umbiegen, ohne den dafür vorgesehenen
`Supersedes`-Weg zu gehen. Bleibt NATS unverändert Wecksignal (wie in
Option C), ist keine Folge-ADR zu `ADR-0055`/`ADR-0056` nötig.

### Teilfrage 2 — Anknüpfungspunkt am Capture-Pfad

| Option | Pro | Contra |
|---|---|---|
| A — bestehenden `ChangeNotificationPort` erweitern (Payload-Parameter ergänzen) | ein Port statt zwei, keine neue Datei | ändert den Vertrag einer `Accepted`-ADR (`ADR-0055` Punkt 3: leerer Payload ist bewusste Entscheidung, `ADR-0056` hat die Signatur bereits einmal um `schema`/`table` erweitert — eine weitere Erweiterung um vollständigen Inhalt widerspräche der in `ADR-0055` explizit begründeten Payload-Leere und wäre eine stillschweigende Korrektur ohne `Supersedes`; bricht zudem den bestehenden `natsnotify`-Adapter, der eine leere `Notify`-Semantik implementiert |
| B — neuer, separater Outbound Port `ChangeStreamPort` + eigener Broadcaster (gewählt) | `ChangeNotificationPort`/NATS bleiben vollständig unverändert (kein Bruch); folgt `ADR-0034`s Prinzip „Ports nach Fähigkeit und Konsistenzgrenze" — Wecksignal und Vollinhalts-Streaming sind zwei verschiedene Fähigkeiten mit verschiedenen Konsistenzanforderungen; `CaptureService` bleibt für beide Ports gleich strukturiert (analog `WithChangeNotification`/`WithChangeStream`) | ein weiterer optionaler Port-Parameter an `CaptureService`, eine weitere Konstruktions-Option |
| C — gRPC-Adapter pollt `ChangeStorePort` periodisch statt Push vom Capture-Pfad | keine Änderung an `CaptureService` nötig | reintroduziert genau das Pollen, das `LH-FA-SST-008`s Beschreibung ausdrücklich vermeiden will („ohne dafür den bestehenden Lesezugriffsweg … pollen zu müssen"); Latenz an ein Poll-Intervall gebunden statt ereignisgetrieben |

`ChangeNotificationPort` (`internal/application/port/outbound/changenotification.go`)
und der `natsnotify`-Adapter bleiben durch diese Entscheidung
**unverändert** — kein Zeilen-Diff an bestehendem Code dieser Dateien ist
Folge dieser ADR.

**Port-/Adapter-Design** (Folgepflicht des umsetzenden Slices, hier
architektonisch festgelegt, analog zur Design-Tiefe in `ADR-0055`/`ADR-0057`):

- Neuer Outbound Port `internal/application/port/outbound/changestream.go`:
  `ChangeStreamPort` mit `Publish(ctx context.Context, change *model.Change)
  error` — ein Aufruf je einzelnem `Change` der committed Transaktion
  (**keine** Deduplizierung nach Tabelle wie bei `ADR-0056`s Wecksignal:
  `LH-FA-SST-008` fordert den vollständigen Inhalt „einer neuen committed
  Änderung", also je Zeilen-Change, nicht je Tabelle), in der
  Reihenfolge von `tx.Changes()`.
- `CaptureService` (`internal/application/usecase/capture/service.go`)
  erhält eine neue, optionale Konstruktions-Option `WithChangeStream(stream
  outbound.ChangeStreamPort) Option`, analog zu `WithChangeNotification`.
  Der Aufruf reiht sich **nach** `ACK Source` und **nach** dem bestehenden
  Notify-Schritt ein (`Receive → Decode → Persist → COMMIT Store → ACK
  Source → Notify (best effort) → Stream-Publish (best effort)`); sein
  Fehler wird an der Aufrufstelle abgefangen (`LogPort`-Warnung) und geht
  nie in den Rückgabewert von `Capture()` ein — dieselbe Isolationsregel
  wie bei `ChangeNotificationPort` (`ADR-0055` Punkt 4). Ungesetzt (`nil`)
  bleibt die Fähigkeit vollständig deaktiviert, bestehendes Verhalten
  bit-identisch.
- Implementiert wird der Port von einem neuen Driven-Adapter
  `internal/adapters/driven/grpcstream/` (`Broadcaster`): ein
  nebenläufigkeitssicherer In-Prozess-Fan-out ohne externes System — jeder
  aktive gRPC-Stream registriert sich über `Subscribe() (<-chan
  *model.Change, cancel func())`, `Publish` verteilt an alle registrierten
  Kanäle, ohne Empfänger wird die Nachricht verworfen (kein Puffer, siehe
  Teilfrage 3).
- Der Driving-Adapter `internal/adapters/driving/grpc/` (gRPC-Server,
  generierter Service-Stub, Auth-Interceptor) referenziert den
  `Broadcaster` **nicht** über einen Import des Driven-Pakets, sondern über
  ein lokal im `grpc`-Paket deklariertes, strukturell erfülltes Interface
  (`changeSubscriber` mit `Subscribe() (<-chan *model.Change, func())`) —
  Bootstrap (`ARC-007`, Composition Root) verdrahtet den konkreten
  `*grpcstream.Broadcaster` in beide Konstruktoren (`capture.WithChangeStream`
  und den gRPC-Server). Damit bleibt die Richtungskonvention
  Driving/Driven aus `spec/architecture.md` §1 gewahrt: Kein Adapter-Paket
  importiert ein anderes Adapter-Paket direkt; nur `internal/bootstrap`
  kennt beide konkret.
- `.a-check.yml` braucht **keine** Änderung — beide neuen Pakete liegen
  vollständig im bestehenden Glob `adapters: ["internal/adapters/**"]`
  bzw. `ports: ["internal/application/port/**"]`; kein neuer
  Hexagon-Schichten-Edge nötig (geprüft: die einzige Kopplung zwischen den
  beiden neuen Paketen läuft über eine strukturell erfüllte Interface im
  Driving-Paket plus Domain-Typ `model.Change`, beides bereits erlaubte
  Kanten `adapters → domain`, `adapters → ports`).

### Teilfrage 3 — Zustellsemantik und Reconnect

| Option | Pro | Contra |
|---|---|---|
| **A — Fire-and-Forget ohne Replay (gewählt)** | erfüllt `LH-FA-SST-008`s Boundary-Kriterium exakt („bleiben über den bestehenden Lesezugriffsweg abrufbar … und über die bestätigte Consumer-Position fortsetzbar") — kein Akzeptanzkriterium verlangt Replay **im** Stream; folgt demselben Architektur-Hebel wie `ADR-0055` Punkt 1 (Core NATS: die Nachvollziehbarkeit liegt bereits vollständig beim `ChangeStorePort`, ein zweiter Mechanismus dafür wäre eine zweite Quelle der Wahrheit); minimale Zustandshaltung im `Broadcaster` (kein Puffer, keine Consumer-Positions-Verwaltung im Stream selbst) | ein Consumer, der zwischen zwei Changes kurz getrennt war, verpasst sie über den Stream ersatzlos — muss selbst über den bestehenden Lesezugriffsweg (`LH-FA-REA-001` ff.) und seine zuletzt bestätigte Position (`LH-FA-CON-003`/`005`) nachholen |
| B — Replay ab zuletzt bestätigter Consumer-Position beim (Re-)Connect | Consumer bekäme lückenlose Zustellung ohne selbst nachzuholen | baut eine zweite Nachvollziehbarkeits-Infrastruktur im Stream auf, die den `ChangeStorePort`-Lesepfad dupliziert (dieselbe Sorge wie bei JetStream in `ADR-0055` Option B: zwei Wahrheiten, ungeklärte Konfliktauflösung); verlangt, dass der Stream-Server bei jedem Connect eine Consumer-Identität kennt und einen Store-Lookup + Live-Übergang synchronisiert — deutlich höherer Aufwand, den kein Akzeptanzkriterium von `LH-FA-SST-008` fordert |
| C — In-Memory Ring-Buffer begrenzter Tiefe im `Broadcaster` | überbrückt sehr kurze Verbindungslücken ohne Store-Zugriff | begrenzte, unvorhersagbare Zusicherung (Tiefe müsste geschätzt werden, ohne dass ein Akzeptanzkriterium eine Zahl vorgibt); löst das Problem nicht grundsätzlich, sondern verschiebt die Lücke nur auf „länger als der Puffer" — zusätzlicher Zustand ohne belastbaren Nutzen |

Die Entscheidung ist bewusst: `LH-FA-SST-008` selbst verlangt kein
Stream-internes Replay, und der Store bleibt die einzige Quelle der
Wahrheit — dieselbe Trennschärfe wie bei `LH-FA-SST-007`/NATS, nur für
eine andere Fähigkeit angewendet.

### Teilfrage 4 — Authentifizierung/Autorisierung

| Option | Pro | Contra |
|---|---|---|
| A — kein Schutz (nichts tun) | kein Zusatzaufwand | `LH-FA-SST-008`s Negative-Akzeptanzkriterium wird strukturell unerfüllbar — „ein Verbindungsversuch ohne gültige Authentifizierung/Autorisierung … wird abgelehnt" |
| B — eigener, neuer Token-Mechanismus (z. B. `CDC_GRPC_TOKEN`) | vom HTTP-Tokenraum entkoppelt, könnte unabhängig rotiert werden | eine dritte Secret-Klasse neben `CDC_API_TOKEN_READER`/`ADMIN` ohne fachlichen Grund — Streaming ist eine Lesefähigkeit wie `GetStatus`/`ListTables` aus `ADR-0057`, keine eigene Rechtsklasse; erhöht die Betriebs-/Rotations-Fläche ohne Mehrwert |
| **C — Wiederverwendung der bestehenden `CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`-Klassen über gRPC-Metadata (gewählt)** | konsistentes Least-Privilege-Modell über beide Netzwerk-Zugriffswege hinweg (`ADR-0057`), keine neue Secret-Klasse; Streaming ist rein lesend — ein gültiges `reader`- **oder** `admin`-Token genügt (kein `403`-Analogon nötig, weil es keine zweite Rechtsklasse für diese eine Fähigkeit gibt); Prüfung über einen gRPC-Unary/Stream-Interceptor, der den Metadata-Key `authorization` gegen dieselben zwei konfigurierten Werte vergleicht | Tokens bleiben Shared Secrets ohne Ablauf-/Widerrufsmechanismus (dieselbe bereits in `ADR-0057` benannte Einschränkung); TLS-Terminierung vor dem gRPC-Endpunkt bleibt Betreiber-Pflicht |

Ein Stream-Öffnungsversuch ohne `authorization`-Metadata oder mit einem
Token, das keiner der beiden konfigurierten Klassen entspricht, endet mit
gRPC-Status `Unauthenticated` — sichtbar, nicht still mit leeren Daten
fortgesetzt (`LH-FA-SST-008` Negative). Anders als bei `ADR-0057` gibt es
hier **keine** zweite Ablehnungsklasse (`PermissionDenied`/HTTP-403-Analogon):
Live-Streaming ist eine einzelne, rein lesende Fähigkeit, für die ein
`reader`-Token bereits hinreichend ist.

### Teilfrage 5 — Adapter-Platzierung

| Option | Pro | Contra |
|---|---|---|
| A — ein einziges Paket `internal/adapters/driving/grpc/`, das sowohl den Server als auch den `ChangeStreamPort`-Broadcaster trägt | keine zweite neue Datei-Gruppe | vermischt Driving- und Driven-Rolle in einem Paket entgegen der in `spec/architecture.md` §1 etablierten Namenskonvention (`driving/`/`driven/` nach Richtung), obwohl `.a-check.yml` das nicht separat prüft — schwerer lesbar, welche Hälfte welche Richtung bedient |
| **B — getrennte Pakete `internal/adapters/driving/grpc/` (Server, Interceptor) und `internal/adapters/driven/grpcstream/` (Broadcaster, `ChangeStreamPort`-Implementierung) (gewählt)** | folgt der bestehenden Richtungskonvention (`ADR-0055`s `natsnotify` unter `driven/`, `ADR-0057`s Adapter unter `driving/`); `ARC-005` nennt gRPC bereits als vorgesehene Driving-Adapter-Technologie; `.a-check.yml`s bestehende Globs erfassen beide Verzeichnisse bereits — **keine Änderung an `.a-check.yml` nötig** | zwei neue Verzeichnisse statt eines; die Kopplung zwischen ihnen läuft über ein lokal deklariertes Interface statt eines direkten Imports (siehe Teilfrage 2) — etwas mehr Entwurfsaufwand als ein direkter Import |
| C — neues Top-Level-Verzeichnis außerhalb `internal/adapters/` (z. B. `internal/streaming/`) | organisatorisch „prominenter" für eine neu eingeführte Technologie | `.a-check.yml`s Glob würde das Verzeichnis nicht als `adapters`-Layer erkennen und bräuchte eine Erweiterung ohne Mehrwert gegenüber der bestehenden Konvention — derselbe Einwand wie in `ADR-0057` Teilfrage 4 Option C |

### Teilfrage 6 — Bootstrap-Aktivierung

| Option | Pro | Contra |
|---|---|---|
| A — Fähigkeit immer aktiv, kein Konfigurations-Gate | einfachste Verdrahtung | jedes bestehende Deployment (`compose.yaml`, jede heutige Env-only-Installation) müsste zwangsläufig einen gRPC-Listener öffnen — Breaking Change für bestehende Betriebsumgebungen ohne diesen Bedarf |
| **B — additive Umgebungsvariable `CDC_GRPC_ADDR`, No-Op bei fehlender Konfiguration (gewählt)** | folgt demselben Muster wie `CDC_HTTP_ADDR` (`ADR-0057`) und `CDC_NATS_URL` (`ADR-0055`); ungesetzt bleibt der gRPC-Server und der `Broadcaster` vollständig deaktiviert, `CaptureService` erhält keinen `ChangeStreamPort` — bestehendes Verhalten jedes heutigen Deployments bleibt bit-identisch | ein weiterer optionaler ENV-Vertrag, den `internal/bootstrap` dokumentieren muss |
| C — ausschließlich über die optionale YAML-Konfigurationsdatei (`ADR-0052`), keine eigene ENV-Variable | konsolidiert Konfiguration an einer Stelle | inkonsistent zum etablierten `CDC_*`-ENV-Präzedenzfall für genau diese Klasse von Netzwerk-Endpunkt-Adressen (`CDC_HTTP_ADDR`, `CDC_NATS_URL`); `ADR-0052` selbst sieht YAML als *Ergänzung* zu ENV vor, nicht als Ersatz |

## Konsequenzen

- Positiv: `LH-FA-SST-008` wird mit einem architektonisch sauber
  abgegrenzten neuen Adapter-Paar erfüllbar, ohne den bestehenden
  `ChangeNotificationPort`/`natsnotify`-Pfad anzufassen — beide
  Mechanismen (Wecksignal, Vollinhalts-Stream) bleiben unabhängig
  austauschbar (`ADR-0034`).
- Positiv: Die Zustellsemantik (Fire-and-Forget, Store bleibt einzige
  Quelle der Wahrheit) hält denselben Architektur-Hebel wie `ADR-0055` für
  das NATS-Wecksignal durch — konsistente Antwort auf dieselbe
  Grundfrage bei zwei verschiedenen Fähigkeiten.
- Positiv: Least-Privilege-Modell (`ADR-0047`, `ADR-0057`) setzt sich über
  einen dritten Zugriffsweg (gRPC) fort, ohne eine dritte Secret-Klasse
  einzuführen.
- Negativ: Neue Build-Toolchain-Abhängigkeit (protoc/buf) und neue
  Laufzeit-Abhängigkeiten (`google.golang.org/grpc`,
  `google.golang.org/protobuf`) — nach `nats.go` (`ADR-0055`) die zweite
  und dritte direkte Nicht-PostgreSQL-Abhängigkeit dieses Repos; die
  Toolchain muss Docker-only laufen (`AGENTS.md` §3.1), eine neue
  Build-Stufe im `Dockerfile` bzw. ein neues `make`-Ziel für die
  Code-Generierung ist Folgepflicht des umsetzenden Slices.
- Negativ: Ohne Stream-internes Replay muss jeder gRPC-Stream-Consumer
  eigenständig einen Resync-Pfad über den bestehenden Lesezugriffsweg
  implementieren — bewusste Entscheidung (Teilfrage 3), keine übersehene
  Lücke, analog zu `ADR-0055`s gleichlautender Konsequenz für NATS.
- Negativ: Tabellen-granulare Filterung des gRPC-Streams (analog zu
  `ADR-0056`s Wecksignal-Granularität) ist **nicht** Gegenstand dieser
  ADR — ein Consumer erhält zunächst alle Changes der konfigurierten
  Tabellen unfiltriert über den Stream; eine spätere Verfeinerung bleibt
  eine eigene Folge-ADR, falls ein konkreter Bedarf entsteht (derselbe
  Vertagungs-Grund wie bei `ADR-0057` Teilfrage 2 für Changes-Lesen/
  Diagnose über HTTP).
- Folgepflicht: Der umsetzende Slice-Schnitt implementiert `ChangeStreamPort`,
  `internal/adapters/driven/grpcstream/` (`Broadcaster`),
  `internal/adapters/driving/grpc/` (Server, Interceptor, Protobuf-Schema),
  die `CaptureService`-Erweiterung (`WithChangeStream`) und die
  Bootstrap-Verdrahtung (`CDC_GRPC_ADDR`, No-Op bei fehlender Adresse —
  additiv, kein Breaking Change, analog zu `ADR-0055` Punkt 5).
- Folgepflicht: `spec/pflichtenheft.md` erhält einen neuen `SPEC-*`-Eintrag
  für das konkrete Protobuf-/Nachrichtenschema, den RPC-Methodennamen und
  die Stream-Semantik (`LH-FA-SST-008`s Out-of-Scope weist das dorthin) —
  Gegenstand des umsetzenden Slices, nicht dieser ADR.
- Folgepflicht: `spec/architecture.md` braucht **keine** Änderung — `ARC-005`
  nennt gRPC bereits generisch als künftige Driving-Adapter-Technologie;
  diese ADR macht die Wahl für die Live-Streaming-Fähigkeit konkret, ohne
  die Sicht selbst zu ändern. Anders als bei `ADR-0055`/`ARC-013` entsteht
  hier keine neue externe Abhängigkeits-Zeile — der `Broadcaster` ist ein
  In-Prozess-Mechanismus ohne externes System.
- Folgepflicht: Erweiterung um Stream-internes Replay oder
  tabellen-granulare Filterung braucht je eine eigene Folge-ADR — nicht
  Gegenstand dieser ADR (siehe Re-Evaluierungs-Trigger).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driven/grpcstream`, Whitebox) | `Publish` verteilt an alle aktiven `Subscribe`-Kanäle; ein `Publish` ohne aktiven Subscriber blockiert nicht und liefert keinen Fehler (Fire-and-Forget-Regressionstest) | `make test` |
| Go-Unit-Test (`internal/application/usecase/capture`, Whitebox) | Ein fehlschlagender `ChangeStreamPort` darf den Rückgabewert von `CaptureService.Capture()` nicht zu einem Fehler machen, wenn `store`/`ack` erfolgreich waren — Regressionstest gegen versehentliche Fehlerpropagation, analog zum bestehenden `ChangeNotificationPort`-Test (`ADR-0055`) | `make test` |
| Go-Unit-Test (`internal/adapters/driving/grpc`, Whitebox) | Ein Stream-Öffnungsversuch ohne oder mit unbekanntem `authorization`-Metadata-Wert endet mit gRPC-Status `Unauthenticated`; ein gültiges `reader`- oder `admin`-Token öffnet den Stream erfolgreich | `make test` |
| `.a-check` | unverändert — beide neuen Pakete liegen vollständig im bestehenden Glob `adapters: ["internal/adapters/**"]`/`ports: ["internal/application/port/**"]`, kein neuer Hexagon-Schichten-Edge nötig | `make a-check` |
| `make test-integration` | realer gRPC-Stream-Rundlauf gegen den laufenden Feed-Container: eine committed Änderung erreicht einen verbundenen Stream-Client mit vollständigem Inhalt; ein Verbindungsversuch ohne gültiges Token wird abgelehnt — Umsetzung Gegenstand des umsetzenden Slice-Schnitts | `make test-integration` |

## Slice-Schnitt-Empfehlung

Regeln dieser Sektion: Empfehlung, endgültiger Schnitt liegt bei der
Planner-Rolle (Baseline-Regelwerk `modul-05-planning-harness.md`).
Analog zu `ADR-0057`s Dreiteilung (Grundgerüst → restliche Fähigkeiten/
Integration → E2E-Beleg):

1. **Slice A** — Adapter-Grundgerüst: Protobuf-/buf-Toolchain-Anbindung
   (neue, Docker-only laufende Build-Stufe), neuer Outbound Port
   `ChangeStreamPort`, `internal/adapters/driven/grpcstream/`
   (`Broadcaster`), `internal/adapters/driving/grpc/`-Server-Grundgerüst
   samt Auth-Interceptor (`Unauthenticated`-Pfad), Bootstrap-Verdrahtung
   (`CDC_GRPC_ADDR`), Unit-Tests.
2. **Slice B** — Capture-Integration: `CaptureService.WithChangeStream`,
   tatsächliche Publish-Kette `Receive → … → ACK Source → Notify →
   Stream-Publish`, Fehlerisolations-Tests (Stream-Fehler beeinflusst
   `Capture()`-Rückgabewert nicht), Fire-and-Forget-Verhalten bei
   getrenntem Client.
3. **Slice C** — Beispiel-Client (`tools/harness/grpcclient/`, Wegwerf-
   Werkzeug außerhalb der Produktionsschichten, analog zu
   `tools/harness/natssub/`/`tools/harness/httpclient/`) und
   E2E-Erweiterung (`compose.yaml`-Port-Exposition `CDC_GRPC_ADDR`,
   `run-integration-tests.sh`-Rundlauf über einen echten gRPC-Stream).

## Re-Evaluierungs-Trigger

Wird ein Bedarf an Stream-internem Replay (Fortsetzung ab einer
Consumer-Position **innerhalb** des gRPC-Streams selbst, ohne Rückgriff
auf den SQL-Lesezugriffsweg) oder an tabellen-granularer Filterung des
Streams (analog zu `ADR-0056` für NATS) konkret benannt: Beide
widersprechen der hier getroffenen Wahl (Teilfrage 3 bzw. der in
Konsequenzen benannten Vertagung) fundamental genug, dass sie eine
eigene Folge-ADR brauchen, die die jeweilige Option neu bewertet — nicht
als Korrektur dieser ADR, sondern als eigene Entscheidung, weil sich die
Prämisse (kein belegter Bedarf über das Boundary-Kriterium hinaus)
geändert hätte. Sonst permanent — die Wahl gRPC-Server-Streaming mit
Fire-and-Forget-Zustellung gilt unabhängig vom Zeitpunkt der Umsetzung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Architect-Entscheidung nach explizitem Auftraggeber-Bedarf (2026-09-14); trifft eine von `ADR-0057` disjunkte Fähigkeit (Live-Streaming statt Verwaltungs-API), keine Supersedes-ADR | [`ADR-0057`](0057-http-grpc-api.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0060` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
