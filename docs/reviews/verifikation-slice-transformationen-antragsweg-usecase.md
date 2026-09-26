# Verifikations-Report: slice-transformationen-antragsweg-usecase — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-transformationen-antragsweg-usecase.md`](review-slice-transformationen-antragsweg-usecase.md);
Formvorbild dieses Reports:
[`verifikation-slice-transformationen-antragsweg-schema.md`](verifikation-slice-transformationen-antragsweg-schema.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`, Gerüst per `cp`
übernommen und durch die Form des Formvorbilds ersetzt).

**Gegenstand:** Slice-Plan `slice-transformationen-antragsweg-usecase` (Welle `welle-transformationen`), Stand
`HEAD` = `86f1ab20`, Diff-Range `2a47cd9a..HEAD`, 14 Commits, 33 Dateien (+4381/−559 einschließlich der beiden
Lifecycle-Moves). Slice-Inhalt: Lifecycle/Plan (`277a1d8b`, `77f75263`, `80eefead`), Implementer-Lauf (`fef790f3`,
`27b3e3c7`, `a8c3d004`, `190417ee`, `ea14cd18`), Review-Report (`c7796d83`, gelesen bis `ea14cd18`), Fixrunde 1
(`cd4dae0e`, `b0506a2b`), Architect-Zug [`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md)
mit dem [`SPEC-019`](../../spec/pflichtenheft.md)-Nachzug (`a69d3853`), Fixrunde 2 (`c3d68df6`, `86f1ab20`). Fixrunde 1,
Fixrunde 2 und `ADR-0127` sind von keinem Reviewer gelesen (der Review-Report ist seit `c7796d83` unverändert); ihre
Wirkung belegen die eigenen Mutationen in §4. Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen
Report. Die Mutationen liefen an Kopien des Baums im Scratchpad (`git archive HEAD`, Ersetzung mit `sed` in eine
Ausgabedatei, kein `sed -i` am Repo); der Arbeitsbaum des Repos blieb unberührt (`git status --short` leer, bis auf
diesen Report).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert, die Logs danach
gelesen. Stand aller Läufe ohne Mutation: `HEAD` = `86f1ab20`, Arbeitsbaum sauber. Es lief ein schwerer Docker-Lauf
zugleich (`free -m` vor dem ersten Lauf: 12,5 GB verfügbar; zwei fremde Container mit Zufallsnamen liefen während der
Läufe, nicht von mir).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test` (Race-Detector) | **EXIT=0** | 44 Zeilen `ok`, keine Zeile `FAIL` |
| `make test-store` | **EXIT=0** | `ok … postgresstorage 9.446s coverage: 78.9% …`, `DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile gemergt: store,replication)`, `db-coverage: OK — DB-Adapter-Coverage 82.56% erfuellt Schwelle 80%`; der Rollout des Laufs druckt sieben `CREATE FUNCTION` |
| `make a-check` | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **EXIT=0** | `coverage-gate: OK — Coverage 84.80% erfüllt Schwelle 80%` |
| `make gates` (einmal) | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1233 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` · `coverage-gate: OK — Coverage 84.80% erfüllt Schwelle 80%` · `gesamt: 0 Befund(e)` (a-check) |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **EXIT=0** | `suchlauf-nachmessen: 32 Zeilen stimmen` |
| `make commit-traceability RANGE=2a47cd9a..HEAD` | **EXIT=0** | `OK — 14 Commit(s) in "2a47cd9a..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=2a47cd9a..HEAD` | **EXIT=0** | `d-check: 1233 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=2a47cd9a..HEAD` | **EXIT=0** | `d-check: 1233 Datei(en) geprüft, 0 Befund(e)` (ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab — Bedienung, kein Befund) |
| `tools/harness/run-schema-rollout-guard-test.sh` | **EXIT=0** | alle sechs Läufe, siehe unten |
| gofmt (Docker, gepinntes Toolchain-Image, `gofmt -l` über die 25 geänderten `*.go`-Dateien) | Exit 0 | eine Datei genannt: `queries/queries.go` — am Parent `2a47cd9a` bereits formabweichend (Umschreibung von `''` in einem Doc-Kommentar der Retention-Abfrage, `gofmt -d`), nicht Teil dieses Diffs; die drei Dateien aus F-3 des Reviews sind konform |
| Mutationen (§4) | 7 Go-Mutationen (G1–G7) und 8 Store-Mutationen (S1–S8) an der Eingabeseite, dazu drei Läufe mit einem Zusatztest (R0–R2) | siehe §4 |

