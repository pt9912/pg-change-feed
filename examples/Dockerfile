# examples/Dockerfile — Bau-/Startform der vier Go-Beispiele (ADR-0098
# Festlegung 1, Supersedes ADR-0076 Festlegung 1/Startform-Bullet). Bau-Kontext
# ist die Repo-Wurzel (`docker build -f examples/Dockerfile .`), NICHT
# `examples/` selbst — die vier Go-Programme sind Wurzelmodul-Code
# (`go.mod`/`go.sum`/`gen/`) und liegen außerhalb eines Sprach-Wurzel-
# Verzeichnisses, anders als bei C#/Kotlin (ADR-0087 Festlegung 3).
#
# Der Bau-Kontext ist über die dateispezifische Ignore-Datei
# `examples/Dockerfile.dockerignore` isoliert (Docker-Mechanismus:
# `<Dockerfile>.dockerignore`, wirkt ausschließlich für Bauten, die genau
# diese Datei über `-f` referenzieren) — die Wurzel-`Dockerfile` und die
# Wurzel-`.dockerignore` bleiben dadurch byte-gleich unberührt (real gemessen,
# ADR-0098 Kontext).
#
# Basis-Digests sind aus der Wurzel-`Dockerfile` wiederverwendet, kein neuer
# Pin: Pin-Hebung fällt automatisch mit der der Wurzel-`Dockerfile` zusammen
# (`make image-stale`).

# --- build: kompiliert alle vier Go-Beispiele in einer gemeinsamen Stufe ---
FROM golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY gen/ gen/
COPY examples/ examples/
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
      -o /out/http-client ./examples/http-client && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
      -o /out/sse-client ./examples/sse-client && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
      -o /out/grpc-client ./examples/grpc-client && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
      -o /out/nats-client ./examples/nats-client

# --- runtime-sse: startbares Image des SSE-Clients ---
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime-sse
COPY --from=build /out/sse-client /sse-client
USER nonroot
ENTRYPOINT ["/sse-client"]

# --- runtime-grpc: startbares Image des gRPC-Clients ---
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime-grpc
COPY --from=build /out/grpc-client /grpc-client
USER nonroot
ENTRYPOINT ["/grpc-client"]

# --- runtime-nats: startbares Image des NATS-Clients ---
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime-nats
COPY --from=build /out/nats-client /nats-client
USER nonroot
ENTRYPOINT ["/nats-client"]

# --- runtime: startbares Image des HTTP-Clients. Bleibt die letzte Stufe der
# Datei, damit ein `docker build` ohne `--target` (mirror der ADR-0087-
# Konvention für C#/Kotlin) dieses Image liefert. ---
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime
COPY --from=build /out/http-client /http-client
USER nonroot
ENTRYPOINT ["/http-client"]
