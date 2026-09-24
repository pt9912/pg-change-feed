# Slice backfill-e2e: E2E — Happy Path, Boundary und Negative am laufenden Feed-Container; gemessene Startposition eines frisch registrierten Consumers

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Happy Path, Boundary, Negative — alle drei am
laufenden System), [`LH-FA-CAP-004`](../../../../spec/lastenheft.md), [`LH-FA-REA-001`](../../../../spec/lastenheft.md), [`LH-FA-CON-005`](../../../../spec/lastenheft.md)
(Anfangsposition eines neu registrierten Consumers), [`LH-QA-POR-003`](../../../../spec/lastenheft.md) (E2E-Ebene),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 1 (Startposition messen und dokumentieren) und
Folgepflicht 2 (`S5`), [`ADR-0030`](../../adr/0030-testpyramide.md) (Testpyramide, E2E-Tier), [`ADR-0058`](../../adr/0058-testansatz-fuenf-luecken.md)
(Testansatz — additive Belege am realen Container), [`ADR-0012`](../../adr/0012-at-least-once.md)
(at-least-once, Consumer arbeiten idempotent), [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2 (Aufnahme einer
`queued`-Zeile beim Prozessstart), [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md)
(Umschreiben der Tabelle im Snapshot-Fenster; Fixrunde dieses Slice, Verdikt
`architect-verdict-backfill-tabellen-rewrite-im-fenster` unter `docs/reviews/`).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md), [`SPEC-022`](../../../../spec/pflichtenheft.md), [`SPEC-029`](../../../../spec/pflichtenheft.md) — gelesen
als Vertrag der Belege, nicht geändert; [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) Absatz „Mechanismus“ trägt
einen Satz zur Lesesperre und zum Ende bei umgeschriebener Tabelle;
`spec/architecture.md` §4 (Sequenz der Run-Ausführung) trägt Lesesperre und
Umschreib-Prüfung des Snapshot-Lesers.

**Verantwortlich:** Implementer-Agent, 2026-09-24.

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die drei Akzeptanzkriterien von [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) sind in
`make test-integration` **am laufenden Feed-Container** belegt, im Muster der
bestehenden Rundläufe (ausschließlich externe Wege: `docker exec`, SQL gegen
`cdc.changes`/`cdc.backfill_status`, HTTP — kein Import aus `internal/**`):

- **Happy Path.** Eine Tabelle mit Bestand wird aktiviert, `SELECT
  cdc.backfill_table(…)` beantragt den Backfill, ein Poll auf
  `cdc.backfill_status` sieht `completed`; jede zum Startzeitpunkt vorhandene
  Zeile ist über `cdc.changes` **und** `GET /changes` als Change lesbar
  (Tabelle, Operation `INSERT`, Row Image, `origin = 'backfill'`), unterscheidbar
  von einer danach WAL-erfassten Änderung derselben Zeile (`origin = 'wal'`); die
  `change_id`s sind unabhängig gegen `cdc.changes` gehalten.
- **Boundary.** Nebenläufige `INSERT`/`UPDATE`/`DELETE` an derselben Tabelle
  während des Runs; danach ergibt das Log ab Log-Anfang, angewandt in
  Lese-Ordnung (`INSERT`/`UPDATE` als Upsert des Row Images, `DELETE` als
  Löschen), je Schlüssel den Quellstand (**Replay-Invariante**); eine leere
  Tabelle endet `completed` mit 0 Zeilen und schreibt keine Transaktion; ein
  zweiter Antrag bei aktivem Run endet `failed` mit Grund.
- **Negative.** `docker kill` des Feed-Containers **mitten im Run**, Neustart:
  der Run steht `interrupted`, **keine** Change des Runs ist sichtbar, ein
  erneuter Antrag ergibt einen neuen Run, der `completed` erreicht — der Bestand
  ist danach **einmal** und vollständig lesbar (keine Dopplung durch den
  Abbruch). Eine zum Abbruchzeitpunkt `queued` wartende Zeile (ein zweiter Antrag
  gegen eine zweite Tabelle, angenommen hinter dem hängenden Run) überlebt den
  Neustart und wird beim Prozessstart aufgenommen und ausgeführt, ohne dass ein
  neuer Antrag nötig ist ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2).
- **Startposition.** Von welcher Position ein frisch registrierter Consumer
  startet, ist in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 1 „erwartet, nicht geprüft": der Slice
  misst sie (Registrierung über den CLI-Weg, Position über den HTTP-Weg und
  `cdc.consumer_status`) und trägt das Ergebnis samt Lauf ins Handbuch ein — als
  Bezugsmuster, wie ein Consumer den Bestand erhält.
