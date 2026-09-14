package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// authorizationMetadataKey trägt den gRPC-Metadata-Schlüssel der
// Authentifizierung (`ADR-0060` Teilfrage 4): der Interceptor vergleicht
// seinen Wert gegen dieselben beiden Token-Klassen, die auch die
// HTTP-Token-Middleware schützt (`SPEC-018`, `CDC_API_TOKEN_READER`/
// `CDC_API_TOKEN_ADMIN`).
const authorizationMetadataKey = "authorization"

// bearerPrefix trägt die Wertform des Metadata-Eintrags: derselbe
// `Bearer `-Vorsprung wie der `Authorization`-Header der HTTP-API
// (`SPEC-018`) — beide Netzwerk-Zugriffswege tragen dieselbe
// Authentifizierungsform, die Token-Klassen bleiben dieselben
// (`ADR-0060` Teilfrage 4).
const bearerPrefix = "Bearer "

// role trägt die zwei Rechtsklassen des Stream-Zugriffs (`ADR-0060`
// Teilfrage 4): wie in der HTTP-Middleware deckt `roleAdmin` implizit
// `roleReader` ab. `roleNone` trägt den fehlenden und den unbekannten Wert
// gemeinsam — beide enden über denselben `Unauthenticated`-Pfad.
type role int

const (
	roleNone role = iota
	roleReader
	roleAdmin
)

// classifyToken ordnet einen Token einer Rechtsklasse zu. Ein leer
// konfiguriertes Token (`readerToken`/`adminToken` ungesetzt) trifft nie ein
// Aufruf-Token — ein ungesetztes Token würde sonst eine dritte, implizite
// Rechtsklasse eröffnen (dieselbe Grenze wie in der HTTP-Token-Middleware,
// `ADR-0057` Teilfrage 3).
//
// Diese Funktion steht als zweite, wortgleiche Fassung in
// `internal/adapters/driving/http/middleware.go` (`role`, die drei
// Konstanten und `classifyToken`). Das `.a-check.yml`-Schichtenmodell führt
// keine `adapters→adapters`-Kante, deshalb trägt jeder Driving-Adapter seine
// eigene Fassung derselben Zuordnung; beide Fassungen sind zusammen zu
// ändern.
func classifyToken(token, readerToken, adminToken string) role {
	if token == "" {
		return roleNone
	}
	if adminToken != "" && token == adminToken {
		return roleAdmin
	}
	if readerToken != "" && token == readerToken {
		return roleReader
	}
	return roleNone
}

// credentialToken liest den Token aus dem `authorization`-Metadata-Eintrag
// des Stream-Kontexts; ein fehlender Eintrag oder eine falsche Wertform
// trägt einen leeren Token zurück, den `classifyToken` wie einen fehlenden
// behandelt.
func credentialToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(authorizationMetadataKey)
	if len(values) == 0 {
		return ""
	}
	if !strings.HasPrefix(values[0], bearerPrefix) {
		return ""
	}
	return strings.TrimPrefix(values[0], bearerPrefix)
}

// authStreamInterceptor schützt alle Streaming-RPCs des Servers
// (`ADR-0060` Teilfrage 4, Fitness Function): ein Öffnungsversuch ohne
// `authorization`-Metadata oder mit einem Wert, der keiner der beiden
// konfigurierten Klassen entspricht, endet mit dem gRPC-Status
// `Unauthenticated` — sichtbar, nicht still mit leeren Daten fortgesetzt
// (`LH-FA-SST-008` Negative). `ChangeStream` trägt ausschließlich
// Streaming-RPCs, deshalb trägt nur der Stream-Interceptor eine Prüfung;
// ein Unary-Interceptor hätte hier keinen Aufruf zu schützen.
func authStreamInterceptor(readerToken, adminToken string) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if classifyToken(credentialToken(stream.Context()), readerToken, adminToken) == roleNone {
			return status.Error(codes.Unauthenticated, "fehlender oder unbekannter authorization-Metadata-Wert")
		}
		return handler(srv, stream)
	}
}
