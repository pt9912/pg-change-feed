# ADR-0162: Festlegungen `SPEC-038`/`SPEC-039` mit Herkunft aus `ADR-0158` bis `ADR-0161`; Umzugs-Commit eines aufgelösten MR-Eintrags vor `make doc-immutable`

**Status:** Accepted

**Datum:** 2026-10-07

**Autor:** pt9912 (Architect-Rolle, Modul 8; Architect-Zug zu
`slice-spec-festlegungen-harness-werkzeuge`, anderer Kontext als der
Planner- und der Implementer-Lauf des Slice)

**Bezug:** [`ADR-0158`](0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) ·
[`ADR-0159`](0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) ·
[`ADR-0160`](0160-teil-range-leer-test-und-hostpaths-home-relativ.md) ·
[`ADR-0161`](0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md) ·
[`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md) (gemessen/hergeleitet) ·
`AGENTS.md` §3.5, §3.6, §3.12 · Architect-Verdikt zum Werkzeug von
`slice-zitat-vergleich-werkzeug` §2 bis §5 · d-check v0.82.0, Modul `vcs`.

**Schärft:** [`SPEC-038`](../../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge),
[`SPEC-039`](../../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
— die Kante von `ADR-0158` bis `ADR-0161` zu diesen Festlegungen läuft über
diese ADR (die vier tragen `Schärft: —` und bleiben unverändert).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`spec/pflichtenheft.md` §7 „Festlegungen der Harness-Werkzeuge“ trägt seit
`slice-spec-festlegungen-harness-werkzeuge` zwei Festlegungen: `SPEC-038`
(`make zitat-vergleich`) und `SPEC-039` (Leer-Test der Teil-Range vor
`make doc-immutable`). Die Spec nennt keine ADR, weil `matrix` die Richtung
Spec → ADR verbietet (gemessen vom Implementer, Slice-Plan §1). Der adoptierte
Stand v6.16.0 (`grundlagen-referenz-richtung.md` §Spec-Straten) verlangt, dass
die ADR eines Werkzeugs mit `Schärft:` auf dessen Spec-Stelle zeigt. Die vier
Quell-ADRs sind `Accepted` und tragen `Schärft: —`; ein Nachtrag dort wäre ein
Überschreiben (`AGENTS.md` §3.5).

Zwei Teile der Quellen sind abgelöst und können nichts mehr schärfen:
`ADR-0159` Entscheidung 2 (MR-Pins, samt dem Absatz „Nicht messbarer Referent an
einem MR-Eintrag“) und `ADR-0160` Entscheidung 2 sind laut Status-Zeile von
`ADR-0161` durch diese ersetzt. Inhalte aus beiden stehen trotzdem in
`SPEC-038` bzw. `SPEC-039`.

`SPEC-038` gibt das Werkzeug wieder, wie es seit `slice-zitat-vergleich-werkzeug`
misst. Das Architect-Verdikt zum Werkzeug hat fünf Abweichungen vom `bash`-Block
in `ADR-0159` Entscheidung 4 freigegeben (gestapelte `id`, `id` in anderer Form,
CommonMark-Fence, Tag-Paar an den Baseline-Bäumen, Locale). Fünf weitere hat der
Vertrag aus den Review-Fixrunden als „Auslegung im Sinn von Verdikt §2“
geführt: `id` in eingerücktem Code und HTML-Kommentar nicht gelesen, Lokator nur
`L<a>-<b>` mit `1 <= a <= b`, schließende `#`-Folge nicht im Slug, mehrdeutige
Stellung endet mit Exit 2, HTML-Kommentar in einer `id`-Zeile gilt als ohne
Inhalt. Verdikt §2 sagt dagegen wörtlich: eine Abweichung außerhalb der fünf
freigegebenen Punkte ist ein Fehler des Werkzeugs. Diese fünf hat also keine
Entscheidung getragen.

Der Slice löst `MR-001` per `git mv` nach `harness/conventions/done/` auf
(Baseline-Regel, `harness/conventions.md` §Adaptions-Block) und legt `MR-006` als
Nachfolger an. `make doc-immutable` meldet den Umzug als `core-drift-vcs`
(„immutable Datei gelöscht oder umbenannt“). Es ist die erste Auflösung eines
MR-Eintrags im Repo.

## Entscheidung

1. **Die Kante.** Diese ADR schärft `SPEC-038` und `SPEC-039`. Je Teil der
   Festlegung gilt diese Herkunft:

   | Teil | Quelle |
   |---|---|
   | `SPEC-038` Gegenstand | `ADR-0158` (Bedingung (b) neu), `ADR-0161` Entscheidung 5 (Adaptions-Durchgang); die Abgrenzung zur Form-Korrektur ist `ADR-0161` Entscheidung 3 und 4 (zweiter Spiegelstrich), dort keine Festlegung |
   | `SPEC-038` Einheit je Verweisform | `ADR-0158` Entscheidung 1 (ohne die Zeile zur HTML-`id`), `ADR-0159` Entscheidung 1 |
   | Heading, Slug, Fence, HTML-`id`, Zeile ohne Inhalt, mehrdeutige Stellung | `ADR-0158` Entscheidung 1, `ADR-0159` Entscheidung 1, Entscheidung 2 dieser ADR |
   | roh und Normalisierung | `ADR-0159` Entscheidung 3 (sie trägt die Sätze von `ADR-0158` Entscheidung 2, die bleiben), Entscheidung 4 vierte Zusage, Entscheidung 5 (Schluss-Umbruch, Tag gleich fremder Version); Tag-Paar an den Baseline-Bäumen: Entscheidung 2 dieser ADR |
   | Stände | `ADR-0158` Entscheidung 3, `ADR-0161` Entscheidung 5 |
   | nicht messbarer Referent | `ADR-0158` Entscheidung 4; der Zweig für einen Adaptions-Eintrag: Entscheidung 1 (c) dieser ADR |
   | Ausgänge und Beleg | `ADR-0158` Entscheidung 5, `ADR-0159` Entscheidung 4 Zusagen eins bis drei, `ADR-0161` Entscheidung 5 (Lesart von `cmp 0`/`cmp 1`/Exit 2 im Durchgang); Locale und Zahl der Argumente: Entscheidung 2 dieser ADR |
   | `SPEC-039` | `ADR-0160` Entscheidung 1, `ADR-0161` Entscheidung 4 erster Spiegelstrich, Entscheidung 3 dieser ADR |

   (a) **Nicht Quelle** sind die abgelösten Teile: `ADR-0158` Entscheidung 2
   und 6, `ADR-0159` Entscheidung 2, `ADR-0160` Entscheidung 2, der `bash`-Block
   in `ADR-0159` Entscheidung 4 (historische Fassung). Der Satz in `SPEC-039`
   „Jede Teil-Range ist entweder leer oder grün; ein Exit 2 … ist ein Befund und
   kein Leerfall“ folgt aus `ADR-0160` Entscheidung 1 (b) und (c), nicht aus
   deren abgelöster Entscheidung 2.
   (b) **Grenzen, die beim Werkzeug bleiben.** Gleicher Körper, Slug mit HTML
   vor dem Zeilenende und NUL-Bytes (`ADR-0159` Entscheidung 5) sind Grenzen der
   Nachbildung. Sie stehen in `harness/targets/zitat-vergleich.md` §Grenzen,
   nicht in §7.
   (c) **Nicht messbarer Referent an einem Adaptions-Eintrag.** Ein
   Adaptions-Eintrag hat die Aussage-Abschnitte einer ADR nicht. Löst die alte
   Adresse an keinem Stand auf, ist eine Korrektur dort Urteil am Diff nach
   `ADR-0073`, und der Beleg nennt „nicht messbar“ mit Grund. Diese ADR trägt
   den Zweig, den `ADR-0159` Entscheidung 2 eingeführt hatte; `ADR-0161` hat
   ihn mit dem Pin-Commit abgelöst, nicht in der Sache verworfen.

2. **`SPEC-038` gibt die Messung des Werkzeugs wieder — das ist Entscheidung,
   keine Auslegung mehr.** Alle zehn Punkte, an denen
   `tools/harness/zitat-vergleich.sh` vom `bash`-Block in `ADR-0159`
   Entscheidung 4 abweicht, gelten als Festlegung: die fünf aus dem
   Architect-Verdikt §3 bis §4 (gestapelte `id`, `id` in anderer Form,
   CommonMark-Fence, strenges Tag-Paar, Locale mit Fähigkeitsprobe) und die fünf
   aus den Review-Fixrunden (`id` in eingerücktem Code und HTML-Kommentar nicht
   gelesen, Lokator nur `L<a>-<b>` mit `1 <= a <= b`, schließende `#`-Folge nicht
   im Slug, mehrdeutige Stellung endet mit Exit 2, HTML-Kommentar in einer
   `id`-Zeile ohne Inhalt). Keiner der zehn widerspricht einer `Accepted`-ADR:
   - Der Lokator ist die Form des Arguments, nicht des Verweises. Eine einzelne
     Zeile („eine Zeile“, `ADR-0158` Entscheidung 1) wird als `L<n>-<n>`
     gemessen.
   - Der Slug ist eine Nachbildung des Renderers. Die Semantik „Anker auf ein
     Heading“ meint den Anker des Renderers, und der setzt die schließende
     `#`-Folge nicht in den Slug.
   - Eingerückter Code und HTML-Kommentar erweitern die Erkennungsregel „eine
     `id` in Fence oder Inline-Code zählt nicht“ (`ADR-0159` Entscheidung 1) auf
     die übrigen Stellen, an denen der Renderer keinen Anker erzeugt.
   - Mehrdeutig mit Exit 2 ist fail-closed und hält die Zusage „nie `cmp 0`“.
   - Der Kommentar in einer `id`-Zeile setzt die Auslegung „Zeile ohne Inhalt“
     aus Verdikt §3 fort.

   Verdikt §2 gilt damit nur noch mit dieser Regel: Weicht das Werkzeug von
   `SPEC-038` ab, ist das Werkzeug falsch. Ändert sich die Festlegung, wird
   `SPEC-038` fortgeschrieben, und diese ADR bleibt ihre Kante. Eine
   Semantik-Änderung gegen eine `Accepted`-ADR bleibt eine Folge-ADR.

   **Wortlaut-Berichtigung.** In der Tabelle „Stellung“ unter `SPEC-038` lautet
   die dritte Zeile künftig: „in einer Zeile ohne Inhalt, die nächste Zeile mit
   Inhalt ist kein Heading | der Block ab der nächsten Zeile mit Inhalt bis vor
   das nächste Heading beliebiger Ebene oder bis zum Dateiende“. Gemessen
   beginnt der Block dort, nicht an der `id`-Zeile (Fitness Function, Zeile 4).

3. **Umzugs-Commit eines aufgelösten MR-Eintrags.** Ein MR-Eintrag wird
   aufgelöst, indem er per `git mv` nach `harness/conventions/done/` wandert.
   Dieser Umzug steht in einem eigenen Commit `M`:
   - `M` ist ein reiner Umzug: jede Zeile von `git diff -M --name-status M~1 M`
     ist `R100`, von `harness/conventions/MR-<NNN>-<titel>.md` nach
     `harness/conventions/done/` unter demselben Dateinamen, und `M` hat genau
     einen Parent.
   - Die Message nennt die aufgelöste Kennung und diese ADR.
   - Der Index in `harness/conventions.md` und der Nachfolge-Eintrag stehen in
     anderen Commits.

   `make doc-immutable` läuft dann um `M` in den zwei Teil-Ranges `B..M~1` und
   `M..H`, je mit dem Leer-Test aus `SPEC-039`. `M` ändert nur Adaptions-Einträge
   und fällt damit unter den Wortlaut von `SPEC-039`. Liegen in einer Range
   mehrere solcher Commits (Form-Commit `F`, Umzugs-Commit `M`), wird an jedem
   geteilt, in der Reihenfolge der Commits. Am Umzugs-Commit selbst prüft der
   Verifier mit dieser Befehlsform (`bash`, `git`, `awk`, `sed`, `wc`;
   Host-Werkzeuge nach `AGENTS.md` §3.1):

   ```bash
   umzug() {  # umzug <M> — Beleg am Umzugs-Commit einer MR-Auflösung
     local m out
     m=$(git rev-parse --verify -q "$1^{commit}") || { echo "umzug: $1 löst nicht auf, Exit 2"; return 2; }
     [ "$(git rev-list --parents -n 1 "$m" | wc -w)" -eq 2 ] || { echo "umzug: $1 hat nicht genau einen Parent, Exit 2"; return 2; }
     out=$(git diff -M --name-status "$m~1" "$m")
     [ -n "$out" ] || { echo "umzug: $1 ändert nichts, Exit 2"; return 2; }
     printf '%s\n' "$out" | sed 's/^/umzug: /'
     if printf '%s\n' "$out" | awk -F'\t' '
         { n = split($2, p, "/"); f = p[n] }
         !($1 == "R100" && $2 ~ /^harness\/conventions\/MR-[0-9]+-[^\/]+\.md$/ && $3 == "harness/conventions/done/" f) { bad = 1 }
         END { exit bad }'; then
       echo "umzug: $1 ist ein reiner Umzug nach harness/conventions/done/, Exit 0"; return 0
     fi
     echo "umzug: $1 ist kein reiner Umzug nach harness/conventions/done/, Exit 1"; return 1
   }
   ```

   Exit 0 ist der Beleg. Exit 1 heißt: `M` ist kein reiner Umzug. Der Umzug
   wird dann neu geschnitten, nicht ausgelassen. Nach dem Umzug bleibt die Datei
   unter `done/` geschützt: `vcs.paths` (`harness/conventions/**/MR-[0-9]*.md`)
   trifft sie weiter. Die Konfiguration in `.d-check.yml` bleibt unverändert.

4. **Der `formnorm`-`cmp` am Form-Commit** (`ADR-0161` Entscheidung 4, zweiter
   Spiegelstrich) gehört nicht zu `SPEC-038` oder `SPEC-039`. Seine Festlegung
   übernimmt `slice-spec-festlegungen-pruefer-hooks`, der den Vertrag
   `harness/targets/pin-stale.md` samt Bump-Ablauf führt. Bis dahin bleibt die
   Befehlsform in `ADR-0161`. Die Befehlsform `umzug` aus Entscheidung 3 bleibt
   in dieser ADR und wird nicht in §7 aufgenommen. Sie liest eine einzige
   Klassifikation von `git` (`R100`), und ihre Randformen stehen vollständig
   oben. Das ist ein akzeptiertes Negativ.

## Verglichene Alternativen

**Frage 1 — die Kante.**

| Option | Pro | Contra |
|---|---|---|
| A — `Schärft:` in `ADR-0158` bis `ADR-0161` nachtragen | kein neues Dokument | überschreibt `Accepted`-ADRs (`AGENTS.md` §3.5); verboten |
| B — die Spec nennt ihre Quell-ADRs | Herkunft direkt lesbar | `matrix` verbietet Spec → ADR (gemessen vom Implementer, zwei Befunde `matrix-forbidden`); falsche Richtung (SDP) |
| C — nichts tun, `SPEC-038`/`SPEC-039` ohne Kante | kein Aufwand | Spec-Stelle ohne Änderungskopplung; die fünf Fixrunden-Punkte blieben ohne Entscheidung, entgegen Verdikt §2 |
| **D — eine neue ADR schärft beide Kennungen und nennt die Herkunft je Teil** | Richtung stimmt, nichts Immutables berührt, ein Muster für die Folge-Slices | ein weiteres Dokument |

**Frage 2 — der Umzug gegen `vcs`.**

| Option | Pro | Contra |
|---|---|---|
| (a) Teil-Ranges um den Umzugs-Commit, Beleg `R100` am Commit | ohne Werkzeug-Änderung, sofort tragfähig, dieselbe Mechanik wie der Form-Commit | eine Ausnahme mehr in der Verifier-Range; sie braucht einen Beleg (`umzug`) |
| (b) `vcs`-Konfiguration ändern | ein Lauf über die ganze Range | d-check v0.82.0 hat keinen Schlüssel für Umzüge (`--print-config`: `paths`, `immutable-when`, `exclude-sections`, `status-line`, `head-allow`). Ein engerer `paths`-Glob hilft nicht, denn der alte Pfad liegt in jedem Glob, der die aktiven Einträge schützt (*hergeleitet*, nicht gefahren). Eine Gate-Änderung ohne Gewinn |
| (c) Change Request an d-check, bis dahin (a) | Werkzeug lernt den Fall | Arbeit außerhalb des Repos für einen Vorgang, der selten anfällt; (a) trägt allein |
| nichts tun | — | `make doc-immutable` bleibt bei jeder Auflösung rot, oder die Range wird still ausgelassen |

Gewählt: (a). Sie trägt allein. Der Change Request (c) bleibt freiwillig und ist
keine Folgepflicht. Bietet d-check einen Umzug nach `done/` an, greift der
Re-Evaluierungs-Trigger.

## Konsequenzen

- Positiv: `SPEC-038` und `SPEC-039` haben eine Kante mit Herkunft je Teil. Die
  zehn Abweichungen des Werkzeugs vom historischen Block stehen auf einer
  Entscheidung. Die Folge-Slices der Festlegungen haben ein Muster: eine ADR je
  Slice mit `Schärft:` auf die neuen Kennungen.
- Positiv: Ein MR-Eintrag lässt sich nach Baseline-Regel auflösen, ohne dass
  `make doc-immutable` rot bleibt oder die Range still schrumpft.
- Negativ: Die Verifier-Range kennt jetzt zwei Arten ausgenommener Commits
  (Form-Commit, Umzugs-Commit). Jede trägt ihren eigenen Beleg.
- Folgepflicht (Implementer, im Slice, vor der Closure; Wortlaut im
  Architect-Bericht an den Implementer):
  - `spec/pflichtenheft.md`: die Wortlaut-Berichtigung aus Entscheidung 2 und ein
    Satz in `SPEC-039` zum Umzugs-Commit und zur Teilung an mehreren Commits.
  - `.claude/agents/verifier.md` und `.claude/agents/implementer.md`: der
    Umzugs-Commit neben dem Form-Commit.
  - `harness/targets/pin-stale.md`: der Umzug im Bump-Ablauf, wenn ein
    Adaptions-Eintrag einen Nachfolger bekommt.
  - `harness/conventions.md` §Adaptions-Block: die Regel zum Auflösen.
  - Der Slice-Plan: der Beleg am Umzugs-Commit `8e00e831` und die Übergabe des
    `formnorm`-`cmp` an `slice-spec-festlegungen-pruefer-hooks`.

## Fitness Function (falls maschinell prüfbar)

Gemessen an einem Klon im Scratchpad des Zugs (`architect-spec7/`), Stand
`54b0d398`, GNU bash, GNU Awk, d-check v0.82.0 über `make doc-immutable`. Die
Probe-Commits gibt es nur im Klon.

| Tooling | Regel | Make-Target |
|---|---|---|
| d-check `vcs` | **Gemessen** an den Teil-Ranges: `fcff30ec..54b0d398` (mit `M` = `8e00e831`) `1 Befund(e)`, `core-drift-vcs` an `MR-001`, Exit 2; `8e89831d..8e00e831` (nur `M`) ebenso Exit 2; `fcff30ec..8e89831d` (= `B..M~1`) `0 Befund(e)`, Exit 0; `8e00e831..54b0d398` (= `M..H`) `0 Befund(e)`, Exit 0 | `make doc-immutable RANGE=…` |
| d-check `vcs` | **Gemessen:** Der Schutz bleibt nach dem Umzug. Ein Probe-Commit, der im Core von `done/MR-001-…` ein Wort ändert, ergibt über `8e00e831..<Probe>` `core-drift-vcs`, Exit 2. Ein Probe-Commit mit zwei reinen Umzügen (`MR-002`, `MR-003`) ergibt über `54b0d398..<Probe>` Exit 2. Die Teilung um `M` ist also nötig, nicht nur zulässig | `make doc-immutable RANGE=…` |
| `umzug` (Entscheidung 3) | **Gemessen** am echten Bestand: `8e00e831` Exit 0 (`R100 harness/conventions/MR-001-… → harness/conventions/done/MR-001-…`), `dc04087e` und `8e89831d` Exit 1, eine nicht auflösende Kennung Exit 2. **Mutationen**, je ein Probe-Commit (Instanz: die Shell-Funktion gegen den Commit, nicht ein Tabellentest): Umzug mit angehängter Zeile (`R` unter 100) Exit 1, Umzug mit neuem Dateinamen Exit 1, Umzug in ein anderes Verzeichnis (`conventions/alt/`) Exit 1, zwei reine Umzüge Exit 0. Je Zweig der Bedingung ist eine Stelle erprobt. Ein Merge-Commit (Parent-Prüfung) ist *hergeleitet*, nicht gefahren | — (Befehlsform) |
| `make zitat-vergleich` | **Gemessen:** `make test-zitat-vergleich` Exit 0, „297 Fälle bestanden (je Runde 99 …)“. An Probe-Dateien im Klon: `L3` ergibt `keine Einheit, Exit 2`, `L3-3` `cmp 0`; `## Eins ##` löst als `#eins` auf (`cmp 0`); gestapelte `id`s vor `## Zwei` messen den Abschnittskörper von `#zwei` (`#a` gegen `#zwei` `cmp 0`). Eine `id`-Zeile ohne Inhalt in Zeile 5, Absatz ab Zeile 7, Heading in Zeile 10: `#x` gegen `L7-9` ergibt `cmp 0`, gegen `L5-9` (ab der `id`-Zeile) `cmp 1`. Der Block beginnt also an der nächsten Zeile mit Inhalt (Wortlaut-Berichtigung) | `make test-zitat-vergleich`, `make zitat-vergleich` |
| d-check `matrix`, `ids` | `Schärft:` dieser ADR auf `SPEC-038`/`SPEC-039` (Richtung ADR → Spec) ergibt keinen Befund. Gemessen mit `make docs-check` im Commit dieser ADR (Ergebnis im Architect-Bericht) | `make docs-check` |

## Re-Evaluierungs-Trigger

- (a) d-check bietet eine Erlaubnis für einen Umzug in ein `done/`-Verzeichnis
  (Schlüssel unter `vcs:` in `--print-config`): dann die Teilung um `M` durch die
  Konfiguration ersetzen, per Folge-ADR.
- (b) Eine Festlegung in `SPEC-038` oder `SPEC-039` widerspricht einer der
  Quell-ADRs: dann Folge-ADR mit `Supersedes` auf die betroffene Entscheidung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-07 | Accepted | `slice-spec-festlegungen-harness-werkzeuge` (Architect-Zug) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
