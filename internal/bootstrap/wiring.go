// Package bootstrap ist die Composition Root (`ADR-0026`): er kennt die
// konkreten Adapter und verdrahtet die Pipeline an genau einer Stelle —
// ChangeStore-Driven-Adapter, Aktivierungs-Driven-Adapter mit dem
// EnableTable Use Case (`ADR-0028`), Heartbeat-Driven-Adapter mit dem
// periodischen Timer-Zug (`ADR-0024`), Retention-Use-Case mit
// periodischem Lösch-Takt über den System-Clock-Adapter (`ADR-0014`,
// `ADR-0040`), Replication-Stream-Driving-Adapter, Capture Service und
// Replication-ACK-Driven-Adapter.
// Die Abhängigkeitsregel (§2 der Architektur-Sicht) bleibt hier lokal
// einhaltbar; `main` referenziert keinen Adapter-Konstruktor.
//
// Die Verdrahtung liest ihre Vorbedingungen aus der Umgebung: DSN, Quelle,
// Publication, Slot-Name und die Tabellen-Aktivierungen (`ConfigFromEnv`).
// Eine optionale YAML-Konfigurationsdatei ergänzt das additiv, mit
// Umgebungsvariable-schlägt-Datei-Feld-für-Feld-Precedence und
// env-var-exklusiven Zugangsdaten (`ConfigFromEnvAndFile`, `config_file.go`,
// `ADR-0052`, `ADR-0088`).
package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nats-io/nats.go"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/grpcstream"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/natsnotify"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/natsstream"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresack"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgressnapshot"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/telemetry"
	apigrpc "github.com/pt9912/pg-change-feed/internal/adapters/driving/grpc"
	apihttp "github.com/pt9912/pg-change-feed/internal/adapters/driving/http"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/acknowledge"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/backfill"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/disable"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/enable"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/excludecolumn"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/includecolumn"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/list"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/position"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/readchanges"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/register"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/remove"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/removetransformation"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/retention"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/settransformation"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/status"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Umgebungs-Namen der Verdrahtungs-Vorbedingungen; sie tragen dieselben
// Werte wie der Container-Vertrag in compose.yaml. Die drei DSN-Variablen
// tragen je eine PostgreSQL-Rolle (`ADR-0047`): `envCaptureDSN` bindet an
// `cdc_capture`, `envAdminDSN` an `cdc_admin`, `envReaderDSN` an
// `cdc_reader` — `CDC_SOURCE_DSN` (eine gemeinsame Instanz-DSN für alle
// Aufrufer) entfällt ersatzlos.
const (
	envCaptureDSN  = "CDC_CAPTURE_DSN"
	envAdminDSN    = "CDC_ADMIN_DSN"
	envReaderDSN   = "CDC_READER_DSN"
	envSource      = "CDC_SOURCE_ID"
	envPublication = "CDC_PUBLICATION"
	envSlot        = "CDC_SLOT"
	envTables      = "CDC_TABLES"
	// envLogLevel trägt den Log-Level der strukturierten Ausgabe
	// (`LH-QA-OPS-004`): anders als die fünf Namen oben ist er
	// keine Start-Vorbedingung — er hat einen Default (parseLogLevel) und
	// eine fehlende oder nicht erkannte Eingabe bricht die Verdrahtung
	// nicht ab.
	envLogLevel = "CDC_LOG_LEVEL"
	// envNatsURL trägt die optionale NATS-Server-URL für das
	// Change-Notification-Wecksignal (`ADR-0055`, `LH-FA-SST-007`): anders
	// als die fünf Namen oben ist sie keine Start-Vorbedingung — ungesetzt
	// bleibt das Feature vollständig deaktiviert, kein
	// `ChangeNotificationPort` wird konstruiert und keine NATS-Verbindung
	// aufgebaut (bestehendes Verhalten bit-identisch, `ADR-0055` Punkt 5).
	// Ist sie gesetzt, ist eine erfolgreiche Verbindung dagegen eine
	// explizite Vorbedingung dieses Laufs (`ErrConfiguration` bei
	// Fehlschlag) — anders als `envLogLevel`, dessen Fehlerfall die
	// Erfassung selbst nicht gefährdet.
	envNatsURL = "CDC_NATS_URL"
	// envNatsStreamToken trägt den optionalen NATS-Verbindungs-Token des
	// dritten, vollinhaltstragenden Zustellwegs (`ADR-0100` Teilfrage 4/5,
	// `LH-FA-SST-008`): eine gesetzte `CDC_NATS_STREAM_TOKEN` aktiviert —
	// zusammen mit `envNatsURL` — sowohl den `natsstream.Publisher` als
	// auch eine Token-Client-Option an der bestehenden
	// `nats.Connect`-Aufrufstelle. Ist nur `envNatsURL` gesetzt (heutiger
	// Zustand jeder bestehenden Installation), bleibt das Wecksignal
	// unverändert aktiv und der dritte Weg unkonstruiert — die
	// Zwei-Bedingungen-Form (`natsStreamEnabled`) schützt bewusst gegen
	// eine versehentliche Ein-Bedingungs-Aktivierung. Ist die Variable
	// gesetzt, aber `envNatsURL` leer, ist das ein Konfigurationsfehler
	// beim Start (`ErrConfiguration`, `validateNatsStreamTokenRequiresURL`)
	// — ein Betreiber, der den dritten Weg aktiviert, aber kein
	// Verbindungsziel angibt, soll das beim Start bemerken.
	envNatsStreamToken = "CDC_NATS_STREAM_TOKEN"
	// envHTTPAddr trägt die optionale Horch-Adresse des HTTP/JSON-Driving-
	// Adapters (`ADR-0057`, `LH-FA-SST-006`): anders als die sechs
	// Vorbedingungen oben ist sie keine Start-Vorbedingung — ungesetzt
	// bleibt die API vollständig deaktiviert, kein `http.Server` wird
	// konstruiert (additiv, kein Breaking Change für bestehende
	// Deployments, analog zu `envNatsURL`/`ADR-0055` Punkt 5).
	envHTTPAddr = "CDC_HTTP_ADDR"
	// envAPITokenReader und envAPITokenAdmin tragen die beiden
	// Rechtsklassen der Token-Middleware (`ADR-0057` Teilfrage 3,
	// analog zum DB-Rollenmodell `ADR-0047`); beide bleiben optional wie
	// `envHTTPAddr` — ein leerer Wert deaktiviert die jeweilige Klasse
	// (`internal/adapters/driving/http`, `classifyToken`).
	envAPITokenReader = "CDC_API_TOKEN_READER"
	envAPITokenAdmin  = "CDC_API_TOKEN_ADMIN"
	// envGRPCAddr trägt die optionale Horch-Adresse des
	// gRPC-Streaming-Driving-Adapters (`ADR-0060`, `LH-FA-SST-008`): wie
	// `envHTTPAddr` ist sie keine Start-Vorbedingung — ungesetzt bleibt der
	// Streaming-Server vollständig deaktiviert, kein Listener wird geöffnet
	// (additiv, unverändertes Bestandsverhalten, `ADR-0060` Teilfrage 6).
	// Ein gesetzter Wert öffnet den Listener in eigener Goroutine (`Run`
	// unten); ein Startfehler wird dort über `log.Error` gemeldet und geht
	// nicht in das Ergebnis von `Run` ein.
	envGRPCAddr = "CDC_GRPC_ADDR"
)

// ErrConfiguration trägt die Fehlerklasse `configuration` der Verdrahtung
// (`SPEC-008`, `ADR-0023`): eine fehlende oder falsch gesetzte
// Vorbedingung endet ohne Start und ohne Fortsetzung im falschen Stand;
// der Prozess-Aufrufer meldet sie als Ausgang.
var ErrConfiguration = errors.New("Fehlerklasse configuration: Verdrahtung ohne vollständige Vorbedingung")

// heartbeatInterval trägt den periodischen Schreib-Zug des
// Heartbeat-Timers (`LH-FA-ADM-002`): ein MVP-Default ohne
// eigene Konfigurationsschicht — dieselbe Minimal-Form wie die übrigen
// Verdrahtungs-Vorbedingungen (Datei-Kommentar oben).
const heartbeatInterval = 5 * time.Second

// heartbeatStaleAfter trägt die Alters-Schwelle des `--healthcheck`-Laufs
// (Healthcheck): älter als das Dreifache des Schreib-Takts gilt als
// unhealthy — ein einzelner verpasster Takt (z. B. durch eine langsame
// Transaktion auf derselben Instanz) bleibt healthy, drei verpasste Takte
// in Folge nicht mehr. Ausführungsdetail der Composition-Root-Verdrahtung
// (`ADR-0026`), keine Architekturentscheidung.
const heartbeatStaleAfter = 3 * heartbeatInterval

// administrationPollInterval trägt den periodischen Fallback-Poll-Takt der
// Administrations-Goroutine (`runAdministration`, `ADR-0050`): derselbe
// Takt wie der Heartbeat-Schreib-Zug (`heartbeatInterval`) — ein
// verpasstes `NOTIFY` (z. B. nach einem Verbindungsabbruch der
// `LISTEN`-Verbindung) bleibt so höchstens einen Heartbeat-Takt lang
// unbemerkt, ohne einen zweiten Konfigurationswert ohne eigenen Bedarf
// einzuführen (Implementer-Entscheidung).
const administrationPollInterval = heartbeatInterval

// retentionInterval trägt den periodischen Lösch-Takt der
// Retention-Goroutine (`runRetentionCleanup`, `LH-FA-RET-002`…`004`,
// `ADR-0014`): ein MVP-Default ohne eigene Konfigurationsschicht, analog
// zu `heartbeatInterval` — Implementer-Entscheidung. Jeder Takt liest die
// Kandidaten der Quelle seitenweise ohne Row Images (`ADR-0124`) und befragt
// `RetentionPolicy.AllowsDeletion`; ein selteneres Intervall als der Heartbeat-Takt
// hält diese breitere Leseoperation von der kritischen Sektion des Capture-Pfads fern.
const retentionInterval = 10 * time.Second

// retentionMinAge trägt das Mindestalter der `RetentionPolicy`
// (`LH-FA-RET-003`): ein MVP-Default ohne eigene Konfigurationsschicht,
// dieselbe Minimal-Form wie `retentionInterval` — Implementer-
// Entscheidung. Eine Laufzeit-Konfigurationsanbindung ist ein anderer
// Vorgang.
const retentionMinAge = 24 * time.Hour

// walRetentionWarnBytes und walRetentionErrorBytes tragen die
// SPEC-013-Startwerte für `cdc_wal_retention_bytes` (`CDC_THRESHOLDS`,
// präzisiert durch `ADR-0049`(b)): unterhalb der Warnschwelle bleibt der
// periodische Schwellen-Vergleich unauffällig, zwischen Warn- und
// Fehlerschwelle setzt sich der Capture-Betrieb sichtbar fort (`SPEC-008`
// „kontrollierte Fortsetzung"), oberhalb der Fehlerschwelle klassifiziert
// `Run` den Lauf als `replication`-Fehler (`outbound.ErrReplication`) und
// bricht über den bestehenden Abbruchpfad ab (`classifyRunError`,
// `reportFault`). `Config.WALRetentionWarnBytes`/`WALRetentionErrorBytes`
// überschreiben diese Werte für den Ende-zu-Ende-Testlauf.
const (
	walRetentionWarnBytes  int64 = 100 * 1024 * 1024  // 100 MiB
	walRetentionErrorBytes int64 = 1024 * 1024 * 1024 // 1 GiB
)

