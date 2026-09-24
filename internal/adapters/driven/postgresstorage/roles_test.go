package postgresstorage_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Die Rollen-Tests belegen die Least-Privilege-Trennung nach
// LH-QA-SEC-001…003 gegen die reale PostgreSQL-Instanz (`make test-store`,
// `ADR-0030`); ohne DSN überspringen sie. Die drei Rollen (`cdc_capture`,
// `cdc_admin`, `cdc_reader`) trägt der d-migrate-Rollout
// (`tools/schema/nacharbeit-roles.sql`), dieselbe Kette wie die drei
// Lese-Views (`sqlviews_test.go`).
//
// Die Rollen bleiben NOLOGIN (Gruppenrollen ohne Anmelde-Identität,
// `tools/schema/nacharbeit-roles.sql`) — die Prüfung wechselt über
// `SET ROLE` innerhalb der Testverbindung, für die der verbindende Nutzer
// (Bootstrap-Superuser des Testcontainers) jede Rolle annehmen darf, ohne
// dass es dafür ein Passwort braucht. `SET ROLE` bindet an die konkrete
// physische Verbindung, nicht an den Pool — die Tests holen deshalb eine
// einzelne Verbindung (`Acquire`) statt über den Pool zu fahren, und setzen
// die Rolle vor jeder Release wieder zurück, damit kein Folgetest eine
// fremde Rollen-Sitzung erbt.
const rolesTestSource = "src-roles"

func newTestRolesConn(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	var roles int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM pg_roles WHERE rolname IN ('cdc_capture', 'cdc_admin', 'cdc_reader')",
	).Scan(&roles); err != nil {
		t.Fatalf("Rollen-Prüfung: %v", err)
	}
	if roles != 3 {
		t.Fatalf("Least-Privilege-Rollen unvollständig (%d von 3) — der Schema-Rollout über make schema-rollout trägt sie (tools/schema/nacharbeit-roles.sql); der test-store-Lauf rollt sie vor dem Testlauf aus", roles)
	}

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Rollen-Quelle') ON CONFLICT (source_id) DO NOTHING",
		rolesTestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	return pool
}

// permissionDenied prüft die PostgreSQL-Fehlerklasse `insufficient_privilege`
// (SQLSTATE 42501) — dieselbe Textprüfung wie store_test.go für SQLSTATE.
func permissionDenied(err error) bool {
	return err != nil && strings.Contains(err.Error(), "SQLSTATE 42501")
}

