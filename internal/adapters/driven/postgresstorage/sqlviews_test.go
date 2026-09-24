package postgresstorage_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// Die SQL-View-Tests belegen den direkten SQL-Zugriff nach
// `LH-FA-SST-002` (Konfiguration, Status, Kernlesezugriffe über SQL) gegen
// die reale PostgreSQL-Instanz (`make test-store`, `ADR-0030`); ohne DSN
// überspringen sie. Die Views (`cdc.active_tables`, `cdc.consumer_status`,
// `cdc.changes`) trägt der d-migrate-Rollout (`tools/schema/schema.yaml`,
// `ADR-0043`) — dieselbe Quelle wie die Consumer-State-Tabellen; die
// handgeschriebene DDL des Store-Adapters (`ApplySchema`) trägt sie nicht.
//
// Kopplung: diese Datei muss vor `store_test.go` laufen — deren Tests
// räumen das Schema per `DROP SCHEMA cdc CASCADE` samt hand-DDL-Neuaufbau
// zurück (`newTestStore`), der die Views nicht mitträgt. `go test` fährt
// die Testdateien eines Pakets in Datei-Namensordnung; „sqlviews_test.go"
// sortiert vor „store_test.go" (`q` < `t`) und läuft deshalb zuerst — wie
// bereits `consumerstate_test.go` (`c` < `s`) für dieselbe Ausgangslage.
const viewsTestSource = "src-views"

// newTestViews prüft die drei Views und liefert den rohen Pool — die
// Views sind reine SQL-Lesezugriffe, kein Go-Adapter-Typ trägt sie.
func newTestViews(t *testing.T) *pgxpool.Pool {
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

	var views int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM information_schema.views WHERE table_schema = 'cdc' AND table_name IN ('active_tables', 'consumer_status', 'changes')",
	).Scan(&views); err != nil {
		t.Fatalf("View-Prüfung: %v", err)
	}
	if views != 3 {
		t.Fatalf("SQL-Views unvollständig (%d von 3) — der Schema-Rollout über make schema-rollout trägt sie (ADR-0043); der test-store-Lauf rollt sie vor dem Testlauf aus", views)
	}

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Views-Quelle') ON CONFLICT (source_id) DO NOTHING",
		viewsTestSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	return pool
}

// TestActiveTablesViewCarriesCurrentSchemaVersion belegt die
// Konfigurations-Sicht (`LH-FA-SST-002`): die gebundene Tabelle liest sich
// mit ihrer höchsten Schema-Version, unabhängig von der Anlage-Reihenfolge
// der Versionen.
func TestActiveTablesViewCarriesCurrentSchemaVersion(t *testing.T) {
	pool := newTestViews(t)
	ctx := context.Background()
	const tableID = "vt-active-1"

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', 'orders')",
		tableID, viewsTestSource,
	); err != nil {
		t.Fatalf("source_table-Zeile: %v", err)
	}
	for _, row := range []struct {
		id      string
		version int
	}{
		{"vt-active-1-sv2", 2},
		{"vt-active-1-sv1", 1},
	} {
		if _, err := pool.Exec(ctx,
			"INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ($1, $2, $3)",
			row.id, tableID, row.version,
		); err != nil {
			t.Fatalf("schema_version-Zeile %s: %v", row.id, err)
		}
	}

	var schemaName, name string
	var version int64
	if err := pool.QueryRow(ctx,
		"SELECT schema_name, table_name, current_schema_version FROM cdc.active_tables WHERE source_table_id = $1",
		tableID,
	).Scan(&schemaName, &name, &version); err != nil {
		t.Fatalf("active_tables-Lesen: %v", err)
	}
	if schemaName != "public" || name != "orders" {
		t.Fatalf("active_tables trägt %s.%s, wollen public.orders", schemaName, name)
	}
	if version != 2 {
		t.Fatalf("active_tables trägt Version %d, wollen die höchste (2)", version)
	}
}

