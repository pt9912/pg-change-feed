# Verifikations-Report: slice-backfill-slot-leerlauf-bestaetigung — 2026-09-25

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates + Sicherheitsbeleg des Capture-kritischen Pfads. Review-Artefakt des
Reviewers: [`review-slice-backfill-slot-leerlauf-bestaetigung.md`](review-slice-backfill-slot-leerlauf-bestaetigung.md);
Formvorbild dieses Reports: [`verifikation-slice-backfill-e2e.md`](verifikation-slice-backfill-e2e.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`); der
Report folgt dem Formvorbild.

**Gegenstand:** Slice-Plan `slice-backfill-slot-leerlauf-bestaetigung` (Welle `welle-backfill-bestand`),
Stand `HEAD` = `80b451c4` (nicht gepusht; `origin/main` steht auf `427f6d1b`), Diff-Range
`f4e32fba..HEAD`: 13 Commits, 26 Dateien (+2188/−550). Davon Slice-Inhalt: Lifecycle-Moves und
Verantwortlich (`1f75a2be`, `52b6d6f0`, `5e84274a`), Umsetzung (`2dabfac3`), Belege in den drei Tiers
(`00caec48`), Träger-Nachzug (`253986e9`), Plan-Nachträge (`3a876362`, `427f6d1b`), Fixrunde nach dem
Review (`ec7e43dd`, `5a5d3422`, `77b9a20f`, `80b451c4`); nicht Slice-Inhalt ist der Review-Report
(`67e331ac`). Dieser Lauf ändert weder Code noch Plan, Spec oder Doku; er schreibt nur diesen Report.
Alle Mutationen liefen im Arbeitsbaum und sind je Lauf per `git checkout` zurückgenommen
(`git status --short` leer nach jedem Lauf); `tools/schema/plan.yaml`/`down.sql` (Nebeneffekt der
Tier-Läufe und des Bench) sind ebenfalls zurückgenommen.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert
(`make … > <log> 2>&1; echo $?`), die Logs danach gelesen. Stand aller Läufe ohne die Mutationen (§4):
`HEAD` = `80b451c4`, Arbeitsbaum sauber.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `db-package-lists-check: OK — … dieselben 4 Pakete` · `coverage-gate: OK — Coverage 83.20% erfüllt Schwelle 80%` (`total: (statements) 83.2%`) · `d-check: 1108 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make test` (Race-Detector) | **EXIT=0** | 42 Pakete `ok`, 0 `FAIL`; `receive`, `replication/mapper`, `usecase/capture` je `ok` |
| `make test-replication` (PostgreSQL 18, Repo-Default) | **EXIT=0** (151 s) | `DB-Adapter-Coverage: 82.51% (gedeckt 873 von 1058 Statements; Profile gemergt: store,replication)`; `db-coverage: OK — DB-Adapter-Coverage 82.51% erfuellt Schwelle 70%`; `--- PASS: TestSlotReserveExhaustedIsConfiguration`; `--- PASS: TestWALRetentionThresholdsFollowGrowthAtInactiveSlot (1.47s)`; Paket `receive` `ok … 25.647s` |
| `make test-replication` (PostgreSQL 17, `PG_TEST_IMAGE` mit dem 17er-Digest aus `.github/workflows/e2e.yml`) | **EXIT=0** | dieselbe DB-Adapter-Coverage-Zeile (82.51 %, 873 von 1058), `--- PASS: TestSlotReserveExhaustedIsConfiguration`, `--- PASS: TestWALRetentionThresholdsFollowGrowthAtInactiveSlot (0.92s)`; das Log nennt das 17er-Image |
| `tools/bench-backfill.sh` einzeln (nach `make image`, Exit 0) | **EXIT=0** (304 s) | Lauf `20260925T053832Z`: acht der neun Runs „WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 0 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff“, ein Run (Stufe 50.000, Lauf 3/3) Spitze 1 MiB im Run, 0 MiB danach; „WAL-Rückstand der größten Stufe (200000 Zeilen) — höchste Spitze im Run 0 MiB, unter der Warnschwelle von 100 MiB“ (Kennung der Schwelle in der gedruckten Zeile: [`SPEC-013`](../../spec/pflichtenheft.md)); vom Slot gehaltenes WAL Median 6/35/141 MiB (Stufen 10.000/50.000/200.000) |
| `make doc-commits RANGE=f4e32fba..HEAD` | **EXIT=0** | `d-check: 1108 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=f4e32fba..HEAD` | **EXIT=0** | `d-check: 1108 Datei(en) geprüft, 0 Befund(e)`; `git diff --stat f4e32fba..HEAD -- docs/plan/adr` ist leer |
| `make commit-traceability RANGE=f4e32fba..HEAD` | **EXIT=0** | `OK — 13 Commit(s) in "f4e32fba..HEAD", Betreffs ohne Struktur-ID` |

**Nicht selbst gefahren:** `make test-integration` (E2E, ~5–6 min) — den Beleg trägt der reale
Post-Push-Lauf von `e2e.yml` für `427f6d1b` (§3); die Nullprobe der E2E-Phase liegt beim Reviewer
(Review-Report, Exit 2 nach 5 min 8 s), von mir nicht wiederholt.

Hygiene: dangling-Volumes (`docker volume ls -q -f dangling=true | wc -l`) vor den Läufen **34**, nach
allen Läufen **34**; keine Container nach den Läufen (`docker ps -a` leer); kein `prune`; den einen
Wegwerf-Container (`docker create` zum Kopieren des Coverage-Profils) mit `docker rm -fv` entfernt. Speicher vor
den Tier-Läufen (`free -m`, verfügbar): 16,2 bis 18,2 GB; die schweren Läufe liefen nacheinander.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Liefer-Punkt 1: Leerlauf-Bestätigung im Capture-Pfad; `make test` mit den Unit-Tests der beiden ersten Fitness-Function-Zeilen; je Zusage eine Eingabeseiten-Mutation | **erfüllt** | `make test` EXIT=0 (§1). Produktionscode im Diff gelesen: `confirmIdle` prüft `TransactionOpen()` und `serverWALEnd <= lastAcked`, ruft `IdleConfirmationInboundPort.ConfirmIdle`, setzt `lastAcked` aus dem Ergebnis, sendet im Bestätigungsfall kein zweites Update; `CaptureService.ConfirmIdle` ruft nur `ReplicationAckPort.Acknowledge`. Sieben eigene Eingabeseiten-Mutationen färben je genau die tragenden Unit-Tests rot (§4, M1–M7), darunter die fünf Mutationen, die der DoD-Text nennt |
| 2 | Liefer-Punkt 2: (a) Store-Tier X1, Nullprobe, Sicherheits-Test; (b) E2E-Phase mit Override der Fehlerschwelle, Nullprobe; die zwei Tests der alten Lage gezogen; Bindung des Sicherheits-Tests an seine Eingabe | **erfüllt, mit V-3 (INFO)** | (a) `make test-replication` EXIT=0 auf PostgreSQL 18 und 17 (§1); Store-Mutation S1 (gemeldete Position `+1 GiB`) färbt `TestStreamIdleConfirmationKeepsOpenTransactionDeliverable` („kein CaptureCommand innerhalb 19.99999958s“) und `TestStreamRestartsOnExistingSlot` rot (§4); die Mutation „Bedingung keine offene Transaktion entfernt“ färbt im Store-Tier **keinen** Store-Test, nur den Unit-Test (§4, S2 — die im Plan §6 benannte Bindung stimmt). (b) `abdeckung_declare "Leerlauf-Bestätigung …"` im Runner (Zeile 3254 von `tools/harness/run-integration-tests.sh`), Zeile in `docs/user/e2e-abdeckung.md` (Zeile 64); Post-Push-Lauf `36092733208` für `427f6d1b`, beide Legs grün (§3); Nullprobe des Reviewers gelesen, nicht wiederholt. Die zwei Tests der alten Lage sind gezogen: `TestStreamKeepaliveReportsAcknowledgedPosition` (mit `declineIdle`) und `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot` (`--- PASS`, beide PostgreSQL-Versionen) |
| 3 | Liefer-Punkt 3: Träger (Spec-Sätze, Architektur-Sicht, Handbuch, `bench-backfill.sh`/Target-Datei, Suchlauf, Bench-Ausgabe); `make gates` grün; Einzel-Lauf des Bench mit Rückstand je Run und Spitze der größten Stufe | **erfüllt** | Spec-Sätze und Architektur-Sicht ohne ADR-/Slice-/Wellen-Bezug (§6); Handbuch `Version: 1.57` mit Historienzeile (§6); eigener Bench-Lauf `20260925T053832Z` EXIT=0 mit der gedruckten Zeile „höchste Spitze im Run 0 MiB, unter der Warnschwelle von 100 MiB“ (§1); `make gates` EXIT=0 |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 5 | Review durchgeführt, Report liegt vor; kein offenes HIGH/MEDIUM | **erfüllt, mit V-2 (LOW)** | Report `review-slice-backfill-slot-leerlauf-bestaetigung.md` (0 HIGH, 2 MEDIUM, 4 LOW, 3 INFO); F-2 behoben (§5), F-1 als *weiter offen* an die Closure adressiert — die Zeile benennt das selbst (Berichtigungs-ADR); strenger gelesen bleibt F-1 ein offenes MEDIUM (V-2) |
| 6 | Verifikation als eigene Rolle mit Report | **dieser Report** | der Haken gehört dem Planner nach der Übergabe |
| 7 | §3.13-Suchlauf: Feld in §3, Gefundenes und Nichtgefundenes, beide Stände gemessen | **erfüllt, mit V-1 (LOW)** | zehn Zeilen des Feldes an beiden Ständen nachgefahren (§7): neun stimmen, das Zählwort „zehn Slices“ der Fixrunden-Zeile driftet um 1 (V-1) |
| 8 | Doku-Update (Handbuch samt Historie; `harness/README.md` nur, was der Lauf fuhr) | **erfüllt** | §6 |
| 9 | Closure-Notiz mit Lerneintrag | **korrekt offen** | §7 des Plans trägt „*(zu tragen bei Closure)*“ |
| 10 | Reconciliation-Register — entfällt | **korrekt entfällt** | Greenfield |
| 11 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht; Vorschläge stehen in §8 des Plans |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | die Ausgänge „bei der Closure“/„weiter offen“ sind Sache der Closure; zwei Ausgänge stehen schon (Startposition der Wache: *entfallen*, Bindung der Kette: *weiter offen* mit Grenze) |
| 13 | Die drei Paarungen | **korrekt offen** | hängen an der Closure von `welle-backfill-bestand` |

`[x]` sind acht Zeilen (Nr. 1–5, 7, 8, 10), `[ ]` fünf (Nr. 6, 9, 11, 12, 13), zusammen 13 — gezählt am
Plan (§2, `grep -c` auf `- [x]` und `- [ ]` im Abschnitt). Kein `[x]` ohne Beleg; kein `[ ]`, das über die
Rollen-Sequenz hinaus belegt wäre.

## 3. Post-Push-Lauf von `e2e.yml` ([`AGENTS.md`](../../AGENTS.md) §3.10)

`427f6d1b` ist der letzte gepushte Commit (`git log -1 origin/main`); `gh run list --commit
427f6d1b42e4915b9cd94760ae5efaa0fc076e0a`: `e2e` (Lauf `36092733208`), `ci` (`36092733158`) und
`examples` (`36092733262`) je `success`. `gh run view 36092733208 --json jobs`: beide Legs
(PostgreSQL 17 und 18) `success`, in beiden die Schritte „Compose-Integrationstest (Black-Box-E2E)“,
„DB-Adapter-Coverage — Store-Teil (Profil)“, „… Replication-Teil, Merge + Schwelle“ und
„Replication-Tier (go test ./... und Slot-Reserve)“ `success`.

**Was sich seit `427f6d1b` geändert hat:** `git diff --stat 427f6d1b..HEAD -- tools/harness compose.yaml
.github Makefile Dockerfile harness/mk internal cmd test` zeigt genau drei Dateien:
`internal/bootstrap/walretention_internal_test.go`, die umbenannte
`internal/bootstrap/walretention_slotgrowth_internal_test.go` und
`tools/harness/run-replication-tests.sh` (nur der Testname im `-run`-Muster und in der PASS-Prüfung).
`git diff --stat 427f6d1b..HEAD -- '*.go' ':!*_test.go'` ist leer: **kein Produktionscode**, kein Workflow,
kein Compose-Zug, kein Runner der E2E-Phase (`run-integration-tests.sh` unverändert) seit dem grünen Lauf.
Der Schritt „Replication-Tier“ ruft `bash tools/harness/run-replication-tests.sh tier`; dessen geänderte
Zeilen decken meine lokalen Läufe auf PostgreSQL 18 **und** 17 (§1, `--- PASS:
TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`). Für `HEAD` selbst fehlt der Post-Push-Lauf
(nicht gepusht); er bestätigt bei der Closure nur die drei Test-/Skript-Dateien, kein neues Verhalten.
[`AGENTS.md`](../../AGENTS.md) §3.10 greift nicht (kein Workflow im Diff); das PostgreSQL-17-Risiko des
Plans (§6) trägt der grüne Lauf `36092733208` für den Produktionsstand.

Laufzeit (Plan §6): Job-Dauer der Legs am Diff-Stand 12 min 36 s / 13 min 0 s gegen 11 min 21 s / 11 min
2 s am Parent (übernommen aus dem Review-Report, `gh run view` dort); mein lokaler `make test-replication`
(PostgreSQL 18) lief 151 s (der Reviewer nennt 92 s; Zeitstempel vor/nach, die Ursache der Differenz ist nicht
untersucht, die Grenze `timeout-minutes: 60` ist nicht berührt).

## 4. Sicherheitsbeleg und Mutationen der Eingabeseite (dieser Lauf)

**Hauptgegenstand — keine noch nicht gelieferte oder nicht persistierte Change wird übersprungen.**
Aufbau je Mutation: Datei im Arbeitsbaum mutiert (`sed` mit Ausgabe in eine Scratch-Datei und `cp`,
kein `sed -i`), Lauf ungefiltert in eine Log-Datei, gedruckte Rot-Meldung gelesen, Datei per
`git checkout` zurückgenommen. Ein erster Versuch brach am Kompilieren ab (M4 traf zwei Zeilen und
ließ `serverWALEnd` in `Process` undefiniert) und zählt nicht; M4 ist danach auf die eine Zeile
eingegrenzt wiederholt.

| # | Mutation | Lauf | Rot gesehen (gedruckt) |
|---|---|---|---|
| M1 | Prüfung `s.assembler.TransactionOpen()` in `confirmIdle` entfernt | `make test`, **EXIT=2** | `--- FAIL: TestRunNoConfirmationInsideOpenTransaction` (allein) |
| M2 | gemeldete Position `serverWALEnd+1` (`receive.go`) | `make test`, **EXIT=2** | `TestRunConfirmsIdleWithoutReplyRequested`, `TestRunIdleConfirmationSendsNoSecondUpdate`, `TestRunNoConfirmationInsideOpenTransaction`, `TestRunIdleConfirmationFailureEndsAsReplication` |
| M3 | `PersistTransaction` im Leerlauf-Pfad von `CaptureService.ConfirmIdle` (`service.go`) | `make test`, **EXIT=2** | `--- FAIL: TestConfirmIdleAcknowledgesOnlyThroughTheAckPort`, `FAIL … usecase/capture` |
| M4 | `lastAcked` aus dem WAL-Ende statt aus dem Ergebnis des Aufrufs (nur `confirmIdle`) | `make test`, **EXIT=2** | `--- FAIL: TestRunIdleConfirmationTakesPositionFromResult` |
| M5 | Rückschritt-Wache `<=` zu `<` | `make test`, **EXIT=2** | `TestRunKeepaliveWithoutReplyRequestedSendsNothing`, `TestRunIdleConfirmationSendsNoSecondUpdate`, `TestRunNoConfirmationBehindAcknowledgedPosition`, `TestRunReportsKeepaliveReplyFailure` |
| M6 | Bestätigung nur bei `ReplyRequested` (Position 0, wenn `replyRequested` falsch) | `make test`, **EXIT=2** | `TestRunConfirmsIdleWithoutReplyRequested`, `TestRunNoConfirmationInsideOpenTransaction`, `TestRunIdleConfirmationTakesPositionFromResult`, `TestRunIdleConfirmationFailureEndsAsReplication` |
| M7 | Rückschritt-Wache entfernt (`|| serverWALEnd <= s.lastAcked` gestrichen) | `make test`, **EXIT=2** | `TestRunAnswersKeepaliveWithAcknowledgedPosition`, `TestRunKeepaliveWithoutReplyRequestedSendsNothing`, `TestRunIdleConfirmationSendsNoSecondUpdate`, `TestRunNoConfirmationBehindAcknowledgedPosition`, `TestRunReportsKeepaliveReplyFailure` |
| S1 | **Store-Tier:** gemeldete Position `+ (1<<30)` (1 GiB hinter dem WAL-Ende) | `make test-replication`, **EXIT=2** (63 s im Paket `receive`) | `--- FAIL: TestStreamIdleConfirmationKeepsOpenTransactionDeliverable (20.68s)` („stream_test.go:1102: kein CaptureCommand innerhalb 19.99999958s“), `--- FAIL: TestStreamRestartsOnExistingSlot (21.20s)` („stream_test.go:821: kein CaptureCommand …“) neben den Unit-Tests aus M2 |
| S2 | **Store-Tier:** Prüfung `TransactionOpen()` entfernt (= M1) | `make test-replication`, **EXIT=2** | **nur** `--- FAIL: TestRunNoConfirmationInsideOpenTransaction`; alle Store-Tests einschließlich des Sicherheits-Tests bleiben grün |

**Sicherheitsurteil.** Die Größe der bestätigten Position trägt der Store-Test am realen Verhalten: der
Sicherheits-Test (offene Quelltransaktion auf einer veröffentlichten Tabelle → Leerlauf-Bestätigung, die
`confirmed_flush_lsn` hinter den ersten Change der offenen Transaktion rückt → Commit → Neustart →
der Change `{"id":"2","name":"Offen"}` wird geliefert) färbt sich bei einer zu großen Position rot (S1). Die
Eigenschaft „inmitten einer Quelltransaktion wird nicht bestätigt“ trägt allein der Unit-Test mit der
Fake-Sitzung (S2 = M1). Das ist die im Plan §6 benannte Lage, nicht eine verdeckte: der Store-Test belegt
„eine Transaktion, die vor P beginnt und danach committet, geht nicht verloren“ (die Quelle liefert jede
Transaktion, deren Commit hinter `confirmed_flush_lsn` liegt), nicht die Mutation der Fitness-Function-Zeile.
Der Reviewer hat den Zustand „Keepalive inmitten der auf der Leitung laufenden Transaktion“ mit einem
Wegwerf-Test hergestellt (Position gleich der Commit-LSN, Neustart lieferte `changes=400000`); den Wegwerf-Test habe
ich nicht wiederholt, seine Aussage steht als übernommen (Review-Report F-1). Ergebnis: kein Weg gefunden, auf dem
die Leerlauf-Bestätigung eine noch nicht gelieferte Change überspringt; die Datensicherheit hängt außerhalb
der Größe der Position am Verhalten der Quelle (V-3).

Eine weitere Prüfung ist äquivalent, keine Bindungslücke: die Leerlauf-Bestätigung schreibt nichts nach
`cdc.transaction`/`cdc.change` und berührt `cdc_capture_lag` nicht — dieser Wert folgt dem `committed_at` der
persistierten Transaktionen (`tools/schema/nacharbeit-observability.sql`); der Verdrahtungs-Test
`TestRealIdleConfirmationReleasesForeignWAL` hält beide Tabellen unverändert (gelesen, nicht gemutet).

## 5. Fixrunde: Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) ADR-Zeile „Store: Mutation inmitten der Transaktion“ nicht erfüllbar | eigene Messung S2 (§4): die Mutation färbt im Store-Tier keinen Store-Test. Plan §6 nennt die gemessene Aussage statt der Herleitung, Ausgang *weiter offen*, Berichtigung als neue ADR mit `Supersedes` bei der Closure; [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) unberührt (`git diff --stat f4e32fba..HEAD -- docs/plan/adr` leer, `make doc-immutable` EXIT=0). Kein Handbuch-, Sensor- oder Spec-Satz behauptet die Store-Bindung (Muster-Suchlauf §7, elf Zeilen, gelesen) | **ehrlich benannt, offen** (V-2) |
| F-2 (MEDIUM) Test unter gleichem Namen abgeschwächt | `TestWALRetentionThresholdEndToEnd` → `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`; Godoc gelesen: „Der Test fährt weder einen Stream noch `Run`“, Abbruchfunktion als Kontext-Abbruch, Prioritätskette in `TestMergeStreamAndWALFaultOutcome…`/`TestClassifyRunError…`; `git grep -n 'TestWALRetentionThresholdEndToEnd\|walretention_endtoend' HEAD` (ohne Records und Plan): **0 Treffer**; PASS auf PostgreSQL 18 und 17 (§1). Kette „Fehlerschwelle → Prozessende“ in Plan §6 auf vier Träger verteilt, Grenze („Schwelle erreicht → Container endet“ trägt nur die Nullprobe des Reviewers) benannt | **behoben, Grenze benannt** |
| F-3 (LOW) Umbenennung und Neufassung in einem Commit | `git show --stat -M ec7e43dd`: `rename … (100%)`, 0 Insertions/0 Deletions — reiner Move; der Inhalt folgt in `5a5d3422`. Der ursprüngliche Bruch bei `00caec48` (`A`/`D`) bleibt in der Historie; im Plan §8 als Register-Vorschlag benannt | **behoben für die Fixrunde**, Altbruch benannt |
| F-4 (LOW) Handbuch: Lauf ohne auflösbaren Träger | Diff `427f6d1b..HEAD` gelesen: Lauf `20260925T032925Z` an beiden Stellen als „übernommen aus dem Lauf-Bericht des Implementers, im Repository nicht auflösbar“ gekennzeichnet; Nachmessung `20260925T043056Z` mit Verweis auf den Review-Report; Spanne der Richtgröße 2.000.000 bis 5.000.000 (Code-Wert 4.000.000 innerhalb). Mein eigener Lauf `20260925T053832Z` (§1): Richtgröße 5.000.000 bei 8.909 Zeilen/s — innerhalb der Spanne; gehaltenes WAL Median 141 MiB bei 200.000 Zeilen wie im Handbuch | **behoben** |
| F-5 (LOW) Suchlauf-Zahl gedriftet | Zeile 1 des Feldes mit Plan-Ausschluss an `67e331ac` und `HEAD`: 85 Zeilen in 26 Dateien (Feld: 85/26); mit Plan-Datei am Stand `427f6d1b`: 94/27 (Feld: 94/27). Restliche Zeilen §7 | **behoben**, mit V-1 für eine Nachbarzeile |
| F-6 (LOW) Konjunktiv-Nebenklausel im Testkommentar | `walretention_slotgrowth_internal_test.go`: „Die Instanz gehört dem Test allein …: der Rückstand misst das WAL der ganzen Instanz, ein gleichzeitiger Schreiber eines anderen Test-Pakets verfälscht ihn“ — Indikativ, Kopplung | **behoben** |
| F-7 (INFO) gedeckte Zahl 2136 gegen 2135 | eigenes Coverage-Profil (`/out/coverage.out` aus dem Image `pg-change-feed:coverage` dieses `make gates`-Laufs, dedupliziert über die Block-Position, `docker cp`): **2136 von 2567 = 83,21 %**, gleich dem Träger `harness/sensors/coverage-gate.md`; der 1-Statement-Unterschied des Reviewers war laufabhängig | **bestätigt, kein Handlungsbedarf** |
| F-8 (INFO) `lastAcked` startet bei 0 | Plan §6 nennt die Bewertung am Code, die Messung des Reviewers und den Ausgang *entfallen* (unschädliche Eigenschaft, kein committeter Test); der Parent `f4e32fba` trägt dieselbe Eigenschaft bei der Keepalive-Antwort. Von mir nicht nachgemessen; die Aussage bleibt eine Messung des Reviewers | **ehrlich benannt** |
| F-9 (INFO) Laufzeit | §3 | **belegt** (übernommen aus `gh run view` des Reviewers; meine Läufe §1) |

## 6. Träger-Nachzüge

| Träger | Beleg | Verdikt |
|---|---|---|
| [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md), [`SPEC-009`](../../spec/pflichtenheft.md) | `git diff f4e32fba..HEAD -- spec`: Schritt 5 und Invariante um die Bestätigung im Leerlauf ergänzt; `SPEC-009` Zeile `cdc_wal_retention_bytes`: „das vom Feed noch nicht bestätigte WAL …“; Änderungshistorie um eine Zeile; die hinzugefügten Zeilen tragen weder `ADR-` noch `slice` noch `welle` (`grep -i`: 0 Treffer außerhalb der ARC-Kennungen des Diagramms) | erfüllt |
| `spec/architecture.md` | Absatz „Leerlauf-Weg“ und Sequenzdiagramm (Replication Stream → `IdleConfirmationInboundPort` → Application → `ReplicationAckPort` → ACK-Adapter, „ohne PersistTransaction“); Diagramm-Teilnehmer tragen die `ARC-`-Kennungen der Komponenten wie die Nachbardiagramme; kein ADR-/Slice-/Wellen-Bezug ([`AGENTS.md`](../../AGENTS.md) §3.4) | erfüllt |
| Handbuch `Version: 1.57` | Zeile 3 `Version: 1.57`; Historienzeilen 1.56 und 1.57 vorhanden (Datei-Ende); Diff `427f6d1b..HEAD` gelesen (F-4); keine neue Betreiber-Oberfläche | erfüllt |
| `harness/README.md` (Zeilen `make test-replication`, `make test-integration`, `make bench`) | Lauf-Aussagen gegen meine Läufe: `make test-replication` (Store-Test-Belege, zweite Instanz, Schwellen-Beleg ohne Stream und ohne `Run`), `make bench` (vier Skripte, das vierte ohne Schwelle, Rückstands-Vergleichszeile) — gelesen, deckt sich mit §1 | erfüllt |
| `harness/sensors/coverage-gate.md`, `db-adapter-coverage.md` | Nenner 2567 und gedeckt 2136 (mein Profil, F-7), gemergter Nenner 1058 und „873 von 1058, 82,51 %“ (§1, beide PostgreSQL-Versionen); die `wiring.go`-Lokatoren `:1287.4,1288.1` (Zeile 1287: `return` unter `if ctx.Err() != nil` in `runAdministration`) und `:1183.5,1184.13` (Zeile 1183: `log.Warn`, 1184: `continue` in `runWALRetentionCheck`) am Quelltext von `HEAD` nachgelesen | erfüllt |
| `docs/user/e2e-abdeckung.md`, `tools/bench-backfill.sh`, `harness/targets/bench-backfill.md` | Zeile 64 der Abdeckungstabelle mit `LH-FA-CAP-009`/`LH-QA-REL-001`; Bench druckt „WAL-Rückstand der größten Stufe … unter der Warnschwelle von 100 MiB“ (§1); Zeile „WAL-Rückstand-Schwellen (abgeleitet)“ entfallen (`git grep` ohne Treffer in `tools/bench-backfill.sh`) | erfüllt |
| [`LH-FA-CAP-009`](../../spec/lastenheft.md) im RTM | Runner-Phase mit `abdeckung_declare` auf `LH-FA-CAP-009,LH-QA-REL-001`, `make gates` (`docs-check` inkl. Matrix) 0 Befunde | erfüllt |

## 7. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent per `git grep … f4e32fba` bzw. `67e331ac` (der Stand vor der Fixrunde, den die Fixrunden-Zeile
nennt), Diff-Stand per `git grep … HEAD`; Ausschlüsse wie im Feld
(`:!docs/reviews`, `:!…/0120-…`, `:!docs/plan/planning/done`, `:!docs/plan/planning/observations`, die
Plan-Datei ausgeschlossen). Gedruckte Zahlen (Zeilen bzw. Dateien nach dem Befehl der Plan-Zeile):

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| Keepalive-Antwort/bestätigte Position (Zeilen/Dateien) | `f4e32fba` **64/26**, `67e331ac` **85/26**, `HEAD` **85/26**; ohne Ausschluss der Plan-Datei `427f6d1b` **94/27** | bestätigt |
| „Rückstand wächst durch fremdes WAL / Live-Commit“ | `f4e32fba` **12/6**, `67e331ac` und `HEAD` je **2/2** | bestätigt |
| Symbolnamen (Dateien) | `f4e32fba` **19**, `67e331ac` und `HEAD` **18**; `427f6d1b` (mit Plan-Datei) **19** | bestätigt |
| `wiring.go`-Lokatoren | `67e331ac` und `HEAD` je **4/2** (`ADR-0082`, `ADR-0088`, `Accepted`); die zwei Blöcke am Quelltext bestätigt (§6) | bestätigt |
| Bench: Schwellen-Text | `f4e32fba` **11/6**, `67e331ac` **10/6**, `HEAD` **9/5** | bestätigt |
| `make test-replication` (`wc -l`) | `f4e32fba` **78**, `67e331ac` und `HEAD` **74**, `427f6d1b` **81** | bestätigt |
| Coverage-Zahlen in `docs/plan/planning/open` | `f4e32fba` und `HEAD` je **0** | bestätigt |
| Name und Datei des Schwellen-Tests | `67e331ac` **9**, `HEAD` **0** | bestätigt |
| Store-Bindung der Bedingung „keine offene Transaktion“ | `67e331ac` und `HEAD` je **11** Zeilen | bestätigt; keine der Zeilen behauptet die Store-Mutation (gelesen im Feld des Plans, §5 F-1) |
| Zählwörter „zehn Slices“ | ohne die Plan-Datei `f4e32fba` und `HEAD` je **3** (`roadmap.md` 48, `welle-transformationen.md` 39 und 84); mit der Plan-Datei 4 bzw. 5 (Selbstverweise) | **abweichend**: die Fixrunden-Zeile nennt „ohne Plan-Datei 4“ (V-1) |

Nichtfund-Suche eigener Art: `git grep -n -i -E 'Live-Commit|nie bestätigt|weder BEGIN noch COMMIT'` über
`HEAD` (ohne Records, Reviews, Plan) trifft nur noch `stream_test.go:984` (Beschreibung des Tests) und
`spec/lastenheft.md:733` (anderer Gegenstand, Consumer). Kein Träger beschreibt die alte Lage (Rückstand
sinkt erst durch einen Live-Commit) als geltend. Ein Träger außerhalb des Feldes bleibt offen: siehe V-2.

## 8. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung 1 (Leerlauf, Meldung, Application bestätigt, nie inmitten, nie zurück, Update sofort, höchstens eines je Keepalive) | Punkte 1–5 | `confirmIdle`/`handleCopyData` im Diff gelesen; M1 (inmitten), M5/M7 (Rückschritt), M6 (auch ohne `ReplyRequested`), M4 (`lastAcked` aus dem Ergebnis), M2 (Größe), `TestRunIdleConfirmationSendsNoSecondUpdate` (höchstens ein Update); Store-Belege X1 und Sicherheits-Test (§1, S1) | konform; die Store-Zeile der Fitness Function („Mutation inmitten der Transaktion → Change fehlt“) nicht erfüllbar (V-3) |
| [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung 2 (Position ohne Change, Schwellen und Formel, gehaltenes WAL und Spill bleiben, Richtgröße unberührt) | nichts in `cdc.transaction`/`cdc.change`, `SPEC-013` unverändert, Handbuch nennt die Grenze | `git diff f4e32fba..HEAD` berührt weder `SPEC-013` noch `harness/mk/**`/`Dockerfile`/`THRESHOLD`; gehaltenes WAL im eigenen Bench-Lauf 141 MiB bei 200.000 Zeilen (Handbuch: Grenze der Ein-Transaktions-Form) | konform |
| [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung 3 (Klasse `replication`, kein neuer Sentinel) | Fehler des Ports geht als `ErrReplication` mit der Ursache dahinter durch | `confirmIdle`: `fmt.Errorf("%w: Leerlauf-Bestätigung: %w", ErrReplication, err)`; `CaptureService.ConfirmIdle` reicht den Fehler unverändert durch; `TestRunIdleConfirmationFailureEndsAsReplication` rot in M2/M6, `TestConfirmIdleForwardsPortErrorUnchanged` im Reviewer-Lauf (M7 dort) | konform |
| [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung 4/5 (Zuschnitt, Port als Arbeitsname) | eigener Slice; Form des Ports vom Slice | eigener Port `IdleConfirmationInboundPort`, `CaptureInboundPort` unberührt (Stand-ins unberührt, `make a-check` `gesamt: 0 Befund(e)`) nach [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md); Verdrahtung: eine `CaptureService`-Instanz an beiden Eingängen (`wiring.go`), `Run` ohne den zweiten Port endet als `ErrConfiguration` | konform |
| [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) Folgepflichten 1–5 | Umsetzung, Spec-Nachzug, Handbuch, Bench, Suchlauf | §1 bis §7 | konform |
| [`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md) (die Application entscheidet, wann bestätigt wird) | der Adapter bestätigt nicht selbst | `receive.go` ruft nur den Inbound-Port; `lastAcked` wird erst aus dem Ergebnis gesetzt (M4 rot); `ack`-Aufruf sitzt in `usecase/capture` (M3 rot bei `PersistTransaction`) | konform |
| [`ADR-0011`](../plan/adr/0011-persist-before-ack.md) / [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) | eine bestätigte Position verdeckt keinen nicht gespeicherten Change | §4: Größe der Position im Store-Tier gebunden (S1), „inmitten“ im Unit-Tier (M1), Herleitung („Quelle liefert jede Transaktion mit Commit hinter `confirmed_flush_lsn`“) vom Reviewer am realen Stream gemessen | konform, mit V-3 |
| [`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md) (Schwellen, Klasse `replication`) | unverändert | `classifyRunError`, `runWALRetentionCheck` unberührt (`git diff` ohne Änderung an den Produktionsdateien außer der Verdrahtungs-Zeile); Schwellen-Test gezogen (§5 F-2) | konform |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) (Capture-kritischer Pfad unberührt für Backfill) | Backfill-Code nicht im Diff | `git diff --stat f4e32fba..HEAD` trägt kein Paket `usecase/backfill` oder `postgressnapshot`; E2E-Phase belegt „Run `completed`, Container läuft“ (§3) | konform |
| [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) | `Accepted`, unberührt; Suchlauf zu Kommentar und Test von `handleCopyData`/`standbyStatus` | Datei nicht im Diff; Kommentar in `receive.go` und die Keepalive-Tests tragen die neue Aussage (Suchlauf-Zeile 1, §7) | konform |

## 9. Plan-vs-Code-Diff

Jede Zeile der Plan-Tabelle §3 ist im Diff (`git diff --stat f4e32fba..HEAD`, 26 Dateien) vertreten:
`receive/receive.go`, `receive/seam_test.go`, `receive/stream_test.go`, `replication/mapper/mapper.go` und
`mapper_test.go`, `usecase/capture/service.go` und `service_test.go`, `port/inbound/idleconfirmation.go`
(neu), `bootstrap/wiring.go`, `bootstrap/replication_stream_test.go`, `bootstrap/walretention_internal_test.go`,
`bootstrap/walretention_slotgrowth_internal_test.go` (Ersatz von `walretention_endtoend_test.go`),
`tools/harness/run-integration-tests.sh`, `tools/harness/run-replication-tests.sh`, `tools/bench-backfill.sh`,
`harness/targets/bench-backfill.md`, `spec/pflichtenheft.md`, `spec/architecture.md`,
`docs/user/benutzerhandbuch.md`, `docs/user/e2e-abdeckung.md`, `harness/README.md`,
`harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md`, Plan und Review-Report.
`compose.yaml` steht im Plan als „prüfen, unverändert“ und ist nicht im Diff. **Über den Plan hinaus**
(im Plan als „Nachzug des Implementers“ und „Fixrunde“ eingetragen, vor dem Sensor-Lauf): der eigene Port,
`ErrMissingIdlePosition`, der Stream verlangt beide Ports, die zweite Instanz und der gesonderte Lauf des
Schwellen-Tests im Runner, die Compose-Override-Datei im Temp-Verzeichnis des Runners, `settle_wal` statt
`release_wal`. Ungeplant und nicht eingetragen: nichts. Lifecycle-Moves rein: `1f75a2be` (`open/` → `next/`)
und `5e84274a` (`next/` → `in-progress/`) tragen im Diff keine Inhaltsänderung
([`AGENTS.md`](../../AGENTS.md) §3.3); `Verantwortlich` steht in `52b6d6f0` als eigener Commit zwischen
beiden; der Fixrunden-Move `ec7e43dd` ist ein reiner Rename (100 %).

## 10. Harte Regeln

- **§3.1** — alle Läufe über `make`, gepinnte Images und `bash tools/…`; Textänderungen der Mutationen per
  `sed` mit Ausgabe auf stdout und `cp`; Coverage-Profil per `docker cp`; kein Host-Compiler.
- **§3.2** — `git diff f4e32fba..HEAD` über `*.go`, `*.sh`, `*.sql`, `*.yaml` (`+`-Zeilen): kein `nolint`.
- **§3.3** — siehe §9 (rein).
- **§3.4** — Architektur-Diff ohne ADR-/Slice-/Wellen-Bezug (§6).
- **§3.5** — `make doc-immutable RANGE=f4e32fba..HEAD` EXIT=0; keine `Accepted`-ADR im Diff.
- **§3.6** — keine Schwelle gesenkt: `Dockerfile`, `harness/mk/**`, `THRESHOLD`, `SPEC-013` nicht im Diff.
- **§3.9** — Exit-Codes nie durch eine Pipe gemessen (§1); Gate-Lauf und Folgehandlung getrennt.
- **§3.10** — der Diff berührt keinen Workflow; der reale Post-Push-Lauf des bestehenden `e2e.yml` für
  `427f6d1b` ist trotzdem herangezogen (§3).
- **§3.11** — `make docs-check` (in `make gates`) 0 Befunde; dieser Report trägt keinen host-lokalen Pfad.
- **§3.12** — jede Zahl dieses Reports ist **gemessen** (Lauf oder Befehl genannt) oder als **übernommen**
  gekennzeichnet (Job-Dauern der Legs, Wegwerf-Test des Reviewers zu F-1, Nullprobe der E2E-Phase, Messung
  zu F-8).
- **§3.13** — Suchlauf-Feld an beiden Ständen nachgefahren (§7).
- **Handbuch-Pflicht** — `Version: 1.56` → `1.57` mit Historienzeile in der Fixrunde.
- **Commit-Traceability** — 13 Commits, jede Message nennt eine Kennung, keine Struktur-Kennung im Betreff (§1).

## 11. Befunde

Kein Befund blockiert. V-1 und V-2 sind LOW, V-3 bis V-5 INFO.

- **V-1 (LOW) — Zählwort „zehn Slices“ der Fixrunden-Zeile driftet um 1.** Die Zeile „Fixrunde nach dem
  Review — Zahlen“ (Plan §3) nennt „Zählwörter „zehn Slices“: 4“ für den Parent `67e331ac` ohne Plan-Datei und
  für den Diff-Stand. Mit dem Befehl der Zeile und dem Plan-Ausschluss messe ich an beiden Ständen **3**
  (`roadmap.md` Zeile 48, `welle-transformationen.md` Zeilen 39 und 84); die vierte Zeile der Zahl 4 ist die
  Selbstreferenz der Plan-Datei (Zeile mit dem Titel „Zählwörter „zehn Slices““), die der Ausschluss
  entfernt. Die Aussage der Zeile (die Zählwörter zählen die andere Welle und bleiben wahr) ist richtig, die
  Zahl gehört zum Stand ohne Ausschluss. Aktion: Zahl in der Closure auf 3 ziehen oder als „mit Selbstverweis“ kennzeichnen.
- **V-2 (LOW) — Träger außerhalb des Diffs führt den überholten „Verdacht“ zu F-1.**
  `docs/plan/planning/welle-backfill-bestand.md` Zeile 503 nennt die Store-Zeile der Fitness Function
  „möglicherweise nicht rot … ein Verdacht, kein Beleg, deshalb nicht gezählt“; nach der Messung des
  Reviewers (F-1) und meiner S2 ist sie ein Beleg. Der Plan meldet den Träger an den Planner (fremde Datei,
  Suchlauf-Feld §3) — der Nachzug steht aus. Dazu die Lesart der DoD-Zeile 5: F-1 bleibt ein offenes
  MEDIUM mit dem Ausgang „Berichtigungs-ADR bei der Closure“ und blockiert die **Closure**, nicht diese
  Verifikation. Aktion Planner/Architect: Berichtigungs-ADR mit `Supersedes` für die Fitness-Function-Zeile
  und die Aussage „hinter Nachrichten, die noch nicht gespeichert sind“ ([`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md)
  Festlegung 1 Punkt 3); Welle-Datei Zeile 503 nachziehen; Register-Evidence zu
  `fitness-function-gegen-eigene-entscheidung` anlegen.
- **V-3 (INFO) — Die Datensicherheit hängt bei „inmitten der Transaktion“ am Verhalten der Quelle.** Der
  Store-Tier bindet nur die Größe der bestätigten Position (S1); die Bedingung „keine offene Transaktion“
  trägt allein der Unit-Test mit der Fake-Sitzung (S2). Was ein Keepalive inmitten einer auf der Leitung
  laufenden Transaktion an Position trägt (Commit-LSN dieser Transaktion, Neustart liefert sie vollständig),
  ist am realen Stream nur vom Reviewer gemessen (Wegwerf-Test, PostgreSQL 18, im Repository nicht
  vorhanden); PostgreSQL 17 hat diese Messung nicht, der Sicherheits-Test und X1 laufen dort grün (§1,
  §3). Die Re-Evaluierungs-Trigger von [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md)
  („neue PostgreSQL-Hauptversion“, „Streaming großer Transaktionen“) sind die benannte Grenze; kein
  committeter Test bindet die Quellseite.
- **V-4 (INFO) — Kein Post-Push-Lauf für `HEAD`.** `origin/main` steht auf `427f6d1b`; die fünf Folge-Commits
  ändern nur zwei Test-Dateien, einen Runner-Skript-Namen und Doku (§3). Der Schritt „Replication-Tier“
  läuft mit dem geänderten Skript erst nach dem Push in CI; lokal auf PostgreSQL 17 und 18 belegt.
- **V-5 (INFO) — Tier-Nebenwirkungen und Hygiene.** `make test-replication` und `tools/bench-backfill.sh`
  überschreiben `tools/schema/plan.yaml` (und `down.sql`); nach jedem Lauf per `git checkout`
  zurückgenommen. `make image` (für den Bench) schreibt `harness/image-hash.txt` (lokal, nicht committet,
  [`ADR-0103`](../plan/adr/0103-image-hash-lokal-statt-committet.md)).

## 12. Grenzen dieses Laufs

- Kein lokaler Lauf von `make test-integration` und keine Wiederholung der E2E-Nullprobe — der Beleg ist der
  reale Post-Push-Lauf `36092733208` (beide Legs) und die Nullprobe des Reviewers (übernommen).
- Keine Wiederholung des Wegwerf-Tests des Reviewers (Keepalive inmitten der Transaktion, F-1) und der
  Messung zu F-8; beide stehen als übernommen.
- Tier-Läufe: `make test-replication` auf PostgreSQL 18 und 17; die Store-Mutation S1/S2 nur auf PostgreSQL 18.
- Mutationen der E2E-Assertions (über die Nullprobe hinaus) nicht gefahren.

## 13. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (8) | **8 von 8 erfüllt**, je mit eigenem Beleg (Nr. 2 mit V-3, Nr. 5 mit V-2, Nr. 7 mit V-1) |
| DoD §2 — übrige `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Verifikation = dieser Report; Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) |
| Sicherheit der Leerlauf-Bestätigung | **belegt**: Größe der Position im Store-Tier gebunden (S1 rot), „inmitten“ im Unit-Tier gebunden (M1/S2), sieben Unit-Mutationen rot; Sicherheits-Test grün auf PostgreSQL 17 und 18 |
| Review-Findings F-1 (MEDIUM), F-2 (MEDIUM), F-3 bis F-6 (LOW), F-7 bis F-9 (INFO) | F-2 bis F-7 **behoben bzw. bestätigt**, je nachgemessen; F-1 **ehrlich als offen benannt** (Berichtigungs-ADR bei der Closure); F-8/F-9 benannt bzw. belegt |
| Plan-vs-Code-Diff | **deckungsgleich**; Abweichungen im Plan als Nachzug/Fixrunde eingetragen, ungeplant nichts |
| Entscheidungen | [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md), [`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md), [`ADR-0011`](../plan/adr/0011-persist-before-ack.md), [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md), [`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md), [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md), [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) **konform** (Fitness-Function-Store-Zeile: V-3, offen) |
| Post-Push-Lauf `e2e.yml` | **grün, beide Legs** für `427f6d1b` (Lauf `36092733208`); seit dem Lauf kein Produktionscode, kein Workflow, kein E2E-Runner geändert (§3) |
| Mutationen der Eingabeseite | **acht rot** (M1–M7, S1) plus S2 als belegte Nicht-Bindung im Store-Tier |
| Harte Regeln | **erfüllt** (§10) |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=f4e32fba..HEAD`) |
| Gates | **`make gates` EXIT=0**, `make test` EXIT=0, `make test-replication` EXIT=0 (PostgreSQL 18 und 17), `tools/bench-backfill.sh` EXIT=0 |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung und den realen Post-Push-Lauf für
den Produktionsstand belegt. `make gates`, `make test` und `make test-replication` laufen am Stand `80b451c4`
mit Exit 0; die Leerlauf-Bestätigung überspringt nach eigener Messung keine noch nicht gelieferte oder nicht
persistierte Change: jede Zusage der Fitness Function färbt bei einer Mutation ihrer Eingabeseite ihren Beleg
rot — im Store-Tier die Größe der bestätigten Position (S1), im Unit-Tier die Bedingung „keine offene
Transaktion“, die Rückschritt-Wache, die Quelle der bestätigten Position, das Fehlen von `PersistTransaction`
und die Bestätigung ohne `ReplyRequested`. Die Store-Zeile der ADR-Fitness-Function („Mutation inmitten der
Transaktion → Change fehlt“) ist nicht erfüllbar (S2); das ist im Plan §6 gemessen benannt und geht als
Berichtigungs-ADR an die Closure — ein offener Ausgang, kein verdeckter. Kein Befund blockiert diese
Verifikation; V-1 und V-2 (LOW) sind Nachzüge für den Planner, V-3 bis V-5 INFO.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: Berichtigungs-ADR mit `Supersedes` zu
[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) (Architect/Planner), Nachzug von
`welle-backfill-bestand.md` Zeile 503 und der Zahl „zehn Slices“ (V-1, V-2), Closure-Notiz mit Lerneintrag
(die Finding-Klassen des Reviews gehen in den Zähler), Ausgänge der §6-Risiken, Beobachtungs-Register,
Welle-Paarungen; nach dem Push der Beleg des Post-Push-Laufs für `HEAD`. Danach der reine `git mv` nach
`done/` ([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2).
