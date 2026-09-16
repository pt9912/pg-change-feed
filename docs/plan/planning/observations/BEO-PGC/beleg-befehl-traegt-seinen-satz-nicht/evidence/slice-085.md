# Beleg: slice-085

Vorgang: `slice-085` — die Naht in `replication/receive`.

Fund: **Derselbe Satz, derselbe Befehl, zweiter Vorgang.** Der Verifier hat in
seiner Verifikation `verify-slice-085` **V-3** erneut gemeldet, dass
`go list -f '{{len .TestGoFiles}}'` als Beleg für „fünf Pakete ohne Testdatei"
25 liefert statt fünf.

**Das Bemerkenswerte an diesem zweiten Beleg:** der Satz lag zu diesem Zeitpunkt
**außerhalb** des Slice-Zuschnitts, und der Implementer hat ihn **bewusst liegen
gelassen** statt ihn mitzunehmen — mit einer Adresse („eigener kleiner Zug").
Er wurde damit korrekt behandelt und **zählte trotzdem**: Der Zähler misst
Wiederholung über Vorgänge, und der Fund trat in zwei unabhängigen
Verifikationen auf.

Quelle: `docs/reviews/verify-slice-085.md` (V-3) ·
`docs/reviews/review-slice-089.md` (F-5: die liegen gelassenen Fundstellen) ·
`harness/sensors/coverage-gate.md` §Grenze Punkt 1.
