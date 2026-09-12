package bootstrap_test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
)

// Die Rollen-Verdrahtungs-Tests belegen ADR-0047 §Fitness Function real
// gegen die Testcontainer-Instanz (`make test-store`): ein echter
// Verbindungsaufbau über eine anmeldefähige Identität, die Mitglied der
// jeweiligen Gruppenrolle ist (`cdc_capture`/`cdc_admin`/`cdc_reader`,
// NOLOGIN, `tools/schema/nacharbeit-roles.sql`) — dieselbe Betriebs-
// Vorbedingung, die ADR-0047 Kontext-Befund 1 offenlässt: „wie der
// Betreiber diese Login-Identität anlegt". Die Login-Rolle ist ein
// Test-Fixture (angelegt und wieder entfernt innerhalb des Tests), kein
// Bestandteil des Rollen-DDL (`nacharbeit-roles.sql`) — dieselbe
// Unterscheidung wie die Eigentümerschafts-Übertragung in
// `postgresstorage/roles_test.go` (Fixture, kein Rollenschnitt).
//
// Diese Tests laufen im Paket `internal/bootstrap`, weil
// `tools/harness/run-store-tests.sh` genau dieses Paket vorgezogen und
// isoliert ausführt (`BEO-PGC/test-isolation-geteilter-zustand`) — der
// testweise Grant-Entzug/-Erneuerung in
// `TestCdcAdminHeartbeatWriteRequiresGrant` bleibt dadurch ohne
// Race-Risiko gegen die `postgresstorage`-Pakettests, die erst danach
// laufen.
const rolesWiringTestSource = "src-roles-wiring"

// adminTestPool baut den Verbindungspool über die Bootstrap-Superuser-DSN
// des Testcontainers (`CDC_STORE_TEST_DSN`) — dieselbe Identität, mit der
// `tools/schema/nacharbeit-roles.sql` die drei Gruppenrollen anlegt und
// die deshalb Rollen-Mitgliedschaft an eine neue Login-Identität granten
// darf.
func adminTestPool(t *testing.T) (*pgxpool.Pool, string) {
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
	return pool, dsn
}

// newTestLoginRole legt eine anmeldefähige Test-Identität an, die
// Mitglied der übergebenen Gruppenrolle ist, und liefert deren DSN
// (gleicher Host/Datenbank/SSL-Modus wie baseDSN, andere
// Anmelde-Identität) — der reale „Verbindungsaufbau", den ADR-0047 §Fitness
// Function verlangt. Aufräumen läuft über t.Cleanup, registriert nach dem
// Aufruf, der adminPool öffnet (LIFO-Reihenfolge: die Login-Rolle
// verschwindet, bevor der adminPool schließt).
func newTestLoginRole(t *testing.T, adminPool *pgxpool.Pool, baseDSN, loginName, groupRole string) string {
	t.Helper()
	ctx := context.Background()
	const password = "test-login-password"

	if _, err := adminPool.Exec(ctx, fmt.Sprintf(`DROP ROLE IF EXISTS %q`, loginName)); err != nil {
		t.Fatalf("Vorab-Aufräumen der Login-Identität %s: %v", loginName, err)
	}
	if _, err := adminPool.Exec(ctx,
		fmt.Sprintf(`CREATE ROLE %q LOGIN PASSWORD '%s' IN ROLE %q`, loginName, password, groupRole),
	); err != nil {
		t.Fatalf("Login-Identität %s anlegen: %v", loginName, err)
	}
	t.Cleanup(func() {
		_, _ = adminPool.Exec(context.Background(), fmt.Sprintf(`DROP ROLE IF EXISTS %q`, loginName))
	})

	parsed, err := url.Parse(baseDSN)
	if err != nil {
		t.Fatalf("Basis-DSN %q nicht parsebar: %v", baseDSN, err)
	}
	parsed.User = url.UserPassword(loginName, password)
	return parsed.String()
}

// permissionDenied prüft dieselbe PostgreSQL-Fehlerklasse wie
// `postgresstorage/roles_test.go` (SQLSTATE 42501,
// insufficient_privilege) — hier auch auf den vom Heartbeat-Adapter
// gewrappten Fehler (`heartbeatStorageFailure`), dessen Text die
// SQLSTATE-Angabe der Ursache weiterträgt.
func permissionDenied(err error) bool {
	return err != nil && strings.Contains(err.Error(), "SQLSTATE 42501")
}

