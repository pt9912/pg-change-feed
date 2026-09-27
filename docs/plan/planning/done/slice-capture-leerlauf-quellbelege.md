# Slice capture-leerlauf-quellbelege: Quellbelege der Capture-Kette — Keepalive inmitten einer Transaktion an PostgreSQL 17 und 18, Fehlerschwelle des WAL-Rückstands beendet den Container

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er geht `slice-transformationen-e2e-abhilfe`
voraus (Start-Trigger dort, §4;
[welle-transformationen](../welle-transformationen.md) §5).

**Bezug:** [`LH-QA-REL-001`](../../../../spec/lastenheft.md) (kein
Datenverlust, Persist-before-ACK),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (WAL-Rückstand sichtbar),
[`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill; die
Leerlauf-Bestätigung entstand an ihm),
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) und
[`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
(Festlegung 2 und ihr Trigger „Ein committeter Test der Quellseite entsteht“),
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) (Schwellen
des WAL-Rückstands), [`ADR-0030`](../../adr/0030-testpyramide.md)
(Testpyramide), Architect-Verdikt
[`architect-verdict-welle-backfill-bestand-lese-schritt`](../../../reviews/architect-verdict-welle-backfill-bestand-lese-schritt.md)
§5 (b) und (e).

**Berührte Spec-Stellen:** [`SPEC-012`](../../../../spec/pflichtenheft.md)
(PostgreSQL 17 und 18 unterstützt),
[`SPEC-013`](../../../../spec/pflichtenheft.md) (Schwellen des
WAL-Rückstands) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-27.

**Autor:** Planner-Agent, Closure der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Zwei Aussagen der Capture-Kette, die an einer einmaligen Messung
des Reviewers hängen, tragen einen committeten Beleg. (1) Ein Keepalive inmitten
einer Quelltransaktion trägt die Commit-LSN dieser Transaktion, und die Quelle
liefert sie nach der Bestätigung vollständig — an PostgreSQL 17 **und** 18, im
Tier `make test-replication`. (2) Erreicht der WAL-Rückstand die Fehlerschwelle,
endet der Container — ein Beleg in `make test-integration` am realen Stream, der
die Verdrahtung in `Run` (`streamCtx`/`stopStream`,
`mergeStreamAndWALFaultOutcome`) trägt. Der Slice endet mit einer Ergänzung von
[`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
Festlegung 2 durch den Architect.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung an Produktionscode der Leerlauf-Bestätigung oder der
  Schwellen-Prüfung.** Der Slice belegt, er ändert nicht; ergibt ein Beleg eine
  Abweichung, ist sie ein Befund für den Architect, kein stiller Umbau.
- **Der Aufbau der Kette „Rückstand wächst, weil die Persistierung hält“ als
  Produktfunktion.** Der Beleg nutzt einen Testaufbau; ein Betreiber-Weg, die
  Persistierung anzuhalten, entsteht nicht.
- **Die Klassen- und Wiederholungsfrage bei `transient`** —
  `slice-capture-transient-wiederholung`; die Schwellen-Kette ist eine andere
  Ursache.
- **Ein Mutations-Harness und ein Gate über die Laufzeit.** Der Beleg ist ein
  Test im bestehenden Tier; er ändert weder Gate noch Workflow-Struktur
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht: die Läufe stehen in
  bestehenden Phasen von `e2e.yml`).

## 2. Definition of Done

- [x] Der Keepalive-Beleg steht als committeter Test im Tier
      `make test-replication` (Phase `tier`): ein **roher Protokoll-Client**
      (kein Stream-Adapter, `proto_version 1`) liest ab `START_REPLICATION` 38 s
      nichts und antwortet nichts, während **eine** Transaktion über 400.000
      Änderungen auf die veröffentlichte Tabelle committet; der Keepalive tritt
      zwischen BEGIN und COMMIT dieser Transaktion auf, der Client bestätigt
      dessen `ServerWALEnd` (gleich der Commit-LSN), beendet die Verbindung und
      startet den Stream ab `confirmed_flush_lsn` neu, und der Neustart liefert
      die Transaktion vollständig. Der Test läuft auf der vorhandenen
      Standard-Instanz `CDC_REPLICATION_TEST_STANDARD_DSN` aus
      `tools/harness/run-replication-tests.sh` (Standardwert von
      `wal_sender_timeout`, kein fremder Schreiber). *Zu belegen durch:* der
      Verifikations-Report
      [`verifikation-slice-capture-leerlauf-quellbelege`](../../../reviews/verifikation-slice-capture-leerlauf-quellbelege.md)
      §1 und §4 K1 — im Lauf des Verifiers **gemessen**: PostgreSQL 18.6
      „Keepalive inmitten der Transaktion nach 11080 von 400000 Änderungen,
      ServerWALEnd 0/12786CB0 gleich Commit-LSN 0/12786CB0“, Neustart lieferte
      400000 Änderungen (40,32 s); PostgreSQL 17.11 (Digest aus `e2e.yml`, über
      `PG_TEST_IMAGE`) nach 11077 Änderungen, 0/124C79A8 (40,37 s); die Mutation
      „bestätigte Position `+ (1 << 30)`“ färbt den Test an beiden rot (je
      78,18 s). In CI je Leg (Lauf 36287009221, Job-Logs mit
      `gh api repos/pt9912/pg-change-feed/actions/jobs/<Job>/logs` am 2026-09-27
      vom Planner **gemessen**): Job 108529548457 (Leg 17.11) nach 11538
      Änderungen, 0/BED26C8 gleich Commit-LSN, `PASS` in 40,72 s; Job
      108531740887 (Leg 18.6, zweiter Versuch) nach 11541 Änderungen, 0/C184090,
      `PASS` in 43,20 s. Die Zahl „etwa 11.000 Änderungen bis zum Keepalive“ ist
      Messung, kein Vertrag.
- [x] Der Beleg „Fehlerschwelle erreicht → Container endet“ steht als Runner-Phase
      „Fehlerschwelle beendet den Container“ in `make test-integration`: die
      Fehlerschwelle wird über den im Runner vorhandenen Compose-Override
      (`wal_retention_error_bytes` klein, Phase „Leerlauf-Bestätigung“ in
      `tools/harness/run-integration-tests.sh`) gesenkt, die Persistierung wird
      gehalten (eine Sitzung des Runners sperrt `cdc.change` exklusiv, die
      Persistierung des Streams wartet an der Sperre), und bei WAL ohne Inhalt für
      die Publication über der Fehlerschwelle endet der Feed-Container mit Ausgang
      1, das Log trägt die Abbruch-Zeile mit einem Rückstand über der
      Fehlerschwelle, und `cdc.process_heartbeat` trägt einen Fehlerzustand. Die
      **Klasse** des Ausgangs ist nicht Teil der Zusage der Phase: sie ist
      `storage` (Ausgabezeile der Phase; **gemessen** im Lauf des Reviewers, im Lauf
      des Verifiers und in den zwei CI-Läufen des Ankers unten — je Ende 2 bis 3 s
      nach der Last; fünf Läufe des Implementers **übernommen** aus dessen
      Bericht), `ADR-0049` legt `replication` fest — Codefehler, Träger
      `slice-wal-fehlerschwelle-ausgangsklasse` (Architect-Verdikt
      [`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md)
      §2 und §4). **Der Ausgang „Grenze bleibt“ ist zulässig** und hier
      eingetreten: die Klasse steht als benannte Grenze mit dem Messergebnis im
      Bericht und als benannter Text in `harness/README.md` §Sensors bei `make
      test-integration`. *Zu belegen durch:* ein realer, grüner `make
      test-integration`-Lauf mit der Abdeckungs-Zeile in
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md) (Erzeugnis
      des Runners; Lauf des Verifiers: Exit 0, 373 s, „E2E-Abdeckungstabelle
      unverändert“) und die benannte Grenze in `harness/README.md` §Sensors; die
      Mutation „Aufruf von `stopStream` in der Schwellen-Prüfung entfernt“ färbt
      die Phase rot — **erprobt** vom Verifier
      ([Verifikations-Report](../../../reviews/verifikation-slice-capture-leerlauf-quellbelege.md)
      §4 P2): Stelle `runWALRetentionCheck` in `internal/bootstrap/wiring.go`, die
      Zeile `stopStream()` zu `_ = stopStream`; Instanz Kopie des Repos, `make image`
      aus der Kopie, voller `make test-integration` dort; Farbe rot: „der
      Feed-Container lief 90 s nach einer Last von 16031584 B WAL über der
      Fehlerschwelle 8388608 B weiter“. In CI je Leg (Lauf 36287009221, Job-Logs
      vom Planner **gemessen**): Job 108529548457 (Leg 17.11) und Job 108531740887
      (Leg 18.6, zweiter Versuch) tragen die Ausgabezeile der Phase, Ende 3 s nach
      der Last, Klasse `storage`.
