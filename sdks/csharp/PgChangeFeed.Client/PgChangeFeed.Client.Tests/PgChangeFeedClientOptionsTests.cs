using Xunit;

namespace PgChangeFeed.Client.Tests;

/// <summary>
/// Construction/validation tests for the shared options class — no wire
/// behavior (that arrives with the surface-specific follow-up slices).
/// </summary>
public class PgChangeFeedClientOptionsTests
{
    private static readonly Uri ValidAddress = new("https://example.invalid:8443");

    [Fact]
    public void Constructor_WithValidArguments_SetsProperties()
    {
        var options = new PgChangeFeedClientOptions(ValidAddress, "test-token");

        Assert.Equal(ValidAddress, options.Address);
        Assert.Equal("test-token", options.ApiToken);
    }

    [Fact]
    public void Constructor_WithNullAddress_Throws()
    {
        Assert.Throws<ArgumentNullException>(() => new PgChangeFeedClientOptions(null!, "test-token"));
    }

    [Theory]
    [InlineData("")]
    [InlineData("   ")]
    public void Constructor_WithEmptyOrWhitespaceToken_Throws(string apiToken)
    {
        Assert.Throws<ArgumentException>(() => new PgChangeFeedClientOptions(ValidAddress, apiToken));
    }

    [Fact]
    public void Constructor_WithNullToken_Throws()
    {
        Assert.Throws<ArgumentException>(() => new PgChangeFeedClientOptions(ValidAddress, null!));
    }
}
