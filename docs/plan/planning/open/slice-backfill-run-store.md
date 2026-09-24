# Slice backfill-run-store: Run-Store — Tabelle `cdc.backfill_run`, Grants, Postgres-Adapter für Run-Zustand und den einen atomaren Schreiber

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Negative: Neubeginn ohne Verlust — ein
abgebrochener Run hinterlässt nichts), [`LH-FA-CAP-004`](../../../../spec/lastenheft.md) (Ordnung),
[`LH-FA-REA-004`](../../../../spec/lastenheft.md) (deterministische Sortierung), [`LH-FA-DAT-006`](../../../../spec/lastenheft.md)
(Metadaten-Erweiterbarkeit), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4/6/7 (Atomarität,
Run-Zustand, Ordnung, Retention), [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) (Schemamigrationen mit
d-migrate), [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md) (rollenspezifische DSN-Verdrahtung), [`ADR-0017`](../../adr/0017-generische-change-tabelle.md)
(generische Change-Tabelle), [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1/3 (Rollenschnitt, Grants, Annahme in einer
Transaktion, Warn-Spalten, `estimated_rows` = `NULL` als „unbekannt").

**Berührte Spec-Stellen:** [`SPEC-029`](../../../../spec/pflichtenheft.md) (Feldform `cdc.backfill_run`, durch
`spec-nachzug`), [`SPEC-001`](../../../../spec/pflichtenheft.md), [`SPEC-002`](../../../../spec/pflichtenheft.md) — gelesen, nicht geändert;
[`ARC-006`](../../../../spec/architecture.md) (Driven Adapter).

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Postgres-Seite des Runs: die Tabelle `cdc.backfill_run`
(`run_id` = die `administration_request_id` des Antrags, `source_id`,
`schema_name`, `table_name`, `status`, `requested_at`, `started_at`,
`finished_at`, `snapshot_position`, `rows_copied`, `estimated_rows`,
`error_message`; `estimated_rows` ist **nullable** — `NULL` heißt „unbekannt",
nie `0` —; dazu die **zwei Warn-Spalten** `warn_estimated_size` und `warn_duration`
für die beiden Warnungen (Typ `boolean NOT NULL DEFAULT false` wie in [`SPEC-029`](../../../../spec/pflichtenheft.md)
durch `spec-nachzug` festgelegt; `false`, solange keine Auswertung sie setzt); die geschlossene `status`-Menge als CHECK **bei
Erstanlage** der Tabelle), die Grants, und die Adapter für die drei Ports aus
`run-usecase`:

- der **Annahme-Adapter** (`Admit`, Paket `postgresstorage`, Pool der
  Administrations-Goroutine aus `CDC_ADMIN_DSN`, Naht `sqlexec.DB` mit `Begin`):
  **eine** Transaktion — Prüfung „kein aktiver Run (`queued`/`running`) derselben
  Tabelle" (Lesen vor Einfügen), `INSERT` der Run-Zeile `queued` (`run_id` = die
  `administration_request_id`, geschätzte Zeilenzahl, ggf. Warnung 1),
  `UPDATE … applied` des Antrags, der genau **eine** `pending`-Zeile treffen muss,
  Commit; jede Abweichung ist Rollback und Fehler, ein aktiver Run ein
  Sentinel-Fehler;
- der **Run-Zustands-Adapter** (Pool aus `CDC_CAPTURE_DSN`, **keine**
  Anlage-Operation): Übergänge, Fortschritt (`rows_copied` je Block
  **außerhalb** der Daten-Transaktion, damit er für Leser sichtbar ist; trägt
  auch die Warnung 2), der Abgleich `running` → `interrupted`, die `queued`-Zeilen
  der eigenen Quelle in Antragsreihenfolge (`ORDER BY requested_at, run_id`);
- der **atomare Schreiber**: **eine** Store-Transaktion über alle Blöcke — je
  Block eine `cdc.transaction`-Zeile (Kennung `0bf-<run-id>-<Block>`, Position
  `X`, `committed_at`) und `cdc.change`-Zeilen (`origin = 'backfill'`,
  `operation = 'INSERT'`, `new_data` = Bild) — und **ein** Commit am Ende
  zusammen mit der Run-Zeile `completed`; Rollback hinterlässt keine Zeile.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die SQL-Funktion, die Antragsart, der Worker, die View `cdc.backfill_status`
  und der Idempotenz-Guard** — `sql-administration`; diese Tabelle ist eine
  d-migrate-fähige Tabelle des neutralen Modells und braucht **keinen**
  Guard-Eintrag (nur Objekte außerhalb des Modells stehen in
  `knownForeignObjects`). *Zu belegen durch:* zweiter `make schema-rollout`.
- **Das Lesen des Bestands** — `cdc.changes` und `GET /changes` liefern ihn
  bereits über die Spalte `origin` aus `change-origin`.
- **Block-Transaktionen und Checkpoint** — Welle §6; die Ein-Transaktions-Form ist
  die Festlegung von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4.
- **Die Retention-Regel für Backfill-Changes** — sie gilt unverändert
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 7, kein Sonderpfad); der `e2e`-Slice trägt keinen
  Retention-Beleg, weil die Regel nichts Neues festlegt.