- [x] Die Ergänzung von
      [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
      Festlegung 2 liegt vor:
      [`ADR-0129`](../../adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md)
      (Commit `4e654153`) ist eine neue ADR des Architects mit teilweisem
      `Supersedes` von `ADR-0121` (Festlegung 2 und Konsequenz-Punkt „Negativ
      (Grenze, benannt)“), nennt das Tier des Keepalive-Tests und ersetzt die Grenze
      „für PostgreSQL 17 liegt die Messung nicht vor“ mit dem gemessenen Ergebnis
      der zwei Pins; jede ihrer Aussagen über eine Menge trägt den Beleg-Anker dieses
      Slice ([`AGENTS.md`](../../../../AGENTS.md) §3.12 „Verfasser einer ADR“). *Zu
      belegen durch:* die ADR und ihre Index-Zeile in
      [`docs/plan/adr/README.md`](../../adr/README.md) (beide liegen vor) und die
      Lese-Prüfung der ADR gegen den Verfasser-Satz durch den Reviewer (frischer
      Kontext, Diff `e4b77a05..4e654153`): Review-Report
      [`review-adr-0129-capture-quellseite-keepalive`](../../../reviews/review-adr-0129-capture-quellseite-keepalive.md)
      (Commit `027533aa`; 0 HIGH · 1 MEDIUM · 2 LOW · 2 INFO, aus dem Report
      **übernommen**) — Menge, gedruckte Zeilen, Mutationszuordnung und Grenzen der
      ADR halten der Nachmessung des Reviewers stand, keine Fixrunde an der ADR;
      die Ränder F-1 bis F-3 stehen als Lesehilfe in §7 (Risiko §6, letzte Zeile).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9): Lauf des
      Planners am Stand `3f107381` (Inhalts-Commit der Closure-Notiz), Exit 0,
      „coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%“, „d-check: 1326
      Datei(en) geprüft, 0 Befund(e)“, „gesamt: 0 Befund(e)“ (a-check); der Lauf des
      Verifiers am Stand `65a63968` steht in dessen Report §1.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8):
      [`review-slice-capture-leerlauf-quellbelege`](../../../reviews/review-slice-capture-leerlauf-quellbelege.md)
      (0 HIGH, 1 MEDIUM, keine Fixrunde am Implementer).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` §Sensors (die Beschreibungen von `make
      test-replication` und `make test-integration` nennen die neuen Belege bzw.
      die benannte Grenze); das Benutzerhandbuch bleibt unberührt (keine
      Betreiber-Oberfläche).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke) — §7 trägt ihn.
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure dieses Slice (§7) und zusätzlich der Closure der nächsten Welle
      (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten; ein Slice ohne Welle wird von ihr
      mitgeprüft).

**Umfang:** M — Schätzung, nicht gemessen: zwei Belege in zwei Tiers, davon
einer mit ungeklärter Stabilität, dazu die ADR-Ergänzung; die Zeit je Lauf des
Keepalive-Tests ist *hergeleitet* (der Wegwerf-Test des Reviewers wartete 55 s,
[`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
§Gemessen) auf etwa eine Minute je Leg.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/replication/receive/sourcekeepalive_test.go` (Ort am Start gelesen: das Paket der Stream-Tests, dessen Helfer `newTestEnvOn`/`readConfirmedFlush`/`awaitSlotInactive` der Test nutzt) | neu (Plan-Nachzug: ein **roher Protokoll-Client** statt des Stream-Adapters) | Der Adapter bestätigt inmitten einer Transaktion nicht (`TestRunNoConfirmationInsideOpenTransaction`, Produktionscode bleibt unberührt); der Beleg der Quellseite liest deshalb das Protokoll selbst: der Client wartet die Hälfte von `wal_sender_timeout` ab, liest den Keepalive zwischen BEGIN und COMMIT einer Transaktion über 400.000 Änderungen, bestätigt dessen `ServerWALEnd`, beendet die Verbindung und startet den Stream ab `confirmed_flush_lsn` neu. |
| `tools/harness/run-replication-tests.sh` | update | Der Test läuft auf der vorhandenen Standard-Instanz (Wahl am Start: sie trägt den Standardwert von `wal_sender_timeout` und im Tier-Lauf keinen fremden Schreiber), als eigener Lauf hinter dem Schwellen-Beleg mit `-run`-Muster und PASS-Wächter (eigene Umgebungsvariable `CDC_SOURCE_KEEPALIVE_TEST_DSN`, nur im Tier-Lauf gesetzt: die Phase `measure` und `go test ./...` überspringen ihn und zahlen die 38 s Wartezeit nicht). Das Tier läuft in CI: der Schritt „Replication-Tier (go test ./... und Slot-Reserve)“ in `.github/workflows/e2e.yml` führt `bash tools/harness/run-replication-tests.sh tier` je Matrix-Leg (PostgreSQL 17 und 18) hinter dem Compose-Integrationstest aus; die gedruckte Zeile des Tests steht im Job-Log beider Legs (DoD 1). |
| `tools/harness/run-integration-tests.sh` | update (Plan-Nachzug: eine **Bash-Phase** statt einer `func TestE2E*`) | Phase „Fehlerschwelle beendet den Container“ hinter der Phase „Leerlauf-Bestätigung“ (deren Compose-Override sie wiederholt), mit `abdeckung_declare`. Die Belegmittel (`docker inspect`, `docker logs`, Compose-Override, `pg_stat_activity`) sind die des Runners, die Gegenseite ist ebenfalls eine Runner-Phase; eine neue `func TestE2E*` entsteht nicht, kein `-run`-Muster ändert sich. |
| `test/integration/integration_test.go` | **entfällt** (Reduktion, Begründung in der Zeile darüber) | Kein Go-Test in diesem Paket. |
| `docs/user/e2e-abdeckung.md` | update (Erzeugnis des Runners) | Eine neue Zeile; die Ort-Zeilen der Runner-Phasen hinter der neuen Phase verschieben sich. |
| `harness/README.md` | update | Beschreibungen der beiden Läufe samt benannter Grenze des zweiten Belegs. |
| `docs/plan/planning/in-progress/slice-capture-leerlauf-quellbelege.md` | update | Dieser Nachzug, DoD-Haken, Suchlauf-Feld. |
| ADR-Ergänzung (Architect) und `docs/plan/adr/README.md` | neu / update — **nicht Teil des Implementer-Laufs** | Architect-Zug nach der Verifikation: [`ADR-0129`](../../adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md) (Commit `4e654153`) und ihre Index-Zeile. |
| Beobachtungs-Register (`observations/BEO-PGC/…`), `open/slice-code-kommentare-bereinigung.md`, `open/slice-wal-fehlerschwelle-ausgangsklasse.md`, `welle-transformationen.md` | update — Planner-Closure | Evidence-Dateien und `state.md` der acht berührten Einträge; Übergabe der zwei Godoc-Kommentare zum durch `ADR-0121` ersetzten Begründungssatz; Kommentar-Übergaben der Runner-Phase und bedingte Kante der Stabilisierung; Links auf Kennungen (§7). |

**Befund der Erprobung (für den Architect, keine Änderung am Produktionscode):**
der zweite Beleg läuft grün, seine Aussage über die **Klasse** des Ausgangs
weicht von [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
ab. Bei gehaltener Persistierung wartet die Persistierung des Streams in
`Capture`; die Schwellen-Prüfung setzt den WAL-Fehler der Klasse `replication`
und ruft `stopStream`, der abgebrochene Kontext lässt `Capture` mit einem Fehler
der Klasse `storage` („Persistenzfehler im ChangeStore: context canceled“)
zurückkehren, und `mergeStreamAndWALFaultOutcome` gibt einen Stream-Fehler jeder
Klasse vor dem WAL-Fehler zurück. Gemessen vom Implementer (fünf Läufe,
**übernommen** aus dessen Bericht: ein Diagnose-Lauf, drei Läufe der Phase in
einem Wegwerf-Aufbau des Runners und ein vollständiger `make test-integration`;
die Ausgabezeile der Phase nennt die Klasse; bestätigt in den Läufen des
Reviewers, des Verifiers und in beiden CI-Läufen, DoD 2): der Container endet
mit Ausgang 1 zwei Sekunden nach der Last (der Diagnose-Lauf: vier; in CI drei),
das Log trägt die Abbruch-Zeile („WAL-Rückstand über
Fehlerschwelle — kontrollierter Abbruch“) und danach „Fehlerklasse storage“,
`cdc.process_heartbeat` trägt in allen fünf Läufen `storage`. Die
Klasse `replication` erreicht `Run` nur, wenn der Stream-Lauf regulär endet
(Kontext-Abbruch in `ReceiveMessage`). Die Phase belegt deshalb Ende, Ausgang,
Abbruch-Zeile und sichtbaren Fehlerzustand, nicht die Klasse; das Handbuch
([`benutzerhandbuch.md`](../../../user/benutzerhandbuch.md), Abschnitt
„Bestand als Backfill überführen“, Punkt „WAL-Rückstand des Capture-Slots“:
„beendet sich der Feed-Container mit der Klasse `replication`“) und die
Kommentare an `mergeStreamAndWALFaultOutcome` und in
`walretention_slotgrowth_internal_test.go` nennen die Klasse als Zusage. Die
Abweichung ist ein Codefehler, kein Fehler des Textes (Architect-Verdikt
[`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md)
§2): die Träger geben `ADR-0049` richtig wieder, der Stream-Fehler nach
`stopStream` verdeckt den WAL-Fehler. Die Korrektur trägt
`slice-wal-fehlerschwelle-ausgangsklasse`; mit ihm trägt die Phase die Klasse als
Zusage.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „welches Tier
belegt die Position eines Keepalive inmitten einer Transaktion“ und „welcher
Test trägt die Kette Fehlerschwelle → Prozessende“).** Suchraum: der ganze Baum
([`AGENTS.md`](../../../../AGENTS.md) §3.13 §Suchform), ausgenommen
`docs/reviews/**`, die Records unter `done/` und `.harness/baseline/**`; die
Plan-Datei schließt das Werkzeug aus. Die zwei Stände sind Commit-Kennungen: der
Parent `dbc4dbe4` (der Commit vor der ersten Änderung des Implementer-Laufs) und
`8611185b` (der letzte Commit vor der Closure, der einen Träger der Muster
ändert; danach ändert nur diese Plan-Datei den Baum). Muster 5 sucht die
Kommentar-Wendungen der Schwellen-Kette im Go-Code und steht deshalb auf
`-- internal`; die Doku-Träger der Klassen-Zusage trägt Muster 4 über den ganzen
Baum. Die Zeilen stehen im Format des Werkzeugs (`make suchlauf-nachmessen
PLAN=<diese Datei>`, Exit 0):

