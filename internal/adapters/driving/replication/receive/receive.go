// Package receive trägt den Replication-Stream-Adapter als Driving
// Adapter (`ARC-005`, `ADR-0006`): er baut die Replication-Verbindung
// auf, verwaltet den Logical Replication Slot (`LH-FA-CFG-001.a`),
// empfängt die `pgoutput`-Nachrichten des Streams (`ADR-0008`,
// `SPEC-010`) und übersetzt sie über Dekodierung und Mapper in Aufrufe
// des `CaptureInboundPort`. Er entscheidet nicht über Persistenz,
// Retention oder Source-ACK (`ADR-0006` Konsequenz) — die bestätigte
// Position kommt ihm ausschließlich als Ergebnis des Capture-Aufrufs
// entgegen (`LH-QA-REL-001.a`).
package receive

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// outputPlugin trägt den Output-Plugin-Standard (`ADR-0008`).
const outputPlugin = "pgoutput"

// ErrReplication trägt die Fehlerklasse `replication` des
// Stream-Adapters (`SPEC-008`, `ADR-0023`): Störungen an Verbindung und
// Slot enden als sichtbare Fehler — kontrollierte Fortsetzung liegt bei
// dem Aufrufer, der den Stream startet. Der Aufrufer klassifiziert über
// `errors.Is`; die technische Ursache bleibt über die zweite Wrappung
// lesbar.
var ErrReplication = errors.New("Fehlerklasse replication: Replication-Stream/Slot-Störung")

// ErrConfiguration trägt die Fehlerklasse `configuration` des
// Stream-Adapters (`SPEC-008`, `ADR-0023`): eine ungültige oder falsch
// gesetzte Konfiguration endet ohne Start und ohne Fortsetzung im
// falschen Stand.
var ErrConfiguration = errors.New("Fehlerklasse configuration: ungültige/falsch gesetzte Konfiguration")

// identifierShape begrenzt Slot- und Publication-Namen auf das
// Bezeichner-Alphabet der Quelle; beide gehen als Bezeichner-Literal in
// Replication- und Katalogabfragen. Diese Definition trägt die Quelle
// des Alphabet-Vertrags; der Aktivierungs-Adapter
// (`postgresstorage.identifierShape`) hält denselben Ausdruck gegen
// denselben Aufrufgegenstand — die Adapter-Schicht importiert keine
// Adapter-Kante (Kopplung).
var identifierShape = regexp.MustCompile(`^[a-z0-9_]{1,63}$`)

// Config trägt die Konfiguration des Stream-Adapters: die
// Replication-Verbindung (DSN), die Quelle, die Verwaltungsnamen
// Publication und Slot (`LH-FA-CFG-001.a`) und die aktivierten Tabellen
// mit ihren Port-Kennungen. Der Capture-Port wird vor dem Lauf
// verdrahtet (BindCapture) — der ACK-Adapter braucht die Verbindung
// (`ADR-0007`), die `NewStream` erst aufbaut. Die Tabellen-Kennungen
// und die Schema-Versionen liegen initial bei der Konfiguration; die
// dynamische Re-Versionierung trägt `SchemaStore` über den Mapper
// (`mapper.Assembler.Consume`, `ADR-0015` Folgepflicht) — die Erkennung
// nicht sicher interpretierbarer Änderungen als Fehlerklasse `schema`
// bleibt `LH-FA-SCH-004.a`.
type Config struct {
	DSN         string
	Source      model.SourceID
	Publication string
	Slot        string
	Tables      map[string]mapper.TableBinding
	Capture     inbound.CaptureInboundPort
	// SchemaStore trägt die Persistenz-Fähigkeit der dynamischen
	// Re-Versionierung (`outbound.SchemaStorePort`, `ADR-0015`
	// Folgepflicht); ungesetzt (`nil`) bleibt die Relation-Behandlung des
	// Assemblers wirkungslos — bestehende Aufrufstellen (Tests), die
	// dieses Feld nicht setzen, bleiben unverändert kompilierbar.
	SchemaStore outbound.SchemaStorePort
	// Log trägt den injizierten `LogPort` (`LH-QA-OPS-004`, `ADR-0024`);
	// ungesetzt (`nil`) fällt `NewStream` auf `outbound.NoopLog` zurück —
	// bestehende Aufrufstellen (Tests), die dieses Feld nicht setzen,
	// bleiben unverändert kompilierbar.
	Log outbound.LogPort
}

