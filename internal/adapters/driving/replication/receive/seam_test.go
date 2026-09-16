// Die Tests dieses Pakets laufen netzlos: sie fahren die Verklebung des
// Adapters gegen einen Fake, der die Naht erfüllt — die Meldungs-Zerlegung
// der Empfangs-Schleife, die Slot-/Publication-Auflösung, die
// Katalog-Zeilen-Übersetzung und die Rückstands-Messung. Sie sind
// ausdrücklich **kein** Ersatz der realen Tests: dass der reale Treiber die
// Nachrichten liefert, der Slot den bestätigten Stand trägt und der
// Rückstand real wächst, prüfen weiterhin die PostgreSQL-gestützten Läufe
// (`stream_test.go`, `make test-replication`, `ADR-0030`).

package receive

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// fakeSession erfüllt dieselbe Naht wie die Treiber-Hülle und hält fest,
// was der Adapter ihm übergibt — die Verklebung, nicht das Protokoll.
// `messages` wird in Reihenfolge abgegeben; ist sie leer, endet
// `ReceiveMessage` mit `receiveErr` (bzw. einem sichtbaren Ersatz).
type fakeSession struct {
	identify     pglogrepl.IdentifySystemResult
	identifyErr  error
	createResult pglogrepl.CreateReplicationSlotResult
	createErr    error
	createCalls  int
	createSlot   string
	createPlugin string
	startErr     error
	startCalls   int
	startSlot    string
	startLSN     pglogrepl.LSN
	sent         []pglogrepl.StandbyStatusUpdate
	sendErr      error
	results      []*pgconn.Result
	execErr      error
	execSQL      []string
	messages     []pgproto3.BackendMessage
	receiveErr   error
	closeCalls   int
}

func (f *fakeSession) IdentifySystem(context.Context) (pglogrepl.IdentifySystemResult, error) {
	return f.identify, f.identifyErr
}

func (f *fakeSession) CreateReplicationSlot(_ context.Context, slot, outputPlugin string, _ pglogrepl.CreateReplicationSlotOptions) (pglogrepl.CreateReplicationSlotResult, error) {
	f.createCalls++
	f.createSlot = slot
	f.createPlugin = outputPlugin
	return f.createResult, f.createErr
}

func (f *fakeSession) StartReplication(_ context.Context, slot string, startLSN pglogrepl.LSN, _ pglogrepl.StartReplicationOptions) error {
	f.startCalls++
	f.startSlot = slot
	f.startLSN = startLSN
	return f.startErr
}

func (f *fakeSession) SendStandbyStatusUpdate(_ context.Context, ssu pglogrepl.StandbyStatusUpdate) error {
	if f.sendErr != nil {
		return f.sendErr
	}
	f.sent = append(f.sent, ssu)
	return nil
}

func (f *fakeSession) Exec(_ context.Context, sql string) ([]*pgconn.Result, error) {
	f.execSQL = append(f.execSQL, sql)
	if f.execErr != nil {
		return nil, f.execErr
	}
	return f.results, nil
}

func (f *fakeSession) ReceiveMessage(context.Context) (pgproto3.BackendMessage, error) {
	if len(f.messages) == 0 {
		if f.receiveErr != nil {
			return nil, f.receiveErr
		}
		return nil, stderrors.New("keine Nachricht mehr")
	}
	message := f.messages[0]
	f.messages = f.messages[1:]
	return message, nil
}

func (f *fakeSession) Close(context.Context) error {
	f.closeCalls++
	return nil
}

var _ driverSession = (*fakeSession)(nil)

// noopCapture trägt den `CaptureInboundPort` der Empfangs-Schleifen-Tests:
// die Verklebung endet vor dem Capture-Pfad, sobald die Meldung kein
// dekodierbarer `pgoutput`-Payload ist.
type noopCapture struct{}

func (noopCapture) Capture(context.Context, inbound.CaptureCommand) (inbound.CaptureResult, error) {
	return inbound.CaptureResult{}, nil
}

