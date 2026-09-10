package postgresstorage_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Aktivierungs-Tests laufen gegen dieselbe reale PostgreSQL-Instanz
// wie die Store-Tests (`make test-store`, `ADR-0030`); ohne DSN
// überspringen sie. Publication und Bindungs-Zeilen räumt der Test je
// Fall ab — die Daten bleiben im Container.
const activationSource = "src-act-1"

// newTestActivation baut den Aktivierungs-Adapter gegen die Test-Instanz
// und setzt den Referenz-Bestand frisch auf: Schema-Rückbau, ApplySchema
// und die Quelle-Zeile — die Aktivierung trägt die Quelle nicht selbst
// (outbound.TableActivationPort, Vorbedingung `SPEC-001`).
func newTestActivation(t *testing.T) (*postgresstorage.TableActivationAdapter, *pgxpool.Pool, string) {
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
	if _, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS cdc CASCADE"); err != nil {
		t.Fatalf("Schema-Rückbau: %v", err)
	}
	if err := postgresstorage.ApplySchema(ctx, pool); err != nil {
		t.Fatalf("ApplySchema: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Aktivierungs-Quelle')", activationSource,
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	activation, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(activation.Close)
	return activation, pool, dsn
}

// activationContext trägt die Test-Zeitspanne eines Aktivierungs-Laufs.
func activationContext(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// TestRegisterAndRegistered trägt die Idempotenz der Aktivierung
// (`LH-FA-CFG-001` Boundary): der erste Aufruf aktiviert, der zweite
// meldet den bestehenden Stand ohne Wirkung.
func TestRegisterAndRegistered(t *testing.T) {
	activation, pool, _ := newTestActivation(t)
	ctx, cancel := activationContext(t)
	defer cancel()

	registered, err := model.NewSourceTable("tbl-act-1", activationSource, "public", "t_act_1")
	if err != nil {
		t.Fatalf("Tabelle: %v", err)
	}
	version, err := model.NewSchemaVersion("sv-act-1", "tbl-act-1", 1)
	if err != nil {
		t.Fatalf("Schema-Version: %v", err)
	}

	if exists, err := activation.TableExists(ctx, "public", "t_act_1"); err != nil || exists {
		t.Fatalf("TableExists vor dem Anlegen: %v (%v)", exists, err)
	}
	if _, err := pool.Exec(ctx, "CREATE TABLE public.t_act_1 (id int PRIMARY KEY)"); err != nil {
		t.Fatalf("Quell-Tabelle: %v", err)
	}

	newly, err := activation.Register(ctx, registered, version)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if !newly {
		t.Fatalf("erste Aktivierung meldet bestehenden Stand")
	}
	repeat, err := activation.Register(ctx, registered, version)
	if err != nil {
		t.Fatalf("Register erneut: %v", err)
	}
	if repeat {
		t.Fatalf("zweite Aktivierung meldet Neustand")
	}
	if _, ok, err := activation.Registered(ctx, activationSource, "public", "t_act_1"); err != nil || !ok {
		t.Fatalf("Registered nach Aktivierung: %v (%v)", ok, err)
	}
}

// TestUnregister trägt die drei Entzugs-Ausgänge: entfernt ohne
// Change-Bestand, belassen mit Bestand (`LH-FA-CFG-002` Out-of-Scope),
// abwesend ohne Aktivierung.
func TestUnregister(t *testing.T) {
	activation, pool, _ := newTestActivation(t)
	ctx, cancel := activationContext(t)
	defer cancel()

	registered, err := model.NewSourceTable("tbl-act-2", activationSource, "public", "t_act_2")
	if err != nil {
		t.Fatalf("Tabelle: %v", err)
	}
	version, err := model.NewSchemaVersion("sv-act-2", "tbl-act-2", 1)
	if err != nil {
		t.Fatalf("Schema-Version: %v", err)
	}
	if _, err := pool.Exec(ctx, "CREATE TABLE public.t_act_2 (id int PRIMARY KEY)"); err != nil {
		t.Fatalf("Quell-Tabelle: %v", err)
	}
	if _, err := activation.Register(ctx, registered, version); err != nil {
		t.Fatalf("Register: %v", err)
	}

	// Change-Bestand: die persistierten Changes tragen die Bindungs-Zeile
	// als Herkunft — der Entzug belässt sie (Retention trägt ihr
	// Verhalten, nicht die Deaktivierung, `LH-FA-CFG-002` Out-of-Scope).
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ('tx-act-2', $1, 1)",
		activationSource,
	); err != nil {
		t.Fatalf("Transaktion: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.change (change_id, transaction_id, source_table_id, sequence, operation, schema_version) VALUES ('ch-act-2', 'tx-act-2', 'tbl-act-2', 1, 'INSERT', 'sv-act-2')",
	); err != nil {
		t.Fatalf("Change: %v", err)
	}
	removal, err := activation.Unregister(ctx, registered)
	if err != nil {
		t.Fatalf("Unregister mit Bestand: %v", err)
	}
	if removal != outbound.ActivationRetained {
		t.Fatalf("Entzug mit Change-Bestand: %s (Erwartung: belassen)", removal)
	}
	if _, ok, err := activation.Registered(ctx, activationSource, "public", "t_act_2"); err != nil || !ok {
		t.Fatalf("Bindungs-Zeile nach belassenem Entzug: %v (%v)", ok, err)
	}

	// Der Entzug nach dem Bestands-Abbau trägt die Bindungs- und
	// Schema-Version-Zeilen; der Change geht nicht mit ihm — der Test
	// räumt ihn selbst, die Retention trägt ihn am Produktionspfad.
	if _, err := pool.Exec(ctx, "DELETE FROM cdc.change WHERE change_id = 'ch-act-2'"); err != nil {
		t.Fatalf("Bestand räumen: %v", err)
	}
	removal, err = activation.Unregister(ctx, registered)
	if err != nil {
		t.Fatalf("Unregister: %v", err)
	}
	if removal != outbound.ActivationRemoved {
		t.Fatalf("Entzug ohne Bestand: %s (Erwartung: entfernt)", removal)
	}
	if _, ok, err := activation.Registered(ctx, activationSource, "public", "t_act_2"); err != nil || ok {
		t.Fatalf("Bindungs-Zeile nach Entzug: %v (%v)", ok, err)
	}
	var versions int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM cdc.schema_version WHERE source_table_id = $1", "tbl-act-2",
	).Scan(&versions); err != nil {
		t.Fatalf("Schema-Version-Stand: %v", err)
	}
	if versions != 0 {
		t.Fatalf("Schema-Version-Zeilen nach Entzug: %d (Erwartung: 0)", versions)
	}

	removal, err = activation.Unregister(ctx, registered)
	if err != nil {
		t.Fatalf("Unregister erneut: %v", err)
	}
	if removal != outbound.ActivationAbsent {
		t.Fatalf("Entzug ohne Aktivierung: %s (Erwartung: nicht aktiviert)", removal)
	}
}

// TestPublishAndUnpublish trägt die Publication-Verwaltung
// (`LH-FA-CFG-001.a`): der erste Aufruf legt die Publication mit der
// Tabelle an, der zweite bleibt ohne Wirkung (`LH-FA-CFG-001` Boundary),
// der Entzug trägt die Tabelle aus und bleibt ohne Wirkung, wenn sie
// fehlt.
func TestPublishAndUnpublish(t *testing.T) {
	activation, pool, _ := newTestActivation(t)
	ctx, cancel := activationContext(t)
	defer cancel()

	if _, err := pool.Exec(ctx, "CREATE TABLE public.t_act_pub (id int PRIMARY KEY)"); err != nil {
		t.Fatalf("Quell-Tabelle: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_, _ = pool.Exec(cleanupCtx, "DROP PUBLICATION IF EXISTS pub_act_test")
		_, _ = pool.Exec(cleanupCtx, "DROP TABLE IF EXISTS public.t_act_pub")
	})

	if err := activation.Publish(ctx, "pub_act_test", "public", "t_act_pub"); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	var publishActions string
	if err := pool.QueryRow(ctx,
		"SELECT pubinsert::text || pubupdate::text || pubdelete::text FROM pg_publication WHERE pubname = 'pub_act_test'",
	).Scan(&publishActions); err != nil {
		t.Fatalf("Publication-Bestand: %v", err)
	}
	if publishActions != "truetruetrue" {
		t.Fatalf("Publication-Operationen: %s (Erwartung: insert, update, delete)", publishActions)
	}

	if err := activation.Publish(ctx, "pub_act_test", "public", "t_act_pub"); err != nil {
		t.Fatalf("Publish erneut (Idempotenz, LH-FA-CFG-001 Boundary): %v", err)
	}
	var members int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM pg_publication_tables WHERE pubname = 'pub_act_test' AND schemaname = 'public' AND tablename = 't_act_pub'",
	).Scan(&members); err != nil {
		t.Fatalf("Mitgliedschaft: %v", err)
	}
	if members != 1 {
		t.Fatalf("Mitgliedschaften: %d (Erwartung: 1)", members)
	}

	if err := activation.Unpublish(ctx, "pub_act_test", "public", "t_act_pub"); err != nil {
		t.Fatalf("Unpublish: %v", err)
	}
	if err := activation.Unpublish(ctx, "pub_act_test", "public", "t_act_pub"); err != nil {
		t.Fatalf("Unpublish erneut: %v", err)
	}
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM pg_publication_tables WHERE pubname = 'pub_act_test' AND schemaname = 'public' AND tablename = 't_act_pub'",
	).Scan(&members); err != nil {
		t.Fatalf("Mitgliedschaft nach Entzug: %v", err)
	}
	if members != 0 {
		t.Fatalf("Mitgliedschaften nach Entzug: %d (Erwartung: 0)", members)
	}
	if err := activation.Unpublish(ctx, "pub_act_test", "public", "t_act_pub"); err != nil {
		t.Fatalf("Unpublish ohne Mitgliedschaft: %v", err)
	}
}