// TestCdcCaptureRoleWritesChangesNotSchema belegt LH-QA-SEC-001/002:
// cdc_capture trägt den Erfassungspfad (INSERT auf transaction), aber
// keine Verwaltungsrechte (INSERT auf source_table bleibt cdc_admin
// vorbehalten).
func TestCdcCaptureRoleWritesChangesNotSchema(t *testing.T) {
	pool := newTestRolesConn(t)
	ctx := context.Background()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Verbindung reservieren: %v", err)
	}
	defer func() {
		_, _ = conn.Exec(ctx, "RESET ROLE")
		conn.Release()
	}()

	if _, err := conn.Exec(ctx, "SET ROLE cdc_capture"); err != nil {
		t.Fatalf("SET ROLE cdc_capture: %v", err)
	}

	if _, err := conn.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ($1, $2, 700)",
		"tx-capture-role", rolesTestSource,
	); err != nil {
		t.Fatalf("cdc_capture INSERT auf transaction: erwartet Erfolg, %v", err)
	}

	if _, err := conn.Exec(ctx,
		"INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', 'orders')",
		"st-capture-role", rolesTestSource,
	); !permissionDenied(err) {
		t.Fatalf("cdc_capture INSERT auf source_table: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

// TestCdcAdminRoleManagesSchemaNotChanges belegt LH-QA-SEC-001/002:
// cdc_admin trägt den Registrierungspfad (INSERT auf source_table), aber
// keinen Erfassungszugriff (INSERT auf change bleibt cdc_capture
// vorbehalten) — Admin ersetzt Capture nicht, dieselbe Trennung in beide
// Richtungen.
func TestCdcAdminRoleManagesSchemaNotChanges(t *testing.T) {
	pool := newTestRolesConn(t)
	ctx := context.Background()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Verbindung reservieren: %v", err)
	}
	defer func() {
		_, _ = conn.Exec(ctx, "RESET ROLE")
		conn.Release()
	}()

	if _, err := conn.Exec(ctx, "SET ROLE cdc_admin"); err != nil {
		t.Fatalf("SET ROLE cdc_admin: %v", err)
	}

	if _, err := conn.Exec(ctx,
		"INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', 'orders')",
		"st-admin-role", rolesTestSource,
	); err != nil {
		t.Fatalf("cdc_admin INSERT auf source_table: erwartet Erfolg, %v", err)
	}

	if _, err := conn.Exec(ctx,
		`INSERT INTO cdc.change
		    (change_id, transaction_id, source_table_id, sequence, operation, old_data, new_data, schema_version)
		 VALUES ($1, 'tx-missing', $2, 1, 'INSERT', NULL, '{}'::jsonb, 'sv-missing')`,
		"c-admin-role", "st-admin-role",
	); !permissionDenied(err) {
		t.Fatalf("cdc_admin INSERT auf change: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

// TestCdcAdminPublicationRequiresTableOwnership belegt: das GRANT
// CREATE ON DATABASE trägt `CREATE PUBLICATION` selbst, aber nicht das
// Hinzufügen einer Tabelle über `FOR TABLE` — PostgreSQL verlangt dafür
// zusätzlich Eigentümerrechte
// an der Zieltabelle (real gegen die Quelltabelle geprüft, nicht gegen
// cdc-Schema-Objekte, denn die zu aktivierenden Tabellen liegen außerhalb
// des cdc-Schemas). Ohne Eigentümerschaft schlägt der Aufruf fehl
// (SQLSTATE 42501 „must be owner of table"); nach Eigentümer-Übertragung
// (die Betriebs-Vorbedingung aus tools/schema/nacharbeit-roles.sql)
// gelingt er — das belegt, dass die dort dokumentierte
// Grant-Strategie tatsächlich trägt.
func TestCdcAdminPublicationRequiresTableOwnership(t *testing.T) {
	pool := newTestRolesConn(t)
	ctx := context.Background()
	const (
		sourceTable = "public.roles_pub_owner_test"
		publication = "roles_admin_owner_test_pub"
	)

	if _, err := pool.Exec(ctx, "DROP TABLE IF EXISTS "+sourceTable); err != nil {
		t.Fatalf("Vorab-Aufräumen der Quelltabelle: %v", err)
	}
	if _, err := pool.Exec(ctx, "CREATE TABLE "+sourceTable+" (id int PRIMARY KEY)"); err != nil {
		t.Fatalf("Quelltabelle anlegen: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_, _ = pool.Exec(cleanupCtx, "DROP PUBLICATION IF EXISTS "+publication)
		_, _ = pool.Exec(cleanupCtx, "DROP TABLE IF EXISTS "+sourceTable)
	})

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Verbindung reservieren: %v", err)
	}
	defer func() {
		_, _ = conn.Exec(ctx, "RESET ROLE")
		conn.Release()
	}()

	if _, err := conn.Exec(ctx, "SET ROLE cdc_admin"); err != nil {
		t.Fatalf("SET ROLE cdc_admin: %v", err)
	}

	// Ohne Eigentümerschaft: CREATE ON DATABASE allein reicht nicht.
	if _, err := conn.Exec(ctx,
		"CREATE PUBLICATION "+publication+" FOR TABLE "+sourceTable+" WITH (publish = 'insert, update, delete')",
	); !permissionDenied(err) {
		t.Fatalf("cdc_admin CREATE PUBLICATION … FOR TABLE ohne Eigentümerschaft: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}

	// Betriebs-Vorbedingung herstellen (nacharbeit-roles.sql-Kommentar):
	// Eigentümerschaft der Quelltabelle an cdc_admin übertragen — dafür
	// zurück zur reservierenden (superuser-artigen) Identität dieser
	// Verbindung.
	if _, err := conn.Exec(ctx, "RESET ROLE"); err != nil {
		t.Fatalf("RESET ROLE vor Eigentümer-Übertragung: %v", err)
	}
	if _, err := conn.Exec(ctx, "ALTER TABLE "+sourceTable+" OWNER TO cdc_admin"); err != nil {
		t.Fatalf("Eigentümer-Übertragung: %v", err)
	}
	if _, err := conn.Exec(ctx, "SET ROLE cdc_admin"); err != nil {
		t.Fatalf("SET ROLE cdc_admin (2. Versuch): %v", err)
	}

	// Mit Eigentümerschaft: derselbe Aufruf gelingt.
	if _, err := conn.Exec(ctx,
		"CREATE PUBLICATION "+publication+" FOR TABLE "+sourceTable+" WITH (publish = 'insert, update, delete')",
	); err != nil {
		t.Fatalf("cdc_admin CREATE PUBLICATION … FOR TABLE nach Eigentümer-Übertragung: erwartet Erfolg, %v", err)
	}
}

// TestCdcAdminAlterPublicationAddTableRequiresTableOwnership belegt
// dieselbe Grenze für den zweiten betroffenen Aufruf
// (tableactivation.go:231, `ALTER PUBLICATION … ADD TABLE`) — eine
// bestehende Publication, eine zweite Quelltabelle ohne Eigentümerschaft.
func TestCdcAdminAlterPublicationAddTableRequiresTableOwnership(t *testing.T) {
	pool := newTestRolesConn(t)
	ctx := context.Background()
	const (
		ownedTable   = "public.roles_pub_add_owned_test"
		secondTable  = "public.roles_pub_add_second_test"
		publication2 = "roles_admin_add_test_pub"
	)

	if _, err := pool.Exec(ctx, "DROP TABLE IF EXISTS "+ownedTable+", "+secondTable); err != nil {
		t.Fatalf("Vorab-Aufräumen der Quelltabellen: %v", err)
	}
	if _, err := pool.Exec(ctx, "CREATE TABLE "+ownedTable+" (id int PRIMARY KEY)"); err != nil {
		t.Fatalf("erste Quelltabelle anlegen: %v", err)
	}
	if _, err := pool.Exec(ctx, "CREATE TABLE "+secondTable+" (id int PRIMARY KEY)"); err != nil {
		t.Fatalf("zweite Quelltabelle anlegen: %v", err)
	}
	if _, err := pool.Exec(ctx, "ALTER TABLE "+ownedTable+" OWNER TO cdc_admin"); err != nil {
		t.Fatalf("Eigentümer-Übertragung erste Tabelle: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_, _ = pool.Exec(cleanupCtx, "DROP PUBLICATION IF EXISTS "+publication2)
		_, _ = pool.Exec(cleanupCtx, "DROP TABLE IF EXISTS "+ownedTable+", "+secondTable)
	})

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Verbindung reservieren: %v", err)
	}
	defer func() {
		_, _ = conn.Exec(ctx, "RESET ROLE")
		conn.Release()
	}()

	if _, err := conn.Exec(ctx, "SET ROLE cdc_admin"); err != nil {
		t.Fatalf("SET ROLE cdc_admin: %v", err)
	}
	if _, err := conn.Exec(ctx,
		"CREATE PUBLICATION "+publication2+" FOR TABLE "+ownedTable+" WITH (publish = 'insert, update, delete')",
	); err != nil {
		t.Fatalf("cdc_admin CREATE PUBLICATION mit eigener Tabelle: erwartet Erfolg, %v", err)
	}

	// Zweite Tabelle ohne Eigentümerschaft hinzufügen: derselbe
	// Eigentümer-Zwang gilt auch für ADD TABLE, nicht nur für FOR TABLE.
	if _, err := conn.Exec(ctx,
		"ALTER PUBLICATION "+publication2+" ADD TABLE "+secondTable,
	); !permissionDenied(err) {
		t.Fatalf("cdc_admin ALTER PUBLICATION … ADD TABLE ohne Eigentümerschaft: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

// TestCdcReaderRoleReadsViewsNotBaseTables belegt LH-QA-SEC-003: der
// Lesezugriff trägt ausschließlich die Views — SELECT direkt auf einer
// Basistabelle scheitert, obwohl dieselben Zeilen über cdc.metrics lesbar
// sind (Definer-Semantik der View, tools/schema/nacharbeit-roles.sql).
func TestCdcReaderRoleReadsViewsNotBaseTables(t *testing.T) {
	pool := newTestRolesConn(t)
	ctx := context.Background()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Verbindung reservieren: %v", err)
	}
	defer func() {
		_, _ = conn.Exec(ctx, "RESET ROLE")
		conn.Release()
	}()

	if _, err := conn.Exec(ctx, "SET ROLE cdc_reader"); err != nil {
		t.Fatalf("SET ROLE cdc_reader: %v", err)
	}

	var rows int
	if err := conn.QueryRow(ctx, "SELECT count(*) FROM cdc.metrics").Scan(&rows); err != nil {
		t.Fatalf("cdc_reader SELECT auf cdc.metrics: erwartet Erfolg, %v", err)
	}
	if rows == 0 {
		t.Fatalf("cdc.metrics trägt keine Zeilen — erwartet mindestens die aggregierten Zähl-Metriken")
	}

	// cdc.retention_blockers (LH-FA-RET-005) trägt dieselbe
	// View-Owner-Lese-Disziplin wie cdc.metrics/cdc.active_tables/
	// cdc.consumer_status — ein leeres Ergebnis ist hier zulässig (kein
	// Consumer hat in dieser Fixtur bestätigt), geprüft wird ausschließlich
	// die Zugriffserlaubnis.
	if _, err := conn.Exec(ctx, "SELECT count(*) FROM cdc.retention_blockers"); err != nil {
		t.Fatalf("cdc_reader SELECT auf cdc.retention_blockers: erwartet Erfolg, %v", err)
	}

	if _, err := conn.Exec(ctx, "SELECT count(*) FROM cdc.transaction"); !permissionDenied(err) {
		t.Fatalf("cdc_reader SELECT auf cdc.transaction (Basistabelle): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}

	if _, err := conn.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ($1, $2, 900)",
		"tx-reader-role", rolesTestSource,
	); !permissionDenied(err) {
		t.Fatalf("cdc_reader INSERT auf transaction: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

// TestMetricsViewCarriesConsumerLag belegt das Metriken-Minimum
// (LH-FA-SST-004/LH-QA-OPS-003, SPEC-009): cdc_consumer_lag trägt die
// Differenz aus bestätigter und letzter Commit-Position je Consumer,
// dieselbe Rechnung wie cdc.consumer_status.
func TestMetricsViewCarriesConsumerLag(t *testing.T) {
	pool := newTestRolesConn(t)
	ctx := context.Background()
	const (
		consumerID    = "mc-lag"
		transactionID = "tx-metrics-lag"
	)

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ($1, $2, 950)",
		transactionID, rolesTestSource,
	); err != nil {
		t.Fatalf("transaction-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.consumer (consumer_id, name) VALUES ($1, 'Metrics-Consumer')",
		consumerID,
	); err != nil {
		t.Fatalf("consumer-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.consumer_position (consumer_id, source_id, acknowledged_position) VALUES ($1, $2, 650)",
		consumerID, rolesTestSource,
	); err != nil {
		t.Fatalf("consumer_position-Zeile: %v", err)
	}

	var lag float64
	if err := pool.QueryRow(ctx,
		"SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_consumer_lag' AND label = $1",
		consumerID,
	).Scan(&lag); err != nil {
		t.Fatalf("cdc_consumer_lag-Lesen: %v", err)
	}
	if lag != 300 {
		t.Fatalf("cdc_consumer_lag = %v, wollen 300 (950 - 650)", lag)
	}
}

// TestMetricsViewCarriesStorageBytes belegt LH-FA-RET-006 (SPEC-009
// cdc_storage_bytes): cdc.metrics trägt einen realen, positiven Wert für
// die physische Speichergröße von cdc.change über pg_relation_size, nach
// dem Einfügen einer vollständigen Change-Zeile (Transaktion,
// Tabellen-/Schema-Referenz, Change selbst) — ein numerischer Wert, kein
// Zeilenzähler.
func TestMetricsViewCarriesStorageBytes(t *testing.T) {
	pool := newTestRolesConn(t)
	ctx := context.Background()
	const (
		sourceTableID   = "st-storage-bytes"
		schemaVersionID = "sv-storage-bytes"
		transactionID   = "tx-storage-bytes"
		changeID        = "ch-storage-bytes"
	)

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', 'storage_bytes')",
		sourceTableID, rolesTestSource,
	); err != nil {
		t.Fatalf("source_table-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ($1, $2, 1)",
		schemaVersionID, sourceTableID,
	); err != nil {
		t.Fatalf("schema_version-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ($1, $2, 970)",
		transactionID, rolesTestSource,
	); err != nil {
		t.Fatalf("transaction-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.change (change_id, transaction_id, source_table_id, sequence, operation, new_data, schema_version) VALUES ($1, $2, $3, 1, 'INSERT', '{}', $4)",
		changeID, transactionID, sourceTableID, schemaVersionID,
	); err != nil {
		t.Fatalf("change-Zeile: %v", err)
	}

	var storageBytes float64
	if err := pool.QueryRow(ctx,
		"SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_storage_bytes'",
	).Scan(&storageBytes); err != nil {
		t.Fatalf("cdc_storage_bytes-Lesen: %v", err)
	}
	if storageBytes <= 0 {
		t.Fatalf("cdc_storage_bytes = %v, wollen > 0 nach dem Einfügen von Testdaten", storageBytes)
	}
}

// Die drei Tests unten belegen den Rollenschnitt auf `cdc.backfill_run`
// (`SPEC-029`, `ADR-0113` Festlegung 1, `LH-QA-SEC-001`…`003`) je Rolle
// real: was die Rolle darf, gelingt; was sie nicht darf, scheitert mit
// SQLSTATE 42501. Rot färbende Mutation je Zeile der Grants in
// `tools/schema/nacharbeit-roles.sql`: `INSERT` an `cdc_capture`, `UPDATE`
// oder `DELETE` an `cdc_admin`, `SELECT` an `cdc_reader` ergänzen bzw. den
// erlaubten Grant streichen — der jeweilige Test meldet die abweichende
// Zeile.
const backfillRoleRun = "roles-backfill-run"

// seedBackfillRoleRun legt als Superuser die Run-Zeile an, an der die
// Rollen-Tests schreiben und lesen, und räumt sie ab.
func seedBackfillRoleRun(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.backfill_run WHERE run_id = $1", backfillRoleRun); err != nil {
		t.Fatalf("Vorab-Aufräumen der Run-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		`INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status)
		 VALUES ($1, $2, 'public', 'roles_backfill', 'queued')`, backfillRoleRun, rolesTestSource,
	); err != nil {
		t.Fatalf("Run-Zeile anlegen: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.backfill_run WHERE run_id LIKE $1", backfillRoleRun+"%")
	})
}

// asRole reserviert eine Verbindung und setzt die Rolle; das Zurücksetzen
// übernimmt die Bereinigung des Tests.
func asRole(t *testing.T, pool *pgxpool.Pool, role string) *pgxpool.Conn {
	t.Helper()
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Verbindung reservieren: %v", err)
	}
	t.Cleanup(func() {
		_, _ = conn.Exec(ctx, "RESET ROLE")
		conn.Release()
	})
	if _, err := conn.Exec(ctx, "SET ROLE "+role); err != nil {
		t.Fatalf("SET ROLE %s: %v", role, err)
	}
	return conn
}

func TestCdcCaptureRoleUpdatesBackfillRunButNeitherInsertsNorDeletes(t *testing.T) {
	pool := newTestRolesConn(t)
	seedBackfillRoleRun(t, pool)
	conn := asRole(t, pool, "cdc_capture")
	ctx := context.Background()

	var status string
	if err := conn.QueryRow(ctx, "SELECT status FROM cdc.backfill_run WHERE run_id = $1", backfillRoleRun).Scan(&status); err != nil {
		t.Fatalf("cdc_capture SELECT auf backfill_run: erwartet Erfolg, %v", err)
	}
	tag, err := conn.Exec(ctx, "UPDATE cdc.backfill_run SET status = 'running', started_at = now() WHERE run_id = $1", backfillRoleRun)
	if err != nil || tag.RowsAffected() != 1 {
		t.Fatalf("cdc_capture UPDATE auf backfill_run: erwartet Erfolg auf 1 Zeile, %v (Zeilen %d)", err, tag.RowsAffected())
	}
	if _, err := conn.Exec(ctx,
		`INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status)
		 VALUES ($1, $2, 'public', 'roles_backfill_capture', 'queued')`, backfillRoleRun+"-capture", rolesTestSource,
	); !permissionDenied(err) {
		t.Fatalf("cdc_capture INSERT auf backfill_run: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
	if _, err := conn.Exec(ctx, "DELETE FROM cdc.backfill_run WHERE run_id = $1", backfillRoleRun); !permissionDenied(err) {
		t.Fatalf("cdc_capture DELETE auf backfill_run: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

func TestCdcAdminRoleInsertsBackfillRunButNeitherUpdatesNorDeletes(t *testing.T) {
	pool := newTestRolesConn(t)
	seedBackfillRoleRun(t, pool)
	conn := asRole(t, pool, "cdc_admin")
	ctx := context.Background()

	if _, err := conn.Exec(ctx,
		`INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status)
		 VALUES ($1, $2, 'public', 'roles_backfill_admin', 'queued')`, backfillRoleRun+"-admin", rolesTestSource,
	); err != nil {
		t.Fatalf("cdc_admin INSERT auf backfill_run: erwartet Erfolg, %v", err)
	}
	var status string
	if err := conn.QueryRow(ctx, "SELECT status FROM cdc.backfill_run WHERE run_id = $1", backfillRoleRun).Scan(&status); err != nil {
		t.Fatalf("cdc_admin SELECT auf backfill_run: erwartet Erfolg, %v", err)
	}
	if _, err := conn.Exec(ctx, "UPDATE cdc.backfill_run SET status = 'running' WHERE run_id = $1", backfillRoleRun); !permissionDenied(err) {
		t.Fatalf("cdc_admin UPDATE auf backfill_run: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
	if _, err := conn.Exec(ctx, "DELETE FROM cdc.backfill_run WHERE run_id = $1", backfillRoleRun); !permissionDenied(err) {
		t.Fatalf("cdc_admin DELETE auf backfill_run: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

func TestCdcReaderRoleCannotReadTheBackfillRunBaseTable(t *testing.T) {
	pool := newTestRolesConn(t)
	seedBackfillRoleRun(t, pool)
	conn := asRole(t, pool, "cdc_reader")

	var status string
	if err := conn.QueryRow(context.Background(), "SELECT status FROM cdc.backfill_run WHERE run_id = $1", backfillRoleRun).Scan(&status); !permissionDenied(err) {
		t.Fatalf("cdc_reader SELECT auf backfill_run: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

// `cdc_reader` liest den Run-Zustand über die View `cdc.backfill_status`
// (`SPEC-029`): die View trägt die Zeile des Runs samt der zwei Warn-Spalten
// (`false`) und der unbekannten Schätzung als NULL, obwohl `cdc_reader` kein
// Recht auf die Basistabelle hat (Definer-Semantik der Views). Rot färbende
// Mutation: `cdc.backfill_status` aus dem Reader-Grant der Rollen-Datei
// streichen — der Lesezugriff endet mit SQLSTATE 42501.
func TestCdcReaderRoleReadsTheBackfillRunThroughTheStatusView(t *testing.T) {
	pool := newTestRolesConn(t)
	seedBackfillRoleRun(t, pool)
	conn := asRole(t, pool, "cdc_reader")

	var (
		status                 string
		estimated              *int64
		warnSize, warnDuration bool
	)
	err := conn.QueryRow(context.Background(),
		"SELECT status, estimated_rows, warn_estimated_size, warn_duration FROM cdc.backfill_status WHERE run_id = $1",
		backfillRoleRun,
	).Scan(&status, &estimated, &warnSize, &warnDuration)
	if err != nil {
		t.Fatalf("cdc_reader SELECT auf cdc.backfill_status: erwartet Erfolg, %v", err)
	}
	if status != "queued" || estimated != nil || warnSize || warnDuration {
		t.Fatalf("cdc.backfill_status trägt %q, Schätzung %v, Warnungen %v/%v — erwartet queued, unbekannt (NULL), false/false",
			status, estimated, warnSize, warnDuration)
	}
}

// Die Tests unten belegen den Rollenschnitt auf
// `cdc.administration_request` (`ADR-0050`, `LH-QA-SEC-001`…`003`) je Rolle
// real: `cdc_admin` liest und vermerkt Anträge (`SELECT`, `UPDATE`), legt
// keine an und löscht keine; `cdc_capture` und `cdc_reader` tragen kein
// Recht auf die Tabelle. Rot färbende Mutation je Zeile der Grants in
// `tools/schema/nacharbeit-roles.sql`: den Grant an `cdc_admin` streichen
// bzw. auf `SELECT` kürzen, `INSERT`/`DELETE` an `cdc_admin` oder `SELECT`
// an `cdc_capture`/`cdc_reader` ergänzen — der jeweilige Test meldet die
// abweichende Anweisung.
const administrationRoleRequest = "roles-administration-request"

// seedAdministrationRoleRequest legt als Superuser den offenen Antrag an,
// an dem die Rollen-Tests lesen und schreiben, und räumt ihn ab.
func seedAdministrationRoleRequest(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.administration_request WHERE administration_request_id LIKE $1", administrationRoleRequest+"%"); err != nil {
		t.Fatalf("Vorab-Aufräumen des Antrags: %v", err)
	}
	if _, err := pool.Exec(ctx,
		`INSERT INTO cdc.administration_request (administration_request_id, source_id, schema_name, table_name, request_kind, status)
		 VALUES ($1, $2, 'public', 'roles_administration', 'enable', 'pending')`, administrationRoleRequest, rolesTestSource,
	); err != nil {
		t.Fatalf("Antrag anlegen: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE administration_request_id LIKE $1", administrationRoleRequest+"%")
	})
}

func TestCdcAdminRoleReadsAndMarksAdministrationRequestsButNeitherInsertsNorDeletes(t *testing.T) {
	pool := newTestRolesConn(t)
	seedAdministrationRoleRequest(t, pool)
	conn := asRole(t, pool, "cdc_admin")
	ctx := context.Background()

	var status string
	if err := conn.QueryRow(ctx,
		"SELECT status FROM cdc.administration_request WHERE administration_request_id = $1", administrationRoleRequest,
	).Scan(&status); err != nil {
		t.Fatalf("cdc_admin SELECT auf administration_request: erwartet Erfolg, %v", err)
	}
	tag, err := conn.Exec(ctx,
		"UPDATE cdc.administration_request SET status = 'applied', error_message = NULL WHERE administration_request_id = $1 AND status = 'pending'",
		administrationRoleRequest)
	if err != nil || tag.RowsAffected() != 1 {
		t.Fatalf("cdc_admin UPDATE auf administration_request: erwartet Erfolg auf 1 Zeile, %v (Zeilen %d)", err, tag.RowsAffected())
	}
	if _, err := conn.Exec(ctx,
		`INSERT INTO cdc.administration_request (administration_request_id, source_id, schema_name, table_name, request_kind, status)
		 VALUES ($1, $2, 'public', 'roles_administration_admin', 'enable', 'pending')`, administrationRoleRequest+"-admin", rolesTestSource,
	); !permissionDenied(err) {
		t.Fatalf("cdc_admin INSERT auf administration_request: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
	if _, err := conn.Exec(ctx, "DELETE FROM cdc.administration_request WHERE administration_request_id = $1", administrationRoleRequest); !permissionDenied(err) {
		t.Fatalf("cdc_admin DELETE auf administration_request: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

func TestCdcCaptureAndReaderRolesCannotTouchAdministrationRequests(t *testing.T) {
	for _, role := range []string{"cdc_capture", "cdc_reader"} {
		t.Run(role, func(t *testing.T) {
			pool := newTestRolesConn(t)
			seedAdministrationRoleRequest(t, pool)
			conn := asRole(t, pool, role)
			ctx := context.Background()

			var status string
			if err := conn.QueryRow(ctx,
				"SELECT status FROM cdc.administration_request WHERE administration_request_id = $1", administrationRoleRequest,
			).Scan(&status); !permissionDenied(err) {
				t.Fatalf("%s SELECT auf administration_request: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", role, err)
			}
			if _, err := conn.Exec(ctx,
				"UPDATE cdc.administration_request SET status = 'applied' WHERE administration_request_id = $1", administrationRoleRequest,
			); !permissionDenied(err) {
				t.Fatalf("%s UPDATE auf administration_request: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", role, err)
			}
		})
	}
}
