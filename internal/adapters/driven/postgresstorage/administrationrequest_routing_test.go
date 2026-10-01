package postgresstorage_test

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Tests der Routing-Antragsarten (`LH-FA-CFG-008`) laufen gegen dieselbe
// reale PostgreSQL-Instanz wie die übrigen Antrags-Queue-Tests (`make
// test-store`); ohne DSN überspringen sie. Die Dateinamen-Ordnung
// (`administrationrequest_` vor `consumerstate_`) trägt dieselbe Kopplung wie
// in `administrationrequest_test.go`.

func describeRoutes(rules []model.RouteRule) string {
	out := ""
	for _, rule := range rules {
		condition, conditional := rule.Condition()
		when := "-"
		if conditional {
			when = condition.Column + "=" + condition.Equals
		}
		out += fmt.Sprintf("%s:%s@%d[%s];", rule.Name(), rule.Target(), rule.Order(), when)
	}
	return out
}

// TestAdministrationRequestRoutingRequestsCarryRuleAndNotify trägt den realen
// Antrags-Weg der Routing-Antragsarten: `cdc.set_route` schreibt eine
// `pending`-Zeile mit Regelname und Regelform (Spalte `jsonb`) und ohne
// Spalte, `cdc.remove_route` eine mit Regelname und ohne Regelform; beide
// senden `pg_notify` mit der Antrags-ID, und der Adapter liest beide Felder
// über `ListPending` zurück. Rot färbende Mutationen (Eingabeseite): die Art
// `'set_route'` in `cdc.set_route` durch `'remove_route'` ersetzen — die Zeile
// trägt die falsche Art; `p_rule_spec::jsonb` in `cdc.set_route` durch
// `NULL::jsonb` ersetzen — die Regelform der Zeile ist leer; `p_rule_name`
// durch `NULL` ersetzen — der Regelname der Zeile ist leer.
func TestAdministrationRequestRoutingRequestsCarryRuleAndNotify(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	var requestIDs []string
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE administration_request_id = ANY($1)", requestIDs)
	})
	adapter, err := postgresstorage.NewAdministrationRequest(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(adapter.Close)
	listener := listenForAdministrationNotify(t, ctx, dsn)

	const ruleSpec = `{"target": "eu", "order": 10, "when": {"column": "region", "equals": "eu"}}`
	var setID, removeID, nullSpecID string
	if err := pool.QueryRow(ctx,
		"SELECT cdc.set_route($1, $2, $3, $4, $5::json)", administrationRequestSource, "public", "orders_routes", "eu_orders", ruleSpec,
	).Scan(&setID); err != nil {
		t.Fatalf("cdc.set_route: %v", err)
	}
	requestIDs = append(requestIDs, setID)
	notifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	notification, err := listener.WaitForNotification(notifyCtx)
	if err != nil {
		t.Fatalf("WaitForNotification: %v", err)
	}
	if notification.Channel != "cdc_administration" || notification.Payload != setID {
		t.Fatalf("Notify = %q/%q, wollen cdc_administration/%q", notification.Channel, notification.Payload, setID)
	}
	if err := pool.QueryRow(ctx,
		"SELECT cdc.remove_route($1, $2, $3, $4)", administrationRequestSource, "public", "orders_routes", "eu_orders",
	).Scan(&removeID); err != nil {
		t.Fatalf("cdc.remove_route: %v", err)
	}
	requestIDs = append(requestIDs, removeID)
	notification, err = listener.WaitForNotification(notifyCtx)
	if err != nil {
		t.Fatalf("WaitForNotification: %v", err)
	}
	if notification.Payload != removeID {
		t.Fatalf("Notify-Payload = %q, wollen die Antrags-ID %q", notification.Payload, removeID)
	}
	// Ein JSON-`null` ist eine Regelform, kein SQL-NULL: die Spalte trägt den
	// Text `null`.
	if err := pool.QueryRow(ctx,
		"SELECT cdc.set_route($1, $2, $3, $4, 'null'::json)", administrationRequestSource, "public", "orders_routes", "json_null",
	).Scan(&nullSpecID); err != nil {
		t.Fatalf("cdc.set_route(JSON null): %v", err)
	}
	requestIDs = append(requestIDs, nullSpecID)

	for _, tc := range []struct {
		id, kind string
		wantSpec bool
	}{{setID, "set_route", true}, {removeID, "remove_route", false}} {
		var kind, status, ruleName string
		var column, spec, message *string
		if err := pool.QueryRow(ctx,
			"SELECT request_kind, status, rule_name, column_name, rule_spec::text, error_message FROM cdc.administration_request WHERE administration_request_id = $1", tc.id,
		).Scan(&kind, &status, &ruleName, &column, &spec, &message); err != nil {
			t.Fatalf("Antrags-Zeile %s lesen: %v", tc.kind, err)
		}
		if kind != tc.kind || status != "pending" || ruleName != "eu_orders" || column != nil || message != nil || (spec != nil) != tc.wantSpec {
			t.Fatalf("Antrags-Zeile %s: Art %q, Status %q, Regelname %q, Spalte gesetzt %t, Regelform gesetzt %t (erwartet %t), Fehlertext gesetzt %t",
				tc.kind, kind, status, ruleName, column != nil, spec != nil, tc.wantSpec, message != nil)
		}
	}

	pending, err := adapter.ListPending(ctx)
	if err != nil {
		t.Fatalf("ListPending: %v", err)
	}
	byID := map[model.AdministrationRequestID]model.AdministrationRequest{}
	for _, row := range pending {
		if row.Rejected == nil {
			byID[row.Request.ID] = row.Request
		}
	}
	setRequest, found := byID[model.AdministrationRequestID(setID)]
	if !found {
		t.Fatalf("ListPending trägt nicht den set_route-Antrag %q: %+v", setID, pending)
	}
	if setRequest.Kind != model.AdministrationRequestSetRoute || setRequest.RuleName != "eu_orders" || setRequest.Column != "" || setRequest.Table != "orders_routes" {
		t.Fatalf("set_route-Antrag: %+v", setRequest)
	}
	var decoded struct {
		Target string `json:"target"`
		Order  int    `json:"order"`
		When   struct {
			Column string `json:"column"`
			Equals string `json:"equals"`
		} `json:"when"`
	}
	if err := json.Unmarshal([]byte(setRequest.RuleSpec), &decoded); err != nil {
		t.Fatalf("Regelform %q ist kein JSON-Objekt: %v", setRequest.RuleSpec, err)
	}
	if decoded.Target != "eu" || decoded.Order != 10 || decoded.When.Column != "region" || decoded.When.Equals != "eu" {
		t.Fatalf("Regelform = %+v, wollen target/order/when wie beantragt", decoded)
	}
	removeRequest, found := byID[model.AdministrationRequestID(removeID)]
	if !found {
		t.Fatalf("ListPending trägt nicht den remove_route-Antrag %q: %+v", removeID, pending)
	}
	if removeRequest.Kind != model.AdministrationRequestRemoveRoute || removeRequest.RuleName != "eu_orders" || removeRequest.RuleSpec != "" {
		t.Fatalf("remove_route-Antrag: %+v", removeRequest)
	}
	if nullRequest := byID[model.AdministrationRequestID(nullSpecID)]; nullRequest.RuleSpec != "null" {
		t.Fatalf("Regelform des JSON-null-Antrags = %q, wollen den Text null", nullRequest.RuleSpec)
	}
}

