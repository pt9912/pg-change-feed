# Verifier-Report: slice-031 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §1 (Ziel/Abgrenzung),
§2 (Definition of Done), §3 (Plan-Nachzug inkl. „Fixrunde nach Review"),
`ADR-0015` (Accepted, `permanent`, ungeändert) und der beiden
Review-Reports `docs/reviews/review-slice-031.md` (2 HIGH, 1 MEDIUM) und
`docs/reviews/review-slice-031-fixrunde.md` (alle drei behoben, keine
Regression).

**Grundsatz:** Keine Behauptung übernommen — jeder Beleg unten wurde in
diesem Lauf selbst gelesen, ausgeführt oder eigenständig reproduziert:
vollständiger Slice-Plan (§1–§8), beide Review-Reports, der tatsächliche
Code (`table_schema.go`, `schemastore.go` Port + Adapter,
`schemastore_test.go`, `tools/schema/schema.yaml`), `make gates` selbst
gestartet, `make test` selbst gestartet, `make test-store` selbst
gestartet (voller Lauf), zusätzlich die beiden namentlich verlangten
Testfälle (`TestSchemaStoreRegisterAndReadRoundTrip`,
`TestSchemaStoreRegisterVersionBackfillsColumns`) in einem eigenständig
aufgesetzten, separaten Testcontainer isoliert mit `-run`/`-v`
reproduziert, Abgrenzungs-Diff eigenständig gegen die drei genannten
Dateien gerechnet.

**Gegenstand:** `ca86249` (`TableSchema`-Modell, `SchemaStorePort`,
Postgres-Adapter), `cebc273` (DoD-Häkchen, Plan-Nachzug), `d25b726`
(Review-Report, 2 HIGH/1 MEDIUM), `100ff2b` (Fixrunde F-1/F-2/F-3),
`ec37b1e` (Review-Bestätigung Fixrunde).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-031-schema-persistenz-faehigkeit.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug inkl. Fixrunde, §6
  Risiken — noch `<bei Closure einzutragen>`, §7 Closure-Notiz — noch
  Platzhalter, §8 Sub-Area-Prüfungen)
- `docs/reviews/review-slice-031.md` (F-1/F-2/F-3, Verdikt: F-1/F-2
  merge-blockierend HIGH, F-3 MEDIUM kein Blocker für diesen Slice)
- `docs/reviews/review-slice-031-fixrunde.md` (alle drei bestätigt behoben,
  keine Regression, `make gates`/`make test-store` im Review-Lauf grün)
- `docs/plan/adr/0015-schema-evolution.md` (Accepted, `permanent`)
- `internal/domain/model/table_schema.go`,
  `internal/domain/model/table_schema_test.go`
- `internal/application/port/outbound/schemastore.go`
- `internal/adapters/driven/postgresstorage/schemastore.go`,
  `internal/adapters/driven/postgresstorage/schemastore_test.go`
- `internal/adapters/driven/postgresstorage/queries/*.go`
  (`InsertTableSchemaColumn`, `SelectTableSchemaColumns`,
  `CountTableSchemaColumns`, `SelectCurrentSchemaVersion`)
- `tools/schema/schema.yaml` (Tabelle `table_schema`)
- `docs/plan/planning/observations/BEO-PGC/schema-evolution-nicht-dynamisch/`
- `internal/adapters/driving/replication/mapper/mapper.go`,
  `internal/application/usecase/enable/service.go`,
  `internal/bootstrap/wiring.go` (Abgrenzungs-Prüfung §1)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify`: v6.5.0, 54 Dateien OK · `d-check` Standardlauf: 270 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 270 Dateien, 0 Befunde · `commit-traceability.sh`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test` | alle Pakete `ok`, inkl. `internal/domain/model` (deckt `table_schema_test.go`) | **0** |
