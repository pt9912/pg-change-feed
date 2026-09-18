package bootstrap

import "testing"

// TestChangeStreamEnabled trägt die Oder-Bedingung der
// Broadcaster-Verdrahtung (`ADR-0061` Teilfrage 5, um die dritte
// Oder-Bedingung erweitert durch `ADR-0100` Teilfrage 5): Die
// Live-Streaming-Fähigkeit ist genau dann aktiv, wenn `CDC_GRPC_ADDR` oder
// `CDC_HTTP_ADDR` gesetzt ist, oder der dritte NATS-Vollinhalts-Zustellweg
// aktiv ist (`natsStreamActive`). Sind alle drei Bedingungen falsch, bleibt
// der `Broadcaster` unkonstruiert und der `CaptureService` trägt keinen
// Stream-Publish-Schritt — dieselbe Aussage, die `Run` beim Aufbau der
// Capture-Optionen zieht.
func TestChangeStreamEnabled(t *testing.T) {
	for _, testcase := range []struct {
		name             string
		grpcAddr         string
		httpAddr         string
		natsStreamActive bool
		want             bool
	}{
		{name: "alle drei aus", grpcAddr: "", httpAddr: "", natsStreamActive: false, want: false},
		{name: "nur gRPC gesetzt", grpcAddr: ":9090", httpAddr: "", natsStreamActive: false, want: true},
		{name: "nur HTTP gesetzt", grpcAddr: "", httpAddr: ":8090", natsStreamActive: false, want: true},
		{name: "nur NATS-Stream aktiv", grpcAddr: "", httpAddr: "", natsStreamActive: true, want: true},
		{name: "alle drei gesetzt", grpcAddr: ":9090", httpAddr: ":8090", natsStreamActive: true, want: true},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			if got := changeStreamEnabled(testcase.grpcAddr, testcase.httpAddr, testcase.natsStreamActive); got != testcase.want {
				t.Fatalf("changeStreamEnabled(%q, %q, %v) = %v, Erwartung %v",
					testcase.grpcAddr, testcase.httpAddr, testcase.natsStreamActive, got, testcase.want)
			}
		})
	}
}

// TestNatsStreamEnabled trägt `ADR-0100` Teilfrage 5s Zwei-Bedingungen-Gate
// (DoD-Regressionsklasse „Zwei-Bedingungen-Gate"): der dritte
// NATS-Vollinhalts-Zustellweg ist genau dann aktiv, wenn **beide**
// `CDC_NATS_URL` und `CDC_NATS_STREAM_TOKEN` gesetzt sind. Ist nur eine der
// beiden gesetzt — insbesondere nur `CDC_NATS_URL`, der heutige Zustand
// jeder bestehenden Installation —, bleibt der Weg deaktiviert.
//
// Rot färbende Mutation: `natsStreamEnabled` auf `natsURL != "" ||
// natsStreamToken != ""` ändern (Oder statt Und) — der Fall „nur
// CDC_NATS_URL gesetzt" färbt dann rot, weil ein bestehendes Deployment
// versehentlich den dritten Weg aktivierte.
func TestNatsStreamEnabled(t *testing.T) {
	for _, testcase := range []struct {
		name            string
		natsURL         string
		natsStreamToken string
		want            bool
	}{
		{name: "beide leer", natsURL: "", natsStreamToken: "", want: false},
		{name: "nur CDC_NATS_URL gesetzt (Bestandsverhalten)", natsURL: "nats://nats:4222", natsStreamToken: "", want: false},
		{name: "nur CDC_NATS_STREAM_TOKEN gesetzt", natsURL: "", natsStreamToken: "s3cr3t", want: false},
		{name: "beide gesetzt", natsURL: "nats://nats:4222", natsStreamToken: "s3cr3t", want: true},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			if got := natsStreamEnabled(testcase.natsURL, testcase.natsStreamToken); got != testcase.want {
				t.Fatalf("natsStreamEnabled(%q, %q) = %v, Erwartung %v",
					testcase.natsURL, testcase.natsStreamToken, got, testcase.want)
			}
		})
	}
}

// TestValidateNatsStreamTokenRequiresURL trägt `ADR-0100` Teilfrage 5s
// Konfigurationsfehler-Zweig: ein gesetzter Token ohne URL ist ein
// Start-Hindernis (`ErrConfiguration`), jede andere Kombination — inklusive
// beider leer und beider gesetzt — ist zulässig.
func TestValidateNatsStreamTokenRequiresURL(t *testing.T) {
	for _, testcase := range []struct {
		name            string
		natsURL         string
		natsStreamToken string
		wantErr         bool
	}{
		{name: "beide leer", natsURL: "", natsStreamToken: "", wantErr: false},
		{name: "nur URL gesetzt", natsURL: "nats://nats:4222", natsStreamToken: "", wantErr: false},
		{name: "beide gesetzt", natsURL: "nats://nats:4222", natsStreamToken: "s3cr3t", wantErr: false},
		{name: "nur Token gesetzt", natsURL: "", natsStreamToken: "s3cr3t", wantErr: true},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			err := validateNatsStreamTokenRequiresURL(testcase.natsURL, testcase.natsStreamToken)
			if testcase.wantErr && err == nil {
				t.Fatal("Fehler erwartet, aber nil erhalten")
			}
			if !testcase.wantErr && err != nil {
				t.Fatalf("kein Fehler erwartet, erhalten: %v", err)
			}
		})
	}
}
