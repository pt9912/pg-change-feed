#!/usr/bin/env bash
# generated-sync.sh — Sync-Gate des generierten Protobuf-/gRPC-Codes
# (ADR-0084 Festlegung 1; loest die Folgepflicht des Generators aus ADR-0060 ein).
#
# Gegenstand ist genau eine Paarung: committetes Erzeugnis = Ausgabe des
# gepinnten Generators aus der committeten Quelle. Nicht Gegenstand: dass der
# erzeugte Code kompiliert oder sich richtig verhaelt (`make test`, `make image`)
# und dass die .proto den Draht-Vertrag richtig beschreibt (die Belege zu
# LH-FA-SST-008).
#
# Der Arbeitsbaum bleibt unberuehrt. Der Generator ist dieselbe
# Dockerfile-Stufe `proto-export`, die `make proto-generate` seit slice-104
# nutzt: die `.proto`-Quelle kommt per `COPY` (kein Bind-Mount) in die Stufe,
# `protoc` laeuft zur Build-Zeit; ihr `ENTRYPOINT` gibt das Erzeugnis als
# `tar`-Stream ueber stdout aus. Dieses Skript baut die Stufe und extrahiert
# host-seitig `docker run --rm --network none <image> | tar -x -C
# <Temp-Verzeichnis>` — anders als `make proto-generate` (das nach `.`
# extrahiert und damit den Baum schreibt) extrahiert dieses Gate in ein
# eigenes Temp-Verzeichnis und vergleicht dort, statt zu schreiben: ein Gate,
# das erst schreibt und dann vergleicht, laesst den Baum schmutzig zurueck
# und ist beim zweiten Lauf gruen.
#
# Frueher lief hier `docker build --target proto` gefolgt von einem
# manuellen `protoc`-Aufruf gegen einen `docker run -v <Temp>:/out --user
# <uid>:<gid>`-Bind-Mount. Der Mount war die Fehlerquelle auf
# Docker-Backends mit eingeschraenktem UID-Mapping ausserhalb des eigenen
# Host-Home (z. B. Colima mit `mounts: []` und `TMPDIR` ausserhalb `$HOME`:
# `Permission denied`, real reproduziert). Die tar-Stream-Extraktion braucht
# keinen Mount, kein `--user`-Workaround, und ist damit strukturell immun
# gegen diese Fehlerklasse.
#
# Modulpfad und `.proto`-Datei traegt die Stufe `proto-export` jetzt fest
# (geteilt mit `make proto-generate`) statt sie hier aus `go.mod`/`find`
# abzuleiten und `protoc` mit diesen Werten selbst aufzurufen. Zwei
# Eigenschaften des Bind-Mount-Mechanismus entfallen damit bewusst:
#   - Modulpfad-Cross-Check: driftet `go.mod`s `module`-Zeile vom in der
#     Stufe hartcodierten Wert weg, faellt das nicht mehr hier auf — es
#     faellt spaetestens bei `make test`/`make image` auf, weil ein
#     tatsaechlicher Modulpfad-Wechsel die Importpfade im gesamten Baum
#     bricht (derselbe Drift war vorher zusaetzlich hier sichtbar, ist aber
#     nicht mehr exklusiv hier sichtbar).
#   - Dynamische `.proto`-Dateierkennung: eine zweite `.proto`-Datei wuerde
#     von der Stufe nicht automatisch mitgeneriert (ihr `RUN`-Schritt nennt
#     die Datei namentlich) — dieselbe Einschraenkung trug `make
#     proto-generate`s Stufe bereits seit slice-104; dieses Gate zieht mit
#     dem Stufen-Wechsel nur nach, was fuer das Erzeugungsziel schon galt,
#     kein neu eingefuehrter Verlust.
# Die verbleibende `find`-Ermittlung unten ist nur noch ein Existenz-/
# Berichts-Check (Exit 2, wenn keine `.proto`-Quelle vorhanden ist) — sie
# beeinflusst nicht mehr, was generiert wird.
#
# Der Befund traegt den Diff: je abweichender Datei die erste abweichende Stelle
# im committeten Erzeugnis und darunter den Unified-Diff mit Kontext. Die Zeile
# kommt aus dem Hunk-Kopf eines `diff -U0` (ohne Kontext) — bei einer Einfuegung
# (`-N,0`) nennt der Befund N+1, weil die Abweichung erst hinter Zeile N beginnt.
# Die gepinnte Stufe laeuft ohne `--no-cache-filter`: ihren Layer-Cache kann
# kein Urteil maskieren, weil das Urteil ausserhalb der Stufe faellt — im
# Vergleich dieses Laufs.
#
# `set -o pipefail` macht die Pipe `docker run | tar -x` sicher (AGENTS.md
# §3.9 — der Exit-Code einer Pipe ist sonst der des letzten Glieds). Dieses
# Skript laeuft wie `tools/harness/proto-generate.sh` explizit unter bash
# (`make` ruft `bash tools/harness/generated-sync.sh`) statt unter dem
# `make`-Default `/bin/sh` (dort `dash`, kein `pipefail`).
#
# Aufruf: `make generated-sync` (haengt an GATE_CHECKS).
# Overrides: GENERATED_SYNC_IMAGE, GENERATED_SYNC_SOURCE_DIR.
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

