package grpc

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
