# Review-Report: slice-048 — Fixrunde — 2026-09-13

**Review-Art:** Code — Bestätigungslauf zu einer Fixrunde nach eigenem
Vorbefund. Geprüft gegen den eigenen vorherigen Report
(`review-slice-048.md`, Finding F-1, MEDIUM), den Fix-Commit `668f2cd`,
`AGENTS.md` §3.7 (Kommentar-Disziplin) und die Commit-Traceability-Regeln
(`harness/README.md` §Traceability rules) — Rollentrennung Modul 8: diese
Prüfung läuft eigenständig gegen Skript, Slice-Plan-Diff und einen real
ausgeführten `make test-integration`-Lauf, nicht als Übernahme der
Implementer-Zusammenfassung aus der Commit-Message.

**Gegenstand:** Commit `668f2cd` (`fix(test-integration):
Retention-Lebenszyklus-Rundlauf — Löschbeleg vor DELETE-Nachbereitung
(LH-FA-RET-004)`) — geändert:
`docs/plan/planning/in-progress/slice-048-retention-lebenszyklus-kombinierter-e2e-rundlauf.md`,
`tools/harness/run-integration-tests.sh`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/reviews/review-slice-048.md` (vollständig — eigener Vorbefund F-1)
- Vollständiger `git show 668f2cd` (beide geänderten Dateien, mit Kontext)
- `internal/domain/model/retention.go` (`RetentionPolicy.AllowsDeletion`,
  bereits im Vorbefund gelesen, zur Bestätigung der Kausal-Unabhängigkeit
  erneut herangezogen)
- `AGENTS.md` §3.7 (Kommentar-Disziplin), §3.3 (`git mv` +
  Inhaltsänderung), `harness/README.md` §Traceability rules
- reale, eigenständig ausgeführte `make test-integration`- und
  `make gates`-Läufe gegen `HEAD` (`668f2cd`), nicht aus dem
  Implementer-Bericht übernommen

---

## Findings

### F-1 (Vorbefund) — Testbeleg-Reihenfolge verdeckte behauptete Kausal-Unabhängigkeit — bestätigt behoben

- `kategorie`: MEDIUM (Vorbefund, jetzt geschlossen)
- `quelle`: `LH-FA-RET-004`
- `pfad`: `tools/harness/run-integration-tests.sh:473-515`
- `befund`: Die Reihenfolge ist jetzt: zweite Bestätigung
  (`acknowledge-consumer … "$lifecycle_position"`, Zeile 473) →
  Polling-Schleife auf die reale Löschung von `id=210` (Zeile 483-496,
  `lifecycle_deleted`) → Log-Zeile, die explizit festhält, dass die
  bestätigte Position zu diesem Zeitpunkt noch real in
  `cdc.consumer_position` vorhanden ist (Zeile 498) → **erst danach**
  `DELETE FROM cdc.consumer_position …` (Zeile 505-506) →
  `blocker_after`-Check (Zeile 508-513). Das ist exakt die im Vorbefund
  verlangte Reihenfolge: Der Testlauf beobachtet jetzt real den Zustand
  „bestätigte, weiterhin vorhandene Position gibt die Löschung frei"
  statt nur „keine Position mehr vorhanden". Die im Vorbefund per
  Code-Lektüre (`RetentionPolicy.AllowsDeletion`,
  `internal/domain/model/retention.go:35-44`) bereits bestätigte
  Kausal-Unabhängigkeit ist damit zusätzlich empirisch durch den Testlauf
  selbst demonstriert, nicht nur durch Domain-Code-Analyse.
- `verifizierbar`: ja — durch eigenen `make test-integration`-Lauf
  bestätigt (siehe unten).
- `klasse`: „Testbeleg-Reihenfolge verdeckt behauptete
  Kausal-Unabhängigkeit" (Vorbefund, jetzt geschlossen)

## Eigener `make test-integration`-Lauf gegen `668f2cd`

Eigenständig (nicht aus der Commit-Message übernommen) ein realer Lauf
gegen den Fix-Commit ausgeführt. Ergebnis: alle Go-Integrationstests und
alle Black-Box-CLI-Abschnitte grün, kein Fehler, Feed-Container lief
durchgehend weiter. Relevanter Ausschnitt, mit Zeilennummern im
Lauf-Log:

```
102: run-integration-tests: Retention-Lebenszyklus-Rundlauf —
     'RetentionLifecycle' (id=210) real entfernt, während die bestätigte
     Position von cli-e2e-lifecycle-consumer noch real in
     cdc.consumer_position vorhanden ist (LH-FA-RET-003/004)
103: run-integration-tests: Retention-Lebenszyklus-Rundlauf —
     cdc.retention_blockers zeigt real keinen Blocker mehr für src-mvp,
     auch nach der bereits erfolgten realen Löschung (LH-FA-RET-005)
