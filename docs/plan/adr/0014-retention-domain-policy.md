# ADR-0014: Retention als Domain Policy

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-RET-004`](../../../spec/lastenheft.md)

**Schärft:** [`LH-FA-RET-004.a`](../../../spec/lastenheft.md),
[`ARC-001`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Welche erfassten Changes entfernt werden dürfen, hängt von
Consumer-Positionen, Mindestaufbewahrungszeiten und der
Safe-Watermark-Logik zusammen — fachliche Regeln über dem gespeicherten
Datenbestand. Läge diese Entscheidung im Storage-Adapter oder außerhalb
des Systems, würde Löschwissen an die Speichertechnik gekoppelt bzw.
verließe das System; die Sicherheit der Retention ([`LH-FA-RET-004`](../../../spec/lastenheft.md),
[`LH-FA-RET-005`](../../../spec/lastenheft.md)) wäre nicht mehr domänentestbar.

## Entscheidung

Der Domain Core entscheidet, was sicher gelöscht werden darf; der
Storage-Adapter führt die physische Löschung aus.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Retention im Storage-Adapter | nah am Datenbestand, keine zusätzliche Indirektion | Löschlogik an Speichertechnik gekoppelt, schwer domänentestbar, Risiko unkontrollierter Löschung |
| B — Retention im Betrieb (externes Werkzeug/Cron) | keine Komplexität im Daemon | Regelwissen verlässt das System; kein Zugriff auf die Safe-Watermark; Löschung außerhalb der Systemkontrolle |
| **C — Retention als Domain Policy** | fachliche Regeln domänentestbar; Trennung von Entscheidung und Ausführung | Löschbefehl und Ausführung sind getrennt (Port-Durchleitung) |

## Konsequenzen

- Positiv: Safe-Watermark-Logik ist domänentestbar; die physische Ausführung
  bleibt austauschbar.
- Negativ: jede Bereinigung läuft über Port und Adapter (Indirektion).
- Folgepflicht: RunRetentionUseCase unter den Inbound Use Cases; blockierende
  Consumer sind sichtbar zu melden ([`LH-FA-RET-005`](../../../spec/lastenheft.md)).

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Trennung von Löschentscheidung und physischer Ausführung
hängt an keine beobachtbare Bedingung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0014` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).