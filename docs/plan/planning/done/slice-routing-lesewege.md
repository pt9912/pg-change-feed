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
`SPEC-031`-Zeile trägt (`ADR-0138` ist `Accepted`). Die Fixrunde des Slice ändert die Spec
an zwei Stellen: `SPEC-022` (Zeile Zustellziel) und `SPEC-031` (Zeile `ReadChanges`) tragen
den Qualifier „bei sonst gültiger Anfrage“ und den Vorrang des Lese-Kontrakts, dazu die
Geschichte-Zeile (Commits `8c3d6e2e` und `699f9f0d`, Satzbezug); `SPEC-020` und `SPEC-021`
bleiben unverändert (die Streams haben keinen Lese-Kontrakt).

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

- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-006`](../../../../spec/lastenheft.md) (A): `GET /changes?...&target=<ziel>`
      liefert genau die Changes im Bereich, deren `route_target` dem Ziel entspricht;
      ohne Parameter liefert der Endpunkt unverändert alle (Regressionstest);
      `target` kombiniert sich mit `schema`/`table` als Konjunktion; ein Parameter
      außerhalb der Menge endet `400`; das Verhalten bei einem syntaktisch
      ungültigen oder nie vergebenen Ziel folgt der Spec (`slice-routing-spec-nachzug`,
      Qualifier der Fixrunde in `SPEC-022`).
      *Zu belegen durch:* Handler-Tabellentest, Store-Test des Lesens mit Filter
      (`make test`, `make test-store`). *Beleg (Verifier):* Verifikations-Report §2
      Zeile 1 (`TestReadChangesTargetGehtAlsFilterAnDenUseCase`,
      `TestReadChangesFiltersByRouteTarget` gegen reale PostgreSQL; `make test` und
      `make test-store` Exit 0; Mutationen M7, M8, M10, S1, S2 rot).
- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (B, C): der gRPC-Stream und
      der SSE-Stream liefern mit `target` nur Changes mit diesem Ziel, ohne `target`
      jede Change (geroutete eingeschlossen); `target` kombiniert sich mit `schema`/`table`;
      der SSE-Parameter außerhalb der Menge endet `400` vor jedem Event; die
      Filterprüfung ist **eine** reine Funktion (`Change.MatchesFilter`), beide
      Handler rufen sie; `make generated-sync` grün (das Proto-Feld ist in
      `gen/cdc/stream/` erzeugt und committet); der Bestand der Aufrufer von
      `MatchesFilter` (Parent: vier Nicht-Test-Zeilen) trägt das Ziel durchgängig.
      *Zu belegen durch:* Handler-Tests mit Fake-Stream (`make test`),
      `make generated-sync`. *Beleg (Verifier):* Verifikations-Report §2 Zeile 2
      (`TestChangeMatchesFilterTarget`, je drei Stream-Tests in gRPC und SSE;
      `make generated-sync` und `make test` Exit 0; Mutationen M4b, M5b, M6, M10, P1 rot).
- [x] **Vorab-Bedingung V1** (Welle §5) ist durch
      [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
      Festlegung 1 entschieden (erfüllt); der Slice setzt sie um: `ReadChangesRequest`
      trägt `string target = 7` (leer = kein Filter, Konjunktion mit `schema`/`table`),
      der Use Case `ReadChangesUseCase` erhält den Filter einmal und bedient
      `GET /changes` und den RPC; die Aussage „der RPC sieht dasselbe wie
      `GET /changes`" ist am Test belegt (gleiche Eingabe, gleiche Changes); ein
      `target` außerhalb des Alphabets (Großbuchstabe, 64 Zeichen, U+0000) liefert auf
      allen Lesewegen eine leere Antwort, keinen Fehler — bei sonst gültiger Anfrage;
      Fehler des Lese-Kontrakts (leere Quelle, `limit < 1`, invertierter Bereich,
      Start-/End-Position einer anderen Quelle) haben Vorrang und enden unabhängig vom
      Ziel mit demselben Fehler wie im Bestand (`SPEC-022`, `SPEC-031`); Festlegung nach
      [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md):
      der gemeinsame Use Case `ReadChangesUseCase` prüft das Alphabet mit der
      Prüffunktion der Domäne und antwortet leer, ohne den Store aufzurufen; Tests:
      Use Case mit Store-Fake, der bei einem Aufruf fehlschlägt (`make test`), `GET
      /changes` und gRPC-`ReadChanges` mit derselben Eingabe, gRPC-Stream und SSE
      liefern keine Nachricht; der SQL-Zugriff `cdc.changes` ist ausgenommen;
      `make generated-sync`
      grün für beide `.proto`-Dateien, `tools/harness/grpcadminclient` trägt `target`.
      *Beleg (Verifier):* Verifikations-Report §2 Zeile 3
      (`TestReadChangesTargetOutsideAlphabetAnswersEmptyWithoutStore`,
      `TestReadChangesContractErrorsWinOverInvalidTarget`,
      `TestReadChangesWegeLiefernFuerDieselbeEingabeDieselbenChanges`,
      `TestReadChangesRequestTraegtTargetAlsFeldSieben`; Mutationen M1, M2, M11b, P2 rot).
      **Rest:** der Flag `-target` des `grpcadminclient` ist übersetzt, ein Aufruf-Test
      fehlt bis zum E2E ([`slice-routing-e2e`](slice-routing-e2e.md), Review F-6).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9). *Beleg
      (Verifier):* Verifikations-Report §1 (Exit 0, `coverage-gate` 81,90 %, `a-check`
      0 Befunde, `generated-sync` OK).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Review
      [`review-slice-routing-lesewege`](../../../reviews/review-slice-routing-lesewege.md)
      (0 HIGH, F-1 und F-2 MEDIUM, F-3 bis F-8 LOW/INFO, Architect-Frage A-1) **und**
      Re-Review der Fixrunde `8c3d6e2e`
      [`review-slice-routing-lesewege-fixrunde-1`](../../../reviews/review-slice-routing-lesewege-fixrunde-1.md)
      (0 HIGH, 0 MEDIUM, F-N1 und F-N2 LOW, F-N3 INFO; nicht merge-blockierend). F-1 und
      F-2 sind in der Fixrunde geschlossen. Die Rangfolge „der Lese-Kontrakt gewinnt“ ist
      Entscheidung des Hauptlaufs der Sitzung, **kein** Architect-Verdikt; das Re-Review
      stuft sie als von [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
      Festlegung 2 gedeckte Auslegung ein (Befund-Abschnitt der ADR: nur Parameterform und
      Bereichsfehler enden `400`; Vorbild `schema`/`table`). Die Norm steht in der Spec
      (`SPEC-022`, `SPEC-031`), in keiner ADR. Die Bestätigung durch den Architect ist
      optional (Re-Review-Frage F-N1, nicht blockierend); Trigger für eine Folge-ADR: der
      nächste Lese-Weg mit einem Filter-Parameter, der ein Alphabet prüft.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-lesewege.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13). *Beleg:*
      Verifikations-Report §1 (Exit 0 mit 16 Zeilen); die Zeilen der zweiten bewegten
      Eigenschaft (Rangfolge) sind bei der Closure ergänzt und nachgemessen (§3).
- [x] Doku-Update: `spec/pflichtenheft.md` — die `SPEC-031`-Zeile `ReadChanges` mit dem
      Feld `target` trägt `slice-routing-spec-nachzug` (`ADR-0138` Folgepflicht); die
      Fixrunde ändert `SPEC-022` (Zeile Zustellziel), `SPEC-031` (Zeile `ReadChanges`) und
      die Geschichte-Zeile (Qualifier „bei sonst gültiger Anfrage“, Vorrang des
      Lese-Kontrakts; `8c3d6e2e`, Satzbezug `699f9f0d`); das Benutzerhandbuch
      (Abschnitte „Zugriff über die
      HTTP-/JSON-API", „…gRPC-Change-Stream", „…gRPC-Verwaltungs-API", „…Server-Sent-Events")
      bleibt unberührt — Adresse: `slice-routing-betriebsdoku` §2 (Parameter `target`
      je Weg, Konjunktion mit `schema`/`table`, Qualifier „bei sonst gültiger Anfrage“, die
      Beispiele der Clients folgen `slice-routing-sdk-beispiel-target`). *Beleg
      (Verifier):* Verifikations-Report §2 Zeile 7, §6.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert (§7: drei weitere `evidence/`-Dateien in
      bestehenden Einträgen, ein Gegenbeleg im `state.md` ohne Datei).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten). *Beleg:*
      `welle-routing-results.md`, Abschnitt „Drei Paarungen“ (Closure 2026-10-02).

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

Festlegungen der Umsetzung (nicht im ursprünglichen Plan; hier vor dem Sensor-Lauf nachgezogen):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `MatchesFilter(schema, table, target string)` | Form | eine Funktion, drei Textargumente (kein Filter-Wert-Typ): beide Stream-Handler rufen sie mit dem Ziel; ein leeres Ziel wählt nicht aus, ein gesetztes vergleicht gleich mit `RouteTarget`. |
| `outbound.ChangeQuery.Target`, `inbound.ReadChangesQuery.Target` | update | das Ziel als Textfeld an beiden Typen (Port-Typ und Query-Typ); der Store reicht es als `$6` (`textArgument`), `LIMIT` wandert auf `$7`. |
| `internal/application/usecase/readchanges/service.go` | update | die Alphabet-Prüfung (`model.IsValidRouteTarget`) steht **nach** der Prüfung der Quelle und der Validierung des Lese-Kontrakts (`outbound.ChangeQuery.Validate`, die eine Stelle, die auch der Port ruft) und **vor** dem Port-Aufruf: ein ungültiges Ziel antwortet bei sonst gültiger Anfrage mit der leeren, gesetzten Liste; ein leeres `Source`, `limit < 1`, ein invertierter Bereich und eine Start-/End-Position einer anderen Quelle enden unabhängig vom Ziel mit demselben Fehler wie im Bestand. **Rangfolge:** der Lese-Kontrakt gewinnt — Entscheidung des Hauptlaufs nach dem Review `review-slice-routing-lesewege` (F-1/A-1), kein Architect-Verdikt; die Spec-Zeilen `SPEC-022`/`SPEC-031` tragen den Qualifier „bei sonst gültiger Anfrage". |
| `internal/bootstrap/readchanges_paritaet_test.go` | neu | die Aussage „der RPC sieht dasselbe wie `GET /changes`" braucht den echten Use Case und beide Adapter in einem Test; ein Adapter-Paket darf den Use Case nicht importieren (`.a-check.yml`), die Composition Root darf es. Der Store ist ein Fake mit der Gleichheits-Auswahl der Lese-Anweisung und einem Aufruf-Zähler. Der gRPC-Server läuft auf einer Loopback-Adresse (`--network none` trägt Loopback); ein Start-Fehler des Servers kommt über einen Kanal in die Fehlermeldung des Tests. Der Fake ruft `ChangeQuery.Validate` wie der reale Store; die Fehlerfälle (`limit < 1`, invertierter Bereich) laufen mit gültigem und ungültigem Ziel über beide Wege. |
| `gen/cdc/stream/v1/changestream_test.go`, `gen/cdc/administration/v1/administration_test.go` | update | Getter `GetTarget` auf dem Nullwert und gesetzt; Draht-Bytes der Felder 3 bzw. 7 (feste Bytes), Lesen ohne das Feld, Lesen mit einem dem Empfänger unbekannten Feld — gemessen an der Protobuf-Bibliothek dieses Repos, nicht an einem ausgelieferten Altserver. |
| `internal/application/usecase/capture/service_test.go` | update | Die Strecke von der Transaktion (das Ziel ist dort per `WithRouteTarget` gesetzt; der Assembler-Schritt `replication/mapper/mapper.go` liegt nicht auf der gefahrenen Strecke) über `CaptureService` zum `ChangeStreamPort`-Fake: dieselbe Change mit demselben `RouteTarget` (Teilbeleg zum Risiko „Kommt das Label an den Handler?" in §6; `Broadcaster` und Assembler bleiben gelesen, nicht gefahren); Test mit gesetztem und leerem Label. |
| `tools/harness/grpcadminclient/main.go` | update | das Ziel ist ein **Flag** `-target` vor den elf Positionsargumenten (leer = kein Filter): der Aufrufer `tools/harness/run-integration-tests.sh` bleibt unverändert; die Nutzung im E2E-Lauf trägt `slice-routing-e2e`. |
| `tools/harness/grpcclient`, `tools/harness/sseclient` | nicht geändert | die Wegwerf-Clients der Streams tragen nur `schema`/`table`; das Ziel am Stream-Beleg am laufenden System gehört `slice-routing-e2e` (gemeldet). |

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
57e4bffe 20 -n -E '`schema`/`table`|schema/table|schema und table' -- spec docs/user
57e4bffe 3 -n -E 'MatchesFilter\([a-zA-Z.]*, [a-zA-Z.]*\)' -- internal
diff 4 -n MatchesFilter -- internal ':!*_test.go'
diff 18 -n MatchesFilter -- internal
diff 2 -n -F 'MatchesFilter(schema, table, target)' -- internal ':!*_test.go'
diff 0 -n -E 'MatchesFilter\([a-zA-Z.]*, [a-zA-Z.]*\)' -- internal
diff 8 -n -E 'StreamChangesRequest' -- internal proto spec docs/user ':!*_test.go' ':!gen'
diff 3 -n -E 'ReadChangesRequest' -- internal proto spec ':!gen' ':!*_test.go'
diff 10 -n -E 'Query-Parameter' -- spec docs/user
diff 20 -n -E '`schema`/`table`|schema/table|schema und table' -- spec docs/user
30fd6cb5 0 -n -E 'sonst gültiger Anfrage' -- spec docs/user docs/plan/planning/welle-routing.md
11cff8d6 0 -n -E 'sonst gültiger Anfrage' -- spec docs/user docs/plan/planning/welle-routing.md
diff 3 -n -E 'sonst gültiger Anfrage' -- spec docs/user docs/plan/planning/welle-routing.md
30fd6cb5 2 -n -E 'außerhalb des Alphabets' -- spec docs/user
11cff8d6 9 -n -E 'außerhalb des Alphabets' -- spec docs/user
diff 9 -n -E 'außerhalb des Alphabets' -- spec docs/user
30fd6cb5 1 -n -E '\.Validate\(\)' -- internal/application/usecase/readchanges internal/adapters/driven/postgresstorage ':!*_test.go'
11cff8d6 1 -n -E '\.Validate\(\)' -- internal/application/usecase/readchanges internal/adapters/driven/postgresstorage ':!*_test.go'
diff 2 -n -E '\.Validate\(\)' -- internal/application/usecase/readchanges internal/adapters/driven/postgresstorage ':!*_test.go'
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Aufrufer der Filterfunktion | Zeilen 1 und 2: 4 Nicht-Test-Zeilen, 10 mit Tests | jeder Aufrufer trägt das Ziel durchgängig; Befund am Diff (Zeilen 8 bis 12): gefunden — die zwei Aufrufer (`grpc/server.go`, `http/sse.go`) rufen mit drei Argumenten, die Definition (`change.go`) trägt `target`, die Tests rufen mit drei Argumenten; nicht gefunden — ein Aufruf mit zwei Argumenten (Zeile 12: 0; am Parent Zeile 8: 3, davon einer im Test). Der Parent `57e4bffe` der Umsetzung liefert dieselben Zahlen wie `30fd6cb5` (4 und 10). |
| Definitionen und Aufrufer der Anfrage-Typen | Zeilen 3 und 4: 8 bzw. 3 Zeilen | jede Stelle gelesen; Befund am Diff (Zeilen 13 und 14: 8 und 3, unverändert): die Zahlen bleiben, weil die Felder in vorhandene Zeilen gelangen (`target` als Feld in `StreamChangesRequest`/`ReadChangesRequest` in den `.proto`-Dateien und in den Aufrufstellen `req.GetTarget()`/`Target:`); die Aufrufer `grpc/server.go` (Stream) und `grpc/administration.go` (RPC) tragen das Ziel; nicht gefunden — eine Aufzählung der Request-Felder, die `target` auslässt. Die erzeugten Dateien unter `gen/` stehen außerhalb dieser Zählung (`generated-sync`). |
| Beschreibungen der Parameter in Spec und Handbuch | Zeilen 5 und 6: 10 bzw. 17 Zeilen („Query-Parameter", „`schema`/`table`") | Spec-Zeilen zieht `slice-routing-spec-nachzug` (bereits geändert, wenn dieser Slice startet); Handbuch gemeldet an `slice-routing-betriebsdoku`, Client-Beispiele an `slice-routing-sdk-beispiel-target`; Befund am Diff (Zeilen 15 und 16: 10 und 20): „Query-Parameter" unverändert 10; „`schema`/`table`" 17 → 20 am Parent `57e4bffe` (Zeile 7) und am Diff gleich 20 — der Zuwachs von drei stammt aus dem Spec-Nachzug (`spec/pflichtenheft.md`, Zeilen mit `schema`/`table`/`target`), nicht aus diesem Diff. Gefunden und gemeldet, nicht mitgeändert: `docs/user/benutzerhandbuch.md` an den Zeilen 1302, 1365, 1378, 1403, 1456, 1558 (Adresse `slice-routing-betriebsdoku` §2), `docs/user/e2e-abdeckung.md` an den Zeilen 60 und 67 (Erzeugnis des Integrationslaufs, beschreibt vorhandene Belege des `schema`/`table`-Filters; das Ziel am Stream belegt `slice-routing-e2e`). Nicht gefunden — eine Spec-Zeile, die `schema`/`table` als einzigen Filter ausgibt. |
| Zweite bewegte Eigenschaft der Fixrunde: „ein `target` außerhalb des Alphabets liefert leer, keinen Fehler“ wird bedingt („bei sonst gültiger Anfrage“, Lese-Kontrakt hat Vorrang) | Zeilen 17 bis 22 (`außerhalb des Alphabets` in `spec`/`docs/user`: 2 am Planungs-Parent `30fd6cb5`, 9 am Stand vor der Fixrunde `11cff8d6` und am Diff; `sonst gültiger Anfrage` in `spec`, `docs/user` und der Welle: 0, 0 und 3), Zeilen 23 bis 25 (Nicht-Test-Aufrufer von `.Validate()` in Use Case und Store: 1, 1 und 2) | Gefunden und nachgezogen: `SPEC-022` und `SPEC-031` (Fixrunde), `welle-routing.md` (Zeile zum Alphabet), die DoD-Zeile 3 dieses Plans, der Kopf dieses Plans (Spec-Änderung), `slice-routing-betriebsdoku` §2 (Adressliste und Qualifier, gemeldet und von der Closure gezogen). Gemeldet, nicht mitgeändert: das Handbuch (`docs/user`) trägt keine Aussage zum Alphabet des Ziels (Muster `Alphabet` in `docs/user`: kein Treffer zum `target`), die Zeile 1149 („gefiltert über `schema` und `table`“) fehlte in der Adressliste und ist aufgenommen. Nicht gefunden: eine Spec-Zeile, die „leer, kein Fehler“ ohne den Qualifier für `GET /changes` oder `ReadChanges` führt; die zwei Stream-Zeilen (`SPEC-020`, `SPEC-021`) tragen ihn zu Recht nicht (kein Lese-Kontrakt). Die ADR-Zeile (`ADR-0139` Festlegung 2) ist `Accepted` und bleibt. |
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
  Felder — an der Protobuf-Bibliothek dieses Repos gemessen, `TestStreamChangesRequestIgnoriertUnbekannteFelder` und
  `TestReadChangesRequestTraegtTargetAlsFeldSieben`; an einem ausgelieferten Altserver *hergeleitet*, nicht gemessen) und am HTTP-/SSE-Weg `400` (die Parameter-
  Menge ist streng, Bestand: `SPEC-021`/`SPEC-022`). — **Ausgang:** weiter offen für den
  ausgelieferten Altserver: gemessen ist das Verhalten der Protobuf-Bibliothek dieses
  Repos (die zwei Tests oben, im Verifier-Lauf `make test` Exit 0), am ausgelieferten
  Altserver ist die Aussage *hergeleitet*; eine Messung ist nirgends eingeplant. Adresse:
  [`slice-routing-betriebsdoku`](slice-routing-betriebsdoku.md) §2 (die Aussage im
  Handbuch trägt den Ursprung *hergeleitet*).
- **Das Ziel ist Auswahl, kein Zugriffsschutz.** Jeder Leser mit `reader`-Token kann
  jedes Ziel wählen oder gar nicht filtern
  ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Entscheidung 2, ungefilterte Leser sehen alles); eine Lesart „Ziel A sieht nur A"
  als Berechtigung ist falsch. — **Ausgang:** entfallen für diesen Slice (er ändert kein
  Handbuch; die Aussage ist Norm der ADR, im Code gilt: ohne `target` liefert jeder Weg
  alle Changes, `TestChangeMatchesFilterTarget`, `TestReadChangesFiltersByRouteTarget`).
  Die Handbuch-Aussage trägt [`slice-routing-betriebsdoku`](slice-routing-betriebsdoku.md)
  §2.
- **Kommt das Label an den Handler?** Der Handler filtert auf `model.Change`; ob der
  Broadcaster dieselbe Change mit gesetztem `RouteTarget` liefert, die persistiert
  wurde, ist *erwartet* (ein Objekt, ein Erzeugungspfad), nicht gelesen. —
  **Ausgang:** teils entfallen, teils weiter offen. Gefahren ist die Strecke von der
  Transaktion über `CaptureService` zum `ChangeStreamPort`-Fake mit gesetztem und leerem
  Label (`TestCapturePublishesTheRouteTargetOfEachChange`, Verifier-Lauf `make test`
  Exit 0); `Broadcaster` und der Assembler-Schritt sind gelesen, nicht gefahren. Die
  Gesamtkette bis zum Handler am laufenden System steht aus. Adresse:
  [`slice-routing-e2e`](slice-routing-e2e.md) §2 (Ziel an den Stream-Wegen, die
  Wegwerf-Clients `grpcclient`/`sseclient` erhalten die Auswahl des Ziels als Flag).
- **Signatur-Ripple.** `MatchesFilter` hat zehn Aufrufstellen (vier ohne Tests); eine
  neue Signatur bricht Tests des Bestands. — **Ausgang:** entfallen. Die Signatur ist
  `MatchesFilter(schema, table, target string)`; Suchlauf in §3: kein Aufruf mit zwei
  Argumenten am Diff (Zeile 12: 0), `make test` Exit 0 (Verifier-Lauf).
- **Rangfolge ungültiges Ziel gegen Lese-Kontrakt.** Weder `ADR-0139` noch die Spec
  entscheiden, ob ein Ziel außerhalb des Alphabets einen Fehler des Lese-Kontrakts
  verdeckt. — **Ausgang:** eingetreten (Review F-1/A-1), entschieden vom Hauptlauf, kein
  Architect-Verdikt: der Lese-Kontrakt gewinnt; Umsetzung in `readchanges/service.go`
  mit `ChangeQuery.Validate` vor der Alphabet-Prüfung, getestet im Use-Case-Test
  (`TestReadChangesContractErrorsWinOverInvalidTarget`) und im Paritätstest, Spec-Qualifier
  in `SPEC-022`/`SPEC-031`; Mutationen M1 und M2 des Verifiers rot (Verifikations-Report §4).
  Die Entscheidung berührt `ADR-0139` Festlegung 2 nicht inhaltlich („leer" gilt für die
  sonst gültige Anfrage); das Re-Review stuft sie als von der Festlegung gedeckte Auslegung
  ein (Befund-Abschnitt der ADR: nur Parameterform und Bereichsfehler enden `400`). Die
  Bestätigung durch den Architect ist optional (F-N1, nicht blockierend); Trigger für eine
  Folge-ADR: der nächste Lese-Weg mit einem Filter-Parameter, der ein Alphabet prüft.
- **Proto-Erzeugnisse.** Die Generierung läuft Docker-only; ein committeter Stand
  außer Takt färbt `make generated-sync`
  (`BEO-PGC/generierte-artefakte-ohne-sync-sensor`, verkörpert, 4×). — **Ausgang:**
  entfallen. `make generated-sync` Exit 0 im Verifier-Lauf; die Mutationen P1 und P2
  (Feldnummer geändert) färben es rot.
- **V1 entschieden (`ADR-0138` Festlegung 1).** Die Gleichwertigkeit der
  Zugriffswege ([`LH-FA-SST-006`](../../../../spec/lastenheft.md)) hängt am Feld
  `target = 7`; dessen Wire-Kompatibilität (alter Server ignoriert das Feld) ist in
  der ADR *hergeleitet*, nicht gefahren. — **Ausgang:** entfallen für die Gleichwertigkeit
  (`TestReadChangesWegeLiefernFuerDieselbeEingabeDieselbenChanges`, das Feld in
  `TestReadChangesRequestTraegtTargetAlsFeldSieben`); für den Altserver siehe das erste
  Risiko (weiter offen, *hergeleitet*).
- **Beleg am laufenden System (aus §1).** Dieser Slice belegt auf Unit-, Handler- und
  Store-Ebene. — **Ausgang:** weiter offen. Adresse
  [`slice-routing-e2e`](slice-routing-e2e.md) §2: das Ziel an den Stream-Wegen und
  an `GET /changes` am laufenden Feed-Container, ein Aufruf-Test für den Flag `-target`
  des `grpcadminclient` (Review F-6).
- **Port-Wahl im Paritätstest (Review F-3).** `freeLoopbackAddr` wählt den Loopback-Port
  per Listen/Close; zwischen Wahl und Start liegt ein Zeitfenster. — **Ausgang:** weiter
  offen, benannter Rest: ein belegter Port zeigt sich seit der Fixrunde als benannter
  Start-Fehler (`failWithStartError`), nicht als Zeitüberschreitung; das Fenster selbst
  bleibt. Keine Adresse außer diesem Test; Trigger: ein Flake des Paritätstests.
- **Parität des Validate-Aufrufs im Store (Re-Review F-N3).** Der Fake des Paritätstests
  ruft `ChangeQuery.Validate`; dass der reale Store es ruft, ist am Code gelesen
  (`store.go`). — **Ausgang:** entfallen. Der Bestandstest `TestReadCarriesPositionsAndRanges`
  (`store_test.go`, invertierter Bereich endet `ErrRangeInverted`) läuft in `make test-store`
  gegen reale PostgreSQL; Verifier-Lauf Exit 0. Er trägt den Aufruf des Stores für den
  invertierten Bereich, nicht jede Fehlerart (*hergeleitet* für `limit < 1`: dieselbe
  Funktion).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die drei Liefer-Punkte tragen: das Ziel an `GET /changes`, am
  gRPC-`ReadChanges` (ein Use Case für beide Wege), am gRPC-Stream und am SSE-Stream, mit
  **einer** Filterfunktion `MatchesFilter(schema, table, target)`; die Proto-Felder (3 und 7)
  sind erzeugt und synchron. Gemessen im Lauf des Verifiers: `make generated-sync`, `make test`,
  `make test-store`, `make gates` (Exit 0, `coverage-gate` 81,90 %), `make docs-check`,
  Suchlauf (16 Zeilen) alle Exit 0 (Verifikations-Report §1). Der Paritätstest in der
  Composition Root (echter Use Case, beide Adapter) trägt die Aussage „der RPC sieht dasselbe
  wie `GET /changes`“; die Auswahl selbst trägt der Store-Test gegen reale PostgreSQL. Das Muster
  der Filter-Parameter `schema`/`table` trug den Schnitt ohne Rückführung nach §4.
- **Was ging anders als geplant:** Eine Fixrunde (`8c3d6e2e`) und ein Spec-Nachzug (`699f9f0d`).
  Das Review fand die Rangfolge „ungültiges Ziel gegen Fehler des Lese-Kontrakts“ (F-1, MEDIUM)
  und den dazu widersprüchlichen Godoc (F-2, MEDIUM); weder `ADR-0139` noch die Spec
  entschieden sie. Der Hauptlauf der Sitzung entschied: der Lese-Kontrakt gewinnt. Die Fixrunde
  setzte `ChangeQuery.Validate` vor die Alphabet-Prüfung und qualifizierte `SPEC-022` und
  `SPEC-031` („bei sonst gültiger Anfrage“). Damit änderte der Slice die Spec, die der Plan als
  unberührt führte (V-1) und ließ zwei Träger mit der unqualifizierten Aussage stehen (V-2);
  die Closure zog den Plan-Kopf, die DoD-Zeilen 3 und 7, `welle-routing.md` und die Adressliste
  von `slice-routing-betriebsdoku` nach (dort zusätzlich Handbuch-Zeile 1149 samt erweitertem
  Suchmuster, V-5). Das Re-Review der Fixrunde (Diff `11cff8d6..8c3d6e2e`) fand keinen Fehler im
  Verhalten. Die Mutationszahlen tragen ihren Ursprung (Instanz A von
  [`AGENTS.md`](../../../../AGENTS.md) §3.12): Implementer rund 17, **übernommen** (Bericht,
  nicht nachgefahren); Reviewer 13 selbst gefahren, **gemessen** (drei weitere Läufe scheiterten
  am Übersetzen und zählen nicht); Verifier 14 selbst gefahren, **gemessen** (alle rot; drei
  Läufe scheiterten am Übersetzen und zählen nicht); Re-Review 4 selbst gefahren, **gemessen**
  (alle rot). Wer die Zahl liest, liest die Läufe unabhängiger Leser, nicht die des Autors.
- **Steering-Loop-Eintrag:** geschärfte Regel, kein neuer Sensor. Eine **Fehlerrangfolge** ist
  ein Vertrag des neuen öffentlichen Parameters: sie gehört vor dem Code in die ADR- oder
  Spec-Zeile des Parameters, nicht in die Auslegung des Implementers oder des Reviews. Prüffrage
  für jeden neuen Filter-Parameter beim Schreiben seiner Zeile: „Was gewinnt bei gleichzeitigem
  Kontraktfehler?“ — also bei einem Wert außerhalb seiner Menge **und** einem Fehler des
  Lese-Kontrakts (`limit`, Bereich, Quelle). Hier stand sie in keiner der beiden Zeilen, der
  Implementer meldete sie als „nicht entschieden“, das Review fand sie als MEDIUM und reichte sie
  als Architect-Frage weiter; die Antwort kam vom Hauptlauf und steht jetzt als Qualifier in der
  Spec statt in einer ADR. Die Handlung davor: beim Schreiben der Festlegung die Zeile gegen den
  Lese-Kontrakt des Ports lesen (`ChangeQuery.Validate`) und beide Fehlerquellen nebeneinander
  stellen. Träger: die Lese-Handlung von Architect und Planner; kein Sensor, weil ob zwei
  Prüfungen eine Rangfolge brauchen, eine Lese-Frage ist; das Register führt es als zweites
  Auftreten von `BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`. Trigger für eine
  Folge-ADR, die die Rangfolge für alle Lese-Wege festschreibt: der nächste Lese-Weg mit einem
  Filter-Parameter, der ein Alphabet prüft. Zweiter Befund der Fixrunde: ein Re-Review nach einer
  Fixrunde, die Anweisungen **und** eine Norm ändert und deren Entscheidung der Leser selbst
  traf, ist keine Formsache (Verifier V-3); er wurde vor der Closure gefahren.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`** (offen): zweites
    Auftreten (F-1, MEDIUM, Ausprägung Fehlerrangfolge am Lesepfad); eine neue `evidence/`-Datei,
    Zähler **1×** → **2×** (`ls evidence | wc -l`), unter der 3×-Schwelle, kein Ausgang. Kein neues
    Verzeichnis: die Klasse „Festlegung, die der Folge-Slice beim ersten Satz Code braucht, fehlt“
    trägt den Fund (Neuformulieren spaltet die Klasse).
  - **`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`** (verkörpert): F-2 (MEDIUM, Godoc neben
    der hinzugefügten Zusage) und F-N2 (LOW, Plan nach der Fixrunde neben der geänderten Norm);
    eine neue `evidence/`-Datei (Schwere ≥ MEDIUM), Zähler **17×** → **18×**.
  - **`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`** (verkörpert): F-4 (LOW, Godoc
    des Capture-Tests breiter als die gefahrene Strecke); eine neue `evidence/`-Datei, Zähler
    **8×** → **9×** (unter dem Deckel von 10×).
  - **`BEO-PGC/fixrunde-ohne-reviewer-lesung`** (offen, 3×, Ausgang beim Lese-Schritt der Closure
    von [welle-routing](../welle-routing.md)): **dritter Gegenbeleg, keine Datei.** Die Fixrunde
    änderte Anweisungen und eine Norm, der Verifier empfahl ein Re-Review (V-3), es wurde gefahren
    (vier Mutationen, F-N3 gefunden) und vor der Closure geschlossen. Im `state.md` vermerkt.
  - Ohne eigene Datei (Deckel gilt für diese Einträge nicht, die Funde sind ≤ LOW und tragen ihre
    Klasse in den Einträgen oben oder entfallen): F-3 (LOW, Port-Wahl per Listen/Close, Rest in
    §6), F-5 (LOW, Aufschub-Adresse ohne Zählwort, in `slice-routing-betriebsdoku` nachgezogen),
    F-6 und F-7 (INFO, als Grenzen in §6 geführt), F-8 (LOW, Umbruchrest, behoben), F-N1 (LOW,
    Entscheidungs-Instanz, in §6), F-N3 (INFO, in §6 entfallen), V-4 (Satzbezug, in `699f9f0d`
    behoben), V-6 und V-7 (INFO).
  - Kein Eintrag steht bei 3× oder mehr ohne Ausgang an, den dieser Slice neu erreichte.
