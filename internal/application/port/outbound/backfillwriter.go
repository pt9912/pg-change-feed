package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// BackfillWriterPort trägt die Fähigkeit, die Blöcke eines Runs in **einer**
// Store-Transaktion zu schreiben (`ARC-004`, Fähigkeits-Port je `ADR-0034`;
// `ADR-0111` Teilfrage 4). Der bestehende `ChangeStorePort.PersistTransaction`
// hält eine ganze Transaktion im Speicher und trägt einen Bestand nicht;
// dieser Port nimmt die Blöcke einzeln entgegen, damit der Aufrufer höchstens
// einen Block hält (`LH-FA-CAP-006.a`).
type BackfillWriterPort interface {
	// Begin öffnet die Schreibtransaktion des Runs. Bis zum Commit ist
	// nichts für Leser sichtbar.
	Begin(ctx context.Context, run model.BackfillRun) (BackfillTransaction, error)
}

// BackfillTransaction ist die eine offene Schreibtransaktion eines Runs.
// Sie wird von einer Goroutine benutzt.
type BackfillTransaction interface {
	// AppendBlock hängt einen Block an: eine committed, synthetische
	// Transaktion (`model.ChangeTransaction`) auf der Position `X` des Runs.
	// Der Port hält den Block nach der Rückkehr nicht.
	AppendBlock(ctx context.Context, block *model.ChangeTransaction) error

	// Commit schreibt die Run-Zeile im Endzustand `completed` und
	// committet sie zusammen mit allen angehängten Blöcken — ein Commit,
	// eine Einheit. Persistenzfehler tragen `ErrBackfillStorage`.
	Commit(ctx context.Context, run model.BackfillRun) error

	// Rollback verwirft alle angehängten Blöcke; ein wiederholter Aufruf
	// und ein Aufruf nach dem Commit bleiben ohne Wirkung. Der Use Case ruft
	// `Rollback` auf einem vom Abbruch gelösten Kontext (ohne Frist des
	// Aufrufers): der Adapter beendet sich bei einem abgelösten Kontext nicht
	// vorzeitig und begrenzt seine Dauer selbst.
	Rollback(ctx context.Context) error
}
