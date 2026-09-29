# Slice kommentar-kennungen-skripte: Kennungs-Form für Kommentare außerhalb von Go — Blockgrenze und Messung, dann Bereinigung in Tranchen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre.

**Bezug:** [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
(Herkunft von Aussagen in Trägern),
[`AGENTS.md`](../../../../AGENTS.md) §3.7 (die Regel gilt für Kommentare in
Code, Konfiguration und Skripten — ihr Werkzeug `make kommentar-kennungen`
liest nur Go-Kommentargruppen) und §3.12.

**Berührte Spec-Stellen:** — (Kommentare in Skripten, Bau- und
Konfigurationsdateien; keine Spec-Stelle).

**Verantwortlich:** Implementer-Agent.

**Autor:** Planner-Agent, angelegt bei der Closure von
`slice-code-kommentare-bereinigung` (deren Plan §1 „Ausdrücklich NICHT“ und
§7 Folge-Slice). **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die §3.7-Form (Herkunft als ein Feld, ein Anker, keine Kette, keine
Kompaktform, kein „ff.“) ist in den Nicht-Go-Kommentaren des Baums messbar:
es gibt eine definierte Blockgrenze je Kommentarform und eine Messung, die
Kandidaten zählt — und der Bestand ist bereinigt oder als begründete, benannte
Restmenge mit Adresse geführt.

**Voraussetzung:** `slice-code-kommentare-bereinigung` (in `done/`) trägt die
Go-Hälfte: Regel, Werkzeug (`tools/harness/kommentar-kennungen`,
`harness/sensors/kommentar-kennungen.md`) und die Go-Restmenge (14
Erzeugnis-Eingabe-Blöcke, §3.7 Ausnahme). Dieser Slice übernimmt Form und
Werkzeuggedanke auf die übrigen Formen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Go-Kommentare.** Tranche des Vorgänger-Slices; `make kommentar-kennungen`
  deckt sie.
- **Erzeugter Code (`gen/`) und SDK-Bäume (`sdks/`).** Erzeugt bzw. von
  `make sdk-public-doc-check` strenger gedeckt — der Slice erfindet keine
  zweite Abdeckung dort.
- **Ein Gate.** Die Gate-Aufnahme-Frage ist in `slice-code-kommentare-bereinigung`
  §7 mit begründetem Ausgang entschieden; dieser Slice stellt sie nicht erneut,
  er schafft nur die Messform. Eine künftige Gate-Frage für Nicht-Go-Formen
  braucht einen eigenen Anlass (Befund, den das Werkzeug verpasst hat).
- **Kommentar-Inhalte.** Der Slice liest Formen (Kennungen, Kette, „ff.“); ob
  ein Satz zutrifft, bleibt Lese-Handlung des Reviewers — die benannte Grenze
  des Werkzeugs.

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — Blockgrenze und Messform.** Je Kommentarform
      (`#` in Shell/Makefile/YAML/Dockerfile, `--` in SQL, Markdown-Kommentare)
      steht die Blockgrenze fest (zusammenhängende Kommentarzeilen vs.
      Grenz-Marker je Form) und ist im Werkzeug oder einem Nebenläufer umgesetzt;
      der Vertrag (`harness/sensors/kommentar-kennungen.md` oder dessen
      Nicht-Go-Pendant) trägt die Definition samt Beispiel. *Zu belegen
      durch:* Tabellentest je Form (Treffer, Blockgrenze, Nicht-Treffer),
      `make test` grün.
- [x] **Liefer-Punkt 2 — Bestands-Messung und Bereinigung.** Basismessung über
      den ganzen Nicht-Go-Baum (zahlt: `tools/`, `harness/mk/`, `Makefile`,
      `tools/schema/*.sql`, `compose.yaml`, `.github/workflows`, Dockerfiles,
      Markdown) mit Befehl, Stand und Zahl je Tranche; Bereinigung in
      Tranchen nach den §3.7-Klassen (nur Kommentarzeilen), Restmenge benannt
      mit `Datei:Zeile`, Klasse und Grund. *Zu belegen durch:* je Tranche
      Diff ohne Nicht-Kommentarzeile, `make test`/`make fmt-check` (soweit
      Formen betreffen) und die Zahlen vor/nach.
- [x] **Liefer-Punkt 3 — Träger.** Living-Docs, die die Kennungs-Form der
      Nicht-Go-Kommentare beschreiben (`AGENTS.md` §3.7,
      `harness/sensors/kommentar-kennungen.md`), tragen die erweiterte
      Geltung; `make docs-check` grün. *Zu belegen durch:* der Diff der
      Träger und der Gate-Lauf.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`); Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6),
      kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: committetes Feld im Plan, Gefundenes und Nichtgefundenes
      je Träger, beide Stände gemessen.
- [x] Doku-Update für den erweiterten Werkzeugvertrag, falls ein öffentlicher
      Vertrag berührt ist (Sensors-Tabelle in `harness/README.md`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — eine
      `evidence/`-Datei oder „kein Anfall“ (in §7 notiert).
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

**Umfang:** M–L — die Basismessung steht noch aus; die Formen-Zahl (Shell,
Make, SQL, YAML, Dockerfile, Markdown) macht Tranchen wahrscheinlich, nicht
Größe.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/kommentar-kennungen/` oder ein Geschwister-Programm | update/neu | Blockgrenze je Form und Kandidaten-Erkennung für Nicht-Go-Kommentare; Go-Pfad bleibt unverändert. |
| `harness/sensors/kommentar-kennungen.md` | update | Vertrag trägt die erweiterte bzw. geteilte Messform. |
| Skripte, `Makefile`, `harness/mk`, `.sql`, `.yml`, Dockerfiles, Markdown des Baums | update (nur Kommentare) | Bereinigung nach §3.7-Klassen in Tranchen. |
| `harness/README.md` §Sensors | update | Zeile `make kommentar-kennungen` trägt die erweiterte Messform (Liefer-Punkt 3). |
| **Plan-Nachzug — Messraum-Entscheidung Markdown:** Markdown ist **nicht** im Messraum. Begründung: Markdown-Prosa ist selbst der Träger — sie hat keine Kommentarform; die einzige Form (`<!-- -->`) ist Inhalts-Form (Zitate, Guidance) und in ADR-/Record-Dateien belegt. Die Kennungs-Linkpflicht in Prosa trägt d-check `ids` als eigene maschinelle Ebene; eine §3.7-Pflicht hier wäre Doppelregulierung. Ausgeschlossen sind außerdem `.harness/baseline/**` (vendored, SHA-gepinnt — kein bearbeitbarer Kommentar-Bestand) und Go (Vorgänger-Slice). Messraum: `#`-Blöcke in `.sh`/`.mk`/`.yml`/`.yaml`/Makefile/Dockerfile und `--`-Blöcke in `.sql`, vollzeilig (Marker am Zeilenanfang nach Leerraum), nachgestellte Kommentaranteile ausgenommen (Trennstelle mehrdeutig — 11 Ketten, 5 Shell + 6 Makefile; übernommen aus der Implementer-Messung, nicht unabhängig nachgemessen — Zählprogramm „Kette" nicht definiert, stützt keine Entscheidung). **Basismessung am Stand `933ea5c0`:** zunächst 149 Kandidaten mit dem Werkzeug-Stand `89d4fa3e` gemessen; das finale Werkzeug (Marker-only-Grenz-Marker, seit T1 im Instrument) misst am selben Stand **160** — die Richtgröße ~150 der Rückführung ist mit dem finalen Instrument überschritten; die geplante Formen-Zerlegung entspricht ihr, die Tranchen sind je Form gefahren und enden je Form bei 0 (T1–T5, Messung mit dem finalen Instrument). | update (Plan) | Rückführung §4: Entscheidung statt Architect-Frage — Markdown ist messbar, aber keine Kommentarform; Schwellen-Ausgang (160 > ~150) wird in §7 geführt. |

**§3.13-Suchlauf (committetes Feld).** Bewegte Eigenschaft: die Kennungs-Zeilen
der Nicht-Go-Kommentarformen (Messraum §3-Nachzug). Gemessen an den Ständen
Werkzeug-Stand `89d4fa3e` (vor der Bereinigung) und Arbeitsbaum (nach T1–T5):

```suchlauf
89d4fa3e 506 -E '^(#|--).*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3})' -- '*.sh' '*.mk' '*.yml' '*.yaml' '*.sql' 'Makefile' 'Dockerfile' 'examples/Dockerfile' 'examples/csharp/Dockerfile' 'examples/kotlin/Dockerfile'
diff 350 -E '^(#|--).*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3})' -- '*.sh' '*.mk' '*.yml' '*.yaml' '*.sql' 'Makefile' 'Dockerfile' 'examples/Dockerfile' 'examples/csharp/Dockerfile' 'examples/kotlin/Dockerfile'
```

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice in `in-progress/`
(WIP-Limit 1); `slice-code-kommentare-bereinigung` liegt in `done/`.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß): die Basismessung meldet mehr als ~150
  Kandidaten oder eine Form, deren Blockgrenze nicht definierbar ist —
  Schnitt entlang der Formen.
- `in-progress` → `open` (blockiert): Markdown als Kommentarform erweist sich
  als nicht messbar (Code-Fences, Zitate) — Architect-Frage, ob Markdown in
  den Messraum gehört.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Review ohne offenes HIGH/MEDIUM +
Verifikation + Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

- **Blockgrenzen sind formenabhängig uneinheitlich** (Shell-Kommentare enden
  an Zeilenenden; zusammenhängende Blöcke sind eine Konvention, keine Syntax).
  *Erwartet, zu belegen durch:* Definition je Form mit Beispiel im Vertrag;
  der Tabellentest hält sie fest.
- **Markdown ist dicht an Doku statt Code** (Code-Fences enthalten `#`-Zeilen,
  Zitate tragen Kennungen). *Erwartet, zu belegen durch:* der Messraum
  entscheidet das explizit (Rückführung `open`).
- **Die Bereinigung berührt Build-Dateien** (`Makefile`, Dockerfiles) — ein
  Kommentarfehler bricht den Bau. *Erwartet, zu belegen durch:*
  `make test`/`make image`-Pfad nach den Tranchen; Diff-Nachweis „nur
  Kommentarzeilen“.
- **Zahlen bewegen sich mit jedem Commit** (Zustandsgröße,
  [`AGENTS.md`](../../../../AGENTS.md) §3.12). *Erwartet, zu belegen durch:*
  je Zahl Befehl, Stand und Lauf.

## 7. Closure-Notiz

- **Was hat funktioniert:** die Formen-Zerlegung — fünf Tranchen je Kommentarform
  (T1 Shell, T2 YAML, T3 Make, T4 SQL, T5 Dockerfile), jede mit Zahlen vor/nach
  in der Commit-Message und dem Diff-Nachweis „nur Kommentarzeilen"; nach T5
  endet die Messung am ganzen Baum bei **14** — exakt der geerbten Go-Restmenge
  des Vorgänger-Slices, 0 Kandidaten in den Nicht-Go-Formen (Verifikation
  `verify-slice-kommentar-kennungen-skripte.md` §2). Die Grenz-Marker-Regel des
  Instruments (eine nur aus `#`/`--` bestehende Zeile beendet den Block) trägt
  die Paragraphen-Struktur des E2E-Runner-Kopfs. Die Messraum-Entscheidung
  (§3-Nachzug) wurde als Rückführung-Entscheidung statt Architect-Frage gefällt
  — die Rückführung `in-progress → open` (Markdown nicht messbar) trat nicht
  ein, und der Schwellen-Ausgang der Rückführung `next` wird durch die
  Formen-Zerlegung getragen: jede Form endet bei 0, die Richtgröße ~150 gilt
  je Form nicht (§Steering-Loop). Das §3.13-Suchlauf-Feld hielt 2/2 Zeilen über
  die Fixrunde hinaus.
- **Was ging anders als geplant:** (1) das Messinstrument wechselte innerhalb
  des Slices — die Grenz-Marker-Regel landete nicht im Werkzeug-Commit
  `89d4fa3e`, sondern in der Tranche T1; dieselbe Basismessung am Stand
  `933ea5c0` liefert 149 mit dem früheren und **160** mit dem finalen
  Instrument, und die Plan-Zahl 149 stand ohne Instrument-Label (Review F-2,
  durch Instrument-Label im §3-Nachzug behoben). (2) Die je-Form-Aufzählung der
  Basismessung ließ die Beispiel-Dockerfiles weg und trug die Summe nicht
  (Review F-4) — in der Fixrunde entfernt; die Totalen sind belastbar.
  (3) Der DIFF-Modus sah Nicht-Go-Kandidaten nicht: der Aufrufer erzeugte den
  Diff-Strom mit `-- '*.go'` (Review F-3) — in der Fixrunde auf die
  Nicht-Go-Formen erweitert; die Vertragszeile `DIFF=<Basis>` dokumentierte die
  Eingabe weiter nur Go und ist mit dieser Closure nachgezogen (LOW-Rest des
  Reviews). (4) Die Zahl „11 gemessene Ketten" im §3-Nachzug ist übernommen
  aus der Implementer-Messung, nicht unabhängig nachgemessen (Zählprogramm
  „Kette" nicht definiert) — dort als übernommen markiert, sie stützt keine
  Entscheidung.
- **Steering-Loop-Eintrag:** die Finding-Klassen des B-Reviews in den Zähler —
  „Nachzug widerspricht dem Nachbarn im selben Träger" (F-1, 1×; Klasse
  bekannt, Heimat `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`),
  „Messwerkzeugwechsel in Tranche — Zahlen nicht komparabel" (F-2, 1×, neu:
  die Basismessung 149 drehte sich mit dem finalen Instrument auf 160; der
  Schwellen-Ausgang 160 > ~150 ist durch die Formen-Zerlegung getragen, je
  Form endet bei 0), „Messform erweitert, Aufrufer-Pfadspec nicht nachgezogen"
  (F-3, 1×, neu); dazu die LOW-Klassen „Aufzählung trägt die Summe nicht"
  (F-4), „Grenze des Sensors unbenannt" (F-5), „Grenzfallestreuung unter
  Formen" (F-6) — je 1×. Lerneintrag: (a) *geschärfte Regel* — der Messraum
  des Werkzeugs ist auf Code- und Konfig-Kommentare festgelegt; Markdown-Prosa
  ist selbst der Träger, die Kennungs-Linkpflicht dort trägt d-check `ids`
  (Vertrag Grenze 2, mit dieser Closure bestätigt und durch die
  Messraum-Entscheidung §3-Nachzug getragen). (b) *neuer Sensor* — keiner: die
  Gate-Frage bleibt entschieden (Plan §1); die Form-Wahrheit bleibt
  Lese-Handlung des Reviewers. (c) *benannte Spec-Lücke* — Messinstrumente
  gehören zur Herkunft einer Zahl: eine Messung trägt neben Befehl und Stand
  das Instrument, sobald ein Lauf das Werkzeug ändert; als Finding-Klasse
  „Messwerkzeugwechsel in Tranche" ins Register eingeschrieben.
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-kommentar-kennungen-skripte.md` in
  `BEO-PGC/kommentar-herkunft-als-kette` ergänzt (die drei MEDIUM-Klassen samt
  Instrument-Korrektur) — Zähler steht damit bei 4×.
- **Restmenge:** 14 (`make kommentar-kennungen COUNT=1` = 14, gemessen am
  Arbeitsbaum, dieser Lauf; 1× `backfill_e2e_test.go`, 13×
  `integration_test.go`); alle Godocs von `func TestE2E*`, alle Klasse
  **Grenze**, alle mit demselben Grund wie im Vorgänger-Slice: die LH- und
  SPEC-Kennungs-Menge des Blocks ist Erzeugnis-Eingabe der
  E2E-Abdeckungstabelle (`abdeckungsAdressiert` liest sie je Testzeile); der
  Schluss-Absatz „Kennungs-Menge:" trägt die Grenze am Block
  (§3.7-Ausnahme). Die Zeilen sind unverändert gegenüber der Restmengen-Tabelle
  des Vorgänger-Slices (`slice-code-kommentare-bereinigung` §7) — geerbter
  Bestand, kein Gegenstand dieses Diffs.

  | # | `Datei:Zeile` | Klasse | Grund |
  |---|---|---|---|
  | 1 | `test/integration/backfill_e2e_test.go:258-273` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 2 | `test/integration/integration_test.go:185-204` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 3 | `test/integration/integration_test.go:308-316` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 4 | `test/integration/integration_test.go:428-443` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 5 | `test/integration/integration_test.go:504-528` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 6 | `test/integration/integration_test.go:639-648` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 7 | `test/integration/integration_test.go:671-680` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 8 | `test/integration/integration_test.go:746-760` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 9 | `test/integration/integration_test.go:819-829` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 10 | `test/integration/integration_test.go:898-915` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 11 | `test/integration/integration_test.go:962-984` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 12 | `test/integration/integration_test.go:1034-1048` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 13 | `test/integration/integration_test.go:1107-1124` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |
  | 14 | `test/integration/integration_test.go:1182-1217` | Grenze | Erzeugnis-Eingabe der E2E-Abdeckungstabelle |

  Der Nicht-Go-Messraum dieses Slices endet bei **0** (alle Formen je Tranche
  bei 0, gemessen am Arbeitsbaum nach T5, Verifikation §2).
- **Folge-Slices:** keiner angelegt — die Messform ist vollständig (Go-Blöcke
  und Nicht-Go-Zeilenkommentare), die Gate-Frage bleibt entschieden. Benannte
  Reste ohne Slice: (a) die Pfadspec-Restlücke `Dockerfile.<Variante>` —
  baumseitig gelesen, diffseitig von keiner Pathspec-Form erfasst; im Bestand
  existiert keine solche Datei, ein Fall wird Befund; (b) die übernommene
  Zahl „11 Ketten" (s. o.).
- **Risiken aus §6:** (1) Blockgrenzen formenabhängig uneinheitlich —
  **entfallen** (Definition je Form im Vertrag §Kandidat samt Blockgrenze und
  Grenz-Marker-Form; Tabellentest `TestLineCommentMarker`, `TestLineBlocks`,
  `TestRunLineForms`; `make test` Exit 0, Verifikation §2). (2) Markdown dicht
  an Doku — **entfallen** (Messraum-Entscheidung §3-Nachzug: Markdown trägt
  den Messraum nicht, Rückführung nicht ausgelöst; d-check `ids` ist die
  maschinelle Ebene für Prosa). (3) Bereinigung berührt Build-Dateien —
  **entfallen** (mechanischer Diff-Nachweis „nur Kommentarzeilen" je Tranche
  mit zwei dokumentierten Sonderfällen: T2-Leerzeile als Blockgrenze,
  T5-gofmt-Ausgleich; `make gates` Exit 0). (4) Zahlen bewegen sich mit jedem
  Commit — **eingetreten, ohne Schaden** (jede Zahl trägt Befehl und Stand;
  die Basismessung trägt seit F-2 das Instrument-Label 149/160 am Stand
  `933ea5c0`).
- **Drei Paarungen:** Anker —
  `harness/sensors/kommentar-kennungen.md` trägt die erweiterte Geltung
  (§Kandidat, Grenze 2, §Test) samt `· seit`-Anker, die Sensors-Zeile in
  `harness/README.md` die Messform; Folge-Slice — keiner (s. o.); Register —
  `BEO-PGC/kommentar-herkunft-als-kette` mit nicht leerem `evidence/`
  (Zähler 4×).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield) quer über die Nicht-Go-Formen; keine eigene Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(2026-09-29): `BEO-PGC/kommentar-herkunft-als-kette` (verkörpert, 3×) —
der Vorgänger-Slice; `BEO-PGC/regel-weiter-als-ihre-sensor` und
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Träger-Pflichten dieses Slices
in Liefer-Punkt 3). Keine weiteren Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
