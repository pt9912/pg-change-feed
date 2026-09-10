package outbound

import (
	"context"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// ActivationRemoval trägt den Ausgang des Bindungs-Zeilen-Entzugs einer
// Deaktivierung (`LH-FA-CFG-002`).
type ActivationRemoval string

const (
	// ActivationRemoved meldet die entfernte Bindungs-Zeile: keine
	// persistierten Changes referenzieren die Tabelle mehr.
	ActivationRemoved ActivationRemoval = "entfernt"

	// ActivationRetained meldet die belassene Bindungs-Zeile: die
	// persistierten Changes tragen sie als Herkunft (`LH-FA-CFG-002`
	// Out-of-Scope: ihr Verhalten folgt der Retention, nicht der
	// Deaktivierung); der Entzug löscht keinen Change und keine
	// Schema-Version.
	ActivationRetained ActivationRemoval = "belassen"

	// ActivationAbsent meldet, dass die Tabelle nicht aktiviert war
	// (`LH-FA-CFG-002` Boundary: idempotentes Verhalten).
	ActivationAbsent ActivationRemoval = "nicht aktiviert"
)

// TableActivationPort trägt die Verwaltungs-Fähigkeit der aktivierten
// Tabellen an der Quelle (`ARC-004`, Fähigkeits-Port je `ADR-0034`): die
// Tabellen-Aktivierung (`LH-FA-CFG-001`, `LH-FA-CFG-002`, `ADR-0028`)
// schreibt die Bindungs-Zeilen der CDC-Referenztabellen (`SPEC-001`,
// `cdc.source_table`, `cdc.schema_version`) und trägt die Publication
// der Quelle (`LH-FA-CFG-001.a`, Schritt 2 — Publication-Verwaltung; den
// Slot trägt der Stream-Adapter, `ADR-0006`).
//
// Die Quelle selbst (Zeile `cdc.source`) registriert der Aufrufer
// vor der ersten Aktivierung — die Fremdschlüssel der DDL setzen sie
// voraus; die Metadaten-Registrierung ist keine Wirkung der
// Tabellen-Aktivierung (Kopplung: `SPEC-001`).
//
// Alle Operationen sind idempotent (`LH-FA-CFG-001`/`002` Boundary):
// ein erneuter Aufruf ändert keinen Stand und meldet keinen Fehler.
type TableActivationPort interface {
	// TableExists prüft die physische Tabelle an der Quelle; die
	// Negative-Pfade der Aktivierung, Deaktivierung und Status-Abfrage
	// (`LH-FA-CFG-001`/`002`/`003`) enden über sie sichtbar statt still.
	TableExists(ctx context.Context, schema, table string) (bool, error)

	// Registered liest die Bindungs-Zeile einer Tabelle; die Abwesenheit
	// ist der Zustand „nicht aktiviert" (`LH-FA-CFG-003` Boundary).
	Registered(ctx context.Context, source model.SourceID, schema, table string) (model.SourceTable, bool, error)

	// Register trägt die Bindungs- und Schema-Version-Zeile einer
	// Aktivierung ein; die Rückkehr meldet, ob der Aufruf neu
	// aktiviert hat. Eine bereits aktivierte Tabelle bleibt unverändert
	// (`LH-FA-CFG-001` Boundary: keine doppelte Erfassung).
	Register(ctx context.Context, table model.SourceTable, version model.SchemaVersion) (bool, error)

	// Unregister entzieht die Bindungs-Zeile einer Tabelle; der Ausgang
	// trägt die drei Zustände der Deaktivierung. Bei Change-Bestand
	// bleibt die Zeile als Herkunft der persistierten Changes bestehen
	// (ActivationRetained) — der Entzug löscht keine Changes und keine
	// Schema-Versionen.
	Unregister(ctx context.Context, table model.SourceTable) (ActivationRemoval, error)

	// List liest die Bindungs-Zeilen einer Quelle — die Liste der
	// aktivierten Tabellen (`LH-FA-CFG-004`); ohne Aktivierung liest sie
	// leer.
	List(ctx context.Context, source model.SourceID) ([]model.SourceTable, error)

	// Publish trägt die Publication der Quelle (`LH-FA-CFG-001.a`): er
	// legt die Publication an, wenn sie fehlt, nimmt die Tabelle
	// auf (INSERT, UPDATE, DELETE — TRUNCATE trägt die Publication
	// nicht, `LH-FA-CFG-001.a` Schritt 3) und bleibt ohne Wirkung, wenn
	// die Tabelle bereits Mitglied ist.
	Publish(ctx context.Context, publication, schema, table string) error

	// Unpublish entzieht die Tabelle der Publication und bleibt ohne
	// Wirkung, wenn die Publication fehlt oder die Tabelle kein Mitglied
	// ist. Der Entzug trägt den Stopp der Erfassung (`LH-FA-CFG-002`
	// Happy Path); seine Wirkung auf einen laufenden Stream-Container
	// liegt am Walsender der Quelle.
	Unpublish(ctx context.Context, publication, schema, table string) error
}
