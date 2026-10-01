package cdcexamples.natsstream

/**
 * SubscribeSubject leitet das Abonnement-Subjekt aus `--source` und
 * `--target` ab (`SPEC-024`): ohne beide der Wurzel-Wildcard, mit beiden das
 * Zusatz-Subjekt eines Zustellziels `cdc.route.<source_id>.<ziel>`. Eines der
 * beiden allein ist ein Eingabefehler — das Zusatz-Subjekt hängt an der Quelle
 * und am Ziel zugleich. Ein Token mit Punkt, Platzhalter oder Leerraum wird
 * abgelehnt, bevor ein Abonnement entsteht. Form-Vorbild: `SubscribeSubject`
 * in `examples/nats-stream-client/subject.go`.
 */
object SubscribeSubject {
    const val ALL: String = "cdc.stream.>"

    private val invalidTokenChars = charArrayOf('.', '*', '>', ' ', '\t', '\n', '\r')

    fun resolve(source: String, target: String): String {
        if (source.isEmpty() && target.isEmpty()) {
            return ALL
        }
        require(source.isNotEmpty() && target.isNotEmpty()) {
            "nats-stream-client: --source und --target gelten nur zusammen"
        }
        validate("--source", source)
        validate("--target", target)
        return "cdc.route.$source.$target"
    }

    private fun validate(flag: String, value: String) {
        require(value.isNotBlank() && value.none { it in invalidTokenChars }) {
            "nats-stream-client: $flag darf weder leer sein noch '.', '*', '>' oder Leerraum enthalten"
        }
    }
}