// TestConsumerStatusViewCarriesBacklog belegt die Status-Sicht
// (`LH-FA-SST-002`, `LH-FA-ADM-005`): ein bestätigender Consumer trägt
// seinen Rückstand als Differenz der beiden Positionsspalten, ein
// unbestätigter Consumer trägt die Positionsspalten als NULL statt eines
// Fehlers.
func TestConsumerStatusViewCarriesBacklog(t *testing.T) {
	pool := newTestViews(t)
	ctx := context.Background()
	const (
		consumerAcked   = "vc-status-acked"
		consumerUnacked = "vc-status-unacked"
		transactionID   = "vt-status-tx"
	)

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ($1, $2, 500)",
		transactionID, viewsTestSource,
	); err != nil {
		t.Fatalf("transaction-Zeile: %v", err)
	}
	for _, id := range []string{consumerAcked, consumerUnacked} {
		if _, err := pool.Exec(ctx,
			"INSERT INTO cdc.consumer (consumer_id, name) VALUES ($1, $2)",
			id, "Consumer "+id,
		); err != nil {
			t.Fatalf("consumer-Zeile %s: %v", id, err)
		}
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.consumer_position (consumer_id, source_id, acknowledged_position) VALUES ($1, $2, 300)",
		consumerAcked, viewsTestSource,
	); err != nil {
		t.Fatalf("consumer_position-Zeile: %v", err)
	}

	var acknowledged, latest int64
	if err := pool.QueryRow(ctx,
		"SELECT acknowledged_position, latest_commit_position FROM cdc.consumer_status WHERE consumer_id = $1",
		consumerAcked,
	).Scan(&acknowledged, &latest); err != nil {
		t.Fatalf("consumer_status-Lesen (acked): %v", err)
	}
	if acknowledged != 300 || latest != 500 {
		t.Fatalf("consumer_status trägt (%d, %d), wollen (300, 500)", acknowledged, latest)
	}
	if backlog := latest - acknowledged; backlog != 200 {
		t.Fatalf("Rückstand = %d, wollen 200", backlog)
	}

	var acknowledgedNull, latestNull *int64
	if err := pool.QueryRow(ctx,
		"SELECT acknowledged_position, latest_commit_position FROM cdc.consumer_status WHERE consumer_id = $1",
		consumerUnacked,
	).Scan(&acknowledgedNull, &latestNull); err != nil {
		t.Fatalf("consumer_status-Lesen (unacked): %v", err)
	}
	if acknowledgedNull != nil || latestNull != nil {
		t.Fatalf("consumer_status trägt für einen unbestätigten Consumer (%v, %v), wollen NULL/NULL", acknowledgedNull, latestNull)
	}
}

