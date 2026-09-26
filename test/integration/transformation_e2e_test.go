package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

// transformationChange trägt eine gelesene Change der Sicht `cdc.changes` mit
// beiden Row Images als Objekte; ein fehlendes Bild ist nil. `raw` hält die
// Textform beider Bilder für den Vergleich zweier Lesungen.
type transformationChange struct {
	operation string
	oldImage  map[string]string
	newImage  map[string]string
	raw       string
}

// requestSetRule legt einen `set_transformation`-Antrag an und liefert seine
// Kennung; die Regelform wird als Text an die SQL-Funktion gereicht.
func (e *backfillEnv) requestSetRule(t *testing.T, table, ruleName, spec string) string {
	t.Helper()
	var requestID string
	if err := e.pool.QueryRow(context.Background(),
		"SELECT cdc.set_transformation($1, 'public', $2, $3, $4::text::json)", e2eSource, table, ruleName, spec,
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.set_transformation(%s, %s): %v", table, ruleName, err)
	}
	return requestID
}

// requestRemoveRule legt einen `remove_transformation`-Antrag an und liefert
// seine Kennung.
func (e *backfillEnv) requestRemoveRule(t *testing.T, table, ruleName string) string {
	t.Helper()
	var requestID string
	if err := e.pool.QueryRow(context.Background(),
		"SELECT cdc.remove_transformation($1, 'public', $2, $3)", e2eSource, table, ruleName,
	).Scan(&requestID); err != nil {
		t.Fatalf("cdc.remove_transformation(%s, %s): %v", table, ruleName, err)
	}
	return requestID
}

