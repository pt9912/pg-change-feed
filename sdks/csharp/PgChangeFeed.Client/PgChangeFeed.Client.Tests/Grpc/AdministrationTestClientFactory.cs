using PgChangeFeed.Client.Grpc;

namespace PgChangeFeed.Client.Tests.Grpc;

/// <summary>
/// Builds a <see cref="PgChangeFeedAdministrationClient"/> wired to a
/// <see cref="FakeUnaryCallInvoker"/> — every test in this project stays
/// network-free (no real server, no real socket).
/// </summary>
internal static class AdministrationTestClientFactory
{
    public static readonly Uri ServerAddress = new("http://localhost:50051");

    public static PgChangeFeedAdministrationClient Create(FakeUnaryCallInvoker invoker, string apiToken = "test-token")
        => new(invoker, new PgChangeFeedClientOptions(ServerAddress, apiToken));
}
