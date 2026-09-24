# Review-Report: slice-backfill-snapshot-reader — 2026-09-24

**Review-Art:** Code — der Diff führt den Outbound Port `TableSnapshotPort`, den
Driven Adapter `postgressnapshot` (temporärer Slot mit Export-Snapshot,
Cursor-Blöcke, Katalog-Schätzung) und die Gate-Zuordnung des neuen Pakets ein;
geprüft gegen Plan, ADRs und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-backfill-snapshot-reader`, Diff-Range
`62d00abe..HEAD` (`6cf9bafa`). Slice-Commits `7673e774`, `1ae3fe3f`, `48aa388a`
(Lifecycle, Verantwortlich), `89b21079` (Port, Adapter, Tests), `d3b47b7e`
(Coverage-Ausschluss, `DB_COVERAGE_PKGS`, Messphase), `9a2efe6a` (Testausgabe
`reltuples`/Version), `641b2ba4` (Sensor-Beschreibungen, Tier-Zeile), `63351d4c`,
`6cf9bafa` (Plan, Messbefunde).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um weitere HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz,
Zusage-ohne-Eingabeseite, Kommentar-Chronik, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-24.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-snapshot-reader` (§1 Ziel, §2 DoD als Prüfmaßstab
  für Plan-Zusagen, §3 Plan, Suchlauf-Feld und Messbefunde, §6 Risiken) und
  Welle `welle-backfill-bestand`
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 1 (Mechanismus, Ablauf), Teilfrage 2 (Row-Image-Parität),
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3 (`reltuples`),
  [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 1/3 und Trigger (a),
  [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
  [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md),
  [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md),
  [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md),
  [`ADR-0030`](../plan/adr/0030-testpyramide.md)
- [`LH-FA-CAP-009`](../../spec/lastenheft.md),
  [`LH-FA-CAP-008`](../../spec/lastenheft.md),
  [`SPEC-008`](../../spec/pflichtenheft.md) (Fehlerklassen),
  [`SPEC-012`](../../spec/pflichtenheft.md) (PostgreSQL 17 und 18)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.6, §3.7, §3.11, §3.12, §3.13),
  `harness/conventions.md` (MR-000/MR-001)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-change-origin.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem
Implementer-Bericht übernommen; Exit-Codes ungepiped gesichert):

