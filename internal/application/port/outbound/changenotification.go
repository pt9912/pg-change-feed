package outbound

import (
	"context"
	stderrors "errors"
)

// ErrNotify trägt die Fehlerklasse `transient` des Wecksignal-Pfads
// (`SPEC-008`, `ADR-0023`, `ADR-0055`): eine unterbrochene NATS-Verbindung
// oder ein fehlgeschlagener Publish-Aufruf ist strukturell dieselbe
// Kategorie wie eine vorübergehend nicht verfügbare Quelle/Speicher — der
// nats.go-Client trägt das in `SPEC-008` für `transient` vorgesehene
// Erneut-versuchen-mit-Backoff bereits über sein eingebautes Reconnect auf
// Verbindungsebene. Application und Betrieb klassifizieren über
// `errors.Is`, ohne einen Treibertyp zu kennen; die technische Ursache
// bleibt über die zweite Wrappung lesbar.
var ErrNotify = stderrors.New("Fehlerklasse transient: Wecksignal fehlgeschlagen")

// ChangeNotificationPort trägt das Wecksignal an verbundene Consumer
// (`ARC-013`, `ADR-0055`): ein reines, verlustbehaftetes Trigger-Signal
// ohne Change-Inhalt und ohne Positionsangabe (`SPEC-017`) — die
// Nachvollziehbarkeit bleibt ausschließlich beim bestehenden
// Lesezugriffsweg (`ChangeStorePort`). Der Port ist optional: ohne
// konfigurierte Verbindung bleibt die Fähigkeit deaktiviert
// (`LH-FA-SST-007` Boundary). Die erste Driven-Implementierung ist der
// `NatsChangeNotificationAdapter`.
type ChangeNotificationPort interface {
	// Notify sendet das Wecksignal für die genannte Quelle; die Rückkehr
	// ohne Fehler meldet den abgeschickten Publish-Versuch, keine
	// Zustellgarantie (`SPEC-017`, Core NATS Fire-and-Forget). Der Aufruf
	// reiht sich als dritter, optionaler Schritt NACH `ACK Source` ein
	// (`ADR-0055`): sein Fehler wird an der Aufrufstelle abgefangen und
	// darf die bereits erfolgte Persistierung oder Bestätigung nicht
	// beeinflussen.
	Notify(ctx context.Context, sourceID string) error
}
