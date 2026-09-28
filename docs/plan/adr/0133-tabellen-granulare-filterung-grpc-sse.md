# ADR-0133: Tabellen-granulare Filterung für gRPC-Change-Stream und SSE (Supersedes ADR-0060/ADR-0061, teilw.)

**Status:** Accepted

**Datum:** 2026-09-28

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-28)

**Bezug:** [`LH-FA-SST-008`](../../../spec/lastenheft.md) (Haupt-Bezug —
Live-Streaming vollständiger Change-Inhalte), [`LH-FA-SST-006`](../../../spec/lastenheft.md)
(Boundary-Kriterium „fachliche Gleichwertigkeit über Zugriffswege" — Maßstab
für Teilfrage 4), [ADR-0060](0060-grpc-streaming-mechanismus.md) (**Supersedes,
teilweise** — ausschließlich die dort unter Konsequenzen benannte Vertagung
der Tabellen-Filterung und ihr Re-Evaluierungs-Trigger; alle sechs
Teilfragen, `ChangeStreamPort` und der `Broadcaster`-Vertrag bleiben
unverändert in Kraft), [ADR-0061](0061-http-sse-zusaetzlich-zu-grpc.md)
(**Supersedes, teilweise** — dieselbe Vertagung, jetzt für SSE; alle sechs
dortigen Teilfragen bleiben unverändert in Kraft), [ADR-0100](0100-nats-dritter-vollinhalts-zustellweg.md)
(Abgrenzung — **nicht berührt**: der NATS-Vollinhalts-Stream trägt bereits
eine tabellen-granulare Filterung über sein Subjekt-Schema, siehe Teilfrage 5),
[ADR-0056](0056-nats-tabellen-granulares-subjekt.md) (Formvorbild — der
Auftrag zitiert wörtlich, dass diese ADR die technische Machbarkeit bereits
zeigt), [ADR-0081](0081-changes-lesen-ueber-die-http-api.md) (Formvorbild —
`GET /changes`s Filterparameter `schema`/`table`, unabhängig optional,
wiederverwendet), [ADR-0034](0034-ports-nach-faehigkeiten.md) (Port-Zuschnitt
nach Fähigkeit — Grundlage dafür, dass diese Filterung keinen neuen Port
braucht)

