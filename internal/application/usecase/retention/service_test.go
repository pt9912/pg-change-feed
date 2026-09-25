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

// fakeStore trägt den `ChangeStorePort` als Fake (`ADR-0030`):
// `ReadRetentionCandidates` meldet einen fest verdrahteten Bestand in der
// Reihenfolge der Liste, die der Fake als Ordnung des Schlüssels führt
// (`after` steht auf der Kennung, hinter der die Seite beginnt),
// `DeleteChanges` trägt die übergebene Menge zur Prüfung, ob sie genau die
// freigegebene ist. `maxPage` und `pageScript` begrenzen die Seite je Aufruf
// unter das übergebene `limit`; `readFailFrom` und `deleteFailFrom`
// (1-basiert) lassen erst die Aufrufe ab dem n-ten scheitern. Die
// `…ErrFor`-Felder binden einen Port-Fehler an genau den Eingabewert, auf
// den er antwortet (LP2): jede andere Quelle bzw. Menge trägt derselbe Fake.
type fakeStore struct {
	records      []outbound.ChangeRecord
	readCalls    int
	readAfters   []model.ChangeID // `after` je Lese-Aufruf
	readLimits   []int            // `limit` je Lese-Aufruf
	deleteCalled bool
	deletedIDs   []model.ChangeID   // Menge des letzten Lösch-Aufrufs
	deleteCalls  [][]model.ChangeID // Menge je Lösch-Aufruf, auch der scheiternden
	removed      []model.ChangeID   // Kennungen der gelungenen Lösch-Aufrufe

	maxPage        int   // größte Seite je Aufruf; 0 = nur `limit` begrenzt
	pageScript     []int // größte Seite je Aufruf n (0-basiert), sonst `maxPage`
	readFailFrom   int
	deleteFailFrom int
	failErr        error

	readErrFor   map[model.SourceID]error // Quell-Kennung → Fehler
	deleteErrFor model.ChangeID           // freigegebene Kennung → Fehler
	deleteErr    error
}

func (f *fakeStore) PersistTransaction(ctx context.Context, transaction *model.ChangeTransaction) error {
	return stderrors.New("Fake-Store trägt nur die Lese-/Löschseite")
}

func (f *fakeStore) ReadChanges(ctx context.Context, query outbound.ChangeQuery) ([]outbound.ChangeRecord, error) {
	return nil, stderrors.New("Fake-Store trägt für die Retention nur ReadRetentionCandidates")
}

func (f *fakeStore) ReadRetentionCandidates(ctx context.Context, source model.SourceID, after model.ChangeID, limit int) ([]outbound.RetentionCandidate, error) {
	call := f.readCalls
	f.readCalls++
	f.readAfters = append(f.readAfters, after)
	f.readLimits = append(f.readLimits, limit)
	if err, ok := f.readErrFor[source]; ok {
		return nil, err
	}
	if f.readFailFrom > 0 && f.readCalls >= f.readFailFrom {
		return nil, f.failErr
	}

	start := 0
	if after != "" {
		start = len(f.records)
		for i, record := range f.records {
			if record.Change.ID == after {
				start = i + 1
				break
			}
		}
	}
	size := limit
	if f.maxPage > 0 && f.maxPage < size {
		size = f.maxPage
	}
	if call < len(f.pageScript) && f.pageScript[call] < size {
		size = f.pageScript[call]
	}
	page := make([]outbound.RetentionCandidate, 0, size)
	for _, record := range f.records[start:] {
		if len(page) == size {
			break
		}
		page = append(page, outbound.RetentionCandidate{
			ChangeID:    record.Change.ID,
			Position:    record.Position,
			CommittedAt: record.CommittedAt,
		})
	}
	return page, nil
}

func (f *fakeStore) DeleteChanges(ctx context.Context, changeIDs []model.ChangeID) error {
	f.deleteCalled = true
	f.deletedIDs = changeIDs
	f.deleteCalls = append(f.deleteCalls, changeIDs)
	if f.deleteFailFrom > 0 && len(f.deleteCalls) >= f.deleteFailFrom {
		return f.failErr
	}
	for _, id := range changeIDs {
		if f.deleteErrFor != "" && id == f.deleteErrFor {
			return f.deleteErr
		}
	}
	f.removed = append(f.removed, changeIDs...)
	return nil
}

