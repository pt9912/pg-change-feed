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
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`); Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6),
      kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: committetes Feld im Plan, Gefundenes und Nichtgefundenes
      je Träger, beide Stände gemessen.
- [x] Doku-Update für den erweiterten Werkzeugvertrag, falls ein öffentlicher
      Vertrag berührt ist (Sensors-Tabelle in `harness/README.md`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — eine
      `evidence/`-Datei oder „kein Anfall“ (in §7 notiert).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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
| **Plan-Nachzug — Messraum-Entscheidung Markdown:** Markdown ist **nicht** im Messraum. Begründung: Markdown-Prosa ist selbst der Träger — sie hat keine Kommentarform; die einzige Form (`<!-- -->`) ist Inhalts-Form (Zitate, Guidance) und in ADR-/Record-Dateien belegt. Die Kennungs-Linkpflicht in Prosa trägt d-check `ids` als eigene maschinelle Ebene; eine §3.7-Pflicht hier wäre Doppelregulierung. Ausgeschlossen sind außerdem `.harness/baseline/**` (vendored, SHA-gepinnt — kein bearbeitbarer Kommentar-Bestand) und Go (Vorgänger-Slice). Messraum: `#`-Blöcke in `.sh`/`.mk`/`.yml`/`.yaml`/Makefile/Dockerfile und `--`-Blöcke in `.sql`, vollzeilig (Marker am Zeilenanfang nach Leerraum), nachgestellte Kommentaranteile ausgenommen (Trennstelle mehrdeutig — 11 gemessene Ketten, 5 Shell + 6 Makefile). **Basismessung am Stand `933ea5c0` mit dem erweiterten Werkzeug:** 149 Kandidaten gesamt, davon 135 Nicht-Go (sh 62, yml 19, sql 15, mk 14, yaml 13, Makefile 4, Dockerfile 3) und 14 Go-Restmenge — unter der Rückführungs-Schwelle (~150). | update (Plan) | Rückführung §4: Entscheidung statt Architect-Frage — Markdown ist messbar, aber keine Kommentarform. |

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

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag:** *(zu tragen bei Closure)*
- **Beobachtungs-Register (`../observations/`):** *(zu tragen bei Closure)*
- **Restmenge:** *(zu tragen bei Closure — `Datei:Zeile`, Klasse, Grund je
  Kandidat, oder „0“)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang, zu tragen bei Closure)*
- **Drei Paarungen:** *(zu tragen bei Closure)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield) quer über die Nicht-Go-Formen; keine eigene Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(2026-09-29): `BEO-PGC/kommentar-herkunft-als-kette` (verkörpert, 3×) —
der Vorgänger-Slice; `BEO-PGC/regel-weiter-als-ihre-sensor` und
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Träger-Pflichten dieses Slices
in Liefer-Punkt 3). Keine weiteren Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
