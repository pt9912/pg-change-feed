# Architect-Review welle-5 — Trigger-Audit (Carveout · Bootstrap-aware Gate · ADR-Re-Evaluierung) und Beobachtungs-Register-Verkörperung (BEO-PGC/dod-checkbox-nachzug, 3×)

**Rolle:** Architect (Modul 8). **Datum:** 2026-09-12.
**Eingang:** [`docs/plan/planning/welle-5.md`](../planning/done/welle-5.md) (§1–§3,
Closure-Trigger noch offen zum Zeitpunkt dieses Laufs) ·
[`docs/plan/planning/done/slice-017-commit-zeitstempel-decoder-mapper-domaene.md`](../planning/done/slice-017-commit-zeitstempel-decoder-mapper-domaene.md),
[`docs/plan/planning/done/slice-018-commit-zeitstempel-store-adapter.md`](../planning/done/slice-018-commit-zeitstempel-store-adapter.md),
[`docs/plan/planning/done/slice-019-cdc-capture-lag-ablösen.md`](../planning/done/slice-019-cdc-capture-lag-ablösen.md)
(je §1–§8) ·
[`docs/reviews/review-slice-017.md`](../../reviews/review-slice-017.md),
[`docs/reviews/review-slice-018.md`](../../reviews/review-slice-018.md),
[`docs/reviews/review-slice-019.md`](../../reviews/review-slice-019.md) ·
[`docs/reviews/verify-slice-017.md`](../../reviews/verify-slice-017.md),
[`docs/reviews/verify-slice-018.md`](../../reviews/verify-slice-018.md),
[`docs/reviews/verify-slice-019.md`](../../reviews/verify-slice-019.md) ·
[`ADR-0040`](0040-clockport.md) (Accepted, `permanent`) ·
`docs/plan/carveouts/` (nur `.gitkeep`) ·
`harness/conventions.md` §Modus-Deklaration (`PGC`, Greenfield, gesamtes
Repo) ·
`docs/plan/planning/observations/BEO-PGC/` (alle zwölf Verzeichnisse,
insbesondere `dod-checkbox-nachzug/` — `observation.md`, `state.md`,
`evidence/{slice-015,slice-016,slice-017}.md`) · Baseline-Regelwerk
`modul-06-roadmap.md` §Wellen-Closure-Prozedur (Closure-Schritt 2
Trigger-Audit, Schritt 3a/3b Lese-Schritt/Verkörperung) ·
`modul-08-agentenrollen.md` §Rollen-Sequenz für eine Welle.

**Ausgang:** Zwei unabhängige Züge, beide mit Verdikt:

1. **Trigger-Audit der Welle (Modul 6, Closure-Schritt 2):** Alle drei
   Artefaktklassen geprüft — **Carveout: 0 aktiv.** **Bootstrap-aware Gate:
   0 fällig** (Repo ist durchgehend Greenfield, keine Reifestufe zum
   Hochschalten). **ADR-Re-Evaluierung:** `ADR-0040` durchgehend eingehalten,
   ihr Trigger ist `permanent` und feuert nicht; keine andere in dieser Welle
   berührte ADR (`ADR-0011`, `ADR-0044`) hat eine fällige Re-Evaluierung.
2. **Beobachtungs-Register — Lese-Schritt (Modul 6, Closure-Schritt 3a/3b):**
   `BEO-PGC/dod-checkbox-nachzug` erreicht mit `evidence/slice-017.md` 3×.
   Ausgang: **verkörpert.** Geschärfte Regel geschrieben in
   [`.claude/commands/implement-slice.md`](../../../.claude/commands/implement-slice.md)
   (Pre-completion-Checkliste, Schritt 18, plus Fixrunden-Ergänzung an
   Schritt 21) · Herkunfts-Anker `seit welle-5`.

