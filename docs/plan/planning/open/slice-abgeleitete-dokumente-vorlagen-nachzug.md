# Slice abgeleitete-dokumente-vorlagen-nachzug: Harness-Einstieg, Konventionen, ADR-Index und Carveout-Ablage entsprechen ihren Vorlagen v6.14.0

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) (Pin-Inventar P8,
Baseline-Bump). Keine `LH-*`-Anforderung ist berührt: der Slice ändert
Harness-Dokumente, nicht das Produkt.

**Berührte Spec-Stellen:** —

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner, Closure von `slice-baseline-6-14-0-dokumente-nachziehen`).
**Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass.** Die Bestandsaufnahme von `slice-baseline-6-14-0-dokumente-nachziehen`
(§3, *Liefer-Punkt 2c*, gemessen am Stand `d79b7ebd`) nennt vier Dokumente mit
**Abweichung mit Beleg** gegen ihre Vorlage, alle älter als der Bump:
`harness/README.md` (Platzhalter `` `<make-target>` ``, `make <mover>`/`<messung>`/`<vorschau>`,
§Safety and scope boundaries zweimal `<…>`, §Leseordnung drei Platzhalter),
`harness/conventions.md` (`<Pfad oder URL>`, `<Pfade zu deinen …>`, MR-000
`**Datum:** <Datum>`, Musterzeilen in §Zusatzklassen und §Glossar),
`docs/plan/adr/README.md` (kein Abschnitt `## Konventionen`, keine Regel zum Feld
`**Schärft:**`) und `docs/plan/carveouts/` (kein `README.md`). Dazu kommt ein
Träger der Kennungsform nach `MR-002`, den die Closure jenes Slice nicht
nachziehen konnte: das `forbid-pattern` der zwei `structure`-Regeln gegen
Slice-Pfad-Links in `.d-check.yml` trifft nur `slice-[0-9]{3}`.

**Ziel:** Die vier Dokumente tragen keinen Vorlagen-Platzhalter mehr und die
Abschnitte ihrer Vorlage (oder eine benannte, repo-spezifische Abweichung mit
Grund), und die Wächter gegen Slice-Pfad-Links decken die Namens-Kennung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Inhalt der Sensor- und Target-Dateien unter `harness/sensors/` und
  `harness/targets/`** (Abschnitte `## Fassung im Gate-Index`, veraltete Aussagen,
  `baseline-verify.md` Grenze 3) — trägt `slice-harness-targets-inhalt-bereinigen`.
  Beide Slices berühren `harness/README.md`: dieser Slice die Vorlagen-Platzhalter
  (§Sensors-Musterzeile, Werkzeug-Musterzeilen, §Safety, §Leseordnung), jener nur
  einen Index-Kurzsatz, der sich durch eine Berichtigung ändert.
- **`MR-002` auflösen.** Die Vorlage `conventions.template.md` v6.14.0 nennt im
  ID-Schema `slice-<Kennung>`; ob damit der Auflösungs-Trigger von `MR-002`
  eingetreten ist, entscheidet die Änderung an `harness/conventions.md`
  §Adaptions-Block — ein anderer Vorgang, Architect-Zug.
- **Die Regex der E2E-Abdeckungstabelle** (`abdeckungKennungMuster` in
  `test/integration/integration_test.go`, `slice-\d{3}|welle-\d{1,2}`) — Bestand
  bleibt: kein `TestE2E*`-Kommentar trägt eine Namens-Kennung (gemessen bei der
  Closure des Auslöser-Slice, `git grep` über `test/integration`, 0 Treffer), und
  der Kommentar beschreibt die Regex richtig; eine Änderung wäre Produkt-Testcode.
- **Records und `Accepted`-ADRs** mit `slice-<NNN>`/`welle-NN` — eingefroren
  (`AGENTS.md` §3.5).

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

- [ ] **Harness-Einstieg und Konventionen (Liefer-Punkt 1).** `harness/README.md`
      und `harness/conventions.md` tragen keinen Platzhalter der Vorlage; jede
      Musterzeile ist ausgefüllt oder gestrichen, und eine bewusst stehende Zeile
      trägt ihren Grund. Zu belegen durch beide Platzhalter-Formen der Regel in
      `harness/targets/pin-stale.md` (*Bump-Ablauf*, Schritt 2) mit 0 Treffern
      bzw. jedem Treffer mit Grund.
