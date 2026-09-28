using PgChangeFeed.Client.Grpc;
using Xunit;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>Constructor argument validation for both constructors.</summary>
public class PgChangeFeedAdministrationClientConstructionTests
{
    [Fact]
    public void OwnedChannelConstructor_WithNullOptions_Throws()
    {
        Assert.Throws<ArgumentNullException>(() => new PgChangeFeedAdministrationClient(null!));
    }

    [Fact]
    public void InjectedInvokerConstructor_WithNullCallInvoker_Throws()
    {
        var options = new PgChangeFeedClientOptions(AdministrationTestClientFactory.ServerAddress, "token");

        Assert.Throws<ArgumentNullException>(() => new PgChangeFeedAdministrationClient(null!, options));
    }

    [Fact]
    public void InjectedInvokerConstructor_WithNullOptions_Throws()
    {
        var invoker = FakeUnaryCallInvoker.WithResponse<object?>(null);

        Assert.Throws<ArgumentNullException>(() => new PgChangeFeedAdministrationClient(invoker, null!));
    }
}
