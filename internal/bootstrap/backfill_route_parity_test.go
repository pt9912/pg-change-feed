package bootstrap

import (
	"context"
	"fmt"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/backfill"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Ziel-Parität von WAL- und Backfill-Pfad (`LH-FA-CFG-008`): dieselbe
// Zeile und dieselbe Regelliste liefern in beiden Erzeugungspfaden dasselbe
// `route_target`. Der WAL-Pfad ist der echte `mapper.Assembler`, der
// Backfill-Pfad der echte `BackfillTableService`; die Ports des Runs sind
// Fakes (Aufbau wie in `backfill_image_parity_test.go`).

func parityRoute(t *testing.T, name, target string, order int64, when *model.RouteCondition) model.RouteRule {
	t.Helper()
	rule, err := model.NewRouteRule(name, target, order, when)
	if err != nil {
		t.Fatalf("NewRouteRule(%q): %v", name, err)
	}
	return rule
}

// TestBackfillAndWALRouteTargetsAreEqual trägt die Parität über vier
// Regellisten (eine bedingte Regel, Auffangregel vor der bedingten in der
// Liste, Regel auf den leeren Text, keine Regel), mit und ohne
// Transformation der Spalte der Bedingung und mit und ohne Ausschluss einer
// anderen Spalte: je Zeile — Treffer, abwesender Wert, leerer Text, Nicht-
// Treffer — ist das Ziel beider Pfade gleich dem erwarteten (die Erwartung ist
// eine feste Tabelle, nicht die Auswertung selbst). Rot färbende Mutationen
// (je eine): im Run `nil` statt der Regelliste an `EvaluateRoute` (Backfill-
// Ziele leer, Parität und Erwartung brechen); die Anwendbarkeitsprüfung des
// Runs gegen `nil` statt der Snapshot-Spalten (der Run endet `failed`).
func TestBackfillAndWALRouteTargetsAreEqual(t *testing.T) {
	const (
		source    = model.SourceID("src-route-parity")
		qualified = "public.parity"
	)
	columns := []string{"id", "region", "note"}
	rows := [][]*string{
		{parityValue("1"), parityValue("eu"), parityValue("plain")},
		{parityValue("2"), nil, parityValue("n")},
		{parityValue("3"), parityValue(""), nil},
		{parityValue("4"), parityValue("us"), parityValue("x")},
	}
	eu := func(t *testing.T) model.RouteRule {
		return parityRoute(t, "eu", "eu_ziel", 10, &model.RouteCondition{Column: "region", Equals: "eu"})
	}
	cases := []struct {
		name  string
		rules func(t *testing.T) []model.RouteRule
		want  []model.RouteTarget
	}{
		{"eine bedingte Regel", func(t *testing.T) []model.RouteRule { return []model.RouteRule{eu(t)} },
			[]model.RouteTarget{"eu_ziel", "", "", ""}},
		{"Auffangregel vor der bedingten in der Liste", func(t *testing.T) []model.RouteRule {
			return []model.RouteRule{parityRoute(t, "alle", "alle", 50, nil), eu(t)}
		}, []model.RouteTarget{"eu_ziel", "alle", "alle", "alle"}},
		{"Regel auf den leeren Text", func(t *testing.T) []model.RouteRule {
			return []model.RouteRule{parityRoute(t, "leer", "leer_ziel", 5, &model.RouteCondition{Column: "region", Equals: ""}), eu(t)}
		}, []model.RouteTarget{"eu_ziel", "", "leer_ziel", ""}},
		{"keine Regel", func(t *testing.T) []model.RouteRule { return nil },
			[]model.RouteTarget{"", "", "", ""}},
	}
	transformations := map[string][]model.Transformation{
		"ohne Transformation":                   nil,
		"rename_column an der Bedingungsspalte": {mustRename(t, "region_rename", "region", "area")},
	}
	for _, tc := range cases {
		for transformationName, transformation := range transformations {
			for _, excluded := range [][]string{nil, {"note"}} {
				t.Run(fmt.Sprintf("%s/%s/ausgeschlossen %v", tc.name, transformationName, excluded), func(t *testing.T) {
					rules := tc.rules(t)
					table, err := model.NewSourceTable("tbl-route-parity", source, "public", "parity")
					if err != nil {
						t.Fatal(err)
					}
					version, err := model.NewSchemaVersion("tbl-route-parity-v1", table.ID, 1)
					if err != nil {
						t.Fatal(err)
					}
					ctx := context.Background()

					assembler, err := mapper.NewAssembler(source, map[string]mapper.TableBinding{
						qualified: {TableID: table.ID, SchemaVersion: version.ID, ExcludedColumns: excluded, Transformations: transformation, Routes: rules},
					}, nil)
					if err != nil {
						t.Fatalf("NewAssembler: %v", err)
					}
					relationColumns := make([]decode.Column, len(columns))
					for i, name := range columns {
						relationColumns[i] = decode.Column{Name: name, Key: i == 0}
					}
					relation := &decode.Relation{Schema: "public", Name: "parity", Columns: relationColumns}
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

					writer := &parityWriter{}
					service := backfill.NewBackfillTableService(backfill.Ports{
						Activation:      &fakeTableActivationPort{registered: map[string]model.SourceTable{qualified: table}},
						Exclusion:       &fakeColumnExclusionPort{excluded: map[string][]string{qualified: excluded}},
						Transformations: &fakeTransformationPort{rules: map[string][]model.Transformation{qualified: transformation}},
						Routing:         &fakeRoutingPort{rules: map[string][]model.RouteRule{qualified: rules}},
						Schemas:         &fakeSchemaStorePort{versions: map[model.SourceTableID]model.SchemaVersion{table.ID: version}},
						Snapshot:        &paritySnapshotPort{snapshot: &paritySnapshot{columns: columns, rows: rows}},
						Runs:            &fakeBackfillRunPort{},
						Writer:          writer,
						Clock:           systemclock.New(),
					})
					queued, err := model.NewQueuedBackfillRun("route-parity-run", source, "public", "parity", systemclock.New().Now())
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
					for i := range rows {
						if walChanges[i].RouteTarget != tc.want[i] || writer.changes[i].RouteTarget != tc.want[i] {
							t.Fatalf("Zeile %d: WAL-Ziel %q, Backfill-Ziel %q, erwartet %q", i+1, walChanges[i].RouteTarget, writer.changes[i].RouteTarget, tc.want[i])
						}
					}
				})
			}
		}
	}
}

func mustRename(t *testing.T, name, column, to string) model.Transformation {
	t.Helper()
	rule, err := model.NewRenameColumn(name, column, to)
	if err != nil {
		t.Fatalf("NewRenameColumn: %v", err)
	}
	return rule
}
