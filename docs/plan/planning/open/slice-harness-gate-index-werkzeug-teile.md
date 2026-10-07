# Slice harness-gate-index-werkzeug-teile: Werkzeug-eigene Teile des Gate-Index übernehmen

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
(Entscheidung 7: Baseline-Aktualisierung als bewusster Bootstrap-Vorgang). Keine
`LH-*`-Anforderung ist berührt: der Slice ändert Harness-Dokumente und die
Gate-Konfiguration, nicht das Produkt. Die Aktivierung des d-check-Moduls
`targets` braucht eine eigene ADR ([`AGENTS.md`](../../../../AGENTS.md) §3.6);
sie gehört zu diesem Slice (§2).

**Berührte Spec-Stellen:** — (keine; Gate-Index und Gate-Konfiguration).

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
Abschnitt „Bump-Ablauf — Belege“) gab den Delta-Punkten R3, R4, R5, R6, R10,
T1, T2, T5 und T9 den Ausgang „Folge-Slice Neu 1“: werkzeug-eigene Teile des
Gate-Index `harness/mk/<werkzeug>.md` (`regelwerk/grundlagen-harness-dateien.md`
„Ein Index, mehrere Eigentümer“; `.d-check.yml` `targets` mit
`makefiles`-Glob und `doc-tables`-/`authority`-Liste, d-check ≥ v0.82.0). Dort
**gemessen**: kein veröffentlichter Stand von `ai-harness-init` (bis v0.2.7)
erzeugt `harness/mk/<werkzeug>.md`; installiert ist v0.2.3. Den Teil schreibt
nach der Regel das Werkzeug, nicht das Repo.

**Ziel:** Sobald `ai-harness-init` die werkzeug-eigenen Teile erzeugt, führt das
Repo seinen Gate-Index als Vereinigung aus `harness/README.md` §Sensors und den
erzeugten `harness/mk/<werkzeug>.md`, ohne Ziel doppelt zu führen; die
Delta-Punkte R3–R6, R10, T1, T2, T5, T9 sind übernommen, und das Modul `targets`
ist per ADR aktiviert oder begründet nicht aktiviert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Festlegungen der Harness-Werkzeuge in der Spezifikation (Neu 2)** —
  Gegenstand von `slice-spec-festlegungen-harness-werkzeuge`: dort berührt die
  Arbeit das Lastenheft und die Spec-Straten, hier den Gate-Index.
- **Die Teile selbst schreiben** — ein anderer Vorgang: die Regel schreibt sie
  dem Werkzeug zu („das Repo schreibt nicht hinein“, R3); eine Änderung an
  `ai-harness-init` liegt nicht in diesem Repo.
- **Ein Werkzeug-Update ohne erzeugten Teil** — Bestand bleibt: ein Bump von
  `ai-harness-init`, der keine `.md`-Teile erzeugt, ist kein Anlass für diesen
  Slice.
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind `harness/`,
  `.d-check.yml`, `AGENTS.md` §4, `Makefile`-Kopf und eine ADR.

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

- [ ] **Werkzeug-Teile im Repo (Liefer-Punkt 1).** `ai-harness-init` auf einem
      Stand, der `harness/mk/<werkzeug>.md` erzeugt; die erzeugten Teile sind
      committet, und kein Ziel steht zugleich in `harness/README.md` §Sensors
      und in einem Teil. *Zu belegen durch:* `--version`, `ls harness/mk/*.md`
      und einen Abgleich der Ziel-Namen beider Quellen (Befehl und Zahl im Plan).
- [ ] **Gate-Entscheidung zu `targets` (Liefer-Punkt 2).** Eine ADR entscheidet
      über die Aktivierung des d-check-Moduls `targets` (`makefiles`,
      `doc-tables`, `authority`); bei Aktivierung `.d-check.yml` danach und
      `make docs-check` Exit 0. *Zu belegen durch:* die ADR (`Accepted`) und den
      Gate-Lauf.
- [ ] **Delta-Punkte übernommen (Liefer-Punkt 3).** R3–R6, R10 (Regelwerk) und
      T1 (`AGENTS.md` §4), T2 (`.d-check.yml`), T5 (Kommentar-Block in
      `harness/README.md` §Sensors), T9 (`Makefile`-Kopf) je mit Ausgang
      „übernommen in `<Datei>`“ oder „betrifft das Repo nicht“ mit Grund.
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
| `harness/mk/<werkzeug>.md` | neu (erzeugt vom Werkzeug) | Liefer-Punkt 1 |
| `harness/README.md` §Sensors | update | Vereinigung, Disjunktheit, T5 |
| `docs/plan/adr/<NNNN>-…` und ADR-Index | neu (Architect) | Liefer-Punkt 2 |
| `.d-check.yml` | update falls aktiviert | T2, Liefer-Punkt 2 |
| `AGENTS.md` §4, `Makefile`-Kopf | update | T1, T9 |

Der Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13) ist beim Start zu
messen; bewegte Eigenschaft ist der Ort des Gate-Index.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): ein veröffentlichter Release von
`ai-harness-init` erzeugt `harness/mk/<werkzeug>.md` (Release-Text oder Quelle
gelesen; v0.2.7 tut es nicht, gemessen in `slice-harness-baseline-v6-16-0`);
`in-progress/` trägt keinen Slice (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die ADR zu `targets` verlangt Nachzüge an
  mehr als einem Gate oder an Sensor-Verträgen über `harness/README.md` hinaus.
- `in-progress` → `open` (blockiert): das Werkzeug erzeugt Teile, die ein Ziel
  doppelt mit `harness/README.md` führen und nur am Werkzeug zu beheben sind.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

Die drei Liefer-Punkte sind abgehakt mit Beleg, `make gates` endet mit Exit 0,
der Review-Report liegt vor und ist aufgelöst, die Closure-Notiz trägt den
Lerneintrag und jedes Risiko aus §6 seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Der Auslöser tritt nie ein** (kein Release erzeugt die Teile). Dann bleibt
  dieser Slice in `open/`; ein Wegfall der Regel in einer späteren Baseline wäre
  ein Ausgang `open → done` mit `Gegenstand: entfallen`. *Zu belegen durch:*
  den Bump-Ablauf des nächsten Baseline-Bumps.
- **Die Aktivierung von `targets` färbt Bestand rot** (Ziele ohne Index-Zeile
  oder doppelt geführt). *Zu belegen durch:* `make docs-check` mit aktiviertem
  Modul am Arbeitsbaum, vor der ADR.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `harness/`,
`.d-check.yml`, `AGENTS.md` und eine ADR; die Modus-Deklaration führt nur die
Default-Sub-Area `*` (Kürzel `PGC`, Greenfield), alle Pfade fallen unter sie.

**Vorgelagert — offene Beobachtungen sichten:** beim Übergang `open → next`
nachzuholen (gemergter Stand des Registers); bei Anlage gelesen:
`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` (verkörpert,
`harness/targets/pin-stale.md` §Bump eines Gate-Werkzeugs) trifft
Liefer-Punkt 2.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
