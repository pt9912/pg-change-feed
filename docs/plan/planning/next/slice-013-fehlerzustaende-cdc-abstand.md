# Slice slice-013: Fehlerzustände sichtbar, CDC-Abstand messbar

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-4.

**Bezug:** [`LH-FA-ADM-003`](../../../../spec/lastenheft.md)/004, [`LH-QA-REL-003`](../../../../spec/lastenheft.md), [`LH-QA-OPS-003`](../../../../spec/lastenheft.md), [`ADR-0023`](../../../../docs/plan/adr/README.md), [`ADR-0046`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md), [`SPEC-013`](../../../../spec/pflichtenheft.md)
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-10.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Fehlerzustände des Capture-Prozesses sichtbar machen (Erfassung
kann nicht fortsetzen, erkennbar und vom Normalbetrieb unterscheidbar) und
den CDC-Abstand (Zeit/Position zwischen Quelländerung und
CDC-Verfügbarkeit) messbar machen — beides über die bestehende
SQL-Lese-Fläche.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Automatisierte Alarmierung/Benachrichtigung bei Fehlerzuständen —
  Sichtbarkeit ist das Ziel dieses Slice, aktive Alarmierung ist ein
  anderer Vorgang (Integration mit externem Monitoring, kein Bedarf
  beobachtet).
- Volle [`SPEC-013`](../../../../spec/pflichtenheft.md)-Latenzschwellen-
  Durchsetzung (Warn-/Fehler-Schwellen als aktive Prüfung) — Bestand
  bleibt bewusst stehen: dieser Slice liefert die Messung, nicht die
  Schwellen-Bewertung; die Schwellen selbst stehen bereits in
  [`SPEC-013`](../../../../spec/pflichtenheft.md).

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

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Fehlerzustands-Sichtbarkeit (Erfassung gestört/normal
      unterscheidbar) — Teil-Beleg zu
      [`LH-FA-ADM-003`](../../../../spec/lastenheft.md),
      [`LH-QA-REL-003`](../../../../spec/lastenheft.md).