// Stream ist der Replication-Stream-Driving-Adapter
// (`PostgresReplicationStreamAdapter`, `ADR-0006`): er hält die
// Replication-Verbindung und streamt committed Quelltransaktionen in den
// Capture-Pfad.
type Stream struct {
	// conn trägt den konkreten Treiber-Typ für den öffentlichen Rand
	// `Conn()` — Stream und ACK sind getrennte Rollen an **einer**
	// technischen Verbindung (`ADR-0007`, Option C).
	conn *pgconn.PgConn
	// session ist die Naht, an der die Logik dieses Adapters hängt
	// (`ADR-0080`): dieselbe Verbindung hinter der Treiber-Hülle, damit
	// die Empfangs-Schleife und die Slot-/Publication-Auflösung netzlos
	// fahrbar sind.
	session   driverSession
	decoder   *decode.Decoder
	assembler *mapper.Assembler
	capture   inbound.CaptureInboundPort
	// lastAcked trägt die letzte bestätigte Position als LSN —
	// ausschließlich Positionsgröße für Keepalive-Antworten und
	// Slot-Feedback: confirmed_flush_lsn rückt nur über bestätigte
	// Positionen (`LH-QA-REL-001.a`), nicht über den Empfangsstand.
	lastAcked pglogrepl.LSN
	// log trägt die strukturierte Protokollierung über den injizierten
	// `LogPort` (`LH-QA-OPS-004`, `ADR-0024`) — nie `nil` (`NewStream`
	// trägt den `outbound.NoopLog`-Default nach).
	log outbound.LogPort
}

// NewStream baut die Replication-Verbindung auf (`LH-QA-REL-001.a`,
// Schritt Receive): die Verbindung streamt im Text-Format mit dem
// `pgoutput`-Plugin, der Logical Replication Slot wird bei Bedarf
// angelegt und die Publication wird verlangt — ein Start ohne
// Publication ist kein Stand zur Fortsetzung. Die Verbindungs-Konfiguration
// setzt den Replication-Modus am Treiber; Treiber-Fehler gehen in die
// Klasse `replication` (`SPEC-008`).
func NewStream(ctx context.Context, cfg Config) (*Stream, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	log := cfg.Log
	if log == nil {
		log = outbound.NoopLog
	}
	conn, err := connectReplication(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}
	session := connSession{conn: conn}
	startLSN, err := ensureSlot(ctx, log, session, cfg.Slot)
	if err != nil {
		conn.Close(ctx)
		return nil, err
	}
	if err := ensurePublication(ctx, session, cfg.Publication); err != nil {
		conn.Close(ctx)
		return nil, err
	}
	assembler, err := mapper.NewAssembler(cfg.Source, cfg.Tables, cfg.SchemaStore)
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("%w: %v", ErrConfiguration, err)
	}
	if err := session.StartReplication(ctx, cfg.Slot, startLSN, pglogrepl.StartReplicationOptions{
		Mode: pglogrepl.LogicalReplication,
		PluginArgs: []string{
			"proto_version '1'",
			"publication_names '" + cfg.Publication + "'",
		},
	}); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("%w: START_REPLICATION: %v", ErrReplication, err)
	}
	log.Info(ctx, "replication: Stream gestartet",
		"source", cfg.Source, "publication", cfg.Publication, "slot", cfg.Slot)
	return newStreamOnSession(session, conn, assembler, cfg.Capture, log), nil
}

// newStreamOnSession verdrahtet den Stream auf eine Naht — der
// paket-interne Einstieg der netzlosen Tests, die einen Fake anstelle des
// Treibers fahren; `NewStream` reicht die Treiber-Hülle durch. Ungesetztes
// `log` (`nil`) fällt auf `outbound.NoopLog` zurück; `conn` trägt der
// netzlose Aufruf nicht (er bedient allein den öffentlichen Rand `Conn()`).
func newStreamOnSession(session driverSession, conn *pgconn.PgConn, assembler *mapper.Assembler, capture inbound.CaptureInboundPort, log outbound.LogPort) *Stream {
	if log == nil {
		log = outbound.NoopLog
	}
	return &Stream{
		conn:      conn,
		session:   session,
		decoder:   decode.NewDecoder(),
		assembler: assembler,
		capture:   capture,
		log:       log,
	}
}

