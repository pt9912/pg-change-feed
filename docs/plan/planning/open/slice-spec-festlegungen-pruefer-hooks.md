# Slice spec-festlegungen-pruefer-hooks: Festlegungen der Prüfer und Hooks ohne Gate im Pflichtenheft

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
(Entscheidung 7: Baseline-Aktualisierung und Pin-Inventar). Kette je Werkzeug:
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) für
`make suchlauf-nachmessen` und `make kommentar-kennungen`;
[`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
für `make pin-stale-all`, `ADR-0051` Entscheidung 7 für die übrigen
`make pin-stale-*`; [`MR-003`](../../../../harness/conventions/MR-003-guard-inplace-textwerkzeug.md)
und [`MR-004`](../../../../harness/conventions/MR-004-guard-host-python-am-kopf.md)
für den PreToolUse-Guard (`make test-command-guard`);
[`ADR-0062`](../../adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
und [`ADR-0069`](../../adr/0069-commit-msg-hook-einseitige-zusage.md) für den
`commit-msg`-Hook; `make fmt-check` hat keine eigene ADR (Vertrag
`harness/sensors/fmt-check.md`).

**Berührte Spec-Stellen:** `spec/pflichtenheft.md` §7 „Festlegungen der
Harness-Werkzeuge“ (neue Zeilen; die Kennungen vergibt der Implementer).

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag, Neuschnitt von
`slice-spec-festlegungen-harness-werkzeuge`). **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung**; die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Herkunft.** Neuschnitt von `slice-spec-festlegungen-harness-werkzeuge` am
2026-10-07 (dort §1 „Planänderung“): der Pilot legt §7 „Festlegungen der
Harness-Werkzeuge“ im Pflichtenheft an; die Festlegungen der übrigen Werkzeuge
gehen je Gruppe an einen eigenen Slice. **Übergabe an diesen Slice:** die
Delta-Punkte R2, R7, R9, T6 und T8 des Bumps auf Baseline v6.16.0
(`slice-harness-baseline-v6-16-0`, Abschnitt „Bump-Ablauf — Belege“) für
`make suchlauf-nachmessen`, `make kommentar-kennungen`, `make fmt-check`,
`make pin-stale-all` und die übrigen `make pin-stale-*`, den PreToolUse-Guard
und den `commit-msg`-Hook. Die DB-Adapter-Coverage war im Auftrag dieser Gruppe
zugeordnet und geht an `slice-spec-festlegungen-coverage-gates` (Grund dort,
§1).

**Übergabe aus `slice-spec-festlegungen-harness-werkzeuge`:** die Festlegung
zum `formnorm`-`cmp` am Form-Commit eines Baseline-Bumps (je MR-Datei der
Vergleich nach Form-Normalisierung, eine Zeile je Datei, jede `cmp 0`; ohne
Zeile ist der Commit falsch bestimmt) —
[`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
Entscheidung 4, zweiter Spiegelstrich, übergeben durch
[`ADR-0162`](../../adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
Entscheidung 4. Dieser Slice führt den Vertrag `harness/targets/pin-stale.md`
samt Bump-Ablauf und trägt die Festlegung als eigene Zeile in §7; bis dahin
bleibt die Befehlsform `formnorm` in der ADR.

**Ziel:** `spec/pflichtenheft.md` §7 trägt je Werkzeug der Gruppe eine
Festlegung mit eigener `SPEC-<NNN>` (was als Treffer, als Kandidat oder als
Drift gilt, Randformen, Ausgänge); die Verträge unter `harness/sensors/` und
`harness/targets/` nennen die Kennung, statt Schwelle und Randform selbst zu
tragen, und die Zeilen in `harness/README.md` tragen die Bindung
„Spec-Kennung“.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Struktur von §7, die `Schärft:`-Konvention im ADR-Index, die
  `MR-001`-Ablösung und der Leer-Test der Teil-Range** — Gegenstand von
  `slice-spec-festlegungen-harness-werkzeuge` (Pilot); dieser Slice fügt
  Zeilen in eine bestehende Tabelle ein.
- **Gates anderer Gruppen** — Folge-Slices:
  `slice-spec-festlegungen-doku-gates`, `slice-spec-festlegungen-kennungs-gates`,
  `slice-spec-festlegungen-code-gates`, `slice-spec-festlegungen-coverage-gates`.
- **Die Adaptionen `MR-003` und `MR-004` ändern** — Bestand bleibt: die
  Festlegung beschreibt, was der Guard liest; die Adaption bleibt der Grund,
  warum er es darf.
- **`Accepted`-ADRs inhaltlich ändern** — Bestand bleibt
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); die `Schärft:`-Kante entsteht
  über eine Architect-ADR (§2, §6).
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind `spec/pflichtenheft.md`,
  `harness/sensors/`, `harness/targets/`, `harness/README.md` und eine ADR.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**. Alle Beleg-Angaben sind **Zusagen**
(„zu belegen durch …“).

- [ ] **Festlegungen in §7 (Liefer-Punkt 1).** `spec/pflichtenheft.md` §7 trägt
      je eine Zeile mit eigener `SPEC-<NNN>` für `make suchlauf-nachmessen`,
      `make kommentar-kennungen`, `make fmt-check`, `make pin-stale-all` samt
      den übrigen `make pin-stale-*`, den PreToolUse-Guard und den
      `commit-msg`-Hook, in der Form, die der Pilot für §7 festlegt (R9).
      *Zu belegen durch:* Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13)
      und `make docs-check` Exit 0.