- [ ] **ADR-Index und Carveout-Ablage (Liefer-Punkt 2).** `docs/plan/adr/README.md`
      trägt den Abschnitt `## Konventionen` der Vorlage samt der Regel zu
      `**Schärft:**`; die Spalten `Datum`/`Datei` statt `Bezug` stehen als
      repo-spezifische Abweichung mit Grund oder sind angeglichen.
      `docs/plan/carveouts/README.md` existiert aus `carveouts/README.template.md`.
      Zu belegen durch den Überschriften-`diff` gegen die Vorlage (Exit 0).
- [ ] **Slice-Pfad-Wächter (Liefer-Punkt 3).** Das `forbid-pattern` der zwei
      `structure`-Regeln in `.d-check.yml` trifft auch einen Link auf
      `…/in-progress/slice-<Name>.md`; Grenze (5) im Kommentar und in
      `harness/sensors/docs-check.md` ist gestrichen. Vor der Umsetzung
      **Übergabe an den Architect**: ob die Erweiterung eine ADR braucht
      (`AGENTS.md` §3.6, `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger`).
      Zu belegen durch eine Mutation (ein solcher Link in einer Kopie) mit Befund.
- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Träger, die die geänderten Abschnitte zitieren, sind
      nachgezogen (Suchlauf, `AGENTS.md` §3.13).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, solange die Roadmap unter *Offene Wellen* keine Welle führt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/README.md`, `harness/conventions.md` | update | Liefer-Punkt 1 |
| `docs/plan/adr/README.md`, `docs/plan/carveouts/README.md` | update / neu | Liefer-Punkt 2 |
| `.d-check.yml` (zwei `structure`-Regeln), `harness/sensors/docs-check.md` Grenze | update | Liefer-Punkt 3 |

Der Suchlauf (§3.13) wird beim Start gesetzt; bewegte Eigenschaften: die
Platzhalter der vier Dokumente und die Reichweite des Slice-Pfad-Musters.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1); läuft `slice-harness-targets-inhalt-bereinigen` zuerst, liest
dieser Slice dessen Stand von `harness/README.md` mit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Architect verlangt
  für Liefer-Punkt 3 eine ADR — dann geht Liefer-Punkt 3 in einen eigenen Slice.
- `in-progress` → `open` (blockiert): eine Leseordnung oder §Safety-Grenze lässt
  sich nicht aus dem Bestand ableiten — Frage an den Auftraggeber.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, und die Closure-Notiz in §7 trägt einen Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Ausgefüllte Leseordnung und §Safety erfinden Inhalt**, den kein Träger trägt.
  Gegenmaßnahme: jede Zeile zeigt auf eine bestehende Regel oder Datei. —
  **Ausgang:** offen bis zur Closure.
- **Das erweiterte Muster meldet Bestand.** Gemessen bei der Anlage: der weite
  Ausdruck trifft in den zwei Dateiklassen eine Zeile, in Inline-Code (für die
  Regel unsichtbar). — **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln. Ging der
Gegenstand an einen anderen Slice oder entfiel er, trägt diese Sektion die Zeile
`Gegenstand:` mit Kennung oder Grund.

*Der Plan füllt diese Sektion nicht; sie wird bei der Closure vor dem
`git mv` nach `done/` geschrieben.*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Der Abschnitt selbst entfällt nie.**

**Vorgelagert — Sub-Area-Wahl prüfen:** [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration führt eine Sub-Area (`*`, Kürzel `PGC`, Greenfield); die
berührten Pfade (`harness/`, `docs/plan/adr/`, `docs/plan/carveouts/`,
`.d-check.yml`) liegen in ihr.

**Vorgelagert — offene Beobachtungen sichten:** beim Start durch den Implementer
(Register `docs/plan/planning/observations/BEO-PGC/`). Bekannt bei der Anlage:
`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` (offen, 2×; Liefer-Punkt 3
wäre das dritte Auftreten, wenn die Erweiterung ohne ADR-Klärung käme).

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
