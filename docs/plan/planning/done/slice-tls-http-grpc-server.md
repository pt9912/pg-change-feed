# Slice tls-http-grpc-server: ein gemeinsames TLS-Paar für die HTTP- (einschließlich SSE) und die gRPC-Schnittstelle

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
in `done/` liegt — **nicht** aus fachlichem Grund (TLS und Token-Wechsel sind
unabhängig), sondern wegen des Konfigurationsgleichlaufs: beide Slices ändern
`internal/bootstrap/wiring.go` (`Config`, `ConfigFromEnv`),
`internal/bootstrap/config_file.go` (`mergeConfig`, `fileConfig`), die
Code-Tabelle der Meldungscodes und den Konfigurationsteil des Benutzerhandbuchs;
die Nummern der neuen Codes werden in der Reihenfolge der Closure vergeben
(Zuweisungsregel im Plan des Vorgängers, §3).

**Bezug:** [`LH-FA-SST-011`](../../../../spec/lastenheft.md) (Scope),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) und
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (die Fähigkeiten der
Schnittstellen sind unter TLS fachlich gleichwertig nutzbar),
[`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) (Festlegung 1 und 5),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
(Meldungscodes der Klasse `configuration`),
[`ADR-0088`](../../adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
(Diskriminator: ein Dateipfad ist keine Zugangsdaten),
[`ADR-0057`](../../adr/0057-http-grpc-api.md) und
[`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) (die zwei Server),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).

**Berührte Spec-Stellen:** [`SPEC-034`](../../../../spec/pflichtenheft.md) ·
[`SPEC-016`](../../../../spec/pflichtenheft.md) (die Datei-Schlüssel
`tls_cert_file` und `tls_key_file`; sie gehören **nicht** zur
Zugangsdaten-Klasse) ·
[`SPEC-018`](../../../../spec/pflichtenheft.md) ·
[`SPEC-020`](../../../../spec/pflichtenheft.md) ·
[`SPEC-021`](../../../../spec/pflichtenheft.md) ·
[`SPEC-031`](../../../../spec/pflichtenheft.md) ·
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Klasse `configuration`) ·
[`ARC-005`](../../../../spec/architecture.md) (Driving Adapter) ·
[`ARC-007`](../../../../spec/architecture.md) (Bootstrap).
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
liegt mit [`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) vor, eine
Schichten-Frage entsteht nicht — TLS ist Sache der Driving Adapter und des
Bootstraps). **Datum:** 2026-10-04.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass:** Der HTTP-Server (einschließlich des SSE-Endpunkts) und der
gRPC-Server laufen im Klartext
([`LH-FA-SST-011`](../../../../spec/lastenheft.md)). Gemessen am Parent: der
HTTP-Adapter startet mit `ListenAndServe` (`internal/adapters/driving/http/server.go`
Zeile 132), der gRPC-Adapter baut `grpc.NewServer` ohne Transport-Credentials
(`internal/adapters/driving/grpc/server.go` Zeile 91) — §3 Suchlauf Zeile 1
zählt 8 Trefferzeilen für die drei Symbole am Parent.

**Ziel:** Sind `CDC_TLS_CERT_FILE` und `CDC_TLS_KEY_FILE` (oder die
Datei-Schlüssel `tls_cert_file` und `tls_key_file`) gesetzt, bedient der
Feed-Container HTTP (einschließlich `GET /changes/stream`) und gRPC
ausschließlich über TLS mit Mindestversion 1.2 — ohne Klartext auf derselben
Adresse, und mit Start-Abbruch der Klasse `configuration` bei unvollständigem,
nicht ladbarem oder nicht zusammenpassendem Paar —; ohne die Konfiguration
bleibt alles unverändert; beides ist am laufenden Container belegt und im
Benutzerhandbuch beschrieben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **TLS-Optionen der SDK-Packages** (C#, Kotlin, Python) — der eigene Slice
  `sdk-tls-optionen`; er wartet auf die in
  [`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) §Konsequenzen verlangte
  Folge-Anforderung (das Lastenheft nennt die SDK-Optionen als „eigene,
  nachgelagerte Anforderung“) und auf diesen Slice als Server-Beleg. Er ist
  **noch keine Datei**: eine Adresse hat er erst mit der Anforderung, bis dahin
  trägt der Plan ihn nur in der Abhängigkeitsübersicht.
