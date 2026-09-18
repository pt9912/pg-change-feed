# Architect-Verdikt — `matrix.status` und die neue `review`-Klasse (ADR-0094-Folgefrage)

**Anlass:** [`ADR-0094`](../plan/adr/0094-review-matrixklasse-kennung-statt-adresse.md)
(Accepted, committet `7809a9c`) beschließt die neue `matrix`-Klasse `review`
und die Regel `{from: adr, to: review, allow: false}`. Die Implementer-
Fixrunde hat den in `ADR-0094` §Entscheidung abgedruckten Diff vorbereitet
(zum Zeitpunkt dieses Verdikts uncommittet im Arbeitsbaum). Auftrag: prüfen,
ob dieser Diff eine unvorhergesehene Nebenwirkung hat, und falls ja,
entscheiden, ob die Korrektur bereits durch `ADR-0094` gedeckt ist oder eine
eigene Folge-Entscheidung braucht.

## Messung 1 — Nebenwirkung real reproduziert

`make docs-check` mit dem vorbereiteten `.d-check.yml`-Diff:

```
d-check: 903 Datei(en) geprüft, 58 Befund(e)
```

Alle 58 Befunde tragen den Code `matrix-inactive`, alle liegen unter
`docs/reviews/**` (Beispiele: der Architect-Review zu `slice-011`,
der Architect-Review zu `welle-1` — 20 Treffer allein dort —,
die Verifikationsberichte zu `slice-001` bis `slice-010`, Review zu `slice-008`,
Review zu `slice-010`).

Ohne den `.d-check.yml`-Diff (per `git stash push -- .d-check.yml`, dann
`make docs-check`, dann `git stash pop`), am unveränderten Bestand:

```
d-check: 903 Datei(en) geprüft, 0 Befund(e)
```

**Ursache, strukturell bestätigt:** `.d-check.yml`s `matrix.status`-Block
(`forbidden: [superseded, deprecated]`) ist ein einziger, klassen-
übergreifender Schlüssel — kein Unter-Schlüssel je `matrix.classes`-Eintrag.
Vor `ADR-0094` war `docs/reviews/**` kein `matrix`-Mitglied und daher von
`matrix.status` strukturell unberührt. Die neue `review`-Klasse macht jeden
Review-Report und jedes Architect-Verdikt erstmals zu einer
**referenzierenden** Datei, deren ausgehende `ADR-\d{4}`-Token gegen
`matrix.status.forbidden` geprüft werden — 58 Bestandsdateien verweisen
legitim auf inzwischen `superseded` ADRs (`ADR-0018`, `ADR-0036`,
`ADR-0038`, `ADR-0039`).

## Präzedenz — dieselbe Begründung existiert bereits im `versions`-Modul

`.d-check.yml`s `versions`-Modul trägt für exakt denselben Pfad
(`docs/reviews/**`) bereits einen `exempt-paths`-Eintrag mit der
Begründung: *„docs/reviews/\*\* sind Lauf-Belege (Modul 10): ein
Review-/Verify-Report zitiert den zum Laufzeitpunkt adoptierten
Baseline-Stand dauerhaft — er wird nie auf einen späteren Stand
nachgezogen (Record-Einfrierung ist zeitlich, nicht Status-basiert,
`ADR-0073`)"*. Diese Begründung überträgt 1:1 auf `matrix.status`: ein
Review-Report zitiert den zum Laufzeitpunkt aktiven ADR-Stand — dass die
zitierte ADR später `superseded` wird, macht das Zitat nicht falsch.

## Messung 2 — Struktur von `matrix.exempt-paths` verifiziert

`matrix.exempt-paths` ist entgegen der ursprünglichen Auftragsannahme
**kein** Unter-Schlüssel von `matrix.status` (`matrix.status.exempt-paths`
existiert nicht) — Einrückung geprüft (`sed -n | cat -A`): `status:` und
`exempt-paths:` stehen auf derselben Einrückungsebene, beide direkte Kinder
von `matrix:`. Der bestehende Eintrag
(`["docs/plan/adr/0039-*.md", "docs/plan/adr/0041-*.md"]`, eingeführt in
Commit `46d2fc6` für genau dieselbe Fallklasse — immutable Alt-ADRs mit
Fitness-Zeiger auf eine inzwischen abgelöste ADR) bestätigt das: er ist
`matrix.exempt-paths`, nicht `matrix.status.exempt-paths`.

