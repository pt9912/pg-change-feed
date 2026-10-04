# Review-Report: slice-tls-http-grpc-server — 2026-10-04

**Review-Art:** Code (gegen Plan, ADR und Hard Rules)

**Gegenstand:** `git diff c45a3cca HEAD` (Commits 8e3fa08b Code, 08a796f2 Runner-Phase, Handbuch 1.96, Plan)

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09 mit den Schärfungen bis `slice-harness-mutationsbild-und-verweigerte-aktion`)
**Modell:** Sonnet 5.5 · **Datum:** 2026-10-04

**Eingangs-Kontext:**

- Plan `slice-tls-http-grpc-server` (in-progress)
- [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md), [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md), [`ADR-0152`](../plan/adr/0152-zugangsdaten-klasse-elf-schluessel.md), [`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-FA-SST-011`](../../spec/lastenheft.md), `SPEC-034` und `SPEC-016` im [Pflichtenheft](../../spec/pflichtenheft.md)
- `AGENTS.md` (Hard Rules), `harness/conventions.md`

Läufe des Reviewers (Exit direkt gelesen, **gemessen** am Stand HEAD): `make gates` Exit 0, `make test` Exit 0 (alle Pakete `ok`), `make fmt-check` Exit 0, `make kommentar-kennungen DIFF=c45a3cca` Exit 0, `make suchlauf-nachmessen PLAN=…` Exit 0 (12 Zeilen stimmen), `make doc-trace`: `83 Anforderung(en), 1 Waise(n).`, Zeile `LH-FA-SST-011 … E2E | ok`. `make a-check`, `make handbuch-public-doc-check`, `make ausgabe-kennungen-check`, `make meldungscodes-check` laufen in `make gates` (Exit 0). `make test-integration` ist **übernommen** vom Implementer (Verifier fährt ihn).

Mutationen (Kopie im Scratchpad, `go test -race` im Toolchain-Container, Rücknahme per `cp`, Echtrepo `git status --short` leer), **gemessen**, drei Stellen, die die Tabelle des Plans nicht führt oder an anderer Stelle mutiert:

| Mutation | Instanz | Farbe |
|---|---|---|
| `newTLSConfig` trägt zusätzlich `ClientAuth: RequireAnyClientCert` (mTLS) | `internal/bootstrap` | rot: `TestNewTLSConfigBautEineKonfiguration`, `TestNewTLSConfigMindestversionAmDraht` |
| `LoadX509KeyPair(keyFile, certFile)` vertauscht | `internal/bootstrap` | rot: fünf Tests, u. a. `TestRunLaedtTLSPaarVorJederVerbindung` |
| `mergeConfig`: Rangfolge Datei/Umgebung für `TLSCertFile` vertauscht | `internal/bootstrap` | rot: `TestTLSPaarWirdGelesen` |

---

## Findings

### F-1 — Handbuch trägt Entwickler-Beleg und eine nicht gemessene Verallgemeinerung zu den Clients

- `kategorie`: MEDIUM
- `quelle`: Skill `nutzerdoku-schreiben` (Betreibersicht), `AGENTS.md` §3.12 Instanz B, Skill „Beleg trägt seinen Satz nicht“
- `pfad`: `docs/user/benutzerhandbuch.md:1788-1795` (Absatz „Clients“) und Änderungshistorie 1.96
- `befund`: Der Absatz belegt mit `git grep`-Befehlen über interne Quelldateien (Entwicklersicht, für den Betreiber nicht ausführbar) und folgert aus „keine TLS-Option in den Options-Dateien“ den Satz „spricht deshalb nur mit einem Server ohne TLS-Paar“. Der Schluss ist nur für die gRPC-Beispiele gemessen; ein HTTP-/SSE-Aufruf gegen einen TLS-Server ist nach Plan-Angabe nicht gefahren, und fehlende Optionen schließen einen `https://`-Zugang über den Standard-Vertrauensspeicher nicht aus. Die Historie-Zeile 1.96 wiederholt die Verallgemeinerung („die Client-Pakete … sprechen nur mit einem Server ohne TLS-Paar“).
- `verifizierbar`: nein — kein Gate liest Sinn; Probe ist ein Lauf der HTTP-/SSE-Clients gegen die TLS-Phase.
- `klasse`: Beleg trägt seinen Satz nicht

### F-2 — Klartext-Prüfung der Negativstarts kann nicht rot werden

