# OCI-Image für PG Change Feed — Reproduzierbarkeits-Anker (Modul 14).
# Rezept-Form: Quellen-Commit + die gepinnten Digests der Eingangs-Images
# bauen die Umgebung neu; der Digest des gebauten Images landet via
# `make image` in harness/image-hash.txt.
# Semantik des Belegs: er ist der Digest des **exportierten Images**
# (Runtime-Stage: Distroless-Basis + Binary). Änderungen, die nur den
# deps-Layer betreffen (go.mod/go.sum, solange das Binary den Import nicht
# trägt), lassen den Digest unverändert — das ist ein gültiger Befund
# ("Beleg, kein Wiederholungs-Schlüssel"), kein Staleness-Zeichen.
# Base-Image-Update = bewusster Commit, der nur die Digest-Zeile anhebt.

# go.mod/go.sum liegen seit dem Go-Modul-Bootstrap im Baum; die
# PostgreSQL-Abhängigkeit (pgx/v5, rein Go, CGO-frei) trägt der deps-Layer
# per `go mod download`/`go mod verify` aus dem gepinnten Stand.

# --- deps: gepinnte Base, Lock-File vor dem Code (Layer-Cache greift) ---
FROM golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS deps
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify

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