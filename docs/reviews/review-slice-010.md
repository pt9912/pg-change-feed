# Review-Report: slice-010 — 2026-09-10

**Review-Art:** Code — geprüft gegen Slice-Plan + ADRs (Maintainability).

**Gegenstand:** Implementer-Commits von slice-010, Range `79c0c12..HEAD`
minus dem Planner-Zug `c57215a` (Register-Nachtrag) — drei Implementer-Commits:
`9a7b3c1` (Plan-Nachzug), `fe1ff63` (SQL-Views), `c534fec` (SQL-View-Tests).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-010 (Lesen-Vollabdeckung und SQL-Schnittstelle,
  `LH-FA-REA-002`…006, `LH-FA-SST-002`, Welle 3)
- `spec/lastenheft.md` §`LH-FA-REA-001`…006, §`LH-FA-SST-002` ·
  `spec/architecture.md` §1/§2 (`ARC-005` Driving, `ARC-006` Driven,
  Import-Constraints)
- [`ADR-0009`](../plan/adr/0009-change-store-outbound-port.md) (Change Store
  Outbound Port), [`ADR-0018`](../plan/adr/0018-sql-driving-adapter.md) (SQL
  als Driving Adapter), [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
  (d-migrate)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.2 Suppression-Verbot, §3.7
  Kommentar-Klassen) · `harness/conventions.md` (MR-000/MR-001) ·
  `.a-check.yml` (Layer-Edges)
- `.claude/commands/implement-slice.md` §„Plan-Nachzug im selben Lauf ·
  seit slice-009"
- Register-Sichtung: `BEO-PGC/d-migrate-nacharbeit` (jetzt 2×,
  `evidence/slice-006.md` + `evidence/slice-010.md`), `BEO-PGC/plan-nachzug`
  (verkörpert, 2×), `BEO-PGC/walsender-wirksamkeit` (1×, Replication-Stream —
  von diesem Slice nicht berührt), `BEO-PGC/adapter-fehler-ausgang` (1×,
  Bootstrap-Verdrahtung — von diesem Slice nicht berührt),
  `BEO-PGC/a-check-null-abdeckung` (verkörpert)
- vorherige Reports: review-slice-006 (F-3, d-migrate-Ausweichform),
  review-slice-007 (Closure-Notiz: „CLI-Adapter und SQL-Funktionen —
  `ADR-0018`/`ADR-0019`-Rest, folgen nach dem MVP"), review-slice-009 (F-1
  Plan-Nachzug 9. Auftreten — Architekt-Sequenz angekündigt; F-5 §8
  unausgefüllt, 1. Auftreten)

---

## Findings

### F-1 — SQL-Views lesen die Driven-Adapter-Tabellen direkt statt Inbound Ports aufzurufen; ADR-0018 ist in Plan und Commits nicht referenziert

