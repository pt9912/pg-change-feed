package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
	"github.com/pt9912/pg-change-feed/internal/application/port/apiauth"
)

// testTokens trägt den Klassifikator der Tests: je ein Token je Klasse.
func testTokens() *apiauth.Classifier {
	return apiauth.FromConfig(testReaderToken, nil, testAdminToken, nil)
}

// TestAuthInterceptorMehrereTokenJeKlasse trägt die Zuordnung der
// Interceptoren für mehrere Token je Klasse (`LH-FA-SST-012`): beide Token
// einer Klasse werden gleich bedient, ein entferntes Token endet mit
// `Unauthenticated`, ein Wert in beiden Klassen erreicht die administrative RPC.
//
// Rot färbende Mutation: in `authUnaryInterceptor` den Rückgabewert von
// `tokens.Classify` durch `apiauth.Reader` ersetzen — die Fälle „entferntes
// Token“ und „administrative RPC“ färben rot.
func TestAuthInterceptorMehrereTokenJeKlasse(t *testing.T) {
	tokens := apiauth.New(
		[]string{"reader-1", "reader-2", "beide"},
		[]string{"admin-1", "admin-2", "beide"},
	)
	unary := authUnaryInterceptor(tokens, administrationRPCRoles)
	stream := authStreamInterceptor(tokens)
	streamInfo := &grpc.StreamServerInfo{FullMethod: streamv1.ChangeStream_StreamChanges_FullMethodName}
	const adminRPC = "/cdc.administration.v1.Administration/RegisterConsumer"
	const readerRPC = "/cdc.administration.v1.Administration/ListTables"

	cases := []struct {
		name       string
		token      string
		method     string
		wantStatus codes.Code
	}{
		{"erstes Reader-Token, lesende RPC", "reader-1", readerRPC, codes.OK},
		{"zweites Reader-Token, lesende RPC", "reader-2", readerRPC, codes.OK},
		{"zweites Reader-Token, administrative RPC", "reader-2", adminRPC, codes.PermissionDenied},
		{"erstes Admin-Token, administrative RPC", "admin-1", adminRPC, codes.OK},
		{"zweites Admin-Token, administrative RPC", "admin-2", adminRPC, codes.OK},
		{"Wert in beiden Klassen, administrative RPC", "beide", adminRPC, codes.OK},
		{"entferntes Token, lesende RPC", "reader-0", readerRPC, codes.Unauthenticated},
		{"entferntes Token, administrative RPC", "admin-0", adminRPC, codes.Unauthenticated},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gerufen, handler := unaryHandlerAufruf()
			_, err := unary(unaryCtxMitToken(bearerPrefix+tc.token), nil, &grpc.UnaryServerInfo{FullMethod: tc.method}, handler)
			if status.Code(err) != tc.wantStatus {
				t.Fatalf("Token %q, %s: Status %v, erwartet %v", tc.token, tc.method, status.Code(err), tc.wantStatus)
			}
			if *gerufen != (tc.wantStatus == codes.OK) {
				t.Fatalf("Token %q, %s: Handler erreicht = %v bei Status %v", tc.token, tc.method, *gerufen, tc.wantStatus)
			}
		})
	}

	for token, want := range map[string]codes.Code{
		"reader-2": codes.OK, "admin-2": codes.OK, "beide": codes.OK, "reader-0": codes.Unauthenticated,
	} {
		ctx := unaryCtxMitToken(bearerPrefix + token)
		err := stream(nil, &fakeServerStream{ctx: ctx}, streamInfo, func(any, grpc.ServerStream) error { return nil })
		if status.Code(err) != want {
			t.Fatalf("Stream mit Token %q: Status %v, erwartet %v", token, status.Code(err), want)
		}
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
// Metadata einen **gültigen** Token-Wert statt `""` zurückgeben
// (`return "reader-token"`) — dann ordnet `Classify` den Aufruf
// `apiauth.Reader` zu, der Handler wird erreicht und dieser Test färbt rot. Ein
// **unbekannter** Wert an derselben Stelle färbt ihn nicht rot:
// `Classify` ordnet ihn `apiauth.None` zu wie den leeren Token, der Handler
// bleibt unerreicht. Das bloße Entfernen des `!ok`-Zweigs färbt diesen Test
// ebenfalls **nicht** rot: der Nullwert der Metadata trägt die leere Kennung
// auf demselben Weg.
func TestAuthStreamInterceptorOhneEingehendeMetadataEndetMitUnauthenticated(t *testing.T) {
	interceptor := authStreamInterceptor(testTokens())
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

// unaryHandlerAufruf trägt einen `grpc.UnaryHandler`-Doppel, der beobachtbar
// macht, ob der Interceptor ihn erreicht hat.
func unaryHandlerAufruf() (*bool, grpc.UnaryHandler) {
	gerufen := false
	return &gerufen, func(context.Context, any) (any, error) { gerufen = true; return "ok", nil }
}

func unaryCtxMitToken(authorization string) context.Context {
	if authorization == "" {
		return context.Background()
	}
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(authorizationMetadataKey, authorization))
}

// TestAuthUnaryInterceptorOhneOderUnbekanntemTokenEndetMitUnauthenticated
// trägt die erste Hälfte der Fitness Function aus `ADR-0130` Teilfrage 4:
// ein Aufruf ohne oder mit unbekanntem `authorization`-Metadata-Wert endet
// mit `Unauthenticated`, der Handler wird nicht erreicht.
//
// Rot färbende Mutation: in `authUnaryInterceptor` den `caller ==
// apiauth.None`-Zweig entfernen — dann erreicht auch ein Aufruf ohne Token den
// Handler, dieser Test färbt rot.
func TestAuthUnaryInterceptorOhneOderUnbekanntemTokenEndetMitUnauthenticated(t *testing.T) {
	interceptor := authUnaryInterceptor(testTokens(), administrationRPCRoles)
	info := &grpc.UnaryServerInfo{FullMethod: "/cdc.administration.v1.Administration/ListTables"}
	for _, authorization := range []string{"", bearerPrefix, bearerPrefix + "unbekannt"} {
		gerufen, handler := unaryHandlerAufruf()
		_, err := interceptor(unaryCtxMitToken(authorization), nil, info, handler)
		if status.Code(err) != codes.Unauthenticated {
			t.Fatalf("authorization %q: Status %v (Erwartung: %v)", authorization, status.Code(err), codes.Unauthenticated)
		}
		if *gerufen {
			t.Fatalf("authorization %q: der Handler wurde ohne gültiges Token erreicht", authorization)
		}
	}
}

// TestAuthUnaryInterceptorReaderTokenGegenAdminRPCEndetMitPermissionDenied
// trägt die zweite Hälfte der Fitness Function: ein gültiges `reader`-Token
// gegen eine `apiauth.Admin`-RPC (`RegisterConsumer`) endet mit
// `PermissionDenied`, der Handler wird nicht erreicht.
//
// Rot färbende Mutation: in `administrationRPCRoles` den Eintrag
// `"RegisterConsumer": apiauth.Admin` auf `apiauth.Reader` ändern — dann erreicht das
// `reader`-Token den Handler, dieser Test färbt rot.
func TestAuthUnaryInterceptorReaderTokenGegenAdminRPCEndetMitPermissionDenied(t *testing.T) {
	interceptor := authUnaryInterceptor(testTokens(), administrationRPCRoles)
	info := &grpc.UnaryServerInfo{FullMethod: "/cdc.administration.v1.Administration/RegisterConsumer"}
	gerufen, handler := unaryHandlerAufruf()
	_, err := interceptor(unaryCtxMitToken(bearerPrefix+testReaderToken), nil, info, handler)
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("Status: %v (Erwartung: %v)", status.Code(err), codes.PermissionDenied)
	}
	if *gerufen {
		t.Fatal("der Handler wurde trotz unzureichender Rechtsklasse erreicht")
	}
}

