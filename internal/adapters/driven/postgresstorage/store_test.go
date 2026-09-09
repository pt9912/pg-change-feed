package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Adapter-Tests laufen gegen eine reale PostgreSQL-Instanz im
// Testcontainer (`ADR-0030`), gepinnt über `make test-store`; ohne DSN
// überspringen sie — die treiberfreie Seite tragen die Unit-Tests.
// Die Instanz gehört dem Container: Schema je Test frisch aufgesetzt,
// Daten bleiben im Container und landen nicht im Arbeitsbaum.

const (
	testSource      = "src-1"
	testTableMain   = "tbl-1"
	testTableOther  = "tbl-2"
	testSchemaMain  = "sv-1"
	testSchemaOther = "sv-2"
)

// newTestStore baut den Adapter gegen die Test-Instanz und setzt das
// CDC-Schema je Test frisch auf: DROP CASCADE löscht den Stand des
// vorherigen Tests, ApplySchema legt die Tabellen an (`Pflichtenheft §2`).
func newTestStore(t *testing.T) (*postgresstorage.PostgresChangeStoreAdapter, *pgxpool.Pool) {
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
	store, err := postgresstorage.New(ctx, dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(store.Close)
	return store, pool
}

// seedReference trägt Quelle, Tabellen und Schema-Versionen vor der ersten
// Persistenz (die Fremdschlüssel der DDL setzen sie voraus).
func seedReference(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	for _, statement := range []string{
		fmt.Sprintf("INSERT INTO cdc.source (source_id, name) VALUES ('%s', 'Quelle')", testSource),
		fmt.Sprintf("INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ('%s', '%s', 'public', 'a')", testTableMain, testSource),
		fmt.Sprintf("INSERT INTO cdc.source_table (source_table_id, source_id, schema_name, table_name) VALUES ('%s', '%s', 'public', 'b')", testTableOther, testSource),
		fmt.Sprintf("INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ('%s', '%s', 1)", testSchemaMain, testTableMain),
		fmt.Sprintf("INSERT INTO cdc.schema_version (schema_version_id, source_table_id, version) VALUES ('%s', '%s', 1)", testSchemaOther, testTableOther),
	} {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Referenz-Zeilen: %v", err)
		}
	}
}

