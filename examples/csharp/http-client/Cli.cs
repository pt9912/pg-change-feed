namespace CdcExamples.Http;

/// <summary>
/// Cli liest Flag-Werte und füllt fehlende Felder aus den im Handbuch
/// dokumentierten Umgebungsvariablen (<c>ADR-0076</c> Festlegung 1). Das
/// Umgebungs-Lookup ist injiziert (<paramref name="getEnv"/> statt
/// <c>Environment.GetEnvironmentVariable</c> direkt) — das macht die Funktion
/// rein und netzlos testbar, ohne echte Prozessumgebung zu setzen. Form-Vorbild:
/// <c>examples/http-client/main.go</c>, <c>parseFlags</c>.
/// </summary>
public static class Cli
{
    public static Config Parse(string[] args, Func<string, string?> getEnv)
    {
        var addr = getEnv("CDC_HTTP_ADDR") ?? "";
        var token = getEnv("CDC_API_TOKEN_READER") ?? "";
        var source = "";
        var publication = "";

        for (var i = 0; i < args.Length; i++)
        {
            var (name, inline) = SplitFlag(args[i]);

            string NextValue()
            {
                if (inline is not null)
                {
                    return inline;
                }
                if (i + 1 >= args.Length)
                {
                    throw new ArgumentException($"http-client: Flag {name} braucht einen Wert");
                }
                return args[++i];
            }

            switch (name)
            {
                case "--addr":
                    addr = NextValue();
                    break;
                case "--token":
                    token = NextValue();
                    break;
                case "--source":
                    source = NextValue();
                    break;
                case "--publication":
                    publication = NextValue();
                    break;
                default:
                    throw new ArgumentException($"http-client: unbekanntes Flag {name}");
            }
        }

        return new Config(addr, token, source, publication);
    }

    /// <summary>
    /// SplitFlag zerlegt <c>--name=wert</c> in Name und Inline-Wert; ein Flag
    /// ohne <c>=</c> liefert keinen Inline-Wert — der nächste Token trägt ihn.
    /// </summary>
    private static (string Name, string? InlineValue) SplitFlag(string arg)
    {
        var eq = arg.IndexOf('=');
        return eq < 0 ? (arg, null) : (arg[..eq], arg[(eq + 1)..]);
    }
}
