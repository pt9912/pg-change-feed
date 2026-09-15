# Beleg: slice-079

Vorgang: `slice-079` — der Scope-Schnitt des Coverage-Gates.

Fund: [`ADR-0071`](../../../../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
legt als Fitness Function fest, der Messgegenstand sei **frei** von den drei
Paketen, deren Testlauf einen externen Dienst voraussetzt. Gewächtert wird das
aber nur **genähert**: die Prozent-Schwelle fängt `postgresstorage` (52,3 %) und
`replication/receive` (64,2 %) — **nicht `postgresack` allein**: käme es zurück
in `-coverpkg`, ergäbe das `(1171 + 2) / (1679 + 23) = 68,92 %` ≥ 65, das Gate
bliebe grün. Die Gegenstands-Hälfte der Fitness Function hat damit **keinen
Sensor**. Ihre Gegenseite, die Testpaket-Liste, ist überhaupt nicht beobachtbar —
sie ändert die Zahl nicht. Beides ist als Grenzpunkt 4 und 5 in
`harness/sensors/coverage-gate.md` §Grenze **benannt**: benannt, nicht
gewächtert.

Quelle: `docs/reviews/review-slice-079.md` F-2 (mit der Arithmetik) ·
`ADR-0071` §Fitness Function · `harness/sensors/coverage-gate.md` §Grenze.
