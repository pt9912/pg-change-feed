# Beleg: slice-sdk-python-projektgeruest

Rückwirkend erfasst bei der Planner-Closure von
`slice-sdk-kotlin-projektgeruest` (2026-09-20) — zum Zeitpunkt des
Python-Slices wurde diese Formulierung von keinem Review als eigene Klasse
benannt.

`sdks/python/README.md:9` trägt „gRPC, SSE and NATS-vollinhalt delivery
remain out of scope for this package's first release …" — das deutsche
Fachwort „vollinhalt" steht unübersetzt mitten im englischen Satz. Real
nachgemessen: `grep -n "vollinhalt" sdks/python/README.md` → Zeile 9.
