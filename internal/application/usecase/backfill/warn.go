package backfill

import "github.com/pt9912/pg-change-feed/internal/domain/model"

// Die Warnungen eines Runs (`ADR-0113` Festlegung 3): zwei Kennzeichnungen an
// der Run-Zeile, keine Ablehnung, kein Abbruch, keine Statusänderung. Diese
// Datei ist die einzige Stelle der Auswertung und der beiden Werte, gegen
// die sie vergleicht.

// copyDurationToleranceMinutes ist die Toleranz der Kopierdauer eines Runs in
// Minuten: der Startwert, eine Setzung ohne Messung. `tools/bench-backfill.sh`
// liest sie hier.
const copyDurationToleranceMinutes = 10

// copyDurationToleranceNanos ist die Toleranz in Nanosekunden, dem Maß von
// `model.Duration`.
const copyDurationToleranceNanos = int64(copyDurationToleranceMinutes) * 60 * 1_000_000_000

// estimatedRowsGuideline ist die Richtgröße in **geschätzten** Zeilen: eine
// Orientierung, keine Grenze. Sie ist aus einer Messung abgeleitet (Rate mal
// Toleranz, auf eine Stelle abgerundet); Host, Lauf und Herleitung stehen im
// Handbuch, Abschnitt „Grenzwerte“.
const estimatedRowsGuideline int64 = 4_000_000

// warnsEstimatedSize meldet Warnung (1): die geschätzte Zeilenzahl liegt über
// der Richtgröße. Genau auf der Richtgröße warnt sie nicht; eine unbekannte
// Schätzung warnt nicht.
func warnsEstimatedSize(estimate model.RowEstimate) bool {
	rows, known := estimate.Rows()
	return known && rows > estimatedRowsGuideline
}

// warnedEstimatedSize liefert den Run mit der Warnung (1), die seine Schätzung
// verlangt.
func warnedEstimatedSize(run model.BackfillRun) model.BackfillRun {
	run.WarnEstimatedSize = warnsEstimatedSize(run.EstimatedRows)
	return run
}

// warnsCopyDuration meldet Warnung (2): die Kopierdauer `now − started` liegt
// über der Toleranz. Genau auf der Toleranz warnt sie nicht. Sie beruht allein
// auf der Uhr und ist von der Schätzung unabhängig.
func warnsCopyDuration(started, now model.TimePoint) bool {
	return now.Sub(started).Nanos > copyDurationToleranceNanos
}

// warnedCopyDuration liefert den Run mit der Warnung (2), wenn seine
// Kopierdauer bis `now` die Toleranz überschreitet. Eine gesetzte Warnung
// bleibt gesetzt; ein Run ohne Beginn der Kopie (`queued`) trägt keine.
func warnedCopyDuration(run model.BackfillRun, now model.TimePoint) model.BackfillRun {
	if run.WarnDuration || run.StartedAt.Unset() {
		return run
	}
	run.WarnDuration = warnsCopyDuration(run.StartedAt, now)
	return run
}
