# Review-Report: slice-backfill-slot-leerlauf-bestaetigung — 2026-09-25

**Review-Art:** Code — der Diff führt die Leerlauf-Bestätigung des Capture-Slots ein (Empfangs-Schleife
`receive`, `Assembler.TransactionOpen`, Inbound-Port `IdleConfirmationInboundPort`,
`CaptureService.ConfirmIdle`, Verdrahtung), die Belege in allen drei Tiers (Unit, Store-Tier, E2E-Phase),
die Bench-Ausgabe und die Träger (Pflichtenheft, Architektur-Sicht, Handbuch 1.56, Sensor-Dateien,
`harness/README.md`, Plan); geprüft gegen Plan, ADRs, Pflichtenheft und `AGENTS.md` Hard Rules (Modul 10
§Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11). Hauptprüfgegenstand ist die
Sicherheit des Capture-kritischen Pfads: die Leerlauf-Bestätigung überspringt keine noch nicht gelieferte
oder nicht persistierte Change.

**Gegenstand:** Slice `slice-backfill-slot-leerlauf-bestaetigung`, Diff-Range `f4e32fba..427f6d1b`
(8 Commits, 25 Dateien, +1732/−528; drei reine `git mv`-Commits `1f75a2be`, `5e84274a` und die
Einzeilen-Änderung `52b6d6f0`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere HIGH-Klassen
ergänzt, u. a. Zahl-im-Träger, Beleg-Satz, Zusage-ohne-Eingabeseite, Kommentar-Chronik, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-25.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-slot-leerlauf-bestaetigung` (§1 Ziel, §2 DoD als Prüfmaßstab für
  Plan-Zusagen, §3 Plan samt Suchlauf-Feld, §6 Risiken) und Welle `welle-backfill-bestand`; Architect-Verdikt
  `architect-verdict-backfill-wal-rueckstand-und-bench-rot`
- [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) (Festlegungen 1–5, Folgepflichten
  1–5, Fitness-Function-Tabelle, Re-Evaluierungs-Trigger),
  [`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md),
  [`ADR-0011`](../plan/adr/0011-persist-before-ack.md),
  [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
  [`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md),
  [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md),
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
  [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-QA-REL-001`](../../spec/lastenheft.md) (und [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md)),
  [`LH-FA-CAP-009`](../../spec/lastenheft.md), [`LH-QA-REL-003`](../../spec/lastenheft.md),
  [`SPEC-009`](../../spec/pflichtenheft.md), [`SPEC-013`](../../spec/pflichtenheft.md)
- `AGENTS.md` (Hard Rules §3.1–§3.13), `harness/conventions.md` (`MR-000`/`MR-001`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-e2e.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert):

- **Gates am Stand `427f6d1b`:** `make test` (`-race`) Exit 0; `make a-check` Exit 0 („gesamt: 0 Befund(e)“);
  `make coverage-gate` Exit 0, gedruckt „total: (statements) 83.2%“ und „coverage-gate: OK — Coverage
  83.20% erfüllt Schwelle 80%“; `make test-replication` (PostgreSQL 18, Pin von `PG_TEST_IMAGE`) Exit 0
  in 92 s (Zeitstempel vor/nach), gedruckt „DB-Adapter-Coverage: 82.51% (gedeckt 873 von 1058 Statements;
  Profile gemergt: store,replication)“, `--- PASS: TestSlotReserveExhaustedIsConfiguration` und
  `--- PASS: TestWALRetentionThresholdEndToEnd (0.93s)`; Paket `receive` 18,8 s. `tools/schema/plan.yaml`
  als Tier-Nebeneffekt danach zurückgenommen.
- **CI zu `427f6d1b`:** `ci` success (2 min 36 s), `examples` success (1 min 49 s), `e2e` success, beide
  Legs (PostgreSQL 17: 12 min 36 s, PostgreSQL 18: 13 min 0 s Job-Dauer); am Parent `f4e32fba` 11 min 21 s
  und 11 min 2 s.
- **Mutationen der Eingabeseite** (`make test`, je danach `git checkout`; alle färben genau die tragenden
  Tests rot):

  | Nr. | Mutation | roter Test |
  |---|---|---|
  | M1 | Prüfung `TransactionOpen()` in `confirmIdle` entfernt | `TestRunNoConfirmationInsideOpenTransaction` (allein) |
  | M2 | Rückschritt-Wache `<=` zu `<` | `TestRunNoConfirmationBehindAcknowledgedPosition` (Fall „gleich“), `TestRunKeepaliveWithoutReplyRequestedSendsNothing`, `TestRunIdleConfirmationSendsNoSecondUpdate`, `TestRunReportsKeepaliveReplyFailure` |
  | M3 | gemeldete Position `P+1` | `TestRunConfirmsIdleWithoutReplyRequested`, `TestRunIdleConfirmationSendsNoSecondUpdate`, `TestRunNoConfirmationInsideOpenTransaction`, `TestRunIdleConfirmationFailureEndsAsReplication` |
  | M4 | `lastAcked` aus dem gemeldeten WAL-Ende statt aus dem Ergebnis | `TestRunIdleConfirmationTakesPositionFromResult` |
  | M5 | zweites Update je Keepalive (`confirmed` ignoriert) | `TestRunIdleConfirmationSendsNoSecondUpdate` |
  | M6 | `PersistTransaction` im Leerlauf-Pfad von `ConfirmIdle` | `TestConfirmIdleAcknowledgesOnlyThroughTheAckPort` |
  | M7 | Fehler des Ack-Ports verschluckt | `TestConfirmIdleForwardsPortErrorUnchanged` |
  | M8 | `TransactionOpen()` gibt immer `false` | `TestTransactionOpenFollowsBeginAndCommit`, `TestRunNoConfirmationInsideOpenTransaction` |
  | S1 | Store-Tier (`make test-replication`): gemeldete Position `+1 GiB` | `TestStreamIdleConfirmationKeepsOpenTransactionDeliverable` („kein CaptureCommand innerhalb 20 s“), `TestStreamRestartsOnExistingSlot` (neben den Unit-Tests von M3) |

- **Nullprobe der E2E-Phase (Mutation `true ||` vor der Bedingung in `confirmIdle`, danach `make image`
  und `make test-integration`):** Exit 2 nach 5 min 8 s, gedruckt „run-integration-tests: Leerlauf-
  Bestätigung — Run 07add76f-f7b3-423f-ad1d-760f86a06a39 endete interrupted statt completed:“. Alle
  Phasen davor liefen grün; die Phase bindet an die Leerlauf-Bestätigung. Mutation danach zurückgenommen,
  `make image` auf dem sauberen Baum neu gefahren.
- **Store-Experiment zu F-1** (Wegwerf-Test im Paket `receive`, danach gelöscht; Mutation M1 aktiv; Instanz
  mit Standard-`wal_sender_timeout`): der Stand-in-`Capture` der ersten Transaktion blockiert, in der Zeit
  committet eine zweite Transaktion 400.000 Änderungen auf die veröffentlichte Tabelle; nach 55 s
  Freigabe. Gedruckt (Lauf 5): „idle P=210904992 open=true confirmed_flush before cancel=210904992
  P=210904992“, `confirmed_flush_lsn` nach dem Abbruch `0/C9227A0`; der Neustart lieferte „commit=210904992
  changes=400000“. Ein Keepalive inmitten der Transaktion auf der Leitung tritt also auf (Zustand
  `open=true` am Assembler), seine Position ist die Commit-LSN der laufenden Transaktion, und die Quelle
  liefert sie nach der Bestätigung dieser Position dennoch vollständig.
- **Flake-Wiederholung** (gepinnte Images, frische Standard-Instanz, direkt nach den X1-Tests): X1
  „Last 53390256 B WAL, Rückstand 0 B nach 259ms, 6 Leerlauf-Bestätigungen“, Nullprobe „Last 53390248 B WAL,
  Rückstand 53390584 B nach 8.065s“; `TestWALRetentionThresholdEndToEnd` mit `-count=10`: alle zehn
  `--- PASS` (0,98 s bis 1,07 s).
- **Bench:** `tools/bench-backfill.sh` einzeln, Lauf `20260925T043056Z`, Exit 0: in allen neun Runs
  „WAL-Rückstand des Slots: Spitze 0 MiB im Run, 0 MiB unmittelbar danach“; gedruckt „WAL-Rückstand der
  größten Stufe (200000 Zeilen) — höchste Spitze im Run 0 MiB, unter der Warnschwelle von 100 MiB
  (Kennung der Schwelle in der gedruckten Zeile nicht wiedergegeben)“; vom Slot gehaltenes WAL Median 7/35/141 MiB (Stufen 10.000/50.000/200.000); Durchsatz Median
  der Stufe 200.000: 4504 Zeilen/s, gedruckt „Richtgröße (abgeleitet) — … = 2702400 Zeilen, abgerundet …:
  2000000 Zeilen“ (siehe F-4).
- **Suchläufe des Plans** an beiden Ständen nachgefahren (Parent `f4e32fba` per `git grep … f4e32fba`,
  Diff-Stand `427f6d1b`; alle sieben Zeilen des Feldes): siehe F-5 und Negativbefunde.
- **Umgebung:** `free -m` vor dem Tier-Lauf 18,4 GB verfügbar; dangling Volumes
  (`docker volume ls -q -f dangling=true | wc -l`) 34 vor und 34 nach allen Läufen; kein `prune`; eigene
  Wegwerf-Container mit `docker rm -fv` entfernt (Runner-Trap bzw. eigenes Skript). Parallel laufende
  fremde Container (`bats/bats`) auf dem Host, nicht von diesem Lauf.
- **Nicht gefahren (Grenze):** die PostgreSQL-17-Legs lokal (die CI-Legs beider Versionen sind grün);
  ein zweiter E2E-Lauf ohne Mutation (die CI zu `427f6d1b` trägt ihn); Mutationen der E2E-Assertions über
  die Nullprobe hinaus.

---

## Findings

### F-1 — ADR-Zeile „Store: Mutation inmitten der Transaktion → Change fehlt“ ist nicht erfüllbar; der geforderte Versuch eines Store-Tests ist nicht ausgeführt

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) Fitness Function (Zeile
  „reale PostgreSQL“) und Festlegung 1 (Begründung „ServerWALEnd liegt dann hinter Nachrichten, die noch nicht
  gespeichert sind“) · Plan §6 zweites Risiko („der Implementer versucht einen deterministischen Store-Test“)
  · Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“ (Abgrenzung)
