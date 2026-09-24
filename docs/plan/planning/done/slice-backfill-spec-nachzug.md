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
Träger dieses Slice), [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Folgepflicht 1 (Grants, Warn-Spalte(n), `estimated_rows` =
`NULL`, Bedeutung von `applied`, Annahme-Sequenz).

**Berührte Spec-Stellen:** [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md), [`SPEC-001`](../../../../spec/pflichtenheft.md) (Tabellenliste),
[`SPEC-002`](../../../../spec/pflichtenheft.md) (`cdc.change`), [`SPEC-019`](../../../../spec/pflichtenheft.md) (Antrags-Datensatz),
[`SPEC-022`](../../../../spec/pflichtenheft.md) (`GET /changes`), eine **neue** Kennung [`SPEC-029`](../../../../spec/pflichtenheft.md) (Feldform
von `cdc.backfill_run`/`cdc.backfill_status`), [`ARC-006`](../../../../spec/architecture.md) (Driven-Adapter-Rolle)
und `spec/architecture.md` §4 (Sequenz). Der Verweis zeigt **aufwärts**: die
Spec nennt diesen Slice nie.

**Verantwortlich:** Implementer-Agent, 2026-09-24.

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
  Run-Zeile entsteht in derselben Transaktion, die Ausführung steht im
  Run-Zustand);
- (d) [`SPEC-022`](../../../../spec/pflichtenheft.md): Feld `origin` in der Antwort; die Anmerkung, dass ein
  Bestandsabzug **eine** Commit-Position teilt und ein `limit` innerhalb einer
  Position nicht fortsetzen kann;
- (e) [`SPEC-029`](../../../../spec/pflichtenheft.md) (neu): Feldform von `cdc.backfill_run` und
  `cdc.backfill_status` samt den **Grants** (`cdc_admin` `SELECT`/`INSERT`,
  `cdc_capture` `SELECT`/`UPDATE`, `cdc_reader` kein Recht auf die Basistabelle,
  `SELECT` auf die View), den **Warn-Spalte(n)** für die beiden Warnungen (Zahl und
  Bezeichner legt dieser Slice fest; `run-store` folgt) — ohne Toleranz und ohne
  Richtgröße — und `estimated_rows` = `NULL` als „unbekannt", nie `0`, jede Stelle
  mit dem Wort „geschätzt"; [`SPEC-001`](../../../../spec/pflichtenheft.md) führt `cdc.backfill_run` in der
  Tabellenliste;
- (f) `spec/architecture.md`: eine Sequenz für den Backfill (Auslösung über die
  Antragsqueue, **Annahme** — Run-Zeile `queued` und Antragsvermerk in einer
  Transaktion —, Wecken des Workers und Aufnahme beim Prozessstart, Snapshot, ein
  Commit) und die Driven-Adapter-Rollen des Snapshot-Lesers und der Annahme in
  [`ARC-006`](../../../../spec/architecture.md) — **ohne** ADR-, Slice- oder Wellen-Bezug ([`AGENTS.md`](../../../../AGENTS.md) §3.4).

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
- **Toleranz und Richtgröße für „große Tabellen"** — die Richtgröße entsteht erst
  aus der Messung in `bench-richtgroesse`, die Toleranz ist ein Startwert im Code
  ohne §3-Eintrag; kein Wert ohne Messung im Pflichtenheft, [`SPEC-029`](../../../../spec/pflichtenheft.md) trägt nur die
  Spalte(n) der Warnungen.

## 2. Definition of Done

