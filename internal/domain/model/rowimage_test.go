package model

import (
	"bytes"
	"sync"
	"testing"
)

// Die Erwartungswerte der Byte-Tabelle sind Referenz-Bytes, gemessen am
// Bild-Erzeuger des Replication-Mappers vor seinem Umzug in die Domäne
// (`ADR-0111` Teilfrage 2: WAL- und Backfill-Pfad teilen eine Funktion,
// das Bild bleibt byte-gleich). Sonderzeichen stehen in Hex-Escapes, damit
// kein Editor oder Werkzeug sie beim Speichern umschreibt: `\x5c` ist ein
// Backslash.

func textValue(value string) *string { return &value }

var rowImageColumns = []string{"id", "secret", "name", "a\"b<"}

func TestBuildRowImageBytes(t *testing.T) {
	cases := []struct {
		name     string
		values   []*string
		excluded []string
		want     string
	}{
		{
			name:   "leere Werteliste ohne tragenden Wert",
			values: []*string{},
			want:   `{}`,
		},
		{
			name:   "NULL-Wert entfällt",
			values: []*string{textValue("1"), nil, textValue("x"), nil},
			want:   `{"id":"1","name":"x"}`,
		},
		{
			name:     "ausgeschlossene Spalte entfällt samt Wert",
			values:   []*string{textValue("1"), textValue("geheim"), textValue("x"), textValue("y")},
			excluded: []string{"secret"},
			want:     "{\"id\":\"1\",\"name\":\"x\",\"a\\\"b\x5cu003c\":\"y\"}",
		},
		{
			name:   "alle Werte abwesend",
			values: []*string{nil, nil, nil, nil},
			want:   `{}`,
		},
		{
			name:   "leerer String ist Wert, keine Abwesenheit",
			values: []*string{textValue(""), nil, nil, nil},
			want:   `{"id":""}`,
		},
		{
			name:   "weniger Werte als Spalten",
			values: []*string{textValue("1"), textValue("s")},
			want:   `{"id":"1","secret":"s"}`,
		},
		{
			name:   "mehr Werte als Spalten",
			values: []*string{textValue("1"), textValue("s"), textValue("n"), textValue("q"), textValue("extra")},
			want:   "{\"id\":\"1\",\"secret\":\"s\",\"name\":\"n\",\"a\\\"b\x5cu003c\":\"q\"}",
		},
		{
			name:   "Maskierung von <, >, &, Anführungszeichen, Umlaut, Steuerzeichen, Backslash, U+2028",
			values: []*string{textValue("<a>&\"b\" \xc3\xbc \x01\n\t\\ \xe2\x80\xa8"), nil, nil, nil},
			want:   "{\"id\":\"\x5cu003ca\x5cu003e\x5cu0026\\\"b\\\" \xc3\xbc \x5cu0001\x5cn\x5ct\x5c\x5c \x5cu2028\"}",
		},
		{
			name:   "ungültiges UTF-8 wird zum Ersatzzeichen",
			values: []*string{textValue("\xff"), nil, nil, nil},
			want:   "{\"id\":\"\xef\xbf\xbd\"}",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := BuildRowImage(rowImageColumns, c.values, c.excluded)
			if err != nil {
				t.Fatalf("BuildRowImage: %v", err)
			}
			if got == nil {
				t.Fatal("Bild ist nil, wollen ein Objekt")
			}
			if string(got) != c.want {
				t.Fatalf("Bild %q, wollen %q", got, c.want)
			}
		})
	}
}

// Eine nil-Werteliste ist Abwesenheit des Bildes (`LH-FA-CAP-008`): weder
// Objekt noch Fehler — unterscheidbar von der leeren Liste (`{}`).
func TestBuildRowImageNilValuesIsAbsent(t *testing.T) {
	got, err := BuildRowImage(rowImageColumns, nil, nil)
	if err != nil {
		t.Fatalf("BuildRowImage: %v", err)
	}
	if got != nil {
		t.Fatalf("Bild %q, wollen nil", got)
	}
}

// Der Ausschluss trägt seine Wirkung an der Spalte, nicht am Wert
// (`LH-QA-SEC-004`): der Schlüssel erscheint nicht, der Wert nirgends, und
// eine nicht geführte Spalte bleibt.
func TestBuildRowImageExcludedValueNowhere(t *testing.T) {
	values := []*string{textValue("1"), textValue("streng-geheim"), textValue("x"), nil}
	got, err := BuildRowImage(rowImageColumns, values, []string{"secret", "unbekannt"})
	if err != nil {
		t.Fatalf("BuildRowImage: %v", err)
	}
	if bytes.Contains(got, []byte("streng-geheim")) || bytes.Contains(got, []byte(`"secret"`)) {
		t.Fatalf("Bild %q trägt die ausgeschlossene Spalte", got)
	}
	if string(got) != `{"id":"1","name":"x"}` {
		t.Fatalf("Bild %q, wollen {\"id\":\"1\",\"name\":\"x\"}", got)
	}
}

// Die Funktion ist rein: die Eingaben bleiben unverändert, das Ergebnis
// teilt keinen Speicher mit ihnen oder mit einem früheren Ergebnis, und
// gleichzeitige Aufrufe (`-race`) liefern dieselben Bytes.
func TestBuildRowImagePureAndConcurrent(t *testing.T) {
	columns := []string{"id", "secret", "name"}
	excluded := []string{"secret"}
	values := []*string{textValue("1"), textValue("s"), textValue("x")}
	first, err := BuildRowImage(columns, values, excluded)
	if err != nil {
		t.Fatalf("BuildRowImage: %v", err)
	}
	want := string(first)
	first[0] = '!'
	second, err := BuildRowImage(columns, values, excluded)
	if err != nil {
		t.Fatalf("BuildRowImage: %v", err)
	}
	if string(second) != want {
		t.Fatalf("zweites Bild %q, wollen %q (Ergebnis teilt Speicher)", second, want)
	}
	if columns[0] != "id" || excluded[0] != "secret" || *values[0] != "1" || len(values) != 3 {
		t.Fatal("Eingabe verändert")
	}
	var wg sync.WaitGroup
	results := make([][]byte, 8)
	for i := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i], _ = BuildRowImage(columns, values, excluded)
		}()
	}
	wg.Wait()
	for i, result := range results {
		if string(result) != want {
			t.Fatalf("Aufruf %d: Bild %q, wollen %q", i, result, want)
		}
	}
}

func BenchmarkBuildRowImage(b *testing.B) {
	columns := []string{"id", "secret", "name", "email", "created_at", "payload"}
	values := []*string{
		textValue("42"), textValue("geheim"), textValue("Anna Beispiel"),
		textValue("anna@example.org"), textValue("2026-09-24 10:00:00+00"), textValue(`{"k":"<v>"}`),
	}
	excluded := []string{"secret"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := BuildRowImage(columns, values, excluded); err != nil {
			b.Fatal(err)
		}
	}
}