// resultWithRow trägt ein Katalogergebnis mit einer Zeile aus Textwerten.
func resultWithRow(values ...string) []*pgconn.Result {
	row := make([][]byte, len(values))
	for i, value := range values {
		row[i] = []byte(value)
	}
	return []*pgconn.Result{{Rows: [][][]byte{row}}}
}

// emptyResult trägt ein Katalogergebnis ohne Zeile.
func emptyResult() []*pgconn.Result {
	return []*pgconn.Result{{}}
}

// keepaliveData baut den CopyData-Payload einer Primary-Keepalive-Nachricht
// (Byte-ID plus die 17 Bytes des `pglogrepl`-Kopfsatzes).
func keepaliveData(replyRequested bool) []byte {
	data := make([]byte, 18)
	data[0] = pglogrepl.PrimaryKeepaliveMessageByteID
	if replyRequested {
		data[17] = 1
	}
	return data
}

// xlogData baut den CopyData-Payload einer XLogData-Nachricht (Byte-ID,
// der 24-byte-Kopfsatz, dann der `pgoutput`-Payload).
func xlogData(payload []byte) []byte {
	data := make([]byte, 25+len(payload))
	data[0] = pglogrepl.XLogDataByteID
	copy(data[25:], payload)
	return data
}

// newTestStream verdrahtet einen Stream auf den Fake — derselbe
// paket-interne Einstieg, den die netzlosen Tests aller Pfade nutzen.
func newTestStream(session driverSession, capture inbound.CaptureInboundPort) *Stream {
	return newStreamOnSession(session, nil, nil, capture, nil)
}

// TestFirstRowTranslatesCatalogRow trägt die Katalog-Zeilen-Übersetzung:
// die erste Zeile des ersten Ergebnisses wird als Textwerte geliefert.
func TestFirstRowTranslatesCatalogRow(t *testing.T) {
	values, exists := firstRow(resultWithRow("0/16B3748"))
	if !exists {
		t.Fatalf("Zeile nicht gefunden")
	}
	if len(values) != 1 || values[0] != "0/16B3748" {
		t.Fatalf("Werte: %q, erwartet [\"0/16B3748\"]", values)
	}
}

// TestFirstRowWithoutRowReportsAbsence trägt den Leerfall der
// Katalog-Zeilen-Übersetzung: ein Ergebnis ohne Zeile meldet Abwesenheit.
func TestFirstRowWithoutRowReportsAbsence(t *testing.T) {
	for name, results := range map[string][]*pgconn.Result{
		"kein Ergebnis":  nil,
		"leeres Ergebnis": emptyResult(),
	} {
		if _, exists := firstRow(results); exists {
			t.Fatalf("%s: Zeile gemeldet, erwartet Abwesenheit", name)
		}
	}
}

// TestStandbyStatusCarriesPositionInAllThreeLSNs trägt die
// Standby-Status-Form der Keepalive-Antwort: die Meldung trägt die Position
// als Write-, Flush- und Apply-Stand, damit der Slot sie als
// confirmed_flush_lsn trägt (`LH-QA-REL-001.a`).
func TestStandbyStatusCarriesPositionInAllThreeLSNs(t *testing.T) {
	status := standbyStatus(pglogrepl.LSN(4711))
	fields := map[string]pglogrepl.LSN{
		"WALWritePosition": status.WALWritePosition,
		"WALFlushPosition": status.WALFlushPosition,
		"WALApplyPosition": status.WALApplyPosition,
	}
	for name, got := range fields {
		if got != pglogrepl.LSN(4711) {
			t.Fatalf("%s: %d, erwartet 4711", name, got)
		}
	}
}