**Schärft:** [`SPEC-020`](../../../spec/pflichtenheft.md) (Protobuf-Schema —
`StreamChangesRequest` erhält zwei neue Felder), [`SPEC-018`](../../../spec/pflichtenheft.md)
(Abschnitt `GET /changes/stream` — zwei neue, optionale Query-Parameter)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Auftraggeber (pt9912) hat am 2026-09-28 explizit den Bedarf an
tabellen-granularer Filterung für den gRPC-Change-Stream
(`ChangeStream/StreamChanges`, [`ADR-0060`](0060-grpc-streaming-mechanismus.md))
und den SSE-Stream (`GET /changes/stream`, [`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md))
benannt, mit einer wörtlichen Begründung: „Dass es technisch geht, zeigt das
eigene Repo bereits [gemeint: `ADR-0056`s NATS-Wecksignal-Subjekt-Filterung].
Dann sollten wir das auch in unsere Schnittstellen anbieten." Das ist exakt
der in `ADR-0060`s Konsequenzen benannte und in seinem Re-Evaluierungs-Trigger
vorgezeichnete Fall („Wird … tabellen-granulare Filterung des Streams …
konkret benannt: … eine eigene Folge-ADR") und `ADR-0061`s gleichlautender,
für SSE wiederholter Trigger — beide jetzt ausgelöst.

Diese ADR ist **keine vollständige Korrektur** von `ADR-0060`/`ADR-0061`: Sie
superseded ausschließlich die eine, schmale Aussage „Tabellen-granulare
Filterung ist nicht Gegenstand dieser ADR" (dort je unter Konsequenzen und im
Re-Evaluierungs-Trigger) — nicht die sechs Teilfragen jeder ADR, nicht
`ChangeStreamPort`, nicht den `Broadcaster`. Beide Textkörper bleiben
ansonsten unangetastet und in Kraft (`AGENTS.md` §3.5: eine `Accepted`-ADR
wird nicht inhaltlich überschrieben, eine Korrektur braucht `Supersedes` —
hier bewusst eng auf die eine Teilaussage begrenzt, analog zur Präzedenz-Form
der `Supers. …, teilw.`-Einträge im ADR-Index, z. B. `ADR-0113`/`ADR-0115`
gegenüber `ADR-0111`).

Vier Konstraints prägen den Lösungsraum:

- **`ChangeStreamPort` und der `Broadcaster` sind bereits protokoll-agnostisch
  und werden von drei Abonnenten geteilt** (gRPC, SSE, NATS-Vollinhalts-Stream,
  `ADR-0060`/`ADR-0061`/`ADR-0100`): `Subscribe() (<-chan *model.Change,
  func())` liefert **jedem** Abonnenten **jeden** Change, ohne Rücksicht auf
  dessen Interesse — Filterung existiert an dieser Stelle heute nicht.
- **Der NATS-Vollinhalts-Stream hat das benannte Problem bereits gelöst, nur
  auf einem anderen Weg:** Sein Publisher (`internal/adapters/driven/natsstream/`)
  veröffentlicht jeden Change auf einem eigenen, tabellen-granularen Subjekt
  `cdc.stream.<source_id>.<schema>.<table>` (`ADR-0100` Teilfrage 2) — ein
  NATS-Client filtert **broker-seitig** über seine Subjekt-Subscription (exakt
  oder per Wildcard, z. B. `cdc.stream.<source_id>.<schema>.>` für eine ganze
  Schema-Gruppe). Für gRPC und SSE gibt es keinen Broker mit Subjekt-Matching
  — die Filterung muss im eigenen Adapter-Code entstehen.
- **`LH-FA-SST-006`s Boundary-Kriterium verlangt fachliche Gleichwertigkeit
  über Zugriffswege** — ein Consumer, der über gRPC mit denselben Parametern
  filtert wie über SSE, bekommt dasselbe Ergebnis.
- **Bestehender Capture-kritischer Pfad ist unantastbar** (`ADR-0011`/`ADR-0027`)
  — dieselbe Grenze wie bei allen bisherigen Zustellwegen; diese ADR fügt
  keinen Schritt zwischen `Publish` und `COMMIT`/`ACK` ein.

## Entscheidung

Wir wählen **ein optionales, unabhängig setzbares `schema`/`table`-Filterpaar
je Stream-Verbindung, ausgewertet im jeweiligen Driving-Handler nach dem
bestehenden `Broadcaster.Subscribe()`-Aufruf — `ChangeStreamPort` und der
`Broadcaster` bleiben unverändert** — mit sechs Teilfragen.

### Teilfrage 1 — Form des Filters

| Option | Pro | Contra |
|---|---|---|
| A — ein Muster/Wildcard (z. B. `schema.*` oder ein Glob über den vollen `schema.table`-String) | maximale Flexibilität, ein Feld statt zwei | keine Präzedenz in diesem Repo für Wildcard-Matching auf Anwendungsebene (das NATS-Vorbild nutzt Broker-native Wildcards, keinen selbst geparsten Glob); eigene Parse-/Escaping-Logik ohne belegten Bedarf — der Auftrag nennt „tabellen-granular", kein Muster über mehrere Tabellen hinweg |
| B — eine Liste von `(schema, table)`-Paaren (repeated Feld) | ein Consumer könnte mehrere Tabellen über eine Verbindung beobachten | Überengineering gegenüber dem konkret benannten Bedarf; kein Akzeptanzkriterium verlangt Mehr-Tabellen-Beobachtung in einer Verbindung, und ein Consumer mit mehreren Interessen kann bereits heute mehrere Verbindungen öffnen (der `Broadcaster` verkraftet beliebig viele Abonnenten, `ADR-0060` Teilfrage 2/`ADR-0066`) — eine Liste bräuchte zusätzlich eine Leerlisten-vs.-ungesetzt-Unterscheidung, die die nächste Option nicht braucht |
| **C — genau ein optionales `(schema, table)`-Paar, beide Felder unabhängig optional (gewählt)** | deckt den benannten Bedarf minimal; **identische Form wie `ADR-0081`s bereits `Accepted`-Filter für `GET /changes`** (`schema`, `table`, unabhängig optional) — ein Consumer, der beide Lesewege kennt, muss keine zweite Filtersemantik lernen; ein leerer/nicht gesetzter Wert je Feld ist unzweideutig „kein Filter auf dieser Dimension", keine Sonderform für „alles" nötig | ein Consumer mit Interesse an mehreren, nicht durch ein gemeinsames Schema abgedeckten Tabellen braucht mehrere Verbindungen — bewusst in Kauf genommen (siehe Contra von Option B) |

Ein gesetztes `schema` ohne `table` filtert auf alle Tabellen dieses Schemas;
ein gesetztes `table` ohne `schema` filtert auf jede Tabelle dieses Namens,
unabhängig vom Schema; beide gesetzt filtert exakt eine Tabelle; keines
gesetzt liefert wie bisher alle Changes der konfigurierten Tabellen —
dieselbe Kombinatorik wie `ADR-0081`s `GET /changes`.

### Teilfrage 2 — Wire-Kompatibilität (Protobuf)

`proto/cdc/stream/v1/changestream.proto` ist seit seiner Einführung
(`ADR-0060`) unverändert geblieben — `ADR-0130` Teilfrage 3/Option B und
`ADR-0131` haben ihre jeweils neuen RPCs bewusst in eine **eigene** `.proto`-
Datei (`administration.proto`) gelegt, damit `changestream.proto` „byte-
identisch unverändert" bleibt (`ADR-0130` Teilfrage 3). Diese ADR ist die
**erste** inhaltliche Änderung an `changestream.proto` seit seiner Einführung.

| Option | Pro | Contra |
|---|---|---|
| A — neue RPC-Methode `StreamChangesFiltered` neben der bestehenden `StreamChanges` | `StreamChanges` selbst bliebe byte-identisch, kein Feld-Zusatz an der bestehenden Request-Message | verdoppelt Service-Definition, Server-Implementierung und Testfläche für eine Fähigkeit, die sich als zwei zusätzliche, additive Felder in der bestehenden Request ausdrücken lässt; ein alter Client, der `StreamChanges` weiter aufruft, bekäme nie Filterung angeboten, obwohl er dafür nicht einmal ein neues Feld setzen müsste |
| **B — zwei neue, optionale `string`-Felder `schema`/`table` in der bestehenden `StreamChangesRequest` (gewählt)** | additiv: Feldnummern 1/2 einer bislang leeren Message (`message StreamChangesRequest {}`) sind rückwärtskompatibel — ein alter, gegen die alte `.proto` kompilierter Client sendet weiterhin eine leere `StreamChangesRequest{}`; der neue Server liest dafür `schema=""`/`table=""` (proto3-Zero-Value) und behandelt das exakt wie „kein Filter" (Teilfrage 1) — bestehendes Verhalten bleibt bit-identisch ohne jede Fallunterscheidung im Server; ein neuer Client setzt die Felder optional | ein alter, noch nicht neu generierter Server-Build ignoriert die Felder eines neuen Clients (kein Schutz gegen Client/Server-Versions-Schere) — dieselbe, bereits an anderer Stelle in Kauf genommene Eigenschaft additiver Proto-Felder |
| C — `optional string schema = 1;`/`optional string table = 2;` mit explizitem Presence-Tracking (proto3 `optional`-Keyword) | unterscheidet „nicht gesetzt" von „explizit leerer String" | kein Akzeptanzkriterium verlangt diese Unterscheidung — ein leerer String als Filterwert wäre ohnehin sinnlos (keine Tabelle heißt ""); zusätzliche Codegen-Komplexität (Wrapper-Typen/Presence-Getter) ohne Gegenwert gegenüber Option B |

`StreamChangesRequest` trägt künftig:

```proto
message StreamChangesRequest {
  string schema = 1;
  string table = 2;
}
```

### Teilfrage 3 — Ansatzpunkt der Filterung

| Option | Pro | Contra |
|---|---|---|
| A — `Broadcaster.Subscribe()` erhält einen Filter-Parameter, filtert selbst vor der Zustellung in die Empfangs-Warteschlange | spart die (geringe) In-Prozess-Kosten unnötiger Kanal-Zustellungen | ändert die Signatur eines von **drei** bereits `Accepted`-ADRs geteilten, lokal deklarierten Interface-Vertrags (`changeSubscriber` in `grpc/server.go`, `http/sse.go`, `natsstream/publisher.go`) und der konkreten `Broadcaster.Subscribe()`-Methode — der NATS-Publisher bräuchte dann einen „Nicht-Filter"-Aufruf, obwohl er (Teilfrage 5) gar keine Anwendungsfilterung braucht; berührt drei Pakete für einen Bedarf, der nur zwei betrifft |
| **B — Filterung im jeweiligen Driving-Handler, direkt vor dem Versand an genau diesen Abonnenten (gewählt)** | `ChangeStreamPort`, `Broadcaster` und `Subscribe()` bleiben **byte-identisch unverändert** — keine der drei bereits `Accepted`-ADRs (`ADR-0060`/`ADR-0061`/`ADR-0100`) wird in ihrem Port-/Broadcaster-Design berührt; der Filter-Vergleich (zwei String-Vergleiche) ist trivial und läuft ohnehin schon in der bestehenden `for { select { … } }`-Empfangsschleife jedes Handlers (`StreamChanges` in `grpc/server.go`, `streamChangesHandler` in `http/sse.go`); ein nicht passender Change wird verworfen, **bevor** er über das Netz an diesen Client geht — „serverseitig vor dem Versand", nicht client-seitiger Nachfilter | jeder der beiden Handler bekommt weiterhin **jeden** Change über seinen Kanal zugestellt (der Broadcaster fan-out bleibt ungefiltert) — eine vernachlässigbare zusätzliche In-Prozess-Zustellung je Change und Abonnent, kein Netzwerk-Overhead |
| C — ein eigener, filternder Wrapper-Kanal zwischen `Broadcaster` und Handler (Adapter-Pattern, kein Signatur-Bruch an `Broadcaster` selbst) | vermeidet Signatur-Bruch wie Option B | zusätzliche Indirektion (ein weiterer Goroutine-/Kanal-Layer je Verbindung) ohne Vorteil gegenüber der einfachen Inline-Prüfung in der bereits bestehenden Empfangsschleife — löst dasselbe Problem mit mehr Code |

Die Filterprüfung ist eine reine, seiteneffektfreie Funktion auf zwei
String-Paaren (`change.Schema`/`change.Table` gegen die beiden
Request-/Query-Werte) — identisch für gRPC und SSE einsetzbar, damit sie
`LH-FA-SST-006`s Gleichwertigkeits-Kriterium (Teilfrage 4) strukturell
erfüllt und nicht zweimal unabhängig nachgebaut wird.

### Teilfrage 4 — Verhältnis zu SSE (fachliche Gleichwertigkeit)

| Option | Pro | Contra |
|---|---|---|
| A — SSE bekommt eine andere Parameterform als gRPC (z. B. ein kombinierter `table=schema.table`-Parameter) | könnte „kürzer" wirken | verletzt `LH-FA-SST-006`s Boundary-Kriterium ohne Not — ein Consumer, der zwischen gRPC und SSE wechselt, müsste zwei verschiedene Filtersyntaxen lernen, obwohl beide dieselbe zugrunde liegende Prüfung ausführen |
| **B — SSE nutzt dieselben Feldnamen `schema`/`table` als Query-Parameter (`GET /changes/stream?schema=<schema>&table=<table>`), dieselbe Kombinatorik wie gRPC und wie `GET /changes` (gewählt)** | ein Consumer, der über gRPC mit `schema`/`table` filtert, bekommt exakt dieselbe Ergebnismenge wie über SSE mit denselben Werten — `LH-FA-SST-006` erfüllt; dieselbe Prüf-Funktion (Teilfrage 3) hinter beiden Handlern; keine dritte Filtersyntax neben `GET /changes`s bereits etablierter Form | keine wesentlichen |
| C — SSE bekommt gar keinen Filter, nur gRPC | halber Implementierungsaufwand | der Auftrag benennt ausdrücklich beide Schnittstellen („unsere Schnittstellen", Plural); verletzt `LH-FA-SST-006` für den Fall, dass ein Consumer beide Wege nutzt |

Ein unbekannter Query-Parameter am `GET /changes/stream`-Endpunkt wird wie
bei `GET /changes` (`ADR-0081`) mit `400` abgelehnt, **vor** Öffnen des
Streams (derselbe Prüfzeitpunkt wie der bestehende `503`-Check bei fehlendem
Broadcaster) — Konsistenz mit dem bereits etablierten Fehlerverhalten
desselben Adapters.

### Teilfrage 5 — Scope-Entscheid: NATS-Vollinhalts-Stream

**Explizit nicht Gegenstand dieser ADR.** Der NATS-Vollinhalts-Stream
(`ADR-0100`) veröffentlicht bereits auf einem tabellen-granularen Subjekt
`cdc.stream.<source_id>.<schema>.<table>` — ein NATS-Client filtert
broker-seitig über seine Subscription (exakt oder Wildcard), ohne dass der
`natsstream.Publisher` irgendeine Anwendungslogik dafür braucht. Der von
dieser ADR gelöste Bedarf — „gRPC und SSE kennen keinen Broker mit
Subjekt-Matching, also fehlt ihnen eine Filter-Fähigkeit, die NATS
strukturell bereits hat" — existiert für NATS nicht. Diese Feststellung ist
**hergeleitet** aus `ADR-0100` Teilfrage 2 (dort bereits `Accepted`
beschrieben), nicht neu geprüft: kein Zeilen-Diff an
`internal/adapters/driven/natsstream/` ist Folge dieser ADR, und `ADR-0100`
wird durch diese ADR **nicht** superseded.

### Teilfrage 6 — Fehlerform bei einem nicht existierenden Schema/Tabellen-Paar

| Option | Pro | Contra |
|---|---|---|
| A — Existenzprüfung gegen die konfigurierten Tabellen bei Verbindungsaufbau, Ablehnung bei Nichttreffer | „schnelles Scheitern" für einen offensichtlichen Tippfehler | bräuchte einen Store-/Konfigurations-Lookup an einer Stelle, die heute rein aus dem Request-Objekt und dem Broadcaster-Kanal liest — neue Abhängigkeit, die kein Akzeptanzkriterium verlangt; zudem inkonsistent zu `GET /changes`s bereits etabliertem Verhalten (Option B) |
| **B — kein Fehler; ein Filter ohne jemals passenden Change liefert dauerhaft keine Events (gewählt)** | konsistent mit `ADR-0081`s bereits `Accepted`-Verhalten für `GET /changes` („kein Treffer ⇒ leere, gesetzte Liste statt eines Fehlers"); ein Stream ist ohnehin ein unbegrenzt offener, ereignisgetriebener Kanal — „aktuell keine Treffer" und „wird nie Treffer haben" sind für den Server ununterscheidbar, ohne die vollständige, sich ändernde Tabellenkonfiguration zu kennen; kein zusätzlicher Code-Pfad | ein Consumer mit einem Tippfehler im Filter bemerkt das nicht über einen Fehler, sondern über Ausbleiben von Events — dieselbe Grenze wie bei jedem Fire-and-Forget-Stream ohne Zustellgarantie (`ADR-0060` Teilfrage 3) |

## Konsequenzen

- Positiv: `LH-FA-SST-008` und der explizit benannte Bedarf werden für gRPC
  und SSE erfüllbar, ohne `ChangeStreamPort`, den `Broadcaster` oder
  `ADR-0100`s NATS-Weg anzufassen — die Filterung bleibt lokal in den zwei
  betroffenen Driving-Handlern.
- Positiv: Dieselbe Filtersyntax (`schema`/`table`, unabhängig optional) über
  drei Lesewege hinweg (`GET /changes`, gRPC-Stream, SSE-Stream) — ein
  Consumer lernt eine Form, nicht drei.
- Positiv: `LH-FA-SST-006`s Gleichwertigkeits-Kriterium wird strukturell
  gestützt, weil beide Handler dieselbe Prüf-Funktion verwenden können
  (Folgepflicht des umsetzenden Slices, hier architektonisch festgelegt).
- Negativ: Der `Broadcaster`-Fan-out bleibt ungefiltert — jeder Abonnent
  bekommt weiterhin jeden Change über seinen In-Prozess-Kanal zugestellt und
  verwirft ihn ggf. selbst; bei sehr vielen gleichzeitigen, eng gefilterten
  Verbindungen ist das eine (nicht gemessene, für die heutige Größenordnung
  als vernachlässigbar eingeschätzte) zusätzliche CPU-Last gegenüber einer
  Broadcaster-seitigen Filterung — bewusst in Kauf genommen für den Erhalt
  des unveränderten Port-/Broadcaster-Vertrags (Teilfrage 3).
- Negativ: Ein Consumer mit einem nicht (mehr) existierenden Schema-/
  Tabellennamen im Filter bemerkt den Fehler nicht über eine Ablehnung,
  sondern über das dauerhafte Ausbleiben von Events (Teilfrage 6) —
  dieselbe, bereits an `GET /changes` akzeptierte Grenze.
- Folgepflicht: Der umsetzende Slice erweitert `changestream.proto` um die
  zwei Felder (Teilfrage 2), implementiert die Filterprüfung in
  `internal/adapters/driving/grpc/server.go` (`StreamChanges`) und
  `internal/adapters/driving/http/sse.go` (`streamChangesHandler`,
  Query-Parameter-Parsing samt `400`-Pfad für unbekannte Parameter,
  Teilfrage 4), und teilt die Prüf-Funktion zwischen beiden Handlern (z. B.
  über ein kleines, gemeinsames internes Paket oder eine an beiden Stellen
  identisch gehaltene Funktion — Wahl liegt beim umsetzenden Slice).
- Folgepflicht: `spec/pflichtenheft.md` (`SPEC-020`, `SPEC-018`) wird um die
  beiden neuen Felder/Parameter samt Kombinatorik (Teilfrage 1) ergänzt.
- Folgepflicht: `docs/user/benutzerhandbuch.md` — die beiden Sätze „Eine
  Filterung nach Tabelle ist nicht Teil dieser Version" (§4 „Zugriff über den
  gRPC-Change-Stream") und die entsprechende Aussage im SSE-Abschnitt werden
  durch die neue Filterbeschreibung ersetzt (Ist-Zustand, keine
  Chronik-Formulierung).
- Folgepflicht: Die drei Beispiel-Clients je Sprache
  (`examples/grpc-client`, `examples/csharp/grpc-client`,
  `examples/kotlin/grpc-client`; `examples/sse-client`,
  `examples/csharp/sse-client`, `examples/kotlin/sse-client`) und die drei
  SDK-Packages (`PgChangeFeed.Client`, `pgchangefeed`, `pgchangefeed-kotlin`)
  erhalten optionale Filter-Parameter an ihren jeweiligen
  `StreamChanges`-Aufrufen — je eigener Folge-Schritt, nicht Gegenstand
  dieser ADR.
- Folgepflicht: `make test-integration` erhält einen Filter-Beleg je Weg
  (gRPC und SSE): zwei Tabellen aktiv, ein gefilterter Consumer sieht nur die
  passenden Changes.
- Folgepflicht: `spec/architecture.md` braucht **keine** Änderung — `ARC-005`
  bleibt unverändert generisch; diese ADR verfeinert ein bereits konkretes
  Protokoll-Detail, ohne die Sicht selbst zu ändern.
- Folgepflicht: `.a-check.yml` braucht **keine** Änderung — keine neue Datei
  außerhalb der bestehenden Glob-Muster, kein neuer Hexagon-Schichten-Edge.

## Fitness Function (falls maschinell prüfbar)

Regeln dieser Sektion: Erwartung an den umsetzenden Slice, nicht selbst
gefahren (`AGENTS.md` §3.12 — der Architect verfasst, er testet nicht).

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driving/grpc`, Whitebox, erwartet) | Ein `StreamChangesRequest` mit gesetztem `schema`/`table` liefert nur Changes der passenden Kombination; eine leere `StreamChangesRequest{}` liefert weiterhin alle Changes (Regressionstest gegen bestehende, ungefilterte Clients) | `make test` |
| Go-Unit-Test (`internal/adapters/driving/http`, Whitebox, erwartet) | `GET /changes/stream?schema=…&table=…` liefert nur Changes der passenden Kombination; ohne Query-Parameter liefert der Endpunkt weiterhin alle Changes; ein unbekannter Query-Parameter endet mit `400`, vor jedem SSE-Event | `make test` |
| Go-Unit-Test (gemeinsame Prüf-Funktion, falls extrahiert, erwartet) | Dieselbe Prüf-Funktion liefert für dasselbe `(schema, table)`-Paar und denselben Change dasselbe Ergebnis, unabhängig vom Aufrufer — Regressionstest gegen eine unabhängig voneinander driftende Doppelimplementierung (`LH-FA-SST-006`) | `make test` |
| `.a-check` (erwartet) | unverändert — keine neue Datei außerhalb der bestehenden Glob-Muster | `make a-check` |
| `make test-integration` (erwartet) | ein Filter-Rundlauf je Weg (gRPC, SSE): zwei aktivierte Tabellen, ein gefilterter Consumer sieht nur die passende, ein ungefilterter Consumer weiterhin beide | `make test-integration` |

## Slice-Schnitt-Empfehlung

Regeln dieser Sektion: Empfehlung, endgültiger Schnitt liegt bei der
Planner-Rolle (Baseline-Regelwerk `modul-05-planning-harness.md`).

Ein einzelner Slice genügt: Die Filterung baut vollständig auf bereits
bestehender Infrastruktur auf (`ChangeStreamPort`/`Broadcaster` unverändert,
beide Driving-Adapter existieren bereits) und bringt keine neue
Toolchain-Abhängigkeit mit. Inhalt: `changestream.proto`-Erweiterung
(Teilfrage 2), Filterprüfung in beiden Handlern (Teilfrage 3), SSE-
Query-Parameter samt `400`-Pfad (Teilfrage 4), Unit-Tests, `SPEC-020`/
`SPEC-018`-Nachzug, Handbuch-Korrektur der beiden „nicht Teil dieser
Version"-Sätze, `make test-integration`-Erweiterung. Beispiel-Clients und
SDK-Erweiterung können als eigener Folge-Slice laufen (analog zur
bestehenden Praxis, z. B. `ADR-0131`s `ReadChanges`, dessen
Beispiel-/SDK-Aufnahme laut Handbuch „ein eigener, noch nicht terminierter
Folge-Schritt" bleibt).

## Re-Evaluierungs-Trigger

Wird ein Bedarf an Mehr-Tabellen-Filterung über eine einzelne Verbindung
(Teilfrage 1 Option B) oder an Broadcaster-seitiger statt Handler-seitiger
Filterung (Teilfrage 3 Option A, z. B. bei belegtem Performance-Druck durch
sehr viele eng gefilterte Verbindungen) konkret benannt: Beide widersprechen
der hier getroffenen Wahl fundamental genug, dass sie eine eigene Folge-ADR
brauchen, die die jeweilige Option neu bewertet — nicht als Korrektur dieser
ADR, sondern als eigene Entscheidung, weil sich die Prämisse (ein Filterpaar
genügt dem konkret benannten Bedarf; die In-Prozess-Kosten ungefilterter
Broadcaster-Zustellung sind vernachlässigbar) geändert hätte. Sonst
permanent — die Wahl eines optionalen `schema`/`table`-Filterpaars,
ausgewertet im jeweiligen Driving-Handler, gilt unabhängig vom Zeitpunkt der
Umsetzung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-28 | Accepted — Architect-Entscheidung nach explizitem Auftraggeber-Bedarf (2026-09-28); superseded `ADR-0060`/`ADR-0061` ausschließlich in der dort benannten Vertagung der Tabellen-Filterung, alle übrigen Teilfragen beider ADRs bleiben unverändert in Kraft; `ADR-0100` wird nicht berührt | [`ADR-0060`](0060-grpc-streaming-mechanismus.md), [`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0133` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
