package bootstrap

import "testing"

// TestChangeStreamEnabled trägt die Oder-Bedingung der
// Broadcaster-Verdrahtung (`ADR-0061` Teilfrage 5): Die Live-Streaming-Fähigkeit
// ist genau dann aktiv, wenn `CDC_GRPC_ADDR` oder `CDC_HTTP_ADDR` gesetzt ist.
// Sind beide leer, bleibt der `Broadcaster` unkonstruiert und der
// `CaptureService` trägt keinen Stream-Publish-Schritt — dieselbe Aussage, die
// `Run` beim Aufbau der Capture-Optionen zieht.
func TestChangeStreamEnabled(t *testing.T) {
	for _, testcase := range []struct {
		name     string
		grpcAddr string
		httpAddr string
		want     bool
	}{
		{name: "beide leer", grpcAddr: "", httpAddr: "", want: false},
		{name: "nur gRPC gesetzt", grpcAddr: ":9090", httpAddr: "", want: true},
		{name: "nur HTTP gesetzt", grpcAddr: "", httpAddr: ":8090", want: true},
		{name: "beide gesetzt", grpcAddr: ":9090", httpAddr: ":8090", want: true},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			if got := changeStreamEnabled(testcase.grpcAddr, testcase.httpAddr); got != testcase.want {
				t.Fatalf("changeStreamEnabled(%q, %q) = %v, Erwartung %v",
					testcase.grpcAddr, testcase.httpAddr, got, testcase.want)
			}
		})
	}
}
