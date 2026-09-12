# Slice slice-025: Metrik `cdc_wal_retention_bytes`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-7`](../welle-7.md) — der Ende-zu-Ende-Beleg (welle-7 §3)
verbindet diese Metrik erst mit der Schwellen-Logik aus `slice-026`; für
sich allein ist die Metrik nur die Hälfte des Zielbilds.

**Bezug:** [`SPEC-009`](../../../../spec/pflichtenheft.md) definiert
`cdc_wal_retention_bytes` bereits als Metrik-Zeile („WAL-Rückstand/
Replication-Slot-Zustand", Quelle „Replication Stream") — dieser Slice liefert
die bisher fehlende Implementierung. Kein aktives ADR wird geändert; das ADR
aus `slice-024` legt die Schwellenwerte fest, die diese Metrik *misst*,
entscheidet aber nicht über die Metrik-Erhebung selbst.

**Berührte Spec-Stellen:** [`SPEC-009`](../../../../spec/pflichtenheft.md)
(Metrik-Zeile `cdc_wal_retention_bytes`) — der Slice erfüllt eine bestehende
Zeile, ändert ihren Text nicht.
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

**Ziel:** Ein periodischer Health-Check (unabhängig vom reaktiven
Fehlerpfad, siehe welle-7 §1) liest — mit der bereits bestehenden
`cdc_capture`-Rollen-Verbindung, derselben, die
[`receive.go`](../../../../internal/adapters/driving/replication/receive/receive.go)
heute schon einmalig für den `START_REPLICATION`-Lookup nutzt — periodisch aus
`pg_replication_slots` den WAL-Rückstand des Capture-Slots (Differenz
zwischen aktuellem WAL-Schreibstand und `confirmed_flush_lsn`) und exponiert
ihn als Metrik `cdc_wal_retention_bytes` in Bytes, real getestet gegen eine
PostgreSQL mit künstlich erzeugtem Rückstand (Slot bewusst inaktiv gehalten,
während weitergeschrieben wird). Diese Metrik läuft **bewusst nicht** über
die bestehende `cdc.metrics`-SQL-View
([`tools/schema/nacharbeit-observability.sql`](../../../../tools/schema/nacharbeit-observability.sql)):
deren eigener Kommentar markiert `cdc_wal_retention_bytes` explizit als
„nicht abgedeckt", weil die Erhebung Systemkatalog-Zugriffe außerhalb des
`cdc`-Schemas braucht (`pg_stat_replication`/`pg_replication_slots`), die die
Least-Privilege-Fläche von `cdc_reader` unnötig erweitern würden — dieselbe
Begründung, mit der der Health-Endpoint dort auf den Heartbeat-Mechanismus
statt auf eine reine Lese-View ausweicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Schwellen-Vergleich und Reaktion (Fortsetzen/Abbrechen)** —
  Folge-Slice `slice-026` übernimmt das; dieser Slice liefert nur den
  gemessenen Wert, keine Entscheidung darüber.
- **Erweiterung von `cdc.metrics`/`cdc_reader`-Rechten um
  `pg_replication_slots`** — Bestand bleibt bewusst stehen: Genau diese
  Erweiterung der Least-Privilege-Fläche zu vermeiden ist der dokumentierte
  Grund, warum `cdc_wal_retention_bytes` prozessintern bleibt statt in die
  SQL-View zu wandern (siehe Ziel oben).
- **Änderung der bestehenden `START_REPLICATION`-Lookup-Nutzung von
  `confirmed_flush_lsn`** in `receive.go` — Bestand bleibt bewusst stehen:
  Der neue Health-Check liest denselben Katalog zusätzlich und periodisch,
  verändert die bestehende Einmal-Lookup-Nutzung beim Verbindungsaufbau
  nicht.

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

- [x] `SPEC-009` (`cdc_wal_retention_bytes`) erfüllt: `receive.WALRetentionChecker`
      (`internal/adapters/driving/replication/receive/walretention.go`) misst
      periodisch gegen `pg_replication_slots`/`IDENTIFY_SYSTEM` und liefert
      einen Bytes-Wert, real getestet
      (`TestWALRetentionMeasuresGrowingBytes`,
      `internal/adapters/driving/replication/receive/stream_test.go`) gegen
      PostgreSQL mit künstlich erzeugtem WAL-Rückstand (Slot bewusst inaktiv
      gehalten, Wert steigt real messbar) — `make test-replication` dreimal in
      Folge grün; rot färbende Mutation (LSN-Subtraktion vertauscht) einmal
      gesehen (rot, Bericht Schritt 8).