// committedTransaction legt eine committed Quelltransaktion mit count
// Changes an; bei mehr als einem Change hängt der Träger die Sequenzen
// verkehrt herum an, damit die Tests die Sequenz-Ordnung gegen die
// Anhang-Reihenfolge prüfen (`LH-FA-DAT-004` Boundary).
func committedTransaction(t *testing.T, id string, offset uint64, count int, tableID model.SourceTableID, schemaVersion string) *model.ChangeTransaction {
	t.Helper()
	tx, err := model.NewOpenTransaction(model.TransactionID(id), testSource)
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	for i := 1; i <= count; i++ {
		sequence := int64(i)
		if count > 1 {
			sequence = int64(count - i + 1)
		}
		change, err := model.NewChange(
			model.ChangeID(fmt.Sprintf("%s-%d", id, sequence)),
			model.TransactionID(id),
			tableID,
			sequence,
			model.OperationInsert,
			nil,
			[]byte(fmt.Sprintf(`{"n":%d}`, sequence)),
			model.SchemaVersionID(schemaVersion),
		)
		if err != nil {
			t.Fatalf("NewChange %d: %v", sequence, err)
		}
		if err := tx.AppendChange(change); err != nil {
			t.Fatalf("AppendChange %d: %v", sequence, err)
		}
	}
	position, err := model.NewSourcePosition(testSource, offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	if err := tx.Commit(position); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	return tx
}

func position(t *testing.T, offset uint64) *model.SourcePosition {
	t.Helper()
	p, err := model.NewSourcePosition(testSource, offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	return &p
}

func query(t *testing.T, start, end *model.SourcePosition, table *model.SourceTableID, limit *int) outbound.ChangeQuery {
	t.Helper()
	q := outbound.ChangeQuery{Source: testSource, Start: start, End: end, Table: table, Limit: limit}
	if err := q.Validate(); err != nil {
		t.Fatalf("ChangeQuery: %v", err)
	}
	return q
}

func tableFilter(t *testing.T, id string) *model.SourceTableID {
	t.Helper()
	filter := model.SourceTableID(id)
	return &filter
}

func limit(t *testing.T, n int) *int {
	t.Helper()
	return &n
}

// persist führt PersistTransaction über alle Transaktionen der Tests; ein
// Fehler bricht den Test ab.
func persist(t *testing.T, store *postgresstorage.PostgresChangeStoreAdapter, transactions ...*model.ChangeTransaction) {
	t.Helper()
	for _, tx := range transactions {
		if err := store.PersistTransaction(context.Background(), tx); err != nil {
			t.Fatalf("PersistTransaction %s: %v", tx.ID, err)
		}
	}
}

// readRecords liest und fasst die Records als Textbeleg (Position, Change)
// auf, damit die Tests über die volle Ordnung urteilen.
func readRecords(t *testing.T, store *postgresstorage.PostgresChangeStoreAdapter, q outbound.ChangeQuery) []string {
	t.Helper()
	records, err := store.ReadChanges(context.Background(), q)
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	lines := make([]string, 0, len(records))
	for _, record := range records {
		lines = append(lines, fmt.Sprintf("%s@%d:INSERT#%d", record.Change.ID, record.Position.Offset, record.Change.Sequence))
	}
	return lines
}

// count trägt eine Tabellen-Zahl über die Test-Verbindung; die
// Deduplizierungs-Behauptungen zählen über sie.
func count(t *testing.T, pool *pgxpool.Pool, statement string, args ...any) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), statement, args...).Scan(&n); err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

// Idempotenz (`ADR-0011`, `LH-QA-REL-001.a`): dieselbe committed
// Transaktion erneut persistiert — der Crash zwischen Persistenz und ACK
// wiederholt sie — ändert keinen Stand und meldet keinen Fehler; die
// Deduplizierungsbasis trägt die interne Transaktions-ID.
func TestPersistTransactionIsIdempotent(t *testing.T) {
	store, pool := newTestStore(t)
	seedReference(t, pool)
	ctx := context.Background()
	tx := committedTransaction(t, "t-1", 100, 3, testTableMain, testSchemaMain)

	if err := store.PersistTransaction(ctx, tx); err != nil {
		t.Fatalf("erste PersistTransaction: %v", err)
	}
	if err := store.PersistTransaction(ctx, tx); err != nil {
		t.Fatalf("erneute PersistTransaction: %v (Idempotenz meldet keinen Fehler)", err)
	}
	if got := count(t, pool, "SELECT count(*) FROM cdc.transaction WHERE source_id = $1", testSource); got != 1 {
		t.Fatalf("Transaktions-Zeilen = %d, wollen 1", got)
	}
	if got := count(t, pool, "SELECT count(*) FROM cdc.change WHERE transaction_id = 't-1'"); got != 3 {
		t.Fatalf("Change-Zeilen = %d, wollen 3", got)
	}
	records := readRecords(t, store, query(t, nil, nil, nil, nil))
	if len(records) != 3 {
		t.Fatalf("Lesen trägt %d Changes, wollen 3 (keine Duplikate)", len(records))
	}
}

