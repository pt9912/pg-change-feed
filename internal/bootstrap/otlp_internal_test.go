package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	colmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	metricspb "go.opentelemetry.io/proto/otlp/metrics/v1"
	"google.golang.org/protobuf/proto"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/otlpexport"
	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/application/usecase/exportmetrics"
	"github.com/pt9912/pg-change-feed/internal/domain/messagecode"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// otlpZugriffswege trägt die zwei Zugriffswege der Konfiguration: den reinen
// Umgebungsweg und den verdrahteten Weg mit geladener Datei, die
// `dateiInhalt` trägt. Beide rufen `applyOTLP`.
func otlpZugriffswege(t *testing.T, dateiInhalt string) map[string]func(map[string]string) (Config, error) {
	t.Helper()
	datei := writeConfigFile(t, "source_id: src-datei\npublication: pub-datei\nslot: slot-datei\ntables:\n  public.t1:\n    table_id: tbl-1\n    schema_version: sv-1\n"+dateiInhalt)
	return map[string]func(map[string]string) (Config, error){
		"ConfigFromEnv": func(env map[string]string) (Config, error) {
			return ConfigFromEnv(func(name string) string { return env[name] })
		},
		"ConfigFromEnvAndFile": func(env map[string]string) (Config, error) {
			env["CDC_CONFIG_FILE"] = datei
			return ConfigFromEnvAndFile(func(name string) string { return env[name] })
		},
	}
}

func otlpEnv(extra map[string]string) map[string]string {
	env := vollständigeEnvOhneDatei()
	for k, v := range extra {
		env[k] = v
	}
	return env
}

