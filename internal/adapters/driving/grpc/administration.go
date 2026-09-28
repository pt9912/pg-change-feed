// administration.go trägt den `Administration`-Service (`ADR-0131`): zehn
// unäre RPCs, jeder ruft denselben Inbound Use Case wie sein
// HTTP-Äquivalent in `internal/adapters/driving/http/{consumer,
// registerconsumer,verwaltung,retention,readchanges}.go` — kein zweiter
// Domänenpfad. Der Adapter importiert ausschließlich Inbound Ports und
// Domain-Typen zur Übersetzung, keine Driven-Adapter- oder
// Application-Interna, dieselbe Grenze wie beim bestehenden Stream-Handler
// und beim HTTP-Adapter.
package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// administrationService implementiert den generierten
// `AdministrationServer` (`SPEC-031`); jedes Feld trägt genau den Inbound
// Use Case, den auch der HTTP-Adapter für dieselbe Fähigkeit aufruft.
type administrationService struct {
	administrationv1.UnimplementedAdministrationServer
	registerConsumer    inbound.RegisterConsumerUseCase
	acknowledgeConsumer inbound.AcknowledgeConsumerUseCase
	getConsumerPosition inbound.GetConsumerPositionUseCase
	removeConsumer      inbound.RemoveConsumerUseCase
	enableTable         inbound.EnableTableUseCase
	disableTable        inbound.DisableTableUseCase
	getStatus           inbound.GetStatusUseCase
	listTables          inbound.ListTablesUseCase
	runRetention        inbound.RunRetentionUseCase
	readChanges         inbound.ReadChangesUseCase
	log                 outbound.LogPort
}

// administrationError bildet einen Use-Case-Fehler auf einen `codes.*`-Wert
// ab (`ADR-0130` Teilfrage 5): dieselbe Fehlerklassifikation wie
// `internal/adapters/driving/http/errors.go`s `writeDomainError`, zwei
// Zieldarstellungen. Ein JSON-Decode-Fehlerfall entfällt strukturell — das
// Nachrichtenschema ist bereits typisiert (Protobuf).
func administrationError(ctx context.Context, log outbound.LogPort, action string, err error) error {
	switch {
	case errors.Is(err, inbound.ErrSourceTableMissing):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domainerrors.ErrEmptyIdentifier),
		errors.Is(err, domainerrors.ErrInvalidPosition),
		errors.Is(err, domainerrors.ErrSourceMismatch),
		errors.Is(err, domainerrors.ErrPositionRegression),
		errors.Is(err, domainerrors.ErrNegativeDuration),
		errors.Is(err, domainerrors.ErrNonPositiveVersion),
		errors.Is(err, outbound.ErrNonPositiveLimit),
		errors.Is(err, outbound.ErrRangeInverted):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		log.Warn(ctx, "grpc: "+action+" fehlgeschlagen", "error", err)
		return status.Error(codes.Internal, "interner Fehler")
	}
}

// RegisterConsumer registriert einen benannten Consumer (`LH-FA-CON-001`);
// ruft denselben Use Case wie `registerConsumerHandler` (HTTP-Adapter).
func (s *administrationService) RegisterConsumer(ctx context.Context, req *administrationv1.RegisterConsumerRequest) (*administrationv1.RegisterConsumerResponse, error) {
	result, err := s.registerConsumer.Register(ctx, inbound.RegisterConsumerCommand{
		Consumer: model.ConsumerID(req.GetConsumerId()),
		Name:     req.GetName(),
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "RegisterConsumer", err)
	}
	return &administrationv1.RegisterConsumerResponse{
		ConsumerId:        string(result.Consumer.ID),
		Name:              result.Consumer.Name,
		AlreadyRegistered: result.AlreadyRegistered,
	}, nil
}

// AcknowledgeConsumer bestätigt die verarbeitete Position eines Consumers
// (`LH-FA-CON-004`); ruft denselben Use Case wie `acknowledgeConsumerHandler`.
func (s *administrationService) AcknowledgeConsumer(ctx context.Context, req *administrationv1.AcknowledgeConsumerRequest) (*administrationv1.AcknowledgeConsumerResponse, error) {
	result, err := s.acknowledgeConsumer.Acknowledge(ctx, inbound.AcknowledgeConsumerCommand{
		Consumer: model.ConsumerID(req.GetConsumerId()),
		Position: model.SourcePosition{SourceID: model.SourceID(req.GetSourceId()), Offset: req.GetOffset()},
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "AcknowledgeConsumer", err)
	}
	return &administrationv1.AcknowledgeConsumerResponse{
		ConsumerId: string(result.Position.ConsumerID),
		SourceId:   string(result.Position.Position.SourceID),
		Offset:     result.Position.Position.Offset,
	}, nil
}

