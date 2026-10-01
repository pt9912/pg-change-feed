// Package readchanges trägt den Changes-Lese-Use-Case (`ARC-002`):
// er ist eine dünne Fassade über dem bestehenden
// `ChangeStorePort` — keine eigene Abfrage, keine zweite Sortierung, kein
// zweiter Lesepfad.
package readchanges

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ReadChangesQuery, ReadChange und ReadChangesResult sind die
// Transport-Typen des Changes-Lese-Use-Cases (`ADR-0042`): ihre Definition
// liegt am Inbound-Port — der Port trägt seinen Vertrag einschließlich der
// Transport-Typen, Ports referenzieren nur die Domain —; die
// Use-Case-Adressen sind Aliase desselben Typs.
type (
	ReadChangesQuery  = inbound.ReadChangesQuery
	ReadChange        = inbound.ReadChange
	ReadChangesResult = inbound.ReadChangesResult
)

// ReadChangesService implementiert `inbound.ReadChangesUseCase` und liest
// die Changes über den `ChangeStorePort`. Die Filterachse
// (Quelle/Schema/Tabelle/Ziel), der Positionsbereich und das Limit gehen in die
// Abfrage des Ports; die Ordnung trägt der bestehende Lesepfad
// (`LH-FA-REA-004`).
type ReadChangesService struct {
	store outbound.ChangeStorePort
}

// NewReadChangesService verdrahtet den Use Case mit dem Change-Store-Port.
func NewReadChangesService(store outbound.ChangeStorePort) *ReadChangesService {
	return &ReadChangesService{store: store}
}

var _ inbound.ReadChangesUseCase = (*ReadChangesService)(nil)

// ReadChanges bildet `ReadChangesQuery` auf `outbound.ChangeQuery` ab und
// ruft `ChangeStorePort.ReadChanges`. Eine leere Quellen-Kennung ist eine
// ungültige Eingabe und endet über den bestehenden Fehlerpfad, keine stille
// Übernahme über alle Quellen. Den Lese-Kontrakt (Limit, Bereich, Quelle der
// Positionen) prüft `outbound.ChangeQuery.Validate`, dieselbe Prüfung, die der
// Port führt; ihr Fehler kommt unverändert zurück, auch bei einem `Target`
// außerhalb des Alphabets. Eine leere Rückgabe trägt eine leere, gesetzte
// Liste. Ein gesetztes `Target` außerhalb des Alphabets des Zielnamens
// (`model.IsValidRouteTarget`) trifft nie ein vergebenes Ziel: bei sonst
// gültiger Anfrage ist die Antwort die leere, gesetzte Liste, ohne Fehler und
// ohne den Port anzufragen (`ADR-0139`).
func (s *ReadChangesService) ReadChanges(ctx context.Context, query ReadChangesQuery) (ReadChangesResult, error) {
	if query.Source == "" {
		return ReadChangesResult{}, domainerrors.ErrEmptyIdentifier
	}
	portQuery := outbound.ChangeQuery{
		Source: query.Source,
		Schema: query.Schema,
		Table:  query.Table,
		Target: query.Target,
		Start:  query.Start,
		End:    query.End,
		Limit:  query.Limit,
	}
	if query.Target != "" && !model.IsValidRouteTarget(query.Target) {
		if err := portQuery.Validate(); err != nil {
			return ReadChangesResult{}, err
		}
		return ReadChangesResult{Changes: []ReadChange{}}, nil
	}
	records, err := s.store.ReadChanges(ctx, portQuery)
	if err != nil {
		return ReadChangesResult{}, err
	}
	changes := make([]ReadChange, 0, len(records))
	for _, record := range records {
		changes = append(changes, ReadChange{
			Position:    record.Position,
			Change:      record.Change,
			CommittedAt: record.CommittedAt,
		})
	}
	return ReadChangesResult{Changes: changes}, nil
}
