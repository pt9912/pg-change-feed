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
	"github.com/pt9912/pg-change-feed/internal/application/usecase/capture"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/enable"
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
	// LogLevel trägt den Level des JSON-Handlers, den der
	// Telemetrie-Driven-Adapter baut (`telemetry.New`, `ADR-0024`);
	// Herkunft ist `envLogLevel`/`parseLogLevel`, mit Default
	// `slog.LevelInfo`.
	LogLevel slog.Level
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

	store, err := postgresstorage.New(ctx, cfg.DSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer store.Close()

	// Die Aktivierung läuft als Use Case (`ADR-0028`, `LH-FA-CFG-001`):
	// die Bindungen der Konfiguration laufen vor dem Stream-Start als
	// EnableTable-Aufrufe — die Publication ist Start-Vorbedingung des
	// Stream-Adapters (`LH-FA-CFG-001.a`) und die Bindungs-Zeilen tragen
	// die Fremdschlüssel der ersten Persistenz (`SPEC-001`). Der Aufruf
	// ist idempotent (`LH-FA-CFG-001` Boundary) und trägt den Stand auch
	// nach einem Container-Neustart nach.
	// Die Verdrahtung trägt vier Verbindungen gegen dieselbe Instanz —
	// Store-Pool, Aktivierungs-Pool, Heartbeat-Pool, Stream-Verbindung;
	// das MVP hält die Adapter-Lebenszyklen getrennt, statt einen Pool
	// über die Adapter zu teilen. Die Instanz trägt Quelle und
	// CDC-Speicher gleichermaßen (Abschnitt 1 Lastenheft); ein geteilter
	// Pool ist keine Wirkung dieses Verdrahtungsstands.
	activation, err := postgresstorage.NewTableActivation(ctx, cfg.DSN, postgresstorage.WithLog(log))
	if err != nil {
		return err
	}
	defer activation.Close()
	// Der Heartbeat-Pool trägt ausschließlich den periodischen
	// Lebenszeichen-Zug (runHeartbeat, unten) — eine eigene Verbindung,
	// getrennt von Store- und Aktivierungs-Pool: der Timer-Zug teilt
	// keine Verbindung und keine Goroutine mit der
	// Capture-Persist-ACK-Schleife (`LH-QA-REL-001.a`).
	heartbeat, err := postgresstorage.NewHeartbeat(ctx, cfg.DSN, postgresstorage.WithLog(log))
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

	stream, err := receive.NewStream(ctx, receive.Config{
		DSN:         cfg.DSN,
		Source:      cfg.Source,
		Publication: cfg.Publication,
		Slot:        cfg.Slot,
		Tables:      cfg.Tables,
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

	streamErr := stream.Run(ctx)
	stopHeartbeat()
	heartbeatDone.Wait()
	return streamErr
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
		errors.Is(err, mapper.ErrTruncateUnsupported):
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
// (`Run` oben).
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
