# tools/harness/lib-pin-compare.sh — gemeinsamer Index-Digest-Vergleich für
# tools/harness/pin-stale.sh und tools/harness/pin-stale-all.sh (ADR-0146).
# Wird per `source` eingebunden, kein eigenes Executable, setzt keine
# Shell-Optionen.
#
# pin_compare_digest <label> <vergleichs-ziel> <gepinnter-digest>
# Fragt den Index-Digest des Ziels ab und druckt genau eine Zeile:
#   OK          <label> (<ziel>) == <digest>              Rückgabe 0
#   DRIFT       <label> (<ziel>): gepinnt …, aktuell …    Rückgabe 1
#   UNBESTIMMT  <label> — <ziel>: Registry nicht erreichbar  Rückgabe 2
# PIN_COMPARE_TIMEOUT (Sekunden, optional) begrenzt den einzelnen
# Registry-Aufruf; ein abgebrochener Aufruf ist UNBESTIMMT. Ohne Wert
# läuft der Aufruf unbegrenzt.
pin_compare_digest() {
  local label="$1" target="$2" pinned="$3" current
  local -a limit=()
  if [ -n "${PIN_COMPARE_TIMEOUT:-}" ]; then
    limit=(timeout "$PIN_COMPARE_TIMEOUT")
  fi

  current=$(${limit[@]+"${limit[@]}"} docker buildx imagetools inspect "$target" --format '{{.Manifest.Digest}}' 2>/dev/null) || {
    echo "UNBESTIMMT  $label — $target: Registry nicht erreichbar"
    return 2
  }

  if [ "$current" = "$pinned" ]; then
    echo "OK          $label ($target) == $pinned"
    return 0
  fi

  echo "DRIFT       $label ($target): gepinnt $pinned, aktuell $current"
  return 1
}
