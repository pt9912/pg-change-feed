# Verifikations-Report: slice-backfill-change-origin — 2026-09-24

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
(`AGENTS.md` §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-backfill-change-origin.md`](review-slice-backfill-change-origin.md);
Hausform dieses Reports:
[`verifikation-slice-backfill-row-image-gemeinsam.md`](verifikation-slice-backfill-row-image-gemeinsam.md).

**Gegenstand:** Slice-Plan `slice-backfill-change-origin` (Welle
`welle-backfill-bestand`), Diff-Range `09386619..HEAD` (`daa86db5`), 20 Commits:
Lifecycle und Verantwortlich (`3b7f828a`, `f8ae1fc6`, `f4ba82ab`), Umsetzung
(`98dab17f` Domäne, `c40e2ead` Store/Schema, `e95937cf` HTTP/Handbuch, `2608df37`,
`5f126971`, `988e8a40`), Fixrunde 1 zu [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)
(`e01d4ba6`, `14f971b7`, `fa954778`, `4744c4a5`, `c1df2ac6`), Review-Report `ae6496a2`,
Fixrunde 2 (`e3d55e2b`, `daa86db5`). Kontext, nicht Prüfling: `ADR-0114` (Accepted),
Architect-Verdikt (`4ba4ef27`, `3acd2c8d`), Planner-Träger `ae97247a`. Slice-Parent für
die Code-Diffs ist `f4ba82ab`. Dieser Lauf ändert keine Spec, keinen Plan und keinen
Code; er schreibt nur diesen Report.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert
gesichert (`make … > <log> 2>&1; echo …=$?`), die Logs danach gelesen.

| Sensor | Ausgang | Beleg aus meinem Lauf |
|---|---|---|
| `make gates` (vor dem Report-Commit) | **EXIT=0** | baseline-verify `v6.9.0 OK — 54 Dateien` · `d-check: 979 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `coverage-gate: OK — Coverage 82.90% erfüllt Schwelle 80%` · generated-sync OK (byte-gleich, Stufe `proto-export`) · a-check `gesamt: 0 Befund(e)` |
| `make test` (Race-Detector) | **EXIT=0** | kein `FAIL`; `ok` für `domain/model`, `postgresstorage/mapper`, `postgresstorage/sqlexec`, `driving/http`, `tools/schema/rolloutguard` |
| `make a-check` | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **EXIT=0** | `total: (statements) 82.9%`, `coverage-gate: OK — Coverage 82.90% erfüllt Schwelle 80%` |
| `make test-store` | **EXIT=0** | `DB-Adapter-Coverage: 77.13% (gedeckt 533 von 691 Statements; Profile gemergt: store,replication)`, `db-coverage: OK — … erfuellt Schwelle 70%` (gedruckte Zeile dieses Laufs) |
| `make test-replication` (nicht vom Plan verlangt, zur Kenntnis) | **EXIT=0** | dieselbe Zeile: `DB-Adapter-Coverage: 77.13% (gedeckt 533 von 691 …)` |
| `make image` + `make test-integration` (nicht in der DoD genannt; der Plan behauptet den Lauf in §3) | **EXIT=0 / EXIT=0** | `E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand`, `Lauf abgeschlossen — E2E-Abdeckungstabelle aus 13 Go-Zeilen und 28 Bash-Zeilen` |
| `bash tools/harness/run-schema-rollout-guard-test.sh` (ungekürzt) | **EXIT=0** | sechs Läufe; Lauf 5 druckt „Lauf 5 OK — Tag v0.1.2: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar"; Schlusszeile „Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit d-migrate-Exit 8". Danach: kein Container, kein Netz, kein Temp-Verzeichnis, `git status --short` leer (die Wiederherstellung von `tools/schema/plan.yaml`/`down.sql` greift) |
| `make commit-traceability RANGE=09386619..HEAD` | **EXIT=0** | `OK — 20 Commit(s) in "09386619..HEAD", Betreffs ohne Struktur-ID`; `git log --format=%s 09386619..HEAD \| grep -E 'SPEC-\|ARC-'` → 0 Treffer |
| `make doc-commits RANGE=09386619..HEAD` | **EXIT=0** | 979 Dateien, 0 Befunde |
| `make doc-immutable RANGE=09386619..HEAD` | **EXIT=0** | 979 Dateien, 0 Befunde; im Diff liegen `docs/plan/adr/0114-…` (neu angelegt und im Range angenommen) und `docs/plan/adr/README.md` (Indexzeile), keine bestehende `Accepted`-ADR wird verändert |
| `make doc-trace` (advisory) | **EXIT=0** | 80 Anforderung(en), 3 Waise(n): `LH-FA-CAP-009`, `LH-FA-CFG-007`, `LH-FA-CFG-008` — Backfill-Waise bis zum Liefer-Slice erwartet, nur zur Kenntnis |

`make doc-immutable` ohne `RANGE` bricht mit „flag needs an argument: --range" ab
(Exit 2); der Aufruf braucht `RANGE=`. Kein Befund am Slice.

Nebenbeobachtung (V-6): `make test-store` und `make test-replication` überschreiben
`tools/schema/plan.yaml` (nur die `target`-Zeile, Test-Datenbank statt der committeten
Compose-Bezeichnung); ich habe die Datei nach jedem Lauf per `git checkout`
zurückgenommen, `git status --short` war danach leer.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Domäne: `model.ChangeOrigin` geschlossene Menge, Default `wal`, anderer Wert abgelehnt, `Change.Origin`, Doc-Kommentar | **erfüllt** | `internal/domain/model/change.go`: `ChangeOrigin` (`wal`/`backfill`), `NewChangeOrigin` (`""` → `wal`, sonst `ErrInvalidChangeOrigin`), `OrDefault`, `NewChange` setzt `ChangeOriginWAL`, `WithOrigin`; Doc-Kommentare von `Operation` und `Change` nennen das Feld und die Live-Wege ohne es. Tests `TestNewChangeDefaultsOriginToWAL`, `TestNewChangeOriginClosedSet` (`""`, `wal`, `backfill`, `snapshot`, `WAL`, ` backfill`), `TestChangeWithOrigin`, `TestChangeOriginOrDefault`; `make test` EXIT=0; Mutation M4 (§5) rot |
| 2 | Store/Schema: Spalte, View, explizite Spaltenlisten, Mapper, `NULL` ≙ `wal`; `make schema-rollout` zweimal Exit 0; `plan.yaml`/`down.sql` neu und committet; Vertragstest unverändert | **erfüllt** (beide benannten Grenzen bestätigt) | `schema.yaml`: `origin: { type: text }` (nullable, ohne Default, ohne CHECK), View mit `COALESCE(c.origin, 'wal') AS origin` als **letzter** Spalte samt `columns:`-Eintrag; `schema.sql` (zweiter Träger) mit derselben Spaltenform; `queries.go`: `InsertChange` (neun Spalten/`$9`) und `SelectChanges` (14. Spalte `COALESCE(c.origin, 'wal')`); `translate.go` scannt `&row.Origin` als 14. Spalte; `mapper.go` (`NewChangeRows`/`ToChange` über `NewChangeOrigin`); `store.go` reicht `row.Origin` durch. `make test-store` EXIT=0 mit `TestPersistAndReadCarryChangeOrigin` (Zeile per SQL auf `NULL`, liest `wal`), `TestPersistRejectsUnknownChangeOrigin`, `TestChangesViewCarriesOriginLikeReadChanges`; Mutation M5 (§5) rot. Guard-Test Lauf 1 und 2: Rollout gegen frisches Ziel und Folgelauf, beide Exit 0. `tools/schema/plan.yaml` trägt `COALESCE(c.origin` (1 Treffer), `operationsTotal":14`, Compose-Ziel-Bezeichnung (`5f126971`); `git diff --stat f4ba82ab..HEAD -- test/integration tools/harness/run-integration-tests.sh docs/user/e2e-abdeckung.md` → leer, `make test-integration` EXIT=0 mit „E2E-Abdeckungstabelle unverändert". Grenze (1) (Parent-Ziel blockiert) aufgelöst durch `ADR-0114`, belegt in §3 dieses Reports; Grenze (2) (Parität trägt der Store-Tier-Test) bestätigt |
| 3 | `GET /changes` trägt `origin`, letzter Feldname, fehlender Wert `wal`; Handbuch zwei Stellen, `Version:`, Historienzeile | **erfüllt** | `readchanges.go`: `Origin string \`json:"origin"\`` als letztes Feld der 13, `string(change.Change.Origin.OrDefault())`; Test `TestReadChangesTraegtOriginAlsLetztesFeld`, Mutation M2 (§5) rot. Handbuch: SQL-Beispiel unter „Änderungen lesen" (Spaltenliste + Bedeutung), Feldliste der `GET /changes`-Antwort; `Version: 1.46`, `Stand: 2026-09-24`, Historienzeilen 1.44/1.45/1.46 im selben Diff (Schritt 17 aus `implement-slice.md`; Kandidatenlauf `git diff … \| grep -E '^\+Version:\|^\+\| [0-9]+\.[0-9]+ \|'` trifft beide Muster) |
| 4 | Upgrade über die View-Signaturänderung (`ADR-0114`, `LH-QA-OPS-005`): Guard, Makefile-Vorlauf, Unit-Tests, Guard-Test mit sechs Läufen, Handbuch, `harness/README.md` | **erfüllt** | siehe §3 (Sicherheitsnaht) und §4; Guard-Test EXIT=0 und Mutation G5 (§5) rot an Lauf 6a; Upgrade-Belege Parent→Arbeitsbaum und Alt-Tag selbst gefahren (§3.2); Handbuch §4 „Schema aktualisieren" (Version 1.45/1.46); `harness/README.md` Zeile `make schema-rollout` (sechs Läufe, Vorlauf, zwei Negativfälle) |
| 5 | `make gates` grün, Exit-Code ungefiltert und gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 6 | Review durchgeführt, Report unter `docs/reviews/` | **erfüllt** (Kästchen wurde im Fixrunden-Commit `daa86db5` gesetzt, Implementer-Schritt 21) | `review-slice-backfill-change-origin.md` liegt vor (0 HIGH / 2 MEDIUM / 3 LOW / 5 INFO). Ein Reviewer-Nachprüfungs-Report existiert nicht und ist für die Zeile nicht tragend: die Zeile belegt, dass der Review-Schritt stattgefunden hat (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug: bei nötiger Fixrunde zieht der Implementer das Kästchen). Damit die Zeile nicht auf einer Behauptung ruht, habe ich die Behebung von F-1…F-8 selbst am Code nachgemessen: F-1 → `guard.go` Zweig `destructiveConfirmationReason` lehnt leere `operationIds` ab, Test `TestDecideRefusesDestructiveBlockerWithoutOperations` (allein / neben Fremdobjekten / neben View-Signatur), Mutation M3 rot; F-2 → Lauf 6a (nicht deklarierte Funktion, Assertion „Error 8", Fortbestehen der Funktion) bindet die Bekannt-Liste, Mutation G5 rot (§5); F-3 → Makefile-Kommentar Absatz 1 nennt „bekanntes Fremdobjekt oder View-Signatur-Änderung … mindestens ein bekanntes Fremdobjekt → `--allow-destructive`", deckt sich mit `decide`; F-4 → Suchraum erweitert (siehe V-1/V-2), Meldung als Zeile in `slice-backfill-sdk-origin`; F-5 → Handbuch „Eigene Rechte" und Makefile-Kommentar nennen den Rechteverlust; F-6 → Skript sichert/restauriert `plan.yaml`/`down.sql` (Arbeitsbaum nach meinem Lauf leer); F-7 → unverändert gemeldet (`run-integration-tests.sh:2711` bleibt); F-8 → Makefile-Kommentar und Handbuch nennen die feste Adressierung von `cdc`. Kein offenes HIGH/MEDIUM |
| 7 | §3.13-Suchlauf: Feld in §3 trägt Gefundenes **und** Nichtgefundenes, beide Stände | **erfüllt, mit V-1 (LOW) und V-2 (INFO)** | Die Suchlauf-Zeilen (Erstlauf und beide Fixrunden-Felder) sind vollständig gegen beide Stände nachgemessen (§3.3): zwei Zahlen stimmen nicht. Substanz und Behandlung der Zeilen bleiben tragfähig |
| 8 | Doku-Update (Handbuch, Änderungshistorie) | **erfüllt** | siehe Zeile 3 und 4 |
| 9 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*" |
| 10 | Reconciliation-Register — „entfällt" | **korrekt offen / entfällt** (V-7 INFO) | Greenfield, keine Reconciliation-Datei; die Zeile trägt „entfällt", steht aber auf `[ ]` (im Vorgänger-Slice `slice-backfill-row-image-gemeinsam` auf `[x]` gesetzt); beim Closure-Nachzug einheitlich setzen |
| 11 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | vier Risiko-Zeilen tragen „**Ausgang:** *(bei Closure)*" |
| 13 | Drei Paarungen | **korrekt offen** | hängen an der Closure von `welle-backfill-bestand` (offen) |

Kein `[x]` ohne Beleg, kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre.

## 3. Kernaussagen, selbst gemessen

### 3.1 `origin` von der Domäne bis zu den Lesewegen

- **Geschlossene Menge, Default, Lesen:** `NewChangeOrigin` lehnt `snapshot`, `WAL`,
  ` backfill` ab und liest `""` als `wal`; `NewChangeRows` und `ToChange` laufen beide
  durch die Menge (unbekannter Wert endet als Domänenfehler vor der Datenbank bzw. beim
  Lesen). Die View und `SelectChanges` tragen denselben `COALESCE(c.origin, 'wal')`,
  `NULL` liest über beide als `wal` (`TestPersistAndReadCarryChangeOrigin`).
- **Live-Wege unverändert:** `git diff --stat f4ba82ab..HEAD -- internal/adapters/driving/grpc internal/adapters/driven/grpcstream gen proto`
  leer; in `natsstream/publisher.go(+_test)` und `driving/http/sse.go` ändern sich nur
  Kommentare („zehn Felder … ohne `Origin`"); `generated-sync` OK.
- **Schema-Träger:** `schema.yaml` und `schema.sql` tragen `origin` in derselben Form
  (nullable, ohne Default, ohne CHECK); ohne die Spalte in `schema.sql` scheitert
  `InsertChange` in den Store-Tests (Plan-Zeile, vom Reviewer bestätigt).
- **Spec:** `spec/pflichtenheft.md` trägt `origin` in `SPEC-002` (Z. 338/345/840) und
  `SPEC-022` (Z. 617, „ein gespeicherter Change ohne das Feld liest als …"); die
  Live-Wege-Aussage steht in `SPEC-021` (Z. 585, „ohne das Feld `origin`"). Der Diff
  setzt das um und ändert die Spec nicht (`git diff --stat f4ba82ab..HEAD -- spec` leer).

### 3.2 Sicherheitsnaht `rolloutguard` und Upgrade-Belege

**Entscheidungslogik (`guard.go`, gelesen):** jeder Blocker muss entweder (a) Grund
`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION` mit nichtleeren `operationIds`, alle auf
`knownForeignObjects`, oder (b) Grund `MANUAL_ACTION_REQUIRED` mit nichtleeren
`operationIds` sein, jede eine `ReplaceView`/`VIEW`-Operation mit
`VIEW_SIGNATURE_INCOMPATIBLE`-Diagnose zur selben `operationId` und einem einzelnen
Pfadsegment aus `^[a-z_][a-z0-9_]*$`; jeder andere Fall liefert eine **leere**
Entscheidung. `allowDestructive` wird nur im Zweig (a) gesetzt; die Klasse (b) allein
setzt es nicht (`TestDecideViewSignatureAlone`). `main.go` gibt Maschinenzeilen
(`allow-destructive`, `drop-view <name>`) auf stdout aus, Exit 0 nur wenn irgendetwas
erlaubt ist. Unit-Tests: 19 `Test…`-Funktionen in `guard_test.go` (gezählt), `make test`
EXIT=0.

**Makefile-Vorlauf (Diff gelesen):** der Vorlauf läuft nur bei Precheck-Exit 8 und
nur über die validierten Namen aus den Maschinenzeilen; `DROP VIEW cdc.$$v` ohne
`CASCADE`, `psql -v ON_ERROR_STOP=1`, `|| exit 1`; ein Guard-Fehlschlag (Exit 1 oder 2)
lässt beides leer (`guard_out=$(…) && guard_ok=1 || guard_ok=0`, keine Pipe um den
`make`-Exit); `--execute` läuft in jedem Fall; die Ausgabe des Vorlaufs trägt je View
die Meldung. Kommentarblock nennt Nicht-Atomarität, Rechteverlust und feste Adressierung
von `cdc`.

**Guard-Test (Lauf für Lauf, aus meinem Lauf):**

| Lauf | Beleg (gedruckt/gemessen) |
|---|---|
| 1/6, 2/6, 3/6 | Rollout frisch Exit 0; zweiter Lauf trägt „`--allow-destructive`" und **keinen** „Vorlauf" (Assertion greift im Skript); Lauf 3 (echte anstehende `DROP COLUMN error_message` neben den sechs Fremdobjekten) kommt zurück, ohne Vorlauf |
| 4/6 (Signaturänderung) | View auf `SELECT change_id, transaction_id` gesetzt; Rollout Exit 0 mit Meldung `schema-rollout: Vorlauf (ADR-0114) - View-Signatur-Aenderung, DROP VIEW cdc.changes`; Soll-Signatur, `has_table_privilege('cdc_reader', …)` = `t`, Zeile `lauf4-ch|wal`; Folgelauf Exit 0 **ohne** Vorlauf |
| 4 (abhängiges Objekt) | View in `public` hängt an `cdc.changes`; Rollout scheitert laut (Exit ≠ 0), Vorlauf-Meldung steht in der Ausgabe, das Objekt besteht (kein `CASCADE`); Wiederholungslauf ohne das Objekt heilt |
| 5/6 (Alt-Tag) | siehe §1, Tag `v0.1.2`, drei Exit 0, Zeile lesbar, Soll-Signatur |
| 6a | nicht deklarierte Funktion; Guard-Ausgabe „unbekannter Blocker DropFunction:FUNCTION:… (DropFunction FUNCTION zz_rolloutguard_unbekannt()) — --execute laeuft ohne --allow-destructive und ohne Vorlauf"; make-Exit 2 mit d-migrate-„Error 8"; Funktion besteht danach |
| 6b | nicht deklarierte Spalte; Guard „unbekannter Blocker … DropColumn COLUMN source._rolloutguard_test_col"; Abbruch Exit 8 |

**Upgrade-Belege unabhängig gefahren (Wegwerf-PostgreSQL 18, `PG_TEST_IMAGE`, Skript im
Scratchpad, `git archive` in einem Verzeichnis unter dem Scratchpad, nie im Repo):**
Parent `f4ba82ab`, ausgerollt mit dessen eigenem Makefile: `PARENT_ROLLOUT_EXIT=0`;
Spalten der View davor: 13, letzte `committed_at`; eine Datenzeile ohne `origin`
geschrieben. Arbeitsbaum, erster Lauf: `WORK_ROLLOUT_1_EXIT=0`, Vorlauf-Meldung 1 Treffer;
zweiter Lauf: `WORK_ROLLOUT_2_EXIT=0`, Vorlauf-Meldung 0 Treffer. Danach
`p-ch|wal` über `cdc.changes`, `has_table_privilege('cdc_reader','cdc.changes','SELECT')`
= `t`, Spalte `cdc.change.origin` vorhanden. Das deckt sich mit den Plan-Aussagen im
„Ausgang der Messung"; die beiden Läufe gegen `f4ba82ab` und `v0.1.2` bestätigen das
Upgrade unabhängig vom Guard-Test.

**Urteil zur Lockerung (`AGENTS.md` §3.6):** die Lockerung geht nicht weiter als
[`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md). `--allow-destructive`
bleibt an bekannte Fremdobjekte gebunden (unverändert zum Parent-Stand, jetzt zusätzlich
gegen den Vakuum-Fall gehärtet); neu ist allein der Vorlauf für die eine Klasse, jeweils
alles-oder-nichts (Entscheidung 2), ohne `CASCADE` (Entscheidung 3), ohne Anlass kein
Fenster (Entscheidung 4/5). Die Härtung (leere `operationIds` → Ablehnung) verschärft.

### 3.3 §3.13-Suchlauf und §3.12-Zahlen — beide Stände selbst gemessen

Erstlauf (Stände `f4ba82ab` und `e95937cf`):

| Plan-Zahl | Meine Messung | Ergebnis |
|---|---|---|
| `committed_at` in `docs/user spec harness`: 6 je Stand, davon 3 im Handbuch | 6 / 6; Handbuch 3 / 3 | bestätigt |
| `origin` (`-w`) in `docs/user`: Parent 1, Diff-Stand 6; `spec`: je 9 | 1 / 6; 9 / 9 | bestätigt |
| `Felder` 33, `zehn Felder` 19, `elf/zwölf/dreizehn Felder` 0 je Stand | 33 / 19 / 0, beide Stände | bestätigt |
| `SELECT *` (`tools test internal`): je Stand 3 | 3 / 3 | bestätigt |
| „dieselben (zehn )?Felder wie" / „Feldern wie der Domain-Typ": je Stand 8 | 8 / 8 | bestätigt |
| `git ls-tree` Schema-Dateien: je Stand 14 | 14 / 14 | bestätigt |
| strikte Dekoder-Muster: je Stand 0 | 0 / 0 | bestätigt |
| E2E-Tabelle/Runner/Testpaket im Diff: leer | 0 Zeilen | bestätigt |
| `examples/`: 0 Treffer für „(ten\|zehn) (fields\|felder)" | 0 / 0 (zeilenweise Suche) | bestätigt; ein zeilenumbrechendes „zehn / Feldern" in `examples/csharp/sse-client/SseStream.cs` (Z. 6) beschreibt einen Live-Weg (SSE) und bleibt wahr |
| **`(ten\|zehn) (fields\|felder)` in `sdks examples`: je Stand 18** | **16 an allen vier gemessenen Ständen** (`09386619`, `f4ba82ab`, `e95937cf`, `HEAD`); gemessen mit `git grep -n -i -E '(ten\|zehn) (fields\|felder)' <Stand> -- sdks examples \| wc -l` (das Muster in der Plan-Zelle steht mit maskiertem `\|` wegen der Tabelle, als reguläres ERE lieferte es 0) | **weicht ab — V-1** |

Fixrunden-Felder:

| Plan-Zahl | Meine Messung | Ergebnis |
|---|---|---|
| Guard-Beschreibungen: Parent `ae97247a` 33 Treffer in 8 Dateien, `fa954778` 42 in 9 | 33 / 8 und 42 / 9 | bestätigt |
| Läufe des Guard-Tests: 9 → 11 | 9 / 11 | bestätigt; die Treffer in `open/slice-backfill-sql-administration.md` (Z. 193) und `slice-transformationen-antragsweg-schema.md` (Z. 171) sind Such-Aufträge in fremden Plänen, gemeldet, nicht geändert |
| „nicht idempotent": 2 → 0 | 2 / 0 | bestätigt |
| Erweiterte Suche über den ganzen Baum: „je Stand 24 (unverändert)" | Parent `09386619`: 24; `HEAD`: **25** — der Zusatztreffer ist die neue Zeile in `open/slice-backfill-sdk-origin.md` (Fixrunde-2-Commit `daa86db5`, sie enthält „twelve"/„zwölf"), die der Plan selbst als Meldung anlegt | Zahl zum Zeitpunkt der Niederschrift plausibel, gegen `HEAD` um eins verschoben — **V-2** |
| vier Feldanzahl-Träger für `GET /changes` in `sdks/` | `Sse/Models/Change.cs` Z. 17, `Nats/Models/Change.cs` Z. 21 (je „eleven fields"), `sse/model/Change.kt` Z. 19–20 („twelve" am Zeilenende), `models.py` Z. 264 („zwölf …") | bestätigt |
| HTTP-Modell trägt am Parent zwölf Felder (`Changes.cs`/`Changes.kt`) | je 12 Feld-Attribute `commit_position` … `committed_at` in `Changes.cs` Z. 15–26 und `Changes.kt` Z. 23–34 | bestätigt |

**Fremde Träger:** die sechs SDK-Fundstellen (vier Zahlen-Träger, zwei „same ten
fields"-Kommentare in `Http/Models/Changes.cs` Z. 8 und `http/model/Changes.kt` Z. 8)
stehen als eine §3-Zeile in `slice-backfill-sdk-origin`; die SDK-Dateien selbst sind im
Diff unberührt. `.proto` Z. 13 und `.pb.go` Z. 31 („dieselben Felder wie der Domain-Typ")
sind unverändert und im Plan gemeldet (Änderung verlangt `make proto-generate`). Die
Zeile `tools/harness/run-integration-tests.sh` Z. 2711 („Exit 8 auf vier Fremdobjekten")
ist unverändert und gemeldet — sie ist Träger der Zeilenanker von
`docs/user/e2e-abdeckung.md` (`make test-integration` lief mit „E2E-Abdeckungstabelle
unverändert"). Alle Meldungen stehen im Plan.

## 4. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2/8 | Feld `origin`, geschlossene Menge, Reichweite in den Lesewegen (`cdc.changes`, `GET /changes`), nicht in den Live-Wegen | §3.1; keine Live-Wege-Änderung; `SPEC-002`/`SPEC-022` gelesen, unverändert | konform |
| [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) | Pflicht-Report, destruktive Operationen default blockiert, `knownForeignObjects` | `--allow-destructive` weiterhin nur bei bekannten Fremdobjekten; `plan.yaml`/`down.sql` neu erzeugt und committet; Guard-Test Lauf 1–3 und 6 unverändert in Nummern und Ausgang | konform |
| [`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md) Trigger 1 | Schema-Upgrade über den Alt-Tag prüfen | Lauf 5 des Guard-Tests (Tag `v0.1.2`) und mein Parent-Lauf | konform (für das Schema; der Container-Tausch prüft es weiterhin nicht, `ADR-0064` Negativ unverändert) |
| [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 1 | Klasse über Operation `ReplaceView`/`VIEW` und Diagnose `VIEW_SIGNATURE_INCOMPATIBLE` | Die Umsetzung verlangt **zusätzlich** den Blocker-Grund `MANUAL_ACTION_REQUIRED` | **konform, zulässig eng** (V-4): die Aussage geht in die sichere Richtung — meldet ein künftiger d-migrate dieselbe Diagnose unter anderem Grund, fällt der Fall in „unbekannt" und endet mit Exit 8 statt mit einem Vorlauf; Entscheidung 2 („jeder andere Blocker lässt beides ausfallen") trägt das. Die Präzisierung steht im Plan (§3, Fixrunde-Zeile `rolloutguard`) und im Kommentar von `guard.go`; real gemessen: der Report der Klasse trägt `MANUAL_ACTION_REQUIRED` (Fixture `realViewSignatureReport`, Guard-Test Lauf 4 fährt den Vorlauf real) |
| ADR-0114 Entscheidung 2 | alles oder nichts | `TestDecideViewSignatureWithUnknownBlockerRefusesEverything`, `…WithKnownForeignObjects`; Vakuum-Fall (leere `operationIds`) abgelehnt | konform |
| ADR-0114 Entscheidung 3 | `DROP VIEW cdc.<name>` ohne `CASCADE`, je View gemeldet, Rechte über `nacharbeit-roles.sql` | Makefile; Guard-Test Lauf 4 „abhängiges Objekt"; `cdc_reader` `t` nach dem Rollout | konform, mit V-3 (ADR-Lücke zur ACL) |
| ADR-0114 Entscheidung 4/5 | kein Vorlauf ohne Anlass, additive Änderungen ohne Vorlauf | Guard-Test: Lauf 2, 3 und der Folgelauf von Lauf 4 melden keinen „Vorlauf" (Assertionen); mein zweiter Parent-Arbeitsbaum-Lauf: 0 Treffer | konform |
| ADR-0114 Entscheidung 6 | Reihenfolge Schema-Rollout vor Container-Tausch im Handbuch | Handbuch §4 „Reihenfolge beim Upgrade" | konform |
| ADR-0114 Entscheidung 7 | Guard-Test-Lauf für abweichende Signatur, Alt-Tag-Lauf, Lauf mit unbekanntem Blocker bleibt der letzte | Lauf 4, 5, 6; Alt-Tag mit gedrucktem Tag und Exit-Codes | konform |
| ADR-0114 Folgepflicht | Umsetzung in `rolloutguard`, Makefile, Guard-Test, Handbuch, `harness/README.md` | alle fünf Träger im Diff (`examples/bootstrap.sh` zusätzlich als Träger nach Suchlauf) | konform |
| [`ADR-0058`](../plan/adr/0058-testansatz-fuenf-luecken.md) Entscheidung 2 | additive Spalte an `cdc.change` als belegter Weg | Spalte nullable, ohne CHECK; Konvergenz gemessen | konform |
| [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md) | Changes lesen über die HTTP-API | Feldreihenfolge unverändert, `origin` zuletzt | konform |

## 5. Mutationen (Eingabeseite, selbst gesehen)

Je Mutation: Edit per `sed`, `make test` (Race-Detector) bzw. der genannte Sensor,
Exit-Code und `FAIL`-Zeilen gelesen, danach `git checkout -- <Datei>`; `git diff --stat`
danach leer.

| # | Zusage | Mutation | Sensor | Rote Tests |
|---|---|---|---|---|
| M1 | Store schreibt die gesetzte Herkunft | `NewChangeRows`: `Origin` immer `"wal"` (zuerst als `Origin: "wal"` versucht → Build-Fehler durch die unbenutzte Variable, deshalb `"wal" + string(origin)[:0]`) | `make test` **EXIT=2** | `TestChangeRowRoundTripsOrigin` |
| M2 | fehlender Wert liest in der HTTP-Antwort als `wal` | `OrDefault()` aus `readchanges.go` entfernt | `make test` **EXIT=2** | `TestReadChangesTraegtOriginAlsLetztesFeld` |
| M3 | Blocker ohne Operationen belegt keine bekannte Operation (F-1) | Prüfung `len(b.OperationIDs) == 0` im Zweig `destructiveConfirmationReason` (drei Zeilen) entfernt | `make test` **EXIT=2** | `TestDecideRefusesDestructiveBlockerWithoutOperations` |
| M4 | `""` liest als `wal` | Fall `""` aus `NewChangeOrigin` entfernt | `make test` **EXIT=2** | `TestNewChangeOriginClosedSet`, `TestToChangeCarriesSchemaAndTable`, `TestToChangeRejectsRowOutsideInvariants`, `TestChangeRowOriginMissingReadsAsWALAndUnknownIsRejected`, `TestReadChangesCarriesOrigin` |
| M5 | View liest `NULL` als `wal` | View in `schema.yaml`: `COALESCE(c.origin, 'wal')` → `c.origin` | `make test-store` **EXIT=2** | `TestChangesViewCarriesOriginLikeReadChanges` |
| G5 | Bekannt-Liste des Guards bindet end-to-end (F-2) | Prüfung `knownForeignObjects` in `guard.go:107` deaktiviert (`if false && …`) | ungekürzter Guard-Test **EXIT=1** | „Lauf 6a lief durch (Exit 0), obwohl ein unbekannter destruktiver Blocker vorlag" — Lauf 6a wird rot; Läufe 1–5 liefen davor grün |

Baum nach jeder Rücknahme: `git status --short` leer, Guard-Test-Lauf räumte Container,
Netz und Temp-Verzeichnis ab und stellte `plan.yaml`/`down.sql` wieder her.

## 6. Harte Regeln

- **§3.3** — die zwei Lifecycle-Moves sind reine Renames: `git show --stat -M 3b7f828a`
  (`{open => next}/… | 0`), `git show --stat -M f4ba82ab` (`{next => in-progress}/… | 0`);
  die Inhaltsänderung „Verantwortlich" (`f8ae1fc6`) liegt in einem eigenen Commit
  dazwischen.
- **§3.5** — keine bestehende `Accepted`-ADR verändert (`make doc-immutable` grün).
- **§3.6** — siehe §3.2, Urteil zur Lockerung; kein Coverage-Schwellenwert, keine
  Linter-/Architekturregel gesenkt.
- **§3.7** — diff-skopierter Kandidatenlauf über die `+`-Zeilen der geänderten `*.go`,
  `Makefile`, `*.sh`, `*.yaml` (Chronik-/Vorher-Vokabular, Konjunktive, Slice-/Wellen-/
  Finding-Kennungen): die Treffer sind Fehlermeldungen des Guard-Test-Skripts („… statt 0"),
  ein Kommentar der `bootstrap.sh` („baut es neu auf") und die Herkunftsanker „real
  gemessen" in den Kommentaren von `guard.go`/`guard_test.go`/`report.go` (§3.12 Ursprung,
  Indikativ); der Satz „statt bei jedem Aufrufer einzeln" in `report.go` steht schon am
  Parent-Stand. Kein Befund.
- **§3.11** — `docs-check` 0 Befunde (`hostpaths`); der Report trägt keinen host-lokalen
  Pfad, auch nicht in Grep-Mustern.
- **§3.12** — Handbuch „Richtwert rund 7 Sekunden" trägt den Ursprung („einzelne
  Architect-Messung auf einer Testinstanz mit einer Zeile (`ADR-0114`) … nicht garantiert");
  Coverage-Zahlen im Plan sind mit dem Lauf-Ursprung geführt (§1 nennt die gedruckten
  Zeilen meiner Läufe); Suchlauf-Zahlen nachgemessen (§3.3), zwei weichen ab (V-1, V-2).
- **§3.13** — Suchlauf-Felder vollständig (Erstlauf, Fixrunde 1, Fixrunde 2), beide
  Stände, Gefundenes und Nichtgefundenes; fremde Träger gemeldet, nicht mitgeändert:
  die sechs SDK-Fundstellen (Zeile in `slice-backfill-sdk-origin`), `.proto`/`.pb.go`,
  `run-integration-tests.sh` Z. 2711, die Such-Aufträge in `slice-backfill-sql-administration`
  und `slice-transformationen-antragsweg-schema`. Die Fremd-Pläne der Backfill-/
  Transformationswelle ändert allein der Planner-Träger `ae97247a`; im Implementer-Diff
  ab `ae97247a` liegt außer dem eigenen Plan nur die Meldungszeile in
  `slice-backfill-sdk-origin` (V-8).
- **Handbuch-Version** — `Version: 1.46`, drei Historienzeilen (1.44 Erstlauf, 1.45 und
  1.46 Fixrunden), `Stand: 2026-09-24`, im selben Diff wie die inhaltliche Änderung.
- **Commit-Traceability** — 20 Commits, jede Message nennt `LH-*`/`ADR-*`, keine
  `SPEC-`/`ARC-`-Kennung im Betreff (§1).

## 7. Befunde

Kein V-Befund blockiert. V-1 ist ein LOW (Zahl im Träger), V-2 bis V-8 sind INFO ohne
erwartete Aktion im Slice.

- **V-1 (LOW) — Plan §3, Fixrunde-2-Suchlauf: „je Stand 18 Treffer" nicht
  reproduzierbar.** Die Zelle „Anzahl-Formulierungen" nennt für
  `git grep -n -i -E '(ten|zehn) (fields|felder)' <Stand> -- sdks examples` **18**
  Treffer je Stand; ich messe **16** an allen vier Ständen (`09386619`, `f4ba82ab`,
  `e95937cf`, `HEAD`). Die Zahl steht als „gemessen" im Träger (`AGENTS.md` §3.12
  Instanz A). Die Substanz — welche Treffer Träger der bewegten Eigenschaft sind (die
  zwei „same ten fields as the domain type" in `Http/Models/Changes.cs`/`http/model/Changes.kt`)
  und dass die übrigen Live-Wege beschreiben — bleibt tragfähig. Nachzug: Zahl
  korrigieren oder den Zählweg nennen (beim Closure-Nachzug des Planners).
- **V-2 (INFO) — Plan §3: „24 (unverändert)" gegen `HEAD` jetzt 25.** Der Zusatztreffer
  ist die vom Plan angelegte Meldungszeile in `slice-backfill-sdk-origin`. Die Zelle
  bleibt für den Stand ihrer Niederschrift wahr; eine Zelle, die ein Suchergebnis über
  Pläne führt, verschiebt sich mit der eigenen Meldung (Klasse wie V-4 im Report zu
  `slice-backfill-row-image-gemeinsam`: Stand-relative Zelle).
- **V-3 (INFO) — ACL-Lücke: Urteil, ADR-Lücke, nicht Slice-Befund.**
  [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 3
  sagt, die Rechte setze `nacharbeit-roles.sql` „im selben Lauf". Das gilt nur für
  `cdc_reader` (`nacharbeit-roles.sql` Z. 105 grantet `SELECT` auf `cdc.changes` an
  genau diese Rolle); `DROP VIEW` verwirft die gesamte ACL der View, ein vom Betreiber an
  eine andere Rolle vergebenes Recht geht verloren. Der Slice setzt die ADR wörtlich um
  und trägt die Grenze an drei Stellen: Handbuch §4 „Eigene Rechte" (1.46), Kommentar am
  `schema-rollout`-Target, Plan-Zeile „`ADR-0114` (F-5) Lücke gemeldet, nicht geändert".
  Die Verletzung liegt nicht im Slice, sondern in der Lücke der (immutable) ADR-Aussage.
  **Notiz an den Architect:** die ADR nennt den Verlust der Rechte anderer Rollen nicht;
  ob das eine Schärfung (neue ADR mit `Supersedes`) oder eine Zitat-Korrektur
  ([`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md), Grenze:
  §Konsequenzen sind unberührbar) braucht, entscheidet der Architect.
- **V-4 (INFO) — `MANUAL_ACTION_REQUIRED` enger als der ADR-Text: zulässig.** Siehe §4;
  Fail-closed-Richtung, im Plan und im Code-Kommentar begründet, real gemessen.
- **V-5 (INFO) — Sensor-Ausgabe: Meldung „nur bekannte Fremdobjekt-Blocker (ADR-0043)"
  bei gemischtem Report.** Bei Läufen mit View-Signatur-Blocker und den sechs bekannten
  Fremdobjekten (Lauf 4 und 5) druckt das Target „nur bekannte Fremdobjekt-Blocker …
  --execute laeuft mit --allow-destructive" und danach den Vorlauf; „nur" ist dort
  ungenau (es blockiert zusätzlich die View-Signatur). Der Guard-Test bindet die
  Meldung als Merkmal des Pfads (Lauf 2, Lauf 6a), nicht ihren Wortlaut im gemischten
  Fall; die Überschrift des Lauf-4-Nebenfalls trägt „Lauf 4 (abhängiges Objekt …)" ohne
  „/6". Kosmetik, keine Zusage betroffen.
- **V-6 (INFO) — `make test-store`/`make test-replication` schreiben `tools/schema/plan.yaml`
  (nur `target`-Zeile) in den Arbeitsbaum.** Eigenschaft des Rollout-Pfads dieser Tiers
  (der Commit `5f126971` stellt die Compose-Bezeichnung wieder her), nicht dieses Slice;
  der Arbeitsbaum ist nach einem Tier-Lauf ohne Nachzug dirty. Zur Kenntnis.
- **V-7 (INFO) — DoD-Zeile „Reconciliation-Register — entfällt" auf `[ ]`.** Im
  Vorgänger-Slice auf `[x]` gesetzt; beim Closure-Nachzug einheitlich behandeln.
- **V-8 (INFO) — Meldung im fremden Plan.** Die Zeile in
  `docs/plan/planning/open/slice-backfill-sdk-origin.md` schreibt der Implementer in
  einen Plan, den der Planner führt; sie ist im Plan des Slice deklariert („update,
  gemeldet") und ändert keine SDK-Datei. Regelkonform als Übergabe-Artefakt
  (`AGENTS.md` §3.13), Inhalt der Zeile (sechs Fundstellen, zwölf Felder am Parent)
  von mir nachgemessen.

## 8. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (8) | **8 von 8 erfüllt**, je mit eigenem Beleg (Zeile 7 mit V-1/V-2) |
| DoD §2 — „Review durchgeführt" | **erfüllt** (Report liegt vor; F-1…F-8 selbst am Code nachgemessen) |
| DoD §2 — übrige `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Closure-/Rollen-Sequenz; Reconciliation „entfällt", V-7) |
| Plan-vs-Code-Diff | **deckungsgleich**: alle Plan-Zeilen im Diff vertreten, ungeplant nur die im Plan benannten Nachträge (`translate.go`, `schema.sql`, drei Live-Weg-Kommentare, `examples/bootstrap.sh`); keine Datei außerhalb; `test/integration` bewusst unberührt (Diff leer) |
| Kernaussagen `origin` (Domäne, Store, View, HTTP, Live-Wege unverändert) | **bestätigt**, 5 Eingabeseiten-Mutationen rot |
| Sicherheitsnaht `rolloutguard` | **bestätigt**: Unit-Tests EXIT=0, Guard-Test EXIT=0, Mutation G5 rot an Lauf 6a, F-1 (M3) rot |
| Upgrade-Belege | **bestätigt**: `f4ba82ab` → Arbeitsbaum (Exit 0 / 0, Vorlauf nur im ersten Lauf, Zeile `wal`, Recht `t`) und `v0.1.2` (Guard-Test Lauf 5) |
| DB-Sensoren | `make test-store` EXIT=0, `make test-replication` EXIT=0 (je `DB-Adapter-Coverage: 77.13% (gedeckt 533 von 691 …)`, Schwelle 70 %), `make test-integration` EXIT=0 |
| Entscheidungen `ADR-0111`/`0043`/`0064`/`0114`/`0058`/`0081` | **konform** (ADR-0114-Blocker-Grund enger: zulässig, V-4; ACL: ADR-Lücke, V-3) |
| Harte Regeln (§3.3, §3.5, §3.6, §3.7, §3.9, §3.11, §3.12, §3.13, Handbuch-Version) | **erfüllt**; §3.12 mit V-1/V-2 |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=09386619..HEAD`) |
| Gates | **`make gates` EXIT=0**, `make test` EXIT=0, `make a-check` EXIT=0, `make coverage-gate` EXIT=0 (`82.90%`) |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung belegt: das
Feld `origin` ist von der Domäne bis zu `cdc.changes` und `GET /changes` durchgängig, die
Live-Wege bleiben bei zehn Feldern, die Sicherheitsnaht des Schema-Rollouts hält (alles
oder nichts, Vakuum-Fall abgelehnt, Bekannt-Liste end-to-end gebunden) und das Upgrade
eines Alt-Schemas läuft über die View-Signaturänderung mit zweimal Exit 0 und lesbarem
Datenstand. Fünf Eingabeseiten-Mutationen und die Guard-Test-Mutation färben rot. Kein
Befund blockiert; V-1 (Zahl „18" gegen gemessene 16) ist ein LOW im Plan-Träger, V-3 ist
als ADR-Lücke an den Architect gemeldet.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: Closure-Notiz mit
Lerneintrag (Finding-Klassen des Reviews in das Beobachtungs-Register), Ausgänge der vier
§6-Risiken (Konvergenz: belegt in §3 des Plans und in diesem Report; Doppelquelle:
`TestChangesViewCarriesOriginLikeReadChanges`; strikte Dekoder: erst mit
`slice-backfill-sdk-origin` real belegbar; E2E-Zeilenanker: Diff leer, Tabelle unverändert),
Reconciliation-Zeile (V-7), Welle-Paarungen; V-1/V-2 im Suchlauf-Feld nachziehen; die
ACL-Notiz V-3 geht an den Architect. Danach der reine `git mv` nach `done/`
(`AGENTS.md` §3.3, Fall 2).
