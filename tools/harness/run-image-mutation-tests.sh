#!/usr/bin/env bash
# run-image-mutation-tests.sh — Tabellentest gegen den Aufrufer
# tools/harness/image-mutation.sh (harness/targets/image-mutation.md), Stub-`docker`
# in einem Wegwerf-Repo im Temp-Verzeichnis (netzlos, kein echter Bau). Fälle:
# gültige Argumente inkl. SRC mit Leerzeichen · `harness/image-hash.txt`/`.raw`
# im Wegwerf-Repo bleiben unverändert · TAG-Zeichenklasse und die reservierten
# Namen `dev`/`latest` · SRC fehlt/kein Verzeichnis/ohne `Dockerfile`/ohne
# `go.mod`/Repo-Wurzel/Verzeichnis unter der Wurzel · Exit-Weitergabe eines
# Docker-Fehlers · `rm` mit genau einem `rmi`-Argument · Make-Ebene ohne SRC
# bzw. TAG (kein Docker-Aufruf, `$(error …)` bricht vor dem Rezept ab). Der
# Prüfling ist per TOOL übersteuerbar (Mutationsläufe gegen eine Kopie).
set -uo pipefail
repo=$(git rev-parse --show-toplevel)
tool=${TOOL:-$repo/tools/harness/image-mutation.sh}

