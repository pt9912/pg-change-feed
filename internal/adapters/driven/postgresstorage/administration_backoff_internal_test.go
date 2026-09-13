package postgresstorage

import (
	"testing"
	"time"
)

// TestNextAdministrationReconnectBackoffDoublesFromZeroAndCaps trägt den
// Backoff-Verlauf des `LISTEN`-Wiederverbindungspfads (Review-Finding F-3,
// `review-slice-037.md`): 0 (noch kein Fehlschlag) springt auf die
// Initial-Backoff, jeder weitere Fehlschlag verdoppelt, die Obergrenze
// deckelt. Rot färbende Mutation: `doubled > administrationReconnectMaxBackoff`
// durch `doubled >= administrationReconnectMaxBackoff` ersetzen — dann
// kippt der Deckel einen Schritt zu früh, und der Fall „knapp unter der
// Obergrenze" in dieser Tabelle liefert einen anderen Wert.
func TestNextAdministrationReconnectBackoffDoublesFromZeroAndCaps(t *testing.T) {
	cases := []struct {
		name    string
		current time.Duration
		want    time.Duration
	}{
		{"kein Fehlschlag bisher", 0, administrationReconnectInitialBackoff},
		{"negativer Ausgangswert (Zero-Value-Absicherung)", -1, administrationReconnectInitialBackoff},
		{"ein Fehlschlag", administrationReconnectInitialBackoff, 2 * administrationReconnectInitialBackoff},
		{"knapp unter der Obergrenze verdoppelt noch", administrationReconnectMaxBackoff - 1, administrationReconnectMaxBackoff},
		{"an der Obergrenze bleibt gedeckelt", administrationReconnectMaxBackoff, administrationReconnectMaxBackoff},
		{"über der Obergrenze (sollte strukturell nicht vorkommen) bleibt gedeckelt", administrationReconnectMaxBackoff * 2, administrationReconnectMaxBackoff},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := nextAdministrationReconnectBackoff(c.current); got != c.want {
				t.Fatalf("nextAdministrationReconnectBackoff(%s) = %s, wollen %s", c.current, got, c.want)
			}
		})
	}
}
