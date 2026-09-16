package bootstrap_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// Diese Datei bindet die Zusagen der Rollen-Verdrahtung an das **reale
// Artefakt**, das sie trägt: `tools/schema/nacharbeit-roles.sql`
// (`BEO-PGC/rollen-test-abdeckungsluecken`, Punkt 1 — 2×).
//
// Die Schwester-Tests in `roles_wiring_test.go` prüfen dieselben Zusagen
// real gegen die Testcontainer-Instanz (`make test-store`) — und
// **entziehen und erteilen die Rechte dabei selbst** per `REVOKE`/`GRANT`.
// Damit messen sie ihren eigenen Aufbau: eine Regression **in der
// Rollout-Datei** (etwa ein zurückgenommener `SELECT`-Grant auf
// `cdc.process_heartbeat`, genau der von `ADR-0048` korrigierte
// Ursprungstext) bliebe dort unsichtbar, und das `Cleanup` des Tests
// räumt die Spur mit weg. Dieser Test liest stattdessen die Datei selbst
// und ist netzlos — er läuft in `make gates` mit.
//
// Er prüft nicht, ob die Rechte auf einer laufenden Instanz wirken (das
// bleibt Sache der `make test-store`-Tests und des Rollen-Rollouts,
// `ADR-0043`); er prüft, dass der **Grant-Text, den der Rollout ausrollt**,
// die Rechte trägt, die die Verdrahtung je Rolle voraussetzt
// (`ADR-0047`, `ADR-0048`, `ADR-0014`, `LH-QA-SEC-001`…`003`).

// rolloutDateiPfad löst `tools/schema/nacharbeit-roles.sql` über die Lage
// dieser Testdatei auf — unabhängig vom Arbeitsverzeichnis des Testlaufs.
func rolloutDateiPfad(t *testing.T) string {
	t.Helper()
	_, quelle, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Rollen-Rollout-Test: Quelldatei nicht auflösbar")
	}
	return filepath.Join(filepath.Dir(quelle), "..", "..", "tools", "schema", "nacharbeit-roles.sql")
}

// grantZeile trägt ein einzelnes `GRANT <Rechte> ON <Objekte> TO <Rollen>;`
// der Rollout-Datei. Zeilen innerhalb der `DO $$ … $$;`-Blöcke beginnen
// nicht mit `GRANT` und matchen darum nicht — die `CREATE ROLE`-Form wird
// separat geprüft.
var grantZeile = regexp.MustCompile(
	`(?im)^[ \t]*GRANT[ \t]+([A-Za-z]+(?:[ \t]*,[ \t]*[A-Za-z]+)*)[ \t]+ON[ \t]+([^;]+?)[ \t]+TO[ \t]+([A-Za-z0-9_]+(?:[ \t]*,[ \t]*[A-Za-z0-9_]+)*)[ \t]*;`)

// createRoleZeile trägt ein `CREATE ROLE <Name> <Attribute…>;` des
// Rollout-DO-Blocks.
var createRoleZeile = regexp.MustCompile(`(?im)^[ \t]*CREATE ROLE[ \t]+([A-Za-z0-9_]+)[ \t]+([^;]+);`)

// rolloutRechte liest den Grant-Text der Rollout-Datei und liefert
// `Rolle → Objekt → Rechte` — nur die tatsächlichen `GRANT`-Anweisungen,
// Kommentarzeilen (`--`) sind vorher entfernt, weil sie Beispiel-Grants in
// Prosa tragen (`GRANT cdc_reader TO ihr_login;`).
func rolloutRechte(t *testing.T, pfad string) (map[string]map[string]map[string]bool, string) {
	t.Helper()
	roh, err := os.ReadFile(pfad)
	if err != nil {
		t.Fatalf("Rollout-Datei %s nicht lesbar: %v", pfad, err)
	}
	inhalt := string(roh)

	// Kommentare entfernen: der Grant-Text ist sonst nicht von den ihn
	// erklärenden Beispiel-Zeilen unterscheidbar.
	var ohneKommentare []string
	for _, zeile := range strings.Split(inhalt, "\n") {
		if index := strings.Index(zeile, "--"); index >= 0 {
			zeile = zeile[:index]
		}
		ohneKommentare = append(ohneKommentare, zeile)
	}
	quelltext := strings.Join(ohneKommentare, "\n")

	rechte := map[string]map[string]map[string]bool{}
	for _, treffer := range grantZeile.FindAllStringSubmatch(quelltext, -1) {
		privilegien := zerlegeListe(treffer[1])
		objekte := zerlegeListe(treffer[2])
		rollen := zerlegeListe(treffer[3])
		for _, rolle := range rollen {
			if rechte[rolle] == nil {
				rechte[rolle] = map[string]map[string]bool{}
			}
			for _, objekt := range objekte {
				if rechte[rolle][objekt] == nil {
					rechte[rolle][objekt] = map[string]bool{}
				}
				for _, privileg := range privilegien {
					rechte[rolle][objekt][privileg] = true
				}
			}
		}
	}
	return rechte, quelltext
}

