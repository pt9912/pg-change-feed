package otlpexport_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	colmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricspb "go.opentelemetry.io/proto/otlp/metrics/v1"
	"google.golang.org/protobuf/proto"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/otlpexport"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
	"github.com/pt9912/pg-change-feed/internal/domain/model"
)

// received ist eine am Empfänger angekommene Anfrage.
type received struct {
	method      string
	path        string
	contentType string
	header      http.Header
	body        []byte
}

// recorder ist ein Wegwerf-Empfänger: er antwortet mit `status` und hält die
// Anfragen fest.
type recorder struct {
	mu       sync.Mutex
	requests []received
	status   int
	server   *httptest.Server
}

func newRecorder(t *testing.T, status int) *recorder {
	t.Helper()
	r := &recorder{status: status}
	r.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.requests = append(r.requests, received{
			method: req.Method, path: req.URL.Path, contentType: req.Header.Get("Content-Type"),
			header: req.Header.Clone(), body: body,
		})
		status := r.status
		r.mu.Unlock()
		w.WriteHeader(status)
	}))
	t.Cleanup(r.server.Close)
	return r
}

func (r *recorder) all() []received {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]received(nil), r.requests...)
}

func intSample(name, label string, v int64, at int64) outbound.MetricSample {
	return outbound.MetricSample{
		Name: name, Label: label, At: model.NewTimePoint(at),
		Value: outbound.MetricValue{Float: float64(v), Int: v, IsInt: true},
	}
}

func realSample(name string, v float64, at int64) outbound.MetricSample {
	return outbound.MetricSample{Name: name, At: model.NewTimePoint(at), Value: outbound.MetricValue{Float: v}}
}

const sampleAt = int64(1_700_000_000_000_000_000)

// fullBatch trägt alle zehn Kennzahlen; die Consumer- und Klassen-Zeilen
// tragen ihr Label. Der Wert der Position liegt jenseits von 2^53.
func fullBatch() outbound.MetricBatch {
	return outbound.MetricBatch{
		Source: "src-1",
		Samples: []outbound.MetricSample{
			intSample("cdc_transactions_total", "", 7, sampleAt),
			intSample("cdc_changes_processed", "", 21, sampleAt),
			realSample("cdc_oldest_change_age_seconds", 12.5, sampleAt),
			realSample("cdc_capture_lag", 0.25, sampleAt),
			intSample("cdc_consumer_position", "consumer-a", 9007199254740993, sampleAt),
			intSample("cdc_consumer_lag", "consumer-a", 3, sampleAt),
			intSample("cdc_changes_pending", "consumer-a", 4, sampleAt),
			intSample("cdc_errors_total", "schema", 1, sampleAt),
			intSample("cdc_storage_bytes", "", 8192, sampleAt),
			intSample("cdc_wal_retention_bytes", "", 1048576, sampleAt+5),
		},
	}
}

func decode(t *testing.T, body []byte) *colmetrics.ExportMetricsServiceRequest {
	t.Helper()
	var req colmetrics.ExportMetricsServiceRequest
	if err := proto.Unmarshal(body, &req); err != nil {
		t.Fatalf("Körper ist keine ExportMetricsServiceRequest: %v", err)
	}
	return &req
}

func attrs(kvs []*commonpb.KeyValue) map[string]string {
	out := map[string]string{}
	for _, kv := range kvs {
		out[kv.Key] = kv.Value.GetStringValue()
	}
	return out
}

func metricsOf(t *testing.T, req *colmetrics.ExportMetricsServiceRequest) []*metricspb.Metric {
	t.Helper()
	if len(req.ResourceMetrics) != 1 || len(req.ResourceMetrics[0].ScopeMetrics) != 1 {
		t.Fatalf("Nachricht trägt %d ResourceMetrics, wollen genau eine mit einem Scope", len(req.ResourceMetrics))
	}
	return req.ResourceMetrics[0].ScopeMetrics[0].Metrics
}

func newExporter(t *testing.T, endpoint string, headers ...otlpexport.Header) *otlpexport.Exporter {
	t.Helper()
	exporter, err := otlpexport.New(endpoint, headers)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return exporter
}

