// Package grpcstream trägt den `Broadcaster` als Driven-Implementierung des
// `ChangeStreamPort` (`ADR-0060` Teilfrage 2/5): ein In-Prozess-Fan-out ohne
// externes System — jeder aktive Live-Stream registriert sich über
// `Subscribe`, `Publish` verteilt jeden Change an alle registrierten
// Empfänger. Die Nachvollziehbarkeit bleibt beim bestehenden
// Lesezugriffsweg (`LH-FA-REA-001` ff.); der Broadcaster trägt keinen
// Zustand über den Verteilungszeitpunkt hinaus (`ADR-0060` Teilfrage 3,
// Fire-and-Forget ohne Replay).
package grpcstream

import (
	"context"
	"fmt"
	"sync"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

var _ outbound.ChangeStreamPort = (*Broadcaster)(nil)

// subscriber trägt einen registrierten Empfänger: `changes` ist der von
// `Subscribe` gelieferte, ungepufferte Kanal, `done` wird beim Abmelden
// geschlossen und löst einen noch wartenden `Publish`-Aufruf für genau
// diesen Empfänger. Der Kanal `changes` wird bewusst nie geschlossen: ein
// Empfänger kann ihn nach seinem `cancel` nicht mehr lesen, ein Schließen
// durch `cancel` träfe aber einen gleichzeitig laufenden Sende-Versuch in
// `Publish` und panikte; die Abmeldung trägt deshalb ausschließlich `done`.
type subscriber struct {
	changes chan *model.Change
	done    chan struct{}
}

// Broadcaster ist ein nebenläufigkeitssicherer In-Prozess-Fan-out
// (`ADR-0060` Teilfrage 2): beliebig viele Empfänger, deren Registrierung
// und Abmeldung gegen laufende Verteilungen über `mu` synchronisiert ist.
// Die Kanalmenge ist ungepuffert (`ADR-0060` Teilfrage 3, kein Puffer):
// `Publish` übergibt jeden Change direkt an einen empfangenden Leser; ein
// Empfänger ohne bereitliegenden Leser hält die Übergabe an, bis er liest
// oder sich abmeldet.
type Broadcaster struct {
	mu   sync.Mutex
	next int
	subs map[int]*subscriber
}

// New legt einen leeren Broadcaster an.
func New() *Broadcaster {
	return &Broadcaster{subs: map[int]*subscriber{}}
}

// Subscribe registriert einen Empfänger und liefert dessen Kanal samt
// Abmelde-Funktion (`ADR-0060` Teilfrage 2). Jeder Aufruf registriert einen
// eigenen Empfänger; `Publish` verteilt an alle. Die Abmelde-Funktion ist
// idempotent, entfernt genau diesen Empfänger und löst einen noch wartenden
// `Publish`-Aufruf für ihn. Nach der Abmeldung darf der Aufrufer den
// gelieferten Kanal nicht mehr lesen — er wird nicht geschlossen.
func (b *Broadcaster) Subscribe() (<-chan *model.Change, func()) {
	sub := &subscriber{changes: make(chan *model.Change), done: make(chan struct{})}
	b.mu.Lock()
	id := b.next
	b.next++
	b.subs[id] = sub
	b.mu.Unlock()

	var once sync.Once
	cancel := func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.subs, id)
			b.mu.Unlock()
			close(sub.done)
		})
	}
	return sub.changes, cancel
}

// Publish verteilt einen Change an alle zum Aufrufzeitpunkt registrierten
// Empfänger (`ChangeStreamPort`, `ADR-0060` Teilfrage 3). Ohne registrierten
// Empfänger unterbleibt jede Übergabe — der Aufruf blockiert dann nicht und
// liefert keinen Fehler (Fire-and-Forget, kein Puffer). Mit einem
// registrierten, aber gerade nicht lesenden Empfänger hält die Übergabe an
// diesen an, bis er liest, sich abmeldet oder `ctx` endet: ohne einen
// solchen Rückhalt ginge ein Change bei jedem noch so kurzen
// Empfangs-Zwischenraum verloren, statt an den verbundenen Empfänger zu
// gelangen. Ein abgemeldeter Empfänger löst die Übergabe sofort über sein
// `done`-Signal, ein beendeter `ctx` über `ctx.Done()`.
func (b *Broadcaster) Publish(ctx context.Context, change *model.Change) error {
	if change == nil {
		return fmt.Errorf("%w: Publish ohne Change", outbound.ErrChangeStream)
	}
	b.mu.Lock()
	targets := make([]*subscriber, 0, len(b.subs))
	for _, sub := range b.subs {
		targets = append(targets, sub)
	}
	b.mu.Unlock()

	for _, sub := range targets {
		select {
		case sub.changes <- change:
		case <-sub.done:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
