package bootstrap_test

import (
	"context"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/retention"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Der Test dieser Datei fährt den Bereinigungslauf über mehr als zwei echte
// Seiten (`ADR-0124`): der echte Use Case über den echten Store-Adapter und
// den echten Consumer-State-Adapter, beide unter einer Login-Identität der
// Gruppenrolle `cdc_admin` (Verdrahtung wie `wiring.go`, ohne Superuser),
// gegen 25.000 Changes einer Quelle. Er liegt im Paket der Verdrahtung, weil
// ein Adapter-Test keinen Use Case importiert (`.a-check.yml`).

const retentionPagesSource = "src-ret-pages"

// countingStore zählt die Lese-Aufrufe des Stores und hält ihre `limit`.
type countingStore struct {
	outbound.ChangeStorePort
	limits []int
}

func (s *countingStore) ReadRetentionCandidates(ctx context.Context, source model.SourceID, after model.ChangeID, limit int) ([]outbound.RetentionCandidate, error) {
	s.limits = append(s.limits, limit)
	return s.ChangeStorePort.ReadRetentionCandidates(ctx, source, after, limit)
}

type fixedClock struct{ now model.TimePoint }

func (c fixedClock) Now() model.TimePoint { return c.now }

// TestRetentionRunOverManyPagesDeletesWhatAnIndependentCountNames belegt: der
// Lauf über 25.000 Changes in 5.000 Transaktionen (drei nichtleere Seiten
// zu je höchstens `retention.PageSize`) löscht genau die Menge, die eine
// unabhängige SQL-Zählung nach denselben Regeln nennt — Mindestalter 24 h
// gegen die feste Uhr und Consumer-Position 3.500 —, lässt alle übrigen
// Changes stehen und löscht im zweiten Lauf nichts mehr. Die Verteilung der
// Commit-Zeitpunkte macht Alter und Consumer-Position je für eine andere
// Teilmenge bindend. Rot färbende Mutation: `AllowsDeletion` an der
// Schleife des Use Cases durch `true` ersetzen (Alter und Consumer-Position
// binden nichts mehr) oder `after` in der Schleife nicht fortschreiben (der
// Lauf liest die erste Seite erneut).
func TestRetentionRunOverManyPagesDeletesWhatAnIndependentCountNames(t *testing.T) {
	adminPool, baseDSN := adminTestPool(t)
	ctx := context.Background()
	const (
		transactions = 5000
		perTx        = 5
		consumerAt   = 3500
	)
	seedAt := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	now := seedAt.Add(60 * time.Hour)

	for _, statement := range []string{
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Retention-Seiten')",
		"INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1::text || '-tbl', $1, 'public', 'ret_pages')",
		"INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ($1::text || '-sv', $1::text || '-tbl', 1)",
	} {
		if _, err := adminPool.Exec(ctx, statement, retentionPagesSource); err != nil {
			t.Fatalf("Referenz-Zeilen: %v", err)
		}
	}
	t.Cleanup(func() {
		for _, statement := range []string{
			"DELETE FROM cdc.consumer_position WHERE source_id = $1",
			"DELETE FROM cdc.consumer WHERE consumer_id = $1::text || '-consumer'",
			"DELETE FROM cdc.change WHERE transaction_id IN (SELECT transaction_id FROM cdc.transaction WHERE source_id = $1)",
			"DELETE FROM cdc.transaction WHERE source_id = $1",
			"DELETE FROM cdc.schema_version WHERE source_table_id = $1::text || '-tbl'",
			"DELETE FROM cdc.source_table WHERE source_id = $1",
			"DELETE FROM cdc.source WHERE source_id = $1",
		} {
			_, _ = adminPool.Exec(context.Background(), statement, retentionPagesSource)
		}
	})

	// Commit-Zeitpunkt der Transaktion g: seedAt + ((g*7) mod 5000) Minuten.
	if _, err := adminPool.Exec(ctx,
		`INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at)
		 SELECT 'rp-tx-' || lpad(g::text, 5, '0'), $1, g, $2::timestamptz + make_interval(mins => (g * 7) % 5000)
		 FROM generate_series(1, $3::int) AS g`,
		retentionPagesSource, seedAt, transactions); err != nil {
		t.Fatalf("Transaktionen: %v", err)
	}
	if _, err := adminPool.Exec(ctx,
		`INSERT INTO cdc.change (change_id, transaction_id, source_table_id, sequence, operation, old_data, new_data, schema_version)
		 SELECT 'rp-tx-' || lpad(g::text, 5, '0') || '-' || s, 'rp-tx-' || lpad(g::text, 5, '0'), $1::text || '-tbl', s, 'INSERT', NULL, '{}'::jsonb, $1::text || '-sv'
		 FROM generate_series(1, $2::int) AS g, generate_series(1, $3::int) AS s`,
		retentionPagesSource, transactions, perTx); err != nil {
		t.Fatalf("Changes: %v", err)
	}

	count := func(sql string, args ...any) int {
		t.Helper()
		var n int
		if err := adminPool.QueryRow(ctx, sql, args...).Scan(&n); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
		return n
	}
	total := count(`SELECT count(*) FROM cdc.change c JOIN cdc.transaction t ON t.transaction_id = c.transaction_id WHERE t.source_id = $1`, retentionPagesSource)
	if total != transactions*perTx {
		t.Fatalf("Bestand = %d Changes, wollen %d", total, transactions*perTx)
	}
	const (
		sourceChanges  = `SELECT count(*) FROM cdc.change c JOIN cdc.transaction t ON t.transaction_id = c.transaction_id WHERE t.source_id = $1`
		oldEnough      = ` AND t.committed_at <= $2::timestamptz - interval '24 hours'`
		behindConsumer = ` AND t.commit_position <= $3`
		consumerOnly   = ` AND t.commit_position <= $2`
	)
	deletable := sourceChanges + oldEnough + behindConsumer
	want := count(deletable, retentionPagesSource, now, consumerAt)
	byAge := count(sourceChanges+oldEnough, retentionPagesSource, now)
	byConsumer := count(sourceChanges+consumerOnly, retentionPagesSource, consumerAt)
	if want == 0 || want >= byAge || want >= byConsumer {
		t.Fatalf("Verteilung nicht bindend: unabhängige Zählung %d, nur nach Alter %d, nur nach Consumer %d", want, byAge, byConsumer)
	}

	loginDSN := newTestLoginRole(t, adminPool, baseDSN, "pgc_test_retention_pages", "cdc_admin")
	store, err := postgresstorage.New(ctx, loginDSN)
	if err != nil {
		t.Fatalf("Store-Adapter unter cdc_admin-Login: %v", err)
	}
	defer store.Close()
	state, err := postgresstorage.NewConsumerState(ctx, loginDSN)
	if err != nil {
		t.Fatalf("Consumer-State-Adapter unter cdc_admin-Login: %v", err)
	}
	defer state.Close()

	consumer, err := model.NewConsumer(retentionPagesSource+"-consumer", "Seiten-Consumer")
	if err != nil {
		t.Fatalf("NewConsumer: %v", err)
	}
	if _, err := state.Register(ctx, consumer); err != nil {
		t.Fatalf("Register: %v", err)
	}
	acked, err := model.NewSourcePosition(retentionPagesSource, consumerAt)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	if _, err := state.Acknowledge(ctx, model.ConsumerPosition{ConsumerID: consumer.ID, Position: acked}); err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}

	policy, err := model.NewRetentionPolicy(model.Duration{Nanos: int64(24 * time.Hour)})
	if err != nil {
		t.Fatalf("NewRetentionPolicy: %v", err)
	}
	counting := &countingStore{ChangeStorePort: store}
	service := retention.NewRunRetentionService(counting, state, fixedClock{now: model.NewTimePoint(now.UnixNano())})

	result, err := service.Run(ctx, retention.RunRetentionCommand{Source: retentionPagesSource, Policy: policy})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Deleted != want {
		t.Fatalf("Deleted = %d, unabhängige Zählung nennt %d", result.Deleted, want)
	}
	if len(counting.limits) != 4 {
		t.Fatalf("Lese-Aufrufe = %d (%v), wollen 4: drei nichtleere Seiten und die leere Endseite", len(counting.limits), counting.limits)
	}
	for _, limit := range counting.limits {
		if limit != retention.PageSize {
			t.Fatalf("limit = %d, wollen %d", limit, retention.PageSize)
		}
	}
	if remaining := count("SELECT count(*) FROM cdc.change c JOIN cdc.transaction t ON t.transaction_id = c.transaction_id WHERE t.source_id = $1", retentionPagesSource); remaining != total-want {
		t.Fatalf("verbleibende Changes = %d, wollen %d", remaining, total-want)
	}
	if left := count(deletable, retentionPagesSource, now, consumerAt); left != 0 {
		t.Fatalf("%d freigegebene Changes stehen noch", left)
	}

	second, err := service.Run(ctx, retention.RunRetentionCommand{Source: retentionPagesSource, Policy: policy})
	if err != nil {
		t.Fatalf("zweiter Run: %v", err)
	}
	if second.Deleted != 0 {
		t.Fatalf("zweiter Run löscht %d Changes, wollen 0", second.Deleted)
	}
}
