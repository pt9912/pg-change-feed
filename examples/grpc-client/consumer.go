package main

import (
	"fmt"

	administrationv1 "github.com/pt9912/pg-change-feed/gen/cdc/administration/v1"
)

// registerConsumer ruft die admin-RPC `RegisterConsumer` auf (`LH-FA-CON-001`).
func registerConsumer(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.adminToken)
	defer cancel()
	resp, err := client.RegisterConsumer(ctx, &administrationv1.RegisterConsumerRequest{ConsumerId: cfg.consumerID, Name: cfg.name})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc-client: consumer_id=%s name=%s already_registered=%v",
		resp.GetConsumerId(), resp.GetName(), resp.GetAlreadyRegistered()), nil
}

// acknowledgeConsumer ruft die admin-RPC `AcknowledgeConsumer` auf
// (`LH-FA-CON-004`).
func acknowledgeConsumer(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.adminToken)
	defer cancel()
	resp, err := client.AcknowledgeConsumer(ctx, &administrationv1.AcknowledgeConsumerRequest{
		ConsumerId: cfg.consumerID, SourceId: cfg.source, Offset: cfg.offset,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc-client: consumer_id=%s source_id=%s offset=%d",
		resp.GetConsumerId(), resp.GetSourceId(), resp.GetOffset()), nil
}

// getConsumerPosition ruft die reader-RPC `GetConsumerPosition` auf
// (`LH-FA-CON-005`).
func getConsumerPosition(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.token)
	defer cancel()
	resp, err := client.GetConsumerPosition(ctx, &administrationv1.GetConsumerPositionRequest{ConsumerId: cfg.consumerID})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc-client: consumer_id=%s source_id=%s offset=%d acknowledged=%v",
		resp.GetConsumerId(), resp.GetSourceId(), resp.GetOffset(), resp.GetAcknowledged()), nil
}

// removeConsumer ruft die admin-RPC `RemoveConsumer` auf (`LH-FA-CON-006`).
func removeConsumer(client administrationv1.AdministrationClient, cfg config) (string, error) {
	ctx, cancel := callCtx(cfg.adminToken)
	defer cancel()
	resp, err := client.RemoveConsumer(ctx, &administrationv1.RemoveConsumerRequest{ConsumerId: cfg.consumerID})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc-client: consumer_id=%s removed=%v", resp.GetConsumerId(), resp.GetRemoved()), nil
}
