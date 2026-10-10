# Slice spec-festlegungen-commit-baseline-gates: Festlegungen von commit-traceability, baseline-verify und doc-immutable im Pflichtenheft

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
ADR-Kette der Gates dieses Slice:
[`ADR-0045`](../../adr/0045-commit-traceability-standing-gate.md) für
`make commit-traceability`;
[`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
(Block `vcs:`) und
[`ADR-0162`](../../adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
(Entscheidung 3, Umzugs-Commit) für `make doc-immutable`;
`make baseline-verify` hat keine eigene ADR (Vertrag
`harness/sensors/baseline-verify.md`). Keine `LH-*`-Anforderung ist berührt:
die Festlegungen gelten Harness-Werkzeugen, nicht dem Produkt.

**Berührte Spec-Stellen:** `spec/pflichtenheft.md` §7 „Festlegungen der
Harness-Werkzeuge“ (neue Zeilen; die Kennungen vergibt der Implementer).

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag, Neuschnitt bei der Rückführung
von `slice-spec-festlegungen-doku-gates`). **Datum:** 2026-10-10.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Herkunft.** Abgetrennt von `slice-spec-festlegungen-doku-gates` bei dessen
Rückführung `in-progress → next` am 2026-10-10 (dort §4, Rückführungs-Grund):
die Quellen der vier Gates zusammen sprengten eine Review-Sitzung. Jener Slice
behält `make docs-check`; dieser übernimmt die drei übrigen Gates samt der
Übergabe-Punkte, die jener Plan für sie trug. Der Text der Übergabe steht in
§2 (Übergabe-Block). Eine `Übernimmt:`-Zeile trägt dieser Slice nicht: der
Geber geht nicht nach `done/`, er arbeitet weiter an `docs-check`.

**Ziel:** `spec/pflichtenheft.md` §7 trägt für `make commit-traceability`,
`make baseline-verify` und `make doc-immutable` je eine Festlegung mit eigener
`SPEC-<NNN>` (was als Treffer gilt, Randform, Ausgänge); die Verträge
`harness/sensors/commit-traceability.md` und
`harness/sensors/baseline-verify.md` nennen die Kennung, statt Schwelle und
Randform selbst zu tragen, und `harness/README.md` §Sensors trägt für alle drei
die Bindung „Spec-Kennung“.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`make docs-check`** (alle acht Module) — bleibt bei
  `slice-spec-festlegungen-doku-gates`; dort §1 und DoD.
- **Der Leer-Test der Teil-Range vor einem `make doc-immutable`-Lauf** —
  Bestand: steht schon in `SPEC-039` (Pilot
  `slice-spec-festlegungen-harness-werkzeuge`); dieser Slice legt den Lauf
  selbst fest, nicht seine Vorbedingung.
- **Die Struktur von §7, die `Schärft:`-Konvention im ADR-Index** —
  Gegenstand des Pilots; dieser Slice fügt Zeilen in eine bestehende Tabelle
  ein.
- **Gates anderer Gruppen** — Folge-Slices:
  `slice-spec-festlegungen-kennungs-gates`,
  `slice-spec-festlegungen-code-gates`,
  `slice-spec-festlegungen-coverage-gates`,
  `slice-spec-festlegungen-pruefer-hooks`.
- **`Accepted`-ADRs inhaltlich ändern** — Bestand bleibt
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); die `Schärft:`-Kante entsteht
  über eine Architect-ADR (§2, §6).
- **Das Verhalten eines Gates ändern** — ein anderer Vorgang: die Festlegung
  beschreibt, was das Werkzeug heute entscheidet; weicht der Vertrag vom
  Werkzeug ab, ist das ein Befund mit eigenem Slice (§6).
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind
  `spec/pflichtenheft.md`, `harness/sensors/`, `harness/README.md` und eine
  ADR.

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
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

**Übergabe-Block** (committeter Text der Übergabe aus
`slice-spec-festlegungen-doku-gates`, Stand `b2451d70`, dort §1 und
Liefer-Punkte 1 und 2):

1. **Delta-Punkte des Bumps auf Baseline v6.16.0** — R2, R7, R9, T6 und T8
   (`slice-harness-baseline-v6-16-0`, Abschnitt „Bump-Ablauf — Belege“) für
   `make commit-traceability` und `make baseline-verify`: R9 die Form der
   §7-Zeile, die der Pilot festlegt; R2, R7 und T8 der Vertrag verweist auf
   die Kennung statt Schwelle und Randform zu tragen; T6 die Bindung
   „Spec-Kennung“ in `harness/README.md` §Sensors.
2. **Review F-5 aus `slice-spec-festlegungen-harness-werkzeuge`** (LOW,
   übergeben in dessen Closure, `0119aeb9`): `make doc-immutable` (d-check
   Modul `vcs`, Block `vcs:` in `.d-check.yml`, Ziel in `d-check.mk`; kein
   Gate, Lauf des Verifiers) bekommt eine Festlegung. Festzulegen ist, was als
   Drift gilt (der Kern einer Datei unter `vcs.paths` ab der Zeile
   `immutable-when` bleibt über die Range unverändert), die zwei Formen des
   Aufrufs (`RANGE=<basis>..<spitze>` und `STAGED=1`) und die Ausgänge.
3. **Review F-13 aus demselben Slice** (INFO, übergeben): `make doc-immutable`
   bekommt eine Zeile im Werkzeug-Teil von `harness/README.md` §Sensors mit
   der Bindung „Spec-Kennung“; heute steht es dort in keiner Zeile. Die Zeile
   entsteht erst mit der Kennung, damit sie Randform und Ausgänge nicht selbst
   trägt.

- [ ] **Festlegungen in §7 (Liefer-Punkt 1).** `spec/pflichtenheft.md` §7 trägt
      je eine Zeile mit eigener `SPEC-<NNN>` für `make commit-traceability`
      (positive Hälfte und Betreff-Grenze, Range), `make baseline-verify`
      (Integrität und Vollständigkeit gegen `SHA256SUMS`) und
      `make doc-immutable` (Drift des Kerns, `RANGE` und `STAGED`, Ausgänge;
      Übergabe 2), in der Form, die der Pilot für §7 festlegt (Übergabe 1, R9).
      *Zu belegen durch:* Gegenprobe der Quellen und Anschluss-Frage (§3),
      Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13) und
      `make docs-check` Exit 0.
- [ ] **Verträge und Index verweisen (Liefer-Punkt 2).**
      `harness/sensors/commit-traceability.md` und
      `harness/sensors/baseline-verify.md` nennen die Kennung ihrer Festlegung
      und tragen Schwelle und Randform nicht mehr selbst (R2, R7, T8); ihre
      zwei Zeilen in `harness/README.md` §Sensors tragen die Bindung
      „Spec-Kennung“ (T6), und `make doc-immutable` bekommt dort eine Zeile im
      Werkzeug-Teil mit derselben Bindung (Übergabe 3).
- [ ] **Bedingung vor der Closure — `Schärft:`-Kante.** Eine Architect-ADR
      (Rollenwechsel, Baseline-Regelwerk `modul-08-agentenrollen.md`) stellt die
      Kante der ADRs aus **Bezug** zu den neuen Kennungen her, nach dem Muster,
      das die ADR des Pilots setzt; Status `Accepted`, Index nachgezogen. Kein
      Liefer-Punkt: der Implementer liefert die Kennungen, der Architect die
      Kante.
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
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/pflichtenheft.md` §7 | update | drei Zeilen: `make commit-traceability`, `make baseline-verify`, `make doc-immutable` (Liefer-Punkt 1) |
| `harness/sensors/commit-traceability.md`, `harness/sensors/baseline-verify.md` | update | Verweis statt Schwelle und Randform (Liefer-Punkt 2) |
| `harness/README.md` §Sensors | update | Bindung „Spec-Kennung“ in zwei Gate-Zeilen, Zeile `make doc-immutable` im Werkzeug-Teil neu (Liefer-Punkt 2) |
| `docs/plan/adr/<NNNN>-…` und ADR-Index | neu (Architect) | `Schärft:`-Kante (Bedingung vor der Closure) |

Kein Test: die Festlegungen beschreiben bestehendes Verhalten, kein Werkzeug
ändert sich (§1).

**Lese-Umfang** (*übernommen* aus der Messung des Implementers von
`slice-spec-festlegungen-doku-gates`, Stand `b2451d70`): die vier Gates
zusammen etwa 1330 Quellzeilen, `docs-check` allein etwa 1130 — für diesen
Slice bleiben also etwa 200 (*abgeleitet*, Differenz).
`harness/sensors/commit-traceability.md` 11 Zeilen,
`harness/sensors/baseline-verify.md` 53 Zeilen. Ob die 200 den Block `vcs:`
und das Ziel in `d-check.mk` schon zählen, ist nicht nachgemessen; beim Start
zu messen.

Der Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13) ist beim Start zu
messen; bewegte Eigenschaft ist der Ort von Schwelle und Randform der zwei
Gates und von `make doc-immutable` (Träger außerhalb der Verträge:
`harness/README.md` §Traceability rules, Kommentare der Blöcke `vcs:` und
`commits:` in `.d-check.yml`, `d-check.mk`, `harness/mk/baseline.mk`).

