# Beleg: slice-archive-wellen-verbleibend

Vorgang: `slice-archive-wellen-verbleibend` — 17 reale `archive-welle`-Läufe
(`welle-2`…`welle-11`, `welle-13`…`welle-17`, `welle-19`, `welle-20`).

Fund: Eine zweite, eigenständige Auslösebedingung derselben
Werkzeug-Schwäche — diesmal auf Welle- statt Slice-Ebene. Für alle Wellen,
deren ursprüngliche Titelzeile die Form `# Welle <N>: <Titel>` trug (ohne
`welle-`-Präfix vor der Nummer), erzeugte das Werkzeug eine Stub-Titelzeile
mit verdoppelter Nummer, z. B. `# welle-14 — 14: Performance-Benchmarks &
Test-Coverage-Gate` statt `# welle-14 — Performance-Benchmarks &
Test-Coverage-Gate` — 13 von 17 archivierten Wellen betroffen
(`welle-5`, `welle-6`, `welle-7`, `welle-8`, `welle-9`, `welle-10`,
`welle-11`, `welle-13`, `welle-14`, `welle-15`, `welle-16`, `welle-17`,
`welle-19`). Die vier Wellen, deren Quelltitel bereits `Welle welle-N: …`
lautete (`welle-2`, `welle-3`, `welle-4`, `welle-20`), waren nicht
betroffen. Inhalt, Archiv-Zeiger und Fußzeilen waren bei allen 17 korrekt;
betroffen war ausschließlich die erste Zeile. Behoben durch direkte
Titelzeilen-Korrektur in derselben Closure-Runde (Review-Finding F-1).

Quelle: Review zu `slice-archive-wellen-verbleibend` F-1 · Commit `301ea0c`.
