#!/usr/bin/env bash
# tools/harness/pin-stale.sh — advisory Sensor für Docker-Digest-Pin-Drift
# außerhalb des Dockerfile (P3–P6, ADR-0051 Pin-Inventar). Analog
# tools/harness/image-stale.sh, aber Quelle ist eine Makefile-/mk-Variable
# (`VAR ?= <image[:tag]>@sha256:<digest>` oder, wenn der Pin selbst kein
# Tag trägt, `VAR ?= <image>@sha256:<digest>` zusammen mit einem
# expliziten Vergleichs-Tag als drittem Argument). Kein Gate: braucht Netz
# und urteilt nicht — er meldet. Ein Update bleibt ein bewusster Commit.
#
# Aufruf: pin-stale.sh <datei> <variable> [vergleichs-tag]
# Exit: 0 = keine Drift · 1 = Drift gefunden · 2 = Registry nicht
# erreichbar, Variable nicht gefunden oder Aufruf-Fehler (Pflichtargument
# fehlt) — in jedem Fall ein unbestimmbares Ergebnis, nie ein DRIFT-Fund.
set -uo pipefail

if [ $# -lt 2 ]; then
  echo "UNBESTIMMT  Aufruf — Pflichtargument fehlt (pin-stale.sh <datei> <variable> [vergleichs-tag])"
  exit 2
fi
file="$1"
var="$2"
compare_tag="${3:-}"

line=$(grep -E "^${var}[[:space:]]*\?=" "$file" || true)
if [ -z "$line" ]; then
  echo "UNBESTIMMT  $var — Variable nicht in $file gefunden"
  exit 2
fi
value=$(printf '%s' "$line" | sed -E "s/^${var}[[:space:]]*\?=[[:space:]]*//" | awk '{print $1}')

case "$value" in
  sha256:*)
    pinned="$value"
    ;;
  *@sha256:*)
    pinned="${value#*@}"
    ;;
  *)
    echo "UNBESTIMMT  $var — Wert $value trägt keinen sha256-Digest"
    exit 2
    ;;
esac

if [ -n "$compare_tag" ]; then
  target="$compare_tag"
else
  target="${value%@*}"
  if [ "$target" = "$value" ]; then
    echo "UNBESTIMMT  $var — kein Vergleichs-Tag ableitbar (Pin trägt nur den Digest, kein drittes Argument übergeben)"
    exit 2
  fi
fi

current=$(docker buildx imagetools inspect "$target" --format '{{.Manifest.Digest}}' 2>/dev/null) || {
  echo "UNBESTIMMT  $var — $target: Registry nicht erreichbar"
  exit 2
}

if [ "$current" = "$pinned" ]; then
  echo "OK          $var ($target) == $pinned"
  exit 0
fi

echo "DRIFT       $var ($target): gepinnt $pinned, aktuell $current"
exit 1
