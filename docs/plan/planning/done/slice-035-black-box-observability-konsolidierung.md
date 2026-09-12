# Slice slice-035: Black-Box-Observability-Konsolidierung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-11`](../welle-11.md) — der Nachweis, dass die
Observability-Signale real über externe Schnittstellen lesbar sind, ist
`welle-11`s Closure-Trigger (§3), kein Einzel-Slice-DoD.

**Bezug:** [`LH-FA-ADM-002`](../../../../spec/lastenheft.md),
[`LH-FA-ADM-005`](../../../../spec/lastenheft.md) — dieser Slice erweitert
die Testabdeckung für bestehende Verträge, ändert sie nicht. Kein aktives
ADR wird geändert.

**Berührte Spec-Stellen:** — (reine Testinfrastruktur, keine
Verhaltensänderung an einer Spec-Stelle).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

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

**Ziel:** Zwei bislang schwarz nicht belegte Observability-Signale real
gegen den Compose-Stack testen: (a)
[`LH-FA-ADM-005`](../../../../spec/lastenheft.md) (Verarbeitungsrückstand)
— ein neuer Testfall liest `cdc.consumer_status` real über SQL (kein
Go-Adapter-Import) und belegt, dass die Differenz
`latest_commit_position - acknowledged_position` einen realen Rückstand
zeigt, solange ein registrierter Consumer nicht bestätigt hat, und auf
den Rückstand `0` fällt, sobald er die aktuelle Position bestätigt (das
Aufrufmuster ist bereits real vorhanden: `register-consumer`/
`acknowledge-consumer` über `docker exec`, `slice-027`). (b)
[`LH-FA-ADM-002`](../../../../spec/lastenheft.md) (Betriebsstatus,
Happy Path) — ein neuer Testfall liest `cdc.heartbeat` real über SQL und
belegt, dass ein betriebsbereites System eine frische Lebenszeichen-Zeile
(niedriges `age_seconds`) ohne `error_class` zeigt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`LH-FA-ADM-003` (sichtbare Fehlerzustände)** — bereits real black-box
  belegt: `TestMVPSchemaChangeIncompatibleTypeChange`s
  `awaitHeartbeatErrorClass`-Muster (`slice-033`) liest exakt dieselbe
  `cdc.heartbeat`-Sicht für den Fehlerfall; kein zweiter Testfall nötig,
  dieser Slice konsolidiert nur den Happy-Path-Gegenpart.
- **`LH-FA-ADM-004` (messbarer CDC-Abstand)** — bereits real black-box
  belegt über den `cdc_capture_lag`-Lasttest-Beleg in
  `tools/harness/run-integration-tests.sh`; anderer Vorgang, kein
  Testauftrag hier.
- **SQL-Administrationsfunktionen, CDC-Deaktivierung, CLI-Diagnose** —
  real fehlende Fähigkeiten, keine Testabdeckungslücke; bereits als
  eigene Feature-Welle vorgemerkt (`BEO-PGC/verwaltung-keine-sql-administration`).
- **Schwellenwert-Entscheidungen (healthy/degraded/unhealthy)** — Bestand
  bleibt bewusst stehen: `nacharbeit-observability.sql`s eigener
  Kommentar zu `cdc.heartbeat` dokumentiert bereits, dass die View nur
  den rohen Zeitstempel/Alter liefert, keine Schwellenwert-Entscheidung
  (Sache des lesenden Systems) — dieser Slice testet die Sichtbarkeit
  der Rohwerte, nicht eine Schwellenwert-Klassifikation, die es noch gar
  nicht gibt.

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

- [x] `LH-FA-ADM-005` erfüllt: neuer Testfall liest `cdc.consumer_status`
      real über SQL für einen registrierten, noch nicht bestätigenden
      Consumer (Rückstand > 0) und nach `acknowledge-consumer` (Rückstand
      = 0). **Präzisierung (Plan-Nachzug, siehe §3):** black-box real
      belegt ist der Rückstand nach einer ersten Bestätigung, gefolgt von
      einer weiteren real erfassten Änderung (Rückstand > 0), und dessen
      Auflösung nach der zweiten Bestätigung (Rückstand = 0) — nicht die
      wörtliche „niemals bestätigt"-Situation, die die Sicht strukturell
      nicht als Zahl ausgibt (siehe Begründung in §3). Beleg: neuer
      Abschnitt in
      [`tools/harness/run-integration-tests.sh`](../../../../tools/harness/run-integration-tests.sh)
      (`BACKLOG_CONSUMER`/`BACKLOG_TABLE`), dreifach grün über
      `make test-integration` — siehe Bericht an den Reviewer.
