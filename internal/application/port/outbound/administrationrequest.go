package outbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrAdministrationStorage trägt die Fehlerklasse `storage` dieses Ports
// (`SPEC-008`, `ADR-0023`): ein Persistenzfehler an der Antrags-Queue
// bleibt über `errors.Is` klassifizierbar, ohne Treibertyp. Wie
// `ErrHeartbeatStorage` trägt dieser Sentinel keine Klasse-Aktion am
// Capture-Pfad — ein Lese- oder Vermerk-Fehler der Administrations-
// Goroutine bricht den Stream-Lauf nicht ab (best-effort, derselbe Fallback-
// Poll deckt einen verpassten Durchlauf ab).
var ErrAdministrationStorage = stderrors.New("Fehlerklasse storage: Persistenzfehler an der Antrags-Queue")

// AdministrationRequestPort trägt die Lese- und Ergebnis-Fähigkeit der
// Antrags-Queue (`ARC-004`, `LH-FA-ADM-001`): `cdc.enable_table`/
// `cdc.disable_table` schreiben den Antrags-Datensatz direkt über SQL
// (kein Go-Aufrufpfad, Fähigkeits-Trennung); dieser Port trägt ausschließlich
// die Gegenrichtung — die Administrations-Goroutine liest offene Anträge
// und vermerkt ihr Ergebnis.
type AdministrationRequestPort interface {
	// ListPending liest die Anträge mit Status `pending` in Anlage-
	// Reihenfolge — sowohl nach `NOTIFY`-Wecksignal als auch periodisch als
	// Fallback für einen verpassten Wecksignal (Verbindungsabbruch der
	// `LISTEN`-Verbindung).
	ListPending(ctx context.Context) ([]model.AdministrationRequest, error)

	// MarkApplied vermerkt einen erfolgreich verarbeiteten Antrag; ein
	// bereits vermerkter Antrag bleibt unverändert (Idempotenz — ein
	// erneuter Aufruf über den Fallback-Poll trägt keinen zweiten Vermerk).
	MarkApplied(ctx context.Context, id model.AdministrationRequestID) error

	// MarkFailed vermerkt einen gescheiterten Antrag samt Fehlertext; der
	// Antrag bleibt damit sichtbar unterscheidbar von einem noch offenen
	// (`pending`) oder erfolgreichen (`applied`) Antrag.
	MarkFailed(ctx context.Context, id model.AdministrationRequestID, message string) error
}
