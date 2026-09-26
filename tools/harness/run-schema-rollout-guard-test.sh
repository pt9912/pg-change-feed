#!/usr/bin/env bash
# run-schema-rollout-guard-test — realer Beleg der zentralen
# `make schema-rollout`-Wache (ADR-0043, ADR-0114, tools/schema/rolloutguard,
# BEO-PGC/schema-rollout-fremdobjekte): sechs Läufe gegen eigenständige
# Ziel-Datenbanken derselben Wegwerf-Instanz.
#
#   1. Frischer Rollout (kein Blocker) — muss durchlaufen.
#   2. Zweiter Lauf gegen dasselbe, jetzt vollständig migrierte Ziel — trägt
#      ausschließlich die neun bekannten Fremdobjekt-Blocker — muss erneut
#      Exit 0 liefern UND den --allow-destructive-Pfad des Guards nehmen
#      (stdout trägt die Meldung), nicht nur zufällig Exit 0 aus anderem
#      Grund; ein Vorlauf (ADR-0114) findet nicht statt.
#   3. Eine echte, gleichzeitig anstehende, NICHT-destruktive Schema-Änderung
#      neben den neun bekannten Blockern: eine von schema.yaml weiterhin
#      deklarierte, nullable Spalte (administration_request.error_message —
#      trägt keine View-/Funktions-Abhängigkeit, anders als
#      process_heartbeat.error_class, das die Sicht cdc.heartbeat trägt) wird
#      per `ALTER TABLE … DROP COLUMN` außerhalb von d-migrate entfernt — der
#      nächste Rollout-Lauf muss sie real zurückbringen, ohne Vorlauf. Das ist
#      der Regressionsschutz für den konkret gefundenen Fehler: ein reines
#      Überspringen von `--execute` bei bekannten Blockern ließ eine solche
#      echte Änderung verlustig gehen.
#   4. Die View cdc.changes trägt eine abweichende Signatur (ADR-0114,
#      Klasse „View-Signatur") — der Rollout muss Exit 0 liefern, den Vorlauf
#      melden, die Soll-Signatur herstellen, cdc_reader das SELECT-Recht
#      zurückgeben und die bestehende Zeile über die View lesbar lassen; ein
#      Folgelauf braucht keinen Vorlauf. Hängt ein fremdes Objekt an der View
#      (eine View im Schema `public`, das d-migrate nicht liest), scheitert
#      der Rollout laut, das Objekt bleibt bestehen (`DROP VIEW` ohne
#      CASCADE), und der Wiederholungslauf ohne das Objekt heilt.
#   5. Alt-Tag-Lauf (ADR-0114 Entscheidung 7): das Schema des jüngsten
#      `v*`-Tags (`git archive` in ein Verzeichnis unter `${TMPDIR:-/tmp}`,
#      nie im Repo-Baum) wird mit dem Makefile dieses Tags ausgerollt, eine
#      Datenzeile geschrieben, danach der Arbeitsbaum zweimal ausgerollt —
#      Exit 0 zweimal, die Zeile über `cdc.changes` unverändert lesbar,
#      Soll-Signatur der View, und der Rollen-Schnitt: die Rechte der drei
#      Rollen auf `cdc.administration_request`, `cdc.backfill_run` und
#      `cdc.backfill_status`, `EXECUTE` auf `cdc.backfill_table`,
#      `cdc.set_transformation` und `cdc.remove_transformation` allein für
#      `cdc_admin` (nicht PUBLIC). Der Stand des Tags trägt weder die Spalten
#      `rule_name`/`rule_spec` noch die zwei Transformations-Funktionen noch
#      die zwei Transformations-Antragsarten (Vorbedingung); nach dem Upgrade
#      trägt die Tabelle die zwei nullable Spalten (`text`, `jsonb`), eine
#      vor dem Upgrade geschriebene Antragszeile trägt dort NULL, die sieben
#      Werte von `chk_administration_request_kind` stehen, `cdc_admin` schreibt
#      über `cdc.set_transformation` einen `pending`-Antrag mit der Regelform
#      in der `jsonb`-Spalte, und `cdc_reader` scheitert mit „permission
#      denied for function“. Tag, Exit-Codes und Zählung stehen in der Ausgabe. Rot
#      färbende Eingabeseiten-Mutationen in
#      tools/schema/nacharbeit-administration.sql:
#      `cdc.set_transformation(text, text, text, text, json)` aus der
#      GRANT-Zeile, die REVOKE-Zeile oder `'set_transformation'` aus dem CHECK;
#      in tools/schema/nacharbeit-roles.sql: `cdc.backfill_status` aus dem
#      GRANT für `cdc_reader`; in tools/schema/schema.yaml: die Spalte
#      `rule_spec`.
#   6. Unbekannte Blocker — der Rollout muss abbrechen (d-migrate-Exit 8,
#      make meldet „Error 8" bzw. lokalisiert „Fehler 8"). Zwei Fälle:
#      a) eine nicht deklarierte Funktion `cdc.zz_rolloutguard_unbekannt()`:
#         dieselbe Klasse wie die neun bekannten Fremdobjekte (Blocker
#         DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION), aber nicht auf der
#         Bekannt-Liste. Ließe die Wache den Blocker durch, riefe das Target
#         `--execute` mit `--allow-destructive` auf und d-migrate löschte die
#         Funktion — der Beleg ist, dass sie nach dem Lauf besteht. Nur dieser
#         Fall bindet die Bekannt-Liste (`knownForeignObjects`) des Guards
#         end-to-end; die Liste selbst trägt der Unit-Test in
#         tools/schema/rolloutguard.
#      b) eine nicht deklarierte Spalte per `ALTER TABLE … ADD COLUMN`:
#         d-migrate meldet einen DropColumn-Blocker und bricht auch unter
#         `--allow-destructive` selbst ab — der Beleg für den Abbruch mit
#         Exit 8 gegen einen real gemeldeten Blocker, nicht für die
#         Bekannt-Liste.
#
# Am Ende (auch beim Fehlschlag) stellt der Lauf die vom Rollout
# überschriebenen Erzeugnisse tools/schema/plan.yaml und tools/schema/down.sql
# wieder her; der Arbeitsbaum bleibt unverändert.
#
# Eigenständiges Docker-Netz/-Container, unabhängig von der
# Wurzel-`compose.yaml` und `examples/compose.yaml` — Daten leben
# ausschließlich im Container, in jedem Ausgang abgeräumt. Kein Gate:
# braucht DB-Zugang, `make gates` bleibt netzlos.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

