package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"
)

// routedChange trägt eine gelesene Change der Sicht `cdc.changes` mit ihrem
// Zustellziel; `hasTarget` unterscheidet ein leeres Label (NULL) von einem
// gesetzten.
type routedChange struct {
	changeID  string
	operation string
	origin    string
	rowID     string
	target    string
	hasTarget bool
	newImage  string
}

// requestSetRoute legt einen `set_route`-Antrag an und liefert seine Kennung.
func (e *backfillEnv) requestSetRoute(t *testing.T, table, ruleName, spec string) string {
	t.Helper()
	var requestID string
	if err := e.pool.QueryRow(context.Background(),
		"SELECT cdc.set_route($1, 'public', $2, $3, $4::text::json)", e2eSource, table, ruleName, spec,
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.set_route(%s, %s): %v", table, ruleName, err)
	}
	return requestID
}

// requestRemoveRoute legt einen `remove_route`-Antrag an und liefert seine
// Kennung.
func (e *backfillEnv) requestRemoveRoute(t *testing.T, table, ruleName string) string {
	t.Helper()
	var requestID string
	if err := e.pool.QueryRow(context.Background(),
		"SELECT cdc.remove_route($1, 'public', $2, $3)", e2eSource, table, ruleName,
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.remove_route(%s, %s): %v", table, ruleName, err)
	}
	return requestID
}

// requestExclude legt einen `exclude_column`-Antrag an und liefert seine
// Kennung.
func (e *backfillEnv) requestExclude(t *testing.T, table, column string) string {
	t.Helper()
	var requestID string
	if err := e.pool.QueryRow(context.Background(),
		"SELECT cdc.exclude_column($1, 'public', $2, $3)", e2eSource, table, column,
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.exclude_column(%s, %s): %v", table, column, err)
	}
	return requestID
}

// setRoute setzt eine Routing-Regel und wartet auf den Vermerk `applied`.
func (e *backfillEnv) setRoute(t *testing.T, table, ruleName, spec string) {
	t.Helper()
	e.awaitRequestApplied(t, e.requestSetRoute(t, table, ruleName, spec))
}

// removeRoute nimmt eine Routing-Regel zurück und wartet auf `applied`.
func (e *backfillEnv) removeRoute(t *testing.T, table, ruleName string) {
	t.Helper()
	e.awaitRequestApplied(t, e.requestRemoveRoute(t, table, ruleName))
}

// mustExec führt eine Anweisung gegen die Quelle aus oder bricht den Test ab.
func (e *backfillEnv) mustExec(t *testing.T, statement string, args ...any) {
	t.Helper()
	if _, err := e.pool.Exec(context.Background(), statement, args...); err != nil {
		t.Fatalf("%v (%s)", err, statement)
	}
}

// readRouted liest die Changes einer Tabelle über `cdc.changes` in
// Lese-Ordnung samt Zustellziel.
func (e *backfillEnv) readRouted(t *testing.T, table string) []routedChange {
	t.Helper()
	rows, err := e.pool.Query(context.Background(), `
		SELECT change_id, operation, origin, coalesce(new_data->>'id', old_data->>'id', ''), route_target, coalesce(new_data::text, '')
		FROM cdc.changes
		WHERE source_id = $1 AND table_name = $2
		ORDER BY commit_position, transaction_id, sequence`, e2eSource, table)
	if err != nil {
		t.Fatalf("cdc.changes lesen: %v", err)
	}
	defer rows.Close()
	var changes []routedChange
	for rows.Next() {
		var change routedChange
		var target *string
		if err := rows.Scan(&change.changeID, &change.operation, &change.origin, &change.rowID, &target, &change.newImage); err != nil {
			t.Fatalf("cdc.changes-Zeile lesen: %v", err)
		}
		if target != nil {
			change.target, change.hasTarget = *target, true
		}
		changes = append(changes, change)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("cdc.changes lesen: %v", err)
	}
	return changes
}

