package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// backfillEnv trägt die Verbindung des Backfill-Belegs gegen die Compose-
// Instanz. Alle Wege des Belegs sind extern: SQL-Funktionen, die Views
// `cdc.backfill_status`/`cdc.changes` und die Quelltabelle selbst.
type backfillEnv struct {
	dsn  string
	pool *pgxpool.Pool
}

func newBackfillEnv(t *testing.T) *backfillEnv {
	t.Helper()
	dsn := os.Getenv("CDC_INTEGRATION_DSN")
	if dsn == "" {
		t.Skip("CDC_INTEGRATION_DSN nicht gesetzt — Compose-Integrationstest läuft über make test-integration gegen die Compose-Umgebung")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)
	return &backfillEnv{dsn: dsn, pool: pool}
}

// awaitRequestApplied wartet, bis die Administrations-Goroutine den Antrag
// vermerkt hat; ein `failed`-Antrag endet den Test mit seinem Fehlertext.
func (e *backfillEnv) awaitRequestApplied(t *testing.T, requestID string) {
	t.Helper()
	ctx := context.Background()
	deadline := time.Now().Add(30 * time.Second)
	var status, message string
	for time.Now().Before(deadline) {
		if err := e.pool.QueryRow(ctx,
			"SELECT status, coalesce(error_message, '') FROM cdc.administration_request WHERE administration_request_id = $1", requestID,
		).Scan(&status, &message); err != nil {
			t.Fatalf("Antrag %s lesen: %v", requestID, err)
		}
		switch status {
		case "applied":
			return
		case "failed":
			t.Fatalf("Antrag %s endete failed: %s", requestID, message)
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("Antrag %s wurde nicht innerhalb der Zeitspanne vermerkt (status=%q)", requestID, status)
}

// enableTable legt eine Tabelle an, aktiviert sie über den SQL-Antragsweg
// und wartet auf den Vermerk.
func (e *backfillEnv) enableTable(t *testing.T, table string) {
	t.Helper()
	ctx := context.Background()
	var requestID string
	if err := e.pool.QueryRow(ctx, "SELECT cdc.enable_table($1, 'public', $2)", e2eSource, table).Scan(&requestID); err != nil {
		t.Fatalf("cdc.enable_table(%s): %v", table, err)
	}
	e.awaitRequestApplied(t, requestID)
}

// awaitRunStatus wartet, bis der Run den Status `want` trägt; ein anderer
// Endzustand endet den Test mit dem Fehlertext des Runs.
func (e *backfillEnv) awaitRunStatus(t *testing.T, runID, want string, within time.Duration) {
	t.Helper()
	ctx := context.Background()
	deadline := time.Now().Add(within)
	var status, message string
	for time.Now().Before(deadline) {
		if err := e.pool.QueryRow(ctx,
			"SELECT status, coalesce(error_message, '') FROM cdc.backfill_run WHERE run_id = $1", runID,
		).Scan(&status, &message); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			t.Fatalf("Run %s lesen: %v", runID, err)
		}
		if status == want {
			return
		}
		if status == "completed" || status == "failed" || status == "interrupted" {
			t.Fatalf("Run %s endete %s statt %s: %s", runID, status, want, message)
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("Run %s erreichte %s nicht innerhalb der Zeitspanne (status=%q)", runID, want, status)
}

// replayWriter ändert die ihm zugeordneten Schlüssel einer Tabelle
// nebenläufig zum Run: aktualisieren, einfügen (auch einen zuvor gelöschten
// Schlüssel wieder) und löschen, je in einer eigenen Transaktion.
type replayWriter struct {
	table   string
	rng     *rand.Rand
	owned   []int
	deleted []int
	next    int
}

// step führt eine Änderung aus; `counter` macht jeden geschriebenen Wert
// eindeutig.
func (w *replayWriter) step(ctx context.Context, pool *pgxpool.Pool, counter int) error {
	value := fmt.Sprintf("w-%d-%d", w.next, counter)
	var note *string
	if counter%2 == 0 {
		n := "n-" + value
		note = &n
	}
	switch draw := w.rng.Intn(10); {
	case draw < 4 && len(w.owned) > 0:
		id := w.owned[w.rng.Intn(len(w.owned))]
		_, err := pool.Exec(ctx, fmt.Sprintf("UPDATE public.%s SET name = $2, note = $3 WHERE id = $1", w.table), id, value, note)
		return err
	case draw < 7:
		var id int
		if len(w.deleted) > 0 && w.rng.Intn(2) == 0 {
			i := w.rng.Intn(len(w.deleted))
			id = w.deleted[i]
			w.deleted = append(w.deleted[:i], w.deleted[i+1:]...)
		} else {
			id = w.next
			w.next++
		}
		if _, err := pool.Exec(ctx, fmt.Sprintf("INSERT INTO public.%s (id, name, note) VALUES ($1, $2, $3)", w.table), id, value, note); err != nil {
			return err
		}
		w.owned = append(w.owned, id)
		return nil
	case len(w.owned) > 1:
		i := w.rng.Intn(len(w.owned))
		id := w.owned[i]
		if _, err := pool.Exec(ctx, fmt.Sprintf("DELETE FROM public.%s WHERE id = $1", w.table), id); err != nil {
			return err
		}
		w.owned = append(w.owned[:i], w.owned[i+1:]...)
		w.deleted = append(w.deleted, id)
	}
	return nil
}

// replayImage ist der Stand einer Zeile im Replay: die JSON-Schlüssel des
// Row Image mit ihren Text-Werten.
type replayImage map[string]string

func parseReplayImage(raw []byte) (replayImage, error) {
	image := replayImage{}
	if len(raw) == 0 {
		return image, nil
	}
	if err := json.Unmarshal(raw, &image); err != nil {
		return nil, err
	}
	return image, nil
}

// applyLog wendet das Log in Lese-Ordnung an: `INSERT` und `UPDATE` als
// Upsert des Row Images, `DELETE` als Löschen des Schlüssels.
func applyLog(rows pgx.Rows) (map[string]replayImage, error) {
	defer rows.Close()
	state := map[string]replayImage{}
	for rows.Next() {
		var operation string
		var oldData, newData []byte
		if err := rows.Scan(&operation, &oldData, &newData); err != nil {
			return nil, err
		}
		switch operation {
		case "INSERT", "UPDATE":
			image, err := parseReplayImage(newData)
			if err != nil {
				return nil, err
			}
			state[image["id"]] = image
		case "DELETE":
			image, err := parseReplayImage(oldData)
			if err != nil {
				return nil, err
			}
			delete(state, image["id"])
		default:
			return nil, fmt.Errorf("unbekannte Operation %q", operation)
		}
	}
	return state, rows.Err()
}

// sourceState liest den Quellstand: je Zeile die Spalten als Text, NULL
// entfällt (dieselbe Form wie das Row Image).
func sourceState(ctx context.Context, pool *pgxpool.Pool, table string) (map[string]replayImage, error) {
	rows, err := pool.Query(ctx, fmt.Sprintf("SELECT id::text, name, note FROM public.%s", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	state := map[string]replayImage{}
	for rows.Next() {
		var id string
		var name, note *string
		if err := rows.Scan(&id, &name, &note); err != nil {
			return nil, err
		}
		image := replayImage{"id": id}
		if name != nil {
			image["name"] = *name
		}
		if note != nil {
			image["note"] = *note
		}
		state[id] = image
	}
	return state, rows.Err()
}

// stateDiff nennt bis zu fünf Schlüssel, an denen zwei Stände abweichen.
func stateDiff(want, got map[string]replayImage) []string {
	keys := map[string]bool{}
	for k := range want {
		keys[k] = true
	}
	for k := range got {
		keys[k] = true
	}
	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	var diffs []string
	for _, k := range sorted {
		w, g := want[k], got[k]
		if fmt.Sprint(w) != fmt.Sprint(g) {
			diffs = append(diffs, fmt.Sprintf("id=%s Quelle=%v Replay=%v", k, w, g))
			if len(diffs) == 5 {
				break
			}
		}
	}
	return diffs
}

// TestE2EBackfillReplayInvariant trägt `LH-FA-CAP-009` (Boundary) am
// laufenden Feed-Container: nebenläufige `INSERT`-, `UPDATE`- und
// `DELETE`-Schreiber ändern die Tabelle vor und nach der Position des
// Snapshots, und das Log ab dem Log-Anfang, in Lese-Ordnung als Upsert und
// Löschen angewandt, ergibt je Schlüssel den Quellstand.
//
// Eine offene Schreibtransaktion hält den Run in `running` (die Anlage des
// temporären Slots wartet auf sie); vor der Freigabe committen die Schreiber
// und liegen damit auf oder vor der Snapshot-Position, danach committen sie
// dahinter. Der Test bricht ab, wenn eine der beiden Seiten leer bleibt. Die
// Lese-Ordnung `(commit_position, transaction_id, sequence)` trägt
// `LH-FA-CAP-004`.
func TestE2EBackfillReplayInvariant(t *testing.T) {
	env := newBackfillEnv(t)
	ctx := context.Background()
	const table = "feed_e2e_backfill_replay"
	const seedRows = 30

	if _, err := env.pool.Exec(ctx, fmt.Sprintf("CREATE TABLE public.%s (id int PRIMARY KEY, name text, note text)", table)); err != nil {
		t.Fatalf("Tabelle anlegen: %v", err)
	}
	for id := 1; id <= seedRows; id++ {
		var note *string
		if id%5 != 0 {
			n := fmt.Sprintf("n-%d", id)
			note = &n
		}
		if _, err := env.pool.Exec(ctx, fmt.Sprintf("INSERT INTO public.%s (id, name, note) VALUES ($1, $2, $3)", table), id, fmt.Sprintf("seed-%d", id), note); err != nil {
			t.Fatalf("Bestand einfügen: %v", err)
		}
	}
	env.enableTable(t, table)

	// Die offene Schreibtransaktion trägt eine Transaktionskennung: die
	// Slot-Anlage des Runs wartet bis zu ihrem Ende.
	hold, err := pgx.Connect(ctx, env.dsn)
	if err != nil {
		t.Fatalf("Verbindung der Haltetransaktion: %v", err)
	}
	t.Cleanup(func() { _ = hold.Close(context.Background()) })
	holdTx, err := hold.Begin(ctx)
	if err != nil {
		t.Fatalf("Haltetransaktion öffnen: %v", err)
	}
	if _, err := holdTx.Exec(ctx, "SELECT pg_current_xact_id()"); err != nil {
		t.Fatalf("Haltetransaktion: %v", err)
	}

	var runID string
	if err := env.pool.QueryRow(ctx, "SELECT cdc.backfill_table($1, 'public', $2)", e2eSource, table).Scan(&runID); err != nil {
		t.Fatalf("cdc.backfill_table: %v", err)
	}
	env.awaitRequestApplied(t, runID)
	env.awaitRunStatus(t, runID, "running", 30*time.Second)
	var snapshotPosition *int64
	if err := env.pool.QueryRow(ctx, "SELECT snapshot_position FROM cdc.backfill_run WHERE run_id = $1", runID).Scan(&snapshotPosition); err != nil {
		t.Fatalf("Run lesen: %v", err)
	}
	if snapshotPosition != nil {
		t.Fatalf("der Run trägt schon eine Snapshot-Position (%d), obwohl die Slot-Anlage auf die Haltetransaktion wartet", *snapshotPosition)
	}

	writers := make([]*replayWriter, 3)
	for w := range writers {
		writer := &replayWriter{table: table, rng: rand.New(rand.NewSource(int64(w + 1))), next: 1000 * (w + 1)}
		for id := 1; id <= seedRows; id++ {
			if id%3 == w {
				writer.owned = append(writer.owned, id)
			}
		}
		writers[w] = writer
	}
	stop := make(chan struct{})
	var wg sync.WaitGroup
	errs := make(chan error, len(writers))
	for _, writer := range writers {
		wg.Add(1)
		go func(w *replayWriter) {
			defer wg.Done()
			for counter := 0; ; counter++ {
				select {
				case <-stop:
					return
				default:
				}
				if err := w.step(ctx, env.pool, counter); err != nil {
					errs <- err
					return
				}
				time.Sleep(10 * time.Millisecond)
			}
		}(writer)
	}

	time.Sleep(1500 * time.Millisecond)
	if err := holdTx.Rollback(ctx); err != nil {
		t.Fatalf("Haltetransaktion beenden: %v", err)
	}
	env.awaitRunStatus(t, runID, "completed", 60*time.Second)
	time.Sleep(500 * time.Millisecond)
	close(stop)
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("Schreiber: %v", err)
	}

	// Ein Marker schließt die Erfassung ab: seine Änderung steht erst im Log,
	// wenn jede davor committete Transaktion persistiert ist.
	if _, err := env.pool.Exec(ctx, fmt.Sprintf("INSERT INTO public.%s (id, name) VALUES (900000, 'marker')", table)); err != nil {
		t.Fatalf("Marker einfügen: %v", err)
	}
	deadline := time.Now().Add(60 * time.Second)
	for {
		var seen int
		if err := env.pool.QueryRow(ctx,
			"SELECT count(*) FROM cdc.changes WHERE source_id = $1 AND table_name = $2 AND new_data->>'id' = '900000'", e2eSource, table,
		).Scan(&seen); err != nil {
			t.Fatalf("Marker lesen: %v", err)
		}
		if seen == 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("die Änderung des Markers wurde nicht innerhalb der Zeitspanne erfasst")
		}
		time.Sleep(200 * time.Millisecond)
	}

	var position int64
	var copied int64
	if err := env.pool.QueryRow(ctx, "SELECT snapshot_position, rows_copied FROM cdc.backfill_run WHERE run_id = $1", runID).Scan(&position, &copied); err != nil {
		t.Fatalf("Run lesen: %v", err)
	}
	rows, err := env.pool.Query(ctx, `
		SELECT operation, old_data, new_data FROM cdc.changes
		WHERE source_id = $1 AND table_name = $2
		ORDER BY commit_position, transaction_id, sequence`, e2eSource, table)
	if err != nil {
		t.Fatalf("Log lesen: %v", err)
	}
	replayed, err := applyLog(rows)
	if err != nil {
		t.Fatalf("Log anwenden: %v", err)
	}
	source, err := sourceState(ctx, env.pool, table)
	if err != nil {
		t.Fatalf("Quellstand lesen: %v", err)
	}
	if diffs := stateDiff(source, replayed); len(diffs) > 0 {
		t.Fatalf("Replay-Invariante verletzt (%d Quellzeilen, %d Replay-Zeilen): %s", len(source), len(replayed), strings.Join(diffs, "; "))
	}

	// Die Überlappung ist Voraussetzung der Aussage: Änderungen auf oder vor
	// der Snapshot-Position und dahinter.
	var backfillChanges, offPosition, walBefore, walAfter int
	if err := env.pool.QueryRow(ctx, `
		SELECT count(*) FILTER (WHERE origin = 'backfill'),
		       count(*) FILTER (WHERE origin = 'backfill' AND commit_position <> $3),
		       count(*) FILTER (WHERE origin = 'wal' AND commit_position <= $3),
		       count(*) FILTER (WHERE origin = 'wal' AND commit_position > $3)
		FROM cdc.changes WHERE source_id = $1 AND table_name = $2`, e2eSource, table, position,
	).Scan(&backfillChanges, &offPosition, &walBefore, &walAfter); err != nil {
		t.Fatalf("Zonen zählen: %v", err)
	}
	if int64(backfillChanges) != copied || backfillChanges == 0 {
		t.Fatalf("Backfill-Changes = %d, rows_copied = %d", backfillChanges, copied)
	}
	if offPosition != 0 {
		t.Fatalf("%d Backfill-Changes liegen nicht auf der Snapshot-Position %d", offPosition, position)
	}
	if walBefore == 0 || walAfter == 0 {
		t.Fatalf("Zonen ohne Überlappung: WAL-Changes auf oder vor der Snapshot-Position %d = %d, dahinter = %d", position, walBefore, walAfter)
	}
	t.Logf("Replay-Invariante: Snapshot-Position %d, %d Backfill-Changes, WAL-Changes davor %d und dahinter %d, %d Zeilen im Quellstand", position, backfillChanges, walBefore, walAfter, len(source))
}
