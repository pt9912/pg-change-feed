# Review-Report: slice-backfill-e2e — 2026-09-24

**Review-Art:** Code — der Diff führt die Backfill-Rundläufe des E2E-Runners (sieben, davon
sechs Runner-Phasen und ein Go-Test), die Erweiterung des Wegwerf-Clients, die Fixrunde zu
[`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) (Lesesperre und
Filenode-Vergleich im Snapshot-Import, drei Test-Tiers), einen Pflichtenheft-Satz und die
Träger (Handbuch 1.51, `harness/README.md`, E2E-Abdeckungstabelle, Plan) ein; geprüft gegen
Plan, ADRs, Pflichtenheft und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein
DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-backfill-e2e`, Diff-Range `e7df5619..360ccb49` (15 Commits,
18 Dateien, +2407/−366; drei reine `git mv`-Commits: `ed0f99a8`, `4a315549`, `baa75c96`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz, Zusage-ohne-Eingabeseite,
Kommentar-Chronik, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-24.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-e2e` (§1 Ziel, §2 DoD als Prüfmaßstab für Plan-Zusagen, §3 Plan
  samt Suchlauf-Feld, §6 Risiken) und Welle `welle-backfill-bestand`; Architect-Verdikte
  `architect-verdict-backfill-schema-klasse-rollen` und
  `architect-verdict-backfill-tabellen-rewrite-im-fenster`
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
  [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md),
  [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md),
  [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md),
  [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) (Festlegungen 1–5,
  Folgepflichten 1–4, Fitness-Function-Tabelle),
  [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
- [`LH-FA-CAP-009`](../../spec/lastenheft.md) (Happy Path, Boundary, Negative),
  [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) (Absatz „Mechanismus“),
  [`SPEC-029`](../../spec/pflichtenheft.md), [`LH-FA-CON-005`](../../spec/lastenheft.md)
- `AGENTS.md` (Hard Rules §3.1–§3.13), `harness/conventions.md` (`MR-000`/`MR-001`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-sql-administration.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht
übernommen; Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert):

- **Gates am Stand `360ccb49`:** `make test` Exit 0 (`-race`); `make a-check` Exit 0 („gesamt: 0
  Befund(e)“); `make coverage-gate` Exit 0, gedruckt „coverage-gate: OK — Coverage 83.00%
  erfüllt Schwelle 80%“; `bash tools/harness/run-replication-tests.sh measure` (PostgreSQL 18,
  Pin von `PG_TEST_IMAGE`) Exit 0, gedruckt „DB-Adapter-Coverage: 82.13% (gedeckt 850 von 1035
  Statements; Profile gemergt: store,replication)“; `make doc-trace` Exit 0, gedruckt „80
  Anforderung(en), 2 Waise(n)“ (`LH-FA-CAP-009` trägt den Nachweis `E2E`).
- **`make gates`** mit dem Report im Baum, ungefiltert: Exit 0, gedruckt „d-check: 1074 Datei(en) geprüft, 0 Befund(e)“, „coverage-gate: OK — Coverage 83.00% erfüllt Schwelle 80%“, „commit-traceability: OK“, „generated-sync: OK“, „gesamt: 0 Befund(e)“.
- **E2E-Lauf:** `make image` Exit 0, danach `make test-integration` Exit 0 in 277 s (Zeitstempel
  vor und nach dem Aufruf); jede der sechs Backfill-Phasen druckt ihre „belegt“-Zeile, der Go-Lauf
  druckt „Replay-Invariante: Snapshot-Position 31463936, 48 Backfill-Changes, WAL-Changes davor
  284 und dahinter 110, 56 Zeilen im Quellstand“, der Runner „E2E-Abdeckungstabelle unverändert“
  und „aus 14 Go-Zeilen und 34 Bash-Zeilen“; die Pausen des Feed-Containers messen 217 ms und
  357 ms (DDL-Fenster). `git status` danach ohne Änderung (`tools/schema/plan.yaml`,
  `tools/schema/down.sql` und `docs/user/e2e-abdeckung.md` unverändert).
- **Mutationen der Eingabeseite** (gegen PostgreSQL 18, interne Tests des Pakets
  `postgressnapshot`, je danach `git checkout` der Datei; alle Mutationen färben rot):

  | Nr. | Mutation | roter Test |
  |---|---|---|
  | M1 | `relfilenode <> …` zu `= …` | `TestRewriteInWindowIsTransient` (vier Formen: `err = nil`), `TestNoRewriteInWindowReadsTheSnapshot`, `TestPartitionedTableIsNoFalseAlarm`, `TestImportWaitsForExclusiveLockAndThenAborts` |
  | M2 | `coalesce(pg_relation_filenode(c.oid), 0)` zu `pg_relation_filenode(c.oid)` | `TestPartitionedTableIsNoFalseAlarm` („Umschreib-Prüfung: unerwartete Antwort“) |
  | M3 | Sperranweisung durch `SELECT 1` ersetzt | `TestImportWaitsForExclusiveLockAndThenAborts` (wartet nicht an `LOCK TABLE`) |
  | M4 | Vergleich vor die Sperre gesetzt | `TestImportWaitsForExclusiveLockAndThenAborts` (`err = nil` nach dem Commit) |
  | M5 | `ClassifyLock`: `3F000` aus der Klasse `configuration` entfernt | `TestConfigurationClass` (fehlendes Schema endet `storage`) |
  | M6 | Klasse des Umschreibens `transient` zu `storage` | `TestRewriteInWindowIsTransient` (alle fünf Formen) |
  | M7 | `lockAndVerify` hinter `readColumns` gesetzt | `TestCatalogQueryFailuresKeepTheirClass` (Phase `Spaltenliste` statt `Umschreib-Prüfung`), `TestRewriteInWindowIsTransient` (`DROP und CREATE` endet `configuration`) |
  | M8 | `MarkRunning` hinter `copyBlocks` gesetzt | `TestExecuteMarksRunningBeforeOpeningSnapshot` („festgehaltene Zustände = [], will [running]“) |

- **Wartegrenze der Sperre** (Scratch-Test, danach gelöscht): `importSnapshot` mit Kontext-Zeitlimit
  2 s gegen eine fremde Transaktion mit `ACCESS EXCLUSIVE` kehrt nach 2,001 s mit Klasse
  `transient` („Tabellensperre: timeout: context deadline exceeded“) zurück, danach
  bleibt keine Sitzung des Runs.
- **`RENAME COLUMN` im Fenster** (Scratch-Test, danach gelöscht): endet am Cursor mit Klasse
  `storage` („column "name" does not exist (SQLSTATE 42703)“).
- **Sperr-Warteschlange** (Wegwerf-Container, PostgreSQL 18, Ergebnis siehe F-1).
- **Suchläufe des Plans** an beiden Ständen nachgefahren (Parent `e7df5619`, Diff-Stand): siehe
  Negativbefunde und F-3.
- **Umgebung:** `free -m` vor dem E2E-Lauf 20,6 GB verfügbar; dangling Volumes
  (`docker volume ls -q -f dangling=true | wc -l`) 34 vor und 34 nach allen Läufen; kein `prune`;
  eigene Wegwerf-Container (`rv-pg`) mit `docker rm -fv` entfernt, kein Container nach dem Lauf.
- **Nicht gefahren (Grenze):** die PostgreSQL-17-Leg der E2E-Matrix und der Store-/Replication-Tier
  gegen 17; Mutationen der E2E-Assertions (ein Lauf dauert 277 s; gelesen statt gemutet, siehe
  F-4); GitHub-Runner-Lauf (kein
  Workflow im Diff, `AGENTS.md` §3.10 greift nicht).

---

## Findings

### F-1 — Handbuch nennt für die Sperr-Warteschlange nur „Leser“; gemessen stauen sich auch Schreiber und die Publication-Abfrage der Administration

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Konsequenzen
  (Sperr-Warteschlange „hergeleitet, nicht gemessen“) · `AGENTS.md` §3.12 Instanz A/B
  (Aussage über den Gegenstand trägt ihren Ursprung) · Reviewer-Skill „Beleg trägt seinen Satz
  nicht“
- `pfad`: `docs/user/benutzerhandbuch.md:452-459` (Absatz „Sperre der Tabelle“)
- `befund`: Der Absatz sagt „Lesen und Schreiben der Tabelle laufen weiter“ und nennt als Folge
  einer wartenden DDL nur „Leser der Tabelle, die nach ihr anfragen“ (als hergeleitet
  gekennzeichnet). Gemessen (PostgreSQL 18, Sitzung A `BEGIN; LOCK TABLE … IN ACCESS SHARE MODE`
  hält, Sitzung B `ALTER TABLE … ALTER COLUMN … TYPE` wartet): ein `INSERT`, ein
  `SELECT count(*)` und `SELECT count(*) FROM pg_publication_tables WHERE pubname = …` (die
  Abfrage der Publication-Mitgliedschaft im Antrag und am Run-Start) liefen je in eine
  4-s-Grenze (gedruckt: „INSERT exit=143 after 4063 ms“, „SELECT exit=143 after 4053 ms“,
  „pg_publication_tables exit=143 after 4056 ms“) — Schreiber der Quelltabelle und die
  Administration stehen für die Dauer des Runs hinter der DDL, nicht nur Leser. Der Absatz
  bleibt als Betreiber-Warnung enger als die Messung.
- `verifizierbar`: ja — die genannte Abfolge mit zwei Sitzungen und einer dritten Anweisung
  reproduziert es.
- `klasse`: Betreiber-Warnung enger als die Messung

### F-2 — Sequenzdarstellung des Snapshot-Lesers in der Architektur-Sicht trägt Sperre und Umschreib-Prüfung nicht; der Suchlauf des Plans hat den Träger nicht gefunden

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Träger einer bewegten Eigenschaft nachziehen) ·
  [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Folgepflicht 4
- `pfad`: `spec/architecture.md:329` (Sequenzdiagramm der Run-Ausführung; Aufzählung der
  Rollback-Fälle im Absatz darunter)
- `befund`: Das Diagramm führt den Import als „Transaktion (REPEATABLE READ) mit importiertem
  Snapshot“, die Absätze darunter nennen als Endzustand `failed` nur die Fail-closed-Abweichung;
  Lesesperre und ein Run-Ende bei umgeschriebener Tabelle stehen nirgends. Das Suchmuster des
  Plans (`Lesesperre`, `Umschreib`, `Tabellensperre` …) trifft die Datei nicht; das Feld nennt sie
  weder unter „gefunden“ noch unter „nicht gefunden“. Die Sicht bleibt wahr, ist aber
  unvollständig; ob sie die Sperre tragen soll, entscheidet der Architect (die Sicht führt keine
  ADR-Bezüge, `AGENTS.md` §3.4).
- `verifizierbar`: ja — `git grep -n -E "REPEATABLE READ|importiert" spec/architecture.md`.
- `klasse`: Träger der bewegten Eigenschaft nicht nachgezogen

### F-3 — Zahl im Suchlauf-Feld des Plans (`make test-integration`-Treffer) driftet gegen den Diff-Stand

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz A (übernommener Wert gegen die eigene Messung) ·
  Reviewer-Skill „Zahl im Träger … gegen die Messung driftend“
- `pfad`: `docs/plan/planning/in-progress/slice-backfill-e2e.md:228` (Zeile „Beschreibungen von
  `make test-integration` außerhalb der README“)
- `befund`: Das Feld nennt für den Diff-Stand „127 Zeilen in 59 Dateien“. Gemessen mit derselben
  Pathspec-Form (ohne die Plan-Dateien): Parent `e7df5619` 124 Zeilen in 57 Dateien (stimmt), Stand
  `26ea0c2e` 127 in 59 (stimmt), Stand `360ccb49` und Arbeitsbaum 129 in 60 — die Fixrunde legt
  [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) an, die zwei
  Treffer in einer neuen Datei trägt. Die Behandlung („unverändert“) bleibt richtig; die Zahl
  gehört zu einem Zwischenstand, nicht zum Diff-Stand des Slice.
- `verifizierbar`: ja — `git grep -n 'test-integration' <Stand> -- . ':!docs/reviews'
  ':!docs/plan/planning/done' ':!docs/plan/planning/observations'
  ':!docs/plan/planning/in-progress/slice-backfill-e2e.md' | wc -l` an beiden Ständen.
- `klasse`: Zahl im Träger gegen den Diff-Stand gedriftet

### F-4 — E2E-Assertion „Fehlerzustand des Erfassungspfads“ kann bei fehlender Heartbeat-Zeile nicht rot werden

- `kategorie`: LOW
- `quelle`: Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“
- `pfad`: `tools/harness/run-integration-tests.sh:3100` (`bf_ddl_window`)
- `befund`: `bf_expect "$(bf_sql "SELECT coalesce(error_class, '') FROM cdc.heartbeat WHERE
  source_id = 'src-e2e'")" ""` liest den leeren String sowohl bei einer Zeile mit NULL als auch
  bei keiner Zeile; eine fehlende Heartbeat-Zeile (Quelle unbekannt, Filter falsch) hielte die
  Zusage „der Run-Fehler ist run-lokal, der Erfassungspfad läuft weiter“ ebenso grün. Die
  Aussage trägt in der Phase Negative die anschließende Erfassung nach dem Neustart, in der Phase
  DDL-Fenster nur `docker inspect … Running`. Gelesen, nicht gemutet (Grenze oben).
- `verifizierbar`: ja — eine Mutation des Quellfilters auf eine nicht vorhandene Quelle lässt die
  Assertion grün.
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-5 — Zwei Kommentare tragen neben Zusage/Kopplung die verworfene Alternative

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (ein Kommentar beschreibt, was da ist) · Reviewer-Skill „Kommentar
  trägt keine der Kommentar-Klassen“
- `pfad`: `test/integration/backfill_e2e_test.go:267-268`;
  `tools/harness/run-integration-tests.sh:3163`
- `befund`: Der Test-Godoc schließt „Der Test bricht ab, wenn eine der beiden Seiten leer bleibt,
  statt ohne Überlappung grün zu werden“ (die Zusage steht vorn, das „statt …“ nennt die
  Alternative); der Runner-Kommentar erklärt „(eine offene Transaktion während der Slot-Anlage
  würde diese verzögern)“ im Konjunktiv über den nicht gewählten Ablauf. In beiden Fällen ist der
  Hauptsatz Zusage bzw. Kopplung und trägt die Stelle; die Nebenklausel ist die Form, die der
  Skill benennt. Kein Produktionscode, keine Slice-/Wellen-Nummer, keine Chronik.
- `verifizierbar`: nein — kein Gate liest Kommentar-Sprache.
- `klasse`: Kommentar trägt keine Kommentar-Klasse

### F-6 — Der Run wartet an der Sperranweisung ohne Zeitgrenze; Handbuch nennt es nicht, kein committeter Test bindet den Kontext-Abbruch

- `kategorie`: LOW
- `quelle`: [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md)
  Festlegung 1 („die Wartezeit trägt der Kontext des Aufrufs, sein Ablauf endet als
  `transient`“) und Konsequenzen („der Run wartet auf eine DDL“)
- `pfad`: `internal/adapters/driven/postgressnapshot/snapshot.go:207-221` (`lockAndVerify`);
  `docs/user/benutzerhandbuch.md:452-459`
- `befund`: Der Kontext, den der Worker übergibt, trägt kein Zeitlimit; eine Transaktion, die
  `ACCESS EXCLUSIVE` hält, lässt den Run in `running` (`rows_copied` 0) stehen, während seine
  Lese-Transaktion den Snapshot-Horizont festhält. Das Handbuch beschreibt nur die Gegenrichtung
  (die DDL wartet auf den Run). Der Kontext-Abbruch an der wartenden Anweisung ist gemessen
  korrekt (Klasse `transient`, keine Sitzung danach, siehe oben), aber in keinem committeten
  Test gebunden (`TestClassifyLock` klassifiziert nur den Fehlerwert). Verhalten ist von der ADR
  angenommen; hier steht nur, dass Betreiber-Text und Test-Bindung fehlen.
- `verifizierbar`: ja — der Scratch-Test aus den Eigenprüfungen (Kontext-Zeitlimit gegen eine
  wartende Sperre) ist die Probe.
- `klasse`: Wartegrenze ungenannt und ungebunden

### F-7 — Der Satz „`RENAME COLUMN` … endet am `DECLARE` … (E2E-belegt)“ hat nur für `DROP COLUMN` einen Beleg im Repo

- `kategorie`: INFO
- `quelle`: [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) Festlegung 5
  · Reviewer-Skill „Beleg trägt seinen Satz nicht“ (Rolle: Architect für die ADR, Planner für das
  Handbuch)
- `pfad`: `docs/plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md:207`;
  `docs/user/benutzerhandbuch.md:509-521`
- `befund`: Die Phase DDL-Fenster und `TestNoRewriteInWindowReadsTheSnapshot` fahren `DROP COLUMN`;
  ein `RENAME COLUMN`-Lauf steht nirgends im Diff. Gemessen ist die Aussage wahr (Klasse `storage`,
  SQLSTATE 42703); der Beleg im Repo fehlt. Ohne Aktion des Implementers; die ADR ist `Accepted`
  und unberührbar.
- `verifizierbar`: nein — kein Lauf des Repos nennt die Form.
- `klasse`: Beleg trägt seinen Satz nicht

### F-8 — Fehlertext des Runs trägt die Klasse doppelt

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: gedruckte Zeile der Phase DDL-Fenster im E2E-Lauf: „failed (transient: Fehlerklasse
  transient: Quelle für den Tabellen-Snapshot vorübergehend nicht verfügbar: Tabelle
  public.feed_e2e_backfill_rewrite wurde nach dem Snapshot-Export umgeschrieben; ein neuer Antrag
  beginnt neu)“
- `befund`: Das Präfix „transient: “ (Run-Fehlertext) und der Sentinel-Text „Fehlerklasse
  transient: …“ nennen dieselbe Klasse zweimal hintereinander; dasselbe Muster trägt die Zeile zu
  `storage`. Bestehende Form des Sentinel-Textes, vom Diff nicht eingeführt, nur sichtbar
  gemacht; das Handbuch nennt nur das erste Präfix.
- `verifizierbar`: ja — die Phase DDL-Fenster druckt es.
- `klasse`: Fehlertext redundant

### F-9 — Zeitannahmen der Haltepunkte sind lokal belegt, für den GitHub-Runner offen

- `kategorie`: INFO
- `quelle`: Plan §6 (Risiko „Laufzeit des erweiterten Testpakets“) · `AGENTS.md` §3.10
- `pfad`: `tools/harness/run-integration-tests.sh:3094` und `:3200` (Pausen-Grenze 1000 ms),
  `.github/workflows/e2e.yml` (unverändert)
- `befund`: Die Grenze „Pause unter 1000 ms, die Hälfte von `wal_sender_timeout`“ ist benannt und
  scheitert sichtbar (kein stilles Grün); lokal gemessen 217 ms und 357 ms, Gesamtlauf 277 s. Ob
  der GitHub-Runner sie hält, ist bis zum ersten Push-Lauf nicht gemessen. Der Workflow ist im
  Diff nicht berührt; das Risiko bleibt im Plan §6 offen und ist Sache von Verifier und Closure.
- `verifizierbar`: nein — nur ein Runner-Lauf.
- `klasse`: Lokal belegt, Runner offen

### F-10 — `gofmt -l` nennt eine Datei des Pakets `test/integration`, die der Diff nicht berührt

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `test/integration/integration_test.go` (nicht im Diff)
- `befund`: `gofmt -l` über die vom Diff berührten Pakete (`postgressnapshot`, `usecase/backfill`,
  `test/integration`, `tools/harness/httpclient`) listet nur diese Datei; alle Dateien des Diffs
  sind formatiert. Bestand, keine Aktion dieses Slice.
- `verifizierbar`: ja — `gofmt -l test/integration` im Toolchain-Container.
- `klasse`: Bestands-Formatierung

## Negativbefunde

- geprüft, ohne Befund: `internal/adapters/driven/postgressnapshot/snapshot.go` und
  `snapshotlogic/logic.go` — Bezeichner-Quoting der Sperranweisung (`QuoteIdent`, Test mit
  eingeschleustem `"; DROP TABLE`), Parameterbindung der Abfrage, Reihenfolge Import → Sperre →
  Vergleich → Spaltenliste → Cursor (M4, M7), Klassen (`42P01`/`3F000` → `configuration`, Umschreiben
  → `transient`, `42501` → `permission` in der Phase `Tabellensperre`, M5, M6), Abbruchpfad
  (`snap.abort()`, keine Sitzung danach), Fehlalarm partitionierte Tabelle (M2).
- geprüft, ohne Befund: `internal/adapters/driven/postgressnapshot/*_test.go` — jede Zusage des
  Tiers an ihrer Eingabeseite gemutet (M1–M7); Tests unter echten Rollen-Logins tragen die
  `SELECT`-Zusage (`TestPermissionClassWithoutSelect`, `TestCatalogQueryFailuresKeepTheirClass`);
  keine Fakes, die nur permanent scheitern.
- geprüft, ohne Befund: `internal/application/usecase/backfill/` — `TestExecuteMarksRunningBeforeOpeningSnapshot`
  ist an seiner Eingabeseite gebunden (M8); der Kommentar an `currentVersion` verweist auf
  [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md) Festlegung 1/2 (in
  der ADR vorhanden), kein Verhaltens-Diff.
- geprüft, ohne Befund: `test/integration/backfill_e2e_test.go` — nur externe Wege (SQL, keine
  Importe aus `internal/**`), Überlappung der beiden Seiten ist Abbruchbedingung des Tests,
  gedruckte Zeile „Replay-Invariante: …“ im Lauf.
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` (Backfill-Abschnitt) — Haltepunkte
  (offene Schreibtransaktion, unbestätigter Schlüssel, `docker pause` unter der Hälfte von
  `wal_sender_timeout`), Cleanup (`docker unpause` und `down -v` im `trap`), Fehlerpfade
  (`bf_fail` mit Exit 1, keine Pipe auf einem Gate-Aufruf, `set +e` nur im Client-Aufruf und
  wiederhergestellt), Superuser-DSNs der Compose-Umgebung (die Rollen-Zusage der Sperre trägt der
  Store-Tier, im Plan §3 benannt), jede Phase mit eigener Tabelle; die
  Startpositions-Werte (Position 0, `acknowledged=false`, 5 Backfill-Changes, 0 hinter der
  bestätigten Position) stimmen mit der gedruckten Zeile des Laufs überein.
- geprüft, ohne Befund: `tools/harness/httpclient/main.go` — neue Modi `changes`/`position`,
  Herkunft geprüft gegen `wal`/`backfill`.
- geprüft, ohne Befund: `spec/pflichtenheft.md` — der Satz in
  [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) präzisiert den Mechanismus, ohne ADR- oder
  Slice-Bezug, das Lastenheft ist unberührt (kein Stratum-Verstoß); Änderungshistorie fortgeschrieben.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — `Version:` 1.49 → 1.51 mit den Zeilen 1.50
  und 1.51 in `### Änderungshistorie`; keine neue Betreiber-Oberfläche (keine `CDC_*`-Variable,
  keine `cdc.*`-Funktion, kein Endpunkt); die gemessene Startposition trägt ihren Lauf als Ursprung,
  hergeleitete Aussagen sind gekennzeichnet (Ausnahme F-1).
- geprüft, ohne Befund: `harness/README.md` — Zeilen `make test-integration`, `make test-replication`
  und `make doc-trace`: „sieben Backfill-Rundläufe“ nachgezählt (6 `abdeckung_declare "Backfill…“` +
  1 `func TestE2EBackfill…`), „drei Haltepunkte“ gegen den Runner-Kopf, „2 Waisen“ gegen `make
  doc-trace` (80/2, `LH-FA-CFG-007`/`LH-FA-CFG-008`).
- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` — Erzeugnis: der Lauf meldet „unverändert“;
  Parent 41 Zeilen / 0 mit `LH-FA-CAP-009`, Diff-Stand 48 / 7 (nachgezählt).
- geprüft, ohne Befund: `docs/plan/adr/` (Datei zu [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) und Index) — Index-Zeile vorhanden, Status
  `Accepted`, „Supersedes“ [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md), teilweise wie bei [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)
  bis [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md).
- geprüft, ohne Befund: Plan-Suchlauf-Feld (§3) an beiden Ständen nachgefahren — Abdeckungstabelle
  41 → 48 und 0 → 7, Zählwort-Suche 12 (11 ohne die Plan-Datei am alten Pfad) → 11, Träger der
  Fixrunde 4 → 37, `httpclient` 39 → 48 in je 14 Dateien, sieben Fremdobjekte in
  `knownForeignObjects` (Ausnahmen F-2, F-3).
- geprüft, ohne Befund: Hard Rules — Docker-only (keine Host-Toolchain im Diff), kein host-lokaler
  absoluter Pfad in Doku, keine Suppression (`nolint`/`noqa`), keine Slice-/Wellen-Chronik in
  Produktionscode-Kommentaren, alle 15 Commits nennen mindestens eine `LH-*`-/`ADR-*`-Kennung, kein
  Betreff trägt eine `SPEC-*`-/`ARC-*`-Kennung, die drei `git mv`-Commits sind rein (0
  Einfügungen/Löschungen), keine Verletzung von `AGENTS.md` §3.9 im Runner.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 5 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Betreiber-Warnung enger als die Messung · Träger der bewegten
Eigenschaft nicht nachgezogen · Zahl im Träger gegen den Diff-Stand gedriftet · Zusage ohne
Bindung an ihre Eingabeseite · Kommentar trägt keine Kommentar-Klasse · Wartegrenze ungenannt
und ungebunden · Beleg trägt seinen Satz nicht · Fehlertext redundant · Lokal belegt, Runner
offen · Bestands-Formatierung

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM): der Handbuch-Absatz zur Sperre bleibt hinter der
gemessenen Wirkung der Sperr-Warteschlange zurück (Schreiber und Administration, nicht nur
Leser); eine Betreiber-Warnung zu einer DDL im Produktionsbetrieb ist der Ort, an dem eine
zu enge Aussage Schaden stiftet. Die Produktions- und Test-Umsetzung von
[`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md) trägt: 0 HIGH, alle
Zusagen der Fixrunde färben an ihrer Eingabeseite rot (M1–M8), der E2E-Lauf ist grün. F-2 bis F-6
(LOW) und F-7 bis F-10 (INFO) blockieren nicht; F-2 braucht eine Entscheidung des Architects, ob
die Sicht die Sperre tragen soll.

**Übergabe:** F-1 sowie F-3 bis F-6 gehen an den Implementer (Fixrunde; F-1 ist eine
Text-Korrektur im Handbuch samt Versionshistorie, F-3 eine Zahl im Plan, F-4 eine Assertion,
F-5 zwei Kommentare, F-6 ein Absatz und ein Test); F-2 geht an den Architect
(Sicht-Frage) und an den Planner (Träger-Nachzug); F-7 an den Architect zur Kenntnis. Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Die
DoD-Zeile „Review durchgeführt“ im Plan bleibt offen, weil eine Fixrunde folgt (Skill
§DoD-Checkbox-Nachzug ohne Fixrunde greift nicht). Dieser Report ist ein **Lauf-Beleg** und
ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
