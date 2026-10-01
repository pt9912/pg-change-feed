// Diese Tests prüfen den erzeugten Protokoll-Stub des `Administration`-
// Service an den Stellen, die eine Zusage des Draht-Vertrags (`ADR-0130`)
// tragen: die Feld-Getter, die Fehler-Weitergabe des
// erzeugten unären Clients, die Vorwärtskompatibilitäts-Stubs
// (`UnimplementedAdministrationServer`) und die Handler-Verklebung des
// Servers (Dekodier-Fehler, Aufruf ohne und mit Interceptor). Der erzeugte
// Code ist nicht von Hand änderbar — `make generated-sync` hält ihn
// byte-gleich zur Ausgabe des gepinnten Generators; geprüft wird deshalb
// sein Verhalten, nicht sein Text — dieselbe Testgrenze wie
// `gen/cdc/stream/v1/changestream_test.go`.
//
// Nicht Gegenstand sind die Proto-Runtime-Interna (`String`, `ProtoMessage`,
// `ProtoReflect`, `Descriptor`, die Deskriptor-Kompression und der
// Wiedereintritt von `init`). Sie sind aufrufbar, tragen aber keine Aussage
// über die Administration-Fähigkeiten: sie beschreiben den erzeugten Code
// gegen sich selbst (`ADR-0082` §Konsequenzen).
package administrationv1_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
)

