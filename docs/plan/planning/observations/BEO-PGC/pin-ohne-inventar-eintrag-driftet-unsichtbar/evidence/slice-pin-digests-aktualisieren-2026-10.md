**Vorgang:** slice-pin-digests-aktualisieren-2026-10
**Fund:** Der Implementer maß den Pin `postgres:17-alpine` in
`.github/workflows/e2e.yml` von Hand: der amd64-Eintrag der Registry trug
`sha256:aa90e97ee862e558111d34cfb8b2c4bec768c2b039fb791341686928560263b3`
(Annotation `17.11-alpine3.24`), der Pin
`sha256:7456ef82e5f5bc43d997f4781bbd7c0d6389bff397564649a356e206ba473aee`;
der Verifier bestätigte die Messung unabhängig (Verifikationsbericht Punkt 4,
Befehl `docker buildx imagetools inspect postgres:17-alpine`). Weder
`make pin-stale-pgtest` (liest nur die Variable `PG_TEST_IMAGE`, PostgreSQL 18)
noch `upstream-drift.yml` (Läufe 37112891025 und alle 13 davor rot, aber ohne
Zeile für PostgreSQL 17) meldeten diese Drift. Zusätzlich gelesen
(Verifikationsbericht §8.6, per `git grep -E '^FROM .*@sha256|image: .*@sha256'`):
fünf Basis-Images in `sdks/` und `examples/` ohne Sensor, ihre Drift ist nicht
gemessen. Erstauftreten.