// TestAdministrationRequestListPendingCarriesRouteRowsWithMissingFields trägt:
// ein `set_route` mit SQL-NULL als Regelname, mit SQL-NULL oder JSON-`null` als
// Regelform und ein `remove_route` mit SQL-NULL als Regelname entstehen als
// `pending`-Antrag (die Funktionen prüfen nichts) und kommen über `ListPending`
// als Antrag zurück, nicht als Lesefehler — SQL- und JSON-`null` bleiben am
// leeren Text bzw. am Text `null` unterscheidbar; die Zeile dahinter bleibt
// lesbar. Rot färbende Mutation (Eingabeseite): die Prüfung `ruleName == ""` in
// `model.NewAdministrationRequest` für `set_route` zurücklegen — die Zeilen mit
// leerem Regelnamen stehen als `Rejected` statt als Antrag im Ergebnis.
func TestAdministrationRequestListPendingCarriesRouteRowsWithMissingFields(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	adapter, err := postgresstorage.NewAdministrationRequest(ctx, dsn)
	if err != nil {
		t.Fatalf("NewAdministrationRequest: %v", err)
	}
	t.Cleanup(adapter.Close)

	const table = "orders_route_pending_rows"
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE table_name = $1", table)
	})
	create := func(call string, args ...any) model.AdministrationRequestID {
		t.Helper()
		var id string
		if err := pool.QueryRow(ctx, call, args...).Scan(&id); err != nil {
			t.Fatalf("%s: %v", call, err)
		}
		return model.AdministrationRequestID(id)
	}
	nullName := create("SELECT cdc.set_route($1, 'public', $2, NULL, '{\"target\": \"x\", \"order\": 1}'::json)", administrationRequestSource, table)
	nullSpec := create("SELECT cdc.set_route($1, 'public', $2, 'regel_a', NULL::json)", administrationRequestSource, table)
	jsonNull := create("SELECT cdc.set_route($1, 'public', $2, 'regel_b', 'null'::json)", administrationRequestSource, table)
	removeNull := create("SELECT cdc.remove_route($1, 'public', $2, NULL)", administrationRequestSource, table)
	valid := create("SELECT cdc.set_route($1, 'public', $2, 'regel_c', '{\"target\": \"x\", \"order\": 2}'::json)", administrationRequestSource, table)

	pending, err := adapter.ListPending(ctx)
	if err != nil {
		t.Fatalf("ListPending = %v, wollen nil (die Zeilen sind Anträge, keine Lesefehler)", err)
	}
	byID := map[model.AdministrationRequestID]model.AdministrationRequest{}
	for _, row := range pending {
		if row.Rejected == nil {
			byID[row.Request.ID] = row.Request
		}
	}
	for id, want := range map[model.AdministrationRequestID][2]string{
		nullName:   {"", `{"order": 1, "target": "x"}`},
		nullSpec:   {"regel_a", ""},
		jsonNull:   {"regel_b", "null"},
		removeNull: {"", ""},
	} {
		request, found := byID[id]
		if !found {
			t.Fatalf("ListPending trägt den Antrag %q nicht: %+v", id, pending)
		}
		if request.RuleName != want[0] || request.RuleSpec != want[1] {
			t.Fatalf("Antrag %q = Regelname %q, Regelform %q, wollen %q und %q", id, request.RuleName, request.RuleSpec, want[0], want[1])
		}
	}
	if request, found := byID[valid]; !found || request.RuleName != "regel_c" || request.Kind != model.AdministrationRequestSetRoute {
		t.Fatalf("ListPending trägt die gültige Zeile nicht: %+v", request)
	}
}

