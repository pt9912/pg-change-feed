package readchanges_test

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/readchanges"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeChangeStore trägt den `ChangeStorePort` als Fälschung (`ADR-0030`)
// und führt dieselbe Kontrakt-Prüfung, die der reale Lesepfad am Port
// führt (`outbound.ChangeQuery.Validate`, `PostgresChangeStoreAdapter.
// ReadChanges`): der Fake beobachtet die übersetzte Abfrage und lässt die
// Port-Sentinels damit real entstehen — die Hälfte des Nachweises, dass
// `limit < 1` und `from > to` über den Use Case ankommen.
type fakeChangeStore struct {
	records []outbound.ChangeRecord
	err     error
	queries []outbound.ChangeQuery
}

func (f *fakeChangeStore) PersistTransaction(context.Context, *model.ChangeTransaction) error {
	return nil
}

func (f *fakeChangeStore) ReadChanges(_ context.Context, query outbound.ChangeQuery) ([]outbound.ChangeRecord, error) {
	f.queries = append(f.queries, query)
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return f.records, f.err
}

func (f *fakeChangeStore) ReadRetentionCandidates(context.Context, model.SourceID, model.ChangeID, int) ([]outbound.RetentionCandidate, error) {
	return nil, nil
}

func (f *fakeChangeStore) DeleteChanges(context.Context, []model.ChangeID) error {
	return nil
}

func readRecord(t *testing.T, id string, offset uint64, schema, table string) outbound.ChangeRecord {
	t.Helper()
	position, err := model.NewSourcePosition("src-1", offset)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	change, err := model.NewChange(
		model.ChangeID(id), model.TransactionID("tx-1"), model.SourceTableID("tbl-1"),
		1, model.OperationInsert, nil, []byte(`{"n":1}`), model.SchemaVersionID("sv-1"),
	)
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	change.Schema = schema
	change.Table = table
	return outbound.ChangeRecord{
		Position:    position,
		Change:      change,
		CommittedAt: model.NewTimePoint(42),
	}
}

// TestReadChangesTranslatesQueryToPort trägt den Kern der Übersetzung
// (`ADR-0081` Teilfrage 1): der Use Case ruft `ChangeStorePort.ReadChanges`
// mit der übersetzten Abfrage auf — dieselbe Delegation, dieselbe
// Filterachse (Quelle/Schema/Tabelle), derselbe Bereich, dasselbe Limit.
// Der Test ist der Regressionstest gegen einen zweiten Lesepfad.
func TestReadChangesTranslatesQueryToPort(t *testing.T) {
	start, err := model.NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	end, err := model.NewSourcePosition("src-1", 200)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	limit := 5
	store := &fakeChangeStore{records: []outbound.ChangeRecord{readRecord(t, "c-1", 150, "public", "orders")}}

	result, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{
		Source: "src-1",
		Schema: "public",
		Table:  "orders",
		Start:  &start,
		End:    &end,
		Limit:  &limit,
	})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if len(store.queries) != 1 {
		t.Fatalf("Port-Aufrufe = %d, wollen 1", len(store.queries))
	}
	got := store.queries[0]
	if got.Source != "src-1" || got.Schema != "public" || got.Table != "orders" {
		t.Fatalf("Abfrage = %+v, wollen Quelle public/orders", got)
	}
	if got.Start == nil || got.End == nil || got.Start.Offset != 100 || got.End.Offset != 200 {
		t.Fatalf("Bereich = %+v, wollen [100,200)", got)
	}
	if got.Limit == nil || *got.Limit != 5 {
		t.Fatalf("Limit = %v, wollen 5", got.Limit)
	}
	if len(result.Changes) != 1 {
		t.Fatalf("Changes = %+v, wollen einen Eintrag", result.Changes)
	}
	entry := result.Changes[0]
	if entry.Position.Offset != 150 || entry.Change.ID != "c-1" || entry.Change.Schema != "public" || entry.Change.Table != "orders" {
		t.Fatalf("Eintrag = %+v, wollen den gelesenen Change samt Identität", entry)
	}
	if entry.CommittedAt.UnixNanos != 42 {
		t.Fatalf("CommittedAt = %d, wollen 42", entry.CommittedAt.UnixNanos)
	}
}

// Ohne Schema/Tabelle ist die Filterachse leer, ohne Limit liest der Aufruf
// unbegrenzt: der Use Case setzt **kein** Default-Limit (`ADR-0081`
// Teilfrage 2) — er normiert nichts.
func TestReadChangesLeavesOptionalFiltersUnset(t *testing.T) {
	store := &fakeChangeStore{}

	if _, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{Source: "src-1"}); err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if got := store.queries[0]; got.Schema != "" || got.Table != "" || got.Limit != nil || got.Start != nil || got.End != nil {
		t.Fatalf("Abfrage = %+v, wollen leere Filter, kein Limit, kein Bereich", got)
	}
}