PG_TEST_IMAGE=${PG_TEST_IMAGE:-postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8}
NETWORK=cdc-schema-rollout-guard-test
CONTAINER=cdc-schema-rollout-guard-test-pg
DB=cdc
ALT_DB=cdc_alttag
USER=postgres
PASSWORD=postgres
TARGET="db:postgres://$USER:$PASSWORD@$CONTAINER:5432/$DB?sslmode=disable"
ALT_TARGET="db:postgres://$USER:$PASSWORD@$CONTAINER:5432/$ALT_DB?sslmode=disable"
ALT_DIR=""
# `make schema-rollout` ohne `-C` überschreibt die committeten Erzeugnisse
# tools/schema/plan.yaml und down.sql; cleanup stellt sie aus dieser Sicherung wieder her.
ARTEFACT_BACKUP=$(mktemp -d "${TMPDIR:-/tmp}/schema-rollout-artefakte.XXXXXX")
cp tools/schema/plan.yaml tools/schema/down.sql "$ARTEFACT_BACKUP"/

docker network inspect "$NETWORK" >/dev/null 2>&1 || docker network create "$NETWORK" >/dev/null

# `-v` entfernt die anonymen Volumes des Containers (das Postgres-Image
# deklariert ein VOLUME).
cleanup() {
  docker rm -fv "$CONTAINER" >/dev/null 2>&1 || true
  docker network rm "$NETWORK" >/dev/null 2>&1 || true
  if [ -n "$ALT_DIR" ]; then
    rm -rf "$ALT_DIR"
  fi
  cp "$ARTEFACT_BACKUP"/plan.yaml "$ARTEFACT_BACKUP"/down.sql tools/schema/
  rm -rf "$ARTEFACT_BACKUP"
}
trap cleanup EXIT

fail() {
  echo "run-schema-rollout-guard-test: FEHLER — $*" >&2
  exit 1
}

# psql_q <db> <sql>: eine Abfrage, ungerahmte Ausgabe.
psql_q() {
  docker exec "$CONTAINER" psql -U "$USER" -d "$1" -v ON_ERROR_STOP=1 -tAc "$2"
}

