# Slice routing-lesewege: Lesewege — Parameter `target` an `GET /changes`, am gRPC-Stream und am SSE-Stream

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](../welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-SST-006`](../../../../spec/lastenheft.md) (Gleichwertigkeit
der Zugriffswege), [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Streaming),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 4 und Teilfrage 5 (ein Parameter je Weg),
[`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
Festlegung 1 (`ReadChangesRequest.target = 7`),
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md)
(Formvorbild: optionale Filterparameter `schema`/`table`),
[`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md) und
[`ADR-0131`](../../adr/0131-grpc-readchanges-zehnter-rpc.md) (`GET /changes`,
gRPC-`ReadChanges`), [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md),
[`ADR-0061`](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md) (Zustellwege).

**Berührte Spec-Stellen:**
[`SPEC-020`](../../../../spec/pflichtenheft.md) (gRPC-Request),
[`SPEC-021`](../../../../spec/pflichtenheft.md) (SSE-Query-Parameter),
[`SPEC-022`](../../../../spec/pflichtenheft.md) (`GET /changes`),
[`SPEC-031`](../../../../spec/pflichtenheft.md) (Zeile `ReadChanges`, Feld `target`).
Die Spec führt: der Slice setzt `slice-routing-spec-nachzug` voraus, der auch die
`SPEC-031`-Zeile trägt (`ADR-0138` ist `Accepted`); dieser Slice ändert die Spec nicht.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Leser wählt das Ziel an den drei Wegen, die der Server mit Parametern
bedient: `GET /changes` (Query-Parameter `target`), der gRPC-Stream
(`StreamChangesRequest` Feld 3 `string target`) und der SSE-Stream (Query-Parameter
`target`); leer heißt kein Filter, `target` kombiniert sich mit `schema`/`table` als
Konjunktion, ein ungefilterter Leser sieht weiterhin alle Changes. Drei Liefer-Punkte:

- (A) **`GET /changes`:** `ReadChangesQuery` erhält das Ziel; der Store filtert
  (`WHERE route_target = $n` in der Lese-Anweisung, reine Auswahl nach
  [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)); unbekannte
  Parameter bleiben `400`; das Nachrichtenschema (dreizehn Felder) bleibt;
  der gRPC-RPC `ReadChanges` trägt `string target = 7` in `ReadChangesRequest`
  (`ADR-0138` Festlegung 1; derselbe Use Case, `ChangeRecord` bleibt dreizehnfeldrig);
- (B) **gRPC-Stream:** das Proto-Feld, die Generierung (`make proto-generate`,
  `make generated-sync`), der Handler filtert über die gemeinsame Funktion;
- (C) **SSE-Stream:** der Query-Parameter mit demselben `400`-Pfad, derselbe Filter;
  die **gemeinsame Filterfunktion** (`Change.MatchesFilter`) trägt das Ziel für beide
  Stream-Handler.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`ChangeStreamPort` und `Broadcaster`** — bleiben unverändert
  ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Teilfrage 5): die Filterprüfung ist eine reine Funktion im jeweiligen Handler;
  ein Filter im Broadcaster wäre ein Eingriff in drei `Accepted` Zustellwege.
- **NATS-Vollinhalt** — `slice-routing-nats-subjekt`: dort gilt ein zusätzliches
  Subjekt statt eines Parameters.
- **Das Wecksignal** — bleibt leer und ohne Ziel
  ([`SPEC-017`](../../../../spec/pflichtenheft.md), Teilfrage 5 der ADR).
- **Das Label im Nachrichtenschema der Wege** — bleibt außerhalb (zehn bzw. dreizehn
  Felder unverändert); der Leser sieht das Label über die SQL-Sicht.
- **SDK- und Beispiel-Clients** — `slice-routing-sdk-beispiel-target`: die drei
  SDK-Packages bauen die gRPC-Stubs aus der `.proto` im eigenen Bau; die Änderung des
  Proto-Felds ist hier, die Nutzung dort.
- **Der Beleg am laufenden System** — `slice-routing-e2e`; dieser Slice belegt auf
  Unit- und Handler-Ebene (`make test`) und die Proto-Synchronität
  (`make generated-sync`).

## 2. Definition of Done

- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-006`](../../../../spec/lastenheft.md) (A): `GET /changes?...&target=<ziel>`
      liefert genau die Changes im Bereich, deren `route_target` dem Ziel entspricht;
      ohne Parameter liefert der Endpunkt unverändert alle (Regressionstest);
      `target` kombiniert sich mit `schema`/`table` als Konjunktion; ein Parameter
      außerhalb der Menge endet `400`; das Verhalten bei einem syntaktisch
      ungültigen oder nie vergebenen Ziel folgt der Spec (`slice-routing-spec-nachzug`).
      *Zu belegen durch:* Handler-Tabellentest, Store-Test des Lesens mit Filter
      (`make test`, `make test-store`).
- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (B, C): der gRPC-Stream und
      der SSE-Stream liefern mit `target` nur Changes mit diesem Ziel, ohne `target`
      jede Change (geroutete eingeschlossen); `target` kombiniert sich mit `schema`/`table`;
      der SSE-Parameter außerhalb der Menge endet `400` vor jedem Event; die
      Filterprüfung ist **eine** reine Funktion (`Change.MatchesFilter`), beide
      Handler rufen sie; `make generated-sync` grün (das Proto-Feld ist in
      `gen/cdc/stream/` erzeugt und committet); der Bestand der Aufrufer von
      `MatchesFilter` (Parent: vier Nicht-Test-Zeilen) trägt das Ziel durchgängig.
      *Zu belegen durch:* Handler-Tests mit Fake-Stream (`make test`),
      `make generated-sync`.
- [ ] **Vorab-Bedingung V1** (Welle §5) ist durch
      [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
      Festlegung 1 entschieden (erfüllt); der Slice setzt sie um: `ReadChangesRequest`
      trägt `string target = 7` (leer = kein Filter, Konjunktion mit `schema`/`table`),
      der Use Case `ReadChangesUseCase` erhält den Filter einmal und bedient
      `GET /changes` und den RPC; die Aussage „der RPC sieht dasselbe wie
      `GET /changes`" ist am Test belegt (gleiche Eingabe, gleiche Changes); ein
      `target` außerhalb des Alphabets (Großbuchstabe, 64 Zeichen, U+0000) liefert auf
      allen Lesewegen eine leere Antwort, keinen Fehler; Festlegung nach
      [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md):
      der gemeinsame Use Case `ReadChangesUseCase` prüft das Alphabet mit der
      Prüffunktion der Domäne und antwortet leer, ohne den Store aufzurufen; Tests:
      Use Case mit Store-Fake, der bei einem Aufruf fehlschlägt (`make test`), `GET
      /changes` und gRPC-`ReadChanges` mit derselben Eingabe, gRPC-Stream und SSE
      liefern keine Nachricht; der SQL-Zugriff `cdc.changes` ist ausgenommen;
      `make generated-sync`
      grün für beide `.proto`-Dateien, `tools/harness/grpcadminclient` trägt `target`.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-lesewege.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `spec/pflichtenheft.md` entfällt — die `SPEC-031`-Zeile
      `ReadChanges` mit dem Feld `target` trägt `slice-routing-spec-nachzug`
      (`ADR-0138` Folgepflicht); das Benutzerhandbuch (Abschnitte „Zugriff über die
      HTTP-/JSON-API", „…gRPC-Change-Stream", „…gRPC-Verwaltungs-API", „…Server-Sent-Events")
      bleibt unberührt — Adresse: `slice-routing-betriebsdoku` §2 (Parameter `target`
      je Weg, Konjunktion mit `schema`/`table`, die Beispiele der Clients folgen
      `slice-routing-sdk-beispiel-target`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/model/change.go` (`MatchesFilter`) | update | die gemeinsame Filterfunktion trägt das Ziel; die Form (Signatur oder Filter-Wert) entscheidet der Implementer, die Aufrufer ziehen mit (Parent: vier Nicht-Test-Aufrufe, `grpc/server.go`, `http/sse.go`, `change.go`). |
| `internal/application/port/inbound/` (`ReadChangesQuery`), `internal/application/usecase/readchanges/` | update | das Ziel im Query-Typ; der Use Case reicht es an den Store. |
| `internal/application/port/outbound/` (Lese-Port des Stores), `internal/adapters/driven/postgresstorage/queries/queries.go` (`SelectChanges`) | update | optionale Auswahl nach `route_target` (Gleichheit, kein `LIKE`). |
| `internal/adapters/driving/http/readchanges.go` | update | Query-Parameter `target`, die strenge Parameter-Menge (unbekannte → `400`) um ein Element erweitert. |
| `internal/adapters/driving/http/sse.go` | update | Query-Parameter `target`; `parseStreamChangesFilter`, gleicher `400`-Pfad. |
| `internal/adapters/driving/grpc/server.go` | update | `StreamChangesRequest.target` in den Filter. |
| `internal/adapters/driving/grpc/administration.go` | update | `ReadChangesRequest.target` an den Use Case (`ADR-0138` Festlegung 1). |
| `tools/harness/grpcadminclient/` | update | Flag `target` für den RPC `ReadChanges` (`ADR-0138` Folgepflicht). |
| `proto/cdc/stream/v1/changestream.proto` und `proto/cdc/administration/v1/administration.proto` (`string target = 7`), `gen/cdc/**` | update (Erzeugnis) | additive Felder; `make proto-generate`, Prüfung `make generated-sync`. |
| `internal/adapters/driving/http/*_test.go`, `internal/adapters/driving/grpc/*_test.go`, `internal/domain/model/change_test.go`, `internal/adapters/driven/postgresstorage/*_test.go` | update | Happy/Boundary/Negative je Weg, Konjunktion, Regression ohne Parameter. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Stream- und
Lese-Anfragen tragen `schema` und `table` als einzige Filter"; Parent ist `30fd6cb5`;
der Implementer ergänzt die `diff`-Zeilen und trägt Gefundenes und Nichtgefundenes ein):**

```suchlauf
30fd6cb5 4 -n MatchesFilter -- internal ':!*_test.go'
30fd6cb5 10 -n MatchesFilter -- internal
30fd6cb5 8 -n -E 'StreamChangesRequest' -- internal proto spec docs/user ':!*_test.go' ':!gen'
30fd6cb5 3 -n -E 'ReadChangesRequest' -- internal proto spec ':!gen' ':!*_test.go'
30fd6cb5 10 -n -E 'Query-Parameter' -- spec docs/user
30fd6cb5 17 -n -E '`schema`/`table`|schema/table|schema und table' -- spec docs/user
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Aufrufer der Filterfunktion | Zeilen 1 und 2: 4 Nicht-Test-Zeilen, 10 mit Tests | jeder Aufrufer trägt das Ziel durchgängig; Befund am Diff: einzutragen |
| Definitionen und Aufrufer der Anfrage-Typen | Zeilen 3 und 4: 8 bzw. 3 Zeilen | jede Stelle lesen: Aufzählung der Felder, die das Ziel braucht — oder eine begründete Auslassung; Befund: einzutragen |
| Beschreibungen der Parameter in Spec und Handbuch | Zeilen 5 und 6: 10 bzw. 17 Zeilen („Query-Parameter", „`schema`/`table`") | Spec-Zeilen zieht `slice-routing-spec-nachzug` (bereits geändert, wenn dieser Slice startet); Handbuch gemeldet an `slice-routing-betriebsdoku`, Client-Beispiele an `slice-routing-sdk-beispiel-target`; Befund: einzutragen |
| Aussagen „gRPC und SSE filtern über `schema`/`table`" in `examples/` und `sdks/` | liegen außerhalb dieses Suchraums | gemeldet an `slice-routing-sdk-beispiel-target` (§2 dort nennt den Gegenstand), nicht mitgeändert |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-kern-label` liegt in `done/`
(die Spalte und das Feld `RouteTarget` an der Change), `slice-routing-spec-nachzug`
liegt in `done/` und kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Die
Frage zu V1 (Welle §5) ist mit `ADR-0138` Festlegung 1 beantwortet, kein Verdikt
steht aus. Der Slice braucht den Antragsweg
nicht: seine Tests setzen das Ziel an hand-gebauten Changes. Er steht in der Tabelle
trotzdem hinter `slice-routing-backfill-pfad` (Welle §4 Abweichung 2).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): zeigt die Proto-Generierung
  oder der Store-Lesepfad mehr Umfang als erwartet, trennt sich (B)/(C) von (A) ab
  (`slice-routing-lesewege-stream`); (A) liefert dann `GET /changes` allein.
- `in-progress` → `open` (blockiert): `make generated-sync` bleibt rot, weil der
  gepinnte Generator das Feld anders erzeugt als die committeten Dateien (Pin-Frage,
  kein Weiterbau).

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert, `generated-sync` eingeschlossen), `make test-store` real grün, Suchlauf-
Block nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Altserver ignoriert das neue Feld still.** `target` ist ein additives Proto-Feld
  und ein zusätzlicher Query-Parameter; ein neuer Client gegen einen Server ohne
  diesen Slice bekommt am gRPC-Weg ungefilterte Changes (proto3 verwirft unbekannte
  Felder — *hergeleitet*, nicht gemessen) und am HTTP-/SSE-Weg `400` (die Parameter-
  Menge ist streng, Bestand: `SPEC-021`/`SPEC-022`). — **Ausgang:** bei der Closure
  einzutragen (Aussage im Handbuch mit Ursprung; Adresse `slice-routing-betriebsdoku`).
- **Das Ziel ist Auswahl, kein Zugriffsschutz.** Jeder Leser mit `reader`-Token kann
  jedes Ziel wählen oder gar nicht filtern
  ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Entscheidung 2, ungefilterte Leser sehen alles); eine Lesart „Ziel A sieht nur A"
  als Berechtigung ist falsch. — **Ausgang:** bei der Closure einzutragen
  (Handbuch-Adresse `slice-routing-betriebsdoku`).
- **Kommt das Label an den Handler?** Der Handler filtert auf `model.Change`; ob der
  Broadcaster dieselbe Change mit gesetztem `RouteTarget` liefert, die persistiert
  wurde, ist *erwartet* (ein Objekt, ein Erzeugungspfad), nicht gelesen. —
  **Ausgang:** bei der Closure einzutragen (erster Schritt: Lesen der Strecke vom
  `Assembler` über `CaptureService` zum `Broadcaster`; Test mit gesetztem Label).
- **Signatur-Ripple.** `MatchesFilter` hat zehn Aufrufstellen (vier ohne Tests); eine
  neue Signatur bricht Tests des Bestands. — **Ausgang:** bei der Closure einzutragen.
- **Proto-Erzeugnisse.** Die Generierung läuft Docker-only; ein committeter Stand
  außer Takt färbt `make generated-sync`
  (`BEO-PGC/generierte-artefakte-ohne-sync-sensor`, verkörpert, 4×). — **Ausgang:**
  bei der Closure einzutragen.
- **V1 entschieden (`ADR-0138` Festlegung 1).** Die Gleichwertigkeit der
  Zugriffswege ([`LH-FA-SST-006`](../../../../spec/lastenheft.md)) hängt am Feld
  `target = 7`; dessen Wire-Kompatibilität (alter Server ignoriert das Feld) ist in
  der ADR *hergeleitet*, nicht gefahren. — **Ausgang:** bei der Closure einzutragen
  (Test mit dem Feld, Aussage zum alten Server mit Ursprung).

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `internal/adapters/driving/`, `internal/domain/model/`, `proto/` und `gen/` —
eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/generierte-artefakte-ohne-sync-sensor` (verkörpert, 4×),
`BEO-PGC/lese-doppelquelle` (verkörpert, 3×),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 21×),
`BEO-PGC/adapter-fehler-ausgang` (3×). Kein offener Eintrag erreicht mit diesem
Slice 3×.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

