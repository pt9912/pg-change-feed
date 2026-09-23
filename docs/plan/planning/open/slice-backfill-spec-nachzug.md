# Slice backfill-spec-nachzug: Spec-Nachzug — Pflichtenheft und Architektur-Sicht tragen den beschlossenen Backfill-Stand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Initial-Snapshot/Backfill des Bestands —
Haupt-Bezug), [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) (die offene, ADR-pflichtige Frage
„Backfill-Mechanismus"), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Folgepflicht 1 (Spec-Nachzug —
Träger dieses Slice).

**Berührte Spec-Stellen:** [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md), [`SPEC-001`](../../../../spec/pflichtenheft.md) (Tabellenliste),
[`SPEC-002`](../../../../spec/pflichtenheft.md) (`cdc.change`), [`SPEC-019`](../../../../spec/pflichtenheft.md) (Antrags-Datensatz),
[`SPEC-022`](../../../../spec/pflichtenheft.md) (`GET /changes`), eine **neue** Kennung [`SPEC-029`](../../../../spec/pflichtenheft.md) (Feldform
von `cdc.backfill_run`/`cdc.backfill_status`), [`ARC-006`](../../../../spec/architecture.md) (Driven-Adapter-Rolle)
und `spec/architecture.md` §4 (Sequenz). Der Verweis zeigt **aufwärts**: die
Spec nennt diesen Slice nie.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die beiden führenden Spec-Straten — Pflichtenheft (Rang 2) und
Architektur-Sicht (Rang 3) — tragen den mit [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) beschlossenen
Backfill-Stand als Technik-Festlegung und als Sequenz, **bevor** der erste
umsetzende Slice startet (die Modus-Deklaration in `harness/conventions.md` ist
Greenfield: die Doku führt). Umfang:

- (a) [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md): die Überschrift verliert das Wort „offen"; der Text
  legt Mechanismus (Bulk-Copy des Tabellenbestands in dem Snapshot eines je Run
  angelegten temporären logischen Slots, in einer Store-Transaktion
  committet), Markierung (`origin`), Überlappungs-Verhalten (alle Blöcke eines
  Runs liegen auf der Slot-Position `X`; begrenzte, idempotente Dopplung im
  Fenster `(Aktivierung, X]`, keine Lücke), Sichtbarkeits-Grenze (der Bestand
  ist ein Zustandsabzug für Consumer vor `X`) und Neubeginn nach Abbruch fest;
- (b) [`SPEC-002`](../../../../spec/pflichtenheft.md): Feld `origin` (`wal` | `backfill`; fehlender Wert liest als
  `wal`; letzte Spalte der View);
- (c) [`SPEC-019`](../../../../spec/pflichtenheft.md): Antragsart `backfill`, geschlossene Menge mit fünf Werten,
  die Bedeutung von `applied` bei dieser Antragsart („angenommen" — die
  Ausführung steht im Run-Zustand);
- (d) [`SPEC-022`](../../../../spec/pflichtenheft.md): Feld `origin` in der Antwort; die Anmerkung, dass ein
  Bestandsabzug **eine** Commit-Position teilt und ein `limit` innerhalb einer
  Position nicht fortsetzen kann;
- (e) [`SPEC-029`](../../../../spec/pflichtenheft.md) (neu): Feldform von `cdc.backfill_run` und
  `cdc.backfill_status`; [`SPEC-001`](../../../../spec/pflichtenheft.md) führt `cdc.backfill_run` in der
  Tabellenliste;
