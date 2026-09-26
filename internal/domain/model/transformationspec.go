package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"unicode/utf8"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// ruleNameShape ist das Alphabet des Regelnamens (`SPEC-030`, Bezeichner):
// 1 bis 63 Zeichen aus `a`–`z`, `0`–`9` und `_`.
var ruleNameShape = regexp.MustCompile(`^[a-z0-9_]{1,63}$`)

// CheckRuleName prüft den Regelnamen eines Antrags gegen das Alphabet
// (`SPEC-030`, Bezeichner); ein leerer Name und ein Name mit Großbuchstaben
// liegen außerhalb. Verletzung: `ErrInvalidRuleName`.
func CheckRuleName(name string) error {
	if !ruleNameShape.MatchString(name) {
		return domainerrors.ErrInvalidRuleName
	}
	return nil
}

// TransformationSpecError trägt zu einem Ablehnungsgrund den Wert, den die
// Adresse des Fehlertextes nennt (`SPEC-019`): den `kind`-Wert bei
// `ErrUnknownTransformationKind`, den Schlüsselnamen bei
// `ErrUnknownRuleSpecKey`. `errors.Is` löst auf den Grund auf.
type TransformationSpecError struct {
	Err    error
	Detail string
}

// Error trägt den Grund und den Wert.
func (e *TransformationSpecError) Error() string {
	return e.Err.Error() + ": " + e.Detail
}

// Unwrap gibt den Grund frei.
func (e *TransformationSpecError) Unwrap() error { return e.Err }

// TransformationSpec trägt die geprüfte Regelform eines Antrags
// (`SPEC-030`): der Regeltyp und die Parameter, die Form ist gültig, die
// Konfliktfreiheit gegen den Regelstand einer Tabelle nicht geprüft. Der
// Zielname darf der Quellspalte gleichen: dieser Fall gehört zu K3
// (`CheckConflicts`), nicht zur Form. Die Felder sind nicht exportiert; ein
// Wert entsteht über `ParseTransformationSpec`.
type TransformationSpec struct {
	kind   TransformationKind
	column string
	to     string
}

// Kind trägt den Regeltyp.
func (s TransformationSpec) Kind() TransformationKind { return s.kind }

// Column trägt die Quellspalte der Regel.
func (s TransformationSpec) Column() string { return s.column }

// Target trägt den Zielnamen des Regeltyps; ein Regeltyp ohne Zielname
// trägt den leeren Wert.
func (s TransformationSpec) Target() string { return s.to }

