// Package sqlexec trägt die schmale Ausführungs-Naht der PostgreSQL-Adapter
// (`ADR-0071` Punkt 5): die Adapter hängen an `Executor` statt am konkreten
// `*pgxpool.Pool`, und die von der Naht getragene Zeilen-Übersetzung samt
// Fehlerklassifikation liegt hier — netzlos prüfbar über einen Träger, der
// die Schnittstelle erfüllt, statt über eine lebende PostgreSQL-Instanz.
// Der reale Pool erfüllt sie ohne Vermittler: die Zusicherung unten ist der
// Kompilier-Beleg dafür.
package sqlexec

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Rows trägt die minimale Lese-Fläche einer Ergebnis-Menge: genau die vier
// Aufrufe, die die Übersetzungen dieses Pakets brauchen. `pgx.Rows` erfüllt
// sie strukturell — ein Träger dieser Naht braucht nicht mehr als diese vier.
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close()
}

// Executor trägt die Ausführungs-Fläche der Adapter: eine Ergebnis-Menge
// absetzen, eine einzelne Zeile lesen, eine Anweisung ausführen. Die
// Methodensignaturen sind die des realen Pools — `*pgxpool.Pool` erfüllt sie
// direkt, `pgx.Tx` ebenfalls, damit ein Transaktions-Träger an derselben
// Naht einsetzbar bleibt.
type Executor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// DB trägt die Fläche, die ein pool-tragender Adapter von seinem Träger
// braucht: die Ausführung (`Executor`), den Transaktions-Einstieg und das
// Lebensende des Pools. `Close` gehört dazu, weil der Adapter den Pool
// besitzt, den er gebaut hat.
type DB interface {
	Executor
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

// Der reale Pool erfüllt die Naht. Die Zusicherung ist der Beleg für „statt
// `*pgxpool.Pool`": ändert eine Pool-Methode ihre Signatur, bricht dieser
// Bau — nicht erst der Aufruf einer Adapter-Methode.
var _ DB = (*pgxpool.Pool)(nil)
