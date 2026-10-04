// Package otlpexport trägt den OTLP/HTTP-Driven-Adapter des Metrik-Exports
// (`ADR-0149`): er bildet die Kennzahlen eines Zyklus auf die
// Nachricht `ExportMetricsServiceRequest` ab und überträgt sie als
// `POST <Basis-URL>/v1/metrics` mit dem Inhaltstyp `application/x-protobuf`.
// Die Bibliothek `go.opentelemetry.io/proto/otlp` liefert nur die Typen, der
// Transport ist `net/http`.
package otlpexport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"

	colmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricspb "go.opentelemetry.io/proto/otlp/metrics/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	"google.golang.org/protobuf/proto"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

const (
	// contentType ist der Inhaltstyp der Protobuf-Variante von OTLP/HTTP.
	contentType = "application/x-protobuf"
	// metricsPath ist der Pfad des Metrik-Signals unter der Basis-URL.
	metricsPath = "v1/metrics"
	// serviceName ist der Wert des Resource-Attributs `service.name`.
	serviceName = "pg-change-feed"
	// drainLimit begrenzt, wie viel vom Antwortkörper der Adapter liest, um die
	// Verbindung wiederzuverwenden.
	drainLimit = 1 << 20
)

// Header ist ein zusätzlicher Header jeder Anfrage (`CDC_OTLP_HEADERS`). Der
// Wert kann Zugangsdaten tragen und steht in keinem Fehlertext.
type Header struct {
	Key   string
	Value string
}

// metricSpec ist die Zeile der Kennzahlen-Tabelle: Einheit und der Schlüssel
// des Datenpunkt-Attributs, das das Label der Sicht trägt (leer: kein
// Attribut). `Real` markiert die Kennzahlen mit Einheit `s`, die immer als
// Fließkommazahl übertragen werden; alle anderen tragen eine ganze Zahl, wenn
// der Wert der Sicht eine ist.
type metricSpec struct {
	name  string
	unit  string
	attr  string
	float bool
}

// specs ist die Reihenfolge und der Inhalt der zehn Kennzahlen. Eine Zeile der
// Sicht mit einem Namen außerhalb dieser Tabelle wird nicht übertragen.
var specs = []metricSpec{
	{name: "cdc_transactions_total", unit: "1"},
	{name: "cdc_changes_processed", unit: "1"},
	{name: "cdc_oldest_change_age_seconds", unit: "s", float: true},
	{name: "cdc_capture_lag", unit: "s", float: true},
	{name: "cdc_consumer_position", unit: "1", attr: "consumer"},
	{name: "cdc_consumer_lag", unit: "By", attr: "consumer"},
	{name: "cdc_changes_pending", unit: "1", attr: "consumer"},
	{name: "cdc_errors_total", unit: "1", attr: "class"},
	{name: "cdc_storage_bytes", unit: "By"},
	{name: "cdc_wal_retention_bytes", unit: "By"},
}

// Exporter implementiert `outbound.MetricExportPort` über OTLP/HTTP.
type Exporter struct {
	url     string
	headers []Header
	client  *http.Client
}

var _ outbound.MetricExportPort = (*Exporter)(nil)

// New baut den Adapter für die Basis-URL des Empfängers (Schema `http` oder
// `https`, die Konfiguration prüft das vorab); der Pfad `/v1/metrics` kommt
// einmal dazu, mit oder ohne abschließenden Schrägstrich der Basis. Eine
// Weiterleitung wird nicht verfolgt: ein `3xx` ist ein Fehlschlag, und kein
// Header erreicht einen anderen Host. Die Frist je Versuch trägt der Kontext
// des Aufrufs.
func New(endpoint string, headers []Header) (*Exporter, error) {
	base, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: Basis-URL nicht lesbar", outbound.ErrMetricExport)
	}
	return &Exporter{
		url:     base.JoinPath(metricsPath).String(),
		headers: append([]Header(nil), headers...),
		client: &http.Client{
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}, nil
}

// exportError ist der Fehler der Übertragung: sein Text trägt Art und Status,
// nie die URL, einen Header-Wert oder den rohen Fehlertext des Transports.
type exportError struct{ detail string }

