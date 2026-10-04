// Package apiauth ordnet ein Bearer-Token der HTTP- oder gRPC-API einer
// Rechtsklasse zu. Beide Driving-Adapter fragen dieselbe Stelle, deshalb
// gilt ein Token auf beiden Wegen gleich (`LH-FA-SST-012`); die Antwort bei
// fehlender oder unzureichender Klasse bleibt Sache des Adapters.
package apiauth

import (
	"crypto/sha256"
	"crypto/subtle"
)

// Role ist die Rechtsklasse eines Aufrufers: `Admin` deckt implizit `Reader`
// ab (Ordnung `Admin > Reader > None`); `None` trägt das fehlende und das
// unbekannte Token gemeinsam.
type Role int

const (
	None Role = iota
	Reader
	Admin
)

// Classifier hält die gültigen Token je Klasse. Er ist nach `New` unveränderlich
// und von mehreren Goroutinen gleichzeitig lesbar.
type Classifier struct {
	reader [][sha256.Size]byte
	admin  [][sha256.Size]byte
}

// New legt den Klassifikator aus den Token je Klasse an. Die Mengen werden
// als SHA-256-Werte gehalten, damit jeder Vergleich über gleich lange Werte
// läuft. Ein leeres Element in einer Menge bleibt wirkungslos: `Classify`
// weist das leere Aufruf-Token zuerst ab.
func New(reader, admin []string) *Classifier {
	return &Classifier{reader: digests(reader), admin: digests(admin)}
}

// FromConfig legt den Klassifikator aus dem Singular-Token und der Liste je
// Klasse an: gültig ist je Klasse die Vereinigung beider. Beide Adapter bilden
// ihre Konfiguration über diese eine Stelle ab.
func FromConfig(readerSingular string, readers []string, adminSingular string, admins []string) *Classifier {
	return New(
		append([]string{readerSingular}, readers...),
		append([]string{adminSingular}, admins...),
	)
}

func digests(tokens []string) [][sha256.Size]byte {
	out := make([][sha256.Size]byte, len(tokens))
	for i, token := range tokens {
		out[i] = sha256.Sum256([]byte(token))
	}
	return out
}

// Classify ordnet ein Token seiner Klasse zu; steht derselbe Wert in beiden
// Mengen, gewinnt `Admin`. Ein leeres Token trifft nie: ein fehlender Header
// ist ein leeres Token und bleibt damit in jeder Klasse unbekannt, auch gegen
// eine ungesetzte.
//
// Der Vergleich läuft gegen alle Token beider Klassen ohne Abbruch beim ersten
// Treffer (`subtle.ConstantTimeCompare` auf SHA-256-Werten gleicher Länge).
// Grenze: die Laufzeit hängt von der Zahl der konfigurierten Token ab, nicht
// vom Inhalt des Aufruf-Tokens; die Zeitkonstanz ist eine Erwartung an diese
// Umsetzung, kein Test misst sie.
func (c *Classifier) Classify(token string) Role {
	if token == "" {
		return None
	}
	sum := sha256.Sum256([]byte(token))
	isReader := matchAny(c.reader, sum)
	isAdmin := matchAny(c.admin, sum)
	if isAdmin == 1 {
		return Admin
	}
	if isReader == 1 {
		return Reader
	}
	return None
}

// matchAny liefert 1, wenn `sum` in `set` vorkommt, sonst 0; es durchläuft
// immer die ganze Menge.
func matchAny(set [][sha256.Size]byte, sum [sha256.Size]byte) int {
	found := 0
	for i := range set {
		found |= subtle.ConstantTimeCompare(set[i][:], sum[:])
	}
	return found
}
