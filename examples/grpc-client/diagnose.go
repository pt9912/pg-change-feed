package main

import (
	"fmt"
	"strings"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
)

// diagnose ruft die reader-RPC `Diagnose` auf (`LH-FA-SST-003`): denselben
// Bericht wie der CLI-Sondermodus und `GET /diagnose` — ein `known`/
// `present`-Feld auf `false` trägt den jeweiligen Abwesenheits-Fall (kein
// Lebenszeichen, kein Blocker, unbekannte Schätzung/unbekannter Rückstand).
func diagnose(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.token)
	defer cancel()
	resp, err := client.Diagnose(ctx, &administrationv1.DiagnoseRequest{Source: cfg.source})
	if err != nil {
		return "", err
	}

	lines := []string{fmt.Sprintf("grpc-client: heartbeat_known=%v heartbeat_age_seconds=%f heartbeat_error_class=%s capture_lag=%f storage_bytes=%f",
		resp.GetHeartbeat().GetKnown(), resp.GetHeartbeat().GetAgeSeconds(), resp.GetHeartbeat().GetErrorClass(),
		resp.GetCaptureLag(), resp.GetStorageBytes())}

	for _, cl := range resp.GetConsumerLags() {
		lines = append(lines, fmt.Sprintf("  consumer_lag: consumer_id=%s known=%v lag=%f", cl.GetConsumerId(), cl.GetKnown(), cl.GetLag()))
	}

	if rb := resp.GetRetentionBlocker(); rb.GetPresent() {
		lines = append(lines, fmt.Sprintf("  retention_blocker: consumer_id=%s name=%s acknowledged_position=%d backlog_known=%v backlog=%d",
			rb.GetConsumerId(), rb.GetName(), rb.GetAcknowledgedPosition(), rb.GetBacklogKnown(), rb.GetBacklog()))
	} else {
		lines = append(lines, "  retention_blocker: kein Blocker")
	}

	for _, b := range resp.GetBackfill() {
		lines = append(lines, fmt.Sprintf("  backfill: schema=%s table=%s status=%s rows_copied=%d estimated_rows_known=%v estimated_rows=%d warn_size=%v warn_duration=%v error_message=%s",
			b.GetSchema(), b.GetTable(), b.GetStatus(), b.GetRowsCopied(), b.GetEstimatedRowsKnown(), b.GetEstimatedRows(),
			b.GetWarnEstimatedSize(), b.GetWarnDuration(), b.GetErrorMessage()))
	}

	return strings.Join(lines, "\n"), nil
}