```suchlauf
dbc4dbe4 16 -n -E 'inmitten (der|einer) (Quell)?[Tt]ransaktion|Quellseite|einmalige Messung|SourceKeepalive|SOURCE_KEEPALIVE' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
8611185b 60 -n -E 'inmitten (der|einer) (Quell)?[Tt]ransaktion|Quellseite|einmalige Messung|SourceKeepalive|SOURCE_KEEPALIVE' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
dbc4dbe4 47 -n -E 'mergeStreamAndWALFaultOutcome|stopStream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
8611185b 58 -n -E 'mergeStreamAndWALFaultOutcome|stopStream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
dbc4dbe4 8 -n -E 'kontrollierter Abbruch|beendet sich der Feed-Container|Fehlerschwelle (beendet|erreicht)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
8611185b 26 -n -E 'kontrollierter Abbruch|beendet sich der Feed-Container|Fehlerschwelle (beendet|erreicht)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
dbc4dbe4 7 -n -E 'mit der Klasse .replication.' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
8611185b 11 -n -E 'mit der Klasse .replication.' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
dbc4dbe4 2 -n -E 'nur zum Zug|Fehler bei Stream-Ende|regulär endete' -- internal
8611185b 2 -n -E 'nur zum Zug|Fehler bei Stream-Ende|regulär endete' -- internal
```

