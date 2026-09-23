# Slice backfill-snapshot-reader: Snapshot-Leser — Port `TableSnapshotPort` und Driven Adapter: temporärer Slot mit Export-Snapshot, Import, Cursor-Blöcke, GUC-Parität

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Bestand über denselben Lesezugriffsweg),
[`LH-FA-CAP-008`](../../../../spec/lastenheft.md) (Row-Image-Abwesenheit), [`LH-FA-CAP-006.a`](../../../../spec/pflichtenheft.md) (keine
unbegrenzte RAM-Haltung), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 1 (Mechanismus: Slot-Snapshot),
[`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Ports nach Fähigkeiten), [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) (Nahtform der
Treiber-Hülle), [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 3 (Messgegenstand der DB-Adapter-Coverage),
[`ADR-0030`](../../adr/0030-testpyramide.md) (Testpyramide).

**Berührte Spec-Stellen:** [`SPEC-012`](../../../../spec/pflichtenheft.md) (PostgreSQL 17 und 18), [`ARC-004`](../../../../spec/architecture.md)
(Outbound Port), [`ARC-006`](../../../../spec/architecture.md) (Driven Adapter), [`ARC-008`](../../../../spec/architecture.md) (PostgreSQL Logical
Replication) — gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Outbound Port `TableSnapshotPort` und sein Driven Adapter
liefern den Bestand einer Tabelle als konsistenten Snapshot: die
Replication-Verbindung legt `CREATE_REPLICATION_SLOT cdc_bf_<run> TEMPORARY
LOGICAL pgoutput EXPORT_SNAPSHOT` an und liefert `consistent_point` **X** und
`snapshot_name`; eine reguläre Verbindung öffnet `BEGIN ISOLATION LEVEL
REPEATABLE READ READ ONLY; SET TRANSACTION SNAPSHOT '…'`, danach endet die
Replication-Verbindung; im Snapshot liest der Adapter die Spaltenliste
(`pg_attribute`, `attnum > 0`, nicht gelöscht, **nicht generiert**) und die
Zeilen über einen `NO SCROLL`-Cursor in Blöcken von `B` Zeilen mit `col::text`
je Spalte (Startwert `B` = 1.000, der Slice legt ihn fest). Die Anlage trägt ein
Zeitlimit; sein Ablauf endet als Fehlerklasse `transient`. Alle Verbindungen
laufen über `CDC_CAPTURE_DSN`. Als Katalog-Fähigkeit derselben Quelle liest der
Adapter außerdem die **geschätzte Zeilenzahl** einer Tabelle
(`pg_class.reltuples`; ein negativer Wert heißt „unbekannt" — erwartet, ungemessen)
für den Antrag des Runs.

Der Slice trägt die real gemessenen Befunde M1–M5 aus [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) als Tests und
die **Bild-Parität**: dasselbe Zeilenbild über WAL-Pfad und Backfill-Pfad
byte-gleich — mit der gemeinsamen Funktion aus `row-image-gemeinsam`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Run** (Vorbedingungen, Fail-closed, Schreiben, Wecksignal) — `run-usecase`;
  der Snapshot-Leser kennt weder Run-Zustand noch Store.
- **Persistenz und Position im Store** — der Adapter liefert `X` als
  Rohwert der LSN; die Abbildung auf die Position ([`ADR-0005`](../../adr/0005-sourceposition-abstrahiert-lsn.md)) und ihre
  Ablage gehören dem Run und dem Store.
- **Die Replay-Invariante und der Abbruch mitten im Lauf** — sie brauchen den
  komponierten Run und gehören dem `e2e`-Slice.
- **Die Betriebs-Vorbedingungen** (`SELECT` auf die Quelltabelle,
  `max_replication_slots`/`max_wal_senders`-Reserve) — Doku im
  `sql-administration`-Slice; dieser Slice belegt nur, dass ihr Fehlen als
  Fehlerklasse (`permission`, `configuration`) sichtbar wird.

## 2. Definition of Done

- [ ] Snapshot-Träger real belegt: der Slot-Snapshot sieht die vor `X`
      committeten Zeilen und keine danach committete (M1); eine Rolle nur mit
      `LOGIN REPLICATION` ohne `SELECT` scheitert mit der Klasse `permission`
      (M2); das Zeitlimit der Slot-Anlage gegen eine offene Schreibtransaktion
      endet als `transient` (M3); nach Import und Ende der Replication-Verbindung
      ist der temporäre Slot weg und der Cursor liest weiter den Snapshot-Stand
      (M4). *Zu belegen durch:* `make test-replication`; die Export-Syntax
      zusätzlich gegen PostgreSQL 17 (`PG_TEST_IMAGE` auf den 17er-Digest der
      CI-Matrix in `.github/workflows/e2e.yml`), damit beide Major-Versionen
      belegt sind ([`SPEC-012`](../../../../spec/pflichtenheft.md)).
- [ ] Blöcke und Spalten: die Zeilen kommen in Blöcken von höchstens `B` (Test
      mit mehr als einem Block und exakt an der Blockgrenze), die Spaltenliste
      lässt generierte und gelöschte Spalten aus, eine leere Tabelle liefert
      keinen Block, die geschätzte Zeilenzahl einer frisch befüllten Tabelle ohne
      `ANALYZE` ist „unbekannt" und nach `ANALYZE` dicht an der tatsächlichen
      Zeilenzahl. *Zu belegen durch:* `make test-replication` (auf 17 und 18).
- [ ] Bild-Parität: für Zeilen mit den Typen `timestamptz`, `float8`, `bytea`,
      `interval`, `numeric`, `jsonb`, Array, `date` und einer generierten Spalte
      ist das über den Backfill-Pfad gebaute Bild byte-gleich dem über den
      WAL-Pfad gebauten — unter einer Rolle mit Sitzungs-GUC ungleich Standard
      (`timezone`, `datestyle`; M5). *Zu belegen durch:* `make test-replication`
      (der WAL-Zweig des Vergleichs läuft im Tier über Publication und Slot).
- [ ] Gate-Läufe der Messungen: das neue Paket ist dem beschlossenen Gegenstand
      zugeordnet (Start-Trigger, §4), `make coverage-gate` bleibt grün, die
      DB-Adapter-Coverage nennt ihre Zahl mit ihrem Lauf ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` §Sensors (`make test-replication` und die Coverage-Zeilen) und die Sensor-Dateien `harness/sensors/coverage-gate.md`/`harness/sensors/db-adapter-coverage.md` tragen das neue Paket, soweit der Architect-Beschluss es verlangt.
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
| `internal/application/port/outbound/tablesnapshot.go` (Arbeitsname) | neu | `TableSnapshotPort`: Snapshot öffnen (Tabelle) → Position `X`, Spalten, Block-Leser → schließen; Fähigkeits-Port ([`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md)). |
| `internal/adapters/driven/postgressnapshot/` (Arbeitsname) | neu | Treiber-Hülle (Form von [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)), Slot-Anlage, Import, Cursor, Spaltenliste, GUC-Parität; Übersetzung oberhalb der Naht netzlos prüfbar. |
| `internal/adapters/driven/postgressnapshot/*_test.go` | neu | Snapshot-Träger M1–M4, Blöcke, Spalten, Bild-Parität; DB-gebundene Tests überspringen ohne `CDC_REPLICATION_TEST_DSN` (Muster der bestehenden DB-Pakete). |
| `tools/harness/run-replication-tests.sh` | update (falls nötig) | Tier-Vorbedingungen (Rollen, zweiter Lauf gegen PostgreSQL 17). |
| `Dockerfile` (Stufe `coverage`), `tools/harness/db-coverage.sh` (`DB_COVERAGE_PKGS`) | update — **nach Architect-Beschluss** | siehe Start-Trigger: das Paket ist ohne DB-Verbindung ungeprüft und darf den netzlosen Nenner nicht füllen. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Menge der Pakete, deren Testlauf einen externen PostgreSQL voraussetzt"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Ausschluss-Muster der Coverage-Stufe | Lesen von `Dockerfile` (`grep -vE` über `postgresstorage`, `postgresack`, `replication/receive`) | *(Implementer trägt ein)* | nach Architect-Beschluss |
| Paketliste der DB-Adapter-Coverage | `grep -n 'DB_COVERAGE_PKGS' tools/harness/db-coverage.sh` | *(Implementer trägt ein)* | nach Architect-Beschluss |
| Sensor-Beschreibungen | `grep -rn 'postgresack' harness docs` | *(Implementer trägt ein)* | aktive Träger nachziehen; `Accepted` ADRs ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md), [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)) nicht ändern |
| Träger der DB-Adapter-Coverage im Workflow | Lesen von `.github/workflows/e2e.yml` auf eine Paketliste | *(Implementer trägt ein)* | trüge der Workflow die Liste, wäre das eine Workflow-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.10: realer Post-Push-Lauf vor Closure) |
| Build-Kontext | Lesen von `.dockerignore` | *(Implementer trägt ein)* | `internal/` ist freigegeben; ein Pfad außerhalb bräuchte eine Freigabe (`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad`) |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `run-usecase` **nicht** vorausgesetzt
ist (dieser Slice geht ihm voraus), `row-image-gemeinsam` in `done/` liegt, kein
anderer Slice in `in-progress/` liegt **und** ein **Architect-Verdikt** unter
`docs/reviews/` die Frage beantwortet, welchem Messgegenstand das neue Paket
zugehört: erweitern die Ausschluss-Regel der Coverage-Stufe und der Gegenstand der
DB-Adapter-Coverage ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 3) sich um dieses Paket (dann mit
Beschluss-Träger, weil die Regel dort namentlich aufzählt), oder liegt der
Snapshot-Leser in einem der bestehenden DB-Pakete. Das ist eine **Vorab-Bedingung**
und steht deshalb am Start, nicht als Rückführung
(`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, offen, 2×). Ohne den
Beschluss füllte ein Paket, dessen Tests ohne Datenbank überspringen, den Nenner
des blockierenden Gates mit ungedeckten Zeilen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Snapshot-Träger
  und Bild-Parität nicht in einem Review tragen — der abtrennbare Teil ist die
  Bild-Parität (dann eigener Slice `backfill-bild-paritaet` mit Start nach
  diesem).
- `in-progress` → `open` (blockiert): falls die Bild-Parität real **nicht**
  besteht und die GUC-Angleichung (`SET` der Sitzungs-GUC auf die Werte des
  Walsenders) eine Entscheidung über den Betriebsvertrag braucht — Architect-Frage.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test-replication` real grün
(beide Phasen) + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Gate-Zuordnung des neuen DB-Pakets** — der Start-Trigger trägt die
  Vorab-Bedingung. *Erwartet, zu belegen durch:* das Architect-Verdikt und der
  grüne `make coverage-gate` am Diff. **Ausgang:** *(bei Closure)*
- **Die GUC-Parität ist ungemessen.** [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) belegt nur, dass der
  Walsender die Rollen-Defaults trägt (M5); dass eine reguläre Verbindung mit
  demselben DSN dieselben Werte trägt, führt die ADR als „erwartet".
  *Erwartet, zu belegen durch:* der Paritätstest unter `ALTER ROLE … SET`.
  Fällt er rot aus, ist der Fallback eine explizite `SET`-Angleichung der
  Sitzung an die Walsender-Werte — dann Rückführung §4. **Ausgang:** *(bei
  Closure)*
- **Generierte Spalten in PostgreSQL 18.** Der WAL-Pfad sendet generierte
  Spalten nicht, solange die Publication sie nicht ausdrücklich veröffentlicht;
  ob das für die Publication dieses Repos unter PostgreSQL 18 gilt, ist ungemessen.
  *Erwartet, zu belegen durch:* der Paritätstest mit generierter Spalte auf 17
  und 18. **Ausgang:** *(bei Closure)*
- **Slot-Anlage wartet auf laufende Schreibtransaktionen** (M3: 5,27 s bei rund
  5 s Restlaufzeit, gemessen in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)) und braucht Reserve in
  `max_replication_slots`/`max_wal_senders`. *Erwartet, zu belegen durch:* der
  Zeitlimit-Test; die Tier-Umgebung setzt beides auf 10
  (`tools/harness/run-replication-tests.sh`). **Ausgang:** *(bei Closure)*
- **Netzlos geprüfter Code im DB-Gegenstand** (`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`,
  offen, 2×): die Übersetzungs-Schicht oberhalb der Naht ([`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)) ist
  netzlos prüfbar und hebt die DB-Zahl ohne DB-Beleg. *Erwartet, zu belegen
  durch:* der Bericht nennt Zähler und Nenner je Messung; ein weiterer Beleg
  der Klasse erreicht die Schwelle 3×. **Ausgang:** *(bei Closure)*
- **`pglogrepl` trägt die Optionen** (`CreateReplicationSlotOptions{Temporary,
  SnapshotAction}`, Ergebnis `SnapshotName`) — in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) als „gelesen"
  geführt, nicht als gefahren. *Erwartet, zu belegen durch:* M1 im Tier.
  **Ausgang:** *(bei Closure)*

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
`*`/`PGC` (Greenfield); der Snapshot-Adapter ist ein Driven Adapter wie die
bestehenden DB-Pakete, keine eigene Sub-Area — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×, einschlägig —
Start-Trigger), `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`
(offen, 2×, einschlägig — Risiko §6),
`BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (verkörpert),
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` (offen, 2×) und
`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×) —
Suchlauf-Zeile Build-Kontext, `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(verkörpert, 26×, Suchlauf §3), `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger`
(offen, 2×, einschlägig — eine Erweiterung des Messgegenstands braucht einen
Träger, Start-Trigger), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(verkörpert, 13×, Zahlen der Coverage tragen den Lauf).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
