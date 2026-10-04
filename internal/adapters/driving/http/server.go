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

	"github.com/pt9912/pg-change-feed/internal/application/port/apiauth"
	"github.com/pt9912/pg-change-feed/internal/application/port/inbound"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// Config trägt die Verdrahtungs-Eingabe des Adapters (`ADR-0057`
// Folgepflicht): die Server-Adresse, die beiden Token-Klassen
// (Teilfrage 3) und die Port-gedeckten Use Cases, die dieser
// Adapter über die API erreichbar macht — die neun aus `ADR-0057` §Umfang
// der ersten API-Version und das Changes-Lesen. Diagnose/Health
// bleibt außerhalb.
type Config struct {
	// Addr trägt die Horch-Adresse (`CDC_HTTP_ADDR`); die Composition
	// Root entscheidet über den Start, dieser Typ trägt nur die Adresse.
	Addr string
	// TokenReader und TokenAdmin tragen das Singular-Token je Rechtsklasse
	// (`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`), TokensReader und
	// TokensAdmin die Liste je Klasse (`CDC_API_TOKENS_READER`/
	// `CDC_API_TOKENS_ADMIN`); gültig ist die Vereinigung beider. Ein leerer
	// Wert ist kein gültiges Token (`apiauth.Classifier`).
	TokenReader  string
	TokenAdmin   string
	TokensReader []string
	TokensAdmin  []string
	// RegisterConsumer trägt die Registrierung eines Consumers
	// (`LH-FA-CON-001`).
	RegisterConsumer inbound.RegisterConsumerUseCase
	// AcknowledgeConsumer, GetConsumerPosition und RemoveConsumer tragen
	// die übrigen Consumer-Fähigkeiten: Bestätigung (`LH-FA-CON-004`),
	// Positions-Lese und administrative Entfernung.
	AcknowledgeConsumer inbound.AcknowledgeConsumerUseCase
	GetConsumerPosition inbound.GetConsumerPositionUseCase
	RemoveConsumer      inbound.RemoveConsumerUseCase
	// EnableTable, DisableTable, GetStatus und ListTables tragen die
	// Verwaltungs-Fähigkeiten (`LH-FA-CFG-001`).
	EnableTable  inbound.EnableTableUseCase
	DisableTable inbound.DisableTableUseCase
	GetStatus    inbound.GetStatusUseCase
	ListTables   inbound.ListTablesUseCase
	// RunRetention trägt den Retention-Lauf (`LH-FA-RET-002`).
	RunRetention inbound.RunRetentionUseCase
	// ReadChanges trägt das Changes-Lesen über die API
	// (`LH-FA-SST-006`) — eine dünne
	// Fassade über demselben `ChangeStorePort`, den der View-Direktzugriff
	// trägt.
	ReadChanges inbound.ReadChangesUseCase
	// Diagnose trägt das Diagnose-Lesen über die API (`ADR-0132`) — derselbe
	// Use Case, den auch der CLI-Sondermodus und der gRPC-Handler aufrufen.
	Diagnose inbound.DiagnoseUseCase
	// Subscriber trägt den `Broadcaster`, von dem der SSE-Stream-Endpunkt
	// seine Changes liest (`ADR-0061` Teilfrage 1/2).
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
	tokens := apiauth.FromConfig(cfg.TokenReader, cfg.TokensReader, cfg.TokenAdmin, cfg.TokensAdmin)
	mux := http.NewServeMux()
	mux.Handle("POST /consumers", withToken(tokens, apiauth.Admin,
		registerConsumerHandler(cfg.RegisterConsumer, log)))
	mux.Handle("POST /consumers/acknowledge", withToken(tokens, apiauth.Admin,
		acknowledgeConsumerHandler(cfg.AcknowledgeConsumer, log)))
	mux.Handle("GET /consumers/position", withToken(tokens, apiauth.Reader,
		getConsumerPositionHandler(cfg.GetConsumerPosition, log)))
	mux.Handle("POST /consumers/remove", withToken(tokens, apiauth.Admin,
		removeConsumerHandler(cfg.RemoveConsumer, log)))
	mux.Handle("POST /tables/enable", withToken(tokens, apiauth.Admin,
		enableTableHandler(cfg.EnableTable, log)))
	mux.Handle("POST /tables/disable", withToken(tokens, apiauth.Admin,
		disableTableHandler(cfg.DisableTable, log)))
	mux.Handle("GET /tables/status", withToken(tokens, apiauth.Reader,
		getStatusHandler(cfg.GetStatus, log)))
	mux.Handle("GET /tables", withToken(tokens, apiauth.Reader,
		listTablesHandler(cfg.ListTables, log)))
	mux.Handle("POST /retention/run", withToken(tokens, apiauth.Admin,
		runRetentionHandler(cfg.RunRetention, log)))
	mux.Handle("GET /changes", withToken(tokens, apiauth.Reader,
		readChangesHandler(cfg.ReadChanges, log)))
	mux.Handle("GET /changes/stream", withToken(tokens, apiauth.Reader,
		streamChangesHandler(cfg.Subscriber, log)))
	mux.Handle("GET /diagnose", withToken(tokens, apiauth.Reader,
		diagnoseHandler(cfg.Diagnose, log)))
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