// TestParseLSNCarriesCatalogValue trägt die LSN-Übersetzung eines
// Katalogwerts in die Positionsform (`ADR-0005`).
func TestParseLSNCarriesCatalogValue(t *testing.T) {
	lsn, err := parseLSN("confirmed_flush_lsn", "0/16B3748")
	if err != nil {
		t.Fatalf("parseLSN: %v", err)
	}
	if lsn != pglogrepl.LSN(0x16B3748) {
		t.Fatalf("LSN: %X, erwartet 16B3748", uint64(lsn))
	}
}

// TestParseLSNRejectsUnreadableValue trägt den Fehlerpfad der
// LSN-Übersetzung: ein unlesbarer Katalogwert endet über die Fehlerklasse
// `replication` (`SPEC-008`) und nennt die Wertquelle.
func TestParseLSNRejectsUnreadableValue(t *testing.T) {
	_, err := parseLSN("confirmed_flush_lsn", "keine-lsn")
	if !stderrors.Is(err, ErrReplication) {
		t.Fatalf("Klasse replication fehlt: %v", err)
	}
	if !strings.Contains(err.Error(), "confirmed_flush_lsn") {
		t.Fatalf("Wertquelle nicht lesbar: %v", err)
	}
}

// TestParseXLogDataCarriesPayload trägt die Meldungs-Zerlegung der
// XLogData-Nachricht: der `pgoutput`-Payload hinter dem Kopfsatz kommt
// unverändert an (`SPEC-010`).
func TestParseXLogDataCarriesPayload(t *testing.T) {
	payload, err := parseXLogData(xlogData([]byte("B-payload")))
	if err != nil {
		t.Fatalf("parseXLogData: %v", err)
	}
	if string(payload) != "B-payload" {
		t.Fatalf("Payload: %q, erwartet \"B-payload\"", payload)
	}
}

// TestParseXLogDataRejectsShortMessage trägt die Grenze der
// Meldungs-Zerlegung: ein zu kurzer Payload endet sichtbar, statt einen
// leeren WAL-Payload vorzutäuschen.
func TestParseXLogDataRejectsShortMessage(t *testing.T) {
	data := make([]byte, 5)
	data[0] = pglogrepl.XLogDataByteID
	if _, err := parseXLogData(data); !stderrors.Is(err, ErrReplication) {
		t.Fatalf("Kurznachricht: %v", err)
	}
}

// TestParseKeepaliveReportsReplyRequested trägt die Meldungs-Zerlegung der
// Keepalive-Nachricht: der Rückgabewert meldet, ob die Quelle eine Antwort
// verlangt.
func TestParseKeepaliveReportsReplyRequested(t *testing.T) {
	for _, want := range []bool{false, true} {
		got, err := parseKeepalive(keepaliveData(want))
		if err != nil {
			t.Fatalf("parseKeepalive(%v): %v", want, err)
		}
		if got != want {
			t.Fatalf("ReplyRequested: %v, erwartet %v", got, want)
		}
	}
}

// TestParseKeepaliveRejectsShortMessage trägt die Grenze der
// Keepalive-Zerlegung.
func TestParseKeepaliveRejectsShortMessage(t *testing.T) {
	if _, err := parseKeepalive([]byte{pglogrepl.PrimaryKeepaliveMessageByteID}); !stderrors.Is(err, ErrReplication) {
		t.Fatalf("Kurznachricht: %v", err)
	}
}

// TestEnsureSlotResumesExistingSlot trägt die Slot-Auflösung des
// bestehenden Slots (`ADR-0012`): der bestätigte Stand
// (`confirmed_flush_lsn`) wird gelesen und als Startposition getragen —
// ohne den Slot neu anzulegen.
func TestEnsureSlotResumesExistingSlot(t *testing.T) {
	session := &fakeSession{results: resultWithRow("0/16B3748")}
	startLSN, err := ensureSlot(context.Background(), outbound.NoopLog, session, "slot_pgc_test")
	if err != nil {
		t.Fatalf("ensureSlot: %v", err)
	}
	if startLSN != pglogrepl.LSN(0x16B3748) {
		t.Fatalf("Startposition: %X, erwartet 16B3748", uint64(startLSN))
	}
	if session.createCalls != 0 {
		t.Fatalf("CreateReplicationSlot-Aufrufe: %d, erwartet 0", session.createCalls)
	}
}

