using Xunit;

namespace CdcExamples.Grpc.Tests;

/// <summary>
/// Prüft <see cref="Cli.Parse"/>: Umgebungs-Default, Flag-Override, beide
/// Flag-Formen (<c>--name=wert</c> und <c>--name wert</c>) und die
/// Fehlerfälle (unbekanntes Flag, fehlender Wert, ungültige Ganzzahl). Diese
/// Funktion prüft nur die Flag-Syntax; Rechtsklassen und Pflichtfelder je
/// Verb prüft <see cref="ValidatorTests"/>. Form-Vorbild:
/// <c>examples/csharp/http-client/HttpClient.Tests/CliTests.cs</c>.
/// </summary>
public class CliTests
{
    private static string? NoEnv(string _) => null;

    private static Func<string, string?> EnvFrom(Dictionary<string, string> values) =>
        key => values.TryGetValue(key, out var value) ? value : null;

    [Fact]
    public void Parse_UsesEnvironmentDefaults()
    {
        var env = EnvFrom(new Dictionary<string, string>
        {
            ["CDC_GRPC_ADDR"] = "feed:9090",
            ["CDC_API_TOKEN_READER"] = "reader-token",
            ["CDC_API_TOKEN_ADMIN"] = "admin-token",
        });

        var cfg = Cli.Parse([], env);

        Assert.Equal("feed:9090", cfg.Addr);
        Assert.Equal("reader-token", cfg.Token);
        Assert.Equal("admin-token", cfg.AdminToken);
    }

    [Fact]
    public void Parse_DefaultVerbIsStream()
    {
        var cfg = Cli.Parse([], NoEnv);

        Assert.Equal("stream", cfg.Verb);
    }

    [Fact]
    public void Parse_FlagOverridesEnvironment()
    {
        var env = EnvFrom(new Dictionary<string, string>
        {
            ["CDC_GRPC_ADDR"] = "feed:9090",
            ["CDC_API_TOKEN_READER"] = "reader-token",
        });

        var cfg = Cli.Parse(["--addr", "other:9091", "--token=other-token"], env);

        Assert.Equal("other:9091", cfg.Addr);
        Assert.Equal("other-token", cfg.Token);
    }

    [Fact]
    public void Parse_NoEnvironmentNoFlags_ReturnsEmpty()
    {
        var cfg = Cli.Parse([], NoEnv);

        Assert.Equal("", cfg.Addr);
        Assert.Equal("", cfg.Token);
        Assert.Equal("", cfg.AdminToken);
    }

    [Fact]
    public void Parse_StreamFilterFlags()
    {
        var cfg = Cli.Parse(["--schema=public", "--table=orders"], NoEnv);

        Assert.Equal("public", cfg.Schema);
        Assert.Equal("orders", cfg.Table);
    }

    [Fact]
    public void Parse_TargetFlagInBothFlagForms()
    {
        Assert.Equal("eu", Cli.Parse(["--target=eu"], NoEnv).Target);
        Assert.Equal("eu", Cli.Parse(["--target", "eu"], NoEnv).Target);
        Assert.Equal("", Cli.Parse([], NoEnv).Target);
    }

    [Fact]
    public void StreamRequest_CarriesTargetWithSchemaAndTable()
    {
        var cfg = Cli.Parse(["--schema=public", "--table=orders", "--target=eu"], NoEnv);

        var request = StreamRequest.Build(cfg);

        Assert.Equal("public", request.Schema);
        Assert.Equal("orders", request.Table);
        Assert.Equal("eu", request.Target);
        Assert.Equal(new Cdc.Stream.V1.StreamChangesRequest(), StreamRequest.Build(Cli.Parse([], NoEnv)));
    }

    [Fact]
    public void ReadChangesRequest_CarriesTarget()
    {
        var cfg = Cli.Parse(["--source=quelle-1", "--schema=public", "--target=eu"], NoEnv);

        var request = ChangesClient.BuildRequest(cfg);

        Assert.Equal("quelle-1", request.Source);
        Assert.Equal("public", request.Schema);
        Assert.Equal("eu", request.Target);
        Assert.Equal("", ChangesClient.BuildRequest(Cli.Parse(["--source=quelle-1"], NoEnv)).Target);
    }

    [Fact]
    public void Parse_EnableTableFields()
    {
        var cfg = Cli.Parse(
            ["--verb=enable-table", "--source=quelle-1", "--schema=public", "--table=orders",
                "--table-id=public.orders", "--schema-version-id=public.orders-v2", "--version=2",
                "--publication=pub_quelle_1"],
            NoEnv);

        Assert.Equal("enable-table", cfg.Verb);
        Assert.Equal("quelle-1", cfg.Source);
        Assert.Equal("public.orders", cfg.TableId);
        Assert.Equal("public.orders-v2", cfg.SchemaVersionId);
        Assert.Equal(2L, cfg.Version);
        Assert.Equal("pub_quelle_1", cfg.Publication);
    }

    [Fact]
    public void Parse_ConsumerFields()
    {
        var cfg = Cli.Parse(["--consumer-id=c-1", "--name=Consumer", "--offset=42"], NoEnv);

        Assert.Equal("c-1", cfg.ConsumerId);
        Assert.Equal("Consumer", cfg.Name);
        Assert.Equal(42UL, cfg.Offset);
    }

    [Fact]
    public void Parse_ReadChangesRangeFields()
    {
        var cfg = Cli.Parse(["--from=10", "--to=20", "--limit=5"], NoEnv);

        Assert.Equal(10UL, cfg.From);
        Assert.Equal(20UL, cfg.To);
        Assert.Equal(5L, cfg.Limit);
    }

    [Fact]
    public void Parse_RunRetentionMinAgeNanos()
    {
        var cfg = Cli.Parse(["--min-age-nanos=1000000000"], NoEnv);

        Assert.Equal(1_000_000_000L, cfg.MinAgeNanos);
    }

    [Fact]
    public void Parse_UnknownFlag_Throws()
    {
        var ex = Assert.Throws<ArgumentException>(() => Cli.Parse(["--bogus"], NoEnv));
        Assert.Contains("unbekanntes Flag", ex.Message);
    }

    [Fact]
    public void Parse_FlagMissingValue_Throws()
    {
        var ex = Assert.Throws<ArgumentException>(() => Cli.Parse(["--addr"], NoEnv));
        Assert.Contains("braucht einen Wert", ex.Message);
    }

    [Fact]
    public void Parse_OffsetNotAnInteger_Throws()
    {
        var ex = Assert.Throws<ArgumentException>(() => Cli.Parse(["--offset=abc"], NoEnv));
        Assert.Contains("nicht-negative Ganzzahl", ex.Message);
    }

    [Fact]
    public void Parse_VersionNotAnInteger_Throws()
    {
        var ex = Assert.Throws<ArgumentException>(() => Cli.Parse(["--version=abc"], NoEnv));
        Assert.Contains("Ganzzahl", ex.Message);
    }
}
