package cdcexamples.http

import java.net.http.HttpClient
import java.time.Duration
import kotlin.system.exitProcess

/**
 * Command http-client ist ein öffentliches Beispiel für den
 * Anfrage/Antwort-Zugriff über die HTTP-/JSON-API (`LH-FA-SST-006`): das
 * `--verb`-Flag ruft eine von zehn dokumentierten Fähigkeiten real gegen
 * den laufenden Feed-Container auf und gibt die Antwort aus. Startform ist
 * ein Container-Aufruf, kein Host-Aufruf (`ADR-0087` Festlegung 3);
 * Default-Verb ist `tables` — die bestehende Startform
 * `ARGS="--source <quelle> --publication <publication>"` bleibt
 * unverändert funktionsfähig.
 *
 * Form-Vorbild: `examples/http-client` (Go), `examples/csharp/http-client`
 * (C#). Dieses Programm trägt keine Zustandsmaschine: es stellt eine
 * Anfrage und endet.
 */
private val connectTimeout: Duration = Duration.ofSeconds(10)

fun main(args: Array<String>) {
    val cfg = try {
        Cli.parse(args) { name -> System.getenv(name) }
    } catch (ex: IllegalArgumentException) {
        System.err.println(ex.message)
        exitProcess(2)
    }

    val validationError = Validator.validate(cfg)
    if (validationError != null) {
        System.err.println("http-client: $validationError")
        exitProcess(2)
    }

    val httpClient = HttpClient.newBuilder()
        .connectTimeout(connectTimeout)
        .build()

    try {
        val body = Dispatcher.dispatch(httpClient, cfg)
        println(body)
    } catch (ex: Exception) {
        // `ex.message` ist bei manchen JDK-Netzwerk-Ausnahmen (z. B.
        // `ConnectException` ohne Text) `null` — `ex` selbst trägt in jedem
        // Fall die Ausnahme-Klasse und, wo vorhanden, den Text.
        System.err.println("http-client: --verb=${cfg.verb} fehlgeschlagen: $ex")
        exitProcess(1)
    }
}
