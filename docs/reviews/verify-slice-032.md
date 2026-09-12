# Verifier-Report: slice-032 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §1 (Ziel/Abgrenzung),
§2 (Definition of Done), §3 (Plan/Plan-Nachzug), `ADR-0015` (Accepted,
`permanent`, ungeändert), `LH-FA-SCH-005` (`spec/lastenheft.md`) und der
Review-Report `docs/reviews/review-slice-032.md` (0 HIGH/MEDIUM, 1 LOW,
kein Merge-Blocker).

**Grundsatz:** Keine Behauptung übernommen — jeder Beleg unten wurde in
diesem Lauf selbst gelesen, ausgeführt oder eigenständig reproduziert:
vollständiger Slice-Plan (§1–§8), Review-Report, der tatsächliche Code
(`decode.go`/`decode_test.go`, `mapper.go`/`mapper_test.go`, `receive.go`,
`wiring.go`), `make gates` selbst gestartet, `make test` selbst gestartet,
`make test-integration` **dreimal** in eigenständigen, unabhängigen Läufen
gestartet (nicht aus dem Review-Report übernommen), mit besonderem
Augenmerk auf die beiden namentlich verlangten Testfälle
(`TestMVPSchemaChangeAddColumn`, `TestMVPSchemaChangeIncompatibleTypeChange`).

**Gegenstand:** `bc40556` (Decoder-Erweiterung, Mapper-Klassifikation und
-Verdrahtung, Rollen-/Testschema-Lücken, `harness/image-hash.txt`),
`5882056` (DoD-Häkchen, Plan-Nachzug), `aa283ec` (Review-Report, 0 HIGH/
MEDIUM, 1 LOW).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-032-dynamische-re-versionierung.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug, §6 Risiken — noch
  `<bei Closure einzutragen>`, §7 Closure-Notiz — noch Platzhalter, §8
  Sub-Area-Prüfungen)
- `docs/reviews/review-slice-032.md` (F-1 LOW, kein Merge-Blocker, Verdikt:
  keins)
- `docs/plan/adr/0015-schema-evolution.md` (Accepted, `permanent`)
- `spec/lastenheft.md` (`LH-FA-SCH-004`, `LH-FA-SCH-005`)
- `internal/adapters/driving/replication/decode/decode.go`,
  `decode_test.go`
- `internal/adapters/driving/replication/mapper/mapper.go`,
  `mapper_test.go`