// Die Bereichs-Grenzen des Port-Kontrakts kommen unverändert zurück
// (`LH-FA-REA-001` Negative): ein invertierter Bereich endet über
// `ErrRangeInverted`, der Use Case normiert ihn nicht zu einer leeren
// Menge.
func TestReadChangesCarriesInvertedRange(t *testing.T) {
	start, err := model.NewSourcePosition("src-1", 300)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	end, err := model.NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	store := &fakeChangeStore{}

	_, err = readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{
		Source: "src-1",
		Start:  &start,
		End:    &end,
	})
	if !stderrors.Is(err, outbound.ErrRangeInverted) {
		t.Fatalf("Fehler = %v, wollen %v", err, outbound.ErrRangeInverted)
	}
}

// Die Limit-Grenze des Port-Kontrakts kommt unverändert zurück
// (`LH-FA-REA-003` Negative): ein Limit unter 1 endet über
// `ErrNonPositiveLimit`.
func TestReadChangesCarriesNonPositiveLimit(t *testing.T) {
	zero := 0
	store := &fakeChangeStore{}

	_, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{
		Source: "src-1",
		Limit:  &zero,
	})
	if !stderrors.Is(err, outbound.ErrNonPositiveLimit) {
		t.Fatalf("Fehler = %v, wollen %v", err, outbound.ErrNonPositiveLimit)
	}
}

// Der Leerfall trägt eine leere, **gesetzte** Liste (`LH-FA-REA-006`
// Boundary): die Antwortform `{"changes": []}` setzt das voraus, nicht
// `null`.
func TestReadChangesEmptyResultIsSet(t *testing.T) {
	store := &fakeChangeStore{}

	result, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{Source: "src-1"})
	if err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if result.Changes == nil || len(result.Changes) != 0 {
		t.Fatalf("Changes = %v, wollen leere, gesetzte Liste", result.Changes)
	}
}

// Eine leere Quellen-Kennung endet als ungültige Eingabe, nicht als stilles
// Lesen über alle Quellen; der Port wird nicht berührt.
func TestReadChangesRejectsMissingSource(t *testing.T) {
	store := &fakeChangeStore{}

	_, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fehler = %v, wollen %v", err, domainerrors.ErrEmptyIdentifier)
	}
	if len(store.queries) != 0 {
		t.Fatalf("Port-Aufrufe = %d, wollen 0", len(store.queries))
	}
}

// Ein Port-Fehler kommt unverändert zurück — der Use Case verdeckt keine
// Klasse (`SPEC-008`).
func TestReadChangesCarriesPortFailure(t *testing.T) {
	cause := stderrors.New("Verbindung abgelehnt")
	store := &fakeChangeStore{err: cause}

	_, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{Source: "src-1"})
	if !stderrors.Is(err, cause) {
		t.Fatalf("Fehler = %v, wollen die Ursache", err)
	}
}

// Das Ziel geht unverändert in die Abfrage des Ports (`ADR-0138`
// Festlegung 1): der Use Case reicht den Filter einmal an den Store, neben
// Schema und Tabelle.
// Rot färbende Mutation: `Target: query.Target` aus der Abbildung streichen
// oder durch eine feste Zeichenkette ersetzen.
func TestReadChangesTranslatesTarget(t *testing.T) {
	store := &fakeChangeStore{}

	if _, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{
		Source: "src-1", Schema: "public", Table: "orders", Target: "eu",
	}); err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if len(store.queries) != 1 {
		t.Fatalf("Port-Aufrufe = %d, wollen 1", len(store.queries))
	}
	if got := store.queries[0]; got.Target != "eu" || got.Schema != "public" || got.Table != "orders" {
		t.Fatalf("Abfrage = %+v, wollen Ziel eu neben public/orders", got)
	}
}

// Ohne Ziel bleibt die Filterachse des Ziels leer (Regression).
func TestReadChangesLeavesTargetUnset(t *testing.T) {
	store := &fakeChangeStore{}

	if _, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{Source: "src-1"}); err != nil {
		t.Fatalf("ReadChanges: %v", err)
	}
	if got := store.queries[0].Target; got != "" {
		t.Fatalf("Ziel = %q, wollen leer", got)
	}
}

// failingStore schlägt bei jedem Aufruf von `ReadChanges` fehl und zählt
// ihn: ein Aufruf des Ports ist in den Tests der Alphabet-Prüfung der Fehler.
type failingStore struct {
	fakeChangeStore
	calls int
}

func (f *failingStore) ReadChanges(context.Context, outbound.ChangeQuery) ([]outbound.ChangeRecord, error) {
	f.calls++
	return nil, stderrors.New("der Store darf nicht angefragt werden")
}

