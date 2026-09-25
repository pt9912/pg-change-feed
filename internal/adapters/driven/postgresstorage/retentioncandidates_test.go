package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Tests dieser Datei führen `ReadRetentionCandidates` (`ADR-0124`)
// gegen die reale PostgreSQL des Testcontainers auf dem ausgerollten Schema:
// die Seiten einer Quelle, Position und Zeitpunkt je Kandidat, die Lage
// einer Zeile hinter und vor dem Cursor, die Eingabe-Grenzen und der Zugriff
// unter den Login-Identitäten der Rollen. Sie laufen vor den Tests, die das
// Schema über `DROP SCHEMA cdc CASCADE` neu aufsetzen (Dateiname-Ordnung).

const (
	retentionSourceA = "ret-cand-a"
	retentionSourceB = "ret-cand-b"
)

// seedRetentionSource legt Quelle, Tabelle und Schema-Version einer Quelle an
// und räumt nach dem Test alle Zeilen der Quelle ab.
func seedRetentionSource(f *backfillFixture, source string) {
	f.t.Helper()
	f.exec("INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Retention-Kandidaten') ON CONFLICT (source_id) DO NOTHING", source)
	f.exec("INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ($1, $2, 'public', 'ret_cand')", source+"-tbl", source)
	f.exec("INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ($1, $2, 1)", source+"-sv", source+"-tbl")
	f.t.Cleanup(func() {
		ctx := context.Background()
		for _, statement := range []string{
			"DELETE FROM cdc.change WHERE transaction_id IN (SELECT transaction_id FROM cdc.transaction WHERE source_id = $1)",
			"DELETE FROM cdc.transaction WHERE source_id = $1",
			"DELETE FROM cdc.schema_version WHERE source_table_id = $1::text || '-tbl'",
			"DELETE FROM cdc.source_table WHERE source_id = $1",
			"DELETE FROM cdc.source WHERE source_id = $1",
		} {
			_, _ = f.pool.Exec(ctx, statement, source)
		}
	})
}

// addRetentionTransaction legt eine Transaktion mit den übergebenen
// Change-Kennungen an; Kennungen mit dem Präfix `0bf-` tragen die Herkunft
// `backfill`, alle anderen keine.
func addRetentionTransaction(f *backfillFixture, source, transactionID string, position int64, committedAt time.Time, changeIDs ...string) {
	f.t.Helper()
	f.exec("INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at) VALUES ($1, $2, $3, $4)",
		transactionID, source, position, committedAt)
	for i, id := range changeIDs {
		var origin any
		if strings.HasPrefix(id, "0bf-") {
			origin = "backfill"
		}
		f.exec(`INSERT INTO cdc.change
		    (change_id, transaction_id, source_table_id, sequence, operation, old_data, new_data, schema_version, origin)
		    VALUES ($1, $2, $3, $4, 'INSERT', NULL, '{}'::jsonb, $5, $6)`,
			id, transactionID, source+"-tbl", i+1, source+"-sv", origin)
	}
}

// candidatePage liest eine Seite und bricht den Test bei einem Fehler ab.
func candidatePage(t *testing.T, store *postgresstorage.PostgresChangeStoreAdapter, source string, after model.ChangeID, limit int) []outbound.RetentionCandidate {
	t.Helper()
	page, err := store.ReadRetentionCandidates(context.Background(), model.SourceID(source), after, limit)
	if err != nil {
		t.Fatalf("ReadRetentionCandidates(after=%q, limit=%d): %v", after, limit, err)
	}
	return page
}

// walkCandidates liest alle Seiten bis zur leeren Seite; die Rückgabe trägt
// die nichtleeren Seiten.
func walkCandidates(t *testing.T, store *postgresstorage.PostgresChangeStoreAdapter, source string, limit int) [][]outbound.RetentionCandidate {
	t.Helper()
	var pages [][]outbound.RetentionCandidate
	var after model.ChangeID
	for i := 0; i < 1000; i++ {
		page := candidatePage(t, store, source, after, limit)
		if len(page) == 0 {
			return pages
		}
		pages = append(pages, page)
		after = page[len(page)-1].ChangeID
	}
	t.Fatalf("der Durchlauf endet nicht nach 1000 Seiten (limit=%d)", limit)
	return nil
}

func candidateIDs(pages [][]outbound.RetentionCandidate) []model.ChangeID {
	var ids []model.ChangeID
	for _, page := range pages {
		for _, candidate := range page {
			ids = append(ids, candidate.ChangeID)
		}
	}
	return ids
}

