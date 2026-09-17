# Review-Report: slice-d-check-tracked-modul — 2026-09-17

**Review-Art:** Code — Diff gegen Plan + Konventionen (Modul 10 §Drei
Review-Arten), nicht gegen DoD (Verifier-Aufgabe).

**Gegenstand:** `git diff 314b6cc..HEAD` (Commits `7c9bb1b`, `a276838`,
`feacf52`), Slice `slice-d-check-tracked-modul`, Welle `welle-d-check`.

**Skill:** `.harness/skills/reviewer.md` @ 2026-09-09-Schärfung (vier
repo-spezifische HIGH-Regeln)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-d-check-tracked-modul.md` (Plan)
- `docs/plan/planning/welle-d-check.md` (Wellen-Kontext, das "Mehr")
- `ADR-0072`, `ADR-0075` (zitierte Präzedenzfälle für `hostpaths`)
- `AGENTS.md` §3.9, §3.12, §3.13 (Hard Rules)
- `harness/sensors/docs-check.md`, `harness/README.md` (Träger)
- Eigene Messung: `make doc-tracked` (isoliert), `make docs-check` (Bündel),
  je einmal am unveränderten Stand und einmal mit einer selbst gesetzten
  Mutation (zwei ungestagte, verlinkte Markdown-Dateien unter
  `docs/reviews/`, danach entfernt)

---

## Findings

### F-1 — Modul-Bündel-Erweiterung nicht in alle lebenden Träger nachgezogen

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.13 ("Eine Arbeit, die eine beschriebene
  Eigenschaft bewegt, zieht ihre Träger nach")
- `pfad`: `harness/README.md:114` (**im selben Diff geänderte Datei** — Zeile
  129, 15 Zeilen darunter, behauptet bereits "das Modul läuft bereits im
  `modules:`-Bündel mit", während Zeile 114 in derselben Datei weiterhin nur
  sieben Module nennt: `links, anchors, ids, matrix, versions, structure,
  hostpaths`); zusätzlich `.claude/agents/verifier.md:40` ("die sieben
  `.d-check.yml`-Module …") und `.claude/agents/implementer.md:44`
  (identische Sieben-Module-Aufzählung) — beide **nicht** Teil dieses Diffs.
- `befund`: Der Diff bewegt die beschriebene Eigenschaft "welche Module
  `docs-check`/`make gates` prüft" von sieben auf acht (`tracked`
  aufgenommen), zieht aber nur einen Teil ihrer Träger nach. In
  `harness/README.md` widerspricht die neue Zeile 129 der unveränderten
  Zeile 114 innerhalb derselben Datei und desselben Commits. Eigene
  Gegenprobe: `grep -rn "links, anchors, ids, matrix, versions" .claude/`
  findet zwei weitere, unveränderte Fundstellen mit derselben veralteten
  Sieben-Module-Aufzählung.
- `verifizierbar`: nein — kein Gate prüft Konsistenz einer Freitext-Zahl über
  mehrere Trägerdateien hinweg.
- `klasse`: „Arbeit überholt stehenden Träger" (deckungsgleich mit
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, dort bereits 3× belegt und in
  `AGENTS.md` §3.13 verkörpert)

### F-2 — Zitierter Präzedenzfall stützt die getroffene Aussage nicht

- `kategorie`: HIGH
- `quelle`: Maintainability (Reviewer-Skill-Klasse „Beleg trägt seinen Satz
  nicht")
- `pfad`: `harness/sensors/docs-check.md:174-178` (§Bindung, neu
  hinzugefügter Absatz zu `tracked`)
- `befund`: Der Absatz behauptet „keine eigene ADR, Präzedenzmuster
  `ADR-0072`/`ADR-0075` für `hostpaths`". Beide zitierten ADRs sind aber
  genau die Dokumente, die **weil** die `hostpaths`-Aktivierung eine eigene
  ADR brauchte, geschrieben wurden (`ADR-0072` §Entscheidung Punkt 1: „Die
  Aktivierung … ist ohne die Zitat-Korrektur der `Accepted`-ADRs nicht
  durchführbar. Der Träger ist diese ADR, keine Aktennotiz."). Als Beleg für
  „keine eigene ADR nötig" tragen sie die gegenteilige Aussage — sie zeigen,
  dass eine vergleichbare Bündel-Aufnahme sehr wohl eine ADR auslösen kann,
  wenn sie eine neue Hard Rule erzwingt oder Korrekturen an `Accepted`
  ADRs verlangt.
- `verifizierbar`: nein — Lese-Handlung am zitierten Beleg, kein Gate.
- `klasse`: „Beleg trägt seinen Satz nicht"

### F-3 — Trigger des §4 Start (Architect-Zug) ohne Übergabe-Artefakt übersprungen

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §4 Trigger / Baseline-Regelwerk `modul-08-agentenrollen.md`
  §Rollen-Regeln (Modul 8)
- `pfad`: `docs/plan/planning/in-progress/slice-d-check-tracked-modul.md`
  §4 Trigger, Abschnitt „Start" — Vergleich mit Commit `262bcda` (`next` →
  `in-progress`)
- `befund`: Der Plan macht den Architect-Zug (Planner→Architect, Modul 8)
  explizit zur Start-Bedingung für `next → in-progress`: er soll vor
  Implementierungsbeginn klären, ob die Bündel-Aufnahme eine eigene ADR
  braucht. Die Commit-Message von `262bcda` nennt nur „WIP-Limit frei.
  Erster Slice der Welle welle-d-check." — kein Architect-Verdikt, keine
  Notiz, kein Verweis auf ein solches Artefakt. Der Slice ging ohne dieses
  Übergabe-Artefakt in `in-progress`, und `harness/sensors/docs-check.md`
  trägt anschließend (F-2) bereits eine fertige Schlussfolgerung dazu, ohne
  dass ein Architect sie getroffen hätte.
- `verifizierbar`: nein — Prozess-/Rollenfrage, kein Gate; Beleg ist die
  Commit-Historie selbst.
- `klasse`: „Trigger ohne Übergabe-Artefakt übersprungen"

### F-4 — Realer Präzedenzfall beantwortet die offene ADR-Frage

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `.d-check.yml` Commit `f9e5a3c` (`docs(harness): structure-Modul
  aktiviert`, 2026-09-09)
- `befund`: Eigene Recherche in der Commit-Historie von `.d-check.yml` zeigt
  einen dritten Präzedenzfall neben `hostpaths` (ADR) und dem initialen
  Fünf-Modul-Bestand: Das `structure`-Modul wurde per einfachem,
  unmittelbarem Commit in die `modules:`-Liste aufgenommen — ohne ADR, ohne
  neue Hard Rule, ohne Korrektur an `Accepted`-Inhalten (Commit-Message
  nennt nur die gemessenen 14 zu langen Zellen, behoben durch
  Grenzen-Anpassung). Dieser Fall gleicht `tracked` in allen relevanten
  Merkmalen (0 reale Befunde am gepinnten Digest, keine neue `AGENTS.md`-Hard
  Rule, keine Korrektur an `Accepted`-ADRs) deutlich enger als der
  `hostpaths`-Fall (31 reale Befunde, neue Hard Rule §3.11, erzwungene
  Zitat-Korrektur zweier `Accepted`-ADRs über `ADR-0073`). Auf Basis dieses
  Präzedenzfalls: Die Bündel-Aufnahme von `tracked` braucht **keine** eigene
  ADR — nicht weil `ADR-0072`/`075` das stützen (F-2), sondern weil der
  `structure`-Präzedenzfall zeigt, dass eine reine, folgenlose
  Modul-Aktivierung ohne ADR läuft.
- `verifizierbar`: ja — `git log -p --follow -- .d-check.yml`,
  Commit `f9e5a3c`.
- `klasse`: „Präzedenzfall für Werkzeug-Konfiguration ohne ADR-Bindung"

## Negativbefunde

- geprüft, ohne Befund: `.d-check.yml` — `tracked` korrekt in `modules:`
  aufgenommen, `tracked:`-Block syntaktisch und inhaltlich konsistent mit
  dem Modul-Vertrag (`internal/hexagon/core/rules/tracked.go` im
  `d-check`-Quellbaum geprüft: `ExtractLinks`-basiert, kein Doppelbefund mit
  `links`); `exempt-targets: []` real gegenmessen — `make doc-tracked`:
  875 Datei(en) geprüft, 0 Befund(e); `git status --ignored` zeigt genau
  drei ignorierte Pfade (`.claude/scheduled_tasks.lock`,
  `.harness/state/`, `.tmp/`), keiner davon als Markdown-Link referenziert.
- geprüft, ohne Befund: Mutationsbeleg-Methode — eigene, unabhängige
  Mutation (zwei neue, ungestagte, gegenseitig verlinkte `.md`-Dateien unter
  `docs/reviews/`, danach entfernt) reproduziert exakt `target-untracked`
  (nicht `target-missing`): `make doc-tracked` und das volle
  `make docs-check`-Bündel melden beide denselben einen Befund, `links`/
  `anchors` bleiben dabei still, weil die Zieldatei real existiert. Die im
  Plan beschriebene Methode ist damit ein valider, unabhängig
  reproduzierbarer Test für genau diese Fund-Klasse.
- geprüft, ohne Befund: `harness/sensors/docs-check.md` §Grenze Punkt 3/4 —
  Vollständigkeit gegen die ganze Datei geprüft (Plan-Risiko 3): beide
  Stellen konsistent auf den neuen Bündel-Status umgestellt, keine weitere
  Fundstelle mit „opt-in"/„tracked" in derselben Datei unaktualisiert
  geblieben (abgesehen vom Präzedenz-Belegfehler in F-2).
- geprüft, ohne Befund: `harness/README.md` — `make doc-tracked` korrekt in
  der „kein Gate"-Werkzeuge-Tabelle eingetragen, nicht in der Gates-Tabelle;
  `grep -n "GATE_CHECKS" Makefile harness/mk/*.mk` bestätigt, dass
  `doc-tracked` nirgends an `GATE_CHECKS` hängt.
- geprüft, ohne Befund: Out-of-Scope-Disziplin — keine Änderung an
  `--trace`/RTM-Config im Diff (`slice-d-check-trace-rtm` bleibt unberührt).
- geprüft, ohne Befund: Commit-Traceability — alle drei Commits
  (`7c9bb1b`, `a276838`, `feacf52`) tragen `ADR-0045` im Betreff/Rumpf,
  keine `SPEC-*`/`ARC-*`-Kennung im Betreff.
- geprüft, ohne Befund: DoD-Checkbox-Nachzug in
  `slice-d-check-tracked-modul.md` — LP1–LP3 und der Doku-Update-Punkt
  entsprechen dem tatsächlich committeten Stand.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Arbeit überholt stehenden Träger · Beleg
trägt seinen Satz nicht · Trigger ohne Übergabe-Artefakt übersprungen ·
Präzedenzfall für Werkzeug-Konfiguration ohne ADR-Bindung

## Verdikt

**Merge-blockierend:** ja — zwei HIGH-Findings (F-1, F-2).

**Übergabe:** F-1 und F-3 gehen an den Implementer/Planner (Träger-Nachzug
in `harness/README.md`, `.claude/agents/verifier.md`,
`.claude/agents/implementer.md`; Nachtrag eines Architect-Verdikts zur
ADR-Frage, ggf. rückwirkend). F-2 ist eine Rückkante Review → Plan-Defekt
(Modul 8): die Präzedenz-Zitation in `harness/sensors/docs-check.md`
§Bindung braucht einen echten Architect-Zug, nicht eine vom Implementer
selbst gezogene Schlussfolgerung. F-4 ist informativ für genau diesen Zug —
der `structure`-Präzedenzfall trägt die Antwort „keine eigene ADR nötig"
tatsächlich, im Unterschied zum fälschlich zitierten `hostpaths`-Fall.

Wegen der HIGH-Findings bleibt die DoD-Checkbox „Review durchgeführt, Report
liegt vor" **offen** — es folgt eine Fixrunde (kein Nachzug ohne Fixrunde,
Reviewer-Skill §DoD-Checkbox-Nachzug).

Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11).
