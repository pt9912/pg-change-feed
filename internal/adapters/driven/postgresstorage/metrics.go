package postgresstorage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/queries"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/sqlexec"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// PostgresMetricsAdapter ist die Referenzimplementierung des
// `MetricsReadPort` (`ADR-0149`): er liest die Zeilen der Sicht `cdc.metrics`
// unverändert, ohne eigene Abfrage über die Basistabellen. `db` trägt die
// Ausführung über die schmale Naht (`sqlexec`).
type PostgresMetricsAdapter struct {
	db  sqlexec.DB
	log outbound.LogPort
}

// NewMetrics legt den Verbindungspool an; der Aufrufer übergibt
// `cfg.ReaderDSN` (Rolle `cdc_reader`, `SELECT` auf der Sicht). Der Pool
// verbindet erst beim ersten Lesen: ein Empfänger-unabhängiger Export darf den
// Start nicht an einer nicht erreichbaren Datenbank scheitern lassen, ein
// Lesefehler ist ein Zyklus-Fehlschlag (`Read`).
func NewMetrics(ctx context.Context, dsn string, opts ...Option) (*PostgresMetricsAdapter, error) {
	o := newOptions(opts)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, metricsReadFailure(ctx, o.log, err)
	}
	return &PostgresMetricsAdapter{db: pool, log: o.log}, nil
}

// metricsReadFailure trägt die Übersetzungsverantwortung dieses Adapters:
// Treiber-Fehler gehen an dieser Grenze über den Port-Sentinel
// (`outbound.ErrMetricsRead`) zurück, die Ursache bleibt über die zweite
// Wrappung lesbar.
func metricsReadFailure(ctx context.Context, log outbound.LogPort, cause error) error {
	log.Debug(ctx, "metrics: Lesefehler", "error", cause)
	return sqlexec.Classify(outbound.ErrMetricsRead, cause)
}

// Close schließt den Verbindungspool.
func (a *PostgresMetricsAdapter) Close() {
	a.db.Close()
}

var _ outbound.MetricsReadPort = (*PostgresMetricsAdapter)(nil)

// Read liest alle Zeilen der Sicht in einem Aufruf. Der Messzeitpunkt bleibt
// ungesetzt, ihn setzt der Use Case. Ein Wert, den die Sicht nicht als Zahl
// liefert, ist ein Lesefehler.
func (a *PostgresMetricsAdapter) Read(ctx context.Context) ([]outbound.MetricSample, error) {
	rows, err := a.db.Query(ctx, queries.SelectMetrics)
	if err != nil {
		return nil, metricsReadFailure(ctx, a.log, err)
	}
	defer rows.Close()
	samples := make([]outbound.MetricSample, 0, 16)
	for rows.Next() {
		var name string
		var label *string
		var text string
		if err := rows.Scan(&name, &label, &text); err != nil {
			return nil, metricsReadFailure(ctx, a.log, err)
		}
		value, err := parseMetricValue(text)
		if err != nil {
			return nil, metricsReadFailure(ctx, a.log, fmt.Errorf("Wert von %s: %w", name, err))
		}
		sample := outbound.MetricSample{Name: name, Value: value}
		if label != nil {
			sample.Label = *label
		}
		samples = append(samples, sample)
	}
	if err := rows.Err(); err != nil {
		return nil, metricsReadFailure(ctx, a.log, err)
	}
	return samples, nil
}

// parseMetricValue liest den Text eines `numeric`: eine ganze Zahl im Bereich
// von `int64` trägt `Int` exakt, jeder andere Wert nur `Float`.
func parseMetricValue(text string) (outbound.MetricValue, error) {
	if i, err := strconv.ParseInt(text, 10, 64); err == nil {
		return outbound.MetricValue{Float: float64(i), Int: i, IsInt: true}, nil
	}
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return outbound.MetricValue{}, err
	}
	return outbound.MetricValue{Float: f}, nil
}
