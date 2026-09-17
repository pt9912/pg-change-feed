package main

import "testing"

// lines liefert eine Zeilenquelle über einen Ausschnitt; jeder Aufruf gibt die
// nächste Zeile zurück, nach der letzten `ok == false` — dieselbe Form, die
// `main` aus dem `bufio.Scanner` über den Response-Body bildet.
func lines(all []string) func() (string, bool) {
	i := 0
	return func() (string, bool) {
		if i >= len(all) {
			return "", false
		}
		line := all[i]
		i++
		return line, true
	}
}

// TestStreamURL prüft den Aufbau der Stream-Adresse (`LH-FA-SST-008`): der
// Endpunkt ist `/changes/stream`, der Host kommt aus `CDC_HTTP_ADDR`.
func TestStreamURL(t *testing.T) {
	got := StreamURL("feed:8080")
	want := "http://feed:8080/changes/stream"
	if got != want {
		t.Fatalf("StreamURL = %q, want %q", got, want)
	}
}

// TestReadEventReadsNameAndPayload prüft das Frame-Lesen (`ADR-0061`): der
// `event:`-Name und die `data:`-Nutzlast kommen aus ihren Zeilen.
func TestReadEventReadsNameAndPayload(t *testing.T) {
	ev, err := readEvent(lines([]string{
		"event: change",
		`data: {"change_id":"c-1","table":"orders"}`,
		"",
	}))
	if err != nil {
		t.Fatalf("readEvent: %v", err)
	}
	if ev.Name != "change" {
		t.Fatalf("Name = %q, want %q", ev.Name, "change")
	}
	want := `{"change_id":"c-1","table":"orders"}`
	if ev.Data != want {
		t.Fatalf("Data = %q, want %q", ev.Data, want)
	}
}

// TestReadEventStopsAtFrameBoundary prüft die Frame-Grenze: die Leerzeile
// schließt ein Frame ab, zwei aufeinanderfolgende Events werden nicht zu
// einem verschmolzen.
func TestReadEventStopsAtFrameBoundary(t *testing.T) {
	next := lines([]string{
		"event: change",
		`data: {"change_id":"c-1"}`,
		"",
		"event: change",
		`data: {"change_id":"c-2"}`,
		"",
	})
	first, err := readEvent(next)
	if err != nil {
		t.Fatalf("erstes readEvent: %v", err)
	}
	second, err := readEvent(next)
	if err != nil {
		t.Fatalf("zweites readEvent: %v", err)
	}
	if first.Data != `{"change_id":"c-1"}` {
		t.Fatalf("erstes Data = %q", first.Data)
	}
	if second.Data != `{"change_id":"c-2"}` {
		t.Fatalf("zweites Data = %q", second.Data)
	}
}

// TestReadEventReportsExhaustedSource prüft das Ende der Quelle: ein
// begonnenes, aber nicht abgeschlossenes Frame ist kein Event.
func TestReadEventReportsExhaustedSource(t *testing.T) {
	next := lines([]string{
		"event: change",
		`data: {"change_id":"c-1"}`,
	})
	if _, err := readEvent(next); err == nil {
		t.Fatal("readEvent lieferte ein unvollständiges Frame als Event")
	}
}
