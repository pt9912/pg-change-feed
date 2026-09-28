package grpc

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// TestAdministrationErrorBildetJedeFehlerklasseAufIhrenCodeAb trägt die
// vollständige Übersetzungstabelle aus `ADR-0130` Teilfrage 5: jede
// benannte Fehlerklasse bildet auf genau ihren `codes.*`-Wert ab, jeder
// übrige Fehler auf `codes.Internal`.
//
// Rot färbende Mutation: einen der `case`-Zweige aus dem
// `errors.Is`-Vergleich in `administrationError` entfernen — der
// betroffene Fall fällt dann auf `codes.Internal` statt auf den
// erwarteten Code, dieser Test färbt für genau diesen Fall rot.
func TestAdministrationErrorBildetJedeFehlerklasseAufIhrenCodeAb(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want codes.Code
	}{
		{"SourceTableMissing", inbound.ErrSourceTableMissing, codes.NotFound},
		{"EmptyIdentifier", domainerrors.ErrEmptyIdentifier, codes.InvalidArgument},
		{"InvalidPosition", domainerrors.ErrInvalidPosition, codes.InvalidArgument},
		{"SourceMismatch", domainerrors.ErrSourceMismatch, codes.InvalidArgument},
		{"PositionRegression", domainerrors.ErrPositionRegression, codes.InvalidArgument},
		{"NegativeDuration", domainerrors.ErrNegativeDuration, codes.InvalidArgument},
		{"NonPositiveVersion", domainerrors.ErrNonPositiveVersion, codes.InvalidArgument},
		{"NonPositiveLimit", outbound.ErrNonPositiveLimit, codes.InvalidArgument},
		{"RangeInverted", outbound.ErrRangeInverted, codes.InvalidArgument},
		{"UnbekannterFehler", stderrors.New("etwas Unerwartetes"), codes.Internal},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := administrationError(context.Background(), outbound.NoopLog, "Test", tc.err)
			if status.Code(err) != tc.want {
				t.Fatalf("administrationError(%v) = %v (Erwartung: %v)", tc.err, status.Code(err), tc.want)
			}
		})
	}
}

// Die folgenden Fälschungen tragen je eine Inbound-Use-Case-Fälschung
// (`ADR-0130` Teilfrage 6): jede zeichnet die empfangene Eingabe auf und
// liefert eine vorgegebene Ausgabe — dieselbe Bauform, mit der
// `server_test.go`s `fakeSubscriber` den Driven-Adapter fälscht.

type fakeRegisterConsumer struct {
	gotCmd inbound.RegisterConsumerCommand
	result inbound.RegisterConsumerResult
	err    error
}

func (f *fakeRegisterConsumer) Register(_ context.Context, cmd inbound.RegisterConsumerCommand) (inbound.RegisterConsumerResult, error) {
	f.gotCmd = cmd
	return f.result, f.err
}

type fakeAcknowledgeConsumer struct {
	gotCmd inbound.AcknowledgeConsumerCommand
	result inbound.AcknowledgeConsumerResult
	err    error
}

func (f *fakeAcknowledgeConsumer) Acknowledge(_ context.Context, cmd inbound.AcknowledgeConsumerCommand) (inbound.AcknowledgeConsumerResult, error) {
	f.gotCmd = cmd
	return f.result, f.err
}

type fakeGetConsumerPosition struct {
	gotQuery inbound.GetConsumerPositionQuery
	result   inbound.GetConsumerPositionResult
	err      error
}

func (f *fakeGetConsumerPosition) Position(_ context.Context, query inbound.GetConsumerPositionQuery) (inbound.GetConsumerPositionResult, error) {
	f.gotQuery = query
	return f.result, f.err
}

type fakeRemoveConsumer struct {
	gotCmd inbound.RemoveConsumerCommand
	result inbound.RemoveConsumerResult
	err    error
}

func (f *fakeRemoveConsumer) Remove(_ context.Context, cmd inbound.RemoveConsumerCommand) (inbound.RemoveConsumerResult, error) {
	f.gotCmd = cmd
	return f.result, f.err
}

type fakeEnableTable struct {
	gotCmd inbound.EnableTableCommand
	result inbound.EnableTableResult
	err    error
}

