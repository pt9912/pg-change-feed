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
// Replication- und Katalogabfragen.
var identifierShape = regexp.MustCompile(`^[a-z0-9_]{1,63}$`)

// Config trägt die Konfiguration des Stream-Adapters: die
// Replication-Verbindung (DSN), die Quelle, die Verwaltungsnamen
// Publication und Slot (`LH-FA-CFG-001.a`) und die aktivierten Tabellen
// mit ihren Port-Kennungen. Der Capture-Port wird vor dem Lauf
// verdrahtet (BindCapture) — der ACK-Adapter braucht die Verbindung
// (`ADR-0007`), die `NewStream` erst aufbaut. Die Tabellen-Kennungen
// und die Schema-Versionen liegen bei der Konfiguration; die
// Schema-Evolution trägt sie über den Metadata-Pfad
// (LH-FA-SCH-004.a).
type Config struct {
	DSN         string
	Source      model.SourceID
	Publication string
	Slot        string
	Tables      map[string]mapper.TableBinding
	Capture     inbound.CaptureInboundPort
}

// Stream ist der Replication-Stream-Driving-Adapter
// (`PostgresReplicationStreamAdapter`, `ADR-0006`): er hält die
// Replication-Verbindung und streamt committed Quelltransaktionen in den
// Capture-Pfad.
type Stream struct {
	conn      *pgconn.PgConn
	decoder   *decode.Decoder
	assembler *mapper.Assembler
	capture   inbound.CaptureInboundPort
	// lastAcked trägt die letzte bestätigte Position als LSN —
	// ausschließlich Positionsgröße für Keepalive-Antworten und
	// Slot-Feedback: confirmed_flush_lsn rückt nur über bestätigte
	// Positionen (`LH-QA-REL-001.a`), nicht über den Empfangsstand.
	lastAcked pglogrepl.LSN
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
	conn, err := connectReplication(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}
	startLSN, err := ensureSlot(ctx, conn, cfg.Slot)
	if err != nil {
		conn.Close(ctx)
		return nil, err
	}
	if err := ensurePublication(ctx, conn, cfg.Publication); err != nil {
		conn.Close(ctx)
		return nil, err
	}
	assembler, err := mapper.NewAssembler(cfg.Source, cfg.Tables)
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("%w: %v", ErrConfiguration, err)
	}
	if err := pglogrepl.StartReplication(ctx, conn, cfg.Slot, startLSN, pglogrepl.StartReplicationOptions{
		Mode: pglogrepl.LogicalReplication,
		PluginArgs: []string{
			"proto_version '1'",
			"publication_names '" + cfg.Publication + "'",
		},
	}); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("%w: START_REPLICATION: %v", ErrReplication, err)
	}
	return &Stream{
		conn:      conn,
		decoder:   decode.NewDecoder(),
		assembler: assembler,
		capture:   cfg.Capture,
	}, nil
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

// querySingle trägt die erste Zeile einer Katalogabfrage; ohne Zeile
// meldet der zweite Ausgang die Abwesenheit.
func querySingle(ctx context.Context, conn *pgconn.PgConn, sql string) ([]string, bool, error) {
	results, err := conn.Exec(ctx, sql).ReadAll()
	if err != nil {
		return nil, false, fmt.Errorf("%w: Katalogabfrage: %v", ErrReplication, err)
	}
	if len(results) == 0 || len(results[0].Rows) == 0 {
		return nil, false, nil
	}
	row := results[0].Rows[0]
	values := make([]string, len(row))
	for i, value := range row {
		values[i] = string(value)
	}
	return values, true, nil
}

