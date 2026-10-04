package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
	"strings"
)

// spec ist eine Zeile der Kennzahlen-Tabelle: Name, Einheit und der Schlüssel
// des Datenpunkt-Attributs (leer: keines). Die Tabelle steht hier als
// eigenes Literal; sie ist die Erwartung an den Empfänger-Ausgang, nicht der
// Code des Senders.
type spec struct {
	name, unit, attr string
}

var specs = []spec{
	{"cdc_transactions_total", "1", ""},
	{"cdc_changes_processed", "1", ""},
	{"cdc_oldest_change_age_seconds", "s", ""},
	{"cdc_capture_lag", "s", ""},
	{"cdc_consumer_position", "1", "consumer"},
	{"cdc_consumer_lag", "By", "consumer"},
	{"cdc_changes_pending", "1", "consumer"},
	{"cdc_errors_total", "1", "class"},
	{"cdc_storage_bytes", "By", ""},
	{"cdc_wal_retention_bytes", "By", ""},
}

const walName = "cdc_wal_retention_bytes"

// Die JSON-Form der Datei des Collectors (`file`-Exporter, eine Zeile je
// Export): die OTLP/JSON-Abbildung der Nachricht.
type record struct {
	ResourceMetrics []struct {
		Resource struct {
			Attributes []keyValue `json:"attributes"`
		} `json:"resource"`
		ScopeMetrics []struct {
			Metrics []metric `json:"metrics"`
		} `json:"scopeMetrics"`
	} `json:"resourceMetrics"`
}

type keyValue struct {
	Key   string `json:"key"`
	Value struct {
		StringValue *string `json:"stringValue"`
	} `json:"value"`
}

type metric struct {
	Name  string          `json:"name"`
	Unit  string          `json:"unit"`
	Gauge *dataPoints     `json:"gauge"`
	Sum   json.RawMessage `json:"sum"`
}

type dataPoints struct {
	DataPoints []dataPoint `json:"dataPoints"`
}

type dataPoint struct {
	Attributes   []keyValue      `json:"attributes"`
	TimeUnixNano json.RawMessage `json:"timeUnixNano"`
	AsInt        json.RawMessage `json:"asInt"`
	AsDouble     json.RawMessage `json:"asDouble"`
}

// Options sind die Eingaben einer Prüfung.
type Options struct {
	// Source ist die erwartete Quellkennung (`cdc.source_id` der Resource).
	Source string
	// ToleranceSeconds ist der erlaubte Abstand der zwei zeitabhängigen
	// Kennzahlen (Einheit `s`) zur Sicht.
	ToleranceSeconds float64
	// WALMeasured sind die Messungen des WAL-Rückstands, die der Prozess
	// selbst gelogen hat (`cdc_wal_retention_bytes` steht nicht in der Sicht):
	// der exportierte Wert ist genau eine von ihnen. nil schaltet die Prüfung ab.
	WALMeasured []int64
	// AuthAttr und AuthValue: ein Resource-Attribut, das der Collector aus dem
	// Anfrage-Header bildet; ein leerer Wert verlangt seine Abwesenheit.
	AuthAttr, AuthValue string
}

// Result ist der Befund einer bestandenen Prüfung.
type Result struct {
	Names        int
	WAL          string
	ConsumerLag  string
	ErrorsLabels []string
	Points       int
	Auth         string
}

// LastRecord liest die letzte vollständige Zeile der Datei des Collectors und
// die Zahl der vollständigen, lesbaren Zeilen. Eine Zeile ohne Zeilenende ist
// noch im Schreiben und zählt nicht.
func LastRecord(path string) (last []byte, count int, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	complete := bytes.Count(data, []byte("\n"))
	for i := 0; sc.Scan() && i < complete; i++ {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		if !json.Valid(line) {
			return nil, count, fmt.Errorf("Zeile %d der Collector-Datei ist kein JSON", i+1)
		}
		last = append(last[:0], line...)
		count++
	}
	return last, count, nil
}