func (f *fakeEnableTable) Enable(_ context.Context, cmd inbound.EnableTableCommand) (inbound.EnableTableResult, error) {
	f.gotCmd = cmd
	return f.result, f.err
}

type fakeDisableTable struct {
	gotCmd inbound.DisableTableCommand
	result inbound.DisableTableResult
	err    error
}

func (f *fakeDisableTable) Disable(_ context.Context, cmd inbound.DisableTableCommand) (inbound.DisableTableResult, error) {
	f.gotCmd = cmd
	return f.result, f.err
}

type fakeGetStatus struct {
	gotQuery inbound.GetStatusQuery
	result   inbound.GetStatusResult
	err      error
}

func (f *fakeGetStatus) Status(_ context.Context, query inbound.GetStatusQuery) (inbound.GetStatusResult, error) {
	f.gotQuery = query
	return f.result, f.err
}

type fakeListTables struct {
	gotQuery inbound.ListTablesQuery
	result   inbound.ListTablesResult
	err      error
}

func (f *fakeListTables) ListTables(_ context.Context, query inbound.ListTablesQuery) (inbound.ListTablesResult, error) {
	f.gotQuery = query
	return f.result, f.err
}

type fakeRunRetention struct {
	gotCmd inbound.RunRetentionCommand
	result inbound.RunRetentionResult
	err    error
}

func (f *fakeRunRetention) Run(_ context.Context, cmd inbound.RunRetentionCommand) (inbound.RunRetentionResult, error) {
	f.gotCmd = cmd
	return f.result, f.err
}

type fakeReadChanges struct {
	called   bool
	gotQuery inbound.ReadChangesQuery
	result   inbound.ReadChangesResult
	err      error
}

func (f *fakeReadChanges) ReadChanges(_ context.Context, query inbound.ReadChangesQuery) (inbound.ReadChangesResult, error) {
	f.called = true
	f.gotQuery = query
	return f.result, f.err
}

type fakeDiagnose struct {
	called   bool
	gotQuery inbound.DiagnoseQuery
	result   inbound.DiagnoseResult
	err      error
}

func (f *fakeDiagnose) Diagnose(_ context.Context, query inbound.DiagnoseQuery) (inbound.DiagnoseResult, error) {
	f.called = true
	f.gotQuery = query
	return f.result, f.err
}

// TestRegisterConsumerRuftUseCaseMitUebersetzterCommandAuf trägt die
// Bindung an die Eingabeseite (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`):
// die Felder des Requests erreichen den Use Case unverändert, das Ergebnis
// erreicht unverändert die Response.
func TestRegisterConsumerRuftUseCaseMitUebersetzterCommandAuf(t *testing.T) {
	fake := &fakeRegisterConsumer{result: inbound.RegisterConsumerResult{
		Consumer:          model.Consumer{ID: "c-1", Name: "Verbraucher"},
		AlreadyRegistered: true,
	}}
	svc := &administrationService{registerConsumer: fake, log: outbound.NoopLog}

	resp, err := svc.RegisterConsumer(context.Background(), &administrationv1.RegisterConsumerRequest{
		ConsumerId: "c-1", Name: "Verbraucher",
	})
	if err != nil {
		t.Fatalf("RegisterConsumer: %v", err)
	}
	if fake.gotCmd.Consumer != "c-1" || fake.gotCmd.Name != "Verbraucher" {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v", fake.gotCmd)
	}
	if resp.GetConsumerId() != "c-1" || resp.GetName() != "Verbraucher" || !resp.GetAlreadyRegistered() {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", resp)
	}
}

// TestRegisterConsumerBildetFehlerAufInvalidArgumentAb trägt den
// Fehlerpfad: eine leere Kennung endet über `domainerrors.ErrEmptyIdentifier`
// als `codes.InvalidArgument`, nicht als `codes.Internal`.
func TestRegisterConsumerBildetFehlerAufInvalidArgumentAb(t *testing.T) {
	fake := &fakeRegisterConsumer{err: domainerrors.ErrEmptyIdentifier}
	svc := &administrationService{registerConsumer: fake, log: outbound.NoopLog}
	_, err := svc.RegisterConsumer(context.Background(), &administrationv1.RegisterConsumerRequest{})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.InvalidArgument)
	}
}

