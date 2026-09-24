# Verifikations-Report: slice-backfill-sql-administration — 2026-09-24

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-backfill-sql-administration.md`](review-slice-backfill-sql-administration.md);
Formvorbild dieses Reports:
[`verifikation-slice-backfill-run-store.md`](verifikation-slice-backfill-run-store.md). Die vendored
Baseline trägt kein eigenes Verifikations-Template (`.harness/baseline/v6.9.0/templates/docs/reviews/`
enthält nur `review-report.template.md`); der Report folgt dem Formvorbild.

**Gegenstand:** Slice-Plan `slice-backfill-sql-administration` (Welle `welle-backfill-bestand`),
Diff-Range `2d47d8a7..HEAD` (`c092efa5`), 21 Commits: Lifecycle und Verantwortlich (`2665e422`,
`9575714d`, `5be20d2f`), Umsetzung (`c561d52d` Speicher, `a7b00be8` Verdrahtung, `1fa9534b` Spec,
`c188be43` Handbuch, `00690fc8` Worker-Test), Plan-Nachzüge (`d524e390`, `d03ae183`), Review
(`667e27da`, `b8b2dfd5`), Fixrunde (`b019cc2e`, `cabc6d14`, `947d9960`, `c0c6e286`, `120139ec`,
`d5e5ea99`, `55085f76`, `4981766a`, `c092efa5`). Dieser Lauf ändert weder Code noch Plan, Spec oder
Doku; er schreibt nur diesen Report. Alle Mutationen liefen im Arbeitsbaum und sind je Lauf per
`git checkout` zurückgenommen (`git status --short` leer nach jedem Lauf); `tools/schema/plan.yaml`
und `tools/schema/down.sql` (Nebeneffekt der Tier-Läufe) sind ebenfalls zurückgenommen.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert
(`make … > <log> 2>&1; echo EXIT=$?`), die Logs danach gelesen. Stand aller Läufe: `HEAD` = `c092efa5`.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `db-package-lists-check: OK — … dieselben 4 Pakete` · `coverage-gate: OK — Coverage 82.90% erfüllt Schwelle 80%` · `d-check: 1050 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make test` (Race-Detector) | **EXIT=0** | 42 Pakete `ok`, 0 `FAIL` |
| `make a-check` | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` (eigener Lauf) | **EXIT=0** | `coverage-gate: OK — Coverage 83.00% erfüllt Schwelle 80%` (der Lauf im `make gates` druckte `82.90%`: lauf-gebundene Streuung) |
| `make test-store` (PostgreSQL 18, Repo-Default), zweimal | **EXIT=0** | `ok … postgresstorage 7.168s`; `DB-Adapter-Coverage: 82.08% (gedeckt 843 von 1027 Statements; Profile gemergt: store,replication)`; `db-coverage: OK — … erfuellt Schwelle 70%` |
| `bash tools/harness/run-schema-rollout-guard-test.sh` (ungekürzt) | **EXIT=0** | `Lauf 5 OK — Tag v0.1.2: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar; 19 Tabellen-/View-Rechte der drei Rollen auf administration_request, backfill_run und backfill_status, EXECUTE auf cdc.backfill_table(text, text, text) allein für cdc_admin (nicht PUBLIC), request_kind-Menge backfill,disable,enable,exclude_column,include_column`; Schlusszeile `OK — alle Belege real erbracht (…)` |
| `make commit-traceability RANGE=2d47d8a7..HEAD` | **EXIT=0** | `OK — 21 Commit(s) in "2d47d8a7..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=2d47d8a7..HEAD` | **EXIT=0** | `d-check: 1050 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=2d47d8a7..HEAD` | **EXIT=0** | `d-check: 1050 Datei(en) geprüft, 0 Befund(e)` |
| `gofmt -l internal tools cmd` (gepinntes Toolchain-Image) | EXIT=0 | listet nur `receive/seam_test.go` und `outbound/log_test.go` — beide Bestand, beide nicht im Diff |

**Übernommen, nicht selbst gefahren:** `make test-replication` (Replication-Teil der
DB-Adapter-Coverage). Die Zahl `82.08% (843 von 1027)` des Store-Laufs merget mit dem zuletzt im
`DB_COVERAGE_DIR` liegenden Replication-Profil; ihren Replication-Anteil habe ich nicht frisch
gemessen. Übernommen von Reviewer (Stand `d03ae183`: `81.76% (838 von 1025)`) und Implementer
(Fixrunden-Stand: `82.08% (843 von 1027)`, dieselbe Zahl wie mein Lauf). Der Replication-Tier trägt
keinen Diff-Anteil dieses Slice.

