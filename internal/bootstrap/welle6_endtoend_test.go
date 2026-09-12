package bootstrap_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/bootstrap"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestWelle6ConsumerFullCycleEndToEnd trägt den repo-weiten
// Verifikations-Beleg für `welle-6.md` §3 (Closure-Trigger, Modul 8
// §Rollen-Sequenz für eine Welle, Closure-Schritt 1): ein externer
// Consumer registriert sich, liest Changes, bestätigt eine Position,
// startet neu und setzt an der bestätigten Position fort —
// ausschließlich über den in `slice-021`/`slice-022` verdrahteten
// Zugriffsweg (`bootstrap.RegisterConsumer`/`bootstrap.AcknowledgeConsumer`)
// und den etablierten Lesezugriffsweg über die SQL-Sicht `cdc.changes`
// (`LH-FA-SST-002`, `LH-FA-SST-007` Boundary „bestehender Lesezugriffsweg
// (SQL/API)"; dasselbe Muster wie
// `postgresstorage/sqlviews_test.go::TestChangesViewCarriesRangeLimitAndFilter`).
// Das ist das *Mehr* gegenüber den einzelnen Slice-DoDs: `slice-022`s
// eigener Test (`acknowledge_test.go::TestAcknowledgeConsumerEndToEnd`)
// endet bei der Vorwärts-Invariante, ohne je über `cdc.changes` zu lesen
// oder einen Neustart zu simulieren.
func TestWelle6ConsumerFullCycleEndToEnd(t *testing.T) {
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

	const (
		consumerName = "welle6-e2e-consumer"
		sourceID     = model.SourceID("welle6-e2e-source")
		tableID      = "welle6-e2e-table"
	)

	// Datenstand-Rückbau (Vorbedingung, kein getesteter Zugriffsweg):
	// dieselbe Reihenfolge wie in `acknowledge_test.go`/`register_test.go`
	// — Positions- vor Consumer-Zeile (Fremdschlüssel) — plus Rückbau der
	// in diesem Test geseedeten Change-/Transaktions-/Tabellen-Zeilen.
	for _, stmt := range []struct {
		sql string
		arg any
	}{
		{"DELETE FROM cdc.consumer_position WHERE consumer_id = $1", consumerName},
		{"DELETE FROM cdc.consumer WHERE consumer_id = $1", consumerName},
		{"DELETE FROM cdc.change WHERE source_table_id = $1", tableID},
		{"DELETE FROM cdc.transaction WHERE source_id = $1", string(sourceID)},
		{"DELETE FROM cdc.schema_version WHERE source_table_id = $1", tableID},
		{"DELETE FROM cdc.source_table WHERE source_table_id = $1", tableID},
	} {
		if _, err := pool.Exec(ctx, stmt.sql, stmt.arg); err != nil {
			t.Fatalf("Datenstand-Rückbau (%s): %v", stmt.sql, err)
		}
	}

	// Vorbedingungs-Einrichtung der Quelle und ihrer Changes — dasselbe
	// Seed-Muster wie `sqlviews_test.go`/`consumerstate_test.go`: erlaubtes
	// Direktschreiben, weil es die Quelle/Changes seedet, nicht den
	// Consumer-Zugriffsweg selbst (Negativ-Beleg siehe Dateikommentar
	// unten).
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Welle-6-E2E-Quelle') ON CONFLICT (source_id) DO NOTHING",
		string(sourceID),
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', 'welle6_e2e')",
		tableID, string(sourceID),
	); err != nil {
		t.Fatalf("source_table-Zeile: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ($1, $2, 1)",
		tableID+"-sv1", tableID,
	); err != nil {
		t.Fatalf("schema_version-Zeile: %v", err)
	}

	type seedRow struct {
		transactionID  string
		commitPosition int64
		changeID       string
	}
	seeds := []seedRow{
		{"welle6-e2e-tx-1000", 1000, "welle6-e2e-c-1000"},
		{"welle6-e2e-tx-2000", 2000, "welle6-e2e-c-2000"},
		{"welle6-e2e-tx-3000", 3000, "welle6-e2e-c-3000"},
		{"welle6-e2e-tx-4000", 4000, "welle6-e2e-c-4000"},
	}
	for _, row := range seeds {
		if _, err := pool.Exec(ctx,
			"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ($1, $2, $3)",
			row.transactionID, string(sourceID), row.commitPosition,
		); err != nil {
			t.Fatalf("transaction-Zeile %s: %v", row.transactionID, err)
		}
		if _, err := pool.Exec(ctx,
			`INSERT INTO cdc.change
			    (change_id, transaction_id, source_table_id, sequence, operation, old_data, new_data, schema_version)
			 VALUES ($1, $2, $3, 1, 'INSERT', NULL, '{}'::jsonb, $4)`,
			row.changeID, row.transactionID, tableID, tableID+"-sv1",
		); err != nil {
			t.Fatalf("change-Zeile %s: %v", row.changeID, err)
		}
	}

	cfg := bootstrap.Config{DSN: dsn, Source: sourceID}

	// 1. Registrieren — der externe Zugriffsweg aus `slice-021`.
	if code := bootstrap.RegisterConsumer(ctx, cfg, consumerName); code != 0 {
		t.Fatalf("Registrierung: Exit-Code = %d, wollen 0", code)
	}

	// 2. Lesen — der etablierte Lesezugriffsweg `cdc.changes`
	// (`LH-FA-SST-002`): erster Ausschnitt, zwei der vier geseedeten
	// Changes (Limit 2, LH-FA-REA-003).
	firstBatch := readChangesAbove(t, pool, sourceID, 0, 2)
	if len(firstBatch) != 2 || firstBatch[0].commitPosition != 1000 || firstBatch[1].commitPosition != 2000 {
		t.Fatalf("erster Lese-Ausschnitt = %v, wollen [1000 2000]", firstBatch)
	}
	lastReadPosition := firstBatch[len(firstBatch)-1].commitPosition

	// 3. Bestätigen — der externe Zugriffsweg aus `slice-022`, mit der
	// Position der zuletzt gelesenen Change.
	if code := bootstrap.AcknowledgeConsumer(ctx, cfg, consumerName, uint64(lastReadPosition)); code != 0 {
		t.Fatalf("erste Bestätigung: Exit-Code = %d, wollen 0", code)
	}

	// 4. Neustart simulieren: ein unabhängiger Aufruf fragt die bestätigte
	// Position erneut aus `cdc.consumer_position` ab, statt sie im
	// Test-Code (`lastReadPosition`) weiterzutragen — genau das, was ein
	// neu gestarteter Consumer-Prozess täte.
	restoredPosition := readAcknowledgedPosition(t, pool, consumerName, string(sourceID))
	if restoredPosition != lastReadPosition {
		t.Fatalf("aus der DB zurückgelesene Position = %d, wollen %d (die zuvor bestätigte)", restoredPosition, lastReadPosition)
	}

	// 5. Fortsetzen — Lesen ab der aus der DB zurückgelesenen Position:
	// nur die noch unbestätigten (neueren) Changes, keine Wiederholung der
	// bereits bestätigten.
	secondBatch := readChangesAbove(t, pool, sourceID, restoredPosition, 10)
	if len(secondBatch) != 2 || secondBatch[0].commitPosition != 3000 || secondBatch[1].commitPosition != 4000 {
		t.Fatalf("zweiter Lese-Ausschnitt (fortgesetzt ab %d) = %v, wollen [3000 4000] — keine Wiederholung der bereits bestätigten", restoredPosition, secondBatch)
	}

	// 6. Abschließende Bestätigung der letzten gelesenen Change — der
	// volle Rundlauf endet dort, wo `LH-FA-CON-004.a` es über den
	// Zugriffsweg verlangt.
	finalPosition := secondBatch[len(secondBatch)-1].commitPosition
	if code := bootstrap.AcknowledgeConsumer(ctx, cfg, consumerName, uint64(finalPosition)); code != 0 {
		t.Fatalf("zweite Bestätigung: Exit-Code = %d, wollen 0", code)
	}
	if got := readAcknowledgedPosition(t, pool, consumerName, string(sourceID)); got != 4000 {
		t.Fatalf("gespeicherte Endposition = %d, wollen 4000", got)
	}

	// Negativ-Beleg (welle-6.md §3, „ohne direktes Schreiben der
	// CDC-Speichertabellen"): Registrierung und Bestätigung liefen oben
	// ausschließlich über `bootstrap.RegisterConsumer`/
	// `bootstrap.AcknowledgeConsumer`. Die einzigen Direktzugriffe dieses
	// Tests auf `cdc.consumer`/`cdc.consumer_position` sind der
	// Rückbau-Block oben (Vorbedingung, keine Registrierung/Bestätigung)
	// und die beiden `readAcknowledgedPosition`-Lesezugriffe (SELECT, kein
	// Schreiben) — der Rundlauf selbst schreibt an keiner Stelle direkt in
	// diese beiden Tabellen.
}

