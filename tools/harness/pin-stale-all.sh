#!/usr/bin/env bash
# tools/harness/pin-stale-all.sh — advisory Sensor P10 (ADR-0146): jede
# Referenz der Form <image>[:<tag>]@sha256:<digest> in einer getrackten Datei
# außerhalb von docs/ und .harness/ wird gegen den Index-Digest der Registry
# verglichen. Die Referenzen kommen aus `git grep`, nicht aus einer Liste;
# ein Pin ohne Tag wird gegen :latest verglichen. Kein Gate: braucht Netz,
# meldet und hebt nicht an. Der Sensor liest die Form, nicht den Sinn.
#
# Meldung je verschiedene Referenz (erster Fundort `datei:zeile`, bei
# mehreren `+N weitere`): OK · DRIFT · UNBESTIMMT — Zeilenform von
# lib-pin-compare.sh.
#
# Exit: 1 = mindestens ein DRIFT · 2 = UNBESTIMMT ohne DRIFT, leerer
# Gegenstand (keine Referenz gefunden) oder `git grep`-Fehler · 0 = alle OK.
# PIN_COMPARE_TIMEOUT (Sekunden je Registry-Aufruf, Default 60).
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"
# shellcheck source=tools/harness/lib-pin-compare.sh
. "$(dirname "$0")/lib-pin-compare.sh"
export PIN_COMPARE_TIMEOUT="${PIN_COMPARE_TIMEOUT:-60}"

pattern='[a-z0-9][a-z0-9./_-]*(:[A-Za-z0-9._-]+)?@sha256:[0-9a-f]{64}'

hits=$(git grep -nHoE "$pattern" -- . ':!docs' ':!.harness')
grep_rc=$?
if [ "$grep_rc" -eq 1 ]; then
  echo "UNBESTIMMT  keine Digest-Pin-Referenz im Baum gefunden — leerer Gegenstand"
  exit 2
elif [ "$grep_rc" -ne 0 ]; then
  echo "UNBESTIMMT  git grep endete mit Exit $grep_rc"
  exit 2
fi

declare -A first_at=()
declare -A extra=()
refs=()
while IFS=: read -r file line ref; do
  if [ -z "${first_at[$ref]+x}" ]; then
    first_at[$ref]="$file:$line"
    extra[$ref]=0
    refs+=("$ref")
  else
    extra[$ref]=$((extra[$ref] + 1))
  fi
done <<< "$hits"

n_ok=0
n_drift=0
n_unbest=0
for ref in "${refs[@]}"; do
  image="${ref%@*}"
  pinned="${ref#*@}"
  case "${image##*/}" in
    *:*) target="$image" ;;
    *) target="$image:latest" ;;
  esac
  label="${first_at[$ref]}"
  if [ "${extra[$ref]}" -gt 0 ]; then
    label="$label +${extra[$ref]} weitere"
  fi
  pin_compare_digest "$label" "$target" "$pinned"
  case $? in
    0) n_ok=$((n_ok + 1)) ;;
    1) n_drift=$((n_drift + 1)) ;;
    *) n_unbest=$((n_unbest + 1)) ;;
  esac
done

echo "pin-stale-all: ${#refs[@]} Referenzen — $n_ok OK, $n_drift DRIFT, $n_unbest UNBESTIMMT"
if [ "$n_drift" -gt 0 ]; then
  exit 1
fi
if [ "$n_unbest" -gt 0 ]; then
  exit 2
fi
exit 0