**Harte Regel eingehalten:** `ADR-0040` (`Accepted`) wird von diesem Lauf
**nicht** inhaltlich geändert — der Trigger-Audit bestätigt sie nur. Die
verkörperte Regel geht in den Implementer-Workflow (`.claude/commands/`),
nicht in eine bestehende Accepted-ADR; das Beobachtungs-Register-`state.md`
wird auf den zugewiesenen Ausgang fortgeschrieben, nicht rückwirkend
umgeschrieben. `welle-5.md` selbst und die drei Slice-Dateien in `done/`
werden von diesem Lauf **nicht** editiert — die restliche Welle-Closure
(Schritte 1, 3c–6 aus Modul 6: Trigger-Prüfung, Closure-Notiz, `git mv`,
Paarungen, Archivierung, Roadmap-Fortschreibung) bleibt Planner-Arbeit.

---

## Zug 1 — Trigger-Audit der Welle

### 1a. Carveout (Modul 7)

`docs/plan/carveouts/` enthält ausschließlich `.gitkeep` — kein aktiver
Carveout im Repo, damit auch keiner mit Bezug zu `welle-5`. Kein Gate lief
in dieser Welle rot: `make gates` wurde in jedem der drei
Verifier-Läufe eigenständig ausgeführt und lieferte durchgehend Exit 0
(`verify-slice-017.md`, `verify-slice-018.md`, `verify-slice-019.md`, je
Sensor-Tabelle). **Feststellung: 0 offen** — weder aufzulösen noch zu
verlängern, weil keiner existiert.

### 1b. Bootstrap-aware Gate (Modul 13)

`harness/conventions.md` §Modus-Deklaration führt genau eine Sub-Area
(`*`/`PGC`, Greenfield) für das gesamte Repo. Alle drei Slice-Pläne
bestätigen in ihrem §8 unabhängig: „Die berührte Sub-Area … ist ein
Segment des GF-Baums der Modus-Deklaration … Reiner GF-Hinweis genügt …
kein Sub-Area-Block." Es gibt keine Brownfield-/Hybrid-Sub-Area und damit
keine Reifestufe (Phase × Modus-Matrix, Modul 2), die diese Welle
hochschalten könnte oder müsste. **Feststellung: 0 offen** — die Frage
stellt sich in einem durchgehend Greenfield-Repo nicht.

### 1c. ADR-Re-Evaluierung (Modul 4), mit Schwerpunkt ADR-0040

**`ADR-0040` (ClockPort, `internal/domain` importiert `time` nicht,
`permanent`).** Geprüft gegen die eigenständige (nicht aus dem
Reviewer-Bericht übernommene) Verifikation in allen drei Slices:

| Slice | Eigener Beleg des Verifiers | Fund |
|---|---|---|
| slice-017 | `grep -rn '"time"' internal/domain/ internal/application/usecase/` über den gesamten Baum | keine Treffer — `mapper.go` (Driving-Adapter) übersetzt `time.Time`→`int64` vor der Domänen-Grenze |
| slice-018 | derselbe `grep`, plus Lesen von `mapper.go` (Driven-Adapter): `time.Unix(0, sourceCommittedAt.UnixNanos).UTC()` | keine Treffer in Domäne/Application; Konvertierung ausschließlich im Store-Adapter-Mapper |
| slice-019 | `git grep -c "cdc_capture_lag_approx"` (Metrik-Umbenennung, kein Zeit-Typ-Bezug) — dieser Slice berührt `ADR-0040` nicht direkt, sondern `ADR-0044` (Image-Beleg) | `ADR-0040` unberührt, kein neuer Treffer möglich (kein Zeit-Code in diesem Diff) |

Die Grenze wurde über alle drei Schichten (Decoder → Mapper → Domäne →
Store-Adapter-Mapper) hinweg durchgehend eingehalten, real geprüft und
nicht nur behauptet — jeweils mit eigenem `grep` in einem frischen
Verifier-Kontext, nicht durch Übernahme des Implementer- oder
Reviewer-Befunds. Der Re-Evaluierungs-Trigger von `ADR-0040` ist
`permanent` (kein beobachtbares Ereignis definiert, das eine
Neubewertung auslöst — die Zeit-Abhängigkeit ist eine Struktur-Zusage des
Lastenhefts, siehe ADR-Text). Es liegt zudem keine Situation vor, aus der
sich ein *neuer*, in der ADR nicht vorgesehener Trigger ableiten ließe:
Die Welle hat die Entscheidung nicht belastet, sondern sie dreimal
unabhängig bestätigt. **Feststellung: kein Trigger eingetreten, `ADR-0040`
bleibt `Accepted`, `permanent`, unverändert. Kein Folge-ADR.**

