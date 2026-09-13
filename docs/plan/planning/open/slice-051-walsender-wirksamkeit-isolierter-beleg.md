# Slice slice-051: Publication-Entzug-Wirksamkeit — isolierter Beleg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Architect-Verdikt
([`docs/reviews/architect-verdict-walsender-wirksamkeit.md`](../../../reviews/architect-verdict-walsender-wirksamkeit.md))
und die vorausgehende Fork-Recherche fanden übereinstimmend: kein Mehr
über die eigene DoD hinaus (ein einzelner isolierter Testabschnitt, keine
Schema-/Domänen-Änderung), siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

**Bezug:** [LH-FA-CFG-002](../../../../spec/lastenheft.md) (CDC-Deaktivierung),
[ADR-0050](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antragsqueue, Live-Reload — nur referenziert, nicht geändert).

**Berührte Spec-Stellen:** — (reine Testabdeckungs-Ergänzung, keine neue
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

**Ziel:** Ein neuer, isolierter Testabschnitt in
`tools/harness/run-integration-tests.sh` belegt real, ob PostgreSQLs
bereits laufende logische Decoding-Session eine Tabelle sofort oder
verzögert ausfiltert, nachdem sie per `ALTER PUBLICATION ... DROP TABLE`
entfernt wurde — **ohne** den regulären `cdc.disable_table`-Antragsweg zu
nutzen, damit die App-seitige `Assembler`-Filterung als
Alternativerklärung ausgeschlossen ist (exakter Ablauf: Architect-Verdikt
oben, Abschnitt „Der isolierte Testansatz"). Je nach realem Ausgang
schließt dieser Slice `BEO-PGC/walsender-wirksamkeit` entweder mit
Ausgang *entfallen* (PostgreSQL filtert sofort) oder trägt einen
benannten Liefer-Punkt für den Fall, dass der Walsender real verzögert
liefert (z. B. Doku-Klarstellung, dass die App-seitige Filterung die
tragende Ebene ist).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ersatz des bestehenden „SQL-Administration Live-Reload-Belegs
  (disable)"** — der Architect-Verdikt stellt klar: beide Belege prüfen
  unterschiedliche Eigenschaften (App-seitige Idempotenz vs. reines
  PostgreSQL-Walsender-Timing); der bestehende Beleg bleibt unverändert
  bestehen.
- **Neuer Introspektions-Mechanismus** (`pg_logical_slot_peek_changes`
  o. ä.) — der Architect-Verdikt hat das geprüft und verworfen (Slot ist
  exklusiv an den laufenden Feed-Prozess gebunden, ein zweiter Slot
  würde die eigentliche Frage verfehlen).
- **Änderung an `Assembler`/`DisableTableUseCase`/`ADR-0050`** — dieser
  Slice ist reine Testabdeckung auf bereits bestehendem Verhalten, keine
  Verhaltensänderung. Ein etwaiger Grace-Wait oder eine
  Doku-Klarstellung (falls der Walsender real verzögert) ist der einzig
  mögliche Code-/Doku-Berührungspunkt und bleibt klein genug, um im
  selben Slice zu bleiben (siehe §2).

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

- [ ] Neuer, isolierter Testabschnitt in `run-integration-tests.sh` real
      ausgeführt: dedizierte, über die reguläre `cdc.enable_table`-Kette
      aktivierte Tabelle, `ALTER PUBLICATION ... DROP TABLE` direkt per
      `psql` (ohne `cdc.disable_table`), neue Zeile eingefügt, reales
      Ergebnis gegen `cdc.changes` dokumentiert (erscheint/erscheint
      nicht).
- [ ] `LH-FA-CFG-002` real belegt in Bezug auf die Walsender-Timing-Frage
      — Ergebnis eindeutig einer der beiden Auswertungen aus dem
      Architect-Verdikt zugeordnet.
- [ ] Falls der Walsender real verzögert liefert: ein benannter
      Liefer-Punkt (Grace-Wait-Doku oder Klarstellung, dass die
      App-seitige Filterung die tragende Ebene ist) — falls PostgreSQL
      real sofort filtert: entfällt dieser Punkt ersatzlos (§1 „Keine
      Mindestzahl"-Prinzip sinngemäß auf DoD-Punkte übertragen; der
      Implementer trägt im Plan-Nachzug nach, welcher Fall eintrat).
- [ ] `make gates` grün, `make test-integration` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: nur falls der Walsender real verzögert (siehe oben) —
      Implementer prüft und begründet im Plan-Nachzug.
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
| `tools/harness/run-integration-tests.sh` | update | neuer isolierter Walsender-Timing-Testabschnitt |
| `docs/user/benutzerhandbuch.md` | update, falls Walsender real verzögert | Implementer entscheidet, siehe §2 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Architect-Verdikt liegt vor
(`docs/reviews/architect-verdict-walsender-wirksamkeit.md`),
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass der reale Walsender-Verzögerungs-Fall eine über eine
  Doku-Klarstellung/einen Grace-Wait hinausgehende Architektur-Änderung
  braucht, gehört das zurück zur Zerlegung (neuer Architect-Zug nötig,
  nicht mehr Testarbeit).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Testabschnitt weist eine Abwesenheit nach (kein Poll auf ein
  eintretendes Ereignis, sondern ein reales, begrenztes `sleep` — Muster
  Architect-Verdikt Punkt 4) — eine zu kurze Wartezeit könnte einen real
  verzögerten Walsender fälschlich als „sofort filternd" auswerten.
  **Ausgang:** <bei Closure einzutragen>
- Isolation gegenüber dem bestehenden „SQL-Administration Live-Reload-
  Beleg (disable)" — beide Abschnitte dürfen sich nicht gegenseitig
  stören (eigene, dedizierte Tabelle nötig, Architect-Verdikt Punkt 1).
  **Ausgang:** <bei Closure einzutragen>
- Bestätigt sich real die Verzögerungs-Annahme (Walsender liefert trotz
  entzogener Publication weiter), könnte der benannte Liefer-Punkt
  (Grace-Wait/Doku-Klarstellung) den realen Umfang unterschätzen — dann
  greift die in §4 vorab benannte Rückführung `in-progress` → `next`.
  **Ausgang:** <bei Closure einzutragen>

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
Treffer für `PGC`: `BEO-PGC/walsender-wirksamkeit` (1×, dieser Slice
liefert den Testbeleg, der ihren Ausgang bestimmt), `BEO-PGC/test-isolation-geteilter-zustand`
(1×, direkt einschlägig — siehe §6), `BEO-PGC/test-runner-stiller-ausschluss`
(1×, jeder neue Abschnitt muss real mitlaufen). Keiner erreicht mit
diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
