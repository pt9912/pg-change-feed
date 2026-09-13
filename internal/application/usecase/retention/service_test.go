package retention_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/retention"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeStore trägt den `ChangeStorePort` als Fake (`ADR-0030`): `ReadChanges`
// meldet einen fest verdrahteten Bestand, `DeleteChanges` trägt die
// übergebene Menge zur Prüfung, ob sie genau die freigegebene ist.
type fakeStore struct {
	records      []outbound.ChangeRecord
	readCalls    int
	deleteCalled bool
	deletedIDs   []model.ChangeID
}

func (f *fakeStore) PersistTransaction(ctx context.Context, transaction *model.ChangeTransaction) error {
	return stderrors.New("Fake-Store trägt nur die Lese-/Löschseite")
}

func (f *fakeStore) ReadChanges(ctx context.Context, query outbound.ChangeQuery) ([]outbound.ChangeRecord, error) {
	f.readCalls++
	return f.records, nil
}

func (f *fakeStore) DeleteChanges(ctx context.Context, changeIDs []model.ChangeID) error {
	f.deleteCalled = true
	f.deletedIDs = changeIDs
	return nil
}

var _ outbound.ChangeStorePort = (*fakeStore)(nil)

// fakeState trägt den `ConsumerStatePort` als Fake (`ADR-0030`): `Positions`
// meldet einen fest verdrahteten Bestand, die übrigen Operationen tragen
// nur die Schnittstelle.
type fakeState struct {
	positions     []model.ConsumerPosition
	positionCalls int
}

func (f *fakeState) Register(ctx context.Context, consumer model.Consumer) (bool, error) {
	return true, nil
}

func (f *fakeState) Position(ctx context.Context, consumer model.ConsumerID) (model.ConsumerPosition, error) {
	return model.ConsumerPosition{}, nil
}

func (f *fakeState) Acknowledge(ctx context.Context, position model.ConsumerPosition) (model.ConsumerPosition, error) {
	return position, nil
}

func (f *fakeState) Remove(ctx context.Context, consumer model.ConsumerID) (bool, error) {
	return true, nil
}

func (f *fakeState) Positions(ctx context.Context, source model.SourceID) ([]model.ConsumerPosition, error) {
	f.positionCalls++
	return f.positions, nil
}

var _ outbound.ConsumerStatePort = (*fakeState)(nil)

// fakeClock trägt eine feste Wanduhr (`ADR-0040`); die Retention-Alters-
// Berechnung liest gegen sie.
type fakeClock struct {
	now model.TimePoint
}

func (f *fakeClock) Now() model.TimePoint {
	return f.now
}

var _ outbound.ClockPort = (*fakeClock)(nil)

// mustChange baut einen Change mit fester Kennung; die übrigen Felder
// tragen gültige, aber für diesen Test irrelevante Werte.
func mustChange(t *testing.T, id string) model.Change {
	t.Helper()
	change, err := model.NewChange(model.ChangeID(id), "t-1", "tbl-1", 1, model.OperationInsert, nil, []byte(`{}`), "sv-1")
	if err != nil {
		t.Fatalf("NewChange %s: %v", id, err)
	}
	return change
}

// mustPosition baut eine Quellposition an einem Offset.
func mustPosition(t *testing.T, offset uint64) model.SourcePosition {
	t.Helper()
	p, err := model.NewSourcePosition("src-1", offset)
	if err != nil {
		t.Fatalf("NewSourcePosition %d: %v", offset, err)
	}
	return p
}

// mustAckedConsumer baut eine bestätigte Consumer-Position an einem Offset.
func mustAckedConsumer(t *testing.T, id string, offset uint64) model.ConsumerPosition {
	t.Helper()
	consumer, err := model.NewConsumerPosition(model.ConsumerID(id))
	if err != nil {
		t.Fatalf("NewConsumerPosition %s: %v", id, err)
	}
	advanced, err := consumer.Advance(mustPosition(t, offset))
	if err != nil {
		t.Fatalf("Advance %s: %v", id, err)
	}
	return advanced
}

