package postgresack_test

import (
	"context"
	stderrors "errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresack"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die ACK-Adapter-Tests tragen die Grenzen des Port-Kontrakts: die
// Null-Position-Grenze und das Fehlerklassen-Wrapping des
// Standby-Status-Aufrufs (`SPEC-008`, Klasse `replication`). Die
// Standby-Status-Wirkung am realen Treiber trägt der Verdrahtungs-Test
// in der Composition-Root (`internal/bootstrap`); die reale Verbindung
// kommt aus dem gepinnten Testcontainer (`make test-replication`,
// `ADR-0030`).

// newTestConn baut die Replication-Verbindung der Tests; ohne DSN
// überspringen die realen Tests — die Konstruktions-Grenze trägt der
// verbindungsfreie Test.
func newTestConn(t *testing.T) *pgconn.PgConn {
	t.Helper()
	dsn := os.Getenv("CDC_REPLICATION_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_REPLICATION_TEST_DSN nicht gesetzt — reale ACK-Adapter-Tests laufen über make test-replication")
	}
	ctx := context.Background()
	connConfig, err := pgconn.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("DSN: %v", err)
	}
	connConfig.RuntimeParams["replication"] = "database"
	conn, err := pgconn.ConnectConfig(ctx, connConfig)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(func() { conn.Close(context.Background()) })
	return conn
}

// TestNewRequiresConnection trägt die Konstruktions-Grenze: ein ACK-
// Adapter ohne Replication-Verbindung endet über die Klasse
// `replication` (`SPEC-008`) — kein ACK-Stand ohne Verbindung.
func TestNewRequiresConnection(t *testing.T) {
	_, err := postgresack.New(nil)
	if !stderrors.Is(err, outbound.ErrReplication) {
		t.Fatalf("New(nil): %v", err)
	}
}

// TestAcknowledgeRejectsZeroPosition trägt die Null-Position-Grenze: die
// Bestätigung ohne Position endet als sichtbarer Fehler, nicht still.
func TestAcknowledgeRejectsZeroPosition(t *testing.T) {
	ack, err := postgresack.New(newTestConn(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	err = ack.Acknowledge(context.Background(), model.SourcePosition{})
	if !stderrors.Is(err, outbound.ErrReplication) {
		t.Fatalf("Null-Position: %v", err)
	}
}

// TestAcknowledgeOnClosedConnection trägt das Fehlerklassen-Wrapping am
// Treiber-Fehler: eine geschlossene Verbindung endet über
// `outbound.ErrReplication`; die technische Ursache bleibt über die
// zweite Wrappung lesbar (`ADR-0023`).
func TestAcknowledgeOnClosedConnection(t *testing.T) {
	conn := newTestConn(t)
	ack, err := postgresack.New(conn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	conn.Close(context.Background())

	position, err := model.NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("Position: %v", err)
	}
	err = ack.Acknowledge(context.Background(), position)
	if !stderrors.Is(err, outbound.ErrReplication) {
		t.Fatalf("Geschlossene Verbindung: %v", err)
	}
}
