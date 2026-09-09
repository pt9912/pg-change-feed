package outbound_test

import (
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Abfrage-Validierung trägt die Grenzen des Lese-Vertrags am Port
// (`ADR-0042`): Bereich (`LH-FA-REA-001`), Limit (`LH-FA-REA-003`) und
// Quellen-Ordnung (`ADR-0005`) enden hier als sichtbare Fehler, bevor ein
// Adapter sie in SQL übersetzt.

func position(source string, offset uint64) *model.SourcePosition {
	p, err := model.NewSourcePosition(model.SourceID(source), offset)
	if err != nil {
		panic(err)
	}
	return &p
}

// Unbegrenztes Lesen über Nil-Felder trägt die valide Abfrage: ohne
// Start, Ende, Limit und Filter liest der Aufruf den vollen Bestand der
// Quelle (`LH-FA-REA-002`).
func TestChangeQueryValidateAcceptsOpenQuery(t *testing.T) {
	query := outbound.ChangeQuery{Source: "src-1"}
	if err := query.Validate(); err != nil {
		t.Fatalf("Validate = %v, wollen keinen Fehler", err)
	}
}

// Ein Endpunkt vor dem Startpunkt endet am Port-Kontrakt (`LH-FA-REA-001`
// Negative: expliziter Fehlerpfad statt unbestimmter leerer Bereich).
func TestChangeQueryValidateRejectsInvertedRange(t *testing.T) {
	query := outbound.ChangeQuery{
		Source: "src-1",
		Start:  position("src-1", 200),
		End:    position("src-1", 100),
	}
	if err := query.Validate(); err != outbound.ErrRangeInverted {
		t.Fatalf("Validate = %v, wollen %v", err, outbound.ErrRangeInverted)
	}
}

// Ein gesetztes Limit unter 1 endet am Port-Kontrakt (`LH-FA-REA-003`
// Negative: expliziter Fehlerpfad statt Normierung).
func TestChangeQueryValidateRejectsNonPositiveLimit(t *testing.T) {
	zero := 0
	negative := -1
	for _, limit := range []*int{&zero, &negative} {
		query := outbound.ChangeQuery{Source: "src-1", Limit: limit}
		if err := query.Validate(); err != outbound.ErrNonPositiveLimit {
			t.Fatalf("Validate mit Limit %d = %v, wollen %v", *limit, err, outbound.ErrNonPositiveLimit)
		}
	}
}

// Positionen einer anderen Quelle tragen keine fachliche Ordnung für
// diese Abfrage (`ADR-0005`); der Port-Kontrakt verwirft sie.
func TestChangeQueryValidateRejectsForeignSource(t *testing.T) {
	start := outbound.ChangeQuery{Source: "src-1", Start: position("src-2", 100)}
	if err := start.Validate(); err != domainerrors.ErrSourceMismatch {
		t.Fatalf("Validate (Start) = %v, wollen %v", err, domainerrors.ErrSourceMismatch)
	}
	end := outbound.ChangeQuery{Source: "src-1", End: position("src-2", 100)}
	if err := end.Validate(); err != domainerrors.ErrSourceMismatch {
		t.Fatalf("Validate (End) = %v, wollen %v", err, domainerrors.ErrSourceMismatch)
	}
}

// Eine Abfrage ohne Quelle trägt keine Ordnungs-Größe; der Kontrakt
// verwirft sie über die Domänen-Kennung.
func TestChangeQueryValidateRejectsMissingSource(t *testing.T) {
	query := outbound.ChangeQuery{}
	if err := query.Validate(); err != domainerrors.ErrEmptyIdentifier {
		t.Fatalf("Validate = %v, wollen %v", err, domainerrors.ErrEmptyIdentifier)
	}
}
