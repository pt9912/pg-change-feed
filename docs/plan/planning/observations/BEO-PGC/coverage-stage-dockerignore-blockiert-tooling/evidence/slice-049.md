# Beleg: slice-049 (Test-Coverage-Gate)

Vorgang: Implementierung von `slice-049` (vierte Docker-Stage `coverage`,
`tools/coverage-gate.sh`,
[ADR-0054](../../../../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)).

Fund 1: Der erste Build-Versuch (`docker build --target coverage`)
scheiterte mit `bash: tools/coverage-gate.sh: No such file or directory`
— `.dockerignore` schloss `tools/` aus (nur `cmd/`, `internal/`,
`go.mod`, `go.sum` waren freigegeben). Behoben durch eine gezielte
`!tools/coverage-gate.sh`-Ausnahme.

Fund 2: Vor dieser Korrektur scheiterte ein noch früherer Versuch mit
`runc run failed: … exec: "/bin/bash": stat /bin/bash: no such file or
directory` — `golang:1.27-alpine` (dieses Repos Basis-Image) trägt kein
`bash`, anders als d-checks Debian-basierte `golang:${GO_VERSION}`.
Behoben durch `RUN apk add --no-cache bash` vor der
`SHELL ["/bin/bash", …]`-Direktive.

Quelle: realer `docker build`-Lauf während der Implementierung,
2026-09-13.
