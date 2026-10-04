#!/bin/bash
# run.sh — Treiber der Python-Kompatibilitätsmessung im Gast-Image (Schritte A1
# bis A3 und A5; A4 ist die Mutationsprobe des Runners über das Verzeichnis
# /neu). Jede Zeile `KOMPAT python …` ist ein Messwert; eine Abweichung von der
# Erwartung des Schritts setzt den Exit auf 1, der Lauf geht bis zum Ende weiter.
set -u

rot=0

# lauf_ok <Schritt> <Skript>: das Gast-Programm endet mit Exit 0 und der Zeile "N Aufrufe ok".
lauf_ok() {
  local schritt=$1 skript=$2 out rc
  out=$(python "$skript" "$schritt" 2>&1)
  rc=$?
  printf '%s\n' "$out"
  if [ "$rc" -ne 0 ] || ! printf '%s\n' "$out" | grep -Eq "^KOMPAT python $schritt: [0-9]+ Aufrufe ok$"; then
    echo "KOMPAT python $schritt: ROT — erwartet Exit 0 und die Zeile Aufrufe ok, erhalten Exit $rc"
    rot=1
  fi
}

# lauf_scheitert <Schritt> <Skript> <Fehlertyp>: das Gast-Programm endet mit Exit != 0 und dem Fehlertyp.
lauf_scheitert() {
  local schritt=$1 skript=$2 typ=$3 out rc
  out=$(python "$skript" "$schritt" 2>&1)
  rc=$?
  printf '%s\n' "$out"
  if [ "$rc" -eq 0 ] || ! printf '%s\n' "$out" | grep -q "Ausnahme $typ:"; then
    echo "KOMPAT python $schritt: ROT — erwartet Exit != 0 mit $typ, erhalten Exit $rc"
    rot=1
  fi
}

python -m venv /tmp/venv
# shellcheck disable=SC1091
. /tmp/venv/bin/activate

pip install -q --no-index --find-links /wheels/05 pgchangefeed==0.5.0
echo "KOMPAT python Umgebung: pgchangefeed $(pip show pgchangefeed | sed -n 's/^Version: //p')"

# A1: Gast gegen das veröffentlichte 0.5.0.
lauf_ok A1 /kompat/gast_alt.py

# A3: Gegenrichtung gegen die Bibliothek 0.5.0 — der Konstruktor nimmt `message_code` nicht an.
lauf_scheitert A3 /kompat/gast_neu.py TypeError

# A5 (Grundlinie): die Signaturen unter 0.5.0.
python /kompat/quelle.py || rot=1

# A2: dieselbe Umgebung, die Bibliothek wird gegen 0.6.0 ausgetauscht, derselbe Gast.
pip install -q --no-index --find-links /wheels/06 pgchangefeed==0.6.0
herkunft=Registry
if ls /neu/pgchangefeed-*.whl >/dev/null 2>&1; then
  herkunft="Artefakt $(basename "$(ls /neu/pgchangefeed-*.whl | head -n 1)")"
fi
echo "KOMPAT python A2: Bibliothek pgchangefeed $(pip show pgchangefeed | sed -n 's/^Version: //p') ($herkunft) ersetzt 0.5.0"
lauf_ok A2 /kompat/gast_alt.py
lauf_ok A3-Grundlage /kompat/gast_neu.py

# A5: Quellseite unter 0.6.0 — die 0.5.x-Form bindet weiter, ein drittes Positional nicht.
python /kompat/quelle.py || rot=1

exit "$rot"
