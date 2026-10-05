# tools/harness/lib-sdk-tls-fixture.sh — gemeinsame Vorbereitung der TLS-Phasen
# für die drei SDK-Realserver-Runner (LH-FA-SST-013): eine Aufrufform für alle
# drei Runner (BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall). Wird per
# `source` eingebunden (Muster tools/harness/lib-sdk-rule-fixture.sh), kein
# eigenes Executable. Setzt COMPOSE, NETWORK, FEED_CONTAINER, TOOLCHAIN_IMAGE
# und GO_MODCACHE_VOLUME voraus (Konvention der drei Runner; das Makefile setzt
# die beiden letzten).
#
# Die Fixture fährt den Feed-Container über eine Compose-Override-Datei im
# Temp-Verzeichnis neu hoch (`compose.yaml` bleibt unverändert): ein
# Zertifikatspaar `good` (Namen `pg-change-feed` und `localhost`) für den
# Server und ein zweites Paar `other` (gleiche Namen, anderer Schlüssel) als
# fremder Vertrauensanker. Der Name `cdc-test-feed` (Container-Name, im
# Docker-Netz auflösbar) steht bewusst nicht im Zertifikat: die Adresse mit
# diesem Namen ist die Namensabweichung gegen denselben Server. Zertifikate und
# Schlüssel erzeugt `tools/harness/certgen` im Toolchain-Container in das
# Temp-Verzeichnis; nichts davon liegt im Arbeitsbaum oder in einem Image.

SDK_TLS_TMP=""
SDK_TLS_MOUNT=/etc/cdc/tls
SDK_TLS_TEST_MOUNT=/tls
SDK_TLS_HTTP_ADDR="https://pg-change-feed:8090"
SDK_TLS_GRPC_ADDR="pg-change-feed:9090"
SDK_TLS_MISMATCH_HOST=cdc-test-feed
# Vertrauensanker der Phasen im Test-Container. Beide Pfade sind übersteuerbar
# (Mutationsläufe: der fremde Anker als das Server-Zertifikat, der Anker als das
# fremde Zertifikat).
SDK_TLS_CA_FILE=${SDK_TLS_CA_FILE:-$SDK_TLS_TEST_MOUNT/good.pem}
SDK_TLS_FOREIGN_CA_FILE=${SDK_TLS_FOREIGN_CA_FILE:-$SDK_TLS_TEST_MOUNT/other.pem}

# sdk_tls_await_healthy <Aufrufer-Präfix> <Beschreibung> — wartet mit Frist auf
# den Compose-Health-Status `healthy` des Feed-Containers.
sdk_tls_await_healthy() {
  local prefix=$1 what=$2 health=""
  for _ in $(seq 1 60); do
    health=$(docker inspect --format '{{.State.Health.Status}}' "$FEED_CONTAINER" 2>/dev/null || echo fehlt)
    if [ "$health" = "healthy" ]; then
      return 0
    fi
    sleep 1
  done
  echo "$prefix: $what — Feed-Container meldet Compose-Health-Status ${health:-fehlt}, wollen healthy: $(docker logs "$FEED_CONTAINER" 2>&1 | tail -n 5)" >&2
  exit 1
}

# sdk_tls_fixture_setup <Aufrufer-Präfix> — erzeugt die beiden Zertifikatspaare,
# startet den Feed-Container mit dem Paar `good` neu, wartet auf `healthy` und
# verlangt, dass beide Server (HTTP und gRPC) im Log melden, über TLS zu
# bedienen.
sdk_tls_fixture_setup() {
  local prefix=$1 name out log
  SDK_TLS_TMP=$(mktemp -d)
  chmod 0755 "$SDK_TLS_TMP"

  for name in good other; do
    out=$(docker run --rm --network none --user "$(id -u):$(id -g)" \
      -v "$(pwd)":/src:ro \
      -v "$GO_MODCACHE_VOLUME":/go/pkg/mod \
      -v "$SDK_TLS_TMP":/out \
      -w /src \
      -e GOCACHE=/tmp/gocache -e HOME=/tmp \
      "$TOOLCHAIN_IMAGE" go run ./tools/harness/certgen /out "$name" pg-change-feed localhost 2>&1) || {
      echo "$prefix: certgen ($name) endete mit einem Fehler: $out" >&2
      exit 1
    }
    if [ "$(printf '%s\n' "$out" | grep -c '^CERTGEN ' || true)" != "1" ]; then
      echo "$prefix: certgen ($name) druckte nicht genau eine CERTGEN-Zeile: $out" >&2
      exit 1
    fi
  done
  # Der Container läuft als `nonroot`, der Besitzer der Dateien ist der
  # Host-Benutzer: Zertifikate und Schlüssel brauchen Leserechte für andere. Die
  # Schlüssel sind Wegwerf-Material dieses Temp-Verzeichnisses.
  chmod 0644 "$SDK_TLS_TMP"/*.pem
  if [ "$(git status --porcelain --untracked-files=all -- . | grep -c -E '\.pem$' || true)" != "0" ]; then
    echo "$prefix: Zertifikatsdateien im Arbeitsbaum" >&2
    exit 1
  fi

  {
    printf 'services:\n  pg-change-feed:\n    environment:\n'
    printf '      CDC_TLS_CERT_FILE: "%s/good.pem"\n' "$SDK_TLS_MOUNT"
    printf '      CDC_TLS_KEY_FILE: "%s/good-key.pem"\n' "$SDK_TLS_MOUNT"
    printf '    volumes:\n      - %s:%s:ro\n' "$SDK_TLS_TMP" "$SDK_TLS_MOUNT"
  } > "$SDK_TLS_TMP/override.yaml"
  chmod 0644 "$SDK_TLS_TMP/override.yaml"
  $COMPOSE -f "$SDK_TLS_TMP/override.yaml" up -d --force-recreate --no-deps pg-change-feed >/dev/null
  sdk_tls_await_healthy "$prefix" "Start mit TLS-Paar"

  log=$(docker logs "$FEED_CONTAINER" 2>&1)
  printf '%s\n' "$log" | grep 'http: Adapter gestartet' | grep '"tls":true' >/dev/null || {
    echo "$prefix: der HTTP-Server meldet nicht, über TLS zu bedienen" >&2
    exit 1
  }
  printf '%s\n' "$log" | grep 'grpc: Adapter gestartet' | grep '"tls":true' >/dev/null || {
    echo "$prefix: der gRPC-Server meldet nicht, über TLS zu bedienen" >&2
    exit 1
  }
}

# sdk_tls_fixture_restore <Aufrufer-Präfix> — stellt den Feed-Container ohne
# Override wieder her (Klartext) und wartet auf `healthy`.
sdk_tls_fixture_restore() {
  local prefix=$1
  $COMPOSE up -d --force-recreate --no-deps pg-change-feed >/dev/null
  sdk_tls_await_healthy "$prefix" "Wiederherstellung ohne TLS-Paar"
}