// Config trägt die Verdrahtungs-Eingabe: drei rollen-spezifische
// Verbindungs-DSNs statt einer gemeinsamen Instanz-DSN (`ADR-0047`) — die
// Quelle und der CDC-Speicher bleiben dieselbe Instanz (MVP-Schnitt,
// Abschnitt 1 Lastenheft), die Rollentrennung greift auf Ebene der
// PostgreSQL-Anmeldung. Dazu die Quelle, die Verwaltungs-Namen
// Publication und Slot (`LH-FA-CFG-001.a`) und die aktivierten Tabellen
// mit ihren Port-Kennungen.
type Config struct {
	// CaptureDSN bindet an die Rolle `cdc_capture` (`ADR-0047`):
	// Store-Adapter, Replication-Stream und der ACK-Adapter (teilt die
	// Stream-Verbindung) verbinden sich darüber.
	CaptureDSN string
	// AdminDSN bindet an die Rolle `cdc_admin` (`ADR-0047`):
	// Tabellen-Aktivierung, Heartbeat-Adapter sowie die beiden
	// CLI-Sondermodi `RegisterConsumer`/`AcknowledgeConsumer` verbinden
	// sich darüber.
	AdminDSN string
	// ReaderDSN bindet an die Rolle `cdc_reader` (`ADR-0047`): der
	// `--healthcheck`-Lauf verbindet sich darüber.
	ReaderDSN   string
	Source      model.SourceID
	Publication string
	Slot        string
	Tables      map[string]mapper.TableBinding
	// LogLevel trägt den Level des JSON-Handlers, den der
	// Telemetrie-Driven-Adapter baut (`telemetry.New`, `ADR-0024`);
	// Herkunft ist `envLogLevel`/`parseLogLevel`, mit Default
	// `slog.LevelInfo`.
	LogLevel slog.Level
	// WALRetentionWarnBytes und WALRetentionErrorBytes überschreiben die
	// SPEC-013-Startwerte (`walRetentionWarnBytes`/`walRetentionErrorBytes`
	// unten, `ADR-0049`(b)) — für Testläufe, die beide Seiten der Schwelle
	// in vertretbarer Testzeit real durchlaufen wollen, statt 100 MiB/1 GiB
	// abzuwarten. Ungesetzt (0 oder negativ) übernimmt `Run` die
	// SPEC-013-Startwerte; kein Umgebungsname liest hierher — kein
	// Operator-Vertrag, siehe `ConfigFromEnv`.
	WALRetentionWarnBytes  int64
	WALRetentionErrorBytes int64
	// NatsURL trägt die optionale NATS-Server-URL des
	// Change-Notification-Wecksignals (`envNatsURL`, `ADR-0055`); leer
	// heißt Feature deaktiviert. Seine Herkunft ist auf beiden Ladepfaden
	// die Umgebungsvariable: die URL-Form kann Zugangsdaten einbetten und
	// gehört damit zur env-var-exklusiven Klasse (`ADR-0088` Festlegung 1).
	NatsURL string
	// HTTPAddr trägt die optionale Horch-Adresse des HTTP/JSON-Driving-
	// Adapters; leer heißt Feature deaktiviert — derselbe additive
	// Zuschnitt wie `NatsURL`. Herkunft ist `envHTTPAddr` oder das
	// Datei-Feld `http_addr`, mit Feld-für-Feld-Vorrang der
	// Umgebungsvariable (`ADR-0088` Festlegung 2/3).
	HTTPAddr string
	// APITokenReader und APITokenAdmin tragen die beiden Rechtsklassen
	// der Token-Middleware (`envAPITokenReader`/`envAPITokenAdmin`,
	// `ADR-0057` Teilfrage 3); leer heißt die jeweilige Klasse
	// deaktiviert. Beide sind Zugangsdaten und damit auf beiden
	// Ladepfaden env-var-exklusiv (`ADR-0088` Festlegung 1).
	APITokenReader string
	APITokenAdmin  string
	// GRPCAddr trägt die optionale Horch-Adresse des
	// gRPC-Streaming-Driving-Adapters; leer heißt Feature deaktiviert —
	// derselbe additive Zuschnitt wie `HTTPAddr`, mit derselben Herkunft
	// aus `envGRPCAddr` oder dem Datei-Feld `grpc_addr`. Der Adapter trägt
	// dieselben beiden Token-Klassen wie der HTTP-Adapter (`ADR-0060`
	// Teilfrage 4).
	GRPCAddr string
	// NatsStreamToken trägt den optionalen NATS-Verbindungs-Token des
	// dritten, vollinhaltstragenden Zustellwegs (`envNatsStreamToken`,
	// `ADR-0100` Teilfrage 4/5); leer heißt der dritte Weg deaktiviert. Er
	// trägt Zugangsdaten und ist damit auf beiden Ladepfaden env-var-exklusiv
	// (`ADR-0088` Festlegung 1), wie `NatsURL`/`APITokenReader`/`APITokenAdmin`.
	NatsStreamToken string
}

