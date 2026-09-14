# Architect-Verdikt: Coverage-Gate — fällige Reifestufe (Hochschalt-Trigger)

**Rolle:** Architect (Modul 8)

**Anlass:** Trigger-Audit der `welle-18`-Closure (Modul 6
§Wellen-Closure-Prozedur, Schritt 2), **Reifestufen-Zweig** — Modul 8
§Rollen-Sequenz für eine Welle, Schritt 2: Planner → **Architect** → Planner.
Der Planner hat den fälligen Trigger **nicht still übergangen**, sondern in
`welle-18-results.md` §Verifikation (Trigger-Audit, Absatz *Bootstrap-aware
Gates*) als offenen Träger benannt und an diese Rolle gereicht.

**Rolleninhaber:** pt9912 (Architect-Zug; anderer Kontext als der
Planner-Lauf der `welle-18`-Closure, als die Implementer-/Reviewer-/Verifier-
Läufe von `slice-068` und als der Implementer-Lauf von `slice-049`)

**Datum:** 2026-09-14

**Bezug:**
[`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a) *Eskalationsklausel* und §Re-Evaluierungs-Trigger (b) ·
[`harness/sensors/coverage-gate.md`](../../harness/sensors/coverage-gate.md)
§Kalibrierungs-Bindung · `harness/mk/coverage.mk` · `AGENTS.md` §3.6 (Gates
ohne ADR), §3.7 (Ist-Zustand) · `welle-14-results.md` (Präzedenz-Lesart
desselben Triggers) · `verify-slice-049.md` §2 (dokumentierte
Lauf-zu-Lauf-Schwankung) · `verify-slice-068.md` §2 (45,80 % unabhängig
bestätigt) · `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`
(Registereintrag derselben Sub-Area, 1×)

**Erzeugtes Artefakt dieses Zugs:** die Adresse der Umsetzung — `slice-076`
(`docs/plan/planning/open/`); die endgültige Planung (Priorisierung,
`Verantwortlich:`, Wellen-Zuordnung) führt der Planner. Der Slice-Plan wird
als **Kennung** zitiert, nicht als Pfad-Link: Ein Slice wechselt die
Lifecycle-Ablage, und die Klasse
`BEO-PGC/slice-pfad-als-link-in-berichten` steht bei 2×.

---

## Verdikt

1. **Der Hochschalt-Trigger ist fällig.** `ADR-0054`s Eskalationsklausel
   bindet die Hochschaltung an den **real gemessenen Ist-Stand gegen die
   nächste 5-%-Stufe** — nicht an eine Wellen-Anzahl, nicht an ein Datum,
   nicht an eine Stabilitätsbedingung. `THRESHOLD` = 35 %, nächste Stufe =
   40 %, realer Ist-Stand = **45,80 %** → die Lücke zur nächsten Stufe ist
   geschlossen.
2. **Neue Stufe: 40 %** — **ein** Schritt, nicht mehrere („die Lücke zur
   **nächsten** Stufe"; „jede Stufe braucht ihren eigenen Beleg
   (Hochschalt-Trigger), keinen Freibrief").
3. **Kein Folge-ADR, kein `Supersedes`.** Die Hochschaltung ist die von
   `ADR-0054` selbst vorgesehene Bewegung ihres bootstrap-aware Gates;
   `AGENTS.md` §3.6 verlangt einen ADR für die **Senkung**, nicht für diese
   Bewegung. `ADR-0054` bleibt `Accepted` und unverändert.
4. **Form: eigener — wellenloser — Slice `slice-076`, kein
   Architect-Commit.** Der Träger ist mit der Adresse entschieden, nicht mehr
   „offen".
5. **Bindung der neuen Stufe:** 40 % ist ein **Boden**, kein Vorschlag.
   Fällt ein künftiger Slice darunter, ist die Antwort Test-Arbeit, nicht
   eine Schwellen-Senkung (die wäre nach §3.6 ADR-pflichtig).

---

## Befundlage (eigene Messung dieses Laufs)

| Lauf | Ergebnis |
|---|---|
| `make coverage-gate THRESHOLD=35` (im `make gates`-Lauf dieses Zugs) | `coverage-gate: OK — Coverage 45.80% erfüllt Schwelle 35%` |
| `make coverage-gate THRESHOLD=40` (eigener, ungepiped ermittelter Lauf; Exit-Code in eigenem Schritt geprüft) | `coverage-gate: OK — Coverage 45.80% erfüllt Schwelle 40%`, Exit **0** |

Der Ist-Stand **45,80 %** ist damit in diesem Lauf real gemessen und deckt
sich mit der unabhängigen Bestätigung aus `verify-slice-068.md` §2 und mit
der Zahl, die `welle-18` §Verifikation führt. Der Hochschalt-Kandidat ist zum
Zeitpunkt der Entscheidung **grün** — die Eskalation kann das Gate nicht
selbst rot färben. Genau das garantiert die Klausel-Bedingung „Ist-Stand ≥
nächste Stufe": fällig ist der Schritt erst, wenn der Wert ihn trägt.

**Grenze des Belegs (benannt, nicht verschwiegen):** Der erste
`make gates`-Lauf dieses Zugs endete **Exit 2** — nicht wegen der hier
entschiedenen Sache, sondern weil `commit-traceability` im 5-Commit-Fenster
eine Struktur-ID im Betreff fand (`SPEC-020` in einem Commit des parallel
laufenden `slice-069`-Zugs, `AGENTS.md` §5). Nach der lokalen Betreff-Korrektur
durch diesen fremden Zug lief `make commit-traceability` allein Exit **0**,
`make a-check` Exit **0**; `baseline-verify` und der `d-check`-Doku-Lauf waren
im ersten Lauf bereits grün (54 Dateien bzw. 538 Dateien / 0 Befunde).

---

## Frage 1 — Ist der Trigger nach dem eigenen Wortlaut fällig?

**Zitat der Klausel** (`ADR-0054` §(a), wörtlich):

> „`THRESHOLD` startet bei einer dokumentierten ersten Stufe (gemessener
> Ist-Stand, abgerundet auf den nächsten vollen 5-%-Schritt, niemals über
> 80 %), mit Hochschalt-Trigger „nächste Coverage-Verbesserung schließt die
> Lücke zur nächsten Stufe" bis 80 % erreicht ist, dokumentiert als
> Kalibrierungs-Bindung in `harness/README.md` §Sensors"

und zwei Sätze später:

> „… jede Stufe braucht ihren eigenen Beleg (Hochschalt-Trigger), keinen
> Freibrief."

Die Sensor-Definition führt denselben Trigger wortwörtlich
(`harness/sensors/coverage-gate.md` §Kalibrierungs-Bindung: „die nächste
Coverage-Verbesserung schließt die Lücke zur **nächsten 5-%-Stufe**").

Drei gelesene Fassungen — entschieden wird an der ersten:

| Lesart | Inhalt | Trägt? |
|---|---|---|
| **W1 — wert-getrieben (gewählt)** | Fällig, wenn der real gemessene Ist-Stand die nächste 5-%-Stufe der geltenden Schwelle erreicht oder überschreitet. | **ja.** Der Wortlaut nennt als Kriterium nur den Wert („die Lücke zur nächsten Stufe") und sonst nichts; die Einstiegsformel daneben bindet ebenso an den Wert (Ist-Stand → 5-%-Schritt). 35 → nächste Stufe 40; Ist 45,80 ≥ 40 → **fällig**. |
| W2 — ereignis-/kalender-getrieben (Wellen-Anzahl, Datum, „stabil über N Läufe") | Fällig erst nach N Wellen, zu einem Termin oder nach einer Stabilitätsbedingung. | **nein.** Keines der drei Elemente steht in der Klausel. „Hochschalten, sobald die Stufe stabil erreicht ist" hätte ein **Maß** nennen müssen; sie nennt keins. W2 ist eine Ergänzung des Wortlauts, keine Lesart. |
| W3 — kein Trigger ohne einen benannten Verbesserungs-**Vorgang** | Fällig nur, wenn ein mit Kennung benannter Vorgang die Lücke geschlossen hat. | **nein, in dieser Form.** Der Wortlaut nennt die „Coverage-Verbesserung" als Auslöser, nicht als Vorgangs-Kennung; ihr Beleg ist die Messung. Der Präzedenzfall liest dieselbe Klausel wert-getrieben: `welle-14-results.md` §Verifikation führt den Trigger als „**noch nicht fällig** … aktueller Ist-Stand (39,6–39,8 %) bleibt über der Einstiegsstufe" — dort lag 39,6 unter der nächsten Stufe 40, hier liegt 45,80 darüber. **Dieselbe Regel, entgegengesetztes Ergebnis.** |

**Fazit Frage 1: fällig.** Es gibt kein Wortlaut-Element, das die
Hochschaltung an eine Wellen-Anzahl, ein Datum oder eine Stabilitätsbedingung
hängt, und es gibt keines, das die reale Messung durch einen Vorgangs-Namen
ersetzt. `THRESHOLD` steht seit `slice-049` unverändert auf 35 %, während der
reale Ist-Stand real von 39,6 % auf 45,80 % gestiegen ist.

## Frage 2 — Auf welche Stufe?

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: `THRESHOLD` bleibt 35 % | kein Eingriff; das Gate ist auf jeden Fall grün | das Gate verliert seine Wirkung: der Ist-Stand liegt **10,8 Prozentpunkte** über der Schwelle — ein Slice könnte mehr als zwei volle Stufen Coverage verlieren, ohne dass ein Gate rot wird. Das ist genau die „stehengebliebene Reifestufe", die Modul 6 §Wellen-Closure-Prozedur Schritt 2 als unzulässigen Wellen-Abschluss nennt. |
| B — nichts tun, aber als Carveout dokumentieren | macht die Nicht-Bewegung sichtbar | ein Carveout ist die Ausnahme für **unerreichbare** Ziele (Modul 7). Hier ist das Ziel erreichbar und der Beleg liegt vor (eigener Grün-Lauf auf 40 %) — der Carveout hätte kein Objekt. |
| **C — ein Schritt auf 40 % (gewählt)** | exakt der Wortlaut „die **nächste** Stufe"; jeder Schritt trägt seinen eigenen Beleg („keinen Freibrief"); 5,8 Prozentpunkte Puffer — rund 29× die einzige dokumentierte Schwankungsquelle (`verify-slice-049.md` §2: 39,6 % vs. 39,8 % bei identischem Code) → die Stufe hält ohne eigene Stabilitätsbedingung; der Schritt ist im Moment seines Vollzugs real grün (eigene Messung, s. o.); folgt der d-check-Reifung, die ebenfalls schrittweise lief (85 → 90 → 93) | ein Schritt weniger, als der Ist-Stand heute trägt — die Schwelle bleibt vorerst unterhalb von `floor(Ist-Stand)` |
| D — auf 45 % (die Stufe, die der Ist-Stand heute trägt) | entspricht der Einstiegsformel, auf heute angewandt (`45,80 → 45`) | 0,8 Prozentpunkte Puffer; ein einzelner Slice mit Coverage-Verlust färbt das Gate rot, und die dann fällige Antwort ist **nicht** die Senkung (§3.6) — die Klausel hebt stufenweise an („nächste Stufe", „jede Stufe … ihren eigenen Beleg"), nicht auf den Grenzwert; die d-check-Reifung sprang nie zwei Stufen |
| E — mehrere Stufen bis 80 % in einem Zug | schnellstes Erreichen der Endstufe | der Direktsprung ist die von `ADR-0054` **verworfene** Option B ihrer eigenen Alternativen-Tabelle („ein Direktsprung ohne Kenntnis des Ist-Stands kann das Gate bei jedem ersten Lauf dauerhaft rot färben"); mit Ist 45,80 % wäre 80 % sofort rot |

**Fazit Frage 2: 40 %.** Der Ist-Stand trägt heute auch 45 %; die Klausel hebt
aber stufenweise an, und jede Stufe braucht ihren **eigenen** Beleg. Die
nächste Stufe ist damit nicht vergessen, sondern an ihren eigenen Trigger
gebunden: **sobald ein weiterer Coverage-Fortschritt die Lücke zur Stufe
45 % schließt (Ist-Stand ≥ 45 %), ist die nächste Hochschaltung nach
derselben Prozedur fällig** — dann wieder über den Reifestufen-Zweig.

## Frage 3 — Nachweis der neuen Stufe, und was bei Coverage-Verlust gilt

- **Nachweis** (dieselbe Beleg-Form wie bei der Einstiegsstufe in
  `slice-049`): ein realer **Grün**-Lauf auf der neuen Stufe
  (`make coverage-gate` mit `THRESHOLD=40`) und ein realer **Rot**-Beleg
  unmittelbar über ihr (`THRESHOLD=46`) — der zeigt, dass die Stufe wirklich
  prüft und nicht leer läuft. Beides ist DoD-Punkt von `slice-076`; der
  Grün-Lauf liegt aus diesem Zug bereits vor (Exit 0, s. §Befundlage).
- **Trägt der Puffer ohne Bedingung? Ja.** 45,80 − 40 = 5,8 Prozentpunkte
  gegen eine dokumentierte Lauf-zu-Lauf-Schwankung von ±0,2 Prozentpunkten
  (`verify-slice-049.md` §2). Eine Bindung „hochschalten erst nach N stabilen
  Läufen" wäre eine Bedingung ohne Objekt — und sie stünde nicht in der
  Klausel (Frage 1, W2).
- **Wenn ein künftiger Slice Coverage verliert:** Die Stufe ist **Boden**.
  Die Antwort ist die Wiederherstellung der Coverage (Tests), nicht die
  Rücknahme der Schwelle. Eine Schwellen-**Senkung** ist nach `AGENTS.md`
  §3.6 ein ADR, kein PR-Kommentar; `ADR-0054`s Klausel schließt sie für die
  Zwischenstufen aus („niemals über 80 %" für den Einstieg, „bis 80 %
  erreicht ist" für die Richtung der Bewegung). Das ist die eine Konsequenz,
  die dieses Verdikt bindend festhält — und sie ist der Grund, warum die
  Stufe nicht auf den Grenzwert 45 % gesetzt wird (Frage 2, Option D).

## Frage 4 — Form: eigener Slice oder Architect-Commit?

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (Trigger weiter als offener Träger führen) | kein Aufwand | genau der Zustand, den dieses Verdikt auflösen soll; der Trigger bliebe „offen", obwohl er entschieden ist |
| B — Architect-Commit („reine Konfigurationszeile": `THRESHOLD ?= 40` + Doku) | ein Commit, kein Lifecycle-Aufwand; Präzedenz: Architect-Verdikte haben **Regel-Text** direkt verkörpert (`AGENTS.md` §3.9, §3.10) | der Grün-/Rot-Beleg bliebe eine Behauptung im Verdikt statt ein DoD-Punkt; keine Review, keine Verifikation an einem Diff; Modul 8, Schritt 2 reicht das Verdikt **an den Planner zurück** — die Umsetzung fiele in dasselbe Kontextfenster, das das Rollen-Modell trennt |
| **C — wellenloser Slice `slice-076` (gewählt)** | die Änderung ist eine **ausgelieferte** Änderung an einem repo-weiten Vertrag (jede künftige Slice-DoD hängt an „`make gates` grün") → sie bekommt DoD, Review und Closure wie jede andere Vertrags-Änderung; der Beleg wird ein prüfbarer DoD-Punkt statt Verdikt-Prosa; Präzedenz: der **Einstiegswert** 35 % wurde in `slice-049` ausgeliefert, nicht als ADR-Anhang — die Hochschaltung ist dieselbe Änderungsklasse (Schwellen-Wert + Beleg); der Träger ist mit der Adresse entschieden statt verschoben | ein Lifecycle-Durchlauf für eine Konfigurationszeile; die neue Stufe wird erst mit der Umsetzung wirksam |

**Fazit Frage 4: Slice.** Der Unterschied zu `AGENTS.md` §3.9/§3.10 ist die
Änderungsklasse: dort ein **Dokument**, hier ein **Gate-Parameter**, dessen
Nachweis ein Lauf ist. Ein Gate-Parameter, der ohne DoD-Punkt und ohne Review
wandert, ist genau die Form, gegen die §3.6 argumentiert.

## Was `slice-076` ändern muss — Datei-Ebene

| Datei | Änderung |
|---|---|
| `harness/mk/coverage.mk` | `THRESHOLD ?= 40`; der Kopf-Kommentar nennt den Ist-Zustand (§3.7), nicht die Stufen-Chronik |
| `harness/sensors/coverage-gate.md` | §Kalibrierungs-Bindung: geltende Stufe 40 %; die Rot-/Grün-Beleg-Zeile dieser Sektion wird auf die neu gemessenen Werte nachgezogen |
| `harness/README.md` | §Sensors-Zeile `make coverage-gate`: die geltende Stufe |
| `AGENTS.md` | §4-Zeile: die **Rampe** nennen (Einstieg → Endstufe), den beweglichen Wert **nicht** doppeln — die geltende Stufe hat mit `harness/mk/coverage.mk` und der Sensor-Definition bereits zwei Orte, ein dritter driftet (§3.7) |

`ADR-0054` bleibt unverändert; seine Fitness-Function-Zeile nennt die Schwelle
bereits als „aktuell gültige" mit Verweis auf die Kalibrierungs-Bindung — sie
trägt die Bewegung ohne Textänderung.

---

## Feststellung zur Closure-Sequenz der `welle-18`

Der Planner hat den Trigger korrekt als fällig erkannt und **nicht** still
übergangen — aber die Welle ist geschlossen worden, während die
Reifestufen-Entscheidung noch offen war. Modul 6 §Wellen-Closure-Prozedur
Schritt 2 nennt „eine stehengebliebene Reifestufe" neben dem stillen roten
Gate und der überfälligen Entscheidung als einen der drei Zustände, mit denen
eine Welle **nicht** schließen darf; die Übergabe Planner → Architect gehört
in **diesen** Schritt, nicht hinter ihn. Der Fall ist damit eingetreten und
**benannt**, und er wird nicht gezählt: Der abgeschlossene Vorgang dieser
Klasse wäre die Welle, und ein Registerbeleg entsteht bei der Slice-Closure
(Modul 6 §Das Beobachtungs-Register). Tritt dasselbe Muster ein zweites Mal
auf, ist es ein Registerkandidat — diese Feststellung ist die Stelle, an der
es dann nachgelesen werden kann.

---

## Was dieser Zug geändert hat — und was nicht

**Geändert:**

- `docs/plan/planning/open/slice-076-…` (neu — Adresse der Umsetzung)
- `docs/plan/planning/done/welle-18-results.md` — **nur** der Trigger-Audit-
  Absatz *Bootstrap-aware Gates*: der offene Träger bekommt den Zeiger auf
  dieses Verdikt und die Kennung `slice-076`; die Notiz wird nicht
  umgeschrieben
- diese Verdikt-Datei

**Nicht geändert:** `harness/mk/coverage.mk`, `harness/sensors/coverage-gate.md`,
`harness/README.md`, `AGENTS.md` §4 (führt `slice-076`), `docs/plan/adr/**`
(kein Folge-ADR — `ADR-0054` bleibt `Accepted`), `internal/**`, `spec/**`.

**Offen für den Planner-Zug:** Priorisierung von `slice-076` (`open` → `next`,
`Verantwortlich:`), Wellen-Zuordnung (wellenlos, solange keine Welle ein
*Mehr* über seine DoD beobachtet) und der Hinweis in der nächsten
Wellen-Closure, dass die Reifestufe 40 % zu diesem Zeitpunkt getragen ist.
