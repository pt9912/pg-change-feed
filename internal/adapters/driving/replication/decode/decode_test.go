package decode_test

import (
	stderrors "errors"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
)

// Die Dekodier-Tests tragen die `pgoutput`-Binärcodes in Handform
// (`ADR-0008`): BEGIN, COMMIT, Relation, Insert, Update, Delete,
// TRUNCATE und die Tupel-Typen werden als bekannte Byte-Stände gefüttert
// — kein Netz, kein Treiber (`ADR-0030`). Die Binärform folgt der
// Logical-Replication-Protokoll-Form des Output-Plugins.

// testColumn trägt eine Relation-Spalte für relationPayload.
type testColumn struct {
	name  string
	flags byte
}

// tupleValue trägt einen Einzelwert für textTuple.
type tupleValue struct {
	text  *string
	toast bool
}

// text trägt einen Text-Wert für textTuple.
func text(value string) tupleValue { return tupleValue{text: &value} }

// null trägt einen NULL-Wert für textTuple.
func null() tupleValue { return tupleValue{} }

// unchangedToast trägt ein unverändertes TOAST-Stützwert-Platzhalterverhalten.
func unchangedToast() tupleValue { return tupleValue{toast: true} }

// beginPayload trägt eine BEGIN-Nachricht (`B`): FinalLSN,
// Commit-Zeitstempel, XID.
func beginPayload(finalLSN uint64, xid uint32) []byte {
	payload := []byte{'B'}
	payload = appendLSN(payload, finalLSN)
	payload = appendMicros(payload, 0)
	return appendUint32(payload, xid)
}

// commitPayload trägt eine COMMIT-Nachricht (`C`): Flags, Commit-LSN,
// End-LSN, Commit-Zeitstempel.
func commitPayload(commitLSN, endLSN uint64) []byte {
	payload := []byte{'C', 0}
	payload = appendLSN(payload, commitLSN)
	payload = appendLSN(payload, endLSN)
	return appendMicros(payload, 0)
}

// relationPayload trägt eine Relation-Nachricht (`R`): Relation-ID,
// Namensraum, Tabellenname, Replica Identity und Spalten.
func relationPayload(oid uint32, namespace, name string, columns ...testColumn) []byte {
	payload := []byte{'R'}
	payload = appendUint32(payload, oid)
	payload = appendString(payload, namespace)
	payload = appendString(payload, name)
	payload = append(payload, 0) // Replica Identity: DEFAULT
	payload = appendUint16(payload, uint16(len(columns)))
	for _, column := range columns {
		payload = append(payload, column.flags)
		payload = appendString(payload, column.name)
		payload = appendUint32(payload, 25) // Typ-OID (text), für die Dekodierung ohne Bedeutung
		payload = appendUint32(payload, 0xFFFFFFFF)
	}
	return payload
}

// insertPayload trägt eine Insert-Nachricht (`I`): Relation-ID,
// Neu-Tupel.
func insertPayload(oid uint32, tuple []byte) []byte {
	payload := appendUint32([]byte{'I'}, oid)
	payload = append(payload, 'N')
	return append(payload, tuple...)
}

// updatePayloadKey trägt eine Update-Nachricht (`U`) mit
// Schlüssel-Tupel (`K`) und Neu-Tupel.
func updatePayloadKey(oid uint32, oldTuple, newTuple []byte) []byte {
	payload := appendUint32([]byte{'U'}, oid)
	payload = append(payload, 'K')
	payload = append(payload, oldTuple...)
	payload = append(payload, 'N')
	return append(payload, newTuple...)
}

// updatePayloadFull trägt eine Update-Nachricht (`U`) mit vollem
// Alt-Tupel (`O`) und Neu-Tupel.
func updatePayloadFull(oid uint32, oldTuple, newTuple []byte) []byte {
	payload := appendUint32([]byte{'U'}, oid)
	payload = append(payload, 'O')
	payload = append(payload, oldTuple...)
	payload = append(payload, 'N')
	return append(payload, newTuple...)
}

// deletePayloadKey trägt eine Delete-Nachricht (`D`) mit
// Schlüssel-Tupel (`K`).
func deletePayloadKey(oid uint32, oldTuple []byte) []byte {
	payload := appendUint32([]byte{'D'}, oid)
	payload = append(payload, 'K')
	return append(payload, oldTuple...)
}

// deletePayloadFull trägt eine Delete-Nachricht (`D`) mit vollem
// Alt-Tupel (`O`).
func deletePayloadFull(oid uint32, oldTuple []byte) []byte {
	payload := appendUint32([]byte{'D'}, oid)
	payload = append(payload, 'O')
	return append(payload, oldTuple...)
}

