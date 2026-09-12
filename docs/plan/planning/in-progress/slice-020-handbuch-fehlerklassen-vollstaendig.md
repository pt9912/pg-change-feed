# Slice slice-020: Handbuch-Fehlerklassen vollständig

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die DoD dieses Slice (Handbuch-Tabelle vollständig,
`make gates` grün) ist bereits die volle Closure-Bedingung; es gibt kein
*Mehr*, das eine Welle beobachten müsste (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (Sichtbare
Fehlerzustände), [`ADR-0023`](../../adr/0023-fehlerklassifikation.md)
(Fehlerklassifikation, Accepted).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md) — die
Handbuch-Tabelle bildet die dortige geschlossene Menge der sieben
Fehlerklassen ab; der Slice ändert `SPEC-008` nicht, nur die Betreiberdoku.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/user/benutzerhandbuch.md` §6 *Fehlerklassen* listet alle
sieben stabilen Kategorien aus `ADR-0023`/`SPEC-008` (`transient`,
`configuration`, `permission`, `schema`, `storage`, `replication`,
`internal`), nicht nur die bisherigen vier.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code-Änderung an `classifyRunError`/`ErrorClass`** — Bestand bleibt
  bewusst stehen: der Fund betrifft ausschließlich die Handbuch-Tabelle, der
  Code implementiert die geschlossene Menge bereits korrekt
  (`internal/domain/model/errorstate.go`).
- **Produktion von `transient`/`permission` durch einen Adapter** — anderer
  Vorgang: Diese beiden Klassen sind Teil der deklarierten Menge, werden
  aber von keinem Adapter aktuell konstruiert; das Handbuch dokumentiert sie
  als Teil der Menge, ohne zu behaupten, dass sie heute beobachtbar sind.
- **Ein Sensor, der Handbuch-Inhalt gegen den Fehlerklassen-Code prüft** —
  anderer Vorgang: Wäre eine Werkzeug-Änderung, kein Doku-Fix; die
  Beobachtungslage dazu steht in §8.

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

- [x] `docs/user/benutzerhandbuch.md` §6 *Fehlerklassen*-Tabelle führt alle
      sieben Klassen aus `ADR-0023`/`SPEC-008`, mit Bedingung und Aktion je
      Zeile; `internal` ist als real erreichbarer Fallback
      (`classifyRunError`, `internal/bootstrap/wiring.go`) erkennbar
      benannt, `transient`/`permission` als deklariert, aber von keinem
      Adapter aktuell konstruiert.
- [x] Änderungshistorie (§9) des Handbuchs trägt einen Eintrag für diese
      Korrektur.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für `error_class`-Sichtbarkeit — bereits der Kern der
      Lieferung (DoD-Punkt 1); Item entfällt hier als Duplikat.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
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
| `docs/user/benutzerhandbuch.md` §6 | update | Fehlerklassen-Tabelle auf alle sieben Klassen aus `ADR-0023`/`SPEC-008` erweitern |
| `docs/user/benutzerhandbuch.md` §9 | update | Änderungshistorie-Zeile für diese Korrektur |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Falls sich beim
  Ausfüllen zeigt, dass `transient`/`permission` doch produziert werden
  (Code-Fund) — dann ist es kein reiner Doku-Fix mehr, sondern berührt
  Adapter-Code.
- `in-progress` → `open` (blockiert — Carveout?): Nicht erwartet; keine
  externe Abhängigkeit.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Handbuch-Tabelle führt alle sieben Fehlerklassen **und** `make gates` ist
grün.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- `classifyRunError` (`internal/bootstrap/wiring.go`) könnte sich zwischen
  dieser Planung und der Umsetzung ändern (neue Sentinel-Fehler, neue
  produzierte Klasse) — der Implementer verifiziert bei Umsetzung gegen den
  dann aktuellen Code, nicht gegen diese Planung. **Ausgang:** wird bei
  Closure eingetragen.

## 7. Closure-Notiz

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührte Sub-Area ist `*` (Default,
`PGC`, Greenfield laut `harness/conventions.md`) — die einzige deklarierte
Sub-Area dieses Repos, Schwelle also trivial erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register (`../observations/`)
durchgegangen — keine der 13 geführten Beobachtungen betrifft
Betreiberdoku-Drift oder die Fehlerklassen-Tabelle; keine Treffer.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (reine
Doku-Korrektur, keine neue Sub-Area-Berührung).
