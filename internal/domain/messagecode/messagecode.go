// Package messagecode trägt die Meldungscodes der Betreiber-Sicht: eine
// Kennung `PCF-<S><NNNN>` je Fehler- und Ablehnungsursache, feiner als die
// Fehlerklasse. Die Schwere `S` ist `E` (Fehler), `W` (Warnung) oder `I`
// (Information, reserviert); bei Fehlern ist die erste Ziffer die
// Fehlerklasse, `8` die Ablehnung einer Aufrufer-Eingabe; bei Warnungen
// (`W`) ist sie der Bereich (`Area`). Die Tabelle (`codes.go`) ist die Quelle
// der Wahrheit.
package messagecode

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Code ist ein Meldungscode in der Form `PCF-<S><NNNN>`.
type Code string

// Class ist die Fehlerklasse eines Codes; `ClassNone` trägt der Bereich der
// Ablehnungen, der keine Fehlerklasse hat. Die Namen sind die der
// Fehlerklassen (`model.ErrorClass`), das Paket importiert sie nicht.
type Class string

const (
	ClassNone          Class = ""
	ClassTransient     Class = "transient"
	ClassConfiguration Class = "configuration"
	ClassPermission    Class = "permission"
	ClassSchema        Class = "schema"
	ClassStorage       Class = "storage"
	ClassReplication   Class = "replication"
	ClassInternal      Class = "internal"
)

// Status ist der Lebensstand eines Codes: ein zurückgezogener Code bleibt in
// Tabelle und Katalog und wird nie neu belegt.
type Status string

const (
	StatusActive    Status = "active"
	StatusWithdrawn Status = "withdrawn"
)

// Entry ist eine Zeile der Tabelle.
type Entry struct {
	Code   Code
	Class  Class
	Status Status
}

// Pattern ist die Form eines Codes (ERE).
const Pattern = `PCF-[EWI][0-9]{4}`

var codeForm = regexp.MustCompile(`^` + Pattern + `$`)

// Valid meldet, ob der Code die Form `PCF-<S><NNNN>` hat.
func Valid(c Code) bool { return codeForm.MatchString(string(c)) }

// DigitClass liest die Klasse aus der ersten Ziffer eines Fehlercodes (`E`):
// 1 bis 7 sind die sieben Fehlerklassen, 8 die Ablehnung (`ClassNone`).
// Ein anderer Code liefert `false`.
func DigitClass(c Code) (Class, bool) {
	if !Valid(c) || c[4] != 'E' {
		return ClassNone, false
	}
	switch c[5] {
	case '1':
		return ClassTransient, true
	case '2':
		return ClassConfiguration, true
	case '3':
		return ClassPermission, true
	case '4':
		return ClassSchema, true
	case '5':
		return ClassStorage, true
	case '6':
		return ClassReplication, true
	case '7':
		return ClassInternal, true
	case '8':
		return ClassNone, true
	default:
		return ClassNone, false
	}
}

// LogKey ist der Name des Log-Attributs, das den Code einer Warnung trägt.
const LogKey = "code"

// Area liest den Bereich aus der ersten Ziffer eines Warncodes (`W`):
// 1 Erfassung und Replikation, 2 Backfill, 3 Retention und Speicher,
// 4 Verwaltung, 5 Konfiguration und Start. Ein anderer Code und der
// reservierte Bereich 9 liefern `false`.
func Area(c Code) (int, bool) {
	if !Valid(c) || c[4] != 'W' || c[5] < '1' || c[5] > '5' {
		return 0, false
	}
	return int(c[5] - '0'), true
}

// Lookup liefert die Tabellenzeile eines Codes.
func Lookup(c Code) (Entry, bool) {
	for _, e := range Table {
		if e.Code == c {
			return e, true
		}
	}
	return Entry{}, false
}