// TestAuthUnaryInterceptorAdminTokenErreichtBeideRechtsklassen trägt die
// dritte Hälfte der Fitness Function: ein gültiges `admin`-Token erreicht
// sowohl eine `apiauth.Reader`- als auch eine `apiauth.Admin`-RPC.
func TestAuthUnaryInterceptorAdminTokenErreichtBeideRechtsklassen(t *testing.T) {
	interceptor := authUnaryInterceptor(testTokens(), administrationRPCRoles)
	for _, method := range []string{"RegisterConsumer", "ListTables"} {
		info := &grpc.UnaryServerInfo{FullMethod: "/cdc.administration.v1.Administration/" + method}
		gerufen, handler := unaryHandlerAufruf()
		if _, err := interceptor(unaryCtxMitToken(bearerPrefix+testAdminToken), nil, info, handler); err != nil {
			t.Fatalf("Methode %s mit admin-Token: %v (Erwartung: kein Fehler)", method, err)
		}
		if !*gerufen {
			t.Fatalf("Methode %s: der Handler wurde mit admin-Token nicht erreicht", method)
		}
	}
}

// TestAuthUnaryInterceptorReadChangesRechtsklasse trägt die Fitness Function
// des zehnten RPC (`ADR-0131` Teilfrage 4): ein Aufruf ohne oder mit
// unbekanntem Token endet mit `Unauthenticated`; ein gültiges `reader`-Token
// erreicht `ReadChanges`, ein gültiges `admin`-Token ebenfalls
// (`apiauth.Admin ≥ apiauth.Reader`).
func TestAuthUnaryInterceptorReadChangesRechtsklasse(t *testing.T) {
	interceptor := authUnaryInterceptor(testTokens(), administrationRPCRoles)
	info := &grpc.UnaryServerInfo{FullMethod: "/cdc.administration.v1.Administration/ReadChanges"}

	for _, authorization := range []string{"", bearerPrefix, bearerPrefix + "unbekannt"} {
		gerufen, handler := unaryHandlerAufruf()
		_, err := interceptor(unaryCtxMitToken(authorization), nil, info, handler)
		if status.Code(err) != codes.Unauthenticated {
			t.Fatalf("authorization %q: Status %v (Erwartung: %v)", authorization, status.Code(err), codes.Unauthenticated)
		}
		if *gerufen {
			t.Fatalf("authorization %q: der Handler wurde ohne gültiges Token erreicht", authorization)
		}
	}

	for _, token := range []string{testReaderToken, testAdminToken} {
		gerufen, handler := unaryHandlerAufruf()
		if _, err := interceptor(unaryCtxMitToken(bearerPrefix+token), nil, info, handler); err != nil {
			t.Fatalf("Token %q: %v (Erwartung: kein Fehler)", token, err)
		}
		if !*gerufen {
			t.Fatalf("Token %q: der Handler wurde nicht erreicht", token)
		}
	}
}

