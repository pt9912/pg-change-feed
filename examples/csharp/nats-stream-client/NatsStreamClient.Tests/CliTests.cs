using Xunit;

namespace CdcExamples.NatsStream.Tests;

/// <summary>
/// Prüft <see cref="Cli.Parse"/>: Umgebungs-Default, Flag-Override, beide
/// Flag-Formen (<c>--name=wert</c> und <c>--name wert</c>) und die
/// Fehlerfälle (unbekanntes Flag, fehlender Wert). Form-Vorbild:
/// <c>examples/csharp/grpc-client/GrpcClient.Tests/CliTests.cs</c>.
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
            ["CDC_NATS_URL"] = "nats://feed:4222",
            ["CDC_NATS_STREAM_TOKEN"] = "stream-token",
        });

        var cfg = Cli.Parse([], env);

        Assert.Equal("nats://feed:4222", cfg.NatsUrl);
        Assert.Equal("stream-token", cfg.Token);
    }

    [Fact]
    public void Parse_FlagOverridesEnvironment()
    {
        var env = EnvFrom(new Dictionary<string, string>
        {
            ["CDC_NATS_URL"] = "nats://feed:4222",
            ["CDC_NATS_STREAM_TOKEN"] = "stream-token",
        });

        var cfg = Cli.Parse(["--nats-url", "nats://override:4222", "--token=other-token"], env);

        Assert.Equal("nats://override:4222", cfg.NatsUrl);
        Assert.Equal("other-token", cfg.Token);
    }

    [Fact]
    public void Parse_NoEnvironmentNoFlags_ReturnsEmpty()
    {
        var cfg = Cli.Parse([], NoEnv);

        Assert.Equal("", cfg.NatsUrl);
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
        var ex = Assert.Throws<ArgumentException>(() => Cli.Parse(["--nats-url"], NoEnv));
        Assert.Contains("braucht einen Wert", ex.Message);
    }
}
