# Beleg: slice-harness-baseline-v6-13-0

Vorgang: `slice-harness-baseline-v6-13-0` — Erst-Review F-1 (HIGH,
`review-slice-harness-baseline-v6-13-0.md`), gemessen statt übernommen.

Fund: Der Drift-Audit (`audit-baseline-v6-13-0-drift.md`,
§2) zählte den Bundle-Diff-Umfang als „36 Dateien — 26 `regelwerk/`-Dateien,
9 Templates, 1× `SHA256SUMS`" und „19 Templates sind byte-gleich"; gemessen
per Blob-Vergleich und `diff -rq` sind es 35 / 8 / 20. Der Audit-eigene
§3-Fundtabelle benannte mit „übrige 5" plus drei genannten Templates
ebenfalls 8 geänderte Templates — §2 widersprach §3 im selben Bericht.
Die Fehlzählung war bei der Niederschrift falsch gebunden (nicht durch
spätere Arbeit gedriftet) und trug kein Mess-Kommando.

Gezogen: `9f1eb320` korrigiert §2 auf 35 / 8 / 20 und 768/285 — byte-gleich
zur Messung des Reviews und zur unabhängigen Reproduktion des Verifiers
(`git ls-tree`-Blob-Ids je relativer Pfad, `git archive` + `diff -ru`); der
Ursprung ist im Bericht deklariert (Blob-Vergleich gegen `d443ee39`, die
Fehlzählung der Erstfassung benannt statt still korrigiert).

Einordnung: Ausgang bleibt **verkörpert** (`AGENTS.md` §3.12 Instanz A,
Reviewer-Skill HIGH-Punkt „Zahl im Träger ohne Ursprung — oder gegen die
Messung driftend"); HIGH → Datei unabhängig vom Deckel. Die Träger waren
die zusammenfassende Diff-Umfang-Zeile und die byte-gleich-Aussage eines
Berichts, dessen operative Teile (Fundtabelle, Nachzugs-Entscheidungen,
Empfehlung) korrekt waren — die Zähler-Zeile driftete gegen die eigene
Fundtabelle des gleichen Trägers.
