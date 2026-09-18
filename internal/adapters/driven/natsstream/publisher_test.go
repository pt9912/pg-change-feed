package natsstream

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die natsstream-Adapter-Tests tragen die beiden Regressionsklassen aus
// `ADR-0100` §Fitness Function, die dieses Paket netzlos (`make test`)
// tragen kann: die Subjekt-Wurzel `cdc.stream` (niemals `cdc.changes`,
// Abgrenzung gegen `ADR-0056`s Wecksignal-Namensraum) und die
// Fire-and-Forget-Eigenschaft (ein Publish ohne verbundenen NATS-Server
// blockiert nicht). Der reale Empfangsbeleg gegen einen laufenden
// NATS-Server (ein abonnierter Test-Client erhält das vollständige
// JSON-Event) trägt `make test-integration` — derselbe Server-Bedarf, den
// `natsnotify`s eigener Empfangsbeleg an `make test-notify` statt an
// `make test` bindet (`internal/adapters/driven/natsnotify/notify_test.go`):
// ein echter Empfangsbeleg braucht einen erreichbaren NATS-Server, den
// `make test`s `--network none`-Lauf strukturell nicht trägt.

// newDisconnectedConn baut eine `*nats.Conn`, die real nie erfolgreich
// verbindet — dasselbe Muster wie `natsnotify`s gleichnamiger Testhelfer
// (eigenständig geführt, kein Test-Utility-Import zwischen Adapter-Paketen).
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

// fakeSubscriber trägt eine minimale `changeSubscriber`-Fälschung: ein
// gepufferter Kanal, den der Test selbst befüllt, und eine
// Aufruf-zählende Abmelde-Funktion.
type fakeSubscriber struct {
	changes      chan *model.Change
	cancelCalled int
}

func newFakeSubscriber() *fakeSubscriber {
	return &fakeSubscriber{changes: make(chan *model.Change, 4)}
}

func (f *fakeSubscriber) Subscribe() (<-chan *model.Change, func()) {
	return f.changes, func() { f.cancelCalled++ }
}

// testChange trägt einen gültigen Change für die Publish-Tests.
func testChange(t *testing.T) *model.Change {
	t.Helper()
	change, err := model.NewChange("change-1", "tx-1", "table-1", 1, model.OperationInsert, nil, []byte(`{"id":1}`), "table-1-v1")
	if err != nil {
		t.Fatalf("Change bauen: %v", err)
	}
	change.Schema = "public"
	change.Table = "orders"
	return &change
}

// recordingLog trägt den `LogPort` der Tests (`LH-QA-OPS-004`, `ADR-0024`):
// er hält die Aufrufe fest, ohne sie auszugeben.
type recordingLog struct {
	infos []string
	warns []string
}

func (l *recordingLog) Debug(context.Context, string, ...any) {}
func (l *recordingLog) Info(_ context.Context, msg string, _ ...any) {
	l.infos = append(l.infos, msg)
}
func (l *recordingLog) Warn(_ context.Context, msg string, _ ...any) {
	l.warns = append(l.warns, msg)
}
func (l *recordingLog) Error(context.Context, string, ...any) {}

var _ outbound.LogPort = (*recordingLog)(nil)

// TestNewRequiresConnection trägt die Konstruktions-Grenze ohne Verbindung.
func TestNewRequiresConnection(t *testing.T) {
	_, err := New(nil, newFakeSubscriber(), "src-1")
	if !stderrors.Is(err, ErrPublish) {
		t.Fatalf("New(nil, ...): %v", err)
	}
}

// TestNewRequiresSubscriber trägt die Konstruktions-Grenze ohne
// Broadcaster: kein Publisher ohne etwas, das er abonnieren könnte.
func TestNewRequiresSubscriber(t *testing.T) {
	_, err := New(newDisconnectedConn(t), nil, "src-1")
	if !stderrors.Is(err, ErrPublish) {
		t.Fatalf("New(conn, nil, ...): %v", err)
	}
}

// TestNewRequiresSourceID trägt die Konstruktions-Grenze ohne Quelle.
func TestNewRequiresSourceID(t *testing.T) {
	_, err := New(newDisconnectedConn(t), newFakeSubscriber(), "")
	if !stderrors.Is(err, ErrPublish) {
		t.Fatalf("New(conn, sub, \"\"): %v", err)
	}
}

// TestSubjectRootIsNeverCdcChanges trägt die erste Regressionsklasse aus
// `ADR-0100` §Fitness Function: das publizierte Subjekt trägt den
// Wurzel-Token `cdc.stream`, niemals `cdc.changes` — Regressionstest gegen
// eine versehentliche Kollision mit dem Wecksignal-Namensraum aus
// `ADR-0056`.
//
// Rot färbende Mutation: `subjectPrefix` auf `"cdc.changes."` ändern —
// dieser Test färbt sofort rot.
func TestSubjectRootIsNeverCdcChanges(t *testing.T) {
	got := subjectFor("src-1", "public", "orders")
	want := "cdc.stream.src-1.public.orders"
	if got != want {
		t.Fatalf("subjectFor(...) = %q, Erwartung %q", got, want)
	}
	if strings.HasPrefix(got, "cdc.changes.") {
		t.Fatalf("subjectFor(...) trägt den Wecksignal-Wurzel-Token: %q", got)
	}
	if !strings.HasPrefix(got, "cdc.stream.") {
		t.Fatalf("subjectFor(...) trägt nicht den Vollinhalts-Wurzel-Token: %q", got)
	}
}

