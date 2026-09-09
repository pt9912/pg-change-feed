# Beleg: slice-002

Vorgang: slice-002 (Domänenkern — Modelle, Invarianten, ClockPort, in
`done/` seit seiner Closure, welle-1).

Fund: die Layer-Globs `domain` und `ports` matchen erstmals echten Content
(`internal/domain/**`, `internal/application/port/**`) — `make a-check`
prüft sie mit 0 Befunden, ohne Abdeckungs-/Auflösungs-Hinweis für diese
Globs. **Die Null-Abdeckung gilt weiter für `app` (`usecase/**`) und
`adapters`** — beide matchen null Dateien bis slice-003; die
Beobachtungsklasse ist damit teils entkräftet, nicht abgeschlossen.

Quelle: `docs/reviews/verify-slice-002.md` (a-check im Gate-Bündel,
Kante `ports → domain` per `a-check-graph` bestätigt).
