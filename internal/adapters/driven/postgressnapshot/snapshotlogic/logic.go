// Package snapshotlogic trägt die netzlos prüfbare Logik des
// `PostgresTableSnapshotAdapter` (`ADR-0111` Teilfrage 1): Konstruktions-
// validierung, Slot-Name, Cursor- und Fetch-Anweisung, Bezeichner-Quoting,
// Zeilenwerte, Schätzungs-Abbildung und Fehlerklassifikation. Das Paket
// öffnet keine Verbindung und läuft im Unit-Gegenstand des Coverage-Gates
// (`make test`, `ADR-0071` Punkt 1); die Verbindungsschritte liegen im
// Elternpaket `postgressnapshot`.
package snapshotlogic

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

const (
	// CursorName ist der Name des `NO SCROLL`-Cursors in der Lese-Transaktion.
	CursorName = "cdc_bf_cursor"

	slotPrefix    = "cdc_bf_"
	maxIdentifier = 63
)

// identifierShape ist das Bezeichner-Alphabet der Quelle für Slot-Namen
// (`receive.identifierShape` hält denselben Ausdruck: die Adapter-Schicht
// importiert keine Adapter-Kante).
var identifierShape = regexp.MustCompile(`^[a-z0-9_]+$`)

// snapshotNameShape begrenzt den vom Server gelieferten Snapshot-Namen auf
// das Alphabet, in dem er als Literal in `SET TRANSACTION SNAPSHOT` steht.
var snapshotNameShape = regexp.MustCompile(`^[0-9A-Fa-f-]+$`)

// Validate prüft die Konstruktionsparameter des Adapters, ohne Verbindung,
// und liefert die geparste Verbindungskonfiguration des DSN: ein leerer oder
// nicht parsbarer DSN, eine Blockgröße unter 1 und ein Zeitlimit ≤ 0 enden
// als Fehlerklasse `configuration`.
func Validate(dsn string, blockSize int, slotTimeout time.Duration) (*pgconn.Config, error) {
	if dsn == "" {
		return nil, fmt.Errorf("%w: DSN fehlt", outbound.ErrSnapshotConfiguration)
	}
	config, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: DSN: %v", outbound.ErrSnapshotConfiguration, err)
	}
	if blockSize < 1 {
		return nil, fmt.Errorf("%w: Blockgröße %d unter 1", outbound.ErrSnapshotConfiguration, blockSize)
	}
	if slotTimeout <= 0 {
		return nil, fmt.Errorf("%w: Zeitlimit der Slot-Anlage nicht positiv", outbound.ErrSnapshotConfiguration)
	}
	return config, nil
}

// SlotName bildet den Slot-Namen `cdc_bf_<run>`: die Run-Kennung geht ohne
// Bindestriche ein und trägt danach das Bezeichner-Alphabet.
func SlotName(runID string) (string, error) {
	name := slotPrefix + strings.ReplaceAll(runID, "-", "")
	if len(name) > maxIdentifier || !identifierShape.MatchString(name) || name == slotPrefix {
		return "", fmt.Errorf("%w: Run-Kennung %q ergibt keinen gültigen Slot-Namen", outbound.ErrSnapshotConfiguration, runID)
	}
	return name, nil
}

// ValidSnapshotName meldet, ob der vom Server gelieferte Snapshot-Name das
// Alphabet trägt, in dem er als Literal in `SET TRANSACTION SNAPSHOT` steht.
func ValidSnapshotName(name string) bool {
	return snapshotNameShape.MatchString(name)
}

