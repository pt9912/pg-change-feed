package mapper_test

import (
	"context"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// benchmarkAssemblerChange misst den heißen Pfad des Erfassungspfads: eine
// Änderung mit zwölf Spalten über `Assembler.Consume` (Begin, Change,
// Commit) gegen eine Bindung mit dem übergebenen Regelstand.
func benchmarkAssemblerChange(b *testing.B, rules ...model.Transformation) {
	columns := make([]decode.Column, 12)
	values := make([]*string, len(columns))
	for i := range columns {
		name := "col_" + string(rune('a'+i))
		columns[i] = decode.Column{Name: name, Key: i == 0}
		value := "wert-" + name
		values[i] = &value
	}
	relationEvent := relation("public", "feed", columns...)
	assembler, err := mapper.NewAssembler("src-1", map[string]mapper.TableBinding{
		"public.feed": {TableID: "tbl-1", SchemaVersion: "sv-1", Transformations: rules},
	}, nil)
	if err != nil {
		b.Fatalf("NewAssembler: %v", err)
	}
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := assembler.Consume(ctx, decode.Begin{XID: uint32(i + 1)}); err != nil {
			b.Fatalf("Begin: %v", err)
		}
		if _, err := assembler.Consume(ctx, decode.Change{Relation: relationEvent, Operation: decode.OpInsert, New: values}); err != nil {
			b.Fatalf("Change: %v", err)
		}
		if _, err := assembler.Consume(ctx, decode.Commit{CommitLSN: uint64(i + 1)}); err != nil {
			b.Fatalf("Commit: %v", err)
		}
	}
}

// BenchmarkAssemblerChange ist der Ausgangswert: Bindung ohne Regelstand.
func BenchmarkAssemblerChange(b *testing.B) {
	benchmarkAssemblerChange(b)
}

// BenchmarkAssemblerChangeWithRules trägt zwei Umbenennungen an der Bindung:
// die Kosten der Anwendbarkeits-Prüfung und der Auswertung je Änderung.
func BenchmarkAssemblerChangeWithRules(b *testing.B) {
	first, err := model.NewRenameColumn("erste", "col_b", "ziel_b")
	if err != nil {
		b.Fatalf("NewRenameColumn: %v", err)
	}
	second, err := model.NewRenameColumn("zweite", "col_g", "ziel_g")
	if err != nil {
		b.Fatalf("NewRenameColumn: %v", err)
	}
	benchmarkAssemblerChange(b, first, second)
}
