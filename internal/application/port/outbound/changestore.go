package outbound

import (
	"context"
	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ChangeRecord trägt einen gelesenen Change mit der Commit-Position seiner
// Quelltransaktion (`LH-FA-REA-004`): die Ordnungs- und Bereichs-Größe
// des Lesens liegt auf der Transaktion, der Change
// trägt sie nicht — der Datensatz am Port bündelt beide.
// CommittedAt trägt den realen Quell-Commit-Zeitpunkt seiner Transaktion;
// die zeitbasierte Retention liest ihr
// Alter gegen diesen Zeitpunkt.
type ChangeRecord struct {
	Position    model.SourcePosition
	Change      model.Change
	CommittedAt model.TimePoint
}

// RetentionCandidate trägt die drei Größen, die die Bereinigung von einem
// Change liest: seine Kennung, die Commit-Position seiner Quelltransaktion
// und deren Commit-Zeitpunkt (`ADR-0124`). Row Images gehören nicht dazu.
type RetentionCandidate struct {
	ChangeID    model.ChangeID
	Position    model.SourcePosition
	CommittedAt model.TimePoint
}

// ChangeQuery trägt die Lese-Eingabe am `ChangeStorePort`: Bereich
// (`LH-FA-REA-001`), Startposition, Limit
// und den Tabellenfilter. Nil- und
// leere Felder grenzen nicht ein: ohne Start liest der Aufruf ab dem ersten
// Change der Quelle, ohne End bis zum letzten, ohne Limit unbegrenzt, ohne
// Schema-/Tabellenfilter über alle Tabellen der Quelle. Der Start ist
// inklusive, das Ende exklusiv (Happy Path: Bereich
// `[p1, p2)` trägt den Change an `p1`, nicht den an `p2`).
//
// Die Filterachse ist **eine** Form: `Schema` und `Table` als Klartext,
// je optional und **unabhängig** — dieselbe Filter-Grammatik wie die
// `WHERE`-Klausel des View-Zugriffs auf `cdc.changes`. Klartext ist die
// Draht-Form: die opake `SourceTableID` ist kein Bezeichner, über den ein
// Consumer dieselbe Tabelle stabil adressiert.
type ChangeQuery struct {
	Source model.SourceID
	Start  *model.SourcePosition
	End    *model.SourcePosition
	Schema string
	Table  string
	Limit  *int
}

// Fehler der Lese-Eingabe: sie tragen die Grenzen des Port-Vertrags, nicht
// Domänen-Invarianten — die Positionen selbst bleiben Domänenwerte
// (`ADR-0029`); diese Sentinels gelten am Port-Kontrakt.
var (
	// ErrRangeInverted: die Endposition liegt vor der Startposition
	// (`LH-FA-REA-001` Negative: expliziter Fehlerpfad, kein unbestimmter
	// leerer Bereich).
	ErrRangeInverted = stderrors.New("Endposition liegt vor der Startposition")

	// ErrNonPositiveLimit: ein gesetztes Limit ist mindestens 1
	// (`LH-FA-REA-003` Negative: expliziter Fehlerpfad statt Normierung;
	// unbegrenzt liest der Aufruf über ein nicht gesetztes Limit).
	ErrNonPositiveLimit = stderrors.New("Lese-Limit ist kleiner als 1")

	// ErrStorage trägt die Fehlerklasse `storage` des Store (`ADR-0023`):
	// ein Persistenzfehler im ChangeStore endet ohne
	// Source-ACK; Application und Betrieb lesen die
	// Klasse über errors.Is und kennen keinen Treibertyp. Der Adapter
	// wickelt seine Treiber-Fehler in dieses Sentinel; die technische
	// Ursache bleibt über `errors.Is` lesbar.
	ErrStorage = stderrors.New("Fehlerklasse storage: Persistenzfehler im ChangeStore")
)

// Validate prüft die Bereichs- und Limit-Grenzen der Abfrage. Positionen
// einer anderen Quelle verwirft sie über die Domänen-Ordnung — die
// fachliche Ordnung gilt innerhalb einer Quelle (`ADR-0005`).
func (q ChangeQuery) Validate() error {
	if q.Limit != nil && *q.Limit < 1 {
		return ErrNonPositiveLimit
	}
	if q.Start != nil && q.End != nil && q.End.Before(*q.Start) {
		return ErrRangeInverted
	}
	if q.Source == "" {
		return domainerrors.ErrEmptyIdentifier
	}
	if q.Start != nil && q.Start.SourceID != q.Source {
		return domainerrors.ErrSourceMismatch
	}
	if q.End != nil && q.End.SourceID != q.Source {
		return domainerrors.ErrSourceMismatch
	}
	return nil
}

// ChangeStorePort trägt die Persistenz- und Lesefähigkeit des Change
// Store (`ARC-009`): er persistiert committed CDC-Transaktionen dauerhaft,
// stellt persistierte Changes bereit und ist die einzige Persistenz-Grenze
// des Capture-Pfads — ein Fähigkeits-Port, der Lesen und
// Schreiben gleichermaßen abdeckt. Die erste Produktionsimplementierung
// ist der `PostgresChangeStoreAdapter`; Tests setzen Fake-Ports ein.
//
// Der Port bestätigt keine Quellpositionen. Bestätigt wird eine Position
// erst, wenn ihre abhängigen Changes dauerhaft gespeichert sind — ACK nur
// über nach Persistenz gefragte Positionen.
// Die Quell-Bestätigung ist die Wirkung des
// `ReplicationAckPort`, keine Speicherwirkung dieses Ports.
//
// Der Port trägt die Idempotenz-Pflicht des Persistierens:
// PersistTransaction ist deduplizierbar — dieselbe Transaktion darf
// erneut persistiert werden, wenn ein Crash zwischen Persistenz und ACK
// sie wiederholt; die Deduplizierungsbasis trägt die interne
// Transaktions-ID (`cdc.transaction`).
//
// Der Port trägt die Fehlerklassen-Grenze: Treiber-Fehler
// gehen an der Adapter-Grenze in die Klassen des Pflichtenhefts; dieser
// Port führt die Klasse `storage` als `ErrStorage`. Application und
// Betrieb klassifizieren über `errors.Is`, nicht über den Treiber.
type ChangeStorePort interface {
	// PersistTransaction persistiert eine committed Quelltransaktion
	// dauerhaft — mit interner ID, Commit-Position und Changes in
	// Sequenz-Reihenfolge (`SPEC-001`, `cdc.transaction`). Die Rückkehr
	// ohne Fehler meldet den abgeschlossenen Store-Commit
	// (Schritt COMMIT Store).
	PersistTransaction(ctx context.Context, transaction *model.ChangeTransaction) error

	// ReadChanges liest persistierte Changes über die Schema-Tabellen und
	// ordnet sie deterministisch nach (Commit-Position der Quelltransaktion,
	// Transaktions-ID, Sequenz innerhalb der Transaktion)
	// (`LH-FA-REA-004`). Lesen verändert keine gespeicherte Position;
	// gelesene Changes bleiben innerhalb der Aufbewahrung
	// erneut lesbar.
	ReadChanges(ctx context.Context, query ChangeQuery) ([]ChangeRecord, error)

	// ReadRetentionCandidates liest eine Seite von Bereinigungs-Kandidaten
	// einer Quelle (`ADR-0124`): höchstens `limit` Changes mit einer Kennung
	// größer `after`, aufsteigend in der Ordnung des Schlüssels `change_id`
	// — nicht in der fachlichen Ordnung von `ReadChanges`. `after` leer
	// beginnt am Anfang; eine leere Seite ist das Ende, eine kürzere als
	// `limit` nicht. Ein `limit` kleiner 1 endet als `ErrNonPositiveLimit`,
	// eine leere Quelle als `ErrEmptyIdentifier`. Die Seite trägt kein Row
	// Image; die Freigabe je Kandidat bleibt beim aufrufenden Use Case.
	ReadRetentionCandidates(ctx context.Context, source model.SourceID, after model.ChangeID, limit int) ([]RetentionCandidate, error)

	// DeleteChanges entfernt physisch genau die übergebenen Changes
	// (`LH-FA-RET-002`): die Freigabe je Change trägt
	// `RetentionPolicy.AllowsDeletion` im aufrufenden Use Case, dieser Port
	// führt nur die bereits freigegebene Menge aus — keine eigene
	// Freigabe-Entscheidung. Eine leere Menge ist ein gültiger Aufruf ohne
	// Wirkung; eine bereits entfernte oder nie vorhandene Kennung bleibt
	// ohne Wirkung (Idempotenz). Treiber-Fehler gehen in die Klasse
	// `storage` (`ErrStorage`).
	DeleteChanges(ctx context.Context, changeIDs []model.ChangeID) error
}
