package snapshotlogic_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgressnapshot/snapshotlogic"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

const validDSN = "postgres://u:p@localhost:5432/db?sslmode=disable"

func TestValidate(t *testing.T) {
	config, err := snapshotlogic.Validate(validDSN, 1, time.Second)
	if err != nil || config == nil || config.Host != "localhost" || config.Database != "db" {
		t.Fatalf("gültige Parameter: %v, %v", config, err)
	}
	for name, run := range map[string]func() error{
		"leerer DSN":        func() error { _, err := snapshotlogic.Validate("", 1, time.Second); return err },
		"DSN nicht lesbar":  func() error { _, err := snapshotlogic.Validate("postgres://%zz", 1, time.Second); return err },
		"Blockgröße 0":      func() error { _, err := snapshotlogic.Validate(validDSN, 0, time.Second); return err },
		"Zeitlimit 0":       func() error { _, err := snapshotlogic.Validate(validDSN, 1, 0); return err },
		"Zeitlimit negativ": func() error { _, err := snapshotlogic.Validate(validDSN, 1, -time.Second); return err },
	} {
		if err := run(); !errors.Is(err, outbound.ErrSnapshotConfiguration) {
			t.Errorf("%s: %v, erwartet Klasse configuration", name, err)
		}
	}
}

func TestSlotName(t *testing.T) {
	got, err := snapshotlogic.SlotName("0a1b-2c3d")
	if err != nil || got != "cdc_bf_0a1b2c3d" {
		t.Fatalf("SlotName = %q, %v; erwartet cdc_bf_0a1b2c3d", got, err)
	}
	if _, err := snapshotlogic.SlotName(strings.Repeat("a", 56)); err != nil {
		t.Errorf("56 Zeichen (Slot-Name genau 63): %v", err)
	}
	for _, bad := range []string{"", "-", "GROSS", "mit leerzeichen", "a.b", strings.Repeat("a", 57)} {
		if _, err := snapshotlogic.SlotName(bad); !errors.Is(err, outbound.ErrSnapshotConfiguration) {
			t.Errorf("SlotName(%q): %v, erwartet Klasse configuration", bad, err)
		}
	}
}

func TestValidSnapshotName(t *testing.T) {
	for _, ok := range []string{"00000003-0000001B-1", "0A-1f"} {
		if !snapshotlogic.ValidSnapshotName(ok) {
			t.Errorf("ValidSnapshotName(%q) = false", ok)
		}
	}
	for _, bad := range []string{"", "x", "00000003-0000001B-1'; DROP TABLE t; --", "0000 0003", "0000003\n"} {
		if snapshotlogic.ValidSnapshotName(bad) {
			t.Errorf("ValidSnapshotName(%q) = true", bad)
		}
	}
}

func TestQuoteIdent(t *testing.T) {
	for in, want := range map[string]string{
		"plain":   `"plain"`,
		"Mixed C": `"Mixed C"`,
		`we"ird`:  `"we""ird"`,
		`""`:      `""""""`,
		"":        `""`,
	} {
		if got := snapshotlogic.QuoteIdent(in); got != want {
			t.Errorf("QuoteIdent(%q) = %s, erwartet %s", in, got, want)
		}
	}
}

// TestCursorStatementSelectsOnlyQuotedNames trägt die Text-Ergebnisformat-
// Zusage an ihrer Anweisungsseite (`ADR-0115` Festlegung 1): an einfachen
// Namen steht kein Cast und kein Funktionsaufruf in der Anweisung.
func TestCursorStatementSelectsOnlyQuotedNames(t *testing.T) {
	got := snapshotlogic.CursorStatement([]string{"a", "b"}, "s", "t")
	want := `DECLARE cdc_bf_cursor NO SCROLL CURSOR FOR SELECT "a", "b" FROM "s"."t"`
	if got != want {
		t.Fatalf("CursorStatement = %s, erwartet %s", got, want)
	}
	for _, forbidden := range []string{"::", "(", ")", "CAST", "CASE"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("die Cursor-Anweisung trägt %q: %s", forbidden, got)
		}
	}
}