- **Gates am Stand `6cf9bafa`:** `make test` (Race-Detector) Exit 0 —
  `postgressnapshot` „ok" (alle Tests `SKIP`, siehe unten); `make a-check`
  Exit 0 („gesamt: 0 Befund(e)"); `make coverage-gate` Exit 0 (drei Läufe
  gedruckt „Coverage 82.80% erfüllt Schwelle 80%"; der Lauf innerhalb von `make gates`
  druckt „82.90%", siehe F-14); `make test-store`
  Exit 0; `make test-replication` (PostgreSQL 18, zweimal) Exit 0, gedruckt
  „DB-Adapter-Coverage: 79.25% (gedeckt 672 von 848 Statements; Profile gemergt:
  store,replication)"; `PG_TEST_IMAGE=<17er-Digest aus e2e.yml> make
  test-replication` Exit 0, dieselbe Zeile.
- **Test-Ausschluss:** `go test -count=1 -v ./internal/adapters/driven/postgressnapshot`
  im gepinnten Toolchain-Image mit `--network none` und ohne
  `CDC_REPLICATION_TEST_DSN`: 13 `func Test…` in `snapshot_test.go`, 13 Zeilen
  `--- SKIP`, 0 `--- PASS`. Jeder Test ruft `testDSN(t)` bzw. prüft
  `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN` als erste Anweisung.
- **Zahlen der DB-Adapter-Doku nachgemessen** (Profil aus dem Lauf oben,
  dedupliziert über die Block-Position): 271 Positionen × 3 = 813 Zeilen im
  Replication-Profil; Anteile 472 · 32 · 157 · 187 = 848; `postgressnapshot`
  139 von 157, `replication/receive` 153 von 187, `postgresack` 32 von 32;
  „nur das erste Vorkommen" ergibt 32 von 376 = 8,51 % — alle Angaben in
  `harness/sensors/db-adapter-coverage.md` stimmen.
- **Unit-Nenner:** die `coverage`-Stufe misst am Parent (`48aa388a`,
  `git archive` in ein Wegwerfverzeichnis) und am Kopf je **2040** Statements,
  gedeckt **1689** (82,80 %) — der Diff bewegt den Nenner nicht (siehe F-14 zur
  Streuung der gedeckten Zahl).
- **Suchlauf-Feld (Plan §3) mit `git grep -c` an den Ständen `48aa388a`,
  `641b2ba4`, `HEAD` nachgefahren:** die Zahlen je Datei stimmen; einzige
  Abweichung: `docs/plan/planning/welle-backfill-bestand.md` trägt den Treffer
  bereits am **Parent** (die Parent-Zeile des Feldes führt die Datei nicht).
  Keine weitere Datei mit einer namentlichen Paketliste außerhalb der genannten
  (`.github/workflows`, `Makefile`, `harness/mk`: kein Treffer).
- **Kommentare (§3.7):** hinzugefügte Zeilen in `*.go`, `Dockerfile`, `*.sh`,
  `*.md` per Textsuche auf `slice-`/`welle-`/„seit "/„jetzt"/„vorher"/„früher"/
  „nur noch"/„bisher"/„nicht mehr"/„wäre"/„würde" durchsucht: keine Chronik in
  Produktionscode-Kommentaren; `slice-…`-Treffer in `harness/sensors/*.md` sind
  Lauf-Kennungen als Ursprung von Zahlen (§3.12, zulässig); die Konjunktive
  stehen in Test-Godoc und Sensor-Grenze (F-8). Kein `//nolint`; keine
  host-lokalen absoluten Pfade im Diff (`git diff` nach den Wurzel-Segmenten der
  Präfixliste des `hostpaths`-Moduls: leer).
- **Commit-Struktur:** `git show -M --stat` über `7673e774` und `48aa388a`: reine
  Renames ohne Zeilenänderung; Inhalt (`1ae3fe3f`) getrennt — §3.3 erfüllt.
  Alle Betreffs tragen `LH-FA-CAP-009`/`ADR-*`, keine `SPEC-`/`ARC-`-Kennung;
  `docs/user/` ohne Diff, keine Betreiber-Oberfläche im Diff
  (`CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN` ist eine Testvariable, keine
  Feed-Container-Variable); `git status`: sauber, nur ignorierte Bestandsdateien.
- **Typ-Sonde zu F-1 (PostgreSQL 18, Wegwerf-Container):** 32 Spalten
  (`bool`, `char(n)`, `varchar`, `inet`, `cidr`, `macaddr`, `numeric`, `uuid`,
  `timestamp(tz)`, `date`, `time(tz)`, `interval`, `bit`, `varbit`, `money`, `oid`,
  `tsvector`, `xml`, `"char"`, `name`, Enum, Range, `bool[]`, `point`, `json`,
  `jsonb`, `float4`, `float8`, `bytea`, `text[]`) je als `col` und als `col::text`
  in `psql -At` gelesen: die beiden Ausgaben unterscheiden sich für drei Typen.

### Mutationen (Eingabeseite, selbst ausgeführt)

Datei nach jeder Mutation per `git checkout` zurückgenommen; Endstand `git diff`
leer (`tools/schema/plan.yaml` als bekannter Tier-Nebeneffekt jeweils
zurückgenommen).

| # | Mutation | Ort | Ergebnis |
|---|---|---|---|
| R1 | `SET TRANSACTION SNAPSHOT` gestrichen | `snapshot.go` (`importSnapshot`) | rot: `TestSnapshotPairedWithPoint` („[1|a 2|b 3|c]"); `TestTemporarySlotEndsWithConnection` bleibt grün |
| R2 | `FETCH FORWARD B+1` | `snapshot.go` (`NextBlock`) | rot: `TestBlocks` (Unterfälle „mehrere Blöcke", „Blockgrenze": [4 3 0], [4 2 0]) |
| R3 | Filter `attgenerated = ''` gestrichen | `snapshot.go` (`readColumns`) | rot: `TestColumnsSkipDroppedAndGenerated`, `TestImageParityWalAndBackfill` |
| R4 | Paritäts-Fixture um `bool`, `inet`, `char(4)` erweitert | `snapshot_test.go` (`parityColumns`, `parityInsert`) | **rot**: WAL-Bild `"bo":"t"`, `"ip":"1.2.3.4"`, `"ch":"ab  "` ≠ Backfill-Bild `"bo":"true"`, `"ip":"1.2.3.4/32"`, `"ch":"ab"` — siehe F-1 |
| R5 | Anmelde-Klasse `28…` aus `classify` gestrichen | `snapshot.go` | rot: `TestWrongPasswordIsPermission` (Klasse `transient`) |
| R6 | Zeitlimit der Slot-Anlage ignoriert (80 s) | `snapshot.go` (`exportSnapshot`) | rot: `TestSlotCreationTimeout` (20 s) |
| R7 | `EstimatedRows`: `value < 0` → `value < 1` | `snapshot.go` | **grün** — siehe F-9 |
| R8 | Schema-/Tabellen-Quoting im `DECLARE` entfernt | `snapshot.go` (`cursorStatement`) | **grün** — siehe F-4 |
| R9 | `53400` aus dem Klassen-Mapping gestrichen | `snapshot.go` (`classify`) | rot nur im manuellen Lauf gegen einen Wegwerf-PostgreSQL mit `max_replication_slots=1`: „Fehlerklasse replication … (SQLSTATE 53400), erwartet Klasse configuration" — die Lesart des Implementers (Klasse `replication` statt `configuration`) stimmt; im Tier kein Träger |
| R10 | `postgressnapshot` aus dem Dockerfile-Ausschluss | `Dockerfile` (Stufe `coverage`) | rot: `make coverage-gate` „FAIL — Coverage 77.00% unter Schwelle 80%", Exit ≠ 0 |
| R11 | `postgressnapshot` aus `DB_COVERAGE_PKGS` | `tools/harness/db-coverage.sh` | **grün**: „DB-Adapter-Coverage: 77.13% (gedeckt 533 von 691 …)", Schwelle 70 — siehe F-3 |
| R12 | `postgressnapshot` aus der `go test`-Liste der Messphase | `tools/harness/run-replication-tests.sh` | **grün**: dieselbe Zeile 533 von 691 — siehe F-3 |

Gebunden an ihrer Eingabeseite sind: Snapshot-Import (R1), Blockgröße (R2),
Generierte-Spalten-Filter (R3), Anmelde-Klasse (R5), Zeitlimit (R6),
Dockerfile-Ausschluss (R10), `53400`-Mapping nur manuell (R9). Nicht gebunden:
Typ-Parität jenseits der Fixture (R4/F-1), Tabellen-/Schema-Quoting (R8/F-4),
Nullgrenze der Schätzung (R7/F-9), zwei der drei namentlichen Listen (R11/R12/F-3).

**Flakiness `TestSlotCreationTimeout`:** in allen Tier-Läufen dieses Reports
(Basis PostgreSQL 18 ×2, PostgreSQL 17 ×1 und je ein Lauf je Mutation R1–R8,
R11, R12) trat außerhalb der eigenen Mutation R6 kein Fehler dieses Tests auf.
Die Zeitgrenze ist einseitig (Ablauf des 1-s-Limits oder späteres Ende führen
beide zu `transient`; der zweite Aufruf trägt 30 s); die offene
Schreibtransaktion (rund eine Sekunde) kann fremde Slot-Anlagen der parallelen
Pakete nur verzögern, nicht scheitern lassen, solange deren Limit größer ist —
für `receive` und `postgresack` nicht gelesen.

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — Bild-Parität bricht für `bool`, `inet` und `char(n)`; die Paritäts-Fixture enthält keinen der drei Typen

- `kategorie`: HIGH
- `quelle`: [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 2 („ein Backfill-Image darf sich von einem WAL-Image derselben Zeile
  nicht unterscheiden", Konsequenz-Abschnitt gleichlautend),
  [`LH-FA-CAP-008`](../../spec/lastenheft.md); Skill-HIGH-Klasse „Zusage ohne
  Bindung an ihre Eingabeseite"
- `pfad`: `internal/adapters/driven/postgressnapshot/snapshot.go:255-264`
  (`cursorStatement`, `::text` in Zeile 260), `snapshot_test.go:579-587`
  (`parityColumns`/`parityInsert`), `internal/application/port/outbound/tablesnapshot.go:71-73`
  („Text-Stand der Quelle")
- `befund`: Selbst gemessen (R4, PostgreSQL 18.6): dieselbe Zeile ergibt im
  WAL-Bild `"bo":"t"`, `"ip":"1.2.3.4"`, `"ch":"ab  "` und im Backfill-Bild
  `"bo":"true"`, `"ip":"1.2.3.4/32"`, `"ch":"ab"` — der Cast `col::text` ruft für
  diese drei Typen eine eigene Cast-Funktion statt der Ausgabefunktion auf, die
  pgoutput sendet (Typ-Sonde: die übrigen 29 gelesenen Typen stimmen überein). Der
  Paritätstest M5 läuft grün, weil seine Fixture (`timestamptz`, `float8`, `bytea`,
  `interval`, `numeric`, `jsonb`, Array, `date`) keinen dieser Typen trägt; die
  Zusage „byte-gleich" ist an ihrer Eingabeseite (Spaltentyp) nur für die acht
  Fixture-Typen gebunden. `boolean` ist als Typ in Quelltabellen üblich; ein
  Consumer, der das Log als Upsert der Row Images anwendet, liest je nach Pfad
  `"t"` oder `"true"`. Die Formulierung `col::text` steht wörtlich in
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 1 Ablauf Nr. 3 und Teilfrage 2 (Accepted); der Diff setzt sie um.
- `verifizierbar`: ja — `make test-replication` mit einer um `bool`/`inet`/`char(n)`
  erweiterten Fixture (R4)
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite (Typ-Parität) · Accepted-ADR-Text
  trägt den Defekt (Weg: Architect, nicht Implementer-Alleingang)

### F-2 — `db-coverage.sh`: „jede Block-Position dort ZWEIMAL" gegen die Messung (dreimal) im selben, vom Diff umgeschriebenen Absatz

- `kategorie`: HIGH
- `quelle`: Skill-HIGH-Klasse „Zahl im Träger … gegen die Messung driftend",
  `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist), §3.12
- `pfad`: `tools/harness/db-coverage.sh:19-23` (Satz „Im Replication-Lauf testet
  `go test` drei Pakete …; jedes der drei Testbinaries instrumentiert alle
  Gegenstands-Pakete, darum erscheint jede Block-Position dort ZWEIMAL — einmal
  mit ihrem count, einmal mit 0")
- `befund`: Der Absatz nennt drei Testbinaries und leitet daraus „zweimal" ab;
  gemessen sind 813 Zeilen für 271 Positionen = **dreimal** (`harness/sensors/db-adapter-coverage.md`
  im selben Diff sagt „dreimal", „271 Positionen × 3"). Auch die Aufteilung „einmal mit
  ihrem count, einmal mit 0" stimmt nicht mehr: je Position trägt nur die Kopie des
  Binaries, das ihr Paket testet, den Zähler; die zwei anderen tragen `0`.
- `verifizierbar`: ja — `awk` über `replication.coverprofile` (Zeilen je Position);
  Textvergleich mit der Sensor-Doku
- `klasse`: Zahl im Träger driftet gegen die Messung (Skript-Kommentar)

### F-3 — DoD (a) „die drei namentlichen Stellen sind gezogen": zwei der drei sind nicht an ihre Eingabeseite gebunden; die Sensor-Doku führt die Grenze nicht

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 1 und Trigger (a) (die namentliche Liste nachziehen), Skill-Klasse „Zusage
  ohne Bindung an ihre Eingabeseite", „Beleg trägt seinen Satz nicht"
- `pfad`: `tools/harness/db-coverage.sh:47` (`DB_COVERAGE_PKGS`),
  `tools/harness/run-replication-tests.sh:99-104` (Paketliste der Messphase),
  `harness/sensors/db-adapter-coverage.md:37-39` („Der **einzige Träger** der
  Gegenstandsliste ist `tools/harness/db-coverage.sh --coverpkg`"), Plan §2 DoD
  Gate-Zuordnung (a) („*Zu belegen durch:* der Suchlauf in §3 … und ein grüner
  `make coverage-gate` am Diff")
- `befund`: Selbst gemessen: das Streichen des Pakets aus `DB_COVERAGE_PKGS` (R11)
  oder aus der `go test`-Liste der Messphase (R12) färbt nichts rot — beide Läufe
  enden Exit 0 und drucken „77.13% (gedeckt 533 von 691 …)" statt 79,25 % von 848;
  nur der Dockerfile-Ausschluss (R10) färbt `make coverage-gate` rot (77,00 %). Der
  genannte Beleg „grüner `make coverage-gate`" trägt deshalb nur eine der drei Stellen;
  für die anderen zwei ist die Bindung der Suchlauf (Lesen). Die Sensor-Doku benennt
  für die Dockerfile-Liste eine Disziplin-Grenze (`coverage-gate.md` §Grenze Nr. 4/5),
  für die zwei DB-Listen nicht, und `db-adapter-coverage.md` nennt `db-coverage.sh`
  „einzigen Träger" der Liste, obwohl die Messphase eine zweite, namentliche Liste
  führt (Zwei-Quellen-Nähe ohne benannten Gewinner).
- `verifizierbar`: ja — R11/R12 mit `make test-replication`; ein maschineller
  Vergleich der Listen existiert nicht
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite (namentliche Liste ohne Wächter)

### F-4 — Schema-/Tabellen-Quoting im `DECLARE` an keiner Eingabe gebunden

- `kategorie`: MEDIUM
- `quelle`: Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite" /
  „fehlende Negativtests bei neuem öffentlichem Vertrag" (SQL-Bezeichner-Sicherheit
  am Adapter-Rand, Schema und Tabelle stammen aus dem Betriebszustand)
- `pfad`: `snapshot.go:257-264` (`quoteIdent(schema)`/`quoteIdent(table)` in Zeile
  263), `snapshot_test.go:461-482`, `:99-107` (`newTable` legt kleingeschriebene,
  alphanumerische Namen an)
- `befund`: Selbst gemessen (R8): `quoteIdent` um Schema und Tabelle im `DECLARE`
  entfernt, `make test-replication` grün. Nur der **Spaltenname** wird mit einem
  Sonderzeichen-Namen geprüft (`"we ""ird"`, `"Id"`); Schema und Tabelle tragen in
  keinem Test Groß-/Sonderzeichen. Der Code ist an dieser Stelle richtig (die
  Katalog-Abfragen gehen über `to_regclass(quote_ident(…))`, Werte werden nicht
  gebaut); gebunden ist die Zusage nicht.
- `verifizierbar`: ja — Test mit einer Tabelle `"Mixed Case ""x"` (`make
  test-replication`) bzw. R8
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite (Bezeichner-Quoting)

### F-5 — DoD (c): netzlos prüfbare Logik liegt im Ausschluss-Paket; die genannte Nenner-Messung stützt den Satz nicht

- `kategorie`: MEDIUM
- `quelle`: Plan §2 DoD Gate-Zuordnung (c) („damit der Nenner des blockierenden
  Gates nicht ohne Not aus dem Gate herausfällt"),
  [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 1, Skill-Klasse „Beleg trägt seinen Satz nicht"
- `pfad`: `snapshot.go:89-107` (`New`-Validierung), `:138-144` (`slotName`),
  `:257-269` (`cursorStatement`, `quoteIdent`), `:384-399` (`classify`); Plan §3
  „Stand der Übergabe" (Nenner „unverändert", „Review des Diffs steht aus")
- `befund`: Die vier Funktionen sind reine Logik ohne Datenbankzugriff (Gegenprobe:
  keine ruft `pgconn` außer `New` über `ParseConfig`); sie stehen im Ausschluss-Paket
  und liegen damit außerhalb des blockierenden Gates. Gemessen: 29 Statements
  (Profil-Zeilenbereiche dieser Funktionen), davon 26 im Tier gedeckt und drei in
  `classify` nirgends (`:393`, `:395`, `:397`, siehe F-6). Der Nenner der
  `coverage`-Stufe bleibt 2040 → 2040 — das ist der vom Plan genannte Beleg, er zeigt
  aber, dass **nichts** in den Gate-Nenner gewandert ist, und stützt damit nicht,
  dass keine netzlos prüfbare Logik im Paket liegt. Die im DoD genannten Beispiele
  (Blockbildung, Row-Image-Konstruktion) liegen tatsächlich außerhalb.
- `verifizierbar`: ja — `git diff`, Profil der Stufe vor/nach (Nenner 2040/2040)
- `klasse`: netzlos prüfbare Logik im Ausschluss-Paket (`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`)

### F-6 — Fehlerpfade des Ports ungetestet: Klassen `replication` und `storage`, Fehler mitten im Lesen, Kontext-Abbruch des Cursors

- `kategorie`: MEDIUM
- `quelle`: Skill-Klasse „fehlende Negativtests bei neuem öffentlichem Vertrag",
  [`SPEC-008`](../../spec/pflichtenheft.md) (jede Störung endet als sichtbarer Fehler
  mit einer Klasse), [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)
- `pfad`: `snapshot.go:174-180` (`consistent_point`/Snapshot-Name nicht lesbar),
  `:194-200` (Lese-Verbindung), `:207-209` (Import-Fehler mit `abort`), `:343-344`
  (Fehler im `FETCH`), `:393`/`:395`/`:397` (`classify`: `53400`/`53300`,
  Verbindungsfehler, Rückfall-Klasse); `outbound/tablesnapshot.go:28-33`
- `befund`: Aus dem gemergten Profil des Tier-Laufs: 15 Blöcke von `snapshot.go` mit
  `count = 0` (18 der 157 Statements), darunter jeder Zweig, der `ErrSnapshotReplication`
  oder die Rückfall-Klasse erzeugt, der `FETCH`-Fehler und die Aufräumung eines
  fehlgeschlagenen Imports. Kein Test setzt `ErrSnapshotReplication` (kein Slot-Fehler
  außer `53400` im manuellen Lauf) oder `ErrSnapshotStorage` durch einen Lesefehler;
  `ErrSnapshotStorage` ist nur über „NextBlock nach Close" belegt (vom Test selbst
  erzeugter Zustand). Der Abbruch mitten im Lauf ist im Plan §1 ausdrücklich dem
  `e2e`-Slice zugeordnet; die Kontext-Abbruch-Zusage des Adapters (`closeConn` mit
  eigenem Zeitlimit) hat keinen eigenen Test.
- `verifizierbar`: ja — `go tool cover` über `replication.coverprofile`
- `klasse`: fehlende Negativtests bei neuem öffentlichem Vertrag (Fehlerklassen des Ports)

### F-7 — `DefaultBlockSize` trägt keine Herkunftsangabe; `DefaultSlotTimeout` schon

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 (Instanz A: Zahl im Träger trägt ihren Ursprung),
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 1
  Ablauf Nr. 3 („Startwert 1.000, Wert legt der Slice fest")
- `pfad`: `snapshot.go:37-39` (gegen `:41-45`)
- `befund`: `DefaultSlotTimeout` steht als „Startwert (Setzung ohne Messung)", die
  Blockgröße 1000 nur mit ihrer Wirkung („Speicherbedarf … begrenzt"); beide sind
  ungemessene Setzungen.
- `verifizierbar`: nein
- `klasse`: Zahl im Träger ohne Ursprung

### F-8 — Kommentar- und Formfehler in geänderten Zeilen: umgebrochene Absätze in `db-coverage.sh`, Konjunktiv-Begründung im Test-Godoc

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7
- `pfad`: `tools/harness/db-coverage.sh:20` (134 Zeichen: „…; Nenner-Stand: Lauf
  `slice-085`). Im Replication-Lauf testet `go test` drei Pakete (postgresack,
  postgressnapshot,"), `:33` (107 Zeichen, „…Jeder Lauf instrumentiert dabei seinen
  Teil; die beiden"); `snapshot_test.go:740-742` („auf dem gemeinsamen Testcontainer
  würde er die Slots paralleler Pakete blockieren")
- `befund`: Die Teilersetzung der Paketliste hat zwei Absätze mitten in der Zeile
  fortgeführt (der übrige Kommentarblock bricht bei rund 80 Zeichen um); der Test-Godoc
  begründet die Voraussetzung über die verworfene Alternative im Konjunktiv (Test-Godoc,
  kein Produktionscode-Pfad der HIGH-Regel).
- `verifizierbar`: nein
- `klasse`: Kommentar mit Teilersetzungs-Rest

### F-9 — `EstimatedRows`: die Grenze `reltuples = 0` ist ungebunden

- `kategorie`: LOW
- `quelle`: Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite",
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3 (unbekannt heißt `−1`, nie `0`)
- `pfad`: `snapshot.go:303-306`; `snapshot_test.go:488-514`
- `befund`: Selbst gemessen (R7): `value < 0` → `value < 1` bleibt grün; der Test
  prüft `−1` (vor `ANALYZE`) und `40` (danach), nie eine analysierte leere Tabelle
  (`reltuples = 0`, „bekannt, null Zeilen").
- `verifizierbar`: ja — R7 bzw. ein Test mit `ANALYZE` einer leeren Tabelle
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite (Grenzwert)

### F-10 — Manueller Träger der `53400`-Zusage ist aus dem Repo nicht reproduzierbar

- `kategorie`: LOW
- `quelle`: Skill-Klasse „Beleg trägt seinen Satz nicht" (Beleg ohne Befehl),
  `AGENTS.md` §3.12 Instanz B
- `pfad`: `harness/sensors/db-adapter-coverage.md:165-175` (Grenze Nr. 7, „einmaliger,
  manueller Lauf"), `snapshot_test.go:738-769`, Plan §3 „Slot-Reserve"
- `befund`: Die Variable `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN` wird von keinem Skript, Make-Target
  oder Workflow gesetzt (`grep` über den Baum: nur Test, Sensor-Doku, Plan); der Weg zum
  Wegwerf-PostgreSQL (`wal_level=logical`, `max_replication_slots=1`, `max_wal_senders`)
  steht nirgends als Befehl. Ich habe den Lauf nachgebaut (PostgreSQL 18, gepinnter Digest,
  `-c wal_level=logical -c max_replication_slots=1 -c max_wal_senders=10`): grün, und R9 rot;
  die Aussage stimmt, ihr Träger ist ohne Rekonstruktion nicht wiederholbar.
- `verifizierbar`: nein
- `klasse`: Beleg ohne reproduzierbaren Befehl

### F-11 — Test-Ausschluss ohne maschinellen Wächter, in der Sensor-Grenze nicht als Disziplin benannt

- `kategorie`: LOW
- `quelle`: [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 1, Plan §6 (Risiko „Der Test-Ausschluss trägt die Ausnahme nicht")
- `pfad`: `snapshot_test.go:25-30`, `:46-53`;
  `harness/sensors/coverage-gate.md` §Grenze Nr. 4 (letzter Absatz), Nr. 5
- `befund`: Der einmalige `go test -v`-Lauf ohne DSN zeigt 13 von 13 `SKIP` (gemessen);
  ein künftig ergänzter Test ohne `testDSN(t)` liefe netzlos und höbe die DB-Zahl ohne
  DB-Beleg — nichts im Gate oder im Tier fängt das. Die Grenze Nr. 5 nennt die
  Dockerfile-Testpaket-Liste als „Disziplin, kein Sensor", die Skip-Eigenschaft des
  Pakets nicht. Das Risiko ist im Plan benannt, der Slice-Bericht nennt den fehlenden
  Wächter.
- `verifizierbar`: nein
- `klasse`: Ausnahme-Eigenschaft ohne Wächter (Disziplin nicht als Grenze geführt)

### F-12 — `harness/README.md`: neue Paketliste in der Zeile `make test-replication`

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (bewegte Eigenschaft: Träger der Menge nachziehen),
  Maintainability (Zwei-Quellen-Nähe)
- `pfad`: `harness/README.md:139`
- `befund`: Die Datei trug an beiden Ständen keine Paketliste (Plan §3: „Zusatz, kein
  Nachzug"); der Diff führt jetzt „drei Pakete: `postgresack`, `postgressnapshot`,
  `replication/receive`" in einer Tabellenzelle ein — eine vierte namentliche Stelle
  der Liste neben `Dockerfile`, `db-coverage.sh`, `run-replication-tests.sh`, ohne
  Wächter (F-3).
- `verifizierbar`: nein
- `klasse`: zusätzlicher Träger einer namentlichen Liste

### F-13 — `ADR-0080`: kein Verstoß durch das Fehlen der Treiber-Hülle

- `kategorie`: INFO
- `quelle`: [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md),
  [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md), `AGENTS.md` §3.5
- `pfad`: Plan §3 (Zeile `seam.go`, „nicht realisiert"); `snapshot.go:150-182`
- `befund`: [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) Festlegung 1/5 und §Kontext (1) gelten ausdrücklich für die zwei
  Pakete `postgresack` und `receive` („Beide Pakete bekommen …"); für ein neues Paket
  legt sie keine Hülle fest, und [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md) verlangt sie nicht. Es wird keine Accepted-ADR
  überschrieben. Die im Plan genannte Begründung („netzlose Fake-Tests widersprächen dem
  Test-Ausschluss") ist eine Wahl, keine Pflicht: [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) Festlegung 4 lässt netzlose
  Fake-Tests in ausgenommenen Paketen zu (Gegenstand wächst um die Hüllen-Statements) — die
  Wahl hat den in F-5/F-6 beschriebenen Preis. Tests `package postgressnapshot` (intern) sind
  durch M1 (zwei Schritte einzeln) begründet und tragen; kein Befund.
- `verifizierbar`: nein
- `klasse`: —

### F-14 — Streuung der gedeckten Unit-Zahl größer als die in `coverage-gate.md` benannte Ein-Statement-Schwankung

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (Ursprung einer Messung), `harness/sensors/coverage-gate.md`
  §Zählbasis
- `pfad`: `harness/sensors/coverage-gate.md:239-244` („ein zweiter Lauf desselben Stands
  druckte `82.80%`: die in §Zählbasis benannte Schwankung"); Plan §3 Messbefunde
- `befund`: Selbst gemessen an Parent und Kopf, je ein bis drei dedupliziert ausgewertete
  Läufe: **1689 von 2040** (82,80 %) in jedem der vier Läufe; ein fünfter Lauf
  (innerhalb von `make gates`, am Kopf) druckt `82.90%`; der Implementer nennt 1692 (Parent) und 1691 (Diff),
  gedruckt 82,90 %. Die Spanne 1689–1692 beträgt bis zu drei Statements; §Zählbasis führt
  „höchstens 1 Statement" (Stand `slice-093`, Nenner 1903). Der Diff bewegt die Zahl
  nicht (Parent und Kopf messen in meiner Umgebung gleich); die Zuschreibung des
  `82.80%`-Laufs an die benannte Schwankung trägt die Bandbreite nicht.
- `verifizierbar`: ja — `make coverage-gate` mehrfach, dedupliziertes Profil
- `klasse`: Zahl im Träger driftet gegen die Messung (Träger außerhalb der Diff-Hunks)

### F-15 — Vorbestehende Drift „Der aktuelle Nenner ist 1936" bestätigt; der Diff hätte sie nicht ziehen müssen

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12, §3.13
- `pfad`: `harness/sensors/coverage-gate.md:55`
- `befund`: Gemessen: Nenner der `coverage`-Stufe **2040** am Parent (`48aa388a`) und am
  Kopf; der Absatz nennt 1936 (Stand `slice-097`). Der Diff bewegt den Unit-Nenner nicht
  (`postgressnapshot` ist ausgenommen, `outbound/tablesnapshot.go` trägt keine
  zählbaren Statements), die bewegte Eigenschaft des Slice ist die Menge der
  ausgenommenen Pakete; die Zeile steht außerhalb der Diff-Hunks und bleibt INFO.
- `verifizierbar`: ja — Profil der `coverage`-Stufe
- `klasse`: Zahl im Träger driftet gegen die Messung (vorbestehend)

### F-16 — Offene Pläne nennen die Ausschlussliste mit drei Paketen — gemeldet, Nachzug beim Planner

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Träger in fremder Datei wird gemeldet)
- `pfad`: `docs/plan/planning/open/slice-transformationen-kern-rename.md:119`,
  `docs/plan/planning/open/slice-transformationen-antragsweg-usecase.md:102`
- `befund`: Beide DoD-Sätze („nur `postgresstorage`, `postgresack` und `replication/receive`")
  sind am Kopf um `postgressnapshot` unvollständig; die Meldung steht im committeten
  Suchlauf-Feld (Plan §3), der Übergabe-Artefakt-Anforderung genügt sie. Weitere Träger
  der Menge außerhalb `internal/`, `docs/reviews`, `done/`, ADRs: keine (`git grep` an
  `48aa388a`/`HEAD`, Zahlen je Datei stimmen mit dem Feld überein; einzige Abweichung
  `welle-backfill-bestand.md`, am Parent bereits mit einem Treffer).
- `verifizierbar`: ja — `git grep`
- `klasse`: —

### F-17 — Fenster zwischen Snapshot-Export und Cursor ohne Tabellensperre (nicht gemessen)

- `kategorie`: INFO
- `quelle`: [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 1 (Konsequenzen benennen das Halten des Snapshots, nicht DDL im Fenster)
- `pfad`: `snapshot.go:189-226`
- `befund`: Spaltenliste und `DECLARE` laufen ohne vorherige Sperre der Tabelle; eine
  gleichzeitige `ALTER TABLE` zwischen Export und `DECLARE` lässt den `DECLARE` mit einem
  Fehler der Klasse `storage` scheitern (`42703`) oder — bei einem Tabellen-Rewrite —
  den Katalog-Stand und den Snapshot auseinanderlaufen. **Nicht gemessen**, aus dem
  Verhalten von `SET TRANSACTION SNAPSHOT` abgeleitet; Hinweis für den `e2e`-Slice,
  keine erwartete Aktion in diesem Diff.
- `verifizierbar`: nein
- `klasse`: —

### F-18 — Klassifikation von `53300` als `configuration`

- `kategorie`: INFO
- `quelle`: [`SPEC-008`](../../spec/pflichtenheft.md), [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)
- `pfad`: `snapshot.go:392`
- `befund`: `53400` (alle Slots belegt) ist als `configuration` richtig, und der Test-Befund
  „Klasse `replication` statt `configuration` bei gestrichenem Mapping" ist richtig gelesen
  (R9). `53300` (`too_many_connections`, u. a. `max_wal_senders`) wird ebenfalls
  `configuration`; derselbe SQLSTATE steht auch für eine ausgeschöpfte `max_connections`
  der Quelle, ein vorübergehender Zustand — die Zuordnung ist vertretbar, im Plan nicht
  begründet, im Tier nicht gefahren (F-6).
- `verifizierbar`: nein
- `klasse`: —

## Negativbefunde

- geprüft, ohne Befund: Ablauf und Lebensdauer der Verbindungen in `snapshot.go`
  (Replication-Verbindung → temporärer Slot mit Export → Import in
  `REPEATABLE READ READ ONLY` → Schließen der Replication-Verbindung vor Spalten und
  Cursor; `closeConn` mit eigenem Zeitlimit idempotent, `defer` und expliziter Aufruf
  doppelt harmlos; jeder Fehlerpfad von `exportSnapshot` und `importSnapshot` schließt
  beide Verbindungen; kein Goroutinen-Start im Adapter; `Close` doppelt ohne Wirkung
  und `NextBlock` danach `storage` — `TestNextBlockAfterClose`, M4, M1, M2, M3 real
  auf PostgreSQL 17 und 18 grün, R1/R6 rot)
- geprüft, ohne Befund: Bezeichner- und Literal-Sicherheit — Slot-Name aus Run-Kennung ohne
  Bindestriche gegen `^[a-z0-9_]+$` und 63 Zeichen; Snapshot-Name gegen `^[0-9A-Fa-f-]+$`
  vor dem Literal in `SET TRANSACTION SNAPSHOT`; Spalten-, Schema-, Tabellennamen über
  `quoteIdent` (Doppelung `""`), Katalog-Abfragen parametrisiert
  (`to_regclass(quote_ident($1) || '.' || quote_ident($2))`); Werte nicht gebaut
  (nur zur Bindung der Zusage siehe F-4)
- geprüft, ohne Befund: Row-Image — der Adapter baut kein Bild, der Test ruft
  `model.BuildRowImage` (eine Konstruktionsstelle); gelöschte und generierte Spalten
  fehlen (R3 rot), NULL fällt weg (`{"id":"2"}`); GUC-Parität: das WAL-Bild trägt die
  Rollen-GUC real (`JST`, `dd.mm.yyyy`, `"1 2:03:04"`), der Vergleich ist nicht leer;
  — Typ-Grenze siehe F-1
- geprüft, ohne Befund: `EstimatedRows`/`reltuples` — `−1` wird `known = false`, nie `0`
  (gemessen `-1` auf 17.11 und 18.6 laut Testausgabe, `40` nach `ANALYZE`); keine
  Warn-Schwelle im Code (`git grep` nach Konstante/Schwelle im Paket: keine); Grenze
  `0` siehe F-9
- geprüft, ohne Befund: Fehlerklassen-Zuordnung `42501`/`28…` → `permission`
  (M2, falsches Passwort, R5), `53400`/`53300` → `configuration`, Kontext-Ende und
  `ConnectError` → `transient`; fünf Sentinels im Port dokumentiert — Testlücken siehe F-6
- geprüft, ohne Befund: Gate-Zuordnung — `Dockerfile:101` (R10 rot), Kommentare in
  `Dockerfile`, `db-coverage.sh`, `run-replication-tests.sh` zählen vier bzw. drei Pakete
  korrekt (bis auf F-2); Beschreibungen `coverage-gate.md`, `db-adapter-coverage.md`,
  `harness/README.md` nennen das Paket, Zahlen stimmen (848, 157, 139, 271×3, 8,51 %);
  §3.6: kein Gate gelockert (Schwellen 70/80 unverändert, Ausschluss folgt der
  Eigenschaft nach [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 1)
- geprüft, ohne Befund: Port `outbound/tablesnapshot.go` — Fähigkeits-Port
  ([`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md)), keine Adapter-Kante
  (`make a-check` Exit 0), keine Warn-Konstante, Erwartung `−1` als „unbekannt"
  abgebildet
- geprüft, ohne Befund: Commit-Struktur (§3.3 rein, Traceability), keine
  host-lokalen Pfade (§3.11), Handbuch unberührt (keine Betreiber-Oberfläche),
  Suppression (§3.2): kein `//nolint`; Docker-only (§3.1): keine Host-Toolchain im Diff

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 4 |
| LOW | 6 |
| INFO | 6 |

**Finding-Klassen dieses Laufs:** Zusage ohne Bindung an ihre Eingabeseite (F-1 Typ-Parität,
F-3 namentliche Listen, F-4 Bezeichner-Quoting, F-9 Nullgrenze) · Zahl im Träger driftet
gegen die Messung (F-2, F-14, F-15) · netzlos prüfbare Logik im Ausschluss-Paket (F-5) ·
fehlende Negativtests bei neuem öffentlichem Vertrag (F-6) · Beleg ohne reproduzierbaren
Befehl (F-10) · Zahl im Träger ohne Ursprung (F-7) · Kommentar mit Teilersetzungs-Rest (F-8)

## Verdikt

**Merge-blockierend:** ja — zwei HIGH (F-1: die Bild-Parität von WAL- und Backfill-Pfad
bricht für `boolean`, `inet`, `char(n)`; F-2: ein Skript-Kommentar gegen die Messung im
selben Diff) und vier MEDIUM (F-3 bis F-6). Die Gates am Stand laufen grün (`make test`,
`make a-check`, `make coverage-gate`, `make test-replication` auf PostgreSQL 18 und 17), die
Befunde liegen an Stellen, die kein Gate liest.

**Übergabe:** F-2 bis F-12 gehen an den Implementer. **F-1 geht zusätzlich an den
Architect**: die Formulierung `col::text` ist Wortlaut einer Accepted ADR
([`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 1 und 2,
§3.5 — eine Korrektur ist eine neue ADR mit `Supersedes`, nicht ein Implementer-Alleingang);
das Übergabe-Artefakt ist dieser Report. Die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" bleibt **offen**, weil eine Fixrunde nötig ist (Skill §DoD-Checkbox-Nachzug).
Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den Zähler; die
Klasse „Zusage ohne Bindung an ihre Eingabeseite" tritt in diesem Lauf viermal auf und ist
in `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` bereits gezählt. Dieser Report ist
ein **Lauf-Beleg**; DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