// TestEnsureSlotCreatesMissingSlot trägt die Slot-Auflösung des fehlenden
// Slots (`LH-FA-CFG-001.a`): er wird mit dem `pgoutput`-Plugin angelegt,
// der konsistente Punkt wird zur Startposition.
func TestEnsureSlotCreatesMissingSlot(t *testing.T) {
	session := &fakeSession{
		results:      emptyResult(),
		createResult: pglogrepl.CreateReplicationSlotResult{SlotName: "slot_pgc_test", ConsistentPoint: "0/16B3748"},
	}
	startLSN, err := ensureSlot(context.Background(), outbound.NoopLog, session, "slot_pgc_test")
	if err != nil {
		t.Fatalf("ensureSlot: %v", err)
	}
	if startLSN != pglogrepl.LSN(0x16B3748) {
		t.Fatalf("Startposition: %X, erwartet 16B3748", uint64(startLSN))
	}
	if session.createCalls != 1 || session.createSlot != "slot_pgc_test" || session.createPlugin != outputPlugin {
		t.Fatalf("Slot-Anlage: %d Aufrufe, Slot %q, Plugin %q", session.createCalls, session.createSlot, session.createPlugin)
	}
}

// TestEnsureSlotRejectsEmptyConfirmedFlush trägt die Grenze der
// Slot-Auflösung: ein bestehender Slot ohne `confirmed_flush_lsn` ist kein
// Stand zur Fortsetzung (`SPEC-008`, Klasse `replication`).
func TestEnsureSlotRejectsEmptyConfirmedFlush(t *testing.T) {
	session := &fakeSession{results: resultWithRow("")}
	if _, err := ensureSlot(context.Background(), outbound.NoopLog, session, "slot_pgc_test"); !stderrors.Is(err, ErrReplication) {
		t.Fatalf("leerer Slot-Stand: %v", err)
	}
}

// TestEnsurePublicationAcceptsExistingPublication trägt den Bestand der
// Publication (`LH-FA-CFG-001.a`).
func TestEnsurePublicationAcceptsExistingPublication(t *testing.T) {
	session := &fakeSession{results: resultWithRow("1")}
	if err := ensurePublication(context.Background(), session, "pub_pgc_test"); err != nil {
		t.Fatalf("ensurePublication: %v", err)
	}
}

// TestEnsurePublicationRejectsMissingPublication trägt die
// Publication-Grenze (`LH-FA-CFG-001.a`): der Stream startet nicht ohne
// sie und legt sie nicht still an (Klasse `configuration`).
func TestEnsurePublicationRejectsMissingPublication(t *testing.T) {
	session := &fakeSession{results: emptyResult()}
	if err := ensurePublication(context.Background(), session, "pub_pgc_test"); !stderrors.Is(err, ErrConfiguration) {
		t.Fatalf("fehlende Publication: %v", err)
	}
}

// TestQuerySingleReportsCatalogFailure trägt den Fehlerpfad der
// Katalogabfrage über die Naht: ein Treiber-Fehler endet über die
// Fehlerklasse `replication` (`SPEC-008`).
func TestQuerySingleReportsCatalogFailure(t *testing.T) {
	session := &fakeSession{execErr: stderrors.New("verbindung weg")}
	_, _, err := querySingle(context.Background(), session, slotLSNQuery("slot_pgc_test"))
	if !stderrors.Is(err, ErrReplication) {
		t.Fatalf("Klasse replication fehlt: %v", err)
	}
	if !strings.Contains(err.Error(), "Katalogabfrage") {
		t.Fatalf("Fehlerort nicht lesbar: %v", err)
	}
}

