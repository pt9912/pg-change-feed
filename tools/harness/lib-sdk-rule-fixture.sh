# tools/harness/lib-sdk-rule-fixture.sh — gemeinsame Vorbereitung der
# Regel-Phase für die drei SDK-Realserver-Runner
# (slice-sdk-regel-realserver-e2e, ADR-0112): legt eine eigene Tabelle an,
# aktiviert sie per `cdc.enable_table` und setzt eine aktive
# `rename_column`-Regel per `cdc.set_transformation` — eine Kopie der
# Aufrufform statt drei (BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall).
# Wird per `source` eingebunden (Muster tools/harness/lib-github-api.sh),
# kein eigenes Executable. Setzt PG_CONTAINER, PG_USER, PG_DB voraus (Konvention
# der drei Runner).

SDK_RULE_TABLE=feed_e2e_sdkrule
SDK_RULE_SOURCE_KEY=name
SDK_RULE_TARGET_KEY=display_name
SDK_RULE_NAME=sdk_rename

# sdk_rule_fixture_setup <Aufrufer-Präfix> — legt $SDK_RULE_TABLE an,
# aktiviert sie und setzt die rename_column-Regel ($SDK_RULE_SOURCE_KEY nach
# $SDK_RULE_TARGET_KEY); wartet je Antrag mit Frist auf `status = 'applied'`,
# `failed` beendet den Runner mit dem Fehlertext der Zeile.
sdk_rule_fixture_setup() {
  local prefix=$1 enable_request_id rule_request_id

  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 >/dev/null <<SQL
CREATE TABLE public.$SDK_RULE_TABLE (id int PRIMARY KEY, name text);
SQL

  enable_request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT cdc.enable_table('src-e2e', 'public', '$SDK_RULE_TABLE')")
  if [ -z "$enable_request_id" ]; then
    echo "$prefix: cdc.enable_table($SDK_RULE_TABLE) lieferte keine Antrags-ID" >&2
    exit 1
  fi
  sdk_rule_fixture_await_applied "$prefix" "$enable_request_id" "cdc.enable_table($SDK_RULE_TABLE)"

  rule_request_id=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
    "SELECT cdc.set_transformation('src-e2e', 'public', '$SDK_RULE_TABLE', '$SDK_RULE_NAME', '{\"kind\": \"rename_column\", \"column\": \"$SDK_RULE_SOURCE_KEY\", \"to\": \"$SDK_RULE_TARGET_KEY\"}')")
  if [ -z "$rule_request_id" ]; then
    echo "$prefix: cdc.set_transformation($SDK_RULE_TABLE) lieferte keine Antrags-ID" >&2
    exit 1
  fi
  sdk_rule_fixture_await_applied "$prefix" "$rule_request_id" "cdc.set_transformation($SDK_RULE_TABLE)"
}

# sdk_rule_fixture_await_applied <Aufrufer-Präfix> <Antrags-ID> <Beschreibung>
sdk_rule_fixture_await_applied() {
  local prefix=$1 request_id=$2 what=$3 status="" message
  for _ in $(seq 1 60); do
    status=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
      "SELECT status FROM cdc.administration_request WHERE administration_request_id = '$request_id'")
    if [ "$status" = "applied" ]; then
      return 0
    fi
    if [ "$status" = "failed" ]; then
      message=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc \
        "SELECT error_message FROM cdc.administration_request WHERE administration_request_id = '$request_id'")
      echo "$prefix: $what-Antrag $request_id scheiterte: $message" >&2
      exit 1
    fi
    sleep 0.5
  done
  echo "$prefix: $what-Antrag $request_id wurde nicht innerhalb der Zeitspanne verarbeitet (status=${status:-leer})" >&2
  exit 1
}
