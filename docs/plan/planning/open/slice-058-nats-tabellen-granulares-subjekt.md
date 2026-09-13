# Slice slice-058: NATS-Subjekt auf Tabellen-Granularität heben (ADR-0056)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-15 — löst die während `slice-053`s Laufzeit entstandene
`ADR-0056`-Folgepflicht ein; muss vor `slice-054`/`055` laufen (deren
Belege gegen das *korrekte*, tabellen-granulare Subjekt geführt werden
müssen).

**Bezug:** [LH-FA-SST-007](../../../../spec/lastenheft.md),
[ADR-0056](../../adr/0056-nats-tabellen-granulares-subjekt.md)
(Supersedes `ADR-0055` Punkt 2 — bindend: Subjekt-Schema, Notify-Kardinalität,
Port-/Modell-Erweiterung — vorab entschieden), [ADR-0055](../../adr/0055-nats-change-notification-wecksignal.md)
(alle übrigen Punkte weiterhin gültig).

**Berührte Spec-Stellen:** [SPEC-017](../../../../spec/pflichtenheft.md)
(Subjekt-Schema-Zeile, bereits durch `ADR-0056` aktualisiert, nur gelesen).

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

**Ziel:** `ChangeNotificationPort.Notify` bekommt die Signatur
`Notify(ctx context.Context, sourceID, schema, table string) error`
(`ADR-0056`); `natsnotify` baut das Subjekt als
`cdc.changes.<source_id>.<schema>.<table>`. `model.Change` (bzw. eine
begleitende Struktur) trägt Schema/Tabellenname zusätzlich zur
`SourceTableID`, durchgereicht aus dem bereits vorhandenen Wissen des
Driving-Adapter-Mappers — **kein** neuer Laufzeit-Lookup in
`CaptureService`. `CaptureService.Capture()` sammelt die distinkten
`(schema, table)`-Paare einer Transaktion und ruft `Notify` genau einmal
je Paar auf (Deduplizierung). Eine defensive Validierung lehnt Schema-/
Tabellennamen mit NATS-reservierten Zeichen (`.`, `*`, `>`) oder
Whitespace vor dem ersten Notify-Versuch ab. `tools/harness/natssub` und
der Happy-Path-Testabschnitt in `run-integration-tests.sh` (aus
`slice-053`) werden auf das neue vier-Ebenen-Subjekt nachgezogen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Neue Boundary-/Negative-Belege** — `slice-054`/`055`; dieser Slice
  liefert nur die korrigierte Grundlage (Port/Adapter/Modell/Kardinalität),
  keine neuen End-zu-Ende-Testfälle über den bereits bestehenden
  Happy-Path-Nachzug hinaus.
- **Zeilen-/Operations-Prädikat-Filterung** (`INSERT`/`UPDATE`/`DELETE`
  als weitere Subjekt-Ebene) — `ADR-0056`s Re-Evaluierungs-Trigger
  benennt das explizit als eigene, künftige Folge-ADR-Frage, keine dieser
  ADR/dieses Slice.
- **Rückwirkende Anpassung von `slice-052`/`053`s Closure-Notizen** —
  `git` hält die Historie; ihre §7-Einträge beschreiben korrekt, was zum
  jeweiligen Zeitpunkt galt (dieselbe Disziplin wie bei Slice-Chronik:
  Historie wird nicht umgeschrieben).

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

- [ ] `ChangeNotificationPort.Notify(ctx, sourceID, schema, table) error`
      real umgesetzt, `natsnotify` publiziert real auf
      `cdc.changes.<source_id>.<schema>.<table>` — Regressionstest gegen
      eine Rückkehr zur Drei-Token-Form (`make test`/`make test-notify`).
- [ ] `model.Change`/Assembler tragen Schema/Tabellenname zusätzlich zur
      `SourceTableID`, ohne neuen Laufzeit-Lookup in `CaptureService` —
      real durch Codeinspektion belegt (keine neue Outbound-Abhängigkeit).
- [ ] `CaptureService.Capture()` dedupliziert real: mehrere Changes
      derselben Tabelle in einer Transaktion lösen genau **ein** Notify
      aus; Changes über zwei Tabellen lösen zwei distinkte Aufrufe aus —
      Regressionstest (`make test`).
- [ ] Defensive Validierung gegen NATS-reservierte Zeichen (`.`, `*`, `>`)
      und Whitespace in Schema-/Tabellennamen vor dem ersten Notify-Versuch
      real getestet (mindestens ein Negativ-Fall).
- [ ] `tools/harness/natssub` und der Happy-Path-Testabschnitt in
      `run-integration-tests.sh` (`slice-053`) real auf das neue
      vier-Ebenen-Subjekt nachgezogen — `make test-integration` grün,
      realer Empfangsbeleg mit dem neuen Subjekt-Format.
- [ ] `make gates` grün, `make test` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: keiner erwartet (`SPEC-017`/`ARC-013` bereits durch
      `ADR-0056` aktualisiert; `docs/user/benutzerhandbuch.md` nennt kein
      Subjekt-Format) — Implementer prüft und begründet im Plan-Nachzug.
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
| `internal/application/port/outbound/changenotification.go` | update | `Notify`-Signatur um `schema`, `table` erweitert |
| `internal/adapters/driven/natsnotify/notify.go` | update | Subjekt-Bildung auf vier Tokens, Validierung reservierter Zeichen |
| `internal/domain/model/change.go` | update | Schema/Tabellenname zusätzlich zur `SourceTableID` |
| `internal/adapters/driving/replication/mapper/mapper.go` | update | Schema/Tabellenname an `model.NewChange`/Assembler-Konstruktion durchreichen |
| `internal/application/usecase/capture/service.go` | update | Deduplizierung je `(schema, table)`-Paar pro Transaktion |
| `tools/harness/natssub/main.go` | update | Subjekt-Parameter auf vier-Ebenen-Form |
| `tools/harness/run-integration-tests.sh` | update | Happy-Path-Subjekt aus `slice-053` nachgezogen |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-053` liegt in `done/`,
`ADR-0056` liegt vor (Accepted), `Verantwortlich:` gesetzt, WIP-Limit
(1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass die `model.Change`-Erweiterung mehr Aufrufer berührt als die in
  §3 gelistete Datei-Menge (z. B. weitere Adapter, die `model.Change`
  direkt konstruieren, außerhalb des Driving-Adapter-Mappers), gehört
  das zurück zur Zerlegung.
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

- Die `model.Change`-Erweiterung um Schema/Tabellenname könnte mehr
  Konstruktionsstellen berühren als nur den Driving-Adapter-Mapper
  (`ADR-0056` nennt ihn als einzige bekannte Stelle, aber das ist eine
  Annahme, keine vollständige Inventur). **Ausgang:** <bei Closure
  einzutragen>
- Die defensive Validierung gegen NATS-reservierte Zeichen könnte
  Aktivierungen ablehnen, die heute (ohne NATS) unauffällig funktionieren
  — ein Bestandsschema/-tabellenname mit einem ungewöhnlichen Zeichen
  würde durch dieses Slice zum ersten Mal sichtbar. **Ausgang:** <bei
  Closure einzutragen>
- `slice-054`/`055` (noch in `open/`) referenzieren in ihrem aktuellen
  Text noch das alte, drei-Token-Subjekt-Schema — sie brauchen einen
  Plan-Nachzug, bevor sie aktiviert werden, sonst driftet ihr Text vom
  tatsächlichen Namensstand. **Ausgang:** <bei Closure einzutragen>

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