| Träger | Befund (Parent → Stand `8611185b`, `-n`-Trefferzeilen) | Behandlung |
|---|---|---|
| Sätze, die die Keepalive-Messung „einmalig“ oder „PostgreSQL 17 nicht gemessen“ nennen | Muster 1: 16 → 60; die 44 neuen Treffer sind: `ADR-0129` (24) mit ihrer Index-Zeile in `docs/plan/adr/README.md` (1), dieser Lauf (Test-Datei 8, `run-replication-tests.sh` 8, `harness/README.md` 1) und der Nachzug des Planners (`state.md` von `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` +1, Übergabe in `slice-code-kommentare-bereinigung` 1). **Gefunden:** die alte Grenze steht in `ADR-0121` (9 Treffer: §Kontext, Festlegung 2, Konsequenz „Grenze, benannt“, Trigger; die Konsequenz und die Festlegung 2 ersetzt `ADR-0129`), `ADR-0120` (1, die Store-Zeile der Fitness Function), in den Dateien des Registers `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` (observation, evidence; die `state.md` trägt den Stand mit dem Test), in einer Evidence-Datei von `BEO-PGC/adr-aussage-breiter-als-ihre-messung` und in einer Zeile der Änderungshistorie von `spec/pflichtenheft.md` (Aussage der Regel „nie inmitten einer Quelltransaktion“, nicht der Messung); `seam_test.go` (1) trägt ein anderes Wort („Updates inmitten der Transaktion“). **Nicht gefunden:** kein Satz in `docs/user`, `harness`, `internal` (ohne die neue Test-Datei), `test` oder `tools`, der die Messung „einmalig“ nennt oder PostgreSQL 17 als nicht gemessen führt; die Wendung „PostgreSQL 17 ist nicht gemessen“ steht außerhalb der Records nur in `ADR-0119` und `ADR-0125` (anderer Gegenstand: Wirkung der Lesesperre bzw. d-migrate). | `ADR-0121` und `ADR-0120` bleiben unberührt (`Accepted`); die Ergänzung ist `ADR-0129`. Das Register `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` trägt den Ausgang „verkörpert“ in seinem `state.md` (Planner-Closure). Die zwei Godoc-Kommentare mit dem durch `ADR-0121` ersetzten Begründungssatz („das WAL-Ende liegt dann hinter Nachrichten, die noch nicht gespeichert sind“: `receive.go` `confirmIdle`, `seam_test.go`) sind an `slice-code-kommentare-bereinigung` übergeben (Tranchen T4 und T8; das Godoc von `confirmIdle` ist kein Kandidat des Werkzeugs `make kommentar-kennungen`). |
| Beschreibungen der Schwellen-Kette in Kommentaren und Doku | Muster 2 (Symbole): 47 → 58; die 11 neuen Treffer sind der Plan `slice-wal-fehlerschwelle-ausgangsklasse` (7), `slice-start-vorlauf-grenze` (1) und drei Register-Dateien (1 je). Muster 3 (Beschreibung): 8 → 26; die 18 neuen Treffer sind dieser Lauf (Runner-Phase 5, `harness/README.md` 1, `docs/user/e2e-abdeckung.md` 1), der Plan `slice-wal-fehlerschwelle-ausgangsklasse` (7), `welle-transformationen` (+1) und drei Register-Dateien. Muster 4 (Klassen-Zusage „mit der Klasse `replication`“): 7 → 11; die 4 neuen Treffer sind der Plan `slice-wal-fehlerschwelle-ausgangsklasse` (3, der Träger-Slice nennt die Klasse als sein Ziel) und eine Register-Evidence-Datei. **Gefunden mit Klassen-Zusage der Kette:** `docs/user/benutzerhandbuch.md` (Abschnitt „Bestand als Backfill überführen“, Punkt „WAL-Rückstand des Capture-Slots“: „beendet sich der Feed-Container mit der Klasse `replication` (Ausgang 1)“) und `internal/bootstrap/wiring.go` (Kommentar an `mergeStreamAndWALFaultOutcome`: „nur zum Zug, wenn der Stream-Lauf regulär endete“; Muster 5: 2 → 2, die zweite Stelle ist der Kommentar von `walretention_slotgrowth_internal_test.go`). Beide gelten nach dem Befund oben nicht für die gehaltene Persistierung. Die Treffer von Muster 4 in `ADR-0128`, `roadmap.md`, `welle-transformationen.md`, `slice-start-vorlauf-grenze` und zwei Register-Dateien beschreiben den Vorlauf vor `stream.Run` (gelesen: je ein Satz zur Kante von `slice-start-vorlauf-grenze` oder zu `ADR-0128`), einen anderen Gegenstand. **Nicht gefunden:** mit den Mustern 4 und 5 kein weiterer Träger der Kette in `internal`, `spec`, `harness` und `test`; die Norm (`ADR-0049` Folgepflicht, `SPEC-008`) ist der Gegenstand der Frage an den Architect, kein nachzuziehender Träger. | Gemeldet, nicht mitgeändert: das Handbuch liegt laut Plan außerhalb dieses Slice, `wiring.go` ist Produktionscode der Schwellen-Prüfung (§1). Entschieden ist „Code“ (Architect-Verdikt, siehe oben): das Handbuch bleibt wahr und unberührt; die zwei Kommentare (`wiring.go` an `mergeStreamAndWALFaultOutcome`, `walretention_slotgrowth_internal_test.go`) trägt `slice-wal-fehlerschwelle-ausgangsklasse` als Adresse (Änderungs-Tabelle dort, `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`); Frist: die Closure jenes Slice, der Planner der Closure zieht nach. Bis dahin trägt `harness/README.md` die Grenze benannt. |

