Ausgang: **verkörpert** ·
[`ADR-0103`](../../../../adr/0103-image-hash-lokal-statt-committet.md) —
`harness/image-hash.txt` ist jetzt lokal geschrieben und nicht mehr
committet (`.gitignore`); das Merge-Rauschen ist damit strukturell
ausgeschlossen, kein Binary-Hash-Ersatz (der in `observation.md` §Kandidat
skizzierte Weg) nötig. Direkt aufgelöst, nicht über den 3×-Lese-Schritt
einer Welle-Closure — der Auftraggeber hat den Fund im Dialog geprüft und
sofort entschieden.

Zähler (abgeleitet): evidence/slice-nats-drittstream-example-go.md — 1×,
Ausgang direkt zugewiesen statt bei 3× im Lese-Schritt.
