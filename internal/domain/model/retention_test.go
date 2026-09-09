package model

import (
	"testing"

	stderrors "errors"

	domainerrors "github.com/pt9912/pg-change-feed/internal/domain/errors"
)

// LH-FA-RET-004 / ADR-0029, Regel 5: Retention löscht keine benötigten
// Changes — ohne die Bestätigung aller relevanten Consumer gibt die Policy
// keine Bereinigung frei.
func TestLHFARET004RetentionRequiresAllConsumersAcknowledged(t *testing.T) {
	policy, err := NewRetentionPolicy(Duration{})
	if err != nil {
		t.Fatalf("Policy: %v", err)
	}
	age, err := NewDuration(1000)
	if err != nil {
		t.Fatalf("Alter: %v", err)
	}
	if policy.AllowsDeletion(age, false) {
		t.Fatal("Bereinigung ohne Bestätigung aller Consumer freigegeben")
	}
	if !policy.AllowsDeletion(age, true) {
		t.Fatal("Bereinigung mit Bestätigung aller Consumer blockiert")
	}
}

// LH-FA-RET-003: unter dem Mindestalter gibt die Policy keine Bereinigung
// frei; ab dem Alter schon.
func TestLHFARET003RetentionRequiresMinimumAge(t *testing.T) {
	policy, err := NewRetentionPolicy(Duration{Nanos: 1000})
	if err != nil {
		t.Fatalf("Policy: %v", err)
	}
	tooYoung, err := NewDuration(999)
	if err != nil {
		t.Fatalf("zu junges Alter: %v", err)
	}
	oldEnough, err := NewDuration(1000)
	if err != nil {
		t.Fatalf("altes genug Alter: %v", err)
	}
	if policy.AllowsDeletion(tooYoung, true) {
		t.Fatal("Bereinigung unter dem Mindestalter freigegeben")
	}
	if !policy.AllowsDeletion(oldEnough, true) {
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
