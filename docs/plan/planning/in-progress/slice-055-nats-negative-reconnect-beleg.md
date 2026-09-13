# Slice slice-055: NATS-Negative-Beleg — Reconnect-Nachholen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-15 — vierter Slice, unabhängig von `slice-054`, baut auf
der laufenden NATS-Verdrahtung aus `slice-053` auf.

**Bezug:** [LH-FA-SST-007](../../../../spec/lastenheft.md) (Negative-Beleg),
[ADR-0055](../../adr/0055-nats-change-notification-wecksignal.md)
(Wecksignal ohne Zustellgarantie, kein Nachliefern bei Core NATS — vorab
entschieden).

**Berührte Spec-Stellen:** [SPEC-017](../../../../spec/pflichtenheft.md)
(nur gelesen, nicht geändert).

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

**Ziel:** Ein realer Beleg zeigt: Nach einem simulierten NATS-
Verbindungsabbruch (der Testablauf trennt real die Verbindung zwischen
Test-Client und NATS-Server, nicht nur eine Attrappe) und anschließender
Wiederverbindung holt ein Consumer die während der Unterbrechung
aufgetretenen Changes **ausschließlich** über den bestehenden
SQL-Lesezugriffsweg `cdc.changes` nach — **nicht** über NATS selbst
(Core NATS liefert fire-and-forget nichts nach, `ADR-0055`). Das ist der
Gegenbeleg zu `slice-053`s Happy Path: Er zeigt, dass NATS als
Wecksignal-Kanal ausfallen darf, ohne dass die Erfassung oder die
Nachvollziehbarkeit der Changes beeinträchtigt wird.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein-Consumer-verbunden-Fall** — `slice-054`; dieser Slice prüft den
  Fall *verbunden gewesen, dann getrennt, dann wiederverbunden*, nicht
  *nie verbunden gewesen*.
- **Automatische NATS-Client-Reconnect-Logik im Produktionscode** — die
  `nats.go`-Bibliothek trägt bereits eingebautes Reconnect-Verhalten für
  den Publisher (`natsnotify`); dieser Slice belegt nur die
  Konsumenten-seitige Nachhol-Eigenschaft über SQL, ändert aber keine
  Publisher-seitige Reconnect-Konfiguration.
- **Neue Compose-/Wiring-Verdrahtung** — `slice-053` liefert sie bereits
  vollständig; dieser Slice nutzt die bestehende Verdrahtung nur als
  Testumgebung.

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

- [x] `LH-FA-SST-007` Negative-Beleg real erfüllt: real gegen den
      NATS-Server-Container die Verbindung des Test-Subscribers trennen
      (nicht bloß simulieren), während der Trennung eine oder mehrere
      Changes erzeugen, danach real wiederverbinden und belegen, dass die
      Changes ausschließlich über `cdc.changes` sichtbar werden — nicht
      über ein nachgeliefertes NATS-Signal — `make test-integration`.
- [x] Der Testablauf belegt real (Log-Beleg), dass während der Trennung
      **kein** Wecksignal beim Subscriber ankommt, sondern erst nach der
      Wiederverbindung überhaupt kein Signal für die verpassten Changes
      mehr eintrifft (Core NATS liefert nichts nach).
- [x] `make gates` grün, `make test-integration` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Kein Doku-Update nötig — kein neuer öffentlicher Vertrag, nur ein
      zusätzlicher Beleg für bestehendes Verhalten.
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
| `tools/harness/run-integration-tests.sh` | update | Negative-Beleg: reale Trennung/Wiederverbindung des Test-Subscribers, SQL-Nachhol-Beleg |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-058` liegt in `done/` (das
Subjekt-Schema muss auf dem tabellen-granularen Stand von `ADR-0056` sein,
bevor dieser Slice dagegen testet), `Verantwortlich:` gesetzt, WIP-Limit
(1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten bei reinem Testablauf-Beleg — falls doch, wäre das ein Zeichen,
  dass eine reale NATS-Verbindungstrennung im Compose-Netz komplexer zu
  erzwingen ist als angenommen (z. B. `docker network disconnect`
  gegenüber einem simplen Prozess-Stopp) und eine eigene Untersuchung
  braucht.
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

- Eine reale NATS-Verbindungstrennung im Compose-Netz zu erzwingen (statt
  sie nur zu simulieren) könnte technisch aufwändiger sein als ein reiner
  Prozess-Stopp des Test-Subscribers — z. B. `docker network disconnect`
  gegen das Compose-Netzwerk-Alias des Test-Containers, statt den
  NATS-Server selbst zu stoppen (ein gestoppter Server würde auch den
  Publisher treffen und den Testfall verfälschen). **Ausgang:** <bei
  Closure einzutragen>
- Der Testablauf könnte fälschlich einen erfolgreichen SQL-Nachhol-Beleg
  zeigen, obwohl in Wahrheit doch (versehentlich) ein NATS-Signal
  nachgeliefert wurde — das Log muss explizit belegen, dass beim
  Subscriber während der Trennung **kein** Frame ankam, nicht nur, dass
  die Changes am Ende über SQL sichtbar sind (sonst bliebe unklar, welcher
  Weg tatsächlich trug). **Ausgang:** <bei Closure einzutragen>

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