- `pfad`: `docs/plan/planning/in-progress/slice-backfill-slot-leerlauf-bestaetigung.md:384-400` (§6);
  `internal/adapters/driving/replication/receive/receive.go:475` (`confirmIdle`)
- `befund`: Der Plan-Bericht nennt die Store-Bindung der Bedingung „keine offene Transaktion“ „nicht gebaut“
  (hergeleitet, nicht ausgeführt); der verlangte Versuch mit langsamem Stand-in ist nicht belegt. Gemessen (siehe
  Experiment oben): der Zustand ist im Store-Tier herstellbar (blockierender `Capture`, große Transaktion,
  Standard-`wal_sender_timeout`, etwa 55 s), das Keepalive inmitten der Transaktion trägt aber eine Position
  gleich der Commit-LSN dieser Transaktion und der Neustart liefert sie trotz bestätigter Position — die
  Mutation lässt sich am Store nicht rot färben, die ADR-Zeile ist wie geschrieben nicht erfüllbar. Die
  Bedingung hält die Invariante `ACK(position) => durable(all changes <= position)` formal, die Datensicherheit
  hängt am Verhalten der Quelle (Keepalive-Position ist der Dekodier-Stand, Wiederlieferung bei Gleichheit),
  das kein Test der Repo-Tiers bindet außer dem Sicherheits-Test der Größe der Position (S1 rot).
  Die Benennung „weiter offen, Adresse Architect“ ist tragfähig; sie braucht die gemessene Aussage statt der
  Herleitung.