# psql_as <db> <rolle> <sql>: eine Abfrage unter `SET ROLE <rolle>` (der
# Superuser der Wegwerf-Instanz darf jede Gruppenrolle annehmen), ungerahmte
# Ausgabe ohne die Befehlsmeldung des `SET`.
psql_as() {
  docker exec "$CONTAINER" psql -q -U "$USER" -d "$1" -v ON_ERROR_STOP=1 -tA -c "SET ROLE $2" -c "$3"
}

# view_signature <db>: Spaltenname:Typ der View cdc.changes in Reihenfolge.
view_signature() {
  psql_q "$1" "SELECT string_agg(column_name || ':' || data_type, ',' ORDER BY ordinal_position) FROM information_schema.columns WHERE table_schema='cdc' AND table_name='changes'"
}

# seed_row <db> <praefix>: eine vollständige Kette Quelle → Change; die
# Zeile trägt keinen origin-Wert (auch ein Schema-Stand ohne die Spalte
# nimmt sie an).
seed_row() {
  docker exec "$CONTAINER" psql -U "$USER" -d "$1" -v ON_ERROR_STOP=1 \
    -c "INSERT INTO cdc.source (source_id, name) VALUES ('$2-src', '$2')" \
    -c "INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ('$2-st', '$2-src', 'public', 't')" \
    -c "INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ('$2-sv', '$2-st', 1)" \
    -c "INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ('$2-tx', '$2-src', 1)" \
    -c "INSERT INTO cdc.change (change_id, transaction_id, source_table_id, sequence, operation, new_data, schema_version) VALUES ('$2-ch', '$2-tx', '$2-st', 1, 'INSERT', '{\"a\": 1}', '$2-sv')" >/dev/null
}

# run_rollout <make-verzeichnis> <ziel-dsn>: ein `make schema-rollout`;
# Ausgabe gedruckt und in RUN_OUT, Exit-Code in RUN_EXIT (kein Abbruch).
run_rollout() {
  set +e
  RUN_OUT=$(make -C "$1" schema-rollout SCHEMA_TARGET="$2" SCHEMA_ROLLOUT_NETWORK="$NETWORK" 2>&1)
  RUN_EXIT=$?
  set -e
  echo "$RUN_OUT"
}

docker rm -fv "$CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$CONTAINER" \
  --network "$NETWORK" \
  -e POSTGRES_DB="$DB" -e POSTGRES_USER="$USER" -e POSTGRES_PASSWORD="$PASSWORD" \
  "$PG_TEST_IMAGE" -c wal_level=logical >/dev/null

ready=0
for _ in $(seq 1 60); do
  if docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -tAc "SELECT 1" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
if [ "$ready" -ne 1 ]; then
  echo "run-schema-rollout-guard-test: $CONTAINER nahm innerhalb 60s keine Verbindung an" >&2
  exit 1
fi

docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "CREATE SCHEMA IF NOT EXISTS cdc" \
  -c "ALTER ROLE $USER IN DATABASE $DB SET search_path = cdc"

echo "run-schema-rollout-guard-test: Lauf 1/6 (frischer Rollout, muss durchlaufen)"
make schema-rollout SCHEMA_TARGET="$TARGET" SCHEMA_ROLLOUT_NETWORK="$NETWORK"

echo "run-schema-rollout-guard-test: Lauf 2/6 (Idempotenz — muss den --allow-destructive-Pfad nehmen, ohne Vorlauf)"
out=$(make schema-rollout SCHEMA_TARGET="$TARGET" SCHEMA_ROLLOUT_NETWORK="$NETWORK" 2>&1)
echo "$out"
if ! grep -q -- "--allow-destructive" <<<"$out"; then
  echo "run-schema-rollout-guard-test: FEHLER — Lauf 2 hat den --allow-destructive-Pfad nicht genommen (Meldung fehlt in der Ausgabe)" >&2
  exit 1
fi
if grep -q "Vorlauf" <<<"$out"; then
  fail "Lauf 2 meldet einen Vorlauf, obwohl keine View ihre Signatur ändert (ADR-0114 Entscheidung 4)"
fi

