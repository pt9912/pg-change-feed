# Slice backfill-change-origin: Feld `origin` von der Domäne bis zu den Lesewegen — `cdc.change`, View `cdc.changes`, `GET /changes`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-DAT-006`](../../../../spec/lastenheft.md) (Metadaten-Erweiterbarkeit — Boundary: gespeicherte
Changes ohne das Feld lesen sich unverändert), [`LH-FA-REA-001`](../../../../spec/lastenheft.md)
(Bereichslesen), [`LH-FA-SST-002`](../../../../spec/lastenheft.md) (SQL-Lesezugriff `cdc.changes`),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (HTTP-Lesezugriff), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2/8 (Feld
`origin`, Reichweite in Lesewegen), [`ADR-0058`](../../adr/0058-testansatz-fuenf-luecken.md) Entscheidung 2 (additive
Spalte an `cdc.change` als belegter Weg), [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md) (Changes lesen über die
HTTP-API).

**Berührte Spec-Stellen:** [`SPEC-002`](../../../../spec/pflichtenheft.md) (`cdc.change`), [`SPEC-022`](../../../../spec/pflichtenheft.md)
(`GET /changes`) — beide durch `spec-nachzug` bereits um `origin` ergänzt, hier
gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Feld `origin` (`wal` | `backfill`, ein fehlender Wert liest als
`wal`) ist von der Domäne bis zu den Lesewegen durchgängig: `model.ChangeOrigin`
(geschlossene Menge, Konstruktor-Default `wal`) und `Change.Origin`, die additive
Spalte `cdc.change.origin` (`text`, nullable, ohne Default und ohne CHECK — die
Menge erzwingt die Domäne, nicht die Datenbank), die View `cdc.changes` mit
`COALESCE(c.origin, 'wal') AS origin` als **letzter** Spalte, und das Feld
`origin` in der Antwort von `GET /changes`. Der WAL-Pfad schreibt `wal`; nichts
schreibt in diesem Slice `backfill` — dafür gibt es noch keinen Schreiber.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Backfill-Schreiber** — `run-store` (dort setzt der atomare Schreiber
  `origin = 'backfill'`); dieser Slice ist eigenständig lieferbar und belegt die
  Boundary von [`LH-FA-DAT-006`](../../../../spec/lastenheft.md) mit alten Zeilen (`NULL` ≙ `wal`).
- **Die Live-Wege** (gRPC, SSE, NATS-Vollinhalt) — sie tragen ausschließlich
  WAL-Changes; ein Feld wäre dort konstant `wal`. Die Schemata
  ([`SPEC-020`](../../../../spec/pflichtenheft.md), [`SPEC-021`](../../../../spec/pflichtenheft.md), [`SPEC-024`](../../../../spec/pflichtenheft.md)) und die Proto-Artefakte bleiben,
  `make generated-sync` ist nicht berührt.
- **Filter `origin` an `GET /changes`** — eine Folge-Entscheidung (Filter-Grammatik
  nach [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)).
- **Die SDK-Lesemodelle** — `sdk-origin`; sie ignorieren das neue Feld bis
  dahin (erwartet, siehe §6).
- **Ein CHECK auf `origin`** — d-migrate konvergiert eine CHECK-Änderung an einer
  bestehenden Tabelle nicht (Kopfkommentar in `tools/schema/schema.yaml`); die
  geschlossene Menge trägt die Domäne, wie bei `request_kind`.

## 2. Definition of Done

- [ ] Domäne: `model.ChangeOrigin` ist eine geschlossene Menge
      (`wal` | `backfill`), der Konstruktor-Default ist `wal`, ein anderer
      Wert wird abgelehnt; `Change.Origin` trägt das Feld, der Doc-Kommentar von
      `model.Change` und der Operationsmenge in
      `internal/domain/model/change.go` nennt es. *Zu belegen durch:* Unit-Tests
      in `change_test.go` (`make test`).
