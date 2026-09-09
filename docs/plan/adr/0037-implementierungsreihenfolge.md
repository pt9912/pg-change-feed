# ADR-0037: Implementierungsreihenfolge

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-MVP-001`](../../../spec/lastenheft.md)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Invarianten (Persist-before-ACK, Retention-Sicherheit) hängen am
Zusammenspiel aller Schichten. Die Reihenfolge der Implementierung
entscheidet, ob die Invarianten früh im Kern entstehen und Adapter später
anklappen — oder ob Infrastruktur zuerst gebaut wird und die Invarianten
hinterher hineinrefaktoriert werden müssen.

## Entscheidung

Die Implementierung folgt dieser Reihenfolge:

1. Domain und Invarianten
2. Inbound/Outbound Ports
3. Application Services
4. Fake Driving/Driven Adapters
5. PostgreSQL Change Store
6. Replication Stream Driving Adapter
7. Replication ACK Driven Adapter
8. E2E
9. Performance/Recovery/Retention

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Adapter zuerst (von der Infrastruktur her) | früh echte Fortschritte sichtbar | Invarianten entstehen in Adaptern statt im Core; späterer Umbau ist teuer |
| B — alle Layer gleichzeitig von mehreren Seiten | hohe Parallelität | Schnittstellen ändern sich unter laufender Arbeit; Merge-Konflikte gegen unverständigte Grenzen |
| **C — Core → Ports → Services → Fakes → echte Adapter → E2E** | Invarianten zuerst; Fakes erlauben volle Application-Tests vor der ersten realen Instanz | sichtbarer Fortschritt außerhalb von Tests kommt später |

## Konsequenzen

- Positiv: Der MVP-Nachweis (robuste persistente CDC-Abstraktion) wächst
  aus getesteten Schichten, nicht aus einem späten Umbau.
- Negativ: Bis Schritt 5 entsteht kein laufendes End-to-End-Verhalten.
- Folgepflicht: Welle- und Slice-Schnitte folgen dieser Reihenfolge; Abwei-
  chungen sind Umplanung (Drift-Eintrag).

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Neu prüfen, wenn sich der MVP-Schnitt des Lastenhefts ändert
(beobachtbar an der Abnahme-Mapping-Tabelle in Lastenheft §1); sonst
permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (überführt aus dem Architektur-Entwurf `architecture-decision-records.md`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).