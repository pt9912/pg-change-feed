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

// PendingAdministrationRequest ist eine gelesene Zeile der Antrags-Queue mit
// Status `pending`: entweder der Antrag (`Rejected == nil`) oder, wenn der
// Antrags-Konstruktor die Zeile verwirft, ihre Kennung und der Fehlertext
// (`Request` trägt dann den Nullwert). Die Lesung lehnt keine einzelne Zeile
// ab; die Verarbeitung vermerkt eine verworfene Zeile `failed` (`SPEC-019`).
type PendingAdministrationRequest struct {
	Request  model.AdministrationRequest
	Rejected *RejectedAdministrationRequest
}

// RejectedAdministrationRequest trägt eine vom Antrags-Konstruktor
// verworfene Zeile: `ID` ist die Antrags-Kennung der Zeile, leer, wenn die
// Zeile keine trägt (sie lässt sich dann nicht vermerken), `Message` der
// Fehlertext des `failed`-Vermerks (`SPEC-019`).
type RejectedAdministrationRequest struct {
	ID      model.AdministrationRequestID
	Message string
}

// AdministrationRequestPort trägt die Lese- und Ergebnis-Fähigkeit der
// Antrags-Queue (`ARC-004`, `LH-FA-ADM-001`): `cdc.enable_table`/
// `cdc.disable_table`/`cdc.exclude_column`/`cdc.include_column`/
// `cdc.backfill_table`/`cdc.set_transformation`/`cdc.remove_transformation`
// schreiben den Antrags-Datensatz direkt über SQL
// (kein Go-Aufrufpfad, Fähigkeits-Trennung); dieser Port trägt ausschließlich
// die Gegenrichtung — die Administrations-Goroutine liest offene Anträge
// und vermerkt ihr Ergebnis.
type AdministrationRequestPort interface {
	// ListPending liest die Anträge mit Status `pending` in Aufruf-
	// Reihenfolge (`requested_at` ist der Aufrufzeitpunkt der schreibenden
	// Funktion, bei gleichem Zeitstempel nach der
	// Antrags-Kennung — dieselbe Ordnung, in der die Ableitung des
	// dauerhaften Standes die `applied`-Zeilen liest) — sowohl nach
	// `NOTIFY`-Wecksignal als auch periodisch als
	// Fallback für einen verpassten Wecksignal (Verbindungsabbruch der
	// `LISTEN`-Verbindung). Ein Fehler betrifft nur die Lesung selbst
	// (Anfrage, Scan, Iteration); eine Zeile, die der Antrags-Konstruktor
	// verwirft, steht als `Rejected` an ihrer Stelle der Ordnung.
	ListPending(ctx context.Context) ([]PendingAdministrationRequest, error)

	// MarkApplied vermerkt einen erfolgreich verarbeiteten Antrag; ein
	// bereits vermerkter Antrag bleibt unverändert (Idempotenz — ein
	// erneuter Aufruf über den Fallback-Poll trägt keinen zweiten Vermerk).
	MarkApplied(ctx context.Context, id model.AdministrationRequestID) error

	// MarkFailed vermerkt einen gescheiterten Antrag samt Fehlertext; der
	// Antrag bleibt damit sichtbar unterscheidbar von einem noch offenen
	// (`pending`) oder erfolgreichen (`applied`) Antrag.
	MarkFailed(ctx context.Context, id model.AdministrationRequestID, message string) error
}
