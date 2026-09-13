# Slice slice-049: Test-Coverage-Gate

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-14 — erster Slice, unabhängig von `slice-050`.

**Bezug:** [ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
(Scope, Schwelle, Eskalationsklausel, Suppression-Verbot — vorab entschieden).

**Berührte Spec-Stellen:** — (Build-/Gate-Infrastruktur, keine neue
Architektur-Sicht-Aussage).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Eine vierte Docker-Multi-Stage-Stufe `coverage` (nach `deps`,
analog `/Development/d-check/Dockerfile`) misst real
`go test -coverpkg=./internal/...,./cmd/... -coverprofile=… -covermode=atomic
./internal/... ./cmd/...`, prüft das Ergebnis über ein Gate-Skript nach dem
Muster von `/Development/d-check/tools/coverage-gate.sh` gegen die in
[ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
festgelegte Schwelle (Endstufe 80 %, oder — falls der real gemessene
Ist-Stand darunter liegt — eine dokumentierte, auf 5 % abgerundete
Eskalationsstufe als bootstrap-aware Gate), und wird als `make
coverage-gate`-Target in `make gates` verdrahtet. `AGENTS.md` §4 bekommt
die neue Zeile, `harness/README.md` §Sensors die Kalibrierungs-Bindung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Einführung eines Linters** — `ADR-0054` §(a) schließt das für diese
  Welle ausdrücklich aus; `AGENTS.md` §3.2 ist bereits vom Architect-Zug
  schmal ausgefüllt (Suppression-Vollverbot mangels Linter), kein
  weiterer Umsetzungsbedarf in diesem Slice.
- **Eine vorab fixierte Ramp-Stufenfolge** — `ADR-0054` entscheidet
  bewusst gegen erfundene Zwischenstufen; die Eskalationsstufe (falls
  nötig) ist Teil dieses Slices selbst (reale Messung beim ersten Lauf),
  nicht ein separater Folge-Slice.
- **Performance-Benchmark-Infrastruktur** — `slice-050`; andere
  Disziplin (Benchmark statt Pass/Fail-Gate) und andere Schicht
  (eigenständige Mess-Skripte statt Build-/Gate-Infrastruktur).

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

- [ ] Neue `coverage`-Docker-Stage misst real die Go-Test-Coverage über
      `internal/...`+`cmd/...` und scheitert real, wenn ein künstlich
      abgesenkter Schwellenwert unterschritten wird (Rot-Beleg), und
      besteht real bei der tatsächlichen Schwelle (Grün-Beleg).
- [ ] Realer Ist-Stand beim ersten Lauf gemessen und dokumentiert
      (Plan-Nachzug) — Endstufe 80 % direkt, oder dokumentierte
      Eskalationsstufe nach `ADR-0054`.
- [ ] `make coverage-gate` in `make gates` verdrahtet, `harness/README.md`
      §Sensors trägt die Kalibrierungs-Bindung, `AGENTS.md` §4 die neue
      Zeile.
- [ ] `make gates` grün (inkl. des neuen Coverage-Gates).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Sensors, `AGENTS.md` §4 (siehe oben).
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
| `Dockerfile` | update | neue `coverage`-Stage nach `deps` |
| `tools/coverage-gate.sh` | neu | Schwellen-Prüfskript, Muster `/Development/d-check/tools/coverage-gate.sh` |
| `Makefile` (bzw. `harness/mk/*.mk`) | update | `coverage-gate`-Target, Einbindung in `gates:` |
| `harness/README.md` | update | §Sensors-Zeile für `make coverage-gate` |
| `AGENTS.md` | update | §4 neue Zeile |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-14` eröffnet,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei —
unabhängig von `slice-050`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass der real gemessene Ist-Stand so weit unter 80 % liegt, dass die
  Eskalationsstufen-Entscheidung selbst eine tiefere Untersuchung
  braucht (z. B. ganze Pakete ohne jeden Test), gehört das zurück zur
  Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün (inkl. `coverage-gate`) **und**
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der reale Ist-Stand könnte deutlich unter 80 % liegen (unbekannt, kein
  Host-Go-Zugriff) — der Implementer müsste dann ein bootstrap-aware Gate
  mit Eskalationsstufe öffnen, statt das Gate direkt auf 80 % zu setzen
  (`ADR-0054`). **Ausgang:** <bei Closure einzutragen>
- Adapter-Tests, die ohne `CDC_*_TEST_DSN`-Umgebungsvariable real
  überspringen (netzloser Docker-Build), könnten die gemessene Coverage
  künstlich drücken, obwohl die Tests real existieren und in
  `make test-store` grün laufen. **Ausgang:** <bei Closure einzutragen>
- Docker-Layer-Caching könnte einen veralteten Coverage-Lauf über einen
  Cache-Hit hinweg überleben lassen (dasselbe Muster, das d-checks
  `NO_CACHE_FILTER_COV` adressiert). **Ausgang:** <bei Closure
  einzutragen>

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
Keine Treffer für `PGC` zu Coverage/Build-Infrastruktur.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