// awaitRequestOutcome wartet auf den Endzustand eines Antrags und liefert
// Status (`applied` oder `failed`) und Fehlertext.
func (e *backfillEnv) awaitRequestOutcome(t *testing.T, requestID string) (string, string) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	var status, message string
	for time.Now().Before(deadline) {
		if err := e.pool.QueryRow(context.Background(),
			"SELECT status, coalesce(error_message, '') FROM cdc.administration_request WHERE administration_request_id = $1", requestID,
		).Scan(&status, &message); err != nil {
			t.Fatalf("Antrag %s lesen: %v", requestID, err)
		}
		if status == "applied" || status == "failed" {
			return status, message
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("Antrag %s wurde nicht innerhalb der Zeitspanne vermerkt (status=%q)", requestID, status)
	return "", ""
}

// setRule setzt eine Regel und wartet auf den Vermerk `applied`.
func (e *backfillEnv) setRule(t *testing.T, table, ruleName, spec string) {
	t.Helper()
	e.awaitRequestApplied(t, e.requestSetRule(t, table, ruleName, spec))
}

// removeRulesAtEnd nimmt die genannten Regeln beim Aufräumen zurück: die
// Tabelle trägt danach keinen Regelstand mehr in spätere Phasen. Eine nicht
// gesetzte Regel meldet die Aufräumung nur, wenn der Test nicht fehlschlug.
func (e *backfillEnv) removeRulesAtEnd(t *testing.T, table string, ruleNames ...string) {
	t.Helper()
	t.Cleanup(func() {
		for _, name := range ruleNames {
			status, message := e.awaitRequestOutcome(t, e.requestRemoveRule(t, table, name))
			if status != "applied" && !t.Failed() {
				t.Errorf("Regel %s von %s nicht zurückgenommen: %s", name, table, message)
			}
		}
	})
}

// readTableChanges liest die Changes einer Tabelle über `cdc.changes` in
// Lese-Ordnung.
func (e *backfillEnv) readTableChanges(t *testing.T, table string) []transformationChange {
	t.Helper()
	rows, err := e.pool.Query(context.Background(), `
		SELECT operation, old_data::text, new_data::text FROM cdc.changes
		WHERE source_id = $1 AND table_name = $2
		ORDER BY commit_position, transaction_id, sequence`, e2eSource, table)
	if err != nil {
		t.Fatalf("cdc.changes lesen: %v", err)
	}
	defer rows.Close()
	var changes []transformationChange
	for rows.Next() {
		var operation string
		var oldText, newText *string
		if err := rows.Scan(&operation, &oldText, &newText); err != nil {
			t.Fatalf("cdc.changes-Zeile lesen: %v", err)
		}
		change := transformationChange{operation: operation}
		if oldText != nil {
			if err := json.Unmarshal([]byte(*oldText), &change.oldImage); err != nil {
				t.Fatalf("Alt-Bild lesen: %v (%s)", err, *oldText)
			}
			change.raw += "old=" + *oldText
		}
		if newText != nil {
			if err := json.Unmarshal([]byte(*newText), &change.newImage); err != nil {
				t.Fatalf("Neu-Bild lesen: %v (%s)", err, *newText)
			}
			change.raw += " new=" + *newText
		}
		changes = append(changes, change)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("cdc.changes lesen: %v", err)
	}
	return changes
}

// awaitTableChanges liest die Changes einer Tabelle, bis mindestens `want`
// erfasst sind.
func (e *backfillEnv) awaitTableChanges(t *testing.T, table string, want int) []transformationChange {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if changes := e.readTableChanges(t, table); len(changes) >= want {
			return changes
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("%s: %d Changes wurden nicht innerhalb der Zeitspanne erfasst", table, want)
	return nil
}

// expectImage bricht den Test, wenn das Bild vom erwarteten Objekt abweicht.
func expectImage(t *testing.T, what string, got, want map[string]string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: Bild %v, erwartet %v", what, got, want)
	}
}

// TestE2ETransformationRulesShapeBothImages trägt `LH-FA-CFG-007` (Happy
// Path) am laufenden Feed-Container: eine per SQL beantragte `rename_column`-
// und eine `map_value`-Regel prägen nach dem Vermerk `applied` beide Row
// Images jeder Operation einer Tabelle mit voller Replica-Identität. Der
// Zielname ersetzt den Quellschlüssel, ein abgebildeter Wert steht unter der
// Spalte, ein nicht abgebildeter Wert bleibt unverändert, und ein fehlender
// Wert (`NULL`) bleibt abwesend, ohne Zielschlüssel. Die vor den Regeln
// erfasste Change bleibt in Rohform, und zwei Lesungen liefern dieselben
// Bilder.
func TestE2ETransformationRulesShapeBothImages(t *testing.T) {
	env := newBackfillEnv(t)
	ctx := context.Background()
	const table = "feed_e2e_transform_images"

	for _, statement := range []string{
		fmt.Sprintf("CREATE TABLE public.%s (id int PRIMARY KEY, name text, status text, note text)", table),
		fmt.Sprintf("ALTER TABLE public.%s REPLICA IDENTITY FULL", table),
	} {
		if _, err := env.pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Tabelle anlegen: %v (%s)", err, statement)
		}
	}
	env.enableTable(t, table)

	if _, err := env.pool.Exec(ctx, fmt.Sprintf("INSERT INTO public.%s (id, name, status, note) VALUES (1, 'Ada', 'o', 'n1')", table)); err != nil {
		t.Fatalf("INSERT vor den Regeln: %v", err)
	}
	env.awaitTableChanges(t, table, 1)

	env.setRule(t, table, "kundenname", `{"kind":"rename_column","column":"name","to":"customer_name"}`)
	env.setRule(t, table, "status_lesbar", `{"kind":"map_value","column":"status","values":{"o":"open","c":"closed"}}`)
	env.removeRulesAtEnd(t, table, "kundenname", "status_lesbar")

	for _, statement := range []string{
		fmt.Sprintf("UPDATE public.%s SET status = 'c' WHERE id = 1", table),
		fmt.Sprintf("INSERT INTO public.%s (id, name, status, note) VALUES (2, 'Bob', 'x', NULL)", table),
		fmt.Sprintf("INSERT INTO public.%s (id, name, status, note) VALUES (3, NULL, 'c', 'n3')", table),
		fmt.Sprintf("DELETE FROM public.%s WHERE id = 2", table),
	} {
		if _, err := env.pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Quelländerung: %v (%s)", err, statement)
		}
	}
	changes := env.awaitTableChanges(t, table, 5)
	if len(changes) != 5 {
		t.Fatalf("Changes der Tabelle: %d, erwartet 5", len(changes))
	}

	// Die Change vor den Regeln behält ihre Rohform.
	expectImage(t, "INSERT vor den Regeln, Neu-Bild", changes[0].newImage,
		map[string]string{"id": "1", "name": "Ada", "status": "o", "note": "n1"})

	// UPDATE: beide Bilder tragen die Form, das Alt-Bild den Wert vor der
	// Änderung.
	expectImage(t, "UPDATE, Alt-Bild", changes[1].oldImage,
		map[string]string{"id": "1", "customer_name": "Ada", "status": "open", "note": "n1"})
	expectImage(t, "UPDATE, Neu-Bild", changes[1].newImage,
		map[string]string{"id": "1", "customer_name": "Ada", "status": "closed", "note": "n1"})

	// INSERT mit einem nicht abgebildeten Wert und einem fehlenden Wert.
	expectImage(t, "INSERT mit nicht abgebildetem Wert und NULL in note, Neu-Bild", changes[2].newImage,
		map[string]string{"id": "2", "customer_name": "Bob", "status": "x"})
	if changes[2].oldImage != nil {
		t.Fatalf("INSERT trägt ein Alt-Bild: %v", changes[2].oldImage)
	}

	// INSERT ohne Wert in der umbenannten Spalte: weder Quell- noch
	// Zielschlüssel.
	expectImage(t, "INSERT mit NULL in name, Neu-Bild", changes[3].newImage,
		map[string]string{"id": "3", "status": "closed", "note": "n3"})

	// DELETE: das Alt-Bild trägt die Form, das Neu-Bild fehlt.
	expectImage(t, "DELETE, Alt-Bild", changes[4].oldImage,
		map[string]string{"id": "2", "customer_name": "Bob", "status": "x"})
	if changes[4].newImage != nil {
		t.Fatalf("DELETE trägt ein Neu-Bild: %v", changes[4].newImage)
	}

	// Zweite Lesung: dieselben Bilder in derselben Ordnung.
	reread := env.readTableChanges(t, table)
	if len(reread) != len(changes) {
		t.Fatalf("zweite Lesung: %d Changes, erste %d", len(reread), len(changes))
	}
	for i := range changes {
		if reread[i].raw != changes[i].raw || reread[i].operation != changes[i].operation {
			t.Fatalf("zweite Lesung, Change %d: %s %s, erste %s %s", i, reread[i].operation, reread[i].raw, changes[i].operation, changes[i].raw)
		}
	}
}

