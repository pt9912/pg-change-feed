package outbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrChangeStream trägt die Fehlerklasse eines fehlgeschlagenen
// Stream-Publish-Aufrufs. Der Kanal ist ein In-Prozess-Fan-out ohne
// externes System (`ADR-0060` Teilfrage 5): ein Fehler entsteht hier aus
// einem ungültigen Aufruf, nicht aus einer vorübergehend nicht
// verfügbaren Quelle oder einem Speicher — anders als `ErrNotify` ist er
// deshalb keiner der sieben `SPEC-008`-Klassen zugeordnet. Der Aufrufer
// (`CaptureService`, `ADR-0060` Teilfrage 2) fängt ihn an der Aufrufstelle
// ab und lässt ihn nicht in den Rückgabewert des Capture-Aufrufs eingehen.
var ErrChangeStream = stderrors.New("Stream-Publish fehlgeschlagen")

// ChangeStreamPort trägt jeden einzelnen Change einer committed
// Quelltransaktion an die verbundenen Live-Stream-Consumer (`ADR-0060`
// Teilfrage 2): anders als `ChangeNotificationPort`, das ein
// tabellen-granulares Wecksignal ohne Inhalt trägt (`SPEC-017`,
// `ADR-0055`), überträgt dieser Port den vollständigen Change-Inhalt je
// Zeile (`LH-FA-SST-008`) — keine Deduplizierung nach Tabelle. Der Port
// ist optional: ohne konfigurierten Adapter bleibt die Fähigkeit
// deaktiviert, bestehendes Verhalten unverändert (`ADR-0060` Teilfrage 6).
// Die erste Driven-Implementierung ist der `Broadcaster` in
// `internal/adapters/driven/grpcstream`.
type ChangeStreamPort interface {
	// Publish verteilt einen Change an alle zum Aufrufzeitpunkt
	// registrierten Stream-Abonnenten (`ADR-0060` Teilfrage 3, `ADR-0066`).
	// Der Aufruf blockiert nie auf einen Abonnenten: jeder Abonnent trägt
	// eine begrenzte Empfangs-Warteschlange; liest er nicht schnell genug,
	// werden die über sie hinausgehenden Changes für ihn verworfen
	// (Drop-Newest, nicht nachgeliefert). Der Aufruf trägt keine
	// Zustellgarantie: ein Abonnent, der nicht verbunden ist oder langsamer
	// liest als Changes eintreffen, verpasst die betroffenen Changes
	// ersatzlos — die Nachvollziehbarkeit bleibt ausschließlich beim
	// Lesezugriffsweg (`LH-FA-REA-001` ff.) und der bestätigten
	// Consumer-Position (`LH-FA-CON-003`/`005`). Ein Fehler entsteht nur aus
	// einem ungültigen Aufruf (`ErrChangeStream`) oder einem bereits
	// beendeten `ctx`; er wird an der Aufrufstelle abgefangen und darf die
	// bereits erfolgte Persistierung oder Bestätigung nicht beeinflussen.
	Publish(ctx context.Context, change *model.Change) error
}
