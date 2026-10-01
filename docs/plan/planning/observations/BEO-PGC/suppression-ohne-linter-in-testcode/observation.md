# BEO-PGC/suppression-ohne-linter-in-testcode

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft Test-Code in den
Nicht-Go-Bäumen, vor allem `sdks/python/`, `sdks/csharp/`, `sdks/kotlin/`).

Die Beobachtung: Ein Lauf schreibt eine Inline-Suppression in Test-Code einer Sprache, deren
Name die Regel nicht nennt. `AGENTS.md` §3.2 verbietet `//nolint` ausnahmslos, weil das Repo
keinen Linter hat; das Python-`# noqa: BLE001` in `route_scenario.py` unterdrückt ebenfalls
keine Warnung eines vorhandenen Werkzeugs (`git grep -n -i -E "ruff|flake8|pylint" -- sdks
Makefile harness/mk tools` ohne Treffer). Der Wortlaut der Regel nennt nur die Go-Form, ihre
Begründung gilt für jede Sprache.

**Warum das zählt:** Eine Suppression ohne Linter ist tote Dekoration oder ein Vorgriff auf ein
Profil, das niemand entschieden hat; der Wortlaut „`//nolint`“ lässt die Lesart zu, die Regel
gelte nur für Go. Kein Sensor liest Suppression-Zeichenfolgen in Nicht-Go-Bäumen.

**Abgrenzung.** `BEO-PGC/kommentar-herkunft-als-kette` zählt Kennungen im Kommentar; hier ist
der Kommentar selbst die Suppression.

## Benannt, nicht gezählt

Keine.
