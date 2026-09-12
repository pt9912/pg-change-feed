# Review-Report: slice-031 (Fixrunde) — 2026-09-12

**Review-Art:** Code — gezielte Fixrunden-Bestätigung gegen den
vorigen Review-Report `docs/reviews/review-slice-031.md` (F-1/F-2/F-3),
kein erneutes Vollreview des Slice.

**Gegenstand:** Commit `100ff2b`
(„fix(schema-store): Review-Fixrunde slice-031 F-1/F-2/F-3 (`ADR-0015`)"),
gegen `docs/reviews/review-slice-031.md` und `ADR-0015`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- `docs/reviews/review-slice-031.md` (F-1 HIGH, F-2 HIGH, F-3 MEDIUM)
- `git show 100ff2b` (Fix-Diff: 6 Dateien)
- `docs/plan/adr/0015-schema-evolution.md` (Accepted, `permanent`)
- `AGENTS.md` §3.7 (Kommentar-Klassen)

---

## F-1 — Vorwärtsverweise auf Folge-Slices in Go-Quellcode-Kommentaren

**Verdikt: behoben.**

`grep -rn "slice-0[0-9][0-9]" --include=*.go internal/` liefert keinen
Treffer mehr in den drei genannten Dateien
(`internal/domain/model/table_schema.go`,
`internal/application/port/outbound/schemastore.go`,
`internal/adapters/driven/postgresstorage/schemastore.go`). Ebenso
`grep -rn "docs/reviews" --include=*.go internal/` — leer; der
Review-Report-Pfad ist aus `outbound/schemastore.go:37` entfernt. Der
einzige verbleibende Treffer `internal/bootstrap/wiring.go:463`
(`slice-026`) ist außerhalb dieses Diffs, unverändert vorbestehend und
war bereits im Original-Report als „außerhalb dieses Diffs" vermerkt.

Die drei neuen Kommentarfassungen (`table_schema.go:6-8`,
`outbound/schemastore.go:35-42`, `postgresstorage/schemastore.go:23-26`)
beschreiben ausschließlich den Ist-Zustand — Adjektiv „ausschließlich
die Persistenz-Fähigkeit", keine Zeit- oder Chronik-Formulierung mehr.
Die Out-of-Scope-Information (dynamische Re-Versionierung, Typ-
Kompatibilitätsprüfung) ist als Sache genannt, ohne Slice-Kennung — das
entspricht `AGENTS.md` §3.7 (Indikativ über den Zustand).

## F-2 — Doc-Kommentar zu `CurrentVersion` behauptet ein empirisch falsches Sichtbarkeits-Verhalten

**Verdikt: behoben.**

Der Kommentar (`postgresstorage/schemastore.go:68-75`) behauptet jetzt:
Version 1 der statischen Erstaktivierung ist „bereits ohne einen
`RegisterVersion`-Aufruf über `CurrentVersion` sichtbar" — das Gegenteil
der vorigen (falschen) Fassung.

Eigenständig real gegen PostgreSQL reproduziert (Testcontainer,
Schema-Rollout über d-migrate, `CDC_STORE_TEST_DSN`): Testfall aktiviert
eine Tabelle über `TableActivationAdapter.Register` und ruft danach
`SchemaStorePort.CurrentVersion` auf, ohne je `RegisterVersion`
aufzurufen. Ausgabe: `ok=true,
version={ID:sv-review-f2-1 SourceTableID:tbl-review-f2 Version:1}`. Der
Kommentar beschreibt jetzt exakt dieses beobachtete Verhalten.

## F-3 — `RegisterVersion` kann die Spaltenform einer bereits bestehenden Schema-Version-Zeile nicht nachtragen

**Verdikt: behoben.**

Der Fix ersetzt die Backfill-Bedingung: statt an der Neuheit der
`schema_version`-Zeile (`RowsAffected() == 1`) hängt die
Spalten-Einfügung jetzt an der Abwesenheit einer bestehenden Spaltenform
selbst (`queries.CountTableSchemaColumns`, neue Query). Code gelesen
(`postgresstorage/schemastore.go:111-148`) — Transaktion: `InsertSchemaVersion`
(Fehler bricht sofort ab, unabhängig vom Ergebnis) → `CountTableSchemaColumns`
für die `SchemaVersionID` → Spalten nur bei `existingColumns == 0`
einfügen → `written` als Rückgabewert.

Eigenständig real gegen PostgreSQL reproduziert (nicht die simulierte
SQL-Zeile aus dem neuen Testfall, sondern die reale Aktivierung über
`TableActivationAdapter.Register`, um denselben Pfad wie in Produktion
zu treffen):

- Aktivierung (Version 1, keine Spaltenform) → `RegisterVersion` mit
  Spaltenform für dieselbe `SchemaVersionID` → `backfilled=true`,
  `TableSchema` liest die Spalte korrekt zurück. Backfill funktioniert.
- Erneuter `RegisterVersion`-Aufruf mit identischer Version/Schema →
  `repeat=false`, `cdc.table_schema`-Zeilenzahl für diese
  `SchemaVersionID` bleibt bei 1 (keine Duplizierung). **Idempotenz
  intakt** — das war der kritischste Punkt, da der F-3-Fix genau an der
  Stelle ansetzt, die vorher die Idempotenz trug.
- Mismatch-Pfad (`version.ID != schema.VersionID`) weiterhin abgelehnt
  (`ErrSchemaVersionMismatch`) — Prüfung steht unverändert vor dem ersten
  SQL-Aufruf.

Zusätzlich: `TestSchemaStoreRegisterVersionBackfillsColumns`
(`schemastore_test.go`) gelesen — deckt denselben Backfill- und
Idempotenz-Fall über eine direkt simulierte Bestandszeile ab und lief in
`make test-store` grün (s. u.).

## Negativbefunde (Regressionsprüfung)

- geprüft, ohne Befund: **Keine Chronik/Slice-Referenz in den neuen
  Kommentaren.** Alle sechs geänderten Kommentarblöcke (`table_schema.go`,
  `outbound/schemastore.go`, `postgresstorage/schemastore.go` ×2 Stellen)
  gelesen — Indikativ über den Ist-Zustand, kein Zeitbezug, keine
  Slice-/Review-Kennung.
- geprüft, ohne Befund: **Idempotenz von `RegisterVersion` nach dem
  F-3-Fix** — eigenständig real reproduziert (Aktivierung → Backfill →
  Wiederholung), keine doppelte Spaltenform, `repeat=false` (s. F-3 oben).
  Dies war der Regressions-Kandidat mit der höchsten Wahrscheinlichkeit,
  da der Fix exakt die vorher Idempotenz-tragende Bedingung ersetzt —
  kein Befund.
- geprüft, ohne Befund: **Mismatch-Fehlerpfad** — `RegisterVersion` mit
  widersprüchlicher `version.ID`/`schema.VersionID`-Kombination weiterhin
  über `ErrSchemaVersionMismatch` abgelehnt, vor jedem SQL-Aufruf; sowohl
  über den bestehenden Testfall
  `TestSchemaStoreRegisterVersionRejectsMismatch` als auch eigenständig
  gegen den Aktivierungspfad reproduziert.
- geprüft, ohne Befund: **`InsertSchemaVersion`-Fehlerbehandlung
  unverändert** — der Fix zieht `InsertSchemaVersion` aus der
  `tag, err :=`-Zuweisung heraus (kein `registered`-Gebrauch mehr), der
  Fehlerpfad selbst (`schemaStoreFailure` bei SQL-Fehler,
  `defer tx.Rollback`) bleibt unangetastet.
- geprüft, ohne Befund: **Zero-Column-Edge-Case.** `model.NewTableSchema`
  verlangt mindestens eine Spalte (`ErrEmptyColumns`) — der theoretische
  Fall „Backfill mit 0 Spalten meldet bei jedem Aufruf erneut `written=true`"
  ist über den Domain-Konstruktor ausgeschlossen, kein eigenständiges
  Finding.
- geprüft, ohne Befund: **Plan-Nachzug.** Der neue Abschnitt „Fixrunde
  nach Review" im Slice-Plan benennt alle drei Findings mit Kategorie,
  Ursache und Fix korrekt und referenziert den Review-Report-Pfad (im
  Plan zulässig, anders als im Quellcode — F-1 betraf nur
  Go-Kommentare/`schema.yaml`).
- geprüft, ohne Befund: **Traceability.** Commit-Betreff trägt
  `(ADR-0015)`, kein `SPEC-*`/`ARC-*`. `make gates` (inkl.
  `commit-traceability` über `HEAD~5..HEAD`) lief grün.
- geprüft, ohne Befund: **Arbeitsbaum nach Ad-hoc-Verifikation.** Die für
  F-2/F-3 genutzten temporären Testfunktionen wurden nach dem Lauf per
  `git checkout --` aus `schemastore_test.go` entfernt; ebenso das durch
  den eigenständigen Schema-Rollout überschriebene `tools/schema/plan.yaml`
  zurückgesetzt. `git status --porcelain` leer am Ende dieser Sitzung.

**Prüfmethodik F-2/F-3:** Beide real gegen PostgreSQL reproduziert, in
zwei Varianten — (a) `make test-store` (voller Lauf, alle Pakete inkl.
`postgresstorage`, grün) mit den unveränderten Testfällen des Commits;
(b) eigenständige, nach dem Lauf per `git checkout --` wieder entfernte
Testfunktionen direkt in `schemastore_test.go` (Dateiwahl bewusst: diese
Datei läuft laut Kopf-Kommentar vor den DROP-SCHEMA-Tests
`sqlviews_test.go`/`store_test.go`/`tableactivation_test.go`), die den
realen `TableActivationAdapter.Register`-Pfad statt der simulierten
SQL-Zeile aus dem neuen Testfall nutzen — verifiziert gegen einen
eigenständig aufgesetzten Testcontainer, `go test -run TestReviewFixrunde -v`,
beide `PASS`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Neue Findings dieser Fixrunde:** keine.

## Verdikt

**Merge-blockierend:** nein. F-1, F-2 und F-3 sind alle drei behoben;
keine Regression gefunden — Idempotenz und Mismatch-Fehlerpfad von
`RegisterVersion` bleiben intakt. `make gates` und `make test-store`
liefen in dieser Sitzung eigenständig grün.

**Übergabe:** Kein Rückgang an den Implementer nötig. Slice-Closure kann
mit dieser Bestätigung fortfahren; die drei ursprünglichen Finding-Klassen
aus `docs/reviews/review-slice-031.md` gehen wie vorgesehen in die
Slice-Closure §7 und von dort in den Zähler des Beobachtungs-Registers.
Dieser Report ist ein Lauf-Beleg und wird über Läufe hinweg nicht wieder
gelesen. Er ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat.
