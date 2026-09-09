package postgresstorage

import (
	"context"
	_ "embed"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
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
func ApplySchema(ctx context.Context, pool *pgxpool.Pool) error {
	for _, statement := range statements(schemaSQL) {
		if _, err := pool.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

// statements teilt die DDL in Einzelanweisungen: pgx führt über das
// Extended-Protokoll eine Anweisung je Exec aus; die DDL trägt keine
// Semikolons in String-Literalen, das einfache Trennen trägt ihren vollen
// Text.
func statements(sql string) []string {
	var result []string
	for _, part := range strings.Split(sql, ";") {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}