// ParseTransformationSpec liest die Regelform aus dem JSON-Text von
// `rule_spec` (`SPEC-030`) strikt und prüft die Formzeilen zwei bis fünf in
// der Reihenfolge der Fehlertext-Tabelle (`SPEC-019`; die erste Formzeile,
// den Regelnamen, prüft `CheckRuleName`):
//  a) der Text ist gültiges UTF-8 und ein JSON-Objekt, `kind` fehlt nicht und
//     ist eine Zeichenkette (der leere Text, JSON-`null` und ein Wert ohne
//     Objekt gehören hierher): `ErrInvalidRuleSpec`;
//  b) `kind` gehört zu `TransformationKinds`: sonst
//     `ErrUnknownTransformationKind` mit dem `kind`-Wert;
//  c) jeder Schlüssel gehört zum Regeltyp: sonst `ErrUnknownRuleSpecKey` mit
//     dem Schlüssel, bei mehreren dem ersten in aufsteigender Ordnung;
//  d) die Pflichtschlüssel des Regeltyps stehen als Zeichenketten und
//     `column`/`to` tragen die Bezeichner-Form (`NewRenameColumn`): sonst
//     `ErrInvalidRuleSpec`.
//
// Ein doppelter Schlüssel ist keine Eingabe dieser Funktion: der Text kommt
// aus einer `jsonb`-Spalte, die den letzten Wert je Schlüssel hält.
func ParseTransformationSpec(text string) (TransformationSpec, error) {
	if !utf8.ValidString(text) {
		return TransformationSpec{}, fmt.Errorf("%w: kein gültiges UTF-8", domainerrors.ErrInvalidRuleSpec)
	}
	// Der leere Text, ein Wert ohne Objekt und ein Text mit Nachlauf enden
	// als Dekodierfehler; JSON-`null` dekodiert zur leeren Map und endet am
	// fehlenden `kind`.
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(text), &fields); err != nil {
		return TransformationSpec{}, fmt.Errorf("%w: %w", domainerrors.ErrInvalidRuleSpec, err)
	}
	kind, ok := jsonString(fields["kind"])
	if !ok {
		return TransformationSpec{}, fmt.Errorf("%w: kind fehlt oder ist keine Zeichenkette", domainerrors.ErrInvalidRuleSpec)
	}
	allowed, known := allowedRuleKeys(TransformationKind(kind))
	if !known {
		return TransformationSpec{}, &TransformationSpecError{Err: domainerrors.ErrUnknownTransformationKind, Detail: kind}
	}
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if !containsName(allowed, key) {
			return TransformationSpec{}, &TransformationSpecError{Err: domainerrors.ErrUnknownRuleSpecKey, Detail: key}
		}
	}
	switch TransformationKind(kind) {
	case TransformationRenameColumn:
		column, columnOK := jsonString(fields["column"])
		to, toOK := jsonString(fields["to"])
		if !columnOK || !toOK {
			return TransformationSpec{}, fmt.Errorf("%w: column und to sind Zeichenketten", domainerrors.ErrInvalidRuleSpec)
		}
		// Die Bezeichner-Form hat eine Quelle, den Konstruktor der Regel;
		// der Platzhalter-Name ist gültig, ein gleicher Zielname ist hier
		// keine Formverletzung (K3).
		if _, err := NewRenameColumn("r", column, to); err != nil && !errors.Is(err, domainerrors.ErrTransformationTargetIsColumn) {
			return TransformationSpec{}, fmt.Errorf("%w: %w", domainerrors.ErrInvalidRuleSpec, err)
		}
		return TransformationSpec{kind: TransformationRenameColumn, column: column, to: to}, nil
	}
	return TransformationSpec{}, &TransformationSpecError{Err: domainerrors.ErrUnknownTransformationKind, Detail: kind}
}

// allowedRuleKeys liefert die Schlüssel eines Regeltyps und ob der Typ zur
// geschlossenen Menge gehört.
func allowedRuleKeys(kind TransformationKind) ([]string, bool) {
	switch kind {
	case TransformationRenameColumn:
		return []string{"kind", "column", "to"}, true
	}
	return nil, false
}

// jsonString liest einen JSON-Wert, der eine Zeichenkette sein muss;
// abwesend, `null` und jeder andere Typ liefern `false`.
func jsonString(raw json.RawMessage) (string, bool) {
	if len(raw) == 0 || raw[0] != '"' {
		return "", false
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", false
	}
	return value, true
}

// CheckConflicts prüft die Konfliktfreiheit K1 bis K3 (`SPEC-019`) der
// Regel `name` gegen den Regelstand `rules` der Tabelle und die
// Spaltennamen `columns` der Quelltabelle (Katalog, ausgeschlossene Spalten
// eingeschlossen), in der Reihenfolge K1, K2, K3:
//   - K1 `ErrRuleNameTaken`: eine Regel trägt den Namen;
//   - K2 `ErrColumnHasRule`: eine Regel trägt dieselbe Quellspalte;
//   - K3 `ErrTargetCollidesWithRule`: der Zielname gleicht dem Zielnamen
//     einer Regel; danach `ErrTargetCollidesWithColumn`: er gleicht der
//     Quellspalte oder einem Spaltennamen der Quelltabelle.
//
// Die Vergleiche sind zeichengenau. K4 (die Quellspalte existiert an der
// Quelle) prüft der Aufrufer nach dieser Funktion. Die Funktion ist rein.
func (s TransformationSpec) CheckConflicts(name string, rules []Transformation, columns []string) error {
	for _, rule := range rules {
		if rule.Name() == name {
			return domainerrors.ErrRuleNameTaken
		}
	}
	for _, rule := range rules {
		if rule.Column() == s.column {
			return domainerrors.ErrColumnHasRule
		}
	}
	if s.to != "" {
		for _, rule := range rules {
			if rule.To() == s.to {
				return domainerrors.ErrTargetCollidesWithRule
			}
		}
		if s.to == s.column || containsName(columns, s.to) {
			return domainerrors.ErrTargetCollidesWithColumn
		}
	}
	return nil
}