Test: `docs/reviews/*.md` zu `matrix.exempt-paths` ergänzt, `make
docs-check` erneut gelaufen:

```
d-check: 903 Datei(en) geprüft, 0 Befund(e)
```

**58 → 0 bestätigt.**

Gegenprobe — schwächt die Ergänzung die `adr → review`-Regel aus
`ADR-0094`? Ein Testverstoß eingefügt (`docs/plan/adr/0094-*.md` verweist
live auf das Review zu `slice-001`), `make docs-check` mit der
`exempt-paths`-Ergänzung erneut gelaufen:

```
docs/plan/adr/0094-review-matrixklasse-kennung-statt-adresse.md:219	../../reviews/review-slice-001.md	matrix-forbidden	Referenz adr → review ist nicht erlaubt
d-check: 903 Datei(en) geprüft, 1 Befund(e)
```

Die Regel bleibt scharf. Testverstoß danach zurückgenommen
(`git diff docs/plan/adr/0094-*.md` zeigt keine Abweichung mehr). `matrix.
exempt-paths` und `matrix.rules` sind zwei unabhängige Mechanismen derselben
Sektion: `exempt-paths` exemptet eine Datei nur als **referenzierende**
Seite von `matrix.status`; die Richtungsregel `adr → review` prüft
unabhängig davon, welche Klasse referenziert wird.

## Entscheidung: (b) — eigene Folge-ADR nötig, nicht (a)

`ADR-0094` deckt diese Korrektur **nicht**. Ihr §Kontext, §Entscheidung und
die Fitness-Function-Zeile erwähnen `matrix.status` an keiner Stelle; die
dort formulierte Zusage lautet exakt „0 `matrix-forbidden`-Befunde für
`adr → review`" — eine engere Zusage, die bereits mit dem reinen
`ADR-0094`-Diff erfüllt ist (die 58 Befunde tragen den andersen Code
`matrix-inactive`). Die `exempt-paths`-Ergänzung ist echter neuer Inhalt
(eine Sachentscheidung: „`docs/reviews/*.md` ist von `matrix.status`
ausgenommen"), keine Form- oder Vollständigkeits-Frage der Zitierform von
`ADR-0094` — eine Zitat-Korrektur (`AGENTS.md` §3.5 Ausnahme, `ADR-0073`)
scheidet damit aus.

Da `ADR-0094`s eigene Klausel (`review`-Klasse, `adr → review`-Regel)
durch diese Ergänzung textlich und inhaltlich unverändert bleibt — real
gegengeprüft (Messung 2) —, ist ein `Supersedes ADR-0094` nicht nötig.
Dieselbe Unterscheidung trifft `ADR-0094` selbst gegenüber `ADR-0073`
(„Da `ADR-0073`s Klausel selbst unverändert bleibt, ist ein `Supersedes`
hier nicht nötig — der Rückverweis lebt stattdessen ausschließlich in
dieser ADR"). [`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md)
folgt demselben Muster: eine eigenständige, `ADR-0094` ergänzende ADR ohne
`Supersedes`-Beziehung, referenziert in ihrem §Bezug.

`AGENTS.md` §3.6 (Gate-Lockerung braucht ADR) greift zusätzlich: die
Ergänzung senkt die Prüfschärfe von `matrix.status` für einen Dateipfad —
eine notwendige Korrektur einer unvollständigen Vorentscheidung, keine
willkürliche Aufweichung, aber dennoch eine Schwellen-Senkung im Sinne der
Regel und damit ADR-pflichtig, unabhängig von der Immutabilitäts-Frage.

## Ausgang

- [`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md)
  geschrieben und `Accepted`: `matrix.exempt-paths` um `docs/reviews/*.md`
  ergänzt (dritter Eintrag neben den beiden Grandfather-ADR-Mustern),
  ADR-Index (`docs/plan/adr/README.md`) aktualisiert.
- `.d-check.yml` bleibt **unverändert gegenüber der vorgefundenen
  Implementer-Fixrunde** — nur der `ADR-0094`-Diff steht dort, meine
  Test-Ergänzung wurde nach der Messung zurückgenommen. Die nächste
  Implementer-Fixrunde trägt sowohl den `ADR-0094`- als auch den
  `ADR-0095`-Diff in einem Commit ein und führt `make gates` vor
  Commit/Closure.
- Committet in diesem Zug: nur `docs/plan/adr/0095-review-klasse-exempt-status-check.md`,
  `docs/plan/adr/README.md` (Index-Zeile) und dieses Verdikt-Dokument.