// zerlegeListe teilt eine kommaseparierte Liste und normalisiert jeden
// Eintrag auf Kleinschreibung — `GRANT SELECT, INSERT ON X TO Y` und
// `GRANT select,insert on X to Y` tragen dieselbe Aussage.
func zerlegeListe(roh string) []string {
	var ergebnis []string
	for _, eintrag := range strings.Split(roh, ",") {
		normalisiert := strings.ToLower(strings.Join(strings.Fields(eintrag), " "))
		if normalisiert != "" {
			ergebnis = append(ergebnis, normalisiert)
		}
	}
	sort.Strings(ergebnis)
	return ergebnis
}

// rechteVon liefert die Rechte einer Rolle auf einem Objekt (leer, wenn
// nicht vorhanden) — das Objekt wird auf die Form der Datei normalisiert
// (`cdc.process_heartbeat`).
func rechteVon(rechte map[string]map[string]map[string]bool, rolle, objekt string) map[string]bool {
	if objekte, ok := rechte[rolle]; ok {
		if privilegien, ok := objekte[strings.ToLower(objekt)]; ok {
			return privilegien
		}
	}
	return map[string]bool{}
}

// TestRolloutDateiTraegtDieRechteDerVerdrahtung ist der Kern dieser Datei:
// je geprüfter Rolle die Rechte, die die Verdrahtung auf dem Objekt
// voraussetzt, gegen den **Rollout-Text** gehalten. Die Rolle-Zuordnung
// selbst steht in `ADR-0047`s Tabelle (Verdrahtung über
// `CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/`CDC_READER_DSN`); dieser Test prüft
// die andere Hälfte derselben Zusage — dass der ausgerollte Grant sie
// trägt.
//
// Rot färbende Mutationen (real gefahren, Exit-Code im Closure-Bericht):
// `SELECT` aus dem `cdc.process_heartbeat`-Grant für `cdc_admin` entfernen
// (der von `ADR-0048` korrigierte Ursprungstext) · den `DELETE`-Grant auf
// `cdc.transaction`/`cdc.change` für `cdc_admin` entfernen · `DELETE` an
// `cdc_capture` auf `cdc.change` ergänzen.
func TestRolloutDateiTraegtDieRechteDerVerdrahtung(t *testing.T) {
	rechte, _ := rolloutRechte(t, rolloutDateiPfad(t))

	// (1) Heartbeat-Pfad — `runHeartbeat` → `postgresstorage.NewHeartbeat`
	// über `cfg.AdminDSN`. `UpsertHeartbeat` schreibt
	// `INSERT … ON CONFLICT (source_id) DO UPDATE …`; PostgreSQL verlangt
	// dafür zusätzlich `SELECT` auf der Zieltabelle (ADR-0048). `cdc_admin`
	// ist die Rolle dieses Pfads (ADR-0047).
	for _, privileg := range []string{"select", "insert", "update"} {
		if !rechteVon(rechte, "cdc_admin", "cdc.process_heartbeat")[privileg] {
			t.Fatalf("cdc_admin fehlt %s auf cdc.process_heartbeat im Rollout-Text — der Heartbeat-Schreibpfad (INSERT … ON CONFLICT DO UPDATE, ADR-0048) scheitert real mit SQLSTATE 42501", strings.ToUpper(privileg))
		}
	}

	// (2) Retention-Pfad — `RunRetentionUseCase` → Store-Adapter über
	// `cfg.AdminDSN`; die Löschausführung braucht `DELETE` auf beiden
	// Tabellen (ADR-0014).
	for _, objekt := range []string{"cdc.transaction", "cdc.change"} {
		if !rechteVon(rechte, "cdc_admin", objekt)["delete"] {
			t.Fatalf("cdc_admin fehlt delete auf %s im Rollout-Text — die Retention-Löschausführung (ADR-0014) scheitert real", objekt)
		}
	}

	// (3) Erfassungspfad — `StorePort.PersistTransaction` über
	// `cfg.CaptureDSN` schreibt `INSERT` auf beide Tabellen; `DELETE` und
	// `UPDATE` bleiben der Admin-Rolle vorbehalten (LH-QA-SEC-002: getrennte
	// Berechtigbarkeit).
	for _, objekt := range []string{"cdc.transaction", "cdc.change"} {
		privilegien := rechteVon(rechte, "cdc_capture", objekt)
		if !privilegien["insert"] {
			t.Fatalf("cdc_capture fehlt insert auf %s im Rollout-Text — der Erfassungspfad scheitert real", objekt)
		}
		for _, verboten := range []string{"delete", "update"} {
			if privilegien[verboten] {
				t.Fatalf("cdc_capture trägt %s auf %s im Rollout-Text — der Rollenschnitt (LH-QA-SEC-002) ist aufgeweicht", strings.ToUpper(verboten), objekt)
			}
		}
	}

	// (4) Kein Erfassungs-Grant auf der Heartbeat-Tabelle: sie ist die
	// Schreibfläche des Admin-Pfads; ein zusätzlicher Grant an `cdc_capture`
	// verbreiterte die Least-Privilege-Fläche (LH-QA-SEC-001) ohne Aufrufer.
	if privilegien := rechteVon(rechte, "cdc_capture", "cdc.process_heartbeat"); len(privilegien) != 0 {
		t.Fatalf("cdc_capture trägt Rechte auf cdc.process_heartbeat im Rollout-Text: %v — die Tabelle gehört zum Admin-Pfad", privilegien)
	}

	// (5) Lesepfad — `cdc_reader` liest ausschließlich über die Views
	// (`LH-FA-SST-002`), kein Grant auf einer Basistabelle (LH-QA-SEC-003:
	// die Definer-Semantik der Views trägt den Lesezugriff).
	for _, basis := range []string{
		"cdc.change", "cdc.transaction", "cdc.source_table", "cdc.schema_version",
		"cdc.consumer", "cdc.consumer_position", "cdc.process_heartbeat", "cdc.table_schema",
	} {
		if privilegien := rechteVon(rechte, "cdc_reader", basis); len(privilegien) != 0 {
			t.Fatalf("cdc_reader trägt Rechte auf der Basistabelle %s: %v — der Lesezugriff läuft über die Views", basis, privilegien)
		}
	}
	leseViews := []string{"cdc.active_tables", "cdc.consumer_status", "cdc.changes"}
	for _, view := range leseViews {
		if !rechteVon(rechte, "cdc_reader", view)["select"] {
			t.Fatalf("cdc_reader fehlt select auf der Lese-View %s im Rollout-Text — der Lesezugriffsweg (LH-FA-SST-002) scheitert real", view)
		}
	}
}

