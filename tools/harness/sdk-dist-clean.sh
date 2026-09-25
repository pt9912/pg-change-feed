#!/usr/bin/env bash
# sdk-dist-clean.sh — Hilfsfunktionen der drei sdk-pack-*-Skripte (gesourced):
# das Verzeichnis sdks/<sprache>/dist/ enthaelt nach einem Lauf genau die
# Artefakte des aktuellen Baus.
#
#   sdk_dist_check <dist>          Exit 0 nur fuer <absolut>/sdks/<name>/dist
#   sdk_dist_stage <dist>          legt ein Staging-Verzeichnis neben <dist> an
#                                  und gibt seinen Pfad aus
#   sdk_dist_swap <stage> <dist>   ersetzt <dist> durch <stage>
#
# Der Export wird zuerst nach <stage> extrahiert; erst danach wird <dist>
# geleert und ersetzt. Ein fehlgeschlagener Bau oder Export laesst <dist>
# unveraendert.

sdk_dist_check() {
  local dist=${1:-} lang
  case "$dist" in
    /*/sdks/*/dist) ;;
    *) return 1 ;;
  esac
  case "$dist" in
    *..*) return 1 ;;
  esac
  lang=${dist%/dist}
  lang=${lang##*/sdks/}
  case "$lang" in
    ''|*/*) return 1 ;;
  esac
}

sdk_dist_stage() {
  local dist=$1
  sdk_dist_check "$dist" || { echo "sdk-dist-clean: unzulaessiger Pfad: '$dist'" >&2; return 1; }
  mkdir -p -- "$(dirname -- "$dist")"
  mktemp -d "$(dirname -- "$dist")/dist-stage.XXXXXX"
}

sdk_dist_swap() {
  local stage=$1 dist=$2
  sdk_dist_check "$dist" || { echo "sdk-dist-clean: unzulaessiger Pfad: '$dist'" >&2; return 1; }
  [ -d "$stage" ] || { echo "sdk-dist-clean: Staging-Verzeichnis fehlt: '$stage'" >&2; return 1; }
  rm -rf -- "$dist"
  mv -- "$stage" "$dist"
  chmod 755 -- "$dist"
}
