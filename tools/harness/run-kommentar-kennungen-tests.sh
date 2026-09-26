#!/usr/bin/env bash
# run-kommentar-kennungen-tests.sh — Tabellentest gegen den Aufrufer
# kommentar-kennungen.sh (harness/sensors/kommentar-kennungen.md): das Programm
# hat seinen Go-Test, dieser Test bindet den Aufrufer: Eingabefehler (Arg-Zahl,
# fehlende TOOLCHAIN_IMAGE, unbekannte Diff-Basis) · Exit-Weitergabe · jedes Wort
# von PATHS/DIFF ein eigenes Argument (die Marker-Datei eines eingeschleusten
# Kommandos entsteht nicht) · gepinnte Form des Diff-Stroms unter fremder
# Git-Konfiguration (Strom byte-gleich zur Konfiguration ohne Eingriff) · keine
# Temp-Datei zurück · vier Läufe mit echtem Docker gegen ein Wegwerf-Repo
# (Kandidat im Diff-Modus mit und ohne `diff.mnemonicPrefix`, TESTS-Wert und Pfad
# als Eingabefehler des Programms).
# Die Fälle mit Stub ersetzen `docker` durch ein Skript, das seine Argumente und
# stdin festhält; die vier echten Läufe brauchen Docker und TOOLCHAIN_IMAGE
# (das Makefile setzt es), netzlos (`--network none`). Der Prüfling ist per TOOL
# übersteuerbar (Mutationsläufe gegen eine Kopie).
set -uo pipefail
repo=$(git rev-parse --show-toplevel)
tool=${TOOL:-$repo/tools/harness/kommentar-kennungen.sh}
image=${TOOLCHAIN_IMAGE:?TOOLCHAIN_IMAGE fehlt (make test-kommentar-kennungen setzt es)}

