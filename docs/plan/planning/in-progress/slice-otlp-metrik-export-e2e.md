# Slice otlp-metrik-export-e2e: der OTLP-Metrik-Export am laufenden Feed-Container gegen einen echten OpenTelemetry-Collector

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

**Abhängigkeit:** [`slice-otlp-metrik-export`](../done/slice-otlp-metrik-export.md)
(Teil a) liegt in `done/` — Konfiguration, Lese-Port, Sender und Takt sind dort
geliefert und auf Unit- und `httptest`-Ebene belegt. **Dieser Slice ist
Voraussetzung des Server-Releases** (Entscheidung des Auftraggebers: erst der
Beleg gegen den echten Collector, dann der Release; Release 0.7.0 und die
Lizenzhinweise für `grpc-gateway/v2` nennt das Releasing, nicht dieser Slice).
**Frist zweier Adressen, die vor dem Release geklärt sein müssen** (sie
stammen aus dem Vorgänger und bleiben dort Architect-/Spec-Fragen, §6): die
Einheit von `cdc_consumer_lag` (`1` oder `By`) und der Spec-Zug zu den
Header-Regeln und zum Benutzerteil der Endpunkt-URL. Der Beleg dieses Slice
misst die Einheit am Wire, er entscheidet sie nicht (§1).

**Bezug:** [`LH-FA-SST-010`](../../../../spec/lastenheft.md) (Scope),
[`LH-FA-SST-004`](../../../../spec/lastenheft.md) und
[`LH-QA-OPS-003`](../../../../spec/lastenheft.md) (die SQL-Sicht als Gegenseite
der Wertgleichheit),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (WAL-Rückstand),
[`LH-QA-POR-003`](../../../../spec/lastenheft.md) (Realserver-Beleg am
Compose-Stack),
[`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) (Fitness-Function-Zeile
„Integrationstest“, Festlegung 8: kein Prometheus-Endpunkt),
[`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Pin-Form `image:tag@sha256:<Index-Digest>`, Sensor P10),
[`ADR-0030`](../../adr/0030-testpyramide.md) (E2E-Tier),
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) (E2E-Workflow),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
(Warn-Bereich 6),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).

**Berührte Spec-Stellen:** [`SPEC-033`](../../../../spec/pflichtenheft.md)
(Drahtform, Kennzahlen, Verhalten — das Test-Soll dieses Slice) ·
[`SPEC-009`](../../../../spec/pflichtenheft.md) (Namen der Kennzahlen) ·
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Warn-Bereich 6) ·
[`SPEC-012`](../../../../spec/pflichtenheft.md) (PostgreSQL-Versionen der
E2E-Matrix).
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
liegt mit [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) vor; der
Empfänger — ein echter Collector statt eines Wegwerf-Empfängers — ist eine
Entscheidung des Auftraggebers). **Datum:** 2026-10-04.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass:** Der Vorgänger belegt den Export gegen einen `httptest`-Empfänger und
gegen die reale Datenbank; er belegt nicht, dass ein **echter** OTLP-Empfänger
die Nachricht liest, und nicht den laufenden Container. Die Fitness-Function-Zeile
von [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) („Wert im
Export = Wert der Sicht im selben Moment“) steht dort als *erwartet, nicht
gefahren*. Die Spec-Stelle [`SPEC-033`](../../../../spec/pflichtenheft.md) macht
die Drahtform verbindlich; ob ein produktiver Empfänger sie so versteht, ist
ohne ihn offen. Am Planungsstand (Parent `5e75d923c2ff2d38a4fcdf839a66f813c1b2db74`)
trägt weder der Integrationsrunner noch `compose.yaml` einen OTLP-Bezug
(§3 Suchlauf Zeile 1: 0 Trefferzeilen).

**Entscheidung des Auftraggebers: der Beleg läuft gegen einen echten
OpenTelemetry-Collector**, nicht gegen einen Wegwerf-Empfänger. Folge: die
Empfänger-Form des Übergabe-Blocks (§2, Punkt (1)) wird hier zum Collector; die
Aussagen darüber, was der Collector kann, sind *nach Kenntnisstand, ungeprüft*
und werden im Prüfpunkt (§2 Liefer-Punkt 1) gemessen, bevor der Runner-Code
entsteht.

**Ziel:** Ein neuer Runner-Abschnitt im Integrationsrunner startet einen
gepinnten Collector im Compose-Netz, setzt am **laufenden** Feed-Container
`CDC_OTLP_ENDPOINT`, `CDC_OTLP_HEADERS` und einen kurzen Takt, und belegt am
Ausgang des Collectors: alle zehn Kennzahlen kommen als `Gauge` mit Einheit,
Datenpunkt-Attributen und Resource-Attributen an, ihre Werte sind gleich der
SQL-Sicht (zeitabhängige mit benannter Toleranz), der konfigurierte Header ist
am Collector wirksam, ein ausgefallener Collector hält den Container nicht auf
und die Übertragung setzt wieder ein, ungültige Konfigurationen beenden den
Start mit Klasse und Code, und ein `https`-Empfänger ist mit dem
Vertrauensspeicher des Prozesses belegt oder die Grenze ist benannt. Die
`abdeckung_declare`-Zeile für [`LH-FA-SST-010`](../../../../spec/lastenheft.md)
macht die Anforderung in `make doc-trace` zur gedeckten (am Parent: 1 Waise,
§3 Suchlauf Zeile 5).

