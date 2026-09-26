// Die Tests dieses Pakets laufen netzlos: sie fahren die Verklebung der
// Zeilen-Übersetzung gegen einen Träger, der die Naht erfüllt — den
// Query-/Exec-Aufruf, die Scan-Schleife, den Leerfall und den Fehlerpfad.
// Sie sind ausdrücklich **kein** Ersatz der realen Datenbank-Tests: das SQL
// prüft weiterhin der PostgreSQL der Adapter-Tests (`make test-store`,
// `make test-replication`).
package sqlexec_test

import (
	"context"
	stderrors "errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/sqlexec"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// --- Fakes der Naht ---

// call trägt einen beobachteten Aufruf des Trägers: die Anweisung und ihre
// Argumente — der Beleg, dass die Übersetzung den Query-/Exec-Aufruf
// tatsächlich absetzt.
type call struct {
	sql  string
	args []any
}

// fakeRows erfüllt `pgx.Rows` aus Werten im Speicher. `scanErrs` trägt je
// Zeilenindex einen Scan-Fehler, `iterErr` den Iterations-Fehler
// (`rows.Err()`).
type fakeRows struct {
	rows     [][]any
	scanErrs map[int]error
	iterErr  error
	index    int
	closed   bool
}

func (r *fakeRows) Next() bool {
	if r.index < len(r.rows) {
		r.index++
		return true
	}
	return false
}

func (r *fakeRows) Scan(dest ...any) error {
	if err, ok := r.scanErrs[r.index-1]; ok {
		return err
	}
	if r.index == 0 || r.index > len(r.rows) {
		return fmt.Errorf("fake: Scan ohne Zeile")
	}
	row := r.rows[r.index-1]
	if len(dest) != len(row) {
		return fmt.Errorf("fake: %d Ziele für %d Werte", len(dest), len(row))
	}
	for i := range dest {
		if err := assign(dest[i], row[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *fakeRows) Err() error                    { return r.iterErr }
func (r *fakeRows) Close()                        { r.closed = true }
func (r *fakeRows) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }
func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}
func (r *fakeRows) Values() ([]any, error) { return nil, nil }
func (r *fakeRows) RawValues() [][]byte    { return nil }
func (r *fakeRows) Conn() *pgx.Conn        { return nil }
func (r *fakeRows) TypeMap() *pgtype.Map   { return nil }

// fakeRow erfüllt `pgx.Row` aus Werten im Speicher oder aus einem
// vorgegebenen Fehler (etwa `pgx.ErrNoRows`).
type fakeRow struct {
	values []any
	err    error
}

func (r *fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return fmt.Errorf("fake: %d Ziele für %d Werte", len(dest), len(r.values))
	}
	for i := range dest {
		if err := assign(dest[i], r.values[i]); err != nil {
			return err
		}
	}
	return nil
}

// fakeExecutor erfüllt `sqlexec.Executor` und nimmt jeden Aufruf auf.
type fakeExecutor struct {
	rows     *fakeRows
	queryErr error
	row      *fakeRow
	tag      pgconn.CommandTag
	execErr  error

	queries    []call
	queryRows  []call
	executions []call
	closed     bool
}

func (e *fakeExecutor) Query(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
	e.queries = append(e.queries, call{sql: sql, args: args})
	if e.queryErr != nil {
		return nil, e.queryErr
	}
	return e.rows, nil
}

func (e *fakeExecutor) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	e.queryRows = append(e.queryRows, call{sql: sql, args: args})
	return e.row
}

func (e *fakeExecutor) Exec(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	e.executions = append(e.executions, call{sql: sql, args: args})
	if e.execErr != nil {
		return pgconn.CommandTag{}, e.execErr
	}
	return e.tag, nil
}

func (e *fakeExecutor) Begin(context.Context) (pgx.Tx, error) {
	return nil, stderrors.New("fake: kein Transaktions-Träger in dieser Naht")
}

func (e *fakeExecutor) Close() { e.closed = true }

// assign trägt einen Wert in das Ziel eines Scan-Aufrufs — das Verhalten des
// Treibers, das die Fakes nachbilden.
func assign(dest any, value any) error {
	target := reflect.ValueOf(dest)
	if target.Kind() != reflect.Pointer || target.IsNil() {
		return fmt.Errorf("fake: Scan-Ziel ohne Zeiger")
	}
	if value == nil {
		target.Elem().Set(reflect.Zero(target.Elem().Type()))
		return nil
	}
	source := reflect.ValueOf(value)
	if !source.Type().AssignableTo(target.Elem().Type()) {
		return fmt.Errorf("fake: %s nicht nach %s zuweisbar", source.Type(), target.Elem().Type())
	}
	target.Elem().Set(source)
	return nil
}

// --- Zusicherungen der Naht selbst ---

func TestClassifyCarriesClassAndCause(t *testing.T) {
	class := stderrors.New("Klasse storage")
	cause := stderrors.New("Verbindung abgelehnt")

	classified := sqlexec.Classify(class, cause)

	if !stderrors.Is(classified, class) {
		t.Fatalf("errors.Is(err, class) = false, Fehler: %v", classified)
	}
	if !stderrors.Is(classified, cause) {
		t.Fatalf("errors.Is(err, cause) = false, Fehler: %v", classified)
	}
}

func TestIsAbsentDistinguishesAbsenceFromErrors(t *testing.T) {
	if !sqlexec.IsAbsent(pgx.ErrNoRows) {
		t.Fatal("pgx.ErrNoRows muss als Abwesenheit gelten")
	}
	if !sqlexec.IsAbsent(fmt.Errorf("umhüllt: %w", pgx.ErrNoRows)) {
		t.Fatal("die Abwesenheit muss durch die Wrappung lesbar bleiben")
	}
	if sqlexec.IsAbsent(stderrors.New("Verbindung abgelehnt")) {
		t.Fatal("ein Treiber-Fehler ist keine Abwesenheit")
	}
	if sqlexec.IsAbsent(nil) {
		t.Fatal("kein Fehler ist keine Abwesenheit")
	}
}

// --- Zeilen-Übersetzung ---

// changeRow trägt eine Ergebnis-Zeile der Change-Abfrage in der
// Projektions-Ordnung von `queries.SelectChanges` — dieselbe Ordnung wie
// die View `cdc.changes` (`tools/schema/schema.yaml`).
func changeRow(id string, sequence int64, position int64) []any {
	return []any{
		"src-1",
		position,
		id,
		"tx-1",
		"tbl-1",
		"public",
		"feed",
		sequence,
		string(model.OperationInsert),
		[]byte(nil),
		[]byte(`{"name":"a"}`),
		"sv-1",
		time.Unix(100, 0),
		"wal",
	}
}

// changeRowWithOrigin trägt eine Ergebnis-Zeile der Change-Abfrage mit der
// übergebenen Herkunft als letzter Spalte (`SPEC-002`, `origin`).
func changeRowWithOrigin(origin string) []any {
	row := changeRow("chg-1", 1, 42)
	row[len(row)-1] = origin
	return row
}

// failRecorder trägt die Ursachen, die der Aufrufer klassifizieren würde —
// die Stelle, an der die Adapter ihre Klasse und ihren Log führen.
type failRecorder struct {
	class  error
	causes []error
}