// TestOTLPKonfigurationWirdGelesen trägt `LH-FA-SST-010` Happy Path und
// Boundary der Konfiguration auf beiden Zugriffswegen: ohne Angabe bleibt der
// Export aus (Endpunkt leer, Takt 60 s, keine Header); Endpunkt und Header aus
// der Umgebung werden gelesen — der Header am ersten `=` geteilt, Leerraum im
// Wert erhalten —; die Grenzwerte 5 und 3600 sind gültig; unter geladener Datei
// trägt das Feld `otlp_interval` den Takt, und die Umgebungsvariable schlägt es;
// eine leere Umgebungsvariable gilt als nicht gesetzt; ein Takt ohne Endpunkt
// ist zulässig.
// Rot färbende Mutationen: in `mergeConfig` den Aufruf von `applyOTLP` mit
// leerem `fileInterval` ausführen — die Fälle des Dateiwegs tragen den Takt
// der Datei nicht; in `parseOTLPHeaders` `strings.Cut` durch ein Teilen am
// letzten `=` ersetzen — der Wert `a=b=` kommt verkürzt an.
func TestOTLPKonfigurationWirdGelesen(t *testing.T) {
	cases := []struct {
		name         string
		datei        string
		env          map[string]string
		nurDatei     bool
		wantEndpoint string
		wantHeaders  []otlpexport.Header
		wantInterval time.Duration
	}{
		{name: "ungesetzt", wantInterval: 60 * time.Second},
		{name: "leere Variablen gelten als ungesetzt", wantInterval: 60 * time.Second,
			env: map[string]string{"CDC_OTLP_ENDPOINT": "", "CDC_OTLP_HEADERS": "", "CDC_OTLP_INTERVAL_SECONDS": ""}},
		{name: "Endpunkt http", env: map[string]string{"CDC_OTLP_ENDPOINT": "http://collector:4318"},
			wantEndpoint: "http://collector:4318", wantInterval: 60 * time.Second},
		{name: "Endpunkt https mit Pfad", env: map[string]string{"CDC_OTLP_ENDPOINT": "https://otel.example.org/otlp/"},
			wantEndpoint: "https://otel.example.org/otlp/", wantInterval: 60 * time.Second},
		{name: "Header am ersten Gleichheitszeichen geteilt",
			env:          map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1", "CDC_OTLP_HEADERS": "Authorization=Bearer abc,x-key=a=b=,leer="},
			wantEndpoint: "http://c:1", wantInterval: 60 * time.Second,
			wantHeaders: []otlpexport.Header{{Key: "Authorization", Value: "Bearer abc"}, {Key: "x-key", Value: "a=b="}, {Key: "leer", Value: ""}}},
		{name: "Takt untere Grenze", env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": "5"}, wantInterval: 5 * time.Second},
		{name: "Takt obere Grenze", env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": "3600"}, wantInterval: 3600 * time.Second},
		{name: "Takt ohne Endpunkt zulässig", env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": "120"}, wantInterval: 120 * time.Second},
		{name: "Takt aus der Datei", datei: "otlp_interval: 30\n", nurDatei: true, wantInterval: 30 * time.Second},
		{name: "Takt aus der Datei als Zeichenkette", datei: "otlp_interval: \"45\"\n", nurDatei: true, wantInterval: 45 * time.Second},
		{name: "Umgebung schlägt Datei", datei: "otlp_interval: 30\n", nurDatei: true,
			env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": "90"}, wantInterval: 90 * time.Second},
		{name: "leere Umgebungsvariable lässt den Datei-Wert stehen", datei: "otlp_interval: 30\n", nurDatei: true,
			env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": ""}, wantInterval: 30 * time.Second},
		{name: "Endpunkt aus der Umgebung unter geladener Datei", datei: "otlp_interval: 30\n", nurDatei: true,
			env: map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1"}, wantEndpoint: "http://c:1", wantInterval: 30 * time.Second},
	}
	for _, tc := range cases {
		for weg, laden := range otlpZugriffswege(t, tc.datei) {
			if tc.nurDatei && weg == "ConfigFromEnv" {
				continue
			}
			t.Run(weg+"/"+tc.name, func(t *testing.T) {
				cfg, err := laden(otlpEnv(tc.env))
				if err != nil {
					t.Fatalf("Konfiguration lädt nicht: %v", err)
				}
				if cfg.OTLPEndpoint != tc.wantEndpoint {
					t.Errorf("Endpunkt = %q, erwartet %q", cfg.OTLPEndpoint, tc.wantEndpoint)
				}
				if cfg.OTLPInterval != tc.wantInterval {
					t.Errorf("Takt = %v, erwartet %v", cfg.OTLPInterval, tc.wantInterval)
				}
				if len(cfg.OTLPHeaders) != len(tc.wantHeaders) {
					t.Fatalf("Header = %v, erwartet %v", cfg.OTLPHeaders, tc.wantHeaders)
				}
				for i, h := range tc.wantHeaders {
					if cfg.OTLPHeaders[i] != h {
						t.Errorf("Header %d = %+v, erwartet %+v", i, cfg.OTLPHeaders[i], h)
					}
				}
			})
		}
	}
}

// TestOTLPKonfigurationUngueltigEndetMitConfiguration trägt `LH-FA-SST-010`
// Negative auf beiden Zugriffswegen: jeder ungültige Wert endet mit der
// Fehlerklasse `configuration`, dem Meldungscode seiner Art und dem
// Klartext, der Variable und — wo es einen gibt — die Position nennt. Die
// Gegenproben (Grenzwerte 5 und 3600, Endpunkt mit beiden Schemata, gültige
// Header) laden in `TestOTLPKonfigurationWirdGelesen`.
// Rot färbende Mutationen: in `validateOTLPEndpoint` die Schema-Prüfung
// streichen — `ftp://…` lädt; in `parseOTLPInterval` die untere Grenze 5 auf 4
// setzen — der Fall „Takt 4“ lädt; die Bedingung `headersRaw != ""` in
// `applyOTLP` streichen — der Fall „Header ohne Endpunkt“ lädt.
func TestOTLPKonfigurationUngueltigEndetMitConfiguration(t *testing.T) {
	cases := []struct {
		name     string
		datei    string
		env      map[string]string
		nurDatei bool
		sentinel error
		code     messagecode.Code
		want     []string
	}{
		{name: "Endpunkt Schema ftp", env: map[string]string{"CDC_OTLP_ENDPOINT": "ftp://c:21"},
			sentinel: ErrOTLPEndpointInvalid, code: messagecode.OTLPEndpointInvalid,
			want: []string{"PCF-E2011", "CDC_OTLP_ENDPOINT", "http oder https"}},
		{name: "Endpunkt ohne Schema", env: map[string]string{"CDC_OTLP_ENDPOINT": "collector:4318"},
			sentinel: ErrOTLPEndpointInvalid, code: messagecode.OTLPEndpointInvalid,
			want: []string{"PCF-E2011", "http oder https"}},
		{name: "Endpunkt ohne Host", env: map[string]string{"CDC_OTLP_ENDPOINT": "http://"},
			sentinel: ErrOTLPEndpointInvalid, code: messagecode.OTLPEndpointInvalid,
			want: []string{"PCF-E2011", "Host fehlt"}},
		{name: "Endpunkt mit Port ohne Host", env: map[string]string{"CDC_OTLP_ENDPOINT": "https://:4318/x"},
			sentinel: ErrOTLPEndpointInvalid, code: messagecode.OTLPEndpointInvalid,
			want: []string{"PCF-E2011", "Host fehlt"}},
		{name: "Endpunkt nicht lesbar", env: map[string]string{"CDC_OTLP_ENDPOINT": "http://[::1"},
			sentinel: ErrOTLPEndpointInvalid, code: messagecode.OTLPEndpointInvalid,
			want: []string{"PCF-E2011", "nicht als URL lesbar"}},
		{name: "Takt 4", env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": "4"},
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013", "CDC_OTLP_INTERVAL_SECONDS", "\"4\"", "5 bis 3600"}},
		{name: "Takt 3601", env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": "3601"},
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013", "\"3601\""}},
		{name: "Takt abc", env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": "abc"},
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013", "\"abc\""}},
		{name: "Takt 60.5", env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": "60.5"},
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013", "\"60.5\""}},
		{name: "Takt negativ", env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": "-5"},
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013", "\"-5\""}},
		{name: "Takt mit Leerraum", env: map[string]string{"CDC_OTLP_INTERVAL_SECONDS": " 60"},
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013"}},
		{name: "Takt aus der Datei 4", datei: "otlp_interval: 4\n", nurDatei: true,
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013", "otlp_interval", "\"4\""}},
		{name: "Takt aus der Datei 0", datei: "otlp_interval: 0\n", nurDatei: true,
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013", "otlp_interval"}},
		{name: "Takt aus der Datei 60.5", datei: "otlp_interval: 60.5\n", nurDatei: true,
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013", "otlp_interval", "\"60.5\""}},
		{name: "Takt aus der Datei abc", datei: "otlp_interval: abc\n", nurDatei: true,
			sentinel: ErrOTLPIntervalInvalid, code: messagecode.OTLPIntervalInvalid,
			want: []string{"PCF-E2013", "otlp_interval", "\"abc\""}},
		{name: "Header ohne Gleichheitszeichen",
			env:      map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1", "CDC_OTLP_HEADERS": "a=b,ohnewert"},
			sentinel: ErrOTLPHeadersInvalid, code: messagecode.OTLPHeadersInvalid,
			want: []string{"PCF-E2012", "CDC_OTLP_HEADERS", "Element 2", "kein ="}},
		{name: "Header mit leerem Schlüssel",
			env:      map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1", "CDC_OTLP_HEADERS": "=wert"},
			sentinel: ErrOTLPHeadersInvalid, code: messagecode.OTLPHeadersInvalid,
			want: []string{"PCF-E2012", "Element 1", "leeren Schlüssel"}},
		{name: "Header mit leerem Element",
			env:      map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1", "CDC_OTLP_HEADERS": "a=b,,c=d"},
			sentinel: ErrOTLPHeadersInvalid, code: messagecode.OTLPHeadersInvalid,
			want: []string{"PCF-E2012", "Element 2"}},
		{name: "Header mit Leerraum im Schlüssel",
			env:      map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1", "CDC_OTLP_HEADERS": "a b=wert"},
			sentinel: ErrOTLPHeadersInvalid, code: messagecode.OTLPHeadersInvalid,
			want: []string{"PCF-E2012", "Element 1", "Leerraum im Schlüssel"}},
		{name: "Header mit unzulässigem Zeichen im Schlüssel",
			env:      map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1", "CDC_OTLP_HEADERS": "a:b=wert"},
			sentinel: ErrOTLPHeadersInvalid, code: messagecode.OTLPHeadersInvalid,
			want: []string{"PCF-E2012", "unzulässiges Zeichen"}},
		{name: "Header mit Steuerzeichen im Wert",
			env:      map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1", "CDC_OTLP_HEADERS": "a=w\nert"},
			sentinel: ErrOTLPHeadersInvalid, code: messagecode.OTLPHeadersInvalid,
			want: []string{"PCF-E2012", "Steuerzeichen im Wert"}},
		{name: "Header ohne Endpunkt",
			env:      map[string]string{"CDC_OTLP_HEADERS": "a=b"},
			sentinel: ErrOTLPHeadersWithoutEndpoint, code: messagecode.OTLPHeadersNoEndpoint,
			want: []string{"PCF-E2014", "CDC_OTLP_HEADERS gesetzt", "CDC_OTLP_ENDPOINT fehlt"}},
	}
	for _, tc := range cases {
		for weg, laden := range otlpZugriffswege(t, tc.datei) {
			if tc.nurDatei && weg == "ConfigFromEnv" {
				continue
			}
			t.Run(weg+"/"+tc.name, func(t *testing.T) {
				_, err := laden(otlpEnv(tc.env))
				if err == nil {
					t.Fatal("ungültige Export-Konfiguration lädt ohne Fehler")
				}
				if !errors.Is(err, ErrConfiguration) || !errors.Is(err, tc.sentinel) {
					t.Fatalf("Fehler ist nicht ErrConfiguration und %v: %v", tc.sentinel, err)
				}
				if code, ok := messagecode.From(err); !ok || code != tc.code {
					t.Fatalf("Code = %q (%v), erwartet %q", code, ok, tc.code)
				}
				for _, want := range append([]string{"configuration"}, tc.want...) {
					if !strings.Contains(err.Error(), want) {
						t.Fatalf("Fehlertext %q trägt %q nicht", err.Error(), want)
					}
				}
			})
		}
	}
}

// TestOTLPKonfigurationsfehlerNennenKeineZugangsdaten trägt, dass der Text
// eines Konfigurationsfehlers weder einen Header-Wert noch den Benutzerteil
// einer URL nennt: der Endpunkt `ftp://nutzer:pw@…` und ein Header mit
// Steuerzeichen im Wert `geheim` enden im Fehler, ohne dass `nutzer`, `pw` oder
// `geheim` im Text stehen.
// Rot färbende Mutation: in `validateOTLPEndpoint` den Endpunkt in den Text
// des Falls „Schema“ aufnehmen — der Test nennt `nutzer`.
func TestOTLPKonfigurationsfehlerNennenKeineZugangsdaten(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  map[string]string
	}{
		{"Endpunkt mit Benutzerteil", map[string]string{"CDC_OTLP_ENDPOINT": "ftp://nutzer:pw@c:21"}},
		{"Header mit Steuerzeichen im Wert", map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1", "CDC_OTLP_HEADERS": "a=geheim\n"}},
		{"Header ohne Gleichheitszeichen", map[string]string{"CDC_OTLP_ENDPOINT": "http://c:1", "CDC_OTLP_HEADERS": "geheim"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ConfigFromEnv(func(name string) string { return otlpEnv(tc.env)[name] })
			if err == nil {
				t.Fatal("ungültige Konfiguration lädt")
			}
			for _, secret := range []string{"nutzer", "pw@", "geheim"} {
				if strings.Contains(err.Error(), secret) {
					t.Fatalf("Fehlertext nennt %q: %v", secret, err)
				}
			}
		})
	}
}

// TestOTLPIntervalGehoertNichtZurZugangsdatenKlasse trägt die Gegenprobe zur
// Zugangsdaten-Klasse: eine Datei mit `otlp_interval: 30` lädt ohne Fehler,
// eine mit `otlp_endpoint:` oder `otlp_headers:` endet mit `ErrConfiguration`
// und der Zeile „Zugangsdaten bleiben env-var-exklusiv“.
// Rot färbende Mutation: `otlp_interval` in `forbiddenFileCredentialKeys`
// aufnehmen — die Datei mit dem Takt endet mit dieser Zeile.
func TestOTLPIntervalGehoertNichtZurZugangsdatenKlasse(t *testing.T) {
	file, err := ConfigFromFile(writeConfigFile(t, "otlp_interval: 30\n"))
	if err != nil {
		t.Fatalf("Datei mit otlp_interval lädt nicht: %v", err)
	}
	if file.OTLPInterval == nil || *file.OTLPInterval != "30" {
		t.Fatalf("Datei-Feld = %v, erwartet \"30\"", file.OTLPInterval)
	}
	for _, key := range []string{"otlp_endpoint", "otlp_headers"} {
		_, err := ConfigFromFile(writeConfigFile(t, key+": \n"))
		if !errors.Is(err, ErrConfiguration) || err == nil || !strings.Contains(err.Error(), "Zugangsdaten bleiben env-var-exklusiv") {
			t.Fatalf("%s in der Datei: Fehler = %v, erwartet die Zeile der Zugangsdaten-Klasse", key, err)
		}
	}
}

// fakeUseCase zählt die Aufrufe des Takts.
type fakeUseCase struct {
	calls atomic.Int64
	err   error
}

func (u *fakeUseCase) Export(context.Context) error {
	u.calls.Add(1)
	return u.err
}

func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("Zeitüberschreitung: %s", what)
}

