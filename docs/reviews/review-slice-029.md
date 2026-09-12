# Review-Report: slice-029 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-029`, §1/§2/§3) und
`AGENTS.md` §3 Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commits `7b7ca3c` (Testfall), `7f6456e` (DoD-Nachzug,
Plan-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-029-black-box-lesepfad-cdc-changes.md` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug, §8 Register-Sichtung)
- `docs/plan/planning/observations/BEO-PGC/lese-doppelquelle/` (observation.md, state.md, evidence/slice-010.md, evidence/slice-011.md)
- `docs/plan/planning/observations/BEO-PGC/test-isolation-geteilter-zustand/` (state.md, evidence/slice-022.md)
- `spec/lastenheft.md` (`LH-FA-REA-002`…`006`)
- `tools/schema/schema.yaml` (View `changes`, Spalten-Signatur)
- `internal/adapters/driven/postgresstorage/queries/queries.go` (`SelectChanges`), `internal/adapters/driven/postgresstorage/store.go` (`ReadChanges`)
- `test/integration/integration_test.go` (Gesamtdatei, für Isolations-Prüfung gegen andere Testfälle)
- `tools/harness/run-integration-tests.sh` (bash-Abschnitte, Zeilen-ID-Kollisionsprüfung)
- `AGENTS.md` §3 Hard Rules

---

## Findings

Keine HIGH- oder MEDIUM-Findings. Ein LOW-Finding zur Vollständigkeit
eines Kopplungs-Kommentars.

### F-1 — Isolations-Kommentar nennt nicht alle koexistierenden Zeilen-IDs auf derselben Tabelle

- `kategorie`: LOW
- `quelle`: Maintainability — Hard Rule 3.7 (Kopplungs-Klasse eines
  Kommentars soll die Kopplung vollständig benennen, nicht nur
  teilweise)
- `pfad`: `test/integration/integration_test.go:378-383` (Doc-Kommentar
  von `TestMVPChangesViewMatchesReadChanges`, Commit `7b7ca3c`)
- `befund`: Der Kommentar begründet die Wahl von `id=60` mit „isoliert
  von den übrigen Testfällen dieser Datei auf derselben Tabelle (id=1 in
  `TestMVPUpdateOldImageWithFullReplicaIdentity`, id=90/91 im
  nachgelagerten Lasttest-Beleg von `run-integration-tests.sh`)" und
  nennt damit zwei der drei tatsächlich koexistierenden ID-Gruppen auf
  `feed_mvp_full`. Der CLI-E2E-Abschnitt desselben Skripts
  (`tools/harness/run-integration-tests.sh:354` `CLI_TABLE=feed_mvp_full`,
  ids 95/96) fehlt in der Aufzählung. Eine reale Kollision entsteht dadurch
  nicht — verifiziert: `id=60` kollidiert mit keiner der drei Gruppen
  (1; 90/91; 95/96) — aber die Kopplungs-Beschreibung ist unvollständig
  und könnte eine künftige ID-Wahl (z. B. `id=61` bei einem weiteren
  Testfall) auf Basis einer unvollständigen Liste treffen lassen.
- `verifizierbar`: ja — `grep -n "CLI_TABLE=feed_mvp_full\|VALUES (9[56]"
  tools/harness/run-integration-tests.sh` zeigt die fehlende dritte Gruppe.
- `klasse`: „Kopplungs-Kommentar unvollständig, aber nicht falsch"

## Negativbefunde

- geprüft, ohne Befund: **Black-Box-Eigenschaft der SQL-View-Seite** —
  `queryChangesView` liest ausschließlich über `env.pool.Query` (rohes
  pgx-SQL gegen `cdc.changes`); kein Import und kein Aufruf von
  `postgresstorage.NewChangeStoreAdapter`/`ReadChanges` auf dieser Seite
  des Vergleichs. Die Spaltenliste der Abfrage (`commit_position,
  sequence, operation, old_data, new_data, schema_version`, gefiltert auf
  `source_id`/`source_table_id`) deckt sich exakt mit der
  View-Definition in `tools/schema/schema.yaml` (Zeilen 232-261).
- geprüft, ohne Befund: **Plan-Konformität der Implementierungs-Wahl** —
  Ziel in §1 erlaubt ausdrücklich „rohes `docker exec … psql`/SQL"; der
  Plan-Nachzug §3 begründet die gewählte Variante (`pool.Query` statt
  Bash/`psql`) mit derselben Außensicht auf die SQL-View, ohne eine
  zweite Sprache im Vergleich zu tragen. Keine Abweichung vom Ziel, die
  nicht im Plan-Nachzug benannt wäre.
- geprüft, ohne Befund: **Realer Vergleich, keine Zeitfenster-Race** —
  `awaitChangesViewRows` wartet, bis die View 3 Zeilen für `id=60`
  zeigt, *bevor* `env.store.ReadChanges` gelesen wird. Beide Lesewege
  fragen dieselben Basistabellen ab (`SelectChanges` in
  `queries.go:37-57`: `cdc.change JOIN cdc.transaction`; die View
  `changes` in `schema.yaml:227-261` auf denselben Tabellen) — sobald die
  View die 3 Zeilen zeigt, sind sie in Postgres committet und für jede
  Folgeabfrage auf denselben Tabellen sofort sichtbar. Kein
  Wartemechanismus fehlt auf der `ReadChanges`-Seite, der zu einem
  falschen Grün führen könnte.
- geprüft, ohne Befund: **Testisolation `id=60`** — kein anderer Testfall
  in `test/integration/integration_test.go` oder in den bash-Abschnitten
  von `tools/harness/run-integration-tests.sh` verwendet `id=60` auf
  `feed_mvp_full` oder einer anderen Tabelle im selben Compose-Lauf; die
  Tabelle `feed_mvp_full` trägt sonst nur `id=1`
  (`TestMVPUpdateOldImageWithFullReplicaIdentity`, andere Tabelle
  `feed_mvp_flow` nutzt ebenfalls `id=1`, aber physisch getrennt) sowie
  `id=90/91` (Lasttest-Beleg) und `id=95/96` (CLI-E2E), beide im
  bash-Teil nach dem Go-Testbinary. Kein `t.Parallel()` im Paket — Läufe
  sind sequenziell. Siehe F-1 zur Unvollständigkeit der Kommentar-Liste
  (kein Isolationsdefekt).
- geprüft, ohne Befund: **DoD-Aktualisierung gegen tatsächlichen Diff**
  (die `slice-028`-Lehre F-1: Closure-Beleg breiter formuliert als
  Testumfang) — die drei gesetzten Häkchen in `7f6456e` sind durch den
  Testfall in `7b7ca3c` gedeckt: Häkchen 1 (Testfall existiert und liest
  real über `cdc.changes`), Häkchen 2 (Vergleich auf Reihenfolge und
  Feldinhalt, alle sechs Assertion-Paare im Diff vorhanden), Häkchen „Doku-
  Update" (zutreffend „keiner berührt" — reine Testabdeckung, keine
  Guide-/Sensor-/Vertragsänderung in den beiden Commits). Keine
  Überbehauptung gefunden.
- geprüft, ohne Befund: **§8-Register-Sichtung** — die im Slice-Kopf §8
  behaupteten Zähler-Stände sind gegen das Register nachvollzogen:
  `BEO-PGC/lese-doppelquelle` trägt genau zwei Belege
  (`evidence/slice-010.md`, `evidence/slice-011.md`, Zustand „weiter
  offen"), `BEO-PGC/test-isolation-geteilter-zustand` genau einen
  (`evidence/slice-022.md`). Beide Angaben stimmen.
- geprüft, ohne Befund: **Kompilierbarkeit** — `go vet
  ./test/integration/...` (Toolchain-Image, gecachtes Modul-Volume,
  netzlos) lief ohne Ausgabe/Fehler; `gofmt -l
  test/integration/integration_test.go` meldete keine Datei (bereits
  formatiert). Beide Läufe eigenständig in dieser Review-Sitzung
  ausgeführt, nicht vom Implementer-Bericht übernommen.
- geprüft, ohne Befund: **Traceability** — beide Commit-Betreffs tragen
  keine `SPEC-*`/`ARC-*`-Kennung; `LH-FA-REA-002` steht im Bodytext
  beider Commits. `make commit-traceability` lief in dieser Sitzung über
  `HEAD~5..HEAD` grün (0 Befunde).
- geprüft, ohne Befund: **Hard Rule 3.3** (git-mv/Inhalt-Trennung) — nicht
  einschlägig, beide Commits sind reine Content-Commits ohne Umbenennung.
- geprüft, ohne Befund: **Hard Rule 3.7** (Kommentar-Klassen) — die neuen
  Kommentare (`changesViewRow`, `queryChangesView`,
  `awaitChangesViewRows`, `changeRowID`, Testfall-Doc-Kommentar)
  beschreiben durchweg den geltenden Zustand indikativ (Zusage-/
  Kopplungs-Klasse), keine verworfene Alternative, kein abgebrochener
  Satz. Einschränkung zur Vollständigkeit einer Kopplungsangabe siehe F-1
  (Vollständigkeit, nicht Klassen-Verstoß).
- geprüft, ohne Befund: **Hard Rule 3.1** (Docker-only), **3.2**
  (Suppression-Verbot), **3.5** (ADR-Immutabilität), **3.6**
  (Gate-Lockerung) — keine Berührung; kein `#noqa`/`//nolint`, kein
  ADR editiert, keine Gate-Schwelle geändert.
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf) —
  `baseline-verify`, `docs-check`, `commit-traceability`, `a-check` alle
  grün, 0 Befunde.
