# Slice slice-045: Sichtbarkeit blockierender Consumer

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-13 — dritter Slice, unabhängig von `slice-043`/`044`
(reine Observability auf bestehendem `cdc.consumer_status`-Zustand).

**Bezug:** [`LH-FA-RET-005`](../../../../spec/lastenheft.md).

**Berührte Spec-Stellen:** — (reine SQL-View-/CLI-Ergänzung auf
bestehenden Daten, keine neue Architektur-Sicht-Aussage).

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

**Ziel:** Eine neue SQL-Lese-View (z. B. `cdc.retention_blockers`,
Implementer entscheidet den Namen und begründet ihn im Plan-Nachzug)
zeigt je Quelle den am weitesten zurückliegenden unbestätigten Consumer
— Name, Position, Rückstand zur letzten Commit-Position — also genau
den Consumer, dessen Position `RetentionPolicy.AllowsDeletion`
(`slice-043`) aktuell blockiert (`LH-FA-RET-004.a`s „Safe Watermark"
sichtbar gemacht, nicht neu berechnet). Der bestehende `diagnose`-
CLI-Befehl (`slice-038`) bekommt optional dieselbe Information
ergänzt, falls das ohne Bruch des bestehenden Ausgabeformats möglich
ist (Implementer-Entscheidung, siehe Plan-Nachzug).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Neue Berechnung des Safe Watermark** — `LH-FA-RET-004.a` beschreibt
  das Verfahren bereits vollständig, `slice-043`s
  `RunRetentionUseCase` implementiert es; dieser Slice macht das
  Ergebnis nur sichtbar, berechnet es nicht neu (keine Duplikat-Logik
  zwischen SQL-View und Go-Code).
- **Automatische Warnung/Alarmierung bei blockierenden Consumern** (z. B.
  über die `diagnose`-Fehlerzustand-Klasse) — ein anderer Vorgang mit
  eigenem Schwellenwert-Bedarf (`AGENTS.md` §3.6, ADR-pflichtig), hier
  nicht angefragt.
- **Löschausführung selbst** — `slice-043`/`044`; dieser Slice ist reine
  Observability auf bereits vorhandenem Zustand.

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

- [ ] Neue SQL-View zeigt je Quelle real den am weitesten
      zurückliegenden unbestätigten Consumer (Name, Position,
      Rückstand) — real gegen mindestens zwei Consumer getestet
      (einer blockiert, einer nicht).
- [ ] `LH-FA-RET-005` real erfüllt: ein Integrationstest liest die neue
      View über SQL (analog zum bestehenden `cdc.active_tables`-/
      `cdc.consumer_status`-Testmuster aus `welle-11`).
- [ ] `make gates` grün, `make test-integration` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` nennt die neue View
      (analog zu „Aktivierte Tabellen auflisten").
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
| `tools/schema/schema.yaml` | update | neue View im `views:`-Knoten (deklarativ, wie `active_tables`/`consumer_status`) |
| `test/integration/integration_test.go` | neu/update | Testfall gegen die neue View |
| `docs/user/benutzerhandbuch.md` | update | neue View dokumentiert |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-13` eröffnet,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei —
unabhängig von `slice-043`/`044`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten bei einer reinen Lese-View auf bereits vorhandenem
  Zustand — falls doch, wäre das ein Zeichen für eine unerwartet
  komplexe Watermark-Berechnung, die eigentlich `slice-043` gehört.
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

- Ein Consumer ohne jede bestätigte Position (nie bestätigt) könnte die
  View entweder fälschlich als „kein Blocker" oder mit einem
  irreführenden Wert zeigen — dieselbe Klasse Randfall wie
  `slice-038`s §6 Risiko 2. **Ausgang:** <bei Closure einzutragen>
- d-migrate könnte für die neue View denselben `POST_EXECUTE_DRIFT`-
  Bootstrap-Bedarf zeigen wie frühere Views vor ihrer deklarativen
  Überführung (`BEO-PGC/d-migrate-nacharbeit`) — mittlerweile aber laut
  `harness/README.md` seit `slice-016` behoben für Views. **Ausgang:**
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
Treffer für `PGC`: `BEO-PGC/retention-keine-loeschausfuehrung` (0×,
benannt nicht gezählt), `BEO-PGC/d-migrate-nacharbeit` (5×, bereits
verkörpert — für Views seit `slice-016` gelöst, siehe §6). Keiner
erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
