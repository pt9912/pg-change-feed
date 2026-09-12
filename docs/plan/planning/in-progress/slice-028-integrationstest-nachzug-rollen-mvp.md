# Slice slice-028: Compose-Integrationstest-Nachzug — Rollen-DSN-Trennung, MVP-Umbenennung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-8`](../welle-8.md) — der Nachweis, dass der
Compose-Integrationstest alle seit `welle-3` hinzugekommenen
post-MVP-Fähigkeiten zusammen abdeckt, ist welle-8s Closure-Trigger (§3),
kein Einzel-Slice-DoD.

**Bezug:** [`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…`003`
(Least-Privilege, getrennte Berechtigbarkeit, Beschränkbarkeit von
CDC-Datenzugriffen) — dieser Slice erweitert die Testabdeckung für einen
bestehenden Vertrag (`ADR-0047`, bereits `slice-023` umgesetzt), ändert ihn
nicht. Kein aktives ADR wird geändert.

**Berührte Spec-Stellen:** — (reine Testinfrastruktur und
Dokumentations-/Namenskorrektur, keine Verhaltensänderung an einer
Spec-Stelle).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

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

**Ziel:** Der Compose-Integrationstest verifiziert real, dass die seit
`slice-023` verdrahtete rollen-spezifische DSN-Trennung
(`cdc_capture`/`cdc_admin`/`cdc_reader`, `ADR-0047`) auch am laufenden,
containerisierten Gesamtsystem greift — insbesondere schließt er
`BEO-PGC/rollen-test-abdeckungsluecken` Punkt (2): Replication-Stream- und
ACK-Adapter (`cdc_capture`-gebunden) werden gegen Rollen-Vertauschung
testgesichert. Zusätzlich: `make test-integration`
([`tools/harness/run-integration-tests.sh`](../../../../tools/harness/run-integration-tests.sh))
und `test/integration/mvp_test.go` werden umbenannt, sobald ihr Scope die
Bezeichnung „MVP" nicht mehr korrekt trägt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`BEO-PGC/rollen-test-abdeckungsluecken` Punkt (1)** (Heartbeat-Grant-Test
  liest nie den tatsächlichen `nacharbeit-roles.sql`-Inhalt) — Bestand
  bleibt bewusst stehen: anderer Vorgang, betrifft einen bestehenden
  Unit-/Adapter-Test (`internal/bootstrap/roles_wiring_test.go`), keine
  neue Compose-Integrationstest-Abdeckung.
- **Neue Rollen oder DDL-Änderungen** — Bestand bleibt bewusst stehen: Die
  drei Rollen existieren bereits (`slice-011`, `slice-023`) und werden nur
  *geprüft*, nicht neu geschnitten.
- **Black-Box-CLI-Aufrufmuster neu erfinden** — Folge-Slice `slice-027`
  liefert das bereits; dieser Slice nutzt es, wo es für die
  Rollen-Verifikation passt (z. B. ein CLI-Aufruf mit absichtlich falscher
  Rolle), erfindet aber kein zweites Muster.
- **WAL-Rückstand-Schwellen-Verhalten zusätzlich gegen den Compose-Stack
  prüfen** — Bestand bleibt bewusst stehen: bereits real gegen PostgreSQL
  bewiesen (`slice-026`, welle-8 §6 Out-of-Scope).

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

- [ ] `LH-QA-SEC-001`…`003` erfüllt: Der Compose-Integrationstest belegt
      real, dass ein Zugriff mit der falschen Rolle (z. B. `cdc_reader` für
      einen schreibenden Aufruf, oder ein Replication-Stream-/ACK-Aufruf
      mit vertauschter Rolle) am laufenden, containerisierten System
      zurückgewiesen wird — schließt `BEO-PGC/rollen-test-abdeckungsluecken`
      Punkt (2).
- [ ] `mvp_test.go` benannt nach tatsächlichem Scope (der Make-**Target**-Name
      `test-integration` bleibt unverändert — er ist bereits scope-neutral,
      nur der Datei- und Helptext-Name trägt „MVP"), und der
      Makefile-Helptext (`## MVP-Integrationstest …`) sowie
      `harness/README.md`/`AGENTS.md`-Erwähnungen entsprechend nachgezogen,
      sobald der tatsächliche Scope (Consumer-Zugriffsweg,
      Rollen-DSN-Trennung, plus der bereits bestehende
      `cdc_capture_lag`-Lasttest-Beleg) die Bezeichnung „MVP" nicht mehr
      korrekt trägt.
- [ ] `make gates` grün, `make test-integration` dreimal in Folge grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für `harness/README.md`/`AGENTS.md` (Sensors-Tabelle,
      Target-Name und -Beschreibung nach der Umbenennung).
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
| `tools/harness/run-integration-tests.sh` | update | neuer Abschnitt: Zugriff mit falscher Rolle (z. B. `cdc_reader`-DSN für einen schreibenden Aufruf) real gegen den Compose-Stack zurückgewiesen |
| `tools/harness/run-replication-tests.sh` oder Äquivalent | update | rollenbeschränkte Login-Test-Identität für den Replication-Stream-/ACK-Rollen-Vertauschungstest (`BEO-PGC/rollen-test-abdeckungsluecken` Punkt 2) |
| `test/integration/mvp_test.go` → neuer Dateiname | rename+update | Umbenennung nach tatsächlichem Scope |
| `Makefile` | update | Helptext-Zeile für `test-integration` (Target-Name bleibt) |
| `harness/README.md`, `AGENTS.md` | update | Sensors-Tabelle/Erwähnungen nach der Umbenennung nachgezogen |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-027` liegt in `done/` (liefert
das Black-Box-CLI-Aufrufmuster, das dieser Slice für die
Rollen-Verifikation mitnutzt) — Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  die Rollen-Verifikation für Replication-Stream/ACK eine grundlegende
  Änderung an `run-replication-tests.sh`s Verbindungsaufbau braucht (mehr
  als eine Schicht), gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): `slice-027` ist noch
  nicht `done`.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
dreimal in Folge grün **und** Closure-Notiz geschrieben (inkl.
`BEO-PGC/rollen-test-abdeckungsluecken` Punkt-(2)-Fortschritt).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Rollenbeschränkte Login-Test-Identitäten für den Replication-Stream
  könnten zusätzliche PostgreSQL-Rollen-Berechtigungen brauchen, die
  `ADR-0047`s Rollenmodell nicht vorsieht — das wäre eine Architect-Frage,
  keine reine Testinfrastruktur-Änderung. **Ausgang:** <bei Closure
  einzutragen>
- Die Umbenennung von `mvp_test.go`/Helptext könnte weitere, hier nicht
  erfasste Erwähnungen von „MVP-Integrationstest" im Repo übersehen (z. B.
  in älteren Closure-Notizen — die bleiben unangetastet, da sie Historie
  sind, aber laufende Doku könnte betroffen sein). **Ausgang:** <bei
  Closure einzutragen>

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
Treffer für `PGC`: `BEO-PGC/rollen-test-abdeckungsluecken` (1×, weiter
offen — dieser Slice schließt Punkt (2) davon),
`BEO-PGC/test-isolation-geteilter-zustand` (1×, nicht einschlägig — anderer
Test-Layer, siehe welle-8 §6 Out-of-Scope). Keiner der Treffer erreicht mit
diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