**Guard-Skript (Tag und gedruckte Exit-Codes).** Sechs Läufe, Exit 0 (die SQL-Funktionen wurden in Fixrunde 2 geändert).
Lauf 5 (Alt-Tag) druckt: „Lauf 5 OK — Tag v0.2.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, ohne Vorlauf), Exit 0
(Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar; 19 Tabellen-/View-Rechte der drei Rollen auf
administration_request, backfill_run und backfill_status, EXECUTE auf 3 Funktionen (cdc.backfill_table(text, text, text),
cdc.set_transformation(text, text, text, text, json), cdc.remove_transformation(text, text, text, text)) allein für
cdc_admin (nicht PUBLIC), Spalten rule_name:text:YES,rule_spec:jsonb:YES (Alt-Zeile alttag-req: NULL),
request_kind-Menge backfill,disable,enable,exclude_column,include_column,remove_transformation,set_transformation,
cdc.set_transformation/cdc.remove_transformation unter cdc_admin schreiben pending-Anträge, cdc_reader: permission denied
for function“; Schlusszeile: „OK — alle Belege real erbracht (Idempotenz-Allow, echte Änderung bleibt wirksam,
View-Signatur-Vorlauf, Alt-Tag v0.2.0, Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit
d-migrate-Exit 8)“. Die Läufe 1 bis 4 und 6 druckten ihre Kopfzeilen ohne `FEHLER`. Der Alt-Tag ist `v0.2.0` (der jüngste
`v*`-Tag laut Skript). Damit ist die als hergeleitet geführte Einordnung von
[`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) §Konsequenzen („`CREATE OR REPLACE` mit
gleicher Signatur, Alt-Tag-Lauf“) an diesem Lauf erprobt; `git status` danach: `tools/schema/plan.yaml` und
`tools/schema/down.sql` unverändert.

Hygiene: dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach dem letzten
**34**, kein `prune`, kein `system prune`. Nicht gefahren: `make test-replication`, `make test-integration`,
`make image`, `make bench` — der Diff berührt weder den Replication-Pfad noch den Compose-Rundlauf noch den
Build-Kontext außerhalb von `internal/` und `tools/schema/`; die Rollout-Kette deckt das Guard-Skript, die reale
Datenbank `make test-store`.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: acht `[x]`-Zeilen, vier `[ ]`-Zeilen, zusammen zwölf (`grep -c '^- \[x\]'` 8, `grep -c '^- \[ \]'` 4).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | K1–K4 gegen Fakes je an ihre Eingabe gebunden; Fehlertext der Spec; Regelstand nach `failed` unverändert; Prüfreihenfolge; Form von `column`/`to` auf `rule_spec ist ungültig`, `to` gleich `column` auf den K3-Text; Regelname ungültig; Übergabe „Queue mit ungültiger und gültiger Zeile“ samt der Formen von `rule_spec`; Mutation je Invariante | **erfüllt** | `make test` Exit 0. Use Cases und Domäne gelesen (`settransformation/service.go`, `transformationspec.go`): Reihenfolge `CheckRuleName` → `ParseTransformationSpec` (a bis d) → erster Port-Zugriff → K1 → K2 → K3 (Regelziel, dann Spalte) → K4; Adressen je Text gegen die Tabelle von [`SPEC-019`](../../spec/pflichtenheft.md) geprüft. Mutationen G1 (K1), G2 (K4 vor K1–K3), G3 (K3 Ziel gleich Quelle), G4 (Regelname-Prüfung) rot (§4). Stelle der Prüfung leerer Regelfelder: **Verarbeiten**; `NewAdministrationRequest` nimmt die Zeile, der Use Case bestimmt den Text (Diff gelesen); der Login-Test läuft real in `make test-store` (Exit 0) |
| 2 | Regelstand-Ableitung und Spaltenliste im Store; Ordnung `requested_at`, bei gleichem Zeitstempel `administration_request_id`; Faltung als reine Funktion in der netzlos gemessenen Fläche | **erfüllt** | `make test-store` Exit 0, `make test` Exit 0. `FoldTransformations` steht in `internal/domain/model` (a-check und coverage-gate zählen die Fläche, Coverage 84.80 %). Statusfilter (S1) und Zweitschlüssel der offenen Queue (S2) rot (§4); die Ordnung der Ableitungs-Abfrage mit dem Zweitschlüssel ist im Review (M25) rot gesehen (übernommen) |
| 3 | Verdrahtung und Dauerhaftigkeit: Nachtrag live ohne Neustart, Remove, Prozessstart, Aktivierungs-Zweig, Idempotenz der Wiederholung, Tabelle ohne Bindung | **erfüllt** | `make test`, `make test-store` Exit 0. Mutationen G5 (Nachtrag `SetTransformation` gestrichen), G6 (Prozessstart, Schlüssel fremd), G7 (Aktivierungs-Zweig, Schlüssel fremd) rot (§4); Idempotenz an der Store-Abfrage (S1) rot im Login-Test |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt, mit V-2 bis V-4 am Rand** | Report liegt vor (0 HIGH, 2 MEDIUM, 5 LOW, 4 INFO; Verdikt merge-blockierend); Finding für Finding in §5 nachgemessen — kein offenes HIGH/MEDIUM; Fixrunden ohne zweiten Reviewer-Lauf (V-6) |
| 6 | §3.13-Suchlauf: committetes Feld in §3, Gefundenes und Nichtgefundenes je Träger, beide Stände | **erfüllt** | 32 Zeilen mit dem Werkzeug Exit 0 (§1); zehn Zeilen von Hand nachgefahren (§6); die Befund-Zellen stimmen |
| 7 | Doku-Update: entfällt, Aufschub mit Adresse `slice-transformationen-betriebsdoku` | **erfüllt (Aufschub mit Adresse, die den Gegenstand trägt)** | Handbuch-Kandidatenlauf §7; die Adresse trägt die Aussage zur Aufruf-Reihenfolge (Diff von `slice-transformationen-betriebsdoku` gelesen) |
| 8 | Reconciliation-Register — entfällt (Greenfield) | **erfüllt (entfällt)** | keine Datei |
| 9 | Closure-Notiz mit Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner) |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: die Zeilen tragen „Ausgang: *(bei Closure …)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure der Welle |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre. Register, Risiko-Ausgänge, Paarungen
und Closure-Notiz sind Planner-Arbeit und nicht Teil dieser Prüfung.

## 3. Plan-vs-Code-Diff

**§1 Ziel:** zwei Use Cases mit K1–K4, ein Outbound Port mit zwei Lese-Methoden, Verdrahtung in
`applyAdministrationRequest`, Regelstand in `activatedTableBindings` und im Aktivierungs-Zweig stehen. Die
„Ausdrücklich NICHT“-Punkte sind eingehalten: kein Schema-/Funktions-Eingriff außer dem Zeitstempel der Fixrunde 2
(Architect-Zug, im Plan §3 benannt), kein Eingriff in `mapper` (Assembler), kein Backfill-Pfad (`git diff --name-only`
nennt kein `backfill`-Paket außerhalb des Bootstrap-Testfixtures `backfill_internal_test.go`, ein Zeilen-Nachzug), keine
Änderung der Startreihenfolge, keine Sicht/`diagnose`-Ausgabe.

**§3-Tabelle, Zeile für Zeile gegen `git diff --stat 2a47cd9a..HEAD`:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `port/inbound/transformation.go` neu | 55 Zeilen, `SetTransformationUseCase`/`RemoveTransformationUseCase` samt Commands |
| `port/outbound/transformation.go` neu | 37 Zeilen, **ein** Port `TransformationPort` mit `TransformationRules` und `SourceColumns` |
| `usecase/settransformation`, `usecase/removetransformation` (+ Tests) neu | beide mit `service.go` und `service_test.go` |
| `domain/model/transformationspec.go` (+ Test) neu | `ParseTransformationSpec`, `CheckRuleName`, `CheckConflicts`, `Build`, `FoldTransformations` |
| `administrationrequest.go` (+ Test), `errors.go` update | Konstruktor ohne Ablehnung leerer Regelfelder, `AdministrationRequestKinds()`, Sentinels (+46 Zeilen) |
| `queries.go`, `sqlexec/translate.go`, `tableactivation.go` (+ Tests) update | `SelectAppliedTransformationRequests`, `SelectTableColumns`, `ReadTransformationRules`, `ReadSourceColumns`, `TableActivationAdapter` implementiert den Port |
| `wiring.go` update | zwei Zweige, drei Felder in `administrationDeps`, `activatedTableBindings`/Aktivierungs-Zweig mit Regelstand, `processedAdministrationKinds()` als Ableitung (Diff gelesen) |
| Bootstrap-Tests (`administration_internal_test.go`, `wiring_rest_…`, `backfill_…`, `administration_roles_…`) | alle im Diff |
| `.dockerignore` prüfen — keine Änderung | nicht im Diff |
| `SPEC-019`-Zeile „Domänen-Invarianten des Antrags-Konstruktors“ — gemeldet, nicht geändert | Zeile 675 von `spec/pflichtenheft.md` unverändert, Meldung im Plan §3 mit Frist „Closure dieses Slice“ (Bewertung V-4) |
| Spalten-Antragsarten mit leerer Spalte — Grenze, benannt | Konstruktor lehnt sie weiter beim Lesen ab, Test-Bindung `TestNewAdministrationRequestRejectsInvariantViolations` (Fall „leere Kennung“) und `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable` (F-6) |
| Fixrunde 1: `queries.go`/`outbound/administrationrequest.go` (F-1, Ordnung), `administrationrequest_order_test.go`, `translate_test.go` (F-6), `wiring.go` (F-2, F-3, F-7), `administration_internal_test.go` (F-4), `transformationspec.go`/`administrationrequest_test.go` (gofmt) | alle im Diff |
| Fixrunde 2: `nacharbeit-administration.sql` (sieben Funktionen, Kopfkommentar), Kommentare in `queries.go`/`administrationrequest.go`/`translate.go`/Port, `administration_callorder_internal_test.go`, `slice-transformationen-betriebsdoku` | alle im Diff; der `INSERT` jeder der sieben Funktionen trägt `requested_at` als zweite Spalte, `clock_timestamp()` als zweiten Wert (`git grep -n clock_timestamp -- tools/schema/nacharbeit-administration.sql`: sieben Wertzeilen plus Kopfkommentar); Signaturen unverändert |
| Handbuch `docs/user/benutzerhandbuch.md` — keine Änderung | nicht im Diff (`git diff --name-only 2a47cd9a HEAD -- docs/user`: 0 Zeilen) |

Nicht im Plan-Feld, im Diff, mit Herkunft im Kopf des Plans: [`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md),
der ADR-Index, `SPEC-019` (Architect-Zug) und der Review-Report. Im Plan, nicht im Diff: nichts.
Keine `Accepted` ADR wurde inhaltlich geändert (`make doc-immutable RANGE=2a47cd9a..HEAD` Exit 0). Keine Schwelle, keine
Gate-Konfiguration verändert ([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6).

**Festlegungen des Plans gegen den Code:**

| Festlegung | Ist |
|---|---|
| Stelle der Prüfung leerer Regelfelder: Verarbeiten, nicht Lesen | `NewAdministrationRequest` nimmt `set_transformation`/`remove_transformation` ohne Ablehnung; einziger Verbraucher von `RuleName`/`RuleSpec` ist der Set-/Remove-Zweig |
| ein Port mit zwei Lese-Methoden | `TransformationPort`, dieselbe Adapter-Instanz wie `ColumnExclusionPort` |
| K4 gegen die Spaltenliste (nicht `ColumnExists`) | `SourceColumns` (`information_schema.columns`, `ORDER BY ordinal_position`), zugleich Grundlage von K3; [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 3 nennt `ColumnExists` als Muster, der Plan §3 begründet den Schnitt |
| Set trägt vor dem Vermerk `applied` nach, Ersatz nach Namen macht die Wiederholung folgenlos | `deps.assembler.SetTransformation` vor dem Vermerk; Grenze (Vermerk-Fehler und kollidierender Folgeantrag) im Kommentar am Set-Zweig und in Plan §6 |
| `requested_at` ist der Aufrufzeitpunkt, Verarbeitung und Ableitung in derselben Ordnung | `clock_timestamp()` in sieben `INSERT`, `ORDER BY requested_at, administration_request_id` in drei Abfragen (`git grep -n 'ORDER BY requested_at' -- internal ':!*_test.go'`: vier Treffer, einer davon die Aufnahme der Backfill-Runs mit `run_id`) |

## 4. Mutationen der Eingabeseite (dieser Lauf)

Go-Mutationen G1–G7 an einer Kopie, Einzellauf `go test` im gepinnten Toolchain-Image ohne Netz über die Pakete
`settransformation`, `model` und `bootstrap` (G5–G7 nur `bootstrap`); Store-Mutationen S1–S8 und R0–R2 mit
`make test-store` in einer Kopie (eigener Lauf je Mutation, die Kopie wurde nach jedem Lauf per `cp` auf den
unmutierten Text zurückgesetzt). Ausgangslauf der unveränderten Bäume: Exit 0 (§1; R0 unten für den Zusatztest).

| # | Zusage | Mutation | Ergebnis (gedruckt) |
|---|---|---|---|
| G1 | K1: ein vergebener Regelname endet `failed` | `if rule.Name() == name` → `if false && …` (`CheckConflicts`) | **rot**: `TestSetTransformationRejectsWithTheSpecTexts`, `TestCheckConflictsBindsK1ToK3AndTheirOrder`, `TestProcessAdministrationRequestsRuleViolationsFailWithSpecTexts` |
| G2 | Prüfreihenfolge: K1–K3 vor K4 | K4-Prüfung (`Spalte existiert nicht an der Quelle`) vor `CheckConflicts` gesetzt | **rot**: `TestSetTransformationRejectsWithTheSpecTexts` |
| G3 | K3: Zielname gleich Quellspalte | `s.to == s.column ||` entfernt | **rot**: `TestSetTransformationRejectsWithTheSpecTexts`, `TestCheckConflictsBindsK1ToK3AndTheirOrder` (`bootstrap` grün: dort deckt die Spaltenliste den Fall) |
| G4 | Formzeile Regelname vor dem ersten Lesen | `CheckRuleName`-Fehler ignoriert (`false && err != nil`) | **rot**: vier Tests, darunter `TestSetTransformationChecksTheFormBeforeReadingTheStore`, `TestProcessAdministrationRequestsInvalidRuleRowsDoNotStallTheQueue` |
| G5 | Nachtrag in die laufende Bindung | `deps.assembler.SetTransformation(qualified, rule)` → `_ = rule` | **rot**: vier Tests, darunter `…SetAndRemoveTransformationTakeEffectLive`, `…SetTransformationIsIdempotent` |
| G6 | Prozessstart trägt den Regelstand | `Transformations: rules[table.QualifiedName()]` → `rules["public.other"]` (die Streichung der Zeile endet am Übersetzer, `declared and not used`, und zählt nicht) | **rot**: `TestActivatedTableBindingsCarriesTransformations` |
| G7 | Aktivierungs-Zweig trägt den Regelstand | `Transformations: rules[qualified]` → `rules["public.other"]` | **rot**: `TestProcessAdministrationRequestsDisableEnableCycleRestoresTransformations` |
| S1 | K1 prüft nur gegen `applied`-Zeilen (Grundlage der Idempotenz) | `status = 'applied'` → `status <> 'failed'` in `SelectAppliedTransformationRequests` | `make test-store` **EXIT=2**: `TestAdministrationPathRunsUnderLeastPrivilegeLogins` („Regelname bereits vergeben: public.admin_roles_path.roles_rule“), `TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived` |
| S2 | die Queue ordnet Gleichzeitige nach `administration_request_id` | Zweitschlüssel aus `SelectPendingAdministrationRequests` gestrichen | **EXIT=2**: `TestAdministrationRequestListPendingOrdersTiesByRequestID` („ListPending-Ordnung der Gleichzeitigen = [queue-order-b-set queue-order-a-remove queue-order-d-include queue-order-c-exclude]“) |
| S3 | [`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md): `set_transformation` trägt den Aufrufzeitpunkt | `clock_timestamp()` → `now()` in `cdc.set_transformation` | **EXIT=2**: `TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived` („Row Image der laufenden Bindung = {"id":"1","secret":"geheim"}, erwartet {"id":"1","second":"geheim"}“) |
| S4 | dieselbe Zusage an `remove_transformation` | `clock_timestamp()` → `now()` in `cdc.remove_transformation` | **EXIT=0**, kein Rot (V-1, V-2) |
| S5 | dieselbe Zusage an `enable_table` | `clock_timestamp()` → `now()` in `cdc.enable_table` | **EXIT=2**: `TestAdministrationRequestSameTransactionCallsKeepCallOrder` |
| S6 | dieselbe Zusage an `disable_table` (der **früher** aufgerufenen Funktion ihres Paares) | `clock_timestamp()` → `now()` in `cdc.disable_table` | **EXIT=2**: `TestAdministrationRequestSameTransactionCallsKeepCallOrder` — **rot**, obwohl Fixrunde 2 die Mutation an der früher aufgerufenen Funktion als äquivalent führt (V-3) |
| S7 | dieselbe Zusage an `backfill_table` | `clock_timestamp()` → `now()` in `cdc.backfill_table` | **EXIT=0**, kein Rot (V-2) |
| S8 | dieselbe Zusage an `exclude_column` (früher aufgerufene Funktion ihres Paares) | `clock_timestamp()` → `now()` in `cdc.exclude_column` | **EXIT=2**: `TestAdministrationRequestSameTransactionCallsKeepCallOrder` — rot |
| R0 | Zusatztest der Verifikation (nur in der Kopie, nicht im Repo): eine Transaktion mit den Aufrufen `set`, `remove`, `include`, `exclude`, `enable`, `disable`, `backfill` (Gegenrichtung der drei Paare, `backfill_table` zuletzt), 60 Transaktionen, Erwartung `ListPending` = Aufruf-Reihenfolge | unmutierte Funktionen | **EXIT=0** — das Produkt hält die Ordnung |
| R1 | Zusatztest wie R0 | `clock_timestamp()` → `now()` in `cdc.remove_transformation` | **EXIT=2**: „Durchlauf 0: ListPending-Ordnung [remove set include exclude enable disable backfill], Aufruf-Reihenfolge [set remove include exclude enable disable backfill]“ |
| R2 | Zusatztest wie R0 | `clock_timestamp()` → `now()` in `cdc.backfill_table` | **EXIT=2**: „ListPending-Ordnung [backfill set remove include exclude enable disable], Aufruf-Reihenfolge [set remove include exclude enable disable backfill]“ |

Aus S4/S6/S8 gegen R1: die Farbe einer Mutation an `clock_timestamp()` hängt an der **Stellung** des Aufrufs in der
Test-Folge (`remove`, `set`, `exclude`, `include`, `disable`, `enable`), nicht an der Rolle „früherer/späterer Aufruf
eines Paares“. Nur `remove_transformation` steht an erster Stelle: vor ihm liegt kein Aufruf, dessen
`clock_timestamp()` die Ordnung stören könnte, deshalb bleibt S4 grün; der Zusatztest R1 zeigt, dass die Mutation für
die Folge „erst `set`, dann `remove`“ nicht äquivalent ist. Die Mutationen des Reviews (M1–M26) habe ich nicht
wiederholt; sie sind **übernommen** aus dem Review-Report (dort an Eingabeseite gefahren, je rot). Meine G1–G7 und S1
sind Stichproben derselben Zusagen an einem neuen Stand.

## 5. Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) Verarbeitungs-Ordnung ungleich Ableitungs-Ordnung bei gleichem Zeitstempel | Zwei Züge gelesen und gefahren: Fixrunde 1 (Zweitschlüssel in `SelectPendingAdministrationRequests`, S2 rot) und der Architect-Zug [`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) mit Fixrunde 2 (`clock_timestamp()` in sieben Funktionen; S3, S5, S6, S8 rot). Live/abgeleitet-Gleichheit selbst geprüft: `TestAdministrationRequestSameTransactionCallsKeepCallOrder` (60 Transaktionen über die realen Funktionen, `ListPending` = Aufruf-Reihenfolge, Regel- und Ausschlussstand der Verarbeitung gleich der Ableitung) und `TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived` (50 Tabellen, laufende Bindung und aus `TransformationRules` neu gebildete Bindung liefern dasselbe Bild) laufen in `make test-store` Exit 0; mein Zusatztest R0 hält die Ordnung auch für die Gegenrichtung der Paare. Die verbleibenden Grenzen (Aufrufzeitpunkt statt Festschreibung bei überlappenden Transaktionen mehrerer Sitzungen, Rückwärtssprung der Uhr, Zeilen vor der Änderung) stehen in der ADR (Festlegung 3), in `SPEC-019` („Ordnung der Verarbeitung“) und im Plan §6; **hergeleitet, von mir nicht erprobt**. Die Bindung der Zusage an die Funktionen hat Lücken (V-2) | **behoben**; die Bindung an `remove_transformation` und `backfill_table` fehlt (V-2, LOW) |
| F-2 (MEDIUM) Kommentar in `Run` beschreibt einen Mechanismus, den der Code nicht trägt | `wiring.go` gelesen: der Kommentar sagt jetzt, `CDC_TABLES` sei ausschließlich der Erstaktivierungs-Seed, die Bindungen von `cfg.Tables` gingen nur in die Aktivierung und erreichten den Assembler nie; `git grep -n 'cfg.Tables' -- internal ':!*_test.go'` bestätigt (Zuweisung, die eine Schleife, `receive.Config.Tables`); die Suchlauf-Zelle im Plan trägt denselben Befund | **behoben** |
| F-3 (LOW) drei Dateien nicht gofmt-konform | `gofmt -l` über alle geänderten Go-Dateien: nur `queries.go`, am Parent formabweichend (§1) | **behoben** |
| F-4 (LOW) Mutationsangabe am Idempotenz-Test nicht setzbar | Doc-Kommentar von `TestProcessAdministrationRequestsSetTransformationIsIdempotent` gelesen: bindet den Store-Fake, nennt die tragende Bindung (`TestAdministrationPathRunsUnderLeastPrivilegeLogins`, Statusfilter); S1 färbt genau diesen Login-Test | **behoben** |
| F-5 (LOW) Plan widerspricht sich im Port-Schnitt | Plan §3 (Zeile zu `outbound/transformation.go`) und §6 („Der Port-Schnitt folgt der ADR-Zählung“) gelesen: beide Aussagen nennen dasselbe | **behoben** |
| F-6 (LOW) Stall-Klasse: leeres Schema/Tabelle beim Lesen, leere `column` bei exclude/include | Plan §3 (Zeile „Spalten-Antragsarten mit leerer Spalte“) und §6 (dritter Punkt) gelesen: benannte Grenze mit Adresse (`BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`, Entscheidung beim Planner, Nachzug in der Closure); Ist-Verhalten der Restklasse gebunden durch `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable` und den Konstruktor-Test; der Register-Zustand (`state.md` des Eintrags) trägt noch den Pfad `open/…` und den Wortlaut „geplant“ — Planner-Nachzug bei Closure | **benannte Grenze mit Adresse** |
| F-7 (LOW) Auseinanderlaufen bei Vermerk-Fehler nur für die Wiederholung belegt | Kommentar am Set-Zweig und Plan §6 (erster Punkt) nennen die Grenze; kein Test (hergeleitet aus dem Quelltext, so im Plan benannt) | **benannte Grenze** |
| F-8 (INFO) Ableitung quellweit, `internal` statt `storage` | Plan §6 („Kenntnis aus dem Review“, Adresse `slice-transformationen-backfill-pfad`) gelesen; der Doc-Kommentar von `activatedTableBindings` sagt weiter „endet den Start wie ein Lesefehler“ (V-5, INFO) | **benannt, mit Adresse** |
| F-9 (INFO) Coverage-Spanne zu eng | Plan nennt die Spanne der acht Läufe 84,70–84,90 %; meine zwei Läufe 84.80 % und 84.80 % liegen darin | **behoben** |
| F-10 (INFO) `sed -i` am Repo | Kenntnis (Prozess); keine inhaltliche Fehlwirkung in meinen Läufen gefunden (Diff und alle Sensoren grün) | **Kenntnis** |
| F-11 (INFO) einseitige Belegtiefe (Statement-Text, Kosten je Block) | Plan §6 nennt (b) mit Adresse `backfill-pfad`; (a) ist über S1 an der Eingabeseite gebunden | **benannt, mit Adresse** |

Kein offenes HIGH, kein offenes MEDIUM.

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent `80eefead` (Slice-Start) für die ersten 16 Messungen, `a69d3853` für die drei Zeilen des Architect-Zugs. Von Hand
mit `git grep -n … | wc -l` (ohne das Werkzeug), zehn der 16 Messpaare; Stand `diff` ist der Arbeitsbaum dieses Laufs
(`HEAD` = `86f1ab20`, sauber; die Plan-Datei ist im Stand `diff` durch das Werkzeug ausgeschlossen, ich habe das Muster der
Aufrufzeitpunkt-Zeile mit `':!docs/plan/planning/in-progress/slice-transformationen-antragsweg-usecase.md'` nachgefahren):

| Zeile | Stand | Soll | Von Hand |
|---|---|---|---|
| `TableBinding{` (Zeilen 1–2) | `80eefead` / `diff` | 5 / 5 | **5** / **5** |
| `ExcludedColumns` (3–4) | `80eefead` / `diff` | 25 / 26 | **25** / **26** |
| `Transformations:` (5–6) | `80eefead` / `diff` | 0 / 4 | **0** / **4** |
| `processedAdministrationKinds` (11–12) | `80eefead` / `diff` | 5 / 3 | **5** / **3** |
| `administrationDeps{` (15–16) | `80eefead` / `diff` | 21 / 21 | **21** / **21** |
| `transformations:` in `internal/bootstrap` (17–18) | `80eefead` / `diff` | 0 / 11 | **0** / **11** |
| `ORDER BY requested_at` (23–24) | `80eefead` / `diff` | 3 / 4 | **3** / **4** |
| `Transaktionszeitstempel\|Transaktionsbeginn\|Anlage-Reihenfolge\|Anlage-Zeitpunkt` (27–28) | `a69d3853` / `diff` | 15 / 14 | **15** (17 roh, 15 ohne die Plan-Datei) / **14** |
| `clock_timestamp` in `tools internal` ohne Tests (29–30) | `a69d3853` / `diff` | 4 / 14 | **4** / **14** |
| `gen_random_uuid` in `tools internal examples cmd` | `diff` | nur `nacharbeit-administration.sql` | **7** Treffer, alle in dieser Datei |

Alle Zahlen stimmen. Die Nichtgefunden-Aussagen habe ich an drei Stellen selbst gesucht: `docs/user` trägt weder
`set_transformation` noch `remove_transformation` noch `requested_at` (0 Treffer); die Ordnung der Queue steht in
`spec/` nur in `SPEC-019`; ein zweiter Träger des Funktionstextes der sieben Funktionen besteht nicht.

## 7. Entscheidungs-Konformität

- **[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 1 (Konfigurationsmechanismus),
  Teilfrage 3 (K1–K4, statisch ausgeschlossen, `failed` mit Fehlertext, Regelstand unverändert), Teilfrage 6
  (Dauerhaftigkeit: `applied`-Zeilen als einzige Herkunft, Prozessstart und Aktivierungs-Zweig) und Folgepflicht 3
  (Use Cases, Port, Verdrahtung, `activatedTableBindings`, Aktivierungs-Zweig, „Belegt durch `make test-store`“):**
  alle Bausteine vorhanden (§2, §3). K4 nutzt die Spaltenliste statt `ColumnExists`, der Wortlaut nennt `ColumnExists`
  als Muster; der Schnitt ist im Plan begründet und belegt (Login-Test liest die Spaltenliste unter dem realen
  `cdc_admin`-Login). Kein Abweichen vom Wortlaut der Folgepflicht.
- **[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md):** der Regelstand ist ein dauerhafter,
  tabellen-scoped Träger nach demselben Muster; G6/G7 färben beide Pfade.
- **[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 5:** Existenz der Spalte an der Quelle über
  den Katalog, dieselbe Sichtbarkeit wie `ColumnExists` (im Review gelesen, Login-Test unter `cdc_admin` real).
- **[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md):** Live-Reload ohne Neustart
  (Nachtrag in die laufende Bindung unter `tablesMu`); `requested_at` als Aufrufzeitpunkt ergänzt die dort
  unbestimmte „Zeitstempel“-Aussage (siehe `ADR-0127`).
- **[`ADR-0028`](../plan/adr/0028-inbound-use-cases.md) / [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md):** zwei
  Inbound Use Cases, ein Outbound Port mit zwei Lese-Methoden, beide Methoden beantworten dieselbe Frage
  (Konfliktprüfung), keine gemeinsame Transaktion nötig; `make a-check` „gesamt: 0 Befund(e)“.
- **[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md):** keine Domänenlogik in SQL — die Funktionen
  schreiben ausschließlich die Antrags-Zeile; die Zeitstempel-Vergabe im `INSERT` ist keine Domänenlogik (Funktionstext
  gelesen).
- **[`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) / [`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md):**
  `rule_spec` als `json`-Parameter, `p_rule_spec::jsonb` im `INSERT` unverändert (Diff der Funktion: nur die
  Spaltenliste und der zweite Wert); Guard-Skript Lauf 5 druckt `…, json)` in der Signatur. Die Formen der Übergabe
  (SQL-`NULL`, JSON-`null`, Wert ohne Objekt) enden `rule_spec ist ungültig`, ein doppelter Schlüssel erreicht Go nicht
  (`jsonb`, Kommentar am Parser); der Store-Test `TestAdministrationRequestListPendingCarriesRuleRowsWithMissingFields`
  läuft in `make test-store` Exit 0.
- **[`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) (Option A: `clock_timestamp()` je Aufruf):**
  Festlegung 1 (sieben Funktionen, Spalten-Default unverändert) im Text gelesen und mit S3, S5, S6, S8 an vier Funktionen
  rot gesehen; Festlegung 2 (Ordnung der Verarbeitung gleich der Ableitung) mit S2 rot; Folgepflicht 1
  (`SPEC-019`-Nachzug im selben Commit `a69d3853`: Spalte, Absatz „Ordnung der Verarbeitung“, K1-Klausel) gelesen,
  erfüllt; Folgepflicht 2 (Funktionstext, Kopfkommentar, Kommentare in `queries.go`) erfüllt; Folgepflicht 3 (Tests,
  die zwei Zeilen der Fitness Function) erfüllt — beide Tests laufen real, die Aufrufreihenfolgen der Zeile sind die des
  Tests (`remove`+`set`, `exclude`+`include`, `disable`+`enable`); Folgepflicht 4 (Handbuch) an `slice-transformationen-betriebsdoku`
  §2 übergeben, gelesen. Messtabelle: die Zeilen 1 bis 5 nicht nachgemessen (PostgreSQL 17.11-Hälfte und Zufallsstichproben
  **übernommen** aus der ADR); ihre Aussage, `now()` in allen sieben Funktionen breche die Ordnung über die zufällige
  Kennung, ist mit meinen Läufen an sechs von sieben Funktionen gemessen stimmig (S3, S5, S6, S8; Plan-Lauf `include_column`).
  **Die Mutationsangabe der Fitness-Function-Zeile ist breiter als ihre Messung** (V-1).
- **[`SPEC-019`](../../spec/pflichtenheft.md) / [`SPEC-030`](../../spec/pflichtenheft.md):** Fehlertext-Tabelle (elf Zeilen),
  Prüfreihenfolge „fünf Formzeilen, dann K1 bis K4“, Adressen (Regelname, Spalte, Zielname), `remove_transformation`
  durchläuft nur Namenszeile und K4 — Use-Case-Code Zeile für Zeile dagegen gelesen und mit G1–G4 gebunden. Der Absatz
  „Ordnung der Verarbeitung“ beschreibt die Zusage und ihre drei Grenzen wie die ADR; die K1-Klausel „beide Aufrufe dürfen
  in einer Transaktion stehen“ ist mit dem Test `…RemoveThenSetLeavesTheNewRuleLiveAndDerived` an der realen Kette gebunden
  (S3). Die Zeile 675 („`rule_name` ist … Pflicht … (Domänen-Invarianten des Antrags-Konstruktors)“) ist ungenau (V-4).
- **[`AGENTS.md`](../../AGENTS.md):** §3.3 (die Moves `277a1d8b` und `80eefead` sind reine Renames, Inhalt in eigenen
  Commits — `git show -M --stat`), §3.5/§3.6 (keine `Accepted` ADR überschrieben, keine Schwelle gesenkt), §3.7 (der
  Kommentar am Set-Zweig nennt die Grenze im Indikativ; zwei Kommentare mit falscher Verallgemeinerung siehe V-3), §3.9
  (meine Läufe ungepiped), §3.12 (Zahlen im Plan tragen Lauf und Stand; Coverage 84.80 % im Plan und in meinem Lauf gleich,
  die Spanne 84,70–84,90 % genannt), §3.13 (fremde Träger gemeldet statt still mitgeändert: `SPEC-019`-Zeile,
  Handbuch). §3.2: kein `//nolint` im Diff (`git grep -n nolint` über die geänderten Dateien: 0 Treffer). §3.10: kein
  Workflow im Diff.
- **Commits:** `make doc-commits` und `make commit-traceability` über die Range Exit 0; jeder der 14 Betreffs nennt
  `LH-FA-CFG-007`, `LH-FA-ADM-001` und/oder `ADR-*`, keiner trägt `SPEC-`/`ARC-` im Betreff.
- **Handbuch-Kandidatenlauf** (`git diff --name-only 80eefead -- internal/bootstrap/ tools/schema/
  internal/adapters/driving/`): sechs Dateien in `internal/bootstrap/` (Verdrahtung und Tests) plus
  `tools/schema/nacharbeit-administration.sql`; keine neue Umgebungsvariable, keine neue SQL-Funktion (der Text von sieben
  bestehenden Funktionen ändert den Zeitstempel, Signaturen unverändert), kein neuer Endpunkt — **keine neue
  Betreiber-Oberfläche**; `docs/user/benutzerhandbuch.md` liegt nicht im Diff, also keine Versionshistorie-Pflicht. Die
  Wirkung auf den Betreiber (Aufruf-Reihenfolge in einer Transaktion) trägt der Aufschub mit Adresse
  `slice-transformationen-betriebsdoku` §2 (Diff gelesen: „Aufrufe einer Transaktion werden in Aufrufreihenfolge
  verarbeitet; Anträge auf dieselbe Regel oder Spalte nicht aus überlappenden Transaktionen absetzen — … die Aussage gilt
  für alle sieben Funktionen und gehört in den Transformations-Abschnitt in §4“). Kein Befund.
- **Kenntnis (Prozess):** der Implementer meldete `sed -i` gegen Repo-Dateien (Review F-10, nur die gofmt-Spur, mit
  Fixrunde 1 bereinigt); keine inhaltliche Fehlwirkung gefunden.

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | LOW | **ADR-Aussage breiter als ihre Messung — Mutationsangabe in der Fitness-Function-Zeile von `ADR-0127`.** Die Zeile nennt als rot färbende Mutation „`requested_at`/`clock_timestamp()` aus **einer** Funktion entfernen (die Zeilen der Funktion tragen wieder den Transaktionsbeginn)“, erprobt am SQL-Text mit `now()` in **allen** sieben Funktionen (Spalte „alt“ der Messtabelle), „am Go-Test noch nicht gefahren“. Gemessen an den Tests (S3–S8 und der Implementer-Lauf für `include_column`): rot bei `set_transformation`, `enable_table`, `disable_table`, `exclude_column`, `include_column`; **grün** bei `remove_transformation` (S4) und `backfill_table` (S7). Der Implementer meldet das für `remove_transformation` und führt es als „nur für die später aufgerufene Funktion eines Paares wirksam“ (Plan §3, Mutationstabelle, Zeile „dieselbe Zusage an der Funktion, die vor einer Nachbarin aufgerufen wird“); die Abgrenzung „früher/später im Paar“ trifft nicht (S6, S8, V-3): die Farbe hängt an der Stellung im Test-Ablauf. **Schwere:** LOW, nicht MEDIUM — die Aussage steht in einer Fitness-Function-Zeile, nicht in §Entscheidung, §Konsequenzen oder `SPEC-019`; die Zeile delegiert die Mutation an den Implementer („noch nicht gefahren“), der das Ergebnis meldete; kein Verbraucher liest daraus ein Verhalten, das die Funktionen nicht halten (R0: das Produkt hält die Ordnung für die Gegenrichtung und für `backfill_table`). **Klasse:** dieselbe wie V-1 der Verifikation zu `antragsweg-schema` (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`), im Register bereits als verkörpert geführt — ein weiteres Auftreten nach der Verkörperung (`AGENTS.md` §3.12, Absatz „Verfasser einer ADR“: eine Fitness-Function-Zeile nennt den Test, der sie trägt, erprobt oder als *hergeleitet* gekennzeichnet). **Empfehlung (nicht geändert):** (a) keine Folge-ADR allein dafür (dieselbe Begründung wie am Vorgänger); (b) die Berichtigung bleibt benannte Grenze im Plan §3 (dort vorhanden, im Wortlaut zu präzisieren, V-3) und Closure-Notiz-Eintrag; (c) die nächste ADR zu `requested_at` oder zur Antrags-Queue trägt die Berichtigung als eigene Klausel; (d) der Planner zählt das Auftreten im Register | `docs/plan/adr/0127-…md` §Fitness Function; Plan §3 (Mutationstabelle, drei Zeilen „Fixrunde 2“); `administrationrequest_order_test.go` (Doc-Kommentar von `…CallsKeepCallOrder`) | ja — S3 bis S8 je `make test-store` |
| V-2 | LOW | **Zusage ohne vollständige Eingabeseiten-Bindung: `clock_timestamp()` in `cdc.remove_transformation` und `cdc.backfill_table` ist an keinen Test gebunden.** `make test-store` bleibt grün, wenn die Funktion `now()` schreibt (S4, S7). Die Tests rufen nur die Folgen `remove`→`set`, `exclude`→`include`, `disable`→`enable` auf (Ablauf `remove`, `set`, `exclude`, `include`, `disable`, `enable`); `remove_transformation` steht an erster Stelle, `backfill_table` fehlt. Mein Zusatztest (nur Kopie, R0–R2: Aufrufe `set`, `remove`, `include`, `exclude`, `enable`, `disable`, `backfill`, 60 Transaktionen) ist am Produkt grün (R0) und färbt sich bei beiden Mutationen (R1, R2): in der Folge „`set`, dann `remove` derselben Regel in einer Transaktion“ sortiert ein `remove_transformation` mit Transaktionsbeginn **vor** das `set` — die Ordnung ist die falsche, und die Regel bliebe live und abgeleitet **gesetzt** statt entfernt. Die Zusage von [`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) Festlegung 1 („Aufrufe derselben Transaktion tragen verschiedene Zeitstempel in der Reihenfolge des Aufrufs“, alle sieben Funktionen) ist für zwei Funktionen ungebunden; die Fitness-Function-Zeile der ADR verlangt nur die drei Paare, der Test entspricht ihr. Kein Produktfehler; die Lücke liegt in der Reichweite der Bindung. **Schwere:** LOW (Produkt richtig, benannte Zusage schmaler gebunden als sie lautet). **Empfehlung:** der Implementer erweitert die Aufrufe des Tests um die Gegenrichtung (`set`→`remove`) und einen `backfill_table`-Aufruf (kleine Fixrunde), oder der Planner führt die zwei Funktionen als benannte Grenze mit Adresse | `internal/adapters/driven/postgresstorage/administrationrequest_order_test.go` (`calls`); `tools/schema/nacharbeit-administration.sql` (Zeilen 174, 212) | ja — S4, S7, R1, R2 je `make test-store` |
| V-3 | LOW | **Kommentar-Zusage breiter als ihre Messung.** Der Doc-Kommentar von `TestAdministrationRequestSameTransactionCallsKeepCallOrder` sagt: „In der früher aufgerufenen Funktion bleibt die Reihenfolge richtig und der Test grün“; die Plan-Zeile „Fixrunde 2: dieselbe Zusage an der Funktion, die **vor** einer Nachbarin aufgerufen wird“ sagt „kein Rot“ und nennt die Mutation „im Fall früherer Aufruf äquivalent“. Für `disable_table` (früher Aufruf im Paar `disable`/`enable`) und `exclude_column` (früher Aufruf im Paar `exclude`/`include`) färbt die Mutation den Test rot (S6, S8), weil vor ihnen Aufrufe mit `clock_timestamp()` stehen. Wahr ist die Aussage nur für die **erste** Funktion des Ablaufs, `remove_transformation` (S4) — und dort ist sie kein Beleg für Äquivalenz der Mutation, sondern für die Enge der Test-Folge (V-2). Der Kommentar von `TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived` („Dieselbe Ersetzung in `cdc.remove_transformation` lässt die Reihenfolge richtig und den Test grün“) ist für diesen Test wahr (S4 grün). **Schwere:** LOW (Kommentar und Plan-Zeile, `AGENTS.md` §3.7 Klasse „Zusage“ und §3.12 Instanz B, Beleg-Anker vorhanden, aber die Verallgemeinerung trägt er nicht) | `administrationrequest_order_test.go` (Kommentar vor `TestAdministrationRequestSameTransactionCallsKeepCallOrder`); Plan §3 (Mutationstabelle, Zeile „dieselbe Zusage an der Funktion, die vor einer Nachbarin aufgerufen wird“) | ja — S6, S8 |
| V-4 | LOW | **`SPEC-019` Zeile 675 ist ungenau.** „`rule_name` ist für `set_transformation` und `remove_transformation` Pflicht, `rule_spec` nur für `set_transformation` (Domänen-Invarianten des Antrags-Konstruktors)“ nennt den Konstruktor als Ort der Prüfung. Seit diesem Slice nimmt `NewAdministrationRequest` beide Arten mit leerem Regelnamen und leerer Regelform an; die Prüfung liegt im Use Case und endet als `failed` mit `Regelname ist ungültig` bzw. `rule_spec ist ungültig` (Tabelle derselben Spec). Der Satz über die Spalten-Antragsarten (Zeile 640, „Domänen-Invariante des Antrags-Konstruktors“) bleibt richtig (Konstruktor lehnt leere `column` weiter beim Lesen ab, Test gebunden). **Bewertung:** die Verhaltens-Aussage der Spec (Pflicht, Fehlertext) bleibt wahr; falsch ist nur der Klammer-Zusatz, der den Ort nennt; niemand liest daraus ein falsches Verhalten, weil die Fehlertext-Tabelle darunter das Verhalten trägt. **Meldung an den Planner** (Frist laut Plan §3: Closure dieses Slice; Spec ist fremde Datei): die Klammer streichen oder ersetzen durch „die Prüfung endet als `failed` mit dem Fehlertext der Tabelle unten, nicht als Ablehnung beim Lesen der Queue“; zugleich den Register-Zustand von `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (teilweise gelöst: Regelfelder gelöst, Rest F-6) und dessen Verweis auf den Pfad `open/…` des Slice nachziehen | `spec/pflichtenheft.md` Zeile 675; Plan §3 (Zeile „`SPEC-019` … gemeldet, nicht geändert“) | ja — `git grep -n 'Domänen-Invarianten des Antrags-Konstruktors' -- spec` |
| V-5 | INFO | Der Doc-Kommentar von `activatedTableBindings` sagt „Eine Regelform, die die Ableitung nicht mehr in eine Regel führt, endet den Start wie ein Lesefehler“; der Lesefehler des Bestands endet in der Startfehlerklasse `storage`, die Faltung endet als `internal` (Review F-8, INFO; Plan §6 nennt die Klasse `internal` korrekt). Erreichbar erst, wenn ein späterer Parser eine bereits vermerkte Regelform ablehnt. Kein Nachzug erfolgt, weil INFO | `internal/bootstrap/wiring.go` (Doc-Kommentar `activatedTableBindings`); Plan §6 | ja — Fehlerklasse am Faltungsfehler prüfen |
| V-6 | INFO | Fixrunde 1, Fixrunde 2 und `ADR-0127` haben keinen zweiten Reviewer-Lauf; ihre Wirkung ist in §4 (S1–S8, G1–G7) und §5 mutiert bestätigt. Fixrunde 2 ändert produktiv nur den Funktionstext (Zeitstempel) und Kommentare. Ein Reviewer-Blick bleibt Rollen-Entscheidung des Planners; der Review-Report ist seit `c7796d83` unverändert und nennt die Fixrunden nicht | `git log --format=%h -- docs/reviews/review-slice-transformationen-antragsweg-usecase.md` (ein Commit) | ja |
| V-7 | INFO | **Übernommen, nicht nachgemessen:** die 26 Mutationen des Reviews (M1–M26) und die Mutationen des Implementers, soweit nicht in §4 wiederholt; die PostgreSQL-17.11-Hälfte und die Zufallsstichproben der Messtabelle von `ADR-0127`; die Grenzen der ADR (überlappende Transaktionen mehrerer Sitzungen, Rückwärtssprung der Uhr) — hergeleitet, nicht erprobt (so in der ADR benannt); die Image-Neutralität und die Läufe `make test-replication`/`make test-integration`/`make image`. Meine Messungen laufen an der PostgreSQL-Version von `make test-store` | §4, §7 | nein |

Kein HIGH, kein MEDIUM. Keine DoD-Verletzung.

## 9. Verdikt

**DoD bestätigt:** ja — jede der acht `[x]`-Zeilen (Nr. 1–8) ist am Ist-Zustand belegt; die vier `[ ]`-Zeilen (Nr. 9–12)
sind korrekt offen (Planner-Closure). **Plan-vs-Code:** keine unbenannte Abweichung; die Meldung der `SPEC-019`-Zeile und die
benannten Grenzen (Restklasse leeres Schema/Tabelle, Vermerk-Fehler mit kollidierendem Folgeantrag, quellweite
Ableitung) stehen im Plan §3/§6; die Meldung ist von mir bewertet (V-4). **Entscheidungs-Konformität:**
[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Teilfrage 1, 3, 6, Folgepflicht 3),
[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md),
[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 5,
[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
[`ADR-0028`](../plan/adr/0028-inbound-use-cases.md),
[`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md),
[`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md),
[`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md) und
[`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) (Festlegung 1 bis 3, Folgepflichten 1 bis 4)
konform; die Mutationsangabe der Fitness-Function-Zeile von `ADR-0127` ist breiter als ihre Messung (V-1, LOW), und die
Tests binden die Zusage für `remove_transformation` und `backfill_table` nicht (V-2, LOW).
[`SPEC-019`](../../spec/pflichtenheft.md)/[`SPEC-030`](../../spec/pflichtenheft.md): Fehlertexte, Prüfreihenfolge, „Ordnung der
Verarbeitung“ und K1-Klausel tragen; eine Klammer ist ungenau (V-4). **F-1 des Reviews (Ordnung):** mit `ADR-0127` und
Fixrunde 2 gelöst — Live- und abgeleiteter Stand sind für die drei getesteten Paare und die Gegenrichtung (R0) gleich;
verbleibende Grenzen benannt. **Gates:** `make test`, `make test-store`, `make a-check`, `make coverage-gate` (84.80 %),
`make gates`, `make suchlauf-nachmessen` (32 Zeilen), `make doc-commits`, `make doc-immutable`,
`make commit-traceability` und das Guard-Skript (sechs Läufe, Alt-Tag `v0.2.0`, Exit 0/0/0) im eigenen Lauf grün; sieben
Go- und acht Store-Eingabeseiten-Mutationen: alle rot bis auf `remove_transformation` und `backfill_table` (V-2).

**Übergabe:** an den Planner — Closure-Notiz mit Lerneintrag (Klassen aus Review und diesem Report:
„Verarbeitungs-Ordnung ungleich Ableitungs-Ordnung“, „ADR-Aussage breiter als ihre Messung“ als weiteres Auftreten nach der
Verkörperung, „Kommentar-Zusage breiter als ihre Messung“), Beobachtungs-Register, Ausgänge der §6-Risiken (darunter die
Ordnungs-Zeile mit den Grenzen der ADR, die Restklasse leeres Schema/Tabelle mit Adresse, die quellweite Ableitung mit
Adresse `slice-transformationen-backfill-pfad`), Paarungen bei der Closure der Welle; Meldung V-4 (`SPEC-019`-Klammer) mit
Frist Closure; V-1 als benannte Grenze und Klausel der nächsten ADR (keine Folge-ADR allein dafür); V-2 und V-3 an den
Implementer (kleine Fixrunde: Test um `set`→`remove` und `backfill_table` erweitern, Kommentar und Plan-Zeile auf die
Stellung im Ablauf präzisieren) **oder** vom Planner als benannte Grenze mit Adresse geführt; V-5 bis V-7 als Anmerkungen.
Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
