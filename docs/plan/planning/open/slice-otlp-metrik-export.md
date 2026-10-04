# Slice otlp-metrik-export: periodischer OTLP/HTTP-Push der Betriebskennzahlen (Teil a: Konfiguration, Lese-Port, Sender, Takt)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice (gemessen: `ls docs/plan/planning/*.md` am Planungsstand nennt nur
`README.md`; die Roadmap führt unter *Offene Wellen* keine Welle), siehe
Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht
(Modul 6).

**Abhängigkeit:** startet erst, wenn
[`slice-api-token-mehrfach-konfiguration`](../done/slice-api-token-mehrfach-konfiguration.md)
in `done/` liegt: dieser Slice hebt die Zugangsdaten-Klasse der
Konfigurationsdatei von neun auf elf Schlüssel und setzt auf der Zahl auf, die der
Vorgänger gesetzt hat (Schnitt und Begründung im Plan des Vorgängers, §1); die
Nummern der neuen Meldungscodes folgen der Zuweisungsregel dort (§3).
[`slice-tls-http-grpc-server`](slice-tls-http-grpc-server.md) ist **keine**
fachliche Voraussetzung; die empfohlene Reihenfolge A → B → C hält die gemeinsam
berührten Dateien (`wiring.go`, `config_file.go`, Code-Tabelle, Handbuch) ohne
Konflikt, und das WIP-Limit 1 serialisiert ohnehin. **Messstand:** die
Suchlauf-Zeilen in §3 sind am Planungsstand `a53f4e75…` gemessen; die Slices davor
bewegen mehrere dieser Eigenschaften (Zugangsdaten-Klasse, Code-Tabelle,
Handbuch) — der Implementer misst beim Start neu (§4).

**Bezug:** [`LH-FA-SST-010`](../../../../spec/lastenheft.md) (Scope),
[`LH-FA-SST-004`](../../../../spec/lastenheft.md) und
[`LH-QA-OPS-003`](../../../../spec/lastenheft.md) (die SQL-Sicht und ihre
Kennzahlen bleiben unverändert),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (WAL-Rückstand),
[`LH-QA-SEC-001`](../../../../spec/lastenheft.md) (Leserolle ohne zusätzliche
Rechte),
[`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) (Festlegung 1 bis 9,
Folgepflicht),
[`ADR-0152`](../../adr/0152-zugangsdaten-klasse-elf-schluessel.md) (die zwei
OTLP-Schlüssel der Klasse, Test auf elf),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
(Warn-Bereich 6, Konfigurationscodes),
[`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md) (Telemetrie
hinter Ports),
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Messgegenstand des Coverage-Gates),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).

**Berührte Spec-Stellen:** [`SPEC-033`](../../../../spec/pflichtenheft.md) ·
[`SPEC-016`](../../../../spec/pflichtenheft.md) (Datei-Schlüssel `otlp_interval`
zulässig; `otlp_endpoint` und `otlp_headers` in der Klasse) ·
[`SPEC-009`](../../../../spec/pflichtenheft.md) (Namen der Kennzahlen) ·
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Warn-Bereich 6, Klasse
`configuration`) ·
[`ARC-002`](../../../../spec/architecture.md) ·
[`ARC-004`](../../../../spec/architecture.md) ·
[`ARC-006`](../../../../spec/architecture.md) ·
[`ARC-007`](../../../../spec/architecture.md) ·
[`ARC-011`](../../../../spec/architecture.md) (OTLP-Empfänger, Sequenz
„Metrik-Export“).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** Planner-Agent, direkt beauftragt (kein Architect: die Entscheidung
liegt mit [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) vor, die
Schichtung folgt `ADR-0024` und der Sequenz der Architektur-Sicht; scheitert der
Bibliothek-Prüfpunkt an beiden Wegen, ist das eine Architect-Frage, §4).
**Datum:** 2026-10-04.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass:** Die Betriebskennzahlen sind heute nur per SQL-Sicht `cdc.metrics`
lesbar (Pull); ein Betreiber mit einem OpenTelemetry-Empfänger übersetzt selbst
([`LH-FA-SST-010`](../../../../spec/lastenheft.md)). Gemessen am Parent: die Sicht
führt 9 Kennzahlen (`grep -o "SELECT 'cdc_[a-z_]*'"
tools/schema/nacharbeit-observability.sql | sort -u`, 9 Zeilen), den zehnten,
`cdc_wal_retention_bytes`, misst nur der Prozess (`runWALRetentionCheck` in
`internal/bootstrap/wiring.go`); kein OTLP-Bezug im Go-Baum (`go.mod` nennt
weder `go.opentelemetry.io` noch ein OTLP-Modul; §3 Suchlauf Zeile 3 zählt 23
Trefferzeilen für „otlp“/„opentelemetry“, alle in Spec und Doku).

**Ziel:** Ist `CDC_OTLP_ENDPOINT` gesetzt, überträgt der Feed-Container die zehn
Kennzahlen wiederkehrend als OTLP/HTTP-Protobuf-Gauges (`POST
<Basis-URL>/v1/metrics`) — gelesen aus `cdc.metrics` über die Leserolle plus dem
zuletzt im Prozess gemessenen `cdc_wal_retention_bytes` —, in eigener Goroutine,
nie blockierend, mit Frist, ohne Warteschlange und mit gedrosselter Warnung;
ohne Endpunkt bleibt alles unverändert. Der Slice belegt das gegen einen
`httptest`-Empfänger und gegen die reale Datenbank (`make test`,
`make test-store`); der Beleg am laufenden Container ist der Folge-Slice.