- **Prüfgrenze — nicht reproduzierbar in dieser Review-Sitzung:** Der
  Implementer-Bericht behauptet einen realen Mutationstest (Vertauschung
  von `old_data`/`new_data` in der `changes`-View-Spaltenprojektion,
  Testfall lief dabei isoliert rot, danach zurückgesetzt) sowie drei
  weitere grüne `make test-integration`-Läufe. `git log -p`,
  `git reflog` und `git stash list` zeigen dazu erwartungsgemäß keine
  Spur (die Mutation wurde lokal vor dem Commit zurückgesetzt,
  `tools/schema/schema.yaml` ist im Diff beider Commits unverändert und
  `git diff` gegen den aktuellen Stand zeigt keine Abweichung). Das ist
  kein Fehler, aber eine echte Prüfgrenze dieses Reviews: Die
  Mutationsprobe und die drei wiederholten `make test-integration`-Läufe
  sind aus dem Repo-Zustand allein nicht nachvollziehbar und bleiben
  unabhängig vom Implementer-Bericht unbestätigt. Eine eigenständige
  Reproduktion von `make test-integration` (braucht eine laufende
  Compose-Umgebung) ist zudem Verifier-Aufgabe (Modul 8/11), nicht Teil
  dieses Maintainability-Reviews — konsistent mit `review-slice-028.md`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Kopplungs-Kommentar unvollständig,
