#!/usr/bin/env bash
# run-pin-stale-all-tests.sh — Tabellentest gegen pin-stale-all.sh: je Zweig
# der Aufzählung, des Vergleichs, der Tag-Ableitung und der Ausgänge ein Fall
# mit Meldungstext. Netzlos: ein Stub-`docker` im Wegwerf-Repo beantwortet die
# Registry-Abfrage aus einer Tabelle. Die Digests der Fixtures entstehen zur
# Laufzeit (hex <Zeichen>); der Test trägt keine vollständige Referenz als
# Literal, weil der Sensor sie im Baum fände. Der Prüfling ist per
# PROG=<Datei> übersteuerbar (Mutationsläufe an Kopien; die gemeinsame
# lib-pin-compare.sh liegt neben dem Prüfling).
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"

PROG=$(realpath "${PROG:-tools/harness/pin-stale-all.sh}")
fail=0
cases=0
tmp=$(mktemp -d)
: "${tmp:?}"
trap 'rm -rf "${tmp:?}"' EXIT

mkdir -p "$tmp/bin" "$tmp/table"
export STUB_TABLE="$tmp/table" STUB_LOG="$tmp/calls.log"
# Stub: `docker buildx imagetools inspect <ziel> --format …` liest den Digest
# aus Tabellendatei <ziel> (`/` und `:` als `_`); ohne Datei endet er mit 1.
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'if [ "$1 $2 $3" != "buildx imagetools inspect" ]; then exit 99; fi' \
  'echo "$4" >> "$STUB_LOG"' \
  'f="$STUB_TABLE/$(printf %s "$4" | tr "/:" "__")"' \
  '[ -f "$f" ] || exit 1' \
  '[ -f "$f.sleep" ] && exec sleep "$(cat "$f.sleep")"' \
  'cat "$f"' > "$tmp/bin/docker"
chmod +x "$tmp/bin/docker"

# hex <zeichen>: 64 Hex-Zeichen
hex() { local i s=""; for i in $(seq 64); do s="$s$1"; done; printf '%s' "$s"; }
A=$(hex a)
B=$(hex b)
C=$(hex c)

