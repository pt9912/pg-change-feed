# Slice slice-044: Hintergrund-Job für die Retention-Löschausführung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-13 — zweiter Slice, baut auf `slice-043`s
`RunRetentionUseCase` auf.

**Bezug:**
[`LH-FA-RET-002`](../../../../spec/lastenheft.md),
[`LH-FA-RET-003`](../../../../spec/lastenheft.md),
[`ADR-0014`](../../../../docs/plan/adr/0014-retention-domain-policy.md)
(nur umgesetzt).

**Berührte Spec-Stellen:** — (reine Verdrahtung eines bereits
bestehenden Use-Case in die Composition Root, keine neue
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

**Ziel:** Ein neuer Hintergrundzug `runRetentionCleanup`
(`internal/bootstrap/wiring.go`, Muster identisch zu `runHeartbeat`/
`runWALRetentionCheck`/`runAdministration`) ruft periodisch
`slice-043`s `RunRetentionUseCase` real gegen die laufende PostgreSQL-
Instanz auf und macht damit `RetentionPolicy.AllowsDeletion` erstmals im
laufenden Prozess wirksam. Ein realer End-zu-End-Beleg
(`tools/harness/run-integration-tests.sh`) zeigt: eine Change-Zeile, die
alle Freigabe-Bedingungen erfüllt, wird real entfernt; eine Zeile, die
das nicht tut (zu jung, oder ein Consumer hängt zurück), bleibt real
erhalten.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Sichtbarkeit blockierender Consumer** (`LH-FA-RET-005`) —
  Folge-Slice `slice-045`; ein anderer Liefer-Fokus.
- **`cdc_storage_bytes`-Metrik** (`LH-FA-RET-006`) — Folge-Slice
  `slice-046`.
- **Konfigurierbarkeit des Job-Takts zur Laufzeit** — Bestand bleibt
  bewusst stehen: ein fester, im Code deklarierter Takt (analog zu
  `heartbeatInterval`) genügt für diesen Slice; eine
  Laufzeit-Konfigurationsanbindung ist `welle-13` §6s ausdrücklicher
  Ausschluss.
- **Black-Box-E2E-Testabdeckung** über die im Compose-Stack real
  laufende Instanz hinaus — `welle-13` §6 verweist das bereits an die
  Folge-Welle „E2E-Abdeckung — Retention"; dieser Slice liefert nur den
  internen End-zu-End-Beleg im `run-integration-tests.sh`-Skript, keine
  eigenständige Black-Box-Testwelle.

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

- [ ] `runRetentionCleanup`-Hintergrundzug verdrahtet (`internal/bootstrap/wiring.go`,
      Muster identisch zu `runHeartbeat`/`runWALRetentionCheck`/
      `runAdministration`), ruft periodisch `RunRetentionUseCase` real
      auf.
- [ ] Realer End-zu-End-Beleg (`tools/harness/run-integration-tests.sh`):
      eine freigegebene Change-Zeile wird am laufenden Prozess real
      entfernt; eine nicht freigegebene (zu jung, oder Consumer hängt
      zurück) bleibt real erhalten — beides ohne Neustart.
- [ ] `make gates` grün, `make test-integration` grün (inkl. des neuen
      Belegs).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` (falls ein
      Betriebs-Aspekt entsteht, den ein Betreiber kennen muss — z. B.
      der Lösch-Takt) — Implementer prüft und begründet im Plan-Nachzug.
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
| `internal/bootstrap/wiring.go` | update | neuer Hintergrundzug `runRetentionCleanup` |
| `tools/harness/run-integration-tests.sh` | update | realer E2E-Beleg (Löschung erfolgt/unterbleibt korrekt) |
| `docs/user/benutzerhandbuch.md` | update, falls zutreffend | Betriebs-Aspekt des Lösch-Takts |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-043` liegt in `done/`,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass der E2E-Beleg eine größere Anpassung am Testskript braucht als
  erwartet (z. B. eine eigene, isolierte Feed-Container-Instanz wie bei
  Fehlerklassen-Tests), gehört das zurück zur Zerlegung.
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

- Ein zu kurzer Lösch-Takt könnte im Compose-Testlauf mit anderen
  Hintergrundzügen (Heartbeat, WAL-Retention-Check, Administration) um
  dieselbe Verbindung/Ressourcen konkurrieren. **Ausgang:** <bei Closure
  einzutragen>
- Der reale E2E-Beleg könnte eine präzise Zeitsteuerung brauchen
  (Change muss „alt genug" sein, `MinAge` real verstreichen lassen),
  was den Testlauf verlangsamt oder flaky machen könnte. **Ausgang:**
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
benannt nicht gezählt — dieser Slice liefert den zweiten Baustein der
Auflösung). Keiner der übrigen Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
