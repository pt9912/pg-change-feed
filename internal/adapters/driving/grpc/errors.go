package grpc

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
)

// errorDomain ist die `domain` des Statusdetails `ErrorInfo`: der Meldungscode
// steht als `reason` unter dieser Domäne.
const errorDomain = "pg-change-feed"

// statusError baut den gRPC-Status mit dem Detail `google.rpc.ErrorInfo`
// (`reason` = Meldungscode, `domain` = `errorDomain`). Der gRPC-Statuscode
// bleibt der übergebene.
func statusError(code codes.Code, message string, reason messagecode.Code) error {
	st := status.New(code, message)
	withInfo, err := st.WithDetails(&errdetails.ErrorInfo{Reason: string(reason), Domain: errorDomain})
	if err != nil {
		return st.Err()
	}
	return withInfo.Err()
}
