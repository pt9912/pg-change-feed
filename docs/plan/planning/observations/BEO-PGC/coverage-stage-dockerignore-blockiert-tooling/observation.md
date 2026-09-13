# BEO-PGC/coverage-stage-dockerignore-blockiert-tooling

**Sub-Area:** Build-/Gate-Infrastruktur (Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: Eine neue Docker-Multi-Stage-Stufe, die ein zusätzliches
Skript unter `tools/` braucht (hier: `tools/coverage-gate.sh`), scheitert
real im Build, weil `.dockerignore` den Kontext auf `cmd/`, `internal/`,
`go.mod`, `go.sum` einschränkt — das Skript landet nicht im Build-Kontext,
obwohl es im Arbeitsbaum liegt, und die Stage bricht mit „No such file or
directory" ab. Zusätzlich trägt die `golang:1.27-alpine`-Basis (im
Unterschied zu d-checks Debian-basierter `golang:${GO_VERSION}`) kein
`bash` — die für `tee`-Exit-Code-Transparenz nötige `SHELL
["/bin/bash", …]`-Direktive scheitert ohne ein vorgeschaltetes `apk add
--no-cache bash` (mit dem Alpine-Default-`sh`, bevor die SHELL-Direktive
umschaltet). Beide Fallstricke traten beim direkten Kopieren des
d-check-Musters auf ein Alpine-basiertes Repo auf und wurden erst durch
den realen, tatsächlich ausgeführten Build sichtbar (nicht durch
Code-Lesen).
