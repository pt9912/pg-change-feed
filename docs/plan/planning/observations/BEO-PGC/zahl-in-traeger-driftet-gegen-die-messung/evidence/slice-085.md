# Beleg: slice-085

Vorgang: `slice-085` — die Naht in `replication/receive`.

Fund: **Derselbe Träger, zum dritten Mal.** `harness/sensors/db-adapter-coverage.md`
führt den Nenner seit der **`slice-084`-Fixrunde** ausdrücklich als
**Zustandsgröße** — und dieser Vorgang hat ihn **nicht** mitgezogen: dort stand
`659` / `receive 155`, gemessen war `691` / `187` (Review `review-slice-085`
F-1). Der Implementer hat also **seine eigene, eine Slice alte Regel** nicht
angewandt.

**Der Ertrag ist die Form — und sie ist schärfer als die der Vorgänger:**

> Jede Zahl eines Doku-Trägers trägt ihren **Ursprung**; ist sie eine Messung,
> trägt sie zusätzlich den **Zeitpunkt** (Lauf/Stand), in dem sie gemessen wurde.

Das ist `ADR-0078`s „gemessen / übernommen / **abgeleitet**" plus die
**`wann`**-Hälfte. Und sie trennt die zwei beweglichen Sorten nach dem, **woran**
sie hängen:

| Sorte | Bindung | Bewegung |
|---|---|---|
| **Nenner** | Code-Stand | derselbe Stand misst denselben Nenner; ein Zug mit Produktionscode einen anderen |
| **gedeckte Zahl** | Lauf | wandert schon bei unverändertem Stand |

**Der Beweis, dass der Zeitpunkt nicht nachträglich zu holen ist:** der
Implementer musste den Ursprung von `112/155`, `132`, `130/178` erst per
`git log -S` rekonstruieren — alle aus **einem** Commit (`93a5cad`), dessen
Message daneben eine **andere** Zahl nennt (`610` neben `650`). Der Träger war
**bei seiner Geburt uneindeutig datiert**. Sein Satz dafür: *„der Schreiber
setzt den Zeitpunkt, der Leser kann ihn nicht erraten."*

**Und die Form ist noch nicht vollständig durchgesetzt:** die Verifikation hat
**gegen** sie gesucht und Fundstellen ohne Zeitpunkt gefunden
(`coverage-gate.md` §Grenze Punkt 1, zwei §Ausgabe-Abschnitte) — **V-1**.

Quelle: `docs/reviews/review-slice-085.md` (F-1) ·
`docs/reviews/verify-slice-085.md` (V-1, eigene Suche gegen die Form) ·
Implementer-Bericht der Fixrunde (`dd29d83`, `61bb5cf`) und der
Auftragserweiterung (`b83217c`) ·
`harness/sensors/db-adapter-coverage.md`, `harness/sensors/coverage-gate.md`,
`tools/harness/db-coverage.sh`.
