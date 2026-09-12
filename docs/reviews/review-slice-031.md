# Review-Report: slice-031 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-031`, §1/§2/§3) und
`ADR-0015` (`docs/reviews/architect-verdict-slice-030-adr-0015.md`
§Umsetzungsskizze, Schritte 1–3) sowie `AGENTS.md` §3 Hard Rules (Modul 10
§Drei Review-Arten).

**Gegenstand:** Commits `ca86249` (`TableSchema`-Modell, `SchemaStorePort`,
Postgres-Adapter), `cebc273` (DoD-Häkchen, Plan-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-031-schema-persistenz-faehigkeit.md` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug, §6 Risiken, §8 Register-Sichtung)
- `docs/plan/adr/0015-schema-evolution.md` (Accepted, `permanent`)
- `docs/reviews/architect-verdict-slice-030-adr-0015.md` (§Umsetzungsskizze Schritte 1–3, Prüfmaßstab dieses Slices; Schritte 4–7 ausdrücklich nicht Gegenstand)
- `spec/lastenheft.md` (`LH-FA-SCH-004`, `LH-FA-SCH-005`), `spec/pflichtenheft.md` (`SPEC-004`)
- `internal/adapters/driving/replication/mapper/mapper.go`, `internal/application/usecase/enable/service.go`, `internal/bootstrap/wiring.go` (Abgrenzungs-Prüfung: laufender Erfassungspfad unberührt)
- `internal/adapters/driven/postgresstorage/tableactivation.go` (statische Erstaktivierung, Version 1)
- `tools/schema/schema.yaml`, `tools/schema/plan.yaml`, `tools/schema/down.sql`
- `docs/plan/planning/observations/BEO-PGC/schema-evolution-nicht-dynamisch/`
- `AGENTS.md` §3 Hard Rules, insbesondere §3.7
- Memory-Notiz „Keine Slice-Chronik in schema.yaml/Quellcode" (Nutzer-Korrektur slice-018, `AGENTS.md` §3.7 geschärft für Go-Quellcode/`schema.yaml`)

---

## Findings

### F-1 — Vorwärtsverweise auf Folge-Slices in Go-Quellcode-Kommentaren

- `kategorie`: HIGH
- `quelle`: Hard Rule `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was da
  ist" — Indikativ über den Zustand, keine Chronik/Zeitbezug), geschärft
  für Go-Quellcode und `schema.yaml` durch die dokumentierte
  Nutzer-Korrektur aus slice-018 (dieselbe Fehlerklasse: `"seit
  slice-018"`/`"nur noch"` in `queries.go`/`mapper.go`/`schema.yaml`
  wurde dort bereits korrigiert)
- `pfad`: `internal/domain/model/table_schema.go:10`;
  `internal/application/port/outbound/schemastore.go:37`, `:42-43`;
  `internal/adapters/driven/postgresstorage/schemastore.go:27`
- `befund`: Drei Go-Quellcode-Kommentare verweisen namentlich auf künftige
  Slices (`` `slice-032` ``, `` `slice-033` ``) statt ausschließlich den
  Ist-Zustand zu beschreiben: „Out-of-Scope dieses Slice, siehe
  `slice-033`" (`table_schema.go:10`), „Folgepflichten anderer Slices
  (`slice-032`, `slice-033`)" (`outbound/schemastore.go:42-43`, wortgleich
  in `postgresstorage/schemastore.go:27`). Der Port-Kommentar verweist
  zusätzlich auf einen Review-Report-Pfad
  (`docs/reviews/architect-verdict-slice-030-adr-0015.md`,
  `outbound/schemastore.go:37`) — ein Lauf-Beleg, der laut
  `.harness/skills/reviewer.md` „über Läufe hinweg nicht gelesen" wird und
  damit als Dauerverweis im Produktionscode ungeeignet ist. Die
  Out-of-Scope-/Folgepflicht-Information gehört bereits vollständig in
  §1/§3 des Slice-Plans (dort mit Kennung Pflicht, siehe Modul 5
  §Ziel-Form: Slice) — die Duplizierung in den Quellcode-Kommentar
  erzeugt genau das Drift-Risiko, das die geschärfte Regel benennt: wird
  `slice-033` später gesplittet, umnummeriert oder storniert, lügt der
  Code-Kommentar unbemerkt weiter, während der Slice-Plan über den
  Lifecycle korrekt bleibt.
