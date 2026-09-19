#!/usr/bin/env bash
# tools/harness/pin-stale-dcheck.sh — advisory Sensor für P7 (ADR-0051
# Pin-Inventar): d-check trägt ZWEI Achsen, anders als P3-P6 — Tag-Frische
# (existiert ein neuerer d-check-Release als der gepinnte Tag?) UND
# Digest-Drift (trägt derselbe Tag noch denselben Bau?). Kein Gate; braucht
# Netz (GitHub-Releases-API + Registry); meldet, hebt nicht an.
#
# Exit: 0 = keine Drift auf beiden Achsen · 1 = mindestens eine Achse
# driftet · 2 = eine Achse unbestimmbar (Netz/API nicht erreichbar).
set -uo pipefail

file="d-check.mk"
rc=0

pinned_tag=$(grep -E '^DCHECK_IMAGE[[:space:]]*\?=' "$file" | sed -E 's/^DCHECK_IMAGE[[:space:]]*\?=[[:space:]]*//' | awk '{print $1}')
pinned_tag="${pinned_tag##*:}"

# Achse A — Digest-Drift: der gepinnte Tag, aktuell inspiziert.
bash "$(dirname "$0")/pin-stale.sh" "$file" DCHECK_DIGEST "ghcr.io/pt9912/d-check:${pinned_tag}"
rc_a=$?
[ "$rc_a" -gt "$rc" ] && rc=$rc_a

# Achse B — Tag-Frische: neuester GitHub-Release von pt9912/d-check ggü.
# dem gepinnten Tag.
latest=$(curl -fsS -m 15 https://api.github.com/repos/pt9912/d-check/releases/latest 2>/dev/null | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name":[[:space:]]*"([^"]+)".*/\1/')
if [ -z "$latest" ]; then
  echo "UNBESTIMMT  DCHECK_IMAGE Tag-Frische — GitHub-Releases-API nicht erreichbar"
  [ 2 -gt "$rc" ] && rc=2
elif [ "$latest" = "$pinned_tag" ]; then
  echo "OK          DCHECK_IMAGE Tag-Frische ($pinned_tag) == neuester Release"
else
  echo "DRIFT       DCHECK_IMAGE Tag-Frische: gepinnt $pinned_tag, neuester Release $latest"
  [ 1 -gt "$rc" ] && rc=1
fi

exit "$rc"
