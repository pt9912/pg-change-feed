package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Schema-Store-Tests laufen gegen dieselbe reale PostgreSQL-Instanz
// wie die Store-Tests (`make test-store`, `ADR-0030`); ohne DSN
// überspringen sie. Die Tabelle `cdc.table_schema` trägt der
// d-migrate-Rollout (`tools/schema/schema.yaml`, `ADR-0043`, `ADR-0015`
// Folgepflicht) — die handgeschriebene DDL des Store-Adapters
// (`ApplySchema`) trägt sie nicht, dieselbe Abgrenzung wie bei den
// Consumer-State-Tabellen (`consumerstate_test.go`).
//
// Kopplung: diese Datei muss vor `store_test.go` laufen — deren Tests
// räumen das Schema per `DROP SCHEMA cdc CASCADE` samt hand-DDL-Neuaufbau
// zurück (`newTestStore`), der `cdc.table_schema` nicht mitträgt. `go test`
// fährt die Testdateien eines Pakets in Datei-Namensordnung;
// „schemastore_test.go" sortiert vor „sqlviews_test.go"/„store_test.go"
// (`c` < `q`/`t`) und läuft deshalb zuerst — wie bereits
// `consumerstate_test.go` und `sqlviews_test.go` für dieselbe Ausgangslage.
const schemaStoreSource = "src-schema"

// newTestSchemaStore baut den Schema-Store-Adapter gegen die Test-Instanz;
// den Quelle-Bestand trägt er idempotent fort (`ON CONFLICT DO NOTHING`,
// wie `newTestConsumerState`). Anders als `newTestActivation`/
// `newTestStore` räumt dieser Aufbau `cdc.schema_version`/`cdc.source*`
// nicht zurück: andere Tests desselben Pakets (`roles_test.go`) tragen
// bereits `cdc.change`-Zeilen, die per Fremdschlüssel an
// `cdc.schema_version` hängen — ein Rückbau hier verletzte deren Bestand.
// Jeder Testfall dieser Datei trägt seine eigene, eindeutige
// Bindungs-Zeile (`registerSchemaStoreTable`); Cross-Test-Kontamination
// bleibt dadurch aus, ohne einen geteilten Tabellen-Bestand zu berühren.
func newTestSchemaStore(t *testing.T) (*postgresstorage.PostgresSchemaStoreAdapter, *pgxpool.Pool) {
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
		"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'cdc' AND table_name = 'table_schema'",
	).Scan(&tables); err != nil {
		t.Fatalf("table_schema-Tabellen-Prüfung: %v", err)
	}
	if tables != 1 {
		t.Fatalf("cdc.table_schema fehlt — der Schema-Rollout über make schema-rollout trägt sie (ADR-0043, ADR-0015 Folgepflicht); der test-store-Lauf rollt sie vor dem Testlauf aus")
	}

	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Schema-Quelle') ON CONFLICT (source_id) DO NOTHING", schemaStoreSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}

	adapter, err := postgresstorage.NewSchemaStore(ctx, dsn)
	if err != nil {
		t.Fatalf("NewSchemaStore: %v", err)
	}
	t.Cleanup(adapter.Close)
	return adapter, pool
}

// registerSchemaStoreTable trägt eine eigene Bindungs-Zeile für einen
// Testfall ein — die Fremdschlüssel-Vorbedingung von `SchemaStorePort`
// (`SPEC-001`, dieselbe Vorbedingung wie bei `TableActivationPort`); jede
// aufrufende Testfunktion trägt eine eindeutige Tabellen-Kennung, damit
// die Testfälle dieser Datei einander nicht über den geteilten
// `cdc.schema_version`-Bestand hinweg beeinflussen.
func registerSchemaStoreTable(t *testing.T, pool *pgxpool.Pool, id string) model.SourceTableID {
	t.Helper()
	if _, err := pool.Exec(context.Background(),
		"INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', $3)",
		id, schemaStoreSource, "t_"+id,
	); err != nil {
		t.Fatalf("Bindungs-Zeile %s: %v", id, err)
	}
	return model.SourceTableID(id)
}

// TestSchemaStoreCurrentVersionMissing trägt die Abwesenheit vor der
// ersten Registrierung (`ADR-0015` Folgepflicht) — kein Fehler, sondern
// der Zustand „keine über diesen Port registrierte Version".
func TestSchemaStoreCurrentVersionMissing(t *testing.T) {
	store, pool := newTestSchemaStore(t)
	table := registerSchemaStoreTable(t, pool, "tbl-schema-missing")
	ctx := context.Background()

	_, ok, err := store.CurrentVersion(ctx, table)
	if err != nil {
		t.Fatalf("CurrentVersion: %v", err)
	}
	if ok {
		t.Fatalf("CurrentVersion vor Registrierung meldet Bestand")
	}
}

