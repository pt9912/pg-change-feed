package cdcexamples.grpc

/**
 * Cli liest Flag-Werte und füllt fehlende Felder aus den im Handbuch
 * dokumentierten Umgebungsvariablen (`ADR-0076` Festlegung 1). Das
 * Umgebungs-Lookup ist injiziert ([getEnv] statt `System.getenv` direkt) —
 * das macht die Funktion rein und netzlos testbar, ohne echte
 * Prozessumgebung zu setzen. Form-Vorbild:
 * `examples/csharp/grpc-client/Cli.cs` (`slice-102`),
 * `examples/kotlin/http-client/Cli.kt`.
 */
object Cli {
    fun parse(args: Array<String>, getEnv: (String) -> String?): Config {
        var addr = getEnv("CDC_GRPC_ADDR") ?: ""
        var token = getEnv("CDC_API_TOKEN_READER") ?: ""

        var i = 0
        while (i < args.size) {
            val (name, inline) = splitFlag(args[i])

            fun nextValue(): String {
                if (inline != null) {
                    return inline
                }
                if (i + 1 >= args.size) {
                    throw IllegalArgumentException("grpc-client: Flag $name braucht einen Wert")
                }
                i += 1
                return args[i]
            }

            when (name) {
                "--addr" -> addr = nextValue()
                "--token" -> token = nextValue()
                else -> throw IllegalArgumentException("grpc-client: unbekanntes Flag $name")
            }
            i += 1
        }

        return Config(addr, token)
    }

    /**
     * splitFlag zerlegt `--name=wert` in Name und Inline-Wert; ein Flag ohne
     * `=` liefert keinen Inline-Wert — der nächste Token trägt ihn.
     */
    private fun splitFlag(arg: String): Pair<String, String?> {
        val eq = arg.indexOf('=')
        return if (eq < 0) arg to null else arg.substring(0, eq) to arg.substring(eq + 1)
    }
}