// ReadView liest die Zeilen der Sicht (`name|label|wert`, eine je Zeile).
func ReadView(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rows := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			return nil, fmt.Errorf("Zeile der Sicht ohne drei Felder: %q", line)
		}
		key := parts[0] + "|" + parts[1]
		if _, dup := rows[key]; dup {
			return nil, fmt.Errorf("Zeile der Sicht doppelt: %q", key)
		}
		rows[key] = parts[2]
	}
	return rows, nil
}

func contains(values []int64, want int64) bool {
	for _, v := range values {
		if v == want {
			return true
		}
	}
	return false
}

func maxOf(values []int64) int64 {
	var m int64
	for i, v := range values {
		if i == 0 || v > m {
			m = v
		}
	}
	return m
}

// ReadInts liest eine Zahl je Zeile (die gelogenen Messungen des Prozesses).
func ReadInts(path string) ([]int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	values := []int64{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		var v int64
		if _, err := fmt.Sscanf(line, "%d", &v); err != nil {
			return nil, fmt.Errorf("Zeile %q ist keine ganze Zahl", line)
		}
		values = append(values, v)
	}
	return values, nil
}

func attrMap(kvs []keyValue) map[string]string {
	out := map[string]string{}
	for _, kv := range kvs {
		if kv.Value.StringValue != nil {
			out[kv.Key] = *kv.Value.StringValue
		} else {
			out[kv.Key] = "<kein String>"
		}
	}
	return out
}

func rawNumber(raw json.RawMessage) (*big.Rat, bool) {
	text := strings.Trim(strings.TrimSpace(string(raw)), `"`)
	if text == "" || text == "null" {
		return nil, false
	}
	r, ok := new(big.Rat).SetString(text)
	return r, ok
}