// TestConsumerNachrichtenGetterTragenDenNullwertUndDenGesetztenWert trägt
// die Getter-Grenze der vier Consumer-Nachrichten (Request und Response von
// RegisterConsumer, AcknowledgeConsumer, GetConsumerPosition,
// RemoveConsumer): auf dem Nullwert-Zeiger antworten alle Getter mit ihrem
// Nullwert, statt zu dereferenzieren; auf einer gesetzten Nachricht tragen
// sie den gesetzten Wert.
// Rot färbende Mutation: in einem Getter den `x != nil`-Zweig fallenlassen
// (Dereferenzierung des Nullwerts) oder einen anderen Nullwert zurückgeben.
func TestConsumerNachrichtenGetterTragenDenNullwertUndDenGesetztenWert(t *testing.T) {
	var registerReq *administrationv1.RegisterConsumerRequest
	if registerReq.GetConsumerId() != "" || registerReq.GetName() != "" {
		t.Fatalf("RegisterConsumerRequest auf dem Nullwert liefert keine leeren Felder")
	}
	gesetzterRegisterReq := &administrationv1.RegisterConsumerRequest{ConsumerId: "c-1", Name: "n"}
	if gesetzterRegisterReq.GetConsumerId() != "c-1" || gesetzterRegisterReq.GetName() != "n" {
		t.Fatalf("RegisterConsumerRequest trägt nicht die gesetzten Felder")
	}

	var registerResp *administrationv1.RegisterConsumerResponse
	if registerResp.GetConsumerId() != "" || registerResp.GetName() != "" || registerResp.GetAlreadyRegistered() {
		t.Fatalf("RegisterConsumerResponse auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterRegisterResp := &administrationv1.RegisterConsumerResponse{ConsumerId: "c-1", Name: "n", AlreadyRegistered: true}
	if gesetzterRegisterResp.GetConsumerId() != "c-1" || gesetzterRegisterResp.GetName() != "n" || !gesetzterRegisterResp.GetAlreadyRegistered() {
		t.Fatalf("RegisterConsumerResponse trägt nicht die gesetzten Felder")
	}

	var ackReq *administrationv1.AcknowledgeConsumerRequest
	if ackReq.GetConsumerId() != "" || ackReq.GetSourceId() != "" || ackReq.GetOffset() != 0 {
		t.Fatalf("AcknowledgeConsumerRequest auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterAckReq := &administrationv1.AcknowledgeConsumerRequest{ConsumerId: "c-1", SourceId: "s-1", Offset: 7}
	if gesetzterAckReq.GetConsumerId() != "c-1" || gesetzterAckReq.GetSourceId() != "s-1" || gesetzterAckReq.GetOffset() != 7 {
		t.Fatalf("AcknowledgeConsumerRequest trägt nicht die gesetzten Felder")
	}

	var ackResp *administrationv1.AcknowledgeConsumerResponse
	if ackResp.GetConsumerId() != "" || ackResp.GetSourceId() != "" || ackResp.GetOffset() != 0 {
		t.Fatalf("AcknowledgeConsumerResponse auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterAckResp := &administrationv1.AcknowledgeConsumerResponse{ConsumerId: "c-1", SourceId: "s-1", Offset: 7}
	if gesetzterAckResp.GetConsumerId() != "c-1" || gesetzterAckResp.GetSourceId() != "s-1" || gesetzterAckResp.GetOffset() != 7 {
		t.Fatalf("AcknowledgeConsumerResponse trägt nicht die gesetzten Felder")
	}

	var positionReq *administrationv1.GetConsumerPositionRequest
	if positionReq.GetConsumerId() != "" {
		t.Fatalf("GetConsumerPositionRequest auf dem Nullwert liefert kein leeres Feld")
	}
	gesetzterPositionReq := &administrationv1.GetConsumerPositionRequest{ConsumerId: "c-1"}
	if gesetzterPositionReq.GetConsumerId() != "c-1" {
		t.Fatalf("GetConsumerPositionRequest trägt nicht das gesetzte Feld")
	}

	var positionResp *administrationv1.GetConsumerPositionResponse
	if positionResp.GetConsumerId() != "" || positionResp.GetSourceId() != "" || positionResp.GetOffset() != 0 || positionResp.GetAcknowledged() {
		t.Fatalf("GetConsumerPositionResponse auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterPositionResp := &administrationv1.GetConsumerPositionResponse{ConsumerId: "c-1", SourceId: "s-1", Offset: 7, Acknowledged: true}
	if gesetzterPositionResp.GetConsumerId() != "c-1" || gesetzterPositionResp.GetSourceId() != "s-1" ||
		gesetzterPositionResp.GetOffset() != 7 || !gesetzterPositionResp.GetAcknowledged() {
		t.Fatalf("GetConsumerPositionResponse trägt nicht die gesetzten Felder")
	}

	var removeReq *administrationv1.RemoveConsumerRequest
	if removeReq.GetConsumerId() != "" {
		t.Fatalf("RemoveConsumerRequest auf dem Nullwert liefert kein leeres Feld")
	}
	gesetzterRemoveReq := &administrationv1.RemoveConsumerRequest{ConsumerId: "c-1"}
	if gesetzterRemoveReq.GetConsumerId() != "c-1" {
		t.Fatalf("RemoveConsumerRequest trägt nicht das gesetzte Feld")
	}

	var removeResp *administrationv1.RemoveConsumerResponse
	if removeResp.GetConsumerId() != "" || removeResp.GetRemoved() {
		t.Fatalf("RemoveConsumerResponse auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterRemoveResp := &administrationv1.RemoveConsumerResponse{ConsumerId: "c-1", Removed: true}
	if gesetzterRemoveResp.GetConsumerId() != "c-1" || !gesetzterRemoveResp.GetRemoved() {
		t.Fatalf("RemoveConsumerResponse trägt nicht die gesetzten Felder")
	}
}

// TestTabellenNachrichtenGetterTragenDenNullwertUndDenGesetztenWert trägt
// dieselbe Grenze für die Tabellen-Verwaltungs-Nachrichten (EnableTable,
// DisableTable, GetTableStatus, ListTables samt `SourceTable`) und
// RunRetention.
func TestTabellenNachrichtenGetterTragenDenNullwertUndDenGesetztenWert(t *testing.T) {
	var enableReq *administrationv1.EnableTableRequest
	if enableReq.GetSource() != "" || enableReq.GetSchema() != "" || enableReq.GetTable() != "" ||
		enableReq.GetTableId() != "" || enableReq.GetSchemaVersionId() != "" || enableReq.GetVersion() != 0 ||
		enableReq.GetPublication() != "" {
		t.Fatalf("EnableTableRequest auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterEnableReq := &administrationv1.EnableTableRequest{
		Source: "src-1", Schema: "public", Table: "orders", TableId: "t-1",
		SchemaVersionId: "t-1-v1", Version: 3, Publication: "pub",
	}
	if gesetzterEnableReq.GetSource() != "src-1" || gesetzterEnableReq.GetSchema() != "public" ||
		gesetzterEnableReq.GetTable() != "orders" || gesetzterEnableReq.GetTableId() != "t-1" ||
		gesetzterEnableReq.GetSchemaVersionId() != "t-1-v1" || gesetzterEnableReq.GetVersion() != 3 ||
		gesetzterEnableReq.GetPublication() != "pub" {
		t.Fatalf("EnableTableRequest trägt nicht die gesetzten Felder")
	}

	var enableResp *administrationv1.EnableTableResponse
	if enableResp.GetTableId() != "" || enableResp.GetSource() != "" || enableResp.GetSchema() != "" ||
		enableResp.GetTable() != "" || enableResp.GetAlreadyEnabled() {
		t.Fatalf("EnableTableResponse auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterEnableResp := &administrationv1.EnableTableResponse{
		TableId: "t-1", Source: "src-1", Schema: "public", Table: "orders", AlreadyEnabled: true,
	}
	if gesetzterEnableResp.GetTableId() != "t-1" || gesetzterEnableResp.GetSource() != "src-1" ||
		gesetzterEnableResp.GetSchema() != "public" || gesetzterEnableResp.GetTable() != "orders" ||
		!gesetzterEnableResp.GetAlreadyEnabled() {
		t.Fatalf("EnableTableResponse trägt nicht die gesetzten Felder")
	}

	var disableReq *administrationv1.DisableTableRequest
	if disableReq.GetSource() != "" || disableReq.GetSchema() != "" || disableReq.GetTable() != "" || disableReq.GetPublication() != "" {
		t.Fatalf("DisableTableRequest auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterDisableReq := &administrationv1.DisableTableRequest{Source: "src-1", Schema: "public", Table: "orders", Publication: "pub"}
	if gesetzterDisableReq.GetSource() != "src-1" || gesetzterDisableReq.GetSchema() != "public" ||
		gesetzterDisableReq.GetTable() != "orders" || gesetzterDisableReq.GetPublication() != "pub" {
		t.Fatalf("DisableTableRequest trägt nicht die gesetzten Felder")
	}

	var disableResp *administrationv1.DisableTableResponse
	if disableResp.GetRemoved() || disableResp.GetRetained() {
		t.Fatalf("DisableTableResponse auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterDisableResp := &administrationv1.DisableTableResponse{Removed: true, Retained: true}
	if !gesetzterDisableResp.GetRemoved() || !gesetzterDisableResp.GetRetained() {
		t.Fatalf("DisableTableResponse trägt nicht die gesetzten Felder")
	}

	var statusReq *administrationv1.GetTableStatusRequest
	if statusReq.GetSource() != "" || statusReq.GetSchema() != "" || statusReq.GetTable() != "" || statusReq.GetPublication() != "" {
		t.Fatalf("GetTableStatusRequest auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterStatusReq := &administrationv1.GetTableStatusRequest{Source: "src-1", Schema: "public", Table: "orders", Publication: "pub"}
	if gesetzterStatusReq.GetSource() != "src-1" || gesetzterStatusReq.GetSchema() != "public" ||
		gesetzterStatusReq.GetTable() != "orders" || gesetzterStatusReq.GetPublication() != "pub" {
		t.Fatalf("GetTableStatusRequest trägt nicht die gesetzten Felder")
	}

	var statusResp *administrationv1.GetTableStatusResponse
	if statusResp.GetEnabled() || statusResp.GetRetained() {
		t.Fatalf("GetTableStatusResponse auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterStatusResp := &administrationv1.GetTableStatusResponse{Enabled: true, Retained: true}
	if !gesetzterStatusResp.GetEnabled() || !gesetzterStatusResp.GetRetained() {
		t.Fatalf("GetTableStatusResponse trägt nicht die gesetzten Felder")
	}

	var listReq *administrationv1.ListTablesRequest
	if listReq.GetSource() != "" || listReq.GetPublication() != "" {
		t.Fatalf("ListTablesRequest auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterListReq := &administrationv1.ListTablesRequest{Source: "src-1", Publication: "pub"}
	if gesetzterListReq.GetSource() != "src-1" || gesetzterListReq.GetPublication() != "pub" {
		t.Fatalf("ListTablesRequest trägt nicht die gesetzten Felder")
	}

	var listResp *administrationv1.ListTablesResponse
	if listResp.GetTables() != nil || listResp.GetRetained() != nil {
		t.Fatalf("ListTablesResponse auf dem Nullwert liefert kein nil")
	}
	var sourceTable *administrationv1.SourceTable
	if sourceTable.GetTableId() != "" || sourceTable.GetSource() != "" || sourceTable.GetSchema() != "" || sourceTable.GetTable() != "" {
		t.Fatalf("SourceTable auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzteSourceTable := &administrationv1.SourceTable{TableId: "t-1", Source: "src-1", Schema: "public", Table: "orders"}
	if gesetzteSourceTable.GetTableId() != "t-1" || gesetzteSourceTable.GetSource() != "src-1" ||
		gesetzteSourceTable.GetSchema() != "public" || gesetzteSourceTable.GetTable() != "orders" {
		t.Fatalf("SourceTable trägt nicht die gesetzten Felder")
	}
	gesetzterListResp := &administrationv1.ListTablesResponse{
		Tables:   []*administrationv1.SourceTable{gesetzteSourceTable},
		Retained: []*administrationv1.SourceTable{gesetzteSourceTable},
	}
	if len(gesetzterListResp.GetTables()) != 1 || len(gesetzterListResp.GetRetained()) != 1 {
		t.Fatalf("ListTablesResponse trägt nicht die gesetzten Listen")
	}

	var retentionReq *administrationv1.RunRetentionRequest
	if retentionReq.GetSource() != "" || retentionReq.GetMinAgeNanos() != 0 {
		t.Fatalf("RunRetentionRequest auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterRetentionReq := &administrationv1.RunRetentionRequest{Source: "src-1", MinAgeNanos: 1000}
	if gesetzterRetentionReq.GetSource() != "src-1" || gesetzterRetentionReq.GetMinAgeNanos() != 1000 {
		t.Fatalf("RunRetentionRequest trägt nicht die gesetzten Felder")
	}

	var retentionResp *administrationv1.RunRetentionResponse
	if retentionResp.GetDeleted() != 0 {
		t.Fatalf("RunRetentionResponse auf dem Nullwert liefert kein 0")
	}
	gesetzterRetentionResp := &administrationv1.RunRetentionResponse{Deleted: 5}
	if gesetzterRetentionResp.GetDeleted() != 5 {
		t.Fatalf("RunRetentionResponse trägt nicht das gesetzte Feld")
	}
}

// TestReadChangesRequestTraegtTargetAlsFeldSieben trägt den Draht-Vertrag
// des Feldes `target` (`ADR-0138` Festlegung 1, Feldnummer 7): auf dem Draht
// steht der Wert unter dem Tag von Feld 7 (`0x3a`, Wire-Typ 2) hinter
// `source` (Feld 1) und `limit` (Feld 6); Bytes eines Clients ohne das Feld
// lesen mit leerem `target`, Bytes mit einem dem Empfänger unbekannten Feld
// (hier 9) lesen ohne Fehler — die Messung der eingesetzten
// Protobuf-Bibliothek, nicht eines ausgelieferten Altservers.
// Rot färbende Mutation: in der `.proto` die Feldnummer von `target` ändern.
func TestReadChangesRequestTraegtTargetAlsFeldSieben(t *testing.T) {
	gesetzt := &administrationv1.ReadChangesRequest{Source: "s", Limit: 5, Target: "eu"}
	draht, err := proto.Marshal(gesetzt)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	wollen := []byte{0x0a, 0x01, 's', 0x30, 0x05, 0x3a, 0x02, 'e', 'u'}
	if !bytes.Equal(draht, wollen) {
		t.Fatalf("Draht-Bytes = %x (Erwartung: %x)", draht, wollen)
	}

	ohneFeld := &administrationv1.ReadChangesRequest{}
	if err := proto.Unmarshal([]byte{0x0a, 0x01, 's', 0x30, 0x05}, ohneFeld); err != nil {
		t.Fatalf("Unmarshal ohne Feld 7: %v", err)
	}
	if ohneFeld.GetSource() != "s" || ohneFeld.GetLimit() != 5 || ohneFeld.GetTarget() != "" {
		t.Fatalf("Request ohne Feld 7 = %+v (Erwartung: target leer)", ohneFeld)
	}

	mitUnbekanntem := &administrationv1.ReadChangesRequest{}
	if err := proto.Unmarshal([]byte{0x0a, 0x01, 's', 0x4a, 0x01, 'x'}, mitUnbekanntem); err != nil {
		t.Fatalf("Unmarshal mit unbekanntem Feld: %v", err)
	}
	if mitUnbekanntem.GetSource() != "s" || mitUnbekanntem.GetTarget() != "" {
		t.Fatalf("Request mit unbekanntem Feld = %+v", mitUnbekanntem)
	}
}

// TestReadChangesNachrichtenGetterTragenDenNullwertUndDenGesetztenWert trägt
// dieselbe Grenze für den zehnten RPC `ReadChanges` (`ADR-0131`): Request,
// `ChangeRecord` und Response.
func TestReadChangesNachrichtenGetterTragenDenNullwertUndDenGesetztenWert(t *testing.T) {
	var req *administrationv1.ReadChangesRequest
	if req.GetSource() != "" || req.GetSchema() != "" || req.GetTable() != "" ||
		req.GetFrom() != 0 || req.GetTo() != 0 || req.GetLimit() != 0 || req.GetTarget() != "" {
		t.Fatalf("ReadChangesRequest auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterReq := &administrationv1.ReadChangesRequest{
		Source: "src-1", Schema: "public", Table: "orders", From: 1, To: 10, Limit: 5, Target: "eu",
	}
	if gesetzterReq.GetSource() != "src-1" || gesetzterReq.GetSchema() != "public" ||
		gesetzterReq.GetTable() != "orders" || gesetzterReq.GetFrom() != 1 ||
		gesetzterReq.GetTo() != 10 || gesetzterReq.GetLimit() != 5 || gesetzterReq.GetTarget() != "eu" {
		t.Fatalf("ReadChangesRequest trägt nicht die gesetzten Felder")
	}

	var record *administrationv1.ChangeRecord
	if record.GetCommitPosition() != 0 || record.GetChangeId() != "" || record.GetTransactionId() != "" ||
		record.GetSourceTableId() != "" || record.GetSchema() != "" || record.GetTable() != "" ||
		record.GetSequence() != 0 || record.GetOperation() != "" || record.GetOldImage() != nil ||
		record.GetNewImage() != nil || record.GetSchemaVersion() != "" || record.GetCommittedAt() != "" ||
		record.GetOrigin() != "" {
		t.Fatalf("ChangeRecord auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterRecord := &administrationv1.ChangeRecord{
		CommitPosition: 1, ChangeId: "c-1", TransactionId: "t-1", SourceTableId: "st-1",
		Schema: "public", Table: "orders", Sequence: 1, Operation: "INSERT",
		OldImage: []byte(`{}`), NewImage: []byte(`{"a":1}`), SchemaVersion: "sv-1",
		CommittedAt: "2026-09-28T00:00:00Z", Origin: "wal",
	}
	if gesetzterRecord.GetCommitPosition() != 1 || gesetzterRecord.GetChangeId() != "c-1" ||
		gesetzterRecord.GetTransactionId() != "t-1" || gesetzterRecord.GetSourceTableId() != "st-1" ||
		gesetzterRecord.GetSchema() != "public" || gesetzterRecord.GetTable() != "orders" ||
		gesetzterRecord.GetSequence() != 1 || gesetzterRecord.GetOperation() != "INSERT" ||
		string(gesetzterRecord.GetOldImage()) != "{}" || string(gesetzterRecord.GetNewImage()) != `{"a":1}` ||
		gesetzterRecord.GetSchemaVersion() != "sv-1" || gesetzterRecord.GetCommittedAt() != "2026-09-28T00:00:00Z" ||
		gesetzterRecord.GetOrigin() != "wal" {
		t.Fatalf("ChangeRecord trägt nicht die gesetzten Felder")
	}

	var resp *administrationv1.ReadChangesResponse
	if resp.GetChanges() != nil {
		t.Fatalf("ReadChangesResponse auf dem Nullwert liefert kein nil")
	}
	gesetzteResp := &administrationv1.ReadChangesResponse{Changes: []*administrationv1.ChangeRecord{gesetzterRecord}}
	if len(gesetzteResp.GetChanges()) != 1 {
		t.Fatalf("ReadChangesResponse trägt nicht die gesetzte Liste")
	}
}

// TestDiagnoseNachrichtenGetterTragenDenNullwertUndDenGesetztenWert trägt
// dieselbe Grenze für den elften RPC `Diagnose` (`ADR-0132`): Request,
// die vier Hilfsnachrichten und Response.
func TestDiagnoseNachrichtenGetterTragenDenNullwertUndDenGesetztenWert(t *testing.T) {
	var req *administrationv1.DiagnoseRequest
	if req.GetSource() != "" {
		t.Fatalf("DiagnoseRequest auf dem Nullwert liefert keinen Nullwert")
	}
	gesetzterReq := &administrationv1.DiagnoseRequest{Source: "src-1"}
	if gesetzterReq.GetSource() != "src-1" {
		t.Fatalf("DiagnoseRequest trägt nicht das gesetzte Feld")
	}

	var heartbeat *administrationv1.HeartbeatStatus
	if heartbeat.GetKnown() || heartbeat.GetAgeSeconds() != 0 || heartbeat.GetErrorClass() != "" {
		t.Fatalf("HeartbeatStatus auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterHeartbeat := &administrationv1.HeartbeatStatus{Known: true, AgeSeconds: 1.5, ErrorClass: "schema"}
	if !gesetzterHeartbeat.GetKnown() || gesetzterHeartbeat.GetAgeSeconds() != 1.5 || gesetzterHeartbeat.GetErrorClass() != "schema" {
		t.Fatalf("HeartbeatStatus trägt nicht die gesetzten Felder")
	}

	var lag *administrationv1.ConsumerLag
	if lag.GetConsumerId() != "" || lag.GetKnown() || lag.GetLag() != 0 {
		t.Fatalf("ConsumerLag auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterLag := &administrationv1.ConsumerLag{ConsumerId: "c-1", Known: true, Lag: 3}
	if gesetzterLag.GetConsumerId() != "c-1" || !gesetzterLag.GetKnown() || gesetzterLag.GetLag() != 3 {
		t.Fatalf("ConsumerLag trägt nicht die gesetzten Felder")
	}

	var blocker *administrationv1.RetentionBlocker
	if blocker.GetPresent() || blocker.GetConsumerId() != "" || blocker.GetName() != "" ||
		blocker.GetAcknowledgedPosition() != 0 || blocker.GetBacklogKnown() || blocker.GetBacklog() != 0 {
		t.Fatalf("RetentionBlocker auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzterBlocker := &administrationv1.RetentionBlocker{
		Present: true, ConsumerId: "c-1", Name: "Consumer 1", AcknowledgedPosition: 10,
		BacklogKnown: true, Backlog: 7,
	}
	if !gesetzterBlocker.GetPresent() || gesetzterBlocker.GetConsumerId() != "c-1" ||
		gesetzterBlocker.GetName() != "Consumer 1" || gesetzterBlocker.GetAcknowledgedPosition() != 10 ||
		!gesetzterBlocker.GetBacklogKnown() || gesetzterBlocker.GetBacklog() != 7 {
		t.Fatalf("RetentionBlocker trägt nicht die gesetzten Felder")
	}

	var table *administrationv1.BackfillTableStatus
	if table.GetSchema() != "" || table.GetTable() != "" || table.GetStatus() != "" ||
		table.GetRowsCopied() != 0 || table.GetEstimatedRowsKnown() || table.GetEstimatedRows() != 0 ||
		table.GetWarnEstimatedSize() || table.GetWarnDuration() || table.GetErrorMessage() != "" {
		t.Fatalf("BackfillTableStatus auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzteTable := &administrationv1.BackfillTableStatus{
		Schema: "public", Table: "orders", Status: "completed", RowsCopied: 5,
		EstimatedRowsKnown: true, EstimatedRows: 10, WarnEstimatedSize: true, WarnDuration: false,
		ErrorMessage: "kaputt",
	}
	if gesetzteTable.GetSchema() != "public" || gesetzteTable.GetTable() != "orders" ||
		gesetzteTable.GetStatus() != "completed" || gesetzteTable.GetRowsCopied() != 5 ||
		!gesetzteTable.GetEstimatedRowsKnown() || gesetzteTable.GetEstimatedRows() != 10 ||
		!gesetzteTable.GetWarnEstimatedSize() || gesetzteTable.GetWarnDuration() ||
		gesetzteTable.GetErrorMessage() != "kaputt" {
		t.Fatalf("BackfillTableStatus trägt nicht die gesetzten Felder")
	}

	var resp *administrationv1.DiagnoseResponse
	if resp.GetHeartbeat() != nil || resp.GetCaptureLag() != 0 || resp.GetConsumerLags() != nil ||
		resp.GetRetentionBlocker() != nil || resp.GetStorageBytes() != 0 || resp.GetBackfill() != nil {
		t.Fatalf("DiagnoseResponse auf dem Nullwert liefert keine Nullwerte")
	}
	gesetzteResp := &administrationv1.DiagnoseResponse{
		Heartbeat: gesetzterHeartbeat, CaptureLag: 2.5, ConsumerLags: []*administrationv1.ConsumerLag{gesetzterLag},
		RetentionBlocker: gesetzterBlocker, StorageBytes: 4096, Backfill: []*administrationv1.BackfillTableStatus{gesetzteTable},
	}
	if gesetzteResp.GetHeartbeat() != gesetzterHeartbeat || gesetzteResp.GetCaptureLag() != 2.5 ||
		len(gesetzteResp.GetConsumerLags()) != 1 || gesetzteResp.GetRetentionBlocker() != gesetzterBlocker ||
		gesetzteResp.GetStorageBytes() != 4096 || len(gesetzteResp.GetBackfill()) != 1 {
		t.Fatalf("DiagnoseResponse trägt nicht die gesetzten Felder")
	}
}

// fakeUnaryClientConn trägt einen `grpc.ClientConnInterface`-Doppel für die
// elf unären RPCs: `Invoke` liefert den injizierten Fehler, ohne dass ein
// Transport im Spiel ist — dieselbe Bauform wie `fakeClientConn` in
// `changestream_test.go`, für unäre statt streamende RPCs.
type fakeUnaryClientConn struct{ err error }

func (f fakeUnaryClientConn) Invoke(context.Context, string, any, any, ...grpc.CallOption) error {
	return f.err
}

func (f fakeUnaryClientConn) NewStream(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, errors.New("administrationClient nutzt keinen Stream")
}

var _ grpc.ClientConnInterface = fakeUnaryClientConn{}

// TestAdministrationClientReichtInvokeErgebnisJeRPCWeiter trägt die
// Fehler-Weitergabe und den Erfolgspfad des erzeugten unären Clients für
// alle elf RPCs (`ADR-0131` Teilfrage 2): scheitert `Invoke`, liefert die
// Methode **genau diesen** Fehler und kein Response; gelingt `Invoke`,
// liefert sie ein Response ohne Fehler.
// Rot färbende Mutation: in einer Client-Methode den `err != nil`-Zweig
// fallenlassen oder den Fehler durch einen anderen ersetzen.
func TestAdministrationClientReichtInvokeErgebnisJeRPCWeiter(t *testing.T) {
	injizierterFehler := errors.New("Verbindung weg")
	cases := []struct {
		name string
		call func(client administrationv1.AdministrationClient) (any, error)
	}{
		{"RegisterConsumer", func(c administrationv1.AdministrationClient) (any, error) {
			return c.RegisterConsumer(context.Background(), &administrationv1.RegisterConsumerRequest{})
		}},
		{"AcknowledgeConsumer", func(c administrationv1.AdministrationClient) (any, error) {
			return c.AcknowledgeConsumer(context.Background(), &administrationv1.AcknowledgeConsumerRequest{})
		}},
		{"GetConsumerPosition", func(c administrationv1.AdministrationClient) (any, error) {
			return c.GetConsumerPosition(context.Background(), &administrationv1.GetConsumerPositionRequest{})
		}},
		{"RemoveConsumer", func(c administrationv1.AdministrationClient) (any, error) {
			return c.RemoveConsumer(context.Background(), &administrationv1.RemoveConsumerRequest{})
		}},
		{"EnableTable", func(c administrationv1.AdministrationClient) (any, error) {
			return c.EnableTable(context.Background(), &administrationv1.EnableTableRequest{})
		}},
		{"DisableTable", func(c administrationv1.AdministrationClient) (any, error) {
			return c.DisableTable(context.Background(), &administrationv1.DisableTableRequest{})
		}},
		{"GetTableStatus", func(c administrationv1.AdministrationClient) (any, error) {
			return c.GetTableStatus(context.Background(), &administrationv1.GetTableStatusRequest{})
		}},
		{"ListTables", func(c administrationv1.AdministrationClient) (any, error) {
			return c.ListTables(context.Background(), &administrationv1.ListTablesRequest{})
		}},
		{"RunRetention", func(c administrationv1.AdministrationClient) (any, error) {
			return c.RunRetention(context.Background(), &administrationv1.RunRetentionRequest{})
		}},
		{"ReadChanges", func(c administrationv1.AdministrationClient) (any, error) {
			return c.ReadChanges(context.Background(), &administrationv1.ReadChangesRequest{})
		}},
		{"Diagnose", func(c administrationv1.AdministrationClient) (any, error) {
			return c.Diagnose(context.Background(), &administrationv1.DiagnoseRequest{})
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name+"/Erfolg", func(t *testing.T) {
			client := administrationv1.NewAdministrationClient(fakeUnaryClientConn{})
			resp, err := tc.call(client)
			if err != nil {
				t.Fatalf("%s: %v (Erwartung: kein Fehler)", tc.name, err)
			}
			if resp == nil {
				t.Fatalf("%s: kein Response bei Erfolg", tc.name)
			}
		})
		t.Run(tc.name+"/Fehler", func(t *testing.T) {
			client := administrationv1.NewAdministrationClient(fakeUnaryClientConn{err: injizierterFehler})
			_, err := tc.call(client)
			if !errors.Is(err, injizierterFehler) {
				t.Fatalf("%s: %v (Erwartung: der injizierte Fehler %v)", tc.name, err, injizierterFehler)
			}
		})
	}
}

// TestUnimplementedAdministrationServerEndetMitUnimplemented trägt den
// Ausgang des Vorwärtskompatibilitäts-Stubs für alle elf RPCs: eine
// Implementierung, die eine Methode nicht trägt, endet mit dem gRPC-Status
// `Unimplemented` — sichtbar, nicht mit einem stillen leeren Response.
// Rot färbende Mutation: den Statuscode eines Stubs ändern.
func TestUnimplementedAdministrationServerEndetMitUnimplemented(t *testing.T) {
	var srv administrationv1.UnimplementedAdministrationServer
	ctx := context.Background()

	if _, err := srv.RegisterConsumer(ctx, &administrationv1.RegisterConsumerRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("RegisterConsumer: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.AcknowledgeConsumer(ctx, &administrationv1.AcknowledgeConsumerRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("AcknowledgeConsumer: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.GetConsumerPosition(ctx, &administrationv1.GetConsumerPositionRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("GetConsumerPosition: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.RemoveConsumer(ctx, &administrationv1.RemoveConsumerRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("RemoveConsumer: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.EnableTable(ctx, &administrationv1.EnableTableRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("EnableTable: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.DisableTable(ctx, &administrationv1.DisableTableRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("DisableTable: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.GetTableStatus(ctx, &administrationv1.GetTableStatusRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("GetTableStatus: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.ListTables(ctx, &administrationv1.ListTablesRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("ListTables: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.RunRetention(ctx, &administrationv1.RunRetentionRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("RunRetention: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.ReadChanges(ctx, &administrationv1.ReadChangesRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("ReadChanges: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
	if _, err := srv.Diagnose(ctx, &administrationv1.DiagnoseRequest{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("Diagnose: %v (Erwartung: %v)", status.Code(err), codes.Unimplemented)
	}
}

// fakeAdministrationServer trägt eine In-Memory-Implementierung von
// `AdministrationServer` für die Handler-Verklebungstests unten: jede
// Methode trägt ein eigenes, unterscheidbares Response, ohne einen echten
// Inbound Use Case zu berühren (Whitebox-Test des erzeugten Codes, kein
// Test der Handler aus `administration.go`).
type fakeAdministrationServer struct {
	administrationv1.UnimplementedAdministrationServer
}

func (fakeAdministrationServer) RegisterConsumer(context.Context, *administrationv1.RegisterConsumerRequest) (*administrationv1.RegisterConsumerResponse, error) {
	return &administrationv1.RegisterConsumerResponse{ConsumerId: "RegisterConsumer"}, nil
}
func (fakeAdministrationServer) AcknowledgeConsumer(context.Context, *administrationv1.AcknowledgeConsumerRequest) (*administrationv1.AcknowledgeConsumerResponse, error) {
	return &administrationv1.AcknowledgeConsumerResponse{ConsumerId: "AcknowledgeConsumer"}, nil
}
func (fakeAdministrationServer) GetConsumerPosition(context.Context, *administrationv1.GetConsumerPositionRequest) (*administrationv1.GetConsumerPositionResponse, error) {
	return &administrationv1.GetConsumerPositionResponse{ConsumerId: "GetConsumerPosition"}, nil
}
func (fakeAdministrationServer) RemoveConsumer(context.Context, *administrationv1.RemoveConsumerRequest) (*administrationv1.RemoveConsumerResponse, error) {
	return &administrationv1.RemoveConsumerResponse{ConsumerId: "RemoveConsumer"}, nil
}
func (fakeAdministrationServer) EnableTable(context.Context, *administrationv1.EnableTableRequest) (*administrationv1.EnableTableResponse, error) {
	return &administrationv1.EnableTableResponse{TableId: "EnableTable"}, nil
}
func (fakeAdministrationServer) DisableTable(context.Context, *administrationv1.DisableTableRequest) (*administrationv1.DisableTableResponse, error) {
	return &administrationv1.DisableTableResponse{}, nil
}
func (fakeAdministrationServer) GetTableStatus(context.Context, *administrationv1.GetTableStatusRequest) (*administrationv1.GetTableStatusResponse, error) {
	return &administrationv1.GetTableStatusResponse{}, nil
}
func (fakeAdministrationServer) ListTables(context.Context, *administrationv1.ListTablesRequest) (*administrationv1.ListTablesResponse, error) {
	return &administrationv1.ListTablesResponse{}, nil
}
func (fakeAdministrationServer) RunRetention(context.Context, *administrationv1.RunRetentionRequest) (*administrationv1.RunRetentionResponse, error) {
	return &administrationv1.RunRetentionResponse{}, nil
}
func (fakeAdministrationServer) ReadChanges(context.Context, *administrationv1.ReadChangesRequest) (*administrationv1.ReadChangesResponse, error) {
	return &administrationv1.ReadChangesResponse{}, nil
}
func (fakeAdministrationServer) Diagnose(context.Context, *administrationv1.DiagnoseRequest) (*administrationv1.DiagnoseResponse, error) {
	return &administrationv1.DiagnoseResponse{}, nil
}

var _ administrationv1.AdministrationServer = fakeAdministrationServer{}

// TestAdministrationHandlerVerklebungJeRPC trägt die Handler-Verklebung des
// erzeugten Servers für alle elf RPCs, über den veröffentlichten
// `Administration_ServiceDesc` angesprochen — denselben Wert, den
// `RegisterAdministrationServer` registriert: ein Dekodier-Fehler endet mit
// **genau diesem** Fehler, statt die Service-Methode aufzurufen; ohne
// Interceptor erreicht der Aufruf die Service-Methode direkt; mit einem
// durchreichenden Interceptor erreicht er sie über dessen Handler-Schluss.
// Rot färbende Mutation: den Dekodier-Fehler verwerfen und die
// Service-Methode dennoch aufrufen, oder den Interceptor-Zweig
// überspringen.
func TestAdministrationHandlerVerklebungJeRPC(t *testing.T) {
	dekodierFehler := errors.New("Nachricht nicht lesbar")
	decOK := func(any) error { return nil }
	decErr := func(any) error { return dekodierFehler }
	durchreichenderInterceptor := func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return handler(ctx, req)
	}
	srv := fakeAdministrationServer{}

	for _, method := range administrationv1.Administration_ServiceDesc.Methods {
		method := method
		t.Run(method.MethodName+"/DekodierFehler", func(t *testing.T) {
			_, err := method.Handler(srv, context.Background(), decErr, nil)
			if !errors.Is(err, dekodierFehler) {
				t.Fatalf("%s: %v (Erwartung: der injizierte Dekodier-Fehler %v)", method.MethodName, err, dekodierFehler)
			}
		})
		t.Run(method.MethodName+"/OhneInterceptor", func(t *testing.T) {
			resp, err := method.Handler(srv, context.Background(), decOK, nil)
			if err != nil {
				t.Fatalf("%s ohne Interceptor: %v (Erwartung: kein Fehler)", method.MethodName, err)
			}
			if resp == nil {
				t.Fatalf("%s ohne Interceptor: kein Response", method.MethodName)
			}
		})
		t.Run(method.MethodName+"/MitInterceptor", func(t *testing.T) {
			resp, err := method.Handler(srv, context.Background(), decOK, durchreichenderInterceptor)
			if err != nil {
				t.Fatalf("%s mit Interceptor: %v (Erwartung: kein Fehler)", method.MethodName, err)
			}
			if resp == nil {
				t.Fatalf("%s mit Interceptor: kein Response", method.MethodName)
			}
		})
	}
}
