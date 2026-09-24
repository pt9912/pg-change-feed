# Verifikations-Report: slice-backfill-e2e — 2026-09-24

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates + realer Post-Push-Lauf ([`AGENTS.md`](../../AGENTS.md) §3.10).
Review-Artefakt des Reviewers: [`review-slice-backfill-e2e.md`](review-slice-backfill-e2e.md);
Formvorbild dieses Reports: [`verifikation-slice-backfill-sql-administration.md`](verifikation-slice-backfill-sql-administration.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`); der
Report folgt dem Formvorbild.

**Gegenstand:** Slice-Plan `slice-backfill-e2e` (Welle `welle-backfill-bestand`), Diff-Range
`e7df5619..HEAD` (`43137ebf`, gepusht), 18 Commits, 23 Dateien (+2910/−370): Lifecycle
(`ed0f99a8`, `ef9c3844`, `4a315549`), E2E-Runner und Go-Test (`6045bb09`, `baa75c96`, `3d756a76`),
Plan- und Träger-Nachzüge (`5ff1b0de`, `10f73f04`, `199bde89`, `c980d13b`, `26ea0c2e`, `54d715f9`,
`360ccb49`), Architect-Verdikt und Entscheidung (`ceea0af3`), Umsetzung der Lesesperre
(`530f8280`), Review (`31946550`), Fixrunde (`4b7e672e`, `43137ebf`). Dieser Lauf ändert weder
Code noch Plan, Spec oder Doku; er schreibt nur diesen Report. Alle Mutationen liefen im
Arbeitsbaum und sind je Lauf per `git checkout` zurückgenommen (`git status --short` leer nach
jedem Lauf); `tools/schema/plan.yaml`/`down.sql` (Nebeneffekt der Tier-Läufe) und
`docs/user/e2e-abdeckung.md` (Nebeneffekt des abgebrochenen E2E-Laufs) sind ebenfalls
zurückgenommen; das lokale `:dev`-Image ist nach dem Mutations-Lauf aus dem sauberen Baum neu
gebaut.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert
(`make … > <log> 2>&1; echo EXIT=$?`), die Logs danach gelesen. Stand aller Läufe (ohne die
Mutationen, §4): `HEAD` = `43137ebf`, Arbeitsbaum sauber.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `db-package-lists-check: OK — … dieselben 4 Pakete` · `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%` · `d-check: 1077 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make test` (Race-Detector) | **EXIT=0** | 42 Pakete `ok`, 0 `FAIL` |
| `make test-replication` (PostgreSQL 18, Repo-Default) | **EXIT=0** (68 s) | `ok … postgressnapshot 12.603s`; `DB-Adapter-Coverage: 82.13% (gedeckt 850 von 1035 Statements; Profile gemergt: store,replication)`; `db-coverage: OK — DB-Adapter-Coverage 82.13% erfuellt Schwelle 70%`; Slot-Reserve-Test `--- PASS: TestSlotReserveExhaustedIsConfiguration` |
| `make doc-trace` | **EXIT=0** | `80 Anforderung(en), 2 Waise(n).`; `LH-FA-CAP-009` … Nachweis `E2E` … `ok`; Waisen `LH-FA-CFG-007`, `LH-FA-CFG-008` |
| `make commit-traceability RANGE=e7df5619..HEAD` | **EXIT=0** | `OK — 18 Commit(s) in "e7df5619..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=e7df5619..HEAD` | **EXIT=0** | `d-check: 1077 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=e7df5619..HEAD` | **EXIT=0** | `d-check: 1077 Datei(en) geprüft, 0 Befund(e)` |
| `make image` + `make test-integration` mit Mutation (§4, M-E2E) | `image=0`, **EXIT=2** an der Zielstelle | `Backfill-DDL-Fenster feed_e2e_backfill_rewrite — Run … endete completed statt failed`; alle vorangehenden Phasen (Happy Path, Schema-Version, Startposition, Boundary, Replay-Invariante) druckten `belegt` |

**Nicht selbst als Grün-Lauf gefahren:** `make test-integration` am unveränderten Stand. Den
Grün-Beleg trägt der reale Post-Push-Lauf (§3): beide Legs, mit gedruckten Lauf-Zeilen. Ich habe
den schweren Lauf nicht zusätzlich lokal wiederholt (Auftragsgrenze: nur ohne CI-Beleg); der eine
lokale E2E-Lauf ist der Mutations-Lauf.

Hygiene: dangling-Volumes (`docker volume ls -q -f dangling=true | wc -l`) vor den Läufen **34**,
nach allen Läufen **34**; eigene Wegwerf-Container (`vf-pg`, Messung der Sperr-Warteschlange) mit
`docker rm -fv` entfernt, kein Container nach den Läufen; kein `prune`. Speicher vor den schweren
Läufen (`free -m`, verfügbar): 20,1 bis 20,8 GB.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Happy Path und Boundary am laufenden Feed-Container (Bestand als `backfill` über beide Lesewege, unterscheidbar von WAL, Replay-Invariante, leere Tabelle, zweiter Antrag; Lauf mit gedruckten `change_id`s) | **erfüllt, mit V-1 (LOW)** | Post-Push-Lauf (§3), beide Legs: Zeile `Backfill-Happy-Path … belegt — Run … übernahm 5 Bestandszeilen … lesbar über cdc.changes und GET /changes als INSERT mit origin backfill, change_id … gehalten (0bf-…-00000001-1,…)`, sechs `READ`-Zeilen (fünf `origin=backfill`, eine `origin=wal` mit `operation=UPDATE`); `Backfill-Boundary … belegt` (leere Tabelle `completed` mit 0 Zeilen ohne Transaktion, zweiter Antrag `failed` „Für die Tabelle besteht bereits ein aktiver Backfill-Run“); `--- PASS: TestE2EBackfillReplayInvariant`, Zeile „Replay-Invariante: Snapshot-Position 28721360, 54 Backfill-Changes, WAL-Changes davor 409 und dahinter 164, 57 Zeilen im Quellstand“ (PostgreSQL 17) bzw. „31875224 … 169 …“ (PostgreSQL 18). Ort der Replay-Invariante: die Bestätigung des Architects, die der Plan verlangt, steht in keinem committeten Artefakt (V-1) |
| 2 | Negative (`docker kill`, `interrupted`, keine sichtbare Change, erneuter Antrag `completed` einmal und vollständig; `queued`-Zeile überlebt); je Kriterium eine Mutation im Bericht | **erfüllt, mit V-2 (LOW)** | Post-Push-Lauf, beide Legs: `Backfill-Negative (docker kill, queued-Aufnahme) belegt — docker kill im zweiten Block von Run … (rows_copied 1000, running): keine Sitzung und kein Slot blieben zurück, keine Change des Runs war sichtbar; nach dem Neustart steht der Run interrupted, der queued wartende Run … wurde ohne neuen Antrag aufgenommen und endete completed (4 Zeilen); der erneute Antrag … übernahm 2500 Zeilen einmal (3 Blöcke), die Erfassung setzte nach dem Neustart fort`. Die Mutationen je Kriterium stehen in keinem committeten Artefakt und liegen für die Negative-Phase nicht in meinem Lauf (V-2) |
| 3 | Startposition gemessen und im Handbuch mit Lauf; Runner deklariert Phasen mit `LH-FA-CAP-009`; `docs/user/e2e-abdeckung.md` trägt Zeilen; `make doc-trace` ohne Waise `LH-FA-CAP-009`; jede neue `TestE2E*` läuft (`-v`) | **erfüllt** | Handbuch §4, Punkt „Startposition eines neuen Consumers“: `offset` 0, `acknowledged` `false`, 5 Backfill-Änderungen, 0 hinter der bestätigten Position, Ursprung „Lauf von `make test-integration`, Phase Backfill-Startposition“ — deckt sich mit der gedruckten Zeile beider CI-Legs („startet an Position 0 (acknowledged=false; … NULL|NULL …), die Snapshot-Position des Bestands ist 27961696: ab der Anfangsposition sind 5 Backfill-Changes lesbar; nach der Bestätigung von 28149704 liegen 0 hinter ihr“). `abdeckung_declare "Backfill…"` 6× (gezählt), `git grep`-Zeilen mit `LH-FA-CAP-009` in `docs/user/e2e-abdeckung.md`: Parent 0, HEAD 7 (die sechs Phasen plus `TestE2EBackfillReplayInvariant`); Runner meldet in beiden Legs „E2E-Abdeckungstabelle unverändert“ und „aus 14 Go-Zeilen und 34 Bash-Zeilen“; `make doc-trace` §1: 80/2, `LH-FA-CAP-009` trägt `E2E`; `=== RUN TestE2EBackfillReplayInvariant` … `--- PASS` in der `-v`-Ausgabe beider Legs. Die Ort-Anker der Abdeckungstabelle (2919, 3141, 3249, 270) zeigen auf die `belegt`-Zeile bzw. die Testfunktion des Standes; die Fixrunde ändert danach nur 1:1 Zeilen (`git diff 3d756a76..HEAD --stat`: 5 Einfügungen/5 Löschungen in Runner und Test) |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 5 | Umschreiben im Snapshot-Fenster ([`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md)): Lesesperre und Filenode-Vergleich, Klasse `transient`, neuer Antrag `completed` mit allen Zeilen; Belege in `make test`, `make test-replication`, Phase DDL-Fenster; je Zusage eine Eingabeseiten-Mutation | **erfüllt** | `make test` und `make test-replication` EXIT=0 (§1); Phase DDL-Fenster, beide CI-Legs: „ALTER COLUMN TYPE endete Run … failed (transient: Quelle für den Tabellen-Snapshot vorübergehend nicht verfügbar: Tabelle public.feed_e2e_backfill_rewrite wurde nach dem Snapshot-Export umgeschrieben; ein neuer Antrag beginnt neu) ohne Change …, der neue Antrag … übernahm 3 Zeilen“, davor `DROP COLUMN` `failed (storage: … SQLSTATE 42703)`; Mutationen der Eingabeseite in jedem der drei Tiers rot (§4, M1–M3, M5, M-E2E) |
| 6 | Review durchgeführt, kein offenes HIGH/MEDIUM | **erfüllt, mit V-3 (INFO)** | Report `review-slice-backfill-e2e.md` liegt vor (0 HIGH, 1 MEDIUM, 5 LOW, 4 INFO); F-1 und F-2 bis F-8 in der Fixrunde behoben, je nachgemessen (§5); kein zweiter Review-Report zur Fixrunde (V-3) |
| 7 | §3.13-Suchlauf: Feld in §3, Gefundenes und Nichtgefundenes, beide Stände gemessen | **erfüllt** | acht Zeilen des Feldes an beiden Ständen nachgefahren (§6): alle gedruckten Zahlen stimmen |
| 8 | Doku-Update (Handbuch §4 mit Startposition und Lauf-Ursprung, Historienzeile) | **erfüllt** | `Version: 1.52`, Historienzeilen 1.50 bis 1.52 (`docs/user/benutzerhandbuch.md`), Träger-Nachzüge §7 |
| 9 | Closure-Notiz mit Lerneintrag | **korrekt offen** | §7 des Plans trägt „*(zu tragen bei Closure)*“ |
| 10 | Reconciliation-Register — entfällt | **korrekt offen / entfällt** | Greenfield |
| 11 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht; im Diff bereits `observations/BEO-PGC/lesesperre-ohne-zeitgrenze/` (`observation.md`, `state.md`, `evidence/slice-backfill-e2e.md`) |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | die Zeilen „Ausgang: *(bei Closure)*“ bzw. „weiter offen“/„eingetreten → behoben“ sind Sache der Closure; der Ausgang des Runner-Zeit-Risikos ist mit §3 belegt |
| 13 | Die drei Paarungen | **korrekt offen** | hängen an der Closure von `welle-backfill-bestand` |

`[x]` sind acht Zeilen (Nr. 1–8), `[ ]` fünf (Nr. 9–13), zusammen 13 — gezählt am Plan (§2 hat 13
Aufzählungspunkte). Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt
wäre.

## 3. Post-Push-Lauf von `e2e.yml` für `43137ebf` (offenes Risiko „Zeitannahmen der Haltepunkte“, [`AGENTS.md`](../../AGENTS.md) §3.10)

Lauf `36065957210` (Event `push`, Kopf `43137ebf8710981adada5a76012193b88e70996d`, gestartet
2026-09-24T22:10:55Z, beendet 22:22:19Z) — `gh run list --commit 43137ebf… --json …` und
`gh api repos/…/actions/runs/36065957210/jobs`, Log über `gh run view 36065957210 --log`
(Exit 0, 1804 Zeilen). **Gesamtergebnis `success`; die Workflows `ci` und `examples` desselben
Commits ebenfalls `success`.**

| Leg | Schritt | Ergebnis | Dauer (Zeitstempel der API) |
|---|---|---|---|
| PostgreSQL 17 | Image bauen | success | 22:11:00 → 22:11:35 |
| PostgreSQL 17 | Compose-Integrationstest (Black-Box-E2E) | **success** | 22:11:35 → 22:18:50 (7 min 15 s) |
| PostgreSQL 17 | DB-Adapter-Coverage Store-Teil / Replication-Teil, Merge + Schwelle | success / success | „82.13% (gedeckt 850 von 1035 Statements; Profile gemergt: store,replication)“, „db-coverage: OK“ |
| PostgreSQL 17 | Replication-Tier (`go test ./...` und Slot-Reserve) | **success** | 22:20:50 → 22:22:09 |
| PostgreSQL 18 | Image bauen | success | 22:11:01 → 22:11:37 |
| PostgreSQL 18 | Compose-Integrationstest (Black-Box-E2E) | **success** | 22:11:37 → 22:18:57 (7 min 20 s) |
| PostgreSQL 18 | DB-Adapter-Coverage Store-Teil / Replication-Teil, Merge + Schwelle | success / success | „82.13% (gedeckt 850 von 1035 Statements …)“ — identisch zu meinem lokalen Lauf (§1) |
| PostgreSQL 18 | Replication-Tier | **success** | 22:20:58 → 22:22:17 |

Gedruckte Haltepunkt-Zeiten der Pause des Feed-Containers (Grenze im Runner: unter 1000 ms, die
Hälfte von `wal_sender_timeout` = 2 s): PostgreSQL 17 — DROP-COLUMN-Lauf **110 ms**,
Umschreib-Lauf **163 ms**; PostgreSQL 18 — **245 ms** und **173 ms**. Keine `FAIL`-Zeile im
Log (`grep -c FAIL` = 0). Die Läufe der beiden Legs ergeben je sechs `belegt`-Zeilen der
Backfill-Runner-Phasen, eine Zeile `Replay-Invariante` und „E2E-Abdeckungstabelle unverändert —
`docs/user/e2e-abdeckung.md` entspricht dem Quelltext-Stand“.

**Befund zum Risiko:** die Zeitannahmen der Haltepunkte tragen auf dem GitHub-Runner in beiden
Legs (Pause 110–245 ms gegen die Grenze 1000 ms; Laufzeit des E2E-Schritts 7 min 15 s bzw.
7 min 20 s, `timeout-minutes: 60`). Aussagegrenze: **ein** Lauf je Leg — er belegt „hält auf dem
Runner“, nicht „flackert nie“. Ausgang für die Closure: *entfallen* mit diesem Lauf als Beleg
(Lauf-Nummer und Kopf-SHA oben); tritt später ein Timeout auf, öffnet er das Risiko erneut. Die
Laufzeit-Zahl des Plans („lokal 4 min 39 s“, Ursprung: der Lauf der Fixrunde) ist von mir nicht
nachgemessen (übernommen); die Zahl des Reviews (277 s) und die des Runners (435 s/440 s) sind
verschiedene Umgebungen.

## 4. Mutationen der Eingabeseite (dieser Lauf)

Aufbau je Mutation: Datei im Arbeitsbaum mutiert (`sed` mit Ausgabe in eine Scratch-Datei und
`cp`, kein `sed -i`), Lauf ungefiltert in eine Log-Datei, gedruckte Rot-Meldung gelesen, Datei
per `git checkout` zurückgenommen.

| # | Mutation | Lauf | Rot gesehen (gedruckt) |
|---|---|---|---|
| M1 | `RewriteQuery`: `c.relfilenode <> coalesce(…)` → `c.relfilenode = coalesce(…)` (`snapshotlogic/logic.go`) | `make test-replication`, **EXIT=2** | `--- FAIL: TestRewriteInWindowIsTransient` (alle fünf Formen: `ALTER_COLUMN_TYPE`, `TRUNCATE`, `VACUUM_FULL`, `CLUSTER`, `DROP_und_CREATE`), `--- FAIL: TestNoRewriteInWindowReadsTheSnapshot` (`keine_DDL`, `ADD_COLUMN`, `DROP_COLUMN`) sowie die Lese-Tests des Pakets (die Umkehrung des Vergleichs endet jeden Import `transient`) |
| M2 | `CheckRewrite`: Klasse des Umschreibens `transient` → `storage` | `make test`, **EXIT=2** | `--- FAIL: TestCheckRewrite`, `FAIL … snapshotlogic` |
| M3 | Sperranweisung durch `SELECT 1` ersetzt (`lockAndVerify`) | `make test-replication`, **EXIT=2** | `--- FAIL: TestImportWaitsForExclusiveLockAndThenAborts (15.15s)`, `TestPermissionClassWithoutSelect`, `TestConfigurationClass` |
| M4 | `failureText` entfernt die Klassen-Angabe nicht mehr (`strings.Replace(…, "", 0)`) — der erste Versuch (Ersatz des Rückgabewerts) brach am ungenutzten Import ab und zählt nicht | `make test`, **EXIT=2** | `--- FAIL: TestExecuteFailureTextCarriesClassOnce` |
| M5 | Wartezeit der Sperranweisung vom Kontext des Aufrufs gelöst (`context.WithTimeout(context.Background(), 15*time.Second)`) | `make test-replication`, **EXIT=2** | `--- FAIL: TestImportLockWaitEndsWithTheContext (15.11s)` — „importSnapshot kehrt nach 15.017430783s zurück, erwartet nach dem Ablauf des Kontexts (2s)“ |
| M-E2E | `CheckRewrite`: bei Antwort `t` (umgeschrieben) `rewritten = false` gesetzt; Image aus dem mutierten Baum gebaut (`make image`, Exit 0) | `make test-integration`, **EXIT=2** nach 236 s | Phasen Happy Path, Schema-Version, Startposition, Boundary, Replay-Invariante `belegt`; danach „Backfill-DDL-Fenster feed_e2e_backfill_rewrite — Run 3f235a7e-… endete completed statt failed“ — die Zusage „Umschreiben endet `failed`/`transient`“ der Phase DDL-Fenster färbt an genau dieser Stelle rot |

Eine weitere Mutation blieb grün und ist **äquivalent**, keine Bindungslücke: `ClassifyLock` mit
`context.Background()` statt `ctx` aufgerufen ergibt weiter `transient`, weil `Classify` neben
`ctx.Err()` auch `errors.Is(cause, context.DeadlineExceeded)` prüft (`snapshotlogic/logic.go`
Kopfzeile von `Classify`). Die Kontext-Bindung trägt M5, nicht dieser Parameter (V-4).

## 5. Fixrunde: Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) Handbuch nennt nur Leser | **Sperr-Warteschlange selbst nachgemessen** (Wegwerf-Container `postgres:18-alpine` mit dem Pin von `PG_TEST_IMAGE`, Sitzung A `LOCK TABLE … IN ACCESS SHARE MODE` mit `pg_sleep(20)`, Sitzung B `ALTER TABLE … ALTER COLUMN name TYPE varchar(64)` wartet): `INSERT INTO t …` → `exit=124 after 4005 ms`, `SELECT count(*) FROM t` → `exit=124 after 4005 ms`, `SELECT count(*) FROM pg_publication_tables WHERE pubname='pub'` → `exit=124 after 4005 ms`; `pg_stat_activity` zeigt alle vier Sitzungen mit `wait_event_type='Lock'`. Handbuch `Version: 1.52`, Absatz „Sperre der Tabelle“ nennt „Schreiber und Leser der Tabelle ebenso wie die Abfrage der Publication-Mitgliedschaft (`pg_publication_tables`)“ mit Messung (PostgreSQL 18, Ursprung Review F-1) und die Gegenrichtung; Historienzeile 1.52 vorhanden | **behoben**, Messung bestätigt |
| F-2 (LOW) Architektur-Sicht | `git diff 31946550..HEAD -- spec/architecture.md`: Diagrammzeile „Lesesperre auf die Tabelle, Umschreib-Prüfung“ und Absatz zu Sperrdauer, wartender DDL, Run-Ende `failed`/`transient` ohne Change; Diff-Zeilen mit `ADR-`, `slice`, `welle`, `SPEC-0`, `ARC-` (`grep -i`): **0 Treffer** ([`AGENTS.md`](../../AGENTS.md) §3.4) | **behoben** |
| F-3 (LOW) Zahl im Suchlauf | Zeile „Beschreibungen von `make test-integration` außerhalb der README“: Parent `e7df5619` 124 Zeilen/57 Dateien, HEAD **130 Zeilen/60 Dateien** — Plan nennt 130/60 | **behoben** |
| F-4 (LOW) Heartbeat-Assertion | Diff gelesen: `SELECT count(*) FROM cdc.heartbeat WHERE source_id = 'src-e2e' AND error_class IS NULL` gleich 1 — bei fehlender Zeile 0, färbt rot; beide CI-Legs grün mit der neuen Form; **nicht** durch einen E2E-Lauf gemutet (Grenze §8) | **behoben** (Lesung + grüner Lauf) |
| F-5 (LOW) zwei Kommentare | Diff gelesen: Godoc von `TestE2EBackfillReplayInvariant` endet „Der Test bricht ab, wenn eine der beiden Seiten leer bleibt.“; Runner-Kommentar „(die Slot-Anlage wartet auf jede offene Schreibtransaktion)“ im Indikativ; je 1:1 Zeilen, keine Anker-Verschiebung | **behoben** |
| F-6 (LOW) Wartegrenze | Handbuch-Absatz „In der Gegenrichtung wartet der Run …“ vorhanden; `TestImportLockWaitEndsWithTheContext` an seiner Eingabeseite gebunden (M5 rot, §4); Register-Eintrag `BEO-PGC/lesesperre-ohne-zeitgrenze` mit Ausgang „weiter offen“ angelegt; Plan §6 trägt den Punkt | **behoben** (Verhalten bleibt als Designfrage offen, benannt) |
| F-7 (INFO) `RENAME COLUMN` | Handbuch kennzeichnet: „gemessen im Review-Report, PostgreSQL 18, Klasse `storage`; kein Beleg im E2E-Runner“; die ADR ist `Accepted` und unberührt | **behoben im Träger**, ADR-Satz „E2E-belegt“ bleibt (V-5) |
| F-8 (INFO) Fehlertext | M4 rot (`TestExecuteFailureTextCarriesClassOnce`); gedruckte Zeile beider CI-Legs: „failed (transient: Quelle für den Tabellen-Snapshot vorübergehend nicht verfügbar: …)“ — die Klasse steht **einmal** | **behoben** |
| F-9 (INFO) Runner-Zeit | §3 | **belegt** |
| F-10 (INFO) `gofmt` | Bestand außerhalb des Diffs | unverändert, kein Anspruch dieses Slice |

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent per `git grep … e7df5619` bzw. `31946550`/`ceea0af3` (die Stände, die der Plan je Zelle
nennt), Diff-Stand per `git grep … HEAD` (Arbeitsbaum = HEAD). Gedruckte Zahlen:

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| Abdeckungs-Tabelle (`^\| \[`-Zeilen; mit `LH-FA-CAP-009`) | Parent **41 / 0**, HEAD **48 / 7** | bestätigt |
| `make test-integration` außerhalb der README (Zeilen/Dateien) | Parent **124 / 57**, HEAD **130 / 60** | bestätigt |
| Zählwörter (Rundläufe, Belege, Phasen, Haltepunkte, Fremdobjekte, Waisen) | Parent **11**, HEAD **11** | bestätigt |
| Träger der Fixrunde (Muster zu Sperre und Umschreiben) | `ceea0af3` **4**, `31946550` **37**, HEAD **44** | bestätigt |
| Sequenzdarstellung in `spec/architecture.md` | `31946550` **9**, HEAD **16** | bestätigt |
| Wirkung der Tabellensperre auf andere Zugriffe | `31946550`: 3 Treffer im Handbuch (Zeilen 456, 458, 518), HEAD 2 (458, 534) | bestätigt |
| `error_message` (Docs, Spec, Harness, README) | `31946550` **15**, HEAD **15** | bestätigt |
| Symbolnamen des Imports (`importSnapshot`, `readColumns`, `snapshotlogic`) außerhalb `*.go` | beide Stände dieselben drei Dateien (`harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md`, `tools/harness/db-coverage.sh`) | bestätigt |
| RTM-Träger | `make doc-trace` §1: 80/2, `LH-FA-CAP-009` `E2E` | bestätigt |
| CI-Träger | `git diff --stat e7df5619 -- .github` leer; `timeout-minutes: 60`, Matrix 17/18 (§3 bestätigt beide Legs) | bestätigt |

Nicht-Fund-Suche eigener Art: `git grep` nach `Lesesperre|Umschreib|umgeschrieben|Tabellensperre|
Sperre der Tabelle|Export und …` über HEAD ohne Reviews, ADRs, `done/`, Beobachtungen, Tests, Plan
dieses Slice: Treffer in `docs/plan/planning/welle-backfill-bestand.md` (1, anderer Gegenstand, im
Plan benannt), `docs/user/benutzerhandbuch.md` (9), `docs/user/e2e-abdeckung.md` (1),
`harness/README.md` (2), `spec/architecture.md` (6), `spec/pflichtenheft.md` (3), Produktions- und
Runner-Dateien — kein weiterer beschreibender Träger, der die Sperre oder das Umschreiben falsch
oder gar nicht führt.

## 7. Träger-Nachzüge

| Träger | Beleg | Verdikt |
|---|---|---|
| Handbuch `Version: 1.52`, Historienzeilen 1.50, 1.51, 1.52 | `docs/user/benutzerhandbuch.md` Zeilen 3 und 1605–1607 (Datei-Ende); Abschnitt „Bestand als Backfill überführen“ trägt Absatz „Sperre der Tabelle“, die Gegenrichtung, den Punkt „Umschreiben der Tabelle im Fenster“ und die Startposition mit Lauf-Ursprung | erfüllt |
| [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) Absatz „Mechanismus“ | Satz „Der Import steht unter einer Lesesperre der Tabelle; wurde die Tabelle zwischen Snapshot-Export und Sperre umgeschrieben, endet der Run `failed` (Fehlerklasse `transient`) ohne Change, ein neuer Antrag beginnt neu.“ — deckt sich mit dem Wortlaut-Vorschlag der ADR-Folgepflicht 2; Änderungshistorie um eine Zeile ergänzt; **kein** ADR-/Slice-/Wellen-Bezug in den hinzugefügten Zeilen (`grep`) | erfüllt |
| `spec/architecture.md` (Sequenz der Run-Ausführung) | siehe §5 F-2; ohne ADR-/Slice-Bezug | erfüllt |
| `harness/README.md`, Zeile `make test-integration` | „Zusätzlich sieben Backfill-Rundläufe …“ — nachgezählt: 6 `abdeckung_declare "Backfill…"` + 1 `func TestE2EBackfill…` = 7; drei Haltepunkte (offene Schreibtransaktion, unbestätigter Schlüssel in `cdc.transaction`, `docker pause`); Phase DDL-Fenster mit `storage` und `transient`; Lauf-Aussagen sind Beschreibung, keine Zahl ohne Lauf | erfüllt |
| `harness/README.md`, Zeile `make test-replication` | „… die Cursor-Blöcke, die Lesesperre und den Filenode-Vergleich bei einem Umschreiben zwischen Export und Import ([`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md)) …“ — deckt sich mit den Tests, die mein Lauf (§1) fährt | erfüllt |
| `harness/README.md`, Zeile `make doc-trace` | „gemessen 2026-09-24 … 80 Anforderungen, **2 Waisen** — `LH-FA-CFG-007`/`LH-FA-CFG-008` …; `LH-FA-CAP-009` trägt den Nachweis `E2E`“ — gegen `make doc-trace` §1 (80/2, dieselben zwei Waisen) | erfüllt |
| `docs/user/e2e-abdeckung.md` | Erzeugnis, sieben Zeilen mit `LH-FA-CAP-009` (Parent 0), Runner meldet „unverändert“ (§3, beide Legs) | erfüllt |
| Register `BEO-PGC/lesesperre-ohne-zeitgrenze` | `observation.md`, `state.md` („offen — Ausgang weiter offen“, Zähler abgeleitet 1×), `evidence/slice-backfill-e2e.md` | erfüllt |

## 8. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Festlegung 1 (Lesesperre zuerst) | `LOCK TABLE <schema>.<table> IN ACCESS SHARE MODE` mit gequoteten Bezeichnern nach dem Import und dem Ende der Replication-Verbindung; Wartezeit trägt der Kontext des Aufrufs | `importSnapshot` → `lockAndVerify` (Import → `closeConn(export.conn)` → Sperre → Vergleich → `readColumns` → `DECLARE`); `LockStatement` mit `QuoteIdent`; M3 (Sperre entfernt: rot), M5 (Kontext-Bindung: rot); `42501` → `permission` durch `TestPermissionClassWithoutSelect` | konform |
| [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Festlegung 2 (Filenode-Vergleich) | Abfrage `SELECT c.relfilenode <> coalesce(pg_relation_filenode(c.oid), 0) … to_regclass(quote_ident($1) \|\| '.' \|\| quote_ident($2))`; Unterschied und keine Zeile enden den Run, `false` setzt fort | `snapshotlogic.RewriteQuery` (Wortlaut identisch), `CheckRewrite` (`t`/keine Zeile → `transient`, `f` → weiter, sonst `storage`); M1 (Operator), M-E2E (Zweig `t`) rot; partitionierte Tabelle: `TestPartitionedTableIsNoFalseAlarm` grün (Review-Mutation M2 des Reviews, von mir nicht wiederholt) | konform |
| [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Festlegung 3 (Klasse `transient`, run-lokal) | `outbound.ErrSnapshotTransient`, Verbindung geschlossen (`snap.abort`), keine Change, Run `failed`, Erfassungspfad unberührt | M2 rot; Phase DDL-Fenster („transient: …“, „ohne Change“, „der Erfassungspfad lief in beiden Fällen weiter“), Run-Fehler ändert nicht den Lebenszeichen-Fehlerzustand (Assertion `error_class IS NULL` gleich 1, §5 F-4) | konform |
| [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Festlegung 4 (Fehlalarme akzeptiert) und 5 (`DROP COLUMN` endet `storage`) | `VACUUM FULL`/`CLUSTER`/`TRUNCATE` enden `transient`; `DROP COLUMN` `storage` (SQLSTATE 42703) | `TestRewriteInWindowIsTransient` (fünf Formen), `TestNoRewriteInWindowReadsTheSnapshot/DROP_COLUMN`; Phase DDL-Fenster („storage: … SQLSTATE 42703“) | konform; der Satz „E2E-belegt“ zu `RENAME COLUMN` in Festlegung 5 hat im Repo keinen Beleg (V-5) |
| [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Folgepflichten 1–4 | Umsetzung, Spec-Satz, Handbuch, Träger nachziehen | §5, §6, §7 | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 1 (Startposition messen), Teilfrage 3/4/5, Folgepflicht 2 (`S5`) | Startposition gemessen und dokumentiert; Replay-Invariante belegt; `running` → `interrupted` beim Start, `interrupted` startet nicht von selbst; Run-Fehler run-lokal | Handbuch-Zeile und CI-Zeile (§2 Nr. 3); `TestE2EBackfillReplayInvariant`; Phase Negative („nach dem Neustart steht der Run interrupted“, erneuter Antrag eigener Run); Ort der Replay-Invariante (E2E statt Tier `make test-replication`): Bestätigung des Architects nicht auffindbar (V-1) | konform, mit V-1 |
| [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2 | `queued`-Zeile überlebt den Neustart und wird beim Prozessstart aufgenommen | Phase Negative, beide Legs: „der queued wartende Run … wurde ohne neuen Antrag aufgenommen und endete completed (4 Zeilen)“ | konform |
| [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) | Bild-Parität von WAL- und Backfill-Pfad über den Typ-Satz | `TestImageParityWalAndBackfill` im `make test-replication`-Lauf (§1, EXIT=0) und im Replication-Tier des CI-Laufs (§3); Happy-Path-Zeilen `new_image={"id":"3","name":"BackfillGamma","note":"n3"}` | konform |
| [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md) | Kennung der Backfill-Changes = zum Antrag aktuelle Version; Bild trägt spätere Spalte | Phase Schema-Version, beide Legs („alle Backfill-Changes tragen die zum Antrag aktuelle Version …-v1, das Row Image trägt die Spalte extra“); Kommentar an `currentVersion` ohne Verhaltens-Diff | konform |
| [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) | Klasse `schema` der Nichtanwendbarkeit einer Regel vorbehalten; keine neue Klasse | Umschreiben endet `transient`, `DROP COLUMN` `storage`; `classifyError` im Diff unberührt (nur `failureText` neu, `git diff … service.go`) | konform |
| [`ADR-0012`](../plan/adr/0012-at-least-once.md), [`ADR-0030`](../plan/adr/0030-testpyramide.md) | at-least-once / Idempotenz; E2E-Tier ausschließlich über externe Wege | `backfill_e2e_test.go` importiert nichts aus `internal/**`; Replay-Anwendung als idempotenter Upsert/Delete | konform |
| [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) | Anweisung und Entscheidung netzlos in `snapshotlogic` | `snapshotlogic` im Unit-Gegenstand; `make coverage-gate` 83.10% (in `make gates`), DB-Adapter-Coverage 82.13% (850/1035) | konform |

## 9. Plan-vs-Code-Diff

Jede Zeile der Plan-Tabelle §3 ist im Diff (`git diff --stat e7df5619..HEAD`, 23 Dateien)
vertreten: `test/integration/backfill_e2e_test.go` (neu, 433 Zeilen; `integration_test.go`
unverändert), `tools/harness/run-integration-tests.sh`, `postgressnapshot/snapshot.go` und
`snapshotlogic/logic.go` mit ihren Tests, `usecase/backfill/service.go` samt Test,
`tools/harness/httpclient/main.go`, `spec/pflichtenheft.md`, `spec/architecture.md`,
`docs/user/benutzerhandbuch.md`, `docs/user/e2e-abdeckung.md`, `harness/README.md`, die ADR zum Umschreiben im Snapshot-Fenster samt
Index-Zeile, Architect-Verdikt, Review, Register-Eintrag. `compose.yaml` steht im Plan als
„prüfen, unverändert“ und ist nicht im Diff (`git diff --stat` ohne `compose.yaml`);
`.github/workflows/e2e.yml` nicht im Diff. Ungeplant im Diff: nichts außer den Lifecycle-Moves
und dem Review-Report. Die Lifecycle-/Move-Commits sind rein: `ed0f99a8` (`open/` → `next/`),
`4a315549` (`next/` → `in-progress/`) und `baa75c96` (`backfill_test.go` →
`backfill_e2e_test.go`) tragen im Diff jeweils keine Inhaltsänderung
([`AGENTS.md`](../../AGENTS.md) §3.3); `Verantwortlich` steht in `ef9c3844` als eigener Commit
zwischen den beiden Moves. Der Plan-Drift im Plan („sechs Fremdobjekte, `knownForeignObjects`
führt sieben“) ist im Plan selbst benannt.

## 10. Harte Regeln

- **§3.1** — alle Läufe über `make` und gepinnte Images; Wegwerf-Container aus dem `PG_TEST_IMAGE`-Pin per
  `docker run`/`docker exec`; keine Host-Toolchain.
- **§3.2** — `git diff e7df5619..HEAD` über `*.go`, `*.sh`, `*.sql`, `*.yaml` (`+`-Zeilen): kein `nolint`.
- **§3.3** — siehe §9.
- **§3.4** — Architektur-Diff ohne ADR-/Slice-/Wellen-Bezug (§5 F-2).
- **§3.5** — `make doc-immutable RANGE=e7df5619..HEAD` EXIT=0; die neue ADR ist als eigene Datei angelegt, keine bestehende `Accepted`-ADR im Diff.
- **§3.6** — keine Schwelle gesenkt: `Dockerfile`, `harness/mk/**`, `THRESHOLD` nicht im Diff.
- **§3.9** — Exit-Codes nie durch eine Pipe gemessen (§1); Gate-Lauf und Folgehandlung getrennt.
- **§3.10** — der Diff berührt keinen Workflow; der reale Post-Push-Lauf des bestehenden `e2e.yml`
  ist für das Risiko „Zeitannahmen der Haltepunkte“ trotzdem gefahren (§3).
- **§3.11** — `make docs-check` (in `make gates`) 0 Befunde; dieser Report trägt keinen host-lokalen Pfad.
- **§3.12** — jede Zahl dieses Reports ist **gemessen** (Lauf oder Befehl genannt) oder als
  **übernommen** gekennzeichnet (Laufzeit „lokal 4 min 39 s“ des Plans; Review-Mutation M2 zur
  partitionierten Tabelle).
- **§3.13** — Suchlauf-Feld an beiden Ständen nachgefahren (§6).
- **Handbuch-Pflicht** — `Version: 1.51` → `1.52` und Historienzeile in der Fixrunde.
- **Commit-Traceability** — 18 Commits, jede Message nennt eine Kennung, keine Struktur-Kennung im Betreff (§1).

## 11. Befunde

Kein V-Befund blockiert. V-1 und V-2 sind LOW, V-3 bis V-6 INFO.

- **V-1 (LOW) — Ortswahl der Replay-Invariante ohne auffindbare Bestätigung.** Der Plan (§2 Punkt 1)
  sagt: „der Architect bestätigt die Ortswahl im Review oder verlangt zusätzlich den Tier-Beleg“
  ([`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt in der
  Fitness-Function-Zeile den Tier `make test-replication`, in Folgepflicht 2 den E2E-Beleg). Weder
  `review-slice-backfill-e2e.md` noch `architect-verdict-backfill-tabellen-rewrite-im-fenster.md`
  noch `architect-verdict-backfill-schema-klasse-rollen.md` nennen die Ortswahl (`grep -i ortswahl` und
  `Replay-Invariante` in den drei Dateien: nur die gedruckte Lauf-Zeile im Review). Der Beleg selbst
  ist erbracht (beide CI-Legs); offen ist die Bestätigung des Orts. Aktion: Architect bestätigt bei
  Closure oder verlangt den Tier-Beleg; bis dahin bleibt die Plan-Zusage unbelegt.
- **V-2 (LOW) — Zusage „je Kriterium eine Mutation im Bericht“ ohne committeten Anker; E2E-Mutationen der Negative-Kriterien nicht nachgefahren.**
  Der Plan (§2 Punkt 2) verlangt für jedes Negative-Kriterium (Kill, `interrupted`, keine sichtbare
  Change, erneuter Antrag einmal, `queued`-Aufnahme) eine Mutation „im Bericht“. Der Bericht des
  Implementers liegt nicht im Repo; ein committetes Artefakt nennt sie nicht (`grep -i mutation` im
  Plan: nur die zwei Zusage-Sätze; Review und Fixrunde tragen die Mutationen M1–M8 für Import und
  Use Case, nicht für die Runner-Phasen). Von mir sind die Mutationen der Zusagen zum Umschreiben im Snapshot-Fenster in allen
  drei Tiers gefahren (§4, davon eine E2E-Mutation der Phase DDL-Fenster); für die Negative-Phase
  (`docker kill`, `queued`) habe ich keine E2E-Mutation gefahren (Lauf-Kosten je ~5 min, Grenze §12).
  Aktion: die Closure-Notiz nennt die Mutationen je Kriterium mit Lauf-Ursprung, oder der Planner
  entscheidet, dass die Runner-Assertions gelesen genügen.
- **V-3 (INFO) — Fixrunde ohne eigenen Review-Report.** Die DoD-Zeile „Review durchgeführt“ ist durch
  den Report des Reviewers erfüllt; die Fixrunde habe ich in §5 selbst gemessen (F-1 durch eigene
  Nachmessung der Sperr-Warteschlange, F-6/F-8 durch Mutation, F-2/F-3/F-7 an den Trägern). Eine
  Reviewer-Durchsicht der Fixrunde wäre ein eigener Schritt.
- **V-4 (INFO) — `ClassifyLock` trägt den `ctx`-Parameter ohne eigene Wirkung für den Kontext-Abbruch.**
  Die Mutation `context.Background()` an dieser Stelle ist äquivalent (§4), weil `Classify` die
  Ursache prüft. Kein Bindungsmangel: die Kontext-Zusage trägt die Anweisung (M5). Entwurfs-Beobachtung
  für den Reviewer, keine Aktion dieses Slice.
- **V-5 (INFO) — `Accepted`-ADR führt zwei Aussagen, die die Messungen dieses Slice überholen; die ADR bleibt unberührt.**
  [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Konsequenzen: „Leser der
  Quelltabelle … stauen sich hinter ihr (hergeleitet, nicht gemessen)“ — gemessen (F-1, §5) stauen
  sich auch Schreiber und die Publication-Abfrage der Administration; Festlegung 5: „`DROP COLUMN` und
  `RENAME COLUMN` … (E2E-belegt)“ — im Repo hat nur `DROP COLUMN` einen Beleg. Beide Träger, die
  Betreiber lesen (Handbuch), sind nachgezogen; der ADR-Text ändert sich nur durch eine Folge-ADR
  ([`AGENTS.md`](../../AGENTS.md) §3.5). Zur Kenntnis des Architects.
- **V-6 (INFO) — Tier-Nebenwirkungen und Hygiene.** `make test-replication` überschreibt
  `tools/schema/plan.yaml` (Zielzeile), der abgebrochene E2E-Lauf `docs/user/e2e-abdeckung.md`
  (Regeneration am mutierten Stand); beides habe ich nach jedem Lauf per `git checkout` zurückgenommen.
  Ein `make image` aus dem mutierten Baum hätte das lokale `:dev`-Image verfälscht; nach der
  Zurücknahme habe ich `make image` aus dem sauberen Baum erneut gebaut (Exit 0).

## 12. Grenzen dieses Laufs

- Kein lokaler Grün-Lauf von `make test-integration` am unveränderten Stand — der Beleg ist der
  reale Post-Push-Lauf (§3), beide PostgreSQL-Legs, mit gedruckten Zeilen.
- Keine E2E-Mutation der Negative-Phase und der Heartbeat-Assertion (V-2, §5 F-4): gelesen, nicht
  gemutet; die eine E2E-Mutation (M-E2E) bindet die Zusage der Phase DDL-Fenster.
- Tier-Läufe (`make test-replication`) nur gegen PostgreSQL 18 (Repo-Default); die Legs mit
  PostgreSQL 17 trägt der CI-Lauf (Replication-Tier `success`, §3), von mir nicht lokal gefahren.
- Übernommen, nicht selbst gefahren: die Review-Mutationen M2/M4–M8 des Reviews (Reviewer-Lauf, hier
  nur benannt) und die Laufzeit „lokal 4 min 39 s“ des Plans.

## 13. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (8) | **8 von 8 erfüllt**, je mit eigenem Beleg (Nr. 1 mit V-1, Nr. 2 mit V-2, Nr. 6 mit V-3) |
| DoD §2 — übrige `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Closure-/Rollen-Sequenz; Reconciliation „entfällt“) |
| Review-Findings F-1 (MEDIUM), F-2 bis F-8 | **alle behoben**, je nachgemessen (F-1 durch eigene Reproduktion, F-6/F-8 durch Mutation, übrige an den Trägern und im grünen Lauf); F-9 belegt (§3), F-10 Bestand |
| Post-Push-Lauf `e2e.yml` für `43137ebf` | **grün, beide Legs** (PostgreSQL 17 und 18; Lauf `36065957210`), Pause 110–245 ms |
| Plan-vs-Code-Diff | **deckungsgleich**; ungeplant nichts außer Lifecycle-Moves und Review-Report |
| Entscheidungen | [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md), [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md), [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md), [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md), [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md), [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md), [`ADR-0012`](../plan/adr/0012-at-least-once.md), [`ADR-0030`](../plan/adr/0030-testpyramide.md), [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) **konform** (Ort der Replay-Invariante: V-1) |
| Mutationen der Eingabeseite | **sechs rot** (M1–M5 und M-E2E), eine äquivalente grün (V-4) |
| Harte Regeln | **erfüllt** (§10) |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=e7df5619..HEAD`) |
| Gates | **`make gates` EXIT=0**, `make test` EXIT=0, `make test-replication` EXIT=0 |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung und den realen
Post-Push-Lauf belegt: `make gates`, `make test` und `make test-replication` laufen am Stand
`43137ebf` mit Exit 0; `e2e.yml` ist auf `43137ebf` in beiden PostgreSQL-Legs grün (die
Haltepunkt-Zeiten halten auf dem Runner: Pause 110–245 ms gegen 1000 ms), das Risiko „Zeitannahmen
der Haltepunkte“ kann bei Closure als *entfallen* mit diesem Lauf als Beleg geführt werden; jede
Zusage von [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) färbt bei einer
Mutation ihrer Eingabeseite rot, in allen drei Tiers einschließlich des E2E; die Findings des Reviews
sind behoben, das MEDIUM F-1 durch eigene Reproduktion der Sperr-Warteschlange bestätigt. Kein Befund
blockiert; V-1 und V-2 (LOW) sind offene Nachweise für die Closure, V-3 bis V-6 INFO.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: Closure-Notiz mit Lerneintrag (die
Finding-Klassen des Reviews gehen in den Zähler), Ausgänge der §6-Risiken (Runner-Zeit: *entfallen*,
Beleg §3; Wartegrenze der Sperre: *weiter offen* mit Register-Eintrag; „Nichtdeterministischer
Abbruch“: die Phase Negative lief in drei grünen Läufen — Review, CI PostgreSQL 17 und 18), Antwort
auf V-1 (Ort der Replay-Invariante, Architect) und V-2 (Mutationen je Negative-Kriterium),
Beobachtungs-Register, Reconciliation-Zeile, Welle-Paarungen. Danach der reine `git mv` nach
`done/` ([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2).
