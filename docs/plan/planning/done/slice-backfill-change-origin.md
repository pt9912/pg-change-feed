# Slice backfill-change-origin: Feld `origin` von der Domäne bis zu den Lesewegen — `cdc.change`, View `cdc.changes`, `GET /changes`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](welle-backfill-bestand.md).

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

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](welle-backfill-bestand.md). **Datum:** 2026-09-23.

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
      ein mit dem Parent-Schema ausgerolltes Ziel blockierte die View-Änderung
      (§3, „Messung zu Risiko §6") — aufgelöst durch
      [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
      (Vorlauf `DROP VIEW` im Target `schema-rollout`, Fixrunde) und belegt
      durch den „Ausgang der Messung" in §3 (Upgrade von `f4ba82ab` und vom
      Tag `v0.1.2`, je zweimal Exit 0). (2) `TestE2EChangesViewMatchesReadChanges`
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
- [x] Upgrade eines Alt-Schemas über die View-Signaturänderung
      ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md),
      [`LH-QA-OPS-005`](../../../../spec/lastenheft.md)): `tools/schema/rolloutguard`
      erkennt die Klasse und liefert die Views, `make schema-rollout` entfernt sie
      vor `--execute` (ohne `CASCADE`, gemeldet), Unit-Tests in `guard_test.go`
      (`make test`), Guard-Test `tools/harness/run-schema-rollout-guard-test.sh`
      mit sechs Läufen (View-Signatur-Lauf, Alt-Tag-Lauf), Handbuch §4 „Schema
      aktualisieren" (Version 1.46) und `harness/README.md` §Sensors.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
      Report `docs/reviews/review-slice-backfill-change-origin.md` (0 HIGH,
      2 MEDIUM, 3 LOW, 5 INFO); F-1 bis F-6 in Fixrunde 2 behoben (§3, Tabelle
      „Fixrunde 2"), F-7 und F-10 als gemeldete Träger geführt, F-8 im Handbuch
      und im Target-Dokument benannt, F-9 bewertet ohne Aktion. Ein
      Reviewer-Nachprüfungs-Report zur Fixrunde liegt nicht vor; die Behebung
      hat der Verifier selbst am Code nachgemessen (Verifikation §2 Zeile 6,
      §5 Mutationen), kein offenes HIGH/MEDIUM.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: siehe dritter Liefer-Punkt (Handbuch, Änderungshistorie); [`SPEC-002`](../../../../spec/pflichtenheft.md)/[`SPEC-022`](../../../../spec/pflichtenheft.md) sind bereits gezogen.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke). *(§7: neuer Sensor-Umfang und
      benannte ADR-Lücke.)*
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert. *(Fünf weitere `evidence/`-Dateien
      in bestehenden Klassen, vier neue Klassen; siehe §7.)*
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen). *(Vier Ausgänge in §6, je am Ort; siehe §7.)*
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](welle-backfill-bestand.md) (die Roadmap führt sie unter
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

**Fixrunde zu [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)** (Architect-Entscheidung zur Grenze (1) aus der „Messung zu Risiko §6"; kein neuer Slice, Umfang laut `docs/reviews/architect-verdict-schema-rollout-view-signatur.md` §Folge-Arbeit 1–6). Die Zeilen sind **vor** dem Sensor-Lauf der Fixrunde eingetragen:

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/rolloutguard/report.go`, `guard.go`, `main.go` (+ `guard_test.go`) | update | `report` dekodiert `diagnostics[].code`/`operationId` (das Feld `blockers[].diagnosticCodes` ist im realen Report leer, gemessen — die Diagnose hängt über `operationId` an der Operation); `decide` liefert eine Entscheidung mit `allowDestructive` **und** der Liste der Views der Klasse „View-Signatur" (Operation `ReplaceView`/`VIEW`, Blocker-Grund `MANUAL_ACTION_REQUIRED`, Diagnose `VIEW_SIGNATURE_INCOMPATIBLE` zur selben Operation); alles oder nichts; ein View-Name muss ein einfacher Bezeichner sein (`[a-z_][a-z0-9_]*`), sonst gilt die Operation als unbekannt. Ausgabe von `main`: Maschinenzeilen (`allow-destructive`, `drop-view <name>`) auf stdout, Begründung auf stderr, Exit 0 wenn irgendetwas erlaubt ist, 1 sonst, 2 bei Report-Fehlern. *Präzisierung gegenüber dem ADR-Wortlaut:* der Blocker-Grund `MANUAL_ACTION_REQUIRED` gehört zur Klasse (im realen Report gemessen; ein View-Signatur-Diagnose unter einem anderen Grund fällt in „unbekannt" und bricht mit Exit 8 ab). |
| `Makefile` (Target `schema-rollout`, Kommentarblock darüber) | update | Vorlauf-Schritt zwischen Precheck und `--execute`: je gemeldeter View `DROP VIEW cdc.<name>` (ohne `CASCADE`) über `psql` im `PG_TEST_IMAGE` wie die `nacharbeit`-Schritte, je View eine Meldung; kein Vorlauf ohne Anlass. Der Kommentarblock beschreibt die Klasse im Indikativ (`AGENTS.md` §3.7). |
| `tools/harness/run-schema-rollout-guard-test.sh` | update | Form entschieden (die ADR lässt sie offen): beide Belege liegen **im** Guard-Test, kein eigenes Make-Target (das Skript hat heute keines; ein neues Target wäre ein neues Sensor-Ziel ohne Gate-Bindung). Läufe: 1–3 unverändert (Nummern bleiben, das Makefile zitiert „Lauf 3"); Lauf 2 und 3 verlangen zusätzlich **keinen** Vorlauf (Entscheidung 4/5); **neu Lauf 4** = View `cdc.changes` auf abweichende Signatur, Rollout Exit 0, Soll-Signatur, `SELECT`-Recht für `cdc_reader`, bestehende Zeile lesbar, Folgelauf ohne Vorlauf; **neu Lauf 5** = Alt-Tag-Lauf in einer zweiten Datenbank derselben Wegwerf-Instanz (`git archive` des jüngsten `v*`-Tags in ein `mktemp`-Verzeichnis außerhalb des Repo-Baums, dessen Makefile rollt aus, Datenzeile, dann der Arbeitsbaum zweimal, Exit 0, Zeile über `cdc.changes` lesbar, Tag und Exit-Codes gedruckt); **Lauf 6** = der bisherige Lauf 4 (unbekannter Blocker, Exit 8) bleibt der letzte; Fixrunde 2 erweitert ihn um den Fall 6a (siehe Tabelle „Fixrunde 2"). |
| `docs/user/benutzerhandbuch.md` | update | §„Schema aktualisieren": Reihenfolge (Schema-Rollout vor Container-Tausch), Vorlauf-Meldung, Lesefenster (mit Ursprung), „ein am Precheck fehlgeschlagener Rollout ändert das Ziel nicht"; `Version:`-Kopf und Änderungshistorie. |
| `harness/README.md` (§Sensors, Zeile `make schema-rollout`) | update | Vorlauf, Läufe des Guard-Tests (sechs statt vier), Alt-Tag-Lauf. |
| `examples/bootstrap.sh` (Kommentar über dem Existenz-Check und Meldung) | update (Träger nach Suchlauf) | der Kommentar „`make schema-rollout` selbst ist NICHT idempotent" ist seit der zentralen Wache falsch; der Existenz-Check bleibt (akzeptiertes Negativ laut [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md): ein weiterlebendes Demo-Volume mit Alt-Schema bekommt das neue Schema nicht, `make example-demo-down` mit `down -v` baut neu auf). Nur der Kommentar/die Meldung ziehen nach, das Verhalten nicht. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | **nicht geändert** | laut Verdikt bleibt das committete Paar das Ergebnis eines Laufs gegen eine leere Datenbank; kein Upgrade-Lauf schreibt es neu (der Guard-Test sichert beide Dateien beim Start und stellt sie in jedem Ausgang wieder her, Fixrunde 2). |
| `docs/plan/planning/open/slice-backfill-sql-administration.md` u. a. (fremde Pläne mit „vier Läufe") | **gemeldet, nicht geändert** | Träger der bewegten Eigenschaft „Läufe des Guard-Tests", fremde Datei (Planner); Befund im Bericht der Fixrunde. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „Feld- und Spaltenmenge von `cdc.change`, `cdc.changes` und `GET /changes`; die Live-Wege bleiben bei zehn Feldern"; beide Stände gemessen: Parent und Diff):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Spaltenlisten der View in Doku | Parent (Commit `f4ba82ab`, Stand vor diesem Slice): `git grep -n 'committed_at' f4ba82ab -- docs/user spec harness`; Diff-Stand (Commit `e95937cf`, Arbeitsbaum sauber): dasselbe mit `e95937cf` | Beide Stände **6** Treffer (gemessen), davon 3 im Handbuch: Parent `docs/user/benutzerhandbuch.md` Z. 377 (Spaltenliste des SQL-Beispiels unter „Änderungen lesen"), Z. 419 (`SELECT change_id, committed_at` — explizite Zwei-Spalten-Abfrage unter „Aufbewahrung", von `origin` nicht berührt), Z. 639 (Feldliste der `GET /changes`-Antwort unter „Changes lesen"); 3 in `spec/pflichtenheft.md` (Z. 236 synthetische Transaktionen, Z. 603 `SPEC-022`-Antwortzelle — trägt `origin` bereits —, Z. 819 Änderungsverlauf), keiner in `harness/`. Ergänzend `git grep -n -w 'origin' <Stand> -- docs/user`: Parent 1 (`docs/user/releasing.md`, das Git-Remote), Diff-Stand 6 (das Remote plus fünf Zeilen im Handbuch); `-- spec`: je Stand 9 (Spec ist gezogen) | Handbuch-Z. 377 und Z. 639 gezogen (die zwei benannten Stellen), Z. 419 unverändert; Version 1.44 samt Historienzeile; Spec unverändert |
| Anzahl-Formulierungen („zehn Felder", „zwölf Felder") an `GET /changes` | Erste Fassung: `git grep -n 'Felder' <Stand> -- docs/user spec` und `git grep -n -E 'elf Felder\|zwölf Felder\|dreizehn Felder' <Stand> -- docs/user spec`. Fixrunde 2 (Suchraum auf den ganzen Baum erweitert, Wortsuche über Zeilenumbrüche hinweg): `git grep -n -i -w -E 'twelve\|eleven\|thirteen\|zwölf\|zwoelf\|elf\|dreizehn' <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` (Parent `09386619`, Diff-Stand `c54f873a`; die Muster sind mit `|` gemessen und stehen in dieser Tabelle als `\|` maskiert) | `Felder` je Stand **33** Treffer (gemessen, `docs/user spec`), davon `zehn Felder` je Stand **19** (Live-Wege, bleiben bei zehn); `elf`/`zwölf`/`dreizehn Felder` in `docs/user spec` je Stand **0**. Erweiterte Suche: **24** Treffer über den Baum am Parent `09386619` sowie an `f4ba82ab` und `e95937cf`, **25** an `c54f873a` (gemessen; der Zusatztreffer an `c54f873a` ist die Meldungszeile in `docs/plan/planning/open/slice-backfill-sdk-origin.md`, die 25 damit abgeleitet: 24 + 1), davon **4** Feldanzahl-Träger für `GET /changes`, alle in `sdks/`: `sdks/csharp/PgChangeFeed.Client/Sse/Models/Change.cs` Z. 17 („eleven fields"), `sdks/csharp/PgChangeFeed.Client/Nats/Models/Change.cs` Z. 21 („eleven fields"), `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse/model/Change.kt` Z. 19–20 („twelve" am Zeilenende, „fields" in der Folgezeile — eine Zeilen-Suche auf „twelve fields" träfe sie nicht), `sdks/python/pgchangefeed/src/pgchangefeed/models.py` Z. 264 („zwölf des HTTP-Lesezugriffs"); die übrigen 20 Treffer sind fremde Zahlen (zwölf Beispiel-Programme, zwölf SDK-Flächen, Baseline-Text). Zusätzlich `git grep -n -i -E '(ten\|zehn) (fields\|felder)' <Stand> -- sdks examples`: je Stand **16** Treffer (gemessen an `09386619`, `f4ba82ab`, `e95937cf` und `c54f873a`; alle 16 in `sdks/`), davon `Http/Models/Changes.cs` Z. 8 und `http/model/Changes.kt` Z. 8 („the same ten fields as the domain type" — das HTTP-Modell trägt am Parent zwölf Felder, gezählt an `Changes.cs`/`Changes.kt`); `examples/`: **0** Treffer. Die übrigen Treffer beschreiben Live-Wege (zehn Felder) | in `docs/user`/`spec`/`harness`/`examples`: nichts zu ändern; die sechs SDK-Kommentare (die vier Zahlen-Träger und die zwei „same ten fields"-Kommentare) gehören zu `slice-backfill-sdk-origin` (fremde Datei, dort nicht geändert) — an dessen §3 gemeldet (eine eigene Zeile „Kommentare mit Feldanzahl in den SDK-Bäumen"), die SDK-Dateien selbst bleiben in diesem Slice unberührt |
| Harness/Tests mit `SELECT *` gegen `cdc.changes` | `git grep -n 'SELECT \*' <Stand> -- tools test internal` | Je Stand **3** Treffer (gemessen): `test/integration/integration_test.go` Z. 1002 und 1052 (Kommentare zur expliziten Spaltenliste, kein SQL) und `tools/bench-batch-vs-single.sh` Z. 42 (`SELECT count(*) FROM (SELECT * FROM cdc.changes … LIMIT $M) t` — zählt Zeilen, liest die zusätzliche Spalte mit, ohne dass eine Aussage davon abhängt) | keine Änderung: keine der drei Stellen setzt eine Spaltenmenge voraus, die `origin` verletzt |
| Kommentare in den Live-Wegen („dieselben Felder wie `model.Change`") | `git grep -n -E 'dieselben (zehn )?Felder wie\|Feldern wie der Domain-Typ' <Stand> -- '*.go' '*.proto' '*.py'` | Je Stand **8** Treffer (gemessen), die `model.Change` als Feldvorbild nennen: `natsstream/publisher.go` Z. 135, `natsstream/publisher_test.go` Z. 261, `driving/http/sse.go` Z. 29 — mit `Change.Origin` nicht mehr wahr —; `gen/cdc/stream/v1/changestream.pb.go` Z. 31 und `proto/cdc/stream/v1/changestream.proto` Z. 13 (dieselbe Aussage über die gRPC-Nachricht); `driving/http/readchanges.go` Z. 46 (`GET /changes`: trägt `origin`, bleibt wahr); `natsstream/publisher.go` Z. 165 (verweist auf das SSE-Event, bleibt wahr); Python-`models.py` Z. 263 | die drei Go-Kommentare ziehen „ohne `Origin`" nach; die `.proto`-Quelle und die generierte Datei **nicht** (die Zeile 13/31 zu ändern verlangt `make proto-generate`, der Plan legt die Proto-Artefakte und `make generated-sync` als unberührt fest) — gemeldet: die Aussage „mit denselben Feldern wie der Domain-Typ" ist an der gRPC-Nachricht ungenau geworden; der Python-Kommentar gehört zu `slice-backfill-sdk-origin` |
| eingebettete DDL, Report, Rollback | `git ls-tree -r --name-only <Stand> -- tools/schema internal/adapters/driven/postgresstorage/schema.sql` | Je Stand **14** Dateien (gemessen, gleiche Menge); `tools/schema/plan.yaml` und `tools/schema/down.sql` ändern sich im Diff (Ergebnis eines frischen `make schema-rollout` gegen eine leere Datenbank), `schema.sql` und `schema.yaml` ebenso. `schema.sql` ist Träger, nicht nur Test-Hilfe: `store_test.go` und `tableactivation_test.go` bauen das Schema je Test über `ApplySchema` auf; ohne die Spalte dort scheitert `InsertChange` (rot gesehen, Mutation Z12) | `plan.yaml`/`down.sql` regeneriert und committet; `schema.sql` gezogen (Entscheidung, siehe Zeile in der Plan-Tabelle oben) |
| Wegwerf-Client `tools/harness/httpclient` (liest `GET /changes`) | Lesen der Dekodierung (`tools/harness/httpclient/main.go` Z. 126–129) und `git grep -n -i -E 'DisallowUnknown\|UnmappedMemberHandling\|FAIL_ON_UNKNOWN\|ignoreUnknownKeys\|extra=.forbid' <Stand> -- '*.go' '*.cs' '*.kt' '*.py'` | Der Client dekodiert über `json.Unmarshal` in eine Struktur ohne `origin` — nicht strikt; strikte Dekoder-Muster je Stand **0** Treffer (gemessen) im ganzen Baum | keine Änderung nötig |
| E2E-Abdeckungs-Tabelle (`Datei:Zeile`-Anker) | `git diff --stat f4ba82ab e95937cf -- test/integration tools/harness/run-integration-tests.sh docs/user/e2e-abdeckung.md` | Leer (gemessen, 0 Zeilen): der Diff berührt weder das Testpaket noch den Runner noch die Tabelle | keine Regeneration nötig; `make test-integration` lief nach dem Diff-Stand real (Exit 0) und meldete „E2E-Abdeckungstabelle unverändert" |

**Messung zu Risiko §6 „d-migrate konvergiert nicht" (Implementer; der Ausgang gehört der Closure).** Alle Läufe gegen einen Wegwerf-PostgreSQL-18-Container (`postgres:18-alpine`, Digest aus `PG_TEST_IMAGE`), Rollout-Werkzeug d-migrate im gepinnten Image (Makefile-Variable `D_MIGRATE_IMAGE`); Ausgabe je Lauf gedruckt, hier die Exit-Codes:

- **Frische Ziel-Datenbank, neues Schema, `make schema-rollout` zweimal hintereinander:** beide Läufe Exit 0 (gemessen). Der zweite Lauf nimmt den `--allow-destructive`-Pfad der Idempotenz-Wache (Meldung in der Ausgabe); die View `cdc.changes` führt danach `origin` als letzte Spalte. Der erste dieser Läufe erzeugt die committeten `plan.yaml`/`down.sql` (14 Operationen, `change` mit der Spalte `origin`, View `changes` mit `COALESCE(c.origin, 'wal') AS origin`).
- **Ziel-Datenbank mit dem Schema des Parent-Commits `f4ba82ab` (zuvor per `make schema-rollout` ausgerollt), danach das neue Schema:** `make schema-rollout` endet mit Exit 2 (der Precheck-Lauf meldet Exit 8). Blocker im Precheck-Report: `MANUAL_ACTION_REQUIRED` für `ReplaceView` der View `changes`, Diagnose `VIEW_SIGNATURE_INCOMPATIBLE` („CREATE OR REPLACE VIEW is only renderable when view columns keep the same count, order, names and visible types"). Die Operation `AddColumn` für `change.origin` steht im selben Plan als regulär renderbar; die Wache (`tools/schema/rolloutguard`) lehnt die Blocker-Klasse `MANUAL_ACTION_REQUIRED` ab und läuft ohne `--allow-destructive`. **Ein bereits mit dem alten Schema ausgerolltes Ziel ließ sich damit (Stand vor [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)) über `make schema-rollout` nicht auf den neuen Stand heben.**
- **Derselbe Ausgangsstand, nach einem manuellen `DROP VIEW cdc.changes` vor dem Rollout:** Exit 0, ein zweiter Lauf Exit 0; eine zuvor ohne Spalte gespeicherte Zeile liest über die View als `wal` (`NULL` in `cdc.change.origin`), die Rechte auf die View sind danach wieder gesetzt (`\dp cdc.changes` zeigt `cdc_reader=r`). Diese Messung führte zur in §4 vorab benannten Architect-Frage; [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) beantwortet sie mit dem Vorlauf-Schritt im Makefile-Target (siehe „Ausgang der Messung" unten).
- **Reichweite der Aussage:** die Test-Tiers (`make test-store`, `make test-replication`, `make test-integration`), `make example-demo-up` und der Compose-Aufbau legen ihr Schema laut `harness/README.md` §Sensors frisch an; der bestehende Rollout-Check in `examples/bootstrap.sh` überspringt den Rollout gegen ein bereits migriertes Ziel — dass ein dort weiterlebendes Alt-Schema `InsertChange` scheitern lässt, ist abgeleitet, nicht gemessen. Akzeptiertes Negativ laut [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md): ein weiterlebendes Demo-Volume mit Alt-Schema bekommt das neue Schema nicht, `make example-demo-down` (mit `down -v`) baut es neu auf; der Kommentar in `examples/bootstrap.sh` nennt das.

**Ausgang der Messung (Fixrunde zu [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md); Diff-Stand der Fixrunde = Commit `fa954778`, gemessen im Lauf, Exit-Codes gedruckt).** Wegwerf-PostgreSQL 18 (`PG_TEST_IMAGE`, `wal_level=logical`), d-migrate im Makefile-Pin (`D_MIGRATE_IMAGE`), Alt-Stände als `git archive`-Kopien in einem Verzeichnis unter dem Scratchpad (nie im Repo-Baum):

- **Parent-Schema `f4ba82ab` → Arbeitsbaum:** der Rollout des Parent-Stands mit dessen eigenem Makefile endet mit Exit 0; eine Zeile (Quelle → Change, ohne `origin`) ist über `cdc.changes` lesbar (Zählung 1). `make schema-rollout` im Arbeitsbaum: **Exit 0**, die Ausgabe trägt `schema-rollout: Vorlauf (ADR-0114) - View-Signatur-Aenderung, DROP VIEW cdc.changes` (1 Treffer); der zweite Lauf: **Exit 0**, ohne Vorlauf-Meldung (0 Treffer). Danach liest `SELECT change_id||'|'||origin FROM cdc.changes` `ch|wal`, `has_table_privilege('cdc_reader','cdc.changes','SELECT')` ist `t`.
- **Alt-Tag `v0.1.2` → Arbeitsbaum** (Lauf 5 des Guard-Tests, dessen gedruckte Zeile): „Tag v0.1.2: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar" — der Guard-Test insgesamt: Exit 0 (sechs Läufe; der Exit-8-Lauf endet als `make` Fehler 8 des d-migrate-Precheck-Laufs, der Skript-Vergleich meldet `make` Exit 2).
- **Compose-Rollout-Pfad** (das umgebaute Rezept läuft dort gegen eine frische Datenbank): `make image` Exit 0, danach `make test-integration` Exit 0 (E2E-Abdeckungstabelle unverändert); `make test-store` Exit 0 (DB-Adapter-Coverage 77,13 % gegen die Schwelle 70 %, gedruckt im Lauf).
- **Unit-Tests** `tools/schema/rolloutguard` (`make test`, gemessen in Fixrunde 2): 19 Tests, Exit 0. **Lauf 6 des Guard-Tests** trägt zwei Fälle: **6a** (nicht deklarierte Funktion) bindet die Bekannt-Liste des Guards end-to-end — mit deaktivierter `knownForeignObjects`-Prüfung endet der Guard-Test mit Exit 1 („Lauf 6a lief durch (Exit 0), obwohl ein unbekannter destruktiver Blocker vorlag", gesehen), unmutiert Exit 0 mit der gedruckten Zeile „Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit d-migrate-Exit 8"; **6b** (`ADD COLUMN`, DropColumn-Blocker) belegt den Abbruch durch d-migrate gegen einen real gemeldeten Blocker, **nicht** die Bekannt-Liste — d-migrate bricht diesen Blocker auch unter `--allow-destructive` selbst ab (mit entfernter Prüfung gesehen: Exit 0 des Guard-Tests in der Fassung vor Fixrunde 2). Die Klassen-Merkmale und der Blocker ohne Operationen tragen die Unit-Tests.

**Fixrunde 2 zu `docs/reviews/review-slice-backfill-change-origin.md`** (0 HIGH, 2 MEDIUM, 3 LOW; die Zeilen sind **vor** dem Sensor-Lauf der Fixrunde eingetragen; Reviewer-Report unverändert, Record):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/rolloutguard/guard.go`, `guard_test.go` | update (F-1) | ein Blocker `DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION` ohne `operationIds` wird wie der `MANUAL_ACTION_REQUIRED`-Zweig abgelehnt (leere Operationsliste belegt keine bekannte Operation; alles oder nichts nach [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 2). Vor der Änderung am Report-Format geprüft: die real gemessenen destruktiven Blocker tragen `operationIds` (Fixture `realViewSignatureReport` in `guard_test.go` — gekürzter realer Precheck-Report — und `knownBlockedReport`; der Idempotenz-Fall der sechs Fremdobjekte läuft im Guard-Test Lauf 2/3 weiter Exit 0). Neuer Test `TestDecideRefusesDestructiveBlockerWithoutOperations` (allein, neben Fremdobjekten, neben View-Signatur). |
| `tools/harness/run-schema-rollout-guard-test.sh` | update (F-2, F-6) | **Entscheidung F-2:** Lauf 6 wird um einen stärkeren Negativfall erweitert, ohne Schema-Änderung im Repo und ohne Eingriff in die Wächter-Semantik: **6a** eine nicht deklarierte Funktion `cdc.zz_rolloutguard_unbekannt()` — dieselbe Klasse wie die sechs bekannten Fremdobjekte (Blocker `DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`), aber nicht auf der Bekannt-Liste; ließe die Wache sie durch, riefe das Target `--execute --allow-destructive` auf und d-migrate löschte die Funktion — der Beleg ist ihr Fortbestehen nach einem Lauf, der mit „Error 8" (d-migrate-Exit 8) endet. **6b** der bisherige `ADD COLUMN`-Fall (DropColumn-Blocker, d-migrate bricht auch unter `--allow-destructive` selbst ab) bleibt als Beleg des Abbruchs gegen einen real gemeldeten Blocker, wird im Kopf aber nicht mehr als Beleg der Bekannt-Liste geführt. Die Assertion prüft „Error 8" statt nur „ungleich 0"; „Exit 2" (make-Exit) und „Exit 8" (d-migrate-Exit) sind getrennt benannt. **F-6:** das Skript sichert `tools/schema/plan.yaml` und `down.sql` beim Start in ein `mktemp`-Verzeichnis und stellt sie in `cleanup` (jeder Ausgang) wieder her, sodass der Arbeitsbaum nach dem Lauf unverändert bleibt; die Läufe des Parent-Stands (1–3) schreiben dieselben Dateien und werden davon nicht verfälscht, weil die Sicherung vor dem ersten Lauf entsteht. Die Zahl „sechs Läufe" bleibt (Lauf 6 trägt zwei Fälle). |
| `harness/README.md` (Zeile `make schema-rollout`) | update (F-2) | Negativ-Abbruch beschreibt die zwei Fälle wahrheitsgemäß; die Bekannt-Liste und der Blocker ohne Operationen sind als Unit-Test-getragen benannt. |
| `Makefile` (Kommentarblock über `schema-rollout`) | update (F-3, F-5, F-8) | F-3: der erste Absatz („ausschließlich die sechs bekannten Objekte … nur dann `--allow-destructive`") wird auf den wahren Stand gebracht (jeder Blocker ist ein bekanntes Fremdobjekt oder eine View-Signatur-Änderung; `--allow-destructive` läuft, wenn mindestens ein bekanntes Fremdobjekt blockiert). F-5: `DROP VIEW` verwirft die ACL der View, nur `cdc_reader` wird durch `nacharbeit-roles.sql` neu gesetzt. F-8: der Vorlauf adressiert das Schema `cdc` fest, Vorbedingung `search_path = cdc`. |
| `docs/user/benutzerhandbuch.md` | update (F-5, F-8) | zwei Aufzählungspunkte in §4 „Schema aktualisieren" (eigene Rechte an anderen Rollen gehen verloren und werden vom Betreiber erneut gesetzt; feste Adressierung des Schemas `cdc`); `Version:` 1.46 und Änderungshistorie. |
| `docs/plan/planning/open/slice-backfill-sdk-origin.md` | update (F-4, gemeldet) | eine §3-Zeile „Kommentare mit Feldanzahl in den SDK-Bäumen" mit den sechs Fundstellen; die SDK-Dateien selbst bleiben unberührt. |
| `tools/harness/run-integration-tests.sh` Z. 2711 (F-7) | **nicht geändert** | wie im Plan begründet: Träger der Zeilenanker von `docs/user/e2e-abdeckung.md`; unverändert gemeldet. |
| `ADR-0114` (F-5) | **Lücke gemeldet, nicht geändert** | die ADR (Accepted) nennt in Entscheidung 3 „die Rechte setzt `nacharbeit-roles.sql` im selben Lauf", benennt aber nicht, dass `DROP VIEW` auch Rechte anderer Rollen als `cdc_reader` verwirft; eine Ergänzung wäre eine neue ADR bzw. Zitat der Konsequenz — Meldung an den Architect. |

**§3.13-Suchlauf der Fixrunde** (bewegte Eigenschaft: „was `make schema-rollout` gegen ein bereits migriertes Ziel tut — Vorlauf, Klasse View-Signatur, `--allow-destructive`, Läufe des Guard-Tests"; Parent = `ae97247a`, Diff-Stand = `fa954778`, beide gemessen):

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Beschreibungen des Guards und des Rollouts | `git grep -n -E 'allow-destructive\|rolloutguard\|knownForeignObjects' <Stand> -- docs/user harness examples Makefile tools AGENTS.md README.md README.de.md` | Parent **33** Treffer in 8 Dateien, Diff-Stand **42** Treffer in 9 Dateien (gemessen); die neue Datei ist `examples/bootstrap.sh` | `harness/README.md` Z. 151, Makefile-Kommentar, Kommentare in `guard.go`/`report.go`/`main.go`, Kopf des Guard-Tests und `examples/bootstrap.sh` gezogen; `AGENTS.md` §3.14 nennt die zentrale Wache, bleibt wahr |
| Zahl der Läufe des Guard-Tests | `git grep -n -E 'vier Läufe\|Lauf [0-9]/4\|Läufe des Guard-Tests\|sechs Läufe\|Lauf [0-9]/6' <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` | Parent **9** Treffer (Skript-Kopf und Läufe 1/4–4/4, `harness/README.md` Z. 151, zwei Pläne in `open/`, ein fremder Treffer in `harness/sensors/coverage-gate.md` Z. 89 zu `slice-093`), Diff-Stand **11** Treffer (gemessen; Läufe 1/6–6/6 statt 1/4–4/4) | Skript und `harness/README.md` gezogen; `docs/plan/planning/open/slice-backfill-sql-administration.md` Z. 193 und `slice-transformationen-antragsweg-schema.md` Z. 171 nennen die Läufe nur als **Such-Auftrag** („Zahl und Beschreibung nachziehen, falls sich die Läufe ändern") — fremde Pläne, gemeldet, nicht geändert; der Treffer in `coverage-gate.md` meint `slice-093`, nicht diesen Guard |
| „nicht idempotent gegen ein migriertes Ziel" | `git grep -n -i -E 'nicht idempotent' <Stand> -- docs/user harness examples Makefile tools` | Parent **2** Treffer (beide `examples/bootstrap.sh`), Diff-Stand **0** | Kommentar und Meldung in `examples/bootstrap.sh` gezogen (der Existenz-Check bleibt, akzeptiertes Negativ laut [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)); `examples/README.md` trägt keine Aussage dazu (Suchlauf `-i -E schema|rollout|volume|down` gegen den Parent: Treffer nur in der Demo-Beschreibung Z. 23–24, bei `example-demo-down` Z. 44 und in den `-schema`-Flags der Beispiel-Aufrufe) |
| Beschreibungen „Blocker/zweiter Rollout" außerhalb der Suchräume oben | `git grep -n -E 'schema-rollout' <Stand> -- '*.sh'` und Lesen von `tools/harness/run-integration-tests.sh` Z. 2711 | Z. 2711 nennt „den in [`ADR-0058`](../../adr/0058-testansatz-fuenf-luecken.md) vorgesehenen, real blockierten zweiten `make schema-rollout`-Lauf (… Exit 8 auf vier Fremdobjekten …)" — die Aussage beschreibt den Ist-Stand der Wache nicht mehr (bereits seit der zentralen Idempotenz-Wache) | nicht geändert: `tools/harness/run-integration-tests.sh` ist das Testpaket-Runner-Skript, dessen Phasen-Deklarationen `docs/user/e2e-abdeckung.md` trägt (Risiko §6, vierter Punkt) — gemeldet |

**Nachzug außerhalb des Slice-Umfangs (Kommentar- und Target-Dokument-Zug,
[`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)):** Der
Kommentarblock über `schema-rollout` im `Makefile` und die Zeile
`make schema-rollout` in `harness/README.md` (die beiden „update"-Zeilen oben)
tragen den Ist-Vertrag nicht mehr selbst; er steht in
`harness/targets/schema-rollout.md`, beide verweisen dorthin. Das Rezept des
Targets ist unverändert bis auf den Meldungstext bei `--allow-destructive`
(„nur" entfällt, weil die Meldung auch bei gleichzeitiger View-Signatur-Änderung
druckt); der Guard-Test zieht seine Text-Erwartung und die Lauf-Kennung der
Nebenfall-Überschrift („Lauf 4/6") gleich.

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

DoD vollständig + `make gates` grün + `make test`, `make test-store`,
`make schema-rollout` (zweimal, frisch und als Upgrade über das Parent-Schema)
und `tools/harness/run-schema-rollout-guard-test.sh` (sechs Läufe, darunter der
Alt-Tag-Lauf) real grün + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **d-migrate konvergiert die additive Spalte oder die View-Änderung nicht.**
  Die Anlage einer Spalte an einer bestehenden Tabelle konvergiert nach dem
  Kopfkommentar von `tools/schema/nacharbeit-administration.sql` (Spalte
  `column_name`, dort als real gemessen, Exit 0); der Folgelauf gegen eine
  bereits migrierte Datenbank mit `ReplaceView` steht in `harness/README.md`
  §Sensors, Zeile `make schema-rollout`. Beides ist für **diese** Spalte und
  **diese** View ungemessen. *Erwartet, zu belegen durch:* zweiter
  `make schema-rollout`. Der Beleg der Fixrunde steht in §3 („Ausgang der
  Messung", [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)).
  **Ausgang:** *eingetreten* — nicht an der Spalte (die additive Spalte
  konvergierte, Operation `AddColumn` regulär renderbar), sondern an der
  Signaturänderung der View: gegen ein Ziel mit dem Parent-Schema blockierte der
  Precheck mit `ReplaceView`/`VIEW_SIGNATURE_INCOMPATIBLE`, `make schema-rollout`
  endete mit Exit 2 (§3, „Messung zu Risiko §6"). Im Slice gelöst, kein
  Carveout und kein Folge-Slice: Architect-Verdikt
  (`docs/reviews/architect-verdict-schema-rollout-view-signatur.md`),
  [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) und die
  Fixrunden (Vorlauf im Target, Wache `tools/schema/rolloutguard`, Guard-Test
  Läufe 4 bis 6 mit dem Alt-Tag-Lauf). Belege: §3 „Ausgang der Messung" (Parent
  `f4ba82ab` → Arbeitsbaum Exit 0 und 0, Tag `v0.1.2` → Arbeitsbaum Exit 0, 0
  und 0) und die eigenen Läufe des Verifiers (Verifikation §3.2: Parent-Lauf
  `WORK_ROLLOUT_1_EXIT=0`, `WORK_ROLLOUT_2_EXIT=0`, Guard-Test EXIT=0). Die
  Klasse trägt `BEO-PGC/d-migrate-nacharbeit` (siebter Beleg).
- **View und Lesepfad driften auseinander** (`BEO-PGC/lese-doppelquelle`,
  verkörpert): die View trägt `origin`, `SelectChanges` nicht (oder umgekehrt).
  *Erwartet, zu belegen durch:* `TestE2EChangesViewMatchesReadChanges` und der
  Store-Test. **Ausgang:** *entfallen* — nicht eingetreten: View und
  `SelectChanges` tragen `origin` beide als letzte Spalte mit demselben
  `COALESCE(…, 'wal')`, und `TestChangesViewCarriesOriginLikeReadChanges`
  (`make test-store`) hält sie gegeneinander; die Mutationen an der View
  (`COALESCE` entfernt, Review M5, Verifikation M5) und an `SelectChanges`
  (Review M6) färben ihn rot. `TestE2EChangesViewMatchesReadChanges` bleibt
  bewusst unverändert (viertes Risiko unten).
- **Ein Client dekodiert strikt und bricht am neuen Feld.** Die SDK-Decoder
  ignorieren unbekannte Felder (gelesen für C# `System.Text.Json`, Kotlin Gson,
  Python `json` mit expliziten Feldern); ein realer Beleg fehlt bis
  `sdk-origin`. *Erwartet, zu belegen durch:* `make test-sdk-*-integration`-Läufe
  bzw. der Slice `sdk-origin`. **Ausgang:** *weiter offen* — der Slice maß nur,
  was er messen kann: 0 strikte Dekoder-Muster im Baum (§3, gemessen), der
  Wegwerf-Client `tools/harness/httpclient` ist nicht strikt (gelesen); ein
  realer Beleg mit den SDK-Decodern fehlt bis `slice-backfill-sdk-origin`, dessen
  §6 (erster Punkt) das Risiko mit Fixture-Test trägt. Register:
  `BEO-PGC/sdk-decoder-verhalten-am-neuen-feld-ungemessen` (offen, 1×).
- **Die Zeilenzahl-Anker der E2E-Abdeckung wandern**, sobald der Zug das
  Testpaket berührt. *Erwartet, zu belegen durch:* der Suchlauf-Eintrag zur
  Abdeckungs-Tabelle. **Ausgang:** *entfallen* — nicht eingetreten: der Diff
  berührt weder das Testpaket noch den Runner noch die Tabelle (§3, gemessen;
  `git diff --stat f4ba82ab..HEAD -- test/integration tools/harness/run-integration-tests.sh docs/user/e2e-abdeckung.md`
  leer, vom Verifier nachgemessen), und `make test-integration` meldete
  „E2E-Abdeckungstabelle unverändert". Der eine Träger, dessen Änderung die
  Anker verschöbe — der Kommentar zum zweiten Rollout-Lauf im Runner (Review
  F-7) —, ist bewusst unverändert und als Zeile in
  `slice-backfill-e2e` (§3) adressiert.

## 7. Closure-Notiz

- **Was hat funktioniert:** das Feld `origin` ist von der Domäne bis zu View
  und `GET /changes` durchgängig, die Live-Wege bleiben bei zehn Feldern
  (Verifikation §3.1). Die Rollen-Kette lief in getrennten Kontexten:
  Review (0 HIGH · 2 MEDIUM · 3 LOW · 5 INFO), Fixrunde 2 (F-1 bis F-6),
  Verifikation „Bestätigt" mit eigenen Läufen — `make gates` EXIT=0, der
  Guard-Test EXIT=0, ein selbst gefahrener Upgrade-Lauf vom Parent-Schema in
  den Arbeitsbaum (Exit 0 und 0), Eingabeseiten-Mutationen am Diff rot
  (Verifikation §5). Der Upgrade-Blocker der View-Signatur wurde nicht im
  Alleingang umgangen: die in §4 vorab benannte Rückführungs-Frage ging als
  Architect-Frage an [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
  und kam als Fixrunde im selben Slice zurück, ohne Rückführung nach `open/`.
- **Was ging anders als geplant:** der Plan wuchs um die Nachträge der
  Umsetzung (`translate.go`, der zweite Schema-Träger `schema.sql`, drei
  Live-Weg-Kommentare, `examples/bootstrap.sh`; §3) und um zwei Fixrunden:
  die Umsetzung von [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
  (Vorlauf im Target, Klasse „View-Signatur" in `tools/schema/rolloutguard`,
  Guard-Test von vier auf sechs Läufe, Handbuch „Schema aktualisieren" bis
  Version 1.46) und die Behebung der Review-Funde F-1 bis F-6. Der Vertragstest
  `TestE2EChangesViewMatchesReadChanges` blieb bewusst unberührt (die
  Zeilen-Anker der E2E-Abdeckung), die `origin`-Parität trägt der
  Store-Tier-Test. Nach der Verifikation kam ein Doku-/Kommentar-Zug
  außerhalb des Slice (Commit `c54f873a`: `harness/targets/schema-rollout.md`
  trägt den Ist-Vertrag des Targets, der Makefile-Kommentar und die README-Zeile
  verweisen dorthin); er ist kein Diff-Bestandteil des geprüften Umfangs, die
  Verweise des Plans (§3, „Nachzug außerhalb des Slice-Umfangs") zeigen auf ihn.
- **Verifier-Beobachtungen (V-1 bis V-8):** *V-1* gezogen: die Suchlauf-Zelle
  „Anzahl-Formulierungen" nennt 16 (statt 18) samt den vier Ständen; der
  Planner maß bei der Closure `git grep -n -i -E '(ten|zehn) (fields|felder)' <Stand> -- sdks examples`
  an `09386619`, `f4ba82ab`, `e95937cf` und `c54f873a`: je 16, alle in `sdks/`.
  *V-2* gezogen: die erweiterte Suche trägt 24 an den ersten drei Ständen und 25
  an `c54f873a`; die 25 ist als abgeleitet gekennzeichnet (24 + Meldungszeile in
  `slice-backfill-sdk-origin`). *V-3* ist die ADR-Lücke (unten). *V-4*
  (Blocker-Grund `MANUAL_ACTION_REQUIRED` enger als der ADR-Wortlaut) zulässig,
  kein Nachzug. *V-5* durch `c54f873a` behoben: der Meldungstext des Rezepts
  druckt „`bekannte Fremdobjekt-Blocker (ADR-0043)`" ohne „nur" (`Makefile`
  Z. 264), der Guard-Test erwartet diesen Text und die Überschrift des
  Nebenfalls heißt „Lauf 4/6" (`git grep -n 'nur bekannte' HEAD -- Makefile tools/harness`:
  0 Treffer). Rest: der Kommentar am Feld `allowDestructive` in
  `tools/schema/rolloutguard/guard.go` sagt weiter „(nur bekannte Fremdobjekte
  blockieren)" und ist neben einem View-Signatur-Blocker ungenau — Kommentar,
  keine Zusage; nicht geändert (Code-Datei, kein Planner-Zug), im Bericht der
  Closure gemeldet. *V-6* siehe Register (`test-schreibt-in-committete-datei`).
  *V-7* gezogen (die Reconciliation-Zeile steht auf `[x]`). *V-8* regelkonforme
  Übergabe, kein Nachzug.
- **Steering-Loop-Eintrag (Lerneintrag):** *neuer Sensor-Umfang:* der Slice
  liefert nicht nur ein Feld, sondern erweitert einen Sensor — die Wache
  `tools/schema/rolloutguard` kennt die Klasse „View-Signatur" (alles oder
  nichts), und der Guard-Test `tools/harness/run-schema-rollout-guard-test.sh`
  trägt Lauf 4 (abweichende Signatur), Lauf 5 (Alt-Bestand: der Schema-Stand des
  letzten Release-Tags wird ausgerollt, danach der Arbeitsbaum) und Lauf 6 (der
  Negativ-Abbruch bleibt der letzte); der Alt-Tag-Lauf schließt für das Schema
  den Trigger 1 von [`ADR-0064`](../../adr/0064-lh-qa-ops-005-testansatz-korrektur.md)
  ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
  Entscheidung 7). *Geschärfte Anwendung (keine neue Regel):* ein Beleg trägt nur
  den Satz, den seine Mutation rot färbt — Lauf 6 nannte sich „Beleg, dass die
  Wache nicht pauschal durchlässt", die Mutation der Bekannt-Liste blieb grün,
  weil d-migrate den Blocker selbst abbricht (Review F-2); nach der Fixrunde
  färbt Lauf 6a sie rot (der Verifier maß Exit 1 des Guard-Tests mit
  deaktivierter Prüfung), und 6b heißt, was er ist. Das ist die verkörperte
  Klasse `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` in der Form
  „Assertion" (neunter Beleg, Reviewer-HIGH-Punkt „Beleg trägt seinen Satz
  nicht" hat gegriffen). *Benannte Spec-Lücke mit Adresse:*
  [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
  Entscheidung 3 sagt, die Rechte setze `nacharbeit-roles.sql` „im selben Lauf";
  das gilt für `cdc_reader` (dort steht das einzige `GRANT` auf `cdc.changes`),
  nicht für Rechte, die ein Betreiber außerhalb des Repos vergibt — `DROP VIEW`
  verwirft die gesamte Rechteliste (abgeleitet aus der PostgreSQL-Semantik,
  nicht am Lauf gemessen). Die ADR ist `Accepted` und wird nicht überschrieben
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); die Grenze steht in
  `docs/user/benutzerhandbuch.md` (§4 „Schema aktualisieren", „Eigene Rechte"),
  in `harness/targets/schema-rollout.md` (§Grenze Punkt 3) und in §3 dieses
  Plans. Die Adresse der Entscheidung ist der Architect (Verifikation V-3);
  das Register führt sie als offenen Eintrag
  `BEO-PGC/adr-aussage-breiter-als-ihre-messung` — dieser Slice ordnet keinen
  ADR-Auftrag an.
- **Beobachtungs-Register (`../observations/`):** je Vorkommen eine
  `evidence/`-Datei `slice-backfill-change-origin.md` (Zähler = Zahl der
  Dateien, real ausgezählt). *Bestehende Klassen:*
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (Review F-2) **9×**;
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (Review F-3, Verifikation V-5)
  **4×** — über der Schwelle, der Ausgang bleibt beim Lese-Schritt der Closure
  von [welle-backfill-bestand](welle-backfill-bestand.md);
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Verifikation V-1, V-2)
  **15×**; `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Review F-4, F-7, F-10)
  **28×**; `BEO-PGC/d-migrate-nacharbeit` (die View-Signaturänderung als fünfte
  Objektklasse) **7×**. *Neue Klassen (je 1×, offen):*
  `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (Review F-5, Verifikation V-3),
  `BEO-PGC/test-fatalf-mit-offener-rows-blockiert-pool-close` (Implementer-Befund,
  übernommen, nicht vom Planner reproduziert; die zwei vorbestehenden
  `pool.Query`-Stellen in `sqlviews_test.go` tragen dieselbe Form),
  `BEO-PGC/test-schreibt-in-committete-datei` (Review F-6, Verifikation V-6),
  `BEO-PGC/sdk-decoder-verhalten-am-neuen-feld-ungemessen` (das dritte Risiko
  aus §6, Ausgang *weiter offen*). *Benannt, nicht gezählt:* die Slice-Bezüge in
  Makefile-Kommentaren, die der Nutzer nach der Verifikation meldete (Zug
  `c54f873a` außerhalb eines Slice, keine `evidence/`-Datei — Vermerk im `state.md`
  von `BEO-PGC/slice-chronik-in-code-kommentar`); F-1 (Blocker ohne Operationen
  gilt als bekannt, ein Vakuum-Zweig der Wache) — im Slice behoben, erstes
  Auftreten ohne Klasse im Register. Keine Klasse erreicht mit diesem Slice
  neu die Schwelle von 3×.
- **Validator (Modul 8):** entfällt ausdrücklich — der Slice ist eine interne
  Erweiterung um ein Herkunftsfeld ohne neues End-Nutzer-Verhalten außer dem
  additiven Feld `origin` in `cdc.changes` und `GET /changes`; das Feld hat
  noch keinen Erzeuger außer `wal` (jede Zeile trägt `wal`, ein fehlender Wert
  liest als `wal`), es gibt also keinen Nutzerfall, den ein Validator gegen einen
  realen Bedarf halten könnte. Der Bedarf aus
  [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) wird erst mit den
  Folge-Slices und dem Wellen-Beleg (`slice-backfill-e2e`) validierbar. Kein
  stilles Überspringen.
- **Closure-Notiz-Review (`.harness/skills/closure-note-reviewer.md`):** eine
  getrennte Rolle im frischen Kontext, kein Schritt der Planner-Closure; der
  Skill prüft Slices in `done/` und greift daher erst nach dem `git mv` — hier
  nicht ausgeführt.
- **Folge-Slices:** keine neuen. Meldungen mit Adresse: `slice-backfill-sdk-origin`
  (die sechs SDK-Kommentare und der reale Decoder-Beleg, §3 und §6 dort),
  `slice-backfill-e2e` (der Kommentar zum zweiten Rollout-Lauf im Runner,
  Review F-7, eigene §3-Zeile). Die Aussage
  „mit denselben Feldern wie der Domain-Typ" an der gRPC-Nachricht
  (`proto/cdc/stream/v1/changestream.proto`, Review F-10) verlangt
  `make proto-generate` und bleibt eine gemeldete Ungenauigkeit ohne Träger-Slice.
- **Risiken aus §6:** je ein Ausgang am Ort — Risiko 1 (Konvergenz)
  **eingetreten**, im Slice gelöst ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md),
  Fixrunden, Guard-Test Läufe 4 bis 6; Klasse `BEO-PGC/d-migrate-nacharbeit`);
  Risiko 2 (Doppelquelle View/Lesepfad) **entfallen**; Risiko 3 (strikter
  Decoder) **weiter offen**, adressiert in `slice-backfill-sdk-origin` §6 und im
  Register; Risiko 4 (Zeilen-Anker der E2E-Abdeckung) **entfallen**.
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure. (a) Anker: der Sensor-Umfang ist
  am Ort verkörpert (Wache und Guard-Test, `harness/targets/schema-rollout.md`,
  [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)); die
  geschärfte Anwendung verkörpert nichts neu. (b) Folge-Slice: keiner neu, die
  Meldungen tragen Adressen (siehe oben). (c) Register: alle genannten Kennungen
  existieren als Verzeichnis, ihre `evidence/`-Verzeichnisse tragen die Datei.

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