var _ outbound.ChangeStorePort = (*fakeStore)(nil)

// fakeState trägt den `ConsumerStatePort` als Fake (`ADR-0030`): `Positions`
// meldet einen fest verdrahteten Bestand, die übrigen Operationen tragen
// nur die Schnittstelle. `positionsErrFor` bindet den Port-Fehler an die
// gelesene Quell-Kennung (LP2).
type fakeState struct {
	positions     []model.ConsumerPosition
	positionCalls int

	positionsErrFor map[model.SourceID]error // Quell-Kennung → Fehler
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
	if err, ok := f.positionsErrFor[source]; ok {
		return nil, err
	}
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

// TestRunPositionsErrorFollowsSource trägt den Fehlerpfad des
// Positions-Lesens: ein Fehler von `Positions` wird unverändert
// durchgereicht, ohne Store-Zug. Die Ablehnung ist an die Quell-Kennung
// gebunden — derselbe Fake liest eine andere Quelle ohne Fehler.
func TestRunPositionsErrorFollowsSource(t *testing.T) {
	policy, err := model.NewRetentionPolicy(model.Duration{})
	if err != nil {
		t.Fatalf("NewRetentionPolicy: %v", err)
	}
	wantErr := stderrors.New("Consumer-Positionen nicht lesbar")
	store := &fakeStore{}
	state := &fakeState{positionsErrFor: map[model.SourceID]error{"src-kaputt": wantErr}}
	service := retention.NewRunRetentionService(store, state, &fakeClock{now: model.NewTimePoint(0)})

	_, err = service.Run(context.Background(), retention.RunRetentionCommand{Source: "src-kaputt", Policy: policy})
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if store.readCalls != 0 {
		t.Fatalf("Positions-Fehler trägt den Store-Lese-Zug: %d", store.readCalls)
	}

	if _, err := service.Run(context.Background(), retention.RunRetentionCommand{Source: "src-1", Policy: policy}); err != nil {
		t.Fatalf("Quelle src-1: %v", err)
	}
	if store.readCalls != 1 {
		t.Fatalf("Store-Lese-Züge: %d (Erwartung: 1)", store.readCalls)
	}
}

// TestRunStoreReadErrorFollowsSource trägt den Fehlerpfad des
// Kandidaten-Lesens: ein Fehler von `ReadRetentionCandidates` wird unverändert
// durchgereicht, ohne Lösch-Zug. Die Ablehnung ist an die Quell-Kennung der
// Abfrage gebunden — derselbe Fake liest eine andere Quelle ohne Fehler.
func TestRunStoreReadErrorFollowsSource(t *testing.T) {
	policy, err := model.NewRetentionPolicy(model.Duration{})
	if err != nil {
		t.Fatalf("NewRetentionPolicy: %v", err)
	}
	wantErr := stderrors.New("Change-Bestand nicht lesbar")
	store := &fakeStore{
		records:    []outbound.ChangeRecord{{Change: mustChange(t, "c-1"), Position: mustPosition(t, 100), CommittedAt: model.NewTimePoint(0)}},
		readErrFor: map[model.SourceID]error{"src-kaputt": wantErr},
	}
	state := &fakeState{}
	service := retention.NewRunRetentionService(store, state, &fakeClock{now: model.NewTimePoint(0)})

	_, err = service.Run(context.Background(), retention.RunRetentionCommand{Source: "src-kaputt", Policy: policy})
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if store.deleteCalled {
		t.Fatalf("Lese-Fehler trägt den Lösch-Zug: %v", store.deletedIDs)
	}

	if _, err := service.Run(context.Background(), retention.RunRetentionCommand{Source: "src-1", Policy: policy}); err != nil {
		t.Fatalf("Quelle src-1: %v", err)
	}
	if !store.deleteCalled {
		t.Fatal("Quelle src-1 ohne Lösch-Zug")
	}
}

// TestRunDeleteErrorFollowsEligibleSet trägt den Fehlerpfad der physischen
// Löschung: ein Fehler von `DeleteChanges` wird unverändert durchgereicht.
// Die Ablehnung ist an die **freigegebene Menge** gebunden — der Fake
// scheitert genau dann, wenn die Kennung in der übergebenen Menge liegt; ob
// sie das tut, entscheidet die Policy des Kommandos: dasselbe Kommando mit
// einem höheren Mindestalter gibt denselben Change nicht frei und die
// Löschung trägt durch.
func TestRunDeleteErrorFollowsEligibleSet(t *testing.T) {
	wantErr := stderrors.New("Löschung nicht ausführbar")
	store := &fakeStore{
		records:      []outbound.ChangeRecord{{Change: mustChange(t, "c-boom"), Position: mustPosition(t, 100), CommittedAt: model.NewTimePoint(8000)}},
		deleteErrFor: "c-boom",
		deleteErr:    wantErr,
	}
	state := &fakeState{}
	service := retention.NewRunRetentionService(store, state, &fakeClock{now: model.NewTimePoint(10000)})

	kleinesMindestalter, err := model.NewRetentionPolicy(model.Duration{Nanos: 1000})
	if err != nil {
		t.Fatalf("NewRetentionPolicy: %v", err)
	}
	_, err = service.Run(context.Background(), retention.RunRetentionCommand{Source: "src-1", Policy: kleinesMindestalter})
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Port-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if len(store.deletedIDs) != 1 || store.deletedIDs[0] != "c-boom" {
		t.Fatalf("freigegebene Menge = %v, wollen [c-boom]", store.deletedIDs)
	}

	grossesMindestalter, err := model.NewRetentionPolicy(model.Duration{Nanos: 5000})
	if err != nil {
		t.Fatalf("NewRetentionPolicy: %v", err)
	}
	result, err := service.Run(context.Background(), retention.RunRetentionCommand{Source: "src-1", Policy: grossesMindestalter})
	if err != nil {
		t.Fatalf("Mindestalter 5000: %v", err)
	}
	if result.Deleted != 0 || len(store.deletedIDs) != 0 {
		t.Fatalf("nicht freigegebener Change: Deleted=%d, Menge=%v (Erwartung: 0 und leer)", result.Deleted, store.deletedIDs)
	}
}

// pagedRecords ist der Bestand der Seiten-Tests in der Ordnung des
// Schlüssels: sieben Changes, davon vier freigegeben (c1, c4, c5, c7) bei
// Mindestalter 1000, Wanduhr 10000 und einem Consumer an Position 200 —
// c2 ist zu jung, c3 und c6 liegen hinter dem Consumer.
func pagedRecords(t *testing.T) []outbound.ChangeRecord {
	t.Helper()
	return []outbound.ChangeRecord{
		{Change: mustChange(t, "c1"), Position: mustPosition(t, 100), CommittedAt: model.NewTimePoint(8000)},
		{Change: mustChange(t, "c2"), Position: mustPosition(t, 150), CommittedAt: model.NewTimePoint(9500)},
		{Change: mustChange(t, "c3"), Position: mustPosition(t, 300), CommittedAt: model.NewTimePoint(8000)},
		{Change: mustChange(t, "c4"), Position: mustPosition(t, 200), CommittedAt: model.NewTimePoint(9000)},
		{Change: mustChange(t, "c5"), Position: mustPosition(t, 50), CommittedAt: model.NewTimePoint(1000)},
		{Change: mustChange(t, "c6"), Position: mustPosition(t, 250), CommittedAt: model.NewTimePoint(8000)},
		{Change: mustChange(t, "c7"), Position: mustPosition(t, 10), CommittedAt: model.NewTimePoint(5000)},
	}
}

// runPaged führt einen Lauf über `pagedRecords` gegen den Fake aus.
func runPaged(t *testing.T, store *fakeStore) (retention.RunRetentionResult, *fakeState, error) {
	t.Helper()
	policy, err := model.NewRetentionPolicy(model.Duration{Nanos: 1000})
	if err != nil {
		t.Fatalf("NewRetentionPolicy: %v", err)
	}
	store.records = pagedRecords(t)
	state := &fakeState{positions: []model.ConsumerPosition{mustAckedConsumer(t, "cons-1", 200)}}
	service := retention.NewRunRetentionService(store, state, &fakeClock{now: model.NewTimePoint(10000)})
	result, err := service.Run(context.Background(), retention.RunRetentionCommand{Source: "src-1", Policy: policy})
	return result, state, err
}

func idsEqual(got, want []model.ChangeID) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// TestRunReleasesSameSetAtEveryPageSize trägt die Mengen-Gleichheit über die
// Seitengrenzen (`ADR-0124`): bei Seiten zu 1, 2, 3, 7 (= N), 8 (> N)
// Kandidaten und bei einer Seite über alles gibt der Lauf genau die
// handgeschriebene Menge [c1 c4 c5 c7] frei und zählt sie über alle Seiten;
// die Zahl der Lese-Aufrufe ist die der nichtleeren Seiten plus die leere
// Endseite. Die Seitengröße kommt vom Fake, `PageSize` ist eine Konstante.
func TestRunReleasesSameSetAtEveryPageSize(t *testing.T) {
	want := []model.ChangeID{"c1", "c4", "c5", "c7"}
	cases := []struct {
		name      string
		maxPage   int
		wantReads int
	}{
		{"Seite 1", 1, 8},
		{"Seite 2", 2, 5},
		{"Seite 3", 3, 4},
		{"Seite N", 7, 2},
		{"Seite größer N", 8, 2},
		{"eine Seite über alles", 0, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeStore{maxPage: tc.maxPage}
			result, _, err := runPaged(t, store)
			if err != nil {
				t.Fatalf("Run: %v", err)
			}
			if !idsEqual(store.removed, want) {
				t.Fatalf("gelöschte Menge = %v, wollen %v", store.removed, want)
			}
			if result.Deleted != len(want) {
				t.Fatalf("Deleted = %d, wollen %d", result.Deleted, len(want))
			}
			if store.readCalls != tc.wantReads {
				t.Fatalf("Lese-Aufrufe = %d, wollen %d", store.readCalls, tc.wantReads)
			}
		})
	}
}