// TestCdcReaderLoginConnectionRejectsWrite belegt ADR-0047 §Fitness
// Function, erste Zeile: ein Verbindungsaufbau mit der
// `cdc_reader`-Login-Identität scheitert an einem schreibenden Aufruf.
func TestCdcReaderLoginConnectionRejectsWrite(t *testing.T) {
	adminPool, baseDSN := adminTestPool(t)
	ctx := context.Background()

	loginDSN := newTestLoginRole(t, adminPool, baseDSN, "pgc_test_reader_login", "cdc_reader")

	readerPool, err := pgxpool.New(ctx, loginDSN)
	if err != nil {
		t.Fatalf("Verbindungsaufbau mit cdc_reader-Login-Identität: %v", err)
	}
	defer readerPool.Close()
	if err := readerPool.Ping(ctx); err != nil {
		t.Fatalf("Ping mit cdc_reader-Login-Identität: %v", err)
	}

	_, err = readerPool.Exec(ctx,
		`INSERT INTO cdc.change
		    (change_id, transaction_id, source_table_id, sequence, operation, old_data, new_data, schema_version)
		 VALUES ('c-reader-login', 'tx-reader-login', 'st-reader-login', 1, 'INSERT', NULL, '{}'::jsonb, 'sv-reader-login')`,
	)
	if !permissionDenied(err) {
		t.Fatalf("cdc_reader-Login-Identität INSERT auf cdc.change: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

// TestCdcCaptureLoginConnectionRejectsAdminWrite belegt ADR-0047 §Fitness
// Function, erste Zeile: ein Verbindungsaufbau mit der
// `cdc_capture`-Login-Identität scheitert an einem administrativen
// Aufruf.
func TestCdcCaptureLoginConnectionRejectsAdminWrite(t *testing.T) {
	adminPool, baseDSN := adminTestPool(t)
	ctx := context.Background()

	loginDSN := newTestLoginRole(t, adminPool, baseDSN, "pgc_test_capture_login", "cdc_capture")

	capturePool, err := pgxpool.New(ctx, loginDSN)
	if err != nil {
		t.Fatalf("Verbindungsaufbau mit cdc_capture-Login-Identität: %v", err)
	}
	defer capturePool.Close()
	if err := capturePool.Ping(ctx); err != nil {
		t.Fatalf("Ping mit cdc_capture-Login-Identität: %v", err)
	}

	_, err = capturePool.Exec(ctx,
		"INSERT INTO cdc.consumer (consumer_id, name) VALUES ('c-capture-login', 'Capture-Login-Consumer')",
	)
	if !permissionDenied(err) {
		t.Fatalf("cdc_capture-Login-Identität INSERT auf cdc.consumer: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
	}
}

// TestCdcAdminHeartbeatWriteRequiresGrant belegt ADR-0047 §Fitness
// Function, zweite Zeile (Kontext-Befund 3): der Heartbeat-Schreibpfad
// gelingt nur, nachdem der ergänzende GRANT auf `cdc.process_heartbeat`
// ausgerollt ist. Der Test entzieht das Recht testweise und stellt es
// danach wieder her (derselbe Vorher/Nachher-Beleg wie
// `roles_test.go::TestCdcAdminPublicationRequiresTableOwnership`), damit
// der Schema-Rollout-Zustand für nachfolgende Tests unverändert bleibt.
func TestCdcAdminHeartbeatWriteRequiresGrant(t *testing.T) {
	adminPool, baseDSN := adminTestPool(t)
	ctx := context.Background()

	if _, err := adminPool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Rollen-Verdrahtungs-Quelle') ON CONFLICT (source_id) DO NOTHING",
		rolesWiringTestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}

	loginDSN := newTestLoginRole(t, adminPool, baseDSN, "pgc_test_admin_login", "cdc_admin")

	if _, err := adminPool.Exec(ctx, "REVOKE SELECT, INSERT, UPDATE ON cdc.process_heartbeat FROM cdc_admin"); err != nil {
		t.Fatalf("Grant testweise entziehen: %v", err)
	}
	t.Cleanup(func() {
		if _, err := adminPool.Exec(context.Background(),
			"GRANT SELECT, INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin",
		); err != nil {
			t.Fatalf("Grant nach Testlauf wiederherstellen: %v", err)
		}
	})

	adapterWithoutGrant, err := postgresstorage.NewHeartbeat(ctx, loginDSN)
	if err != nil {
		t.Fatalf("NewHeartbeat mit cdc_admin-Login-Identität: %v", err)
	}
	beatErr := adapterWithoutGrant.Beat(ctx, rolesWiringTestSource)
	adapterWithoutGrant.Close()
	if !permissionDenied(beatErr) {
		t.Fatalf("Beat ohne Grant: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", beatErr)
	}

	if _, err := adminPool.Exec(ctx, "GRANT SELECT, INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin"); err != nil {
		t.Fatalf("Grant wiederherstellen (vor dem zweiten Beat): %v", err)
	}

	adapterWithGrant, err := postgresstorage.NewHeartbeat(ctx, loginDSN)
	if err != nil {
		t.Fatalf("NewHeartbeat mit cdc_admin-Login-Identität (2. Versuch): %v", err)
	}
	defer adapterWithGrant.Close()
	if err := adapterWithGrant.Beat(ctx, rolesWiringTestSource); err != nil {
		t.Fatalf("Beat nach Grant: erwartet Erfolg, erhalten %v", err)
	}
}