func TestAcknowledgeConsumerRuftUseCaseMitUebersetzterCommandAuf(t *testing.T) {
	fake := &fakeAcknowledgeConsumer{result: inbound.AcknowledgeConsumerResult{
		Position: model.ConsumerPosition{ConsumerID: "c-1", Position: model.SourcePosition{SourceID: "src-1", Offset: 7}},
	}}
	svc := &administrationService{acknowledgeConsumer: fake, log: outbound.NoopLog}

	resp, err := svc.AcknowledgeConsumer(context.Background(), &administrationv1.AcknowledgeConsumerRequest{
		ConsumerId: "c-1", SourceId: "src-1", Offset: 7,
	})
	if err != nil {
		t.Fatalf("AcknowledgeConsumer: %v", err)
	}
	if fake.gotCmd.Consumer != "c-1" || fake.gotCmd.Position.SourceID != "src-1" || fake.gotCmd.Position.Offset != 7 {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v", fake.gotCmd)
	}
	if resp.GetConsumerId() != "c-1" || resp.GetSourceId() != "src-1" || resp.GetOffset() != 7 {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", resp)
	}
}

func TestAcknowledgeConsumerBildetFehlerAufInvalidArgumentAb(t *testing.T) {
	fake := &fakeAcknowledgeConsumer{err: domainerrors.ErrPositionRegression}
	svc := &administrationService{acknowledgeConsumer: fake, log: outbound.NoopLog}
	_, err := svc.AcknowledgeConsumer(context.Background(), &administrationv1.AcknowledgeConsumerRequest{})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.InvalidArgument)
	}
}

func TestGetConsumerPositionRuftUseCaseMitUebersetzterQueryAuf(t *testing.T) {
	position, err := model.NewSourcePosition("src-1", 3)
	if err != nil {
		t.Fatalf("SourcePosition bauen: %v", err)
	}
	fake := &fakeGetConsumerPosition{result: inbound.GetConsumerPositionResult{
		Position: model.ConsumerPosition{ConsumerID: "c-1", Position: position},
	}}
	svc := &administrationService{getConsumerPosition: fake, log: outbound.NoopLog}

	resp, err := svc.GetConsumerPosition(context.Background(), &administrationv1.GetConsumerPositionRequest{ConsumerId: "c-1"})
	if err != nil {
		t.Fatalf("GetConsumerPosition: %v", err)
	}
	if fake.gotQuery.Consumer != "c-1" {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v", fake.gotQuery)
	}
	if resp.GetConsumerId() != "c-1" || resp.GetSourceId() != "src-1" || resp.GetOffset() != 3 || !resp.GetAcknowledged() {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", resp)
	}
}

func TestGetConsumerPositionBildetUnbekanntenFehlerAufInternalAb(t *testing.T) {
	fake := &fakeGetConsumerPosition{err: stderrors.New("Speicher nicht erreichbar")}
	svc := &administrationService{getConsumerPosition: fake, log: outbound.NoopLog}
	_, err := svc.GetConsumerPosition(context.Background(), &administrationv1.GetConsumerPositionRequest{ConsumerId: "c-1"})
	if status.Code(err) != codes.Internal {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.Internal)
	}
}

func TestRemoveConsumerRuftUseCaseMitUebersetzterCommandAuf(t *testing.T) {
	fake := &fakeRemoveConsumer{result: inbound.RemoveConsumerResult{Removed: true}}
	svc := &administrationService{removeConsumer: fake, log: outbound.NoopLog}

	resp, err := svc.RemoveConsumer(context.Background(), &administrationv1.RemoveConsumerRequest{ConsumerId: "c-1"})
	if err != nil {
		t.Fatalf("RemoveConsumer: %v", err)
	}
	if fake.gotCmd.Consumer != "c-1" {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v", fake.gotCmd)
	}
	if resp.GetConsumerId() != "c-1" || !resp.GetRemoved() {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", resp)
	}
}

