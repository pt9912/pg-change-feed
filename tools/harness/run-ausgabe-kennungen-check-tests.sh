#!/usr/bin/env bash
# run-ausgabe-kennungen-check-tests.sh — Tabellentest gegen
# ausgabe-kennungen-check.sh: je Kennungsart ein Treffer im Go-Literal und im
# echo, Literal mit Schraegstrich/URL davor, Raw-String, Printf; Nicht-Treffer
# (Kommentarzeile, nachgestellter Kommentar ohne Anfuehrungszeichen, Struct-Tag,
# Test-Datei, Laeufer-Verzeichnis); der benannte Falsch-Positiv (nachgestellter
# Kommentar mit zitiertem Literal) und die benannte Grenze (mehrzeiliges
# Raw-String-Literal); die Lesefehler (Exit 2) und die saubere Flaeche
# (Exit 0). Je Exit-Klasse prueft mindestens ein Fall den Meldungstext.
# Netzlos, kein Docker noetig.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0
cases=0
tmp=$(mktemp -d)
: "${tmp:?}"
trap 'chmod -R u+rwx "$tmp" 2>/dev/null; rm -rf "${tmp:?}"' EXIT

CLEAN_GO='package x

func f() { println("Der Betreiber startet den Feed-Container.") }
'
CLEAN_SH='echo "schema-rollout: Vorlauf - View-Signatur-Aenderung"'

# fresh: Wurzel mit je einer sauberen Go-Datei je Wurzel und je einem Skript
fresh() {
  chmod -R u+rwx "${tmp:?}/root" 2>/dev/null
  rm -rf "${tmp:?}/root"
  mkdir -p "$tmp/root/internal/a" "$tmp/root/cmd/b" "$tmp/root/tools/schema/c" \
           "$tmp/root/examples" "$tmp/root/tools/harness"
  printf '%s\n' "$CLEAN_GO" > "$tmp/root/internal/a/a.go"
  printf '%s\n' "$CLEAN_GO" > "$tmp/root/cmd/b/b.go"
  printf '%s\n' "$CLEAN_GO" > "$tmp/root/tools/schema/c/c.go"
  printf '%s\n' "$CLEAN_SH" > "$tmp/root/tools/schema/rollout.sh"
  printf '%s\n' "$CLEAN_SH" > "$tmp/root/examples/demo.sh"
}

# run: gibt den Exit aus, legt stdout+stderr in $tmp/out ab
run() {
  bash tools/harness/ausgabe-kennungen-check.sh "$tmp/root" > "$tmp/out" 2>&1
  echo $?
}

check() {
  local want=$1 what=$2 got=$3
  cases=$((cases + 1))
  if [ "$got" -ne "$want" ]; then
    echo "FEHLER: $what -> Exit $got, erwartet $want: $(head -c 300 "$tmp/out")" >&2
    fail=1
  fi
}

# msg <Beschreibung> <ERE>: die letzte Ausgabe traegt den Meldungstext
msg() {
  cases=$((cases + 1))
  if ! grep -qE -e "$2" "$tmp/out"; then
    echo "FEHLER: $1 -> Meldung ohne /$2/: $(head -c 300 "$tmp/out")" >&2
    fail=1
  fi
}

# expect <Exit> <Beschreibung> <Datei relativ zur Wurzel> <Inhalt> [<Meldungs-ERE>]
expect() {
  fresh
  mkdir -p "$(dirname "$tmp/root/$3")"
  printf '%s\n' "$4" > "$tmp/root/$3"
  check "$1" "$2" "$(run)"
  if [ "$#" -ge 5 ]; then msg "$2" "$5"; fi
}

fresh
check 0 "saubere Wurzel" "$(run)"
msg "saubere Wurzel" '^ausgabe-kennungen-check: keine interne Kennung in Ausgabe-Literalen von 3 Go-Dateien und 2 Skripten$'

# Treffer im Go-Literal, je Kennungsart
expect 1 "LH im Println-Literal" internal/a/a.go 'func f() { fmt.Println("Betriebsstatus (LH-FA-ADM-002): ok") }' \
  '^internal/a/a\.go:1:.*LH-FA-ADM-002'
msg "Sammelzeile Exit 1" '^ausgabe-kennungen-check: interne Kennung in einem Ausgabe-Literal \(siehe oben\)$'
expect 1 "ADR im Errorf-Literal" internal/a/a.go 'var e = fmt.Errorf("%w: Text (ADR-0088 Festlegung 1)", x)'
expect 1 "SPEC im Literal" cmd/b/b.go 'var s = "siehe SPEC-016"'
expect 1 "ARC im Literal" tools/schema/c/c.go 'var s = "siehe ARC-004"'
expect 1 "Literal mit Schraegstrich davor" internal/a/a.go 'x := f("a/b") + "ADR-0043"'
expect 1 "Literal mit URL" internal/a/a.go 'errors.New("http://x/y ADR-0043")'
expect 1 "Raw-String in Backticks" internal/a/a.go 'var s = `Text ADR-0043`'
expect 1 "Printf mit Kennung nach Formatzeichen" internal/a/a.go 'fmt.Printf("  Zeile (LH-FA-RET-006): %.0f\n", v)'
expect 1 "Treffer in eingerueckter Zeile" internal/a/a.go $'func f() {\n\t\tfmt.Println("x (LH-FA-ADM-003): y")\n}'