// TestAuthUnaryInterceptorDiagnoseRechtsklasse trägt die Fitness Function
// des elften RPC (`ADR-0132` Teilfrage 6): ein Aufruf ohne oder mit
// unbekanntem Token endet mit `Unauthenticated`; ein gültiges `reader`-Token
// erreicht `Diagnose`, ein gültiges `admin`-Token ebenfalls
// (`apiauth.Admin ≥ apiauth.Reader`).
func TestAuthUnaryInterceptorDiagnoseRechtsklasse(t *testing.T) {
	interceptor := authUnaryInterceptor(testTokens(), administrationRPCRoles)
	info := &grpc.UnaryServerInfo{FullMethod: "/cdc.administration.v1.Administration/Diagnose"}

	for _, authorization := range []string{"", bearerPrefix, bearerPrefix + "unbekannt"} {
		gerufen, handler := unaryHandlerAufruf()
		_, err := interceptor(unaryCtxMitToken(authorization), nil, info, handler)
		if status.Code(err) != codes.Unauthenticated {
			t.Fatalf("authorization %q: Status %v (Erwartung: %v)", authorization, status.Code(err), codes.Unauthenticated)
		}
		if *gerufen {
			t.Fatalf("authorization %q: der Handler wurde ohne gültiges Token erreicht", authorization)
		}
	}

	for _, token := range []string{testReaderToken, testAdminToken} {
		gerufen, handler := unaryHandlerAufruf()
		if _, err := interceptor(unaryCtxMitToken(bearerPrefix+token), nil, info, handler); err != nil {
			t.Fatalf("Token %q: %v (Erwartung: kein Fehler)", token, err)
		}
		if !*gerufen {
			t.Fatalf("Token %q: der Handler wurde nicht erreicht", token)
		}
	}
}

// TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin
// trägt den Fail-closed-Zweig: ein Methodenname ohne Eintrag in der
// Rechtsklassen-Tabelle fällt auf `apiauth.Admin` — ein `reader`-Token erreicht
// ihn nicht.
//
// Rot färbende Mutation: in `authUnaryInterceptor` den `!ok`-Zweig
// entfernen — der Nullwert von `apiauth.Role` ist `apiauth.None` (nicht `apiauth.Admin`),
// jedes bekannte Token würde die dann implizit offene RPC erreichen, dieser
// Test färbt rot.
func TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin(t *testing.T) {
	interceptor := authUnaryInterceptor(testTokens(), administrationRPCRoles)
	info := &grpc.UnaryServerInfo{FullMethod: "/cdc.administration.v1.Administration/UnbekannteMethode"}
	gerufen, handler := unaryHandlerAufruf()
	_, err := interceptor(unaryCtxMitToken(bearerPrefix+testReaderToken), nil, info, handler)
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("Status: %v (Erwartung: %v, fail-closed auf apiauth.Admin)", status.Code(err), codes.PermissionDenied)
	}
	if *gerufen {
		t.Fatal("der Handler wurde für eine unbekannte Methode mit reader-Token erreicht")
	}
}
