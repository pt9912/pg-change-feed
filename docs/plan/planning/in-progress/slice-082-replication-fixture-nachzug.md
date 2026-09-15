# Slice slice-082: Replication-Fixture nachziehen — `make test-replication` wieder grün

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD (ein Beleg
läuft wieder grün), kein repo-weites *Mehr*. Nicht Teil von `welle-20`: deren
Gegenstand ist die netzlos prüfbare Fläche, dieser Slice die DB-gestützte Ebene.

**Bezug:** [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(`bootstrap.Run` liest `cdc.administration_request` und `cdc.process_heartbeat`
seit dieser Entscheidung — die Tabellen, die dem Fixture fehlen);
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 3 (die DB-Adapter-Coverage, deren CI-Schritt deswegen derzeit rot endet);
`BEO-PGC/roter-test-ohne-leser` (die Klasse: ein Beleg, den kein Gate abholt).

**Berührte Spec-Stellen:** — (Test-Fixture und Testpyramide,
[`ADR-0030`](../../adr/0030-testpyramide.md); kein Spec-Stratum).

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

**Ziel:** `make test-replication`s tier-weites `go test ./...` läuft **grün**. Das
Fixture von `internal/bootstrap · TestWALRetentionThresholdEndToEnd` fährt heute
`DROP SCHEMA cdc CASCADE` und ein eigenes `ApplySchema`, das die Tabellen
`cdc.administration_request` und `cdc.process_heartbeat` nicht mitbringt — beide
liest `bootstrap.Run` seit [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md).
Der Test scheitert darum mit `42P01`, und weil `make test-replication` **weder im
Gate noch in der CI** läuft, war das über viele Slices unsichtbar. Sichtbar wurde
es erst, als die DB-Adapter-Coverage den Lauf als Träger aufrief.

**Der Ist-Zustand ist belegt:** der Fehler reproduziert am **unveränderten**
Runner (`9fa9be3`), er kann nicht aus einem späteren Diff stammen, und
`internal/**` ist in jenem Diff unberührt (Review `review-slice-080` F-2/F-4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Maskierung.** Kein `|| true`, kein Entfernen des Tests, kein Abschwächen
  der Zusicherung: der Beleg soll **grün werden**, nicht **unsichtbar**.
- **Produktionscode jenseits des Fixtures** — geändert wird der **Test**-Aufbau,
  nicht `bootstrap.Run`; liest der Bootstrap etwas, das ein Fixture nicht
  mitbringen kann, ist das ein Befund dieses Slice, kein Auftrag.
- **Die DB-Adapter-Coverage selbst** (`slice-080`): ihre Zahl und ihre Schwelle
  bleiben, wie sie sind; dieser Slice macht ihren CI-Schritt nur
  beobachtbar-grün.
- **Ein neuer Gate-Aufruf für `test-replication`** — ob dieses Target ins
  Bündel gehört, ist die *übergeordnete* Frage (die Klasse
  `BEO-PGC/roter-test-ohne-leser`), und sie gehört als eigene Entscheidung
  geführt, nicht als Beigabe.
- **Die netzlos prüfbare Fläche** (`welle-20`): unberührt.

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

**Liefer-Punkt 1 — der Fixture-Aufbau kommt aus einer Quelle.**

- [ ] Das Fixture der betroffenen Testdatei bringt den Schema-Stand mit, den
      `bootstrap.Run` liest — und zwar aus **derselben** Quelle wie der Betrieb
      (der Schema-Rollout), nicht aus einer zweiten, handgepflegten Liste. Eine
      zweite Liste ist genau die Drift, die diesen Befund erzeugt hat.
- [ ] Die zwei Tabellen, deren Fehlen belegt ist (`cdc.administration_request`,
      `cdc.process_heartbeat`), sind nach dem Nachzug **real** vorhanden — der
      Lauf beweist es, nicht der Kommentar.

**Liefer-Punkt 2 — der Lauf ist grün.**

- [ ] `make test-replication` (die **Tier**-Hälfte) läuft real **Exit 0** —
      Exit-Code direkt gelesen und ungepiped (`AGENTS.md` §3.9).
- [ ] **Keine Maskierung:** der Test ist unverändert in seiner Zusicherung, kein
      `|| true`, kein `t.Skip` als Ersatz für den Nachzug.

**Liefer-Punkt 3 — der Beleg ist wieder lesbar.**

- [ ] Der CI-Schritt (`measure` **und** `tier`) ist damit **grün beobachtbar**;
      die DB-Adapter-Coverage-Zahl bleibt unverändert (sie war nie rot).
- [ ] `make gates` grün.

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
| `internal/bootstrap/**` (die betroffene Testdatei samt Fixture) | update | der Schema-Aufbau kommt aus der einen Quelle statt aus einer eigenen Teilmenge |
| `tools/schema/**` — **nur falls** der Nachzug eine geteilte Schema-Anwendung braucht | update | dann gehört sie an eine Stelle, die Betrieb und Test gemeinsam nutzen; der Implementer entscheidet es am Artefakt und begründet es |

**Nicht in dieser Liste:** `internal/bootstrap/wiring.go` (Produktionscode —
liest der Bootstrap etwas, das ein Fixture nicht mitbringen kann, ist das ein
Befund, kein Auftrag), `tools/harness/db-coverage.sh` und die
Träger-Läufe aus `slice-080` (ihre Zahl und ihre Schwelle bleiben), `Makefile`
(kein neuer Gate-Aufruf), `spec/**`.

**Der genaue Zuschnitt entsteht im ersten Implementer-Lauf** — die Liste nennt
die Träger, nicht jede Zeile.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der rote Lauf ist **belegt** (Review
`review-slice-080` F-2/F-4 samt Nachweis am unveränderten Runner), die
Container-Umgebung steht (`make test-replication` aufrufbar), `Verantwortlich:`
gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): zeigt sich, dass der
  Nachzug **mehrere** Fixtures derselben veralteten Form betrifft, gehört er
  zurück zur Zerlegung — **nach Fixture**, nicht als Sammelumbau.
- `in-progress` → `open` (blockiert — Carveout?): erweist sich, dass der Test nur
  grün wird, wenn `bootstrap.Run` **anders liest** (also Produktionscode geändert
  werden müsste), ist das ein Blocker mit Entscheidung — die Zusage dieses Slice
  ist „Fixture nachziehen, Verhalten unberührt".

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** die Tier-Hälfte läuft real **grün** (Exit 0) **und** die
Zusicherung des Tests ist unverändert **und** `make gates` grün **und** die
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Das Fixture könnte erneut driften** — dieselbe Klasse, die diesen Befund
  erzeugt hat: ein Test-Aufbau, der den Schema-Stand selbst pflegt, driftet
  gegen den Betrieb. Die Antwort ist die **eine** Quelle (der Schema-Rollout);
  bleibt daneben eine zweite Liste stehen, ist der Befund nur verschoben. —
  **Ausgang:** <bei Closure>
- **Die Grüne könnte durch Abschwächen entstehen.** In §1 ausgeschlossen; der
  Wächter ist das Review (Liefer-Punkt 2 verlangt die unveränderte Zusicherung).
  — **Ausgang:** <bei Closure>
- **Der Beleg braucht mehrere reale Läufe** (Container-Aufbau je Lauf). Ein Lauf,
  der nur grün wird, weil ein Container-Rest stehen blieb, ist kein Beleg. —
  **Ausgang:** <bei Closure>

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
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt `internal/bootstrap/**`
(die Testdatei) und `tools/schema/**` in **einem** Kürzel. Eine feinere Sub-Area
ist nicht zu bilden: „Test-Infrastruktur" ist keine deklarierte Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen, gemergter Stand; zwei Treffer:

- `BEO-PGC/roter-test-ohne-leser` (1×, offen): **Treffer** — dieser Slice ist
  der **Träger** des Instanz-Fixes. Die **Klasse** (ein Target, das kein Gate
  abholt) bleibt offen und wird hier **nicht** gezählt: die Behebung ist kein
  zweites Auftreten („ein Vorgang zählt einmal").
- `BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg` (2×, offen): **kein**
  Treffer — dieser Slice erklärt keinen Werkzeug-Mechanismus; er zieht einen
  Schema-Aufbau nach.

**Kein** Eintrag erreicht mit diesem Slice die 3×-Schwelle; es entsteht kein
neues Verzeichnis.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (die Testpyramide ist über
[`ADR-0030`](../../adr/0030-testpyramide.md) und
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) verankert),
**Phase-Reife** hoch, **Evidenz-/Diskrepanz-Risiko** **niedrig** — der
Ist-Zustand ist **nachgewiesen** (der rote Lauf am unveränderten Runner) —,
**Reconciliation-Aufwand** null.
