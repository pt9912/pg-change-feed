package model

import (
	"testing"

	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// changeArgs trägt die Argumente eines gültigen Changes; je Test weicht
// genau ein Feld ab.
type changeArgs struct {
	id      ChangeID
	tx      TransactionID
	table   SourceTableID
	seq     int64
	op      Operation
	oldData []byte
	newData []byte
	sv      SchemaVersionID
}

func validChangeArgs() changeArgs {
	return changeArgs{
		id:      "ch-1",
		tx:      "tx-1",
		table:   "tbl-1",
		seq:     1,
		op:      OperationUpdate,
		oldData: []byte(`{"a":1}`),
		newData: []byte(`{"a":2}`),
		sv:      "sv-1",
	}
}

func buildChange(t *testing.T, args changeArgs) Change {
	t.Helper()
	change, err := NewChange(args.id, args.tx, args.table, args.seq, args.op, args.oldData, args.newData, args.sv)
	if err != nil {
		t.Fatalf("Change-Konstruktor: %v", err)
	}
	return change
}

// LH-FA-DAT-001, Happy Path: zwei erfasste Changes sind über ihre
// Identifikatoren unterscheidbar.
func TestLHFADAT001ChangesAreDistinguishable(t *testing.T) {
	first := buildChange(t, validChangeArgs())
	second := buildChange(t, changeArgs{id: "ch-2", tx: "tx-1", table: "tbl-1", seq: 2, op: OperationUpdate, oldData: []byte(`{"a":1}`), newData: []byte(`{"a":2}`), sv: "sv-1"})
	if first.ID == second.ID {
		t.Fatalf("zwei erfasste Changes tragen dieselbe Kennung: %q", first.ID)
	}
}

// LH-FA-DAT-001, Boundary: derselbe Datensatz, zweimal identisch geändert,
// ergibt zwei unterscheidbare Changes — identische Bilder, verschiedene
// Kennung und Sequenz.
func TestLHFADAT001SameRowChangedTwiceYieldsDistinctChanges(t *testing.T) {
	first := buildChange(t, validChangeArgs())
	second := buildChange(t, changeArgs{id: "ch-2", tx: "tx-1", table: "tbl-1", seq: 2, op: OperationUpdate, oldData: []byte(`{"a":1}`), newData: []byte(`{"a":2}`), sv: "sv-1"})
	if string(first.NewImage) != string(second.NewImage) || string(first.OldImage) != string(second.OldImage) {
		t.Fatalf("Boundary verletzt: die Änderung ist identisch, Bilder %s/%s vs %s/%s", first.OldImage, first.NewImage, second.OldImage, second.NewImage)
	}
	if first.ID == second.ID || first.Sequence == second.Sequence {
		t.Fatalf("zwei identische Änderungen sind nicht unterscheidbar: Kennung %q und Sequenz %d doppelt", first.ID, first.Sequence)
	}
}

// Die Konstruktoren erzwingen die Change-Invarianten (`ADR-0029`, Regel 6
// und 7): illegale Zustände scheitern am Konstruktor.
func TestNewChangeRejectsInvariantViolations(t *testing.T) {
	cases := []struct {
		name      string
		mutate    func(args *changeArgs)
		wantError error
	}{
		{"leere Change-Kennung", func(a *changeArgs) { a.id = "" }, domainerrors.ErrEmptyIdentifier},
		{"leere Transaktions-Kennung", func(a *changeArgs) { a.tx = "" }, domainerrors.ErrEmptyIdentifier},
		{"leere Tabellen-Kennung", func(a *changeArgs) { a.table = "" }, domainerrors.ErrEmptyIdentifier},
		{"leere Schema-Version", func(a *changeArgs) { a.sv = "" }, domainerrors.ErrEmptyIdentifier},
		{"Sequenz 0", func(a *changeArgs) { a.seq = 0 }, domainerrors.ErrNonPositiveSequence},
		{"negative Sequenz", func(a *changeArgs) { a.seq = -3 }, domainerrors.ErrNonPositiveSequence},
		{"unbekannte Operation", func(a *changeArgs) { a.op = Operation("TRUNCATE") }, domainerrors.ErrInvalidOperation},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := validChangeArgs()
			tc.mutate(&args)
			_, err := NewChange(args.id, args.tx, args.table, args.seq, args.op, args.oldData, args.newData, args.sv)
			if !stderrors.Is(err, tc.wantError) {
				t.Fatalf("Fehler = %v, wollen %v", err, tc.wantError)
			}
		})
	}
}

// SPEC-002: der Konstruktor legt einen Change mit der
// Herkunft `wal` an.
func TestNewChangeDefaultsOriginToWAL(t *testing.T) {
	change := buildChange(t, validChangeArgs())
	if change.Origin != ChangeOriginWAL {
		t.Fatalf("Origin = %q, wollen %q", change.Origin, ChangeOriginWAL)
	}
}

// Die Herkunft ist eine geschlossene Menge (`SPEC-002`): `wal` und
// `backfill` werden angenommen, jeder andere Wert abgelehnt; die leere
// Zeichenkette (fehlender Wert, `NULL`) liest als `wal`.
func TestNewChangeOriginClosedSet(t *testing.T) {
	cases := []struct {
		raw     string
		want    ChangeOrigin
		wantErr error
	}{
		{"wal", ChangeOriginWAL, nil},
		{"backfill", ChangeOriginBackfill, nil},
		{"", ChangeOriginWAL, nil},
		{"snapshot", "", domainerrors.ErrInvalidChangeOrigin},
		{"WAL", "", domainerrors.ErrInvalidChangeOrigin},
		{" backfill", "", domainerrors.ErrInvalidChangeOrigin},
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := NewChangeOrigin(tc.raw)
			if !stderrors.Is(err, tc.wantErr) {
				t.Fatalf("Fehler = %v, wollen %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("Herkunft = %q, wollen %q", got, tc.want)
			}
		})
	}
}

