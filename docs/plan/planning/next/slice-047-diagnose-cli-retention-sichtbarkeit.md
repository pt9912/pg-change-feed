# Slice slice-047: Diagnose-CLI-Erweiterung für Retention-Sichtbarkeit

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Eröffnungs-Recherche (Fork) zur ursprünglich in der
Roadmap vorgemerkten Welle „E2E-Abdeckung — Retention" fand real: Die
SQL-View-Black-Box-Ebene ist bereits vollständig geliefert
(`TestMVPRetentionBlockersViewShowsFurthestBehindConsumer`,
`TestMVPMetricsCarriesStorageBytes`, beide `slice-045`/`046`, beide reine
SQL-Lesetests ohne Go-Domain-Import — exakt die Testebene, die `welle-11`
für `cdc.active_tables`/`cdc.consumer_status` erst nachliefern musste). Die
einzige real verbleibende Lücke ist die fehlende `docker exec`-CLI-
Sichtbarkeit (`internal/bootstrap/wiring.go`s `Diagnose`-Funktion zeigt
nichts zu Retention) — ein einzelner Slice ohne Closure-Bedingung jenseits
seiner eigenen DoD, siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

**Bezug:** [`LH-FA-SST-003`](../../../../spec/lastenheft.md) (CLI-Diagnose),
[`LH-FA-RET-005`](../../../../spec/lastenheft.md),
[`LH-FA-RET-006`](../../../../spec/lastenheft.md) (Sichtbarkeit, jetzt auch
über die CLI statt nur über SQL).

**Berührte Spec-Stellen:** — (CLI-Ausgabe-Erweiterung auf bereits
bestehenden Daten, keine neue Architektur-Sicht-Aussage).

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

**Ziel:** Der bestehende `diagnose`-CLI-Befehl
(`internal/bootstrap/wiring.go`, `slice-038`) bekommt einen zusätzlichen
Abschnitt: den ggf. aktuell blockierenden Consumer je Quelle (aus
`cdc.retention_blockers`, `slice-045`) und den aktuellen
`cdc_storage_bytes`-Wert (aus `cdc.metrics`, `slice-046`) — dieselbe
Lese-Disziplin wie die bestehenden Abschnitte (direkte SQL-Abfrage über
den `CDC_READER_DSN`-Pool, kein neuer Port). Ein neuer, externer
`docker exec`-Diagnose-Beleg in `tools/harness/run-integration-tests.sh`
(Muster: bestehender „CLI-Diagnose-Beleg (Normalbetrieb)"-Abschnitt)
zeigt real: kein Blocker (leere Quelle), dann ein realer Blocker
(zurückhängender Consumer), dann der reale `cdc_storage_bytes`-Wert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Neue Berechnungslogik** — die CLI liest ausschließlich bereits
  bestehende SQL-Sichten (`cdc.retention_blockers`, `cdc.metrics`); keine
  neue Domain-/Use-Case-Logik entsteht.
- **Automatische Warnung/Alarmierung bei blockierenden Consumern** —
  bereits in `slice-045` §1 als „anderer Vorgang" ausgeschlossen, gilt
  unverändert.
- **Kombinierter End-zu-Ende-Rundlauf über alle vier `welle-13`-
  Fähigkeiten hinweg** (Consumer blockiert → sichtbar → bestätigt →
  gelöscht → Storage-Bytes reflektiert das) — der bestehende
  `slice-044`-Retention-Beleg in `run-integration-tests.sh` deckt die
  Löschausführung selbst bereits ab; dieser Slice fügt nur die
  CLI-Sichtbarkeit hinzu, keinen neuen kombinierten Testfall. Bestand
  bleibt bewusst stehen: ein eigener Vorgang, falls ein Betriebs-Bedarf
  dafür entsteht.

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

- [ ] `diagnose`-CLI zeigt real den aktuell blockierenden Consumer je
      Quelle (aus `cdc.retention_blockers`) und den `cdc_storage_bytes`-Wert
      (aus `cdc.metrics`) — real gegen mindestens einen Zustand ohne
      Blocker und einen mit realem Blocker getestet.
- [ ] `LH-FA-SST-003` real erweitert: ein externer `docker exec`-Beleg in
      `tools/harness/run-integration-tests.sh` zeigt beide Zustände in der
      `diagnose`-Ausgabe, ohne den laufenden Feed-Container zu beenden
      (analog zum bestehenden CLI-Diagnose-Beleg-Muster).
- [ ] `make gates` grün, `make test-integration` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` nennt die erweiterte
      `diagnose`-Ausgabe (Abschnitt „Aufbewahrung (Retention)").
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
| `internal/bootstrap/wiring.go` (`Diagnose`) | update | neuer Retention-Abschnitt (Blocker, Storage-Bytes) |
| `tools/harness/run-integration-tests.sh` | update | neuer `docker exec`-Diagnose-Beleg für Retention |
| `docs/user/benutzerhandbuch.md` | update | erweiterte `diagnose`-Ausgabe dokumentiert |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-13` liegt in `done/`
(`slice-045`/`046` liefern `cdc.retention_blockers`/`cdc.metrics`),
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten bei einer reinen CLI-Ausgabe-Erweiterung auf bereits
  bestehenden SQL-Sichten — falls doch, wäre das ein Zeichen für eine
  unerwartet komplexe Ausgabeformat-Änderung.
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

- Eine Quelle ohne aktuellen Blocker (`cdc.retention_blockers` liefert
  keine Zeile) könnte die CLI-Ausgabe fälschlich als Fehlerzustand statt
  als „kein Blocker" zeigen — dieselbe Klasse Randfall wie `slice-038`s
  §6 Risiko 2 und `slice-045`s §6 Risiko 1. **Ausgang:** <bei Closure
  einzutragen>
- Die neue Retention-Sektion könnte das bestehende, stabile
  `diagnose`-Ausgabeformat so verändern, dass `run-integration-tests.sh`s
  bereits bestehende Text-Assertions gegen `LH-FA-ADM-002`…`005`
  (Normalbetrieb/Fehlerzustand-Belege) brechen. **Ausgang:** <bei Closure
  einzutragen>

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
Treffer für `PGC`: `BEO-PGC/test-runner-stiller-ausschluss` (1×, ein neuer
Testfall im `run-integration-tests.sh`-Textmuster muss real ins
Ausgabe-Matching aufgenommen werden), `BEO-PGC/test-isolation-geteilter-zustand`
(1×, kein direkter Bezug — dieser Slice fügt keinen neuen Zustands-
teilenden Testfall hinzu). Keiner erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