func TestEnableTableRuftUseCaseMitUebersetzterCommandAuf(t *testing.T) {
	fake := &fakeEnableTable{result: inbound.EnableTableResult{
		Table:          model.SourceTable{ID: "t-1", SourceID: "src-1", Schema: "public", Table: "orders"},
		AlreadyEnabled: false,
	}}
	svc := &administrationService{enableTable: fake, log: outbound.NoopLog}

	resp, err := svc.EnableTable(context.Background(), &administrationv1.EnableTableRequest{
		Source: "src-1", Schema: "public", Table: "orders",
		TableId: "t-1", SchemaVersionId: "t-1-v1", Version: 1, Publication: "pub",
	})
	if err != nil {
		t.Fatalf("EnableTable: %v", err)
	}
	if fake.gotCmd.Source != "src-1" || fake.gotCmd.Schema != "public" || fake.gotCmd.Table != "orders" ||
		fake.gotCmd.TableID != "t-1" || fake.gotCmd.SchemaVersionID != "t-1-v1" || fake.gotCmd.Version != 1 ||
		fake.gotCmd.Publication != "pub" {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v", fake.gotCmd)
	}
	if resp.GetTableId() != "t-1" || resp.GetSource() != "src-1" || resp.GetSchema() != "public" ||
		resp.GetTable() != "orders" || resp.GetAlreadyEnabled() {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", resp)
	}
}

func TestEnableTableBildetSourceTableMissingAufNotFoundAb(t *testing.T) {
	fake := &fakeEnableTable{err: inbound.ErrSourceTableMissing}
	svc := &administrationService{enableTable: fake, log: outbound.NoopLog}
	_, err := svc.EnableTable(context.Background(), &administrationv1.EnableTableRequest{})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.NotFound)
	}
}

func TestDisableTableRuftUseCaseMitUebersetzterCommandAuf(t *testing.T) {
	fake := &fakeDisableTable{result: inbound.DisableTableResult{Removed: true, Retained: false}}
	svc := &administrationService{disableTable: fake, log: outbound.NoopLog}

	resp, err := svc.DisableTable(context.Background(), &administrationv1.DisableTableRequest{
		Source: "src-1", Schema: "public", Table: "orders", Publication: "pub",
	})
	if err != nil {
		t.Fatalf("DisableTable: %v", err)
	}
	if fake.gotCmd.Source != "src-1" || fake.gotCmd.Schema != "public" || fake.gotCmd.Table != "orders" || fake.gotCmd.Publication != "pub" {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v", fake.gotCmd)
	}
	if !resp.GetRemoved() || resp.GetRetained() {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", resp)
	}
}

func TestDisableTableBildetSourceTableMissingAufNotFoundAb(t *testing.T) {
	fake := &fakeDisableTable{err: inbound.ErrSourceTableMissing}
	svc := &administrationService{disableTable: fake, log: outbound.NoopLog}
	_, err := svc.DisableTable(context.Background(), &administrationv1.DisableTableRequest{})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.NotFound)
	}
}

func TestGetTableStatusRuftUseCaseMitUebersetzterQueryAuf(t *testing.T) {
	fake := &fakeGetStatus{result: inbound.GetStatusResult{Enabled: true, Retained: false}}
	svc := &administrationService{getStatus: fake, log: outbound.NoopLog}

	resp, err := svc.GetTableStatus(context.Background(), &administrationv1.GetTableStatusRequest{
		Source: "src-1", Schema: "public", Table: "orders", Publication: "pub",
	})
	if err != nil {
		t.Fatalf("GetTableStatus: %v", err)
	}
	if fake.gotQuery.Source != "src-1" || fake.gotQuery.Schema != "public" || fake.gotQuery.Table != "orders" || fake.gotQuery.Publication != "pub" {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v", fake.gotQuery)
	}
	if !resp.GetEnabled() || resp.GetRetained() {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", resp)
	}
}

func TestGetTableStatusBildetSourceTableMissingAufNotFoundAb(t *testing.T) {
	fake := &fakeGetStatus{err: inbound.ErrSourceTableMissing}
	svc := &administrationService{getStatus: fake, log: outbound.NoopLog}
	_, err := svc.GetTableStatus(context.Background(), &administrationv1.GetTableStatusRequest{})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.NotFound)
	}
}

