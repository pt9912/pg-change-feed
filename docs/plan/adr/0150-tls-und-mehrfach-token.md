# ADR-0150: TLS der HTTP- und gRPC-Schnittstellen und mehrere API-Token je Klasse

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-04

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-FA-SST-011`](../../../spec/lastenheft.md) (TLS),
[`LH-FA-SST-012`](../../../spec/lastenheft.md) (Token-Wechsel),
[`LH-FA-SST-006`](../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../spec/lastenheft.md),
[`LH-QA-SEC-002`](../../../spec/lastenheft.md),
[ADR-0057](0057-http-grpc-api.md),
[ADR-0083](0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.5, §3.12.

**Schärft:** [`SPEC-034`](../../../spec/pflichtenheft.md#spec-034--tls-der-http--und-grpc-schnittstellen),
[`SPEC-035`](../../../spec/pflichtenheft.md#spec-035--mehrere-api-token-je-klasse),
[`SPEC-016`](../../../spec/pflichtenheft.md#spec-016--konfigurationsdatei-cdc_config_file)
(Schlüssel `tls_cert_file`, `tls_key_file`, ausgeschlossene Klasse),
[`SPEC-018`](../../../spec/pflichtenheft.md#spec-018--http-api-endpunkte-und-token-header-form).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die HTTP- und gRPC-Schnittstellen laufen im Klartext, mit je einem Token je
Klasse (`CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`). Ein Token lässt sich nur
mit einem Zeitfenster wechseln, in dem Clients abgewiesen werden.

**Gemessen** (`sed -n 30,50p internal/adapters/driving/http/middleware.go`,
`sed -n 45,65p internal/adapters/driving/grpc/interceptor.go`, 2026-10-04): der
HTTP-Adapter prüft mit `==` (`classifyToken`), der gRPC-Adapter trägt eine zweite
Fassung derselben Funktion. `.a-check.yml` führt `adapters` als eine Schicht ohne
Kante von Adapter zu Adapter; beide Fassungen sind deshalb zusammen zu ändern.

Der Auftraggeber hat TLS plus Token-Rotation gewählt, kein mTLS.

## Entscheidung

Wir wählen **ein gemeinsames TLS-Paar für HTTP und gRPC und zusätzliche
Token-Listen je Klasse, ohne Neuladen und ohne mTLS**.

1. **TLS:** `CDC_TLS_CERT_FILE`/`CDC_TLS_KEY_FILE` (Datei-Schlüssel
   `tls_cert_file`, `tls_key_file`), ein Paar für HTTP (einschließlich SSE) und
   gRPC. Nur eines gesetzt, nicht ladbar oder nicht passend: Klasse
   `configuration`, kein Start. Mindestversion TLS 1.2, explizit gesetzt. Kein
   Klartext auf derselben Adresse; kein mTLS; kein Neuladen. Ein Bind-Fehler bleibt
   wie bisher ein Log-Eintrag. NATS-TLS bleibt Sache der NATS-Konfiguration.
   Die Pfade sind keine Zugangsdaten und dürfen in die Konfigurationsdatei.
2. **Token:** neue env-exklusive Variablen `CDC_API_TOKENS_READER`/
   `CDC_API_TOKENS_ADMIN` (Plural, kommagetrennt). Die wirksame Menge einer Klasse
   ist die Vereinigung mit dem Singular; der Singular bleibt ohne Listen-Syntax,
   damit ein bestehendes Token mit Komma nicht zerschnitten wird. Ein leeres
   Element oder Whitespace im Element ist ein Startfehler der Klasse
   `configuration`; ein leeres Token matcht nie; bei gleichem Wert in beiden Klassen
   gewinnt `admin`.
3. **Prüfung:** zeitkonstant über alle Token beider Klassen (`crypto/subtle`) —
   *Erwartung, nicht erprobt*. Beide Fassungen (HTTP und gRPC) sind zusammen zu
   ändern.
4. **Wechsel:** zwei Neustarts (neues Token ergänzen, Clients umstellen, altes
   entfernen). Ein Neuladen zur Laufzeit gibt es nicht.
5. **Abgrenzung:** TLS-Optionen der SDKs (C#, Kotlin, Python) sind eine eigene,
   nachgelagerte Anforderung und hier nicht gelöst; das NATS-Stream-Token ist
   außerhalb.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, TLS über einen vorgeschalteten Proxy | kein Code | gRPC und SSE hinter einem Proxy brauchen eigene Konfiguration; kein Beleg im Produkt; Token-Wechsel bleibt offen |
| B — mTLS statt Token | stärkere Authentisierung | Zertifikatsverwaltung je Client; vom Auftraggeber abgewählt |
| C — Neuladen von Zertifikat und Token zur Laufzeit | kein Neustart | Dateiüberwachung und Nebenläufigkeit im Adapter; für den Zuschnitt nicht gefordert |
| **D — ein TLS-Paar und Token-Listen mit Neustart (gewählt)** | kleinster Zuschnitt, Bestandskonfiguration unverändert | Wechsel braucht zwei Neustarts; ein Zertifikat gilt für beide Schnittstellen |

## Konsequenzen

- Positiv: verschlüsselter Betrieb ohne Proxy; ein Token ist ohne Abweisung
  wechselbar; Bestandskonfiguration läuft unverändert.
- Negativ: Zertifikat und Token wirken erst nach einem Neustart; ein SDK-Client
  kann TLS vorerst nicht selbst konfigurieren.
- Folgepflicht: Slice(s) für TLS und Token (`api-tls`, `api-token-mehrfach`,
  Code-Namen ohne Link), die beiden `classifyToken`-Fassungen im selben Zug; neue
  Meldungscodes der Klasse `configuration` in der Code-Tabelle; Handbuch
  (kennungsfrei) um Konfiguration und den Wechselablauf ergänzen; eine Folge-Anforderung
  für SDK-TLS-Optionen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test | nur ein TLS-Pfad gesetzt, unlesbar oder nicht passend: Klasse `configuration`, kein Start — *erwartet*, nicht gefahren | `make test` |
| Go-Test | Token-Liste mit leerem Element oder Whitespace: Startfehler; derselbe Wert in beiden Klassen: `admin`; leeres Token nie gültig — *erwartet*, nicht gefahren | `make test` |
| Integrationstest (Compose) | mit TLS: Klartext-Zugriff auf dieselbe Adresse wird nicht bedient; beide Token einer Klasse gelten; nach Entfernen des alten Tokens `401` — *erwartet*, nicht gefahren | `make test-integration` |

## Re-Evaluierungs-Trigger

Ein Betreiber, der das Zertifikat ohne Neustart erneuern muss, oder die Forderung
nach gegenseitiger Authentisierung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Accepted — Anlass: Auftraggeber-Entscheidung „TLS und Token-Rotation, kein mTLS“ | `LH-FA-SST-011`, `LH-FA-SST-012` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0150` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
