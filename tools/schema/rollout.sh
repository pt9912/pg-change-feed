#!/usr/bin/env bash
# rollout.sh — Rezeptur von `make schema-validate` und `make schema-rollout`
# ohne Bind-Mount des Arbeitsbaums (ADR-0142). Ablauf, Reihenfolge, Exit-Codes
# und Grenzen: harness/targets/schema-rollout.md.
#
# Die Eingabe reist per COPY in zwei Images (tools/schema/Dockerfile): das
# Schema-YAML in `rollout` (d-migrate), die Wache als statisches Binary in
# `guard`. Die Nacharbeit-Dateien lesen die psql-Container über stdin. Die
# Erzeugnisse (plan.yaml, down.sql, rollout-precheck.yaml) verlassen den
# d-migrate-Container als tar-Stream (`docker cp … - | tar -x`) in
# SCHEMA_ARTEFACT_DIR außerhalb des versionierten Baums; der Arbeitsbaum
# bleibt unberührt.
#
# `set -o pipefail` macht die Export-Pipe sicher: ein `docker cp`-Fehlschlag
# verschwindet sonst hinter einem erfolgreichen, aber leeren `tar -x`. Das
# Skript läuft unter bash, das make-Rezept selbst unter /bin/sh.
#
# Kopplung: die Reihenfolge der vier Nacharbeit-Dateien unten und die
# Bekannt-Liste `knownForeignObjects` in tools/schema/rolloutguard/guard.go
# werden gemeinsam geändert (neue Datei: Aufrufzeile und Listeneintrag im
# selben Commit).
#
# Aufruf: rollout.sh validate | rollout
# Umgebung (vom Makefile gesetzt): D_MIGRATE_IMAGE, TOOLCHAIN_IMAGE,
# PG_TEST_IMAGE, SCHEMA_SOURCE, SCHEMA_TARGET, SCHEMA_ROLLOUT_NETWORK,
# SCHEMA_ARTEFACT_DIR.
set -uo pipefail

mode=${1:-}
case $mode in
  validate | rollout) ;;
  *)
    echo "Aufruf: rollout.sh validate | rollout" >&2
    exit 2
    ;;
esac

: "${D_MIGRATE_IMAGE:?}" "${TOOLCHAIN_IMAGE:?}" "${PG_TEST_IMAGE:?}"
: "${SCHEMA_SOURCE:?}" "${SCHEMA_TARGET:?}" "${SCHEMA_ROLLOUT_NETWORK:?}"
: "${SCHEMA_ARTEFACT_DIR:?}"

cd "$(dirname "$0")/../.." || exit 2