// TestChangesViewCarriesRangeLimitAndFilter belegt den Kernlesezugriff
// (`LH-FA-SST-002`): dieselben Bereichs- (`LH-FA-REA-001`), Limit-
// (`LH-FA-REA-003`), Ordnungs- (`LH-FA-REA-004`) und Filter-Zusagen
// (`LH-FA-REA-006`) wie am Go-Port, hier über WHERE/LIMIT der SQL-Sicht
// statt über `ChangeQuery`.
func TestChangesViewCarriesRangeLimitAndFilter(t *testing.T) {
	pool := newTestViews(t)
	ctx := context.Background()
	const (
		tableMain  = "vt-changes-main"
		tableOther = "vt-changes-other"
	)

	for _, table := range []string{tableMain, tableOther} {
		if _, err := pool.Exec(ctx,
			"INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', $3)",
			table, viewsTestSource, table,
		); err != nil {
			t.Fatalf("source_table-Zeile %s: %v", table, err)
		}
		if _, err := pool.Exec(ctx,
			"INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ($1, $2, 1)",
			table+"-sv", table,
		); err != nil {
			t.Fatalf("schema_version-Zeile %s: %v", table, err)
		}
	}

	type seedRow struct {
		transactionID  string
		commitPosition int
		changeID       string
		table          string
	}
	rows := []seedRow{
		{"vt-changes-tx-100", 1100, "vt-changes-c-100", tableMain},
		{"vt-changes-tx-200", 1200, "vt-changes-c-200", tableOther},
		{"vt-changes-tx-300", 1300, "vt-changes-c-300", tableMain},
	}
	for _, row := range rows {
		if _, err := pool.Exec(ctx,
			"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ($1, $2, $3)",
			row.transactionID, viewsTestSource, row.commitPosition,
		); err != nil {
			t.Fatalf("transaction-Zeile %s: %v", row.transactionID, err)
		}
		if _, err := pool.Exec(ctx,
			`INSERT INTO cdc.change
			    (change_id, transaction_id, source_table_id, sequence, operation, old_data, new_data, schema_version)
			 VALUES ($1, $2, $3, 1, 'INSERT', NULL, '{}'::jsonb, $4)`,
			row.changeID, row.transactionID, row.table, row.table+"-sv",
		); err != nil {
			t.Fatalf("change-Zeile %s: %v", row.changeID, err)
		}
	}

	// Bereich [1100, 1300) mit Limit 1 (LH-FA-REA-001, LH-FA-REA-003):
	// trägt genau den ersten Change in Commit-Position-Ordnung.
	rangeRows, err := pool.Query(ctx,
		`SELECT change_id FROM cdc.changes
		 WHERE source_id = $1 AND commit_position >= 1100 AND commit_position < 1300
		 ORDER BY commit_position, transaction_id, sequence
		 LIMIT 1`,
		viewsTestSource,
	)
	if err != nil {
		t.Fatalf("Bereichs-Lesen: %v", err)
	}
	var rangeIDs []string
	for rangeRows.Next() {
		var id string
		if err := rangeRows.Scan(&id); err != nil {
			t.Fatalf("Bereichs-Scan: %v", err)
		}
		rangeIDs = append(rangeIDs, id)
	}
	rangeRows.Close()
	if len(rangeIDs) != 1 || rangeIDs[0] != "vt-changes-c-100" {
		t.Fatalf("Bereich [1100,1300) Limit 1 = %v, wollen genau [vt-changes-c-100]", rangeIDs)
	}

	// Tabellenfilter (LH-FA-REA-006): nur Changes der gefilterten Tabelle.
	filterRows, err := pool.Query(ctx,
		"SELECT change_id FROM cdc.changes WHERE source_id = $1 AND table_name = $2 ORDER BY commit_position",
		viewsTestSource, tableMain,
	)
	if err != nil {
		t.Fatalf("Filter-Lesen: %v", err)
	}
	var filterIDs []string
	for filterRows.Next() {
		var id string
		if err := filterRows.Scan(&id); err != nil {
			t.Fatalf("Filter-Scan: %v", err)
		}
		filterIDs = append(filterIDs, id)
	}
	filterRows.Close()
	if len(filterIDs) != 2 || filterIDs[0] != "vt-changes-c-100" || filterIDs[1] != "vt-changes-c-300" {
		t.Fatalf("Filter %s = %v, wollen [vt-changes-c-100 vt-changes-c-300] in Commit-Ordnung", tableMain, filterIDs)
	}
}