Tier-Hygiene: dangling-Volumes vor dem Lauf **34**, nach dem Lauf **34**. Die sechs anonymen Volumes
(`7a09dcab…`, `43657134…`, `6fde172c…`, `6a8a9f3a…`, `cd8512e7…`, `1f8f00ef…`), die meine Tier-Läufe
erzeugten (Erzeugungszeit im Lauf-Fenster), habe ich je einzeln per `docker volume rm <name>`
entfernt; kein `prune`. Speicher vor den Tier-Läufen: 19,5 GB verfügbar (`free -m`).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Antragsweg (`cdc.backfill_table`, `PUBLIC` ohne `EXECUTE`, fünf `request_kind`-Werte, Guard kennt die Funktion, zweimal Rollout, Alt-Tag-Lauf, `plan.yaml`/`down.sql`) | **erfüllt** | Guard-Skript EXIT=0, Lauf 5 gedruckt (§1): Alt-Tag `v0.1.2` Exit 0/0/0, `EXECUTE` allein für `cdc_admin`, fünf Werte, 19 Rechte; Lauf 2 und 5 tragen den zweiten Rollout (Exit 0) — ein eigener `make schema-rollout`-Doppellauf gegen eine dauerhafte DB ist nicht Teil meines Laufs (der Guard-Test ist sein Beleg, wie im Plan §2 benannt). `plan.yaml`/`down.sql`: der frische Rollout im Store-Tier erzeugt `down.sql` byte-gleich und `plan.yaml` gleich bis auf die `target`-Zeile (`git diff` gedruckt: genau diese eine Zeile); `down.sql` trägt `DROP VIEW "backfill_status";`. `guard_test.go` grün (`make test`) |
| 2 | Verarbeitung (Zweig `backfill`, Reihenfolge, Neustart, Signal, Aufnahme ohne Bindung, fremde Quelle `pending`, Aufrufer-Vertrag `Execute`, Bezug Antrag/Run) | **erfüllt** | `make test` EXIT=0 und `make test-store` EXIT=0; Testnamen gelesen und gefahren: `TestBackfillWorkerRunsQueuedRunsAtStartInOrderAndSkipsTheOthers`, `…WakesForARunAdmittedAfterItsStart`, `…ReadsAgainAfterEachRun`, `…KeepsASignalThatArrivesBetweenReadAndWait`, `TestSignalBackfillWorkerNeverBlocksAndCoalesces`, `…WritesNoEndStateAndFilesNoRequestFromAFailedResult` (Fake-Use-Case liefert `failed`), `TestBackfillRunDoesNotBlockTheAdministrationGoroutine`, `TestReconcileBackfillRunsAgainstPostgreSQL`, `TestBackfillWorkerEndsARunOfAnUnboundTableAsConfigurationFailure` (reale Run-Tabelle, Klasse `configuration`), `TestAdministrationBackfillBranchLeavesTheRequestOfAnotherSourcePending`, Login-Test `TestAdministrationPathRunsUnderLeastPrivilegeLogins` (zweiter Antrag `failed` „aktiver Backfill-Run“, keine zweite Run-Zeile, nicht aktivierte Tabelle `failed`); Bezug von Antrag und Run: `…RefusesARequestThatIsNotTheBackfillOfThisRun` je Art/Quelle/Schema/Tabelle (Reviewer-Mutationen S2/S5/S6 rot, von mir nicht wiederholt). Der Zweig „Publication-Mitgliedschaft fehlt“ ist im Use Case getestet (`usecase/backfill/service_test.go`), im Worker-Test nur der Bindungs-Zweig real (V-4) |
| 3 | Sichtbarkeit (View, Grant, `diagnose`, Handbuch-Abschnitt, Fortsetzungs-Idiome) | **erfüllt** | `make test-store` grün: `TestBackfillStatusViewShowsTheLatestRunPerTable`, `TestCdcReaderRoleReadsTheBackfillRunThroughTheStatusView`, `TestDiagnoseReportsTheLatestBackfillRunPerTable`, `TestChangesViewKeysetContinuesInsideOnePosition`; Format-Test `TestFormatBackfillRunNamesTheEstimateAsEstimatedAndUnknownAsUnknown`; Guard-Lauf 5 trägt `SELECT` auf `cdc.backfill_status` für `cdc_reader` (19 Rechte); Handbuch `Version: 1.49`, Historienzeilen 1.48 und 1.49 (§6) |
| 4 | Spec-Zug (`SPEC-019` Grants, `LH-FA-CAP-009.a` Schema-Version) | **erfüllt** | `git diff 2d47d8a7..HEAD -- spec/pflichtenheft.md`: Absatz „Grants“ (`cdc_admin` `SELECT`/`UPDATE`, `cdc_capture` und `cdc_reader` keines, niemand `INSERT`/`DELETE`) deckt sich mit `nacharbeit-roles.sql` und dem Rechte-Lauf 5; Satz zur Schema-Version; Historienzeile; keine ADR-/Slice-/Wellen-Kennung im Spec-Diff; `make docs-check` 0 Befunde |
| 5 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 6 | Review durchgeführt, kein offenes HIGH/MEDIUM nach der Fixrunde | **erfüllt, mit V-2 (INFO)** | Report liegt vor (3 HIGH, 3 MEDIUM); Behebung je Finding selbst nachgemessen (§4); kein zweiter Review-Report zur Fixrunde |
| 7 | §3.13-Suchlauf (Feld in §3, beide Stände) | **erfüllt, mit V-1 (LOW)** | neun Zeilen des Feldes nachgefahren (§5): alle gedruckten Zahlen stimmen; ein Zählwort-Rest („beiden Tabellen-Antragsarten“) steht außerhalb des Suchmusters (V-1) |
| 8 | Doku-Update | **erfüllt** | siehe Zeile 3, §6 (Träger-Nachzüge) |
| 9 | Closure-Notiz mit Lerneintrag | **korrekt offen** | §7 des Plans trägt „*(zu tragen bei Closure)*“ |
| 10 | Reconciliation-Register — entfällt | **korrekt offen / entfällt** | Greenfield; `[ ]` wie bei den Vorgängern |
| 11 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | alle neun Risiko-Zeilen tragen „**Ausgang:** *(bei Closure)*“ (gezählt am Plan) |
| 13 | Die drei Paarungen | **korrekt offen** | hängen an der Closure von `welle-backfill-bestand` |

