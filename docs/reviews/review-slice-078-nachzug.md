# Review-Report: slice-078 — Nachzug zu N-1 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** der Ein-Satz-Nachzug `859357b` (N-1 aus
[`review-slice-078-delta`](review-slice-078-delta.md)). Frisches Artefakt:
geprüft wird der Commit-Diff, nicht die Zusage, dass er trägt.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd`.
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext:**

- `git diff 900c6c6..859357b` (vollständig, alle drei Dateien)
- [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md)
  §Entscheidung Punkt 1–2 · [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung Punkt 4 (Records) · `AGENTS.md` §3.7, §3.11 (Fassung 2)
- `harness/sensors/docs-check.md` §Vertrag, §Grenze Punkt 8 · `.d-check.yml`

---

## Antworten

**1 — Der Selbstwiderspruch ist weg.** `AGENTS.md` schließt jetzt mit „Dieser
Abschnitt trägt die **Regel und ihre Reichweite**; die Präfixliste führt allein
das Modul …, und was der Sensor deckt und was nicht, führt `docs-check.md` aus
seiner Sicht" (`AGENTS.md:360-363`); die Sensordoku spiegelt umgekehrt („Dieser
Abschnitt trägt, **was der Sensor deckt und was nicht**; die Reichweite der
Regel steht in `AGENTS.md` §3.11", `docs-check.md:113-115`). Kein Träger
behauptet mehr, der andere wiederhole die Grenzen nicht. Die verbleibenden
„stehen in"-Sätze zeigen auf **verschiedene Gegenstände** (`.d-check.yml` =
Vertrags-Deklaration, `§3.11` = Regel-Reichweite, Sensordoku = Sensor-Deckung)
und widerlegen einander nicht. Die Begründung für Weg B trägt: die
Perspektiven-Zuweisung ist umformulierungsstabil, eine
Nicht-Wiederholungs-Behauptung wäre es nicht.

**2 — Die Grenzen stehen weiterhin vollständig an beiden Orten.**
`AGENTS.md:340-347` (Fenced-Code-Blöcke, relative Pfade, Nicht-Markdown —
`Makefile`, `tools/**`, `harness/mk/**` —, `scan.ignore`, plus die
Fenced-Regel-Lücke) und `docs-check.md:102-114` (vier Ränder plus
`scan.ignore`, plus dieselbe Lücke). `ADR-0075` Punkt 2 verlangt den
Sensor-Block in `§3.11` — er steht. Beim Streichen ist nichts Substantives
mitverschwunden (der Diff ersetzt genau die Schluss-Klausel).

**3 — Keine dritte Formulierung.** Die Grenz-Liste lebt an den drei
vorgesehenen Orten: der Entscheidung (`ADR-0075` Punkt 2), der Hard Rule
`§3.11` und der Sensordoku `docs-check.md` — Entscheidung plus zwei Träger,
keine weitere. `ADR-0072`s Entwurf ist supersedet.

**4 — Abnahme unverändert.** Modul **0**, Grep **0**; der Diff zeigt **drei**
Dateien, nicht zwei (N-6).

---

## Findings

### N-1 — geschlossen

- `kategorie`: INFO (Closure des LOW aus dem Delta-Lauf)
- `befund`: Der Selbstwiderspruch ist beseitigt; die Grenzen sind an beiden
  Orten vollständig. Weg B (Behauptung ersatzlos gestrichen) trägt.

### N-6 — Der Nachzug-Commit ändert zusätzlich den Delta-Report

- `kategorie`: INFO
- `quelle`: [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung Punkt 4 (Records: nur das Zitat-Gerüst) · `AGENTS.md` §3.5
- `pfad`: `docs/reviews/review-slice-078-delta.md:121` (Commit `859357b`)
- `befund`: `859357b` berührt drei Dateien: `AGENTS.md`, `docs-check.md` und
  den Delta-Report — dort wird eine nackte Kennung im `befund` von N-4 auf
  einen Link gezogen. Nachgemessen: die nackte Form ergab
  `id-unlinked` (`make docs-check` Exit 2, `review-slice-078-delta.md:121`);
  die Link-Form ergibt 0. Der Eingriff ist damit **notwendig** und eine
  Zitat-Korrektur am Record-Gerüst (Referent unverändert), die `ADR-0073`
  Punkt 4 zulässt; der Nachzug-Commit ist der Beleg. Die Abweichung von der
  erwarteten Zwei-Datei-Stat ist **mein** Defekt aus dem Delta-Report, nicht
  einer des Nachzugs.
- `verifizierbar`: ja — `make docs-check` (`ids`)
- `klasse`: „Record-Zitatgerüst nachgezogen (nackte Kennung → Link)"

### N-7 — Die Sensordoku nennt den `§3.11`-Zeiger zweimal

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `harness/sensors/docs-check.md:112-115`
- `befund`: Innerhalb von vier Zeilen steht zweimal, die Reichweite der Regel
  liege in `AGENTS.md` §3.11 („ihre Reichweite und diese benannte Lücke stehen
  in `AGENTS.md` §3.11" und „die Reichweite der Regel steht in `AGENTS.md`
  §3.11"). Rein redaktionelle Doppelung ohne semantische Wirkung.
- `verifizierbar`: nein
- `klasse`: „Doppelter Querverweis ohne Aussage-Zuwachs"

## Eigene Messungen

| Messung | Ergebnis | Exit |
|---|---|---|
| `hostpaths`-Modul über das Repo | 625 Dateien, **0** Befund(e) | 0 |
| repo-weiter Grep über die Präfix-Muster | **0** Fundstellen | 1 |
| `git diff 900c6c6..859357b --stat` | **3** Dateien (`AGENTS.md`, `docs-check.md`, Delta-Report) | — |
| Rücksetzung der Link-Form im Delta-Report (Experiment) | `id-unlinked` auf `:121` | 2 |
| Experiment zurückgenommen (`git checkout`), `git status --porcelain` | leer | — |
| `make gates` | baseline-verify OK, coverage 49,30 % ≥ 40 %, a-check 0, traceability OK | 0 |

## Negativbefunde

- geprüft, ohne Befund: `AGENTS.md` §3.11 Fassung 2 nach `859357b` — Reichweite
  Fences eingeschlossen, Platzhalter-Form, benannte Lücke, Präfixliste nur als
  Zeiger; kein Widerspruch zur Sensordoku mehr
- geprüft, ohne Befund: `docs-check.md` Punkt 8 — vier Ränder plus
  `scan.ignore` vollständig, Fenced-Regel-Lücke benannt, Perspektive markiert
- geprüft, ohne Befund: keine dritte normative Grenz-Formulierung im Repo
- geprüft, ohne Befund: `859357b` fasst **keine** Accepted-ADR an

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** „Record-Zitatgerüst nachgezogen (nackte
Kennung → Link)" · „Doppelter Querverweis ohne Aussage-Zuwachs"

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW; N-1 ist geschlossen.
Die N-6-Stat-Abweichung ist ein nachgemessener, `ADR-0073`-konformer Nachzug
an einem Record und **mein** Delta-Defekt; N-7 ist redaktionell.

**DoD-Nachzug:** ja — die Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" in §2 des Slice-Plans ist auf `[x]` nachgezogen,
mit Verweis auf Erstlauf, Delta-Lauf und diesen Nachzug
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Nur
diese eine Zeile.

Der Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation (Modul 11).