- **Folge-Slices:** keine neuen. Übergaben: `slice-routing-e2e` (Ziel an den Stream-Wegen und an
  `GET /changes` am laufenden System, Aufruf-Test für `grpcadminclient -target`, die Gesamtkette
  Assembler → `Broadcaster` → Handler), `slice-routing-betriebsdoku` (Handbuch-Zeilen des
  Übergabe-Blocks samt Zeile 1149 und dem Qualifier), `slice-routing-sdk-beispiel-target` (die
  Client-Beispiele). Die Fremddatei `slice-routing-betriebsdoku` wurde bei der Closure an einer
  Stelle berührt: Adressliste und ein Suchmuster (V-5), das Handbuch selbst nicht.
- **Risiken aus §6:** Altserver **weiter offen** (am ausgelieferten Altserver *hergeleitet*;
  Adresse `slice-routing-betriebsdoku` §2); Auswahl-kein-Zugriffsschutz **entfallen** für diesen
  Slice (Handbuch-Aussage bei `slice-routing-betriebsdoku`); Label am Handler **teils entfallen**
  (Transaktion → Stream-Port-Fake gefahren), die Gesamtkette **weiter offen** (Adresse
  `slice-routing-e2e` §2); Signatur-Ripple **entfallen**; Rangfolge **eingetreten** (Hauptlauf
  entschied, Spec trägt die Norm, Architect-Bestätigung optional); Proto-Erzeugnisse
  **entfallen**; V1/Wire-Kompatibilität **entfallen** für die Gleichwertigkeit, Altserver siehe
  oben; Beleg am laufenden System **weiter offen** (Adresse `slice-routing-e2e` §2); Port-Wahl im
  Paritätstest **weiter offen** (benannter Rest); Validate-Aufruf im Store **entfallen**.
- **Drei Paarungen:** der Slice gehört zu [welle-routing](../welle-routing.md) (offen) — die
  Prüfung läuft bei deren Closure; die DoD-Zeile bleibt deshalb `[ ]`. (a) Anker: der Lerneintrag
  verkörpert nichts neu in `AGENTS.md` (Prüffrage als Lese-Handlung im Register, zweites Auftreten);
  (b) Folge-Slice: keiner neu, die genannten Pläne liegen unter `open/`; (c) Register: die
  genannten Kennungen existieren als Verzeichnis, jede trägt ein nicht leeres `evidence/`.
- **Validator (Modul 8):** entfällt — der Filter `target` an den Lesewegen ist ohne Backfill-Beleg am
  System und ohne Handbuch für Betreiber noch nicht als Ganzes nutzbar; der Nutzer-Bedarf
  ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) wird erst durch den Wellen-Beleg validierbar.

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

