package cdcexamples.grpc

/**
 * Cli liest Flag-Werte und füllt fehlende Felder aus den im Handbuch
 * dokumentierten Umgebungsvariablen (`ADR-0076` Festlegung 1). Das
 * Umgebungs-Lookup ist injiziert ([getEnv] statt `System.getenv` direkt) —
 * das macht die Funktion rein und netzlos testbar, ohne echte
 * Prozessumgebung zu setzen. Diese Funktion prüft nur die Flag-Syntax; ob
 * ein Verb bekannt ist oder seine Pflichtfelder trägt, prüft [Validator].
 * Form-Vorbild: `examples/csharp/grpc-client/Cli.cs`,
 * `examples/kotlin/http-client/Cli.kt`, `parseFlags` in
 * `examples/grpc-client/main.go`.
 */
object Cli {
    fun parse(args: Array<String>, getEnv: (String) -> String?): Config {
        var addr = getEnv("CDC_GRPC_ADDR") ?: ""
        var token = getEnv("CDC_API_TOKEN_READER") ?: ""
        var adminToken = getEnv("CDC_API_TOKEN_ADMIN") ?: ""
        var verb = "stream"
        var schema = ""
        var table = ""
        var target = ""
        var consumerId = ""
        var name = ""
        var offset = 0L
        var tableId = ""
        var schemaVersionId = ""
        var version = 1L
        var source = ""
        var publication = ""
        var from = 0L
        var to = 0L
        var limit = 0L
        var minAgeNanos = 0L

        var i = 0
        while (i < args.size) {
            val (flagName, inline) = splitFlag(args[i])

            fun nextValue(): String {
                if (inline != null) {
                    return inline
                }
                if (i + 1 >= args.size) {
                    throw IllegalArgumentException("grpc-client: Flag $flagName braucht einen Wert")
                }
                i += 1
                return args[i]
            }

            when (flagName) {
                "--addr" -> addr = nextValue()
                "--token" -> token = nextValue()
                "--admin-token" -> adminToken = nextValue()
                "--verb" -> verb = nextValue()
                "--schema" -> schema = nextValue()
                "--table" -> table = nextValue()
                "--target" -> target = nextValue()
                "--consumer-id" -> consumerId = nextValue()
                "--name" -> name = nextValue()
                "--offset" -> offset = parseNonNegativeLong(flagName, nextValue())
                "--table-id" -> tableId = nextValue()
                "--schema-version-id" -> schemaVersionId = nextValue()
                "--version" -> version = parseLong(flagName, nextValue())
                "--source" -> source = nextValue()
                "--publication" -> publication = nextValue()
                "--from" -> from = parseNonNegativeLong(flagName, nextValue())
                "--to" -> to = parseNonNegativeLong(flagName, nextValue())
                "--limit" -> limit = parseLong(flagName, nextValue())
                "--min-age-nanos" -> minAgeNanos = parseLong(flagName, nextValue())
                else -> throw IllegalArgumentException("grpc-client: unbekanntes Flag $flagName")
            }
            i += 1
        }

        return Config(
            addr = addr, token = token, adminToken = adminToken, verb = verb,
            schema = schema, table = table,
            consumerId = consumerId, name = name, offset = offset,
            tableId = tableId, schemaVersionId = schemaVersionId, version = version,
            source = source, publication = publication,
            from = from, to = to, limit = limit, minAgeNanos = minAgeNanos,
            target = target,
        )
    }

    private fun parseLong(flagName: String, value: String): Long =
        value.toLongOrNull()
            ?: throw IllegalArgumentException("grpc-client: Flag $flagName braucht eine Ganzzahl, erhielt \"$value\"")

    private fun parseNonNegativeLong(flagName: String, value: String): Long {
        val parsed = value.toLongOrNull()
        if (parsed == null || parsed < 0) {
            throw IllegalArgumentException("grpc-client: Flag $flagName braucht eine nicht-negative Ganzzahl, erhielt \"$value\"")
        }
        return parsed
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
