# Review-Report: slice-032 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-032`, §1/§2/§3) und
`ADR-0015` (`Accepted`, `permanent`) sowie `AGENTS.md` §3 Hard Rules
(Modul 10 §Drei Review-Arten).

**Gegenstand:** Commits `bc40556` (Decoder-Erweiterung, Mapper-Klassifikation
und -Verdrahtung, Rollen-/Testschema-Lücken, `harness/image-hash.txt`),
`5882056` (DoD-Häkchen, Plan-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-032-dynamische-re-versionierung.md` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug, §6 Risiken, §8 Register-Sichtung)
- `docs/plan/adr/0015-schema-evolution.md` (Accepted, `permanent`)
- `docs/reviews/review-slice-031.md` (F-1/F-2/F-3 — dieselben Kommentar-/API-Klassen als Prüfmaßstab für Wiederholung)
- `spec/lastenheft.md` (`LH-FA-SCH-004`, `LH-FA-SCH-005`), `spec/pflichtenheft.md` (`SPEC-004`)
- `internal/adapters/driving/replication/decode/decode.go`, `decode_test.go`
- `internal/adapters/driving/replication/mapper/mapper.go`, `mapper_test.go`
- `internal/adapters/driving/replication/receive/receive.go`
- `internal/bootstrap/wiring.go`, `internal/bootstrap/walretention_endtoend_test.go`
- `internal/adapters/driven/postgresstorage/schemastore.go` (F-3-Fix aus `100ff2b`, Grundlage für den Backfill-Pfad dieses Slices)
- `test/integration/integration_test.go` (`TestMVPSchemaChangeAddColumn`, `TestMVPSchemaChangeIncompatibleTypeChange`)
- `tools/schema/nacharbeit-roles.sql`
- `AGENTS.md` §3 Hard Rules, insbesondere §3.7
- Memory-Notiz „Keine Slice-Chronik in schema.yaml/Quellcode" (Nutzer-Korrektur slice-018, geschärft für slice-031 F-1)

---

## Findings

Keine HIGH/MEDIUM/LOW-Findings mit Merge-Relevanz. Ein LOW-Hinweis zur
Testabdeckung.

### F-1 — Zwei der drei benannten `relationOther`-Unterfälle sind nur über den Code-Pfad, nicht über einen eigenen Test belegt

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `internal/adapters/driving/replication/mapper/mapper_test.go:452-479` (`TestConsumeRelationOtherChangeStaysConservative`)
- `befund`: §1/§3 des Slice-Plans und der Doc-Kommentar auf `relationOther`
  (`mapper.go:239-244`) benennen explizit drei Unterfälle, die konservativ
  bleiben müssen: Spalte entfernt, Spaltentyp geändert, Spalte umbenannt.
  Der einzige Test für diesen Zweig
  (`TestConsumeRelationOtherChangeStaysConservative`) deckt ausschließlich
  den Typ-Änderungs-Fall ab (`amount` behält den Namen, OID wechselt von
  25 auf 23). Spalte-entfernt und Spalte-umbenannt laufen zwar
  nachweislich über denselben Code-Pfad (`classifyRelationColumns`,
  `mapper.go:283-286`: eine bekannte Spalte, die in der eingehenden
  Relation nicht unter ihrem Namen auftaucht, liefert unabhängig vom
  Grund sofort `relationOther`) — das wurde in dieser Review-Sitzung durch
  Code-Lektüre nachvollzogen, nicht aber durch einen eigenen Testfall im
  Diff belegt.
- `verifizierbar`: ja — `go test ./internal/adapters/driving/replication/mapper/... -run TestConsumeRelation -v` zeigt fünf Tests, keiner davon variiert den Fall „Spaltenzahl gleich, aber eine bekannte Spalte fehlt unter jedem Namen" oder „Spalte umbenannt, Namen tauschen".
- `klasse`: „Testfall pro benanntem Unterfall fehlt bei mehrdeutigem Zweig"

## Negativbefunde

- geprüft, ohne Befund: **Klassifikationslogik konservativ für alle drei
  benannten `default`-Fälle.** Eigenständige Code-Lektüre von
  `classifyRelationColumns` (`mapper.go:265-286`): Der Namens→OID-Abgleich
  über `known` läuft vor jedem Längenvergleich; fehlt eine bekannte Spalte
  unter ihrem Namen (Entfernung, Umbenennung) oder weicht ihre OID ab
  (Typänderung), kehrt die Funktion sofort mit `relationOther` zurück,
  unabhängig von der resultierenden Spaltenzahl. Das
  Risiko aus §6 („Spalte entfernt und mit gleichem Namen, anderem Typ
  wieder hinzugefügt") fällt ebenfalls in diesen Zweig (OID-Mismatch unter
  gleichem Namen) — bereits durch das bestehende Design konservativ
  behandelt, nicht nur durch Zufall. `observeRelation` (`mapper.go:311-345`)
  schreibt bei `relationOther` nicht in den `SchemaStorePort` und meldet
  keinen Fehler — wie in §1 verlangt.
- geprüft, ohne Befund: **`TestMVPSchemaChangeIncompatibleTypeChange`
  bleibt unverändert im Verhalten.** Diff (`git show bc40556 --
  test/integration/integration_test.go`) berührt diese Testfunktion nicht.
  Eigenständig durchdacht und durch realen Lauf bestätigt: Fall 2 dieses
  Tests (`ALTER COLUMN amount TYPE integer` mit konvertierbaren
  Bestandsdaten) erzeugt eine neue Relation-Nachricht mit geänderter
  `amount`-OID — `classifyRelationColumns` liefert dafür `relationOther`
  (OID-Mismatch bei gleichem Namen), `observeRelation` bleibt wirkungslos;
  der Test prüft ohnehin nur das Row-Image, nicht die Schema-Version, und
  bestand real dreimal in Folge (siehe Testläufe unten).
- geprüft, ohne Befund: **Wiring-Änderung ohne Nebenwirkung auf den
  bestehenden Erfassungspfad.** Alle Aufrufstellen von `Consume` und
  `NewAssembler` sind konsistent auf die neuen Signaturen migriert
  (`grep -rn "\.Consume(\|NewAssembler("`, neun Konsumenten in
  `receive.go`, `decode_test.go`, `mapper_test.go`) — kein
  Alt-Aufrufer ohne `ctx`/`schemaStore`-Parameter übrig.
  `receive.go:process` reicht den bereits vorhandenen `ctx` nur durch
  (kein neuer Kontext-Ursprung); `wiring.go` legt einen eigenen Pool
  (`postgresstorage.NewSchemaStore`) mit eigenem `defer …Close()` an,
  analog zu `store`/`activation`/`heartbeat` — kein geteilter Pool, keine
  fehlende Freigabe.
- geprüft, ohne Befund: **Rollen-Grants konsistent mit dem
  Least-Privilege-Muster.** `GRANT INSERT ON cdc.schema_version TO
  cdc_capture` ergänzt exakt die fehlende Hälfte (bestehendes `SELECT` war
  bereits vorhanden, Zeile 57); `GRANT SELECT, INSERT ON
  cdc.table_schema TO cdc_capture` trägt genau die beiden vom
  Erfassungspfad benötigten Operationen (`CurrentVersion`/`TableSchema`
  lesen, `RegisterVersion` schreibt über `ON CONFLICT DO NOTHING` ohne
  `UPDATE`-Bedarf) — kein `UPDATE`/`DELETE` gegrantet, wo keins gebraucht
  wird. Kein Neu-Zuschnitt der bestehenden `cdc_admin`/`cdc_reader`-Grants.
- geprüft, ohne Befund: **`walretention_endtoend_test.go`-DDL vollständig
  und formkonsistent.** Die inline `CREATE TABLE cdc.table_schema`-DDL
  (Spalten, PK, UNIQUE, FK) stimmt exakt mit `tools/schema/schema.yaml`s
  `table_schema`-Knoten überein (`schema_version_id text` mit FK auf
  `schema_version(schema_version_id)`, `ordinal_position bigint CHECK >=
  1`, `column_name text`, `column_oid bigint`, PK
  `(schema_version_id, ordinal_position)`, UNIQUE
  `(schema_version_id, column_name)`) — keine Abweichung, die die beiden
  DDL-Quellen später divergieren ließe.
- geprüft, ohne Befund: **F-3 aus `review-slice-031.md` ist Grundlage,
  nicht offenes Risiko.** `RegisterVersion` (`schemastore.go:111-148`)
  wurde in `100ff2b` bereits vor diesem Slice auf `CountTableSchemaColumns`
  statt auf den `registered`-Rückgabewert der Versions-Zeile umgestellt —
  der Backfill-Pfad dieses Slices (`observeRelation`, `ErrSchemaVersionUnknown`
  → `RegisterVersion` mit unveränderter `SchemaVersionID`) schreibt die
  Spaltenform deshalb tatsächlich, unabhängig davon, ob die
  Schema-Version-Zeile bereits vor diesem Aufruf bestand. Real bestätigt
  durch `TestConsumeRelationBackfillsMissingTableSchema` (`registrations
  == 1`) und durch den grünen `TestMVPSchemaChangeAddColumn`-Lauf.
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7), dritte
  Gelegenheit nach slice-018/slice-031-F-1.** Alle neuen/geänderten
  Kommentare in `decode.go`, `mapper.go`, `receive.go`, `wiring.go`
  beschreiben den geltenden Zustand im Indikativ und referenzieren
  ausschließlich auflösbare Anker (`ADR-0015`, `LH-FA-SCH-004.a/005`,
  `SPEC-004`) — kein Vorwärtsverweis auf eine konkrete Folge-Slice-Kennung
  (`slice-033`) im Produktionscode; die einzige `slice-033`-Nennung im
  gesamten Diff steht im Slice-Plan selbst (§1, §Welle-Feld), wo eine
  Kennung Pflicht ist. `grep -n "slice-0[0-9][0-9]" --include=*.go -r
  internal/` zeigt für diesen Diff keinen neuen Treffer (nur den bereits
  vor `slice-031` bestehenden `wiring.go:463`, außerhalb dieses Diffs).
  Kein Kommentar bricht mitten im Satz ab oder beschreibt eine verworfene
  Alternative.
- geprüft, ohne Befund: **DoD-Aktualisierung ohne Overclaiming.**
  `LH-FA-SCH-004` wird nirgends als geschlossen behauptet; Slice-Kopf und
  DoD-Punkt 3 benennen `LH-FA-SCH-004`s Negative-Fall ausdrücklich als
  `slice-033`-Aufgabe. Die drei letzten DoD-Punkte (Closure-Notiz,
  Risiko-Ausgänge, Paarungen) sind bewusst offen gelassen
  (Plan-Nachzug-Commit-Message: „Planner-Closure") — konsistent mit dem
  aktuellen Lifecycle-Stand (`in-progress/`).
- geprüft, ohne Befund: **`tools/schema/nacharbeit-roles.sql`-Kommentar
  sachlich korrekt.** Die Begründung „`ON CONFLICT DO NOTHING` verlangt
  kein zusätzliches `SELECT`, anders als der `ON-CONFLICT-DO-UPDATE`-Zweig
  oben" ist konsistent mit `queries.InsertSchemaVersion`/
  `InsertTableSchemaColumn` (beide `ON CONFLICT DO NOTHING`, siehe
  `queries.go:78-81`, `:196-201`) und mit dem bereits verifizierten
  Gegenbeispiel `process_heartbeat` (`ON CONFLICT … DO UPDATE`, empirisch
  mit `SELECT`-Bedarf belegt).
- geprüft, ohne Befund: **Traceability.** Beide Commit-Betreffs tragen
  `LH-FA-SCH-005`/`ADR-0015`, kein `SPEC-*`/`ARC-*` im Betreff; `make
  commit-traceability` lief in dieser Sitzung über `HEAD~5..HEAD` grün (0
  Befunde, „Betreffs ohne Struktur-ID").
- geprüft, ohne Befund: **Hard Rule 3.1** (Docker-only) — alle Läufe über
  `make`/Docker, kein lokales Toolchain-Install. **3.2**
  (Suppression-Verbot) — kein `#noqa`/`//nolint`/`[SuppressMessage]` im
  Diff. **3.5** (ADR-Immutabilität) — `docs/plan/adr/0015-schema-evolution.md`
  im Diff unverändert. **3.6** (Gate-Lockerung) — keine Gate-Schwelle
  berührt. **3.3** (git-mv/Inhalt-Trennung) — nicht einschlägig, keine
  Umbenennung in diesem Diff.
- geprüft, ohne Befund: **`harness/image-hash.txt`-Nachzug korrekt
  ausgelöst.** Der Diff ändert Build-Kontext-Dateien (`internal/…`); der
  neue Digest (`sha256:b81a6407…`) ist Teil desselben Commits, konsistent
  mit `harness/README.md` §Werkzeuge (`make image`).
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf) —
  `baseline-verify` (54 Dateien), `docs-check` (272 Dateien, 0 Befunde),
  `commit-traceability` (5 Commits, `HEAD~5..HEAD`, OK), `a-check` (0
  Befunde) — alle grün.
