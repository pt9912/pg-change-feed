package outbound_test

import (
	"context"
	"testing"

	"github.com/pt9912/pg-change-feed/internal/application/port/outbound"
)

// fakeLog trägt einen Testdouble für `LogPort` — genau die Substitution,
// die ein globaler `slog`-Singleton nicht erlaubt (`ADR-0024`): der Port
// ist eine Schnittstelle, kein paketweiter Zustand.
type fakeLog struct {
	messages []string
}

func (f *fakeLog) Debug(_ context.Context, msg string, _ ...any) { f.messages = append(f.messages, "DEBUG:"+msg) }
func (f *fakeLog) Info(_ context.Context, msg string, _ ...any)  { f.messages = append(f.messages, "INFO:"+msg) }
func (f *fakeLog) Warn(_ context.Context, msg string, _ ...any)  { f.messages = append(f.messages, "WARN:"+msg) }
func (f *fakeLog) Error(_ context.Context, msg string, _ ...any) { f.messages = append(f.messages, "ERROR:"+msg) }

// TestLogPortAcceptsSubstitution trägt die Testdouble-Fähigkeit, die
// `ADR-0024` für Telemetrie verlangt: ein Aufrufer, der gegen `LogPort`
// programmiert, läuft unverändert gegen einen Fake.
func TestLogPortAcceptsSubstitution(t *testing.T) {
	fake := &fakeLog{}
	var port outbound.LogPort = fake
	port.Debug(context.Background(), "d")
	port.Info(context.Background(), "i")
	port.Warn(context.Background(), "w")
	port.Error(context.Background(), "e")
	want := []string{"DEBUG:d", "INFO:i", "WARN:w", "ERROR:e"}
	if len(fake.messages) != len(want) {
		t.Fatalf("Aufrufe: %v, wollen %v", fake.messages, want)
	}
	for i := range want {
		if fake.messages[i] != want[i] {
			t.Fatalf("Aufruf %d: %s, wollen %s", i, fake.messages[i], want[i])
		}
	}
}

// TestNoopLogSwallowsCalls trägt den Default für Aufrufer ohne injizierten
// Adapter (Tests, `postgresstorage`/`postgresack` ohne `WithLog`-Option):
// alle vier Stufen laufen ohne Wirkung und ohne Panic.
func TestNoopLogSwallowsCalls(t *testing.T) {
	ctx := context.Background()
	outbound.NoopLog.Debug(ctx, "d", "k", "v")
	outbound.NoopLog.Info(ctx, "i")
	outbound.NoopLog.Warn(ctx, "w")
	outbound.NoopLog.Error(ctx, "e", "err", nil)
}