- `verifizierbar`: ja — Wegwerf-Test wie beschrieben, Mutation M1 aktiv, Ausgabe der Ereignisse.
- `klasse`: Fitness-Function-Zeile nicht erfüllbar wie geschrieben (`BEO-PGC/fitness-function-gegen-eigene-entscheidung`)

### F-2 — `TestWALRetentionThresholdEndToEnd` fährt weder Stream noch `bootstrap.Run`; kein committeter Test bindet „Fehlerschwelle beendet den Prozess“

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md) (b) ·
  [`LH-QA-REL-003`](../../spec/lastenheft.md) · Plan §3 (Zeile „walretention_endtoend_internal_test.go“)
- `pfad`: `internal/bootstrap/walretention_endtoend_internal_test.go:34-140`
- `befund`: Der Vorgänger (`package bootstrap_test`) trieb `bootstrap.Run` mit realem Stream und beobachtete
  das Ende des Laufs; der Ersatz wächst an einem Slot ohne Stream und prüft `runWALRetentionCheck` und
  `mergeStreamAndWALFaultOutcome(nil, …)` direkt, `stopStream` ist eine bloße Abbruchfunktion. Name („EndToEnd“)
  und Meldungen („der Stream-Lauf wurde … nicht beendet“) beschreiben einen Stream, den es im Test nicht gibt.
  Die Kette „Schwelle überschritten → Stream endet → `Run` gibt `replication` zurück → Ausgang 1“ trägt danach
  nur die Herleitung des Plans (Priorität in `mergeStreamAndWALFaultOutcome`); die neue E2E-Phase belegt die
  Gegenseite (Schwelle nicht erreicht), die Nullprobe (Ende bei ausgeschalteter Bestätigung) ist eine Mutation und
  nicht committet — gemessen im Lauf oben, aber ohne Träger im Baum.