// validateConfig trägt die Konfigurationsgrenzen vor dem
// Verbindungsaufbau: nichtleere Kennungen und Quelle; Publication- und
// Slot-Namen tragen das Bezeichner-Alphabet (`SPEC-008`, Klasse
// `configuration` — kein Start im falschen Stand). Der Capture-Port
// wird vor `Run` verdrahtet (BindCapture) — die Composition Root
// verdrahtet Stream- und ACK-Adapter über dieselbe Verbindung
// (`ADR-0007`, `ADR-0026`).
func validateConfig(cfg Config) error {
	if cfg.DSN == "" {
		return fmt.Errorf("%w: DSN fehlt", ErrConfiguration)
	}
	if cfg.Source == "" {
		return fmt.Errorf("%w: Quelle ohne Kennung", ErrConfiguration)
	}
	if !identifierShape.MatchString(cfg.Publication) {
		return fmt.Errorf("%w: Publication-Name %q trägt nicht das Bezeichner-Alphabet", ErrConfiguration, cfg.Publication)
	}
	if !identifierShape.MatchString(cfg.Slot) {
		return fmt.Errorf("%w: Slot-Name %q trägt nicht das Bezeichner-Alphabet", ErrConfiguration, cfg.Slot)
	}
	return nil
}

// BindCapture verdrahtet den `CaptureInboundPort` vor dem Lauf; ein
// Lauf ohne Port endet über die Konfigurationsklasse.
func (s *Stream) BindCapture(capture inbound.CaptureInboundPort) error {
	if capture == nil {
		return fmt.Errorf("%w: CaptureInboundPort fehlt", ErrConfiguration)
	}
	s.capture = capture
	return nil
}

// Assembler trägt den laufenden `mapper.Assembler` dieses Streams nach
// außen (`ADR-0050`): die Administrations-Goroutine der Composition Root
// trägt über ihn eine neue oder entfallene `TableBinding` synchronisiert
// nach (`Assembler.AddBinding`/`RemoveBinding`), nachdem eine über SQL
// beantragte Aktivierung/Deaktivierung real ausgeführt wurde — derselbe
// laufende Übersetzer, den `Consume` in `Run` verwendet, kein zweiter.
func (s *Stream) Assembler() *mapper.Assembler {
	return s.assembler
}

// connectReplication baut die Replication-Verbindung über den
// Treiber-Aufrufparameter `replication=database` auf; der DSN bleibt
// unverändert für die übrigen Verbindungszüge (`ADR-0032`).
func connectReplication(ctx context.Context, dsn string) (*pgconn.PgConn, error) {
	connConfig, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReplication, err)
	}
	connConfig.RuntimeParams["replication"] = "database"
	conn, err := pgconn.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: Verbindungsaufbau: %v", ErrReplication, err)
	}
	return conn, nil
}

// slotLSNQuery trägt die Katalogabfrage des Slot-Stands
// (`LH-FA-CFG-001.a`): `confirmed_flush_lsn` des logischen Slots. Der
// Slot-Name trägt das Bezeichner-Alphabet (`identifierShape`) und geht
// deshalb als Bezeichner-Literal in die Abfrage.
func slotLSNQuery(slot string) string {
	return "SELECT confirmed_flush_lsn::text FROM pg_replication_slots WHERE slot_name = '" + slot + "' AND slot_type = 'logical'"
}

// querySingle trägt die erste Zeile einer Katalogabfrage über die Naht;
// ohne Zeile meldet der zweite Ausgang die Abwesenheit.
func querySingle(ctx context.Context, session driverSession, sql string) ([]string, bool, error) {
	results, err := session.Exec(ctx, sql)
	if err != nil {
		return nil, false, fmt.Errorf("%w: Katalogabfrage: %v", ErrReplication, err)
	}
	values, exists := firstRow(results)
	return values, exists, nil
}

// firstRow trägt die Katalog-Zeilen-Übersetzung: die erste Zeile des
// ersten Ergebnisses in Textwerte übersetzt; ohne Zeile meldet der zweite
// Ausgang die Abwesenheit.
func firstRow(results []*pgconn.Result) ([]string, bool) {
	if len(results) == 0 || len(results[0].Rows) == 0 {
		return nil, false
	}
	row := results[0].Rows[0]
	values := make([]string, len(row))
	for i, value := range row {
		values[i] = string(value)
	}
	return values, true
}

// parseLSN trägt die LSN-Übersetzung eines Katalogwerts in die
// Positionsform (`ADR-0005`); ein unlesbarer Wert endet über die
// Fehlerklasse `replication` (`SPEC-008`). `what` benennt die Quelle des
// Werts in der Fehlermeldung.
func parseLSN(what, text string) (pglogrepl.LSN, error) {
	lsn, err := pglogrepl.ParseLSN(text)
	if err != nil {
		return 0, fmt.Errorf("%w: %s %q: %v", ErrReplication, what, text, err)
	}
	return lsn, nil
}