// TestRunMetricExportTicksCallsOncePerTick trägt den Takt: jeder Takt ruft den
// Use Case genau einmal — ein Fehlschlag wird nicht innerhalb des Taktes
// wiederholt und ein ausgefallener Takt nicht nachgeholt —, und der Kontext
// beendet die Schleife.
// Rot färbende Mutation: in `runMetricExportTicks` den Aufruf bei einem Fehler
// ein zweites Mal ausführen — die Zählung nennt doppelt so viele Aufrufe wie
// Takte.
func TestRunMetricExportTicksCallsOncePerTick(t *testing.T) {
	useCase := &fakeUseCase{err: errors.New("Fehlschlag")}
	ticks := make(chan time.Time)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { runMetricExportTicks(ctx, useCase, ticks); close(done) }()

	for i := 1; i <= 3; i++ {
		ticks <- time.Now()
		waitFor(t, fmt.Sprintf("Aufruf %d", i), func() bool { return useCase.calls.Load() >= int64(i) })
	}
	time.Sleep(50 * time.Millisecond)
	if got := useCase.calls.Load(); got != 3 {
		t.Fatalf("Aufrufe = %d, erwartet 3 (einer je Takt)", got)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("die Schleife endet nicht mit dem Kontext")
	}
}

// TestRecordingWALMeasurerHoldsOnlySuccessfulMeasurements trägt den Halter:
// eine erfolgreiche Messung wird unverändert durchgereicht und mit dem
// Zeitpunkt der Uhr festgehalten, eine fehlgeschlagene gibt den Fehler
// unverändert zurück und lässt den Halter unberührt; vor der ersten Messung ist
// `ok` falsch.
// Rot färbende Mutation: in `recordingWALMeasurer.Measure` die Bedingung
// `err == nil` streichen — der Fall der fehlgeschlagenen Messung hält 0 Byte
// fest.
func TestRecordingWALMeasurerHoldsOnlySuccessfulMeasurements(t *testing.T) {
	gauge := &walRetentionGauge{}
	if _, _, ok := gauge.LastWALRetention(); ok {
		t.Fatal("der Halter trägt vor der ersten Messung einen Wert")
	}
	inner := &stubMeasurer{value: 4096}
	clock := fixedClock{at: model.NewTimePoint(1234)}
	measurer := recordingWALMeasurer{inner: inner, gauge: gauge, clock: clock}

	got, err := measurer.Measure(context.Background())
	if err != nil || got != 4096 {
		t.Fatalf("Measure = %d, %v, erwartet 4096, nil", got, err)
	}
	if bytes, at, ok := gauge.LastWALRetention(); !ok || bytes != 4096 || at.UnixNanos != 1234 {
		t.Fatalf("Halter = %d, %v, %v, erwartet 4096 zum Zeitpunkt 1234", bytes, at, ok)
	}

	inner.value, inner.err = 99, errors.New("Messung fehlgeschlagen")
	clock.at = model.NewTimePoint(5678)
	measurer.clock = clock
	if _, err := measurer.Measure(context.Background()); err == nil {
		t.Fatal("der Fehler der Messung geht verloren")
	}
	if bytes, at, _ := gauge.LastWALRetention(); bytes != 4096 || at.UnixNanos != 1234 {
		t.Fatalf("Halter nach Fehlschlag = %d zum Zeitpunkt %d, erwartet unverändert 4096 / 1234", bytes, at.UnixNanos)
	}
}