// TestE2ETransformationConflictsFailWithSpecText trägt `LH-FA-CFG-007`
// (Boundary) am laufenden Feed-Container: je eine Verletzung der
// Konfliktfreiheit endet real `failed` mit dem Klartext der Spec und der
// Adresse, und der Regelstand der Tabelle bleibt unverändert — die danach
// erfasste Change trägt weiter die Form der gültigen Regeln. Jeder
// Negativfall weicht in genau einem Feld von der Gegenprobe ab, die
// `applied` endet und ab der nächsten Change wirkt.
func TestE2ETransformationConflictsFailWithSpecText(t *testing.T) {
	env := newBackfillEnv(t)
	ctx := context.Background()
	const table = "feed_e2e_transform_conflict"
	address := func(name string) string { return "public." + table + "." + name }

	if _, err := env.pool.Exec(ctx, fmt.Sprintf("CREATE TABLE public.%s (id int PRIMARY KEY, name text, status text, note text)", table)); err != nil {
		t.Fatalf("Tabelle anlegen: %v", err)
	}
	env.enableTable(t, table)
	env.setRule(t, table, "kundenname", `{"kind":"rename_column","column":"name","to":"customer_name"}`)
	env.setRule(t, table, "status_lesbar", `{"kind":"map_value","column":"status","values":{"o":"open","c":"closed"}}`)
	env.removeRulesAtEnd(t, table, "kundenname", "status_lesbar", "notiz_memo")

	// Die Gegenprobe ist die Regel `notiz_memo` (note wird zu memo); jeder
	// Negativfall ändert genau ein Feld daran.
	const validSpec = `{"kind":"rename_column","column":"note","to":"memo"}`
	violations := []struct {
		name     string
		ruleName string
		spec     string
		remove   bool
		want     string
	}{
		{"K1 Regelname vergeben", "kundenname", validSpec, false,
			"Regelname bereits vergeben: " + address("kundenname")},
		{"K2 Spalte trägt bereits eine Regel", "notiz_memo", `{"kind":"rename_column","column":"name","to":"memo"}`, false,
			"Spalte trägt bereits eine Regel: " + address("name")},
		{"K3 Zielname gleicht dem Zielnamen einer anderen Regel", "notiz_memo", `{"kind":"rename_column","column":"note","to":"customer_name"}`, false,
			"Zielname kollidiert mit einer anderen Regel: " + address("customer_name")},
		{"K3 Zielname gleicht einer Spalte der Tabelle", "notiz_memo", `{"kind":"rename_column","column":"note","to":"id"}`, false,
			"Zielname kollidiert mit einer Spalte der Tabelle: " + address("id")},
		{"K4 Spalte fehlt an der Quelle", "notiz_memo", `{"kind":"rename_column","column":"nicht_vorhanden","to":"memo"}`, false,
			"Spalte existiert nicht an der Quelle: " + address("nicht_vorhanden")},
		{"K4 Regelname nicht geführt", "gibt_es_nicht", "", true,
			"Regelname nicht geführt: " + address("gibt_es_nicht")},
	}
	for _, violation := range violations {
		var requestID string
		if violation.remove {
			requestID = env.requestRemoveRule(t, table, violation.ruleName)
		} else {
			requestID = env.requestSetRule(t, table, violation.ruleName, violation.spec)
		}
		status, message := env.awaitRequestOutcome(t, requestID)
		if status != "failed" || message != violation.want {
			t.Fatalf("%s: Antrag endete %s mit %q, erwartet failed mit %q", violation.name, status, message, violation.want)
		}
	}

	// Der Regelstand ist unverändert: die Change trägt die Form der zwei
	// gültigen Regeln und kein Ergebnis eines abgelehnten Antrags.
	if _, err := env.pool.Exec(ctx, fmt.Sprintf("INSERT INTO public.%s (id, name, status, note) VALUES (1, 'Ada', 'o', 'n1')", table)); err != nil {
		t.Fatalf("INSERT nach den abgelehnten Anträgen: %v", err)
	}
	changes := env.awaitTableChanges(t, table, 1)
	expectImage(t, "Change nach den abgelehnten Anträgen", changes[0].newImage,
		map[string]string{"id": "1", "customer_name": "Ada", "status": "open", "note": "n1"})

	// Gegenprobe: dieselbe Regel ohne die Verletzung endet `applied` und
	// wirkt ab der nächsten Change.
	status, message := env.awaitRequestOutcome(t, env.requestSetRule(t, table, "notiz_memo", validSpec))
	if status != "applied" {
		t.Fatalf("Gegenprobe endete %s (%s), erwartet applied", status, message)
	}
	if _, err := env.pool.Exec(ctx, fmt.Sprintf("INSERT INTO public.%s (id, name, status, note) VALUES (2, 'Bea', 'c', 'n2')", table)); err != nil {
		t.Fatalf("INSERT nach der Gegenprobe: %v", err)
	}
	changes = env.awaitTableChanges(t, table, 2)
	expectImage(t, "Change nach der Gegenprobe", changes[1].newImage,
		map[string]string{"id": "2", "customer_name": "Bea", "status": "closed", "memo": "n2"})
}
