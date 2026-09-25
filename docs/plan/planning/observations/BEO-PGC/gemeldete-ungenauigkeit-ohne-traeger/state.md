Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen. Die zwei
Funde des Erstauftretens haben je eine Adresse:

- Kommentar an der gRPC-Nachricht `Change` in `proto/cdc/stream/v1/changestream.proto`
  („mit denselben Feldern wie der Domain-Typ“; der Domain-Typ trägt `origin`, die
  Live-Nachricht nach dem Pflichtenheft nicht): **offen**, Adresse dieser Eintrag. Die
  Korrektur verlangt `make proto-generate` und ist ein Implementer-Zug; sie geht mit dem
  nächsten Slice, der die `.proto` ändert, oder als eigener kleiner Zug.
- Kommentar am Feld `allowDestructive` in `tools/schema/rolloutguard/guard.go`
  („nur bekannte Fremdobjekte blockieren“, neben der Klasse „View-Signatur“ ungenau):
  Adresse ist der Plan von `slice-transformationen-antragsweg-schema`, der dieselbe Datei
  ändert (§3-Zeile zu `guard.go`).

Zähler (abgeleitet): **1×** (evidence/slice-backfill-change-origin.md).