| `make test-store` (voller Lauf) | alle Pakete `ok`, inkl. `internal/adapters/driven/postgresstorage` (3.052s), real gegen PostgreSQL-Testcontainer, Schema über d-migrate ausgerollt | **0** |
| eigenständiger isolierter Lauf: `go test ./internal/adapters/driven/postgresstorage/... -run 'TestSchemaStoreRegisterAndReadRoundTrip\|TestSchemaStoreRegisterVersionBackfillsColumns' -v` gegen einen selbst aufgesetzten, separaten Testcontainer (eigenes Docker-Netz `cdc-verify-031`, eigener Schema-Rollout) | `--- PASS: TestSchemaStoreRegisterAndReadRoundTrip (0.03s)` / `--- PASS: TestSchemaStoreRegisterVersionBackfillsColumns (0.03s)` | **0** |
| `git diff ca86249~1..HEAD -- internal/adapters/driving/replication/mapper/mapper.go internal/application/usecase/enable/service.go internal/bootstrap/wiring.go` | leer (kein Treffer, keine Zeile) | — |
| `git status`/`git diff` (am Ende dieses Laufs) | sauber, keine Restspur durch die Verifikation selbst (Testcontainer/-netz abgeräumt, `tools/schema/plan.yaml`/`down.sql` nach eigenem Rollout per `git checkout --` zurückgesetzt) | — |

`make doc-commits`/`make doc-immutable` existieren in diesem Repo nicht
(eigene Prüfung: `grep -n "^doc-commits\|^doc-immutable" -r Makefile
harness/mk/` — kein Treffer); die entsprechende Traceability-Deckung läuft
über `make commit-traceability` (Teil von `make gates`, oben grün). Eine
dedizierte ADR-Immutabilitäts-Prüfung ist mechanisch nicht verdrahtet (Hard
Rule 3.5 ist Prozessregel, kein Gate) — hier ohnehin ohne Berührungspunkt:
keiner der fünf slice-031-Commits ändert `docs/plan/adr/0015-schema-evolution.md`
(eigene `git show --stat` über alle fünf Commits).

## Prüfpunkt 1 — Abgrenzung §1: laufender Erfassungspfad wirklich unberührt

**Bestanden — eigenständig, nicht aus den Berichten übernommen.**

```
git diff ca86249~1..HEAD -- \
  internal/adapters/driving/replication/mapper/mapper.go \
  internal/application/usecase/enable/service.go \
  internal/bootstrap/wiring.go
```