`[x]` sind acht Zeilen (Nr. 1–8), `[ ]` fünf (Nr. 9–13), zusammen 13 — gezählt am Plan (§2 hat 13
Aufzählungspunkte). Kein `[x]` ohne Beleg, kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre.

## 3. Kernaussagen, selbst gemessen

### 3.1 Alt-Bestand-Beleg im Repo (Review F-1)

Lauf 5 des Guard-Skripts prüft jetzt im Skript, was `harness/targets/schema-rollout.md` §Belege
(5) zusagt: Vorbedingungen am Alt-Stand (Funktion fehlt, CHECK ohne `backfill`, `cdc_admin` ohne
`UPDATE`), nach dem Upgrade 19 Rechte der drei Rollen, `EXECUTE` allein für `cdc_admin` (auch nicht
`PUBLIC`) und die fünf Werte der CHECK-Menge. Gedruckte Zeile: §1. Bindung an die Eingabeseite: siehe
§4, M-F1.

### 3.2 `plan.yaml`, `down.sql`

Zweiter Store-Lauf (frische Wegwerf-DB): `git diff` nach dem Lauf zeigt allein die `target`-Zeile von
`plan.yaml` (`postgres://postgres:***@cdc-test-postgres:5432/cdc?…` committet, im Tier
`postgres://cdc:***@cdc-store-test-pg:5432/cdc_test?…`); `down.sql` unverändert. Der committete Stand
ist damit das Ergebnis eines Rollouts gegen eine leere Datenbank.

### 3.3 Umfang

`git diff --shortstat 2d47d8a7 c188be43`: 33 Dateien, 2106 Einfügungen; `git diff --numstat` über `*.go`:
Produktion 415 Zeilen in neun Dateien (413 in acht ohne `guard.go`), Test 1377 in 14 Dateien — deckt sich mit der
Umfangsentscheidung in Plan §3. Die Lifecycle-Moves `2665e422` und `5be20d2f` sind reine Renames
(`{open => next}` bzw. `{next => in-progress}`, je `0 insertions, 0 deletions`); `Verantwortlich` steht in
`9575714d` (1 Zeile, in `next/`) vor dem zweiten Move ([`AGENTS.md`](../../AGENTS.md) §3.3).

## 4. Fixrunde: Review-Findings nachgemessen (Mutationen an meinen Läufen)

Aufbau je Mutation: Datei im Arbeitsbaum mutiert, Lauf ungefiltert in eine Log-Datei, gedruckte
Rot-Meldung gelesen, Datei per `git checkout` zurückgenommen (Endstand `git status --short` leer).

