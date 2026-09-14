// Package http trägt den HTTP/JSON-Driving-Adapter (`ADR-0057`): Routing,
// Server-Lebenszyklus und die Zusammensetzung aus Token-Middleware und
// Handlern. Der Adapter importiert ausschließlich Inbound Ports
// (`internal/application/port/inbound`) und Domain-Typen zur
// Request-/Response-Übersetzung — keine Driven-Adapter-Interna, keine
// Application-Interna (`ADR-0057` Teilfrage 4, bereits maschinell über
// `a-check` durchgesetzt: der Glob `adapters: ["internal/adapters/**"]`
// erfasst dieses Verzeichnis bereits, keine `.a-check.yml`-Änderung
// nötig).
package http

import (
	"context"
	"net/http"

	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// Config trägt die Verdrahtungs-Eingabe des Adapters (`ADR-0057`
// Folgepflicht): die Server-Adresse, die beiden Token-Klassen
// (`ADR-0057` Teilfrage 3) und die neun Port-gedeckten Use Cases, die
// dieser Adapter über die API erreichbar macht (`ADR-0057` §Umfang der
// ersten API-Version) — Changes-Lesen und Diagnose/Health bleiben
// außerhalb, siehe `ADR-0057` §Konsequenzen/Folgepflicht.
type Config struct {
	// Addr trägt die Horch-Adresse (`CDC_HTTP_ADDR`); die Composition
	// Root entscheidet über den Start, dieser Typ trägt nur die Adresse.
	Addr string
	// TokenReader und TokenAdmin tragen die beiden Rechtsklassen
	// (`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`); ein leerer Wert
	// deaktiviert die jeweilige Klasse, statt ihn als gültiges Token zu
	// behandeln (`middleware.go`, `classifyToken`).
	TokenReader string
	TokenAdmin  string
	// RegisterConsumer trägt die Registrierung eines Consumers
	// (`LH-FA-CON-001`, `ADR-0028`).
	RegisterConsumer inbound.RegisterConsumerUseCase
	// AcknowledgeConsumer, GetConsumerPosition und RemoveConsumer tragen
	// die übrigen Consumer-Fähigkeiten: Bestätigung (`LH-FA-CON-004`),
	// Positions-Lese (`LH-FA-CON-003`/`-005`) und administrative Entfernung
	// (`LH-FA-CON-006`).
	AcknowledgeConsumer inbound.AcknowledgeConsumerUseCase
	GetConsumerPosition inbound.GetConsumerPositionUseCase
	RemoveConsumer      inbound.RemoveConsumerUseCase
	// EnableTable, DisableTable, GetStatus und ListTables tragen die
	// Verwaltungs-Fähigkeiten (`LH-FA-CFG-001`…`004`).
	EnableTable  inbound.EnableTableUseCase
	DisableTable inbound.DisableTableUseCase
	GetStatus    inbound.GetStatusUseCase
	ListTables   inbound.ListTablesUseCase
	// RunRetention trägt den Retention-Lauf (`LH-FA-RET-002`…`004`).
	RunRetention inbound.RunRetentionUseCase
	// Subscriber trägt den `Broadcaster`, von dem der SSE-Stream-Endpunkt
	// seine Changes liest (`LH-FA-SST-008`, `ADR-0061` Teilfrage 1/2).
	// Ohne ihn antwortet `GET /changes/stream` mit `503` (`sse.go`).
	Subscriber changeSubscriber
	// Log trägt den Telemetrie-Port (`ADR-0024`); ein nicht gesetzter
	// Wert fällt auf `outbound.NoopLog` zurück.
	Log outbound.LogPort
}

// Server trägt den `net/http`-Server und dessen Lebenszyklus (`ADR-0057`
// Teilfrage 1: Go-Standardbibliothek genügt, keine neue
// Fremd-Abhängigkeit).
type Server struct {
	httpServer *http.Server
	log        outbound.LogPort
}

// New verdrahtet Routing und Token-Middleware; der zurückgegebene Server
// ist noch nicht gestartet (`Start`).
func New(cfg Config) *Server {
	log := cfg.Log
	if log == nil {
		log = outbound.NoopLog
	}
	mux := http.NewServeMux()
	mux.Handle("POST /consumers", withToken(cfg.TokenReader, cfg.TokenAdmin, roleAdmin,
		registerConsumerHandler(cfg.RegisterConsumer, log)))
	mux.Handle("POST /consumers/acknowledge", withToken(cfg.TokenReader, cfg.TokenAdmin, roleAdmin,
		acknowledgeConsumerHandler(cfg.AcknowledgeConsumer, log)))
	mux.Handle("GET /consumers/position", withToken(cfg.TokenReader, cfg.TokenAdmin, roleReader,
		getConsumerPositionHandler(cfg.GetConsumerPosition, log)))
	mux.Handle("POST /consumers/remove", withToken(cfg.TokenReader, cfg.TokenAdmin, roleAdmin,
		removeConsumerHandler(cfg.RemoveConsumer, log)))
	mux.Handle("POST /tables/enable", withToken(cfg.TokenReader, cfg.TokenAdmin, roleAdmin,
		enableTableHandler(cfg.EnableTable, log)))
	mux.Handle("POST /tables/disable", withToken(cfg.TokenReader, cfg.TokenAdmin, roleAdmin,
		disableTableHandler(cfg.DisableTable, log)))
	mux.Handle("GET /tables/status", withToken(cfg.TokenReader, cfg.TokenAdmin, roleReader,
		getStatusHandler(cfg.GetStatus, log)))
	mux.Handle("GET /tables", withToken(cfg.TokenReader, cfg.TokenAdmin, roleReader,
		listTablesHandler(cfg.ListTables, log)))
	mux.Handle("POST /retention/run", withToken(cfg.TokenReader, cfg.TokenAdmin, roleAdmin,
		runRetentionHandler(cfg.RunRetention, log)))
	mux.Handle("GET /changes/stream", withToken(cfg.TokenReader, cfg.TokenAdmin, roleReader,
		streamChangesHandler(cfg.Subscriber, log)))
	return &Server{
		httpServer: &http.Server{Addr: cfg.Addr, Handler: mux},
		log:        log,
	}
}

// Handler trägt den verdrahteten `http.Handler` ohne Bindung an einen
// realen Port — die Grundlage für `httptest`-Whitebox-Tests gegen den
// Adapter, ohne einen realen Socket zu öffnen.
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}

// Start blockiert, bis der Server über `Shutdown` beendet wird oder ein
// Bind-/Laufzeitfehler auftritt; `http.ErrServerClosed` ist dabei der
// reguläre Ausgang eines geordneten `Shutdown`, keine Fehlerklasse
// (`SPEC-008`). Additiv zur bestehenden Verdrahtung
// (`internal/bootstrap`): ein Aufrufer, der `CDC_HTTP_ADDR` nicht setzt,
// ruft `Start` nicht auf — das No-Op bei fehlender Adresse bleibt Sache
// der Composition Root, nicht dieses Adapters.
func (s *Server) Start() error {
	s.log.Info(context.Background(), "http: Adapter gestartet", "addr", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown beendet den Server geordnet (`http.Server.Shutdown`): laufende
// Aufrufe laufen zu Ende, neue Verbindungen werden nicht mehr
// angenommen.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