echo "run-schema-rollout-guard-test: Lauf 3/6 (echte anstehende Änderung neben den neun bekannten Blockern — muss real zurückkommen, ohne Vorlauf)"
docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "ALTER TABLE cdc.administration_request DROP COLUMN error_message"
still_missing=$(docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -tAc \
  "SELECT count(*) FROM information_schema.columns WHERE table_schema='cdc' AND table_name='administration_request' AND column_name='error_message'")
if [ "$still_missing" != "0" ]; then
  echo "run-schema-rollout-guard-test: FEHLER — Vorbedingung fehlgeschlagen, error_message ist nach DROP COLUMN noch da" >&2
  exit 1
fi
out=$(make schema-rollout SCHEMA_TARGET="$TARGET" SCHEMA_ROLLOUT_NETWORK="$NETWORK" 2>&1)
echo "$out"
if grep -q "Vorlauf" <<<"$out"; then
  fail "Lauf 3 meldet einen Vorlauf, obwohl die anstehende Änderung additiv ist (ADR-0114 Entscheidung 5)"
fi
restored=$(docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -tAc \
  "SELECT count(*) FROM information_schema.columns WHERE table_schema='cdc' AND table_name='administration_request' AND column_name='error_message'")
if [ "$restored" != "1" ]; then
  echo "run-schema-rollout-guard-test: FEHLER — Lauf 3 hat die echte anstehende Spalten-Wiederherstellung nicht angewendet (Regression des behobenen Fehlers)" >&2
  exit 1
fi

echo "run-schema-rollout-guard-test: Lauf 4/6 (View-Signatur-Änderung, ADR-0114 — Vorlauf, Soll-Signatur, Recht, Zeile lesbar)"
sig_ref=$(view_signature "$DB")
[ -n "$sig_ref" ] || fail "Lauf 4: Soll-Signatur der View cdc.changes nicht lesbar"
seed_row "$DB" lauf4
docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "DROP VIEW cdc.changes" \
  -c "CREATE VIEW cdc.changes AS SELECT change_id, transaction_id FROM cdc.change"
[ "$(view_signature "$DB")" != "$sig_ref" ] || fail "Lauf 4: Vorbedingung fehlgeschlagen, die View trägt schon die Soll-Signatur"
[ "$(psql_q "$DB" "SELECT has_table_privilege('cdc_reader', 'cdc.changes', 'SELECT')")" = "f" ] \
  || fail "Lauf 4: Vorbedingung fehlgeschlagen, cdc_reader hat das SELECT-Recht schon vor dem Rollout"
run_rollout . "$TARGET"
[ "$RUN_EXIT" -eq 0 ] || fail "Lauf 4: Rollout gegen die abweichende View-Signatur endete mit Exit $RUN_EXIT statt 0"
grep -q "Vorlauf (ADR-0114) - View-Signatur-Aenderung, DROP VIEW cdc.changes" <<<"$RUN_OUT" \
  || fail "Lauf 4: die Vorlauf-Meldung für cdc.changes fehlt in der Ausgabe"
[ "$(view_signature "$DB")" = "$sig_ref" ] || fail "Lauf 4: die View trägt nach dem Rollout nicht die Soll-Signatur"
[ "$(psql_q "$DB" "SELECT has_table_privilege('cdc_reader', 'cdc.changes', 'SELECT')")" = "t" ] \
  || fail "Lauf 4: cdc_reader hat nach dem Rollout kein SELECT-Recht auf cdc.changes"
[ "$(psql_q "$DB" "SELECT change_id || '|' || origin FROM cdc.changes WHERE change_id = 'lauf4-ch'")" = "lauf4-ch|wal" ] \
  || fail "Lauf 4: die bestehende Zeile ist über cdc.changes nach dem Rollout nicht (als wal) lesbar"
run_rollout . "$TARGET"
[ "$RUN_EXIT" -eq 0 ] || fail "Lauf 4: der Folgelauf endete mit Exit $RUN_EXIT statt 0"
if grep -q "Vorlauf" <<<"$RUN_OUT"; then
  fail "Lauf 4: der Folgelauf meldet einen Vorlauf, obwohl die View die Soll-Signatur trägt (ADR-0114 Entscheidung 4)"
fi