- `verifizierbar`: ja — `git grep -n "bootstrap.Run(" internal/bootstrap` gegen die Tests der Schwelle; die
  Nullprobe (Mutation `true ||`) färbt die E2E-Phase, nicht einen Test dieser Kette.
- `klasse`: Test-Zusage unter gleichem Namen abgeschwächt

### F-3 — Datei ersetzt statt verschoben: Umbenennung und Neufassung in einem Commit

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.3 (git mv + Inhaltsänderung = zwei Commits)
- `pfad`: `internal/bootstrap/walretention_endtoend_test.go` → `internal/bootstrap/walretention_endtoend_internal_test.go`
  (Commit `00caec48`)
- `befund`: `git log --name-status -M` zeigt `A` und `D` statt `R` (Ähnlichkeit unter der Schwelle: Paket
  `bootstrap_test` → `bootstrap`, 249 → 150 Zeilen); der Name des Tests bleibt, seine Historie reißt bei
  `git log --follow`. Ein reiner `git mv` vorab hätte sie gehalten. Einordnung: es ist keine Verschiebung im
  engen Sinn, sondern ein Ersatz unter neuem Dateinamen — deshalb LOW, obwohl der Wortlaut der Regel berührt ist.
- `verifizierbar`: ja — `git log --format=%h --name-status -M f4e32fba..427f6d1b -- internal/bootstrap`.
- `klasse`: Move und Inhaltsänderung in einem Commit

