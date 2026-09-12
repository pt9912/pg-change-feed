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
	"github.com/pt9912/pg-change-feed/internal/bootstrap"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
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
// Function, zweite Zeile (Kontext-Befund 3), und ADR-0048s korrigierten
// Grant-Text real gegen die exakte Fehlerklasse, die ADR-0048 dokumentiert:
// drei Stufen — kein Recht (muss scheitern), `INSERT, UPDATE` ohne `SELECT` (der
// ursprüngliche, von ADR-0048 korrigierte ADR-0047-Text — muss ebenfalls
// scheitern, sonst zeigt der Test eine Regression auf genau diesen Text
// nicht an) und der vollständige `SELECT, INSERT, UPDATE`-Grant (muss
// gelingen). Der Test entzieht/erteilt das Recht testweise und stellt den
// vollständigen Grant danach wieder her (derselbe Vorher/Nachher-Beleg wie
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

	// Stufe 1: kein Recht — Beat muss scheitern.
	adapterWithoutGrant, err := postgresstorage.NewHeartbeat(ctx, loginDSN)
	if err != nil {
		t.Fatalf("NewHeartbeat mit cdc_admin-Login-Identität: %v", err)
	}
	beatErr := adapterWithoutGrant.Beat(ctx, rolesWiringTestSource)
	adapterWithoutGrant.Close()
	if !permissionDenied(beatErr) {
		t.Fatalf("Beat ohne Grant: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", beatErr)
	}

	// Stufe 2 (die dokumentierte Regressions-Stufe): exakt der
	// ursprüngliche ADR-0047-Text (`INSERT, UPDATE` ohne `SELECT`) — dieser
	// Grant lässt den `ON CONFLICT … DO UPDATE`-Zweig weiterhin mit
	// SQLSTATE 42501 scheitern (ADR-0048 §Kontext). Ein Test, der diese
	// Stufe nicht real prüft, bliebe grün, wenn `nacharbeit-roles.sql`
	// versehentlich auf diesen Text zurückfiele.
	if _, err := adminPool.Exec(ctx, "GRANT INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin"); err != nil {
		t.Fatalf("Grant auf den ADR-0047-Ursprungstext setzen: %v", err)
	}
	adapterInsertUpdateOnly, err := postgresstorage.NewHeartbeat(ctx, loginDSN)
	if err != nil {
		t.Fatalf("NewHeartbeat mit cdc_admin-Login-Identität (Stufe INSERT/UPDATE ohne SELECT): %v", err)
	}
	insertUpdateOnlyErr := adapterInsertUpdateOnly.Beat(ctx, rolesWiringTestSource)
	adapterInsertUpdateOnly.Close()
	if !permissionDenied(insertUpdateOnlyErr) {
		t.Fatalf("Beat mit INSERT/UPDATE ohne SELECT (ADR-0047-Ursprungstext, von ADR-0048 korrigiert): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", insertUpdateOnlyErr)
	}

	// Stufe 3: der von ADR-0048 korrigierte vollständige Grant — Beat muss
	// gelingen.
	if _, err := adminPool.Exec(ctx, "GRANT SELECT, INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin"); err != nil {
		t.Fatalf("Grant wiederherstellen (vor dem dritten Beat): %v", err)
	}

	adapterWithGrant, err := postgresstorage.NewHeartbeat(ctx, loginDSN)
	if err != nil {
		t.Fatalf("NewHeartbeat mit cdc_admin-Login-Identität (3. Versuch): %v", err)
	}
	defer adapterWithGrant.Close()
	if err := adapterWithGrant.Beat(ctx, rolesWiringTestSource); err != nil {
		t.Fatalf("Beat nach vollständigem Grant: erwartet Erfolg, erhalten %v", err)
	}
}