## 2. Definition of Done

- [ ] Schema und Grants: `cdc.backfill_run` liegt in `tools/schema/schema.yaml`
      (CHECK auf `status` bei Erstanlage, `estimated_rows` nullable, die
      zwei Warn-Spalten nach [`SPEC-029`](../../../../spec/pflichtenheft.md)), die Grants stehen in
      `tools/schema/nacharbeit-roles.sql` — `cdc_admin` mit `SELECT`, `INSERT`,
      `cdc_capture` mit `SELECT`, `UPDATE`, niemand mit `DELETE`, `cdc_reader` ohne
      Recht auf die Basistabelle ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1, Punkt 3). *Zu belegen
      durch:* `make schema-rollout` zweimal
      hintereinander mit Exit 0, `tools/schema/plan.yaml` und
      `tools/schema/down.sql` neu erzeugt und mitcommittet, der Rollen-Test
      `internal/bootstrap/roles_rollout_file_internal_test.go` (netzlos, Teil der
      Gates) angepasst, und ein Rollen-Test im Store-Tier (`roles_test.go`-Muster,
      je Rolle ein Login): `cdc_capture` — `INSERT` scheitert mit SQLSTATE `42501`,
      `UPDATE` gelingt; `cdc_admin` — `INSERT` gelingt, `UPDATE` und `DELETE`
      scheitern; `cdc_reader` — `SELECT` auf die Basistabelle scheitert.
- [ ] Annahme-Adapter: ein Store-Test gegen reale PostgreSQL zeigt, dass eine
      Annahme Run-Zeile `queued` und Antragsvermerk `applied` **zugleich**
      hinterlässt, dass ein Antragsvermerk ohne `pending`-Zeile (Rollback-Fall)
      **weder** Run-Zeile **noch** Vermerk hinterlässt, dass ein zweiter Antrag bei
      aktivem Run ohne zweite Zeile endet, und dass `estimated_rows` als `NULL` und
      nicht als `0` gespeichert und gelesen wird. *Zu belegen durch:*
      `make test-store`.
- [ ] Atomarität: ein Store-Test gegen reale PostgreSQL zeigt, dass ein zweiter
      Leser vor dem Commit **keine** Zeile des Runs sieht und danach alle
      Blöcke zugleich, dass ein Rollback **keine** `cdc.transaction`- und keine
      `cdc.change`-Zeile hinterlässt, und dass jeder Change `origin = 'backfill'`,
      `operation = 'INSERT'` und `old_data IS NULL` trägt. *Zu belegen durch:*
      `make test-store`.