// ensureSlot legt den Logical Replication Slot an, wenn er fehlt
// (`LH-FA-CFG-001.a`: Publication- und Slot-Verwaltung, Output Plugin
// `pgoutput`), und trägt die Startposition: der Stand von
// confirmed_flush_lsn des bestehenden Slots bzw. der konsistente Punkt
// des neu angelegten.
func ensureSlot(ctx context.Context, log outbound.LogPort, session driverSession, slot string) (pglogrepl.LSN, error) {
	values, exists, err := querySingle(ctx, session, slotLSNQuery(slot))
	if err != nil {
		return 0, err
	}
	if exists {
		if values[0] == "" {
			return 0, fmt.Errorf("%w: Slot %q trägt keine confirmed_flush_lsn", ErrReplication, slot)
		}
		startLSN, err := parseLSN("confirmed_flush_lsn", values[0])
		if err != nil {
			return 0, err
		}
		log.Debug(ctx, "replication: bestehender Slot fortgesetzt", "slot", slot, "start_lsn", startLSN)
		return startLSN, nil
	}
	created, err := session.CreateReplicationSlot(ctx, slot, outputPlugin, pglogrepl.CreateReplicationSlotOptions{
		Mode:           pglogrepl.LogicalReplication,
		SnapshotAction: "NOEXPORT_SNAPSHOT",
	})
	if err != nil {
		return 0, fmt.Errorf("%w: CREATE_REPLICATION_SLOT: %v", ErrReplication, err)
	}
	startLSN, err := parseLSN("ConsistentPoint", created.ConsistentPoint)
	if err != nil {
		return 0, err
	}
	log.Info(ctx, "replication: Slot angelegt", "slot", slot, "start_lsn", startLSN)
	return startLSN, nil
}

// ensurePublication verlangt die Publication (`LH-FA-CFG-001.a`): ihr
// Bestand liegt an der Aktivierung der Tabellen; der Stream-Adapter
// startet nicht ohne sie und legt sie nicht still an.
func ensurePublication(ctx context.Context, session driverSession, publication string) error {
	_, exists, err := querySingle(ctx, session,
		"SELECT 1 FROM pg_publication WHERE pubname = '"+publication+"'")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: Publication %q fehlt an der Quelle", ErrConfiguration, publication)
	}
	return nil
}

// Conn trägt die technische Replication-Verbindung; der
// PostgresReplicationAckAdapter führt seine Bestätigung über dieselbe
// Verbindung — Stream und ACK sind getrennte Rollen an einer technischen
// Verbindung (`ADR-0007`, Option C).
func (s *Stream) Conn() *pgconn.PgConn {
	return s.conn
}

// Run streamt, bis der Kontext endet oder ein Fehler auftritt, und
// schließt die Verbindung bei seiner Rückkehr. Jede `pgoutput`-Änderung
// läuft über Dekodierung und Mapper in den `CaptureInboundPort`; die
// Rückkehr ohne Fehler meldet das reguläre Stream-Ende. Ein Fehler
// bricht den Stream ab — Persist-before-ACK und das Verbot des stillen
// Überspringens enden an keiner stillen Fortsetzung (`LH-QA-REL-001.a`,
// `SPEC-008`).
func (s *Stream) Run(ctx context.Context) (err error) {
	if s.capture == nil {
		return fmt.Errorf("%w: CaptureInboundPort fehlt", ErrConfiguration)
	}
	defer s.session.Close(ctx)
	// Ein Log je Rückkehr des Stream-Laufs (`LH-QA-OPS-004`) — der
	// benannte Rückgabewert `err` trägt den Ausgang über alle
	// `return`-Stellen unten hinweg zu diesem einen `defer`, dasselbe
	// Muster wie `bootstrap.Run`/`reportFault`.
	defer func() {
		if err != nil {
			s.log.Error(ctx, "replication: Stream beendet mit Fehler", "error", err)
			return
		}
		s.log.Info(ctx, "replication: Stream regulär beendet")
	}()
	for {
		rawMessage, err := s.session.ReceiveMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("%w: Empfang: %v", ErrReplication, err)
		}
		switch message := rawMessage.(type) {
		case *pgproto3.CopyData:
			if err := s.handleCopyData(ctx, message.Data); err != nil {
				return err
			}
		case *pgproto3.CopyDone:
			// Der Stream ist regulär beendet.
			return nil
		case *pgproto3.ErrorResponse:
			return fmt.Errorf("%w: %v", ErrReplication, pgconn.ErrorResponseToPgError(message))
		default:
			// Übrige Backend-Nachrichten (NoticeResponse u. a.) tragen
			// keine Stream-Änderung.
		}
	}
}

