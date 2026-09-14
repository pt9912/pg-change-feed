package http

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// fakeChangeSubscriber trägt eine In-Memory-Fälschung des
// `changeSubscriber` (`ADR-0030`): Whitebox-Test des Adapters ohne den
// Driven-Broadcaster — der Adapter-Test darf kein Adapter-Paket importieren
// (Richtungskonvention `spec/architecture.md` §1, maschinell über
// `a-check`). `ready` schließt `Subscribe` einmalig und macht die
// Registrierung beobachtbar, `cancelled` die Abmeldung. Die
// Empfangs-Warteschlange ist begrenzt wie die des `Broadcaster` (`ADR-0066`).
type fakeChangeSubscriber struct {
	changes    chan *model.Change
	ready      chan struct{}
	readyOnce  sync.Once
	cancelled  chan struct{}
	cancelOnce sync.Once
}

func newFakeChangeSubscriber(capacity int) *fakeChangeSubscriber {
	return &fakeChangeSubscriber{
		changes:   make(chan *model.Change, capacity),
		ready:     make(chan struct{}),
		cancelled: make(chan struct{}),
	}
}

func (f *fakeChangeSubscriber) Subscribe() (<-chan *model.Change, func()) {
	f.readyOnce.Do(func() { close(f.ready) })
	return f.changes, func() { f.cancelOnce.Do(func() { close(f.cancelled) }) }
}

// neuerStreamTestChange baut einen Change mit vollständigem Inhalt (Schema,
// Tabelle, beide Row Images), wie ihn der Capture-Pfad veröffentlicht.
func neuerStreamTestChange(t *testing.T) model.Change {
	t.Helper()
	change, err := model.NewChange("change-1", "tx-1", "table-1", 2, model.OperationInsert,
		[]byte(`{"id":1}`), []byte(`{"id":1,"name":"StreamE2ESentinel"}`), "table-1-v1")
	if err != nil {
		t.Fatalf("Change bauen: %v", err)
	}
	change.Schema = "public"
	change.Table = "orders"
	return change
}

// newTestSSEServer verdrahtet den Adapter mit dem übergebenen Subscriber auf
// einem `httptest`-Server — kein realer Port (`AGENTS.md` §3.1 im Testlauf).
func newTestSSEServer(t *testing.T, subscriber changeSubscriber) *httptest.Server {
	t.Helper()
	srv := New(Config{
		Addr:        "unused:0",
		TokenReader: testReaderToken,
		TokenAdmin:  testAdminToken,
		Subscriber:  subscriber,
	})
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts
}

// streamMitToken öffnet den Stream; ein leerer Token lässt den
// `Authorization`-Header weg.
func streamMitToken(t *testing.T, ts *httptest.Server, ctx context.Context, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/changes/stream", nil)
	if err != nil {
		t.Fatalf("Request bauen: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Request senden: %v", err)
	}
	return resp
}

// readSSEEvent liest ein vollständiges SSE-Event (bis zur Leerzeile) mit
// Frist — ein hängender Stream ist ein Fehlschlag, kein Warten.
func readSSEEvent(t *testing.T, reader *bufio.Reader) (event, data string) {
	t.Helper()
	type ergebnis struct {
		event string
		data  string
		err   error
	}
	out := make(chan ergebnis, 1)
	go func() {
		var ev, da string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				out <- ergebnis{err: err}
				return
			}
			line = strings.TrimRight(line, "\r\n")
			if line == "" {
				out <- ergebnis{event: ev, data: da}
				return
			}
			switch {
			case strings.HasPrefix(line, "event: "):
				ev = strings.TrimPrefix(line, "event: ")
			case strings.HasPrefix(line, "data: "):
				da = strings.TrimPrefix(line, "data: ")
			}
		}
	}()
	select {
	case r := <-out:
		if r.err != nil {
			t.Fatalf("SSE-Event lesen: %v", r.err)
		}
		return r.event, r.data
	case <-time.After(3 * time.Second):
		t.Fatal("kein SSE-Event innerhalb der Frist")
		return "", ""
	}
}

