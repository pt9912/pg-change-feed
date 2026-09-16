// main_test.go trägt den Argument-Dispatch von `main` als Gegenstand:
// seine Zusage ist der **Ausgang des Prozesses** und die Ausgabe, die er
// dazu trägt — nicht ein Rückgabewert, den ein Test im eigenen Prozess
// lesen könnte. `main` endet in `os.Exit`; der Re-Exec-Harness startet
// darum dasselbe Binary mit anderen Argumenten und wertet Exit-Code,
// stdout und stderr des Kindprozesses — ohne Naht im Produktionscode
// (`ADR-0082` §Kontext (4a), §Entscheidung 2).
package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// reexecMarker schaltet denselben Binary-Stand zwischen seinen zwei
// Rollen um: gesetzt ruft `TestMain` `main` auf, sonst fährt es die
// Testfälle.
const reexecMarker = "PGC_FEED_MAIN_REEXEC"

// TestMain trägt beide Rollen: als Kindprozess führt dieses Binary `main`
// mit den übergebenen Argumenten aus und endet über dessen `os.Exit`
// (bzw. über den eigenen Ausgang 0, wo `main` zurückkehrt); als
// Elternprozess fährt es die Testfälle.
func TestMain(m *testing.M) {
	if os.Getenv(reexecMarker) == "1" {
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// prozessLauf trägt, was der Kindprozess hinterlässt.
type prozessLauf struct {
	ausgang int
	stdout  string
	stderr  string
}

// nichtErreichbareQuelle trägt die Adresse des Prozess-Starts, die ohne
// Dienst bis zum Verbindungsaufbau kommt und dort scheitert: der lokale
// Port 1 ist geschlossen, `connect_timeout` begrenzt den Versuch am
// Adapter (`postgresstorage.New`, `ADR-0082` §Kontext (3)). Der
// Datenbankname trägt die Rolle des DSN (`ADR-0047`) und macht die drei
// Vorbedingungen in der Fehlerzeile unterscheidbar.
const nichtErreichbareQuelle = "postgres://x:x@127.0.0.1:1/%s?sslmode=disable&connect_timeout=1"

// starteVollstaendig trägt die Umgebungs-Vorbedingungen des Prozess-
// Starts (`ConfigFromEnv`): die drei rollen-spezifischen DSNs
// (`ADR-0047`) und Quelle, Publication, Slot, Tabellen. Die Adressen sind
// nicht erreichbar — der Lauf kommt bis zur Verdrahtung, nicht weiter.
func starteVollstaendig() []string {
	return []string{
		"CDC_CAPTURE_DSN=" + fmt.Sprintf(nichtErreichbareQuelle, "capture-db"),
		"CDC_ADMIN_DSN=" + fmt.Sprintf(nichtErreichbareQuelle, "admin-db"),
		"CDC_READER_DSN=" + fmt.Sprintf(nichtErreichbareQuelle, "reader-db"),
		"CDC_SOURCE_ID=src-1",
		"CDC_PUBLICATION=pub_1",
		"CDC_SLOT=slot_1",
		"CDC_TABLES=public.t1=tbl-1:sv-1",
	}
}

// ohneVorbedingung entfernt genau eine der Vorbedingungen aus dem
// vollständigen Satz — die Eingabeseite, an die die Fehlerzeile des
// Prozesses gebunden wird.
func ohneVorbedingung(vollstaendig []string, name string) []string {
	gefiltert := make([]string, 0, len(vollstaendig))
	for _, eintrag := range vollstaendig {
		if !strings.HasPrefix(eintrag, name+"=") {
			gefiltert = append(gefiltert, eintrag)
		}
	}
	return gefiltert
}

// kindUmgebung setzt die Umgebung des Kindprozesses selbst zusammen: die
// Vorbedingungen stellen die Testfälle, `PATH` wird aus der
// Eltern-Umgebung übernommen — daran löst ein nicht absoluter
// `os.Args[0]` auf. Jede weitere Variable der Eltern-Umgebung reist
// nicht mit. `GOCOVERDIR` wird durchgereicht, sobald es gesetzt ist:
// der Kindprozess ist dasselbe instrumentierte Binary, und `go test`
// mergt die Zähler-Dateien aller Prozesse, die es tragen, in das Profil
// des Gate-Laufs; ohne die Weitergabe bleibt der Kindprozess ungezählt
// (`ADR-0082` §Kontext (4a)).
func kindUmgebung(vorbedingungen ...string) []string {
	umgebung := []string{"PATH=" + os.Getenv("PATH"), reexecMarker + "=1"}
	if verzeichnis := os.Getenv("GOCOVERDIR"); verzeichnis != "" {
		umgebung = append(umgebung, "GOCOVERDIR="+verzeichnis)
	}
	return append(umgebung, vorbedingungen...)
}

// fahreProzess startet das Binary erneut und wartet auf sein **Ende** —
// gewertet wird der Ausgang des abgeschlossenen Prozesses, nicht eine
// verstrichene Dauer.
func fahreProzess(t *testing.T, vorbedingungen []string, args ...string) prozessLauf {
	t.Helper()

	kind := exec.Command(os.Args[0], args...)
	kind.Env = kindUmgebung(vorbedingungen...)
	var stdout, stderr bytes.Buffer
	kind.Stdout = &stdout
	kind.Stderr = &stderr

	err := kind.Run()
	ausgang := 0
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("Prozess %v ließ sich nicht starten: %v", args, err)
		}
		ausgang = exitErr.ExitCode()
	}
	return prozessLauf{ausgang: ausgang, stdout: stdout.String(), stderr: stderr.String()}
}

