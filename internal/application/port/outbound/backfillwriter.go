package outbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrBackfillBlockInvalid: ein Block, der Commit eines Runs oder eine
// Operation der Schreibtransaktion verletzt den Vertrag des Schreibers — ein
// Block ohne Changes, mit falscher Quelle, Position, Transaktions- oder
// Change-Kennung, ein Change außerhalb von `INSERT` der Herkunft `backfill`
// oder mit `old_data`, ein Commit, dessen Zähler nicht die Zahl der
// angehängten Changes trägt. Der Fehler schreibt nichts; die Transaktion
// bleibt für `Rollback` offen.
var ErrBackfillBlockInvalid = stderrors.New("Block oder Commit des Backfill-Runs verletzt den Vertrag des Schreibers")

// BackfillWriterPort trägt die Fähigkeit, die Blöcke eines Runs in **einer**
// Store-Transaktion zu schreiben (`ARC-004`, Fähigkeits-Port je `ADR-0034`;
// `ADR-0111` Teilfrage 4). Der bestehende `ChangeStorePort.PersistTransaction`
// hält eine ganze Transaktion im Speicher und trägt einen Bestand nicht;
// dieser Port nimmt die Blöcke einzeln entgegen, damit der Aufrufer höchstens
// einen Block hält (`LH-FA-CAP-006.a`).
type BackfillWriterPort interface {
	// Begin öffnet die Schreibtransaktion des Runs, der `running` ist und
	// die Position `X` trägt (sonst `ErrBackfillRunInvalid`). Bis zum Commit
	// ist nichts für Leser sichtbar.
	Begin(ctx context.Context, run model.BackfillRun) (BackfillTransaction, error)
}

// BackfillTransaction ist die eine offene Schreibtransaktion eines Runs.
// Sie wird von einer Goroutine benutzt.
type BackfillTransaction interface {
	// AppendBlock hängt einen Block an: eine committed, synthetische
	// Transaktion (`model.ChangeTransaction`) auf der Position `X` des Runs.
	// Die Blöcke tragen die Kennungen `model.BackfillTransactionID` in der
	// Reihenfolge 1, 2, 3, …, ihre Changes `model.ChangeIDFor`, die Herkunft
	// `backfill`, die Operation `INSERT` und kein `old_data`; jede Abweichung
	// endet als `ErrBackfillBlockInvalid` ohne Schreibvorgang. Der Port hält
	// den Block nach der Rückkehr nicht.
	AppendBlock(ctx context.Context, block *model.ChangeTransaction) error

	// Commit schreibt die Run-Zeile im Endzustand `completed` und
	// committet sie zusammen mit allen angehängten Blöcken — ein Commit,
	// eine Einheit. `run.RowsCopied` trägt die Zahl der angehängten Changes
	// (sonst `ErrBackfillBlockInvalid`). Persistenzfehler tragen
	// `ErrBackfillStorage`.
	Commit(ctx context.Context, run model.BackfillRun) error

	// Rollback verwirft alle angehängten Blöcke; ein wiederholter Aufruf
	// und ein Aufruf nach dem Commit bleiben ohne Wirkung. Der Use Case ruft
	// `Rollback` auf einem vom Abbruch gelösten Kontext (ohne Frist des
	// Aufrufers): der Adapter beendet sich bei einem abgelösten Kontext nicht
	// vorzeitig und begrenzt seine Dauer selbst.
	Rollback(ctx context.Context) error
}