// dbOrderedIDs liest die Kennungen einer Quelle über eine unabhängige
// Anweisung in der Ordnung der Datenbank.
func dbOrderedIDs(f *backfillFixture, source string) []model.ChangeID {
	f.t.Helper()
	rows, err := f.pool.Query(context.Background(),
		`SELECT c.change_id FROM cdc.change c JOIN cdc.transaction t ON t.transaction_id = c.transaction_id
		 WHERE t.source_id = $1 ORDER BY c.change_id`, source)
	if err != nil {
		f.t.Fatalf("Kennungen der Quelle: %v", err)
	}
	defer rows.Close()
	var ids []model.ChangeID
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			f.t.Fatalf("Kennungen der Quelle: %v", err)
		}
		ids = append(ids, model.ChangeID(id))
	}
	return ids
}

func sameIDs(got, want []model.ChangeID) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func newRetentionStore(t *testing.T, dsn string) *postgresstorage.PostgresChangeStoreAdapter {
	t.Helper()
	store, err := postgresstorage.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(store.Close)
	return store
}

// TestRetentionCandidatePagesCoverExactlyTheSourceChanges belegt den
// Seitenvertrag (`ADR-0124` Festlegung 2): bei Seiten zu 1, 2, 5 und 100
// decken die Seiten einer Quelle genau ihre fünf Changes ab — keine Kennung
// doppelt, keine der Fremdquelle, keine Seite über `limit`, aufsteigend in
// der Ordnung des Schlüssels (verglichen mit einer unabhängigen Abfrage).
// Die Zeilen der Quelle stehen in ungeordneter Einfügereihenfolge, damit
// eine fehlende Sortierung die Seitenfolge bricht. Rot färbende Mutationen
// der Eingabeseite: `c.change_id >= $2` (die letzte Kennung kommt doppelt),
// `t.source_id = $1` entfernt (Fremdquelle in den Seiten), `LIMIT $3`
// entfernt (eine Seite liefert alles), `ORDER BY c.change_id` entfernt (der
// Cursor überspringt oder wiederholt Kennungen).
func TestRetentionCandidatePagesCoverExactlyTheSourceChanges(t *testing.T) {
	f := newBackfillFixture(t)
	seedRetentionSource(f, retentionSourceA)
	seedRetentionSource(f, retentionSourceB)
	at := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	addRetentionTransaction(f, retentionSourceA, "ret-a-tx9", 90, at, "tx9-1")
	addRetentionTransaction(f, retentionSourceB, "ret-b-tx4", 40, at, "tx4-b1")
	addRetentionTransaction(f, retentionSourceA, "ret-a-bf2", 1000, at, "0bf-run-2")
	addRetentionTransaction(f, retentionSourceA, "ret-a-tx1", 10, at, "tx1-1")
	addRetentionTransaction(f, retentionSourceB, "ret-b-tx7", 70, at, "tx7-b1", "tx7-b2")
	addRetentionTransaction(f, retentionSourceA, "ret-a-tx5", 50, at, "tx5-1")
	addRetentionTransaction(f, retentionSourceA, "ret-a-bf1", 1000, at, "0bf-run-1")
	store := newRetentionStore(t, f.dsn)

	want := map[model.ChangeID]bool{"0bf-run-1": true, "0bf-run-2": true, "tx1-1": true, "tx5-1": true, "tx9-1": true}
	independent := dbOrderedIDs(f, retentionSourceA)
	if len(independent) != len(want) {
		t.Fatalf("unabhängige Abfrage nennt %d Kennungen, wollen %d: %v", len(independent), len(want), independent)
	}

	for _, limit := range []int{1, 2, 5, 100} {
		pages := walkCandidates(t, store, retentionSourceA, limit)
		for i, page := range pages {
			if len(page) > limit {
				t.Fatalf("limit=%d: Seite %d trägt %d Kandidaten", limit, i+1, len(page))
			}
		}
		got := candidateIDs(pages)
		seen := map[model.ChangeID]bool{}
		for _, id := range got {
			if seen[id] {
				t.Fatalf("limit=%d: Kennung %q doppelt in %v", limit, id, got)
			}
			seen[id] = true
			if !want[id] {
				t.Fatalf("limit=%d: Kennung %q gehört nicht zur Quelle: %v", limit, id, got)
			}
		}
		if len(got) != len(want) {
			t.Fatalf("limit=%d: %d Kandidaten, wollen %d: %v", limit, len(got), len(want), got)
		}
		if !sameIDs(got, independent) {
			t.Fatalf("limit=%d: Folge %v, unabhängige Ordnung %v", limit, got, independent)
		}
	}
}