// Offene Transaktionen sind nicht konsumierbar (`ADR-0029`, Regel 3); der
// Adapter persistiert sie nicht und trägt keinen Stand.
func TestPersistRejectsOpenTransaction(t *testing.T) {
	store, pool := newTestStore(t)
	seedReference(t, pool)
	tx, err := model.NewOpenTransaction("t-open", testSource)
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}

	if err := store.PersistTransaction(context.Background(), tx); !stderrors.Is(err, domainerrors.ErrTransactionNotCommitted) {
		t.Fatalf("Fehler = %v, wollen %v", err, domainerrors.ErrTransactionNotCommitted)
	}
	if got := count(t, pool, "SELECT count(*) FROM cdc.transaction"); got != 0 {
		t.Fatalf("Transaktions-Zeilen = %d, wollen 0", got)
	}
}

// Eine committed Transaktion ohne Changes persistiert ihre Commit-Position
// (`LH-FA-CAP-006.a`); das Lesen trägt keinen Change.
func TestPersistCarriesEmptyCommittedTransaction(t *testing.T) {
	store, pool := newTestStore(t)
	seedReference(t, pool)
	tx := committedTransaction(t, "t-empty", 50, 0, testTableMain, testSchemaMain)

	if err := store.PersistTransaction(context.Background(), tx); err != nil {
		t.Fatalf("PersistTransaction: %v", err)
	}
	if got := count(t, pool, "SELECT count(*) FROM cdc.transaction WHERE commit_position = 50"); got != 1 {
		t.Fatalf("Transaktions-Zeilen an Position 50 = %d, wollen 1", got)
	}
	if records := readRecords(t, store, query(t, nil, nil, nil, nil)); len(records) != 0 {
		t.Fatalf("Lesen trägt %v, wollen leer", records)
	}
}

// Deterministische Ordnung (`LH-FA-REA-004.a`): die Sortierung läuft über
// (Commit-Position, Transaktions-ID, Sequenz) — quer über die
// Persistierungs-Reihenfolge, mit Transaktions-ID als Tiebreak bei gleicher
// Position, unabhängig von der Anhang-Reihenfolge innerhalb der
// Transaktion. Dieselbe Abfrage trägt denselben Stand in zwei Läufen
// (`LH-FA-REA-004` Happy Path).
func TestReadIsDeterministicallySorted(t *testing.T) {
	store, pool := newTestStore(t)
	seedReference(t, pool)

	// Persistierungs-Reihenfolge absichtlich verkehrt herum: die Ordnung
	// trägt die Commit-Position, nicht die Einfüge-Reihenfolge.
	first := committedTransaction(t, "t-zeta", 300, 2, testTableMain, testSchemaMain)
	second := committedTransaction(t, "t-alpha", 300, 1, testTableMain, testSchemaMain)
	third := committedTransaction(t, "t-mid", 200, 2, testTableMain, testSchemaMain)
	persist(t, store, first, second, third)

	q := query(t, nil, nil, nil, nil)
	got := readRecords(t, store, q)
	want := []string{
		// Position 200: beide Changes von t-mid in Sequenz-Reihenfolge.
		"t-mid-1@200:INSERT#1",
		"t-mid-2@200:INSERT#2",
		// Position 300: Transaktions-ID als Tiebreak — t-alpha vor t-zeta.
		"t-alpha-1@300:INSERT#1",
		"t-zeta-1@300:INSERT#1",
		"t-zeta-2@300:INSERT#2",
	}
	if len(got) != len(want) {
		t.Fatalf("Ordnung trägt %d Zeilen, wollen %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Ordnung an Stelle %d = %q, wollen %q (voll: %v)", i, got[i], want[i], got)
		}
	}
	again := readRecords(t, store, q)
	for i := range got {
		if got[i] != again[i] {
			t.Fatalf("Zweiter Lauf weicht an Stelle %d ab: %q gegen %q", i, again[i], got[i])
		}
	}
}