# Treffer im echo/printf der Skripte
expect 1 "ADR im echo (tools/schema)" tools/schema/rollout.sh 'echo "schema-rollout: Vorlauf (ADR-0114) - x"' \
  '^tools/schema/rollout\.sh:1:.*ADR-0114'
expect 1 "LH im echo (examples)" examples/demo.sh '  echo "Lauf (LH-FA-CAP-009)" >&2'
expect 1 "SPEC im printf (examples)" examples/demo.sh 'printf "%s\n" "siehe SPEC-016"'
expect 1 "ARC im echo" tools/schema/rollout.sh 'echo "ARC-004"'

# Nicht-Treffer
expect 0 "Kommentarzeile mit Kennung" internal/a/a.go '// siehe ADR-0049 Teilfrage 1'
expect 0 "eingerueckte Kommentarzeile mit Literal und Kennung" internal/a/a.go $'func f() {\n\t// Fehlertext "x (LH-FA-ADM-002)" frueher\n}'
expect 0 "nachgestellter Kommentar ohne Anfuehrungszeichen" internal/a/a.go 'x := 1 // ADR-0049'
expect 0 "Struct-Tag in Backticks" internal/a/a.go 'Name string `json:"name" ADR-0049`'
expect 0 "Shell-Kommentar mit Kennung" examples/demo.sh '# Ausgabe (ADR-0043) im Kommentar'
expect 0 "Shell-Zeile ohne echo/printf" tools/schema/rollout.sh 'SCHEMA_SOURCE=x # ADR-0043'
expect 0 "Kennung ohne Folgezeichen aus der Klasse" internal/a/a.go 'x := "ADR-"'
expect 0 "Kennung in Test-Datei" internal/a/a_test.go 'var s = "Betriebsstatus (LH-FA-ADM-002)"'
expect 0 "Kennung in erzeugter Datei" internal/a/a.pb.go 'var s = "Betriebsstatus (LH-FA-ADM-002)"'
expect 0 "Kennung im Laeufer-Verzeichnis tools/harness" tools/harness/run-x.sh 'echo "Lauf (LH-FA-ADM-002)"'
expect 0 "Kennung in Go-Datei ausserhalb der Wurzeln" docs/x.go 'var s = "ADR-0043"'
expect 0 "Kennung in Skript ausserhalb der Wurzeln" docs/x.sh 'echo "ADR-0043"'

# Benannte Grenzen: der Falsch-Positiv (laut) und das Negativ (leise)
expect 1 "Falsch-Positiv: nachgestellter Kommentar mit zitiertem Literal" internal/a/a.go 'x := 1 // Fehler "ADR-0049" frueher'
expect 0 "Grenze: mehrzeiliges Raw-String-Literal wird nicht gelesen" internal/a/a.go $'var s = `\nZeile mit ADR-0043\n`'

# Lesefehler und fehlender Gegenstand
fresh
rm -rf "${tmp:?}/root/cmd"
check 2 "fehlende Wurzel cmd" "$(run)"
msg "fehlende Wurzel" '^ausgabe-kennungen-check: Lesefehler: Verzeichnis fehlt: cmd$'

fresh
rm "$tmp/root/tools/schema/rollout.sh" "$tmp/root/examples/demo.sh"
check 2 "kein Skript im Gegenstand" "$(run)"
msg "kein Gegenstand" '^ausgabe-kennungen-check: Lesefehler: kein Gegenstand \(3 Go-Dateien, 0 Skripte\); leer ist nicht bestanden$'

# Datei mit NUL-Byte: wird als Text gelesen, der Treffer bleibt sichtbar
fresh
printf 'a\0b\nvar s = "ADR-0043"\n' > "$tmp/root/internal/a/a.go"
check 1 "Kennung in Datei mit NUL-Byte" "$(run)"

# Lesefehler: als root sind Rechte wirkungslos, der Fall entfaellt dort
if [ "$(id -u)" -ne 0 ]; then
  fresh
  chmod 000 "$tmp/root/internal/a/a.go"
  check 2 "nicht lesbare Go-Datei" "$(run)"
  msg "nicht lesbare Go-Datei" '^ausgabe-kennungen-check: Lesefehler: internal/a/a\.go '

  fresh
  chmod 000 "$tmp/root/examples/demo.sh"
  check 2 "nicht lesbares Skript" "$(run)"
  msg "nicht lesbares Skript" '^ausgabe-kennungen-check: Lesefehler: examples/demo\.sh '

  fresh
  mkdir "$tmp/root/internal/gesperrt"
  chmod 000 "$tmp/root/internal/gesperrt"
  check 2 "nicht lesbares Unterverzeichnis" "$(run)"
  msg "nicht lesbares Unterverzeichnis" '^ausgabe-kennungen-check: Lesefehler beim Durchsuchen von internal '
else
  echo "run-ausgabe-kennungen-check-tests: Lesefehler-Faelle uebersprungen (root liest trotz chmod 000)"
fi

if [ "$fail" -eq 0 ]; then
  echo "run-ausgabe-kennungen-check-tests: alle $cases Fälle bestanden"
else
  echo "run-ausgabe-kennungen-check-tests: mindestens ein Fall gescheitert" >&2
  exit 1
fi
