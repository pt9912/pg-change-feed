package bootstrap

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/backfill"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Bild-Parität von WAL- und Backfill-Pfad mit Regelstand (`LH-FA-CFG-007`,
// `ADR-0112` Folgepflicht 7, `ADR-0115` Festlegung 4): dieselbe Zeile und
// derselbe Regelstand liefern in beiden Erzeugungspfaden byte-gleiche
// Row Images am Ausgang der gemeinsamen Bild-Konstruktion. Die Composition
// Root ist der einzige Ort, der beide Pfade zugleich importieren darf
// (`.a-check.yml`, `composition_root`): der WAL-Pfad ist der echte
// `mapper.Assembler`, der Backfill-Pfad der echte `BackfillTableService`; die
// Ports des Runs sind Fakes. Der reale Rundlauf gegen PostgreSQL steht in
// `make test-integration`.

// paritySnapshot liefert die Zeilen in einem Block.
type paritySnapshot struct {
	columns []string
	rows    [][]*string
	served  bool
}

func (s *paritySnapshot) Offset() uint64    { return 0x100 }
func (s *paritySnapshot) Columns() []string { return s.columns }
func (s *paritySnapshot) NextBlock(context.Context) ([][]*string, error) {
	if s.served {
		return nil, nil
	}
	s.served = true
	return s.rows, nil
}
func (s *paritySnapshot) Close(context.Context) error { return nil }

type paritySnapshotPort struct{ snapshot *paritySnapshot }

func (p *paritySnapshotPort) OpenSnapshot(context.Context, string, string, string) (outbound.TableSnapshot, error) {
	return p.snapshot, nil
}
func (p *paritySnapshotPort) EstimatedRows(context.Context, string, string) (int64, bool, error) {
	return 0, false, nil
}

// parityWriter hält die Changes der angehängten Blöcke.
type parityWriter struct{ changes []model.Change }

func (w *parityWriter) Begin(context.Context, model.BackfillRun) (outbound.BackfillTransaction, error) {
	return w, nil
}
func (w *parityWriter) AppendBlock(_ context.Context, block *model.ChangeTransaction) error {
	changes, err := block.Changes()
	if err != nil {
		return err
	}
	w.changes = append(w.changes, changes...)
	return nil
}
func (w *parityWriter) Commit(context.Context, model.BackfillRun) error { return nil }
func (w *parityWriter) Rollback(context.Context) error                  { return nil }

// parityRule bildet einen Regeltyp der Domäne auf eine anwendbare Regel für
// `column` ab; ein Regeltyp ohne Fall bricht den Test ab.
func parityRule(t *testing.T, kind model.TransformationKind, column string) model.Transformation {
	t.Helper()
	switch kind {
	case model.TransformationRenameColumn:
		rule, err := model.NewRenameColumn("parity-"+column, column, column+"_umbenannt")
		if err != nil {
			t.Fatalf("NewRenameColumn: %v", err)
		}
		return rule
	}
	t.Fatalf("Regeltyp %q ohne Fall in diesem Test", kind)
	return model.Transformation{}
}

func parityValue(v string) *string { return &v }

