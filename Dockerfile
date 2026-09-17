# OCI-Image für PG Change Feed — Reproduzierbarkeits-Anker (Modul 14).
# Rezept-Form: Quellen-Commit + die gepinnten Digests der Eingangs-Images
# bauen die Umgebung neu; der Digest des gebauten Images landet via
# `make image` in harness/image-hash.txt.
# Semantik des Belegs (ADR-0044): der Digest ist der **Lauf-Beleg** des
# letzten offiziellen `make image`-Laufs — builder- und lauf-gebunden,
# kein Inhalts-Fingerabdruck; ein Digest-Vergleich über Umgebungen oder
# Läufe entscheidet Staleness nicht. Inhalts-Streits werden über den
# sha256 des extrahierten Binaries entschieden (Container-Export;
# Verfahren: docs/reviews/verify-slice-004.md F-2-Schiedsspruch,
# docs/reviews/review-slice-005.md F-7).
# Base-Image-Update = bewusster Commit, der nur die Digest-Zeile anhebt.

# go.mod/go.sum liegen seit dem Go-Modul-Bootstrap im Baum; die
# PostgreSQL-Abhängigkeit (pgx/v5, rein Go, CGO-frei) trägt der deps-Layer
# per `go mod download`/`go mod verify` aus dem gepinnten Stand.

# --- deps: gepinnte Base, Lock-File vor dem Code (Layer-Cache greift) ---
FROM golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS deps
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# --- proto: Docker-only Toolchain-Stufe der Protobuf-/gRPC-Codegenerierung
# (`ADR-0060` Folgepflicht, `LH-FA-SST-008`). `protoc` und die beiden
# `protoc-gen-*`-Plugins laufen ausschliesslich hier (`AGENTS.md` §3.1: kein
# Host-`protoc`/`buf`); die Plugin-Versionen sind gepinnt, die Basis ist der
# bereits digest-gepinnte Toolchain-Stand aus `deps`. Diese Stufe traegt nur
# das Werkzeug; die eigentliche Erzeugung liegt in der Folgestufe
# `proto-export`. Der normale Build (`make image`) braucht sie nicht — der
# erzeugte Go-Code liegt committet im Baum und wird von `build`/`coverage`
# mitkompiliert. `make generated-sync` baut diese Stufe eigenstaendig
# (`tools/harness/generated-sync.sh`) und vergleicht in einem eigenen
# Temp-Verzeichnis, ohne den Arbeitsbaum zu schreiben. ---
FROM deps AS proto
RUN apk add --no-cache protobuf-dev=31.1-r1 \
 && GOBIN=/usr/local/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.12 \
 && GOBIN=/usr/local/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2

# --- proto-export: erzeugt den Go-/gRPC-Code zur Build-Zeit (seit slice-104;
# vormals `docker run -v` in den Bind-Mount des Arbeitsbaums, siehe
# ADR-0060). Die `.proto`-Quelle kommt per `COPY` (Build-Kontext) statt per
# Mount in die Stufe; `protoc` laeuft als `RUN`-Schritt, das Erzeugnis liegt
# im Image-Layer unter `/out`. `tools/harness/proto-generate.sh` baut diese
# Stufe und liest ihr `ENTRYPOINT` per `docker run --rm --network none
# <image> | tar -x -C .` aus: Es gibt `/out` als `tar`-Stream ueber stdout
# aus, die Extraktion laeuft host-seitig — kein `docker run -v`, kein
# `--user`-Workaround, die extrahierten Dateien gehoeren dadurch automatisch
# dem aufrufenden Nutzer (Host-Prozess, kein Container-Schreibzugriff auf
# den Baum). Die Pipe laeuft unter `bash`/`set -o pipefail` im Skript, nicht
# unter dem `make`-Default `/bin/sh` (`AGENTS.md` §3.9). ---
FROM proto AS proto-export
COPY proto/ proto/
RUN mkdir -p /out && \
    protoc -I proto \
      --go_out=/out --go_opt=module=github.com/pt9912/pg-change-feed \
      --go-grpc_out=/out --go-grpc_opt=module=github.com/pt9912/pg-change-feed \
      proto/cdc/stream/v1/changestream.proto
ENTRYPOINT ["tar", "-cf", "-", "-C", "/out", "."]

# --- coverage: Go-Test-Coverage ueber die netzlos pruefbare Flaeche des
# Baums und Gate-Skript gegen COVERAGE_THRESHOLD (ADR-0071, ADR-0054;
# Kalibrierungs-Bindung harness/README.md §Sensors). Die Flaeche ist
# internal/...+cmd/...+gen/... ohne die Pakete, deren Testlauf einen externen
# Dienst voraussetzt — der Filter unten nennt die drei namentlich, die
# tragende Regel ist die Eigenschaft, nicht die Liste; die Paketliste selbst
# kommt aus `go list` und zieht neue Pakete mit. `gen/...` ist seit
# `slice-097` Teil der Liste: der Umzug der erzeugten Vertragsflaeche
# (`ADR-0076`) bewegt den Traeger, nicht den Gegenstand — die Eigenschaft
# (netzlos pruefbar, `ADR-0071` Punkt 1) bleibt erfuellt, siehe
# architect-verdict-slice-097-coverage-gegenstand.md. `-coverpkg` misst ueber
# die Paketgrenzen hinweg, sonst zaehlt nur paket-lokale Abdeckung.
# `test/integration/` bleibt ausgeschlossen (eigene Black-Box-Paketwurzel
# gegen einen laufenden Compose-Container, kein Unit-Coverage-Kandidat).
# Adapter-Tests ohne gesetzte CDC_*_TEST_DSN-Variable skippen real in diesem
# netzlosen Lauf, ohne den Build zu brechen. `pipefail` via SHELL, damit
# `go test … | tee` bzw. `go tool cover … | tee` den Exit-Code nicht maskiert. ---
FROM deps AS coverage

# golang:1.27-alpine traegt kein bash (anders als d-checks Debian-basierte
# golang:${GO_VERSION}); die SHELL-Direktive unten braucht ein installiertes
# bash, also VOR der Umstellung installieren (noch mit dem Alpine-Default-sh).
RUN apk add --no-cache bash

SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

ARG COVERAGE_THRESHOLD
ENV COVERAGE_THRESHOLD=${COVERAGE_THRESHOLD}

COPY . .
RUN mkdir -p /out && \
    pkgs=( $(go list ./internal/... ./cmd/... ./gen/... \
        | grep -vE '(^|/)(postgresstorage|postgresack|replication/receive)$') ) && \
    go test \
        -coverpkg="$(IFS=,; echo "${pkgs[*]}")" \
        -coverprofile=/out/coverage.out \
        -covermode=atomic \
        "${pkgs[@]}" && \
    go tool cover -func=/out/coverage.out | tee /out/coverage-func.txt && \
    bash tools/coverage-gate.sh /out/coverage-func.txt "$COVERAGE_THRESHOLD"

# --- build: Kompilierung getrennt vom Cache-sensiblen Layer; CGO aus
# (ADR-0042: die Struktur-Regeln der abgeloesten Kette ADR-0038/0039 bleiben fortgeltend) ---
FROM deps AS build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
      -o /out/pg-change-feed ./cmd/pg-change-feed

# --- runtime: distroless, nonroot, nur Artefakte — keine Shell, kein Paketmanager ---
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime
COPY --from=build /out/pg-change-feed /pg-change-feed
USER nonroot
ENTRYPOINT ["/pg-change-feed"]
