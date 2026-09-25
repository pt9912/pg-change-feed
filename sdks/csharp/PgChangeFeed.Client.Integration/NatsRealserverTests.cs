using PgChangeFeed.Client;
using PgChangeFeed.Client.Nats;
using PgChangeFeed.Client.Nats.Models;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// Real-server phase for the NATS stream client: subscribes the source's
/// namespace against the running feed container and
/// receives a change committed afterwards as a full JSON event; a second
/// connect with a wrong token is rejected by the NATS server, surfacing as a
/// NATS error carrying the server's
/// <c>Authorization Violation</c> wording — not swallowed as an empty
/// stream.
/// </summary>
public sealed class NatsRealserverTests
{
    [Fact]
    public async Task ReceivesACommittedChangeOverTheStream()
    {
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.NatsUrl), PhaseEnvironment.NatsStreamToken);
        await using var client = new PgChangeFeedNatsStreamClient(options);
        var subject = PgChangeFeedNatsStreamClient.BuildSourceSubject(PhaseEnvironment.SourceId);
        PhaseEnvironment.Print("READY");

        var received = await ReceiveSentinelAsync(
            client.StreamChangesAsync(subject, PhaseEnvironment.ReceiveCts.Token));

        Assert.NotEqual(string.Empty, received.ChangeId);
        Assert.NotEqual(string.Empty, received.TransactionId);
        Assert.NotEqual(string.Empty, received.SourceTableId);
        Assert.NotEqual(string.Empty, received.SchemaVersion);
        Assert.Equal("INSERT", received.Operation);
        Assert.Null(received.OldImage);
        Assert.Equal(PhaseEnvironment.Table, received.Table);
        Assert.Equal("public", received.Schema);
        Assert.Contains(PhaseEnvironment.Sentinel, received.NewImage?.ToString() ?? string.Empty);

        PhaseEnvironment.Print(
            $"RECEIVED change_id={received.ChangeId} table={received.Table} " +
            $"operation={received.Operation} new_image={received.NewImage}");
    }

    [Fact]
    public async Task ConnectWithWrongTokenIsRejectedByTheServer()
    {
        var options = new PgChangeFeedClientOptions(new Uri(PhaseEnvironment.NatsUrl), "wrong-token");
        await using var client = new PgChangeFeedNatsStreamClient(options);

        var exception = await Assert.ThrowsAnyAsync<Exception>(async () =>
        {
            await foreach (var _ in client.StreamChangesAsync(
                PgChangeFeedNatsStreamClient.BuildSourceSubject(PhaseEnvironment.SourceId),
                PhaseEnvironment.RejectCts.Token)) { }
        });

        Assert.NotEqual(typeof(TimeoutException), exception.GetType());
        // Die Ablehnungsursache (Server-Rohwortlaut "-ERR Authorization
        // Violation") traegt die Ausnahme-Kette (ToString schliesst die
        // inneren Meldungen ein); der NATS-.NET-Wrapper haellt sie in eine
        // eigene Connect-Meldung ("can not start to connect nats server: ...").
        Assert.Contains("Authorization", exception.ToString());
        PhaseEnvironment.Print($"REJECTED token-rejected: {exception.ToString()}");
    }

    private static async Task<Change> ReceiveSentinelAsync(IAsyncEnumerable<Change> stream)
    {
        var deadline = DateTime.UtcNow + TimeSpan.FromSeconds(90);
        await foreach (var change in stream)
        {
            if (change.Table == PhaseEnvironment.Table &&
                change.Operation == "INSERT" &&
                (change.NewImage?.ToString() ?? string.Empty).Contains(PhaseEnvironment.Sentinel))
            {
                return change;
            }
            Assert.True(DateTime.UtcNow < deadline,
                $"kein Event mit dem Sentinel {PhaseEnvironment.Sentinel} innerhalb der Frist empfangen");
        }
        throw new InvalidOperationException("Der Stream endete, bevor der Sentinel empfangen wurde.");
    }
}