echo "run-schema-rollout-guard-test: Lauf 4/6 (abhängiges Objekt — der Vorlauf scheitert laut, das Objekt bleibt bestehen, der Wiederholungslauf heilt)"
view_def=$(psql_q "$DB" "SELECT pg_get_viewdef('cdc.changes'::regclass)")
docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "CREATE VIEW public.zz_lauf4_abhaengig AS SELECT change_id FROM cdc.changes" \
  -c "CREATE OR REPLACE VIEW cdc.changes AS SELECT c.*, 1 AS zz_lauf4_zusatz FROM (${view_def%;}) c"
[ "$(view_signature "$DB")" != "$sig_ref" ] || fail "Lauf 4: Vorbedingung fehlgeschlagen, die View trägt trotz Zusatzspalte die Soll-Signatur"
run_rollout . "$TARGET"
[ "$RUN_EXIT" -ne 0 ] || fail "Lauf 4: der Rollout lief durch (Exit 0), obwohl ein fremdes Objekt an der View hängt — der Vorlauf darf nicht kaskadieren"
grep -q "DROP VIEW cdc.changes" <<<"$RUN_OUT" || fail "Lauf 4: die Vorlauf-Meldung fehlt in der Ausgabe des scheiternden Laufs"
[ "$(psql_q "$DB" "SELECT to_regclass('public.zz_lauf4_abhaengig') IS NOT NULL")" = "t" ] \
  || fail "Lauf 4: das abhängige Objekt wurde mitgelöscht"
docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 -c "DROP VIEW public.zz_lauf4_abhaengig"
run_rollout . "$TARGET"
[ "$RUN_EXIT" -eq 0 ] || fail "Lauf 4: der Wiederholungslauf nach dem gescheiterten Vorlauf endete mit Exit $RUN_EXIT statt 0"
[ "$(view_signature "$DB")" = "$sig_ref" ] || fail "Lauf 4: die View trägt nach dem Wiederholungslauf nicht die Soll-Signatur"

echo "run-schema-rollout-guard-test: Lauf 5/6 (Alt-Tag-Lauf — Schema des jüngsten v*-Tags, danach der Arbeitsbaum)"
ALT_TAG=$(git tag --list 'v[0-9]*' --sort=-v:refname | head -n 1)
[ -n "$ALT_TAG" ] || fail "Lauf 5: kein v*-Tag im Repository"
ALT_DIR=$(mktemp -d "${TMPDIR:-/tmp}/schema-rollout-alt-tag.XXXXXX")
git archive "$ALT_TAG" | tar -x -C "$ALT_DIR"
docker exec "$CONTAINER" psql -U "$USER" -d postgres -v ON_ERROR_STOP=1 -c "CREATE DATABASE $ALT_DB"
docker exec "$CONTAINER" psql -U "$USER" -d "$ALT_DB" -v ON_ERROR_STOP=1 \
  -c "CREATE SCHEMA IF NOT EXISTS cdc" \
  -c "ALTER ROLE $USER IN DATABASE $ALT_DB SET search_path = cdc"
run_rollout "$ALT_DIR" "$ALT_TARGET"
alt_exit=$RUN_EXIT
[ "$alt_exit" -eq 0 ] || fail "Lauf 5: der Rollout des Schema-Stands von $ALT_TAG endete mit Exit $alt_exit statt 0"
seed_row "$ALT_DB" alttag
alt_rows_before=$(psql_q "$ALT_DB" "SELECT count(*) FROM cdc.changes")
[ "$alt_rows_before" = "1" ] || fail "Lauf 5: Vorbedingung fehlgeschlagen, cdc.changes trägt $alt_rows_before statt 1 Zeile im Stand von $ALT_TAG"
# alt_privilege <rolle> <objekt> <recht>: Tabellen-/View-Recht der Rolle im Alt-Tag-Ziel.
alt_privilege() { psql_q "$ALT_DB" "SELECT has_table_privilege('$1', '$2', '$3')"; }
ALT_FUNCTION="cdc.backfill_table(text, text, text)"
ALT_SET_FUNCTION="cdc.set_transformation(text, text, text, text, json)"
ALT_REMOVE_FUNCTION="cdc.remove_transformation(text, text, text, text)"
for alt_new_function in "$ALT_SET_FUNCTION" "$ALT_REMOVE_FUNCTION"; do
  [ "$(psql_q "$ALT_DB" "SELECT to_regprocedure('$alt_new_function') IS NULL")" = "t" ] \
    || fail "Lauf 5: Vorbedingung fehlgeschlagen, $alt_new_function besteht im Stand von $ALT_TAG schon"
