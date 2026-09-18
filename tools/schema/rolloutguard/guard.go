package main

import (
	"fmt"
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

// destructiveConfirmationReason ist der einzige Blocker-Grund, den dieser
// Guard automatisch auflöst — die einzige Klasse, in der d-migrate
// überhaupt Blocker-Operationen (mit operationIds) statt eines
// pauschalen Fehlers meldet (real gemessen, --plan-only-Report).
const destructiveConfirmationReason = "DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION"

// decide meldet, ob der nachfolgende `--execute`-Schritt zusätzlich mit
// `--allow-destructive` laufen darf: nur wenn der Report überhaupt
// blockiert war und JEDE Blocker-Operation (a) den Grund
// destructiveConfirmationReason trägt und (b) auf der Liste
// knownForeignObjects steht. Das Ergebnis entscheidet bewusst NICHT, ob
// `--execute` überhaupt läuft — es läuft immer, damit eine echte,
// gleichzeitig anstehende Schema-Änderung (z. B. eine neue Spalte) nicht
// verlustig geht, nur weil die sechs bekannten Fremdobjekte ebenfalls im
// Plan stehen (real geprüft: ein reines Überspringen von `--execute` ließ
// eine per ALTER TABLE … DROP COLUMN entfernte, von schema.yaml weiterhin
// deklarierte Spalte nicht zurückkommen). Ein einziger unbekannter
// Blocker — real geprüft mit einer künstlich per ALTER TABLE … ADD COLUMN
// hinzugefügten, nicht deklarierten Spalte, die d-migrate als
// unbekannten DropColumn-Blocker neben den sechs bekannten meldet — lässt
// decide false liefern; der `--execute`-Lauf läuft dann ohne
// `--allow-destructive` und bricht mit demselben Blocker real mit Exit 8
// ab.
func decide(r report) (allowDestructive bool, reason string) {
	if len(r.Blockers) == 0 {
		return false, "kein Blocker im Report — kein Grund fuer --allow-destructive"
	}

	opByID := make(map[string]foreignObject, len(r.Operations))
	for _, op := range r.Operations {
		opByID[op.ID] = foreignObject{kind: op.Kind, objectType: op.ObjectType, path: strings.Join(op.Path, ".")}
	}

	for _, b := range r.Blockers {
		if b.Reason != destructiveConfirmationReason {
			return false, fmt.Sprintf("unbekannte Blocker-Klasse %q", b.Reason)
		}
		for _, id := range b.OperationIDs {
			obj, ok := opByID[id]
			if !ok {
				return false, fmt.Sprintf("Blocker-Operation %s nicht im Report gefunden", id)
			}
			if !knownForeignObjects[obj] {
				return false, fmt.Sprintf("unbekannter Blocker %s (%s %s %s)", id, obj.kind, obj.objectType, obj.path)
			}
		}
	}
	return true, "alle Blocker auf der bekannten Fremdobjekt-Liste (ADR-0043) — --allow-destructive ist sicher"
}