func (f *failRecorder) fail(cause error) error {
	f.causes = append(f.causes, cause)
	return sqlexec.Classify(f.class, cause)
}

func TestReadChangesIssuesQueryAndTranslatesRows(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		changeRow("chg-1", 1, 42),
		changeRow("chg-2", 2, 43),
	}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	records, err := sqlexec.ReadChanges(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT changes",
		Args: []any{"src-1", int64(1), nil, nil, 100},
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}

	if len(exec.queries) != 1 {
		t.Fatalf("erwartete 1 Query-Aufruf, gesehen %d", len(exec.queries))
	}
	if exec.queries[0].sql != "SELECT changes" {
		t.Fatalf("SQL = %q", exec.queries[0].sql)
	}
	if !reflect.DeepEqual(exec.queries[0].args, []any{"src-1", int64(1), nil, nil, 100}) {
		t.Fatalf("Argumente = %v", exec.queries[0].args)
	}
	if !exec.rows.closed {
		t.Fatal("die Ergebnis-Menge muss geschlossen werden")
	}
	if len(records) != 2 {
		t.Fatalf("erwartete 2 Records, gesehen %d", len(records))
	}
	if records[0].Change.ID != "chg-1" || records[1].Change.ID != "chg-2" {
		t.Fatalf("Changes = %v", records)
	}
	if records[0].Position.Offset != 42 || records[0].Position.SourceID != "src-1" {
		t.Fatalf("Position = %+v", records[0].Position)
	}
	// Die Zeile trägt die Klartext-Identität der Tabelle (Join auf
	// `cdc.source_table`); der Lesepfad setzt sie am Change (`ADR-0081`
	// Teilfrage 3).
	if records[0].Change.Schema != "public" || records[0].Change.Table != "feed" {
		t.Fatalf("Schema/Table = %q/%q, wollen public/feed", records[0].Change.Schema, records[0].Change.Table)
	}
	if want := int64(time.Unix(100, 0).UnixNano()); records[0].CommittedAt.UnixNanos != want {
		t.Fatalf("CommittedAt = %d, erwartet %d", records[0].CommittedAt.UnixNanos, want)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("der Erfolgsfall darf die Klasse nicht berühren: %v", recorder.causes)
	}
}

// Die letzte Spalte der Projektion trägt die Herkunft (`SPEC-002`): ein
// gelesener Wert erreicht den Change, ein Wert außerhalb der geschlossenen
// Menge endet als Domänen-Fehler statt als gefälschter Change (`ADR-0029`).
func TestReadChangesCarriesOrigin(t *testing.T) {
	for _, tc := range []struct {
		origin  string
		want    model.ChangeOrigin
		wantErr error
	}{
		{"wal", model.ChangeOriginWAL, nil},
		{"backfill", model.ChangeOriginBackfill, nil},
		{"", model.ChangeOriginWAL, nil},
		{"snapshot", "", domainerrors.ErrInvalidChangeOrigin},
	} {
		t.Run(tc.origin, func(t *testing.T) {
			exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{changeRowWithOrigin(tc.origin)}}}
			records, err := sqlexec.ReadChanges(context.Background(), exec, sqlexec.Statement{
				SQL:  "SELECT changes",
				Fail: (&failRecorder{class: outbound.ErrStorage}).fail,
			})
			if !stderrors.Is(err, tc.wantErr) {
				t.Fatalf("Fehler = %v, wollen %v", err, tc.wantErr)
			}
			if tc.wantErr == nil && records[0].Change.Origin != tc.want {
				t.Fatalf("Origin = %q, wollen %q", records[0].Change.Origin, tc.want)
			}
		})
	}
}

func TestReadChangesEmptyResultIsNoError(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	records, err := sqlexec.ReadChanges(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT changes",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if records == nil || len(records) != 0 {
		t.Fatalf("erwartete leere, gesetzte Rückgabe, gesehen %v", records)
	}
}

