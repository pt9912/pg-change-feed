package mapper_test

import (
	"math"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/mapper"
	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Zeilen-Übersetzung läuft ohne PostgreSQL-Instanz: die
// Positions-Abbildung (`SPEC-003`) und die Bild-Form (`SPEC-002`) sind
// reine Übersetzungsregeln und tragen ihre Grenzen ohne Treiber.

// Eine Position im bigint-Bereich trägt ihre Abbildung in beide
// Richtungen (`SPEC-003`).
func TestTransactionRowRoundTripsPosition(t *testing.T) {
	tx, err := model.NewOpenTransaction("t-1", "src-1")
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	position, err := model.NewSourcePosition("src-1", 12345)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	row, err := mapper.NewTransactionRow(tx, position)
	if err != nil {
		t.Fatalf("NewTransactionRow: %v", err)
	}
	if row.TransactionID != "t-1" || row.SourceID != "src-1" || row.CommitPosition != 12345 {
		t.Fatalf("Zeile = %+v, wollen t-1/src-1/12345", row)
	}
	restored, err := mapper.ToPosition(row.SourceID, row.CommitPosition)
	if err != nil {
		t.Fatalf("ToPosition: %v", err)
	}
	if restored != position {
		t.Fatalf("Abbildung = %+v, wollen %+v", restored, position)
	}
}

// Eine Position jenseits des bigint-Bereichs trägt keine
// Spalten-Abbildung; der Mapper meldet die Grenze sichtbar
// (`SPEC-003`, Adapter-Abbildung).
func TestTransactionRowRejectsPositionBeyondBigint(t *testing.T) {
	tx, err := model.NewOpenTransaction("t-1", "src-1")
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	position, err := model.NewSourcePosition("src-1", math.MaxInt64+1)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	if _, err := mapper.NewTransactionRow(tx, position); err != mapper.ErrPositionOutOfRange {
		t.Fatalf("Fehler = %v, wollen %v", err, mapper.ErrPositionOutOfRange)
	}
}

// Der Zeitstempel der Zeile trägt den realen Quell-Commit-Zeitpunkt
// (`LH-FA-ADM-004`) aus dem committed Domänenobjekt — nicht
// irgendeine Instanzzeit des Mappers selbst.
func TestTransactionRowCarriesSourceCommittedAt(t *testing.T) {
	tx, err := model.NewOpenTransaction("t-1", "src-1")
	if err != nil {
		t.Fatalf("NewOpenTransaction: %v", err)
	}
	position, err := model.NewSourcePosition("src-1", 100)
	if err != nil {
		t.Fatalf("NewSourcePosition: %v", err)
	}
	sourceCommittedAt := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	if err := tx.Commit(position, model.NewTimePoint(sourceCommittedAt.UnixNano())); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	row, err := mapper.NewTransactionRow(tx, position)
	if err != nil {
		t.Fatalf("NewTransactionRow: %v", err)
	}
	if !row.CommittedAt.Equal(sourceCommittedAt) {
		t.Fatalf("CommittedAt = %s, wollen den Quell-Commit-Zeitpunkt %s", row.CommittedAt, sourceCommittedAt)
	}
}

// Ein fehlendes Row Image geht als nil in die Spalte — Abwesenheit, kein
// Fehler (`LH-FA-CAP-008` Boundary).
func TestChangeRowsCarryMissingImageAsNil(t *testing.T) {
	change, err := model.NewChange("c-1", "t-1", "tbl-1", 1, model.OperationInsert, nil, nil, "sv-1")
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	rows, err := mapper.NewChangeRows([]model.Change{change})
	if err != nil {
		t.Fatalf("NewChangeRows: %v", err)
	}
	if rows[0].OldData != nil || rows[0].NewData != nil {
		t.Fatalf("Bilder = %v/%v, wollen nil/nil", rows[0].OldData, rows[0].NewData)
	}
}

// Ein nicht-leeres Bild geht als JSON-Text in die Spalte und liest sich
// über die Domänen-Konstruktoren unverändert zurück (`SPEC-002`).
func TestChangeRowRoundTripsImage(t *testing.T) {
	image := []byte(`{"x":1}`)
	change, err := model.NewChange("c-1", "t-1", "tbl-1", 1, model.OperationUpdate, image, image, "sv-1")
	if err != nil {
		t.Fatalf("NewChange: %v", err)
	}
	rows, err := mapper.NewChangeRows([]model.Change{change})
	if err != nil {
		t.Fatalf("NewChangeRows: %v", err)
	}
	restored, err := mapper.ToChange(rows[0])
	if err != nil {
		t.Fatalf("ToChange: %v", err)
	}
	if string(restored.OldImage) != string(image) || string(restored.NewImage) != string(image) {
		t.Fatalf("Bilder = %q/%q, wollen %q/%q", restored.OldImage, restored.NewImage, image, image)
	}
	if restored.Operation != model.OperationUpdate || restored.Sequence != 1 {
		t.Fatalf("Change = %+v, wollen UPDATE mit Sequenz 1", restored)
	}
}

// Eine Zeile unterhalb der Positionsschwelle verwirft der Mapper über die
// Domänen-Kante (`ADR-0029`); die CHECK-Kante der Spalte hält denselben
// Stand (`schema.sql`).
func TestToPositionRejectsNonPositiveColumn(t *testing.T) {
	if _, err := mapper.ToPosition("src-1", 0); err != domainerrors.ErrInvalidPosition {
		t.Fatalf("Fehler = %v, wollen %v", err, domainerrors.ErrInvalidPosition)
	}
}

// Der Lesepfad trägt die Klartext-Identität der Tabelle in den Change
// (`ADR-0081` Teilfrage 3): die Zeile trägt Schema- und Tabellenname aus
// dem Join auf `cdc.source_table`, `ToChange` setzt beide am Ergebnis — die
// opake `SourceTableID` bleibt daneben erhalten.
func TestToChangeCarriesSchemaAndTable(t *testing.T) {
	row := mapper.ChangeRow{
		ChangeID:      "c-1",
		TransactionID: "t-1",
		SourceTableID: "tbl-1",
		Sequence:      1,
		Operation:     string(model.OperationInsert),
		SchemaVersion: "sv-1",
		Schema:        "public",
		Table:         "orders",
	}
	change, err := mapper.ToChange(row)
	if err != nil {
		t.Fatalf("ToChange: %v", err)
	}
	if change.Schema != "public" || change.Table != "orders" {
		t.Fatalf("Schema/Table = %q/%q, wollen public/orders", change.Schema, change.Table)
	}
	if change.SourceTableID != "tbl-1" {
		t.Fatalf("SourceTableID = %q, wollen tbl-1", change.SourceTableID)
	}
}

// Ein fehlendes Bild und ein leeres Byte-Slice sind Abwesenheit und lesen
// sich beide als NULL; ein nicht-leeres Bild geht als JSON-Text in die
// `jsonb`-Spalte (`SPEC-002`, `LH-FA-CAP-008` Boundary).
func TestJSONImageCarriesAbsenceAndText(t *testing.T) {
	if got := mapper.JSONImage(nil); got != nil {
		t.Fatalf("JSONImage(nil) = %v, wollen nil (Abwesenheit)", got)
	}
	if got := mapper.JSONImage([]byte{}); got != nil {
		t.Fatalf("JSONImage(leeres Bild) = %v, wollen nil (kein gültiges JSON)", got)
	}
	image := []byte(`{"name":"a"}`)
	got, isText := mapper.JSONImage(image).(string)
	if !isText || got != string(image) {
		t.Fatalf("JSONImage(%q) = %v (%T), wollen den Text-Stand", image, got, got)
	}
}

// Eine Zeile außerhalb der Change-Invarianten endet über den
// Domänen-Konstruktor (`ADR-0029`): die Sequenz unter 1 trägt keinen Change,
// der Fehler kommt unverändert zurück.
func TestToChangeRejectsRowOutsideInvariants(t *testing.T) {
	row := mapper.ChangeRow{
		ChangeID:      "c-1",
		TransactionID: "t-1",
		SourceTableID: "tbl-1",
		Sequence:      0,
		Operation:     string(model.OperationInsert),
		SchemaVersion: "sv-1",
		Schema:        "public",
		Table:         "orders",
	}

	if _, err := mapper.ToChange(row); err != domainerrors.ErrNonPositiveSequence {
		t.Fatalf("Fehler = %v, wollen %v", err, domainerrors.ErrNonPositiveSequence)
	}
}