| # | Finding | Mutation | Rot gesehen (gedruckt) |
|---|---|---|---|
| M-F1 | F-1 (HIGH) `EXECUTE` allein für `cdc_admin` von Lauf 5 gebunden | `cdc.backfill_table(text, text, text)` aus der `GRANT EXECUTE`-Zeile von `nacharbeit-administration.sql` gestrichen, `bash tools/harness/run-schema-rollout-guard-test.sh` | **EXIT=1**, `FEHLER — Lauf 5: cdc_admin trägt nach dem Upgrade über v0.1.2 auf cdc.backfill_table(text, text, text) das Recht EXECUTE nicht als t` |
| M-F2 | F-2 (HIGH) Quellfilter von `diagnose` | `WHERE source_id = $1` → `WHERE $1::text IS NOT NULL` in `internal/bootstrap/backfill.go` (`make test-store`) | **EXIT=2**, `--- FAIL: TestDiagnoseReportsTheLatestBackfillRunPerTable` (nur dieser Test; vor der Fixrunde blieb dieselbe Mutation grün) |
| M-F9 | F-9 (LOW) „Antrag erledigt“ nur bei vermerkter Zeile | `tag.RowsAffected() > 0` → `>= 0` in `MarkApplied` (`make test-store`) | **EXIT=2**, `--- FAIL: TestAdministrationRequestAdapterMarkAppliedLogsOnlyARequestItMarked` — `Meldungen "administrationrequest: Antrag erledigt" nach dem zweiten MarkApplied = 2, erwartet 1` |
| M-F10 | F-10 (INFO→umgesetzt) fremde Quelle bleibt `pending` | Bedingung `request.Kind == … && request.Source != deps.source` durch `if false` ersetzt (`make test`) | **EXIT=2**, `--- FAIL: TestAdministrationBackfillBranchLeavesTheRequestOfAnotherSourcePending` |

Die Mutation M-F2 lief zusammen mit M-F9 im ersten Store-Lauf; das `set -e` des Runners endete nach dem
roten Paket `internal/bootstrap`, `postgresstorage` lief nicht mehr — M-F9 habe ich deshalb einzeln wiederholt
(Zeile oben).

Die übrigen Findings, ohne Mutation, am Träger nachgemessen:

- **F-3** (Suchlauf-Zeile): der Befehl `grep -n -e 'letzte-gelesene' -e 'letzte gelieferte' -e 'LIMIT' -e 'backfill'
  docs/user/benutzerhandbuch.md` druckt an `2d47d8a7` 10, an `c188be43` 40 und an `55085f76` 41 Zeilen; `-e 'backfill'`
  allein 7, 30 und 31 — die Plan-Zeile nennt genau diese Paare (§5). **Behoben.**
- **F-4** (Status-Abfrage unter `cdc_admin`): das Handbuch nennt die `cdc_reader`-Mitgliedschaft (`CDC_READER_DSN`)
  für das Lesen des Runs und sagt, dass `cdc_admin` kein `SELECT` auf die View trägt (`git diff c188be43..HEAD`);
  Lauf 5 belegt `SELECT` auf `cdc.backfill_status` allein für `cdc_reader` (19 Rechte). **Behoben.**
- **F-5** (Stichtag): `grep -n 'Zeitpunkt des Antrags' docs/user/benutzerhandbuch.md` druckt 0 Zeilen
  (an `d5e5ea99`: 2); beide Stellen nennen den Startzeitpunkt des Runs. **Behoben.**
- **F-6** (Zählwörter): `grep -rn -e 'vier Antragsarten' -e 'die vier Views' -e 'vier Antrags' -e 'Folge-Slice' internal harness tools`
  druckt an `2d47d8a7` 18, an `c188be43` 16, an `4981766a` 10 Zeilen; die zehn Reste sind unabhängige Begriffe (Carveout-Text in
  `harness/README.md`/`harness/conventions.md`, Retention, drei SDK-Läufer). `harness/targets/schema-rollout.md` trägt „die fünf Views“
  (Zeile 12) und „fünf deklarierten Views“; `queries.go:306-307` zählt fünf Antragsarten und drei Tabellen-Antragsarten.
  **Behoben, mit V-1** (Zählwort „beiden“ in `queries.go:324`).
- **F-7** (Fremdobjekte): `grep -rn -e 'sechs' harness Makefile tools docs/user` druckt an `HEAD` 7 Zeilen; im Skript
  stehen „sieben“ (Zeilen 9, 14, 50, Ausgabe von Lauf 3 in Zeile 177), die verbleibende „sechs“ (Kopf, Zeile 4) meint die sechs Läufe.
  **Behoben.**
- **F-8** (`gofmt`): gedruckt (§1) — die vier vom Review gemeldeten Dateien fehlen in der Liste. **Behoben.**

