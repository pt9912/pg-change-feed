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
// die beiden Tabellen-Antragsarten `enable`/`disable`, die beiden
// Spalten-Antragsarten `exclude_column`/`include_column`
// (`LH-FA-CFG-005`), die Bestands-Antragsart `backfill`
// (`LH-FA-CAP-009`, `ADR-0111`) und die beiden Transformations-Antragsarten
// `set_transformation`/`remove_transformation` (`LH-FA-CFG-007`,
// `ADR-0112`). Bei `backfill` heißt der Status `applied` „angenommen": die
// Ausführung steht in `cdc.backfill_run` (`SPEC-019`).
type AdministrationRequestKind string

const (
	AdministrationRequestEnable               AdministrationRequestKind = "enable"
	AdministrationRequestDisable              AdministrationRequestKind = "disable"
	AdministrationRequestExcludeColumn        AdministrationRequestKind = "exclude_column"
	AdministrationRequestIncludeColumn        AdministrationRequestKind = "include_column"
	AdministrationRequestBackfill             AdministrationRequestKind = "backfill"
	AdministrationRequestSetTransformation    AdministrationRequestKind = "set_transformation"
	AdministrationRequestRemoveTransformation AdministrationRequestKind = "remove_transformation"
)

// AdministrationRequest trägt einen offenen (`pending`) Antrags-Datensatz,
// wie ihn die Administrations-Goroutine liest und verarbeitet: Quelle,
// Schema und Tabellenname adressieren dieselbe Tabelle wie
// `EnableTableCommand`/`DisableTableCommand`, ohne deren
// Bindungs-Kennungen — die vergibt die Verarbeitung selbst (`ARC-007`).
// `Column` trägt den Ziel-Spaltennamen der beiden Spalten-Antragsarten; die
// übrigen Antragsarten tragen dort den leeren Wert. `RuleName` trägt den
// Regelnamen der beiden Transformations-Antragsarten, `RuleSpec` die
// Regelform der Antragsart `set_transformation` als JSON-Text; die übrigen
// Antragsarten tragen beide als leeren Wert. Die Regelform prüft der Use
// Case, nicht dieser Typ (`ADR-0046`).
type AdministrationRequest struct {
	ID       AdministrationRequestID
	Source   SourceID
	Schema   string
	Table    string
	Column   string
	RuleName string
	RuleSpec string
	Kind     AdministrationRequestKind
}

// NewAdministrationRequest legt einen Antrags-Datensatz an und verlangt
// nichtleere Kennungen (ID, Quelle, Schema, Tabelle) sowie eine Antragsart
// aus der geschlossenen Menge `enable`/`disable`/`exclude_column`/
// `include_column`/`backfill`/`set_transformation`/`remove_transformation` —
// dasselbe Konstruktor-Muster wie die übrigen zehn Domänentypen in diesem
// Paket (z. B. `NewSchemaVersion`); die Prüfung der geschlossenen Menge liegt
// am Domain-Core-Rand, wie es die Architektur-Sicht für Domänenobjekte und
// ihre Invarianten vorsieht (`ARC-001`).
// Die beiden Spalten-Antragsarten tragen eine nichtleere Spalte — ohne sie
// adressiert der Antrag kein Ziel; die beiden Transformations-Antragsarten
// tragen einen nichtleeren Regelnamen, `set_transformation` zusätzlich eine
// nichtleere Regelform.
func NewAdministrationRequest(id AdministrationRequestID, source SourceID, schema, table, column, ruleName, ruleSpec string, kind AdministrationRequestKind) (AdministrationRequest, error) {
	if id == "" || source == "" || schema == "" || table == "" {
		return AdministrationRequest{}, domainerrors.ErrEmptyIdentifier
	}
	switch kind {
	case AdministrationRequestEnable, AdministrationRequestDisable, AdministrationRequestBackfill:
	case AdministrationRequestExcludeColumn, AdministrationRequestIncludeColumn:
		if column == "" {
			return AdministrationRequest{}, domainerrors.ErrEmptyIdentifier
		}
	case AdministrationRequestSetTransformation:
		if ruleName == "" || ruleSpec == "" {
			return AdministrationRequest{}, domainerrors.ErrEmptyIdentifier
		}
	case AdministrationRequestRemoveTransformation:
		if ruleName == "" {
			return AdministrationRequest{}, domainerrors.ErrEmptyIdentifier
		}
	default:
		return AdministrationRequest{}, domainerrors.ErrInvalidAdministrationRequestKind
	}
	return AdministrationRequest{ID: id, Source: source, Schema: schema, Table: table, Column: column, RuleName: ruleName, RuleSpec: ruleSpec, Kind: kind}, nil
}