- [ ] Ordnung und Zustand: mit einem WAL-Commit **auf derselben Position** `X`
      liest `cdc.changes` die Backfill-Blöcke **vor** dem WAL-Commit
      (`(commit_position, transaction_id, sequence)`); die Übergänge, der
      Fortschritt außerhalb der Daten-Transaktion und der Abgleich `running` →
      `interrupted` verhalten sich wie im Use-Case-Vertrag. *Zu belegen durch:*
      `make test-store` in der Test-Datenbank; der Bericht nennt die
      Kollation der Datenbank (`datcollate`), weil die Ordnung `0bf-…` vor
      Ziffern-Kennungen eine Textsortierung ist.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: die Sensor-Zeile `make test-store` in `harness/README.md` und die Schema-Kommentare, soweit sie die Tabellenliste oder die Rollenverteilung aufzählen; das Benutzerhandbuch bleibt bis `sql-administration` unberührt.
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
| `tools/schema/schema.yaml` | update | Tabelle `backfill_run`; CHECK auf `status` bei Erstanlage (spätere Erweiterung über Nacharbeit-SQL, nach dem Muster von `request_kind`); `estimated_rows` nullable; die zwei Warn-Spalten. |
| `tools/schema/nacharbeit-roles.sql` | update | Grants: `cdc_admin` `SELECT`, `INSERT`; `cdc_capture` `SELECT`, `UPDATE`. |
| `internal/bootstrap/roles_rollout_file_internal_test.go` | update | hält den Rollenschnitt gegen die Rollout-Datei; die neue Tabelle darf für `cdc_reader` **nicht** als Basistabellen-Grant erscheinen. |
| `internal/adapters/driven/postgresstorage/` (Arbeitsname `backfilladmission.go`, `backfillrun.go`, `backfillwriter.go`, + Tests) | neu | die drei Adapter über die schmale Ausführungs-Naht des Pakets (`sqlexec`); explizite Spaltenlisten; der Annahme-Adapter nutzt `Begin`. |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | die neuen Anweisungen. |
| `internal/adapters/driven/postgresstorage/schema.sql` | prüfen | wie in `change-origin`: ob die eingebettete DDL eine Träger-Rolle hat. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | regeneriert | Ergebnis von `make schema-rollout`, committet. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Tabellenliste des `cdc`-Schemas und die Rollenverteilung der Grants"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Tabellenlisten in Doku und Kommentaren | `grep -rn 'cdc.administration_request' docs spec harness tools --include=*.md --include=*.sql --include=*.yaml` (Aufzählungen von Tabellen) | *(Implementer trägt ein)* | Aufzählungen nachziehen; [`SPEC-001`](../../../../spec/pflichtenheft.md) trägt `spec-nachzug` |
| Rollenverteilung in Kommentaren | Lesen der Kopfkommentare von `nacharbeit-roles.sql` | *(Implementer trägt ein)* | nachziehen, wo die Aufzählung `cdc_capture` beschreibt |
| Rollen-Test der Rollout-Datei | Lesen von `internal/bootstrap/roles_rollout_file_internal_test.go` (Regeln zu Views/Basistabellen) | *(Implementer trägt ein)* | anpassen |
| Test-Bereinigung anderer Store-Tests | Lesen der `Cleanup`-Blöcke in `postgresstorage/*_test.go` | *(Implementer trägt ein)* | neue Tests bereinigen nur ihre eigenen Zeilen (`BEO-PGC/test-isolation-geteilter-zustand`) |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `run-usecase` und `change-origin` in
`done/` liegen und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Schema/Grants
  und die beiden Adapter nicht in einem Review tragen — der abtrennbare Teil
  ist Schema samt Grants (dann als eigener Schritt vor den Adaptern).
