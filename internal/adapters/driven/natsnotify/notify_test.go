package natsnotify_test

import (
	"context"
	stderrors "errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/natsnotify"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// Die natsnotify-Adapter-Tests tragen die Grenzen des Port-Kontrakts: den
// Konstruktions-Fehler ohne Verbindung, die Empty-Source-Grenze und das
// Fehlerklassen-Wrapping des Publish-Aufrufs (`SPEC-008`, Klasse
// `transient`). Der echte Publish-Erfolgsbeleg gegen einen laufenden
// NATS-Server trägt `make test-notify` (`CDC_NATS_TEST_URL`,
// `tools/harness/run-notify-tests.sh`); die übrigen Tests brauchen keinen
// Server (Verbindungsfehler bleiben auf Loopback, auch unter `--network
// none`, `make test`).

// newDisconnectedConn baut eine `*nats.Conn`, die real nie erfolgreich
// verbindet: `RetryOnFailedConnect` lässt `Connect` trotzdem einen Conn im
// Reconnecting-Zustand zurückgeben (Ziel-Adresse Port 1, kein Server dort).
// Das trägt die Konstruktions- und Empty-Source-Tests ohne echten Server;
// die Verbindung selbst wird für den Wrapping-Test anschließend geschlossen.
func newDisconnectedConn(t *testing.T) *nats.Conn {
	t.Helper()
	conn, err := nats.Connect("nats://127.0.0.1:1",
		nats.RetryOnFailedConnect(true),
		nats.NoReconnect(),
		nats.Timeout(200*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("nats.Connect (erwartungsgemäß ohne Server, reconnecting): %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

// TestNewRequiresConnection trägt die Konstruktions-Grenze: ein
// Notify-Adapter ohne NATS-Verbindung endet über die Klasse `transient`
// (`SPEC-008`) — kein Publish-Versuch ohne Verbindung.
func TestNewRequiresConnection(t *testing.T) {
	_, err := natsnotify.New(nil)
	if !stderrors.Is(err, outbound.ErrNotify) {
		t.Fatalf("New(nil): %v", err)
	}
}

// TestNotifyRejectsEmptySourceID trägt die Empty-Source-Grenze: die
// Prüfung läuft vor jedem Zugriff auf die Verbindung, ein
// nie-verbundener Conn genügt.
func TestNotifyRejectsEmptySourceID(t *testing.T) {
	adapter, err := natsnotify.New(newDisconnectedConn(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := adapter.Notify(context.Background(), "", "public", "tbl"); !stderrors.Is(err, outbound.ErrNotify) {
		t.Fatalf("Notify(\"\", ...): %v", err)
	}
}

// TestNotifyRejectsEmptySchemaOrTable trägt dieselbe Grenze für Schema und
// Tabelle: beide sind Pflichtangaben des vier-Ebenen-Subjekts (`ADR-0056`).
func TestNotifyRejectsEmptySchemaOrTable(t *testing.T) {
	adapter, err := natsnotify.New(newDisconnectedConn(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := adapter.Notify(context.Background(), "src-1", "", "tbl"); !stderrors.Is(err, outbound.ErrNotify) {
		t.Fatalf("Notify(..., \"\", \"tbl\"): %v", err)
	}
	if err := adapter.Notify(context.Background(), "src-1", "public", ""); !stderrors.Is(err, outbound.ErrNotify) {
		t.Fatalf("Notify(..., \"public\", \"\"): %v", err)
	}
}

// TestNotifyRejectsReservedSubjectCharacters trägt die defensive
// Validierung gegen NATS-reservierte Zeichen und Whitespace in Schema-
// oder Tabellenname (`ADR-0056` Folgepflicht) — vor jedem Publish-Versuch,
// ein nie-verbundener Conn genügt.
func TestNotifyRejectsReservedSubjectCharacters(t *testing.T) {
	adapter, err := natsnotify.New(newDisconnectedConn(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	cases := []struct {
		name   string
		schema string
		table  string
	}{
		{"Punkt im Schema", "pub.lic", "tbl"},
		{"Punkt in der Tabelle", "public", "tbl.name"},
		{"Stern in der Tabelle", "public", "tbl*"},
		{"Größer-als im Schema", "pub>lic", "tbl"},
		{"Whitespace in der Tabelle", "public", "tbl name"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := adapter.Notify(context.Background(), "src-1", tc.schema, tc.table); !stderrors.Is(err, outbound.ErrNotify) {
				t.Fatalf("Notify(%q, %q): %v", tc.schema, tc.table, err)
			}
		})
	}
}

// TestNotifyWrapsPublishFailureAsTransient trägt das Fehlerklassen-Wrapping
// am Treiber-Fehler: eine geschlossene Verbindung endet über
// `outbound.ErrNotify`; die technische Ursache bleibt über die zweite
// Wrappung lesbar (`ADR-0023`).
func TestNotifyWrapsPublishFailureAsTransient(t *testing.T) {
	conn := newDisconnectedConn(t)
	adapter, err := natsnotify.New(conn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	conn.Close()

	if err := adapter.Notify(context.Background(), "src-1", "public", "tbl"); !stderrors.Is(err, outbound.ErrNotify) {
		t.Fatalf("Notify nach Close: %v", err)
	}
}

// TestNotifyPublishesEmptyPayloadOnSubject trägt den echten
// Publish-Erfolgsbeleg gegen einen laufenden NATS-Server (`SPEC-017`,
// `ADR-0056`): Subjekt-Schema `cdc.changes.<source_id>.<schema>.<table>`,
// leerer Payload — kein Change-Inhalt, keine Positionsangabe. Läuft nur
// über `make test-notify` (`CDC_NATS_TEST_URL` gesetzt); ohne Server
// übersprungen.
func TestNotifyPublishesEmptyPayloadOnSubject(t *testing.T) {
	url := os.Getenv("CDC_NATS_TEST_URL")
	if url == "" {
		t.Skip("CDC_NATS_TEST_URL nicht gesetzt — reale natsnotify-Adapter-Tests laufen über make test-notify")
	}
	conn, err := nats.Connect(url)
	if err != nil {
		t.Fatalf("nats.Connect(%q): %v", url, err)
	}
	t.Cleanup(func() { conn.Close() })

	sub, err := conn.SubscribeSync("cdc.changes.>")
	if err != nil {
		t.Fatalf("SubscribeSync: %v", err)
	}
	t.Cleanup(func() { _ = sub.Unsubscribe() })
	if err := conn.Flush(); err != nil {
		t.Fatalf("Flush nach Subscribe: %v", err)
	}

	adapter, err := natsnotify.New(conn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := adapter.Notify(context.Background(), "src-1", "public", "tbl"); err != nil {
		t.Fatalf("Notify: %v", err)
	}

	msg, err := sub.NextMsg(3 * time.Second)
	if err != nil {
		t.Fatalf("NextMsg: %v", err)
	}
	if msg.Subject != "cdc.changes.src-1.public.tbl" {
		t.Fatalf("Subjekt = %q, wollen %q", msg.Subject, "cdc.changes.src-1.public.tbl")
	}
	if len(msg.Data) != 0 {
		t.Fatalf("Payload = %q, wollen leer (SPEC-017: kein Change-Inhalt)", msg.Data)
	}
}

// recordingLog trägt den `LogPort` der Tests (`LH-QA-OPS-004`, `ADR-0024`):
// er hält die Aufrufe fest, ohne sie auszugeben.
type recordingLog struct {
	infos  []string
	errors []string
}

func (l *recordingLog) Debug(context.Context, string, ...any) {}
func (l *recordingLog) Info(_ context.Context, msg string, _ ...any) {
	l.infos = append(l.infos, msg)
}
func (l *recordingLog) Warn(context.Context, string, ...any) {}
func (l *recordingLog) Error(_ context.Context, msg string, _ ...any) {
	l.errors = append(l.errors, msg)
}

var _ outbound.LogPort = (*recordingLog)(nil)

// TestNewWithLogReichtDenLogPortDurch trägt beide Kettenglieder der
// Konstruktions-Option (`ADR-0024`, `LH-QA-OPS-004`): der über `WithLog`
// übergebene `LogPort` ist der, über den der Adapter protokolliert —
// erstens im Konstruktionsaufruf selbst, zweitens von seinem **Adapterfeld**
// aus an einem späteren Aufruf (`notifyFailure`). Das zweite Glied ist das
// mittlere der Kette: trüge der Adapter an dieser Stelle den Default
// `outbound.NoopLog`, käme die Konstruktionszeile weiterhin an, die
// Signalfehler-Zeile nicht. Beide Aufzeichnungen hängen am eingegebenen Wert
// — die Verbindung ist getrennt, der Publish scheitert, der Fehlerpfad läuft
// netzlos.
//
// Rot färbende Mutation: in `New` das Adapterfeld auf `outbound.NoopLog`
// festlegen (`log: outbound.NoopLog`) — dann fehlt die Fehler-Aufzeichnung
// und dieser Test färbt rot. Rot ebenfalls: in `newOptions` die
// Options-Schleife fallenlassen — dann fehlen beide Aufzeichnungen.
func TestNewWithLogReichtDenLogPortDurch(t *testing.T) {
	log := &recordingLog{}
	conn := newDisconnectedConn(t)
	adapter, err := natsnotify.New(conn, natsnotify.WithLog(log))
	if err != nil {
		t.Fatalf("New mit WithLog: %v", err)
	}
	if len(log.infos) != 1 || !strings.Contains(log.infos[0], "natsnotify") {
		t.Fatalf("Aufzeichnung des Konstruktionsaufrufs: %q (Erwartung: die Konstruktionszeile über den injizierten LogPort)", log.infos)
	}

	conn.Close()
	if err := adapter.Notify(context.Background(), "src-1", "public", "tbl"); !stderrors.Is(err, outbound.ErrNotify) {
		t.Fatalf("Notify nach Close: %v", err)
	}
	if len(log.errors) != 1 || !strings.Contains(log.errors[0], "natsnotify") {
		t.Fatalf("Aufzeichnung des Adapterfelds: %q (Erwartung: die Signalfehler-Zeile über denselben LogPort)", log.errors)
	}
}
