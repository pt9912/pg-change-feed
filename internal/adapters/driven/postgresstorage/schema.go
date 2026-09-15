package postgresstorage

import (
	"context"
	_ "embed"
	"strings"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage/sqlexec"
)

// schema.sql trägt die CDC-Tabellen-DDL (Pflichtenheft §2); sie wird
// eingebettet und über ApplySchema gegen die verbundene Instanz ausgeführt.
//
//go:embed schema.sql
var schemaSQL string

// SchemaSQL trägt die eingebettete Schema-DDL; sie bleibt unverändert
// ausgeliefert, die Datei ist ihre einzige Quelle.
func SchemaSQL() string {
	return schemaSQL
}

// ApplySchema legt die CDC-Tabellen an (`Pflichtenheft §2`). Die DDL ist
// idempotent (IF NOT EXISTS); der Aufruf auf einer Instanz mit Bestand
// ändert diesen nicht. Der Aufrufer trägt die Registrierung von Quelle,
// Tabelle und Schema-Version vor der ersten Persistenz — die
// Fremdschlüssel der DDL setzen sie voraus (schema.sql, Kopf-Kommentar).
// Der Träger kommt über die schmale Ausführungs-Naht (`sqlexec`,
// `ADR-0071` Punkt 5) — der reale Pool erfüllt sie ohne Vermittler.
func ApplySchema(ctx context.Context, exec sqlexec.Executor) error {
	for _, statement := range statements(schemaSQL) {
		if _, err := exec.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

// statements teilt die DDL in Einzelanweisungen: pgx führt über das
// Extended-Protokoll eine Anweisung je Exec aus. Kommentarzeilen werden
// vor dem Trennen abgetrennt — ihre Texte können Semikolons tragen —; die
// SQL-Blöcke selbst tragen keine Semikolons in String-Literalen.
func statements(sql string) []string {
	var withoutComments []string
	for _, line := range strings.Split(sql, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}
		withoutComments = append(withoutComments, line)
	}
	var result []string
	for _, part := range strings.Split(strings.Join(withoutComments, "\n"), ";") {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}