type stubMeasurer struct {
	value int64
	err   error
}

func (m *stubMeasurer) Measure(context.Context) (int64, error) { return m.value, m.err }

// captureLog hält jede Log-Zeile samt Attributen fest.
type captureLog struct {
	mu      sync.Mutex
	entries []string
	warns   int
}

func (l *captureLog) add(level, msg string, attrs []any) {
	line := level + " " + msg + fmt.Sprint(attrs...)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, line)
	if level == "WARN" {
		l.warns++
	}
}

func (l *captureLog) Debug(_ context.Context, msg string, attrs ...any) { l.add("DEBUG", msg, attrs) }
func (l *captureLog) Info(_ context.Context, msg string, attrs ...any)  { l.add("INFO", msg, attrs) }
func (l *captureLog) Warn(_ context.Context, msg string, attrs ...any)  { l.add("WARN", msg, attrs) }
func (l *captureLog) Error(_ context.Context, msg string, attrs ...any) { l.add("ERROR", msg, attrs) }

func (l *captureLog) text() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Join(l.entries, "\n")
}

func (l *captureLog) warnCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.warns
}

type stubReader struct {
	samples []outbound.MetricSample
	err     error
}

func (r stubReader) Read(context.Context) ([]outbound.MetricSample, error) {
	if r.err != nil {
		return nil, r.err
	}
	return append([]outbound.MetricSample(nil), r.samples...), nil
}