- (f) `spec/architecture.md`: eine Sequenz für den Backfill (Auslösung über die
  Antragsqueue, Übergabe an den Worker, Snapshot, ein Commit) und die
  Driven-Adapter-Rolle des Snapshot-Lesers in [`ARC-006`](../../../../spec/architecture.md) — **ohne**
  ADR-, Slice- oder Wellen-Bezug ([`AGENTS.md`](../../../../AGENTS.md) §3.4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Lastenheft** — bleibt unverändert; das Pflichtenheft präzisiert, erweitert
  nie (Kopf von `spec/pflichtenheft.md`).
- **Benutzerhandbuch** — beschreibt Betreiber-Oberfläche und darf keine
  Funktion nennen, die es noch nicht gibt; jeder Slice mit Betreiber-Oberfläche
  zieht seinen Abschnitt nach (`change-origin`, `sql-administration`, `e2e`,
  `bench-richtgroesse`, siehe Welle §4 Abweichung 4).
- **Code, Schema, Skripte** — Umsetzung gehört den Folge-Slices.
- **Transformationen** ([`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md)) — die Überschrift „… offen" von
  [`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md) und die Regeltypen in [`SPEC-019`](../../../../spec/pflichtenheft.md) bleiben Gegenstand
  der Transformations-Umsetzung; dieser Slice fügt [`SPEC-019`](../../../../spec/pflichtenheft.md) nur die
  Antragsart `backfill` hinzu und lässt Raum für weitere Werte (Welle §5, K3).
- **Eine Richtgröße für „große Tabellen"** — die Zahl entsteht erst aus der
  Messung in `bench-richtgroesse`; kein Wert ohne Messung im Pflichtenheft.

## 2. Definition of Done

- [ ] [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) ist beantwortet: die Überschrift trägt kein „offen"
      mehr, der Text nennt Mechanismus, Markierung, Überlappungs-Verhalten,
      Sichtbarkeits-Grenze und Neubeginn als Zusagen (Zukunfts-Form: was die
      Umsetzung liefern **muss**, nicht was gemessen wurde — die Zusagen sind
      bis zu den Test-Slices unbelegt und stehen nicht als geprüfte Aussage,
      [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) führt Lückenfreiheit und Replay-Invariante als
      „hergeleitet"). *Zu belegen durch:* Lesen des Abschnitts und
      `make docs-check`.
- [ ] Die Datenstrukturen stehen: [`SPEC-002`](../../../../spec/pflichtenheft.md) (`origin`), [`SPEC-019`](../../../../spec/pflichtenheft.md)
      (Antragsart `backfill`, fünf Werte), [`SPEC-022`](../../../../spec/pflichtenheft.md) (`origin`, Positions-
      Anmerkung), [`SPEC-001`](../../../../spec/pflichtenheft.md) (`cdc.backfill_run`) und [`SPEC-029`](../../../../spec/pflichtenheft.md) (Feldform
      `cdc.backfill_run`/`cdc.backfill_status`, die Kennung ist die nächste
      freie: höchste vergebene Kennung vor diesem Slice ist `SPEC-028` —
      *zu belegen durch* `grep -o 'SPEC-0[0-9][0-9]' spec/pflichtenheft.md | sort -u`
      am Parent-Stand); §7 Historie trägt je Änderung eine Zeile ohne ADR-/
      Slice-Bezug.
- [ ] `spec/architecture.md` trägt die Backfill-Sequenz und die
      Snapshot-Leser-Rolle in [`ARC-006`](../../../../spec/architecture.md), ohne ADR-/Slice-/Wellen-Bezug;
      die Aufzählungen der Antragsarten in derselben Datei sind auf die
      Menge mit `backfill` gezogen. *Zu belegen durch:* `make docs-check`
      (`matrix`-Modul) und der Suchlauf in §3.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt als eigener Punkt — der Slice **ist** das Doku-Update der Spec; das Benutzerhandbuch bleibt unberührt (§1).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/pflichtenheft.md` §1 (`LH-FA-CAP-009.a`) | update | Überschrift und Text auf den beantworteten Stand; Zusagen-Form. |
| `spec/pflichtenheft.md` §2 ([`SPEC-001`](../../../../spec/pflichtenheft.md), [`SPEC-002`](../../../../spec/pflichtenheft.md), [`SPEC-019`](../../../../spec/pflichtenheft.md), [`SPEC-022`](../../../../spec/pflichtenheft.md), [`SPEC-029`](../../../../spec/pflichtenheft.md)) | update / neu | Datenstrukturen; [`SPEC-029`](../../../../spec/pflichtenheft.md) neu (Datenstruktur, §2 — keine Breiten-Regel wie in §3). |
| `spec/pflichtenheft.md` §7 Historie | update | je Änderung eine Zeile, ohne ADR-/Slice-Bezug. |
| `spec/architecture.md` §1/§4 | update | Rolle des Snapshot-Lesers in [`ARC-006`](../../../../spec/architecture.md); Sequenz „Bestand als Backfill überführen"; Antragsarten-Aufzählung im bestehenden Sequenz-Abschnitt zur SQL-Aktivierung. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`LH-FA-CAP-009.a` ist beantwortet; die Antragsarten-Menge trägt `backfill`; `cdc.change` trägt `origin`"; beide Stände gemessen: Parent und Diff — der Implementer trägt Gefundenes und Nichtgefundenes je Zeile ein):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Anker auf die Überschrift „… offen" (`#lh-fa-cap-009a--backfill-mechanismus-offen`) in anderen Dokumenten | `grep -rn 'backfill-mechanismus' --include=*.md .` | *(Implementer trägt ein)* | Anker mitziehen; Records (`docs/reviews/**`, `done/**`) nur als Zitat-Korrektur nach [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) |
| Aufzählungen der Antragsarten („vier Antragsarten", `enable`/`disable`/`exclude_column`/`include_column`, „die drei übrigen Antragsarten") | `grep -rn 'exclude_column' spec docs/user harness README.md` | *(Implementer trägt ein)* | Spec-Stellen dieses Slice ziehen; Handbuch-Stellen an `sql-administration` melden (Welle §4) |
| Feldlisten von `cdc.change`/`GET /changes` mit Anzahl-Formulierung (z. B. „zwölf Felder") | `grep -rn 'Felder' spec/pflichtenheft.md spec/architecture.md` | *(Implementer trägt ein)* | Spec-Stellen dieses Slice ziehen; Handbuch an `change-origin` melden |
| `LH-FA-CAP-009` in Trägern der Abdeckung | `grep -rn 'CAP-009' docs harness .d-check.yml` | *(Implementer trägt ein)* | unverändert bis `e2e` (der Runner schreibt die Abdeckungs-Zeile) |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn die Welle [welle-backfill-bestand](../welle-backfill-bestand.md) eröffnet ist und
kein anderer Slice in `in-progress/` liegt (WIP-Limit 1). Dieser Slice
kommt zuerst: jeder Folge-Slice liest die hier festgelegten Zusagen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet —
  ein reiner Doku-Zug über zwei Dateien; sprengte er den Umfang, wäre die
  Architektur-Sicht (Punkt f) der abtrennbare Teil.
- `in-progress` → `open` (blockiert): falls das Übertragen eine Aussage von
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) als widersprüchlich zum Bestand der Spec zeigt — eine
  `Accepted` ADR wird nicht überschrieben ([`AGENTS.md`](../../../../AGENTS.md) §3.5), der Widerspruch
  ginge als Frage an den Architect und ggf. in eine Folge-ADR.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Spec-Straten dürfen keine ADR- und keine Slice-Kennung tragen**
  (`matrix`-Modul in `.d-check.yml`, [`AGENTS.md`](../../../../AGENTS.md) §3.4) — die Begründung der
  Entscheidung bleibt in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md). *Erwartet, zu belegen durch:*
  `make docs-check` grün und ein `grep -n 'ADR-0\|slice-\|welle-'` über beide
  Spec-Dateien am Diff. **Ausgang:** *(bei Closure)*
- **Zusagen als Tatsachen missverstanden** ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B):
  Lückenfreiheit, Replay-Invariante und GUC-Parität sind in
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) „hergeleitet" bzw. „erwartet"; das Pflichtenheft darf sie nur
  als Zusage an die Umsetzung formulieren. *Erwartet, zu belegen durch:* Review
  liest jeden Satz der neuen Abschnitte auf Zukunfts-Form. **Ausgang:** *(bei
  Closure)*
