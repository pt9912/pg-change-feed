package main

import "testing"

func noEnv(string) string { return "" }

// TestFlagsWireIntoTheStreamURL prüft die Verdrahtung Flag → Anfrage mit den
// echten Flag-Argumenten: `-schema`, `-table` und `-target` landen als
// Query-Parameter in der Stream-Adresse, ohne Flags bleibt die Adresse ohne
// Query.
func TestFlagsWireIntoTheStreamURL(t *testing.T) {
	const base = "http://feed:8080/changes/stream"
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"keine Flags", []string{"-addr", "feed:8080"}, base},
		{"alle drei", []string{"-addr", "feed:8080", "-schema", "public", "-table", "orders", "-target", "eu"}, base + "?schema=public&table=orders&target=eu"},
		{"nur table", []string{"-addr", "feed:8080", "-table", "orders"}, base + "?table=orders"},
		{"nur target", []string{"-addr", "feed:8080", "-target", "eu"}, base + "?target=eu"},
		{"Inline-Form", []string{"-addr=feed:8080", "-schema=public"}, base + "?schema=public"},
	} {
		cfg, err := parseConfig(tc.args, noEnv)
		if err != nil {
			t.Fatalf("%s: parseConfig: %v", tc.name, err)
		}
		if got := cfg.streamURL(); got != tc.want {
			t.Fatalf("%s: streamURL = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// TestParseConfigUsesEnvironmentAndRejectsUnknownFlags prüft die Defaults aus
// der Umgebung und die Ablehnung eines unbekannten Flags.
func TestParseConfigUsesEnvironmentAndRejectsUnknownFlags(t *testing.T) {
	env := func(name string) string {
		switch name {
		case "CDC_HTTP_ADDR":
			return "feed:8080"
		case "CDC_API_TOKEN_READER":
			return "tok-env"
		}
		return ""
	}
	cfg, err := parseConfig(nil, env)
	if err != nil {
		t.Fatalf("parseConfig: %v", err)
	}
	if cfg.addr != "feed:8080" || cfg.token != "tok-env" {
		t.Fatalf("cfg = %+v, want addr feed:8080 und token tok-env", cfg)
	}
	if _, err := parseConfig([]string{"-unknown"}, noEnv); err == nil {
		t.Fatal("parseConfig nahm ein unbekanntes Flag an")
	}
}