// truncatePayload trägt eine Truncate-Nachricht (`T`): Relation-Anzahl,
// Optionen, Relation-IDs.
func truncatePayload(oids ...uint32) []byte {
	payload := appendUint32([]byte{'T'}, uint32(len(oids)))
	payload = append(payload, 0)
	for _, oid := range oids {
		payload = appendUint32(payload, oid)
	}
	return payload
}

// textTuple trägt ein Tupel (`TupleData`): Spaltenanzahl und je Spalte
// der Typ — Text mit Länge und Bytes, NULL, unverändertes TOAST.
func textTuple(values ...tupleValue) []byte {
	tuple := appendUint16(nil, uint16(len(values)))
	for _, value := range values {
		switch {
		case value.text != nil:
			tuple = append(tuple, 't')
			tuple = appendUint32(tuple, uint32(len(*value.text)))
			tuple = append(tuple, []byte(*value.text)...)
		case value.toast:
			tuple = append(tuple, 'u')
		default:
			tuple = append(tuple, 'n')
		}
	}
	return tuple
}

// appendString trägt ein C-String-Feld.
func appendString(payload []byte, value string) []byte {
	return append(append(payload, value...), 0)
}

// appendLSN trägt einen LSN im Big-Endian-Format.
func appendLSN(payload []byte, lsn uint64) []byte {
	raw := make([]byte, 8)
	for i := range raw {
		raw[7-i] = byte(lsn >> (8 * uint(i)))
	}
	return append(payload, raw...)
}

// appendUint32 trägt ein uint32 im Big-Endian-Format.
func appendUint32(payload []byte, value uint32) []byte {
	raw := make([]byte, 4)
	for i := range raw {
		raw[3-i] = byte(value >> (8 * uint(i)))
	}
	return append(payload, raw...)
}

// appendUint16 trägt ein uint16 im Big-Endian-Format.
func appendUint16(payload []byte, value uint16) []byte {
	raw := make([]byte, 2)
	for i := range raw {
		raw[1-i] = byte(value >> (8 * uint(i)))
	}
	return append(payload, raw...)
}

// appendMicros trägt einen pgoutput-Zeitstempel im
// Mikrosekunden-Format.
func appendMicros(payload []byte, micros int64) []byte {
	raw := make([]byte, 8)
	for i := range raw {
		raw[7-i] = byte(micros >> (8 * uint(i)))
	}
	return append(payload, raw...)
}

// decodeOne dekodiert einen Payload und bricht den Test bei Fehler ab.
func decodeOne(t *testing.T, decoder *decode.Decoder, payload []byte) decode.Event {
	t.Helper()
	event, err := decoder.Decode(payload)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	return event
}

func decodeErr(t *testing.T, payload []byte) error {
	t.Helper()
	decoder := decode.NewDecoder()
	_, err := decoder.Decode(payload)
	if err == nil {
		t.Fatalf("Decode endete ohne Fehler")
	}
	return err
}

func pointerValue(value *string) string {
	if value == nil {
		return "<nil>"
	}
	return *value
}

// TestDecodeBeginCommit trägt BEGIN und COMMIT mit ihren LSN-Ständen.
func TestDecodeBeginCommit(t *testing.T) {
	decoder := decode.NewDecoder()
	begin := decodeOne(t, decoder, beginPayload(0x1A2B3C, 42))
	event, isBegin := begin.(decode.Begin)
	if !isBegin {
		t.Fatalf("BEGIN endete als %T", begin)
	}
	if event.XID != 42 {
		t.Fatalf("XID: %d != 42", event.XID)
	}
	commit := decodeOne(t, decoder, commitPayload(0xABCDEF01, 0xABCDFFFF))
	commitEvent, isCommit := commit.(decode.Commit)
	if !isCommit {
		t.Fatalf("COMMIT endete als %T", commit)
	}
	if commitEvent.CommitLSN != 0xABCDEF01 || commitEvent.EndLSN != 0xABCDFFFF {
		t.Fatalf("Commit-LSN: %x/%x", commitEvent.CommitLSN, commitEvent.EndLSN)
	}
}

