#!/usr/bin/env bash
# working-tree-hash — inhaltsbasierter Hash ueber alle getrackten UND untracked
# Dateien (exclude-standard, respektiert .gitignore). Gemeinsame Quelle der
# Wahrheit fuer record-gates und den Stop-Hook — keine Logik-Dopplung.
#
# INHALTSBASIERT statt diff-basiert: ein Commit ohne Inhaltsaenderung laesst den
# Hash unveraendert — der Gate-Nachweis bleibt gueltig. Ein Commit OHNE
# vorherigen Gate-Lauf erzeugt dagegen keinen passenden Nachweis und wird vom
# Stop-Hook nicht mehr durchgewunken.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

# sha256sum ist GNU-Coreutils (Linux, WSL); macOS traegt es nicht, wohl aber
# das formatkompatible shasum (Perl, Teil der Basisinstallation seit jeher —
# keine Fremd-Laufzeit im Sinne von node/jq/python).
if command -v sha256sum >/dev/null 2>&1; then
  hash_stdin() { sha256sum; }
  hash_file() { sha256sum -- "$1"; }
else
  hash_stdin() { shasum -a 256; }
  hash_file() { shasum -a 256 -- "$1"; }
fi

# `sort -z` (NUL-getrennt) ist GNU-only — BSD-/macOS-sort kennt das Flag
# ueberhaupt nicht. GRENZE des Fallbacks: er sortiert ueber einen
# Zeilenumbruch-getrennten Zwischenschritt und ist deshalb nicht mehr
# NUL-sicher gegen einen Dateinamen mit eingebettetem Zeilenumbruch (beim
# GNU-Pfad weiterhin sicher) — ein in diesem Repo bislang nie aufgetretener
# Dateiname-Fall, keine neue Fremd-Laufzeit noetig.
if printf '' | sort -z >/dev/null 2>&1; then
  sort_null() { sort -zu; }
else
  sort_null() {
    local f
    while IFS= read -r -d '' f; do printf '%s\n' "$f"; done \
      | sort -u \
      | while IFS= read -r f; do printf '%s\0' "$f"; done
  }
fi

git ls-files -z --cached --others --exclude-standard \
  | sort_null \
  | while IFS= read -r -d '' f; do
      if [ -L "$f" ]; then
        printf 'LINK %s -> %s\n' "$f" "$(readlink -- "$f")"
      elif [ -f "$f" ]; then
        hash_file "$f"
      else
        printf 'GONE %s\n' "$f"
      fi
    done | hash_stdin | awk '{print $1}'
