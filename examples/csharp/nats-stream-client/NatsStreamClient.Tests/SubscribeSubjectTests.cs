using Xunit;

namespace CdcExamples.NatsStream.Tests;

/// <summary>
/// Prüft die Ableitung des Abonnement-Subjekts aus <c>--source</c> und
/// <c>--target</c> — netzlos, reine Funktion. Form-Vorbild:
/// <c>examples/nats-stream-client/subject_test.go</c>.
/// </summary>
public class SubscribeSubjectTests
{
    [Fact]
    public void Resolve_WithoutSourceAndTarget_IsTheRootWildcard()
    {
        Assert.Equal("cdc.stream.>", SubscribeSubject.Resolve("", ""));
    }

    [Fact]
    public void Resolve_SourceAndTarget_IsTheRouteSubject()
    {
        Assert.Equal("cdc.route.quelle-1.eu", SubscribeSubject.Resolve("quelle-1", "eu"));
    }

    [Theory]
    [InlineData("quelle-1", "")]
    [InlineData("", "eu")]
    [InlineData("quelle-1", "   ")]
    [InlineData("quelle-1", "a.b")]
    [InlineData("quelle-1", "*")]
    [InlineData("quelle-1", ">")]
    [InlineData("quelle-1", "a b")]
    [InlineData("a.b", "eu")]
    [InlineData("*", "eu")]
    [InlineData(">", "eu")]
    public void Resolve_HalfOrInvalidInput_Throws(string source, string target)
    {
        Assert.Throws<ArgumentException>(() => SubscribeSubject.Resolve(source, target));
    }

    [Fact]
    public void Parse_SourceAndTargetFlagsFillTheConfig()
    {
        var cfg = Cli.Parse(["--source=quelle-1", "--target", "eu"], _ => null);

        Assert.Equal("quelle-1", cfg.Source);
        Assert.Equal("eu", cfg.Target);
        Assert.Equal("", Cli.Parse([], _ => null).Target);
    }
}
