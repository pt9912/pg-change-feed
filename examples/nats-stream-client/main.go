// Command nats-stream-client ist ein öffentliches Beispiel für den dritten,
// vollinhaltstragenden NATS-Zustellweg (LH-FA-SST-008, ADR-0100, SPEC-024):
// es verbindet mit einem gültigen CDC_NATS_STREAM_TOKEN, abonniert den
// Vollinhalts-Namensraum cdc.stream.> real gegen den laufenden Feed-Container
// und gibt jede empfangene Change aus. Startform ist
// `go run ./examples/nats-stream-client`; der Zugriffs-Abschnitt des
// Benutzerhandbuchs ist „Zugriff über den NATS-Vollinhalts-Stream"
// (docs/user/benutzerhandbuch.md).
//
// Dieses Programm ist das Vorbild (minimal, lesbar); der E2E-Belegträger zu
// LH-FA-SST-008 ist der Wegwerf-Client tools/harness/natsstreamsub. Das
// Beispiel trägt keine Zustandsmaschine: der Stream kennt kein Replay
// (ADR-0100), verpasste Changes holt der bestehende Lesezugriffsweg nach.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// subject ist der Wurzel-Wildcard des Vollinhalts-Namensraums (ADR-0100
// Teilfrage 2) — dieses Beispiel zeigt den Zugriffsweg über alle Quellen und
// Tabellen; ein Consumer, der nur eine Tabelle verfolgt, engt das Subjekt
// entsprechend ein (siehe Benutzerhandbuch).
const subject = "cdc.stream.>"

// waitTimeout begrenzt das Warten auf das nächste Ereignis; ein großzügiger,
// aber endlicher Wert hält das Beispiel terminierend, ohne den Regelfall
// (mehrere Changes innerhalb der Demo-Umgebung) frühzeitig zu beenden —
// dasselbe Muster wie examples/nats-client.
const waitTimeout = 24 * time.Hour

// config trägt die Laufzeit-Eingabe des Beispiels: die NATS-Server-URL und
// der Verbindungs-Token des dritten Zustellwegs. Beide kommen aus denselben
// Umgebungsvariablen, die das Benutzerhandbuch führt (CDC_NATS_URL,
// CDC_NATS_STREAM_TOKEN), und lassen sich per Flag übersteuern (ADR-0076
// Festlegung 1).
type config struct {
	natsURL string
	token   string
}

func main() {
	cfg := parseFlags()
	if cfg.natsURL == "" {
		fmt.Fprintln(os.Stderr, "nats-stream-client: keine NATS-URL gesetzt — CDC_NATS_URL (oder -nats-url) ist nötig, um den Vollinhalts-Stream zu öffnen")
		os.Exit(2)
	}
	if cfg.token == "" {
		fmt.Fprintln(os.Stderr, "nats-stream-client: kein Token gesetzt — CDC_NATS_STREAM_TOKEN (oder -token) ist nötig, um den dritten Zustellweg zu abonnieren")
		os.Exit(2)
	}

	conn, err := nats.Connect(cfg.natsURL, nats.Token(cfg.token))
	if err != nil {
		fmt.Fprintf(os.Stderr, "nats-stream-client: Verbindung zu NATS (%s) fehlgeschlagen: %v\n", cfg.natsURL, err)
		os.Exit(1)
	}
	defer conn.Close()

	sub, err := conn.SubscribeSync(subject)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nats-stream-client: Abonnement auf %s fehlgeschlagen: %v\n", subject, err)
		os.Exit(1)
	}
	defer func() { _ = sub.Unsubscribe() }()
	if err := conn.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "nats-stream-client: Flush nach Abonnement fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("nats-stream-client: lauscht auf %s\n", subject)

	for {
		msg, err := sub.NextMsg(waitTimeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "nats-stream-client: Empfang beendet: %v\n", err)
			os.Exit(1)
		}
		var change streamMessage
		if err := json.Unmarshal(msg.Data, &change); err != nil {
			fmt.Fprintf(os.Stderr, "nats-stream-client: Event nicht dekodierbar: %v (payload=%s)\n", err, msg.Data)
			os.Exit(1)
		}
		fmt.Println(formatChange(change))
	}
}

// parseFlags liest die Flag-Werte und füllt fehlende Felder aus den im
// Handbuch dokumentierten Umgebungsvariablen (ADR-0076 Festlegung 1).
func parseFlags() config {
	var cfg config
	flag.StringVar(&cfg.natsURL, "nats-url", os.Getenv("CDC_NATS_URL"), "NATS-Server-URL des Vollinhalts-Streams (Default: CDC_NATS_URL)")
	flag.StringVar(&cfg.token, "token", os.Getenv("CDC_NATS_STREAM_TOKEN"), "Verbindungs-Token des dritten Zustellwegs (Default: CDC_NATS_STREAM_TOKEN)")
	flag.Parse()
	return cfg
}
