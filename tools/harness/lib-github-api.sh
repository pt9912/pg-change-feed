# tools/harness/lib-github-api.sh — gemeinsame containerisierte
# GitHub-API-GET-Funktion für tools/harness/pin-stale-*.sh (P7-P9,
# ADR-0051 Pin-Inventar). Analog tools/harness/ci-matrix-abdeckung.sh:
# die Abfrage selbst laeuft in einem Container (AGENTS.md §3.1) — kein
# curl auf dem Host noetig. Wird per `source` eingebunden, kein eigenes
# Executable.
GITHUB_API_TOOLCHAIN_IMAGE="${GITHUB_API_TOOLCHAIN_IMAGE:-golang:1.27-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414}"

# github_api_get <pfad> — GET https://api.github.com/<pfad>, rohe
# JSON-Antwort auf stdout (leer bei HTTP-Fehler/Timeout).
github_api_get() {
  docker run --rm "$GITHUB_API_TOOLCHAIN_IMAGE" sh -c "
    apk add --no-cache curl >/dev/null 2>&1 &&
    curl -fsS -m 15 -H 'Accept: application/vnd.github+json' 'https://api.github.com/$1'
  " 2>/dev/null
}