func (e *exportError) Error() string         { return outbound.ErrMetricExport.Error() + ": " + e.detail }
func (e *exportError) Unwrap() error         { return outbound.ErrMetricExport }
func (e *exportError) FailureDetail() string { return e.detail }

// Export überträgt einen Zyklus. `2xx` ist angenommen; jeder andere Status,
// jeder Verbindungs- und jeder Zeitfehler ist ein Fehlschlag.
func (e *Exporter) Export(ctx context.Context, batch outbound.MetricBatch) error {
	body, err := proto.Marshal(buildRequest(batch))
	if err != nil {
		return &exportError{detail: "Kodierung fehlgeschlagen"}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.url, bytes.NewReader(body))
	if err != nil {
		return &exportError{detail: "Anfrage nicht bildbar"}
	}
	for _, h := range e.headers {
		req.Header.Add(h.Key, h.Value)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := e.client.Do(req)
	if err != nil {
		return &exportError{detail: transportDetail(ctx, err)}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, drainLimit))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &exportError{detail: "Status " + strconv.Itoa(resp.StatusCode)}
	}
	return nil
}

// transportDetail benennt die Art eines Transportfehlers ohne dessen Text: der
// Text eines `*url.Error` nennt die URL samt Benutzerteil.
func transportDetail(ctx context.Context, err error) string {
	var netErr net.Error
	switch {
	case errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) ||
		(errors.As(err, &netErr) && netErr.Timeout()):
		return "Frist überschritten"
	case errors.Is(err, context.Canceled):
		return "abgebrochen"
	default:
		return "Verbindung fehlgeschlagen"
	}
}

// buildRequest bildet einen Zyklus auf die Nachricht ab: alle Kennzahlen als
// `Gauge` (auch die mit dem Namensteil `_total`; ein Wert der Sicht ist ein
// Momentstand), je Name eine Metrik mit Einheit der Tabelle, je Zeile ein
// Datenpunkt mit dem Messzeitpunkt der Zeile und, wo die Tabelle es nennt, dem
// Label als Attribut.
func buildRequest(batch outbound.MetricBatch) *colmetrics.ExportMetricsServiceRequest {
	metrics := make([]*metricspb.Metric, 0, len(specs))
	for _, spec := range specs {
		var points []*metricspb.NumberDataPoint
		for _, sample := range batch.Samples {
			if sample.Name != spec.name {
				continue
			}
			points = append(points, dataPoint(spec, sample))
		}
		if len(points) == 0 {
			continue
		}
		metrics = append(metrics, &metricspb.Metric{
			Name: spec.name,
			Unit: spec.unit,
			Data: &metricspb.Metric_Gauge{Gauge: &metricspb.Gauge{DataPoints: points}},
		})
	}
	return &colmetrics.ExportMetricsServiceRequest{
		ResourceMetrics: []*metricspb.ResourceMetrics{{
			Resource: &resourcepb.Resource{Attributes: []*commonpb.KeyValue{
				stringAttr("service.name", serviceName),
				stringAttr("cdc.source_id", string(batch.Source)),
			}},
			ScopeMetrics: []*metricspb.ScopeMetrics{{
				Scope:   &commonpb.InstrumentationScope{Name: serviceName},
				Metrics: metrics,
			}},
		}},
	}
}

func dataPoint(spec metricSpec, sample outbound.MetricSample) *metricspb.NumberDataPoint {
	point := &metricspb.NumberDataPoint{TimeUnixNano: uint64(sample.At.UnixNanos)}
	if spec.attr != "" {
		point.Attributes = []*commonpb.KeyValue{stringAttr(spec.attr, sample.Label)}
	}
	if sample.Value.IsInt && !spec.float {
		point.Value = &metricspb.NumberDataPoint_AsInt{AsInt: sample.Value.Int}
	} else {
		point.Value = &metricspb.NumberDataPoint_AsDouble{AsDouble: sample.Value.Float}
	}
	return point
}

func stringAttr(key, value string) *commonpb.KeyValue {
	return &commonpb.KeyValue{Key: key, Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: value}}}
}