## 4. Trigger

**Start** (`next` → `in-progress`): kein weiterer Slice in `in-progress/`
(WIP-Limit 1). Der Slice muss `done` sein, **bevor**
`slice-transformationen-e2e-abhilfe` startet (Start-Trigger dort): beide tragen
eine Container-Ende-Grenze im selben Runner
`tools/harness/run-integration-tests.sh`, und der Belegaufbau „Prozess endet“
entsteht einmal.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Keepalive-Test
  und E2E-Beleg nicht in einem Review tragen — der abtrennbare Teil ist der
  E2E-Beleg (zweiter Liefer-Punkt) als eigener Slice mit Start vor
  `slice-transformationen-e2e-abhilfe`.
- `in-progress` → `open` (blockiert): falls der Keepalive-Test an PostgreSQL 17
  ein anderes Verhalten zeigt als an 18 (Architect-Frage: eine Folge-ADR zur
  Bedingung der Leerlauf-Bestätigung, kein Anpassen des Tests an das Ergebnis).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer `make test-replication`-Lauf an
beiden PostgreSQL-Versionen + ein realer, grüner `make test-integration`-Lauf
(oder die benannte Grenze des zweiten Belegs) + die ADR-Ergänzung des Architects
liegt vor + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Keepalive-Test ist zeitabhängig** (Standard-`wal_sender_timeout` 60 s;
  ein Keepalive tritt nur innerhalb der Wartezeit auf). *Erwartet, zu belegen
  durch:* mehrere Läufe je Version ohne Ausfall; ein Ausfall ohne Ursache im
  Code ist Grund für die Rückführung nach `open/`. **Ausgang: entfallen.** Der
  Test lief in allen sechs Läufen mit gedruckter Dauer bis zum Ende grün:
  Reviewer 40,52 s (18.6) und 40,42 s (17.11), Verifier 40,32 s (18.6) und
  40,37 s (17.11) (je aus dem Report **übernommen**), CI 40,72 s (17.11, Job
  108529548457) und 43,20 s (18.6, Job 108531740887) (vom Planner **gemessen**);
  dazu zwei Mutationsläufe des Verifiers, die den Test grün ließen (§4 K6, K7).
  Ein Ausfall ohne Ursache im Code ist nicht aufgetreten. Ein Runner, auf dem
  400.000 Einfügungen länger als 38 s brauchen, endet laut mit „COMMIT … ohne
  Keepalive“ (Mutationen M3 und M7 des Reviewers, **übernommen**), nicht mit einer
  falschen Aussage. Der rote Versuch in CI (Lauf 36287009221, Leg 18) liegt in der
  Phase „Leerlauf-Bestätigung“; der Schritt dieses Tests lief dort nicht
  (`skipped`) und im Wiederholungsversuch grün (Register, Risiko unten).
- **PostgreSQL 17 liefert eine andere Position als 18.** *Erwartet, zu belegen
  durch:* der Lauf an beiden Versionen; eine Abweichung ist ein Befund für den
  Architect (Rückführung nach `open/`, §4). **Ausgang: entfallen.** An 17.11
  ist die Keepalive-Position gleich der Commit-LSN und der Neustart liefert
  400.000 Änderungen (Verifikations-Report §1, **übernommen**; CI Job 108529548457,
  vom Planner **gemessen**: `ServerWALEnd 0/BED26C8 gleich Commit-LSN 0/BED26C8`).
- **Der Aufbau „Persistierung halten“ ist instabil** (der erste Ansatz eines
  Tier-Tests scheiterte am `wal_sender_timeout` der Testinstanz,
  [`architect-verdict-welle-backfill-bestand-lese-schritt`](../../../reviews/architect-verdict-welle-backfill-bestand-lese-schritt.md)
  §5 (b)). *Erwartet, zu belegen durch:* wiederholte Läufe; der Ausgang „Grenze
  bleibt“ ist zulässig und steht mit dem Messergebnis im Bericht. **Ausgang:
  entfallen.** Die Phase endete in jedem Lauf, in dem sie lief, 2 bis 3 s nach der
  Last mit Ausgang 1: fünf Läufe des Implementers (**übernommen**), Reviewer und
  Verifier (je aus dem Report **übernommen**), zwei CI-Läufe (vom Planner
  **gemessen**, je 3 s). Ein Fehlschlag der Vorbedingungen endet als benannter
  Abbruch: die Mutation „Sperre nicht gesetzt“ endete an der Vorbedingung
  (Verifikations-Report §4 P1, rot nach 5 min 12 s); „Sperre weg **und**
  Vorbedingung weg“ bleibt *hergeleitet*, nicht erprobt. Im roten Versuch in CI
  lief die Phase nicht (Schritte des Legs `skipped`): das ist der Ausfall der Phase
  davor, kein Befund zum Aufbau.