- **Kennungs-Vergabe [`SPEC-029`](../../../../spec/pflichtenheft.md)** gegen die Transformations-Umsetzung: sie
  erweitert [`SPEC-019`](../../../../spec/pflichtenheft.md) und braucht ggf. eigene Kennungen. Die Vergabe ist
  fortlaufend je Datei ([`SPEC-028`](../../../../spec/pflichtenheft.md) ist die höchste vorhandene); dieser Slice
  vergibt genau eine. *Erwartet, zu belegen durch:* die Suche im DoD-Punkt 2.
  **Ausgang:** *(bei Closure)*
- **Überholter Text im selben Dokument** (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`,
  offen, 2×): die Überschrift „… offen" und die Einleitungssätze der Abschnitte
  können nach dem Nachzug noch die alte Aussage tragen. *Erwartet, zu belegen
  durch:* Lesen beider Abschnitte von oben nach unten und Suchlauf §3, Zeile 1.
  **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Modus Greenfield, `harness/conventions.md` §Modus-Deklaration);
die Spec-Dateien bilden keine eigene Sub-Area — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (offen, 2×, einschlägig —
Risiko §6), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×,
Suchlauf §3), `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`
(verkörpert, 6×, einschlägig — Risiko §6, zweiter Punkt),
`BEO-PGC/zitat-nennt-die-falsche-stelle` (verkörpert, 7×, Review-Leser: jede
zitierte Stelle der Spec ist geprüft), `BEO-PGC/adr-folgepflicht-ohne-traeger-slice`
(offen, 1×, dieser Slice ist der Träger der Folgepflicht 1).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
