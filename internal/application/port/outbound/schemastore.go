package outbound

import (
	"context"
	stderrors "errors"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ErrSchemaStoreStorage trägt die Fehlerklasse `storage` dieses Ports
// (`SPEC-008`, `ADR-0023`) — derselbe Aufbau wie `ErrConsumerStateStorage`
// (`consumerstate.go`): ein Persistenzfehler am Schema Store endet
// sichtbar, Application und Betrieb klassifizieren über `errors.Is` und
// kennen keinen Treibertyp.
var ErrSchemaStoreStorage = stderrors.New("Fehlerklasse storage: Persistenzfehler im Schema Store")

// ErrSchemaVersionUnknown trägt die Abwesenheit einer TableSchema-Zeile zu
// einer `SchemaVersionID` (`ADR-0015` Folgepflicht, `LH-FA-SCH-004`
// Negative-Fall): eine referenzierte, aber nicht registrierte Version
// endet sichtbar — keine stille Fehlinterpretation.
var ErrSchemaVersionUnknown = stderrors.New("Schema-Version trägt keine TableSchema-Registrierung")

// ErrSchemaVersionMismatch trägt die Aufrufvertrags-Grenze von
// `RegisterVersion`: `model.SchemaVersion` und `model.TableSchema` müssen
// dieselbe `SchemaVersionID` tragen — ein widersprüchlicher Aufruf endet
// vor dem ersten SQL-Aufruf sichtbar, kein Schreiben im falschen Stand.
var ErrSchemaVersionMismatch = stderrors.New("SchemaVersion und TableSchema tragen unterschiedliche Kennungen")

// SchemaStorePort trägt die Persistenz-Fähigkeit der historisch stabilen
// Schema-Interpretation (`ARC-004`, `ADR-0015` Folgepflicht — Architektur-
// Sicht §4, Sequenzdiagramm `LH-FA-CFG-001.a`, `SchemaStorePort (ARC-004)`):
// jeder Change referenziert eine `model.SchemaVersion` (`SPEC-004`), deren
// zugehörige Spaltenform (`model.TableSchema`) dieser Port hält — ohne sie
// bleiben `LH-FA-SCH-004` (inkompatible Typänderungen erkennbar melden)
// und `LH-FA-SCH-005` (Changes einer Schema-Version zuordenbar,
// unterscheidbar) unerfüllbar.
//
// Dieser Port trägt ausschließlich die Persistenz-Fähigkeit: weder die
// dynamische Re-Versionierung im laufenden Erfassungspfad
// (`Assembler.Consume`) noch die Typ-Kompatibilitätsprüfung mit der
// Fehlerklasse `schema` gehören zu ihm — er trägt weder Vergleichs- noch
// Kompatibilitätslogik. Die statische Erstaktivierung
// (`TableActivationPort.Register`, Version 1) bleibt unverändert und
// schreibt ihre Schema-Version-Zeile weiterhin über jenen Port, nicht
// über diesen.
type SchemaStorePort interface {
	// CurrentVersion liest die höchste registrierte Schema-Version einer
	// Tabelle; die Abwesenheit (`bool` false) meldet, dass die Tabelle
	// noch keine über diesen Port registrierte Version trägt.
	CurrentVersion(ctx context.Context, table model.SourceTableID) (model.SchemaVersion, bool, error)

	// RegisterVersion trägt eine Schema-Version und ihre Spaltenform in
	// EINEM Store-Commit ein; die Idempotenz trägt die Anwesenheit der
	// Spaltenform selbst — eine `SchemaVersionID`, die bereits eine
	// Spaltenform trägt, bleibt ohne Wirkung (Rückkehr `false`). Eine
	// `SchemaVersionID` ohne bestehende Spaltenform bekommt sie
	// geschrieben und die Rückkehr meldet `true` — unabhängig davon, ob
	// die `SchemaVersion`-Zeile selbst neu ist oder bereits über einen
	// anderen Schreibpfad besteht. `version.ID` und `schema.VersionID`
	// müssen übereinstimmen (`ErrSchemaVersionMismatch`).
	RegisterVersion(ctx context.Context, version model.SchemaVersion, schema model.TableSchema) (bool, error)

	// TableSchema liest die Spaltenform einer Schema-Version — die
	// Grundlage der historisch stabilen Interpretation älterer Changes.
	// Die Abwesenheit endet über `ErrSchemaVersionUnknown` sichtbar.
	TableSchema(ctx context.Context, versionID model.SchemaVersionID) (model.TableSchema, error)
}
