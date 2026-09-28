using Cdc.Administration.V1;

namespace CdcExamples.Grpc;

/// <summary>
/// TablesAdminClient ruft die vier Tabellen-Verwaltungs-RPCs des
/// <c>Administration</c>-Diensts auf. Form-Vorbild:
/// <c>examples/grpc-client/tables_admin.go</c>.
/// </summary>
public static class TablesAdminClient
{
    /// <summary>EnableTableAsync ruft die admin-RPC <c>EnableTable</c> auf (<c>LH-FA-CFG-001</c>).</summary>
    public static async Task<string> EnableTableAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.EnableTableAsync(
            new EnableTableRequest
            {
                Source = cfg.Source, Schema = cfg.Schema, Table = cfg.Table,
                TableId = cfg.TableId, SchemaVersionId = cfg.SchemaVersionId,
                Version = cfg.Version, Publication = cfg.Publication,
            },
            headers: CallMetadata.Headers(cfg.AdminToken), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatEnableTable(resp);
    }

    /// <summary>DisableTableAsync ruft die admin-RPC <c>DisableTable</c> auf (<c>LH-FA-CFG-002</c>).</summary>
    public static async Task<string> DisableTableAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.DisableTableAsync(
            new DisableTableRequest { Source = cfg.Source, Schema = cfg.Schema, Table = cfg.Table, Publication = cfg.Publication },
            headers: CallMetadata.Headers(cfg.AdminToken), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatDisableTable(resp);
    }

    /// <summary>GetTableStatusAsync ruft die reader-RPC <c>GetTableStatus</c> auf (<c>LH-FA-CFG-003</c>).</summary>
    public static async Task<string> GetTableStatusAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.GetTableStatusAsync(
            new GetTableStatusRequest { Source = cfg.Source, Schema = cfg.Schema, Table = cfg.Table, Publication = cfg.Publication },
            headers: CallMetadata.Headers(cfg.Token), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatTableStatus(resp);
    }

    /// <summary>ListTablesAsync ruft die reader-RPC <c>ListTables</c> auf (<c>LH-FA-CFG-004</c>).</summary>
    public static async Task<string> ListTablesAsync(
        Administration.AdministrationClient client, Config cfg, CancellationToken cancellationToken = default)
    {
        var resp = await client.ListTablesAsync(
            new ListTablesRequest { Source = cfg.Source, Publication = cfg.Publication },
            headers: CallMetadata.Headers(cfg.Token), deadline: CallMetadata.Deadline(),
            cancellationToken: cancellationToken);
        return Format.FormatListTables(resp);
    }
}