func TestReadChangesClassifiesQueryFailure(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	exec := &fakeExecutor{queryErr: cause}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadChanges(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT changes",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrStorage) {
		t.Fatalf("Fehler trägt nicht die Klasse storage: %v", err)
	}
	if !stderrors.Is(err, cause) {
		t.Fatalf("die technische Ursache muss lesbar bleiben: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadChangesClassifiesScanFailure(t *testing.T) {
	cause := stderrors.New("Spaltentyp passt nicht")
	exec := &fakeExecutor{rows: &fakeRows{
		rows:     [][]any{changeRow("chg-1", 1, 42)},
		scanErrs: map[int]error{0: cause},
	}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadChanges(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT changes",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrStorage) {
		t.Fatalf("Scan-Fehler trägt nicht die Klasse storage: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadChangesClassifiesIterationFailure(t *testing.T) {
	cause := stderrors.New("Ergebnis-Menge abgebrochen")
	exec := &fakeExecutor{rows: &fakeRows{iterErr: cause}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadChanges(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT changes",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrStorage) {
		t.Fatalf("Iterations-Fehler trägt nicht die Klasse storage: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadChangesLeavesDomainFailureUnclassified(t *testing.T) {
	// Die Zeile trägt die Commit-Position 0 — der Mapper lehnt sie ab; das
	// ist ein Domänen-Fehler der Übersetzung, keine Treiber-Klasse.
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{changeRow("chg-1", 1, 0)}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadChanges(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT changes",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, domainerrors.ErrInvalidPosition) {
		t.Fatalf("erwartete die Domänen-Invariante, gesehen: %v", err)
	}
	if stderrors.Is(err, outbound.ErrStorage) {
		t.Fatalf("ein Domänen-Fehler darf die Treiber-Klasse nicht tragen: %v", err)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("der Übersetzungspunkt darf hier nicht laufen: %v", recorder.causes)
	}
}

func TestStatementWithoutClassReturnsCause(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	exec := &fakeExecutor{queryErr: cause}

	_, err := sqlexec.ReadPendingRequests(context.Background(), exec, sqlexec.Statement{
		SQL: "SELECT requests",
	})

	if err != cause {
		t.Fatalf("ohne Übersetzungspunkt muss die Ursache unverändert zurückkommen: %v", err)
	}
}

func TestReadConsumerPositionReadsAcknowledgedPosition(t *testing.T) {
	exec := &fakeExecutor{row: &fakeRow{values: []any{"src-1", int64(77)}}}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	position, err := sqlexec.ReadConsumerPosition(context.Background(), exec, model.ConsumerID("consumer-1"), sqlexec.Statement{
		SQL:  "SELECT position",
		Args: []any{"consumer-1"},
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadConsumerPosition: %v", err)
	}
	if position.ConsumerID != "consumer-1" {
		t.Fatalf("ConsumerID = %q", position.ConsumerID)
	}
	if position.Position.Offset != 77 || position.Position.SourceID != "src-1" {
		t.Fatalf("Position = %+v", position.Position)
	}
	if len(exec.queryRows) != 1 || exec.queryRows[0].sql != "SELECT position" {
		t.Fatalf("QueryRow-Aufrufe = %v", exec.queryRows)
	}
}

func TestReadConsumerPositionReadsAbsenceAsZero(t *testing.T) {
	exec := &fakeExecutor{row: &fakeRow{err: pgx.ErrNoRows}}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	position, err := sqlexec.ReadConsumerPosition(context.Background(), exec, model.ConsumerID("consumer-1"), sqlexec.Statement{
		SQL:  "SELECT position",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("die Abwesenheit ist kein Fehler: %v", err)
	}
	if position.ConsumerID != "" || !position.Position.IsZero() {
		t.Fatalf("erwartete den Nullwert, gesehen %+v", position)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("die Abwesenheit darf die Klasse nicht berühren: %v", recorder.causes)
	}
}

func TestReadConsumerPositionClassifiesReadFailure(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	exec := &fakeExecutor{row: &fakeRow{err: cause}}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	_, err := sqlexec.ReadConsumerPosition(context.Background(), exec, model.ConsumerID("consumer-1"), sqlexec.Statement{
		SQL:  "SELECT position",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrConsumerStateStorage) {
		t.Fatalf("Fehler trägt nicht die Klasse des Ports: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadConsumerPositionsTranslatesRows(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"consumer-1", int64(10)},
		{"consumer-2", int64(20)},
	}}}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	positions, err := sqlexec.ReadConsumerPositions(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{
		SQL:  "SELECT positions",
		Args: []any{"src-1"},
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadConsumerPositions: %v", err)
	}
	if len(positions) != 2 {
		t.Fatalf("erwartete 2 Positionen, gesehen %d", len(positions))
	}
	if positions[1].ConsumerID != "consumer-2" || positions[1].Position.Offset != 20 {
		t.Fatalf("Position = %+v", positions[1])
	}
	if positions[1].Position.SourceID != "src-1" {
		t.Fatalf("die Quelle kommt aus dem Aufruf: %+v", positions[1].Position)
	}
}

func TestReadConsumerPositionsEmptyResult(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{}}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	positions, err := sqlexec.ReadConsumerPositions(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{
		SQL:  "SELECT positions",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadConsumerPositions: %v", err)
	}
	if positions == nil || len(positions) != 0 {
		t.Fatalf("erwartete leere, gesetzte Rückgabe, gesehen %v", positions)
	}
}

func TestReadSourceTablesTranslatesRows(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"tbl-1", "src-1", "public", "feed"},
	}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	tables, err := sqlexec.ReadSourceTables(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT tables",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadSourceTables: %v", err)
	}
	if len(tables) != 1 {
		t.Fatalf("erwartete 1 Tabelle, gesehen %d", len(tables))
	}
	if tables[0].QualifiedName() != "public.feed" || tables[0].SourceID != "src-1" {
		t.Fatalf("Tabelle = %+v", tables[0])
	}
}

func TestReadSourceTablesLeavesDomainFailureUnclassified(t *testing.T) {
	// Die Zeile trägt keine Tabellen-Kennung — der Domänen-Konstruktor lehnt
	// sie ab.
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"", "src-1", "public", "feed"},
	}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadSourceTables(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT tables",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("erwartete die Domänen-Invariante, gesehen: %v", err)
	}
	if stderrors.Is(err, outbound.ErrStorage) {
		t.Fatalf("ein Domänen-Fehler darf die Treiber-Klasse nicht tragen: %v", err)
	}
}

func TestReadExcludedColumnsCarriesExclusionAndInclusion(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"public", "feed", string(model.AdministrationRequestExcludeColumn), "secret"},
		{"public", "feed", string(model.AdministrationRequestExcludeColumn), "secret"},
		{"public", "feed", string(model.AdministrationRequestExcludeColumn), "token"},
		{"public", "feed", string(model.AdministrationRequestEnable), "ignored"},
	}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	excluded, err := sqlexec.ReadExcludedColumns(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT applied column requests",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadExcludedColumns: %v", err)
	}
	if got := excluded["public.feed"]; !reflect.DeepEqual(got, []string{"secret", "token"}) {
		t.Fatalf("Ausschlussstand = %v", got)
	}
}

func TestReadExcludedColumnsRemovesIncludedColumn(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"public", "feed", string(model.AdministrationRequestExcludeColumn), "secret"},
		{"public", "feed", string(model.AdministrationRequestExcludeColumn), "token"},
		{"public", "feed", string(model.AdministrationRequestIncludeColumn), "secret"},
	}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	excluded, err := sqlexec.ReadExcludedColumns(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT applied column requests",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadExcludedColumns: %v", err)
	}
	if got := excluded["public.feed"]; !reflect.DeepEqual(got, []string{"token"}) {
		t.Fatalf("Ausschlussstand = %v", got)
	}
}

func TestReadExcludedColumnsLastInclusionRemovesEntry(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"public", "feed", string(model.AdministrationRequestExcludeColumn), "secret"},
		{"public", "feed", string(model.AdministrationRequestIncludeColumn), "secret"},
	}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	excluded, err := sqlexec.ReadExcludedColumns(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT applied column requests",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadExcludedColumns: %v", err)
	}
	if len(excluded) != 0 {
		t.Fatalf("ein Stand ohne Ausschluss ist ein fehlender Eintrag: %v", excluded)
	}
}

func TestReadTableSchemaTranslatesColumns(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"id", int64(23)},
		{"name", int64(25)},
	}}}
	recorder := &failRecorder{class: outbound.ErrSchemaStoreStorage}

	schema, err := sqlexec.ReadTableSchema(context.Background(), exec, model.SchemaVersionID("sv-1"), sqlexec.Statement{
		SQL:  "SELECT columns",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadTableSchema: %v", err)
	}
	if schema.VersionID != "sv-1" {
		t.Fatalf("VersionID = %q", schema.VersionID)
	}
	if len(schema.Columns) != 2 || schema.Columns[0].Name != "id" || schema.Columns[1].OID != 25 {
		t.Fatalf("Spalten = %+v", schema.Columns)
	}
}

func TestReadTableSchemaReportsUnknownVersion(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{}}
	recorder := &failRecorder{class: outbound.ErrSchemaStoreStorage}

	_, err := sqlexec.ReadTableSchema(context.Background(), exec, model.SchemaVersionID("sv-1"), sqlexec.Statement{
		SQL:  "SELECT columns",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrSchemaVersionUnknown) {
		t.Fatalf("erwartete den Sentinel der unbekannten Version, gesehen: %v", err)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("die leere Spaltenform ist kein Treiber-Fehler: %v", recorder.causes)
	}
}