PROTO_SOURCE_DIR=${GENERATED_SYNC_SOURCE_DIR:-proto}
GENERATED_SYNC_IMAGE=${GENERATED_SYNC_IMAGE:-pg-change-feed:proto-sync}

# Nur noch Existenz-/Berichts-Check (siehe Kopf-Kommentar): die eigentliche
# Erzeugung nennt ihre Quelle(n) selbst, in der Dockerfile-Stufe `proto-export`.
mapfile -t proto_sources < <(find "$PROTO_SOURCE_DIR" -type f -name '*.proto' | sort)
if [ "${#proto_sources[@]}" -eq 0 ]; then
  echo "generated-sync: FAIL — keine .proto-Quelle unter $PROTO_SOURCE_DIR" >&2
  exit 2
fi

out_dir=$(mktemp -d "${TMPDIR:-/tmp}/pg-change-feed-generated-sync.XXXXXX")
trap 'rm -rf "$out_dir"' EXIT

docker build --target proto-export -t "$GENERATED_SYNC_IMAGE" .
docker run --rm --network none "$GENERATED_SYNC_IMAGE" | tar -x -C "$out_dir"

report_diff() {
  local rel=$1 generated=$2 committed=$3 diff_text zero_ctx first_line
  diff_text=$(diff -u --label "committet: $rel" --label "Generatorausgabe: $rel" \
    "$committed" "$generated" || true)
  zero_ctx=$(diff -U0 "$committed" "$generated" || true)
  first_line=$(printf '%s\n' "$zero_ctx" | awk '
    /^@@ / {
      split($2, p, ",")
      start = p[1]; sub(/^-/, "", start)
      count = (p[2] == "" ? 1 : p[2] + 0)
      print (count == 0 ? start + 1 : start)
      exit
    }')
  printf 'generated-sync: FAIL — %s weicht von der Generatorausgabe ab (erste abweichende Stelle im committeten Erzeugnis: Zeile %s)\n' \
    "$rel" "${first_line:-unbekannt}" >&2
  printf '%s\n' "$diff_text" >&2
}

status=0

mapfile -t generated_files < <(cd "$out_dir" && find . -type f | sed 's|^\./||' | sort)

# Richtung 1: jede erzeugte Datei gegen ihren committeten Gegenpart.
for rel in "${generated_files[@]}"; do
  committed="$repo_root/$rel"
  if [ ! -f "$committed" ]; then
    printf 'generated-sync: FAIL — %s fehlt im Baum (der Generator erzeugt sie)\n' "$rel" >&2
    status=1
    continue
  fi
  if ! cmp -s "$out_dir/$rel" "$committed"; then
    report_diff "$rel" "$out_dir/$rel" "$committed"
    status=1
  fi
done

# Richtung 2: jede committete Datei, die sich als Erzeugnis dieses Generators
# kennzeichnet, muss der Lauf erzeugt haben (eine zurueckgebliebene
# Erzeugnis-Datei ohne Quelle ist ebenso eine Abweichung wie ein abweichender
# Inhalt).
while IFS= read -r committed; do
  rel=${committed#"$repo_root"/}
  if [ ! -f "$out_dir/$rel" ]; then
    printf 'generated-sync: FAIL — %s ist als Erzeugnis gekennzeichnet (protoc-gen-go), der Lauf erzeugt sie aber nicht\n' \
      "$rel" >&2
    status=1
  fi
done < <(grep -rlF --include='*.pb.go' --exclude-dir=.git \
  -e '// Code generated by protoc-gen-go' "$repo_root" | sort)

if [ "$status" -ne 0 ]; then
  echo "generated-sync: FAIL — das committete Erzeugnis stimmt nicht mit der Ausgabe des gepinnten Generators ueberein" >&2
  exit 1
fi

printf 'generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)\n'
printf 'generated-sync:   Quelle: %s\n' "${proto_sources[@]}"
printf 'generated-sync:   geprueft: %s\n' "${generated_files[@]}"
