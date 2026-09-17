# Review-Report: slice-041 (Fixrunde 2) — 2026-09-13

**Review-Art:** Code — gezielte Bestätigungsprüfung des dritten,
unbehandelten Vorkommens der Kommentar-Chronik-Klasse
(`docs/reviews/review-slice-041-fixrunde.md`, HIGH-Finding zu
`internal/bootstrap/config_file.go:115`), **kein** vollständiges
Re-Review des Slice. Geprüft gegen `AGENTS.md` §3.7 (Kommentar-
Disziplin) und `ADR-0052`.

**Gegenstand:** Commit `8ceee6b` (`fix(bootstrap): drittes
Chronik-Vorkommen in Kommentar behoben (ADR-0052,
review-slice-041-fixrunde)`) auf `25b7a6c`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/reviews/review-slice-041-fixrunde.md` (vollständig, eigener
  vorheriger Lauf — nennt den Fund im Detail)
- `AGENTS.md` §3.7 (Kommentar-Disziplin), Baseline-Regelwerk
  `grundlagen-harness-dateien.md` §Was ein Kommentar trägt
- `git show 8ceee6b` (vollständiger Diff, einzige geänderte Datei)
- `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar/`
  (vollständig — `observation.md`, `state.md`, beide `evidence/`-Dateien),
  angelegt vom Planner in Commit `7fc8b90`
- eigene reale Werkzeugläufe (nicht aus dem Commit-Text übernommen):
  `make gates`, `make test`
- eigenständiger `grep -niE` über `internal/bootstrap/config_file.go`
  **und** `internal/bootstrap/config_file_internal_test.go` gegen ein
  breiteres Chronik-Vokabular als der Implementer-Commit nannte
  ("slice", "vor diesem", "seit ", "nur noch", "jetzt ", "heutigen",
  "nicht mehr", "zuvor", "früher", "vormals")

---

## Prüfung 1 — Ist der gemeldete Kommentar jetzt Ist-Zustand-Beschreibung?

**Verdikt: ja, behoben.**

Vorher (Zeilen 111–118 vor `8ceee6b`):

```
// `CDC_CONFIG_FILE` leer/unbenannt trägt exakt den heutigen
// `ConfigFromEnv`-Pfad, unverändert — jeder bestehende Env-only-Aufruf
// zeigt identisches Verhalten wie vor diesem Slice.
```

Nachher:

```
// `CDC_CONFIG_FILE` leer/unbenannt delegiert vollständig an
// `ConfigFromEnv` — derselbe Env-only-Pfad, keine Datei-Berührung.
```

Der neue Satz beschreibt, was `ConfigFromEnvAndFile` bei leerem
`CDC_CONFIG_FILE` tut (Delegation, keine Datei-Berührung) — kein
Vorher/Nachher-Vergleich, keine Slice-Referenz. Er trägt die
Kommentar-Klasse *Zusage* (was bei leerer Variable gilt). Der Rest des
Blocks (Fehlerpfad bei nicht ladbarer Datei, `ADR-0052`-Zeiger) ist
unverändert und war bereits vorher unbeanstandet.

Der Diff selbst ist minimal und ausschließlich Kommentar: `git show
8ceee6b --stat` zeigt eine Datei, 4 Insertions/5 Deletions, keine
Funktionscode-Zeile betroffen.

## Prüfung 2 — Vollständigkeits-Grep über beide Dateien

**Verdikt: keine weiteren Vorkommen.**

```
grep -niE "slice|vor diesem|seit |nur noch|jetzt |heutigen|nicht mehr|zuvor|früher|vormals" \
  internal/bootstrap/config_file.go internal/bootstrap/config_file_internal_test.go
