# Beleg: slice-080

Vorgang: `slice-080` — die neue DB-Adapter-Coverage und ihre Verdrahtung in
`.github/workflows/e2e.yml`.

Fund: `make test-replication` ist **seit Längerem rot**, und niemand sieht es.
Sein tier-weites `go test ./...` scheitert an
`internal/bootstrap · TestWALRetentionThresholdEndToEnd`: dessen Fixture fährt
`DROP SCHEMA cdc CASCADE` und ein eigenes `ApplySchema`, das
`cdc.administration_request` und `cdc.process_heartbeat` nicht mitbringt — beide
liest `bootstrap.Run` seit `ADR-0050`. Der Fehler ist **vorbestehend**, am
**unveränderten** Runner reproduziert (`42P01`), und er war unsichtbar, weil
`make test-replication` **weder Gate noch CI-Schritt** ist: es läuft nur, wenn
jemand es von Hand aufruft.

Sichtbar wurde er erst, als ein neuer Träger ihn aufrief — die DB-Adapter-Coverage
braucht den Replication-Lauf, um ihre Zahl zu messen. Damit ist der Befund
zugleich der Beleg für die Klasse: **ein Beleg, den es gibt, aber niemand liest,
ist keiner.** Das ist die Umkehrung der Klasse aus
`BEO-PGC/regel-weiter-als-ihr-sensor` (dort ist der Sensor enger als die Regel,
hier ist der Lauf da und wird von keinem Gate abgeholt).

Quelle: Review zu `slice-080`, F-2/F-4 (mit dem unabhängigen Nachweis
am unveränderten Runner) · `.github/workflows/e2e.yml` · `tools/harness/run-replication-tests.sh`.
