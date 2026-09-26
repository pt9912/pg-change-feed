package outbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrChangeStream markiert einen fehlgeschlagenen Stream-Publish-Aufruf
// (`ADR-0060`) für `errors.Is`. Der Kanal ist ein In-Prozess-Fan-out ohne
// externes System: ein Fehler entsteht hier aus einem ungültigen Aufruf, nicht
// aus einer vorübergehend nicht verfügbaren Quelle oder einem Speicher.
// Anders als `ErrNotify` (`changenotification.go`, Klasse `transient`) ordnet
// er sich deshalb keiner Fehlerklasse zu. Der Aufrufer (`CaptureService`)
// fängt ihn an der Aufrufstelle ab und lässt ihn nicht in den Rückgabewert des
// Capture-Aufrufs eingehen.
var ErrChangeStream = stderrors.New("Stream-Publish fehlgeschlagen")

// ChangeStreamPort trägt jeden einzelnen Change einer committed
// Quelltransaktion an die verbundenen Live-Stream-Consumer (`ADR-0060`).
// Anders als `ChangeNotificationPort`, das ein tabellen-granulares
// Wecksignal ohne Inhalt trägt, überträgt dieser Port den vollständigen
// Change-Inhalt je Zeile, ohne Deduplizierung nach Tabelle. Der Port ist
// optional: ohne konfigurierten Adapter bleibt die Fähigkeit deaktiviert und
// das bestehende Verhalten unverändert. Die erste Driven-Implementierung ist
// der `Broadcaster` in `internal/adapters/driven/grpcstream`.
type ChangeStreamPort interface {
	// Publish verteilt einen Change an alle zum Aufrufzeitpunkt
	// registrierten Stream-Abonnenten (`ADR-0060`). Der Aufruf blockiert nie
	// auf einen Abonnenten: jeder Abonnent trägt eine begrenzte
	// Empfangs-Warteschlange; liest er nicht schnell genug, werden die über
	// sie hinausgehenden Changes für ihn verworfen (Drop-Newest, nicht
	// nachgeliefert). Der Aufruf trägt keine Zustellgarantie: ein Abonnent,
	// der nicht verbunden ist oder langsamer liest als Changes eintreffen,
	// verpasst die betroffenen Changes ersatzlos — die Nachvollziehbarkeit
	// tragen ausschließlich der Lesezugriffsweg (`ChangeStorePort`, View
	// `cdc.changes`) und die bestätigte Consumer-Position. Ein Fehler
	// entsteht nur aus einem ungültigen Aufruf (`ErrChangeStream`) oder einem
	// bereits beendeten `ctx`; er wird an der Aufrufstelle abgefangen und darf
	// die bereits erfolgte Persistierung oder Bestätigung nicht beeinflussen.
	Publish(ctx context.Context, change *model.Change) error
}