// TestList trägt die Tabellen-Liste (`LH-FA-CFG-004`): die Quelle trägt
// die aktivierten Tabellen in Schema- und Tabellen-Ordnung, ohne
// Aktivierung liest die Liste leer.
func TestList(t *testing.T) {
	activation, _, _ := newTestActivation(t)
	ctx, cancel := activationContext(t)
	defer cancel()

	tables, err := activation.List(ctx, activationSource)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tables) != 0 {
		t.Fatalf("Liste ohne Aktivierung: %d Tabellen (Erwartung: 0)", len(tables))
	}

	for _, entry := range []struct{ id, schema, table string }{
		{"tbl-act-l2", "public", "t_act_b"},
		{"tbl-act-l1", "public", "t_act_a"},
	} {
		table, err := model.NewSourceTable(model.SourceTableID(entry.id), activationSource, entry.schema, entry.table)
		if err != nil {
			t.Fatalf("Tabelle: %v", err)
		}
		version, err := model.NewSchemaVersion(model.SchemaVersionID("sv-"+entry.id), table.ID, 1)
		if err != nil {
			t.Fatalf("Schema-Version: %v", err)
		}
		if _, err := activation.Register(ctx, table, version); err != nil {
			t.Fatalf("Register %s: %v", entry.id, err)
		}
	}

	tables, err = activation.List(ctx, activationSource)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tables) != 2 {
		t.Fatalf("Tabellen-Liste: %d Tabellen (Erwartung: 2)", len(tables))
	}
	if tables[0].QualifiedName() != "public.t_act_a" || tables[1].QualifiedName() != "public.t_act_b" {
		t.Fatalf("Tabellen-Ordnung: %v, %v", tables[0].QualifiedName(), tables[1].QualifiedName())
	}
}