// Check prüft einen Export des Collectors gegen die Sicht und die Tabelle.
// Alle Abweichungen werden gesammelt und als ein Fehler gemeldet.
func Check(recordJSON []byte, view map[string]string, opts Options) (*Result, error) {
	var rec record
	if err := json.Unmarshal(recordJSON, &rec); err != nil {
		return nil, fmt.Errorf("Export nicht lesbar: %v", err)
	}
	var problems []string
	fail := func(format string, args ...any) { problems = append(problems, fmt.Sprintf(format, args...)) }

	if len(rec.ResourceMetrics) != 1 || len(rec.ResourceMetrics[0].ScopeMetrics) != 1 {
		return nil, fmt.Errorf("Export trägt %d ResourceMetrics, wollen genau eine mit einem Scope", len(rec.ResourceMetrics))
	}
	resource := attrMap(rec.ResourceMetrics[0].Resource.Attributes)
	wantResource := map[string]string{"service.name": "pg-change-feed", "cdc.source_id": opts.Source}
	auth := "-"
	if opts.AuthAttr != "" {
		got, present := resource[opts.AuthAttr]
		switch {
		case opts.AuthValue == "" && present:
			fail("Resource-Attribut %s ist vorhanden (%q), erwartet abwesend", opts.AuthAttr, got)
		case opts.AuthValue != "" && !present:
			fail("Resource-Attribut %s fehlt, erwartet %q", opts.AuthAttr, opts.AuthValue)
		case opts.AuthValue != "" && got != opts.AuthValue:
			fail("Resource-Attribut %s = %q, erwartet %q", opts.AuthAttr, got, opts.AuthValue)
		}
		auth = "absent"
		if present {
			auth = "present"
			delete(resource, opts.AuthAttr)
		}
	}
	for k, want := range wantResource {
		if resource[k] != want {
			fail("Resource-Attribut %s = %q, erwartet %q", k, resource[k], want)
		}
	}
	if len(resource) != len(wantResource) {
		fail("Resource-Attribute %v, erwartet genau %v", resource, wantResource)
	}

	byName := map[string]metric{}
	for _, m := range rec.ResourceMetrics[0].ScopeMetrics[0].Metrics {
		if _, dup := byName[m.Name]; dup {
			fail("Metrik %s kommt doppelt vor", m.Name)
		}
		byName[m.Name] = m
	}
	known := map[string]bool{}
	for _, s := range specs {
		known[s.name] = true
	}
	for name := range byName {
		if !known[name] {
			fail("Metrik %s steht nicht in der Tabelle", name)
		}
	}

	res := &Result{Auth: auth}
	seen := map[string]bool{}
	for _, s := range specs {
		m, ok := byName[s.name]
		if !ok {
			fail("Metrik %s fehlt im Export", s.name)
			continue
		}
		res.Names++
		if m.Gauge == nil || len(m.Sum) > 0 {
			fail("%s: Typ ist kein Gauge", s.name)
			continue
		}
		if m.Unit != s.unit {
			fail("%s: Einheit %q, erwartet %q", s.name, m.Unit, s.unit)
		}
		if s.name == "cdc_consumer_lag" {
			res.ConsumerLag = m.Unit
		}
		if len(m.Gauge.DataPoints) == 0 {
			fail("%s: kein Datenpunkt", s.name)
		}
		for _, p := range m.Gauge.DataPoints {
			res.Points++
			attrs := attrMap(p.Attributes)
			label := ""
			if s.attr == "" {
				if len(attrs) != 0 {
					fail("%s: Attribute %v, erwartet keine", s.name, attrs)
				}
			} else {
				if len(attrs) != 1 {
					fail("%s: Attribute %v, erwartet genau %s", s.name, attrs, s.attr)
				}
				v, ok := attrs[s.attr]
				if !ok {
					fail("%s: Attribut %s fehlt (Attribute %v)", s.name, s.attr, attrs)
				}
				label = v
			}
			if ts, ok := rawNumber(p.TimeUnixNano); !ok || ts.Sign() <= 0 {
				fail("%s[%s]: Zeitpunkt fehlt", s.name, label)
			}
			var value *big.Rat
			isInt := len(p.AsInt) > 0
			if isInt {
				value, _ = rawNumber(p.AsInt)
			} else {
				value, _ = rawNumber(p.AsDouble)
			}
			if value == nil {
				fail("%s[%s]: kein Wert", s.name, label)
				continue
			}
			if s.unit == "s" && isInt {
				fail("%s[%s]: Einheit s geht als ganze Zahl hinaus, erwartet Fließkommazahl", s.name, label)
			}
			if s.name == walName {
				res.WAL = value.RatString()
				if _, inView := view[walName+"|"+label]; inView {
					fail("%s steht in der Sicht, erwartet nur im Export", walName)
				}
				if opts.WALMeasured != nil {
					if !value.IsInt() || !contains(opts.WALMeasured, value.Num().Int64()) {
						fail("%s = %s, keine der %d Messungen des Prozesses (größte %d)", walName, value.RatString(), len(opts.WALMeasured), maxOf(opts.WALMeasured))
					}
				}
				continue
			}
			key := s.name + "|" + label
			seen[key] = true
			text, inView := view[key]
			if !inView {
				fail("%s[%s]: keine Zeile in der Sicht", s.name, label)
				continue
			}
			viewValue, ok := new(big.Rat).SetString(strings.TrimSpace(text))
			if !ok {
				fail("%s[%s]: Wert der Sicht %q nicht lesbar", s.name, label, text)
				continue
			}
			if s.unit == "s" {
				a, _ := value.Float64()
				b, _ := viewValue.Float64()
				if math.Abs(a-b) > opts.ToleranceSeconds {
					fail("%s[%s]: Export %v, Sicht %v, Toleranz %v s", s.name, label, a, b, opts.ToleranceSeconds)
				}
				continue
			}
			if value.Cmp(viewValue) != 0 {
				fail("%s[%s]: Export %s, Sicht %s", s.name, label, value.RatString(), viewValue.RatString())
			}
			if viewValue.IsInt() && !isInt {
				fail("%s[%s]: ganzzahliger Wert geht als Fließkommazahl hinaus", s.name, label)
			}
			if s.name == "cdc_errors_total" {
				res.ErrorsLabels = append(res.ErrorsLabels, label)
			}
		}
	}
	for key := range view {
		if !seen[key] {
			fail("Zeile %s der Sicht fehlt im Export", key)
		}
	}
	sort.Strings(res.ErrorsLabels)
	if len(problems) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(problems, "\n"))
	}
	return res, nil
}