// GetConsumerPosition liest die bestätigte Position eines Consumers
// (`LH-FA-CON-003`); ruft denselben Use Case wie
// `getConsumerPositionHandler`.
func (s *administrationService) GetConsumerPosition(ctx context.Context, req *administrationv1.GetConsumerPositionRequest) (*administrationv1.GetConsumerPositionResponse, error) {
	consumerID := req.GetConsumerId()
	result, err := s.getConsumerPosition.Position(ctx, inbound.GetConsumerPositionQuery{
		Consumer: model.ConsumerID(consumerID),
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "GetConsumerPosition", err)
	}
	return &administrationv1.GetConsumerPositionResponse{
		ConsumerId:   consumerID,
		SourceId:     string(result.Position.Position.SourceID),
		Offset:       result.Position.Position.Offset,
		Acknowledged: result.Position.Acknowledged(),
	}, nil
}

// RemoveConsumer entfernt einen Consumer administrativ (`LH-FA-CON-006`);
// ruft denselben Use Case wie `removeConsumerHandler`.
func (s *administrationService) RemoveConsumer(ctx context.Context, req *administrationv1.RemoveConsumerRequest) (*administrationv1.RemoveConsumerResponse, error) {
	consumerID := req.GetConsumerId()
	result, err := s.removeConsumer.Remove(ctx, inbound.RemoveConsumerCommand{
		Consumer: model.ConsumerID(consumerID),
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "RemoveConsumer", err)
	}
	return &administrationv1.RemoveConsumerResponse{
		ConsumerId: consumerID,
		Removed:    result.Removed,
	}, nil
}

// EnableTable aktiviert CDC für eine einzelne Tabelle (`LH-FA-CFG-001`);
// ruft denselben Use Case wie `enableTableHandler`.
func (s *administrationService) EnableTable(ctx context.Context, req *administrationv1.EnableTableRequest) (*administrationv1.EnableTableResponse, error) {
	result, err := s.enableTable.Enable(ctx, inbound.EnableTableCommand{
		Source:          model.SourceID(req.GetSource()),
		Schema:          req.GetSchema(),
		Table:           req.GetTable(),
		TableID:         model.SourceTableID(req.GetTableId()),
		SchemaVersionID: model.SchemaVersionID(req.GetSchemaVersionId()),
		Version:         req.GetVersion(),
		Publication:     req.GetPublication(),
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "EnableTable", err)
	}
	return &administrationv1.EnableTableResponse{
		TableId:        string(result.Table.ID),
		Source:         string(result.Table.SourceID),
		Schema:         result.Table.Schema,
		Table:          result.Table.Table,
		AlreadyEnabled: result.AlreadyEnabled,
	}, nil
}

// DisableTable deaktiviert CDC für eine einzelne Tabelle (`LH-FA-CFG-002`);
// ruft denselben Use Case wie `disableTableHandler`.
func (s *administrationService) DisableTable(ctx context.Context, req *administrationv1.DisableTableRequest) (*administrationv1.DisableTableResponse, error) {
	result, err := s.disableTable.Disable(ctx, inbound.DisableTableCommand{
		Source:      model.SourceID(req.GetSource()),
		Schema:      req.GetSchema(),
		Table:       req.GetTable(),
		Publication: req.GetPublication(),
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "DisableTable", err)
	}
	return &administrationv1.DisableTableResponse{
		Removed:  result.Removed,
		Retained: result.Retained,
	}, nil
}

// GetTableStatus meldet den CDC-Zustand einer Tabelle (`LH-FA-CFG-003`);
// ruft denselben Use Case wie `getStatusHandler`.
func (s *administrationService) GetTableStatus(ctx context.Context, req *administrationv1.GetTableStatusRequest) (*administrationv1.GetTableStatusResponse, error) {
	result, err := s.getStatus.Status(ctx, inbound.GetStatusQuery{
		Source:      model.SourceID(req.GetSource()),
		Schema:      req.GetSchema(),
		Table:       req.GetTable(),
		Publication: req.GetPublication(),
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "GetTableStatus", err)
	}
	return &administrationv1.GetTableStatusResponse{
		Enabled:  result.Enabled,
		Retained: result.Retained,
	}, nil
}

// toAdministrationSourceTables übersetzt Domain-Tabellen in ihre
// Protobuf-Form (`ListTables`) — dieselben vier Felder wie
// `toSourceTableResponse` im HTTP-Adapter, nie `nil` für eine leere Menge.
func toAdministrationSourceTables(tables []model.SourceTable) []*administrationv1.SourceTable {
	out := make([]*administrationv1.SourceTable, 0, len(tables))
	for _, table := range tables {
		out = append(out, &administrationv1.SourceTable{
			TableId: string(table.ID),
			Source:  string(table.SourceID),
			Schema:  table.Schema,
			Table:   table.Table,
		})
	}
	return out
}