- **Die Beispielprogramme unter `examples/`** (Go, C#, Kotlin) — sie brechen
  nicht am Server-Vertrag, weil der Server ohne TLS-Konfiguration unverändert
  läuft; gegen einen TLS-Server sind sie nicht Gegenstand. Gemessen am Parent:
  das Go-gRPC-Beispiel verbindet mit `insecure.NewCredentials()`
  (`examples/grpc-client/main.go` Zeile 102; §3 Suchlauf Zeile 2 zählt 9
  Trefferzeilen im Baum). Das Handbuch nennt nur, was am Arbeitsstand gemessen
  ist (§2 Liefer-Punkt 3).
- **Client-Zertifikate (mTLS), Neuladen von Zertifikat und Schlüssel,
  Beschaffung und Erneuerung** — Out-of-Scope von
  [`LH-FA-SST-011`](../../../../spec/lastenheft.md); ein Zertifikatswechsel ist
  ein Neustart.
- **TLS zu NATS und zu PostgreSQL** — Out-of-Scope derselben Anforderung (NATS-
  Konfiguration, DSN-Parameter).
- **Ein Umbau von `--healthcheck`** — gemessen am Parent: der Aufruf
  (`cmd/pg-change-feed/main.go`, `bootstrap.Healthcheck(ctx, cfg.ReaderDSN,
  cfg.Source)`) liest das Alter des letzten Lebenszeichens über die
  Reader-Verbindung zur Datenbank und erreicht den HTTP- oder gRPC-Server
  **nicht**; der Compose-Healthcheck ruft dasselbe Binary
  (`compose.yaml`, `examples/compose.yaml`). TLS ändert ihn deshalb nicht; der
  Slice belegt das (§2 Liefer-Punkt 2), er baut nichts um. Läuft `--healthcheck`
  in einem Zug über `ConfigFromEnvAndFile`, liest er die neue Konfiguration mit:
  das Paar wird dort nur auf Vollständigkeit geprüft, nicht geladen (§3).
- **Eine Ablauf- oder Namensprüfung des Zertifikats beim Start** — die Spec
  nennt „nicht lesbar, ungültig oder passt nicht zusammen“; ein abgelaufenes
  Zertifikat lehnt der Client ab, der Server lädt es (benannte Grenze, §6).

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

- [x] **TLS in beiden Servern, Konfiguration, Codes (Liefer-Punkt 1,
      [`LH-FA-SST-011`](../../../../spec/lastenheft.md) Happy Path, Boundary,
      Negative).**
      (a) Konfiguration: `CDC_TLS_CERT_FILE`/`CDC_TLS_KEY_FILE` und die
      Datei-Schlüssel `tls_cert_file`/`tls_key_file` (env schlägt Datei, Feld für
      Feld wie `http_addr`); die zwei Schlüssel stehen **nicht** in
      `forbiddenFileCredentialKeys` (die Klasse bleibt bei der Zahl, die der
      Vorgänger-Slice gesetzt hat; ein Test lädt eine Datei mit beiden Schlüsseln
      ohne Fehler — die Gegenprobe zur Klasse). Nach dem Merge gilt: keines
      gesetzt → unverschlüsselt wie bisher; **nur eines** gesetzt (auch: Zertifikat
      aus der Datei, Schlüssel fehlt in beiden Quellen) → `ErrConfiguration` mit neuem
      Meldungscode (nächste freie Nummer ab `PCF-E2008`, Zuweisungsregel im Plan von
      [`slice-api-token-mehrfach-konfiguration`](../done/slice-api-token-mehrfach-konfiguration.md)
      §3). Beide Zugriffswege (`ConfigFromEnv`, `mergeConfig`) tragen die
      Prüfung aus **einer** Funktion.
      (b) Laden: ein Pfad nicht lesbar, keine gültigen PEM-Daten oder Zertifikat
      und Schlüssel bilden kein Paar → `ErrConfiguration` mit neuem Code, **vor** dem
      ersten Listener und vor dem Anlegen des Replication-Slots (der Ort des Ladens
      in `internal/bootstrap` ist so gewählt, dass bei einem Ladefehler weder eine
      Datenbankverbindung entstanden noch ein Slot angelegt ist; zu belegen durch
      den Verdrahtungs-Test und die Runner-Phase, nicht durch eine Lesung; zu einem
      geöffneten Port macht der Plan nach der Fixrunde keine Aussage, siehe den
      Fixrunden-Beleg in diesem Abschnitt). `--healthcheck` lädt die
      Dateien **nicht**.
      (c) Betrieb: HTTP (`ListenAndServeTLS`-Weg, einschließlich SSE) und gRPC
      (Server-Credentials) bedienen mit einer `tls.Config` mit explizit gesetzter
      `MinVersion` TLS 1.2, ohne Client-Zertifikat-Pflicht; ein Bind-Fehler bleibt
      ein Log-Eintrag. Unit-Tests (`make test`, netzlos, Loopback) mit einem im Test
      erzeugten Zertifikat (`crypto/x509`, keine Datei im Repo, kein Host-Werkzeug):
      ein TLS-Client mit dem Zertifikat als Vertrauensanker bekommt jede Fähigkeit
      (HTTP: ein Endpunkt je Rechtsklasse und ein SSE-Event; gRPC: ein Stream-Event
      und ein Verwaltungs-RPC); ein Klartext-Client bekommt **keine** Antwort der
      API (Statuscode `2xx` oder gRPC-Antwort ausgeschlossen — das konkrete
      Verhalten von `net/http` und `grpc-go` gegenüber einem Klartext-Client ist
      *nach Kenntnisstand, ungeprüft*; der Test misst es und der Bericht nennt es);
      ein Client, der höchstens TLS 1.1 spricht, scheitert (ist die Bibliothek dazu
      nicht fähig, prüft der Test stattdessen die `MinVersion` der gebauten
      Konfiguration unmittelbar — der Bericht nennt den gewählten Weg).
      (d) **Mutationsprobe (Zusage, an einer Kopie im Scratchpad, nicht an der
      Datei im Arbeitsbaum — [`AGENTS.md`](../../../../AGENTS.md) §3.1; Instanz:
      der Go-Test des jeweiligen Pakets im Toolchain-Container mit der Kopie als
      Mount):** vier Stellen einzeln mutiert, je die rote Farbe im Bericht — die
      Paar-Bedingung („nur eines gesetzt“, `||` zu `&&`), das Verschlucken des
      Ladefehlers (Server startet danach im Klartext), die `MinVersion` (entfernt
      oder auf 1.0 gesetzt) und der Klartext-Pfad im HTTP-Server (TLS gesetzt, Start
      trotzdem ohne `tls.Config`). Die Verallgemeinerung auf den gRPC-Server ist
      **hergeleitet**, wo die Mutation nur am HTTP-Server gefahren ist; der Bericht
      nennt je Mutation Stelle und Instanz.
- [x] **Realserver-Beleg am laufenden Container (Liefer-Punkt 2,
      [`LH-FA-SST-011`](../../../../spec/lastenheft.md) Happy Path, Boundary,
      Negative).** Eine neue Phase von `tools/harness/run-integration-tests.sh`
      (`make test-integration`), über eine Compose-Override-Datei im Temp-Verzeichnis
      (Muster: die Phasen zur Leerlauf-Bestätigung; `compose.yaml` bleibt
      unverändert; der Container ist am Phasen-Ende ohne Override wiederhergestellt):
      (1) Ein Hilfsprogramm unter `tools/harness/` (Arbeitsname `certgen`, nur
      Standardbibliothek, `tooling`-Gruppe von `.a-check.yml`) erzeugt Zertifikat
      und Schlüssel im Toolchain-Container (Docker-only, kein Host-Werkzeug über
      die Basis hinaus; Alternativnamen und Dateirechte so, dass der Prozess des
      Containers beide liest); (2) Happy Path: die Wegwerf-Clients
      (`tools/harness/httpclient`, `sseclient`, `grpcclient`, `grpcadminclient`)
      verbinden über TLS mit dem Zertifikat als Vertrauensanker und erreichen
      `GET /tables`, ein SSE-Event einer danach committeten Änderung (Tabelle,
      Operation, Wert, `change_id` gegen `cdc.changes` gehalten), einen gRPC-Stream-
      Event und einen Verwaltungs-RPC; (3) **Boundary — kein stiller Rückfall:** erst
      nach einem belegten erfolgreichen TLS-Aufruf versucht der Runner Klartext auf
      dieselbe HTTP- und dieselbe gRPC-Adresse (ein Negativ über Abwesenheit
      beobachtet zuerst, dass der Server steht,
      [`BEO-PGC/ready-ist-nicht-verbunden`](../observations/BEO-PGC/ready-ist-nicht-verbunden/observation.md)),
      und beide werden nicht bedient (keine Antwort der API, auch nicht mit gültigem
      Token); (4) `--healthcheck` im Container endet unter TLS mit Exit 0 und der
      Docker-Gesundheitszustand ist `healthy` (die Aussage „TLS ändert den
      Healthcheck nicht“ ist damit **gemessen**, nicht nur gelesen); (5) **Negative:**
      drei Starts mit je einer Verletzung — nur ein Pfad gesetzt, ein nicht lesbarer
      Pfad, ein nicht zusammenpassendes Paar — beenden den Container mit der
      Fehlerklasse `configuration` und dem jeweiligen Meldungscode in der Log-Zeile
      (per `docker logs` gelesen), und es entsteht kein Slot (die Prüfung der
      HTTP-Adresse nach einem beendeten Container entfällt, siehe den
      Fixrunden-Beleg in diesem Abschnitt) — die
      Gegenprobe ist der Start mit dem
      gültigen Paar aus (2); (6) ohne die Konfiguration laufen **alle** übrigen
      Phasen unverändert im Klartext (sie sind der Boundary-Beleg „unverändert wie
      ohne diese Anforderung“). Die Phase trägt eine `abdeckung_declare`-Zeile für
      [`LH-FA-SST-011`](../../../../spec/lastenheft.md);
      `docs/user/e2e-abdeckung.md` ist vom Runner neu geschrieben (Erzeugnis,
      committet); `make doc-trace` druckt die Zeile mit der Waisen-Zahl, und
      `LH-FA-SST-011` ist darin **keine** Waise (gedruckte Zeile im Bericht).
- [x] **Benutzerhandbuch und benannte Grenzen zu den Clients (Liefer-Punkt 3).**
      `docs/user/benutzerhandbuch.md`: Umgebungsvariablen-Tabelle (§5) um die zwei
      Variablen, die Datei-Felder `tls_cert_file`/`tls_key_file` samt
      YAML-Beispiel (und der Satz, dass es Pfade sind, keine Zugangsdaten — der
      private Schlüssel selbst steht in keiner Umgebungsvariable und keiner
      Konfigurationsdatei), ein Abschnitt zur Verschlüsselung der HTTP- und
      gRPC-Schnittstelle (ein Paar für beide, auch für den SSE-Endpunkt;
      Mindestversion TLS 1.2; kein Klartext auf derselben Adresse; kein
      Client-Zertifikat; Zertifikatswechsel nur mit Neustart; NATS und PostgreSQL
      unberührt; die Start-Fehler mit Klasse `configuration` und den neuen Codes),
      die Katalog-Zeilen der neuen Codes, und der Satz zu den Clients: er nennt
      **nur das am Arbeitsstand Gemessene** — für die Beispielprogramme unter
      `examples/` und die Client-Packages ein Aufruf bzw. `git grep` gegen einen
      TLS-Server (Befehl und Ergebnis im Bericht), nicht eine Erwartung;
      Änderungshistorie-Zeile mit der nächsten freien Nummer
      ([`BEO-PGC/handbuch-versionshistorie-uebersprungen`](../observations/BEO-PGC/handbuch-versionshistorie-uebersprungen/observation.md)).
      Kennungsfrei: `make handbuch-public-doc-check` Exit 0; die neuen
      Meldungstexte im Go-Quelltext tragen keine Kennung
      (`make ausgabe-kennungen-check` Exit 0, `make meldungscodes-check` Exit 0).

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [x] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0
      (Kennungen in diesem Plan verlinkt), `make test`, `make test-store`
      und `make test-integration` Exit 0, `make image` Exit 0 (der Zug ändert
      Build-Kontext-Dateien), `make a-check` Exit 0.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      ([`review-slice-tls-http-grpc-server`](../../../reviews/review-slice-tls-http-grpc-server.md),
      F-1 und F-2 MEDIUM in der Fixrunde behoben;
      `.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors
      (Zeile `make test-integration`: ein Satz zum TLS-Rundlauf); gemeldete
      Träger fremder Dateien mit der Closure nachgezogen (§3 Suchlauf,
      [`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7: benannte Spec-Lücke).

**Beleg des Implementers** (Lauf am Arbeitsstand dieses Zuges; jede Angabe
ist **gemessen** am genannten Befehl, die Verallgemeinerungen tragen das Wort
*hergeleitet*; Review und Verifikation stehen aus):

- **Liefer-Punkt 1.** `make test` Exit 0 (52 Pakete `ok`, darunter
  `internal/adapters/driving/http`, `…/grpc`, `internal/bootstrap`,
  `internal/domain/messagecode`); `make test-store` Exit 0, gedruckt
  `db-coverage: OK — DB-Adapter-Coverage 83.03% erfuellt Schwelle 80%`;
  `make gates` Exit 0, gedruckt `coverage-gate: OK — Coverage 83.10% erfüllt
  Schwelle 80%`, `d-check: 1662 Datei(en) geprüft, 0 Befund(e)`,
  `meldungscodes-check: 100 Codes in Tabelle und Katalog gleich`, `gesamt: 0
  Befund(e)` (a-check); `make a-check` Exit 0 (`.a-check.yml` unverändert),
  `make fmt-check` Exit 0 (`342 Go-Dateien geprüft, alle formatiert`).
  Codes: `PCF-E2009` (unvollständiges Paar), `PCF-E2010` (nicht ladbar).
  Gemessenes Verhalten von `net/http` und `grpc-go` gegenüber einem
  Klartext-Client (Unit-Tests, Loopback): `net/http` antwortet auf einen
  `http://`-Aufruf mit Status `400` und dem Text `Client sent an HTTP request to
  an HTTPS server.` (eine Antwort des Servers, keine der API); `grpc-go` beendet
  den Verbindungsaufbau, der RPC endet mit `Unavailable` (`error reading server
  preface: connection reset by peer`; in der Runner-Phase `…: EOF`). Ein
  TLS-1.1-Client scheitert mit `remote error: tls: protocol version not
  supported` (am `httptest`-Server mit der gebauten Konfiguration; die zwei
  Adapter tragen die Konfiguration unverändert, ein zweiter Aufbau am Draht
  steht nicht aus — *hergeleitet*).
- **Mutationen** (an einer Kopie im Scratchpad, Instanz: `go test -race` des
  Pakets im Toolchain-Container mit der Kopie als Mount; je Mutation einzeln,
  Datei danach per `cp` zurückgesetzt):

  | Zusage | Stelle · mutierte Eingabe | Instanz | gesehenes Rot |
  |---|---|---|---|
  | nur Zertifikat gesetzt endet mit `PCF-E2009` | `validateTLSPair`, Fall 1 `&&` zu `\|\|` | `internal/bootstrap` | rot: `TestTLSPaarWirdGelesen`, `TestNewTLSConfigBautEineKonfiguration`, `TestNewTLSConfigMindestversionAmDraht`, `TestNewTLSConfigLadefehlerEndenMitConfiguration`, `TestNewTLSConfigPruftAblaufUndNamenNicht`, `TestRunLaedtTLSPaarVorJederVerbindung` |
  | nur Schlüssel gesetzt endet mit `PCF-E2009` | `validateTLSPair`, Fall 2 zu `case false` | `internal/bootstrap` | rot: `TestTLSPaarUnvollstaendigEndetMitConfiguration` |
  | Ladefehler beendet den Start vor jeder Verbindung | `Run`: Fehler von `newTLSConfig` verworfen (`tlsConfig = nil`) | `internal/bootstrap` | rot: `TestRunLaedtTLSPaarVorJederVerbindung` |
  | Schlüsselpfad wird gelesen | `newTLSConfig`: `LoadX509KeyPair(certFile, certFile)` | `internal/bootstrap` | rot: `TestNewTLSConfigBautEineKonfiguration`, `…MindestversionAmDraht`, `…PruftAblaufUndNamenNicht`, `TestRunLaedtTLSPaarVorJederVerbindung` |
  | Mindestversion TLS 1.2 explizit | `newTLSConfig`: `MinVersion` entfernt | `internal/bootstrap` | rot: nur `TestNewTLSConfigBautEineKonfiguration` (die Feldprüfung); der Draht-Test bleibt grün, weil `net/http` ohne `MinVersion` selbst TLS 1.2 verlangt — die Zusage „explizit“ hängt allein an der Feldprüfung |
  | Mindestversion TLS 1.2 am Draht | `newTLSConfig`: `MinVersion` auf `tls.VersionTLS10` | `internal/bootstrap` | rot: `TestNewTLSConfigBautEineKonfiguration`, `TestNewTLSConfigMindestversionAmDraht` |
  | Datei-Feld `tls_key_file` wirkt | `mergeConfig`: `cfg.TLSKeyFile` ohne Datei-Wert | `internal/bootstrap` | rot: `TestTLSPaarWirdGelesen`, `TestTLSPaarUnvollstaendigEndetMitConfiguration` |
  | die zwei Datei-Schlüssel gehören nicht zur Zugangsdaten-Klasse | `forbiddenFileCredentialKeys` um `tls_key_file` ergänzt | `internal/bootstrap` | rot: `TestZugangsdatenKlasseCodeUndTestSindMengengleich`, `TestTLSDateiSchluesselGehoerenNichtZurZugangsdatenKlasse`, beide Paar-Tests |
  | HTTP-Server bedient mit TLS-Konfiguration nur über TLS | `http.Server.serve`: Verzweigung `TLSConfig != nil` zu `false` (Klartext trotz Konfiguration) | `internal/adapters/driving/http` | rot: `TestServerTLSBedientJedeFaehigkeit`, `TestServerTLSLehntKlartextClientAb` (beide am Handshake der TLS-Gegenprobe) |
  | HTTP-Server ohne Konfiguration im Klartext | `http.Server.serve`: Verzweigung zu `true` | `internal/adapters/driving/http` | rot: `TestServerOhneTLSKonfigurationBedientKlartext` |
  | gRPC-Server bedient mit TLS-Konfiguration nur über TLS | `grpc.New`: Bedingung `cfg.TLSConfig != nil` zu `false` | `internal/adapters/driving/grpc` | rot: `TestServerTLSBedientStreamUndVerwaltung`, `TestServerTLSLehntKlartextClientAb` |

  Die Mutation „Klartext-Pfad im HTTP-Server“ und die des gRPC-Servers sind
  **je am eigenen Server** gefahren (keine Herleitung von HTTP auf gRPC).
  Der Zweig „Klartext-Client bekommt eine Antwort“ von
  `TestServerTLSLehntKlartextClientAb` färbt sich in beiden Paketen erst über
  die TLS-Gegenprobe davor rot; die Klartext-Aussage allein (Status `2xx`) ist
  an der Eingabeseite nicht einzeln mutiert — der Gegenfall liegt in
  `TestServerOhneTLSKonfigurationBedientKlartext` (Klartext bedient, wenn keine
  Konfiguration gesetzt ist).
- **Liefer-Punkt 2.** `make image` Exit 0 (`ghcr.io/pt9912/pg-change-feed:dev`);
  `make test-integration` Exit 0, gedruckt: `run-integration-tests: TLS der
  HTTP-, SSE- und gRPC-Schnittstelle (LH-FA-SST-011) belegt — fünf Starts mit
  einer Verletzung (nur Zertifikat, nur Schlüssel, nicht lesbarer Pfad, Datei
  ohne PEM-Daten, nicht zusammenpassendes Paar) beendeten den Container mit
  Fehlerklasse configuration und PCF-E2009 bzw. PCF-E2010 im Log, ohne Slot;
  … unter TLS (Umgebung) blieb --healthcheck gesund
  (Ausgang 0, Zustand healthy), HTTP, gRPC, SSE (change_id=2365-1) und
  gRPC-Stream (change_id=2369-1) wurden über TLS bedient …, der Klartext-Versuch
  danach bekam keine Antwort der API (HTTP: `CLEARTEXT status=400
  body="Client sent an HTTP request to an HTTPS server."`; gRPC: `CLEARTEXT
  code=Unavailable message="connection error: desc = \"error reading server
  preface: EOF\""`) …`. Der Ausgang der Starts ist 2 bei `PCF-E2009` (die
  Konfiguration endet vor `Run`) und 1 bei `PCF-E2010` (`Run` lädt das Paar);
  der Slot `slot_tls_neg` ist nach jedem der fünf Negativstarts abwesend, und
  der Start mit gültigem Paar über die Datei-Schlüssel legt ihn an (Gegenprobe;
  der Hilfs-Slot wird danach entfernt). Die Wiederherstellung ohne Override
  bedient beide Adressen wieder im Klartext (`probe` auf `http://` und über
  gRPC ohne Vertrauensanker). `git status --porcelain` nennt am Phasen-Ende
  keine `.pem`-Datei. `make doc-trace`: `83 Anforderung(en), 1 Waise(n).`,
  Zeile `LH-FA-SST-011 | … | E2E | ok`; die verbleibende Waise ist
  `LH-FA-SST-010`. Die Zahl am Parent ist nicht nachgemessen.
- **Liefer-Punkt 3.** Handbuch 1.96 (Umgebungsvariablen-Tabelle, Datei-Felder
  samt YAML-Beispiel, Abschnitt „Schnittstellen mit TLS verschlüsseln“,
  Katalogzeilen `PCF-E2009`/`PCF-E2010`, Fehlersuche, Änderungshistorie);
  `make handbuch-public-doc-check` Exit 0 (`keine interne Kennung in 3
  Nutzerdokumenten unter docs/user/`), `make ausgabe-kennungen-check` Exit 0,
  `make meldungscodes-check` Exit 0. Die Aussage zu den Clients ist je Sprache
  gemessen: `git grep -n -i -E 'tls|ssl|certificate'` über
  `sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs`,
  `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/PgChangeFeedClientOptions.kt`
  und `sdks/python/pgchangefeed/src/pgchangefeed/options.py` liefert Exit 1
  (kein Treffer); `git grep -c -E 'insecure\.NewCredentials|usePlaintext|plaintext|Plaintext|Klartext'`
  über `examples/grpc-client`, `examples/csharp/grpc-client`,
  `examples/kotlin/grpc-client` trifft je Programm (Go 1, C# 2, Kotlin 4 in
  `Main.kt` und 1 im Test). Die HTTP- und SSE-Beispiele tragen keinen
  TLS-spezifischen Code (`git grep -i -E 'tls|https://|ssl'` in
  `examples/http-client`, `examples/sse-client`: 0 Treffer); das Handbuch nennt
  für sie nichts, weil ein Aufruf gegen einen TLS-Server nicht gefahren ist.
  Gemessen sind nur die zwei Befunde der `git grep`-Läufe (Optionen ohne
  TLS-Einstellung, gRPC-Beispiele fest im Klartext); sie sind Entwicklersicht
  und stehen im Handbuch nicht als Beleg. Der Handbuch-Absatz sagt nach der
  Fixrunde nur, dass die gRPC-Beispiele nicht mit einem TLS-Server verbinden,
  und sagt für die Client-Pakete und für HTTP-/SSE-Clients nichts zu.
- **Kommentar-Läufe (Schritt 20).** `make kommentar-kennungen DIFF=e1db9c30`
  Exit 0 (kein Kandidat); der Chronik-Lauf trifft nichts (`xargs`/`grep` Exit
  123 ohne Treffer); der Konjunktiv-Lauf trifft fünf Zeilen, jede geprüft:
  zwei Mutationsbeschreibungen in Test-Godocs (`statt`), eine
  Test-Godoc-Bindung („wäre die Reihenfolge anders, endete `Run` …“) und zwei
  normale Zweige („sonst im Klartext“) — zulässig.
- **Benannte Grenzen.** (1) Ablauf und Name des Zertifikats prüft der Server
  beim Start nicht: `TestNewTLSConfigPruftAblaufUndNamenNicht` lädt ein
  abgelaufenes Zertifikat ohne Fehler (Spec-Lücke zu „ungültig“, im Handbuch
  benannt). (2) Ein Fehler beim Laden steht nicht im Heartbeat-Fehlerzustand:
  `Run` endet vor dem Aufbau des Heartbeat-Adapters (wie jeder Fehler vor
  diesem Punkt); der Betreiber liest ihn im Log. (3) Ein `Addr` ohne Wert
  bindet über `net.Listen` einen zufälligen Port statt `:http`/`:https`; die
  Verdrahtung ruft `Start` nur mit gesetzter Adresse. (4) Die Datei mit
  Zertifikat und Kette ist nur an `crypto/tls` gemessen
  (`TestNewTLSConfigLiestDieKetteDerZertifikatsdatei`: zwei Glieder), nicht an
  einem Client mit Zwischenzertifikat.
- **Fixrunde (Review `review-slice-tls-http-grpc-server`).** F-1: der
  Clients-Absatz des Handbuchs und die Historie-Zeile 1.96 sagen nur noch, was
  gemessen ist (gRPC-Beispiele im Klartext, Optionen der Pakete ohne
  TLS-Einstellung), und treffen für Client-Pakete und HTTP-/SSE-Clients keine
  Aussage; ein Lauf der SDK-HTTP-Clients gegen einen TLS-Server ist nicht
  gefahren (*offen*, keine Behauptung). F-2: die Prüfung „kein Klartext-Server“
  in `tl_expect_start_refused` las einen bereits beendeten Container und konnte
  an keiner Eingabe rot werden; sie entfällt, und mit ihr die Aussage „ohne
  geöffneten Port“ in Plan, Handbuch, `harness/README.md` und Runner-Zeile.
  Getragen bleibt: Ausgang, Fehlerklasse, Code und Text im Log, kein
  Schlüsselmaterial, kein Slot; in `Run` entsteht vor dem Laden keine
  Datenbankverbindung. Die Aussage „vor der Datenbankverbindung“ steht im Code-
  Kommentar von `Run` (angepasst). Mutation (Kopie im Scratchpad, Instanz
  `go test -race -run TestRunLaedtTLSPaarVorJederVerbindung ./internal/bootstrap`
  im Toolchain-Container, **gemessen**): `Run` lädt das Paar erst nach
  `postgresstorage.New` → rot: `tls_internal_test.go:421: Run mit nicht ladbarem
  Paar: Fehlerklasse storage [PCF-E5001] …, erwartet ErrTLSPairUnusable ohne
  Datenbankfehler`. Der Slot-Zähler des Runners färbt sich bei *dieser*
  Mutation nicht rot (der Slot entsteht erst in `NewStream`); er ist nur an
  einer Mutation gebunden, die das Laden hinter `NewStream` zieht — diese ist am
  Runner nicht gefahren, *hergeleitet*. F-3: Zahl der Starts in den
  Runner-Kommentaren auf fünf. F-4: `chmod 0644` auf die `*.pem` bleibt, weil
  der Container als `nonroot` läuft und der Besitzer der Dateien der
  Host-Benutzer ist (der Kommentar im Runner nennt den Grund; Schlüssel sind
  Wegwerf-Material). F-5: die Kommentare der C#- und Kotlin-gRPC-Beispiele
  nennen den Server als TLS-fähig und das Beispiel als Klartext-Client.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt (gemessen: `docs/plan/planning/` trägt keine flache Welle-Datei) und „die nächste Welle-Closure“ damit keine Adresse ist.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/http/server.go` | update | `Config` trägt eine optionale `*tls.Config`; `Start` wählt den TLS-Weg, wenn sie gesetzt ist, sonst den bisherigen; ein Hilfsweg mit fertigem Listener (Muster `serve(listener)` im gRPC-Adapter) macht den Server im Test auf Port 0 erreichbar |
| `internal/adapters/driving/grpc/server.go` | update | `Config` trägt dieselbe optionale `*tls.Config`; `grpc.NewServer` bekommt Transport-Credentials nur, wenn sie gesetzt ist |
| `internal/adapters/driving/http/*_test.go`, `internal/adapters/driving/grpc/*_test.go` | neu / update | Tests aus §2 (c): TLS-Client je Fähigkeit, Klartext-Client, `MinVersion`; ein im Test erzeugtes Zertifikat über einen kleinen gemeinsamen Test-Helfer je Paket (keine Datei im Repo) |
| `internal/bootstrap/wiring.go` | update | `Config`-Felder, Konstanten der zwei Variablen, `ConfigFromEnv`; Laden des Paars und Bau der `tls.Config` **vor** `NewStream` und vor den Servern; Übergabe an beide Adapter |
| `internal/bootstrap/config_file.go` | update | `fileConfig` (`tls_cert_file`, `tls_key_file`), `mergeConfig`; **kein** Eintrag in `forbiddenFileCredentialKeys` |
| `internal/bootstrap/config_file_internal_test.go`, `internal/bootstrap/wiring_test.go` | update | Paar-Fälle für beide Zugriffswege (keines, eines je Quelle, beide, gemischt Datei/Env), Ladefehler (nicht lesbar, kein PEM, nicht zusammenpassend) mit Gegenprobe des gültigen Paars, Beleg „kein Slot und kein Listener bei Ladefehler“ |
| `internal/bootstrap/tls.go` | neu (Implementer-Nachzug) | trägt `validateTLSPair` (die eine Funktion, die beide Konfigurations-Zugriffswege rufen), `newTLSConfig` (Laden und Bau der einen `tls.Config`, `MinVersion` TLS 1.2 explizit) und die zwei Sentinels `ErrTLSPairIncomplete`/`ErrTLSPairUnusable` samt `tlsConfigError` (Form wie `apiTokenListError`, ist zugleich `ErrConfiguration`); eigene Datei statt Einbau in `wiring.go` (2353 Zeilen) |
| `internal/bootstrap/tls_internal_test.go`, `internal/bootstrap/messagecodes_internal_test.go`, `internal/adapters/driving/http/tls_test.go`, `internal/adapters/driving/grpc/tls_test.go` | neu / update (Implementer-Nachzug) | die Tests aus §2 (a) bis (c): Konfigurationswege, Ladefehler, `MinVersion`, Ablauf-Grenze, Reihenfolge in `Run`, TLS-Client je Fähigkeit und Klartext-Client je Server; die zwei Sentinels stehen in der Sentinel-Tabelle des Meldungscode-Tests |
| `internal/domain/messagecode/codes.go` (und der Test des Pakets, falls er die Menge zählt) | update | neue Codes der Klasse `configuration` (unvollständiges Paar `PCF-E2009`; Laden fehlgeschlagen `PCF-E2010`), Zuweisungsregel im Plan des Vorgängers; der Test des Pakets zählt die Menge nicht |
| `tools/harness/certgen/` (Arbeitsname) | neu | erzeugt Zertifikat und Schlüssel (Standardbibliothek), läuft im Toolchain-Container; `tooling`-Gruppe, keine Kante nach `internal/` |
| `tools/harness/httpclient`, `sseclient`, `grpcclient`, `grpcadminclient` | update | TLS-Weg mit Vertrauensanker aus einer Datei (Umgebungsvariable `HARNESS_TLS_CA_FILE`; ohne sie bleibt die Aufrufform der bestehenden Phasen unverändert); dazu je ein Modus `cleartext` in `httpclient` und `grpcadminclient`, der einen Klartext-Aufruf mit gültigem Token sendet und die Antwort meldet |
| `tools/harness/run-integration-tests.sh` | update | neue Phase „TLS“ (fünf Negativstarts statt drei: nur Zertifikat, nur Schlüssel, nicht lesbarer Pfad, Datei ohne PEM-Daten, nicht zusammenpassendes Paar; dazu ein Start über die Datei-Schlüssel mit frischem Slot als Gegenprobe zur Abwesenheit des Slots), `abdeckung_declare` für [`LH-FA-SST-011`](../../../../spec/lastenheft.md), Aufräumen des Temp-Verzeichnisses im `cleanup`-Trap; jede neue `func TestE2E*` stünde im `-run` (Vollständigkeits-Test `test/integration/runner_vollstaendigkeit_test.go`) — es entsteht keine |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | vom Runner geschrieben, committet |
| `docs/user/benutzerhandbuch.md` | update | §2 Liefer-Punkt 3 |
| `harness/README.md` | update | Zeile `make test-integration` |

- **Ort des Ladens (Entscheidung dieses Plans).** `Config` trägt zwei
  Pfade; die **Vollständigkeit** des Paars prüft die Konfiguration (so liest sie
  auch der Lauf `--healthcheck`, der dieselbe `ConfigFromEnvAndFile` ruft und die
  Dateien nicht braucht), das **Laden** (`tls.LoadX509KeyPair`) macht die
  Verdrahtung ganz am Anfang von `Run`, bevor `NewStream` den Slot anlegt
  (`internal/bootstrap/wiring.go`) — ein Ladefehler beendet den Start, ohne
  Datenbankverbindung und ohne Slot. Der Implementer misst die Reihenfolge am Arbeitsstand
  (`git grep -n 'NewStream' -- internal/bootstrap/wiring.go`) und belegt sie mit
  dem Test aus §2 (b).
- **Die zwei Adapter teilen keine TLS-Fassung.** Beide nehmen eine fertige
  `*tls.Config` entgegen; gebaut wird sie **einmal** in `internal/bootstrap`
  (Composition Root, [`ADR-0026`](../../adr/0026-composition-root.md)). Damit
  entsteht — anders als bei der Token-Prüfung — keine Doppel-Fassung in den
  Adaptern und keine Kante zwischen ihnen; `.a-check.yml` bleibt unverändert,
  `make a-check` Exit 0 ist DoD.
- **Zertifikat im Test.** Unit-Tests erzeugen das Zertifikat mit `crypto/x509`
  zur Laufzeit (selbstsigniert, Alternativname `127.0.0.1`/`localhost`); der
  Runner erzeugt es mit `certgen` im Toolchain-Container. Kein Zertifikat, kein
  Schlüssel liegt im Repo (ein committeter privater Schlüssel — auch ein
  Test-Schlüssel — wäre ein Befund für den Reviewer).
- **Reihenfolge:** (1) Startmessung: Suchlauf unten, dazu `--healthcheck`-Pfad
  lesen; (2) Konfiguration, Codes, Verdrahtungs-Test des Ladens; (3) die zwei
  Adapter mit den Unit-Tests; (4) Mutationsläufe an Kopien; (5) `certgen`, die
  Wegwerf-Clients, Runner-Phase, `make image`, `make test-integration`;
  (6) Handbuch, README, `make gates`.
- **Suchlauf (§3.13 der Regeln, [`AGENTS.md`](../../../../AGENTS.md)).** Bewegte
  Eigenschaft: **wie die zwei Server ihre Verbindungen annehmen** (Symbole
  `ListenAndServe`, `grpc.NewServer`, `net.Listen`; die Gegenseite
  `insecure.NewCredentials`), das **Healthcheck-Verhalten** (`--healthcheck`)
  und die Beschreibung der Schnittstellen als **unverschlüsselt** (Hedge
  „unverschlüsselt“, „ohne TLS“, „im Klartext“) samt der neuen Namen
  (`tls_cert_file`, `tls_key_file`, `CDC_TLS_`). Parent ist der Stand
  `a53f4e75…` (`git rev-parse HEAD` am Planungsstand, vor dem Plan-Commit; nie
  `HEAD`). Suchraum: ganzer Baum ohne `docs/reviews`, `done/`, `observations/`,
  `.harness/baseline`, die offenen Pläne und `docs/plan/adr` (Accepted ADRs
  halten ihren Stand, [`AGENTS.md`](../../../../AGENTS.md) §3.5).
  **Parent-Zeilen gemessen; die `diff`-Zeilen tragen den Planungsstand** — der
  Implementer setzt sie nach der Arbeit auf den Endstand und liest jede
  Trefferzeile. **Erwartung am Endstand** (hergeleitet, nicht gemessen): Zeile 1
  wächst um die TLS-Wege; Zeile 4 trägt neben der Spec das Handbuch; Zeile 5
  steigt um Handbuch, Code und Tests; Zeile 6 bleibt stehen, weil
  die übrigen Phasen und SDK-Läufe Klartext sprechen (sie sind der
  Boundary-Beleg); kein Treffer der Zeile 3 ändert sich im Verhalten — jeder
  Träger, der `--healthcheck` beschreibt, bleibt richtig und wird nicht
  umformuliert.

```suchlauf
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 8 -n -E 'ListenAndServe|grpc\.NewServer|net\.Listen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 9 -n -E 'insecure\.NewCredentials|WithInsecure' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 22 -n -E '\-\-healthcheck' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 6 -n -E 'unverschlüsselt|ohne TLS|im Klartext' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 5 -n -E 'tls_cert_file|tls_key_file|CDC_TLS_' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
a53f4e75ca51bdaa4319aeabba96f89c48b6448f 29 -n -E 'http://(cdc-test-feed|pg-change-feed)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 12 -n -E 'ListenAndServe|grpc\.NewServer|net\.Listen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 11 -n -E 'insecure\.NewCredentials|WithInsecure' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 30 -n -E '\-\-healthcheck' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 30 -n -E 'unverschlüsselt|ohne TLS|im Klartext' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 52 -n -E 'tls_cert_file|tls_key_file|CDC_TLS_' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 30 -n -E 'http://(cdc-test-feed|pg-change-feed)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
```

**Startmessung und Endstand (Implementer).** Am Start (Arbeitsstand `c45a3cca`
ohne Änderung) wich genau eine der zwölf Zeilen ab: Zeile 8 der `diff`-Gruppe
(`insecure.NewCredentials`) zählte 10 statt 9 — der Vorgänger-Slice hat dem
Wegwerf-Client `grpcadminclient` den Modus `probe` hinzugefügt (ein dritter
Aufbau der Verbindung); die Parent-Zeilen am Planungsstand stimmten alle. Am
Endstand (Arbeitsbaum, `make suchlauf-nachmessen` mit den hier stehenden Zahlen)
sind alle Trefferzeilen der sechs Muster gelesen: Zeile 1 (12) — `ListenAndServe`
fällt aus `internal/adapters/driving/http/server.go`, dort und im gRPC-Adapter
stehen `net.Listen` und `serve(listener net.Listener)`, dazu je ein Test pro
Paket mit `net.Listen` auf Loopback; kein Träger außerhalb des Codes
beschreibt den Annahmeweg. Zeile 2 (11) — hinzugekommen sind die zwei neuen gRPC-Tests und `grpcadminclient` mit
dem Klartext-Modus und der Transportwahl; die Beispiele selbst sind
unverändert. Zeile 3 (30) — die acht neuen Trefferzeilen sind neue Texte
(Handbuch, `tls.go`, Kommentare, Runner); jeder bestehende Träger, der
`--healthcheck` beschreibt, bleibt richtig und ist nicht umformuliert.
Zeile 4 (30) — die 24 neuen Trefferzeilen sind die neuen Texte des Zuges
(Handbuch, Kommentare, Tests, Runner, `harness/README.md` trägt `Klartext`
nur in der neuen Sensors-Beschreibung); die Spec-Zeilen (`spec/lastenheft.md`,
`spec/pflichtenheft.md`) sind unverändert und tragen die Aussage dieses
Slice. Zeile 5 (52) — Handbuch, `config_file.go`, `wiring.go`, die Tests und
der Runner; `spec/pflichtenheft.md` trägt die Namen wie zuvor. Zeile 6 (30) —
eine Zeile mehr als erwartet: die Konstante `TLS_HTTP_URL` der neuen
Runner-Phase (`http://pg-change-feed:8090` für den Klartext-Versuch); die
übrigen Phasen und SDK-Läufe sprechen unverändert Klartext (der Boundary-Beleg
„unverändert wie ohne diese Anforderung“).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1), **und**
[`slice-api-token-mehrfach-konfiguration`](../done/slice-api-token-mehrfach-konfiguration.md)
liegt in `done/` (Konfigurationsgleichlauf, Nummernvergabe der Codes), **und**
der Implementer hat die Startmessung gefahren: die Suchlauf-Zeilen aus §3 am
Arbeitsstand neu gemessen (neue Commit-Kennung als Parent, `make
suchlauf-nachmessen PLAN=` mit diesem Plan Exit 0 nach Anpassung von Parent und
Soll; jede Abweichung zum Planungsstand mit Ursache im Bericht — der Vorgänger
fügt in `tools/harness/run-integration-tests.sh` eine Phase hinzu und ändert
`wiring.go`, Zeilen und Zahlen verschieben sich).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Runner-Phase
  verlangt einen Umbau der Wegwerf-Clients über den TLS-Weg hinaus, oder der
  Slice wächst über die drei Punkte aus §2 (etwa SDK- oder
  Beispiel-Anpassungen; die gehören in `sdk-tls-optionen`).
- `in-progress` → `open` (blockiert — Carveout?): die Spec 0.15.0 ändert den
  Wortlaut von [`SPEC-034`](../../../../spec/pflichtenheft.md), oder der
  Compose-Stack ist am Messhost nicht startbar; kein Carveout, weil dann kein
  Gate rot ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, `make test-integration` und `make a-check` enden mit Exit 0,
`make suchlauf-nachmessen PLAN=` mit diesem Plan endet am Endstand mit Exit 0
(die `diff`-Zeilen auf den Endstand gesetzt, jede Trefferzeile gelesen), und die
Closure-Notiz in §7 trägt einen Lerneintrag (geschärfte Regel, neuer Sensor oder
benannte Spec-Lücke; Kandidaten: das gemessene Verhalten von `net/http` und
`grpc-go` gegenüber einem Klartext-Client, die Grenze „Ablauf des Zertifikats
wird beim Start nicht geprüft“ als benannte Spec-Lücke). Ein Gate, das am Stand
der Closure rot ist, geht nur mit dokumentiertem Carveout nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen
(Grund) · weiter offen (BEO-Eintrag im Register). Alle Ausgänge sind bis zur
Closure **offen** (Platzhalter `Ausgang: offen bis Closure`).

- **Das Verhalten gegenüber einem Klartext-Client ist ungeprüft.** `net/http`
  und `grpc-go` können einem Klartext-Client auf einem TLS-Port eine Fehlerantwort
  oder einen Verbindungsabbruch geben; „nicht bedient“
  ([`LH-FA-SST-011`](../../../../spec/lastenheft.md)) heißt: keine Antwort der API,
  auch nicht mit gültigem Token. Das konkrete Verhalten ist *nach Kenntnisstand*,
  die Unit-Tests und die Runner-Phase messen es. — **Ausgang: entfallen.**
  Gemessen: `net/http` antwortet auf `http://` mit `400` und dem Text „Client sent
  an HTTP request to an HTTPS server.“ (Antwort des Servers, keine der API),
  `grpc-go` endet mit `Unavailable`; Beleg: die Klartext-Unit-Tests beider
  Pakete, die gedruckte Zeile `CLEARTEXT …` des Runners und der eigene Lauf der
  Verifikation (§1, `make test-integration` Exit 0).
- **Der Negativ-Beleg „kein Klartext“ trifft einen noch nicht bereiten Server**
  und belegt dann eine Abwesenheit statt einer Auswahl
  ([`BEO-PGC/ready-ist-nicht-verbunden`](../observations/BEO-PGC/ready-ist-nicht-verbunden/observation.md),
  zwei Evidenz-Dateien am Planungsstand gezählt): gegen den Vorsatz steht die
  Reihenfolge in §2 (3), TLS-Erfolg zuerst. — **Ausgang: entfallen.** Der
  Klartext-Versuch folgt im Runner dem belegten TLS-Erfolg (Runner gelesen, eigener
  Lauf der Verifikation grün); die Prüfung „Port nach beendetem Container“ ist in
  der Fixrunde entfernt (Review F-2, das Register trägt die Klasse als weiteren
  Beleg bei `negativtest-ohne-bindung-an-seine-eingabe`).
- **Zertifikat im Image oder im Repo:** ein erzeugter privater Schlüssel landet
  in einer getrackten Datei oder im Image. Gegenmaßnahme: `certgen` schreibt in
  das Temp-Verzeichnis des Runners, nichts davon wird committet; der Reviewer
  prüft `git status` am Ende der Phase. — **Ausgang: entfallen.** Der
  Runner-Zähler der `.pem`-Dateien im Arbeitsbaum steht auf 0, und
  `git status --porcelain` ist nach dem Lauf der Verifikation leer.
- **Der Prozess liest die Zertifikatsdateien nicht** (Dateirechte, Benutzer des
  distroless-Images): die Runner-Phase scheitert dann am Start, nicht an der
  Spec. Die Rechte der Hilfsdateien folgen dem Muster der
  Konfigurationsdatei der Leerlauf-Phasen (`chmod`). — **Ausgang: entfallen.** `chmod 0644` auf die Wegwerf-Dateien
  (der Container läuft als `nonroot`, Kommentar im Runner); der Start mit gültigem
  Paar lief im Container bis `healthy` (Lauf der Verifikation).
- **Spec-Lücke zu „ungültig“:** ein abgelaufenes oder für den falschen Namen
  ausgestelltes Zertifikat lädt `tls.LoadX509KeyPair`; der Server startet, der
  Client lehnt ab. Der Plan nennt die Grenze (§1), die Spec benennt „ungültig“
  nicht näher. — **Ausgang: weiter offen** als benannte Spec-Lücke (Lerneintrag in §7;
  Handbuch und der Test `TestNewTLSConfigPruftAblaufUndNamenNicht` nennen die
  Grenze). Ein Spec-Nachtrag zu „ungültig“ ist ein Zug des Auftraggebers oder des
  Architects; die Adresse ist die Meldung im Bericht dieser Closure.
- **Zwei Konfigurationswege lesen verschieden:** `ConfigFromEnv` und
  `mergeConfig` bilden die Konfiguration getrennt; die Paar-Prüfung steht in
  **einer** Funktion, beide Wege rufen sie
  ([`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`](../observations/BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke/observation.md),
  fünf Evidenz-Dateien am Planungsstand; der Container-Start der Runner-Phase
  ist der Beleg über den echten Weg). — **Ausgang: entfallen.** Eine Funktion
  `validateTLSPair` in `internal/bootstrap/tls.go`, von beiden Wegen gerufen; die
  Mutation `&&` zu `||` färbt beide Wege rot (Verifikation M1), und der Runner
  startet den Container über die Umgebung **und** über die Datei.
- **Die SDK- und Beispiel-Aussage im Handbuch ist breiter als ihre Messung:**
  „die Client-Packages bieten keine TLS-Einstellung“ gälte für drei Sprachen;
  §2 Liefer-Punkt 3 verlangt je Sprache eine Messung
  ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
  Instanz B). — **Ausgang: entfallen.** Das Handbuch ist in der Fixrunde auf das
  Gemessene zurückgenommen (Review F-1); der Schluss „die gRPC-Beispiele
  verbinden deshalb nicht mit einem TLS-Server“ steht als *hergeleitet*, für
  Client-Pakete und HTTP-/SSE-Clients trifft es keine Aussage. Der Beispiel-Nachzug
  ist der Folge-Slice
  [`slice-examples-grpc-tls`](../open/slice-examples-grpc-tls.md).
- **Die neue Runner-Phase stört die Folge-Phasen** (Container mit Override neu
  erzeugt; Zustand muss am Phasen-Ende ohne Override wiederhergestellt sein). —
  **Ausgang: entfallen.** Alle Folge-Phasen liefen im Lauf der Verifikation grün
  im Klartext, der Container ist ohne Override wiederhergestellt
  (`make test-integration` Exit 0, keine `cdc-*`-Reste).
- **Die Bindung des Slot-Zählers der Runner-Phase ist hergeleitet** (in der
  Fixrunde und der Verifikation sichtbar geworden, kein Risiko des Planungsstands):
  die Mutation „Laden hinter `NewStream`“ ist am Runner nicht gefahren; die Zusage
  „kein Slot bei Ladefehler“ ruht auf dem Code-Aufbau und dem Unit-Test hinter der
  gefahrenen Mutation (`Run` lädt nach `postgresstorage.New`: rot). —
  **Ausgang: weiter offen**
  ([`BEO-PGC/runner-zaehler-ohne-gefahrene-mutation`](../observations/BEO-PGC/runner-zaehler-ohne-gefahrene-mutation/observation.md),
  1×).

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

- **Was hat funktioniert:** Der Schnitt trug: TLS ist Sache von Bootstrap und
  Driving Adaptern, `.a-check.yml` blieb unverändert, die Adapter nehmen eine
  fertige `*tls.Config`. Die Mutationen je Zusage (zehn Zeilen in §2, dazu die
  sieben der Verifikation, HTTP und gRPC je am eigenen Server) färbten rot; der
  Reviewer fand mit einer eigenen Mutation die Prüfung, die an keiner Eingabe rot
  werden konnte (F-2), vor dem Merge. Die Reihenfolge „erst TLS-Erfolg, dann
  Klartext-Versuch“ im Runner hielt die Klasse
  [`BEO-PGC/ready-ist-nicht-verbunden`](../observations/BEO-PGC/ready-ist-nicht-verbunden/observation.md)
  fern (vermieden, nicht eingetreten; keine neue Datei, bleibt 2×).
- **Was ging anders als geplant:** (1) Der Plan sagte „ohne geöffneten Port“ zu;
  die Prüfung las einen beendeten Container und entfiel in der Fixrunde samt der
  Aussage in Plan, Handbuch, `harness/README.md` und Runner-Zeile (Review F-2).
  (2) Das Handbuch trug zunächst eine Verallgemeinerung über Client-Pakete und
  HTTP-/SSE-Clients, die nicht gemessen war; es sagt nach der Fixrunde nur das
  Gemessene (Review F-1). (3) Die Datei `internal/bootstrap/tls.go` entstand statt
  eines Einbaus in `wiring.go` (§3, Implementer-Nachzug). (4) Der Beispiel-Nachzug
  (Review F-5) war im Plan nicht als Slice benannt; er steht jetzt als Datei.
  Die Verifikation fuhr die Runner-Mutation „Laden hinter `NewStream`“ nicht
  (rund fünfzehn Minuten); sie bleibt *hergeleitet*.
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke, keine verkörperte Regel (ohne
  `liegt in`): [`SPEC-034`](../../../../spec/pflichtenheft.md) nennt ein
  Zertifikat „nicht lesbar, ungültig oder passt nicht zusammen“, sagt aber nicht,
  ob Ablauf und Name beim Start geprüft werden; `tls.LoadX509KeyPair` prüft
  beides nicht (gemessen: `TestNewTLSConfigPruftAblaufUndNamenNicht` lädt ein
  abgelaufenes Zertifikat ohne Fehler), der Server startet, der Client lehnt ab.
  Handbuch und Plan nennen die Grenze. Ein Spec-Nachtrag ist ein Zug des
  Auftraggebers oder des Architects, nicht dieser Closure; die Adresse ist die
  Meldung im Bericht an den Hauptlauf. Zweiter Lerninhalt, **gemessen** (Unit-Tests
  und Runner): `net/http` beantwortet einen Klartext-Aufruf auf einem TLS-Port mit
  `400` und Text, `grpc-go` beendet den Verbindungsaufbau (`Unavailable`) — „nicht
  bedient“ heißt hier „keine Antwort der API“, nicht „keine Antwort“.
- **Beobachtungs-Register (`../observations/`):** Weitere Datei
  `evidence/slice-tls-http-grpc-server.md` in
  [`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`](../observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/observation.md)
  (Review F-2, MEDIUM, daher Datei trotz Deckel; **neue Form**: Abwesenheits-Prüfung
  gegen einen beendeten Container, ein Beleg über einen Gegenstand, der beim Lesen
  nicht mehr existiert; der Mechanismus ist derselbe wie dort — eine Aussage ohne
  Bindung an ihre Eingabeseite —, ein eigenes Verzeichnis würde die Klasse
  spalten); der Zähler folgt aus den Dateien und steht bei 26 (`ls evidence`).
  Neu angelegt
  [`BEO-PGC/runner-zaehler-ohne-gefahrene-mutation`](../observations/BEO-PGC/runner-zaehler-ohne-gefahrene-mutation/observation.md)
  (1×, Ausgang des Risikos „weiter offen“). Kein Zuwachs bei
  [`BEO-PGC/ready-ist-nicht-verbunden`](../observations/BEO-PGC/ready-ist-nicht-verbunden/observation.md)
  (bleibt 2×, siehe oben). Die in §8 genannten Einträge mit mindestens drei
  Dateien sind mit diesem Slice berührt und haben ihre Gegenmaßnahme getragen
  (Container-Start der Runner-Phase für `adapter-unittest-verdeckt-bootstrap-luecke`,
  Gegenprobe je Negativfall für `negativtest-ohne-bindung-an-seine-eingabe`,
  Handbuch-Liefer-Punkt samt Historienzeile 1.96); für sie entstand keine weitere
  Datei, außer der genannten. Der Eintrag `negativtest-ohne-bindung-an-seine-eingabe`
  war schon vor dieser Closure über der Schwelle (Ausgang *verkörpert*); die Closure
  eines Slice liest ihn nicht neu, kein Kandidat der Schärfung. Kein Eintrag mit
  `evidence/` ≥ 3 ist neu über die Schwelle gewachsen.
- **Validator-Feststellung (Modul 8):** Der Slice liefert Betreiber-Wert
  (verschlüsselte HTTP-/gRPC-Schnittstelle, [`LH-FA-SST-011`](../../../../spec/lastenheft.md)).
  Belegt am Realserver (Runner-Phase, Verifikation: TLS-Erfolg, Klartext nicht
  bedient, `--healthcheck` gesund, fünf Negativstarts mit Code im Log). Der
  Nutzerbedarf selbst (ein Betreiber betreibt den Server hinter eigenem
  Zertifikat, ohne Proxy) ist **nicht** am realen Bedarf geprüft; ein
  Validator-Lauf ist nach dem nächsten Release sinnvoll, mit Token-Wechsel und
  OTLP (Betreiber-Oberfläche als Ganzes). Release-Folge dieser Closure: keine
  (kein SDK-Zug, ein Release braucht eine Freigabe).
- **Folge-Slices:**
  [`slice-examples-grpc-tls`](../open/slice-examples-grpc-tls.md) (TLS-fähige
  gRPC-Beispiele in Go, C# und Kotlin; Review F-5 und Verifikation) — Datei in
  `open/`, keine harte Abhängigkeit, kein SDK-Release; die Beispiele bauen ihren
  Kanal selbst, deshalb muss er nicht auf den SDK-Slice warten. Und
  `sdk-tls-optionen` — **keine Datei**: er wartet auf eine Folge-Anforderung im
  Lastenheft für die TLS-Optionen der SDK-Packages
  ([`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) §Konsequenzen; das
  Lastenheft nennt sie „eigene, nachgelagerte Anforderung“). Voraussetzung: dieser
  Slice als Server-Beleg (erfüllt) und die Anforderung selbst — ein Spec-Zug des
  Auftraggebers oder des Architects, den der Planner nicht anlegt. Adresse dieses
  Aufschubs ist diese Notiz; der Slice wird mit der Anforderung als Datei in
  `open/` angelegt.
- **Risiken aus §6:** neun, je ein Ausgang in §6 selbst: Klartext-Verhalten
  entfallen (gemessen) · Negativ trifft nicht bereiten Server entfallen ·
  Zertifikat im Repo entfallen · Dateirechte entfallen · Spec-Lücke Ablauf und Name
  weiter offen (Lerneintrag) · zwei Konfigurationswege entfallen · SDK-/Beispiel-Aussage
  entfallen (Folge-Slice `examples-grpc-tls`) · Runner-Phase stört Folge-Phasen
  entfallen · Runner-Slot-Zähler hergeleitet weiter offen (BEO angelegt).
- **Drei Paarungen:** (a) *Anker* — kein Eintrag trägt `liegt in`
  (nichts verkörpert); entfällt. (b) *Folge-Slice* —
  [`slice-examples-grpc-tls`](../open/slice-examples-grpc-tls.md) existiert in
  `open/`; `sdk-tls-optionen` ist als Aufschub mit Adresse benannt, **keine**
  Datei, und trägt deshalb keine Paarung. (c) *Register* —
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
  `BEO-PGC/runner-zaehler-ohne-gefahrene-mutation`,
  `BEO-PGC/ready-ist-nicht-verbunden` und die in §6/§8 genannten Kennungen
  existieren als Verzeichnisse mit nicht leerem `evidence/` (Messung nach dem
  `git mv`). Träger fremder Dateien: der Suchlauf (§3) ist am Endstand gelesen und
  vom Implementer berichtet; die Links der Verweise auf diese Datei sind im
  Reconcile-Commit nach dem Move nachgezogen (Slice
  [`slice-otlp-metrik-export`](../open/slice-otlp-metrik-export.md), Plan und
  Verifikation des Vorgängers `slice-api-token-mehrfach-konfiguration`, deren
  Review/Verifikation).

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
berührten Pfade (`internal/…`, `tools/harness/`, `docs/user/`, `harness/`)
liegen alle in ihr, es entsteht keine neue Sub-Area — die Schwelle
(≥ 2 von 3 Achsen) wird nicht angewandt, weil kein Pfad ausdifferenziert werden muss.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/BEO-PGC/`) am Planungsstand durchgegangen;
Treffer für die berührte Sub-Area, Zähler = Zahl der `evidence/`-Dateien:
`adapter-unittest-verdeckt-bootstrap-luecke` (5, erreicht 3× — Gegenmaßnahme:
Container-Start der Runner-Phase, §2/§6), `ready-ist-nicht-verbunden` (2, §2 (3)),
`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (3, eigener
Liefer-Punkt), `handbuch-versionshistorie-uebersprungen` (3, Historien-Zeile),
`zwei-quellen-drift-handbuch-gegen-pflichtenheft` (3, Handbuch und
Pflichtenheft werden für die Datei-Schlüssel gegeneinander gelesen),
`intern-kennungen-in-ausgelieferten-texten` (1, die neuen Meldungstexte tragen
keine Kennung, `make ausgabe-kennungen-check`), `kommentar-herkunft-als-kette`
(4, höchstens eine Kennung je Kommentar im neuen Code,
[`AGENTS.md`](../../../../AGENTS.md) §3.7),
`negativtest-ohne-bindung-an-seine-eingabe` (25, Gegenprobe je Negativfall:
gültiges Paar neben jedem Ladefehler). Die Einträge mit mindestens drei Dateien
sind **mit** diesem Slice berührt und tragen deshalb hier ihre Gegenmaßnahme.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