func TestReadPendingRequestsTranslatesRequests(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"req-1", "src-1", "public", "feed", "", "", "", string(model.AdministrationRequestEnable)},
		{"req-2", "src-1", "public", "feed", "secret", "", "", string(model.AdministrationRequestExcludeColumn)},
		{"req-3", "src-1", "public", "feed", "", "umbenennung", `{"kind": "rename_column"}`, string(model.AdministrationRequestSetTransformation)},
		{"req-4", "src-1", "public", "feed", "", "umbenennung", "", string(model.AdministrationRequestRemoveTransformation)},
	}}}
	recorder := &failRecorder{class: outbound.ErrAdministrationStorage}

	requests, err := sqlexec.ReadPendingRequests(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT pending",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadPendingRequests: %v", err)
	}
	if len(requests) != 4 {
		t.Fatalf("erwartete 4 Anträge, gesehen %d", len(requests))
	}
	for i, row := range requests {
		if row.Rejected != nil {
			t.Fatalf("Zeile %d ist verworfen: %+v", i, row.Rejected)
		}
	}
	if requests[1].Request.Kind != model.AdministrationRequestExcludeColumn || requests[1].Request.Column != "secret" {
		t.Fatalf("Antrag = %+v", requests[1])
	}
	if requests[2].Request.Kind != model.AdministrationRequestSetTransformation || requests[2].Request.RuleName != "umbenennung" || requests[2].Request.RuleSpec != `{"kind": "rename_column"}` {
		t.Fatalf("Antrag = %+v, wollen set_transformation mit Regelname und Regelform", requests[2])
	}
	if requests[3].Request.Kind != model.AdministrationRequestRemoveTransformation || requests[3].Request.RuleName != "umbenennung" || requests[3].Request.RuleSpec != "" {
		t.Fatalf("Antrag = %+v, wollen remove_transformation mit Regelname ohne Regelform", requests[3])
	}
}