func TestListTablesRuftUseCaseMitUebersetzterQueryAuf(t *testing.T) {
	fake := &fakeListTables{result: inbound.ListTablesResult{
		Tables:   []model.SourceTable{{ID: "t-1", SourceID: "src-1", Schema: "public", Table: "orders"}},
		Retained: []model.SourceTable{{ID: "t-2", SourceID: "src-1", Schema: "public", Table: "legacy"}},
	}}
	svc := &administrationService{listTables: fake, log: outbound.NoopLog}

	resp, err := svc.ListTables(context.Background(), &administrationv1.ListTablesRequest{Source: "src-1", Publication: "pub"})
	if err != nil {
		t.Fatalf("ListTables: %v", err)
	}
	if fake.gotQuery.Source != "src-1" || fake.gotQuery.Publication != "pub" {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v", fake.gotQuery)
	}
	if len(resp.GetTables()) != 1 || resp.GetTables()[0].GetTableId() != "t-1" ||
		len(resp.GetRetained()) != 1 || resp.GetRetained()[0].GetTableId() != "t-2" {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", resp)
	}
}

// TestListTablesLeereMengenSindNieNil trägt dieselbe Zusage wie
// `toSourceTableResponse` im HTTP-Adapter: eine leere Menge trägt eine
// leere, gesetzte Liste, nie `nil`.
func TestListTablesLeereMengenSindNieNil(t *testing.T) {
	fake := &fakeListTables{result: inbound.ListTablesResult{}}
	svc := &administrationService{listTables: fake, log: outbound.NoopLog}
	resp, err := svc.ListTables(context.Background(), &administrationv1.ListTablesRequest{Source: "src-1", Publication: "pub"})
	if err != nil {
		t.Fatalf("ListTables: %v", err)
	}
	if resp.GetTables() == nil || resp.GetRetained() == nil {
		t.Fatalf("leere Mengen sind nil statt einer leeren Liste: %+v", resp)
	}
}

func TestRunRetentionRuftUseCaseMitUebersetzterCommandAuf(t *testing.T) {
	fake := &fakeRunRetention{result: inbound.RunRetentionResult{Deleted: 5}}
	svc := &administrationService{runRetention: fake, log: outbound.NoopLog}

	resp, err := svc.RunRetention(context.Background(), &administrationv1.RunRetentionRequest{Source: "src-1", MinAgeNanos: 1000})
	if err != nil {
		t.Fatalf("RunRetention: %v", err)
	}
	if fake.gotCmd.Source != "src-1" || fake.gotCmd.Policy.MinAge.Nanos != 1000 {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v", fake.gotCmd)
	}
	if resp.GetDeleted() != 5 {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", resp)
	}
}

// TestRunRetentionNegativeDauerEndetMitInvalidArgumentOhneUseCaseAufruf
// trägt die Domänen-Invariante vor dem Use-Case-Aufruf: eine negative
// `min_age_nanos` endet über `model.NewDuration`s Invariante, der Use Case
// wird nicht erreicht.
//
// Rot färbende Mutation: den `model.NewDuration`-Fehlerpfad in `RunRetention`
// verwerfen und stattdessen den Use Case mit der negativen Dauer aufrufen —
// `fake.gotCmd` bliebe dann nicht auf dem Nullwert, dieser Test färbt rot.
func TestRunRetentionNegativeDauerEndetMitInvalidArgumentOhneUseCaseAufruf(t *testing.T) {
	fake := &fakeRunRetention{}
	svc := &administrationService{runRetention: fake, log: outbound.NoopLog}
	_, err := svc.RunRetention(context.Background(), &administrationv1.RunRetentionRequest{Source: "src-1", MinAgeNanos: -1})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.InvalidArgument)
	}
	if fake.gotCmd != (inbound.RunRetentionCommand{}) {
		t.Fatalf("der Use Case wurde trotz ungültiger Dauer aufgerufen: %+v", fake.gotCmd)
	}
}

