// Command otlpcheck ist das Wegwerf-Werkzeug des OTLP-Export-Belegs in
// tools/harness/run-integration-tests.sh (LH-FA-SST-010): es liest die Datei,
// die der `file`-Exporter eines echten OpenTelemetry-Collectors schreibt (eine
// JSON-Zeile je empfangenem Export), und die Zeilen der Sicht `cdc.metrics`
// und prüft, was am Ausgang des Empfängers ankommt.
//
//	otlpcheck check -collector <Datei> -view <Datei> -source <Kennung>
//	    [-tolerance-seconds <n>] [-wal-log <Datei>]
//	    [-auth-attr <Schlüssel> -auth-value <Wert>]
//	otlpcheck count -collector <Datei>
//
// `check` prüft den letzten vollständigen Export: zehn Namen, Typ `Gauge`,
// Einheit, Datenpunkt-Attribut, Resource-Attribute und Wertgleichheit mit der
// Sicht (`cdc_wal_retention_bytes` steht nicht in der Sicht: sein Wert ist
// genau eine der Messungen, die der Prozess in `-wal-log` geloggt hat); ein
// Befund endet mit Exit 1 und einer Zeile je Abweichung auf stderr.
// `count` druckt die Zahl der vollständigen Exporte.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: otlpcheck check|count ...")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "check":
		os.Exit(runCheck(os.Args[2:]))
	case "count":
		os.Exit(runCount(os.Args[2:]))
	default:
		fmt.Fprintf(os.Stderr, "otlpcheck: unbekannter Modus %q\n", os.Args[1])
		os.Exit(2)
	}
}

func runCount(args []string) int {
	fs := flag.NewFlagSet("count", flag.ContinueOnError)
	collector := fs.String("collector", "", "Datei des Collectors")
	if err := fs.Parse(args); err != nil || *collector == "" {
		return 2
	}
	_, count, err := LastRecord(*collector)
	if err != nil {
		fmt.Fprintf(os.Stderr, "otlpcheck: %v\n", err)
		return 1
	}
	fmt.Printf("COUNT %d\n", count)
	return 0
}

func runCheck(args []string) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	collector := fs.String("collector", "", "Datei des Collectors")
	view := fs.String("view", "", "Zeilen der Sicht (name|label|wert)")
	source := fs.String("source", "", "erwartete Quellkennung")
	tolerance := fs.Float64("tolerance-seconds", 12, "Toleranz der Kennzahlen mit Einheit s")
	walLog := fs.String("wal-log", "", "Datei mit den gelogenen Messungen des WAL-Rückstands, eine Zahl je Zeile")
	authAttr := fs.String("auth-attr", "", "Resource-Attribut, das der Collector aus dem Header bildet")
	authValue := fs.String("auth-value", "", "erwarteter Wert; leer: das Attribut fehlt")
	if err := fs.Parse(args); err != nil || *collector == "" || *view == "" || *source == "" {
		fmt.Fprintln(os.Stderr, "otlpcheck check: -collector, -view und -source sind Pflicht")
		return 2
	}
	last, _, err := LastRecord(*collector)
	if err != nil {
		fmt.Fprintf(os.Stderr, "otlpcheck: %v\n", err)
		return 1
	}
	if last == nil {
		fmt.Fprintln(os.Stderr, "otlpcheck: die Collector-Datei trägt keinen Export")
		return 1
	}
	rows, err := ReadView(*view)
	if err != nil {
		fmt.Fprintf(os.Stderr, "otlpcheck: %v\n", err)
		return 1
	}
	opts := Options{
		Source: *source, ToleranceSeconds: *tolerance,
		AuthAttr: *authAttr, AuthValue: *authValue,
	}
	if *walLog != "" {
		if opts.WALMeasured, err = ReadInts(*walLog); err != nil {
			fmt.Fprintf(os.Stderr, "otlpcheck: %v\n", err)
			return 1
		}
		if len(opts.WALMeasured) == 0 {
			fmt.Fprintln(os.Stderr, "otlpcheck: der Prozess hat keine Messung des WAL-Rückstands geloggt")
			return 1
		}
	}
	res, err := Check(last, rows, opts)
	if err != nil {
		for _, line := range strings.Split(err.Error(), "\n") {
			fmt.Fprintf(os.Stderr, "otlpcheck: %s\n", line)
		}
		return 1
	}
	fmt.Printf("OTLPCHECK names=%d points=%d consumer_lag_unit=%s wal=%s errors_labels=%s auth=%s\n",
		res.Names, res.Points, res.ConsumerLag, res.WAL, strings.Join(res.ErrorsLabels, ","), res.Auth)
	return 0
}