// TestExportCarriesAllTenMetricsAsGauges belegt die Drahtform (`SPEC-033`):
// `POST /v1/metrics` mit `application/x-protobuf`, im Körper alle zehn Namen
// als `Gauge` — auch die mit dem Namensteil `_total` —, Einheit und
// Datenpunkt-Attribut je Zeile der Tabelle, die Resource-Attribute
// `service.name` und `cdc.source_id` und der Messzeitpunkt je Datenpunkt. Die
// erwartete Tabelle steht hier als Literal und nicht aus dem Code des Adapters.
// Rot färbende Mutationen: den `Gauge` in `buildRequest` durch `Sum` ersetzen —
// der Typ-Fall wird rot; die Einheit `By` von `cdc_storage_bytes` in `1` ändern
// — die Einheit-Prüfung nennt den Namen; das Attribut `class` in `consumer`
// ändern — die Attribut-Prüfung nennt `cdc_errors_total`.
func TestExportCarriesAllTenMetricsAsGauges(t *testing.T) {
	rec := newRecorder(t, http.StatusOK)
	exporter := newExporter(t, rec.server.URL)

	if err := exporter.Export(context.Background(), fullBatch()); err != nil {
		t.Fatalf("Export: %v", err)
	}

	got := rec.all()
	if len(got) != 1 {
		t.Fatalf("Empfänger erhielt %d Anfragen, wollen 1", len(got))
	}
	if got[0].method != http.MethodPost || got[0].path != "/v1/metrics" || got[0].contentType != "application/x-protobuf" {
		t.Fatalf("Anfrage = %s %s (%s), wollen POST /v1/metrics (application/x-protobuf)", got[0].method, got[0].path, got[0].contentType)
	}
	req := decode(t, got[0].body)

	resource := attrs(req.ResourceMetrics[0].Resource.Attributes)
	if resource["service.name"] != "pg-change-feed" || resource["cdc.source_id"] != "src-1" || len(resource) != 2 {
		t.Fatalf("Resource-Attribute = %v, wollen service.name=pg-change-feed und cdc.source_id=src-1", resource)
	}

	want := []struct {
		name, unit, attr string
	}{
		{"cdc_transactions_total", "1", ""},
		{"cdc_changes_processed", "1", ""},
		{"cdc_oldest_change_age_seconds", "s", ""},
		{"cdc_capture_lag", "s", ""},
		{"cdc_consumer_position", "1", "consumer"},
		{"cdc_consumer_lag", "1", "consumer"},
		{"cdc_changes_pending", "1", "consumer"},
		{"cdc_errors_total", "1", "class"},
		{"cdc_storage_bytes", "By", ""},
		{"cdc_wal_retention_bytes", "By", ""},
	}
	wantLabel := map[string]string{"consumer": "consumer-a", "class": "schema"}
	metrics := metricsOf(t, req)
	if len(metrics) != len(want) {
		t.Fatalf("Nachricht trägt %d Metriken, wollen %d", len(metrics), len(want))
	}
	for i, w := range want {
		m := metrics[i]
		if m.Name != w.name {
			t.Errorf("Metrik %d heißt %q, wollen %q", i, m.Name, w.name)
			continue
		}
		gauge, ok := m.Data.(*metricspb.Metric_Gauge)
		if !ok {
			t.Errorf("%s: Typ %T, wollen Gauge", w.name, m.Data)
			continue
		}
		if m.Unit != w.unit {
			t.Errorf("%s: Einheit %q, wollen %q", w.name, m.Unit, w.unit)
		}
		if len(gauge.Gauge.DataPoints) != 1 {
			t.Errorf("%s: %d Datenpunkte, wollen 1", w.name, len(gauge.Gauge.DataPoints))
			continue
		}
		point := gauge.Gauge.DataPoints[0]
		pointAttrs := attrs(point.Attributes)
		if w.attr == "" {
			if len(pointAttrs) != 0 {
				t.Errorf("%s: Attribute %v, wollen keine", w.name, pointAttrs)
			}
		} else if len(pointAttrs) != 1 || pointAttrs[w.attr] != wantLabel[w.attr] {
			t.Errorf("%s: Attribute %v, wollen genau %s=%s", w.name, pointAttrs, w.attr, wantLabel[w.attr])
		}
		wantAt := uint64(sampleAt)
		if w.name == "cdc_wal_retention_bytes" {
			wantAt += 5
		}
		if point.TimeUnixNano != wantAt {
			t.Errorf("%s: Zeitpunkt %d, wollen %d", w.name, point.TimeUnixNano, wantAt)
		}
	}
}

