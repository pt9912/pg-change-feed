// Der Backfill-Use-Case trägt diese Datei an einer Stelle (`ARC-003`): der
// Administrations-Hintergrundzug nimmt einen Antrag über `Request` an, der
// Backfill-Worker führt einen angenommenen Run über `Execute` aus
// (`ADR-0028`, `ADR-0111` Teilfrage 5, `ADR-0113` Festlegung 1); die
// Orchestrierung liegt im Application Service (`ARC-002`).

package inbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// BackfillRequestCommand trägt die Eingabe der Annahme (`LH-FA-CAP-009`):
// die Kennung des Antrags — sie wird die Run-Kennung (`SPEC-029`) —, die
// Tabelle und die Publication der Quelle, deren Mitgliedschaft die
// Vorbedingung mitträgt.
type BackfillRequestCommand struct {
	RequestID   model.AdministrationRequestID
	Source      model.SourceID
	Schema      string
	Table       string
	Publication string
}

// BackfillRequestResult trägt den angenommenen Run im Zustand `queued`.
type BackfillRequestResult struct {
	Run model.BackfillRun
}

// BackfillExecuteCommand trägt die Eingabe der Ausführung: den
// angenommenen Run (Zustand `queued`) und die Publication der Quelle, gegen
// die die Vorbedingungen erneut geprüft werden.
type BackfillExecuteCommand struct {
	Run         model.BackfillRun
	Publication string
}

// BackfillExecuteResult trägt den Run in dem Zustand, in dem die Ausführung
// ihn verlässt.
type BackfillExecuteResult struct {
	Run model.BackfillRun
}

// BackfillTableUseCase überführt den Bestand einer aktivierten Tabelle in
// Backfill-Changes (`LH-FA-CAP-009`, `ADR-0111`).
type BackfillTableUseCase interface {
	// Request prüft die Vorbedingungen (Tabelle aktiviert mit laufender
	// Bindung und Mitglied der Publication), liest die geschätzte
	// Zeilenzahl und nimmt den Antrag als letzten Schritt an: die Run-Zeile
	// `queued` und der Antragsvermerk entstehen als eine Einheit. Eine
	// verletzte Vorbedingung und ein aktiver Run derselben Tabelle enden
	// ohne Run-Zeile.
	Request(ctx context.Context, command BackfillRequestCommand) (BackfillRequestResult, error)

	// Execute führt einen angenommenen Run aus. Ein Fehler des Runs endet
	// den Run `failed` mit seiner Fehlerklasse, ein beendeter eigener
	// Kontext `interrupted`; beides ist ein Ergebnis mit dem Run im
	// Endzustand, kein Fehler dieses Aufrufs. Der Fehler dieses Aufrufs
	// meldet, dass der Run-Zustand nicht festgehalten werden konnte oder
	// dass der Kontext endete, bevor der Run begann — das Ergebnis trägt
	// dann den zuletzt bekannten Zustand.
	Execute(ctx context.Context, command BackfillExecuteCommand) (BackfillExecuteResult, error)
}
