package main

import (
	"reflect"
	"testing"
)

// knownBlockedReport spiegelt real gemessene Report-Felder (--plan-only
// gegen ein bereits vollständig migriertes Ziel): genau die neun
// bekannten Fremdobjekt-Blocker, sonst nichts. Die Signatur von
// `set_transformation` trägt die Schreibweise des Reports (`in:json` für den
// `json`-Parameter).
func knownBlockedReport() report {
	return report{
		Status: "blocked",
		Blockers: []blocker{
			{Reason: destructiveConfirmationReason, OperationIDs: []string{
				"DropFunction:FUNCTION:a1:a2",
				"DropFunction:FUNCTION:b1:b2",
				"DropFunction:FUNCTION:c1:c2",
				"DropFunction:FUNCTION:d1:d2",
				"DropFunction:FUNCTION:d3:d4",
				"DropFunction:FUNCTION:d5:d6",
				"DropFunction:FUNCTION:d7:d8",
				"DropView:VIEW:e1:e2",
				"DropView:VIEW:f1:f2",
			}},
		},
		Operations: []operation{
			{ID: "DropFunction:FUNCTION:a1:a2", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"enable_table(in:text,in:text,in:text)"}},
			{ID: "DropFunction:FUNCTION:b1:b2", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"disable_table(in:text,in:text,in:text)"}},
			{ID: "DropFunction:FUNCTION:c1:c2", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"exclude_column(in:text,in:text,in:text,in:text)"}},
			{ID: "DropFunction:FUNCTION:d1:d2", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"include_column(in:text,in:text,in:text,in:text)"}},
			{ID: "DropFunction:FUNCTION:d3:d4", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"backfill_table(in:text,in:text,in:text)"}},
			{ID: "DropFunction:FUNCTION:d5:d6", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"set_transformation(in:text,in:text,in:text,in:text,in:json)"}},
			{ID: "DropFunction:FUNCTION:d7:d8", Kind: "DropFunction", ObjectType: "FUNCTION", Path: []string{"remove_transformation(in:text,in:text,in:text,in:text)"}},
			{ID: "DropView:VIEW:e1:e2", Kind: "DropView", ObjectType: "VIEW", Path: []string{"heartbeat"}},
			{ID: "DropView:VIEW:f1:f2", Kind: "DropView", ObjectType: "VIEW", Path: []string{"metrics"}},
		},
	}
}

// withViewSignature hängt an einen Report die Klasse „View-Signatur" für
// die View `name` an: eine ReplaceView-Operation, einen Blocker
// MANUAL_ACTION_REQUIRED und die Diagnose VIEW_SIGNATURE_INCOMPATIBLE zur
// selben Operation — die Form des real gemessenen Precheck-Reports gegen
// ein Ziel mit abweichender View-Signatur.
func withViewSignature(r report, name string) report {
	id := "ReplaceView:VIEW:" + name + ":h1"
	r.Operations = append(r.Operations, operation{ID: id, Kind: "ReplaceView", ObjectType: "VIEW", Path: []string{name}})
	r.Blockers = append(r.Blockers, blocker{Reason: manualActionReason, OperationIDs: []string{id}})
	r.Diagnostics = append(r.Diagnostics, diagnostic{Code: viewSignatureCode, OperationID: id})
	return r
}

func viewSignatureOnlyReport(name string) report {
	return withViewSignature(report{Status: "blocked"}, name)
}

// TestDecideAllowsDestructiveWhenAllBlockersKnown prüft den Regelfall: ein
// zweiter Lauf gegen ein bereits vollständig migriertes Ziel trägt
// ausschließlich die neun bekannten Fremdobjekt-Blocker — decide erlaubt
// --allow-destructive und verlangt keinen Vorlauf.
func TestDecideAllowsDestructiveWhenAllBlockersKnown(t *testing.T) {
	d := decide(knownBlockedReport())
	if !d.allowDestructive {
		t.Fatalf("allowDestructive = false, reason %q — wollte true", d.reason)
	}
	if len(d.dropViews) != 0 {
		t.Fatalf("dropViews = %v — wollte keinen Vorlauf ohne View-Signatur-Blocker", d.dropViews)
	}
}

// TestDecideRefusesUnknownDestructiveBlocker prüft den Negativfall: eine
// künstlich per ALTER TABLE … ADD COLUMN hinzugefügte, nicht deklarierte
// Spalte erzeugt real einen zusätzlichen, unbekannten DropColumn-Blocker
// neben den neun bekannten — decide bleibt leer, auch im Mischfall.
func TestDecideRefusesUnknownDestructiveBlocker(t *testing.T) {
	r := knownBlockedReport()
	r.Blockers[0].OperationIDs = append(r.Blockers[0].OperationIDs, "DropColumn:COLUMN:g1:g2")
	r.Operations = append(r.Operations, operation{ID: "DropColumn:COLUMN:g1:g2", Kind: "DropColumn", ObjectType: "COLUMN", Path: []string{"source", "_scratch_test_col"}})

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer, weil ein Blocker nicht auf der bekannten Liste steht", d)
	}
}