// TestReadChangesRuftUseCaseMitUebersetzterQueryAuf trägt die Bindung an
// die Eingabeseite (`ADR-0131` Teilfrage 3): Quelle, Klartext-Filter,
// Positionsbereich und Limit erreichen den Use Case unverändert übersetzt,
// das Ergebnis erreicht unverändert die Response — dieselben dreizehn
// Felder wie `readChangeResponse` im HTTP-Adapter.
func TestReadChangesRuftUseCaseMitUebersetzterQueryAuf(t *testing.T) {
	committedAt := time.Date(2026, 9, 28, 10, 30, 0, 123456789, time.UTC)
	position, err := model.NewSourcePosition("src-1", 4711)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	change, err := model.NewChange(
		model.ChangeID("c-1"), model.TransactionID("tx-1"), model.SourceTableID("tbl-1"),
		3, model.OperationUpdate, []byte(`{"id":1}`), []byte(`{"id":1,"name":"neu"}`), model.SchemaVersionID("sv-1"),
	)
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	change.Schema = "public"
	change.Table = "orders"
	fake := &fakeReadChanges{result: inbound.ReadChangesResult{Changes: []inbound.ReadChange{{
		Position:    position,
		Change:      change,
		CommittedAt: model.NewTimePoint(committedAt.UnixNano()),
	}}}}
	svc := &administrationService{readChanges: fake, log: outbound.NoopLog}

	resp, err := svc.ReadChanges(context.Background(), &administrationv1.ReadChangesRequest{
		Source: "src-1", Schema: "public", Table: "orders", From: 100, To: 500, Limit: 10,
	})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if !fake.called {
		t.Fatalf("der Use Case wurde nicht aufgerufen")
	}
	if fake.gotQuery.Source != "src-1" || fake.gotQuery.Schema != "public" || fake.gotQuery.Table != "orders" {
		t.Fatalf("Abfrage = %+v, wollen Quelle/schema/table übersetzt", fake.gotQuery)
	}
	if fake.gotQuery.Start == nil || fake.gotQuery.Start.Offset != 100 ||
		fake.gotQuery.End == nil || fake.gotQuery.End.Offset != 500 {
		t.Fatalf("Bereich = %+v, wollen [100,500)", fake.gotQuery)
	}
	if fake.gotQuery.Limit == nil || *fake.gotQuery.Limit != 10 {
		t.Fatalf("Limit = %v, wollen 10", fake.gotQuery.Limit)
	}
	if len(resp.GetChanges()) != 1 {
		t.Fatalf("Response trägt nicht genau einen Change: %+v", resp)
	}
	record := resp.GetChanges()[0]
	if record.GetCommitPosition() != 4711 || record.GetChangeId() != "c-1" || record.GetTransactionId() != "tx-1" ||
		record.GetSourceTableId() != "tbl-1" || record.GetSchema() != "public" || record.GetTable() != "orders" ||
		record.GetSequence() != 3 || record.GetOperation() != "UPDATE" ||
		string(record.GetOldImage()) != `{"id":1}` || string(record.GetNewImage()) != `{"id":1,"name":"neu"}` ||
		record.GetSchemaVersion() != "sv-1" || record.GetOrigin() != "wal" {
		t.Fatalf("Response trägt nicht das Use-Case-Ergebnis: %+v", record)
	}
	if want := committedAt.Format(time.RFC3339Nano); record.GetCommittedAt() != want {
		t.Fatalf("CommittedAt = %q, wollen %q", record.GetCommittedAt(), want)
	}
}

// TestReadChangesOhneBereichTraegtKeineGrenze trägt die Nullwert-Semantik
// aus `ADR-0131` Teilfrage 3: `from`/`to`/`limit` auf `0` bleiben
// unbegrenzt (`nil`) — keine `model.NewSourcePosition`-Invariante wird für
// eine nicht gesetzte Grenze geprüft.
func TestReadChangesOhneBereichTraegtKeineGrenze(t *testing.T) {
	fake := &fakeReadChanges{}
	svc := &administrationService{readChanges: fake, log: outbound.NoopLog}
	_, err := svc.ReadChanges(context.Background(), &administrationv1.ReadChangesRequest{Source: "src-1"})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if fake.gotQuery.Start != nil || fake.gotQuery.End != nil || fake.gotQuery.Limit != nil {
		t.Fatalf("Abfrage trägt eine Grenze trotz unbegrenzter Anfrage: %+v", fake.gotQuery)
	}
}