```

Zeile 102 (Löschbeleg, Position noch vorhanden) erscheint **vor** Zeile 103
(Sichtbarkeits-Check nach `DELETE`) — nicht nur behauptet, sondern im
tatsächlichen `stdout`-Strom des Laufs so beobachtet. Damit ist die im
Vorbefund verlangte zeitliche Reihenfolge real bestätigt, nicht nur aus dem
Skript-Text abgeleitet. Der Implementer-Bericht behauptet drei
aufeinanderfolgende grüne Läufe während der Fixrunde; dieser Report
bestätigt einen weiteren, eigenständig ausgeführten Lauf danach — die
geforderte Mindestanzahl (≥ 1) ist damit erfüllt, unabhängig von den drei
bereits vom Implementer dokumentierten.

## Prüfung: Interferenz mit dem nachfolgenden „Zustand 1 — kein Blocker"-Abschnitt

Frage: Verschiebt die neu eingefügte Polling-Schleife (Zeile 483-496) das
`DELETE` (Zeile 505-506) so weit nach hinten, dass es mit dem
`slice-047`-Abschnitt „Zustand 1" (beginnt Zeile 532) kollidiert?

Befund: nein. Das `DELETE` steht weiterhin **vor** Zeile 532 — die
Umordnung hat nur die Polling-Schleife *zwischen* die zweite Bestätigung
(vorher Zeile 472) und das `DELETE` (vorher Zeile 481) eingefügt; sie hat
das `DELETE` nicht *hinter* den „Zustand 1"-Abschnitt verschoben. Der
„Zustand 1"-Abschnitt setzt weiterhin voraus, dass `cdc.consumer_position`
für `src-mvp` bereits leer ist (dessen eigener Kommentar, Zeile 532-543,
unverändert gegenüber dem Vorbefund) — diese Voraussetzung ist zum
Zeitpunkt von Zeile 532 unverändert erfüllt, weil das `DELETE` bei Zeile
505-506 bereits gelaufen ist. Der reale `make test-integration`-Lauf
bestätigt das zusätzlich empirisch: Der Diagnose-Beleg „kein Blocker" (Log:
„Retention-Sichtbarkeits-Beleg (CLI, Zustand 1) — 'kein Blocker' vor jeder
Consumer-Bestätigung …") lief im selben Lauf unverändert grün, direkt im
Anschluss an die Retention-Lebenszyklus-Kette.

## Prüfung: Kommentar-Disziplin (`AGENTS.md` §3.7)

Die neu/geänderten Kommentarzeilen in `run-integration-tests.sh` (Zeilen
385-388, 478-482, 500-504) beschreiben durchgehend den **Ist-Zustand** an
der jeweiligen Skript-Stelle im Indikativ („bleibt … bestehen", „ist an
dieser Stelle bereits real gelöscht", „dient ausschließlich …") — keine
Konjunktiv-Aussage über eine verworfene Alternative, kein Verweis auf
abwesenden Text, kein `slice-\d+`-Zitat. Eigener `grep -nE
"slice-[0-9]+|früher|vorher stand|wurde entfernt"` gegen den Diff-Ausschnitt
von `tools/harness/run-integration-tests.sh` in `668f2cd`: kein Treffer.

Die Prosa im Slice-Plan (`§3 Plan-Nachzug Punkt 2`, `§7 Closure-Notiz`)
nennt dagegen ausdrücklich „Reviewer-Finding F-1" und
`docs/reviews/review-slice-048.md` — das ist an dieser Stelle **zulässig
und vorgesehen**: §7 ist die Closure-Notiz, deren Aufgabe laut
Baseline-Regelwerk `modul-05-planning-harness.md` §Closure- und
Lerneintrag-Regeln ausdrücklich ist, eine „wiederkehrende Finding-Klasse
aus dem Review" als eine von drei Quellen aufzunehmen. Hard Rule 3.7 gilt
für Code-, Konfigurations- und Skript-Kommentare (und für
Stand-/Status-Zellen in Roadmap/Beobachtungs-Register/Meilenstein-Tabellen)
— nicht für die Fließtext-Closure-Notiz eines Slice-Plans. Kein Verstoß.

## Prüfung: Commit-Traceability

`git show 668f2cd --format="%s" -s`:
`fix(test-integration): Retention-Lebenszyklus-Rundlauf — Löschbeleg vor
DELETE-Nachbereitung (LH-FA-RET-004)` — trägt `LH-FA-RET-004` im Betreff,
keine `SPEC-*`/`ARC-*`-ID. Eigener Lauf `RANGE=11889c2..668f2cd make
commit-traceability`: `d-check --enable commits` 0 Befunde,
`commit-traceability.sh` OK (1 Commit, Betreff ohne Struktur-ID).

## Prüfung: Plan-Nachzug / §7 gegen die tatsächliche Umsetzung

`§3 Plan-Nachzug Punkt 2` beschreibt jetzt exakt die im Skript umgesetzte
Reihenfolge (Löschbeleg zuerst, `DELETE`/Sichtbarkeits-Check danach) und
benennt konkret, dass „alle drei Zwischenzustände" real geprüft sind:
`blocker_before`, der Löschbeleg selbst (`lifecycle_deleted`, real während
die Position noch besteht) und `blocker_after` — das deckt sich mit dem
tatsächlichen Skript (Zeilen 452-459, 483-498, 505-515). `§7 Closure-Notiz`
trägt einen eigenen Absatz „Fixrunde (Reviewer-Finding F-1, MEDIUM)", der
die alte und neue Reihenfolge, die Begründung und den realen Log-Beleg aus
drei Implementer-Läufen benennt — konsistent mit dem Diff. Kein
Widerspruch zwischen Plan-Text und Code gefunden.

## Negativbefunde

- geprüft, ohne Befund: **weitere Reihenfolge-Regression** — der
  bestehende, unveränderte Kommentar zum Abschnittsstart (Zeile 381-393)
  wurde konsistent mit der neuen Reihenfolge nachgezogen (nennt jetzt
  „die reale Löschung … danach die reale Abwesenheit jeder Zeile"); kein
  veralteter Text stehen geblieben.
- geprüft, ohne Befund: **Interferenz mit dem „Zustand 1"-Abschnitt**
  (`slice-047`) — siehe eigener Abschnitt oben; `DELETE` bleibt vor
  Zeile 532, reale Bestätigung im Testlauf.
- geprüft, ohne Befund: **Isolation** (`BEO-PGC/test-isolation-geteilter-zustand`)
  — keine neuen IDs/Consumer-Namen eingeführt, nur Reihenfolge innerhalb
  desselben bestehenden Abschnitts verändert.
- geprüft, ohne Befund: **`git mv` + Inhaltsänderung** (`AGENTS.md` §3.3)
  — keine Datei-Verschiebung in diesem Commit, reine Inhaltsänderung an
  zwei bestehenden Dateien; Regel nicht einschlägig.
- geprüft, ohne Befund: **DoD-Checkbox „Review durchgeführt"** — bereits
  `[x]` seit dem ursprünglichen Review-Commit (`11889c2`), unverändert
  durch `668f2cd` (Diff berührt DoD-Sektion nicht). Kein Nachzug nötig
  (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde) —
  die Checkbox war zum Zeitpunkt des ursprünglichen Verdikts bereits
  korrekt gesetzt, weil dieses Finding von Anfang an ohne
  Reviewer→Implementer-Rückgabe-Pfeil weitergereicht wurde.
- geprüft, ohne Befund: **`make gates` real ausgeführt, grün** —
  `baseline-verify` (54 Dateien), `docs-check` (377 Dateien, 0 Befunde),
  `commit-traceability` (`HEAD~5..HEAD`, 5 Commits, Betreffs ohne
  Struktur-ID), `a-check` (0 Befunde).
- geprüft, ohne Befund: **Docker-only** (`AGENTS.md` §3.1) — keine neue
  lokale Toolchain, `DELETE`/Polling laufen unverändert über
  `docker exec`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine neuen; F-1 aus dem Vorbefund ist
hier dessen Bestätigung (geschlossen), kein neues Auftreten.

## Verdikt

**F-1: bestätigt behoben.** Die Reihenfolge in
`tools/harness/run-integration-tests.sh` ist jetzt: zweite Bestätigung →
Poll auf reale Löschung (bestätigte Position bleibt dabei bestehen) →
Log-Beleg → **erst danach** `DELETE FROM cdc.consumer_position` →
`blocker_after`-Check. Ein eigenständig ausgeführter
`make test-integration`-Lauf bestätigt real und im tatsächlichen
Log-Strom (nicht nur laut Skript-Text), dass die Log-Zeile zum
Löschbeleg zeitlich vor der Log-Zeile zum Sichtbarkeits-Check nach dem
`DELETE` erscheint — der Testlauf demonstriert die im Vorbefund per
Code-Lektüre bereits bestätigte Kausal-Unabhängigkeit jetzt zusätzlich
empirisch. Keine neue Interferenz mit dem nachfolgenden
„Zustand 1"-Abschnitt (`DELETE` bleibt vor diesem Abschnitt, nur die
Polling-Schleife wurde davor eingefügt). Kommentar-Disziplin
(`AGENTS.md` §3.7) gewahrt — Ist-Zustand-Beschreibungen im Skript, kein
Chronik-Zitat; die Fixrunden-Referenz im Slice-Plan steht zulässig in
der Closure-Notiz (§7), nicht im Code. Commit-Traceability erfüllt
(`LH-FA-RET-004` im Betreff, keine `SPEC-*`/`ARC-*`-ID). Plan-Nachzug
Punkt 2 und §7 sind konsistent mit der tatsächlichen Umsetzung.

**DoD-Checkbox:** bereits `[x]` seit dem ursprünglichen Review-Commit,
kein Nachzug nötig — geprüft und bestätigt.

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, `668f2cd` ist bereits
committet, `make gates`/`make test-integration` real grün.

**Weitere Fixrunde:** nicht nötig.

**Übergabe:** Der Slice ist aus Reviewer-Sicht abgeschlossen geprüft und
kann an die Verifier-Rolle übergeben werden (DoD-/ADR-Konformität gegen
Spec, Plan-vs-Code-Diff — Modul 11). Dieser Report ist Lauf-Beleg und wird
über Läufe hinweg nicht wieder gelesen.
