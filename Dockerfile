# OCI-Image für PG Change Feed — Reproduzierbarkeits-Anker (Modul 14).
# Rezept-Form: Quellen-Commit + die gepinnten Digests der Eingangs-Images
# bauen die Umgebung neu; der Digest des gebauten Images landet via
# `make image` in harness/image-hash.txt (Beleg, kein Wiederholungs-Schlüssel).
# Base-Image-Update = bewusster Commit, der nur die Digest-Zeile anhebt.

# go.mod/go.sum liegen seit dem Go-Modul-Bootstrap im Baum; solange keine
# externen Dependencies eingehen, bleibt go.sum leer und der deps-Layer
# verifiziert eine leere Modulliste (`go mod download`/`go mod verify`).

# --- deps: gepinnte Base, Lock-File vor dem Code (Layer-Cache greift) ---
FROM golang:1.26-alpine@sha256:ce864e7223ac17b1775e6fd0b4c0db580c2eb50e7953a427916379e4b92a1628 AS deps
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# --- build: Kompilierung getrennt vom Cache-sensiblen Layer; CGO aus (ADR-0038) ---
FROM deps AS build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
      -o /out/pg-change-feed ./cmd/pg-change-feed

# --- runtime: distroless, nonroot, nur Artefakte — keine Shell, kein Paketmanager ---
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime
COPY --from=build /out/pg-change-feed /pg-change-feed
USER nonroot
ENTRYPOINT ["/pg-change-feed"]