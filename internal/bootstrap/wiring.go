// Package bootstrap ist die Composition Root (`ADR-0026`): er kennt die
// konkreten Adapter und verdrahtet die Pipeline an genau einer Stelle —
// ChangeStore-Driven-Adapter, Replication-Stream-Driving-Adapter, Capture
// Service und Replication-ACK-Driven-Adapter. Die Abhängigkeitsregel (§2
// der Architektur-Sicht) bleibt hier lokal einhaltbar; `main` referenziert
// keinen Adapter-Konstruktor.
//
// Die Verdrahtung liest ihre Vorbedingungen als Minimal-Form aus der
// Umgebung: DSN, Quelle, Publication, Slot-Name und die
// Tabellen-Aktivierungen. Eine vollständige Konfigurationsschicht mit
// Format-Wahl ist nicht Teil dieses Verdrahtungsstands.
package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresack"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Umgebungs-Namen der Verdrahtungs-Vorbedingungen; sie tragen dieselben
// Werte wie der Container-Vertrag in compose.yaml.
const (
	envDSN         = "CDC_SOURCE_DSN"
	envSource      = "CDC_SOURCE_ID"
	envPublication = "CDC_PUBLICATION"
	envSlot        = "CDC_SLOT"
	envTables      = "CDC_TABLES"
)

// ErrConfiguration trägt die Fehlerklasse `configuration` der Verdrahtung
// (`SPEC-008`, `ADR-0023`): eine fehlende oder falsch gesetzte
// Vorbedingung endet ohne Start und ohne Fortsetzung im falschen Stand;
// der Prozess-Aufrufer meldet sie als Ausgang.
var ErrConfiguration = errors.New("Fehlerklasse configuration: Verdrahtung ohne vollständige Vorbedingung")

// Config trägt die Verdrahtungs-Eingabe: die Verbindung zur Instanz, die
// Quelle und die CDC-Speicherrollen gleichermaßen trägt (MVP-Schnitt,
// Abschnitt 1 Lastenheft), die Quelle, die Verwaltungs-Namen Publication
// und Slot (`LH-FA-CFG-001.a`) und die aktivierten Tabellen mit ihren
// Port-Kennungen.
type Config struct {
	DSN         string
	Source      model.SourceID
	Publication string
	Slot        string
	Tables      map[string]mapper.TableBinding
}

// ConfigFromEnv liest die Verdrahtungs-Vorbedingungen über die
// übergebene Umgebungs-Lese-Funktion; eine fehlende Vorbedingung endet
// über die Klasse `configuration` (ErrConfiguration).
func ConfigFromEnv(getenv func(string) string) (Config, error) {
	cfg := Config{
		DSN:         getenv(envDSN),
		Source:      model.SourceID(getenv(envSource)),
		Publication: getenv(envPublication),
		Slot:        getenv(envSlot),
	}
	if cfg.DSN == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envDSN)
	}
	if cfg.Source == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envSource)
	}
	if cfg.Publication == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envPublication)
	}
	if cfg.Slot == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envSlot)
	}
	tables, err := parseTables(getenv(envTables))
	if err != nil {
		return Config{}, err
	}
	cfg.Tables = tables
	return cfg, nil
}

// parseTables liest die Tabellen-Aktivierungen in der Form
// `schema.table=tabelle-id:schema-version-id`, Komma-getrennt; die
// Bindungs-Kennungen tragen die CDC-Referenztabellen (`LH-FA-CFG-001`).
// Ohne Aktivierung startet die Verdrahtung nicht — ein Feed ohne aktivierte
// Tabelle ist kein Stand zur Fortsetzung.
func parseTables(raw string) (map[string]mapper.TableBinding, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("%w: %s trägt keine Tabellen-Aktivierung", ErrConfiguration, envTables)
	}
	tables := map[string]mapper.TableBinding{}
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		qualified, binding, found := strings.Cut(entry, "=")
		if !found || qualified == "" {
			return nil, fmt.Errorf("%w: Aktivierung %q trägt nicht die Form schema.table=tabelle-id:schema-version-id", ErrConfiguration, entry)
		}
		tableID, versionID, found := strings.Cut(binding, ":")
		if !found || tableID == "" || versionID == "" {
			return nil, fmt.Errorf("%w: Bindung %q trägt nicht die Form tabelle-id:schema-version-id", ErrConfiguration, binding)
		}
		tables[qualified] = mapper.TableBinding{
			TableID:       model.SourceTableID(tableID),
			SchemaVersion: model.SchemaVersionID(versionID),
		}
	}
	return tables, nil
}

// Run verdrahtet die Pipeline (`ADR-0026`) und trägt den Stream-Lauf bis
// zum Kontext-Ende: der Stream baut die Replication-Verbindung, der
// ACK-Adapter bestätigt über dieselbe Verbindung (`ADR-0007`, Option C)
// und der Capture Service orchestriert Persist-before-ACK
// (`LH-QA-REL-001.a`). Die Rückkehr ohne Fehler meldet das reguläre
// Lauf-Ende; ein Fehler wird durchgereicht, nicht still fortgesetzt
// (`SPEC-008`).
func Run(ctx context.Context, cfg Config) error {
	store, err := postgresstorage.New(ctx, cfg.DSN)
	if err != nil {
		return err
	}
	defer store.Close()

	stream, err := receive.NewStream(ctx, receive.Config{
		DSN:         cfg.DSN,
		Source:      cfg.Source,
		Publication: cfg.Publication,
		Slot:        cfg.Slot,
		Tables:      cfg.Tables,
	})
	if err != nil {
		return err
	}
	ack, err := postgresack.New(stream.Conn())
	if err != nil {
		return err
	}
	if err := stream.BindCapture(capture.NewCaptureService(store, ack)); err != nil {
		return err
	}
	return stream.Run(ctx)
}