// TestChangesViewCarriesOriginLikeReadChanges belegt das Feld `origin` an
// der View (`SPEC-002`, `LH-FA-CAP-009`, `LH-FA-DAT-006` Boundary) und die
// Spaltenmenge-Parität von View und `ReadChanges`
// (`BEO-PGC/lese-doppelquelle`): `origin` ist die **letzte** Spalte von
// `cdc.changes`; eine gespeicherte Zeile mit `wal`, mit `backfill` und eine
// Zeile ohne Wert (`NULL`) lesen über die View und über `ReadChanges` mit
// derselben Herkunft — die Zeile ohne Wert als `wal`.
func TestChangesViewCarriesOriginLikeReadChanges(t *testing.T) {
	pool := newTestViews(t)
	ctx := context.Background()
	const table = "vt-origin-table"

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', $1)",
		table, viewsTestSource,
	); err != nil {
		t.Fatalf("source_table-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ($1, $2, 1)",
		table+"-sv", table,
	); err != nil {
		t.Fatalf("schema_version-Zeile: %v", err)
	}

	type seedRow struct {
		changeID string
		position int
		origin   any
		want     string
	}
	seeds := []seedRow{
		{"vt-origin-c-wal", 2100, "wal", "wal"},
		{"vt-origin-c-backfill", 2200, "backfill", "backfill"},
		{"vt-origin-c-null", 2300, nil, "wal"},
	}
	for _, seed := range seeds {
		txID := seed.changeID + "-tx"
		if _, err := pool.Exec(ctx,
			"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ($1, $2, $3)",
			txID, viewsTestSource, seed.position,
		); err != nil {
			t.Fatalf("transaction-Zeile %s: %v", txID, err)
		}
		if _, err := pool.Exec(ctx,
			`INSERT INTO cdc.change
			    (change_id, transaction_id, source_table_id, sequence, operation, old_data, new_data, schema_version, origin)
			 VALUES ($1, $2, $3, 1, 'INSERT', NULL, '{}'::jsonb, $4, $5)`,
			seed.changeID, txID, table, table+"-sv", seed.origin,
		); err != nil {
			t.Fatalf("change-Zeile %s: %v", seed.changeID, err)
		}
	}

	// `origin` steht als letzte Spalte der View.
	var lastColumn string
	if err := pool.QueryRow(ctx,
		`SELECT column_name FROM information_schema.columns
		 WHERE table_schema = 'cdc' AND table_name = 'changes'
		 ORDER BY ordinal_position DESC LIMIT 1`,
	).Scan(&lastColumn); err != nil {
		t.Fatalf("Spalten-Prüfung: %v", err)
	}
	if lastColumn != "origin" {
		t.Fatalf("letzte Spalte von cdc.changes = %q, wollen origin", lastColumn)
	}

	// Über die View gelesen: NULL liest als wal.
	viewOrigins := map[string]string{}
	viewRows, err := pool.Query(ctx,
		"SELECT change_id, origin FROM cdc.changes WHERE source_id = $1 AND table_name = $2",
		viewsTestSource, table,
	)
	if err != nil {
		t.Fatalf("View-Lesen: %v", err)
	}
	// Die Verbindung geht auch bei einem t.Fatalf in der Schleife zurück in
	// den Pool — sonst blockiert pool.Close() im Cleanup den Testlauf.
	defer viewRows.Close()
	for viewRows.Next() {
		var id string
		var origin *string
		if err := viewRows.Scan(&id, &origin); err != nil {
			t.Fatalf("View-Scan: %v", err)
		}
		if origin == nil {
			t.Fatalf("%s: die View liefert NULL statt wal", id)
		}
		viewOrigins[id] = *origin
	}
	viewRows.Close()

	// Über ReadChanges gelesen: dieselbe Herkunft je Change.
	store, err := postgresstorage.New(ctx, os.Getenv("CDC_STORE_TEST_DSN"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(store.Close)
	records, err := store.ReadChanges(ctx, outbound.ChangeQuery{Source: viewsTestSource, Table: table})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if len(records) != len(seeds) {
		t.Fatalf("ReadChanges = %d Records, wollen %d", len(records), len(seeds))
	}
	for _, record := range records {
		id := string(record.Change.ID)
		if viewOrigins[id] != string(record.Change.Origin) {
			t.Fatalf("%s: origin View=%q ReadChanges=%q", id, viewOrigins[id], record.Change.Origin)
		}
	}
	for _, seed := range seeds {
		if viewOrigins[seed.changeID] != seed.want {
			t.Fatalf("%s: origin über die View = %q, wollen %q", seed.changeID, viewOrigins[seed.changeID], seed.want)
		}
	}
}
