Zustand: offen (1×) — unter der Schwelle, kein Ausgang zugewiesen. Ein Träger
ist nicht vorgeschlagen: die Klasse wäre über einen Sensor auf dem Artefakt
schwer zu fassen — verlangte er, dass jede Zahl ihren Ursprung trägt, wäre er
eine Formpflicht auf Prosa und erzeugte Pflichterfüllung. Die naheliegende
Antwort ist die **Herkunfts-Regel**, die der Architect-Zug zu `ADR-0078` für die
Schwester-Klasse benannt hat (gemessen / übernommen / **abgeleitet**); ob sie
diese Klasse mitdeckt, ist beim Lese-Schritt zu entscheiden.

Zähler (abgeleitet): **1×** (evidence/slice-081.md). Die zwei Fundstellen
(`db-coverage.sh:18` und `ADR-0071`/`welle-20.md`) liegen **im selben** Vorgang
und sind damit *eine* Gelegenheit — der Zähler misst Wiederholung über Vorgänge,
nicht die Zahl der Funde.

**Nicht zu verwechseln** mit dem benachbarten Eintrag
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (2×): dort sagt ein
Kommentar einen **Fehlerpfad** zu, den der Code nicht trägt (die Aussage kehrt
das Verhalten um); hier nennt ein Träger eine **Zahl**, die gegen die Messung
driftet. Der Review zu `slice-081` hatte seinen F-1 jenem Eintrag zugeordnet —
das ist nach Prüfung der beiden Vorgänger-Findings (`review-slice-070` F-2,
`review-slice-077` F-1, beide Fehlerpfad-Fälle) **nicht** dieselbe Beobachtung.
