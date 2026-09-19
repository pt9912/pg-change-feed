# CI-Matrix-Abdeckung je Lastenheft-Kennung

Erzeugt von `tools/harness/ci-matrix-abdeckung.sh` (`make doc-ci-matrix`,
[`ADR-0105`](../plan/adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md)):
fragt reale, bereits abgeschlossene GitHub-Actions-Läufe über die
GitHub-REST-API ab. Diese Datei zitiert den jeweils letzten **erfolgreichen**
Lauf zum Zeitpunkt der Erzeugung — kein Beleg für den aktuellen HEAD, ein
erneuter Aufruf überschreibt sie mit dem dann aktuellsten Lauf.

| Lastenheft-Kennung | Kurzbeschreibung | Realer Beleg | Ort |
| --- | --- | --- | --- |
| [`LH-QA-POR-001`](../../spec/lastenheft.md) | Mehrere aktiv unterstützte PostgreSQL-Major-Versionen | Lauf [35404712097](https://github.com/pt9912/pg-change-feed/actions/runs/35404712097), Commit `22a3b63fa93bd070a94d9a2192ba8d5e96bc79da` — beide Matrix-Legs (PostgreSQL 17/18) `success` | `.github/workflows/e2e.yml` |
| [`LH-QA-POR-002`](../../spec/lastenheft.md) | Linux als primäre Zielplattform | Lauf [35404712089](https://github.com/pt9912/pg-change-feed/actions/runs/35404712089), Commit `22a3b63fa93bd070a94d9a2192ba8d5e96bc79da` — Schritt „Linux-Plattform-Assertion" `success` | `.github/workflows/ci.yml` |