- `kategorie`: MEDIUM
- `quelle`: Skill „Zusage ohne Bindung an ihre Eingabeseite“, [`LH-FA-SST-011`](../../spec/lastenheft.md) Negative
- `pfad`: `tools/harness/run-integration-tests.sh:4635-4646` (`tl_expect_start_refused`), Handbuch „Start-Fehler: … ohne geöffneten Port“
- `befund`: Die Prüfung „auf den Adressen antwortet kein Klartext-Server“ läuft nach `tl_await_stopped`, also gegen einen beendeten Container; `transport_error`/`Unavailable` folgen dann aus dem Stillstand, unabhängig davon, was `Run` vor dem Ende getan hat. Gebunden ist die Reihenfolge nur gegen die Datenbankverbindung (`TestRunLaedtTLSPaarVorJederVerbindung`, der Slot-Zähler im Runner); die Aussage „ohne geöffneten Port“ trägt weder der Runner noch eine Mutation (ein vor dem Laden gebundener Listener färbt keinen der beiden rot, solange der Lauf danach scheitert).
- `verifizierbar`: ja — Mutation: Listener-Aufbau vor `newTLSConfig` ziehen; Unit-Test und Runner-Zeile bleiben dort grün.
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-3 — Runner-Kopf und Kommentar nennen vier statt fünf Starts

- `kategorie`: LOW
- `quelle`: Maintainability, `AGENTS.md` §3.12
- `pfad`: `tools/harness/run-integration-tests.sh:4561` und `:4738`
- `befund`: Die Kommentare beschreiben „vier Starts mit je einer Verletzung“ bzw. „Vier Starts“; der Code fährt fünf (nur Zertifikat, nur Schlüssel, nicht lesbarer Pfad, Datei ohne PEM, nicht zusammenpassendes Paar), und Plan, Handbuch und Abschlusszeile sagen fünf. Der Kommentar nennt zudem „Ausgang 2 / Ausgang 1“ für vier Fälle.
- `verifizierbar`: ja — Zählung der Aufrufe `tl_expect_start_refused` (5).
- `klasse`: Zahl im Träger driftet gegen die Messung

### F-4 — Gesamtes Schlüsselmaterial der Phase wird world-readable

- `kategorie`: LOW
- `quelle`: Sicherheit (Hygiene)
- `pfad`: `tools/harness/run-integration-tests.sh` (`chmod 0755 "$TLS_TMP"`, `chmod 0644 "$TLS_TMP"/*.pem`), `tools/harness/certgen/main.go` (Schlüssel `0600`)
- `befund`: `certgen` schreibt den Schlüssel bewusst mit `0600`, der Runner setzt danach `0644` auf alle `*.pem` inklusive der Schlüssel (nötig, damit der unprivilegierte Container-Benutzer sie liest). Die Schlüssel sind Wegwerf-Material (24 h, im Temp-Verzeichnis, vom `cleanup`-Trap entfernt); der Handbuch-Hinweis „beide Dateien müssen für ihn lesbar sein“ ist konsistent.
- `verifizierbar`: nein
- `klasse`: Testmaterial-Rechte

### F-5 — Beispielprogramme tragen veraltete Aussage „Der Feed-Container spricht Klartext-gRPC (kein TLS)“

- `kategorie`: INFO (Träger außerhalb des Diffs, `AGENTS.md` §3.13)
- `quelle`: [`LH-FA-SST-011`](../../spec/lastenheft.md), `AGENTS.md` §3.13
- `pfad`: `examples/csharp/grpc-client/Program.cs:49`, `examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/Main.kt:56`
- `befund`: Der Feed-Container kann jetzt TLS sprechen; der Kommentar beschreibt nur noch die Demo-Umgebung. Der Plan hält Beispiele ausdrücklich unverändert (Suchlauf-Zeile 2), das Handbuch benennt die Grenze; eine Adresse für einen Beispiel-Nachzug ist im Plan nicht als Folge-Slice benannt.
- `verifizierbar`: ja — `git grep -n 'kein TLS' examples`.
- `klasse`: Träger-Nachzug

### F-6 — Fehlerzustand beim Laden steht nicht im Heartbeat; Ausgänge 2 und 1

