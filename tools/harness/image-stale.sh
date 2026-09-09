#!/usr/bin/env bash
# tools/harness/image-stale.sh — advisory Sensor für Base-Image-Drift
# (Modul 14). Vergleicht die FROM-Digests des Dockerfile gegen die
# aktuellen Registry-Digests derselben Tags. Kein Gate: er braucht Netz
# und urteilt nicht — er meldet. Ein Update bleibt ein bewusster Commit,
# der die Digest-Zeile anhebt (Modul 14); Major-Version-Wechsel sind
# Entscheidungen, kein automatischer Nachlauf.
#
# Exit: 0 = keine Drift · 1 = Drift gefunden · 2 = Registry nicht
# erreichbar (Ergebnis unbestimmbar — keine still grüne Aussage).
set -uo pipefail

dockerfile="${1:-Dockerfile}"
rc=0

while read -r tag pinned; do
  current=$(docker buildx imagetools inspect "$tag" --format '{{.Manifest.Digest}}' 2>/dev/null) || { echo "UNBESTIMMT  $tag — Registry nicht erreichbar"; rc=2; continue; }
  if [ "$current" = "$pinned" ]; then
    echo "OK          $tag == $pinned"
  else
    echo "DRIFT       $tag: gepinnt $pinned, aktuell $current"
    rc=1
  fi
done < <(grep -E '^FROM .*@' "$dockerfile" | awk '{split($2,a,"@"); print a[1], a[2]}')

# Major-Versions-Abdeckung: für Versionstag der Form <repo>:<MAJOR>.<MINOR><suffix>
# prüfen, ob ein Tag mit der nächsten MAJOR-Version existiert — sonst meldet der
# Digest-Vergleich oben nur Patch-Drift innerhalb des Tags, nie einen
# Major-Wechsel (Lücke, im Implementer-/Review-Kontext benannt).
while read -r tag; do
  case "$tag" in
    *:[0-9].[0-9]*) ;;
    *) continue ;;
  esac
  base="${tag%%:*}"
  ver="${tag##*:}"
  suffix=$(printf '%s' "$ver" | sed -E 's/^[0-9]+\.[0-9]+//')
  minor=$(printf '%s' "$ver" | sed -E 's/^[0-9]+\.([0-9]+).*/\1/')
  major=$(printf '%s' "$ver" | sed -E 's/^([0-9]+)\..*/\1/')
  next_maj="$major.$((minor + 1))"
  if docker manifest inspect "$base:$next_maj$suffix" >/dev/null 2>&1; then
    echo "MAJOR-DRIFT  $tag: neuer Major-Tag existiert: $base:$next_maj$suffix — Bump ist ein bewusster Commit (Modul 14)"
    rc=1
  fi
done < <(grep -E '^FROM ' "$dockerfile" | awk '{split($2,a,"@"); print a[1]}')

exit "$rc"