package natsstream

import (
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// Die Tests dieser Datei laufen gegen einen echten NATS-Server und über
// `make test-notify` (`CDC_NATS_TEST_URL`, `tools/harness/run-notify-tests.sh`);
// ohne Server werden sie übersprungen. Das Skript fährt alle Tests des Pakets
// mit gesetzter URL und endet mit Exit 1, sobald ein Test sich überspringt;
// ein neuer Real-Server-Test braucht deshalb keinen bestimmten Namenspräfix.

func realConn(t *testing.T) *nats.Conn {
	t.Helper()
	url := os.Getenv("CDC_NATS_TEST_URL")
	if url == "" {
		t.Skip("CDC_NATS_TEST_URL nicht gesetzt — reale natsstream-Tests laufen über make test-notify")
	}
	conn, err := nats.Connect(url)
	if err != nil {
		t.Fatalf("nats.Connect(%q): %v", url, err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func subscribeSync(t *testing.T, conn *nats.Conn, subject string) *nats.Subscription {
	t.Helper()
	sub, err := conn.SubscribeSync(subject)
	if err != nil {
		t.Fatalf("SubscribeSync(%q): %v", subject, err)
	}
	t.Cleanup(func() { _ = sub.Unsubscribe() })
	return sub
}

// drain liest alle Nachrichten, die binnen 300 ms nach dem Flush eintreffen,
// und liefert ihre Subjekte sortiert.
func drain(sub *nats.Subscription) []string {
	var subjects []string
	for {
		msg, err := sub.NextMsg(300 * time.Millisecond)
		if err != nil {
			break
		}
		subjects = append(subjects, msg.Subject)
	}
	sort.Strings(subjects)
	return subjects
}

func changeWithTarget(t *testing.T, id string, table string, target model.RouteTarget) *model.Change {
	t.Helper()
	change := testChange(t)
	change.ID = model.ChangeID(id)
	change.Table = table
	change.RouteTarget = target
	return change
}

func equalSubjects(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// TestRealServerRouteSubjects belegt am Server: ein Abonnent auf dem
// Ziel-Subjekt empfängt nur die Changes dieses Ziels, ein Zielname mit `-`
// und `_` ist ein einzelnes gültiges Token, ein Wildcard-Abonnent empfängt
// beide Ziele, und die Abonnenten des Tabellen-Subjekts und der
// Wecksignal-Wurzel sehen nur das Ihre.
func TestRealServerRouteSubjects(t *testing.T) {
	conn := realConn(t)
	const source = "src-route"
	subEU := subscribeSync(t, conn, "cdc.route."+source+".eu-west_1")
	subUS := subscribeSync(t, conn, "cdc.route."+source+".us")
	subAll := subscribeSync(t, conn, "cdc.route."+source+".*")
	subRest := subscribeSync(t, conn, "cdc.route."+source+".>")
	subStream := subscribeSync(t, conn, "cdc.stream."+source+".>")
	subWake := subscribeSync(t, conn, "cdc.changes.>")
	if err := conn.Flush(); err != nil {
		t.Fatalf("Flush nach Subscribe: %v", err)
	}

	sub := newFakeSubscriber()
	p, err := New(conn, sub, source)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		p.Run(ctx)
		close(done)
	}()
	sub.changes <- changeWithTarget(t, "c1", "orders", "eu-west_1")
	sub.changes <- changeWithTarget(t, "c2", "orders", "us")
	sub.changes <- changeWithTarget(t, "c3", "orders", "")
	sub.changes <- changeWithTarget(t, "c4", "items", "eu-west_1")
	// Der Kanal ist gepuffert (4); die Veröffentlichung läuft in Run. Ein
	// Flush nach einer kurzen Wartezeit stellt sicher, dass der Server alles
	// verteilt hat, bevor gelesen wird.
	deadline := time.Now().Add(3 * time.Second)
	for len(sub.changes) > 0 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	if err := conn.Flush(); err != nil {
		t.Fatalf("Flush nach Publish: %v", err)
	}
	cancel()
	<-done

	for _, testcase := range []struct {
		name string
		sub  *nats.Subscription
		want []string
	}{
		{"Abonnent eu-west_1", subEU, []string{"cdc.route.src-route.eu-west_1", "cdc.route.src-route.eu-west_1"}},
		{"Abonnent us", subUS, []string{"cdc.route.src-route.us"}},
		{"Wildcard .*", subAll, []string{"cdc.route.src-route.eu-west_1", "cdc.route.src-route.eu-west_1", "cdc.route.src-route.us"}},
		{"Wildcard .>", subRest, []string{"cdc.route.src-route.eu-west_1", "cdc.route.src-route.eu-west_1", "cdc.route.src-route.us"}},
		{"Tabellen-Subjekt cdc.stream", subStream, []string{
			"cdc.stream.src-route.public.items", "cdc.stream.src-route.public.orders",
			"cdc.stream.src-route.public.orders", "cdc.stream.src-route.public.orders",
		}},
		{"Wecksignal-Wurzel cdc.changes", subWake, nil},
	} {
		got := drain(testcase.sub)
		if !equalSubjects(got, testcase.want) {
			t.Fatalf("%s: Subjekte %q, Erwartung %q", testcase.name, got, testcase.want)
		}
		t.Logf("natsstream-real: %s empfängt %d Nachricht(en) %q", testcase.name, len(got), got)
	}
}

// TestRealServerRoutePayloadIsByteEqual belegt am Server, dass beide
// Nachrichten einer Change denselben Payload tragen.
func TestRealServerRoutePayloadIsByteEqual(t *testing.T) {
	conn := realConn(t)
	subStream := subscribeSync(t, conn, "cdc.stream.src-bytes.>")
	subRoute := subscribeSync(t, conn, "cdc.route.src-bytes.>")
	if err := conn.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}
	p, err := New(conn, newFakeSubscriber(), "src-bytes")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	p.publish(context.Background(), changeWithTarget(t, "c1", "orders", "eu"))
	if err := conn.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}
	a, err := subStream.NextMsg(3 * time.Second)
	if err != nil {
		t.Fatalf("cdc.stream: %v", err)
	}
	b, err := subRoute.NextMsg(3 * time.Second)
	if err != nil {
		t.Fatalf("cdc.route: %v", err)
	}
	if string(a.Data) != string(b.Data) || len(a.Data) == 0 {
		t.Fatalf("Payloads unterschiedlich: %s / %s", a.Data, b.Data)
	}
	t.Logf("natsstream-real: Payload auf cdc.stream und cdc.route byte-gleich (%d Bytes)", len(a.Data))
}

// TestPublishCostWithAndWithoutTarget misst die Zeit je 10 000 Changes mit
// und ohne Ziel (Median von fünf Läufen) an `publish` samt Flush zum Server,
// ohne Abonnenten. Die Messung trägt keine Schwelle.
func TestPublishCostWithAndWithoutTarget(t *testing.T) {
	conn := realConn(t)
	p, err := New(conn, newFakeSubscriber(), "src-cost")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	const changes = 10000
	const runs = 5
	measure := func(target model.RouteTarget) time.Duration {
		change := changeWithTarget(t, "cost", "orders", target)
		durations := make([]time.Duration, 0, runs)
		for run := 0; run < runs; run++ {
			start := time.Now()
			for i := 0; i < changes; i++ {
				p.publish(context.Background(), change)
			}
			if err := conn.Flush(); err != nil {
				t.Fatalf("Flush: %v", err)
			}
			durations = append(durations, time.Since(start))
		}
		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
		return durations[runs/2]
	}
	without := measure("")
	with := measure("eu")
	t.Logf("natsstream-kosten: je %d Changes, Median von %d Läufen: ohne Ziel %v (%.0f Changes/s), mit Ziel %v (%.0f Changes/s)",
		changes, runs, without, changes/without.Seconds(), with, changes/with.Seconds())
	t.Logf("natsstream-kosten (abgeleitet): mit/ohne = %.2f", with.Seconds()/without.Seconds())
}