// TestRunReadsEachPageAtPageSizeFromTheLastKey trägt den Lese-Vertrag je
// Aufruf: jeder Aufruf trägt `limit` = `retention.PageSize`, `after` ist beim
// ersten leer und danach die letzte Kennung der vorigen nichtleeren Seite
// (handgeschrieben für Seiten zu 3: c3, c6, c7).
func TestRunReadsEachPageAtPageSizeFromTheLastKey(t *testing.T) {
	store := &fakeStore{maxPage: 3}
	if _, _, err := runPaged(t, store); err != nil {
		t.Fatalf("Run: %v", err)
	}
	wantAfters := []model.ChangeID{"", "c3", "c6", "c7"}
	if !idsEqual(store.readAfters, wantAfters) {
		t.Fatalf("after je Aufruf = %q, wollen %q", store.readAfters, wantAfters)
	}
	if retention.PageSize != 10000 {
		t.Fatalf("PageSize = %d, wollen 10000", retention.PageSize)
	}
	for i, limit := range store.readLimits {
		if limit != 10000 {
			t.Fatalf("limit des Aufrufs %d = %d, wollen 10000", i+1, limit)
		}
	}
}

// TestRunDeletesOnlyTheReleasedIDsOfEachPage trägt die Löschung je Seite:
// jeder `DeleteChanges`-Aufruf trägt nur Kennungen genau seiner Seite, eine
// Seite ohne Freigabe geht mit leerer Menge (handgeschrieben für Seiten zu 1
// und zu 3).
func TestRunDeletesOnlyTheReleasedIDsOfEachPage(t *testing.T) {
	cases := []struct {
		name    string
		maxPage int
		want    [][]model.ChangeID
	}{
		{"Seiten zu 3", 3, [][]model.ChangeID{{"c1"}, {"c4", "c5"}, {"c7"}}},
		{"Seiten zu 1", 1, [][]model.ChangeID{{"c1"}, {}, {}, {"c4"}, {"c5"}, {}, {"c7"}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeStore{maxPage: tc.maxPage}
			if _, _, err := runPaged(t, store); err != nil {
				t.Fatalf("Run: %v", err)
			}
			if len(store.deleteCalls) != len(tc.want) {
				t.Fatalf("Lösch-Aufrufe = %v, wollen %v", store.deleteCalls, tc.want)
			}
			for i := range tc.want {
				if !idsEqual(store.deleteCalls[i], tc.want[i]) {
					t.Fatalf("Lösch-Aufruf %d = %v, wollen %v", i+1, store.deleteCalls[i], tc.want[i])
				}
			}
		})
	}
}

