#!/usr/bin/env bash
# run-fmt-check-tests.sh — Tabellentest gegen den Aufrufer fmt-check.sh
# (harness/sensors/fmt-check.md), echter Docker-Lauf gegen Wegwerf-Verzeichnisse
# im Temp-Verzeichnis (der Container des Werkzeugs läuft mit `--network none`; der
# Fall Docker-Fehler löst einen Image-Zugriff des Docker-Daemons aus). Fälle:
# formatierte Datei · unformatierte Datei (genannt, Exit 1) · nur die
# unformatierte von zwei Dateien genannt · Datei in einem Unterverzeichnis ·
# Datei mit Syntaxfehler (Exit 2) · Verzeichnis ohne Go-Datei (Exit 2) · Eingabe
# bleibt byte-gleich (lesender Mount) · Pfad mit Leerzeichen · `.go`-Symlink
# (gezählt und genannt) · Eingabefehler des Aufrufers (Argumentzahl, fehlendes
# Verzeichnis, fehlende TOOLCHAIN_IMAGE) · nicht auflösbares Image (Docker-Fehler,
# Exit 2) · Argumente des Docker-Aufrufs (`--network none`, Mount `:ro`, Image;
# ein Stub-`docker` hält sie fest). Der Prüfling ist per TOOL übersteuerbar
# (Mutationsläufe gegen eine Kopie); TOOLCHAIN_IMAGE setzt das Makefile.
set -uo pipefail
repo=$(git rev-parse --show-toplevel)
tool=${TOOL:-$repo/tools/harness/fmt-check.sh}
image=${TOOLCHAIN_IMAGE:?TOOLCHAIN_IMAGE fehlt (make test-fmt-check setzt es)}