// TestSchemaStoreRegisterAndReadRoundTrip trägt den Round-Trip:
// Registrierung, aktuelle Version, Spaltenform (`SPEC-004`,
// `LH-FA-SCH-005`).
func TestSchemaStoreRegisterAndReadRoundTrip(t *testing.T) {
	store, pool := newTestSchemaStore(t)
	table := registerSchemaStoreTable(t, pool, "tbl-schema-roundtrip")
	ctx := context.Background()

	version, err := model.NewSchemaVersion("sv-schema-1", table, 1)
	if err != nil {
		t.Fatalf("NewSchemaVersion: %v", err)
	}
	schema, err := model.NewTableSchema("sv-schema-1", []model.Column{
		{Name: "id", OID: 23},
		{Name: "name", OID: 25},
	})
	if err != nil {
		t.Fatalf("NewTableSchema: %v", err)
	}

	registered, err := store.RegisterVersion(ctx, version, schema)
	if err != nil {
		t.Fatalf("RegisterVersion: %v", err)
	}
	if !registered {
		t.Fatalf("erste Registrierung meldet bestehenden Stand")
	}

	current, ok, err := store.CurrentVersion(ctx, table)
	if err != nil {
		t.Fatalf("CurrentVersion: %v", err)
	}
	if !ok || current.ID != "sv-schema-1" || current.Version != 1 {
		t.Fatalf("CurrentVersion nach Registrierung: %+v (ok=%v)", current, ok)
	}

	read, err := store.TableSchema(ctx, "sv-schema-1")
	if err != nil {
		t.Fatalf("TableSchema: %v", err)
	}
	if len(read.Columns) != 2 || read.Columns[0].Name != "id" || read.Columns[1].Name != "name" || read.Columns[1].OID != 25 {
		t.Fatalf("gelesene Spaltenform: %+v", read.Columns)
	}

	repeat, err := store.RegisterVersion(ctx, version, schema)
	if err != nil {
		t.Fatalf("RegisterVersion erneut: %v", err)
	}
	if repeat {
		t.Fatalf("zweite Registrierung meldet Neustand (Idempotenz)")
	}
	var columnRows int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM cdc.table_schema WHERE schema_version_id = $1", "sv-schema-1",
	).Scan(&columnRows); err != nil {
		t.Fatalf("Spalten-Zeilen: %v", err)
	}
	if columnRows != 2 {
		t.Fatalf("Spalten-Zeilen nach erneuter Registrierung: %d (Erwartung: 2, keine doppelte Spaltenform)", columnRows)
	}
}

// TestSchemaStoreCurrentVersionHighest trägt die „aktuelle" Version als
// die höchste registrierte, unabhängig von der Registrierungs-Reihenfolge.
func TestSchemaStoreCurrentVersionHighest(t *testing.T) {
	store, pool := newTestSchemaStore(t)
	table := registerSchemaStoreTable(t, pool, "tbl-schema-highest")
	ctx := context.Background()

	for _, v := range []struct {
		id      string
		version int64
	}{
		{"sv-schema-v1", 1},
		{"sv-schema-v3", 3},
		{"sv-schema-v2", 2},
	} {
		version, err := model.NewSchemaVersion(model.SchemaVersionID(v.id), table, v.version)
		if err != nil {
			t.Fatalf("NewSchemaVersion %s: %v", v.id, err)
		}
		schema, err := model.NewTableSchema(model.SchemaVersionID(v.id), []model.Column{{Name: "id", OID: 23}})
		if err != nil {
			t.Fatalf("NewTableSchema %s: %v", v.id, err)
		}
		if _, err := store.RegisterVersion(ctx, version, schema); err != nil {
			t.Fatalf("RegisterVersion %s: %v", v.id, err)
		}
	}

	current, ok, err := store.CurrentVersion(ctx, table)
	if err != nil {
		t.Fatalf("CurrentVersion: %v", err)
	}
	if !ok || current.ID != "sv-schema-v3" || current.Version != 3 {
		t.Fatalf("CurrentVersion = %+v (ok=%v), wollen sv-schema-v3/3", current, ok)
	}
}

// TestSchemaStoreTableSchemaUnknown trägt die Fehlerklasse-Grenze
// (`LH-FA-SCH-004` Negative-Fall): eine nicht registrierte Version endet
// sichtbar über `ErrSchemaVersionUnknown`, keine stille Fehlinterpretation.
func TestSchemaStoreTableSchemaUnknown(t *testing.T) {
	store, _ := newTestSchemaStore(t)
	ctx := context.Background()

	if _, err := store.TableSchema(ctx, "sv-schema-nicht-vorhanden"); !stderrors.Is(err, outbound.ErrSchemaVersionUnknown) {
		t.Fatalf("TableSchema mit unbekannter Version: %v (Erwartung: ErrSchemaVersionUnknown)", err)
	}
}

// TestSchemaStoreRegisterVersionRejectsMismatch trägt die
// Aufrufvertrags-Grenze: `version.ID` und `schema.VersionID` müssen
// übereinstimmen.
func TestSchemaStoreRegisterVersionRejectsMismatch(t *testing.T) {
	store, pool := newTestSchemaStore(t)
	table := registerSchemaStoreTable(t, pool, "tbl-schema-mismatch")
	ctx := context.Background()

	version, err := model.NewSchemaVersion("sv-schema-mismatch", table, 1)
	if err != nil {
		t.Fatalf("NewSchemaVersion: %v", err)
	}
	schema, err := model.NewTableSchema("sv-schema-andere-id", []model.Column{{Name: "id", OID: 23}})
	if err != nil {
		t.Fatalf("NewTableSchema: %v", err)
	}
	if _, err := store.RegisterVersion(ctx, version, schema); !stderrors.Is(err, outbound.ErrSchemaVersionMismatch) {
		t.Fatalf("RegisterVersion mit widersprüchlichen Kennungen: %v (Erwartung: ErrSchemaVersionMismatch)", err)
	}
}