// TestStreamOhneTokenEndetMit401 trägt die erste Hälfte der Fitness Function
// aus `ADR-0061` Teilfrage 4: ein Aufruf ohne Bearer-Token endet mit `401`,
// bevor der Handler läuft und damit vor jedem geschriebenen SSE-Event
// (`LH-FA-SST-008` Negative).
func TestStreamOhneTokenEndetMit401(t *testing.T) {
	subscriber := newFakeChangeSubscriber(1)
	ts := newTestSSEServer(t, subscriber)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp := streamMitToken(t, ts, ctx, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Status: %d (Erwartung: 401)", resp.StatusCode)
	}
	select {
	case <-subscriber.ready:
		t.Fatal("der Handler registrierte trotz 401 eine Subskription am Broadcaster")
	default:
	}
}

// TestStreamUnbekannterTokenEndetMit401 trägt denselben Pfad für einen
// Token, der keiner konfigurierten Klasse entspricht — nicht nur einen
// fehlenden Header.
func TestStreamUnbekannterTokenEndetMit401(t *testing.T) {
	ts := newTestSSEServer(t, newFakeChangeSubscriber(1))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp := streamMitToken(t, ts, ctx, "unbekannt")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Status: %d (Erwartung: 401)", resp.StatusCode)
	}
}

// TestStreamReaderTokenOeffnetTraegtChange trägt die zweite Hälfte der
// Fitness Function: ein gültiges `reader`-Token öffnet den Stream als
// `text/event-stream`, und ein über den Broadcaster verteilter Change
// erreicht den Client als vollständiges Event (`LH-FA-SST-008` Happy Path,
// `SPEC-018`).
func TestStreamReaderTokenOeffnetTraegtChange(t *testing.T) {
	traegtTokenOeffnetStreamUndTraegtChange(t, testReaderToken)
}

// TestStreamAdminTokenOeffnetTraegtChange trägt denselben Pfad für die
// administrative Klasse — ein `admin`-Token deckt die lesende
// Streaming-Fähigkeit ab (`ADR-0061` Teilfrage 4).
func TestStreamAdminTokenOeffnetTraegtChange(t *testing.T) {
	traegtTokenOeffnetStreamUndTraegtChange(t, testAdminToken)
}

func traegtTokenOeffnetStreamUndTraegtChange(t *testing.T, token string) {
	t.Helper()
	subscriber := newFakeChangeSubscriber(4)
	ts := newTestSSEServer(t, subscriber)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp := streamMitToken(t, ts, ctx, token)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status: %d (Erwartung: 200)", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("Content-Type: %q (Erwartung: text/event-stream)", ct)
	}
	select {
	case <-subscriber.ready:
	case <-time.After(3 * time.Second):
		t.Fatal("der Stream hat sich nicht am Broadcaster registriert")
	}

	change := neuerStreamTestChange(t)
	subscriber.changes <- &change

	event, data := readSSEEvent(t, bufio.NewReader(resp.Body))
	if event != "change" {
		t.Fatalf("event-Typ: %q (Erwartung: %q)", event, "change")
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(data), &got); err != nil {
		t.Fatalf("Event-Daten sind kein JSON: %v (%q)", err, data)
	}
	if got["change_id"] != "change-1" || got["transaction_id"] != "tx-1" ||
		got["source_table_id"] != "table-1" || got["sequence"] != float64(2) ||
		got["operation"] != string(model.OperationInsert) ||
		got["schema"] != "public" || got["table"] != "orders" ||
		got["schema_version"] != "table-1-v1" {
		t.Fatalf("Event trägt nicht die erwarteten Felder: %v", got)
	}
	newImage, ok := got["new_image"].(map[string]any)
	if !ok || newImage["name"] != "StreamE2ESentinel" {
		t.Fatalf("new_image trägt nicht den erwarteten Inhalt: %v", got["new_image"])
	}
	oldImage, ok := got["old_image"].(map[string]any)
	if !ok || oldImage["id"] != float64(1) {
		t.Fatalf("old_image trägt nicht den erwarteten Inhalt: %v", got["old_image"])
	}
}

