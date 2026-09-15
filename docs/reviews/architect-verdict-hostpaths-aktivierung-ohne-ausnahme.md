# Architect-Verdikt: `hostpaths` aktivieren — ohne Ausnahme, nur die Form ist offen

**Rolle:** Architect (Modul 8)

**Anlass:** Anordnung des Auftraggebers, das `hostpaths`-Modul in
`.d-check.yml` einzubauen, und die Nachschärfung: „Entscheidend ist, dass alle
*hostpaths* entfernt werden müssen" — **31 → 0**, Modul aktiv in `modules`,
**kein** `scope`-/`ignore`-/`exempt`-Block. Der Planner hat gemessen; die
Frage, in welcher **Form** die Korrektur legitim und belegbar ist, ist eine
Regel-Entscheidung und gehört in die Architect-Rolle
(`modul-08-agentenrollen.md` §Rollen-Regeln: „ADR-Änderung: Architect
schreibt").

**Rolleninhaber:** pt9912 (Architect-Zug; anderer Kontext als der Planner-Lauf,
der die 31 Befunde gemessen hat, und als die Implementer-/Reviewer-/Verifier-
Läufe der betroffenen Slices)

**Datum:** 2026-09-15

**Bezug:**
[`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
(dieses Zugs Entscheidung — die Aktivierung, ohne Ausnahme) ·
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
(dieses Zugs zweite Entscheidung — die §3.5-Klasse) ·
`AGENTS.md` §3.1, §3.5, §3.6, §3.7 · `.d-check.yml` (`modules`) ·
`d-check.mk` (`DCHECK_DIGEST`) · `harness/sensors/docs-check.md` ·
`docs/plan/adr/0051-cicd-pipeline-github-actions.md`,
`docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md` (die zwei
betroffenen `Accepted` ADRs) · gepinntes Modul-Image `pt9912/d-check` `v0.75.0`
(`internal/hexagon/core/rules/hostpaths.go`).

**Erzeugte Artefakte dieses Zugs:** die Entscheidungen `ADR-0072` und
`ADR-0073` + ihre Index-Zeilen; die Folgearbeit (Zielrolle, unten) ist als
Adresse benannt, nicht als Slice angelegt — Priorisierung und Schnitt führt
der Planner.

---

## Frage

Der Auftraggeber hat „ob" und „welche Klassen bleiben" beantwortet (31 → 0,
kein Ausnahmeblock). Offen ist allein die **Form**:

1. **Aktivieren?** — ja; mit welcher Lesart von `AGENTS.md` §3.6 (Verschärfung
   braucht keinen Träger nach dem Wortlaut) und mit welcher Konsequenz?
2. **Wie wird die Korrektur der 11 Stellen in den zwei `Accepted` ADRs
   legitimiert?** — §3.5 erlaubt nur `Supersedes`. „Ohne Ausnahme" verlangt
   einen Nachzug von §3.5 für eine enge Klasse.
3. **Wie wird ein Schwester-Repo host-pfad-frei zitiert?** — eine Form für
   alle 31 Stellen.
4. **Wie wird die Korrektur belegt?** — Commit-Kennung genügt, oder braucht es
   eine §Geschichte-Zeile in der betroffenen ADR?
5. **Wie werden die `done/`- und `docs/reviews/`-Korrekturen legitimiert?**
   — nicht durch §3.5 geschützt, aber durch die Beleg-Konvention.
6. **Welchen Träger hat die Zusage?** — eine Hard Rule, mit ihren Grenzen.

## Befund — eigene Messung dieses Zugs

Der Planner hat 31 Befunde gemeldet; dieser Zug hat **nachgemessen**, nicht
übernommen. Verfahren: das gepinnte Modul-Image (`pt9912/d-check` `v0.75.0`,
`--enable hostpaths --disable <alle übrigen Module>`) über
`scan.roots: ["."]`.

**Ergebnis: 31 Befunde `hostpath-forbidden` über 616 geprüfte Dateien — deckungsgleich mit dem Planner-Befund.** Verteilung (eigene Zählung):

| Dokumentklasse | Dateien | Befunde |
|---|---|---|
| `Accepted` ADRs | `0051-…md` (1), `0054-…md` (10) | **11** |
| Zeitdokumente in `done/` | `slice-049-…md` (5), `slice-050-…md` (4), `welle-14.md` (1), `welle-14-results.md` (1) | **11** |
| Lauf-Belege | `review-slice-036.md` (2), `review-slice-039.md` (2), `review-slice-049.md` (1), `review-slice-073.md` (2) | **7** |
| lebende Doku | `harness/sensors/coverage-gate.md` (2) | **2** |
| **Summe** | | **31** |

Zitate: `pt9912/d-check` (27), `pt9912/d-migrate` (2), `pt9912/ai-harness-init` (1),
`pt9912/d-check/.github/` (1). Wörtlich:

```text
0054:41  Real geprüftes Vorbild: /Development/d-check/Dockerfile (Stage `coverage`, Zeilen 69–93)
0051:51  Die Musterquelle (/Development/d-check/.github/, gelesen, nicht kopiert)
review-slice-036:30  /Development/d-migrate
0054:54  (/Development/KI/ai-harness-init/.golangci.yml, exclusions.rules ...)
```

**Nachgeprüft, nicht übernommen (drei Modul-Eigenschaften):**

1. **Fenced-Code ist frei.** Eigene Probe (`t.md`): ein host-lokaler Pfad in
   Prosa und in Inline-Code wird gemeldet; derselbe Pfad in einem
   Backtick-Fence **nicht**. Ein **eingerückter** Code-Block (vier Leerzeichen)
   wird dagegen gemeldet — nur der Fence ist ausgenommen. Ein
   Windows-Laufwerks- oder UNC-Pfad wird ebenfalls gemeldet (die Muster sind
   fest).
2. **Kein Opt-out-Marker.** Der Modul-Quelltext
   (`internal/hexagon/core/rules/hostpaths.go`) sagt wörtlich: „Fenced-Code-Blöcke
   sind ausgenommen … es gibt keinen Opt-out-Marker." Die Ventil-Übersicht der
   Modul-Doku bestätigt: `hostpaths` hat **kein** `exempt-paths`, und den
   Zeilen-Marker `d-check:ignore` kennen nur `codepaths`, `ids`, `versions`,
   `diagrams`. Für einen nicht korrigierbaren Befund gibt es **kein** Ventil.
3. **Nur `.md`.** `walkMarkdown` filtert auf `.md`; die zwei bekannten
   Nicht-Markdown-Stellen (`Makefile`, `tools/coverage-gate.sh` — je ein
   Skriptkommentar) werden **nicht** gemeldet.

Zusätzlich: `make docs-check` ist auf dem Ist-Stand **grün** (616 Dateien,
0 Befunde) — die Aktivierung fügt die 31 Befunde also **neu** hinzu.

## Verdikt 1 — Aktivieren: ja, mit Träger

**Aktivieren.** `hostpaths` kommt in die `modules:`-Liste. Der Config-Block ist
**eine Zeile**; ein `hostpaths:`-Abschnitt entsteht **nicht** (keine
Präfix-Änderung, kein Scope).

**Lesart von §3.6.** `AGENTS.md` §3.6 regelt im Wortlaut die **Lockerung**
(„Jede Schwellen-Senkung … ist ein ADR"). Die Aktivierung ist eine
**Verschärfung** — sie fügt eine Prüfung hinzu, senkt nichts. Nach dem Wortlaut
bräuchte sie also **kein** ADR. Dieses Repo liest auch die Verschärfung als
**Entscheidung mit Träger**, nicht als Handgriff: Sie prägt, welche Zitate
künftig zulässig sind, welchen Wortlaut Hard Rule §3.11 trägt, und sie ist
ohne den §3.5-Nachzug (Verdikt 2) nicht durchführbar. Das Dokument ist deshalb
eine **ADR**, keine Aktennotiz.

**Config-Block (wörtlich, zur Übernahme durch den Planner):**

```yaml
modules: [links, anchors, ids, matrix, versions, structure, hostpaths]
```

Kein `hostpaths:`-Knoten, kein `scope`, kein `ignore`, kein `exempt-paths`.
Die Einzel-Targets (`doc-immutable`, `doc-commits`, …) führen bereits
`--disable hostpaths` und sind nicht betroffen.

## Verdikt 2 — Die §3.5-Klasse (keine Ausnahmeliste)

**Kein Ausschluss.** Es gibt **keine** Klasse, die ausgenommen bleibt. „Ohne
Ausnahme" wird erreicht, indem `AGENTS.md` §3.5 für eine **eng umrissene
Klasse** nachzieht: die **Zitat-Korrektur** — eine Änderung, die ausschließlich
die **Form eines Verweises auf einen unveränderten Referenten** ändert
(host-lokale Pfade, Linkziele, Zeilen-/Bereichs-Lokatoren, gebrochene
Referenz-Form).

**Was immutabel bleibt** (sonst wäre die Klasse ein Einfallstor): bei einer ADR
§Entscheidung, §Konsequenzen der Aussage nach, §Verglichene Alternativen,
§Status, die `Supersedes`-Kette, Fitness-Function-Regeln, Re-Evaluierungs-
Trigger, `Datum`/`Autor`, die Aussage-Semantik von §Bezug/§Schärft; bei einem
Record die **Funde, Beobachtungen, Closure-Aussagen, DoD-Haken**. Kurzform:
*das Gerüst darf sich ändern, die Aussage nie.*

**Die `Supersedes`-Alternative trägt nicht.** Eine Folge-ADR mit `Supersedes`,
die die betroffenen Klauseln host-pfad-frei neu fasst, **entfernt die Pfade
nicht**: der alte, host-lokale Text bleibt in der abgelösten ADR stehen und
wird weiter gescannt — im Umfang „Entfernen" ist sie gleich „nichts tun".
Zudem ist **nichts Normatives falsch**; `Supersedes` ist das Werkzeug für
*unwahre* Klauseln (`ADR-0070`/`ADR-0071`), nicht für eine Zitat-Form.

**Beleg der Korrektur.** (a) Die Commit-Message nennt `ADR-0073`
(Traceability, `ADR-0045`). (b) Jede betroffene `Accepted` ADR erhält **eine**
§Geschichte-Zeile (Datum, „Zitat-Korrektur — host-lokale Pfade ersetzt
(`ADR-0073`)", Commit-Verweis). Die Zeile ist das **in-Format-Provenienzfeld**
der ADR und selbst eine Zitat-Korrektur; sie hält die Unterscheidung
„Zitat-Fix vs. Substanz-Edit" **im Dokument**, ohne `git` — die richtige
Leseordnung für die immutabile Entscheidungsschicht. Ein Commit allein genügt
**nicht**.

## Verdikt 3 — Zielform des Zitats (eine für alle 31)

Ein Schwester-Artefakt wird **besitzer-qualifiziert** zitiert:

- **Kanonische Form:** `pt9912/<repo>` optional mit dem Pfad **im** Repo —
  `pt9912/d-check/Dockerfile`, `pt9912/ai-harness-init/.golangci.yml`.
- **Als Link** (wenn anklickbar gewünscht): zeigt auf
  `https://github.com/pt9912/<repo>` bzw. `/blob/main/<pfad>` — die bereits
  geübte Form (`docs/user/benutzerhandbuch.md` zitiert `pt9912/d-migrate` so).
- **Zeilen-/Bereichs-Lokatoren** (`Zeilen 69–93`) werden durch den **stabilen
  benannten Anker** ersetzt (`Stage coverage`, `bench:`-Target) — der Referent
  bleibt derselbe.

Beispiel (Fence, weil die Regel sonst ihre eigene Aussage verletzte):

```text
vorher:  Real geprüftes Vorbild: /Development/d-check/Makefile Zeile 84 (bench:-Target)
nachher: Real geprüftes Vorbild: `pt9912/d-check/Makefile` (`bench:`-Target)
```

Eine Form, keine 31 Einzelfälle. Die drei Repos sind alle
`github.com/pt9912/…` — die Form greift für alle.

## Verdikt 4 — Legitimation der `done/`- und `docs/reviews/`-Korrekturen

Diese 18 Stellen sind **nicht** durch §3.5 geschützt — sie sind
**Zeitdokumente** (Modul 5/6) bzw. **Lauf-Belege** (Modul 10), deren Freiheit
eine **Konvention** ist („Belege werden nicht umgeschrieben"), keine Hard Rule.
Der Träger ihrer Korrektur ist **dieselbe Klasse** (`ADR-0073`). Der Schnitt:

- **Erlaubt bei Records, verlangt bei ADRs:** Records brauchen **keine**
  §Geschichte-Zeile (sie haben keine) und **keinen** `Accepted`-Träger je
  Datei — der Commit genügt. Bei ADRs ist die Zeile **Pflicht** (Leseordnung
  ADR-vor-`git`).
- **Verboten bei beiden:** eine Änderung am **Record-Inhalt** (Fund,
  Beobachtung, Closure-Aussage, DoD-Haken) bzw. an der **ADR-Aussage**.

Hinweis: `docs/plan/planning/done/**` trägt in diesem Repo **keine** Archive
(kein `archiv.zip`); es gibt also keine zweite Kopie, die auseinanderlaufen
könnte. Sollte künftig eine Welle archiviert werden, ist die Korrektur **davor**
zu erledigen (Re-Evaluierungs-Trigger (d) in `ADR-0073`).

## Verdikt 5 — Grenzen der Zusage (benannt, nicht überzogen)

Der Satz „keine absoluten Pfade in diesem Repo" trägt **genau die gescannte
Markdown-Fläche**:

| Rand | Warum nicht gedeckt |
|---|---|
| **Fenced-Code** | Modul-Design — dort gehören bewusste Beispiel-Pfade hin; es gibt **keinen** Marker. Ein host-lokaler Pfad im Fence bleibt unbehelligt. |
| **Relative Pfade** | Das Modul prüft nur host-lokale **absolute** Pfade; ein relatives Ziel (`../…`), das auf einen Nachbarn oder aus dem Repo zeigt, wird **nicht** geprüft. Die Zusage ist **nicht** „kein Pfad verlässt das Repo". |
| **Windows-Muster fest** | Laufwerks-/UNC-Muster sind nicht konfigurierbar; sie sind immer an. |
| **Nicht-Markdown** | `Makefile`, `tools/**`, `harness/mk/**` liest der Scan nicht — die zwei bekannten Skriptkommentare werden in derselben Folgearbeit mitkorrigiert, aber **kein Gate hält sie**. |
| **`scan.ignore`** | `.harness/**` (vendored Baseline/Vorlagen) und `**/*.template.md` sind ausgenommen — der Scanner-Ausschnitt, nicht das Repo. |

Die **Gegenrichtung**: Die Zusage sagt nichts über *relative* Verweise; ein
relativer Pfad auf ein Schwester-Verzeichnis wäre grün und trotzdem ein
Repo-Escape. Das ist eine benannte Grenze, kein Defekt.

## Verdikt 6 — Keine Priorität/Rangordnung nötig (die gestellte Frage)

Der Konflikt zwischen §3.5 und der neuen Zusage ist **Umfang, nicht Rang**:
beide überschneiden sich allein in der Frage, was „inhaltlich" umfasst. Die
Lösung ist die **klassenweise Eingrenzung** von §3.5 (`ADR-0073`), nicht eine
Priorität. Eine Rangordnung zwischen Hard Rules wäre Seniorität in
Tabellenform — von `modul-08-agentenrollen.md` §Konflikt-Pfad ausdrücklich als
Auflösungsmittel **verworfen**. `ADR-0051` und `ADR-0054` bleiben in **jeder**
ihrer Entscheidungen unberührt; geändert wird an ihnen ausschließlich die
Zitat-Form.

## Verdikt 7 — Hard Rule §3.11 (Entwurf, nicht eingebaut)

Der Wortlaut steht in [`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
§Entwurf der Hard Rule §3.11. Er trägt: **Aussage** in einem Satz, **Grenzen**
(Fenced frei, Windows fest, relative ungeprüft, `tools/**`/`Makefile`/`harness/mk/**`
ungescannt, `scan.ignore`), ein **Falsch/Richtig**-Paar mit **echten** Beispielen
aus den 31 Befunden (im Fence, weil die Regel sonst sich selbst verletzte),
**Begründung** (warum Regel, nicht bloß Sensor) und **Träger/Anker**
(Rang-Zeiger auf das Modul und auf `harness/sensors/docs-check.md`, damit keine
zweite Quelle für dieselbe Aussage entsteht). Nächste freie Nummer am
Artefakt geprüft: `AGENTS.md` §3 endet bei §3.10 → **§3.11**.

## Folgearbeit mit Zielrolle

| Vorgang | Zielrolle | Art |
|---|---|---|
| `.d-check.yml`: `hostpaths` in `modules` | Implementer | ein Token, keine Ausnahme |
| 31 Zitat-Korrekturen (2 lebend · 7 Belege · 11 `done/` · 11 `Accepted` ADRs) + 2 Skriptkommentare | Implementer | Doku, kein Produkt-Code |
| je betroffener ADR **eine** §Geschichte-Zeile | Implementer | Belegform (`ADR-0073`) |
| `AGENTS.md` §3.5 (Wortlaut) + §3.11 (Wortlaut aus `ADR-0072`) | Implementer | Hard Rules |
| `docs/plan/adr/README.md` Kopf-Satz („Zitat-Korrektur ausgenommen") | Implementer | Index |
| `harness/sensors/docs-check.md` §Vertrag (`hostpath-forbidden`) + §Grenze | Implementer | Sensor-Doc |
| `harness/README.md` §Sensors-Zeile (`docs-check`) | Implementer | Harness-Einstieg |
| Priorisierung/Schnitt der Folgearbeit | Planner | Slice |

## Was dieser Zug geändert hat — und was nicht

**Geändert:**

- `docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md` (neu — die Aktivierung)
- `docs/plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md` (neu — die §3.5-Klasse)
- `docs/plan/adr/README.md` (zwei neue Zeilen)
- diese Verdikt-Datei

**Nicht geändert:** `.d-check.yml` (kein Modul aktiviert — Folgearbeit),
`AGENTS.md` (§3.5/§3.11 nicht eingebaut — Folgearbeit), die 31 Dokumente
(nicht korrigiert), `docs/plan/adr/0051-…`/`0054-…` (unangetastet — sie
bekommen erst in der Folgearbeit ihre §Geschichte-Zeile), jeder Slice-Plan,
`spec/**`, `internal/**`. Die Umsetzung ist Folgearbeit, nicht Teil dieses Zugs.
