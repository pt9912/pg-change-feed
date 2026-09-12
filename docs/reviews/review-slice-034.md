# Review-Report: slice-034 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-034`, §1/§2/§3/§6) und
`AGENTS.md` §3 Hard Rules (Modul 10 §Drei Review-Arten). Kein aktives ADR
berührt (Slice-Kopf: „Kein aktives ADR wird geändert").

**Gegenstand:** Commits `f3a45f6` (neuer Testfall
`TestMVPActiveTablesViewMatchesActivationState`, `-run`-Muster-Nachzug in
`tools/harness/run-integration-tests.sh`), `6ef4aec` (DoD-Häkchen,
Plan-Nachzug §3).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-034-black-box-status-liste.md` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug, §6 Risiken, §8 Register-Sichtung)
- `spec/lastenheft.md` (`LH-FA-CFG-003`, `LH-FA-CFG-004`, `LH-FA-SST-002`, `LH-FA-SST-003`, `LH-FA-CFG-002`)
- `test/integration/integration_test.go` (vollständig, insbesondere `newMVPEnv`, `TestMVPActivationState`, `TestMVPDisableRetainedState`, `TestMVPActiveTablesViewMatchesActivationState`)
- `tools/harness/run-integration-tests.sh` (`-run`-Muster, `feed_mvp_idle`-Anlage, `CDC_TABLES`-Referenz)
- `compose.yaml` (`CDC_TABLES`-Bindung)
- `tools/schema/schema.yaml` (`views.active_tables`-Spaltensignatur und Kommentar zur bewussten aktiviert/retained-Nichtunterscheidung)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7
- `harness/conventions.md` (MR-000 ID-Schema, Modus-Deklaration `PGC`/GF)
- `docs/plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/`, `.../verwaltung-keine-sql-administration/` (im Slice-Plan §8 zitiert)
- `harness/README.md` §Werkzeuge (`make image`-Semantik), [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) (Accepted)

---

## Findings

Keine HIGH/MEDIUM/LOW-Findings mit Merge-Relevanz.

## Negativbefunde

- geprüft, ohne Befund: **Black-Box-Reinheit gegenüber der SQL-Sicht.**
  Beide neuen SQL-Lesungen (`test/integration/integration_test.go:571-578`
  Happy Path, `:592-599` Boundary) laufen als rohes `env.pool.QueryRow`
  gegen `cdc.active_tables`, ohne `status`/`list`-Paket-Beteiligung am
  Lesepfad selbst. Der `status.NewGetStatusService`-Aufruf
  (`:579-587`, `:601-608`) ist als Kontrollwert sauber *nach* der
  jeweiligen SQL-Lesung platziert und in eigenen Variablen
  (`enabled`, `idle`) gehalten — keine Vermischung der beiden Lesewege in
  einer Assertion. `newMVPEnv`s eigene SQL-Lesung gegen
  `cdc.source_table` (Zeile 110-114) ist Infrastruktur zur
  Tabellen-ID-Auflösung, existierte bereits vor diesem Diff und wird von
  keinem anderen bestehenden Testfall als Verstoß gegen Black-Box-Prinzip
  gewertet.
- geprüft, ohne Befund: **`feed_mvp_idle` ist eine korrekte
  Boundary-Wahl.** `tools/harness/run-integration-tests.sh:163` legt die
  Tabelle physisch an (`CREATE TABLE public.feed_mvp_idle`);
  `compose.yaml:68` bindet in `CDC_TABLES` ausschließlich
  `feed_mvp_flow`, `feed_mvp_full`, `feed_mvp_schema` — `feed_mvp_idle`
  wird nie per `EnableTable` aktiviert. Dieselbe Tabelle dient bereits
  `TestMVPActivationState` (Zeile 505-513, 538) als Boundary-Fall; die
  Wahl ist konsistent mit etabliertem Muster, keine Neuerfindung.
- geprüft, ohne Befund: **Testisolation gegen die übrigen
  `feed_mvp_*`-Aktivierungen im selben Compose-Lauf.** `grep -n
  "Disable\|feed_mvp_full"` über die gesamte Testdatei zeigt: nur
  `feed_mvp_flow` wird deaktiviert (`TestMVPDisableRetainedState`,
  Zeile 619ff.); `feed_mvp_full` wird an keiner Stelle deaktiviert. Die
  Happy-Path-Abfrage filtert zusätzlich hart auf
  `source_table_id = $1` (die Bindungs-ID von `feed_mvp_full` selbst),
  ist also auch unabhängig von der Zahl gleichzeitig aktivierter Zeilen
  eindeutig. Go führt Testfunktionen innerhalb eines Pakets in
  Quelltext-Reihenfolge aus (unabhängig von der Reihenfolge im
  `-run`-Alternationsmuster); der neue Testfall steht im Quelltext vor
  `TestMVPDisableRetainedState`, ist also auch bei einer künftigen
  Umsortierung des `-run`-Musters durch die Deaktivierung nicht
  gefährdet, solange die Quelltextreihenfolge erhalten bleibt — real
  bestätigt durch drei aufeinanderfolgende `make test-integration`-Läufe
  in dieser Review-Sitzung (siehe unten), alle mit
  `TestMVPActiveTablesViewMatchesActivationState` PASS vor
  `TestMVPDisableRetainedState`.
- geprüft, ohne Befund: **Mutationsprobe des Implementers plausibel,
  eigenständig durch Assertion-Lektüre nachvollzogen.** Beide
  Fatalf-Bedingungen (`activeCount != 1`, `idleCount != 0`) sind scharfe
  Punktvergleiche ohne Toleranzband; eine Drift zwischen SQL-Sicht und
  Use-Case (z. B. Sicht liefert eine Zeile zu viel oder zu wenig)
  schlägt exakt an der Grenze fehl, die die Boundary-Bedingung aus §1
  benennt. Die Spaltenfilterung (`source_table_id` bzw.
  `source_id`+`schema_name`+`table_name`) verhindert einen
  Falsch-Positiv-Grün-Zustand durch Zeilen anderer aktivierter Tabellen.
  Kein Git-Trace der Mutationsprobe selbst nötig (Aufgabenstellung); die
  Schärfe der Assertions selbst trägt die Aussage unabhängig vom
  Bericht.
- geprüft, ohne Befund: **DoD-Aktualisierung ohne Overclaiming.** Die vier
  in `6ef4aec` auf `[x]` gesetzten Punkte sind durch den Diff und durch
  eigene Testläufe in dieser Sitzung gedeckt (Testfall existiert und lief
  grün; SQL-Sicht wird real gegen `status.NewGetStatusService` gehalten;
  `make gates` und `make test-integration` liefen real, siehe unten; kein
  öffentlicher Vertrag berührt — reine Testdatei- und
  Runner-Skript-Änderung, keine API-/CLI-/Schema-Änderung). Die drei
  offen gelassenen Punkte (Review, Closure-Notiz, Register-/Risiko-/
  Paarungs-Punkte) sind korrekt unbeansprucht — sie sind Rollen nach dem
  Reviewer vorbehalten (Verifier/Planner-Closure), kein Self-Review nach
  Modul 8. Der Reconciliation-Punkt ist korrekt als „entfällt" markiert
  (Repo ist laut `harness/conventions.md` durchgehend GF, keine
  `reconciliation.md` vorhanden).
- geprüft, ohne Befund: **Plan-Nachzug §3 ehrlich benannt.** Die
  Ergänzung der Runner-Skript-Zeile als eigene §3-Zeile mit Verweis auf
  `BEO-PGC/test-runner-stiller-ausschluss` beschreibt zutreffend, dass die
  Änderung Teil des ursprünglichen Auftrags war, aber nicht als eigene
  Plan-Zeile geführt wurde — keine rückwirkende Umdeutung des Slice-Ziels.
- geprüft, ohne Befund: **Kein Image-Rebuild-Trigger.** Der Diff berührt
  ausschließlich `test/integration/integration_test.go` und
  `tools/harness/run-integration-tests.sh` — keine Datei unter `internal/`
  oder `cmd/`, also keinen Pfad, den `Dockerfile`s Build-Stage
  (`go build … ./cmd/pg-change-feed`) tatsächlich kompiliert. Der
  Laufzeit-Image-Layer (`COPY --from=build /out/pg-change-feed …`) kann
  sich bei unverändertem Binary-Inhalt nicht ändern; `harness/image-hash.txt`
  bleibt unangetastet — konsistent mit `ADR-0044`/`harness/README.md`
  §Werkzeuge.
- geprüft, ohne Befund: **AGENTS.md §3.7 (Kommentar-Disziplin), vierte
  Gelegenheit nach slice-031/slice-033.** Der neue Doc-Kommentar auf
  `TestMVPActiveTablesViewMatchesActivationState`
  (`integration_test.go:543-554`) beschreibt den geltenden Testzweck im
  Indikativ und referenziert auflösbare Anker (`LH-FA-SST-002`,
  `LH-FA-CFG-003`/`004`) sowie die Kopplungs-Begründung zur
  Testreihenfolge (Unabhängigkeit von `TestMVPDisableRetainedState`) —
  das ist die zulässige Kopplungs-Klasse, keine Chronik einer verworfenen
  Alternative oder eines abwesenden Textzustands. Kein Kommentar bricht
  mitten im Satz ab. Die einzige Änderung in
  `tools/harness/run-integration-tests.sh` ist eine reine
  `-run`-Muster-Zeile ohne begleitenden neuen Kommentar; der bestehende
  Kommentar oberhalb der Zeile („Eine künftig ergänzte Testfunktion
  gehört in dieses `-run`-Muster …") bleibt unverändert und trifft
  weiterhin zu.
- geprüft, ohne Befund: **Traceability.** Beide Commit-Betreffs tragen
  `LH-FA-CFG-003`/`004`, kein `SPEC-*`/`ARC-*` im Betreff; `make
  commit-traceability` lief in dieser Sitzung über `HEAD~5..HEAD` grün (0
  Befunde).
- geprüft, ohne Befund: **Hard Rule 3.1** (Docker-only) — beide neuen
  Zeilen laufen ausschließlich über den bestehenden
  `docker run … go test`-Aufruf in `run-integration-tests.sh`, kein
  lokales Toolchain-Install. **3.2** (Suppression-Verbot) — kein
  `#noqa`/`//nolint`/`[SuppressMessage]` im Diff. **3.3**
  (git-mv/Inhalt-Trennung) — nicht einschlägig, keine Umbenennung in
  diesem Diff. **3.5**/**3.6** — kein ADR im Diff berührt, keine
  Gate-Schwelle gelockert.
- geprüft, ohne Befund: **§8 Sub-Area-Sichtung konsistent mit dem
  Register-Bestand.** Beide im Plan §8 zitierten Beobachtungsverzeichnisse
  (`BEO-PGC/test-runner-stiller-ausschluss`,
  `BEO-PGC/verwaltung-keine-sql-administration`) existieren im
  Beobachtungs-Register. `BEO-PGC/test-isolation-geteilter-zustand`
  (Race-Risiko bei paralleler Paketausführung gegen `CDC_STORE_TEST_DSN`
  in `postgresstorage`-Paket-Tests) ist ein anderer Gegenstand — dieser
  Slice betrifft einen sequenziellen Compose-Testlauf innerhalb *eines*
  Testbinaries, nicht parallele Paketausführung gegen eine geteilte
  Nicht-Compose-Instanz; die Nicht-Zitierung ist damit keine Lücke.
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf) —
  `baseline-verify` (54 Dateien), `docs-check` (287 Dateien, 0 Befunde),
  `commit-traceability` (5 Commits, `HEAD~5..HEAD`, OK), `a-check` (0
  Befunde) — alle grün.
- geprüft, ohne Befund: **`make test-integration`, real dreimal in Folge
  ausgeführt** (dieser Review-Sitzung, nicht aus dem Implementer-Bericht
  übernommen) — `TestMVPActiveTablesViewMatchesActivationState` PASS in
  allen drei Läufen (0,02s/0,03s/0,04s), alle sieben Testfälle des
  Compose-Laufs PASS in jedem der drei Durchläufe, kein Arbeitsbaum-Rest
  außer dem erwarteten `tools/schema/plan.yaml`-Rollout-Artefakt (per
  `git checkout --` zurückgesetzt); `git status --porcelain` am Ende
  dieser Sitzung leer.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine.

## Verdikt

**Merge-blockierend:** keins. Der Diff ist klein, einzeln lieferbar und
hält die im Slice-Plan §1 versprochene Black-Box-Grenze zur SQL-Sicht
sauber ein; die Boundary-Wahl (`feed_mvp_idle`) und die
Isolation gegen die übrigen `feed_mvp_*`-Aktivierungen im selben
Compose-Lauf sind eigenständig gegen `compose.yaml`,
`run-integration-tests.sh` und den vollständigen Testdatei-Bestand
nachvollzogen, nicht aus dem Implementer-Bericht übernommen.

**Ausdrücklich unabhängig geprüft, nicht vom Implementer-Bericht
übernommen:** die Black-Box-Reinheit der beiden neuen SQL-Lesungen
(Code-Lektüre), die physische Existenz und Nicht-Aktivierung von
`feed_mvp_idle` (grep gegen `compose.yaml`/`run-integration-tests.sh`),
die Testisolation gegen `TestMVPDisableRetainedState` (grep über die
gesamte Testdatei plus Quelltext-Reihenfolge-Argument), die Schärfe der
Boundary-/Happy-Path-Assertions, die Nicht-Auslösung des
Image-Rebuild-Triggers (Dockerfile-Lektüre gegen die geänderten Pfade),
und `make gates`/`make test-integration` (real dreimal) als eigenständige
Läufe in dieser Sitzung.

**Übergabe:** Keine Findings, kein Rollen-Konflikt, keine
Architect-Sequenz nötig (Modul 8). Dieser Report ist ein Lauf-Beleg und
wird über Läufe hinweg nicht wieder gelesen. Er ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