// TestRunReadsPositionsOncePerRun trägt den Ablauf von `ADR-0124`
// Festlegung 4: über acht Lese-Aufrufe liest der Lauf die
// Consumer-Positionen einmal.
func TestRunReadsPositionsOncePerRun(t *testing.T) {
	store := &fakeStore{maxPage: 1}
	_, state, err := runPaged(t, store)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if store.readCalls != 8 || state.positionCalls != 1 {
		t.Fatalf("Lese-Aufrufe = %d, Positions-Lese-Züge = %d, wollen 8 und 1", store.readCalls, state.positionCalls)
	}
}

// TestRunReadFailureLeavesPagesBefore trägt den Abbruch zwischen Seiten
// (`ADR-0124` Festlegung 6, „ab Aufruf n“): scheitert das Lesen ab dem
// dritten Aufruf, bleiben die Löschungen der zwei Seiten davor bestehen, der
// Lauf liefert den Fehler und ein leeres Ergebnis.
func TestRunReadFailureLeavesPagesBefore(t *testing.T) {
	wantErr := stderrors.New("Lesen ab Aufruf 3 nicht möglich")
	store := &fakeStore{maxPage: 3, readFailFrom: 3, failErr: wantErr}
	result, _, err := runPaged(t, store)
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Lese-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if want := []model.ChangeID{"c1", "c4", "c5"}; !idsEqual(store.removed, want) {
		t.Fatalf("gelöschte Menge = %v, wollen genau die Seiten davor %v", store.removed, want)
	}
	if store.readCalls != 3 || result.Deleted != 0 {
		t.Fatalf("Lese-Aufrufe = %d, Deleted = %d, wollen 3 und 0", store.readCalls, result.Deleted)
	}
}

