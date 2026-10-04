package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixture baut einen Export im JSON-Format des `file`-Exporters. `edit`
// verändert die Metrik-Zeilen vor dem Zusammensetzen.
type fixture struct {
	unit      map[string]string
	attrKey   map[string]string
	gaugeName map[string]bool
	value     map[string]string
	asDouble  map[string]bool
	asInt     map[string]bool
	resource  string
}

func newFixture() *fixture {
	return &fixture{
		unit:      map[string]string{},
		attrKey:   map[string]string{},
		gaugeName: map[string]bool{},
		value:     map[string]string{},
		asDouble:  map[string]bool{},
		asInt:     map[string]bool{},
		resource:  `{"key":"service.name","value":{"stringValue":"pg-change-feed"}},{"key":"cdc.source_id","value":{"stringValue":"src-1"}}`,
	}
}

var baseValues = map[string]string{
	"cdc_transactions_total":        "7",
	"cdc_changes_processed":         "21",
	"cdc_oldest_change_age_seconds": "12.5",
	"cdc_capture_lag":               "0.25",
	"cdc_consumer_position":         "9007199254740993",
	"cdc_consumer_lag":              "3",
	"cdc_changes_pending":           "4",
	"cdc_errors_total":              "1",
	"cdc_storage_bytes":             "8192",
	"cdc_wal_retention_bytes":       "1000",
}

func (f *fixture) build() []byte {
	var metrics []string
	for _, s := range specs {
		unit := s.unit
		if u, ok := f.unit[s.name]; ok {
			unit = u
		}
		attr := s.attr
		if a, ok := f.attrKey[s.name]; ok {
			attr = a
		}
		value := baseValues[s.name]
		if v, ok := f.value[s.name]; ok {
			value = v
		}
		number := fmt.Sprintf(`"asInt":"%s"`, value)
		if (s.unit == "s" || f.asDouble[s.name]) && !f.asInt[s.name] {
			number = fmt.Sprintf(`"asDouble":%s`, value)
		}
		attrs := ""
		switch s.attr {
		case "consumer":
			attrs = fmt.Sprintf(`"attributes":[{"key":"%s","value":{"stringValue":"consumer-a"}}],`, attr)
		case "class":
			attrs = fmt.Sprintf(`"attributes":[{"key":"%s","value":{"stringValue":"schema"}}],`, attr)
		}
		kind := "gauge"
		if f.gaugeName[s.name] {
			kind = "sum"
		}
		metrics = append(metrics, fmt.Sprintf(`{"name":"%s","unit":"%s","%s":{"dataPoints":[{%s"timeUnixNano":"1700000000000000000",%s}]}}`,
			s.name, unit, kind, attrs, number))
	}
	return []byte(fmt.Sprintf(`{"resourceMetrics":[{"resource":{"attributes":[%s]},"scopeMetrics":[{"scope":{"name":"pg-change-feed"},"metrics":[%s]}]}]}`,
		f.resource, strings.Join(metrics, ",")))
}

func baseView() map[string]string {
	view := map[string]string{}
	for _, s := range specs {
		if s.name == walName {
			continue
		}
		label := ""
		switch s.attr {
		case "consumer":
			label = "consumer-a"
		case "class":
			label = "schema"
		}
		view[s.name+"|"+label] = baseValues[s.name]
	}
	return view
}

var baseOpts = Options{Source: "src-1", ToleranceSeconds: 12, WALMeasured: []int64{0, 640, 1000}}

// TestCheckAcceptsMatchingExport belegt den Gut-Fall und die gedruckten Werte.
func TestCheckAcceptsMatchingExport(t *testing.T) {
	res, err := Check(newFixture().build(), baseView(), baseOpts)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if res.Names != 10 || res.ConsumerLag != "By" || res.WAL != "1000" || res.Points != 10 {
		t.Fatalf("Ergebnis = %+v", res)
	}
}