// TestDecideRefusesUnknownBlockerReason prüft, dass decide nur die zwei
// bekannten Blocker-Klassen automatisch auflöst — jede andere gemeldete
// Blocker-Ursache bleibt unaufgelöst.
func TestDecideRefusesUnknownBlockerReason(t *testing.T) {
	// ein zusätzlicher Blocker unbekannter Klasse neben den bekannten
	// Fremdobjekten: alles oder nichts, auch --allow-destructive entfällt
	r := knownBlockedReport()
	r.Blockers = append(r.Blockers, blocker{Reason: "SOME_OTHER_REASON", OperationIDs: []string{"DropFunction:FUNCTION:a1:a2"}})

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer bei unbekannter Blocker-Klasse", d)
	}
}

// TestDecideRefusesMissingOperation prüft die Verteidigungsgrenze: eine
// Blocker-operationId ohne zugehörigen Eintrag in operations[] gilt als
// unbekannt, nicht als übersehbar.
func TestDecideRefusesMissingOperation(t *testing.T) {
	r := knownBlockedReport()
	r.Blockers[0].OperationIDs = append(r.Blockers[0].OperationIDs, "DropTable:TABLE:zz:zz")

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer bei einer operationId ohne Eintrag in operations[]", d)
	}
}

// TestDecideRefusesEmptyBlockers prüft die Randbedingung: ein Report ohne
// jeden Blocker trägt weder einen Grund für --allow-destructive noch für
// einen Vorlauf (ADR-0114 Entscheidung 4/5: kein Vorlauf ohne Anlass —
// additive Änderungen erzeugen keinen Blocker).
func TestDecideRefusesEmptyBlockers(t *testing.T) {
	d := decide(report{Status: "ok"})
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer bei leerem Blockers[]", d)
	}
}

// TestDecideViewSignatureWithKnownForeignObjects prüft den real
// gemessenen Fall des Upgrades über die View-Signaturänderung: die Klasse
// „View-Signatur" neben den neun bekannten Fremdobjekten — decide erlaubt
// --allow-destructive UND meldet die View für den Vorlauf.
func TestDecideViewSignatureWithKnownForeignObjects(t *testing.T) {
	d := decide(withViewSignature(knownBlockedReport(), "changes"))
	if !d.allowDestructive {
		t.Fatalf("allowDestructive = false, reason %q — wollte true (bekannte Fremdobjekte im selben Report)", d.reason)
	}
	if !reflect.DeepEqual(d.dropViews, []string{"changes"}) {
		t.Fatalf("dropViews = %v — wollte [changes]", d.dropViews)
	}
}

// TestDecideViewSignatureAlone prüft, dass die Klasse ohne Fremdobjekte nur
// den Vorlauf auslöst, nicht --allow-destructive.
func TestDecideViewSignatureAlone(t *testing.T) {
	d := decide(viewSignatureOnlyReport("changes"))
	if d.allowDestructive {
		t.Fatal("allowDestructive = true — wollte false, kein destruktiver Blocker im Report")
	}
	if !reflect.DeepEqual(d.dropViews, []string{"changes"}) {
		t.Fatalf("dropViews = %v — wollte [changes]", d.dropViews)
	}
}

// TestDecideViewSignatureWithUnknownBlockerRefusesEverything prüft die
// Alles-oder-nichts-Regel (ADR-0114 Entscheidung 2): die Klasse neben
// einem unbekannten Blocker lässt auch den Vorlauf ausfallen.
func TestDecideViewSignatureWithUnknownBlockerRefusesEverything(t *testing.T) {
	r := withViewSignature(knownBlockedReport(), "changes")
	r.Blockers[0].OperationIDs = append(r.Blockers[0].OperationIDs, "DropColumn:COLUMN:g1:g2")
	r.Operations = append(r.Operations, operation{ID: "DropColumn:COLUMN:g1:g2", Kind: "DropColumn", ObjectType: "COLUMN", Path: []string{"source", "_scratch_test_col"}})

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer (alles oder nichts)", d)
	}
}

// TestDecideManualActionWithoutSignatureDiagnostic prüft, dass
// MANUAL_ACTION_REQUIRED allein keine View freigibt: ohne die Diagnose
// VIEW_SIGNATURE_INCOMPATIBLE zur Operation bleibt es ein unbekannter
// Blocker.
func TestDecideManualActionWithoutSignatureDiagnostic(t *testing.T) {
	r := viewSignatureOnlyReport("changes")
	r.Diagnostics = nil

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer ohne Diagnose", d)
	}
}

