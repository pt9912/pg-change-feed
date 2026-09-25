# Slice retention-lauf-speicher-begrenzung: Der Retention-Lauf liest je Takt nicht mehr alle Changes einer Quelle in den Speicher des Feed-Containers

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er hat keine Kante zu einer offenen Welle.

**Bezug:** [`LH-FA-RET-002`](../../../../spec/lastenheft.md),
[`LH-FA-RET-003`](../../../../spec/lastenheft.md) und
[`LH-FA-RET-004`](../../../../spec/lastenheft.md) (Retention: Aufbewahrung,
Mindestalter, Consumer-Positionen),
[`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill des Bestands, der
`cdc.change` in einem Zug füllt),
[`ADR-0014`](../../adr/0014-retention-domain-policy.md) (Retention als Domain
Policy),
[`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
(Richtgröße; sein Trigger „Eine Messung liegt vor“), Messbericht
[`messbericht-slice-backfill-speicher-untersuchung`](../../../reviews/messbericht-slice-backfill-speicher-untersuchung.md).

**Berührte Spec-Stellen:** [`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md)
(Safe Watermark der Retention) — gelesen, nicht geändert: die Bereinigungsmenge
bleibt dieselbe, geändert wird die Art, sie zu bestimmen;
[`SPEC-005`](../../../../spec/pflichtenheft.md) (Zeile 72: große Transaktionen
nicht unbegrenzt im RAM) — gelesen.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Implementer-Agent, Ausgang von `slice-backfill-speicher-untersuchung`.
**Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Speicher des Feed-Containers in Ruhe und nach einem Backfill hängt
nicht mehr an der Zahl der Changes in `cdc.change`: der periodische Retention-Lauf
bestimmt die Bereinigungsmenge, ohne alle Changes der Quelle samt Row Images in
den Speicher zu lesen.

**Befund (Herkunft: Messbericht, gemessen und aus dem Code gelesen).**
`RunRetentionService.Run` (`internal/application/usecase/retention/service.go`)
ruft je Takt (`retentionInterval`, 10 s, `internal/bootstrap/wiring.go`)
`ChangeStorePort.ReadChanges` mit der Quelle als einzigem Filter auf; die Abfrage
`SelectChanges` (`internal/adapters/driven/postgresstorage/queries/queries.go`) trägt
`LIMIT NULL`, liefert je Change beide Row Images und wird vollständig in eine Liste
gelesen, bevor `RetentionPolicy.AllowsDeletion` je Change befragt wird. Der
Speicher des Feed-Containers wächst deshalb mit der Zahl der Changes, unabhängig
davon, ob sie ein Backfill oder die laufende Erfassung geschrieben hat. Die
gemessenen Größen stehen im Messbericht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Bytegrenze oder Zeilenbreiten-Wache im Backfill** (`NextBlock`,
  `DefaultBlockSize`). Der Messbericht belegt, dass der Speicher im Run selbst flach
  bleibt; die Blockgrenze in Zeilen ist damit nicht die Ursache.
- **Ein Wert der Richtgröße von 4.000.000 geschätzten Zeilen in
  `internal/application/usecase/backfill/warn.go`.** Ihre Neubewertung ist ein
  Liefer-Punkt dieses Slice, aber erst nach der Messung *nach* der Änderung: vorher
  fehlt die Größe, an der sie zu bemessen wäre.
- **Andere Aufrufer von `ReadChanges` ohne Limit.** Der Slice liest nur den
  Retention-Pfad; ob weitere Aufrufer unbegrenzt lesen, ist nicht untersucht.
- **Ein Gate oder eine Schwelle für den Speicher.** Die Messung bleibt ohne
  Pass/Fail (`harness/targets/bench-backfill.md`); eine Schwelle brauchte eine ADR
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6).

## 2. Definition of Done

- [ ] Der Retention-Lauf liest je Takt eine begrenzte Menge (Seiten oder eine
      Projektion ohne Row Images, Entscheidung des Architects, ob eine ADR nötig
      ist — der Slice berührt den Vertrag von `ChangeStorePort`) und löscht
      dieselbe Menge wie zuvor: die Tests zu `LH-FA-RET-002` bis `LH-FA-RET-004`
      laufen unverändert grün, ein neuer Test bindet die Seitengrenze an ihre
      Eingabe. *Zu belegen durch:* `make test` und ein Test, der bei
      Seitengröße 1 dieselbe Bereinigungsmenge liefert wie bei unbegrenzter Größe.
- [ ] Die Messung nach der Änderung liegt vor: `tools/bench-backfill-memory.sh` mit
      `BENCH_MEM_STAGES=1000000` (Default-Image, gedruckte Zeilen im Bericht des
      Slice); *erwartet:* die Spitze nach dem Run liegt unabhängig von der Zahl der
      Changes in `cdc.change` im Bereich der Spitze im Run des Messberichts. Trifft
      das nicht zu, steht der Befund im Bericht und der Slice geht nicht nach `done/`.
- [ ] Das Benutzerhandbuch trägt unter „Grenzwerte“ die Messung nach der Änderung
      und die Bewertung der Richtgröße nach dem Trigger von
      [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
      („Eine Messung liegt vor“); Version und Änderungshistorie tragen eine Zeile.
- [ ] `make gates` grün — Exit-Code ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

**Umfang:** M — Schätzung, nicht gemessen: eine Änderung an Use Case und Port mit
Adapter-Test gegen PostgreSQL und eine Messung der Stufe mit 1.000.000 Zeilen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/usecase/retention/service.go` | update | die Bereinigungsmenge wird seitenweise bestimmt |
| `internal/application/port/outbound/changestore.go` und `internal/adapters/driven/postgresstorage/store.go` | update | Vertrag und Adapter des Lesens der Kandidaten (Seite, Schlüssel der Ordnung) |
| Tests des Use Cases und des Adapters | update | Seitengrenze an Eingabe gebunden; gleiche Menge bei jeder Seitengröße |
| `docs/user/benutzerhandbuch.md` (§Grenzwerte, Version, Änderungshistorie) | update | Messung nach der Änderung, Bewertung der Richtgröße |

**§3.13-Suchlauf (Träger, die die Aussage „Speicher des Feed-Containers“ und die
Retention-Lesung beschreiben):**

```text
git grep -n -i -E 'alle Changes der Quelle|liest alle|ReadChanges|retentionInterval' -- docs internal harness spec
git grep -n -i -E 'Speicher des Feed-Containers|Spitze im Run' -- docs harness
```

| Träger | Befund | Behandlung |
|---|---|---|
| *(Implementer trägt ein)* | | |

## 4. Trigger

**Start** (`next` → `in-progress`): kein weiterer Slice in `in-progress/`
(WIP-Limit 1). Der Slice muss `done` sein, **bevor** ein Server-Release
veröffentlicht wird, der den Backfill trägt und dessen Handbuch den Speicherbedarf
nicht als bekannte Grenze führt; die Entscheidung, ob der Release mit der
dokumentierten Grenze vorgezogen wird, liegt beim Nutzer (Empfehlung im Messbericht).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls die Änderung des
  Vertrags von `ChangeStorePort` eine ADR verlangt — dann geht der Architect-Zug
  der Umsetzung voraus.
- `in-progress` → `open` (blockiert): falls die Entscheidung des Architects aussteht.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + die Messung nach der Änderung mit gedruckten
Zeilen im Bericht + Handbuch nachgezogen + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Bereinigungsmenge ändert sich** (Sicherheit: kein Change eines aktiven
  Consumers darf gelöscht werden, `LH-FA-RET-004`). *Erwartet, zu belegen durch:* der
  Test „gleiche Menge bei jeder Seitengröße“ und die unverändert grünen Tests der
  Anforderungen. **Ausgang:** *(bei Closure)*
- **Die Abfrage selbst wächst mit der Zahl der Changes** (Sortierung über die ganze
  Menge in PostgreSQL, im Messbericht bei 3.000.000 Zeilen als ausbleibende
  Bereinigungs-Takte beobachtet, Ursache nicht belegt). *Erwartet, zu belegen
  durch:* die Messung nach der Änderung bei 1.000.000 und 3.000.000 Changes.
  **Ausgang:** *(bei Closure)*
- **Ein Backfill während einer laufenden Bereinigung.** Die Schreibtransaktion des
  Backfills ist bis zum Commit unsichtbar; die Seitengrenze darf keine Zeile
  überspringen, die zwischen zwei Seiten sichtbar wird. *Erwartet, zu belegen durch:*
  ein Test mit einem Commit zwischen zwei Seiten. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** *(Anker · Folge-Slice · Register, Ergebnis)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Use Case der Retention und Store-Adapter sind keine eigenen
Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** `BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes`
(Ausgang: `slice-backfill-speicher-untersuchung`, dessen Befund diesen Slice trägt),
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (die Speicher-Zahlen des
Handbuchs).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
