# Slice slice-052: `ChangeNotificationPort` und `natsnotify`-Adapter

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-15 — erster Slice, `slice-053` baut auf ihm auf.

**Bezug:** [ADR-0055](../../adr/0055-nats-change-notification-wecksignal.md)
(Port-/Adapter-Design, Fehlerklasse, best-effort-Einreihung — vorab
entschieden).

**Berührte Spec-Stellen:** [SPEC-017](../../../../spec/pflichtenheft.md)
(Subjekt-/Nachrichtenform, nur gelesen, nicht geändert).

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

**Ziel:** Ein neuer Outbound-Port `ChangeNotificationPort`
(`internal/application/port/outbound/changenotification.go`, Methode
`Notify(ctx, sourceID) error`) und ein neuer Driven-Adapter
`internal/adapters/driven/natsnotify/` (Konstruktions-/Options-Muster
analog `postgresack`, Fehler-Wrapping über einen neuen Sentinel in die
Klasse `transient`) setzen `ADR-0055`s Port-/Adapter-Design um.
`CaptureService.Capture()` bekommt einen dritten, **optionalen**
Port-Parameter — ein Notify-Aufruf reiht sich NACH `ACK Source` ein, sein
Fehler wird an der Aufrufstelle abgefangen und niemals in den
Rückgabewert von `Capture()` übernommen (Regressionstest: ein
fehlschlagender `ChangeNotificationPort` darf den Erfolg von `store`/`ack`
nicht verdecken).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Reale NATS-Verbindung/Compose-Verdrahtung** — `slice-053`; dieser
  Slice liefert Port/Adapter/Verdrahtung in `CaptureService`, testet den
  Adapter aber gegen einen lokalen, dediziert für den Unit-Test
  gestarteten NATS-Server (oder eine Testdouble-Verbindung), nicht gegen
  die produktive Compose-Umgebung.
- **`CDC_NATS_URL`-Bootstrap-Verdrahtung** (`internal/bootstrap/wiring.go`)
  — `slice-053`; dieser Slice liefert den Adapter isoliert, seine
  Verdrahtung in die Composition Root ist ein eigener Liefer-Punkt in
  der Folge-Slice.
- **Boundary-/Negative-Belege** — `slice-054`/`055`; dieser Slice liefert
  nur das Grundgerüst (Happy-Path-Publish auf Adapter-Ebene), keine
  End-zu-Ende-Belege gegen einen laufenden Compose-Stack.

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

- [ ] `ChangeNotificationPort` real definiert, `natsnotify`-Adapter
      implementiert `Notify` gegen einen echten NATS-Server (Unit-/
      Adapter-Test, kein reiner Mock) — Fehlerklasse `transient` real
      über einen Sentinel geprüft.
- [ ] `CaptureService.Capture()` ruft den Port real als dritten,
      optionalen Schritt NACH `ACK Source` auf — Regressionstest belegt
      real: ein fehlschlagender `ChangeNotificationPort` lässt
      `Capture()` trotzdem erfolgreich zurückkehren, wenn `store`/`ack`
      erfolgreich waren (Rot-Beleg: Test schlägt fehl, wenn die
      Fehlerpropagation versehentlich eingeführt wird — real
      demonstriert durch temporäres Entfernen des Error-Swallowing).
- [x] `make gates` grün, `make test` grün. Verifier hat beide real
      erneut ausgeführt (`docs/reviews/verify-slice-052.md`).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: **entfällt** — kein Betriebs-Vertrag entsteht, bevor
      `slice-053` `CDC_NATS_URL` verdrahtet (Verifier bestätigt: kein
      `internal/bootstrap`/`compose.yaml`-Diff in diesem Slice).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. **Entfällt** — Repo ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine `reconciliation.md` vorhanden.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). **Verschoben auf `welle-15`-Closure** (dieser Slice trägt `Welle: welle-15`).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/outbound/changenotification.go` | neu | `ChangeNotificationPort` |
| `internal/adapters/driven/natsnotify/` | neu | Adapter, `ErrNotify`-Sentinel |
| `internal/application/usecase/capture/service.go` | update | dritter, optionaler Port-Parameter, best-effort nach ACK |
| `go.mod`/`go.sum` | update | `github.com/nats-io/nats.go` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-15` eröffnet,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass Port+Adapter+`CaptureService`-Erweiterung zusammen mehr als drei
  Liefer-Punkte oder mehr als zwei Schichten berühren (z. B. weil die
  `CaptureService`-Signaturänderung alle Aufrufer in mehreren Paketen
  gleichzeitig anfassen muss), gehört das zurück zur Zerlegung — Architect-
  Verdikt bereits vorab auf genau dieses Risiko hingewiesen.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test` grün **und**
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Adapter-Test gegen einen „echten NATS-Server" braucht eine
  netzlose, Docker-only-taugliche Testumgebung (analog zu
  `make test-store`s Testcontainer-Muster) — ohne das würde der Test
  entweder gegen einen Mock laufen (schwächerer Beleg) oder Netzzugriff
  brauchen (verboten, `AGENTS.md` §3.1). **Ausgang: entfallen** —
  `tools/harness/run-notify-tests.sh` + `make test-notify` lösen es real
  (Verifier hat den Lauf selbst reproduziert, `docs/reviews/verify-slice-052.md`).
- Die `CaptureService.Capture()`-Signaturänderung (dritter Parameter)
  könnte mehr Aufrufer/Tests berühren als erwartet und den Slice über
  die Drei-Liefer-Punkte-Grenze heben. **Ausgang: entfallen** —
  Functional-Option-Muster (`WithChangeNotification`) hielt alle fünf
  bestehenden Aufrufstellen unverändert kompilierbar, kein Aufrufer
  musste angefasst werden.
- `github.com/nats-io/nats.go` als erste Nicht-PostgreSQL-Abhängigkeit
  könnte unerwartete transitive Abhängigkeiten in `go.sum` einführen, die
  gegen die bisher schlanke Abhängigkeitsfläche des Repos abgewogen
  werden müssen. **Ausgang: entfallen** — nur fünf kleine, erwartbare
  `// indirect`-Einträge (`klauspost/compress`, `nats-io/nkeys`,
  `nats-io/nuid`, `golang.org/x/crypto`, `golang.org/x/sys`), vom
  Reviewer und Verifier unabhängig als unproblematisch bewertet.

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