**Festlegung zu `https` (Entscheidung dieses Plans, wo der Auftrag sie offenließ):**
`https` ist **Teil** dieses Slice, soweit es gemessen werden kann. Zwei getrennte
Aussagen: (a) *Enthält das Runtime-Image ein CA-Bündel?* — wird am Image
gemessen (Export des Dateisystems ohne Start, kein Netz), nicht aus Kenntnisstand
übernommen; (b) *Verifiziert der Sender gegen den Vertrauensspeicher des
Prozesses?* — belegt am Collector mit TLS-Empfänger und einem vom Runner erzeugten
Zertifikat (`tools/harness/certgen`, Muster der TLS-Phase des Vorgänger-Slice
[`slice-tls-http-grpc-server`](../done/slice-tls-http-grpc-server.md)): die
Verbindung gelingt, wenn der Prozess das Zertifikat über die Umgebungsvariable
des Go-Laufzeitsystems `SSL_CERT_FILE` als Vertrauensanker kennt (Mechanismus *nach
Kenntnisstand, ungeprüft*, im Slice zu messen), und scheitert ohne sie mit
`PCF-W6001` — die Gegenprobe belegt, dass wirklich verifiziert wird. Beides
sagt **nichts** über öffentliche Zertifizierungsstellen und über eine
konfigurierbare eigene CA (das bleibt Nicht-Ziel, §1 unten); das Handbuch trägt
nur, was gemessen ist. Gelingt (b) im Slice nicht (Collector-TLS-Konfiguration,
Zertifikatsform), wird der Teil nicht still gestrichen: die Grenze steht als
benannte Lücke mit Adresse in §7.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Prometheus-Scrape, Dashboards, Alarmierung, Traces, Logs** — Out-of-Scope von
  [`LH-FA-SST-010`](../../../../spec/lastenheft.md),
  [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) Festlegung 8.
  Dass Prometheus über einen Collector mit OTLP-Receiver erreichbar ist, bleibt
  *nach Kenntnisstand, nicht gemessen*: der Beleg dieses Slice fährt keinen
  Prometheus-Exporter und keinen Scrape, das Handbuch trägt keine Zusage dazu.
- **Die Entscheidung über die Einheit von `cdc_consumer_lag`** — eine Frage an den
  Architect (Folge-ADR bei `By`, weil [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md)
  `Accepted` ist). Dieser Slice **misst** die Einheit, die am Wire ankommt (der
  Collector zeigt sie), und trägt sie in die Abdeckungszeile ein; er ändert Code,
  Spec und SQL-Kommentar nicht.
- **Eine Änderung am Export-Code** (Adapter, Use Case, Konfiguration): der Slice
  ist ein Beleg. Ein Befund am Wire ist ein Fix-Slice oder eine Fixrunde mit
  eigenem Plan-Vermerk, kein stilles Mitändern (Rückführung §4).
- **Eine eigene CA, Client-Zertifikat, Proxy-Einstellungen für den Empfänger** —
  Out-of-Scope des Vorgängers; ein Betreiber-Bedarf wäre eine neue Anforderung.
- **Last- und Langzeitverhalten des Exports** (Fußabdruck im Betrieb, Verhalten
  bei sehr vielen Consumern) — keine Anforderung nennt eine Schwelle;
  Messen ohne Schwelle ist ein anderer Vorgang (Muster `make bench`).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

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

**Übergabe-Block des Vorgängers** (aus
[`slice-otlp-metrik-export`](../done/slice-otlp-metrik-export.md) §1,
als committeter Text übernommen; Punkt (1) und (4) sind durch die Entscheidung des
Auftraggebers und die Festlegung zu `https` in §1 konkretisiert, die Liefer-Punkte
unten tragen die geltende Form):

> Der Gegenstand, den der Realserver-Beleg dort tragen muss und den dieser Slice
> **nicht** trägt: (1) ein Wegwerf-Empfänger im Compose-Netz nimmt `POST /v1/metrics`
> des laufenden Feed-Containers entgegen, und ein übertragener Wert ist gleich dem
> Wert der Sicht im selben Moment (zeitabhängige Kennzahlen bis auf den Abstand der
> Messzeitpunkte; `cdc_wal_retention_bytes` aus dem Prozess) — die Fitness-Function-Zeile
> „Integrationstest (Compose, Empfänger-Wegwerf-Client)“ von `ADR-0149`; (2) Empfänger
> nicht erreichbar oder antwortet `500`: der Container läuft, Erfassung und
> Bestätigung unverändert, die Übertragung setzt beim nächsten Takt wieder ein;
> (3) ein ungültiger Endpunkt, Header oder Takt beendet den Start mit der Klasse
> `configuration` und dem Meldungscode in der Log-Zeile; (4) ein `https`-Endpunkt
> gegen den Vertrauensspeicher des Runtime-Images (nach Kenntnisstand enthält das
> distroless-Image CA-Zertifikate, ungeprüft); (5) eine `abdeckung_declare`-Zeile für
> `LH-FA-SST-010`, `docs/user/e2e-abdeckung.md` vom Runner geschrieben.

