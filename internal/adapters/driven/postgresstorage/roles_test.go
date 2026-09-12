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
