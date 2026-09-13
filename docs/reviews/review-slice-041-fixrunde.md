# Review-Report: slice-041 (Fixrunde) — 2026-09-13

**Review-Art:** Code — gezielte Bestätigungsprüfung der vier Findings
(F-1 HIGH, F-2 MEDIUM, F-3 LOW, F-4 LOW) aus
`docs/reviews/review-slice-041.md` (Modul 10), **kein** vollständiges
Re-Review des Slice. Geprüft gegen den ursprünglichen Befund-Text,
`ADR-0052`, `SPEC-016` und `AGENTS.md` §3.7 (Kommentar-Disziplin).

**Gegenstand:** Commit `35a5279` (`fix(bootstrap): Slice-Chronik aus
Kommentaren entfernt, Precedence-Entscheidung dokumentiert (ADR-0052,
review-slice-041 F-1..F-4)`) auf `4e4c7bb`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/reviews/review-slice-041.md` (F-1..F-5, vollständig, eigener
  vorheriger Lauf)
- `git show 35a5279 --stat` und `git show 35a5279` (vollständiger Diff,
  alle drei geänderten Dateien)
- `AGENTS.md` §3.7 (Kommentar-Disziplin), Baseline-Regelwerk
  `grundlagen-harness-dateien.md` §Was ein Kommentar trägt
- `docs/plan/planning/in-progress/slice-041-yaml-konfigurationsdatei.md`
  §3 (Plan-Nachzug, vollständig)
- eigene reale Werkzeugläufe (nicht aus dem Commit-Text übernommen):
  `make gates`, `make test`
- eigenständige `grep -rni "slice"` über `internal/bootstrap/config_file.go`
  und `internal/bootstrap/config_file_internal_test.go` (Vollständigkeits-
  Check über den Diff hinaus, nicht auf die vom Implementer genannten
  Zeilen beschränkt)

---

## F-1 — Slice-Chronik in `mergeTables`-Doc-Kommentar (HIGH)

**Verdikt: an der gemeldeten Stelle behoben.**

- Der Kommentar lautet jetzt: „`mergeTables` trägt die
  `tables`-Merge-Precedence (`SPEC-016`): eine gesetzte `CDC_TABLES`
  schlägt die gesamte Datei-`tables`-Mapping vollständig …" — keine
  Slice-Kennung mehr, Verweis auf `SPEC-016`/`ADR-0052` (stabile,
  kanonische Quellen statt einer Adresse, die bei Welle-Closure ins
  Archiv wandert).
- `grep -n "slice-[0-9]\|Slice-Plan" internal/bootstrap/config_file.go
  internal/bootstrap/config_file_internal_test.go` → kein Treffer mehr.

**Aber — bei eigenständiger, über den gemeldeten Fundort hinausgehender
Prüfung derselben Datei gefunden (nicht Teil des Findings-Umfangs dieser
Fixrunde, aber derselbe Verstoß, unverändert seit `af94b22`):**

`internal/bootstrap/config_file.go:111–118`, Doc-Kommentar von
`ConfigFromEnvAndFile`:

```
// ConfigFromEnvAndFile ist der verdrahtete Zugriffsweg (`ADR-0052`
// Entscheidung 5, verdrahtet in `cmd/pg-change-feed/main.go`):
// `CDC_CONFIG_FILE` leer/unbenannt trägt exakt den heutigen
// `ConfigFromEnv`-Pfad, unverändert — jeder bestehende Env-only-Aufruf
// zeigt identisches Verhalten wie vor diesem Slice. Ist die Variable
// gesetzt, aber die Datei unter diesem Pfad nicht ladbar, ist das ein
// `ErrConfiguration`-Fehler — kein stiller Fallback auf Env-only
// (`ADR-0052` Entscheidung 5).
```

„zeigt identisches Verhalten wie vor diesem Slice" ist ein
Vorher/Nachher-Vergleich (Chronik-Sprache), derselbe Verstoß gegen
`AGENTS.md` §3.7 wie der ursprüngliche F-1-Fund — nur ohne explizite
Slice-Kennung. Der Satz trägt keine der zulässigen Kommentar-Klassen
(Zusage · Kopplung · Abgrenzung · Rang-Zeiger · Grenze): Er beschreibt
nicht den Ist-Zustand von `ConfigFromEnvAndFile`, sondern einen Vergleich
zu einem Zustand vor dieser Arbeit. Die Zusage selbst („Env-only-Pfad
unverändert") ist mit `git show af94b22 -- .../wiring.go` bereits
unabhängig verifiziert (Negativbefund im Ursprungs-Report) — der Satz
ist damit auch inhaltlich redundant zur eigentlichen Garantie, die zwei
Sätze vorher schon steht.

Diese Stelle war **nicht** Teil von F-1..F-4 und ist damit kein
Nichterfüllen der beauftragten Fixrunde — aber sie ist dieselbe
Fehler-Klasse „Slice-Chronik in Go-Quellcode-Kommentar" (bereits 2×
in `review-slice-041.md` als F-1/F-4 gezählt), jetzt ein drittes
Vorkommen in demselben Slice/derselben Datei. Das erreicht die
Steering-Loop-Schwelle aus `.harness/skills/reviewer.md` §Pflege
(„Bei dreimaligem Auftreten desselben Findings … gibt es eine Fitness
Function, die das prüfen würde?") und aus `grundlagen-klassifikation.md`
§Steering Loop (1× notieren · 2× Symptom · 3× Lücke).

- `kategorie`: HIGH (neues Finding dieser Fixrunde, gleiche Klasse wie
  F-1)
- `quelle`: `AGENTS.md` §3.7, Baseline-Regelwerk
  `grundlagen-harness-dateien.md` §Was ein Kommentar trägt
- `pfad`: `internal/bootstrap/config_file.go:115`
- `verifizierbar`: ja — Code-Inspektion, Diff gegen `af94b22`
- `klasse`: „Slice-Chronik in Go-Quellcode-Kommentar" (3. Vorkommen)

## F-2 — Tables-Merge-Precedence ohne Folge-ADR (MEDIUM)

**Verdikt: prozessual geschlossen, kein neuer Fehler.**

- Neuer Absatz „Rollen-Einordnung dieser Entscheidung (`review-slice-041.md`
  F-2, MEDIUM)" in §3 des Slice-Plans, geprüft vollständig.
- Der Absatz übernimmt mein ursprüngliches Verdikt korrekt und
  unverfälscht: kein Rollen-Widerspruch (die offene Frage wurde in §6
  Risiko 1 selbst benannt, nicht bestritten), Entscheidung wird als
  hinreichend begründetes Implementierungsdetail akzeptiert (Konsistenz
  zur Ganzwert-Precedence der Stringfelder), mit explizitem
  Folge-ADR-Vorbehalt, falls die Precedence-Frage künftig erneut
  strittig wird.
- Das deckt sich mit meiner eigenen Einschätzung aus `review-slice-041.md`
  Zeilen 292–296 wortwörtlich in der Sache: Ich hatte selbst notiert,
  dass die Architect-Sequenz aus Modul 8 „nicht zwingend" greift, weil
  kein bestrittener Rollen-Widerspruch vorliegt, und dass Annahme als
  Implementierungsdetail „mit Vermerk in der Closure-Notiz" eine
  legitime Option ist. Der Implementer hat diesen Vermerk vorgezogen in
  den Plan-Nachzug gesetzt statt ihn erst bei Closure zu schreiben — das
  ist strenger, nicht schwächer, als gefordert.
- Kein Widerspruch, keine neue Bewertung nötig.

## F-3 — Doppelte DSN-Ablehnung, Kommentar-Klarheit (LOW)

**Verdikt: behoben.**

- Neuer Kommentar über `ConfigFromFile` (Zeilen 68–82) beschreibt jetzt
  korrekt die tatsächliche Reihenfolge: der explizite
  `forbiddenFileDSNKeys`-Check „im aktuellen Kontrollfluss immer
  zuerst", `KnownFields(true)` „würde sie ebenfalls ablehnen, sollte der
  explizite Check je entfallen — im jetzigen Kontrollfluss ist dieser
  Pfad für die drei DSN-Schlüssel nicht erreichbar, weil der explizite
  Check vorher zurückkehrt".
- Gegen den Code verifiziert (`config_file.go:93–107`): der
  `forbiddenFileDSNKeys`-Loop läuft auf `raw` vor `decoder.Decode`, mit
  frühem `return` bei Treffer — deckt sich exakt mit der neuen
  Kommentar-Aussage. Keine Suggestion mehr von zwei aktuell gleichzeitig
  wirksamen Linien.
- Kein Funktionsfehler, keine Verhaltensänderung — reine
  Kommentar-Präzisierung, wie gefordert.

## F-4 — Testkommentar referenziert „Slice-Plan" (LOW)

**Verdikt: behoben.**

- `TestMergeConfigTabellenCDCTablesSchlaegtDatei`-Doc-Kommentar in
  `config_file_internal_test.go:202–207` lautet jetzt: „trägt die
  `tables`-Merge-Precedence (`SPEC-016`): …" — kein Verweis auf „des
  Slice-Plans" mehr.
- `grep -n "Slice-Plan" internal/bootstrap/*_test.go` → kein Treffer.

## Gate- und Testläufe (real ausgeführt, nicht übernommen)

- `make gates`: `baseline-verify` (54 Dateien OK), `d-check`/`docs-check`
  (328 Dateien, 0 Befunde), `commit-traceability` (5 Commits im aktuellen
  `HEAD~5..HEAD`-Fenster, „Betreffs ohne Struktur-ID" — OK; der zuvor
  bekannte rote Befund aus `4e4c7bb` ist wie erwartet aus dem
  5-Commit-Fenster gerutscht, siehe Registereintrag
  `BEO-PGC/commit-traceability-kein-vorab-hook`, nicht Gegenstand dieser
  Prüfung), `a-check` (0 Befunde) — alle grün.
- `make test`: alle Pakete grün, inklusive
  `github.com/pt9912/pg-change-feed/internal/bootstrap` (1.031s, mit
  `-race`).

## Negativbefunde

- geprüft, ohne Befund: kein Funktions-/Verhaltenscode geändert — der
  Fix-Commit ändert ausschließlich Kommentare (Go) und den Plan-Nachzug
  (Markdown); `git show 35a5279` bestätigt: keine Zeile außerhalb von
  `//`-Kommentaren und der Plan-Datei betroffen.
- geprüft, ohne Befund: keine neue Slice-Kennung oder sonstige Chronik
  (`seit `, `nur noch`, `jetzt `) in den geänderten Kommentarblöcken
  außer dem in F-1 benannten Altbestand.
- geprüft, ohne Befund: Commit-Betreff trägt `ADR-0052`, kein
  `SPEC-*`/`ARC-*` im Betreff (`fix(bootstrap): Slice-Chronik aus
  Kommentaren entfernt, Precedence-Entscheidung dokumentiert
  (ADR-0052, review-slice-041 F-1..F-4)`); die Erwähnung
  „review-slice-041 F-1..F-4" ist keine Struktur-ID im Sinne des
  Musters.

## Summary

| Finding | Verdikt |
|---|---|
| F-1 (HIGH) | an gemeldeter Stelle behoben — **neues, verwandtes Vorkommen derselben Klasse gefunden** (siehe oben) |
| F-2 (MEDIUM) | prozessual geschlossen, keine Beanstandung |
| F-3 (LOW) | behoben |
| F-4 (LOW) | behoben |

**Neue Findings dieser Fixrunde:** 1× HIGH — drittes Vorkommen der
Klasse „Slice-Chronik in Go-Quellcode-Kommentar"
(`internal/bootstrap/config_file.go:115`, `ConfigFromEnvAndFile`),
außerhalb des ursprünglichen F-1-Fundorts, bei eigenständiger Prüfung
über den gemeldeten Diff hinaus gefunden. Erreicht die
Steering-Loop-Schwelle (3×: F-1, F-4, dieses Vorkommen).

## Verdikt

**Merge-blockierend:** ja — F-1..F-4 aus `review-slice-041.md` sind an
den ursprünglich gemeldeten Stellen real behoben (Kommentare korrigiert,
kein Funktionscode berührt, alle Gates und Tests grün), aber die
eigenständige Nachprüfung deckt ein drittes, unbehandeltes Vorkommen
derselben HIGH-Klasse in derselben Datei auf. Der Slice darf aus
Reviewer-Sicht **noch nicht** nach `done/` übergehen, bevor
`internal/bootstrap/config_file.go:115` ebenfalls korrigiert ist
(Ist-Zustand statt Vorher/Nachher-Vergleich, ohne Slice-/Vorgangsbezug).

**Empfehlung für die Planner-Rolle bei Closure:** Da dies das dritte
Vorkommen der Klasse „Slice-Chronik in Go-Quellcode-Kommentar" ist
(F-1, F-4, dieser Fund), ist die Steering-Loop-Schwelle erreicht — ein
Eintrag im Beobachtungs-Register (`BEO-PGC/…`) oder eine geschärfte
Prüfregel (z. B. `grep -rni "vor diesem\|vor dies\|seit \|nur noch\|jetzt "`
als Teil des 8-Schritt-Workflows vor Handoff) ist fällig, nicht nur eine
punktuelle Korrektur.

**Übergabe:** Eine weitere, kleine Fixrunde ist nötig (ein Kommentar,
kein Funktionscode), bevor Verifikation (Modul 11) sinnvoll greift.
Dieser Report ist ein Lauf-Beleg und wird über Läufe hinweg nicht wieder
gelesen; die Summary-Zeile speist bei Bedarf den Closure-Eintrag
(Modul 5). Er ersetzt keine Verifikation gegen DoD/Spec.
