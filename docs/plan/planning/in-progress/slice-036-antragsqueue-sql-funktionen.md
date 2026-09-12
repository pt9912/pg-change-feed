# Slice slice-036: Antrags-Queue und schreibende SQL-Funktionen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-12`](../welle-12.md) — der Nachweis, dass ein
Administrator über SQL real aktivieren/deaktivieren kann, ist `welle-12`s
Closure-Trigger (§3); dieser Slice liefert nur die Antrags-Seite, ohne
die sich `slice-037`s Goroutine noch nichts zu verarbeiten hätte.

**Bezug:** [`LH-FA-ADM-001`](../../../../spec/lastenheft.md),
[`ADR-0050`](../../../../docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(nur umgesetzt — keine aktive ADR wird geändert, `ADR-0050` bleibt
`Accepted`).

**Berührte Spec-Stellen:** [`ARC-005`](../../../../spec/architecture.md)
(„SQL-Funktionen/Views" als Driving-Adapter-Fläche, durch `ADR-0050`
geschärft).

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

**Ziel:** `ADR-0050`s Antrags-Seite real bauen: eine neue Tabelle im
`cdc`-Schema (z. B. `cdc.administration_request` — exakter Name
Implementer-Entscheidung, Plan-Nachzug), die einen Antrag trägt (Quelle,
Schema, Tabelle, Art `enable`/`disable`, Zeitstempel, Status
`pending`/`applied`/`failed`, Fehlertext), und zwei SQL-Funktionen
`cdc.enable_table(...)`/`cdc.disable_table(...)`, die **ausschließlich**
einen Antrags-Datensatz schreiben und `pg_notify` auf einem
Administrations-Kanal senden. Real gegen PostgreSQL getestet: der
Funktionsaufruf legt real eine Zeile mit Status `pending` an, das
`NOTIFY` ist real über `LISTEN` beobachtbar.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Verarbeitung der Anträge** (Administrations-Goroutine, Aufruf von
  `EnableTableUseCase`/`DisableTableUseCase`, `Assembler`-Live-Reload) —
  Folge-Slice `slice-037` übernimmt das explizit; dieser Slice liefert
  nur die Antrags-Seite, ohne dass irgendetwas die Anträge bislang
  abholt (Anträge bleiben real `pending`, bis `slice-037` existiert).
- **Direktes Schreiben von `cdc.source_table`/`cdc.schema_version`/
  Publication durch die SQL-Funktionen** — Bestand bleibt bewusst außen
  vor: `ADR-0050`s zentrale Entscheidung ist, dass diese Zielzeilen
  ausschließlich über `TableActivationAdapter`
  (aufgerufen vom Inbound Port) geschrieben werden; eine SQL-Funktion,
  die das umginge, widerspräche der ADR direkt.
- **CLI-Diagnose, Consumer-Verwaltung über SQL** — bereits als eigener
  Slice (`slice-038`) bzw. Out-of-Scope der Welle benannt.

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

- [x] Neue Antrags-Tabelle im neutralen Schema (`tools/schema/schema.yaml`
      oder, falls d-migrate keine SQL-Funktionen aus dem Schemamodell
      generieren kann, `tools/schema/nacharbeit-*.sql`-Ausweichform nach
      dem Muster von `nacharbeit-heartbeat.sql`/`nacharbeit-observability.sql`
      — Implementer-Entscheidung, Plan-Nachzug), ausgerollt über
      `make schema-rollout`. Beleg: `tools/schema/schema.yaml`
      (`administration_request`), real ausgerollt über
      `bash tools/harness/run-store-tests.sh` (`make schema-rollout`
      intern, Exit 0, Tabelle real angelegt).
- [x] `cdc.enable_table(...)`/`cdc.disable_table(...)` real als SQL-
      Funktionen angelegt; schreiben ausschließlich einen Antrags-
      Datensatz (Status `pending`) und senden `pg_notify` auf einem
      Administrations-Kanal — real gegen PostgreSQL getestet
      (`make test-store`-Muster: Funktionsaufruf, resultierende Zeile,
      `LISTEN`-Beobachtung des `NOTIFY`). Beleg:
      `tools/schema/nacharbeit-administration.sql`,
      `internal/adapters/driven/postgresstorage/administrationrequest_test.go`
      (`TestAdministrationRequestEnableTableWritesPendingRequestAndNotifies`,
      `TestAdministrationRequestDisableTableWritesPendingRequestAndNotifies`
      — beide real grün, `go test -v -run TestAdministrationRequest`
      gegen eine frische PostgreSQL-18-Instanz).
- [x] `make gates` grün. Beleg: `baseline-verify` OK (54 Dateien),
      `d-check` 0 Befunde (296 Dateien), `commit-traceability` OK,
      `a-check` 0 Befunde.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für `harness/README.md` §Sensors/`AGENTS.md`, falls ein
      neuer Sensor/Vertrag entsteht — Implementer entscheidet und
      begründet im Plan-Nachzug. Entscheidung: kein Update nötig —
      korrigierte Begründung (Review-Finding F-3, `review-slice-036.md`):
      Die `make schema-rollout`-Zeile in `harness/README.md` nennt
      `nacharbeit-views.sql` tatsächlich namentlich, aber ausschließlich
      im Kontext einer **aufgelösten** Post-Compare-Drift-Geschichte
      („zurückgebaut · seit slice-016"). Dieselbe Zeile nennt weder
      `nacharbeit-roles.sql` noch `nacharbeit-observability.sql` noch
      `nacharbeit-heartbeat.sql`, obwohl alle drei seit ihrer jeweiligen
      Einführung aktiv und im Makefile verankert sind — der Präzedenzfall
      trägt Auflösungen, nicht aktive Ausweichformen.
      `nacharbeit-administration.sql` ist mit diesem Slice neu entstanden
      und bleibt aktiv (offener dritter Fall im Beobachtungs-Register
      `BEO-PGC/d-migrate-nacharbeit`) — sie fällt damit unter dieselbe
      Nicht-Nennung wie die drei bestehenden aktiven Ausweichformen, nicht
      unter den Views-Präzedenzfall. Sobald d-migrate die Funktionsklasse
      ebenso wie Views seit 1.3.1 aus dem Post-Compare-Fingerabdruck
      ausblendet und die Ausweichform zurückgebaut wird, bekommt die Zeile
      denselben Nachtrag wie bei den Views; kein neues Gate/Target
      entstanden.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `PGC`), `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Beleg:
      `../observations/BEO-PGC/d-migrate-nacharbeit/evidence/slice-036.md`
      (weitere Datei im bestehenden Verzeichnis ergänzt).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/schema.yaml` oder `tools/schema/nacharbeit-administration.sql` | neu/update | Antrags-Tabelle |
| `tools/schema/nacharbeit-administration.sql` (oder gleichwertig) | neu | `cdc.enable_table`/`cdc.disable_table`-SQL-Funktionen |
| `internal/adapters/driven/postgresstorage/` | update (Test) | Test gegen die neuen SQL-Funktionen real gegen PostgreSQL |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einem
anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass Antrags-Tabelle und beide SQL-Funktionen zusammen mehr als drei
  Liefer-Punkte oder mehr als zwei Schichten in einer Review-Sitzung
  nicht mehr prüfbar machen, gehört das zurück zur Zerlegung (z. B.
  `enable`/`disable` als getrennte Slices).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** Adapter-Test gegen reale
PostgreSQL grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- d-migrate könnte keine SQL-Funktionen aus dem neutralen Schemamodell
  generieren können (nur Tabellen/Views bekannt) — dann muss die
  Antrags-Tabelle über `schema.yaml` laufen, die Funktionen aber über die
  bestehende `nacharbeit-*.sql`-Ausweichform (`ADR-0043`
  Re-Evaluierungs-Trigger). **Ausgang:** <bei Closure einzutragen>
- Die Rollentrennung (`ADR-0047`) könnte verlangen, dass die neuen
  SQL-Funktionen nur einer bestimmten Rolle (`cdc_admin`?) gewährt
  werden, nicht `cdc_reader`/`cdc_capture` — das muss beim Grants-Schritt
  (`nacharbeit-roles.sql`-Muster) bedacht werden. **Ausgang:** <bei
  Closure einzutragen>

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
benannt nicht gezählt — dieser Slice liefert den ersten Baustein der
Auflösung, ohne sie selbst zu vollenden). Keiner der übrigen Treffer
erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).

