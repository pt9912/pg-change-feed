package model

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
