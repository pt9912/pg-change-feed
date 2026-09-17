package cdcexamples.http

import java.net.http.HttpClient
import java.time.Duration
import kotlin.system.exitProcess

/**
 * Command http-client ist ein öffentliches Beispiel für den
 * Anfrage/Antwort-Zugriff über die HTTP-/JSON-API (`LH-FA-SST-006`,
 * `ADR-0057`, `ADR-0090`): es ruft den `reader`-Endpunkt `GET /tables` real
 * gegen den laufenden Feed-Container auf und gibt die Antwort aus. Startform
 * ist ein Container-Aufruf, kein Host-Aufruf (`ADR-0087` Festlegung 3);
 * Adresse und Token kommen aus `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER` und
 * lassen sich per Flag übersteuern (`--addr`, `--token`, `--source`,
 * `--publication`).
 *
 * Form-Vorbild: `examples/http-client` (Go), `examples/csharp/http-client`
 * (C#). Dieses Programm trägt keine Zustandsmaschine: es stellt eine Anfrage
 * und endet.
 */
private val connectTimeout: Duration = Duration.ofSeconds(10)

fun main(args: Array<String>) {
    val cfg = try {
        Cli.parse(args) { name -> System.getenv(name) }
    } catch (ex: IllegalArgumentException) {
        System.err.println(ex.message)
        exitProcess(2)
    }

    if (cfg.addr.isEmpty()) {
        System.err.println(
            "http-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder --addr) ist noetig, um die Verwaltungs-API zu erreichen",
        )
        exitProcess(2)
    }
    if (cfg.token.isEmpty()) {
        System.err.println(
            "http-client: kein Token gesetzt — CDC_API_TOKEN_READER (oder --token) ist noetig, um ueber die reader-Rechtsklasse zu lesen",
        )
        exitProcess(2)
    }
    if (cfg.source.isEmpty() || cfg.publication.isEmpty()) {
        System.err.println(
            "http-client: --source und --publication sind Pflicht — sie sind die zwei Pflichtfelder von GET /tables",
        )
        exitProcess(2)
    }

    val httpClient = HttpClient.newBuilder()
        .connectTimeout(connectTimeout)
        .build()

    try {
        val body = TablesClient.listTables(httpClient, cfg)
        println(body)
    } catch (ex: Exception) {
        // `ex.message` ist bei manchen JDK-Netzwerk-Ausnahmen (z. B.
        // `ConnectException` ohne Text) `null` — `ex` selbst trägt in jedem
        // Fall die Ausnahme-Klasse und, wo vorhanden, den Text.
        System.err.println("http-client: GET /tables fehlgeschlagen: $ex")
        exitProcess(1)
    }
}