// TestCdcWiringCallerRejectsWrongRoleAssignment belegt ADR-0047
// §Entscheidung gegen Vertauschung: für jeden DML-fähigen Aufrufer aus der
// Zuordnungstabelle verbindet dieser Test mit der FALSCHEN Rolle und prüft,
// dass der reale Produktionsaufruf scheitert — dass er mit der richtigen
// Rolle gelingt, belegt für sich keine Vertauschungssicherheit.
//
// Zwei Prüftiefen, je nachdem, ob der Aufrufer eigenständig aufrufbar ist:
//   - RegisterConsumer, AcknowledgeConsumer, Healthcheck rufen die
//     tatsächliche `wiring.go`-Funktion auf (`bootstrap.RegisterConsumer`/
//     `AcknowledgeConsumer`/`Healthcheck`) mit einer `Config`, deren
//     einschlägiges DSN-Feld die falsche Rolle trägt — eine Vertauschung in
//     `wiring.go` selbst (z. B. `cfg.CaptureDSN` statt `cfg.AdminDSN`)
//     würde dieser Test real anzeigen.
//   - Store-Adapter, Tabellen-Aktivierung und Heartbeat-Adapter laufen nur
//     über `bootstrap.Run`, das den vollen Erfassungspfad braucht
//     (Replication-Stream, `wal_level=logical`) und dessen
//     Heartbeat-Fehlerpfad bewusst unterdrückt wird (`runHeartbeat`,
//     `_ = port.Beat(...)` — kein Abbruch des Capture-Pfads durch einen
//     Heartbeat-Fehler, `SPEC-008`). Diese drei Fälle bauen deshalb den
//     jeweiligen Adapter direkt mit der falschen Rolle und rufen seine
//     reale Produktionsmethode auf — das belegt „mit dieser Rolle scheitert
//     der Aufruf real", nicht „`wiring.go` weist diesem Aufrufer die
//     richtige Rolle zu" (Letzteres bleibt Code-Lektüre). Eine Vertauschung
//     ausschließlich dieser drei Zeilen in `wiring.go` bliebe von diesem
//     Test unentdeckt.
//
// Replication-Stream und ACK-Adapter (`cdc_capture`, REPLICATION-Attribut)
// bleiben ganz außerhalb dieser Tabelle: ihre reale Prüfung braucht eine
// Instanz mit `wal_level=logical` (`make test-replication`), die die drei
// Login-Test-Identitäten dieser Datei nicht bereitstellt — eine offene,
// benannte Lücke, kein stiller Auslassungsfall.
func TestCdcWiringCallerRejectsWrongRoleAssignment(t *testing.T) {
	adminPool, baseDSN := adminTestPool(t)
	ctx := context.Background()

	if _, err := adminPool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Rollen-Verdrahtungs-Quelle') ON CONFLICT (source_id) DO NOTHING",
		rolesWiringTestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}

	cases := []struct {
		name      string
		wrongRole string
		loginName string
		action    func(t *testing.T, wrongDSN string)
	}{
		{
			name:      "Store-Adapter (ADR-0047: cdc_capture) mit cdc_admin-Login",
			wrongRole: "cdc_admin",
			loginName: "pgc_test_wrongrole_store",
			action: func(t *testing.T, wrongDSN string) {
				store, err := postgresstorage.New(ctx, wrongDSN)
				if err != nil {
					t.Fatalf("Store-Adapter mit cdc_admin-Login verbinden: %v", err)
				}
				defer store.Close()
				tx, err := model.NewOpenTransaction("tx-wrongrole-store", rolesWiringTestSource)
				if err != nil {
					t.Fatalf("NewOpenTransaction: %v", err)
				}
				position, err := model.NewSourcePosition(rolesWiringTestSource, 100001)
				if err != nil {
					t.Fatalf("NewSourcePosition: %v", err)
				}
				if err := tx.Commit(position, model.NewTimePoint(1)); err != nil {
					t.Fatalf("Commit: %v", err)
				}
				if err := store.PersistTransaction(ctx, tx); !permissionDenied(err) {
					t.Fatalf("Store-Adapter PersistTransaction mit cdc_admin-Login: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
				}
			},
		},
		{
			name:      "Tabellen-Aktivierung (ADR-0047: cdc_admin) mit cdc_capture-Login",
			wrongRole: "cdc_capture",
			loginName: "pgc_test_wrongrole_activation",
			action: func(t *testing.T, wrongDSN string) {
				activation, err := postgresstorage.NewTableActivation(ctx, wrongDSN)
				if err != nil {
					t.Fatalf("Aktivierungs-Adapter mit cdc_capture-Login verbinden: %v", err)
				}
				defer activation.Close()
				table, err := model.NewSourceTable("st-wrongrole-activation", rolesWiringTestSource, "public", "wrongrole_activation")
				if err != nil {
					t.Fatalf("NewSourceTable: %v", err)
				}
				version, err := model.NewSchemaVersion("sv-wrongrole-activation", "st-wrongrole-activation", 1)
				if err != nil {
					t.Fatalf("NewSchemaVersion: %v", err)
				}
				if _, err := activation.Register(ctx, table, version); !permissionDenied(err) {
					t.Fatalf("Tabellen-Aktivierung Register mit cdc_capture-Login: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
				}
			},
		},
		{
			name:      "Heartbeat-Adapter (ADR-0047: cdc_admin) mit cdc_capture-Login",
			wrongRole: "cdc_capture",
			loginName: "pgc_test_wrongrole_heartbeat",
			action: func(t *testing.T, wrongDSN string) {
				heartbeat, err := postgresstorage.NewHeartbeat(ctx, wrongDSN)
				if err != nil {
					t.Fatalf("Heartbeat-Adapter mit cdc_capture-Login verbinden: %v", err)
				}
				defer heartbeat.Close()
				if err := heartbeat.Beat(ctx, rolesWiringTestSource); !permissionDenied(err) {
					t.Fatalf("Heartbeat-Adapter Beat mit cdc_capture-Login: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", err)
				}
			},
		},
		{
			// Ruft die tatsächliche `wiring.go`-Aufrufer-Funktion auf (nicht
			// nur den darunterliegenden Adapter direkt) — sie liest
			// `cfg.AdminDSN` selbst; eine Vertauschung in `wiring.go` (z. B.
			// `cfg.CaptureDSN` statt `cfg.AdminDSN`) würde dieser Test real
			// als Regression zeigen.
			name:      "RegisterConsumer (ADR-0047: cdc_admin) mit cdc_capture-Login als cfg.AdminDSN",
			wrongRole: "cdc_capture",
			loginName: "pgc_test_wrongrole_register",
			action: func(t *testing.T, wrongDSN string) {
				stderr := captureStderr(t, func() {
					if code := bootstrap.RegisterConsumer(ctx, bootstrap.Config{AdminDSN: wrongDSN}, "c-wrongrole-register"); code == 0 {
						t.Fatalf("RegisterConsumer mit cdc_capture-Login als cfg.AdminDSN: erwartet Exit-Code != 0, erhalten 0")
					}
				})
				if !strings.Contains(stderr, "SQLSTATE 42501") {
					t.Fatalf("RegisterConsumer mit cdc_capture-Login als cfg.AdminDSN: erwartet SQLSTATE 42501 (insufficient_privilege) in stderr, erhalten %q", stderr)
				}
			},
		},
		{
			// Dieselbe Begründung wie oben — ruft `bootstrap.AcknowledgeConsumer`
			// direkt auf, nicht nur `postgresstorage.NewConsumerState`.
			name:      "AcknowledgeConsumer (ADR-0047: cdc_admin) mit cdc_capture-Login als cfg.AdminDSN",
			wrongRole: "cdc_capture",
			loginName: "pgc_test_wrongrole_acknowledge",
			action: func(t *testing.T, wrongDSN string) {
				cfg := bootstrap.Config{AdminDSN: wrongDSN, Source: rolesWiringTestSource}
				stderr := captureStderr(t, func() {
					if code := bootstrap.AcknowledgeConsumer(ctx, cfg, "c-wrongrole-acknowledge", 1); code == 0 {
						t.Fatalf("AcknowledgeConsumer mit cdc_capture-Login als cfg.AdminDSN: erwartet Exit-Code != 0, erhalten 0")
					}
				})
				if !strings.Contains(stderr, "SQLSTATE 42501") {
					t.Fatalf("AcknowledgeConsumer mit cdc_capture-Login als cfg.AdminDSN: erwartet SQLSTATE 42501 (insufficient_privilege) in stderr, erhalten %q", stderr)
				}
			},
		},
		{
			name:      "Healthcheck (ADR-0047: cdc_reader) mit cdc_capture-Login",
			wrongRole: "cdc_capture",
			loginName: "pgc_test_wrongrole_healthcheck",
			action: func(t *testing.T, wrongDSN string) {
				stderr := captureStderr(t, func() {
					if code := bootstrap.Healthcheck(ctx, wrongDSN, rolesWiringTestSource); code == 0 {
						t.Fatalf("Healthcheck mit cdc_capture-Login: erwartet Exit-Code != 0, erhalten 0")
					}
				})
				if !strings.Contains(stderr, "nicht lesbar") {
					t.Fatalf("Healthcheck mit cdc_capture-Login: erwartet die Permission-Denied-Meldung (\"nicht lesbar\"), stderr=%q", stderr)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wrongDSN := newTestLoginRole(t, adminPool, baseDSN, tc.loginName, tc.wrongRole)
			tc.action(t, wrongDSN)
		})
	}
}
