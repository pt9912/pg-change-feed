package postgresstorage_test

import (
	"context"
	stderrors "errors"
	"math"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Consumer-State-Tests laufen gegen dieselbe reale PostgreSQL-Instanz
// wie die Store-Tests (`make test-store`, `ADR-0030`); ohne DSN
// überspringen sie. Die Consumer-State-Tabellen (`cdc.consumer`,
// `cdc.consumer_position`) trägt der d-migrate-Rollout (`ADR-0043`) — die
// handgeschriebene DDL des Store-Adapters trägt sie nicht; der
// `test-store`-Lauf rollt das Schema vor dem Testlauf aus.
const (
	consumerStateSource      = "src-consumer"
	consumerStateSourceOther = "src-consumer-other"
	testConsumer             = "con-1"
	testConsumerOther        = "con-2"
)

// newTestConsumerState baut den Consumer-State-Adapter gegen die
// Test-Instanz. Der Datenstand räumt der Test je Fall ab — die Zeilen
// bleiben im Container; das Schema räumt der d-migrate-Rollout des
// Runner-Laufs, der Test trägt es nicht zurück (die handgeschriebene DDL
// des Store-Adapters trägt die Consumer-State-Tabellen nicht, `ADR-0043`).
func newTestConsumerState(t *testing.T) (*postgresstorage.PostgresConsumerStateAdapter, *pgxpool.Pool) {
	t.Helper()
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	var tables int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'cdc' AND table_name IN ('consumer', 'consumer_position')",
	).Scan(&tables); err != nil {
		t.Fatalf("Consumer-State-Tabellen-Prüfung: %v", err)
	}
	if tables != 2 {
		t.Fatalf("Consumer-State-Tabellen unvollständig (%d von 2) — der Schema-Rollout über make schema-rollout trägt sie (ADR-0043); der test-store-Lauf rollt sie vor dem Testlauf aus", tables)
	}
	for _, statement := range []string{
		"DELETE FROM cdc.consumer_position",
		"DELETE FROM cdc.consumer",
	} {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("Datenstand-Rückbau: %v", err)
		}
	}
	for _, source := range []string{consumerStateSource, consumerStateSourceOther} {
		if _, err := pool.Exec(ctx,
			"INSERT INTO cdc.source (source_id, name) VALUES ($1, 'Consumer-Quelle') ON CONFLICT (source_id) DO NOTHING",
			source,
		); err != nil {
			t.Fatalf("Quelle-Zeile: %v", err)
		}
	}
	adapter, err := postgresstorage.NewConsumerState(ctx, dsn)
	if err != nil {
		t.Fatalf("NewConsumerState: %v", err)
	}
	t.Cleanup(adapter.Close)
	return adapter, pool
}

// registerConsumer registriert einen Consumer vor der Bestätigung; die
// Registrierung trägt die Fremdschlüssel-Vorbedingung der Position-Zeile
// (`SPEC-001`).
func registerConsumer(t *testing.T, adapter *postgresstorage.PostgresConsumerStateAdapter, id string) {
	t.Helper()
	registered, err := adapter.Register(context.Background(), mustConsumer(t, id, "Consumer "+id))
	if err != nil {
		t.Fatalf("Register %s: %v", id, err)
	}
	if !registered {
		t.Fatalf("Register %s: bereits registriert", id)
	}
}

// mustConsumer legt einen Consumer an; ein Fehler bricht den Test ab.
func mustConsumer(t *testing.T, id, name string) model.Consumer {
	t.Helper()
	consumer, err := model.NewConsumer(model.ConsumerID(id), name)
	if err != nil {
		t.Fatalf("NewConsumer %s: %v", id, err)
	}
	return consumer
}

// sourcePosition trägt eine Position der Quelle; ein Offset von 0 und
// darüber hinaus trägt der Konstruktor nicht (`SPEC-003`).
func consumerPosition(t *testing.T, source string, offset uint64) model.SourcePosition {
	t.Helper()
	p, err := model.NewSourcePosition(model.SourceID(source), offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	return p
}

// storedOffset liest die bestätigte Position über den Adapter; die
// Bestätigungs-Behauptungen urteilen über den gespeicherten Stand.
func storedOffset(t *testing.T, adapter *postgresstorage.PostgresConsumerStateAdapter, id string) (uint64, bool) {
	t.Helper()
	carried, err := adapter.Position(context.Background(), model.ConsumerID(id))
	if err != nil {
		t.Fatalf("Position %s: %v", id, err)
	}
	if !carried.Acknowledged() {
		return 0, false
	}
	return carried.Position.Offset, true
}

// TestRegisterIsIdempotent trägt die Registrierung (`LH-FA-CON-001`):
// Happy Path und Boundary — die erneut registrierte Kennung bleibt ohne
// Wirkung, die Consumer-Zeilen bleiben eine.
func TestRegisterIsIdempotent(t *testing.T) {
	adapter, pool := newTestConsumerState(t)
	registerConsumer(t, adapter, testConsumer)

	again, err := adapter.Register(context.Background(), mustConsumer(t, testConsumer, "Consumer erneut"))
	if err != nil {
		t.Fatalf("erneutes Register: %v", err)
	}
	if again {
		t.Fatalf("erneute Registrierung meldet neu registriert (Idempotenz, LH-FA-CON-001)")
	}
	var count int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM cdc.consumer WHERE consumer_id = $1", testConsumer).Scan(&count); err != nil {
		t.Fatalf("Consumer-Zeilen: %v", err)
	}
	if count != 1 {
		t.Fatalf("Consumer-Zeilen = %d, wollen 1", count)
	}
}