- `kategorie`: INFO
- `quelle`: [`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md), `LH-QA-REL-001`
- `pfad`: `internal/bootstrap/wiring.go` (`Run`, `newTLSConfig` am Anfang), `internal/bootstrap/tls.go`
- `befund`: `PCF-E2009` endet mit Ausgang 2 vor `Run`, `PCF-E2010` mit Ausgang 1 aus `Run`; beides ist im Plan und im Handbuch benannt und folgt der bestehenden Aufteilung der Konfigurationsfehler (jeder Modus — `--healthcheck`, `diagnose`, CLI — lädt `ConfigFromEnvAndFile`, prüft also die Vollständigkeit, aber lädt das Paar nicht). Kein Heartbeat-Fehlerzustand, weil `Run` vor dem Heartbeat-Aufbau endet, wie jeder Fehler an dieser Stelle. Als tragfähig gelesen, keine Spec-Lücke.
- `verifizierbar`: ja — Runner-Phase (übernommen).
- `klasse`: Einordnung

## Negativbefunde (geprüft, ohne Befund)

- `internal/bootstrap/tls.go`: `MinVersion` TLS 1.2 explizit, kein `ClientAuth`, keine `CipherSuites`, kein `InsecureSkipVerify`; Zertifikat/Schlüssel in der Reihenfolge `LoadX509KeyPair(cert, key)` (Vertauschung färbt fünf Tests rot); Fehlertext nennt Variablen und Pfade, nie Schlüsselinhalt (der Ursachetext von `crypto/tls` trägt keinen).
- `internal/bootstrap/wiring.go` `Run`: zwischen Start und `newTLSConfig` stehen nur Telemetrie-Aufbau und `defer`; vor `postgresstorage.New`/`NewStream`/Slot-Anlage kein Port, kein Schreibzugriff.
- `internal/adapters/driving/http/server.go`: Verhalten ohne TLS unverändert (`Serve`); Bind-Fehler nun synchron in `Start`, aber der Aufrufer ist dieselbe Goroutine mit Log wie zuvor; `Shutdown` vor `Serve` liefert `ErrServerClosed`, der Listener wird von `Serve` geschlossen; `Addr` leer nur durch die Verdrahtung ausgeschlossen (Aufruf nur mit gesetzter Adresse).
- `internal/adapters/driving/grpc/server.go`: `grpc.Creds` nur bei gesetzter Konfiguration; keine adapters→adapters-Kante (`make a-check` in `make gates` Exit 0, `.a-check.yml` unverändert).
- Konfiguration: `tls_cert_file`/`tls_key_file` nicht in `forbiddenFileCredentialKeys` (Klasse bleibt, Test vorhanden, Mutation des Implementers übernommen); beide Zugriffswege rufen `validateTLSPair`; Rangfolge Umgebung > Datei durch eigene Mutation belegt.
- Meldungscodes `PCF-E2009`/`PCF-E2010` in `codes.go`, Katalog und Sentinel-Tabelle; `make meldungscodes-check`, `make ausgabe-kennungen-check`, `make handbuch-public-doc-check` grün (`make gates`).
- Runner-Phase (Lesen, nicht gefahren): Temp-Verzeichnis im `cleanup`-Trap, `git status --porcelain`-Prüfung auf `.pem`, Klartext-Versuch erst nach belegtem TLS-Erfolg, Wiederherstellung ohne Override mit Klartext-Probe, Abdeckungszeile `LH-FA-SST-011` in `docs/user/e2e-abdeckung.md` (`make doc-trace` ok). Hilfs-Slot `slot_tls_neg` wird nur im Erfolgspfad entfernt; bei Abbruch trägt die Wegwerf-Umgebung ihn nicht weiter (kein Befund).
- Mutationstabelle des Plans: Kennzeichnung „hergeleitet“ für den Draht-Test an `httptest` statt an den Adaptern und für `MinVersion` (nur die Feldprüfung belegt) ist korrekt.
- Grenze „Zertifikat und Kette“: im Plan an `crypto/tls` begrenzt (`TestNewTLSConfigLiestDieKetteDerZertifikatsdatei`), im Handbuch nur als „das Zertifikat und dahinter seine Kette“ ohne Client-Messaggrenze — ohne Befund.
- Handbuch-Versionskopf 1.96 und Historie-Zeile ohne Kennungen; Ablauf/Name-Grenze benannt.
- Suchlauf §3.13: Träger „ohne TLS“/„Klartext“ in README, Compose, SDK-READMEs: keine Treffer; Ausnahmen siehe F-5.
- Kommentare im Diff: keine Kennungsketten (`make kommentar-kennungen`), keine Chronik in Produktionscode.
- Docker-only/§3.1/§3.9: keine Host-Toolchain, keine in-place-Werkzeuge, keine Umleitung in Repo-Dateien (die Redirects des Runners gehen in `$TLS_TMP`).

## Verdikt

**Merge-blockierend: nein.** 0 HIGH. Der Sicherheitskern (TLS-Konfiguration, Reihenfolge in `Run`, keine Klartext-Rückfallstufe, kein mTLS, Pfad statt Schlüsselmaterial in Fehlern) trägt, drei eigene Mutationen färben rot. F-1 und F-2 (MEDIUM) sind Nachbesserungen an Handbuch-Wortlaut bzw. Belegtiefe, F-3/F-4 LOW, F-5/F-6 INFO. Fixrunde empfohlen, nicht zwingend. DoD-Zeile „Review durchgeführt“ bleibt bewusst unberührt (Auftrag).

**Übergabe:** F-1/F-2/F-3 an den Implementer; F-5 an den Planner (Adresse für den Beispiel-Nachzug); F-4/F-6 ohne Aktion.
