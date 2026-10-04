package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/pt9912/pg-change-feed/internal/application/port/apiauth"
)

// authorizationMetadataKey trägt den gRPC-Metadata-Schlüssel der
// Authentifizierung (`ADR-0060` Teilfrage 4): der Interceptor vergleicht
// seinen Wert gegen dieselben beiden Token-Klassen, die auch die
// HTTP-Token-Middleware schützt (`CDC_API_TOKEN_READER`/
// `CDC_API_TOKEN_ADMIN`).
const authorizationMetadataKey = "authorization"

// bearerPrefix trägt die Wertform des Metadata-Eintrags: derselbe
// `Bearer `-Vorsprung wie der `Authorization`-Header der HTTP-API
// — beide Netzwerk-Zugriffswege tragen dieselbe
// Authentifizierungsform, die Token-Klassen bleiben dieselben
// (`ADR-0060` Teilfrage 4).
const bearerPrefix = "Bearer "

// credentialToken liest den Token aus dem `authorization`-Metadata-Eintrag
// des Stream-Kontexts; ein fehlender Eintrag oder eine falsche Wertform
// trägt einen leeren Token zurück, den `apiauth.Classifier.Classify` wie
// einen fehlenden behandelt.
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
// `Unauthenticated` — sichtbar, nicht still mit leeren Daten fortgesetzt.
// `ChangeStream` trägt ausschließlich
// Streaming-RPCs, deshalb trägt nur der Stream-Interceptor eine Prüfung;
// ein Unary-Interceptor hätte hier keinen Aufruf zu schützen.
func authStreamInterceptor(tokens *apiauth.Classifier) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if tokens.Classify(credentialToken(stream.Context())) == apiauth.None {
			return status.Error(codes.Unauthenticated, "fehlender oder unbekannter authorization-Metadata-Wert")
		}
		return handler(srv, stream)
	}
}

// administrationRPCRoles trägt die Rechtsklassen-Tabelle der elf
// `Administration`-RPCs (`ADR-0132` Teilfrage 6): dieselbe Rollen-Zuordnung
// wie `withToken` in `internal/adapters/driving/http/middleware.go`.
var administrationRPCRoles = map[string]apiauth.Role{
	"RegisterConsumer":    apiauth.Admin,
	"AcknowledgeConsumer": apiauth.Admin,
	"RemoveConsumer":      apiauth.Admin,
	"EnableTable":         apiauth.Admin,
	"DisableTable":        apiauth.Admin,
	"RunRetention":        apiauth.Admin,
	"GetConsumerPosition": apiauth.Reader,
	"GetTableStatus":      apiauth.Reader,
	"ListTables":          apiauth.Reader,
	"ReadChanges":         apiauth.Reader,
	"Diagnose":            apiauth.Reader,
}

// methodName liest das letzte Pfadsegment aus `info.FullMethod`
// (Form `/cdc.administration.v1.Administration/<Methode>`, `ADR-0130`
// Teilfrage 4).
func methodName(fullMethod string) string {
	if idx := strings.LastIndex(fullMethod, "/"); idx >= 0 {
		return fullMethod[idx+1:]
	}
	return fullMethod
}

// authUnaryInterceptor schützt alle unären RPCs des `Administration`-Service
// (`ADR-0130` Teilfrage 4, Fitness Function): ein Aufruf ohne
// `authorization`-Metadata oder mit einem Wert, der keiner der beiden
// konfigurierten Klassen entspricht, endet mit `Unauthenticated`; ein
// bekanntes Token unterhalb der für die Methode verlangten Rechtsklasse
// endet mit `PermissionDenied`. Ein Methodenname ohne Eintrag in `rights`
// fällt fail-closed auf `apiauth.Admin` — am fest verdrahteten RPC-Satz sollte
// das nicht vorkommen.
func authUnaryInterceptor(tokens *apiauth.Classifier, rights map[string]apiauth.Role) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		required, ok := rights[methodName(info.FullMethod)]
		if !ok {
			required = apiauth.Admin
		}
		caller := tokens.Classify(credentialToken(ctx))
		if caller == apiauth.None {
			return nil, status.Error(codes.Unauthenticated, "fehlender oder unbekannter authorization-Metadata-Wert")
		}
		if caller < required {
			return nil, status.Error(codes.PermissionDenied, "Rechtsklasse unzureichend für diese RPC")
		}
		return handler(ctx, req)
	}
}