// TestAdministrationRequestRoutingFunctionsRequireCdcAdminMembership belegt das
// `REVOKE … FROM PUBLIC`/`GRANT … TO cdc_admin`-Paar für `cdc.set_route` und
// `cdc.remove_route` (`LH-QA-SEC-002`): eine Rolle ohne
// `cdc_admin`-Mitgliedschaft scheitert mit SQLSTATE 42501 („permission denied
// for function“), und der gescheiterte Aufruf hinterlässt keine Zeile. Der
// Aufruf unter `cdc_admin` gelingt (Gegenprobe). Rot färbende Mutation je
// Funktion: ihre Signatur aus der `REVOKE`-Zeile der Nacharbeit-Datei
// streichen — `PUBLIC` behält `EXECUTE`, der Aufruf unter `cdc_reader` bzw.
// `cdc_capture` gelingt; sie aus der `GRANT`-Zeile streichen — die Gegenprobe
// unter `cdc_admin` scheitert.
func TestAdministrationRequestRoutingFunctionsRequireCdcAdminMembership(t *testing.T) {
	pool, _ := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	const table = "orders_routes_denied"
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE table_name = $1", table)
	})

	callAs := func(role string) (setErr, removeErr error) {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			t.Fatalf("Verbindung reservieren: %v", err)
		}
		defer func() {
			_, _ = conn.Exec(ctx, "RESET ROLE")
			conn.Release()
		}()
		if _, err := conn.Exec(ctx, "SET ROLE "+role); err != nil {
			t.Fatalf("SET ROLE %s: %v", role, err)
		}
		var requestID string
		setErr = conn.QueryRow(ctx,
			"SELECT cdc.set_route($1, $2, $3, $4, $5::json)", administrationRequestSource, "public", table, "verboten", `{"target": "x", "order": 1}`,
		).Scan(&requestID)
		removeErr = conn.QueryRow(ctx,
			"SELECT cdc.remove_route($1, $2, $3, $4)", administrationRequestSource, "public", table, "verboten",
		).Scan(&requestID)
		return setErr, removeErr
	}
	for _, role := range []string{"cdc_reader", "cdc_capture"} {
		setErr, removeErr := callAs(role)
		if !permissionDenied(setErr) {
			t.Fatalf("%s SELECT cdc.set_route(...): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", role, setErr)
		}
		if !permissionDenied(removeErr) {
			t.Fatalf("%s SELECT cdc.remove_route(...): erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", role, removeErr)
		}
	}
	var rows int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM cdc.administration_request WHERE table_name = $1", table).Scan(&rows); err != nil {
		t.Fatalf("Zeilen zählen: %v", err)
	}
	if rows != 0 {
		t.Fatalf("die abgelehnten Aufrufe hinterließen %d Antragszeile(n)", rows)
	}
	setErr, removeErr := callAs("cdc_admin")
	if setErr != nil || removeErr != nil {
		t.Fatalf("Aufruf unter cdc_admin: set_route %v, remove_route %v, wollen beide nil", setErr, removeErr)
	}
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM cdc.administration_request WHERE table_name = $1 AND status = 'pending'", table).Scan(&rows); err != nil {
		t.Fatalf("Zeilen zählen: %v", err)
	}
	if rows != 2 {
		t.Fatalf("Aufruf unter cdc_admin hinterließ %d pending-Zeile(n), wollen 2", rows)
	}
}