// TestTableExistsMissing trägt die Abwesenheit der Quelle-Tabelle; die
// Negative-Pfade der Aktivierung, Deaktivierung und Status-Abfrage
// (`LH-FA-CFG-001`/`002`/`003`) lesen darüber.
func TestTableExistsMissing(t *testing.T) {
	activation, _, _ := newTestActivation(t)
	ctx, cancel := activationContext(t)
	defer cancel()

	exists, err := activation.TableExists(ctx, "public", "t_act_nicht_vorhanden")
	if err != nil {
		t.Fatalf("TableExists: %v", err)
	}
	if exists {
		t.Fatalf("fehlende Tabelle meldet Bestand")
	}
}

// TestIdentifierVerweigerung trägt den configuration-Vertrag der
// Aktivierung an der Adapter-Grenze (`SPEC-008`): Bezeichner außerhalb
// des Alphabets enden vor dem ersten SQL-Aufruf über
// `ErrActivationConfiguration` — die Publication bleibt unberührt.
func TestIdentifierVerweigerung(t *testing.T) {
	activation, pool, _ := newTestActivation(t)
	ctx, cancel := activationContext(t)
	defer cancel()

	if _, err := activation.TableExists(ctx, "Public", "t_act_1"); !errors.Is(err, postgresstorage.ErrActivationConfiguration) {
		t.Fatalf("TableExists mit Schema außerhalb des Alphabets: %v (Erwartung: Fehlerklasse configuration)", err)
	}
	if _, err := activation.Published(ctx, "mit-Bindestrich", "public", "t_act_1"); !errors.Is(err, postgresstorage.ErrActivationConfiguration) {
		t.Fatalf("Published mit Publication außerhalb des Alphabets: %v (Erwartung: Fehlerklasse configuration)", err)
	}
	if err := activation.Publish(ctx, "mit-Bindestrich", "public", "t_act_1"); !errors.Is(err, postgresstorage.ErrActivationConfiguration) {
		t.Fatalf("Publish mit Publication außerhalb des Alphabets: %v (Erwartung: Fehlerklasse configuration)", err)
	}
	if err := activation.Publish(ctx, "pub_act_verweigert", "Public", "t_act_1"); !errors.Is(err, postgresstorage.ErrActivationConfiguration) {
		t.Fatalf("Publish mit Schema außerhalb des Alphabets: %v (Erwartung: Fehlerklasse configuration)", err)
	}
	if err := activation.Unpublish(ctx, "pub_act_verweigert", "public", "T_act_1"); !errors.Is(err, postgresstorage.ErrActivationConfiguration) {
		t.Fatalf("Unpublish mit Tabelle außerhalb des Alphabets: %v (Erwartung: Fehlerklasse configuration)", err)
	}

	var publications int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM pg_publication").Scan(&publications); err != nil {
		t.Fatalf("Publication-Bestand: %v", err)
	}
	if publications != 0 {
		t.Fatalf("Verweigerung hinterließ Publicationen: %d (Erwartung: 0)", publications)
	}
}

