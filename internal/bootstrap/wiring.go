// Package bootstrap ist die Composition Root (`ADR-0026`): er kennt die
// konkreten Adapter und verdrahtet die Pipeline an genau einer Stelle —
// ChangeStore-Driven-Adapter, Aktivierungs-Driven-Adapter mit dem
// EnableTable Use Case (`ADR-0028`), Heartbeat-Driven-Adapter mit dem
// periodischen Timer-Zug (`ADR-0024`), Replication-Stream-
// Driving-Adapter, Capture Service und Replication-ACK-Driven-Adapter.
// Die Abhängigkeitsregel (§2 der Architektur-Sicht) bleibt hier lokal
// einhaltbar; `main` referenziert keinen Adapter-Konstruktor.
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
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresack"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/telemetry"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/decode"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/mapper"
	"github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/acknowledge"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/disable"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/enable"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/register"
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
	return cfg, nil
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
func activatedTableBindings(ctx context.Context, activation outbound.TableActivationPort, schemaStore outbound.SchemaStorePort, source model.SourceID) (map[string]mapper.TableBinding, error) {
	registered, err := activation.List(ctx, source)
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
		tables[table.QualifiedName()] = mapper.TableBinding{TableID: table.ID, SchemaVersion: current.ID}
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
	enableTables := enable.NewEnableTableService(activation)
	disableTables := disable.NewDisableTableService(activation)
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
	// Datenbank fort (`TableActivationPort.List`/`SchemaStorePort.CurrentVersion`),
	// nicht ausschließlich aus `cfg.Tables`: `CDC_TABLES` bleibt der
	// Erstaktivierungs-Seed oben (schreibt eine leere Datenbank fort), die
	// Grundlage des Stream-Starts ist der committed Stand — ein
	// Prozess-Neustart verliert damit keine zwischenzeitlich per SQL
	// aktivierte Tabelle, deren Kennung `CDC_TABLES` nicht trägt.
	assemblerTables, err := activatedTableBindings(ctx, activation, schemaStore, cfg.Source)
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
	if err := stream.BindCapture(capture.NewCaptureService(store, ack)); err != nil {
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
			requests:      adminRequests,
			listener:      adminListener,
			activation:    activation,
			enableTables:  enableTables,
			disableTables: disableTables,
			schemaStore:   schemaStore,
			assembler:     stream.Assembler(),
			publication:   cfg.Publication,
			pollInterval:  administrationPollInterval,
			log:           log,
		})
	}()

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
// testbar ist (`slice-026` Fixrunde, Review F-1).
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
	requests      outbound.AdministrationRequestPort
	listener      administrationListener
	activation    outbound.TableActivationPort
	enableTables  inbound.EnableTableUseCase
	disableTables inbound.DisableTableUseCase
	schemaStore   outbound.SchemaStorePort
	assembler     *mapper.Assembler
	publication   string
	pollInterval  time.Duration
	log           outbound.LogPort
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
// ändert.
func processAdministrationRequests(ctx context.Context, deps administrationDeps) {
	pending, err := deps.requests.ListPending(ctx)
	if err != nil {
		deps.log.Warn(ctx, "administration: Anträge lesen fehlgeschlagen", "error", err)
		return
	}
	for _, request := range pending {
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
		deps.assembler.AddBinding(qualified, mapper.TableBinding{TableID: registered.ID, SchemaVersion: current.ID})
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
	default:
		return fmt.Errorf("Antragsart %q trägt nicht die geschlossene Menge enable/disable", request.Kind)
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
		errors.Is(err, mapper.ErrIncompatibleSchemaChange):
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
