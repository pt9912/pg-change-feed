# ADR-0049: Replication-Fehlerklassen-Trennung und WAL-Rückstand-Schwellen

**Status:** Accepted

**Datum:** 2026-09-12

**Autor:** pt9912

**Bezug:** [`LH-QA-REL-003`](../../../spec/lastenheft.md), [`LH-FA-CAP-004`](../../../spec/lastenheft.md), [ADR-0023](0023-fehlerklassifikation.md)

**Schärft:** [`SPEC-008`](../../../spec/pflichtenheft.md) (Zeile `replication`), [`SPEC-013`](../../../spec/pflichtenheft.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`SPEC-008` fordert für die Fehlerklasse `replication`: „Überwachung über
Schwellen (§5, WAL-Rückstand); kontrollierte Fortsetzung." Der tatsächliche
Code (`internal/bootstrap/wiring.go`, `classifyRunError`/`Run`) kennt bisher
keine Schwelle und keinen Fortsetzungspfad — jeder als `replication`
klassifizierte Fehler beendet den Prozess sofort
(`docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/`).

Beim Nachvollziehen des Code-Bestands zeigt sich, dass `replication`
technisch zwei strukturell verschiedene Fehlerquellen bündelt:

- `mapper.ErrChangeWithoutBegin`, `mapper.ErrCommitWithoutBegin`,
  `mapper.ErrBeginWithoutCommit`
  (`internal/adapters/driving/replication/mapper/mapper.go`): Der Assembler
  hat bereits vollständig empfangene, unbeschädigte `pgoutput`-Ereignisse
  decodiert und stellt dabei fest, dass ihre Reihenfolge dem Vertrag
  widerspricht (Change ohne offenes BEGIN, Commit ohne offenes BEGIN, BEGIN
  über einem bereits offenen BEGIN). Das ist keine Störung der Übertragung,
  sondern ein Widerspruch zur eindeutigen logischen Reihenfolge
  (`LH-FA-CAP-004`) — der Zustand des Assemblers ist ab diesem Punkt nicht
  mehr vertrauenswürdig interpretierbar, unabhängig davon, wie es um den
  WAL-Rückstand steht.
- `receive.ErrReplication`
  (`internal/adapters/driving/replication/receive/receive.go`) und
  `outbound.ErrReplication`
  (`internal/adapters/driven/postgresack/ack.go`,
  `internal/application/port/outbound/replicationack.go`): Diese Sentinels
  tragen Störungen der Übertragung selbst — Verbindungsaufbau
  (`connectReplication`), Start des Streams (`START_REPLICATION`),
  Katalogabfragen und Slot-Verwaltung (`ensureSlot`/`ensurePublication`),
  Keepalive-Empfang/-Antwort, Nachrichtenempfang (`ReceiveMessage`) und die
  Quell-Bestätigung (`Acknowledge`). Sie zeigen an, dass die Verbindung zur
  Quelle momentan nicht funktioniert — nicht, dass bereits empfangene Daten
  widersprüchlich sind. Genau eine unterbrochene oder ausbleibende
  Verbindung ist der reale Mechanismus, über den ein Replication-Slot
  inaktiv wird und sein WAL-Rückstand unbegrenzt wächst
  (`cdc_wal_retention_bytes`, `SPEC-009`) — PostgreSQL erzwingt ohne
  `max_slot_wal_keep_size` keine Deckelung.

[`SPEC-013`](../../../spec/pflichtenheft.md) legt bisher nur für
`cdc_capture_lag` einen Warn-/Fehlerwert fest (`CDC_LAG_THRESHOLDS`); für den
WAL-Rückstand, den [`SPEC-008`](../../../spec/pflichtenheft.md)s Prosa
ausdrücklich nach §5 verortet, fehlt der Initialwert — eine Spec-interne
Inkonsistenz: [`SPEC-013`](../../../spec/pflichtenheft.md) trägt bereits den
Anspruch, Initialwert für WAL-Rückstand *und* Capture-Lag zu sein (§5, Absatz
nach der Metriken-Tabelle), löst ihn für WAL-Rückstand aber noch nicht ein.

## Entscheidung

Wir treffen zwei Festlegungen:

**(a) Sentinel-Trennung innerhalb der Klasse `replication`:**
`mapper.ErrChangeWithoutBegin`, `mapper.ErrCommitWithoutBegin` und
`mapper.ErrBeginWithoutCommit` gelten als **Stream-Ordnungs-Verletzung** —
sie bleiben unverändert hart abbrechend, unabhängig vom WAL-Rückstand.
`receive.ErrReplication` und `outbound.ErrReplication` gelten als
**Transport-/Verbindungsstörung** — sie sind der Kandidatenkreis, den die
Schwellen-Überwachung aus `SPEC-008` (kontrollierte Fortsetzung unterhalb der
Warnschwelle, kontrollierter Abbruch oberhalb der Fehlerschwelle) betrifft;
die Umsetzung im Capture-Pfad liefert `slice-026`. <!-- d-check:status-provenance -->

**(b) WAL-Rückstand-Schwellen (`cdc_wal_retention_bytes`, `SPEC-009`):**
Warnschwelle **100 MiB**, Fehlerschwelle **1 GiB** — Startwerte, die
`SPEC-013` als `CDC_THRESHOLDS`-Eintrag neben `cdc_capture_lag` trägt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### (a) Sentinel-Trennung

| Option | Pro | Contra |
|---|---|---|
| A — `replication` bleibt undifferenziert; alle sieben Sentinels weiter hart abbrechend (nichts tun) | keine Übersetzungsarbeit, kein Analyse-Aufwand | `SPEC-008`s Zusage „kontrollierte Fortsetzung" bleibt für **jeden** `replication`-Fehler unerfüllbar, auch für reine Verbindungswackler; die Beobachtung `BEO-PGC/spec008-replication-luecke` bleibt offen |
| B — alle fünf Sentinels werden zu Schwellen-Kandidaten (auch die Mapper-Sentinels) | maximale Vereinheitlichung, ein einziger Prüfpfad in `slice-026` <!-- d-check:status-provenance --> | Eine Stream-Ordnungs-Verletzung bedeutet, dass der Assembler-Zustand bereits inkonsistent ist (offene Transaktion mit undefiniertem Inhalt) — eine Fortsetzung liest oder schreibt aus einem Zustand, dessen Korrektheit nicht mehr feststeht, unabhängig vom WAL-Rückstand; das widerspricht `LH-FA-CAP-004` |
| **C — zwei Klassen nach Fehlerquelle: Stream-Ordnungs-Verletzung (`mapper.*`) hart, Transport-/Verbindungsstörung (`receive.ErrReplication`/`outbound.ErrReplication`) Schwellen-Kandidat (gewählt)** | folgt der bestehenden Adaptergrenze (Mapper übersetzt bereits decodierte Ereignisse, Stream/ACK sprechen mit der Verbindung); `LH-FA-CAP-004` bleibt für den Assembler-Zustand unangetastet; genau die Fehlerquelle, die WAL-Rückstand real verursacht (unterbrochene Verbindung), wird schwellen-überwacht | zwei Sentinel-Gruppen statt einer, Übersetzungsentscheidung muss bei jedem neuen Sentinel wiederholt werden |

### (b) WAL-Rückstand-Schwellenwerte

| Option | Pro | Contra |
|---|---|---|
| A — kein Initialwert, jede Instanz konfiguriert ohne Default (nichts tun) | keine vorab geratene Zahl | Betrieb hat am ersten Tag keine Überwachung; widerspricht dem in `SPEC-013` bereits etablierten Stil (Initialwert + „über ADR schärfbar", wie bei `cdc_capture_lag`) |
| B — eng: Warnschwelle 10 MiB / Fehlerschwelle 100 MiB | früheste mögliche Warnung | normale kurzzeitige Lastspitzen oder ein Wartungsfenster lösen ständig Fehlalarme aus (Alert-Fatigue), ohne dass eine Disk-Erschöpfung droht |
| **C — Warnschwelle 100 MiB / Fehlerschwelle 1 GiB, Faktor 10 (gewählt)** | rund und im Stil von `SPEC-013`s bestehenden Zeitschwellen (`cdc_capture_lag`: Faktor 12 zwischen Warn und Fehler); 1 GiB ist klein genug, um vor einer Disk-Erschöpfung zu reagieren, aber groß genug, um kurzzeitige Lastspitzen nicht fälschlich zu melden | nicht aus Produktionsmessungen abgeleitet, reine Stil-Analogie — im Slice-Plan (`slice-024` §6) als Risiko vermerkt <!-- d-check:status-provenance --> |

## Konsequenzen

- Positiv: `SPEC-008`s Zusage „Überwachung über Schwellen; kontrollierte
  Fortsetzung" bekommt einen konkreten, umsetzbaren Kandidatenkreis
  (Transport-/Verbindungsstörung) statt der gesamten `replication`-Klasse;
  Stream-Ordnungs-Verletzungen bleiben unmissverständlich hart abbrechend.
  `SPEC-013` trägt ab sofort einen Initialwert für WAL-Rückstand, wie es §5
  bereits unterstellte.
- Negativ: Die Zwei-Klassen-Trennung ist an den heutigen fünf Sentinels
  festgemacht; ein neuer `replication`-Sentinel in einem künftigen Adapter
  braucht dieselbe Einordnungsentscheidung erneut (siehe Re-Evaluierungs-
  Trigger).
- Folgepflicht: `slice-025` liefert die Metrik `cdc_wal_retention_bytes`; <!-- d-check:status-provenance -->
  `slice-026` liefert die Schwellen-Überwachung im Capture-Pfad und den in <!-- d-check:status-provenance -->
  `welle-7` §3 geforderten Ende-zu-Ende-Beleg (WAL-Rückstand unterhalb der
  Warnschwelle → unveränderte Fortsetzung; oberhalb der Fehlerschwelle →
  kontrollierter Abbruch mit `replication`-Klassifikation).

## Fitness Function (falls maschinell prüfbar)

Noch nicht anwendbar auf dieser Entscheidungs-Ebene — diese ADR liefert
keine Implementierung (`slice-024` §1 Out-of-Scope). Der maschinell prüfbare <!-- d-check:status-provenance -->
Beleg entsteht mit `slice-026`: der Ende-zu-Ende-Test aus `welle-7` §3, der <!-- d-check:status-provenance -->
beide Seiten der Schwelle real durchläuft.

| Tooling | Regel | Make-Target |
|---|---|---|
| — | wird mit `slice-026` benannt <!-- d-check:status-provenance --> | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Für (a), die Sentinel-Trennung: **permanent** — die Grenze folgt der
Adapterstruktur (Mapper vs. Stream/ACK) und ist unabhängig von konkreten
Zahlenwerten.

Für (b), die Byte-Schwellen: **nicht permanent.** Sobald der Ende-zu-Ende-
Test aus `welle-7` §3 oder späterer Produktionsbetrieb zeigt, dass
Warnschwelle oder Fehlerschwelle zu grob oder zu fein greifen (z. B. der Test
schlägt mit den hier gesetzten Startwerten fehl, oder reale WAL-Wachstumsraten
weichen deutlich ab) — Folge-ADR mit `Supersedes ADR-0049` für die
Zahlenwerte in (b); Entscheidung (a) bleibt davon unberührt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-12 | Accepted | `docs/plan/planning/in-progress/slice-024-adr-fehlerklassen-schwellen-praezisierung.md` <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0049` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