**Schnitt der Zugangsdaten-Klasse.** Dieser Slice trägt `otlp_endpoint` und
`otlp_headers` in `forbiddenFileCredentialKeys` ein und hebt Test und Handbuch
von **neun auf elf** — genau seine zwei Schlüssel; die zwei anderen trägt
[`slice-api-token-mehrfach-konfiguration`](../done/slice-api-token-mehrfach-konfiguration.md)
(Begründung dort, §1). `otlp_interval` ist ein **zulässiges** Datei-Feld und kein
Klassenmitglied (eine Zahl trägt keine Zugangsdaten).

**Übergabe-Block an `otlp-metrik-export-e2e`** (Folge-Slice, **noch keine Datei**;
Adresse ist die Closure dieses Slice, in der der Planner ihn anlegt — §7). Der
Gegenstand, den der Realserver-Beleg dort tragen muss und den dieser Slice
**nicht** trägt: (1) ein Wegwerf-Empfänger im Compose-Netz nimmt `POST /v1/metrics`
des laufenden Feed-Containers entgegen, und ein übertragener Wert ist gleich dem
Wert der Sicht im selben Moment (zeitabhängige Kennzahlen bis auf den Abstand der
Messzeitpunkte; `cdc_wal_retention_bytes` aus dem Prozess) — die Fitness-Function-Zeile
„Integrationstest (Compose, Empfänger-Wegwerf-Client)“ von
[`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md); (2) Empfänger
nicht erreichbar oder antwortet `500`: der Container läuft, Erfassung und
Bestätigung unverändert, die Übertragung setzt beim nächsten Takt wieder ein;
(3) ein ungültiger Endpunkt, Header oder Takt beendet den Start mit der Klasse
`configuration` und dem Meldungscode in der Log-Zeile; (4) ein `https`-Endpunkt
gegen den Vertrauensspeicher des Runtime-Images (nach Kenntnisstand enthält das
distroless-Image CA-Zertifikate, ungeprüft); (5) eine `abdeckung_declare`-Zeile für
[`LH-FA-SST-010`](../../../../spec/lastenheft.md), `docs/user/e2e-abdeckung.md`
vom Runner geschrieben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Realserver-Beleg am Compose-Stack** — `otlp-metrik-export-e2e`: er braucht
  einen Wegwerf-Empfänger und eine Runner-Phase und würde diesen Slice über die
  drei Liefer-Punkte treiben; die Tier-Trennung hält dieser Slice auf
  Unit- und Verdrahtungs-Ebene (der Übergabe-Block oben trägt den Gegenstand).
- **Ein Pull-Endpunkt (Prometheus-Scrape), Traces, Logs, Dashboards,
  Alarmierung** — Out-of-Scope von
  [`LH-FA-SST-010`](../../../../spec/lastenheft.md),
  [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) Festlegung 8.
- **Eine Änderung der SQL-Sicht `cdc.metrics`, neue Rechte oder neue Tabellen** —
  [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) Festlegung 2:
  der Export liest, was die Sicht liefert, unter der Rolle `cdc_reader`; der
  WAL-Rückstand bleibt außerhalb der Sicht
  ([`LH-QA-SEC-001`](../../../../spec/lastenheft.md)).
- **Eine eigene CA, ein Client-Zertifikat oder Proxy-Einstellungen für den
  Empfänger** — die Spec nennt nur Schema `http`/`https` und Host; `https` nutzt
  den Vertrauensspeicher des Prozesses. Ein Betreiber-Bedarf wäre eine neue
  Anforderung (Re-Evaluierungs-Trigger von
  [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) liest nur
  Wiederholung und Typ).
- **TLS des eigenen Servers und der Token-Wechsel** —
  [`slice-tls-http-grpc-server`](slice-tls-http-grpc-server.md) und
  [`slice-api-token-mehrfach-konfiguration`](../done/slice-api-token-mehrfach-konfiguration.md);
  andere Anforderungen, andere Eingriffsstellen.
- **Die Änderung des WAL-Prüfzugs** (`runWALRetentionCheck`): sein
  Schwellenverhalten und seine Abbruchpfade bleiben; der Slice liest nur den
  zuletzt gemessenen Wert aus einem gemeinsamen Halter (die bestehenden Tests des
  Zugs bleiben ohne Änderung grün).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Alle Beleg-Angaben dieser Liste sind **Zusagen** („zu belegen durch …“): der
Planungsstand hat keinen davon gefahren
([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
Instanz B). Die gedruckte Zeile bzw. der Exit-Code steht im Bericht des
Implementers und in §7.

- [ ] **Konfiguration, Meldungscodes, Zugangsdaten-Klasse elf (Liefer-Punkt 1,
      [`LH-FA-SST-010`](../../../../spec/lastenheft.md) Boundary und Negative).**
      (a) `CDC_OTLP_ENDPOINT` und `CDC_OTLP_HEADERS` sind env-exklusiv,
      `CDC_OTLP_INTERVAL_SECONDS` schlägt das Datei-Feld `otlp_interval` (Default
      60), beide Zugriffswege (`ConfigFromEnv`, `mergeConfig`) rufen **eine**
      Validierungsfunktion. Gültig: Endpunkt mit Schema `http` oder `https` und
      Host; Intervall eine ganze Zahl von 5 bis 3600; Header `k=v,k2=v2` mit je nicht
      leerem Schlüssel und einem `=`. Ungültig endet mit `ErrConfiguration` und
      einem **neuen** Meldungscode (nächste freie Nummer ab `PCF-E2008` zum
      Zeitpunkt der Arbeit, Zuweisungsregel im Plan von
      [`slice-api-token-mehrfach-konfiguration`](../done/slice-api-token-mehrfach-konfiguration.md)
      §3): Endpunkt mit Schema `ftp`, ohne Host, leer-aber-gesetzt-Fälle der
      Umgebung wie „ungesetzt“; Intervall 4, 3601, `abc`, `60.5`; Header ohne `=`,
      mit leerem Schlüssel. Jeder Negativfall prüft den Code **und** den Klartext
      und trägt eine Gegenprobe ohne die Verletzung (Grenzwerte 5 und 3600 gültig;
      [`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`](../observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/observation.md)).
      **Festlegungen dieses Plans, wo die Spec schweigt** (Spec-Lücke-Kandidaten,
      §6): der Header wird am **ersten** `=` geteilt (der Wert darf `=` tragen,
      etwa Base64-Auffüllung; ein Komma im Wert gibt es nicht, es trennt); Leerraum
      im Schlüssel ist ein Fehler, Leerraum im Wert bleibt erhalten (`Bearer abc`);
      `CDC_OTLP_HEADERS` **ohne** Endpunkt endet mit `ErrConfiguration` (Vorbild
      `validateNatsStreamTokenRequiresURL`: ein Wert, der ohne Endpunkt nie wirkt,
      ist keine stille Auslassung); ein Intervall ohne Endpunkt ist zulässig (die
      Datei darf es tragen).
      (b) **Warn-Bereich 6.** Die Menge der Warn-Bereiche wächst von 1 bis 5 auf
      1 bis 6 an **allen** Trägern der Menge (Suchlauf Zeile 1): der Kommentar in
      `internal/domain/messagecode/codes.go` (heute „… 5 Konfiguration und Start,
      9 reserviert“) und die Konstanten `WarnExportFailed` (`PCF-W6001`) und
      `WarnMetricsReadFailed` (`PCF-W6002`) (Namen Arbeitsnamen) samt Tabelle;
      `Area` in `internal/domain/messagecode/messagecode.go` (die Grenze
      `c[5] > '5'`) und ihr Kommentar; `messagecode_test.go`
      (`TestWarningCodesCarryTheirArea`: Kommentar, Meldung „außerhalb 1 bis 5“, die
      benannte Menge je Bereich um Bereich 6 — sonst färbt die Gleichheit
      „Tabelle gleich benannte Konstanten“ rot); der Satz zu den Bereichen im
      Handbuch (Abschnitt „Meldungscodes“, heute: „ein Bereich 5 … ist vorgesehen
      und trägt keine Warnung“) mit **Bereich 6 Beobachtbarkeit und Transport**;
      die Katalog-Zeilen beider Warn-Codes und der neuen `configuration`-Codes im
      Handbuch; `make meldungscodes-check` Exit 0. Die Info-Zeile bei
      Wiederaufnahme trägt keinen Code (die Schwere `I` ist reserviert,
      [`SPEC-008`](../../../../spec/pflichtenheft.md)).
      (c) **Zugangsdaten-Klasse elf:** `forbiddenFileCredentialKeys` trägt
      `otlp_endpoint` und `otlp_headers` zusätzlich zu den neun des Vorgängers;
      `TestConfigFromFileLehntZugangsdatenAb` iteriert über die elf (Liste **und**
      Kommentar „neun“); der Mengengleichheit-Test des Vorgängers bleibt grün; eine
      Datei mit `otlp_interval: 30` lädt ohne Fehler, eine mit `otlp_endpoint:`
      endet mit `ErrConfiguration` und der Zeile „Zugangsdaten bleiben
      env-var-exklusiv“; der Klassen-Satz im Handbuch §5 nennt dieselben **elf**
      (von Hand nachgezählt,
      [`ADR-0152`](../../adr/0152-zugangsdaten-klasse-elf-schluessel.md) Fitness
      Function).
- [ ] **Export-Pipeline (Liefer-Punkt 2, [`LH-FA-SST-010`](../../../../spec/lastenheft.md)
      Happy Path, Boundary).**
      (a) **Bibliothek-Prüfpunkt als erster Schritt**
      ([`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) Festlegung 1:
      Bibliothek nicht bindend, *nach Kenntnisstand, ungeprüft*). Gemessen und im
      Bericht genannt: (i) Existiert `go.opentelemetry.io/proto/otlp` als Modul, in
      welcher Version, und trägt es die Nachricht `ExportMetricsServiceRequest`
      samt Gauge; (ii) die transitiven Anforderungen gegen `go.mod`
      (`google.golang.org/grpc v1.84.0`, `google.golang.org/protobuf v1.36.12` am
      Parent): hebt die Aufnahme vorhandene Abhängigkeiten an, welche; (iii) der
      Fußabdruck: neue Zeilen in `go.sum` und die Größe des Produktions-Binaries
      vor und nach (Toolchain-Container, gemessene Bytes); (iv) die Aufnahme in
      `go.mod`/`go.sum` läuft über einen Toolchain-Container-Lauf, nicht über ein
      Host-`go` ([`AGENTS.md`](../../../../AGENTS.md) §3.1; es gibt kein
      `make`-Ziel dafür, gemessen: `git grep -n 'go get\|tidy' -- Makefile
      harness/mk` ohne Treffer); `make mod-download` bereitet den netzlosen
      `make test` vor. **Entscheidungsregel:** die Typen-Bibliothek mit `net/http`,
      wenn (i) bis (iii) keinen Befund tragen, den der Bericht als Grund gegen sie
      nennt; sonst der Rückfall (das offizielle SDK) mit demselben Messsatz; tragen
      beide einen Befund, hält der Slice an (§4, Architect). Der Wire-Vertrag
      ([`SPEC-033`](../../../../spec/pflichtenheft.md) Drahtform) ist unabhängig von
      der Wahl und gilt als Test-Soll.
      (b) **Lese-Port und Adapter:** ein Outbound-Port für die Zeilen der Sicht
      (Name, Label, Wert) und sein Adapter im Paket `postgresstorage` (Muster
      `DiagnosticsPort`; **kein** neues DB-Paket, damit die namentlichen
      Paketlisten von `tools/harness/db-package-lists-check.sh` unberührt
      bleiben — zu belegen durch `make coverage-gate` Exit 0 und den Lauf des
      Prüfers). Test in `make test-store` gegen die reale Datenbank unter einem
      `cdc_reader`-Login (keine weitere Berechtigung nötig): die Ausgabe des
      Adapters ist gleich einem `SELECT metric_name, label, value FROM cdc.metrics`
      derselben Datenbank mit Consumern und Fehlerzuständen als Zeilen mit Label
      (zeitabhängige Kennzahlen mit benannter Toleranz).
      (c) **Export-Port und OTLP-Adapter** (Driven Adapter, `ADR-0024`): die zehn
      Kennzahlen als `Gauge` — auch die mit dem Namensteil `_total` —, Namen,
      Einheiten und Datenpunkt-Attribute (`consumer`, `class`) genau wie die Tabelle
      in [`SPEC-033`](../../../../spec/pflichtenheft.md), Resource-Attribute
      `service.name` = `pg-change-feed` und `cdc.source_id`; `POST
      <Basis-URL>/v1/metrics`, `Content-Type: application/x-protobuf`, die Header
      der Konfiguration; `2xx` ist angenommen, jeder andere Status, jeder
      Verbindungs- und jeder Zeitfehler ein Fehlschlag; die Basis-URL mit oder ohne
      abschließenden Schrägstrich ergibt genau **einen** `/v1/metrics`. Tests gegen
      `httptest` (`make test`, netzlos, Loopback): der empfangene Körper wird mit
      der Protobuf-Bibliothek zurückgelesen und trägt alle zehn Namen, den Typ
      Gauge, Einheit und Attribut je Zeile der Tabelle; ohne gemessenen
      WAL-Rückstand entfällt genau dieser Datenpunkt (neun Namen); Status `200`,
      `202` angenommen, `404`, `500`, Verbindung abgewiesen und überschrittene Frist
      Fehlschlag.
      (d) **Use Case und Takt:** ein Export-Use-Case in der Application (liest,
      ergänzt den zuletzt gemessenen WAL-Rückstand aus einem gemeinsamen Halter, den
      der unveränderte WAL-Prüfzug schreibt, überträgt) und die Verdrahtung in
      `internal/bootstrap` mit **eigener Goroutine**, Frist je Versuch
      `min(Takt, 10 s)`, ohne Warteschlange. Zu belegen durch Tests: ein
      Empfänger, der nicht antwortet, hält Erfassung und Heartbeat nicht auf und der
      Versuch endet an der Frist; der nächste Takt wiederholt, ein ausgefallener Takt
      wird nicht nachgeholt (nach Wiederaufnahme genau ein Request je Takt); der
      erste Fehlschlag seit Start oder Wiederaufnahme erzeugt **eine** Warn-Zeile
      mit dem Code `PCF-W6001` (Übertragung) bzw. `PCF-W6002` (Lesen der Sicht),
      weitere höchstens alle 5 Minuten (Fake-Uhr, `outbound.Clock`, ohne Wartezeit),
      die Wiederaufnahme eine Info-Zeile; Health-Zustand und `error_class` des
      Heartbeats bleiben unberührt; **weder ein Header-Wert noch die
      Zugangsdaten-Teile einer URL erscheinen im Log** (ein Test mit einem Header
      `Authorization=geheim` und einer URL mit Benutzer-Passwort-Teil prüft die
      Log-Ausgabe auf deren Abwesenheit). Verdrahtungs-Test in `internal/bootstrap`
      (`make test-store`, reale Datenbank, `httptest`-Empfänger): der Empfänger
      erhält einen Request, dessen `cdc_changes_processed` gleich dem Wert der Sicht
      ist; ohne Endpunkt bleibt die Zahl der Requests über ein Beobachtungsfenster
      **null** und kein Sender wird konstruiert (Boundary „keine Verbindung“).
      (e) **Mutationsprobe (Zusage, an einer Kopie im Scratchpad, nicht an der
      Datei im Arbeitsbaum — [`AGENTS.md`](../../../../AGENTS.md) §3.1; Instanz: der
      Go-Test des jeweiligen Pakets im Toolchain-Container mit der Kopie als
      Mount):** fünf Stellen einzeln mutiert, je die rote Farbe im Bericht — `Gauge`
      zu `Sum` (der Typ-Fall), die Frist je Versuch entfernt (der Blockier-Fall),
      die 5-Minuten-Drosselung entfernt (der Drossel-Fall), die untere Intervallgrenze
      5 auf 4 (der Grenzfall), und der WAL-Datenpunkt ohne gemessenen Wert
      trotzdem ausgegeben (der Fall „noch keiner gemessen“). Die Verallgemeinerung
      auf weitere Stellen ist **hergeleitet**, nicht erprobt; der Bericht nennt je
      Mutation Stelle und Instanz.