- [x] **Collector-Prüfpunkt, Pin und Happy Path (Liefer-Punkt 1,
      [`LH-FA-SST-010`](../../../../spec/lastenheft.md) Happy Path).**
      (a) **Prüfpunkt als erster Schritt, bevor Runner-Code entsteht** (alle Aussagen
      über den Collector sind bis dahin *nach Kenntnisstand, ungeprüft*): welches
      Image trägt die benötigten Komponenten. Kandidaten: `otel/opentelemetry-collector`
      (Kern) und `otel/opentelemetry-collector-contrib`. Gemessen und im Bericht
      genannt (Befehl `components` des Collector-Binarys im Container, gedruckte
      Zeile je Komponente): ein OTLP-HTTP-Empfänger; ein Exporter, dessen Ausgabe der
      Runner lesen kann — bevorzugt `file` (eine JSON-Zeile je Export, maschinell
      prüfbar) oder, wenn das gewählte Image ihn nicht trägt, `debug` (Text in
      `docker logs`); und für die Header-Wirkung eine Authentifizierungs-Erweiterung
      (Kandidat `bearertokenauth`, Bearer-Token am Empfänger geprüft) oder ein
      Weg, die Anfrage-Header sichtbar zu machen. **Entscheidungsregel:** das kleinste
      Image, das alle drei trägt; der Bericht nennt je Kandidat die Messung und
      den Grund. Trägt keines die Header-Beobachtung, hält der Slice an (§4).
      (b) **Pin nach [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md):**
      `image:tag@sha256:<Index-Digest>`, der Digest im Slice gemessen
      (`docker buildx imagetools inspect <image:tag> --format '{{.Manifest.Digest}}'`, die
      gedruckte Zeile im Bericht; **kein Digest-Literal in diesem Plan**, nur das
      Verfahren); der Pin steht als Variable mit Default im Runner (Muster `PG_TEST_IMAGE`
      in `compose.yaml`) — ein Pin an genau einer Stelle. Zu belegen:
      `make pin-stale-all` druckt für diesen Pin weder `DRIFT` noch `UNBESTIMMT`
      (Netz nötig; der Lauf steht im Bericht, die Ausgabezeile des Pins gelesen).
      (c) **Collector-Konfiguration und Aufbau:** OTLP/HTTP-Empfänger auf einer Adresse
      im Compose-Netz (`cdc-feed-test`), der Exporter nach (a), die
      Konfigurationsdatei im Temp-Verzeichnis des Runners (Muster `TLS_TMP`:
      `cleanup`-Trap räumt es, nichts im Arbeitsbaum). Der Aufbau läuft über eine
      Compose-Override-Datei im Temp-Verzeichnis (`compose.yaml` bleibt unverändert);
      der Feed-Container wird am Phasen-Ende **ohne** Override wiederhergestellt
      (Muster `tl_recreate`/Wiederherstellung der Phasen „API-Token-Wechsel“ und
      „TLS-Schnittstellen“ im Runner).
      (d) **Runner-Phase** im Integrationsrunner (`tools/harness/run-integration-tests.sh`):
      Feed-Container mit `CDC_OTLP_ENDPOINT` auf den Collector, `CDC_OTLP_INTERVAL_SECONDS`
      = 5 (die untere Grenze), `CDC_OTLP_HEADERS` mit einem Authorization-Header
      (Wert ein Platzhalter-Token der Phase). Ein Wegwerf-Werkzeug (Go, unter
      `tools/harness/`, Muster der anderen Wegwerf-Clients, im Toolchain-Container,
      `--network none` wo der Zugriff es erlaubt) liest die Ausgabe des Collectors und
      die Zeilen von `SELECT metric_name, label, value FROM cdc.metrics` (der Runner
      schreibt sie in eine Datei im Temp-Verzeichnis) und prüft: **zehn Namen** (die der
      Tabelle in [`SPEC-033`](../../../../spec/pflichtenheft.md)), je Name Typ `Gauge`,
      Einheit, Datenpunkt-Attribut (`consumer` bzw. `class`) und die
      Resource-Attribute `service.name` = `pg-change-feed` und `cdc.source_id` = die
      Quellkennung des Containers; **Werte gleich der Sicht** mit einem Consumer und
      einem Fehlerzustand als Zeilen mit Label; die Toleranz für `cdc_capture_lag` und
      `cdc_oldest_change_age_seconds` (hängen an `now()`) nennt der Test mit Zahl und
      Grund (Abstand der Messzeitpunkte, höchstens ein Takt plus Frist), alle anderen
      Werte sind **exakt gleich**, `cdc_wal_retention_bytes` wird gegen die Messung des
      Prozesses gehalten, nicht gegen die Sicht (der Weg der Gegenlesung —
      `pg_replication_slots` des Slots — steht im Bericht). Eine **gemessene**
      Einheit von `cdc_consumer_lag` wird gedruckt und in die Abdeckungszeile
      übernommen. Gegenprobe zum Header: ohne Header (zweiter Start, Collector mit
      Authentifizierung) antwortet der Collector abweisend und der Container warnt mit
      `PCF-W6001`; mit Header kommen Daten an — der Test nennt beide Richtungen.
      (e) **Abdeckung:** eine `abdeckung_declare`-Zeile je Phase für
      [`LH-FA-SST-010`](../../../../spec/lastenheft.md); der Runner schreibt
      `docs/user/e2e-abdeckung.md` (Erzeugnis, **kein** Lauf-Beleg); zu belegen durch
      `make doc-trace`: 0 Waisen (am Parent 1, Zeile `LH-FA-SST-010`; §3 Suchlauf
      Zeile 5).
      (f) **Mutationsproben (Zusage, an Kopien im Scratchpad, nicht an Dateien im
      Arbeitsbaum — [`AGENTS.md`](../../../../AGENTS.md) §3.1; §3.12):** vier Stellen
      einzeln, je Stelle, Instanz und gesehene Farbe im Bericht — ein Wert im
      Werkzeug-Soll um 1 verschoben (die Wertgleichheit färbt rot), die Einheit eines
      Namens im Soll geändert (rot), das Datenpunkt-Attribut `class` im Soll durch
      `consumer` ersetzt (rot) und der Header-Beleg ohne die Gegenprobe (ein Start ohne
      Header zählt als Erfolg: der Test bleibt grün, weil die Gegenprobe fehlt — das
      belegt, dass sie trägt; **erwartet**, im Lauf zu bestätigen). Die Verallgemeinerung
      auf weitere Stellen ist **hergeleitet**, nicht erprobt.
- [x] **Ausfall, Wiederaufnahme und ungültige Konfiguration (Liefer-Punkt 2,
      [`LH-FA-SST-010`](../../../../spec/lastenheft.md) Negative, Boundary).**
      (a) *Collector angehalten:* `docker stop` des Collectors bei laufendem
      Container: der Container bleibt `running` über eine Beobachtung von mindestens
      zwei Takten plus Frist, `--healthcheck` bleibt gesund, die Erfassung läuft
      unverändert (eine danach committete Zeile ist über `cdc.changes` lesbar), der
      Heartbeat trägt keine `error_class`; das Log trägt **genau eine** Warn-Zeile mit
      `PCF-W6001` je Beobachtungsfenster (die Drosselung von 5 Minuten, am echten
      Takt: im Fenster keine zweite) und **keinen** Header-Wert und keinen Teil der
      URL mit Benutzerteil (Test auf Abwesenheit des Platzhalter-Tokens im gesamten
      Log). (b) *Wiederaufnahme:* `docker start` des Collectors: beim nächsten Takt kommt
      ein Export an (die Ausgabe des Collectors wächst), das Log trägt die
      Wiederaufnahme-Info-Zeile ohne Meldungscode, die Zahl der Exporte je Takt ist
      **eins** (kein Nachholen der ausgefallenen Takte: die Zahl der Zeilen seit
      Wiederaufnahme über ein Fenster von k Takten ist k oder k+1, nicht die der
      verpassten). (c) *Antwort `500`:* der Empfänger antwortet abweisend (Collector mit
      Authentifizierung und falschem Token, oder ein vorgeschalteter Fehler, wie der
      Prüfpunkt es erlaubt): `PCF-W6001`, Container läuft; die Form der Abweisung
      (Status) steht im Bericht. (d) *Ungültige Konfiguration:* je ein Start mit
      einem ungültigen Endpunkt (Schema `ftp`), einem Header ohne `=` und einem Takt
      4 beendet den Container mit Ausgang 2, der Log-Zeile `Fehlerklasse
      configuration [<Code>]` mit dem jeweiligen Meldungscode aus der Code-Tabelle
      (nicht im Plan festgeschrieben: der Test liest ihn aus der Tabelle des Vorgängers)
      und **ohne** Slot (Muster `tl_expect_start_refused`); Gegenprobe mit gültigem
      Wert neben jedem Negativfall
      ([`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`](../observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/observation.md)).
      Wiederherstellung am Phasen-Ende: Container ohne Override, die Folge-Phasen
      laufen unverändert.
      (e) **Mutationsprobe (Zusage, Kopie im Scratchpad):** zwei Stellen — die
      Beobachtungsfrist ohne den Ausfall (der Collector bleibt oben: die Prüfung auf
      `PCF-W6001` muss rot werden) und die Gegenprobe des Konfigurationsfalls
      (gültiger Wert wird abgelehnt: rot). Instanz und Farbe im Bericht.