// Build legt die Regel unter dem Namen an. Der Name gehört zuvor durch
// `CheckRuleName`; ein Zielname gleich der Quellspalte endet mit
// `ErrTransformationTargetIsColumn`, den `CheckConflicts` vorher als K3
// ablehnt.
func (s TransformationSpec) Build(name string) (Transformation, error) {
	switch s.kind {
	case TransformationRenameColumn:
		return NewRenameColumn(name, s.column, s.to)
	}
	return Transformation{}, &TransformationSpecError{Err: domainerrors.ErrUnknownTransformationKind, Detail: string(s.kind)}
}

// TransformationRecord trägt eine `applied`-Zeile der beiden
// Transformations-Antragsarten einer Tabelle, wie der Regelstand sie
// auswertet (`ADR-0112` Teilfrage 6): Art, Regelname und Regelform als
// JSON-Text (bei `remove_transformation` leer).
type TransformationRecord struct {
	Kind AdministrationRequestKind
	Name string
	Spec string
}

// FoldTransformations leitet den Regelstand einer Tabelle aus ihren
// `applied`-Zeilen ab, die in der Ordnung `requested_at`, bei gleichem
// Zeitstempel nach der Antrags-Kennung stehen (`SPEC-019`):
// `set_transformation` trägt die Regel ein — eine Regel gleichen Namens wird
// an ihrer Stelle ersetzt, sonst hinten angehängt —, `remove_transformation`
// nimmt sie unter dem Namen heraus; jede andere Art trägt keinen Stand. Eine
// Zeile, deren Regelform nicht mehr zu einer Regel führt, endet als Fehler:
// der Stand wird nie um eine Zeile verkürzt. Die Rückgabe ist eine frische
// Liste (`nil` ohne Regel), die der Aufrufer nicht mit dem Speicher der
// Eingabe teilt. Die Funktion ist rein.
func FoldTransformations(records []TransformationRecord) ([]Transformation, error) {
	var rules []Transformation
	for _, record := range records {
		switch record.Kind {
		case AdministrationRequestSetTransformation:
			spec, err := ParseTransformationSpec(record.Spec)
			if err != nil {
				return nil, fmt.Errorf("Regelstand: Regel %q: %w", record.Name, err)
			}
			rule, err := spec.Build(record.Name)
			if err != nil {
				return nil, fmt.Errorf("Regelstand: Regel %q: %w", record.Name, err)
			}
			rules = replaceOrAppend(rules, rule)
		case AdministrationRequestRemoveTransformation:
			rules = dropByName(rules, record.Name)
		}
	}
	return rules, nil
}

// replaceOrAppend liefert eine neue Liste mit der Regel: eine Regel gleichen
// Namens wird an ihrer Stelle ersetzt, sonst hinten angehängt.
func replaceOrAppend(rules []Transformation, rule Transformation) []Transformation {
	next := make([]Transformation, 0, len(rules)+1)
	replaced := false
	for _, existing := range rules {
		if existing.name == rule.name {
			next = append(next, rule)
			replaced = true
			continue
		}
		next = append(next, existing)
	}
	if !replaced {
		next = append(next, rule)
	}
	return next
}

// dropByName liefert eine neue Liste ohne die Regel unter dem Namen; sie ist
// `nil`, wenn keine Regel bleibt.
func dropByName(rules []Transformation, name string) []Transformation {
	var next []Transformation
	for _, rule := range rules {
		if rule.name != name {
			next = append(next, rule)
		}
	}
	return next
}
