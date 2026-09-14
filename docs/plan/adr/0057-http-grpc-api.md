# ADR-0057: HTTP/JSON-API mit Token-Authn (Supersedes ADR-0020)

**Status:** Accepted — Supersedes [`ADR-0020`](0020-http-grpc-optional.md) (Re-Evaluierungs-Trigger „beobachtbarer Bedarf eines API-Consumers" eingetreten)

**Datum:** 2026-09-14

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-14)

**Bezug:** [`LH-FA-SST-006`](../../../spec/lastenheft.md) (Haupt-Bezug —
konkrete API), [`LH-FA-SST-005`](../../../spec/lastenheft.md)
(Vorbereitung, bereits erfüllt), [`LH-QA-SEC-001`](../../../spec/lastenheft.md),
[`LH-QA-SEC-002`](../../../spec/lastenheft.md),
[`LH-QA-SEC-003`](../../../spec/lastenheft.md), [ADR-0020](0020-http-grpc-optional.md),
[ADR-0028](0028-inbound-use-cases.md) (Inbound Use Cases, die diese API
exponiert), [ADR-0034](0034-ports-nach-faehigkeiten.md) (Port-Zuschnitt),
[ADR-0046](0046-sql-driving-adapter-lese-schreib-trennung.md) (Lese-/
Schreib-Trennung als Präzedenzfall für „nicht jeder Zugriffsweg braucht
sofort jede Fähigkeit"), [ADR-0047](0047-rollenspezifische-dsn-verdrahtung.md)
(Rollen-Modell, das die Token-Klassen hier spiegeln)

**Schärft:** [`ARC-005`](../../../spec/architecture.md) (macht das dort
bereits generisch genannte „später HTTP-/gRPC" konkret: HTTP/JSON, jetzt)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0020`](0020-http-grpc-optional.md) verschob die Wahl zwischen
HTTP/gRPC auf einen späteren Zeitpunkt und trug einen expliziten
Re-Evaluierungs-Trigger: „Beobachtbarer Bedarf eines API-Consumers —
sichtbar als Anforderung im Lastenheft-Change oder als Eintrag im
Beobachtungs-Register." Dieser Trigger ist eingetreten — nicht als
Ableitung aus „MVP ist vorbei", sondern weil der Auftraggeber (pt9912) den
Bedarf am 2026-09-14 explizit als Auftrag benannt hat, samt der zusätzlich
geforderten Nachweisform (E2E-Test, Beispiel-Client). Das Lastenheft trägt
bereits seit Version 0.4.0 die dafür vorbereitete Anforderung
[`LH-FA-SST-006`](../../../spec/lastenheft.md) — sie grenzt Protokoll-Wahl,
Authn/Authz-Verfahren sowie Endpunkt-Details ausdrücklich als
Architektur-/Spezifikationsfrage aus. Diese ADR ist diese Entscheidung.

Vier Konstraints prägen den Lösungsraum:

- **Bestehender Adapter-Zuschnitt.** Der einzige heute existierende
  Driving Adapter neben CLI (`cmd/pg-change-feed`) und SQL
  (`ADR-0046`-Views) ist `internal/adapters/driving/replication` — ein
  schlanker Zuschnitt ohne Fremd-Framework. `ARC-005` nennt HTTP/gRPC
  bereits generisch als künftige Driving-Adapter-Technologie.
- **Nicht jede Fähigkeit liegt hinter einem Inbound Port.** Über
  `internal/application/port/inbound/` vorhanden: Consumer-Registrierung/
  -Bestätigung/-Position/-Entfernung (`consumer.go`), Tabellen-Aktivierung/
  -Deaktivierung/-Status/-Liste (`verwaltung.go`), Retention-Lauf
  (`retention.go`). `Diagnose` und `Healthcheck`
  (`internal/bootstrap/wiring.go`, Zeilen ~1144/~1199) sind freie
  Funktionen ohne Anwendungsschicht-Abstraktion — für einen Driving
  Adapter, der laut `ARC-002`/`ARC-003` nur über Inbound Ports mit der
  Application-Schicht sprechen darf, fehlt hier die Vermittlung. Das
  Lesen über `cdc.changes` (`LH-FA-REA-001`ff.) läuft nach `ADR-0046`
  bewusst direkt über eine SQL-View, ohne Anwendungsdienst dahinter.
- **Sicherheitsanforderungen sind bisher DB-Rollen-gebunden.**
  `LH-QA-SEC-001`…`003` sind ausschließlich über PostgreSQL-Rollen/DSNs
  gelöst (`ADR-0047`: `cdc_capture`/`cdc_admin`/`cdc_reader`, je eine
  eigene DSN). Ein Netzwerk-Zugriffsweg ohne DB-Login braucht einen
  eigenen Authn-Mechanismus; `LH-QA-SEC-003` lässt das ausdrücklich zu
  („PostgreSQL-Berechtigungen **oder einen gleichwertigen Mechanismus**").
- **`LH-FA-SST-006` Negative-Akzeptanzkriterium ist bindend:** ein nicht
  autorisierter Aufruf über die API „wird abgelehnt, nicht stillschweigend
  ignoriert" — ohne jeden Authn-Mechanismus ist „nicht autorisiert"
  strukturell nicht abbildbar.

## Entscheidung

Wir wählen **HTTP/JSON (REST-artig) über einen neuen Driving Adapter
`internal/adapters/driving/http/`**, in der ersten Version beschränkt auf
die bereits Port-gedeckten Fähigkeiten, geschützt durch statische
Bearer-Tokens in zwei Rechtsklassen. Vier Teilfragen, vier Festlegungen:

### Teilfrage 1 — Protokoll

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (`ADR-0020` bleibt in Kraft, kein konkreter Adapter) | kein Aufwand | widerspricht dem eingetretenen Trigger und dem expliziten Auftrag |
| B — gRPC | typed Contracts, Streaming eingebaut | braucht eine protoc-/buf-Toolchain im Build (neue Docker-Build-Stufe, verstößt gegen den bisher schlanken Abhängigkeits-Fußabdruck — vgl. `ADR-0055`s Beobachtung, dass `nats.go` „die erste direkte Nicht-PostgreSQL-Abhängigkeit" ist); Konsumenten brauchen Codegen; kein bekannter gRPC-Consumer benennt diesen Bedarf |
| C — beide parallel (HTTP + gRPC) | deckt beide Konsumenten-Erwartungen ab | doppelter Implementierungs- und Wartungsaufwand für dieselbe Fähigkeitsmenge, ohne dass ein zweiter Consumer-Typ diesen Aufwand rechtfertigt |
| **D — HTTP/JSON, REST-artig (gewählt)** | Go-Standardbibliothek (`net/http`, `encoding/json`) genügt — keine neue Fremd-Abhängigkeit; universell konsumierbar (curl, jede Sprache ohne Codegen); passt zum bisher schlanken Adapter-Zuschnitt (CLI/SQL sind ebenfalls ohne komplexe Frameworks) | kein eingebauter typed-Contract-/Streaming-Mechanismus; eine spätere gRPC-Ergänzung bleibt eine eigene Entscheidung, falls ein konkreter Bedarf entsteht |

### Teilfrage 2 — Umfang der ersten API-Version

| Option | Pro | Contra |
|---|---|---|
| A — nichts exponieren (nur Adapter-Grundgerüst ohne Fähigkeit) | kleinster erster Schritt | erfüllt kein Akzeptanzkriterium von `LH-FA-SST-006`; kein sinnvoller Liefer-Schnitt |
| **B — nur die bereits Port-gedeckten Fähigkeiten (Consumer-Registrierung/-Bestätigung/-Position/-Entfernung, Tabellen-Aktivierung/-Deaktivierung/-Status/-Liste, Retention-Lauf) (gewählt)** | reiner Adapter über bestehende Inbound Ports (`ARC-003`) — keine neue Anwendungsschicht-Entscheidung nötig, volle Schichten-Konformität ohne Zusatzaufwand; deckt `LH-FA-SST-006`s Happy-Path/Boundary/Negative bereits exemplarisch ab (z. B. Consumer-Registrierung über CLI **und** API mit fachlich gleichwertigem Ergebnis) | Changes lesen (`LH-FA-REA-*`) und Diagnose/Health bleiben vorerst CLI-/SQL-exklusiv |
| C — zusätzlich Changes lesen (`LH-FA-REA-*`) über einen neuen `ReadChangesUseCase`/Inbound Port | deckt den vermutlich wichtigsten Consumer-Bedarf (Lesen) sofort | `cdc.changes` ist nach `ADR-0046` bewusst ein SQL-View-Direktzugriff ohne Anwendungsdienst; eine neue Port-Abstraktion dafür (Filter-/Pagination-Signatur, Konsistenz zur View) ist eine eigene architektonische Entscheidung, die diese ADR nicht nebenbei treffen sollte |
| D — zusätzlich Diagnose/Healthcheck über eine neue Inbound-Port-Abstraktion | deckt Monitoring-Bedarf sofort | `Diagnose`/`Healthcheck` sind heute freie Funktionen ohne Anwendungsschicht-Abstraktion — dieselbe Argumentation wie C: eigene Entscheidung wert (welche Query-Form, welcher Rückgabetyp) |
| E — voller Umfang (B + C + D) in einer ersten Version | ein einziger großer Wurf | zu groß für einen ersten, in einer Review-Sitzung prüfbaren Schnitt; mehrere neue Ports auf einmal ohne graduellen Beleg |

Changes-Lesen und Diagnose/Health bleiben damit **bewusst ausgeschlossen** —
nicht verworfen, sondern vertagt: Beide brauchen vorher eine eigene
Port-Design-Entscheidung (Teilfrage siehe Konsequenzen/Folgepflicht), die
diese ADR nicht mitentscheidet, um den Umfang dieser Entscheidung nicht zu
sprengen.

### Teilfrage 3 — Authentifizierung/Autorisierung

| Option | Pro | Contra |
|---|---|---|
| A — kein Schutz, Netzwerk-Perimeter vorausgesetzt (nichts tun) | kein Zusatzaufwand | `LH-FA-SST-006`s Negative-Akzeptanzkriterium wird strukturell unerfüllbar — ohne jeden Mechanismus kann kein Aufruf „nicht autorisiert" sein |
| B — mTLS (Client-Zertifikate) | starke Authentifizierung ohne Klartext-Secret im Request | hoher Betriebsaufwand (eigene CA/Zertifikatsverwaltung), passt nicht zum bisherigen Betriebsmodell (ENV-DSNs, kein PKI-Bestand im Repo), unverhältnismäßiger Aufwand für einen ersten Schnitt |
| **C — statische Bearer-Tokens in zwei Rechtsklassen, analog zum DB-Rollenmodell (gewählt)** | minimaler Aufwand (Header-Vergleich gegen konfigurierte Tokens, kein neuer Fremd-Dienst); spiegelt das bewährte Least-Privilege-Muster aus `ADR-0047` (`cdc_reader` vs. `cdc_admin`); erfüllt das Negative-Akzeptanzkriterium durch klare Ablehnung (401 fehlend/ungültig, 403 falsche Klasse) | Tokens sind Shared Secrets ohne Ablauf-/Widerrufsmechanismus in dieser ersten Version; Klartext-Transport verlangt TLS-Terminierung vor der API — Betreiber-Pflicht, nicht Gegenstand dieser ADR |
| D — OAuth2/JWT mit externem Identity Provider | Industriestandard, delegierbare Autorisierung | für ein Werkzeug mit heute einem Auftraggeber/Consumer-Kreis unverhältnismäßig; neue externe Abhängigkeit (IdP) ohne belegten Bedarf |

Zwei Token-Klassen, benannt analog zu den bestehenden DB-Rollen
(`ADR-0047`), aber **orthogonal** zu ihnen — die API-Token-Prüfung
entscheidet an der HTTP-Schicht, welche Use Cases ein Aufruf erreichen
darf; welche DSN der darunterliegende Adapter tatsächlich benutzt, bleibt
unverändert die bei der Verdrahtung (`ADR-0047`) fixierte:

- **`CDC_API_TOKEN_READER`** — lesende Aufrufe: `GetConsumerPosition`,
  `GetStatus`, `ListTables`.
- **`CDC_API_TOKEN_ADMIN`** — schreibende/administrative Aufrufe:
  `RegisterConsumer`, `AcknowledgeConsumer`, `RemoveConsumer`,
  `EnableTable`, `DisableTable`, `RunRetention`. Ein `admin`-Token deckt
  implizit auch die Leserechte ab (dieselbe Hierarchie wie zwischen
  `cdc_admin` und `cdc_reader`, `ADR-0047`).

Ein Aufruf ohne `Authorization: Bearer <token>`-Header oder mit einem
Token, das keiner konfigurierten Klasse entspricht, endet mit `401`; ein
gültiges `reader`-Token gegen einen schreibenden Endpunkt endet mit `403`
— beides sichtbar, nicht still verworfen.

### Teilfrage 4 — Adapter-Platzierung und Schichten-Konformität

| Option | Pro | Contra |
|---|---|---|
| **A — neuer Driving Adapter `internal/adapters/driving/http/` (gewählt)** | folgt dem bestehenden Muster (`internal/adapters/driving/replication`); `ARC-005` nennt HTTP/gRPC bereits als vorgesehene Driving-Adapter-Technologie; `.a-check.yml`s bestehender Glob `adapters: ["internal/adapters/**"]` erfasst das neue Verzeichnis bereits — **keine Änderung an `.a-check.yml` nötig** (geprüft: der Glob ist rekursiv, kein neuer Hexagon-Schichten-Edge) | keine wesentlichen |
| B — HTTP-Handler direkt in `internal/bootstrap`/`cmd/pg-change-feed` | kein neues Verzeichnis | verstößt gegen `ARC-007` (Bootstrap ist Composition Root, keine Adapter-Logik); vermischt Verdrahtung und Protokoll-Handling, schwerer isoliert testbar |
| C — neues Top-Level-Verzeichnis außerhalb `internal/adapters/` (z. B. `internal/api/`) | organisatorisch „prominenter" | `.a-check.yml`s Glob würde das Verzeichnis nicht als `adapters`-Layer erkennen und bräuchte eine Erweiterung ohne Mehrwert gegenüber der bestehenden Konvention |

Der Adapter importiert ausschließlich Inbound Ports
(`internal/application/port/inbound/`) und Domain-Typen zur
Request-/Response-Übersetzung — keine Driven-Adapter-Interna, keine
Application-Interna (`ARC-002`/`ARC-003`/`ARC-005`-Constraint, bereits
maschinell durch `a-check` durchgesetzt).

### Testabdeckung (Erwartung, keine abschließende Festlegung)

Zwei Nachweisformen sind für den umsetzenden Slice-Schnitt verbindlich zu
benennen, ihre genaue Ausgestaltung ist Implementer-Detail:

- **Beispiel-Client** als Wegwerf-Werkzeug außerhalb der
  Produktionsschichten, analog zu `tools/harness/natssub/` — Vorschlag:
  `tools/harness/httpclient/` (von `a-check` bewusst ausgeschlossen, kein
  Bestandteil von `internal/`).
- **E2E-Beleg** als Erweiterung von `make test-integration`/
  `tools/harness/run-integration-tests.sh`: ein realer Rundlauf über die
  neue API gegen den laufenden Feed-Container, analog zum bestehenden
  `docker exec`-Rundlauf für CLI (Register-/Acknowledge-Consumer,
  `LH-FA-CON-*`) — diesmal über einen echten HTTP-Request statt
  `docker exec`, was eine Port-Exposition in `compose.yaml` voraussetzt
  (neue Umgebungsvariable nach dem bestehenden `CDC_*`-Schema, z. B.
  `CDC_HTTP_ADDR`).

## Konsequenzen

- Positiv: `LH-FA-SST-006` wird mit minimalem neuen Abhängigkeits-
  Fußabdruck erfüllbar — ein `net/http`-Adapter, zwei Umgebungsvariablen
  für Tokens, keine neue Fremd-Bibliothek.
- Positiv: Das Least-Privilege-Muster aus `ADR-0047` setzt sich auf der
  API-Ebene konsistent fort (zwei Rechtsklassen statt einer
  Alles-oder-nichts-Grenze).
- Positiv: Kein Eingriff in die bestehende Schichtentrennung nötig —
  `.a-check.yml` bleibt unverändert, weil der neue Adapter vollständig im
  bestehenden `adapters`-Glob liegt.
- Negativ: Changes-Lesen (`LH-FA-REA-*`) und Diagnose/Health bleiben nach
  dieser ADR weiterhin CLI-/SQL-exklusiv — ein Netzwerk-Consumer, der genau
  das braucht, muss auf einen Folge-Slice warten.
- Negativ: Statische Bearer-Tokens ohne Ablauf-/Widerrufsmechanismus sind
  ein Betriebsrisiko bei Kompromittierung (Rotation erfordert einen
  Neustart mit neuer Umgebungsvariable); TLS-Terminierung vor der API ist
  Betreiber-Pflicht, nicht durch diese ADR abgedeckt.
- Folgepflicht: Der umsetzende Slice-Schnitt implementiert Adapter,
  Token-Middleware und die Verdrahtung (`internal/bootstrap`:
  `CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`, No-Op bei
  fehlender Adresse — additiv, kein Breaking Change für bestehende
  Deployments, analog zu `ADR-0055` Punkt 5).
- Folgepflicht: `spec/pflichtenheft.md` erhält einen neuen `SPEC-*`-Eintrag
  für die konkrete Endpunkt-/Methoden-/JSON-Schema-Zuordnung
  (`LH-FA-SST-006`s Out-of-Scope weist das dorthin) — Gegenstand des
  umsetzenden Slices, nicht dieser ADR.
- Folgepflicht: `spec/architecture.md` braucht **keine** Änderung — `ARC-005`
  nennt HTTP/gRPC bereits generisch als Driving-Adapter-Technologie; diese
  ADR macht die Wahl konkret, ohne die Sicht selbst zu ändern.
- Folgepflicht: Erweiterung um Changes-Lesen oder Diagnose/Health braucht
  je eine eigene Port-Design-Entscheidung (neuer Inbound Port bzw.
  Anwendungsschicht-Abstraktion um die heute freien Funktionen
  `Diagnose`/`Healthcheck`) und damit eine eigene Folge-ADR oder einen
  Architect-Zug vor der Implementierung — nicht Gegenstand dieser ADR.
- Folgepflicht (Slice-Schnitt-Empfehlung, endgültiger Schnitt liegt bei der
  Planner-Rolle):
  1. **Slice A** — Adapter-Grundgerüst (`internal/adapters/driving/http/`),
     Token-Middleware (401/403-Pfade), Bootstrap-Verdrahtung, eine erste
     Fähigkeit end-to-end (z. B. `RegisterConsumer`) samt Unit-Tests.
  2. **Slice B** — restliche Port-gedeckte Fähigkeiten
     (Acknowledge/Position/Remove-Consumer, Enable/Disable/Status/
     List-Table, Retention-Lauf) über denselben Adapter, einheitliches
     Fehler-Mapping (400/401/403/404/500), Unit-Tests.
  3. **Slice C** — Beispiel-Client (`tools/harness/httpclient/`) und
     E2E-Erweiterung (`compose.yaml`-Port-Exposition,
     `run-integration-tests.sh`-Rundlauf über echten HTTP-Request).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (neuer `internal/adapters/driving/http`-Adapter, Whitebox) | Ein Aufruf ohne oder mit unbekanntem Bearer-Token endet mit `401`; ein gültiges `reader`-Token gegen einen schreibenden Endpunkt endet mit `403`; ein gültiges `admin`-Token erreicht sowohl lesende als auch schreibende Endpunkte | `make test` |
| `.a-check` | unverändert — der neue Adapter liegt vollständig im bestehenden Glob `adapters: ["internal/adapters/**"]`, kein neuer Hexagon-Schichten-Edge nötig | `make a-check` |
| `make test-integration` | realer HTTP-Rundlauf gegen den laufenden Feed-Container über mindestens eine Fähigkeit je Token-Klasse — Umsetzung Gegenstand von Slice C | `make test-integration` |

## Re-Evaluierungs-Trigger

An Slice C gebunden — sobald der dort gelieferte E2E-Beleg in `done/`
liegt: Erweitert sich der beobachtete Consumer-Bedarf um typed Contracts
oder Streaming (z. B. ein konkreter gRPC-Consumer benennt sich), braucht
die Protokoll-Wahl aus Teilfrage 1 eine eigene Folge-ADR, die gRPC (Option
B) neu bewertet — nicht als Korrektur, sondern als eigene Entscheidung,
weil sich die Prämisse (kein bekannter gRPC-Consumer) geändert hätte.
Ebenso: Wird ein Bedarf an Changes-Lesen oder Diagnose/Health über die API
konkret benannt, löst das die unter Konsequenzen benannte Folge-ADR aus.
Sonst permanent — die Wahl HTTP/JSON mit Token-Authn gilt unabhängig vom
Zeitpunkt der Umsetzung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Architect-Entscheidung nach explizitem Auftraggeber-Bedarf; Re-Evaluierungs-Trigger aus `ADR-0020` eingetreten (Supersedes `ADR-0020`) | [`ADR-0020`](0020-http-grpc-optional.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0057` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
