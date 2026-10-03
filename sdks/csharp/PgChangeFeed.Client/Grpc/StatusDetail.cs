using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;

namespace PgChangeFeed.Client.Grpc;

/// <summary>
/// Reads the message code out of the <c>grpc-status-details-bin</c> trailer of
/// a failed call. The trailer carries a serialized <c>google.rpc.Status</c>;
/// its repeated field 3 holds <see cref="Any"/> entries, and an entry whose
/// type is <c>google.rpc.ErrorInfo</c> carries the message code as
/// <c>reason</c> (field 1) next to its <c>domain</c> (field 2). The package
/// depends on no <c>google.rpc</c> classes, so this class reads exactly these
/// three levels of fields with the Protobuf runtime that is already a
/// dependency.
/// </summary>
internal static class StatusDetail
{
    private const string TrailerKey = "grpc-status-details-bin";
    private const string ErrorInfoTypeSuffix = "google.rpc.ErrorInfo";
    private const string ServerDomain = "pg-change-feed";

    // Tags: (field number << 3) | wire type 2 (length-delimited).
    private const uint StatusDetailsTag = (3 << 3) | 2;
    private const uint ErrorInfoReasonTag = (1 << 3) | 2;
    private const uint ErrorInfoDomainTag = (2 << 3) | 2;

    /// <summary>
    /// Returns the <c>reason</c> of the first <c>ErrorInfo</c> of domain
    /// <c>pg-change-feed</c> in the status detail of the exception, or
    /// <see langword="null"/>. Never throws.
    /// </summary>
    internal static string? ReadMessageCode(RpcException exception)
    {
        try
        {
            foreach (var entry in exception.Trailers)
            {
                if (entry.IsBinary && string.Equals(entry.Key, TrailerKey, StringComparison.OrdinalIgnoreCase))
                {
                    return ReadFromStatus(entry.ValueBytes);
                }
            }
        }
        catch (Exception)
        {
            return null;
        }

        return null;
    }

    private static string? ReadFromStatus(byte[] status)
    {
        var input = new CodedInputStream(status);
        uint tag;
        while ((tag = input.ReadTag()) != 0)
        {
            if (tag != StatusDetailsTag)
            {
                input.SkipLastField();
                continue;
            }

            var code = ReadFromAny(Any.Parser.ParseFrom(input.ReadBytes()));
            if (code is not null)
            {
                return code;
            }
        }

        return null;
    }

    private static string? ReadFromAny(Any detail)
    {
        if (!detail.TypeUrl.EndsWith(ErrorInfoTypeSuffix, StringComparison.Ordinal))
        {
            return null;
        }

        var input = new CodedInputStream(detail.Value.ToByteArray());
        var reason = string.Empty;
        var domain = string.Empty;
        uint tag;
        while ((tag = input.ReadTag()) != 0)
        {
            switch (tag)
            {
                case ErrorInfoReasonTag:
                    reason = input.ReadString();
                    break;
                case ErrorInfoDomainTag:
                    domain = input.ReadString();
                    break;
                default:
                    input.SkipLastField();
                    break;
            }
        }

        return domain == ServerDomain && reason.Length > 0 ? reason : null;
    }
}