// ConfigFromEnv liest die Verdrahtungs-Vorbedingungen über die
// übergebene Umgebungs-Lese-Funktion; eine fehlende Vorbedingung endet
// über die Klasse `configuration` (ErrConfiguration).
func ConfigFromEnv(getenv func(string) string) (Config, error) {
	cfg := Config{
		CaptureDSN:  getenv(envCaptureDSN),
		AdminDSN:    getenv(envAdminDSN),
		ReaderDSN:   getenv(envReaderDSN),
		Source:      model.SourceID(getenv(envSource)),
		Publication: getenv(envPublication),
		Slot:        getenv(envSlot),
	}
	if cfg.CaptureDSN == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envCaptureDSN)
	}
	if cfg.AdminDSN == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envAdminDSN)
	}
	if cfg.ReaderDSN == "" {
		return Config{}, fmt.Errorf("%w: %s fehlt", ErrConfiguration, envReaderDSN)
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
	cfg.LogLevel = parseLogLevel(getenv(envLogLevel))
	cfg.NatsURL = getenv(envNatsURL)
	cfg.HTTPAddr = getenv(envHTTPAddr)
	cfg.APITokenReader = getenv(envAPITokenReader)
	cfg.APITokenAdmin = getenv(envAPITokenAdmin)
	cfg.GRPCAddr = getenv(envGRPCAddr)
	cfg.NatsStreamToken = getenv(envNatsStreamToken)
	if err := validateNatsStreamTokenRequiresURL(cfg.NatsURL, cfg.NatsStreamToken); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// validateNatsStreamTokenRequiresURL trägt `ADR-0100` Teilfrage 5s
// Konfigurationsfehler-Zweig: eine gesetzte `CDC_NATS_STREAM_TOKEN` ohne
// `CDC_NATS_URL` ist kein stiller Halbzustand — ein Betreiber, der den
// dritten Zustellweg aktiviert, aber kein Verbindungsziel angibt, bemerkt
// das beim Start. Geteilt zwischen `ConfigFromEnv` und `mergeConfig`
// (`config_file.go`): beide Ladepfade lesen beide Felder ausschließlich aus
// der Umgebung (`ADR-0088` Festlegung 1), die Vorbedingung gilt auf beiden
// gleich.
func validateNatsStreamTokenRequiresURL(natsURL, natsStreamToken string) error {
	if natsStreamToken != "" && natsURL == "" {
		return fmt.Errorf("%w: %s gesetzt, aber %s fehlt", ErrConfiguration, envNatsStreamToken, envNatsURL)
	}
	return nil
}

// natsStreamEnabled meldet, ob der dritte, vollinhaltstragende
// NATS-Zustellweg aktiv ist (`ADR-0100` Teilfrage 5): **beide** Bedingungen
// müssen gesetzt sein — `CDC_NATS_URL` und `CDC_NATS_STREAM_TOKEN`. Ist nur
// eine der beiden gesetzt, bleibt der Weg deaktiviert — Regressionsschutz
// gegen eine versehentliche Ein-Bedingungs-Aktivierung (`ADR-0100`
// §Fitness Function).
func natsStreamEnabled(natsURL, natsStreamToken string) bool {
	return natsURL != "" && natsStreamToken != ""
}

// parseLogLevel liest den Log-Level der strukturierten Ausgabe
// (`envLogLevel`, `LH-QA-OPS-004`): die Textformen trägt `slog.Level`
// selbst (`DEBUG`/`INFO`/`WARN`/`ERROR`, case-insensitive,
// `UnmarshalText`). Eine leere oder nicht erkannte Eingabe bleibt beim
// Default `Info` — anders als die Vorbedingungen oben (`ConfigFromEnv`)
// ist ein falsch gesetzter Level kein Start-Hindernis: er gefährdet die
// laufende Erfassung nicht, nur ihre Beobachtbarkeit.
func parseLogLevel(raw string) slog.Level {
	if raw == "" {
		return slog.LevelInfo
	}
	var level slog.Level
	if err := level.UnmarshalText([]byte(raw)); err != nil {
		return slog.LevelInfo
	}
	return level
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

// activatedTableBindings liest den committed Bindungsstand einer Quelle
// (`TableActivationPort.List`) und trägt ihn als `mapper.TableBinding`-Map
// fort — die Grundlage des Stream-Starts (`ADR-0050`): `CDC_TABLES` bleibt
// der Erstaktivierungs-Seed einer leeren Datenbank (oben in `Run`), dieser
// Aufruf liest danach den tatsächlichen Stand, einschließlich jeder
// zwischenzeitlich über SQL aktivierten Tabelle, deren Kennung `CDC_TABLES`
// nicht trägt. Eine aktivierte Tabelle ohne registrierte Schema-Version
// (widerspräche `TableActivationPort.Register`s Vertrag: beide Zeilen in
// einem Commit) bleibt ohne Bindung — keine Erfassung ohne
// Interpretationsgrundlage.
//
// Der Ausschlussstand kommt aus seiner dauerhaften Herkunft
// (`ColumnExclusionPort.ExcludedColumns`, `ADR-0065`) und geht je Tabelle
// in die Bindung ein — der Bindungs-Neuaufbau trägt damit auch einen
// früher beantragten Spaltenausschluss. Die Lese-Fähigkeit trägt dieselbe
// Adapter-Instanz wie die Aktivierung (im MVP eine Instanz, `ARC-004`);
// ein Lesefehler endet vor dem Stream-Start in der Startfehlerklasse des
// bestehenden Pfads (`storage`, `SPEC-008`).
//
// Der Regelstand kommt aus seiner dauerhaften Herkunft
// (`TransformationPort.TransformationRules`, `ADR-0112` Teilfrage 6) und
// geht wie der Ausschlussstand je Tabelle in die Bindung ein; eine Tabelle
// ohne Regel trägt eine Bindung ohne Regelstand. Die Lesung baut je Aufruf
// frische Listen, die die Bindung behält und die niemand mehr schreibt
// (`mapper.TableBinding`). Eine Regelform, die die Ableitung nicht mehr in
// eine Regel führt, endet den Start vor dem Stream-Start als Fehler der
// Klasse `internal` (`classifyRunError`); ein Lesefehler des Bestands endet
// dagegen in `storage`.
func activatedTableBindings(ctx context.Context, activation outbound.TableActivationPort, schemaStore outbound.SchemaStorePort, columnExclusion outbound.ColumnExclusionPort, transformations outbound.TransformationPort, source model.SourceID) (map[string]mapper.TableBinding, error) {
	registered, err := activation.List(ctx, source)
	if err != nil {
		return nil, err
	}
	excluded, err := columnExclusion.ExcludedColumns(ctx, source)
	if err != nil {
		return nil, err
	}
	rules, err := transformations.TransformationRules(ctx, source)
	if err != nil {
		return nil, err
	}
	tables := make(map[string]mapper.TableBinding, len(registered))
	for _, table := range registered {
		current, found, err := schemaStore.CurrentVersion(ctx, table.ID)
		if err != nil {
			return nil, err
		}
		if !found {
			continue
		}
		tables[table.QualifiedName()] = mapper.TableBinding{
			TableID:         table.ID,
			SchemaVersion:   current.ID,
			ExcludedColumns: excluded[table.QualifiedName()],
			Transformations: rules[table.QualifiedName()],
		}
	}
	return tables, nil
}

// splitQualifiedName teilt den qualifizierten Tabellennamen am ersten
// Punkt in Schema und Tabellenname; die Bindungs-Zeile trägt beide Teile
// getrennt (`LH-FA-DAT-002`), die Umgebung trägt den Namen als
// Schlüssel.
func splitQualifiedName(qualified string) (string, string, error) {
	schema, table, found := strings.Cut(qualified, ".")
	if !found || schema == "" || table == "" {
		return "", "", fmt.Errorf("%w: Aktivierung %q trägt nicht die Form schema.table", ErrConfiguration, qualified)
	}
	return schema, table, nil
}

// changeStreamEnabled meldet, ob die Live-Streaming-Fähigkeit aktiv ist
// (`ADR-0061` Teilfrage 5, um eine dritte Oder-Bedingung erweitert durch
// `ADR-0100` Teilfrage 5): Der `Broadcaster` wird konstruiert und über
// `CaptureService.WithChangeStream` verdrahtet, sobald mindestens einer von
// drei Zustellwegen aktiv ist — `CDC_GRPC_ADDR` gesetzt, `CDC_HTTP_ADDR`
// gesetzt, oder `natsStreamActive` (beide `CDC_NATS_URL` und
// `CDC_NATS_STREAM_TOKEN` gesetzt, `natsStreamEnabled`). Sind alle drei
// Bedingungen falsch, bleibt die Fähigkeit vollständig deaktiviert: kein
// Broadcaster, kein Stream-Publish-Schritt, unverändertes
// Bestandsverhalten.
func changeStreamEnabled(grpcAddr, httpAddr string, natsStreamActive bool) bool {
	return grpcAddr != "" || httpAddr != "" || natsStreamActive
}

// Run verdrahtet die Pipeline (`ADR-0026`) und trägt den Stream-Lauf bis
// zum Kontext-Ende: der Stream baut die Replication-Verbindung, der
// ACK-Adapter bestätigt über dieselbe Verbindung (`ADR-0007`, Option C)
// und der Capture Service orchestriert Persist-before-ACK
// (`LH-QA-REL-001.a`). Die Rückkehr ohne Fehler meldet das reguläre
// Lauf-Ende. Ein Adapter-Fehler wird durchgereicht, nicht still
// fortgesetzt (`SPEC-008`): der Prozess-Aufrufer endet auf jeden
// Adapter-Fehler mit Ausgang 1 — die Fortsetzung nach
// Verbindungsabbruch trägt der Prozess-Neustart, der Slot liest seinen
// Start über confirmed_flush_lsn (`ADR-0012`); die `transient`-Aktion
// (Erneut versuchen mit begrenztem Backoff, `SPEC-008`) trägt dieser
// Pfad nicht. Ein nicht-`nil`-Ausgang meldet zusätzlich
// den Fehlerzustand über den Heartbeat (`reportFault` unten,
// `LH-FA-ADM-003`, `LH-QA-REL-003`), bevor der Prozess-Aufrufer beendet —
// der benannte Rückgabewert `runErr` trägt dafür den Fehler über die
// `defer`-Kette hinweg.
func Run(ctx context.Context, cfg Config) (runErr error) {
	// Der Telemetrie-Driven-Adapter (`ADR-0024`: „Logging-/Metrics-
	// Frameworks bleiben Infrastruktur. … werden durch Driven Adapters
	// implementiert.", geschärft durch `ARC-011`) baut den JSON-Handler
	// (`telemetry.New`, `LH-QA-OPS-004`) — die Composition Root hält ihn
	// als lokale Variable und injiziert ihn über `WithLog`/`Config.Log`
	// in jeden Adapter-Konstruktor unten; kein Paket-globaler
	// Logging-Zustand — ein `slog.SetDefault` an dieser Stelle verletzt
	// `ADR-0024`, Alternative B.
	var log outbound.LogPort = telemetry.New(cfg.LogLevel)
	log.Info(ctx, "pg-change-feed: Verdrahtung gestartet", "source", string(cfg.Source))
	defer func() {
		if runErr != nil {
			log.Error(ctx, "pg-change-feed: Lauf beendet mit Fehler", "error", runErr)
			return
		}
		log.Info(ctx, "pg-change-feed: Lauf regulär beendet")
	}()

	store, err := postgresstorage.New(ctx, cfg.CaptureDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer store.Close()

	// Die dynamische Re-Versionierung im laufenden Erfassungspfad
	// (`mapper.Assembler.Consume`, `ADR-0015` Folgepflicht) liest und
	// schreibt über dieselbe Rolle wie Store und Stream (`cdc_capture`,
	// `ADR-0047`) — derselbe DSN, ein eigener Pool.
	schemaStore, err := postgresstorage.NewSchemaStore(ctx, cfg.CaptureDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer schemaStore.Close()

	// Die Aktivierung läuft als Use Case (`ADR-0028`, `LH-FA-CFG-001`):
	// die Bindungen der Konfiguration laufen vor dem Stream-Start als
	// EnableTable-Aufrufe — die Publication ist Start-Vorbedingung des
	// Stream-Adapters (`LH-FA-CFG-001.a`) und die Bindungs-Zeilen tragen
	// die Fremdschlüssel der ersten Persistenz (`SPEC-001`). Der Aufruf
	// ist idempotent (`LH-FA-CFG-001` Boundary) und trägt den Stand auch
	// nach einem Container-Neustart nach.
	// Die Verdrahtung trägt vier Verbindungen gegen dieselbe Instanz —
	// Store-Pool, Aktivierungs-Pool, Heartbeat-Pool, Stream-Verbindung —,
	// verteilt auf zwei PostgreSQL-Rollen (`ADR-0047`): Store, Stream und
	// der ACK-Adapter (der die Stream-Verbindung teilt) binden an
	// `cdc_capture`, Aktivierung und Heartbeat an `cdc_admin`. Das MVP
	// hält die Adapter-Lebenszyklen getrennt, statt einen Pool über die
	// Adapter zu teilen.
	activation, err := postgresstorage.NewTableActivation(ctx, cfg.AdminDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer activation.Close()
	// Der Heartbeat-Pool trägt ausschließlich den periodischen
	// Lebenszeichen-Zug (runHeartbeat, unten) — eine eigene Verbindung,
	// getrennt von Store- und Aktivierungs-Pool: der Timer-Zug teilt
	// keine Verbindung und keine Goroutine mit der
	// Capture-Persist-ACK-Schleife (`LH-QA-REL-001.a`).
	heartbeat, err := postgresstorage.NewHeartbeat(ctx, cfg.AdminDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer heartbeat.Close()
	// reportFault läuft vor heartbeat.Close() (LIFO-Reihenfolge der
	// `defer`-Kette: zuletzt registriert, zuerst ausgeführt) — der
	// Schreibversuch startet, bevor der Pool schließt. Kein
	// Zustellungs-Erfolg: reportFault ist best-effort (siehe dort) und
	// dasselbe Auslöse-Szenario (Quell-Instanz nicht erreichbar) kann den
	// nachfolgenden Fault-Schreibversuch am selben DSN ebenfalls scheitern
	// lassen.
	defer reportFault(heartbeat, cfg.Source, &runErr)
	// Der Retention-Pool trägt ausschließlich die periodische
	// Lösch-Ausführung (`runRetentionCleanup` unten) — eine eigene
	// Verbindung, gebunden an `cdc_admin` (`ADR-0047`): `DeleteChanges` und
	// die Waisen-Transaktions-Bereinigung sind ein Verwaltungs-, kein
	// Erfassungs-Nutzlast-Schreibzug (`tools/schema/nacharbeit-roles.sql`).
	retentionStore, err := postgresstorage.New(ctx, cfg.AdminDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer retentionStore.Close()
	// Eine eigene `ConsumerStatePort`-Verbindung, getrennt von den
	// kurzlebigen Verbindungen der CLI-Sondermodi
	// (`RegisterConsumer`/`AcknowledgeConsumer` unten): der Retention-Zug
	// liest die bestätigten Consumer-Positionen der Quelle wiederholt, über
	// die Lebensdauer des Prozesses.
	retentionConsumerState, err := postgresstorage.NewConsumerState(ctx, cfg.AdminDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer retentionConsumerState.Close()
	retentionMinAgeDuration, err := model.NewDuration(int64(retentionMinAge))
	if err != nil {
		return err
	}
	retentionPolicy, err := model.NewRetentionPolicy(retentionMinAgeDuration)
	if err != nil {
		return err
	}
	retentionUseCase := retention.NewRunRetentionService(retentionStore, retentionConsumerState, systemclock.New())
	enableTables := enable.NewEnableTableService(activation)
	disableTables := disable.NewDisableTableService(activation)
	excludeColumns := excludecolumn.NewExcludeColumnService(activation)
	includeColumns := includecolumn.NewIncludeColumnService(activation)
	setTransformations := settransformation.NewSetTransformationService(activation)
	removeTransformations := removetransformation.NewRemoveTransformationService(activation)
	for qualified, binding := range cfg.Tables {
		schema, table, err := splitQualifiedName(qualified)
		if err != nil {
			return err
		}
		if _, err := enableTables.Enable(ctx, inbound.EnableTableCommand{
			Source:          cfg.Source,
			Schema:          schema,
			Table:           table,
			TableID:         binding.TableID,
			SchemaVersionID: binding.SchemaVersion,
			// Die Umgebung trägt die Bindungs-Kennungen; die
			// Anfangs-Version trägt die Verdrahtung als erste Version
			// (`SPEC-004`) — spätere Versionen trägt die Schema-Evolution
			// über den Metadata-Pfad (`LH-FA-SCH-004.a`).
			Version:     1,
			Publication: cfg.Publication,
		}); err != nil {
			return err
		}
	}

	// Der laufende Bindungsstand des Assemblers trägt sich aus der
	// Datenbank fort (`TableActivationPort.List`/`SchemaStorePort.CurrentVersion`
	// samt `ColumnExclusionPort.ExcludedColumns` und
	// `TransformationPort.TransformationRules`), nicht ausschließlich aus
	// `cfg.Tables`: `CDC_TABLES` ist ausschließlich der Erstaktivierungs-Seed
	// oben (schreibt eine leere Datenbank fort; die Bindungen von
	// `cfg.Tables` gehen nur in die Aktivierung und erreichen den Assembler
	// nie), die Grundlage des
	// Stream-Starts ist der committed Stand — ein Prozess-Neustart verliert
	// damit weder eine zwischenzeitlich per SQL aktivierte Tabelle, deren
	// Kennung `CDC_TABLES` nicht trägt, noch einen dauerhaft vermerkten
	// Spaltenausschluss (`ADR-0065`) oder eine Transformationsregel
	// (`ADR-0112`). Die Aktivierungs-Instanz trägt die Lese-Fähigkeiten
	// (`ARC-004`).
	assemblerTables, err := activatedTableBindings(ctx, activation, schemaStore, activation, activation, cfg.Source)
	if err != nil {
		return err
	}

	stream, err := receive.NewStream(ctx, receive.Config{
		DSN:         cfg.CaptureDSN,
		Source:      cfg.Source,
		Publication: cfg.Publication,
		Slot:        cfg.Slot,
		Tables:      assemblerTables,
		SchemaStore: schemaStore,
		Log:         log,
	})
	if err != nil {
		return err
	}
	ack, err := postgresack.New(stream.Conn(), postgresack.WithLog(log))
	if err != nil {
		return err
	}

	// Der In-Prozess-`Broadcaster` (`ADR-0060` Teilfrage 2/5, dritter
	// Abonnent seit `ADR-0100`) ist die eine Stelle, an der die Zustellwege
	// zusammenlaufen: die beiden Driving-Adapter (gRPC-Server und
	// SSE-Endpunkt) und seit `ADR-0100` zusätzlich der Driven-Adapter
	// `natsstream.Publisher` lesen aus ihm, der `CaptureService` schreibt
	// über den Outbound Port `ChangeStreamPort` in ihn. Er wird
	// konstruiert, sobald mindestens einer der drei Zustellwege aktiv ist
	// (`changeStreamEnabled`, `ADR-0061` Teilfrage 5, `ADR-0100`
	// Teilfrage 5); sind alle drei Bedingungen falsch, entsteht kein
	// Broadcaster und der `CaptureService` trägt keinen Stream-Publish-
	// Schritt (additiv, unverändertes Bestandsverhalten). Die Bindung ist
	// keine Start-Vorbedingung: die Server starten weiter unten in eigener
	// Goroutine, ein Startfehler wird dort über `log.Error` gemeldet und
	// geht nicht in das Ergebnis von `Run` ein — wie beim HTTP-Adapter.
	natsStreamActive := natsStreamEnabled(cfg.NatsURL, cfg.NatsStreamToken)
	var changeBroadcaster *grpcstream.Broadcaster
	captureOpts := []capture.Option{capture.WithLog(log)}
	if changeStreamEnabled(cfg.GRPCAddr, cfg.HTTPAddr, natsStreamActive) {
		changeBroadcaster = grpcstream.New()
		captureOpts = append(captureOpts, capture.WithChangeStream(changeBroadcaster))
	}
	// Das Change-Notification-Wecksignal (`ADR-0055`, `LH-FA-SST-007`)
	// bleibt vollständig deaktiviert, solange `envNatsURL` leer ist — kein
	// Verbindungsversuch, kein `ChangeNotificationPort`. Ist die
	// Umgebungsvariable gesetzt, ist die Verbindung eine explizite
	// Vorbedingung dieses Laufs: anders als der Notify-Aufruf selbst später
	// (best-effort nach ACK, `ADR-0055` Punkt 4) meldet ein
	// Verbindungsfehler an dieser Stelle die Klasse `configuration`
	// (`ErrConfiguration`) — ein Betreiber, der das Feature einschaltet,
	// aber die Server-Adresse falsch trägt, soll das beim Start bemerken,
	// nicht durch ein unauffällig ausbleibendes Wecksignal. Dieselbe
	// Verbindung trägt seit `ADR-0100` zusätzlich den dritten,
	// vollinhaltstragenden Zustellweg — keine zweite Verbindung (Teilfrage
	// 5 Option B): eine gesetzte `CDC_NATS_STREAM_TOKEN` erweitert diesen
	// `nats.Connect`-Aufruf um eine Token-Client-Option, die der NATS-Server
	// serverweit erzwingt (bewusst benannter Nebeneffekt, `ADR-0100`
	// §Konsequenzen) — auch für diese bislang anonyme Wecksignal-Verbindung.
	stopNatsStreamPublisher := func() {}
	var natsStreamPublisherDone sync.WaitGroup
	// changeNotification trägt dasselbe Wecksignal für den Backfill-Run
	// (`ADR-0111` Teilfrage 5): `nil`, solange `envNatsURL` leer ist.
	var changeNotification outbound.ChangeNotificationPort
	if cfg.NatsURL != "" {
		var natsConnOpts []nats.Option
		if cfg.NatsStreamToken != "" {
			natsConnOpts = append(natsConnOpts, nats.Token(cfg.NatsStreamToken))
		}
		natsConn, err := nats.Connect(cfg.NatsURL, natsConnOpts...)
		if err != nil {
			return fmt.Errorf("%w: %s (%s) fehlgeschlagen: %v", ErrConfiguration, envNatsURL, cfg.NatsURL, err)
		}
		defer natsConn.Close()
		notify, err := natsnotify.New(natsConn, natsnotify.WithLog(log))
		if err != nil {
			return err
		}
		captureOpts = append(captureOpts, capture.WithChangeNotification(notify))
		changeNotification = notify

		// Der `natsstream.Publisher` (`ADR-0100` Teilfrage 1/5) entsteht nur
		// unter beiden Bedingungen (`natsStreamActive`) — dieselbe
		// Zwei-Bedingungen-Form wie `changeStreamEnabled` oben. Er
		// abonniert denselben `changeBroadcaster`, den `changeStreamEnabled`
		// bereits für genau diesen Fall konstruiert hat, und läuft als
		// eigener Hintergrund-Zug (eigene Goroutine, eigener WaitGroup-
		// Eintrag) — dieselbe Struktur wie Heartbeat/Administration/
		// Retention: kein Eingriff in die kritische Sektion des
		// Capture-Persist-ACK-Pfads (`LH-QA-REL-001.a`), weil er
		// ausschließlich aus dem bereits isolierten Broadcaster-Kanal
		// liest, nicht aus dem `CaptureService` selbst.
		if natsStreamActive {
			publisher, err := natsstream.New(natsConn, changeBroadcaster, string(cfg.Source), natsstream.WithLog(log))
			if err != nil {
				return err
			}
			natsStreamCtx, cancel := context.WithCancel(ctx)
			stopNatsStreamPublisher = cancel
			natsStreamPublisherDone.Add(1)
			go func() {
				defer natsStreamPublisherDone.Done()
				publisher.Run(natsStreamCtx)
			}()
		}
	}
	// Ein `CaptureService` trägt beide Eingänge des Streams: den Capture-Pfad
	// und die Leerlauf-Bestätigung (`ADR-0120`), beide über denselben
	// `ReplicationAckPort`.
	captureService := capture.NewCaptureService(store, ack, captureOpts...)
	if err := stream.BindCapture(captureService); err != nil {
		return err
	}
	if err := stream.BindIdleConfirmation(captureService); err != nil {
		return err
	}

	// Die Antrags-Queue-Verbindung (`cdc.administration_request`,
	// `ADR-0050`) trägt dieselbe Rolle wie Aktivierung und Heartbeat
	// (`cdc_admin`) — ein eigener Pool, getrennt von beiden: die
	// Administrations-Goroutine liest und schreibt unabhängig von ihrem
	// Schreib-Takt.
	adminRequests, err := postgresstorage.NewAdministrationRequest(ctx, cfg.AdminDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer adminRequests.Close()
	// Die `LISTEN`-Verbindung trägt eine eigene, dedizierte Verbindung
	// (keine Pool-Verbindung, siehe `AdministrationListener`) — dieselbe
	// Rolle `cdc_admin`.
	adminListener, err := postgresstorage.NewAdministrationListener(ctx, cfg.AdminDSN)
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = adminListener.Close(closeCtx)
	}()

	// Der Backfill (`LH-FA-CAP-009`, `ADR-0111`, `ADR-0113`): die Annahme
	// eines Antrags läuft über den Pool der Administrations-Goroutine
	// (`cdc_admin`), Run-Zustand und Schreiber laufen je über einen eigenen
	// Pool der Rolle `cdc_capture`, der Snapshot-Adapter liest über dieselbe
	// DSN — auch die geschätzte Zeilenzahl im Antrag liest er darüber, nicht
	// über den Pool der Administrations-Goroutine. Blockgröße und Zeitlimit
	// des Snapshot-Adapters tragen seine Startwerte.
	backfillAdmission, err := postgresstorage.NewBackfillAdmission(ctx, cfg.AdminDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer backfillAdmission.Close()
	backfillRuns, err := postgresstorage.NewBackfillRun(ctx, cfg.CaptureDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer backfillRuns.Close()
	backfillWriter, err := postgresstorage.NewBackfillWriter(ctx, cfg.CaptureDSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer backfillWriter.Close()
	backfillSnapshot, err := postgressnapshot.New(cfg.CaptureDSN)
	if err != nil {
		return err
	}
	clock := systemclock.New()
	backfillOpts := []backfill.Option{backfill.WithLog(log)}
	if changeNotification != nil {
		backfillOpts = append(backfillOpts, backfill.WithChangeNotification(changeNotification))
	}
	backfillTables := backfill.NewBackfillTableService(backfill.Ports{
		Activation:      activation,
		Exclusion:       activation,
		Transformations: activation,
		Schemas:         schemaStore,
		Snapshot:        backfillSnapshot,
		Admission:       backfillAdmission,
		Runs:            backfillRuns,
		Writer:          backfillWriter,
		Clock:           clock,
	}, backfillOpts...)

	// Start-Reihenfolge des Backfills (`ADR-0113` Festlegung 2): zuerst der
	// Bindungsaufbau der Composition Root (`assemblerTables`,
	// `stream.BindCapture`, oben), dann der Abgleich `running` → `interrupted`,
	// dann startet der Worker und liest die `queued`-Zeilen seiner Quelle.
	if err := reconcileBackfillRuns(ctx, backfillRuns, clock, cfg.Source, log); err != nil {
		return err
	}
	backfillWake := newBackfillWake()
	backfillCtx, stopBackfill := context.WithCancel(ctx)
	var backfillDone sync.WaitGroup
	backfillDone.Add(1)
	go func() {
		defer backfillDone.Done()
		runBackfillWorker(backfillCtx, backfillWorkerDeps{
			runs:        backfillRuns,
			useCase:     backfillTables,
			source:      cfg.Source,
			publication: cfg.Publication,
			wake:        backfillWake,
			retryAfter:  backfillRetryInterval,
			log:         log,
		})
	}()
	// Ein früher Rücksprung aus `Run` beendet den Worker, bevor die Pools
	// (`defer`-Kette oben) schließen.
	defer func() {
		stopBackfill()
		backfillDone.Wait()
	}()

	// Die Administrations-Goroutine läuft wie Heartbeat und WAL-Retention
	// über den eigenen Pool und die eigene Goroutine — kein Eingriff in
	// die kritische Sektion des Capture-Persist-ACK-Pfads
	// (`LH-QA-REL-001.a`). Sie trägt die laufende `Assembler`-Bindung
	// desselben Streams nach (`stream.Assembler()`), den `stream.Run`
	// unten konsumiert — kein zweiter Übersetzer.
	administrationCtx, stopAdministration := context.WithCancel(ctx)
	var administrationDone sync.WaitGroup
	administrationDone.Add(1)
	go func() {
		defer administrationDone.Done()
		runAdministration(administrationCtx, administrationDeps{
			requests:              adminRequests,
			listener:              adminListener,
			activation:            activation,
			enableTables:          enableTables,
			disableTables:         disableTables,
			excludeColumns:        excludeColumns,
			includeColumns:        includeColumns,
			schemaStore:           schemaStore,
			columnExclusion:       activation,
			transformations:       activation,
			assembler:             stream.Assembler(),
			setTransformations:    setTransformations,
			removeTransformations: removeTransformations,
			backfill:              backfillTables,
			backfillWake:          backfillWake,
			source:                cfg.Source,
			publication:           cfg.Publication,
			pollInterval:          administrationPollInterval,
			log:                   log,
		})
	}()

	// Die Retention-Goroutine läuft wie Heartbeat, Administration und
	// WAL-Retention über den eigenen Pool und die eigene Goroutine — kein
	// Eingriff in die kritische Sektion des Capture-Persist-ACK-Pfads
	// (`LH-QA-REL-001.a`). Der Lösch-Takt trägt `retentionInterval`, die
	// Freigabe je Change `RetentionPolicy.AllowsDeletion` über
	// `retentionPolicy` (`LH-FA-RET-002`…`004`, `ADR-0014`).
	retentionCtx, stopRetention := context.WithCancel(ctx)
	var retentionDone sync.WaitGroup
	retentionDone.Add(1)
	go func() {
		defer retentionDone.Done()
		runRetentionCleanup(retentionCtx, retentionUseCase, cfg.Source, retentionInterval, retentionPolicy, log)
	}()

	// Der HTTP/JSON-Driving-Adapter (`ADR-0057`, `LH-FA-SST-006`) bleibt
	// vollständig deaktiviert, solange `cfg.HTTPAddr` leer ist — kein
	// `http.Server` wird konstruiert, keine zusätzliche Verbindung
	// geöffnet (additiv, unverändertes Bestandsverhalten, analog zum
	// Change-Notification-Wecksignal oben, `ADR-0055` Punkt 5). Ist die
	// Adresse gesetzt, trägt der Adapter eine eigene `cdc_admin`-
	// Verbindung (`ADR-0047`) für die Consumer-Fähigkeiten
	// (Registrierung/Bestätigung/Position/Entfernung); die
	// Tabellen-Verwaltung und der Retention-Lauf nutzen dieselben
	// Use-Case-Instanzen wie die übrige Verdrahtung (`activation`,
	// `enableTables`, `disableTables`, `retentionUseCase` oben). Der
	// Adapter trägt daneben den SSE-Stream-Endpunkt (`LH-FA-SST-008`,
	// `ADR-0061`): er liest aus `changeBroadcaster`, demselben Broadcaster
	// wie der gRPC-Server.
	var httpServer *apihttp.Server
	var httpDone sync.WaitGroup
	if cfg.HTTPAddr != "" {
		apiConsumerState, err := postgresstorage.NewConsumerState(ctx, cfg.AdminDSN, postgresstorage.WithLog(log))
		if err != nil {
			return err
		}
		defer apiConsumerState.Close()
		// Der lesende Endpunkt `GET /changes` liest über denselben
		// `ChangeStorePort` wie der View-Direktzugriff (`ADR-0081`) — eine
		// eigene Verbindung derselben Rolle `cdc_admin`, die als einzige
		// `SELECT` auf `cdc.change`/`cdc.transaction` trägt
		// (`tools/schema/nacharbeit-roles.sql`); `cdc_reader` liest
		// ausschließlich die Views. Kein zweiter Lesepfad: derselbe
		// Adapter-Typ, derselbe Port.
		apiChangeStore, err := postgresstorage.New(ctx, cfg.AdminDSN, postgresstorage.WithLog(log))
		if err != nil {
			return err
		}
		defer apiChangeStore.Close()
		httpServer = apihttp.New(apihttp.Config{
			Addr:                cfg.HTTPAddr,
			TokenReader:         cfg.APITokenReader,
			TokenAdmin:          cfg.APITokenAdmin,
			RegisterConsumer:    register.NewRegisterConsumerService(apiConsumerState),
			AcknowledgeConsumer: acknowledge.NewAcknowledgeConsumerService(apiConsumerState),
			GetConsumerPosition: position.NewGetConsumerPositionService(apiConsumerState),
			RemoveConsumer:      remove.NewRemoveConsumerService(apiConsumerState),
			EnableTable:         enableTables,
			DisableTable:        disableTables,
			GetStatus:           status.NewGetStatusService(activation),
			ListTables:          list.NewListTablesService(activation),
			RunRetention:        retentionUseCase,
			ReadChanges:         readchanges.NewReadChangesService(apiChangeStore),
			Subscriber:          changeBroadcaster,
			Log:                 log,
		})
		httpDone.Add(1)
		go func() {
			defer httpDone.Done()
			if err := httpServer.Start(); err != nil {
				log.Error(ctx, "http: Adapter beendet mit Fehler", "error", err)
			}
		}()
	}

	// Der gRPC-Streaming-Driving-Adapter (`ADR-0060`, `LH-FA-SST-008`)
	// bleibt vollständig deaktiviert, solange `cfg.GRPCAddr` leer ist — kein
	// `grpc.Server`, kein Listener (additiv, unverändertes
	// Bestandsverhalten, analog zum HTTP-Adapter oben, `ADR-0060`
	// Teilfrage 6). Ist die Adresse gesetzt, liest der Server aus demselben
	// `changeBroadcaster`, den der `CaptureService` oben über
	// `WithChangeStream` bedient (`ADR-0060` Teilfrage 2). Der Adapter trägt
	// dieselben beiden Token-Klassen wie der HTTP-Adapter (`ADR-0060`
	// Teilfrage 4).
	var grpcServer *apigrpc.Server
	var grpcDone sync.WaitGroup
	if cfg.GRPCAddr != "" {
		grpcServer = apigrpc.New(apigrpc.Config{
			Addr:        cfg.GRPCAddr,
			TokenReader: cfg.APITokenReader,
			TokenAdmin:  cfg.APITokenAdmin,
			Subscriber:  changeBroadcaster,
			Log:         log,
		})
		grpcDone.Add(1)
		go func() {
			defer grpcDone.Done()
			if err := grpcServer.Start(); err != nil {
				log.Error(ctx, "grpc: Adapter beendet mit Fehler", "error", err)
			}
		}()
	}

	// Der WAL-Rückstand-Health-Check (`SPEC-009` `cdc_wal_retention_bytes`,
	// `ADR-0049`) braucht eine eigene Verbindung derselben Rolle
	// (`cdc_capture`) — die Stream-Verbindung steht während `stream.Run` im
	// COPY-Modus des Replication-Protokolls und nimmt keine Abfragen mehr
	// entgegen. Der Slot besteht an dieser Stelle bereits (`NewStream`
	// oben).
	walRetention, err := receive.NewWALRetentionChecker(ctx, cfg.CaptureDSN, cfg.Slot)
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = walRetention.Close(closeCtx)
	}()

	// Der Heartbeat-Zug läuft in einer eigenen Goroutine über den eigenen
	// Pool (oben) — kein Eingriff in die kritische Sektion des
	// Capture-Persist-ACK-Pfads (`LH-QA-REL-001.a`). `heartbeatCtx` endet
	// spätestens mit `stream.Run`; das Warten auf die Goroutine läuft
	// synchron vor der Rückkehr, damit der deferred `heartbeat.Close()`
	// oben nicht gegen einen noch schreibenden Aufruf läuft.
	heartbeatCtx, stopHeartbeat := context.WithCancel(ctx)
	var heartbeatDone sync.WaitGroup
	heartbeatDone.Add(1)
	go func() {
		defer heartbeatDone.Done()
		runHeartbeat(heartbeatCtx, heartbeat, cfg.Source, heartbeatInterval)
	}()

	// Derselbe Aufbau wie der Heartbeat-Zug: eigene Goroutine, eigene
	// Verbindung, an denselben Schreibtakt gebunden (`heartbeatInterval`)
	// — keine zweite Konfigurationsachse für ein zweites periodisches
	// Intervall ohne eigenen Bedarf. `streamCtx` ist eigens für den
	// Stream-Lauf abgeleitet (statt `ctx` direkt): der Schwellen-Vergleich
	// unten kann darüber den Stream-Lauf gezielt beenden (`stopStream`),
	// ohne den gesamten Prozess-Kontext zu kappen — ein regulärer
	// Prozess-Abbruch über `ctx` bricht `streamCtx` unverändert mit.
	streamCtx, stopStream := context.WithCancel(ctx)
	defer stopStream()

	warnBytes, errorBytes := resolveWALRetentionThresholds(cfg.WALRetentionWarnBytes, cfg.WALRetentionErrorBytes)
	var walFault walRetentionFault
	walRetentionCtx, stopWALRetention := context.WithCancel(ctx)
	var walRetentionDone sync.WaitGroup
	walRetentionDone.Add(1)
	go func() {
		defer walRetentionDone.Done()
		runWALRetentionCheck(walRetentionCtx, walRetention, log, heartbeatInterval, warnBytes, errorBytes, stopStream, &walFault)
	}()

	streamErr := stream.Run(streamCtx)
	stopHeartbeat()
	heartbeatDone.Wait()
	stopWALRetention()
	walRetentionDone.Wait()
	stopAdministration()
	administrationDone.Wait()
	stopBackfill()
	backfillDone.Wait()
	stopRetention()
	retentionDone.Wait()
	stopNatsStreamPublisher()
	natsStreamPublisherDone.Wait()
	if httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_ = httpServer.Shutdown(shutdownCtx)
		cancel()
		httpDone.Wait()
	}
	if grpcServer != nil {
		grpcServer.Shutdown()
		grpcDone.Wait()
	}
	return mergeStreamAndWALFaultOutcome(streamErr, &walFault)
}

// runHeartbeat schreibt das Lebenszeichen der Quelle periodisch fort, bis
// ctx endet (`LH-FA-ADM-002`). Ein Persistenzfehler des
// Heartbeats bricht den Aufruf nicht ab und wird verworfen: ein
// Schreibfehler des Heartbeats ist keine Fehlerklasse des Capture-Pfads
// (`SPEC-008`) — seine Abwesenheit zeigt sich stattdessen über das Alter
// der Lebenszeichen-Zeile (`cdc.heartbeat`, Healthcheck unten), nicht über
// einen abgebrochenen Stream-Lauf.
func runHeartbeat(ctx context.Context, port outbound.HeartbeatPort, source model.SourceID, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = port.Beat(ctx, source)
		}
	}
}

// walRetentionMeasurer entkoppelt den periodischen Schwellen-Vergleich von
// der konkreten Verbindung: `*receive.WALRetentionChecker` erfüllt dieses
// Interface bereits über seine `Measure`-Methode — Whitebox-Tests
// (`package bootstrap`) belegen `runWALRetentionCheck` damit ohne reale
// PostgreSQL-Instanz.
type walRetentionMeasurer interface {
	Measure(ctx context.Context) (int64, error)
}

// walRetentionLevel trägt die drei Zustände des Schwellen-Vergleichs
// (`SPEC-008` „kontrollierte Fortsetzung", `ADR-0049`(b)).
type walRetentionLevel int

const (
	walRetentionOK walRetentionLevel = iota
	walRetentionWarn
	walRetentionError
)

// resolveWALRetentionThresholds löst die effektiven Warn-/Fehlerschwellen
// auf: ein nicht gesetzter Override (0 oder negativ) übernimmt die
// SPEC-013-Startwerte (`walRetentionWarnBytes`/`walRetentionErrorBytes`)
// statt einer Schwelle von 0 — 0 würde jeden gemessenen Byte-Wert sofort als
// Fehlerschwellen-Überschreitung klassifizieren (`classifyWALRetention`) und
// den Lauf beim ersten Tick abbrechen lassen. Eigene Funktion statt Inline-
// Code in `Run`, damit der Fallback ohne reale PostgreSQL-Verbindung
// testbar ist.
func resolveWALRetentionThresholds(warnOverride, errorOverride int64) (warnBytes, errorBytes int64) {
	warnBytes, errorBytes = warnOverride, errorOverride
	if warnBytes <= 0 {
		warnBytes = walRetentionWarnBytes
	}
	if errorBytes <= 0 {
		errorBytes = walRetentionErrorBytes
	}
	return warnBytes, errorBytes
}

// classifyWALRetention vergleicht den gemessenen WAL-Rückstand gegen die
// beiden Schwellen aus `SPEC-013` (strikt größer als, wie dort formuliert:
// „Warn > 100 MiB · Fehler > 1 GiB") — dieselbe `>`-Semantik wie
// `healthcheckVerdict` oben (an der Schwelle selbst gilt die niedrigere
// Stufe).
func classifyWALRetention(bytes, warnBytes, errorBytes int64) walRetentionLevel {
	switch {
	case bytes > errorBytes:
		return walRetentionError
	case bytes > warnBytes:
		return walRetentionWarn
	default:
		return walRetentionOK
	}
}

// walRetentionFault trägt den WAL-Schwellen-Fehler thread-sicher von der
// periodischen Prüf-Goroutine (`runWALRetentionCheck`) zum Rückgabewert von
// `Run` (`mergeStreamAndWALFaultOutcome`): Die Fehlerschwelle wird aus einer
// eigenen Goroutine heraus festgestellt, `Run` liest sie erst, nachdem
// `stream.Run` zurückgekehrt ist. Nur der erste gesetzte Fehler bleibt
// erhalten — ein zweiter Tick über der Fehlerschwelle (vor der Rückkehr der
// bereits über `stopStream` beendeten Stream-Goroutine) trägt keine neue
// Information.
type walRetentionFault struct {
	mu  sync.Mutex
	err error
}

func (f *walRetentionFault) set(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err == nil {
		f.err = err
	}
}

func (f *walRetentionFault) get() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.err
}

// mergeStreamAndWALFaultOutcome trägt die Priorität zwischen dem
// Stream-Ausgang und einem aufgelaufenen WAL-Schwellen-Fehler: Ein
// Stream-Fehler — jede Klasse, ausdrücklich einschließlich einer
// Stream-Ordnungs-Verletzung (`mapper.ErrChangeWithoutBegin` u. ä.) —
// erreicht `Run`s Rückgabewert unverändert; ein WAL-Schwellen-Fehler kommt
// nur zum Zug, wenn der Stream-Lauf regulär endete (`nil`, ausgelöst über
// `stopStream` durch dieselbe Schwellen-Prüfung). Diese Reihenfolge trägt
// die Sentinel-Trennung aus `ADR-0049`(a) auf Ebene der Rückgabewert-
// Priorität: Der Schwellen-Fortsetzungspfad kann eine Stream-Ordnungs-
// Verletzung nie überschreiben oder verdecken.
func mergeStreamAndWALFaultOutcome(streamErr error, fault *walRetentionFault) error {
	if streamErr != nil {
		return streamErr
	}
	return fault.get()
}

// runWALRetentionCheck misst den WAL-Rückstand des Capture-Slots periodisch
// und vergleicht ihn gegen die Warn-/Fehlerschwelle (`SPEC-009`
// `cdc_wal_retention_bytes`, `SPEC-008` „kontrollierte Fortsetzung",
// `ADR-0049`): unterhalb der Warnschwelle protokolliert der Zug den Wert
// unauffällig — ein Betreiber liest ihn aus dem Log, ohne dass die Erhebung
// über `cdc.metrics`/`cdc_reader` läuft (dieselbe Begründung wie beim
// Health-Endpoint: `pg_replication_slots`/`IDENTIFY_SYSTEM` liegen
// außerhalb der Least-Privilege-Fläche von `cdc_reader`). Zwischen Warn-
// und Fehlerschwelle setzt der Capture-Betrieb sich sichtbar fort
// (Log-Warnung, kein Abbruch). Oberhalb der Fehlerschwelle setzt der Zug
// den WAL-Schwellen-Fehler (`fault.set`, `outbound.ErrReplication`), beendet
// den Stream-Lauf über `stopStream` und kehrt zurück — der bestehende
// Abbruchpfad (`classifyRunError`, `reportFault`) trägt den Rest. Eine
// fehlgeschlagene Messung selbst bricht den Aufruf nicht ab und wird
// verworfen — dieselbe best-effort-Haltung wie beim Heartbeat-Schreib-Zug
// oben. Eine dauerhaft gestörte eigene Verbindung des Checkers (z. B.
// Idle-Timeout zwischen zwei Ticks) muss diese Schleife nicht selbst
// behandeln: `WALRetentionChecker.Measure` ersetzt eine so gestörte
// Verbindung intern, sodass ein folgender Tick erneut misst statt dauerhaft
// zu verstummen.
func runWALRetentionCheck(ctx context.Context, checker walRetentionMeasurer, log outbound.LogPort, interval time.Duration, warnBytes, errorBytes int64, stopStream context.CancelFunc, fault *walRetentionFault) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			bytes, err := checker.Measure(ctx)
			if err != nil {
				log.Warn(ctx, "replication: WAL-Rückstand-Messung fehlgeschlagen", "error", err)
				continue
			}
			switch classifyWALRetention(bytes, warnBytes, errorBytes) {
			case walRetentionError:
				log.Error(ctx, "replication: WAL-Rückstand über Fehlerschwelle — kontrollierter Abbruch",
					"metric", "cdc_wal_retention_bytes", "bytes", bytes, "threshold_bytes", errorBytes)
				fault.set(fmt.Errorf("%w: WAL-Rückstand %d Bytes über Fehlerschwelle %d Bytes", outbound.ErrReplication, bytes, errorBytes))
				stopStream()
				return
			case walRetentionWarn:
				log.Warn(ctx, "replication: WAL-Rückstand über Warnschwelle — kontrollierte Fortsetzung",
					"metric", "cdc_wal_retention_bytes", "bytes", bytes, "threshold_bytes", warnBytes)
			default:
				log.Info(ctx, "replication: WAL-Rückstand gemessen", "metric", "cdc_wal_retention_bytes", "bytes", bytes)
			}
		}
	}
}

// runRetentionCleanup ruft die Retention-Bereinigung periodisch für die
// konfigurierte Quelle auf, bis `ctx` endet (`LH-FA-RET-002`…`004`,
// `ADR-0014`): jeder Tick befragt `RunRetentionUseCase.Run` — die Freigabe
// je Change trägt `RetentionPolicy.AllowsDeletion` im Use Case, diese
// Schleife trägt nur den periodischen Auslöser. Ein Fehler des Aufrufs
// bricht den Lauf nicht ab und wird protokolliert: eine gescheiterte
// Bereinigung ist kein Fehlerzustand des Capture-Pfads (`SPEC-008`) —
// dieselbe best-effort-Haltung wie beim Heartbeat-Schreib-Zug und der
// WAL-Rückstand-Messung oben.
func runRetentionCleanup(ctx context.Context, useCase inbound.RunRetentionUseCase, source model.SourceID, interval time.Duration, policy model.RetentionPolicy, log outbound.LogPort) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result, err := useCase.Run(ctx, inbound.RunRetentionCommand{Source: source, Policy: policy})
			if err != nil {
				log.Warn(ctx, "retention: Bereinigung fehlgeschlagen", "error", err)
				continue
			}
			log.Info(ctx, "retention: Bereinigung gelaufen", "deleted", result.Deleted)
		}
	}
}