// ensureSlot legt den Logical Replication Slot an, wenn er fehlt
// (`LH-FA-CFG-001.a`: Publication- und Slot-Verwaltung, Output Plugin
// `pgoutput`), und trägt die Startposition: der Stand von
// confirmed_flush_lsn des bestehenden Slots bzw. der konsistente Punkt
// des neu angelegten.
func ensureSlot(ctx context.Context, conn *pgconn.PgConn, slot string) (pglogrepl.LSN, error) {
	values, exists, err := querySingle(ctx, conn,
		"SELECT confirmed_flush_lsn::text FROM pg_replication_slots WHERE slot_name = '"+slot+"' AND slot_type = 'logical'")
	if err != nil {
		return 0, err
	}
	if exists {
		if values[0] == "" {
			return 0, fmt.Errorf("%w: Slot %q trägt keine confirmed_flush_lsn", ErrReplication, slot)
		}
		startLSN, err := pglogrepl.ParseLSN(values[0])
		if err != nil {
			return 0, fmt.Errorf("%w: confirmed_flush_lsn %q: %v", ErrReplication, values[0], err)
		}
		return startLSN, nil
	}
	created, err := pglogrepl.CreateReplicationSlot(ctx, conn, slot, outputPlugin, pglogrepl.CreateReplicationSlotOptions{
		Mode:           pglogrepl.LogicalReplication,
		SnapshotAction: "NOEXPORT_SNAPSHOT",
	})
	if err != nil {
		return 0, fmt.Errorf("%w: CREATE_REPLICATION_SLOT: %v", ErrReplication, err)
	}
	startLSN, err := pglogrepl.ParseLSN(created.ConsistentPoint)
	if err != nil {
		return 0, fmt.Errorf("%w: ConsistentPoint %q: %v", ErrReplication, created.ConsistentPoint, err)
	}
	return startLSN, nil
}

// ensurePublication verlangt die Publication (`LH-FA-CFG-001.a`): ihr
// Bestand liegt an der Aktivierung der Tabellen; der Stream-Adapter
// startet nicht ohne sie und legt sie nicht still an.
func ensurePublication(ctx context.Context, conn *pgconn.PgConn, publication string) error {
	_, exists, err := querySingle(ctx, conn,
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
func (s *Stream) Run(ctx context.Context) error {
	if s.capture == nil {
		return fmt.Errorf("%w: CaptureInboundPort fehlt", ErrConfiguration)
	}
	defer s.conn.Close(ctx)
	for {
		rawMessage, err := s.conn.ReceiveMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("%w: Empfang: %v", ErrReplication, err)
		}
		switch message := rawMessage.(type) {
		case *pgproto3.CopyData:
			if len(message.Data) == 0 {
				return fmt.Errorf("%w: leeres CopyData", ErrReplication)
			}
			switch message.Data[0] {
			case pglogrepl.XLogDataByteID:
				xlogData, err := pglogrepl.ParseXLogData(message.Data[1:])
				if err != nil {
					return fmt.Errorf("%w: XLogData: %v", ErrReplication, err)
				}
				if err := s.process(ctx, xlogData.WALData); err != nil {
					return err
				}
			case pglogrepl.PrimaryKeepaliveMessageByteID:
				keepalive, err := pglogrepl.ParsePrimaryKeepaliveMessage(message.Data[1:])
				if err != nil {
					return fmt.Errorf("%w: Keepalive: %v", ErrReplication, err)
				}
				if keepalive.ReplyRequested {
					// Die Keepalive-Antwort meldet die letzte
					// bestätigte Position — nie den Empfangsstand;
					// die Bestätigung entscheidet die Application
					// nach Persistenz (`ADR-0007`, `LH-QA-REL-001.a`).
					if err := pglogrepl.SendStandbyStatusUpdate(ctx, s.conn, pglogrepl.StandbyStatusUpdate{
						WALWritePosition: s.lastAcked,
						WALFlushPosition: s.lastAcked,
						WALApplyPosition: s.lastAcked,
					}); err != nil {
						return fmt.Errorf("%w: Keepalive-Antwort: %v", ErrReplication, err)
					}
				}
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
	command, err := s.assembler.Consume(event)
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
