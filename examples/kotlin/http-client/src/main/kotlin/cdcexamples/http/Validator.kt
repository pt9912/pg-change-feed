package cdcexamples.http

/**
 * Validator prüft eine [Config] gegen die Pflichtfelder des gewählten
 * Verbs — vor jedem Netzwerkaufruf. Jedes Verb braucht die Horch-Adresse
 * und die zu seiner Rechtsklasse passende Token-Variable; die übrigen
 * Felder folgen dem jeweiligen Endpunkt-Vertrag. Form-Vorbild: `validate`
 * in `examples/http-client/main.go`.
 */
object Validator {
    /** knownVerbs trägt die geschlossene Menge der `--verb`-Werte. */
    val knownVerbs: Set<String> = setOf(
        "tables", "changes", "register-consumer", "acknowledge", "consumer-position",
        "remove-consumer", "enable-table", "disable-table", "table-status", "retention-run",
    )

    /**
     * validate liefert eine Fehlerbeschreibung, wenn [cfg] nicht gegen den
     * gewählten Verb-Vertrag genügt, sonst `null`.
     */
    fun validate(cfg: Config): String? {
        if (cfg.verb !in knownVerbs) {
            return "unbekanntes --verb \"${cfg.verb}\""
        }
        if (cfg.addr.isEmpty()) {
            return "keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder --addr) ist noetig, um die Verwaltungs-API zu erreichen"
        }

        return when (cfg.verb) {
            "tables" -> requireReaderToken(cfg) ?: requireFields(
                cfg.source, cfg.publication,
                "--source und --publication sind Pflicht — sie sind die zwei Pflichtfelder von GET /tables",
            )
            "changes" -> requireReaderToken(cfg) ?: requireFields(
                cfg.source,
                "--source ist Pflicht — es ist das einzige Pflichtfeld von GET /changes",
            )
            "register-consumer" -> requireAdminToken(cfg) ?: requireFields(
                cfg.consumerId, cfg.name,
                "--consumer-id und --name sind Pflicht — sie sind die zwei Pflichtfelder von POST /consumers",
            )
            "acknowledge" -> requireAdminToken(cfg) ?: requireFields(
                cfg.consumerId, cfg.source,
                "--consumer-id und --source sind Pflicht — sie sind zwei der drei Pflichtfelder von POST /consumers/acknowledge",
            )
            "consumer-position" -> requireReaderToken(cfg) ?: requireFields(
                cfg.consumerId,
                "--consumer-id ist Pflicht — es ist das einzige Pflichtfeld von GET /consumers/position",
            )
            "remove-consumer" -> requireAdminToken(cfg) ?: requireFields(
                cfg.consumerId,
                "--consumer-id ist Pflicht — es ist das einzige Pflichtfeld von POST /consumers/remove",
            )
            "enable-table" -> requireAdminToken(cfg) ?: requireFields(
                cfg.source, cfg.schema, cfg.table, cfg.publication,
                "--source, --schema, --table und --publication sind Pflicht — table-id und schema-version-id haben einen Default",
            )
            "disable-table" -> requireAdminToken(cfg) ?: requireFields(
                cfg.source, cfg.schema, cfg.table, cfg.publication,
                "--source, --schema, --table und --publication sind Pflicht — sie sind die vier Pflichtfelder von POST /tables/disable",
            )
            "table-status" -> requireReaderToken(cfg) ?: requireFields(
                cfg.source, cfg.schema, cfg.table, cfg.publication,
                "--source, --schema, --table und --publication sind Pflicht — sie sind die vier Pflichtfelder von GET /tables/status",
            )
            "retention-run" -> requireAdminToken(cfg) ?: requireFields(
                cfg.source,
                "--source ist Pflicht — es ist eines der zwei Pflichtfelder von POST /retention/run",
            )
            else -> null
        }
    }

    private fun requireFields(field: String, errorMessage: String): String? =
        if (field.isEmpty()) errorMessage else null

    private fun requireFields(field1: String, field2: String, errorMessage: String): String? =
        if (field1.isEmpty() || field2.isEmpty()) errorMessage else null

    private fun requireFields(field1: String, field2: String, field3: String, field4: String, errorMessage: String): String? =
        if (field1.isEmpty() || field2.isEmpty() || field3.isEmpty() || field4.isEmpty()) errorMessage else null

    private fun requireReaderToken(cfg: Config): String? =
        if (cfg.token.isEmpty()) {
            "kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen"
        } else {
            null
        }

    private fun requireAdminToken(cfg: Config): String? =
        if (cfg.adminToken.isEmpty()) {
            "kein Admin-Token gesetzt — CDC_API_TOKEN_ADMIN (oder --admin-token) ist noetig, um ueber die admin-Rechtsklasse zu schreiben"
        } else {
            null
        }
}
