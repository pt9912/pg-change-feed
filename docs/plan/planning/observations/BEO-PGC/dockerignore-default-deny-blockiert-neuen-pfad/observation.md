# `.dockerignore`-Default-Deny blockiert einen neu benötigten Pfad

**Sub-Area:** `*`/`PGC` (Build-/Gate-Infrastruktur, repo-weite Default-Sub-Area)

Die Beobachtung: Eine neue Docker-Multi-Stage-Stufe, die einen bisher **nicht**
in `.dockerignore`s Allow-Liste geführten Pfad braucht — Verzeichnis oder
Datei, unabhängig davon, ob es ein Werkzeug-Skript, ein Quellverzeichnis oder
etwas anderes ist —, scheitert real im Build: Der Pfad landet nicht im
Build-Kontext, obwohl er im Arbeitsbaum liegt (`COPY … : "…": not found`).

**Abgrenzung zu `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`:**
Derselbe zugrunde liegende Mechanismus (`.dockerignore`s Allow-Listen-
Default-Deny bricht bei jeder neuen Docker-Stufe, die einen bisher
ungelisteten Pfad braucht) — aber jener Eintrag ist über **Pfad und Text**
eng auf „eine neue Docker-Stage, die ein **Skript unter `tools/`** braucht"
gezogen (beide bisherigen Belege, `slice-049`/`slice-093`, betreffen genau
das, plus eine bei beiden mitlaufende, hier unabhängige Alpine/`bash`-Facette;
`state.md` dort benennt den künftigen Trigger wörtlich als „ein Skript unter
`tools/`"). Dieser Eintrag führt den **generalisierten** Fall — ein
beliebiger, bisher ungelisteter Pfad, der kein Skript unter `tools/` sein
muss (Erstbeleg: ein Quellverzeichnis für eine neue Erzeugungsstufe). Beide
Klassen bleiben bewusst getrennt, damit die aus `…-blockiert-tooling`
abgeleitete Regel („bei einem Skript unter `tools/` vorab prüfen") nicht auf
einen Fall gedehnt wird, den ihr eigener Wortlaut nicht trägt (Reviewer- und
Verifier-Urteil zu `slice-104` F-2, unabhängig übereinstimmend).
