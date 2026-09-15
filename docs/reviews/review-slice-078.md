# Review-Report: slice-078 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-078`, Diff `3f4c747..f6008e7` (6 Züge:
`df47282` · `9ba35e8` · `f46c610` · `aab7aa8` · `540861d` · Implementer;
`e37694c` Planner-Plan-Nachzug; `f6008e7` §3.5-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14 — die Regeln `Handbuch-Versionshistorie` und
`Neue Betreiber-Oberfläche ohne Handbuch-Zug` liegen **vor** diesem Lauf).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-078-hostpfade-entfernen.md`
  vollständig (§1–§8), einschließlich des §2/§3-Nachzugs aus `e37694c`
- [`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
  vollständig (§Kontext, §Entscheidung Punkt 1–5 samt §3.11-Entwurf,
  §Verglichene Alternativen, §Konsequenzen, §Fitness Function,
  §Re-Evaluierungs-Trigger, §Geschichte)
- [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  vollständig (die Klassen-Grenze in §Entscheidung Punkt 1 und der
  vorgegebene §3.5-Wortlaut in Punkt 2)
- [`ADR-0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md)
  vollständig (Supersedes-Klausel — welcher Teil `ADR-0072` abgelöst ist)
- `AGENTS.md` §3.5 (nachgezogen), §3.11 (neu), §3.6, §3.7, §5 ·
  `harness/conventions.md` (MR-000)
- `.d-check.yml` (`modules`, `scan`), `d-check.mk` (`DCHECK_DIGEST`),
  gepinntes Modul-Image `pt9912/d-check` `v0.75.0`
- `docs/reviews/architect-verdict-hostpaths-aktivierung-ohne-ausnahme.md`
  (die drei Verdikte, die die Korrektur tragen)
- `harness/sensors/docs-check.md`, `harness/README.md` §Sensors

---

## Findings

### F-1 — Die Zitat-Korrektur ändert `§Entscheidung` von `ADR-0072`; `§3.5` und `ADR-0073` Punkt 1 schließen genau diesen Abschnitt aus

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.5 (nachgezogen durch
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung Punkt 2) · `ADR-0073` §Entscheidung Punkt 1 (Klassen-Grenze:
  „**keine** Zitat-Korrektur, wenn sie … bei einer ADR: §Entscheidung …")
- `pfad`: `docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md:173-190`
  (§Entscheidung, Unterabschnitt `Entwurf der Hard Rule §3.11`) · Commit
  `aab7aa8`
- `befund`: In `ADR-0072` `§Entscheidung` wurde der Unterabschnitt
  `Entwurf der Hard Rule §3.11` in-place geändert — sowohl die
  Falsch-/Richtig-Beispiele als auch ihr erklärender Klammertext (Zeile 173:
  „der wörtliche Befund …" → „Form: der Pfad beginnt mit dem Wurzel-Segment
  …"). `ADR-0072` Punkt 5 erklärt diesen Unterabschnitt ausdrücklich zum
  „Teil dieser Entscheidung"; `AGENTS.md` §3.5 und `ADR-0073` Punkt 1 nehmen
  `§Entscheidung` von der Zitat-Korrektur aus, und der geänderte Klammertext
  ist keine Form eines Verweises. Folge im selben Dokument: `§Entscheidung`
  Punkt 3 führt weiter die alte besitzer-qualifizierte Form, während das
  Beispiel darunter die Hausform zeigt.
- `verifizierbar`: nein — die Klassenzugehörigkeit ist ein Urteil
  (`ADR-0073` §Fitness Function, Zeile 167)
- `klasse`: „Zitat-Korrektur überschreitet die Klassengrenze (§Entscheidung)"

### F-2 — Zeilen-/Bereichs-Lokatoren bleiben stehen, obwohl `ADR-0072` Punkt 3 und das Verdikt-Beispiel ihre Ersetzung verlangen

- `kategorie`: MEDIUM
- `quelle`: `ADR-0072` §Entscheidung Punkt 3 („Zeilen-/Bereichs-Lokatoren
  … werden durch den **stabilen benannten Anker** ersetzt") ·
  `docs/reviews/architect-verdict-hostpaths-aktivierung-ohne-ausnahme.md:179-190`
  (Verdikt 3 und dessen vorher/nachher-Beispiel, das den Lokator fallen lässt)
  · [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung Punkt 1 („Zeilen-/Bereichs-Lokatoren" stehen **in** der Klasse)
- `pfad`: `docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md:42,44,75,128,154`
  · `docs/plan/planning/done/slice-050-performance-benchmark-infrastruktur.md:274`
- `befund`: Die Lokatoren (`Zeilen 69–93`, `Zeilen 35–39`, `Zeile 84`,
  `Zeile 35–38`) stehen unverändert im Bestand, teils **neben** dem bereits
  vorhandenen benannten Anker (`Stage coverage`, `bench:`-Target). `ADR-0074`
  löst `ADR-0072` Punkt 3 ab, restituiert die Lokator-Klausel in seinem
  Ersatztext aber nicht und führt sie auch nicht unter „Der Rest … bleibt"
  auf — die Disposition der Klausel ist damit nirgends ausgesprochen, während
  `ADR-0072:120` sie weiter wörtlich führt. Die vorgebrachte Begründung
  („fällt nicht unter die in-place-Klasse") steht gegen `ADR-0073` Punkt 1,
  das Lokatoren ausdrücklich in die Klasse aufnimmt.
- `verifizierbar`: nein
- `klasse`: „Zeilen-Lokator nicht durch benannten Anker ersetzt"

### F-3 — „kein Host-Pfad im Repo" (Plan) und „nur Prosa + Inline-Code" (§3.11) beschreiben dieselbe Regel verschieden

- `kategorie`: MEDIUM
- `quelle`: Plan §1 („der Modul lässt Fenced-Blöcke frei, die Regel dieses
  Repos nicht") · `AGENTS.md` §3.11 („in Prosa oder Inline-Code" / „Fenced-
  Code-Blöcke … dort sind Beispiel-Pfade erlaubt") · `ADR-0072`
  §3.11-Entwurf (gleicher Wortlaut)
- `pfad`: `docs/plan/planning/in-progress/slice-078-hostpfade-entfernen.md:39,44-50`
  · `AGENTS.md:317,341`
- `befund`: Plan §Ziel und §1 sagen, das Repo entferne den Pfad auch dort, wo
  der Sensor ihn nicht liest (Fences) — „sie stehen nicht mehr da". Die
  ausgelieferte Hard Rule §3.11 grenzt ihre eigene Aussage auf „Prosa oder
  Inline-Code" ein und nennt Fenced-Beispiele ausdrücklich erlaubt; ihr
  ADR-Entwurf sagt dasselbe. Die beiden Stellen legen dieselbe Regel
  unterschiedlich weit aus; welche gilt, ist nicht deklariert.
- `verifizierbar`: nein
- `klasse`: „Regel-Reichweite für Fenced-Blöcke widersprüchlich"

### F-4 — Der nachgezogene `§3.5`-Wortlaut weicht von der in `ADR-0073` Punkt 2 vorgegebenen Liste ab

- `kategorie`: LOW
- `quelle`: [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung Punkt 2 (vorgegebener §3.5-Wortlaut)
- `pfad`: `AGENTS.md:144-150`
- `befund`: Vorgegeben war „das Zitat- und Verweisgerüst — host-lokale Pfade,
  **Linkziele, Zeilen-Lokatoren** — bei unverändertem Referenten"; ausgeliefert
  ist „host-lokale Pfade, **gebrochene Linkziele, Formfehler der Zitation**".
  `Zeilen-Lokatoren` fehlen, `Formfehler der Zitation` ist weiter gefasst als
  `ADR-0073`s „Form einer gebrochenen Referenz"; die Klasse wird damit an
  dieser Stelle enger (Klassen-Kategorie entfällt) **und** unschärfer. Der
  Satz verweist auf `ADR-0073`, die Definition bleibt dort gebunden.
- `verifizierbar`: nein
- `klasse`: „Nachgezogener Hard-Rule-Wortlaut weicht von der ADR-Vorgabe ab"

### F-5 — `harness/sensors/docs-check.md` erklärt die Config zur Deklaration und führt die Modulliste daneben

- `kategorie`: INFO
- `quelle`: Maintainability (`AGENTS.md` §3.7, Zwei-Quellen-Disziplin)
- `pfad`: `harness/sensors/docs-check.md:18-21`
- `befund`: Der Absatz sagt „Die Module … stehen in `.d-check.yml`; die
  Konfiguration ist die Deklaration dieses Vertrags, nicht dieses Dokument" und
  listet im nächsten Satz den Ist-Zustand derselben Modulliste erneut auf. Die
  Liste lebt damit in drei Dateien (`.d-check.yml`, `harness/README.md`
  §Sensors, hier). `AGENTS.md` §3.11 vermeidet genau diese Wiederholung für
  die Präfixliste.
- `verifizierbar`: nein
- `klasse`: „Modulliste doppelt geführt"

### F-6 — `ADR-0072` spricht weiter von 31 Stellen, Plan und Messung von 42

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md:59`
  · `docs/plan/planning/in-progress/slice-078-hostpfade-entfernen.md:39,114`
- `befund`: `ADR-0072` führt durchgehend 31 (die Modul-Befunde zum
  Entscheidungszeitpunkt), der Plan und die Messung 42 (Fenced-,
  Nicht-Markdown- und danach selbst erzeugte Stellen). Der Unterschied ist in
  `ADR-0072` nicht abgeglichen; „Korrektur **aller** 31 Stellen" liest sich
  dort als Gesamtumfang.
- `verifizierbar`: ja — der `hostpaths`-Modul-Aufruf und der repo-weite Grep
  belegen die Ist-Zahl 0 (siehe §Eigene Messungen)
- `klasse`: „Zähl-Differenz zwischen ADR und Plan unkommentiert"

### F-7 — Das Falsch-Beispiel in §3.11 ist ein Platzhalter und in dieser Form kein Befund

- `kategorie`: INFO
- `quelle`: Plan §2, Liefer-Punkt 3 („Falsch/Richtig an **realen** Beispielen")
- `pfad`: `AGENTS.md:322-326`
- `befund`: Das Falsch-Beispiel nennt `<Host-Wurzel>` statt des Wurzel-Segments
  und steht im Fence (den §3.11 selbst freigibt). Als Text würde das
  `hostpaths`-Modul daran nichts melden; das Richtig-Beispiel ist dagegen die
  reale korrigierte Zitation. Die Regel „auch das Beispiel deckt" erklärt den
  Platzhalter, die DoD-Formulierung „reale Beispiele" trifft auf die
  Falsch-Hälfte damit nur der Form nach zu.
- `verifizierbar`: nein
- `klasse`: „Falsch-Beispiel ohne realen Befund"

## Negativbefunde

- geprüft, ohne Befund: `.d-check.yml` — `hostpaths` steht in `modules`; kein
  `hostpaths:`-Knoten, kein `scope`, kein `ignore`, kein `exempt-paths` für
  das Modul; keine Mutation des Ausschlussblocks gegenüber dem Vorzustand
- geprüft, ohne Befund: `Makefile:77` und `tools/coverage-gate.sh:4` — die
  beiden einzigen Nicht-Markdown-Fundstellen sind korrigiert; weitere
  Nicht-Markdown-Dateien trugen vor dem Slice keinen Host-Pfad
- geprüft, ohne Befund: `docs/plan/adr/0051-…md` — genau ein Hunk (Zitatgerüst)
  plus §Geschichte-Zeile; keine Stelle aus §Entscheidung, §Konsequenzen,
  §Verglichene Alternativen, §Status oder `Supersedes`-Kette berührt
- geprüft, ohne Befund: `docs/plan/adr/0054-…md` — acht Hunk-Zeilen, alle
  Zitatgerüst (Host-Präfix → Hausform) plus §Geschichte-Zeile; der
  Entscheidungstext ist unberührt (die Locator-Frage steht als F-2 gesondert)
- geprüft, ohne Befund: `docs/plan/planning/done/slice-049`, `-050`,
  `welle-14`, `welle-14-results` sowie `docs/reviews/review-slice-036/039/049/073`
  — nur Zitatgerüst; keine Fund-, Beobachtungs- oder Closure-Aussage geändert
- geprüft, ohne Befund: `docs/plan/adr/README.md`, `harness/README.md`,
  `harness/sensors/coverage-gate.md` — die Kopf-Satz-Ausnahme und die
  `hostpaths`-Nennung sind Ist-Zustand, ohne Chronik
- geprüft, ohne Befund: Fenced-, Nicht-Markdown- und `.harness/`-Formen des
  Host-Pfads — repo-weit keine Fundstelle (die benannten Ränder tragen)

## Eigene Messungen

Alle Läufe netzlos gegen das in `d-check.mk` gepinnte Image
`sha256:18e9cd85…`; Exit-Codes ungepiped (`AGENTS.md` §3.9).

| Messung | Ergebnis | Exit |
|---|---|---|
| `hostpaths`-Modul über das Repo (`--enable hostpaths --disable <übrige>`) | 621 Dateien geprüft, **0** Befund(e) | 0 |
| repo-weiter Grep über die Präfix-Muster (`.md` + Nicht-`.md`, inkl. Fences, `.harness/`, Binär via `-a`) | **0** Fundstellen in 15 ehemaligen Dateien | 1 (grep: kein Treffer) |
| dieselben Muster **am Basisstand** `3f4c747` (`git grep -n`) | **42** Zeilen in **15** Dateien | 0 |
| Mutation A — Host-Pfad in einer lebenden `.md` im Repo-Wurzel | 622 Dateien, **1** Befund `hostpath-forbidden` | 1 |
| Mutation B — Host-Pfad **im Fence** + relativer Pfad | 622 Dateien, **0** Befund(e) | 0 |
| Mutationen exakt zurückgenommen, `git status --porcelain` | leer | — |
| `make docs-check` (aktives Modul, ohne Ausschlussblock) | 621 Dateien, 0 Befunde | 0 |
| `make gates` (voller Satz, Nachweis-Stempel zuletzt) | Coverage 49,30 % ≥ 40 %, a-check 0, traceability OK | 0 |

Der Basisstand belegt zugleich die Zählung: `0051` 1 · `0054` 10 · `0072` 4 ·
Verdikt 5 · `done/` 11 · `docs/reviews/` 7 · `coverage-gate.md` 2 · `Makefile`
1 · `tools/coverage-gate.sh` 1 = **42**. Die zwei Mutationen bestätigen genau
die in §3.11/`docs-check.md` benannten Ränder: lebende Prosa wird gemeldet,
Fence und relativer Pfad nicht.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** „Zitat-Korrektur überschreitet die
Klassengrenze (§Entscheidung)" · „Zeilen-Lokator nicht durch benannten Anker
ersetzt" · „Regel-Reichweite für Fenced-Blöcke widersprüchlich" ·
„Nachgezogener Hard-Rule-Wortlaut weicht von der ADR-Vorgabe ab" ·
„Modulliste doppelt geführt" · „Zähl-Differenz zwischen ADR und Plan
unkommentiert" · „Falsch-Beispiel ohne realen Befund"

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) und F-2/F-3 (MEDIUM) blockieren.

**Übergabe:** F-1 ist eine Abweichung von `AGENTS.md` §3.5/`ADR-0073`; ob der
Kanal benigne war, ist ein **Architect-Urteil** (Verdikt-Pfad, Modul 8) — bei
„benigne" ist `ADR-0074` um die Lokator-/Beispiel-Klausel zu ergänzen, sonst
ist die Änderung aus `ADR-0072` `§Entscheidung` zurückzunehmen. F-2 hängt an
derselben Frage: der Lokator-Ersatz ist nach `ADR-0073` in-place zulässig;
bleibt die Klausel, zieht der Implementer nach. F-3 liegt als Plan-Defekt beim
Planner (Rückkante Review → Plan). F-4 ist ein Implementer-Nachzug am
`§3.5`-Satz. **Kein DoD-Nachzug:** der Slice braucht eine Fixrunde, die
DoD-Zeile „Review durchgeführt …" bleibt daher offen (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde).

Der Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft
der Verifier separat (Modul 11).
