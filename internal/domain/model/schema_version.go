package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// SchemaVersionID referenziert eine Schema-Version (`SPEC-001`, Tabelle
// `cdc.schema_version`).
type SchemaVersionID string

// SchemaVersion trägt eine Version des Tabellenschemas (`SPEC-004`);
// jeder Change referenziert eine solche Version (`LH-FA-SCH-005`,
// `ADR-0029`, Regel 7).
type SchemaVersion struct {
	ID            SchemaVersionID
	SourceTableID SourceTableID
	Version       int64
}

// NewSchemaVersion legt eine Schema-Version an und verlangt nichtleere
// Kennungen und eine Versionsnummer mindestens 1.
func NewSchemaVersion(id SchemaVersionID, table SourceTableID, version int64) (SchemaVersion, error) {
	if id == "" || table == "" {
		return SchemaVersion{}, domainerrors.ErrEmptyIdentifier
	}
	if version < 1 {
		return SchemaVersion{}, domainerrors.ErrNonPositiveVersion
	}
	return SchemaVersion{ID: id, SourceTableID: table, Version: version}, nil
}
