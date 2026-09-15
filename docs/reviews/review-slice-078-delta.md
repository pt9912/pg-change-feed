# Review-Report: slice-078 — Delta-Nachlauf (Fixrunde) — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** Fixrunde `9c33855..f0f2a8a` (Architect-Zug `9c33855`; Implementer
`f0f2a8a`). Neues Artefakt, nicht Bestätigung: geprüft werden
`ADR-0075` (`docs/plan/adr/0075-hostpaths-reichweite-und-wortlaut.md`), das
Verdikt `docs/reviews/architect-verdict-slice-078-konfliktpfad.md` und der
Implementer-Diff — **frisch gegen den Plan**, nicht gegen die Findings des
ersten Reports (`review-slice-078.md`, Commit `b651441`).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd`.
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext:**

- `ADR-0075` vollständig (§Kontext, §Entscheidung Punkt 1–3, §Verglichene
  Alternativen, §Konsequenzen, §Fitness Function, §Re-Evaluierungs-Trigger,
  §Geschichte)
- `docs/reviews/architect-verdict-slice-078-konfliktpfad.md` vollständig
- [`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md),
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
  [`ADR-0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md)
- `AGENTS.md` §3.5, §3.11 (Fassung 2) · `harness/sensors/docs-check.md`
  (§Vertrag, §Grenze Punkt 8, §Bindung) · `harness/README.md` §Sensors
- Slice-Plan `docs/plan/planning/in-progress/slice-078-hostpfade-entfernen.md`
- `.d-check.yml`, `d-check.mk`, gepinntes Modul-Image `pt9912/d-check` `v0.75.0`

---

## Antworten auf die fünf Prüfpunkte

**1 — Legitimierung (F-1).** `ADR-0075` Punkt 2 deckt **genau** den Eingriff:
es benennt die drei geänderten Elemente beim Namen („Platzhalter-Beispiel,
Hausform-Beispiel, Form-Beschreibung im Klammertext") und stellt sie als
**beschlossenen** Text fest; die §Status-Zeile supersedet `ADR-0072`
§Entscheidung Punkt 5. Die **Commit-Attribution ist korrigiert und
nachgeprüft**: `git show df47282` enthält den §Entscheidung-Hunk
(`0072:173-190`) **und** den §Kontext-Hunk; `git show aab7aa8` berührt nur
`### Der gemessene Befund` (§Kontext). Der Verdikt-Eintrag trägt also — meine
erste Attribution (`aab7aa8`) war falsch. **Dritte Reichweiten-Aussage:**
siehe N-1/N-5; die einzige *nicht* identische steht in `ADR-0072`s supersedetem
§3.11-Entwurf und ist über `ADR-0075` §Status plus Index-Zeile 87
(`→ ADR-0074/0075, teilw.`) deklariert.

**2 — F-3/F-7.** Die Lücke („die Regel deckt die Fenced-Fläche voll, der Sensor
nicht") steht an **beiden** Orten, an denen sie wirkt: `AGENTS.md:345-346` und
`docs-check.md:111-114`. Der Platzhalter `<Host-Wurzel>` ist grep- und
modul-sicher (beide Läufe 0); Mutation A/B belegen das (siehe Messungen).

**3 — F-4/F-5.** `AGENTS.md` §3.5 ist **sinngleich, nicht wortgleich** zur
Vorgabe aus `ADR-0073` Punkt 2 (siehe N-3). Die Modulliste ist **nicht** auf
eine Stelle reduziert (N-2).

**4 — Abnahme.** Modul **0** (Exit 0) und Grep **0** (Exit 1); beide Mutationen
laufen wie vorhergesagt (siehe §Eigene Messungen).

**5 — Grenze.** Der Implementer-Diff `9c33855..f0f2a8a` berührt **nur**
`AGENTS.md` und `harness/sensors/docs-check.md` — beide vom Verdikt
freigegeben; **keine** `ADR-0051/0054/0072/0073/0074/0075` angefasst, **keine**
`§Entscheidung` eines Accepted-Dokuments geändert, `.d-check.yml` unverändert.
Kein Übergriff.

---

## Findings

### N-1 — `§3.11` und `docs-check.md` behaupten je, der andere nenne die Modul-Grenzen nicht erneut — beide tun es

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Zwei-Quellen-Disziplin) · `ADR-0075`
  §Entscheidung Punkt 2 (Fassung 2 führt diesen Satz **nicht**)
- `pfad`: `AGENTS.md:360-363` · `harness/sensors/docs-check.md:114`
- `befund`: `§3.11` zählt die Sensor-Grenzen auf (Fenced-Code-Blöcke, relative
  Pfade, Nicht-Markdown, `scan.ignore`, `:342-345`) und schließt mit „Die
  Modul-Grenzen stehen einmal in `harness/sensors/docs-check.md` — dieser
  Abschnitt nennt sie nicht erneut"; `docs-check.md` Punkt 8 zählt dieselben
  Grenzen auf und schließt mit „Die Modul-Grenzen stehen einmal hier — `§3.11`
  nennt sie nicht erneut". Beide Nicht-Wiederholungs-Aussagen sind durch die
  jeweils andere Datei widerlegt; nur die Modul-*Semantik* (Präfixliste,
  Windows-Muster, Zeiger-Register) steht tatsächlich einmal.
- `verifizierbar`: nein
- `klasse`: „§3.7-Nicht-Wiederholungs-Aussage in beiden Trägern falsch"

### N-2 — Die Modulliste steht weiter in zwei Dateien

- `kategorie`: INFO
- `quelle`: Maintainability (`AGENTS.md` §3.7)
- `pfad`: `harness/README.md:114` · `.d-check.yml:` (`modules`)
- `befund`: `docs-check.md` hat die Ist-Zustand-Modulliste korrekt abgegeben;
  die Liste steht aber nach wie vor in `harness/README.md` §Sensors
  („links, anchors, ids, matrix, versions, structure, hostpaths"). Die Zeile
  ist plan-gefordert (DoD Liefer-Punkt 2), also kein Defekt der Fixrunde —
  aber „nur noch an einer Stelle" trifft nicht zu.
- `verifizierbar`: nein
- `klasse`: „Modulliste in Deklaration und Übersicht doppelt"

### N-3 — `AGENTS.md` §3.5 ist nicht wortgleich zu `ADR-0073` Punkt 2

- `kategorie`: INFO
- `quelle`: [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung Punkt 2
- `pfad`: `AGENTS.md:144-148`
- `befund`: Vorgegeben war die Parenthese „host-lokale Pfade, Linkziele,
  Zeilen-Lokatoren"; ausgeliefert ist „host-lokale Pfade, Linkziele,
  Zeilen-Lokatoren, die Form einer gebrochenen Referenz". Das vierte Element
  ist `ADR-0073` Punkt 1s vierte Klasse — die Fassung ist damit **die volle
  Klassen-Enumeration aus Punkt 1**, keine Ausweitung über die ADR hinaus
  (F-4 des ersten Laufs ist damit behoben). Sinngleich, nicht wortgleich.
- `verifizierbar`: nein
- `klasse`: „Nachgezogene Hard Rule spiegelt Punkt 1 statt Punkt 2"

### N-4 — Die Index-Zeile von `ADR-0075` nennt den Teil-Supersede nicht

- `kategorie`: INFO
- `quelle`: Maintainability (ADR-Index-Konvention, `docs/plan/adr/README.md` Kopf)
- `pfad`: `docs/plan/adr/README.md:90` (vs. `:89`)
- `befund`: `ADR-0074`s Zeile trägt „(Supersedes ADR-0072, teilweise)";
  `ADR-0075`s Zeile trägt nur den Titel. Die Lineage ist über Zeile 87
  (`ADR-0072 … (→ ADR-0074/0075, teilw.)`) und `ADR-0075` §Status dennoch
  auffindbar; die `structure`-Regel (`cell-max-chars: 80` für `Titel`) begrenzt
  die Zelle.
- `verifizierbar`: ja — `make docs-check` (Modul `structure`)
- `klasse`: „Supersede-Vermerk fehlt in der Index-Zeile der ablösenden ADR"

### N-5 — Der DoD-Punkt §2 nennt weiter „Fenced-Blöcke frei"

- `kategorie`: INFO
- `quelle`: Plan §2, Liefer-Punkt 2 · [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md)
  §Konsequenzen (Slice-Plan trägt keine Regel)
- `pfad`: `docs/plan/planning/in-progress/slice-078-hostpfade-entfernen.md:131`
- `befund`: Der DoD-Punkt charakterisiert die ausgelieferte Sensordoku als
  „Fenced-Blöcke frei"; gemeint ist die Sensor-Grenze — der ausgelieferte Text
  sagt zusätzlich, dass die **Regel** die Fenced-Fläche voll deckt. Als
  Modul-Grenze richtig, als Regel-Aussage überholt; der Plan ist nach
  `ADR-0075` Zeitdokument und trägt keine Regel-Fassung.
- `verifizierbar`: nein
- `klasse`: „DoD-Charakterisierung nach Reichweiten-Festlegung nicht nachgezogen"

---

## Eigene Messungen

| Messung | Ergebnis | Exit |
|---|---|---|
| `hostpaths`-Modul über das Repo | 624 Dateien, **0** Befund(e) | 0 |
| repo-weiter Grep über die Präfix-Muster | **0** Fundstellen | 1 |
| Mutation A — Host-Pfad in **Prosa** | 625 Dateien, **1** Befund | 1 |
| Mutation A — Grep | findet die Stelle | 0 |
| Mutation B — Host-Pfad im **Fence** | 625 Dateien, **0** Befund(e) | 0 |
| Mutation B — Grep | findet die Stelle (Regel greift, Sensor nicht) | 0 |
| Mutationen zurückgenommen, `git status --porcelain` | leer | — |
| `make gates` | Coverage 49,30 % ≥ 40 %, baseline-verify OK, a-check 0, traceability OK | 0 |

Die zwei Läufe `git show df47282` / `git show aab7aa8` belegen die
Attributions-Korrektur des Verdikts (df47282 trug den §Entscheidung-Hunk).

## Closure der Findings aus dem ersten Lauf

| Erstlauf | Stand |
|---|---|
| F-1 HIGH | **geschlossen** — `ADR-0075` Punkt 2 (Folge-ADR, Verdikt „Lockerung legitim, aber undokumentiert"); Attribution auf `df47282` korrigiert |
| F-2 MEDIUM | **geschlossen** — `ADR-0075` Punkt 3: Lokator-Klausel supersedet-nicht-restitiert, keine Ersatzpflicht |
| F-3 MEDIUM | **geschlossen** — `ADR-0075` Punkt 1: Gewinner „Fences eingeschlossen", an `§3.11` festgezogen |
| F-4 LOW | **geschlossen** — `§3.5` auf `ADR-0073` Punkt 1s Klassen-Liste gezogen (N-3: sinngleich, nicht wortgleich) |
| F-5 INFO | **teilweise geschlossen** — Modulliste aus `docs-check.md` raus; `harness/README.md:114` bleibt (N-2) |
| F-6 INFO | **Adresse benannt** — Planner-Closure (`ADR-0075` §Kontext führt 31↔42) |
| F-7 INFO | **geschlossen** — Platzhalter ist beschlossene Form (`ADR-0075` Punkt 2) |

## Negativbefunde

- geprüft, ohne Befund: `ADR-0075` §Entscheidung Punkt 1–3 — deckt genau den
  Text, den `df47282` in `ADR-0072` §Entscheidung Punkt 5 gesetzt hat; keine
  weitere In-place-Änderung an einer Accepted-ADR in `9c33855`
- geprüft, ohne Befund: Diff `9c33855..f0f2a8a` — nur `AGENTS.md` und
  `harness/sensors/docs-check.md`; keine Accepted-ADR, `.d-check.yml`
  unverändert, kein Ausschlussblock
- geprüft, ohne Befund: `AGENTS.md` §3.11 Fassung 2 — Reichweite Fences
  eingeschlossen, Platzhalter-Form, benannte Lücke; Präfixliste nur als
  Zeiger (`hostpaths.prefixes`)
- geprüft, ohne Befund: die Lücke steht an beiden Wirkorten (`§3.11`,
  `docs-check.md` Punkt 8); der Platzhalter ist grep- und modul-sicher
- geprüft, ohne Befund: `harness/sensors/docs-check.md` §Bindung — Zeiger auf
  `ADR-0075`/`ADR-0072`/`ADR-0073`/`ADR-0074` Ist-Zustand, ohne Chronik
- geprüft, ohne Befund: `make gates` mit aktivem Modul und ohne
  Ausschlussblock

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** „§3.7-Nicht-Wiederholungs-Aussage in beiden
Trägern falsch" · „Modulliste in Deklaration und Übersicht doppelt" ·
„Nachgezogene Hard Rule spiegelt Punkt 1 statt Punkt 2" · „Supersede-Vermerk
fehlt in der Index-Zeile der ablösenden ADR" · „DoD-Charakterisierung nach
Reichweiten-Festlegung nicht nachgezogen"

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM; die drei Regelfragen des
Erstlaufs sind durch `ADR-0075` entschieden.

**Übergabe:** N-1 ist ein Ein-Satz-Nachzug am Implementer (Wortlaut in
`AGENTS.md` §3.11 und `harness/sensors/docs-check.md`), N-2…N-5 sind Hinweise
ohne Aktion (die Index-Zeile N-4 ist durch die 80-Zeichen-Zelle der
`structure`-Regel begrenzt).

**Kein DoD-Nachzug.** N-1 trägt einen Reviewer→Implementer-Pfeil; der Slice
braucht damit eine (letzte) Fixrunde, und die DoD-Zeile „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" bleibt offen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). F-6 (31↔42)
geht als Closure-Notiz-Zeile an den Planner, nicht in den Plan.

Der Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation (Modul 11).