// administrationListener entkoppelt `runAdministration` von der konkreten
// `LISTEN`-Verbindung (`*postgresstorage.AdministrationListener` erfüllt
// dieses Interface) — dasselbe Whitebox-Test-Muster wie
// `walRetentionMeasurer`: eine Fälschung belegt die Fallback-Poll-Schleife
// ohne reale PostgreSQL-Instanz.
type administrationListener interface {
	WaitForNotification(ctx context.Context) error
}

// administrationDeps bündelt die Abhängigkeiten der Administrations-
// Goroutine (`runAdministration`, `ADR-0050`) — ein eigener Typ statt
// einer langen Parameterliste, wie bei den übrigen periodischen Zügen
// dieser Datei.
type administrationDeps struct {
	requests       outbound.AdministrationRequestPort
	listener       administrationListener
	activation     outbound.TableActivationPort
	enableTables   inbound.EnableTableUseCase
	disableTables  inbound.DisableTableUseCase
	excludeColumns inbound.ExcludeColumnUseCase
	includeColumns inbound.IncludeColumnUseCase
	schemaStore    outbound.SchemaStorePort
	// columnExclusion trägt die dauerhafte Herkunft des Ausschlussstandes
	// (`ADR-0065`): der Aktivierungs-Zweig liest sie und trägt sie in die
	// neu angelegte `Assembler`-Bindung — derselbe Mechanismus wie der
	// Prozessstart (`activatedTableBindings`).
	columnExclusion outbound.ColumnExclusionPort
	// transformations trägt die dauerhafte Herkunft des Regelstandes
	// (`ADR-0112` Teilfrage 6): der Aktivierungs-Zweig liest sie und trägt
	// sie in die neu angelegte `Assembler`-Bindung — derselbe Mechanismus
	// wie der Prozessstart (`activatedTableBindings`).
	transformations outbound.TransformationPort
	assembler       *mapper.Assembler
	// setTransformations und removeTransformations prüfen die beiden
	// Transformations-Antragsarten (`LH-FA-CFG-007`); die Verarbeitung trägt
	// danach die Regel in die laufende `Assembler`-Bindung nach.
	setTransformations    inbound.SetTransformationUseCase
	removeTransformations inbound.RemoveTransformationUseCase
	// backfill nimmt einen Antrag der Art `backfill` an (`Request`) und
	// weckt danach den Backfill-Worker über `backfillWake`; die Ausführung
	// des Runs trägt der Worker, nicht diese Goroutine (`ADR-0111`
	// Teilfrage 5, `ADR-0113` Festlegung 2).
	backfill     inbound.BackfillTableUseCase
	backfillWake chan<- struct{}
	// source ist die Quelle dieser Instanz: ein `backfill`-Antrag einer
	// anderen Quelle bleibt `pending` für die Instanz, die ihn annimmt
	// (`ADR-0113` Festlegung 1, je Quelle nimmt eine Instanz Anträge an).
	source       model.SourceID
	publication  string
	pollInterval time.Duration
	log          outbound.LogPort
}

