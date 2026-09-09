# ADR-0040: ClockPort

**Status:** Accepted

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-FA-RET-003`](../../../spec/lastenheft.md) (zeitbasierte
Retention), [`LH-FA-RET-005`](../../../spec/lastenheft.md) (Alter der
blockierenden Retention), [`LH-QA-OPS-003`](../../../spec/lastenheft.md)
(Alter des ältesten Changes), [`LH-FA-ADM-004`](../../../spec/lastenheft.md)
(CDC-Lag)

**Schärft:** [`ARC-004`](../../../spec/architecture.md) (Outbound Ports als
Fähigkeits-Schnittstellen) und [architecture.md §3](../../../spec/architecture.md)
(SYSTEM `ARC-012` — Zeit)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Drei Zusagen des Lastenhefts hängen an der Systemzeit: die zeitbasierte
Retention ([`LH-FA-RET-003`](../../../spec/lastenheft.md)) und die
Sichtbarkeit alternder Changes bzw. ihres blockierenden Einflusses
([`LH-FA-RET-006`](../../../spec/lastenheft.md)), die Metrik
`cdc_oldest_change_age`
([`SPEC-009`](../../../spec/pflichtenheft.md)) und der messbare CDC-Abstand
([`LH-FA-ADM-004`](../../../spec/lastenheft.md)). Importiert Domain oder
Application die globale Wanduhr direkt (`time.Now()`), sind Retention-Tests
nicht deterministisch, die Alter-Metriken schwanken im Testlauf, und die
Schichten-Disziplin (Domain ohne Treiber) bekommt eine versteckte Ausnahme.
Der Port-Bestand (ChangeStorePort, ReplicationAckPort, ConsumerStatePort,
SchemaStorePort, TransactionBufferPort, MetricsPort) enthält bisher keine
Zeit-Fähigkeit.

## Entscheidung

Zeit ist eine Fähigkeit hinter einem Outbound Port: der Application Layer
empfängt die Wanduhr ausschließlich über den **`ClockPort`**
(`Now() domain.TimePoint`); die einzige Produktionsimplementierung ist der
`SystemClockAdapter` (Driven, im Bootstrap verdrahtet). Domain- und
Application-Code rufen `time.Now()` nicht direkt auf; Tests setzen einen
Fake Clock ein.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — `time.Now()` direkt in Domain/Application | kein Boilerplate, kürzeste Signaturen | nicht deterministisch testbar; versteckte Schichten-Ausnahme; Retention-Tests hängen von Realzeit ab |
| B — Zeit als Parameter an jeden Use Case (`now`-Argument, Aufrufer trägt) | explizit, kein Port | Signaturen-Noise über alle neun Use Cases; die Zeit müsste der Driving Adapter liefern — falsche Schicht |
| **C — ClockPort (gewählt)** | Zeit ist eine testbare Fähigkeit wie Speicher/Ack; deterministische Retention- und Lag-Tests; konsistent mit Ports-nach-Fähigkeiten | ein weiterer Port und eine Adapter-Implementierung; geringes Boilerplate je Use-Case-Konstruktor |
| D — nichts tun (Wanduhr belassen) | — | Testpyramide bekommt flaky Retention-Tests; die Ältest-Change-Metrik wird unprüfbar |

## Konsequenzen

- Positiv: Retention-Bereinigung und Change-Alter sind in
  Application-Tests deterministisch prüfbar; die Fake-Clock passt in die
  Application-Tests mit Fake Ports (Testpyramide);
  [`ARC-012`](../../../spec/architecture.md) bekommt die
  Zeit als benannten Berührungspunkt der Sicht.
- Negativ: ein Port mehr; jeder Retention-/Lag-Codepfad braucht die
  injizierte Uhr statt des Direktaufrufs.
- Folgepflicht: `SystemClockAdapter` als Driven Adapter
  ([`ARC-006`](../../../spec/architecture.md)-Familie); die Metrik
  `cdc_oldest_change_age` ([`SPEC-009`](../../../spec/pflichtenheft.md))
  liest ihr Alter über den Port, nicht über eine globale Uhr.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| depguard | `internal/domain` und `internal/application/usecase/*` importieren das Paket `time` nicht | `— (Gate geplant, siehe [`ADR-0036`](0036-architekturpruefung-ci.md))` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

`permanent` — die Zeit-Abhängigkeit der Retention ist eine Struktur-Zusage
des Lastenhefts, kein Wandel-Event; eine Ablösung wäre eine neue
Entscheidung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted (Anlass: Paketstruktur-Ideen, ADR-0039; Lücke im Port-Bestand) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).