// TestDecideRefusesOtherDiagnosticCode prüft, dass die Diagnose den
// Code VIEW_SIGNATURE_INCOMPATIBLE tragen muss — ein anderer Code zur
// selben Operation gibt die View nicht frei.
func TestDecideRefusesOtherDiagnosticCode(t *testing.T) {
	r := viewSignatureOnlyReport("changes")
	r.Diagnostics[0].Code = "VIEW_SIGNATURE_UNKNOWN"

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer bei anderem Diagnosecode", d)
	}
}

// TestDecideRefusesDiagnosticOfOtherOperation prüft die Bindung der
// Diagnose an ihre Operation: eine VIEW_SIGNATURE_INCOMPATIBLE-Diagnose zu
// einer anderen operationId gibt die blockierende Operation nicht frei.
func TestDecideRefusesDiagnosticOfOtherOperation(t *testing.T) {
	r := viewSignatureOnlyReport("changes")
	r.Diagnostics[0].OperationID = "ReplaceView:VIEW:andere:h9"

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer, weil die Diagnose an einer anderen Operation hängt", d)
	}
}

// TestDecideRefusesSignatureDiagnosticOnNonViewOperation prüft, dass die
// Klasse an der Operation ReplaceView/VIEW hängt: dieselbe Diagnose an
// einer anderen Operationsart gibt nichts frei.
func TestDecideRefusesSignatureDiagnosticOnNonViewOperation(t *testing.T) {
	r := viewSignatureOnlyReport("changes")
	r.Operations[0].Kind = "AlterColumnType"
	r.Operations[0].ObjectType = "COLUMN"

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer bei einer Operation, die kein ReplaceView auf einer View ist", d)
	}
}

// TestDecideRefusesManualActionWithoutOperations prüft, dass ein
// MANUAL_ACTION_REQUIRED-Blocker ohne operationIds nicht vakuum als
// „View-Signatur" durchgeht.
func TestDecideRefusesManualActionWithoutOperations(t *testing.T) {
	// neben einem gültigen Blocker: alles oder nichts
	r := knownBlockedReport()
	r.Blockers = append(r.Blockers, blocker{Reason: manualActionReason})

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer bei einem Blocker ohne Operation", d)
	}
}

// TestDecideRefusesDestructiveBlockerWithoutOperations prüft, dass ein
// Blocker DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION ohne operationIds
// nicht vakuum als „bekanntes Fremdobjekt" durchgeht — allein, neben den
// bekannten Fremdobjekten und neben der Klasse „View-Signatur" (alles oder
// nichts: weder --allow-destructive noch Vorlauf).
func TestDecideRefusesDestructiveBlockerWithoutOperations(t *testing.T) {
	empty := blocker{Reason: destructiveConfirmationReason}
	alone := report{Status: "blocked", Blockers: []blocker{empty}}
	withForeign := knownBlockedReport()
	withForeign.Blockers = append(withForeign.Blockers, empty)
	withSignature := viewSignatureOnlyReport("changes")
	withSignature.Blockers = append(withSignature.Blockers, empty)

	for _, tc := range []struct {
		name string
		r    report
	}{
		{"allein", alone},
		{"neben Fremdobjekten", withForeign},
		{"neben View-Signatur", withSignature},
	} {
		d := decide(tc.r)
		if d.allowDestructive || len(d.dropViews) != 0 {
			t.Errorf("%s: decision = %+v — wollte leer bei einem destruktiven Blocker ohne Operation", tc.name, d)
		}
	}
}

// TestDecideRefusesViewNameThatIsNotAnIdentifier prüft die
// Eingabe-Validierung des Namens, der in `DROP VIEW cdc.<name>` gelangt:
// Anführungszeichen, Semikolon, Punkt und Großbuchstaben lehnt decide ab.
func TestDecideRefusesViewNameThatIsNotAnIdentifier(t *testing.T) {
	for _, name := range []string{"changes; DROP TABLE cdc.change", "a.b", "Changes", `"changes"`, "", "1changes", "changes cascade"} {
		d := decide(viewSignatureOnlyReport(name))
		if d.allowDestructive || len(d.dropViews) != 0 {
			t.Errorf("name %q: decision = %+v — wollte leer", name, d)
		}
	}
}

// TestDecideRefusesMultiSegmentViewPath prüft, dass der Pfad genau ein
// Segment trägt (alle Objekte liegen im einen Ziel-Schema `cdc`).
func TestDecideRefusesMultiSegmentViewPath(t *testing.T) {
	r := viewSignatureOnlyReport("changes")
	r.Operations[0].Path = []string{"public", "changes"}

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer bei einem Pfad mit zwei Segmenten", d)
	}
}

