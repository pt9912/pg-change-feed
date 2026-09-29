# Welle welle-sdk-grpc-administration-flaeche: SDK-Erweiterung um die gRPC-Verwaltungs-API

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-sdk-grpc-administration-flaeche-results.md`). Der
Zustand ist die Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug (`LH-FA-SST-009` selbst trägt
keinen Meilenstein-Eintrag im Lastenheft).

**Verantwortlich:** —. **Datum:** 2026-09-28.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die drei SDK-Packages ([`PgChangeFeed.Client`](../../../../sdks/csharp/README.md)
(C#), [`pgchangefeed`](../../../../sdks/python/README.md) (Python),
[`pgchangefeed-kotlin`](../../../../sdks/kotlin/pgchangefeed-kotlin/README.md) (Kotlin),
[`LH-FA-SST-009`](../../../../spec/lastenheft.md)) decken die gRPC-Fläche heute
nur mit dem Live-Change-Stream (`StreamChanges`) ab. Vier bereits
umgesetzte, servers­eitig vollständig gepushte ADRs erweitern die
gRPC-Fläche seither um elf Administration-RPCs und einen Stream-Filter:
[`ADR-0130`](../../adr/0130-grpc-verwaltungs-api-neun-rpcs.md) (neun
Verwaltungsfähigkeiten), [`ADR-0131`](../../adr/0131-grpc-readchanges-zehnter-rpc.md)
(`ReadChanges`), [`ADR-0132`](../../adr/0132-diagnose-ueber-inbound-port-http-grpc.md)
(`Diagnose`), [`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md)
(`schema`/`table`-Filter an `StreamChanges`). Alle vier benennen die
SDK-Erweiterung ausdrücklich als eigene, hier nicht umgesetzte Folgepflicht.
Die drei Beispiel-Client-Sprachen (`examples/grpc-client` (Go),
`examples/csharp/grpc-client`, `examples/kotlin/grpc-client`) decken die
volle Fläche bereits ab (implementiert, reviewt, verifiziert, gepusht) und
dienen als direktes fachliches Vorbild für Nachrichtenschema, Rechtsklassen
und Fehlerform.

Diese Welle bringt alle drei SDK-Packages auf denselben Fähigkeitsstand wie
die Beispiel-Clients: ein `PgChangeFeedAdministrationClient` (oder analog
benannt) je Sprache mit allen elf RPCs, dazu die optionalen `schema`/`table`-
Filter-Parameter am bestehenden `PgChangeFeedGrpcClient`.

## 2. Trigger (Welle startet)

- `ADR-0130`, `ADR-0131`, `ADR-0132`, `ADR-0133` sind `Accepted` und
  serverseitig vollständig implementiert, reviewt, verifiziert, gepusht,
  CI grün — erfüllt.
- Die Beispiel-Client-Verb-Matrix (Go/C#/Kotlin) für dieselbe Fläche ist
  implementiert, reviewt, verifiziert, gepusht — erfüllt (fachliches
  Vorbild dieser Welle).

## 3. Closure-Trigger (Welle schließt)

- Alle drei Slices in `done/`.
- `make gates` grün.
- `make sdk-pack-csharp`, `make sdk-pack-python`, `make sdk-pack-kotlin`
  grün für alle drei Packages.
- `docs/user/benutzerhandbuch.md` nennt in beiden gRPC-Abschnitten
  („Zugriff über den gRPC-Change-Stream", „Zugriff über die
  gRPC-Verwaltungs-API") alle drei SDK-Packages als abdeckend, nicht mehr
  als offenen Folge-Schritt.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| slice-sdk-csharp-grpc-administration-flaeche | C#-SDK: `PgChangeFeedAdministrationClient` + Stream-Filter | [`LH-FA-SST-009`](../../../../spec/lastenheft.md) |
| slice-sdk-python-grpc-administration-flaeche | Python-SDK: `PgChangeFeedAdministrationClient` + Stream-Filter | [`LH-FA-SST-009`](../../../../spec/lastenheft.md) |
| slice-sdk-kotlin-grpc-administration-flaeche | Kotlin-SDK: `PgChangeFeedAdministrationClient` + Stream-Filter | [`LH-FA-SST-009`](../../../../spec/lastenheft.md) |

## 5. Abhängigkeiten

- Wird blockiert von: keiner laufenden Welle (Roadmap führt aktuell keine
  weitere offene Welle).
- Blockiert: keine bekannte Folge-Welle.

## 6. Out-of-Scope für diese Welle

- **Neue ADR oder Änderung an `ADR-0130`–`0133`** — diese Welle setzt eine
  bereits getroffene, akzeptierte Entscheidung um, sie überprüft sie nicht.
- **Beispiel-Clients (Go/C#/Kotlin)** — bereits vollständig umgesetzt,
  außerhalb dieser Welle (fachliches Vorbild, nicht Gegenstand).
- **Server-seitige Änderungen** — die vier ADRs sind serverseitig
  abgeschlossen; findet ein SDK-Slice einen Server-Fehler (wie bereits
  einmal bei der Beispiel-Client-Arbeit geschehen, siehe
  [`docs/reviews/review-example-kotlin-grpc-client-verbmatrix.md`](../../../reviews/review-example-kotlin-grpc-client-verbmatrix.md)
  INFO-1), wird er gemeldet, nicht in dieser Welle repariert.
- **Diagnose/Health über HTTP** — bereits eigenständig umgesetzt
  (`ADR-0132`), kein Gegenstand der SDK-Erweiterung dieser Welle über den
  gRPC-Weg hinaus.
- **`GET /diagnose`-SDK-Methode** — dieselbe Fläche wie `Diagnose` über
  gRPC; ob eine HTTP-Diagnose-SDK-Methode folgt, ist Gegenstand einer
  eigenen, hier nicht geplanten Folge-Welle.
- **Breaking Changes an bestehenden SDK-Methoden** (`StreamChangesAsync`/
  `stream_changes`/`streamChanges`) — die neuen Filter-Parameter sind
  additiv/optional, kein Vertragsbruch.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln.

Ergebnis: `done/welle-sdk-grpc-administration-flaeche-results.md`
Zähler: `../observations/README.md`