// TestAdministrationRequestSetRouteAcceptanceSet trägt die Annahmemenge des
// Parameters `rule_spec` von `cdc.set_route` gegen die reale Instanz: dieselbe
// Menge wie bei `cdc.set_transformation` (`SPEC-019`, Gleichheit hergeleitet,
// hier für die Beispiele des Transformations-Tests gemessen): jede Form läuft
// als Text, der über `::text::json` zum Parameter wird. Eine angenommene Form
// schreibt genau eine `pending`-Zeile und liest ihren Wert als
// `rule_spec::text` zurück; eine abgelehnte Form endet mit einem Fehler des
// Aufrufs und hinterlässt keine Zeile. Rot färbende Mutation (Eingabeseite):
// den Ausdruck `p_rule_spec::jsonb` in `cdc.set_route` durch `NULL::jsonb`
// ersetzen — jede abgelehnte Form wird angenommen, jede angenommene Form mit
// Wert liest NULL.
func TestAdministrationRequestSetRouteAcceptanceSet(t *testing.T) {
	pool, _ := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	const table = "route_spec_acceptance"
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE table_name = $1", table)
	})

	call := func(ruleName string, spec *string) error {
		var requestID string
		return pool.QueryRow(ctx,
			"SELECT cdc.set_route($1, $2, $3, $4, $5::text::json)", administrationRequestSource, "public", table, ruleName, spec,
		).Scan(&requestID)
	}
	stored := func(ruleName string) (rows int, spec *string) {
		if err := pool.QueryRow(ctx,
			"SELECT count(*), min(rule_spec::text) FROM cdc.administration_request WHERE table_name = $1 AND rule_name = $2", table, ruleName,
		).Scan(&rows, &spec); err != nil {
			t.Fatalf("Zeilen von %q lesen: %v", ruleName, err)
		}
		return rows, spec
	}
	text := func(s string) *string { return &s }

	accepted := []struct {
		name string
		spec *string
		want *string
	}{
		{"sql_null", nil, nil},
		{"json_null", text(`null`), text(`null`)},
		{"leeres_objekt", text(`{}`), text(`{}`)},
		{"leeres_array", text(`[]`), text(`[]`)},
		{"zahl", text(`1`), text(`1`)},
		{"regelform", text(`{"target":"eu","order":10,"when":{"column":"region","equals":"eu"}}`), text(`{"when": {"column": "region", "equals": "eu"}, "order": 10, "target": "eu"}`)},
		{"doppelter_schluessel", text(`{"order":1,"order":2}`), text(`{"order": 2}`)},
		{"order_mit_nachkommastelle", text(`{"target":"a","order":10.0}`), text(`{"order": 10.0, "target": "a"}`)},
	}
	for _, tc := range accepted {
		if err := call(tc.name, tc.spec); err != nil {
			t.Errorf("Form %s: erwartet angenommen, erhalten %v", tc.name, err)
			continue
		}
		rows, spec := stored(tc.name)
		if rows != 1 || (spec == nil) != (tc.want == nil) || (spec != nil && *spec != *tc.want) {
			t.Errorf("Form %s: %d Zeile(n), rule_spec gelesen %s, erwartet 1 Zeile mit %s", tc.name, rows, describeSpec(spec), describeSpec(tc.want))
		}
	}

	rejected := []struct {
		name string
		spec string
	}{
		{"nul_im_wert", `{"target":"\u0000"}`},
		{"nul_in_when", `{"when":{"column":"a\u0000","equals":"x"}}`},
		{"zahl_ausserhalb_des_bereichs", `{"order":1e200000}`},
		{"syntaxfehler", `{oops`},
		{"leerer_text", ``},
	}
	for _, tc := range rejected {
		if err := call(tc.name, &tc.spec); err == nil {
			t.Errorf("Form %s: erwartet abgelehnt, der Aufruf gelang", tc.name)
		}
		if rows, _ := stored(tc.name); rows != 0 {
			t.Errorf("Form %s: der abgelehnte Aufruf hinterließ %d Zeile(n)", tc.name, rows)
		}
	}
}

