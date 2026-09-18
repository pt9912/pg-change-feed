# ADR-0095: `review`-Klasse von der `matrix.status`-Prüfung ausgenommen (ergänzt ADR-0094)

**Status:** Accepted

**Datum:** 2026-09-18

**Autor:** Architect (pt9912; anderer Kontext als die Implementer-Fixrunde,
die den ADR-0094-Diff in `.d-check.yml` einträgt — `modul-08-agentenrollen.md`
§Rollen-Regeln)

**Bezug:** [`ADR-0094`](0094-review-matrixklasse-kennung-statt-adresse.md)
(die `review`-Matrixklasse und die Regel `adr → review: false`, die diese
ADR unverändert lässt) · [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
(Record-Einfrierung ist zeitlich, nicht Status-basiert — dieselbe Begründung,
hier auf `matrix.status` statt auf `versions.pin-pattern` angewandt) ·
`AGENTS.md` §3.5 (ADR-Immutabilität) · `AGENTS.md` §3.6 (Gate-Lockerung
braucht ADR).

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0094`; ändert
`.d-check.yml`, nicht `spec/`).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0094` aktiviert die neue `matrix`-Klasse `review`
(`paths: ["docs/reviews/*.md"]`) und die Regel
`{from: adr, to: review, allow: false}`. Der dort gezeigte Diff wurde von
der nachfolgenden Implementer-Fixrunde unverändert in `.d-check.yml`
eingetragen — **aktuell uncommittet im Arbeitsbaum**, `git diff
.d-check.yml` zeigt exakt den in `ADR-0094` §Entscheidung abgedruckten
Diff, keine Abweichung.

Real gemessen (`make docs-check`, Docker-Digest
`ghcr.io/pt9912/d-check@sha256:18e9cd857f…`) bricht dieser Diff jedoch an
einer Stelle, die `ADR-0094` nicht betrachtet hat: **58 neue
`matrix-inactive`-Befunde**, alle unter `docs/reviews/**`. Ursache,
strukturell bestätigt: `.d-check.yml`s `matrix.status`-Block
(`forbidden: [superseded, deprecated]`) gilt **für alle** in `matrix.classes`
eingetragenen Klassen gemeinsam — er ist kein Unter-Schlüssel einzelner
Klassen. Vor `ADR-0094` war `docs/reviews/**` in keiner `matrix`-Klasse
Mitglied und damit von `matrix.status` strukturell unberührt. Die neue
`review`-Klasse macht jeden Review-Report und jedes Architect-Verdikt
erstmals zu einem `matrix`-Mitglied — und damit erstmals zu einer
**referenzierenden** Datei, deren ausgehende `ADR-\d{4}`-Token gegen
`matrix.status.forbidden` geprüft werden. 58 Bestandsdateien verweisen
legitim auf inzwischen `superseded` ADRs (`ADR-0018`, `ADR-0036`,
`ADR-0038`, `ADR-0039`) — sie frieren den zum Laufzeitpunkt geltenden Stand
ein, exakt das Prinzip, das `ADR-0073`/`ADR-0072` für „Records" bereits als
„Einfrierung ist zeitlich, nicht Status-basiert" etabliert haben
(`.d-check.yml`s `versions`-Modul trägt für `docs/reviews/**` bereits
dieselbe Begründung wörtlich: *„docs/reviews/\*\* sind Lauf-Belege
(Modul 10): ein Review-/Verify-Report zitiert den zum Laufzeitpunkt
adoptierten Baseline-Stand dauerhaft — er wird nie auf einen späteren
Stand nachgezogen (Record-Einfrierung ist zeitlich, nicht Status-basiert,
ADR-0073)"*).

Real gemessen, dass diese Begründung 1:1 überträgt: Nach Ergänzung von
`docs/reviews/*.md` in `matrix.exempt-paths` (demselben Schlüssel, der für
die beiden Grandfather-ADRs `docs/plan/adr/0039-*.md` und
`docs/plan/adr/0041-*.md` bereits existiert, Präzedenz-Commit `46d2fc6`)
sinkt die Befund-Zahl exakt auf 0 — bei unveränderter `review`-Klasse und
unveränderter `adr → review`-Regel. Die Trennung ist strukturell real, nicht
nur behauptet: `matrix.exempt-paths` und `matrix.rules` sind zwei
unabhängige Schlüssel derselben `matrix`-Sektion; eine gezielte
Gegenprobe — ein Testeintrag `ADR-0094` → ein bestehender Review-Report
als Live-Link (testweise eingefügt, sofort wieder zurückgenommen, kein
Bestandteil dieser ADR) — liefert weiterhin real `matrix-forbidden: Referenz
adr → review ist nicht erlaubt`, **trotz** der Exempt-Paths-Ergänzung. Die
Ergänzung schwächt die von `ADR-0094` eingeführte Regel nicht; sie exemptet
`docs/reviews/*.md` ausschließlich als **referenzierende** Datei von der
`status`-Prüfung — dem einzigen Prüfteil, der vor `ADR-0094` nicht griff.

`ADR-0094` selbst hat den `status`-Block an keiner Stelle erwähnt oder
geprüft — weder in §Kontext noch in §Entscheidung noch in der
Fitness-Function-Zeile (die dort formulierte Zusage lautet exakt „0
`matrix-forbidden`-Befunde für `adr → review`", nicht „0 Befunde
insgesamt", und diese engere Zusage ist bereits mit dem reinen
`ADR-0094`-Diff erfüllt — die 58 Befunde tragen den Code
`matrix-inactive`, einen anderen Befund-Code). Diese ADR schließt eine
Lücke, die `ADR-0094` nicht betrachtet hat; sie ändert an `ADR-0094`s
eigener Entscheidung — Klasse `review`, Regel `adr → review: false` — nichts.

## Entscheidung

Wir ergänzen `matrix.exempt-paths` in `.d-check.yml` um `docs/reviews/*.md`:

```yaml
matrix:
  # ...
  exempt-paths: ["docs/plan/adr/0039-*.md", "docs/plan/adr/0041-*.md", "docs/reviews/*.md"]
```

Diese Ergänzung ist **unabhängig** von `ADR-0094`s `classes`/`rules`-Diff
(zwei getrennte Schlüssel derselben `matrix`-Sektion) und lässt ihn
unverändert. Sie exemptet `docs/reviews/*.md` ausschließlich von der
`matrix.status`-Prüfung, nicht von `matrix.rules` — die Regel
`{from: adr, to: review, allow: false}` bleibt scharf (real gegengeprüft,
§Kontext).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — `review`-Klasse ohne jede `status`-Anwendung, etwa über einen Klassen-lokalen `status: off`-Schalter | sauberer als eine Pfad-Ausnahme, adressiert die Ursache (Klassen-Zugehörigkeit) statt eines Pfad-Musters | technisch nicht verfügbar: `.d-check.yml`s `matrix.status` ist ein einziger, klassen-übergreifender Block (kein `status`-Unterschlüssel je `classes`-Eintrag, geprüft am bestehenden Bestand und an der Baseline-Vorlage `.harness/baseline/v6.9.0/templates/.d-check.yml`, die `status: {forbidden: […]}` ebenfalls nur auf Matrix-Ebene zeigt) — eine Klassen-lokale Abschaltung existiert im Werkzeug nicht |
| B — `superseded`/`deprecated` repo-weit aus `matrix.status.forbidden` streichen | löst die 58 Befunde ebenfalls, ein Zeilen-Diff | schwächt die Klasse für `adr`/`spec`/`slice` mit — genau dort bleibt die Prüfung sinnvoll (eine lebende ADR soll keine abgelöste ADR referenzieren); overreach gegenüber dem real beobachteten Problem, das ausschließlich `docs/reviews/**` betrifft |
| **C — `docs/reviews/*.md` in `matrix.exempt-paths` — gewählt** | exakt dieselbe Begründung und derselbe Mechanismus, den `.d-check.yml`s `versions`-Modul für denselben Pfad bereits trägt (Record-Einfrierung ist zeitlich, `ADR-0073`); real auf 0 Befunde gemessen; die `adr → review`-Regel bleibt real scharf (Gegenprobe) | ein weiterer Eintrag in einer bereits zwei Einträge tragenden Liste — Pfad-Muster statt Klassen-Eigenschaft, siehe A |

**Fazit:** C. A ist technisch nicht verfügbar (kein Klassen-lokaler
`status`-Schalter im Werkzeug); B schwächt die Klasse dort, wo sie weiterhin
trägt.

## Konsequenzen

- Positiv: `make docs-check` liefert nach dieser Ergänzung real 0 Befunde
  für den `ADR-0094`-Diff (58 → 0, gemessen) — `ADR-0094`s eigentliches Ziel
  (die `adr → review`-Regel) bleibt scharf und real gegengeprüft.
- Positiv: Konsistente Begründung mit dem bereits bestehenden
  `versions.exempt-paths`-Eintrag für denselben Pfad — keine neue
  Konzept-Klasse, dieselbe „Record-Einfrierung ist zeitlich"-Linie
  (`ADR-0073`).
- Negativ mit Grenze: `matrix.exempt-paths` exemptet `docs/reviews/*.md`
  vollständig von `matrix.status` — nicht nur von Referenzen auf die vier
  aktuell betroffenen ADRs. Ein künftiger Review-Report, der auf eine
  inzwischen `superseded`/`deprecated` ADR verweist, wird von `matrix.status`
  ebenfalls nicht mehr gemeldet. Das ist beabsichtigt (dieselbe
  Record-Einfrierungs-Logik gilt für jeden künftigen Report gleichermaßen),
  nicht nur für den aktuellen Bestand — anders als die beiden
  ADR-spezifischen Grandfather-Einträge (`0039-*.md`, `0041-*.md`), die auf
  genau eine unveränderliche Datei zielen.
- Folgepflicht: Die Implementer-Fixrunde trägt diese `.d-check.yml`-Zeile
  zusammen mit dem bereits vorbereiteten `ADR-0094`-Diff in einem Commit;
  `make gates` grün vor Commit/Closure.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` Modul `matrix` (`exempt-paths` inkl. `docs/reviews/*.md`) | `make docs-check` liefert mit dem vollständigen `.d-check.yml`-Diff (ADR-0094-Diff + diese Ergänzung) **0** Befunde — real gemessen (903 Dateien geprüft, 0 Befunde); ohne diese Ergänzung, nur mit dem reinen `ADR-0094`-Diff: real gemessen 58 `matrix-inactive`-Befunde; mit dieser Ergänzung UND einem injizierten `adr → review`-Testverstoß: real gemessen weiterhin 1 `matrix-forbidden`-Befund (Regel bleibt scharf) | `make docs-check` (in `make gates`) |

## Re-Evaluierungs-Trigger

Ein künftiger `d-check`-Stand führt einen Klassen-lokalen `status`-Schalter
ein (Option A wird technisch verfügbar) — dann Folge-ADR, die die
Pfad-Ausnahme durch die sauberere Klassen-Eigenschaft ersetzt. Andernfalls
permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — `docs/reviews/*.md` in `matrix.exempt-paths` ergänzt, Anlass: 58 real gemessene `matrix-inactive`-Befunde nach dem `ADR-0094`-Diff | Architect-Verdikt zur `matrix.status`-Ausnahme für die `review`-Klasse (2026-09-18) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0095` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für
Accepted-ADRs).
