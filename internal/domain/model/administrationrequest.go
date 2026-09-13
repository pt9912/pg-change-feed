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
// (`chk_administration_request_kind`, Tabelle `cdc.administration_request`).
type AdministrationRequestKind string

const (
	AdministrationRequestEnable  AdministrationRequestKind = "enable"
	AdministrationRequestDisable AdministrationRequestKind = "disable"
)

// AdministrationRequest trägt einen offenen (`pending`) Antrags-Datensatz,
// wie ihn die Administrations-Goroutine liest und verarbeitet: Quelle,
// Schema und Tabellenname adressieren dieselbe Tabelle wie
// `EnableTableCommand`/`DisableTableCommand`, ohne deren
// Bindungs-Kennungen — die vergibt die Verarbeitung selbst (`ARC-007`).
type AdministrationRequest struct {
	ID     AdministrationRequestID
	Source SourceID
	Schema string
	Table  string
	Kind   AdministrationRequestKind
}

// NewAdministrationRequest legt einen Antrags-Datensatz an und verlangt
// nichtleere Kennungen (ID, Quelle, Schema, Tabelle) sowie eine Antragsart
// aus der geschlossenen Menge `enable`/`disable` — dasselbe
// Konstruktor-Muster wie die übrigen zehn Domänentypen in diesem Paket
// (z. B. `NewSchemaVersion`), bislang die einzige Ausnahme davon
// (Review-Finding F-4, `review-slice-037.md`). Die Prüfung der
// geschlossenen Menge lag zuvor außerhalb der Domänenschicht (`default`-Zweig
// von `applyAdministrationRequest`, `internal/bootstrap/wiring.go`); dieser
// Konstruktor trägt sie jetzt am Domain-Core-Rand, wie es die
// Architektur-Sicht für Domänenobjekte und ihre Invarianten vorsieht.
func NewAdministrationRequest(id AdministrationRequestID, source SourceID, schema, table string, kind AdministrationRequestKind) (AdministrationRequest, error) {
	if id == "" || source == "" || schema == "" || table == "" {
		return AdministrationRequest{}, domainerrors.ErrEmptyIdentifier
	}
	switch kind {
	case AdministrationRequestEnable, AdministrationRequestDisable:
	default:
		return AdministrationRequest{}, domainerrors.ErrInvalidAdministrationRequestKind
	}
	return AdministrationRequest{ID: id, Source: source, Schema: schema, Table: table, Kind: kind}, nil
}
