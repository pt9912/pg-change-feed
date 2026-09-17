// Package grpc trägt den gRPC-Driving-Adapter (`ADR-0060`): Server,
// Auth-Interceptor und die Übersetzung Domain-Change ↔ Protobuf-Nachricht
// (`SPEC-020`). Der Adapter kennt den `Broadcaster` nicht über einen Import
// des Driven-Pakets, sondern über ein hier deklariertes, strukturell
// erfülltes Interface (`changeSubscriber`, `ADR-0060` Teilfrage 2) — kein
// Adapter-Paket importiert ein anderes, nur `internal/bootstrap` kennt
// beide konkret (Richtungskonvention aus `spec/architecture.md` §1).
package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// changeSubscriber trägt die eine Fähigkeit, die dieser Adapter vom
// `Broadcaster` braucht (`ADR-0060` Teilfrage 2): registrieren und einen
// Kanal samt Abmelde-Funktion erhalten. `*grpcstream.Broadcaster` erfüllt
// dieses Interface strukturell; die Composition Root verdrahtet es.
type changeSubscriber interface {
	Subscribe() (<-chan *model.Change, func())
}

// Config trägt die Verdrahtungs-Eingabe des Adapters (`ADR-0060`
// Folgepflicht): die Horch-Adresse, die beiden Token-Klassen aus
// `ADR-0057` Teilfrage 3 und den Broadcaster.
type Config struct {
	// Addr trägt die Horch-Adresse (`CDC_GRPC_ADDR`); die Composition Root
	// entscheidet über den Start, dieser Typ trägt nur die Adresse.
	Addr string
	// TokenReader und TokenAdmin tragen die beiden Rechtsklassen
	// (`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`); ein leerer Wert
	// deaktiviert die jeweilige Klasse, statt ihn als gültiges Token zu
	// behandeln (`interceptor.go`, `classifyToken`).
	TokenReader string
	TokenAdmin  string
	// Subscriber trägt den Broadcaster, von dem der Stream seine Changes
	// liest. Ohne ihn öffnet kein Stream — der RPC endet sichtbar statt
	// still leer (`StreamChanges`).
	Subscriber changeSubscriber
	// Log trägt den Telemetrie-Port (`ADR-0024`); ein nicht gesetzter Wert
	// fällt auf `outbound.NoopLog` zurück.
	Log outbound.LogPort
}

// Server trägt den `grpc.Server` und dessen Lebenszyklus (`ADR-0060`).
type Server struct {
	grpcServer *grpc.Server
	addr       string
	log        outbound.LogPort
}

// New verdrahtet den Stream-Service und den Auth-Interceptor; der
// zurückgegebene Server ist noch nicht gestartet (`Start`).
func New(cfg Config) *Server {
	log := cfg.Log
	if log == nil {
		log = outbound.NoopLog
	}
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(authStreamInterceptor(cfg.TokenReader, cfg.TokenAdmin)))
	streamv1.RegisterChangeStreamServer(grpcServer, &changeStreamService{subscriber: cfg.Subscriber, log: log})
	return &Server{grpcServer: grpcServer, addr: cfg.Addr, log: log}
}

// Start bindet die konfigurierte Adresse und blockiert, bis der Server über
// `Shutdown` beendet wird oder ein Bind-/Laufzeitfehler auftritt
// (`grpc.ErrServerStopped` ist der reguläre Ausgang eines geordneten
// `Shutdown`). Additiv zur bestehenden Verdrahtung
// (`internal/bootstrap`): ein Aufrufer, der `CDC_GRPC_ADDR` nicht setzt,
// ruft `Start` nicht auf — das No-Op bei fehlender Adresse bleibt Sache der
// Composition Root, nicht dieses Adapters.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc: Horch-Adresse %q: %w", s.addr, err)
	}
	return s.serve(listener)
}

// serve trägt den Server-Lauf auf einem fertigen Listener: `Start` bindet
// darüber die konfigurierte Adresse, der Whitebox-Test einen
// In-Memory-Listener (`bufconn`).
func (s *Server) serve(listener net.Listener) error {
	s.log.Info(context.Background(), "grpc: Adapter gestartet", "addr", listener.Addr().String())
	if err := s.grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}
	return nil
}

// Shutdown beendet den Server. `grpc.Server.Stop` statt `GracefulStop`: die
// Streaming-Verbindungen dieses Adapters enden nicht von selbst, ein
// `GracefulStop` wartete deshalb auf ihr Ende statt zu beenden.
func (s *Server) Shutdown() {
	s.grpcServer.Stop()
}

// changeStreamService implementiert den generierten `ChangeStreamServer`
// (`SPEC-020`).
type changeStreamService struct {
	streamv1.UnimplementedChangeStreamServer
	subscriber changeSubscriber
	log        outbound.LogPort
}

// StreamChanges registriert einen Empfänger am Broadcaster und überträgt
// jeden eintreffenden Change, bis der Client die Verbindung beendet
// (`ADR-0060` Teilfrage 3): kein Replay — ein ab Verbindungsaufbau
// eintreffender Change geht über den Stream, ältere bleiben dem
// Lesezugriffsweg vorbehalten (`LH-FA-SST-008` Boundary).
func (s *changeStreamService) StreamChanges(_ *streamv1.StreamChangesRequest, stream grpc.ServerStreamingServer[streamv1.Change]) error {
	if s.subscriber == nil {
		return status.Error(codes.Internal, "ChangeStream ohne Broadcaster verdrahtet")
	}
	changes, cancel := s.subscriber.Subscribe()
	defer cancel()
	s.log.Info(stream.Context(), "grpc: Stream geöffnet")
	defer s.log.Info(context.Background(), "grpc: Stream beendet")
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case change := <-changes:
			if err := stream.Send(toProtoChange(change)); err != nil {
				return err
			}
		}
	}
}

// toProtoChange übersetzt einen Domain-Change in seine Protobuf-Form
// (`SPEC-020`): dieselben Felder, Row Images unverändert übernommen.
func toProtoChange(change *model.Change) *streamv1.Change {
	return &streamv1.Change{
		ChangeId:      string(change.ID),
		TransactionId: string(change.TransactionID),
		SourceTableId: string(change.SourceTableID),
		Sequence:      change.Sequence,
		Operation:     string(change.Operation),
		OldImage:      change.OldImage,
		NewImage:      change.NewImage,
		SchemaVersion: string(change.SchemaVersion),
		Schema:        change.Schema,
		Table:         change.Table,
	}
}
