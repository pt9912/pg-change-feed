// Die Verwaltungs-Use-Cases trägt diese Datei an einer Stelle (`ARC-003`):
// die Driving-Adapter (`ARC-005`) rufen die Tabellen-Aktivierung, die
// Deaktivierung, den Status und die Liste über sie auf (`ADR-0028`); die
// Orchestrierung liegt in den Application Services (`ARC-002`).

package inbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrSourceTableMissing meldet, dass die angesprochene Tabelle an der
// Quelle nicht existiert; die Negative-Pfade von Aktivierung,
// Deaktivierung und Status-Abfrage enden über dieses Sentinel sichtbar
// (`LH-FA-CFG-001`/`002`/`003`) statt still.
var ErrSourceTableMissing = stderrors.New("Tabelle existiert nicht an der Quelle")

// EnableTableCommand trägt die Eingabe der Aktivierung (`LH-FA-CFG-001`):
// die Tabelle mit ihrer Bindung — Tabellen- und Schema-Version-Kennung
// (`SPEC-001`, `cdc.source_table`/`cdc.schema_version`) — und die
// Publication, die die Tabelle aufnimmt (`LH-FA-CFG-001.a`). Die
// Anfangs-Version trägt das erste Schema der Tabelle; spätere Versionen
// trägt die Schema-Evolution (`LH-FA-SCH-004.a`, Metadata-Pfad).
type EnableTableCommand struct {
	Source          model.SourceID
	Schema          string
	Table           string
	TableID         model.SourceTableID
	SchemaVersionID model.SchemaVersionID
	Version         int64
	Publication     string
}

// EnableTableResult trägt den Ausgang der Aktivierung: die aktivierte
// Tabelle und den idempotenten Ausgang (`LH-FA-CFG-001` Boundary) — ein
// erneuter Aufruf einer aktiven Tabelle meldet das Ergebnis, ohne den
// Stand zu ändern.
type EnableTableResult struct {
	Table          model.SourceTable
	AlreadyEnabled bool
}

// DisableTableCommand trägt die Eingabe der Deaktivierung
// (`LH-FA-CFG-002`): die Tabelle und die Publication, aus der sie
// entzogen wird.
type DisableTableCommand struct {
	Source      model.SourceID
	Schema      string
	Table       string
	Publication string
}

// DisableTableResult trägt den Ausgang der Deaktivierung: der Stopp der
// Erfassung läuft über den Publication-Entzug; der Bindungs-Zeilen-Ausgang
// trägt den Change-Bestand — bei persistierten Changes bleibt die Zeile als
// Herkunft bestehen (`LH-FA-CFG-002` Out-of-Scope: ihr Verhalten folgt der
// Retention), ohne Bestand wird sie entfernt.
type DisableTableResult struct {
	Removed  bool
	Retained bool
}

// GetStatusQuery trägt die Eingabe der Status-Abfrage (`LH-FA-CFG-003`).
type GetStatusQuery struct {
	Source model.SourceID
	Schema string
	Table  string
}

// GetStatusResult trägt den Ausgang der Status-Abfrage: der Zustand
// „aktiviert" einer Tabelle; „nicht aktiviert" meldet die Rückkehr mit
// Enabled=false (`LH-FA-CFG-003` Happy Path und Boundary).
type GetStatusResult struct {
	Enabled bool
}

// ListTablesQuery trägt die Eingabe der Tabellen-Liste (`LH-FA-CFG-004`).
type ListTablesQuery struct {
	Source model.SourceID
}

// ListTablesResult trägt die aktivierten Tabellen der Quelle
// (`LH-FA-CFG-004`); ohne Aktivierung trägt die Liste keine Tabellen.
type ListTablesResult struct {
	Tables []model.SourceTable
}

// EnableTableUseCase aktiviert CDC für eine einzelne Tabelle
// (`LH-FA-CFG-001`, `ADR-0028`): die Bindungs-Zeilen und die Publication
// der Quelle tragen die Aktivierung.
type EnableTableUseCase interface {
	Enable(ctx context.Context, command EnableTableCommand) (EnableTableResult, error)
}

// DisableTableUseCase deaktiviert CDC für eine einzelne Tabelle
// (`LH-FA-CFG-002`, `ADR-0028`); das Verhalten der persistierten Changes
// trägt die Retention, nicht die Deaktivierung.
type DisableTableUseCase interface {
	Disable(ctx context.Context, command DisableTableCommand) (DisableTableResult, error)
}

// GetStatusUseCase meldet den CDC-Zustand einer Tabelle
// (`LH-FA-CFG-003`, `ADR-0028`).
type GetStatusUseCase interface {
	Status(ctx context.Context, query GetStatusQuery) (GetStatusResult, error)
}

// ListTablesUseCase listet die aktivierten Tabellen einer Quelle
// (`LH-FA-CFG-004`, `ADR-0028`).
type ListTablesUseCase interface {
	ListTables(ctx context.Context, query ListTablesQuery) (ListTablesResult, error)
}