// TestSlotLSNQueryNamesTheSlot trägt die Abfrage des Slot-Stands: sie
// fragt den logischen Slot beim gegebenen Namen ab. Der Fake beantwortet
// jeden SQL-Text gleich — das Prädikat `slot_type` dieser Abfrage deckt
// der Test damit nicht ab.
func TestSlotLSNQueryNamesTheSlot(t *testing.T) {
	query := slotLSNQuery("slot_pgc_test")
	if !strings.Contains(query, "'slot_pgc_test'") || !strings.Contains(query, "confirmed_flush_lsn") {
		t.Fatalf("Abfrage: %q", query)
	}
}

// TestMeasureComputesWALBacklog trägt die Rückstands-Messung (`SPEC-009`
// `cdc_wal_retention_bytes`, `ADR-0049`): die Differenz zwischen der
// aktuellen WAL-Schreibposition der Quelle und `confirmed_flush_lsn` des
// Slots.
func TestMeasureComputesWALBacklog(t *testing.T) {
	session := &fakeSession{
		identify: pglogrepl.IdentifySystemResult{XLogPos: pglogrepl.LSN(2000)},
		results:  resultWithRow("0/400"),
	}
	checker := newWALRetentionCheckerOnSession(session, "", "slot_pgc_test")
	backlog, err := checker.Measure(context.Background())
	if err != nil {
		t.Fatalf("Measure: %v", err)
	}
	if backlog != 976 {
		t.Fatalf("Rückstand: %d, erwartet 976", backlog)
	}
}

// TestMeasureReportsMissingSlot trägt die Grenze der Rückstands-Messung:
// ein Slot ohne `confirmed_flush_lsn` endet sichtbar statt mit einem
// stillen Nullwert (`SPEC-008`).
func TestMeasureReportsMissingSlot(t *testing.T) {
	session := &fakeSession{results: emptyResult()}
	checker := newWALRetentionCheckerOnSession(session, "", "slot_pgc_test")
	if _, err := checker.Measure(context.Background()); !stderrors.Is(err, ErrReplication) {
		t.Fatalf("fehlender Slot: %v", err)
	}
}

// TestMeasureReportsIdentifySystemFailure trägt den Fehlerpfad der
// Rückstands-Messung: ein fehlgeschlagenes `IDENTIFY_SYSTEM` endet über die
// Fehlerklasse `replication` (`SPEC-008`).
func TestMeasureReportsIdentifySystemFailure(t *testing.T) {
	session := &fakeSession{identifyErr: stderrors.New("verbindung weg")}
	checker := newWALRetentionCheckerOnSession(session, "", "slot_pgc_test")
	if _, err := checker.Measure(context.Background()); !stderrors.Is(err, ErrReplication) {
		t.Fatalf("IDENTIFY_SYSTEM-Fehler: %v", err)
	}
}

// TestWALRetentionCheckerCloseClosesTheSession trägt den Rand des
// Checkers: `Close` schließt die Verbindung hinter der Naht.
func TestWALRetentionCheckerCloseClosesTheSession(t *testing.T) {
	session := &fakeSession{}
	checker := newWALRetentionCheckerOnSession(session, "", "slot_pgc_test")
	if err := checker.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if session.closeCalls != 1 {
		t.Fatalf("Close-Aufrufe: %d, erwartet 1", session.closeCalls)
	}
}

// TestRunWithoutCaptureIsConfigurationError trägt die Konstruktions-Grenze
// des Laufs: ohne `CaptureInboundPort` endet er über die
// Konfigurationsklasse, bevor eine Nachricht empfangen wird.
func TestRunWithoutCaptureIsConfigurationError(t *testing.T) {
	session := &fakeSession{}
	stream := newTestStream(session, nil)
	if err := stream.Run(context.Background()); !stderrors.Is(err, ErrConfiguration) {
		t.Fatalf("Lauf ohne Capture-Port: %v", err)
	}
	if session.closeCalls != 0 {
		t.Fatalf("Close-Aufrufe: %d, erwartet 0", session.closeCalls)
	}
}

