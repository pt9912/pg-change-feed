#!/usr/bin/env bash
# kommentar-kennungen — Aufrufer von `make kommentar-kennungen`: startet das Go-
# Programm tools/harness/kommentar-kennungen im gepinnten Toolchain-Container
# (netzlos, Repo lesend gemountet). Mit einer Diff-Basis liest das Programm
# einen `git diff -U0`-Strom; der Diff entsteht auf dem Host (`git`) in einer
# Temp-Datei, weil `pipefail` in einer Pipe den Exit des rechten Glieds vor dem
# des linken meldet und ein fehlgeschlagenes `git` hinter Exit 1 des Programms
# verschwände. Vertrag und Exit-Codes: harness/sensors/kommentar-kennungen.md
# (AGENTS.md §3.7).
#
# Aufruf: kommentar-kennungen.sh <Pfade> <COUNT> <TESTS> <DIFF>
#   Pfade  Leerzeichen-getrennt, relativ zur Repo-Wurzel; leer = der Baum
#   COUNT  nicht leer = nur die Zahl der Kandidaten drucken
#   TESTS  leer oder all | exclude | only
#   DIFF   leer oder eine Commit-Kennung/ein Ref als Basis
# Umgebung: TOOLCHAIN_IMAGE (Pflicht, setzt das Makefile)
# Exit: 0 kein Kandidat (oder Zahl gedruckt) · 1 mindestens ein Kandidat · 2 Eingabefehler
set -uo pipefail

if [ $# -ne 4 ]; then
  echo "Aufruf: kommentar-kennungen.sh <Pfade> <COUNT> <TESTS> <DIFF>" >&2
  exit 2
fi
paths=$1
count=$2
tests=$3
diff_base=$4
if [ -z "${TOOLCHAIN_IMAGE:-}" ]; then
  echo "kommentar-kennungen: TOOLCHAIN_IMAGE fehlt (das Makefile setzt es)" >&2
  exit 2
fi

root=$(git rev-parse --show-toplevel) || exit 2
cd "$root" || exit 2

args=()
if [ -n "$count" ]; then args+=(-count); fi
if [ -n "$tests" ]; then args+=(-tests "$tests"); fi

difffile=''
if [ -n "$diff_base" ]; then
  if ! git rev-parse --verify --quiet "$diff_base^{commit}" >/dev/null; then
    echo "kommentar-kennungen: DIFF '$diff_base' ist kein Commit dieses Repos" >&2
    exit 2
  fi
  difffile=$(mktemp "${TMPDIR:-/tmp}/kommentar-kennungen.XXXXXX") || exit 2
  trap 'rm -f "$difffile"' EXIT
  if ! git diff -U0 --no-color --no-ext-diff "$diff_base" -- '*.go' >"$difffile"; then
    echo "kommentar-kennungen: git diff gegen '$diff_base' schlug fehl" >&2
    exit 2
  fi
  args+=(-diff)
fi

# Pfade zerlegen (Leerzeichen-getrennt, ohne Expansion).
read -r -a path_words <<<"$paths"
args+=(${path_words[@]+"${path_words[@]}"})

# `go build` + `exec` statt `go run`: `go run` meldet den Exit des Programms nicht
# unverändert weiter, die Exit-Codes 1 und 2 des Vertrags blieben ununterscheidbar.
run=(docker run --rm -i --network none -v "$root":/src:ro -w /src
  -e GOCACHE=/tmp/gocache -e GOFLAGS=-buildvcs=false "$TOOLCHAIN_IMAGE"
  sh -c 'go build -o /tmp/kommentar-kennungen ./tools/harness/kommentar-kennungen && exec /tmp/kommentar-kennungen "$@"'
  kommentar-kennungen "${args[@]}")

if [ -n "$difffile" ]; then
  "${run[@]}" <"$difffile"
else
  "${run[@]}" </dev/null
fi
