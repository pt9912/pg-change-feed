# ADR-0023: Fehlerklassifikation

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-REL-003`](../../../spec/lastenheft.md)

**Schärft:** [`SPEC-008`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Adapter melden technisch sehr unterschiedliche Fehler
(Bibliotheks-Exceptions, Protokollfehler, Datenbankfehler). Application
und Betrieb brauchen eine stabile Klassifikation, um Wiederholung, Abort
und Sichtbarkeit strategisch zu entscheiden ([`LH-QA-REL-003`](../../../spec/lastenheft.md)), ohne
Bibliotheksdetails zu kennen.

## Entscheidung

Adapter übersetzen technische Fehler in stabile Kategorien:
`transient`, `configuration`, `permission`, `schema`, `storage`,
`replication`, `internal`.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun“ ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — rohe technische Fehler weiterreichen | keine Übersetzungsarbeit | Application müsste Bibliotheksdetails kennen; verletzt die Abhängigkeitsregel der Architektur-Sicht (§2) |
| B — nur generisch/nicht-generisch | minimaler Aufwand | keine differenzierte Strategie (Retry vs. Abort vs. Sichtbarkeit) möglich |
| **C — sieben stabile Fehlerkategorien** | strategische Entscheidungen je Klasse; Metriken je Klasse (`cdc_errors_total{class}`) | Übersetzungspflicht in jedem Adapter |

## Konsequenzen

- Positiv: Application entscheidet je Klasse (z. B. `storage` → kein
  Source-ACK, `transient` → Backoff); Betrieb sieht Fehler je Klasse.
- Negativ: jeder Adapter trägt eine Übersetzungsverantwortung.
- Folgepflicht: —

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

permanent — die Kategorien sind die Vertragsschnittstelle zwischen
Adaptern und Application ([`SPEC-008`](../../../spec/pflichtenheft.md)).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0023` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).