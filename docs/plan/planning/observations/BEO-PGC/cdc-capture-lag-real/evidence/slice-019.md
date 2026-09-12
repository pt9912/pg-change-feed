# Beleg: slice-019 (`cdc_capture_lag` ablösen — Lasttest-Beleg)

Vorgang: slice-019 (dritter und letzter Slice von `welle-5`).

Fund: Der reale Quell-Commit-Zeitstempel (slice-017: Decoder/Mapper/
Domäne, slice-018: Store-Adapter) trägt jetzt bis in
`cdc.transaction.committed_at`. Dieser Slice löst
`cdc_capture_lag_approx` durch den kanonischen `cdc_capture_lag` ab und
belegt real per Ende-zu-Ende-Lasttest (`make test-integration`,
künstliche `docker pause`-Verzögerung um einen Quell-Commit): Baseline
≈0,1–0,16 s, verzögert (1 s Pause) ≈1,15–1,25 s, stabil über mehrere
unabhängige Läufe (Implementer, Reviewer, Verifier je eigenständig
reproduziert).

**Ausgang: eingetreten.** Die Beobachtung ist technisch vollständig
aufgelöst — `LH-FA-ADM-004` liefert jetzt einen realen, nicht nur
angenäherten CDC-Abstand.

Quelle: `docs/reviews/review-slice-019.md`, `docs/reviews/verify-slice-019.md`.