// TestRunStopsOnCopyDone trägt den regulären Ausgang der Empfangs-Schleife:
// ein `CopyDone` beendet den Lauf ohne Fehler und schließt die Verbindung.
func TestRunStopsOnCopyDone(t *testing.T) {
	session := &fakeSession{messages: []pgproto3.BackendMessage{&pgproto3.CopyDone{}}}
	stream := newTestStream(session, noopCapture{})
	if err := stream.Run(context.Background()); err != nil {
		t.Fatalf("CopyDone: %v", err)
	}
	if session.closeCalls != 1 {
		t.Fatalf("Close-Aufrufe: %d, erwartet 1", session.closeCalls)
	}
}

// TestRunIgnoresOtherBackendMessages trägt den Leerfall der
// Meldungs-Zerlegung: übrige Backend-Nachrichten (NoticeResponse u. a.)
// tragen keine Stream-Änderung, der Lauf setzt fort.
func TestRunIgnoresOtherBackendMessages(t *testing.T) {
	session := &fakeSession{messages: []pgproto3.BackendMessage{
		&pgproto3.NoticeResponse{},
		&pgproto3.CopyDone{},
	}}
	stream := newTestStream(session, noopCapture{})
	if err := stream.Run(context.Background()); err != nil {
		t.Fatalf("NoticeResponse: %v", err)
	}
}

// TestRunAnswersKeepaliveWithAcknowledgedPosition trägt die
// Keepalive-Behandlung der Empfangs-Schleife (`ADR-0007`,
// `LH-QA-REL-001.a`): die Antwort meldet die letzte bestätigte Position —
// hier den Stand des Adapters, nicht den Empfangsstand.
func TestRunAnswersKeepaliveWithAcknowledgedPosition(t *testing.T) {
	session := &fakeSession{messages: []pgproto3.BackendMessage{
		&pgproto3.CopyData{Data: keepaliveData(true)},
		&pgproto3.CopyDone{},
	}}
	stream := newTestStream(session, noopCapture{})
	stream.lastAcked = pglogrepl.LSN(4711)
	if err := stream.Run(context.Background()); err != nil {
		t.Fatalf("Lauf: %v", err)
	}
	if len(session.sent) != 1 {
		t.Fatalf("Keepalive-Antworten: %d, erwartet 1", len(session.sent))
	}
	got := session.sent[0]
	if got.WALWritePosition != pglogrepl.LSN(4711) ||
		got.WALFlushPosition != pglogrepl.LSN(4711) ||
		got.WALApplyPosition != pglogrepl.LSN(4711) {
		t.Fatalf("Meldung trägt %d/%d/%d, erwartet 4711/4711/4711",
			got.WALWritePosition, got.WALFlushPosition, got.WALApplyPosition)
	}
}

// TestRunKeepaliveWithoutReplyRequestedSendsNothing trägt die Grenze der
// Keepalive-Behandlung: ohne `ReplyRequested` setzt der Adapter nichts ab.
func TestRunKeepaliveWithoutReplyRequestedSendsNothing(t *testing.T) {
	session := &fakeSession{messages: []pgproto3.BackendMessage{
		&pgproto3.CopyData{Data: keepaliveData(false)},
		&pgproto3.CopyDone{},
	}}
	stream := newTestStream(session, noopCapture{})
	if err := stream.Run(context.Background()); err != nil {
		t.Fatalf("Lauf: %v", err)
	}
	if len(session.sent) != 0 {
		t.Fatalf("Keepalive-Antworten: %d, erwartet 0", len(session.sent))
	}
}