// Lesen ab bekannter Position (`LH-FA-REA-002`) und Bereich
// (`LH-FA-REA-001`): `[p1, p2)` ist start-inklusiv und end-exklusiv, ein
// leerer Bereich trägt eine leere Menge, ein invertierter Bereich endet am
// Port-Kontrakt (Negative).
func TestReadCarriesPositionsAndRanges(t *testing.T) {
	store, pool := newTestStore(t)
	seedReference(t, pool)
	var transactions []*model.ChangeTransaction
	for _, item := range []struct {
		id     string
		offset uint64
	}{
		{"t-1", 100},
		{"t-2", 200},
		{"t-3", 300},
	} {
		transactions = append(transactions, committedTransaction(t, item.id, item.offset, 1, testTableMain, testSchemaMain))
	}
	persist(t, store, transactions...)

	got := readRecords(t, store, query(t, position(t, 100), position(t, 200), nil, nil))
	if len(got) != 1 || got[0] != "t-1-1@100:INSERT#1" {
		t.Fatalf("Bereich [100,200) = %v, wollen genau t-1", got)
	}
	got = readRecords(t, store, query(t, position(t, 150), position(t, 200), nil, nil))
	if len(got) != 0 {
		t.Fatalf("Leerer Bereich [150,200) = %v, wollen leer", got)
	}
	got = readRecords(t, store, query(t, position(t, 200), nil, nil, nil))
	if len(got) != 2 {
		t.Fatalf("Lesen ab 200 = %v, wollen t-2 und t-3", got)
	}
	got = readRecords(t, store, query(t, position(t, 500), nil, nil, nil))
	if len(got) != 0 {
		t.Fatalf("Lesen ab 500 = %v, wollen leer", got)
	}
	if _, err := store.ReadChanges(context.Background(), outbound.ChangeQuery{
		Source: testSource,
		Start:  position(t, 300),
		End:    position(t, 100),
	}); !stderrors.Is(err, outbound.ErrRangeInverted) {
		t.Fatalf("Invertierter Bereich: Fehler = %v, wollen %v", err, outbound.ErrRangeInverted)
	}
}

// Lese-Limit (`LH-FA-REA-003`): höchstens `n` Changes in Ordnung, bei
// weniger als `n` alle verfügbaren; ein Limit unter 1 endet am
// Port-Kontrakt.
func TestReadCarriesLimit(t *testing.T) {
	store, pool := newTestStore(t)
	seedReference(t, pool)
	var transactions []*model.ChangeTransaction
	for i, id := range []string{"t-1", "t-2", "t-3", "t-4", "t-5"} {
		transactions = append(transactions, committedTransaction(t, id, uint64(100*(i+1)), 1, testTableMain, testSchemaMain))
	}
	persist(t, store, transactions...)

	got := readRecords(t, store, query(t, nil, nil, nil, limit(t, 2)))
	if len(got) != 2 || got[0] != "t-1-1@100:INSERT#1" || got[1] != "t-2-1@200:INSERT#1" {
		t.Fatalf("Limit 2 = %v, wollen die zwei ersten in Ordnung", got)
	}
	got = readRecords(t, store, query(t, nil, nil, nil, limit(t, 50)))
	if len(got) != 5 {
		t.Fatalf("Limit 50 trägt %d, wollen alle 5 verfügbaren", len(got))
	}
	if _, err := store.ReadChanges(context.Background(), outbound.ChangeQuery{
		Source: testSource,
		Limit:  limit(t, 0),
	}); !stderrors.Is(err, outbound.ErrNonPositiveLimit) {
		t.Fatalf("Limit 0: Fehler = %v, wollen %v", err, outbound.ErrNonPositiveLimit)
	}
}

// Tabellenfilter (`LH-FA-REA-006`): nur Changes der gefilterten Tabelle;
// ein Filter ohne Treffer trägt eine leere Menge.
func TestReadCarriesTableFilter(t *testing.T) {
	store, pool := newTestStore(t)
	seedReference(t, pool)
	first := committedTransaction(t, "t-1", 100, 2, testTableMain, testSchemaMain)
	second := committedTransaction(t, "t-2", 200, 1, testTableOther, testSchemaOther)
	persist(t, store, first, second)

	got := readRecords(t, store, query(t, nil, nil, tableFilter(t, testTableMain), nil))
	if len(got) != 2 {
		t.Fatalf("Filter tbl-1 = %v, wollen nur seine 2 Changes", got)
	}
	got = readRecords(t, store, query(t, nil, nil, tableFilter(t, "tbl-fehlt"), nil))
	if len(got) != 0 {
		t.Fatalf("Filter ohne Treffer = %v, wollen leer", got)
	}
}