// TestRetentionCandidateCarriesPositionAndTimeOfReadChanges belegt Position
// und Commit-Zeitpunkt je Kandidat: sie sind gleich denen des Datensatzes von
// `ReadChanges` für dieselbe Kennung — für eine Backfill-Kennung (`0bf-…`) und
// eine WAL-Kennung an derselben Position — und gleich den eingefügten Werten.
// Rot färbende Mutation: die Projektion liest `t.commit_position` und
// `t.committed_at` vertauscht oder aus einer anderen Spalte.
func TestRetentionCandidateCarriesPositionAndTimeOfReadChanges(t *testing.T) {
	f := newBackfillFixture(t)
	seedRetentionSource(f, retentionSourceA)
	backfillAt := time.Date(2026, 9, 1, 8, 30, 15, 123456000, time.UTC)
	walAt := time.Date(2026, 9, 1, 9, 45, 0, 654321000, time.UTC)
	addRetentionTransaction(f, retentionSourceA, "ret-a-bf", 1000, backfillAt, "0bf-run-1", "0bf-run-2")
	addRetentionTransaction(f, retentionSourceA, "ret-a-wal", 1000, walAt, "wal-tx-1")
	store := newRetentionStore(t, f.dsn)

	candidates := candidateIDs([][]outbound.RetentionCandidate{candidatePage(t, store, retentionSourceA, "", 100)})
	if len(candidates) != 3 {
		t.Fatalf("%d Kandidaten, wollen 3: %v", len(candidates), candidates)
	}
	wantTime := map[model.ChangeID]int64{
		"0bf-run-1": backfillAt.UnixNano(),
		"0bf-run-2": backfillAt.UnixNano(),
		"wal-tx-1":  walAt.UnixNano(),
	}
	records, err := store.ReadChanges(context.Background(), outbound.ChangeQuery{Source: retentionSourceA})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	byID := map[model.ChangeID]outbound.ChangeRecord{}
	for _, record := range records {
		byID[record.Change.ID] = record
	}
	for _, candidate := range candidatePage(t, store, retentionSourceA, "", 100) {
		record, ok := byID[candidate.ChangeID]
		if !ok {
			t.Fatalf("Kandidat %q fehlt in ReadChanges", candidate.ChangeID)
		}
		if candidate.Position != record.Position || candidate.CommittedAt != record.CommittedAt {
			t.Fatalf("%s: Kandidat %+v, ReadChanges %+v/%+v", candidate.ChangeID, candidate, record.Position, record.CommittedAt)
		}
		if candidate.Position.SourceID != retentionSourceA || candidate.Position.Offset != 1000 {
			t.Fatalf("%s: Position %+v, wollen %s/1000", candidate.ChangeID, candidate.Position, retentionSourceA)
		}
		if candidate.CommittedAt.UnixNanos != wantTime[candidate.ChangeID] {
			t.Fatalf("%s: Zeitpunkt %d, wollen %d", candidate.ChangeID, candidate.CommittedAt.UnixNanos, wantTime[candidate.ChangeID])
		}
	}
}

// TestRetentionCandidateBehindCursorAppearsInNextRun belegt die
// Seitengrenze (`ADR-0124` Festlegung 6): eine Zeile, die zwischen zwei
// Seiten **hinter** dem Cursor (Kennung nicht größer als `after`) committet,
// fehlt im laufenden Durchlauf und steht im nächsten; eine Zeile **vor** dem
// Cursor (Kennung größer als `after`) steht im laufenden. Rot färbende
// Mutation: der Cursor ignoriert `after` (Kennungen wiederholen sich) oder
// wird `>=` (die Kennung am Cursor wiederholt sich).
func TestRetentionCandidateBehindCursorAppearsInNextRun(t *testing.T) {
	f := newBackfillFixture(t)
	seedRetentionSource(f, retentionSourceA)
	at := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	addRetentionTransaction(f, retentionSourceA, "ret-a-tx2", 20, at, "c-20")
	addRetentionTransaction(f, retentionSourceA, "ret-a-tx4", 40, at, "c-40")
	store := newRetentionStore(t, f.dsn)

	first := candidatePage(t, store, retentionSourceA, "", 1)
	if len(first) != 1 || first[0].ChangeID != "c-20" {
		t.Fatalf("erste Seite = %+v, wollen [c-20]", first)
	}
	addRetentionTransaction(f, retentionSourceA, "ret-a-tx1", 10, at, "c-10") // hinter dem Cursor
	addRetentionTransaction(f, retentionSourceA, "ret-a-tx3", 30, at, "c-30") // vor dem Cursor

	var running []model.ChangeID
	after := first[0].ChangeID
	running = append(running, after)
	for i := 0; i < 10; i++ {
		page := candidatePage(t, store, retentionSourceA, after, 1)
		if len(page) == 0 {
			break
		}
		running = append(running, page[0].ChangeID)
		after = page[0].ChangeID
	}
	if want := []model.ChangeID{"c-20", "c-30", "c-40"}; !sameIDs(running, want) {
		t.Fatalf("laufender Durchlauf = %v, wollen %v", running, want)
	}

	next := candidateIDs(walkCandidates(t, store, retentionSourceA, 1))
	if want := []model.ChangeID{"c-10", "c-20", "c-30", "c-40"}; !sameIDs(next, want) {
		t.Fatalf("nächster Durchlauf = %v, wollen %v", next, want)
	}
}