tmp=$(mktemp -d "${TMPDIR:-/tmp}/image-mutation-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT

# Wegwerf-Repo: eigene Repo-Wurzel für den Wurzel-Vergleich, mit den Fühler-
# Dateien harness/image-hash.txt und .raw (sie dürfen der Lauf nicht anfassen).
wrepo="$tmp/repo"
mkdir -p "$wrepo/harness" "$wrepo/sub"
git -C "$wrepo" init -q
printf 'sha256:vorher-unveraendert\n' >"$wrepo/harness/image-hash.txt"
printf '{"fuehler":"vorher"}\n' >"$wrepo/harness/image-hash.raw"
printf 'FROM scratch\n' >"$wrepo/Dockerfile"
printf 'module wrepo\n' >"$wrepo/go.mod"
printf 'FROM scratch\n' >"$wrepo/sub/Dockerfile"
printf 'module wrepo-sub\n' >"$wrepo/sub/go.mod"

# Gültige SRC-Kopie außerhalb der Wegwerf-Wurzel, mit Leerzeichen im Pfad.
srcdir="$tmp/mutcopy mit leerzeichen"
mkdir -p "$srcdir"
printf 'FROM scratch\n' >"$srcdir/Dockerfile"
printf 'module mutcopy\n' >"$srcdir/go.mod"
src_abs=$(realpath -- "$srcdir")

# Stub-`docker`: hält die Argumente fest (eine Zeile je Argument) und endet
# mit DOCKER_EXIT (Default 0). Die Zusagen betreffen die Argumente, mit denen
# der Aufrufer Docker ruft; dass der Daemon sie einhält, belegt dieser Test
# nicht (Nachbild: run-fmt-check-tests.sh Fall 11).
mkdir -p "$tmp/bin" "$tmp/stub"
cat >"$tmp/bin/docker" <<'STUB'
#!/usr/bin/env bash
printf '%s\n' "$@" >"$STUB_DIR/args"
exit "${DOCKER_EXIT:-0}"
STUB
chmod +x "$tmp/bin/docker"

fail=0
out=""
rc=0

# run <Argumente…>: der Aufrufer läuft mit dem Arbeitsverzeichnis der
# Wegwerf-Wurzel (git rev-parse --show-toplevel liest die aktuelle Wurzel)
# und dem Stub-`docker` vor dem echten `docker` auf dem PATH.
run() {
  rm -f "$tmp/stub/args"
  out=$(cd "$wrepo" && env PATH="$tmp/bin:$PATH" STUB_DIR="$tmp/stub" DOCKER_EXIT="${DOCKER_EXIT:-0}" bash "$tool" "$@" 2>&1)
  rc=$?
}
expect() { # <Fallname> <erwarteter Exit> <Muster in der Ausgabe oder ->
  local name=$1 want_rc=$2 pattern=$3
  if [ "$rc" -ne "$want_rc" ]; then
    echo "FEHLER: $name — Exit $rc, erwartet $want_rc; Ausgabe:" >&2
    printf '%s\n' "$out" >&2
    fail=1
  elif [ "$pattern" != "-" ] && ! printf '%s\n' "$out" | grep -qE -- "$pattern"; then
    echo "FEHLER: $name — Ausgabe trägt '$pattern' nicht:" >&2
    printf '%s\n' "$out" >&2
    fail=1
  fi
}
no_docker_call() { # <Fallname>: der Stub wurde nicht aufgerufen
  if [ -f "$tmp/stub/args" ]; then
    echo "FEHLER: $1 — Docker wurde trotz Eingabefehler aufgerufen:" >&2
    cat "$tmp/stub/args" >&2
    fail=1
  fi
}
args_line() { # <Fallname> <Zeilennummer 1-basiert> <erwartete Zeile>
  local got
  got=$(sed -n "${2}p" "$tmp/stub/args" 2>/dev/null || true)
  if [ "$got" != "$3" ]; then
    echo "FEHLER: $1 — Zeile $2 der Docker-Argumente ist '$got', erwartet '$3':" >&2
    cat "$tmp/stub/args" >&2
    fail=1
  fi
}
args_count() { # <Fallname> <erwartete Zeilenzahl>
  local n
  n=$(wc -l <"$tmp/stub/args" 2>/dev/null || echo 0)
  if [ "$n" -ne "$2" ]; then
    echo "FEHLER: $1 — Docker-Aufruf trägt $n Argumente, erwartet $2:" >&2
    cat "$tmp/stub/args" >&2
    fail=1
  fi
}

tag="probe-tag"
imgref="pg-change-feed-mutation:$tag"

# Fall 1 — gültige Argumente (SRC mit Leerzeichen): Exit 0, genau die sechs
# Docker-Argumente von `docker buildx build --load -t <repo>:<tag> <SRC>`,
# nichts sonst (kein --metadata-file, --push, --platform, --build-arg).
run build "$srcdir" "$tag"
expect "gültige Argumente: Lauf" 0 "$imgref"
args_count "gültige Argumente" 6
args_line "gültige Argumente" 1 "buildx"
args_line "gültige Argumente" 2 "build"
args_line "gültige Argumente" 3 "--load"
args_line "gültige Argumente" 4 "-t"
args_line "gültige Argumente" 5 "$imgref"
args_line "gültige Argumente" 6 "$src_abs"

# Fall 2 — harness/image-hash.txt und .raw im Wegwerf-Repo bleiben nach dem
# Lauf byte-gleich (Fühler gegen "-t auf :dev" oder "schreibt die Hash-Datei").
if [ "$(cat "$wrepo/harness/image-hash.txt")" != "sha256:vorher-unveraendert" ]; then
  echo "FEHLER: harness/image-hash.txt wurde verändert" >&2
  fail=1
fi
if [ "$(cat "$wrepo/harness/image-hash.raw")" != '{"fuehler":"vorher"}' ]; then
  echo "FEHLER: harness/image-hash.raw wurde verändert" >&2
  fail=1
fi

# Fall 3 — TAG-Eingabefehler: Zeichenklasse, reservierte Namen; jeder endet
# mit Exit 2, ohne dass Docker aufgerufen wird.
for bad_tag in "" dev latest Foo "a b" "-x" "a:b" "../x"; do
  run build "$srcdir" "$bad_tag"
  expect "TAG '$bad_tag'" 2 -
  no_docker_call "TAG '$bad_tag'"
done

# Fall 4 — SRC-Eingabefehler: fehlt, kein Verzeichnis, ohne Dockerfile, ohne
# go.mod, Repo-Wurzel, Verzeichnis unter der Wurzel; jeder endet mit Exit 2,
# ohne Docker-Aufruf.
run build "" "$tag"
expect "SRC fehlt (leer)" 2 -
no_docker_call "SRC fehlt (leer)"

run build "$tag"
expect "SRC fehlt (Argumentzahl)" 2 -
no_docker_call "SRC fehlt (Argumentzahl)"

run build "$tmp/gibt-es-nicht" "$tag"
expect "SRC kein Verzeichnis" 2 'kein Verzeichnis'
no_docker_call "SRC kein Verzeichnis"

mkdir -p "$tmp/ohne-dockerfile"
printf 'module x\n' >"$tmp/ohne-dockerfile/go.mod"
run build "$tmp/ohne-dockerfile" "$tag"
expect "SRC ohne Dockerfile" 2 'Dockerfile'
no_docker_call "SRC ohne Dockerfile"

mkdir -p "$tmp/ohne-gomod"
printf 'FROM scratch\n' >"$tmp/ohne-gomod/Dockerfile"
run build "$tmp/ohne-gomod" "$tag"
expect "SRC ohne go.mod" 2 'go\.mod'
no_docker_call "SRC ohne go.mod"

run build "$wrepo" "$tag"
expect "SRC ist Repo-Wurzel" 2 'Repo-Wurzel'
no_docker_call "SRC ist Repo-Wurzel"

run build "$wrepo/sub" "$tag"
expect "SRC unter der Wurzel" 2 'Repo-Wurzel'
no_docker_call "SRC unter der Wurzel"

# Fall 5 — Exit-Weitergabe eines Docker-Fehlers: Exit 1, die Zeile nennt den
# Exit-Code des Aufrufs.
DOCKER_EXIT=17 run build "$srcdir" "$tag"
expect "Docker-Fehler (build)" 1 'Exit 17'
DOCKER_EXIT=0

# Fall 6 — rm mit genau einem rmi-Argument, kein -f, kein prune.
run rm "$tag"
expect "rm: Lauf" 0 "$imgref"
args_count "rm" 2
args_line "rm" 1 "rmi"
args_line "rm" 2 "$imgref"

DOCKER_EXIT=5 run rm "$tag"
expect "Docker-Fehler (rm)" 1 'Exit 5'
DOCKER_EXIT=0

run rm ""
expect "rm: TAG fehlt (leer)" 2 -
no_docker_call "rm: TAG fehlt (leer)"

run rm
expect "rm: TAG fehlt (Argumentzahl)" 2 -
no_docker_call "rm: TAG fehlt (Argumentzahl)"

# Fall 7 — unbekanntes Verb: Exit 2, kein Docker-Aufruf.
run wat "$srcdir" "$tag"
expect "unbekanntes Verb" 2 'Aufruf:'
no_docker_call "unbekanntes Verb"

# Fall 8 — Make-Ebene ohne SRC bzw. TAG: `$(error …)` bricht vor dem Rezept
# ab, kein Docker-Aufruf (echtes Makefile, ohne Stub nötig — der Fehler
# feuert vor jedem Befehl).
out=$(cd "$repo" && env -u SRC -u TAG make image-mutation 2>&1)
rc=$?
expect "make image-mutation ohne SRC" 2 'SRC fehlt'
out=$(cd "$repo" && env -u TAG make image-mutation SRC="$srcdir" 2>&1)
rc=$?
expect "make image-mutation ohne TAG" 2 'TAG fehlt'
out=$(cd "$repo" && env -u TAG make image-mutation-rm 2>&1)
rc=$?
expect "make image-mutation-rm ohne TAG" 2 'TAG fehlt'

if [ "$fail" -eq 0 ]; then
  echo "run-image-mutation-tests: alle Fälle bestanden"
else
  echo "run-image-mutation-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