### F-4 — Handbuch: Lauf `20260925T032925Z` als „gemessen, gedruckt im Lauf“ ohne auflösbaren Träger; Richtgröß-Spanne wird von der eigenen Nachmessung verlassen

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz A · [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
  · Reviewer-Skill „Zahl im Träger … gegen die Messung driftend“ (Kriterien der HIGH-Klasse nicht erfüllt:
  Ursprung und Lauf-Kennung stehen, die Kernzahlen reproduzieren)
- `pfad`: `docs/user/benutzerhandbuch.md:1606-1610` (Richtgröße), `:1659-1690` (WAL-Rückstand, gehaltenes WAL)
- `befund`: Der Lauf steht nirgends im Repository außer im Handbuch; andere nicht auflösbare Läufe derselben
  Stelle tragen „übernommen … im Repository nicht auflösbar“, dieser „gemessen, gedruckt im Lauf“. Nachgemessen
  (Lauf `20260925T043056Z`): 0 MiB Spitze in neun Runs (stimmt), gehaltenes WAL 7/35/141 MiB gegen 7/31/140 MiB
  (Handbuch), Durchsatz 4504 Zeilen/s gegen 5532 → Richtgröße 2.000.000 gegen die genannte Spanne „3.000.000 bis
  5.000.000“; der Code-Wert (`estimatedRowsGuideline` in `internal/application/usecase/backfill/warn.go`)
  ist 4.000.000 und mit der Handbuch-Aussage („innerhalb der Spanne dieser Läufe“) widerspruchsfrei. Der Host
  trug bei der Nachmessung fremde Last (acht `bats`-Container), ein Zusammenhang mit dem Slice ist nicht belegt.
- `verifizierbar`: ja — `tools/bench-backfill.sh` einzeln.
- `klasse`: Zahl im Träger: Lauf nicht auflösbar

### F-5 — Suchlauf-Feld im Plan: Diff-Stand „92 Zeilen“ driftet gegen den Stand `427f6d1b`

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz A und §3.13 · Reviewer-Skill „Zahl im Träger … gegen die Messung
  driftend“
- `pfad`: `docs/plan/planning/in-progress/slice-backfill-slot-leerlauf-bestaetigung.md:300` (erste Zeile des Feldes)
- `befund`: Das Feld nennt für den Diff-Stand „92 Zeilen in 27 Dateien“ (Pathspec mit den fünf Ausschlüssen).
  Gemessen mit derselben Form: Parent `f4e32fba` 64 Zeilen in 26 Dateien (stimmt), `427f6d1b` 94 Zeilen in 27
  Dateien — die Plan-Datei ist im Suchraum (`in-progress/` ist nicht ausgeschlossen) und ihre späteren Edits
  tragen zwei weitere Treffer. Die Behandlung bleibt richtig; die Zahl gehört zu einem Zwischenstand. Die
  übrigen sechs Zeilen des Feldes stimmen an beiden Ständen (12→8 Zeilen bzw. „ohne die Plan-Datei 2“; 19→19
  Dateien; 78→81; 4→4; die `wiring.go`-Lokatoren `:1287.4,1288.1` und `:1183.5,1184.13` im Profil nachgemessen).
- `verifizierbar`: ja — `git grep -n -i -E 'bestätigten Position|lastAcked|Keepalive-Antwort|Keepalive-Position' <Stand> -- . ':!docs/reviews' ':!docs/plan/adr/0120-capture-slot-leerlauf-bestaetigung.md' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' | wc -l`.
- `klasse`: Zahl im Träger gegen den Diff-Stand gedriftet

### F-6 — Testkommentar mit Nebenklausel im Konjunktiv über den nicht gewählten Aufbau

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 · Reviewer-Skill „Kommentar trägt keine der Kommentar-Klassen“
- `pfad`: `internal/bootstrap/walretention_endtoend_internal_test.go:27-30`
- `befund`: „ein gleichzeitiger Schreiber eines anderen Test-Pakets läge sonst in derselben Größenordnung wie
  die Schwellen“ — der Hauptsatz (die Instanz gehört dem Test allein) ist Kopplung und trägt die Stelle, die
  Nebenklausel beschreibt die verworfene Alternative. Kein Produktionscode, keine Slice-/Wellen-Nummer.
- `verifizierbar`: nein — kein Gate liest Kommentar-Sprache.
- `klasse`: Kommentar trägt keine Kommentar-Klasse

### F-7 — Gedeckte Zahl im Sensor-Träger: 2136 gegen 2135 beim Nachmessen

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A (bewegliche Zahlen: die gedeckte Zahl ist der Beleg eines Laufs)
- `pfad`: `harness/sensors/coverage-gate.md:65-72`
- `befund`: Das Profil des eigenen Laufs (dedupliziert über die Block-Position) ergibt 2135 von 2567 = 83,17 %
  (Nenner 2567 stimmt, gedruckt in beiden Läufen „83.2%“); der Träger nennt 2136 = 83,21 %. Die Zahl hängt am
  Lauf und ist als solche gestempelt; die Differenz von einem Statement ist laufabhängig. Die DB-Adapter-Zahl
  (873 von 1058, 82,51 %) und die drei Produktionsdateien (47/47, 168/175, 237/604) stimmen.
- `verifizierbar`: ja — Profil aus der Stufe `coverage`, Auszählung wie im Träger beschrieben.
- `klasse`: Zahl im Träger gegen den Lauf (laufabhängig)

### F-8 — `lastAcked` startet bei 0; die erste Leerlauf-Bestätigung nach einem Neustart kann unter der Slot-Position liegen

- `kategorie`: INFO
- `quelle`: [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung 1 Punkt 4 (die
  Bestätigung geht nie zurück)
- `pfad`: `internal/adapters/driving/replication/receive/receive.go:183-191` (`newStreamOnSession`),
  `:475-478`
- `befund`: `newStreamOnSession` setzt `lastAcked` nicht auf die Startposition des Slots; nach einem Neustart
  kann ein frühes Keepalive mit einer Position hinter dem Dekodier-Stand der Wiederaufnahme, aber unter dem
  `confirmed_flush_lsn` des Slots, bestätigt werden. Ein Rückwärtsschritt des Slots ist nicht gemessen (die
  Quelle verwirft ihn nach meiner Kenntnis, der Quelltext ist nicht gelesen); der Store-Test
  `TestStreamRestartsOnExistingSlot` läuft grün. Höchstens eine Wiederlieferung, kein Verlust.
- `verifizierbar`: nein — kein Test führt einen Neustart mit hohem `confirmed_flush_lsn` und frühem Keepalive.
- `klasse`: Rückschritt-Wache ohne Startposition

### F-9 — Laufzeit von `make test-replication` und der E2E-Läufe

- `kategorie`: INFO
- `quelle`: Plan §6 (Laufzeit von `make test-integration` und `make test-replication`)
- `pfad`: `tools/harness/run-replication-tests.sh`, `.github/workflows/e2e.yml` (Lauf)
- `befund`: `make test-replication` 92 s (Lauf oben); die Job-Dauer der E2E-Legs stieg von 11 min 21 s/11 min
  2 s (PostgreSQL 17/18, Parent) auf 12 min 36 s/13 min 0 s (`427f6d1b`), also etwa +1,3 bis +2 min; die
  Grenze `timeout-minutes: 60` ist nicht berührt.
- `verifizierbar`: ja — `gh run view` beider Läufe.
- `klasse`: Laufzeit des Tiers (Beobachtung)

## Negativbefunde

- geprüft, ohne Befund: `internal/adapters/driving/replication/receive` (Sicherheit) — (a) `START_REPLICATION`
  trägt `proto_version '1'` und keine Option `streaming`, große Transaktionen laufen nicht vorab auf die
  Leitung; (b) `Begin` öffnet, `Commit` schließt den `Assembler`-Zustand, `Relation`-Nachrichten stehen bei
  `proto_version 1` zwischen `BEGIN` und `COMMIT`; `TransactionOpen` und `Consume` laufen in derselben
  Goroutine (Race-Detector in `make test` grün); (c) Keepalive und Commit stehen im selben Strom, `confirmIdle`
  läuft synchron zwischen zwei Nachrichten, kein `Capture` ist offen; die gemessene Position eines Keepalive
  inmitten der Transaktion ist die Commit-LSN; (d) Fehler von `Acknowledge` mitten im Stream endet in
  `ErrReplication` mit der Ursache dahinter, `lastAcked` bleibt (`TestRunIdleConfirmationFailureEndsAsReplication`,
  M7); (e) Neustart: `TestStreamIdleConfirmationKeepsOpenTransactionDeliverable` (S1 rot bei falscher Position)
  und das Experiment zu F-1 (Change vor P begonnen, nach P committet, vollständig geliefert).
- geprüft, ohne Befund: `internal/application/usecase/capture` und `port/inbound/idleconfirmation.go` —
  `ConfirmIdle` ruft ausschließlich `ReplicationAckPort.Acknowledge` (M6), leere Position erreicht den Port nicht,
  Ports nach [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md) geschnitten, `CaptureInboundPort` und seine
  Stand-ins unberührt; `make a-check` ohne Befund.
- geprüft, ohne Befund: `internal/adapters/driving/replication/mapper`, `internal/bootstrap/wiring.go` — additive
  lesende Methode, eine `captureService`-Instanz an beiden Eingängen; `Run` endet ohne den zweiten Port in der Klasse
  `configuration` (`TestRunWithoutIdleConfirmationIsConfigurationError`).
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` — Compose-Override im Temp-Verzeichnis, `compose.yaml`
  unverändert, Bind-Mount und `CDC_CONFIG_FILE`, Wiederherstellung am Phasen-Ende, `rm -rf` im `trap` (leerer
  Wert harmlos), Startzeitpunkt-Prüfung gegen Neustart, Last wird gemessen und bricht bei Unterschreitung ab;
  Nullprobe rot (siehe oben).
- geprüft, ohne Befund: `tools/harness/run-replication-tests.sh` — zweite Instanz mit Standardwert, Isolation
  (`CDC_REPLICATION_TEST_STANDARD_DSN` nur im Paket `receive`), `--- PASS`-Prüfung des Schwellen-Tests, Exit-Codes
  ungepiped (§3.9), `docker rm -fv` und Netz im `trap`; `TestWALRetentionThresholdEndToEnd` 10× grün.
- geprüft, ohne Befund: `tools/bench-backfill.sh` und `harness/targets/bench-backfill.md` — `settle_wal` ohne
  Schreibzugriff, Zeile „WAL-Rückstand-Schwellen“ entfallen, neue Vergleichszeile ohne Pass/Fail, Exit 0
  (Lauf oben).
- geprüft, ohne Befund: `spec/pflichtenheft.md`, `spec/architecture.md` — Sätze zu `LH-QA-REL-001.a` und
  `SPEC-009` ohne ADR-/Slice-Bezug (Stratum präzisiert, nichts Neues bindend), Sicht sprach- und meilensteinfrei
  (§3.4), Änderungshistorie fortgeschrieben.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Version 1.56 und Zeile der Änderungshistorie im selben
  Diff, keine neue Betreiber-Oberfläche (keine `CDC_*`-Variable, kein Endpunkt), die drei Stellen des Verdikts
  nachgezogen, die entfallenen Sätze („kann bei weniger Zeilen greifen“) an beiden Stellen weg (Ausnahme F-4).
- geprüft, ohne Befund: `harness/README.md`, `harness/sensors/*.md`, `docs/user/e2e-abdeckung.md` — die drei
  Zeilen tragen, was gefahren wurde; Lokatoren und Nenner (2567, 1058) nachgemessen (Ausnahme F-7).
- geprüft, ohne Befund: Hard Rules — §3.1 Docker-only (nur `docker`/`make`), §3.2 kein `nolint` im Diff, §3.5
  `Accepted`-ADRs unberührt (Zitate in `ADR-0080`/`ADR-0049` bleiben), §3.7 in Produktionscode (keine Slice-/Wellen-
  Chronik, keine Konjunktiv-Kommentare; Ausnahme F-6 im Test), §3.11 kein host-lokaler Pfad im Diff; Commit-
  Betreffs tragen `LH-*`/`ADR-*` und keine `SPEC-*`/`ARC-*`; die Lifecycle-Moves sind reine `git mv` (Ausnahme F-3).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Fitness-Function-Zeile nicht erfüllbar wie geschrieben · Test-Zusage unter
gleichem Namen abgeschwächt · Move und Inhaltsänderung in einem Commit · Zahl im Träger: Lauf nicht auflösbar ·
Zahl im Träger gegen den Diff-Stand gedriftet · Kommentar trägt keine Kommentar-Klasse · Zahl im Träger gegen den
Lauf (laufabhängig) · Rückschritt-Wache ohne Startposition · Laufzeit des Tiers (Beobachtung)

## Verdikt

**Merge-blockierend:** nein für die Sicherheit des Capture-Pfads — kein HIGH, die Leerlauf-Bestätigung
überspringt nach Messung keine Change (Unit-, Store- und E2E-Mutationen färben je die tragenden Belege rot; die
Herleitung „P ist die Dekodier-Position, Wiederlieferung bei Gleichheit“ ist im Store-Tier am realen Verhalten
belegt). Ja für die Schließung des Slice: zwei MEDIUM — F-1 geht als Architect-Frage (Fitness-Function-Zeile
und Aussage „hinter noch nicht gespeicherten Nachrichten“ gegen die gemessene Lage, ggf. neue ADR mit
`Supersedes`), F-2 an den Implementer (Bindung der Kette „Schwelle → Prozessende“ oder benannter Aufschub mit
Adresse, Name und Meldungen des Tests). Die LOW-Findings zieht der Implementer in derselben Runde nach. Die
DoD-Zeile „Review durchgeführt“ bleibt offen, weil eine Fixrunde folgt.

**Übergabe:** Findings gehen an den Implementer (F-2 bis F-6) und den Architect (F-1); die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg (dieser Diff,
dieser Skill, dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11).
