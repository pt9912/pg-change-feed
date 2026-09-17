# BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Form, in
der ein Implementer-Zug eine Lösung für ein bekanntes Problem entwirft, ohne
zuerst zu prüfen, ob derselbe Bestand (dieselbe Datei, derselbe Ordner) das
Problem bereits an anderer Stelle löst).

Die Beobachtung: Ein Slice steht vor einer Shell-/Werkzeug-Einschränkung
(hier: `make`-Rezepte laufen unter `/bin/sh`, auf diesem Host `dash`, das
`set -o pipefail` nicht trägt — eine Pipe `docker run | tar -x` würde einen
`docker run`-Fehlschlag hinter `tar`s Exit-Code verstecken, `AGENTS.md`
§3.9). Der Implementer löste das korrekt, aber umständlich: `docker run`
schreibt in eine Zwischendatei, `tar -xf` extrahiert getrennt, `rm -f` räumt
auf — drei Rezept-Zeilen statt einer. Die einfachere Lösung stand im selben
Makefile bereits neunfach etabliert: `@bash tools/harness/*.sh` ruft für
jedes andere Gate-Skript explizit `bash` auf (nicht `/bin/sh`), wodurch
`set -o pipefail` real trägt und eine echte Pipe sicher wird. Ein einziger
`grep bash Makefile` hätte das gezeigt.

Aufgefallen ist es nicht im Review/Verifikations-Schritt (beide prüften nur,
ob die gewählte Lösung *korrekt* ist — was sie war — nicht, ob eine
*einfachere, bereits im selben Bestand etablierte* Lösung existiert), sondern
erst im Gespräch mit dem Nutzer nach der Slice-Closure.

Deklaration: `slice-104` (Commit `338cfe0`, korrigiert in `99b1af9`,
Orchestrator-Dialog 2026-09-17).
