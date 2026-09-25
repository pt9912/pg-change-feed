**Vorgang:** slice-backfill-change-origin (Closure-Notiz §7, Closure-Note-Review von
`welle-backfill-bestand`, F-2)

**Fund:** Die Closure-Notiz meldet zwei Ungenauigkeiten, die der Slice nicht änderte:
der Kommentar an der gRPC-Nachricht in `proto/cdc/stream/v1/changestream.proto`
(Zeilen 12 und 13, „mit denselben Feldern wie der Domain-Typ“; der Domain-Typ trägt
`origin`) — eine Änderung verlangt `make proto-generate`, die Notiz nennt sie „gemeldete
Ungenauigkeit ohne Träger-Slice“ — und der Kommentar am Feld `allowDestructive` in
`tools/schema/rolloutguard/guard.go` (Zeilen 53 und 54, „nur bekannte Fremdobjekte
blockieren“, neben der Klasse „View-Signatur“ ungenau). Weder die Results-Notiz noch ein
Folge-Slice noch ein Register-Eintrag führte sie; beide Stellen standen zu HEAD
`6839f738` unverändert (gemessen mit `sed -n 12,13p` und `sed -n 53,54p`).

Quelle: `docs/reviews/review-closure-notes-welle-backfill-bestand.md` (F-2). <!-- d-check:status-provenance -->
