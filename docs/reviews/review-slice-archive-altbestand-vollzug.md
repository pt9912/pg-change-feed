# Review-Report: slice-archive-altbestand-vollzug — 2026-09-18

**Review-Art:** Code + Tooling — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.9.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität ist Verifier-Aufgabe (Modul 11)
und nicht Gegenstand dieses Reports.

**Gegenstand:** `slice-archive-altbestand-vollzug`
(`docs/plan/planning/in-progress/slice-archive-altbestand-vollzug.md`),
Commit-Range `5a29cca..7cb0333` — der erste reale, schreibende
`archive-welle`-Lauf dieses Repos (Schlüssel `altbestand` und
`welle-d-check`).

**Skill:** `.harness/skills/reviewer.md` @ geschärft 2026-09-09 (vier
repo-spezifische HIGH-Regeln)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-archive-altbestand-vollzug.md` (Plan, DoD, §6, §7)
- `docs/plan/adr/0096-altbestand-schluessel-fuer-wellenlosen-archiv-bestand.md` (`Accepted`)
- `docs/plan/planning/welle-archive-altbestand.md`
- `docs/plan/planning/observations/BEO-PGC/externes-werkzeug-committet-ohne-kennung/*`
- Commits `5a29cca..7cb0333` (voller Vorgang inkl. gescheitertem erstem
  Versuch `f46526e`/`ef84d48`)
- `AGENTS.md` §3.3, §3.5, §3.9, §3.12
- `harness/conventions.md` (MR-000, MR-002)
- `.harness/baseline/v6.9.0/templates/docs/plan/planning/archiv-stub-{slice,welle}.template.md`

---

## Findings

### F-1 — Archiv-Stub-Titel bei namensbasierter Slice-Kennung malformt

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Template-Konformität
  `archiv-stub-slice.template.md`)
- `pfad`: `docs/plan/planning/done/welle-d-check/slice-d-check-trace-rtm.md:1`,
  `docs/plan/planning/done/welle-d-check/slice-d-check-tracked-modul.md:1`,
  `docs/plan/planning/done/altbestand/slice-gate-index-konsolidierung.md:1`
- `befund`: Für alle drei archivierten Slices mit namensbasierter Kennung
  (MR-002-Fälle, nicht `slice-<NNN>`) lautet die generierte Titelzeile
  `# slice- — slice-<name>: <Titel>` statt der Template-Form
  `# slice-<Kennung> — <Titel>` — das Kennung-Feld ist leer und der Slug
  wiederholt sich im Titeltext. Der externe `ai-harness-init`-Stand parst die
  Kennung offenbar nach dem Muster `slice-<NNN>` und versagt bei
  Namens-Kennungen. Inhalt, Archiv-Zeiger und Fußzeilen (`Welle:`,
  `Archiviert mit:`, `Geschlossen:`) sind bei allen drei korrekt; betroffen
  ist ausschließlich die erste Zeile. Die Closure-Notiz des Slice-Plans
  (§7) behauptet „ohne inhaltlichen Fehler" — dieser Titel-Defekt
  widerspricht dem in der Sache, auch wenn die Auswirkung klein ist.
- `verifizierbar`: nein — kein Gate prüft Stub-Titel-Form gegen das Template
  (der `structure`-Modul-Blindfleck aus `ADR-0096` deckt genau diesen
  Pfad nicht mehr ab, siehe unabhängig bestätigt unten).
- `klasse`: `archiv-stub-titel-malformed-bei-namensbasierter-kennung`
  (bereits 3× real in diesem einen Lauf — erreicht die Pflege-Schwelle des
  Reviewer-Skills unmittelbar)

### F-2 — Hook-Bypass-Mechanismus widersprüchlich beschrieben

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 (Aussage trägt ihren Ursprung / Beleg trägt
  seinen Satz)
- `pfad`: `docs/plan/planning/in-progress/slice-archive-altbestand-vollzug.md:82-85`
  (§2 LP2) vs.
  `docs/plan/planning/observations/BEO-PGC/externes-werkzeug-committet-ohne-kennung/evidence/slice-archive-altbestand-vollzug.md`
- `befund`: Der Slice-Plan (§2 LP2) behauptet, beide Läufe hätten „den vom
  Werkzeug selbst dokumentierten `--no-verify`-Pfad" genutzt. Der
  Beobachtungs-Beleg beschreibt denselben Vorgang anders und in sich
  uneinheitlich: „wurde `git commit --no-verify` … verwendet (über eine
  prozess-scoped `GIT_CONFIG_*`-Umgebungsvariable, die `core.hooksPath` …
  auf einen leeren Pfad setzt …)" — das nennt eine CLI-Flag-Umgehung und
  eine Config-Override-Umgehung im selben Satz als dieselbe Handlung, ohne
  zu klären, welche der beiden tatsächlich lief. Ein Hook-Bypass hinterlässt
  keine `git`-History-Spur, die das nachträglich entscheidbar machen würde
  — die Aussage ist damit weder durch den Beleg gedeckt noch selbst in sich
  widerspruchsfrei.
- `verifizierbar`: nein — kein Reflog-/Log-Beleg kann rückwirkend
  entscheiden, welcher Mechanismus tatsächlich griff.
- `klasse`: `beleg-widerspricht-sich-bei-hook-bypass-mechanismus`

## Negativbefunde

- geprüft, ohne Befund: Archiv-Integrität `docs/plan/planning/done/altbestand/archiv.zip`
  (141 Dateien: 41 Slices + 1 Plan + 99 Reviews; `unzip -t` fehlerfrei;
  4 Stichproben inhaltlich intakt)
- geprüft, ohne Befund: Archiv-Integrität `docs/plan/planning/done/welle-d-check/archiv.zip`
  (3 Dateien: 2 Slices + 1 Plan; `unzip -t` fehlerfrei; Stichprobe intakt)
- geprüft, ohne Befund: Stub-Form gegen Template (Fußzeilen, keine
  Abschnittsüberschriften, Archiv-Zeiger) — abgesehen von F-1s Titelzeile
- geprüft, ohne Befund: Review-Report-Entfernung — `altbestand` entfernte
  real 99 (`git show --numstat 0caee7f`, reine Löschungen unter
  `docs/reviews/`), deckungsgleich mit den 99 im Zip; `welle-d-check`
  entfernte real 0, korrekt, weil alle 6 zu seinen zwei Slices
  existierenden Review-/Verdikt-/Verify-Reports nachweislich noch anderswo
  zitiert werden (`grep -rl`, je 1–6 Fundstellen)
- geprüft, ohne Befund: Verweis-Nachzug — `make docs-check` reproduziert:
  Exit 0, 0 Befunde, 818 Dateien; kein Rest-Vorkommen alter Pfade in
  `.md`/`Makefile`/`*.mk`/`*.sh`/`*.yml`
- geprüft, ohne Befund: Git-Historie — beide Move-Commits (`b96e3e7`,
  `a39bbcd`) 0 Insertions/0 Deletions (reine Renames); beide
  Inhalts-Commits (`0caee7f`, `a787203`) betreffen ausschließlich Archiv,
  Stubs, `AGENTS.md`-Zeile (Kennung, Zitat-Ziel) und nachgezogene
  Referenzen — kein `-A`-Scope-Creep
- geprüft, ohne Befund: gescheiterter erster Versuch (`f46526e`) und sein
  `git revert` (`ef84d48`) — exakt komplementäre Rename-Paare, keine
  Restspuren
- geprüft, ohne Befund: `commit-traceability` — real reproduziert, aktuell
  rot (4 `commit-untraceable`-Funde für `0caee7f`/`a39bbcd`/`a787203`/
  `b96e3e7` im `HEAD~5..HEAD`-Fenster), exakt wie im Slice-Plan §2 LP3
  dokumentiert und als vorübergehend erwartet
- geprüft, ohne Befund: §6-Risiko-Ausgänge — alle drei „entfallen"
  unabhängig bestätigt: (1) `make docs-check` 0 Befunde, (2) `structure`-Modul
  nutzt tatsächlich einen flachen Glob (`docs/plan/planning/done/slice-*.md`,
  kein `**`, `.d-check.yml:153`), archivierte Stubs liegen außerhalb, (3)
  die `[haenger]`-Bereinigung (`a12c62a`, `22aacd9` u. a.) liegt vollständig
  vor dem Start dieses Slices (`6dab665`)
- geprüft, ohne Befund: ADR-0096-Konformität — Schlüsselwahl `altbestand`
  umgesetzt wie entschieden, die zwei Hand-Dokumente
  (`altbestand.md`/`altbestand-results.md`) vor dem realen Lauf angelegt
  (`587f24d` vor `b96e3e7`), `welle-d-check` unangetastet als eigener,
  zweiter Schlüssel
- geprüft, ohne Befund: Ergebnisnotiz-Semantik — `altbestand-results.md`
  bleibt flach in `done/`, nicht archiviert (Template-Vorgabe eingehalten)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:**
`archiv-stub-titel-malformed-bei-namensbasierter-kennung` ·
`beleg-widerspricht-sich-bei-hook-bypass-mechanismus`

## Verdikt

**Merge-blockierend:** nein für die DoD-Zeile „Review durchgeführt" (0
HIGH, Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde) — **aber** beide
MEDIUM-Findings sollten vor der Welle-Closure `welle-archive-altbestand`
behoben werden: F-1 durch eine kleine Korrektur der drei Stub-Titelzeilen
(betrifft nur den außerhalb der `archiv.zip` liegenden Stub-Text, nicht das
eingefrorene Zip), F-2 durch eine Präzisierung des Beobachtungs-Belegs
oder eine erneute, eindeutige Feststellung des tatsächlichen Mechanismus.
Keine der beiden Findings hat einen Rollen-Widerspruch oder ein drittes
Auftreten über Reviews hinweg ausgelöst — kein Architect-Konfliktpfad
nötig.

**Übergabe:** Beide Findings gehen an den Planner (kein
Implementer-Rückkanal — reiner Planning-/Tooling-Slice ohne Produktcode);
die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den
Beobachtungs-Zähler. Dieser Report ist Lauf-Beleg; DoD-Substanz und
Closure-Trigger prüft der Verifier separat (Modul 11).
