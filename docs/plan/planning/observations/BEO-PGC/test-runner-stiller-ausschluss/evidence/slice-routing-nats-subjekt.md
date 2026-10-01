**Vorgang:** slice-routing-nats-subjekt (Review F-3, LOW)

**Fund:** `tools/harness/run-notify-tests.sh` fuhr den `natsstream`-Teil mit `-run 'TestRealServer|TestPublishCost'`. Ein neuer Real-Server-Test mit anderem Präfix liefe nicht im Skript, und es käme kein Fehler: `go test -run` ohne Treffer endet mit Exit 0, und ein Test, der ohne `CDC_NATS_TEST_URL` läuft, überspringt sich. Weder das Skript noch der Kopfkommentar der Testdatei nannten die Namenskonvention als Kopplung. Die Fixrunde nahm das `-run`-Muster zurück (das Skript fährt alle Tests des Pakets) und lässt den Lauf bei jedem `--- SKIP` mit Exit 1 enden; der Verifier führte das Skript in drei Zuständen aus (grün, `FAIL`, `SKIP`).

**Form (Ausprägung):** dieselbe Klasse wie die `-run`-Muster von `tools/harness/run-integration-tests.sh`, in einem zweiten Runner-Skript (`run-notify-tests.sh`); der Fund kam vor dem Merge vom Reviewer, die Lücke ist für dieses Skript geschlossen. Für `run-integration-tests.sh` bleibt der Eintrag offen. Schwere LOW; der Eintrag ist offen und hat keinen Deckel, daher eine Datei.

Quelle: `docs/reviews/review-slice-routing-nats-subjekt.md` (F-3) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-nats-subjekt.md` (§4 M4, §5 F-3). <!-- d-check:status-provenance -->
