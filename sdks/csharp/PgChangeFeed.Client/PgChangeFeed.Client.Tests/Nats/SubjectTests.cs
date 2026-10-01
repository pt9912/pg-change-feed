using PgChangeFeed.Client.Nats;
using Xunit;

namespace PgChangeFeed.Client.Tests.Nats;

/// <summary>
/// Subject formatting for the NATS stream namespace — network-free,
/// pure string construction, no NATS connection involved. Covers the two
/// builder helpers (<see cref="PgChangeFeedNatsStreamClient.BuildSubject"/>,
/// <see cref="PgChangeFeedNatsStreamClient.BuildSourceSubject"/>) and the
/// <see cref="PgChangeFeedNatsStreamClient.AllSourcesSubject"/> constant.
/// </summary>
public class SubjectTests
{
    [Fact]
    public void AllSourcesSubject_IsRootWildcard()
    {
        Assert.Equal("cdc.stream.>", PgChangeFeedNatsStreamClient.AllSourcesSubject);
    }

    [Fact]
    public void BuildSubject_JoinsAllFourTokens()
    {
        var subject = PgChangeFeedNatsStreamClient.BuildSubject("source-1", "public", "orders");

        Assert.Equal("cdc.stream.source-1.public.orders", subject);
    }

    [Fact]
    public void BuildSourceSubject_JoinsSourceThenWildcard()
    {
        var subject = PgChangeFeedNatsStreamClient.BuildSourceSubject("source-1");

        Assert.Equal("cdc.stream.source-1.>", subject);
    }

    [Theory]
    [InlineData("")]
    [InlineData("   ")]
    [InlineData(null)]
    public void BuildSubject_RejectsBlankSourceId(string? sourceId)
    {
        Assert.Throws<ArgumentException>(() => PgChangeFeedNatsStreamClient.BuildSubject(sourceId!, "public", "orders"));
    }

    [Theory]
    [InlineData(".")]
    [InlineData("*")]
    [InlineData(">")]
    [InlineData("has.dot")]
    [InlineData("has space")]
    public void BuildSubject_RejectsTokenSeparatorsAndWildcards(string invalidToken)
    {
        Assert.Throws<ArgumentException>(() => PgChangeFeedNatsStreamClient.BuildSubject(invalidToken, "public", "orders"));
        Assert.Throws<ArgumentException>(() => PgChangeFeedNatsStreamClient.BuildSubject("source-1", invalidToken, "orders"));
        Assert.Throws<ArgumentException>(() => PgChangeFeedNatsStreamClient.BuildSubject("source-1", "public", invalidToken));
    }

    [Fact]
    public void BuildTargetSubject_JoinsRouteSourceAndTarget()
    {
        Assert.Equal("cdc.route.source-1.eu", PgChangeFeedNatsStreamClient.BuildTargetSubject("source-1", "eu"));
    }

    [Fact]
    public void BuildSourceTargetsSubject_JoinsSourceThenWildcard()
    {
        Assert.Equal("cdc.route.source-1.>", PgChangeFeedNatsStreamClient.BuildSourceTargetsSubject("source-1"));
    }

    [Theory]
    [InlineData("")]
    [InlineData("   ")]
    [InlineData(".")]
    [InlineData("*")]
    [InlineData(">")]
    [InlineData("has.dot")]
    [InlineData("has space")]
    public void BuildTargetSubject_RejectsBlankTokensSeparatorsAndWildcards(string invalidToken)
    {
        Assert.Throws<ArgumentException>(() => PgChangeFeedNatsStreamClient.BuildTargetSubject("source-1", invalidToken));
        Assert.Throws<ArgumentException>(() => PgChangeFeedNatsStreamClient.BuildTargetSubject(invalidToken, "eu"));
        Assert.Throws<ArgumentException>(() => PgChangeFeedNatsStreamClient.BuildSourceTargetsSubject(invalidToken));
    }

    [Fact]
    public async Task StreamChangesAsync_SubscribesToTheTargetSubject()
    {
        var fake = FakeNatsClient.WithPayloads();
        var client = new PgChangeFeedNatsStreamClient(fake);

        await foreach (var _ in client.StreamChangesAsync(
            PgChangeFeedNatsStreamClient.BuildTargetSubject("source-1", "eu")))
        {
        }

        Assert.Equal("cdc.route.source-1.eu", fake.LastSubject);
    }

    [Fact]
    public void BuildSourceSubject_RejectsBlankOrInvalidSourceId()
    {
        Assert.Throws<ArgumentException>(() => PgChangeFeedNatsStreamClient.BuildSourceSubject(""));
        Assert.Throws<ArgumentException>(() => PgChangeFeedNatsStreamClient.BuildSourceSubject("has.dot"));
    }
}