- [x] **`https`, Handbuch und Abdeckung (Liefer-Punkt 3).**
      (a) *CA-Bündel im Image, gemessen:* die Aussage „das Runtime-Image enthält
      CA-Zertifikate“ wird am gebauten `:dev`-Image gelesen (Export des
      Dateisystems ohne Start; der Pfad der Bündel-Datei und ein Treffer im Export,
      gedruckte Zeile im Bericht), nicht aus dem Kenntnisstand übernommen.
      (b) *TLS-Empfänger:* der Collector bedient `https` mit einem von `certgen`
      erzeugten Paar (Temp-Verzeichnis, nicht im Arbeitsbaum); der Feed-Container mit
      `https`-Endpunkt und dem Zertifikat als Vertrauensanker über die Umgebung (nach (a)
      der Mechanismus der Wahl, im Slice gemessen) überträgt, die Werte sind gleich
      der Sicht; **Gegenprobe:** derselbe Endpunkt ohne den Vertrauensanker endet
      im Warn-Zustand (`PCF-W6001`), der Container läuft, die Ausgabe des Collectors
      wächst nicht. Gelingt (b) nicht, steht die Lücke in §7 mit Adresse (§1).
      (c) *Handbuch* (`docs/user/benutzerhandbuch.md`, kennungsfrei,
      [`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`](../observations/BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/observation.md)):
      **nur falls der Beleg etwas Neues zeigt** — der Satz „Ob ein `https`-Empfänger
      erreichbar ist, hängt vom Vertrauensspeicher des Container-Images ab“ wird durch
      das Gemessene ersetzt (CA-Bündel vorhanden ja/nein, Wirkung der Umgebungsvariable
      ja/nein, gemessen mit Collector-Version); die gemessene Einheit von
      `cdc_consumer_lag`, wenn sie von der Handbuch-Tabelle abweicht, **nicht**
      still ändern, sondern als Befund an den Architect (§6); keine Zusage zu
      Prometheus; Änderungshistorie-Zeile mit der nächsten freien Nummer
      ([`BEO-PGC/handbuch-versionshistorie-uebersprungen`](../observations/BEO-PGC/handbuch-versionshistorie-uebersprungen/observation.md)).
      Zeigt der Beleg nichts Neues, bleibt das Handbuch unverändert und der Bericht
      nennt das. `make handbuch-public-doc-check` Exit 0.
      (d) `harness/README.md` §Sensors, Zeile `make test-integration`: ein Satz zur
      neuen Phase (Collector-Image als Variable, Gegenstand der Belege); die Zeile ist
      lang, der Satz wird angehängt, nicht eingefügt.

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [x] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0,
      `make test-integration` Exit 0 (der Runner trägt die neue Phase; die
      gedruckte Schlusszeile mit der Zeilenzahl der Abdeckungstabelle steht im
      Bericht), `make doc-trace` Exit 0 mit 0 Waisen, `make pin-stale-all` für den
      neuen Pin ohne `DRIFT`, `make fmt-check` Exit 0 für das Wegwerf-Werkzeug,
      `make test` Exit 0 (das Werkzeug läuft in `make test`, wenn es Tests trägt).
