# ADR-0020: HTTP/gRPC optional

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-SST-005`](../../../spec/lastenheft.md)

**Schärft:** [`ARC-005`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Lastenheft fordert nur, dass eine spätere HTTP-/gRPC-API möglich ist,
ohne das interne CDC-Modell grundlegend zu verändern ([`LH-FA-SST-005`](../../../spec/lastenheft.md)).
Ein konkreter API-Consumer-Bedarf liegt zum Entscheidungszeitpunkt nicht
vor.

## Entscheidung

HTTP/gRPC können später als zusätzliche Driving Adapters ergänzt werden;
sie gehören nicht zwingend zum MVP.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — HTTP/gRPC ab MVP | API-Consumer werden früh bedient | MVP-Wachstum ohne beobachtbaren Bedarf; [`LH-FA-SST-005`](../../../spec/lastenheft.md) fordert keine API |
| B — HTTP/gRPC nie | geringste Komplexität | schließt die von [`LH-FA-SST-005`](../../../spec/lastenheft.md) geforderte Erweiterungsfähigkeit aus |
| **C — HTTP/gRPC optional, später** | MVP fokussiert; Architektur lädt die Erweiterung ([`ARC-005`](../../../spec/architecture.md)-Liste offen) | API-Anfragen müssen auf einen späteren Stand warten |

## Konsequenzen

- Positiv: MVP-Umfang bleibt auf den Nachweis der CDC-Abstraktion gerichtet.
- Negativ: ein früh auftretender API-Bedarf braucht einen neuen
  Entscheidungsanlass.
- Folgepflicht: —

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbarer Bedarf eines API-Consumers — sichtbar als Anforderung im
Lastenheft-Change oder als Eintrag im Beobachtungs-Register; nicht an einem
Datum festgemacht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0020` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).