// TestRunReportsKeepaliveReplyFailure trägt den Fehlerpfad der
// Keepalive-Behandlung: eine fehlgeschlagene Antwort endet über die
// Fehlerklasse `replication` (`SPEC-008`).
func TestRunReportsKeepaliveReplyFailure(t *testing.T) {
	session := &fakeSession{
		sendErr: stderrors.New("verbindung weg"),
		messages: []pgproto3.BackendMessage{
			&pgproto3.CopyData{Data: keepaliveData(true)},
		},
	}
	stream := newTestStream(session, noopCapture{})
	err := stream.Run(context.Background())
	if !stderrors.Is(err, ErrReplication) {
		t.Fatalf("Antwortfehler: %v", err)
	}
	if !strings.Contains(err.Error(), "Keepalive-Antwort") {
		t.Fatalf("Fehlerort nicht lesbar: %v", err)
	}
}

// TestRunRejectsEmptyCopyData trägt die Grenze der Meldungs-Zerlegung: ein
// leeres `CopyData` ist ein sichtbarer Fehler, keine stille Fortsetzung.
func TestRunRejectsEmptyCopyData(t *testing.T) {
	session := &fakeSession{messages: []pgproto3.BackendMessage{
		&pgproto3.CopyData{Data: []byte{}},
	}}
	stream := newTestStream(session, noopCapture{})
	err := stream.Run(context.Background())
	if !stderrors.Is(err, ErrReplication) {
		t.Fatalf("leeres CopyData: %v", err)
	}
	if !strings.Contains(err.Error(), "leeres CopyData") {
		t.Fatalf("Fehlerort nicht lesbar: %v", err)
	}
}

// TestRunDispatchesXLogDataToTheDecoder trägt den XLogData-Zweig der
// Meldungs-Zerlegung: der Payload läuft in den Dekodier-Pfad, ein nicht
// interpretierbarer `pgoutput`-Payload endet als Schema-Klasse.
func TestRunDispatchesXLogDataToTheDecoder(t *testing.T) {
	session := &fakeSession{messages: []pgproto3.BackendMessage{
		&pgproto3.CopyData{Data: xlogData([]byte("kein-pgoutput"))},
	}}
	stream := newTestStream(session, noopCapture{})
	if err := stream.Run(context.Background()); !stderrors.Is(err, decode.ErrSchema) {
		t.Fatalf("XLogData-Zweig: %v", err)
	}
}

// TestRunReportsErrorResponse trägt die Fehlerklasse der Empfangs-Schleife:
// eine Fehlermeldung der Quelle endet über `replication` (`SPEC-008`).
func TestRunReportsErrorResponse(t *testing.T) {
	session := &fakeSession{messages: []pgproto3.BackendMessage{
		&pgproto3.ErrorResponse{Severity: "ERROR", Message: "boom"},
	}}
	stream := newTestStream(session, noopCapture{})
	if err := stream.Run(context.Background()); !stderrors.Is(err, ErrReplication) {
		t.Fatalf("ErrorResponse: %v", err)
	}
}

// TestRunReportsReceiveFailure trägt den Empfangs-Fehlerpfad: eine
// Störung an der Verbindung endet über `replication`, solange der Kontext
// nicht endet (`SPEC-008`).
func TestRunReportsReceiveFailure(t *testing.T) {
	session := &fakeSession{receiveErr: stderrors.New("verbindung weg")}
	stream := newTestStream(session, noopCapture{})
	err := stream.Run(context.Background())
	if !stderrors.Is(err, ErrReplication) {
		t.Fatalf("Empfangsfehler: %v", err)
	}
	if !strings.Contains(err.Error(), "Empfang") {
		t.Fatalf("Fehlerort nicht lesbar: %v", err)
	}
}

// TestRunEndsWithoutErrorWhenContextEnded trägt das reguläre
// Kontext-Ende der Empfangs-Schleife: endet der Kontext, meldet der Lauf
// keinen Fehler — die Störung an der Verbindung ist dann der erwartete
// Ausgang.
func TestRunEndsWithoutErrorWhenContextEnded(t *testing.T) {
	session := &fakeSession{receiveErr: context.Canceled}
	stream := newTestStream(session, noopCapture{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := stream.Run(ctx); err != nil {
		t.Fatalf("Kontext-Ende: %v", err)
	}
}
