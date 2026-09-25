**Vorgang:** slice-backfill-bench-richtgroesse (Review F-3, Verifikation §5 F-3)

**Fund:** Die Fixrunde des Slice ergänzte im Skript `tools/bench-backfill.sh` die Messung des WAL-Rückstands (Rückstand je Run, gehaltenes WAL, Freigabe durch einen Live-Commit, abgeleitete Schwellen-Zeile) und den Speicher 20 s nach dem letzten Run. Der Vertrag `harness/targets/bench-backfill.md` und die `make bench`-Zeile in `harness/README.md` beschrieben diese Größen nicht (F-3, LOW); beide Träger stehen nicht im Skript-Diff und kein Sensor liest sie. Gefunden vom Reviewer; die Fixrunde zog beide nach, der Verifier prüfte, dass jede beschriebene Größe in der gedruckten Ausgabe eines eigenen Laufs erscheint (Verifikation §5 F-3).

Quelle: `docs/reviews/review-slice-backfill-bench-richtgroesse.md` (F-3) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-bench-richtgroesse.md` (§5 F-3). <!-- d-check:status-provenance -->