- [ ] **Realer Post-Push-Lauf** ([`AGENTS.md`](../../../../AGENTS.md) §3.10): der
      Runner ändert den Inhalt von `e2e.yml`s Ziel, nicht die Datei; der Lauf holt in
      beiden PostgreSQL-Legs ein zusätzliches Image von Docker Hub und läuft länger.
      Das Risiko (§6) bleibt **offen**, bis ein Lauf von `e2e.yml` am Stand der
      Closure beide Legs `success` zeigt (`gh run list --workflow e2e.yml`); der
      Ausgang wird nachgetragen, bevor Closure erfolgt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update (siehe Liefer-Punkt 3 (c), (d)); gemeldete Träger fremder Dateien
      mit der Closure nachgezogen (§3 Suchlauf,
      [`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` | update | neue Phase(n) mit `abdeckung_declare`, Pin-Variable des Collector-Images, Temp-Verzeichnis im `cleanup`-Trap (Muster `TLS_TMP`), Wiederherstellung am Phasen-Ende |
| `tools/harness/otlpcheck/` (`main.go`, `check.go`, `check_test.go`) | neu | Wegwerf-Werkzeug: liest die Ausgabe des Collectors und die Zeilen der Sicht, prüft Namen, Typ, Einheit, Attribute, Resource-Attribute, Wertgleichheit samt Toleranz; Muster der anderen Wegwerf-Clients unter `tools/harness/`; die Tabellentests laufen in `make test` |
| Collector-Konfiguration (Temp-Verzeichnis des Runners, im Runner als Text erzeugt) | neu, nicht im Baum | Empfänger, Exporter, Header-Beobachtung; kein Dateiname im Arbeitsbaum |
| `docs/user/e2e-abdeckung.md` | update (Erzeugnis des Runners) | die Zeilen für `LH-FA-SST-010`; **kein** Lauf-Beleg |
| `docs/user/benutzerhandbuch.md` | update (nicht bedingt) | der Beleg zeigt Neues (CA-Bündel, `SSL_CERT_FILE`, Einheit `By`, Benutzerteil der URL), §2 Liefer-Punkt 3 (c); Version 1.98 |
| `harness/README.md` | update | Zeile `make test-integration`: ein Satz zur neuen Phase |
| `internal/adapters/driven/otlpexport/exporter.go` und `exporter_test.go` | update (Auftrag des Auftraggebers nach Spec-Stand, nicht im ursprünglichen Plan) | die Einheit von `cdc_consumer_lag` ist `By`; ein Test für den Benutzerteil der Endpunkt-URL (`Authorization: Basic …`, ein Header der Konfiguration gewinnt) — die Entscheidung liegt in den Folge-ADRs zu Einheit und Spec-Lücken, dieser Slice führt sie aus |

- **Reihenfolge:** (1) Startmessung: Suchlauf unten am Arbeitsstand neu (die
  Slices davor bewegen Handbuch und Runner); (2) Collector-Prüfpunkt (§2 Liefer-Punkt 1
  (a)) und Entscheidung des Images; Pin messen; (3) Collector-Konfiguration und Aufbau,
  Werkzeug mit Tabellentests; (4) Happy-Path-Phase; (5) Ausfall-, Wiederaufnahme- und
  Konfigurations-Phase; (6) `https`; (7) Mutationsläufe an Kopien; (8) Handbuch (bedingt),
  README, `make gates`, `make test-integration`, `make pin-stale-all`, `make doc-trace`.
- **Laufzeit:** der Runner läuft heute lang (am Parent ca. 5961 Zeilen,
  `wc -l tools/harness/run-integration-tests.sh`); die neue Phase hängt am Ende der
  Reihe vor der Schluss-Wiederherstellung und darf Phasen davor nicht stören (Risiko §6).
- **Zeitverhalten der Beobachtung:** der Takt ist 5 s, die Beobachtungsfenster
  der Ausfall-Phase sind ein Mehrfaches davon; die erste Übertragung erfolgt nach
  einem vollen Takt (Festlegung des Vorgängers, nicht dieser Slice); der Test wartet
  auf die Ausgabe mit einer Frist, nicht mit einem festen Schlaf.

```suchlauf
5e75d923c2ff2d38a4fcdf839a66f813c1b2db74 0 -n -i -E 'otlp|opentelemetry' -- tools/harness compose.yaml examples/compose.yaml .github Makefile harness/mk
5e75d923c2ff2d38a4fcdf839a66f813c1b2db74 0 -n 'LH-FA-SST-010' -- docs/user harness tools .d-check.yml
5e75d923c2ff2d38a4fcdf839a66f813c1b2db74 56 -n -E '^abdeckung_declare "' -- tools/harness/run-integration-tests.sh
5e75d923c2ff2d38a4fcdf839a66f813c1b2db74 1 -n -E '@sha256:[0-9a-f]{64}' -- tools/harness/run-integration-tests.sh
5e75d923c2ff2d38a4fcdf839a66f813c1b2db74 1 -n -E 'Waise' -- harness/README.md
5e75d923c2ff2d38a4fcdf839a66f813c1b2db74 2 -n -E 'Vertrauensspeicher|otel-collector' -- docs/user/benutzerhandbuch.md
diff 67 -n -i -E 'otlp|opentelemetry' -- tools/harness compose.yaml examples/compose.yaml .github Makefile harness/mk
diff 12 -n 'LH-FA-SST-010' -- docs/user harness tools .d-check.yml
diff 59 -n -E '^abdeckung_declare "' -- tools/harness/run-integration-tests.sh
diff 2 -n -E '@sha256:[0-9a-f]{64}' -- tools/harness/run-integration-tests.sh
diff 1 -n -E 'Waise' -- harness/README.md
diff 1 -n -E 'Vertrauensspeicher|otel-collector' -- docs/user/benutzerhandbuch.md
```

- **Suchlauf (§3.13 der Regeln, [`AGENTS.md`](../../../../AGENTS.md)).** Bewegte
  Eigenschaften: der **OTLP-Bezug im Runner** (Hedge „otlp“, „opentelemetry“), die
  **Deckung von `LH-FA-SST-010`** (Symbol), die **Zahl der Abdeckungs-Phasen** im Runner
  (Zählwort: die Zeilen `abdeckung_declare`), die **Digest-Pins im Runner** (Muster
  `@sha256:`), die **Waisen-Aussage** in `harness/README.md` (Beschreibung) und die
  **Empfänger- und `https`-Aussage** im Handbuch (Hedge „Vertrauensspeicher“, Beispiel
  „otel-collector“). Parent ist `5e75d923…` (`git rev-parse HEAD` am Planungsstand,
  nie `HEAD`). Suchraum je Zeile ein benannter Pathspec, nicht der ganze Baum: die
  bewegte Eigenschaft liegt in Runner, Compose-Dateien, Workflows, Harness-Nachschlage-
  Doku und Handbuch; **Grund der Einschränkung:** die Eigenschaften sind
  Zeilenmengen dieser Dateien, ein Baum-Lauf träfe Plan- und Review-Texte, die der
  Suchraum ausschließt. **Soll-Zahlen sind am Planungsstand gemessen (`make suchlauf-nachmessen`,
  Parent `5e75d923…`: 0 · 0 · 56 · 1 · 1 · 2; Zeile 3 zählt die Aufrufe mit
  Anführungszeichen — das Muster ohne es trifft auch die Definition der Funktion
  `abdeckung_declare()` und zählt 57); die `diff`-Zeilen tragen den Planungsstand
  und werden am Endstand durch den Implementer auf die Messung gesetzt** — Erwartung am Endstand (hergeleitet): Zeile 1
  wächst um die Phase und die Variable, Zeile 2 um die Abdeckungs-Zeile(n) in
  `docs/user/e2e-abdeckung.md`, Zeile 3 um die neuen Phasen, Zeile 4 um den Pin,
  Zeile 5 sinkt nicht ohne Lesen der Treffer, Zeile 6 ändert sich nur, wenn das
  Handbuch mitzieht. Jede Trefferzeile wird gelesen, nicht gezählt; der Implementer
  misst beim Start am Arbeitsstand neu (neue Commit-Kennung als Parent, Abweichung mit
  Ursache im Bericht).

### Umsetzungsnachzug (Implementer)

Alle Angaben stammen aus Läufen dieses Slice am Arbeitsstand (gemessen); wo eine
Aussage über den Lauf hinausgeht, steht sie als *hergeleitet*.

**Auftrag über den Plan hinaus** (Entscheidung des Auftraggebers nach den Folge-ADRs zu
Einheit und Spec-Lücken; die Ausschlüsse in §1 „Entscheidung über die Einheit“ und
„Änderung am Export-Code“ sind damit für genau diese Punkte überholt, der Rest von §1
gilt): die Einheit von `cdc_consumer_lag` ist `By` in `exporter.go` und
`exporter_test.go`; ein Go-Test für den Benutzerteil der Endpunkt-URL
(`TestExportSendsBasicAuthorizationFromEndpointUserinfo`); das Handbuch zieht nach
(Einheit `By`, WAL-Strecke in Bytes, `https`, Benutzerteil der URL, Version 1.98).

**Collector-Prüfpunkt (Liefer-Punkt 1 (a)).** `docker run --rm otel/opentelemetry-collector:0.162.0 components`
nennt den Empfänger `otlp`, die Exporter `file` und `debug`, die Prozessoren `resource` und
`attributes` und als Erweiterungen nur `health_check`, `pprof` und `zpages` — **keine**
Authentifizierungs-Erweiterung. Entscheidung: das Kern-Image (126 578 012 Byte,
`docker image inspect`), nach der Entscheidungsregel das kleinste, das Empfänger,
auslesbaren Exporter und einen Weg zur Header-Beobachtung trägt: der Empfänger mit
`include_metadata: true` und der `resource`-Prozessor mit `from_context: authorization`
bilden den Anfrage-Header als Resource-Attribut `http.authorization`, das der
`file`-Exporter mitschreibt. Das Contrib-Image wurde nicht abgerufen und nicht gemessen.
Folge für den Plan: eine abweisende Authentifizierung (Plan (d) Gegenprobe, (c) Status
`500`) gibt es in diesem Image nicht; die Gegenprobe zum Header ist das **Fehlen** des
Attributs ohne `CDC_OTLP_HEADERS`, die abweisende Antwort ein unbekannter Pfad
(Status `404`, `detail` der Warnung `Status 404`).

**Pin (Liefer-Punkt 1 (b)).** Der Index-Digest ist mit
`docker buildx imagetools inspect otel/opentelemetry-collector:0.162.0 --format '{{.Manifest.Digest}}'`
gemessen und steht an genau einer Stelle: Variable `OTLP_COLLECTOR_IMAGE` mit Default im
Runner. `make pin-stale-all`: `pin-stale-all: 17 Referenzen — 17 OK, 0 DRIFT, 0 UNBESTIMMT`
(Exit 0, die Zeile des Pins `OK`).

**Abweichungen vom Plan.** (1) Die Phase hängt **nach** der TLS-Phase, nicht am Ende der
Reihe: der letzte Rundlauf (Routing-Nichtanwendbarkeit) lässt den Erfassungspfad beendet
stehen, die Phase braucht einen gesunden Container. (2) Die Header-Beobachtung und die
abweisende Antwort wie oben. (3) `cdc_wal_retention_bytes` steht nicht in der Sicht; die
Gegenlesung ist die **Messreihe des Prozesses** (die Log-Zeilen `replication: WAL-Rückstand
gemessen`): der exportierte Wert ist exakt eine dieser Messungen. Damit die Reihe nicht
trivial null ist, hält eine Sitzung die Persistierung an (SHARE-Sperre auf `cdc.change`,
die Sicht bleibt lesbar) und 30.000 Zeilen in einer nicht aktivierten Tabelle bilden den
Rückstand; der Runner verlangt eine größte Messung von mindestens 1 MiB. (4) Phase 1 fährt
vier Starts (mit Header, ohne Header, Benutzerteil der URL, Benutzerteil und Header). (5)
Ein Consumer mit bestätigter Position und eine erfasste Change sind Vorbedingung der Sicht
und werden von der Phase angelegt. (6) Der Ausgang `Laufzeit der drei OTLP-Phasen` steht in
der Ausgabe des Runners.

**Gedruckte Zeilen des Laufs `make test-integration` (Exit 0).**

```text
OTLP-Metrik-Export — Collector otelcol version 0.162.0 (Komponenten otlp, file, resource vorhanden)
OTLP-Metrik-Export — CA-Bündel im Image: etc/ssl/certs/ca-certificates.crt
mit Header: OTLPCHECK names=10 points=19 consumer_lag_unit=By wal=4831552 errors_labels=schema auth=present gemessen_im_prozess=4831552
ohne Header: OTLPCHECK names=10 points=19 consumer_lag_unit=By wal=0 errors_labels=schema auth=absent gemessen_im_prozess=0
https mit Vertrauensanker: OTLPCHECK names=10 points=19 consumer_lag_unit=By wal=0 errors_labels=schema auth=present gemessen_im_prozess=0
nach docker start kamen 3 Exporte in 16 s an
Laufzeit der drei OTLP-Phasen 252 s
Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 59 Bash-Zeilen
```

Die gemessene Einheit von `cdc_consumer_lag` am Draht ist `By`; `make doc-trace`
druckt `83 Anforderung(en), 0 Waise(n).` (am Parent: 1 Waise, `LH-FA-SST-010`).
Exit-Codes der Läufe am Stand der vier Commits des Slice (je Aufruf ungefiltert, der Exit
direkt gelesen): `make gates` 0 (darin `docs-check` 0 mit `d-check: 1675 Datei(en) geprüft,
0 Befund(e)`, `coverage-gate: OK — Coverage 83.30% erfüllt Schwelle 80%`),
`make test-integration` 0, `make test` 0, `make test-store` 0, `make doc-trace` 0,
`make pin-stale-all` 0, `make fmt-check` 0 (`357 Go-Dateien geprüft, alle formatiert`),
`make handbuch-public-doc-check` 0, `make suchlauf-nachmessen PLAN=` 0 (`12 Zeilen
stimmen`), `make kommentar-kennungen DIFF=<Parent>` 0 ohne Kandidat.
Gesamtzeit von `make test-integration`: 1037 s; die Gesamtzeit ohne die Phasen wurde
nicht gemessen, der Vergleich vor und nach (§6) bleibt offen. Collector-Version, CA-Bündel
und die Wirkung von `SSL_CERT_FILE` sind an Version 0.162.0 gemessen; ob öffentliche
Zertifizierungsstellen verifiziert werden, ist **nicht** gemessen.

**Mutationen (Zusage · mutierte Eingabe · gesehenes Rot), an Kopien im Scratchpad.**

| Zusage | mutierte Eingabe | gesehenes Rot |
|---|---|---|
| Wertgleichheit mit der Sicht | `otlpcheck`-Soll: Wert der Sicht um 1 verschoben, Instanz: Phase 1, Start mit Header | `sechs Versuche ohne bestandene Prüfung … Export 1, Sicht 1` |
| Einheit je Name | `otlpcheck`-Soll: Einheit von `cdc_storage_bytes` `By` zu `1` | `Einheit "By", erwartet "1"` |
| Datenpunkt-Attribut | `otlpcheck`-Soll: Attribut von `cdc_errors_total` `class` zu `consumer` | `Attribut consumer fehlt` |
| Header-Gegenprobe trägt | Collector setzt `http.authorization` unbedingt: „Start mit Header“ bleibt grün (erwartet, bestätigt), „Start ohne Header“ rot | `Resource-Attribut http.authorization ist vorhanden … erwartet abwesend` |
| Ausfall-Fenster | `docker stop` des Collectors entfernt | `Warn-Zeilen PCF-W6001 im Beobachtungsfenster () — erwartet '1', gelesen '0'` |
| Konfigurations-Gegenprobe | gültiger Endpunkt der Gegenprobe durch `ftp://…` ersetzt | `Gegenprobe 'Endpunkt mit Schema http' — … unhealthy, wollen healthy` |
| Basic aus dem Benutzerteil | Go-Test: `base.User = nil` in `New` | `Authorization = "", wollen Basic dXNlcjpwYXNz` |
| Vorrang des Konfigurations-Headers | Go-Test: `Authorization` aus den Headern der Konfiguration verworfen | `Authorization = "Basic dXNlcjpwYXNz", wollen Bearer abc` |
| Einheit `By` | Go-Test: Einheit von `cdc_consumer_lag` zu `1` | `cdc_consumer_lag: Einheit "1", wollen "By"` |
| Fehlertext ohne Benutzerteil | Go-Test: roher Transportfehler im Detail | `TestExportErrorsCarryNoCredentials` rot |

Nicht erprobt (hergeleitet): die Gegenprobe der `https`-Phase (Vertrauensanker), die
Untergrenze von 1 MiB und die Beobachtungsfenster von 16 s; die Verallgemeinerung auf
weitere Stellen eines Werkzeugs ist hergeleitet.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1), **und**
[`slice-otlp-metrik-export`](../done/slice-otlp-metrik-export.md) liegt in
`done/`, **und** der Implementer hat die Startmessung gefahren (die Suchlauf-Zeilen
am Arbeitsstand neu gemessen, `make suchlauf-nachmessen PLAN=` mit diesem Plan
Exit 0 nach Anpassung von Parent und Soll), **und** der Collector-Prüfpunkt (§2
Liefer-Punkt 1 (a)) hat eine Entscheidung getragen, bevor Runner-Code entsteht.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice wächst über die
  drei Punkte aus §2 — Kandidat für die Teilung: Liefer-Punkt 3 (`https`) als eigener
  Slice nach Punkt 1 und 2; der Planner teilt dann die Datei.
- `in-progress` → `open` (blockiert — Carveout?): (1) der Prüfpunkt findet kein Image,
  das Empfänger, auslesbaren Exporter und eine Header-Beobachtung trägt — dann ist es
  eine Auftraggeber-/Architect-Frage (Wegwerf-Empfänger nach
  [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md) als Rückfall, mit der
  Entscheidung des Auftraggebers); (2) ein **Befund am Wire** (der Collector liest die
  Nachricht anders als [`SPEC-033`](../../../../spec/pflichtenheft.md) sie vorgibt, oder
  ein Wert weicht von der Sicht ab) — der Befund ist ein Fix im Export-Code und damit
  ein eigener Slice oder eine Fixrunde des Vorgängers, nicht Teil dieses Beleg-Slice;
  der Slice hält an und der Planner schneidet den Fix. Kein Carveout, weil dann kein
  Gate rot ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, `make test-integration` endet mit Exit 0, `make doc-trace` meldet 0 Waisen,
ein realer Lauf von `e2e.yml` am Stand der Closure zeigt beide Legs `success`,
`make suchlauf-nachmessen PLAN=` mit diesem Plan endet am Endstand mit Exit 0, und die
Closure-Notiz in §7 trägt einen Lerneintrag (Kandidaten: die Erprobung der
„nach Kenntnisstand“-Annahmen über Collector und Vertrauensspeicher, die gemessene
Einheit von `cdc_consumer_lag`, eine benannte Spec-Lücke zu `https`/Vertrauensanker).
Ein Gate, das am Stand der Closure rot ist, geht nur mit dokumentiertem Carveout nach
`done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen
(Grund) · weiter offen (BEO-Eintrag im Register). Alle Ausgänge sind bis zur
Closure **offen** (Platzhalter `Ausgang: offen bis Closure`).

- **Docker-Hub-Abruflimit für das Collector-Image:** der Runner holt ein
  zusätzliches Image (der Kern-Collector ist klein, das Contrib-Image groß; Größen
  *nach Kenntnisstand, ungemessen*), lokal und in beiden Legs von `e2e.yml`. Ein
  anonymer Abruf kann am Limit scheitern, der Lauf endet dann rot ohne Aussage über
  den Export. Gegenmaßnahme: das kleinste tragende Image (Entscheidungsregel §2),
  Pin mit Digest (ein Abruf je Lauf), der Fehlschlag des Abrufs ist im Runner eine
  eigene Meldung, nicht ein Fehlschlag der Belege. — **Ausgang:** offen bis Closure.
- **Laufzeit des Runners:** der Runner ist bereits lang; die neue Phase (Takt 5 s,
  mehrere Takte je Beobachtung, zwei Neustarts des Feed-Containers) fügt Minuten
  hinzu, und eine Verlängerung von `e2e.yml` über ein Job-Limit hinaus wäre ein
  Fehlschlag ohne Bezug zum Export. Gegenmaßnahme: Fristen statt fester Schlafzeiten,
  die Gesamtzeit der Phase im Bericht (gemessen, Lauf), Vergleich der Gesamtzeit
  vor und nach (gedruckte Zeilen). — **Ausgang:** offen bis Closure.
- **Das Ausgabeformat des Collectors ist ungeprüft** (`file`-Exporter: JSON-Form,
  Zahlentypen als Zeichenkette oder Zahl, Zeitstempel; `debug`: Textform, die sich
  zwischen Versionen ändert): der Pin hält die Version, das Werkzeug liest nur,
  was der Prüfpunkt an der gepinnten Version gemessen hat; ein
  Versionswechsel durch `make pin-stale-all` und Hebung ist ein bewusster Commit mit
  erneutem Lauf. — **Ausgang:** offen bis Closure.
- **Port-/Netz-Kollision und Dateirechte:** der Collector läuft im Netz
  `cdc-feed-test`; ein Port auf dem Host wird nicht veröffentlicht (der Zugriff
  erfolgt über das Compose-Netz), die Ausgabe-Datei des Collector-Prozesses (ein
  fremder Benutzer im Image, *nach Kenntnisstand*) braucht ein beschreibbares
  Temp-Verzeichnis (Rechte im Runner, `cleanup` räumt es). Eine Kollision mit dem
  Container-Namen einer Vorphase oder ein nicht beschreibbares Verzeichnis ließe die
  Phase ohne Aussage enden. Gegenmaßnahme: eigener Container-Name mit Präfix der Phase
  im `cleanup`-Trap, Rechte im Runner gesetzt und gelesen (der Test schreibt vor dem
  Beleg eine Probe). — **Ausgang:** offen bis Closure.
- **Die Einheit von `cdc_consumer_lag` ist offen** (`1` laut
  [`SPEC-033`](../../../../spec/pflichtenheft.md) und
  [`ADR-0149`](../../adr/0149-otlp-metrik-export-mechanismus.md), „LSN-Byte-Abstand“ im
  SQL-Kommentar der Sicht): der Beleg misst nur, was am Wire ankommt, er entscheidet
  nicht. **Adresse:** Architect-Frage mit Frist **vor dem Release** (§7 des Vorgängers
  trägt dieselbe Adresse). — **Ausgang:** offen bis Closure (übergeben an den Architect).
- **Spec-Lücken, die im Vorgänger offen blieben:** die Header-Regeln von
  `otlp_headers` und der Benutzerteil der Endpunkt-URL (`http://user:pass@host`: wird
  nicht abgelehnt; `net/http` sendet ihn nach Herleitung als Basic-Authorization, **nicht
  gefahren**). Der Beleg dieses Slice kann das Letzte belegen oder widerlegen, **wenn**
  der Prüfpunkt es mit dem Collector (Authentifizierung `basicauth`) kostengünstig
  erlaubt; sonst bleibt es hergeleitet. **Adresse:** Spec-Zug des Auftraggebers/Architects,
  Frist **vor dem Release**. — **Ausgang:** offen bis Closure.