- `in-progress` → `open` (blockiert — Carveout?): falls d-migrate die Tabelle
  oder den CHECK bei Erstanlage nicht konvergiert (Exit 5) — dann Nacharbeit-SQL
  nach [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) §Re-Evaluierungs-Trigger und ein Eintrag in
  `knownForeignObjects` im selben Commit; Architect-Frage.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test-store` und
`make schema-rollout` (zweimal) real grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Textsortierung trägt die Ordnung `0bf-…` vor Ziffern-Kennungen**
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 6 leitet sie her: WAL-Kennungen beginnen nie mit
  `0`). Sie hängt an der Kollation der Datenbank (`ORDER BY t.commit_position,
  c.transaction_id, c.sequence` in der View). *Erwartet, zu belegen durch:* der
  Ordnungs-Test in der Test-Datenbank samt genannter Kollation; weicht eine
  gängige Kollation ab, wird der Betreiber-Hinweis Teil des Handbuchs.
  **Ausgang:** *(bei Closure)*
- **Der CHECK bei Erstanlage konvergiert** (die Zusage von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage
  4): die Tabelle `change` trägt CHECKs bei Erstanlage bereits (`chk_change_operation`
  in `tools/schema/schema.yaml`); für **diese** Tabelle ungemessen. *Erwartet, zu
  belegen durch:* Erst- und Folgelauf des Rollouts. **Ausgang:** *(bei Closure)*
- **Der Fortschritt außerhalb der Daten-Transaktion** braucht eine zweite
  Verbindung und kollidiert nicht mit dem Commit derselben Run-Zeile. *Erwartet,
  zu belegen durch:* der Zustands-Test mit gleichzeitigem Fortschritts-Update
  und Commit. **Ausgang:** *(bei Closure)*
- **Die Rollenlage** ist mit [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 gesetzt (`cdc_admin` legt an, `cdc_capture`
  ändert, `cdc_reader` liest nur die View); ein `INSERT` für `cdc_capture`, ein
  `UPDATE` oder `DELETE` für `cdc_admin` oder ein `SELECT`-Grant für `cdc_reader`
  auf die Basistabelle wäre ein Rollen-Verstoß ([`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…[`LH-QA-SEC-003`](../../../../spec/lastenheft.md)).
  *Erwartet, zu belegen durch:* der Rollen-Test im Store-Tier und der netzlose
  Rollen-Test der Rollout-Datei. **Ausgang:** *(bei Closure)*
- **Die Annahme-Transaktion ist nicht atomar oder greift zu weit.** Der
  Antragsvermerk muss genau eine `pending`-Zeile treffen, sonst Rollback; der
  Verzicht auf eine Sperre zwischen mehreren Annehmenden beruht auf der Annahme,
  dass `processAdministrationRequests` Anträge sequenziell in **einer** Goroutine
  verarbeitet und je Quelle eine Instanz Anträge annimmt ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1,
  Punkt 2). *Erwartet, zu belegen durch:* der Rollback-Test des zweiten
  Liefer-Punkts. **Ausgang:** *(bei Closure)*
- **Geteilter Zustand in den Store-Tests** (`BEO-PGC/test-isolation-geteilter-zustand`,
  offen, 1×): unskopierte Bereinigung in den bestehenden Tests würde die Zeilen
  eines parallelen Tests treffen. *Erwartet, zu belegen durch:* die neuen Tests
  bereinigen nur ihre Run-Kennung. **Ausgang:** *(bei Closure)*
- **DB-Adapter-Coverage** ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 3): neuer, DB-gestützt geprüfter
  Code in `postgresstorage` bewegt Zähler und Nenner. *Erwartet, zu belegen
  durch:* der Bericht nennt die Zahl mit ihrem Lauf ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A).
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
`*`/`PGC` (Greenfield); Schema, Grants und der Store-Adapter sind keine
eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/d-migrate-nacharbeit` (verkörpert, 6×, einschlägig — Rückführung §4),
`BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×, Tabelle bleibt im
neutralen Modell — kein Guard-Eintrag nötig, zu belegen durch den zweiten Lauf),
`BEO-PGC/rollen-verdrahtung` (eingetreten) und
`BEO-PGC/rollen-test-abdeckungsluecken` (offen, 2×, einschlägig — Rollen-Test
im DoD; ein weiterer Beleg erreicht 3×),
`BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×, Risiko §6),
`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (verkörpert — der
Run-Zustand ist eine Tabelle, kein Prozessspeicher),
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (offen, 2×, gesichtet
— hier wächst der DB-Gegenstand mit DB-gestütztem Code),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
