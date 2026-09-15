# Slice slice-080: DB-Adapter-Coverage — eigene, subjekt-qualifizierte Messung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD (eine **zweite,
getrennte** Messung mit eigener Schwelle), kein repo-weites *Mehr*. Ausdrücklich
**nicht** Teil von `welle-20`: deren §6 schließt sie als eigenen Vorgang aus.

**Bezug:** [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 3 (die Messung, ihr Träger und ihre Schwelle, unverändert nach dem
bootstrap-aware-Muster aus [`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a)); [`ADR-0030`](../../adr/0030-testpyramide.md) (die Test-Tiers);
`harness/README.md` §Werkzeuge (die beiden Träger-Läufe).

**Berührte Spec-Stellen:** — (kein Spec-Stratum: eine Messung an bestehenden
Tests, ohne Vertragsänderung).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die drei Pakete, die aus dem Unit-Messgegenstand ausgeschlossen sind
(`postgresstorage` **ohne** das Unterpaket `mapper`, `postgresack`,
`replication/receive`), bekommen eine **eigene Coverage-Zahl** — gemessen dort,
wo ihre Tests ohnehin gegen einen echten PostgreSQL laufen (`make test-store`,
`make test-replication`), mit eigenem `-coverprofile` und eigener, nach
demselben bootstrap-aware-Muster kalibrierter Schwelle. Die Zahl trägt **immer
ihr Subjekt** („DB-Adapter-Coverage") und heißt nie „die Coverage" des Repos
([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Container in `make gates`** — der von `ADR-0071` verworfene Preis
  (Option C): das Gate läuft bei jedem Commit.
- **Die Unit-Zahl** — sie bleibt die Gate-getragene; diese Messung ist eine
  **zweite, getrennte**, und ihre Schwelle gilt nur für ihr Subjekt.
- **Die Tests selbst** — der Slice ändert keinen Test; er **misst** die
  bestehenden. Fehlt Abdeckung, ist das ein Befund dieses Slice, kein Auftrag.
- **Eine Zusammenführung beider Zahlen** — zwei Zahlen, zwei Gegenstände, ein
  Name je Gegenstand; eine gemeinsame Zahl wäre die Doppelquelle, die
  `ADR-0071` Punkt 4 ausschließt.
- **Produktionscode unter `internal/**`** — Schicht-Abgrenzung: Messung, Runner,
  Workflow.

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

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — die Messung existiert.**

- [ ] Die beiden Träger-Läufe erzeugen **je ein** `-coverprofile` für ihren
      Testbestand; die Profile werden **gemergt** und als **eine** Zahl
      ausgegeben — das Merge-Verfahren ist benannt.
- [ ] Die **Herkunft der Zahl** steht dabei (welcher Lauf, welches Profil,
      welche Deduplizierung): die Zählbasis-Regel, die `slice-079` für die
      Unit-Zahl verkörpert hat, gilt hier analog.

**Liefer-Punkt 2 — die Schwelle ist kalibriert.**

- [ ] Der erste reale Wert wird **gemessen** und die Stufe nach dem
      bootstrap-aware-Muster gesetzt (Ist-Stand, abgerundet auf die volle
      5-%-Stufe; `ADR-0054` §(a)) — keine erfundene Zahl, keine Senkung.
- [ ] Ihr **Träger** ist der nicht-blockierende Workflow
      (`.github/workflows/e2e.yml`), **nicht** `make gates`; die Zahl trägt ihr
      Subjekt im Namen und heißt nie „die Coverage".

**Liefer-Punkt 3 — der Beleg.**

- [ ] Ein **realer** Lauf zeigt die Zahl, Exit direkt gelesen und ungepiped
      (`AGENTS.md` §3.9).
- [ ] `make gates` bleibt **unverändert** grün — kein Container dort, keine
      zweite Schwelle im Bündel.

- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
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
| `tools/harness/run-store-tests.sh` | update | `-coverprofile` für den Store-Testbestand, Merge, Ausgabe der Zahl |
| `tools/harness/run-replication-tests.sh` | update | dito für den Replication-Testbestand |
| `harness/sensors/` (neue Datei) | neu | die Bindung dieser **zweiten** Zahl: Subjekt, Schwelle, Träger, Grenzen |
| `harness/README.md` §Werkzeuge | update | die beiden Läufe nennen die Zahl, die sie erzeugen |
| `.github/workflows/e2e.yml` | update | der Träger: die Messung läuft dort, nicht in `make gates` |

**Nicht in dieser Liste:** `Makefile` (kein neues Gate-Target), `internal/**`
(kein Produktionscode), die Tests selbst (sie werden gemessen, nicht geändert).

**Der genaue Zuschnitt der Runner-Änderungen entsteht im ersten Implementer-Lauf**
— die Liste nennt die Träger, nicht jede Zeile.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 3 ist entschieden, die beiden Träger-Läufe existieren
(`make test-store`, `make test-replication`), `.github/workflows/e2e.yml` läuft,
`Verantwortlich:` gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): zeigt sich, dass die
  zwei Läufe **keine gemeinsame** Zahl tragen können (verschiedene Testbestände,
  verschiedene Nenner, kein sinnvolles Merge), ist der Schnitt zu grob — dann
  gehört er zurück zur Zerlegung: **eine Zahl je Lauf** statt einer gemeinsamen.
- `in-progress` → `open` (blockiert — Carveout?): erweist sich, dass die Messung
  den nicht-blockierenden Workflow über seine Laufzeitgrenze treibt (der
  Container-Aufbau steht schon, die Profile nicht), braucht sie eine Entscheidung
  (Carveout oder eigener Träger) — **nicht** stilles Weglassen.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** ein **realer** Lauf zeigt die kalibrierte Zahl **und**
`make gates` unverändert grün **und** die Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die zwei Läufe tragen verschiedene Testbestände.** Ein gemeinsamer Nenner
  über beide kann eine Zahl erzeugen, die keinen der beiden Gegenstände
  beschreibt. Antwort: eine Zahl **je Lauf** oder ein ausdrücklich benannter
  gemeinsamer Nenner samt Verfahren. — **Ausgang:** <bei Closure>
- **Die erste Kalibrierung könnte sehr niedrig liegen.** Das ist zulässig
  (bootstrap-aware: die Stufe folgt der Messung) — aber der Wert gehört
  **benannt**, samt dem, was ihn drückt. — **Ausgang:** <bei Closure>
- **`postgresstorage/mapper` liegt in beiden Messungen** — es bleibt im
  Unit-Gegenstand. Die zwei Zahlen **überlappen** in diesem Paket; wer sie
  addiert oder vergleicht, muss das wissen. — **Ausgang:** <bei Closure>

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
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt `tools/harness/**`,
`harness/sensors/**`, `harness/README.md` und `.github/workflows/` in **einem**
Kürzel. Eine feinere Sub-Area ist nicht zu bilden: „Test-Infrastruktur" ist keine
deklarierte Sub-Area, und das Register führt den Runner unter demselben Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen, gemergter Stand; zwei Treffer:

- `BEO-PGC/regel-weiter-als-ihr-sensor` (2×, offen): **kein** Treffer — die
  beiden Belege betreffen die host-lokale Pfad-Regel und die **Unit**-Fitness-
  Function. Dieser Slice fügt eine **dritte, eigene** Messung hinzu und ändert
  keine der beiden Zusagen.
- `BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (1×,
  `verkörpert` → `ADR-0071`): der **Vorläufer** — dieser Slice setzt seinen
  Punkt 3 um. **Kein neuer Beleg**: derselbe Gegenstand, derselbe Träger.

**Kein** Eintrag erreicht mit diesem Slice die 3×-Schwelle; es entsteht kein
neues Verzeichnis.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (die beiden Träger-Läufe existieren,
`harness/sensors/` ist die Bindungsform), **Phase-Reife** hoch (das Gate läuft
seit `slice-049`, die Träger-Läufe seit `slice-050`), **Evidenz-/Diskrepanz-Risiko**
niedrig — die Zahl wird **gemessen**, nicht gegen einen Bestand inventarisiert —,
**Reconciliation-Aufwand** null.
