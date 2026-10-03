#!/usr/bin/env bash
# run-meldungscodes-check-tests.sh — Tabellentest gegen meldungscodes-check.sh:
# je Zweig des Wächters ein Fall mit Meldungstext. Befunde (Exit 1): ein Code im
# Go-Quelltext, im Skript oder im Handbuch ohne Eintrag in der Tabelle, ein
# Tabellen-Code ohne Katalog-Zeile, eine Katalog-Zeile ohne Tabellen-Code, ein
# Code in falscher Form, ein Code nur in Prosa statt in der Katalog-Zeile.
# Nicht-Befunde (Exit 0): saubere Menge, zurückgezogener Code in Tabelle und
# Katalog, Test-Datei und erzeugte Datei, das Paket der Tabelle selbst, ein Code
# der Tabelle im Quelltext. Lesefehler (Exit 2): fehlende Wurzel, fehlende
# Tabelle, fehlender Katalog, Tabelle ohne Code, leerer Gegenstand, nicht
# lesbare Datei. Netzlos, kein Docker nötig.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

fail=0
cases=0
tmp=$(mktemp -d)
: "${tmp:?}"
trap 'chmod -R u+rwx "$tmp" 2>/dev/null; rm -rf "${tmp:?}"' EXIT

TABLE='package messagecode

const (
	DecodeUnreadable Code = "PCF-E4001"
	StorageFallback  Code = "PCF-E5000"
)
'
CATALOG='# Handbuch

| Code | Klasse | Bedeutung | Maßnahme |
|---|---|---|---|
| `PCF-E4001` | `schema` | Dekodierung | prüfen |
| `PCF-E5000` | `storage` | Rückfall | prüfen |
'
CLEAN_GO='package x

func f() { println("Der Betreiber startet den Feed-Container.") }
'
CLEAN_SH='echo "schema-rollout: Vorlauf"'

fresh() {
  chmod -R u+rwx "${tmp:?}/root" 2>/dev/null
  rm -rf "${tmp:?}/root"
  mkdir -p "$tmp/root/internal/a" "$tmp/root/internal/domain/messagecode" "$tmp/root/cmd/b" \
           "$tmp/root/tools/schema" "$tmp/root/examples" "$tmp/root/docs/user"
  printf '%s\n' "$TABLE" > "$tmp/root/internal/domain/messagecode/codes.go"
  printf '%s\n' "$CATALOG" > "$tmp/root/docs/user/benutzerhandbuch.md"
  printf '%s\n' "$CLEAN_GO" > "$tmp/root/internal/a/a.go"
  printf '%s\n' "$CLEAN_GO" > "$tmp/root/cmd/b/b.go"
  printf '%s\n' "$CLEAN_SH" > "$tmp/root/tools/schema/rollout.sh"
  printf '%s\n' "$CLEAN_SH" > "$tmp/root/examples/demo.sh"
}

# CHECK übersteuert den Prüfling (Mutationsläufe an einer Kopie).
CHECK=${CHECK:-tools/harness/meldungscodes-check.sh}

run() {
  bash "$CHECK" "$tmp/root" > "$tmp/out" 2>&1
  echo $?
}

check() {
  local want=$1 what=$2 got=$3
  cases=$((cases + 1))
  if [ "$got" -ne "$want" ]; then
    echo "FEHLER: $what -> Exit $got, erwartet $want: $(head -c 400 "$tmp/out")" >&2
    fail=1
  fi
}

msg() {
  cases=$((cases + 1))
  if ! grep -qE -e "$2" "$tmp/out"; then
    echo "FEHLER: $1 -> Meldung ohne /$2/: $(head -c 400 "$tmp/out")" >&2
    fail=1
  fi
}

# put <Datei relativ zur Wurzel> <Inhalt>
put() {
  mkdir -p "$(dirname "$tmp/root/$1")"
  printf '%s\n' "$2" > "$tmp/root/$1"
}

# --- sauber ----------------------------------------------------------------
fresh
check 0 "saubere Menge" "$(run)"
msg "saubere Menge" '^meldungscodes-check: 2 Codes in Tabelle und Katalog gleich, Quelltext \(2 Go-Dateien, 2 Skripte\)'

fresh
put internal/a/a.go 'package x

func f() error { return messagecode.New(messagecode.DecodeUnreadable, "Ursache") } // PCF-E4001'
check 0 "Code der Tabelle im Quelltext" "$(run)"

fresh
put tools/schema/rollout.sh 'echo "FEHLER [PCF-E5000]: Quelle fehlt" >&2'
check 0 "Code der Tabelle im Skript" "$(run)"

# --- Code im Quelltext ohne Tabelle (Befund) --------------------------------
fresh
put internal/a/a.go 'package x

