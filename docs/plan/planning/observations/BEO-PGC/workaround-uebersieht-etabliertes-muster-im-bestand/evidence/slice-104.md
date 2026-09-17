# Beleg: slice-104

Vorgang: `slice-104` (mount-loser Umbau von `make proto-generate`).

Fund: Implementer, Reviewer und Verifier akzeptierten übereinstimmend eine
Zwischendatei-Lösung (`docker run > Datei; tar -xf; rm -f`) gegen die
Pipe-Exit-Code-Maskierung (`AGENTS.md` §3.9), weil das Rezept unter dem
`make`-Default `/bin/sh` (`dash`, kein `pipefail`) läuft. Alle drei Rollen
prüften die Shell-Frage isoliert ("welche Shell, trägt sie `pipefail`"),
keine prüfte, ob dasselbe Makefile das Problem bereits an anderer Stelle
löst — es tut es: neun Ziele rufen bereits explizit `@bash
tools/harness/*.sh` auf. Der Fund kam erst nach der Slice-Closure im
Orchestrator-Dialog mit dem Nutzer auf, der direkt nach den bestehenden
`bash`-Aufrufen fragte. Korrigiert in Commit `99b1af9`:
`tools/harness/proto-generate.sh` (bash, `set -euo pipefail`, echte Pipe),
Makefile ruft `bash tools/harness/proto-generate.sh`.

Quelle: `Makefile` vor/nach Commit `338cfe0`/`99b1af9`; Orchestrator-Dialog
2026-09-17.
