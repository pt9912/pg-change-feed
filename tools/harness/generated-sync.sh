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
# Der Arbeitsbaum bleibt unberuehrt. Der Generator ist die Dockerfile-Stufe
# `proto` (protoc 31.1-r1, protoc-gen-go v1.36.12, protoc-gen-go-grpc v1.6.2 —
# gepinnt); er schreibt in ein Temp-Verzeichnis, der Baum haengt als
# `:ro`-Bind-Mount im Container. Die Modul-Layout-Relation stellt dieselbe
# Option her, die `make proto-generate` setzt (`--go_opt=module=<Modulpfad>`):
# der Generator legt seine Ausgabe damit unter dem Temp-Wurzel genau unter den
# Pfaden ab, unter denen das Erzeugnis im Baum liegt. Das Ziel
# `make proto-generate` schreibt in den Bind-Mount des Baums und ist deshalb
# kein Pruef-Schritt — ein Gate, das erst schreibt und dann vergleicht, laesst
# den Baum schmutzig zurueck und ist beim zweiten Lauf gruen.
#
# Der Befund traegt den Diff: je abweichender Datei die erste abweichende Stelle
# im committeten Erzeugnis und darunter den Unified-Diff mit Kontext. Die Zeile
# kommt aus dem Hunk-Kopf eines `diff -U0` (ohne Kontext) — bei einer Einfuegung
# (`-N,0`) nennt der Befund N+1, weil die Abweichung erst hinter Zeile N
# beginnt. Die gepinnte Stufe laeuft ohne
# `--no-cache-filter`: ihren Layer-Cache kann kein Urteil maskieren, weil das
# Urteil ausserhalb der Stufe faellt — im Vergleich dieses Laufs.
#
# Aufruf: `make generated-sync` (haengt an GATE_CHECKS).
# Overrides: GENERATED_SYNC_IMAGE, GENERATED_SYNC_SOURCE_DIR,
# GENERATED_SYNC_MODULE, GENERATED_SYNC_RUN_USER.
set -euo pipefail

repo_root=$(git rev-parse --show-toplevel)
cd "$repo_root"

PROTO_SOURCE_DIR=${GENERATED_SYNC_SOURCE_DIR:-proto}
GENERATED_SYNC_IMAGE=${GENERATED_SYNC_IMAGE:-pg-change-feed:proto-sync}
# Der Modulpfad steht in go.mod und nirgends sonst: er ist der Praefix, den
# `--go_opt=module=…` von den Ausgabepfaden abschneidet. Stimmt er nicht mit dem
# go_package-Praefix der .proto ueberein, bricht der Generator ab (`generated
# file does not match prefix …`) — die Ausgabe landet dann auf keinem falschen
# Pfad, sondern der Lauf ist an dieser Stelle rot.
GENERATED_SYNC_MODULE=${GENERATED_SYNC_MODULE:-$(awk '$1 == "module" { print $2; exit }' go.mod || true)}
# Der Generator laeuft als Aufrufer-uid, damit er in das Temp-Verzeichnis
# schreiben darf (dasselbe Muster wie PROTO_RUN_USER im Makefile).
RUN_USER=${GENERATED_SYNC_RUN_USER:-$(id -u):$(id -g)}

if [ -z "$GENERATED_SYNC_MODULE" ]; then
  echo "generated-sync: FAIL — kein Modulpfad (go.mod ohne module-Zeile, GENERATED_SYNC_MODULE leer)" >&2
  exit 2
fi

# Die Quellen sind die .proto-Dateien unter dem Quellverzeichnis; eine neue
# Datei tritt damit von selbst in den Lauf ein.
mapfile -t proto_sources < <(find "$PROTO_SOURCE_DIR" -type f -name '*.proto' | sort)
if [ "${#proto_sources[@]}" -eq 0 ]; then
  echo "generated-sync: FAIL — keine .proto-Quelle unter $PROTO_SOURCE_DIR" >&2
  exit 2
fi

out_dir=$(mktemp -d "${TMPDIR:-/tmp}/pg-change-feed-generated-sync.XXXXXX")
trap 'rm -rf "$out_dir"' EXIT

docker build --target proto -t "$GENERATED_SYNC_IMAGE" .
docker run --rm --user "$RUN_USER" --network none \
  -v "$repo_root":/src:ro -v "$out_dir":/out -w /src "$GENERATED_SYNC_IMAGE" \
  protoc -I "$PROTO_SOURCE_DIR" \
    --go_out=/out --go_opt=module="$GENERATED_SYNC_MODULE" \
    --go-grpc_out=/out --go-grpc_opt=module="$GENERATED_SYNC_MODULE" \
    "${proto_sources[@]}"

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

printf 'generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)\n'
printf 'generated-sync:   Quelle: %s\n' "${proto_sources[@]}"
printf 'generated-sync:   geprueft: %s\n' "${generated_files[@]}"
