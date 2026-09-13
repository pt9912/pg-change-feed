# Slice slice-054: NATS-Boundary-Beleg — nicht verbundener Consumer

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-15 — dritter Slice, unabhängig von `slice-055`, baut auf
der laufenden NATS-Verdrahtung aus `slice-053` auf.

**Bezug:** [LH-FA-SST-007](../../../../spec/lastenheft.md) (Boundary),
[ADR-0055](../../adr/0055-nats-change-notification-wecksignal.md)
(leerer Payload, best-effort Notify — vorab entschieden).

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

**Ziel:** Ein realer Beleg zeigt: Eine Change, die auftritt, während **kein**
Consumer gegen `cdc.changes.<source_id>.<schema>.<table>` abonniert ist
(tabellen-granulares Subjekt, [ADR-0056](../../adr/0056-nats-tabellen-granulares-subjekt.md)),
geht nicht verloren — sie bleibt über den bestehenden SQL-Lesezugriffsweg
`cdc.changes`
vollständig lesbar. Das belegt strukturell, dass NATS ausschließlich ein
Wecksignal ist und niemals zur Wahrheitsquelle wird: Ein fehlendes oder zu
spätes Wecksignal darf die Erfassung selbst nicht beeinflussen
(`ADR-0055`, `CaptureService.Capture()`s Notify-Fehlerschluckung).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Reconnect-Verhalten eines zuvor verbundenen, dann getrennten
  Consumers** — `slice-055`; dieser Slice prüft den Fall *nie verbunden
  gewesen*, nicht *Verbindung verloren und wiederhergestellt*.
- **Änderung an `ChangeNotificationPort`/`natsnotify`** — `slice-052`
  liefert die Fehlerschluckung bereits vollständig (Regressionstest); dieser
  Slice belegt sie nur zusätzlich end-to-end gegen einen echten NATS-Server.
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

- [x] `LH-FA-SST-007` Boundary real erfüllt: **ohne** einen abonnierten
      Test-Client eine Change erzeugen, danach real gegen `cdc.changes`
      belegen, dass sie vollständig vorhanden ist — `make test-integration`.
- [x] Derselbe Lauf belegt real, dass `CaptureService.Capture()` in diesem
      Fall **nicht** blockiert oder fehlschlägt (Notify hat keinen
      Empfänger, aber das ist kein Fehler — `ADR-0055`).
- [x] `make gates` grün, `make test-integration` grün. Verifier hat beide
      real erneut ausgeführt, Exit-Code separat geprüft.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Kein Doku-Update nötig — kein neuer öffentlicher Vertrag, nur ein
      zusätzlicher Beleg für bestehendes Verhalten.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. **Entfällt** — Repo ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine `reconciliation.md` vorhanden.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. **Keine Beobachtung angefallen** — siehe §7 (Pipe-Exit-Code-Fallstrick unter der Zähl-Schwelle, kein formaler Vorgang).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). **Verschoben auf `welle-15`-Closure** (dieser Slice trägt `Welle: welle-15`).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` | update | Boundary-Beleg: Change ohne abonnierten Client, danach `cdc.changes`-Lesebeleg |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-058` liegt in `done/` (das
Subjekt-Schema muss auf dem tabellen-granularen Stand von `ADR-0056` sein,
bevor dieser Slice dagegen testet), `Verantwortlich:` gesetzt, WIP-Limit
(1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten — reiner Testablauf-Beleg ohne Produktionscode-Änderung.
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

- Der Testablauf könnte versehentlich doch einen (Rest-)Subscriber aus
  einem vorherigen Testschritt aktiv lassen und damit den Boundary-Fall
  gar nicht real herstellen — der Testablauf muss explizit belegen, dass
  zum Zeitpunkt der Change kein Client verbunden ist (z. B. über eine
  reale `nats server list-connections`/Server-Statistik-Abfrage statt
  einer bloßen Annahme). **Ausgang: entfallen** — reale Abfrage gegen den
  NATS-Server-Monitor-Endpunkt `/subsz` vor UND nach der Change belegt
  real das Fehlen des Subjekts in den aktiven Subscriptions; Reviewer
  und Verifier haben den Mechanismus unabhängig voneinander gegen einen
  frischen, gleich gepinnten Testcontainer reproduziert.
- `CaptureService.Capture()`s Notify-Aufruf könnte bei fehlendem Subscriber
  einen NATS-Client-seitigen Fehler zurückgeben, der versehentlich doch
  propagiert wird (Regressions-Risiko trotz `slice-052`s Unit-Test) —
  muss end-to-end, nicht nur unit-testseitig, real widerlegt werden.
  **Ausgang: entfallen** — realer `make test-integration`-Lauf zeigt: die
  Change bleibt über `cdc.changes` vollständig lesbar, der Feed-Container
  läuft danach unverändert weiter (`feed_running=true`).

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

- **Was hat funktioniert:** Der NATS-Server-Monitor-Endpunkt `/subsz`
  (Port 8222, bereits aus `slice-053`s Healthcheck bekannt) erwies sich
  als robuster, mechanismus-basierter Beleg für „kein Subscriber
  verbunden" — kein Rückgriff auf eine bloße Testablauf-Annahme nötig.
  Das bestehende Muster aus dem Happy-Path-Beleg (`natssub`,
  `docker logs`-Polling) diente als Stilvorbild für den neuen Abschnitt,
  ohne selbst gebraucht zu werden (hier läuft bewusst kein Subscriber).
- **Was ging anders als geplant:** Während der Verifikation trat real
  ein Pipe-Exit-Code-Fallstrick auf: ein `make gates | tail` maskiert
  einen Fehlschlag von `make` durch den Exit-Code von `tail` (0) — das
  führte zu einem real gepushten Commit mit einer nicht verlinkten
  Kennung im eigenen Reviewer-Report, bevor es bemerkt und sofort
  korrigiert wurde. Kein Produktionscode-Fehler, aber ein Planungs-/
  Verifikationsablauf-Fallstrick, der bereits einmal zuvor (beim
  `slice-058`-Implementer) auftrat.
- **Steering-Loop-Eintrag:** keiner — zwei Vorkommen (`slice-058`-
  Implementer, diese Slice-Closure) sind unter der 3×-Schwelle für einen
  Architect-Zug; benannt für den Fall eines dritten Auftretens.
- **Beobachtungs-Register (`../observations/`):** keine neue Beobachtung
  angelegt — der Pipe-Exit-Code-Fallstrick ist (noch) kein formal
  gezählter Vorgang (kein abgeschlossener Slice/Review trägt ihn als
  Fund), sondern eine im Fließtext benannte Beobachtung; bei einem
  dritten Auftreten wird ein Register-Eintrag fällig.
- **Folge-Slices:** keine.
- **Risiken aus §6:** beide mit Ausgang *entfallen* — siehe §6.
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
