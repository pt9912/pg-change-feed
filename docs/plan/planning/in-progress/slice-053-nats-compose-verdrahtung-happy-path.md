# Slice slice-053: NATS-Compose-Verdrahtung und Happy-Path-Beleg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-15 — zweiter Slice, baut auf `slice-052`s Port/Adapter
auf.

**Bezug:** [LH-FA-SST-007](../../../../spec/lastenheft.md) (Happy Path),
[ADR-0055](../../adr/0055-nats-change-notification-wecksignal.md)
(`CDC_NATS_URL`, Image-Pin — vorab entschieden).

**Berührte Spec-Stellen:** [SPEC-017](../../../../spec/pflichtenheft.md)
(Subjekt-Schema, nur gelesen, nicht geändert).

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

**Ziel:** `internal/bootstrap/wiring.go` verdrahtet `CDC_NATS_URL`
(vollständig optional — ungesetzt heißt Feature deaktiviert, kein
`ChangeNotificationPort` wird konstruiert). `compose.yaml` bekommt einen
neuen, digest-gepinnten NATS-Server-Service
(`nats:2-alpine@sha256:065e8355c20a5575b3c77224be1855e8103fd148b68fba05130b9b8ddfa40ccc`,
`ADR-0055`) und `CDC_NATS_URL` für den Feed-Container. Ein realer
`make test-integration`-Beleg zeigt: ein gegen `cdc.changes.<source_id>`
abonnierter Test-Client erhält real ein Wecksignal, sobald eine neue
Change erfasst wird (Happy Path, `LH-FA-SST-007`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Boundary-/Negative-Belege** — `slice-054`/`055`; dieser Slice liefert
  nur den Happy-Path-Beleg bei durchgehend verbundenem Consumer.
- **Neue `ChangeNotificationPort`-Logik** — `slice-052` liefert Port und
  Adapter bereits vollständig; dieser Slice verdrahtet sie nur in die
  Composition Root.
- **Consumer-seitiges NATS-Client-Werkzeug** — `welle-15` §6 schließt das
  bereits repo-weit aus (Betreiber-/Nutzer-Sache, kein Repo-Gegenstand).

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

- [x] `CDC_NATS_URL` real in `internal/bootstrap/wiring.go` verdrahtet —
      real getestet: gesetzt → `ChangeNotificationPort` konstruiert und real
      benachrichtigt (`make test-integration`, NATS-Happy-Path-Beleg);
      ungesetzt → kein Notify-Versuch, bestehendes Verhalten unverändert,
      real belegt über `make test-replication`
      (`TestWALRetentionThresholdEndToEnd` ruft `bootstrap.Run` mit einer
      `Config` ohne `NatsURL` auf, unverändert grün) (Boundary-Vorprüfung
      für `slice-054`, hier nur die Verdrahtung selbst).
- [x] `compose.yaml` trägt real einen digest-gepinnten NATS-Service samt
      `CDC_NATS_URL` für den Feed-Container.
- [x] `LH-FA-SST-007` Happy Path real erfüllt: ein Testclient, der
      `cdc.changes.<source_id>` abonniert, empfängt real ein Wecksignal
      nach einer neuen Change — `make test-integration`.
- [x] `make gates` grün, `make test-integration` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Report: `docs/reviews/review-slice-053.md` (0 HIGH, 1 MEDIUM ohne
      Fixrunde, 2 INFO).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` §„Umgebungsvariablen
      des Feed-Containers" um `CDC_NATS_URL` ergänzt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. **Entfällt** — Repo ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine `reconciliation.md` vorhanden.
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
| `internal/bootstrap/wiring.go` | update | `CDC_NATS_URL`-Verdrahtung, No-Op bei ungesetzt |
| `compose.yaml` | update | neuer NATS-Service, `CDC_NATS_URL`-Env für Feed-Container |
| `tools/harness/run-integration-tests.sh` | update | Happy-Path-Beleg (Test-Subscriber) |
| `tools/harness/natssub/main.go` | neu | Wegwerf-Testclient (Subscribe-vor-Change, `go run` im Toolchain-Container) — Plan-Nachzug: der Happy-Path-Beleg brauchte einen eigenständigen NATS-Testclient-Prozess, den `run-integration-tests.sh` allein (Shell/`psql`) nicht bereitstellen kann |
| `docs/user/benutzerhandbuch.md` | update | `CDC_NATS_URL` in der Env-Var-Tabelle |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-052` liegt in `done/`,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten bei reiner Verdrahtungs-/Compose-Arbeit auf bereits fertigem
  Port/Adapter — falls doch, wäre das ein Zeichen für unerwartete
  NATS-Compose-Startreihenfolge-Probleme, die eine eigene Untersuchung
  brauchen.
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

- Der Feed-Container könnte beim Compose-Start versuchen, sich mit NATS
  zu verbinden, bevor der NATS-Server bereit ist (dieselbe Klasse
  Startreihenfolge-Risiko wie bei PostgreSQL, dort über einen
  Healthcheck gelöst, `compose.yaml`). **Ausgang:** <bei Closure
  einzutragen>
- Ein Test-Subscriber, der erst NACH der Change abonniert, verpasst das
  Fire-and-Forget-Signal strukturell (kein Nachliefern bei Core NATS,
  `ADR-0055`) — der Testablauf muss die Subscription real vor dem
  Erfassungs-Ereignis aufbauen. **Ausgang:** <bei Closure einzutragen>

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
Keine Treffer für `PGC` zu NATS/Notification/Messaging.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
