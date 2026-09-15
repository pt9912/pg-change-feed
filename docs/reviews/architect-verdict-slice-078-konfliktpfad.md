# Architect-Verdikt: `slice-078` — Konflikt-Pfad (F-1, F-2, F-3)

**Rolle:** Architect (Baseline-Regelwerk `modul-08-agentenrollen.md`
§Konflikt-Pfad als Rollen-Sequenz)

**Anlass:** [`review-slice-078`](review-slice-078.md) meldet einen HIGH und zwei
MEDIUM, alle drei Regelfragen — `slice-078` hat 42 host-lokale Pfade entfernt
und das `hostpaths`-Modul ohne Ausnahme aktiviert. Der Reviewer hat die
Regel-Lage als **Architect-Urteil** übergeben (Review §Verdikt). Dieser Zug
entscheidet sie und schreibt die Entscheidungen.

**Rolleninhaber:** pt9912 (Architect-Zug; anderer Kontext als der
Implementer-Lauf, der korrigiert hat, und als der Reviewer-Lauf)

**Datum:** 2026-09-15

**Bezug:** [`review-slice-078`](review-slice-078.md) (Übergabe-Artefakt) ·
[`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md),
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
[`ADR-0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md),
[`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) (dieses
Zugs Entscheidung) · `AGENTS.md` §3.5, §3.7, §3.11 ·
[`harness/sensors/docs-check.md`](../../harness/sensors/docs-check.md) ·
`.d-check.yml` (`modules`, `structure`)

**Erzeugte Artefakte dieses Zugs:** die Entscheidung `ADR-0075` + ihre
Index-Zeile; dieses Verdikt. Die Folgearbeit (Zielrolle, unten) ist als Adresse
benannt, nicht als Slice angelegt — Priorisierung und Schnitt führt der Planner.

---

## Befund-Korrektur vorab — die Commit-Attribution von F-1

Die Review schreibt den Eingriff in `ADR-0072` §Entscheidung Punkt 5 dem Commit
`aab7aa8` zu. Das ist **nicht** der tragende Commit: `aab7aa8` ändert in
`ADR-0072` nur den Abschnitt `### Der gemessene Befund` (§Kontext, Zeile 67 ff.)
und die Zählung des Verdikts — ein Vorgang **innerhalb** der Zitat-Klasse, denn
`ADR-0073` Punkt 1 nimmt §Kontext nicht von der Klasse aus. Der angezeigte
Unterabschnitt `Entwurf der Hard Rule §3.11` (`ADR-0072:173-190`) wurde dagegen
in `df47282` geändert (Falsch-Beispiel, Richtig-Beispiel, Klammertext) — belegt
durch `git log -S"Host-Wurzel"` gegen `ADR-0072` und den Commit-Diff.

**Der Befund trägt, die Attribution wird korrigiert:** der Eingriff überschreitet
die Klassengrenze — er berührt §Entscheidung, und der geänderte Klammertext ist
keine Form eines Verweises, sondern Erklärung.

## Verdikt je Finding

- **F-1 (HIGH) — bestätigt, am korrigierten Commit `df47282`; Verdikt
  „Lockerung legitim, aber undokumentiert" → Folge-ADR.** Der Eingriff in
  `ADR-0072` §Entscheidung Punkt 5 ist in der Sache richtig — die reale Form
  muss weg (Auftrag), die Hausform ist
  [`ADR-0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md), und
  der Klammertext bleibt nur als Form-Beschreibung wahr (`AGENTS.md` §3.7) —,
  also wird er per [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md)
  Punkt 2 als **beschlossener** Text festgestellt, nicht zurückgenommen. Die
  anderen zwei Wege tragen nicht: **Rücknahme** setzt eine reale Form in den
  Fence zurück, der repo-weite Grep wäre nicht mehr 0 (der Sensor fängt es
  nicht, der Grep schon — der Auftrag „entfernt" wäre verfehlt); die
  **Klassen-Erweiterung** öffnet den schmalen Kanal, den `ADR-0073` eng hält,
  und heilt den Klammertext nicht (er ist keine Form eines Verweises).

- **F-2 (MEDIUM) — bestätigt; die Disposition entscheidet: keine
  Lokator-Ersatzpflicht.** Der Implementer-Grund („fällt nicht unter die
  in-place-Klasse") war falsch — `ADR-0073` Punkt 1 nimmt Zeilen-/
  Bereichs-Lokatoren ausdrücklich in die Klasse auf —, aber die Klausel, die
  ihren Ersatz **anordnete** (`ADR-0072` §Entscheidung Punkt 3, letzter Satz),
  ist mit Punkt 3 durch
  [`ADR-0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md)
  supersedet und nicht restituiert; ihr Ersatztext fasst nur die Zitationsform,
  und „Der Rest bleibt" führt die Lokator-Klausel nicht auf.
  [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) Punkt 3
  stellt die Disposition **supersedet und nicht restituiert** fest — die
  verbleibenden Lokatoren (`ADR-0054`, ein `done/`-Zeitdokument) sind kein
  Defekt dieses Vorgangs, kein Nachzug.

- **F-3 (MEDIUM) — bestätigt; Gewinner „Fences eingeschlossen", festgezogen an
  der Hard Rule.** `AGENTS.md` §3.11 ist die Regel, ein Slice-Plan ist ein
  Zeitdokument (Modul 5/6) und trägt keine Regel;
  [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) Punkt 1
  deklariert die Reichweite an dieser einen Stelle, die Folgearbeit zieht §3.11
  nach. Eine verbotene Form zeigt ein Dokument nur als Platzhalter; die
  Konsequenz — die Regel ist **strenger als ihr Sensor** (das Modul lässt
  Fences frei) — ist als Lücke benannt, nicht still.

- **F-4 (LOW) — bestätigt; `ADR-0073` Punkt 2 ist die präzisere Fassung und
  gilt.** `ADR-0073` Punkt 2 gibt den §3.5-Wortlaut vor („host-lokale Pfade,
  Linkziele, Zeilen-Lokatoren"); die ausgelieferte §3.5-Form lässt
  `Zeilen-Lokatoren` weg und fasst `Formfehler der Zitation` weiter als die
  ADR. Der Implementer gleicht `AGENTS.md` §3.5 an die ADR-Vorgabe an.

- **F-5 (INFO) — Handlungsbedarf ja, klein.** Die Ist-Zustand-Modulliste in
  [`harness/sensors/docs-check.md`](../../harness/sensors/docs-check.md) (Zeile
  18-21) wiederholt die Deklaration aus `.d-check.yml`, obwohl der Absatz selbst
  sagt, dass die Konfiguration die Deklaration ist (`AGENTS.md` §3.7,
  Zwei-Quellen-Disziplin). Die Liste entfällt, der Zeiger auf `.d-check.yml`
  bleibt.

- **F-6 (INFO) — Handlungsbedarf ja; Abgleich im nicht-immutablen Träger.**
  `ADR-0072` führt 31 (die Modul-Befunde zum Entscheidungszeitpunkt), die
  Korrektur umfasst 42 Vorkommen. Die immutabile ADR wird nicht geändert; der
  Abgleich steht im §Kontext von
  [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) und
  gehört in die Closure-Notiz.

- **F-7 (INFO) — der Platzhalter ist die richtige Form.** Unter der Reichweite
  „Fences eingeschlossen" (F-3) deckt die Regel ihr eigenes Beispiel; ein
  reales, sensor-fangbares Beispiel wäre selbst ein Verstoß und bräche den
  repo-weiten Grep 0. Die DoD-Formulierung „reale Beispiele" trägt der
  **Referent** (die Hausform, der benannte Anker), nicht das verbotene Segment.

## Folgearbeit mit Zielrolle

| Vorgang | Zielrolle | Art |
|---|---|---|
| `AGENTS.md` §3.11 auf die **Fassung 2** aus [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) Punkt 2: Reichweite Fences eingeschlossen, benannte Lücke (Regel deckt Fences, das Modul nicht); das Falsch/Richtig-Paar bleibt Platzhalter/Hausform | Implementer | Hard Rule |
| `AGENTS.md` §3.5 auf den Wortlaut von `ADR-0073` Punkt 2 angleichen (`Zeilen-Lokatoren` aufnehmen, `Formfehler der Zitation` auf die ADR-Form zurückführen) — F-4 | Implementer | Hard Rule |
| [`harness/sensors/docs-check.md`](../../harness/sensors/docs-check.md) §Grenze (Punkt 8) um „die Regel deckt Fences, das Modul nicht" ergänzen | Implementer | Sensor-Doc |
| [`harness/sensors/docs-check.md`](../../harness/sensors/docs-check.md) Zeile 18-21: die Ist-Zustand-Modulliste entfällt, der Zeiger auf `.d-check.yml` bleibt — F-5 | Implementer | Sensor-Doc |
| Closure-Notiz: 31→42-Abgleich (F-6), Risiko-Ausgänge, Register | Planner | Slice-Closure |

**An welcher Stelle die neue Grenze steht:** die Reichweite der Regel lebt ab
diesem Zug in **`AGENTS.md` §3.11** (getragen von
[`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) Punkt 1);
der Slice-Plan trägt keine Regel-Fassung — er ist Zeitdokument.

## Die Grenze, die der Implementer in der Fixrunde anwenden darf

- **Reichweite:** die Regel deckt die ganze Markdown-Fläche, Fences
  eingeschlossen; eine verbotene Form wird **nur** als Platzhalter gezeigt, die
  Hausform steht für das richtige Zitat, die reale Form nirgends.
- **Anfassen darf er:** `AGENTS.md` §3.5 und §3.11;
  [`harness/sensors/docs-check.md`](../../harness/sensors/docs-check.md).
- **Nicht anfassen:** `ADR-0051`, `ADR-0054`, `ADR-0072`, `ADR-0073`,
  `ADR-0074` (Accepted, immutabel — nur Zitat-Korrektur, und die ist hier nicht
  nötig); **kein Lokator-Ersatz** (F-2); **keine** Scope-/`ignore`-/
  `exempt-paths`-Ausnahme; **keine** §Entscheidung eines `Accepted`-Dokuments;
  kein Register, kein Plan, `.d-check.yml` unverändert.

## Kein DoD-Nachzug

Der Slice braucht die **Fixrunde**; die DoD-Zeile „Review durchgeführt, Report
unter `docs/reviews/` liegt vor" bleibt offen, bis die Folgearbeit gelaufen und
die Review-Auflagen abgearbeitet sind (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde).

## Was dieser Zug geändert hat — und was nicht

**Geändert:**

- [`docs/plan/adr/0075-hostpaths-reichweite-und-wortlaut.md`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md)
  (neu — die Reichweite, der §3.11-Entwurf Fassung 2, die Lokator-Disposition)
- [`docs/plan/adr/README.md`](../plan/adr/README.md) (eine neue Zeile;
  `ADR-0072`-Zeile um den zweiten Teil-Supersede annotiert)
- diese Verdikt-Datei

**Nicht geändert:** `AGENTS.md` (§3.5/§3.11 — Folgearbeit),
[`harness/sensors/docs-check.md`](../../harness/sensors/docs-check.md)
(Folgearbeit), `.d-check.yml` (kein Ausschlussblock, unverändert),
`ADR-0051`/`0054`/`0072`/`0073`/`0074` (immutabel), der Slice-Plan, jedes
Register, `spec/**`, `internal/**`. Die Umsetzung der Folgearbeit ist
Implementer-/Planner-Arbeit, nicht Teil dieses Zugs.
