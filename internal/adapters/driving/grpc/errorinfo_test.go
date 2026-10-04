package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	streamv1 "github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// TestStatusErrorTraegtErrorInfoMitCodeUndDomaene trägt: der Status behält
// Code und Text, das Detail `ErrorInfo` trägt den Meldungscode als `reason`
// und `pg-change-feed` als `domain`. Eingabeseite: der übergebene Code.
func TestStatusErrorTraegtErrorInfoMitCodeUndDomaene(t *testing.T) {
	err := statusError(codes.InvalidArgument, "Klartext", messagecode.RejectedValueInvalid)
	if got := status.Code(err); got != codes.InvalidArgument {
		t.Fatalf("Status %v (Erwartung %v)", got, codes.InvalidArgument)
	}
	if got := status.Convert(err).Message(); got != "Klartext" {
		t.Fatalf("Text %q (Erwartung %q)", got, "Klartext")
	}
	info := errorInfoOf(t, err)
	if info == nil || info.GetReason() != string(messagecode.RejectedValueInvalid) || info.GetDomain() != "pg-change-feed" {
		t.Fatalf("ErrorInfo = %v", info)
	}
}

// TestStreamOhneBroadcasterTraegtDenCodeDerVerdrahtung trägt: ein Stream ohne
// verdrahteten Broadcaster endet mit `Internal` und dem Code der
// Verdrahtungs-Vorbedingung.
func TestStreamOhneBroadcasterTraegtDenCodeDerVerdrahtung(t *testing.T) {
	service := &changeStreamService{log: outbound.NoopLog}
	err := service.StreamChanges(&streamv1.StreamChangesRequest{}, nil)
	if status.Code(err) != codes.Internal {
		t.Fatalf("Status %v (Erwartung %v)", status.Code(err), codes.Internal)
	}
	info := errorInfoOf(t, err)
	if info == nil || info.GetReason() != string(messagecode.WiringPrecondition) {
		t.Fatalf("ErrorInfo = %v, Erwartung reason %q", info, messagecode.WiringPrecondition)
	}
}

// TestAuthStatusTraegtKeinErrorInfo trägt: ein fehlendes oder unbekanntes
// Token (`Unauthenticated`) und eine unzureichende Rechtsklasse
// (`PermissionDenied`) tragen kein Statusdetail `ErrorInfo` — die Tabelle der
// Meldungscodes führt für sie keinen Eintrag. Eingabeseite: die drei
// Interceptor-Pfade (Stream, unär ohne Token, unär mit unzureichendem Token).
func TestAuthStatusTraegtKeinErrorInfo(t *testing.T) {
	streamInterceptor := authStreamInterceptor(testTokens())
	streamErr := streamInterceptor(nil, &fakeServerStream{ctx: context.Background()},
		&grpc.StreamServerInfo{FullMethod: streamv1.ChangeStream_StreamChanges_FullMethodName},
		func(any, grpc.ServerStream) error { return nil })

	unary := authUnaryInterceptor(testTokens(), administrationRPCRoles)
	info := &grpc.UnaryServerInfo{FullMethod: "/cdc.administration.v1.Administration/RegisterConsumer"}
	handler := func(context.Context, any) (any, error) { return nil, nil }
	_, unauthenticated := unary(context.Background(), nil, info, handler)
	_, denied := unary(metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(authorizationMetadataKey, bearerPrefix+testReaderToken)), nil, info, handler)

	for _, tc := range []struct {
		name string
		err  error
		want codes.Code
	}{
		{"Stream ohne Token", streamErr, codes.Unauthenticated},
		{"unär ohne Token", unauthenticated, codes.Unauthenticated},
		{"unär mit Reader-Token gegen Admin-RPC", denied, codes.PermissionDenied},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if status.Code(tc.err) != tc.want {
				t.Fatalf("Status %v (Erwartung %v)", status.Code(tc.err), tc.want)
			}
			if got := errorInfoOf(t, tc.err); got != nil {
				t.Fatalf("Status trägt ein ErrorInfo: %v", got)
			}
		})
	}
}