// Ein gesetztes Ziel außerhalb des Alphabets des Zielnamens (`ADR-0139`
// Festlegung 2) liefert die leere, gesetzte Liste ohne Fehler, und der Port
// wird nicht angefragt (der Store-Fake endet bei jedem Aufruf mit einem
// Fehler und zählt ihn). Die Eingabe ist das Ziel: Großbuchstabe, 64 Zeichen, U+0000, Trennzeichen am
// Anfang, Leerzeichen, Punkt.
// Rot färbende Mutation: die Alphabet-Prüfung streichen (der Store-Fehler
// kommt zurück, der Port-Zähler steigt) oder `!model.IsValidRouteTarget`
// zu `model.IsValidRouteTarget` kehren (die Fälle laufen in den Store).
func TestReadChangesTargetOutsideAlphabetAnswersEmptyWithoutStore(t *testing.T) {
	cases := []struct {
		name   string
		target string
	}{
		{"Großbuchstabe", "Eu"},
		{"64 Zeichen", strings.Repeat("a", 64)},
		{"U+0000 am Ende", "eu\x00"},
		{"nur U+0000", "\x00"},
		{"Bindestrich am Anfang", "-eu"},
		{"Leerzeichen", "e u"},
		{"Punkt", "eu.de"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &failingStore{}

			result, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{
				Source: "src-1", Target: tc.target,
			})
			if err != nil {
				t.Fatalf("Fehler = %v, wollen keinen", err)
			}
			if result.Changes == nil || len(result.Changes) != 0 {
				t.Fatalf("Changes = %v, wollen leere, gesetzte Liste", result.Changes)
			}
			if store.calls != 0 {
				t.Fatalf("Port-Aufrufe = %d, wollen 0", store.calls)
			}
		})
	}
}

// Gegenprobe zur Alphabet-Prüfung: ein Ziel im Alphabet (auch an der Grenze
// von 63 Zeichen) erreicht den Port — die Prüfung weist nicht jedes Ziel ab.
func TestReadChangesTargetInsideAlphabetReachesStore(t *testing.T) {
	for _, target := range []string{"eu", "a", "0-9_z", strings.Repeat("a", 63)} {
		store := &failingStore{}

		_, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), inbound.ReadChangesQuery{
			Source: "src-1", Target: target,
		})
		if err == nil || store.calls != 1 {
			t.Fatalf("Ziel %q: Fehler = %v, Port-Aufrufe = %d, wollen den Store-Fehler nach einem Aufruf", target, err, store.calls)
		}
	}
}

// Eine leere Quelle endet vor der Alphabet-Prüfung als ungültige Eingabe:
// ein ungültiges Ziel verdeckt den Pflichtfeld-Fehler nicht.
func TestReadChangesMissingSourceWinsOverInvalidTarget(t *testing.T) {
	_, err := readchanges.NewReadChangesService(&failingStore{}).ReadChanges(context.Background(), inbound.ReadChangesQuery{Target: "EU"})
	if !stderrors.Is(err, domainerrors.ErrEmptyIdentifier) {
		t.Fatalf("Fehler = %v, wollen %v", err, domainerrors.ErrEmptyIdentifier)
	}
}

// Ein Fehler des Lese-Kontrakts hat Vorrang vor der Antwort „leer" für ein
// Ziel außerhalb des Alphabets: dieselbe fehlerhafte Anfrage endet mit
// demselben Fehler wie ohne Ziel (`ADR-0139`), und der Port wird nicht
// angefragt.
// Rot färbende Mutation: die Alphabet-Prüfung vor die Validierung setzen
// (Rückgabe der leeren Liste ohne `Validate`).
func TestReadChangesContractErrorsWinOverInvalidTarget(t *testing.T) {
	position := func(source model.SourceID, offset uint64) *model.SourcePosition {
		p, err := model.NewSourcePosition(source, offset)
		if err != nil {
			t.Fatalf("NewSourcePosition: %v", err)
		}
		return &p
	}
	zero := 0
	cases := []struct {
		name  string
		query inbound.ReadChangesQuery
		want  error
	}{
		{"Limit unter 1", inbound.ReadChangesQuery{Source: "src-1", Limit: &zero}, outbound.ErrNonPositiveLimit},
		{"invertierter Bereich", inbound.ReadChangesQuery{Source: "src-1", Start: position("src-1", 300), End: position("src-1", 100)}, outbound.ErrRangeInverted},
		{"Start einer anderen Quelle", inbound.ReadChangesQuery{Source: "src-1", Start: position("src-2", 5)}, domainerrors.ErrSourceMismatch},
	}
	for _, tc := range cases {
		for _, target := range []string{"EU", ""} {
			t.Run(tc.name+" Ziel "+target, func(t *testing.T) {
				store := &fakeChangeStore{}
				query := tc.query
				query.Target = target

				_, err := readchanges.NewReadChangesService(store).ReadChanges(context.Background(), query)
				if !stderrors.Is(err, tc.want) {
					t.Fatalf("Fehler = %v, wollen %v", err, tc.want)
				}
				wantCalls := 1 // ohne Ziel erzeugt der Port-Fake den Fehler
				if target != "" {
					wantCalls = 0
				}
				if len(store.queries) != wantCalls {
					t.Fatalf("Port-Aufrufe = %d, wollen %d", len(store.queries), wantCalls)
				}
			})
		}
	}
}
