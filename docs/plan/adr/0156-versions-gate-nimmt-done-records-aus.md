# ADR-0156: Das `versions`-Modul nimmt die Records unter `docs/plan/planning/done/**` aus

**Status:** Accepted

**Datum:** 2026-10-06

**Autor:** pt9912 (Architect-Rolle, Modul 8; der Weg ist vom Auftraggeber
gewählt, diese ADR trägt seine Begründung und Grenze)

**Bezug:** [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
(Record-Einfrierung ist zeitlich, Zitat-Korrektur als Klasse — bleibt
unverändert in Kraft, kein `Supersedes`) ·
[`ADR-0095`](0095-review-klasse-exempt-status-check.md) (dieselbe Begründung
für `docs/reviews/*.md` in `matrix.exempt-paths`) ·
[`ADR-0051`](0051-cicd-pipeline-github-actions.md) Entscheidung 7 (Pin-Inventar,
Achse P8 Kurs-Baseline) · [`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md)
· `AGENTS.md` §3.5, §3.6, §3.12 · `.d-check.yml` Block `versions:` ·
`harness/targets/pin-stale.md` §`make pin-stale-baseline` Schritt 3 · Slice
`slice-harness-baseline-v6-14-1` (Anlass).

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0073` und `ADR-0095`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Das d-check-Modul `versions` meldet jeden Verweis in die vendored Baseline,
dessen Versions-Segment (`pin-pattern: '\.harness/baseline/(v\d+\.\d+\.\d+)/'`)
vom adoptierten Stand (`harness/conventions.md` §Baseline) abweicht, als
`version-stale`. Ausgenommen sind heute `harness/conventions/done/**` und
`docs/reviews/**`. Für Letzteres steht im Kommentar am Block die Begründung
„Record-Einfrierung ist zeitlich, nicht Status-basiert, `ADR-0073`“. Für die
Records unter `docs/plan/planning/done/**` fehlt die Ausnahme, obwohl für sie
dieselbe Einfrierung gilt (`ADR-0073` Entscheidung 4) und der Bump-Ablauf sie
ausdrücklich nicht rückwirkend umschreibt (`harness/targets/pin-stale.md`,
Schritt 3).

**Gemessene Fälle** (Quelle: Plan von `slice-harness-baseline-v6-14-1`,
Abschnitt „Verweise, Pins und Records — Entscheidung je Zeile“, und eigene
Messung unten):

- Jeder Baseline-Bump färbt `versions` an `done/`-Plänen rot, **ohne dass sich
  die Datei ändert**; der Befund entsteht allein durch den neuen Stand in
  `harness/conventions.md`. Beantwortet wurde das bisher mit Zitat-Korrekturen
  je Zeile: `996e6231` (Bump auf v6.14.0, zwei `done/`-Pläne, drei Zeilen) und
  `54c6e632` (Bump auf v6.14.1, vier `done/`-Pläne, fünf Zeilen; gelesen mit
  `git show --stat`).
- Zwei Zeilen lassen sich **nicht** korrigieren, ohne die Aussage des Records
  zu ändern. Beide stehen in `done/slice-baseline-6-14-0-dokumente-nachziehen.md`:
  - **Zeile 277** ist eine `suchlauf`-Messzeile. Das Pfad-Argument zeigt am
    Commit `d79b7ebd` in das Tag-Verzeichnis v6.14.0. Mit dem v6.14.1-Pfad
    druckt dieselbe Messung nichts (Exit 1) statt 1, eine Korrektur würde also
    das Soll der Zeile fälschen (*übernommen* aus dem Slice-Plan).
  - **Zeile 337** zitiert das Report-Gerüst `review-report.template.md`. Der
    Referent hat sich zwischen v6.14.0 und v6.14.1 geändert (`cmp` = 1,
    *übernommen* aus dem Slice-Plan). `ADR-0073` verlangt einen unveränderten
    Referenten, die Zitat-Korrektur ist hier also nicht zulässig.
- **Gemessen** (d-check `sha256:b4b8756b…`, Klon von `aa0a44bc` im Scratchpad):
  Ohne Änderung liefert der Lauf `2 Befund(e)`, beide `version-stale`, Zeilen
  277 und 337.

Die Lage ist dieselbe wie bei `docs/reviews/**`. Ein Record nennt den Stand,
der zum Laufzeitpunkt adoptiert war, und dieser Satz wird durch den Zeitablauf
nicht falsch. Ein `version-stale` in einem Record ist deshalb kein Defekt des
Records, sondern eine Folge des Zeitablaufs.

**Abgrenzung zum Kommentar am `matrix`-Block.** Dort steht ausdrücklich, dass
`docs/plan/planning/done/**` **kein** Freibrief für `matrix` ist. Das bleibt
richtig und wird hier nicht berührt. Ein `matrix`-Befund in einem `done/`-Plan
entsteht durch den Text der Datei, ein `version-stale` durch einen Bump
außerhalb der Datei. Die Ausnahme gilt deshalb nur für das Modul, dessen Befund
ohne Dateiänderung entsteht.

## Entscheidung

1. **`versions.exempt-paths` nimmt `docs/plan/planning/done/**` auf**, mit der
   Begründung des `docs/reviews/**`-Eintrags. Sonst ändert sich am Modul nichts:
   `pin-pattern` und `current-from` bleiben, und `open/`, `next/`,
   `in-progress/` sowie die flachen Welle-Dateien bleiben geprüft.
2. **`links` und `anchors` bleiben ohne Ausnahme.** Gemessen: Nach dem
   Entfernen von v6.14.0 meldet `links` im ganzen Baum genau einen Befund
   (`target-missing`). Er betrifft den Architect-Verdikt-Report zur
   verfallenden Aufschub-Adresse unter `docs/reviews/`, Zeile 114, mit einem
   Link auf die Slice-Vorlage; der Pfad steht in der Entscheidungstabelle des
   Slice-Plans. Unter `done/` steht kein Markdown-Link mehr in ein altes
   Tag-Verzeichnis: `git grep "baseline/v6\.14\.0/"` trifft dort nur die Zeilen
   277 und 337, beide nicht in Linkform. Ältere Versionen trifft die Suche in
   `done/` gar nicht. Eine Ausnahme für `links`/`anchors` hätte also nichts zu
   tun und würde echte Linkbrüche verdecken (Fitness Function, M3–M5).
3. **Ein Link in ein altes Tag-Verzeichnis, der nach dem Löschen bricht**, wird
   per Zitat-Korrektur der **Form** nach `ADR-0073` behandelt, in `done/` wie
   in `docs/reviews/`: aus dem Markdown-Link wird Inline-Code, der Pfad samt
   alter Version bleibt **unverändert**. Dieser Weg ist hiermit bestätigt,
   nicht nur erlaubt. Er trägt auch bei geändertem Referenten, denn der
   Referent ist weiter dieselbe Datei des alten Tags (in der Git-Historie
   auflösbar). Es ändert sich nur die Form der gebrochenen Referenz, die
   Aussage des Records bleibt wortgleich. Das Linkziel auf das neue Tag
   **umzuziehen** ist dagegen nur zulässig, wenn der Referent zwischen beiden
   Tags gleich ist (`cmp`). Für den Anlassfall Zeile 114 (`cmp` = 1) bleibt
   deshalb nur die Form-Korrektur. Gemessen: Mit ihr und der Ausnahme aus 1
   meldet der Klon nach dem Entfernen von v6.14.0 `0 Befund(e)` (V3).
4. **Die Zeilen 277 und 337 bleiben unverändert.** Mit der Ausnahme sind sie
   kein Befund mehr.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A: Records bei jedem Bump korrigieren (heutiger Weg) | `versions` bleibt über `done/` scharf; kein neuer Eintrag in einer Ausnahmeliste | An zwei Zeilen dieses Bumps **unmöglich**, ohne die Aussage zu ändern (277: Messung, 337: Referent geändert). Jeder Bump macht Arbeit an eingefrorenen Dateien, die nichts verbessert (`996e6231`, `54c6e632`). Widerspricht `harness/targets/pin-stale.md` Schritt 3 („nicht rückwirkend umgeschrieben“). |
| B: Ausnahme-Marker je Zeile | Punktgenau, die übrigen `done/`-Zeilen blieben geprüft | d-check kennt für `versions` keinen Zeilen-Marker. Zeile 277 steht im Fence, wo kein Kommentar-Marker möglich ist (*übernommen* aus dem Auftrag, nicht gegen die d-check-Doku geprüft). Jeder Bump legte neue Marker an, die Pflege wächst also wie bei A mit dem Bestand. |
| C: `docs/plan/planning/done/**` aus `matrix` oder dem ganzen Scan ausnehmen | Ein Eintrag für alles | Verdeckt `matrix`-, `links`-, `anchors`- und `ids`-Befunde, deren Ursache im Text liegt. Widerspricht dem Kommentar am `matrix`-Block; für `links`/`anchors` gemessen ohne Objekt (Entscheidung 2) |
| D: nichts tun | Kein Eingriff | `make docs-check` bleibt rot; der Slice kann v6.14.0 nicht löschen, ohne die Records zu fälschen |
| **E: nur `versions`, nur `done/**` (gewählt)** | Dieselbe Begründung und derselbe Mechanismus wie bei `docs/reviews/**`, kein neues Konzept. Deckt sich mit dem dokumentierten Bump-Ablauf. Ein Eintrag, der beim nächsten Bump nicht wächst | Ein veralteter Pin in einem `done/`-Plan bleibt ungesehen (siehe Grenze unter §Konsequenzen). Die Ausnahme hängt an einem Pfad-Muster, nicht an der Eigenschaft „Record“. |

## Konsequenzen

- **Positiv:** Ein Baseline-Bump erzeugt in `done/` keine `version-stale` mehr.
  An Records ist nur noch dort Arbeit nötig, wo ein Link nach dem Löschen
  bricht (`links`, Entscheidung 3). Gate und Bump-Ablauf (`pin-stale.md`
  Schritt 3) sagen jetzt dasselbe.
- **Die Zitat-Korrekturen von `996e6231` und `54c6e632` bleiben stehen.** Sie
  genügen `ADR-0073`: der Referent ist je Zeile `cmp`-gleich, die Belegtabelle
  steht im Plan von `slice-harness-baseline-v6-14-1`. Records dürfen sie
  tragen. Ein Rückbau wäre eine weitere Änderung an eingefrorenen Dateien ohne
  Gewinn.
- **Negativ, also was das Gate verliert:**
  - *Gemessen:* Einen Pin mit falscher Version in Inline-Code einer
    `done/`-Datei meldet `versions` nicht mehr, weder auf oberster Ebene (M2)
    noch in einem Unterverzeichnis (`done/altbestand/`, M6). Ohne die Ausnahme
    meldet `versions` denselben Pin (M7).
  - *Hergeleitet:* Dasselbe gilt für einen Pin im Fence und für eine erfundene
    Version (z. B. `v6.41.0`), solange sie nicht in Linkform steht. Tatsächlich
    unentdeckt bleibt damit nur ein Pin, der **nach** dem `git mv` nach `done/`
    in die Datei kommt, etwa durch eine spätere Zitat-Korrektur mit Tippfehler.
    Ein Pin, der schon beim Schreiben falsch war, wird bereits in `in-progress/`
    gemeldet, denn dort bleibt `versions` scharf (M1 gemessen für
    `in-progress/`; für `open/`/`next/` hergeleitet, weil das Muster sie nicht
    nennt). Das ist ein akzeptiertes Negativ: Einen Pin in einem Record liest
    niemand mehr als Arbeitsanweisung.
  - *Gemessen, bleibt scharf:* Einen Link aus `done/` auf eine fehlende Datei
    eines Tags (M3, M4) und auf einen fehlenden Anker (M5) melden `links` bzw.
    `anchors` weiterhin.
- **Folgepflicht (Implementer, im Slice `slice-harness-baseline-v6-14-1`):**
  - `.d-check.yml`, Block `versions:`: `exempt-paths` wird zu
    `["harness/conventions/done/**", "docs/reviews/**", "docs/plan/planning/done/**"]`.
    Der Kommentar darüber wird um den `done/`-Fall erweitert, mit Anker auf
    diese ADR.
  - Die Report-Zeile 114 aus Entscheidung 2 bekommt die Zitat-Korrektur der
    Form nach Entscheidung 3; der Commit nennt `ADR-0073`.
  - In der Entscheidungstabelle des Slice-Plans werden für die drei offenen
    Zeilen Entscheidung und Kennung nachgetragen.

## Fitness Function (falls maschinell prüfbar)

**Instanz aller Zeilen:** ein d-check-Lauf
(`ghcr.io/pt9912/d-check@sha256:b4b8756b40d3dcd2670a3f83526cb5e5d727d1a850571f73be31edba248abb40`,
`docker run --rm --network none -v <Klon>:/repo:ro`) gegen einen Klon von
`aa0a44bc` im Scratchpad, dessen Konfiguration per Kopie um den Eintrag aus
Entscheidung 1 ergänzt ist. Jede Mutation ist an **einer** Stelle gefahren
(eine angehängte Zeile in der genannten Datei) und danach zurückgesetzt. Der
Lauf im Arbeitsbaum mit der echten `.d-check.yml` ist eine **Erwartung** an
den Implementer und noch nicht erprobt.

| Lauf | Zustand | Gesehen |
|---|---|---|
| V0 | Konfiguration unverändert | Exit 1, `2 Befund(e)`: `version-stale` Zeilen 277, 337 |
| V1 | + Ausnahme | Exit 0, `0 Befund(e)` |
| V2 | + Ausnahme, v6.14.0 per `git rm -r` entfernt | Exit 1, `1 Befund(e)`: `target-missing` an der Report-Zeile aus Entscheidung 2 |
| V3 | wie V2 + Form-Korrektur dieser Report-Zeile (Link → Inline-Code, Pfad unverändert) | Exit 0, `0 Befund(e)` |
| M1 | V3 + Inline-Pin `v6.13.0` in `in-progress/slice-harness-baseline-v6-14-1.md` | rot: 1 `version-stale` (die Ausnahme reicht nicht über `done/` hinaus) |
| M2 | V3 + Inline-Pin `v6.13.0` in `done/slice-baseline-6-14-0-dokumente-nachziehen.md` | grün: 0 Befunde (der Verlust) |
| M3 | V3 + Link auf eine fehlende Datei unter v6.14.1, in derselben `done/`-Datei | rot: 1 `target-missing` |
| M4 | V3 + Link in das entfernte Tag v6.13.0, in derselben `done/`-Datei | rot: 1 `target-missing` |
| M5 | V3 + Link mit fehlendem Anker in `regelwerk/modul-06-roadmap.md`, in derselben `done/`-Datei | rot: 1 `anchor-missing` |
| M6 | V3 + Inline-Pin `v6.13.0` in `done/altbestand/altbestand.md` | grün: 0 Befunde (`**` greift ins Unterverzeichnis) |
| M7 | wie M6, aber Konfiguration ohne Ausnahme | rot: 3 `version-stale` (277, 337, die mutierte Zeile) |

| Tooling | Regel | Make-Target |
|---|---|---|
| d-check `versions` | `exempt-paths` enthält `docs/plan/planning/done/**`; Pfade außerhalb bleiben geprüft (M1) | `make docs-check` (in `make gates`) |
| d-check `links`/`anchors` | keine Ausnahme für `done/`; Linkbruch und fehlender Anker bleiben rot (M3–M5) | `make docs-check` (in `make gates`) |

## Re-Evaluierungs-Trigger

Beobachtbare Trigger:

- **(a)** d-check bekommt für `versions` einen Zeilen- oder Fence-Marker oder
  eine Klasse „Record“. Dann wird der Pfad-Eintrag durch den Mechanismus
  ersetzt, der den Verlust (M2) enger fasst.
- **(b)** In einer `done/`-Datei taucht ein falscher Pin auf, der nach dem
  `git mv` hineinkam und Schaden anrichtete, etwa einen Leser fehlleitete. Dann
  ist die Annahme des akzeptierten Negativs widerlegt; es folgt eine
  Folge-ADR.
- **(c)** Die Ablage der Records ändert sich, etwa durch Archivierung nach
  `done/<welle>/archiv.zip` mit Stubs an anderem Ort. Dann wird das Muster
  geprüft.

Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-06 | Accepted. Anlass: zwei `version-stale` in `done/slice-baseline-6-14-0-dokumente-nachziehen.md` (Zeilen 277, 337), die keine Zitat-Korrektur zulassen; Weg vom Auftraggeber gewählt | `slice-harness-baseline-v6-14-1` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0156` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