- **Umschreiben im Snapshot-Fenster** ([`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md)).
  Ein Umschreiben der Tabelle zwischen Snapshot-Export und Import endet den Run
  `failed` (Klasse `transient`) ohne Change; ein neuer Antrag danach endet
  `completed` mit allen Zeilen. Der Import steht dafür unter einer Lesesperre
  mit Filenode-Vergleich (`postgressnapshot`, `snapshotlogic`); die Belege
  liegen in `make test`, `make test-replication` und der Phase DDL-Fenster von
  `make test-integration`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Fail-closed-Prüfung am Container** (`exclude_column` während des Runs) — die
  Negativ-Tests mit Fakes in `run-usecase` tragen sie; ein Zeitfenster-Test am
  Container wäre nicht deterministisch.
- **Retention-Beleg für Backfill-Changes** — die Regel gilt unverändert
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 7); ein Beleg über 24 Stunden `MinAge` wäre keiner
  Laufzeit tragbar.
- **Bench und Warn-Richtgröße** — `bench-richtgroesse`.
- **Die Live-Wege** — Backfill-Changes gehen nicht in gRPC, SSE oder
  NATS-Vollinhalt (Welle §6); ein Nicht-Beleg dort ist kein Kriterium.
- **Die SDK-Läufe** (`make test-sdk-*-integration`) — `sdk-origin`.

## 2. Definition of Done

- [x] Happy Path und Boundary am laufenden Feed-Container: Bestand lesbar und als
      `backfill` erkennbar über beide Lesewege, unterscheidbar von einer
      WAL-Änderung derselben Zeile, Replay-Invariante nach nebenläufigen
      Schreibern, leere Tabelle, zweiter Antrag. *Zu belegen durch:* ein realer
      `make test-integration`-Lauf mit Exit 0; der Bericht nennt den Lauf mit
      den gedruckten `change_id`s. **Ort der Replay-Invariante:** die
      Fitness-Function-Zeile in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt den Tier `make test-replication`,
      Folgepflicht 2 den E2E-Beleg (`S5`); dieser Slice belegt sie im E2E, weil erst
      dort Snapshot, Schreiber und WAL-Pfad komponiert laufen. Die Ortswahl ist
      eine benannte Entscheidung des Planners (§3, letzte Zeile); kein
      Architect-Verdikt trägt sie.
- [x] Negative: `docker kill` im laufenden Run → nach dem Neustart `interrupted`,
      keine sichtbare Change des Runs, erneuter Antrag erreicht `completed`, Bestand
      einmal und vollständig; eine zum Abbruchzeitpunkt `queued` wartende Zeile
      überlebt den Neustart und wird ausgeführt (`completed`, Bestand ihrer Tabelle
      lesbar). *Zu belegen durch:* `make test-integration`; **jedes** dieser
      Kriterien trägt je eine Mutation im Bericht (die Prüfung gegen die
      Eingabe gelenkt, der Lauf färbt rot —
      `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert). Im Repo
      belegt sind die Mutationen der Zusagen zu [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md) (Verifikations-Report §4);
      die der Negative-Phase trägt kein committetes Artefakt (§7).
- [x] Startposition gemessen und im Handbuch mit ihrem Lauf genannt; der
      Runner deklariert die Phase(n) über `abdeckung_declare` mit der Kennung
      [`LH-FA-CAP-009`](../../../../spec/lastenheft.md), und `docs/user/e2e-abdeckung.md` trägt nach dem Lauf eine
      Zeile dafür (Erzeugnis des Runners, kein Lauf-Beleg, mitcommittet);
      `make doc-trace` führt [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) nicht mehr unter den Waisen. Jede
      neue `func TestE2E*` steht in einem `-run`-Muster des Runners oder ist eine
      deklarierte Runner-Phase — *zu belegen durch:* der Lauf zeigt jede in der
      `-v`-Ausgabe (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Umschreiben im Snapshot-Fenster ([`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md)):
      Lesesperre und Filenode-Vergleich im Import, Abbruch als Klasse
      `transient`, ein neuer Antrag erreicht `completed` mit allen Zeilen.
      *Zu belegen durch:* `make test` (Anweisungs- und Entscheidungslogik),
      `make test-replication` (Umschreiben, Kontrolle, partitionierte Tabelle,
      Wartezeit an der Sperre), `make test-integration` (Phase DDL-Fenster,
      zweiter Lauf); je Zusage eine Eingabeseiten-Mutation im Bericht.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
      Kein offenes HIGH/MEDIUM: F-1 (MEDIUM) ist in der Fixrunde behoben, der Report trägt 0 HIGH.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: Handbuch §4 (Abschnitt „Bestand als Backfill überführen“) trägt die gemessene Startposition mit Lauf-Ursprung; die Änderungshistorie eine Zeile.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | unverändert | die Go-Hälfte liegt in der neuen Datei `test/integration/backfill_e2e_test.go` (eigene Helfer; der Abdeckungs-Erzeuger liest alle `.go`-Dateien des Pakets). Go-seitig laufen nur die Replay-Anwendung und der Bestandsvergleich, die übrigen Belege sind Runner-Phasen. |
| `test/integration/backfill_e2e_test.go` | neu | `TestE2EBackfillReplayInvariant` (Boundary: drei nebenläufige `INSERT`/`UPDATE`/`DELETE`-Schreiber vor und nach der Snapshot-Position, Replay des Logs gegen den Quellstand, Überlappung der beiden Seiten als Abbruchbedingung) samt Hilfsfunktionen; ausschließlich externe Wege (SQL, keine Importe aus `internal/**`). Der Godoc der Testfunktion nennt die Zusage der Überlappung ohne Nebenklausel zur verworfenen Alternative (Review F-5). |
| `tools/harness/run-integration-tests.sh` | update | sieben Rundläufe (sechs Runner-Phasen mit `abdeckung_declare` auf `LH-FA-CAP-009`, ein eigener Go-Aufruf): Happy Path, Schema-Version, Startposition, Boundary (leere Tabelle, zweiter Antrag), Replay-Invariante, DDL-Fenster, Negative (`docker kill`, `queued`-Aufnahme); drei Haltepunkte (offene Schreibtransaktion, unbestätigter Schlüssel im zweiten Block, `docker pause`); der Tabellensperren-Ansatz ist nicht realisiert (siehe Ansatz-Ergebnis). Die Phase DDL-Fenster fährt zwei Läufe über die Funktion `bf_ddl_window`: `DROP COLUMN` endet `failed`/`storage`, `ALTER COLUMN … TYPE` endet `failed`/`transient` ([`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md)); nach jedem Lauf übernimmt ein neuer Antrag den Bestand (beim Umschreiben `rows_copied` 3). Fixrunde zum Review: die Assertion zum Erfassungspfad in `bf_ddl_window` liest die Zeile von `cdc.heartbeat` als `count(*)` mit `error_class IS NULL` gleich 1 und färbt bei fehlender Zeile rot (F-4); der Kommentar zum Haltepunkt der Konflikt-Sitzung nennt die Kopplung im Indikativ (F-5). |
| `internal/adapters/driven/postgressnapshot/snapshot.go` | update | Fixrunde ([`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md)): `importSnapshot` ruft nach dem Import und dem Ende der Replication-Verbindung `lockAndVerify` (Lesesperre, Filenode-Vergleich) vor `readColumns`; Paketkopf und Kommentare tragen die Abfolge. |
| `internal/adapters/driven/postgressnapshot/snapshotlogic/logic.go` | update | `LockStatement`, `RewriteQuery`, `CheckRewrite` (Entscheidung: keine Zeile und `t` enden `transient`, `f` setzt fort, alles andere `storage`) und `ClassifyLock` (`42P01`/`3F000` bleiben `configuration`) — netzlos, im Unit-Gegenstand des Coverage-Gates ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)). |
| `internal/adapters/driven/postgressnapshot/snapshotlogic/logic_test.go` | update | Unit-Tests der Anweisungen (Quoting), der Abfrage (`coalesce`), der Entscheidung und der Sperr-Klassifikation. |
| `internal/adapters/driven/postgressnapshot/snapshot_test.go` | update | Store-Tier: `TestRewriteInWindowIsTransient` (fünf Formen), `TestNoRewriteInWindowReadsTheSnapshot` (Kontrolle, kompatible Erweiterung, `DROP COLUMN` → `storage`), `TestPartitionedTableIsNoFalseAlarm`, `TestImportWaitsForExclusiveLockAndThenAborts`; `TestCatalogQueryFailuresKeepTheirClass`, `TestPermissionClassWithoutSelect` und `TestConfigurationClass` tragen die Phase der Sperr- und Prüfungs-Anweisung. Fixrunde zum Review (F-6): `TestImportLockWaitEndsWithTheContext` bindet den Kontext-Abbruch an der wartenden Sperranweisung (Klasse `transient`, keine Sitzung danach). |
| `spec/pflichtenheft.md` | update | `LH-FA-CAP-009.a` Absatz „Mechanismus“: ein Satz (Wortlaut aus [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Folgepflicht 2) und eine Zeile der Änderungshistorie. |
| Kommentar zum zweiten `make schema-rollout`-Lauf in `tools/harness/run-integration-tests.sh` (gemeldet von `slice-backfill-change-origin`, Review F-7) | update | der Kommentar („… real blockierten zweiten `make schema-rollout`-Lauf (… Exit 8 auf vier Fremdobjekten …)") beschrieb den Stand der Wache nicht: er nennt jetzt den Ist-Zustand (der Tausch belegt das Upgrade des Prozesses; den zweiten Rollout belegt `tools/harness/run-schema-rollout-guard-test.sh`, idempotent über die Fremdobjekte des Guards). **Plan-Drift:** der Plan nennt sechs Fremdobjekte, `knownForeignObjects` in `tools/schema/rolloutguard/guard.go` führt sieben (fünf Funktionen, zwei Views); der Kommentar nennt sieben. Die Zeilen-Anker von `docs/user/e2e-abdeckung.md` verschieben sich hier ohnehin (Erzeugnis, mitcommittet) — der Slice, der den Kommentar ändert, ist der Slice, der die Anker regeneriert. |
| `tools/harness/httpclient/main.go` | update | geliefert: `origin` in der READ-Zeile (geprüft gegen `wal`/`backfill`, der Beleg der Feldform am Wire) und die Modi `changes` (`GET /changes` ohne Registrierung) und `position` (`GET /consumers/position`). |
| `compose.yaml` | prüfen | unverändert: `max_replication_slots=10` und `max_wal_senders=10` tragen die Reserve für den temporären Slot; `wal_sender_timeout=2000` begrenzt die Pause des `docker pause`-Haltepunkts (der Runner hält sie unter der Hälfte); die drei DSNs sind Superuser-Verbindungen, die Rollen-Zusagen tragen die Store-Tier-Tests (Architect-Verdikt, „Akzeptiertes Negativ"). |
| `internal/application/usecase/backfill/service.go` | update | Grenze-Kommentar an `currentVersion` (Auflage aus dem Architect-Verdikt, [`ADR-0116`](../../adr/0116-backfill-schema-version-referenz-reichweite.md) Folgepflicht 1; kein Verhaltens-Diff). Fixrunde zum Review (F-8): `failureText` entfernt die Klassen-Angabe, die ein Fehlerwert der Ports selbst trägt („Fehlerklasse transient: …“), vor dem Setzen des Fehlertexts; `error_message` nennt die Klasse einmal. |
| `internal/application/usecase/backfill/service_test.go` | update | `TestExecuteMarksRunningBeforeOpeningSnapshot` bindet die Reihenfolge `MarkRunning` vor `OpenSnapshot` (Übergabe aus `slice-backfill-run-usecase`). Fixrunde (F-8): `TestExecuteFailureTextCarriesClassOnce` bindet den Fehlertext an die einmalige Klasse; die Erwartung der Fail-closed-Fälle vergleicht gegen den Text ohne Klassen-Angabe. |
| `docs/user/e2e-abdeckung.md` | regeneriert | Erzeugnis des Runners, mitcommittet. |
| `docs/user/benutzerhandbuch.md` | update | geliefert: Punkt „Startposition eines neuen Consumers“ (Ursprung: der Lauf), `rows_copied` im Zustand `interrupted`; Absatz „Sperre der Tabelle“ und Punkt „Umschreiben der Tabelle im Fenster“; Version 1.52 samt Änderungshistorie. Fixrunde zum Review: der Absatz „Sperre der Tabelle“ nennt die gestauten Zugriffe vollständig (Schreiber, Leser, Publication-Abfrage der Administration; gemessen im Review, PostgreSQL 18) und die Gegenrichtung (der Run wartet ohne eigene Zeitgrenze auf eine offene `ACCESS EXCLUSIVE`-Transaktion; F-1, F-6); `RENAME COLUMN` im Fenster steht als im Review gemessen, ohne Beleg im E2E-Runner (F-7). |
| `spec/architecture.md` | update | Fixrunde zum Review (F-2, Entscheidung des Auftraggebers: die Sicht trägt sie): das Sequenzdiagramm der Run-Ausführung führt „Lesesperre auf die Tabelle, Umschreib-Prüfung“; der Absatz darunter nennt Sperrdauer, wartende DDL, Ende des Runs bei umgeschriebener Tabelle (`failed`, Klasse `transient`, ohne Change) und die Aufgabe des Snapshot-Lesers — ohne ADR-, Slice- oder Wellen-Bezug (`AGENTS.md` §3.4). |
| `docs/plan/planning/observations/BEO-PGC/lesesperre-ohne-zeitgrenze/` | neu | Register-Eintrag zur Wartegrenze der Lesesperre (Review F-6, Ausgang *weiter offen*): `observation.md`, `state.md`, `evidence/slice-backfill-e2e.md`. |
| `harness/README.md` §Sensors | update | die Zeile `make test-integration` trägt die sieben Backfill-Rundläufe und das Umschreiben in der Phase DDL-Fenster; die Zeile `make test-replication` nennt Lesesperre und Filenode-Vergleich. |
| `docs/plan/adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md` und `docs/plan/adr/README.md` | neu / update | Closure: [`ADR-0119`](../../adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md) (`Supersedes` [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md), teilweise) setzt zwei Aussagen auf die gemessene Reichweite (Wirkung der Lesesperre auf wartende DDL; `RENAME COLUMN` ohne E2E-Beleg — Review F-1/F-7, Verifikation V-5); Index-Zeile. |
| `harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md` | update | Closure: Nenner und gedeckte Zahl mit Lauf-Ursprung — Coverage-Gate 2541 Statements (54 in `snapshotlogic`), gedeckt 2112, gedruckt 83.10 %; DB-Adapter-Coverage 1035 Statements (`postgressnapshot` 130), gedeckt 850, gedruckt 82.13 % (Messung im Bericht der Closure, §7). |
| Ortswahl der Replay-Invariante (Entscheidung des Planners, Closure; Verifikation V-1) | Entscheidung | Ort: `make test-integration` (`TestE2EBackfillReplayInvariant`). Anker der Wahl: [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Folgepflicht 2 führt die Replay-Invariante im Umfang von `S5` (E2E); die Fitness-Function-Zeile derselben ADR nennt zusätzlich den Tier `make test-replication` (`git grep -n 'Replay-Invariante' docs/plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md`: die Zeilen zu Herleitung, Folgepflicht und Fitness Function). Begründung: die Invariante verlangt Snapshot, nebenläufige Schreiber und WAL-Pfad zugleich; das leistet erst der komponierte Lauf am Feed-Container. Ein Tier-Beleg ist nicht geliefert (`git grep -n -i replay -- '*.go'`: nur `backfill_e2e_test.go` trägt die Invariante); weder der Review noch ein Architect-Verdikt hat die Ortswahl bestätigt oder den Tier-Beleg verlangt. Die Fitness-Function-Zeile nennt damit einen Beleg, den kein Test trägt — *benannte Lücke*, Adresse: Architect im Lese-Schritt der Closure von [welle-backfill-bestand](../welle-backfill-bestand.md) (Zeile schärfen oder Tier-Beleg verlangen). |
| Register-Einträge der Closure (`observations/BEO-PGC/…`) | neu / update | `evidence/slice-backfill-e2e.md` in fünf bestehenden Einträgen (`negativtest-ohne-bindung-an-seine-eingabe`, `zahl-in-traeger-driftet-gegen-die-messung`, `arbeit-ueberholt-stehenden-traeger`, `adr-aussage-breiter-als-ihre-messung`, `vorher-nachher-sprache-in-test-harness-kommentar`) und zwei neuen Einträgen (`plan-zusage-erfuellung-ohne-committeten-anker`, `run-fehlertext-traegt-klasse-doppelt`); `lesesperre-ohne-zeitgrenze` trägt sie aus der Fixrunde. Zähler und Ausgänge in §7. |

**Ansatz-Vorschlag, zu belegen (nicht bindend):** Ein „Abbruch mitten im Run" ist
nur mit einem festen Haltepunkt deterministisch. Der Slot-Anlage blockiert bis
eine beim Aufruf laufende Schreibtransaktion endet (M3 in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md), gemessen
5,27 s bei rund 5 s Restlaufzeit): hält der Test eine offene Schreibtransaktion
auf der Quelle, steht der Run in `running`, ohne dass er weiterkommt — dann
`docker kill`, Transaktion beenden, Neustart. Das setzt voraus, dass der Run
vor der Slot-Anlage `running` trägt; der Implementer prüft es am
Use-Case-Vertrag. Für den `queued`-Beleg beantragt der Test, während der erste Run
hängt, einen zweiten Run gegen eine zweite Tabelle (derselben Tabelle würde der
zweite Antrag als `failed` enden): der Ein-Worker-Betrieb hält ihn `queued`; nach
dem Neustart beobachtet der Test, dass er ohne neuen Antrag `completed` erreicht.

**Ansatz-Ergebnis (Implementer; gemessen, Ursprung je Punkt):**

- **Haltepunkt „offene Schreibtransaktion“ trägt.** Der Run steht `running` mit `snapshot_position` NULL an der Slot-Anlage (`make test-integration`, Phase Boundary und `TestE2EBackfillReplayInvariant`); die Voraussetzung „`running` vor der Slot-Anlage“ gilt am Container und ist am Use Case gebunden (`TestExecuteMarksRunningBeforeOpeningSnapshot`).
- **Haltepunkt „Tabellensperre“ trägt nicht** (Wegwerf-Läufe gegen `PG_TEST_IMAGE`, PostgreSQL 18): `pg_publication_tables` — die Abfrage der Publication-Mitgliedschaft im Antrag und am Run-Start — wartet auf jede Sperre der Tabelle (7,06 s Wartezeit an einer 8-s-Sperre); eine `ACCESS EXCLUSIVE`-Sperre weist ihrer Transaktion eine Transaktionskennung zu, auf die die Slot-Anlage wartet (der Run endete `transient` nach 30 s). Ersatz: ein unbestätigter Schlüssel `0bf-<Run>-00000002` in `cdc.transaction` hält den Schreiber im zweiten Block (Blockgröße 1000, Tabelle mit 2500 Zeilen), und `docker pause` des Feed-Containers, während die Haltetransaktion endet, öffnet ein Fenster nach dem Snapshot-Export, in dem die Quelle sperr- und DDL-frei zugänglich ist.
- **Abbruch mitten im Lauf.** `docker kill` im zweiten Block (`rows_copied` 1000, Status `running`): danach keine Sitzung (`client backend`, `walsender`) und kein Slot `cdc_bf_*`, keine sichtbare Change des Runs. Der Kill trifft eine offene Lese- **und** Schreibtransaktion; die Prüfung „keine sichtbare Change“ ist an die Ein-Transaktions-Form gebunden (ein Commit je Block färbt sie rot).
- **Fenster zwischen Export und Sperre.** `DROP COLUMN` im Fenster: `DECLARE` scheitert `42703`, der Run endet `failed` mit Klasse `storage`, ohne Change (Phase DDL-Fenster). Ein Umschreiben der Tabelle (`ALTER COLUMN … TYPE`) endet den Run `failed` mit Klasse `transient`, ohne Change (Phase DDL-Fenster, zweiter Lauf): die Lese-Transaktion sperrt die Tabelle und vergleicht den `relfilenode` im Snapshot mit dem aktuellen Katalog ([`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md)); ohne beide Schritte endet der Run `completed` mit `rows_copied` 0 bei drei Zeilen der Tabelle (gemessen, `ADR-0118` §Gemessen) — siehe §6.
- **Ergebnis gegen Zeile.** Die Belege lesen den Zustand aus `cdc.backfill_run` und `cdc.backfill_status`, also aus der Zeile.

**Übergaben aus `slice-backfill-snapshot-reader`** (gemeldet, aus dem Reader-Slice
gemessen oder benannt; die Komposition trägt erst dieser Slice):

- **Zeitlimit des Haltepunkts.** Die Slot-Anlage trägt ein Zeitlimit
  (`DefaultSlotTimeout`, 30 s, Startwert ohne Messung); sein Ablauf endet als
  `transient` (M3 im Reader-Test: 1-s-Limit gegen eine offene Schreibtransaktion). Hält
  der Test die Schreibtransaktion länger als das Limit, endet der Run `failed` statt zu
  hängen — der Haltepunkt des Ansatz-Vorschlags hält kürzer als das Limit (der Test
  läuft gegen den Container, nicht gegen den Adapter im Prozess).
- **Abbruch mitten im Lauf.** Der Reader belegt, dass der temporäre Slot mit der
  Replication-Verbindung endet (M4) und dass der Cursor den Snapshot-Stand liest; der
  Abbruch des komponierten Runs (`docker kill` mit offener Lese-Transaktion) ist nicht
  gefahren. *Erwartet, zu belegen durch:* nach dem Kill steht kein Slot `cdc_bf_*` in
  `pg_replication_slots` und keine Sitzung des Runs in `pg_stat_activity`.
- **Fenster zwischen Export und Sperre.** Ein gleichzeitiges `ALTER TABLE` in
  diesem Fenster lässt den `DECLARE` mit `storage` (`42703`) scheitern (`DROP COLUMN`)
  oder — bei einem Umschreiben der Tabelle — Katalog-Stand und Snapshot
  auseinanderlaufen; das Umschreiben endet den Run `failed`/`transient`
  ([`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md)). Der
  Haltepunkt (`docker pause` nach dem Export) misst beide Formen in der Phase
  DDL-Fenster; die Dauer des Fensters ist nicht gemessen.

**Übergaben aus `slice-backfill-run-usecase`** (gemeldet, kein zusätzlicher Umfang):

- **`running` vor der Slot-Anlage.** Der Ansatz-Vorschlag setzt voraus, dass der Run vor
  der Slot-Anlage `running` trägt. Am Code gelesen: `Execute` ruft `MarkRunning` vor
  `copyBlocks`, und erst dort öffnet `OpenSnapshot` den Slot
  (`internal/application/usecase/backfill/service.go`); kein Test des Use Cases bindet
  diese Reihenfolge, sie ist nicht gemessen.
- **Ergebnis gegen Zeile.** Bei einem Commit mit unbekanntem Ausgang kann die Zeile
  `completed` tragen, während das Ergebnis von `Execute` `failed` meldet
  (`BEO-PGC/execute-ergebnis-widerspricht-persistierter-zeile`); die Belege dieses Slice
  lesen den Zustand aus `cdc.backfill_status`, also aus der Zeile.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Menge der E2E-belegten Kennungen und ihre Zeilen-Anker", „die Beschreibung von `make test-integration`"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Abdeckungs-Tabelle | `git show e7df5619:docs/user/e2e-abdeckung.md \| grep -c '^\| \['` (Parent); Lauf von `make test-integration`, danach `git diff docs/user/e2e-abdeckung.md` und `grep -c` (Diff-Stand) | Parent: 41 Zeilen, 0 mit `LH-FA-CAP-009` (13 Go-, 28 Bash-Zeilen; ein Lauf des Parent-Runners im Arbeitsbaum dieses Slice meldet „aus 14 Go-Zeilen und 28 Bash-Zeilen“, weil `backfill_test.go` im Paket liegt). Diff-Stand: 48 Zeilen (Lauf-Zeile: „aus 14 Go-Zeilen und 34 Bash-Zeilen“), 7 mit `LH-FA-CAP-009`; `git diff --stat`: 35 Einfügungen, 28 Löschungen — 7 neue Zeilen, 28 Zeilen mit verschobenem Ort-Anker | mitcommittet; Zeilen-Anker früherer Zeilen verschieben sich |
| Zeile `make test-integration` in `harness/README.md` (lange Aufzählung der Rundläufe) | Lesen; Zählwort „sieben“ nachgezählt: `grep -c '^abdeckung_declare "Backfill' tools/harness/run-integration-tests.sh` und `grep -c '^func TestE2EBackfill' test/integration/backfill_e2e_test.go` | Aufzählung endet vor dem Satz zum Erzeugnis `docs/user/e2e-abdeckung.md`; Zählung: 6 Runner-Deklarationen + 1 Testfunktion = 7 Rundläufe, drei Haltepunkte (offene Schreibtransaktion, unbestätigter Schlüssel, `docker pause`); der Kopf des Abschnitts im Runner nennt dieselben drei | ein Satz ergänzt (sieben Backfill-Rundläufe), Rest unverändert; der Runner-Kopf-Kommentar nennt `docker pause` statt der Tabellensperre |
| Aussage zu den Waisen in `harness/README.md` (Zeile `make doc-trace`) | `make -C <Worktree von e7df5619> doc-trace` (Parent), `make doc-trace` (Diff-Stand); `git grep -n -i waise <Stand>` | Parent: „80 Anforderung(en), 3 Waise(n)“ (`LH-FA-CAP-009`, `LH-FA-CFG-007`, `LH-FA-CFG-008`); Diff-Stand: „80 Anforderung(en), 2 Waise(n)“, `LH-FA-CAP-009` mit Nachweis `E2E`. Träger der Waisen-Aussage: die README-Zeile (3 Waisen mit Datum 2026-09-23), sonst nur Pläne mit Zukunftsaussagen (`welle-backfill-bestand`, `welle-transformationen`, zwei offene Slices) | README-Zeile nachgezogen (2 Waisen, Datum 2026-09-24, `LH-FA-CFG-007`/`LH-FA-CFG-008`); Pläne unverändert |
| Beschreibungen von `make test-integration` außerhalb der README | `git grep -n 'test-integration' <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':(exclude,glob)docs/plan/planning/*/slice-backfill-e2e.md' \| wc -l` (Zeilen) und dieselbe Abfrage mit `-l` (Dateien); Parent `e7df5619`, Diff-Stand = Arbeitsbaum der Fixrunde (ohne `<Stand>`); dazu Lesen der beschreibenden Treffer (`Makefile`, `harness/sensors/docs-check.md`, `AGENTS.md`, `examples/README.md`, `docs/plan/planning/in-progress/roadmap.md`) | Parent: 124 Zeilen in 57 Dateien; Diff-Stand: 130 Zeilen in 60 Dateien (zusätzlich `docs/user/benutzerhandbuch.md` mit 3, `docs/plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md` mit 2 und `test/integration/backfill_e2e_test.go` mit 1 Zeile); beide Stände ohne die Plan-Datei dieses Slice, die sich selbst zitiert. Beschreibend: `Makefile` (Hilfezeile nennt drei Beispiele ohne Vollständigkeitsanspruch), `harness/sensors/docs-check.md` (Erzeugung der Abdeckungstabelle, unverändert wahr); die übrigen sind Verweise auf das Ziel ohne Beschreibung des Inhalts | unverändert |
| Ausgabe und Modi des Wegwerf-Clients | `git grep -n httpclient <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` und Lesen der Treffer außerhalb von `tools/` | Parent: 39 Treffer-Zeilen in 14 Dateien, Diff-Stand: 48 Zeilen in 14 Dateien (die Mehrzeilen liegen in `tools/harness/`); beide ohne die Plan-Datei dieses Slice; beschreibende Träger der Ausgabeform: `examples/http-client/main.go` (Verweis „Wegwerf-Client als Belegträger“), `examples/README.md` (Abgrenzung), `docs/plan/planning/open/slice-transformationen-e2e-wirkung.md` (Plan-Zeile „geben Schlüssel und Werte der Row Images aus“) — keine nennt die READ-Zeile oder deren Feldliste | unverändert |
| Zählwörter zu den bewegten Eigenschaften | `git grep -n -E '(zwei\|drei\|vier\|fünf\|sechs\|sieben\|acht\|beiden\|beide) (Backfill-\|weitere \|zusätzliche )?(Rundl\|Belege\|Phasen\|Haltepunkt\|Fremdobjekt\|Waise)' <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/adr'` (ohne die Plan-Datei dieses Slice) | Parent: 11 Treffer; Diff-Stand: 11 Treffer. Entfallen: „vier Fremdobjekten“ im Runner-Kommentar (ersetzt); hinzugekommen: „sieben Backfill-Rundläufe“ in `harness/README.md`. Die übrigen zehn betreffen Kurs-Dokumente (`.claude/commands/close-welle.md`, `.harness/baseline`), die SDK-Werkzeuge (`harness/mk/sdk.mk`, „vier Phasen“) und die Coverage-/Replikations-Skripte (zwei Phasen) und sind von diesem Slice nicht berührt | unverändert |
| RTM-Träger | `make doc-trace` | siehe Zeile zu den Waisen: `LH-FA-CAP-009` trägt `E2E`, 2 Waisen | Deklarations-Anker tragen |
| Träger der Fixrunde [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md) — bewegte Eigenschaft: „das Verhalten des Snapshot-Imports bei einem Umschreiben der Tabelle (Sperre, Vergleich, Klasse `transient`)“ | `git grep -n -i -E 'ohne (Tabellen)?sperre\|Tabellensperre\|Sperre der Tabelle\|Lesesperre\|Fenster zwischen (Snapshot-)?Export\|Export und (Snapshot-)?(Import\|Cursor\|Sperre)\|umgeschrieben\|Umschreib' <Stand> -- . ':!docs/reviews' ':!docs/plan/adr' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/in-progress/slice-backfill-e2e.md' ':!*_test.go' ':!.harness' ':!AGENTS.md' ':!.d-check.yml' \| wc -l` und `git grep -n -i -E 'Snapshot-Fenster\|Rewrite\|Fenster zwischen\|Tabellen-Rewrite\|DDL-Fenster' <Stand> -- docs/plan/planning ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` (ohne die Plan-Datei dieses Slice); Parent `ceea0af3` (`git grep … ceea0af3`) und Diff-Stand (Arbeitsbaum) | Parent: 4 Zeilen (`welle-backfill-bestand`: „`request_kind`-Menge würde zweimal umgeschrieben“ — anderer Gegenstand; `docs/user/e2e-abdeckung.md` Zeile 62, Runner-Deklaration und Runner-Ausgabe der Phase DDL-Fenster — beschreibend, nennen nur die Spaltenentfernung); zweite Abfrage: 1 Treffer (`slice-transformationen-antragsweg-usecase`, „Fenster zwischen Regel-Setzbarkeit und Backfill-Bindung“ — anderer Gegenstand). Diff-Stand: 37 Zeilen (neu: Code, Spec, Handbuch, README, Runner), zweite Abfrage unverändert 1 Treffer. Fixrunde zum Review (Parent `31946550`: 37 Zeilen, Arbeitsbaum ohne Stand-Argument: 44 Zeilen): zusätzlich `spec/architecture.md` mit 6 und `docs/user/benutzerhandbuch.md` mit 1 Zeile; zweite Abfrage unverändert 1 Treffer | Träger nachgezogen: Runner-Deklaration und -Ausgabe (mitgeändert), `docs/user/e2e-abdeckung.md` (Erzeugnis, regeneriert), `spec/pflichtenheft.md` `LH-FA-CAP-009.a`, Handbuch (1.52), `spec/architecture.md` (Sequenz der Run-Ausführung), beide Zeilen der `harness/README.md`, §3 und §6 dieses Plans; die zwei Treffer anderer Gegenstände unverändert |
| Sequenzdarstellung des Snapshot-Lesers in der Architektur-Sicht (bewegte Eigenschaft: „Ablauf des Imports: Sperre und Umschreib-Prüfung“; Review F-2) | `git grep -n -i -E 'snapshot\|REPEATABLE\|umschreib\|sperre\|filenode' <Stand> -- spec/architecture.md \| wc -l`; Parent `31946550`, Diff-Stand = Arbeitsbaum der Fixrunde; Lesen von §4 | Parent: 9 Zeilen (Komponenten-Nennungen in Zeile 41 und 59, Diagramm und Absatz der Run-Ausführung); Diff-Stand: 16 Zeilen (zusätzlich Diagrammzeile und Absatz zu Sperre und Umschreib-Prüfung). Nicht gefunden: eine weitere Sicht-Datei — `spec/architecture.md` ist die einzige | Diagramm und Absatz nachgezogen; Komponenten-Nennungen unverändert wahr |
| Aussage zur Wirkung der Tabellensperre auf andere Zugriffe (bewegte Eigenschaft: „wer sich hinter einer wartenden DDL staut“; Review F-1) | `git grep -n -i -E 'Leser der Tabelle\|stauen\|hinter (ihr\|der DDL)\|wartende DDL\|Sperr-Warteschlange\|nicht gemessen' <Stand> -- . ':!docs/reviews' ':!docs/plan/adr' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/in-progress/slice-backfill-e2e.md' ':!*_test.go' ':!.harness'`, gefiltert auf Treffer mit `sperre`, `DDL` oder `Leser`; Parent `31946550`, Diff-Stand = Arbeitsbaum | Parent: 3 Treffer, alle im Handbuch (Zeilen 456/458 „Leser der Tabelle … nicht gemessen“, Zeile 518 „Dauer … nicht gemessen“ zum Fenster); Diff-Stand: 2 Treffer im Handbuch (Schreiber und Leser, Fenster-Dauer). Weitere Träger: §6 dieses Plans (Zeile zur Wirkung „auf andere Leser … nicht gemessen“) — nachgezogen; `docs/plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md` (`Accepted`, Konsequenzen „hergeleitet, nicht gemessen“) bleibt unberührt, die Korrektur trägt der Review-Report | Handbuch und §6 nachgezogen; ADR unberührt |
| Text des Run-Fehlers `error_message` (bewegte Eigenschaft: „die Klasse steht einmal vor der Ursache“; Review F-8) | `git grep -n -E 'error_message' <Stand> -- docs/user spec harness README.md`, Lesen der Treffer zum Backfill; Parent `31946550`, Diff-Stand = Arbeitsbaum | beide Stände 15 Treffer; beschreibend: Handbuch (Zustandstabelle des Backfills: „trägt die Fehlerklasse vor der Ursache“; Diagnose-Abschnitt), `spec/pflichtenheft.md` Zeile 706 („mit der Fehlerklasse aus `SPEC-008`“) — alle bleiben wahr, keine nennt den Text mit doppelter Klasse | unverändert |
| Symbolnamen des Imports (`importSnapshot`, `readColumns`, `snapshotlogic`) | `git grep -l -E 'importSnapshot\|readColumns\|snapshotlogic' <Stand> -- . ':!docs/reviews' ':!docs/plan/adr' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/in-progress/slice-backfill-e2e.md' ':!*.go'` | beide Stände: dieselben drei Dateien (`harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md`, `tools/harness/db-coverage.sh`); sie nennen die Paketzugehörigkeit und Statement-Zahlen mit Lauf-Stempel (der Stempel bindet die Zahl an ihren Lauf, sie behaupten keinen aktuellen Stand) | unverändert |
| Zahl der Statements des Gegenstands `postgressnapshot`/`snapshotlogic` | Lesen der beiden Sensor-Dokumente an den Stellen mit Lauf-Stempel | gestempelte Werte (`snapshotlogic` 42 Statements, DB-Adapter-Coverage 82,08 % bei Stand `02b3059d`) bleiben wahr für ihren Lauf; die Zahlen dieses Laufs stehen im Bericht, nicht in den Sensor-Dokumenten | unverändert; nicht gefunden: ein Träger, der die Zahl ungestempelt als Ist-Stand führt |
| CI-Träger der Läufe | Lesen von `.github/workflows/e2e.yml` (Trigger, Matrix, `timeout-minutes`); `git diff --stat e7df5619 -- .github` | Trigger Pull Request und Push, Matrix über die PostgreSQL-Versionen 17 und 18, `timeout-minutes: 60`; `.github` ohne Diff | unverändert; die Laufzeit-Frage steht in §6 |
| Träger der Closure — bewegte Eigenschaften: „Nenner und gedeckte Zahl der Coverage-Messungen“ (Produktionscode in `postgressnapshot`/`snapshotlogic`/`service.go` bewegt), „Reichweite der Sätze zu Sperr-Warteschlange und `RENAME COLUMN`“ | `git grep -n -E '\b(2527\|1027)\b\|82[.,]08\|82[.,]9[04]?%' <Stand> -- . ':!docs/reviews' ':!docs/plan/adr' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/in-progress/slice-backfill-e2e.md'` (Parent `cf7f2d02`, Diff-Stand = Arbeitsbaum der Closure); `git grep -n -i -E 'stauen\|Sperr-Warteschlange\|E2E-belegt' -- docs/plan/adr`; `git grep -n -i 'RENAME COLUMN' -- . ':!docs/reviews' ':!docs/plan/planning/observations'` (Arbeitsbaum) | Parent: 7 Treffer-Zeilen in zwei Dateien (`harness/sensors/coverage-gate.md` 4, `harness/sensors/db-adapter-coverage.md` 3); Diff-Stand: 0. ADR-Kette: `ADR-0118` trägt fünf Treffer-Zeilen (Kontext-Zeile zu `DROP COLUMN` mit „E2E-belegt“, „Nicht gemessen“, Option F, Festlegung 5, Konsequenzen); `RENAME COLUMN` steht in `ADR-0118` (Festlegung 5), im Handbuch (Fenster-Absatz, Historienzeile 1.52) und in der neuen ADR. Nicht gefunden: ein weiterer Träger, der die Zahlen 2527/1027 oder „nur Leser“ führt | Sensor-Dokumente nachgezogen (Nenner mit Lauf); `ADR-0118` `Accepted`, die drei betroffenen Stellen berichtigt `ADR-0119`; Option F und die Kontext-Zeile zu `DROP COLUMN` bleiben wahr |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `sql-administration` in `done/` liegt,
`make image` real gelaufen ist (`compose.yaml` trägt keinen `build:`-Block,
[`ADR-0044`](../../adr/0044-image-beleg-semantik.md)), kein anderer Slice in `in-progress/` liegt (WIP-Limit 1)
**und** ein Architect-Verdikt unter `docs/reviews/` zur Schema-Version der
Backfill-Changes vorliegt, **das der Übergangs-Commit `next` → `in-progress` nennt**
(`BEO-PGC/start-trigger-ohne-uebergabe-artefakt`, offen, 1×). Die Frage des Verdikts
(Register: `BEO-PGC/backfill-schema-version-hinter-snapshot-spalten`, aus dem Review von
`slice-backfill-run-usecase`, F-8): der Run liest die aktuelle Schema-Version der Tabelle
vor dem Öffnen des Snapshots; sie wechselt allein mit der nächsten Relation-Nachricht des
WAL-Pfads und kann hinter den Spalten des Snapshots liegen (kompatible
Spalten-Erweiterung ohne WAL-Änderung) oder auf eine Version ohne `TableSchema`
verweisen (statische Erstaktivierung). Das akzeptierte Negativ von
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt nur die
Erweiterung „während des Runs". Ob der Wortlaut ausreicht oder eine Schärfung nötig ist
(Folge-ADR, ggf. eine Änderung des Zeitpunkts, an dem der Run die Version liest), ist eine
Entscheidung, keine Auslegung dieses Slice; der Beleg der Bild-Form über
`cdc.changes` hängt an ihr.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Happy Path/Boundary
  und Negative nicht in einem Review tragen — der abtrennbare Teil ist das
  Negative (dann als zweiter Slice `backfill-e2e-abbruch` mit Start nach
  diesem).
- `in-progress` → `open` (blockiert): falls ein Abbruch mitten im Run ohne
  Eingriff in die Produktion nicht deterministisch herstellbar ist (dann
  Architect-Frage nach einem Test-Haltepunkt im Run); falls die Replay-Invariante
  real **nicht** gilt (dann ein Befund gegen [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 3, keine
  Testanpassung).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test-integration` real grün (Exit
ungefiltert gesichert) + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Replay-Invariante gilt real nicht** — [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) führt sie als
  „hergeleitet, im Slice als Eigenschaftstest zu belegen". Ein Rot wäre kein
  Testfehler, sondern ein Befund gegen die Entscheidung; die Rückführung §4
  benennt das. *Erwartet, zu belegen durch:* der Lauf. **Ausgang:** **entfallen** — die
  Invariante gilt real: `--- PASS: TestE2EBackfillReplayInvariant` in beiden Legs von
  `e2e.yml` Lauf `36065957210` (Kopf `43137ebf`), gedruckt „Replay-Invariante:
  Snapshot-Position 28721360, 54 Backfill-Changes, WAL-Changes davor 409 und dahinter
  164, 57 Zeilen im Quellstand“ (PostgreSQL 17) und „… 31875224, 54 … 409 … 169, 57 …“
  (PostgreSQL 18) — gemessen an diesem Lauf (`gh run view 36065957210 --log`).
- **Nichtdeterministischer Abbruch** — siehe Ansatz-Vorschlag; ohne Haltepunkt
  bliebe der Test flackernd. *Erwartet, zu belegen durch:* mehrere
  Wiederholungen im Bericht. **Ausgang:** **entfallen** — die Haltepunkte sind
  mechanisch fest (offene Schreibtransaktion, unbestätigter Schlüssel im zweiten Block,
  `docker pause` mit gemessener Grenze); die Phase Negative lief grün im E2E-Lauf des
  Reviews (Review-Report, 277 s) und in beiden Legs von `e2e.yml` Lauf `36065957210`
  (gedruckt „Backfill-Negative (docker kill, queued-Aufnahme) belegt“, gemessen an
  diesem Lauf). Aussagegrenze: ein Lauf je Leg und ein lokaler Lauf sind kein Beweis
  gegen Flackern (`AGENTS.md` §3.10); tritt ein Flackern auf, öffnet es das Risiko erneut.
- **Randfall-Belege nur beim Reviewer** (`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt`,
  offen, 1×): die leere Tabelle und der zweite Antrag sind committete
  Testfälle, nicht Reviewer-Scratch-Läufe. *Erwartet, zu belegen durch:* die
  Testfunktionen im Diff. **Ausgang:** **entfallen** — die Phase `Backfill-Boundary (leere
  Tabelle, zweiter Antrag)` ist eine committete Runner-Phase (`abdeckung_declare`) und
  druckt in beiden Legs von Lauf `36065957210` „belegt“ (leere Tabelle `completed` mit 0
  Zeilen ohne Transaktion, zweiter Antrag `failed`); Beleg: der Lauf, gemessen mit
  `gh run view 36065957210 --log`.
- **Stiller Ausschluss aus dem Runner** (`BEO-PGC/test-runner-stiller-ausschluss`,
  offen, 2×): eine neue `TestE2E*`-Funktion ohne `-run`-Muster läuft nie.
  *Erwartet, zu belegen durch:* die `-v`-Ausgabe des Laufs. Ein weiterer
  Auftritt der Klasse erreicht 3×. **Ausgang:** **entfallen** — die einzige neue
  `TestE2E*`-Funktion, `TestE2EBackfillReplayInvariant`, steht in der `-v`-Ausgabe beider
  Legs von Lauf `36065957210` (`=== RUN` und `--- PASS`); die sechs übrigen Rundläufe sind
  deklarierte Runner-Phasen. Kein Auftritt der Klasse, der Zähler des Eintrags bleibt 2×.
- **Geteilter Zustand zwischen Rundläufen** (`BEO-PGC/test-isolation-geteilter-zustand`,
  offen, 1×): eigene Tabellennamen und eigene Quelle je Phase; die
  Bestands-Tabelle wird nach dem Lauf abgeräumt. *Erwartet, zu belegen durch:*
  Lesen der Phasen. **Ausgang:** **entfallen** — der Review las jede Phase mit eigener
  Tabelle (Review-Report, Negativbefunde: „jede Phase mit eigener Tabelle“), und beide
  Legs von `e2e.yml` Lauf `36065957210` fuhren die sieben Rundläufe nacheinander ohne
  Kollision.
- **Laufzeit des erweiterten Testpakets.** `e2e.yml` fährt `make test-integration`
  je PostgreSQL-Version der Matrix; der Zuwachs verlängert jeden Lauf. Der
  Workflow bleibt unverändert ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht). *Erwartet, zu
  belegen durch:* die lokale Laufzeit vor und nach dem Zug im Bericht (mit
  Lauf-Ursprung); eine Aussage über den GitHub-Runner ist erst nach dem ersten
  Push-Lauf möglich und bleibt **weiter offen**, falls ein Timeout auftritt.
  **Ausgang:** **entfallen** — der Schritt „Compose-Integrationstest (Black-Box-E2E)“
  lief in `e2e.yml` Lauf `36065957210` (Kopf `43137ebf`) 7 min 15 s (PostgreSQL 17,
  22:11:35 bis 22:18:50) und 7 min 20 s (PostgreSQL 18, 22:11:37 bis 22:18:57) bei
  `timeout-minutes: 60`, beide `success` (gemessen: `gh api
  repos/pt9912/pg-change-feed/actions/runs/36065957210/jobs`, Zeitstempel je Schritt).
  Aussagegrenze: ein Lauf je Leg; eine „Laufzeit vor dem Zug“ ist auf dem Runner nicht
  gemessen (kein Lauf vor dem Zug mit diesem Schritt im selben Aufbau).
- **Die Startposition** ist ein gemessener Wert mit Ursprung (der Lauf); eine
  Übernahme aus [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) oder aus dem Gedächtnis ist nicht zulässig
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). *Erwartet, zu belegen durch:* das Handbuch nennt
  den Lauf. **Ausgang:** **entfallen** — Handbuch §4, Punkt „Startposition eines neuen
  Consumers“ nennt den Lauf von `make test-integration` (Phase Backfill-Startposition);
  der Verifikations-Report bestätigt die Werte (Position 0, `acknowledged` `false`, 5
  Backfill-Changes, 0 hinter der bestätigten Position) gegen die gedruckte Zeile beider
  Legs von `e2e.yml` Lauf `36065957210`.
- **Umschreiben der Tabelle im Fenster zwischen Export und Sperre** (gemessen,
  [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md) §Gemessen,
  PostgreSQL 17.11 und 18.6): ohne Sperre und Filenode-Vergleich endet der Run bei
  `ALTER TABLE … ALTER COLUMN name TYPE varchar(64)` im Fenster `completed` mit
  `rows_copied` 0 bei drei Zeilen der Tabelle; die neu geschriebene Datei ist für den
  älteren Snapshot leer (die umschreibenden Formen von `ALTER TABLE` sind nach der
  PostgreSQL-Dokumentation nicht MVCC-sicher). Der Ablauf des Imports trägt Sperre und
  Vergleich; das Umschreiben endet den Run `failed`/`transient`, und die Tests binden
  das Verhalten in drei Tiers (siehe §3). Fehlalarme bei `VACUUM FULL`, `CLUSTER` und
  `TRUNCATE` sind akzeptiert (`ADR-0118` Festlegung 4). Die Dauer des Fensters ist nicht
  gemessen. Die Wirkung der Sperre auf andere Zugriffe, solange eine DDL auf sie wartet,
  ist gemessen (Review-Report F-1, PostgreSQL 18): ein `INSERT`, ein `SELECT` auf die
  Tabelle und die Abfrage von `pg_publication_tables` liefen je in ein 4-s-Limit; das
  Handbuch nennt Schreiber, Leser und Administration. **Ausgang:** **eingetreten** → in
  diesem Slice behoben ([`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md),
  Fixrunde; Berichtigung der Aussagen zu Sperr-Warteschlange und `RENAME COLUMN`:
  [`ADR-0119`](../../adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md)); kein Carveout
  und kein Folge-Slice nötig. Die Fehlalarm-Frage ist akzeptiertes Negativ der ADR
  (Festlegung 4) und an deren Re-Evaluierungs-Trigger gebunden, kein offenes Risiko dieses
  Slice.
- **Wartegrenze der Lesesperre** (Review F-6): der Run wartet an der Sperranweisung,
  solange eine fremde Transaktion `ACCESS EXCLUSIVE` hält; der Kontext, den der Worker
  übergibt, trägt kein Zeitlimit, und `ADR-0118` Festlegung 1 nimmt das an. Gebunden ist
  das Ende des Wartens durch einen ablaufenden Kontext
  (`TestImportLockWaitEndsWithTheContext`, Klasse `transient`, keine Sitzung danach); das
  Handbuch nennt die Gegenrichtung. Eine Zeitgrenze (`lock_timeout`) ist eine
  Designänderung und in diesem Slice nicht umgesetzt. **Ausgang:** *weiter offen* →
  Beobachtungs-Register `BEO-PGC/lesesperre-ohne-zeitgrenze` (`observations/`).
- **Zeitannahmen der Haltepunkte auf dem GitHub-Runner** (Review F-9): die Pause des
  Feed-Containers bleibt unter 1000 ms, der Lauf scheitert sichtbar darüber (lokal
  gemessen 377 ms und 248 ms in den zwei Läufen der Phase DDL-Fenster, Gesamtlauf 4 min 39 s, Lauf-Ursprung: der
  Lauf dieser Fixrunde); `e2e.yml` ist unverändert. **Ausgang:** **entfallen** — der reale
  Post-Push-Lauf von `e2e.yml` ([`AGENTS.md`](../../../../AGENTS.md) §3.10), Lauf
  `36065957210` am Kopf `43137ebf`, beide Legs `success`; gedruckt „Pause des
  Feed-Containers“ 110 ms und 163 ms (PostgreSQL 17, Zeilen `DROP COLUMN` und
  Umschreiben) sowie 245 ms und 173 ms (PostgreSQL 18), Grenze im Runner 1000 ms;
  `ci` (Lauf `36065957245`) und `examples` (Lauf `36065957304`) desselben Kopfes ebenfalls
  `success` (gemessen: `gh run view 36065957210 --json jobs`, `gh run list --commit`,
  `gh run view --log`). Aussagegrenze: ein Lauf je Leg, kein Beweis gegen Flackern; ein
  späteres Timeout öffnet das Risiko erneut.
- **Die Abdeckungs-Zeilen-Anker** verschieben sich mit jeder Einfügung oberhalb
  bestehender Phasen; die regenerierte Datei wird committet. **Ausgang:** **entfallen** —
  `docs/user/e2e-abdeckung.md` ist regeneriert und committet; der Runner meldet in beiden
  Legs von `e2e.yml` Lauf `36065957210` „E2E-Abdeckungstabelle unverändert —
  `docs/user/e2e-abdeckung.md` entspricht dem Quelltext-Stand“ (gemessen an diesem Lauf).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Rollen-Kette lief in getrennten Kontexten und fand, was kein
  Gate las. Die E2E-Phase DDL-Fenster fand den Vertragsbruch des Slice: ein Umschreiben der
  Tabelle zwischen Snapshot-Export und Import ließ den Run `completed` mit `rows_copied` 0 bei
  drei vorhandenen Zeilen enden (gemessen, [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md)
  §Gemessen); die Tests zum Umschreiben im Fenster in `make test` und `make test-replication` entstehen erst
  mit der Fixrunde (§3). Das
  Architect-Verdikt setzte Lesesperre und Filenode-Vergleich als Fixrunde in diesen Slice; jede
  Zusage der Fixrunde färbt bei einer Mutation ihrer Eingabeseite rot, in allen drei Tiers
  (Verifikations-Report §4: sechs rote Mutationen, eine äquivalente grün). Der Review (0 HIGH · 1
  MEDIUM · 5 LOW · 4 INFO, Summary des Reports) fand F-1 durch Nachmessen der Sperr-Warteschlange an
  einem Wegwerf-Container (Schreiber und Publication-Abfrage stauen sich, nicht nur Leser), F-3 durch
  Nachzählen des Suchlauf-Felds und F-4 durch Lesen einer Assertion an ihrer Eingabeseite; der Verifier
  reproduzierte F-1 selbst und bestätigte mit `make gates`, `make test` und `make test-replication`
  (je Exit 0) und dem realen Post-Push-Lauf: `e2e.yml` Lauf `36065957210` am Kopf `43137ebf`, beide
  Legs `success` (Schritt „Compose-Integrationstest“ 7 min 15 s und 7 min 20 s, Pausen des Haltepunkts
  110/163 ms und 245/173 ms gegen die Grenze 1000 ms; gemessen an diesem Lauf, `gh run view
  36065957210 --json jobs` und Log, Ursprung je Zahl in §6). Es ist ein Lauf je Leg: der Beleg sagt
  „hält auf dem Runner“, nicht „flackert nie“.
- **Was ging anders als geplant:** (1) Der Slice ist nicht „nur Test“: die Fixrunde legt Produktionscode an
  (`postgressnapshot/snapshot.go`, `snapshotlogic/logic.go`, `usecase/backfill/service.go`); der Diff
  umfasst 18 Commits und 23 Dateien (+2910/−370, **übernommen** aus dem Verifikations-Report, Range
  `e7df5619..43137ebf`). (2) Der Ansatz „Tabellensperre als Haltepunkt“ trägt nicht: `pg_publication_tables`
  wartet auf jede Sperre der Tabelle, und eine `ACCESS EXCLUSIVE`-Sperre weist ihrer Transaktion eine
  Kennung zu, auf die die Slot-Anlage wartet (§3, Ansatz-Ergebnis, gemessen an PostgreSQL 18); Ersatz sind
  der unbestätigte Schlüssel im zweiten Block und `docker pause`. (3) Die Fixrunde lief ohne eigenen
  Review-Report — benannte Grenze (V-3): der Verifier las den Fixrunden-Diff, mutierte ihn und
  reproduzierte F-1 selbst, ein Reviewer-Durchgang über den Fixrunden-Diff ist nicht gefahren.
  (4) Der Plan nennt sechs Fremdobjekte, `knownForeignObjects` führt sieben (Plan-Drift, in §3 benannt).
- **Verifier-Beobachtungen (V-1 bis V-6):** *V-1* (LOW): die Ortswahl der Replay-Invariante ist eine
  benannte Entscheidung des Planners (§3, letzte Zeile der Tabelle mit Begründung und Anker); kein Architect-Verdikt
  trägt sie, und ein Tier-Beleg ist nicht geliefert — die Fitness-Function-Zeile von
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt einen Beleg, den kein Test trägt
  (*benannte Lücke*, Adresse: Architect im Lese-Schritt der Closure von
  [welle-backfill-bestand](../welle-backfill-bestand.md)). *V-2* (LOW): die Mutationen je Negative-Kriterium
  sind im Repo nicht belegt — der Bericht des Implementers liegt nicht im Repo, kein committetes Artefakt
  nennt sie, der Verifier fuhr für die Negative-Phase keine E2E-Mutation (Lauf-Kosten je rund 5 Minuten;
  Verifikations-Report §12). Die Closure fährt sie nicht nach; die Runner-Assertions der Negative-Phase sind
  gelesen und im grünen Post-Push-Lauf bestätigt, nicht gemutet — eine benannte Grenze, keine Erfüllung der
  Zusage. *V-3* (INFO) siehe „Was ging anders“ (3). *V-4* (INFO): der `ctx`-Parameter von `ClassifyLock`
  ist wirkungsgleich zu `Classify` (äquivalente Mutation), die Kontext-Bindung trägt
  `TestImportLockWaitEndsWithTheContext`; keine Aktion. *V-5* (INFO): die zwei überholten Aussagen von
  `ADR-0118` sind mit [`ADR-0119`](../../adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md) berichtigt.
  *V-6* (INFO): Tier-Nebenwirkungen (`tools/schema/plan.yaml`, `docs/user/e2e-abdeckung.md`) sind
  zurückgenommen, keine Aktion.
- **Steering-Loop-Eintrag (Lerneintrag):** *Lerneintrag (Ausprägung, in der Phase verkörpert):* ein E2E-Beleg
  für ein DDL-Fenster fährt nicht nur eine Form, die sichtbar scheitert (`DROP COLUMN`, `failed`/`storage`),
  sondern auch eine Rewrite-DDL (`ALTER COLUMN … TYPE`): die sichtbar scheiternde Form belegt das Fenster
  nur für Fehler mit Meldung, die stille Form deckt den Verlust ohne Meldung auf. Die Phase DDL-Fenster fährt beide Läufe
  über `bf_ddl_window`. *Neuer Sensor:* die Phase DDL-Fenster (zwei Läufe, ein neuer Antrag danach) und im
  Store-Tier `TestRewriteInWindowIsTransient`, `TestNoRewriteInWindowReadsTheSnapshot`,
  `TestImportWaitsForExclusiveLockAndThenAborts`, `TestImportLockWaitEndsWithTheContext`; gebunden an die
  Eingabeseite (Mutationen M1 bis M5 und M-E2E rot, Verifikations-Report §4). *Benannte Spec-Lücke:* die
  Fitness-Function-Zeile von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt für die
  Replay-Invariante den Tier `make test-replication`; der Beleg liegt im E2E
  (`TestE2EBackfillReplayInvariant`), der Tier-Beleg fehlt (Adresse: V-1 oben). *Geschärfte Regel
  (Kandidat, nicht entschieden):* eine DoD-Zusage an eine nachgelagerte Rolle („der Architect bestätigt“,
  „je Kriterium eine Mutation im Bericht“) nennt den committeten Ort ihrer Erfüllung
  (`BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker`, 1×; ein zweites Auftreten liegt unter der
  Schwelle, der Träger ist Sache des Lese-Schritts). *Ansatz-Ergebnis:* der Haltepunkt „Tabellensperre“ trägt
  nicht (§3); der Ersatz ist gemessen und belegt. *Sensor-Dokumente:* Nenner und gedeckte Zahl tragen ihren
  Lauf — Coverage-Gate: `make coverage-gate` (Closure-Lauf am Stand `cf7f2d02`, Exit 0), gedruckt
  `total: (statements) 83.1%` und `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%`, dedupliziert
  2112 von 2541 = 83,12 % (Awk über `/out/coverage.out` des Images, **abgeleitet**), darunter `snapshotlogic`
  54 von 54; DB-Adapter-Coverage: `make test-replication` (Closure-Lauf, PostgreSQL 18, Exit 0) mit dem
  Store-Profil des Verifikations-Laufs (der Slice berührt `postgresstorage` nicht:
  `git diff --stat e7df5619..HEAD` über dessen Verzeichnis ist leer), gedruckt `DB-Adapter-Coverage: 82.13%
  (gedeckt 850 von 1035 Statements; Profile gemergt: store,replication)`, dieselbe Zeile wie in beiden Legs
  des Post-Push-Laufs.
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei `evidence/slice-backfill-e2e.md`,
  Zähler = Zahl der Dateien (real ausgezählt: `ls …/evidence | wc -l`). *Bestehende Klassen:*
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (F-4) **10×**, `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
  (F-3) **18×**, `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (F-2) **30×** und
  `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (F-1, F-7, V-5) **4×** stehen **über** der Schwelle 3×; ihr Ausgang
  gehört dem Lese-Schritt der Closure von [welle-backfill-bestand](../welle-backfill-bestand.md) (die ersten drei sind
  verkörpert, der vierte ohne zugewiesenen Ausgang; die drei ersten erweitern den Bestand um eine weitere
  Ausprägung, die state-Dateien tragen den Vermerk). `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar`
  (F-5) **2×** und `BEO-PGC/formatierungs-drift-ohne-gate` (F-10) **2×**, offen; `BEO-PGC/lesesperre-ohne-zeitgrenze` (F-6)
  **1×**, offen (Risiko §6 „Wartegrenze der Lesesperre“). *Neue Klassen:*
  `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker` (V-1, V-2) **1×**, offen;
  `BEO-PGC/run-fehlertext-traegt-klasse-doppelt` (F-8) **1×**, verkörpert (`failureText`,
  `TestExecuteFailureTextCarriesClassOnce`). F-9 ist kein Register-Anfall: das Risiko hat seinen Ausgang in §6.
- **Folge-Slices:** keine neuen. Die offenen Slices der Welle (`slice-backfill-bench-richtgroesse`,
  `slice-backfill-sdk-origin`) sind von diesem Slice nicht abhängig; der Tier-Beleg der Replay-Invariante
  ist eine Entscheidung des Architects im Lese-Schritt der Welle-Closure, ein Slice entsteht erst mit ihr.
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Entfallen:* Die Replay-Invariante gilt real nicht ·
  Nichtdeterministischer Abbruch · Randfall-Belege nur beim Reviewer · Stiller Ausschluss aus dem Runner ·
  Geteilter Zustand zwischen Rundläufen · Laufzeit des erweiterten Testpakets · Die Startposition ·
  Zeitannahmen der Haltepunkte auf dem GitHub-Runner · Die Abdeckungs-Zeilen-Anker. *Eingetreten:*
  Umschreiben der Tabelle im Fenster — in diesem Slice behoben ([`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md),
  [`ADR-0119`](../../adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md)). *Weiter offen:* Wartegrenze der
  Lesesperre → Register `BEO-PGC/lesesperre-ohne-zeitgrenze`.
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); die E2E-Werkzeuge unter `tools/harness/` und `test/`
sind keine eigene Sub-Area — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, einschlägig — DoD),
`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt` (offen, 1×, einschlägig —
Risiko §6), `BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×, Risiko §6),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3× — Zeitfenster
in Polls tragen ihre Begründung), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
(verkörpert, 6×, DoD), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(verkörpert, 13×, Startposition mit Lauf),
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 8× — die
Zeile in `harness/README.md` beschreibt nur, was der Lauf gefahren hat),
`BEO-PGC/anforderung-ohne-erkennbaren-nachweis` (offen, 1×, einschlägig: dieser
Slice löst die Waise `LH-FA-CAP-009` auf), `BEO-PGC/limit-fortsetzung-innerhalb-einer-position`
(offen, 0×, einschlägig — die Belege lesen den Bestand ohne `Limit`),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3),
`BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert, 7× — nicht
einschlägig: kein Workflow-Zug).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
