# Welle welle-20: Coverage 80 % über der netzlos prüfbaren Fläche

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<NN>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-15.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

**„Die Coverage" dieses Repos erreicht 80 %** — die Gate-getragene Unit-Zahl
über der **netzlos prüfbaren Fläche** ([`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md);
1679 Statements, und „die Coverage" ohne Subjekt-Zusatz meint genau diese Zahl,
nie die DB-Adapter-Coverage). Der Weg dorthin hat zwei Hälften: zuerst eine
**Präzisierung des Gegenstands** (`slice-079`: die drei Pakete, deren Tests ohne
externen Dienst überspringen, verlassen den Nenner — der Ist-Stand springt real
von 49,3 % auf ~69,7 %, **ohne eine Zeile Test**) und danach **Test-Arbeit** an
dem, was ungedeckt bleibt. Getragen wird das Maß von `make coverage-gate` gegen
`THRESHOLD`; die Schwelle wandert nach dem unveränderten
bootstrap-aware-Mechanismus ([`ADR-0054`](../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a)) stufenweise bis 80.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `slice-079` (Scope-Schnitt und Neukalibrierung des Coverage-Gates) liegt in
  `done/` — **noch nicht erfüllt.** Ohne ihn misst die Welle gegen einen Nenner,
  den [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  ersetzt hat, und ihr eigener Fortschritt wäre nicht belegbar.
- [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  ist `Accepted` — **erfüllt** (der Gegenstand ist entschieden).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/`.
- `make coverage-gate` ist **real grün bei `THRESHOLD=80`** über der netzlos
  prüfbaren Fläche, mit eigenem Grün-Beleg dieses Laufs. Das ist das *Mehr*:
  keine einzelne Slice-DoD kann eine repo-weite Schwelle belegen.
- Ein **Rot-Beleg unmittelbar über der Endstufe** ist geführt (`THRESHOLD=85`,
  Exit ≠ 0) — sonst prüfte die Stufe nicht real.
- `make gates` grün.
- Closure-Notiz in `welle-20-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| `slice-079` | Coverage-Gate: Scope-Schnitt und Neukalibrierung | [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) |

**Als Nächstes zu schneiden — der Schnitt folgt der Messung.** Erst wenn
`slice-079` liegt, steht die Zahl über dem neuen Nenner; dann entstehen die
Test-Slices nach dem **Größenmaß aus dem Schnittvorschlag des Verdikts**
(Statement-Anteil): `internal/bootstrap` 591 Statements / 327 ungedeckt (der
Hebel, ggf. zwei Slices — der zweite erst nach dem ersten), `cmd/pg-change-feed`
49 / 49, Rest-Tail ~130. **Ihre Kennungen entstehen mit ihren Dateien** — kein
Name ohne Adresse, die die Sendung annimmt (die Klasse aus
`BEO-PGC/aufschub-adresse-verfaellt`).

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- **Wird blockiert von:** `slice-079` (der Scope-Schnitt) — ohne ihn misst die
  Welle gegen einen Nenner, den [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  ersetzt hat.
- **Blockiert:** nichts.
- **Verwandt, aber nicht Teil:** die DB-gestützte Messung (`ADR-0071` Punkt 3)
  und die Executor-Naht (Punkt 5) — beide eigene Vorgänge. Die Naht würde die
  Decke heben, ohne das Verfahren zu ändern; sie ist deshalb ausdrücklich
  **nicht** der Weg zu diesen 80 %.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Die Endstufe 80 % selbst** — sie steht ([`ADR-0054`](../adr/0054-coverage-gate-und-benchmark-infrastruktur.md),
  [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md));
  sie zu ändern wäre eine Schwellen-Senkung und nach `AGENTS.md` §3.6
  ADR-pflichtig. Diese Welle **erreicht** sie, sie verhandelt sie nicht.
- **Die DB-gestützte Ebene** (`ADR-0071` Punkt 3): eigene, subjekt-qualifizierte
  Messung im nicht-blockierenden `e2e.yml`, eigener Vorgang. `make gates` bekommt
  in dieser Welle **keinen** Container.
- **Die Executor-Naht** (`ADR-0071` Punkt 5): design-begründeter eigener Vorgang,
  nicht Beigabe. Ein Slice dieser Welle ändert **kein Produktionsverhalten**, um
  Coverage zu gewinnen — er prüft Verhalten.
- **Der Messmechanismus** (`-coverpkg`, Docker-Stage, Gate-Skript): mit
  `slice-079` eingestellt. Wer ihn weiter ändern will, braucht eine Entscheidung,
  keinen Test-Slice.
- **Die Integrations-/E2E-Fläche** (`test/integration/**`, die Compose-Kette):
  sie hat ihren eigenen Beleg-Träger (`make test-integration`), und ihre Coverage
  ist nicht Gegenstand dieses Maßes (`ADR-0054` §(a), `ADR-0071` Punkt 1).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: <Zeiger auf `welle-<NN>-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