- `kategorie`: HIGH
- `quelle`: [`ADR-0018`](../plan/adr/0018-sql-driving-adapter.md) (Entscheidung:
  „SQL-Funktionen/Views sind Driving Adapter und rufen Inbound Ports auf.
  Businesslogik wird nicht in SQL dupliziert.") · `spec/architecture.md` §2
  (Zeile `ARC-005` Driving: „Darf importieren: Inbound Ports (`ARC-003`)" /
  „Darf NICHT importieren: … Driven Adapters")
- `pfad`: `tools/schema/nacharbeit-views.sql:9-56` (`active_tables`,
  `consumer_status`, `changes`) · `tools/schema/schema.yaml:12-27`
  (Kommentar-Block) · Commit `fe1ff63`
- `befund`: Die drei neuen SQL-Views selektieren direkt aus
  `cdc.source_table`, `cdc.consumer`, `cdc.consumer_position`,
  `cdc.transaction`, `cdc.change` — denselben Tabellen, die der
  `PostgresChangeStoreAdapter` (`ARC-006`, Driven) verwaltet — ohne über
  einen Inbound Port zu laufen. Das widerspricht dem Wortlaut der
  Entscheidung C aus `ADR-0018` (gewählt gegen Option A „Logik in
  PL/pgSQL-Funktionen"), die für „SQL-Funktionen/Views" explizit „rufen
  Inbound Ports auf" verlangt. Weder der Slice-Kopf (`Bezug: … ADR-0009`,
  nicht `ADR-0018`) noch die Commit-Messages `9a7b3c1`/`fe1ff63`/`c534fec`
  noch der Kommentar-Block in `schema.yaml` erwähnen `ADR-0018` — obwohl
  slice-007s eigene Closure-Notiz genau diese Lieferung dorthin verwiesen
  hatte („CLI-Adapter und SQL-Funktionen — `ADR-0018`/`ADR-0019`-Rest, folgen
  nach dem MVP"). `ADR-0018` erklärt die Regel selbst als „Review-Prüfpflicht"
  (kein Gate deckt Go-Import-Verstöße in reinem SQL ab; `.a-check.yml`
  kennt nur Go-Globs) — der Konflikt ist also genau der Fall, den kein
  Sensor fängt.
- `verifizierbar`: nein — kein Gate deckt SQL-Layer-Edges; Lese der Views
  gegen `ADR-0018` und `spec/architecture.md` §2
- `klasse`: ADR-Verstoß — SQL-Driving-Adapter ohne Inbound-Port-Aufruf
  (1. Auftreten)

**Konflikt-Pfad (Modul 8):** HIGH mit möglichem Rollen-Widerspruch (der
Implementer könnte einwenden, reine Lesezugriffe ohne Zustandsänderung seien
von „Businesslogik wird nicht dupliziert" gedeckt und die Ports-Klausel ziele
auf Administrations-Funktionen). Dieser Report benennt den Konflikt, löst ihn
aber nicht — das ist die Architect-Sequenz (`ADR-0018` bestätigen / per
Folge-ADR schärfen / Lockerung als Folge-ADR nachziehen), kein Implementer-
Ermessen.

### F-2 — §8 vorgelagerte Prüfungen bleiben Vorlagen-Platzhalter (2. Auftreten, außerhalb dieses Diffs)

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Zwei Schritte
  vor der Modus-Begründung · review-slice-009 F-5 (1. Auftreten)
- `pfad`: `docs/plan/planning/in-progress/slice-010-lese-vollabdeckung-sql.md`
  §8 („**Vorgelagert — Sub-Area-Wahl prüfen:** <je berührter Sub-Area …>" ·
  „**Vorgelagert — offene Beobachtungen sichten:** <Register durchgegangen
  …>")
- `befund`: Beide unbedingten §8-Prüfungen stehen unverändert als
  Vorlagen-Platzhalter — schon in der Fassung, die `9a7b3c1` vorfand (Datei
  entstand bei der Welle-3-Eröffnung `7fc9d23`, vor review-slice-009 F-5).
  Der Plan-Nachzug-Commit ändert nur §3, wie die Regel in
  `.claude/commands/implement-slice.md` es verlangt; §8 ist nicht Teil des
  Implementer-Plan-Nachzugs. **Nicht Teil dieses Diffs** — außerhalb des
  Review-Gegenstands (Implementer-Commits) — aber zweites Auftreten derselben
  Klasse und damit ein Hinweis für den **Planner**, der die Sichtung ohnehin
  faktisch durchgeführt hat (§4 zitiert `BEO-PGC/d-migrate-nacharbeit`), sie
  aber nirgends im Plan trägt.
- `verifizierbar`: ja — Lese der Plan-Datei gegen die Vorlage
- `klasse`: Vorgelagerte §8-Prüfungen unausgefüllt (2. Auftreten)

---

## Negativbefunde

- geprüft, ohne Befund: **Plan-Nachzug im selben Lauf, erster Praxistest der
  Regel seit slice-009** — `9a7b3c1` (11:47:08) liegt vor `fe1ff63`
  (11:47:19) und `c534fec` (11:47:31): der Plan-Commit ist tatsächlich vor
  den Code-Commits im selben Lauf entstanden, nicht als Nachtrag danach. §3
  trägt nach dem Nachzug alle vier abweichenden Punkte (Streichung der
  `queries/*.go`-Zeile, `nacharbeit-views.sql` neu, `Makefile`-Update,
  `sqlviews_test.go` neu) — der gelieferte Datei-Umfang von `fe1ff63`/`c534fec`
  ist damit vollständig durch §3 gedeckt. `BEO-PGC/plan-nachzug` bleibt bei
  2× stehen, erreicht mit diesem Vorgang **nicht** die 3×-Schwelle
- geprüft, plausibel: **Streichung der geplanten `queries/*.go`-Zeile** — die
  fünf zitierten Tests (`TestReadCarriesPositionsAndRanges`,
  `TestReadCarriesLimit`, `TestReadCarriesTableFilter`,
  `TestReadIsDeterministicallySorted`, `TestReadLeavesPersistedStateUnchanged`
  in `internal/adapters/driven/postgresstorage/store_test.go`, seit
  slice-004) decken Bereich (`LH-FA-REA-001`), Startposition/Wiederholbarkeit
  (`LH-FA-REA-002`/005), Limit (`LH-FA-REA-003`), Ordnung inkl. Tiebreak
  (`LH-FA-REA-004`) und Tabellenfilter (`LH-FA-REA-006`) real gegen
  PostgreSQL, inklusive Grenzfällen (leerer Bereich, invertierter Bereich,
  Limit 0, Filter ohne Treffer) — der Befund trägt
- geprüft, konsistent: **d-migrate-Ausweichform, 2. Beleg** — dieselbe Form
  wie bei `chk_change_operation` (slice-006): benannte Grenze im
  `schema.yaml`-Kommentar, Nacharbeit-Datei mit `CREATE OR REPLACE VIEW`
  (wiederholbar), zweiter psql-Schritt im `schema-rollout`-Target,
  Register-Beleg vorbereitet; `ADR-0043`s Re-Evaluierungs-Trigger feuert
  nicht, weil eine Ausweichform („berichtete manuelle Nacharbeit") existiert
  — Text des Triggers nennt Exit 8, real beobachtet sind Exit 5 („Post-execute
  compare detected drift") bzw. Exit 7 (slice-006); diese Diskrepanz ist
  vorbestehend (schon bei slice-006 nicht Exit 8) und wird hier nicht neu
  eingeführt
- geprüft, ohne Befund: **`active_tables`-Konsistenz-Grenze** — der
  Kommentar in `tools/schema/schema.yaml:16-22` benennt explizit und
  indikativ, dass die View die Publication-Mitgliedschaft (`LH-FA-CFG-003`
  „aktiviert" vs. „retained") nicht nachbildet und diese Prüfung beim
  `TableActivationPort` bleibt — trägt die Grenze-Kommentarklasse
  (`AGENTS.md` §3.7)
- geprüft, ohne Befund: **Kommentar-Klassen** (`AGENTS.md` §3.7) in
  `tools/schema/nacharbeit-views.sql`, `tools/schema/schema.yaml`,
  `Makefile`, `sqlviews_test.go` — indikativ über den geltenden Zustand,
  keine Konjunktiv-Klausel über verworfene Alternativen, keine abgebrochenen
  Sätze
- geprüft, ohne Befund: **Docker-only** — beide `schema-rollout`-Schritte
  laufen über `docker run` (d-migrate-Image bzw. `PG_TEST_IMAGE`), kein
  lokales Toolchain-Install
- geprüft, ohne Befund: **Suppression-Verbot** — keine `nolint`/`noqa`/
  `SuppressMessage`-Marker in den drei Commits
- geprüft, ohne Befund: **Traceability der drei Implementer-Commits** — je
  Betreff mindestens eine `LH-*`-/`ADR-*`-Kennung, keine Struktur-ID
  (`SPEC-*`/`ARC-*`) im Betreff
- geprüft, plausibel, nicht im Review-Lauf reproduziert: **Mutationsproben**
  — die drei berichteten Mutationen (`source_table` per `WHERE 1=0`
  ausgeschlossen; `consumer_position`-Join per `ON false` gekappt;
  `source_table`-Join in `changes` per `ON false` gekappt) entsprechen exakt
  den realen Join-/Subquery-Konstrukten der drei Views (Scalar-Subquery in
  `active_tables`, `LEFT JOIN cdc.consumer_position` in `consumer_status`,
  `JOIN cdc.source_table` in `changes`) — strukturell plausibel, aber nicht
  in diesem Review-Lauf erneut ausgeführt; Bestätigung ist Verifier-Sache
  (`make test-store`)
- geprüft, ohne Befund: **a-check-Edges** — `sqlviews_test.go` liegt im
  `adapters`-Layer-Glob und referenziert nur `pgxpool`/Standardbibliothek,
  keine Rückwärts-Kante; die SQL-Views selbst liegen außerhalb der
  Go-Layer-Globs von `.a-check.yml` und sind für a-check unsichtbar (siehe
  F-1)
- geprüft, ohne Befund: **WIP-Limit und Lifecycle** — genau ein Slice in
  `in-progress/`; die drei Lifecycle-Commits vor dem Implementer-Range sind
  reine Feld-/Move-Commits mit Kennung im Betreff
- geprüft, ohne Befund: **Register-Sichtung außerhalb dieses Slices** —
  `BEO-PGC/walsender-wirksamkeit` (Replication-Stream) und
  `BEO-PGC/adapter-fehler-ausgang` (Bootstrap-Verdrahtung) betreffen
  Sub-Areas, die slice-010 nicht berührt; `BEO-PGC/a-check-null-abdeckung`
  ist bereits verkörpert

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** ADR-Verstoß — SQL-Driving-Adapter ohne
Inbound-Port-Aufruf (1. Auftreten) · Vorgelagerte §8-Prüfungen unausgefüllt
(2. Auftreten, außerhalb dieses Diffs)

**Sequenz-Beobachtung (Steering-Loop):** F-1 ist das erste Auftreten seiner
Klasse — kein automatischer Architekt-Zug nach der Drei-Zähler-Regel, aber
als **HIGH mit möglichem Rollen-Widerspruch** läuft der Konflikt-Pfad ab
sofort als Sequenz mit Übergabe-Artefakten (Modul 8), nicht als Herabstufung
im Implementer-Dialog. F-2 ist das zweite Auftreten der §8-Klasse aus
review-slice-009 F-5 — beide Vorkommen liegen an Planner-erzeugten
Plan-Dateien, nicht an Implementer-Commits; die Klasse gehört in die
Planner-seitige Schärfung (Vorlage bei Slice-Anlage konsequent mit den zwei
Prüfungen befüllen), nicht in diesen Implementer-Bericht.

**Plan-Nachzug-Regel (erster Praxistest, seit slice-009):** getragen — der
Plan-Commit `9a7b3c1` liegt zeitlich vor den Code-Commits `fe1ff63`/`c534fec`
im selben Lauf, und §3 deckt nach dem Nachzug den vollständigen
Lieferumfang. Kein drittes Auftreten der Klasse „Plan-Erweiterung ohne
Plan-Nachzug".

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH). Die Architektur-Frage, ob reine
Lesezugriffe über SQL-Views `ADR-0018`s „rufen Inbound Ports auf" erfüllen
müssen oder ob die Klausel auf Administrations-Funktionen zielt, ist nicht
implementierbar-lokal zu entscheiden; sie braucht die Architect-Sequenz
(Modul 8: `ADR-0018` bestätigen mit Plan-Korrektur / Folge-ADR mit
`supersedes` / Lockerung per Folge-ADR nachgezogen).

**Übergabe:** F-1 geht an Implementer **und** Architect (Konflikt-Pfad,
Modul 8) — kein stilles Weiterimplementieren auf den drei Views, bis eines
der drei Verdikte vorliegt. F-2 geht als Hinweis an den Planner (außerhalb
des Implementer-Diffs, keine Closure-Blockade durch diesen Report).

---

**Gate-Beleg:** `make gates` nach diesem Report-Commit (Lauf 2026-09-10,
Range-Head); Ergebnis im Commit-Text vermerkt.
