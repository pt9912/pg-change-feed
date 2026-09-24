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
`queued`-Zeile beim Prozessstart).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md), [`SPEC-022`](../../../../spec/pflichtenheft.md), [`SPEC-029`](../../../../spec/pflichtenheft.md) — gelesen
als Vertrag der Belege, nicht geändert.

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
      dort Snapshot, Schreiber und WAL-Pfad komponiert laufen — der Architect
      bestätigt die Ortswahl im Review oder verlangt zusätzlich den Tier-Beleg.
- [x] Negative: `docker kill` im laufenden Run → nach dem Neustart `interrupted`,
      keine sichtbare Change des Runs, erneuter Antrag erreicht `completed`, Bestand
      einmal und vollständig; eine zum Abbruchzeitpunkt `queued` wartende Zeile
      überlebt den Neustart und wird ausgeführt (`completed`, Bestand ihrer Tabelle
      lesbar). *Zu belegen durch:* `make test-integration`; **jedes** dieser
      Kriterien trägt je eine Mutation im Bericht (die Prüfung gegen die
      Eingabe gelenkt, der Lauf färbt rot —
      `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert).
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
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: Handbuch §4 (Abschnitt „Bestand als Backfill überführen“) trägt die gemessene Startposition mit Lauf-Ursprung; die Änderungshistorie eine Zeile.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | unverändert | die Go-Hälfte liegt in der neuen Datei `test/integration/backfill_e2e_test.go` (eigene Helfer; der Abdeckungs-Erzeuger liest alle `.go`-Dateien des Pakets). Go-seitig laufen nur die Replay-Anwendung und der Bestandsvergleich, die übrigen Belege sind Runner-Phasen. |
| `test/integration/backfill_e2e_test.go` | neu | `TestE2EBackfillReplayInvariant` (Boundary: drei nebenläufige `INSERT`/`UPDATE`/`DELETE`-Schreiber vor und nach der Snapshot-Position, Replay des Logs gegen den Quellstand, Überlappung der beiden Seiten als Abbruchbedingung) samt Hilfsfunktionen; ausschließlich externe Wege (SQL, keine Importe aus `internal/**`). |
| `tools/harness/run-integration-tests.sh` | update | sieben Rundläufe (sechs Runner-Phasen mit `abdeckung_declare` auf `LH-FA-CAP-009`, ein eigener Go-Aufruf): Happy Path, Schema-Version, Startposition, Boundary (leere Tabelle, zweiter Antrag), Replay-Invariante, DDL-Fenster, Negative (`docker kill`, `queued`-Aufnahme); drei Haltepunkte (offene Schreibtransaktion, unbestätigter Schlüssel im zweiten Block, `docker pause`); der Tabellensperren-Ansatz ist nicht realisiert (siehe Ansatz-Ergebnis). |
| Kommentar zum zweiten `make schema-rollout`-Lauf in `tools/harness/run-integration-tests.sh` (gemeldet von `slice-backfill-change-origin`, Review F-7) | update | der Kommentar („… real blockierten zweiten `make schema-rollout`-Lauf (… Exit 8 auf vier Fremdobjekten …)") beschrieb den Stand der Wache nicht: er nennt jetzt den Ist-Zustand (der Tausch belegt das Upgrade des Prozesses; den zweiten Rollout belegt `tools/harness/run-schema-rollout-guard-test.sh`, idempotent über die Fremdobjekte des Guards). **Plan-Drift:** der Plan nennt sechs Fremdobjekte, `knownForeignObjects` in `tools/schema/rolloutguard/guard.go` führt sieben (fünf Funktionen, zwei Views); der Kommentar nennt sieben. Die Zeilen-Anker von `docs/user/e2e-abdeckung.md` verschieben sich hier ohnehin (Erzeugnis, mitcommittet) — der Slice, der den Kommentar ändert, ist der Slice, der die Anker regeneriert. |
| `tools/harness/httpclient/main.go` | update | geliefert: `origin` in der READ-Zeile (geprüft gegen `wal`/`backfill`, der Beleg der Feldform am Wire) und die Modi `changes` (`GET /changes` ohne Registrierung) und `position` (`GET /consumers/position`). |
| `compose.yaml` | prüfen | unverändert: `max_replication_slots=10` und `max_wal_senders=10` tragen die Reserve für den temporären Slot; `wal_sender_timeout=2000` begrenzt die Pause des `docker pause`-Haltepunkts (der Runner hält sie unter der Hälfte); die drei DSNs sind Superuser-Verbindungen, die Rollen-Zusagen tragen die Store-Tier-Tests (Architect-Verdikt, „Akzeptiertes Negativ"). |
| `internal/application/usecase/backfill/service.go` | update | Grenze-Kommentar an `currentVersion` (Auflage aus dem Architect-Verdikt, [`ADR-0116`](../../adr/0116-backfill-schema-version-referenz-reichweite.md) Folgepflicht 1; kein Verhaltens-Diff). |
| `internal/application/usecase/backfill/service_test.go` | update | `TestExecuteMarksRunningBeforeOpeningSnapshot` bindet die Reihenfolge `MarkRunning` vor `OpenSnapshot` (Übergabe aus `slice-backfill-run-usecase`). |
| `docs/user/e2e-abdeckung.md` | regeneriert | Erzeugnis des Runners, mitcommittet. |
| `docs/user/benutzerhandbuch.md` | update | geliefert: Punkt „Startposition eines neuen Consumers“ (Ursprung: der Lauf), `rows_copied` im Zustand `interrupted`, Version 1.50 samt Änderungshistorie. |
| `harness/README.md` §Sensors | update | die Zeile `make test-integration` trägt die sieben Backfill-Rundläufe. |

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
- **Fenster zwischen Export und Cursor.** `DROP COLUMN` im Fenster: `DECLARE` scheitert `42703`, der Run endet `failed` mit Klasse `storage`, ohne Change (Phase DDL-Fenster). Ein Tabellen-Rewrite (`ALTER COLUMN … TYPE` statt `DROP COLUMN` in derselben Phase, Wegwerf-Lauf, nicht committet): der Run endet `completed` mit `rows_copied` 0 bei drei Zeilen der Tabelle — siehe §6.
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
- **Fenster zwischen Export und Cursor.** Spaltenliste und `DECLARE` laufen ohne
  Tabellensperre im importierten Snapshot; ein gleichzeitiges `ALTER TABLE` in diesem
  Fenster lässt den `DECLARE` mit `storage` (`42703`) scheitern oder — bei einem
  Tabellen-Rewrite — Katalog-Stand und Snapshot auseinanderlaufen (aus dem Verhalten
  von `SET TRANSACTION SNAPSHOT` abgeleitet, **nicht gemessen**). Ist ein Haltepunkt
  dafür deterministisch herstellbar, misst der Slice das Fenster; sonst nennt der
  Bericht es als ungemessen.

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
| Beschreibungen von `make test-integration` außerhalb der README | `git grep -l 'test-integration' <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` und Lesen der beschreibenden Treffer (`Makefile`, `harness/sensors/docs-check.md`, `AGENTS.md`, `examples/README.md`, `docs/plan/planning/in-progress/roadmap.md`) | Parent: 124 Zeilen in 57 Dateien; Diff-Stand: 127 Zeilen in 59 Dateien (zusätzlich `docs/user/benutzerhandbuch.md` mit 2 und `test/integration/backfill_e2e_test.go` mit 1 Zeile); beide Stände ohne die Plan-Datei dieses Slice, die sich selbst zitiert. Beschreibend: `Makefile` (Hilfezeile nennt drei Beispiele ohne Vollständigkeitsanspruch), `harness/sensors/docs-check.md` (Erzeugung der Abdeckungstabelle, unverändert wahr); die übrigen sind Verweise auf das Ziel ohne Beschreibung des Inhalts | unverändert |
| Ausgabe und Modi des Wegwerf-Clients | `git grep -n httpclient <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` und Lesen der Treffer außerhalb von `tools/` | Parent: 39 Treffer-Zeilen in 14 Dateien, Diff-Stand: 48 Zeilen in 14 Dateien (die Mehrzeilen liegen in `tools/harness/`); beide ohne die Plan-Datei dieses Slice; beschreibende Träger der Ausgabeform: `examples/http-client/main.go` (Verweis „Wegwerf-Client als Belegträger“), `examples/README.md` (Abgrenzung), `docs/plan/planning/open/slice-transformationen-e2e-wirkung.md` (Plan-Zeile „geben Schlüssel und Werte der Row Images aus“) — keine nennt die READ-Zeile oder deren Feldliste | unverändert |
| Zählwörter zu den bewegten Eigenschaften | `git grep -n -E '(zwei\|drei\|vier\|fünf\|sechs\|sieben\|acht\|beiden\|beide) (Backfill-\|weitere \|zusätzliche )?(Rundl\|Belege\|Phasen\|Haltepunkt\|Fremdobjekt\|Waise)' <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/adr'` (ohne die Plan-Datei dieses Slice) | Parent: 11 Treffer; Diff-Stand: 11 Treffer. Entfallen: „vier Fremdobjekten“ im Runner-Kommentar (ersetzt); hinzugekommen: „sieben Backfill-Rundläufe“ in `harness/README.md`. Die übrigen zehn betreffen Kurs-Dokumente (`.claude/commands/close-welle.md`, `.harness/baseline`), die SDK-Werkzeuge (`harness/mk/sdk.mk`, „vier Phasen“) und die Coverage-/Replikations-Skripte (zwei Phasen) und sind von diesem Slice nicht berührt | unverändert |
| RTM-Träger | `make doc-trace` | siehe Zeile zu den Waisen: `LH-FA-CAP-009` trägt `E2E`, 2 Waisen | Deklarations-Anker tragen |
| CI-Träger der Läufe | Lesen von `.github/workflows/e2e.yml` (Trigger, Matrix, `timeout-minutes`); `git diff --stat e7df5619 -- .github` | Trigger Pull Request und Push, Matrix über die PostgreSQL-Versionen 17 und 18, `timeout-minutes: 60`; `.github` ohne Diff | unverändert; die Laufzeit-Frage steht in §6 |

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
  benennt das. *Erwartet, zu belegen durch:* der Lauf. **Ausgang:** *(bei Closure)*
- **Nichtdeterministischer Abbruch** — siehe Ansatz-Vorschlag; ohne Haltepunkt
  bliebe der Test flackernd. *Erwartet, zu belegen durch:* mehrere
  Wiederholungen im Bericht. **Ausgang:** *(bei Closure)*
- **Randfall-Belege nur beim Reviewer** (`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt`,
  offen, 1×): die leere Tabelle und der zweite Antrag sind committete
  Testfälle, nicht Reviewer-Scratch-Läufe. *Erwartet, zu belegen durch:* die
  Testfunktionen im Diff. **Ausgang:** *(bei Closure)*
- **Stiller Ausschluss aus dem Runner** (`BEO-PGC/test-runner-stiller-ausschluss`,
  offen, 2×): eine neue `TestE2E*`-Funktion ohne `-run`-Muster läuft nie.
  *Erwartet, zu belegen durch:* die `-v`-Ausgabe des Laufs. Ein weiterer
  Auftritt der Klasse erreicht 3×. **Ausgang:** *(bei Closure)*
- **Geteilter Zustand zwischen Rundläufen** (`BEO-PGC/test-isolation-geteilter-zustand`,
  offen, 1×): eigene Tabellennamen und eigene Quelle je Phase; die
  Bestands-Tabelle wird nach dem Lauf abgeräumt. *Erwartet, zu belegen durch:*
  Lesen der Phasen. **Ausgang:** *(bei Closure)*
- **Laufzeit des erweiterten Testpakets.** `e2e.yml` fährt `make test-integration`
  je PostgreSQL-Version der Matrix; der Zuwachs verlängert jeden Lauf. Der
  Workflow bleibt unverändert ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht). *Erwartet, zu
  belegen durch:* die lokale Laufzeit vor und nach dem Zug im Bericht (mit
  Lauf-Ursprung); eine Aussage über den GitHub-Runner ist erst nach dem ersten
  Push-Lauf möglich und bleibt **weiter offen**, falls ein Timeout auftritt.
  **Ausgang:** *(bei Closure)*
- **Die Startposition** ist ein gemessener Wert mit Ursprung (der Lauf); eine
  Übernahme aus [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) oder aus dem Gedächtnis ist nicht zulässig
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). *Erwartet, zu belegen durch:* das Handbuch nennt
  den Lauf. **Ausgang:** *(bei Closure)*
- **Tabellen-Rewrite im Fenster zwischen Export und Cursor** (gemessen, Wegwerf-Lauf
  der Phase DDL-Fenster mit `ALTER TABLE … ALTER COLUMN name TYPE varchar(64)` statt
  `DROP COLUMN`, PostgreSQL 18): der Run endet `completed` mit `rows_copied` 0 bei drei
  Zeilen der Tabelle; die neu geschriebene Tabelle ist für den älteren Snapshot leer
  (die tabellenumschreibenden Formen von `ALTER TABLE` sind nach der
  PostgreSQL-Dokumentation nicht MVCC-sicher). Ein Verlust, den der Run nicht meldet —
  ein Befund gegen die Aussage „ohne Verlust“ von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) für dieses
  Fenster, kein Testfehler; der Slice bindet ihn nicht als Test (ein Test der Fehlform
  bände den Fehler), er meldet ihn an den Architect (Folge-ADR oder Schärfung). Das
  Fenster liegt zwischen Snapshot-Export und dem Cursor-Aufbau des Runs; die Dauer des
  Fensters ist nicht gemessen. Ab dem Cursor-Aufbau hält der Run eine
  `ACCESS SHARE`-Sperre der Tabelle, die ein `ALTER TABLE` bis zum Ende des Runs warten
  lässt (aus der Sperrmatrix hergeleitet, nicht gemessen). **Ausgang:** *(bei Closure)*
- **Die Abdeckungs-Zeilen-Anker** verschieben sich mit jeder Einfügung oberhalb
  bestehender Phasen; die regenerierte Datei wird committet. **Ausgang:** *(bei
  Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
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