// TestCursorStatementQuotesEveryIdentifier bindet das Quoting von Schema,
// Tabelle und Spalten an Namen mit Anführungszeichen, Leerzeichen und
// Großbuchstaben.
func TestCursorStatementQuotesEveryIdentifier(t *testing.T) {
	got := snapshotlogic.CursorStatement([]string{`Id`, `we "ird`}, `Sch "ema x`, `Mixed Case "y`)
	want := `DECLARE cdc_bf_cursor NO SCROLL CURSOR FOR SELECT "Id", "we ""ird" FROM "Sch ""ema x"."Mixed Case ""y"`
	if got != want {
		t.Fatalf("CursorStatement = %s, erwartet %s", got, want)
	}
}

func TestFetchStatement(t *testing.T) {
	if got, want := snapshotlogic.FetchStatement(1000), "FETCH FORWARD 1000 FROM cdc_bf_cursor"; got != want {
		t.Fatalf("FetchStatement = %s, erwartet %s", got, want)
	}
}

func TestValues(t *testing.T) {
	values := snapshotlogic.Values([][]byte{[]byte("t"), nil, {}, []byte("a b")})
	if len(values) != 4 || values[0] == nil || *values[0] != "t" || values[1] != nil ||
		values[2] == nil || *values[2] != "" || values[3] == nil || *values[3] != "a b" {
		t.Fatalf("Values = %v, erwartet [t nil \"\" \"a b\"] (NULL und leerer Text bleiben verschieden)", values)
	}
}

func TestEstimate(t *testing.T) {
	cases := []struct {
		in    string
		rows  int64
		known bool
	}{
		{"-1", 0, false},
		{"-0.5", 0, false},
		{"0", 0, true}, // analysierte leere Tabelle: bekannt null, nie „unbekannt"
		{"40", 40, true},
		{"39.6", 40, true},
		{"1.5e+06", 1500000, true},
	}
	for _, tc := range cases {
		rows, known, err := snapshotlogic.Estimate(tc.in)
		if err != nil || rows != tc.rows || known != tc.known {
			t.Errorf("Estimate(%q) = %d, %v, %v; erwartet %d, %v", tc.in, rows, known, err, tc.rows, tc.known)
		}
	}
	for _, bad := range []string{"", "abc"} {
		if _, _, err := snapshotlogic.Estimate(bad); !errors.Is(err, outbound.ErrSnapshotStorage) {
			t.Errorf("Estimate(%q): %v, erwartet Klasse storage", bad, err)
		}
	}
}

func TestClassify(t *testing.T) {
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	live := context.Background()
	fallback := outbound.ErrSnapshotReplication
	cases := []struct {
		name  string
		ctx   context.Context
		cause error
		want  error
	}{
		{"Kontext beendet", cancelled, errors.New("x"), outbound.ErrSnapshotTransient},
		{"Zeitlimit als Ursache", live, fmt.Errorf("wrap: %w", context.DeadlineExceeded), outbound.ErrSnapshotTransient},
		{"Abbruch als Ursache", live, context.Canceled, outbound.ErrSnapshotTransient},
		{"42501", live, &pgconn.PgError{Code: "42501"}, outbound.ErrSnapshotPermission},
		{"28P01", live, &pgconn.PgError{Code: "28P01"}, outbound.ErrSnapshotPermission},
		{"28000", live, &pgconn.PgError{Code: "28000"}, outbound.ErrSnapshotPermission},
		{"57P01", live, &pgconn.PgError{Code: "57P01"}, outbound.ErrSnapshotTransient},
		{"57P02", live, &pgconn.PgError{Code: "57P02"}, outbound.ErrSnapshotTransient},
		{"57P03", live, &pgconn.PgError{Code: "57P03"}, outbound.ErrSnapshotTransient},
		{"53400", live, &pgconn.PgError{Code: "53400"}, outbound.ErrSnapshotConfiguration},
		{"53300", live, &pgconn.PgError{Code: "53300"}, outbound.ErrSnapshotConfiguration},
		{"Verbindungsfehler", live, &pgconn.ConnectError{Config: &pgconn.Config{}}, outbound.ErrSnapshotTransient},
		{"anderer SQLSTATE", live, &pgconn.PgError{Code: "42P01"}, fallback},
		{"Fehler ohne SQLSTATE", live, errors.New("boom"), fallback},
	}
	for _, tc := range cases {
		err := snapshotlogic.Classify(tc.ctx, tc.cause, fallback, "Phase")
		if !errors.Is(err, tc.want) {
			t.Errorf("%s: %v, erwartet %v", tc.name, err, tc.want)
		}
		if !strings.Contains(err.Error(), "Phase") {
			t.Errorf("%s: die Meldung nennt die Phase nicht: %v", tc.name, err)
		}
	}
}
