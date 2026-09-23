/// <summary>
/// Shared environment plumbing for the realserver integration phases
/// (ADR-0110 §Entscheidung Festlegung 2, mirrored to the C# tree): every
/// phase names its test class explicitly through the runner's
/// <c>PGCHANGEFEED_TEST_NAME</c> filter (no silent exclusion), and every
/// marker (<c>READY</c>/<c>RECEIVED</c>/<c>REJECTED</c>) is what
/// <c>tools/harness/run-sdk-csharp-integration-tests.sh</c> asserts on via
/// <c>docker logs</c> while the test still runs — so every print is flushed
/// immediately.
/// </summary>
internal static class PhaseEnvironment
{
    private static string Required(string name)
        => Environment.GetEnvironmentVariable(name)
           ?? throw new InvalidOperationException(
               $"Umgebungsvariable {name} fehlt — die Runner-Phase trägt sie.");

    /// <summary>Prints a runner marker and flushes it immediately (the
    /// runner polls <c>docker logs</c> while the test is still running).</summary>
    internal static void Print(string marker)
    {
        Console.WriteLine(marker);
        Console.Out.Flush();
    }

    internal static string Table => Required("PGCHANGEFEED_E2E_TABLE");
    internal static string Sentinel => Required("PGCHANGEFEED_E2E_SENTINEL");
    internal static string ApiToken => Required("PGCHANGEFEED_API_TOKEN");
    internal static string AdminToken => Required("PGCHANGEFEED_API_TOKEN_ADMIN");
    internal static string ReaderToken => Required("PGCHANGEFEED_API_TOKEN_READER");
    internal static string GrpcAddr => Required("PGCHANGEFEED_GRPC_ADDR");
    internal static string HttpAddr => Required("PGCHANGEFEED_HTTP_ADDR");
    internal static string NatsUrl => Required("PGCHANGEFEED_NATS_URL");
    internal static string NatsStreamToken => Required("PGCHANGEFEED_NATS_STREAM_TOKEN");
    internal static string SourceId => Required("PGCHANGEFEED_SOURCE_ID");
    internal static string HttpPublication => Required("PGCHANGEFEED_HTTP_PUBLICATION");

    internal static CancellationTokenSource ReceiveCts => new(TimeSpan.FromSeconds(90));
    internal static CancellationTokenSource RejectCts => new(TimeSpan.FromSeconds(15));
}