// ListTables listet die aktivierten Tabellen einer Quelle (`LH-FA-CFG-004`);
// ruft denselben Use Case wie `listTablesHandler`.
func (s *administrationService) ListTables(ctx context.Context, req *administrationv1.ListTablesRequest) (*administrationv1.ListTablesResponse, error) {
	result, err := s.listTables.ListTables(ctx, inbound.ListTablesQuery{
		Source:      model.SourceID(req.GetSource()),
		Publication: req.GetPublication(),
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "ListTables", err)
	}
	return &administrationv1.ListTablesResponse{
		Tables:   toAdministrationSourceTables(result.Tables),
		Retained: toAdministrationSourceTables(result.Retained),
	}, nil
}

// RunRetention führt eine Bereinigung für eine Quelle aus
// (`LH-FA-RET-002`); ruft denselben Use Case wie `runRetentionHandler`.
// Die Policy entsteht über den Domänen-Konstruktor — eine negative Dauer
// endet über die Invariante, bevor der Use Case aufgerufen wird.
func (s *administrationService) RunRetention(ctx context.Context, req *administrationv1.RunRetentionRequest) (*administrationv1.RunRetentionResponse, error) {
	minAge, err := model.NewDuration(req.GetMinAgeNanos())
	if err != nil {
		return nil, administrationError(ctx, s.log, "RunRetention", err)
	}
	policy, err := model.NewRetentionPolicy(minAge)
	if err != nil {
		return nil, administrationError(ctx, s.log, "RunRetention", err)
	}
	result, err := s.runRetention.Run(ctx, inbound.RunRetentionCommand{
		Source: model.SourceID(req.GetSource()),
		Policy: policy,
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "RunRetention", err)
	}
	return &administrationv1.RunRetentionResponse{Deleted: int64(result.Deleted)}, nil
}

// readChangesPosition übersetzt einen `from`/`to`-Wert des Requests in eine
// optionale Positions-Grenze: `0` trägt „nicht gesetzt" (`ADR-0131`
// Teilfrage 3), eine gültige Quellposition ist immer `≥ 1`.
func readChangesPosition(source model.SourceID, offset uint64) (*model.SourcePosition, error) {
	if offset == 0 {
		return nil, nil
	}
	position, err := model.NewSourcePosition(source, offset)
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// toChangeRecord übersetzt einen gelesenen Change in seine Protobuf-Form
// (`ADR-0131` Teilfrage 3) — dieselben dreizehn Felder wie
// `readChangeResponse` im HTTP-Adapter, in derselben Reihenfolge.
func toChangeRecord(change inbound.ReadChange) *administrationv1.ChangeRecord {
	return &administrationv1.ChangeRecord{
		CommitPosition: int64(change.Position.Offset),
		ChangeId:       string(change.Change.ID),
		TransactionId:  string(change.Change.TransactionID),
		SourceTableId:  string(change.Change.SourceTableID),
		Schema:         change.Change.Schema,
		Table:          change.Change.Table,
		Sequence:       change.Change.Sequence,
		Operation:      string(change.Change.Operation),
		OldImage:       change.Change.OldImage,
		NewImage:       change.Change.NewImage,
		SchemaVersion:  string(change.Change.SchemaVersion),
		CommittedAt:    time.Unix(0, change.CommittedAt.UnixNanos).UTC().Format(time.RFC3339Nano),
		Origin:         string(change.Change.Origin.OrDefault()),
	}
}

// ReadChanges liest persistierte Changes über einen begrenzten Bereich
// (`ADR-0131`); ruft denselben Inbound Use Case wie `readChangesHandler`
// (HTTP-Adapter) — kein zweiter Lesepfad.
func (s *administrationService) ReadChanges(ctx context.Context, req *administrationv1.ReadChangesRequest) (*administrationv1.ReadChangesResponse, error) {
	source := model.SourceID(req.GetSource())
	start, err := readChangesPosition(source, req.GetFrom())
	if err != nil {
		return nil, administrationError(ctx, s.log, "ReadChanges", err)
	}
	end, err := readChangesPosition(source, req.GetTo())
	if err != nil {
		return nil, administrationError(ctx, s.log, "ReadChanges", err)
	}
	var limit *int
	if l := req.GetLimit(); l != 0 {
		v := int(l)
		limit = &v
	}
	result, err := s.readChanges.ReadChanges(ctx, inbound.ReadChangesQuery{
		Source: source,
		Schema: req.GetSchema(),
		Table:  req.GetTable(),
		Start:  start,
		End:    end,
		Limit:  limit,
	})
	if err != nil {
		return nil, administrationError(ctx, s.log, "ReadChanges", err)
	}
	changes := make([]*administrationv1.ChangeRecord, 0, len(result.Changes))
	for _, change := range result.Changes {
		changes = append(changes, toChangeRecord(change))
	}
	return &administrationv1.ReadChangesResponse{Changes: changes}, nil
}
