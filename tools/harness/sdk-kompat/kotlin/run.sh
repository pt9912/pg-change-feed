#!/bin/bash
# run.sh — Treiber der Kotlin-Kompatibilitätsmessung im Gast-Image (Schritte A1
# bis A3 und A5; A4 ist die Mutationsprobe des Runners über das Verzeichnis
# /neu). Jede Zeile `KOMPAT kotlin …` ist ein Messwert; eine Abweichung von der
# Erwartung des Schritts setzt den Exit auf 1, der Lauf geht bis zum Ende weiter.
set -u

# NEU_VERSION und REGISTRY_VERSION setzt das Gast-Image (Dockerfile, ENV).
: "${NEU_VERSION:?}" "${REGISTRY_VERSION:?}"

BIN=/kompat/bin
LIB05=$BIN/alt05/lib/pgchangefeed-kotlin-0.5.0.jar
LIB06=$BIN/alt06/lib/pgchangefeed-kotlin-$REGISTRY_VERSION.jar
GAST=kompat.GastKt
rot=0

kurz() { sha256sum "$1" | cut -c1-12; }

# lauf_ok <Schritt> <Klassenpfad>: das Gast-Programm endet mit Exit 0 und der Zeile "N Aufrufe ok".
lauf_ok() {
  local schritt=$1 cp=$2 out rc
  out=$(java -cp "$cp" "$GAST" "$schritt" 2>&1)
  rc=$?
  printf '%s\n' "$out"
  if [ "$rc" -ne 0 ] || ! printf '%s\n' "$out" | grep -Eq "^KOMPAT kotlin $schritt: [0-9]+ Aufrufe ok$"; then
    echo "KOMPAT kotlin $schritt: ROT — erwartet Exit 0 und die Zeile Aufrufe ok, erhalten Exit $rc"
    rot=1
  fi
}

# lauf_scheitert <Schritt> <Klassenpfad> <Fehlertyp>: das Gast-Programm endet mit Exit != 0 und dem Fehlertyp.
lauf_scheitert() {
  local schritt=$1 cp=$2 typ=$3 out rc
  out=$(java -cp "$cp" "$GAST" "$schritt" 2>&1)
  rc=$?
  printf '%s\n' "$out"
  if [ "$rc" -eq 0 ] || ! printf '%s\n' "$out" | grep -q "Ausnahme $typ:"; then
    echo "KOMPAT kotlin $schritt: ROT — erwartet Exit != 0 mit $typ, erhalten Exit $rc"
    rot=1
  fi
}

# Der gebundene Deskriptor des Gast-Aufrufs und die Konstruktoren der Bibliothek je Version.
echo "KOMPAT kotlin javap Gast (alt05): $(javap -c -p -cp "$BIN/alt05/klassen" "$GAST" | grep -m1 'PgChangeFeedBadRequestException."<init>"' | sed 's/^ *//')"
echo "KOMPAT kotlin javap Gast (neu06): $(javap -c -p -cp "$BIN/neu06/klassen" "$GAST" | grep -m1 'PgChangeFeedBadRequestException."<init>"' | sed 's/^ *//')"
for v in 05 06; do
  jar=LIB$v
  echo "KOMPAT kotlin javap Bibliothek ${v}: $(javap -s -cp "${!jar}" io.github.pt9912.pgchangefeed.http.PgChangeFeedBadRequestException | grep -A1 'public io.github.pt9912.pgchangefeed.http.PgChangeFeedBadRequestException(' | grep -E 'descriptor' | sed 's/^ *descriptor: //' | tr '\n' ' ')"
done

# A1: Gast gegen das veröffentlichte 0.5.0 gebaut und gegen 0.5.0 gelaufen.
lauf_ok A1 "$BIN/alt05/klassen:$BIN/alt05/lib/*"

# A2: dieselben Gast-Klassen, Klassenpfad mit der veröffentlichten Laufzeit-Menge;
# liegt in /neu das Jar der neuen Version (Modus dist), tritt es an die Stelle
# des veröffentlichten Jars.
rm -rf /tmp/a2lib
cp -r "$BIN/alt06/lib" /tmp/a2lib
neu_jar=/tmp/a2lib/pgchangefeed-kotlin-$REGISTRY_VERSION.jar
herkunft="Registry $REGISTRY_VERSION"
if [ -f "/neu/pgchangefeed-kotlin-$NEU_VERSION.jar" ]; then
  rm "$neu_jar"
  neu_jar=/tmp/a2lib/pgchangefeed-kotlin-$NEU_VERSION.jar
  cp "/neu/pgchangefeed-kotlin-$NEU_VERSION.jar" "$neu_jar"
  herkunft="Artefakt $(basename "$neu_jar")"
fi
echo "KOMPAT kotlin A2: Bibliothek $NEU_VERSION ($herkunft) $(kurz "$neu_jar") ersetzt 0.5.0 $(kurz "$LIB05")"
lauf_ok A2 "$BIN/alt05/klassen:/tmp/a2lib/*"

# A3: Gegenrichtung. Grundlage: der gegen die veröffentlichte Bibliothek
# REGISTRY_VERSION übersetzte Gast läuft gegen sie; danach mit dem Klassenpfad
# 0.5.0 — er muss mit NoSuchMethodError scheitern.
lauf_ok A3-Grundlage "$BIN/neu06/klassen:$BIN/neu06/lib/*"
lauf_scheitert A3 "$BIN/neu06/klassen:$BIN/alt05/lib/*" NoSuchMethodError

# A5: Quellseite. Laufzeit-Mengen der zwei Versionen, derselbe Gast-Quelltext gegen beide übersetzt.
d05=$(ls "$BIN/alt05/lib" | grep -v '^pgchangefeed-kotlin-' || true)
d06=$(ls "$BIN/alt06/lib" | grep -v '^pgchangefeed-kotlin-' || true)
if [ "$d05" = "$d06" ]; then
  echo "KOMPAT kotlin A5 abhaengigkeiten: 0.5.0 und $REGISTRY_VERSION gleich ($(printf '%s\n' "$d05" | wc -l) Jars)"
else
  echo "KOMPAT kotlin A5 abhaengigkeiten: verschieden"
  diff <(printf '%s\n' "$d05") <(printf '%s\n' "$d06")
fi
for v in 0.5.0 "$REGISTRY_VERSION"; do
  if ./gradlew --no-daemon -q --rerun-tasks -PpgcfVersion="$v" -PgastQuelle=alt compileKotlin >"/tmp/a5-$v.log" 2>&1; then
    echo "KOMPAT kotlin A5 quelle $v: Gast-Quelltext (0.5.x-Formen, erschöpfendes when) übersetzt"
  else
    echo "KOMPAT kotlin A5 quelle $v: ROT — Übersetzung scheitert"
    grep -E '^e: ' "/tmp/a5-$v.log" | head -n 10
    rot=1
  fi
done

exit "$rot"