// Erneutes Lesen innerhalb der Aufbewahrung (`LH-FA-REA-005`) und Lesen
// ohne Positionsänderung (`LH-FA-REA-002`): dieselbe Abfrage trägt
// denselben Stand, der gespeicherte Stand bleibt unverändert.
func TestReadLeavesPersistedStateUnchanged(t *testing.T) {
	store, pool := newTestStore(t)
	seedReference(t, pool)
	persist(t, store, committedTransaction(t, "t-1", 100, 3, testTableMain, testSchemaMain))
	q := query(t, nil, nil, nil, limit(t, 2))

	first := readRecords(t, store, q)
	second := readRecords(t, store, q)
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("Erneutes Lesen weicht an Stelle %d ab: %q gegen %q", i, second[i], first[i])
		}
	}
	if got := count(t, pool, "SELECT count(*) FROM cdc.change WHERE transaction_id = 't-1'"); got != 3 {
		t.Fatalf("Change-Zeilen nach dem Lesen = %d, wollen 3", got)
	}
	if got := count(t, pool, "SELECT count(*) FROM cdc.transaction WHERE commit_position = 100"); got != 1 {
		t.Fatalf("Transaktions-Zeilen nach dem Lesen = %d, wollen 1", got)
	}
}

// Treiber-Fehler tragen die Klasse `storage` (`outbound.ErrStorage`,
// `SPEC-008`, `ADR-0023`); Application und Betrieb klassifizieren über
// errors.Is und kennen keinen Treibertyp. Die technische Ursache bleibt
// über errors.Is hinter der Klasse lesbar.

// Ein Persistenzfehler (hier: Fremdschlüssel-Verstoß gegen die nicht
// registrierte Schema-Version) trägt die Klasse `storage` und die
// Treiber-Ursache in derselben Meldung.
func TestPersistCarriesStorageClass(t *testing.T) {
	store, pool := newTestStore(t)
	seedReference(t, pool)
	tx, err := model.NewOpenTransaction("t-fk", testSource)
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	change, err := model.NewChange("c-fk", "t-fk", testTableMain, 1, model.OperationInsert, nil, []byte(`{"n":1}`), "sv-fehlt")
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	if err := tx.AppendChange(change); err != nil {
		t.Fatalf("AppendChange: %v", err)
	}
	position, err := model.NewSourcePosition(testSource, 100)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	if err := tx.Commit(position); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	err = store.PersistTransaction(context.Background(), tx)
	if !stderrors.Is(err, outbound.ErrStorage) {
		t.Fatalf("Fehler = %v, wollen Klasse storage (%v)", err, outbound.ErrStorage)
	}
	if !strings.Contains(err.Error(), "SQLSTATE") {
		t.Fatalf("Meldung = %q, wollen die Treiber-Ursache hinter der Klasse", err.Error())
	}
}

// Eine nicht erreichbare Instanz meldet der Aufbau als Fehler der Klasse
// `storage`; der Test braucht keine Datenbank (Port 1 verwirft lokal).
func TestNewCarriesStorageClass(t *testing.T) {
	_, err := postgresstorage.New(context.Background(), "postgres://cdc:cdc@127.0.0.1:1/cdc_test?sslmode=disable")
	if !stderrors.Is(err, outbound.ErrStorage) {
		t.Fatalf("Fehler = %v, wollen Klasse storage (%v)", err, outbound.ErrStorage)
	}
}
