# Beleg: slice-sdk-csharp-projektgeruest

Rückwirkend erfasst bei der Planner-Closure von
`slice-sdk-kotlin-projektgeruest` (2026-09-20) — zum Zeitpunkt des
C#-Slices wurde diese Formulierung von keinem Review als eigene Klasse
benannt.

`sdks/csharp/README.md:9` trägt „SSE and NATS-vollinhalts delivery remain
uncovered by this package (`ADR-0106` Festlegung 1)" — das deutsche
Fachwort „vollinhalts" (unflektierte Form von „Vollinhalts-Zustellweg",
vgl. `ADR-0100`) steht unübersetzt mitten im englischen Satz. Real
nachgemessen: `grep -n "vollinhalt" sdks/csharp/README.md` → Zeile 9.
