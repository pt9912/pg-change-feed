package main

import (
	"fmt"
	"os"
)

// Ausgabevertrag (das Makefile-Target schema-rollout liest ihn): stdout
// trägt je erlaubter Erweiterung genau eine Zeile — `allow-destructive`
// bzw. `drop-view <name>` —, stderr die Begründung. Exit 0: mindestens eine
// Erweiterung ist erlaubt; 1: keine; 2: Report nicht lesbar/dekodierbar.
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: rolloutguard <plan-only-report.json>")
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "rolloutguard: Report nicht lesbar: %v\n", err)
		os.Exit(2)
	}

	r, err := parseReport(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rolloutguard: Report nicht dekodierbar: %v\n", err)
		os.Exit(2)
	}

	d := decide(r)
	if !d.allowDestructive && len(d.dropViews) == 0 {
		fmt.Fprintf(os.Stderr, "rolloutguard: %s — --execute laeuft ohne --allow-destructive und ohne Vorlauf\n", d.reason)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "rolloutguard: %s\n", d.reason)
	if d.allowDestructive {
		fmt.Fprintln(os.Stdout, "allow-destructive")
	}
	for _, v := range d.dropViews {
		fmt.Fprintf(os.Stdout, "drop-view %s\n", v)
	}
	os.Exit(0)
}