var e = messagecode.New(messagecode.Other, "Ursache") // PCF-E4999'
check 1 "Go-Code ohne Tabelle" "$(run)"
msg "Go-Code ohne Tabelle" '^Code ohne Eintrag in der Tabelle internal/domain/messagecode/codes\.go: internal/a/a\.go:3:PCF-E4999$'
msg "Sammelzeile" '^meldungscodes-check: Quelltext, Tabelle und Katalog sind nicht gleich \(siehe oben\)$'

fresh
put cmd/b/b.go 'package x

var e = "PCF-E3001"'
check 1 "Go-Code in cmd ohne Tabelle" "$(run)"
msg "Go-Code in cmd ohne Tabelle" 'cmd/b/b\.go:3:PCF-E3001$'

fresh
put tools/schema/rollout.sh 'echo "FEHLER [PCF-E2007]: Quelle fehlt" >&2'
check 1 "Skript-Code ohne Tabelle" "$(run)"
msg "Skript-Code ohne Tabelle" 'tools/schema/rollout\.sh:1:PCF-E2007$'

fresh
put examples/demo.sh 'echo "[PCF-E2007]"'
check 1 "Skript-Code unter examples ohne Tabelle" "$(run)"
msg "Skript-Code unter examples ohne Tabelle" 'examples/demo\.sh:1:PCF-E2007$'

# --- falsche Form (Befund) --------------------------------------------------
fresh
put internal/a/a.go 'package x

var e = "PCF-X1234"'
check 1 "Buchstabe außerhalb E/W/I" "$(run)"
msg "falsche Form im Quelltext" '^Code in falscher Form: internal/a/a\.go:3:PCF-X1234 '

fresh
put internal/a/a.go 'package x

var e = "PCF-E401"'
check 1 "drei Ziffern" "$(run)"
msg "drei Ziffern" 'a\.go:3:PCF-E401 '

fresh
put internal/a/a.go 'package x

var e = "PCF-E40010"'
check 1 "fünf Ziffern" "$(run)"
msg "fünf Ziffern" 'a\.go:3:PCF-E40010 '

fresh
put internal/a/a.go 'package x

// Präfix PCF- ohne Code
var e = 1'
check 1 "Präfix ohne Code" "$(run)"
msg "Präfix ohne Code" 'a\.go:3:PCF- '

# --- Tabelle gegen Katalog (Befund) ----------------------------------------
fresh
put docs/user/benutzerhandbuch.md '| Code | Klasse | Bedeutung | Maßnahme |
|---|---|---|---|
| `PCF-E4001` | `schema` | Dekodierung | prüfen |'
check 1 "Tabellen-Code ohne Katalog-Zeile" "$(run)"
msg "Tabellen-Code ohne Katalog-Zeile" '^Tabellen-Code ohne Katalog-Zeile in docs/user/benutzerhandbuch\.md: PCF-E5000$'

fresh
put docs/user/benutzerhandbuch.md '| Code | Klasse | Bedeutung | Maßnahme |
|---|---|---|---|
| `PCF-E4001` | `schema` | Dekodierung | prüfen |
| `PCF-E5000` | `storage` | Rückfall | prüfen |
| `PCF-E6001` | `replication` | Strom | prüfen |'
check 1 "Katalog-Zeile ohne Tabellen-Code" "$(run)"
msg "Katalog-Zeile ohne Tabellen-Code (Katalog)" '^Katalog-Zeile ohne Eintrag in der Tabelle internal/domain/messagecode/codes\.go: PCF-E6001$'
msg "Katalog-Zeile ohne Tabellen-Code (Handbuch-Code)" '^Code ohne Eintrag in der Tabelle internal/domain/messagecode/codes\.go: docs/user/benutzerhandbuch\.md:5:PCF-E6001$'

fresh
put docs/user/benutzerhandbuch.md '# Handbuch

Prosa nennt `PCF-E4001` und `PCF-E5000`, aber eine Tabellenzeile fehlt.'
check 1 "Codes nur in Prosa" "$(run)"
msg "Codes nur in Prosa (E4001)" '^Tabellen-Code ohne Katalog-Zeile in docs/user/benutzerhandbuch\.md: PCF-E4001$'
msg "Codes nur in Prosa (E5000)" '^Tabellen-Code ohne Katalog-Zeile in docs/user/benutzerhandbuch\.md: PCF-E5000$'

fresh
put docs/user/benutzerhandbuch.md '| Code | Klasse | Bedeutung | Maßnahme |
|---|---|---|---|
| `PCF-E4001` | `schema` | siehe auch `PCF-E5000` | prüfen |'
check 1 "Code in einer späteren Zelle zählt nicht" "$(run)"
msg "Code in einer späteren Zelle" '^Tabellen-Code ohne Katalog-Zeile in docs/user/benutzerhandbuch\.md: PCF-E5000$'

