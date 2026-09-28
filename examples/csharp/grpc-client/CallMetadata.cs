using Grpc.Core;

namespace CdcExamples.Grpc;

/// <summary>
/// CallMetadata baut den <c>authorization</c>-Metadata-Eintrag, den der
/// Stream-Aufruf und alle elf Administration-RPCs teilen (<c>SPEC-020</c>) —
/// derselbe <c>Bearer</c>-Vorsprung wie der Auth-Interceptor des Adapters.
/// <see cref="RequestTimeout"/> begrenzt einen einzelnen unären RPC-Aufruf;
/// der Stream trägt keine eigene Frist. Form-Vorbild:
/// <c>authorizationMetadataKey</c>/<c>bearerPrefix</c>/<c>callCtx</c> in
/// <c>examples/grpc-client/main.go</c>.
/// </summary>
public static class CallMetadata
{
    private const string AuthorizationMetadataKey = "authorization";
    private const string BearerPrefix = "Bearer ";

    public static readonly TimeSpan RequestTimeout = TimeSpan.FromSeconds(10);

    public static Metadata Headers(string token) => new() { { AuthorizationMetadataKey, BearerPrefix + token } };

    public static DateTime Deadline() => DateTime.UtcNow.Add(RequestTimeout);
}
