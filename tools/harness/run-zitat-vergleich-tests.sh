#!/usr/bin/env bash
# run-zitat-vergleich-tests.sh — Tabellentest gegen zitat-vergleich.sh
# (harness/targets/zitat-vergleich.md): je Fall ein Aufruf gegen ein
# Wegwerf-Repo im Temp-Verzeichnis, mit Soll-Exit und Muster der Ausgabe.
# Die Fall-Tabelle läuft dreimal: ohne Shell-Option, unter nullglob und unter
# failglob (über BASHOPTS des Prüflings). Fallgruppen: HTML-id je Stellung ·
# fail-closed · set -euo pipefail · roh · Normalisierung und Tag-Paar ·
# Regressionsfälle (Anker-Wechsel, Lokator, git mv) · Heading mit Inline-Code ·
# Fence nach CommonMark · gestapelte id · id in anderer Form · Locale und
# Fähigkeitsprobe von awk · Argumente.
# Netzlos, kein Docker (bash, git, awk, sed, cmp, env). Der Prüfling ist per
# PROG übersteuerbar (Mutationsläufe gegen eine Kopie).
set -uo pipefail
repo=$(git rev-parse --show-toplevel)
prog=${PROG:-$repo/tools/harness/zitat-vergleich.sh}
case "$prog" in /*) ;; *) prog=$PWD/$prog ;; esac
[ -f "$prog" ] || { echo "run-zitat-vergleich-tests: Prüfling $prog fehlt" >&2; exit 2; }
real_awk=$(command -v awk)

tmp=$(mktemp -d "${TMPDIR:-/tmp}/zitat-vergleich-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT
cd "$tmp" || exit 1

git init -q .
git config user.name test
git config user.email test@example.invalid
git config commit.gpgsign false

# commit <Nachricht> — committet den Arbeitsbaum und druckt die kurze Kennung;
# ein Fehler beim Aufbau beendet den Test mit Exit 2.
commit() {
  git add -A && git commit -q -m "$1" >/dev/null || { echo "run-zitat-vergleich-tests: Aufbau-Commit '$1' gescheitert" >&2; exit 2; }
  git rev-parse --short HEAD
}

# --- Stand c0: alle Fixtures in Ausgangsform ---------------------------------
mkdir -p .harness/baseline/v6.14.0
printf 'Baseline v6.14.0\n' >.harness/baseline/v6.14.0/README.md
g_base='# Modul 13\n\n<a id="gh"></a>\n\n## Guard-Härtung: Wächter reifen\n\nKörper eins.\n\n## Danach\n\nRest.\n'
printf "$g_base" >g1.md
printf "$g_base" >g2.md
printf "$g_base" >g3.md
printf '# Tabelle\n\n| MR | Titel |\n|---|---|\n| MR-003 <a id="mr-003"></a> | alpha |\n| MR-004 <a id="mr-004"></a> | beta |\n' >t.md
printf '# P\n\n<a id="absatz"></a>\nErster Absatz alpha.\n\nZweiter Absatz.\n\n## Weiter\n\nx\n' >p.md
printf '# H\n\n## 2x2 Matrix <a id="m22"></a>\n\nInhalt alpha.\n' >h.md
printf '# I\n\nSiehe `<a id="code"></a>` hier.\n' >i.md
printf 'x\n' >r1.md
printf '## S\n\ntext' >s.md
printf '# A\n\n## Teilfrage 2\n\nzwei\n\n## Teilfrage 3\n\ndrei\n' >ad.md
printf 'eins\nzwei\ndrei\nvier\nfuenf\nsechs\n' >lz.md
printf '# Mv\n\ninhalt\n' >mv.md
printf '# N1\n\n## Der `make gates`-Lauf\n\nalpha\n' >n1.md
printf '## Ungerade\n\nvor\n\n````\n```\n## im fence\n````\n\nnachher alpha\n\n## Ende\n\nschluss\n' >fe.md
printf '## Gerade\n\nvor\n\n````\n```\ncode\n```\n## im fence\n````\n\nnachher alpha\n' >fe2.md
printf '## Tilde\n\n```\n~~~\n## nicht\n```\n\nnachher alpha\n' >ft.md
printf '## Eins\n\n```a`b\n## Zwei\n\nzwei alpha\n' >fb.md
printf '# S\n\n<a id="oben"></a>\n<a id="unten"></a>\n\n## Abschnitt\n\nKörper alpha.\n' >st.md
printf '<a id="s1"></a>\n<a id="s2"></a>\n\nAbsatz alpha.\n\n## H\n\nh\n' >sp.md
printf '# F\n\n<a id="attr" class="k"></a>\n\n## A\n\nalpha\n' >fa.md
printf '# F\n\n<a id="selbst"/>\n\n## B\n\nalpha\n' >fs.md
printf '# F\n\n## Doppel\n\nalpha\n\n<a id="doppel" class="x"></a>\n' >fh.md
printf '## Fi\n\nText mit `<a id="fi" class="k">` hier.\n' >fi.md
printf '<a id="fz"></a>\n\n```\ncode\n```\ntext\n\n## Next\n\nn\n' >fz.md
printf '## Pins\n\nPfad .harness/baseline/v6.14.0/regelwerk/x.md und https://github.com/pt9912/ai-harness-course/blob/v6.14.0/x.md\n' >n.md
printf '## Pins\n\nPfad .harness/baseline/v6.14.0/x.md und https://example.org/tool/v0.79.0/bin\n' >fp.md
printf '## Pins\n\nWerkzeug https://example.org/tool/v0.79.0/bin\n' >fo.md
printf '## Pins\n\nPfad .harness/baseline/v6.14.0/x und https://example.org/tool-v6.14.0/x\n' >seg.md
printf '## Pins\n\nPfad .harness/baseline/v6.14.0/x und https://example.org/v6-14-0/x\n' >dot.md
printf '## Eins\n\n   ```\n## im fence\n   ```\n\neins alpha\n\n## Zwei\n\nz\n' >fi3.md
printf '## Eins\n\n    ```\n## Zwei\n\nzwei alpha\n' >fi4.md
printf '## Code <code>x</code>\n\nalpha\n' >code.md
printf '## Dup\n\neins\n\n## Dup\n\nzwei alpha\n\n## Dup\n\ndrei\n' >dup.md
printf '###### Sechs\n\nalpha\n\n####### sieben\n\nnach alpha\n' >lvl.md
printf '1. Liste\n\n    ```\n    <a id="ein"></a>\n    ```\n\n<a id="ein"></a>\n\n## Ziel\n\nziel alpha\n' >indent.md
printf '<!-- Kommentar\n<a id="kom"></a>\n-->\n<!-- <a id="k2"></a> -->\n\n<a id="kom"></a>\n<a id="k2"></a>\n\n## Ziel\n\nziel alpha\n' >kom.md
printf 'Text mit `<!--` in Inline-Code.\n\n<a id="k3"></a>\n\n## Z3\n\nz3 alpha\n' >ik.md
printf '## Eins ##\n\nalpha\n\n## Zwei\n\nz\n' >close.md
printf 'Ein ` Backtick <!-- <a id="m4"></a> -->\n\n<a id="m4"></a>\n\n## Ziel\n\nziel alpha\n' >n4.md
printf '``a`b`` <!-- <a id="m5"></a> --> `\n\n<a id="m5"></a>\n\n## Ziel\n\nziel alpha\n' >n5.md
printf 'Text `beginnt\nendet `<!-- <a id="m6"></a> -->` x\n\n<a id="m6"></a>\n\n## Ziel\n\nziel alpha\n' >n6.md
printf 'Ein ` Backtick <!--\n\n<a id="m7"></a>\n\n-->\n\n<a id="m7"></a>\n\n## Ziel\n\nziel alpha\n' >un.md
printf '  \t<a id="m3"></a>\n\n<a id="m3"></a>\n\n## Ziel\n\nziel alpha\n' >mi.md
printf '<!-- k --> <a id="m8"></a>\n\n## Ziel\n\nziel alpha\n' >hc.md
printf '| MR | Titel | Text |\n|---|---|---|\n| MR-009 <a id="mr-009"></a> | `x` | y alpha |\n' >tb.md
printf 'z1\nz2\nz3\nz4\nz5\nz6\nz7\nz8\nz9\nz10\nz11\nz12\n' >l12.md
printf 'Absatz `siehe <a id="q1"></a>` hier.\n\n<a id="q1"></a>\n\n## Ziel\n\nziel alpha\n' >e1.md
printf -- '- `<!-- k --> <a id="q2"></a>` gilt als ohne Inhalt.\n\n<a id="q2"></a>\n\n## Ziel\n\nziel alpha\n' >e2.md
printf '| A | B |\n|---|---|\n| x | `Form <a id="q6"></a>` |\n\n<a id="q6"></a>\n\n## Ziel\n\nziel alpha\n' >e6.md
printf '## Kopf\n\nNur `in <a id="q7"></a> Code`.\n' >e7.md
printf 'Erst `code` dann <a id="q8"></a> Text alpha\n\n## Z\n\nz\n' >e8.md
c0=$(commit c0) || exit 2

# --- je Commit eine Änderung -------------------------------------------------
printf '# Modul 13\n\n<a id="gh"></a>\n\n## Guard-Härtung: Wächter reifen neu\n\nKörper eins.\n\n## Danach\n\nRest.\n' >g1.md
c_head=$(commit "nur Heading-Zeile") || exit 2
printf '# Modul 13\n\n<a id="gh"></a>\n\n## Guard-Härtung: Wächter reifen\n\nKörper zwei.\n\n## Danach\n\nRest.\n' >g2.md
c_body=$(commit "Wort im Körper") || exit 2
printf '# Modul 13\n\n<a id="gh"></a>\n\n## Guard-Härtung: Wächter reifen\n\nKörper eins.\n\n\n## Danach\n\nRest.\n' >g3.md
c_leer=$(commit "Leerzeile im Abschnitt") || exit 2
printf '# Tabelle\n\n| MR | Titel |\n|---|---|\n| MR-003 <a id="mr-003"></a> | alpha |\n| MR-004 <a id="mr-004"></a> | gamma |\n' >t.md
c_tab=$(commit "Wort in Zeile MR-004") || exit 2
printf '# P\n\n<a id="absatz"></a>\nErster Absatz beta.\n\nZweiter Absatz.\n\n## Weiter\n\nx\n' >p.md
printf '# H\n\n## 2x2 Matrix <a id="m22"></a>\n\nInhalt beta.\n' >h.md
printf 'x\n\n\n\n' >r1.md
printf '## S\n\ntext\n' >s.md
printf 'neu1\nneu2\nneu3\neins\nzwei\ndrei\nvier\nfuenf\nsechs\n' >lz.md
printf '# N1\n\n## Der `make gates`-Lauf\n\nbeta\n' >n1.md
printf '## Ungerade\n\nvor\n\n````\n```\n## im fence\n````\n\nnachher beta\n\n## Ende\n\nschluss\n' >fe.md
printf '## Gerade\n\nvor\n\n````\n```\ncode\n```\n## im fence\n````\n\nnachher beta\n' >fe2.md
printf '## Tilde\n\n```\n~~~\n## nicht\n```\n\nnachher beta\n' >ft.md
printf '## Eins\n\n```a`b\n## Zwei\n\nzwei beta\n' >fb.md
printf '# S\n\n<a id="oben"></a>\n<a id="unten"></a>\n\n## Abschnitt\n\nKörper beta.\n' >st.md
printf '<a id="s1"></a>\n<a id="s2"></a>\n\nAbsatz beta.\n\n## H\n\nh\n' >sp.md
for f in fi3.md fi4.md code.md dup.md lvl.md indent.md kom.md ik.md close.md n4.md n5.md n6.md un.md mi.md hc.md tb.md e1.md e2.md e6.md e8.md; do
  sed 's/alpha/beta/' "$f" >"$f.neu" && mv "$f.neu" "$f"
done
c_wort=$(commit "Wortänderungen") || exit 2
git mv mv.md mv2.md || exit 2
c_mv=$(commit "git mv") || exit 2
# Bump: neuer Baseline-Baum neben dem alten, Tag in Pfad und Kurs-URL bewegt;
# in fp.md zusätzlich ein fremder Pin, in fo.md nur ein fremder Pin.
mkdir -p .harness/baseline/v6.14.1
printf 'Baseline v6.14.1\n' >.harness/baseline/v6.14.1/README.md
printf '## Pins\n\nPfad .harness/baseline/v6.14.1/regelwerk/x.md und https://github.com/pt9912/ai-harness-course/blob/v6.14.1/x.md\n' >n.md
printf '## Pins\n\nPfad .harness/baseline/v6.14.1/x.md und https://example.org/tool/v0.80.0/bin\n' >fp.md
printf '## Pins\n\nWerkzeug https://example.org/tool/v0.80.0/bin\n' >fo.md
printf '## Pins\n\nPfad .harness/baseline/v6.14.1/x und https://example.org/tool-v6.14.1/x\n' >seg.md
printf '## Pins\n\nPfad .harness/baseline/v6.14.1/x und https://example.org/v6-14-0/x\n' >dot.md
c_bump=$(commit "Bump") || exit 2

fail=0
count=0
out=""
rc=0
opts=""

# run <Argumente des Prüflings> — Ausgabe (stdout und stderr) und Exit; die
# Shell-Option der Runde steht in BASHOPTS des Prüflings.
run() {
  if [ -n "$opts" ]; then
    out=$(env BASHOPTS="$opts" bash "$prog" "$@" 2>&1)
  else
    out=$(bash "$prog" "$@" 2>&1)
  fi
  rc=$?
}

expect() { # <Fallname> <erwarteter Exit> <Muster in der Ausgabe oder ->
  local name=$1 want_rc=$2 pattern=$3
  count=$((count + 1))
  if [ "$rc" -ne "$want_rc" ]; then
    echo "FEHLER [${opts:-ohne Option}]: $name — Exit $rc, erwartet $want_rc; Ausgabe:" >&2
    printf '%s\n' "$out" >&2
    fail=1
  elif [ "$pattern" != "-" ] && ! printf '%s\n' "$out" | grep -qE -- "$pattern"; then
    echo "FEHLER [${opts:-ohne Option}]: $name — Ausgabe trägt '$pattern' nicht:" >&2
    printf '%s\n' "$out" >&2
    fail=1
  fi
}

# check <Fallname> <Exit> <Muster> <Argumente…> — run und expect in einem.
check() {
  local name=$1 want_rc=$2 pattern=$3
  shift 3
  run "$@"
  expect "$name" "$want_rc" "$pattern"
}

cases() {
  # HTML-id vor einem Heading: Einheit ist der Abschnittskörper.
  check "id vor Heading, nur Heading-Zeile geändert" 0 'cmp 0' "$c0" g1.md '#gh' "$c_head" g1.md '#gh'
  check "id vor Heading, Wort im Körper (A1)" 1 'vergleich roh: .*cmp 1' "$c0" g2.md '#gh' "$c_body" g2.md '#gh'
  check "id vor Heading gegen sich selbst" 0 'cmp 0' "$c_body" g2.md '#gh' "$c_body" g2.md '#gh'
  check "Zeile der id unverändert" 0 'cmp 0' "$c0" g2.md 'L3-3' "$c_body" g2.md 'L3-3'
  check "id-Einheit gleich Einheit des Heading-Slugs mit Umlaut" 0 'cmp 0' "$c0" g2.md '#gh' "$c0" g2.md '#guard-härtung-wächter-reifen'
  check "Leerzeile im Abschnitt (A4b)" 1 'cmp 1' "$c0" g3.md '#gh' "$c_leer" g3.md '#gh'
  # HTML-id in einer Tabellenzeile: Einheit ist die Zeile.
  check "Tabellenzeile, Wort in der eigenen Zeile" 1 'cmp 1' "$c0" t.md '#mr-004' "$c_tab" t.md '#mr-004'
  check "Tabellenzeile, Wort in der Nachbarzeile" 0 'cmp 0' "$c0" t.md '#mr-003' "$c_tab" t.md '#mr-003'
  check "Tabellenzeile gegen Abschnitt (N2)" 1 'cmp 1' "$c0" t.md '#mr-004' "$c0" t.md '#tabelle'
  # HTML-id vor einem Absatz und in einer Heading-Zeile.
  check "id vor Absatz, Wort im Absatz" 1 'cmp 1' "$c0" p.md '#absatz' "$c_wort" p.md '#absatz'
  check "id in Heading-Zeile, Wort im Abschnitt" 1 'cmp 1' "$c0" h.md '#m22' "$c_wort" h.md '#m22'
  check "id in Heading-Zeile, Änderung in anderer Datei" 0 'cmp 0' "$c0" h.md '#m22' "$c_tab" h.md '#m22'
  check "id nur in Inline-Code" 2 'keine Einheit, Exit 2' "$c0" i.md '#code' "$c0" i.md '#code'
  # fail-closed.
  check "verschiedene Einheiten (A2)" 1 'cmp 1' "$c0" g2.md '#gh' "$c0" p.md '#absatz'
  check "fehlende Datei" 2 'keine Einheit, Exit 2' "$c0" fehlt.md '' "$c0" g1.md ''
  check "fehlende Datei, Meldung" 2 'nicht lesbar' "$c0" fehlt.md '' "$c0" g1.md ''
  check "unbekannter Anker links" 2 "^vergleich: $c0:g1.md#unbekannt keine Einheit, Exit 2" "$c0" g1.md '#unbekannt' "$c0" g1.md '#gh'
  check "unbekannter Anker rechts" 2 "^vergleich: $c0:g1.md#unbekannt keine Einheit, Exit 2" "$c0" g1.md '#gh' "$c0" g1.md '#unbekannt'
  check "Lokator X1-2" 2 'einheit: Lokator X1-2' "$c0" g1.md 'X1-2' "$c0" g1.md 'X1-2'
  check "Lokator L0-0" 2 'Exit 2' "$c0" g1.md 'L0-0' "$c0" g1.md 'L0-0'
  check "Tag-Paar mit Bindestrich" 2 'Tag-Paar v6.14.0-v6.14.1, Exit 2' "$c_wort" n.md '' "$c_bump" n.md '' 'v6.14.0-v6.14.1'
  check "Tag-Paar mit Suffix" 2 'Tag-Paar v6.14.0:v6.14.1x, Exit 2' "$c_wort" n.md '' "$c_bump" n.md '' 'v6.14.0:v6.14.1x'
  # roh: byte-gleich, abschließende Leerzeilen eingeschlossen.
  check "abschließende Leerzeilen der Datei (A4)" 1 'cmp 1' "$c0" r1.md '' "$c_wort" r1.md ''
  check "Leerzeile als Lokator ist eine Einheit" 0 'cmp 0' "$c_wort" r1.md 'L2-2' "$c_wort" r1.md 'L3-3'
  check "Schluss-Umbruch, Abschnitt" 0 'cmp 0' "$c0" s.md '#s' "$c_wort" s.md '#s'
  check "Schluss-Umbruch, ganze Datei" 1 'cmp 1' "$c0" s.md '' "$c_wort" s.md ''
  # Normalisierung nur des bewegten Tags; Tag-Paar an den Baseline-Bäumen.
  check "Bump roh" 1 'vergleich roh: .*cmp 1' "$c_mv" n.md '#pins' "$c_bump" n.md '#pins'
  check "Bump mit Tag-Paar" 0 'vergleich norm v6.14.0:v6.14.1: .*cmp 0' "$c_mv" n.md '#pins' "$c_bump" n.md '#pins' 'v6.14.0:v6.14.1'
  check "Bump, ganze Datei mit Tag-Paar" 0 'cmp 0' "$c_mv" n.md '' "$c_bump" n.md '' 'v6.14.0:v6.14.1'
  check "fremder Pin neben dem Bump bleibt roh (A5)" 1 'vergleich norm v6.14.0:v6.14.1: .*cmp 1' "$c_mv" fp.md '#pins' "$c_bump" fp.md '#pins' 'v6.14.0:v6.14.1'
  check "nur fremder Pin roh (N7)" 1 'cmp 1' "$c_mv" fo.md '#pins' "$c_bump" fo.md '#pins'
  check "nur fremder Pin, Bump-Paar (N7)" 1 'cmp 1' "$c_mv" fo.md '#pins' "$c_bump" fo.md '#pins' 'v6.14.0:v6.14.1'
  check "fremdes Paar ohne Baseline-Baum (N7)" 2 "$c_mv trägt .harness/baseline/v0.79.0 nicht, Exit 2" "$c_mv" fo.md '#pins' "$c_bump" fo.md '#pins' 'v0.79.0:v0.80.0'
  check "Tag-Paar mit gleichen Tags" 2 'mit gleichen Tags, Exit 2' "$c_bump" n.md '' "$c_bump" n.md '' 'v6.14.1:v6.14.1'
  check "alter Stand ohne alten Baum" 2 "$c_mv trägt .harness/baseline/v6.14.1 nicht, Exit 2" "$c_mv" n.md '' "$c_bump" n.md '' 'v6.14.1:v6.14.0'
  check "neuer Stand ohne neuen Baum" 2 "$c_mv trägt .harness/baseline/v6.14.1 nicht, Exit 2" "$c_bump" n.md '' "$c_mv" n.md '' 'v6.14.0:v6.14.1'
  check "leeres Tag-Paar" 2 'Tag-Paar , Exit 2' "$c_mv" n.md '' "$c_bump" n.md '' ''
  # Regressionsfälle: Anker-Wechsel, Lokator, git mv.
  check "Anker-Wechsel auf anderen Abschnitt" 1 'cmp 1' "$c0" ad.md '#teilfrage-2' "$c_wort" ad.md '#teilfrage-3'
  check "gleicher Anker" 0 'cmp 0' "$c0" ad.md '#teilfrage-2' "$c_wort" ad.md '#teilfrage-2'
  check "Lokator nachgezogen" 0 'cmp 0' "$c0" lz.md 'L2-4' "$c_wort" lz.md 'L5-7'
  check "Lokator nicht nachgezogen" 1 'cmp 1' "$c0" lz.md 'L2-4' "$c_wort" lz.md 'L2-4'
  check "git mv, alte Adresse am neuen Stand" 2 'keine Einheit, Exit 2' "$c_mv" mv.md '' "$c_mv" mv2.md ''
  check "git mv, alte Adresse am Parent" 0 'cmp 0' "$c_wort" mv.md '' "$c_mv" mv2.md ''
  # Heading mit Inline-Code (N1).
  check "Heading mit Inline-Code (N1)" 1 'cmp 1' "$c0" n1.md '#der-make-gates-lauf' "$c_wort" n1.md '#der-make-gates-lauf'
  # Fence nach CommonMark (F-3).
  check "Fence aus vier Backticks, ein innerer Fence (F-3)" 1 'cmp 1' "$c0" fe.md '#ungerade' "$c_wort" fe.md '#ungerade'
  check "Abschnitt hinter dem Fence (F-3)" 0 'cmp 0' "$c0" fe.md '#ende' "$c_wort" fe.md '#ende'
  check "Fence aus vier Backticks, zwei innere Fences (N6)" 1 'cmp 1' "$c0" fe2.md '#gerade' "$c_wort" fe2.md '#gerade'
  check "~~~ in einem Fence aus Backticks" 1 'cmp 1' "$c0" ft.md '#tilde' "$c_wort" ft.md '#tilde'
  check "Backtick im Info-String öffnet keinen Fence" 0 'cmp 0' "$c0" fb.md '#eins' "$c_wort" fb.md '#eins'
  check "Fence direkt nach id-Zeile gehört zur Einheit" 0 'cmp 0' "$c0" fz.md '#fz' "$c0" fz.md 'L3-7'
  # Gestapelte id (F-1).
  check "gestapelte id, obere (F-1)" 1 'cmp 1' "$c0" st.md '#oben' "$c_wort" st.md '#oben'
  check "gestapelte id, untere (F-1)" 1 'cmp 1' "$c0" st.md '#unten' "$c_wort" st.md '#unten'
  check "gestapelte id adressiert den Abschnitt" 0 'cmp 0' "$c0" st.md '#oben' "$c0" st.md '#abschnitt'
  check "gestapelte id vor Absatz, Einheit ist der Block" 0 'cmp 0' "$c0" sp.md '#s1' "$c0" sp.md 'L4-5'
  check "gestapelte id vor Absatz, Wort im Absatz" 1 'cmp 1' "$c0" sp.md '#s1' "$c_wort" sp.md '#s1'
  # id in anderer Form (F-2).
  check "id mit Attribut (F-2)" 2 'anderer Form .*Exit 2|keine Einheit, Exit 2' "$c0" fa.md '#attr' "$c0" fa.md '#attr'
  check "id mit Attribut, Meldung (F-2)" 2 '<a id="attr" in anderer Form' "$c0" fa.md '#attr' "$c0" fa.md '#attr'
  check "id selbstschließend (F-2)" 2 '<a id="selbst" in anderer Form' "$c0" fs.md '#selbst' "$c0" fs.md '#selbst'
  check "id mit Attribut neben gleichnamigem Heading-Slug (F-2)" 2 '<a id="doppel" in anderer Form' "$c0" fh.md '#doppel' "$c0" fh.md '#doppel'
  check "andere Form in Inline-Code bleibt ungelesen" 0 'cmp 0' "$c0" fi.md '#fi' "$c0" fi.md '#fi'
  # Fence-Einzug: bis 3 Leerzeichen ein Fence, ab 4 keiner.
  check "Fence mit 3 Leerzeichen Einzug" 1 'cmp 1' "$c0" fi3.md '#eins' "$c_wort" fi3.md '#eins'
  check "4 Leerzeichen Einzug öffnet keinen Fence" 0 'cmp 0' "$c0" fi4.md '#eins' "$c_wort" fi4.md '#eins'
  # Normalisierung nur am Segment /<tag>/ mit maskierten Punkten.
  check "Tag ohne Segment-Schrägstrich bleibt roh" 1 'vergleich norm v6.14.0:v6.14.1: .*cmp 1' "$c_mv" seg.md '#pins' "$c_bump" seg.md '#pins' 'v6.14.0:v6.14.1'
  check "Punkt im Tag ist kein Platzhalter" 0 'vergleich norm v6.14.0:v6.14.1: .*cmp 0' "$c_mv" dot.md '#pins' "$c_bump" dot.md '#pins' 'v6.14.0:v6.14.1'
  # Slug: ohne HTML-Tags, Dubletten-Suffix, Ebene höchstens 6, schließende #-Folge.
  check "Slug ohne HTML-Tags" 1 'cmp 1' "$c0" code.md '#code-x' "$c_wort" code.md '#code-x'
  check "Dublette -1" 1 'cmp 1' "$c0" dup.md '#dup-1' "$c_wort" dup.md '#dup-1'
  check "Dublette ohne Suffix" 0 'cmp 0' "$c0" dup.md '#dup' "$c_wort" dup.md '#dup'
  check "Dublette -2" 0 'cmp 0' "$c0" dup.md '#dup-2' "$c_wort" dup.md '#dup-2'
  check "sieben # sind kein Heading" 2 'leere Einheit' "$c0" lvl.md '#sieben' "$c0" lvl.md '#sieben'
  check "Ebene 6 ist ein Heading" 1 'cmp 1' "$c0" lvl.md '#sechs' "$c_wort" lvl.md '#sechs'
  check "schließende #-Folge gehört nicht zum Slug" 1 'cmp 1' "$c0" close.md '#eins' "$c_wort" close.md '#eins'
  check "Slug mit Bindestrich der #-Folge löst nicht auf" 2 'leere Einheit' "$c0" close.md '#eins-' "$c0" close.md '#eins-'
  # id in eingerücktem Code und in HTML-Kommentar wird nicht gelesen.
  check "id in eingerücktem Code zählt nicht" 1 'cmp 1' "$c0" indent.md '#ein' "$c_wort" indent.md '#ein'
  check "id in mehrzeiligem HTML-Kommentar zählt nicht" 1 'cmp 1' "$c0" kom.md '#kom' "$c_wort" kom.md '#kom'
  check "id in einzeiligem HTML-Kommentar zählt nicht" 1 'cmp 1' "$c0" kom.md '#k2' "$c_wort" kom.md '#k2'
  check "<!-- in Inline-Code öffnet keinen Kommentar" 1 'cmp 1' "$c0" ik.md '#k3' "$c_wort" ik.md '#k3'
  # Mehrdeutige Stellung endet mit Exit 2.
  check "Kommentar hinter einzelnem Backtick ist mehrdeutig (n4)" 2 'm4" mehrdeutig' "$c0" n4.md '#m4' "$c_wort" n4.md '#m4'
  check "Kommentar neben Doppel-Backtick-Span ist mehrdeutig (n5)" 2 'm5" mehrdeutig' "$c0" n5.md '#m5' "$c_wort" n5.md '#m5'
  check "Kommentar im Absatz mit offenem Code-Span ist mehrdeutig (n6)" 2 'm6" mehrdeutig' "$c0" n6.md '#m6' "$c_wort" n6.md '#m6'
  check "id nach unsicherer Kommentar-Grenze ist mehrdeutig" 2 'm7" mehrdeutig' "$c0" un.md '#m7' "$c_wort" un.md '#m7'
  check "Einzug aus Leerzeichen und Tab ist mehrdeutig" 2 'm3" mehrdeutig' "$c0" mi.md '#m3' "$c_wort" mi.md '#m3'
  check "Tabellenzeile mit Inline-Code bleibt gelesen" 1 'cmp 1' "$c0" tb.md '#mr-009' "$c_wort" tb.md '#mr-009'
  check "Kommentar in der id-Zeile vor Heading, Einheit ist der Abschnitt" 1 'cmp 1' "$c0" hc.md '#m8' "$c_wort" hc.md '#m8'
  # id in einem Code-Span mit Text davor wird nicht gelesen.
  check "id in Code-Span im Absatz zählt nicht" 1 'cmp 1' "$c0" e1.md '#q1' "$c_wort" e1.md '#q1'
  check "id und Kommentar in Code-Span zählen nicht" 1 'cmp 1' "$c0" e2.md '#q2' "$c_wort" e2.md '#q2'
  check "id in Code-Span einer Tabellenzelle zählt nicht" 1 'cmp 1' "$c0" e6.md '#q6' "$c_wort" e6.md '#q6'
  check "id nur in Code-Span mit Text davor" 2 'keine Einheit, Exit 2' "$c0" e7.md '#q7' "$c0" e7.md '#q7'
  check "id nach geschlossenem Code-Span wird gelesen" 1 'cmp 1' "$c0" e8.md '#q8' "$c_wort" e8.md '#q8'
  check "Lokator mit führenden Nullen ist dezimal" 0 'cmp 0' "$c0" l12.md 'L010-012' "$c0" l12.md 'L10-12'
  # Lokator strikt.
  check "Lokator ohne Bindestrich" 2 'Lokator L7 \(Form' "$c0" lz.md 'L7' "$c0" lz.md 'L7-7'
  check "Lokator mit drittem Teil" 2 'Lokator L1-2-3 \(Form' "$c0" lz.md 'L1-2-3' "$c0" lz.md 'L1-2'
  check "Lokator a > b" 2 'Lokator L5-3 \(1 <= a <= b\)' "$c0" lz.md 'L5-3' "$c0" lz.md 'L5-5'
  check "Lokator ab 0" 2 'Lokator L0-1 \(1 <= a <= b\)' "$c0" lz.md 'L0-1' "$c0" lz.md 'L1-1'
  # source ändert die Shell-Optionen des Aufrufers nicht.
  out=$(env ${opts:+BASHOPTS="$opts"} bash -c 'v=$-; p=$(set -o | grep "^pipefail"); source "$1"; if [ "$v" = "$-" ] && [ "$p" = "$(set -o | grep "^pipefail")" ] && declare -F einheit >/dev/null; then echo "Optionen unverändert: $-"; else echo "Optionen geändert: $v -> $-"; exit 1; fi' x "$prog" 2>&1)
  rc=$?
  expect "source lässt die Shell-Optionen unverändert" 0 'Optionen unverändert'
  # Locale und Fähigkeitsprobe (F-6).
  out=$(env ${opts:+BASHOPTS="$opts"} LC_ALL=C bash "$prog" "$c0" g2.md '#guard-härtung-wächter-reifen' "$c0" g2.md '#guard-härtung-wächter-reifen' 2>&1)
  rc=$?
  expect "Umlaut-Slug unter LC_ALL=C beim Aufrufer (L)" 0 'cmp 0'
  mkdir -p "$tmp/stub"
  printf '#!/bin/sh\nLC_ALL=C exec %s "$@"\n' "$real_awk" >"$tmp/stub/awk"
  chmod +x "$tmp/stub/awk"
  out=$(env ${opts:+BASHOPTS="$opts"} PATH="$tmp/stub:$PATH" bash "$prog" "$c0" g1.md '' "$c0" g1.md '' 2>&1)
  rc=$?
  expect "awk ohne Multibyte" 2 'awk ohne Multibyte/UTF-8'
  # Argumente.
  check "fünf Argumente" 2 '^Aufruf: .*5 Argumente' "$c0" g1.md '#gh' "$c0" g1.md
  check "acht Argumente" 2 '^Aufruf: .*8 Argumente' "$c0" g1.md '#gh' "$c0" g1.md '#gh' 'v6.14.0:v6.14.1' x
  # set -euo pipefail beim Aufrufer: die Zeile steht vor dem Abbruch.
  out=$(env ${opts:+BASHOPTS="$opts"} SHELLOPTS=errexit:nounset:pipefail bash "$prog" "$c0" g2.md '#gh' "$c_body" g2.md '#gh' 2>&1)
  rc=$?
  expect "roter Vergleich unter set -euo pipefail" 1 'vergleich roh: .*cmp 1'
}

for opts in "" nullglob failglob; do
  cases
done

if [ "$fail" -eq 0 ]; then
  echo "run-zitat-vergleich-tests: $count Fälle bestanden (je Runde $((count / 3)), Runden: ohne Option, nullglob, failglob)"
else
  echo "run-zitat-vergleich-tests: mindestens ein Fall gescheitert ($count Fälle)" >&2
  exit 1
fi
