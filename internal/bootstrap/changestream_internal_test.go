package bootstrap

import (
	"errors"
	"testing"
)

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

// TestConfigFromEnvBindetTokenOhneURL trägt den Aufrufort von
// `validateNatsStreamTokenRequiresURL` in `ConfigFromEnv` (`wiring.go`):
// der Start-Fehlervertrag hängt an diesem Aufruf, nicht allein an der
// reinen Funktion (`TestValidateNatsStreamTokenRequiresURL`). Beide
// Richtungen sind vertreten — der abgelehnte Halbzustand und der
// zulässige Vollzustand —, damit der Test nicht auch bei einem
// entfernten Aufruf grün bliebe.
//
// Rot färbende Mutation: den `validateNatsStreamTokenRequiresURL`-Aufruf
// aus `ConfigFromEnv` streichen — der Fall „nur Token gesetzt" färbt rot,
// weil das Laden dann fehlerfrei durchläuft, obwohl das Handbuch den
// Container-Start für genau diese Kombination ausschließt.
func TestConfigFromEnvBindetTokenOhneURL(t *testing.T) {
	for _, testcase := range []struct {
		name    string
		zusatz  map[string]string
		wantErr bool
	}{
		{name: "nur Token gesetzt", zusatz: map[string]string{"CDC_NATS_STREAM_TOKEN": "s3cr3t"}, wantErr: true},
		{name: "Token mit URL", zusatz: map[string]string{"CDC_NATS_URL": "nats://nats:4222", "CDC_NATS_STREAM_TOKEN": "s3cr3t"}, wantErr: false},
		{name: "keine NATS-Variable gesetzt", zusatz: nil, wantErr: false},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			values := vollständigeEnvOhneDatei()
			for name, wert := range testcase.zusatz {
				values[name] = wert
			}
			_, err := ConfigFromEnv(func(name string) string { return values[name] })
			if testcase.wantErr && !errors.Is(err, ErrConfiguration) {
				t.Fatalf("Fehlerklasse configuration erwartet, erhalten: %v", err)
			}
			if !testcase.wantErr && err != nil {
				t.Fatalf("kein Fehler erwartet, erhalten: %v", err)
			}
		})
	}
}

// TestMergeConfigBindetTokenOhneURL trägt denselben Aufrufort auf dem
// Datei-Ladepfad (`mergeConfig`, `config_file.go`): auch unter gesetzter
// `CDC_CONFIG_FILE` bricht ein Token ohne URL den Start mit
// `ErrConfiguration` ab, und die zulässige Kombination lädt fehlerfrei.
//
// Rot färbende Mutation: den `validateNatsStreamTokenRequiresURL`-Aufruf
// aus `mergeConfig` streichen — der Fall „nur Token gesetzt" färbt rot.
func TestMergeConfigBindetTokenOhneURL(t *testing.T) {
	path := writeConfigFile(t, dateiOhneAdressen)

	t.Run("nur Token gesetzt", func(t *testing.T) {
		values := envMitDatei(path, map[string]string{"CDC_NATS_STREAM_TOKEN": "s3cr3t"})
		_, err := ConfigFromEnvAndFile(func(name string) string { return values[name] })
		if !errors.Is(err, ErrConfiguration) {
			t.Fatalf("Fehlerklasse configuration erwartet, erhalten: %v", err)
		}
	})

	t.Run("Token mit URL", func(t *testing.T) {
		values := envMitDatei(path, map[string]string{
			"CDC_NATS_URL":          "nats://nats:4222",
			"CDC_NATS_STREAM_TOKEN": "s3cr3t",
		})
		cfg, err := ConfigFromEnvAndFile(func(name string) string { return values[name] })
		if err != nil {
			t.Fatalf("Token mit URL unter geladener Datei: %v", err)
		}
		if cfg.NatsStreamToken != "s3cr3t" {
			t.Fatalf("NatsStreamToken: Feld trägt %q, Erwartung: %q", cfg.NatsStreamToken, "s3cr3t")
		}
	})
}

// TestValidateNatsStreamTokenRequiresURL trägt `ADR-0100` Teilfrage 5s
// Konfigurationsfehler-Zweig: ein gesetzter Token ohne URL ist ein
// Start-Hindernis (`ErrConfiguration`), jede andere Kombination — inklusive
// beider leer und beider gesetzt — ist zulässig. Geprüft ist hier die reine
// Funktion; ihre zwei Aufruforte (`ConfigFromEnv`, `mergeConfig`) binden
// `TestConfigFromEnvBindetTokenOhneURL` und
// `TestMergeConfigBindetTokenOhneURL`.
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