// TestAcknowledgeCarriesPosition trägt die Persistierung der
// Verarbeitungsposition (`LH-FA-CON-003` Happy Path): der bestätigte
// Wert liest sich zurück.
func TestAcknowledgeCarriesPosition(t *testing.T) {
	adapter, _ := newTestConsumerState(t)
	registerConsumer(t, adapter, testConsumer)

	carried, err := adapter.Acknowledge(context.Background(), model.ConsumerPosition{
		ConsumerID: testConsumer,
		Position:   consumerPosition(t, consumerStateSource, 100),
	})
	if err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}
	if carried.Position.Offset != 100 || carried.Position.SourceID != consumerStateSource {
		t.Fatalf("fortgeführte Position: %+v", carried)
	}
	if got, ok := storedOffset(t, adapter, testConsumer); !ok || got != 100 {
		t.Fatalf("gespeicherte Position = %d (bestätigt: %v), wollen 100", got, ok)
	}
}

// TestAcknowledgeIsMonotonic trägt den monotonen ACK-Vertrag
// (`ADR-0029`, Regel 2; `LH-FA-CON-003` Boundary, `LH-FA-CON-004`
// Boundary): eine frühere Position endet über die Invariante, die
// Wiederholung derselben Position ist idempotent, der gespeicherte
// Fortschritt bleibt bei der abgewiesenen Bestätigung stehen.
func TestAcknowledgeIsMonotonic(t *testing.T) {
	adapter, _ := newTestConsumerState(t)
	registerConsumer(t, adapter, testConsumer)
	ctx := context.Background()
	ack := func(offset uint64) error {
		_, err := adapter.Acknowledge(ctx, model.ConsumerPosition{
			ConsumerID: testConsumer,
			Position:   consumerPosition(t, consumerStateSource, offset),
		})
		return err
	}

	if err := ack(200); err != nil {
		t.Fatalf("erste Bestätigung 200: %v", err)
	}
	// Eine frühere Position endet über die Invariante; der gespeicherte
	// Fortschritt bleibt stehen (`ADR-0029`, Regel 2).
	if err := ack(100); !stderrors.Is(err, domainerrors.ErrPositionRegression) {
		t.Fatalf("rückläufige Bestätigung 100: %v (Erwartung: ErrPositionRegression)", err)
	}
	if got, ok := storedOffset(t, adapter, testConsumer); !ok || got != 200 {
		t.Fatalf("gespeicherte Position nach abgewiesener Bestätigung = %d (bestätigt: %v), wollen 200", got, ok)
	}
	// Die Wiederholung derselben Position ist idempotent
	// (`LH-FA-CON-004` Boundary).
	if err := ack(200); err != nil {
		t.Fatalf("wiederholte Bestätigung 200: %v (Idempotenz meldet keinen Fehler)", err)
	}
	if got, ok := storedOffset(t, adapter, testConsumer); !ok || got != 200 {
		t.Fatalf("gespeicherte Position nach Wiederholung = %d (bestätigt: %v), wollen 200", got, ok)
	}
	if err := ack(300); err != nil {
		t.Fatalf("fortlaufende Bestätigung 300: %v", err)
	}
	if got, ok := storedOffset(t, adapter, testConsumer); !ok || got != 300 {
		t.Fatalf("gespeicherte Position = %d (bestätigt: %v), wollen 300", got, ok)
	}
}

