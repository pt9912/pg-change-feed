using Xunit;

namespace CdcExamples.Nats.Tests;

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
            "CDC_NATS_URL" => "nats://feed:4222",
            "CDC_HTTP_ADDR" => "feed:8080",
            "CDC_API_TOKEN_READER" => "tok-env",
            _ => null,
        });

        Assert.Equal(new Config("nats://feed:4222", "feed:8080", "tok-env", "", "", ""), got);
    }

    [Fact]
    public void ParseFlagsOverrideEnvironment()
    {
        var got = Cli.Parse(
            new[]
            {
                "--nats-url", "nats://override:4222",
                "--addr", "override:9090",
                "--token=tok-flag",
                "--source", "quelle-1",
                "--schema", "public",
                "--table", "orders",
            },
            name => name == "CDC_NATS_URL" ? "nats://feed:4222" : null);

        Assert.Equal(new Config("nats://override:4222", "override:9090", "tok-flag", "quelle-1", "public", "orders"), got);
    }

    [Fact]
    public void ParseRejectsUnknownFlag()
    {
        Assert.Throws<ArgumentException>(() => Cli.Parse(new[] { "--unknown" }, EmptyEnv));
    }

    [Fact]
    public void ParseRejectsMissingFlagValue()
    {
        Assert.Throws<ArgumentException>(() => Cli.Parse(new[] { "--nats-url" }, EmptyEnv));
    }
}
