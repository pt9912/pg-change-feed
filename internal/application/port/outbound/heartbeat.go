package outbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrHeartbeatStorage trägt die Fehlerklasse `storage` dieses Ports
// (`SPEC-008`, `ADR-0023`): ein Persistenzfehler am Heartbeat-Speicher
// bleibt über `errors.Is` klassifizierbar, ohne Treibertyp. Abgrenzung:
// anders als der ChangeStore-Sentinel (`ErrStorage`, Klasse-Aktion „kein
// Source-ACK", `LH-QA-REL-001.a`) trägt dieser Sentinel keine
// Klasse-Aktion am Capture-Pfad — ein Heartbeat-Schreibfehler bricht den
// Stream-Lauf nicht ab (slice-012 §1: Liveness, keine Fehlerdetails;
// Folge-Slice slice-013 trägt Erfassungs-Fehlerklassen).
var ErrHeartbeatStorage = stderrors.New("Fehlerklasse storage: Persistenzfehler im Heartbeat-Speicher")

// HeartbeatPort trägt die Lebenszeichen-Fähigkeit des Capture-Prozesses
// (`ARC-004`, `ADR-0024` „Betriebsinformationen … müssen aus allen
// Schichten ankommen"): der periodische Schreib-Zug der Composition Root
// (`ADR-0026`) meldet, dass die Instanz läuft. Die Quelle trägt zugleich
// die Prozess-Kennung — im MVP-Schnitt trägt eine Instanz die Quelle und
// den CDC-Speicher gleichermaßen (Abschnitt 1 Lastenheft), ein eigenes
// Prozess-Kennungsfeld trägt dieser Port nicht.
type HeartbeatPort interface {
	// Beat trägt das Lebenszeichen der Quelle fort; die Instanzzeit der
	// Speicherseite trägt den Zeitstempel (wie `transaction.committed_at`)
	// — der Aufrufer übergibt keine Uhr.
	Beat(ctx context.Context, source model.SourceID) error
}
