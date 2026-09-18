package main

import (
	"fmt"
	"os"
)

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

	skip, reason := decide(r)
	if skip {
		fmt.Fprintf(os.Stdout, "rolloutguard: %s — --execute wird uebersprungen\n", reason)
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "rolloutguard: %s — --execute laeuft regulaer\n", reason)
	os.Exit(1)
}