- `internal/adapters/driving/replication/receive/receive.go`
- `internal/bootstrap/wiring.go`
- `test/integration/integration_test.go`
  (`TestMVPSchemaChangeAddColumn`, `TestMVPSchemaChangeIncompatibleTypeChange`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify`: v6.5.0, 54 Dateien OK · `d-check` Standardlauf: 273 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 273 Dateien, 0 Befunde · `commit-traceability.sh`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test` (voller Lauf, alle Pakete) | alle Pakete `ok`, inkl. `internal/adapters/driving/replication/decode`, `.../mapper`, `.../receive`, `internal/bootstrap` | **0** |
| `make test-integration`, Lauf 1/3 (eigenständig gestartet) | `TestMVPSchemaChangeAddColumn` PASS (0.24s) · `TestMVPSchemaChangeIncompatibleTypeChange` PASS (0.25s, PostgreSQL lehnt Fall 1 selbst ab, Fall 2 übernimmt konvertierten Wert) · alle übrigen 5 Tests PASS | **0** |
| `make test-integration`, Lauf 2/3 (eigenständig gestartet) | dieselben Ergebnisse, `TestMVPSchemaChangeAddColumn` PASS (0.24s), `TestMVPSchemaChangeIncompatibleTypeChange` PASS (0.25s) | **0** |
| `make test-integration`, Lauf 3/3 (eigenständig gestartet) | dieselben Ergebnisse, `TestMVPSchemaChangeAddColumn` PASS (0.23s), `TestMVPSchemaChangeIncompatibleTypeChange` PASS (0.24s) | **0** |
| `git diff bc40556 -- test/integration/integration_test.go \| grep -c TestMVPSchemaChangeIncompatibleTypeChange` | `0` — dieser Test ist im gesamten Diff des Slices unberührt | — |
| `git status`/`git diff` (am Ende dieses Laufs) | sauber, keine Restspur durch die Verifikation selbst (`tools/schema/plan.yaml`, das Rollout-Artefakt jedes `test-integration`-Laufs, jeweils per `git checkout --` zurückgesetzt) | — |

## Prüfpunkt 1 — DoD Punkt für Punkt gegen tatsächlichen Code

| # | DoD-Punkt | Verdikt | Beleg |
|---|---|---|---|
| 1 | Decoder trägt Spalten-Typ-OID (`decode.Column` erweitert), unit-getestet | **bestätigt** | `decode.go`: `Column{Name, Key, TypeOID}`, `relationColumns` liest `source.DataType` in `TypeOID` — eigenständig gelesen (Zeilen 64–68, 239–250). `decode_test.go`: `TestDecodeRelationColumnTypeOID` prüft zwei Spalten mit unterschiedlichen OIDs (23/25) auf Unterscheidbarkeit nach der Dekodierung — eigenständig gelesen, Test liegt vor. `make test` in diesem Lauf grün, Paket `decode` explizit `ok`. |
| 2 | `Assembler.Consume` behandelt `*decode.Relation` real (unverändert → No-op, kompatible Erweiterung → neue `SchemaVersionID`, `TableBinding` aktualisiert), unit-getestet | **bestätigt** | `mapper.go`: `observeRelation` (Zeilen 303–347), `classifyRelationColumns` (250–269) — eigenständig nachvollzogen: Namens→OID-Abgleich über `known` läuft vor jedem Längenvergleich, fehlende/abweichende bekannte Spalte liefert sofort `relationOther`, sonst entscheidet die Spaltenzahl zwischen `relationUnchanged`/`relationCompatibleExtension`. Fünf Tests in `mapper_test.go` gelesen (`TestConsumeRelationUnchanged`, `TestConsumeRelationCompatibleExtension`, `TestConsumeRelationOtherChangeStaysConservative`, `TestConsumeRelationBackfillsMissingTableSchema`, `TestConsumeRelationNilSchemaStoreStaysNoop`) — decken die drei Klassifikationsfälle, den Backfill-Pfad und den `nil`-Store-Default ab. `make test` grün, Paket `mapper` explizit `ok`. |
| 3 | `LH-FA-SCH-005` real geschlossen: `TestMVPSchemaChangeAddColumn` läuft ohne Konzeptänderung grün, unterscheidbare `schema_version`-Werte vor/nach `ALTER TABLE ADD COLUMN`; `make gates`/`make test-integration` dreimal grün | **bestätigt, eigenständig dreimal reproduziert** | `git show bc40556 -- test/integration/integration_test.go` eigenständig gelesen: einzige inhaltliche Änderung ist die letzte Assertion, invertiert von „Schema-Version gleich" (Tatsachenbeleg des bis dahin unimplementierten Zustands) zu „Schema-Version unterscheidet sich" (Erwartung, `LH-FA-SCH-005` Boundary) — Testaufbau (INSERT vor Änderung, `ALTER TABLE ADD COLUMN`, INSERT danach, Row-Image- und Lesbarkeits-Prüfungen) unverändert, keine Konzeptänderung. In drei eigenständigen, unabhängigen `make test-integration`-Läufen dieses Verifikations-Laufs **PASS** in allen drei Durchläufen. `make gates` in diesem Lauf ebenfalls grün (Tabelle oben). |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **bestätigt** (Checkbox im Plan noch offen — korrekt für diesen Lauf-Zeitpunkt) | `docs/reviews/review-slice-032.md` liegt vor und ist committet (`aa283ec`), 0 HIGH/MEDIUM, 1 LOW (F-1), kein Merge-Blocker. Kein Self-Review — der Review-Report trägt einen eigenen Skill-Kopf und prüft unabhängig gegen Plan/ADR/Hard Rules. |
| 5 | Doku-Update, falls öffentlicher Vertrag berührt | **bestätigt** (kein Update nötig, Begründung trägt) | `git diff 631045d..aa283ec --stat -- compose.yaml spec/ harness/ docs/plan/adr/` liefert nur `harness/image-hash.txt` (erwarteter Digest-Nachzug bei geändertem Build-Kontext, kein öffentlicher Vertrag). `CDC_TABLES`-Format, `SchemaStorePort`-Signatur selbst und `ADR-0015` bleiben unverändert; `Config.SchemaStore`/`Consume(ctx, …)` sind interne Adapter-/Assembler-Verdrahtung — eigenständig anhand des Diffs nachvollzogen. |
| 6–10 | Closure-Notiz, Reconciliation-Register (entfällt, Greenfield), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **korrekt offen** | §6 trägt beide Risiken noch mit `<bei Closure einzutragen>`, §7 ist vollständig Platzhalter. Das ist der erwartete Zustand an dieser Stelle der Rollen-Sequenz (Implementer→Reviewer→Verifier, **vor** Planner-Closure, Modul 8): Diese Punkte sind Planner-Arbeit bei der Closure nach `done/`, nicht Teil der Implementer→Verifier-Übergabe. Kein DoD-Defekt. |

**Zwischenstand:** Alle drei vom Implementer bereits als erledigt
markierten technischen DoD-Punkte (1–3) sind eigenständig vollständig
nachgeprüft und **bestätigt**, einschließlich beider namentlich
verlangten Testfälle in drei unabhängigen Reproduktionen. Punkte 4 und 5
sind ebenfalls bestätigt. Die fünf Closure-Pflichten (6–10) stehen
konsistent mit dem aktuellen Lifecycle-Zeitpunkt (noch nicht per Planner
geschlossen) offen — keine dieser Lücken ist eine DoD-**Verletzung**,
sondern der korrekte Zwischenstand vor Closure.

## Prüfpunkt 2 — `LH-FA-SCH-005`s Boundary-Kriterium: real oder nur behauptet?

**Real erfüllt, eigenständig verifiziert.**

`LH-FA-SCH-005` Boundary (`spec/lastenheft.md`): „Given zwei Changes vor
und nach einer Schemaänderung, when ihre Schema-Versionen gelesen werden,
then lassen sie sich unterscheiden." Vor diesem Slice war die
`TableBinding`-Schema-Version statisch bei der Aktivierung gebunden und
lief über die gesamte Container-Laufzeit unverändert — `mapper.go`s
Default-Zweig verwarf jede `*decode.Relation`-Nachricht (Review-Report
`review-slice-032.md`, Diff-Kommentar in `git show bc40556` zitiert dies
wörtlich als vorherigen Zustand). Mit diesem Slice:

1. **Code-Pfad real durchlaufen:** `observeRelation` liest über
   `SchemaStorePort.CurrentVersion`/`TableSchema` die bekannte Spaltenform,
   klassifiziert die eingehende Relation-Nachricht, und registriert bei
   `relationCompatibleExtension` real eine neue `SchemaVersionID`
   (`nextSchemaVersionID`) über `RegisterVersion` — anschließend hebt
   `a.tables[...]` die `TableBinding` auf die neue Version. Eigenständig
   gelesen, `mapper.go:330–346`.
2. **Rollen-Grants real vorhanden:** `tools/schema/nacharbeit-roles.sql`
   trägt `GRANT INSERT ON cdc.schema_version TO cdc_capture` und
   `GRANT SELECT, INSERT ON cdc.table_schema TO cdc_capture` — ohne diese
   Grants würde `RegisterVersion` an der Datenbank scheitern, nicht nur
   im Unit-Test. Eigenständig gelesen.
3. **End-to-End real bestätigt:** `TestMVPSchemaChangeAddColumn` führt
   real ein `ALTER TABLE … ADD COLUMN` gegen eine laufende
   PostgreSQL-Instanz im Testcontainer aus, liest beide Changes über
   `cdc.changes` zurück und vergleicht `schema_version`. In drei
   eigenständigen Läufen dieser Verifikation: **PASS** in allen drei —
   die beiden Werte sind tatsächlich unterscheidbar (Assertion würde bei
   Gleichheit `t.Fatalf` auslösen).

Das Boundary-Kriterium ist damit nicht nur im Code plausibel, sondern
real gegen eine laufende Instanz demonstriert — dreifach reproduziert,
nicht einmalig zufällig grün.

## Prüfpunkt 3 — Trägt F-1 aus dem Review-Report tatsächlich keinen
Merge-Blocker?

**Ja, eigenständig am Code nachvollzogen, nicht aus dem Review-Report
übernommen.**

`classifyRelationColumns` (`mapper.go:250–269`):

```go
for _, column := range known {
    oid, present := incomingByName[column.Name]
    if !present || oid != uint32(column.OID) {
        return relationOther
    }
}
```

- **Spalte entfernt:** eine bekannte Spalte fehlt in `incomingByName` →
  `present == false` → `return relationOther`.
- **Spalte umbenannt:** dieselbe Bedingung — die bekannte Spalte taucht
  unter ihrem alten Namen nicht mehr auf → `present == false` →
  `return relationOther`.
- **Spaltentyp geändert** (der einzige per Test belegte Fall,
  `TestConsumeRelationOtherChangeStaysConservative`): `present == true`,
  aber `oid != uint32(column.OID)` → `return relationOther`.

Alle drei Unterfälle enden in derselben Zeile (`return relationOther`)
derselben Schleifeniteration — kein bedingter Zweig danach unterscheidet
sie weiter. Der Testfall für den dritten Unterfall deckt damit
tatsächlich denselben Ausführungspfad ab wie die beiden ungetesteten. Die
Einstufung des Reviewers als LOW/Testabdeckungs-Hinweis ohne semantische
Auswirkung ist bei eigenständiger Prüfung korrekt — kein Verifikations-
Widerspruch.

## Negativbefunde

- geprüft, ohne Befund: **Kein ADR-Inhalt geändert** — `git show --stat`
  über alle drei slice-032-Commits: keiner berührt
  `docs/plan/adr/0015-schema-evolution.md`.
- geprüft, ohne Befund: **Kein Overclaiming ggü. `LH-FA-SCH-004`.**
  Slice-Kopf, §1 und DoD-Punkt 3 benennen ausdrücklich, dass
  `LH-FA-SCH-004`s Negative-Fall (Fehlerklasse `schema` bei inkompatibler
  Änderung) nicht in diesem Slice geschlossen wird — Folge-Slice
  `slice-033`.
- geprüft, ohne Befund: **`TestMVPSchemaChangeIncompatibleTypeChange`
  unverändert und nicht durch diesen Slice beschädigt.** `git show bc40556
  -- test/integration/integration_test.go` berührt diese Funktion
  nachweislich nicht (0 Treffer bei `grep -c`); in allen drei eigenständig
  gestarteten `make test-integration`-Läufen PASS.
- geprüft, ohne Befund: **Wiring-Konsistenz.** `wiring.go` legt einen
  eigenen `SchemaStore`-Pool (`postgresstorage.NewSchemaStore`) mit
  eigenem `defer schemaStore.Close()` an, gebunden an dieselbe Rolle wie
  Store/Stream (`cdc_capture`) — kein geteilter Pool, konsistent mit dem
  Plan-Nachzug.
- geprüft, ohne Befund: **Traceability** — alle drei Commit-Betreffs
  tragen `LH-FA-SCH-005`/`ADR-0015`, kein `SPEC-*`/`ARC-*` im Betreff;
  `make commit-traceability` (Teil von `make gates`) in diesem Lauf grün.
- geprüft, ohne Befund: **`git status`/`git diff`** am Ende dieses Laufs
  — sauber, keine Restspur durch die Verifikation selbst.

## Eigene Befunde

Keine.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Zusammenfassung DoD:** 5/5 als „erledigt" markierte DoD-Punkte (1–5)
eigenständig vollständig nachgeprüft und bestanden, einschließlich beider
namentlich verlangten Testfälle in drei unabhängigen Reproduktionen. 5/10
Closure-Pflichten (6–10) stehen korrekt offen — sie sind Planner-Arbeit
bei der Closure, kein Verifikations-Defekt an dieser Stelle der
Rollen-Sequenz.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt.** Alle fünf als erledigt
markierten DoD-Punkte sind real erfüllt, nicht nur behauptet — der
Decoder trägt die Typ-OID, der Mapper klassifiziert und re-versioniert
real über den `SchemaStorePort`, und `LH-FA-SCH-005`s Boundary-Kriterium
ist gegen eine laufende PostgreSQL-Instanz dreifach reproduziert real
erfüllt, nicht nur im Unit-Test plausibel. `ADR-0015` bleibt `Accepted`
und unverändert; `LH-FA-SCH-004` wird nirgends fälschlich als
mitgeschlossen dargestellt. Der einzige Review-Fund (F-1, LOW) ist bei
eigenständiger Code-Prüfung tatsächlich ohne semantische Auswirkung — die
beiden ungetesteten Unterfälle laufen nachweislich über denselben
Return-Pfad wie der getestete dritte.

**Plan-vs-Code-Diff:** Keine Abweichung zwischen §3-Plan-Nachzug und
tatsächlichem Code gefunden — `Consume`-/`NewAssembler`-Signaturen,
Vergleichslogik, Versions-ID-Format, Backfill-Verhalten und
`wiring.go`-Verdrahtung sind exakt so umgesetzt, wie der Plan-Nachzug es
dokumentiert.

**Closure-Bereitschaft:** Aus Verifikations-Sicht steht der
Fähigkeits-Teil (DoD 1–5) closure-bereit. Vor dem `git mv` nach `done/`
fehlen noch die fünf Planner-Pflichten (§6 Risiko-Ausgänge, §7
Closure-Notiz mit Steering-Loop-Lerneintrag — F-1 aus
`review-slice-032.md` gehört dort als Beobachtung hinein — und die
Register-/Paarungs-Prüfung), die dieser Bericht bewusst nicht vorwegnimmt.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert (`git status`/`git diff` am
Ende sauber).

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: v6.5.0, 54 Dateien OK; `d-check`
Standardlauf: 273 Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD`
Modul `commits`: 273 Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5
Commits; `a-check`: 0 Befunde). `make test` einmal ausgeführt (voller
Lauf), Exit 0, alle Pakete `ok`. `make test-integration` **dreimal**
eigenständig und unabhängig voneinander ausgeführt, Exit 0 in allen drei
Läufen — `TestMVPSchemaChangeAddColumn` PASS in allen drei
(unterscheidbare `schema_version`-Werte real bestätigt),
`TestMVPSchemaChangeIncompatibleTypeChange` PASS in allen drei
(Verhalten unverändert, Diff berührt diese Funktion nachweislich nicht).
`git status`/`git diff` am Ende dieses Laufs sauber (keine
Arbeitsverzeichnis-Änderung durch die Verifikation selbst).