// TestAcknowledgeCarriesSourceBinding trägt die Quell-Bindung der
// bestätigten Position (`ADR-0005`): die erste Bestätigung bindet die
// Quelle, eine Position einer anderen Quelle ist keine
// Bestätigungs-Größe (`ADR-0029` Regel 2 über die Domänen-Ordnung).
func TestAcknowledgeCarriesSourceBinding(t *testing.T) {
	adapter, _ := newTestConsumerState(t)
	registerConsumer(t, adapter, testConsumer)

	carried, err := adapter.Acknowledge(context.Background(), model.ConsumerPosition{
		ConsumerID: testConsumer,
		Position:   consumerPosition(t, consumerStateSource, 100),
	})
	if err != nil {
		t.Fatalf("erste Bestätigung: %v", err)
	}
	if carried.Position.SourceID != consumerStateSource {
		t.Fatalf("gebundene Quelle: %q", carried.Position.SourceID)
	}
	if _, err := adapter.Acknowledge(context.Background(), model.ConsumerPosition{
		ConsumerID: testConsumer,
		Position:   consumerPosition(t, consumerStateSourceOther, 200),
	}); !stderrors.Is(err, domainerrors.ErrSourceMismatch) {
		t.Fatalf("Bestätigung einer anderen Quelle: %v (Erwartung: ErrSourceMismatch)", err)
	}
}

// TestAcknowledgeRejectsUnregisteredConsumer trägt die
// Registrierungs-Vorbedingung der Bestätigung: der Fremdschlüssel der
// Position-Zeile endet den Aufruf sichtbar über die Klasse `storage`
// (`SPEC-001`, `SPEC-008`), nicht still.
func TestAcknowledgeRejectsUnregisteredConsumer(t *testing.T) {
	adapter, _ := newTestConsumerState(t)

	_, err := adapter.Acknowledge(context.Background(), model.ConsumerPosition{
		ConsumerID: testConsumer,
		Position:   consumerPosition(t, consumerStateSource, 100),
	})
	if !stderrors.Is(err, outbound.ErrStorage) {
		t.Fatalf("Bestätigung ohne Registrierung: %v (Erwartung: Klasse storage)", err)
	}
}

// TestConsumersCarryIndependentPositions trägt die Unabhängigkeit der
// Consumer (`LH-FA-CON-002`): die Bestätigungen des einen Consumers
// beflussen den Fortschritt des anderen nicht.
func TestConsumersCarryIndependentPositions(t *testing.T) {
	adapter, _ := newTestConsumerState(t)
	registerConsumer(t, adapter, testConsumer)
	registerConsumer(t, adapter, testConsumerOther)

	if _, err := adapter.Acknowledge(context.Background(), model.ConsumerPosition{
		ConsumerID: testConsumer,
		Position:   consumerPosition(t, consumerStateSource, 100),
	}); err != nil {
		t.Fatalf("Bestätigung %s: %v", testConsumer, err)
	}
	if _, err := adapter.Acknowledge(context.Background(), model.ConsumerPosition{
		ConsumerID: testConsumerOther,
		Position:   consumerPosition(t, consumerStateSource, 50),
	}); err != nil {
		t.Fatalf("Bestätigung %s: %v", testConsumerOther, err)
	}
	if got, ok := storedOffset(t, adapter, testConsumer); !ok || got != 100 {
		t.Fatalf("Fortschritt %s = %d (bestätigt: %v), wollen 100", testConsumer, got, ok)
	}
	if _, err := adapter.Acknowledge(context.Background(), model.ConsumerPosition{
		ConsumerID: testConsumerOther,
		Position:   consumerPosition(t, consumerStateSource, 400),
	}); err != nil {
		t.Fatalf("Bestätigung %s: %v", testConsumerOther, err)
	}
	if got, ok := storedOffset(t, adapter, testConsumerOther); !ok || got != 400 {
		t.Fatalf("Fortschritt %s = %d (bestätigt: %v), wollen 400", testConsumerOther, got, ok)
	}
	if got, ok := storedOffset(t, adapter, testConsumer); !ok || got != 100 {
		t.Fatalf("Fortschritt %s nach Bestätigung des anderen = %d (bestätigt: %v), wollen 100", testConsumer, got, ok)
	}
}

// TestPositionCarriesRestartPersistence trägt die Fortsetzung nach
// Neustart (`LH-FA-CON-005`): der bestätigte Stand liest sich über einen
// neuen Adapter-Verbindungs-Aufbau — derselbe Bestand, dieselbe Position.
func TestPositionCarriesRestartPersistence(t *testing.T) {
	dsn := os.Getenv("CDC_STORE_TEST_DSN")
	if dsn == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	first, pool := newTestConsumerState(t)
	registerConsumer(t, first, testConsumer)
	if _, err := first.Acknowledge(ctx, model.ConsumerPosition{
		ConsumerID: testConsumer,
		Position:   consumerPosition(t, consumerStateSource, 250),
	}); err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}
	first.Close()
	pool.Close()

	// Der Neustart baut die Verbindung neu auf; der Stand bleibt
	// (`LH-FA-CON-005` Happy Path).
	restarted, err := postgresstorage.NewConsumerState(ctx, dsn)
	if err != nil {
		t.Fatalf("NewConsumerState nach Neustart: %v", err)
	}
	t.Cleanup(restarted.Close)
	carried, err := restarted.Position(ctx, testConsumer)
	if err != nil {
		t.Fatalf("Position nach Neustart: %v", err)
	}
	if !carried.Acknowledged() || carried.Position.Offset != 250 {
		t.Fatalf("Position nach Neustart: %+v (Erwartung: 250)", carried)
	}
}