- `verifizierbar`: ja — `grep -n "slice-0[0-9][0-9]" --include=*.go -r internal/` zeigt die drei Fundstellen (plus einen vorbestehenden Treffer `wiring.go:463`, außerhalb dieses Diffs).
- `klasse`: „Slice-Chronik/-Vorwärtsverweis in Go-Quellcode-Kommentar"

### F-2 — Doc-Kommentar zu `CurrentVersion` behauptet ein empirisch falsches Sichtbarkeits-Verhalten

- `kategorie`: HIGH
- `quelle`: Hard Rule `AGENTS.md` §3.7 (Kommentar muss den geltenden
  Zustand zutreffend beschreiben — Zusage-Klasse)
- `pfad`: `internal/adapters/driven/postgresstorage/schemastore.go:69-75`
- `befund`: Der Doc-Kommentar behauptet: „ihre [der statischen
  Erstaktivierung] Version 1 wird über `CurrentVersion` erst nach der
  ersten `RegisterVersion` sichtbar". Das ist empirisch falsch: `cdc.schema_version`
  ist eine einzige, von `TableActivationAdapter.Register` und
  `PostgresSchemaStoreAdapter.RegisterVersion` gemeinsam beschriebene
  Tabelle; `SelectCurrentSchemaVersion` selektiert unabhängig davon, wer
  die Zeile geschrieben hat. In einer eigenständigen Review-Probe (Fixtur:
  `TableActivationAdapter.Register` aktiviert eine Tabelle real gegen
  PostgreSQL, danach direkt `SchemaStorePort.CurrentVersion` ohne je
  `RegisterVersion` aufzurufen) meldet `CurrentVersion` sofort `ok=true,
  version=1` — nicht `ok=false` wie der Kommentar suggeriert. Das ist kein
  Fehler im aktuellen Verhalten (die Sichtbarkeit über eine gemeinsame
  Tabelle ist plausibel und für `slice-032` sogar hilfreich), aber der
  Kommentar beschreibt das Gegenteil dessen, was der Code tatsächlich tut
  — und genau dieser Kommentar ist der Ort, an dem `slice-032` das
  Sichtbarkeits-Verhalten des Ports nachschlagen wird.
- `verifizierbar`: ja — eigenständig reproduziert in dieser Review-Sitzung
  (ad-hoc Testdatei gegen reale PostgreSQL, danach entfernt, siehe
  Negativbefunde/Prüfmethodik unten); Ausgabe: „CurrentVersion nach reiner
  statischer Erstaktivierung (kein RegisterVersion-Aufruf): ok=true,
  version={ID:sv-review-vis-1 SourceTableID:tbl-review-vis Version:1}".
- `klasse`: „Kommentar beschreibt Gegenteil des tatsächlichen Verhaltens"

### F-3 — `RegisterVersion` kann die Spaltenform einer bereits bestehenden Schema-Version-Zeile nicht nachtragen (betrifft insbesondere Version 1)

- `kategorie`: MEDIUM
- `quelle`: Maintainability — Architektursoundness der als
  „Persistenz-Grundlage für die beiden Folge-Slices" deklarierten
  Fähigkeit (Slice-Ziel §1)
- `pfad`: `internal/adapters/driven/postgresstorage/schemastore.go:105-129`
  (`RegisterVersion`, insbesondere die `if registered { … }`-Bedingung um
  die Spalten-Einfügeschleife)
