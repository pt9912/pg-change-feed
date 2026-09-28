using Xunit;

namespace CdcExamples.Http.Tests;

/// <summary>
/// Prüft den Aufbau der Laufzeit-Konfiguration aus Flags und einem
/// injizierten Umgebungs-Lookup — netzlos, reine Funktion.
/// </summary>
public class CliTests
{
    private static string? EmptyEnv(string _) => null;

    private static Config ExpectedDefaults(string addr = "", string token = "", string adminToken = "", string verb = "tables", string source = "", string publication = "") => new(
        Addr: addr, Token: token, AdminToken: adminToken, Verb: verb,
        Source: source, Publication: publication,
        ConsumerId: "", Name: "", Offset: 0,
        Schema: "", Table: "", TableId: "", SchemaVersionId: "", Version: 1,
        From: "", To: "", Limit: "", MinAgeNanos: 0);

    [Fact]
    public void ParseUsesEnvironmentDefaults()
    {
        var got = Cli.Parse(Array.Empty<string>(), name => name switch
        {
            "CDC_HTTP_ADDR" => "feed:8080",
            "CDC_API_TOKEN_READER" => "tok-env",
            "CDC_API_TOKEN_ADMIN" => "tok-admin-env",
            _ => null,
        });

        Assert.Equal(ExpectedDefaults(addr: "feed:8080", token: "tok-env", adminToken: "tok-admin-env"), got);
    }

    [Fact]
    public void ParseFlagsOverrideEnvironment()
    {
        var got = Cli.Parse(
            new[] { "--addr", "override:9090", "--token=tok-flag", "--source", "quelle-1", "--publication", "pub-1" },
            name => name == "CDC_HTTP_ADDR" ? "feed:8080" : null);

        Assert.Equal(
            ExpectedDefaults(addr: "override:9090", token: "tok-flag", source: "quelle-1", publication: "pub-1"),
            got);
    }

    [Fact]
    public void ParseVerbAndCapabilityFieldsFillAllTenFieldsIndependently()
    {
        var got = Cli.Parse(
            new[]
            {
                "--verb", "enable-table", "--admin-token", "admin-tok",
                "--source", "src", "--schema", "public", "--table", "orders",
                "--table-id", "public.orders", "--schema-version-id", "public.orders-v1",
                "--version", "2", "--publication", "pub",
                "--consumer-id", "c-1", "--name", "n-1", "--offset", "42",
                "--from", "1", "--to", "10", "--limit", "5", "--min-age-nanos", "-3",
            },
            EmptyEnv);

        Assert.Equal(new Config(
            Addr: "", Token: "", AdminToken: "admin-tok", Verb: "enable-table",
            Source: "src", Publication: "pub",
            ConsumerId: "c-1", Name: "n-1", Offset: 42,
            Schema: "public", Table: "orders", TableId: "public.orders", SchemaVersionId: "public.orders-v1", Version: 2,
            From: "1", To: "10", Limit: "5", MinAgeNanos: -3), got);
    }

    [Fact]
    public void ParseInlineFlagValueSplitsOnFirstEquals()
    {
        var got = Cli.Parse(new[] { "--verb=changes" }, EmptyEnv);
        Assert.Equal("changes", got.Verb);
    }

    [Fact]
    public void ParseRejectsUnknownFlag()
    {
        Assert.Throws<ArgumentException>(() => Cli.Parse(new[] { "--unknown" }, EmptyEnv));
    }

    [Fact]
    public void ParseRejectsMissingFlagValue()
    {
        Assert.Throws<ArgumentException>(() => Cli.Parse(new[] { "--addr" }, EmptyEnv));
    }

    [Fact]
    public void ParseRejectsNonNumericOffset()
    {
        Assert.Throws<ArgumentException>(() => Cli.Parse(new[] { "--offset", "nicht-zahl" }, EmptyEnv));
    }

    [Fact]
    public void ParseRejectsNonNumericVersion()
    {
        Assert.Throws<ArgumentException>(() => Cli.Parse(new[] { "--version", "nicht-zahl" }, EmptyEnv));
    }

    [Fact]
    public void ParseRejectsNonNumericMinAgeNanos()
    {
        Assert.Throws<ArgumentException>(() => Cli.Parse(new[] { "--min-age-nanos", "nicht-zahl" }, EmptyEnv));
    }
}
