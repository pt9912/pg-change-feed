#!/usr/bin/env bash
# tools/harness/dockerfile-from.sh — liest die Image-Referenzen der
# FROM-Zeilen eines Dockerfile für tools/harness/image-stale.sh.
#
# Aufruf: dockerfile-from.sh [dockerfile]   (Default: Dockerfile)
# Ausgabe: je FROM-Zeile eine Zeile `<referenz> <digest>` — <referenz> ist
# die Image-Angabe ohne `@<digest>`, <digest> der Teil hinter `@` (leer,
# wenn die Zeile nicht digest-gepinnt ist; die Zeile endet dann auf das
# Leerzeichen hinter <referenz>). Optionen (`--platform=…`) und `AS <name>`
# gehören nicht zur Referenz. Netzlos, reine Textlogik.
set -uo pipefail

dockerfile="${1:-Dockerfile}"

awk '
  toupper($1) == "FROM" {
    for (i = 2; i <= NF; i++) {
      if ($i ~ /^--/) continue
      n = split($i, a, "@")
      print a[1], (n > 1 ? a[2] : "")
      break
    }
  }
' "$dockerfile"