func TestRegisterConsumerReportsOutcome(t *testing.T) {
	exec := &fakeExecutor{tag: pgconn.NewCommandTag("INSERT 0 1")}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	registered, err := sqlexec.RegisterConsumer(context.Background(), exec, sqlexec.Statement{
		SQL:  "INSERT consumer",
		Args: []any{"consumer-1", "Name"},
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("RegisterConsumer: %v", err)
	}
	if !registered {
		t.Fatal("eine betroffene Zeile heißt registriert")
	}
	if len(exec.executions) != 1 || exec.executions[0].sql != "INSERT consumer" {
		t.Fatalf("Exec-Aufrufe = %v", exec.executions)
	}
	if !reflect.DeepEqual(exec.executions[0].args, []any{"consumer-1", "Name"}) {
		t.Fatalf("Argumente = %v", exec.executions[0].args)
	}
}

func TestRegisterConsumerReportsRepeat(t *testing.T) {
	exec := &fakeExecutor{tag: pgconn.NewCommandTag("INSERT 0 0")}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	registered, err := sqlexec.RegisterConsumer(context.Background(), exec, sqlexec.Statement{
		SQL:  "INSERT consumer",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("RegisterConsumer: %v", err)
	}
	if registered {
		t.Fatal("keine betroffene Zeile heißt: war schon da")
	}
}

func TestRegisterConsumerClassifiesExecFailure(t *testing.T) {
	cause := stderrors.New("Fremdschlüssel verletzt")
	exec := &fakeExecutor{execErr: cause}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	_, err := sqlexec.RegisterConsumer(context.Background(), exec, sqlexec.Statement{
		SQL:  "INSERT consumer",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrConsumerStateStorage) {
		t.Fatalf("Fehler trägt nicht die Klasse des Ports: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

// Der Change der Zeile endet außerhalb der Change-Invarianten (Sequenz 0):
// der Domänen-Konstruktor lehnt ihn ab, der Übersetzungspunkt des Aufrufers
// bleibt unberührt (`ADR-0029`) — der Fehler trägt seine Klasse schon selbst.
func TestReadChangesLeavesChangeTranslationFailureUnclassified(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{changeRow("chg-1", 0, 42)}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadChanges(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT changes",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, domainerrors.ErrNonPositiveSequence) {
		t.Fatalf("erwartete die Domänen-Invariante, gesehen: %v", err)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("der Übersetzungspunkt darf hier nicht laufen: %v", recorder.causes)
	}
}

// Eine bestätigte Position der Spalte 0 ist keine Position (`SPEC-003`): der
// Mapper lehnt sie ab, der Übersetzungspunkt bleibt unberührt — derselbe
// Domänen-Fehler wie im Lesepfad, keine Treiber-Klasse.
func TestReadConsumerPositionRejectsNonPositiveColumn(t *testing.T) {
	exec := &fakeExecutor{row: &fakeRow{values: []any{"src-1", int64(0)}}}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	_, err := sqlexec.ReadConsumerPosition(context.Background(), exec, model.ConsumerID("consumer-1"), sqlexec.Statement{
		SQL:  "SELECT position",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, domainerrors.ErrInvalidPosition) {
		t.Fatalf("erwartete die Domänen-Invariante, gesehen: %v", err)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("der Übersetzungspunkt darf hier nicht laufen: %v", recorder.causes)
	}
}

func TestReadConsumerPositionsClassifiesQueryFailure(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	exec := &fakeExecutor{queryErr: cause}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	_, err := sqlexec.ReadConsumerPositions(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{
		SQL:  "SELECT positions",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrConsumerStateStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadConsumerPositionsClassifiesScanFailure(t *testing.T) {
	cause := stderrors.New("Spaltentyp passt nicht")
	exec := &fakeExecutor{rows: &fakeRows{
		rows:     [][]any{{"consumer-1", int64(10)}},
		scanErrs: map[int]error{0: cause},
	}}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	_, err := sqlexec.ReadConsumerPositions(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{
		SQL:  "SELECT positions",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrConsumerStateStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadConsumerPositionsLeavesDomainFailureUnclassified(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"consumer-1", int64(10)},
		{"consumer-2", int64(0)},
	}}}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	_, err := sqlexec.ReadConsumerPositions(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{
		SQL:  "SELECT positions",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, domainerrors.ErrInvalidPosition) {
		t.Fatalf("erwartete die Domänen-Invariante, gesehen: %v", err)
	}
	if stderrors.Is(err, outbound.ErrConsumerStateStorage) {
		t.Fatalf("ein Domänen-Fehler darf die Treiber-Klasse nicht tragen: %v", err)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("der Übersetzungspunkt darf hier nicht laufen: %v", recorder.causes)
	}
}

func TestReadConsumerPositionsClassifiesIterationFailure(t *testing.T) {
	cause := stderrors.New("Ergebnis-Menge abgebrochen")
	exec := &fakeExecutor{rows: &fakeRows{iterErr: cause}}
	recorder := &failRecorder{class: outbound.ErrConsumerStateStorage}

	_, err := sqlexec.ReadConsumerPositions(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{
		SQL:  "SELECT positions",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrConsumerStateStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadSourceTablesClassifiesQueryFailure(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	exec := &fakeExecutor{queryErr: cause}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadSourceTables(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT tables",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadSourceTablesClassifiesScanFailure(t *testing.T) {
	cause := stderrors.New("Spaltentyp passt nicht")
	exec := &fakeExecutor{rows: &fakeRows{
		rows:     [][]any{{"tbl-1", "src-1", "public", "feed"}},
		scanErrs: map[int]error{0: cause},
	}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadSourceTables(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT tables",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadSourceTablesClassifiesIterationFailure(t *testing.T) {
	cause := stderrors.New("Ergebnis-Menge abgebrochen")
	exec := &fakeExecutor{rows: &fakeRows{iterErr: cause}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadSourceTables(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT tables",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadExcludedColumnsClassifiesQueryFailure(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	exec := &fakeExecutor{queryErr: cause}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadExcludedColumns(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT applied column requests",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadExcludedColumnsClassifiesScanFailure(t *testing.T) {
	cause := stderrors.New("Spaltentyp passt nicht")
	exec := &fakeExecutor{rows: &fakeRows{
		rows: [][]any{
			{"public", "feed", string(model.AdministrationRequestExcludeColumn), "secret"},
		},
		scanErrs: map[int]error{0: cause},
	}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadExcludedColumns(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT applied column requests",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadExcludedColumnsClassifiesIterationFailure(t *testing.T) {
	cause := stderrors.New("Ergebnis-Menge abgebrochen")
	exec := &fakeExecutor{rows: &fakeRows{iterErr: cause}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadExcludedColumns(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT applied column requests",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadTableSchemaClassifiesQueryFailure(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	exec := &fakeExecutor{queryErr: cause}
	recorder := &failRecorder{class: outbound.ErrSchemaStoreStorage}

	_, err := sqlexec.ReadTableSchema(context.Background(), exec, model.SchemaVersionID("sv-1"), sqlexec.Statement{
		SQL:  "SELECT columns",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrSchemaStoreStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadTableSchemaClassifiesScanFailure(t *testing.T) {
	cause := stderrors.New("Spaltentyp passt nicht")
	exec := &fakeExecutor{rows: &fakeRows{
		rows:     [][]any{{"id", int64(23)}},
		scanErrs: map[int]error{0: cause},
	}}
	recorder := &failRecorder{class: outbound.ErrSchemaStoreStorage}

	_, err := sqlexec.ReadTableSchema(context.Background(), exec, model.SchemaVersionID("sv-1"), sqlexec.Statement{
		SQL:  "SELECT columns",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrSchemaStoreStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadTableSchemaClassifiesIterationFailure(t *testing.T) {
	cause := stderrors.New("Ergebnis-Menge abgebrochen")
	exec := &fakeExecutor{rows: &fakeRows{iterErr: cause}}
	recorder := &failRecorder{class: outbound.ErrSchemaStoreStorage}

	_, err := sqlexec.ReadTableSchema(context.Background(), exec, model.SchemaVersionID("sv-1"), sqlexec.Statement{
		SQL:  "SELECT columns",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrSchemaStoreStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadPendingRequestsClassifiesScanFailure(t *testing.T) {
	cause := stderrors.New("Spaltentyp passt nicht")
	exec := &fakeExecutor{rows: &fakeRows{
		rows: [][]any{
			{"req-1", "src-1", "public", "feed", "", "", "", string(model.AdministrationRequestEnable)},
		},
		scanErrs: map[int]error{0: cause},
	}}
	recorder := &failRecorder{class: outbound.ErrAdministrationStorage}

	_, err := sqlexec.ReadPendingRequests(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT pending",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrAdministrationStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

func TestReadPendingRequestsClassifiesIterationFailure(t *testing.T) {
	cause := stderrors.New("Ergebnis-Menge abgebrochen")
	exec := &fakeExecutor{rows: &fakeRows{iterErr: cause}}
	recorder := &failRecorder{class: outbound.ErrAdministrationStorage}

	_, err := sqlexec.ReadPendingRequests(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT pending",
		Fail: recorder.fail,
	})

	if !stderrors.Is(err, outbound.ErrAdministrationStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	if len(recorder.causes) != 1 || recorder.causes[0] != cause {
		t.Fatalf("Übersetzungspunkt gesehen: %v", recorder.causes)
	}
}

// Die Naht selbst: der reale Pool erfüllt sie, ein Träger der Fakes ebenso —
// die Zusicherung steht in `seam.go` (Kompilier-Beleg), dieser Test hält die
// Fake-Seite dagegen.
var _ sqlexec.Executor = (*fakeExecutor)(nil)

var _ pgx.Rows = (*fakeRows)(nil)

// --- Backfill-Runs ---

// backfillRunRow trägt eine Ergebnis-Zeile der Run-Abfrage in der Spalten-
// Reihenfolge von `ReadBackfillRuns`.
func backfillRunRow(runID, status string, started *time.Time, position, estimate *int64) []any {
	return []any{
		runID, "src-1", "public", "orders", status, time.Unix(100, 0).UTC(), started, (*time.Time)(nil),
		position, int64(3), estimate, "", false, true,
	}
}

func TestReadBackfillRunsTranslatesRuns(t *testing.T) {
	started := time.Unix(200, 0).UTC()
	position := int64(7000)
	estimate := int64(50)
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		backfillRunRow("run-1", "running", &started, &position, &estimate),
		backfillRunRow("run-2", "queued", nil, nil, nil),
	}}}
	recorder := &failRecorder{class: outbound.ErrBackfillStorage}

	runs, err := sqlexec.ReadBackfillRuns(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT runs",
		Args: []any{"src-1"},
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadBackfillRuns: %v", err)
	}
	if len(runs) != 2 || runs[0].ID != "run-1" || runs[1].ID != "run-2" {
		t.Fatalf("Runs = %+v", runs)
	}
	if runs[0].Status != model.BackfillRunRunning || runs[0].StartedAt.UnixNanos != started.UnixNano() ||
		runs[0].SnapshotPosition.Offset != 7000 || runs[0].RowsCopied != 3 || !runs[0].WarnDuration {
		t.Fatalf("Run 1 = %+v", runs[0])
	}
	if rows, known := runs[0].EstimatedRows.Rows(); !known || rows != 50 {
		t.Fatalf("Schätzung Run 1 = %d, bekannt %v", rows, known)
	}
	if _, known := runs[1].EstimatedRows.Rows(); known {
		t.Fatal("NULL-Schätzung liest als bekannt")
	}
	if len(exec.queries) != 1 || exec.queries[0].sql != "SELECT runs" || !reflect.DeepEqual(exec.queries[0].args, []any{"src-1"}) {
		t.Fatalf("Query-Aufrufe = %v", exec.queries)
	}
}

func TestReadBackfillRunsLeavesDomainFailureUnclassified(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		backfillRunRow("run-1", "paused", nil, nil, nil),
	}}}
	recorder := &failRecorder{class: outbound.ErrBackfillStorage}

	_, err := sqlexec.ReadBackfillRuns(context.Background(), exec, sqlexec.Statement{SQL: "SELECT runs", Fail: recorder.fail})

	if !stderrors.Is(err, mapper.ErrUnknownBackfillStatus) {
		t.Fatalf("erwartete die Status-Invariante, gesehen: %v", err)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("ein Domänen-Fehler läuft nicht durch den Übersetzungspunkt: %v", recorder.causes)
	}
}

func TestReadBackfillRunsClassifiesDriverFailures(t *testing.T) {
	cause := stderrors.New("Treiber-Fehler")
	cases := map[string]*fakeExecutor{
		"Abfrage":   {queryErr: cause},
		"Scan":      {rows: &fakeRows{rows: [][]any{backfillRunRow("run-1", "queued", nil, nil, nil)}, scanErrs: map[int]error{0: cause}}},
		"Iteration": {rows: &fakeRows{iterErr: cause}},
	}
	for name, exec := range cases {
		recorder := &failRecorder{class: outbound.ErrBackfillStorage}
		_, err := sqlexec.ReadBackfillRuns(context.Background(), exec, sqlexec.Statement{SQL: "SELECT runs", Fail: recorder.fail})
		if !stderrors.Is(err, outbound.ErrBackfillStorage) || !stderrors.Is(err, cause) {
			t.Fatalf("%s: Fehler trägt nicht Klasse und Ursache: %v", name, err)
		}
		if len(recorder.causes) != 1 || recorder.causes[0] != cause {
			t.Fatalf("%s: Übersetzungspunkt gesehen: %v", name, recorder.causes)
		}
	}
}

// --- Retention-Kandidaten ---

func TestReadRetentionCandidatesIssuesQueryAndTranslatesRows(t *testing.T) {
	committed := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"0bf-1", int64(1000), committed},
		{"tx-9-1", int64(1000), committed.Add(time.Second)},
	}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	candidates, err := sqlexec.ReadRetentionCandidates(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{
		SQL:  "SELECT candidates",
		Args: []any{"src-1", "", 10},
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadRetentionCandidates: %v", err)
	}
	if len(exec.queries) != 1 || exec.queries[0].sql != "SELECT candidates" || !reflect.DeepEqual(exec.queries[0].args, []any{"src-1", "", 10}) {
		t.Fatalf("abgesetzter Aufruf = %+v", exec.queries)
	}
	want := []outbound.RetentionCandidate{
		{ChangeID: "0bf-1", Position: model.SourcePosition{SourceID: "src-1", Offset: 1000}, CommittedAt: model.NewTimePoint(committed.UnixNano())},
		{ChangeID: "tx-9-1", Position: model.SourcePosition{SourceID: "src-1", Offset: 1000}, CommittedAt: model.NewTimePoint(committed.Add(time.Second).UnixNano())},
	}
	if !reflect.DeepEqual(candidates, want) {
		t.Fatalf("Kandidaten = %+v, wollen %+v", candidates, want)
	}
}

func TestReadRetentionCandidatesEmptyResult(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	candidates, err := sqlexec.ReadRetentionCandidates(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{SQL: "SELECT candidates", Fail: recorder.fail})
	if err != nil {
		t.Fatalf("ReadRetentionCandidates: %v", err)
	}
	if candidates == nil || len(candidates) != 0 {
		t.Fatalf("erwartete leere, gesetzte Rückgabe, gesehen %v", candidates)
	}
}

func TestReadRetentionCandidatesLeavesDomainFailureUnclassified(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{{"tx-1-1", int64(0), time.Unix(0, 0)}}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	_, err := sqlexec.ReadRetentionCandidates(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{SQL: "SELECT candidates", Fail: recorder.fail})

	if !stderrors.Is(err, domainerrors.ErrInvalidPosition) {
		t.Fatalf("erwartete die Positions-Invariante, gesehen: %v", err)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("ein Domänen-Fehler läuft nicht durch den Übersetzungspunkt: %v", recorder.causes)
	}
}

func TestReadRetentionCandidatesClassifiesDriverFailures(t *testing.T) {
	cause := stderrors.New("Treiber-Fehler")
	cases := map[string]*fakeExecutor{
		"Abfrage":   {queryErr: cause},
		"Scan":      {rows: &fakeRows{rows: [][]any{{"tx-1-1", int64(1), time.Unix(0, 0)}}, scanErrs: map[int]error{0: cause}}},
		"Iteration": {rows: &fakeRows{iterErr: cause}},
	}
	for name, exec := range cases {
		recorder := &failRecorder{class: outbound.ErrStorage}
		_, err := sqlexec.ReadRetentionCandidates(context.Background(), exec, model.SourceID("src-1"), sqlexec.Statement{SQL: "SELECT candidates", Fail: recorder.fail})
		if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, cause) {
			t.Fatalf("%s: Fehler trägt nicht Klasse und Ursache: %v", name, err)
		}
		if len(recorder.causes) != 1 || recorder.causes[0] != cause {
			t.Fatalf("%s: Übersetzungspunkt gesehen: %v", name, recorder.causes)
		}
	}
}

// ReadPendingRequests liefert eine Zeile der Transformations-Antragsarten mit
// leerem Regelnamen und leerer Regelform (SQL-NULL, vom `COALESCE` der Abfrage
// normalisiert) als Antrag statt als verworfene Zeile — der Antrag endet im
// Use Case als `failed` mit dem Fehlertext der Spec (`SPEC-019`), und die
// gültige Zeile dahinter bleibt lesbar. Rot färbende Mutation: die Prüfung
// `ruleName == ""` in den Konstruktor `model.NewAdministrationRequest`
// zurücklegen — die Zeile steht als `Rejected` statt als Antrag im Ergebnis.
func TestReadPendingRequestsCarriesRuleRowsWithEmptyRuleFields(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"req-1", "src-1", "public", "feed", "", "", "", string(model.AdministrationRequestSetTransformation)},
		{"req-2", "src-1", "public", "feed", "", "", "null", string(model.AdministrationRequestSetTransformation)},
		{"req-3", "src-1", "public", "feed", "", "", "", string(model.AdministrationRequestRemoveTransformation)},
		{"req-4", "src-1", "public", "feed", "", "gueltig", `{"kind": "rename_column", "column": "name", "to": "x"}`, string(model.AdministrationRequestSetTransformation)},
	}}}
	recorder := &failRecorder{class: outbound.ErrAdministrationStorage}

	requests, err := sqlexec.ReadPendingRequests(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT pending",
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadPendingRequests = %v, wollen nil (die Zeilen sind Anträge, keine Lesefehler)", err)
	}
	if len(requests) != 4 {
		t.Fatalf("erwartete 4 Anträge, gesehen %d", len(requests))
	}
	for i, row := range requests {
		if row.Rejected != nil {
			t.Fatalf("Zeile %d ist verworfen: %+v", i, row.Rejected)
		}
	}
	if requests[0].Request.RuleName != "" || requests[0].Request.RuleSpec != "" || requests[1].Request.RuleSpec != "null" || requests[2].Request.RuleName != "" {
		t.Fatalf("Anträge = %+v, wollen leeren Regelnamen und leere Regelform bzw. den Text null unverändert", requests[:3])
	}
	if requests[3].Request.RuleName != "gueltig" {
		t.Fatalf("gültiger Antrag hinter den ungültigen = %+v", requests[3])
	}
}

// ReadPendingRequests reicht eine vom Antrags-Konstruktor verworfene Zeile mit
// Kennung und Fehlertext durch (`SPEC-019`), statt die Lesung zu beenden: je
// Grund — leere Quelle, leeres Schema, leerer Tabellenname,
// `exclude_column`/`include_column` mit leerer Spalte, Antragsart außerhalb der
// geschlossenen Menge, ein Grund ohne eigenen Klartext — steht die Zeile als
// `Rejected` an ihrer Stelle der Ordnung, mit einer gültigen Zeile davor und
// einer dahinter; die Lesung endet ohne Fehler und der Übersetzungspunkt läuft
// nicht. Der Text je Grund ist der Klartext, Doppelpunkt, Leerzeichen und die
// Antrags-Kennung; mehrere Gründe an einer Zeile nennen den ersten der
// Konstruktor-Reihenfolge. Rot färbende Mutationen (Eingabeseite, die Zeilen
// der Fälle): den Konstruktor-Fehler wieder zurückgeben statt die Zeile
// durchzureichen — jeder Fall endet mit einem Fehler; die verworfenen Zeilen
// ans Ende des Ergebnisses stellen — die Zeile zwischen den gültigen ist keine
// verworfene; in `rejectionMessage` die Fälle `schema` und `table` vertauschen
// — der Fall mit beiden Gründen nennt den Tabellennamen; den Text durch einen
// festen ersetzen — der Text je Grund weicht ab.
func TestReadPendingRequestsPassesRejectedRowsThrough(t *testing.T) {
	const (
		remove  = string(model.AdministrationRequestRemoveTransformation)
		set     = string(model.AdministrationRequestSetTransformation)
		exclude = string(model.AdministrationRequestExcludeColumn)
		include = string(model.AdministrationRequestIncludeColumn)
		ruleRow = `{"kind": "rename_column", "column": "a", "to": "b"}`
	)
	for _, tc := range []struct {
		name    string
		row     []any
		wantID  model.AdministrationRequestID
		wantMsg string
	}{
		{"leere Quelle", []any{"req-x", "", "public", "feed", "", "regel", ruleRow, set}, "req-x", "Quelle ist leer: req-x"},
		{"leeres Schema", []any{"req-x", "src-1", "", "feed", "", "regel", ruleRow, set}, "req-x", "Schemaname ist leer: req-x"},
		{"leerer Tabellenname", []any{"req-x", "src-1", "public", "", "", "regel", ruleRow, set}, "req-x", "Tabellenname ist leer: req-x"},
		{"leeres Schema und leerer Tabellenname", []any{"req-x", "src-1", "", "", "", "", "", remove}, "req-x", "Schemaname ist leer: req-x"},
		{"exclude_column mit leerer Spalte", []any{"req-x", "src-1", "public", "feed", "", "", "", exclude}, "req-x", "Spaltenname ist leer: req-x"},
		{"include_column mit leerer Spalte", []any{"req-x", "src-1", "public", "feed", "", "", "", include}, "req-x", "Spaltenname ist leer: req-x"},
		{"Antragsart außerhalb der Menge", []any{"req-x", "src-1", "public", "feed", "", "", "", "unbekannt"}, "req-x", "Antragsart ist unbekannt: req-x"},
		{"Zeile ohne Kennung", []any{"", "src-1", "public", "feed", "", "", "", remove}, "", "Kennung ist leer: public.feed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
				{"req-1", "src-1", "public", "feed", "", "regel_a", ruleRow, set},
				tc.row,
				{"req-2", "src-1", "public", "feed", "", "regel_a", "", remove},
			}}}
			recorder := &failRecorder{class: outbound.ErrAdministrationStorage}

			requests, err := sqlexec.ReadPendingRequests(context.Background(), exec, sqlexec.Statement{
				SQL:  "SELECT pending",
				Fail: recorder.fail,
			})
			if err != nil {
				t.Fatalf("ReadPendingRequests = %v, wollen nil (die Lesung lehnt keine Zeile ab)", err)
			}
			if len(recorder.causes) != 0 {
				t.Fatalf("die verworfene Zeile ist kein Lesefehler: %v", recorder.causes)
			}
			if len(requests) != 3 {
				t.Fatalf("erwartete 3 Zeilen in der Ordnung der Abfrage, gesehen %+v", requests)
			}
			if requests[0].Rejected != nil || requests[0].Request.ID != "req-1" || requests[2].Rejected != nil || requests[2].Request.ID != "req-2" {
				t.Fatalf("die gültigen Zeilen davor und dahinter sind Anträge: %+v", requests)
			}
			rejected := requests[1].Rejected
			if rejected == nil {
				t.Fatalf("die Zeile ist ein Antrag statt einer verworfenen Zeile: %+v", requests[1])
			}
			if rejected.ID != tc.wantID || rejected.Message != tc.wantMsg {
				t.Fatalf("verworfene Zeile = Kennung %q, Text %q, wollen %q und %q", rejected.ID, rejected.Message, tc.wantID, tc.wantMsg)
			}
		})
	}
}

// Ein Fehler des Lesens selbst bleibt ein Fehler: Anfrage, Scan und
// Iteration enden die Lesung mit der Klasse des Aufrufers, auch wenn Zeilen
// davor verworfen wurden (`TestReadPendingRequestsClassifiesScanFailure`,
// `TestReadPendingRequestsClassifiesIterationFailure` tragen die beiden
// anderen Stellen). Rot färbende Mutation: die Rückgabe von `rows.Err()`
// durch `nil` ersetzen — die Lesung endet ohne Fehler.
func TestReadPendingRequestsKeepsReadFailureAfterRejectedRow(t *testing.T) {
	cause := stderrors.New("Verbindung getrennt")
	exec := &fakeExecutor{rows: &fakeRows{
		rows:    [][]any{{"req-a", "src-1", "", "feed", "", "", "", string(model.AdministrationRequestEnable)}},
		iterErr: cause,
	}}
	recorder := &failRecorder{class: outbound.ErrAdministrationStorage}

	requests, err := sqlexec.ReadPendingRequests(context.Background(), exec, sqlexec.Statement{SQL: "SELECT pending", Fail: recorder.fail})

	if !stderrors.Is(err, outbound.ErrAdministrationStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Fehler = %v, wollen Klasse und Ursache der Iteration", err)
	}
	if len(requests) != 0 {
		t.Fatalf("Ergebnis = %+v, wollen keines neben dem Fehler", requests)
	}
}

// ReadTransformationRules faltet die Zeilen je Tabelle in der Ordnung der
// Abfrage: Set trägt ein, Remove nimmt heraus (Zyklus endet mit der letzten
// Regel), eine Tabelle ohne verbleibende Regel trägt keinen Eintrag, jede
// Tabelle hat ihre eigene Liste. Rot färbende Mutation: die Ordnung der
// Zeilen vor der Faltung umdrehen — die Zyklus-Tabelle endet ohne Regel.
func TestReadTransformationRulesFoldsPerTable(t *testing.T) {
	const (
		setKind    = string(model.AdministrationRequestSetTransformation)
		removeKind = string(model.AdministrationRequestRemoveTransformation)
	)
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"public", "orders", setKind, "kundenname", `{"kind": "rename_column", "column": "name", "to": "customer_name"}`},
		{"public", "orders", setKind, "statusname", `{"kind": "rename_column", "column": "status", "to": "state"}`},
		{"public", "cycle", setKind, "a", `{"kind": "rename_column", "column": "x", "to": "y"}`},
		{"public", "cycle", removeKind, "a", ""},
		{"public", "cycle", setKind, "a", `{"kind": "rename_column", "column": "x", "to": "z"}`},
		{"public", "gone", setKind, "a", `{"kind": "rename_column", "column": "x", "to": "y"}`},
		{"public", "gone", removeKind, "a", ""},
	}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	state, err := sqlexec.ReadTransformationRules(context.Background(), exec, sqlexec.Statement{
		SQL:  "SELECT applied transformation requests",
		Args: []any{"src-1"},
		Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadTransformationRules: %v", err)
	}
	names := func(rules []model.Transformation) []string {
		out := make([]string, 0, len(rules))
		for _, rule := range rules {
			out = append(out, rule.Name()+":"+rule.Column()+">"+rule.To())
		}
		return out
	}
	if got := names(state["public.orders"]); !reflect.DeepEqual(got, []string{"kundenname:name>customer_name", "statusname:status>state"}) {
		t.Fatalf("Regelstand public.orders = %v", got)
	}
	if got := names(state["public.cycle"]); !reflect.DeepEqual(got, []string{"a:x>z"}) {
		t.Fatalf("Regelstand public.cycle = %v, wollen die Regel des letzten Set", got)
	}
	if got, present := state["public.gone"]; present {
		t.Fatalf("Regelstand public.gone = %v, wollen keinen Eintrag (Remove hat die letzte Regel genommen)", got)
	}
	if len(state) != 2 {
		t.Fatalf("Regelstand = %v, wollen genau zwei Tabellen", state)
	}
	if len(exec.queries) != 1 || len(exec.queries[0].args) != 1 || exec.queries[0].args[0] != "src-1" {
		t.Fatalf("Abfragen = %v, wollen eine Abfrage mit der Quelle src-1", exec.queries)
	}
}

// Eine Quelle ohne Regel-Anträge liefert eine leere Map; eine Zeile, deren
// Regelform nicht zu einer Regel führt, endet sichtbar mit der Tabelle im
// Fehler, statt den Stand zu verkürzen.
func TestReadTransformationRulesEmptyAndUnparsable(t *testing.T) {
	recorder := &failRecorder{class: outbound.ErrStorage}
	state, err := sqlexec.ReadTransformationRules(context.Background(), &fakeExecutor{rows: &fakeRows{}}, sqlexec.Statement{
		SQL: "SELECT applied transformation requests", Fail: recorder.fail,
	})
	if err != nil || len(state) != 0 {
		t.Fatalf("ohne Anträge: Stand = %v, Fehler = %v, wollen leer und nil", state, err)
	}

	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{
		{"public", "orders", string(model.AdministrationRequestSetTransformation), "regel", `{"kind": "unbekannt"}`},
	}}}
	_, err = sqlexec.ReadTransformationRules(context.Background(), exec, sqlexec.Statement{
		SQL: "SELECT applied transformation requests", Fail: recorder.fail,
	})
	if !stderrors.Is(err, domainerrors.ErrUnknownTransformationKind) {
		t.Fatalf("Fehler = %v, wollen ErrUnknownTransformationKind", err)
	}
	if !strings.Contains(err.Error(), "public.orders") {
		t.Fatalf("Fehler = %v, wollen die Tabelle genannt", err)
	}
	if len(recorder.causes) != 0 {
		t.Fatalf("die nicht führende Regelform ist kein Treiber-Fehler: %v", recorder.causes)
	}
}

func TestReadTransformationRulesClassifiesFailures(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	recorder := &failRecorder{class: outbound.ErrStorage}
	_, err := sqlexec.ReadTransformationRules(context.Background(), &fakeExecutor{queryErr: cause}, sqlexec.Statement{
		SQL: "SELECT applied transformation requests", Fail: recorder.fail,
	})
	if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Abfrage-Fehler trägt nicht Klasse und Ursache: %v", err)
	}

	scanCause := stderrors.New("Spaltentyp passt nicht")
	recorder = &failRecorder{class: outbound.ErrStorage}
	_, err = sqlexec.ReadTransformationRules(context.Background(), &fakeExecutor{rows: &fakeRows{
		rows:     [][]any{{"public", "orders", "set_transformation", "regel", "{}"}},
		scanErrs: map[int]error{0: scanCause},
	}}, sqlexec.Statement{SQL: "SELECT applied transformation requests", Fail: recorder.fail})
	if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, scanCause) {
		t.Fatalf("Scan-Fehler trägt nicht Klasse und Ursache: %v", err)
	}

	iterCause := stderrors.New("Ergebnis-Menge abgebrochen")
	recorder = &failRecorder{class: outbound.ErrStorage}
	_, err = sqlexec.ReadTransformationRules(context.Background(), &fakeExecutor{rows: &fakeRows{iterErr: iterCause}}, sqlexec.Statement{
		SQL: "SELECT applied transformation requests", Fail: recorder.fail,
	})
	if !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, iterCause) {
		t.Fatalf("Iterations-Fehler trägt nicht Klasse und Ursache: %v", err)
	}
}

// ReadSourceColumns trägt die Spaltennamen in der Ordnung der Abfrage; die
// Übergabe der Argumente der Abfrage und die Fehlerklassifikation laufen wie
// bei den übrigen Lesern.
func TestReadSourceColumnsCarriesNamesInOrder(t *testing.T) {
	exec := &fakeExecutor{rows: &fakeRows{rows: [][]any{{"id"}, {"Name"}, {"status"}}}}
	recorder := &failRecorder{class: outbound.ErrStorage}

	columns, err := sqlexec.ReadSourceColumns(context.Background(), exec, sqlexec.Statement{
		SQL: "SELECT columns", Args: []any{"public", "orders"}, Fail: recorder.fail,
	})
	if err != nil {
		t.Fatalf("ReadSourceColumns: %v", err)
	}
	if !reflect.DeepEqual(columns, []string{"id", "Name", "status"}) {
		t.Fatalf("Spalten = %v", columns)
	}
	if len(exec.queries) != 1 || !reflect.DeepEqual(exec.queries[0].args, []any{"public", "orders"}) {
		t.Fatalf("Abfragen = %v", exec.queries)
	}

	empty, err := sqlexec.ReadSourceColumns(context.Background(), &fakeExecutor{rows: &fakeRows{}}, sqlexec.Statement{SQL: "SELECT columns", Fail: recorder.fail})
	if err != nil || empty == nil || len(empty) != 0 {
		t.Fatalf("Tabelle ohne Spalten: %v, %v, wollen leere Liste und nil", empty, err)
	}
}

func TestReadSourceColumnsClassifiesFailures(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	recorder := &failRecorder{class: outbound.ErrStorage}
	if _, err := sqlexec.ReadSourceColumns(context.Background(), &fakeExecutor{queryErr: cause}, sqlexec.Statement{SQL: "SELECT columns", Fail: recorder.fail}); !stderrors.Is(err, outbound.ErrStorage) || !stderrors.Is(err, cause) {
		t.Fatalf("Abfrage-Fehler trägt nicht Klasse und Ursache: %v", err)
	}
	scanCause := stderrors.New("Spaltentyp passt nicht")
	if _, err := sqlexec.ReadSourceColumns(context.Background(), &fakeExecutor{rows: &fakeRows{rows: [][]any{{"id"}}, scanErrs: map[int]error{0: scanCause}}}, sqlexec.Statement{SQL: "SELECT columns", Fail: recorder.fail}); !stderrors.Is(err, scanCause) {
		t.Fatalf("Scan-Fehler trägt nicht die Ursache: %v", err)
	}
	iterCause := stderrors.New("Ergebnis-Menge abgebrochen")
	if _, err := sqlexec.ReadSourceColumns(context.Background(), &fakeExecutor{rows: &fakeRows{iterErr: iterCause}}, sqlexec.Statement{SQL: "SELECT columns", Fail: recorder.fail}); !stderrors.Is(err, iterCause) {
		t.Fatalf("Iterations-Fehler trägt nicht die Ursache: %v", err)
	}
}
