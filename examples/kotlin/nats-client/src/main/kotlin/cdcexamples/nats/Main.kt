package cdcexamples.nats

import io.nats.client.Message
import io.nats.client.Nats
import java.net.http.HttpClient
import java.time.Duration
import kotlin.system.exitProcess

/**
 * Command nats-client ist ein öffentliches Beispiel für den zweiseitigen
 * Zugriffsweg über das NATS-Wecksignal (`LH-FA-SST-007`, `ADR-0079`,
 * `ADR-0090`): es abonniert das tabellen-granulare Subjekt
 * `cdc.changes.<source_id>.<schema>.<table>` (`SPEC-017`, `ADR-0056`) und
 * holt beim Weckruf die Änderung selbst über die HTTP-/JSON-API
 * (`LH-FA-SST-006`) — das Signal trägt per Vertrag einen leeren Payload
 * (`ADR-0055`), es sagt nur „lies erneut über den bestehenden Zugriffsweg".
 * Startform ist ein Container-Aufruf, kein Host-Aufruf (`ADR-0087`
 * Festlegung 3); NATS-URL, Adresse und Token kommen aus `CDC_NATS_URL`/
 * `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER` und lassen sich per Flag
 * übersteuern (`--nats-url`, `--addr`, `--token`, `--source`, `--schema`,
 * `--table`).
 *
 * Form-Vorbild: `examples/nats-client` (Go), `examples/csharp/nats-client`
 * (C#). Dieses Programm trägt keine Zustandsmaschine: eine
 * Reconnect-/Dedup-/Rückstand-Logik ist bewusst nicht seine Aufgabe
 * (`SPEC-023`). Das Warten auf das Wecksignal benutzt
 * [io.nats.client.Subscription.nextMessage] mit [Duration.ZERO] — jnats
 * wartet damit real unbegrenzt (`MessageQueueBase#poll`,
 * `io.nats:jnats`-Quelltext, `ADR-0090` Festlegung 4). **Keine** dieselbe
 * Semantik wie `sub.NextMsg(0)` im Go-Vorbild: `nats.go` startet bei einem
 * Timeout-Wert von `0` einen sofort ablaufenden Timer statt unbegrenzt zu
 * warten — das Go-Beispiel braucht deshalb einen expliziten, endlichen
 * `waitTimeout` (siehe dortigen Kommentar). Beide Beispiele erreichen im
 * Ergebnis ein praktisch unbegrenztes Warten, aber über unterschiedliche
 * Bibliotheks-Semantik, nicht über identische Aufrufformen.
 */
fun main(args: Array<String>) {
    val cfg = try {
        Cli.parse(args) { name -> System.getenv(name) }
    } catch (ex: IllegalArgumentException) {
        System.err.println(ex.message)
        exitProcess(2)
    }

    // `CDC_HTTP_ADDR` ungesetzt heißt: die HTTP-API ist deaktiviert
    // (`SPEC-018`). Ohne sie kann der Weckruf keine Änderung holen — das
    // Beispiel scheitert sichtbar, statt still nichts zu tun (`ADR-0079`
    // Festlegung 2).
    if (cfg.addr.isEmpty()) {
        System.err.println(
            "nats-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder --addr) ist noetig, um beim Weckruf die Aenderung zu holen",
        )
        exitProcess(2)
    }
    if (cfg.natsUrl.isEmpty()) {
        System.err.println(
            "nats-client: keine NATS-URL gesetzt — CDC_NATS_URL (oder --nats-url) ist noetig, um das Wecksignal zu abonnieren",
        )
        exitProcess(2)
    }
    if (cfg.source.isEmpty() || cfg.schema.isEmpty() || cfg.table.isEmpty()) {
        System.err.println("nats-client: --source, --schema und --table sind Pflicht")
        exitProcess(2)
    }
    if (cfg.token.isEmpty()) {
        System.err.println(
            "nats-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen",
        )
        exitProcess(2)
    }

    val subject = Subject.build(cfg.source, cfg.schema, cfg.table)

    val nc = try {
        Nats.connect(cfg.natsUrl)
    } catch (ex: Exception) {
        System.err.println("nats-client: Verbindung zu NATS (${cfg.natsUrl}) fehlgeschlagen: $ex")
        exitProcess(1)
    }

    val sub = try {
        nc.subscribe(subject)
    } catch (ex: Exception) {
        System.err.println("nats-client: Abonnement auf $subject fehlgeschlagen: $ex")
        exitProcess(1)
    }

    println(
        "nats-client: lauscht auf $subject — beim naechsten Weckruf wird die Aenderung ueber ${cfg.addr} geholt",
    )

    val msg: Message? = try {
        sub.nextMessage(Duration.ZERO)
    } catch (ex: Exception) {
        System.err.println("nats-client: Empfang des Wecksignals fehlgeschlagen: $ex")
        exitProcess(1)
    }
    if (msg == null) {
        System.err.println("nats-client: Abonnement endete ohne Wecksignal")
        exitProcess(1)
    }

    // Der Payload trägt keine Änderungsdaten (`SPEC-017`); die Änderung
    // kommt ausschließlich über den HTTP-Lesezugriff.
    println("nats-client: Weckruf auf ${msg.subject} (Payload ${msg.data.size} Byte) — hole die Aenderung ueber HTTP")

    val httpClient = HttpClient.newBuilder().build()
    try {
        val body = ChangesClient.fetch(httpClient, cfg)
        println(body)
    } catch (ex: Exception) {
        System.err.println("nats-client: HTTP-Abfrage fehlgeschlagen: $ex")
        exitProcess(1)
    } finally {
        nc.close()
    }
}