- `befund`: `RegisterVersion` schreibt die `cdc.table_schema`-Zeilen nur,
  wenn `INSERT … ON CONFLICT DO NOTHING` auf `cdc.schema_version`
  tatsächlich eine neue Zeile anlegt (`registered == true`). Für eine
  `SchemaVersionID`, die bereits existiert — was für Version 1 nach jeder
  `EnableTable`-Aktivierung der Fall ist, da `TableActivationAdapter.Register`
  dieselbe `cdc.schema_version`-Tabelle über dieselbe SQL-Konstante
  (`queries.InsertSchemaVersion`) beschreibt —, meldet `RegisterVersion`
  `registered=false` und überspringt die Spalten-Einfügung vollständig. Ein
  späterer Versuch, für exakt diese `SchemaVersionID` die Spaltenform
  nachzutragen (der naheliegende Weg für `slice-032`, Version 1s
  Spaltenform zu befüllen, da die statische Erstaktivierung selbst keine
  Spaltenform schreibt, §1 dieses Slices), bleibt wirkungslos:
  `TableSchema(v1)` liefert `ErrSchemaVersionUnknown` dauerhaft. In einer
  eigenständigen Review-Probe (reale PostgreSQL: `Register` für Version 1,
  danach `RegisterVersion` mit derselben ID und einer nichtleeren
  Spaltenform) meldete der Aufruf `registered=false` und
  `TableSchema(v1)` weiterhin `ErrSchemaVersionUnknown`. Das ist keine
  DoD-Verletzung dieses Slices (die Verdrahtung ist ausdrücklich
  ausgeschlossen), aber eine offene Lücke in der als „Grundlage" gedachten
  API, die `slice-032` entweder über einen anderen Schreibpfad lösen muss
  oder die hier nachgebessert werden sollte, bevor darauf gebaut wird.
- `verifizierbar`: ja — eigenständig reproduziert (ad-hoc Testdatei, siehe
  unten); Ausgabe: „RegisterVersion (Backfill-Versuch) meldet
  registered=false" / „TableSchema(v1) nach Backfill-Versuch:
  err=Schema-Version trägt keine TableSchema-Registrierung
  (ErrSchemaVersionUnknown=true)".
- `klasse`: „API erlaubt kein Nachtragen einer bereits bestehenden
  Versions-Zeile"

## Negativbefunde

- geprüft, ohne Befund: **Abgrenzung §1 — laufender Erfassungspfad
  unberührt.** `grep -rn "SchemaStorePort\|schemastore\|SchemaStore"` gegen
  `internal/adapters/driving/replication/mapper/mapper.go`,
  `internal/application/usecase/enable/service.go`,
  `internal/bootstrap/wiring.go` liefert keinen Treffer; `git diff
  ca86249~1 ca86249 -- <diese drei Dateien>` ist leer — keiner der drei
  Pfade wird vom Commit überhaupt berührt. Eigenständig geprüft, nicht aus
  dem Implementer-Bericht übernommen.
- geprüft, ohne Befund: **Wiring-Verzicht ehrlich dokumentiert.** Der
  Plan-Nachzug benennt explizit, dass die §3-Plan-Zeile zu `wiring.go`
  „nicht eingelöst" ist und begründet dies (keine Aufrufstelle, sonst
  ungenutzte DB-Verbindung) — keine stillschweigende Abweichung vom
  Vorab-Plan.
- geprüft, ohne Befund: **DoD-Aktualisierung ohne Overclaiming.** Weder
  im DoD noch im Plan-Nachzug wird behauptet, `LH-FA-SCH-004`/`005` seien
  bereits geschlossen; der Slice-Kopf benennt das ausdrücklich als Aufgabe
  von `slice-032`/`slice-033`.
- geprüft, ohne Befund: **`TableSchema`/`Column`-Form als Grundlage für
  `slice-033`.** `Column{Name, OID}` trägt die PostgreSQL-Typ-OID roh;
  eine OID identifiziert den exakten PostgreSQL-Datentyp eindeutig und ist
  ausreichendes Rohmaterial für eine spätere
  Typ-Kompatibilitätsprüfung (OID-Gleichheit/-Differenz), auch wenn die
  Übersetzung in eine Kompatibilitätsentscheidung bewusst nicht Teil
  dieses Modells ist (konsistent mit Architect-Verdikt Schritt 4).