- **Der Collector bedient `https` nicht mit dem erzeugten Zertifikat oder die
  Umgebungsvariable `SSL_CERT_FILE` wirkt nicht** (beides *nach Kenntnisstand*):
  dann bleibt die Grenze benannt statt gestrichen (§1 Festlegung zu `https`). —
  **Ausgang:** offen bis Closure.
- **Ein Befund am Wire** (der Collector liest eine Zahlenform anders, eine Einheit
  weicht ab, ein Attribut fehlt): der Beleg ist dafür da; Folge ist ein Fix-Slice
  (§4), kein Weiterschreiben. — **Ausgang:** offen bis Closure.
- **Der Post-Push-Lauf von `e2e.yml` bleibt unbestätigt** (§3.10 von
  [`AGENTS.md`](../../../../AGENTS.md), Klasse
  [`BEO-PGC/github-actions-unverifizierbar-lokal`](../observations/BEO-PGC/github-actions-unverifizierbar-lokal/observation.md)):
  ein lokaler grüner Lauf sagt nichts über den Runner mit Docker-Hub-Zugriff und
  Zeitlimit. — **Ausgang:** offen bis Closure (nachzutragen mit dem Lauf).
- **Das Handbuch zieht nicht mit** (die `https`- und Einheits-Aussage):
  [`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`](../observations/BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/observation.md)
  — eigener Liefer-Punkt 3 (c), bedingt. — **Ausgang:** offen bis Closure.
