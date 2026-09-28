using Cdc.Administration.V1;

namespace CdcExamples.Grpc;

/// <summary>
/// DiagnoseClient ruft die reader-RPC <c>Diagnose</c> des
/// <c>Administration</c>-Diensts auf. Form-Vorbild:
/// <c>examples/grpc-client/diagnose.go</c>.
/// </summary>
public static class DiagnoseClient
{
    /// <summary>
    /// DiagnoseAsync ruft <c>Diagnose</c> auf (<c>LH-FA-SST-003</c>): denselben
    /// Bericht wie der CLI-Sondermodus und <c>GET /diagnose</c>.
    /// </summary>
    public static async Task<string> DiagnoseAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.DiagnoseAsync(
            new DiagnoseRequest { Source = cfg.Source },
            headers: CallMetadata.Headers(cfg.Token), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatDiagnose(resp);
    }
}
