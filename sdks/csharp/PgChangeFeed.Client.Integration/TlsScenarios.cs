using System.Security.Authentication;
using Grpc.Core;
using PgChangeFeed.Client;
using Xunit;

namespace PgChangeFeed.Client.Integration;

/// <summary>
/// The three TLS refusals every TLS phase proves against the running feed
/// container: no trust anchor (the certificate of the server is unknown to the
/// system), a foreign trust anchor, and a server name that is not in the
/// certificate (the container name, which resolves in the Docker network but
/// is not listed in the certificate; the anchor is the right one). Each must
/// end with a failed TLS verification — a refusal for any other reason (an
/// unreachable name, a refused port) does not count.
/// </summary>
internal static class TlsScenarios
{
    internal static Uri WithMismatchHost(Uri address)
        => new UriBuilder(address) { Host = PhaseEnvironment.TlsMismatchHost }.Uri;

    /// <summary>
    /// Runs <paramref name="call"/> with the three option sets and asserts a TLS
    /// verification failure for each; prints the runner marker afterwards.
    /// </summary>
    internal static async Task ExpectRefusalsAsync(Uri address, string token, Func<PgChangeFeedClientOptions, Task> call)
    {
        var cases = new (string Name, PgChangeFeedClientOptions Options)[]
        {
            ("no_anchor", new PgChangeFeedClientOptions(address, token)),
            ("foreign_anchor", new PgChangeFeedClientOptions(address, token, PhaseEnvironment.TlsForeignCaFile)),
            ("name_mismatch", new PgChangeFeedClientOptions(WithMismatchHost(address), token, PhaseEnvironment.TlsCaFile)),
        };
        foreach (var (name, options) in cases)
        {
            var error = await Record.ExceptionAsync(() => call(options));
            Assert.True(error is not null, $"{name}: the call succeeded");
            Assert.True(IsVerificationFailure(error), $"{name}: not a TLS verification failure: {error}");
        }

        PhaseEnvironment.Print("REJECTED tls no_anchor=ok foreign_anchor=ok name_mismatch=ok");
    }

    /// <summary>True when the exception chain holds the <see cref="AuthenticationException"/> of a failed TLS check.</summary>
    internal static bool IsVerificationFailure(Exception? error)
    {
        for (var current = error; current is not null; current = current.InnerException)
        {
            if (current is AuthenticationException)
            {
                return true;
            }

            if (current is RpcException rpc && IsVerificationFailure(rpc.Status.DebugException))
            {
                return true;
            }
        }

        return false;
    }
}
