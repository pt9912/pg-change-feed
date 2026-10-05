namespace CdcExamples.Grpc;

/// <summary>
/// Cli liest Flag-Werte und füllt fehlende Felder aus den im Handbuch
/// dokumentierten Umgebungsvariablen (<c>ADR-0076</c> Festlegung 1). Das
/// Umgebungs-Lookup ist injiziert (<paramref name="getEnv"/> statt
/// <c>Environment.GetEnvironmentVariable</c> direkt) — das macht die Funktion
/// rein und netzlos testbar, ohne echte Prozessumgebung zu setzen. Diese
/// Funktion prüft nur die Flag-Syntax; ob ein Verb bekannt ist oder seine
/// Pflichtfelder trägt, prüft <see cref="Validator"/>. Form-Vorbild:
/// <c>examples/csharp/http-client/Cli.cs</c>, <c>examples/grpc-client/main.go</c>
/// (<c>parseFlags</c>).
/// </summary>
public static class Cli
{
    public static Config Parse(string[] args, Func<string, string?> getEnv)
    {
        var addr = getEnv("CDC_GRPC_ADDR") ?? "";
        var token = getEnv("CDC_API_TOKEN_READER") ?? "";
        var adminToken = getEnv("CDC_API_TOKEN_ADMIN") ?? "";
        var caFile = getEnv("CDC_TLS_CA_FILE") ?? "";
        var verb = "stream";
        var schema = "";
        var table = "";
        var target = "";
        var consumerId = "";
        var name = "";
        var offset = 0UL;
        var tableId = "";
        var schemaVersionId = "";
        var version = 1L;
        var source = "";
        var publication = "";
        var from = 0UL;
        var to = 0UL;
        var limit = 0L;
        var minAgeNanos = 0L;

        for (var i = 0; i < args.Length; i++)
        {
            var (flagName, inline) = SplitFlag(args[i]);

            string NextValue()
            {
                if (inline is not null)
                {
                    return inline;
                }
                if (i + 1 >= args.Length)
                {
                    throw new ArgumentException($"grpc-client: Flag {flagName} braucht einen Wert");
                }
                return args[++i];
            }

            switch (flagName)
            {
                case "--addr":
                    addr = NextValue();
                    break;
                case "--token":
                    token = NextValue();
                    break;
                case "--admin-token":
                    adminToken = NextValue();
                    break;
                case "--ca-file":
                    caFile = NextValue();
                    break;
                case "--verb":
                    verb = NextValue();
                    break;
                case "--schema":
                    schema = NextValue();
                    break;
                case "--target":
                    target = NextValue();
                    break;
                case "--table":
                    table = NextValue();
                    break;
                case "--consumer-id":
                    consumerId = NextValue();
                    break;
                case "--name":
                    name = NextValue();
                    break;
                case "--offset":
                    offset = ParseUInt64(flagName, NextValue());
                    break;
                case "--table-id":
                    tableId = NextValue();
                    break;
                case "--schema-version-id":
                    schemaVersionId = NextValue();
                    break;
                case "--version":
                    version = ParseInt64(flagName, NextValue());
                    break;
                case "--source":
                    source = NextValue();
                    break;
                case "--publication":
                    publication = NextValue();
                    break;
                case "--from":
                    from = ParseUInt64(flagName, NextValue());
                    break;
                case "--to":
                    to = ParseUInt64(flagName, NextValue());
                    break;
                case "--limit":
                    limit = ParseInt64(flagName, NextValue());
                    break;
                case "--min-age-nanos":
                    minAgeNanos = ParseInt64(flagName, NextValue());
                    break;
                default:
                    throw new ArgumentException($"grpc-client: unbekanntes Flag {flagName}");
            }
        }

        return new Config(
            Addr: addr, Token: token, AdminToken: adminToken, Verb: verb,
            Schema: schema, Table: table,
            ConsumerId: consumerId, Name: name, Offset: offset,
            TableId: tableId, SchemaVersionId: schemaVersionId, Version: version,
            Source: source, Publication: publication,
            From: from, To: to, Limit: limit, MinAgeNanos: minAgeNanos, Target: target, CaFile: caFile);
    }

    private static ulong ParseUInt64(string flagName, string value)
    {
        if (!ulong.TryParse(value, out var parsed))
        {
            throw new ArgumentException($"grpc-client: Flag {flagName} braucht eine nicht-negative Ganzzahl, erhielt \"{value}\"");
        }
        return parsed;
    }

    private static long ParseInt64(string flagName, string value)
    {
        if (!long.TryParse(value, out var parsed))
        {
            throw new ArgumentException($"grpc-client: Flag {flagName} braucht eine Ganzzahl, erhielt \"{value}\"");
        }
        return parsed;
    }

    /// <summary>
    /// SplitFlag zerlegt <c>--name=wert</c> in Name und Inline-Wert; ein Flag
    /// ohne <c>=</c> liefert keinen Inline-Wert — der nächste Token trägt ihn.
    /// </summary>
    private static (string Name, string? InlineValue) SplitFlag(string arg)
    {
        var eq = arg.IndexOf('=');
        return eq < 0 ? (arg, null) : (arg[..eq], arg[(eq + 1)..]);
    }
}
