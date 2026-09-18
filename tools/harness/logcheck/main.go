// Command logcheck ist ein Wegwerf-Testwerkzeug für den E2E-Beleg der
// maschinenlesbaren Log-Struktur (LH-QA-OPS-004, ADR-0024): es liest die
// Zeilen von stdin — der Aufrufer reicht `docker logs` des laufenden
// Feed-Containers ein — und prüft jede nicht-leere Zeile als eigenständiges
// JSON-Objekt mit den drei vom `slog.NewJSONHandler`-Standard garantierten
// Feldern `time`, `level`, `msg`. Träger ist
// tools/harness/run-integration-tests.sh — der Aufrufer liest die
// stdout-Zeile dieses Prozesses.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lines := 0
	levels := map[string]bool{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var decoded map[string]any
		if err := json.Unmarshal([]byte(line), &decoded); err != nil {
			fmt.Fprintf(os.Stderr, "logcheck: Zeile %d ist kein gültiges JSON: %v: %s\n", lines+1, err, line)
			os.Exit(1)
		}
		for _, field := range []string{"time", "level", "msg"} {
			value, ok := decoded[field]
			if !ok {
				fmt.Fprintf(os.Stderr, "logcheck: Zeile %d trägt kein Feld %q: %s\n", lines+1, field, line)
				os.Exit(1)
			}
			if text, ok := value.(string); !ok || text == "" {
				fmt.Fprintf(os.Stderr, "logcheck: Zeile %d trägt Feld %q nicht als nicht-leeren String: %s\n", lines+1, field, line)
				os.Exit(1)
			}
		}
		lines++
		levels[decoded["level"].(string)] = true
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "logcheck: stdin-Lesefehler: %v\n", err)
		os.Exit(1)
	}
	if lines == 0 {
		fmt.Fprintln(os.Stderr, "logcheck: keine Zeile gelesen — leere Eingabe belegt keine Struktur")
		os.Exit(1)
	}

	sortedLevels := make([]string, 0, len(levels))
	for level := range levels {
		sortedLevels = append(sortedLevels, level)
	}
	sort.Strings(sortedLevels)

	fmt.Printf("STRUCTURED_LOG_OK lines=%d levels=%s\n", lines, strings.Join(sortedLevels, ","))
}