// runAdministration verarbeitet offene Anträge der Antrags-Queue
// (`cdc.administration_request`, `ADR-0050`), bis `ctx` endet: jeder
// Durchlauf verarbeitet zuerst die offenen Anträge, dann wartet er auf das
// nächste Wecksignal (`NOTIFY`) — mit dem Fallback-Poll-Takt als Timeout
// desselben Warte-Aufrufs, kein zweiter Ticker neben `WaitForNotification`.
// Ein abgelaufener Timeout und ein reales Wecksignal lösen denselben
// nächsten Verarbeitungs-Durchlauf aus; ein Verbindungsfehler der
// `LISTEN`-Verbindung wird protokolliert und bleibt bis zum nächsten
// erfolgreichen Wiederaufbau (`AdministrationListener`) durch den
// Fallback-Poll gedeckt.
func runAdministration(ctx context.Context, deps administrationDeps) {
	for {
		processAdministrationRequests(ctx, deps)
		if ctx.Err() != nil {
			return
		}
		waitCtx, cancel := context.WithTimeout(ctx, deps.pollInterval)
		err := deps.listener.WaitForNotification(waitCtx)
		cancel()
		if err != nil && ctx.Err() != nil {
			return
		}
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			deps.log.Warn(ctx, "administration: Wecksignal gestört — Fallback-Poll übernimmt", "error", err)
		}
	}
}