func oneSample() []outbound.MetricSample {
	return []outbound.MetricSample{{Name: "cdc_changes_processed", Value: outbound.MetricValue{Float: 7, Int: 7, IsInt: true}}}
}

// TestMetricExportLogsCarryNoCredentials trägt `LH-FA-SST-010` (Fehlschlag
// sichtbar) ohne Zugangsdaten im Log: ein Empfänger mit Status 500 und ein
// abgewiesener Empfänger, jeweils mit einem Header `Authorization=geheim` und
// einer URL mit Benutzer-Passwort-Teil, führen zu einer Warn-Zeile, und in
// keiner Log-Zeile des Laufs (alle Stufen, alle Attribute) steht der
// Header-Wert, der Benutzer oder das Passwort.
// Rot färbende Mutation: in `exportmetrics.failed` `"error", err` an die Attribute
// hängen und im Adapter den Fehler des Transports mit `%v` in den Text
// aufnehmen — die Zeile der abgewiesenen Verbindung nennt Benutzer und Adresse.
func TestMetricExportLogsCarryNoCredentials(t *testing.T) {
	status500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(status500.Close)
	refused := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	refusedURL := refused.URL
	refused.Close()

	for name, endpoint := range map[string]string{"Status 500": status500.URL, "abgewiesen": refusedURL} {
		t.Run(name, func(t *testing.T) {
			withUser := strings.Replace(endpoint, "http://", "http://nutzer:geheimespasswort@", 1)
			log := &captureLog{}
			cfg := Config{Source: "src-1", OTLPEndpoint: withUser, OTLPInterval: 20 * time.Millisecond,
				OTLPHeaders: []otlpexport.Header{{Key: "Authorization", Value: "geheim-header"}}}
			stop, err := startMetricExportLoop(context.Background(), cfg, stubReader{samples: oneSample()}, log, systemclock.New(), &walRetentionGauge{})
			if err != nil {
				t.Fatalf("startMetricExportLoop: %v", err)
			}
			waitFor(t, "erste Warnung", func() bool { return log.warnCount() >= 1 })
			stop()
			text := log.text()
			if !strings.Contains(text, string(messagecode.WarnExportFailed)) {
				t.Fatalf("Log trägt PCF-W6001 nicht:\n%s", text)
			}
			for _, secret := range []string{"geheim-header", "geheimespasswort", "nutzer"} {
				if strings.Contains(text, secret) {
					t.Fatalf("Log nennt %q:\n%s", secret, text)
				}
			}
		})
	}
}