// TestDecodeRelationInsert trägt Relation und Insert mit
// Katalog-Auflösung: die Relation-Nachricht füllt den Katalog, der
// Insert löst seine Spalten darüber auf.
func TestDecodeRelationInsert(t *testing.T) {
	decoder := decode.NewDecoder()
	relation := decodeOne(t, decoder, relationPayload(30001, "public", "feed",
		testColumn{name: "id", flags: 1}, testColumn{name: "name"}))
	relationEvent, isRelation := relation.(*decode.Relation)
	if !isRelation {
		t.Fatalf("Relation endete als %T", relation)
	}
	if relationEvent.QualifiedName() != "public.feed" {
		t.Fatalf("QualifiedName: %q", relationEvent.QualifiedName())
	}
	if len(relationEvent.Columns) != 2 || relationEvent.Columns[0].Name != "id" || !relationEvent.Columns[0].Key || relationEvent.Columns[1].Key {
		t.Fatalf("Spalten: %+v", relationEvent.Columns)
	}

	insert := decodeOne(t, decoder, insertPayload(30001, textTuple(text("7"), text("Wert"))))
	change, isChange := insert.(decode.Change)
	if !isChange {
		t.Fatalf("Insert endete als %T", insert)
	}
	if change.Operation != decode.OpInsert || change.Relation != relationEvent {
		t.Fatalf("Change: %+v", change)
	}
	if pointerValue(change.New[0]) != "7" || pointerValue(change.New[1]) != "Wert" {
		t.Fatalf("Tupel-Werte: %v", change.New)
	}
}

// TestDecodeUpdateDelete trägt Update mit Schlüssel- und vollem
// Alt-Tupel sowie Delete: der `K`-Tupel trägt nur Schlüssel-Spalten, der
// `O`-Tupel alle (`LH-FA-CAP-008`).
func TestDecodeUpdateDelete(t *testing.T) {
	decoder := decode.NewDecoder()
	if _, err := decoder.Decode(relationPayload(30002, "public", "feed",
		testColumn{name: "id", flags: 1}, testColumn{name: "name"})); err != nil {
		t.Fatalf("Relation: %v", err)
	}

	update := decodeOne(t, decoder, updatePayloadKey(30002, textTuple(text("7")), textTuple(text("7"), text("neu"))))
	updateEvent, isUpdate := update.(decode.Change)
	if !isUpdate {
		t.Fatalf("Update endete als %T", update)
	}
	if updateEvent.Operation != decode.OpUpdate {
		t.Fatalf("Operation: %v", updateEvent.Operation)
	}
	if len(updateEvent.Old) != 2 || pointerValue(updateEvent.Old[0]) != "7" || pointerValue(updateEvent.Old[1]) != "<nil>" {
		t.Fatalf("Update-Alt-Tupel: %+v", updateEvent.Old)
	}
	if pointerValue(updateEvent.New[1]) != "neu" {
		t.Fatalf("Update-Neu-Tupel: %+v", updateEvent.New)
	}

	updateFull := decodeOne(t, decoder, updatePayloadFull(30002, textTuple(text("7"), text("alt")), textTuple(text("7"), text("neu"))))
	updateFullEvent := updateFull.(decode.Change)
	if pointerValue(updateFullEvent.Old[1]) != "alt" {
		t.Fatalf("Volles Alt-Tupel: %+v", updateFullEvent.Old)
	}

	deleteEvent := decodeOne(t, decoder, deletePayloadKey(30002, textTuple(text("7"))))
	deleteChange, isDelete := deleteEvent.(decode.Change)
	if !isDelete {
		t.Fatalf("Delete endete als %T", deleteEvent)
	}
	if deleteChange.Operation != decode.OpDelete || deleteChange.New != nil || pointerValue(deleteChange.Old[0]) != "7" {
		t.Fatalf("Delete-Change: %+v", deleteChange)
	}

	deleteFull := decodeOne(t, decoder, deletePayloadFull(30002, textTuple(text("7"), text("alt"))))
	deleteFullChange := deleteFull.(decode.Change)
	if pointerValue(deleteFullChange.Old[1]) != "alt" {
		t.Fatalf("Delete-Alt-Tupel: %+v", deleteFullChange.Old)
	}
}

// TestDecodeToast trägt unverändertes TOAST als Abwesenheit
// (`LH-FA-CAP-008` Boundary).
func TestDecodeToast(t *testing.T) {
	decoder := decode.NewDecoder()
	if _, err := decoder.Decode(relationPayload(30003, "public", "toast",
		testColumn{name: "id", flags: 1}, testColumn{name: "blob"})); err != nil {
		t.Fatalf("Relation: %v", err)
	}
	insert := decodeOne(t, decoder, insertPayload(30003, textTuple(text("1"), unchangedToast())))
	change := insert.(decode.Change)
	if change.New[1] != nil {
		t.Fatalf("TOAST-Wert ist Abwesenheit, trägt %v", change.New[1])
	}
}