- [ ] CDC-Abstand messbar (`cdc_capture_lag`-Kennzahl in der
      Metriken-View) — Teil-Beleg zu
      [`LH-FA-ADM-004`](../../../../spec/lastenheft.md),
      [`LH-QA-OPS-003`](../../../../spec/lastenheft.md).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für `harness/README.md` falls berührt (Sensors-Tabelle,
      falls ein neuer Rollout-Schritt entsteht).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Plan-Nachzug im selben Lauf · seit slice-009 (`AGENTS.md` §6,
`.claude/commands/implement-slice.md` Schritt 14):** Die §4-Prüfung ergab
eine Rückführung (`in-progress` → `next`, siehe §4/§6) — die Zeile für
`tools/schema/nacharbeit-observability.sql` unten trägt deshalb eine
**Reduktion** (Näherung statt echtem `cdc_capture_lag`, mit Begründung, nicht
still gestrichen); vier Zeilen sind gegenüber der ursprünglichen Planung
**neu hinzugekommen** (Domänenfehler-Sentinel, Store-Adapter, SQL-Text,
Composition Root) — sie tragen dieselben Liefer-Punkte (Fehlerzustands-
Sichtbarkeit, CDC-Abstand-Näherung), keinen dritten.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/model/errorstate.go` | neu | Domänen-Typ für Fehlerklassen-Sichtbarkeit (`ErrorClass`, sieben Kategorien), orientiert an [`ADR-0023`](../../../../docs/plan/adr/README.md)s Fehlerklassifikation |
| `internal/domain/errors/errors.go` | update | *(Plan-Nachzug)* Sentinel `ErrInvalidErrorClass` für die Invariante des neuen Domänen-Typs — derselbe Konstruktor-Pfad wie die übrigen Domänenfehler |
| `internal/application/port/outbound/heartbeat.go` | update | Heartbeat-Schreiber (slice-012) trägt zusätzlich den Fehlerzustand über eine neue Port-Methode `Fault`, statt eine zweite Tabelle einzuführen |
| `internal/adapters/driven/postgresstorage/heartbeat.go` | update | *(Plan-Nachzug)* `Fault`-Implementierung des Store-Adapters; `Beat` löscht einen zuvor gemeldeten Fehlerzustand wieder |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | *(Plan-Nachzug)* `UpsertHeartbeat` löscht `error_class` mit, neue Query `UpsertHeartbeatFault` |
| `tools/schema/schema.yaml` | update | *(Plan-Nachzug)* `process_heartbeat.error_class` (nullable text, kein Default = NULL = Normalbetrieb) |
| `tools/schema/nacharbeit-heartbeat.sql` | update | Heartbeat-View erweitert um Fehlerzustands-Spalte |
| `tools/schema/nacharbeit-observability.sql` | update, **reduziert** | `cdc.metrics`-View um `cdc_capture_lag` ergänzt — **nicht** wie ursprünglich geplant als Differenz Quell-Commit-Zeit/Persistenz-Zeit (siehe §4 Rückführung: das bräuchte Commit-Zeitstempel-Wiring durch die Replication-Adapter-Schicht), sondern als dokumentierte Näherung `now() − max(committed_at)` (Pipeline-Frische über die letzte persistierte Transaktion, Kommentar-Klasse Grenze in der Datei) |
| `internal/bootstrap/wiring.go` | update | *(Plan-Nachzug)* `classifyRunError`/`reportFault`: die Composition Root übersetzt den Lauf-Fehler in eine `ADR-0023`-Klasse und meldet ihn über `Fault`, bevor der Prozess auf einen Adapter-Fehler endet |
| `internal/domain/model/errorstate_test.go`, `internal/adapters/driven/postgresstorage/heartbeat_test.go`, `internal/bootstrap/heartbeat_internal_test.go` | neu/update | *(Plan-Nachzug)* Tests für Konstruktor-Invariante, Store-Adapter (`make test-store`) und Composition-Root-Klassifikation |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): slice-012 liegt in `done/`; kein anderes
Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): `committed_at`
  trägt heute die **Persistenz-Zeit** (DB-`DEFAULT now()`), nicht den
  tatsächlichen Quell-Commit-Zeitpunkt aus dem WAL — ein reales
  `cdc_capture_lag` bräuchte den Commit-Zeitstempel aus dem
  Replication-Stream gespiegelt bis in die Persistenz, was die
  Replication-Adapter-Schicht zusätzlich zu Application/Store-Adapter
  berührt (mehr als zwei Schichten) — Rückzug mit Zerlegung (WAL-
  Zeitstempel-Wiring als eigener Slice, CDC-Abstand-Metrik danach).
- `in-progress` → `open` (blockiert — Carveout?): die Heartbeat-Tabelle aus
  slice-012 existiert noch nicht (WIP-Reihenfolge verletzt) — Blocker,
  Priorität offen.

**Ausgang der Prüfung (Implementer-Lauf, vor jeder Code-Änderung, wie
`.claude/commands/implement-slice.md` Schritt 12 verlangt):** Die
`in-progress` → `next`-Bedingung ist **eingetreten**, mit einer Präzisierung
gegenüber der vorab benannten Schätzung. `pglogrepl` (Treiber, `ADR-0032`)
dekodiert `BEGIN`/`COMMIT`-Nachrichten bereits mit einem `CommitTime
time.Time`-Feld (`go doc github.com/jackc/pglogrepl.CommitMessage`,
netzlos gegen den vendored Modul-Cache geprüft) — ein reales
`cdc_capture_lag` bräuchte **keine neue Wire-Verbindung**, wie die
ursprüngliche Formulierung nahelegte. Trotzdem bleibt die Schichten-Zahl
unverändert zu groß: den Zeitstempel bis in `cdc.transaction.committed_at`
zu spiegeln, verlangt eine Änderung an `decode.Commit` (Replication-
Driving-Adapter), am `ChangeTransaction`/`SourcePosition`-Domänenmodell
(Domain), am `CaptureCommand`/`ChangeStorePort`-Vertrag (Application/Ports)
und an `InsertTransaction` (Store-Driven-Adapter) — vier Schichten statt
höchstens zwei. Die Rückführung trägt deshalb weiterhin: **Rückzug mit
Zerlegung**, mit einer kleineren, klareren Folge-Slice-Schätzung (reines
Zeitstempel-Durchreichen, kein neuer Protokoll-Zugriff) als eigenem
Liefer-Punkt. Der unabhängige Teil — Fehlerzustands-Sichtbarkeit
([`LH-FA-ADM-003`](../../../../spec/lastenheft.md),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md)) und eine als Grenze
dokumentierte `cdc_capture_lag`-Näherung auf Persistenz-Zeit-Basis (kein
echter [`LH-FA-ADM-004`](../../../../spec/lastenheft.md)-Abstand zur
Quelländerung) — bleibt vollständig innerhalb von Application/Store-Adapter
und wird in diesem Lauf trotzdem geliefert (§3 Plan-Nachzug), während dieser
Plan nach `next/` zurückgeht.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD vollständiges Häkchen + `make gates` grün + Review-Schluss ohne
  offenes HIGH-Finding; Fehlerzustands-Sichtbarkeit und CDC-Abstand-
  Messung am realen Adapter belegt (`make test-store`).
- Lerneintrag §7: geschärfte Regel oder benannte Spec-Lücke — ohne ihn
  bleibt der Slice abgelegt, nicht fertig.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- `committed_at` misst Persistenz-Zeit statt Quell-Commit-Zeit (siehe §4
  Rückführung) — **eingetreten**: die Rückführung trägt den Ausgang, dieser
  Slice geht nach `next/` zurück, kein `done/`-Übergang in diesem Lauf. Kein
  Ausgang im Sinn von §Offene Risiken werden bei Closure aufgelöst
  (Baseline-Regelwerk `modul-05-planning-harness.md`) — die drei dort
  benannten Ausgänge (eingetreten/entfallen/weiter offen) sind an den
  `in-progress → done`-Übergang gebunden, den dieser Lauf nicht vollzieht;
  die Zeile bleibt bis zur nächsten `in-progress`-Runde offen und wandert mit
  dem wiederaufgenommenen Plan. Die Näherung (Persistenz-Zeit als Proxy, mit
  Kommentar-Klasse Grenze) ist trotzdem umgesetzt — als Teil des
  unabhängigen, gelieferten Teils (§3, §4), nicht als Risiko-Ausgang.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Die berührten Sub-Areas
(Domänenmodell, Application/Store-Adapter) sind Segmente des GF-Baums der
Modus-Deklaration (`PGC` Greenfield, Doc führt) — jede erfüllt die Schwelle
≥ 2 von 3 Achsen. Keine Sub-Area zu grob.

**Vorgelagert — offene Beobachtungen sichten:** Register gelesen (neun
Einträge zum Zeitpunkt der Planung): `lese-doppelquelle` 2× (berührt — die
Metriken-View-Erweiterung verstärkt dieselbe Klasse potenziell weiter,
Beleg bei Closure prüfen), `d-migrate-nacharbeit` 2× (berührt — View-
Erweiterung über dieselbe Ausweichform). Übrige sieben ohne Bezug. Kein
Eintrag erreicht mit diesem Slice neu 3× — keine Lücke, aber
`lese-doppelquelle` beobachten (steht bei 2×, ein dritter Beleg würde die
Schwelle erreichen).

**Modus-Begründungsblock — Umfang.** Reiner GF-Hinweis genügt (siehe oben);
kein Sub-Area-Block.
