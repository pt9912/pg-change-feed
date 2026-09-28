package bootstrap

import (
	"context"
	"fmt"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// syncAssemblerAddBinding trägt eine erfolgreich aktivierte Tabelle in die
// laufende `Assembler`-Bindung nach — dieselbe Lese-Folge wie
// `activatedTableBindings` beim Prozessstart, auf eine einzelne Tabelle
// zugeschnitten: die registrierte Bindungs-Zeile, ihre aktuelle
// Schema-Version, der dauerhafte Ausschluss- und Regelstand der Quelle. Dies
// ist der eine Mechanismus, den sowohl `applyAdministrationRequest`s
// Enable-Zweig (SQL-Antragsqueue) als auch `enableTableWithAssemblerSync`
// (direkter HTTP-/gRPC-Zugriffsweg) aufrufen — beide riefen bislang denselben
// `EnableTableUseCase` auf, aber nur der Antragsqueue-Pfad trug den Nachtrag
// in den laufenden Prozess; eine über HTTP oder gRPC aktivierte Tabelle blieb
// bis zum nächsten Neustart unerfasst (`LH-FA-CFG-001`).
func syncAssemblerAddBinding(ctx context.Context, assembler *mapper.Assembler, activation outbound.TableActivationPort, schemaStore outbound.SchemaStorePort, columnExclusion outbound.ColumnExclusionPort, transformations outbound.TransformationPort, source model.SourceID, schema, table string) error {
	qualified := schema + "." + table
	registered, found, err := activation.Registered(ctx, source, schema, table)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("Aktivierung ohne Bindungs-Zeile: %s", qualified)
	}
	current, found, err := schemaStore.CurrentVersion(ctx, registered.ID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("Aktivierung ohne registrierte Schema-Version: %s", qualified)
	}
	excluded, err := columnExclusion.ExcludedColumns(ctx, source)
	if err != nil {
		return err
	}
	rules, err := transformations.TransformationRules(ctx, source)
	if err != nil {
		return err
	}
	assembler.AddBinding(qualified, mapper.TableBinding{
		TableID:         registered.ID,
		SchemaVersion:   current.ID,
		ExcludedColumns: excluded[qualified],
		Transformations: rules[qualified],
	})
	return nil
}

// syncAssemblerRemoveBinding trägt eine erfolgreich deaktivierte Tabelle aus
// der laufenden `Assembler`-Bindung aus — das Gegenstück zu
// `syncAssemblerAddBinding`, derselbe gemeinsame Mechanismus für beide
// Aufrufpfade.
func syncAssemblerRemoveBinding(assembler *mapper.Assembler, schema, table string) {
	assembler.RemoveBinding(schema + "." + table)
}

// enableTableWithAssemblerSync dekoriert den `EnableTableUseCase` des
// direkten HTTP-/gRPC-Zugriffswegs: nach erfolgreichem `Enable` trägt sie die
// Bindung über `syncAssemblerAddBinding` in den laufenden `Assembler` nach —
// ohne diesen Nachtrag verliert eine über diesen Weg aktivierte Tabelle jede
// Änderung bis zum nächsten Prozess-Neustart (`LH-FA-CFG-001`). Ein Fehler
// des Nachtrags trägt den Aufruf insgesamt als Fehlschlag, mit leerem
// Ergebnis — derselbe Ausgang wie ein Fehler des inneren Use Cases selbst.
type enableTableWithAssemblerSync struct {
	inbound.EnableTableUseCase
	assembler       *mapper.Assembler
	activation      outbound.TableActivationPort
	schemaStore     outbound.SchemaStorePort
	columnExclusion outbound.ColumnExclusionPort
	transformations outbound.TransformationPort
}

var _ inbound.EnableTableUseCase = enableTableWithAssemblerSync{}

// Enable ruft den inneren Use Case auf und trägt den laufenden Assembler bei
// Erfolg nach; ein gescheiterter innerer Aufruf lässt den Assembler
// unangetastet.
func (e enableTableWithAssemblerSync) Enable(ctx context.Context, command inbound.EnableTableCommand) (inbound.EnableTableResult, error) {
	result, err := e.EnableTableUseCase.Enable(ctx, command)
	if err != nil {
		return result, err
	}
	if err := syncAssemblerAddBinding(ctx, e.assembler, e.activation, e.schemaStore, e.columnExclusion, e.transformations, command.Source, command.Schema, command.Table); err != nil {
		return inbound.EnableTableResult{}, err
	}
	return result, nil
}

// disableTableWithAssemblerSync spiegelt `enableTableWithAssemblerSync` für
// den Deaktivierungs-Pfad: nach erfolgreichem `Disable` entfernt sie die
// Bindung über `syncAssemblerRemoveBinding` aus dem laufenden `Assembler`.
type disableTableWithAssemblerSync struct {
	inbound.DisableTableUseCase
	assembler *mapper.Assembler
}

var _ inbound.DisableTableUseCase = disableTableWithAssemblerSync{}

// Disable ruft den inneren Use Case auf und entfernt die Bindung bei Erfolg
// aus dem laufenden Assembler; ein gescheiterter innerer Aufruf lässt den
// Assembler unangetastet.
func (d disableTableWithAssemblerSync) Disable(ctx context.Context, command inbound.DisableTableCommand) (inbound.DisableTableResult, error) {
	result, err := d.DisableTableUseCase.Disable(ctx, command)
	if err != nil {
		return result, err
	}
	syncAssemblerRemoveBinding(d.assembler, command.Schema, command.Table)
	return result, nil
}
