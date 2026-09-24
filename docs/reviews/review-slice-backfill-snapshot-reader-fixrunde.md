# Review-Report: slice-backfill-snapshot-reader, Nachprüfung der Fixrunde — 2026-09-24

**Review-Art:** Code — die Fixrunde behebt den Erst-Review
(`docs/reviews/review-slice-backfill-snapshot-reader.md`, 2 HIGH, 4 MEDIUM, 6 LOW,
6 INFO; Record) und führt dabei neuen, unabhängig ungeprüften Code ein
(Text-Ergebnisformat des Lesepfads, Unterpaket `snapshotlogic`, Typ-Paritätstest,
Listen-Prüfskript, Slot-Reserve-Phase im Tier). Geprüft gegen Plan, ADRs und
`AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist
Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-backfill-snapshot-reader`, Diff-Range
`d7539e2c..HEAD` (`406810e4`): Commits `d308b3a5` (Code, Tests), `eda2e41f`
(Harness, Sensor-Doku), `bf05067d` (Plan), `406810e4` (Suchlauf-Zahlen), dazu der
Slice-Plan am Stand `HEAD`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09" (seither um
weitere HIGH-Klassen ergänzt). **Modell:** claude-sonnet-5 · **Datum:** 2026-09-24.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-snapshot-reader` (§1 Ziel, §2 DoD als Prüfmaßstab für
  Plan-Zusagen, §3 Plan, Suchlauf-Feld der Fixrunde und Messbefunde, §6 Risiken)
- [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)
  (Festlegung 1 bis 4, Folgepflichten 1 bis 5, Fitness Function; `Accepted`, Constraint),
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 1
  und 2, [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3,
  [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 1 und 3, [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md),
  [`ADR-0030`](../plan/adr/0030-testpyramide.md)
- [`LH-FA-CAP-009`](../../spec/lastenheft.md),
  [`LH-FA-CAP-008`](../../spec/lastenheft.md),
  [`SPEC-008`](../../spec/pflichtenheft.md) (Fehlerklassen),
  [`SPEC-012`](../../spec/pflichtenheft.md) (PostgreSQL 17 und 18)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.6, §3.7, §3.10, §3.11, §3.12, §3.13),
  `harness/conventions.md` (MR-000/MR-001)
- Vorheriger Review desselben Moduls: `docs/reviews/review-slice-backfill-snapshot-reader.md`
  (Record); Report-Gerüst `docs/reviews/review-report.template.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht
übernommen; Exit-Codes ungepiped gesichert):

