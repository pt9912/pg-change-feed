namespace CdcExamples.Grpc;

/// <summary>
/// Validator prüft eine <see cref="Config"/> gegen die Pflichtfelder des
/// gewählten Verbs — vor jedem Netzwerkaufruf. Jedes Verb braucht die
/// Horch-Adresse und die zu seiner Rechtsklasse passende Token-Variable
/// (Rechtsklassen-Tabelle: <c>ADR-0130</c>). Form-Vorbild:
/// <c>examples/csharp/http-client/Validator.cs</c>,
/// <c>validate</c> in <c>examples/grpc-client/main.go</c>.
/// </summary>
public static class Validator
{
    /// <summary>KnownVerbs trägt die geschlossene Menge der <c>--verb</c>-Werte.</summary>
    public static readonly IReadOnlySet<string> KnownVerbs = new HashSet<string>
    {
        "stream", "register-consumer", "acknowledge-consumer", "get-consumer-position", "remove-consumer",
        "enable-table", "disable-table", "get-table-status", "list-tables", "run-retention",
        "read-changes", "diagnose",
    };

    /// <summary>
    /// Validate liefert eine Fehlerbeschreibung, wenn <paramref name="cfg"/>
    /// nicht gegen den gewählten Verb-Vertrag genügt, sonst <c>null</c>.
    /// </summary>
    public static string? Validate(Config cfg)
    {
        if (!KnownVerbs.Contains(cfg.Verb))
        {
            return $"unbekanntes --verb \"{cfg.Verb}\"";
        }
        if (string.IsNullOrEmpty(cfg.Addr))
        {
            return "keine gRPC-Adresse gesetzt — CDC_GRPC_ADDR (oder --addr) ist noetig, um den Server zu erreichen";
        }

        return cfg.Verb switch
        {
            "stream" => RequireReaderToken(cfg),
            "register-consumer" => RequireAdminToken(cfg) ?? RequireFields(
                cfg.ConsumerId, cfg.Name,
                "--consumer-id und --name sind Pflicht — sie sind die zwei Pflichtfelder von RegisterConsumer"),
            "acknowledge-consumer" => RequireAdminToken(cfg) ?? RequireFields(
                cfg.ConsumerId, cfg.Source,
                "--consumer-id und --source sind Pflicht — sie sind zwei der drei Pflichtfelder von AcknowledgeConsumer"),
            "get-consumer-position" => RequireReaderToken(cfg) ?? RequireFields(
                cfg.ConsumerId,
                "--consumer-id ist Pflicht — es ist das einzige Pflichtfeld von GetConsumerPosition"),
            "remove-consumer" => RequireAdminToken(cfg) ?? RequireFields(
                cfg.ConsumerId,
                "--consumer-id ist Pflicht — es ist das einzige Pflichtfeld von RemoveConsumer"),
            "enable-table" => RequireAdminToken(cfg) ?? RequireFields(
                cfg.Source, cfg.Schema, cfg.Table, cfg.Publication,
                "--source, --schema, --table und --publication sind Pflicht — table-id und schema-version-id haben einen Default"),
            "disable-table" => RequireAdminToken(cfg) ?? RequireFields(
                cfg.Source, cfg.Schema, cfg.Table, cfg.Publication,
                "--source, --schema, --table und --publication sind Pflicht — sie sind die vier Pflichtfelder von DisableTable"),
            "get-table-status" => RequireReaderToken(cfg) ?? RequireFields(
                cfg.Source, cfg.Schema, cfg.Table, cfg.Publication,
                "--source, --schema, --table und --publication sind Pflicht — sie sind die vier Pflichtfelder von GetTableStatus"),
            "list-tables" => RequireReaderToken(cfg) ?? RequireFields(
                cfg.Source, cfg.Publication,
                "--source und --publication sind Pflicht — sie sind die zwei Pflichtfelder von ListTables"),
            "run-retention" => RequireAdminToken(cfg) ?? RequireFields(
                cfg.Source,
                "--source ist Pflicht — es ist eines der zwei Pflichtfelder von RunRetention"),
            "read-changes" => RequireReaderToken(cfg) ?? RequireFields(
                cfg.Source,
                "--source ist Pflicht — es ist das einzige Pflichtfeld von ReadChanges"),
            "diagnose" => RequireReaderToken(cfg) ?? RequireFields(
                cfg.Source,
                "--source ist Pflicht — es ist das einzige Pflichtfeld von Diagnose"),
            _ => null,
        };
    }

    private static string? RequireFields(string field, string errorMessage) =>
        string.IsNullOrEmpty(field) ? errorMessage : null;

    private static string? RequireFields(string field1, string field2, string errorMessage) =>
        string.IsNullOrEmpty(field1) || string.IsNullOrEmpty(field2) ? errorMessage : null;

    private static string? RequireFields(string field1, string field2, string field3, string field4, string errorMessage) =>
        string.IsNullOrEmpty(field1) || string.IsNullOrEmpty(field2) || string.IsNullOrEmpty(field3) || string.IsNullOrEmpty(field4)
            ? errorMessage
            : null;

    private static string? RequireReaderToken(Config cfg) =>
        string.IsNullOrEmpty(cfg.Token)
            ? "kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen"
            : null;

    private static string? RequireAdminToken(Config cfg) =>
        string.IsNullOrEmpty(cfg.AdminToken)
            ? "kein Admin-Token gesetzt — CDC_API_TOKEN_ADMIN (oder --admin-token) ist noetig, um ueber die admin-Rechtsklasse zu schreiben"
            : null;
}