// TestVersionMeldetDenLieferstandOhneDiagnose trägt den Ausgang 0 des
// `--version`-Zweigs: der Prozess schreibt den Lieferstand nach stdout
// und endet ohne Diagnose-Zeile — die Form `pg-change-feed: …` tragen
// alle Fehlerpfade von `main`.
func TestVersionMeldetDenLieferstandOhneDiagnose(t *testing.T) {
	lauf := fahreProzess(t, nil, "--version")

	if lauf.ausgang != 0 {
		t.Errorf("--version: Ausgang = %d, wollen 0; stderr = %q", lauf.ausgang, lauf.stderr)
	}
	if want := "pg-change-feed " + version + "\n"; lauf.stdout != want {
		t.Errorf("--version: stdout = %q, wollen %q", lauf.stdout, want)
	}
	if strings.Contains(lauf.stderr, "pg-change-feed:") {
		t.Errorf("--version: stderr = %q, wollen keine Diagnose-Zeile", lauf.stderr)
	}
}

// TestVorbedingungFehltTraegtJeModusIhrenAusgang bindet die zwei
// Ausgänge, die derselbe Fehler je Modus trägt: die Sondermodi enden bei
// fehlender Umgebungs-Vorbedingung mit Ausgang 1, der argumentlose
// CDC-Lauf mit Ausgang 2. Die Ausgabe wird an die **Eingabeseite**
// geprüft — die Fehlerzeile nennt die in diesem Fall fehlende Variable,
// nicht irgendeine.
func TestVorbedingungFehltTraegtJeModusIhrenAusgang(t *testing.T) {
	for _, fall := range []struct {
		modus   string
		args    []string
		fehlt   string
		ausgang int
	}{
		{"--healthcheck", []string{"--healthcheck"}, "CDC_READER_DSN", 1},
		{"register-consumer", []string{"register-consumer", "c-1"}, "CDC_SLOT", 1},
		{"acknowledge-consumer", []string{"acknowledge-consumer", "c-1", "42"}, "CDC_TABLES", 1},
		{"diagnose", []string{"diagnose"}, "CDC_CAPTURE_DSN", 1},
		{"argumentloser Lauf", nil, "CDC_PUBLICATION", 2},
	} {
		t.Run(fall.modus, func(t *testing.T) {
			vorbedingungen := ohneVorbedingung(starteVollstaendig(), fall.fehlt)
			lauf := fahreProzess(t, vorbedingungen, fall.args...)

			if lauf.ausgang != fall.ausgang {
				t.Errorf("%s ohne %s: Ausgang = %d, wollen %d; stderr = %q",
					fall.modus, fall.fehlt, lauf.ausgang, fall.ausgang, lauf.stderr)
			}
			if !strings.Contains(lauf.stderr, fall.fehlt) {
				t.Errorf("%s ohne %s: stderr = %q, wollen die fehlende Variable %q in der Fehlerzeile",
					fall.modus, fall.fehlt, lauf.stderr, fall.fehlt)
			}
		})
	}
}

