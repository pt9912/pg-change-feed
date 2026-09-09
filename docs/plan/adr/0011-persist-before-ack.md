# ADR-0011: Persist-before-ACK

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-REL-001`](../../../spec/lastenheft.md),
[`LH-FA-RET-001`](../../../spec/lastenheft.md)

**Schärft:** [`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Quelle wartet auf Bestätigungen: Bestätigt der Capture-Pfad zu früh und
stürzt danach ab, sind Changes still verloren (LH-QA-REL-001). Zu späte
Bestätigung kostet dagegen nur Wiederholung. Die beiden Fehlerklassen sind
nicht gleich schwer — das legt die Ordnung fest.

## Entscheidung

Wir wählen **Persistenz vor Bestätigung**: Die WAL-Position wird erst
bestätigt, wenn alle abhängigen Changes dauerhaft gespeichert sind.

    Receive -> Decode -> Persist -> COMMIT Store -> ACK Source

Formal:

    ACK(position) => durable(all changes <= position)

Wiederholung wird gegenüber Datenverlust bevorzugt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Sofortiges ACK (Throughput zuerst) | Höchster Durchsatz, geringster WAL-Rückstand | Crash nach ACK, vor Persistenz: stille Datenlücke — verletzt LH-QA-REL-001 direkt |
| B — Periodisches Sammel-ACK | Weniger ACK-Verkehr | Zeitfenster für Verlust bleibt; Grenze „welche Changes sind sicher?" wird unbestimmt |
| **C — Persist-before-ACK** | Keine stille Lücke möglich; Restrisiko ist doppelt verarbeitete, deduplizierbare Changes | Mehr Write-Ahead-Rückstand; Durchsatz hängt am Speicher-Commit |

## Konsequenzen

- Positiv: Neustart setzt ohne Datenlücke fort (LH-FA-RET-001); der
  kritische Test „Crash nach Persistenz vor ACK erzeugt keine Lücke" ist
  direkt aus der Invariante ableitbar.
- Negativ: Idempotente Persistenz (Deduplizierung wiederholter WAL-Daten)
  ist Pflicht, kein Optimismus.
- Folgepflicht: Rollen-Trennung Stream/ACK (ADR-0007); Invariantenliste
  (ADR-0029); Abnahme im MVP-Integrationstest.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Invariante trägt die Vertragszusage LH-QA-REL-001.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).