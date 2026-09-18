# ADR-0099: `{from: slice/welle, to: review}` — Regel zurückgenommen, `ADR-0097`/`ADR-0094`s Selbst-Zitat-Begründung bestätigt

**Status:** Accepted — Supersedes [`ADR-0097`](0097-observation-matrixklasse-review-verboten.md)
(nur die Frage `{from: slice, to: review}`/`{from: welle, to: review}`; die
übrigen Festlegungen dort — `observation`-Matrixklasse,
`{from: observation, to: review, allow: false}` — bleiben unverändert)

**Datum:** 2026-09-18

**Autor:** Architect (pt9912)

**Bezug:** [`ADR-0094`](0094-review-matrixklasse-kennung-statt-adresse.md)
(Entscheidung, `{from: slice, to: review}` bewusst **nicht** zu regeln —
Begründung: Selbst-Zitat, gemeinsame Archivierung, kein Hänger-Risiko),
[`ADR-0097`](0097-observation-matrixklasse-review-verboten.md) (bestätigt
dieselbe Auslassung ein zweites Mal, benennt einen konkreten
Re-Evaluierungs-Trigger — **Supersedes** hier ausschließlich auf diesem
einen Punkt), `AGENTS.md` §3.5 (ADR-Immutabilität — der Anlass dieser ADR
ist ihre Verletzung), §3.6 (Gate-Lockerung/-Verschärfung ohne ADR
verboten), Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-
Prozedur Schritt 4 (Review-Reports archivieren gemeinsam mit ihrem Slice
in `done/<welle-id>/archiv.zip`, ohne eigenen Stub) · Beobachtungs-Register
`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` (zweites Auftreten
derselben Fehlerklasse „Gate-Scope wächst durch die Konfigurationsdatei
selbst, ohne den `AGENTS.md` §3.6-Träger").

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0094`/`ADR-0097`).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Commit `01b7b09` (Anlass: ein zu weiter `**`-Glob der `slice`-/`welle`-
Matrixklassen in `.d-check.yml`, der versehentlich in
`docs/plan/planning/observations/**/evidence/slice-*.md` hineinreichte)
hat **zwei Änderungen in einem Commit** vorgenommen: (1) die Pfad-Listen
von `slice` eng gefasst und eine neue Klasse `welle` mit ebenso engen
Pfaden eingeführt — eine reine Bugfix-Verfeinerung, unstrittig — und (2)
zwei neue Matrix-Regeln ergänzt, `{from: slice, to: review, allow: false}`
und `{from: welle, to: review, allow: false}`. Die Commit-Message nennt
als Grund einzig einen „Nutzerhinweis nach eigenem Fehler in mehreren
Slice-Dateien dieser Session" — kein Bezug auf eine ADR, keine
`Supersedes`-Zeile.

Das ist kein Formfehler, sondern eine inhaltliche Kollision mit **zwei**
bereits `Accepted` ADRs:

- `ADR-0094` §Entscheidung, wörtlich: „Bewusst **keine** neue Regel für
  `slice → review` (ein Slice zitiert seinen eigenen Review-Report
  routinemäßig und beide werden gemeinsam archiviert, kein
  Hänger-Risiko)."
- `ADR-0097` §Entscheidung, wörtlich: „Bewusst keine Regel für
  `{from: slice, to: review}` in dieser ADR: […] zitiert teils den
  eigenen, teils einen fremden Review — die Unterscheidung braucht
  Einzelfallprüfung, kein blanket Verbot." Und §Re-Evaluierungs-Trigger:
  Die Regel wird erst erneut geprüft, „[wenn] die verbleibenden
  `{from: slice, to: review}`-Fälle […] wiederholtes Auftreten
  (Beobachtungs-Register, 3×)" erreichen — **nicht** durch eine
  Ad-hoc-Ergänzung während eines unabhängigen Bugfixes.

Die Folgerunde (Commit `b88dbf6`) hat die neue Regel real durchgesetzt:
21 Dateien, 46 Befunde, davon **11 Dateien** ausschließlich wegen der
beiden neuen `slice`/`welle`-Regeln betroffen (die übrigen fallen unter
die bereits gültige `observation → review`-Regel aus `ADR-0097` und
bleiben von dieser ADR unberührt). Diese ADR entscheidet, ob die neue
Regel Bestand hat, indem sie die realen 11 Fälle inhaltlich prüft — nicht
indem sie die Frage erneut abstrakt stellt.

### Befund an den 11 realen Fällen

Alle 11 betroffenen Dateien geprüft (vor und nach `b88dbf6`,
`git show b88dbf6~1:<pfad>` gegen den Fixrunden-Diff):

| Datei | Zitierte Review-Datei(en) | Beziehung |
|---|---|---|
| `slice-001-bootstrap.md` | `docs/reviews/review-report.template.md` | Selbst — eigenes Plan-Nachzug-Artefakt | <!-- d-check:status-provenance -->
| `slice-036-antragsqueue-sql-funktionen.md` | `review-slice-036.md`, `-fixrunde.md`, `verify-slice-036.md` | Selbst | <!-- d-check:status-provenance -->
| `slice-037-administrations-goroutine-live-reload.md` | `review-slice-037.md`, `-fixrunde.md`, `verify-slice-037.md` | Selbst | <!-- d-check:status-provenance -->
| `slice-038-cli-diagnose.md` | `review-slice-038.md`, `-fixrunde.md`, `verify-slice-038.md` | Selbst | <!-- d-check:status-provenance -->
| `slice-archive-altbestand-adr.md` | `review-slice-archive-altbestand-adr.md` | Selbst | <!-- d-check:status-provenance -->
| `slice-archive-altbestand-vollzug.md` | `review-slice-archive-altbestand-vollzug.md` | Selbst | <!-- d-check:status-provenance -->
| `welle-13-results.md` | `architect-verdict-retention-loeschausfuehrung.md`, `architect-verdict-slice-chronik-in-code-kommentar.md` | Selbst — Reviews/Verdikte der eigenen Welle bzw. eines ihrer eigenen Slices | <!-- d-check:status-provenance -->
| `welle-6-results.md` | `architect-review-welle-6.md` | Selbst — Review der eigenen Welle | <!-- d-check:status-provenance -->
| `welle-archive-altbestand-results.md` | `review-slice-archive-altbestand-adr.md`, `-vollzug.md` | Selbst — Reviews der eigenen zwei Slices | <!-- d-check:status-provenance -->
| `welle-beispiele-start-ueber-make-results.md` | vier `review-slice-beispiele-*.md` | Selbst — Reviews der eigenen vier Slices | <!-- d-check:status-provenance -->
| `welle-d-check-results.md` | zwei `verify-slice-*.md`, `verify-welle-d-check.md`, zwei `architect-verdict-*.md` | Selbst — Belege der eigenen Welle bzw. ihrer eigenen Slices | <!-- d-check:status-provenance -->

**Kein einziger der 11 Fälle ist eine Fremd-Zitation** (kein Slice
zitiert den Review eines *anderen* Slices, keine Welle den Review einer
*anderen* Welle). Das ist genau die Kategorie, die `ADR-0094` als
risikofrei eingestuft hat: Slice-Datei, Welle-Ergebnisnotiz und ihre(r)
eigene(n) Review(s)/Verdikt(e) sammelt dieselbe Welle-Closure ein und
archiviert sie **gemeinsam** in dieselbe `done/<welle-id>/archiv.zip`
(Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur
Schritt 4). Verifiziert: `docs/reviews/*.md` selbst ist bislang nie
verschoben worden (`git log --follow` liefert keine Rename-Historie für
z. B. `review-slice-036.md`) — es liegt weiterhin unter `docs/reviews/`, <!-- d-check:status-provenance -->
wie `AGENTS.md` §3.5 letzter Absatz es als **Record** ohne eigene
Wanderung beschreibt; die Wanderung, vor der `ADR-0094`/`ADR-0097` warnen,
ist die **künftige** Zip-Archivierung der Welle, noch nicht eingetreten.

Der `Re-Evaluierungs-Trigger` aus `ADR-0097` — „3× im Beobachtungs-
Register" — ist zum Zeitpunkt von Commit `01b7b09` **nicht** eingetreten:
Es existiert keine Beobachtungs-Registerzeile zu wiederholten
`slice → review`-Fällen; die 11 Fälle wurden nicht durch drei getrennte,
protokollierte Beobachtungen aufgedeckt, sondern in einem einzigen,
ungeplanten Aufwasch während eines unabhängigen Glob-Bugfixes gefunden
und sofort mechanisiert.

### Was die Mechanisierung real gekostet hat

Die Fixrunde hat jede der 11 Selbst-Zitationen von einem klickbaren
relativen Markdown-Link (Beispiel einer der obigen Zeilen: Link auf den
zugehörigen Review-Report) in unverlinkte Prosa (Form „Review zu
<Kennung>") umgeschrieben — dieselbe Zielform, die
`ADR-0094`/`ADR-0095` für **Fremd-Zitationen** (`adr →
review`, `observation → review`) richtig gewählt haben, weil dort der
Basisname sonst tatsächlich hängt. Für eine Selbst-Zitation entsteht
dadurch aber ein reiner Navigierbarkeits-Verlust ohne Gegenwert: der
Leser eines abgeschlossenen Slice- oder Welle-Dokuments verliert den
Ein-Klick-Pfad zu dessen eigenem Review, ohne dass dafür irgendein
Hänger-Risiko entfiele — es gab keins.

## Entscheidung

Wir nehmen die beiden Regeln `{from: slice, to: review, allow: false}`
und `{from: welle, to: review, allow: false}` aus `.d-check.yml` zurück
und stellen `ADR-0094`s ursprüngliche Position wieder her — mit derselben
Begründung, jetzt mit den 11 realen Fällen belegt statt nur angenommen.
Die begleitende Bugfix-Verfeinerung aus Commit `01b7b09` (die engen
Pfad-Listen der `slice`-Klasse, die neue `welle`-Klasse selbst) bleibt
bestehen — sie ist unabhängig von der zurückgenommenen Regel richtig: die
`welle`-Klasse trennt Welle-Dokumente korrekt von `slice`/`observation`,
unabhängig davon, ob eine `to: review`-Regel an ihr hängt.

Die 11 betroffenen Dateien werden auf ihre vorige, verlinkte Form
zurückgesetzt (`git checkout b88dbf6~1 -- <pfad>` je Datei, geprüft: keine
der 11 Dateien wurde seit `b88dbf6` durch einen anderen Commit erneut
berührt). Die 10 unter der weiterhin gültigen `observation → review`-Regel
(`ADR-0097`) korrekt umgeschriebenen Beobachtungs-Register-Dateien bleiben
unverändert in Prosa-Kennung-Form — sie sind von dieser ADR nicht
betroffen.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Regel bestätigen: Selbst- vs. Fremd-Zitat spielt in der Praxis keine Rolle, Prosa-Kennung deckt beide gleich gut | ein Weniger an Ausnahmen, konsistente Zielform über alle Quellklassen | widerspricht den realen 11 Fällen — kein einziger war Fremd-Zitat, der Navigierbarkeits-Verlust hat keinen Gegenwert; kassiert zudem `ADR-0094`s explizit begründete Grenze ohne neuen Gegenbeleg |
| B — Regel verfeinern: Ausnahme für eine erkennbare Selbst-Zitat-Form (z. B. `slice-NNN` zitiert nur Reviews mit `-slice-NNN` im Basisnamen) statt vollständiger Rücknahme | behält die Fremd-Zitat-Deckung, die `ADR-0097` bewusst offengelassen hat, mechanisch | keiner der 11 realen Fälle war je ein Fremd-Zitat-Risiko; eine Erkennungsregel für „passt der Basisname zur eigenen Slice-/Welle-Kennung" ist selbst fehleranfällig (z. B. `welle-13-results.md` zitiert `architect-verdict-*`-Dateien ohne `welle-13` im Namen) und würde echte, harmlose Selbst-Zitate weiter falsch meldet — mehr Mechanik für ein Risiko, das die Beleglage nicht zeigt |
| **C — Regel vollständig zurücknehmen, `ADR-0094`s Position bestätigen — gewählt** | deckt sich mit der realen Beleglage (11/11 Selbst-Zitate); stellt die durch zwei `Accepted`-ADRs bereits getroffene, richtige Entscheidung wieder her; kein Mehraufwand für ein nicht beobachtetes Risiko | Fremd-Zitat-Fälle (sollten sie künftig auftreten) bleiben wie schon in `ADR-0097` reine Review-Prüfpflicht statt Gate — unverändert gegenüber dem Vorher-Zustand, kein neuer Verlust |

**Fazit:** C.

## Konsequenzen

- Positiv: Die durch `01b7b09`/`b88dbf6` verursachte, ungedeckte
  ADR-Reversion ist behoben; 11 Dateien haben ihre klickbaren
  Selbst-Zitat-Links zurück; `ADR-0094`s Begründung ist jetzt mit realer
  Beleglage bestätigt statt nur behauptet.
- Negativ: Keiner — die zurückgenommene Regel deckte in ihrer kurzen
  Lebensdauer keinen einzigen realen Fremd-Zitat-Fall; ihr einziger realer
  Effekt war der oben beschriebene Navigierbarkeits-Verlust.
- Folgepflicht: (1) `.d-check.yml` bereinigt (diese ADR nennt den exakten
  Diff im Kontext-Abschnitt), (2) die 11 Dateien zurückgesetzt, (3)
  `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` um einen zweiten
  Beleg ergänzt (zweites reales Auftreten derselben Fehlerklasse — Zähler
  2×, noch nicht 3×-fällig), (4) `make gates` grün, dann Commit.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` Modul `matrix` | `.d-check.yml`s `rules:`-Liste enthält **keine** Regel `{from: slice, to: review}` oder `{from: welle, to: review}` — geprüft durch Lesen der Datei (kein Sensor kann die *Abwesenheit* einer möglichen Regel positiv erzwingen; die Fitness-Function dieser ADR ist deshalb der reale Config-Stand selbst, review-geprüft) | `make docs-check` (bleibt grün, weil die 11 Dateien wieder verlinkt sind und keine `matrix-forbidden`-Regel mehr greift) |

## Re-Evaluierungs-Trigger

Beobachtbarer Trigger: ein **echter** Fremd-Zitat-Fall (ein Slice zitiert
den Review eines *anderen* Slices, oder eine Welle den Review einer
*anderen* Welle) tritt real auf **und** betrifft eine Datei, die noch
nicht archiviert ist. Dann Einzelfallprüfung im Review (wie von
`ADR-0097` bereits vorgesehen); erreicht diese Fehlerklasse 3× im
Beobachtungs-Register, Folge-ADR mit `Supersedes ADR-0099` zur erneuten
Mechanisierungsfrage — diesmal ggf. mit einer Regel, die tatsächlich nur
Fremd-Zitate trifft, nicht pauschal jede `slice`/`welle → review`-Kante.
Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — `{from: slice, to: review}`/`{from: welle, to: review}` zurückgenommen, `ADR-0094`s Selbst-Zitat-Begründung mit 11 realen Fällen bestätigt; Supersedes `ADR-0097` auf diesem einen Punkt | Architect-Zug (Korrektur einer ungedeckten Ad-hoc-Regeländerung in Commit `01b7b09`) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0099`.