// TestSondermodiMitVollstaendigerUmgebungNennenIhreRolle trägt den
// Zweig, der die vier Sondermodi durchlässt statt sie abzuweisen: mit
// vollständiger Umgebung und nicht erreichbarer Instanz. Der **Aufruf**
// ist netzlos — nur der Rumpf der Modi braucht eine erreichbare Instanz.
//
// Gebunden an die Eingabeseite: jeder Modus erreicht seine **eigene**
// Rolle (`ADR-0047`) und die Zeile trägt deren Datenbanknamen — die zwei
// Lese-Modi den aus `CDC_READER_DSN`, die zwei Verwaltungs-Modi den aus
// `CDC_ADMIN_DSN`; der Name der je anderen Rolle und der aus
// `CDC_CAPTURE_DSN` darf darin nicht vorkommen. Die Zeile wird zusätzlich
// auf die moduseigene Präfix-Form geprüft, nicht auf den blossen Modus-Namen.
func TestSondermodiMitVollstaendigerUmgebungNennenIhreRolle(t *testing.T) {
	for _, fall := range []struct {
		modus string
		args  []string
		rolle string
		fremd string
		nennt string
	}{
		{"--healthcheck", []string{"--healthcheck"}, "reader-db", "admin-db", "healthcheck: Instanz nicht erreichbar"},
		{"register-consumer", []string{"register-consumer", "c-1"}, "admin-db", "reader-db", "register-consumer: Fehlerklasse storage"},
		{"acknowledge-consumer", []string{"acknowledge-consumer", "c-1", "42"}, "admin-db", "reader-db", "acknowledge-consumer: Fehlerklasse storage"},
		{"diagnose", []string{"diagnose"}, "reader-db", "admin-db", "diagnose: Instanz nicht erreichbar"},
	} {
		t.Run(fall.modus, func(t *testing.T) {
			lauf := fahreProzess(t, starteVollstaendig(), fall.args...)

			if lauf.ausgang != 1 {
				t.Errorf("%s: Ausgang = %d, wollen 1 (nicht erreichbare Instanz); stderr = %q",
					fall.modus, lauf.ausgang, lauf.stderr)
			}
			if !strings.Contains(lauf.stderr, fall.nennt) {
				t.Errorf("%s: stderr = %q, wollen die moduseigene Zeile %q", fall.modus, lauf.stderr, fall.nennt)
			}
			if !strings.Contains(lauf.stderr, fall.rolle) {
				t.Errorf("%s: stderr = %q, wollen den Datenbanknamen der eigenen Rolle (%q)",
					fall.modus, lauf.stderr, fall.rolle)
			}
			if strings.Contains(lauf.stderr, fall.fremd) || strings.Contains(lauf.stderr, "capture-db") {
				t.Errorf("%s: stderr = %q, wollen keine fremde Rolle — erwartet nur %q",
					fall.modus, lauf.stderr, fall.rolle)
			}
		})
	}
}