// TestReadChangesLeereQuelleMitPositionEndetMitInvalidArgumentOhneUseCaseAufruf
// trägt die Domänen-Invariante vor dem Use-Case-Aufruf
// (`domainerrors.ErrEmptyIdentifier`, `model.NewSourcePosition`): eine leere
// Quelle mit gesetztem `from` endet als `codes.InvalidArgument`, der Use
// Case wird nicht erreicht.
//
// Rot färbende Mutation: den Fehlerpfad von `readChangesPosition` verwerfen
// und stattdessen den Use Case mit der leeren Quelle aufrufen — `fake.called`
// bliebe dann `true`, dieser Test färbt rot.
func TestReadChangesLeereQuelleMitPositionEndetMitInvalidArgumentOhneUseCaseAufruf(t *testing.T) {
	fake := &fakeReadChanges{}
	svc := &administrationService{readChanges: fake, log: outbound.NoopLog}
	_, err := svc.ReadChanges(context.Background(), &administrationv1.ReadChangesRequest{From: 5})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.InvalidArgument)
	}
	if fake.called {
		t.Fatalf("der Use Case wurde trotz ungültiger Quelle aufgerufen: %+v", fake.gotQuery)
	}
}

// TestReadChangesLeereTrefferlisteBleibtGesetzteListe trägt dieselbe Zusage
// wie `TestListTablesLeereMengenSindNieNil`: eine leere Menge trägt eine
// leere, gesetzte Liste, nie `nil` (`LH-FA-REA-006` Boundary).
func TestReadChangesLeereTrefferlisteBleibtGesetzteListe(t *testing.T) {
	fake := &fakeReadChanges{result: inbound.ReadChangesResult{}}
	svc := &administrationService{readChanges: fake, log: outbound.NoopLog}
	resp, err := svc.ReadChanges(context.Background(), &administrationv1.ReadChangesRequest{Source: "src-1"})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if resp.GetChanges() == nil {
		t.Fatalf("leere Trefferliste ist nil statt einer leeren Liste: %+v", resp)
	}
}

// TestDiagnoseRuftUseCaseMitDerQuelleAufUndUebersetztDenBericht trägt die
// Bindung an die Eingabeseite (`ADR-0132` Teilfrage 5): die Quelle erreicht
// den Use Case unverändert, das Ergebnis erreicht unverändert die
// Response — die drei Präsenz-Flags (`Known`/`Present`/`EstimatedRowsKnown`)
// tragen die drei Abwesenheits-Fälle.
func TestDiagnoseRuftUseCaseMitDerQuelleAufUndUebersetztDenBericht(t *testing.T) {
	age := 1.5
	errorClass := "schema"
	lag := 3.0
	backlog := int64(7)
	estimated := int64(10)
	fake := &fakeDiagnose{result: inbound.DiagnoseResult{
		HeartbeatAgeSeconds: &age,
		ErrorClass:          &errorClass,
		CaptureLag:          2.5,
		ConsumerLags:        []inbound.ConsumerLag{{ConsumerID: "c-1", Lag: &lag}},
		RetentionBlocker: &inbound.RetentionBlocker{
			ConsumerID: "c-1", Name: "Consumer 1", AcknowledgedPosition: 10, Backlog: &backlog,
		},
		StorageBytes: 4096,
		Backfill: []inbound.BackfillTableStatus{{
			Schema: "public", Table: "orders", Status: "completed", RowsCopied: 5,
			EstimatedRows: &estimated, WarnEstimatedSize: true, WarnDuration: false,
		}},
	}}
	svc := &administrationService{diagnose: fake, log: outbound.NoopLog}

	resp, err := svc.Diagnose(context.Background(), &administrationv1.DiagnoseRequest{Source: "src-1"})
	if err != nil {
		t.Fatalf("Diagnose: %v", err)
	}
	if !fake.called || fake.gotQuery.Source != "src-1" {
		t.Fatalf("Use Case erhielt nicht die Request-Felder: %+v (called=%v)", fake.gotQuery, fake.called)
	}
	if !resp.GetHeartbeat().GetKnown() || resp.GetHeartbeat().GetAgeSeconds() != age || resp.GetHeartbeat().GetErrorClass() != errorClass {
		t.Fatalf("Heartbeat = %+v", resp.GetHeartbeat())
	}
	if resp.GetCaptureLag() != 2.5 {
		t.Fatalf("CaptureLag = %v, wollen 2.5", resp.GetCaptureLag())
	}
	if len(resp.GetConsumerLags()) != 1 || !resp.GetConsumerLags()[0].GetKnown() || resp.GetConsumerLags()[0].GetLag() != lag {
		t.Fatalf("ConsumerLags = %+v", resp.GetConsumerLags())
	}
	if !resp.GetRetentionBlocker().GetPresent() || resp.GetRetentionBlocker().GetConsumerId() != "c-1" ||
		!resp.GetRetentionBlocker().GetBacklogKnown() || resp.GetRetentionBlocker().GetBacklog() != backlog {
		t.Fatalf("RetentionBlocker = %+v", resp.GetRetentionBlocker())
	}
	if resp.GetStorageBytes() != 4096 {
		t.Fatalf("StorageBytes = %v, wollen 4096", resp.GetStorageBytes())
	}
	if len(resp.GetBackfill()) != 1 || !resp.GetBackfill()[0].GetEstimatedRowsKnown() || resp.GetBackfill()[0].GetEstimatedRows() != estimated {
		t.Fatalf("Backfill = %+v", resp.GetBackfill())
	}
}