// TestRetainedStateView trägt den Zustand nach der Deaktivierung mit
// Change-Bestand (`LH-FA-CFG-002` Out-of-Scope): die Bindungs-Zeile
// bleibt als Herkunft der persistierten Changes, die Publication trägt
// die Tabelle nicht mehr — Status- und Listen-Abfragen trennen darüber
// den Erfassungs-Zustand von der Herkunft.
func TestRetainedStateView(t *testing.T) {
	activation, pool, _ := newTestActivation(t)
	ctx, cancel := activationContext(t)
	defer cancel()

	registered, err := model.NewSourceTable("tbl-act-3", activationSource, "public", "t_act_state")
	if err != nil {
		t.Fatalf("Tabelle: %v", err)
	}
	version, err := model.NewSchemaVersion("sv-act-3", "tbl-act-3", 1)
	if err != nil {
		t.Fatalf("Schema-Version: %v", err)
	}
	if _, err := pool.Exec(ctx, "CREATE TABLE public.t_act_state (id int PRIMARY KEY)"); err != nil {
		t.Fatalf("Quell-Tabelle: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_, _ = pool.Exec(cleanupCtx, "DROP PUBLICATION IF EXISTS pub_act_state")
		_, _ = pool.Exec(cleanupCtx, "DROP TABLE IF EXISTS public.t_act_state")
	})
	if _, err := activation.Register(ctx, registered, version); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := activation.Publish(ctx, "pub_act_state", "public", "t_act_state"); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	// Change-Bestand: die persistierten Changes tragen die Bindungs-Zeile
	// als Herkunft.
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position) VALUES ('tx-act-3', $1, 1)",
		activationSource,
	); err != nil {
		t.Fatalf("Transaktion: %v", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.change (change_id, transaction_id, source_table_id, sequence, operation, schema_version) VALUES ('ch-act-3', 'tx-act-3', 'tbl-act-3', 1, 'INSERT', 'sv-act-3')",
	); err != nil {
		t.Fatalf("Change: %v", err)
	}

	if published, err := activation.Published(ctx, "pub_act_state", "public", "t_act_state"); err != nil || !published {
		t.Fatalf("Published vor dem Entzug: %v (%v)", published, err)
	}
	if err := activation.Unpublish(ctx, "pub_act_state", "public", "t_act_state"); err != nil {
		t.Fatalf("Unpublish: %v", err)
	}
	removal, err := activation.Unregister(ctx, registered)
	if err != nil {
		t.Fatalf("Unregister mit Bestand: %v", err)
	}
	if removal != outbound.ActivationRetained {
		t.Fatalf("Entzug mit Change-Bestand: %s (Erwartung: belassen)", removal)
	}

	// Der Zustands-View nach der Deaktivierung: Bindungs-Zeile steht,
	// Mitgliedschaft fehlt — GetStatus liest darüber Retained, ListTables
	// führt die Tabelle in der Retained-Liste.
	if _, ok, err := activation.Registered(ctx, activationSource, "public", "t_act_state"); err != nil || !ok {
		t.Fatalf("Bindungs-Zeile nach der Deaktivierung: %v (%v)", ok, err)
	}
	if published, err := activation.Published(ctx, "pub_act_state", "public", "t_act_state"); err != nil || published {
		t.Fatalf("Mitgliedschaft nach der Deaktivierung: %v (%v)", published, err)
	}
}

// TestPublishedMissingPublication liest die Abwesenheit der
// Mitgliedschaft bei fehlender Publication.
func TestPublishedMissingPublication(t *testing.T) {
	activation, pool, _ := newTestActivation(t)
	ctx, cancel := activationContext(t)
	defer cancel()

	if _, err := pool.Exec(ctx, "CREATE TABLE public.t_act_pub2 (id int PRIMARY KEY)"); err != nil {
		t.Fatalf("Quell-Tabelle: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_, _ = pool.Exec(cleanupCtx, "DROP TABLE IF EXISTS public.t_act_pub2")
	})
	published, err := activation.Published(ctx, "pub_act_fehlt", "public", "t_act_pub2")
	if err != nil {
		t.Fatalf("Published: %v", err)
	}
	if published {
		t.Fatalf("fehlende Publication meldet Mitgliedschaft")
	}
}