- geprüft, ohne Befund: **Tabellenform `cdc.table_schema`.** PK
  `(schema_version_id, ordinal_position)`, UNIQUE
  `(schema_version_id, column_name)`, Fremdschlüssel auf
  `schema_version.schema_version_id` (Spaltenname stimmt mit der
  referenzierten Tabelle überein). `column_oid` als `biginteger`
  korrekt begründet: eine PostgreSQL-OID ist ein vorzeichenloses
  32-Bit-Feld (0…4.294.967.295) und passt nicht verlustfrei in ein
  vorzeichenbehaftetes `integer` (Bereich bis ca. 2,1 Mrd.). `tools/schema/plan.yaml`
  zeigt die erzeugte DDL mit exakt dieser Form
  (`CREATE TABLE "table_schema" ( … REFERENCES "schema_version"("schema_version_id") …
  "column_oid" BIGINT NOT NULL … )`).
- geprüft, ohne Befund: **Adapter-Test real gegen PostgreSQL, nicht nur
  Mocks.** `schemastore_test.go` nutzt dasselbe Skip-Pattern
  (`CDC_STORE_TEST_DSN`) wie die übrigen Store-Tests; kein Mock-Interface
  im Testfile. Eigenständig reproduziert in dieser Review-Sitzung: `make
  test-store` (voller Lauf, Testcontainer, Schema-Rollout über d-migrate)
  grün; zusätzlich gezielt `go test
  ./internal/adapters/driven/postgresstorage/... -run TestSchemaStore -v`
  gegen einen eigenständig aufgesetzten Testcontainer — alle 5 Tests
  (`TestSchemaStoreCurrentVersionMissing`,
  `TestSchemaStoreRegisterAndReadRoundTrip`,
  `TestSchemaStoreCurrentVersionHighest`,
  `TestSchemaStoreTableSchemaUnknown`,
  `TestSchemaStoreRegisterVersionRejectsMismatch`) PASS. Idempotenz
  (keine doppelte Spaltenform bei erneuter Registrierung) und die beiden
  Fehlerfälle (unbekannte Version, Versions-/Schema-ID-Mismatch) sind
  damit real belegt, nicht nur behauptet.
- geprüft, ohne Befund: **`InsertTableSchemaColumn`/`InsertSchemaVersion`
  Idempotenz-Konsistenz.** Beide Statements nutzen `ON CONFLICT DO
  NOTHING` auf dem jeweiligen Primärschlüssel; die
  `RowsAffected()==1`-Prüfung in `RegisterVersion` bildet exakt diese
  Semantik ab (siehe aber F-3 zur Kehrseite dieser Idempotenz-Wahl).
- geprüft, ohne Befund: **Traceability.** Beide Commit-Betreffs tragen
  `(ADR-0015)`, kein `SPEC-*`/`ARC-*` im Betreff; `make
  commit-traceability` lief in dieser Sitzung über `HEAD~5..HEAD` grün (0
  Befunde, „Betreffs ohne Struktur-ID").
- geprüft, ohne Befund: **Hard Rule 3.1** (Docker-only) — Testskip nur
  über Umgebungsvariable, kein lokales Toolchain-Install; alle Läufe über
  `make`/Docker. **3.2** (Suppression-Verbot) — kein `#noqa`/`//nolint`.
  **3.5** (ADR-Immutabilität) — `docs/plan/adr/0015-schema-evolution.md`
  im Diff unverändert. **3.6** (Gate-Lockerung) — keine Gate-Schwelle
  berührt. **3.3** (git-mv/Inhalt-Trennung) — nicht einschlägig, keine
  Umbenennung in diesem Diff.
- geprüft, ohne Befund: **Beobachtungs-Register-Zeiger aus §8.**
  `docs/plan/planning/observations/BEO-PGC/schema-evolution-nicht-dynamisch/`
  existiert wie im Slice-Kopf §8 zitiert.
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf) —
  `baseline-verify` (54 Dateien), `docs-check` (268 Dateien, 0 Befunde),
  `commit-traceability` (5 Commits, `HEAD~5..HEAD`, OK), `a-check` (0
  Befunde) — alle grün.
