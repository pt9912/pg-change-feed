#!/usr/bin/env bash
# handbuch-public-doc-check.sh — prueft, dass die Nutzerdokumente unter
# docs/user/ keine interne Kennung und keinen Link nach docs/plan/ oder
# docs/reviews/ tragen. Die Nutzerdokumentation ist die Betreiber- und
# Integrator-Sicht; Kennungen der Spezifikations-, Entscheidungs- und
# Anforderungsdokumente, Slice-/Welle-/Register-Namen gehoeren in Plan und
# Commit, nicht in das Handbuch.
#
# Reichweite: zwei benannte Listen plus Vollstaendigkeitspruefung. Geprueft
# werden die Dateien der Liste `checked`; `excluded` nennt die von Runnern
# geschriebenen Erzeugnisse, deren Kennungen gewollt sind (die Maintainer-Doku
# liegt unter docs/maintainer/, ausserhalb des Gegenstands). Eine *.md unter docs/user/,
# die in keiner Liste steht, und eine genannte, nicht vorhandene Datei enden
# mit Exit 2: eine neue Datei erzwingt die Klassifikation.
#
# Aufruf: `make handbuch-public-doc-check` (netzlos, nur bash/git/grep).
# Optional eine abweichende Wurzel als erstes Argument (fuer den Tabellentest
# run-handbuch-public-doc-check-tests.sh). Exit 0 ohne Treffer, Exit 1 mit
# Treffern (Ausgabe `datei:zeile: text`), Exit 2 bei Klassifikationsfehlern und
# bei Lesefehlern (nicht lesbare Datei oder Verzeichnis; nie "sauber").
set -euo pipefail

if [ "$#" -ge 1 ]; then
  root=$1
else
  root=$(git rev-parse --show-toplevel)
fi
dir="$root/docs/user"

checked=(benutzerhandbuch.md benutzerhandbuch-standard.md version.md)
excluded=(bench-abdeckung.md ci-matrix-abdeckung.md
          e2e-abdeckung.md sdk-e2e-abdeckung.md)

# Muster P: Kennungen; Muster L: Links nach docs/plan/ und docs/reviews/.
pat_id='\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]'
pat_link='docs/(reviews|plan)/|\.\./(reviews|plan)/'

scratch=$(mktemp -d)
trap 'rm -rf "${scratch:?}"' EXIT

# find laeuft vor der Schleife: ein nicht lesbares Verzeichnis ist Exit 2, nicht
# eine stillschweigend kuerzere Liste.
frc=0
find "$dir" -type f -name '*.md' -print0 > "$scratch/list.raw" 2> "$scratch/find.err" || frc=$?
if [ "$frc" -ne 0 ]; then
  echo "handbuch-public-doc-check: Lesefehler beim Durchsuchen von docs/user/ ($(head -n 1 "$scratch/find.err"))" >&2
  exit 2
fi
sort -z "$scratch/list.raw" > "$scratch/list"

class_err=0
for f in "${checked[@]}" "${excluded[@]}"; do
  if [ ! -f "$dir/$f" ]; then
    echo "handbuch-public-doc-check: genannte Datei fehlt: docs/user/$f" >&2
    class_err=1
  fi
done

while IFS= read -r -d '' path; do
  rel=${path#"$dir"/}
  known=0
  for f in "${checked[@]}" "${excluded[@]}"; do
    if [ "$rel" = "$f" ]; then known=1; fi
  done
  if [ "$known" -eq 0 ]; then
    echo "handbuch-public-doc-check: unklassifizierte Datei: docs/user/$rel (in checked oder excluded aufnehmen)" >&2
    class_err=1
  fi
done < "$scratch/list"

if [ "$class_err" -ne 0 ]; then
  exit 2
fi

hits=""
for f in "${checked[@]}"; do
  # -a: eine Datei mit NUL-Byte wird als Text gelesen statt still uebersprungen.
  # grep-Exit 1 = kein Treffer (Erfolg), >= 2 = Lesefehler (Exit 2).
  grc=0
  h=$(grep -anE -e "$pat_id" -e "$pat_link" "$dir/$f" 2>"$scratch/grep.err") || grc=$?
  if [ "$grc" -ge 2 ]; then
    echo "handbuch-public-doc-check: Lesefehler: docs/user/$f ($(head -n 1 "$scratch/grep.err"))" >&2
    exit 2
  fi
  if [ -n "$h" ]; then
    hits+=$(printf '%s\n' "$h" | sed "s|^|docs/user/$f:|")$'\n'
  fi
done

if [ -n "$hits" ]; then
  printf '%s' "$hits" >&2
  echo "handbuch-public-doc-check: interne Kennung oder Link nach docs/plan/ bzw. docs/reviews/ in den Nutzerdokumenten (siehe oben)" >&2
  exit 1
fi
echo "handbuch-public-doc-check: keine interne Kennung in ${#checked[@]} Nutzerdokumenten unter docs/user/"
