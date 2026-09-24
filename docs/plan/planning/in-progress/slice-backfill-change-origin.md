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

**Verantwortlich:** Implementer-Agent, 2026-09-24.

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

- [x] Domäne: `model.ChangeOrigin` ist eine geschlossene Menge
      (`wal` | `backfill`), der Konstruktor-Default ist `wal`, ein anderer
      Wert wird abgelehnt; `Change.Origin` trägt das Feld, der Doc-Kommentar von
      `model.Change` und der Operationsmenge in
      `internal/domain/model/change.go` nennt es. *Zu belegen durch:* Unit-Tests
      in `change_test.go` (`make test`).
- [x] Store und Schema: `cdc.change.origin` und die View-Spalte existieren,
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
      *Stand des Implementers, zwei Grenzen benannt:* (1) `make schema-rollout`
      zweimal Exit 0 ist gegen eine **frische** Ziel-Datenbank gemessen; gegen
      ein mit dem Parent-Schema ausgerolltes Ziel blockiert die View-Änderung
      (§3, „Messung zu Risiko §6"). (2) `TestE2EChangesViewMatchesReadChanges`
      bleibt unverändert; die `origin`-Parität von View und `ReadChanges` trägt
      der Store-Tier-Test `TestChangesViewCarriesOriginLikeReadChanges`
      (§3, Zeile zu `test/integration`).
- [x] `GET /changes` trägt `origin` je Change; ein fehlender Wert steht als
      `wal`, die Feldreihenfolge bleibt, `origin` steht zuletzt. *Zu belegen
      durch:* Handler-Test in `readchanges_test.go` (`make test`); das
      Benutzerhandbuch (§4 „Änderungen lesen" und der `GET /changes`-Teil der
      HTTP-Beschreibung) nennt das Feld — zwei benannte Stellen: die
      Spaltenliste des SQL-Beispiels unter „Änderungen lesen" und die Feldliste
      der `GET /changes`-Antwort unter „Changes lesen" —, der `Version:`-Kopf
      ist hochgezogen, die Änderungshistorie trägt eine Zeile
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
      `BEO-PGC/handbuch-versionshistorie-uebersprungen`).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: siehe dritter Liefer-Punkt (Handbuch, Änderungshistorie); [`SPEC-002`](../../../../spec/pflichtenheft.md)/[`SPEC-022`](../../../../spec/pflichtenheft.md) sind bereits gezogen.
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
| `internal/domain/errors/errors.go` | update | `ErrInvalidChangeOrigin` für einen unzulässigen Wert (Gemessen: nötig — die Konstruktoren lehnen über einen Sentinel ab). |
| `tools/schema/schema.yaml` | update | Spalte `origin` an `change`; View `changes` mit `COALESCE(...)` als letzter Spalte samt `columns:`-Signatur. |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | `InsertChange`, `SelectChanges`: explizite Spaltenlisten. |
| `internal/adapters/driven/postgresstorage/mapper/mapper.go` (+ `mapper_test.go`) | update | Zeile ↔ Change trägt `origin`. |
| `internal/adapters/driven/postgresstorage/store.go` (+ `store_test.go`, `sqlviews_test.go`) | update | Schreiben/Lesen; Testfall `NULL` ≙ `wal`. |
| `internal/adapters/driven/postgresstorage/schema.sql` | update | Entscheidung nach dem Suchlauf: die eingebettete DDL ist Test-Schema-Loader (`ApplySchema` in `store_test.go` und `tableactivation_test.go`, `DROP SCHEMA cdc CASCADE` + `ApplySchema` je Test) und damit zweiter Schema-Träger der Spalte — ohne `origin text` scheitert `InsertChange` in diesen Tests. Die Spalte steht dort mit derselben Form wie in `schema.yaml` (nullable, ohne Default, ohne CHECK). |
| `internal/adapters/driven/postgresstorage/sqlexec/translate.go` (+ `translate_test.go`) | update (nicht im ursprünglichen Plan) | die Scan-Schleife von `ReadChanges` liest die Projektion von `SelectChanges` Spalte für Spalte; `origin` ist die 14. Spalte und kommt hier in `mapper.ChangeRow`. Test: Zeilen-Fake mit Herkunft (`wal`, `backfill`, leer, unbekannt). |
| `internal/adapters/driven/natsstream/publisher.go` (+ `publisher_test.go`), `internal/adapters/driving/http/sse.go` | update (Kommentar, nicht im ursprünglichen Plan) | drei Kommentare behaupten „dieselben (zehn) Felder wie `model.Change`" — mit `Change.Origin` stimmt das nicht mehr; sie nennen jetzt „ohne das Feld `origin`" (§3.13, Suchlauf-Zeile „Kommentare in den Live-Wegen"). |
| Tests: `mapper_test.go`, `store_test.go`, `sqlviews_test.go` | update | Herkunft-Rundlauf im Mapper; Store-Test (`wal`/`backfill`/`NULL`, unbekannter Wert wird abgelehnt); View-Test (`origin` letzte Spalte, `NULL` ≙ `wal`, Parität View ↔ `ReadChanges`). |
| `test/integration/integration_test.go` (`TestE2EChangesViewMatchesReadChanges`) | **nicht geändert** | der Vertragstest bleibt auf den bisherigen Spalten; die `origin`-Parität von View und `ReadChanges` trägt der Store-Tier-Test `TestChangesViewCarriesOriginLikeReadChanges` (`make test-store`). Grund: das Testpaket zu berühren verschöbe die `Datei:Zeile`-Anker in `docs/user/e2e-abdeckung.md` (Risiko §6, vierter Punkt). |
| `internal/adapters/driving/http/readchanges.go` (+ `readchanges_test.go`) | update | Antwortfeld `origin`. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | regeneriert | Ergebnis von `make schema-rollout`, committet. |
| `docs/user/benutzerhandbuch.md` | update | §4 „Änderungen lesen" (Spaltenliste des SQL-Beispiels), HTTP-Beschreibung von `GET /changes` (Feldliste der Antwort), `Version:`-Kopf und Änderungshistorie. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „Feld- und Spaltenmenge von `cdc.change`, `cdc.changes` und `GET /changes`; die Live-Wege bleiben bei zehn Feldern"; beide Stände gemessen: Parent und Diff):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Spaltenlisten der View in Doku | Parent (Commit `f4ba82ab`, Stand vor diesem Slice): `git grep -n 'committed_at' f4ba82ab -- docs/user spec harness`; Diff-Stand (Commit `e95937cf`, Arbeitsbaum sauber): dasselbe mit `e95937cf` | Beide Stände **6** Treffer (gemessen), davon 3 im Handbuch: Parent `docs/user/benutzerhandbuch.md` Z. 377 (Spaltenliste des SQL-Beispiels unter „Änderungen lesen"), Z. 419 (`SELECT change_id, committed_at` — explizite Zwei-Spalten-Abfrage unter „Aufbewahrung", von `origin` nicht berührt), Z. 639 (Feldliste der `GET /changes`-Antwort unter „Changes lesen"); 3 in `spec/pflichtenheft.md` (Z. 236 synthetische Transaktionen, Z. 603 `SPEC-022`-Antwortzelle — trägt `origin` bereits —, Z. 819 Änderungsverlauf), keiner in `harness/`. Ergänzend `git grep -n -w 'origin' <Stand> -- docs/user`: Parent 1 (`docs/user/releasing.md`, das Git-Remote), Diff-Stand 6 (das Remote plus fünf Zeilen im Handbuch); `-- spec`: je Stand 9 (Spec ist gezogen) | Handbuch-Z. 377 und Z. 639 gezogen (die zwei benannten Stellen), Z. 419 unverändert; Version 1.44 samt Historienzeile; Spec unverändert |
| Anzahl-Formulierungen („zehn Felder", „zwölf Felder") an `GET /changes` | `git grep -n 'Felder' <Stand> -- docs/user spec` und `git grep -n -E 'elf Felder\|zwölf Felder\|dreizehn Felder' <Stand> -- docs/user spec` | `Felder` je Stand **33** Treffer (gemessen), davon `zehn Felder` je Stand **19** (gRPC-, SSE- und NATS-Vollinhalts-Abschnitte, deren SDK-Absätze und `SPEC-020`/`SPEC-021`/`SPEC-024` — Live-Wege, bleiben bei zehn); `elf`/`zwölf`/`dreizehn Felder` je Stand **0** — keine Anzahl-Formulierung für `GET /changes` in `docs/user` oder `spec`. Außerhalb des Suchraums: `sdks/python/pgchangefeed/src/pgchangefeed/models.py` Z. 264 („nicht die zwölf des HTTP-Lesezugriffs (`SPEC-022`)") — der HTTP-Lesezugriff trägt mit `origin` dreizehn Felder | keine Änderung im Suchraum nötig; der Python-Kommentar gehört zu `slice-backfill-sdk-origin` (fremde Datei, gemeldet statt still geändert) |
| Harness/Tests mit `SELECT *` gegen `cdc.changes` | `git grep -n 'SELECT \*' <Stand> -- tools test internal` | Je Stand **3** Treffer (gemessen): `test/integration/integration_test.go` Z. 1002 und 1052 (Kommentare zur expliziten Spaltenliste, kein SQL) und `tools/bench-batch-vs-single.sh` Z. 42 (`SELECT count(*) FROM (SELECT * FROM cdc.changes … LIMIT $M) t` — zählt Zeilen, liest die zusätzliche Spalte mit, ohne dass eine Aussage davon abhängt) | keine Änderung: keine der drei Stellen setzt eine Spaltenmenge voraus, die `origin` verletzt |
| Kommentare in den Live-Wegen („dieselben Felder wie `model.Change`") | `git grep -n -E 'dieselben (zehn )?Felder wie\|Feldern wie der Domain-Typ' <Stand> -- '*.go' '*.proto' '*.py'` | Je Stand **8** Treffer (gemessen), die `model.Change` als Feldvorbild nennen: `natsstream/publisher.go` Z. 135, `natsstream/publisher_test.go` Z. 261, `driving/http/sse.go` Z. 29 — mit `Change.Origin` nicht mehr wahr —; `gen/cdc/stream/v1/changestream.pb.go` Z. 31 und `proto/cdc/stream/v1/changestream.proto` Z. 13 (dieselbe Aussage über die gRPC-Nachricht); `driving/http/readchanges.go` Z. 46 (`GET /changes`: trägt `origin`, bleibt wahr); `natsstream/publisher.go` Z. 165 (verweist auf das SSE-Event, bleibt wahr); Python-`models.py` Z. 263 | die drei Go-Kommentare ziehen „ohne `Origin`" nach; die `.proto`-Quelle und die generierte Datei **nicht** (die Zeile 13/31 zu ändern verlangt `make proto-generate`, der Plan legt die Proto-Artefakte und `make generated-sync` als unberührt fest) — gemeldet: die Aussage „mit denselben Feldern wie der Domain-Typ" ist an der gRPC-Nachricht ungenau geworden; der Python-Kommentar gehört zu `slice-backfill-sdk-origin` |
| eingebettete DDL, Report, Rollback | `git ls-tree -r --name-only <Stand> -- tools/schema internal/adapters/driven/postgresstorage/schema.sql` | Je Stand **14** Dateien (gemessen, gleiche Menge); `tools/schema/plan.yaml` und `tools/schema/down.sql` ändern sich im Diff (Ergebnis eines frischen `make schema-rollout` gegen eine leere Datenbank), `schema.sql` und `schema.yaml` ebenso. `schema.sql` ist Träger, nicht nur Test-Hilfe: `store_test.go` und `tableactivation_test.go` bauen das Schema je Test über `ApplySchema` auf; ohne die Spalte dort scheitert `InsertChange` (rot gesehen, Mutation Z12) | `plan.yaml`/`down.sql` regeneriert und committet; `schema.sql` gezogen (Entscheidung, siehe Zeile in der Plan-Tabelle oben) |
| Wegwerf-Client `tools/harness/httpclient` (liest `GET /changes`) | Lesen der Dekodierung (`tools/harness/httpclient/main.go` Z. 126–129) und `git grep -n -i -E 'DisallowUnknown\|UnmappedMemberHandling\|FAIL_ON_UNKNOWN\|ignoreUnknownKeys\|extra=.forbid' <Stand> -- '*.go' '*.cs' '*.kt' '*.py'` | Der Client dekodiert über `json.Unmarshal` in eine Struktur ohne `origin` — nicht strikt; strikte Dekoder-Muster je Stand **0** Treffer (gemessen) im ganzen Baum | keine Änderung nötig |
| E2E-Abdeckungs-Tabelle (`Datei:Zeile`-Anker) | `git diff --stat f4ba82ab e95937cf -- test/integration tools/harness/run-integration-tests.sh docs/user/e2e-abdeckung.md` | Leer (gemessen, 0 Zeilen): der Diff berührt weder das Testpaket noch den Runner noch die Tabelle | keine Regeneration nötig; `make test-integration` lief nach dem Diff-Stand real (Exit 0) und meldete „E2E-Abdeckungstabelle unverändert" |

**Messung zu Risiko §6 „d-migrate konvergiert nicht" (Implementer; der Ausgang gehört der Closure).** Alle Läufe gegen einen Wegwerf-PostgreSQL-18-Container (`postgres:18-alpine`, Digest aus `PG_TEST_IMAGE`), Rollout-Werkzeug d-migrate im gepinnten Image (Makefile-Variable `D_MIGRATE_IMAGE`); Ausgabe je Lauf gedruckt, hier die Exit-Codes:

- **Frische Ziel-Datenbank, neues Schema, `make schema-rollout` zweimal hintereinander:** beide Läufe Exit 0 (gemessen). Der zweite Lauf nimmt den `--allow-destructive`-Pfad der Idempotenz-Wache (Meldung in der Ausgabe); die View `cdc.changes` führt danach `origin` als letzte Spalte. Der erste dieser Läufe erzeugt die committeten `plan.yaml`/`down.sql` (14 Operationen, `change` mit der Spalte `origin`, View `changes` mit `COALESCE(c.origin, 'wal') AS origin`).
- **Ziel-Datenbank mit dem Schema des Parent-Commits `f4ba82ab` (zuvor per `make schema-rollout` ausgerollt), danach das neue Schema:** `make schema-rollout` endet mit Exit 2 (der Precheck-Lauf meldet Exit 8). Blocker im Precheck-Report: `MANUAL_ACTION_REQUIRED` für `ReplaceView` der View `changes`, Diagnose `VIEW_SIGNATURE_INCOMPATIBLE` („CREATE OR REPLACE VIEW is only renderable when view columns keep the same count, order, names and visible types"). Die Operation `AddColumn` für `change.origin` steht im selben Plan als regulär renderbar; die Wache (`tools/schema/rolloutguard`) lehnt die Blocker-Klasse `MANUAL_ACTION_REQUIRED` ab und läuft ohne `--allow-destructive`. **Ein bereits mit dem alten Schema ausgerolltes Ziel lässt sich damit über `make schema-rollout` nicht auf den neuen Stand heben.**
- **Derselbe Ausgangsstand, nach einem manuellen `DROP VIEW cdc.changes` vor dem Rollout:** Exit 0, ein zweiter Lauf Exit 0; eine zuvor ohne Spalte gespeicherte Zeile liest über die View als `wal` (`NULL` in `cdc.change.origin`), die Rechte auf die View sind danach wieder gesetzt (`\dp cdc.changes` zeigt `cdc_reader=r`). Das ist eine Messung, kein umgesetzter Weg: einen Vorlauf-Schritt im Makefile-Target oder eine Sicht außerhalb des neutralen Modells zu führen, berührt [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) und `knownForeignObjects` — die in §4 vorab benannte Architect-Frage; dieser Slice setzt sie nicht um.
- **Reichweite der Aussage:** die Test-Tiers (`make test-store`, `make test-replication`, `make test-integration`), `make example-demo-up` und der Compose-Aufbau legen ihr Schema laut `harness/README.md` §Sensors frisch an; der bestehende Rollout-Check in `examples/bootstrap.sh` überspringt den Rollout gegen ein bereits migriertes Ziel — dass ein dort weiterlebendes Alt-Schema `InsertChange` scheitern lässt, ist abgeleitet, nicht gemessen.

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