// TestExportKeepsValuesExactAndByUnit belegt die Zahlenform: eine ganze Zahl
// jenseits von 2^53 geht als ganze Zahl (`as_int`) hinaus, ohne Rundung; die
// Kennzahlen mit Einheit `s` gehen immer als Fließkommazahl hinaus, auch bei
// einem ganzzahligen Wert, damit eine Reihe ihren Typ nicht wechselt; ein
// Bruch einer Kennzahl mit Einheit `1` geht als Fließkommazahl.
// Rot färbende Mutation: in `dataPoint` die Bedingung `!spec.float` streichen —
// der ganzzahlige Wert von `cdc_capture_lag` geht als ganze Zahl hinaus; die
// Zweige tauschen — die Position verliert die Genauigkeit.
func TestExportKeepsValuesExactAndByUnit(t *testing.T) {
	rec := newRecorder(t, http.StatusOK)
	exporter := newExporter(t, rec.server.URL)
	batch := outbound.MetricBatch{Source: "src-1", Samples: []outbound.MetricSample{
		intSample("cdc_consumer_position", "c", 9007199254740993, sampleAt),
		intSample("cdc_capture_lag", "", 0, sampleAt),
		realSample("cdc_oldest_change_age_seconds", 1.5, sampleAt),
		realSample("cdc_changes_processed", 2.5, sampleAt),
	}}
	if err := exporter.Export(context.Background(), batch); err != nil {
		t.Fatalf("Export: %v", err)
	}
	byName := map[string]*metricspb.NumberDataPoint{}
	for _, m := range metricsOf(t, decode(t, rec.all()[0].body)) {
		gauge, ok := m.Data.(*metricspb.Metric_Gauge)
		if !ok {
			t.Fatalf("%s: Typ %T, wollen Gauge", m.Name, m.Data)
		}
		byName[m.Name] = gauge.Gauge.DataPoints[0]
	}
	if v, ok := byName["cdc_consumer_position"].Value.(*metricspb.NumberDataPoint_AsInt); !ok || v.AsInt != 9007199254740993 {
		t.Errorf("cdc_consumer_position = %v, wollen as_int 9007199254740993", byName["cdc_consumer_position"].Value)
	}
	if _, ok := byName["cdc_capture_lag"].Value.(*metricspb.NumberDataPoint_AsDouble); !ok {
		t.Errorf("cdc_capture_lag (Einheit s, Wert 0) = %T, wollen as_double", byName["cdc_capture_lag"].Value)
	}
	if v, ok := byName["cdc_oldest_change_age_seconds"].Value.(*metricspb.NumberDataPoint_AsDouble); !ok || v.AsDouble != 1.5 {
		t.Errorf("cdc_oldest_change_age_seconds = %v, wollen as_double 1.5", byName["cdc_oldest_change_age_seconds"].Value)
	}
	if v, ok := byName["cdc_changes_processed"].Value.(*metricspb.NumberDataPoint_AsDouble); !ok || v.AsDouble != 2.5 {
		t.Errorf("cdc_changes_processed (Bruch) = %v, wollen as_double 2.5", byName["cdc_changes_processed"].Value)
	}
}

// TestExportWithoutWALMeasurementCarriesNineNames belegt `SPEC-033` (Quelle der
// Werte): ohne `cdc_wal_retention_bytes` im Zyklus entfällt genau dieser
// Datenpunkt, die neun Namen der Sicht gehen hinaus. Ein Name außerhalb der
// Tabelle geht nicht hinaus. Gegenprobe: `TestExportCarriesAllTenMetricsAsGauges`.
// Rot färbende Mutation: in `buildRequest` für einen fehlenden Datenpunkt eine
// leere Metrik anlegen (die Bedingung `len(points) == 0` streichen) — die
// Nachricht trägt zehn Namen.
func TestExportWithoutWALMeasurementCarriesNineNames(t *testing.T) {
	rec := newRecorder(t, http.StatusOK)
	exporter := newExporter(t, rec.server.URL)
	batch := fullBatch()
	batch.Samples = batch.Samples[:9]
	batch.Samples = append(batch.Samples, intSample("cdc_unbekannt", "", 1, sampleAt))
	if err := exporter.Export(context.Background(), batch); err != nil {
		t.Fatalf("Export: %v", err)
	}
	metrics := metricsOf(t, decode(t, rec.all()[0].body))
	if len(metrics) != 9 {
		t.Fatalf("Nachricht trägt %d Metriken, wollen 9", len(metrics))
	}
	for _, m := range metrics {
		if m.Name == "cdc_wal_retention_bytes" || m.Name == "cdc_unbekannt" {
			t.Errorf("Metrik %q darf nicht hinausgehen", m.Name)
		}
	}
}

