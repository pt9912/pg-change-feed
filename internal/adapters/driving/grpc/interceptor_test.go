package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/pt9912/pg-change-feed/internal/adapters/driving/grpc/streamv1"
)

// TestClassifyTokenLeereKonfigurationTrifftKeinToken trägt die Grenze der
// Token-Konfiguration (`ADR-0060` Teilfrage 4): ein leer konfiguriertes
// Token darf kein Aufruf-Token treffen — sonst würde ein fehlendes
// `authorization` (leerer Token) gegen eine ebenfalls ungesetzte
// Token-Klasse eine dritte, implizite Rechtsklasse eröffnen.
func TestClassifyTokenLeereKonfigurationTrifftKeinToken(t *testing.T) {
	if got := classifyToken("", "", ""); got != roleNone {
		t.Fatalf("classifyToken(\"\", \"\", \"\") = %v (Erwartung: roleNone)", got)
	}
	if got := classifyToken("beliebig", "", ""); got != roleNone {
		t.Fatalf("classifyToken(\"beliebig\", \"\", \"\") = %v (Erwartung: roleNone)", got)
	}
}

// TestStreamChangesLeereTokenKonfigurationEndetMitUnauthenticated trägt
// dieselbe Grenze am laufenden Adapter: ohne konfigurierte Token
// (`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN` ungesetzt) öffnet kein
// Aufruf-Token den Stream — der Interceptor ist fail-closed, nicht still
// durchlässig (`LH-FA-SST-008` Negative).
func TestStreamChangesLeereTokenKonfigurationEndetMitUnauthenticated(t *testing.T) {
	client := startTestServerMitTokenKonfiguration(t, newFakeSubscriber(), "", "")
	for _, authorization := range []string{"", bearerPrefix, bearerPrefix + "beliebig"} {
		stream := streamMitToken(t, client, authorization)
		_, err := recvChange(t, stream)
		if status.Code(err) != codes.Unauthenticated {
			t.Fatalf("authorization %q: Status %v (Erwartung: %v)", authorization, status.Code(err), codes.Unauthenticated)
		}
	}
}

// TestAuthStreamInterceptorOhneEingehendeMetadataEndetMitUnauthenticated
// trägt den Interceptor an der Eingabeseite seines Kontexts: ein Kontext
// **ohne** eingehende Metadata trägt keinen Token, der Handler wird nicht
// erreicht. Über eine reale gRPC-Verbindung setzt der Server die eingehende
// Metadata immer, dieser Zweig ist dort nicht erreichbar.
//
// Die Ablehnung hängt am übergebenen Kontext: derselbe Interceptor reicht
// einen Aufruf mit `Bearer <reader-token>` an den Handler durch — der
// Negativfall hängt damit an der Eingabe, nicht am Interceptor allein
// (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`).
//
// Rot färbende Mutation: in `credentialToken` an der Stelle ohne eingehende
// Metadata einen Wert statt `""` zurückgeben — dann wird der Handler erreicht.
// Das bloße Entfernen des `!ok`-Zweigs färbt diesen Test **nicht** rot: der
// Nullwert der Metadata trägt die leere Kennung auf demselben Weg.
func TestAuthStreamInterceptorOhneEingehendeMetadataEndetMitUnauthenticated(t *testing.T) {
	interceptor := authStreamInterceptor(testReaderToken, testAdminToken)
	info := &grpc.StreamServerInfo{FullMethod: streamv1.ChangeStream_StreamChanges_FullMethodName}

	gerufen := false
	handler := func(any, grpc.ServerStream) error { gerufen = true; return nil }

	ohneMetadata := &fakeServerStream{ctx: context.Background()}
	err := interceptor(nil, ohneMetadata, info, handler)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("Status ohne eingehende Metadata: %v (Erwartung: %v)", status.Code(err), codes.Unauthenticated)
	}
	if gerufen {
		t.Fatal("der Handler wurde ohne Token erreicht")
	}

	mitToken := &fakeServerStream{ctx: metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(authorizationMetadataKey, bearerPrefix+testReaderToken))}
	if err := interceptor(nil, mitToken, info, handler); err != nil {
		t.Fatalf("Interceptor mit gültigem Token: %v (Erwartung: kein Fehler)", err)
	}
	if !gerufen {
		t.Fatal("der Handler wurde mit gültigem Token nicht erreicht")
	}
}