- **Die Laufzeit von `make test-integration` und `make test-replication` wächst**
  (`BEO-PGC/test-integration-retention-timing-flake`, 3×, verkörpert). *Erwartet,
  zu belegen durch:* die gedruckte Laufzeit je Lauf; die Zahl trägt ihren Lauf
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). **Ausgang: weiter offen**
  — im Register `BEO-PGC/test-integration-retention-timing-flake` (4×; `state.md`
  trägt die Frage an den Architect). Gemessen, Läufe je einmal: lokal
  `make test-integration` 373 s (Verifier, Exit 0, aus dem Report **übernommen**)
  und der Keepalive-Test 40,3 s je Tier-Lauf; in CI, aus den Schritt-Zeitstempeln
  der Jobs **abgeleitet**, Lauf 36279515849 (Stand `7367c483`, ohne diesen Slice)
  gegen Lauf 36287009221: Schritt „Replication-Tier“ 91 s → 167 s (Leg 18, zweiter
  Versuch) und 104 s → 159 s (Leg 17); Schritt „Compose-Integrationstest“ 535 s →
  606 s (Leg 18, zweiter Versuch) und 568 s → 578 s (Leg 17). Das Wachstum ist
  eingetreten; als „eingetreten“ führt der Plan es nicht, weil es weder einen
  Carveout noch einen Folge-Slice trägt (kein Gate wird rot, kein Zeitlimit steht im
  Runner oder im Workflow), und seine mögliche Folge — eine intermittierende Phase auf
  dem gehosteten Runner — ist die offene Frage des Registers (die Phase
  „Leerlauf-Bestätigung“, 1 von 40 Läufen).
- **Die Ergänzung von `ADR-0121` behauptet mehr als der Beleg trägt**
  (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, 5×, verkörpert). *Erwartet,
  zu belegen durch:* der Reviewer liest die ADR im Diff gegen
  [`AGENTS.md`](../../../../AGENTS.md) §3.12 „Verfasser einer ADR“. **Ausgang:
  entfallen** für die Aussagen, die der Review gebunden hat.
  [`ADR-0129`](../../adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md)
  (seit `4e654153`) wurde vom Reviewer im frischen Kontext gegen den Diff
  `e4b77a05..4e654153` gelesen
  ([Review-Report](../../../reviews/review-adr-0129-capture-quellseite-keepalive.md),
  aus dem Report **übernommen**): Menge (zwei Pins mit Digest-Anker), gedruckte
  Zeilen, Mutationszuordnung (Stellen, Instanz, Farbe je Mutation; die Verallgemeinerung
  „auch an 17 rot“ steht als *hergeleitet*), Grenzen und Träger-Liste halten der
  Nachmessung stand; keine Fixrunde. Am Rand steht Reichweite, die die Messung nicht
  ganz trägt: F-1 (MEDIUM, CI-Anker ohne Versuchsnummer), F-2 und F-3 (je LOW, „ungebunden“
  und „nicht gemessen“ ohne Menge). Sie berühren keine Festlegung; als Ausgang
  „eingetreten“ führt der Plan sie nicht, weil kein Träger und keine Entscheidung an ihnen
  hängt und die ADR (`Accepted`) nicht berichtigt wird — die Berichtigung trägt eine
  spätere Berührung der ADR (Architect), die Lesehilfe steht in §7. Die Mutationen R1 bis
  R3 des Reviewers sind **übernommen**, nicht nachgefahren.

## 7. Closure-Notiz

Stand dieser Notiz: nach Review, Verifikation, Architect-Zug (`ADR-0129`), den
Post-Push-Läufen zu `e4b77a05` und dem ADR-Review zu `ADR-0129` (Commit `027533aa`);
alle Risiken aus §6 tragen einen Ausgang.

- **Was hat funktioniert:** Die Rollen-Kette trug. Der Review (0 HIGH · 1 MEDIUM ·
  3 LOW · 3 INFO, aus dem Report **übernommen**) fuhr sieben Mutationen am
  Keepalive-Test, drei davon an der Eingabeseite der Bestätigung; der Verifier
  fuhr elf verschiedene in 13 Läufen (neun rot, zwei grün mit Bedeutung: die
  Startposition des Neustarts und die Wartezeit auf die Inaktivität des Slots
  bindet der Test nicht), darunter zwei Mutationen, die zuvor nicht erprobt
  standen: „`stopStream` entfernt“, im Plan als „übernommen“ geführt (rot, der
  Container lief 90 s weiter), und den `--- PASS`-Wächter des Runners (rot), den
  der Review nur hergeleitet hatte. Der
  Aufbau „Persistierung halten“ trug in jedem Lauf, in dem die Phase lief (§6).
  Der Architect entschied den Befund der Klasse als Codefehler und beauftragte den
  Träger-Slice, ohne diesen Slice zu blockieren: die Grenze stand benannt im
  Plan und in `harness/README.md`. Die Quellseite ist an beiden Pins ein committeter
  Wächter, der seit dem Push in beiden Legs von `e2e.yml` läuft (DoD 1). Der ADR-Review
  (0 HIGH · 1 MEDIUM · 2 LOW · 2 INFO, aus dem Report **übernommen**) las `ADR-0129` im
  frischen Kontext gegen den Verfasser-Satz von `AGENTS.md` §3.12, verglich jede gedruckte
  Zeile mit dem Verifikations-Report und mutierte an PostgreSQL 17 an drei Stellen
  (R1 bis R3, **übernommen**, nicht nachgefahren): keine Fixrunde an der ADR.
- **Was ging anders als geplant:** (1) Der Keepalive-Test ist ein roher
  Protokoll-Client über **eine** Transaktion statt Stream-Adapter mit erster und
  zweiter Transaktion (§3); der Wortlaut von DoD 1 zog erst die Closure nach
  (Review F-1, Verifikation V-1). (2) Der zweite Beleg ist eine Bash-Phase statt
  einer `func TestE2E*` (§3, Reduktion). (3) Die Klasse des Ausgangs ist `storage`
  statt `replication`: „Grenze bleibt“, Träger
  `slice-wal-fehlerschwelle-ausgangsklasse`. (4) Der erste `e2e.yml`-Lauf nach dem
  Push war im Leg PostgreSQL 18 im ersten Versuch rot, in der **bestehenden**
  Phase „Leerlauf-Bestätigung“ (Lauf 36287009221, Versuch 1, Job 108529548391); die
  Schritte dahinter liefen dort nicht, der Wiederholungsversuch (Versuch 2, Job
  108531740887) war grün. Ursache offen:
  Register, Frage an den Architect. [`AGENTS.md`](../../../../AGENTS.md) §3.10
  greift nicht (kein Workflow im Diff), der Lauf war Beobachtung, keine
  Closure-Bedingung. (5) Das Suchlauf-Feld mit `diff` als zweitem Stand wurde mit
  `ADR-0129` rot (`make suchlauf-nachmessen`, gemessen am 2026-09-27 vor dem
  Nachzug: Zeile 2 soll 33, ist 58, Exit 2); die Stände stehen jetzt als
  Commit-Kennung. (6) Der Satz des Reviews „das Tier läuft nicht in
  `ci.yml`/`e2e.yml`“ war falsch (Verifikation V-4); der Workflow und der erste
  CI-Lauf haben ihn berichtigt.