// TestExportSendsConfiguredHeadersButKeepsContentType belegt, dass die Header
// der Konfiguration jede Anfrage tragen (ein Schlüssel mehrfach: beide Werte)
// und dass der Inhaltstyp der Drahtform nicht übersteuerbar ist.
// Rot färbende Mutation: in `Export` den Inhaltstyp vor die Schleife der
// Header setzen — der übersteuernde Header gewinnt und der Test nennt den Typ.
func TestExportSendsConfiguredHeadersButKeepsContentType(t *testing.T) {
	rec := newRecorder(t, http.StatusOK)
	exporter := newExporter(t, rec.server.URL,
		otlpexport.Header{Key: "Authorization", Value: "Bearer abc"},
		otlpexport.Header{Key: "X-Mehrfach", Value: "eins"},
		otlpexport.Header{Key: "X-Mehrfach", Value: "zwei"},
		otlpexport.Header{Key: "Content-Type", Value: "text/plain"},
	)
	if err := exporter.Export(context.Background(), fullBatch()); err != nil {
		t.Fatalf("Export: %v", err)
	}
	got := rec.all()[0]
	if got.header.Get("Authorization") != "Bearer abc" {
		t.Errorf("Authorization = %q", got.header.Get("Authorization"))
	}
	if values := got.header.Values("X-Mehrfach"); len(values) != 2 || values[0] != "eins" || values[1] != "zwei" {
		t.Errorf("X-Mehrfach = %v, wollen eins und zwei", values)
	}
	if got.contentType != "application/x-protobuf" {
		t.Errorf("Content-Type = %q, wollen application/x-protobuf", got.contentType)
	}
}

// TestExportBuildsExactlyOneMetricsPath belegt die Basis-URL mit und ohne
// abschließenden Schrägstrich und mit einem Pfadanteil: der Pfad
// `/v1/metrics` steht genau einmal am Ende.
// Rot färbende Mutation: in `New` `JoinPath` durch eine Verkettung mit
// `"/v1/metrics"` ersetzen — die Basis mit Schrägstrich ergibt `//v1/metrics`.
func TestExportBuildsExactlyOneMetricsPath(t *testing.T) {
	cases := []struct{ suffix, wantPath string }{
		{"", "/v1/metrics"},
		{"/", "/v1/metrics"},
		{"/otlp", "/otlp/v1/metrics"},
		{"/otlp/", "/otlp/v1/metrics"},
	}
	for _, tc := range cases {
		t.Run("Basis"+tc.suffix, func(t *testing.T) {
			rec := newRecorder(t, http.StatusOK)
			exporter := newExporter(t, rec.server.URL+tc.suffix)
			if err := exporter.Export(context.Background(), fullBatch()); err != nil {
				t.Fatalf("Export: %v", err)
			}
			if got := rec.all()[0].path; got != tc.wantPath {
				t.Fatalf("Pfad = %q, wollen %q", got, tc.wantPath)
			}
		})
	}
}

// TestExportStatusClasses belegt die Annahme: `200`, `202` und `204` sind
// angenommen; `301` (nicht verfolgt), `404` und `500` sind Fehlschläge der
// Klasse `ErrMetricExport` mit dem Status im Detail.
// Rot färbende Mutationen: in `Export` die Grenze `> 299` auf `> 399` setzen —
// der Fall `301` gilt als angenommen; die Bedingung auf `!= 200` setzen — `202`
// und `204` werden Fehlschläge.
func TestExportStatusClasses(t *testing.T) {
	cases := []struct {
		status int
		ok     bool
	}{
		{200, true}, {202, true}, {204, true},
		{301, false}, {404, false}, {500, false}, {503, false},
	}
	for _, tc := range cases {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			rec := newRecorder(t, tc.status)
			exporter := newExporter(t, rec.server.URL)
			err := exporter.Export(context.Background(), fullBatch())
			if tc.ok {
				if err != nil {
					t.Fatalf("Status %d: %v", tc.status, err)
				}
				return
			}
			if !errors.Is(err, outbound.ErrMetricExport) {
				t.Fatalf("Status %d: Fehler %v, wollen ErrMetricExport", tc.status, err)
			}
			var detail outbound.FailureDetail
			if !errors.As(err, &detail) || !strings.Contains(detail.FailureDetail(), "Status") {
				t.Fatalf("Status %d: Detail fehlt oder nennt keinen Status: %v", tc.status, err)
			}
		})
	}
}

