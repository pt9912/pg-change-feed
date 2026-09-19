#!/usr/bin/env bash
# ci-matrix-abdeckung.sh — LH-QA-POR-001/002-Beleg (ADR-0105): fragt den
# jeweils letzten erfolgreichen Lauf von .github/workflows/e2e.yml
# (PostgreSQL-Versionsmatrix) und ci.yml (Linux-Plattform-Assertion) über
# die oeffentliche GitHub-REST-API ab und schreibt bei realem Erfolg
# docs/user/ci-matrix-abdeckung.md. Kein Gate (braucht Netz, wie
# make image-stale). Die Abfrage selbst laeuft in einem Container
# (AGENTS.md §3.1) — kein gh/curl auf dem Host noetig.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

REPO="pt9912/pg-change-feed"
TOOLCHAIN_IMAGE="${TOOLCHAIN_IMAGE:-golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125}"
OUT=docs/user/ci-matrix-abdeckung.md

# Fuehrt genau EINEN Container aus: installiert curl+jq (Netz ohnehin
# noetig), holt den API-Pfad ($1) und wertet ihn mit dem jq-Filter ($2)
# aus. Ergebnis auf stdout (leer bei fehlendem Feld/HTTP-Fehler).
# `set -o pipefail` in der inneren sh -c-Pipeline: ohne sie liefert `jq`
# bei leerem stdin (z. B. nach einem gescheiterten `curl`) Exit 0, und ein
# API-Fehlschlag saehe wie ein leeres Ergebnis statt wie ein Fehler aus.
api_query() {
  docker run --rm "$TOOLCHAIN_IMAGE" sh -c "
    set -o pipefail &&
    apk add --no-cache curl jq >/dev/null 2>&1 &&
    curl -fsS -H 'Accept: application/vnd.github+json' 'https://api.github.com/$1' | jq -r '$2'
  "
}

# Lauf-IDs aus der API-Antwort werden in den API-Pfad eines FOLGENDEN
# api_query-Aufrufs eingesetzt (Jobs-Abfrage) — vor dieser Wiederverwendung
# geprueft, dass der Wert rein numerisch ist: ein API-Feld mit `'` wuerde
# sonst aus der einfach gequoteten curl-URL im sh -c-String ausbrechen.
require_numeric_run_id() {
  case "$1" in
    ''|*[!0-9]*)
      echo "ci-matrix-abdeckung: $2 lieferte keine rein numerische Lauf-ID ('$1')" >&2
      exit 1
      ;;
  esac
}

echo "ci-matrix-abdeckung: frage e2e.yml (LH-QA-POR-001) ab..."
e2e_run_line=$(api_query "repos/$REPO/actions/workflows/e2e.yml/runs?status=success&per_page=1" \
  '[.workflow_runs[0].id, .workflow_runs[0].head_sha, .workflow_runs[0].html_url] | @tsv')
IFS=$'\t' read -r e2e_run_id e2e_sha e2e_url <<<"$e2e_run_line"
if [ -z "$e2e_run_id" ] || [ "$e2e_run_id" = "null" ]; then
  echo "ci-matrix-abdeckung: kein erfolgreicher e2e.yml-Lauf gefunden" >&2
  exit 1
fi
require_numeric_run_id "$e2e_run_id" "e2e.yml-Lauf-ID"

e2e_job_line=$(api_query "repos/$REPO/actions/runs/$e2e_run_id/jobs" \
  '[([.jobs[] | select(.name | contains("PostgreSQL 17")) | .conclusion][0] // "fehlt"), ([.jobs[] | select(.name | contains("PostgreSQL 18")) | .conclusion][0] // "fehlt")] | @tsv')
IFS=$'\t' read -r pg17_ok pg18_ok <<<"$e2e_job_line"
if [ "$pg17_ok" != "success" ] || [ "$pg18_ok" != "success" ]; then
  echo "ci-matrix-abdeckung: e2e.yml-Lauf $e2e_run_id — PostgreSQL-17-Job=$pg17_ok, PostgreSQL-18-Job=$pg18_ok (beide muessen success sein)" >&2
  exit 1
fi
echo "ci-matrix-abdeckung: e2e.yml-Lauf $e2e_run_id ($e2e_sha) — beide Matrix-Legs success"

echo "ci-matrix-abdeckung: frage ci.yml (LH-QA-POR-002) ab..."
ci_run_line=$(api_query "repos/$REPO/actions/workflows/ci.yml/runs?status=success&per_page=1" \
  '[.workflow_runs[0].id, .workflow_runs[0].head_sha, .workflow_runs[0].html_url] | @tsv')
IFS=$'\t' read -r ci_run_id ci_sha ci_url <<<"$ci_run_line"
if [ -z "$ci_run_id" ] || [ "$ci_run_id" = "null" ]; then
  echo "ci-matrix-abdeckung: kein erfolgreicher ci.yml-Lauf gefunden" >&2
  exit 1
fi
require_numeric_run_id "$ci_run_id" "ci.yml-Lauf-ID"

linux_step_ok=$(api_query "repos/$REPO/actions/runs/$ci_run_id/jobs" \
  '[.jobs[].steps[] | select(.name | contains("Linux-Plattform-Assertion")) | .conclusion][0] // "fehlt"')
if [ "$linux_step_ok" != "success" ]; then
  echo "ci-matrix-abdeckung: ci.yml-Lauf $ci_run_id — Linux-Plattform-Assertion-Schritt=$linux_step_ok (muss success sein)" >&2
  exit 1
fi
echo "ci-matrix-abdeckung: ci.yml-Lauf $ci_run_id ($ci_sha) — Linux-Plattform-Assertion success"

cat > "$OUT" <<EOF
# CI-Matrix-Abdeckung je Lastenheft-Kennung

Erzeugt von \`tools/harness/ci-matrix-abdeckung.sh\` (\`make doc-ci-matrix\`,
[\`ADR-0105\`](../plan/adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md)):
fragt reale, bereits abgeschlossene GitHub-Actions-Läufe über die
GitHub-REST-API ab. Diese Datei zitiert den jeweils letzten **erfolgreichen**
Lauf zum Zeitpunkt der Erzeugung — kein Beleg für den aktuellen HEAD, ein
erneuter Aufruf überschreibt sie mit dem dann aktuellsten Lauf.

| Lastenheft-Kennung | Kurzbeschreibung | Realer Beleg | Ort |
| --- | --- | --- | --- |
| [\`LH-QA-POR-001\`](../../spec/lastenheft.md) | Mehrere aktiv unterstützte PostgreSQL-Major-Versionen | Lauf [$e2e_run_id]($e2e_url), Commit \`$e2e_sha\` — beide Matrix-Legs (PostgreSQL 17/18) \`success\` | \`.github/workflows/e2e.yml\` |
| [\`LH-QA-POR-002\`](../../spec/lastenheft.md) | Linux als primäre Zielplattform | Lauf [$ci_run_id]($ci_url), Commit \`$ci_sha\` — Schritt „Linux-Plattform-Assertion" \`success\` | \`.github/workflows/ci.yml\` |
EOF

echo "ci-matrix-abdeckung: $OUT geschrieben"