// TestRolloutDateiLegtGruppenrollenOhneAnmeldungAn trägt die Grenze des
// Rollenschnitts (`nacharbeit-roles.sql` §Grenze, `ADR-0047`
// Kontext-Befund 1): die drei Rollen sind Gruppenrollen — sie bündeln
// Privilegien und tragen **kein** `LOGIN`. Die anmeldefähige Identität legt
// der Betreiber je Rolle selbst an (`GRANT cdc_reader TO ihr_login;`,
// `docs/user/benutzerhandbuch.md`). Eine `LOGIN`-Rolle im Rollout-Text
// wäre eine Anmelde-Identität ohne Passwort-Verwaltung — die Datei würde
// damit eine Betriebsentscheidung treffen, die nicht ihre ist.
//
// `cdc_capture` trägt zusätzlich das `REPLICATION`-Attribut: der
// Replication-Stream verlangt es (`ADR-0008`); PostgreSQL vererbt Attribute
// nicht über Mitgliedschaft, der Betreiber setzt es darum zusätzlich auf
// die Login-Identität (ADR-0047 Kontext-Befund 2).
//
// Rot färbende Mutation: `CREATE ROLE cdc_admin NOLOGIN;` auf `LOGIN`
// umstellen — dann bricht der Test.
func TestRolloutDateiLegtGruppenrollenOhneAnmeldungAn(t *testing.T) {
	_, quelltext := rolloutRechte(t, rolloutDateiPfad(t))

	attribute := map[string]string{}
	for _, treffer := range createRoleZeile.FindAllStringSubmatch(quelltext, -1) {
		attribute[strings.ToLower(treffer[1])] = strings.ToLower(strings.Join(strings.Fields(treffer[2]), " "))
	}

	for _, rolle := range []string{"cdc_capture", "cdc_admin", "cdc_reader"} {
		attrs, ok := attribute[rolle]
		if !ok {
			t.Fatalf("Rollout-Text legt die Gruppenrolle %s nicht an — die Verdrahtung verbindet sich mit ihr (ADR-0047)", rolle)
		}
		if !strings.Contains(attrs, "nologin") {
			t.Fatalf("Gruppenrolle %s trägt kein NOLOGIN (%q) — der Rollenschnitt sieht keine Anmelde-Identität im Rollout vor", rolle, attrs)
		}
	}
	if !strings.Contains(attribute["cdc_capture"], "replication") {
		t.Fatalf("cdc_capture trägt kein REPLICATION-Attribut (%q) — der Replication-Stream des Erfassungspfads verlangt es (ADR-0008)", attribute["cdc_capture"])
	}
}