type welle6ChangeRow struct {
	changeID       string
	commitPosition int64
}

// readChangesAbove liest über den etablierten Lesezugriffsweg `cdc.changes`
// (`LH-FA-SST-002`, Bereich/Limit/Ordnung wie `LH-FA-REA-001/003/004`) —
// derselbe SQL-Zugriff wie `postgresstorage/sqlviews_test.go`, hier als
// externer Consumer-Lesezugriff auf die von diesem Test geseedeten
// Changes.
func readChangesAbove(t *testing.T, pool *pgxpool.Pool, source model.SourceID, above int64, limit int) []welle6ChangeRow {
	t.Helper()
	rows, err := pool.Query(context.Background(),
		`SELECT change_id, commit_position FROM cdc.changes
		 WHERE source_id = $1 AND commit_position > $2
		 ORDER BY commit_position, transaction_id, sequence
		 LIMIT $3`,
		string(source), above, limit,
	)
	if err != nil {
		t.Fatalf("cdc.changes-Lesen: %v", err)
	}
	defer rows.Close()
	var result []welle6ChangeRow
	for rows.Next() {
		var row welle6ChangeRow
		if err := rows.Scan(&row.changeID, &row.commitPosition); err != nil {
			t.Fatalf("cdc.changes-Scan: %v", err)
		}
		result = append(result, row)
	}
	return result
}

// readAcknowledgedPosition liest die bestätigte Position direkt aus
// `cdc.consumer_position` — der Neustart-Beleg: ein unabhängiger
// Lesezugriff (kein Schreiben), der keinen im-Prozess gehaltenen Zustand
// nutzt, sondern die Position aus der Datenbank zurückliest.
func readAcknowledgedPosition(t *testing.T, pool *pgxpool.Pool, consumer, source string) int64 {
	t.Helper()
	var position int64
	if err := pool.QueryRow(context.Background(),
		"SELECT acknowledged_position FROM cdc.consumer_position WHERE consumer_id = $1 AND source_id = $2",
		consumer, source,
	).Scan(&position); err != nil {
		t.Fatalf("acknowledged_position lesen: %v", err)
	}
	return position
}