- **Was hat funktioniert:** Das Functional-Option-Muster
  (`WithChangeNotification`, analog `postgresack.WithLog`) hielt die
  `CaptureService`-Erweiterung additiv — alle fünf bestehenden
  `NewCaptureService`-Aufrufstellen kompilierten unverändert, kein
  Rückführungsrisiko trat ein. Der Referenz-Adapter `postgresack.go` als
  Kopiervorlage für `natsnotify` trug die Konstruktions-/Options-/
  Fehler-Wrapping-Form vollständig. Der Rot-Beleg-Nachweis (Notify-Fehler
  darf `Capture()` nicht scheitern lassen) wurde real erbracht und vom
  Verifier unabhängig reproduziert (Error-Swallowing temporär entfernt,
  Test schlug real fehl, danach zurückgesetzt).
- **Was ging anders als geplant:** Der Godoc-Kommentar über `natsnotify.New`
  trug versehentlich eine Slice-Chronik ("Folge-Slice `slice-053`") statt
  eines ADR-Bezugs — das vierte Auftreten der bereits 3×-verkörperten
  Klasse `BEO-PGC/slice-chronik-in-code-kommentar`, trotz der seit der
  3×-Verkörperung geltenden Implementer-Selbstprüf-Instruktion
  (`.claude/commands/implement-slice.md` Schritt 20). Fixrunde behob den
  Fund; ein vorgezogener Architect-Zug (siehe Steering-Loop-Eintrag)
  bewertete das vierte Auftreten.
- **Steering-Loop-Eintrag:** `.harness/skills/reviewer.md` geschärft: ein
  eigener, benannter HIGH-Unterpunkt „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar" ersetzt die bisher implizite Subsumtion
  unter „Kommentar trägt keine der Kommentar-Klassen"; zusätzlich trägt
  `.claude/commands/implement-slice.md` Schritt 20 seither eine
  Grenz-Klarstellung (Implementer-Selbstprüfung ist erste, nicht
  tragende Verteidigungslinie — der unabhängige Reviewer bleibt die
  tragende Instanz, Modul 8 §Kernidee)
  — liegt in `.harness/skills/reviewer.md` (HIGH-Liste) und
  `.claude/commands/implement-slice.md` (Schritt 20).
  Auslöser: `BEO-PGC/slice-chronik-in-code-kommentar` (bereits 3×
  verkörpert; dieser Slice liefert den vorgezogen bewerteten 4. Beleg,
  siehe Architect-Verdikt
  `docs/reviews/architect-verdict-slice-chronik-in-code-kommentar-4x.md`
  — Status quo bestätigt, kein neuer mechanischer Sensor, aber die
  Rollenteilung Implementer-Selbstprüfung/Reviewer-Sicherheitsnetz
  explizit verkörpert).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-052.md`
  in `BEO-PGC/slice-chronik-in-code-kommentar/` ergänzt — 4. Beleg (Klasse
  war bereits bei 3× verkörpert; dieses Auftreten löste den oben
  genannten vorgezogenen Architect-Zug aus, unabhängig vom normalen
  Lese-Schritt bei Welle-Closure).
- **Folge-Slices:** `slice-053` (Compose-Verdrahtung und
  Happy-Path-Beleg) — liegt als Datei in `open/`.
- **Risiken aus §6:** alle drei mit Ausgang *entfallen* — siehe §6.
- **Drei Paarungen:** verschoben auf `welle-15`-Closure (dieser Slice
  trägt `Welle: welle-15`, siehe DoD-Item).

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
