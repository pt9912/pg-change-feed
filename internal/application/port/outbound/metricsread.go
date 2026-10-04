package outbound

import (
	"context"
	"errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrMetricsRead trägt das Scheitern des Lesens der Kennzahlen-Sicht
// (`ADR-0149`): ein Lesefehler bleibt über `errors.Is` erkennbar, ohne
// Treibertyp. Der Fehler trägt keine Warnung-Klasse; der Export ist ein
// Zugriffsweg neben der Sicht und ändert weder Health noch `error_class`.
var ErrMetricsRead = errors.New("Kennzahlen-Sicht nicht lesbar")

// ErrMetricExport trägt das Scheitern der Übertragung an den Empfänger: jeder
// Status außerhalb `2xx`, jeder Verbindungs- und jeder Zeitfehler.
var ErrMetricExport = errors.New("Kennzahlen-Übertragung fehlgeschlagen")

// FailureDetail ist die Eigenschaft eines Fehlers, einen Text zu tragen, der
// bedenkenlos in eine Log-Zeile gehört: er ist aus Art und Statuscode
// gebaut und trägt weder URL noch Header-Wert. Ein Aufrufer, der die Ursache
// protokolliert, liest ihn; der rohe Fehlertext einer Bibliothek gehört nicht in
// eine Log-Zeile.
type FailureDetail interface {
	FailureDetail() string
}

// MetricValue trägt den Wert einer Kennzahl in zwei Lesarten: `Float` ist
// immer belegt; `Int` ist belegt und `IsInt` wahr, wenn der Wert der Sicht
// eine ganze Zahl im Bereich von `int64` ist. Positionen und Zähler behalten
// so ihre volle 64-Bit-Genauigkeit, die ein `float64` ab 2^53 verliert.
type MetricValue struct {
	Float float64
	Int   int64
	IsInt bool
}

// MetricSample ist eine Zeile der Sicht der Betriebsschnittstelle
// (`SPEC-009`): Name, optionales Label und Wert. `Label` ist leer, wenn die
// Zeile keines trägt (`NULL` in der Sicht); die Bedeutung des Labels
// (Consumer oder Fehlerklasse) folgt aus dem Namen. `At` ist der
// Messzeitpunkt, den der Use Case setzt.
type MetricSample struct {
	Name  string
	Label string
	Value MetricValue
	At    model.TimePoint
}

// MetricsReadPort liest die Zeilen der Sicht `cdc.metrics`. Ein Adapter
// liest unter der Leserolle (`cdc_reader`, keine zusätzlichen Rechte) alle
// Zeilen in einem Aufruf, ohne eigenen Zustand zwischen zwei Aufrufen. Er
// beachtet die Frist des übergebenen Kontexts und meldet jeden Fehler über
// `ErrMetricsRead`.
type MetricsReadPort interface {
	Read(ctx context.Context) ([]MetricSample, error)
}