// TestMetricExportThrottlesWarningsAcrossCycles trägt die Drosselung am echten
// Takt: viele Fehlschläge in kurzer Folge erzeugen genau eine Warnung. Die
// Fake-Uhr-Fälle trägt der Test des Use Cases.
func TestMetricExportThrottlesWarningsAcrossCycles(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)
	log := &captureLog{}
	cfg := Config{Source: "src-1", OTLPEndpoint: server.URL, OTLPInterval: 15 * time.Millisecond}
	stop, err := startMetricExportLoop(context.Background(), cfg, stubReader{samples: oneSample()}, log, systemclock.New(), &walRetentionGauge{})
	if err != nil {
		t.Fatalf("startMetricExportLoop: %v", err)
	}
	waitFor(t, "fünf Anfragen", func() bool { return requests.Load() >= 5 })
	stop()
	if got := log.warnCount(); got != 1 {
		t.Fatalf("Warnungen = %d bei %d Anfragen, erwartet 1", got, requests.Load())
	}
}

type countingHeartbeat struct{ beats atomic.Int64 }

func (h *countingHeartbeat) Beat(context.Context, model.SourceID) error {
	h.beats.Add(1)
	return nil
}

func (h *countingHeartbeat) Fault(context.Context, model.SourceID, model.ErrorClass, messagecode.Code) error {
	return nil
}

// TestUnresponsiveReceiverDoesNotHoldHeartbeatAndEndsAtDeadline trägt
// `LH-FA-SST-010` Boundary: ein Empfänger, der nicht antwortet, hält den
// Heartbeat-Zug nicht auf (er schlägt weiter im eigenen Takt), jeder Versuch
// endet an der Frist (der Empfänger sieht mehrere Anfragen, jede von der
// Gegenseite abgebrochen), und der Fehlzustand erzeugt eine Warnung mit
// `PCF-W6001`; Health-Zustand und `error_class` bleiben unberührt (der
// Heartbeat trägt keinen Fehlerzustand).
// Rot färbende Mutation: in `startMetricExportLoop` `WithTimeout` weglassen —
// der erste Versuch kehrt nie zurück, der Empfänger sieht eine einzige Anfrage.
func TestUnresponsiveReceiverDoesNotHoldHeartbeatAndEndsAtDeadline(t *testing.T) {
	release := make(chan struct{})
	var requests, abandoned atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Der Körper wird gelesen, damit der Server das Ende der Verbindung
		// bemerkt und den Kontext der Anfrage beendet.
		_, _ = io.Copy(io.Discard, r.Body)
		requests.Add(1)
		select {
		case <-r.Context().Done():
			abandoned.Add(1)
		case <-release:
		}
	}))
	t.Cleanup(func() { close(release); server.Close() })

	log := &captureLog{}
	cfg := Config{Source: "src-1", OTLPEndpoint: server.URL, OTLPInterval: 100 * time.Millisecond}
	stop, err := startMetricExportLoop(context.Background(), cfg, stubReader{samples: oneSample()}, log, systemclock.New(), &walRetentionGauge{})
	if err != nil {
		t.Fatalf("startMetricExportLoop: %v", err)
	}
	heartbeat := &countingHeartbeat{}
	hbCtx, hbCancel := context.WithCancel(context.Background())
	hbDone := make(chan struct{})
	go func() { runHeartbeat(hbCtx, heartbeat, "src-1", 10*time.Millisecond); close(hbDone) }()

	waitFor(t, "drei Anfragen an der Frist", func() bool { return requests.Load() >= 3 && abandoned.Load() >= 2 })
	beats := heartbeat.beats.Load()
	if beats < 10 {
		t.Fatalf("Heartbeat schlug %d-mal während des blockierten Exports, erwartet mindestens 10", beats)
	}
	stop()
	hbCancel()
	<-hbDone
	if !strings.Contains(log.text(), string(messagecode.WarnExportFailed)) {
		t.Fatalf("keine Warnung PCF-W6001:\n%s", log.text())
	}
	if got := log.warnCount(); got != 1 {
		t.Fatalf("Warnungen = %d, erwartet 1 (gedrosselt)", got)
	}
}