// processAdministrationRequests liest die offenen Anträge und verarbeitet
// jeden einzeln; ein Lesefehler bleibt best-effort (derselbe nächste
// Durchlauf versucht erneut). Ein gescheiterter Antrag wird als `failed`
// vermerkt, statt `pending` zu bleiben — ein `pending` bleibender Antrag
// würde jeden Durchlauf erneut versuchen, ohne dass sich der Fehlerzustand
// ändert. Ein `backfill`-Antrag einer anderen Quelle als `deps.source` bleibt
// unberührt `pending`.
func processAdministrationRequests(ctx context.Context, deps administrationDeps) {
	pending, err := deps.requests.ListPending(ctx)
	if err != nil {
		deps.log.Warn(ctx, "administration: Anträge lesen fehlgeschlagen", "error", err)
		return
	}
	for _, request := range pending {
		if request.Kind == model.AdministrationRequestBackfill && request.Source != deps.source {
			continue
		}
		if err := applyAdministrationRequest(ctx, deps, request); err != nil {
			deps.log.Warn(ctx, "administration: Antrag fehlgeschlagen",
				"request_id", request.ID, "kind", request.Kind, "error", err)
			if markErr := deps.requests.MarkFailed(ctx, request.ID, err.Error()); markErr != nil {
				deps.log.Warn(ctx, "administration: Fehlschlag nicht vermerkt", "request_id", request.ID, "error", markErr)
			}
			continue
		}
		if err := deps.requests.MarkApplied(ctx, request.ID); err != nil {
			deps.log.Warn(ctx, "administration: Erfolg nicht vermerkt", "request_id", request.ID, "error", err)
		}
	}
}

// processedAdministrationKinds nennt die Antragsarten, die
// `applyAdministrationRequest` verarbeitet — die geschlossene Menge der
// Domäne (`model.AdministrationRequestKinds`); der Fehlertext des
// `default`-Zweigs trägt sie. Dass jede dieser Arten einen Zweig des
// `switch` trägt und nicht im `default` endet, bindet
// `TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet`.
func processedAdministrationKinds() string {
	kinds := model.AdministrationRequestKinds()
	names := make([]string, 0, len(kinds))
	for _, kind := range kinds {
		names = append(names, string(kind))
	}
	return strings.Join(names, "/")
}