**Gegenprobe der Quellen und Anschluss-Frage** (Beleg zu Liefer-Punkt 1;
Regel: `.claude/commands/plan-welle.md` Schritt 6 · seit
slice-spec-festlegungen-harness-werkzeuge,
`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`): Vor dem
Review trägt dieser Abschnitt eine Tabelle mit einer Zeile je normativem Satz
der Quellen — die zwei Verträge aus Liefer-Punkt 2, für `make doc-immutable`
der Block `vcs:` in `.d-check.yml` samt Kommentar und das Ziel in
`d-check.mk`, und die Entscheidungen der ADRs aus **Bezug** — und ihrem
Ausgang *steht in `SPEC-<NNN>` · bleibt im Vertrag (Grund) · entfällt
(Grund)*; darunter je Ausgang und Randfall der neuen Zeilen die Folge für den
Anwender (besteht oder nicht, welcher Fall gewinnt). Eine Folge, die keine
Quelle trägt, entscheidet der Architect vor der Closure, nicht eine Lesung im
Auftrag.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): `slice-spec-festlegungen-harness-werkzeuge`
liegt in `done/` (§7 und das Muster der `Schärft:`-Kante existieren); danach
Priorisierung durch den Auftraggeber, `in-progress/` trägt keinen Slice
(WIP-Limit 1). Von `slice-spec-festlegungen-doku-gates` hängt der Start nicht
ab: beide fügen unabhängige Zeilen in §7 ein.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die Gegenprobe der Quellen (§3) für
  `make doc-immutable` zieht die Befehlsformen aus `ADR-0161` und `ADR-0162`
  (Teil-Ranges, Umzugs-Commit, Form-Commit) als eigene Festlegungs-Gegenstände
  nach sich — dann `make doc-immutable` in einen eigenen Slice.