// TestRetentionCandidatesRejectInvalidInput belegt die Eingabe-Grenzen am
// Adapter: `limit` kleiner 1 endet als `ErrNonPositiveLimit`, eine leere
// Quelle als `ErrEmptyIdentifier` — kein leerer Erfolg. Rot färbende
// Mutation: eine der beiden Prüfungen in `ReadRetentionCandidates`
// entfernen.
func TestRetentionCandidatesRejectInvalidInput(t *testing.T) {
	f := newBackfillFixture(t)
	store := newRetentionStore(t, f.dsn)
	ctx := context.Background()

	for _, limit := range []int{0, -1} {
		if _, err := store.ReadRetentionCandidates(ctx, retentionSourceA, "", limit); !stderrors.Is(err, outbound.ErrNonPositiveLimit) {
			t.Fatalf("limit=%d: Fehler %v, wollen ErrNonPositiveLimit", limit, err)
		}
	}
	if _, err := store.ReadRetentionCandidates(ctx, "", "", 10); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Quelle: Fehler %v, wollen ErrEmptyIdentifier", err)
	}
}

// TestRetentionCandidatesRunUnderTheirRoles belegt den Zugriff unter den
// Login-Identitäten (`ADR-0053`, `ADR-0047`): die Login-Identität in der
// Gruppenrolle `cdc_admin` — ohne Superuser-Recht — liest die Seiten; die
// Identitäten in `cdc_capture` und `cdc_reader` scheitern am Lesen mit
// SQLSTATE 42501. Rot färbende Mutation: `SELECT` an `cdc_admin` auf
// `cdc.transaction` oder `cdc.change` in `tools/schema/nacharbeit-roles.sql`
// streichen — der Lesezugriff unter `cdc_admin` endet mit SQLSTATE 42501.
func TestRetentionCandidatesRunUnderTheirRoles(t *testing.T) {
	f := newBackfillFixture(t)
	seedRetentionSource(f, retentionSourceA)
	at := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	addRetentionTransaction(f, retentionSourceA, "ret-a-tx1", 10, at, "c-10", "c-11")
	ctx := context.Background()

	adminStore := newRetentionStore(t, f.login("pgc_test_retcand_admin", "cdc_admin"))
	page, err := adminStore.ReadRetentionCandidates(ctx, retentionSourceA, "", 10)
	if err != nil {
		t.Fatalf("ReadRetentionCandidates unter cdc_admin: %v", err)
	}
	if got := candidateIDs([][]outbound.RetentionCandidate{page}); !sameIDs(got, []model.ChangeID{"c-10", "c-11"}) {
		t.Fatalf("unter cdc_admin gelesen: %v, wollen [c-10 c-11]", got)
	}

	for _, group := range []string{"cdc_capture", "cdc_reader"} {
		store := newRetentionStore(t, f.login("pgc_test_retcand_"+group, group))
		_, err := store.ReadRetentionCandidates(ctx, retentionSourceA, "", 10)
		if err == nil || !strings.Contains(err.Error(), "SQLSTATE 42501") {
			t.Fatalf("ReadRetentionCandidates unter %s: erwartet SQLSTATE 42501 (insufficient_privilege), erhalten %v", group, err)
		}
		if !stderrors.Is(err, outbound.ErrStorage) {
			t.Fatalf("ReadRetentionCandidates unter %s: Fehler nicht in der Klasse storage: %v", group, err)
		}
	}
}
