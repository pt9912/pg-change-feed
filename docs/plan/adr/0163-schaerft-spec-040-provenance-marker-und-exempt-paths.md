# ADR-0163: Festlegung `SPEC-040` (`make docs-check`) mit Herkunft aus der ADR-Kette; Provenance-Marker und `matrix.exempt-paths` nach dem Werkzeug

**Status:** Accepted — Supersedes [`ADR-0094`](0094-review-matrixklasse-kennung-statt-adresse.md)
(nur der Absatz „Bewusst **kein** Provenance-Marker-Escape für `adr → review`“
in §Entscheidung; die Klasse `review`, die Regel `adr → review` und die übrigen
Absätze bleiben in Kraft) · Supersedes
[`ADR-0097`](0097-observation-matrixklasse-review-verboten.md) (nur der Satz
„Kein Provenance-Marker-Escape, aus demselben Grund wie in `ADR-0094`“ in
§Entscheidung; die Klasse `observation` und die Regel `observation → review`
bleiben in Kraft)

**Datum:** 2026-10-10

**Autor:** pt9912 (Architect-Rolle, Modul 8; Architect-Zug zu
`slice-spec-festlegungen-doku-gates`, anderer Kontext als der Planner- und der
Implementer-Lauf des Slice)

**Bezug:** [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) ·
[`ADR-0074`](0074-zitationsform-schwester-repo-hausform.md) ·
[`ADR-0075`](0075-hostpaths-reichweite-und-wortlaut.md) ·
[`ADR-0094`](0094-review-matrixklasse-kennung-statt-adresse.md) ·
[`ADR-0095`](0095-review-klasse-exempt-status-check.md) ·
[`ADR-0097`](0097-observation-matrixklasse-review-verboten.md) ·
[`ADR-0099`](0099-slice-welle-review-regel-zurueckgenommen.md) ·
[`ADR-0156`](0156-versions-gate-nimmt-done-records-aus.md) ·
[`ADR-0160`](0160-teil-range-leer-test-und-hostpaths-home-relativ.md) ·
[`ADR-0161`](0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md) ·
[`ADR-0162`](0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
(Muster der Kante) · [`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md)
(gemessen/übernommen/hergeleitet) · `AGENTS.md` §3.5, §3.12 · d-check v0.82.0
(Digest aus `d-check.mk`).

**Schärft:** [`SPEC-040`](../../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
— die Kante der ADR-Kette von `make docs-check` zu dieser Festlegung läuft über
diese ADR (die Quell-ADRs tragen `Schärft: —` und bleiben unverändert).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`spec/pflichtenheft.md` §7 trägt seit Commit `9e43f968` die Festlegung
`SPEC-040` für `make docs-check`: Fläche, Befund je Modul, Randformen,
`structure`, `hostpaths`, Ausgänge. Die Spec nennt keine ADR (`matrix` verbietet
Spec → ADR). Nach dem Muster von `ADR-0162` braucht die Festlegung eine ADR, die
mit `Schärft:` auf sie zeigt; die Quell-ADRs sind `Accepted`, ein Nachtrag dort
wäre ein Überschreiben (`AGENTS.md` §3.5).

Der Implementer hat vier Fragen gestellt (Slice-Plan §6). Zwei davon sind
Abweichungen zwischen Werkzeug und ADR-Wortlaut:

- **Marker.** Der Kommentar `<!-- d-check:status-provenance -->` hebt einen
  Token-Befund von `matrix` auf. `ADR-0094` schließt das für `adr → review` aus,
  `ADR-0097` für `observation → review`. Der Bestand nutzt den Marker genau
  dort: 488 Zeilen in 271 Dateien unter `docs/plan/planning/observations/`
  tragen den Marker und ein Token der Klasse `review`, dazu eine Zeile in
  `ADR-0099` (gemessen am Stand `ff391aa0`, Befehle in §Fitness Function,
  Zeile 1). Ohne den Marker meldet das Werkzeug diese Zeilen (gemessen an einer
  Datei, Zeile 3; für die übrigen *hergeleitet*). Die Evidence-Dateien des
  Registers sind ab Merge unveränderlich (Baseline-Regelwerk
  `modul-06-roadmap.md` §Das Beobachtungs-Register), eine Umschreibung auf die
  Kennung-Form ist dort nicht zulässig.
- **`matrix.exempt-paths`.** Das Werkzeug prüft eine Datei unter diesem
  Schlüssel als Quelle gar nicht, weder auf Ziele mit verbotenem Status noch
  gegen die Regeln. `ADR-0095` nennt für `docs/reviews/*.md` nur die
  Status-Prüfung, `ADR-0073` nennt den Schlüssel „`matrix.status.exempt-paths`“.
  Die Liste trägt vier Einträge: `ADR-0039`, `ADR-0041`, `docs/reviews/*.md`,
  `done/welle-3-results.md` (`.d-check.yml`, Schlüssel `matrix.exempt-paths`).

Die zwei übrigen Fragen sind die Form der Kante und zwei Sätze in
`Accepted`-ADRs, die einen anderen Ort der Modul-Semantik nennen
(`ADR-0072` §Entscheidung, Absatz „Träger und Anker“ im Entwurf von §3.11;
`ADR-0160` §Konsequenzen, Folgepflicht 2).

## Entscheidung

1. **Eine ADR für `SPEC-040`, nicht eine je Modul; `make doc-immutable` nicht
   hier.** Die Kante folgt dem Muster von `ADR-0162`: eine ADR je Festlegung
   des Slice, Herkunft je Teil in einer Tabelle. `docs-check` ist ein Gate mit
   einer Spec-Zeile; eine ADR je Modul gäbe acht Kanten auf dieselbe Kennung.
   `vcs:`/`make doc-immutable` gehört zu
   `slice-spec-festlegungen-commit-baseline-gates`; dessen Festlegung gibt es
   noch nicht, und eine Kante auf eine nicht vergebene Kennung zeigt ins Leere.
   Der Folge-Slice plant seine eigene Architect-ADR bereits (dort §2,
   Bedingung vor der Closure).

   | Teil von `SPEC-040` | Quelle |
   |---|---|
   | Fläche und Eingabe (`scan.roots`, `scan.ignore`, Liste `modules:`) | `ADR-0161` Entscheidung 7 (alle `.md` unter `scan.roots`, `.harness/**` nicht gelesen); `ADR-0072` §Config-Block (`hostpaths` in `modules:`, kein eigener Block); im Übrigen diese ADR |
   | Befund je Modul `links`, `anchors`, `ids`, `tracked`; `structure` samt Regelarten; Ausgänge und Exit-Codes | diese ADR — keine Quell-ADR trägt sie; die Festlegung gibt das Werkzeug mit der Eingabe `.d-check.yml` wieder (Messungen M0 bis M9, M17, M18, M1a/M1b des Implementers, Slice-Plan §3, *übernommen*) |
   | `matrix`: Klassen `review`, `observation`, Regeln `adr → review`, `observation → review`, keine Regel `slice`/`welle → review`, Token fängt Pfad ohne Link | `ADR-0094` (ohne den abgelösten Absatz), `ADR-0097` (ohne den abgelösten Satz), `ADR-0099` |
   | `matrix`: Status-Prüfung, `allow-supersede-lineage`, `exempt-paths` | `ADR-0095`; Entscheidung 3 dieser ADR |
   | `matrix`: Provenance-Marker | Entscheidung 2 dieser ADR |
   | `versions` samt `versions.exempt-paths` | `ADR-0156` Entscheidung 1 und 2, `ADR-0161` Entscheidung 2 und 3 (erster Teil: `links`/`anchors` ohne Ausnahme) |
   | `hostpaths` | `ADR-0072` Entscheidung 2 (kein Ventil, kein Zeilen-Marker), `ADR-0075` Entscheidung 1 (Fence kein Treffer), `ADR-0160` Entscheidung 3 (Formen) und 4 (Ventil ungenutzt) |
   | `d-check:ignore` | `ADR-0072` Entscheidung 2 (`ids`, `versions` ja, `hostpaths` nein); für `matrix` Entscheidung 2 dieser ADR |

   **Nicht Quelle** sind: `ADR-0074` (Zitierform eines Schwester-Repos — eine
   Regel der Doku, kein Prüfverhalten; sie bleibt bei `AGENTS.md` §3.11 und im
   Vertrag); `ADR-0072` Entscheidung 1, 3, 4 und 5 (Aktivierung, abgelöste
   Zitierform, einmalige Korrektur, Entwurf von §3.11 — Punkt 5 ist durch
   `ADR-0075` abgelöst); `ADR-0075` Entscheidung 2 und 3; `ADR-0094` letzter
   Absatz (Semantik einer Umformulierung bleibt Review-Prüfpflicht);
   `ADR-0156` Entscheidung 3 und 4 und `ADR-0161` Entscheidung 3 zweiter Teil
   (Form-Korrektur, ein Doku-Verfahren); `ADR-0161` Entscheidung 7 die Liste
   des Bumps (Träger `harness/targets/pin-stale.md`).

   Wie bei `ADR-0162` Entscheidung 2 gilt: Weicht das Werkzeug von `SPEC-040`
   ab, ist das Werkzeug falsch. Ändert sich die Festlegung, wird `SPEC-040`
   fortgeschrieben, und diese ADR bleibt ihre Kante.

2. **Provenance-Marker: Token ist Herkunft, Link ist Adresse.** Der Marker
   `<!-- d-check:status-provenance -->` hebt in seiner Zeile den
   **Token-Befund** von `matrix` auf; einen **Link-Befund** hebt er nicht auf.
   Das gilt für jede Regel, auch `adr → review` und `observation → review`
   (gemessen je an einer Stelle, Zeilen 4, 5 und 7; für alle Regeln
   *hergeleitet*). Ein Pfad unter `docs/reviews/` in Inline-Code mit dem Marker
   ist eine Herkunftsangabe, keine Adresse, und in einer ADR und im
   Beobachtungs-Register zulässig. Der Grund aus `ADR-0094` bleibt gültig und
   wird hier eingelöst, nicht verworfen: ein Review-Bericht hat keine
   dauerhafte Adresse, deshalb bleibt der **Link** verboten (gemessen: Link mit
   Marker meldet `matrix-forbidden`, Zeilen 4 und 5). Die Zeile muss ohne die
   Datei lesbar bleiben; die Kennung des Slice daneben ist die bevorzugte Form,
   nicht die vorgeschriebene. `<!-- d-check:ignore -->` hebt keinen Befund von
   `matrix` auf (Zeile 6).

   **Grenze: die Spec.** In `spec/` ist der Marker kein zulässiger Weg. Die
   Decken-Regel (kein Spec-Stratum nennt eine ADR oder einen Slice) hat keinen
   Provenance-Ausweg; das Werkzeug lässt ihn trotzdem zu (Zeile 7). Der Bestand
   trägt in `spec/` keinen Marker (Zeile 1). Wächter ist das Review jeder
   Spec-Änderung: der Marker ist ein wörtlicher Kommentar im Diff. Das ist ein
   akzeptiertes Negativ — die Spec wird nur in Slices mit Review geändert, und
   ein Sensor dafür wäre eine Gate-Änderung ohne bisheriges Auftreten.

3. **`matrix.exempt-paths` nimmt die Datei als Quelle ganz aus.** Eine Datei
   unter dem Schlüssel prüft `matrix` als Quelle weder auf Ziele mit
   verbotenem Status noch gegen die Regeln; als **Ziel** einer Regel bleibt sie
   Gegenstand der Regel (Zeilen 8 und 2). Für `ADR-0095` deckt der Wortlaut die
   Arbeit: keine Regel in `.d-check.yml` hat `review` als Quelle (gelesen,
   Schlüssel `matrix.rules`, fünf Regeln mit den Quellen `spec`, `adr`,
   `observation`), also hat die Ausnahme für `docs/reviews/*.md` über die
   Status-Prüfung hinaus kein Objekt; dasselbe gilt für
   `done/welle-3-results.md` (Klasse `welle`, keine Regel als Quelle). Für
   `ADR-0039` und `ADR-0041` ist die Ausnahme weiter, als `ADR-0073` sie nennt:
   sie nimmt die beiden auch von `adr → slice` und `adr → review` aus. Das wird
   hiermit dokumentiert und angenommen: beide sind `Accepted` und ändern sich
   nur per Zitat-Korrektur mit Review, und heute trägt keine der beiden ein
   Token oder einen Link der Klassen `slice` oder `review` (Zeile 9). Wer einen
   Pfad neu unter `matrix.exempt-paths` aufnimmt, nennt in seiner ADR, dass die
   Datei damit auch von den Regeln ausgenommen ist.

4. **Ort der Modul-Semantik: `SPEC-040`; keine weitere Folge-ADR.** Was ein
   Modul von `make docs-check` als Befund meldet und an welcher Randform es
   entscheidet, steht in `SPEC-040`; der Vertrag
   `harness/sensors/docs-check.md` führt Grenzen und Sperren aus seiner Sicht.
   Die zwei Sätze, die einen anderen Ort nennen, brauchen keine eigene
   Ablösung: der Satz in `ADR-0072` steht im Entwurf von §3.11 (§Entscheidung
   Punkt 5), den `ADR-0075` laut seiner Status-Zeile bereits ersetzt hat; der
   Satz in `ADR-0160` ist eine Folgepflicht — eine einmalige Anweisung an
   `slice-dcheck-v0-82-0`, die erledigt ist und keine stehende Regel trägt.
   Beide bleiben als Geschichte stehen.

## Verglichene Alternativen

**Frage 1 — Form der Kante.**

| Option | Pro | Contra |
|---|---|---|
| A — eine ADR je Modul | feine Herkunft | acht ADRs auf eine Kennung; die Tabelle in Entscheidung 1 trägt dieselbe Feinheit in einer Datei |
| B — eine ADR für `SPEC-040` und die künftige Festlegung von `doc-immutable` | ein Dokument weniger | `Schärft:` auf eine nicht vergebene Kennung; die Herkunft des Folge-Slice ist noch nicht gelesen |
| C — nichts tun | kein Aufwand | `SPEC-040` ohne Kante, Closure-Bedingung des Slice offen |
| **D — eine ADR für `SPEC-040`, Herkunft je Teil (gewählt)** | Muster von `ADR-0162`, nichts Immutables berührt | ein weiteres Dokument |

**Frage 2 — Marker gegen `ADR-0094`/`ADR-0097`.**

| Option | Pro | Contra |
|---|---|---|
| (a) ADR gilt, Bestand umschreiben | Wortlaut bleibt | 271 Evidence-Dateien sind ab Merge unveränderlich; eine `Accepted`-ADR (`ADR-0099`) wäre in-place zu ändern |
| (b) ADR gilt, Change Request an d-check (Marker je Regel abschaltbar) | Werkzeug fängt den Fall | der Bestand würde mit dem Werkzeug rot; Arbeit außerhalb des Repos ohne Gewinn, weil die Praxis legitim ist |
| (c) nichts tun, `SPEC-040` schweigt | kein Aufwand | Festlegung unvollständig, Werkzeug und ADR widersprechen sich still |
| **(d) Lockerung nachziehen: Token mit Marker ist Herkunft, Link bleibt verboten (gewählt)** | deckt Bestand und Werkzeug, hält den Grund von `ADR-0094` am Link fest | Spec-Ausweg bleibt ohne Sensor (akzeptiertes Negativ) |

**Frage 3 — `matrix.exempt-paths` gegen `ADR-0095`.**

| Option | Pro | Contra |
|---|---|---|
| (a) Werkzeug soll nur die Status-Prüfung auslassen (Change Request) | Wortlaut von `ADR-0073` wird wahr | d-check v0.82.0 kennt keinen solchen Schlüssel; ändert heute nichts |
| (b) `ADR-0039`/`ADR-0041` aus der Liste nehmen | keine Regel-Ausnahme | dann melden sie `matrix-inactive` auf ihre abgelösten Ziele, die sie unveränderlich tragen (*hergeleitet* aus dem Grund der Grandfather-Zeile in `.d-check.yml`, nicht gefahren) |
| (c) nichts tun, `SPEC-040` nennt nur die Status-Prüfung | kein Aufwand | Festlegung beschreibt das Werkzeug falsch (Zeile 8) |
| **(d) Wirkung festlegen und für die zwei ADRs annehmen (gewählt)** | Festlegung = Werkzeug; Ausnahme ohne Objekt heute (Zeile 9) | eine weitere, benannte Lücke |

**Frage 4 — überholte Sätze.**

| Option | Pro | Contra |
|---|---|---|
| (a) Folge-ADR mit `Supersedes` je Satz | formal vollständig | der Satz in `ADR-0072` ist schon abgelöst, der in `ADR-0160` war eine Einmal-Anweisung; zwei Dokumente ohne Wirkung |
| (b) Zitat-Korrektur in-place | kein neues Dokument | der Referent ändert sich, also keine Zitat-Korrektur (`AGENTS.md` §3.5) |
| **(c) Entscheidung 4: Ort festlegen, Sätze als Geschichte (gewählt)** | eine Zeile, kein Eingriff | — |

## Konsequenzen

- Positiv: `SPEC-040` hat eine Kante mit Herkunft je Teil. Werkzeug,
  Festlegung und Entscheidungen widersprechen sich beim Marker und bei
  `matrix.exempt-paths` nicht mehr.
- Negativ: der Marker in `spec/` und die Regel-Ausnahme für `ADR-0039`/`ADR-0041`
  haben keinen Sensor; Wächter ist das Review.
- Folgepflicht (Implementer, in `slice-spec-festlegungen-doku-gates`, vor der
  Closure):
  1. `spec/pflichtenheft.md`, Absatz „Zu `SPEC-040` — Randformen der
     Verweis-Module“: der Satz zu `matrix.exempt-paths` lautet künftig „Eine
     Datei unter `matrix.exempt-paths` prüft das Modul als Quelle nicht, weder
     auf Ziele mit verbotenem Status noch gegen die Regeln; als Ziel einer Regel
     bleibt sie Gegenstand der Regel.“ Nach dem Satz zu `d-check:ignore` steht
     neu: „Eine Zeile, die den Kommentar `<!-- d-check:status-provenance -->`
     trägt, ist für ein Token von `matrix` kein Treffer, für jede Regel; ein
     Link in derselben Zeile bleibt ein Befund. `<!-- d-check:ignore -->` wirkt
     auf `matrix` nicht.“
  2. Slice-Plan §3: Gegenprobe C15, C17 und B8 mit Ausgang S und Verweis auf
     diese ADR; Anschluss-Frage, Zeilen `d-check:ignore`,
     `d-check:status-provenance` und `matrix.exempt-paths` mit der Folge nach
     Entscheidung 2 und 3; §6 die Architect-Frage und die drei Befunde mit
     Ausgang *entfallen* (Grund: diese ADR); das Bezug-Feld um diese ADR
     ergänzt; der DoD-Punkt „Bedingung vor der Closure — `Schärft:`-Kante“ mit
     dem Commit dieser ADR als Beleg.
  3. Reviewer des Slice: prüft die zwei Sätze aus Punkt 1 gegen Entscheidung 2
     und 3 und die Herkunftstabelle in Entscheidung 1 gegen die Gegenprobe des
     Slice-Plans.

## Fitness Function (falls maschinell prüfbar)

Gemessen von diesem Architect-Zug an einem Klon von `ff391aa0` im Scratchpad,
d-check v0.82.0 über `docker run --rm --network none -v <Klon>:/repo:ro <Digest>`
(Digest aus `d-check.mk`), Basislauf `0 Befund(e)`, Exit 0. Jede Mutation ist
**eine** Zeile, angehängt an **eine** Datei (Instanz: der d-check-Lauf über den
Klon, kein Tabellentest), danach `git checkout -- .`. Die Ausweitung auf alle
Dateien einer Klasse ist *hergeleitet*.

| # | Tooling | Regel | Make-Target |
|---|---|---|---|
| 1 | `git grep` | **Gemessen** (Stand `ff391aa0`): Zeilen mit Marker und Token der Klasse `review` unter `docs/plan/planning/observations/` 488, Dateien 271; unter `docs/plan/adr/[0-9]*.md` 1 Zeile (`ADR-0099`); unter `spec/` 0. Befehle im Codeblock unter der Tabelle | — |
| 2 | d-check `matrix` | **Gemessen:** Token `docs/reviews/<Bericht>.md` in Inline-Code ohne Marker, angehängt an `ADR-0094` → `matrix-forbidden`, Exit 1 (das Ziel liegt unter `exempt-paths` und bleibt Gegenstand der Regel) | `make docs-check` |
| 3 | d-check `matrix` | **Gemessen:** Marker aus einer Evidence-Datei (`BEO-PGC/adapter-fehler-ausgang`) entfernt → 2 × `matrix-forbidden`, Exit 1 | `make docs-check` |
| 4 | d-check `matrix` | **Gemessen:** Link auf einen Bericht unter `docs/reviews/` mit Marker, angehängt an `ADR-0094` → `matrix-forbidden`, Exit 1; derselbe Link ohne Marker ebenso | `make docs-check` |
| 5 | d-check `matrix` | **Gemessen:** Link auf einen Bericht unter `docs/reviews/` mit Marker in einer Evidence-Datei → `matrix-forbidden`, Exit 1; Token mit Marker in `ADR-0094` → `0 Befund(e)`, Exit 0 | `make docs-check` |
| 6 | d-check `matrix` | **Gemessen:** Token wie Zeile 2 mit `<!-- d-check:ignore -->` → `matrix-forbidden`, Exit 1 | `make docs-check` |
| 7 | d-check `matrix` | **Gemessen:** Tokens einer ADR- und einer Slice-Kennung mit Marker, angehängt an `spec/architecture.md` → `0 Befund(e)`, Exit 0 (die Lücke aus Entscheidung 2) | `make docs-check` |
| 8 | d-check `matrix` | **Gemessen:** Link und, getrennt, Token auf einen Bericht unter `docs/reviews/`, angehängt an `ADR-0039` (unter `exempt-paths`) → je `0 Befund(e)`, Exit 0. Die Slice-Hälfte (Token mit Nummer in `ADR-0039`, Exit 0) ist die Messung M16 des Implementers, *übernommen* | `make docs-check` |
| 9 | `git grep` | **Gemessen** (Stand `ff391aa0`): Tokens der Klassen `slice` und `review` sowie Pfade nach `planning/` oder `reviews/` in `ADR-0039` und `ADR-0041` → 0 Treffer. Befehle im Codeblock unter der Tabelle | — |
| 10 | d-check `matrix`, `ids` | `Schärft:` dieser ADR auf `SPEC-040` ergibt keinen Befund — gemessen mit `make gates` nach dem Commit dieser ADR (Ergebnis im Architect-Bericht) | `make docs-check` |

Befehle zu Zeile 1 und 9 (Stand `ff391aa0`):

```text
git grep -h "status-provenance" -- 'docs/plan/planning/observations/**' | grep -cE 'docs/reviews/[A-Za-z0-9_.-]+\.md'
git grep -n "status-provenance" -- 'docs/plan/planning/observations/**' | grep -E 'docs/reviews/[A-Za-z0-9_.-]+\.md' | cut -d: -f1 | sort -u | wc -l
git grep -h "status-provenance" -- 'docs/plan/adr/[0-9]*.md' | grep -cE 'docs/reviews/[A-Za-z0-9_.-]+\.md'
git grep -c "status-provenance" -- spec/
git grep -nE 'slice-[0-9]{3}|docs/reviews/[A-Za-z0-9_.-]+\.md' -- 'docs/plan/adr/0039-*.md' 'docs/plan/adr/0041-*.md'
grep -n "planning/\|reviews/" docs/plan/adr/0039-*.md docs/plan/adr/0041-*.md
```

## Re-Evaluierungs-Trigger

- (a) d-check bietet einen Schlüssel, der `matrix.exempt-paths` auf die
  Status-Prüfung beschränkt, oder den Marker je Regel abschaltbar macht: dann
  per Folge-ADR die Regel-Ausnahme für `ADR-0039`/`ADR-0041` und den
  Marker-Ausweg in `spec/` schließen.
- (b) Ein Review findet den Marker in `spec/` oder einen neuen Eintrag unter
  `matrix.exempt-paths` ohne die Angabe aus Entscheidung 3: dann den Sensor aus
  (a) oder eine `structure`-Regel per Folge-ADR.
- (c) Eine Festlegung in `SPEC-040` widerspricht einer Quell-ADR: Folge-ADR mit
  `Supersedes` auf die betroffene Entscheidung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-10 | Accepted | `slice-spec-festlegungen-doku-gates` (Architect-Zug) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