// applyAdministrationRequest führt einen einzelnen Antrag über den
// passenden Inbound Port aus (einziger Schreibpfad auf Bindungs-Zeile und
// Publication bleibt der Port, `ADR-0018`/`ADR-0046` unverändert) und
// trägt bei Erfolg die laufende `Assembler`-Bindung nach. Die
// Bindungs-Kennungen einer SQL-beantragten Aktivierung liegen nicht am
// Antrags-Datensatz (anders als bei `CDC_TABLES`) — `administrationTableID`
// trägt eine deterministische, wiederholbare Kennung: ein erneuter
// Durchlauf über denselben (bereits verarbeiteten) Antrag — etwa nach
// einem Prozess-Neustart, bevor der vorige Durchlauf den Vermerk schreiben
// konnte — trägt dieselbe Kennung und trifft über `EnableTableUseCase`s
// Idempotenz (`LH-FA-CFG-001` Boundary) dieselbe Zeile. Nach `Enable`
// liest der Aufruf die tatsächlich registrierte Bindung über
// `TableActivationPort.Registered` zurück, statt der soeben übergebenen
// Kennung blind zu vertrauen — bereits vor diesem Antrag über `CDC_TABLES`
// aktivierte Tabellen tragen sonst die falsche Kennung in der
// nachgetragenen `Assembler`-Bindung.
//
// Die beiden Spalten-Antragsarten rufen ihren Use Case auf und tragen
// danach den Ausschlussstand in die laufende `Assembler`-Bindung nach: ihr
// Ziel ist der Filterzustand der laufenden Erfassung, nicht die Bindungs-
// oder Publication-Menge, die die beiden Tabellen-Antragsarten tragen
// (`ADR-0059` Teilfrage 3). Die Nachträge greifen unter `tablesMu` —
// derselbe synchronisierte Schreibpfad wie `AddBinding`/`RemoveBinding`.
//
// Der Aktivierungs-Zweig liest den dauerhaften Ausschlussstand der Tabelle
// (`ADR-0065`) und ihren Regelstand (`ADR-0112` Teilfrage 6) und übergibt
// beide an `AddBinding`: `RemoveBinding` hat den Bindungs-Eintrag samt
// Ausschluss- und Regelstand entfernt, dieser Zweig legt sie über die
// Herkunft neu an — die Stände überleben den `disable`/`enable`-Zyklus ohne
// Neustart. Ein Lese-Fehler endet wie jeder Antrags-Fehler im
// `failed`-Vermerk (`processAdministrationRequests`).
//
// Der `backfill`-Zweig nimmt den Antrag über `BackfillTableUseCase.Request`
// an und weckt danach den Backfill-Worker; die Kopie läuft dort, nicht in
// dieser Goroutine. Ein Antrag, den `Request` annimmt, ist damit `applied` im
// Sinn von „angenommen“ (die Annahme vermerkt ihn in ihrer Transaktion; das
// anschließende `MarkApplied` trifft keine `pending`-Zeile mehr und ist
// kein Fehler). Eine verletzte Vorbedingung und ein aktiver Run derselben
// Tabelle enden als Fehler und damit im `failed`-Vermerk.
//
// Die Antragsarten `set_transformation`/`remove_transformation`
// (`LH-FA-CFG-007`) rufen ihren Use Case auf, der Form und Konfliktfreiheit
// K1 bis K4 (`SPEC-019`) prüft, und tragen danach die geprüfte Regel bzw. das
// Herausnehmen in den Regelstand der laufenden `Assembler`-Bindung nach — wie
// die Spalten-Antragsarten unter `tablesMu`, ohne Bindungs- oder
// Publication-Änderung. Eine Tabelle ohne laufende Bindung lässt den Nachtrag
// wirkungslos; der Antrag endet `applied`, und der Regelstand entsteht beim
// Anlegen der Bindung aus der dauerhaften Herkunft (Prozessstart,
// Aktivierungs-Zweig unten). Die Wiederholung eines bereits nachgetragenen,
// noch `pending` stehenden Antrags ist folgenlos: der Nachtrag ersetzt die
// Regel nach ihrem Namen, K1 prüft nur gegen `applied`-Zeilen. Eine Zeile mit
// leerem Regelnamen oder leerer Regelform erreicht diesen Zweig als Antrag
// (`model.NewAdministrationRequest`) und endet im Use Case als Fehler mit dem
// Fehlertext der Spec, ohne die Queue anzuhalten.
//
// Der `default`-Zweig endet als Fehler und damit im `failed`-Vermerk, dessen
// Text die verarbeiteten Antragsarten nennt; er trifft nur eine Antragsart
// außerhalb der geschlossenen Menge der Domäne.
func applyAdministrationRequest(ctx context.Context, deps administrationDeps, request model.AdministrationRequest) error {
	qualified := request.Schema + "." + request.Table
	switch request.Kind {
	case model.AdministrationRequestEnable:
		tableID := administrationTableID(request.Schema, request.Table)
		if _, err := deps.enableTables.Enable(ctx, inbound.EnableTableCommand{
			Source:          request.Source,
			Schema:          request.Schema,
			Table:           request.Table,
			TableID:         tableID,
			SchemaVersionID: administrationSchemaVersionID(tableID),
			Version:         1,
			Publication:     deps.publication,
		}); err != nil {
			return err
		}
		registered, found, err := deps.activation.Registered(ctx, request.Source, request.Schema, request.Table)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("Aktivierung ohne Bindungs-Zeile: %s", qualified)
		}
		current, found, err := deps.schemaStore.CurrentVersion(ctx, registered.ID)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("Aktivierung ohne registrierte Schema-Version: %s", qualified)
		}
		excluded, err := deps.columnExclusion.ExcludedColumns(ctx, request.Source)
		if err != nil {
			return err
		}
		rules, err := deps.transformations.TransformationRules(ctx, request.Source)
		if err != nil {
			return err
		}
		deps.assembler.AddBinding(qualified, mapper.TableBinding{
			TableID:         registered.ID,
			SchemaVersion:   current.ID,
			ExcludedColumns: excluded[qualified],
			Transformations: rules[qualified],
		})
		return nil
	case model.AdministrationRequestDisable:
		if _, err := deps.disableTables.Disable(ctx, inbound.DisableTableCommand{
			Source:      request.Source,
			Schema:      request.Schema,
			Table:       request.Table,
			Publication: deps.publication,
		}); err != nil {
			return err
		}
		deps.assembler.RemoveBinding(qualified)
		return nil
	case model.AdministrationRequestExcludeColumn:
		if err := deps.excludeColumns.Exclude(ctx, inbound.ExcludeColumnCommand{
			Source: request.Source,
			Schema: request.Schema,
			Table:  request.Table,
			Column: request.Column,
		}); err != nil {
			return err
		}
		deps.assembler.ExcludeColumn(qualified, request.Column)
		return nil
	case model.AdministrationRequestIncludeColumn:
		if err := deps.includeColumns.Include(ctx, inbound.IncludeColumnCommand{
			Source: request.Source,
			Schema: request.Schema,
			Table:  request.Table,
			Column: request.Column,
		}); err != nil {
			return err
		}
		deps.assembler.IncludeColumn(qualified, request.Column)
		return nil
	case model.AdministrationRequestBackfill:
		if deps.backfill == nil {
			return errBackfillNotWired
		}
		if _, err := deps.backfill.Request(ctx, inbound.BackfillRequestCommand{
			RequestID:   request.ID,
			Source:      request.Source,
			Schema:      request.Schema,
			Table:       request.Table,
			Publication: deps.publication,
		}); err != nil {
			return err
		}
		signalBackfillWorker(deps.backfillWake)
		return nil
	case model.AdministrationRequestSetTransformation:
		// K3 hat gegen die Spaltenliste zum Antragszeitpunkt geprüft: eine
		// spätere Spalten-Erweiterung, die den Zielnamen kollidieren lässt,
		// erreicht diese Prüfung nicht — der Assembler fängt sie als nicht
		// anwendbare Regel im Erfassungspfad (`SPEC-030`, Anwendbarkeit).
		// Der Nachtrag steht vor dem Vermerk `applied`: scheitert der Vermerk,
		// trägt die laufende Bindung die Regel weiter, und der nächste
		// Durchlauf verarbeitet denselben Antrag erneut (der Ersatz nach Namen
		// macht die Wiederholung folgenlos). Grenze: verarbeitet ein Durchlauf
		// dazwischen einen Antrag, der zum noch nicht vermerkten in K2 oder K3
		// steht (dieselbe Spalte, dasselbe Ziel), prüft er gegen einen
		// Regelstand ohne den ersten und wird `applied`; die Wiederholung des
		// ersten endet danach `failed`, und die laufende Bindung trägt bis zum
		// Neustart beide Regeln.
		rule, err := deps.setTransformations.Set(ctx, inbound.SetTransformationCommand{
			Source:   request.Source,
			Schema:   request.Schema,
			Table:    request.Table,
			RuleName: request.RuleName,
			RuleSpec: request.RuleSpec,
		})
		if err != nil {
			return err
		}
		deps.assembler.SetTransformation(qualified, rule)
		return nil
	case model.AdministrationRequestRemoveTransformation:
		if err := deps.removeTransformations.Remove(ctx, inbound.RemoveTransformationCommand{
			Source:   request.Source,
			Schema:   request.Schema,
			Table:    request.Table,
			RuleName: request.RuleName,
		}); err != nil {
			return err
		}
		deps.assembler.RemoveTransformation(qualified, request.RuleName)
		return nil
	default:
		return fmt.Errorf("Antragsart %q gehört nicht zu den verarbeiteten Antragsarten %s", request.Kind, processedAdministrationKinds())
	}
}

// administrationTableID trägt die deterministische Tabellen-Kennung einer
// über SQL beantragten Aktivierung: der qualifizierte Name selbst — anders
// als `CDC_TABLES`, das eine frei gewählte Kennung aus der Umgebung trägt,
// hält ein Antrags-Datensatz nur Schema und Tabellenname.
func administrationTableID(schema, table string) model.SourceTableID {
	return model.SourceTableID(schema + "." + table)
}

// administrationSchemaVersionID trägt die Anfangs-Version-Kennung einer
// über SQL beantragten Aktivierung — dasselbe Kennungsformat wie die
// dynamische Re-Versionierung des Mappers (Tabellen-Kennung, getrennt
// durch `-v`, Versionsnummer), hier für die erste Version.
func administrationSchemaVersionID(table model.SourceTableID) model.SchemaVersionID {
	return model.SchemaVersionID(fmt.Sprintf("%s-v1", table))
}

// reportFault meldet einen nicht-`nil` Lauf-Fehler als Fehlerzustand über
// den Heartbeat (`LH-FA-ADM-003`, `LH-QA-REL-003`) — die
// Sichtbarkeit gilt für „Erfassung kann nicht fortsetzen" (`Run` endet auf
// jeden Adapter-Fehler, Dateikommentar oben), nicht für einen regulären
// Lauf-Abschluss (`runErr == nil`). Der Schreib-Zug trägt eine eigene,
// kurzlebige Frist, unabhängig von `ctx`: das Beenden des Prozesses selbst
// hat `ctx` bereits abgebrochen, und genau dann muss der Fehlerzustand noch
// geschrieben werden können (dasselbe Muster wie `Healthcheck` unten). Ein
// Persistenzfehler dieses Schreib-Zugs bleibt best-effort und unterdrückt
// nicht den eigentlichen Lauf-Fehler, den `*runErr` weiterhin trägt.
func reportFault(port outbound.HeartbeatPort, source model.SourceID, runErr *error) {
	if runErr == nil || *runErr == nil {
		return
	}
	faultCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = port.Fault(faultCtx, source, classifyRunError(*runErr))
}

// classifyRunError übersetzt den Lauf-Fehler in eine der sieben stabilen
// Kategorien aus `ADR-0023`/`SPEC-008`: die Composition Root kennt die
// Sentinel-Fehler aller beteiligten Adapter (`.a-check.yml`
// `composition_root` — kein Hexagon-Schichten-Edge, der diese Referenz
// einschränkt) und übersetzt sie in den Fehlerzustand, den `reportFault`
// oben fortträgt. Ein nicht erkannter Fehler bleibt in der Kategorie
// `internal` (`SPEC-008`: „unerwarteter interner Fehler").
func classifyRunError(err error) model.ErrorClass {
	switch {
	case errors.Is(err, ErrConfiguration),
		errors.Is(err, receive.ErrConfiguration),
		errors.Is(err, postgresstorage.ErrActivationConfiguration):
		return model.ErrorClassConfiguration
	case errors.Is(err, decode.ErrSchema),
		errors.Is(err, mapper.ErrTruncateUnsupported),
		errors.Is(err, mapper.ErrIncompatibleSchemaChange),
		errors.Is(err, mapper.ErrTransformationNotApplicable):
		return model.ErrorClassSchema
	case errors.Is(err, receive.ErrReplication),
		errors.Is(err, outbound.ErrReplication),
		errors.Is(err, mapper.ErrChangeWithoutBegin),
		errors.Is(err, mapper.ErrCommitWithoutBegin),
		errors.Is(err, mapper.ErrBeginWithoutCommit):
		return model.ErrorClassReplication
	case errors.Is(err, outbound.ErrStorage),
		errors.Is(err, outbound.ErrHeartbeatStorage),
		errors.Is(err, outbound.ErrConsumerStateStorage):
		return model.ErrorClassStorage
	default:
		return model.ErrorClassInternal
	}
}

