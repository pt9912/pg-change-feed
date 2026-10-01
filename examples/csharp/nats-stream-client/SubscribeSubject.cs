namespace CdcExamples.NatsStream;

/// <summary>
/// SubscribeSubject leitet das Abonnement-Subjekt aus <c>--source</c> und
/// <c>--target</c> ab (<c>SPEC-024</c>): ohne beide der Wurzel-Wildcard, mit
/// beiden das Zusatz-Subjekt eines Zustellziels
/// <c>cdc.route.&lt;source_id&gt;.&lt;ziel&gt;</c>. Eines der beiden allein ist ein
/// Eingabefehler — das Zusatz-Subjekt hängt an der Quelle und am Ziel zugleich.
/// Ein Token mit Punkt, Platzhalter oder Leerraum wird abgelehnt, bevor ein
/// Abonnement entsteht. Form-Vorbild: <c>SubscribeSubject</c> in
/// <c>examples/nats-stream-client/subject.go</c>.
/// </summary>
public static class SubscribeSubject
{
    public const string All = "cdc.stream.>";

    private static readonly char[] InvalidTokenChars = ['.', '*', '>', ' ', '\t', '\n', '\r'];

    public static string Resolve(string source, string target)
    {
        if (source.Length == 0 && target.Length == 0)
        {
            return All;
        }
        if (source.Length == 0 || target.Length == 0)
        {
            throw new ArgumentException("nats-stream-client: --source und --target gelten nur zusammen");
        }
        Validate("--source", source);
        Validate("--target", target);
        return $"cdc.route.{source}.{target}";
    }

    private static void Validate(string flag, string value)
    {
        if (string.IsNullOrWhiteSpace(value) || value.IndexOfAny(InvalidTokenChars) >= 0)
        {
            throw new ArgumentException(
                $"nats-stream-client: {flag} darf weder leer sein noch '.', '*', '>' oder Leerraum enthalten");
        }
    }
}