- [ ] **Verträge und Index verweisen (Liefer-Punkt 2).** Die Verträge
      `harness/sensors/suchlauf-nachmessen.md`, `kommentar-kennungen.md`,
      `fmt-check.md`, `pin-stale-all.md`, `harness/targets/pin-stale.md` und
      `command-guard.md` nennen die Kennung ihrer Festlegung und tragen Schwelle
      und Randform nicht mehr selbst (R2, R7, T8); die Zeilen in
      `harness/README.md` tragen die Bindung „Spec-Kennung“ (T6).
- [ ] **Bedingung vor der Closure — `Schärft:`-Kante.** Eine Architect-ADR
      (Rollenwechsel, Baseline-Regelwerk `modul-08-agentenrollen.md`) stellt die
      Kante der ADRs dieser Gruppe zu den neuen Kennungen her, nach dem Muster,
      das die ADR des Pilots setzt; für die Werkzeuge, die eine `MR` statt einer
      ADR tragen, entscheidet der Architect die Form. Status `Accepted`, Index
      nachgezogen. Kein Liefer-Punkt: der Implementer liefert die Kennungen, der
      Architect die Kante.
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
| `spec/pflichtenheft.md` §7 | update | eine Zeile je Werkzeug (Liefer-Punkt 1) |
| `harness/sensors/suchlauf-nachmessen.md`, `kommentar-kennungen.md`, `fmt-check.md`, `pin-stale-all.md` | update | Verweis statt Schwelle und Randform (Liefer-Punkt 2) |
| `harness/targets/pin-stale.md`, `command-guard.md` | update | dito (Liefer-Punkt 2) |
| `harness/README.md` §Sensors samt Werkzeug-Tabelle | update | Bindung „Spec-Kennung“ (Liefer-Punkt 2) |
| `docs/plan/adr/<NNNN>-…` und ADR-Index | neu (Architect) | `Schärft:`-Kante (Bedingung vor der Closure) |

Der Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13) ist beim Start zu
messen; bewegte Eigenschaft ist der Ort von Schwelle und Randform der
Werkzeuge (Träger außerhalb der Verträge: `AGENTS.md` §3.1 „Durchsetzung“,
§3.7, §3.13 „Nachmessen“, die Traceability-Regeln in `harness/README.md`).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): `slice-spec-festlegungen-harness-werkzeuge`
liegt in `done/` — §7 und das Muster der `Schärft:`-Kante existieren; danach
Priorisierung durch den Auftraggeber, `in-progress/` trägt keinen Slice
(WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die sechs Werkzeuge sprengen eine
  Review-Sitzung — dann die Prüfer (`suchlauf-nachmessen`,
  `kommentar-kennungen`, `fmt-check`) und die Pins mit den Hooks
  (`pin-stale-*`, Guard, `commit-msg`) als zwei Slices.
- `in-progress` → `open` (blockiert): der Architect entscheidet, dass die Kante
  für die MR-getragenen Werkzeuge eine eigene Form braucht, und die ADR liegt
  nicht vor.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

Die zwei Liefer-Punkte sind abgehakt mit Beleg, die Architect-ADR ist
`Accepted`, `make gates` endet mit Exit 0, der Review-Report liegt vor und ist
aufgelöst, die Closure-Notiz trägt den Lerneintrag und jedes Risiko aus §6
seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Architect-Frage zur `Schärft:`-Kante.** Die Bestands-ADRs dieser Gruppe
  sind `Accepted` und tragen ihr `Schärft:`-Feld unveränderlich; der Guard und
  `make fmt-check` haben keine ADR. Ob eine Sammel-ADR die Kante herstellt und
  wie ein MR-getragenes Werkzeug zu seiner Festlegung kommt, entscheidet der
  Architect. *Zu belegen durch:* die ADR.
- **Die Größe** — sechs Werkzeuge, Verträge mit zusammen rund 700 Zeilen
  (abgeleitet aus `wc -l` bei der Anlage: 112, 194, 133, 134, 130 und 9).
  *Zu belegen durch:* die Rückführungs-Bedingung in §4, beim Start gegen den
  Diff-Umfang geprüft.
- **Die Spec darf nicht auf ADRs zeigen** (gemessen im Pilot: `matrix-forbidden`
  für Link und Kennung im Inline-Code). Die Festlegung trägt ihren Inhalt
  selbst. *Zu belegen durch:* `make docs-check` Exit 0.

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
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `harness/` und
der ADR-Index; die Modus-Deklaration führt nur die Default-Sub-Area `*`
(Kürzel `PGC`, Greenfield), alle Pfade fallen unter sie.

**Vorgelagert — offene Beobachtungen sichten:** beim Übergang `open → next`
nachzuholen (gemergter Stand des Registers); bei Anlage gelesen:
`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` (2×) trifft
Liefer-Punkt 1; `BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar` (2×)
betrifft die Festlegung der `pin-stale`-Achsen;
`BEO-PGC/werkzeugvertrag-zusage-ohne-testfall` (1×) betrifft jede Festlegung,
deren Randform kein Fall des Tabellentests trägt.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