done
[ "$(psql_q "$ALT_DB" "SELECT count(*) FROM information_schema.columns WHERE table_schema='cdc' AND table_name='administration_request' AND column_name IN ('rule_name', 'rule_spec')")" = "0" ] \
  || fail "Lauf 5: Vorbedingung fehlgeschlagen, cdc.administration_request trägt im Stand von $ALT_TAG schon rule_name oder rule_spec"
alt_kinds_sql="SELECT string_agg(k, ',' ORDER BY k) FROM (SELECT unnest(regexp_matches(pg_get_constraintdef(oid), '''([a-z_]+)''', 'g')) AS k FROM pg_constraint WHERE conname = 'chk_administration_request_kind' AND conrelid = 'cdc.administration_request'::regclass) s"
alt_kinds_before=$(psql_q "$ALT_DB" "$alt_kinds_sql")
case "$alt_kinds_before" in
  *transformation*) fail "Lauf 5: Vorbedingung fehlgeschlagen, chk_administration_request_kind trägt im Stand von $ALT_TAG schon eine Transformations-Antragsart ($alt_kinds_before)" ;;
esac
# Eine Antragszeile im Stand des Tags: nach dem Upgrade trägt sie in den zwei
# neuen Spalten NULL und bleibt lesbar.
docker exec "$CONTAINER" psql -U "$USER" -d "$ALT_DB" -v ON_ERROR_STOP=1 \
  -c "INSERT INTO cdc.administration_request (administration_request_id, source_id, schema_name, table_name, column_name, request_kind, status) VALUES ('alttag-req', 'alttag-src', 'public', 't', 'a', 'exclude_column', 'applied')" >/dev/null
run_rollout . "$ALT_TARGET"
work_exit_1=$RUN_EXIT
work_out_1=$RUN_OUT
[ "$work_exit_1" -eq 0 ] || fail "Lauf 5: der erste Rollout des Arbeitsbaums über $ALT_TAG endete mit Exit $work_exit_1 statt 0"
run_rollout . "$ALT_TARGET"
work_exit_2=$RUN_EXIT
[ "$work_exit_2" -eq 0 ] || fail "Lauf 5: der zweite Rollout des Arbeitsbaums endete mit Exit $work_exit_2 statt 0"
[ "$(psql_q "$ALT_DB" "SELECT change_id FROM cdc.changes")" = "alttag-ch" ] \
  || fail "Lauf 5: die Zeile aus dem Stand von $ALT_TAG ist über cdc.changes nach dem Upgrade nicht unverändert lesbar"
[ "$(view_signature "$ALT_DB")" = "$sig_ref" ] || fail "Lauf 5: die View trägt nach dem Upgrade nicht die Soll-Signatur"
alt_privilege_count=0
for expected in \
  "cdc_admin cdc.administration_request SELECT t" "cdc_admin cdc.administration_request UPDATE t" \
  "cdc_admin cdc.administration_request INSERT f" "cdc_admin cdc.administration_request DELETE f" \
  "cdc_capture cdc.administration_request SELECT f" "cdc_capture cdc.administration_request UPDATE f" \
  "cdc_reader cdc.administration_request SELECT f" "cdc_reader cdc.administration_request UPDATE f" \
  "cdc_admin cdc.backfill_run SELECT t" "cdc_admin cdc.backfill_run INSERT t" "cdc_admin cdc.backfill_run UPDATE f" \
  "cdc_capture cdc.backfill_run SELECT t" "cdc_capture cdc.backfill_run UPDATE t" "cdc_capture cdc.backfill_run INSERT f" \
  "cdc_reader cdc.backfill_run SELECT f" "cdc_reader cdc.backfill_run UPDATE f" \
  "cdc_reader cdc.backfill_status SELECT t" \
  "cdc_admin cdc.backfill_status SELECT f" "cdc_capture cdc.backfill_status SELECT f"; do
  read -r alt_role alt_object alt_right alt_want <<<"$expected"
  [ "$(alt_privilege "$alt_role" "$alt_object" "$alt_right")" = "$alt_want" ] \
    || fail "Lauf 5: $alt_role trägt nach dem Upgrade über $ALT_TAG auf $alt_object das Recht $alt_right nicht als $alt_want"
  alt_privilege_count=$((alt_privilege_count + 1))
