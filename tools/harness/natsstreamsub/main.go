// Command natsstreamsub ist ein Wegwerf-Testclient für den dritten,
// vollinhaltstragenden NATS-Zustellweg (LH-FA-SST-008, ADR-0100): er belegt
// zuerst, dass ein Verbindungsversuch mit einem falschen Token vom
// NATS-Server bereits auf Verbindungsebene abgelehnt wird ("REJECTED",
// ADR-0100 Teilfrage 4) — unabhängig vom Subjekt-Namensraum, weil Core NATS
// den Token serverweit prüft. Danach verbindet er sich real mit dem
// gültigen Token, abonniert das tabellen-granulare Subjekt
// cdc.stream.<source_id>.<schema>.<table>, bevor die auslösende Change
// entsteht ("READY"), und meldet den Empfang eines vollständigen
// JSON-Change-Events ("RECEIVED") — dieselbe Zeilenform wie
// tools/harness/grpcclient/tools/harness/sseclient, damit der Aufrufer
// (tools/harness/run-integration-tests.sh) dieselben Auswertungsmuster
// wiederverwenden kann. Träger ist run-integration-tests.sh — der Aufrufer
// liest die stdout-Zeilen dieses Prozesses über `docker logs`, nicht über
// einen Exit-Code allein, weil "READY" vor der auslösenden Change
// beobachtbar sein muss.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// streamMessage trägt dasselbe Nachrichtenschema wie
// `natsstream.streamMessage` (`SPEC-021`/`SPEC-024`) — hier eigenständig
// geführt: dieses Kommando importiert kein internes Adapter-Paket.
type streamMessage struct {
	ChangeID      string          `json:"change_id"`
	TransactionID string          `json:"transaction_id"`
	SourceTableID string          `json:"source_table_id"`
	Sequence      int64           `json:"sequence"`
	Operation     string          `json:"operation"`
	OldImage      json.RawMessage `json:"old_image"`
	NewImage      json.RawMessage `json:"new_image"`
	SchemaVersion string          `json:"schema_version"`
	Schema        string          `json:"schema"`
	Table         string          `json:"table"`
}

// wrongToken ist bewusst nie der real konfigurierte Test-Token
// (compose.yaml, CDC_NATS_STREAM_TOKEN) — der Rundlauf braucht nur einen
// Wert, der garantiert vom Server abgelehnt wird, keinen bestimmten.
const wrongToken = "natsstreamsub-wrong-token"

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: natsstreamsub <nats-url> <subject> <token>")
		os.Exit(2)
	}
	url, subject, token := os.Args[1], os.Args[2], os.Args[3]

	if err := assertConnectRejectedWithWrongToken(url); err != nil {
		fmt.Fprintf(os.Stderr, "natsstreamsub: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("REJECTED")

	conn, err := nats.Connect(url, nats.Token(token))
	if err != nil {
		fmt.Fprintf(os.Stderr, "natsstreamsub: Verbindung (gültiger Token, %s) fehlgeschlagen: %v\n", url, err)
		os.Exit(1)
	}
	defer conn.Close()

	sub, err := conn.SubscribeSync(subject)
	if err != nil {
		fmt.Fprintf(os.Stderr, "natsstreamsub: SubscribeSync(%s) fehlgeschlagen: %v\n", subject, err)
		os.Exit(1)
	}
	defer func() { _ = sub.Unsubscribe() }()
	if err := conn.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "natsstreamsub: Flush nach Subscribe fehlgeschlagen: %v\n", err)
		os.Exit(1)
	}

	// Die Subscription ist server-seitig bestätigt (Flush oben) — erst ab
	// hier verpasst dieser Prozess kein zwischen Subscribe und Bereitschaft
	// gesendetes Event (dieselbe Fire-and-Forget-Grenze wie beim
	// Wecksignal-Testclient natssub).
	fmt.Println("READY")

	msg, err := sub.NextMsg(30 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "natsstreamsub: kein Event auf %s innerhalb der Frist: %v\n", subject, err)
		fmt.Println("TIMEOUT")
		os.Exit(1)
	}

	var change streamMessage
	if err := json.Unmarshal(msg.Data, &change); err != nil {
		fmt.Fprintf(os.Stderr, "natsstreamsub: Event nicht dekodierbar: %v (payload=%s)\n", err, msg.Data)
		os.Exit(1)
	}
	fmt.Printf("RECEIVED change_id=%s table=%s operation=%s new_image=%s\n",
		change.ChangeID, change.Table, change.Operation, change.NewImage)
}

// assertConnectRejectedWithWrongToken versucht real eine NATS-Verbindung
// mit einem falschen Token; der NATS-Server lehnt jede Verbindung ohne den
// konfigurierten Token bereits auf Verbindungsebene ab (`ADR-0100`
// Teilfrage 4), nicht erst beim Abonnieren eines Subjekts — die Ablehnung
// ist deshalb unabhängig vom Subjekt-Namensraum testbar, real gegen den
// laufenden Compose-NATS-Server geprüft (`docker run` mit einem falschen
// Token gegen `--auth <test-token>`).
func assertConnectRejectedWithWrongToken(url string) error {
	conn, err := nats.Connect(url,
		nats.Token(wrongToken),
		nats.RetryOnFailedConnect(false),
		nats.Timeout(5*time.Second),
	)
	if err == nil {
		conn.Close()
		return fmt.Errorf("Verbindungsversuch mit falschem Token wurde vom NATS-Server akzeptiert")
	}
	return nil
}