// QuoteIdent setzt einen Bezeichner in Anführungszeichen; ein
// Anführungszeichen im Namen wird verdoppelt.
func QuoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// CursorStatement bildet die Cursor-Deklaration: jede Spalte steht nur beim
// gequoteten Namen — kein Cast, kein Funktionsaufruf —, sodass der Server
// das Ergebnis der Ausgabefunktion des Spaltentyps im Text-Ergebnisformat
// liefert, dieselbe Erzeugung wie der Text, den `pgoutput` je Spalte sendet
// (`ADR-0115` Festlegung 1). `NO SCROLL` hält keinen Rückwärtspuffer.
func CursorStatement(columns []string, schema, table string) string {
	quoted := make([]string, len(columns))
	for i, column := range columns {
		quoted[i] = QuoteIdent(column)
	}
	return "DECLARE " + CursorName + " NO SCROLL CURSOR FOR SELECT " +
		strings.Join(quoted, ", ") + " FROM " + QuoteIdent(schema) + "." + QuoteIdent(table)
}

// FetchStatement holt die nächsten höchstens `blockSize` Zeilen des Cursors.
func FetchStatement(blockSize int) string {
	return "FETCH FORWARD " + strconv.Itoa(blockSize) + " FROM " + CursorName
}

// Values übernimmt die Rohbytes der Felder einer Zeile als Text: `nil` ist
// NULL, ein leeres Feld ist der leere Text (`ADR-0115` Festlegung 1).
func Values(fields [][]byte) []*string {
	values := make([]*string, len(fields))
	for i, field := range fields {
		if field != nil {
			text := string(field)
			values[i] = &text
		}
	}
	return values
}

// Estimate bildet `pg_class.reltuples` (als Text von `float8`) auf die
// Schätzung ab: ein negativer Wert heißt „unbekannt" (`known` falsch), null
// ist eine bekannte Schätzung (`ADR-0113` Festlegung 3, Punkt 5).
func Estimate(reltuples string) (rows int64, known bool, err error) {
	value, err := strconv.ParseFloat(reltuples, 64)
	if err != nil {
		return 0, false, fmt.Errorf("%w: Zeilenschätzung %q nicht lesbar", outbound.ErrSnapshotStorage, reltuples)
	}
	if value < 0 {
		return 0, false, nil
	}
	return int64(math.Round(value)), true, nil
}

// Classify übersetzt einen Treiber-Fehler an der Adapter-Grenze in die
// Fehlerklasse (`ADR-0023`): ein Kontextende und ein Verbindungsfehler sind
// `transient`, ebenso das Beenden der Sitzung durch den Server (`57P01`,
// `57P02`, `57P03`); `42501` und die Anmelde-Klasse `28…` sind `permission`,
// die Konfigurationsgrenzen der Quelle (`53400`, `53300`) sind
// `configuration`; der Rest trägt die Klasse der Phase (`fallback`).
func Classify(ctx context.Context, cause, fallback error, phase string) error {
	var pgErr *pgconn.PgError
	var connectErr *pgconn.ConnectError
	switch {
	case ctx.Err() != nil || errors.Is(cause, context.DeadlineExceeded) || errors.Is(cause, context.Canceled) || pgconn.Timeout(cause):
		return fmt.Errorf("%w: %s: %v", outbound.ErrSnapshotTransient, phase, cause)
	case errors.As(cause, &pgErr) && (pgErr.Code == "57P01" || pgErr.Code == "57P02" || pgErr.Code == "57P03"):
		return fmt.Errorf("%w: %s: %v", outbound.ErrSnapshotTransient, phase, cause)
	case errors.As(cause, &pgErr) && (pgErr.Code == "42501" || strings.HasPrefix(pgErr.Code, "28")):
		return fmt.Errorf("%w: %s: %v", outbound.ErrSnapshotPermission, phase, cause)
	case errors.As(cause, &pgErr) && (pgErr.Code == "53400" || pgErr.Code == "53300"):
		return fmt.Errorf("%w: %s: %v", outbound.ErrSnapshotConfiguration, phase, cause)
	case errors.As(cause, &connectErr):
		return fmt.Errorf("%w: %s: %v", outbound.ErrSnapshotTransient, phase, cause)
	default:
		return fmt.Errorf("%w: %s: %v", fallback, phase, cause)
	}
}