done
alt_function_count=0
for alt_function in "$ALT_FUNCTION" "$ALT_SET_FUNCTION" "$ALT_REMOVE_FUNCTION"; do
  for expected in "cdc_admin t" "cdc_capture f" "cdc_reader f"; do
    read -r alt_role alt_want <<<"$expected"
    [ "$(psql_q "$ALT_DB" "SELECT has_function_privilege('$alt_role', '$alt_function', 'EXECUTE')")" = "$alt_want" ] \
      || fail "Lauf 5: $alt_role trägt nach dem Upgrade über $ALT_TAG auf $alt_function das Recht EXECUTE nicht als $alt_want"
  done
  [ "$(psql_q "$ALT_DB" "SELECT coalesce(bool_or(a.grantee = 0), false) FROM pg_proc p, LATERAL aclexplode(coalesce(p.proacl, acldefault('f', p.proowner))) a WHERE p.oid = '$alt_function'::regprocedure")" = "f" ] \
    || fail "Lauf 5: $alt_function trägt nach dem Upgrade über $ALT_TAG ein EXECUTE-Recht für PUBLIC"
  alt_function_count=$((alt_function_count + 1))
done
alt_kinds_after=$(psql_q "$ALT_DB" "$alt_kinds_sql")
[ "$alt_kinds_after" = "backfill,disable,enable,exclude_column,include_column,remove_transformation,set_transformation" ] \
  || fail "Lauf 5: chk_administration_request_kind trägt nach dem Upgrade über $ALT_TAG die Menge '$alt_kinds_after' statt der sieben Antragsarten"
alt_rule_columns=$(psql_q "$ALT_DB" "SELECT string_agg(column_name || ':' || data_type || ':' || is_nullable, ',' ORDER BY column_name) FROM information_schema.columns WHERE table_schema='cdc' AND table_name='administration_request' AND column_name IN ('rule_name', 'rule_spec')")
[ "$alt_rule_columns" = "rule_name:text:YES,rule_spec:jsonb:YES" ] \
  || fail "Lauf 5: cdc.administration_request trägt nach dem Upgrade über $ALT_TAG die Spalten '$alt_rule_columns' statt rule_name:text:YES,rule_spec:jsonb:YES"
[ "$(psql_q "$ALT_DB" "SELECT COALESCE(rule_name, '<NULL>') || '|' || COALESCE(rule_spec::text, '<NULL>') || '|' || status FROM cdc.administration_request WHERE administration_request_id = 'alttag-req'")" = "<NULL>|<NULL>|applied" ] \
  || fail "Lauf 5: die Antragszeile aus dem Stand von $ALT_TAG trägt nach dem Upgrade nicht NULL in rule_name und rule_spec"
# Die Funktion schreibt den Antrag unter cdc_admin; ein Login ohne diese
# Mitgliedschaft scheitert an EXECUTE.
alt_set_id=$(psql_as "$ALT_DB" cdc_admin "SELECT cdc.set_transformation('alttag-src', 'public', 't', 'umbenennung', '{\"kind\": \"rename_column\"}'::json)")
[ "$(psql_q "$ALT_DB" "SELECT request_kind || '|' || status || '|' || rule_name || '|' || (rule_spec->>'kind') || '|' || COALESCE(column_name, '<NULL>') FROM cdc.administration_request WHERE administration_request_id = '$alt_set_id'")" = "set_transformation|pending|umbenennung|rename_column|<NULL>" ] \
  || fail "Lauf 5: cdc.set_transformation schreibt unter cdc_admin nach dem Upgrade über $ALT_TAG nicht den erwarteten pending-Antrag"
alt_remove_id=$(psql_as "$ALT_DB" cdc_admin "SELECT cdc.remove_transformation('alttag-src', 'public', 't', 'umbenennung')")
[ "$(psql_q "$ALT_DB" "SELECT request_kind || '|' || status || '|' || rule_name || '|' || COALESCE(rule_spec::text, '<NULL>') FROM cdc.administration_request WHERE administration_request_id = '$alt_remove_id'")" = "remove_transformation|pending|umbenennung|<NULL>" ] \
  || fail "Lauf 5: cdc.remove_transformation schreibt unter cdc_admin nach dem Upgrade über $ALT_TAG nicht den erwarteten pending-Antrag"