# fresh: leeres Wegwerf-Repo ohne Stub-Tabelle
fresh() {
  rm -rf "${tmp:?}/repo" "${tmp:?}/table"/* "$STUB_LOG"
  mkdir -p "$tmp/repo"
  git -C "$tmp/repo" init -q
  : > "$STUB_LOG"
}
# put <pfad> <inhalt>: Datei im Wegwerf-Repo
put() {
  mkdir -p "$tmp/repo/$(dirname "$1")"
  printf '%s\n' "$2" > "$tmp/repo/$1"
}
# serve <ziel> <digest>: die Registry kennt das Ziel
serve() { printf '%s\n' "$2" > "$STUB_TABLE/$(printf %s "$1" | tr '/:' '__')"; }
# run: Exit nach stdout, Ausgabe in $tmp/out
run() {
  git -C "$tmp/repo" add -A
  (cd "$tmp/repo" && PATH="$tmp/bin:$PATH" bash "$PROG" > "$tmp/out" 2>&1)
  echo $?
}
check() {
  cases=$((cases + 1))
  if [ "$3" -ne "$1" ]; then
    echo "FEHLER: $2 -> Exit $3, erwartet $1" >&2
    fail=1
  fi
}
# msg <Beschreibung> <ERE>: die Ausgabe trägt den Text
msg() {
  cases=$((cases + 1))
  if ! grep -qE -e "$2" "$tmp/out"; then
    echo "FEHLER: $1 -> Ausgabe ohne /$2/: $(head -c 400 "$tmp/out")" >&2
    fail=1
  fi
}
# nomsg <Beschreibung> <ERE>: die Ausgabe trägt den Text nicht
nomsg() {
  cases=$((cases + 1))
  if grep -qE -e "$2" "$tmp/out"; then
    echo "FEHLER: $1 -> Ausgabe trägt /$2/: $(head -c 400 "$tmp/out")" >&2
    fail=1
  fi
}
# calls <Beschreibung> <Ziel> <Anzahl>: Aufrufe des Stubs mit genau diesem Ziel
calls() {
  local n
  n=$(grep -cxF -e "$2" "$STUB_LOG")
  cases=$((cases + 1))
  if [ "$n" -ne "$3" ]; then
    echo "FEHLER: $1 -> $n Registry-Aufrufe für $2, erwartet $3" >&2
    fail=1
  fi
}

# --- Aufzählung (Festlegung 1): yml, Shell-Default, compose.yaml ---
fresh
put .github/workflows/e2e.yml "        pg: img-yml:1@sha256:$A"
put tools/run.sh "IMG=\${IMG:-img-sh:2@sha256:$A}"
put compose.yaml "    image: \${IMG:-img-compose:3@sha256:$A}"
serve img-yml:1 "sha256:$A"
serve img-sh:2 "sha256:$A"
serve img-compose:3 "sha256:$A"
check 0 "Treffer in yml, Shell-Default und compose.yaml, alle gleich" "$(run)"
msg "Fundort yml" '^OK +\.github/workflows/e2e\.yml:1 \(img-yml:1\) == sha256:a{64}$'
msg "Fundort Shell-Default" '^OK +tools/run\.sh:1 \(img-sh:2\) == sha256:a{64}$'
msg "Fundort compose.yaml" '^OK +compose\.yaml:1 \(img-compose:3\) == sha256:a{64}$'
msg "Zusammenfassung" '^pin-stale-all: 3 Referenzen — 3 OK, 0 DRIFT, 0 UNBESTIMMT$'

# --- Vergleich: Drift, Einzelplattform-Digest als Drift (Festlegung 2) ---
fresh
put Makefile "IMG ?= img-d:1@sha256:$A"
serve img-d:1 "sha256:$B"
check 1 "Digest weicht ab" "$(run)"
msg "DRIFT mit beiden Digests" "^DRIFT +Makefile:1 \(img-d:1\): gepinnt sha256:a{64}, aktuell sha256:b{64}$"
msg "Zusammenfassung Drift" '^pin-stale-all: 1 Referenzen — 0 OK, 1 DRIFT, 0 UNBESTIMMT$'

fresh
put e2e.yml "pg: img-plat:18@sha256:$C"
serve img-plat:18 "sha256:$B"
check 1 "Einzelplattform-Digest gegen Index-Digest" "$(run)"
msg "Einzelplattform-Digest als DRIFT" "^DRIFT +e2e\.yml:1 \(img-plat:18\): gepinnt sha256:c{64}, aktuell sha256:b{64}$"

# --- Tag-Ableitung (Festlegung 3) ---
fresh
put a.mk "T ?= reg.example/ns/img-tag:7@sha256:$A"
put b.mk "U ?= ns/img-notag@sha256:$B"
serve reg.example/ns/img-tag:7 "sha256:$A"
serve ns/img-notag:latest "sha256:$B"
check 0 "Pin mit Tag (Registry-Host) und Pin ohne Tag" "$(run)"
msg "Pin mit Tag gegen den eigenen Tag" '^OK +a\.mk:1 \(reg\.example/ns/img-tag:7\) =='
msg "Pin ohne Tag gegen :latest" '^OK +b\.mk:1 \(ns/img-notag:latest\) =='
calls "Pin ohne Tag" ns/img-notag:latest 1
calls "Pin ohne Tag, kein Aufruf ohne Tag" ns/img-notag 0

# --- Ausschlüsse: docs/ und .harness/ ---
fresh
put ok.mk "I ?= img-ok:1@sha256:$A"
put docs/plan/adr/x.md "Beispiel img-docs:1@sha256:$B im Fließtext"
put .harness/baseline/v1/regelwerk/y.md "img-harness:1@sha256:$B"
serve img-ok:1 "sha256:$A"
check 0 "docs/ und .harness/ ausgenommen" "$(run)"
nomsg "docs/-Referenz nicht gemeldet" 'img-docs'
nomsg ".harness/-Referenz nicht gemeldet" 'img-harness'
msg "nur die Referenz außerhalb" '^pin-stale-all: 1 Referenzen — 1 OK, 0 DRIFT, 0 UNBESTIMMT$'

# --- Deduplizierung: dieselbe Referenz an zwei Fundorten ---
fresh
put a.mk "I ?= img-dup:1@sha256:$A"
put b/c.sh "I=\${I:-img-dup:1@sha256:$A}"
serve img-dup:1 "sha256:$A"
check 0 "dieselbe Referenz in zwei Dateien" "$(run)"
msg "ein Eintrag mit Fundort und Zähler" '^OK +a\.mk:1 \+1 weitere \(img-dup:1\) =='
calls "ein Registry-Aufruf je verschiedene Referenz" img-dup:1 1

# --- fail-open: Registry nicht erreichbar (Festlegung 4) ---
fresh
put a.mk "I ?= img-up:1@sha256:$A"
put b.mk "J ?= img-down:1@sha256:$A"
put c.mk "K ?= img-up2:1@sha256:$A"
serve img-up:1 "sha256:$A"
serve img-up2:1 "sha256:$A"
check 2 "eine Registry-Abfrage scheitert, die übrigen laufen weiter" "$(run)"
msg "UNBESTIMMT je Referenz" '^UNBESTIMMT +b\.mk:1 — img-down:1: Registry nicht erreichbar$'
msg "Referenz davor bleibt OK" '^OK +a\.mk:1 \(img-up:1\)'
msg "Referenz danach bleibt OK" '^OK +c\.mk:1 \(img-up2:1\)'
msg "Zusammenfassung UNBESTIMMT" '^pin-stale-all: 3 Referenzen — 2 OK, 0 DRIFT, 1 UNBESTIMMT$'

fresh
put a.mk "I ?= img-drift:1@sha256:$A"
put b.mk "J ?= img-down:1@sha256:$A"
serve img-drift:1 "sha256:$B"
check 1 "DRIFT neben UNBESTIMMT" "$(run)"
msg "DRIFT bleibt sichtbar" '^DRIFT +a\.mk:1 \(img-drift:1\)'
msg "UNBESTIMMT bleibt sichtbar" '^UNBESTIMMT +b\.mk:1 — img-down:1:'

# --- Zeitlimit je Registry-Aufruf (Festlegung 4): Ablauf ist UNBESTIMMT ---
# Der Stub der langsamen Referenz schläft 6 s (`exec`, damit `timeout` den
# schlafenden Prozess selbst trifft); das Limit ist 1 s.
fresh
put a.mk "I ?= img-slow:1@sha256:$A"
put b.mk "J ?= img-fast:1@sha256:$A"
serve img-slow:1 "sha256:$A"
printf '6\n' > "$STUB_TABLE/img-slow_1.sleep"
serve img-fast:1 "sha256:$A"
export PIN_COMPARE_TIMEOUT=1
t0=$SECONDS
rc=$(run)
elapsed=$((SECONDS - t0))
unset PIN_COMPARE_TIMEOUT
check 2 "Registry-Aufruf überschreitet das Zeitlimit" "$rc"
msg "Zeitüberschreitung ist UNBESTIMMT" '^UNBESTIMMT +a\.mk:1 — img-slow:1: Registry nicht erreichbar$'
nomsg "Zeitüberschreitung ist kein OK" '^OK +a\.mk:1'
nomsg "Zeitüberschreitung ist kein DRIFT" '^DRIFT'
msg "schnelle Referenz bleibt OK" '^OK +b\.mk:1 \(img-fast:1\)'
msg "Zusammenfassung Zeitlimit" '^pin-stale-all: 2 Referenzen — 1 OK, 0 DRIFT, 1 UNBESTIMMT$'
cases=$((cases + 1))
if [ "$elapsed" -ge 5 ]; then
  echo "FEHLER: Lauf dauerte ${elapsed} s, das Limit von 1 s begrenzt den Aufruf nicht" >&2
  fail=1
fi

# Gegenfall: gleiches Limit, Antwort ohne Verzögerung
fresh
put a.mk "I ?= img-fast:1@sha256:$A"
serve img-fast:1 "sha256:$A"
export PIN_COMPARE_TIMEOUT=1
rc=$(run)
unset PIN_COMPARE_TIMEOUT
check 0 "Antwort innerhalb des Zeitlimits" "$rc"
msg "Antwort innerhalb des Limits ist OK" '^OK +a\.mk:1 \(img-fast:1\)'

# --- Registry-Adresse mit Port (Grenze 3 des Vertrags) ---
# Der Treffer beginnt am Port selbst; das Abfrageziel ist der Rest ab dort.
fresh
put a.mk "I ?= localhost:5000/ns/img-port:1@sha256:$A"
serve 5000/ns/img-port:1 "sha256:$A"
check 0 "Referenz mit Registry-Port" "$(run)"
msg "Treffer beginnt am Port" '^OK +a\.mk:1 \(5000/ns/img-port:1\) =='
calls "Abfrageziel ohne Host" 5000/ns/img-port:1 1
calls "Host mit Port wird nicht abgefragt" localhost:5000/ns/img-port:1 0

# --- leerer Gegenstand (Festlegung 4): keine Referenz ist nicht bestanden ---
fresh
put a.mk "I ?= img-nodigest:1"
put b.mk "J ?= img-short:1@sha256:${A:0:63}"
check 2 "keine Referenz im Baum" "$(run)"
msg "Meldung leerer Gegenstand" '^UNBESTIMMT +keine Digest-Pin-Referenz im Baum gefunden — leerer Gegenstand$'

fresh
put docs/x.md "img-docs:1@sha256:$A"
check 2 "einzige Referenz liegt unter docs/" "$(run)"
msg "Meldung leerer Gegenstand unter docs/" 'leerer Gegenstand$'

if [ "$fail" -eq 0 ]; then
  echo "run-pin-stale-all-tests: alle $cases Prüfungen bestanden"
else
  echo "run-pin-stale-all-tests: mindestens eine Prüfung gescheitert" >&2
  exit 1
fi
