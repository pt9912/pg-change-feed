using Cdc.Administration.V1;

namespace CdcExamples.Grpc;

/// <summary>
/// ConsumerClient ruft die vier Consumer-Verwaltungs-RPCs des
/// <c>Administration</c>-Diensts auf. Form-Vorbild:
/// <c>examples/grpc-client/consumer.go</c>.
/// </summary>
public static class ConsumerClient
{
    /// <summary>RegisterConsumerAsync ruft die admin-RPC <c>RegisterConsumer</c> auf (<c>LH-FA-CON-001</c>).</summary>
    public static async Task<string> RegisterConsumerAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.RegisterConsumerAsync(
            new RegisterConsumerRequest { ConsumerId = cfg.ConsumerId, Name = cfg.Name },
            headers: CallMetadata.Headers(cfg.AdminToken), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatRegisterConsumer(resp);
    }

    /// <summary>AcknowledgeConsumerAsync ruft die admin-RPC <c>AcknowledgeConsumer</c> auf (<c>LH-FA-CON-004</c>).</summary>
    public static async Task<string> AcknowledgeConsumerAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.AcknowledgeConsumerAsync(
            new AcknowledgeConsumerRequest { ConsumerId = cfg.ConsumerId, SourceId = cfg.Source, Offset = cfg.Offset },
            headers: CallMetadata.Headers(cfg.AdminToken), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatAcknowledgeConsumer(resp);
    }

    /// <summary>GetConsumerPositionAsync ruft die reader-RPC <c>GetConsumerPosition</c> auf (<c>LH-FA-CON-005</c>).</summary>
    public static async Task<string> GetConsumerPositionAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.GetConsumerPositionAsync(
            new GetConsumerPositionRequest { ConsumerId = cfg.ConsumerId },
            headers: CallMetadata.Headers(cfg.Token), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatConsumerPosition(resp);
    }

    /// <summary>RemoveConsumerAsync ruft die admin-RPC <c>RemoveConsumer</c> auf (<c>LH-FA-CON-006</c>).</summary>
    public static async Task<string> RemoveConsumerAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.RemoveConsumerAsync(
            new RemoveConsumerRequest { ConsumerId = cfg.ConsumerId },
            headers: CallMetadata.Headers(cfg.AdminToken), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatRemoveConsumer(resp);
    }
}
