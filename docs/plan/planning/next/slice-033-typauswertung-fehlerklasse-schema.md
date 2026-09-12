# Slice slice-033: Typ-Auswertung und Fehlerklasse `schema` für inkompatible Typänderungen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-10`](../welle-10.md) — der Nachweis, dass die
`ADR-0015`-Folgepflicht real eingelöst ist, ist `welle-10`s Closure-Trigger
(§3); dieser Slice ist der letzte der drei geplanten Slices und schließt
`LH-FA-SCH-004`s Negative-Fall real.

**Bezug:** [`LH-FA-SCH-004`](../../../../spec/lastenheft.md), `ADR-0015`
(nur gelesen — keine aktive ADR wird geändert, `ADR-0015` bleibt
`Accepted`).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Fehlerklasse `schema`), `LH-FA-SCH-004.a`.

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-09-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `Assembler.observeRelation`
(`internal/adapters/driving/replication/mapper/mapper.go`) meldet den
Fall `relationOther` (Spalte entfernt, Typ einer bestehenden Spalte
geändert, Spalte umbenannt — jede nicht sicher als reine
Obermengen-Erweiterung erkennbare Relation-Änderung) ab sofort als
**sichtbaren Fehler der Fehlerklasse `schema`**
([`SPEC-008`](../../../../spec/pflichtenheft.md)), statt ihn wie bisher
(`slice-032`) konservativ stillschweigend zu ignorieren. Das schließt
`LH-FA-SCH-004`s Negative-Fall real: „Given eine inkompatible
Typänderung, when sie eintritt, then ist sie erkennbar gemeldet; die
Daten werden nicht still fehlinterpretiert." `TestMVPSchemaChangeIncompatibleTypeChange`
(`slice-030`) läuft danach mit angepasster Erwartung an seinen zweiten
Testfall (PostgreSQL-seitig zugelassene Typänderung) grün: **statt**
stillschweigender Übernahme erwartet er jetzt einen sichtbaren
`schema`-Fehler. `TestMVPSchemaChangeAddColumn` (kompatible Erweiterung,
bereits von `slice-032` real geschlossen) bleibt unverändert im
Verhalten — dieser Slice rührt nur den `relationOther`-Zweig an.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Entfernte Spalten differenziert behandeln** (`LH-FA-SCH-003`) — Bestand
  bleibt bewusst stehen: eine entfernte Spalte fällt unter denselben
  `relationOther`-Zweig wie jede andere nicht sicher interpretierbare
  Änderung; eine feinere, `LH-FA-SCH-003`-spezifische Behandlung (z. B.
  definiertes Verhalten für historische Werte) ist ein eigener, hier
  ausgeschlossener Vorgang (bereits in `slice-030` §1 als Out-of-Scope
  benannt).
- **Wiederherstellung/Recovery nach einem `schema`-Fehler** (z. B.
  automatischer Neustart mit neuer Baseline-Version) — anderer Vorgang:
  dieser Slice meldet den Fehler sichtbar, definiert aber nicht, wie ein
  Operator danach fortsetzt (das ist bereits heute das etablierte
  Verhalten anderer `schema`-Klassen-Fehler in diesem Repo, z. B.
  `ErrTruncateUnsupported` — konsistentes Verhalten, keine neue
  Recovery-Fähigkeit).
- **Änderung der Klassifikationslogik selbst
  (`classifyRelationColumns`)** — Bestand bleibt bewusst stehen: die
  Drei-Fälle-Unterscheidung aus `slice-032` ist korrekt und bereits
  real verifiziert; dieser Slice ändert nur, **was im dritten Fall
  passiert**, nicht die Klassifikation selbst.

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

- [ ] `observeRelation` meldet `relationOther` als sichtbaren Fehler der
      Fehlerklasse `schema` (Muster: `Fehlerklasse schema: …`, analog zu
      `ErrTruncateUnsupported`/`ErrChangeWithoutBegin` in derselben
      Datei) statt `nil` zurückzugeben. Unit-getestet (Assembler-Tests).
- [ ] `LH-FA-SCH-004`s Negative-Fall real geschlossen:
      `TestMVPSchemaChangeIncompatibleTypeChange` (`slice-030`) läuft mit
      angepasster Erwartung grün — der PostgreSQL-seitig zugelassene
      Typänderungs-Fall erwartet jetzt einen sichtbaren `schema`-Fehler
      statt stillschweigender Übernahme. `TestMVPSchemaChangeAddColumn`
      bleibt unverändert grün (Regressionsschutz für `slice-032`s
      Ergebnis). `make gates` und `make test-integration` dreimal in
      Folge grün.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt wird —
      Implementer entscheidet und begründet im Plan-Nachzug.
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
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `observeRelation` meldet `relationOther` als sichtbaren `schema`-Fehler |
| `test/integration/integration_test.go` | update | `TestMVPSchemaChangeIncompatibleTypeChange`s zweiter Testfall erwartet jetzt einen sichtbaren Fehler |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-032` liegt in `done/`,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass eine differenzierte Fehlerbehandlung je Unterfall (entfernte
  Spalte vs. Typänderung vs. Umbenennung) statt einer einheitlichen
  `schema`-Fehlerklasse nötig ist, gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
dreimal in Folge grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein sichtbarer `schema`-Fehler bei jeder nicht als Obermenge erkennbaren
  Änderung (auch harmlose Umbenennungen ohne Datenverlust) könnte in der
  Praxis zu häufig auslösen und den Erfassungspfad unnötig hart stoppen —
  eine feinere Unterscheidung wäre eine größere, hier nicht geleistete
  Fähigkeit. **Ausgang:** <bei Closure einzutragen>
- Der bestehende Fehlerbehandlungspfad (wie ein `schema`-Fehler aus
  `Consume` den Erfassungspfad tatsächlich beendet/meldet) könnte
  Annahmen treffen, die für einen während des laufenden Streams
  auftretenden Fehler (statt eines Fehlers beim initialen Decode) nicht
  zutreffen. **Ausgang:** <bei Closure einzutragen>

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
Treffer für `PGC`: `BEO-PGC/schema-evolution-nicht-dynamisch` (1×, weiter
offen — dieser Slice liefert den dritten und letzten Baustein; die
Auflösung selbst entscheidet der Planner bei `welle-10`s Closure). Keiner
der übrigen Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