- **Gates am Stand `406810e4`:** `make test` Exit 0 (`postgressnapshot` ok,
  `postgressnapshot/snapshotlogic` ok); `make a-check` Exit 0 („gesamt: 0 Befund(e)");
  `make coverage-gate` Exit 0 (gedruckt „Coverage 83.30% erfüllt Schwelle 80%"; der
  Vorlauf druckt „db-package-lists-check: OK — … dieselben 4 Pakete"); `make gates` Exit 0
  (Arbeitsbaum mit dem angelegten Report-Gerüst, `docs-check` „gesamt: 0 Befund(e)").
- **Unit-Nenner** (die Stufe `coverage` nachgestellt: `go test -coverpkg=<Paketliste der
  Stufe> -covermode=atomic`, dedupliziert über die Block-Position, Parent und Kopf je aus
  `git archive` in ein Wegwerfverzeichnis außerhalb des Repos): Parent `d7539e2c`
  **1691 von 2040** (gedruckt `82.9%`), Kopf **1733 von 2082** (gedruckt `83.2%`); das
  Unterpaket `snapshotlogic` trägt **42 von 42**. Der END-verankerte Filter
  `(^|/)(postgresstorage|postgresack|postgressnapshot|replication/receive)$` trifft
  `…/postgressnapshot/snapshotlogic` nicht (die Paketliste der Stufe im Testlauf nennt es
  unter `-coverpkg`). Gegenprobe Filter ohne `postgressnapshot`: 1734 von 2204 = 78,68 %
  (gedruckt `78.7%`), unter der Schwelle 80.
- **DB-Adapter-Coverage** (frischer `make test-store`, danach `make test-replication`,
  PostgreSQL 18): gedruckt „80.07% (gedeckt 651 von 813 Statements; Profile gemergt:
  store,replication)"; Anteile aus dem gemergten Profil abgeleitet: `postgresstorage`
  348/472, `postgresack` 32/32, `postgressnapshot` 118/122, `replication/receive` 153/187;
  Replication-Profil 245 Positionen × 3 = 735 Zeilen; „nur das erste Vorkommen"
  ergibt `postgresack` 32, `postgressnapshot` 0/122, `receive` 0/187 (32/341 = 9,38 %),
  mit Dedupe 303/341 = 88,86 %. Derselbe Wert 651/813 druckt jeder der drei
  PostgreSQL-17-Läufe (siehe Flake-Untersuchung).
- **Skip-Eigenschaft:** `go test -count=1 -v ./internal/adapters/driven/postgressnapshot`
  im Toolchain-Image mit `--network none` und ohne `CDC_REPLICATION_TEST_DSN`: 21 `--- SKIP`,
  0 `--- PASS` (21 `func Test…` in `snapshot_test.go`).
- **Suchlauf-Feld (Plan §3) nachgefahren:** `::tex[t]`-Suchlauf: 9 Zeilen in 4 Dateien am
  Vorgänger `5e5d0b89`, 21 Zeilen in 4 Dateien an `bf05067d` und `HEAD` (`ADR-0111` 2,
  `ADR-0115` 13, `done/`-Record 1, Erst-Review 5) — beide Zahlen stimmen. Namensuchlauf
  (`git grep -c -E 'postgresack|replication/receive|postgressnapshot'` mit den Pfad-Ausschlüssen
  des Feldes): Parent `d7539e2c` und Kopf je Datei wie im Feld genannt (`coverage-gate.md`
  16 → 20, `db-adapter-coverage.md` 11 → 14, `db-coverage.sh` 7 → 4, `run-replication-tests.sh`
  7 → 8, `kern-rename` 1 → 2); kein Treffer in `.github/workflows`, `Makefile`, `harness/mk`.
- **Kommentare (§3.7):** hinzugefügte Zeilen in `*.go`, `*.sh`, `*.mk` per Textsuche auf
  `slice-`/`welle-`/„seit "/„jetzt"/„vorher"/„früher"/„nur noch"/„bisher"/„nicht mehr"/
  „wäre"/„würde"/„statt"/„Review"/„Fixrunde"/`F-<n>`: kein Treffer. In den Markdown-Trägern siehe F-5.
  Kein `//nolint`, keine host-lokalen absoluten Pfade im Diff, keine Konstruktion außerhalb
  Docker (§3.1: Host-Werkzeuge im Skript sind `bash`/`git`/`awk`/`sed`, wie in den
  Nachbar-Skripten).
- **Commit-Struktur:** vier Commits ohne Renames; Betreffs tragen `LH-FA-CAP-009`/`ADR-*`,
  keine `SPEC-`/`ARC-`-Kennung (`make commit-traceability` im Lauf von `make gates` grün).

### Mutationen (Eingabeseite, selbst ausgeführt)

Datei nach jeder Mutation per `git checkout` zurückgenommen; Endstand `git diff` leer
(`tools/schema/plan.yaml` als bekannter Tier-Nebeneffekt jeweils zurückgenommen).

| # | Mutation | Ort | Ergebnis |
|---|---|---|---|
| N1 | `::text` an jeden Spaltennamen der Cursor-Anweisung angehängt | `snapshotlogic/logic.go` (`CursorStatement`) | rot netzlos: `TestCursorStatementSelectsOnlyQuotedNames`, `TestCursorStatementQuotesEveryIdentifier`; rot im Tier (PostgreSQL 17): `TestImageParityWalAndBackfill` mit 33 Meldungen, je Spalte und Typ — `bo` boolean „t" ≠ „true", `ch5` char(5) „ab   " ≠ „ab", `ch1` char(1) „ " ≠ „", `ia` inet „1.2.3.4" ≠ „1.2.3.4/32", `xm` xml, `db` (Domain über bool); dieselben sechs Typen wie in `ADR-0115` |
| N2 | Binär-Ergebnisformat angefordert (`ExecParams` mit `resultFormats = {1}` statt `Exec`) | `snapshot.go` (`NextBlock`) | netzlos **grün** (nicht erreichbar, wie `ADR-0115` §Fitness Function es führt); rot im Tier: `TestColumnsSkipDroppedAndGenerated`, `TestImageParityWalAndBackfill` („Zeile \x00\x00\x00\x01 fehlt im WAL"; Bild leer) — die Meldung nennt hier keine Spalte und keinen Typ, weil schon der Schlüssel `id` binär ist |
| N3 | Rückfall-Klasse bei beendeter Verbindung immer `storage` (`IsClosed`-Zweig gestrichen) | `snapshot.go:292` | rot im Tier: `TestReadAfterBackendTermination` („zweiter NextBlock auf der beendeten Verbindung: … storage, erwartet Klasse transient") |
| N4 | `57P01`/`57P02`/`57P03` aus `Classify` gestrichen | `snapshotlogic/logic.go:143` | rot netzlos: `TestClassify` (drei Fälle); im Tier bleibt es dort grün (N3-Zweig) |
| N5 | `postgressnapshot` aus `DB_COVERAGE_PKGS` | `db-coverage.sh` | `db-package-lists-check.sh` Exit 1 („Filter-Eintrag 'postgressnapshot' trifft nicht genau ein Paket") |
| N6 | `postgressnapshot` aus dem Dockerfile-Filter | `Dockerfile:101` | Exit 1, Diff der Listen genannt |
| N7 | `postgressnapshot` aus der Paketliste der Replication-Messphase | `run-replication-tests.sh` | Exit 1, Diff genannt |
| N8 | `postgresstorage` in der Store-Liste durch `postgresack` ersetzt | `run-store-tests.sh:138` | Exit 1, Diff genannt |
| N9 | unbekannter Eintrag `foo` im Filter | `Dockerfile:101` | Exit 1 („Filter-Eintrag 'foo' trifft nicht genau ein Paket") |
| N10 | zusätzliches Paket in der Replication-Liste | `run-replication-tests.sh` | Exit 1 (`> ./internal/bootstrap`) |
| N11 | Filter-Muster mit Doppel- statt Einzelquote | `Dockerfile:101` | Exit 1 („nicht genau ein Ausschluss-Muster"): das Skript scheitert bei Formatänderung geschlossen, nicht still |

Gebunden an ihrer Eingabeseite sind: Cast-Freiheit der Anweisung (N1, doppelt), Text-Format
(N2, nur im Tier), Fehlerklasse beendeter Sitzung (N3, N4), die vier namentlichen Listen (N5 bis N10).

### Flake-Untersuchung (`TestStreamRestartsOnExistingSlot`, SQLSTATE 55006)

Volle Läufe (`make test-replication`, Exit ungepiped gesichert, `tools/schema/plan.yaml`
danach zurückgenommen):

| Lauf | Image | Exit | Dauer | Befund |
|---|---|---|---|---|
| 1 | PostgreSQL 17 (Digest aus `e2e.yml`, über `PG_TEST_IMAGE`) | 0 | 76 s | kein Fehlschlag; 80.07 % (651/813); Slot-Reserve `--- PASS` (0,16 s) |
| 2 | PostgreSQL 17 | 0 | 77 s | kein Fehlschlag; 80.07 %; `--- PASS` (0,14 s) |
| 3 | PostgreSQL 17 | 0 | 84 s | kein Fehlschlag; 80.07 %; `--- PASS` (0,30 s) |
| 4 | PostgreSQL 18 (Default) | 2 (make; `go test` Exit 1) | — | **`TestStreamRestartsOnExistingSlot` rot in der Messphase**: „replication slot "slot_pgc_test_restart" is active for PID 180 (SQLSTATE 55006)"; kein DB-Coverage-Verdikt gedruckt |
| 5 | PostgreSQL 18 | 0 | — | kein Fehlschlag; 80.07 %; `--- PASS` |

Isolierte Wiederholung gegen ein Wegwerf-PostgreSQL 17 mit derselben Konfiguration wie
der Tier (`wal_level=logical`, `max_replication_slots=10`, `wal_sender_timeout=2000`), nur
dieser Test, `-count=300`:

| Bedingung | Fehlschläge |
|---|---|
| allein | **0 von 300** |
| parallel dazu Paket `postgressnapshot` (Stand `HEAD`, `-count=6`, Lastfenster 49 s) | **11 von 300** |
| dasselbe, zweiter Lauf (Lastfenster 55 s) | **15 von 300** |
| parallel dazu Paket `postgressnapshot` am Parent-Stand `d7539e2c` (`-count=6`, Lastfenster 28 s) | **4 von 300** |

Jede Fehlermeldung ist dieselbe (`START_REPLICATION: … is active for PID …`, 55006, an
`stream_test.go:765`). Lesart: die Ursache liegt im Test selbst — der Restart legt den
Slot auf, bevor der Walsender des beendeten ersten Laufs ihn serverseitig freigegeben hat;
`receive` kennt keinen Wiederholungsversuch. Sie besteht am Parent-Stand ebenfalls (4 von 300
unter der Last des Parent-Pakets); die neuen Tests verändern die Rate je Überlappung nicht
erkennbar (rund 5 % gegen 8 bis 10 % je Iteration im Lastfenster, kleine Zahlen, keine
Signifikanzprüfung), verlängern aber das Lastfenster gegenüber dem Parent-Stand auf etwa das
Doppelte (Slot-Anlagen, `CREATE DATABASE`, 81-Spalten-Tabelle je GUC-Lage, `pg_terminate_backend`).
Die Ursache im neuen Code: nein; die Häufigkeit im Tier: mitbedingt. Slot-Namen der neuen Tests
tragen Prozess- und Zähler-Suffix (`suffix()`), Kollisionen mit `receive`-Slots sind
ausgeschlossen; kein Test verwendet `t.Parallel`; der 10-Slot-Rahmen des Tier-Containers wurde
in keinem Lauf überschritten (kein `53400` in fünf Tier-Läufen).

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — Typ-Satz ohne `hstore`/`citext`/`ltree`: die Begründung („der Testcontainer bringt sie nicht mit") widerspricht der Messung und der ADR

- `kategorie`: HIGH
- `quelle`: [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)
  Festlegung 4 („Der Paritätstest … trägt den Typ-Satz aus §Kontext als Daten") und
  Fitness Function Zeile 1 („je Spalte des Typ-Satzes"), `AGENTS.md` §3.12 Instanz B
  (Tatsachenbehauptung im Träger trägt ihren Beleg-Anker), §3.7 (ein Kommentar beschreibt,
  was da ist); Skill-Klassen „ADR-Verstoß" und „Beleg trägt seinen Satz nicht"
- `pfad`: `internal/adapters/driven/postgressnapshot/parity_types_test.go:26-29`
  (Godoc `typeTable`), `docs/plan/planning/in-progress/slice-backfill-snapshot-reader.md:107-109`
  (DoD Bild-Parität: „Extension-Typen, die der Testcontainer nicht mitbringt") und `:219`
  (Messbefund: „Nicht anlegbar im Testcontainer sind die Extension-Typen `hstore`, `citext`,
  `ltree`")
- `befund`: Selbst gemessen: `pg_available_extensions` führt `citext`, `hstore` und `ltree`
  in `postgres:17-alpine` (17.11) und `postgres:18-alpine` (18.6), und der Testnutzer
  (`POSTGRES_USER`) ist Superuser; eine Kopie des Paritätstests mit `CREATE EXTENSION` und vier
  zusätzlichen Spalten (`hstore`, `citext`, `ltree`, `hstore[]`; 85 statt 81 Typ-Spalten) läuft
  auf beiden Versionen grün (17.11 und 18.6, `--- PASS`, Roh-Text und Bild byte-gleich). Die
  ADR-Probe hatte genau diese Typen auf denselben Images gemessen (86 Spalten), der Test trägt
  81; die Begründung des Ausschlusses in Godoc, DoD und Messbefund ist damit als gemessen
  dargestellt und trifft nicht zu. Die Folgepflicht 2 der ADR nennt die Extension-Typen in
  ihrer „mindestens"-Liste nicht — Festlegung 4 und die Fitness Function beziehen den Typ-Satz
  aus §Kontext, der sie trägt.
- `verifizierbar`: ja — `psql -Atc "select name from pg_available_extensions where name in
  ('hstore','citext','ltree')"` im Testcontainer; `make test-replication` mit den vier Spalten
- `klasse`: Tatsachenbehauptung im Träger widerspricht der Messung · Typ-Satz der Accepted-ADR
  nur teilweise geliefert

### F-2 — Flake der Tier-Läufe: `TestStreamRestartsOnExistingSlot` (55006) ohne Träger, durch Last mitbedingt, die der Diff vergrößert

- `kategorie`: MEDIUM
- `quelle`: Maintainability; `AGENTS.md` §3.13 (bewegte Eigenschaft: Last auf dem geteilten
  Testcontainer), Plan §6 (Risiken), [`ADR-0030`](../plan/adr/0030-testpyramide.md)
- `pfad`: `internal/adapters/driving/replication/receive/stream_test.go:703-765` (nicht im
  Diff), `tools/harness/run-replication-tests.sh:57-62` (ein Container für alle Pakete beider
  Phasen), `harness/sensors/db-adapter-coverage.md` §Grenze Nr. 7, Plan §6
- `befund`: Gemessen (siehe Flake-Untersuchung): 1 von 5 vollen Läufen rot (PostgreSQL 18,
  Messphase, 55006; das DB-Adapter-Verdikt entfällt dann); isoliert 0 von 300 gegen 11, 15
  und 4 von 300 unter der Last des Pakets `postgressnapshot` (Stand `HEAD` bzw. Parent).
  Weder der Plan §6 noch die Sensor-Dokumentation nennen die Fehlerklasse, obwohl der
  Implementer sie im ersten PostgreSQL-17-Lauf beobachtet hat; Grenze Nr. 7 beschreibt nur,
  dass `TestSlotCreationTimeout` fremde Slot-Anlagen verzögert, nicht die Wirkung der Last auf
  den Restart-Test. `e2e.yml` fährt beide Phasen je Matrix-Leg (siehe F-9).
- `verifizierbar`: ja — `go test -count=300 -run '^TestStreamRestartsOnExistingSlot$'` des
  Pakets `replication/receive` mit und ohne parallele Last auf demselben PostgreSQL
- `klasse`: Flake im geteilten Testcontainer ohne Träger (Last-abhängiger Wettlauf an einer
  Walsender-Freigabe)

### F-3 — Beschreibung der Phase `tier` in Skript und Workflow trägt den Slot-Reserve-Lauf nicht; der Suchlauf suchte sie nicht

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Träger der bewegten Eigenschaft nachziehen; Suchlauf nach der
  Eigenschaft, nicht nach dem eigenen Diff), §3.7 (Zusage stimmt)
- `pfad`: `tools/harness/run-replication-tests.sh:111-112` („Phase `tier` — der Tier-weite
  `go test ./...`; der Exit dieses Aufrufs ist das Verdikt dieser Phase" — die Phase führt jetzt
  zwei Aufrufe mit je eigenem Exit), `.github/workflows/e2e.yml:109` und `:122` (Schrittname
  „Replication-Tier (go test ./...)"); Plan §3 Suchlauf-Feld „Träger der Slot-Reserve-Zusage"
- `befund`: Die Kopfzeile des Skripts (Zeile 21) trägt die zweite Aufgabe der Phase, der
  Kommentar an der Phase selbst und der Workflow-Kommentar beschreiben `tier` weiter als
  den einen Aufruf; der Suchlauf des Plans greift nur die Variablen- und Testnamen
  (`SNAPSHOT_TEST_EXCLUSIVE_DSN|TestSlotReserve`), nicht die Bezeichnung der Phase.
- `verifizierbar`: ja — `git grep -n -i "tier-weite\|Replication-Tier"`
- `klasse`: bewegte Eigenschaft ohne Nachzug in zwei Trägern (Suchform trifft nur Symbolnamen)

### F-4 — Test-Hygiene: eine Rolle je Lauf bleibt zurück; zwei Tests lassen bei `t.Fatalf` die Replication-Verbindung offen

- `kategorie`: LOW
- `quelle`: Maintainability (Aufräumen im Testcontainer)
- `pfad`: `snapshot_test.go:1035` (Cleanup `DROP DATABASE … WITH (FORCE)`) gegen `:1044-1047`
  (`newRole`, Cleanup `DROP ROLE`); `:234-241` und `:274-281` (`exportSnapshot` ohne
  Cleanup der `export.conn` vor dem `importSnapshot`)
- `befund`: Gemessen im Wegwerf-Container: `TestCatalogQueryFailuresKeepTheirClass` lässt je
  Lauf eine Rolle `snap_role_…` zurück (Zählung 0 → 1 → 2 über zwei Läufe; die übrigen Tests
  hinterlassen weder Rolle noch Tabelle, Typ, Publication, Schema, Datenbank oder Slot):
  die Rolle trägt einen Grant in der eigenen Datenbank und wird vor ihr gedroppt (Cleanup-Reihenfolge
  LIFO), `bestEffort` verschluckt den Fehler. Aus dem Code gelesen (nicht ausgelöst): ein
  `t.Fatalf` zwischen `exportSnapshot` und `importSnapshot` in den zwei genannten Tests lässt die
  Replication-Verbindung samt temporärem Slot bis zum Prozessende offen.
- `verifizierbar`: ja — `select count(*) from pg_roles where rolname like 'snap_role%'` nach
  dem Paketlauf
- `klasse`: Aufräum-Reihenfolge im Test (Cleanup ohne Wirkung)

### F-5 — Doku-Prosa: Vorher-Nachher-Satz in `coverage-gate.md`, Vorbedingungs-Satz in `db-adapter-coverage.md` gegen den Test

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Doku-Prosa beschreibt, was da ist), §3.12
- `pfad`: `harness/sensors/coverage-gate.md:69` („Der Stand davor misst im selben Verfahren
  **1691 von 2040**…"), `harness/sensors/db-adapter-coverage.md:207` („jeder Test eines
  ausgenommenen Pakets überspringt ohne Datenbank (`testDSN(t)` als erste Anweisung)")
- `befund`: Der erste Satz benennt den Vorgänger-Stand als Zustand einer Vergleichsmessung in
  einem Sensor-Dokument, dessen Zahlen sonst Lauf-Ursprünge tragen (Messwert richtig, gemessen
  1691/2040; die Form ist Chronik). Der zweite nennt `testDSN(t)` als erste Anweisung jedes
  Tests; `TestSlotReserveExhaustedIsConfiguration` beginnt mit einer eigenen `os.Getenv`-Prüfung
  der zweiten Variablen (`snapshot_test.go:828-832`), nicht mit `testDSN(t)`.
- `verifizierbar`: nein
- `klasse`: Chronik-Prosa in Sensor-Doku · Beleg-Satz trägt die Ausnahme nicht ganz

### F-6 — Fehlerklassen `57P01`/`57P02`/`57P03` und beendete Verbindung als `transient`: durch Port und Plan gedeckt, kein weiterer Aufrufer

- `kategorie`: INFO
- `quelle`: [`SPEC-008`](../../spec/pflichtenheft.md) (`transient`: vorübergehend nicht
  verfügbare Quelle oder Speicher), [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)
  (`transient` → Backoff)
- `pfad`: `snapshotlogic/logic.go:133-145`, `snapshot.go:283-296`,
  `internal/application/port/outbound/tablesnapshot.go:22-25`
- `befund`: Der Port führte „die Verbindung ist abgebrochen" schon in der Erstlieferung unter
  `transient`; die Zuordnung `storage` bei `57P01` widersprach ihm, die Fixrunde stellt Code und
  Port-Doku gleich. Der Plan §3 („Fehlerklasse einer beendeten Sitzung") trägt die Messung
  und die Zuordnung; eine ADR ist nicht nötig, weil der Port die Klasse bereits benennt.
  Verbraucher der `ErrSnapshot*`-Sentinels außerhalb des Pakets und seiner Tests: keine
  (`git grep`), die Änderung hat keine Wirkung auf bestehende Aufrufer. Aus der Erst-Review-Reihe
  bleibt unverändert: `53300` als `configuration` ist im Plan nicht begründet (ausgeschöpfte
  `max_connections` ist ein vorübergehender Zustand).
- `verifizierbar`: ja — N3/N4
- `klasse`: —

### F-7 — Zahlen der Fixrunde nachgemessen: sie stimmen; die Lauf-Bindung der gedeckten Zahl spannt jetzt 83,1 bis 83,3 %

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12, `harness/sensors/coverage-gate.md` §Zählbasis
- `pfad`: `harness/sensors/coverage-gate.md:62-75`, `harness/sensors/db-adapter-coverage.md`
  §Zählbasis, Plan §3 Messbefunde
- `befund`: Nachgemessen und übereinstimmend: Nenner 2040 → 2082, `snapshotlogic` 42/42, 813 mit
  den vier Anteilen 472 · 32 · 122 · 187, gedeckt 348 · 32 · 118 · 153, 651/813 (80,07 %) auf
  PostgreSQL 17 (dreimal) und 18, 245 × 3 = 735, 303/341 und 32/341, 21 von 21 `SKIP`. Die
  gedeckte Unit-Zahl ist lauf-gebunden, wie die Sensor-Doku es führt (Implementer: 1731 und
  Ausdruck `83.1%`/`83.2%`; hier 1733, `make coverage-gate` druckt zweimal `83.30%`); die
  Rücknahme-Zahl „78.60%" hier 78,68 % (dedupliziert, gedruckt `78.7%`), dasselbe Verhalten.
- `verifizierbar`: ja
- `klasse`: —

### F-8 — Speicher eines Blocks: die Grenze `B` zählt Zeilen, nicht Bytes; der Block liegt kurzzeitig doppelt vor

- `kategorie`: INFO
- `quelle`: [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md), Port-Doku
  `outbound/tablesnapshot.go:77-79`
- `pfad`: `snapshot.go:287-301` (`ReadAll` über den einfachen Query-Pfad, dann
  `snapshotlogic.Values` je Zeile)
- `befund`: Aus dem Code und dem Treiber-Quelltext (`pgconn` v5.11.0, `MultiResultReader.ReadAll`
  kopiert die Zeilen) abgeleitet, **nicht gemessen**: ein `NextBlock` hält bis zur Rückgabe die
  Zeilenbytes des Treibers und die Zeichenketten der Werte, also etwa das Doppelte der
  Blocknutzdaten; bei Zeilen im MB-Bereich (`jsonb`, `bytea`) ist der Bedarf `B` × Zeilenbreite,
  nicht durch `B` allein begrenzt. Die Blockgröße ist ein als solcher benannter Startwert. Keine
  erwartete Aktion in diesem Diff; Hinweis für die Ausbaustufe.
- `verifizierbar`: nein
- `klasse`: —

### F-9 — Der Tier läuft in `e2e.yml`; der Plan nennt keine Workflow-Wirkung, und drei Zeitannahmen der neuen Phase sind nur auf dem lokalen Rechner gemessen

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.10, §3.13; Plan §3 Suchlauf-Feld („keine Workflow-Änderung")
- `pfad`: `.github/workflows/e2e.yml:120-123` (`run-replication-tests.sh measure` und `tier`,
  je Matrix-Leg PostgreSQL 17 und 18), `tools/harness/run-replication-tests.sh:132-155`
- `befund`: Die Datei `e2e.yml` ist unverändert, ihr Schritt `tier` fährt aber seither einen
  zweiten Container (`max_replication_slots=1`, `wait_ready` über `pg_isready` im Container,
  dann `go test -v`); §3.10 verlangt den Post-Push-Lauf für neue oder strukturell geänderte
  Workflows — hier ändert sich das Verhalten eines bestehenden Schritts, der Plan führt die
  Wirkung nicht als Risiko. Gemessen: Container-Start bis `pg_isready` rund 4 s, Test 0,13 bis
  0,30 s. Nicht ausgelöst, aus dem Code gelesen: (a) `pg_isready` im Container antwortet schon
  während der Init-Phase des Images über den Unix-Socket, das Bestandsmuster des Hauptcontainers
  trägt dasselbe; (b) der WAL-Zweig des Paritätstests liest ohne Keepalive-Antwort und braucht
  je GUC-Lage rund 0,4 bis 0,6 s gegen `wal_sender_timeout=2000`; (c) `e2e.yml` ist nicht
  blockierend.
- `verifizierbar`: nein (erst ein Lauf auf dem Runner)
- `klasse`: —

### F-10 — Fremde offene Pläne in der Fixrunde nachgezogen statt gemeldet

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 („ein Träger, der eine fremde Datei betrifft, wird gemeldet statt
  still mitgeändert")
- `pfad`: `docs/plan/planning/open/slice-transformationen-kern-rename.md`,
  `docs/plan/planning/open/slice-transformationen-antragsweg-usecase.md` (je eine Aufzählung um
  `postgressnapshot` ergänzt; Plan §3 nennt beide Änderungen)
- `befund`: Die Änderung ist im committeten Suchlauf-Feld benannt und minimal (nicht still); sie
  liegt trotzdem an Plänen anderer, noch nicht gestarteter Slices, deren Planner die Aufzählung
  beim Start prüft. Kein inhaltlicher Fehler.
- `verifizierbar`: ja — `git diff d7539e2c..HEAD -- docs/plan/planning/open/`
- `klasse`: —

## Negativbefunde

- geprüft, ohne Befund: **Umsetzung von [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) im Lesepfad** — `CursorStatement` wählt nur
  gequotete Namen (kein `::`, kein `(`; Unit-Test N1 rot, Namens- und Schema-Quoting an
  `"Sch ""ema x"`/`"Mixed Case ""y"` gebunden); `NextBlock` liest über den einfachen
  `pgconn`-Query-Pfad (`Exec`), dessen Ergebnisformat Text ist — Binär-Anforderung färbt den Tier
  rot (N2); `Values` übernimmt die Rohbytes ohne Trimmen oder Umkodierung (`char(5)` „ab   ",
  Backslashes in `text` und `bytea` unter `escape`, Unicode, leerer Text gegen NULL:
  `TestValues`, Paritätstest); der Adapter setzt keine Sitzungs-GUC (`SET`/`options`), nur
  `BEGIN … REPEATABLE READ READ ONLY` und `SET TRANSACTION SNAPSHOT` mit geprüftem Namen;
  `client_encoding` setzt `pgconn` nicht (Treiber-Quelltext v5.11.0), beide Pfade laufen mit
  demselben DSN. Ungültiges UTF-8 (Datenbank mit anderer Kodierung) ist nicht getestet; die Bytes
  gehen unverändert in den Go-String.
- geprüft, ohne Befund: **Paritätstest** — 81 Typ-Spalten (gemeldet in der Testausgabe, PostgreSQL
  17.11 und 18.6), drei GUC-Lagen; der Vergleich läuft gegen das reale WAL-Bild (Publication, Slot,
  Walsender der Rolle: gedruckt „Bild Zeile 1 in der Lage Tokyo … `24.09.2026 17:11:12.5 JST`,
  `+1-2 +3 +4:05:06.7`, `bytea` im Escape-Format), nicht gegen eine aus derselben Funktion abgeleitete
  Erwartung; die Meldung nennt Spalte und Typ (N1: 33 Zeilen); NULL-Zeile, Rand-Zeile, gelöschte und
  generierte Spalte (`STORED`, unter 18 `VIRTUAL`) belegt; Typ-Lücke siehe F-1.
- geprüft, ohne Befund: **Unterpaket `snapshotlogic`** — liegt im Gate-Nenner (2040 → 2082, 42/42);
  Imports beschränkt auf `pgconn` und den Outbound Port, `make a-check` „gesamt: 0 Befund(e)"; Tests
  (`Validate` inkl. Grenzfälle, `SlotName` an 56/57 Zeichen, `ValidSnapshotName` mit Injection-Form,
  `QuoteIdent`, `CursorStatement`, `Estimate` inkl. `0` und `-0.5`, `Classify` mit 14 Fällen) an
  ihrer Eingabeseite gebunden (N1, N4); nicht gebunden bleibt `pgconn.Timeout(cause)` als
  eigener Zweig von `Classify` (kein Fall mit reinem Timeout-Fehler ohne Kontext-Ende).
- geprüft, ohne Befund: **Listen-Prüfskript `db-package-lists-check.sh` und `coverage.mk`** —
  Extraktion aus dem Filter, `--coverpkg`, Store- und Replication-Liste; alle Mutationen N5 bis N11
  enden mit Exit 1 und nennen die abweichenden Pakete; Formatänderung des Filters scheitert
  geschlossen; keine Ausnahme, kein Überspringen; die erste Rezeptzeile von `coverage-gate` läuft
  vor `docker build`, im Lauf von `make gates` (Ausgabe: `db-package-lists-check: OK`); netzlos,
  read-only. §3.6: keine Schwelle gesenkt (`THRESHOLD`, `DB_COVERAGE_THRESHOLD` unverändert), der
  Diff verschärft (Nenner der Stufe wächst, Listengleichheit gewächtert).
- geprüft, ohne Befund: **Tier-Änderungen** — `trap cleanup EXIT` räumt beide Container und das Netz
  in jedem Ausgang (nach den Läufen keine Container und kein Netz `cdc-repl-test` übrig);
  `wait_ready` bricht bei Nicht-Bereitschaft ab; der Lauf ist rot ohne `--- PASS` des Tests;
  Laufzeitkosten der Phase rund 4 s (Container) plus unter 1 s (Test) von 76 bis 84 s Gesamtdauer;
  die Zusatzwerte `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN` und `max_replication_slots=1` stehen nur in der
  Phase `tier`, die Messphase läuft ohne sie (kein Einfluss auf die DB-Adapter-Zahl).
- geprüft, ohne Befund: **Sensor-Doku-Zahlen und Grenze Nr. 8** (`db-adapter-coverage.md`,
  `coverage-gate.md`, `harness/README.md`): Zahlen mit Lauf-Ursprung, alle nachgemessen (F-7); die
  Aussagen zu N5 bis N10 stimmen mit meinen Läufen; Grenze Nr. 8 trennt Listengleichheit
  (gewächtert) von Eigenschaft und Skip-Eigenschaft (Disziplin) sachlich richtig; `harness/README.md`
  verweist statt eine vierte Paketliste zu führen.
- geprüft, ohne Befund: **Port-Doku** (`outbound/tablesnapshot.go`) — die Beschreibung der
  Zeilenwerte nennt die Ausgabefunktion ohne Cast als Erzeugung (Zusage, `ADR-0115`); §3.7-Klassen
  gewahrt.
- geprüft, ohne Befund: **Commit-Struktur, Traceability, Docker-only, Suppression, host-lokale
  Pfade** — siehe Prüfungen oben; `docs/user/` unberührt (keine Betreiber-Oberfläche:
  `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN` ist eine Testvariable).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Tatsachenbehauptung im Träger widerspricht der Messung (F-1) ·
Flake im geteilten Testcontainer ohne Träger (F-2) · bewegte Eigenschaft ohne Nachzug, Suchform trifft
nur Symbolnamen (F-3) · Aufräum-Reihenfolge im Test (F-4) · Chronik-Prosa in Sensor-Doku (F-5)

## Verdikt

**Merge-blockierend:** ja — ein HIGH (F-1: die Ausschluss-Begründung der Extension-Typen ist als
gemessen dargestellt und widerlegt; der Typ-Satz der Accepted-ADR ist um vier Spalten unvollständig,
die Zusage selbst hält an ihnen, gemessen auf 17.11 und 18.6) und ein MEDIUM (F-2: ein Flake im Tier,
der in einem von fünf vollen Läufen rot war und keinen Träger führt). Die Umsetzung von
[`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) im Lesepfad ist
konform und an ihrer Eingabeseite gebunden (N1, N2); das Unterpaket, die Klassifikation, das
Listen-Prüfskript und die Sensor-Zahlen tragen; die Gates am Stand laufen grün (`make test`,
`make a-check`, `make coverage-gate`, `make gates`). Die Fixrunde lockert keine Gate-Schwelle.

**Übergabe:** F-1 bis F-5 gehen an den Implementer; über den Ort der F-2-Behebung (Test, Tier-Aufbau,
Beobachtungs-Eintrag) entscheidet er mit dem Planner, der Report ist das Übergabe-Artefakt. F-6 bis F-10
sind ohne erwartete Aktion in diesem Diff. Die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" ist im Plan bereits `[x]` (erster Review); dieser Report legt keinen
weiteren Haken an, weil eine weitere Fixrunde nötig ist. Die **Finding-Klassen** gehen in die
Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg**; DoD- und
Spec-Konformität prüft der Verifier separat (Modul 11).
