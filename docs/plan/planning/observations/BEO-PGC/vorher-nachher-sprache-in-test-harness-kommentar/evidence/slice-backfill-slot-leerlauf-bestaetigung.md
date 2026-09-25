**Vorgang:** slice-backfill-slot-leerlauf-bestaetigung (Review-Fund F-6, Fixrunde)

**Fund:** Der Kommentar der exklusiven Instanz in `internal/bootstrap/walretention_endtoend_internal_test.go` trug neben der Kopplung („die Instanz gehört dem Test allein“) eine Nebenklausel im Konjunktiv über den nicht gewählten Aufbau: „ein gleichzeitiger Schreiber eines anderen Test-Pakets läge sonst in derselben Größenordnung wie die Schwellen“ (F-6, LOW). Der Hauptsatz trägt die Stelle (Kopplung), die Nebenklausel ist die Form aus `AGENTS.md` §3.7. Ein Test-Paket, keine Slice-/Wellen-Nummer, kein Produktionscode. Die Fixrunde formuliert im Indikativ („… der Rückstand misst das WAL der ganzen Instanz, ein gleichzeitiger Schreiber eines anderen Test-Pakets verfälscht ihn“); der Verifier las die Datei in `walretention_slotgrowth_internal_test.go`.

Quelle: `docs/reviews/review-slice-backfill-slot-leerlauf-bestaetigung.md` (F-6) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-slot-leerlauf-bestaetigung.md` (§5 F-6). <!-- d-check:status-provenance -->