// TestCheckRejectsEachDeviation belegt je Abweichung den Befund mit Namen
// der Stelle: jede Mutation am Export bzw. an der Sicht färbt den Fall rot.
func TestCheckRejectsEachDeviation(t *testing.T) {
	cases := []struct {
		name string
		edit func(*fixture, map[string]string, *Options)
		want string
	}{
		{"Wert um eins verschoben", func(f *fixture, _ map[string]string, _ *Options) { f.value["cdc_changes_pending"] = "5" }, "cdc_changes_pending[consumer-a]: Export 5, Sicht 4"},
		{"Einheit geändert", func(f *fixture, _ map[string]string, _ *Options) { f.unit["cdc_consumer_lag"] = "1" }, `cdc_consumer_lag: Einheit "1", erwartet "By"`},
		{"Attribut class durch consumer ersetzt", func(f *fixture, _ map[string]string, _ *Options) { f.attrKey["cdc_errors_total"] = "consumer" }, "cdc_errors_total: Attribut class fehlt"},
		{"Sum statt Gauge", func(f *fixture, _ map[string]string, _ *Options) { f.gaugeName["cdc_storage_bytes"] = true }, "cdc_storage_bytes: Typ ist kein Gauge"},
		{"Quellkennung falsch", func(_ *fixture, _ map[string]string, o *Options) { o.Source = "src-2" }, `Resource-Attribut cdc.source_id = "src-1", erwartet "src-2"`},
		{"Zeile der Sicht fehlt im Export", func(_ *fixture, v map[string]string, _ *Options) { v["cdc_changes_pending|consumer-b"] = "1" }, "Zeile cdc_changes_pending|consumer-b der Sicht fehlt im Export"},
		{"Zeile des Exports fehlt in der Sicht", func(_ *fixture, v map[string]string, _ *Options) { delete(v, "cdc_storage_bytes|") }, "cdc_storage_bytes[]: keine Zeile in der Sicht"},
		{"Zeitabhängige Kennzahl außerhalb der Toleranz", func(f *fixture, _ map[string]string, _ *Options) { f.value["cdc_capture_lag"] = "40.5" }, "cdc_capture_lag[]: Export 40.5, Sicht 0.25"},
		{"WAL-Rückstand ist keine Messung des Prozesses", func(f *fixture, _ map[string]string, _ *Options) { f.value["cdc_wal_retention_bytes"] = "1001" }, "cdc_wal_retention_bytes = 1001, keine der 3 Messungen des Prozesses (größte 1000)"},
		{"WAL-Rückstand steht in der Sicht", func(_ *fixture, v map[string]string, _ *Options) { v["cdc_wal_retention_bytes|"] = "1000" }, "cdc_wal_retention_bytes steht in der Sicht"},
		{"Fließkommazahl für ganzzahligen Wert", func(f *fixture, _ map[string]string, _ *Options) { f.asDouble["cdc_changes_processed"] = true }, "cdc_changes_processed[]: ganzzahliger Wert geht als Fließkommazahl hinaus"},
		{"Ganze Zahl für Einheit s", func(f *fixture, _ map[string]string, _ *Options) { f.asInt["cdc_capture_lag"] = true }, "cdc_capture_lag[]: Einheit s geht als ganze Zahl hinaus"},
		{"Zusätzliches Resource-Attribut", func(f *fixture, _ map[string]string, _ *Options) {
			f.resource += `,{"key":"x","value":{"stringValue":"y"}}`
		}, "Resource-Attribute"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, view, opts := newFixture(), baseView(), baseOpts
			tc.edit(f, view, &opts)
			_, err := Check(f.build(), view, opts)
			if err == nil {
				t.Fatalf("Abweichung %q ohne Befund", tc.name)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("Befund %q nennt nicht %q", err.Error(), tc.want)
			}
		})
	}
}

// TestCheckAuthAttribute belegt beide Richtungen der Header-Beobachtung:
// ein erwarteter Wert muss ankommen, ein erwartet abwesendes Attribut darf
// nicht ankommen.
func TestCheckAuthAttribute(t *testing.T) {
	withAuth := newFixture()
	withAuth.resource += `,{"key":"http.authorization","value":{"stringValue":"Bearer tok"}}`
	opts := baseOpts
	opts.AuthAttr, opts.AuthValue = "http.authorization", "Bearer tok"
	if res, err := Check(withAuth.build(), baseView(), opts); err != nil || res.Auth != "present" {
		t.Fatalf("mit Header: res=%+v err=%v", res, err)
	}
	opts.AuthValue = "Bearer andere"
	if _, err := Check(withAuth.build(), baseView(), opts); err == nil || !strings.Contains(err.Error(), "erwartet") {
		t.Fatalf("falscher Wert: err=%v", err)
	}
	opts.AuthValue = ""
	if _, err := Check(withAuth.build(), baseView(), opts); err == nil || !strings.Contains(err.Error(), "erwartet abwesend") {
		t.Fatalf("Header da, erwartet abwesend: err=%v", err)
	}
	if res, err := Check(newFixture().build(), baseView(), opts); err != nil || res.Auth != "absent" {
		t.Fatalf("ohne Header: res=%+v err=%v", res, err)
	}
	opts.AuthValue = "Bearer tok"
	if _, err := Check(newFixture().build(), baseView(), opts); err == nil || !strings.Contains(err.Error(), "fehlt") {
		t.Fatalf("Header erwartet, nicht angekommen: err=%v", err)
	}
}

// TestLastRecordCountsOnlyCompleteLines belegt, dass eine Zeile ohne
// Zeilenende (noch im Schreiben) weder zählt noch als letzter Export gilt.
func TestLastRecordCountsOnlyCompleteLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.jsonl")
	if err := os.WriteFile(path, []byte("{\"a\":1}\n{\"b\":2}\n{\"c\":"), 0o600); err != nil {
		t.Fatal(err)
	}
	last, count, err := LastRecord(path)
	if err != nil || count != 2 || string(last) != `{"b":2}` {
		t.Fatalf("LastRecord = %q, %d, %v", last, count, err)
	}
	if _, count, _ = LastRecord(filepath.Join(t.TempDir(), "fehlt")); count != 0 {
		t.Fatalf("fehlende Datei zählt %d", count)
	}
}

// TestReadView belegt die Zerlegung der Sicht-Zeilen samt leerem Label.
func TestReadView(t *testing.T) {
	path := filepath.Join(t.TempDir(), "view.txt")
	if err := os.WriteFile(path, []byte("cdc_storage_bytes||8192\ncdc_changes_pending|c|4\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	rows, err := ReadView(path)
	if err != nil || rows["cdc_storage_bytes|"] != "8192" || rows["cdc_changes_pending|c"] != "4" {
		t.Fatalf("ReadView = %v, %v", rows, err)
	}
	if err := os.WriteFile(path, []byte("kaputt\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadView(path); err == nil {
		t.Fatal("Zeile ohne drei Felder ohne Fehler")
	}
}

// TestReadInts belegt das Lesen der gelogenen Messungen und die Ablehnung
// einer Zeile, die keine ganze Zahl ist.
func TestReadInts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wal.txt")
	if err := os.WriteFile(path, []byte("0\n4096\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	values, err := ReadInts(path)
	if err != nil || len(values) != 2 || values[1] != 4096 {
		t.Fatalf("ReadInts = %v, %v", values, err)
	}
	if err := os.WriteFile(path, []byte("zwölf\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadInts(path); err == nil {
		t.Fatal("Zeile ohne ganze Zahl ohne Fehler")
	}
}
