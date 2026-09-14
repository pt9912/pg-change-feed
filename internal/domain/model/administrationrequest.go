package model

import (
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// AdministrationRequestID identifiziert einen Antrags-Datensatz der
// schreibenden SQL-Administration (`cdc.administration_request`,
// `LH-FA-ADM-001`): `cdc.enable_table`/`cdc.disable_table` vergeben sie
// beim Schreiben, die Administrations-Goroutine trägt sie unverändert
// durch Verarbeitung und Ergebnis-Vermerk.
type AdministrationRequestID string

// AdministrationRequestKind trägt die geschlossene Menge der Antragsarten
// (`chk_administration_request_kind`, Tabelle `cdc.administration_request`):
// die beiden Tabellen-Antragsarten `enable`/`disable` und die beiden
// Spalten-Antragsarten `exclude_column`/`include_column`
// (`LH-FA-CFG-005`).
type AdministrationRequestKind string

const (
	AdministrationRequestEnable        AdministrationRequestKind = "enable"
	AdministrationRequestDisable       AdministrationRequestKind = "disable"
	AdministrationRequestExcludeColumn AdministrationRequestKind = "exclude_column"
	AdministrationRequestIncludeColumn AdministrationRequestKind = "include_column"
)

// AdministrationRequest trägt einen offenen (`pending`) Antrags-Datensatz,
// wie ihn die Administrations-Goroutine liest und verarbeitet: Quelle,
// Schema und Tabellenname adressieren dieselbe Tabelle wie
// `EnableTableCommand`/`DisableTableCommand`, ohne deren
// Bindungs-Kennungen — die vergibt die Verarbeitung selbst (`ARC-007`).
// `Column` trägt den Ziel-Spaltennamen der beiden Spalten-Antragsarten; die
// beiden Tabellen-Antragsarten tragen dort den leeren Wert.
type AdministrationRequest struct {
	ID     AdministrationRequestID
	Source SourceID
	Schema string
	Table  string
	Column string
	Kind   AdministrationRequestKind
}

// NewAdministrationRequest legt einen Antrags-Datensatz an und verlangt
// nichtleere Kennungen (ID, Quelle, Schema, Tabelle) sowie eine Antragsart
// aus der geschlossenen Menge `enable`/`disable`/`exclude_column`/
// `include_column` — dasselbe Konstruktor-Muster wie die übrigen zehn
// Domänentypen in diesem Paket (z. B. `NewSchemaVersion`); die Prüfung der
// geschlossenen Menge liegt am Domain-Core-Rand, wie es die
// Architektur-Sicht für Domänenobjekte und ihre Invarianten vorsieht
// (`ARC-001`).
// Die beiden Spalten-Antragsarten tragen eine nichtleere Spalte — ohne sie
// adressiert der Antrag kein Ziel; die beiden Tabellen-Antragsarten tragen
// keine Spalte.
func NewAdministrationRequest(id AdministrationRequestID, source SourceID, schema, table, column string, kind AdministrationRequestKind) (AdministrationRequest, error) {
	if id == "" || source == "" || schema == "" || table == "" {
		return AdministrationRequest{}, domainerrors.ErrEmptyIdentifier
	}
	switch kind {
	case AdministrationRequestEnable, AdministrationRequestDisable:
	case AdministrationRequestExcludeColumn, AdministrationRequestIncludeColumn:
		if column == "" {
			return AdministrationRequest{}, domainerrors.ErrEmptyIdentifier
		}
	default:
		return AdministrationRequest{}, domainerrors.ErrInvalidAdministrationRequestKind
	}
	return AdministrationRequest{ID: id, Source: source, Schema: schema, Table: table, Column: column, Kind: kind}, nil
}
