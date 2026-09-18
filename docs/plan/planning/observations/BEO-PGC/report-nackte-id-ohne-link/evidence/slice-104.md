**Vorgang:** slice-104

**Fund:** Der Review-Report zu `slice-104` trug an Zeile 34 eine nackte
`ADR-0060`-Kennung ohne Link/Backticks im Fließtext (Zitat einer
Dockerfile-Kommentarzeile: „…siehe `ADR-0060`" — im ursprünglichen Report
folgte die schließende Anführung direkt auf die nackte Kennung, ohne
Backticks; alle sechs anderen `ADR-0060`-Stellen
desselben Reports standen korrekt in Backticks). Der Verifier fand den Fund
über seinen eigenen `make gates`-Lauf (`docs-check` brach mit `id-unlinked`
ab, Exit 2) und dokumentierte ihn als Blocker, bevor die DoD-Zeile
„`make gates` grün" angehakt werden durfte. Behoben in einem eigenen Commit
(`aabbe01`, Backticks ergänzt) außerhalb der reinen
Implementer→Reviewer→Verifier-Sequenz, danach vom Planner mit einem eigenen
`make gates`-Lauf (Exit 0) bestätigt.

**Abgrenzung zu den Belegen 4/5 (Sequenzierungs-Klasse):** Hier lag **keine**
Sequenzierungs-Verletzung vor — der rote Exit-Code wurde korrekt gemessen und
blockierte die Closure tatsächlich, bis der Fund behoben war; der `git mv`
nach `done/` erfolgte erst danach. Der Fund trifft stattdessen den in
`observation.md` ursprünglich benannten **Basis**-Fehler: eine nackte Kennung
im Fließtext eines neu geschriebenen Review-/Verifikationsberichts. Zählt als
6. Beleg; löst **keinen** neuen Lese-Schritt aus — die Regel ist bereits seit
`slice-063` in `AGENTS.md` §3.9 verkörpert, dieser Beleg bestätigt sie.

Quelle: Verifikationsbericht zu `slice-104` §1 (Lauf 6), §4 · Commit `aabbe01`
· Planner-Closure-Entscheidung.