// ClassOf liefert die Fehlerklasse eines Codes der Tabelle; ein unbekannter
// Code und eine Ablehnung liefern `ClassNone`.
func ClassOf(c Code) Class {
	e, _ := Lookup(c)
	return e.Class
}

// Fallback liefert den Rückfall einer Fehlerklasse (`…000`): jeder
// klassifizierte Fehler ohne Einzelursache trägt ihn, nie keinen Code.
func Fallback(class Class) Code {
	switch class {
	case ClassTransient:
		return TransientFallback
	case ClassConfiguration:
		return ConfigurationFallback
	case ClassPermission:
		return PermissionFallback
	case ClassSchema:
		return SchemaFallback
	case ClassStorage:
		return StorageFallback
	case ClassReplication:
		return ReplicationFallback
	default:
		return InternalFallback
	}
}

// Head ist der Kopf eines Fehlertexts: `Fehlerklasse <klasse> [<code>]: `.
func Head(c Code) string {
	return fmt.Sprintf("Fehlerklasse %s [%s]: ", ClassOf(c), c)
}

// Error ist ein klassifizierter Fehlerwert: sein Text beginnt mit dem Kopf
// seines Codes. Er wird als Sentinel vereinbart und über `errors.Is`
// verglichen.
type Error struct {
	code   Code
	reason string
}

// New legt den Fehlerwert zum Code einer Fehlerklasse an. Ein Code außerhalb
// der Tabelle oder der Bereich der Ablehnungen ist ein Programmierfehler und
// endet mit einer Panik, beim Sentinel schon im Paket-Init.
func New(code Code, reason string) error {
	e, ok := Lookup(code)
	if !ok || e.Class == ClassNone {
		panic("messagecode: " + string(code) + " ist kein Code einer Fehlerklasse")
	}
	return &Error{code: code, reason: reason}
}

// Error liefert Kopf und Ursache.
func (e *Error) Error() string { return Head(e.code) + e.reason }

// Code liefert den Meldungscode des Fehlerwerts.
func (e *Error) Code() Code { return e.code }

// From liefert den Code des ersten klassifizierten Fehlerwerts der Kette.
func From(err error) (Code, bool) {
	var target *Error
	if errors.As(err, &target) {
		return target.code, true
	}
	return "", false
}

// Codes liefert die Codes aller klassifizierten Fehlerwerte der Kette in der
// Reihenfolge eines Tiefendurchlaufs; auch ein Fehler mit mehreren
// `Unwrap`-Zielen (`errors.Join`, mehrere `%w`) wird ganz gelesen.
func Codes(err error) []Code {
	var out []Code
	var walk func(error)
	walk = func(e error) {
		if e == nil {
			return
		}
		if c, ok := e.(*Error); ok {
			out = append(out, c.code)
		}
		switch u := e.(type) {
		case interface{ Unwrap() error }:
			walk(u.Unwrap())
		case interface{ Unwrap() []error }:
			for _, w := range u.Unwrap() {
				walk(w)
			}
		}
	}
	walk(err)
	return out
}

// WithoutHead liefert den Text des Fehlers ohne den Kopf seines Codes; trägt
// die Kette keinen klassifizierten Fehlerwert, ist es der Text selbst.
func WithoutHead(err error) string {
	text := err.Error()
	if code, ok := From(err); ok {
		return strings.Replace(text, Head(code), "", 1)
	}
	return text
}

// RunMessage ist der Text einer Fehlerzeile, die den Code als Feld des Textes
// trägt (`error_message` eines Backfill-Runs): `<klasse> [<code>]: <text>`.
func RunMessage(c Code, text string) string {
	return fmt.Sprintf("%s [%s]: %s", ClassOf(c), c, text)
}

// RejectionMessage ist der Text eines abgelehnten Antrags:
// `abgelehnt [<code>]: <text>`.
func RejectionMessage(c Code, text string) string {
	return fmt.Sprintf("abgelehnt [%s]: %s", c, text)
}
