package postgresstorage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Ordnung der offenen Anträge (`ListPending`) und die Ordnung der
// Ableitung des dauerhaften Standes (`ExcludedColumns`,
// `TransformationRules`) sind dieselbe: `requested_at`, bei gleichem
// Zeitstempel `administration_request_id` (`SPEC-019`, `ADR-0127`). Die
// Verarbeitung führt live und beim Prozessstart nur dann zum selben Stand,
// wenn beide Stellen sie in derselben Reihenfolge lesen.

// describeRules gibt einen Regelstand als vergleichbaren Text aus.
func describeRules(rules []model.Transformation) string {
	out := ""
	for _, rule := range rules {
		out += fmt.Sprintf("%s:%s>%s;", rule.Name(), rule.Column(), rule.To())
	}
	return out
}

// pendingOrder liest die offenen Anträge und gibt die Kennungen der Zeilen
// mit der Vorsilbe zurück, in der gelesenen Reihenfolge.
func pendingOrder(t *testing.T, adapter *postgresstorage.AdministrationRequestAdapter, prefix string) ([]model.AdministrationRequestID, []model.AdministrationRequest) {
	t.Helper()
	pending, err := adapter.ListPending(context.Background())
	if err != nil {
		t.Fatalf("ListPending: %v", err)
	}
	var ids []model.AdministrationRequestID
	var requests []model.AdministrationRequest
	for _, request := range pending {
		if len(request.ID) >= len(prefix) && string(request.ID)[:len(prefix)] == prefix {
			ids = append(ids, request.ID)
			requests = append(requests, request)
		}
	}
	return ids, requests
}

// TestAdministrationRequestListPendingOrdersTiesByRequestID trägt die
// Ordnung der Queue gegen die reale PostgreSQL (`LH-FA-CFG-007`,
// `LH-FA-CFG-005`, `SPEC-019`): vier offene Anträge mit demselben
// `requested_at` — ein Remove und ein Set derselben Regel, ein Exclude und
// ein Include derselben Spalte — kommen in der Ordnung der Antrags-Kennung
// aus `ListPending`, unabhängig von der Einfüge-Ordnung, und die Verarbeitung
// in dieser Ordnung führt zum selben Stand wie die Ableitung aus den
// `applied`-Zeilen (Regelstand und Ausschlussstand).
//
// Die Zeilen sind in absteigender Kennungs-Ordnung eingefügt (`b` vor `a`,
// `d` vor `c`): ohne Zweitschlüssel liest die Abfrage die Gleichzeitigen in
// Einfüge-Ordnung. Rot färbende Mutation: `administration_request_id` aus dem
// `ORDER BY` von `SelectPendingAdministrationRequests` streichen — die
// Reihenfolge der Kennungen kehrt sich um, und der Regelstand der Verarbeitung
// (Set, dann Remove: keine Regel) weicht von der Ableitung (Remove, dann Set:
// Regel `r`) ab.
func TestAdministrationRequestListPendingOrdersTiesByRequestID(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	queue, err := postgresstorage.NewAdministrationRequest(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(queue.Close)
	activation, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(activation.Close)

	const (
		ruleTable   = "orders_queue_order_rules"
		columnTable = "orders_queue_order_columns"
		insertRow   = `INSERT INTO cdc.administration_request
    (administration_request_id, source_id, schema_name, table_name, column_name, rule_name, rule_spec, request_kind, requested_at, status)
VALUES ($1, $2, 'public', $3, $4, $5, $6::text::jsonb, $7, $8, $9)`
	)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(),
			"DELETE FROM cdc.administration_request WHERE administration_request_id LIKE 'queue-order-%'")
	})
	insert := func(id, table, kind string, column, rule, spec *string, status string, at time.Time) {
		t.Helper()
		if _, err := pool.Exec(ctx, insertRow, id, administrationRequestSource, table, column, rule, spec, kind, at, status); err != nil {
			t.Fatalf("Antrags-Zeile %q schreiben: %v", id, err)
		}
	}
	text := func(value string) *string { return &value }

	tiedAt := time.Date(2026, 9, 26, 13, 0, 0, 0, time.UTC)
	spec := func(to string) *string {
		return text(fmt.Sprintf(`{"kind": "rename_column", "column": "name", "to": %q}`, to))
	}
	// Der Stand vor dem Paar: die Regel `r` (Ziel `first`) und die
	// Spalte `secret` sind vermerkt.
	insert("queue-order-0-set", ruleTable, "set_transformation", nil, text("r"), spec("first"), "applied", tiedAt.Add(-time.Hour))
	insert("queue-order-0-exclude", columnTable, "exclude_column", text("secret"), nil, nil, "applied", tiedAt.Add(-time.Hour))
	// Die vier gleichzeitigen, offenen Anträge, absteigend eingefügt.
	insert("queue-order-b-set", ruleTable, "set_transformation", nil, text("r"), spec("second"), "pending", tiedAt)
	insert("queue-order-a-remove", ruleTable, "remove_transformation", nil, text("r"), nil, "pending", tiedAt)
	insert("queue-order-d-include", columnTable, "include_column", text("secret"), nil, nil, "pending", tiedAt)
	insert("queue-order-c-exclude", columnTable, "exclude_column", text("secret"), nil, nil, "pending", tiedAt)

	ids, requests := pendingOrder(t, queue, "queue-order-")
	want := []model.AdministrationRequestID{"queue-order-a-remove", "queue-order-b-set", "queue-order-c-exclude", "queue-order-d-include"}
	if fmt.Sprint(ids) != fmt.Sprint(want) {
		t.Fatalf("ListPending-Ordnung der Gleichzeitigen = %v, wollen %v (Antrags-Kennung als Zweitschlüssel)", ids, want)
	}

	// Der Stand der Verarbeitung: die vermerkte Vorgeschichte, dann die offenen
	// Anträge in der gelesenen Ordnung.
	records := []model.TransformationRecord{{Kind: model.AdministrationRequestSetTransformation, Name: "r", Spec: *spec("first")}}
	excluded := map[string]bool{"secret": true}
	for _, request := range requests {
		switch request.Kind {
		case model.AdministrationRequestSetTransformation, model.AdministrationRequestRemoveTransformation:
			records = append(records, model.TransformationRecord{Kind: request.Kind, Name: request.RuleName, Spec: request.RuleSpec})
		case model.AdministrationRequestExcludeColumn:
			excluded[request.Column] = true
		case model.AdministrationRequestIncludeColumn:
			delete(excluded, request.Column)
		}
	}
	live, err := model.FoldTransformations(records)
	if err != nil {
		t.Fatalf("FoldTransformations: %v", err)
	}

	// Die Ableitung: dieselben Anträge vermerkt, Stand aus den `applied`-Zeilen.
	for _, id := range ids {
		if err := queue.MarkApplied(ctx, id); err != nil {
			t.Fatalf("MarkApplied %q: %v", id, err)
		}
	}
	derivedRules, err := activation.TransformationRules(ctx, administrationRequestSource)
	if err != nil {
		t.Fatalf("TransformationRules: %v", err)
	}
	if got, want := describeRules(derivedRules["public."+ruleTable]), describeRules(live); got != want || got != "r:name>second;" {
		t.Fatalf("Regelstand abgeleitet = %q, live = %q, wollen beide %q (Remove vor Set der Kennungs-Ordnung)", got, want, "r:name>second;")
	}
	derivedColumns, err := activation.ExcludedColumns(ctx, administrationRequestSource)
	if err != nil {
		t.Fatalf("ExcludedColumns: %v", err)
	}
	if got := derivedColumns["public."+columnTable]; len(got) != 0 || len(excluded) != 0 {
		t.Fatalf("Ausschlussstand abgeleitet = %v, live = %v, wollen beide leer (Exclude vor Include der Kennungs-Ordnung)", got, excluded)
	}
}

