# ADR-0076: Beispiel-Clients unter `examples/` — der öffentliche Draht-Vertrag als Vorbild

**Status:** Accepted — Supersedes [`ADR-0060`](0060-grpc-streaming-mechanismus.md)
in **einer** Klausel und [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
in **einer** Klausel.

Die erste: die Zuordnung des erzeugten Protobuf-Bindings zum Adapter-Paket
`internal/adapters/driving/grpc/` — in `ADR-0060` §Entscheidung Teilfrage 2
(„Der Driving-Adapter `internal/adapters/driving/grpc/` (gRPC-Server,
generierter Service-Stub, Auth-Interceptor) …") und in deren
§Konsequenzen-Folgepflicht („`internal/adapters/driving/grpc/` (Server,
Interceptor, Protobuf-Schema)"). Die zweite: die Kante-Aussage der `ADR-0068`
§Entscheidung Festlegung 2 („**und** genau eine Kante
`{from: tooling, to: adapters}`") samt der sie wiederholenden
§Fitness-Function-Zeile („Die Gruppe `tooling` … importiert **nur**
`adapters` … Zugleich bleibt der Stub-Import
(`internal/adapters/driving/grpc/streamv1`) grün").

Alles Übrige beider ADRs bleibt **hiermit bestätigt** und wird hier nicht
wiederholt: die sechs Teilfragen der [`ADR-0060`](0060-grpc-streaming-mechanismus.md)
(gRPC-Server-Streaming, `ChangeStreamPort`, `Broadcaster`,
Fire-and-Forget, Authentifizierung über die bestehenden Token-Klassen,
`CDC_GRPC_ADDR`), die [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
Festlegungen 1, 3 und 4 (`composition_root` ohne Werkzeug-Pfad, Verfeinerung
gegen Erweiterung, `test/integration/**` bestätigt) sowie deren
§Re-Evaluierungs-Trigger 2, der mit dieser ADR **eintritt** und dessen dort
vorgezeichnete Folge („die Kante wird auf diesen Bereich umgestellt; die
Berechtigung wird feiner, nicht breiter") hier ausgeführt wird.

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf explizite
Anforderung des Auftraggebers vom 2026-09-15; anderer Kontext als die
Implementer-Läufe, die `tools/harness/{httpclient,sseclient,grpcclient}` als
Belegträger geschrieben haben, und als der Planner-Lauf, der die Anforderung
aufnimmt — Modul 8 §Rollen-Regeln: „Architect schreibt")

**Bezug:** [`LH-FA-SST-006`](../../../spec/lastenheft.md) (HTTP-/JSON-API —
der Draht, den das HTTP-Beispiel anspricht), [`LH-FA-SST-008`](../../../spec/lastenheft.md)
(Live-Streaming vollständiger Change-Inhalte — der Draht beider Stream-Beispiele),
[`ARC-005`](../../../spec/architecture.md) (Driving Adapters),
[`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(§Entscheidung Festlegung 2, §Fitness Function — in der Kante-Aussage
superseded; §Re-Evaluierungs-Trigger 2 tritt ein),
[`ADR-0060`](0060-grpc-streaming-mechanismus.md) (in der
Binding-Zuordnung superseded; sonst bestätigt),
[`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md) (SSE-Endpunkt),
[`ADR-0057`](0057-http-grpc-api.md) (HTTP-API und Token-Klassen),
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Messgegenstand des Coverage-Gates — von der Folgepflicht des Umzugs berührt),
[`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md)
(Maschinenform der §2-Constraints), [`ADR-0026`](0026-composition-root.md)
(Composition Root), [`SPEC-020`](../../../spec/pflichtenheft.md)
(Protobuf-/Stream-Schema), [`AGENTS.md`](../../../AGENTS.md) §3.5/§3.6,
`.a-check.yml`, `a-check.mk`, `Makefile`, `Dockerfile` (Stufe `coverage`),
`harness/sensors/a-check.md`, `harness/sensors/coverage-gate.md`,
`docs/user/benutzerhandbuch.md` (`### Zugriff über die HTTP-/JSON-API`,
`### Zugriff über den gRPC-Change-Stream`,
`### Zugriff über Server-Sent-Events`),
`proto/cdc/stream/v1/changestream.proto`,
`tools/harness/{httpclient,sseclient,grpcclient}/main.go`,
`tools/harness/run-integration-tests.sh`,
`docs/plan/planning/observations/BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/observation.md`

**Schärft:** — (Prozess-ADR mit Architektur-Geltung, wie
[`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)/[`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md);
sie ändert keine Zusage des Lastenhefts oder Pflichtenhefts und keinen Satz
der Sicht — `examples/**` und der öffentliche Vertrags-Pfad sind
Gate-Scope-Gruppen, keine Komponenten der §2-Sicht)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Auftraggeber hat am 2026-09-15 angefragt, einen `./examples`-Ordner mit
„richtigen" Clients anzulegen, „die man auch im Handbuch erwähnen kann".
Auslöser war die Frage, wo die Client-Beispiele für HTTP, SSE und gRPC liegen.

**Der Ist-Stand.** Die drei Clients liegen unter
`tools/harness/{httpclient,sseclient,grpcclient}`. Sie sind **Wegwerf-Clients**:
ihr Träger ist `tools/harness/run-integration-tests.sh`, das sie per
`go run` im Toolchain-Container gegen den **laufenden** Compose-Feed-Container
startet und ihren **stdout** über `docker logs` auswertet (das
`READY`/`RECEIVED`/`REJECTED`-Protokoll) — nicht ihren Exit-Code allein. Sie
sind **Belegträger** eines E2E-Rundlaufs und deshalb zu Recht **nicht** im
Handbuch: Die Zeichenketten `httpclient`/`sseclient`/`grpcclient` kommen in
`docs/user/benutzerhandbuch.md` **0×** vor.

Das Handbuch dokumentiert dagegen die **Schnittstellen** vollständig — die
HTTP-/JSON-API samt `Authorization: Bearer <token>`, die zwei Token-Klassen
und die Statuscodes, den gRPC-Server-Stream mit der Zehn-Feld-Nachricht, den
SSE-Endpunkt `GET /changes/stream` — und die Konfiguration
(`CDC_HTTP_ADDR`, `CDC_GRPC_ADDR`, `CDC_API_TOKEN_READER`/`ADMIN`). Was fehlt,
ist ein **vorzeigbares Programm**: etwas, das ein Integrator lesen, kopieren
und starten kann.

**Die Anforderung schneidet zwei Programmkassen, nicht eine.** Ein Beispiel
ist kein Belegträger mit anderem Namen. Die beiden haben verschiedene
Optimierungsziele: der Wegwerf-Client ist auf *Ausgewertetwerden* optimiert
(Maschinen-protokoll auf stdout, Negativpfad-Assertions, Nicht-Null-Exit bei
Zeitüberschreitung), das Beispiel auf *Gelesen- und Nachgebautwerden*.

**Drei Bindungen prägen den Lösungsraum.**

1. **Ein Beispiel, das `internal/` importiert, ist für seinen Leser
   unbrauchbar.** Go verbietet `internal/`-Importe außerhalb des Baums, der
   das `internal/`-Verzeichnis trägt — die interne Regel greift auf den
   **Import-Pfad** zu, nicht auf die Modulgrenze. Ein solches Programm
   kompiliert in diesem Repo und ist für einen externen Integrator nicht
   übersetzbar; es wäre ein Vorbild nur dem Anschein nach.
2. **Der gRPC-Baustein ist heute nicht öffentlich.** Der Go-Code aus
   `proto/cdc/stream/v1/changestream.proto` liegt als
   `internal/adapters/driving/grpc/streamv1/` unter `internal/` — die
   **Quelle** (`.proto`) ist öffentlich, das daraus erzeugte Go-Binding nicht.
   Ein Go-gRPC-Beispiel kann es deshalb nicht importieren.
3. **Ein Pfad-Bereich mit Import-Berechtigung in `.a-check.yml` ist
   ADR-pflichtig.** Nach [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
   §Entscheidung Festlegung 3 braucht eine ADR, wer einen Pfad-Bereich, den
   `spec/architecture.md §2` nicht als Komponente führt, über `layers`/`edges`
   oder `composition_root` **mit einer Import-Berechtigung** aufnimmt. Diese
   ADR ist der geforderte Träger für `examples/**` und für den öffentlichen
   Vertrags-Pfad.

**Und eine bereits verkörperte Regel greift.** Das Beobachtungs-Register führt
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` als
**verkörpert**: Wächst eine benutzer- oder betreiber-sichtbare Oberfläche,
ohne dass das Handbuch im selben Zug mitgezogen wird, entsteht eine Lücke, die
kein Sensor fängt. Für diese ADR heißt das: Ein Beispiel ohne seine
Handbuch-Zeile ist nicht „später dokumentiert", sondern undokumentiert.

**Was diese Entscheidung zusätzlich berührt — real gemessen in diesem Zug.**
Der erzeugte Vertrags-Code liegt heute in der **Messfläche** des
Coverage-Gates (`Dockerfile`, Stufe `coverage`: `go list ./internal/...
./cmd/...`, [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md));
`harness/sensors/coverage-gate.md` nennt
`internal/adapters/driving/grpc/streamv1` ausdrücklich als Beispiel eines
Pakets ohne eigene Testdatei, das „über fremde Testpakete gedeckt" ist. Ein
Umzug dieses Pakets aus `internal/` heraus verkürzt diese Fläche **still** —
die Zahl ändert sich, ohne dass jemand entschieden hätte, dass sich der
Messgegenstand ändert. Das ist keine Nebensache: sie ist der Grund, warum der
Umzug eine benannte Folgepflicht bekommt und nicht nur ein `git mv`.

## Entscheidung

Wir wählen: **`examples/` auf der Repo-Wurzel führt drei öffentliche
Beispiel-Clients — HTTP, SSE und gRPC —, die ausschließlich den öffentlichen
Draht-Vertrag benutzen, nicht unter `internal/` importieren und deshalb im
Benutzerhandbuch zitierbar sind. Der erzeugte Protobuf-Go-Code zieht dafür aus
dem Adapter-Paket an einen öffentlichen Vertrags-Pfad.**

Sieben Festlegungen:

### 1 — Ort und Form

- **`examples/` liegt auf der Repo-Wurzel**, nicht unter `tools/` und nicht
  unter `internal/` oder `docs/`: Es ist ein Erzeugnis **für externe Leser**;
  `tools/**` ist die interne Werkzeug-Gruppe ([`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)),
  `internal/**` ist privat, `docs/**` ist Doku ohne ausführbaren Code.
- **Ein Verzeichnis je Protokoll, je ein `main`-Paket:**
  `examples/http-client/` (ein echter Anfrage/Antwort-Aufruf gegen die
  Verwaltungs-API, z. B. `GET /tables` mit dem `reader`-Token),
  `examples/sse-client/` (öffnet `GET /changes/stream`, gibt jedes Event aus),
  `examples/grpc-client/` (öffnet `ChangeStream/StreamChanges`, gibt jede
  Nachricht aus). Die Namen tragen `-client`, um die Klasse zu benennen; das
  Wort *Consumer* bleibt dem Domänenbegriff (`cdc.consumer`) vorbehalten.
- **Startform (Form-Vorgabe, Detail der umsetzende Slice):** je Programm ein
  `go run ./examples/<name>`; Adresse und Token kommen aus **denselben
  Umgebungsvariablen, die das Handbuch bereits dokumentiert**
  (`CDC_HTTP_ADDR` bzw. `CDC_GRPC_ADDR`, `CDC_API_TOKEN_READER`), und lassen
  sich per Flag (`-addr`, `-token`) übersteuern. Damit zitiert das Handbuch
  **eine** Befehlsform, die innerhalb und außerhalb des Compose-Netzes
  funktioniert.
- **Jedes Beispiel trägt einen Doc-Kommentar**, der die angesprochene
  Schnittstelle und ihren Handbuch-Abschnitt, die tragende `LH-*`-Kennung und
  die entscheidende ADR nennt — und der ausspricht, dass es **nicht** der
  E2E-Belegträger ist (das sind die Clients unter `tools/harness/**`).

### 2 — Import-Grenze und ihr Träger

- **Ein Beispiel importiert ausschließlich den öffentlichen Draht-Vertrag:**
  die Go-Standardbibliothek, öffentliche Fremdmodule
  (`google.golang.org/grpc`, `google.golang.org/protobuf`) und — für gRPC —
  das öffentliche Vertrags-Paket aus Festlegung 3. **Kein** Import-Pfad
  enthält `/internal/`.
- **Begründung:** Die interne Regel macht ein `internal/`-importierendes
  Programm außerhalb dieses Repos unübersetzbar (Kontext Punkt 1). Ein
  Beispiel hat genau einen Leser — den externen Integrator; ein Programm,
  das er nicht bauen kann, verfehlt seine einzige Aufgabe.
- **Träger ist `.a-check.yml`** (Festlegung 4) — **nicht** ein zweites
  Go-Modul. Das ist die entscheidende Präzisierung, und sie ist **real
  gemessen, nicht angenommen** (dieser Zug, im Toolchain-Container):

  | Probe | Ergebnis |
  |---|---|
  | `examples/` als **eigenes Modul** mit Pfad `github.com/pt9912/pg-change-feed/examples`, Import von `github.com/pt9912/pg-change-feed/internal/secret` | **erlaubt** (Exit 0) — die interne Regel kennt den **Import-Pfad**, nicht die Modulgrenze; der Pfad liegt weiter unter dem `internal/`-Wurzelbaum |
  | dasselbe, Modulpfad `example.com/cdc-consumer` (außerhalb des Repo-Präfix) | `use of internal package … not allowed` |
  | dasselbe, Import eines **öffentlichen** Pfades (`…/gen/streamv1`) | **erlaubt** (Exit 0) |

  Ein zweites `go.mod` allein ist also **kein** Zaun. Echte Zähne bekommt die
  Modul-Route nur über einen **fremden** Modulpfad — und der bräuchte für die
  öffentlichen Pakete dieses Repos `require`+`replace`, eine zweite `go.sum`
  und nähme die Beispiele aus dem `go test ./...` der Wurzel (Festlegung 5).
  Die Modul-Route ist damit eine **zweite** Mechanik für dieselbe Frage, die
  `.a-check.yml` bereits als Maschinenform führt
  ([`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md),
  [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)) — und
  die schwächere, weil sie nicht das ist, wonach sie aussieht.
- **Die Grenze des Trägers, benannt statt behauptet** (gleiche Arbeitsteilung
  wie `harness/sensors/a-check.md` Grenze 1): a-check fängt einen Import aus
  `examples/**` in jede **Schicht** (`domain`, `ports`, `app`, `adapters`) —
  real belegt mit einer Probe-Datei unter `examples/`:
  `wrong-direction: examples -> adapters`, Exit ≠ 0. Es fängt **nicht** einen
  Import in einen `internal/**`-Pfad außerhalb der Schichten — `internal/bootstrap/**`
  ist `composition_root`; eine Probe-Datei unter `examples/` mit Import von
  `internal/bootstrap` blieb **grün** (0 Befunde). In der Sache ist diese
  Lücke für Beispiele unerheblich; sie ist **named**, nicht still, und der
  Wächter dort ist das Review.

### 3 — Der gRPC-Baustein: das Binding zieht heraus

- **Entscheidung:** Der erzeugte Go-Code wandert von
  `internal/adapters/driving/grpc/streamv1/` nach **`gen/cdc/stream/v1/`**
  (Modul-Pfad `github.com/pt9912/pg-change-feed/gen/cdc/stream/v1`, Paket
  `streamv1`). `gen/` ist „erzeugt, nicht handgeschrieben" und spiegelt den
  `.proto`-Pfad `proto/cdc/stream/v1/`.
- **Warum nicht gRPC aus `examples/` heraushalten.** Die Alternative ist
  zulässig und billiger (siehe §Verglichene Alternativen, Option B). Sie
  verliert, weil das Handbuch **drei** Schnittstellen führt: Ein
  Beispiel-Ordner, der zwei zeigt, wäre eine halbe Antwort auf eine
  Anforderung, deren einziges Hindernis ein **reparierbarer** Artefakt-Ort
  ist.
- **Warum das Binding überhaupt öffentlich gehört.** Die `.proto`-Quelle ist
  bereits öffentlich; ein Binding, das aus einem öffentlichen Vertrag erzeugt
  wird und unter dem privaten Baum eines Adapters liegt, ist für den Leser
  eines Beispiels unerreichbar. **Und der Repo-Beschluss steht schon:**
  [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md) nennt
  den Umzug in §Verglichene Alternativen Option F als „die genaueste Antwort
  auf die Ursache" und zieht in §Re-Evaluierungs-Trigger 2 die Folge
  ausdrücklich vor: *„Der generierte Protokoll-Stub wandert aus
  `internal/adapters/**` heraus … Dann wird die Kante auf diesen Bereich
  umgestellt; die Berechtigung wird feiner, nicht breiter."* Diese ADR löst
  den dort benannten Trigger **ein** und führt die dort vorgezeichnete Folge
  aus — sie erfindet keinen neuen Pfad, sie nimmt einen vorbereiteten.
- **Was der Umzug nach sich zieht** (Folgepflichten im Einzelnen unter
  §Konsequenzen): `option go_package` im `.proto`, `make proto-generate`, der
  Import in `internal/adapters/driving/grpc/server.go` (und `server_test.go`)
  und in `tools/harness/grpcclient` sowie die Messfläche des Coverage-Gates.
- **Was der Umzug *nicht* ist:** keine Änderung an `ADR-0060`s sechs
  Teilfragen. Der Server, der Interceptor, der `ChangeStreamPort`, der
  `Broadcaster`, die Zustellsemantik und die Token-Klassen bleiben exakt wie
  dort entschieden; nur der **Ort eines Artefakts** ändert sich. Genau diese
  eine Klausel ist superseded (§Status).

### 4 — `.a-check.yml`: zwei Gate-Scope-Gruppen, vier Kanten

- **`layers`:** `examples: ["examples/**"]` und `contract: ["gen/**"]`.
  Beide sind **keine** Komponenten der §2-Sicht, sondern Gate-Scope-Gruppen —
  dieselbe Rolle, die `tooling` seit [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  trägt; sie stehen im Kommentar der Datei als solche benannt.
- **`edges`:** `{from: examples, to: contract}` (das gRPC-Beispiel),
  `{from: adapters, to: contract}` (der gRPC-Server importiert sein
  erzeugtes Binding), `{from: tooling, to: contract}` (der Wegwerf-Client
  `tools/harness/grpcclient`, umgehängt).
- **`{from: tooling, to: adapters}` wird zurückgenommen.** Ihre einzige
  Begründung war der Stub-Import ([`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  Festlegung 2); mit dem Umzug hat sie kein Objekt mehr. Das ist die
  Richtung, die derselbe Trigger 2 vorschreibt: **feiner, nicht breiter.**
- **`examples` bekommt seine Kante erst mit ihrem Objekt.** Das HTTP- und das
  SSE-Beispiel brauchen **keine** Kante (sie benutzen nur die
  Standardbibliothek); die Kante `examples → contract` entsteht mit dem
  gRPC-Beispiel, nicht vorher — dieselbe Disziplin, die
  [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md) gegen
  „Aufnahme ohne Objekt" durchgesetzt hat.
- **ADR-pflichtig: ja.** Beide Gruppen sind Pfad-Bereiche, die
  `spec/architecture.md §2` nicht als Komponenten führt, und beide werden
  **mit einer Import-Berechtigung** aufgenommen ([`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  Festlegung 3, [`AGENTS.md`](../../../AGENTS.md) §3.6). Diese ADR ist der
  Träger; die `.a-check.yml`-Änderung ist damit kein PR-Kommentar.

### 5 — Verifikation: der Kompilierpfad trägt, ein eigener Lauf-Beleg nicht

- **Der Anti-Verrottungs-Träger ist der bereits existierende Kompilierpfad.**
  Die Beispiele liegen im **Wurzelmodul**; `make test` (`go test -race ./...`)
  übersetzt `examples/**` mit und scheitert an einem Kompilierfehler — **real
  gemessen** in diesem Zug (Probe: ein absichtlich kaputtes
  `examples/http/main.go` lässt `go test ./...` mit Exit 1 abbrechen, Meldung
  `examples/http/main.go:5:27: undefined: …`; der intakte Stand läuft grün
  durch). Derselbe Lauf läuft in der CI. **Ein neues Gate wird dafür nicht
  eingeführt** — die Bindung existiert, sie muss nur nicht ausgeschlossen
  werden. Daraus folgt die Gegenpflicht: die Beispiele dürfen aus `go test
  ./...` nicht herausfallen (kein zweites Modul, kein Ausschluss).
- **Ein eigener realer Lauf gegen die Compose-Umgebung wird *nicht*
  eingeführt.** Er würde die Draht-Eigenschaft, die
  `make test-integration` bereits über die Wegwerf-Clients belegt, ein
  zweites Mal belegen — **Duplikation, kein zusätzlicher Beleg**. Die
  Belegklassen sind verschieden: der E2E-Rundlauf sagt „die Schnittstelle
  funktioniert", ein Lauf des Beispiels sagte „dieses Programm funktioniert".
  Nur die erste Aussage trägt eine Anforderung des Lastenhefts.
- **Die Beispiele sind damit kein Lauf-Beleg, sondern Doku mit
  Kompilier-Bindung** — als solches deklariert, nicht als Beleg ausgegeben.
  Sie erscheinen nicht in `docs/user/e2e-abdeckung.md`.
- **Die Grenze, benannt:** „kompiliert" ist nicht „läuft". Der Rest-Fall —
  ein Beispiel, das übersetzt, aber dessen Argumente oder Ablauf falsch sind —
  ist von keinem Sensor gedeckt. Er wird getragen von der Tatsache, dass das
  Beispiel **dieselben** Draht-Aufrufe macht wie der E2E-belegte
  Wegwerf-Client, und vom Review. Wächst ein Beispiel über einen
  Lesbarkeits-Artefakt hinaus, greift Re-Evaluierungs-Trigger 3 und schließt
  die Lücke mit einem Smoke-Lauf, statt sie weiter zu benennen.

### 6 — Verhältnis zu den Wegwerf-Clients: nebeneinander, mit Grund

- **Beide bleiben.** Die Wegwerf-Clients unter `tools/harness/**` sind
  **Orakel**: ihr stdout wird geparst, sie asserten Negativpfade
  (`REJECTED code=401`/`Unauthenticated`) und laufen bei Zeitüberschreitung
  mit Nicht-Null-Exit. Die Beispiele sind **Vorbilder**: minimal, lesbar,
  ohne Orakel-Protokoll. Zwei Optimierungsziele — zwei Erzeugnisse.
- **Zusammenlegen wäre schlechter in beide Richtungen:** ein Beispiel mit
  Orakel-Maschinerie ist ein schlechteres Vorbild; ein E2E-Client ohne
  Negativpfad und Zeitüberschreitungs-Semantik ist ein schwächerer Beleg.
  Das ist die **Begründung** für das Nebeneinander, nicht dessen Duldung.
- **Der Preis, benannt:** derselbe Draht-Aufruf existiert in zwei Programmen.
  Die Drift-Fläche ist begrenzt: beide liegen im Wurzelmodul und werden von
  `make test` übersetzt, und der Draht selbst ist einmal real belegt.

### 7 — Handbuch

- **Die Beispiele werden im Benutzerhandbuch zitiert** — in den drei
  Schnittstellen-Abschnitten (`### Zugriff über die HTTP-/JSON-API`,
  `### Zugriff über den gRPC-Change-Stream`,
  `### Zugriff über Server-Sent-Events`), je mit Programm-Pfad und
  Startbefehl. Zitiert wird das **Beispiel**, nicht der Wegwerf-Client.
- **Die Handbuch-Zeile gehört in denselben Slice wie das Beispiel.** Das ist
  die verkörperte Regel aus
  `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`; eine
  neue zitierbare Oberfläche und ihre Handbuch-Zeile landen zusammen.
- **Ein kurzer Abgrenzungssatz im Handbuch** nennt, dass die Beispiele zum
  Lesen und Nachbauen gedacht sind und die E2E-Testclients des Harness unter
  `tools/harness/` liegen und kein Vorbild sind — sonst verwechselt der
  nächste Leser die beiden Klassen.
- **Die Beispiele stehen in den Schnittstellen-Abschnitten, nicht in §4
  (Betriebs-Aufgaben).** Ihr Leser ist der Integrator, nicht der Betreiber.

### Was diese ADR nicht ändert

- **`spec/architecture.md` bleibt unberührt.** `examples/**` und `gen/**` sind
  Gate-Scope-Gruppen wie `tooling`; die §2-Komponenten und ihre Constraints
  ändern sich nicht, kein `ARC-*` kommt hinzu.
- **`internal/**` bleibt sonst unberührt** — kein Umbau am Hexagon, nur der
  Ort eines erzeugten Artefakts wandert.
- **[`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md), [`ADR-0057`](0057-http-grpc-api.md)
  und die Spec-Straten bleiben unverändert.** Diese ADR fügt kein Programm der
  Laufzeit hinzu; sie schafft ein Erzeugnis für Leser.
- **Kein neues Gate, keine Schwellen-Senkung.** Der Kompilierpfad ist
  vorhanden; `make gates` bleibt unverändert.
- **Keine neue Abdeckungs- oder Auflösungs-Ausnahme in `.a-check.yml`.** Die
  Änderung *nimmt zurück* (`tooling → adapters`) und *benennt* (`contract`),
  sie weitet nicht.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Gesamtform: wo die Beispiele leben und wer die Grenze trägt

| Option | Pro | Contra |
|---|---|---|
| A1 — nichts tun: die Wegwerf-Clients bleiben die einzigen Client-Beispiele; das Handbuch bleibt ohne zitierbares Programm | kein Eingriff; der E2E-Beleg läuft unverändert grün | die Anforderung bleibt offen; das Handbuch führt drei Schnittstellen und kein Programm, das sie benutzt — der Integrator liest die Felder-Tabelle, nicht ein laufendes Beispiel |
| A2 — Beispiele als **zweites Go-Modul** unter `examples/` (eigene `go.mod`) | suggeriert „externe Nutzer"-Kapselung | ein zweites `go.mod` ist **kein** Zaun (real gemessen, Festlegung 2): unter dem Repo-Präfix bleibt der `internal/`-Import erlaubt; echte Zähne bräuchte ein **fremder** Modulpfad plus `require`+`replace` und eine zweite `go.sum`; die Beispiele fielen aus `go test ./...` der Wurzel (Festlegung 5) und bräuchten einen eigenen Kompilier-Pfad — mehr Infrastruktur für eine schwächere Garantie als die `.a-check.yml`-Kante, die es schon gibt |
| A3 — Beispiele **nicht** als Ordner, sondern als Code-Snippets im Handbuch | kein neuer Pfad-Bereich, keine `.a-check.yml`-Änderung | Snippets sind nicht kompilierbar und verrotten stillschweigend — ein Snippet widerspricht dem Zweck („richtige Clients", die man starten kann) und dem Kriterium „ein Beispiel, das nicht läuft, ist schlechter als keines" |
| **A4 — `examples/` auf der Wurzel, im Wurzelmodul, Grenze als `.a-check.yml`-Kante (gewählt)** | Beispiele sind kompilierbar, kopierbar und im Handbuch zitierbar; die Import-Grenze liegt bei der Mechanik, die dieses Repo für Import-Rechte führt ([`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md)/[`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)); der bestehende Kompilierpfad verhindert Verrottung ohne neues Gate; die Grenze des Trägers ist gemessen und benannt | die Grenze ist nur bis zur Schicht-Ebene mechanisch (composition_root-Pfade bleiben ungeprüft — benannte Lücke); ein neuer Pfad-Bereich in `.a-check.yml` braucht diese ADR |

### B — gRPC-Baustein: mitnehmen oder heraushalten

| Option | Pro | Contra |
|---|---|---|
| B1 — generiertes Binding bleibt unter `internal/`, **kein** Go-gRPC-Beispiel; `examples/` trägt HTTP+SSE, gRPC zeigt das Handbuch über die öffentliche `.proto`-Quelle | kein supersede, keine `.a-check.yml`-Erweiterung über `examples` hinaus, kein Umzug; nichts am Coverage-Messgegenstand ändert sich; die `.proto` **ist** der öffentliche Vertrag — ein Fremdsprachen-Consumer erzeugt sein Binding ohnehin selbst | ein Beispiel-Ordner zeigt zwei von drei dokumentierten Schnittstellen; für Go-Consumer bleibt der häufigste Fall („ich will die Nachrichtentypen benutzen") unbedient, obwohl das Hindernis ein reparierbarer Artefakt-Ort ist; die Antwort auf die Auftraggeber-Frage bliebe halb |
| **B2 — Binding zieht nach `gen/cdc/stream/v1`, gRPC-Beispiel wird möglich (gewählt)** | alle drei Schnittstellen bekommen ein echtes Vorbild; der öffentliche Vertrag bekommt ein öffentliches Go-Binding; [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md) Trigger 2 war genau dafür vorbereitet und nennt den Umzug „die genaueste Antwort auf die Ursache" | supersedes je eine Klausel in [`ADR-0060`](0060-grpc-streaming-mechanismus.md) und [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md); Import- und Generierungspfade ändern sich; die Messfläche des Coverage-Gates ist berührt (Folgepflicht) — der **größere** der beiden Wege |
| B3 — das gRPC-Beispiel importiert den internen Stub, als in-Repo-Illustration gekennzeichnet | kein Umzug; das Beispiel ist im Repo lauffähig | verletzt Festlegung 2 an genau der Stelle, an der sie zählt: das Programm ist für den einzigen Leser, den es hat, nicht übersetzbar — ein Vorbild nur dem Anschein nach; verlangte eine `examples → adapters`-Kante, also die unbeschränkte Berechtigung, die [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md) gerade zurückgenommen hat |

### C — Verifikation

| Option | Pro | Contra |
|---|---|---|
| C1 — nichts (kein Beleg, nicht einmal Übersetzung) | kein Aufwand | die Beispiele verrotten still: die `.pb.go`-Typen, die Flags, die API ändern sich, und niemand merkt es — der schlechteste der Fälle |
| **C2 — der vorhandene Kompilierpfad (`go test ./...`) trägt (gewählt)** | real belegt; kein neues Gate, kein neuer Bindungsträger; greift in CI wie lokal | fängt nur Kompilier-Verrottung, nicht einen falsch verdrahteten Ablauf (benannte Grenze, Festlegung 5) |
| C3 — eigener realer Lauf gegen die Compose-Umgebung (eigenes `make`-Ziel oder Erweiterung von `make test-integration`) | deckt auch den Ablauf-Fall | dupliziert den Draht-Beleg, den `make test-integration` über die Wegwerf-Clients schon führt; ein weiterer Lauf im ohnehin langen Compose-Stack für eine Aussage, die keine Anforderung trägt — der Aufwand steht in keinem Verhältnis zum Rest-Fall |
| C4 — eigenes Gate über `examples/...` in `make gates` | sichtbar im Gate-Bündel | dasselbe Ergebnis wie C2 (Übersetzen), nur als zweiter Mechanismus an einer Stelle, die den Kompilierpfad bereits fährt — Zeremonie |

### D — Verhältnis zu den Wegwerf-Clients

| Option | Pro | Contra |
|---|---|---|
| D1 — Beispiele **ersetzen** die Wegwerf-Clients (E2E läuft künftig die Beispiele) | kein zweites Programm je Protokoll | nimmt dem E2E sein Orakel: `READY`-Protokoll, Negativpfad-Assertions und Nicht-Null-Exit bei Zeitüberschreitung müssten in das Beispiel wandern — es wäre kein Vorbild mehr; oder der Beleg verlöre diese Hälften (stille Beleg-Schwächung) |
| **D2 — beide bleiben, mit benannten Optimierungszielen (gewählt)** | Beispiel bleibt minimal und lesbar; E2E behält sein Orakel; beide sind übersetzt und der Draht ist einmal real belegt | derselbe Aufruf existiert zweimal — begrenzte Drift-Fläche (Festlegung 6) |
| D3 — Beispiel als dünne Hülle um den Wegwerf-Client | ein Programm statt zwei | der Wegwerf-Client liegt unter `tools/harness/**` und importiert den internen Stub bzw. das Werkzeug-Umfeld — die Hülle erbte die falsche Schicht und wäre für einen externen Leser wieder unbrauchbar |

## Konsequenzen

- Positiv: Die Auftraggeber-Anforderung ist eingelöst — es gibt ein
  zitierbares Programm je dokumentierter Schnittstelle, und das Handbuch kann
  es nennen, ohne die Belegträger des E2E zu verwechseln.
- Positiv: Die Import-Grenze ist **begrenzt und deklariert**, nicht
  stillgelegt: `examples` darf genau den öffentlichen Vertrag erreichen, jeder
  Import in eine Schicht ist `wrong-direction` (real belegt); `gen/**` ist als
  Gate-Scope-Gruppe sichtbar statt „Datei ohne Schicht".
- Positiv: Die Berechtigung des Werkzeug-Baums wird **feiner**: die grobe
  Kante `tooling → adapters` fällt mit ihrem Objekt, `tooling → contract`
  ersetzt sie — genau die von [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  Trigger 2 vorgezeichnete Folge.
- Positiv: Kein neues Gate, keine Schwellen-Senkung, kein Umbau am Hexagon.
  Die Verrottungs-Bindung ist der Pfad, den `make test` bereits fährt.
- Negativ: [`ADR-0060`](0060-grpc-streaming-mechanismus.md) und
  [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md) müssen
  je in einer Klausel **zusammengelesen** werden (Muster der übrigen engen
  Klausel-Korrekturen: [`ADR-0048`](0048-heartbeat-grant-korrektur-select-ergaenzung.md),
  [`ADR-0063`](0063-lh-fa-sch-003-testform-korrektur.md),
  [`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md),
  [`ADR-0070`](0070-supersede-reichweite-und-klassengrenze.md)).
- Negativ mit Grenze: Die mechanische Grenze reicht bis zur **Schicht**-Ebene.
  Ein `examples/`-Import von `internal/bootstrap` (composition_root) oder eines
  anderen schicht-fremden `internal/**`-Pfades bleibt von a-check unentdeckt
  (real gemessen, Festlegung 2). Benannt, nicht verschwiegen; der Wächter ist
  das Review.
- Negativ mit Grenze: „kompiliert" ist nicht „läuft" (Festlegung 5). Der
  Ablauf-Fall bleibt für den Sensor offen und ist an das Review gebunden, bis
  Re-Evaluierungs-Trigger 3 ihn schließt.
- Folgepflicht (Umzugs-Slice): `option go_package` in
  `proto/cdc/stream/v1/changestream.proto` auf den neuen Pfad;
  `make proto-generate` schreibt nach `gen/cdc/stream/v1/`; die Importe in
  `internal/adapters/driving/grpc/server.go`, `server_test.go` und
  `tools/harness/grpcclient/main.go` werden gezogen; die `.a-check.yml`
  erhält `contract`, die drei Kanten, und die Rücknahme `tooling → adapters`;
  der Kommentar in `.a-check.yml` und `harness/sensors/a-check.md` (Grenze 1,
  Bindung) werden nachgezogen.
- Folgepflicht (Umzugs-Slice, eigener Beleg): **die Messfläche des
  Coverage-Gates wird nachgezogen, nicht still verkürzt.** `Dockerfile` (Stufe
  `coverage`) listet `./internal/... ./cmd/...`; mit dem Umzug fällt
  `gen/cdc/stream/v1` aus dieser Liste. Der Messgegenstand
  ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md))
  bleibt derselbe, indem `./gen/...` in die `go list`-Zeile aufgenommen wird —
  sonst ändert sich die Zahl, ohne dass jemand den Gegenstand geändert hätte.
  `harness/sensors/coverage-gate.md` (der `streamv1` als Beispiel „über fremde
  Testpakete gedeckt" nennt) zieht den Pfad nach.
- Folgepflicht (Beispiel-Slices): `.a-check.yml` `examples`-Gruppe (mit der
  Kante erst am gRPC-Beispiel), die drei Programme, und je Slice die
  zugehörigen Handbuch-Abschnitte.
- Hinweis: Der Wegwerf-Client `tools/harness/grpcclient` bleibt an seinem Ort
  und behält seinen Zweck; nur sein Import zeigt nach dem Umzug auf den
  öffentlichen Vertrags-Pfad.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check (Modulzeile `wrong-direction`) | Die Gruppe `examples` (`examples/**`) erreicht **nur** `contract`: ein Import aus `domain`, `ports`, `app` oder `adapters` in `examples/**` ist `wrong-direction` (real belegt: `wrong-direction: examples -> adapters`, Exit ≠ 0). Die Gruppe `contract` (`gen/**`) wird von `adapters`, `tooling` und `examples` erreicht | `make a-check` (im Gate-Bündel) |
| Go-Toolchain (`go test ./...`, im gepinnten Container) | `examples/**` liegt im Wurzelmodul und wird von `make test` übersetzt: ein Kompilierfehler in einem Beispiel bricht den Lauf (real belegt: Exit 1, `examples/http/main.go:5:27: undefined: …`). Kein eigenes Gate — derselbe Kompilierpfad, den alle Pakete fahren | `make test` |
| Review-Prüfpflicht (nicht maschinell) | Ein `examples/`-Import in einen `internal/**`-Pfad **außerhalb** der Schichten (namentlich `internal/bootstrap/**`, `composition_root`) ist von a-check nicht fassbar (real gemessen: 0 Befunde). Ebenso: „kompiliert" ist nicht „läuft" — ein Ablauf-Fehler eines Beispiels ohne Kompilierfehler hat keinen Sensor | — |

## Slice-Schnitt-Empfehlung

Regeln dieser Sektion: Empfehlung, endgültiger Schnitt liegt bei der
Planner-Rolle (Baseline-Regelwerk `modul-05-planning-harness.md`).

Drei Slices, in dieser Reihenfolge. Der Umzug zuerst — er ist die
Voraussetzung des gRPC-Beispiels und der riskanteste Teil (er berührt
ausgelieferten, E2E-belegten Code und die Messfläche eines Gates); er soll
allein stehen.

1. **Vertragsfläche zieht um** (kein benutzer-sichtbares Erzeugnis, keine
   Handbuch-Zeile). Inhalt: `gen/cdc/stream/v1/` als neuer Ort des erzeugten
   Go-Codes; `option go_package` im `.proto`; `make proto-generate`; die
   Importe in `server.go`/`server_test.go`/`tools/harness/grpcclient`;
   `.a-check.yml` (`contract`-Gruppe, Kanten `adapters → contract` und
   `tooling → contract`, Rücknahme `tooling → adapters`); die
   Coverage-Messfläche (`./gen/...` in der `go list`-Zeile des `Dockerfile`,
   `harness/sensors/coverage-gate.md`). Beleg: `make test`, `make a-check`,
   `make coverage-gate`, `make test-integration` (der gRPC-Rundlauf bleibt
   unverändert grün). Drei Liefer-Punkte (Vertrags-Ort+Generierung ·
   Importe · Gate-Konfiguration), zwei Berührungsflächen (Build/Tooling und
   Adapter) — ein mechanischer Umzug, in einer Sitzung prüfbar.
2. **HTTP- und SSE-Beispiel + Handbuch.** Inhalt: `examples/http-client/`,
   `examples/sse-client/` (nur Standardbibliothek), `.a-check.yml`
   `examples`-Gruppe **ohne** Kante (ihr Objekt fehlt noch), die beiden
   Handbuch-Abschnitte samt Abgrenzungssatz. Berührt keine Produktionsschicht.
3. **gRPC-Beispiel + Handbuch.** Inhalt: `examples/grpc-client/`, die Kante
   `examples → contract` (jetzt, mit ihrem Objekt), der Handbuch-Abschnitt
   zum gRPC-Stream.

Abhängigkeit, explizit benannt: Slice 3 setzt Slice 1 voraus (ohne den
öffentlichen Vertrags-Pfad gibt es kein importierbares gRPC-Binding). Slice 2
ist von beiden unabhängig.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Vier benannte Trigger, sonst permanent:

1. **Ein Beispiel braucht einen Import unter `internal/`** (etwa für eine
   Fähigkeit ohne öffentlichen Draht-Vertrag). Dann ist der öffentliche
   Vertrag unvollständig — die Antwort ist eine **Vertrags-Erweiterung** per
   Folge-ADR, **nicht** eine Ausnahme an der `examples`-Gruppe; die Grenze
   wird nicht stillschweigend aufgeweicht.
2. **Der öffentliche Vertrag bekommt einen weiteren Konsumenten** außerhalb
   von `adapters`, `tooling` und `examples` (weitere Sprache, weiteres
   Werkzeug). Dann werden die `contract`-Kanten per Folge-ADR neu gezogen;
   sie wachsen nicht stillschweigend mit.
3. **Ein Beispiel verrottet trotz Kompilier-Bindung** — es übersetzt, bricht
   aber im Ablauf, und das fällt erst einem Leser auf. Dann wird der in
   Festlegung 5 benannte Rest-Fall **geschlossen** (eigener Smoke-Lauf mit
   eigener Bindung), nicht weiter benannt.
4. **[`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
   §Re-Evaluierungs-Trigger 1 tritt ein** (ein zweiter Werkzeug-Baum außerhalb
   `tools/harness/**` braucht denselben Vertrags-Zugriff). Dann gilt er jetzt
   für die `contract`-Gruppe; die neue Kante wird per Folge-ADR gezogen.

Sonst permanent — die Wahl „öffentliche Beispiel-Clients nur auf dem
öffentlichen Draht-Vertrag, Grenze als a-check-Kante" gilt unabhängig von der
Zahl der Beispiele und der Protokolle.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: explizite Auftraggeber-Anforderung (2026-09-15), einen `./examples`-Ordner mit im Handbuch zitierbaren Clients anzulegen. Unabhängiger Architect-Zug entscheidet Ort, Import-Grenze, gRPC-Baustein, `.a-check.yml`-Gruppen, Verifikation, Verhältnis zu den Wegwerf-Clients und Handbuch-Bindung; supersedes je eine Klausel der [`ADR-0060`](0060-grpc-streaming-mechanismus.md) und [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md) und löst deren Trigger 2 ein | [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md) §Re-Evaluierungs-Trigger 2 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0076` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