// awaitRouted liest die Changes einer Tabelle, bis mindestens `want`
// erfasst sind.
func (e *backfillEnv) awaitRouted(t *testing.T, table string, want int) []routedChange {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if changes := e.readRouted(t, table); len(changes) >= want {
			return changes
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("%s: %d Changes wurden nicht innerhalb der Zeitspanne erfasst", table, want)
	return nil
}

// targetsOf liefert das Ziel je Change in Lese-Ordnung; NULL erscheint als
// leere Zeichenkette.
func targetsOf(changes []routedChange) []string {
	targets := make([]string, len(changes))
	for i, change := range changes {
		targets[i] = change.target
	}
	return targets
}

// idsWithTarget liefert die Kennungen der Changes, deren Ziel `target`
// trägt, sortiert; `NULL` wählt die Changes ohne Ziel.
func idsWithTarget(changes []routedChange, target string, isNull bool) []string {
	ids := []string{}
	for _, change := range changes {
		if (isNull && !change.hasTarget) || (!isNull && change.hasTarget && change.target == target) {
			ids = append(ids, change.changeID)
		}
	}
	sort.Strings(ids)
	return ids
}

// sqlIDsWithTarget liest die Kennungen einer Tabelle mit `WHERE route_target`
// aus der Sicht, sortiert.
func (e *backfillEnv) sqlIDsWithTarget(t *testing.T, table, target string, isNull bool) []string {
	t.Helper()
	query := "SELECT change_id FROM cdc.changes WHERE source_id = $1 AND table_name = $2 AND route_target = $3"
	args := []any{e2eSource, table, target}
	if isNull {
		query = "SELECT change_id FROM cdc.changes WHERE source_id = $1 AND table_name = $2 AND route_target IS NULL"
		args = args[:2]
	}
	rows, err := e.pool.Query(context.Background(), query, args...)
	if err != nil {
		t.Fatalf("cdc.changes mit route_target lesen: %v", err)
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("change_id lesen: %v", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("cdc.changes lesen: %v", err)
	}
	sort.Strings(ids)
	return ids
}

// serverVersion liest die Version des PostgreSQL-Servers der Quelle.
func (e *backfillEnv) serverVersion(t *testing.T) string {
	t.Helper()
	var version string
	if err := e.pool.QueryRow(context.Background(), "SELECT current_setting('server_version')").Scan(&version); err != nil {
		t.Fatalf("server_version lesen: %v", err)
	}
	return version
}

// httpChangeIDs liest `GET /changes` mit dem reader-Token über den Lesezugriff
// der Feed-Container-API; ein leeres `target` lässt den Parameter weg. Die
// Adresse und das Token reicht der Runner über die Umgebung.
func httpChangeIDs(t *testing.T, table, target string) []string {
	t.Helper()
	base, token := os.Getenv("CDC_INTEGRATION_HTTP_URL"), os.Getenv("CDC_INTEGRATION_HTTP_READER_TOKEN")
	if base == "" || token == "" {
		t.Fatal("CDC_INTEGRATION_HTTP_URL und CDC_INTEGRATION_HTTP_READER_TOKEN sind nicht gesetzt — der Runner reicht beide an den Testlauf")
	}
	query := url.Values{}
	query.Set("source", e2eSource)
	query.Set("schema", "public")
	query.Set("table", table)
	if target != "" {
		query.Set("target", target)
	}
	request, err := http.NewRequest(http.MethodGet, base+"/changes?"+query.Encode(), nil)
	if err != nil {
		t.Fatalf("Request bauen: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("GET /changes: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Antwort lesen: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /changes?%s: Status %d: %s", query.Encode(), response.StatusCode, body)
	}
	var parsed struct {
		Changes []struct {
			ChangeID string `json:"change_id"`
		} `json:"changes"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("Antwort ist kein JSON: %v (%s)", err, body)
	}
	ids := []string{}
	for _, change := range parsed.Changes {
		ids = append(ids, change.ChangeID)
	}
	sort.Strings(ids)
	return ids
}

// expectIDs bricht den Test, wenn zwei Kennungsmengen abweichen.
func expectIDs(t *testing.T, what string, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: change_id %v, erwartet %v", what, got, want)
	}
}

// TestE2ERoutingSelectsTargetsOverChangesView belegt Happy Path und Boundary am
// laufenden Feed-Container (`LH-FA-CFG-008`): eine per SQL beantragte
// Herkunftsregel (ohne `when`) lenkt jede Change ihrer Tabelle auf ihr Ziel,
// zwei Inhaltsregeln (`when`) lenken nach dem Wert der Spalte, und die
// Sicht `cdc.changes` wählt jedes Ziel über `WHERE route_target`. Eine Change
// ohne Treffer trägt `route_target IS NULL` und erscheint nur ungefiltert, ein
// ungefilterter Leser sieht jede Change einschließlich der gerouteten, und bei
// zwei treffenden Regeln gewinnt die kleinere `order`, auch wenn die Regel mit
// der größeren `order` zuerst gesetzt wurde (`SPEC-032`).
func TestE2ERoutingSelectsTargetsOverChangesView(t *testing.T) {
	env := newBackfillEnv(t)
	const content = "feed_e2e_route_content"
	const origin = "feed_e2e_route_origin"
	const ordered = "feed_e2e_route_order"

	for _, table := range []string{content, origin, ordered} {
		env.mustExec(t, fmt.Sprintf("CREATE TABLE public.%s (id int PRIMARY KEY, name text, region text)", table))
		env.enableTable(t, table)
	}

	// Inhaltsregeln: zwei Ziele, kein Standardziel.
	env.setRoute(t, content, "eu_orders", `{"target":"eu","order":10,"when":{"column":"region","equals":"eu"}}`)
	env.setRoute(t, content, "us_orders", `{"target":"us","order":20,"when":{"column":"region","equals":"us"}}`)
	for _, statement := range []string{
		"INSERT INTO public.%s (id, name, region) VALUES (1, 'Ada', 'eu')",
		"INSERT INTO public.%s (id, name, region) VALUES (2, 'Bob', 'us')",
		"INSERT INTO public.%s (id, name, region) VALUES (3, 'Cy', 'asia')",
		"INSERT INTO public.%s (id, name, region) VALUES (4, 'Di', NULL)",
		"INSERT INTO public.%s (id, name, region) VALUES (5, 'Ed', 'eu')",
		"UPDATE public.%s SET region = 'eu' WHERE id = 3",
	} {
		env.mustExec(t, fmt.Sprintf(statement, content))
	}
	changes := env.awaitRouted(t, content, 6)
	if len(changes) != 6 {
		t.Fatalf("Changes der Inhaltsregel-Tabelle: %d, erwartet 6", len(changes))
	}
	if got, want := targetsOf(changes), []string{"eu", "us", "", "", "eu", "eu"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ziele der Inhaltsregeln in Lese-Ordnung: %q, erwartet %q (leer = NULL)", got, want)
	}
	// Ohne Treffer ist das Ziel NULL, nicht eine leere Zeichenkette.
	for _, index := range []int{2, 3} {
		if changes[index].hasTarget {
			t.Fatalf("Change %d (%s) ohne Treffer trägt das Ziel %q statt NULL", index, changes[index].rowID, changes[index].target)
		}
	}
	for _, target := range []string{"eu", "us"} {
		expectIDs(t, "WHERE route_target = "+target, env.sqlIDsWithTarget(t, content, target, false), idsWithTarget(changes, target, false))
	}
	expectIDs(t, "WHERE route_target IS NULL", env.sqlIDsWithTarget(t, content, "", true), idsWithTarget(changes, "", true))
	// Ziel A sieht nur A: kein Ziel wählt eine Change eines anderen.
	if euCount, usCount := len(env.sqlIDsWithTarget(t, content, "eu", false)), len(env.sqlIDsWithTarget(t, content, "us", false)); euCount != 3 || usCount != 1 {
		t.Fatalf("Zahl der Changes je Ziel: eu=%d us=%d, erwartet 3 und 1", euCount, usCount)
	}
	// Ein ungefilterter Leser sieht jede Change, geroutete eingeschlossen.
	var total int
	if err := env.pool.QueryRow(context.Background(),
		"SELECT count(*) FROM cdc.changes WHERE source_id = $1 AND table_name = $2", e2eSource, content).Scan(&total); err != nil {
		t.Fatalf("ungefilterte Lesung: %v", err)
	}
	if total != 6 {
		t.Fatalf("ungefilterte Lesung: %d Changes, erwartet 6 (3 eu, 1 us, 2 ohne Ziel)", total)
	}

	// Herkunftsregel: jede Operation der Tabelle trägt das Ziel, auch DELETE.
	env.setRoute(t, origin, "alle", `{"target":"herkunft","order":1}`)
	for _, statement := range []string{
		"INSERT INTO public.%s (id, name, region) VALUES (1, 'Ada', 'eu')",
		"UPDATE public.%s SET name = 'Ada2' WHERE id = 1",
		"DELETE FROM public.%s WHERE id = 1",
	} {
		env.mustExec(t, fmt.Sprintf(statement, origin))
	}
	originChanges := env.awaitRouted(t, origin, 3)
	if got, want := targetsOf(originChanges), []string{"herkunft", "herkunft", "herkunft"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ziele der Herkunftsregel (INSERT, UPDATE, DELETE): %q, erwartet %q", got, want)
	}

	// Auflösung: die kleinere order gewinnt, nicht die Reihenfolge des Setzens.
	env.setRoute(t, ordered, "spaet", `{"target":"spaet","order":30,"when":{"column":"region","equals":"eu"}}`)
	env.setRoute(t, ordered, "frueh", `{"target":"frueh","order":5,"when":{"column":"name","equals":"Ada"}}`)
	env.setRoute(t, ordered, "rest", `{"target":"sonstige","order":100}`)
	for _, statement := range []string{
		"INSERT INTO public.%s (id, name, region) VALUES (1, 'Ada', 'eu')",
		"INSERT INTO public.%s (id, name, region) VALUES (2, 'Bob', 'eu')",
		"INSERT INTO public.%s (id, name, region) VALUES (3, 'Bob', 'us')",
	} {
		env.mustExec(t, fmt.Sprintf(statement, ordered))
	}
	orderedChanges := env.awaitRouted(t, ordered, 3)
	if got, want := targetsOf(orderedChanges), []string{"frueh", "spaet", "sonstige"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Auflösung bei zwei treffenden Regeln: %q, erwartet %q (order 5 vor order 30, Abschlussregel zuletzt)", got, want)
	}
	t.Logf("ROUTING-HAPPY-PATH Inhaltsregeln %q, Herkunftsregel %q, Auflösung %q", targetsOf(changes), targetsOf(originChanges), targetsOf(orderedChanges))
}

// TestE2ERoutingConflictsFailWithSpecText belegt die Boundary am laufenden
// Feed-Container (`LH-FA-CFG-008`, `LH-QA-SEC-004`): je eine Verletzung von R1
// bis R6 endet real `failed` mit dem Klartext der Spec und der Adresse
// (`SPEC-019`), der Regelstand der Tabelle bleibt unverändert, und die
// Gegenprobe ohne Verletzung endet `applied` und wirkt ab der nächsten Change.
// R3 gilt in beide Richtungen: eine Regel gegen eine ausgeschlossene Spalte
// und ein Ausschluss gegen die Spalte einer Routing-Bedingung enden `failed`;
// nach dem Ausschluss einer Spalte trägt keine Change ihren Wert, und kein Ziel
// hängt an ihr.
func TestE2ERoutingConflictsFailWithSpecText(t *testing.T) {
	env := newBackfillEnv(t)
	const table = "feed_e2e_route_conflict"
	address := func(name string) string { return "public." + table + "." + name }

	env.mustExec(t, fmt.Sprintf("CREATE TABLE public.%s (id int PRIMARY KEY, name text, region text, secret text)", table))
	env.enableTable(t, table)
	env.setRoute(t, table, "eu_orders", `{"target":"eu","order":10,"when":{"column":"region","equals":"eu"}}`)

	type violation struct {
		name     string
		ruleName string
		spec     string
		remove   bool
		want     string
	}
	check := func(violations []violation) {
		t.Helper()
		for _, v := range violations {
			var requestID string
			if v.remove {
				requestID = env.requestRemoveRoute(t, table, v.ruleName)
			} else {
				requestID = env.requestSetRoute(t, table, v.ruleName, v.spec)
			}
			status, message := env.awaitRequestOutcome(t, requestID)
			if status != "failed" || message != v.want {
				t.Fatalf("%s: Antrag endete %s mit %q, erwartet failed mit %q", v.name, status, message, v.want)
			}
		}
	}

	// Stand: eine Regel mit Bedingung.
	check([]violation{
		{"R1 Regelname vergeben", "eu_orders", `{"target":"x","order":20,"when":{"column":"name","equals":"a"}}`, false,
			"Regelname bereits vergeben: " + address("eu_orders")},
		{"R2 order vergeben", "andere", `{"target":"x","order":10,"when":{"column":"name","equals":"a"}}`, false,
			"order bereits vergeben: " + address("10")},
		{"R3 Spalte fehlt an der Quelle", "ohne_spalte", `{"target":"x","order":20,"when":{"column":"nicht_vorhanden","equals":"a"}}`, false,
			"Spalte existiert nicht an der Quelle: " + address("nicht_vorhanden")},
		{"R4 Regel ohne when trägt nicht die höchste order", "fruehrest", `{"target":"x","order":5}`, false,
			"Regel ohne when trägt nicht die höchste order: " + address("fruehrest")},
		{"R5 Bedingung vergeben", "eu_doppelt", `{"target":"eu2","order":30,"when":{"column":"region","equals":"eu"}}`, false,
			"Bedingung bereits vergeben: " + address("region")},
		{"R6 Regelname nicht geführt", "gibt_es_nicht", "", true,
			"Regelname nicht geführt: " + address("gibt_es_nicht")},
		{"Zielname außerhalb des Alphabets", "schlecht", `{"target":"Eu.Bad","order":40,"when":{"column":"name","equals":"a"}}`, false,
			"Zielname ist ungültig: " + address("Eu.Bad")},
		{"unbekannter Schlüssel", "unbekannt", `{"target":"x","order":50,"extra":1}`, false,
			"unbekannter Schlüssel in rule_spec: extra"},
	})

	// Stand mit Abschlussregel: eine zweite Regel ohne when und eine Regel
	// hinter ihr enden failed.
	env.setRoute(t, table, "rest", `{"target":"sonstige","order":100}`)
	check([]violation{
		{"R4 zweite Regel ohne when", "rest2", `{"target":"x","order":200}`, false,
			"Regel ohne when bereits vorhanden: " + address("rest")},
		{"R4 Regel mit when hinter der Abschlussregel", "hinter_rest", `{"target":"x","order":150,"when":{"column":"name","equals":"Bob"}}`, false,
			"order liegt hinter der Regel ohne when: " + address("rest")},
	})

	// R3 in beide Richtungen: der Ausschluss der Spalte secret gelingt, eine
	// Regel auf sie scheitert; der Ausschluss der Bedingungsspalte region
	// scheitert.
	env.awaitRequestApplied(t, env.requestExclude(t, table, "secret"))
	check([]violation{
		{"R3 Regel auf ausgeschlossene Spalte", "auf_geheim", `{"target":"x","order":60,"when":{"column":"secret","equals":"a"}}`, false,
			"Spalte ist ausgeschlossen: " + address("secret")},
	})
	status, message := env.awaitRequestOutcome(t, env.requestExclude(t, table, "region"))
	if want := "Spalte trägt eine Routing-Bedingung: " + address("region"); status != "failed" || message != want {
		t.Fatalf("R3 Ausschluss der Bedingungsspalte: Antrag endete %s mit %q, erwartet failed mit %q", status, message, want)
	}

	// Der Regelstand ist unverändert: nur eu_orders und rest wirken, kein
	// Ziel eines abgelehnten Antrags erscheint, und der Wert der
	// ausgeschlossenen Spalte steht nirgends.
	env.mustExec(t, fmt.Sprintf("INSERT INTO public.%s (id, name, region, secret) VALUES (1, 'Ada', 'eu', 'GeheimWert1')", table))
	env.mustExec(t, fmt.Sprintf("INSERT INTO public.%s (id, name, region, secret) VALUES (2, 'Bob', 'us', 'GeheimWert2')", table))
	changes := env.awaitRouted(t, table, 2)
	if got, want := targetsOf(changes), []string{"eu", "sonstige"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ziele nach den abgelehnten Anträgen: %q, erwartet %q", got, want)
	}
	for _, change := range changes {
		var leaked int
		if err := env.pool.QueryRow(context.Background(), `
			SELECT count(*) FROM cdc.changes
			WHERE source_id = $1 AND change_id = $2
			  AND (new_data::text LIKE '%GeheimWert%' OR coalesce(old_data::text, '') LIKE '%GeheimWert%'
			       OR jsonb_exists(new_data, 'secret') OR coalesce(route_target, '') LIKE '%GeheimWert%')`,
			e2eSource, change.changeID).Scan(&leaked); err != nil {
			t.Fatalf("Wert der ausgeschlossenen Spalte suchen: %v", err)
		}
		if leaked != 0 {
			t.Fatalf("Change %s trägt die ausgeschlossene Spalte secret oder ihren Wert (Bild %s, Ziel %q)", change.changeID, change.newImage, change.target)
		}
		if !jsonHasKey(change.newImage, "region") {
			t.Fatalf("Change %s verlor die Spalte region trotz abgelehntem Ausschluss: %s", change.changeID, change.newImage)
		}
	}

	// Gegenprobe: dieselbe Regelform ohne Verletzung endet applied und wirkt
	// ab der nächsten Change.
	env.setRoute(t, table, "us_orders", `{"target":"us","order":20,"when":{"column":"region","equals":"us"}}`)
	env.mustExec(t, fmt.Sprintf("INSERT INTO public.%s (id, name, region, secret) VALUES (3, 'Cy', 'us', 'GeheimWert3')", table))
	changes = env.awaitRouted(t, table, 3)
	if got, want := changes[2].target, "us"; got != want {
		t.Fatalf("Gegenprobe: Ziel %q, erwartet %q", got, want)
	}
}

// jsonHasKey sagt, ob der Text ein JSON-Objekt mit dem Schlüssel ist.
func jsonHasKey(text, key string) bool {
	var object map[string]any
	if err := json.Unmarshal([]byte(text), &object); err != nil {
		return false
	}
	_, found := object[key]
	return found
}

// TestE2ERoutingReplayKeepsLabelAndBackfillCarriesIt belegt am laufenden
// Feed-Container, dass das Label einer Change bei der Erfassung festgeschrieben
// ist (`LH-FA-CFG-008`, `LH-FA-REA-005`). Dieselbe `change_id` liefert nach
// einer Regeländerung dasselbe Ziel über `cdc.changes` und `GET /changes`, eine
// vor der Regel erfasste Change bleibt ohne Ziel (keine Rückwirkung), und ein
// danach beantragter Backfill-Run trägt das Label des aktuellen Regelstands an
// jeder Bestandszeile mit `origin = 'backfill'`.
func TestE2ERoutingReplayKeepsLabelAndBackfillCarriesIt(t *testing.T) {
	env := newBackfillEnv(t)
	const table = "feed_e2e_route_replay"

	env.mustExec(t, fmt.Sprintf("CREATE TABLE public.%s (id int PRIMARY KEY, name text, region text)", table))
	env.enableTable(t, table)

	// Vor der Regel erfasst: kein Ziel.
	env.mustExec(t, fmt.Sprintf("INSERT INTO public.%s (id, name, region) VALUES (1, 'Ada', 'eu')", table))
	env.awaitRouted(t, table, 1)

	env.setRoute(t, table, "eu_orders", `{"target":"eu","order":10,"when":{"column":"region","equals":"eu"}}`)
	env.mustExec(t, fmt.Sprintf("INSERT INTO public.%s (id, name, region) VALUES (2, 'Bob', 'eu')", table))
	env.mustExec(t, fmt.Sprintf("INSERT INTO public.%s (id, name, region) VALUES (3, 'Cy', 'us')", table))
	before := env.awaitRouted(t, table, 3)
	if got, want := targetsOf(before), []string{"", "eu", ""}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Ziele vor der Regeländerung: %q, erwartet %q (die vor der Regel erfasste Change bleibt ohne Ziel)", got, want)
	}
	expectIDs(t, "GET /changes?target=eu vor der Regeländerung", httpChangeIDs(t, table, "eu"), idsWithTarget(before, "eu", false))

	// Regeländerung: dasselbe Kriterium führt künftig auf ein anderes Ziel.
	env.removeRoute(t, table, "eu_orders")
	env.setRoute(t, table, "eu_neu", `{"target":"europa","order":10,"when":{"column":"region","equals":"eu"}}`)
	env.mustExec(t, fmt.Sprintf("INSERT INTO public.%s (id, name, region) VALUES (4, 'Di', 'eu')", table))
	after := env.awaitRouted(t, table, 4)
	for i := range before {
		if after[i].changeID != before[i].changeID || after[i].target != before[i].target || after[i].hasTarget != before[i].hasTarget {
			t.Fatalf("Change %d nach der Regeländerung: %s/%q/%v, davor %s/%q/%v — das Label ändert sich nicht rückwirkend",
				i, after[i].changeID, after[i].target, after[i].hasTarget, before[i].changeID, before[i].target, before[i].hasTarget)
		}
	}
	if after[3].target != "europa" {
		t.Fatalf("Change nach der Regeländerung: Ziel %q, erwartet europa", after[3].target)
	}
	expectIDs(t, "GET /changes?target=eu nach der Regeländerung", httpChangeIDs(t, table, "eu"), idsWithTarget(before, "eu", false))
	expectIDs(t, "GET /changes?target=europa nach der Regeländerung", httpChangeIDs(t, table, "europa"), []string{after[3].changeID})
	if reread := targetsOf(env.readRouted(t, table)); !reflect.DeepEqual(reread, targetsOf(after)) {
		t.Fatalf("zweite Lesung über cdc.changes: Ziele %q, erste Lesung %q", reread, targetsOf(after))
	}

	// Backfill mit dem aktuellen Regelstand (eu_neu): die Bestandszeilen 1, 2
	// und 4 (region eu) tragen europa, die Zeile 3 (us) kein Ziel.
	var runID string
	if err := env.pool.QueryRow(context.Background(), "SELECT cdc.backfill_table($1, 'public', $2)", e2eSource, table).Scan(&runID); err != nil {
		t.Fatalf("cdc.backfill_table(%s): %v", table, err)
	}
	env.awaitRequestApplied(t, runID)
	env.awaitRunStatus(t, runID, "completed", 60*time.Second)
	all := env.readRouted(t, table)
	backfilled := map[string]routedChange{}
	for _, change := range all {
		if change.origin == "backfill" {
			backfilled[change.rowID] = change
		}
	}
	if len(backfilled) != 4 {
		t.Fatalf("Backfill-Changes: %d, erwartet 4 (Bestand der Tabelle)", len(backfilled))
	}
	for rowID, want := range map[string]string{"1": "europa", "2": "europa", "3": "", "4": "europa"} {
		got := backfilled[rowID]
		if got.target != want || (want != "" && !got.hasTarget) || (want == "" && got.hasTarget) {
			t.Fatalf("Backfill-Change der Zeile %s: Ziel %q (gesetzt %v), erwartet %q", rowID, got.target, got.hasTarget, want)
		}
	}
	// Die Changes des WAL-Pfads behalten ihr Label neben dem Bestand.
	walTargets := []string{}
	for _, change := range all {
		if change.origin == "wal" {
			walTargets = append(walTargets, change.target)
		}
	}
	if want := []string{"", "eu", "", "europa"}; !reflect.DeepEqual(walTargets, want) {
		t.Fatalf("WAL-Changes nach dem Backfill: %q, erwartet %q", walTargets, want)
	}
	t.Logf("ROUTING-REPLAY WAL-Ziele %q, Backfill-Ziele zeilenweise 1=%q 2=%q 3=%q 4=%q",
		walTargets, backfilled["1"].target, backfilled["2"].target, backfilled["3"].target, backfilled["4"].target)
}

// TestE2ERoutingDeleteWithoutFullReplicaIdentity misst den Abwesenheitsfall an
// der laufenden PostgreSQL-Version (`LH-QA-POR-001`, `LH-FA-CFG-008`): eine
// Inhaltsregel auf eine Nicht-Schlüsselspalte trifft bei DELETE ohne volle
// Replica-Identität nicht, weil das Alt-Bild nur den Schlüssel trägt; die
// nächste Regel bestimmt das Ziel, ohne sie bleibt es leer. Die Gegenprobe unter
// voller Replica-Identität trifft. Der Test druckt die Version des
// PostgreSQL-Servers und das gemessene Ziel je Operation; er bricht, wenn die
// Messung von dieser Erwartung abweicht.
func TestE2ERoutingDeleteWithoutFullReplicaIdentity(t *testing.T) {
	env := newBackfillEnv(t)
	version := env.serverVersion(t)

	type measurement struct {
		table    string
		identity string
		rules    bool
		wantDel  string
		wantNull bool
	}
	cases := []measurement{
		{"feed_e2e_route_del_key", "DEFAULT", true, "sonstige", false},
		{"feed_e2e_route_del_bare", "DEFAULT", false, "", true},
		{"feed_e2e_route_del_full", "FULL", true, "eu", false},
	}
	var failures []string
	for _, c := range cases {
		env.mustExec(t, fmt.Sprintf("CREATE TABLE public.%s (id int PRIMARY KEY, name text, region text)", c.table))
		if c.identity == "FULL" {
			env.mustExec(t, fmt.Sprintf("ALTER TABLE public.%s REPLICA IDENTITY FULL", c.table))
		}
		env.enableTable(t, c.table)
		env.setRoute(t, c.table, "eu_orders", `{"target":"eu","order":10,"when":{"column":"region","equals":"eu"}}`)
		if c.rules {
			env.setRoute(t, c.table, "rest", `{"target":"sonstige","order":100}`)
		}
		env.mustExec(t, fmt.Sprintf("INSERT INTO public.%s (id, name, region) VALUES (1, 'Ada', 'eu')", c.table))
		env.mustExec(t, fmt.Sprintf("UPDATE public.%s SET name = 'Ada2' WHERE id = 1", c.table))
		env.mustExec(t, fmt.Sprintf("DELETE FROM public.%s WHERE id = 1", c.table))
		changes := env.awaitRouted(t, c.table, 3)
		var oldKeys string
		if err := env.pool.QueryRow(context.Background(), `
			SELECT coalesce(old_data::text, '') FROM cdc.changes
			WHERE source_id = $1 AND table_name = $2 AND operation = 'DELETE'`, e2eSource, c.table).Scan(&oldKeys); err != nil {
			t.Fatalf("Alt-Bild des DELETE lesen: %v", err)
		}
		insert, update, deletion := changes[0], changes[1], changes[2]
		t.Logf("ROUTING-DELETE-MESSUNG PostgreSQL %s, Replica-Identität %s, Regeln eu%s: INSERT=%q UPDATE=%q DELETE=%q (gesetzt %v), Alt-Bild des DELETE %s",
			version, c.identity, map[bool]string{true: "+rest", false: ""}[c.rules], insert.target, update.target, deletion.target, deletion.hasTarget, oldKeys)
		if insert.target != "eu" || update.target != "eu" {
			failures = append(failures, fmt.Sprintf("%s: INSERT=%q UPDATE=%q, erwartet eu und eu (Neu-Bild trägt region)", c.table, insert.target, update.target))
		}
		if deletion.target != c.wantDel || deletion.hasTarget == c.wantNull {
			failures = append(failures, fmt.Sprintf("%s: DELETE=%q (gesetzt %v), erwartet %q (gesetzt %v)", c.table, deletion.target, deletion.hasTarget, c.wantDel, !c.wantNull))
		}
	}
	if len(failures) > 0 {
		t.Fatalf("DELETE-Messung an PostgreSQL %s weicht von der Erwartung ab: %v", version, failures)
	}
}
