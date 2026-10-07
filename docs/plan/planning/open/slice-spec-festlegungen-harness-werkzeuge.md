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
(Entscheidung 7: Baseline-Aktualisierung als bewusster Bootstrap-Vorgang);
Kandidaten für das `Schärft:`-Ziel einer Gate-ADR:
[`ADR-0158`](../../adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md),
[`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md),
[`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
(Auswahl ist Teil der Arbeit, §6). Die Adaption
[`MR-001`](../../../../harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md)
(Rang-2-Datei heißt Pflichtenheft) ist berührt. Welche `LH-*`-Anforderung der
neue Abschnitt verfeinert, steht beim Planen fest, nicht hier.

**Berührte Spec-Stellen:** `spec/pflichtenheft.md` §7 (Historie, wird §8) und
ein neuer §7 „Festlegungen der Harness-Werkzeuge“; `spec/lastenheft.md`
(Zeile mit „Abschnitte 1–7“, Prüfauftrag §6).

**Verantwortlich:** —
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

**Ziel:** `spec/pflichtenheft.md` trägt den Abschnitt „Festlegungen der
Harness-Werkzeuge“ vor der Historie, die Sensor-Verträge und Gate-Zeilen
verweisen für Schwelle und Randform auf ihn, Gate-ADRs schärfen ihn über
`Schärft:`, und die Aussagen „Abschnitte 1–7“ in Spec und `MR-001` haben je
einen Ausgang.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Werkzeug-eigene Teile des Gate-Index (Neu 1)** — Gegenstand von
  `slice-harness-gate-index-werkzeug-teile`: dort der Ort des Gate-Index, hier
  der Ort der Festlegung.
- **`Accepted`-ADRs inhaltlich ändern** — Bestand bleibt: ein `Schärft:`-Feld
  an einer `Accepted`-ADR nachzutragen wäre ein Überschreiben
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); neue Gate-ADRs folgen der neuen
  Vorlage (T3), bestehende bekommen das Feld nur über eine Folge-ADR, falls der
  Architect das entscheidet.
- **Bestehende Sensor-Verträge über den Verweis hinaus umschreiben** — Bestand
  bleibt; der Slice verschiebt Schwelle und Randform in die Spec und lässt im
  Vertrag den Verweis, er schreibt die Verträge nicht neu.
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind `spec/`,
  `harness/sensors/`, `harness/targets/`, `harness/README.md`, `MR-001` bzw. ein
  Nachfolge-Eintrag und der ADR-Index.

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

- [ ] **Abschnitt im Pflichtenheft (Liefer-Punkt 1).** `spec/pflichtenheft.md`
      trägt §7 „Festlegungen der Harness-Werkzeuge“, die Historie ist §8, jeder
      Verweis auf die alte §7 ist nachgezogen; je Gate des Index eine
      Festlegung mit Kennung (R9, T10). *Zu belegen durch:* Suchlauf
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13) auf „§7“, „Historie“,
      „Abschnitte 1–7“ an beiden Ständen und `make docs-check` Exit 0.
- [ ] **Verträge und Index verweisen (Liefer-Punkt 2).** Sensor- und
      Target-Verträge sowie `harness/README.md` §Sensors nennen die
      Spec-Kennung der Festlegung (R2, R7, T6, T7, T8); die Konvention zu
      `Schärft:` im ADR-Index zeigt auf eine Stelle, die es gibt (T4).
- [ ] **Ausgänge für „Abschnitte 1–7“ und `MR-001` (Liefer-Punkt 3).** `MR-001`
      bleibt gültig, wird per Nachfolge-Eintrag abgelöst oder bekommt einen
      anderen Ausgang des Adaptions-Durchgangs, mit Grund (R8); die
      Lastenheft-Zeile bekommt ihren Ausgang (§6).
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
| `spec/pflichtenheft.md` | update | neuer §7, Historie §8 (Liefer-Punkt 1) |
| `harness/sensors/*.md`, `harness/targets/*.md` | update | Verweis auf die Festlegung (Liefer-Punkt 2) |
| `harness/README.md` §Sensors, `docs/plan/adr/README.md` | update | Bindung „Spec-Kennung“, `Schärft:`-Konvention (T4, T6, T7) |
| `harness/conventions/MR-001-…` bzw. Nachfolge-Eintrag | Ausgang des Durchgangs | R8, Liefer-Punkt 3 |
| `spec/lastenheft.md` | Prüfauftrag | Zeile mit „Abschnitte 1–7“ (§6) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): Priorisierung durch den
Auftraggeber; `in-progress/` trägt keinen Slice (WIP-Limit 1). Der Slice hängt
nicht an `slice-harness-gate-index-werkzeug-teile`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die Festlegungen aller Gates sprengen eine
  Review-Sitzung — dann je Gate-Gruppe ein Slice.
- `in-progress` → `open` (blockiert): die Änderung am Lastenheft (Rang 1,
  vertraglich) verlangt eine Freigabe, die nicht vorliegt.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

Die drei Liefer-Punkte sind abgehakt mit Beleg, `make gates` endet mit Exit 0,
der Review-Report liegt vor und ist aufgelöst, die Closure-Notiz trägt den
Lerneintrag und jedes Risiko aus §6 seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Berührung des Lastenhefts (Rang 1).** Der Bump-Plan zählte die Nennung von
  „Abschnitte 1–7“ in `spec/lastenheft.md` als betroffen. **Gemessen** bei der
  Anlage (`grep -n 'Abschnitte 1–7' spec/lastenheft.md`): eine Zeile, 1517, in
  der Historie-Tabelle (Version 0.3.0, „Überführung in
  Lastenheft-Vorlagen-Struktur“) — sie beschreibt die Gliederung des
  Lastenhefts selbst, nicht die des Pflichtenhefts. Ob sie betroffen ist, ist
  am Start zu entscheiden; erwartet: nicht betroffen, Historie bleibt.
- **Umnummerierung bricht Verweise auf §7 des Pflichtenhefts** (Anker,
  Prosa-Nennungen). *Zu belegen durch:* Suchlauf an beiden Ständen und
  `make docs-check`.
- **Welche Gate-ADR schärft welche Festlegung** — Auswahl unter `ADR-0158` bis
  `ADR-0160` und älteren Gate-ADRs; eine `Accepted`-ADR bekommt kein
  nachgetragenes Feld (§1). *Zu belegen durch:* Architect-Artefakt, falls eine
  Folge-ADR nötig ist.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `harness/` und
der ADR-Index; die Modus-Deklaration führt nur die Default-Sub-Area `*`
(Kürzel `PGC`, Greenfield), alle Pfade fallen unter sie.

**Vorgelagert — offene Beobachtungen sichten:** beim Übergang `open → next`
nachzuholen (gemergter Stand des Registers).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
