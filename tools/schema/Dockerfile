# tools/schema/Dockerfile — Eingabe von `make schema-validate` und
# `make schema-rollout` ohne Bind-Mount des Arbeitsbaums (ADR-0142
# Festlegung 3): das neutrale Schema-YAML kommt per COPY in das Image der
# Stufe `rollout`, die Idempotenz-Wache als statisches Binary in das Image
# der Stufe `guard`. Gebaut wird mit `-f tools/schema/Dockerfile` und dem
# Kontext der Repo-Wurzel; die Allow-Liste des Kontexts steht in
# tools/schema/Dockerfile.dockerignore.
#
# Die Pins D_MIGRATE_IMAGE und TOOLCHAIN_IMAGE stehen allein im Makefile und
# kommen als --build-arg; dieses Dockerfile führt keinen eigenen Digest.
# Beide Build-Args tragen beide Bauten, weil ein FROM ohne Wert die Stufe
# auch dann verwirft, wenn sie nicht das Ziel ist.
#
# Aufrufer: tools/schema/rollout.sh. Ablauf und Erzeugnisse:
# harness/targets/schema-rollout.md.

ARG D_MIGRATE_IMAGE
ARG TOOLCHAIN_IMAGE

# --- rollout: d-migrate mit dem Schema-YAML unter /work/schema.yaml ----------
# Das Basis-Image setzt ENTRYPOINT d-migrate, USER dmigrate und WORKDIR /work;
# die Läufe schreiben Reports und das Rollback-Artefakt nach /work, der
# Aufrufer holt sie per `docker cp`.
FROM ${D_MIGRATE_IMAGE} AS rollout
ARG SCHEMA_SOURCE=tools/schema/schema.yaml
COPY ${SCHEMA_SOURCE} /work/schema.yaml

# --- guard-build: statisches Binary der Idempotenz-Wache ---------------------
FROM ${TOOLCHAIN_IMAGE} AS guard-build
WORKDIR /src
COPY go.mod go.sum ./
COPY tools/schema/rolloutguard/ tools/schema/rolloutguard/
RUN CGO_ENABLED=0 go build -o /out/rolloutguard ./tools/schema/rolloutguard

# --- guard: nur das Binary; der Report kommt über stdin ----------------------
FROM scratch AS guard
COPY --from=guard-build /out/rolloutguard /rolloutguard
ENTRYPOINT ["/rolloutguard"]
