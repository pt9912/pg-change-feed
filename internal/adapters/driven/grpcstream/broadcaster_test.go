package grpcstream

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// testChange trägt einen gültigen Change für die Verteilungs-Tests.
func testChange(t *testing.T) *model.Change {
	t.Helper()
	change, err := model.NewChange("change-1", "tx-1", "table-1", 1, model.OperationInsert, nil, []byte(`{"id":1}`), "table-1-v1")
	if err != nil {
		t.Fatalf("Change bauen: %v", err)
	}
	change.Schema = "public"
	change.Table = "orders"
	return &change
}

// TestPublishVerteiltAnAlleAbonnenten trägt die erste Hälfte der Fitness
// Function aus `ADR-0060`: `Publish` verteilt denselben Change an alle zum
// Aufrufzeitpunkt registrierten Empfänger. Die Leser laufen nebenläufig,
// weil die Übergabe ungepuffert an genau einen lesenden Empfänger geht.
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

// TestAbmeldenLoestWartendenPublish trägt den Rückhalt der Übergabe: ein
// registrierter, aber nicht lesender Empfänger hält `Publish` an, seine
// Abmeldung löst den Aufruf ohne Fehler. Ohne diesen Pfad bliebe der
// Aufrufer eines abgemeldeten Streams für immer hängen.
func TestAbmeldenLoestWartendenPublish(t *testing.T) {
	b := New()
	_, cancel := b.Subscribe()

	fertig := make(chan error, 1)
	go func() { fertig <- b.Publish(context.Background(), testChange(t)) }()

	cancel()
	select {
	case err := <-fertig:
		if err != nil {
			t.Fatalf("Publish nach Abmeldung: %v (Erwartung: kein Fehler)", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Abmeldung löst den wartenden Publish-Aufruf nicht")
	}
}

// TestAbmeldenIstIdempotent trägt den Vertrag der Abmelde-Funktion: ein
// zweiter Aufruf ist folgenlos — er trifft keinen Empfänger mehr und
// schließt kein zweites Mal.
func TestAbmeldenIstIdempotent(t *testing.T) {
	b := New()
	ch, cancel := b.Subscribe()
	cancel()
	cancel()

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
	// Der abgemeldete Kanal trägt nichts mehr.
	select {
	case wert := <-ch:
		t.Fatalf("abgemeldeter Kanal trug %v", wert)
	default:
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

// TestPublishMitBeendetemKontext trägt den Ausweg aus einer wartenden
// Übergabe: endet der Kontext, kehrt `Publish` mit dessen Fehler zurück,
// statt für immer zu warten.
func TestPublishMitBeendetemKontext(t *testing.T) {
	b := New()
	_, cancel := b.Subscribe()
	defer cancel()

	ctx, abbrechen := context.WithCancel(context.Background())
	fertig := make(chan error, 1)
	go func() { fertig <- b.Publish(ctx, testChange(t)) }()

	abbrechen()
	select {
	case err := <-fertig:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Publish mit beendetem Kontext: %v (Erwartung: %v)", err, context.Canceled)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("beendeter Kontext löst den wartenden Publish-Aufruf nicht")
	}
}