tmp=$(mktemp -d "${TMPDIR:-/tmp}/fmt-check-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT

fail=0
out=""
rc=0

# run <Verzeichnis>: Aufrufer mit echtem Docker; stdout und stderr zusammen.
run() {
  out=$(env TOOLCHAIN_IMAGE="$image" bash "$tool" "$@" 2>&1)
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
refuse() { # <Fallname> <Muster, das die Ausgabe nicht tragen darf>
  if printf '%s\n' "$out" | grep -qE -- "$2"; then
    echo "FEHLER: $1 — Ausgabe trägt '$2':" >&2
    printf '%s\n' "$out" >&2
    fail=1
  fi
}

formatted='package a\n\nfunc F() {}\n'
unformatted='package a\n\nfunc   F( ) {\n}\n'
syntaxerr='package a\n\nfunc F( {\n'

# Fall 1 — formatierte Datei: Exit 0, die Zählung nennt eine Datei.
d="$tmp/ok"
mkdir -p "$d"
printf "$formatted" >"$d/a.go"
run "$d"
expect "formatierte Datei" 0 '1 Go-Dateien geprüft, alle formatiert'

# Fall 2 — unformatierte Datei: Exit 1, die Datei wird genannt.
d="$tmp/unformatiert"
mkdir -p "$d"
printf "$unformatted" >"$d/u.go"
run "$d"
expect "unformatierte Datei" 1 '^u\.go$'

# Fall 3 — zwei Dateien, eine unformatiert: nur diese wird genannt.
d="$tmp/gemischt"
mkdir -p "$d"
printf "$formatted" >"$d/gut.go"
printf "$unformatted" >"$d/schlecht.go"
run "$d"
expect "gemischt: unformatierte genannt" 1 '^schlecht\.go$'
refuse "gemischt: formatierte genannt" '^gut\.go$'

# Fall 4 — Datei in einem Unterverzeichnis (die Wurzel trägt nur eine
# formatierte Datei): Exit 1, der Pfad trägt das Unterverzeichnis.
d="$tmp/unterverzeichnis"
mkdir -p "$d/sub/tief"
printf "$formatted" >"$d/a.go"
printf "$unformatted" >"$d/sub/tief/x.go"
run "$d"
expect "Unterverzeichnis" 1 '^sub/tief/x\.go$'

# Fall 5 — Datei mit Syntaxfehler: Exit 2, kein Exit 0 und kein Exit 1.
d="$tmp/syntax"
mkdir -p "$d"
printf "$syntaxerr" >"$d/s.go"
run "$d"
expect "Syntaxfehler" 2 -

# Fall 6 — Verzeichnis ohne Go-Datei: leer ist nicht bestanden (Exit 2), auch
# mit einer Nicht-Go-Datei und mit einer Datei, deren Name mit einem Punkt beginnt.
d="$tmp/leer"
mkdir -p "$d"
run "$d"
expect "leeres Verzeichnis" 2 'keine Go-Datei'
printf 'text\n' >"$d/notiz.txt"
printf "$unformatted" >"$d/.versteckt.go"
run "$d"
expect "Verzeichnis ohne Go-Datei" 2 'keine Go-Datei'

# Fall 7 — die Eingabe bleibt byte-gleich (lesender Mount, kein Schreibpfad).
d="$tmp/unveraendert"
mkdir -p "$d"
printf "$unformatted" >"$d/u.go"
cp "$d/u.go" "$tmp/u.go.vorher"
run "$d"
expect "Eingabe unverändert: Lauf" 1 '^u\.go$'
if ! cmp -s "$d/u.go" "$tmp/u.go.vorher"; then
  echo "FEHLER: die Eingabe wurde umgeschrieben" >&2
  fail=1
fi

# Fall 8 — Pfad mit Leerzeichen: ein Verzeichnis-Argument, ein Mount.
d="$tmp/mit leerzeichen"
mkdir -p "$d"
printf "$unformatted" >"$d/u.go"
run "$d"
expect "Pfad mit Leerzeichen" 1 '^u\.go$'

# Fall 8b — ein `.go`-Symlink: die Zählung nimmt ihn wie `gofmt`, der Lauf nennt
# ihn (das Ziel des Links trägt keine `.go`-Endung, der Link ist die einzige Go-Datei).
d="$tmp/symlink"
mkdir -p "$d"
printf "$unformatted" >"$d/ziel.txt"
ln -s ziel.txt "$d/l.go"
run "$d"
expect "Symlink mit Endung .go" 1 '^l\.go$'

# Fall 9 — Eingabefehler des Aufrufers: zu viele Argumente, kein Verzeichnis,
# fehlende TOOLCHAIN_IMAGE; jeder endet mit Exit 2.
run "$tmp/ok" "$tmp/ok"
expect "zu viele Argumente" 2 'Aufruf:'
run "$tmp/gibt-es-nicht"
expect "kein Verzeichnis" 2 'kein Verzeichnis'
out=$(env -u TOOLCHAIN_IMAGE bash "$tool" "$tmp/ok" 2>&1)
rc=$?
expect "TOOLCHAIN_IMAGE fehlt" 2 'TOOLCHAIN_IMAGE fehlt'

# Fall 10 — Docker-Fehler (das Image ist nicht auflösbar): Exit 2.
out=$(env TOOLCHAIN_IMAGE=fmt-check-test.invalid/gibt-es-nicht:0 bash "$tool" "$tmp/ok" 2>&1)
rc=$?
expect "Docker-Fehler" 2 'Exit'

# Fall 11 — Argumente des Docker-Aufrufs: ein Stub-`docker` hält sie fest (eine je
# Zeile) und endet mit Exit 0. Der Container läuft ohne Netz (`--network none`),
# das Verzeichnis liegt lesend unter /src (`<absoluter Pfad>:/src:ro`), das Image
# ist der Wert von TOOLCHAIN_IMAGE. Die Zusagen betreffen die Argumente, mit denen
# der Aufrufer Docker ruft; dass der Daemon sie einhält, belegt dieser Fall nicht.
mkdir -p "$tmp/bin" "$tmp/stub"
cat >"$tmp/bin/docker" <<'STUB'
#!/usr/bin/env bash
printf '%s\n' "$@" >"$STUB_DIR/args"
exit 0
STUB
chmod +x "$tmp/bin/docker"
d="$tmp/mit leerzeichen"
abs=$(realpath -- "$d")
out=$(env PATH="$tmp/bin:$PATH" STUB_DIR="$tmp/stub" TOOLCHAIN_IMAGE=stub-image:1 bash "$tool" "$d" 2>&1)
rc=$?
expect "Stub-docker: Lauf" 0 -
args_have() { grep -Fxq -- "$1" "$tmp/stub/args"; }
args_check() { # <Fallname> <Muster als ganze Zeile>
  if ! args_have "$2"; then
    echo "FEHLER: $1 — die Argumente des Docker-Aufrufs tragen '$2' nicht:" >&2
    cat "$tmp/stub/args" >&2
    fail=1
  fi
}
args_check "Docker-Argument: --network" '--network'
if ! grep -A1 -Fx -- '--network' "$tmp/stub/args" | grep -Fxq -- 'none'; then
  echo "FEHLER: Docker-Argument: --network trägt nicht den Wert none:" >&2
  cat "$tmp/stub/args" >&2
  fail=1
fi
args_check "Docker-Argument: lesender Mount" "$abs:/src:ro"
args_check "Docker-Argument: Image" 'stub-image:1'

if [ "$fail" -eq 0 ]; then
  echo "run-fmt-check-tests: alle Fälle bestanden"
else
  echo "run-fmt-check-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