// TestArgumentFehlerEndenMitAusgang2 trägt die Argument-Fehler des
// Dispatches: fehlende, überzählige und nicht deutbare Argumente enden
// vor jedem Dienst-Zugriff mit Ausgang 2. Die Ausgabe-Hälfte ist an die
// **Eingabeseite** gebunden: jede Zeile wird auf ihre eigene, den
// konkreten Verstoß nennende Formulierung geprüft — der blosse
// Modus-Name genügt nicht als Prüfung, weil ihn auch die generische
// Zeile für ein unbekanntes Argument trägt (`main.go:104`), und der
// nicht deutbare Wert allein genügt nicht, weil ihn auch die Ursache des
// Deutungsfehlers trägt.
func TestArgumentFehlerEndenMitAusgang2(t *testing.T) {
	for _, fall := range []struct {
		name  string
		args  []string
		nennt string
	}{
		{"register-consumer ohne Namen", []string{"register-consumer"}, "register-consumer erwartet genau einen Namen als Argument"},
		{"register-consumer mit zwei Namen", []string{"register-consumer", "a", "b"}, "register-consumer erwartet genau einen Namen als Argument"},
		{"acknowledge-consumer ohne Position", []string{"acknowledge-consumer", "c-1"}, "acknowledge-consumer erwartet Consumer-Kennung und Position als Argumente"},
		{"acknowledge-consumer mit überzähligem Argument", []string{"acknowledge-consumer", "c-1", "42", "x"}, "acknowledge-consumer erwartet Consumer-Kennung und Position als Argumente"},
		{"acknowledge-consumer mit nicht deutbarer Position", []string{"acknowledge-consumer", "c-1", "42x"}, `"42x" ist kein gültiger Offset`},
		{"acknowledge-consumer mit leerer Position", []string{"acknowledge-consumer", "c-1", ""}, `"" ist kein gültiger Offset`},
		{"unbekanntes Argument", []string{"--unbekannt"}, "unbekanntes Argument"},
	} {
		t.Run(fall.name, func(t *testing.T) {
			lauf := fahreProzess(t, starteVollstaendig(), fall.args...)

			if lauf.ausgang != 2 {
				t.Errorf("%s: Ausgang = %d, wollen 2; stderr = %q", fall.name, lauf.ausgang, lauf.stderr)
			}
			if !strings.Contains(lauf.stderr, fall.nennt) {
				t.Errorf("%s: stderr = %q, wollen %q in der Fehlerzeile", fall.name, lauf.stderr, fall.nennt)
			}
		})
	}
}

// TestLaufMitNichtErreichbarerQuelleEndetMitAusgang1 trägt den
// argumentlosen Lauf bis in die Verdrahtung: die Vorbedingungen sind
// vollständig, `bootstrap.Run` scheitert am Verbindungsaufbau, `main`
// reicht den Fehler als Diagnose-Zeile weiter und endet mit Ausgang 1
// (ein Adapter-Fehler beendet den Prozess-Aufrufer mit Ausgang 1).
// Gebunden ist hier, **welche Rolle** der Lauf erreicht: die Zeile trägt
// den Datenbanknamen aus `CDC_CAPTURE_DSN` und nicht den aus
// `CDC_ADMIN_DSN` — der Prozess bleibt vor dem Aktivierungs-Pool stehen.
// Dass er schon am **ersten** Konstruktor endet, bindet der
// Schwester-Test `run_test.go` über den Sentinel.
func TestLaufMitNichtErreichbarerQuelleEndetMitAusgang1(t *testing.T) {
	lauf := fahreProzess(t, starteVollstaendig())

	if lauf.ausgang != 1 {
		t.Errorf("Lauf ohne erreichbare Quelle: Ausgang = %d, wollen 1; stderr = %q", lauf.ausgang, lauf.stderr)
	}
	if !strings.HasPrefix(lauf.stderr, "pg-change-feed: ") {
		t.Errorf("Lauf ohne erreichbare Quelle: stderr = %q, wollen eine Diagnose-Zeile", lauf.stderr)
	}
	if !strings.Contains(lauf.stderr, "capture-db") {
		t.Errorf("Lauf ohne erreichbare Quelle: stderr = %q, wollen den Datenbanknamen aus CDC_CAPTURE_DSN", lauf.stderr)
	}
	if strings.Contains(lauf.stderr, "admin-db") {
		t.Errorf("Lauf ohne erreichbare Quelle: stderr = %q, wollen keinen Zugriff über den Aktivierungs-Pool", lauf.stderr)
	}
}