// `WithOrigin` setzt die Herkunft und lehnt einen Wert außerhalb der Menge
// ab; der Change bleibt dann unverändert.
func TestChangeWithOrigin(t *testing.T) {
	change := buildChange(t, validChangeArgs())

	backfill, err := change.WithOrigin(ChangeOriginBackfill)
	if err != nil {
		t.Fatalf("WithOrigin(backfill): %v", err)
	}
	if backfill.Origin != ChangeOriginBackfill {
		t.Fatalf("Origin = %q, wollen %q", backfill.Origin, ChangeOriginBackfill)
	}
	if change.Origin != ChangeOriginWAL {
		t.Fatalf("der Ursprungs-Change bleibt unverändert: Origin = %q", change.Origin)
	}

	rejected, err := change.WithOrigin(ChangeOrigin("snapshot"))
	if !stderrors.Is(err, domainerrors.ErrInvalidChangeOrigin) {
		t.Fatalf("Fehler = %v, wollen %v", err, domainerrors.ErrInvalidChangeOrigin)
	}
	if rejected.Origin != ChangeOriginWAL {
		t.Fatalf("abgelehnter Wert verändert den Change: Origin = %q", rejected.Origin)
	}
}

// Ein fehlender Wert liest als `wal`, ein gesetzter bleibt.
func TestChangeOriginOrDefault(t *testing.T) {
	if got := ChangeOrigin("").OrDefault(); got != ChangeOriginWAL {
		t.Fatalf("leere Herkunft liest als %q, wollen %q", got, ChangeOriginWAL)
	}
	if got := ChangeOriginBackfill.OrDefault(); got != ChangeOriginBackfill {
		t.Fatalf("gesetzte Herkunft bleibt: %q", got)
	}
}

// TestChangeMatchesFilter trägt die Kombinatorik aus `ADR-0133` Teilfrage 1:
// die Dimensionen sind unabhängig optional, alle leer lässt jeden Change
// passieren, gesetzte Dimensionen verknüpfen sich als Konjunktion.
// Rot färbende Mutation: eine der drei `if`-Bedingungen streichen — dann
// passiert ein nicht passender Change trotz gesetztem Filter auf dieser
// Dimension.
func TestChangeMatchesFilter(t *testing.T) {
	change := buildChange(t, validChangeArgs())
	change.Schema = "public"
	change.Table = "orders"

	cases := []struct {
		name         string
		filterSchema string
		filterTable  string
		want         bool
	}{
		{"kein Filter", "", "", true},
		{"Schema passt, kein Table-Filter", "public", "", true},
		{"Schema passt nicht", "other", "", false},
		{"Table passt, kein Schema-Filter", "", "orders", true},
		{"Table passt nicht", "", "customers", false},
		{"beide passen", "public", "orders", true},
		{"Schema passt, Table nicht", "public", "customers", false},
		{"Schema passt nicht, Table passt", "other", "orders", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := change.MatchesFilter(tc.filterSchema, tc.filterTable, ""); got != tc.want {
				t.Fatalf("MatchesFilter(%q, %q, \"\") = %v, wollen %v", tc.filterSchema, tc.filterTable, got, tc.want)
			}
		})
	}
}

// TestChangeMatchesFilterTarget trägt die Zieldimension (`ADR-0137`
// Teilfrage 5): ein leeres `target` wählt nicht aus (auch eine geroutete
// Change passiert), ein gesetztes trifft genau das Ziel der Change; eine
// Change ohne Ziel trifft nie ein gesetztes `target`; die Dimension
// verknüpft sich als Konjunktion mit `schema`/`table`.
// Rot färbende Mutation: die `target`-Bedingung streichen (Fälle „anderes
// Ziel", „kein Ziel") oder `!=` zu `==` kehren (alle Fälle).
func TestChangeMatchesFilterTarget(t *testing.T) {
	routed := buildChange(t, validChangeArgs())
	routed.Schema = "public"
	routed.Table = "orders"
	routed.RouteTarget = "eu"
	unrouted := routed
	unrouted.RouteTarget = ""

	cases := []struct {
		name   string
		change Change
		schema string
		table  string
		target string
		want   bool
	}{
		{"geroutet, kein Filter", routed, "", "", "", true},
		{"ungeroutet, kein Filter", unrouted, "", "", "", true},
		{"Ziel passt", routed, "", "", "eu", true},
		{"anderes Ziel", routed, "", "", "us", false},
		{"Präfix des Ziels ist kein Treffer", routed, "", "", "e", false},
		{"kein Ziel an der Change", unrouted, "", "", "eu", false},
		{"Ziel passt, Tabelle passt", routed, "public", "orders", "eu", true},
		{"Ziel passt, Tabelle passt nicht", routed, "public", "customers", "eu", false},
		{"Ziel passt nicht, Tabelle passt", routed, "public", "orders", "us", false},
		{"Name außerhalb des Alphabets", routed, "", "", "EU", false},
		{"Name mit U+0000", routed, "", "", "eu\x00", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.change.MatchesFilter(tc.schema, tc.table, tc.target); got != tc.want {
				t.Fatalf("MatchesFilter(%q, %q, %q) = %v, wollen %v", tc.schema, tc.table, tc.target, got, tc.want)
			}
		})
	}
}