- **Steering-Loop-Eintrag (Lerneintrag):** *(a) Neuer Sensor — zwei Belege, kein
  Gate.* Der Test `TestSourceKeepaliveInsideTransactionDeliversWholeTransaction`
  (Tier `make test-replication`, Phase `tier`, `--- PASS`-Wächter in
  `tools/harness/run-replication-tests.sh`) bindet die Quellseite der
  Leerlauf-Bestätigung an PostgreSQL 17.11 und 18.6, und die Runner-Phase
  „Fehlerschwelle beendet den Container“ in `make test-integration` bindet Ende,
  Ausgang, Abbruch-Zeile und Fehlerzustand der Schwellen-Kette am realen Prozess ·
  seit slice-capture-leerlauf-quellbelege; beide stehen in `harness/README.md`
  §Sensors, kein Gate (Laufzeit, Docker und Datenbank). *(b) Geschärfte Regel —
  keine, mit Grund.* Die Befunde dieses Slice sind Klassen mit Träger im Register
  (Nachzug lässt den Nachbarn stehen, Haken ohne Anker, Kommentar-Allaussage,
  Tatsachenbehauptung im Report, Aussage einer ADR breiter als ihre Messung, Beleg-Adresse
  ohne Versuchsnummer); ein neuer Wortlaut in `AGENTS.md` oder im Skill
  fügt keiner davon eine Linie hinzu. Der ADR-Review hat den Verfasser-Satz von
  `AGENTS.md` §3.12 als Leser gebraucht und an drei Rändern Reichweite gefunden, die die
  Messung nicht trug (Lesehilfe unten): die Klasse hat ihren Leser, der Satz bleibt.
  *(c) Benannte Lücken, keine Spec-Lücke.*
  (i) Nicht gebunden im Keepalive-Test: die Startposition des Neustarts **nach unten**
  (gefahren ist `START_REPLICATION` ab 0, grün; ab `confirmed_flush_lsn + 1` ist der Test an
  PostgreSQL 17 rot, Reviewer R2, **übernommen**), die Wartezeit auf die Inaktivität des
  Slots (ein grüner Lauf der Streichung), mehr als ein Keepalive je Transaktion,
  `proto_version 2` mit Streaming, eine zweite gleichzeitige Quelltransaktion **im Test**
  (`ADR-0129` Festlegung 2; die Lesehilfe unten schärft die Reichweite). (ii) Die Klasse `storage` der Fehlerschwellen-Kette ist ein
  Codefehler mit Träger. (iii) Es gibt keine Regel dafür, ob der erste CI-Lauf eines
  Tests, der in einer **bestehenden** Workflow-Phase mitläuft, Closure-Bedingung ist:
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 gilt für einen neuen oder strukturell
  geänderten Workflow. Die Praxis dieses Slice ist Beobachtung mit Anker (`ADR-0129`
  Festlegung 1 und Folgepflicht 5, Verifikation V-4); ein Lauf ist kein Muster, der
  Trigger für eine Regel ist ein zweiter Slice, dessen CI-Beleg an einer
  intermittierenden Vor-Phase hängt. (iv) Ein Beleg hinter einer intermittierenden
  Phase ist im roten Versuch „nicht gelaufen“, nicht „grün“: in Lauf 36287009221
  Versuch 1 standen die Coverage-Schritte und „Replication-Tier“ des Legs 18 auf
  `skipped`. Die Aussage über die Quellseite ist kein Spec-Satz: `SPEC-012` und
  `SPEC-013` sind gelesen, nicht geändert.
- **Beobachtungs-Register (`../observations/`):** je Anfall geschrieben. Neue
  `evidence/slice-capture-leerlauf-quellbelege.md` (Zähler gemessen mit `ls
  evidence | wc -l` am 2026-09-27): `BEO-PGC/test-integration-retention-timing-flake`
  (4×; die Resthälfte „`make test-integration`-Lauf rot bei unverändertem Stand,
  Wiederholung grün“ steht bei 2× ohne Ausgang, Frage an den Architect im
  `state.md`), `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (7×; Review
  F-4, F-5, Verifikation V-5, V-6), `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
  (14×; F-1, V-1, MEDIUM), `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker`
  (4×; F-2, V-2), `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (11×;
  V-4), `BEO-PGC/slice-pfad-als-link-in-berichten` (5×; F-7), und aus dem ADR-Review zu
  `ADR-0129` (dessen Findings hier als „ADR-Review F-<n>“ geführt, um sie von den Findings
  des Slice-Reviews zu trennen): `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (9×;
  ADR-Review F-2 und F-3, beide LOW, eine Datei; der Zähler stand vor dem Review bei 8×, der
  Deckel greift erst ab 10×) und `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (16×;
  ADR-Review F-1, MEDIUM, Form „Adresse“, daher Datei trotz Deckel bei 14×). `state.md`
  fortgeschrieben (zusätzlich zu den beiden neuen Dateien):
  `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` (Ausgang „verkörpert“, 1×),
  `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (Falsifikation der Phase,
  unverändert 5×). **Ohne Datei (Deckel, `../observations/README.md`):** F-3 und V-3
  (Suchlauf-Feld bindet `diff`, LOW) an `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
  (Deckel bei 23×), F-5 und V-6 (Godoc von `confirmIdle`, INFO; Übergabe an
  `slice-code-kommentare-bereinigung`) an `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
  (Deckel); alle vor dem Merge von Lesern gefunden. **Ohne Klasse:** F-6 (Bindung und
  Grenze des Tests, steht in `ADR-0129` Festlegung 2), V-7 (dieselbe Bindung), V-8
  (übernommen und nicht gefahren, im Report benannt), ADR-Review F-4 (INFO: Aussage (a)
  bindet die Zeile, an der keine Mutation der ADR ansetzt; Reviewer-Probe R1 trägt, keine
  Aktion erwartet) und ADR-Review F-5 (INFO: „etwa 11,1 Tausend Änderungen“ ist über die
  Läufe nicht stabil, der Vorbehalt „kein Vertrag“ steht in der ADR).