- [x] Messintervall an den bestehenden Heartbeat-Takt gebunden
      (`heartbeatInterval`, `internal/bootstrap/wiring.go`) — keine neue
      Konfigurationsachse.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für `docs/user/benutzerhandbuch.md` (§4 „WAL-Rückstand
      prüfen", neu; §9 Grenzwerte und Änderungshistorie 1.3) — die Metrik
      wird erstmals real erhoben, exponiert über strukturiertes Log statt
      `cdc.metrics` (Begründung §3).
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
| `internal/adapters/driving/replication/receive/walretention.go` (neu, Paket `receive`) | neu | `WALRetentionChecker`: eigene, von der Stream-Verbindung getrennte `replication=database`-Verbindung derselben Rolle `cdc_capture` (die Stream-Verbindung steht während `Run` im COPY-Modus und nimmt keine Abfragen mehr entgegen); `Measure` bildet die Differenz aus `IDENTIFY_SYSTEM` (aktuelle WAL-Position) und `confirmed_flush_lsn` (`pg_replication_slots`) — beide LSN-Werte sind derselbe 64-Bit-Byte-Offset, die Subtraktion in Go erspart den sonst nötigen Aufruf der auf `pg_monitor` beschränkten Funktion `pg_current_wal_lsn()` |
| `internal/bootstrap/wiring.go` | update | **Plan-Nachzug:** periodischer Health-Check-Zug `runWALRetentionCheck` — eigene Goroutine, eigener `WALRetentionChecker`, gebunden an `heartbeatInterval` (keine neue Konfigurationsachse), analog zu `runHeartbeat` |
| Expositionsweg: **strukturiertes Log** (`outbound.LogPort`, `metric=cdc_wal_retention_bytes`) — **kein** `cdc.metrics`-View-Eintrag, siehe §1 | neu | **Plan-Nachzug — Entscheidung nachgetragen:** von den beiden im Slice-Kopf skizzierten Optionen (Tabellenspalte analog `cdc.process_heartbeat`, oder Log) fiel die Wahl auf Log, nicht auf eine Tabellenspalte: `slice-026` liest den Messwert direkt im selben Prozess (In-Memory-Rückgabewert von `Measure`) für die Schwellen-Entscheidung — ein persistierter Zwischenstand hätte keinen weiteren Leser, bräuchte aber einen Schema-Rollout (`tools/schema/schema.yaml` + d-migrate) und eine zusätzliche Schreibrolle/-verbindung; das Log macht den Wert für einen Betreiber lesbar, ohne die `cdc_reader`-Rechtefläche zu erweitern und ohne neues Schema |
| `internal/adapters/driving/replication/receive/stream_test.go` | update | **Plan-Nachzug (Ort statt neuer Datei):** `TestWALRetentionMeasuresGrowingBytes` realer Test gegen PostgreSQL-Testcontainer mit künstlich erzeugtem Rückstand (Slot inaktiv nach einer bestätigten Transaktion, weitere unbestätigte Inserts lassen die Messung real wachsen) — im bestehenden `stream_test.go` statt einer neuen Datei, weil die Test-Infrastruktur (`newTestEnv`, `fakeCapture`, `readConfirmedFlush`) dort bereits liegt |
| `docs/user/benutzerhandbuch.md` | update | **Plan-Nachzug:** neuer Abschnitt „WAL-Rückstand prüfen" (§4), Abgrenzung gegen `cdc.metrics` (§4 Metriken lesen), Grenzwerte-Hinweis zu `max_wal_senders` (§9), Änderungshistorie 1.3 |

**Plan-Nachzug (Fixrunde nach Review, [`review-slice-025.md`](../../../reviews/review-slice-025.md) F-1/F-2 behoben):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/replication/receive/walretention.go` | update | F-2: `WALRetentionChecker` trägt jetzt die eigene `dsn`; `Measure` ersetzt die Verbindung über `reconnectAfterError` bei einem Protokoll-/Katalogfehler (`IDENTIFY_SYSTEM`, Katalogabfrage) — der nächste Aufruf misst wieder, statt die Metrik dauerhaft verstummen zu lassen |
| `internal/bootstrap/wiring.go` | update | Kommentar-Nachzug an `runWALRetentionCheck`: die Schleife muss den Reconnect nicht selbst behandeln, das übernimmt `Measure` intern — keine Verhaltensänderung der Schleife selbst |
| `internal/adapters/driving/replication/receive/stream_test.go` | update | F-1: drei neue reale Fehlerpfad-Tests (`TestWALRetentionMeasureInvalidSlotName`, `TestWALRetentionMeasureConnectionRefusedFails`, `TestWALRetentionMeasureMissingSlot`) für ungültigen/fehlenden Slot und Verbindungsfehler; F-2: `TestWALRetentionMeasureReconnectsAfterConnectionLoss` terminiert die Checker-Verbindung serverseitig (`pg_terminate_backend`, Backend über `application_name` identifiziert) und belegt real, dass der übernächste `Measure`-Aufruf nach dem sichtbar gescheiterten wieder erfolgreich misst |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-024` liegt in `done/` (ADR
`Accepted`, legt die Sentinel-Trennung und Schwellenwerte fest, die diese
Metrik später bedient) — Priorisiert, `Verantwortlich:` gesetzt, WIP-Limit
(1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  kein bestehender Metrik-Erhebungs-Mechanismus existiert und einer komplett
  neu gebaut werden müsste (mehr als eine Schicht: Erhebung + Export +
  Konfiguration), gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): `slice-024` ist noch nicht
  `done` (Sentinel-Trennung/Schwellenwerte fehlen als Referenzpunkt).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** der reale
Rückstand-Test (`make test-replication` oder Äquivalent) grün **und**
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein künstlich erzeugter WAL-Rückstand in einem Testcontainer-Lauf könnte
  durch PostgreSQLs eigene WAL-Rotation/Checkpoint-Verhalten schwer
  reproduzierbar oder flaky sein. **Ausgang:** <bei Closure einzutragen>
- Der konkrete prozessinterne Expositionsweg (eigene Tabelle analog
  `cdc.process_heartbeat`, reine Log-Zeile, oder ein neuer schlanker
  Lese-Pfad) ist hier noch nicht entschieden — nur, dass es **nicht** die
  `cdc.metrics`-View sein darf. **Ausgang:** <bei Closure einzutragen>

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
Treffer für `PGC`: `BEO-PGC/spec008-replication-luecke` (1×, weiter offen —
dieser Slice ist Teil der Antwort), `BEO-PGC/walsender-wirksamkeit` (nicht
einschlägig), `BEO-PGC/rollen-test-abdeckungsluecken` (nicht einschlägig).
Keiner der Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
