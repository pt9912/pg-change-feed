#!/bin/bash
# run.sh — Treiber der C#-Kompatibilitätsmessung im Gast-Image (Schritte A1 bis
# A3 und A5; A4 ist die Mutationsprobe des Runners über das Verzeichnis /neu). Jede Zeile
# `KOMPAT csharp …` ist ein Messwert; eine Abweichung von der Erwartung des
# Schritts setzt den Exit auf 1, der Lauf geht bis zum Ende weiter.
set -u

# NEU_VERSION und REGISTRY_VERSION setzt das Gast-Image (Dockerfile, ENV).
: "${NEU_VERSION:?}" "${REGISTRY_VERSION:?}"

BIN=/kompat/bin
rot=0

kurz() { sha256sum "$1" | cut -c1-12; }

# lauf_ok <Schritt> <Verzeichnis>: das Gast-Programm endet mit Exit 0 und der Zeile "N Aufrufe ok".
lauf_ok() {
  local schritt=$1 dir=$2 out rc
  out=$(cd "$dir" && dotnet Gast.dll "$schritt" 2>&1)
  rc=$?
  printf '%s\n' "$out"
  if [ "$rc" -ne 0 ] || ! printf '%s\n' "$out" | grep -Eq "^KOMPAT csharp $schritt: [0-9]+ Aufrufe ok$"; then
    echo "KOMPAT csharp $schritt: ROT — erwartet Exit 0 und die Zeile Aufrufe ok, erhalten Exit $rc"
    rot=1
  fi
}

# lauf_scheitert <Schritt> <Verzeichnis> <Ausnahmetypen als Alternativen>: das
# Gast-Programm endet mit Exit != 0 und einem der Ausnahmetypen.
lauf_scheitert() {
  local schritt=$1 dir=$2 typen=$3 out rc
  out=$(cd "$dir" && dotnet Gast.dll "$schritt" 2>&1)
  rc=$?
  printf '%s\n' "$out"
  if [ "$rc" -eq 0 ] || ! printf '%s\n' "$out" | grep -Eq "Ausnahme ($typen):"; then
    echo "KOMPAT csharp $schritt: ROT — erwartet Exit != 0 mit $typen, erhalten Exit $rc"
    rot=1
  fi
}

# A1: Gast gegen das veröffentlichte 0.5.0 gebaut und gegen 0.5.0 gelaufen.
lauf_ok A1 "$BIN/alt05"

# A2: dieselben Gast-Binärdateien, nur die Bibliotheksdatei ist die neue
# (NEU_VERSION). Trägt /neu das Artefakt dieser Version (Modus dist), stammt die
# Datei aus ihm; fehlt die daraus gewonnene Datei, ist A2 rot und der Gast
# läuft nicht. Über den Runner ist dieser Zweig nicht erreichbar: die
# Bedingung ist dieselbe wie im Dockerfile (RUN-Schicht lib-dist), dort bricht
# ein fehlendes `cp` den Bau unter `set -eu` ab. Kopplung: wer die Bedingung
# oder den Zielpfad /kompat/lib-dist der RUN-Schicht lib-dist im Dockerfile
# ändert, ändert die Bedingung `[ -f "$neu_nupkg" ]` und `neu_dll` hier mit.
rm -rf /tmp/a2
cp -r "$BIN/alt05" /tmp/a2
neu_nupkg="/neu/PgChangeFeed.Client.$NEU_VERSION.nupkg"
neu_dll="$BIN/altreg/PgChangeFeed.Client.dll"
herkunft="Registry $REGISTRY_VERSION"
if [ -f "$neu_nupkg" ]; then
  neu_dll=/kompat/lib-dist/PgChangeFeed.Client.dll
  herkunft="Artefakt $(basename "$neu_nupkg")"
fi
if [ -f "$neu_dll" ]; then
  cp "$neu_dll" /tmp/a2/PgChangeFeed.Client.dll
  echo "KOMPAT csharp A2: Bibliothek $NEU_VERSION ($herkunft) $(kurz /tmp/a2/PgChangeFeed.Client.dll) ersetzt 0.5.0 $(kurz "$BIN/alt05/PgChangeFeed.Client.dll")"
  lauf_ok A2 /tmp/a2
else
  echo "KOMPAT csharp A2: ROT — Bibliotheksdatei $neu_dll fehlt ($herkunft)"
  rot=1
fi

