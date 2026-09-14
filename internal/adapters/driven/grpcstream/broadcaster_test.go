package grpcstream

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// testChangeMitSequenz trägt einen gültigen Change mit der übergebenen
// Sequenz; die Sequenz identifiziert ihn für die Reihenfolge-Prüfungen.
func testChangeMitSequenz(t *testing.T, sequenz int64) *model.Change {
	t.Helper()
	id := model.ChangeID(fmt.Sprintf("change-%d", sequenz))
	change, err := model.NewChange(id, "tx-1", "table-1", sequenz, model.OperationInsert, nil, []byte(`{"id":1}`), "table-1-v1")
	if err != nil {
		t.Fatalf("Change bauen: %v", err)
	}
	change.Schema = "public"
	change.Table = "orders"
	return &change
}

// testChange trägt einen gültigen Change für die Verteilungs-Tests.
func testChange(t *testing.T) *model.Change {
	t.Helper()
	return testChangeMitSequenz(t, 1)
}

// publishMitFrist ruft `Publish` in einer Goroutine auf und erzwingt eine
// Frist — ein blockierender `Publish` wird sofort als Fehlschlag sichtbar,
// statt den Testlauf bis zum globalen Timeout hängen zu lassen.
func publishMitFrist(t *testing.T, b *Broadcaster, change *model.Change) error {
	t.Helper()
	fertig := make(chan error, 1)
	go func() { fertig <- b.Publish(context.Background(), change) }()
	select {
	case err := <-fertig:
		return err
	case <-time.After(2 * time.Second):
		t.Fatal("Publish blockiert an einem registrierten, nicht lesenden Abonnenten")
		return nil
	}
}

// TestPublishVerteiltAnAlleAbonnenten trägt die erste Hälfte der Fitness
// Function aus `ADR-0060`: `Publish` verteilt denselben Change an alle zum
// Aufrufzeitpunkt registrierten Empfänger.
func TestPublishVerteiltAnAlleAbonnenten(t *testing.T) {
	b := New()
	const anzahl = 3
	kanäle := make([]<-chan *model.Change, 0, anzahl)
	for i := 0; i < anzahl; i++ {
		ch, cancel := b.Subscribe()
		defer cancel()
		kanäle = append(kanäle, ch)
	}

	change := testChange(t)
	type empfangen struct {
		index int
		wert  *model.Change
	}
	ergebnisse := make(chan empfangen, anzahl)
	for i, ch := range kanäle {
		go func(i int, ch <-chan *model.Change) {
			select {
			case wert := <-ch:
				ergebnisse <- empfangen{index: i, wert: wert}
			case <-time.After(2 * time.Second):
				ergebnisse <- empfangen{index: i}
			}
		}(i, ch)
	}

	if err := b.Publish(context.Background(), change); err != nil {
		t.Fatalf("Publish: %v (Erwartung: kein Fehler)", err)
	}

	for i := 0; i < anzahl; i++ {
		ergebnis := <-ergebnisse
		if ergebnis.wert == nil {
			t.Fatalf("Empfänger %d hat keinen Change empfangen", ergebnis.index)
		}
		if ergebnis.wert.ID != change.ID {
			t.Fatalf("Empfänger %d: Change-ID %q (Erwartung: %q)", ergebnis.index, ergebnis.wert.ID, change.ID)
		}
	}
}

