package receive

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

// TestServerFaultKlassifiziertNachSQLState bindet die Einwicklung an den
// SQLSTATE der Eingabe: Berechtigungsfehler tragen ErrPermission und nicht
// ErrReplication, nicht transiente Abweisungen tragen ErrRejected, transiente
// Zustände und Fehler ohne SQLSTATE tragen nur ErrReplication.
func TestServerFaultKlassifiziertNachSQLState(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		permission bool
		rejected   bool
		replicat   bool
	}{
		{"42501 insufficient_privilege", &pgconn.PgError{Code: "42501"}, true, false, false},
		{"28000 invalid_authorization", &pgconn.PgError{Code: "28000"}, true, false, false},
		{"28P01 invalid_password", &pgconn.PgError{Code: "28P01"}, true, false, false},
		{"42704 undefined_object", &pgconn.PgError{Code: "42704"}, false, true, true},
		{"XX000 internal_error", &pgconn.PgError{Code: "XX000"}, false, true, true},
		{"55006 object_in_use", &pgconn.PgError{Code: "55006"}, false, false, true},
		{"57P01 admin_shutdown", &pgconn.PgError{Code: "57P01"}, false, false, true},
		{"08006 connection_failure", &pgconn.PgError{Code: "08006"}, false, false, true},
		{"53300 too_many_connections", &pgconn.PgError{Code: "53300"}, false, false, true},
		{"umwickelter PgError", fmt.Errorf("Ebene: %w", &pgconn.PgError{Code: "42501"}), true, false, false},
		{"Fehler ohne SQLSTATE", errors.New("connection reset by peer"), false, false, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := serverFault("Test", c.err)
			if errors.Is(got, ErrPermission) != c.permission {
				t.Errorf("ErrPermission = %v, erwartet %v: %v", !c.permission, c.permission, got)
			}
			if errors.Is(got, ErrRejected) != c.rejected {
				t.Errorf("ErrRejected = %v, erwartet %v: %v", !c.rejected, c.rejected, got)
			}
			if errors.Is(got, ErrReplication) != c.replicat {
				t.Errorf("ErrReplication = %v, erwartet %v: %v", !c.replicat, c.replicat, got)
			}
		})
	}
}
