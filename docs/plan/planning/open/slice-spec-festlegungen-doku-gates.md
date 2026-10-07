# Slice spec-festlegungen-doku-gates: Festlegungen von docs-check, commit-traceability und baseline-verify im Pflichtenheft

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
(Entscheidung 7: Baseline-Aktualisierung als bewusster Bootstrap-Vorgang). ADR-Kette
der Gates dieser Gruppe:
[`ADR-0072`](../../adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md),
[`ADR-0075`](../../adr/0075-hostpaths-reichweite-und-wortlaut.md),
[`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
(Entscheidung 3, Home-relative Pfade in `hostpaths`) für `make docs-check`;
[`ADR-0045`](../../adr/0045-commit-traceability-standing-gate.md) für
`make commit-traceability`; `make baseline-verify` hat keine eigene ADR (Vertrag
`harness/sensors/baseline-verify.md`). Keine `LH-*`-Anforderung ist berührt: die
Festlegungen gelten Harness-Werkzeugen, nicht dem Produkt.

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
Harness-Werkzeuge“ im Pflichtenheft an und trägt die Festlegungen von
`make zitat-vergleich` und des Leer-Tests der Teil-Range; die Festlegungen der
übrigen Gates gehen je Gruppe an einen eigenen Slice. **Übergabe an diesen
Slice:** die Delta-Punkte R2, R7, R9, T6 und T8 des Bumps auf Baseline v6.16.0
(`slice-harness-baseline-v6-16-0`, Abschnitt „Bump-Ablauf — Belege“) für die
Gates `make docs-check` (alle Module der Liste `modules:` in `.d-check.yml`,
samt `hostpaths`), `make commit-traceability` und `make baseline-verify`.

**Ziel:** `spec/pflichtenheft.md` §7 trägt für `make docs-check`,
`make commit-traceability` und `make baseline-verify` je eine Festlegung mit
eigener `SPEC-<NNN>` (was als Treffer gilt, Randform, Ausgänge); die Verträge
unter `harness/sensors/` nennen die Kennung, statt Schwelle und Randform selbst
zu tragen, und die Zeilen in `harness/README.md` §Sensors tragen die Bindung
„Spec-Kennung“.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Struktur von §7, die `Schärft:`-Konvention im ADR-Index und die
  `MR-001`-Ablösung** — Gegenstand von `slice-spec-festlegungen-harness-werkzeuge`
  (Pilot); dieser Slice fügt Zeilen in eine bestehende Tabelle ein.
- **Gates anderer Gruppen** — Folge-Slices:
  `slice-spec-festlegungen-kennungs-gates` (sdk-, handbuch-public-doc-check,
  ausgabe-kennungen-check, meldungscodes-check),
  `slice-spec-festlegungen-code-gates` (a-check, generated-sync),
  `slice-spec-festlegungen-coverage-gates` (coverage-gate, db-adapter-coverage),
  `slice-spec-festlegungen-pruefer-hooks` (Prüfer und Hooks ohne Gate).
- **`Accepted`-ADRs inhaltlich ändern** — Bestand bleibt
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); die `Schärft:`-Kante entsteht
  über eine Architect-ADR (§2, §6).
- **Das Verhalten eines Gates ändern** — ein anderer Vorgang: die Festlegung
  beschreibt, was das Werkzeug heute entscheidet; weicht der Vertrag vom
  Werkzeug ab, ist das ein Befund mit eigenem Slice (§6).
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind `spec/pflichtenheft.md`,
  `harness/sensors/`, `harness/README.md` und eine ADR.

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
      je eine Zeile mit eigener `SPEC-<NNN>` für `make docs-check` (je Modul
      der Liste `modules:` was als Treffer gilt, samt der Reichweite von
      `hostpaths`), `make commit-traceability` (positive Hälfte und
      Betreff-Grenze, Range) und `make baseline-verify` (Integrität und
      Vollständigkeit gegen `SHA256SUMS`), in der Form, die der Pilot für §7
      festlegt (R9). *Zu belegen durch:* Suchlauf
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13) und `make docs-check` Exit 0.
- [ ] **Verträge und Index verweisen (Liefer-Punkt 2).**
      `harness/sensors/docs-check.md`, `commit-traceability.md` und
      `baseline-verify.md` nennen die Kennung ihrer Festlegung und tragen
      Schwelle und Randform nicht mehr selbst (R2, R7, T8); die drei Zeilen in
      `harness/README.md` §Sensors tragen die Bindung „Spec-Kennung“ (T6).
- [ ] **Bedingung vor der Closure — `Schärft:`-Kante.** Eine Architect-ADR
      (Rollenwechsel, Baseline-Regelwerk `modul-08-agentenrollen.md`) stellt die
      Kante der Gate-ADRs dieser Gruppe zu den neuen Kennungen her, nach dem
      Muster, das die ADR des Pilots setzt; Status `Accepted`, Index
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
| `spec/pflichtenheft.md` §7 | update | drei Zeilen (Liefer-Punkt 1) |
| `harness/sensors/docs-check.md`, `commit-traceability.md`, `baseline-verify.md` | update | Verweis statt Schwelle und Randform (Liefer-Punkt 2) |
| `harness/README.md` §Sensors | update | Bindung „Spec-Kennung“ (Liefer-Punkt 2) |
| `docs/plan/adr/<NNNN>-…` und ADR-Index | neu (Architect) | `Schärft:`-Kante (Bedingung vor der Closure) |

Der Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13) ist beim Start zu
messen; bewegte Eigenschaft ist der Ort von Schwelle und Randform der drei
Gates (Träger außerhalb der Verträge: `AGENTS.md` §3.11, `.d-check.yml`
Kommentare).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): `slice-spec-festlegungen-harness-werkzeuge`
liegt in `done/` — §7 und das Muster der `Schärft:`-Kante existieren; danach
Priorisierung durch den Auftraggeber, `in-progress/` trägt keinen Slice
(WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die Festlegungen der acht `docs-check`-Module
  (`.d-check.yml` Zeile `modules:`, gezählt bei der Anlage) sprengen eine Review-Sitzung — dann `docs-check` allein und die zwei übrigen
  Gates in einem zweiten Slice.
- `in-progress` → `open` (blockiert): der Architect entscheidet, dass die Kante
  für diese Gruppe eine andere Form braucht als im Pilot, und die ADR liegt
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
  sind `Accepted` und tragen ihr `Schärft:`-Feld unveränderlich; ob eine
  Sammel-ADR je Gruppe oder eine ADR je Gate die Kante herstellt, entscheidet
  der Architect. *Zu belegen durch:* die ADR.
- **Vertrag und Werkzeug weichen ab**, sobald die Randform als Festlegung
  formuliert wird (etwa eine Reichweite von `hostpaths`, die
  [`ADR-0075`](../../adr/0075-hostpaths-reichweite-und-wortlaut.md) anders
  nennt als der Vertrag). *Zu belegen durch:* Gegenlesen der Festlegung gegen
  Vertrag und ADR-Kette; ein Befund wird ein eigener Slice.
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
Liefer-Punkt 1 — eine Festlegung, die eine Randform offen lässt, wäre das
dritte Auftreten.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
