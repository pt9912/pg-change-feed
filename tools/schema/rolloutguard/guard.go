package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// knownForeignObjects trägt die aktuell sechs Objekte, die außerhalb des
// neutralen Modells (tools/schema/schema.yaml) über
// tools/schema/nacharbeit-*.sql angelegt werden und deshalb bei jedem
// zweiten `schema migrate`-Lauf gegen ein bereits migriertes Ziel als
// destruktive Blocker erscheinen (real gemessen, --plan-only-Report gegen
// eine frisch migrierte Ziel-DB). Ein neues nacharbeit-*.sql-Skript trägt
// seinen Eintrag hier im selben Commit wie seine Aufrufzeile im
// Makefile-Target `schema-rollout` — dieselbe Kolokation.
var knownForeignObjects = map[foreignObject]bool{
	{kind: "DropFunction", objectType: "FUNCTION", path: "enable_table(in:text,in:text,in:text)"}:           true,
	{kind: "DropFunction", objectType: "FUNCTION", path: "disable_table(in:text,in:text,in:text)"}:          true,
	{kind: "DropFunction", objectType: "FUNCTION", path: "exclude_column(in:text,in:text,in:text,in:text)"}: true,
	{kind: "DropFunction", objectType: "FUNCTION", path: "include_column(in:text,in:text,in:text,in:text)"}: true,
	{kind: "DropView", objectType: "VIEW", path: "heartbeat"}:                                               true,
	{kind: "DropView", objectType: "VIEW", path: "metrics"}:                                                 true,
}

// destructiveConfirmationReason ist der Blocker-Grund, unter dem d-migrate
// die Abbau-Operationen der bekannten Fremdobjekte meldet (real gemessen,
// --plan-only-Report).
const destructiveConfirmationReason = "DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION"

// manualActionReason und viewSignatureCode bilden zusammen die Klasse
// „View-Signatur" (ADR-0114 Entscheidung 1): eine Operation ReplaceView
// (Objekttyp VIEW) unter dem Blocker-Grund MANUAL_ACTION_REQUIRED, zu der
// der Report die Diagnose VIEW_SIGNATURE_INCOMPATIBLE trägt (real
// gemessen: d-migrate rendert ein CREATE OR REPLACE VIEW nur bei gleicher
// Spaltenzahl, -reihenfolge, -namen und -typen).
const (
	manualActionReason = "MANUAL_ACTION_REQUIRED"
	viewSignatureCode  = "VIEW_SIGNATURE_INCOMPATIBLE"
)

// viewIdentifier ist die einzige Form, in der ein Name aus dem Report in
// `DROP VIEW cdc.<name>` gelangt: ein einfacher, kleingeschriebener
// Bezeichner. Alles andere gilt als unbekannter Blocker.
var viewIdentifier = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

// decision ist das Ergebnis von decide. Leer (kein --allow-destructive,
// keine Views) heißt: der Aufrufer läuft ohne jede Erweiterung weiter und
// bricht bei einem Blocker mit Exit 8 ab.
type decision struct {
	// allowDestructive: der `--execute`-Schritt läuft zusätzlich mit
	// `--allow-destructive` (nur bekannte Fremdobjekte blockieren).
	allowDestructive bool
	// dropViews: die Views der Klasse „View-Signatur", die der Aufrufer vor
	// `--execute` per `DROP VIEW cdc.<name>` (ohne CASCADE) entfernt;
	// d-migrate legt sie anschließend neu an. Sortiert, ohne Duplikate.
	dropViews []string
	reason    string
}

// decide wertet den Precheck-Report aus. Alles oder nichts (ADR-0114
// Entscheidung 2): jeder Blocker muss entweder zur Klasse „View-Signatur"
// gehören oder unter destructiveConfirmationReason eine Operation auf
// knownForeignObjects sein; schon ein einziger anderer Blocker — real
// geprüft mit einer künstlich per ALTER TABLE … ADD COLUMN hinzugefügten,
// nicht deklarierten Spalte, die d-migrate als unbekannten
// DropColumn-Blocker meldet — lässt die Entscheidung leer: kein Vorlauf,
// kein --allow-destructive, der --execute-Lauf bricht mit Exit 8 ab.
//
// Das Ergebnis entscheidet bewusst NICHT, ob `--execute` überhaupt läuft —
// der Aufrufer (Makefile-Target schema-rollout) lässt `--execute` in jedem
// Fall laufen, damit eine echte, gleichzeitig anstehende Schema-Änderung im
// selben Plan wirksam bleibt (Regressionsbeleg:
// tools/harness/run-schema-rollout-guard-test.sh Lauf 3).
func decide(r report) decision {
	if len(r.Blockers) == 0 {
		return decision{reason: "kein Blocker im Report — weder --allow-destructive noch Vorlauf noetig"}
	}

	opByID := make(map[string]operation, len(r.Operations))
	for _, op := range r.Operations {
		opByID[op.ID] = op
	}
	signatureDiagnosed := make(map[string]bool)
	for _, d := range r.Diagnostics {
		if d.Code == viewSignatureCode && d.OperationID != "" {
			signatureDiagnosed[d.OperationID] = true
		}
	}

	var d decision
	views := make(map[string]bool)
	for _, b := range r.Blockers {
		switch b.Reason {
		case destructiveConfirmationReason:
			d.allowDestructive = true
			for _, id := range b.OperationIDs {
				op, ok := opByID[id]
				if !ok {
					return decision{reason: fmt.Sprintf("Blocker-Operation %s nicht im Report gefunden", id)}
				}
				if !knownForeignObjects[foreignObject{kind: op.Kind, objectType: op.ObjectType, path: strings.Join(op.Path, ".")}] {
					return decision{reason: fmt.Sprintf("unbekannter Blocker %s (%s %s %s)", id, op.Kind, op.ObjectType, strings.Join(op.Path, "."))}
				}
			}
		case manualActionReason:
			if len(b.OperationIDs) == 0 {
				return decision{reason: fmt.Sprintf("Blocker %q ohne Operation", b.Reason)}
			}
			for _, id := range b.OperationIDs {
				op, ok := opByID[id]
				if !ok {
					return decision{reason: fmt.Sprintf("Blocker-Operation %s nicht im Report gefunden", id)}
				}
				if op.Kind != "ReplaceView" || op.ObjectType != "VIEW" || !signatureDiagnosed[id] {
					return decision{reason: fmt.Sprintf("unbekannter Blocker %s (%s %s %s) — keine View-Signatur-Aenderung", id, op.Kind, op.ObjectType, strings.Join(op.Path, "."))}
				}
				if len(op.Path) != 1 || !viewIdentifier.MatchString(op.Path[0]) {
					return decision{reason: fmt.Sprintf("View-Name %q der Operation %s ist kein einfacher Bezeichner", strings.Join(op.Path, "."), id)}
				}
				views[op.Path[0]] = true
			}
		default:
			return decision{reason: fmt.Sprintf("unbekannte Blocker-Klasse %q", b.Reason)}
		}
	}

	for v := range views {
		d.dropViews = append(d.dropViews, v)
	}
	sort.Strings(d.dropViews)
	d.reason = "alle Blocker sind bekannte Fremdobjekte (ADR-0043) oder View-Signatur-Aenderungen (ADR-0114)"
	return d
}
