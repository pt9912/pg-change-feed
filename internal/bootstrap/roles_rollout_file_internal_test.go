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
//
// **Benannte Grenzen dieses Tests — was er strukturell nicht sieht.** Der
// Parser liest literale `GRANT … ;`-Anweisungen der Datei (Zeilenanfang,
// Kommentare vorher entfernt). Drei Stellen bleiben damit außerhalb seiner
// Reichweite, und keine davon ist still:
//
//   - **Der dynamische Grant im `DO $$ … $$;`-Block** (`nacharbeit-roles.sql`
//     §„CREATE PUBLICATION/ALTER PUBLICATION"): `EXECUTE format('GRANT CREATE
//     ON DATABASE %I TO cdc_admin', current_database())` ist eine zur
//     Laufzeit zusammengesetzte Zeichenkette, keine Anweisung der Datei; der
//     Datenbankname steht erst im Rollout fest. Der Parser trifft sie
//     bewusst nicht — ein Mustertreffer auf `GRANT` innerhalb eines
//     `format(...)`-Strings wäre kein Beleg für einen ausgerollten Grant.
//     Der Grant trägt `CREATE PUBLICATION`/`ALTER PUBLICATION` (ADR-0047);
//     sein realer Beleg ist der DB-gestützte Tier (`roles_test.go`).
//   - **Ein Lesegrant auf einem Objekt, das keine schreibende Rolle hält**
//     (z. B. ein zusätzliches `GRANT SELECT ON cdc.<fremd> TO cdc_reader`):
//     die Regel unten fängt solche Grants nur, wenn eine schreibende Rolle
//     dasselbe Objekt trägt. „Ist dieses Objekt eine View oder eine
//     Basistabelle?" ist hier nicht entscheidbar — die Klassifikation steht
//     im `views:`-Knoten von `tools/schema/schema.yaml`, und die liegt
//     **nicht** im Build-Kontext der `coverage`-Stufe. Die DDL in
//     `internal/adapters/driven/postgresstorage/schema.sql` liegt dort, trägt
//     aber nur die Store-Seite (fünf Tabellen, keine View, kein
//     `cdc.metrics`) und taugt darum nicht als Klassifikator.
//   - **Grants, die eine andere Rollout-Datei ausrollt** (etwa `cdc.metrics`
//     in `nacharbeit-observability.sql`): Gegenstand dieses Tests ist genau
//     eine Datei.
//
// Eine benannte Lücke ist zulässig, eine stille nicht — die drei stehen
// deshalb hier und nicht in der Stille.

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

// istKlassenGrant meldet, ob eine Grant-Objektklausel nicht **ein** Objekt
// benennt, sondern eine ganze Objektklasse (`ALL TABLES IN SCHEMA cdc`,
// `ALL SEQUENCES IN SCHEMA …`). Der Schema-USAGE-Grant trägt die Form
// `SCHEMA cdc` und ist damit **keine** Klasse — er ist die Vorbedingung
// jedes Objektzugriffs und wird von der Regel unten nicht getroffen.
func istKlassenGrant(objekt string) bool {
	return strings.HasPrefix(objekt, "all ") || strings.Contains(" "+objekt+" ", " all ") ||
		strings.Contains(objekt, " in schema ")
}