// TestAdministrationRequestSameTransactionCallsKeepCallOrder trägt die
// Zusage von `ADR-0127` gegen die realen SQL-Funktionen (nicht gegen von
// Hand geschriebene Zeilen): sechs Aufrufe einer Transaktion —
// `remove_transformation` und `set_transformation` derselben Regel,
// `exclude_column` und `include_column` derselben Spalte, `disable_table` und
// `enable_table` derselben Tabelle — tragen `requested_at` je Aufruf, und
// `ListPending` liefert sie in der Aufruf-Reihenfolge. Die Verarbeitung in
// dieser Reihenfolge führt zum selben Regel- und Ausschlussstand wie die
// Ableitung aus den vermerkten Zeilen, und dieser Stand ist der der
// Aufrufe: die neue Regel `r` (Ziel `second`), die Spalte `secret` nicht
// ausgeschlossen.
//
// Die Kennung ordnet gleichzeitige Zeilen nach einer zufälligen UUID; ein
// einzelner Durchlauf färbte sich bei falscher Ordnung nur mit
// Wahrscheinlichkeit 1/2 je Paar rot. Der Test läuft deshalb `iterations`
// Transaktionen mit je eigener Tabelle.
//
// Rot färbende Mutation: in der später aufgerufenen Funktion eines Paares
// (`cdc.set_transformation`, `cdc.include_column` oder `cdc.enable_table`
// in `tools/schema/nacharbeit-administration.sql`) `clock_timestamp()` durch
// `now()` ersetzen — ihre Zeile trägt den Transaktionsbeginn und sortiert
// vor die Zeile des früheren Aufrufs. In der früher aufgerufenen Funktion
// bleibt die Reihenfolge richtig und der Test grün. Sind alle sieben
// Funktionen so ersetzt, bricht die Reihenfolge über die zufällige Kennung
// bei etwa jeder zweiten Transaktion.
func TestAdministrationRequestSameTransactionCallsKeepCallOrder(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	queue, err := postgresstorage.NewAdministrationRequest(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(queue.Close)
	activation, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(activation.Close)

	const (
		iterations = 60
		firstSpec  = `{"kind": "rename_column", "column": "name", "to": "first"}`
		secondSpec = `{"kind": "rename_column", "column": "name", "to": "second"}`
	)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE table_name LIKE 'orders_call_order_%'")
	})

	type call struct {
		function string
		query    string
		args     func(table string) []any
	}
	tableArgs := func(table string) []any { return []any{administrationRequestSource, table} }
	calls := []call{
		{"remove_transformation", "SELECT cdc.remove_transformation($1, 'public', $2, 'r')", tableArgs},
		{"set_transformation", "SELECT cdc.set_transformation($1, 'public', $2, 'r', $3::json)", func(table string) []any {
			return []any{administrationRequestSource, table, secondSpec}
		}},
		{"exclude_column", "SELECT cdc.exclude_column($1, 'public', $2, 'secret')", tableArgs},
		{"include_column", "SELECT cdc.include_column($1, 'public', $2, 'secret')", tableArgs},
		{"disable_table", "SELECT cdc.disable_table($1, 'public', $2)", tableArgs},
		{"enable_table", "SELECT cdc.enable_table($1, 'public', $2)", tableArgs},
	}

	for i := 0; i < iterations; i++ {
		table := fmt.Sprintf("orders_call_order_%02d", i)
		if _, err := pool.Exec(ctx, `INSERT INTO cdc.administration_request
    (administration_request_id, source_id, schema_name, table_name, rule_name, rule_spec, request_kind, requested_at, status)
VALUES ($1, $2, 'public', $3, 'r', $4::text::jsonb, 'set_transformation', current_timestamp - interval '1 hour', 'applied')`,
			fmt.Sprintf("call-order-%02d-0-set", i), administrationRequestSource, table, firstSpec); err != nil {
			t.Fatalf("Durchlauf %d: Vorgeschichte schreiben: %v", i, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatalf("Durchlauf %d: Begin: %v", i, err)
		}
		var called []model.AdministrationRequestID
		for _, c := range calls {
			var id string
			if err := tx.QueryRow(ctx, c.query, c.args(table)...).Scan(&id); err != nil {
				_ = tx.Rollback(ctx)
				t.Fatalf("Durchlauf %d: cdc.%s: %v", i, c.function, err)
			}
			called = append(called, model.AdministrationRequestID(id))
		}
		if err := tx.Commit(ctx); err != nil {
			t.Fatalf("Durchlauf %d: Commit: %v", i, err)
		}

		pending, err := queue.ListPending(ctx)
		if err != nil {
			t.Fatalf("Durchlauf %d: ListPending: %v", i, err)
		}
		inCall := make(map[model.AdministrationRequestID]bool, len(called))
		for _, id := range called {
			inCall[id] = true
		}
		var mine []model.AdministrationRequest
		var mineIDs []model.AdministrationRequestID
		for _, request := range pending {
			if inCall[request.ID] {
				mine = append(mine, request)
				mineIDs = append(mineIDs, request.ID)
			}
		}
		if fmt.Sprint(mineIDs) != fmt.Sprint(called) {
			t.Fatalf("Durchlauf %d: ListPending-Ordnung = %v, wollen die Aufruf-Reihenfolge %v", i, mineIDs, called)
		}

		// Der Stand der Verarbeitung: die vermerkte Vorgeschichte, dann die
		// Anträge in der gelesenen Ordnung.
		records := []model.TransformationRecord{{Kind: model.AdministrationRequestSetTransformation, Name: "r", Spec: firstSpec}}
		excluded := map[string]bool{}
		for _, request := range mine {
			switch request.Kind {
			case model.AdministrationRequestSetTransformation, model.AdministrationRequestRemoveTransformation:
				records = append(records, model.TransformationRecord{Kind: request.Kind, Name: request.RuleName, Spec: request.RuleSpec})
			case model.AdministrationRequestExcludeColumn:
				excluded[request.Column] = true
			case model.AdministrationRequestIncludeColumn:
				delete(excluded, request.Column)
			}
		}
		live, err := model.FoldTransformations(records)
		if err != nil {
			t.Fatalf("Durchlauf %d: FoldTransformations: %v", i, err)
		}

		for _, request := range mine {
			if err := queue.MarkApplied(ctx, request.ID); err != nil {
				t.Fatalf("Durchlauf %d: MarkApplied %q: %v", i, request.ID, err)
			}
		}
		derivedRules, err := activation.TransformationRules(ctx, administrationRequestSource)
		if err != nil {
			t.Fatalf("Durchlauf %d: TransformationRules: %v", i, err)
		}
		const wantRules = "r:name>second;"
		if got, liveText := describeRules(derivedRules["public."+table]), describeRules(live); got != liveText || got != wantRules {
			t.Fatalf("Durchlauf %d: Regelstand abgeleitet = %q, live = %q, wollen beide %q", i, got, liveText, wantRules)
		}
		derivedColumns, err := activation.ExcludedColumns(ctx, administrationRequestSource)
		if err != nil {
			t.Fatalf("Durchlauf %d: ExcludedColumns: %v", i, err)
		}
		if got := derivedColumns["public."+table]; len(got) != 0 || len(excluded) != 0 {
			t.Fatalf("Durchlauf %d: Ausschlussstand abgeleitet = %v, live = %v, wollen beide leer", i, got, excluded)
		}
	}
}