- geprüft, ohne Befund: **`make test`** (voller Lauf, alle Pakete) — grün,
  inkl. `mapper`, `decode`, `receive`, `bootstrap`.
- geprüft, ohne Befund: **`make test-store`** (voller Lauf, dieser
  Review-Sitzung) — grün, `postgresstorage` real gegen PostgreSQL im
  Testcontainer, Schema über d-migrate ausgerollt.
- geprüft, ohne Befund: **`make test-replication`** (voller Lauf, dieser
  Review-Sitzung) — grün, `internal/bootstrap` (18,7s) inkl.
  `TestWALRetentionThresholdEndToEnd` mit der neu ergänzten
  `cdc.table_schema`-Inline-DDL.
- geprüft, ohne Befund: **`make test-integration`, dreimal in Folge
  grün** (dieser Review-Sitzung, wie im Slice-DoD verlangt) —
  `TestMVPSchemaChangeAddColumn` PASS in allen drei Läufen (belegt
  unterscheidbare `schema_version`-Werte vor/nach `ALTER TABLE ADD
  COLUMN`), `TestMVPSchemaChangeIncompatibleTypeChange` PASS in allen drei
  Läufen (Verhalten unverändert: PostgreSQL lehnt Fall 1 selbst ab, Fall 2
  übernimmt den konvertierten Wert ohne Fehlerklasse `schema`). Nach jedem
  Lauf blieb außer `tools/schema/plan.yaml` (Rollout-Artefakt mit
  wechselndem Ziel-Hostnamen, per `git checkout --` zurückgesetzt) kein
  Arbeitsbaum-Rest; `git status --porcelain` am Ende dieser Sitzung leer.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Testfall pro benanntem Unterfall fehlt