// TestTableActivationRoutingRulesDeriveAppliedRouteRequests trägt die
// Ableitung des dauerhaften Routing-Regelstandes gegen die reale PostgreSQL
// (`SPEC-019`): die `applied`-Zeilen der beiden Routing-Antragsarten tragen den
// Stand, `set_route` trägt eine Regel ein, `remove_route` nimmt sie unter dem
// Namen heraus. Der Test setzt die Antrags-Zeilen direkt (`requested_at` und
// Antrags-ID sind Prüfgegenstand); Bereinigung und Zeilen sind auf die
// Kennungs-Vorsilbe `routes-derivation-` begrenzt.
//
// Die beiden ersten Zeilen tragen denselben `requested_at`: sie stehen in der
// Reihenfolge `b-set` vor `a-remove` eingefügt; die `administration_request_id`
// als Zweitschlüssel ordnet `a-remove` zuerst, die Regel bleibt geführt.
// Zeilen ohne Anteil am Stand: `pending`, `failed`, eine Transformationsregel
// gleichen Namens und eine andere Quelle. Rot färbende Mutationen
// (Eingabeseite): `administration_request_id` aus dem `ORDER BY` von
// `SelectAppliedRoutingRequests` streichen — die gleichzeitigen Zeilen folgen
// der Einfügung, die Regel `tie` fehlt; `'applied'` im `WHERE` durch
// `'pending'` ersetzen — der Stand ist leer; `source_id = $1` streichen —
// die Zeile der anderen Quelle erscheint; `'set_route'` aus dem `IN` streichen
// — der Stand ist leer.
func TestTableActivationRoutingRulesDeriveAppliedRouteRequests(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()

	adapter, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(adapter.Close)
	const otherSource = "src-administration-routes-other"
	if _, err := pool.Exec(ctx, "INSERT INTO cdc.source (source_id, name) VALUES ($1, 'routes-other') ON CONFLICT (source_id) DO NOTHING", otherSource); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE administration_request_id LIKE 'routes-derivation-%'")
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.source WHERE source_id = $1", otherSource)
	})

	const (
		table      = "orders_routes_derivation"
		otherTable = "orders_routes_other"
		insertRule = `INSERT INTO cdc.administration_request
    (administration_request_id, source_id, schema_name, table_name, rule_name, rule_spec, request_kind, requested_at, status)
VALUES ($1, $2, $3, $4, $5, $6::text::jsonb, $7, $8, $9)`
	)
	tiedAt := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	insert := func(id, source, targetTable, kind, ruleName string, spec *string, status string, requestedAt time.Time) {
		t.Helper()
		if _, err := pool.Exec(ctx, insertRule, id, source, "public", targetTable, ruleName, spec, kind, requestedAt, status); err != nil {
			t.Fatalf("Antrags-Zeile %q schreiben: %v", id, err)
		}
	}
	spec := func(target string, order int, column string) *string {
		text := fmt.Sprintf(`{"target": %q, "order": %d}`, target, order)
		if column != "" {
			text = fmt.Sprintf(`{"target": %q, "order": %d, "when": {"column": %q, "equals": "v"}}`, target, order, column)
		}
		return &text
	}

	// Gleicher Zeitstempel, zwei Paare: Set eingefügt vor Remove, das Remove
	// ordnet nach der Antrags-ID zuerst (`a-…` vor `b-…`), die Regel bleibt.
	insert("routes-derivation-b-set", administrationRequestSource, table, "set_route", "tie", spec("eu", 10, "region"), "applied", tiedAt)
	insert("routes-derivation-a-remove", administrationRequestSource, table, "remove_route", "tie", nil, "applied", tiedAt)
	insert("routes-derivation-d-set", administrationRequestSource, table, "set_route", "tie_two", spec("us", 20, "status"), "applied", tiedAt)
	insert("routes-derivation-c-remove", administrationRequestSource, table, "remove_route", "tie_two", nil, "applied", tiedAt)
	// Eine Regel, die wieder herausgenommen ist, lässt keinen Eintrag zurück.
	insert("routes-derivation-f-set", administrationRequestSource, otherTable, "set_route", "gone", spec("x", 5, ""), "applied", tiedAt)
	// Zyklus über drei Zeitpunkte: Set, Remove, Set mit anderem Ziel.
	insert("routes-derivation-g-remove", administrationRequestSource, otherTable, "remove_route", "gone", nil, "applied", tiedAt.Add(time.Second))
	insert("routes-derivation-h-set", administrationRequestSource, table, "set_route", "cycle", spec("a", 30, "note"), "applied", tiedAt.Add(time.Second))
	insert("routes-derivation-i-remove", administrationRequestSource, table, "remove_route", "cycle", nil, "applied", tiedAt.Add(2*time.Second))
	insert("routes-derivation-j-set", administrationRequestSource, table, "set_route", "cycle", spec("b", 40, ""), "applied", tiedAt.Add(3*time.Second))
	// Zeilen ohne Anteil am Stand: anderer Ausgang, andere Antragsart, andere Quelle.
	insert("routes-derivation-k-pending", administrationRequestSource, table, "set_route", "pending_rule", spec("p", 50, ""), "pending", tiedAt)
	insert("routes-derivation-l-failed", administrationRequestSource, table, "set_route", "failed_rule", spec("f", 60, ""), "failed", tiedAt)
	transformation := `{"kind": "rename_column", "column": "name", "to": "title"}`
	insert("routes-derivation-m-transformation", administrationRequestSource, table, "set_transformation", "trans_rule", &transformation, "applied", tiedAt)
	insert("routes-derivation-n-source", otherSource, table, "set_route", "fremd", spec("s", 70, ""), "applied", tiedAt)

	state, err := adapter.RoutingRules(ctx, administrationRequestSource)
	if err != nil {
		t.Fatalf("RoutingRules: %v", err)
	}
	if got, want := describeRoutes(state["public."+table]), "tie:eu@10[region=v];tie_two:us@20[status=v];cycle:b@40[-];"; got != want {
		t.Fatalf("Regelstand public.%s = %q, wollen %q (Zweitschlüssel, Zyklus, ohne pending/failed/fremde Art/fremde Quelle)", table, got, want)
	}
	if got, present := state["public."+otherTable]; present {
		t.Fatalf("Regelstand public.%s = %v, wollen keinen Eintrag (remove_route hat die letzte Regel genommen)", otherTable, got)
	}
	foreign, err := adapter.RoutingRules(ctx, otherSource)
	if err != nil || describeRoutes(foreign["public."+table]) != "fremd:s@70[-];" {
		t.Fatalf("Regelstand der anderen Quelle = %q (%v), wollen nur die eigene Regel fremd", describeRoutes(foreign["public."+table]), err)
	}

	// Ein Remove nach dem letzten Set lässt die Tabelle ohne Eintrag.
	insert("routes-derivation-o-remove", administrationRequestSource, table, "remove_route", "tie", nil, "applied", tiedAt.Add(4*time.Second))
	insert("routes-derivation-p-remove", administrationRequestSource, table, "remove_route", "cycle", nil, "applied", tiedAt.Add(5*time.Second))
	insert("routes-derivation-q-remove", administrationRequestSource, table, "remove_route", "tie_two", nil, "applied", tiedAt.Add(6*time.Second))
	state, err = adapter.RoutingRules(ctx, administrationRequestSource)
	if err != nil {
		t.Fatalf("RoutingRules nach dem Herausnehmen: %v", err)
	}
	if got, present := state["public."+table]; present {
		t.Fatalf("Regelstand public.%s = %v, wollen keinen Eintrag", table, got)
	}

	// Eine Quelle ohne Antrag trägt keinen Stand.
	state, err = adapter.RoutingRules(ctx, model.SourceID("src-administration-ohne-antraege"))
	if err != nil || len(state) != 0 {
		t.Fatalf("Regelstand einer Quelle ohne Anträge = %v (%v), wollen leer", state, err)
	}
}

