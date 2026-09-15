# Review-Report: slice-076 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-076`, Diff-Range `7670da4..HEAD` — die beiden
slice-076-Commits `542f4ff` (die Änderung) und `9f98b7a` (DoD-Haken +
Plan-Nachzug). `7670da4` ist der reine `next → in-progress`-Move und
ausdrücklich **nicht** Teil des Gegenstands. Die Range enthält zusätzlich zwei
**fremde** Commits (`0ef7f13`, `6f1d01d` — Planung von `slice-077`); sie sind
nicht Gegenstand dieses Reports. `git diff --name-only 7670da4..HEAD` nennt
fünf Pfade, davon vier zu `slice-076`: `harness/mk/coverage.mk` ·
`harness/sensors/coverage-gate.md` · `harness/README.md` · der Slice-Plan.
`AGENTS.md` ist **nicht** im Diff (der Implementer begründet den Verzicht; s.
Negativbefund-Zeile).

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09 — vier
repo-spezifische HIGH-Regeln; auf diesen Diff kommt keine davon zur Anwendung).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-076-coverage-gate-reifestufe-40.md`
  vollständig (§1–§8), sowohl im Stand von `7670da4` als auch nach dem
  §3/§6-Nachzug aus `9f98b7a`
- [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
  vollständig — §(a) *Eskalationsklausel*, §Re-Evaluierungs-Trigger (b) und
  (c), §Fitness Function, §Konsequenzen
- `architect-verdict-coverage-gate-reifestufe` — insbesondere §Verdikt 2/3/5,
  §Frage 2 (Option C gegen D), §Frage 3 (Nachweisform), §Was `slice-076`
  ändern muss
- `AGENTS.md` §3.6 (Schwellen-Senkung nur per ADR), §3.7
  (Kommentar/Zustandsfeld nennt den Ist-Zustand), §3.9 (Exit-Code nie durch
  eine Pipe), §4 (Quality Gates), §5 (Traceability)
- `harness/conventions.md` (MR-000 ID-Schema) ·
  `harness/mk/coverage.mk` · `harness/sensors/coverage-gate.md` ·
  `tools/coverage-gate.sh` · `harness/sensors/docs-check.md`
- `v6.5.0` · `regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad als
  Rollen-Sequenz · `regelwerk/modul-10-review-harness.md`
- `verify-slice-049.md` §2 (dokumentierte Lauf-zu-Lauf-Schwankung) ·
  `welle-14-results.md` §Verifikation (Präzedenz-Lesart desselben Triggers)

---

## Findings

### F-1 — Einzelquellen-Aussage der Sensor-Doku durch ihr eigenes Beleg-Zitat widerlegt

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 (ein geltender Wert, ein Ort)
- `pfad`: `harness/sensors/coverage-gate.md:25-28` gegen `harness/sensors/coverage-gate.md:68-73`
- `befund`: Der neue Absatz sagt, der bewegliche Wert stehe „ausschließlich in
  `harness/mk/coverage.mk`“ und diese Sektion führe ihn nicht — der Beleg
  derselben Sektion nennt drei Zeilen darunter `THRESHOLD=40` und „Schwelle
  40%“ als real bestehende Stufe. Beide Aussagen stehen in derselben Datei;
  die zweite widerlegt die erste.
- `verifizierbar`: nein — kein Gate deckt einen Doku-Selbstwiderspruch
  (eigener Mutationslauf, s. F-4: `make docs-check` bleibt grün, auch wenn die
  Doku einen abweichenden Wert führt).
- `klasse`: Einzelquellen-Aussage durch Beleg-Zitat in derselben Sektion widerlegt

### F-2 — Risiko-Herleitung mit veralteter Messgröße

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 (Zustandsfeld nennt den Ist-Zustand)
- `pfad`: `docs/plan/planning/in-progress/slice-076-coverage-gate-reifestufe-40.md:177-182`
- `befund`: §6 Risiko 3 begründet die Wahl `THRESHOLD=50` damit, dass „bei
  `THRESHOLD=46` der Abstand zum Ist-Stand 0,2 Prozentpunkte“ beträgt. §6
  Risiko 1 und der reale Beleg führen den Ist-Stand mit 49,30 % — der Abstand
  zu 46 wäre 3,3 Prozentpunkte (der Lauf wäre dort grün, nicht „könnte grün
  ausfallen“), und der tatsächlich gewählte Rot-Beleg liegt 0,7 Prozentpunkte
  über der Messung. Die Zahl stammt aus dem Stand des Architect-Verdikts
  (45,80 %); Risiko 1 wurde nachgezogen, Risiko 3 nicht.
