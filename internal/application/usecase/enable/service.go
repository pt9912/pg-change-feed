// Package enable trägt den EnableTable Use Case (`ARC-002`, `ADR-0028`):
// die Aktivierung einer Tabelle schreibt die Bindungs-Zeilen der
// CDC-Referenztabellen und trägt die Publication der Quelle
// (`LH-FA-CFG-001`).
package enable

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// EnableTableCommand und EnableTableResult sind die Transport-Typen des
// EnableTable Use Cases (`ADR-0042`): ihre Definition liegt am
// Inbound-Port — der Port trägt seinen Vertrag einschließlich der
// Transport-Typen, Ports referenzieren nur die Domain —; die Use-Case-
// Adressen sind Aliase desselben Typs.
type (
	EnableTableCommand = inbound.EnableTableCommand
	EnableTableResult  = inbound.EnableTableResult
)

// EnableTableService implementiert `inbound.EnableTableUseCase` und trägt
// die Aktivierung über den `TableActivationPort`: die Existenz der
// physischen Tabelle ist die Vorbedingung (`LH-FA-CFG-001` Negative:
// expliziter Fehlerpfad), die Bindungs-Zeilen und die Publication tragen
// den Stand — die `REPLICA IDENTITY` bleibt vom Aktivieren unberührt
// (`LH-FA-CFG-001.a` Schritt 1: nicht ungefragt aktiviert).
type EnableTableService struct {
	activation outbound.TableActivationPort
}

// NewEnableTableService verdrahtet den Use Case mit dem Aktivierungs-Port.
func NewEnableTableService(activation outbound.TableActivationPort) *EnableTableService {
	return &EnableTableService{activation: activation}
}

var _ inbound.EnableTableUseCase = (*EnableTableService)(nil)

// Enable aktiviert die Tabelle idempotent (`LH-FA-CFG-001` Boundary): die
// Bindungs-Zeile läuft vor der Publication — eine fehlgeschlagene
// Publication hinterlässt die Bindung ohne Erfassung, ein erneuter Aufruf
// trägt die Publication nach; umgekehrt würde eine Publication ohne
// Bindungs-Zeile die ersten Changes an der Fremdschlüssel-Kante verlieren.
// Eine bereits aktivierte Tabelle bleibt unverändert; die Rückkehr meldet
// den Ausgang über `AlreadyEnabled`.
func (s *EnableTableService) Enable(ctx context.Context, command EnableTableCommand) (EnableTableResult, error) {
	table, err := model.NewSourceTable(command.TableID, command.Source, command.Schema, command.Table)
	if err != nil {
		return EnableTableResult{}, err
	}
	version, err := model.NewSchemaVersion(command.SchemaVersionID, command.TableID, command.Version)
	if err != nil {
		return EnableTableResult{}, err
	}
	exists, err := s.activation.TableExists(ctx, command.Schema, command.Table)
	if err != nil {
		return EnableTableResult{}, err
	}
	if !exists {
		return EnableTableResult{}, fmt.Errorf("%w: %s.%s", inbound.ErrSourceTableMissing, command.Schema, command.Table)
	}
	registered, err := s.activation.Register(ctx, table, version)
	if err != nil {
		return EnableTableResult{}, err
	}
	if err := s.activation.Publish(ctx, command.Publication, command.Schema, command.Table); err != nil {
		return EnableTableResult{}, err
	}
	return EnableTableResult{Table: table, AlreadyEnabled: !registered}, nil
}
