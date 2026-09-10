// cmd/pg-change-feed/main.go — Einstiegspunkt des Binarys. `main` bleibt
// dünn: er liest die Verdrahtungs-Vorbedingungen aus der Umgebung und
// trägt den Prozess-Ausgang; die Verdrahtung der konkreten Adapter liegt
// an genau einer Stelle im Bootstrap (`ADR-0026`).
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pt9912/pg-change-feed/internal/bootstrap"
)

// version trägt den Lieferstand des Binarys.
const version = "0.2.0-verdrahtung"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf("pg-change-feed %s\n", version)
		return
	}
	if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, "pg-change-feed: unbekanntes Argument; der CDC-Lauf läuft ohne Argumente, --version zeigt den Lieferstand")
		os.Exit(2)
	}
	// Der Lauf endet kontrolliert auf SIGINT/SIGTERM: der Stream-Lauf
	// kehrt ohne Fehler zurück und der Prozess trägt Ausgang 0.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg, err := bootstrap.ConfigFromEnv(os.Getenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: %v\n", err)
		os.Exit(2)
	}
	if err := bootstrap.Run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "pg-change-feed: %v\n", err)
		os.Exit(1)
	}
}