## Plan-Nachzug (Implementer-Entscheidungen)

Nachgetragen während der Implementierung — §6/§7 bleiben davon unberührt
(deren Ausgänge/Lerneintrag sind Sache des Übergangs nach `done/`).

- **Tabellenname:** `cdc.administration_request` — der in §1 vorgeschlagene
  Name übernommen, kein abweichender Bedarf aufgetreten.
- **d-migrate-Funktionsunterstützung:** Ja, mit einer real geprüften
  Einschränkung. Der `functions:`-Knoten des neutralen Schemamodells
  generiert die DDL korrekt (`schema generate --target postgresql`,
  inklusive `security: definer`/`search_path`) — aber `schema migrate
  --execute` bricht für **jede** dort deklarierte Funktion mit
  `POST_EXECUTE_DRIFT` (Exit 5) ab, real isoliert mit einer trivialen
  No-Arg-Funktion gegen eine frische PostgreSQL-18-Instanz (dieselbe
  Objektklasse betroffen, nicht auf `cdc.enable_table`/`cdc.disable_table`
  beschränkt; eine Tabelle im selben Lauf rollt dagegen Exit 0). Die beiden
  Funktionen laufen deshalb über die etablierte Ausweichform
  `tools/schema/nacharbeit-administration.sql` (`ADR-0043`
  Re-Evaluierungs-Trigger, dieselbe Klasse wie `nacharbeit-roles.sql`); die
  Antrags-Tabelle selbst läuft unverändert über `tools/schema/schema.yaml`.
  Beleg: `../observations/BEO-PGC/d-migrate-nacharbeit/evidence/slice-036.md`.
- **Rollenwahl:** `cdc_admin` (`ADR-0047`-Verwaltungspfad) — `EXECUTE` wird
  zuerst explizit `REVOKE … FROM PUBLIC` (PostgreSQL vergibt `EXECUTE` auf
  neue Funktionen sonst standardmäßig an `PUBLIC`) und danach gezielt
  `GRANT … TO cdc_admin`. Real geprüft: ein Login ohne
  `cdc_admin`-Mitgliedschaft scheitert an `SELECT
  cdc.enable_table(...)` mit „permission denied for function", ein Login
  mit Mitgliedschaft gelingt.
- **Security-Modell:** `SECURITY DEFINER` mit gepinntem `SET search_path =
  cdc, pg_temp` — dieselbe Definer-Semantik, die die drei Lese-Views bereits
  implizit tragen (`nacharbeit-roles.sql`); die Funktion schreibt mit den
  Rechten ihres Eigentümers (des Rollout-Läufers), `cdc_admin` braucht kein
  zusätzliches `INSERT`-Grant auf `cdc.administration_request`.
- **Kanal-Name:** `cdc_administration` — real über `LISTEN`/`pgx`
  `WaitForNotification` bestätigt
  (`administrationrequest_test.go`).