- `verifizierbar`: ja — `make coverage-gate` nennt den Ist-Stand (49,30 %);
  `49,30 − 46 = 3,3`.
- `klasse`: Risiko-Herleitung mit veralteter Messgröße

### F-3 — Trigger-Bedingung mit unauflösbarer Bezugszahl

- `kategorie`: LOW
- `quelle`: Maintainability · `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Trigger je Lifecycle-Übergang und WIP-Limit
- `pfad`: `docs/plan/planning/in-progress/slice-076-coverage-gate-reifestufe-40.md:137-140` gegen `:109-117` und `:169-176`
- `befund`: §4 benennt die Rückführungs-Bedingung als „an **mehr** Orten
  geführt … als den drei geplanten“, während §3 vier Datei-Zeilen führt und §6
  vier Bindung-Orte nennt; nach dem Nachzug führt genau **einer** davon die
  geltende Stufe. Gegen welchen Zählstand die Bedingung zu messen ist, ist dem
  Plan damit nicht mehr zu entnehmen.
- `verifizierbar`: nein (Plan-interne Zählung; kein Gate)
- `klasse`: Trigger-Bedingung mit unauflösbarer Bezugszahl

### F-4 — Benannte Lücke bestätigt: kein Gate gegen Wert-Dopplung

- `kategorie`: INFO
- `quelle`: Maintainability · `harness/sensors/docs-check.md` (Geltungsbereich)
- `pfad`: `harness/mk/coverage.mk:8-14` (Behauptung) · `harness/README.md:117` (Mutationsziel)
- `befund`: Ein absichtlich divergierender Wert in `harness/README.md`
  („geltende Stufe 35 %“ bei `THRESHOLD ?= 40`) lässt `make docs-check` mit
  Exit 0 und 0 Befunden durch — die vom Implementer benannte Lücke ist real und
  nicht nur behauptet. Bewertung: **akzeptable, benannte Grenze**. Der Slice
  hat den beweglichen Wert auf einen Träger zusammengezogen; ein Sensor hätte
  heute kein Objekt (es gibt keine zweite Stelle, die ihn führt) und wäre als
  Textmuster-Sensor dieselbe Werkzeugklasse, die der Skill für die
  Chronik-Klasse geprüft und verworfen hat. Als eigene Beobachtung nur dann
  ein Registerkandidat, wenn die Dopplung ein zweites Mal auftritt.
- `verifizierbar`: ja — Mutation + `make docs-check` (Exit 0, 0 Befunde).
- `klasse`: Wert-Dopplung ohne Sensor (benannt)

### F-5 — Verweis auf einen Träger, der den Wert nicht mehr führt

- `kategorie`: INFO
- `quelle`: Maintainability · `AGENTS.md` §3.7
- `pfad`: `tools/coverage-gate.sh:2-3`
- `befund`: Der Skriptkopf verweist für „aktuelle Schwelle und Historie“ auf
  `harness/README.md` §Sensors; seit diesem Diff führt README §Sensors den
  beweglichen Wert nicht mehr, sondern verweist auf `harness/mk/coverage.mk`
  (zweiter Sprung). Die Datei ist nicht Teil des Diff, und README trug auch
  vorher nur die Rampe — kein Nachzug verlangt, nur benannt.
- `verifizierbar`: nein
- `klasse`: Zeiger-Kette um ein Glied verlängert

### F-6 — Rollen-Zeiger: DoD-Bewertung der Ein-Ort-Zeile

- `kategorie`: INFO
- `quelle`: `v6.5.0` · `regelwerk/modul-08-agentenrollen.md` §Die neun Übergaben und ihre Artefakte
- `pfad`: `docs/plan/planning/in-progress/slice-076-coverage-gate-reifestufe-40.md:75-78`
- `befund`: Ob die DoD-Zeile „die geltende Stufe steht widerspruchsfrei an
  **einem** beweglichen Ort und wird an den übrigen Stellen nur als Rampe bzw.
  als Verweis genannt“ angesichts von F-1 erfüllt ist, ist eine
  DoD-/Spec-Frage und liegt beim **Verifier** (Modul 11) — hier nur benannt,
  nicht bewertet.
- `verifizierbar`: nein
- `klasse`: Rollen-Zeiger (DoD-Bewertung nicht Gegenstand des Reviews)

## Negativbefunde

- geprüft, ohne Befund: `harness/mk/coverage.mk` — `THRESHOLD ?= 40`; neuer
  Kopf-Kommentar trägt die Kommentar-Klassen Zusage („die geltende Stufe ist
  `THRESHOLD` unten“), Kopplung (Rampe, Endstufe, Trigger) und Grenze
  („Senkung … nur per ADR“); die Einstiegs-Herleitung („39.6 %“) ist entfernt,
  keine Stufen-Chronik (§3.7). Die `## help`-Zeile nennt die Rampe, keinen
  beweglichen Wert.
