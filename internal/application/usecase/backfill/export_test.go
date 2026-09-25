package backfill

// Prüf-Zugriff der Tests des Pakets `backfill_test` auf die beiden Konstanten
// der Warn-Auswertung (`warn.go`); kompiliert nur im Test.
const (
	EstimatedRowsGuidelineForTest     = estimatedRowsGuideline
	CopyDurationToleranceNanosForTest = copyDurationToleranceNanos
)
