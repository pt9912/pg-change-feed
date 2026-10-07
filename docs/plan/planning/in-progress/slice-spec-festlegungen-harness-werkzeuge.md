# Slice spec-festlegungen-harness-werkzeuge: Festlegungen der Harness-Werkzeuge in das Pflichtenheft

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
(Entscheidung 7: Baseline-Aktualisierung als bewusster Bootstrap-Vorgang).
Quellen der zwei Festlegungen:
[`ADR-0158`](../../adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md),
[`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md),
[`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
(Entscheidungen 4 und 5) für `make zitat-vergleich`;
[`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
(Entscheidung 1) für den Leer-Test der Teil-Range. Alle vier tragen
`Schärft: —` und sind `Accepted`; die Kante zu den neuen Kennungen stellt eine
Architect-ADR her (§2). Die Adaption
[`MR-001`](../../../../harness/conventions.md#mr-001)
(Rang-2-Datei heißt Pflichtenheft) wird durch `MR-006` abgelöst. Keine
`LH-*`-Anforderung ist berührt: die Festlegungen gelten Harness-Werkzeugen,
nicht dem Produkt.

**Berührte Spec-Stellen:** `spec/pflichtenheft.md` §1 (Verweis „§2 bis §6“
wird „§2 bis §7“), ein neuer §7 „Festlegungen der Harness-Werkzeuge“ mit zwei
`SPEC-<NNN>` und §7 Historie, die §8 wird. `spec/lastenheft.md` ist nicht
berührt (§6, gemessen).

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag, Closure von
`slice-harness-baseline-v6-16-0`). **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung**; die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Herkunft.** Der Bump auf Baseline v6.16.0 (`slice-harness-baseline-v6-16-0`,
Abschnitt „Bump-Ablauf — Belege“) gab den Delta-Punkten R2, R7, R8, R9, T3, T4,
T6, T7, T8 und T10 den Ausgang „Folge-Slice Neu 2“: was ein Harness-Werkzeug
prüft (Schwelle, Randform, Spec-Kennung), steht in der Spezifikation, neuer
Abschnitt der Spezifikations-Vorlage vor der Historie
(`regelwerk/grundlagen-referenz-richtung.md` §Spec-Straten,
`regelwerk/modul-03-spec.md`, `templates/harness/sensors/gate.template.md`,
`templates/docs/plan/adr/NNNN-titel.template.md`).

**Planänderung (2026-10-07) — Neuschnitt in Arbeit, keine Rückführung.**
*Grund:* Die Rückführungs-Bedingung „zu groß“ aus §4 trat beim Start ein; der
Auftraggeber hat statt der Rückführung `in-progress → next` den Zuschnitt A
freigegeben: dieser Slice bleibt in `in-progress/` und wird zum Pilot, die
Festlegungen der übrigen Werkzeuge gehen an fünf Folge-Slices in `open/`.
*Messung des Implementers* (Lauf am Stand `abbe11b4`, Probe als Kopie des
Arbeitsbaums im Scratchpad, nicht committet):

- Die Spec darf nicht auf ADRs zeigen: ein §7 mit Link auf
  [`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
  und mit `ADR-0158` im Inline-Code ergibt in `make docs-check` zwei Befunde
  `matrix-forbidden` (gedruckte Zeile `d-check: 1814 Datei(en) geprüft,
  2 Befund(e)`). Die Festlegung trägt ihren Inhalt deshalb selbst, und die
  Kante läuft aufwärts über `Schärft:` einer ADR.
- Die vier Quell-ADRs tragen `Schärft: —` und sind `Accepted`
  (`grep -n 'Schärft' docs/plan/adr/015[89]*.md docs/plan/adr/016[01]*.md`):
  die Kante braucht eine neue ADR.
- `spec/lastenheft.md` ist nicht betroffen: die einzige Nennung von
  „Abschnitte 1–7“ steht in Zeile 1517 in der Historie-Tabelle und beschreibt
  die Gliederung des Lastenhefts selbst.
- Die nächste freie Kennung im Pflichtenheft ist `SPEC-038` (höchste vergebene
  `SPEC-037`, `grep -o 'SPEC-[0-9]\{3\}' spec/pflichtenheft.md | sort -u`); der
  Implementer misst sie beim Schreiben erneut.
- Größe der Verträge (`wc -l harness/sensors/*.md`, bei der Planänderung
  nachgemessen): `coverage-gate.md` 385 Zeilen mit acht ADRs,
  `db-adapter-coverage.md` 313, `docs-check.md` 247,
  `harness/targets/zitat-vergleich.md` 285 — alle Gates in einem Slice sprengen
  eine Review-Sitzung.

**Ziel (Zuschnitt A):** `spec/pflichtenheft.md` trägt §7 „Festlegungen der
Harness-Werkzeuge“ in der Struktur der Vorlage
(`.harness/baseline/v6.16.0/templates/spec/spezifikation.template.md` §7) mit
zwei Festlegungen — `make zitat-vergleich` und der Leer-Test der Teil-Range
`teilrange` —, die Historie ist §8, und §1 verweist auf „§2 bis §7“; der
Vertrag von `make zitat-vergleich` und seine Zeile in `harness/README.md`
verweisen auf die Kennung; die `Schärft:`-Konvention im ADR-Index zeigt auf
§7; eine Architect-ADR stellt die `Schärft:`-Kante der Quell-ADRs her; `MR-006`
löst `MR-001` ab.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Festlegungen der übrigen Gates, Prüfer und Hooks** — Folge-Slices (je in
  `open/`, Start-Trigger: dieser Slice in `done/`):
  `slice-spec-festlegungen-doku-gates` (docs-check samt aller Module,
  commit-traceability, baseline-verify),
  `slice-spec-festlegungen-kennungs-gates` (sdk-public-doc-check,
  handbuch-public-doc-check, ausgabe-kennungen-check, meldungscodes-check),
  `slice-spec-festlegungen-code-gates` (a-check, generated-sync),
  `slice-spec-festlegungen-coverage-gates` (coverage-gate,
  DB-Adapter-Coverage) und `slice-spec-festlegungen-pruefer-hooks`
  (suchlauf-nachmessen, kommentar-kennungen, fmt-check, pin-stale-*,
  PreToolUse-Guard, commit-msg-Hook). Jeder trägt die Delta-Punkte R2, R7, R9,
  T6 und T8 für seine Werkzeuge in §1 als Übergabe.
- **Werkzeug-eigene Teile des Gate-Index (Neu 1)** — Gegenstand von
  `slice-harness-gate-index-werkzeug-teile`: dort der Ort des Gate-Index, hier
  der Ort der Festlegung.
- **`Accepted`-ADRs inhaltlich ändern** — Bestand bleibt: ein `Schärft:`-Feld
  an einer `Accepted`-ADR nachzutragen wäre ein Überschreiben
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); die Kante entsteht über die
  Architect-ADR (§2).
- **Die Träger des Leer-Tests außerhalb des ADR-Index umschreiben**
  (`AGENTS.md` §3.5, `.claude/agents/implementer.md`, `verifier.md`) — Bestand
  bleibt: sie sind Briefings, keine Werkzeug-Verträge (R2, R7 und T8 gelten
  Verträgen), zitieren die Funktion `teilrange` und bleiben nach der Festlegung
  gültig.
- **`spec/lastenheft.md` ändern** — Bestand bleibt: die Nennung von
  „Abschnitte 1–7“ dort betrifft das Lastenheft selbst (gemessen, oben).
- **`formnorm`-`cmp` am Form-Commit** — übernimmt
  `slice-spec-festlegungen-pruefer-hooks` (führt `harness/targets/pin-stale.md`),
  [`ADR-0162`](../../adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
  Entscheidung 4.
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind `spec/pflichtenheft.md`,
  `harness/targets/zitat-vergleich.md`, `harness/README.md`, der ADR-Index, eine
  neue ADR und `harness/conventions/`.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**. Alle Beleg-Angaben sind **Zusagen**
(„zu belegen durch …“).

- [x] **§7 im Pflichtenheft mit zwei Festlegungen (Liefer-Punkt 1).**
      `spec/pflichtenheft.md` trägt §7 „Festlegungen der Harness-Werkzeuge“ in
      der Struktur der Vorlage (Regel-Absatz, Tabelle `ID · Werkzeug ·
      Festlegung`), die Historie ist §8 mit einer neuen Zeile, §1 verweist auf
      „§2 bis §7“ (T10, R9 für die zwei Werkzeuge). Zwei Zeilen mit eigener
      `SPEC-<NNN>` (erwartet `SPEC-038` und `SPEC-039`, der Implementer misst):
      (a) `make zitat-vergleich` — Vergleichseinheit je Verweisform,
      Normalisierung, die zwei Stände, der nicht messbare Referent und die
      Ausgänge; (b) der Leer-Test der Teil-Range `teilrange` am Pin-Commit.
      Die Zeilen tragen ihren Inhalt selbst und nennen keine ADR (Messung in
      §1). **Dazu, als Bedingung vor der Closure:** eine Architect-ADR
      (Rollenwechsel, Baseline-Regelwerk `modul-08-agentenrollen.md`) mit
      `Schärft:` auf die zwei Kennungen stellt die Kante von
      [`ADR-0158`](../../adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
      bis [`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
      her; Status `Accepted`, ADR-Index nachgezogen. Sie setzt das Muster, dem
      die Folge-Slices folgen. *Zu belegen durch:* Suchlauf (§3),
      `make docs-check` Exit 0 und die ADR. *Stand Implementer:* §7 mit
      `SPEC-038` und `SPEC-039` geliefert (`14bbd8dc`, Nachtrag `8e89831d`),
      Suchlauf und `make docs-check` im Befund unter §3; die Architect-ADR ist
      [`ADR-0162`](../../adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
      (`Accepted`, `f2d6f194`), ihre Wortlaut-Berichtigung an `SPEC-038` und
      `SPEC-039` steht in `d2a863a1`.
- [x] **Vertrag, Index und Konvention verweisen (Liefer-Punkt 2).**
      `harness/targets/zitat-vergleich.md` nennt die Kennung (a) und trägt
      Einheit, Normalisierung und Ausgänge nicht mehr selbst (R2, R7, T8 für
      dieses Werkzeug); seine Zeile in `harness/README.md` trägt die Bindung
      „Spec-Kennung“ (T6 für diese Zeile); der Kommentar-Block in
      `harness/README.md` §Sensors trägt den Satz, dass was ein Werkzeug prüft
      in der Spezifikation steht (T7); die `Schärft:`-Konvention im ADR-Index
      nennt §7 als Ziel einer Gate-ADR (T4).
- [x] **`MR-006` löst `MR-001` ab (Liefer-Punkt 3).** `harness/conventions/MR-006-…`
      (per `cp` aus `.harness/baseline/v6.16.0/templates/harness/conventions/MR-NNN-titel.template.md`)
      übernimmt die Adaption „Technik-Dokument heißt Pflichtenheft“ mit der
      Struktur §1–§8 und den Feldern `Löst auf` und `Ausgelöst durch
      Baseline-Stand`; `MR-001` geht per `git mv` nach
      `harness/conventions/done/` (eigener Commit); der Index in
      `harness/conventions.md` führt `MR-006` unter *Aktive* und `MR-001` unter
      *Aufgelöste Adaptionen* (R8). Freigabe des Auftraggebers liegt vor
      (2026-10-07).
- [ ] `make gates` grün, Exit-Code ungefiltert gesichert
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/pflichtenheft.md` | update | §1-Verweis, neuer §7 mit zwei Zeilen, Historie §8 (Liefer-Punkt 1) |
| `docs/plan/adr/<NNNN>-…` und ADR-Index | neu (Architect) | `Schärft:`-Kante von `ADR-0158` bis `ADR-0161` (Bedingung in Liefer-Punkt 1) |
| `harness/targets/zitat-vergleich.md` | update | Verweis auf Kennung (a) (Liefer-Punkt 2) |
| `harness/README.md` (Zeile `make zitat-vergleich`, Kommentar-Block §Sensors) | update | T6, T7 (Liefer-Punkt 2) |
| `docs/plan/adr/README.md` §Konventionen | update | `Schärft:`-Konvention (T4, Liefer-Punkt 2) |
| `harness/conventions/MR-006-…`, `MR-001-…` → `done/`, `harness/conventions.md` | neu, `git mv`, update | R8 (Liefer-Punkt 3) |
| `docs/plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/state.md` | Prüfauftrag | nennt „`spec/pflichtenheft.md` §7 Historie“ (Suchlauf unten); Zustandsfeld, nachziehen oder mit Grund stehen lassen — **nachgezogen** auf §8 (Implementer, `14bbd8dc`) |
| `harness/targets/pin-stale.md` (Bump-Ablauf Schritt 1, Absatz „Form-Commit vor dem Löschen“) | update (Plan-Nachzug, Implementer) | der Vertrag wiederholte die Lesung des Adaptions-Durchgangs und die drei Fälle des Leer-Tests; er verweist jetzt auf `SPEC-038` bzw. `SPEC-039` (R2, R7, T8 für diese zwei Werkzeuge; Befehlsform `teilrange` bleibt mit ADR-Zeiger) |
| `harness/README.md` Kommentar-Block §Sensors, Liste „Bindung“ | update (Plan-Nachzug, Implementer) | „Spec-Kennung“ in der Bindungs-Liste wie in der Vorlage (T6, im Kommentar statt nur in der Zeile) |
| `.claude/agents/architect.md`, `planner.md`, `reviewer.md`, `.harness/skills/reviewer.md` | update (Plan-Nachzug, Implementer) | nannten `MR-001` als geltende Adaption; jetzt `MR-006` (LP3, Suchlauf `git grep -n MR-001`) |
| `spec/pflichtenheft.md` §7 | update (Plan-Nachzug, Implementer) | vier Randformen, die nur im Vertrag standen (Inline-Code mit Text vor dem Tag, Fence in einem Listenpunkt, Locale, leeres Tag-Paar), in `SPEC-038` übernommen, bevor der Vertrag sie abgab (`8e89831d`) |
| `formnorm`-`cmp` am Form-Commit | **nicht realisiert**, mit Grund | ist weder `make zitat-vergleich` noch der Leer-Test; `SPEC-038` nennt ihn nur als Beleg der Form-Korrektur, die Festlegung übernimmt `slice-spec-festlegungen-pruefer-hooks` (`ADR-0162` Entscheidung 4, Ausschluss in §1) |
| `docs/plan/adr/0162-…` und Index | neu (Architect) | erledigt |
| Umzugs-Commit `8e00e831` (`MR-001` nach `harness/conventions/done/`) | `git mv`, eigener Commit (Implementer) | reiner Umzug nach [`ADR-0162`](../../adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 3, belegt mit `umzug` (Beleg unten); die Message nennt `MR-001` und `ADR-0161`, nicht `ADR-0162` — der Commit liegt vor der ADR (`f2d6f194`) |
| `.claude/agents/verifier.md`, `.claude/agents/implementer.md`, `harness/targets/pin-stale.md` (Absatz „Umzug eines abgelösten Eintrags“), `harness/conventions.md` §Adaptions-Block | update (Implementer, wörtliche Vorgabe des Architects nach `ADR-0162`, `d2a863a1`) | Teilung und Beleg am Umzugs-Commit; ändert den Ausschluss „Träger des Leer-Tests außerhalb des ADR-Index umschreiben“ in §1 für die zwei Briefings — Planänderung auf Vorgabe des Architects |
| Abschnitt „Test“ in `harness/targets/zitat-vergleich.md` | **nicht geändert**, mit Grund | die Deckung des Werkzeugs gehört nach `gate.template.md` in ADR oder Skriptkopf, Liefer-Punkt 2 verlangt nur, dass Einheit, Normalisierung und Ausgänge den Vertrag verlassen; der Umzug der Testbeschreibung ist ein eigener Vorgang |

**Suchlauf — bewegte Eigenschaften:** die Nummer der Historie (§7 → §8), der
Abschnitts-Umfang „§2 bis §6“ / „Abschnitte 1–7“ und das Ziel der
`Schärft:`-Konvention. Suchraum: der ganze Baum ohne `docs/reviews/**`,
`docs/plan/planning/done/**` und `.harness/baseline/**`. Gemessen beim
Neuschnitt am Parent `abbe11b4`; die `diff`-Zeilen gelten dem Arbeitsbaum
nach dem Commit der Folge-Slices. Der Implementer ergänzt die `diff`-Zahlen
nach seiner Arbeit und nennt Gefundenes und Nichtgefundenes.

```suchlauf
abbe11b4 3 -F 'Abschnitte 1–7' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 4 -F 'Abschnitte 1–7' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
abbe11b4 1 -F '§2 bis §6' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 0 -F '§2 bis §6' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
abbe11b4 1 -E 'pflichtenheft[^ ]*`? §7' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
abbe11b4 0 -F '7-historie' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
abbe11b4 1 -F 'welche Spec-Stelle' -- docs/plan/adr/README.md
diff 2 -F '§2 bis §7' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 22 -E 'pflichtenheft[^ ]*`? §7' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 1 -E 'pflichtenheft[^ ]*`? §8' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 0 -F '7-historie' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 14 -F '7-festlegungen-der-harness-werkzeuge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 1 -F 'Spec-Stelle in' -- docs/plan/adr/README.md
```

Befund am Parent: „Abschnitte 1–7“ in `spec/lastenheft.md` (Historie, nicht
berührt), `MR-001` (abgelöst durch `MR-006`) und
[`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
Zeile 187 (`Accepted`, bleibt); „§2 bis §6“ nur in `spec/pflichtenheft.md` §1;
die Historie als §7 nur im Zustandsfeld oben; kein Anker `#7-historie`.

Befund im Arbeitsbaum nach `dc04087e` (Implementer, `make suchlauf-nachmessen`
unten): „Abschnitte 1–7“ 4 — die drei vom Parent (`spec/lastenheft.md`
bleibt, `MR-001` jetzt unter `harness/conventions/done/`, `ADR-0161` bleibt)
und das Zitat im Feld `Löst auf` von `MR-006`; „§2 bis §6“ 0, „§2 bis §7“ 2
(`spec/pflichtenheft.md` §1 und die neue Historie-Zeile); „pflichtenheft … §7“
21 — einmal der neue Verweis im ADR-Index, zwanzigmal die fünf Folge-Slices in
`open/`, alle im neuen Sinn (§7 = Festlegungen); „pflichtenheft … §8“ 1 (das
nachgezogene Zustandsfeld); der Anker `#7-festlegungen-der-harness-werkzeuge`
11 (ADR-Index 1, `harness/README.md` 1, `pin-stale.md` 2,
`zitat-vergleich.md` 7); kein `#7-historie`. Nach `d2a863a1`: „pflichtenheft …
§7“ 22 (dazu `ADR-0162` Zeile 31), der Anker 14 (dazu `ADR-0162` 2 im
`Schärft:`-Feld und der neue Absatz in `pin-stale.md`); die übrigen Zeilen
unverändert. Nicht gefunden: ein Verweis auf
die alte §7 als Historie außerhalb des Zustandsfelds; in den `Accepted`-ADRs
`ADR-0090` und `ADR-0098` steht „§7“ für die Historie (Zeilen 407, 562 bzw.
449), sie bleiben eingefroren (`AGENTS.md` §3.5). `MR-001` als geltende
Adaption nannten vier Briefings (`.claude/agents/architect.md`,
`planner.md`, `reviewer.md`, `.harness/skills/reviewer.md`), nachgezogen auf
`MR-006`; `MR-005` nennt `MR-001` als Anlass und bleibt (immutabel ab Datum).
`make docs-check` am Arbeitsbaum: Exit 0, `1820 Datei(en) geprüft, 0 Befund(e)`.

**Beleg zur Verifier-Range (Implementer, am Stand `d2a863a1`).** Die Range des
Slice `fcff30ec..HEAD` enthält den Umzugs-Commit `8e00e831`; ein Lauf über die
ganze Range endet mit `core-drift-vcs` an `MR-001` (Exit 2, gemessen vor
`ADR-0162`). Geteilt nach `SPEC-039`, mit `teilrange` aus
[`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
Entscheidung 1 und `umzug` aus
[`ADR-0162`](../../adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
Entscheidung 3, beide wörtlich aus dem `bash`-Block der ADR gezogen
(`awk`-Ausschnitt in eine Datei im Scratchpad, `source`):

```text
teilrange fcff30ec 8e00e831~1; echo "Exit $?"
teilrange: fcff30ec..8e00e831~1 enthält 2 Commit(s), Lauf
d-check: 1821 Datei(en) geprüft, 0 Befund(e)
Exit 0
teilrange 8e00e831 HEAD; echo "Exit $?"
teilrange: 8e00e831..HEAD enthält 4 Commit(s), Lauf
d-check: 1821 Datei(en) geprüft, 0 Befund(e)
Exit 0
umzug 8e00e831; echo "Exit $?"
umzug: R100	harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md	harness/conventions/done/MR-001-technik-dokument-heisst-pflichtenheft.md
umzug: 8e00e831 ist ein reiner Umzug nach harness/conventions/done/, Exit 0
Exit 0
```

Die Commits nach `d2a863a1` ändern keinen Adaptions-Eintrag; der Verifier fährt
die drei Aufrufe am Endstand nach.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): Priorisierung durch den
Auftraggeber; `in-progress/` trägt keinen Slice (WIP-Limit 1). Der Slice hängt
nicht an `slice-harness-gate-index-werkzeug-teile`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die Festlegungen aller Gates sprengen eine
  Review-Sitzung — dann je Gate-Gruppe ein Slice. **Eingetreten am
  2026-10-07**; Ausgang ist die Planänderung in §1 (Zuschnitt A, fünf
  Folge-Slices), freigegeben durch den Auftraggeber, keine Rückführung.
- `in-progress` → `next` (Zuschnitt A zu groß): die Festlegung von
  `make zitat-vergleich` allein sprengt eine Review-Sitzung — dann
  `MR-006` als eigener Slice.
- `in-progress` → `open` (blockiert): der Architect lehnt eine ADR ab, die
  `Schärft:` auf eine `SPEC-<NNN>` in §7 setzt, und benennt keinen anderen
  Weg der Kante.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

Die drei Liefer-Punkte sind abgehakt mit Beleg, die Architect-ADR ist
`Accepted`, `make gates` endet mit Exit 0,
der Review-Report liegt vor und ist aufgelöst, die Closure-Notiz trägt den
Lerneintrag und jedes Risiko aus §6 seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Berührung des Lastenhefts (Rang 1).** Der Bump-Plan zählte die Nennung von
  „Abschnitte 1–7“ in `spec/lastenheft.md` als betroffen. **Gemessen** bei der
  Anlage und vom Implementer bestätigt (`grep -n 'Abschnitte 1–7'
  spec/lastenheft.md`): eine Zeile, 1517, in der Historie-Tabelle (Version
  0.3.0) — sie beschreibt die Gliederung des Lastenhefts selbst. Erwarteter
  Ausgang: *entfallen*, mit diesem Grund.
- **Umnummerierung bricht Verweise auf §7 des Pflichtenhefts** (Anker,
  Prosa-Nennungen). *Zu belegen durch:* Suchlauf (§3) an beiden Ständen und
  `make docs-check`.
- **Kante von `ADR-0158` bis `ADR-0161` zu den neuen Kennungen.** Die vier
  tragen `Schärft: —` und sind immutabel; die Spec darf nicht auf ADRs zeigen
  (Messung in §1). Die Kante entsteht nur über eine Architect-ADR; lehnt der
  Architect sie ab, greift die Rückführung `in-progress → open` (§4).
  *Zu belegen durch:* die ADR.
- **Lücke in §7 bis zu den Folge-Slices.** Nach diesem Slice trägt §7 zwei
  Festlegungen; die übrigen Gates, Prüfer und Hooks führen Schwelle und
  Randform weiter in ihren Verträgen, bis
  `slice-spec-festlegungen-doku-gates`, `-kennungs-gates`, `-code-gates`,
  `-coverage-gates` und `-pruefer-hooks` geschlossen sind. Die Regel der
  Baseline („steht in der Spezifikation“) gilt bis dahin nur für zwei
  Werkzeuge — benannt, nicht still. Erwarteter Ausgang: *eingetreten* mit den
  fünf Kennungen der Folge-Slices.
- **Festlegung lässt eine Randform offen** (`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`,
  2×; ein drittes Auftreten wäre eine Lücke). *Zu belegen durch:* Gegenlesen
  der zwei Zeilen gegen den Vertrag von `make zitat-vergleich` und die vier
  Quell-ADRs; die Folge-Slices brauchen die Form von §7 ohne Rückfrage.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Gegenstand:** <übernommen von `slice-<Kennung>` | entfallen: <Grund>>
  *(nur beim Ausgang ohne Arbeit; sonst Zeile löschen)*
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<KUERZEL>/<slug>` (<slice-kennung-a>, <slice-kennung-b>, <slice-kennung-c> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-<Kennung>.md` | `evidence/slice-<Kennung>.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-<Kennung> (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Der Abschnitt selbst entfällt nie.**

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `harness/`
(samt `harness/conventions/`), der ADR-Index und eine neue ADR; die Modus-Deklaration führt nur die Default-Sub-Area `*`
(Kürzel `PGC`, Greenfield), alle Pfade fallen unter sie.

**Vorgelagert — offene Beobachtungen sichten:** beim Neuschnitt (2026-10-07)
am Stand `abbe11b4` gelesen:
`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` (2×) trifft
Liefer-Punkt 1 (Risiko in §6);
`BEO-PGC/messwerkzeug-grenze-unbenannt-fail-open` (1×) betrifft
`make zitat-vergleich`, das Markdown nachbildet — die Festlegung nennt die
Grenze; `BEO-PGC/werkzeugvertrag-zusage-ohne-testfall` (1×) betrifft die
Randformen der Festlegung (a), die der Tabellentest `make test-zitat-vergleich`
tragen muss. Keiner erreicht mit diesem Slice 3×, solange die Festlegungen
vollständig sind.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
