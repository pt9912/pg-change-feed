package main

import (
	"fmt"
	"strings"
)

// allSubject ist der Wurzel-Wildcard des Vollinhalts-Namensraums
// (`ADR-0100`) und das Subjekt ohne `-target`.
const allSubject = "cdc.stream.>"

// invalidTokenChars trägt die Zeichen, die in einem Subjekt-Token ein
// Trennzeichen oder einen NATS-Platzhalter bilden und das Muster still
// veränderten.
const invalidTokenChars = ".*> \t\n\r"

// SubscribeSubject leitet das Abonnement-Subjekt aus `-source` und `-target`
// ab (`SPEC-024`): ohne beide der Wurzel-Wildcard, mit beiden das
// Zusatz-Subjekt eines Zustellziels `cdc.route.<source_id>.<ziel>`. Eines der
// beiden allein ist ein Eingabefehler — das Zusatz-Subjekt hängt an der Quelle
// und am Ziel zugleich. Ein Token mit Punkt, Platzhalter oder Leerraum wird
// abgelehnt, bevor ein Abonnement entsteht.
func SubscribeSubject(source, target string) (string, error) {
	switch {
	case source == "" && target == "":
		return allSubject, nil
	case source == "" || target == "":
		return "", fmt.Errorf("-source und -target gelten nur zusammen")
	}
	for name, value := range map[string]string{"-source": source, "-target": target} {
		if strings.TrimSpace(value) == "" || strings.ContainsAny(value, invalidTokenChars) {
			return "", fmt.Errorf("%s darf weder leer sein noch '.', '*', '>' oder Leerraum enthalten", name)
		}
	}
	return "cdc.route." + source + "." + target, nil
}
