package postgresstorage

import (
	"context"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/sqlexec"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// Die drei Backfill-Adapter (`BackfillAdmissionAdapter`, `BackfillRunAdapter`,
// `BackfillWriterAdapter`) teilen zwei Frist-Startwerte und die
// Fehlerübersetzung; jeder Adapter hält einen eigenen Verbindungspool
// (`ADR-0113` Festlegung 1: die Annahme läuft über `CDC_ADMIN_DSN`, Run-Zustand
// und Schreiber über `CDC_CAPTURE_DSN`), damit der Fortschritt des Runs auf
// einer anderen Verbindung als die offene Schreibtransaktion läuft.
//
// Jede Operation ist adapterseitig zeitbegrenzt (`outbound.BackfillRunPort`,
// `outbound.BackfillTransaction`): der Use Case ruft `Finish`,
// `InterruptRunning` und `Rollback` auf einem vom Abbruch gelösten Kontext
// ohne Frist, die Frist setzt dieser Adapter selbst. Die Werte sind
// Startwerte ohne Messung.
const (
	// backfillStateTimeout begrenzt jede Operation auf Run-Zeilen und die
	// Verwaltung der Schreibtransaktion (Öffnen, Rollback).
	backfillStateTimeout = 30 * time.Second
	// backfillBlockTimeout begrenzt das Anhängen eines Blocks und den
	// Commit der Schreibtransaktion; beide tragen mehr Zeilen als eine
	// Zustands-Operation.
	backfillBlockTimeout = 5 * time.Minute
)

// backfillStorageFailure trägt die Übersetzungsverantwortung der drei
// Backfill-Adapter (`ADR-0023`, `SPEC-008`): Treiber-Fehler gehen an dieser
// Grenze in die Klasse `storage` über `outbound.ErrBackfillStorage`, die
// technische Ursache bleibt über die zweite Wrappung lesbar.
func backfillStorageFailure(ctx context.Context, log outbound.LogPort, cause error) error {
	log.Error(ctx, "backfill: Datenbankfehler", "error", cause)
	return sqlexec.Classify(outbound.ErrBackfillStorage, cause)
}

// boundedContext leitet den Kontext einer Operation mit ihrer Frist ab. Die
// Frist gilt auch für einen Kontext ohne eigene Frist, den der Aufrufer vom
// Abbruch gelöst hat (`context.WithoutCancel`).
func boundedContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}
