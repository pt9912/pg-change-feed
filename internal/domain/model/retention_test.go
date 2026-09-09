package model

import (
	"testing"

	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// LH-FA-RET-004 / ADR-0029, Regel 5: Retention löscht keine benötigten
// Changes — eine unbestätigte Consumer-Position oder eine Position hinter
// der Change-Position blockiert die Bereinigung.
func TestLHFARET004RetentionRequiresAcknowledgedPositions(t *testing.T) {
	policy, err := NewRetentionPolicy(Duration{})
	if err != nil {
		t.Fatalf("Policy: %v", err)
	}
	changePosition := mustSourcePosition(t, "src-1", 100)
	consumer, err := NewConsumerPosition("c-1")
	if err != nil {
		t.Fatalf("Consumer-Position: %v", err)
	}
	confirmed, err := consumer.Advance(mustSourcePosition(t, "src-1", 200))
	if err != nil {
		t.Fatalf("Bestätigung: %v", err)
	}
	if !policy.AllowsDeletion(Duration{Nanos: 1}, changePosition, []ConsumerPosition{confirmed}) {
		t.Fatal("Bereinigung hinter der bestätigten Position blockiert")
	}
	if policy.AllowsDeletion(Duration{Nanos: 1}, changePosition, []ConsumerPosition{consumer}) {
		t.Fatal("Bereinigung mit unbestätigter Consumer-Position freigegeben")
	}
	behind, err := NewConsumerPosition("c-2")
	if err != nil {
		t.Fatalf("Consumer-Position: %v", err)
	}
	behind, err = behind.Advance(mustSourcePosition(t, "src-1", 50))
	if err != nil {
		t.Fatalf("Bestätigung: %v", err)
	}
	if policy.AllowsDeletion(Duration{Nanos: 1}, changePosition, []ConsumerPosition{behind}) {
		t.Fatal("Bereinigung hinter einer benötigten Position freigegeben")
	}
}

// LH-FA-RET-003: unter dem Mindestalter gibt die Policy keine Bereinigung
// frei; ab dem Alter schon — auch ohne blockierende Consumer-Positionen.
func TestLHFARET003RetentionRequiresMinimumAge(t *testing.T) {
	policy, err := NewRetentionPolicy(Duration{Nanos: 1000})
	if err != nil {
		t.Fatalf("Policy: %v", err)
	}
	changePosition := mustSourcePosition(t, "src-1", 100)
	if policy.AllowsDeletion(Duration{Nanos: 999}, changePosition, nil) {
		t.Fatal("Bereinigung unter dem Mindestalter freigegeben")
	}
	if !policy.AllowsDeletion(Duration{Nanos: 1000}, changePosition, nil) {
		t.Fatal("Bereinigung am Mindestalter blockiert")
	}
}

// Die Policy verlangt ein nichtnegatives Mindestalter (`LH-FA-RET-003`).
func TestNewRetentionPolicyRejectsNegativeDuration(t *testing.T) {
	if _, err := NewRetentionPolicy(Duration{Nanos: -1}); !stderrors.Is(err, domainerrors.ErrNegativeDuration) {
		t.Fatalf("Fehler = %v, wollen ErrNegativeDuration", err)
	}
	if _, err := NewDuration(-1); !stderrors.Is(err, domainerrors.ErrNegativeDuration) {
		t.Fatalf("Fehler = %v, wollen ErrNegativeDuration", err)
	}
}