// handleCopyData trägt die Meldungs-Zerlegung des CopyData-Payloads
// (`ADR-0006`): eine XLogData-Nachricht läuft über `process` in den
// Capture-Pfad, eine Primary-Keepalive-Nachricht mit `ReplyRequested`
// wird mit der letzten bestätigten Position beantwortet, ein leeres
// CopyData ist ein sichtbarer Fehler. Unbekannte Byte-IDs tragen keine
// Stream-Änderung.
func (s *Stream) handleCopyData(ctx context.Context, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("%w: leeres CopyData", ErrReplication)
	}
	switch data[0] {
	case pglogrepl.XLogDataByteID:
		payload, err := parseXLogData(data)
		if err != nil {
			return err
		}
		return s.process(ctx, payload)
	case pglogrepl.PrimaryKeepaliveMessageByteID:
		replyRequested, err := parseKeepalive(data)
		if err != nil {
			return err
		}
		if !replyRequested {
			return nil
		}
		// Die Keepalive-Antwort meldet die letzte bestätigte Position —
		// nie den Empfangsstand; die Bestätigung entscheidet die
		// Application nach Persistenz (`ADR-0007`, `LH-QA-REL-001.a`).
		if err := s.session.SendStandbyStatusUpdate(ctx, standbyStatus(s.lastAcked)); err != nil {
			return fmt.Errorf("%w: Keepalive-Antwort: %v", ErrReplication, err)
		}
		return nil
	default:
		return nil
	}
}

// parseXLogData trägt die Zerlegung einer XLogData-Nachricht: der
// Kopfsatz mit der Byte-ID wird geprüft, der Rückgabewert ist der
// `pgoutput`-Payload (`SPEC-010`).
func parseXLogData(data []byte) ([]byte, error) {
	xlogData, err := pglogrepl.ParseXLogData(data[1:])
	if err != nil {
		return nil, fmt.Errorf("%w: XLogData: %v", ErrReplication, err)
	}
	return xlogData.WALData, nil
}

// parseKeepalive trägt die Zerlegung einer Primary-Keepalive-Nachricht:
// der Rückgabewert meldet, ob die Quelle eine Antwort verlangt.
func parseKeepalive(data []byte) (bool, error) {
	keepalive, err := pglogrepl.ParsePrimaryKeepaliveMessage(data[1:])
	if err != nil {
		return false, fmt.Errorf("%w: Keepalive: %v", ErrReplication, err)
	}
	return keepalive.ReplyRequested, nil
}

// standbyStatus trägt die Standby-Status-Form der Keepalive-Antwort: die
// drei Positionen tragen denselben Stand, damit der Slot ihn als
// confirmed_flush_lsn trägt (`LH-QA-REL-001.a`).
func standbyStatus(lsn pglogrepl.LSN) pglogrepl.StandbyStatusUpdate {
	return pglogrepl.StandbyStatusUpdate{
		WALWritePosition: lsn,
		WALFlushPosition: lsn,
		WALApplyPosition: lsn,
	}
}

// process dekodiert einen `pgoutput`-Payload, übersetzt ihn über den
// Mapper und ruft den `CaptureInboundPort` auf; das Capture-Ergebnis
// trägt die bestätigte Position und setzt den Stand des Slot-Feedbacks.
// Ein Fehler endet der Stream — Persist-before-ACK und das Verbot des
// stillen Überspringens gelten an der Adapter-Grenze (`LH-QA-REL-001.a`,
// `SPEC-008`).
func (s *Stream) process(ctx context.Context, payload []byte) error {
	event, err := s.decoder.Decode(payload)
	if err != nil {
		return err
	}
	if event == nil {
		return nil
	}
	command, err := s.assembler.Consume(ctx, event)
	if err != nil {
		return err
	}
	if command == nil {
		return nil
	}
	result, err := s.capture.Capture(ctx, *command)
	if err != nil {
		// Persistenzfehler enden ohne Source-ACK (`LH-QA-REL-001.a`);
		// die Fehlerklasse bleibt am Port-Kontrakt.
		return err
	}
	if !result.Acknowledged.IsZero() {
		s.lastAcked = pglogrepl.LSN(result.Acknowledged.Offset)
	}
	return nil
}
