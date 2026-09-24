**Vorgang:** slice-backfill-e2e (Review-Fund F-5, Fixrunde)

**Fund:** Der Godoc von `TestE2EBackfillReplayInvariant` (`test/integration/backfill_e2e_test.go`) schloss mit „Der Test bricht ab, wenn eine der beiden Seiten leer bleibt, statt ohne Überlappung grün zu werden“ — die Zusage steht vorn, das „statt …“ nennt die Alternative; der Runner-Kommentar in `tools/harness/run-integration-tests.sh` erklärte „(eine offene Transaktion während der Slot-Anlage würde diese verzögern)“ im Konjunktiv über den nicht gewählten Ablauf (F-5, LOW). Der Hauptsatz trägt in beiden Fällen die Stelle (Zusage bzw. Kopplung); die Nebenklausel ist die Form aus `AGENTS.md` §3.7. Kein Produktionscode, keine Slice-/Wellen-Nummer. Die Fixrunde formulierte beide im Indikativ („Der Test bricht ab, wenn eine der beiden Seiten leer bleibt.“; „(die Slot-Anlage wartet auf jede offene Schreibtransaktion)“); der Verifier las beide Diffs (je 1:1 Zeilen).

Quelle: `docs/reviews/review-slice-backfill-e2e.md` (F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-e2e.md` (§5 F-5). <!-- d-check:status-provenance -->