// TestPublishOhneAbonnentenBlockiertNicht trägt die zweite Hälfte der
// Fitness Function aus `ADR-0060`: ein `Publish` ohne registrierten
// Empfänger blockiert nicht und liefert keinen Fehler.
func TestPublishOhneAbonnentenBlockiertNicht(t *testing.T) {
	b := New()
	fertig := make(chan error, 1)
	go func() { fertig <- b.Publish(context.Background(), testChange(t)) }()
	select {
	case err := <-fertig:
		if err != nil {
			t.Fatalf("Publish ohne Abonnenten: %v (Erwartung: kein Fehler)", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Publish ohne Abonnenten blockiert")
	}
}

// TestNichtLesenderAbonnentHaeltPublishNichtAn trägt die Kernzusage aus
// `ADR-0066` Festlegung 1: ein registrierter, gerade nicht lesender
// Empfänger hält `Publish` nicht an — auch dann nicht, wenn mehr Changes
// eintreffen, als seine begrenzte Warteschlange fasst. Ohne den
// nicht-blockierenden Send wartete der Aufruf hier bis zum Fristablauf.
func TestNichtLesenderAbonnentHaeltPublishNichtAn(t *testing.T) {
	b := New()
	_, cancel := b.Subscribe()
	defer cancel() // der Empfänger liest nie

	for i := 1; i <= queueCapacity+1; i++ {
		if err := publishMitFrist(t, b, testChangeMitSequenz(t, int64(i))); err != nil {
			t.Fatalf("Publish %d: %v (Erwartung: kein Fehler)", i, err)
		}
	}
}

// TestWarteschlangeVerwirftUeberlaufOhneAndereZuStauen trägt die Festlegungen
// 2/3 aus `ADR-0066`: die Empfangs-Warteschlange je Abonnent ist begrenzt,
// über die Kapazität hinausgehende Changes werden für diesen Empfänger
// verworfen (Drop-Newest), die eingereihten kommen in Reihenfolge an, und der
// volle Zustand staut keinen weiteren Empfänger.
func TestWarteschlangeVerwirftUeberlaufOhneAndereZuStauen(t *testing.T) {
	b := New()
	langsam, cancelLangsam := b.Subscribe()
	defer cancelLangsam()
	_, cancelSchnell := b.Subscribe() // liest ebenfalls nicht
	defer cancelSchnell()

	for i := 1; i <= queueCapacity; i++ {
		if err := publishMitFrist(t, b, testChangeMitSequenz(t, int64(i))); err != nil {
			t.Fatalf("Publish %d: %v", i, err)
		}
	}
	if err := publishMitFrist(t, b, testChangeMitSequenz(t, int64(queueCapacity+1))); err != nil {
		t.Fatalf("Publish über die Kapazität: %v", err)
	}

	for i := 1; i <= queueCapacity; i++ {
		select {
		case got := <-langsam:
			if got.Sequence != int64(i) {
				t.Fatalf("Reihenfolge: Sequenz %d (Erwartung: %d)", got.Sequence, i)
			}
		default:
			t.Fatalf("Warteschlange trug nur %d von %d Changes", i-1, queueCapacity)
		}
	}
	select {
	case got := <-langsam:
		t.Fatalf("über die Kapazität hinausgehender Change wurde zugestellt: Sequenz %d", got.Sequence)
	default:
	}
}

// TestAbmeldenIstIdempotent trägt den Vertrag der Abmelde-Funktion: ein
// zweiter Aufruf ist folgenlos — er trifft keinen Empfänger mehr und schließt
// kein zweites Mal; nach der Abmeldung wird der Kanal nicht mehr bedient.
func TestAbmeldenIstIdempotent(t *testing.T) {
	b := New()
	_, cancel := b.Subscribe()
	cancel()
	cancel()

	if len(b.subs) != 0 {
		t.Fatalf("Broadcaster führt nach der Abmeldung noch %d Empfänger", len(b.subs))
	}

	fertig := make(chan error, 1)
	go func() { fertig <- b.Publish(context.Background(), testChange(t)) }()
	select {
	case err := <-fertig:
		if err != nil {
			t.Fatalf("Publish nach doppelter Abmeldung: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Publish blockiert nach doppelter Abmeldung")
	}
}

// TestPublishOhneChangeLiefertFehler trägt den Fehler-Ausgang des Ports:
// ein fehlender Change ist ein ungültiger Aufruf, kein stiller
// Fehlerschluck (`outbound.ErrChangeStream`).
func TestPublishOhneChangeLiefertFehler(t *testing.T) {
	b := New()
	err := b.Publish(context.Background(), nil)
	if !errors.Is(err, outbound.ErrChangeStream) {
		t.Fatalf("Publish(nil): %v (Erwartung: %v)", err, outbound.ErrChangeStream)
	}
}

// TestPublishMitBeendetemKontext trägt den Vorab-Ausgang aus `ADR-0066`
// Festlegung 1: ein bereits beendetes `ctx` beendet den Aufruf mit dessen
// Fehler, statt den Change noch zu verteilen.
func TestPublishMitBeendetemKontext(t *testing.T) {
	b := New()
	ch, cancel := b.Subscribe()
	defer cancel()

	ctx, abbrechen := context.WithCancel(context.Background())
	abbrechen()

	err := b.Publish(ctx, testChange(t))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Publish mit beendetem Kontext: %v (Erwartung: %v)", err, context.Canceled)
	}
	select {
	case got := <-ch:
		t.Fatalf("beendeter Kontext verteilte dennoch einen Change: %v", got)
	default:
	}
}