// TestPositionWithoutAcknowledgement trägt die definierte Anfangsposition
// (`LH-FA-CON-005` Boundary): ein Consumer ohne Bestätigung liest den
// Nullwert, ohne einen Fehler zu tragen.
func TestPositionWithoutAcknowledgement(t *testing.T) {
	adapter, _ := newTestConsumerState(t)
	registerConsumer(t, adapter, testConsumer)

	carried, err := adapter.Position(context.Background(), testConsumer)
	if err != nil {
		t.Fatalf("Position: %v", err)
	}
	if carried.Acknowledged() {
		t.Fatalf("Position ohne Bestätigung liest sich bestätigt: %+v", carried)
	}
}

// TestRemoveCarriesConsumerAndPosition trägt die administrative
// Entfernung (`LH-FA-CON-006`): der Consumer und seine bestätigte
// Position gehen gemeinsam — die dokumentierte Grenze —, ein erneuter
// Aufruf bleibt ohne Wirkung.
func TestRemoveCarriesConsumerAndPosition(t *testing.T) {
	adapter, pool := newTestConsumerState(t)
	registerConsumer(t, adapter, testConsumer)
	if _, err := adapter.Acknowledge(context.Background(), model.ConsumerPosition{
		ConsumerID: testConsumer,
		Position:   consumerPosition(t, consumerStateSource, 100),
	}); err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}

	removed, err := adapter.Remove(context.Background(), testConsumer)
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if !removed {
		t.Fatalf("Entfernung meldet Removed=false")
	}
	if _, ok := storedOffset(t, adapter, testConsumer); ok {
		t.Fatalf("bestätigte Position nach der Entfernung gelesen (LH-FA-CON-006 Boundary: sie geht mit)")
	}
	var count int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM cdc.consumer WHERE consumer_id = $1", testConsumer).Scan(&count); err != nil {
		t.Fatalf("Consumer-Zeilen: %v", err)
	}
	if count != 0 {
		t.Fatalf("Consumer-Zeilen nach der Entfernung = %d, wollen 0", count)
	}
	again, err := adapter.Remove(context.Background(), testConsumer)
	if err != nil {
		t.Fatalf("erneutes Remove: %v", err)
	}
	if again {
		t.Fatalf("Entfernung ohne Bestand meldet Removed=true (Idempotenz)")
	}
}

// Die Kennungs-Grenzen des Adapters enden vor dem ersten SQL-Aufruf über
// die Domänen-Invarianten (`ADR-0029`); ein Offset außerhalb des
// bigint-Bereichs endet über die PostgreSQL-Abbildung (`SPEC-003`).
func TestConsumerStateCarriesContractBounds(t *testing.T) {
	adapter, _ := newTestConsumerState(t)
	ctx := context.Background()

	if _, err := adapter.Position(ctx, ""); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Position ohne Kennung: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
	if _, err := adapter.Remove(ctx, ""); !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Remove ohne Kennung: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
	if _, err := adapter.Acknowledge(ctx, model.ConsumerPosition{ConsumerID: testConsumer}); !stderrors.Is(err, domainerrors.ErrInvalidPosition) {
		t.Fatalf("Bestätigung ohne Position: %v (Erwartung: ErrInvalidPosition)", err)
	}
	beyond, err := model.NewSourcePosition(consumerStateSource, math.MaxUint64)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	if _, err := adapter.Acknowledge(ctx, model.ConsumerPosition{
		ConsumerID: testConsumer,
		Position:   beyond,
	}); !stderrors.Is(err, mapper.ErrPositionOutOfRange) {
		t.Fatalf("Bestätigung außerhalb des bigint-Bereichs: %v (Erwartung: ErrPositionOutOfRange)", err)
	}
}

// Eine nicht erreichbare Instanz meldet der Aufbau als Fehler der Klasse
// `storage`; der Test braucht keine Datenbank (Port 1 verwirft lokal).
func TestNewConsumerStateCarriesStorageClass(t *testing.T) {
	_, err := postgresstorage.NewConsumerState(context.Background(), "postgres://cdc:cdc@127.0.0.1:1/cdc_test?sslmode=disable")
	if !stderrors.Is(err, outbound.ErrStorage) {
		t.Fatalf("Fehler = %v, wollen Klasse storage (%v)", err, outbound.ErrStorage)
	}
}