// TestPublishWithoutConnectedServerDoesNotBlock trägt die zweite
// Regressionsklasse aus `ADR-0100` §Fitness Function: ein `Publish` ohne
// verbundenen NATS-Server blockiert nicht (Fire-and-Forget), dieselbe
// Zusicherung wie beim `Broadcaster` selbst
// (`grpcstream.TestPublishOhneAbonnentenBlockiertNicht`).
func TestPublishWithoutConnectedServerDoesNotBlock(t *testing.T) {
	sub := newFakeSubscriber()
	p, err := New(newDisconnectedConn(t), sub, "src-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	fertig := make(chan struct{})
	go func() {
		p.publish(context.Background(), testChange(t))
		close(fertig)
	}()
	select {
	case <-fertig:
	case <-time.After(2 * time.Second):
		t.Fatal("publish ohne verbundenen NATS-Server blockiert")
	}
}

// TestRunPublishesUntilContextEnds trägt den Run-Loop: ein über den
// gefälschten Subscriber eintreffender Change löst einen (nicht
// blockierenden) Publish-Versuch aus, und `Run` kehrt zurück, sobald `ctx`
// endet — die Abmelde-Funktion wird dabei genau einmal aufgerufen.
func TestRunPublishesUntilContextEnds(t *testing.T) {
	sub := newFakeSubscriber()
	p, err := New(newDisconnectedConn(t), sub, "src-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		p.Run(ctx)
		close(done)
	}()

	sub.changes <- testChange(t)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run kehrt nach Kontext-Ende nicht zurück")
	}
	if sub.cancelCalled != 1 {
		t.Fatalf("Abmelde-Funktion %d mal aufgerufen, Erwartung 1", sub.cancelCalled)
	}
}

// TestRunSkipsNilChange trägt die defensive Behandlung eines nil-Werts auf
// dem Kanal: `Run` überspringt ihn, statt zu blockieren oder zu paniken.
func TestRunSkipsNilChange(t *testing.T) {
	sub := newFakeSubscriber()
	p, err := New(newDisconnectedConn(t), sub, "src-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		p.Run(ctx)
		close(done)
	}()

	sub.changes <- nil
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run kehrt nach nil-Change und Kontext-Ende nicht zurück")
	}
}

// TestPublishSkipsReservedSubjectToken trägt dieselbe defensive Validierung
// wie `natsnotify` (`ADR-0056` Folgepflicht): ein Schema/Tabellenname mit
// NATS-reserviertem Zeichen wird ohne Publish-Versuch übersprungen und
// protokolliert, statt das Subjekt unbeabsichtigt aufzuspalten.
func TestPublishSkipsReservedSubjectToken(t *testing.T) {
	log := &recordingLog{}
	sub := newFakeSubscriber()
	p, err := New(newDisconnectedConn(t), sub, "src-1", WithLog(log))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	change := testChange(t)
	change.Table = "tbl.name"

	p.publish(context.Background(), change)

	if len(log.warns) != 1 || !strings.Contains(log.warns[0], "reserviertes Zeichen") {
		t.Fatalf("Warn-Aufzeichnung: %q (Erwartung: Hinweis auf reserviertes Zeichen)", log.warns)
	}
}

// TestToStreamMessageCarriesAllTenFields trägt die Nachrichtenform
// (`SPEC-021`/`SPEC-024`): dieselben zehn Felder wie `model.Change`, ein
// fehlendes Row Image wird zu JSON `null`.
func TestToStreamMessageCarriesAllTenFields(t *testing.T) {
	change := testChange(t)
	change.OldImage = nil

	msg := toStreamMessage(change)
	payload, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	for _, field := range []string{
		"change_id", "transaction_id", "source_table_id", "sequence", "operation",
		"old_image", "new_image", "schema_version", "schema", "table",
	} {
		if _, found := decoded[field]; !found {
			t.Fatalf("Feld %q fehlt in der Nachrichtenform: %s", field, payload)
		}
	}
	if string(decoded["old_image"]) != "null" {
		t.Fatalf("old_image = %s, Erwartung null bei fehlendem Bild", decoded["old_image"])
	}
	if string(decoded["schema"]) != `"public"` || string(decoded["table"]) != `"orders"` {
		t.Fatalf("schema/table = %s/%s, Erwartung \"public\"/\"orders\"", decoded["schema"], decoded["table"])
	}
}

// TestNewWithLogReichtDenLogPortDurch trägt beide Kettenglieder der
// Konstruktions-Option (`ADR-0024`, `LH-QA-OPS-004`): der über `WithLog`
// übergebene `LogPort` ist der, über den der Publisher protokolliert —
// erstens im Konstruktionsaufruf, zweitens von seinem Adapterfeld aus an
// einem späteren Aufruf.
func TestNewWithLogReichtDenLogPortDurch(t *testing.T) {
	log := &recordingLog{}
	p, err := New(newDisconnectedConn(t), newFakeSubscriber(), "src-1", WithLog(log))
	if err != nil {
		t.Fatalf("New mit WithLog: %v", err)
	}
	if len(log.infos) != 1 || !strings.Contains(log.infos[0], "natsstream") {
		t.Fatalf("Aufzeichnung des Konstruktionsaufrufs: %q", log.infos)
	}

	change := testChange(t)
	change.Schema = "pub lic"
	p.publish(context.Background(), change)
	if len(log.warns) != 1 || !strings.Contains(log.warns[0], "natsstream") {
		t.Fatalf("Aufzeichnung des Adapterfelds: %q", log.warns)
	}
}
