package sqlexec

// Statement trägt eine Anweisung an die Naht samt der Übersetzungs-
// Verantwortung ihres Aufrufers: `SQL`/`Args` gehen an den Executor, `Fail`
// klassifiziert einen Treiber-Fehler und protokolliert ihn über den
// `LogPort` des Adapters (`LH-QA-OPS-004`, `ADR-0024`) — der eine
// Übersetzungspunkt, den jeder Adapter führt (`storageFailure` und
// Geschwister). Ein Domänen-Fehler der Übersetzung läuft **nicht** durch
// `fail`: er trägt seine Klasse schon selbst und wird unverändert
// zurückgegeben.
type Statement struct {
	SQL  string
	Args []any
	Fail func(error) error
}

// fail reicht einen Treiber-Fehler an die Übersetzung des Aufrufers; ohne
// gesetztes `Fail` bleibt die Ursache unverändert — der Aufrufer hat dann
// keine Klasse zu vergeben.
func (s Statement) fail(cause error) error {
	if s.Fail == nil {
		return cause
	}
	return s.Fail(cause)
}
