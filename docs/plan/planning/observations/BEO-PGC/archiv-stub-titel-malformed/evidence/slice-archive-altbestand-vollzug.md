# Beleg: slice-archive-altbestand-vollzug

Vorgang: `slice-archive-altbestand-vollzug` — realer `archive-welle`-Lauf
für die Schlüssel `altbestand` und `welle-d-check`.

Fund: Für alle drei archivierten Slices mit namensbasierter Kennung
(`MR-002`-Fälle, nicht `slice-<NNN>`) lautete die generierte Titelzeile
`# slice- — slice-<name>: <Titel>` statt der Template-Form
`# slice-<Kennung> — <Titel>` — das Kennung-Feld war leer und der Slug
wiederholte sich im Titeltext (`docs/plan/planning/done/welle-d-check/slice-d-check-trace-rtm.md`,
`slice-d-check-tracked-modul.md`, `docs/plan/planning/done/altbestand/slice-gate-index-konsolidierung.md`).
Inhalt, Archiv-Zeiger und Fußzeilen waren korrekt; betroffen war
ausschließlich die erste Zeile. Behoben durch direkte Titelzeilen-Korrektur
in derselben Closure-Runde (Review-Finding F-1).

Quelle: `docs/reviews/review-slice-archive-altbestand-vollzug.md` F-1 ·
Commit `641eaf0`.