// TestRolloutDateiTraegtDieRechteDerVerdrahtung ist der Kern dieser Datei:
// je geprüfter Rolle die Rechte, die die Verdrahtung auf dem Objekt
// voraussetzt, gegen den **Rollout-Text** gehalten. Die Rolle-Zuordnung
// selbst steht in `ADR-0047`s Tabelle (Verdrahtung über
// `CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/`CDC_READER_DSN`); dieser Test prüft
// die andere Hälfte derselben Zusage — dass der ausgerollte Grant sie
// trägt.
//
// Drei Prüfungen sind **Regeln über den geparsten Grant-Bestand**, keine
// Namenslisten: (6) kein Rollen-Grant über eine ganze Objektklasse,
// (6a) auf einem Schema-Objekt nur `USAGE` und (7) der Leser teilt kein
// Objekt mit einer schreibenden Rolle. Eine Namensliste könnte das nicht
// tragen: sie müsste die Objekte kennen, die es geben **darf**, und ließe
// jedes nicht aufgeführte Objekt durch — auch das, welches eine schreibende
// Rolle hält. Die Regeln leiten den Umfang aus dem Rollen-Bild der Datei
// selbst ab und ziehen damit mit jedem Objekt mit; geprüft wird, was ein
// Grant **tut**, nicht wie er heißt.
//
// Rot färbende Mutationen (real gefahren; die Exit-Codes führt §7 der
// Closure-Notiz von `slice-093`): `SELECT` aus dem
// `cdc.process_heartbeat`-Grant für `cdc_admin` entfernen (der von
// `ADR-0048` korrigierte Ursprungstext) · den `DELETE`-Grant auf
// `cdc.transaction`/`cdc.change` für `cdc_admin` entfernen · `DELETE` an
// `cdc_capture` auf `cdc.change` ergänzen · den Schema-USAGE-Grant
// entfernen · die Lese-View `cdc.retention_blockers` bzw. `cdc.backfill_status`
// aus dem Reader-Grant
// streichen · `GRANT ALL ON ALL TABLES IN SCHEMA cdc TO cdc_reader`
// anhängen · `GRANT CREATE ON SCHEMA cdc TO cdc_reader` anhängen. Zu (4a):
// `INSERT` an `cdc_capture` auf `cdc.backfill_run` anhängen · `UPDATE` an
// `cdc_admin` anhängen · den `UPDATE`-Grant von `cdc_capture` streichen ·
// `SELECT` an `cdc_reader` anhängen. Zu (4b): den Grant
// `SELECT, UPDATE ON cdc.administration_request TO cdc_admin` streichen bzw.
// auf `SELECT` kürzen · `INSERT` oder `DELETE` an `cdc_admin` anhängen ·
// `SELECT` an `cdc_capture` oder `cdc_reader` anhängen.
func TestRolloutDateiTraegtDieRechteDerVerdrahtung(t *testing.T) {
	rechte, _ := rolloutRechte(t, rolloutDateiPfad(t))

	// (0) Die Vorbedingung jedes anderen Grants der Datei: ohne `USAGE` auf
	// dem Schema kann **keine** der drei Rollen ein Objekt darin berühren —
	// ein entzogener USAGE-Grant färbt alle Pfade der Verdrahtung rot, nicht
	// nur einen. Der Grant trägt die Form `SCHEMA cdc`.
	for _, rolle := range []string{"cdc_capture", "cdc_admin", "cdc_reader"} {
		if !rechteVon(rechte, rolle, "schema cdc")["usage"] {
			t.Fatalf("%s fehlt USAGE auf dem Schema cdc im Rollout-Text — ohne ihn ist jeder andere Grant dieser Datei wirkungslos (LH-QA-SEC-001)", rolle)
		}
	}

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

	// (4a) Backfill-Run-Zustand (`SPEC-029`, `ADR-0113` Festlegung 1): die
	// Annahme (`cdc_admin`) prüft und legt die Zeile an — `SELECT`,
	// `INSERT`; der Worker (`cdc_capture`) führt sie fort — `SELECT`,
	// `UPDATE`. Niemand trägt `DELETE`, `cdc_admin` kein `UPDATE`,
	// `cdc_capture` kein `INSERT`, `cdc_reader` kein Recht auf die
	// Basistabelle (er liest über die View der Sichtbarkeit).
	for _, erwartet := range []struct {
		rolle, privileg string
	}{
		{"cdc_admin", "select"}, {"cdc_admin", "insert"},
		{"cdc_capture", "select"}, {"cdc_capture", "update"},
	} {
		if !rechteVon(rechte, erwartet.rolle, "cdc.backfill_run")[erwartet.privileg] {
			t.Fatalf("%s fehlt %s auf cdc.backfill_run im Rollout-Text — die Annahme bzw. der Worker des Backfills scheitert real mit SQLSTATE 42501 (ADR-0113)", erwartet.rolle, strings.ToUpper(erwartet.privileg))
		}
	}
	for _, verboten := range []struct {
		rolle, privileg string
	}{
		{"cdc_admin", "update"}, {"cdc_admin", "delete"},
		{"cdc_capture", "insert"}, {"cdc_capture", "delete"},
		{"cdc_reader", "select"}, {"cdc_reader", "insert"}, {"cdc_reader", "update"}, {"cdc_reader", "delete"},
	} {
		if rechteVon(rechte, verboten.rolle, "cdc.backfill_run")[verboten.privileg] {
			t.Fatalf("%s trägt %s auf cdc.backfill_run im Rollout-Text — der Rollenschnitt des Backfills (ADR-0113 Festlegung 1, LH-QA-SEC-001…003) ist aufgeweicht", verboten.rolle, strings.ToUpper(verboten.privileg))
		}
	}

	// (4b) Antrags-Queue (`ADR-0050`, `LH-FA-ADM-001`): die
	// Administrations-Goroutine und die Annahme eines Backfills laufen über
	// `CDC_ADMIN_DSN` und lesen offene Anträge (`ListPending`), leiten den
	// Spaltenausschluss-Stand ab (`ExcludedColumns`) und vermerken den
	// Ausgang (`MarkApplied`/`MarkFailed`) — `SELECT`, `UPDATE`. Angelegt
	// werden Anträge nur von den SECURITY-DEFINER-Funktionen
	// (`nacharbeit-administration.sql`): `cdc_admin` trägt weder `INSERT`
	// noch `DELETE`, die beiden anderen Rollen kein Recht auf die Tabelle.
	for _, privileg := range []string{"select", "update"} {
		if !rechteVon(rechte, "cdc_admin", "cdc.administration_request")[privileg] {
			t.Fatalf("cdc_admin fehlt %s auf cdc.administration_request im Rollout-Text — die Antragsverarbeitung (ListPending, Vermerk applied/failed, ADR-0050) und die Annahme des Backfills scheitern unter einem Least-Privilege-Login real mit SQLSTATE 42501", strings.ToUpper(privileg))
		}
	}
	for _, verboten := range []struct {
		rolle, privileg string
	}{
		{"cdc_admin", "insert"}, {"cdc_admin", "delete"},
		{"cdc_capture", "select"}, {"cdc_capture", "insert"}, {"cdc_capture", "update"}, {"cdc_capture", "delete"},
		{"cdc_reader", "select"}, {"cdc_reader", "insert"}, {"cdc_reader", "update"}, {"cdc_reader", "delete"},
	} {
		if rechteVon(rechte, verboten.rolle, "cdc.administration_request")[verboten.privileg] {
			t.Fatalf("%s trägt %s auf cdc.administration_request im Rollout-Text — der Rollenschnitt der Antrags-Queue (ADR-0050, LH-QA-SEC-001…003) ist aufgeweicht", verboten.rolle, strings.ToUpper(verboten.privileg))
		}
	}

	// (5) Lesepfad — `cdc_reader` liest über die **fünf** Lese-Views der
	// Datei (`LH-FA-SST-002`, mit `retention_blockers` aus
	// `LH-FA-RET-005` und `backfill_status` aus `SPEC-029`). Die fünf sind der
	// erklärte Lese-Umfang des Rollouts;
	// fehlt eine, scheitert der jeweilige Lesezugriffsweg real.
	for _, view := range []string{
		"cdc.active_tables", "cdc.consumer_status", "cdc.changes", "cdc.retention_blockers", "cdc.backfill_status",
	} {
		if !rechteVon(rechte, "cdc_reader", view)["select"] {
			t.Fatalf("cdc_reader fehlt select auf der Lese-View %s im Rollout-Text — der Lesezugriffsweg (LH-FA-SST-002) scheitert real", view)
		}
	}

	// (6) Keine Rolle bekommt Rechte über eine ganze Objektklasse: eine
	// `ALL TABLES IN SCHEMA`-Zeile ist keine Aufzählung, sondern eine
	// Rundum-Vergabe über jedes zu diesem Zeitpunkt existierende Objekt des
	// Schemas — auch über jedes Schreibpfad-Objekt. Sie lässt die übrigen
	// Zeilen der Datei stehen; aufgehoben ist die **Trennung**, die (7)
	// prüft, weil der Umfang des Lesers dann nicht mehr an seinen fünf
	// Views hängt (LH-QA-SEC-001).
	for rolle, objekte := range rechte {
		for objekt := range objekte {
			if istKlassenGrant(objekt) {
				t.Fatalf("%s trägt im Rollout-Text einen Grant über eine ganze Objektklasse (%q) — der Rollenschnitt (LH-QA-SEC-001) ist damit aufgehoben", rolle, objekt)
			}
		}
	}

	// (6a) Auf einem Schema-Objekt trägt **keine** Rolle mehr als `USAGE`:
	// das Schema-Objekt ist die Vorbedingung jedes Objektzugriffs, kein
	// Objekt selbst — `CREATE` oder `ALL` darauf ist die Rundum-Vergabe über
	// das ganze Schema, nicht ein Zugriff auf eines seiner Objekte.
	for rolle, objekte := range rechte {
		for objekt, privilegien := range objekte {
			if !strings.HasPrefix(objekt, "schema ") {
				continue
			}
			for privileg := range privilegien {
				if privileg != "usage" {
					t.Fatalf("%s trägt %s auf dem Schema-Objekt %s im Rollout-Text — über `USAGE` hinaus ist das eine Rundum-Vergabe über das ganze Schema (LH-QA-SEC-001)", rolle, strings.ToUpper(privileg), objekt)
				}
			}
		}
	}

	// (7) Der Leser teilt kein Objekt mit einer schreibenden Rolle: was eine
	// der beiden schreibenden Rollen hält, ist ein Schreibpfad-Objekt und
	// gehört nicht in den Leseumfang (LH-QA-SEC-003; die Definer-Semantik
	// der Views trägt den Lesezugriff, nicht ein Grant auf der Basistabelle).
	// Ausgenommen ist ein Schema-Objekt: (6a) hält jedes von ihnen auf
	// `USAGE` fest, und die Vorbedingung beider Seiten ist kein Objektzugriff.
	for _, schreibend := range []string{"cdc_capture", "cdc_admin"} {
		for objekt := range rechte["cdc_reader"] {
			if strings.HasPrefix(objekt, "schema ") {
				continue
			}
			if _, geteilt := rechte[schreibend][objekt]; geteilt {
				t.Fatalf("cdc_reader und %s tragen beide Rechte auf %s — der Leseumfang (LH-QA-SEC-003) ist nicht mehr vom Schreibpfad getrennt", schreibend, objekt)
			}
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