// TestStartMetricExportWithoutEndpointStartsNothing trägt `LH-FA-SST-010`
// Boundary „kein Empfänger konfiguriert“: ohne Endpunkt entsteht kein Sender
// (`started` falsch), der Aufruf der Stopp-Funktion ist ein No-Op, und ein
// vorhandener Empfänger sieht über ein Beobachtungsfenster keine Anfrage.
// Rot färbende Mutation: in `startMetricExport` die Bedingung
// `cfg.OTLPEndpoint == ""` streichen — der Aufruf versucht den Lese-Pool und
// den Sender zu bauen (`started` wahr).
func TestStartMetricExportWithoutEndpointStartsNothing(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { requests.Add(1) }))
	t.Cleanup(server.Close)

	stop, started, err := startMetricExport(context.Background(), Config{Source: "src-1", ReaderDSN: "postgres://x:x@127.0.0.1:1/db", OTLPInterval: 10 * time.Millisecond},
		&captureLog{}, systemclock.New(), &walRetentionGauge{})
	if err != nil || started {
		t.Fatalf("startMetricExport = started %v, %v, erwartet false, nil", started, err)
	}
	time.Sleep(200 * time.Millisecond)
	stop()
	if got := requests.Load(); got != 0 {
		t.Fatalf("Empfänger sah %d Anfragen, erwartet 0", got)
	}
}