# A3: Gegenrichtung. Grundlage: der gegen die veröffentlichte Bibliothek
# REGISTRY_VERSION gebaute Gast läuft gegen sie; danach mit der
# Bibliotheksdatei 0.5.0 — er muss mit einer Bindungsausnahme scheitern. Die Laufzeit weist die ältere Assembly-Version vor dem Aufruf ab
# (FileNotFoundException); eine fehlende Signatur bei gleicher Assembly-Version
# zeigt erst die Mutationsprobe A4 (MissingMethodException).
lauf_ok A3-Grundlage "$BIN/neureg"
rm -rf /tmp/a3
cp -r "$BIN/neureg" /tmp/a3
cp "$BIN/alt05/PgChangeFeed.Client.dll" /tmp/a3/PgChangeFeed.Client.dll
lauf_scheitert A3 /tmp/a3 'MissingMethodException|FileNotFoundException|FileLoadException'

# A5: Quellseite. Abhängigkeiten der zwei Package-Versionen, derselbe Gast-Quelltext
# gegen beide Versionen übersetzt, die null-Matrix je Version.
nuspec() { printf '%s/.nuget/packages/pgchangefeed.client/%s/pgchangefeed.client.nuspec' "$HOME" "$1"; }
deps() { grep -o '<dependency [^>]*>' "$(nuspec "$1")" | sort; }
if [ "$(deps 0.5.0)" = "$(deps "$REGISTRY_VERSION")" ]; then
  echo "KOMPAT csharp A5 abhaengigkeiten: 0.5.0 und $REGISTRY_VERSION gleich ($(deps 0.5.0 | wc -l) Zeilen)"
else
  echo "KOMPAT csharp A5 abhaengigkeiten: ROT — 0.5.0 und $REGISTRY_VERSION verschieden"
  diff <(deps 0.5.0) <(deps "$REGISTRY_VERSION")
  rot=1
fi

for v in 0.5.0 "$REGISTRY_VERSION"; do
  rm -rf "/tmp/a5-$v"
  cp -r /kompat/src/Gast "/tmp/a5-$v"
  if dotnet build "/tmp/a5-$v/Gast.csproj" -c Release -p:PgcfVersion="$v" -p:GastQuelle=alt -o "/tmp/a5-$v/out" >"/tmp/a5-$v.log" 2>&1; then
    echo "KOMPAT csharp A5 quelle $v: Gast-Quelltext (0.5.x-Formen) übersetzt"
  else
    echo "KOMPAT csharp A5 quelle $v: ROT — Übersetzung scheitert"
    grep -E 'error ' "/tmp/a5-$v.log" | sort -u | head -n 10
    rot=1
  fi
done

fall_zeilen=$(grep -n '// CASE ' /kompat/src/NullMatrix/NullMatrix.cs | awk -F: '{ name = $0; sub(/.*\/\/ CASE /, "", name); print $1 " " name }')
for v in 0.5.0 "$REGISTRY_VERSION"; do
  rm -rf "/tmp/nm-$v"
  cp -r /kompat/src/NullMatrix "/tmp/nm-$v"
  dotnet build "/tmp/nm-$v/NullMatrix.csproj" -c Release -p:PgcfVersion="$v" >"/tmp/nm-$v.log" 2>&1
  if ! grep -q 'Build succeeded\.' "/tmp/nm-$v.log" && ! grep -q 'error CS' "/tmp/nm-$v.log"; then
    echo "KOMPAT csharp A5 null-matrix $v: ROT — weder Erfolg noch Compilerfehler im Log"
    tail -n 5 "/tmp/nm-$v.log"
    rot=1
  fi
done
declare -A ergebnis
while read -r zeile fall; do
  [ -n "$fall" ] || continue
  for v in 0.5.0 "$REGISTRY_VERSION"; do
    code=$(grep -oE "NullMatrix\.cs\(${zeile},[0-9]+\): error CS[0-9]+" "/tmp/nm-$v.log" | head -n 1 | grep -oE 'CS[0-9]+' || true)
    ergebnis[$v]=${code:-ok}
  done
  erwartet_neu=ok
  if [ "$fall" = "http-basis-3-argumente-null-literal" ]; then
    erwartet_neu=CS0121
  fi
  urteil=wie-erwartet
  if [ "${ergebnis[0.5.0]}" != ok ] || [ "${ergebnis[$REGISTRY_VERSION]}" != "$erwartet_neu" ]; then
    urteil=ROT
    rot=1
  fi
  echo "KOMPAT csharp A5 null-matrix $fall: 0.5.0 ${ergebnis[0.5.0]}, $REGISTRY_VERSION ${ergebnis[$REGISTRY_VERSION]} (erwartet 0.5.0 ok, $REGISTRY_VERSION $erwartet_neu) $urteil"
done < <(printf '%s\n' "$fall_zeilen")

exit "$rot"
