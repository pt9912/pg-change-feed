package bootstrap

import (
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
)

// TestErrorStateLineNamesClassAndCode trägt die Zeile „Fehlerzustand“ des
// CLI-Berichts: Normalbetrieb ohne Klasse, Klasse und Meldungscode in eckigen
// Klammern, Klasse allein bei einem Heartbeat ohne Code. Rot färbende
// Mutation (Eingabeseite: `ErrorCode`): den Zweig mit Code entfernen — der Fall
// mit Code liest die Zeile ohne Code.
func TestErrorStateLineNamesClassAndCode(t *testing.T) {
	class, code := "schema", "PCF-E4003"
	cases := []struct {
		name   string
		result inbound.DiagnoseResult
		want   string
	}{
		{"Normalbetrieb", inbound.DiagnoseResult{}, "  Fehlerzustand: keiner (Normalbetrieb)"},
		{"Klasse und Code", inbound.DiagnoseResult{ErrorClass: &class, ErrorCode: &code}, "  Fehlerzustand: schema [PCF-E4003]"},
		{"Klasse ohne Code", inbound.DiagnoseResult{ErrorClass: &class}, "  Fehlerzustand: schema"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := errorStateLine(tc.result); got != tc.want {
				t.Fatalf("errorStateLine = %q, wollen %q", got, tc.want)
			}
		})
	}
}
