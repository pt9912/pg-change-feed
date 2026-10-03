package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrHeartbeatStorage trägt die Fehlerklasse `storage` dieses Ports
// (`ADR-0023`): ein Persistenzfehler am Heartbeat-Speicher
// bleibt über `errors.Is` klassifizierbar, ohne Treibertyp. Abgrenzung:
// anders als der ChangeStore-Sentinel (`ErrStorage`, Klasse-Aktion „kein
// Source-ACK") trägt dieser Sentinel keine
// Klasse-Aktion am Capture-Pfad — ein Heartbeat-Schreibfehler bricht den
// Stream-Lauf nicht ab.
var ErrHeartbeatStorage = messagecode.New(messagecode.HeartbeatStoreFailed, "Persistenzfehler im Heartbeat-Speicher")

// HeartbeatPort trägt die Lebenszeichen-Fähigkeit des Capture-Prozesses
// (`ARC-004`): der periodische Schreib-Zug der Composition Root
// meldet, dass die Instanz läuft. Die Quelle trägt zugleich
// die Prozess-Kennung — im MVP-Schnitt trägt eine Instanz die Quelle und
// den CDC-Speicher gleichermaßen (Abschnitt 1 Lastenheft), ein eigenes
// Prozess-Kennungsfeld trägt dieser Port nicht. Derselbe Port trägt auch
// den zuletzt beobachteten Fehlerzustand — dieselbe Ablage, statt eine
// zweite Tabelle
// einzuführen.
type HeartbeatPort interface {
	// Beat trägt das Lebenszeichen der Quelle fort; die Instanzzeit der
	// Speicherseite trägt den Zeitstempel (wie `transaction.committed_at`)
	// — der Aufrufer übergibt keine Uhr. Ein erfolgreicher Beat löscht
	// einen zuvor gemeldeten Fehlerzustand (`Fault`) wieder — der
	// Fehlerzustand endet dadurch selbst erkennbar (`LH-FA-ADM-003`
	// Boundary).
	Beat(ctx context.Context, source model.SourceID) error

	// Fault trägt den zuletzt beobachteten Fehlerzustand der Quelle fort
	// (`LH-FA-ADM-003`): die Composition Root ruft ihn
	// auf, bevor der Capture-Prozess auf einen Adapter-Fehler endet.
	// Der Zeitstempel läuft mit fort (wie `Beat`) — ein
	// Fehlerzustand ist damit ebenso ein Lebenszeichen, nur mit Klasse und
	// Meldungscode. Der Code ist ein Code dieser Klasse (`messagecode.ClassOf`);
	// ein Widerspruch erreicht keinen SQL-Aufruf.
	Fault(ctx context.Context, source model.SourceID, class model.ErrorClass, code messagecode.Code) error
}