// TestDecideListsSeveralViewsSortedWithoutDuplicates prüft die
// Ausgabeform der Liste: sortiert, je View einmal.
func TestDecideListsSeveralViewsSortedWithoutDuplicates(t *testing.T) {
	r := withViewSignature(withViewSignature(viewSignatureOnlyReport("consumer_status"), "changes"), "active_tables")
	// dieselbe View unter einem zweiten Blocker: ein Eintrag in der Liste
	r.Blockers = append(r.Blockers, r.Blockers[0])

	d := decide(r)
	want := []string{"active_tables", "changes", "consumer_status"}
	if !reflect.DeepEqual(d.dropViews, want) {
		t.Fatalf("dropViews = %v — wollte %v", d.dropViews, want)
	}
}

// realViewSignatureReport ist ein gekürzter, real gemessener
// Precheck-Report (--plan-only gegen ein Ziel mit dem Schema vor der
// Spalte `origin`): der Blocker trägt `diagnosticCodes` leer, der Code
// steht in `diagnostics[]` mit `operationId`.
const realViewSignatureReport = `{
  "status": "blocked",
  "exitCode": 8,
  "blockers": [{"reason":"MANUAL_ACTION_REQUIRED","operationIds":["ReplaceView:VIEW:496cb540710d:de089b34b29f"],"diagnosticCodes":[]},{"reason":"DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION","operationIds":["DropView:VIEW:e652b961e18f:0b921c47c4f7"],"diagnosticCodes":[]}],
  "diagnostics": [{"code":"VIEW_SIGNATURE_INCOMPATIBLE","severity":"BLOCKER","message":"x","operationId":"ReplaceView:VIEW:496cb540710d:de089b34b29f"},{"code":"UNSAFE_DEPENDENCY_PAIR","severity":"WARNING","message":"y"}],
  "operations": [{"id":"ReplaceView:VIEW:496cb540710d:de089b34b29f","kind":"ReplaceView","objectType":"VIEW","path":["changes"],"phase":"VIEWS","rendered":false,"skipped":true},{"id":"DropView:VIEW:e652b961e18f:0b921c47c4f7","kind":"DropView","objectType":"VIEW","path":["heartbeat"],"phase":"VIEWS","rendered":true,"skipped":false}]
}`

// TestParseReportDecodesRealViewSignatureReport bindet die JSON-Feldnamen
// (`diagnostics[].code`, `diagnostics[].operationId`) an die reale
// Report-Form: ein umbenanntes Feld ließe die Klasse „View-Signatur"
// unerkannt.
func TestParseReportDecodesRealViewSignatureReport(t *testing.T) {
	r, err := parseReport([]byte(realViewSignatureReport))
	if err != nil {
		t.Fatalf("parseReport: %v", err)
	}
	d := decide(r)
	if !reflect.DeepEqual(d.dropViews, []string{"changes"}) {
		t.Fatalf("dropViews = %v, reason %q — wollte [changes]", d.dropViews, d.reason)
	}
	if !d.allowDestructive {
		t.Fatal("allowDestructive = false — wollte true: der destruktive Blocker des Reports (DropView heartbeat) steht auf der bekannten Liste")
	}
}

// TestDecideRefusesManualActionOperationMissingFromReport prüft die
// Verteidigungsgrenze der Klasse „View-Signatur": eine operationId eines
// MANUAL_ACTION_REQUIRED-Blockers ohne Eintrag in operations[] gibt nichts
// frei, auch neben bekannten Fremdobjekten.
func TestDecideRefusesManualActionOperationMissingFromReport(t *testing.T) {
	r := knownBlockedReport()
	r.Blockers = append(r.Blockers, blocker{Reason: manualActionReason, OperationIDs: []string{"ReplaceView:VIEW:zz:zz"}})

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer bei einer operationId ohne Eintrag in operations[]", d)
	}
}

// TestDecideRefusesOtherSignatureSpelling prüft die Bindung der Bekannt-Liste
// an die Schreibweise des Reports: die Funktion `set_transformation` mit dem
// Parametertyp `in:jsonb` (statt der gemessenen Schreibweise `in:json`) ist
// kein bekanntes Fremdobjekt — decide bleibt leer.
func TestDecideRefusesOtherSignatureSpelling(t *testing.T) {
	r := knownBlockedReport()
	for i, op := range r.Operations {
		if op.ID == "DropFunction:FUNCTION:d5:d6" {
			r.Operations[i].Path = []string{"set_transformation(in:text,in:text,in:text,in:text,in:jsonb)"}
		}
	}

	d := decide(r)
	if d.allowDestructive || len(d.dropViews) != 0 {
		t.Fatalf("decision = %+v — wollte leer, weil die Signatur nicht die gemessene Schreibweise trägt", d)
	}
}