- [x] [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) ist beantwortet: die Überschrift trägt kein „offen"
      mehr, der Text nennt Mechanismus, Markierung, Überlappungs-Verhalten,
      Sichtbarkeits-Grenze und Neubeginn als Zusagen (Zukunfts-Form: was die
      Umsetzung liefern **muss**, nicht was gemessen wurde — die Zusagen sind
      bis zu den Test-Slices unbelegt und stehen nicht als geprüfte Aussage,
      [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) führt Lückenfreiheit und Replay-Invariante als
      „hergeleitet"). *Zu belegen durch:* Lesen des Abschnitts und
      `make docs-check`.
- [x] Die Datenstrukturen stehen: [`SPEC-002`](../../../../spec/pflichtenheft.md) (`origin`), [`SPEC-019`](../../../../spec/pflichtenheft.md)
      (Antragsart `backfill`, fünf Werte), [`SPEC-022`](../../../../spec/pflichtenheft.md) (`origin`, Positions-
      Anmerkung), [`SPEC-001`](../../../../spec/pflichtenheft.md) (`cdc.backfill_run`) und [`SPEC-029`](../../../../spec/pflichtenheft.md) (Feldform
      `cdc.backfill_run`/`cdc.backfill_status` samt Grants, Warn-Spalte(n) und
      `estimated_rows` = `NULL` als „unbekannt"; [`SPEC-019`](../../../../spec/pflichtenheft.md) nennt `applied`
      bei `backfill` „angenommen"; die Kennung ist die nächste
      freie: höchste vergebene Kennung vor diesem Slice ist `SPEC-028` —
      *zu belegen durch* `grep -o 'SPEC-0[0-9][0-9]' spec/pflichtenheft.md | sort -u`
      am Parent-Stand); §7 Historie trägt je Änderung eine Zeile ohne ADR-/
      Slice-Bezug.
- [x] `spec/architecture.md` trägt die Backfill-Sequenz (einschließlich der
      Annahme in einer Transaktion und der Aufnahme beim Start) und die Rollen des
      Snapshot-Lesers und der Annahme in [`ARC-006`](../../../../spec/architecture.md), ohne ADR-/Slice-/Wellen-Bezug;
      die Aufzählungen der Antragsarten in derselben Datei sind auf die
      Menge mit `backfill` gezogen. *Zu belegen durch:* `make docs-check`
      (`matrix`-Modul) und der Suchlauf in §3.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
      Report `docs/reviews/review-slice-backfill-spec-nachzug.md` (0 HIGH, 1 MEDIUM,
      3 LOW, 2 INFO); F-1 bis F-4 in der Fixrunde behoben (§3, Zeilen „Fixrunde
      nach Review"), kein offenes HIGH/MEDIUM.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt als eigener Punkt — der Slice **ist** das Doku-Update der Spec; das Benutzerhandbuch bleibt unberührt (§1).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke). *(§7: geschärfte Regel — Suchlauf
      nach dem Hedge des offenen Punkts; benannte Lücken.)*
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert. *(Zwei weitere `evidence/`-Dateien,
      je `slice-backfill-spec-nachzug.md`: `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
      und `BEO-PGC/arbeit-ueberholt-stehenden-traeger`; siehe §7.)*
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen). *(Sechs Ausgänge in §6, je am Ort; siehe §7.)*
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/pflichtenheft.md` §1 (`LH-FA-CAP-009.a`) | update | Überschrift und Text auf den beantworteten Stand; Zusagen-Form. |
| `spec/pflichtenheft.md` §2 ([`SPEC-001`](../../../../spec/pflichtenheft.md), [`SPEC-002`](../../../../spec/pflichtenheft.md), [`SPEC-019`](../../../../spec/pflichtenheft.md), [`SPEC-022`](../../../../spec/pflichtenheft.md), [`SPEC-029`](../../../../spec/pflichtenheft.md)) | update / neu | Datenstrukturen; [`SPEC-029`](../../../../spec/pflichtenheft.md) neu (Datenstruktur, §2 — keine Breiten-Regel wie in §3) mit Grants und Warn-Spalte(n). |
| `spec/pflichtenheft.md` §7 Historie | update | je Änderung eine Zeile, ohne ADR-/Slice-Bezug. |
| `spec/architecture.md` §1/§4 | update | Rolle des Snapshot-Lesers in [`ARC-006`](../../../../spec/architecture.md); Sequenz „Bestand als Backfill überführen"; Antragsarten-Aufzählung im bestehenden Sequenz-Abschnitt zur SQL-Aktivierung. |
| `spec/pflichtenheft.md` §2 ([`SPEC-020`](../../../../spec/pflichtenheft.md), [`SPEC-021`](../../../../spec/pflichtenheft.md)) | update (Plan-Nachzug im Lauf) | Fund des §3.13-Suchlaufs (Zeile 3): „dieselben Felder wie der Domain-Typ `model.Change`" bleibt nach der Einführung von `origin` an `model.Change` nur wahr, wenn die Live-Nachricht das Feld ausdrücklich ausnimmt — je ein Halbsatz, keine Änderung der Nachrichtenschemata. |
| `spec/architecture.md` Kopf und Komponenten-Diagramm | update (Plan-Nachzug im Lauf) | der Frische-Marker „Letzte Änderung" und der Knoten „Driven Adapters" des Mermaid-Diagramms tragen die neuen Rollen aus [`ARC-006`](../../../../spec/architecture.md) mit; die Sicht bliebe sonst gegen ihre eigene Tabelle uneinheitlich. |
| `spec/pflichtenheft.md` §2 [`SPEC-029`](../../../../spec/pflichtenheft.md) — **Festlegung der Warn-Spalten** (Auftrag aus §1 (e)) | Festlegung | **Zwei** Spalten vom Typ `boolean`, `NOT NULL`, Default `false`: `warn_estimated_size` (Warnung 1: die **geschätzte** Zeilenzahl liegt über der Richtgröße; schreibt die Annahme) und `warn_duration` (Warnung 2: Kopierdauer über der Toleranz; schreibt der Worker). Grund für zwei statt einer: die beiden Schreiber sind verschiedene Rollen mit verschiedenem Grant (`cdc_admin` nur `INSERT`, `cdc_capture` nur `UPDATE`) — eine gemeinsame Spalte hätte zwei Schreiber für einen Wert. Grund für `boolean` statt Text/Zahl: die Spalten tragen nur das Ergebnis, weder Toleranz noch Richtgröße (`ADR-0113` Festlegung 3 Punkt 4: sie stehen an genau einer Stelle im Code). `run-store` folgt diesen Namen und Typen. |
| `spec/pflichtenheft.md` §2 [`SPEC-022`](../../../../spec/pflichtenheft.md) Zeile „Reihenfolge" (Fixrunde nach Review, F-1) | update (Plan-Nachzug im Lauf) | die unbedingte Fortsetzungs-Regel `from = <letzte gelieferte commit_position> + 1` trägt den Vorbehalt „wenn das Lesen die letzte Position vollständig erfasst hat" und verweist auf die Zeile „Position und `limit`"; die Zeilen widersprechen sich sonst im selben Abschnitt. |
| `spec/pflichtenheft.md` §2 [`SPEC-029`](../../../../spec/pflichtenheft.md) (Absatz nach der Spaltentabelle) und `spec/architecture.md` §4 (Annahme in einer Transaktion) (Fixrunde nach Review, F-3) | update (Plan-Nachzug im Lauf) | zwei Begründungssätze im Konjunktiv über die verworfene Alternative („hätte …", „sonst … stünde") stehen im Indikativ über den geltenden Zustand (je Spalte genau eine schreibende Rolle; Annahme und Antragsvermerk gemeinsam oder gar nicht). |
| `spec/architecture.md` §4 (Ausführung eines Runs, Sichtbarkeit) (Fixrunde nach Review, F-4) | update (Plan-Nachzug im Lauf) | „Bis zum einen Commit ist nichts für Leser sichtbar" widersprach dem sichtbaren `rows_copied` aus [`SPEC-029`](../../../../spec/pflichtenheft.md); der Satz trennt jetzt den sichtbaren Run-Fortschritt von den erst mit dem einen Commit sichtbaren Changes. |
| `docs/plan/planning/open/slice-backfill-run-store.md`, `…/slice-backfill-run-usecase.md`, `…/slice-backfill-sql-administration.md` (Fixrunde nach Review, F-2) | update (Meldung, minimal mitgezogen) | „Warn-Spalte(n)/Warn-Kennzeichnung leer, solange keine Auswertung sie setzt" → „`false`, …" (Text der Absicht unverändert; [`SPEC-029`](../../../../spec/pflichtenheft.md) legt `boolean NOT NULL DEFAULT false` fest). |
| `docs/plan/planning/open/slice-transformationen-spec-nachzug.md`, `…/slice-transformationen-antragsweg-schema.md` (Fixrunde nach Review, F-2) | update (Meldung, minimal mitgezogen) | „nennt die vier SQL-Funktionen"/„Satz ‚vier Arten'"/„den vier bestehenden Funktionen" → ohne Zahl („die bestehenden …"), weil [`SPEC-019`](../../../../spec/pflichtenheft.md) fünf Antragsarten führt; die Absicht beider Pläne (zwei weitere Arten) bleibt unberührt. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`LH-FA-CAP-009.a` ist beantwortet; die Antragsarten-Menge trägt `backfill`; `cdc.change` trägt `origin`"; beide Stände gemessen: Parent und Diff — der Implementer trägt Gefundenes und Nichtgefundenes je Zeile ein):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Anker auf die Überschrift „… offen" (`#lh-fa-cap-009a--backfill-mechanismus-offen`) in anderen Dokumenten | `grep -rn 'backfill-mechanismus' --include=*.md .` | **Parent** (`77c60bd1`): 1 Treffer — die Suchbefehl-Zelle dieser Tabelle selbst (Zitat des Ankers in Inline-Code, kein Link); die Überschrift `### LH-FA-CAP-009.a — Backfill-Mechanismus offen` stand in `spec/pflichtenheft.md` Z. 144. **Diff:** dieselbe 1 Treffer (die eigene Zelle), die Überschrift lautet `### LH-FA-CAP-009.a — Backfill-Mechanismus`. **Nicht gefunden:** jeder Link auf den Anker; ergänzender Lauf `grep -rnE 'cap-009a\|CAP-009\.a' --include=*.md --include=*.yml --include=*.yaml .` außerhalb der Backfill-Plan-/ADR-Dateien: nur die nackte Kennung ohne Anker in [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Z. 393, 411) und in der Historie des Pflichtenhefts. | keine Anker zu ziehen; `ADR-0112` ist `Accepted` und trägt die Kennung ohne Anker, unberührt. |
| Aufzählungen der Antragsarten („vier Antragsarten", `enable`/`disable`/`exclude_column`/`include_column`, „die drei übrigen Antragsarten") | `grep -rn 'exclude_column' spec docs/user harness README.md` | **Parent:** 9 Treffer (`spec/architecture.md` Z. 228; `spec/pflichtenheft.md` Z. 403, 414, 433, 658, 661; `docs/user/benutzerhandbuch.md` Z. 259, 1250; `harness/README.md` Z. 141). **Diff:** 9 Treffer, dieselbe Menge, Zeilennummern in `spec/pflichtenheft.md` verschoben. **Gefunden und gezogen:** `spec/architecture.md` („vier Arten" → „fünf", Tabelle + `backfill`, „die drei übrigen Antragsarten" umformuliert, „beiden Tabellen-Antragsarten" → `enable`/`disable`); `spec/pflichtenheft.md` `SPEC-019` (Funktionsliste, `request_kind`-Menge, `column_name`-Zelle „die drei übrigen", „beiden Tabellen-Antragsarten" → `enable`/`disable`). **Nicht gefunden:** keine weitere Zählformulierung („vier/fünf Antragsarten") in `spec`, `docs/user`, `harness`, `README.md`, `README.de.md`, `AGENTS.md`. **Nicht gezogen (Records):** Historie-Zeilen Z. 658/661 (Zeitpunkt-Records, wahr zu ihrem Datum). | Spec-Stellen gezogen; **Meldung an `slice-backfill-sql-administration`:** `docs/user/benutzerhandbuch.md` nennt `backfill` nicht; Z. 302 „Tabellen-Antragsarten (`enable`/`disable`)" bleibt wahr; `harness/README.md` Z. 141 beschreibt den Spaltenausschluss-Rundlauf, unberührt. |
| Feldlisten von `cdc.change`/`GET /changes` mit Anzahl-Formulierung (z. B. „zwölf Felder") | `grep -rn 'Felder' spec/pflichtenheft.md spec/architecture.md` | **Parent:** 11 Treffer in `spec/pflichtenheft.md`, 0 in `spec/architecture.md`. **Diff:** 11 (Anzahl unverändert; Wortlaut an zwei Stellen erweitert). **Gefunden:** „zehn Felder" in `SPEC-021` und `SPEC-024` und „dieselben Felder wie der Domain-Typ `model.Change`" in `SPEC-020` — Live-Nachrichten, die `origin` nicht tragen (Zusage in `SPEC-002`); die Aussage „wie der Domain-Typ" würde durch das neue Feld an `model.Change` falsch → `SPEC-020`/`SPEC-021` um die Ausnahme ergänzt (`SPEC-024` verweist auf `SPEC-021`s Feldliste und bleibt wahr). **Nicht gefunden:** eine Anzahl-Formulierung für `cdc.change` oder für die `GET /changes`-Antwort in `spec/`. | Handbuch an `slice-backfill-change-origin` melden: `docs/user/benutzerhandbuch.md` Z. 639 (Feldliste der `GET /changes`-Antwort) und Z. 377 (Spaltenliste der View-Abfrage) tragen `origin` noch nicht; die „zehn Feldern" der Live-Wege (Z. 721 ff.) bleiben wahr. Die Fortsetzungs-Idiome `commit_position > <letzte-gelesene-position>` (Z. 375 ff.) und `from = <letzte> + 1` (Z. 641 ff.) tragen die neue Position-und-`limit`-Regel aus `SPEC-022` noch nicht — Meldung an `slice-backfill-sql-administration` (Handbuch „Lesen ohne `Limit`"). |
| Aufzählungen der Views und Rollen-Grants in `spec/architecture.md` ([`ARC-005`](../../../../spec/architecture.md) „SQL-Funktionen/Views", [`ARC-006`](../../../../spec/architecture.md) Adapter-Liste) | `grep -n 'ARC-005\|ARC-006' spec/architecture.md` | **Parent:** 20 Zeilen (Komponenten-Tabelle Z. 58/59, Schichten-Tabelle Z. 77/78, Teilnehmer der bestehenden Sequenzen). **Diff:** 24 Zeilen (die vier zusätzlichen sind Teilnehmer der neuen Sequenzen). **Gefunden und gezogen:** [`ARC-006`](../../../../spec/architecture.md) (Adapterliste um Snapshot-Leser und Annahme) und der Knoten „Driven Adapters" im Diagramm. **Nicht gefunden:** eine Aufzählung der Views oder der Rollen-Grants in `spec/architecture.md` — [`ARC-005`](../../../../spec/architecture.md) nennt „SQL-Funktionen/Views" in Sammelform, die Schichten-Tabelle (Z. 77/78) ist rollen-, nicht adapterlistenbezogen. | keine weitere Stelle zu ziehen; die Grants stehen in `SPEC-029`, nicht in der Sicht. |
| `LH-FA-CAP-009` in Trägern der Abdeckung | `grep -rn 'CAP-009' docs harness .d-check.yml` | **Parent** (`git grep` am Stand `77c60bd1`, ohne die Backfill-Plan-/ADR-Dateien): 13 Zeilen — Roadmap, Transformations-Plan und -Welle, ADR-Index, `ADR-0112`, drei Verifikationsberichte (Records), `harness/README.md` (`make doc-trace`-Zeile mit der Waisen-Momentaufnahme). **Diff:** dieselbe Menge — der Diff berührt keine dieser Dateien. **Nicht gefunden:** ein Träger in `docs/user/` (weder Abdeckungstabelle noch Handbuch) und in `.d-check.yml` (`trace.coverage`). | unverändert bis `e2e` (der Runner schreibt die Abdeckungs-Zeile); die Waisen-Momentaufnahme in `harness/README.md` gehört dem Slice, der sie ändert. |
| Fortsetzungs-Aussage „`from = <letzte> + 1`" (Fixrunde nach Review, F-1; bewegte Eigenschaft: die Fortsetzung über eine Position gilt nur bei vollständig erfasster Position) | `git grep -nE 'letzte gelieferte\|letzte.*commit_position\|\+ 1' -- spec` (Parent: `git grep … 77c60bd1 -- spec`) | **Parent** (`77c60bd1`): 1 Treffer — `spec/pflichtenheft.md` Z. 506 (`SPEC-022` „Reihenfolge"). **Diff** (Stand `88bee004`, vor der Fixrunde): 2 Treffer — Z. 614 („Reihenfolge", unbedingt) und Z. 618 („Position und `limit`", mit Grenze). **Nach der Fixrunde:** dieselben 2 Zeilen, Z. 614 trägt den Vorbehalt und den Verweis auf die Zeile „Position und `limit`". **Nicht gefunden:** dieselbe Aussage in `spec/architecture.md` (0 Treffer; dort nur „Position fortsetzen (regulär nur vorwärts)", Z. 182, eine Consumer-Position, keine Lese-Fortsetzung) und an anderer Stelle des Pflichtenhefts (Z. 216 sagt bereits „ein `limit` kann innerhalb einer Position nicht fortsetzen"). | Z. 614 gezogen. Handbuch Z. 379/642 bleiben bei der Meldung an `slice-backfill-sql-administration` (Zeile oben). |
| Plan-Träger der Warn-Spalten und der Antragsarten-Menge (Fixrunde nach Review, F-2; bewegte Eigenschaften: Warn-Spalten sind `boolean NOT NULL DEFAULT false` statt „leer"; fünf Antragsarten/SQL-Funktionen statt vier) | `git grep -cE 'leer, solange\|das leere Feld\|Spalte\(n\) leer\|Spalte\(n\) \(leer\)\|vier SQL-Funktionen\|vier Arten\|vier bestehenden\|die vier$' <Stand> -- docs/plan/planning/open 'docs/plan/planning/welle-*.md'` | **Parent** (`77c60bd1`) und **Diff-Stand** (`88bee004`): je 9 Zeilen in 5 Dateien — `slice-backfill-run-store` 1, `slice-backfill-run-usecase` 2, `slice-backfill-sql-administration` 3, `slice-transformationen-antragsweg-schema` 1, `slice-transformationen-spec-nachzug` 2; zusätzlich eine Umbruch-Form („der vier / bestehenden", `slice-transformationen-antragsweg-schema` Z. 200), die das Muster nicht trifft. **Nach der Fixrunde:** 0 Treffer, Umbruch-Form ebenfalls gezogen. **Nicht gefunden:** eine Zählung der Antragsarten oder eine Warn-Spalten-Beschreibung in `welle-backfill-bestand.md`, `slice-backfill-bench-richtgroesse.md`, `slice-backfill-snapshot-reader.md`, `slice-backfill-change-origin.md`, `slice-backfill-e2e.md`, `slice-backfill-sdk-origin.md`, `slice-backfill-row-image-gemeinsam.md`; „zehn Felder" in `slice-transformationen-e2e-wirkung`/`-spec-nachzug`, `slice-backfill-sdk-origin` und `welle-transformationen` beschreiben die Live-Nachrichten und bleiben wahr; „die zwei weiteren Arten" der Transformations-Pläne bleibt wahr (relativ zur Bestandsmenge). | in fünf Fremd-Plänen minimal gezogen (Zeilen in §3 oben); den Fund im Slice `slice-backfill-run-usecase` (Z. 42, 87), den der Review nicht nannte, trägt derselbe Suchlauf. Die Absichten der Pläne bleiben unberührt. |

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
  Spec-Dateien am Diff. **Ausgang:** *entfallen* — nicht eingetreten:
  `make docs-check` lief im Verifier-Lauf mit 0 Befunden, und das `grep` über
  beide Spec-Dateien nach ADR-, Slice- und Wellen-Kennungen fand 0 Treffer
  (Verifikation §2 Zeile 3).
- **Zusagen als Tatsachen missverstanden** ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B):
  Lückenfreiheit, Replay-Invariante und GUC-Parität sind in
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) „hergeleitet" bzw. „erwartet"; das Pflichtenheft darf sie nur
  als Zusage an die Umsetzung formulieren. *Erwartet, zu belegen durch:* Review
  liest jeden Satz der neuen Abschnitte auf Zukunfts-Form. **Ausgang:**
  *entfallen* — nicht eingetreten: der Review nennt die Zusagen-Form im
  Negativbefund (Lückenfreiheit und Replay-Invariante stehen als Zusage, nicht
  als geprüfte Aussage), und der Verifier las den Abschnitt gegen die
  Entscheidungen (Verifikation §4, Zeile „Mechanismus, Markierung, Position …").
- **Zahl und Bezeichner der Warn-Spalte(n)** legt dieser Slice in [`SPEC-029`](../../../../spec/pflichtenheft.md) fest;
  `run-store` folgt ihnen. Weicht das Schema davon ab, ist es ein Plan-Nachzug
  dieses Slice, kein Alleingang des Schemas. *Erwartet, zu belegen durch:* der
  Abgleich von [`SPEC-029`](../../../../spec/pflichtenheft.md) und `tools/schema/schema.yaml` im Review von `run-store`.
  **Ausgang:** *entfallen* — die Zahl (zwei) und die Bezeichner
  (`warn_estimated_size`, `warn_duration`, `boolean NOT NULL DEFAULT false`)
  stehen in [`SPEC-029`](../../../../spec/pflichtenheft.md) fest (Verifikation
  §4, erste Zeile); der Abgleich mit dem Schema ist DoD-Gegenstand von
  `slice-backfill-run-store` (die Schema-und-Grants-Zeile trägt „zwei
  Warn-Spalten nach [`SPEC-029`](../../../../spec/pflichtenheft.md)", auf die festgelegte Form gezogen im Zug
  `95bd14a8`).
- **Kennungs-Vergabe [`SPEC-029`](../../../../spec/pflichtenheft.md)** gegen die Transformations-Umsetzung: sie
  erweitert [`SPEC-019`](../../../../spec/pflichtenheft.md) und braucht ggf. eigene Kennungen. Die Vergabe ist
  fortlaufend je Datei ([`SPEC-028`](../../../../spec/pflichtenheft.md) ist die höchste vorhandene); dieser Slice
  vergibt genau eine. *Erwartet, zu belegen durch:* die Suche im DoD-Punkt 2.
  **Ausgang:** *entfallen* — nicht eingetreten: `grep -o 'SPEC-0[0-9][0-9]'
  spec/pflichtenheft.md | sort -u` liefert am Parent (`77c60bd1`) als höchste
  Kennung `SPEC-028` und am Diff `SPEC-029`; genau eine Kennung vergeben, und
  die offenen Transformations-Pläne führen `SPEC-030` nicht.
- **Überholter Text im selben Dokument** (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`,
  offen, 2×): die Überschrift „… offen" und die Einleitungssätze der Abschnitte
  können nach dem Nachzug noch die alte Aussage tragen. *Erwartet, zu belegen
  durch:* Lesen beider Abschnitte von oben nach unten und Suchlauf §3, Zeile 1.
  **Ausgang:** *eingetreten* — nicht an der Überschrift (Suchlauf §3, Zeile 1:
  kein Link auf den Anker, Überschrift ohne „offen"), sondern an `SPEC-022`:
  die unbedingte Fortsetzungs-Regel der Zeile „Reihenfolge" stand neben der
  neuen Zeile „Position und `limit`" (Review F-1, MEDIUM). Im Slice behoben
  (Fixrunde `1034840f`), kein Carveout und kein Folge-Slice nötig; die Klasse
  trägt jetzt den dritten Beleg im Register (§7).

- **Übergabe aus dem Review** (Deutungsfragen, keine Änderung in diesem Slice):
  F-5 (die zwei Warnungen für große Tabellen tragen im Pflichtenheft keinen
  Lastenheft-Anker; Präzisierung oder neue Zusage?) — an Verifier/Architect
  gemeldet. F-6 (die Architektur-Sicht beschreibt Annahme, Worker und
  Snapshot-Leser im Indikativ, obwohl die Umsetzung in den Folge-Slices liegt)
  — an Verifier/Architect gemeldet. **Ausgang:** *entfallen* — der Verifier
  urteilt beide als kein Befund (Verifikation §6): F-5 ist eine zulässige
  Präzisierung von [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (die Warnung ändert keinen Ausgang der
  drei Akzeptanzkriterien), F-6 folgt dem Greenfield-Stil der Sicht, ein
  Zukunftsmarker wäre nach [`AGENTS.md`](../../../../AGENTS.md) §3.4 regelwidrig. Der Re-Evaluierungs-Trigger für F-5
  (die Warnung wird zum Vertrag) liegt bei [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md).

## 7. Closure-Notiz

- **Was hat funktioniert:** der Zug blieb ein reiner Doku-Zug über die zwei
  Spec-Dateien ohne Rückführung (§4): Lastenheft, Benutzerhandbuch, Code und
  Schema sind unverändert (`git diff --stat`, Verifikation §3). Die
  Rollen-Kette lief unabhängig — Review (Summary aus dem Report übernommen:
  0 HIGH · 1 MEDIUM · 3 LOW · 2 INFO), Fixrunde für F-1 bis F-4 (`1034840f`,
  `6d6ee7ea`), Verifikation „Bestätigt": sieben `[x]`-Zeilen je mit eigenem
  Beleg, die Spec gegen [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) und
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) konform, an den zwei überholten Stellen der älteren ADR
  (Grants, `applied` = „angenommen") folgt die Spec der jüngeren. Der
  Verifier maß das §3.13-Suchlauf-Feld an beiden Ständen nach (sieben Zeilen,
  alle bestätigt) und fuhr `make gates` real (EXIT=0, Stempel gleich —
  Verifikation §1). Die Zusagen-Form trug: Lückenfreiheit und
  Replay-Invariante stehen als Zusage an die Umsetzung, nicht als geprüfte
  Tatsache.
- **Was ging anders als geplant:** der Plan wuchs im Lauf um drei Zeilen
  (`SPEC-020`/`SPEC-021`: die Live-Nachricht nimmt `origin` aus; Kopf und
  Diagramm der Sicht; die Warn-Spalten-Festlegung — der Auftrag „Zahl und
  Bezeichner legt dieser Slice fest" ist mit zwei Spalten eingelöst). Die
  Fixrunde zog F-2 direkt in fünf Fremd-Plänen statt nur zu melden
  ([`AGENTS.md`](../../../../AGENTS.md) §3.13: „gemeldet statt still mitgeändert" — hier nicht still, die
  Plan-Zeilen und die Suchlauf-Zeile 7 nennen es). Danach maß der Verifier
  einen Rest: das Muster der Fixrunde suchte die „leer"-Formen, nicht den
  Hedge „Warn-Spalte(n)" (V-1) — im Planner-Zug `95bd14a8` gezogen.
- **Verifier-Beobachtungen (V-1 bis V-3):** *V-1* gezogen: 17 Zeilen in fünf
  Plan-Dateien (`slice-backfill-run-store`, `-run-usecase`,
  `-sql-administration`, `-bench-richtgroesse`, `welle-backfill-bestand`) am
  Stand vor `95bd14a8`, 0 danach (`git grep -cE 'Warn-Spalte\(n\)|Spalte\(n\)|ein Feld für die Warn|das Feld der Warn' -- docs/plan/planning`,
  ohne die Datei dieses Slice); der Verifier nannte 14 Zeilen in vier Dateien,
  die übrigen drei trägt dasselbe Muster; die Absichten der Pläne bleiben
  unberührt. Dieser Plan selbst behält „Warn-Spalte(n)" in §1 und §2 als
  Wortlaut des ursprünglichen Auftrags; die Festlegung steht in §3. *V-2*
  ist ein **Fehlalarm**: `harness/README.md` (Zeile `make doc-trace`) nennt
  bereits „80 Anforderungen, **3 Waisen**" — am Parent und an `HEAD`
  (`git grep -o '[0-9]* Anforderungen, \*\*[0-9]* Waisen\*\*' <Stand> -- harness/README.md`,
  seit `712dc26d`); nichts zu ändern, der Verifikations-Report bleibt als
  Record unberührt. *V-3* ohne Handlungsbedarf: die Messung der Startposition
  eines frisch registrierten Consumers trägt `slice-backfill-e2e`.
- **Steering-Loop-Eintrag (Lerneintrag):** geschärfte Regel (Anwendungs-Schärfung
  der verkörperten Klasse `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, Anker
  [`AGENTS.md`](../../../../AGENTS.md) §3.13): *ein Slice, der eine bis dahin offene Zahl oder Menge festlegt,
  sucht den Hedge des offenen Punkts („(n)", „offen", „noch nicht") über alle
  Plan-Träger — nicht nur den Vorzustand des Werts, den er ersetzt.* Beleg:
  die Träger sagten „Warn-Spalte(n)", und erst die Suche nach dem Hedge
  (Verifikation V-1) fand sie; die „leer"-Suche der Fixrunde hatte nur einen
  Teil getroffen, der committete Suchlauf des Implementers `docs/plan/planning`
  gar nicht (Review F-2). Zweitens erreicht `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
  mit diesem Slice **3×** (die Tabellenzeile „Reihenfolge" von `SPEC-022`
  trug die Fortsetzungs-Regel ohne den Vorbehalt der neuen Zeile „Position
  und `limit`", Review F-1) und wandert in die Steering-Loop-Einträge der
  Closure von [welle-backfill-bestand](../welle-backfill-bestand.md); die Regelschärfungs-Frage bleibt eine
  Architect-Entscheidung. Kein neuer Sensor: der Reviewer und der Verifier
  fanden die Funde als Leser, kein Gate liest Prosa-Kohärenz. Benannte
  Spec-Lücken (jeweils ein Träger außerhalb dieses Slice, der Zug ist
  gemeldet, nicht gezogen): die Handbuch-Träger `docs/user/benutzerhandbuch.md`
  (Feldliste der `GET /changes`-Antwort und Spaltenliste der View-Abfrage ohne
  `origin` → `slice-backfill-change-origin`; Fortsetzungs-Idiome ohne die
  Position-und-`limit`-Regel → `slice-backfill-sql-administration`; beide
  Meldungen stehen im Suchlauf-Feld §3, Zeile 3) und die Startposition eines
  frisch registrierten Consumers, in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) „erwartet, nicht geprüft"
  (→ `slice-backfill-e2e`).
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`** — weitere Datei
    `evidence/slice-backfill-spec-nachzug.md` (Review F-1); Zähler **3×**
    (Datei-Anzahl unter `evidence/`, real ausgezählt), Schwelle erreicht,
    Ausgang noch nicht zugewiesen — `state.md` trägt den Übergang in den
    Lese-Schritt der Closure von [welle-backfill-bestand](../welle-backfill-bestand.md).
  - **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** — weitere Datei
    `evidence/slice-backfill-spec-nachzug.md` (Review F-2, Verifikation
    V-1); Zähler **27×** (real ausgezählt), Ausgang bleibt verkörpert.
  - **F-3 (Konjunktiv über die verworfene Alternative in Spec-Prosa)** —
    benannt, nicht angelegt: keine Register-Klasse trägt Doku-Prosa
    (`BEO-PGC/slice-chronik-in-code-kommentar` ist auf Code-Kommentare
    skopiert), und andere Reports nennen Konjunktiv nur als bestandene
    Negativprüfung; ein Auftreten, kein zweites.
  - **F-4 (unpräzise Sichtbarkeitsaussage neben einer Gegenstelle)** und der
    **Host-Python-Einsatz** für Markdown-Ersetzungen (Review: ohne Befund,
    kein Artefakt im Repo) — benannt, nicht angelegt: erstes Auftreten ohne
    Klasse im Register.
  - **F-5, F-6, V-2, V-3** — kein Befund bzw. Fehlalarm (siehe §6 und oben),
    keine Beobachtung angefallen.
- **Validator (Modul 8):** entfällt ausdrücklich — der Slice ist ein reiner
  Doku-/Spec-Nachzug ohne End-Nutzer-Wert; der Nutzer-Bedarf (`LH-FA-CAP-009`
  laut Lastenheft) wird erst durch die Umsetzung der Folge-Slices und den
  Wellen-Beleg (`slice-backfill-e2e`) validierbar. Kein stilles Überspringen.
- **Closure-Notiz-Review (`.harness/skills/closure-note-reviewer.md`):**
  eine getrennte Rolle im frischen Kontext, kein Schritt der Planner-Closure;
  der Skill prüft Slices in `done/` (Kontext-Eingang: alle Closure-Notizen
  dort) und greift daher erst nach dem `git mv` — hier nicht ausgeführt.
- **Folge-Slices:** keine neuen — die neun Folge-Slices der Welle
  [welle-backfill-bestand](../welle-backfill-bestand.md) liegen als Dateien in `open/`; die zwei Handbuch-Meldungen
  gehen an `slice-backfill-change-origin` und `slice-backfill-sql-administration`.
- **Risiken aus §6:** je ein Ausgang am Ort — Risiko 1 (Spec-Straten ohne
  Kennung) **entfallen**; Risiko 2 (Zusagen als Tatsachen) **entfallen**;
  Risiko 3 (Warn-Spalten) **entfallen** (festgelegt; Schema-Abgleich beim
  Folge-Slice `run-store`); Risiko 4 (Kennungs-Vergabe) **entfallen**;
  Risiko 5 (überholter Text) **eingetreten**, im Slice behoben, Klasse im
  Register; Übergabe F-5/F-6 **entfallen** (kein Befund).
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure. (a) Anker: der Lerneintrag
  verkörpert nichts neu, er schärft die Anwendung einer bestehenden Regel
  ([`AGENTS.md`](../../../../AGENTS.md) §3.13, am Ort existent); (b) Folge-Slice: keiner neu; (c) Register:
  beide genannten Kennungen existieren als Verzeichnis, und ihre
  `evidence/`-Verzeichnisse tragen die neue Datei (3× bzw. 27×, real
  ausgezählt).

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
