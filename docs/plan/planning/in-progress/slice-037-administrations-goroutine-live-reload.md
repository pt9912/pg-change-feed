# Slice slice-037: Administrations-Goroutine, Assembler-Live-Reload und Boot-Wechsel

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-12`](../welle-12.md) — der Nachweis, dass eine über SQL
beantragte Aktivierung/Deaktivierung den **laufenden** Erfassungspfad real
erreicht, ist `welle-12`s Closure-Trigger (§3); dieser Slice liefert genau
das, nachdem `slice-036` die Antrags-Seite geliefert hat.

**Bezug:** [`LH-FA-CFG-001`](../../../../spec/lastenheft.md),
[`LH-FA-CFG-002`](../../../../spec/lastenheft.md),
[`ADR-0050`](../../../../docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(nur umgesetzt — keine aktive ADR wird geändert, `ADR-0050` bleibt
`Accepted`).

**Berührte Spec-Stellen:** [`ARC-005`](../../../../spec/architecture.md)
(„SQL-Funktionen/Views"); `spec/architecture.md`s Sequenzdiagramm zu
`LH-FA-CFG-001.a` — `ADR-0050`s Folgepflicht verlangt eine Korrektur, die
SQL- und CLI-Pfad trennt.

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

**Ziel:** Der bereits laufende Capture-Prozess bekommt eine neue
Hintergrund-Goroutine (Muster: `runHeartbeat`/`runWALRetentionCheck` in
`internal/bootstrap/wiring.go`), die `LISTEN` auf dem
Administrations-Kanal hält (aus `slice-036`) und periodisch als Fallback
offene `pending`-Anträge aus `cdc.administration_request` liest. Für
jeden Antrag ruft sie den passenden Inbound Port
(`EnableTableUseCase`/`DisableTableUseCase`) auf, vermerkt das Ergebnis
im Antrags-Datensatz (`applied`/`failed` + Fehlertext) und trägt — **im
selben Prozess** — die laufende `Assembler`-Bindung nach
(`internal/adapters/driving/replication/mapper/mapper.go`): eine neue
synchronisierte Methode ergänzt/entfernt eine `TableBinding` in
`a.tables`, analog zu `observeRelation`s bestehender Laufzeit-Mutation
für bereits aktivierte Tabellen — hier erstmals für eine **neue**
Bindung, aus einer zweiten Goroutine. Zusätzlich baut
`internal/bootstrap/wiring.go` `Assembler.tables` beim Prozessstart
künftig aus `cdc.source_table` (`TableActivationPort.List`) statt
ausschließlich aus `CDC_TABLES` — `CDC_TABLES` bleibt als
Erstaktivierungs-Seed für eine leere Datenbank bestehen. Real gegen den
Compose-Stack belegt: `SELECT cdc.enable_table(...)` → Goroutine
verarbeitet den Antrag → eine danach ausgeführte Änderung an der neu
aktivierten Tabelle wird vom laufenden Prozess real erfasst, ohne
Neustart. Derselbe Nachweis für `cdc.disable_table(...)`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **CLI-Diagnose** (`slice-038`) — unabhängig lösbar, kein Reload-Bezug,
  eigener Slice.
- **Replica-Identity-Prüfung real umsetzen** — `ADR-0050`s Kontext
  benennt sie als separaten, hier nicht zu schließenden Befund; dieser
  Slice ruft den bestehenden `EnableTableUseCase` unverändert auf, ändert
  seine interne Prüftiefe nicht.
- **Consumer-Verwaltung über SQL** — bereits als Out-of-Scope der
  gesamten Welle benannt (`welle-12` §6).
- **`spec/architecture.md`s Sequenzdiagramm-Korrektur real ausführen
  als Teil dieses Slices bleibt bewusst inkludiert** — `ADR-0050`s
  Folgepflicht wird HIER erledigt (nicht ausgeschlossen), weil erst
  dieser Slice den asynchronen SQL-Pfad real existieren lässt, den das
  Diagramm bisher nicht zeigt.

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

- [x] Neue Administrations-Goroutine (`internal/bootstrap/wiring.go`,
      Muster `runHeartbeat`/`runWALRetentionCheck`) verarbeitet reale
      `pending`-Anträge aus `cdc.administration_request`: ruft
      `EnableTableUseCase`/`DisableTableUseCase` auf, schreibt
      `applied`/`failed` zurück. Beleg: `runAdministration`/
      `processAdministrationRequests`/`applyAdministrationRequest`
      (`internal/bootstrap/wiring.go`), real gegen den Compose-Stack
      geprüft (nächstes Item) und dreifach grün in `make test-integration`.
- [x] `Assembler` bekommt eine neue, synchronisierte Methode, die eine
      `TableBinding` zur Laufzeit hinzufügt/entfernt; `a.tables`-Zugriff
      aus zwei Goroutinen ist race-frei (`go test -race`). Beleg:
      `Assembler.AddBinding`/`RemoveBinding` (`sync.RWMutex`,
      `internal/adapters/driving/replication/mapper/mapper.go`),
      `TestAssemblerLiveReloadIsRaceFree` (`mapper_test.go`) — real rot
      gesehen ohne die Sperren (Mutation entfernt, `go test -race` meldete
      die Data Race exakt an `lookupBinding`/`AddBinding`), danach wieder
      grün mit den Sperren.
- [x] `internal/bootstrap/wiring.go` baut `Assembler.tables` beim Start
      aus `cdc.source_table` (`TableActivationPort.List`); `CDC_TABLES`
      bleibt Erstaktivierungs-Seed. Beleg: `activatedTableBindings`
      (`internal/bootstrap/wiring.go`), `receive.Config.Tables` liest jetzt
      diesen Rückgabewert statt `cfg.Tables` direkt.
- [x] Real gegen den Compose-Stack belegt: `cdc.enable_table(...)` →
      Goroutine verarbeitet → neu aktivierte Tabelle wird vom laufenden
      Prozess ohne Neustart erfasst. Derselbe Nachweis für
      `cdc.disable_table(...)`. Beleg: `tools/harness/run-integration-tests.sh`
      Abschnitt „SQL-Administration Live-Reload-Beleg" gegen
      `feed_mvp_sql_admin` (bewusst nicht in `CDC_TABLES`) — dreifach grün
      in `make test-integration`, Feed-Container läuft nach beiden
      Anträgen unverändert weiter (kein `docker restart`).
- [x] `spec/architecture.md`s Sequenzdiagramm zu `LH-FA-CFG-001.a`
      korrigiert: SQL-Pfad (asynchron, Antrag → Queue → Goroutine → Port)
      von CLI-Pfad (synchron, direkter Aufruf) getrennt (`ADR-0050`
      Folgepflicht). Beleg: zwei Sequenzdiagramme unter `LH-FA-CFG-001.a`
      (CLI/SQL), keine Wellen-/Slice-/ADR-Bezüge im Diagramm selbst
      (Hard Rule 3.4).
- [x] `make gates` grün, `make test` (Race-Detector) und
      `make test-integration` dreimal in Folge grün. Beleg:
      `make test` läuft jetzt mit `-race` über `TOOLCHAIN_RACE_IMAGE`
      (Makefile-Plan-Nachzug, Debian-basiert wegen `gcc`); alle vier
      Kommandos real ausgeführt, siehe Bericht.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt wird —
      Implementer entscheidet und begründet im Plan-Nachzug (die
      Architektur-Sicht-Korrektur oben zählt bereits als eigenes
      DoD-Item, nicht doppelt hier). Entscheidung: `harness/README.md`
      §Werkzeuge aktualisiert (`make test`-Zeile: Race-Detector +
      Debian-Image-Begründung; `make test-integration`-Zeile: neuer
      Live-Reload-Beleg) — beide sind öffentliche Sensor-Beschreibungen,
      deren Verhalten sich mit diesem Slice geändert hat.
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
| `internal/bootstrap/wiring.go` | update | neue Administrations-Goroutine; Boot-Wechsel auf `TableActivationPort.List` |
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `sync.RWMutex` + `AddBinding`/`RemoveBinding` für neue/entfallende `TableBinding`-Einträge zur Laufzeit |
| `internal/adapters/driven/postgresstorage/administrationrequest.go` | neu | `AdministrationRequestAdapter` (Port-Implementierung) + `AdministrationListener` (`LISTEN`-Wecksignal, eigene Verbindung mit Reconnect) |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | SQL-Texte der Antrags-Queue (`SelectPendingAdministrationRequests`, `UpdateAdministrationRequestApplied/Failed`) |
| `internal/application/port/outbound/administrationrequest.go` | neu | `AdministrationRequestPort` (`ListPending`/`MarkApplied`/`MarkFailed`) |
| `internal/domain/model/administrationrequest.go` | neu | `AdministrationRequest`, `AdministrationRequestID`, `AdministrationRequestKind` |
| `internal/adapters/driving/replication/receive/receive.go` | update | `Stream.Assembler()`-Zugriffsmethode — die Administrations-Goroutine trägt die Bindung des laufenden Streams nach, kein zweiter Übersetzer |
| `spec/architecture.md` | update | Sequenzdiagramm zu `LH-FA-CFG-001.a`: SQL-Pfad (asynchron) von CLI-Pfad (synchron) getrennt, zwei Diagramme statt eines |
| `tools/harness/run-integration-tests.sh` | update | realer End-zu-End-Nachweis: `cdc.enable_table`/`cdc.disable_table` → Goroutine verarbeitet → laufender Feed-Container erfasst/stoppt ohne Neustart |
| `Makefile` | update | `make test` läuft jetzt mit `-race` (`TOOLCHAIN_RACE_IMAGE`, Debian-basiert — der Race-Detector braucht `gcc`, das Alpine-Toolchain-Image trägt keinen) |
| `internal/adapters/driven/postgresstorage/administrationrequest_test.go` | update | Adapter-/Listener-Tests (`ListPending`/`MarkApplied`/`MarkFailed`, `WaitForNotification`) — ergänzt die bereits vorhandenen Funktions-Tests aus `slice-036` |
| `internal/adapters/driving/replication/mapper/mapper_test.go` | update | `AddBinding`/`RemoveBinding`-Tests + `TestAssemblerLiveReloadIsRaceFree` (`ADR-0050` Fitness Function) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-036` liegt in `done/`,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass Goroutine, Assembler-Synchronisation, Boot-Wechsel und
  Architektur-Sicht-Korrektur zusammen mehr als drei Liefer-Punkte oder
  mehr als zwei Schichten in einer Review-Sitzung nicht mehr prüfbar
  machen, gehört das zurück zur Zerlegung (z. B. Boot-Wechsel als eigener
  Folge-Slice).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test`
(Race-Detector) **und** `make test-integration` dreimal in Folge grün
**und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der `Assembler.tables`-Zugriff aus zwei Goroutinen (Capture-Stream,
  Administrations-Goroutine) könnte eine Data Race einführen, wenn die
  Synchronisation unvollständig ist — `ADR-0050`s eigene Fitness Function
  benennt `go test -race` als Prüfpflicht. **Ausgang:** <bei Closure
  einzutragen>
- Ein offener Antrag über einen Prozess-Neustart hinweg (`ADR-0050`s
  Konsequenz: „bleibt pending, wird beim nächsten Boot/Poll erneut
  abgeholt") könnte doppelt verarbeitet werden, wenn `EnableTableUseCase`
  nicht real idempotent ist. **Ausgang:** <bei Closure einzutragen>
- Der Fallback-Poll-Intervall könnte zu grob gewählt werden und die
  Latenz zwischen SQL-Antrag und realer Wirksamkeit unnötig verlängern,
  wenn `NOTIFY` verpasst wird (z. B. nach einem Verbindungsabbruch der
  `LISTEN`-Verbindung). **Ausgang:** <bei Closure einzutragen>

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
Treffer für `PGC`: `BEO-PGC/verwaltung-keine-sql-administration` (0×,
benannt nicht gezählt — dieser Slice liefert den zweiten und
architektonisch entscheidenden Baustein der Auflösung). Keiner der
übrigen Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
