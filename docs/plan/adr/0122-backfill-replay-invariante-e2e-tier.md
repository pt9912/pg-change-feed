# ADR-0122: Backfill — Tier der Replay-Invariante (Supersedes ADR-0111, teilweise)

**Status:** Accepted — Supersedes [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
in genau einer Stelle: die Tier-Angabe der dritten Zeile der Fitness Function
(Replay-Invariante), also die Spalten „Tooling“ und „Make-Target“ dieser Zeile.
Alles Übrige von `ADR-0111` bleibt in Kraft, insbesondere die **Regel** der Zeile
(das Log ab Log-Anfang, angewandt als Upsert/Delete, ergibt den Quellstand), die
Herleitung der Replay-Invariante in Teilfrage 3 und die übrigen Zeilen der
Fitness Function.

**Datum:** 2026-09-25

**Autor:** pt9912 (Architect-Rolle, Modul 8; Lese-Schritt der Closure von
`welle-backfill-bestand`, Vollmacht des Auftraggebers); jede Tatsachenaussage
trägt ihren Beleg-Anker (siehe §Kontext).

**Bezug:** [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Boundary: Überlappung
ohne stille Lücke),
[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) (teilweise superseded —
Haupt-Bezug), [`ADR-0121`](0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
(dieselbe Klasse: eine Fitness-Function-Zeile, deren Tier den Beleg nicht trägt)

**Schärft:** —

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) ist `Accepted` und
unberührbar (`AGENTS.md` §3.5). Seine Fitness Function nennt für die
Replay-Invariante einen Tier, den kein Test trägt; die Berichtigung ändert den
Referenten und ist keine Zitat-Korrektur.

### Gemessen

**Zwei Orte in `ADR-0111`.** `git grep -n 'Replay-Invariante'
docs/plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md` (Stand `c82d3333`)
nennt die Herleitung in Teilfrage 3 (Zeile 256, „im Slice als Eigenschaftstest zu
belegen“), den Umfang des Slice `S5` in §Folgepflichten 2 (Zeile 507: „Boundary
(nebenläufige Schreiber während des Runs, Replay-Invariante)“, der Slice trägt
`make test-integration`) und die dritte Zeile der Fitness Function (Zeile 537:
Tooling „Go-Test, reale PostgreSQL (Eigenschaftstest)“, Make-Target
`make test-replication`). Die Folgepflicht und die Fitness-Function-Zeile nennen
verschiedene Tiers.

**Wo der Beleg liegt.** `git grep -n -i -E 'replayinvariant|Replay-Invariante'
c82d3333 -- '*.go'` liefert vier Treffer, alle in
`test/integration/backfill_e2e_test.go` (Zeilen 258, 270, 408, 432): der Test
`TestE2EBackfillReplayInvariant`, ein Eigenschaftstest gegen den laufenden
Feed-Container mit drei nebenläufigen `INSERT`/`UPDATE`/`DELETE`-Schreibern vor
und nach der Snapshot-Position. In `make test-replication` liegt kein Test der
Invariante (derselbe Suchlauf: kein Treffer in `internal/**`). Der Lauf `36065957210`
von `e2e.yml` (Kopf `43137ebf`) druckte in beiden Legs `--- PASS:
TestE2EBackfillReplayInvariant`, „Replay-Invariante: Snapshot-Position 28721360,
54 Backfill-Changes, WAL-Changes davor 409 und dahinter 164, 57 Zeilen im
Quellstand“ (PostgreSQL 17) und „… 31875224, 54 … 409 … 169, 57 …“ (PostgreSQL 18)
— **übernommen** aus dem Plan von `slice-backfill-e2e` (§6, dort mit `gh run view
36065957210 --log` gemessen).

**Begründung der Ortswahl.** Die Invariante verlangt Snapshot-Position,
nebenläufige Schreiber und den WAL-Pfad zugleich; das leistet erst der komponierte
Lauf am Feed-Container (Plan von `slice-backfill-e2e`, §3, Zeile „Ortswahl der
Replay-Invariante“, **übernommen**). Ob ein Tier-Test derselben Eigenschaft
herstellbar wäre, ist nicht gemessen.

### Konstraints

- Eine `Accepted`-ADR ist Beleg für nachgeordnete Träger; ihre
  Fitness-Function-Zeilen nennen den Tier, der sie trägt (`AGENTS.md` §3.12
  Instanz B).
- `make test-integration` läuft in beiden PostgreSQL-Versionen von `SPEC-012` (die
  Matrix von `e2e.yml`); der Tier belegt die Invariante damit an 17 und 18.

## Entscheidung

Wir wählen **eine neue ADR mit teilweisem `Supersedes`**, die den Tier der Zeile auf
den Ort setzt, an dem der Beleg liegt. Eine Festlegung.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, der Plan von `slice-backfill-e2e` trägt die Ortswahl | keine Änderung an der ADR | die Fitness-Function-Zeile nennt einen Tier ohne Test; wer `ADR-0111` liest, sucht die Invariante in `make test-replication` |
| B — die Zeile in `ADR-0111` in-place ändern | ein Ort | ändert die Fitness Function, keine Zitat-Korrektur (`AGENTS.md` §3.5, [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)) |
| C — einen Tier-Test der Invariante zusätzlich bauen | die Zeile gilt wörtlich | Snapshot, Schreiber und WAL-Pfad zugleich sind der Aufbau des E2E-Belegs; ein Tier-Test wiederholt ihn in einem zweiten Tier ohne neue Aussage |
| **D — neue ADR mit teilweisem `Supersedes` (gewählt)** | die Zeile trägt den Tier, an dem der Beleg liegt; `ADR-0111` bleibt unberührt | eine weitere ADR für eine Tabellenzeile |

### Festlegung 1 — Tier der Replay-Invariante

Die Replay-Invariante belegt der E2E-Tier: `TestE2EBackfillReplayInvariant` in
`make test-integration`, an PostgreSQL 17 und 18. Diese Festlegung ersetzt in der
dritten Zeile der Fitness Function von `ADR-0111` das Tooling „Go-Test, reale
PostgreSQL (Eigenschaftstest)“ durch „E2E, Eigenschaftstest gegen den laufenden
Feed-Container“ und das Make-Target `make test-replication` durch
`make test-integration`. Die Regel der Zeile bleibt.

## Konsequenzen

- Positiv: die Fitness Function trägt nur Zeilen, deren Tier den Beleg enthält;
  `ADR-0111` §Folgepflichten 2 und die Fitness Function nennen denselben Tier.
- Negativ: `ADR-0111` trägt die berichtigte Tier-Angabe weiter im Text; wer sie
  liest, liest diese ADR über den Index (Titel „Supers. ADR-0111, teilw.“).
- Folgepflicht: keine — der Plan und die Closure von `slice-backfill-e2e` tragen die
  Ortswahl, der Test liegt im Baum.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| E2E, Eigenschaftstest gegen den laufenden Feed-Container | **Replay-Invariante:** Backfill mit nebenläufigen `INSERT`/`UPDATE`/`DELETE` an der Tabelle; das Log ab Log-Anfang, angewandt als Upsert/Delete, ergibt den Quellstand (`TestE2EBackfillReplayInvariant`) | `make test-integration` |

Die übrigen Zeilen der Fitness Function von `ADR-0111` bleiben.

## Re-Evaluierungs-Trigger

- **Ein Test der Replay-Invariante entsteht im Tier `make test-replication`:**
  Festlegung 1 um diesen Tier ergänzen.
- **Der Aufbau des E2E-Belegs ändert sich** (Snapshot-Position, Schreiber oder
  WAL-Pfad nicht mehr zugleich am Feed-Container): den Tier neu bewerten.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-25 | Accepted — Berichtigung der Tier-Angabe einer Fitness-Function-Zeile von `ADR-0111` an dem Ort, an dem der Beleg liegt | [`LH-FA-CAP-009`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0122` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