// TestDiagnoseUebersetztAbwesenheitenAlsUnbekanntFlags trägt die drei
// Abwesenheits-Fälle (`ADR-0132` Teilfrage 5): kein Lebenszeichen, kein
// Blocker, unbekannte Schätzung/Rückstand — jeweils über das
// Präsenz-Flag, nicht über einen Sentinel-Wert wie `0`.
//
// Rot färbende Mutation: in `toHeartbeatStatus` den frühen `nil`-Zweig
// entfernen und stattdessen `AgeSeconds: 0` zurückgeben — `Known` bliebe
// `true` trotz fehlendem Lebenszeichen, dieser Test färbt rot.
func TestDiagnoseUebersetztAbwesenheitenAlsUnbekanntFlags(t *testing.T) {
	fake := &fakeDiagnose{result: inbound.DiagnoseResult{
		ConsumerLags: []inbound.ConsumerLag{{ConsumerID: "c-1", Lag: nil}},
		Backfill:     []inbound.BackfillTableStatus{{Schema: "public", Table: "orders", EstimatedRows: nil}},
	}}
	svc := &administrationService{diagnose: fake, log: outbound.NoopLog}

	resp, err := svc.Diagnose(context.Background(), &administrationv1.DiagnoseRequest{Source: "src-1"})
	if err != nil {
		t.Fatalf("Diagnose: %v", err)
	}
	if resp.GetHeartbeat().GetKnown() {
		t.Fatalf("Heartbeat.Known = true, wollen false (kein Lebenszeichen)")
	}
	if resp.GetRetentionBlocker().GetPresent() {
		t.Fatalf("RetentionBlocker.Present = true, wollen false (kein Blocker)")
	}
	if resp.GetConsumerLags()[0].GetKnown() {
		t.Fatalf("ConsumerLags[0].Known = true, wollen false (unbekannt)")
	}
	if resp.GetBackfill()[0].GetEstimatedRowsKnown() {
		t.Fatalf("Backfill[0].EstimatedRowsKnown = true, wollen false (unbekannt)")
	}
}

// TestDiagnoseBildetUnbekanntenFehlerAufInternalAb trägt die Fehlerform
// (`ADR-0132` Teilfrage 6): ein Fehler der Klasse `storage`
// (`outbound.ErrDiagnosticsStorage`) endet als `codes.Internal`.
func TestDiagnoseBildetUnbekanntenFehlerAufInternalAb(t *testing.T) {
	fake := &fakeDiagnose{err: outbound.ErrDiagnosticsStorage}
	svc := &administrationService{diagnose: fake, log: outbound.NoopLog}
	_, err := svc.Diagnose(context.Background(), &administrationv1.DiagnoseRequest{Source: "src-1"})
	if status.Code(err) != codes.Internal {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.Internal)
	}
}