// TestRunDistinguishesEligibleChangesFromMixedSet trägt die Freigabe je
// betrachtetem Change (`LH-FA-RET-002`…`004`): eine gemischte Menge aus
// löschbar (alt genug, Consumer bereits vorbei), nicht löschbar wegen
// Alters (zu jung, obwohl der Consumer bereits vorbei ist) und nicht
// löschbar wegen eines zurückhängenden Consumers (alt genug, Consumer noch
// nicht so weit) — nur die freigegebene Teilmenge geht an
// `DeleteChanges`.
func TestRunDistinguishesEligibleChangesFromMixedSet(t *testing.T) {
	policy, err := model.NewRetentionPolicy(model.Duration{Nanos: 1000})
	if err != nil {
		t.Fatalf("NewRetentionPolicy: %v", err)
	}
	store := &fakeStore{records: []outbound.ChangeRecord{
		{ // löschbar: alt genug (Alter 2000 >= 1000), Consumer bereits vorbei (200 >= 100).
			Change:      mustChange(t, "c-eligible"),
			Position:    mustPosition(t, 100),
			CommittedAt: model.NewTimePoint(8000),
		},
		{ // nicht löschbar wegen Alters: Alter 500 < 1000, obwohl der Consumer bereits vorbei ist.
			Change:      mustChange(t, "c-too-young"),
			Position:    mustPosition(t, 150),
			CommittedAt: model.NewTimePoint(9500),
		},
		{ // nicht löschbar wegen eines zurückhängenden Consumers: alt genug, aber der Consumer hat die Position noch nicht bestätigt (200 < 300).
			Change:      mustChange(t, "c-consumer-behind"),
			Position:    mustPosition(t, 300),
			CommittedAt: model.NewTimePoint(8000),
		},
		{ // löschbar am Boundary: Alter exakt am Mindestalter, Consumer exakt an der Change-Position.
			Change:      mustChange(t, "c-boundary"),
			Position:    mustPosition(t, 200),
			CommittedAt: model.NewTimePoint(9000),
		},
	}}
	state := &fakeState{positions: []model.ConsumerPosition{mustAckedConsumer(t, "cons-1", 200)}}
	clock := &fakeClock{now: model.NewTimePoint(10000)}
	service := retention.NewRunRetentionService(store, state, clock)

	result, err := service.Run(context.Background(), retention.RunRetentionCommand{Source: "src-1", Policy: policy})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Deleted != 2 {
		t.Fatalf("Deleted = %d, wollen 2 (c-eligible, c-boundary)", result.Deleted)
	}
	if !store.deleteCalled {
		t.Fatal("DeleteChanges wurde nicht aufgerufen")
	}
	if len(store.deletedIDs) != 2 || store.deletedIDs[0] != "c-eligible" || store.deletedIDs[1] != "c-boundary" {
		t.Fatalf("freigegebene Menge = %v, wollen genau [c-eligible c-boundary]", store.deletedIDs)
	}
	if state.positionCalls != 1 {
		t.Fatalf("Positions-Lese-Züge = %d, wollen 1", state.positionCalls)
	}
}

// TestRunWithoutEligibleChangesCallsDeleteWithEmptySet trägt den Fall ohne
// Kandidaten: `DeleteChanges` geht mit einer leeren Menge — ein gültiger
// Aufruf ohne Wirkung, kein übersprungener Zug.
func TestRunWithoutEligibleChangesCallsDeleteWithEmptySet(t *testing.T) {
	policy, err := model.NewRetentionPolicy(model.Duration{Nanos: 1000})
	if err != nil {
		t.Fatalf("NewRetentionPolicy: %v", err)
	}
	store := &fakeStore{records: []outbound.ChangeRecord{
		{Change: mustChange(t, "c-too-young"), Position: mustPosition(t, 100), CommittedAt: model.NewTimePoint(9999)},
	}}
	state := &fakeState{}
	clock := &fakeClock{now: model.NewTimePoint(10000)}
	service := retention.NewRunRetentionService(store, state, clock)

	result, err := service.Run(context.Background(), retention.RunRetentionCommand{Source: "src-1", Policy: policy})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Deleted != 0 {
		t.Fatalf("Deleted = %d, wollen 0", result.Deleted)
	}
	if !store.deleteCalled || len(store.deletedIDs) != 0 {
		t.Fatalf("DeleteChanges-Aufruf = %v mit %v, wollen aufgerufen mit leerer Menge", store.deleteCalled, store.deletedIDs)
	}
}

// TestRunRejectsEmptySource trägt `LH-FA-RET-002` Negative: eine ungültige
// Konfiguration (leere Quellen-Kennung) endet über einen expliziten
// Fehlerpfad, ohne einen der Ports zu berühren — keine stille Übernahme.
func TestRunRejectsEmptySource(t *testing.T) {
	policy, err := model.NewRetentionPolicy(model.Duration{})
	if err != nil {
		t.Fatalf("NewRetentionPolicy: %v", err)
	}
	store := &fakeStore{}
	state := &fakeState{}
	clock := &fakeClock{now: model.NewTimePoint(0)}
	service := retention.NewRunRetentionService(store, state, clock)

	_, err = service.Run(context.Background(), retention.RunRetentionCommand{Policy: policy})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("leere Quelle: %v (Erwartung: ErrEmptyIdentifier)", err)
	}
	if store.readCalls != 0 || store.deleteCalled || state.positionCalls != 0 {
		t.Fatalf("ungültige Eingabe berührt Ports: readCalls=%d, deleteCalled=%v, positionCalls=%d", store.readCalls, store.deleteCalled, state.positionCalls)
	}
}
