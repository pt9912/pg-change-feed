# Beleg: review-slice-079

Vorgang: das Review von `slice-079` — gefunden im Bestätigungslauf
(`docs/reviews/review-slice-079-fixrunde-2.md`, F-1).

Fund: Die neue Sektion §Zählbasis in `harness/sensors/coverage-gate.md` erklärte
die **gedruckte** Prozentzeile mit einem Mechanismus, den das Werkzeug nicht hat
(„jede Block-Position zählt so oft, wie sie vorkommt"). Gemessen: `go tool cover`
führt identische Block-Positionen zu **einem** Block zusammen — der Wert ist
gegen die Duplikate **invariant**; die per-Zeile-Naivsumme ergäbe 3,95 % statt
69,9 %. Korrigiert: die gedruckte Zeile ist **keine** eigene Größe, sondern
dieselbe Messung in anderer Ausgabepräzision.

Quelle: `docs/reviews/review-slice-079-fixrunde-2.md` F-1 · Commit `0dd531d`.