- [x] `LH-FA-ADM-002` erfüllt (Happy Path): neuer Testfall liest
      `cdc.heartbeat` real über SQL während des laufenden Betriebs und
      belegt eine frische, fehlerfreie Lebenszeichen-Zeile. Beleg:
      `TestMVPHeartbeatHealthy` in
      [`test/integration/integration_test.go`](../../../../test/integration/integration_test.go),
      Schwelle 15s (Plan-Nachzug, siehe §3), dreifach grün über
      `make test-integration`.
- [x] `make gates` grün, `make test-integration` dreimal in Folge grün —
      beide Läufe lokal ausgeführt (tatsächlich vier grüne
      `make test-integration`-Läufe plus zwei gezielt rot geführte
      Mutationsläufe zum Beleg der Wächter-Wirksamkeit — siehe Bericht an
      den Reviewer).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-035.md`](../../../reviews/review-slice-035.md)
      (1 LOW, kein Merge-Blocker; Plan-Abweichung bei `LH-FA-ADM-005`
      eigenständig geprüft und als legitime Präzisierung bestätigt, keine
      Verwässerung). Verifikation in
      [`docs/reviews/verify-slice-035.md`](../../../reviews/verify-slice-035.md)
      (DoD eigenständig nachgeprüft, dreifach real reproduziert).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt wird — kein
      öffentlicher Vertrag berührt (reine Testabdeckung bestehender,
      bereits dokumentierter Lesezugriffswege); keine Doku-Änderung.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. — entfällt: `../reconciliation.md` existiert nicht (Repo ist Greenfield, kein Brownfield-Bootstrap).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      Keine Beobachtung angefallen — reine Testabdeckung ohne neuen Fund.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb — Prüfung läuft bei der `welle-11`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | update | **Plan-Nachzug:** neuer Testfall `TestMVPHeartbeatHealthy` (`LH-FA-ADM-002` Happy Path) — reine SQL-Lesung ohne CLI-Bedarf, dasselbe Muster wie `awaitHeartbeatErrorClass` (`slice-033`). Läuft im ersten `go test`-Aufruf des Runner-Skripts, vor `TestMVPSchemaChangeIncompatibleTypeChange` (jener setzt `error_class` dauerhaft). **Schwellenwert-Entscheidung:** 15s, identisch mit der Produktions-Healthcheck-Schwelle `heartbeatStaleAfter` (`internal/bootstrap/wiring.go`, 3 × `heartbeatInterval` = 15s) — derselbe Wert, den der Compose-Healthcheck bereits gegen dieselbe `cdc.heartbeat`-Zeile prüft, statt einer neu erfundenen Zahl. |
| `tools/harness/run-integration-tests.sh` | update | **Plan-Nachzug:** (1) neuer Bash-Abschnitt für `LH-FA-ADM-005` — Implementer-Entscheidung gegen einen Go-Testfall, weil `register-consumer`/`acknowledge-consumer` nur extern per `docker exec` gegen den laufenden Feed-Container aufrufbar sind (der Toolchain-Container, in dem die Go-Tests laufen, hat keinen Docker-Socket-Zugriff) — dasselbe Muster wie der bestehende Black-Box-CLI-Rundlauf (`slice-027`). **Design-Entscheidung/Präzisierung ggü. der ursprünglichen Formulierung:** `cdc.consumer_status` liest `latest_commit_position` über eine auf `cp.source_id` korrelierte Unterabfrage; ein Consumer, der *noch nie* bestätigt hat, hat keine `cdc.consumer_position`-Zeile, wodurch die Sicht dafür `NULL` statt eines Rückstands liefert (kein SQL-Fehler, aber auch keine belegbare Zahl > 0). Real umgesetzt und belegt: eine erste externe Bestätigung bindet `source_id`, eine danach real erfasste weitere Änderung hebt `latest_commit_position` über die bestätigte Position — der Rückstand wird über eine reine SQL-Lesung sichtbar (> 0); die zweite externe Bestätigung senkt ihn auf 0. (2) `TestMVPHeartbeatHealthy` ins bestehende `-run`-Muster des ersten `go test`-Aufrufs aufgenommen (`BEO-PGC/test-runner-stiller-ausschluss`, weiter offen, 1× — dieser Slice vermeidet die Lücke erneut, wie bereits `slice-034`). |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einem
anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass die beiden Observability-Signale (Rückstand, Betriebsstatus)
  jeweils eigene, größere Testinfrastruktur brauchen, gehört das zurück
  zur Zerlegung — je ein Slice pro Signal.
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

- Der Rückstandstestfall könnte durch Timing/Nebenläufigkeit (Erfassung
  läuft weiter, während der Test die Position prüft) flaky werden, wenn
  nicht sauber auf „alle erwarteten Changes erfasst" gewartet wird, bevor
  die Rückstands-Messung gelesen wird. **Ausgang: entfallen** — Reviewer
  und Verifier haben den Testfall unabhängig voneinander dreifach real
  reproduziert (identischer Rückstandswert 1464 → 0 in jedem Lauf), keine
  Flakiness gefunden.
- Der Betriebsstatus-Testfall könnte gegen eine bereits durch einen
  anderen Testfall (z. B. `slice-033`s Fehlerklassen-Test, der den
  Feed-Container real beendet) beeinträchtigte `cdc.heartbeat`-Zeile
  laufen, wenn die Testreihenfolge im Runner-Skript nicht beachtet wird.
  **Ausgang: entfallen** — Reviewer prüfte die Platzierung eigenständig:
  `TestMVPHeartbeatHealthy` läuft im vorderen `-run`-Muster, vor dem
  containerbeendenden Testfall aus `slice-033`; dreifach real grün ohne
  Beeinträchtigung.

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

- **Was hat funktioniert:** Der Implementer fand real, dass
  `cdc.consumer_status`s Rückstands-Semantik von der wörtlichen
  Ziel-Formulierung abweicht (LEFT JOIN liefert `NULL` statt eines
  positiven Rückstands für einen nie bestätigenden Consumer), passte den
  Testablauf präzise an den tatsächlichen Lastenheft-Wortlaut an und
  dokumentierte die Abweichung transparent im Plan-Nachzug statt sie zu
  verschweigen. Reviewer und Verifier prüften diese Präzisierung beide
  unabhängig gegen den Lastenheft-Text und den Code (LEFT JOIN, Go-Pfad
  `UpsertConsumerPosition`) und bestätigten: legitime Präzisierung, keine
  Verwässerung.
- **Was ging anders als geplant:** Der Reviewer fand ein LOW-Finding
  (F-1: ein unveränderter Kommentar im Runner-Skript zählt den neuen
  Rückstands-Beleg-Block nicht unter den Voraussetzungen für den
  nachfolgenden containerbeendenden Test auf) — kein Merge-Blocker, keine
  Fixrunde nötig.
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3× —
  der Normalfall. `BEO-PGC/test-runner-stiller-ausschluss` bleibt
  unverändert bei 1× (dieser Slice hat die Lücke für seine beiden
  Testfälle korrekt vermieden).
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen — siehe §2-Begründung.
- **Folge-Slices:** keine — `welle-11` schließt mit diesem Slice.
- **Risiken aus §6:** beide *entfallen* — siehe §6 für Begründung.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-11` offen) —
  Prüfung läuft bei der `welle-11`-Closure.

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
Treffer für `PGC`: `BEO-PGC/verwaltung-keine-sql-administration` (0×,
benannt nicht gezählt — nicht Gegenstand dieses reinen Testabdeckungs-
Slices), `BEO-PGC/test-runner-stiller-ausschluss` (1×, weiter offen —
relevant: neue Testfälle müssen im `-run`-Muster des Runner-Skripts
erfasst werden). Keiner der übrigen Treffer erreicht mit diesem Slice
3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
