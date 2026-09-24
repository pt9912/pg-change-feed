package outbound

import (
	"context"
	stderrors "errors"
)

// Fehlerklassen des Snapshot-Lesers (`SPEC-008`, `ADR-0023`): jede
// Störung endet als sichtbarer Fehler mit einer dieser Klassen. Der
// Aufrufer klassifiziert über `errors.Is`; die technische Ursache bleibt
// über die zweite Wrappung lesbar.
var (
	// ErrSnapshotPermission: der Capture-Login darf die Quelltabelle nicht
	// lesen oder den temporären Slot nicht anlegen (`SELECT` und
	// `REPLICATION` sind Betriebs-Vorbedingungen, `ADR-0047`).
	ErrSnapshotPermission = stderrors.New("Fehlerklasse permission: fehlende Berechtigung für den Tabellen-Snapshot")

	// ErrSnapshotConfiguration: die Anfrage oder die Quelle steht im
	// falschen Stand — Tabelle nicht vorhanden, ungültige Kennung, keine
	// Reserve in `max_replication_slots`/`max_wal_senders`.
	ErrSnapshotConfiguration = stderrors.New("Fehlerklasse configuration: Tabellen-Snapshot im falschen Stand der Konfiguration")

	// ErrSnapshotTransient: die Quelle war für den Snapshot vorübergehend
	// nicht verfügbar — das Zeitlimit der Slot-Anlage ist abgelaufen, die
	// Verbindung ist abgebrochen oder der Kontext endete.
	ErrSnapshotTransient = stderrors.New("Fehlerklasse transient: Quelle für den Tabellen-Snapshot vorübergehend nicht verfügbar")

	// ErrSnapshotReplication: die Störung liegt an Replication-Verbindung
	// oder Slot-Anlage, ohne dass eine spezifischere Klasse zutrifft.
	ErrSnapshotReplication = stderrors.New("Fehlerklasse replication: Slot-/Replication-Störung beim Tabellen-Snapshot")

	// ErrSnapshotStorage: das Lesen der Zeilen im Snapshot ist gescheitert.
	ErrSnapshotStorage = stderrors.New("Fehlerklasse storage: Lesefehler im Tabellen-Snapshot")
)

// TableSnapshotPort trägt die Fähigkeit, den Bestand einer Tabelle als
// konsistenten Snapshot zu lesen (`ARC-004`, Fähigkeits-Port je
// `ADR-0034`; `LH-FA-CAP-009`). Er kennt weder Run-Zustand noch Store: die
// Abbildung der Snapshot-Position auf eine `SourcePosition` (`ADR-0005`),
// die Bild-Konstruktion (`model.BuildRowImage`) und das Schreiben liegen
// beim Aufrufer.
type TableSnapshotPort interface {
	// OpenSnapshot öffnet den Snapshot einer Tabelle für den Run
	// `runID`: Snapshot und Position sind gepaart — jeder Commit mit
	// Position ≤ `Offset()` steckt im Bestand, jeder spätere nicht.
	// Der Aufrufer schließt den Snapshot, auch nach einem Fehler beim
	// Lesen.
	OpenSnapshot(ctx context.Context, runID, schema, table string) (TableSnapshot, error)

	// EstimatedRows liest die vom Katalog **geschätzte** Zeilenzahl der
	// Tabelle. `known` ist falsch, solange der Katalog keine Schätzung
	// führt (nie analysierte Tabelle): die Schätzung ist dann unbekannt,
	// nie 0 (`ADR-0113` Festlegung 3, Punkt 5). Die Zahl ist eine
	// Orientierung, keine Grenze.
	EstimatedRows(ctx context.Context, schema, table string) (rows int64, known bool, err error)
}

// TableSnapshot ist ein geöffneter Tabellen-Snapshot: eine
// Lese-Transaktion im Stand des Slot-Punkts. Ein Snapshot wird von einer
// Goroutine gelesen.
type TableSnapshot interface {
	// Offset trägt die Position `X` des Snapshots als Rohwert der LSN
	// (`ADR-0005`: der Offset einer `SourcePosition`).
	Offset() uint64

	// Columns trägt die gelesenen Spaltennamen in Tabellenordnung — ohne
	// gelöschte und ohne generierte Spalten, denn der WAL-Pfad sendet
	// sie nicht.
	Columns() []string

	// NextBlock liest die nächsten höchstens `B` Zeilen. Zeile `i` trägt
	// je Spalte den Text-Stand der Quelle in der Ordnung von `Columns()`:
	// das Ergebnis der Ausgabefunktion des Spaltentyps unter den
	// Sitzungs-GUC, ohne Cast und ohne Funktion auf dem Spaltenwert — dieselbe
	// Erzeugung wie der Text, den `pgoutput` je Spalte sendet (`ADR-0115`);
	// `nil` ist NULL. Ein leerer Block meldet das Ende des Bestands —
	// auch beim ersten Aufruf einer leeren Tabelle. Die Zeilenzahl eines
	// Aufrufs ist durch `B` begrenzt (`LH-FA-CAP-006.a`); der Speicherbedarf
	// ist `B` mal die Zeilenbreite, weil die Grenze Zeilen zählt, nicht Bytes,
	// und der Block bis zur Rückgabe als Treiberbytes und als Zeichenketten
	// vorliegt (aus dem Treiber-Quelltext abgeleitet, nicht gemessen).
	NextBlock(ctx context.Context) ([][]*string, error)

	// Close beendet die Lese-Transaktion und die Verbindung; ein
	// wiederholter Aufruf bleibt ohne Wirkung. Der Aufrufer ruft `Close` auf
	// einem vom Abbruch gelösten Kontext (ohne Frist des Aufrufers): der
	// Adapter beendet sich bei einem abgelösten Kontext nicht vorzeitig und
	// begrenzt seine Dauer selbst.
	Close(ctx context.Context) error
}
