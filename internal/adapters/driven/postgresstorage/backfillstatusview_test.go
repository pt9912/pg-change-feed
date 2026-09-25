package postgresstorage_test

import (
	"strconv"
	"testing"
	"time"
)

// `cdc.backfill_status` liefert je Tabelle den zuletzt beantragten Run
// (`SPEC-029`): das größte `requested_at` gewinnt, bei gleichem Zeitpunkt die
// größere Run-Kennung; verschiedene Tabellen tragen je eine Zeile. Die
// geschätzte Zeilenzahl bleibt NULL („unbekannt“) und wird nie zu 0, die
// bekannte Schätzung 0 bleibt 0, und die zwei Warn-Spalten tragen den Wert der
// Run-Zeile. Rot färbende Mutationen (je eine, am `query:`-Knoten der View in
// `tools/schema/schema.yaml`): `requested_at DESC` zu `ASC` ändern — die
// Auswahl liefert den älteren Run; `run_id DESC` zu `ASC` ändern — der
// Gleichstand liefert die kleinere Kennung; das `DISTINCT ON` entfernen — die
// Tabelle trägt mehrere Zeilen; `r.warn_estimated_size` durch `false` ersetzen
// — die Sicht zeigt die Warnung „Größe“ nicht.
func TestBackfillStatusViewShowsTheLatestRunPerTable(t *testing.T) {
	f := newBackfillFixture(t)
	base := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	t.Cleanup(func() { f.exec("DELETE FROM cdc.backfill_run WHERE run_id LIKE 'vwst-%'") })

	insert := func(id, table, status string, requestedAt time.Time, estimate any, warnDuration bool) {
		t.Helper()
		f.exec(`INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status, requested_at, estimated_rows, warn_duration)
		        VALUES ($1, $2, 'public', $3, $4, $5, $6, $7)`, id, backfillTestSource, table, status, requestedAt, estimate, warnDuration)
	}
	insert("vwst-a-1", "vwst_a", "completed", base, nil, false)
	insert("vwst-a-2", "vwst_a", "failed", base.Add(time.Hour), nil, true)
	insert("vwst-b-1", "vwst_b", "completed", base, nil, false)
	insert("vwst-b-2", "vwst_b", "queued", base, nil, false)
	insert("vwst-c-1", "vwst_c", "queued", base, nil, false)
	insert("vwst-d-1", "vwst_d", "queued", base, int64(0), false)
	f.exec(`INSERT INTO cdc.backfill_run (run_id, source_id, schema_name, table_name, status, requested_at, estimated_rows, warn_estimated_size)
	        VALUES ('vwst-e-1', $1, 'public', 'vwst_e', 'queued', $2, 7000000, true)`, backfillTestSource, base)

	latest := func(table string) (id, status string, estimate *int64, warnSize, warnDuration bool, rows int) {
		t.Helper()
		rows = f.count("SELECT count(*) FROM cdc.backfill_status WHERE source_id = $1 AND table_name = $2", backfillTestSource, table)
		f.scan(`SELECT run_id, status, estimated_rows, warn_estimated_size, warn_duration
		        FROM cdc.backfill_status WHERE source_id = $1 AND table_name = $2 LIMIT 1`,
			[]any{backfillTestSource, table}, &id, &status, &estimate, &warnSize, &warnDuration)
		return
	}
	// Die Meldung zeigt den Wert der Schätzung, nicht den Zeiger.
	text := func(estimate *int64) string {
		if estimate == nil {
			return "unbekannt (NULL)"
		}
		return strconv.FormatInt(*estimate, 10)
	}
	if id, status, _, _, warnDuration, rows := latest("vwst_a"); rows != 1 || id != "vwst-a-2" || status != "failed" || !warnDuration {
		t.Errorf("vwst_a: %d Zeile(n), Run %q, Status %q, warn_duration %v — erwartet 1, vwst-a-2, failed, true", rows, id, status, warnDuration)
	}
	if id, _, _, _, _, rows := latest("vwst_b"); rows != 1 || id != "vwst-b-2" {
		t.Errorf("vwst_b (Gleichstand in requested_at): %d Zeile(n), Run %q — erwartet 1, vwst-b-2", rows, id)
	}
	if _, _, estimate, warnSize, warnDuration, _ := latest("vwst_c"); estimate != nil || warnSize || warnDuration {
		t.Errorf("vwst_c: Schätzung %s, Warnungen %v/%v — erwartet unbekannt (NULL), false/false", text(estimate), warnSize, warnDuration)
	}
	if _, _, estimate, _, _, _ := latest("vwst_d"); estimate == nil || *estimate != 0 {
		t.Errorf("vwst_d: Schätzung %s — erwartet die bekannte Schätzung 0", text(estimate))
	}
	if _, _, estimate, warnSize, warnDuration, _ := latest("vwst_e"); estimate == nil || *estimate != 7000000 || !warnSize || warnDuration {
		t.Errorf("vwst_e: Schätzung %s, Warnungen %v/%v — erwartet 7000000, true/false", text(estimate), warnSize, warnDuration)
	}
}