- `in-progress` → `open` (blockiert): der Architect entscheidet, dass die Kante
  für diese Gruppe eine andere Form braucht als im Pilot, und die ADR liegt
  nicht vor.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die zwei Liefer-Punkte sind abgehakt mit Beleg, die Architect-ADR ist
`Accepted`, `make gates` endet mit Exit 0, der Review-Report liegt vor und ist
aufgelöst, die Closure-Notiz trägt den Lerneintrag und jedes Risiko aus §6
seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Architect-Frage zur `Schärft:`-Kante.** Die ADRs aus **Bezug** sind
  `Accepted` und tragen ihr `Schärft:`-Feld unveränderlich; ob eine Sammel-ADR
  oder eine ADR je Gate die Kante herstellt — und ob sie mit der Kante von
  `slice-spec-festlegungen-doku-gates` zusammengeht —, entscheidet der
  Architect. *Zu belegen durch:* die ADR. — **Ausgang:** bei Closure.
- **`make doc-immutable` ist kein Gate und hat keinen eigenen Vertrag** unter
  `harness/sensors/`; Randform und Ausgänge stehen heute verteilt im Block
  `vcs:`, in `d-check.mk` und in Befehlsformen der ADRs aus **Bezug**. Eine
  Festlegung kann einen Satz verlieren, den nur eine dieser Stellen trägt.
  *Zu belegen durch:* die Gegenprobe der Quellen (§3) mit einer Zeile je
  Quelle. — **Ausgang:** bei Closure.
- **Vertrag und Werkzeug weichen ab**, sobald die Randform als Festlegung
  formuliert wird (etwa die Betreff-Grenze von `make commit-traceability`
  gegen `tools/harness/commit-traceability.sh`). *Zu belegen durch:*
  Gegenlesen der Festlegung gegen Vertrag, Skript und ADR; ein Befund wird ein
  eigener Slice. — **Ausgang:** bei Closure.
- **Die Spec darf nicht auf ADRs zeigen** (gemessen im Pilot:
  `matrix-forbidden` für Link und Kennung im Inline-Code). Die Festlegung
  trägt ihren Inhalt selbst. *Zu belegen durch:* `make docs-check` Exit 0. —
  **Ausgang:** bei Closure.

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

**Vorgelagert — offene Beobachtungen sichten** (gemergter Stand `b2451d70`,
`docs/plan/planning/observations/BEO-PGC/`, Zustand `offen` gelesen):

- `BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` (3×,
  verkörpert) trifft Liefer-Punkt 1; Prüfschritt in §3.
- `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (verkörpert) trifft die
  Anlage dieses Slice selbst: geprüft mit `git grep` der Kernbegriffe in
  dieser Datei vor dem Commit (`doc-immutable`, `STAGED`, `immutable-when`,
  `Betreff-Grenze`, `SHA256SUMS`, `F-5`, `F-13`).
- `BEO-PGC/messwerkzeug-grenze-unbenannt-fail-open` (1×) trifft
  `make doc-immutable`: endet es an einer Eingabeform, die es nicht kennt, still
  grün, gehört diese Grenze in die Festlegung oder als Befund in §6.
- `BEO-PGC/git-mv-und-inhalt-in-einem-commit` (2×) berührt den Umzugs-Commit
  aus `ADR-0162` Entscheidung 3, den `make doc-immutable` in Teil-Ranges
  umgeht; ein drittes Auftreten in diesem Slice wäre eine Lücke.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