aber nicht falsch

## Verdikt

**Merge-blockierend:** nein. Der neue Testfall ist real black-box
gegenüber dem Go-Adapter, deckt denselben Rundlauf auf beiden
Lesewegen ab, ist kollisionsfrei isoliert und die DoD-Häkchen sind durch
den Diff gedeckt (keine Wiederholung der `slice-028`-Überbehauptungslehre).
F-1 ist eine kleine Präzisierung ohne Merge-Blockade.

**Ausdrücklich unabhängig geprüft, nicht vom Implementer-Bericht
übernommen:** Black-Box-Eigenschaft der SQL-Seite, Abwesenheit einer
Zeitfenster-Race (durch Vergleich der zugrundeliegenden Tabellen beider
Lesewege), Testisolation gegen alle drei koexistierenden ID-Gruppen auf
`feed_mvp_full`, Kompilierbarkeit (`go vet`, `gofmt`), Register-Zähler-
Stände aus §8, Traceability. Die Mutationstest-Behauptung bleibt als
Prüfgrenze explizit unbestätigt (siehe Negativbefunde) — sie wird hier
weder bestätigt noch verworfen, sondern als außerhalb der Nachprüfbarkeit
dieses Review-Laufs benannt.

**Übergabe:** F-1 geht an den Implementer/Planner zur freien Entscheidung
(LOW, keine Sequenz nach Modul 8 nötig). Die **Finding-Klasse** geht
zusätzlich in die Slice-Closure §7 und von dort in den Zähler des
Beobachtungs-Registers. Dieser Report ist ein Lauf-Beleg und wird über
Läufe hinweg nicht wieder gelesen. Er ersetzt keine Verifikation —
DoD-/Spec-Konformität (inkl. Reproduktion von `make test-integration`
und Bewertung der Mutationstest-Behauptung) prüft der Verifier separat.
