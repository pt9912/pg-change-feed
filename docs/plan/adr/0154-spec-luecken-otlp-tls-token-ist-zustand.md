# ADR-0154: Spec-Lücken zu OTLP-Export, TLS und API-Token — der Ist-Zustand wird festgeschrieben

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-04

**Autor:** pt9912 (Architect-Rolle, Modul 8; Anlass: der Review-Report
`review-slice-otlp-metrik-export` (Befund F-4), der Verifikations-Report
`verifikation-slice-otlp-metrik-export` und die Closure-Notizen der Slices
`otlp-metrik-export`, `tls-http-grpc-server`, `api-token-mehrfach-konfiguration`)

**Bezug:** [`LH-FA-SST-010`](../../../spec/lastenheft.md),
[`LH-FA-SST-011`](../../../spec/lastenheft.md),
[`LH-FA-SST-012`](../../../spec/lastenheft.md),
[`ADR-0149`](0149-otlp-metrik-export-mechanismus.md),
[`ADR-0150`](0150-tls-und-mehrfach-token.md),
[`ADR-0152`](0152-zugangsdaten-klasse-elf-schluessel.md),
[`ADR-0153`](0153-otlp-einheit-consumer-lag-byte.md),
[`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.5, §3.12.

**Schärft:** [`SPEC-033`](../../../spec/pflichtenheft.md#spec-033--otlp-metrik-export-drahtform-konfiguration-verhalten),
[`SPEC-034`](../../../spec/pflichtenheft.md#spec-034--tls-der-http--und-grpc-schnittstellen),
[`SPEC-035`](../../../spec/pflichtenheft.md#spec-035--mehrere-api-token-je-klasse).
`ADR-0149` und `ADR-0150` bleiben unverändert: nichts hier widerspricht ihrer
§Entscheidung, jeder Satz ergänzt, wo sie schwiegen.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die drei Slices haben Verhalten umgesetzt, das die Spec nicht nennt. Jede Lücke
stand als Adresse mit der Frist „vor dem Release“ in einer Closure-Notiz. Diese
ADR legt fest, welche Lücken in die Spec gehören, und schreibt den **Ist-Zustand**
fest — es entsteht kein neues Verhalten. Die Herkunft jeder Aussage:

- **gelesen** — Quelltext am Stand `b03ff591` (2026-10-04): `internal/bootstrap/otlp.go`
  (`applyOTLP`, `validateOTLPEndpoint`, `parseOTLPHeaders`, `parseOTLPInterval`),
  `internal/bootstrap/wiring.go` (`parseAPITokenList`, `applyAPITokens`),
  `internal/bootstrap/tls.go` (`newTLSConfig`),
  `internal/application/usecase/exportmetrics/service.go` (`Export`),
  `internal/adapters/driven/otlpexport/exporter.go` (`New`, `specs`);
- **gemessen** — ein Wegwerf-Go-Programm im Toolchain-Image (`--network none`,
  Server und Client im selben Prozess auf Loopback), gelaufen am 2026-10-04: (1)
  eine URL `http://user:pass@127.0.0.1:<port>` mit `http.Client.Post` liefert dem
  Server `Authorization: Basic dXNlcjpwYXNz`; (2) setzt die Anfrage zusätzlich
  `Authorization: Bearer abc`, sieht der Server `Bearer abc` — der explizite
  Header gewinnt; (3) der Fehlertext von `net/http` bei einer nicht erreichbaren
  URL mit Benutzerteil trägt den Benutzernamen und `***` statt des Passworts;
- **übernommen** — Aussagen der Review- und Verifikations-Berichte, jeweils
  benannt.

## Entscheidung

Wir schreiben die folgenden Sätze in die Spec (Wortlaut in der Anwende-Anleitung).

### SPEC-033

1. **Header-Regeln (`CDC_OTLP_HEADERS`).** Das Komma trennt die Elemente, das
   **erste** `=` Schlüssel und Wert (der Wert darf weitere `=` tragen). Der
   Schlüssel ist ein gültiger HTTP-Header-Name (RFC-9110-`token`, **Menge der
   zulässigen Zeichen:** Ziffern, ASCII-Buchstaben und ``!#$%&'*+-.^_`|~``; Leerraum
   im Schlüssel ist ein Fehler); der Wert trägt kein Steuerzeichen (Menge:
   `U+0000` bis `U+001F` ohne Tabulator, und `U+007F`), Leerraum im Wert bleibt
   erhalten. Eine leere Zeichenkette liefert keine Header. Ein Verstoß endet den
   Start mit der Klasse `configuration`; der Fehlertext nennt die Position des
   Elements, nie dessen Inhalt. *Gelesen:* `parseOTLPHeaders` und
   `isHeaderTokenRune`/`isHeaderControlRune`.
2. **Header ohne Endpunkt.** `CDC_OTLP_HEADERS` gesetzt und `CDC_OTLP_ENDPOINT`
   nicht gesetzt endet den Start mit der Klasse `configuration` (ein Wert, der
   ohne Endpunkt nie wirkt, ist keine stille Auslassung). Der Takt ohne
   Endpunkt ist zulässig. Eine leer gesetzte Umgebungsvariable gilt als nicht
   gesetzt. *Gelesen:* `applyOTLP`.
3. **Benutzerteil der Endpunkt-URL — wird zugelassen.** `http://user:pass@host`
   wird nicht abgelehnt; der Benutzerteil geht als `Authorization: Basic` hinaus,
   **außer** `CDC_OTLP_HEADERS` setzt selbst einen `Authorization`-Header (der
   gewinnt). Weder ein Log-Eintrag noch ein Fehlertext des Exports nennt die URL,
   den Benutzerteil oder einen Header-Wert; die Zugangsdaten gehören wie bei
   `nats_url` zur Zugangsdaten-Klasse (`ADR-0152`: `otlp_endpoint` und
   `otlp_headers` sind Mitglieder, env-exklusiv). *Gemessen:* Basic und Vorrang
   (oben). *Übernommen* aus dem Review (INFO-1, Mutation an einer Kopie, Go-Test
   des Pakets `otlpexport`, rot): kein Header- und kein Benutzerteil-Leck im
   Adapter; *gelesen:* der Start-Log des Exports nennt allein den Takt
   (`internal/bootstrap/otlp.go` Zeile 292), die Konfigurationsfehler nennen
   Variable und Position (`otlpConfigError`).
4. **Erster Versuch.** Der erste Versuch folgt einen vollen Takt nach dem Start,
   nicht beim Start. *Gelesen:* `time.NewTicker` ohne Sofortlauf.
5. **Frist je Versuch** gilt für Lesen der Sicht und Übertragen **zusammen**
   (ein Kontext je Versuch). *Gelesen:* `Service.Export`.
6. **Gemeinsamer Fehlerzustand.** Lesen und Übertragen teilen den Zustand der
   gedrosselten Warnung; die Wiederaufnahme ist ein Versuch, in dem beide Stufen
   gelungen sind. *Übernommen* aus dem Plan §3 und dem Review (die Mutation der
   Drosselung färbt `TestReadAndTransmitShareOneFailureState` rot).
7. **Weiterleitungen** (`3xx`) werden nicht verfolgt; sie sind ein Fehlschlag, damit
   kein Header einen anderen Host erreicht. *Übernommen:* Review (Mutation
   `CheckRedirect` entfernt, `TestExportDoesNotFollowRedirects` rot).
8. **Zahlenform.** Eine Kennzahl der **zehn** der Tabelle geht als ganze Zahl
   hinaus, wenn ihr Wert der Sicht eine ist **und** ihre Einheit nicht `s` ist;
   die zwei Kennzahlen mit Einheit `s` (`cdc_oldest_change_age_seconds`,
   `cdc_capture_lag`) gehen immer als Fließkommazahl. Eine Zeile der Sicht mit
   einem Namen außerhalb der Tabelle wird nicht übertragen. *Gelesen:*
   `metricSpec.float`, Zeilen 198 bis 201 von `exporter.go`.
9. **Datei-Feld `otlp_interval`:** dieselbe Prüfung (5 bis 3600) wie die Variable,
   die Variable geht vor. Wie der Wert im YAML steht (Zahl oder Zeichenkette), ist
   Sache der Umsetzung und nicht Teil des Vertrags.

### SPEC-034

10. **„Ladbar“ ist die ganze Prüfung beim Start:** beide Dateien sind lesbares
    PEM, der Schlüssel gehört zum Zertifikat (Laden als Paar). Der Start prüft
    **nicht** den Gültigkeitszeitraum (Ablauf, Beginn), **nicht** die Namen
    (Hostname, SAN) und **nicht** die Kette zu einer Wurzel. Ein abgelaufenes
    oder für den Namen unpassendes Zertifikat startet den Server; der Client
    lehnt es beim Verbinden ab. *Erprobt* für den Ablauf: Go-Test
    `TestNewTLSConfigPruftAblaufUndNamenNicht` (Paket `bootstrap`, ein Zertifikat
    mit Ende vor einer Minute lädt ohne Fehler — Testtext gelesen, Lauf **übernommen**
    aus dem Review); für Name und Kette *hergeleitet* (`newTLSConfig` ruft allein
    `tls.LoadX509KeyPair`, gelesen).

### SPEC-035

11. **Leere Variable.** Eine gesetzt-leere Plural-Variable gilt als ungesetzt (liefert
    keine Token); ebenso eine leere Singular-Variable. *Gelesen:*
    `parseAPITokenList`, `applyAPITokens`.
12. **Komma im Token:** in der Liste nicht ausdrückbar (das Komma trennt); ein Token
    mit Komma steht nur im Singular.
13. **Leerraum:** der Singular wird nicht auf Leerraum geprüft; jedes Element der
    zwei Plural-Variablen wird geprüft (Menge: vier Variablen, zwei davon
    prüfend). *Gelesen:* `applyAPITokens` weist den Singular ohne Prüfung zu;
    `git grep -n 'IsSpace\|TrimSpace' -- internal/bootstrap internal/adapters/driving`
    (ohne Tests) trifft die Plural-Zerlegung, die OTLP-Header und zwei Stellen der
    Tabellen-Aktivierung (`CDC_TABLES`), keine Prüfung eines Singular-Tokens.

### Lastenheft

14. **„ungültig“ in `LH-FA-SST-011` Negative** meint ein Zertifikat, das sich nicht als
    PEM-Zertifikat samt passendem Schlüssel laden lässt; die Prüfung von
    Gültigkeitszeitraum, Namen und Kette beim Start ist **Out-of-Scope**. Das ist
    eine **Lockerung gegenüber der weitesten Lesart** des Lastenheft-Worts
    („ungültig“ könnte „abgelaufen“ einschließen) und damit das dritte Verdikt des
    Konflikt-Pfads: die Lockerung ist tragfähig (der Client lehnt ab; der Betreiber
    sieht das Ablaufdatum an seiner Zertifikatsverwaltung; Beschaffung und
    Erneuerung sind bereits Out-of-Scope), aber undokumentiert — sie wird im
    Lastenheft nachgezogen (Anwende-Anleitung, Patch-Version `0.15.1`). Weil das
    Lastenheft vertraglich bindet, braucht dieser Punkt die Bestätigung des
    Auftraggebers; ohne sie bleibt als Alternative eine Prüfung des Ablaufs beim
    Start im Code (eigener Slice).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Verhalten nur im Handbuch | kein Aufwand | die Spec (Rang 2) gewinnt gegen das Handbuch (Rang 6): sie sagt weniger als der Code tut; ein späterer Umbauer liest nur die Spec |
| B — Benutzerteil der URL ablehnen (`configuration`, `PCF-E2011`) | keine Zugangsdaten in einer URL | Code, Test und Handbuch müssen vor dem Release mitziehen; der Gewinn ist klein, weil `CDC_OTLP_ENDPOINT` ohnehin zur Zugangsdaten-Klasse gehört und `CDC_OTLP_HEADERS` denselben Schutz trägt; ein Betreiber mit einem Basic-Auth-Empfänger verlöre einen gangbaren Weg |
| **C — Ist-Zustand festschreiben, Benutzerteil zulassen (gewählt)** | keine Code-Änderung; Spec und Code sind gleich; das Leck-Risiko (Log, Fehlertext) ist im Review an Mutationen belegt (übernommen) | die Zugangsdaten stehen in der Prozess-Umgebung — wie bei jedem Header-Wert und jedem Token dieser Klasse |
| D — Ablauf des TLS-Zertifikats beim Start prüfen | schlösse die weite Lesart von „ungültig“ | neues Verhalten; Zeitquelle, Fehlercode und Handbuch; ein abgelaufenes Zertifikat zu starten ist beobachtbar (Client lehnt ab) |

## Konsequenzen

- Positiv: Spec und Code sagen dasselbe; ein Umbauer findet die Grenzen im Rang 2.
  Kein Code-Eingriff.
- Negativ: ein Zertifikat mit Ablauf in der Vergangenheit startet den Server
  (benannte Grenze, das Handbuch nennt sie bereits); ein Singular-Token mit
  Leerraum wird nicht beanstandet.
- Folgepflicht: das Handbuch nennt den Benutzerteil der Endpunkt-URL noch nicht
  (Basic, Vorrang eines `Authorization`-Headers aus `CDC_OTLP_HEADERS`) — Nachzug
  durch den Slice `otlp-metrik-export-e2e` vor dem Tag; die Aussage zu den
  Weiterleitungen steht dort bereits. Ein Test, der den Benutzerteil fährt (der
  Fake-Empfänger liest `Authorization: Basic …`), gehört in das Paket
  `otlpexport` — *erwartet*, nicht gefahren.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test `otlpexport` (Vorschlag) | Endpunkt mit Benutzerteil: der Empfänger-Fake liest `Authorization: Basic <base64(user:pass)>`; mit gesetztem `Authorization`-Header aus `CDC_OTLP_HEADERS` liest er diesen — *erwartet*, nicht gefahren; die Aussage ist am Wegwerf-Programm oben gemessen, nicht an einem Test des Repos | `make test` |
| Go-Test `bootstrap` (bestehend) | `TestNewTLSConfigPruftAblaufUndNamenNicht` trägt Satz 10 für den Ablauf: ein Zertifikat mit Ablauf in der Vergangenheit lädt — Lauf übernommen; Name und Kette: hergeleitet | `make test` |
| Handbuch-Gate | `make handbuch-public-doc-check` für den Handbuch-Nachzug — *erwartet* | `make handbuch-public-doc-check` |

## Re-Evaluierungs-Trigger

Ein Betreiber meldet einen Fall, in dem ein abgelaufenes Zertifikat den Betrieb
beeinträchtigt hat (dann Option D als eigener Slice), oder ein Leck des
Benutzerteils in einer Log-Zeile oder einem Fehlertext (dann Option B).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Accepted (Punkt 14 unter Vorbehalt der Bestätigung des Auftraggebers, siehe dort) | Closure-Notizen der drei Slices, Review-Report `review-slice-otlp-metrik-export` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