- [ ] Store und Schema: `cdc.change.origin` und die View-Spalte existieren,
      `InsertChange`/`SelectChanges` tragen sie in **expliziten** Spaltenlisten,
      der Mapper (Zeile ↔ Change) trägt sie; eine bestehende Zeile ohne Wert
      (`NULL`) liest über View **und** `ReadChanges` als `wal`
      ([`LH-FA-DAT-006`](../../../../spec/lastenheft.md) Boundary). *Zu belegen durch:* `make test-store` (ein
      Test setzt eine Zeile per SQL mit `NULL`) und `make schema-rollout` zweimal
      hintereinander gegen dieselbe Ziel-Datenbank mit Exit 0 (die additive
      Spalte und die geänderte View konvergieren); der committete Report
      `tools/schema/plan.yaml` und das Rollback-Artefakt `tools/schema/down.sql`
      sind neu erzeugt und mitcommittet (Muster der bisherigen Schema-Slices).
      Der Vertragstest `TestE2EChangesViewMatchesReadChanges` hält View und
      `ReadChanges` auf derselben Spaltenmenge (`BEO-PGC/lese-doppelquelle`).
- [ ] `GET /changes` trägt `origin` je Change; ein fehlender Wert steht als
      `wal`, die Feldreihenfolge bleibt, `origin` steht zuletzt. *Zu belegen
      durch:* Handler-Test in `readchanges_test.go` (`make test`); das
      Benutzerhandbuch (§4 „Änderungen lesen" und der `GET /changes`-Teil der
      HTTP-Beschreibung) nennt das Feld — zwei benannte Stellen: die
      Spaltenliste des SQL-Beispiels unter „Änderungen lesen" und die Feldliste
      der `GET /changes`-Antwort unter „Changes lesen" —, der `Version:`-Kopf
      ist hochgezogen, die Änderungshistorie trägt eine Zeile
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
      `BEO-PGC/handbuch-versionshistorie-uebersprungen`).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: siehe dritter Liefer-Punkt (Handbuch, Änderungshistorie); [`SPEC-002`](../../../../spec/pflichtenheft.md)/[`SPEC-022`](../../../../spec/pflichtenheft.md) sind bereits gezogen.
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
| `internal/domain/model/change.go` (+ `change_test.go`) | update | `ChangeOrigin`, `Change.Origin`, Konstruktor-Default; Kommentar-Nachzug. |
| `internal/domain/errors/` | update (falls nötig) | Fehler für einen unzulässigen Wert. |
| `tools/schema/schema.yaml` | update | Spalte `origin` an `change`; View `changes` mit `COALESCE(...)` als letzter Spalte samt `columns:`-Signatur. |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | `InsertChange`, `SelectChanges`: explizite Spaltenlisten. |
| `internal/adapters/driven/postgresstorage/mapper/mapper.go` (+ `mapper_test.go`) | update | Zeile ↔ Change trägt `origin`. |
| `internal/adapters/driven/postgresstorage/store.go` (+ `store_test.go`, `sqlviews_test.go`) | update | Schreiben/Lesen; Testfall `NULL` ≙ `wal`. |
| `internal/adapters/driven/postgresstorage/schema.sql` | prüfen | eingebettete DDL (`ApplySchema`); ob sie noch ein Träger der Spalte ist oder nur Test-Hilfe, klärt der Suchlauf. |
| `internal/adapters/driving/http/readchanges.go` (+ `readchanges_test.go`) | update | Antwortfeld `origin`. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | regeneriert | Ergebnis von `make schema-rollout`, committet. |
| `docs/user/benutzerhandbuch.md` | update | §4 „Änderungen lesen" (Spaltenliste des SQL-Beispiels), HTTP-Beschreibung von `GET /changes` (Feldliste der Antwort), `Version:`-Kopf und Änderungshistorie. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „Feld- und Spaltenmenge von `cdc.change`, `cdc.changes` und `GET /changes`; die Live-Wege bleiben bei zehn Feldern"; beide Stände gemessen: Parent und Diff):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Spaltenlisten der View in Doku | `grep -rn 'committed_at' docs/user spec harness` | *(Implementer trägt ein)* | Handbuch-Stellen ziehen; Spec ist gezogen. Übergabe aus dem Spec-Nachzug, am Stand `89053d3b` nachgemessen (`grep -n 'committed_at' docs/user/benutzerhandbuch.md`): die Spaltenliste des SQL-Beispiels unter „Änderungen lesen" (Z. 377) und die Feldliste der `GET /changes`-Antwort unter „Changes lesen" (Z. 639) tragen `origin` nicht; Zeilennummern sind der Stand dieser Messung, am Start neu messen |
| Anzahl-Formulierungen („zehn Felder", „zwölf Felder") an `GET /changes` | `grep -rn 'Felder' docs/user spec` | *(Implementer trägt ein)* | nur Stellen, die `GET /changes` betreffen, ziehen; die Live-Wege bleiben bei zehn |
| Harness/Tests mit `SELECT *` gegen `cdc.changes` | `grep -rn 'SELECT \*' tools test internal` | *(Implementer trägt ein)* | prüfen, ob die zusätzliche Spalte die Aussage ändert |
| eingebettete DDL, Report, Rollback | `git ls-files tools/schema internal/adapters/driven/postgresstorage/schema.sql` und Lesen | *(Implementer trägt ein)* | `plan.yaml`/`down.sql` regenerieren; `schema.sql` entscheiden und im Bericht nennen |
| Wegwerf-Client `tools/harness/httpclient` (liest `GET /changes`) | Lesen der Dekodierung | *(Implementer trägt ein)* | nur anpassen, wenn er strikt dekodiert |
| E2E-Abdeckungs-Tabelle (`Datei:Zeile`-Anker) | `git diff --stat` auf `test/integration/**` und `tools/harness/run-integration-tests.sh` | *(Implementer trägt ein)* | berührt der Zug den Runner oder das Testpaket, regeneriert `make test-integration` `docs/user/e2e-abdeckung.md` und der Zug committet sie |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `spec-nachzug` in `done/` liegt und
kein anderer Slice in `in-progress/` liegt (WIP-Limit 1). Technisch unabhängig
von `row-image-gemeinsam`; die Reihenfolge folgt der Tabelle der Welle.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Schema, Store
  und HTTP-Handler nicht in einem Review tragen — der abtrennbare Teil ist der
  HTTP-Handler mit Handbuch.