- [ ] **Benutzerhandbuch (Liefer-Punkt 3).** `docs/user/benutzerhandbuch.md`:
      Umgebungsvariablen-Tabelle (§5) um `CDC_OTLP_ENDPOINT`,
      `CDC_OTLP_HEADERS`, `CDC_OTLP_INTERVAL_SECONDS` (Gültigkeit, Default, Fehlerfall
      mit Klasse `configuration`, ohne Endpunkt kein Export und keine Verbindung);
      das Datei-Feld `otlp_interval` samt YAML-Beispiel und der Satz, dass Endpunkt und
      Header env-exklusiv sind (Zugangsdaten); der Klassen-Satz mit den **elf**
      Schlüsseln; ein Abschnitt zum Export unter „Metriken lesen“ mit der
      Kennzahlen-Tabelle (Name, Einheit, Attribut), dem **Hinweis, dass alle
      Kennzahlen Momentstände (Gauge) sind, auch die mit dem Namensteil `_total`**
      ([`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md)
      Folgepflicht), Takt, Frist, „kein Nachholen“, den Warnungen mit ihren Codes,
      der Gleichheit mit der SQL-Sicht je Messzeitpunkt und dem Verhalten bei
      Ausfall des Empfängers; die Katalog-Zeilen der neuen Codes und der Satz zu
      Bereich 6; Änderungshistorie-Zeile mit der nächsten freien Nummer
      ([`BEO-PGC/handbuch-versionshistorie-uebersprungen`](../observations/BEO-PGC/handbuch-versionshistorie-uebersprungen/observation.md)).
      Kennungsfrei: `make handbuch-public-doc-check` Exit 0; die neuen
      Meldungstexte im Go-Quelltext tragen keine Kennung
      (`make ausgabe-kennungen-check` Exit 0).

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9; darin `make coverage-gate` für die
      neue, netzlos prüfbare Fläche, `make a-check`, `make meldungscodes-check`),
      `make docs-check` Exit 0 (Kennungen in diesem Plan verlinkt), `make test`,
      `make test-store` und `make test-integration` Exit 0 (der Container ohne
      Endpunkt läuft unverändert — Boundary am echten Prozess), `make image`
      Exit 0 (der Zug ändert `go.mod` und Build-Kontext-Dateien).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors
      (Zeile `make test-store`: ein Satz zum Lese-Adapter und zum
      Verdrahtungs-Test, `ADR-0149` Folgepflicht); gemeldete Träger fremder
      Dateien mit der Closure nachgezogen (§3 Suchlauf,
      [`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Der Folge-Slice `otlp-metrik-export-e2e` ist als Datei in `open/` angelegt
      (Planner, mit der Closure) und trägt den Übergabe-Block aus §1 als
      committeten Text in seinem §2 — kein Verweis auf diesen Plan allein.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt (gemessen: `docs/plan/planning/` trägt keine flache Welle-Datei) und „die nächste Welle-Closure“ damit keine Adresse ist.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `go.mod`, `go.sum` | update | die OTLP-Abhängigkeit nach dem Prüfpunkt (§2 Liefer-Punkt 2 (a)); Aufnahme über einen Toolchain-Container-Lauf |
| `internal/application/port/outbound/` (Arbeitsnamen `metricsread.go`, `metricexport.go`) | neu | Lese-Port der Sicht und Export-Port; der Doc-Kommentar am Port nennt die Pflichten des Adapters (Frist unter vom Abbruch gelöstem Kontext, Fehlerklasse), der Adapter-Teil des DoD trägt dieselben Pflichten ([`BEO-PGC/adapter-pflicht-eines-ports-ohne-traeger-im-adapter-slice`](../observations/BEO-PGC/adapter-pflicht-eines-ports-ohne-traeger-im-adapter-slice/observation.md)) |
| `internal/application/usecase/exportmetrics/` (Arbeitsname) | neu | Export-Use-Case: lesen, den WAL-Rückstand ergänzen, übertragen; Warn-Zustand mit Fake-Uhr testbar |
| `internal/adapters/driven/postgresstorage/` (`metrics.go`, Anfrage in `queries/queries.go`) | neu / update | Lese-Adapter der Sicht neben `PostgresDiagnosticsAdapter`; Test in `make test-store` unter einem `cdc_reader`-Login |
| `internal/adapters/driven/otlpexport/` (Arbeitsname) | neu | OTLP/HTTP-Adapter (Protobuf, `net/http`), Abbildung der zehn Gauges; Tests gegen `httptest` |
| `internal/bootstrap/wiring.go` | update | `Config`-Felder, Konstanten, `ConfigFromEnv`, Halter des zuletzt gemessenen WAL-Rückstands (vom unveränderten Prüfzug geschrieben), Start und Ende der Export-Goroutine |
| `internal/bootstrap/config_file.go` | update | `fileConfig` (`otlp_interval`), `mergeConfig`, `forbiddenFileCredentialKeys` (9 → 11) samt Kommentar |
| `internal/bootstrap/config_file_internal_test.go`, `internal/bootstrap/wiring_test.go` und neue Verdrahtungs-Tests | update / neu | Zugangsdaten-Klasse auf elf, Parser-Fälle beider Zugriffswege, Verdrahtung gegen reale Datenbank und `httptest` |
| `internal/domain/messagecode/codes.go`, `messagecode.go`, `messagecode_test.go` | update | Warn-Bereich 6: Konstanten und Tabelle, `Area`, Kommentare, benannte Menge im Test; neue `configuration`-Codes |
| `docs/user/benutzerhandbuch.md` | update | §2 Liefer-Punkt 3 |
| `harness/README.md` | update | Zeile `make test-store` |

- **Zuweisung der Codes.** `PCF-W6001` (Übertragung fehlgeschlagen) und
  `PCF-W6002` (Lesen der Sicht fehlgeschlagen) sind durch
  [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) Festlegung 7
  benannt; die Konstanten-Namen legt der Implementer fest. Die Konfigurationscodes
  folgen der Zuweisungsregel im Plan von
  [`slice-api-token-mehrfach-konfiguration`](../done/slice-api-token-mehrfach-konfiguration.md)
  §3 (nächste freie Nummer ab `PCF-E2008`, am Arbeitsstand gemessen). Die Zahl der
  Codes ist offen: ein Code für „Endpunkt ungültig“, einer für „Header ungültig“,
  einer für „Intervall außerhalb“ ist ein Vorschlag, kein Soll.
- **Warn-Zustand (Festlegung dieses Plans, wo die Spec schweigt).** Ein
  gemeinsamer Fehlerzustand: der erste Fehlschlag eines Zyklus seit Start oder
  seit der letzten Wiederaufnahme — beim Lesen der Sicht oder beim Übertragen —
  warnt mit dem Code seiner Stufe; weitere Fehlschläge warnen höchstens alle
  5 Minuten, wieder mit dem Code der Stufe des jeweiligen Fehlschlags;
  Wiederaufnahme ist ein vollständig erfolgreicher Zyklus und erzeugt die
  Info-Zeile.
- **Erster Versuch.** Der erste Versuch erfolgt nach einem vollen Takt, nicht beim
  Start (Festlegung dieses Plans: der WAL-Messwert liegt erst nach dem ersten
  Prüf-Takt des Zugs vor, `heartbeatInterval` 5 s, die untere Intervallgrenze ist
  ebenfalls 5 s). Die Spec nennt den Zeitpunkt des ersten Versuchs nicht.
- **Zeitstempel.** Jeder Datenpunkt trägt den Messzeitpunkt seines Zyklus
  (Beginn des Lesens der Sicht, bzw. der Zeitpunkt der WAL-Messung für
  `cdc_wal_retention_bytes`); die Gleichheit mit der Sicht gilt je Messzeitpunkt
  ([`SPEC-033`](../../../../spec/pflichtenheft.md), „Quelle der Werte“).
- **Reihenfolge:** (1) Startmessung: Suchlauf unten neu am Arbeitsstand,
  `grep -o "SELECT 'cdc_[a-z_]*'" tools/schema/nacharbeit-observability.sql |
  sort -u` (die 9 der ADR); (2) Bibliothek-Prüfpunkt und Entscheidung, `go.mod`;
  (3) Ports, Adapter (Lese- und OTLP-Adapter) mit ihren Tests;
  (4) Use Case, Verdrahtung, Konfiguration, Codes, Warn-Bereich 6,
  Zugangsdaten-Klasse; (5) Mutationsläufe an Kopien; (6) `make test-store`,
  `make test-integration`, `make image`; (7) Handbuch, README, `make gates`.
- **Suchlauf (§3.13 der Regeln, [`AGENTS.md`](../../../../AGENTS.md)).** Bewegte
  Eigenschaften: die **Menge der Warn-Bereiche** (Zählwort „1 bis 5“, „Bereich 5“,
  `c[5] > '5'`), die **Zugangsdaten-Klasse** (Symbole `otlp_endpoint`,
  `otlp_headers`, `otlp_interval`, `CDC_OTLP_`), der **OTLP-Bezug im Baum**
  (Hedge „otlp“, „opentelemetry“), die **Kennzahl** `cdc_wal_retention_bytes` und
  die **Aufzählungen der Kennzahlen** (`cdc_oldest_change_age_seconds`,
  `cdc_consumer_lag`). Parent ist der Stand `a53f4e75…` (`git rev-parse HEAD`
  am Planungsstand, vor dem Plan-Commit; nie `HEAD`). Suchraum: ganzer Baum ohne
  `docs/reviews`, `done/`, `observations/`, `.harness/baseline`, die offenen
  Pläne und `docs/plan/adr` (Accepted ADRs halten ihren Stand,
  [`AGENTS.md`](../../../../AGENTS.md) §3.5). **Parent-Zeilen gemessen; die
  `diff`-Zeilen tragen den Planungsstand.** Beim Start misst der Implementer neu
  (die Vorgänger-Slices ändern Handbuch, Code-Tabelle und Konfiguration, Zeilen und
  Zahlen verschieben sich): neue Commit-Kennung als Parent, die Abweichung zum
  Plan mit Ursache im Bericht.
  **Erwartung am Endstand** (hergeleitet, nicht gemessen): Zeile 1 hat keinen
  Träger mehr, der nur „1 bis 5“ nennt (`Area`, Kommentare, Test, Handbuch tragen
  die 6); Zeile 2 wächst um Code, Test, Handbuch; Zeile 3 um Adapter und Doku; Zeile 4
  um Adapter und Test; Zeile 5 um Adapter, Abbildung und Handbuch. Jede
  Trefferzeile wird gelesen, nicht gezählt.

```suchlauf
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 7 -n -E "Bereich 1 bis 5|1 bis 5 in der ersten|Bereich 5 für Konfiguration|außerhalb 1 bis 5|5 Konfiguration und Start|c\[5\] > '5'" -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 10 -n -E 'otlp_endpoint|otlp_headers|CDC_OTLP_|otlp_interval' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 23 -n -i -E 'otlp|opentelemetry' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 27 -n -E 'cdc_wal_retention_bytes' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 26 -n -E 'cdc_oldest_change_age_seconds|cdc_consumer_lag' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 7 -n -E "Bereich 1 bis 5|1 bis 5 in der ersten|Bereich 5 für Konfiguration|außerhalb 1 bis 5|5 Konfiguration und Start|c\[5\] > '5'" -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 10 -n -E 'otlp_endpoint|otlp_headers|CDC_OTLP_|otlp_interval' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 23 -n -i -E 'otlp|opentelemetry' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 27 -n -E 'cdc_wal_retention_bytes' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 26 -n -E 'cdc_oldest_change_age_seconds|cdc_consumer_lag' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
```

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1), **und**
[`slice-api-token-mehrfach-konfiguration`](../done/slice-api-token-mehrfach-konfiguration.md)
liegt in `done/`, **und** der Implementer hat die Startmessung gefahren: die
Suchlauf-Zeilen aus §3 am Arbeitsstand neu gemessen (neue Commit-Kennung als
Parent, `make suchlauf-nachmessen PLAN=` mit diesem Plan Exit 0 nach Anpassung
von Parent und Soll; jede Abweichung zum Planungsstand mit Ursache im
Bericht — der Vorgänger bewegt die Zugangsdaten-Klasse und das Handbuch), **und**
der Bibliothek-Prüfpunkt (§2 Liefer-Punkt 2 (a)) hat eine Entscheidung
getragen, bevor Adapter-Code entsteht.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice wächst über
  die drei Punkte aus §2 — Kandidat für die Teilung: Konfiguration, Codes und
  Zugangsdaten-Klasse (Liefer-Punkt 1) als eigener Slice vor der Pipeline
  (Liefer-Punkt 2); der Planner teilt dann die Datei.
- `in-progress` → `open` (blockiert — Carveout?): der Bibliothek-Prüfpunkt trägt
  an **beiden** Wegen (Typen-Bibliothek und offizielles SDK) einen Befund, den
  [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) nicht
  vorsieht — dann ist es eine Architect-Frage (Festlegung 1 nennt nur diese zwei
  Wege); oder die Spec 0.15.0 ändert den Wortlaut von
  [`SPEC-033`](../../../../spec/pflichtenheft.md). Kein Carveout, weil dann kein
  Gate rot ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, `make test-store`, `make test-integration` und `make image` enden mit
Exit 0, `make suchlauf-nachmessen PLAN=` mit diesem Plan endet am Endstand mit
Exit 0 (die `diff`-Zeilen auf den Endstand gesetzt, jede Trefferzeile gelesen),
der Folge-Slice `otlp-metrik-export-e2e` liegt als Datei in `open/`, und die
Closure-Notiz in §7 trägt einen Lerneintrag (geschärfte Regel, neuer Sensor oder
benannte Spec-Lücke; Kandidaten: der Befund des Bibliothek-Prüfpunkts als
Erprobung der „ungeprüften“ Annahme von
[`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md), die vier
Festlegungen dieses Plans, wo die Spec schweigt — Header-Form, erster Versuch,
gemeinsamer Fehlerzustand, Header ohne Endpunkt — als benannte Spec-Lücke).
Ein Gate, das am Stand der Closure rot ist, geht nur mit dokumentiertem
Carveout nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen
(Grund) · weiter offen (BEO-Eintrag im Register). Alle Ausgänge sind bis zur
Closure **offen** (Platzhalter `Ausgang: offen bis Closure`).

- **Die Bibliothek ist ungeprüft.** Modul-Pfad, Version, transitive
  Anforderungen und Fußabdruck von `go.opentelemetry.io/proto/otlp` sind *nach
  Kenntnisstand* ([`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md)
  Festlegung 1); eine Aufnahme kann `google.golang.org/grpc` oder
  `google.golang.org/protobuf` anheben und damit SDK-Tiers, Generator-Pins und den
  `generated-sync`-Lauf berühren. Gegenmaßnahme: der Prüfpunkt vor dem Code (§2),
  `make generated-sync` im Gate-Lauf. — **Ausgang:** offen bis Closure.
- **Spec-Lücken, die der Plan festlegt** (Header-Teilung am ersten `=`, Leerraum,
  Header ohne Endpunkt, Zeitpunkt des ersten Versuchs, gemeinsamer Fehlerzustand
  für Lesen und Übertragen): die Spec nennt sie nicht; jede Festlegung kann der
  Reviewer anders lesen, und eine Spec-Klarstellung wäre ein Zug des
  Planners/Architect, kein Zug dieses Slice. — **Ausgang:** offen bis Closure
  (Kandidat: benannte Spec-Lücke im Lerneintrag).
- **Ein Wert des Exports weicht von der Sicht ab, ohne dass ein Test es merkt**
  (Skalierung, Einheit, Attribut): der Test am Wire liest den Körper zurück und
  vergleicht gegen die Sicht derselben Datenbank, nicht gegen ein
  Test-Literal derselben Quelle
  ([`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt`](../observations/BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt/observation.md)
  nennt die Form; der Realserver-Beleg ist der Folge-Slice). — **Ausgang:** offen
  bis Closure.
- **Die Warnung trägt eine Zugangsdaten-Information** (Header-Wert, Userinfo
  der URL im Fehlertext des HTTP-Clients: `net/http` nennt die URL im Fehler): die
  Log-Zeile wird aus Code und Klasse gebaut, nicht aus dem rohen Fehlertext
  der Bibliothek; ein Test (§2) belegt die Abwesenheit. — **Ausgang:** offen bis
  Closure.
- **Die Aufzählung der Warn-Bereiche bleibt an einem Träger bei 1 bis 5**
  (`Area`, Kommentar, Test oder Handbuch): dann färbt der Test rot oder, schlimmer,
  der Code `PCF-W6001` gilt als kein Warncode. Gegenmaßnahme: der Suchlauf
  Zeile 1 und `make meldungscodes-check`. — **Ausgang:** offen bis Closure.
- **Der Slice ist größer als drei Punkte tragen** (Konfiguration plus Pipeline plus
  Bibliothek-Prüfpunkt): Rückführung nach §4 mit dem benannten
  Teilungsvorschlag. — **Ausgang:** offen bis Closure.
- **Das Runtime-Image trägt für `https` keine CA-Zertifikate** (distroless;
  *nach Kenntnisstand* enthalten, ungeprüft): im Folge-Slice zu messen (Übergabe-Block
  Punkt 4); dieser Slice liefert dazu keine Aussage und das Handbuch nennt für
  `https` keine Zusage. — **Ausgang:** offen bis Closure (übergeben an
  `otlp-metrik-export-e2e`).
- **Das Handbuch zieht nicht mit** (Betreiber-Oberfläche wächst um drei
  Variablen, ein Datei-Feld, zwei Codes und einen Bereich):
  [`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`](../observations/BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/observation.md)
  (drei Evidenz-Dateien am Planungsstand) — eigener Liefer-Punkt in §2. —
  **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*Der Plan füllt diese Sektion nicht; sie wird bei der Closure vor dem
`git mv` nach `done/` geschrieben (§2: die Paarungs-Zeile nach dem `git mv`).*

- **Was hat funktioniert:** (bei der Closure zu füllen)
- **Was ging anders als geplant:** (bei der Closure zu füllen)
- **Steering-Loop-Eintrag:** (bei der Closure zu füllen; Kandidaten §5)
- **Beobachtungs-Register (`../observations/`):** (bei der Closure zu füllen;
  zu lesen sind die in §6 und §8 genannten Einträge)
- **Folge-Slices:** `otlp-metrik-export-e2e` — Realserver-Beleg des Exports
  (Gegenstand im Übergabe-Block §1); **noch keine Datei**, der Planner legt sie
  mit der Closure dieses Slice an (DoD §2, Gate- und Lauf-Pflichten)
- **Risiken aus §6:** (bei der Closure zu füllen, je genau ein Ausgang)
- **Drei Paarungen:** (bei der Closure zu füllen: Anker · Folge-Slice · Register)

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration führt **eine** Sub-Area (`*`, Kürzel `PGC`, Greenfield); die
berührten Pfade (`internal/…`, `go.mod`, `docs/user/`, `harness/`) liegen alle in
ihr, es entsteht keine neue Sub-Area — die Schwelle (≥ 2 von 3 Achsen) wird nicht
angewandt, weil kein Pfad ausdifferenziert werden muss.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/BEO-PGC/`) am Planungsstand durchgegangen;
Treffer für die berührte Sub-Area, Zähler = Zahl der `evidence/`-Dateien:
`adapter-unittest-verdeckt-bootstrap-luecke` (5, erreicht 3× — Gegenmaßnahme:
Verdrahtungs-Test gegen reale Datenbank, §2 (d); der Container-Beleg ist der
Folge-Slice), `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (3,
eigener Liefer-Punkt), `handbuch-versionshistorie-uebersprungen` (3,
Historien-Zeile), `zwei-quellen-drift-handbuch-gegen-pflichtenheft` (3, Handbuch
und Pflichtenheft werden für die Kennzahlen-Tabelle gegeneinander gelesen),
`zahl-in-traeger-driftet-gegen-die-messung` (Zahlen im Handbuch tragen ihren
Ursprung, [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
Instanz A: die „neun“ und „zehn“ Kennzahlen stehen mit Messung oder als
übernommen), `intern-kennungen-in-ausgelieferten-texten` (1, Meldungstexte ohne
Kennung, `make ausgabe-kennungen-check`), `kommentar-herkunft-als-kette` (4,
höchstens eine Kennung je Kommentar im neuen Code,
[`AGENTS.md`](../../../../AGENTS.md) §3.7), `negativtest-ohne-bindung-an-seine-eingabe`
(25, Gegenprobe je Negativfall), `adapter-pflicht-eines-ports-ohne-traeger-im-adapter-slice`
(1, die Pflichten am Port stehen im Doc-Kommentar **und** im DoD des Adapters),
`e2e-metrik-boundary-nur-reviewer-belegt` (Zähler am Planungsstand nicht gezählt;
der Eintrag gehört dem Folge-Slice). Die Einträge mit mindestens drei Dateien
sind **mit** diesem Slice berührt und tragen deshalb hier ihre Gegenmaßnahme.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
