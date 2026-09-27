Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.12 Instanz A
(*„die gedeckte Zahl hängt am Lauf … sie nennt ihn und nie ‚der Ist-Stand‘“*) und
`harness/sensors/coverage-gate.md` §Zählbasis + §Grenze 4 (das gemessene **Band**
dieses Gegenstands, sein verbleibender Block, die abgeleitete Schranke 1 Statement)
· seit slice-089 / seit slice-093. **Nicht gestrichen** — die Beobachtung kann noch
auttreten. **Verworfen** wurde der Kandidat „die Slice-Pläne als Risiko“: ein
Zeitdokument trägt keine Regel (`ADR-0083` §Verglichene Alternativen).
**Benannte Resthälfte:** der Erstauftreten-Fall `slice-057` hat keine Regel und
keinen Fix — er gehört eigenständig geführt, nicht in diesen Eintrag.

Die drei Belege treffen **denselben Gegenstand über verschiedene Träger**:
`slice-057` den `make test-integration`-Exit, `slice-090` die **Coverage-Zahl**
desselben Test-Objekts (`WAL-Retention`, `runWALRetentionCheck`) — dort brach ein
Lauf real ab, hier liefert `go tool cover -func` in 7 von 8 Läufen 87,5 % und in
einem 100,0 %, bei durchweg grünen Tests. Beide Male wechselt das Ergebnis bei
**unverändertem Stand**.

**Vierter Beleg (`slice-capture-leerlauf-quellbelege`) — die Resthälfte hat ein
zweites Auftreten.** Der erste `e2e.yml`-Lauf nach dem Push (Lauf 36287009221, Leg
PostgreSQL 18, erster Versuch) endete rot in der Phase „Leerlauf-Bestätigung“ von
`make test-integration`: der Feed-Container endete mit dem WAL-Schwellen-Fehler
(Rückstand 15.238.216 B über der Fehlerschwelle 8.388.608 B), im Wiederholungsversuch
lief dieselbe Phase grün. Gemessen (`evidence/architect-verdict-leerlauf-bestaetigung-intermittenz.md`):
1 rotes erstes Ergebnis unter **80** Ausführungen der Phase (40 abgeschlossene Läufe von
`e2e.yml` auf Commits mit der Phase, je zwei Legs; 79 grün).
Zähler: **4×**. Die Zählung teilt sich nach dem, was sie trägt: die **Zahl-Hälfte**
(Coverage, `slice-090`/`slice-091`) bleibt **verkörpert**; die **Resthälfte** — ein
`make test-integration`-Lauf endet rot bei unverändertem Stand und die Wiederholung
ist grün — steht bei **2×** (`slice-057`, `slice-capture-leerlauf-quellbelege`) mit
**verschiedenen Ursachen**: `slice-057` (Retention-Lebenszyklus-Timing) bleibt **1×
ohne Ursache und ohne Ausgang**; der Leerlauf-Fall hat Ursache und Ausgang (unten).

**Ausgang des Leerlauf-Falls — Testaufbau, Slice `slice-leerlauf-phase-last-in-stuecken`
(Plan: `slice-leerlauf-phase-last-in-stuecken`; beauftragt im Architect-Verdikt `docs/reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md`).** <!-- d-check:status-provenance -->
Die Anweisung des Runners — ein einzelnes `INSERT` mit
16,0 MB WAL, dem 1,9-fachen der Fehlerschwelle der Phase — ist ein Stoß: der Rückstand steht
danach kurz bei der ganzen Last, bevor die Bestätigung folgt (am echten Stream lokal 42 bis
63 ms über 8 MiB), und ein Prüf-Takt (5 s) in diesem Fenster beendet den Container. Kein Mangel
der Leerlauf-Bestätigung (`ADR-0120`); die Phase behauptete breiter, als der Mechanismus trägt.
Abhilfe: die Last in Stücken unter der Warnschwelle mit Wartebedingung auf `confirmed_flush_lsn`
zwischen den Stücken, Summe über der Fehlerschwelle (lokal erprobt: höchster Rückstand 32 % der
Fehlerschwelle bei gleicher Gesamtlast). Kette: `slice-capture-leerlauf-quellbelege` (done) →
`slice-leerlauf-phase-last-in-stuecken` → `slice-wal-fehlerschwelle-ausgangsklasse` →
`slice-start-vorlauf-grenze` → `slice-transformationen-e2e-abhilfe`: ein Rot in dieser Phase
lässt jede Phase dahinter und die Schritte des Legs dahinter ungelaufen (Lauf 36287009221,
Versuch 1, Leg PostgreSQL 18: Coverage- und „Replication-Tier“-Schritte `skipped`) — ihre
Belege wären nicht widerlegt, sondern nicht gelaufen.

**Trigger (Neubewertung durch den Architect):** ein weiteres Rot mit der Signatur „Fehlerklasse
`replication` … WAL-Rückstand … über Fehlerschwelle“ in „Leerlauf-Bestätigung“ **nach** dem Slice
— dann ist die Ursache falsch bestimmt, und die Karenz der Fehlerschwelle im Produkt wird neu
entschieden. Erkennung ohne Sensor (beide am 2026-09-27 gefahren; Ausgabe des ersten: genau
Lauf 36287009221, des zweiten für diesen Lauf: `1`):

```text
gh api --paginate 'repos/pt9912/pg-change-feed/actions/workflows/e2e.yml/runs?per_page=100' --jq '.workflow_runs[] | select(.run_attempt > 1) | [.id, .head_sha[0:8], .run_attempt, .created_at] | @tsv'
gh run view <Lauf> --attempt 1 --log-failed | grep -c 'Leerlauf-Bestätigung — nach dem Schreiber'
```

Der Eintrag nimmt darüber hinaus keinen Ausgang vorweg: die Zahl-Hälfte und `slice-057`
bleiben wie oben.