```

→ kein Treffer in beiden Dateien (breiteres Vokabular als der
Implementer-Commit, der nur "vor diesem", "seit ", "nur noch", "jetzt ",
"heutigen" und Slice-Kennungen nannte — hier zusätzlich "nicht mehr",
"zuvor", "früher", "vormals" ergänzt, ebenfalls ohne Treffer).

Volltext-Review beider Dateien bestätigt das Grep-Ergebnis: Die
verbleibenden Kommentare mit Zeitbezug sind Kopplungs-Aussagen zwischen
Code-Pfaden, keine Vorher/Nachher-Vergleiche zu einem früheren
Repo-Zustand — z. B. Zeile 121 in
`config_file_internal_test.go` ("zeigt identisches Verhalten zu
`ConfigFromEnv`") vergleicht zwei *aktuell koexistierende* Funktionen,
nicht einen Zustand vor/nach einem Slice, und ist damit keine Chronik im
Sinne von `AGENTS.md` §3.7. Ebenso die Formulierungen "im aktuellen
Kontrollfluss" / "im jetzigen Kontrollfluss" in `config_file.go:74–78`
(Ist-Zustand-Beschreibung einer Bedingung, kein Vorher/Nachher-Vergleich
zu einem früheren Slice-Stand) — unverändert seit vor dieser Fixrunde
und nicht Gegenstand des gemeldeten Findings, bei dieser Prüfung aber
mitgeprüft und für unauffällig befunden.

## Prüfung 3 — Gate- und Testläufe (real ausgeführt)

- `make gates`: `baseline-verify` (54 Dateien OK), `d-check`/`docs-check`
  (333 Dateien, 0 Befunde), `commit-traceability` (5 Commits im
  `HEAD~5..HEAD`-Fenster, Betreffs ohne Struktur-ID — OK), `a-check`
  (0 Befunde) — alle grün.
- `make test`: alle Pakete grün, inklusive
  `github.com/pt9912/pg-change-feed/internal/bootstrap` (1.038s, mit
  `-race`).

## Prüfung 4 — Beobachtungs-Register-Zähler

**Verdikt: nachvollziehbar und korrekt angewendet.**

`BEO-PGC/slice-chronik-in-code-kommentar/` trägt zwei Evidence-Dateien:

- `evidence/review-slice-041.md` — **ein** Vorgang (der Review-Report
  `review-slice-041.md`), der zwei Funde innerhalb desselben Laufs bündelt
  (F-1 HIGH, F-4 LOW). Nach der Register-Regel "Ein Vorgang zählt
  einmal — zwei Funde im selben Vorgang sind eine Gelegenheit, kein
  zweites Auftreten" ([Modul 6](../../.harness/baseline/v6.9.0/regelwerk/modul-06-roadmap.md#das-beobachtungs-register-modul-6))
  ist das korrekt **eine** Evidence-Datei für F-1+F-4 zusammen, nicht
  zwei.
- `evidence/review-slice-041-fixrunde.md` — ein zweiter, eigenständiger
  Vorgang (der Fixrunden-Review-Report, ein anderer abgeschlossener
  Lauf), der das dritte Vorkommen fand.

Zähler (abgeleitet aus der Zahl der Evidence-Dateien): **2** — unter der
3×-Schwelle aus [Modul 6 §Das Beobachtungs-Register](../../.harness/baseline/v6.9.0/regelwerk/modul-06-roadmap.md#das-beobachtungs-register-modul-6).
Damit ist `state.md`s Eintrag "Zustand: offen … kein Ausgang fällig"
korrekt: Ein dritter *Vorgang* (dritte Evidence-Datei) läge erst vor,
wenn eine künftige, unabhängige Prüfung ein weiteres Vorkommen fände.
Der mit `8ceee6b` behobene Fund selbst erzeugt **keine** dritte
Evidence-Datei — er ist der Inhalt der bereits vorhandenen
`evidence/review-slice-041-fixrunde.md` und wurde jetzt behoben, nicht
neu gefunden. Die historische `slice-018`-Korrektur bleibt zu Recht
unter "Historie vor der Registrierung — benannt, nicht gezählt"
(Präzedenzfall `BEO-PGC/plan-nachzug`, vor der Registrierung des
Registers selbst).

Der Zähler-Stand und die Nichtauslösung eines Ausgangs (verkörpert /
geplant / gestrichen) sind damit für den Reviewer nachvollziehbar und
formal korrekt — urteilsfrei geprüft (Vorgangs-vs-Fund-Zählung,
Schwellenvergleich), kein Widerspruch zum Planner-Verdikt.

## Negativbefunde

- geprüft, ohne Befund: `git show 8ceee6b` ändert ausschließlich
  Kommentarzeilen in `internal/bootstrap/config_file.go`, keine
  Funktionscode-Zeile.
- geprüft, ohne Befund: keine weitere Chronik-Sprache in
  `config_file.go` oder `config_file_internal_test.go` (Grep gegen
  erweitertes Vokabular, siehe Prüfung 2).
- geprüft, ohne Befund: Register-Zähler (2×) korrekt aus den
  vorhandenen Evidence-Dateien abgeleitet, keine dritte Datei fällig
  durch diesen Fix.
- geprüft, ohne Befund: Commit-Betreff von `8ceee6b` trägt `ADR-0052`,
  kein `SPEC-*`/`ARC-*` im Betreff.

## Summary

| Prüfung | Verdikt |
|---|---|
| Gemeldeter Kommentar (`config_file.go:111–118`) | behoben — Ist-Zustand, kein Vorher/Nachher-Vergleich, keine Slice-Referenz |
| Vollständigkeits-Grep (erweitertes Vokabular) | keine weiteren Vorkommen in beiden Dateien |
| Gates (`make gates`) | grün |
| Tests (`make test`, inkl. `internal/bootstrap` mit `-race`) | grün |
| Register-Zähler `BEO-PGC/slice-chronik-in-code-kommentar` | 2× — unter der 3×-Schwelle, korrekt kein Ausgang fällig |

**Verdikt:** Die Fixrunde ist bestätigt. Kein neues Finding. Der
ursprünglich blockierende HIGH-Befund aus
`review-slice-041-fixrunde.md` ist real behoben, ohne dass die Prüfung
über den gemeldeten Umfang hinaus weitere Vorkommen der Klasse
"Slice-Chronik in Go-Quellcode-Kommentar" findet. Reviewer-seitig keine
Einwände gegen einen Übergang nach Verifikation (Modul 11).

Dieser Report ist ein Lauf-Beleg und wird über Läufe hinweg nicht wieder
gelesen; die Summary-Zeile speist bei Bedarf den Closure-Eintrag
(Modul 5). Er ersetzt keine Verifikation gegen DoD/Spec.
