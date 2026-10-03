package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// ErrNotify trägt einen Code der Klasse `internal` des Wecksignal-Pfads
// (`ADR-0023`; die Composition Root ordnet ihn keiner anderen Klasse zu): eine
// unterbrochene NATS-Verbindung
// oder ein fehlgeschlagener Publish-Aufruf ist strukturell dieselbe
// Kategorie wie eine vorübergehend nicht verfügbare Quelle/Speicher — der
// nats.go-Client trägt das für `transient` vorgesehene
// Erneut-versuchen-mit-Backoff bereits über sein eingebautes Reconnect auf
// Verbindungsebene. Application und Betrieb klassifizieren über
// `errors.Is`, ohne einen Treibertyp zu kennen; die technische Ursache
// bleibt über die zweite Wrappung lesbar.
var ErrNotify = messagecode.New(messagecode.NotifyFailed, "Wecksignal fehlgeschlagen")

// ChangeNotificationPort trägt das Wecksignal an verbundene Consumer
// (`ARC-013`): ein reines, verlustbehaftetes
// Trigger-Signal ohne Change-Inhalt und ohne Positionsangabe
// — die Nachvollziehbarkeit bleibt ausschließlich beim bestehenden
// Lesezugriffsweg (`ChangeStorePort`). Der Port ist optional: ohne
// konfigurierte Verbindung bleibt die Fähigkeit deaktiviert.
// Die erste Driven-Implementierung ist der
// `NatsChangeNotificationAdapter`.
type ChangeNotificationPort interface {
	// Notify sendet das Wecksignal für die genannte Quelle und Tabelle;
	// Schema und Tabelle gehen getrennt ein, nicht vorkombiniert
	// (`ADR-0056`) — der Adapter, nicht der Aufrufer, setzt sie zu einem
	// eindeutigen Subjekt zusammen. Die Rückkehr ohne Fehler meldet den
	// abgeschickten Publish-Versuch, keine Zustellgarantie (Core NATS
	// Fire-and-Forget). Der Aufruf reiht sich als dritter,
	// optionaler Schritt NACH `ACK Source` ein: sein Fehler
	// wird an der Aufrufstelle abgefangen und darf die bereits erfolgte
	// Persistierung oder Bestätigung nicht beeinflussen.
	Notify(ctx context.Context, sourceID, schema, table string) error
}