// RegisterConsumer verdrahtet den `RegisterConsumerUseCase` (`ADR-0028`)
// für den `register-consumer`-Sondermodus
// (`cmd/pg-change-feed/main.go`, `LH-FA-CON-001.a`) und trägt dessen
// Prozess-Ausgang: 0 nach Registrierung — neu oder bereits registriert,
// beide Fälle tragen `LH-FA-CON-001`s Boundary-Kriterium, die Ausgabe
// unterscheidet sie —, 1 bei Verdrahtungs- oder Domänenfehler. Der Aufruf
// öffnet eine eigene, kurzlebige Verbindung über denselben
// `ConsumerStatePort`-Adapter, den die laufende Verdrahtung (`Run` oben)
// nutzen würde — kein Bestandteil des Dauerbetriebs, gebunden an die Rolle
// `cdc_admin` (`ADR-0047`: DML auf `consumer`/`consumer_position`).
// Consumer-Kennung und -Name tragen denselben Wert (`name`); dieser
// Zugriffsweg trennt beide (noch) nicht.
func RegisterConsumer(ctx context.Context, cfg Config, name string) int {
	state, err := postgresstorage.NewConsumerState(ctx, cfg.AdminDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: register-consumer: %v\n", err)
		return 1
	}
	defer state.Close()

	result, err := register.NewRegisterConsumerService(state).Register(ctx, inbound.RegisterConsumerCommand{
		Consumer: model.ConsumerID(name),
		Name:     name,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: register-consumer: %v\n", err)
		return 1
	}
	if result.AlreadyRegistered {
		fmt.Printf("pg-change-feed: Consumer %q bereits registriert\n", result.Consumer.ID)
		return 0
	}
	fmt.Printf("pg-change-feed: Consumer %q registriert\n", result.Consumer.ID)
	return 0
}

// AcknowledgeConsumer verdrahtet den `AcknowledgeConsumerUseCase`
// (`ADR-0028`) für den `acknowledge-consumer`-Sondermodus
// (`cmd/pg-change-feed/main.go`, `LH-FA-CON-004.a`) und trägt dessen
// Prozess-Ausgang: 0 nach Bestätigung — neu vorgerückt oder eine
// Wiederholung derselben Position, beide Fälle tragen `LH-FA-CON-004`s
// Boundary-Kriterium (Idempotenz), die Ausgabe unterscheidet sie nicht
// gesondert —, 1 bei Verdrahtungs- oder Domänenfehler, unter anderem ein
// echter Rückschritt (`domainerrors.ErrPositionRegression`) oder eine
// nicht registrierte Kennung (`outbound.ErrConsumerUnregistered`). Der
// Aufruf öffnet eine eigene, kurzlebige Verbindung über denselben
// `ConsumerStatePort`-Adapter, den die laufende Verdrahtung (`Run` oben)
// nutzen würde — kein Bestandteil des Dauerbetriebs, gebunden an die Rolle
// `cdc_admin` (`ADR-0047`: DML auf `consumer`/`consumer_position`). Die
// bestätigte Position trägt dieselbe Quelle wie die laufende Erfassung
// (`cfg.Source`); der Zugriffsweg unterscheidet keine zweite Quelle.
func AcknowledgeConsumer(ctx context.Context, cfg Config, consumerID string, offset uint64) int {
	state, err := postgresstorage.NewConsumerState(ctx, cfg.AdminDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: acknowledge-consumer: %v\n", err)
		return 1
	}
	defer state.Close()

	result, err := acknowledge.NewAcknowledgeConsumerService(state).Acknowledge(ctx, inbound.AcknowledgeConsumerCommand{
		Consumer: model.ConsumerID(consumerID),
		Position: model.SourcePosition{SourceID: cfg.Source, Offset: offset},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: acknowledge-consumer: %v\n", err)
		return 1
	}
	fmt.Printf("pg-change-feed: Consumer %q Position bestätigt (Quelle %q, Offset %d)\n",
		consumerID, result.Position.Position.SourceID, result.Position.Position.Offset)
	return 0
}

// healthcheckVerdict trägt die binäre Healthcheck-Entscheidung (Compose
// Exit-Code 0/1) über das Alter des letzten Lebenszeichens — reiner
// Vergleich ohne Verbindungsversuch, testbar ohne reale Instanz.
func healthcheckVerdict(age time.Duration) int {
	if age > heartbeatStaleAfter {
		return 1
	}
	return 0
}

// Healthcheck liest das Alter des letzten Lebenszeichens über die
// SQL-Lese-View `cdc.heartbeat` (`LH-FA-ADM-002`, `LH-QA-OPS-002`,
// `ADR-0046` Kategorie C) und trägt den Prozess-Ausgang des
// `--healthcheck`-Laufs (`cmd/pg-change-feed/main.go`): 0 (healthy)
// unterhalb der Schwelle, 1 sonst — ein Verbindungsfehler, eine fehlende
// Zeile (die Instanz hat noch nie geschlagen) und eine nicht lesbare View
// (z. B. Schema-Rollout nicht gelaufen) gelten als nicht gesund. Jede der
// drei Fehlerklassen trägt eine eigene, kurze `stderr`-Zeile — ein
// Docker-Healthcheck-Fail unterscheidet sich sonst nicht von echter
// Staleness; kein neues Fehlerklassen-Schema, nur Diagnose-Text. Die Klassifikation `SPEC-007`
// (`HEALTH_STATES`) bleibt Sache des lesenden Systems; der Exit-Code
// trägt nur die binäre Compose-Semantik. Der Aufruf öffnet eine eigene,
// kurzlebige Verbindung — kein Bestandteil der laufenden Verdrahtung
// (`Run` oben); der Aufrufer übergibt `cfg.ReaderDSN` (`ADR-0047`: die
// View-Lesung deckt sich mit dem `SELECT`-Grant der Rolle `cdc_reader`).
func Healthcheck(ctx context.Context, dsn string, source model.SourceID) int {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: healthcheck: DSN ungültig: %v\n", err)
		return 1
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: healthcheck: Instanz nicht erreichbar: %v\n", err)
		return 1
	}
	var ageSeconds float64
	err = pool.QueryRow(ctx, "SELECT age_seconds FROM cdc.heartbeat WHERE source_id = $1", string(source)).Scan(&ageSeconds)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Fprintf(os.Stderr, "pg-change-feed: healthcheck: kein Lebenszeichen für Quelle %q — Instanz hat noch nie geschlagen\n", source)
		} else {
			fmt.Fprintf(os.Stderr, "pg-change-feed: healthcheck: cdc.heartbeat nicht lesbar (Schema-Rollout gelaufen?): %v\n", err)
		}
		return 1
	}
	age := time.Duration(ageSeconds * float64(time.Second))
	if verdict := healthcheckVerdict(age); verdict != 0 {
		fmt.Fprintf(os.Stderr, "pg-change-feed: healthcheck: Lebenszeichen veraltet (%s, Schwelle %s)\n", age, heartbeatStaleAfter)
		return verdict
	}
	return 0
}

// Diagnose liest die in `LH-FA-SST-003`s Boundary genannten
// Status-/Diagnosesignale über dieselben SQL-Lese-Views wie `Healthcheck`
// oben (`cdc.heartbeat`, `cdc.metrics`) sowie `cdc.retention_blockers` und
// gibt sie menschenlesbar auf `stdout` aus: Betriebsstatus und
// Fehlerzustand (`LH-FA-ADM-002`/`003`, aus `cdc.heartbeat`), CDC-Abstand
// (`LH-FA-ADM-004`, `cdc_capture_lag`), Verarbeitungsrückstand je Consumer
// (`LH-FA-ADM-005`, `cdc_consumer_lag{consumer}`), der aktuell die Löschung
// blockierende Consumer je Quelle (`LH-FA-RET-005`, `cdc.retention_blockers`)
// und der Speicherverbrauch (`LH-FA-RET-006`, `cdc_storage_bytes`) sowie der
// zuletzt beantragte Backfill-Run je Tabelle (`LH-FA-CAP-009`,
// `cdc.backfill_status`, `diagnoseBackfillStatus`). Anders
// als `Healthcheck` trifft dieser Befehl keine binäre Verdikt-Entscheidung —
// er gibt die Rohwerte aller Views unverändert weiter, keine
// Schwellenwert-Klassifikation (`SPEC-007` bleibt Sache des lesenden
// Systems, wie bei den Views selbst). Der Prozess-Ausgang trägt nur den
// Lese-Erfolg: 0 nach vollständig gelesenen Views — unabhängig vom Inhalt,
// ein gemeldeter Fehlerzustand oder Rückstand ist Berichtsinhalt, kein
// Befehlsfehler —, 1 bei Verbindungs- oder Query-Fehler (dieselben ersten
// beiden Fehlerklassen wie `Healthcheck`: DSN ungültig, Instanz nicht
// erreichbar, plus eine dritte je nicht lesbarer View). Eine Quelle ohne
// Zeile in `cdc.retention_blockers` (kein Consumer hat je gegen sie
// bestätigt) meldet „kein Blocker" — dieselbe Abwesenheits-Lesart wie beim
// Betriebsstatus oben, kein Fehlerzustand. Der Aufruf öffnet eine eigene,
// kurzlebige Verbindung — kein Bestandteil der laufenden Verdrahtung
// (`Run` oben); der Aufrufer übergibt `cfg.ReaderDSN` (`ADR-0047`: alle
// vier Views tragen ein `SELECT`-Grant an `cdc_reader`).
func Diagnose(ctx context.Context, dsn string, source model.SourceID) int {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: DSN ungültig: %v\n", err)
		return 1
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: Instanz nicht erreichbar: %v\n", err)
		return 1
	}

	fmt.Printf("pg-change-feed diagnose: Quelle %q\n", source)

	var ageSeconds float64
	var errorClass *string
	err = pool.QueryRow(ctx, "SELECT age_seconds, error_class FROM cdc.heartbeat WHERE source_id = $1", string(source)).Scan(&ageSeconds, &errorClass)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		fmt.Println("  Betriebsstatus (LH-FA-ADM-002): kein Lebenszeichen — Instanz hat noch nie geschlagen")
		fmt.Println("  Fehlerzustand (LH-FA-ADM-003): unbekannt (kein Lebenszeichen)")
	case err != nil:
		fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: cdc.heartbeat nicht lesbar (Schema-Rollout gelaufen?): %v\n", err)
		return 1
	default:
		fmt.Printf("  Betriebsstatus (LH-FA-ADM-002): Lebenszeichen vor %.3fs\n", ageSeconds)
		if errorClass == nil {
			fmt.Println("  Fehlerzustand (LH-FA-ADM-003): keiner (Normalbetrieb)")
		} else {
			fmt.Printf("  Fehlerzustand (LH-FA-ADM-003): %s\n", *errorClass)
		}
	}

	var captureLag float64
	if err := pool.QueryRow(ctx, "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_capture_lag'").Scan(&captureLag); err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: cdc.metrics (cdc_capture_lag) nicht lesbar: %v\n", err)
		return 1
	}
	fmt.Printf("  CDC-Abstand cdc_capture_lag (LH-FA-ADM-004): %.3fs\n", captureLag)

	rows, err := pool.Query(ctx, "SELECT label, value FROM cdc.metrics WHERE metric_name = 'cdc_consumer_lag' ORDER BY label")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: cdc.metrics (cdc_consumer_lag) nicht lesbar: %v\n", err)
		return 1
	}
	defer rows.Close()
	// Nur Consumer mit mindestens einer bestätigten Position tragen eine
	// Zeile (cdc.metrics-Definition, nacharbeit-observability.sql): ein
	// frisch registrierter, noch nie bestätigender Consumer erscheint
	// hier nicht — der Text unten benennt das ausdrücklich, statt sein
	// Fehlen als "kein Rückstand" lesbar zu lassen. `value` selbst ist für
	// eine solche Zeile trotzdem NULL, wenn die gebundene Quelle noch nie
	// eine Transaktion trug (`latest_commit_position`-Unterabfrage liefert
	// dann NULL, `WHERE cs.acknowledged_position IS NOT NULL` filtert das
	// nicht heraus) — der Scan liest deshalb über einen Zeiger, statt auf
	// diesen Fall mit einem Lesefehler zu enden.
	fmt.Println("  Verarbeitungsrückstand cdc_consumer_lag je Consumer (LH-FA-ADM-005, nur Consumer mit mindestens einer bestätigten Position):")
	found := false
	for rows.Next() {
		var consumer string
		var lag *float64
		if err := rows.Scan(&consumer, &lag); err != nil {
			fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: cdc.metrics (cdc_consumer_lag) nicht lesbar: %v\n", err)
			return 1
		}
		if lag == nil {
			fmt.Printf("    %s: unbekannt (Quelle trug noch nie eine Transaktion)\n", consumer)
		} else {
			fmt.Printf("    %s: %.0f\n", consumer, *lag)
		}
		found = true
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: cdc.metrics (cdc_consumer_lag) nicht lesbar: %v\n", err)
		return 1
	}
	if !found {
		fmt.Println("    (keiner — kein Consumer mit bestätigter Position)")
	}

	var blockerConsumer, blockerName string
	var blockerAckPos int64
	var blockerBacklog *int64
	err = pool.QueryRow(ctx,
		"SELECT consumer_id, name, acknowledged_position, backlog FROM cdc.retention_blockers WHERE source_id = $1",
		string(source),
	).Scan(&blockerConsumer, &blockerName, &blockerAckPos, &blockerBacklog)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		fmt.Println("  Blockierender Consumer (LH-FA-RET-005): kein Blocker (kein Consumer hat je gegen diese Quelle bestätigt)")
	case err != nil:
		fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: cdc.retention_blockers nicht lesbar: %v\n", err)
		return 1
	default:
		if blockerBacklog == nil {
			fmt.Printf("  Blockierender Consumer (LH-FA-RET-005): %s (%s), bestätigte Position %d, Rückstand unbekannt (Quelle trug noch nie eine Transaktion)\n", blockerName, blockerConsumer, blockerAckPos)
		} else {
			fmt.Printf("  Blockierender Consumer (LH-FA-RET-005): %s (%s), bestätigte Position %d, Rückstand %d\n", blockerName, blockerConsumer, blockerAckPos, *blockerBacklog)
		}
	}

	var storageBytes float64
	if err := pool.QueryRow(ctx, "SELECT value FROM cdc.metrics WHERE metric_name = 'cdc_storage_bytes'").Scan(&storageBytes); err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: cdc.metrics (cdc_storage_bytes) nicht lesbar: %v\n", err)
		return 1
	}
	fmt.Printf("  Speicherverbrauch cdc_storage_bytes (LH-FA-RET-006): %.0f Bytes\n", storageBytes)

	if err := diagnoseBackfillStatus(ctx, pool, source); err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: diagnose: cdc.backfill_status nicht lesbar: %v\n", err)
		return 1
	}

	return 0
}
