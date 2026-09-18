# Beleg: slice-088

Vorgang: `slice-088` — Coverage-Tail „Reine Übersetzung" (Cluster B der
`welle-20`).

Fund: **Vier Planner-Kanten-Findings und zwei Register-Köpfe in einem Vorgang** —
alle dieselbe Bewegung: eine **ungezählte Übernahme**.

| Fund | Träger | stand dort | richtig |
|---|---|---|---|
| Review F-2 | Slice-Plan §1 | „**fünf** Funktionen ins falsche Paket" (nannte sechs Namen) | **acht** |
| Review F-3 | Slice-Plan §8 | „Zähler am Register **nachgezählt**", `db-gegenstand-…` **1×**, `zahl-in-traeger-…` **2×** | **2×** / **3×** (Schwelle) |
| Delta D-1 | Slice-Plan §1 | ein **Vorher/Nachher über den Text selbst** („die Zahl ‚fünf' in einer früheren Fassung …") | die Historie hält `git` (`AGENTS.md` §3.7) |
| Delta D-2 | Slice-Plan §2 LP1 | `JSONImage` als „**dabei** der erste Fall" der Stream-Übersetzung | es liegt in `postgresstorage/mapper`, LP3 |
| Delta D-3 | `db-gegenstand-…/state.md` | Zeile 1 `Zustand: offen (1×)`, Zeile 8 `Zähler: 2×` — **Selbstwiderspruch einer Datei** | 2× |
| Klassen-Check | `dod-begruendung-…/state.md` | Kopf `3× erreicht`, Belegliste **4** | 4× |

**Der F-3-Fall ist der schärfste:** §8 behauptete, die Zähler am Register
*nachgezählt* zu haben — die Sichtung stammt aber vom **Anlegen** des Plans, und
die Zähler liefen mit `slice-084`/`085` weiter. Die Behauptung blieb stehen,
während die Stände alterten: die Form, die `slice-085` gerade entschieden hatte
(Ursprung **und Zeitpunkt**), war hier nicht angewandt.

**Und die Gegenprobe, die F-2 gefangen hätte:** die **Summen-Spalte** in §1
(14/25/17/4 = 60) — die falsche Zuordnung erfüllte sie nicht.

**Vierter Vorgang:** `slice-081` (F-1, F-6), `slice-084` (F-1, vier Träger),
`slice-085` (F-1, die Form), `slice-088` (sechs Funde in einem Zug).

Quelle: Review zu `slice-088` (F-2, F-3) ·
Delta-Review zu `slice-088` (D-1, D-2, D-3) ·
`docs/plan/planning/in-progress/slice-088-coverage-tail-uebersetzung.md`
(berichtigt in `11ba45b`, `eb68126`, `8506539`) ·
`docs/plan/planning/observations/BEO-PGC/*/state.md`.