- **Lesehilfe zu `ADR-0129` (ADR-Review F-1 bis F-3; die ADR ist `Accepted` und bleibt
  unberührt, die Berichtigung trägt ihre nächste Berührung durch den Architect):**
  (F-1) Der CI-Anker der ADR nennt Lauf 36287009221 ohne Versuchsnummer. Am 2026-09-27
  vom Planner **gemessen** (`gh run view 36287009221 --attempt 1|2 --json jobs`, Job-Logs
  mit `gh api repos/pt9912/pg-change-feed/actions/jobs/<Job>/logs`): *Versuch 1* —
  PostgreSQL 18, Job 108529548391, `failure` im Schritt „Compose-Integrationstest
  (Black-Box-E2E)“, die Schritte „Replication-Tier“ und beide Coverage-Schritte `skipped`
  (Log der Phase „Leerlauf-Bestätigung“: „WAL-Rückstand 15238216 Bytes über Fehlerschwelle
  8388608 Bytes“); PostgreSQL 17, Job 108529548457, `success`, `PASS` 40,72 s, 11538 von
  400000 Änderungen; *Versuch 2* — PostgreSQL 18, Job 108531740887, `success`, `PASS`
  43,20 s, 11541 Änderungen; PostgreSQL 17, Job 108531741788, `success` (der zweite
  Versuch führt beide Legs auf). Der Satz der ADR „in CI weder belegt noch widerlegt“ gilt
  für ihren Stand gegen 02:07 UTC; belegt ist der Test in CI durch diese Läufe. Ohne
  `--attempt` liefert der Befehl heute Versuch 2. (F-2) „Ungebunden“ in Festlegung 2
  Punkt 4 heißt: gefahren ist ein Wert; die Startposition des Neustarts ist nach oben
  gebunden (`confirmed_flush_lsn + 1` ist an PostgreSQL 17 rot, Reviewer R2,
  **übernommen**), die Wartezeit stützt sich auf einen grünen Lauf. (F-3) „Eine zweite
  committende Transaktion ist nicht gemessen“ in Festlegung 2 Punkt 2 gilt **für den
  committeten Test**; `ADR-0121` §Gemessen (in Kraft) berichtet den Aufbau einmalig an
  PostgreSQL 18 (Stand-in hält die erste Transaktion, eine zweite committet 400.000
  Änderungen).
- **Folge-Slices:** keiner neu angelegt. Adressen: `slice-wal-fehlerschwelle-ausgangsklasse`
  (Klasse der Fehlerschwelle, Kommentare der Runner-Phase), `slice-start-vorlauf-grenze`,
  `slice-code-kommentare-bereinigung` (Übergabe der zwei Godoc-Kommentare zum durch
  `ADR-0121` ersetzten Begründungssatz, Tranchen T4 und T8) — alle in `open/`. Ein
  Slice zur Stabilisierung der Phase „Leerlauf-Bestätigung“ ist nicht angelegt: die
  Register-Regel verlangt ihn erst ab 3× ohne Ausgang (die Resthälfte steht bei 2×), und
  ob er nötig ist, entscheidet die Frage an den Architect
  (`BEO-PGC/test-integration-retention-timing-flake`, `state.md`). Wird er beauftragt,
  geht er `slice-wal-fehlerschwelle-ausgangsklasse` und `slice-start-vorlauf-grenze`
  voraus (bedingte Kante in [welle-transformationen](../welle-transformationen.md) §5 und
  im Start-Trigger von `slice-wal-fehlerschwelle-ausgangsklasse`).
- **Folgepflichten aus `ADR-0129` (mit Adresse):** (1) Register
  `beleg-nur-als-einmalige-reviewer-messung`: nachgezogen. (2) Träger der alten
  Grenze: die Records bleiben, `ADR-0121` bleibt `Accepted`, dieser Plan trägt die
  neue Fassung (§3 Suchlauf-Feld). (3) Code-Kommentare: an
  `slice-code-kommentare-bereinigung` übergeben. (4) Slice-Plan: DoD 1 und DoD 3,
  Stände des Suchlauf-Feldes nachgezogen. (5) CI-Beleg und Befund: DoD 1 und DoD 2
  tragen die Läufe, das Register die Zuordnung der roten Phase.
- **Risiken aus §6:** fünf Ausgänge gesetzt, je Risiko einer — vier entfallen mit Messung
  (der Test ist nicht zeitabhängig ausgefallen, 17.11 liefert dieselbe Position, der Aufbau
  der Phase trug in jedem Lauf, `ADR-0129` hält der Lese-Prüfung des Reviewers an ihren
  tragenden Aussagen stand, mit den Rändern der Lesehilfe oben), einer weiter offen im
  Register (Laufzeit, `BEO-PGC/test-integration-retention-timing-flake`).
- **Drei Paarungen:** dieser Slice hat keine Welle; die Roadmap führt
  [welle-transformationen](../welle-transformationen.md) unter *Offene Wellen*, das
  Ereignis kann eintreten: die Closure dieser Welle prüft die Paarungen mit. Die
  Slice-Closure trägt sie zusätzlich jetzt: *Anker:* der Test, die Runner-Phase und
  die ADR existieren als Dateien
  (`internal/adapters/driving/replication/receive/sourcekeepalive_test.go`,
  `tools/harness/run-integration-tests.sh` mit `abdeckung_declare` der Phase,
  `docs/plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md` samt
  Index-Zeile), `harness/README.md` §Sensors nennt beide Belege
  (`git grep -o 'TestSourceKeepaliveInsideTransactionDeliversWholeTransaction' --
  harness/README.md` trifft 1). *Folge-Slice:* die drei Adressen oben existieren als
  Dateien in `open/`; kein Versprechen ohne Adresse. *Register:* jede genannte
  Kennung `BEO-PGC/<slug>` existiert als Verzeichnis mit nicht leerem `evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Replication-Adapter, Composition Root und Test-Runner
sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Zähler gemessen am 2026-09-25 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` (1×, Ausgang: dieser Slice),
`BEO-PGC/adr-aussage-breiter-als-ihre-messung` (verkörpert, 5×, die
ADR-Ergänzung), `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (3×,
verkörpert, Kommentare zur Schwellen-Kette),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×, Risiko §6),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, Plan-Zeile Runner).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
