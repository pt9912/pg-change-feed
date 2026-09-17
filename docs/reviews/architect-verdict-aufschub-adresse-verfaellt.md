# Architect-Verdikt: `BEO-PGC/aufschub-adresse-verfaellt` — die Ereignis-Adresse muss eintreten können

**Rolle:** Architect (Modul 8)

**Anlass:** Der Register-Eintrag `BEO-PGC/aufschub-adresse-verfaellt`
erreicht mit `slice-077` die **3×-Schwelle**. Der Planner-Lauf konnte den
fälligen Architect-Zug nicht selbst dispatchen und hat den Fall beim
Priorisieren von `slice-077` vermerkt (`§8` jenes Plans). Dieser Lauf führt
den Rollenwechsel nach (Modul 8 §Rollen-Sequenz für eine Welle, Schritt 3b
— Verkörperung ist eine **Entscheidung**, keine Planung).

**Rolleninhaber:** pt9912 (Architect-Zug; anderer Eingabe-Kontext als die
Planner-Läufe, die die drei Pläne geschrieben und die Adresse nachgezogen
haben — die drei Belege sind genau deshalb aus jenem Kontext nicht
aufgefallen).

**Datum:** 2026-09-15

**Bezug:** `BEO-PGC/aufschub-adresse-verfaellt` (dieses Zugs Gegenstand) ·
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (verwandter Eintrag,
andere Klasse — unten abgegrenzt) · [`LH-FA-CFG-005`](../../spec/lastenheft.md)
(Bezug des auslösenden Vorgangs `slice-077`) ·
[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
(zweiter Bezug von `slice-077`) ·
[`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
(Bezug von `slice-076`) · Baseline-Regelwerk
`modul-06-roadmap.md` §Was der wellenlose Betrieb selbst auslöst ·
`.claude/commands/plan-welle.md` (künftiger Zielort der Regel) ·
`.claude/commands/implement-slice.md` (Träger der verwandten, bereits
verkörperten Regel) ·
`.harness/baseline/v6.5.0/templates/docs/plan/planning/slice.template.md`
(Quelle der verfallenden Adresse — immutable).

---

## Frage

Der Register-Eintrag steht über der Schwelle und ihm fehlt der Ausgang
(Modul 6 §Das Beobachtungs-Register). Zu entscheiden ist:

1. Ist die neue Klasse eine **zweite Hälfte** der bereits verkörperten Regel
   „Aufschub mit Adresse benennen", eine **eigene** Regel, oder etwas
   Drittes?
2. An **welchem Träger** steht die Regel künftig — erreicht er den
   **Schreiber** der Adresse?
3. Welchen **Ausgang** trägt `state.md` (`verkörpert` / `geplant`)?

## Prüfung der drei Belege — sind sie dieselbe Klasse?

Nicht übernommen, am Bestand nachgeprüft (die zitierten Stellen liegen
teils in `done/`, teils in der Historie):

| Vorgang | Form beim Schreiben | Korrigiert in |
|---|---|---|
| `slice-074` | §2/§7: „Repo **mit** Wellen-Betrieb (`welle-18` offen): die Prüfung läuft bei der nächsten Welle-Closure, auch ohne Wellen-Zugehörigkeit dieses Slice." | `a6cbed9` (Priorisierung) |
| `slice-076` | §2: „im Repo **mit** Wellen-Betrieb von der nächsten Welle-Closure geprüft"; §7: „Repo **mit** Wellen-Betrieb — Prüfung läuft bei der [nächsten Welle-Closure]" | `7f71c60` (Priorisierung) |
| `slice-077` | §2: „im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure" | `0ef7f13` (Priorisierung) |

**Befund: dieselbe Klasse.** Alle drei Fälle sind die **Vorlagen-Default**-
Formulierung der §2/§7-Paarungszeile, die den Träger als **Ereignis** benennt
(„die nächste Welle-Closure"). In allen drei Fällen war die Form beim
Schreiben zulässig (es gab eine offene Welle) und wurde es später nicht mehr:
Nach dem Schließen der Wellen führt die Roadmap keine offene Welle, und eine
Welle-Closure, die die Paarungen einsammeln könnte, tritt nicht mehr ein. Die
drei unterscheiden sich nur in §2 gegen §7 — **nicht** in der Klasse.

**Eine Präzisierung am Beleg:** Das Feld `Quelle` in
`evidence/slice-074.md` nennt als Nachzug-Commit `f4164a9`; die
Paarungszeile wurde tatsächlich in `a6cbed9` gezogen. Die Klasse ist davon
unberührt; die Zitat-Ungenauigkeit gehört zur Kenntnis genommen, nicht
zum Anlass genommen.

## Verdikt 1 — eigene Regel (b), keine zweite Hälfte

**Gewählt: (b) eine eigene Regel.**

**Warum die verwandte Fassung nicht trägt.** Die Regel „Aufschub mit Adresse
benennen" ist in [`.claude/commands/implement-slice.md`](../../.claude/commands/implement-slice.md)
verkörpert und gerichtet an den **Implementer**; die Adresse, die sie verlangt,
ist eine **Folge-Slice-ID** („eine Folge-Slice-ID, die die Doku nachholt"). Der
neue Fall hat eine andere Adress-Klasse (ein **Ereignis**, nicht eine
Slice-Kennung), einen anderen **Schreiber** (den Planner in §2/§7 des
Slice-Plans, nicht den Implementer) und einen anderen Zielort. Eine zweite
Hälfte an der bestehenden Zeile säße damit im **falschen Träger** und erreichte
den falschen Leser.

**Warum nicht trotz Verwandtschaft zusammenlegen.** Das Register hat die
Trennung bereits entschieden: `state.md` des Eintrags hält fest, dass es „zwei
verschiedene Regeln, deshalb zwei Einträge statt eines Zählers" sind. Zwei
Adress-Gültigkeiten sind voneinander unabhängig — die Adresse war hier **nicht
zu eng** (das ist die Umfangs-Klasse des verwandten Eintrags), sondern ihr
**Träger-Ereignis tritt nicht ein**. Und die für den verwandten Eintrag selbst
angedachte zweite Hälfte („… und die Adresse prüft, dass ihr Umfang den
Gegenstand deckt") ist eine **dritte**, noch nicht fällige Bedingung — sie
jetzt mit einer Eintretens-Hälfte zu verschmelzen, verknäuelte zwei Zählräume
und zwei Träger.

**Zur Randnotiz des Auftrags.** Der verwandte Eintrag wird im Auftrag als
„3×, verkörpert" geführt; sein `state.md` weist ihn bei **2×** aus. Verkörpert
ist die **Regel** (in `implement-slice.md`, seit `slice-077`), nicht der
Zählstand des Eintrags. Die Unterscheidung trägt die obige Abgrenzung: die
Regel existiert, ihre Ausweitung auf den Eintretens-Fall ist trotzdem eine
eigene Regel.

## Verdikt 2 — Träger: die Planner-Instruktion, nicht die Vorlage

**Gewählt: [`.claude/commands/plan-welle.md`](../../.claude/commands/plan-welle.md),
Schritt 6 („Slices bereitstellen").**

Der **Schreiber** der Adresse ist der Planner, wenn er §2/§7 eines Slice-Plans
füllt. Genau dort reicht ihn `plan-welle.md` an — Schritt 6 trägt bereits die
„§2-Form-Prüfung nach dem Füllen · seit slice-010", dieselbe Bauform (eine
dateiskopierte Selbstprüfung mit Herkunfts-Anker). Die vendorte
[Slice-Vorlage](../../.harness/baseline/v6.9.0/templates/docs/plan/planning/slice.template.md)
scheidet aus: sie ist Teil der vendorten Baseline, von `make baseline-verify`
gegen `SHA256SUMS` gepinnt, und eine In-Place-Änderung färbte das Gate rot —
dieselbe Begründung wie beim
[`handbuch-versionshistorie`](architect-verdict-handbuch-versionshistorie-uebersprungen.md)-Zug.
`AGENTS.md` scheidet aus, weil die Regel keine repo-weite Hard Rule ist,
sondern eine **Füll-Regel** an einer bestimmten Stelle eines bestimmten
Artefakts; sie gehört zu der Instruktion, die das Füllen beschreibt.

**Zweiter Leser (benannt, nicht eingebaut).** Die Closure-Seite prüft die drei
Paarungen ohnehin ([`.claude/commands/close-welle.md`](../../.claude/commands/close-welle.md)
Schritt 4; im wellenlosen Betrieb trägt sie die Slice-Closure). Eine
symmetrische Zeile dort ist möglich, aber nicht Teil dieses Verdikts — der
Auftrag verlangt **den** Träger, der den Schreiber erreicht.

**Kein Sensor.** Die Prüfung „ist dieses Ereignis noch eingetreten?" hängt am
Zustand der Roadmap (führt sie eine offene Welle?), nicht an einer
Text-Eigenschaft des Plans. Ein Doku-Gate auf die Zeichenkette „Welle-Closure"
wäre blind gegen den legitimen Fall (eine Welle **ist** offen) und müsste die
Roadmap mitlesen — eine Zustands-Kopplung, die kein Modul dieses Repos
trägt. Der Wächter bleibt die Instruktion am Schreiber, die tragende zweite
Linie der unabhängige Reviewer.

## Verdikt 3 — der Ausgang von `state.md`

**`geplant`** — nicht `verkörpert`. Die Regel ist mit diesem Verdikt
**beschlossen** und ihr Träger **benannt**, aber noch **nicht geschrieben**:
Der Träger liegt unter `.claude/commands/**`, und die Setzung dort ist
Folgearbeit (unten). Geschrieben wird sie vom **Planner** im Zuge der
`slice-077`-Closure; sobald die Zeile steht, trägt `state.md`
**`verkörpert`** mit Zielort `.claude/commands/plan-welle.md` und Anker
**`seit slice-077`**.

## Der Regel-Wortlaut (wörtlich, zur Übernahme durch den Planner)

Als neuer Absatz in Schritt 6 von `.claude/commands/plan-welle.md`, neben der
bestehenden §2-Form-Prüfung:

> **Ereignis-Adresse muss eintreten können · seit slice-077**
> (`BEO-PGC/aufschub-adresse-verfaellt`, 3×; Architect-Zug des Lese-Schritts,
> Modul 6): Wird ein Aufschub oder eine Closure-Pflicht an einen Träger
> gebunden, **der ein Ereignis ist** — „die nächste Welle-Closure", „mit dem
> nächsten Release" —, dann muss dieses Ereignis noch **eintreten können**. Ein
> Ereignis-Träger ohne gesichertes Eintreten ist **keine Adresse**: die Form war
> beim Schreiben zulässig und verfällt still, wenn das Ereignis ausbleibt. So
> traf es die §2/§7-Paarungszeile der Slice-Pläne, solange sie „die nächste
> Welle-Closure" nannte und die Roadmap keine offene Welle mehr führte.
> **Ersatz-Träger, wenn er fehlt:** Führt die Roadmap *Offene Wellen* keine
> Welle, trägt die **Slice-Closure selbst** die drei Paarungen
> (Baseline-Regelwerk `modul-06-roadmap.md` §Was der wellenlose Betrieb selbst
> auslöst); §2/§7 nennt dann „die Slice-Closure selbst" statt „die nächste
> Welle-Closure". **Kandidatenlauf beim Füllen:** einen Blick in die Roadmap
> *Offene Wellen* werfen — führt sie keine Welle, ist „die nächste
> Welle-Closure" keine Adresse. **Grenze:** dieselbe Selbstprüfung im selben
> schreibenden Kontext — erste, nicht tragende Linie; die tragende ist der
> unabhängige Reviewer (`.harness/skills/reviewer.md`).

Die Zeile trägt den **Herkunfts-Anker** `· seit slice-077` unmittelbar an der
Regel (die Form, die `implement-slice.md` für die verkörperten Regeln führt).

## Folgearbeit mit Zielrolle

| Vorgang | Zielrolle | Art |
|---|---|---|
| Den Regel-Wortlaut oben in `.claude/commands/plan-welle.md` Schritt 6 einsetzen, mit Anker `· seit slice-077` | Planner | Command (Instruktions-Schärfung) |
| `evidence/slice-077.md` in `BEO-PGC/aufschub-adresse-verfaellt/` anlegen → Zähler steht dann bei 3× | Planner | Register |
| `state.md` auf `verkörpert` setzen (Zielort + `seit slice-077`), **nachdem** die Zeile steht | Planner | Register |
| Optional und nur falls gewünscht: die symmetrische Zeile in `.claude/commands/close-welle.md` Schritt 4 | Planner | Command |
| Optional: die Zitat-Ungenauigkeit in `evidence/slice-074.md` (`f4164a9` → `a6cbed9`) berichtigen | Planner | Register-Beleg |

## Was dieses Verdikt NICHT tut

- **Kein neues ADR.** Dies ist eine Prozess-/Instruktions-Schärfung, kein
  Architektur- oder Vertragsgegenstand; Produktionscode ist nicht berührt —
  dieselbe Einordnung wie bei den Präzedenzfällen.
- **Keine Änderung an der vendorten Baseline** (`.harness/baseline/v6.5.0/…`):
  die Vorlage ist referenziert und `baseline-verify`-gepinnt.
- **Kein Sensor/Gate** — siehe Verdikt 2.
- **Keine Setzung am Träger, am Register, am Slice-Plan oder an der Roadmap** —
  allesamt Folgearbeit (Tabelle oben).