**Sonstige in dieser Welle berührte ADRs.** `ADR-0011` (Persist-before-ACK)
wurde in `review-slice-018.md` als Kontext für das Idempotenz-Risiko
(§6-Risiko-2) herangezogen — ihr Re-Evaluierungs-Trigger ist `permanent`
(„die Invariante trägt die Vertragszusage `LH-QA-REL-001`"), und die Welle
liefert keinen Beleg, der diese Invariante infrage stellt; im Gegenteil,
die WAL-Determinismus-Begründung (Retry liefert denselben
Quell-Commit-Zeitstempel) **setzt** `ADR-0011` voraus, statt sie zu
berühren. `ADR-0044` (Image-Beleg-Semantik) wurde in slice-019 angewendet,
nicht neu bewertet: Der Image-Rebuild (`e7232b5`) folgt exakt der in der
ADR beschriebenen Pflicht („Build-Kontext-relevanter Zug → `make image`
vor Closure, Digest-Commit nur bei geändertem Digest") — ihr eigener
Re-Evaluierungs-Trigger („Lauf-Zweig braucht Inhalts-Vergleich über
Umgebungen/Läufe") ist in dieser Welle nicht eingetreten. **Feststellung:
keine weitere ADR mit fälliger Re-Evaluierung.**

---

## Zug 2 — Beobachtungs-Register, Lese-Schritt (`BEO-PGC/dod-checkbox-nachzug`)

### Zähler und Belege

Drei `evidence/`-Dateien: `slice-015.md`, `slice-016.md` (beide
rückwirkend erfasst, siehe `observation.md` §„Benannt, nicht gezählt"),
`slice-017.md` (bei der Closure von slice-017 neu angelegt, dritter
Beleg). Eigene Nachzählung bestätigt: `find
docs/plan/planning/observations/BEO-PGC/dod-checkbox-nachzug/evidence/`
liefert genau drei Dateien. `state.md` sagt korrekt „Zähler (abgeleitet):
3× … Schwelle erreicht", Ausgang bislang unzugewiesen, dem Lese-Schritt der
`welle-5`-Closure zugeordnet — exakt der Punkt, an dem dieser Zug greift.

### Belegte Stabilität des Musters (Modul 6 §Beobachtungs-Register, "1× notieren · 2× Symptom · 3× Lücke")

Über die drei gezählten Belege hinaus bestätigen alle drei
Verifier-Reports dieser Welle dieselbe Klasse zusätzlich, ohne den Zähler
erneut zu erhöhen (Modul 6: „Ein Vorgang zählt einmal"):

- `verify-slice-017.md` VF-1: fünf von zehn DoD-Checkboxen unchecked trotz
  materieller Erledigung — explizit als „dieselbe Klasse wie
  `verify-slice-016.md` VF-2 (und davor `verify-slice-015.md` V-1)"
  benannt, dies **ist** der dritte Beleg (`evidence/slice-017.md`).
- `verify-slice-018.md` (Zusammenfassung DoD, VF-1/VF-2 zu anderen
  Themen; DoD-Checkboxen dort explizit als „regulär noch unchecked, da
  Datei in `in-progress/`" gewertet — kein neuer Fund, weil zum
  Prüfzeitpunkt noch keine Closure stattfand) — bestätigt indirekt, dass
  das Muster nur beim `git mv` nach `done/` sichtbar wird.
- `verify-slice-019.md` VF-2: „DoD-Checkbox „Review durchgeführt" bleibt
  unchecked, obwohl Review und Fixrunde bereits abgeschlossen sind" —
  ausdrücklich als „exakt die Musterklasse der bereits registrierten
  Beobachtung `BEO-PGC/dod-checkbox-nachzug` … kein neuer Fund" markiert,
  **kein** zusätzlicher Zähler-Beitrag (korrekt nach Modul 6: derselbe
  Vorgang, slice-019, zählt nicht doppelt neben einem etwaigen eigenen
  Beleg für diesen Slice — und da `slice-019` selbst keine eigene
  `evidence/`-Datei für diese Klasse führt, bleibt der Zähler bei 3×, wie
  in `state.md` vermerkt).

Sieben Vorkommen über fünf verschiedene Slices (015, 016, 017, 018, 019) in
derselben Session, davon drei gezählt (unterschiedliche Vorgänge) und vier
kommentierend bestätigt (keine neuen Vorgänge oder nicht separat
zählbar): Das Muster ist **stabil und reproduzierbar**, nicht ein
Einzelfall. Es trägt die Schwelle klar, nicht nur formal.

### Prüfung: trägt „verkörpert"?

Modul 6 nennt zwei Bedingungen, die an der Schwelle hängen —
*verkörpert* (Zielort + Herkunfts-Anker) und *geplant* (Kennung eines noch
zu schreibenden Slice/Welle) — sowie *gestrichen* unabhängig davon. Geprüft:

- **Ist die Ursache strukturell behebbar, ohne ein neues Gate zu
  erfinden?** Ja. Das Muster entsteht, weil das *Behaupten* der DoD
  (Implementer, Schritt 18 des Workflows) und das *Ankreuzen* der
  Checkbox als zwei getrennte Handlungen behandelt werden, obwohl sie
  dieselbe Information tragen. Eine Workflow-Zeile, die beide Handlungen
  im selben Schritt zusammenführt, behebt die Ursache direkt — dieselbe
  Form wie die bereits verkörperte „Plan-Nachzug im selben Lauf"-Regel
  (`.claude/commands/implement-slice.md`, Schritt 14, seit slice-009):
  auch dort wurde eine wiederkehrende Nachtrags-Klasse durch eine
  Workflow-Pflicht am Entstehungsort geschlossen, nicht durch ein Gate.
- **Ist eine mechanische Absicherung (Sensor/Gate) sinnvoll oder
  möglich?** Geprüft und **verneint**, aus demselben Grund wie in der
  Aufgabenstellung vermutet: Ein Sensor müsste den Satz „diese Checkbox
  ist `[ ]`, aber der referenzierte Liefer-Punkt ist materiell erledigt"
  automatisch entscheiden — das verlangt ein Urteil über den *fachlichen*
  Erledigungsstand (z. B. „ist `TestDecodeFlowToCapture` wirklich der
  Beleg für DoD-Punkt 1?"), nicht nur eine Textform. Das ist dieselbe
  Grenze, die Modul 6 für den Beobachtungs-Register-Beleg selbst zieht
  („die *Existenz* der Datei wird nicht verlangt … das ist die Grenze der
  Deklaration"): Ein Gate kann *Checkbox-Syntax* prüfen (`[ ]` vs. `[x]`
  als Zeichen), aber nicht *Checkbox-Wahrheit* gegen beliebigen
  Freitext-DoD-Inhalt. **Bestätigt**, keine mechanische Absicherung über
  eine reine Workflow-Disziplin hinaus.
- **Widerspricht die Verkörperung einer bereits Accepted-Entscheidung?**
  Nein — es existiert keine ADR zu DoD-Checkbox-Handhabung; der Zielort
  ist ein Workflow-Dokument (`.claude/commands/`), kein ADR-Eingriff.

**Verdikt Zug 2:** Ausgang **verkörpert**. Die Regel wird an zwei Stellen
desselben Workflow-Dokuments ergänzt (Details siehe Commit-Diff):

1. **Schritt 18** (Pre-completion-Checkliste — hier wird die DoD "Punkt für
   Punkt behauptet"): Jede zu diesem Zeitpunkt materiell erfüllte oder
   korrekt entfallene DoD-Zeile wird **im selben Lauf** in §2 von `[ ]` auf
   `[x]` gesetzt, statt nur im Bericht behauptet zu werden. Punkte, die die
   Rollen-Sequenz zu diesem Zeitpunkt noch nicht durchlaufen haben (Review,
   Verifikation, Register-/Risiko-Ausgänge — Planner-Closure-Arbeit),
   bleiben regulär `[ ]`.
2. **Schritt 21** (Übergabe an den Reviewer, Fixrunde nach Findings): Löst
   eine Fixrunde einen bislang offenen Punkt auf (typischerweise „Review
   durchgeführt … kein offenes HIGH"), wird die zugehörige Checkbox im
   Fixrunden-Commit mitgesetzt — das schließt exakt die von
   `verify-slice-019.md` VF-2 benannte Lücke (Checkbox blieb nach
   abgeschlossener Fixrunde `[ ]`).

Zielort: [`.claude/commands/implement-slice.md`](../../../.claude/commands/implement-slice.md).
Herkunfts-Anker: `seit welle-5`.

---

## Disposition (Gesamt)

| Frage | Ergebnis |
|---|---|
| Aktiver Carveout mit Bezug zu `welle-5`? | Nein — 0 aktiv im gesamten Repo |
| Reifestufe (Bootstrap-aware Gate) hochzuschalten? | Nein — durchgehend Greenfield, keine Stufe vorhanden |
| `ADR-0040`-Re-Evaluierungs-Trigger eingetreten? | Nein — `permanent`, durchgehend eingehalten (eigenständig in allen drei Slices verifiziert) |
| Andere ADR mit fälliger Re-Evaluierung? | Nein — `ADR-0011`/`ADR-0044` beide `permanent` bzw. Trigger nicht eingetreten, in dieser Welle nur angewendet, nicht belastet |
| `BEO-PGC/dod-checkbox-nachzug` (3×) — Ausgang | **verkörpert** → `.claude/commands/implement-slice.md` Schritt 18 + 21, `seit welle-5` |
| Mechanische Absicherung für die neue Regel? | Geprüft, verneint — Checkbox-Wahrheit gegen Freitext-DoD ist kein Gate-Gegenstand; reine Workflow-Disziplin, wie bei „Plan-Nachzug" (seit slice-009) |
| Auswirkung auf die übrige `welle-5`-Closure | Keine Blocker — Schritte 1, 3c–6 (Trigger-Prüfung, Closure-Notiz, `git mv`, Paarungen, Archivierung, Roadmap) bleiben Planner-Arbeit |

---

## Beleg-Anker (Kurzfassung)

| Aussage | Beleg |
|---|---|
| 0 aktive Carveouts | `docs/plan/carveouts/` (nur `.gitkeep`) |
| Durchgehend Greenfield, keine Reifestufe | `harness/conventions.md` §Modus-Deklaration; slice-017/018/019 §8 je „Reiner GF-Hinweis genügt" |
| `ADR-0040` real eingehalten, alle drei Slices | `verify-slice-017.md` „ADR-0040-Konformität (eigenständig verifiziert)"; `verify-slice-018.md` dieselbe Sektion; `review-slice-017.md`/`review-slice-018.md` Prüfung 1 |
| `ADR-0040`-Trigger `permanent` | [`ADR-0040`](0040-clockport.md) §Re-Evaluierungs-Trigger |
| `ADR-0011`/`ADR-0044` nur angewendet, nicht belastet | `review-slice-018.md` Prüfung 4 (Idempotenz/`ADR-0011`); `verify-slice-019.md` Sensor-Tabelle + Prüfung 4 (`ADR-0044`-Digest-Commit) |
| Zähler `dod-checkbox-nachzug` = 3× | `BEO-PGC/dod-checkbox-nachzug/evidence/{slice-015,slice-016,slice-017}.md`, `state.md` |
| Wiederholtes Auftreten über die drei Belege hinaus (kein neuer Zähler-Beitrag) | `verify-slice-017.md` VF-1, `verify-slice-019.md` VF-2 |
| Keine mechanische Prüfbarkeit von Checkbox-Wahrheit | Modul 6 §Das Beobachtungs-Register, Abschnitt „Grenze: Die *Existenz* der Datei wird nicht verlangt …" (analoge Argumentationsform) |
| Präzedenzfall „Workflow-Disziplin statt Gate" | `BEO-PGC/plan-nachzug/state.md` (verkörpert seit slice-009, `.claude/commands/implement-slice.md` Schritt 14) |
