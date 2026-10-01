# Verifikations-Report: slice-routing-lesewege — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Teilfrage 5,
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
Festlegung 1,
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
Festlegung 2,
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md),
[`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md),
[`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)) und Plan-vs-Code-Diff.
Review-Artefakt: [`review-slice-routing-lesewege.md`](review-slice-routing-lesewege.md).
Formvorbild: [`verifikation-slice-routing-backfill-pfad.md`](verifikation-slice-routing-backfill-pfad.md).

**Gegenstand:** Slice-Plan [`slice-routing-lesewege`](../plan/planning/done/slice-routing-lesewege.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), weiter
[`LH-FA-SST-006`](../../spec/lastenheft.md), [`LH-FA-SST-008`](../../spec/lastenheft.md); Welle
[`welle-routing`](../plan/planning/welle-routing.md)). Diff `26c16275~1..HEAD` (`8c3d6e2e`, Fixrunde):
sieben Commits, 30 Dateien (`git diff --stat 26c16275~1 HEAD`: +1652/−98). Fixrunde = der einzige
Code-Commit nach dem Review `11cff8d6`: `8c3d6e2e` (acht Dateien, +174/−31).

Dieser Lauf ändert weder Code noch Plan noch Spec (keine DoD-Häkchen); er schreibt nur diesen Report.
Alle Mutationen liefen an einem `git clone` im Scratchpad (Mutation per `sed … > Temp-Datei`, danach
`cp` innerhalb des Scratchpads auf die Klon-Datei; nie `sed -i`, nie eine Umleitung auf eine Repo-Datei),
Unit-Läufe im gepinnten Toolchain-Image (Aufrufform von `make test`: `-race`, `--network none`), Store-
und Generator-Läufe über `make test-store` bzw. `make generated-sync` im Klon; Rücknahme je Mutation
`git checkout`. Der Arbeitsbaum des Repos blieb sauber. Keine verweigerte Aktion
([`AGENTS.md`](../../AGENTS.md) §3.15).

## 1. Eigene Sensor-Belege (ungefiltert in Log-Dateien, Exit-Code je Lauf einzeln gesichert, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | Ausgabe (gedruckte Zeile) |
|---|---|---|
| `make generated-sync` | 0 | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators`, geprüft beide Quellen (`changestream.proto`, `administration.proto`) und vier `.pb.go`-Dateien |
| `make test` | 0 | `ok …/internal/application/usecase/readchanges 1.013s`, `ok …/internal/bootstrap 1.672s`, `ok …/internal/adapters/driving/http 3.353s`, `ok …/internal/adapters/driving/grpc 2.657s`, `ok …/internal/domain/model 1.157s`, `ok …/gen/cdc/stream/v1`, `ok …/gen/cdc/administration/v1`; `-race` |
| `make test-store` | 0 | `ok …/internal/adapters/driven/postgresstorage 14.102s` (darin `TestReadChangesFiltersByRouteTarget`: Mutationen S1 und S2 färben ihn rot, er läuft also real gegen PostgreSQL), `ok …/internal/bootstrap 5.696s` |
| `make gates` | 0 | `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung unter sdks`, `coverage-gate: OK — Coverage 81.90% erfüllt Schwelle 80%`, `a-check … gesamt: 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)` |
| `make docs-check` (vor Anlage dieses Reports) | 0 | `d-check: 1497 Datei(en) geprüft, 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-lesewege.md` | 0 | `suchlauf-nachmessen: 16 Zeilen stimmen` (letzte Zeile `OK soll=20 ist=20 diff 20 … schema/table …`) |
| `make kommentar-kennungen DIFF=26c16275~1` | 0 | keine Ausgabe, kein Kandidat (Form-Probe, kein Beleg der Wahrheit) |
| `make fmt-check` | 0 | `fmt-check: 317 Go-Dateien geprüft, alle formatiert` |
| `make commit-traceability RANGE=26c16275~1..HEAD` | 0 | `OK — 7 Commit(s) in "26c16275~1..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=26c16275~1..HEAD` | 0 | `d-check: 1497 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=26c16275~1..HEAD` (ohne `RANGE` bricht das Ziel mit Exit 2 ab: `flag needs an argument: --range`) | 0 | `d-check: 1497 Datei(en) geprüft, 0 Befund(e)` |