// TestDecodeTruncate trägt TRUNCATE über mehrere Relationen.
func TestDecodeTruncate(t *testing.T) {
	decoder := decode.NewDecoder()
	if _, err := decoder.Decode(relationPayload(30004, "public", "a")); err != nil {
		t.Fatalf("Relation a: %v", err)
	}
	if _, err := decoder.Decode(relationPayload(30005, "public", "b")); err != nil {
		t.Fatalf("Relation b: %v", err)
	}
	event := decodeOne(t, decoder, truncatePayload(30004, 30005))
	truncate, isTruncate := event.(decode.Truncate)
	if !isTruncate {
		t.Fatalf("TRUNCATE endete als %T", event)
	}
	if len(truncate.Relations) != 2 || truncate.Relations[0].Name != "a" || truncate.Relations[1].Name != "b" {
		t.Fatalf("TRUNCATE-Relationen: %+v", truncate.Relations)
	}
}

// TestDecodeNotRelevant trägt die für den Capture-Pfad irrelevanten
// Nachrichtentypen: sie liefern kein Ereignis.
func TestDecodeNotRelevant(t *testing.T) {
	decoder := decode.NewDecoder()
	// Typ-Nachricht (`Y`): OID, C-String, C-String.
	typePayload := appendString(appendString(appendUint32([]byte{'Y'}, 25), "text"), "text")
	event, err := decoder.Decode(typePayload)
	if err != nil {
		t.Fatalf("Typ-Nachricht: %v", err)
	}
	if event != nil {
		t.Fatalf("Typ-Nachricht endete als %T", event)
	}
}

// TestDecodeUnknownType trägt einen unbekannten Nachrichtentyp als
// sichtbaren Dekodierfehler der Klasse `schema` (`SPEC-008`).
func TestDecodeUnknownType(t *testing.T) {
	err := decodeErr(t, []byte{'Z', 0, 0, 0, 0})
	if !stderrors.Is(err, decode.ErrSchema) {
		t.Fatalf("Unbekannter Typ: %v", err)
	}
}

// TestDecodeChangeBeforeRelation trägt eine Änderung ohne
// Relation-Nachricht als sichtbaren Dekodierfehler (`LH-FA-SCH-004.a`).
func TestDecodeChangeBeforeRelation(t *testing.T) {
	err := decodeErr(t, insertPayload(99999, textTuple(text("x"))))
	if !stderrors.Is(err, decode.ErrSchema) {
		t.Fatalf("Änderung ohne Relation: %v", err)
	}
}

// TestDecodeFlowToCapture trägt den vollständigen Byte-Durchlauf:
// BEGIN, Relation, Insert, COMMIT — das COMMIT meldet die committed
// Quelltransaktion als CaptureCommand (`LH-FA-CAP-006.a`,
// `LH-QA-REL-001.a` Schritte Receive und Decode).
func TestDecodeFlowToCapture(t *testing.T) {
	assembler, err := mapper.NewAssembler("src-1", map[string]mapper.TableBinding{
		"public.feed": {TableID: "tbl-1", SchemaVersion: "sv-1"},
	})
	if err != nil {
		t.Fatalf("NewAssembler: %v", err)
	}
	decoder := decode.NewDecoder()

	steps := [][]byte{
		beginPayload(0x1000, 99),
		relationPayload(40001, "public", "feed", testColumn{name: "id", flags: 1}),
		insertPayload(40001, textTuple(text("5"))),
		commitPayload(0x2000, 0x2004),
	}
	var command *inbound.CaptureCommand
	for i, payload := range steps {
		event, err := decoder.Decode(payload)
		if err != nil {
			t.Fatalf("Schritt %d: %v", i, err)
		}
		consumed, err := assembler.Consume(event)
		if err != nil {
			t.Fatalf("Schritt %d: %v", i, err)
		}
		if consumed != nil {
			command = consumed
		}
	}
	if command == nil {
		t.Fatalf("Durchlauf meldet kein CaptureCommand")
	}
	position, committed := command.Transaction.CommitPosition()
	if !committed || position.Offset != 0x2000 {
		t.Fatalf("Commit-Position: %+v committed=%v", position, committed)
	}
	changes, err := command.Transaction.Changes()
	if err != nil {
		t.Fatalf("Changes: %v", err)
	}
	if len(changes) != 1 || string(changes[0].NewImage) != `{"id":"5"}` {
		t.Fatalf("Changes: %+v", changes)
	}
}