- **Der Slice ist größer als drei Punkte tragen:** Rückführung nach §4 mit dem
  benannten Teilungsvorschlag. — **Ausgang:** offen bis Closure.

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
- **Validator-Feststellung (Modul 8):** (bei der Closure zu füllen; der Slice
  liefert Betreiber-Wert, ein Validator-Lauf ist nach dem Release sinnvoll)
- **Folge-Slices:** (bei der Closure zu füllen)
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
berührten Pfade (`tools/harness/`, `docs/user/`, `harness/`) liegen alle in ihr, es
entsteht keine neue Sub-Area — die Schwelle (≥ 2 von 3 Achsen) wird nicht
angewandt, weil kein Pfad ausdifferenziert werden muss.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/BEO-PGC/`) am Planungsstand durchgegangen;
die Zähler sind die Zahl der `evidence/`-Dateien und beim Start neu zu zählen
(`ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence | wc -l`); am
Planungsstand **nicht einzeln gezählt**, deshalb ohne Zahl genannt (Instanz B von
[`AGENTS.md`](../../../../AGENTS.md) §3.12): Treffer für die berührte Sub-Area sind
[`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt`](../observations/BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt/observation.md)
(Gegenstand dieses Slice: Metrik-Boundary im E2E belegt, nicht vom Reviewer),
[`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`](../observations/BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke/observation.md)
(Gegenmaßnahme: der Container-Start der Runner-Phase),
[`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`](../observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/observation.md)
(Gegenprobe je Negativfall),
[`BEO-PGC/runner-zaehler-ohne-gefahrene-mutation`](../observations/BEO-PGC/runner-zaehler-ohne-gefahrene-mutation/observation.md)
(die Mutationsproben der Runner-Phase),
[`BEO-PGC/ready-ist-nicht-verbunden`](../observations/BEO-PGC/ready-ist-nicht-verbunden/observation.md)
(erst der Beleg, dass der Empfänger steht, dann die Aussage über seine Abwesenheit),
[`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md)
(der neue Pin steht an einer Stelle und wird von P10 gelesen),
[`BEO-PGC/github-actions-unverifizierbar-lokal`](../observations/BEO-PGC/github-actions-unverifizierbar-lokal/observation.md)
(Post-Push-Lauf),
[`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`](../observations/BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/observation.md)
und
[`BEO-PGC/handbuch-versionshistorie-uebersprungen`](../observations/BEO-PGC/handbuch-versionshistorie-uebersprungen/observation.md)
(Handbuch-Zug, bedingt). Der Planner liest die Zähler **nicht** hier ab; der
Implementer zählt beim Start und trägt sie in den Bericht.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
