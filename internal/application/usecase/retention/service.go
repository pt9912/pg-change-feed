// Package retention trägt den RunRetention Use Case (`ARC-002`,
// `ADR-0014`): die Freigabe je Change trägt
// `model.RetentionPolicy.AllowsDeletion`, die physische Löschung der
// freigegebenen Menge trägt der `ChangeStorePort`.
package retention

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// RunRetentionCommand und RunRetentionResult sind die Transport-Typen des
// RunRetention Use Cases (`ADR-0042`): ihre Definition liegt am
// Inbound-Port — der Port trägt seinen Vertrag einschließlich der
// Transport-Typen, Ports referenzieren nur die Domain —; die Use-Case-
// Adressen sind Aliase desselben Typs.
type (
	RunRetentionCommand = inbound.RunRetentionCommand
	RunRetentionResult  = inbound.RunRetentionResult
)

// RunRetentionService implementiert `inbound.RunRetentionUseCase` und
// orchestriert die Bereinigung über drei Ports: der `ConsumerStatePort`
// trägt die bestätigten Consumer-Positionen der Quelle, der `ClockPort`
// die Wanduhr für das Alter jedes Changes (`ADR-0040`), der
// `ChangeStorePort` liest die Kandidaten und führt die physische Löschung
// der freigegebenen Menge aus.
type RunRetentionService struct {
	store outbound.ChangeStorePort
	state outbound.ConsumerStatePort
	clock outbound.ClockPort
}

// NewRunRetentionService verdrahtet den Use Case mit seinen drei Ports.
func NewRunRetentionService(store outbound.ChangeStorePort, state outbound.ConsumerStatePort, clock outbound.ClockPort) *RunRetentionService {
	return &RunRetentionService{store: store, state: state, clock: clock}
}

var _ inbound.RunRetentionUseCase = (*RunRetentionService)(nil)

// PageSize ist die Zahl der Kandidaten je Lese-Aufruf eines Laufs
// (`ADR-0124`): der Arbeitsspeicher eines Laufs hängt an ihr, nicht an der
// Zahl der gespeicherten Changes.
const PageSize = 10_000

// Run liest die bestätigten Consumer-Positionen und die Wanduhr einmal und
// dann die Kandidaten der Quelle seitenweise (`PageSize`, ohne Row Images).
// Je Seite befragt der Lauf `RetentionPolicy.AllowsDeletion` je Kandidat
// (Alter aus `CommittedAt` gegen die Wanduhr, Change-Position, alle
// bestätigten Consumer-Positionen der Quelle, `LH-FA-RET-002`…`004`) und
// übergibt ausschließlich die freigegebene Menge dieser Seite an
// `ChangeStorePort.DeleteChanges`; nur eine leere Seite beendet den Lauf.
// Der Lauf ist über die Seiten nicht atomar: ein Fehler hinterlässt die
// Löschungen der Seiten davor, der nächste Lauf setzt fort (`ADR-0124`).
// Eine Seite, deren letzte Kennung der Cursor ist, verletzt den Seitenvertrag
// des Ports und endet als Fehler der Klasse `storage`, nicht als Endlosschleife.
// Eine leere Quellen-Kennung ist eine ungültige Konfiguration und endet über
// einen expliziten Fehlerpfad, keine stille Übernahme (`LH-FA-RET-002`
// Negative).
func (s *RunRetentionService) Run(ctx context.Context, command RunRetentionCommand) (RunRetentionResult, error) {
	if command.Source == "" {
		return RunRetentionResult{}, domainerrors.ErrEmptyIdentifier
	}

	consumerPositions, err := s.state.Positions(ctx, command.Source)
	if err != nil {
		return RunRetentionResult{}, err
	}
	now := s.clock.Now()

	deleted := 0
	var after model.ChangeID
	for {
		page, err := s.store.ReadRetentionCandidates(ctx, command.Source, after, PageSize)
		if err != nil {
			return RunRetentionResult{}, err
		}
		if len(page) == 0 {
			return RunRetentionResult{Deleted: deleted}, nil
		}

		eligible := make([]model.ChangeID, 0, len(page))
		for _, candidate := range page {
			age := now.Sub(candidate.CommittedAt)
			if command.Policy.AllowsDeletion(age, candidate.Position, consumerPositions) {
				eligible = append(eligible, candidate.ChangeID)
			}
		}
		if err := s.store.DeleteChanges(ctx, eligible); err != nil {
			return RunRetentionResult{}, err
		}
		deleted += len(eligible)

		last := page[len(page)-1].ChangeID
		if last == after {
			return RunRetentionResult{}, fmt.Errorf("%w: Kandidaten-Seite ohne Fortschritt hinter %q", outbound.ErrStorage, after)
		}
		after = last
	}
}
