package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Heartbeat-Tests laufen gegen dieselbe reale PostgreSQL-Instanz wie
// die Store- und Consumer-State-Tests (`make test-store`, `ADR-0030`);
// ohne DSN überspringen sie. Die Tabelle `cdc.process_heartbeat` trägt der
// d-migrate-Rollout (`tools/schema/schema.yaml`, `ADR-0043`), die View
// `cdc.heartbeat` die Nacharbeit (`tools/schema/nacharbeit-heartbeat.sql`)
// — beide trägt die handgeschriebene DDL des Store-Adapters nicht.
//
// Kopplung: diese Datei muss vor `store_test.go` laufen — dieselbe
// Ordnungs-Begründung wie bei `consumerstate_test.go`/`sqlviews_test.go`
// (`h` < `s`/`t`).
const heartbeatTestSource = "src-heartbeat"

// newTestHeartbeat baut den Heartbeat-Adapter gegen die Test-Instanz und
// prüft Tabelle und View vor dem ersten Aufruf.
func newTestHeartbeat(t *testing.T) (*postgresstorage.PostgresHeartbeatAdapter, *pgxpool.Pool) {
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

	var tables int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'cdc' AND table_name = 'process_heartbeat'",
	).Scan(&tables); err != nil {
		t.Fatalf("process_heartbeat-Tabellen-Prüfung: %v", err)
	}
	if tables != 1 {
		t.Fatalf("cdc.process_heartbeat fehlt — der Schema-Rollout über make schema-rollout trägt sie (ADR-0043); der test-store-Lauf rollt sie vor dem Testlauf aus")
	}
	var views int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM information_schema.views WHERE table_schema = 'cdc' AND table_name = 'heartbeat'",
	).Scan(&views); err != nil {
		t.Fatalf("heartbeat-View-Prüfung: %v", err)
	}
	if views != 1 {
		t.Fatalf("cdc.heartbeat fehlt — tools/schema/nacharbeit-heartbeat.sql trägt sie im schema-rollout-Lauf")
	}

	if _, err := pool.Exec(ctx, "DELETE FROM cdc.process_heartbeat"); err != nil {
		t.Fatalf("Datenstand-Rückbau: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Heartbeat-Quelle') ON CONFLICT (source_id) DO NOTHING",
		heartbeatTestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}

	adapter, err := postgresstorage.NewHeartbeat(ctx, dsn)
	if err != nil {
		t.Fatalf("NewHeartbeat: %v", err)
	}
	t.Cleanup(adapter.Close)
	return adapter, pool
}

// TestBeatWritesHeartbeatRow belegt den Happy Path (`LH-FA-ADM-002`): der
// erste Schreib-Zug legt genau eine Zeile mit Zeitstempel an.
func TestBeatWritesHeartbeatRow(t *testing.T) {
	adapter, pool := newTestHeartbeat(t)
	ctx := context.Background()

	if err := adapter.Beat(ctx, heartbeatTestSource); err != nil {
		t.Fatalf("Beat: %v", err)
	}

	var rows int
	var heartbeatAt time.Time
	if err := pool.QueryRow(ctx,
		"SELECT count(*), max(heartbeat_at) FROM cdc.process_heartbeat WHERE source_id = $1",
		heartbeatTestSource,
	).Scan(&rows, &heartbeatAt); err != nil {
		t.Fatalf("process_heartbeat-Lesen: %v", err)
	}
	if rows != 1 {
		t.Fatalf("process_heartbeat trägt %d Zeilen für %s, wollen 1", rows, heartbeatTestSource)
	}
	if heartbeatAt.IsZero() {
		t.Fatalf("heartbeat_at ist der Nullwert, wollen einen gesetzten Zeitstempel")
	}
}