- geprüft, ohne Befund: `harness/README.md` §Sensors-Zeile `make coverage-gate`
  — nennt die Rampe und verweist auf den einen beweglichen Ort; kein
  Einzelwert.
- geprüft, ohne Befund: `AGENTS.md` (nicht im Diff) — §4 Zeile 311 nennt
  `Einstiegsstufe 35 % → Endstufe 80 %` und **keinen** beweglichen Wert; §3.6
  (Z. 93–99) ebenso. Der Verzicht auf den Edit ist grep-verifiziert korrekt,
  kein unterlassener Nachzug (eigener `grep` über `*.md`/`*.mk`/`*.sh`).
- geprüft, ohne Befund: `docs/plan/adr/**` — `ADR-0054` ist im Range
  **unverändert**, bleibt `Accepted`, kein `Supersedes`; die
  Fitness-Function-Zeile trägt die Bewegung mit „aktuell gültiger Schwelle“
  ohne Textänderung. Genau **eine** Stufe (35 → 40); keine zweite Hochschaltung
  (45 %) mitgenommen — repo-weiter `grep` nach `40 %`/`45 %`: kein weiterer
  Träger außer Plan und Review-Belegen.
- geprüft, ohne Befund: Commit-Range — beide Betreffe nennen `ADR-0054`, keine
  `SPEC-*`/`ARC-*` im Betreff, keine Attributions-Trailer, `git commit -F`.
- geprüft, ohne Befund: `harness/sensors/coverage-gate.md` §Ausgabe und
  Ausgänge — die Exit-Tabelle ist unverändert und meint die Skript-Exits; die
  Beleg-Zeile nennt die Träger getrennt („das Gate-Skript endet Exit 1, `make`
  meldet … Exit 2“) und ist damit präziser als der Stand davor.
- geprüft, ohne Befund: Rot-Beleg — `make coverage-gate THRESHOLD=50` endet
  real rot; `docker build` meldet für den RUN-Schritt `exit code: 1` (der
  Skript-Exit, durch `pipefail` durchgereicht), `make` selbst Exit 2. Die
  Nachzug-Formulierung „Exit ≠ 0 — Skript 1, `make` 2“ trifft den Ist-Stand.
- nicht Gegenstand: `0ef7f13` / `6f1d01d` (`slice-077`-Planung) im Range —
  fremder Slice, fremder Review.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Einzelquellen-Aussage durch Beleg-Zitat in
derselben Sektion widerlegt · Risiko-Herleitung mit veralteter Messgröße ·
Trigger-Bedingung mit unauflösbarer Bezugszahl · Wert-Dopplung ohne Sensor
(benannt) · Zeiger-Kette um ein Glied verlängert · Rollen-Zeiger
(DoD-Bewertung)

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Die drei LOW sind
Plan- bzw. Doku-Wortlaut-Findings ohne Wirkung auf Gate oder Produkt-Code.

**Übergabe:** Die drei LOW nehmen die **Rückkante Review → Plan** (Modul 8,
Plan-Defekt): F-3 und F-2 liegen in §4/§6 des Slice-Plans, F-1 in einem Satz
der Sensor-Doku, dessen absolute Formulierung aus der §3-Zeile
„die geltende Stufe wird dort nur als Verweis geführt“ neben „die Beleg-Zeile
auf die neu gemessenen Werte nachziehen“ folgt — die Spannung steht damit
bereits im Plan. **Kein Reviewer → Implementer-Pfeil**: kein Fix an Code, Gate
oder Sensor-Semantik ist nötig; der Wortlaut-Nachzug kann im
Closure-Zug des Planners mitlaufen. Dieser Report ist ein **Lauf-Beleg** und
ersetzt keine Verifikation (Modul 11); der Rollen-Zeiger F-6 nennt die eine
Stelle, an der der Verifier zuständig ist.

**DoD-Nachzug (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne
Fixrunde):** Die Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt
vor“ ist im Slice-Plan auf `[x]` nachgezogen, im selben Commit wie dieser
Report.
