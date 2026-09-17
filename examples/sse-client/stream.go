package main

import (
	"io"
	"net/url"
	"strings"
)

// Event trägt ein Frame des Change-Streams (`LH-FA-SST-008`): den Namen aus
// der `event:`-Zeile und die Nutzlast aus der `data:`-Zeile. Der Server
// schreibt je Change ein JSON-Objekt mit den zehn Nachrichtenfeldern in die
// `data:`-Zeile (`ADR-0061`).
type Event struct {
	Name string
	Data string
}

// StreamURL baut die Adresse des SSE-Endpunkts (`LH-FA-SST-008`): ein `GET`
// auf `/changes/stream`. `addr` ist die Horch-Adresse des Feed-Containers
// (`CDC_HTTP_ADDR`, Form `host:port`).
func StreamURL(addr string) string {
	u := url.URL{Scheme: "http", Host: addr, Path: "/changes/stream"}
	return u.String()
}

// readEvent liest ein vollständiges Frame über next und liefert es zurück.
// Ein Frame endet mit der Leerzeile, die der Server nach der `event:`- und der
// `data:`-Zeile schreibt (`ADR-0061`); sie trennt zwei aufeinanderfolgende
// Events. Ist die Quelle vor dem Frame-Abschluss erschöpft, liefert readEvent
// `io.EOF`, und ein begonnenes Frame wird verworfen — ein unvollständiges
// Frame ist kein Event.
func readEvent(next func() (string, bool)) (Event, error) {
	var ev Event
	for {
		line, ok := next()
		if !ok {
			return Event{}, io.EOF
		}
		if line == "" {
			if ev.Name == "" && ev.Data == "" {
				continue
			}
			return ev, nil
		}
		switch {
		case strings.HasPrefix(line, "event: "):
			ev.Name = strings.TrimPrefix(line, "event: ")
		case strings.HasPrefix(line, "data: "):
			ev.Data = strings.TrimPrefix(line, "data: ")
		}
	}
}
