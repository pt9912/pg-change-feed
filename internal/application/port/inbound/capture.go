// Package inbound bündelt die Inbound-Ports der Use Cases (`ARC-003`):
// die Driving-Adapter (`ARC-005`) rufen die fachlichen Fähigkeiten über
// sie auf (`ADR-0028`); die Orchestrierung liegt in den Application
// Services (`ARC-002`).
package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// CaptureCommand trägt die Eingabe des Capture Use Cases: die
// vollständige, committed Quelltransaktion (`LH-FA-CAP-006.a`). Offene
// Transaktionen sind nicht konsumierbar und zurückgerollte erzeugen keine
// committed Changes (`LH-FA-CAP-006`, `LH-FA-CAP-007`); die Ordnungs-Logik
// lehnt sie über die Transaktion ab.
type CaptureCommand struct {
	Transaction *model.ChangeTransaction
}

// CaptureResult trägt den Ausgang des Capture Use Cases: die bestätigte
// Quellposition. Ein Persistenzfehler liefert kein Ergebnis und keinen
// Source-ACK (`LH-QA-REL-001.a`).
type CaptureResult struct {
	Acknowledged model.SourcePosition
}

// CaptureInboundPort nimmt eine vollständige, committed Quelltransaktion
// auf (`ADR-0028`): der Replication-Stream-Adapter übersetzt seine
// dekodierten `pgoutput`-Ereignisse in seine Aufrufe (`ARC-003`,
// `ADR-0027`).
//
// Ordnung (`LH-QA-REL-001.a`): Receive → Decode → Persist → COMMIT Store
// → ACK Source. Die Implementierung bestätigt die Commit-Position erst,
// wenn die abhängigen Changes dauerhaft gespeichert sind; ein
// Persistenzfehler endet ohne Source-ACK (`SPEC-008`, Klasse `storage`).
//
// Eine committed Quelltransaktion darf ohne Changes auftreten — in
// `pgoutput` treten BEGIN/COMMIT-Paare ohne relevante Relation-Nachricht
// als realer Input auf (`LH-FA-CAP-006.a`). Der Pfad persistiert sie mit
// ihrer Commit-Position und bestätigt die Position; das Ablehnen würde
// dieselbe leere Transaktion in einer Wiederholung (`ADR-0012`) endlos
// neu liefern.
type CaptureInboundPort interface {
	Capture(ctx context.Context, command CaptureCommand) (CaptureResult, error)
}