Zwei frühere Lauf-Versuche (eine Kette, deren Subshell mit dem Aufruf endete, und eine versehentlich
zweite parallele Kette) habe ich abgebrochen und verworfen; die Tabelle gibt ausschließlich Läufe einer
einzigen, sauber gefahrenen Kette wieder (Zeitstempel der Logs). `make test-integration` und `make
test-replication` habe ich nicht gefahren (laut Plan §1 gehört der Beleg am laufenden System
[`slice-routing-e2e`](../plan/planning/in-progress/slice-routing-e2e.md)).

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | (A) `GET /changes?target=` liefert genau die Changes mit dem Ziel; ohne Parameter unverändert alle; Konjunktion mit `schema`/`table`; unbekannter Parameter `400`; Verhalten bei ungültigem/unvergebenem Ziel folgt der Spec; Store-Filter | **getragen** | Handler: `TestReadChangesTargetGehtAlsFilterAnDenUseCase` (Fälle `target=eu`, mit `schema`/`table` als Konjunktion, `target=` leer, `EU`, `eu%00`), `TestReadChangesTargetNameBleibtStrengeParameterMenge`, `TestReadChangesLeereAntwortDesUseCaseBleibtGesetzteListe`; Regression: die Bestandstests laufen unverändert grün, ihre geänderten Zeilen sind nur Ergänzungen um `Target` (Diff der Testdateien gelesen, kein abgeschwächter Vergleich). Store real: `TestReadChangesFiltersByRouteTarget` (Gleichheit, `e` und `e%` treffen nicht, Konjunktion mit Tabelle und mit Limit, „kein Ziel wählt nicht aus“). Mutationen M7, M8, M10, S1, S2 rot (§4). Hinweis zum Spec-Teil: V-1 |
| 2 | (B, C) gRPC- und SSE-Stream mit `target`; `target` kombiniert; SSE-Parameter außerhalb der Menge `400` vor jedem Event; **eine** Funktion `MatchesFilter`; `make generated-sync` grün; Aufrufer tragen das Ziel | **getragen** | `MatchesFilter(schema, table, target)` ist die einzige Filterfunktion (`change.go`); beide Handler rufen sie mit drei Argumenten (Suchlauf-Zeilen: 4 Nicht-Test-Zeilen am Parent und am Diff, `MatchesFilter(…, …)` mit zwei Argumenten: 0 am Diff, Nachmessen Exit 0). Tests: `TestChangeMatchesFilterTarget` (elf Fälle), `TestStreamChangesTargetFilterLaesstNurDasZiel`/`…OhneTrefferLiefertKeineNachricht`/`…OhneTargetLiefertGeroutetEbenfalls` (gRPC), `TestStreamTargetFilterLaesstNurDasZiel`/`…OhneTrefferLiefertKeinEvent`/`…OhneTargetLiefertGeroutetEbenfalls`/`…ParameterNameBleibtStrengeMenge` (SSE). `make generated-sync` Exit 0 (§1). Mutationen M4b, M5b, M6, M10, P1 rot |
| 3 | V1 ([`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md) Festlegung 1) und [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Festlegung 2: `ReadChangesRequest.target = 7`; ein Use Case für `GET /changes` und RPC, Gleichheit am Test; Alphabet im Use Case ohne Store-Aufruf; Stream/SSE keine Nachricht; `grpcadminclient` trägt `target` | **getragen, mit Qualifier und Plan-Widerspruch (V-1, V-2)** | Feld 7 in `administration.proto`, Draht-Bytes `TestReadChangesRequestTraegtTargetAlsFeldSieben`; P2 (Feldnummer 8) färbt `generated-sync` rot. Use Case: `TestReadChangesTargetOutsideAlphabetAnswersEmptyWithoutStore` (sieben Fälle, Store-Fake endet bei jedem Aufruf mit Fehler und zählt, M11b rot). Gleichwertigkeit: `TestReadChangesWegeLiefernFuerDieselbeEingabeDieselbenChanges` (echter Use Case, beide Adapter, Composition Root). `grpcadminclient`: Flag `-target`, übersetzt (`go test ./...` kompiliert das Programm), ein Aufruf-Test fehlt bis zum E2E (Review F-6, benannt). **Einschränkung:** die Zusage „leer, kein Fehler“ gilt seit der Fixrunde nur „bei sonst gültiger Anfrage“ (siehe §6); der Text dieser DoD-Zeile steht unqualifiziert (V-2) |
| 4 | `make gates` grün, Exit-Code ungefiltert | **getragen** | §1, eigener Lauf, Exit 0 |
| 5 | Review durchgeführt, Report liegt vor, kein offenes HIGH/MEDIUM | **offen — Re-Review nötig (V-3)** | Report liegt vor (0 HIGH, 2 MEDIUM, 4 LOW, 2 INFO). F-1/F-2 sind am Code, Test und Text geschlossen (§5), aber die Fixrunde ändert Anweisungen **und** eine Norm und ist von keinem anderen Kontext als dem entscheidenden gelesen (§7) |
| 6 | §3.13-Suchlauf: Feld mit Gefundenem und Nichtgefundenem, beide Stände; Nachmessen Exit 0 | **getragen für die erste Eigenschaft, Lücke für die der Fixrunde (V-2)** | Exit 0, 16 Zeilen; die Befund-Spalte nennt Treffer, Auslassungen und Meldungen. Die Fixrunde bewegt eine zweite Eigenschaft (die Zusage „leer, kein Fehler“ wird bedingt), für die der Plan keine Suchlauf-Zeile führt |
| 7 | Doku-Update: Spec entfällt (Träger `slice-routing-spec-nachzug`), Handbuch unberührt, Adresse [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md) | **teilweise — Plan widerspricht dem Diff (V-1)** | Handbuch (`docs/user/`) im Diff nicht berührt (`git diff --stat`). Die Spec **wurde** berührt (`spec/pflichtenheft.md`, +3/−3 in `8c3d6e2e`), die Zeile sagt weiter „entfällt“. Adresse gelesen und im Handbuch nachgelesen (§5, F-5) |
| 8–13 | Closure-Notiz, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §7 des Plans trägt Platzhalter, §6 hat fünf Ausgänge „bei der Closure einzutragen“ (der Ausgang zur Rangfolge ist eingetragen) |

Die DoD-Häkchen im Plan stehen unverändert alle auf `[ ]`; ich setze keine.

## 3. Plan-vs-Code-Diff

Plan-Tabelle §3 gegen `git diff --stat 26c16275~1 HEAD`:

- **Im Plan, im Diff:** `change.go`/`change_test.go` (`MatchesFilter`), die beiden Query-Typen und die Store-Anweisung samt `store.go`, `readchanges.go`, `sse.go`, `server.go`, `administration.go` (je ein Feld bzw. Argument), beide `.proto`-Dateien samt vier erzeugten Dateien unter `gen/`, `grpcadminclient/main.go`, die Tests je Weg.
- **Im Plan als Umsetzungs-Festlegung geführt, im Diff:** `readchanges_paritaet_test.go` (neu), `capture/service_test.go` (Godoc seit der Fixrunde ehrlich auf die gefahrene Strecke begrenzt, gelesen), Draht-Tests in `gen/`.
- **Im Diff, nicht in der Plan-Tabelle:** `spec/pflichtenheft.md` (Fixrunde, V-1), der Review-Report, die Fremddatei [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md) (Übergabe-Block, in §3 als Meldung geführt).
- **Im Plan, nicht im Diff:** nichts gefunden. Wecksignal, `ChangeStreamPort`, `Broadcaster`, NATS-Adapter, Nachrichtenschema (`toProtoChange`, `toStreamChange`, `ChangeRecord`) sind im Diff unberührt, wie §1 des Plans es zusagt.
- **Nullpfad ohne `target`:** das neue Prädikat `$6::text IS NULL OR …` trifft bei einem leeren Ziel nicht ein (`textArgument("")`), der Use Case nimmt den ersten Zweig nicht; die Bestandstests sind unverändert grün.
- **Docker-only (§3.1):** im Diff kein `sed -i`, kein Host-Interpreter, keine Umleitung auf Repo-Dateien erkennbar; `make fmt-check` Exit 0.

## 4. Mutationen, selbst nachgefahren (Klon im Scratchpad, Rücknahme `git checkout`)

Gefahren: 14 gültige Läufe, alle rot (10 Unit im Toolchain-Image mit `-race`, 2 über `make test-store`, 2 über `make generated-sync`). Drei Läufe (M4, M5, M11) scheiterten am Übersetzen (unbenutzte Variable bzw. Import) und zählen nicht; sie wurden in kompilierbarer Form wiederholt (M4b, M5b, M11b).

| # | Mutation | Ziel | Ergebnis |
|---|---|---|---|
| M1 (PFLICHT 1) | `service.go`: `portQuery.Validate()` im Zweig des ungültigen Ziels wirkungslos (`error(nil)`) — Alphabet-Prüfung faktisch vor den Lese-Kontrakt gesetzt | `make test`-Pakete | **rot** — `TestReadChangesContractErrorsWinOverInvalidTarget`, `TestReadChangesWegeLiefernFuerDieselbeEingabeDieselbenChanges` |
| M2 (PFLICHT 2) | `changestore.go`: Limit-Prüfung in `ChangeQuery.Validate` entfernt | dieselben Pakete | **rot** — `TestReadChangesCarriesNonPositiveLimit`, `TestReadChangesContractErrorsWinOverInvalidTarget`, Paritätstest. Eingabe-Bindung: ungültiges Ziel **und** `limit < 1` endet am unmutierten Baum mit `ErrNonPositiveLimit` (Fälle `Limit unter 1 Ziel EU` und `Ziel ` leer, Port-Aufrufe 0 bzw. 1), HTTP `400`, gRPC `InvalidArgument` (Paritätstest, vier Fehlerfälle) |
| S1 (PFLICHT 3) | `queries.go`: `OR TRUE` an der Zeile `route_target` | `make test-store` (reale PostgreSQL) | **rot** — `TestReadChangesFiltersByRouteTarget` (`FAIL … postgresstorage 13.545s`) |
| M4b (PFLICHT 4) | `sse.go`: `MatchesFilter(schema, table, target[:0])` — SSE reicht das Ziel nicht durch | Unit | **rot** — `TestStreamTargetFilterLaesstNurDasZiel`, `…OhneTrefferLiefertKeinEvent` |
| P1 (PFLICHT 5) | `changestream.proto`: Feldnummer 3 → 4 | `make generated-sync` im Klon | **rot** — Exit 2, `weicht von der Generatorausgabe ab (… Zeile 159)` |
| P2 | `administration.proto`: Feldnummer 7 → 8 | `make generated-sync` im Klon | **rot** — Exit 2, `… Zeile 1161` |
| M5b | `server.go`: gRPC-Stream reicht das Ziel nicht durch | Unit | **rot** — `TestStreamChangesTargetFilterLaesstNurDasZiel`, `…OhneTrefferLiefertKeineNachricht` |
| M6 | `change.go`: Vergleich `!=` → `==` | Unit | **rot** — `TestChangeMatchesFilterTarget`, beide Stream-Tests |
| M7 | `service.go`: Use Case reicht `Target` nicht an den Port | Unit | **rot** — `TestReadChangesTranslatesTarget`, Paritätstest |
| M8 | `readchanges.go`: `GET /changes` liest `target` nicht | Unit | **rot** — `TestReadChangesTargetGehtAlsFilterAnDenUseCase`, Paritätstest |
| M9 | `administration.go`: RPC reicht `Target` nicht | Unit | **rot** — `TestReadChangesRuftUseCaseMitUebersetzterQueryAuf`, `TestReadChangesTargetWirdUnveraendertUebergeben`, Paritätstest |
| M10 | `sse.go`: `target` aus der geschlossenen Parameter-Menge | Unit | **rot** — beide Stream-Tests |
| M11b | `service.go`: Alphabet-Prüfung wirkungslos (`&& false`) | Unit | **rot** — `TestReadChangesTargetOutsideAlphabetAnswersEmptyWithoutStore`, Kontrakt-Test, Paritätstest |
| S2 | `store.go`: `textArgument("")` statt `query.Target` | `make test-store` | **rot** — `TestReadChangesFiltersByRouteTarget` |

Die Mutationen der Review (13 Läufe) und des Implementers sind **übernommen**, nicht nachgemessen. Die
Verallgemeinerung auf andere Stellen (etwa die Konjunktion mit `schema` am Stream-Handler) ist
*hergeleitet*: der Review hat dort bereits festgehalten, dass sie allein der Domänentest trägt.

## 5. Review-Findings F-1 bis F-8 — an Code, Tests, Plan und Handbuch geprüft, nicht am Fixrunden-Bericht

| Finding | Verdikt | Beleg |
|---|---|---|
| F-1 (MEDIUM) Rangfolge ungültiges Ziel gegen Lese-Kontrakt | **in der Sache geschlossen, Entscheidungs-Instanz offen (V-1, V-3)** | `service.go` ruft im Zweig des ungültigen Ziels `portQuery.Validate()` (die eine Stelle, die auch `store.go` ruft); `Validate` ist einmal definiert, es gibt keine zweite Prüflogik. Reihenfolge: leere Quelle → Lese-Kontrakt → Alphabet → Port, damit endet jede fehlerhafte Anfrage unabhängig vom Ziel mit demselben Fehler wie im Bestand. Tests und M1/M2 binden das. Spec-Qualifier „bei sonst gültiger Anfrage“ in [`SPEC-022`](../../spec/pflichtenheft.md) und [`SPEC-031`](../../spec/pflichtenheft.md) gelesen. Die Review nannte die Frage A-1 ausdrücklich eine Architect-Sache; die Entscheidung traf der Hauptlauf |
| F-2 (MEDIUM) Godoc widerspricht dem Nachbarsatz | **geschlossen** | Godoc von `ReadChanges` neu gelesen: der Satz „kommt unverändert zurück“ steht jetzt mit dem Zweig des ungültigen Ziels zusammen und nennt `Validate` als die gemeinsame Prüfung; kein Widerspruch zum Code |
| F-3 (LOW) Paritätstest: Port per Listen/Close, Start-Fehler verworfen | **geschlossen für die Sichtbarkeit, das Zeitfenster bleibt** | `startGRPC` liefert den Start-Ergebnis-Kanal, `failWithStartError` nennt den Start-Fehler in der Meldung; `freeLoopbackAddr` wählt den Port weiter per Listen/Close (Godoc nennt es ehrlich). Ein belegter Port zeigt sich jetzt benannt, nicht als Zeitüberschreitung |
| F-4 (LOW) Godoc des Capture-Tests breiter als die Strecke | **geschlossen** | Godoc und Plan-Zeile nennen „von der Transaktion … zum Stream-Port-Fake“, der Assembler-Schritt ist ausdrücklich ausgenommen; das Risiko „Kommt das Label an den Handler?“ in §6 bleibt als „bei der Closure einzutragen“ (Planner), der Teilbeleg ist als solcher benannt |
| F-5 (LOW) Handbuch-Zählwort „zwei“, Adresse ohne Zählwort | **geschlossen, mit Rest (V-5)** | Adresse in [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md) nennt Zählwort und Lokatoren; im Handbuch nachgelesen: Zeile 1301 („trägt zwei optionale, unabhängig setzbare Felder“), 1364–1365 („zwei“ am Ende von 1364, „optionale Parameter“ in 1365, C#), 1378 (Kotlin), 1403 (Python), 1456–1457 (`ReadChanges`-Filterliste ohne `target`), 1557 („Zwei optionale“, SSE) stimmen. Rest: Zeile 1149 (Abschnitt HTTP-/JSON-API, „optional gefiltert über `schema` und `table`“) steht nicht in der Liste |
| F-6 (INFO) `grpcadminclient -target` ungetestet | **als Grenze geführt** | Plan §3 trägt die Zeile, der E2E-Lauf ist die Adresse ([`slice-routing-e2e`](../plan/planning/in-progress/slice-routing-e2e.md)) |
| F-7 (INFO) Paritätstest: Fake bildet die Auswahl nach | **geschlossen, Lesehinweis (siehe §6)** | Der Testname trägt „Wege liefern für dieselbe Eingabe dieselben Changes“, nicht „die SQL-Auswahl ist gleich“; die reale Auswahl trägt der Store-Test (S1, S2 rot). Der Fake ruft seit der Fixrunde `ChangeQuery.Validate` wie der reale Store |
| F-8 (LOW) Umbruchrest im Kommentar von `SelectChanges` | **geschlossen** | Absatz in `queries.go` neu gelesen: „…trifft nie ein gesetztes `$6`.“ / „Der Join auf …“ steht in eigenen Sätzen |

## 6. Entscheidungs-Konformität und Inhalt der Fixrunde `8c3d6e2e`

- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Teilfrage 5 (ein Parameter je Weg, kein Eingriff in `ChangeStreamPort`/`Broadcaster`): getragen, Diff berührt beide nicht.
- [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md) Festlegung 1 (`ReadChangesRequest.target = 7`, Konjunktion, ein Use Case für beide Wege): getragen (Zeile 3 der DoD).
- [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Festlegung 2 (`target` außerhalb des Alphabets, leer, kein `400`, Prüfung im Use Case ohne Store-Aufruf, SQL-Zugriff ausgenommen): **in der Sache getragen**; der Wortlaut „leer … kein `400`“ der ADR (und der Test-Zeile 224 dort) ist unbedingt, die Fixrunde macht ihn bedingt („bei sonst gültiger Anfrage“). Die Begründung der ADR trägt die Bedingung (Befund-Abschnitt der ADR: „nur Parameterform … und Bereichsfehler enden `400`“; ihr Vorbild ist das Verhalten von `schema`/`table`, deren unbekannter Wert ebenfalls einen Bereichsfehler nicht verdeckt). Ein Widerspruch zu einer `Accepted` Aussage liegt damit nicht vor, aber die Lücke zwischen dem Wortlaut und der Spec-Zeile ist real.
- **Bestandsverhalten ohne `target`:** unverändert. Ablesbar am Code (Zweig nur bei gesetztem Ziel; sonst derselbe Port-Aufruf mit leerem Zusatzfeld) und an den unveränderten Bestandstests (nur Ergänzungen um `Target`, kein gelockerter Vergleich — alle entfernten Testzeilen des Diffs haben eine erweiterte Entsprechung).
- **Eine Stelle für die Validierung:** `ChangeQuery.Validate` ist einmal definiert; zwei Aufrufer (Use Case im Zweig des ungültigen Ziels, Store für alles andere) rufen dieselbe Funktion. Kein Doppel der Logik. Das Preis-Detail: für ein gültiges Ziel läuft die Prüfung erst im Store; die Reihenfolge der Fehler bleibt trotzdem gleich, weil beide Pfade `Validate` in derselben Position vor jedem Zugriff rufen.
- **Spec-Qualifier:** [`SPEC-022`](../../spec/pflichtenheft.md) (Zeile „Zustellziel“) und [`SPEC-031`](../../spec/pflichtenheft.md) (Zeile `ReadChanges`) tragen „bei sonst gültiger Anfrage“ und den Vorrang des Lese-Kontrakts; [`SPEC-020`](../../spec/pflichtenheft.md) und [`SPEC-021`](../../spec/pflichtenheft.md) bleiben unqualifiziert, zu Recht (kein Lese-Kontrakt auf den Streams). Sprachliche Schwäche in der `SPEC-031`-Zeile: der Satzteil „das gilt auch für einen Wert mit dem Zeichen U+0000“ steht nach dem eingefügten Vorrang-Satz und hat seinen Bezug (die leere Liste) verloren (V-4).
- **Paritätstest:** `parityStore` ruft jetzt `ChangeQuery.Validate` wie der reale Store. Der Satz des Tests („`GET /changes` und RPC reichen für dieselbe Eingabe dasselbe über denselben Use Case an den Port“, plus der Fehlerfall-Satz) bleibt wahr, und M1/M2 zeigen, dass die Fehlerfälle mit **ungültigem** Ziel den Use Case binden. Die Fehlerfälle mit **gültigem** Ziel binden dagegen nur den Fake (seine `Validate`-Zeile); den realen Store trägt dafür der bestehende Store-Test. Das ist die richtige Aufteilung, aber die Zeilen „gültiges Ziel“ sind Kontrolle, kein Beleg des Produktionswegs (V-6).
- **Plan-Aussagen gegen den Diff:** (a) Kopf des Plans und DoD-Zeile 7 sagen „dieser Slice ändert die Spec nicht“ bzw. „entfällt“ — der Diff ändert drei Spec-Zeilen und die Geschichte-Zeile (V-1). (b) DoD-Zeile 3 sagt „liefert auf allen Lesewegen eine leere Antwort, keinen Fehler“ unqualifiziert, ebenso [`welle-routing`](../plan/planning/welle-routing.md) in der Zeile zum Alphabet (V-2).
- [`AGENTS.md`](../../AGENTS.md) §3.12/§3.13: der Suchlauf trägt beide Stände und das Nichtgefundene; für die Eigenschaft der Fixrunde fehlt er (V-2). §3.7: `make kommentar-kennungen` ohne Kandidat; die neuen Kommentare tragen höchstens eine Kennung und beschreiben den Ist-Zustand.

## 7. Bewertung der Fixrunde und Re-Review

Inhaltlich gelesen: Code (`service.go`: Reihenfolge und eine Funktion), Tests (`TestReadChangesContractErrorsWinOverInvalidTarget`: drei Fehlerarten × ungültiges/leeres Ziel mit Port-Aufruf-Zähler; Paritätstest: vier Fehlerfälle über beide Wege), Spec, Plan, Handbuch-Adresse. Mutationen M1 und M2 färben die neuen Tests rot. **Kein Fehler im Verhalten gefunden.**

**Re-Review: ja — und nicht nur als Formsache.** Begründung:

1. Die Fixrunde ändert **Anweisungen** (neue Verzweigung samt `Validate`-Aufruf im Produktionscode) und eine **Norm** (Qualifier in zwei Spec-Zeilen). Die Entscheidung, die sie umsetzt, hat der Hauptlauf getroffen, den die Review selbst als Architect-Sache gekennzeichnet hat; der Leser derselben Entscheidung ist bisher der Entscheider.
2. Das Register [`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/observation.md) trägt zwei Gegenbelege, in denen ein Re-Review nach einer Fixrunde mit Anweisungen etwas fand, das die Verifier-Lesung nicht fand — einmal ein MEDIUM (Widerspruch zum Bestand). Hier ist die Wahrscheinlichkeit nicht geringer: Verhalten an einer öffentlichen Schnittstelle, Norm im Spiel.
3. Die DoD-Zeile 5 („kein offenes HIGH/MEDIUM“) hängt sonst allein an einem Verifier-Beleg über Text und Diff.

Umfang des Re-Reviews: `git diff 11cff8d6 8c3d6e2e` (acht Dateien, +174/−31). Leitfragen: (a) trägt die Reihenfolge (leere Quelle → Lese-Kontrakt → Alphabet) jeden Fall der drei Fehlerarten? (b) widerspricht der Spec-Qualifier einer Zusage des Lastenhefts? (c) Wording der `SPEC-031`-Zeile (V-4). Der Regel-Kandidat des Registers („Re-Review verlangen, sobald die Fixrunde Anweisungen ändert“) trägt genau diesen Fall; dieses Auftreten ist ein weiterer Beleg dafür, dass er greift, kein viertes Auftreten der Klasse „ungelesen“ (der Re-Review ist die Abhilfe, nicht der Befund).

**Offene Architect-Frage (Folge-ADR?) — Bewertung, keine Entscheidung.** Der Implementer vermerkt zu Recht, dass die Rangfolge nur in Plan und Spec steht. Meine Einschätzung:

- **Die Spec genügt als führende Stelle.** Sie steht im Rang über jeder ADR (Source Precedence), die Rangfolge ist eine Zusage an die Umsetzung (Fehler-Antwortform und Wertebereich), keine Architektur-Entscheidung zwischen Alternativen mit Preis. `ADR-0139` Festlegung 2 wird nicht inhaltlich überschrieben: ihre Begründung trägt die Lesart (siehe §6), die Festlegung schweigt zur Rangfolge. Eine Folge-ADR ist deshalb **nicht erforderlich**.
- **Was fehlt, ist die Bestätigung durch die richtige Rolle,** nicht ein neues Dokument. Eine kurze Bestätigung des Architect (oder eine Zeile des Planners mit Verweis auf die Begründung aus der ADR) schließt die Lücke zu A-1. Wählt er Option (a) („das ungültige Ziel gewinnt“), müsste die Fixrunde zurückgenommen werden — deshalb gehört die Bestätigung **vor** die Closure.
- **§3.5 und §3.12:** die `Accepted` ADR bleibt unberührt (kein Zitat-Korrektur-Fall). Der Satz im Plan §6 („berührt `ADR-0139` Festlegung 2 nicht inhaltlich“) ist eine Tatsachenbehauptung mit dem Beleg im ADR-Text (Befund-Abschnitt, Begründung); sie wäre durch diese Verweisstellen zu stützen, nicht nur zu behaupten.
- **Folge-ADR nur dann,** wenn der Architect die Rangfolge als eigene, auch für künftige Lesewege geltende Regel festschreiben will (etwa für weitere Filter-Dimensionen mit Alphabet-Prüfung). Für den Gegenstand dieses Slice ist das nicht nötig.

## 8. Benannte Grenzen — ehrlich geführt?

| Grenze | Befund |
|---|---|
| Altserver ignoriert das Feld | **ehrlich geführt:** gemessen an der Protobuf-Bibliothek dieses Repos (`TestStreamChangesRequestIgnoriertUnbekannteFelder`, beide existieren), am ausgelieferten Altserver *hergeleitet* (Plan §6) |
| Beleg am laufenden System | Gehört [`slice-routing-e2e`](../plan/planning/in-progress/slice-routing-e2e.md); ich habe `make test-integration` nicht gefahren. Die Gesamtkette vom Assembler über den Broadcaster bis zum Handler ist nicht gefahren (der Capture-Test endet am Stream-Port-Fake; Plan benennt es) |
| `grpcadminclient -target` | benannt, bis zum E2E ungetestet (F-6) |
| Konjunktion mit `schema` am Stream-Handler | trägt allein der Domänentest (Bestand von [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)), von mir nicht an den Handlern mutiert — *hergeleitet* aus dem Review |
| Port-Wahl im Paritätstest | Zeitfenster bleibt (Listen/Close), Start-Fehler sichtbar (F-3) |

## 9. Befunde

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | MEDIUM (Plan-Konformität) | Der Plan sagt zweimal, dieser Slice ändere die Spec nicht (Kopf; DoD-Zeile „Doku-Update: `spec/pflichtenheft.md` entfällt“); die Fixrunde ändert [`SPEC-022`](../../spec/pflichtenheft.md), [`SPEC-031`](../../spec/pflichtenheft.md) und die Geschichte-Zeile. Dazu: die Norm-Entscheidung stammt vom Hauptlauf, die Review hatte sie dem Architect zugewiesen (A-1) | Plan Kopf und DoD-Zeile 7; Spec |
| V-2 | LOW | Die unbedingte Aussage „leer, keinen Fehler“ steht weiter unqualifiziert in der DoD-Zeile 3 des Plans und in [`welle-routing`](../plan/planning/welle-routing.md); für die in der Fixrunde bewegte Eigenschaft gibt es keine Suchlauf-Zeile und keine Nachmessung (§3.13). Die ADR-Aussage selbst ist `Accepted` und bleibt | Plan DoD-Zeile 3; Welle |
| V-3 | Prozess | Re-Review der Fixrunde nötig (§7) — Anweisungen und Norm geändert, Entscheider = Leser | DoD-Zeile 5 |
| V-4 | LOW | `SPEC-031`-Zeile `ReadChanges`: Satzteil „das gilt auch für einen Wert mit dem Zeichen U+0000“ hat nach dem eingefügten Vorrang-Satz keinen eindeutigen Bezug | Spec, Zeile `ReadChanges` |
| V-5 | LOW | Handbuch-Zeile 1149 („optional gefiltert über `schema` und `table`“, Abschnitt HTTP-/JSON-API) steht nicht in der Adressliste von `slice-routing-betriebsdoku`; das Muster der Suchlauf-Zeilen („`schema`/`table`|schema/table|schema und table“) trifft die Schreibweise mit Backticks um beide Wörter nicht | Handbuch Zeile 1149; Suchlauf-Muster |
| V-6 | INFO | Die Zeilen „gültiges Ziel“ der Fehlerfälle im Paritätstest binden nur den Fake (dessen `Validate`-Zeile), nicht den Produktionspfad; der Store-Test trägt ihn | Paritätstest |
| V-7 | INFO | `make doc-immutable` ohne `RANGE` bricht mit Exit 2 ab; im Verifier-Beleg des Repos ist die Aufrufform mit `RANGE` zu führen | Sensor-Aufruf |

Kein HIGH. V-1 ist MEDIUM als Plan-Konformität, nicht als Defekt im Code.

## 10. Verdikt

**DoD getragen: ja für die Zeilen 1, 2, 4; Zeile 3 getragen mit Qualifier und Plan-Widerspruch (V-1, V-2); Zeile 6 getragen für die erste Eigenschaft, Lücke für die der Fixrunde; Zeile 7 teilweise (V-1); Zeile 5 offen (V-3); Zeilen 8 bis 13 offen, gehören dem Planner.** Die Review-Findings F-1 bis F-8 sind in der Sache geschlossen (F-1 mit offener Entscheidungs-Instanz, F-6 als Grenze benannt, F-3 bleibt als Zeitfenster, nicht als Unsichtbarkeit). Die Mutationen (14 gültige Läufe, alle rot), die Sensor-Läufe und die Bestandstests tragen das Verhalten; ohne `target` ist es unverändert.

**Nötiger Nachzug:**

1. **Reviewer:** Re-Review der Fixrunde `8c3d6e2e` (Diff `11cff8d6..8c3d6e2e`, Leitfragen §7) — ja.
2. **Architect (oder Planner mit Verweis auf die ADR-Begründung):** Bestätigung der Rangfolge „Lese-Kontrakt gewinnt“ vor der Closure; eine Folge-ADR ist dafür **nicht** nötig (§7).
3. **Planner:** Kopf und DoD-Zeile 7 des Plans an die Spec-Änderung der Fixrunde anpassen (V-1); die unbedingte Aussage in DoD-Zeile 3 und in der Welle qualifizieren und eine Suchlauf-Zeile für die zweite Eigenschaft ergänzen (V-2); §6-Ausgänge eintragen, Closure-Notiz mit Lerneintrag, Register-Vermerk zu [`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/observation.md).
4. **Spec-Verfasser:** V-4 (Satzbezug in der `ReadChanges`-Zeile).
5. **Folge-Slice [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md):** Handbuch-Zeile 1149 in die Adressliste aufnehmen (V-5); Aussage „leer, kein Fehler“ dort mit dem Qualifier „bei sonst gültiger Anfrage“ führen.