// TestBackfillAndWALImagesAreByteEqualWithRules trägt die Parität: über die
// Regeltypen der Domäne (`model.TransformationKinds`), jede Spalte als Ziel der
// Regel und jeden Ausschluss (keiner, jede Spalte) sind die Bilder des Runs
// und der `Assembler`-Bilder derselben Zeilen — mit NULL-Wert, Anführungszeichen,
// Backslash, Zeilenumbruch und Nicht-ASCII-Zeichen — byte-gleich; die Regel
// wirkt im Bild (der Zielname steht dort, wo die Spalte einen Wert und keinen
// Ausschluss trägt). Rot färbende Mutationen (je eine): `nil` statt des
// Regelsatzes an `BuildRowImage` im Run (Bild ohne Zielname, Paritätsfehler);
// `nil` statt `binding.Transformations` an `BuildRowImage` im `Assembler`
// (dasselbe von der anderen Seite).
func TestBackfillAndWALImagesAreByteEqualWithRules(t *testing.T) {
	const (
		source    = model.SourceID("src-parity")
		qualified = "public.parity"
	)
	columns := []string{"id", "name", "note"}
	rows := [][]*string{
		{parityValue("1"), parityValue("Ada"), parityValue("plain")},
		{parityValue("2"), nil, parityValue(`quote " backslash \ newline` + "\n" + `umlaut äöü €`)},
		{parityValue("3"), parityValue(""), nil},
	}
	exclusions := append([][]string{nil}, [][]string{{"id"}, {"name"}, {"note"}}...)
	for _, kind := range model.TransformationKinds() {
		for _, ruleColumn := range columns {
			for _, excluded := range exclusions {
				t.Run(fmt.Sprintf("%s/Regel an %s/ausgeschlossen %v", kind, ruleColumn, excluded), func(t *testing.T) {
					rule := parityRule(t, kind, ruleColumn)
					rules := []model.Transformation{rule}

					table, err := model.NewSourceTable("tbl-parity", source, "public", "parity")
					if err != nil {
						t.Fatal(err)
					}
					version, err := model.NewSchemaVersion("tbl-parity-v1", table.ID, 1)
					if err != nil {
						t.Fatal(err)
					}

					// WAL-Pfad: der echte Assembler mit der Bindung des Regelstands.
					assembler, err := mapper.NewAssembler(source, map[string]mapper.TableBinding{
						qualified: {TableID: table.ID, SchemaVersion: version.ID, ExcludedColumns: excluded, Transformations: rules},
					}, nil)
					if err != nil {
						t.Fatalf("NewAssembler: %v", err)
					}
					relationColumns := make([]decode.Column, len(columns))
					for i, name := range columns {
						relationColumns[i] = decode.Column{Name: name, Key: i == 0}
					}
					relation := &decode.Relation{Schema: "public", Name: "parity", Columns: relationColumns}
					ctx := context.Background()
					if _, err := assembler.Consume(ctx, decode.Begin{XID: 7}); err != nil {
						t.Fatalf("Begin: %v", err)
					}
					for _, row := range rows {
						if _, err := assembler.Consume(ctx, decode.Change{Relation: relation, Operation: decode.OpInsert, New: row}); err != nil {
							t.Fatalf("Change: %v", err)
						}
					}
					command, err := assembler.Consume(ctx, decode.Commit{CommitLSN: 700})
					if err != nil {
						t.Fatalf("Commit: %v", err)
					}
					walChanges, err := command.Transaction.Changes()
					if err != nil {
						t.Fatalf("Changes: %v", err)
					}

					// Backfill-Pfad: der echte Use Case mit dem Regelstand über den Port.
					writer := &parityWriter{}
					service := backfill.NewBackfillTableService(backfill.Ports{
						Activation:      &fakeTableActivationPort{registered: map[string]model.SourceTable{qualified: table}},
						Exclusion:       &fakeColumnExclusionPort{excluded: map[string][]string{qualified: excluded}},
						Transformations: &fakeTransformationPort{rules: map[string][]model.Transformation{qualified: rules}},
						Schemas:         &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{table.ID: version}},
						Snapshot:        &paritySnapshotPort{snapshot: &paritySnapshot{columns: columns, rows: rows}},
						Runs:            &fakeBackfillRunPort{},
						Writer:          writer,
						Clock:           systemclock.New(),
					})
					queued, err := model.NewQueuedBackfillRun("parity-run", source, "public", "parity", systemclock.New().Now())
					if err != nil {
						t.Fatal(err)
					}
					result, err := service.Execute(ctx, inbound.BackfillExecuteCommand{Run: queued, Publication: "pub"})
					if err != nil || result.Run.Status != model.BackfillRunCompleted {
						t.Fatalf("Execute = %+v, %v, will completed", result.Run, err)
					}

					if len(walChanges) != len(rows) || len(writer.changes) != len(rows) {
						t.Fatalf("WAL-Changes %d, Backfill-Changes %d, will je %d", len(walChanges), len(writer.changes), len(rows))
					}
					for i, row := range rows {
						walImage, backfillImage := walChanges[i].NewImage, writer.changes[i].NewImage
						if !bytes.Equal(walImage, backfillImage) {
							t.Fatalf("Zeile %d: WAL-Bild %s ≠ Backfill-Bild %s", i+1, walImage, backfillImage)
						}
						carriesTarget := bytes.Contains(backfillImage, []byte(`"`+rule.To()+`"`))
						wantTarget := row[indexOf(columns, ruleColumn)] != nil && !(len(excluded) == 1 && excluded[0] == ruleColumn)
						if carriesTarget != wantTarget {
							t.Fatalf("Zeile %d: Bild %s, Zielname %q vorhanden = %t, will %t", i+1, backfillImage, rule.To(), carriesTarget, wantTarget)
						}
					}
				})
			}
		}
	}
}

func indexOf(names []string, name string) int {
	for i, candidate := range names {
		if candidate == name {
			return i
		}
	}
	return -1
}
