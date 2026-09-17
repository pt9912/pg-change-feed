using Xunit;

namespace CdcExamples.Grpc.Tests;

/// <summary>
/// Prüft <see cref="Cli.Parse"/>: Umgebungs-Default, Flag-Override, beide
/// Flag-Formen (<c>--name=wert</c> und <c>--name wert</c>) und die
/// Fehlerfälle (unbekanntes Flag, fehlender Wert). Form-Vorbild:
/// <c>examples/csharp/sse-client/SseClient.Tests/CliTests.cs</c>.
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
        });

        var cfg = Cli.Parse([], env);

        Assert.Equal("feed:9090", cfg.Addr);
        Assert.Equal("reader-token", cfg.Token);
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
}
