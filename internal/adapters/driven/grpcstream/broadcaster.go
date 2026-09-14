// Package grpcstream trägt den `Broadcaster` als Driven-Implementierung des
// `ChangeStreamPort` (`ADR-0060` Teilfrage 2/5, `ADR-0066`): ein
// In-Prozess-Fan-out ohne externes System — jeder aktive Live-Stream
// registriert sich über `Subscribe`, `Publish` verteilt jeden Change an alle
// registrierten Empfänger. Die Verteilung ist nicht-blockierend: jeder
// Empfänger trägt eine begrenzte Empfangs-Warteschlange, deren Überlauf
// verworfen wird, der Erzeuger hält nie auf einen Empfänger an (`ADR-0066`).
// Die Nachvollziehbarkeit bleibt beim bestehenden Lesezugriffsweg
// (`LH-FA-REA-001` ff.); der Broadcaster trägt keinen Zustand über den
// Verteilungszeitpunkt hinaus (`ADR-0060` Teilfrage 3, Fire-and-Forget ohne
// Replay).
package grpcstream

import (
	"context"
	"fmt"
	"sync"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

var _ outbound.ChangeStreamPort = (*Broadcaster)(nil)

// queueCapacity ist die Kapazität der Empfangs-Warteschlange je Empfänger
// (`ADR-0066` Festlegung 2). Sie ist eine Einstellgröße, keine Zusicherung:
// die Warteschlange absorbiert einen kurzen Erzeuger-Burst, während der
// Abonnent zwischen zwei Reads steht (etwa in einem blockierenden
// `stream.Send`), und macht keine Zustellzusage — über sie hinausgehende
// Changes werden für diesen Empfänger verworfen. 64 trägt bei einem Zeiger je
// Eintrag rund einen halben Kilobyte je Abonnent, liegt deutlich über der
// Spanne eines gewöhnlichen Scheduling-/Flow-Control-Zwischenraums und bleibt
// trotzdem eine kleine, feste Schranke. Eine Zustellzusage über Fire-and-Forget
// hinaus bräuchte eine eigene Entscheidung (`ADR-0066`
// Re-Evaluierungs-Trigger 2).
const queueCapacity = 64

// Broadcaster ist ein nebenläufigkeitssicherer In-Prozess-Fan-out
// (`ADR-0060` Teilfrage 2): beliebig viele Empfänger, deren Registrierung und
// Abmeldung gegen laufende Verteilungen über `mu` synchronisiert ist. Jeder
// Empfänger trägt eine begrenzte Empfangs-Warteschlange (`queueCapacity`);
// `Publish` legt einen Change nicht-blockierend ab und verwirft ihn für einen
// Empfänger, dessen Warteschlange voll ist (`ADR-0066`). Der Erzeuger hält
// damit nie auf einen Empfänger an.
type Broadcaster struct {
	mu   sync.Mutex
	next int
	subs map[int]chan *model.Change
}

// New legt einen leeren Broadcaster an.
func New() *Broadcaster {
	return &Broadcaster{subs: map[int]chan *model.Change{}}
}

// Subscribe registriert einen Empfänger und liefert dessen begrenzte
// Empfangs-Warteschlange samt Abmelde-Funktion (`ADR-0060` Teilfrage 2,
// `ADR-0066` Festlegung 2). Jeder Aufruf registriert einen eigenen Empfänger;
// `Publish` verteilt an alle. Die Abmelde-Funktion ist idempotent, entfernt
// genau diesen Empfänger und verwirft dabei die noch in seiner Warteschlange
// liegenden Changes. Der gelieferte Kanal wird bewusst nie geschlossen: ein
// Schließen träfe einen gleichzeitig laufenden Sende-Versuch in `Publish` und
// panikte, `select` schützt davor nicht. Nach der Abmeldung werden dem Aufrufer
// keine weiteren Changes zugestellt; er liest den gelieferten Kanal nicht
// weiter.
func (b *Broadcaster) Subscribe() (<-chan *model.Change, func()) {
	changes := make(chan *model.Change, queueCapacity)
	b.mu.Lock()
	id := b.next
	b.next++
	b.subs[id] = changes
	b.mu.Unlock()

	var once sync.Once
	cancel := func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.subs, id)
			b.mu.Unlock()
		})
	}
	return changes, cancel
}

// Publish verteilt einen Change an alle zum Aufrufzeitpunkt registrierten
// Empfänger (`ChangeStreamPort`, `ADR-0060` Teilfrage 3, `ADR-0066`). Der
// Aufruf blockiert nie auf einen Empfänger: er legt den Change je Empfänger
// nicht-blockierend in dessen begrenzte Empfangs-Warteschlange; ist sie voll,
// wird der eintreffende Change für diesen Empfänger verworfen (Drop-Newest,
// nicht nachgeliefert), die übrigen Empfänger bleiben unberührt. Ohne
// registrierten Empfänger unterbleibt jede Übergabe — der Aufruf blockiert
// dann nicht und liefert keinen Fehler (Fire-and-Forget). Ein bereits
// beendetes `ctx` beendet den Aufruf vorab mit dessen Fehler; der
// `Publish`-Fehler entsteht sonst nur aus einem ungültigen Aufruf
// (`outbound.ErrChangeStream`).
func (b *Broadcaster) Publish(ctx context.Context, change *model.Change) error {
	if change == nil {
		return fmt.Errorf("%w: Publish ohne Change", outbound.ErrChangeStream)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	b.mu.Lock()
	targets := make([]chan *model.Change, 0, len(b.subs))
	for _, changes := range b.subs {
		targets = append(targets, changes)
	}
	b.mu.Unlock()

	for _, changes := range targets {
		select {
		case changes <- change:
		default:
		}
	}
	return nil
}