// TestBeatIsIdempotentAndAdvancesTimestamp belegt die Fortschreibung
// (UPSERT): ein zweiter Schreib-Zug derselben Quelle trägt keine zweite
// Zeile, der Zeitstempel rückt nicht zurück.
func TestBeatIsIdempotentAndAdvancesTimestamp(t *testing.T) {
	adapter, pool := newTestHeartbeat(t)
	ctx := context.Background()

	if err := adapter.Beat(ctx, heartbeatTestSource); err != nil {
		t.Fatalf("erster Beat: %v", err)
	}
	var first time.Time
	if err := pool.QueryRow(ctx,
		"SELECT heartbeat_at FROM cdc.process_heartbeat WHERE source_id = $1", heartbeatTestSource,
	).Scan(&first); err != nil {
		t.Fatalf("erstes Lesen: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	if err := adapter.Beat(ctx, heartbeatTestSource); err != nil {
		t.Fatalf("zweiter Beat: %v", err)
	}

	var rows int
	var second time.Time
	if err := pool.QueryRow(ctx,
		"SELECT count(*), max(heartbeat_at) FROM cdc.process_heartbeat WHERE source_id = $1",
		heartbeatTestSource,
	).Scan(&rows, &second); err != nil {
		t.Fatalf("zweites Lesen: %v", err)
	}
	if rows != 1 {
		t.Fatalf("process_heartbeat trägt %d Zeilen nach zwei Beat-Aufrufen, wollen weiterhin 1", rows)
	}
	if second.Before(first) {
		t.Fatalf("zweiter Zeitstempel %v liegt vor dem ersten %v", second, first)
	}
}

// TestBeatRejectsEmptySource trägt die Port-Grenze: eine leere Quelle
// erreicht keinen SQL-Aufruf.
func TestBeatRejectsEmptySource(t *testing.T) {
	adapter, _ := newTestHeartbeat(t)

	if err := adapter.Beat(context.Background(), ""); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Beat(\"\") = %v, wollen %v", err, domainerrors.ErrEmptyIdentifier)
	}
}

// TestHeartbeatViewProjectsAge belegt die Lese-Seite (`LH-QA-OPS-002`,
// `ADR-0046` Kategorie C): `cdc.heartbeat` liest Zeitstempel und Alter der
// Quelle, ohne eine Schwellenwert-Entscheidung zu treffen (das bleibt
// Sache des lesenden Systems, wie bei `cdc.metrics`).
func TestHeartbeatViewProjectsAge(t *testing.T) {
	adapter, pool := newTestHeartbeat(t)
	ctx := context.Background()

	if err := adapter.Beat(ctx, heartbeatTestSource); err != nil {
		t.Fatalf("Beat: %v", err)
	}

	var sourceID string
	var heartbeatAt time.Time
	var ageSeconds float64
	if err := pool.QueryRow(ctx,
		"SELECT source_id, heartbeat_at, age_seconds FROM cdc.heartbeat WHERE source_id = $1",
		heartbeatTestSource,
	).Scan(&sourceID, &heartbeatAt, &ageSeconds); err != nil {
		t.Fatalf("cdc.heartbeat-Lesen: %v", err)
	}
	if sourceID != heartbeatTestSource {
		t.Fatalf("cdc.heartbeat trägt source_id %q, wollen %q", sourceID, heartbeatTestSource)
	}
	if ageSeconds < 0 {
		t.Fatalf("age_seconds = %v, wollen einen nichtnegativen Wert direkt nach dem Schreiben", ageSeconds)
	}
	if ageSeconds > 5 {
		t.Fatalf("age_seconds = %v, wollen einen kleinen Wert direkt nach dem Schreiben", ageSeconds)
	}

	var absentRows int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM cdc.heartbeat WHERE source_id = 'src-heartbeat-never-beaten'",
	).Scan(&absentRows); err != nil {
		t.Fatalf("cdc.heartbeat-Lesen (unbekannte Quelle): %v", err)
	}
	if absentRows != 0 {
		t.Fatalf("cdc.heartbeat trägt %d Zeilen für eine nie geschlagene Quelle, wollen 0", absentRows)
	}
}

// TestFaultWritesErrorClass belegt den Happy Path (`LH-FA-ADM-003`,
// `LH-QA-REL-003`): der Fehlerzustand landet in derselben Zeile wie das
// Lebenszeichen, unterscheidbar vom Normalbetrieb (NULL).
func TestFaultWritesErrorClass(t *testing.T) {
	adapter, pool := newTestHeartbeat(t)
	ctx := context.Background()

	if err := adapter.Fault(ctx, heartbeatTestSource, model.ErrorClassStorage); err != nil {
		t.Fatalf("Fault: %v", err)
	}

	var errorClass *string
	if err := pool.QueryRow(ctx,
		"SELECT error_class FROM cdc.process_heartbeat WHERE source_id = $1", heartbeatTestSource,
	).Scan(&errorClass); err != nil {
		t.Fatalf("process_heartbeat-Lesen: %v", err)
	}
	if errorClass == nil || *errorClass != string(model.ErrorClassStorage) {
		t.Fatalf("error_class = %v, wollen %q", errorClass, model.ErrorClassStorage)
	}
}