fresh
put docs/user/benutzerhandbuch.md '| Code | Klasse | Bedeutung | Maßnahme |
|---|---|---|---|
| PCF-E4001 | schema | Dekodierung | prüfen |
| PCF-E5000 | storage | Rückfall | prüfen |'
check 0 "Katalog-Zeile ohne Backticks" "$(run)"

fresh
put docs/user/benutzerhandbuch.md "$CATALOG
Der Text nennt \`PCF-E9000\` in Prosa."
check 1 "Prosa-Code ohne Tabelle" "$(run)"
msg "Prosa-Code ohne Tabelle" 'benutzerhandbuch\.md:[0-9]+:PCF-E9000$'

# --- zurückgezogener Code bleibt in beiden Mengen ---------------------------
fresh
put internal/domain/messagecode/codes.go 'package messagecode

const (
	DecodeUnreadable Code = "PCF-E4001"
	StorageFallback  Code = "PCF-E5000"
	// Withdrawn: Status StatusWithdrawn, die Kennung bleibt vergeben.
	OldCode Code = "PCF-E4009"
)'
put docs/user/benutzerhandbuch.md '| Code | Klasse | Bedeutung | Maßnahme |
|---|---|---|---|
| `PCF-E4001` | `schema` | Dekodierung | prüfen |
| `PCF-E4009` | `schema` | zurückgezogen | entfallen |
| `PCF-E5000` | `storage` | Rückfall | prüfen |'
check 0 "zurückgezogener Code in Tabelle und Katalog" "$(run)"

fresh
put internal/domain/messagecode/codes.go 'package messagecode

const (
	DecodeUnreadable Code = "PCF-E4001"
	StorageFallback  Code = "PCF-E5000"
	OldCode          Code = "PCF-E4009"
)'
check 1 "zurückgezogener Code fehlt im Katalog" "$(run)"
msg "zurückgezogener Code fehlt im Katalog" '^Tabellen-Code ohne Katalog-Zeile in docs/user/benutzerhandbuch\.md: PCF-E4009$'

# --- Gegenstand ausgenommen -------------------------------------------------
fresh
put internal/a/a_test.go 'package x

var e = "PCF-E4999"'
check 0 "Test-Datei ausgenommen" "$(run)"

fresh
put internal/a/a.pb.go 'package x

var e = "PCF-E4999"'
check 0 "erzeugte Datei ausgenommen" "$(run)"

fresh
put internal/domain/messagecode/messagecode.go 'package messagecode

// Form: PCF- plus Schwere und vier Ziffern, Beispiel PCF-E4999.'
check 0 "Paket der Tabelle ausgenommen" "$(run)"

fresh
put tools/harness/x.sh 'echo "PCF-E4999"'
check 0 "Läufer unter tools/harness nicht gelesen" "$(run)"

# --- Datei mit NUL-Byte: grep liest sie als Text (-a) -----------------------
fresh
printf 'package x\n\0var e = "PCF-E4999"\n' > "$tmp/root/internal/a/a.go"
check 1 "Go-Datei mit NUL-Byte und Code ohne Tabelle" "$(run)"
msg "NUL-Byte-Datei" '^Code ohne Eintrag in der Tabelle internal/domain/messagecode/codes\.go: internal/a/a\.go:2:PCF-E4999$'

# --- Mengenbildung und -vergleich: ein Fehler des Werkzeugs ist Exit 2 ------
# Stubs vor dem echten Werkzeug im PATH; jeder lässt den Aufruf der Mengen-
# bildung scheitern und alle übrigen Aufrufe durch.
mkdir -p "$tmp/bin"
REAL_SED=$(command -v sed)
REAL_SORT=$(command -v sort)
REAL_COMM=$(command -v comm)
printf '#!/bin/sh\ncase "$1" in -n) exec %s "$@";; *) echo "sed: Stub-Fehler" >&2; exit 2;; esac\n' "$REAL_SED" > "$tmp/bin/sed"
printf '#!/bin/sh\ncase "$1" in -zu) exec %s "$@";; *) echo "sort: Stub-Fehler" >&2; exit 2;; esac\n' "$REAL_SORT" > "$tmp/bin/sort"
printf '#!/bin/sh\necho "comm: Stub-Fehler" >&2\nexit 2\n' > "$tmp/bin/comm"
chmod +x "$tmp/bin/sed" "$tmp/bin/sort" "$tmp/bin/comm"
stubbed() {
  local only=$1 dir="$tmp/bin-$1"
  rm -rf "${dir:?}"
  mkdir -p "$dir"
  cp "$tmp/bin/$only" "$dir/$only"
  PATH="$dir:$PATH" bash "$CHECK" "$tmp/root" > "$tmp/out" 2>&1
  echo $?
}

fresh
check 2 "Mengenbildung der Tabelle: sed scheitert" "$(stubbed sed)"
msg "sed der Tabellenmenge" '^meldungscodes-check: Lesefehler: Mengenbildung der Tabelle \(sed: Stub-Fehler\)$'

fresh
check 2 "Mengenbildung: sort scheitert" "$(stubbed sort)"
msg "sort der Tabellenmenge" '^meldungscodes-check: Lesefehler: Mengenbildung der Tabelle \(sort\)$'

fresh
check 2 "Mengenvergleich: comm scheitert" "$(stubbed comm)"
msg "comm Tabelle gegen Katalog" '^meldungscodes-check: Lesefehler: Mengenvergleich Tabelle gegen Katalog \(comm: Stub-Fehler\)$'

# --- Lesefehler (Exit 2) ----------------------------------------------------
check 2 "fehlende Wurzel" "$(bash "$CHECK" "$tmp/fehlt" > "$tmp/out" 2>&1; echo $?)"
msg "fehlende Wurzel" '^meldungscodes-check: Lesefehler: Wurzel fehlt: '

fresh
rm "$tmp/root/internal/domain/messagecode/codes.go"
check 2 "fehlende Tabelle" "$(run)"
msg "fehlende Tabelle" '^meldungscodes-check: Lesefehler: Code-Tabelle fehlt: internal/domain/messagecode/codes\.go$'

fresh
rm "$tmp/root/docs/user/benutzerhandbuch.md"
check 2 "fehlender Katalog" "$(run)"
msg "fehlender Katalog" '^meldungscodes-check: Lesefehler: Katalog fehlt: docs/user/benutzerhandbuch\.md$'

fresh
put internal/domain/messagecode/codes.go 'package messagecode
'
check 2 "Tabelle ohne Code" "$(run)"
msg "Tabelle ohne Code" 'Lesefehler: Code-Tabelle internal/domain/messagecode/codes\.go trägt keinen Code; leer ist nicht bestanden'

fresh
rm -r "${tmp:?}/root/cmd"
check 2 "fehlendes Verzeichnis cmd" "$(run)"
msg "fehlendes Verzeichnis" '^meldungscodes-check: Lesefehler: Verzeichnis fehlt: cmd$'

fresh
rm "$tmp/root/examples/demo.sh" "$tmp/root/tools/schema/rollout.sh"
check 2 "kein Skript" "$(run)"
msg "kein Skript" 'Lesefehler: kein Gegenstand \(2 Go-Dateien, 0 Skripte\); leer ist nicht bestanden'

fresh
rm "$tmp/root/internal/a/a.go" "$tmp/root/cmd/b/b.go"
check 2 "keine Go-Datei" "$(run)"
msg "keine Go-Datei" 'Lesefehler: kein Gegenstand \(0 Go-Dateien, 2 Skripte\)'

if [ "$(id -u)" -ne 0 ]; then
  fresh
  chmod 000 "$tmp/root/internal/a/a.go"
  check 2 "nicht lesbare Go-Datei" "$(run)"
  msg "nicht lesbare Go-Datei" '^meldungscodes-check: Lesefehler: internal/a/a\.go \('
  chmod 644 "$tmp/root/internal/a/a.go"

  fresh
  chmod 000 "$tmp/root/docs/user/benutzerhandbuch.md"
  check 2 "nicht lesbarer Katalog" "$(run)"
  msg "nicht lesbarer Katalog" '^meldungscodes-check: Lesefehler: docs/user/benutzerhandbuch\.md \('
  chmod 644 "$tmp/root/docs/user/benutzerhandbuch.md"

  fresh
  chmod 000 "$tmp/root/internal/domain/messagecode/codes.go"
  check 2 "nicht lesbare Tabelle" "$(run)"
  msg "nicht lesbare Tabelle" '^meldungscodes-check: Lesefehler: internal/domain/messagecode/codes\.go \('

  fresh
  chmod 000 "$tmp/root/tools/schema"
  check 2 "nicht lesbares Verzeichnis" "$(run)"
  msg "nicht lesbares Verzeichnis" '^meldungscodes-check: Lesefehler: '
  chmod 755 "$tmp/root/tools/schema"
else
  echo "run-meldungscodes-check-tests: unter root entfallen die vier Lesefehler-Fälle mit chmod 000"
fi

if [ "$fail" -ne 0 ]; then
  exit 1
fi
echo "run-meldungscodes-check-tests: $cases Prüfungen bestanden"