liefert keine Ausgabe. Alle drei Dateien sind über den gesamten
Slice-Verlauf (Erst-Commit bis Fixrunde) unverändert. Der Plan-Nachzug
(„Wiring — entfällt bewusst") behauptet exakt das für `wiring.go` und
begründet es zusätzlich (keine Aufrufstelle, kein Gegenwert einer
ungenutzten DB-Verbindung) — Begründung nachvollziehbar und nicht
overclaimend: Die §3-Plan-Zeile zu `wiring.go` wird explizit als *nicht
eingelöst* markiert statt stillschweigend als erledigt geführt.

## Prüfpunkt 2 — DoD Punkt für Punkt gegen tatsächlichen Code

| # | DoD-Punkt | Verdikt | Beleg |
|---|---|---|---|
| 1 | `TableSchema`-Domänenmodell mit Spaltenmenge inkl. Typ-/OID-Information je `SchemaVersionID`, Unit-getestet | **bestätigt** | `table_schema.go`: `Column{Name, OID}`, `TableSchema{VersionID, Columns}`, `NewTableSchema` mit drei Invarianten (`ErrEmptyIdentifier`, `ErrEmptyColumns`, leerer Spaltenname); `table_schema_test.go` deckt Happy-Path, alle drei Fehlerfälle und Kopiersemantik (5 Testfälle). `make test` in diesem Lauf grün, Paket `internal/domain/model` explizit `ok`. |
| 2 | `SchemaStorePort` neu definiert, Postgres-Adapter real gegen `cdc.table_schema` (d-migrate), real gegen PostgreSQL getestet (Schreiben/Lesen/Round-Trip) | **bestätigt, eigenständig reproduziert** | Port: `CurrentVersion`/`RegisterVersion`/`TableSchema`, exakte Signatur wie im Plan-Nachzug dokumentiert. Adapter: eigene Datei `postgresstorage/schemastore.go`, `cdc.table_schema` in `tools/schema/schema.yaml` (PK `(schema_version_id, ordinal_position)`, UNIQUE `(schema_version_id, column_name)`, `column_oid biginteger`, FK auf `schema_version`) — stimmt mit den drei neuen Queries (`InsertTableSchemaColumn`, `SelectTableSchemaColumns`, `CountTableSchemaColumns`) und den fünf Testfällen in `schemastore_test.go` überein. `TestSchemaStoreRegisterAndReadRoundTrip` und `TestSchemaStoreRegisterVersionBackfillsColumns` in diesem Lauf **selbst** in einem separaten, eigenständig aufgesetzten Testcontainer mit `-v` reproduziert — beide `PASS`. Zusätzlich `make test-store` (voller Lauf) grün, inkl. aller fünf Testfälle des Pakets. |
| 3 | `make gates` grün | **bestätigt, eigenständig reproduziert** | `baseline-verify` (54 Dateien), `docs-check` (270 Dateien/0 Befunde — DoD nennt „268", Differenz durch die beiden Review-Reports d25b726/ec37b1e, die nach der DoD-Beleg-Formulierung entstanden; kein Sachfehler, nur ein seither gewachsener Dateibestand), `commit-traceability` (5 Commits, `HEAD~5..HEAD`, OK), `a-check` (0 Befunde) — alle in diesem Lauf grün. |
| 4 | Review durchgeführt, Report unter `docs/reviews/` | **bestätigt** (Checkbox im Plan noch offen — korrekt für diesen Lauf-Zeitpunkt) | Zwei Reports vorliegend: `review-slice-031.md` (Erstbefund) und `review-slice-031-fixrunde.md` (Bestätigung). Kein Self-Review — anderer Modell-/Sitzungs-Kontext laut Skill-Kopf beider Reports. |
| 5 | Doku-Update `harness/README.md`/`AGENTS.md`, falls neuer Sensor/Vertrag | **offen, korrekt unbeansprucht** | Kein neuer Sensor/Vertrag in diesem Slice (reine Fähigkeits-Lieferung ohne Live-Verdrahtung); Checkbox im Plan konsistent unmarkiert, keine Entscheidung dazu im Plan-Nachzug dokumentiert — das ist eine offene Closure-Pflicht für den Planner, kein Verifikations-Defekt. |
| 6–10 | Closure-Notiz, Reconciliation-Register (falls einschlägig), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **korrekt offen** | §6 trägt beide Risiken noch mit `<bei Closure einzutragen>`, §7 ist vollständig Platzhalter. Das ist der erwartete Zustand an dieser Stelle der Rollen-Sequenz (Implementer→Verifier, **vor** Planner-Closure, Modul 8): Diese fünf Punkte sind Planner-Arbeit bei der Closure nach `done/`, nicht Teil der Implementer→Verifier-Übergabe. Kein DoD-Defekt. |

**Zwischenstand:** Die drei vom Implementer bereits als erledigt markierten
DoD-Punkte (1–3) sind alle drei eigenständig nachgeprüft und **bestätigt**,
einschließlich der beiden namentlich verlangten Testfälle in einer
isolierten, eigenständigen Reproduktion. Die sieben Closure-Pflichten
(4–10) stehen konsistent mit dem aktuellen Lifecycle-Zeitpunkt (noch nicht
per Planner geschlossen) offen — keine dieser Lücken ist eine
DoD-**Verletzung**, sondern der korrekte Zwischenstand vor Closure.

## Prüfpunkt 3 — Trägt die Review-Fixrunde tatsächlich, was sie behauptet?

**Ja, eigenständig am Code nachvollzogen, nicht aus dem Fixrunden-Report
übernommen.**

- **F-1 (Kommentar-Chronik):** `grep -rn "slice-0[0-9][0-9]"
  internal/domain/model/table_schema.go
  internal/application/port/outbound/schemastore.go
  internal/adapters/driven/postgresstorage/schemastore.go` — eigenständig
  ausgeführt, kein Treffer in diesem Lauf. Alle drei Dateien vollständig
  gelesen (oben zitiert): Kommentare beschreiben durchgängig den
  Ist-Zustand („trägt ausschließlich die Persistenz-Fähigkeit"), keine
  Zeit-/Chronik-Formulierung, kein Review-Report-Pfad im Quellcode mehr.
- **F-2 (`CurrentVersion`-Kommentar):** Der Kommentar behauptet jetzt
  „macht ihre Version 1 dadurch bereits ohne einen
  `RegisterVersion`-Aufruf über `CurrentVersion` sichtbar" —
  `SelectCurrentSchemaVersion` liest unabhängig vom schreibenden Adapter
  aus derselben Tabelle `cdc.schema_version`; das Code-Verhalten stützt
  diese Aussage (`TableActivationAdapter.Register` und
  `PostgresSchemaStoreAdapter.RegisterVersion` schreiben dieselbe
  SQL-Konstante `queries.InsertSchemaVersion` in dieselbe Tabelle,
  `CurrentVersion` selektiert ohne Herkunfts-Unterscheidung).
- **F-3 (Backfill):** Code gelesen (`RegisterVersion`,
  `schemastore.go:111-148`): Die Spalten-Einfügung hängt jetzt an
  `existingColumns == 0` (`CountTableSchemaColumns`), nicht mehr an
  `RowsAffected()`der `schema_version`-Zeile. Das eigenständig
  reproduzierte `TestSchemaStoreRegisterVersionBackfillsColumns` legt die
  `schema_version`-Zeile direkt an (simuliert die Erstaktivierungs-Lage),
  registriert danach die Spaltenform nach — `backfilled=true`, liest sie
  über `TableSchema` zurück, und eine erneute Registrierung bleibt
  idempotent (`repeat=false`, `columnRows == 1`). **PASS** in diesem Lauf.

## Negativbefunde

- geprüft, ohne Befund: **Kein ADR-Inhalt geändert** — `git show --stat`
  über alle fünf slice-031-Commits: keiner berührt
  `docs/plan/adr/0015-schema-evolution.md`.
- geprüft, ohne Befund: **Kein Overclaiming ggü. `ADR-0015`-Folgepflicht.**
  Slice-Kopf und §1 benennen ausdrücklich, dass `LH-FA-SCH-004`/`005`
  hierdurch *nicht* geschlossen werden (das leisten `slice-032`/`slice-033`).
- geprüft, ohne Befund: **Bootstrap-Wiring-Verzicht ehrlich dokumentiert**
  — Plan-Nachzug benennt die Abweichung von der ursprünglichen §3-Plan-Zeile
  namentlich, keine stillschweigende Lücke.
- geprüft, ohne Befund: **Tabellenform `cdc.table_schema`** — PK, UNIQUE,
  FK und `biginteger`-Wahl für `column_oid` (OID ist vorzeichenloses
  32-Bit-Feld) stimmen mit der DoD-Beschreibung und den Test-Erwartungen
  überein.
- geprüft, ohne Befund: **Idempotenz-Regression nach F-3-Fix.**
  `TestSchemaStoreRegisterAndReadRoundTrip` (Erstregistrierung → Round-Trip
  → erneute Registrierung meldet `false`, keine doppelte Spaltenform) lief
  in derselben eigenständigen Reproduktion wie der Backfill-Test — kein
  Rückschritt an der bestehenden Idempotenz-Semantik.
- geprüft, ohne Befund: **§8 Sub-Area-Prüfung.** Einzige berührte Sub-Area
  `*`/`PGC` (Repo-Default), GF; Register-Sichtung nennt korrekt den
  bestehenden Treffer `BEO-PGC/schema-evolution-nicht-dynamisch`
  (eigenständig geprüft: Verzeichnis existiert mit `observation.md`,
  `state.md`, genau einer Evidence-Datei `evidence/slice-030.md` — „1×,
  weiter offen" ist damit korrekt, kein 3×-Übertritt).
- geprüft, ohne Befund: **Traceability** — alle fünf Commit-Betreffs tragen
  `ADR-0015`, kein `SPEC-*`/`ARC-*` im Betreff; `make commit-traceability`
  (Teil von `make gates`) in diesem Lauf grün.
- geprüft, ohne Befund: **`git status`/`git diff`** am Ende dieses Laufs —
  sauber, keine Restspur durch die Verifikation selbst.

## Eigene Befunde

Keine. Die einzige beobachtete Abweichung (`docs-check`-Dateizahl 270 statt
der im DoD-Beleg genannten 268) ist reine Zeitdrift durch danach entstandene
Dateien (die beiden Review-Reports selbst) und keine Aussage über den
Prüfzustand — nicht als eigener Befund geführt, da sie nichts über
DoD-Konformität aussagt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Zusammenfassung DoD:** 3/3 als „erledigt" markierte DoD-Punkte (1–3)
eigenständig vollständig nachgeprüft und bestanden, einschließlich beider
namentlich verlangten Testfälle in isolierter Reproduktion. 7/10
Closure-Pflichten (4–10) stehen korrekt offen — sie sind Planner-Arbeit bei
der Closure, kein Verifikations-Defekt an dieser Stelle der Rollen-Sequenz.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt.** Alle drei technischen
DoD-Punkte sind real erfüllt, nicht nur behauptet — Modell, Port und
Adapter sind vollständig und real gegen PostgreSQL getestet, wie die DoD
es verlangt. Die Fixrunde hat alle drei Review-Findings (F-1/F-2 HIGH,
F-3 MEDIUM) tatsächlich behoben, eigenständig am Code und über reale
Testläufe nachvollzogen — keine Regression. Die Abgrenzung §1 (laufender
Erfassungspfad unberührt) hält: `mapper.go`, `enable/service.go`,
`wiring.go` sind über den gesamten Slice-Verlauf unverändert. `ADR-0015`
bleibt `Accepted` und unverändert; keine Anforderung wird fälschlich als
geschlossen dargestellt.

**Plan-vs-Code-Diff:** Keine Abweichung zwischen §3-Plan-Nachzug und
tatsächlichem Code gefunden — Modellform, Port-Signatur, Adapter-Ort,
Tabellenform und der bewusste Wiring-Verzicht sind exakt so umgesetzt, wie
der Plan-Nachzug es dokumentiert.

**Closure-Bereitschaft:** Aus Verifikations-Sicht steht der Fähigkeits-Teil
(DoD 1–3) closure-bereit. Vor dem `git mv` nach `done/` fehlen noch die
sieben Planner-Pflichten (§6 Risiko-Ausgänge, §7 Closure-Notiz mit
Steering-Loop-Lerneintrag — die drei Finding-Klassen aus
`review-slice-031.md` gehören dort hinein — und die Register-/Paarungs-Prüfung),
die dieser Bericht bewusst nicht vorwegnimmt.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert (`git status`/`git diff` am
Ende sauber).

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: v6.5.0, 54 Dateien OK; `d-check`
Standardlauf: 270 Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD`
Modul `commits`: 270 Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5
Commits; `a-check`: 0 Befunde). `make test` einmal ausgeführt, Exit 0,
alle Pakete `ok`. `make test-store` einmal ausgeführt (voller Lauf), Exit
0, alle Pakete `ok`, inkl. `postgresstorage` real gegen PostgreSQL.
Zusätzlich `TestSchemaStoreRegisterAndReadRoundTrip` und
`TestSchemaStoreRegisterVersionBackfillsColumns` einzeln mit `-v` in einem
eigenständig aufgesetzten, separaten Testcontainer reproduziert, beide
`PASS`. `make doc-commits`/`make doc-immutable` existieren in diesem Repo
nicht (eigene Makefile-Prüfung) — Traceability-Deckung läuft über
`commit-traceability` (Teil von `make gates`); ADR-Immutabilität ist hier
ohne Berührungspunkt (kein Commit ändert `docs/plan/adr/0015-schema-evolution.md`).
`git status`/`git diff` am Ende dieses Laufs sauber (keine
Arbeitsverzeichnis-Änderung durch die Verifikation selbst).
