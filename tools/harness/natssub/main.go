// Command natssub ist ein Wegwerf-Testclient für den NATS-Happy-Path-Beleg
// (LH-FA-SST-007, ADR-0055): er abonniert ein Subjekt real, bevor die
// auslösende Change entsteht, meldet die Bereitschaft über die Zeile
// "READY" auf stdout und wartet danach auf genau ein Wecksignal. Träger ist
// tools/harness/run-integration-tests.sh — der Aufrufer liest die stdout-
// Zeilen dieses Prozesses über `docker logs`, nicht über einen Exit-Code
// allein, weil "READY" vor der auslösenden Change beobachtbar sein muss.
//
// Der optionale vierte Aufrufparameter <token> trägt den seit `ADR-0100`
// serverweiten NATS-Verbindungs-Token (Teilfrage 4/5): sobald
// `compose.yaml`s `nats`-Service mit `--auth` läuft, verlangt der Server ihn
// von **jeder** Verbindung — auch von diesem bislang anonymen
// Wecksignal-Testclient. Ein leerer/fehlender Wert verbindet weiterhin ohne
// Token (Bestandsverhalten gegen einen NATS-Server ohne `--auth`).
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: natssub <nats-url> <subject> [token]")
		os.Exit(2)
	}
	url, subject := os.Args[1], os.Args[2]
	var opts []nats.Option
	if len(os.Args) == 4 && os.Args[3] != "" {
		opts = append(opts, nats.Token(os.Args[3]))
	}

	conn, err := nats.Connect(url, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "natssub: Verbindung (%s) fehlgeschlagen: %v\n", url, err)
		os.Exit(1)
	}
	defer conn.Close()

	sub, err := conn.SubscribeSync(subject)
	if err != nil {
		fmt.Fprintf(os.Stderr, "natssub: SubscribeSync(%s) fehlgeschlagen: %v\n", subject, err)
		os.Exit(1)
	}
	defer func() { _ = sub.Unsubscribe() }()
	if err := conn.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "natssub: Flush nach Subscribe fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}

	// Die Subscription ist server-seitig bestätigt (Flush oben) — erst ab
	// hier verpasst dieser Prozess kein zwischen Subscribe und Bereitschaft
	// gesendetes Signal.
	fmt.Println("READY")

	msg, err := sub.NextMsg(30 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "natssub: kein Signal auf %s innerhalb der Frist: %v\n", subject, err)
		fmt.Println("TIMEOUT")
		os.Exit(1)
	}
	fmt.Printf("RECEIVED subject=%s payload_len=%d\n", msg.Subject, len(msg.Data))
}