// TestRunDeleteFailureStopsAtThatPage trägt dieselbe Grenze für die
// Löschung: scheitert `DeleteChanges` ab dem zweiten Aufruf, bleibt die
// Löschung der ersten Seite bestehen, der Lauf liest keine weitere Seite und
// liefert den Fehler.
func TestRunDeleteFailureStopsAtThatPage(t *testing.T) {
	wantErr := stderrors.New("Löschen ab Aufruf 2 nicht möglich")
	store := &fakeStore{maxPage: 3, deleteFailFrom: 2, failErr: wantErr}
	result, _, err := runPaged(t, store)
	if !stderrors.Is(err, wantErr) {
		t.Fatalf("Lösch-Fehler: %v (Erwartung: %v)", err, wantErr)
	}
	if want := []model.ChangeID{"c1"}; !idsEqual(store.removed, want) {
		t.Fatalf("gelöschte Menge = %v, wollen genau die Seite davor %v", store.removed, want)
	}
	if len(store.deleteCalls) != 2 || store.readCalls != 2 || result.Deleted != 0 {
		t.Fatalf("Lösch-Aufrufe = %d, Lese-Aufrufe = %d, Deleted = %d, wollen 2, 2 und 0", len(store.deleteCalls), store.readCalls, result.Deleted)
	}
}

// TestRunContinuesPastShortPages trägt: eine Seite, die kürzer ist als
// `limit`, ist nicht das Ende — nur die leere Seite beendet den Lauf. Der
// Fake liefert Seiten zu 2, 1 und 4 Kandidaten, der Lauf erreicht c7.
func TestRunContinuesPastShortPages(t *testing.T) {
	store := &fakeStore{pageScript: []int{2, 1, 4}}
	result, _, err := runPaged(t, store)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if want := []model.ChangeID{"c1", "c4", "c5", "c7"}; !idsEqual(store.removed, want) || result.Deleted != 4 {
		t.Fatalf("gelöschte Menge = %v (Deleted %d), wollen %v", store.removed, result.Deleted, want)
	}
	if wantAfters := []model.ChangeID{"", "c2", "c3", "c7"}; !idsEqual(store.readAfters, wantAfters) {
		t.Fatalf("after je Aufruf = %q, wollen %q", store.readAfters, wantAfters)
	}
}
