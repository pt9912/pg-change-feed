using Cdc.Administration.V1;

namespace CdcExamples.Grpc;

/// <summary>
/// RetentionClient ruft die admin-RPC <c>RunRetention</c> des
/// <c>Administration</c>-Diensts auf. Form-Vorbild:
/// <c>examples/grpc-client/retention.go</c>.
/// </summary>
public static class RetentionClient
{
    /// <summary>RunRetentionAsync ruft <c>RunRetention</c> auf (<c>LH-FA-RET-002</c>): <c>MinAgeNanos</c> 0 heißt „kein zeitliches Mindestalter" und ist gültig.</summary>
    public static async Task<string> RunRetentionAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.RunRetentionAsync(
            new RunRetentionRequest { Source = cfg.Source, MinAgeNanos = cfg.MinAgeNanos },
            headers: CallMetadata.Headers(cfg.AdminToken), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatRunRetention(resp);
    }
}