- `in-progress` → `open` (blockiert — Carveout?): falls d-migrate die additive
  Spalte oder die geänderte View nicht konvergiert (Exit 5) — dann der
  Nacharbeit-Weg nach [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) §Re-Evaluierungs-Trigger, der
  `knownForeignObjects` berührt; das wäre eine Architect-Frage.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test`, `make test-store` und
`make schema-rollout` (zweimal) real grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **d-migrate konvergiert die additive Spalte oder die View-Änderung nicht.**
  Die Anlage einer Spalte an einer bestehenden Tabelle konvergiert nach dem
  Kopfkommentar von `tools/schema/nacharbeit-administration.sql` (Spalte
  `column_name`, dort als real gemessen, Exit 0); der Folgelauf gegen eine
  bereits migrierte Datenbank mit `ReplaceView` steht in `harness/README.md`
  §Sensors, Zeile `make schema-rollout`. Beides ist für **diese** Spalte und
  **diese** View ungemessen. *Erwartet, zu belegen durch:* zweiter
  `make schema-rollout`. **Ausgang:** *(bei Closure)*
- **View und Lesepfad driften auseinander** (`BEO-PGC/lese-doppelquelle`,
  verkörpert): die View trägt `origin`, `SelectChanges` nicht (oder umgekehrt).
  *Erwartet, zu belegen durch:* `TestE2EChangesViewMatchesReadChanges` und der
  Store-Test. **Ausgang:** *(bei Closure)*
- **Ein Client dekodiert strikt und bricht am neuen Feld.** Die SDK-Decoder
  ignorieren unbekannte Felder (gelesen für C# `System.Text.Json`, Kotlin Gson,
  Python `json` mit expliziten Feldern); ein realer Beleg fehlt bis
  `sdk-origin`. *Erwartet, zu belegen durch:* `make test-sdk-*-integration`-Läufe
  bzw. der Slice `sdk-origin`. **Ausgang:** *(bei Closure)*
- **Die Zeilenzahl-Anker der E2E-Abdeckung wandern**, sobald der Zug das
  Testpaket berührt. *Erwartet, zu belegen durch:* der Suchlauf-Eintrag zur
  Abdeckungs-Tabelle. **Ausgang:** *(bei Closure)*

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
`*`/`PGC` (Greenfield); Domäne, Store-Adapter, Schema und HTTP-Handler sind
keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/lese-doppelquelle` (verkörpert, 3×, einschlägig — Risiko §6),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3),
`BEO-PGC/d-migrate-nacharbeit` (verkörpert, 6×, einschlägig — Rückführung §4),
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` und
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert, je 3×, im DoD),
`BEO-PGC/limit-fortsetzung-innerhalb-einer-position` (offen, 0×, Bezug der
Welle — dieser Slice ändert das Lesen nicht),
`BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×, neue Store-Tests
begrenzen ihre Bereinigung auf die eigenen Zeilen).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