if alt_denied=$(psql_as "$ALT_DB" cdc_reader "SELECT cdc.set_transformation('alttag-src', 'public', 't', 'verboten', '{}'::json)" 2>&1); then
  fail "Lauf 5: cdc_reader ruft cdc.set_transformation nach dem Upgrade über $ALT_TAG auf, ohne dass es scheitert"
fi
grep -q "permission denied for function" <<<"$alt_denied" \
  || fail "Lauf 5: der Aufruf unter cdc_reader scheitert nicht mit „permission denied for function“ ($alt_denied)"
if grep -q "Vorlauf" <<<"$work_out_1"; then
  alt_vorlauf="mit Vorlauf"
else
  alt_vorlauf="ohne Vorlauf"
fi
echo "run-schema-rollout-guard-test: Lauf 5 OK — Tag $ALT_TAG: Exit $alt_exit (Rollout des Tags), Exit $work_exit_1 (Arbeitsbaum, $alt_vorlauf), Exit $work_exit_2 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar; $alt_privilege_count Tabellen-/View-Rechte der drei Rollen auf administration_request, backfill_run und backfill_status, EXECUTE auf $alt_function_count Funktionen ($ALT_FUNCTION, $ALT_SET_FUNCTION, $ALT_REMOVE_FUNCTION) allein für cdc_admin (nicht PUBLIC), Spalten $alt_rule_columns (Alt-Zeile alttag-req: NULL), request_kind-Menge $alt_kinds_after, cdc.set_transformation/cdc.remove_transformation unter cdc_admin schreiben pending-Anträge, cdc_reader: permission denied for function"

echo "run-schema-rollout-guard-test: Lauf 6/6 (unbekannte Blocker, müssen mit Exit 8 abbrechen)"
echo "run-schema-rollout-guard-test: Lauf 6a (nicht deklarierte Funktion — die Wache lässt sie nicht unter --allow-destructive löschen)"
docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "CREATE FUNCTION cdc.zz_rolloutguard_unbekannt() RETURNS integer LANGUAGE sql AS 'SELECT 1'"
run_rollout . "$TARGET"
run6a_exit=$RUN_EXIT
[ "$run6a_exit" -ne 0 ] || fail "Lauf 6a lief durch (Exit 0), obwohl ein unbekannter destruktiver Blocker vorlag"
grep -qE "(Error|Fehler) 8" <<<"$RUN_OUT" || fail "Lauf 6a: der Abbruch trägt nicht den d-migrate-Exit 8 (Error 8/Fehler 8 fehlt in der Ausgabe)"
if grep -q "bekannte Fremdobjekt-Blocker (ADR-0043)" <<<"$RUN_OUT"; then
  fail "Lauf 6a: die Wache hat --allow-destructive für einen unbekannten Blocker freigegeben"
fi
[ "$(psql_q "$DB" "SELECT to_regprocedure('cdc.zz_rolloutguard_unbekannt()') IS NOT NULL")" = "t" ] \
  || fail "Lauf 6a: die nicht deklarierte Funktion wurde gelöscht — --allow-destructive lief für einen unbekannten Blocker"
docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 -c "DROP FUNCTION cdc.zz_rolloutguard_unbekannt()"

echo "run-schema-rollout-guard-test: Lauf 6b (nicht deklarierte Spalte — DropColumn-Blocker, Abbruch mit Exit 8)"
docker exec "$CONTAINER" psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 \
  -c "ALTER TABLE cdc.source ADD COLUMN _rolloutguard_test_col text"
run_rollout . "$TARGET"
run6b_exit=$RUN_EXIT
[ "$run6b_exit" -ne 0 ] || fail "Lauf 6b lief durch (Exit 0), obwohl ein unbekannter Blocker vorlag"
grep -qE "(Error|Fehler) 8" <<<"$RUN_OUT" || fail "Lauf 6b: der Abbruch trägt nicht den d-migrate-Exit 8 (Error 8/Fehler 8 fehlt in der Ausgabe)"

echo "run-schema-rollout-guard-test: OK — alle Belege real erbracht (Idempotenz-Allow, echte Änderung bleibt wirksam, View-Signatur-Vorlauf, Alt-Tag $ALT_TAG, Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit $run6a_exit/$run6b_exit mit d-migrate-Exit 8)"