// TestTableActivationRoutingRulesFailVisiblyOnUnparsableRow trägt: eine
// vermerkte Zeile, deren Regelform nicht mehr zu einer Regel führt, endet
// sichtbar, statt den Stand um sie zu verkürzen. Rot färbende Mutation: den
// Fehler von `model.FoldRoutes` in `sqlexec.ReadRoutingRules` verwerfen — der
// Stand erscheint ohne die Zeile.
func TestTableActivationRoutingRulesFailVisiblyOnUnparsableRow(t *testing.T) {
	pool, dsn := newTestAdministrationRequestPool(t)
	ctx := context.Background()
	adapter, err := postgresstorage.NewTableActivation(ctx, dsn)
	if err != nil {
		t.Fatalf("NewTableActivation: %v", err)
	}
	t.Cleanup(adapter.Close)
	const source = "src-administration-routes-unparsable"
	if _, err := pool.Exec(ctx, "INSERT INTO cdc.source (source_id, name) VALUES ($1, 'unparsable') ON CONFLICT (source_id) DO NOTHING", source); err != nil {
		t.Fatalf("Quelle-Zeile: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE source_id = $1", source)
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.source WHERE source_id = $1", source)
	})
	if _, err := pool.Exec(ctx, `INSERT INTO cdc.administration_request
    (administration_request_id, source_id, schema_name, table_name, rule_name, rule_spec, request_kind, status)
VALUES ('routes-unparsable-1', $1, 'public', 'orders_unparsable', 'regel', '{"target": "x"}'::jsonb, 'set_route', 'applied')`, source); err != nil {
		t.Fatalf("Antrags-Zeile: %v", err)
	}
	if _, err := adapter.RoutingRules(ctx, source); !stderrors.Is(err, domainerrors.ErrInvalidRuleSpec) {
		t.Fatalf("RoutingRules = %v, wollen ErrInvalidRuleSpec", err)
	}
}

