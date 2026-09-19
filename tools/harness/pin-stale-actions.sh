#!/usr/bin/env bash
# tools/harness/pin-stale-actions.sh — advisory Sensor für P9 (ADR-0051
# Pin-Inventar): jede SHA-gepinnte `uses:`-Zeile über alle
# `.github/workflows/*.yml` (AGENTS.md §3.8) gegen zwei Achsen — Tag-
# Mutation (zeigt der im Kommentar genannte Tag noch auf denselben SHA?)
# und Tag-Frische (existiert ein neuerer Release des Action-Repos?). Kein
# Gate; braucht Netz (git ls-remote + GitHub-Releases-API); meldet, hebt
# nicht an — eine Pin-Hebung ist ein bewusster Commit (Modul 14).
#
# Exit: 0 = keine Drift auf jeder Achse jeder Zeile · 1 = mindestens eine
# Achse driftet · 2 = mindestens eine Achse unbestimmbar (Netz nicht
# erreichbar). Ein doppelt vorkommendes Repo@Tag (mehrere Workflows) wird
# nur einmal geprüft.
set -uo pipefail

rc=0
seen=""

while read -r repo sha tag; do
  key="$repo@$tag"
  case " $seen " in *" $key "*) continue ;; esac
  seen="$seen $key"

  current_sha=$(git ls-remote "https://github.com/${repo}" "refs/tags/${tag}" 2>/dev/null | awk '{print $1}')
  if [ -z "$current_sha" ]; then
    echo "UNBESTIMMT  $repo@$tag — Tag nicht auflösbar (git ls-remote)"
    [ 2 -gt "$rc" ] && rc=2
  elif [ "$current_sha" = "$sha" ]; then
    echo "OK          $repo@$tag Tag-Mutation — SHA unverändert"
  else
    echo "DRIFT       $repo@$tag Tag-Mutation: gepinnt $sha, Tag zeigt jetzt auf $current_sha"
    [ 1 -gt "$rc" ] && rc=1
  fi

  latest=$(curl -fsS -m 15 "https://api.github.com/repos/${repo}/releases/latest" 2>/dev/null | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name":[[:space:]]*"([^"]+)".*/\1/')
  if [ -z "$latest" ]; then
    echo "UNBESTIMMT  $repo@$tag Tag-Frische — GitHub-Releases-API nicht erreichbar"
    [ 2 -gt "$rc" ] && rc=2
  elif [ "$latest" = "$tag" ]; then
    echo "OK          $repo@$tag Tag-Frische == neuester Release"
  else
    echo "DRIFT       $repo@$tag Tag-Frische: gepinnt $tag, neuester Release $latest"
    [ 1 -gt "$rc" ] && rc=1
  fi
done < <(grep -hoE 'uses: [A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+@[0-9a-f]{40} # v[0-9][^[:space:]]*' .github/workflows/*.yml \
  | sed -E 's/^uses: ([^@]+)@([0-9a-f]{40}) # (v[0-9.][^[:space:]]*)/\1 \2 \3/')

exit "$rc"