// TestStartMetricExportDeliversTheValuesOfTheView trägt `LH-FA-SST-010` Happy
// Path über die Verdrahtung gegen die reale Datenbank (`make test-store`): der
// Empfänger erhält einen Request, dessen `cdc_changes_processed` und
// `cdc_transactions_total` gleich den Werten der Sicht sind (Lesen unter einem
// `cdc_reader`-Login), dessen Resource die Quellkennung trägt und dessen
// `cdc_wal_retention_bytes` der Wert des Halters mit dessen Zeitpunkt ist.
// Rot färbende Mutation: in `startMetricExportLoop` den Halter `wal` durch einen
// leeren Halter ersetzen — die Nachricht trägt keinen
// `cdc_wal_retention_bytes`; in `exportmetrics.Service.Export` die Zeilen der
// Sicht verwerfen — die Zählwerte fehlen.
func TestStartMetricExportDeliversTheValuesOfTheView(t *testing.T) {
	baseDSN := os.Getenv("CDC_STORE_TEST_DSN")
	if baseDSN == "" {
		t.Skip("CDC_STORE_TEST_DSN nicht gesetzt — reale PostgreSQL-Tests laufen über make test-store")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, baseDSN)
	if err != nil {
		t.Fatalf("Verbindungsaufbau: %v", err)
	}
	t.Cleanup(pool.Close)

	const source, login, password = "src-otlp-wiring", "pgc_test_otlp_reader", "test-login-password"
	for _, stmt := range []string{
		"DELETE FROM cdc.transaction WHERE source_id = '" + source + "'",
		"INSERT INTO cdc.source (source_id, name) VALUES ('" + source + "', 'OTLP-Verdrahtung') ON CONFLICT (source_id) DO NOTHING",
		"INSERT INTO cdc.transaction (transaction_id, source_id, commit_position, committed_at) VALUES ('tx-otlp-wiring-1', '" + source + "', 11, current_timestamp)",
		"DROP ROLE IF EXISTS " + login,
		"CREATE ROLE " + login + " LOGIN PASSWORD '" + password + "' IN ROLE cdc_reader",
	} {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			t.Fatalf("Fixture (%s): %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM cdc.transaction WHERE source_id = '"+source+"'")
		_, _ = pool.Exec(context.Background(), "DROP ROLE IF EXISTS "+login)
	})
	parsed, err := url.Parse(baseDSN)
	if err != nil {
		t.Fatalf("Basis-DSN nicht parsebar: %v", err)
	}
	parsed.User = url.UserPassword(login, password)

	var mu sync.Mutex
	var bodies [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		bodies = append(bodies, body)
		mu.Unlock()
	}))
	t.Cleanup(server.Close)

	gauge := &walRetentionGauge{}
	gauge.set(5120, model.NewTimePoint(42))
	cfg := Config{Source: source, ReaderDSN: parsed.String(), OTLPEndpoint: server.URL, OTLPInterval: 50 * time.Millisecond}
	stop, started, err := startMetricExport(ctx, cfg, &captureLog{}, systemclock.New(), gauge)
	if err != nil || !started {
		t.Fatalf("startMetricExport = started %v, %v", started, err)
	}
	t.Cleanup(stop)
	waitFor(t, "ein Request", func() bool { mu.Lock(); defer mu.Unlock(); return len(bodies) > 0 })

	mu.Lock()
	body := bodies[0]
	mu.Unlock()
	var req colmetrics.ExportMetricsServiceRequest
	if err := proto.Unmarshal(body, &req); err != nil {
		t.Fatalf("Körper: %v", err)
	}
	var sourceAttr string
	for _, kv := range req.ResourceMetrics[0].Resource.Attributes {
		if kv.Key == "cdc.source_id" {
			sourceAttr = kv.Value.GetStringValue()
		}
	}
	if sourceAttr != source {
		t.Errorf("cdc.source_id = %q, erwartet %q", sourceAttr, source)
	}
	got := map[string]*metricspb.NumberDataPoint{}
	for _, m := range req.ResourceMetrics[0].ScopeMetrics[0].Metrics {
		if len(m.Data.(*metricspb.Metric_Gauge).Gauge.DataPoints) > 0 {
			got[m.Name] = m.Data.(*metricspb.Metric_Gauge).Gauge.DataPoints[0]
		}
	}
	for _, name := range []string{"cdc_changes_processed", "cdc_transactions_total"} {
		var want int64
		if err := pool.QueryRow(ctx, "SELECT value::bigint FROM cdc.metrics WHERE metric_name = $1", name).Scan(&want); err != nil {
			t.Fatalf("Sicht (%s): %v", name, err)
		}
		point, ok := got[name]
		if !ok {
			t.Errorf("%s fehlt in der Nachricht", name)
			continue
		}
		if v, ok := point.Value.(*metricspb.NumberDataPoint_AsInt); !ok || v.AsInt != want {
			t.Errorf("%s = %v, die Sicht trägt %d", name, point.Value, want)
		}
	}
	if want := int64(1); got["cdc_transactions_total"] == nil || got["cdc_transactions_total"].Value.(*metricspb.NumberDataPoint_AsInt).AsInt < want {
		t.Errorf("cdc_transactions_total trägt nicht die Fixture-Transaktion")
	}
	wal, ok := got["cdc_wal_retention_bytes"]
	if !ok || wal.Value.(*metricspb.NumberDataPoint_AsInt).AsInt != 5120 || wal.TimeUnixNano != 42 {
		t.Errorf("cdc_wal_retention_bytes = %v, erwartet 5120 zum Zeitpunkt 42", wal)
	}
}

// TestExportServiceWithFailingThenRecoveringReceiver trägt „nach
// Wiederaufnahme genau ein Request je Takt“ über den Takt: ein Empfänger, der
// dreimal ausfällt und danach annimmt, sieht über fünf Takte genau fünf
// Aufrufe (kein Nachholen, keine Wiederholung innerhalb eines Taktes), die
// Wiederaufnahme erzeugt eine Info-Zeile und der Fehlzustand eine Warnung.
func TestExportServiceWithFailingThenRecoveringReceiver(t *testing.T) {
	var calls atomic.Int64
	exporter := exporterFunc(func(context.Context, outbound.MetricBatch) error {
		if calls.Add(1) <= 3 {
			return outbound.ErrMetricExport
		}
		return nil
	})
	log := &captureLog{}
	service := exportmetrics.New("src-1", stubReader{samples: oneSample()}, &walRetentionGauge{}, exporter, systemclock.New(), exportmetrics.WithLog(log))
	ticks := make(chan time.Time)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { runMetricExportTicks(ctx, service, ticks); close(done) }()
	for i := 1; i <= 5; i++ {
		ticks <- time.Now()
		waitFor(t, fmt.Sprintf("Aufruf %d", i), func() bool { return calls.Load() >= int64(i) })
	}
	cancel()
	<-done
	if got := calls.Load(); got != 5 {
		t.Fatalf("Anfragen = %d, erwartet 5 (eine je Takt)", got)
	}
	text := log.text()
	if log.warnCount() != 1 || !strings.Contains(text, "INFO metrics-export: Übertragung wieder aufgenommen") {
		t.Fatalf("erwartet eine Warnung und die Info der Wiederaufnahme:\n%s", text)
	}
}

type exporterFunc func(context.Context, outbound.MetricBatch) error

func (f exporterFunc) Export(ctx context.Context, b outbound.MetricBatch) error { return f(ctx, b) }
