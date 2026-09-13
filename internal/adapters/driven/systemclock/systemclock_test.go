package systemclock_test

import (
	"testing"
	"time"

	"github.com/pt9912/pg-change-feed/internal/adapters/driven/systemclock"
	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// TestAdapterImplementsClockPort belegt die Port-Kante (`ARC-004`,
// `ADR-0040`): der Adapter erfüllt `outbound.ClockPort` als Wert, keinen
// Zeiger nötig — er trägt keinen Zustand.
func TestAdapterImplementsClockPort(t *testing.T) {
	var port outbound.ClockPort = systemclock.New()
	_ = port
}

// TestNowReflectsWallClock belegt, dass `Now` die reale Wanduhr liest
// (Unix-Nanosekunden) statt eines festen oder abgeleiteten Werts: der
// gemeldete Zeitpunkt liegt innerhalb einer großzügigen Toleranz um
// `time.Now()` zum Zeitpunkt des Aufrufs.
func TestNowReflectsWallClock(t *testing.T) {
	adapter := systemclock.New()
	before := time.Now().UnixNano()
	got := adapter.Now()
	after := time.Now().UnixNano()

	if got.UnixNanos < before || got.UnixNanos > after {
		t.Fatalf("Now = %d, wollen einen Wert zwischen %d und %d", got.UnixNanos, before, after)
	}
}

// TestNowIsNonDecreasing belegt, dass zwei aufeinanderfolgende Aufrufe
// keinen rückwärts laufenden Zeitpunkt melden — dieselbe Eigenschaft, die
// `RetentionPolicy.AllowsDeletion` über das Alter voraussetzt.
func TestNowIsNonDecreasing(t *testing.T) {
	adapter := systemclock.New()
	first := adapter.Now()
	second := adapter.Now()
	if second.Before(first) {
		t.Fatalf("zweiter Aufruf (%d) liegt vor dem ersten (%d)", second.UnixNanos, first.UnixNanos)
	}
}
