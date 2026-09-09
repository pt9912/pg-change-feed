// cmd/pg-change-feed/main.go — Einstiegspunkt des Binarys (Composition-Root-
// Komponente der Architektur-Sicht). Das Binary verdrahtet keine Adapter:
// es kennt ausschließlich das --version-Flag, das den Container-Build-Vertrag
// beobachtbar macht (reproduzierbares Deployment).
package main

import (
	"fmt"
	"os"
)

// version trägt den Lieferstand des Binarys: der Bootstrap liefert den
// Build-Vertrag, keine fachliche Fähigkeit.
const version = "0.1.0-bootstrap"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf("pg-change-feed %s\n", version)
		return
	}
	fmt.Fprintln(os.Stderr, "pg-change-feed: keine Argumente unterstützt; --version zeigt den Lieferstand")
	os.Exit(2)
}