- geprüft, ohne Befund: **`make test-store`** (voller Lauf, dieser
  Review-Sitzung) — alle Pakete grün, inkl. `postgresstorage` (2,997s),
  real gegen PostgreSQL im Testcontainer, Schema über d-migrate
  ausgerollt.

**Prüfmethodik zu F-2/F-3:** Beide Befunde wurden durch zwei temporäre,
nach dem Lauf wieder entfernte Testdateien im Paket
`postgresstorage_test` gegen einen eigenständig aufgesetzten
PostgreSQL-Testcontainer verifiziert (nicht Teil des Slice-Diffs, nicht
committet). `tools/schema/plan.yaml` wurde durch die Schema-Rollouts
dieser Review-Sitzung mehrfach mit einem abweichenden `target`-Hostnamen
überschrieben und danach per `git checkout -- tools/schema/plan.yaml`
zurückgesetzt; der Arbeitsbaum ist am Ende dieser Review-Sitzung wieder
sauber (`git status --porcelain` leer).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Slice-Chronik/-Vorwärtsverweis in
Go-Quellcode-Kommentar · Kommentar beschreibt Gegenteil des tatsächlichen
Verhaltens · API erlaubt kein Nachtragen einer bereits bestehenden
Versions-Zeile

## Verdikt

**Merge-blockierend:** F-1 und F-2 (HIGH) — Kommentar-Klassen-Verstöße
gegen die Hard Rule `AGENTS.md` §3.7. F-1 wiederholt eine bereits einmal
korrigierte Fehlerklasse (Nutzer-Korrektur slice-018) an neuer Stelle;
F-2 ist eine sachlich falsche Verhaltensbeschreibung genau an der Stelle,
die `slice-032` als Referenz nutzen wird. Beide sind durch reine
Kommentar-Änderungen behebbar, ohne die Fähigkeit selbst (Modell, Port,
Adapter, Tabellenform, Tests) anzufassen — kein Rollen-Konflikt nach
Modul 8 (keine ADR-Frage, kein Plan-Widerspruch), daher genügt die direkte
Rückgabe an den Implementer ohne Architect-Sequenz.

F-3 (MEDIUM) ist kein Merge-Blocker für diesen Slice (Verdrahtung ist
§1-Ausschluss), gehört aber vor `slice-032` geklärt — entweder als
Risiko-Ergänzung in dessen Plan oder als kleine Nachbesserung an
`RegisterVersion` in einem eigenen Slice/Commit.

**Ausdrücklich unabhängig geprüft, nicht vom Implementer-Bericht
übernommen:** die Abgrenzungs-Behauptung „laufender Erfassungspfad
unberührt" (per Grep und Diff auf die drei genannten Dateien), die
Sichtbarkeits-Semantik von `CurrentVersion` (F-2, empirisch gegen reale
PostgreSQL), die Nachtragbarkeit der Spaltenform für Version 1 (F-3,
ebenso empirisch), die reale (nicht gemockte) Testausführung von
`schemastore_test.go`, die Korrektheit der `biginteger`-Begründung für
`column_oid`, und `make gates`/`make test-store` als eigenständige Läufe.

**Übergabe:** F-1/F-2 gehen an den Implementer zur Korrektur (Kommentare
präzisieren: Ist-Zustand statt Slice-Chronik, Sichtbarkeits-Beschreibung
an das tatsächliche Verhalten anpassen). F-3 geht als Risiko-Hinweis an
den Planner für `slice-032`. Die drei **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler des
Beobachtungs-Registers. Dieser Report ist ein Lauf-Beleg und wird über
Läufe hinweg nicht wieder gelesen. Er ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat.