ROLLOUT_IMG=pg-change-feed-schema:rollout
GUARD_IMG=pg-change-feed-schema:guard
DSN=${SCHEMA_TARGET#db:}
NET=$SCHEMA_ROLLOUT_NETWORK
DIR=$SCHEMA_ARTEFACT_DIR
CONTAINERS=()

cleanup() {
  local c
  for c in ${CONTAINERS[@]+"${CONTAINERS[@]}"}; do
    docker rm -f "$c" >/dev/null 2>&1 || true
  done
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

if [ ! -f "$SCHEMA_SOURCE" ]; then
  echo "FEHLER: $SCHEMA_SOURCE fehlt — das neutrale Schema-YAML ist die Erstlieferung des d-migrate-Einbaus (ADR-0043); Überführungsquelle ist internal/adapters/driven/postgresstorage/schema.sql" >&2
  exit 2
fi

build_stage() {
  docker build -q -f tools/schema/Dockerfile --target "$1" -t "$2" \
    --build-arg "D_MIGRATE_IMAGE=$D_MIGRATE_IMAGE" \
    --build-arg "TOOLCHAIN_IMAGE=$TOOLCHAIN_IMAGE" \
    --build-arg "SCHEMA_SOURCE=$SCHEMA_SOURCE" . >/dev/null
}

# run_dmigrate <name> <netz> <d-migrate-Argumente …>: ein d-migrate-Lauf als
# create/start; der Exit-Code ist der des Containers, Reports bleiben bis zum
# Export (und zum trap) im gestoppten Container.
run_dmigrate() {
  local name=$1 net=$2
  shift 2
  name="pg-change-feed-schema-$$-$name"
  docker create --name "$name" --network "$net" "$ROLLOUT_IMG" "$@" >/dev/null || return 125
  CONTAINERS+=("$name")
  docker start -a "$name"
}

# export_file <name> <Pfad im Container>: eine Datei als tar-Stream nach DIR.
export_file() {
  docker cp "pg-change-feed-schema-$$-$1:$2" - | tar -x --no-same-owner -C "$DIR"
}

build_stage rollout "$ROLLOUT_IMG" || exit 2

if [ "$mode" = validate ]; then
  docker run --rm --network none "$ROLLOUT_IMG" schema validate --source /work/schema.yaml
  exit $?
fi

build_stage guard "$GUARD_IMG" || exit 2
mkdir -p "$DIR" || exit 2
rm -f "$DIR/plan.yaml" "$DIR/down.sql" "$DIR/rollout-precheck.yaml"

# Precheck: sein Exit-Code wird gelesen, nicht durchgereicht.
plan_exit=0
run_dmigrate precheck "$NET" schema migrate --source /work/schema.yaml \
  --target "$SCHEMA_TARGET" --plan-only --report /work/rollout-precheck.yaml || plan_exit=$?
if [ "$plan_exit" = 0 ] || [ "$plan_exit" = 8 ]; then
  export_file precheck /work/rollout-precheck.yaml || exit 2
fi

execute_flags=()
drop_views=""
if [ "$plan_exit" = 8 ]; then
  # Die Wache liest den Precheck-Report über stdin; ihr Exit-Code wird
  # gelesen: ohne Erweiterung läuft alles Weitere ohne sie.
  if guard_out=$(docker run --rm -i --network none "$GUARD_IMG" /dev/stdin <"$DIR/rollout-precheck.yaml"); then
    if printf '%s\n' "$guard_out" | grep -qx 'allow-destructive'; then
      echo "schema-rollout: bekannte Fremdobjekt-Blocker (ADR-0043) - --execute laeuft mit --allow-destructive"
      execute_flags=(--allow-destructive)
    fi
    drop_views=$(printf '%s\n' "$guard_out" | sed -n 's/^drop-view //p')
  fi
fi

for v in $drop_views; do
  echo "schema-rollout: Vorlauf (ADR-0114) - View-Signatur-Aenderung, DROP VIEW cdc.$v"
  docker run --rm --network "$NET" "$PG_TEST_IMAGE" psql "$DSN" -v ON_ERROR_STOP=1 -c "DROP VIEW cdc.$v" || exit 1
done

exec_exit=0
run_dmigrate execute "$NET" schema migrate --source /work/schema.yaml \
  --target "$SCHEMA_TARGET" --execute ${execute_flags[@]+"${execute_flags[@]}"} \
  --report /work/plan.yaml --generate-rollback --rollback-output /work/down.sql || exec_exit=$?
if [ "$exec_exit" -ne 0 ]; then
  export_file execute /work/plan.yaml 2>/dev/null || true
  exit "$exec_exit"
fi
export_file execute /work/plan.yaml || exit 2
export_file execute /work/down.sql || exit 2

# Nacharbeit: `-i` hält stdin offen; ohne ihn liest `psql -f -` sofort EOF
# und endet mit Exit 0, ohne etwas auszuführen.
for f in nacharbeit-roles nacharbeit-observability nacharbeit-heartbeat nacharbeit-administration; do
  docker run --rm -i --network "$NET" "$PG_TEST_IMAGE" psql "$DSN" -v ON_ERROR_STOP=1 -f - <"tools/schema/$f.sql" || exit $?
done

echo "schema-rollout: Erzeugnisse (plan.yaml, down.sql, rollout-precheck.yaml) in $DIR"