tmp=$(mktemp -d "${TMPDIR:-/tmp}/kommentar-kennungen-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/bin" "$tmp/stub" "$tmp/tmpdir" "$tmp/repo"

# Stub für docker: hält Argumente (eins je Zeile) und stdin fest, endet mit STUB_RC.
cat >"$tmp/bin/docker" <<'STUB'
#!/usr/bin/env bash
printf '%s\n' "$@" >"$STUB_DIR/args"
cat >"$STUB_DIR/stdin"
exit "${STUB_RC:-0}"
STUB
chmod +x "$tmp/bin/docker"

# Wegwerf-Repo: der Parent trägt einen Kandidaten (Zeilen 3-4); der Arbeitsbaum
# fügt einen zweiten Kandidaten hinzu (Zeilen 7-8). Das Programm und go.mod
# liegen darin, damit der echte Docker-Lauf es bauen kann.
cd "$tmp/repo" || exit 1
git init -q .
git config user.name test
git config user.email test@example.invalid
git config commit.gpgsign false
mkdir -p a tools/harness/kommentar-kennungen
cp "$repo/go.mod" go.mod
cp "$repo/tools/harness/kommentar-kennungen/main.go" tools/harness/kommentar-kennungen/main.go
printf 'package a\n\n// ADR-0001 und\n// LH-FA-CAP-001\nfunc F() {}\n' >a/x.go
git add -A
git commit -q -m parent
parent=$(git rev-parse --short HEAD)
printf 'package a\n\n// ADR-0001 und\n// LH-FA-CAP-001\nfunc F() {}\n\n// SPEC-003 und\n// ARC-002\nfunc G() {}\n' >a/x.go
printf 'a\n' >nichts.txt # Datei für den Glob-Fall

fail=0
out=""
rc=0

# call <Pfade> <COUNT> <TESTS> <DIFF>: Aufrufer mit Stub-docker; Umgebung über
# die Variable ENV_EXTRA (Wörter der Form NAME=Wert).
ENV_EXTRA=()
call() {
  rm -f "$tmp/stub/args" "$tmp/stub/stdin"
  out=$(env PATH="$tmp/bin:$PATH" STUB_DIR="$tmp/stub" TMPDIR="$tmp/tmpdir" TOOLCHAIN_IMAGE=stub \
    ${ENV_EXTRA[@]+"${ENV_EXTRA[@]}"} bash "$tool" "$@" 2>&1)
  rc=$?
}
# real <Pfade> <COUNT> <TESTS> <DIFF>: Aufrufer mit echtem Docker.
real() {
  out=$(env TMPDIR="$tmp/tmpdir" TOOLCHAIN_IMAGE="$image" \
    ${ENV_EXTRA[@]+"${ENV_EXTRA[@]}"} bash "$tool" "$@" 2>&1)
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
check() { # <Fallname> <Bedingung als Befehl ...>
  local name=$1
  shift
  if ! "$@"; then
    echo "FEHLER: $name" >&2
    fail=1
  fi
}
docker_not_called() { [ ! -e "$tmp/stub/args" ]; }
tmpdir_empty() { [ -z "$(ls -A "$tmp/tmpdir")" ]; }

# Fall 1 — Eingabefehler des Aufrufers: Argumentzahl, fehlende TOOLCHAIN_IMAGE,
# unbekannte Diff-Basis; jede endet mit Exit 2, ohne dass docker läuft.
call a b c
expect "Argumentzahl" 2 'Aufruf:'
check "Argumentzahl: docker lief" docker_not_called
out=$(env -u TOOLCHAIN_IMAGE PATH="$tmp/bin:$PATH" STUB_DIR="$tmp/stub" bash "$tool" "" "" "" "" 2>&1)
rc=$?
expect "TOOLCHAIN_IMAGE fehlt" 2 'TOOLCHAIN_IMAGE fehlt'
call "" "" "" nicht-vorhanden
expect "unbekannte Diff-Basis" 2 "kein Commit"
check "unbekannte Diff-Basis: docker lief" docker_not_called
check "unbekannte Diff-Basis: Temp-Datei" tmpdir_empty

# Fall 2 — Exit-Weitergabe: der Exit des Programms erreicht den Aufrufer
# unverändert (0, 1, 2 bleiben unterscheidbar).
for want in 0 1 2; do
  STUB_RC=$want
  export STUB_RC
  call "" "" "" ""
  expect "Exit $want weitergegeben" "$want" -
done
unset STUB_RC

# Fall 3 — jedes Wort von PATHS ein eigenes Argument, keine Expansion, kein
# Kommando: die Marker-Dateien der eingeschleusten Kommandos entstehen nicht.
call "x; touch $tmp/marker1 \$(touch $tmp/marker2) \`touch $tmp/marker3\` *" "" "" ""
expect "PATHS mit Kommando-Syntax" 0 -
check "PATHS: Marker 1 entstand" test ! -e "$tmp/marker1"
check "PATHS: Marker 2 entstand" test ! -e "$tmp/marker2"
check "PATHS: Marker 3 entstand" test ! -e "$tmp/marker3"
check "PATHS: 'x;' ist ein eigenes Argument" grep -Fxq -- 'x;' "$tmp/stub/args"
check "PATHS: 'touch' ist ein eigenes Argument" grep -Fxq -- 'touch' "$tmp/stub/args"
check "PATHS: '*' bleibt unexpandiert" grep -Fxq -- '*' "$tmp/stub/args"
check "PATHS: nichts.txt wurde nicht expandiert" bash -c '! grep -Fxq -- nichts.txt "$1"' _ "$tmp/stub/args"
call "" "" "" "--output=$tmp/diffout"
expect "DIFF als Option" 2 'kein Commit'
check "DIFF als Option: Datei entstand" test ! -e "$tmp/diffout"
call "" "" "" "$parent; touch $tmp/marker4"
expect "DIFF mit Kommando-Syntax" 2 'kein Commit'
check "DIFF: Marker 4 entstand" test ! -e "$tmp/marker4"

# Fall 4 — gepinnte Form des Diff-Stroms: unter fremder Git-Konfiguration ist der
# Strom byte-gleich zu dem ohne Eingriff, und die Zieldatei-Zeile trägt `b/`.
call "" "" "" "$parent"
expect "Diff-Strom (Vergleichslauf)" 0 -
cp "$tmp/stub/stdin" "$tmp/stream-plain"
check "Diff-Strom trägt die Zieldatei-Zeile" grep -Fxq -- '+++ b/a/x.go' "$tmp/stream-plain"
check "Diff-Strom trägt die Quelldatei-Zeile" grep -Fxq -- '--- a/a/x.go' "$tmp/stream-plain"
for kv in diff.mnemonicPrefix=true diff.noprefix=true diff.dstPrefix=X/ diff.srcPrefix=Y/ \
  diff.external=/bin/false color.diff=always diff.algorithm=patience; do
  ENV_EXTRA=(GIT_CONFIG_COUNT=1 "GIT_CONFIG_KEY_0=${kv%%=*}" "GIT_CONFIG_VALUE_0=${kv#*=}")
  call "" "" "" "$parent"
  expect "Diff-Strom unter $kv" 0 -
  check "Diff-Strom unter $kv weicht ab" cmp -s "$tmp/stub/stdin" "$tmp/stream-plain"
done
ENV_EXTRA=()

# Fall 5 — ein fehlschlagendes git diff endet mit Exit 2, ohne dass docker läuft
# (ein ungültiger Wert von diff.algorithm lässt git scheitern).
ENV_EXTRA=(GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=diff.algorithm GIT_CONFIG_VALUE_0=gibt-es-nicht)
call "" "" "" "$parent"
expect "git diff schlägt fehl" 2 'schlug fehl'
check "git diff schlägt fehl: docker lief" docker_not_called
ENV_EXTRA=()

# Fall 6 — keine Temp-Datei bleibt zurück, weder nach Exit 0 noch nach 1 und 2.
tmpdir_empty_after() { # <Exit des Stubs>
  STUB_RC=$1
  export STUB_RC
  call "" "" "" "$parent"
  unset STUB_RC
  tmpdir_empty
}
for want in 0 1 2; do
  check "Temp-Datei nach Exit $want" tmpdir_empty_after "$want"
done

# Fall 7 — echter Docker-Lauf, Diff-Modus: nur der hinzugefügte Kandidat
# (Zeilen 7-8) wird gemeldet, der Bestandskandidat (3-4) nicht; dasselbe unter
# `diff.mnemonicPrefix=true`.
real "" 1 "" "$parent"
expect "echter Lauf: Zahl" 0 '^1$'
real "" "" "" "$parent"
expect "echter Lauf: Kandidat" 1 '^a/x\.go:7-8  SPEC-003, ARC-002$'
ENV_EXTRA=(GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=diff.mnemonicPrefix GIT_CONFIG_VALUE_0=true)
real "" 1 "" "$parent"
expect "echter Lauf unter diff.mnemonicPrefix: Zahl" 0 '^1$'
ENV_EXTRA=()

# Fall 8 — echter Docker-Lauf, Eingabefehler des Programms: ein ungültiger
# TESTS-Wert und ein fehlender Pfad enden mit Exit 2 und einer Meldung.
real "" "" alle ""
expect "TESTS=alle" 2 "weder all, exclude noch only"
real "gibt-es-nicht" "" "" ""
expect "fehlender Pfad" 2 'gibt-es-nicht'

check "Temp-Datei nach allen Läufen" tmpdir_empty

if [ "$fail" -eq 0 ]; then
  echo "run-kommentar-kennungen-tests: alle Fälle bestanden"
else
  echo "run-kommentar-kennungen-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