// TestAdministrationRequestRoutingCallsKeepCallOrderInOneTransaction trägt die
// Zusage von `SPEC-019` gegen die realen SQL-Funktionen: `cdc.remove_route` vor
// `cdc.set_route` derselben Regel und umgekehrt, in einer Transaktion
// aufgerufen, werden in der Aufruf-Reihenfolge gelesen, und der abgeleitete
// Stand ist der der Aufrufe — vorwärts trägt die Tabelle die neue Regel (Ziel
// `second`), in der Umkehrung ist sie entfernt. Die Kennung ordnet gleichzeitige
// Zeilen nach einer zufälligen UUID; der Test läuft deshalb 40 Durchläufe mit
// je eigenen Tabellen. Rot färbende Mutation je Funktion: in
// `tools/schema/nacharbeit-administration.sql` `clock_timestamp()` durch
// `now()` ersetzen — in `cdc.set_route` färbt sich die Vorwärts-Folge rot (das
// Set trägt den Transaktionsbeginn und sortiert vor das Remove), in
// `cdc.remove_route` die Umkehrung.
func TestAdministrationRequestRoutingCallsKeepCallOrderInOneTransaction(t *testing.T) {
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
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.administration_request WHERE table_name LIKE 'orders_route_order_%'")
	})

	const (
		iterations = 40
		firstSpec  = `{"target": "first", "order": 10}`
		secondSpec = `{"target": "second", "order": 10}`
	)
	type step struct{ query string }
	setStep := step{"SELECT cdc.set_route($1, 'public', $2, 'r', '" + secondSpec + "'::json)"}
	removeStep := step{"SELECT cdc.remove_route($1, 'public', $2, 'r')"}
	sequences := []struct {
		name  string
		steps []step
		want  string
	}{
		{"vorwärts (remove, set)", []step{removeStep, setStep}, "r:second@10[-];"},
		{"umgekehrt (set, remove)", []step{setStep, removeStep}, ""},
	}
	for i := 0; i < iterations; i++ {
		for s, seq := range sequences {
			label := fmt.Sprintf("Durchlauf %d (%s)", i, seq.name)
			table := fmt.Sprintf("orders_route_order_%02d_%d", i, s)
			if _, err := pool.Exec(ctx, `INSERT INTO cdc.administration_request
    (administration_request_id, source_id, schema_name, table_name, rule_name, rule_spec, request_kind, requested_at, status)
VALUES ($1, $2, 'public', $3, 'r', $4::text::jsonb, 'set_route', current_timestamp - interval '1 hour', 'applied')`,
				fmt.Sprintf("route-order-%02d-%d-0-set", i, s), administrationRequestSource, table, firstSpec); err != nil {
				t.Fatalf("%s: Vorgeschichte schreiben: %v", label, err)
			}
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatalf("%s: Begin: %v", label, err)
			}
			var called []model.AdministrationRequestID
			for _, st := range seq.steps {
				var id string
				if err := tx.QueryRow(ctx, st.query, administrationRequestSource, table).Scan(&id); err != nil {
					_ = tx.Rollback(ctx)
					t.Fatalf("%s: %s: %v", label, st.query, err)
				}
				called = append(called, model.AdministrationRequestID(id))
			}
			if err := tx.Commit(ctx); err != nil {
				t.Fatalf("%s: Commit: %v", label, err)
			}

			pending, err := queue.ListPending(ctx)
			if err != nil {
				t.Fatalf("%s: ListPending: %v", label, err)
			}
			inCall := map[model.AdministrationRequestID]bool{called[0]: true, called[1]: true}
			var mineIDs []model.AdministrationRequestID
			var mine []model.AdministrationRequest
			for _, row := range pending {
				if row.Rejected == nil && inCall[row.Request.ID] {
					mineIDs = append(mineIDs, row.Request.ID)
					mine = append(mine, row.Request)
				}
			}
			if fmt.Sprint(mineIDs) != fmt.Sprint(called) {
				t.Fatalf("%s: ListPending-Ordnung = %v, wollen die Aufruf-Reihenfolge %v", label, mineIDs, called)
			}
			for _, request := range mine {
				if err := queue.MarkApplied(ctx, request.ID); err != nil {
					t.Fatalf("%s: MarkApplied %q: %v", label, request.ID, err)
				}
			}
			derived, err := activation.RoutingRules(ctx, administrationRequestSource)
			if err != nil {
				t.Fatalf("%s: RoutingRules: %v", label, err)
			}
			if got := describeRoutes(derived["public."+table]); got != seq.want {
				t.Fatalf("%s: abgeleiteter Stand = %q, wollen %q (der Stand der Aufrufe)", label, got, seq.want)
			}
		}
	}
}