// TestStreamLeeresRowImageTraegtNull trägt die Grenze des Nachrichtenschemas
// (`LH-FA-CAP-008` Boundary): ein fehlendes Bild bleibt gültiges JSON und
// wird als `null` übertragen, nicht als leerer Wert.
func TestStreamLeeresRowImageTraegtNull(t *testing.T) {
	subscriber := newFakeChangeSubscriber(4)
	ts := newTestSSEServer(t, subscriber)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp := streamMitToken(t, ts, ctx, testReaderToken)
	defer resp.Body.Close()
	select {
	case <-subscriber.ready:
	case <-time.After(3 * time.Second):
		t.Fatal("der Stream hat sich nicht am Broadcaster registriert")
	}

	change, err := model.NewChange("change-2", "tx-2", "table-1", 1, model.OperationDelete,
		[]byte(`{"id":1}`), nil, "table-1-v1")
	if err != nil {
		t.Fatalf("Change bauen: %v", err)
	}
	subscriber.changes <- &change

	_, data := readSSEEvent(t, bufio.NewReader(resp.Body))
	var got map[string]any
	if err := json.Unmarshal([]byte(data), &got); err != nil {
		t.Fatalf("Event-Daten sind kein JSON: %v (%q)", err, data)
	}
	value, present := got["new_image"]
	if !present || value != nil {
		t.Fatalf("fehlendes new_image ist nicht null: %q", data)
	}
}

// TestStreamOhneBroadcasterAntwortetMit503 trägt den Fehler-Ausgang einer
// fehlenden Verdrahtung (`ADR-0061` Teilfrage 5): ohne Broadcaster antwortet
// der Endpunkt sichtbar mit `503`, statt eine leere, nie endende Verbindung
// offenzuhalten.
func TestStreamOhneBroadcasterAntwortetMit503(t *testing.T) {
	ts := newTestSSEServer(t, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp := streamMitToken(t, ts, ctx, testReaderToken)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("Status: %d (Erwartung: 503)", resp.StatusCode)
	}
}

// TestStreamOhneVerbundenenClientBlockiertNicht trägt die
// Fire-and-Forget-Hälfte der Fitness Function aus `ADR-0061`: ohne
// verbundenen SSE-Client liegt am Adapter keine Subskription, auf die ein
// `Publish` warten könnte — der nicht-blockierende Handoff an die begrenzte
// Empfangs-Warteschlange kehrt sofort zurück (`ADR-0066` Festlegung 1/3).
func TestStreamOhneVerbundenenClientBlockiertNicht(t *testing.T) {
	subscriber := newFakeChangeSubscriber(1)
	_ = newTestSSEServer(t, subscriber)

	select {
	case <-subscriber.ready:
		t.Fatal("der Adapter registriert ohne laufenden Request eine Subskription")
	default:
	}

	change := neuerStreamTestChange(t)
	fertig := make(chan struct{})
	go func() {
		// Der Handoff des Broadcaster: nicht-blockierender Send, der bei
		// voller Warteschlange verwirft statt zu warten.
		select {
		case subscriber.changes <- &change:
		default:
		}
		close(fertig)
	}()
	select {
	case <-fertig:
	case <-time.After(3 * time.Second):
		t.Fatal("Zustellung ohne verbundenen SSE-Client blockierte den Erzeuger")
	}
}

// TestStreamAbgemeldeterClientGibtSubskriptionFrei trägt die zweite Hälfte
// derselben Zusage: endet die Verbindung, gibt der Handler seine Subskription
// frei — ein getrennter Client hält über die begrenzte Warteschlange weder
// einen Erzeuger noch Speicher.
func TestStreamAbgemeldeterClientGibtSubskriptionFrei(t *testing.T) {
	subscriber := newFakeChangeSubscriber(1)
	ts := newTestSSEServer(t, subscriber)
	ctx, cancel := context.WithCancel(context.Background())

	resp := streamMitToken(t, ts, ctx, testReaderToken)
	defer resp.Body.Close()
	select {
	case <-subscriber.ready:
	case <-time.After(3 * time.Second):
		cancel()
		t.Fatal("der Stream hat sich nicht am Broadcaster registriert")
	}

	cancel()
	select {
	case <-subscriber.cancelled:
	case <-time.After(3 * time.Second):
		t.Fatal("der Handler gab seine Subskription nach dem Verbindungsende nicht frei")
	}
}
