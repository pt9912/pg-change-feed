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
	"strconv"
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
	if len(os.Args) == 2 && os.Args[1] == "--healthcheck" {
		// Der Compose-Healthcheck des Feed-Containers (`compose.yaml`,
		// `LH-FA-ADM-002`, `LH-QA-OPS-002`): das Runtime-Image
		// ist distroless (kein Shell, kein `psql`, Dockerfile) — der
		// einzige Aufruf, den Compose innerhalb dieses Containers
		// ausführen kann, ist das Binary selbst (`CMD`-Form ohne Shell).
		// Dieselben Umgebungs-Vorbedingungen wie der reguläre Lauf
		// (`ConfigFromEnv`) tragen die drei rollen-spezifischen DSNs und
		// die Quelle (`ADR-0047`) — der Healthcheck-Lauf selbst nutzt nur
		// `cfg.ReaderDSN`; Publication/Slot/Tabellen bleiben ungenutzt, die
		// Vorbedingungsprüfung teilt sich beide Läufe trotzdem, statt eine
		// zweite Lese-Funktion zu pflegen.
		cfg, err := bootstrap.ConfigFromEnv(os.Getenv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pg-change-feed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(bootstrap.Healthcheck(context.Background(), cfg.ReaderDSN, cfg.Source))
	}
	if len(os.Args) >= 2 && os.Args[1] == "register-consumer" {
		// Der Sondermodus registriert einen Consumer über
		// `RegisterConsumerUseCase` (`LH-FA-CON-001.a`) und beendet sich,
		// ohne je den Capture-Loop (`bootstrap.Run`) zu erreichen —
		// dasselbe Muster wie `--healthcheck` oben. Kennung und Name des
		// Consumers tragen denselben Wert; dieser Zugriffsweg trennt
		// beide (noch) nicht.
		if len(os.Args) != 3 {
			fmt.Fprintln(os.Stderr, "pg-change-feed: register-consumer erwartet genau einen Namen als Argument")
			os.Exit(2)
		}
		cfg, err := bootstrap.ConfigFromEnv(os.Getenv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pg-change-feed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(bootstrap.RegisterConsumer(context.Background(), cfg, os.Args[2]))
	}
	if len(os.Args) >= 2 && os.Args[1] == "acknowledge-consumer" {
		// Der Sondermodus bestätigt die Position eines Consumers über
		// `AcknowledgeConsumerUseCase` (`LH-FA-CON-004.a`) und beendet
		// sich, ohne je den Capture-Loop (`bootstrap.Run`) zu erreichen —
		// dasselbe Muster wie `register-consumer` oben. Die Quelle der
		// bestätigten Position ist die konfigurierte `CDC_SOURCE_ID`
		// (`cfg.Source`); das zweite Argument trägt nur den Offset
		// innerhalb dieser Quelle.
		if len(os.Args) != 4 {
			fmt.Fprintln(os.Stderr, "pg-change-feed: acknowledge-consumer erwartet Consumer-Kennung und Position als Argumente")
			os.Exit(2)
		}
		offset, parseErr := strconv.ParseUint(os.Args[3], 10, 64)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "pg-change-feed: acknowledge-consumer: Position %q ist kein gültiger Offset: %v\n", os.Args[3], parseErr)
			os.Exit(2)
		}
		cfg, err := bootstrap.ConfigFromEnv(os.Getenv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pg-change-feed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(bootstrap.AcknowledgeConsumer(context.Background(), cfg, os.Args[2], offset))
	}
	if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, "pg-change-feed: unbekanntes Argument; der CDC-Lauf läuft ohne Argumente, --version und --healthcheck zeigen bzw. prüfen den Lieferstand, register-consumer <name> registriert einen Consumer, acknowledge-consumer <consumer-id> <position> bestätigt eine Position")
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
