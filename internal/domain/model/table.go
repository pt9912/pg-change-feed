package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// SourceTableID identifiziert eine aktivierte Tabelle einer Quelle
// (`SPEC-001`, Tabelle `cdc.source_table`).
type SourceTableID string

// SourceTable ist eine aktivierte Tabelle; Schema und Tabellenname
// identifizieren sie innerhalb der Quelle (`LH-FA-DAT-002`).
type SourceTable struct {
	ID       SourceTableID
	SourceID SourceID
	Schema   string
	Table    string
}

// NewSourceTable legt eine aktivierte Tabelle an und verlangt nichtleere
// Kennungen; die Zuordnung zur Quelle ist Teil des Objekts (`LH-FA-DAT-002`,
// Boundary: gleichnamige Tabellen in verschiedenen Schemata sind
// unterscheidbar).
func NewSourceTable(id SourceTableID, source SourceID, schema, table string) (SourceTable, error) {
	if id == "" || source == "" || schema == "" || table == "" {
		return SourceTable{}, domainerrors.ErrEmptyIdentifier
	}
	return SourceTable{ID: id, SourceID: source, Schema: schema, Table: table}, nil
}

// QualifiedName trägt schema-qualifizierten Tabellennamen in der Form
// `schema.table`.
func (t SourceTable) QualifiedName() string {
	return t.Schema + "." + t.Table
}
