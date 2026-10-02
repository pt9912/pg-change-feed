using Xunit;

namespace CdcExamples.Sse.Tests;

/// <summary>
/// Prüft den Aufbau der Laufzeit-Konfiguration aus Flags und einem
/// injizierten Umgebungs-Lookup — netzlos, reine Funktion. Form-Vorbild:
/// <c>examples/csharp/http-client/HttpClient.Tests/CliTests.cs</c>.
/// </summary>
public class CliTests
{
    private static string? EmptyEnv(string _) => null;

    [Fact]
    public void ParseUsesEnvironmentDefaults()
    {
        var got = Cli.Parse(Array.Empty<string>(), name => name switch
        {
            "CDC_HTTP_ADDR" => "feed:8080",
            "CDC_API_TOKEN_READER" => "tok-env",
            _ => null,
        });

        Assert.Equal(new Config("feed:8080", "tok-env"), got);
    }

    [Fact]
    public void ParseFlagsOverrideEnvironment()
    {
        var got = Cli.Parse(
            new[] { "--addr", "override:9090", "--token=tok-flag" },
            name => name == "CDC_HTTP_ADDR" ? "feed:8080" : null);

        Assert.Equal(new Config("override:9090", "tok-flag"), got);
    }

    [Fact]
    public void ParseTargetFlagFillsTheConfigInBothFlagForms()
    {
        Assert.Equal("eu", Cli.Parse(new[] { "--target", "eu" }, EmptyEnv).Target);
        Assert.Equal("eu", Cli.Parse(new[] { "--target=eu" }, EmptyEnv).Target);
        Assert.Equal("", Cli.Parse(Array.Empty<string>(), EmptyEnv).Target);
    }

    [Fact]
    public void ParseSchemaAndTableFlagsFillTheConfigInBothFlagForms()
    {
        var got = Cli.Parse(new[] { "--schema", "public", "--table=orders" }, EmptyEnv);

        Assert.Equal("public", got.Schema);
        Assert.Equal("orders", got.Table);
        var none = Cli.Parse(Array.Empty<string>(), EmptyEnv);
        Assert.Equal("", none.Schema);
        Assert.Equal("", none.Table);
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
}
