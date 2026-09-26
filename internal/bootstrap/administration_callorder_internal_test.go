package bootstrap

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/removetransformation"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/settransformation"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived
// trägt die Zusage von `ADR-0127` durch die Verarbeitung: ein
// `cdc.remove_transformation` und ein `cdc.set_transformation` derselben Regel
// in **einer** Transaktion (die Folge, die der Fehlertext von K1 nahelegt)
// werden von der Administrations-Goroutine in dieser Folge verarbeitet — das
// Remove nimmt die Regel heraus, das Set trägt sie danach ohne K1-Verstoß ein.
// Die laufende `Assembler`-Bindung liefert danach den Zielnamen der neuen
// Regel, und eine Bindung, die der Prozessstart aus dem abgeleiteten
// Regelstand (`TransformationRules`) neu bildet, liefert dasselbe Bild.
//
// Die Kennung ordnet gleichzeitige Zeilen nach einer zufälligen UUID; ein
// einzelner Durchlauf färbte sich bei falscher Ordnung nur mit
// Wahrscheinlichkeit 1/2 rot. Der Test läuft `iterations` Transaktionen mit
// je eigener Tabelle.
//
// Rot färbende Mutation: in `cdc.set_transformation`
// (`tools/schema/nacharbeit-administration.sql`) `clock_timestamp()` durch
// `now()` ersetzen — das Set trägt den Transaktionsbeginn und läuft vor dem
// Remove, endet an K1 `failed`, und das Remove nimmt die Regel heraus: die
// Bindung liefert die Rohform. Dieselbe Ersetzung in `cdc.remove_transformation`
// lässt die Reihenfolge richtig und den Test grün.
func TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived(t *testing.T) {
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
		iterations = 50
		sourceID   = model.SourceID("src-admin-callorder")
		tablePart  = "admin_callorder_"
		renamed    = `{"kind": "rename_column", "column": "secret", "to": "%s"}`
	)
	if _, err := pool.Exec(ctx,
		"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Aufruf-Reihenfolge der Antrags-Queue') ON CONFLICT (source_id) DO NOTHING",
		string(sourceID),
	); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	bindings := make(map[string]mapper.TableBinding, iterations)
	tables := make([]string, iterations)
	for i := range tables {
		tables[i] = fmt.Sprintf("%s%02d", tablePart, i)
		if _, err := pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS public."+tables[i]+" (id integer PRIMARY KEY, secret text)"); err != nil {
			t.Fatalf("Quell-Tabelle %s anlegen: %v", tables[i], err)
		}
		bindings["public."+tables[i]] = mapper.TableBinding{TableID: model.SourceTableID("tbl-" + tables[i]), SchemaVersion: "sv-callorder"}
	}
	t.Cleanup(func() {
		cleanup := context.Background()
		_, _ = pool.Exec(cleanup, "DELETE FROM cdc.administration_request WHERE source_id = $1", string(sourceID))
		for _, table := range tables {
			_, _ = pool.Exec(cleanup, "DROP TABLE IF EXISTS public."+table)
		}
		_, _ = pool.Exec(cleanup, "DELETE FROM cdc.source WHERE source_id = $1", string(sourceID))
	})

	requests, err := postgresstorage.NewAdministrationRequest(ctx, dsn, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(requests.Close)
	activation, err := postgresstorage.NewTableActivation(ctx, dsn, postgresstorage.WithLog(&recordingLog{}))
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(activation.Close)

	assembler, err := mapper.NewAssembler(sourceID, bindings, nil)
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	log := &recordingLog{}
	deps := administrationDeps{
		requests:              requests,
		transformations:       activation,
		setTransformations:    settransformation.NewSetTransformationService(activation),
		removeTransformations: removetransformation.NewRemoveTransformationService(activation),
		assembler:             assembler,
		source:                sourceID,
		log:                   log,
	}

	var xid uint32
	nextXID := func() uint32 { xid++; return xid }
	for i, table := range tables {
		// Erster Antrag: die Regel `r` mit dem Ziel `first`, einzeln verarbeitet.
		if _, err := pool.Exec(ctx, "SELECT cdc.set_transformation($1, 'public', $2, 'r', $3::json)", string(sourceID), table, fmt.Sprintf(renamed, "first")); err != nil {
			t.Fatalf("Durchlauf %d: cdc.set_transformation: %v", i, err)
		}
		processAdministrationRequests(ctx, deps)
		if got := assemblerRowImage(t, assembler, nextXID(), "public", table); got != `{"id":"1","first":"geheim"}` {
			t.Fatalf("Durchlauf %d: Row Image nach dem ersten Set = %s, erwartet den Zielnamen first — Log: %v", i, got, log.messages)
		}

		// Zweite Runde: Remove und Set derselben Regel in einer Transaktion.
		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatalf("Durchlauf %d: Begin: %v", i, err)
		}
		if _, err := tx.Exec(ctx, "SELECT cdc.remove_transformation($1, 'public', $2, 'r')", string(sourceID), table); err != nil {
			_ = tx.Rollback(ctx)
			t.Fatalf("Durchlauf %d: cdc.remove_transformation: %v", i, err)
		}
		if _, err := tx.Exec(ctx, "SELECT cdc.set_transformation($1, 'public', $2, 'r', $3::json)", string(sourceID), table, fmt.Sprintf(renamed, "second")); err != nil {
			_ = tx.Rollback(ctx)
			t.Fatalf("Durchlauf %d: cdc.set_transformation: %v", i, err)
		}
		if err := tx.Commit(ctx); err != nil {
			t.Fatalf("Durchlauf %d: Commit: %v", i, err)
		}
		processAdministrationRequests(ctx, deps)

		const want = `{"id":"1","second":"geheim"}`
		if got := assemblerRowImage(t, assembler, nextXID(), "public", table); got != want {
			t.Fatalf("Durchlauf %d: Row Image der laufenden Bindung = %s, erwartet %s (Remove vor Set) — Log: %v", i, got, want, log.messages)
		}
	}

	// Der Prozessstart bildet die Bindungen aus dem abgeleiteten Regelstand
	// neu: dasselbe Bild je Tabelle.
	derived, err := activation.TransformationRules(ctx, sourceID)
	if err != nil {
		t.Fatalf("TransformationRules: %v", err)
	}
	rebuilt := make(map[string]mapper.TableBinding, iterations)
	for _, table := range tables {
		qualified := "public." + table
		rebuilt[qualified] = mapper.TableBinding{TableID: model.SourceTableID("tbl-" + table), SchemaVersion: "sv-callorder", Transformations: derived[qualified]}
	}
	restarted, err := mapper.NewAssembler(sourceID, rebuilt, nil)
	if err != nil {
		t.Fatalf("NewAssembler (Neubildung): %v", err)
	}
	for _, table := range tables {
		if got := assemblerRowImage(t, restarted, nextXID(), "public", table); got != `{"id":"1","second":"geheim"}` {
			t.Fatalf("Tabelle %s: Row Image nach der Neubildung aus dem abgeleiteten Regelstand = %s, erwartet die neue Regel", table, got)
		}
	}
}
