package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// ColumnOID trägt die rohe PostgreSQL-Typ-OID einer Spalte (`SPEC-004`,
// `ADR-0015` Folgepflicht); die technologieunabhängige Übersetzung dieser
// OID in eine Vergleichs-/Kompatibilitätsentscheidung trägt dieses Modell
// nicht.
type ColumnOID uint32

// Column trägt eine Spalte eines TableSchema: Name und PostgreSQL-Typ-OID
// (`ADR-0015`, Option C). Die Spalten-Reihenfolge innerhalb eines
// TableSchema trägt die Anlage-Reihenfolge des Slice — kein separates
// Ordnungsfeld auf diesem Wertobjekt.
type Column struct {
	Name string
	OID  ColumnOID
}

// TableSchema trägt die Spaltenmenge einer Schema-Version (`ADR-0015`,
// `SPEC-004`, `LH-FA-SCH-005`): jeder Change referenziert eine
// SchemaVersionID (siehe `SchemaVersion`, `schema_version.go`); TableSchema
// trägt die dazugehörige Spaltenform, über die die historisch stabile
// Interpretation eines Change läuft — ohne sie bleibt eine inkompatible
// Typänderung unerkennbar (`LH-FA-SCH-004`).
type TableSchema struct {
	VersionID SchemaVersionID
	Columns   []Column
}

// NewTableSchema legt ein TableSchema an und verlangt eine nichtleere
// Versions-Kennung, mindestens eine Spalte und nichtleere Spaltennamen je
// Spalte (`ErrEmptyColumns`, `ErrEmptyIdentifier`).
func NewTableSchema(versionID SchemaVersionID, columns []Column) (TableSchema, error) {
	if versionID == "" {
		return TableSchema{}, domainerrors.ErrEmptyIdentifier
	}
	if len(columns) == 0 {
		return TableSchema{}, domainerrors.ErrEmptyColumns
	}
	copied := make([]Column, len(columns))
	for i, column := range columns {
		if column.Name == "" {
			return TableSchema{}, domainerrors.ErrEmptyIdentifier
		}
		copied[i] = column
	}
	return TableSchema{VersionID: versionID, Columns: copied}, nil
}
