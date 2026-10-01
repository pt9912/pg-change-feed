using Cdc.Administration.V1;

namespace CdcExamples.Grpc;

/// <summary>
/// ChangesClient ruft die reader-RPC <c>ReadChanges</c> des
/// <c>Administration</c>-Diensts auf. Form-Vorbild:
/// <c>examples/grpc-client/changes.go</c>.
/// </summary>
public static class ChangesClient
{
    /// <summary>
    /// BuildRequest bildet die Anfrage von <c>ReadChanges</c> aus der
    /// Konfiguration; ein leeres <c>Target</c> ist kein Filter.
    /// </summary>
    public static ReadChangesRequest BuildRequest(Config cfg) => new()
    {
        Source = cfg.Source, Schema = cfg.Schema, Table = cfg.Table, Target = cfg.Target,
        From = cfg.From, To = cfg.To, Limit = cfg.Limit,
    };

    /// <summary>
    /// ReadChangesAsync ruft <c>ReadChanges</c> auf (<c>LH-FA-SST-006</c>):
    /// einen begrenzten Bereich persistierter Änderungen, dieselbe Filter-
    /// und Bereichs-Semantik wie <c>GET /changes</c> — ein leerer
    /// <c>From</c>/<c>To</c>/<c>Limit</c>-Wert trägt <c>0</c> (nicht gesetzt).
    /// </summary>
    public static async Task<string> ReadChangesAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.ReadChangesAsync(
            BuildRequest(cfg),
            headers: CallMetadata.Headers(cfg.Token), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatReadChanges(resp);
    }
}