bei mehrdeutigem Zweig

## Verdikt

**Merge-blockierend:** keins. F-1 (LOW) ist eine
Testabdeckungs-Empfehlung ohne semantische Auswirkung — die Korrektheit
des `relationOther`-Zweigs für alle drei benannten Unterfälle wurde in
dieser Sitzung durch Code-Lektüre unabhängig nachvollzogen (siehe
Negativbefund oben) und ist keine offene Frage, nur keine per Test
demonstrierte für zwei der drei Unterfälle.

**Ausdrücklich unabhängig geprüft, nicht vom Implementer-Bericht
übernommen:** die Korrektheit von `classifyRelationColumns` für alle drei
`relationOther`-Unterfälle (Code-Lektüre), die Unveränderlichkeit von
`TestMVPSchemaChangeIncompatibleTypeChange`s Verhalten (Code-Lektüre plus
realer dreifacher Testlauf), die Konsistenz aller `Consume`-/
`NewAssembler`-Aufrufstellen (grep über den vollständigen Baum), die
Least-Privilege-Korrektheit der neuen Grants, die Formgleichheit der
`walretention_endtoend_test.go`-Inline-DDL gegen `schema.yaml`, dass F-3
aus `review-slice-031.md` bereits vor diesem Diff behoben war, die
Kommentar-Disziplin nach `AGENTS.md` §3.7 über alle vier berührten
Go-Dateien, und `make gates`/`make test`/`make test-store`/
`make test-replication`/`make test-integration` (dreimal) als
eigenständige Läufe.

**Übergabe:** F-1 (LOW) geht als optionaler Hinweis an den Implementer
oder direkt in die Slice-Closure §7 als Beobachtung — kein
Rollen-Konflikt, keine Architect-Sequenz nötig (Modul 8). Dieser Report
ist ein Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen. Er
ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier
separat.
