# Welle 16: HTTP/JSON-API mit Token-Authn (`LH-FA-SST-006`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-16-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-14.

---

## 1. Welle-Ziel

Ein Netzwerk-Client erreicht die bereits Port-gedeckten Fähigkeiten
(Consumer-Registrierung/-Bestätigung/-Position/-Entfernung,
Tabellen-Aktivierung/-Deaktivierung/-Status/-Liste, Retention-Lauf) über
einen neuen HTTP/JSON-Adapter (`internal/adapters/driving/http/`), geschützt
durch statische Bearer-Tokens in zwei Rechtsklassen
(`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`) — `LH-FA-SST-006`s Happy-Path/
Boundary/Negative-Akzeptanzkriterien werden über einen real ausgeführten
End-to-End-Rundlauf gegen den laufenden Feed-Container belegt, nicht nur
über Unit-Tests. `ADR-0057` (Accepted) trifft die Architektur-Entscheidung
bereits vollständig; diese Welle setzt sie in drei Slices um.

## 2. Trigger (Welle startet)

- `ADR-0057` liegt Accepted vor (bereits erfüllt — Architect-Entscheidung
  vom 2026-09-14, Supersedes `ADR-0020`).

## 3. Closure-Trigger (Welle schließt)

- `slice-059`, `slice-060`, `slice-061` liegen in `done/`.
- `make gates` grün.
- Ein real belegter End-to-End-Rundlauf über den Beispiel-Client
  (`tools/harness/httpclient/`) gegen mindestens eine Fähigkeit je
  Token-Klasse (reader und admin), gegen den laufenden Feed-Container in
  der Compose-Umgebung — das *Mehr* gegenüber den einzelnen Slice-DoDs:
  keiner der drei Slices allein belegt den vollen Rundlauf über alle drei
  Liefer-Schichten (Adapter, restliche Fähigkeiten, Client+E2E-Verdrahtung).
- Closure-Notiz in `welle-16-results.md`.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| slice-059 | HTTP-API-Adapter-Grundgerüst, Token-Middleware, `RegisterConsumer` | [LH-FA-SST-006](../../../spec/lastenheft.md), [LH-FA-CON-001](../../../spec/lastenheft.md) |
| slice-060 | HTTP-API — restliche Port-gedeckte Fähigkeiten | [LH-FA-SST-006](../../../spec/lastenheft.md), `LH-FA-CON-*`, `LH-FA-CFG-*`, `LH-FA-RET-*` |
| slice-061 | HTTP-API — Beispiel-Client und E2E-Rundlauf | [LH-FA-SST-006](../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

- `slice-060` baut auf `slice-059` auf (derselbe Adapter, dieselbe
  Token-Middleware, dasselbe Fehler-Mapping-Muster).
- `slice-061` baut auf `slice-059` und `slice-060` auf (der Beispiel-Client
  und der E2E-Rundlauf brauchen den vollständigen Fähigkeitsumfang, um je
  eine Fähigkeit pro Token-Klasse real aufzurufen).
- Blockiert: keine andere offene Welle (`Nächste Wellen` ist derzeit leer).
- Wird blockiert von: keiner (`ADR-0057` liegt bereits Accepted vor).

## 6. Out-of-Scope für diese Welle

- **Changes-Lesen über die API (`LH-FA-REA-*`)** — `ADR-0057` schließt das
  in Teilfrage 2 (Option C) ausdrücklich aus: `cdc.changes` ist nach
  `ADR-0046` bewusst ein SQL-View-Direktzugriff ohne Anwendungsdienst; eine
  neue Port-Abstraktion dafür ist eine eigene architektonische Entscheidung
  (Folge-ADR, siehe `ADR-0057` §Konsequenzen/Folgepflicht).
- **Diagnose/Health über die API** — `ADR-0057` schließt das in Teilfrage 2
  (Option D) ausdrücklich aus: `Diagnose`/`Healthcheck` sind heute freie
  Funktionen ohne Anwendungsschicht-Abstraktion; eine Erweiterung braucht
  eine eigene Port-Design-Entscheidung (Folge-ADR, siehe `ADR-0057`
  §Konsequenzen/Folgepflicht).
- **gRPC** — `ADR-0057` Teilfrage 1 entscheidet HTTP/JSON; eine gRPC-Ergänzung
  bleibt an den in `ADR-0057` §Re-Evaluierungs-Trigger benannten, bislang
  nicht eingetretenen Bedarf gebunden.
- **Token-Ablauf-/Widerrufsmechanismus, mTLS, OAuth2/JWT** — `ADR-0057`
  Teilfrage 3 wählt statische Bearer-Tokens bewusst als minimalen ersten
  Schritt; Rotation bleibt Betreiber-Pflicht (Neustart mit neuer
  Umgebungsvariable).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [welle-16-results.md](welle-16-results.md), Geschwister im Ruheort `done/`
Zähler: [../observations/](../observations/)`BEO-PGC/`, eine Ebene über dem Ruheort