// TestBeatClearsPriorFault belegt die Boundary (`LH-FA-ADM-003`): endet
// der Fehlerzustand, ist das über den nächsten erfolgreichen Beat
// erkennbar — dieselbe Zeile trägt wieder NULL.
func TestBeatClearsPriorFault(t *testing.T) {
	adapter, pool := newTestHeartbeat(t)
	ctx := context.Background()

	if err := adapter.Fault(ctx, heartbeatTestSource, model.ErrorClassReplication); err != nil {
		t.Fatalf("Fault: %v", err)
	}
	if err := adapter.Beat(ctx, heartbeatTestSource); err != nil {
		t.Fatalf("Beat: %v", err)
	}

	var errorClass *string
	if err := pool.QueryRow(ctx,
		"SELECT error_class FROM cdc.process_heartbeat WHERE source_id = $1", heartbeatTestSource,
	).Scan(&errorClass); err != nil {
		t.Fatalf("process_heartbeat-Lesen: %v", err)
	}
	if errorClass != nil {
		t.Fatalf("error_class = %q nach Beat, wollen NULL (Fehlerzustand endet erkennbar)", *errorClass)
	}
}

// TestFaultRejectsEmptySource trägt dieselbe Port-Grenze wie
// TestBeatRejectsEmptySource: eine leere Quelle erreicht keinen SQL-Aufruf.
func TestFaultRejectsEmptySource(t *testing.T) {
	adapter, _ := newTestHeartbeat(t)

	if err := adapter.Fault(context.Background(), "", model.ErrorClassInternal); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fault(\"\", …) = %v, wollen %v", err, domainerrors.ErrEmptyIdentifier)
	}
}

// TestFaultRejectsUnknownClass trägt dieselbe Invarianten-Grenze wie
// `model.NewErrorClass`: die sieben ADR-0023-Kategorien sind eine
// geschlossene Menge, kein Freitext — auch am Adapter, nicht nur am
// Domänen-Konstruktor.
func TestFaultRejectsUnknownClass(t *testing.T) {
	adapter, _ := newTestHeartbeat(t)

	if err := adapter.Fault(context.Background(), heartbeatTestSource, model.ErrorClass("unbekannt")); !stderrors.Is(err, domainerrors.ErrInvalidErrorClass) {
		t.Fatalf("Fault(…, \"unbekannt\") = %v, wollen %v", err, domainerrors.ErrInvalidErrorClass)
	}
}

// TestHeartbeatViewProjectsErrorClass belegt die Lese-Seite
// (`cdc.heartbeat`, `LH-FA-ADM-003`): die View projiziert den
// Fehlerzustand ungefiltert, NULL bleibt NULL.
func TestHeartbeatViewProjectsErrorClass(t *testing.T) {
	adapter, pool := newTestHeartbeat(t)
	ctx := context.Background()

	if err := adapter.Fault(ctx, heartbeatTestSource, model.ErrorClassPermission); err != nil {
		t.Fatalf("Fault: %v", err)
	}

	var errorClass *string
	if err := pool.QueryRow(ctx,
		"SELECT error_class FROM cdc.heartbeat WHERE source_id = $1", heartbeatTestSource,
	).Scan(&errorClass); err != nil {
		t.Fatalf("cdc.heartbeat-Lesen: %v", err)
	}
	if errorClass == nil || *errorClass != string(model.ErrorClassPermission) {
		t.Fatalf("cdc.heartbeat.error_class = %v, wollen %q", errorClass, model.ErrorClassPermission)
	}

	if err := adapter.Beat(ctx, heartbeatTestSource); err != nil {
		t.Fatalf("Beat: %v", err)
	}
	if err := pool.QueryRow(ctx,
		"SELECT error_class FROM cdc.heartbeat WHERE source_id = $1", heartbeatTestSource,
	).Scan(&errorClass); err != nil {
		t.Fatalf("cdc.heartbeat-Lesen nach Beat: %v", err)
	}
	if errorClass != nil {
		t.Fatalf("cdc.heartbeat.error_class = %q nach Beat, wollen NULL", *errorClass)
	}
}

var _ outbound.HeartbeatPort = (*postgresstorage.PostgresHeartbeatAdapter)(nil)