// TestExportRefusedConnectionFails belegt die Verbindung, die abgewiesen
// wird: der Fehlschlag trägt `ErrMetricExport` und die Art, nicht die URL.
func TestExportRefusedConnectionFails(t *testing.T) {
	rec := newRecorder(t, http.StatusOK)
	url := rec.server.URL
	rec.server.Close()
	exporter := newExporter(t, url)
	err := exporter.Export(context.Background(), fullBatch())
	if !errors.Is(err, outbound.ErrMetricExport) {
		t.Fatalf("Fehler = %v, wollen ErrMetricExport", err)
	}
	if strings.Contains(err.Error(), strings.TrimPrefix(url, "http://")) {
		t.Fatalf("Fehlertext nennt die Adresse des Empfängers: %v", err)
	}
}

// TestExportExceedsDeadlineFails belegt die Frist: ein Empfänger, der nicht
// antwortet, hält den Aufruf nur bis zur Frist des Kontexts; der Fehlschlag
// nennt die Frist. Der Test hängt nicht, wenn die Frist nicht greift: er
// bricht nach zwei Sekunden mit einem Befund ab.
// Rot färbende Mutation: in `Export` `NewRequestWithContext` durch
// `NewRequest` ersetzen — der Aufruf kehrt erst mit der Antwort zurück, der Test
// meldet „Frist greift nicht“.
func TestExportExceedsDeadlineFails(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		<-release
	}))
	t.Cleanup(func() { close(release); server.Close() })
	exporter := newExporter(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- exporter.Export(ctx, fullBatch()) }()
	select {
	case err := <-done:
		if !errors.Is(err, outbound.ErrMetricExport) {
			t.Fatalf("Fehler = %v, wollen ErrMetricExport", err)
		}
		var detail outbound.FailureDetail
		if !errors.As(err, &detail) || detail.FailureDetail() != "Frist überschritten" {
			t.Fatalf("Detail = %v, wollen „Frist überschritten“", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Frist greift nicht: der Aufruf kehrt nach zwei Sekunden nicht zurück")
	}
}

// TestExportErrorsCarryNoCredentials belegt, dass kein Fehlertext einen
// Header-Wert oder den Benutzerteil der URL trägt — nicht beim Status, nicht bei
// der abgewiesenen Verbindung.
// Rot färbende Mutation: in `Export` den Fehler des Transports mit `%v` in den
// Text aufnehmen — der Text der abgewiesenen Verbindung nennt Benutzer und
// Adresse.
func TestExportErrorsCarryNoCredentials(t *testing.T) {
	rec := newRecorder(t, http.StatusInternalServerError)
	withUser := strings.Replace(rec.server.URL, "http://", "http://nutzer:geheimespasswort@", 1)
	exporter := newExporter(t, withUser, otlpexport.Header{Key: "Authorization", Value: "geheim-header"})
	statusErr := exporter.Export(context.Background(), fullBatch())

	rec.server.Close()
	connErr := exporter.Export(context.Background(), fullBatch())

	for name, err := range map[string]error{"Status": statusErr, "Verbindung": connErr} {
		if err == nil {
			t.Fatalf("%s: kein Fehler", name)
		}
		for _, secret := range []string{"geheim-header", "geheimespasswort", "nutzer"} {
			if strings.Contains(err.Error(), secret) {
				t.Errorf("%s: Fehlertext nennt %q: %v", name, secret, err)
			}
		}
	}
}

// TestExportDoesNotFollowRedirects belegt, dass eine Weiterleitung ein
// Fehlschlag ist und kein Header einen anderen Host erreicht.
// Rot färbende Mutation: `CheckRedirect` in `New` streichen — der zweite Host
// erhält die Anfrage samt Authorization-Header.
func TestExportDoesNotFollowRedirects(t *testing.T) {
	other := newRecorder(t, http.StatusOK)
	redirecting := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.server.URL+"/v1/metrics", http.StatusTemporaryRedirect)
	}))
	t.Cleanup(redirecting.Close)
	exporter := newExporter(t, redirecting.URL, otlpexport.Header{Key: "Authorization", Value: "geheim-header"})

	err := exporter.Export(context.Background(), fullBatch())
	if !errors.Is(err, outbound.ErrMetricExport) {
		t.Fatalf("Fehler = %v, wollen ErrMetricExport", err)
	}
	if n := len(other.all()); n != 0 {
		t.Fatalf("der Ziel-Host der Weiterleitung erhielt %d Anfragen, wollen 0", n)
	}
}

// TestNewRejectsUnparsableEndpoint belegt den Konstruktionsfehler.
func TestNewRejectsUnparsableEndpoint(t *testing.T) {
	if _, err := otlpexport.New("http://[::1", nil); !errors.Is(err, outbound.ErrMetricExport) {
		t.Fatalf("Fehler = %v, wollen ErrMetricExport", err)
	}
}
