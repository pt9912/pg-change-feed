# Beleg: slice-022 (Positions-Bestätigung über denselben Zugriffsweg)

Vorgang: slice-022 (zweiter Slice von `welle-6`).

Fund: `review-slice-022.md` F-2 — mehrere `postgresstorage`-Tests
(`consumerstate_test.go`, `store_test.go`, `tableactivation_test.go`)
mutieren den geteilten `cdc`-Schema-/Tabellen-Zustand unskopiert. Mit
dem unveränderten `run-store-tests.sh` (alle Pakete in einem `go test
./...`-Aufruf) reproduzierte sich real und reproduzierbar eine Race
gegen den neuen `internal/bootstrap`-Test (4 von 5 Läufen scheiterten).
Der Fix (`tools/harness/run-store-tests.sh`, paketweise
Vorzieh-Isolierung von `internal/bootstrap`) behebt das Symptom, nicht
die Ursache — Reviewer-Verdikt: kein Merge-Blocker, aber zweiter
unabhängiger Workaround für dasselbe Grundmuster im Repo (erster:
`sqlviews_test.go`s Dateinamen-Sortierung, vor der Registrierung dieser
Beobachtung entstanden — nachrichtlich, nicht gezählt).

**Ausgang: weiter offen.** Die Ursachenbehebung (unskopierte
`DELETE`/`DROP SCHEMA CASCADE` in den Adapter-Tests selbst) berührt
mehrere bestehende Testdateien zugleich; kein Slice dafür existiert.

Quelle: `docs/reviews/review-slice-022.md` F-2, `docs/reviews/verify-slice-022.md`.
