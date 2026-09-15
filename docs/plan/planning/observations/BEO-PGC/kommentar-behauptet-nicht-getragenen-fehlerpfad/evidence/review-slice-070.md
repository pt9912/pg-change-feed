# Beleg: review-slice-070

Vorgang: das Review von `slice-070` (F-2, HIGH) — gefunden beim Bau der
gRPC-Capture-Integration.

Fund: Der neue `wiring.go`-Kommentar sagte zu, „der gRPC-Server-Startfehler
erreicht das Ergebnis von `Run` weiter unten auf demselben Pfad wie jeder andere
Adapter-Startfehler". Der Code startet den Server in eigener Goroutine und
verwirft den Fehler per `log.Error`; `Run` gibt ausschließlich
`mergeStreamAndWALFaultOutcome(streamErr, &walFault)` zurück. Eine
fehlkonfigurierte Adresse beendet den Lauf nicht — der Kommentar legte das
Gegenteil nahe. Behoben in jenem Slice („eine Zeile, kein Verhalten"), aber an
**einem** Block; ein zweiter mit derselben Behauptung blieb stehen
(s. `review-slice-077.md`).

Quelle: `docs/reviews/review-slice-070.md` F-2 ·
`docs/plan/planning/done/slice-070-grpc-capture-integration.md`.
