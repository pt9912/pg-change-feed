# Slice slice-050: Performance-Benchmark-Infrastruktur

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-14 — zweiter Slice, unabhängig von `slice-049`.

**Bezug:** [LH-QA-PER-001](../../../../spec/lastenheft.md),
[LH-QA-PER-002](../../../../spec/lastenheft.md),
[LH-QA-PER-003](../../../../spec/lastenheft.md),
[ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
(Muster, kein Gate).

**Berührte Spec-Stellen:** — (Mess-Skript-Ergänzung, keine neue
Architektur-Sicht-Aussage).

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Drei eigenständige Bench-Skripte nach dem Stil von
`/Development/d-check/tools/bench-fixture.sh`
([ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(b)), gebündelt hinter einem gemeinsamen `make bench`-Target (analog
`/Development/d-check/Makefile`), das explizit **kein Gate** ist:
(1) `LH-QA-PER-001` Quell-Impact mit/ohne CDC (dieselbe Schreiblast auf
die Quelltabelle, einmal mit aktivem Replication-Slot/Capture-Prozess,
einmal ohne), (2) `LH-QA-PER-002` Skalierung über die drei
[SPEC-014](../../../../spec/pflichtenheft.md)-Lastenstufen (klein ≤10/s,
mittel 100/s×30min, groß 1000/s×60min) als feste Eingabeparameter,
(3) `LH-QA-PER-003` Batch- vs. Einzelabruf-Effizienz beim Lesen über
`cdc.changes`. Jedes Skript dokumentiert Aufwand/Ergebnis (Report-Datei
oder stdout), ohne einen Pass/Fail-Schwellenwert durchzusetzen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Aufnahme in `make gates`/`fullbuild` als Pass/Fail-Bedingung** —
  `ADR-0054` §(b) legt Benchmark ausdrücklich als dokumentierten Beleg
  fest, nicht als Gate; andere Disziplin als das Coverage-Gate.
- **Neue Lastenstufen-Definition** — `SPEC-014` legt die drei Stufen
  bereits fest; dieser Slice übernimmt sie als Eingabeparameter, ändert
  sie nicht.
- **Test-Coverage-Gate** — `slice-049`; andere Schicht (Build-/Gate-
  Infrastruktur statt eigenständiger Mess-Skripte).
- **`LH-QA-PER-004`-Ausbau** (Commit→CDC-Latenz über den bestehenden
  Lasttest-Beleg hinaus) — `ADR-0054`s Re-Evaluierungs-Trigger (c) benennt
  einen eigenen Folge-Vorgang, falls das je nötig wird; bleibt hier
  unverändert Bestand.

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

- [ ] `LH-QA-PER-001` real erfüllt: Bench-Skript zeigt real einen
      messbaren Unterschied (oder dessen Abwesenheit) zwischen
      Schreiblast mit und ohne aktivem Replication-Slot/Capture-Prozess,
      dokumentiertes Ergebnis.
- [ ] `LH-QA-PER-002` real erfüllt: Bench-Skript durchläuft real alle
      drei `SPEC-014`-Lastenstufen, dokumentiertes Ergebnis je Stufe.
- [ ] `LH-QA-PER-003` real erfüllt: Bench-Skript vergleicht real
      Batch- vs. Einzelabruf beim Lesen über `cdc.changes`,
      dokumentiertes Ergebnis.
- [ ] `make bench` startet alle drei Skripte, kein Gate (Aufnahme in
      `make gates` explizit unterlassen), `harness/README.md`
      §Werkzeuge trägt die neue Zeile.
- [ ] `make gates` grün (unverändert, da kein Gate hinzukommt).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Werkzeuge (siehe oben).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
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
| `tools/bench-source-impact.sh` | neu | `LH-QA-PER-001`-Beleg |
| `tools/bench-scaling.sh` | neu | `LH-QA-PER-002`-Beleg |
| `tools/bench-batch-vs-single.sh` | neu | `LH-QA-PER-003`-Beleg |
| `Makefile` (bzw. `harness/mk/*.mk`) | update | `bench`-Target, kein Gate |
| `harness/README.md` | update | §Werkzeuge-Zeile für `make bench` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-14` eröffnet,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei —
unabhängig von `slice-049`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Drei
  eigenständige Bench-Skripte plus Makefile-Wiring nähern sich der
  Drei-Liefer-Punkte-Grenze — zeigt sich, dass eines der drei Skripte
  selbst schon einen eigenen Liefer-Punkt umfangreicher wird (z. B.
  `LH-QA-PER-002`s drei Lastenstufen brauchen eine eigene
  Fixture-Erzeugung je Stufe statt eines Parameters), gehört das zurück
  zur Zerlegung (ein Skript je Slice statt aller drei zusammen).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make bench` liefert real alle drei Belege **und**
`make gates` unverändert grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Benchmark-Ergebnisse könnten in einer geteilten/virtualisierten
  Docker-Umgebung (kein dediziertes Hardware-Budget) real streuen —
  dasselbe Problem, das d-checks `bench:`-Target über N=3-Läufe und
  Median statt Einzelmessung adressiert. **Ausgang:** <bei Closure
  einzutragen>
- `LH-QA-PER-002`s „groß"-Lastenstufe (1.000/s × 60 Minuten,
  `SPEC-014`) könnte `make bench` für einen schnellen, wiederholten
  Implementer-/Reviewer-Lauf unpraktikabel lang machen. **Ausgang:**
  <bei Closure einzutragen>

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
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

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Keine Treffer für `PGC` zu Performance/Benchmark.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