## 5. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent `2d47d8a7` per `git grep -n … 2d47d8a7 -- <Pfade>`, Diff-Stand per `grep -rn` bzw. `git grep … c188be43` und
`… 4981766a`/`55085f76` (die Stände, die der Plan je Zelle nennt); gedruckte Zeilenzahlen:

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| Antragsarten (`exclude_column`, `internal tools docs spec harness`): Parent 103, Diff-Stand 105 | Parent **103**, `c188be43` **105** (am `HEAD` 111: zählt Review-Report, Plan und Fixrunde mit) | bestätigt am benannten Stand |
| „sechs“ (`harness Makefile tools docs/user`): Parent 19, Diff-Stand 11, Fixrunde 7 | Parent **19**, `c188be43` **11**, `4981766a` und `HEAD` **7** | bestätigt |
| Zählwörter (siehe §4 F-6): Parent 18, Diff-Stand 16, Fixrunde 10 | **18**, **16**, **10** | bestätigt |
| Stichtag: vorher 2, danach 0 | `d5e5ea99` **2**, `55085f76` und `HEAD` **0** | bestätigt |
| `diagnose` (`docs/user internal/bootstrap`): Parent 52, Diff-Stand 76 | Parent **52**, `c188be43` **76** (`HEAD` 78) | bestätigt |
| Replication-Verbindungen (`max_wal_senders`, …): Parent 1, Diff-Stand 3 | **1**, **3** (`c188be43` und `HEAD`) | bestätigt |
| Fortsetzungs-Idiome: 10 → 40 bzw. 7 → 30; Fixrunde 41 und 31 | Parent **10**/**7**, `c188be43` **40**/**30**, `55085f76` **41**/**31** | bestätigt (F-3 behoben) |
| Läufe des Guard-Tests: gedruckte Zeile des Endlaufs | identisch zu meinem Lauf (§1) | bestätigt |
| Nenner der Coverage-Messungen: `1025`, `82.90%`; Fixrunde `83.00%`, `82.08% (843 von 1027)` | mein Lauf: `82.90%` (`make gates`), `83.00%` (`make coverage-gate`), `82.08% (843 von 1027)` (`make test-store`) | bestätigt; die Streuung `82.90`/`83.00` ist lauf-gebunden |
| Umfang (Plan §3): 33 Dateien, 2106 Zeilen, 413 in acht, 1377 in 14 | §3.3 | bestätigt |

Eigener Zusatz-Suchlauf nach den Zählwörtern „beiden“/„zwei“ im Umfeld von Antragsarten
(`grep -rnE 'beiden (Tabellen|Spalten)?-?Antrags|…' internal tools harness spec docs/user Makefile`): ein Träger
trägt eine überholte Zahl — `internal/adapters/driven/postgresstorage/queries/queries.go:324` („Die beiden
Tabellen-Antragsarten bleiben außen vor“ vor der Abfrage `request_kind IN ('exclude_column', 'include_column')`; außen vor
bleiben `enable`, `disable` **und** `backfill`) — V-1. Die übrigen Treffer sind zutreffend: „beiden Spalten-Antragsarten“
(zwei), `model/administrationrequest.go` zählt `enable`/`disable`, zwei Spalten-Antragsarten und `backfill` einzeln auf,
`wiring.go:1345` nennt `enable`/`disable` als Träger von Bindungs- und Publication-Menge.

## 6. Träger-Nachzüge

| Träger | Beleg | Verdikt |
|---|---|---|
| Handbuch `Version: 1.49`, Historienzeilen 1.48 und 1.49 | `docs/user/benutzerhandbuch.md` Zeilen 3, 1548, 1549; Abschnitt „Bestand als Backfill überführen“, §2 Rollen, §4 Diagnose, Glossar, beide Fortsetzungs-Idiome | erfüllt |
| Handbuch trägt keine Aussage ohne Beleg (Startposition eines frisch registrierten Consumers, Richtgröße) | `grep -n -i 'Richtgr\|Startposition' docs/user/benutzerhandbuch.md`: kein Treffer | erfüllt |
| `harness/README.md` (Zeile `make example-demo-up`) | „sieben bekannten Fremdobjekte“ | erfüllt |
| `harness/targets/schema-rollout.md` | fünf Views mit `backfill_status`, fünf Funktionen, „sieben Objekte der Nacharbeit-Dateien (fünf Funktionen, die Views `metrics` und `heartbeat`)“, §Belege (5) beschreibt Lauf 5 so, wie er läuft (§1) | erfüllt |
| Spec-Grants ([`SPEC-019`](../../spec/pflichtenheft.md)) | §2 Zeile 4 | erfüllt |
| Guard-Eintrag | `guard.go` trägt `backfill_table(in:text,in:text,in:text)` als siebten Eintrag; `guard_test.go` grün | erfüllt |

## 7. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 5 | Antragsart `backfill`, Funktion `cdc.backfill_table` schreibt nur den Antrag und sendet `pg_notify`, Administrations-Goroutine blockiert nicht, Run-Zeile `queued` + Antrag `applied` („angenommen“), ein Worker, Diagnose über `cdc.backfill_status`, `failed`/`interrupted` ist Berichtsinhalt | `nacharbeit-administration.sql` (Guard-Lauf 5, Reviewer-Mutation S1); `TestBackfillRunDoesNotBlockTheAdministrationGoroutine`; `applyAdministrationRequest` (Zweig `backfill` ruft `Request`, weckt danach); `backfill.go` (`diagnoseBackfillStatus`, `formatBackfillRun`) | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4 | `running` → `interrupted` beim Start, `interrupted` startet nicht von selbst | `reconcileBackfillRuns`, `TestReconcileBackfillRunsAgainstPostgreSQL` (Store-Tier); `Queued` liest nur `queued` | konform |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 | Annahme in **einer** Transaktion über den Pool von `cdc_admin`, je Quelle nimmt eine Instanz Anträge an | Login-Test unter `cdc_admin`/`cdc_capture`; die zweite Hälfte der Festlegung ist mit der Fixrunde für die Art `backfill` am Code gebunden (M-F10); für die vier übrigen Antragsarten bleibt der Bestand aus [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md) (Grenze im Plan §3 und §6 benannt) | konform, Grenze benannt |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2 | Worker arbeitet zuerst alle `queued`-Zeilen in `(requested_at, run_id)` ab, wartet dann; Aufnahme beim Start und bei Wecksignal (Kapazität 1, verschmelzend); Start-Reihenfolge Bindung → Abgleich → Worker | `runBackfillWorker`/`drainBackfillQueue`/`newBackfillWake`; `wiring.go` Start-Reihenfolge; die sieben Worker-Tests (§2 Zeile 2); Reviewer-Mutationen M1–M6 | konform |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3 | `NULL` = unbekannt (nie 0), zwei Warn-Spalten in View und `diagnose`, `false` bis zur Auswertung | View-Test (Ordnung, Gleichstand, NULL/0, Warn-Spalten); `TestFormatBackfillRun…`; Reviewer-Mutation D1 | konform |
| [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) | Alt-Tag-Lauf, Guard-Eintrag, View mit endgültiger Signatur von Anfang an | Guard-Lauf 5 Exit 0/0/0, Guard-Test EXIT=0 | konform |
| [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md) Folgepflicht 2/3 | Spec- und Handbuch-Satz zur Schema-Version einer Backfill-Change | `spec/pflichtenheft.md` Absatz „Markierung“ (Satz zur Schema-Version); Handbuch-Abschnitt trägt „von Änderungen einer später registrierten Version“ (Reviewer-Negativbefund, von mir gelesen) | konform (der Wortlaut „unterscheidet“ ohne Objekt ist der ADR-Wortlaut, Review F-13 b) |
| [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md) | Antrags-Queue, Live-Reload, Funktionen schreiben nur Antrag + `pg_notify` | Login-Test unter Least-Privilege-Logins (grün im Store-Lauf) | konform |
| [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md) | keine Domänenlogik in SQL | `cdc.backfill_table` legt nur die Antragszeile an (Reviewer-Negativbefund, `nacharbeit-administration.sql` gelesen) | konform |
| [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) | rollenspezifische DSN, Least Privilege | Pools je Rolle in `wiring.go` (`AdminDSN` Annahme, `CaptureDSN` Run-Zustand/Schreiber/Snapshot); Rechte-Matrix Lauf 5; `cdc_reader` liest nur die View | konform |
| [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) | neutrales Modell wo konvergent, sonst Nacharbeit-SQL + Guard-Eintrag | View im neutralen Modell (`schema.yaml`), CHECK und Funktion in `nacharbeit-administration.sql`, Guard-Eintrag | konform |

## 8. Plan-vs-Code-Diff

Jede Zeile der Plan-Tabelle §3 ist im Diff vertreten (`git diff --stat 2d47d8a7..HEAD`, 37 Dateien):
Modell/Adapter/Queries mit Tests, `internal/bootstrap` (`backfill.go` neu, `wiring.go`, sieben Testdateien),
Schema-Dateien (`schema.yaml`, `nacharbeit-administration.sql`, `nacharbeit-roles.sql`, `plan.yaml`, `down.sql`,
`rolloutguard`), Guard-Skript, Spec, Handbuch, `harness/README.md`, `harness/targets/schema-rollout.md`. Im Plan als
„nicht im Plan“ markierte Zusätze (`errors.go`, `outbound/administrationrequest.go`, `backfill.go` als eigene Datei,
`administrationDeps.source`) stehen im Diff und in der Tabelle. Ungeplant im Diff: nichts außer dem Plan selbst
(Lifecycle-Moves) und dem Review-Report. Kein Fremdverzeichnis im Baum.

## 9. Harte Regeln

- **§3.1** — alle Läufe über `make`/gepinnte Images; `gofmt` im gepinnten Toolchain-Image.
- **§3.2** — `git diff 2d47d8a7..HEAD -- '*.go' '*.sql' '*.sh' '*.yaml'` (Zeilen mit `+`) enthält kein `nolint`.
- **§3.3** — die zwei Moves sind reine Renames, die Inhaltsänderung `Verantwortlich` liegt dazwischen in eigenem Commit (§3.3).
- **§3.5** — `make doc-immutable RANGE=2d47d8a7..HEAD` EXIT=0; keine bestehende `Accepted`-ADR im Diff.
- **§3.6** — keine Schwelle gesenkt: `THRESHOLD`/`DB_COVERAGE_THRESHOLD`, `Dockerfile`, `harness/mk` nicht im Diff (`git diff --stat`).
- **§3.7** — diff-skopierter Kandidatenlauf über die `+`-Zeilen der geänderten `.go`, `.sql`, `.sh`, `.yaml` nach
  `Fixrunde`, `Review-`, `F-<n>`, `früher`, `bisher`, `vorher`, `nachher`, `wäre`/`würde`/`hätte`, `slice-`, `welle-`: ein Treffer
  („Ein früher Rücksprung“, gemeint ist „vorzeitig“), kein Befund.
- **§3.9** — Exit-Codes nie durch eine Pipe gemessen (§1).
- **§3.11** — `docs-check` 0 Befunde (`hostpaths`); dieser Report trägt keinen host-lokalen Pfad.
- **§3.12** — jede Zahl dieses Reports ist **gemessen** (Lauf in §1 bzw. Befehl in §4/§5 genannt) oder als **übernommen**
  gekennzeichnet (Replication-Anteil, Reviewer-Mutationen S1–S7, K1, D1, M1–M6).
- **§3.13** — Suchlauf-Feld an beiden Ständen nachgefahren (§5).
- **Handbuch-Pflicht** — `Version: 1.48` → `1.49` und Historienzeile je Zug.
- **Commit-Traceability** — 21 Commits, jede Message nennt eine Kennung, keine Struktur-Kennung im Betreff (§1).
- **§3.10** — der Diff berührt keinen GitHub-Actions-Workflow.

## 10. Befunde

Kein V-Befund blockiert. V-1 ist LOW, V-2 bis V-5 sind INFO.

- **V-1 (LOW) — Zählwort „beiden“ in `queries.go:324` steht falsch, der Plan-Suchlauf fand es nicht.** Der Kommentar
  vor `SelectAppliedColumnRequests` sagt „Die beiden Tabellen-Antragsarten bleiben außen vor“; die Abfrage schließt
  `enable`, `disable` und `backfill` aus (`request_kind IN ('exclude_column', 'include_column')`). Die Datei ist im Diff
  (die Fixrunde hat Zeile 306-307 derselben Datei nachgezogen); das Suchmuster der Plan-Zeile „Zählwörter“ enthält „vier“, nicht
  „beiden“, und die Zelle behauptet „Nicht gefunden: weitere Zählwörter „vier“ …“ — wahr für „vier“, nicht für „beiden“.
  Aktion: bei Closure den Satz auf „Die drei Tabellen-Antragsarten“ setzen. Kein Gate betroffen.
- **V-2 (INFO) — Fixrunde ohne eigenen Review-Report.** Die DoD-Zeile „Review durchgeführt“ ist durch den Report des
  Reviewers erfüllt; die Fixrunde habe ich in §4 selbst gemessen (vier Eingabeseiten-Mutationen rot, sechs Träger-Nachmessungen).
  Wer eine Reviewer-Durchsicht der Fixrunde verlangt, verlangt sie als eigenen Schritt.
- **V-3 (INFO) — Nenner-Träger.** `harness/sensors/db-adapter-coverage.md` (Nenner 1016, `81.69%`) und
  `harness/sensors/coverage-gate.md` (Nenner 2398) tragen ihre Zahlen mit ihrem Lauf (`slice-backfill-run-store`); der Lauf dieses
  Slice druckt `82.08% (gedeckt 843 von 1027 Statements …)` und `82.90%`/`83.00%`. Der Plan führt es als „gemeldet, nicht
  nachgezogen“ (§3, Zeile „Nenner“), die Träger nennen ihren Lauf und bleiben als lauf-gebundener Beleg wahr; das
  §6-Risiko „DB-Adapter-Coverage bewegt Zähler und Nenner“ trägt bei Closure seinen Ausgang.
- **V-4 (INFO) — Aufnahme ohne Publication-Mitgliedschaft nicht auf Worker-Ebene real getestet.** Der DoD-Satz „fehlen Bindung oder
  Publication-Mitgliedschaft bei der Aufnahme, endet der Run `failed` (Klasse `configuration`) ohne Slot“ ist am Worker gegen die
  reale Run-Tabelle nur für die Bindung belegt (`TestBackfillWorkerEndsARunOfAnUnboundTableAsConfigurationFailure`); die
  Publication-Mitgliedschaft prüft der Use Case (`usecase/backfill/service_test.go`, netzlos). Beide Zweige liegen in derselben
  Vorbedingungs-Prüfung des Use Cases; keine Lücke im Verhalten, eine Lücke in der Tiefe des Beleges.
- **V-5 (INFO) — Tier-Nebenwirkungen.** `make test-store` überschreibt `tools/schema/plan.yaml` (Zielzeile) und lässt je Lauf ein
  anonymes Volume zurück; beides habe ich nach jedem Lauf zurückgenommen bzw. gezielt entfernt (§1). Eigenschaft der Tier-Pfade,
  nicht dieses Slice.

## 11. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (8) | **8 von 8 erfüllt**, je mit eigenem Beleg (Zeile 6 mit V-2, Zeile 7 mit V-1) |
| DoD §2 — übrige `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Closure-/Rollen-Sequenz; Reconciliation „entfällt“) |
| Review-Findings F-1 bis F-6 (HIGH/MEDIUM), F-10 | **alle behoben**, je nachgemessen (Mutation rot bzw. Träger-Zahl gedruckt); F-6 mit Rest V-1 (LOW) |
| Plan-vs-Code-Diff | **deckungsgleich**; ungeplant nichts außer Plan und Review-Report |
| Alt-Tag-/Guard-Lauf | **EXIT=0**, Lauf 5 `v0.1.2` Exit 0/0/0, 19 Rechte, `EXECUTE` allein für `cdc_admin`, fünf `request_kind`-Werte; unter der Grant-Mutation EXIT=1 an Lauf 5 |
| Store-Tier | `make test-store` **EXIT=0** (zweimal), `82.08% (843 von 1027)`; `plan.yaml`/`down.sql` gleich dem frischen Rollout bis auf die `target`-Zeile |
| Entscheidungen | [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md), [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md), [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md), [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md), [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md), [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md), [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md), [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) **konform** |
| Harte Regeln | **erfüllt** (§9) |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=2d47d8a7..HEAD`) |
| Gates | **`make gates` EXIT=0**, `make test` EXIT=0, `make a-check` EXIT=0, `make coverage-gate` EXIT=0 (`83.00%`; im `make gates` `82.90%`) |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung belegt: `make gates`, `make test`, `make test-store` und
der Guard-Test laufen am Stand `c092efa5` mit Exit 0; der Alt-Tag-Lauf von `v0.1.2` über den Arbeitsbaum endet dreimal mit Exit 0
und belegt jetzt im Repo, was `harness/targets/schema-rollout.md` zusagt (Rechte der drei Rollen, `EXECUTE` allein für `cdc_admin`,
fünf `request_kind`-Werte); die vier Eingabeseiten-Mutationen der Fixrunde (Grant der Funktion, Quellfilter von `diagnose`,
Erfolgsmeldung von `MarkApplied`, Quellvergleich des Antragszweigs) färben ihre Tests rot; die Suchlauf-Zahlen des Plans stimmen
an beiden Ständen. Kein Befund blockiert; V-1 (LOW) ist ein Zählwort-Rest in einem Kommentar, V-2 bis V-5 sind INFO.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: Closure-Notiz mit Lerneintrag (die Finding-Klassen des Reviews
gehen in den Zähler, darunter „Nachzug lässt überholten Text stehen“ mit V-1 als weiterem Fall der Klasse „Suchmuster deckt das
Zählwort nicht“), Ausgänge der neun §6-Risiken (darunter der DB-Adapter-Nenner mit diesem Lauf, V-3), Beobachtungs-Register,
Reconciliation-Zeile, Welle-Paarungen. Danach der reine `git mv` nach `done/` ([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2).
