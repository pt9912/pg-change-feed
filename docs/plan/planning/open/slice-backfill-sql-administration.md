# Slice backfill-sql-administration: SQL-Auslösung — Antragsart `backfill`, `cdc.backfill_table`, Worker, Start-Abgleich, `cdc.backfill_status`, `diagnose`, Idempotenz-Guard

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (Administration über SQL), [`LH-FA-SST-003`](../../../../spec/lastenheft.md)
(CLI-Diagnose), [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Auslösung: „ein Backfill wird ausgelöst";
Negative: Neubeginn nach Abbruch), [`LH-FA-CON-004`](../../../../spec/lastenheft.md) (Positionen wandern nur
vorwärts — die Sichtbarkeits-Grenze), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 5 (Auslösung,
Administration, Sichtbarkeit) und Festlegung 1/2/3, [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue und Live-Reload), [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) (Schemamigrationen — Idempotenz-
Guard), [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md) (Rollen), [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1/2/3 (Annahme in einer
Transaktion, Aufnahme beim Start und bei Wecksignal, Warn-Spalten in View und
`diagnose`).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md) (Antrags-Datensatz, durch
`spec-nachzug`), [`SPEC-029`](../../../../spec/pflichtenheft.md) (Run-Zustand), [`ARC-005`](../../../../spec/architecture.md), [`ARC-007`](../../../../spec/architecture.md) — gelesen,
nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Backfill ist über SQL auslösbar, läuft im Prozess und ist
sichtbar. Drei Teile:

- **Antragsweg.** Die Antragsart `backfill` (`model.AdministrationRequestKind`,
  Store-Abbildung), die geschlossene `request_kind`-Menge mit fünf Werten
  (CHECK in `tools/schema/nacharbeit-administration.sql` — dort steht die Menge
  außerhalb des neutralen Modells, weil d-migrate eine CHECK-Änderung an
  einer bestehenden Tabelle nicht konvergiert), die SQL-Funktion
  `cdc.backfill_table(p_source_id, p_schema_name, p_table_name) RETURNS text`
  (`SECURITY DEFINER`, gepinnter `search_path`, `REVOKE … FROM PUBLIC`,
  `GRANT EXECUTE` an `cdc_admin`), die **nur** einen Antrag schreibt und
  `pg_notify` sendet, und der Eintrag der neuen Funktion in
  `knownForeignObjects` (`tools/schema/rolloutguard/guard.go`) — sonst bricht
  ein zweiter Rollout mit Exit 8.
- **Verarbeitung.** Ein Zweig `backfill` in `applyAdministrationRequest`
  (`internal/bootstrap/wiring.go`) ruft `BackfillTableUseCase.Request`:
  Vorbedingungen, die beim Antrag **geschätzte** Zeilenzahl und `Admit` — Run-Zeile
  `queued` und Antrag `applied` im Sinn von „angenommen" in **einer** Transaktion
  über den Pool der Administrations-Goroutine. Nach einem erfolgreichen `Admit`
  sendet der Zweig ein nicht blockierendes Wecksignal (Kapazität 1, Signale
  verschmelzen) an **einen** Backfill-Worker (eine Goroutine, ein Run zugleich;
  die Administrations-Goroutine blockiert **nicht** auf der Kopie) mit eigenem
  Pool aus `CDC_CAPTURE_DSN`. Die Worker-Schleife arbeitet **zuerst** alle
  `queued`-Zeilen ihrer Quelle in der Reihenfolge `(requested_at, run_id)` ab und
  wartet **dann** auf das nächste Signal ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2). Beim Prozessstart
  läuft zuerst der Bindungsaufbau der Composition Root, dann der Abgleich
  `running` → `interrupted` je `running`-Run der Quelle, dann startet der Worker
  und liest die `queued`-Zeilen; eine `queued`-Zeile überlebt einen Neustart und
  wird ausgeführt, ein `interrupted`-Run wird **nicht** aufgenommen (**kein**
  automatischer Neustart).
- **Sichtbarkeit.** Die View `cdc.backfill_status` (letzter Run je Tabelle:
  Status, Zeilen, Zeiten, Fehlertext, **geschätzte** Zeilenzahl, die
  zwei Warn-Spalten aus `run-store`) mit `SELECT` für `cdc_reader`, die Ausgabe in
  `diagnose` (`bootstrap.Diagnose`: je Tabelle Status und Fortschritt, die
  die zwei Warn-Spalten — `false`, solange keine Auswertung sie setzt —, eine unbekannte
  Schätzung als „unbekannt"; ein `failed`/`interrupted`-Run ist Berichtsinhalt,
  kein Befehlsfehler) und der Handbuch-Abschnitt der Oberfläche.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Auswertung der Warnungen** (Toleranz, Richtgröße) — sie entsteht mit
  `bench-richtgroesse`; hier zeigen View und `diagnose` die zwei Warn-Spalten (`false`),
  Schätzung und Laufzeit **ohne** Schwellenwert, und die Schätzung trägt an jedem
  Träger das Wort „geschätzt" (`BEO-PGC/geschaetzter-wert-als-grenze`).
- **Die gemessene Startposition eines frisch registrierten Consumers** —
  ungeprüft ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 1: „erwartet, nicht geprüft"); das Handbuch
  dieses Slice trägt sie nicht, `e2e` misst sie und trägt sie nach.
- **HTTP-Endpunkt und CLI-Auslösung** — Welle §6; ein Auslöser ohne SQL-Zugang
  ist ein eigener Slice über denselben Use Case.
- **Der E2E-Beleg am Container** — `e2e`; dieser Slice belegt netzlos (Fakes) und
  im Store-Tier.
- **Regelauswertung/Transformationen** — der Antragsweg der Transformationen
  erweitert die hier gesetzte `request_kind`-Menge additiv (Welle §5, K3).

## 2. Definition of Done

- [ ] Antragsweg: `cdc.backfill_table` schreibt ausschließlich einen Antrag der
      Art `backfill` (`column_name` leer) und sendet `pg_notify`; `PUBLIC` hat kein
      `EXECUTE`, ein Login ohne `cdc_admin`-Mitgliedschaft scheitert mit
      „permission denied for function"; die `request_kind`-Menge trägt genau die
      fünf Werte; der Idempotenz-Guard kennt die Funktion. *Zu belegen durch:*
      `make test-store`, `make schema-rollout` zweimal hintereinander (Exit 0),
      `tools/harness/run-schema-rollout-guard-test.sh` (alle Läufe) und der
      Unit-Test in `tools/schema/rolloutguard/guard_test.go`; `plan.yaml` und
      `down.sql` regeneriert, falls der Rollout sie verändert. Der
      **Alt-Tag-Lauf** desselben Skripts
      ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
      Entscheidung 7: das Schema des jüngsten `v*`-Tags per `git archive`
      ausrollen, danach den Arbeitsbaum — Exit 0 zweimal, der zuvor eingefügte
      Datenstand über `cdc.changes` lesbar; der Bericht nennt den Tag und die
      gedruckten Exit-Codes) trägt den Upgrade-Beleg für die CHECK-Menge und die
      Funktion in `nacharbeit-administration.sql` und für die neue View über
      einen Alt-Bestand.
- [ ] Verarbeitung: ein `backfill`-Antrag gegen eine aktivierte Tabelle ruft
      `Request`, legt den Run `queued` mit Schätzung an und vermerkt ihn `applied`
      („angenommen") in einer Transaktion, weckt den Worker und blockiert die
      Administrations-Goroutine nicht; eine nicht aktivierte Tabelle und ein
      zweiter Antrag bei aktivem Run enden als `failed` mit Text; der Worker führt
      Runs nacheinander in Antragsreihenfolge `(requested_at, run_id)`; nach einem
      Prozessneustart ist ein vorheriger `running`-Run `interrupted` und ein
      `queued`-Run wird beim Start aufgenommen und ausgeführt; ein `interrupted`-Run
      startet nicht von selbst; ein Signal, das während eines Runs eintrifft, geht
      nicht verloren; fehlen Bindung oder Publication-Mitgliedschaft bei der
      Aufnahme, endet der Run `failed` (Klasse `configuration`) ohne Slot. *Zu
      belegen durch:* `make test` (Whitebox in `internal/bootstrap`, mit Fakes; je
      Regel ein Test mit Mutation) und ein Test im Store-Tier für den Abgleich und
      für zwei aufeinanderfolgende `backfill`-Anträge derselben Tabelle (der zweite
      endet `failed`; das trägt die Annahme von [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1, dass
      Anträge sequenziell in **einer** Goroutine verarbeitet werden).
      **Aufrufer-Vertrag von `Execute`** (Übergabe aus `slice-backfill-run-usecase`,
      Verifikation V-1): das Ergebnis trägt den Run im Endzustand **des Aufrufs**;
      die Zeile in `cdc.backfill_run` ist der dauerhafte Zustand und kann `completed`
      tragen, während das Ergebnis `failed` meldet (Commit mit unbekanntem Ausgang;
      das Wecksignal des Use Cases entfällt dann, obwohl die Daten sichtbar sind).
      Der Worker verwendet das Ergebnis für Log und Fehlerklasse, schreibt daraus
      keinen weiteren Endzustand und stellt keinen neuen Antrag;
      `cdc.backfill_status` und `diagnose` lesen die Zeile. Der Test-Beleg ist ein
      Whitebox-Fall des Worker-Tests (ein Fake-Use-Case liefert ein Ergebnis
      `failed`; der Worker schreibt keinen Endzustand und beantragt nichts).
      **Bezug von Antrag und Run** (Übergabe aus `slice-backfill-run-store`,
      Review F-4): `Admit` nimmt einen `pending`-Antrag über seine Kennung an und
      liest weder die Antragsart noch Quelle, Schema und Tabelle des Antrags; die
      Verarbeitung prüft die Antragsart `backfill` und dass Quelle, Schema und
      Tabelle des Antrags die des Runs sind, bevor sie `Admit` ruft (oder der
      Vermerk trägt beide Bedingungen in seiner `WHERE`-Klausel). *Zu belegen
      durch:* je Abweichung (Art, Quelle/Schema/Tabelle) ein Test mit Mutation
      (Prüfung entfernen → rot).
- [ ] Sichtbarkeit: `cdc.backfill_status` liefert je Tabelle den letzten Run mit
      Status, Zeilen, Zeiten, Fehlertext, der als geschätzt geführten
      Zeilenzahl (`NULL`/unbekannt bleibt „unbekannt", nie `0`) und der
      zwei Warn-Spalten — **schon in dieser Signatur** (`warn_estimated_size`,
      `warn_duration`, die `false`-Werte der Run-Zeile), damit
      `slice-backfill-bench-richtgroesse` die View nur mit Werten füllt und ihre
      Signatur nicht erneut ändert (Auflage aus dem Architect-Verdikt
      `architect-verdict-schema-rollout-view-signatur`; eine Signaturänderung
      einer bestehenden View kostet nach
      [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) ein
      Lesefenster im Rollout); `cdc_reader` liest sie, die Basistabelle nicht; `diagnose`
      gibt sie aus (beide Warn-Spalten `false`, solange keine Auswertung sie setzt),
      ein `failed`-Run ist Berichtsinhalt (Exit 0). Das Benutzerhandbuch trägt den neuen Abschnitt
      „Bestand als Backfill überführen" (Auslösung, Betriebs-Vorbedingungen —
      `SELECT`-Recht der Login-Identität von `CDC_CAPTURE_DSN` auf die
      Quelltabelle, Reserve in `max_replication_slots`/`max_wal_senders`,
      Snapshot-Haltedauer an der Quelle —, Sichtbarkeits-Grenze, Lesen ohne
      `Limit`, Status und Diagnose), die Rollen-Beschreibung, das Glossar, der
      `Version:`-Kopf und die Änderungshistorie. Die zwei Fortsetzungs-Idiome des
      Handbuchs tragen die Regel „Position und `limit`" aus [`SPEC-022`](../../../../spec/pflichtenheft.md):
      `commit_position > <letzte-gelesene-position>` mit `LIMIT` im SQL-Beispiel
      unter „Änderungen lesen" und `from = <letzte gelieferte commit_position> + 1`
      unter „Changes lesen". *Zu belegen durch:* `make test-store` (View, Grant),
      `make test` (`diagnose`), Review des Handbuchs.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: siehe dritter Liefer-Punkt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/model/administrationrequest.go` (+ Test) | update | Antragsart `backfill`; der Doc-Kommentar zählt die Menge auf. |
| `internal/adapters/driven/postgresstorage/administrationrequest.go` (+ Test) | update | Abbildung der Antragsart. |
| `internal/adapters/driven/postgresstorage/backfilladmission.go`, `queries/queries.go` (+ Test) oder die Verarbeitung in `internal/bootstrap/wiring.go` | update (Übergabe aus `slice-backfill-run-store`) | Art- und Bezugs-Prüfung von Antrag und Run vor bzw. bei `Admit` (DoD „Bezug von Antrag und Run"). |
| `tools/schema/nacharbeit-administration.sql` | update | CHECK-Menge (fünf Werte), Funktion `cdc.backfill_table`, Kopfkommentar. |
| `tools/schema/rolloutguard/guard.go` (+ `guard_test.go`) | update | Eintrag der Funktion; der Kommentar „aktuell sechs Objekte" zählt neu. |
| `tools/schema/schema.yaml` | update | View `backfill_status` im neutralen Modell mit ihrer endgültigen Spaltenliste, die zwei Warn-Spalten eingeschlossen (Ausweichform: Nacharbeit-SQL, dann Guard-Eintrag). |
| `tools/schema/nacharbeit-roles.sql` (+ `roles_rollout_file_internal_test.go`) | update | `SELECT` auf die View für `cdc_reader`. |
| `internal/bootstrap/wiring.go` (+ Tests) | update | Zweig `backfill` (ruft `Request`, sendet das Wecksignal), Worker-Goroutine mit Start-Aufnahme und Schleife „erst abarbeiten, dann warten", Pool, Start-Reihenfolge (Bindungsaufbau, Abgleich, Worker), `Diagnose`-Ausgabe; der Snapshot-Adapter (`postgressnapshot.New`, aus `slice-backfill-snapshot-reader`) entsteht hier mit `CDC_CAPTURE_DSN` und seinen Startwerten für Blockgröße und Zeitlimit (Setzung ohne Messung) — auch die Schätzung im Antrag liest er über diesen DSN, nicht über den Pool der Administrations-Goroutine. |
| `tools/harness/run-schema-rollout-guard-test.sh` | prüfen | trägt die Läufe des Guards; ein Lauf gegen die neue Funktion; der Alt-Tag-Lauf ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7) wird ausgeführt, nicht geändert. |
| `docs/user/benutzerhandbuch.md` | update | neuer Abschnitt, §2 Rollen, §4 Diagnose, Glossar, die zwei Fortsetzungs-Idiome (§4 „Änderungen lesen", „Changes lesen"), `Version:`-Kopf und Änderungshistorie. |
| `harness/README.md` §Sensors | update | Zeile `make example-demo-up`, die die Fremdobjekt-Aufzählung wiederholt; die Zeile `make schema-rollout` trägt keine Aufzählung, sie verweist auf `harness/targets/schema-rollout.md`. |
| `harness/targets/schema-rollout.md` | update | trägt die Zahl und die Objektklassen der Fremdobjekte; der Makefile-Kommentar über `schema-rollout` trägt sie nicht. |

**Übergaben aus `slice-backfill-run-usecase`** (gemeldet, kein zusätzlicher Umfang; die
Ports liegen in `internal/application/port/outbound/backfill*.go`, der Inbound Port in
`internal/application/port/inbound/backfill.go`, der Use Case in
`internal/application/usecase/backfill/service.go`):

- **Pflicht-Ports.** `Ports` des Use Cases trägt acht Pflicht-Ports, darunter
  `SchemaStorePort` (`model.NewChange` verlangt eine Schema-Version-Referenz); die
  Verdrahtung in `wiring.go` reicht ihn durch. Optional sind `WithChangeNotification`
  ([`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md)) und `WithLog`.
- **Publication im Command.** `BackfillRequestCommand` und `BackfillExecuteCommand`
  tragen die Publication der Quelle; der Antragszweig und der Worker übergeben sie
  (wie bei `EnableTableCommand`).
- **Ergebnis von `Execute`.** Ein Run-Fehler ist ein Ergebnis mit dem Run im
  Endzustand, kein Fehler des Aufrufs; der Fehler des Aufrufs meldet „Endzustand nicht
  festgehalten" oder „Kontext endete vor dem Beginn" (ein noch `queued` Run bleibt dann
  `queued`). Der Aufrufer-Vertrag steht im DoD „Verarbeitung".
- **Aufnahme und Abgleich.** `Queued` (Aufnahme) und `InterruptRunning` (Start-Abgleich)
  sind Operationen des Run-Zustands-Ports; der Use Case ruft keine von beiden, Worker
  und Prozessstart tragen sie.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die geschlossene `request_kind`-Menge (vier → fünf)", „die Menge der Fremdobjekte außerhalb des neutralen Modells (sechs → sieben)", „die Ausgabe von `diagnose`"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Aufzählungen der Antragsarten | `grep -rn 'exclude_column' internal tools docs spec harness` | *(Implementer trägt ein)* | Fehlertext in `applyAdministrationRequest` („… geschlossene Menge enable/disable/exclude_column/include_column"), Doc-Kommentare, Handbuch, Architektur-Sicht |
| Zahl der Fremdobjekte („sechs") | `grep -rn 'sechs' harness Makefile tools docs/user` | *(Implementer trägt ein)* | Guard-Kommentar, `harness/targets/schema-rollout.md`, die `harness/README.md`-Zeile `make example-demo-up`, Sensor-Dateien nachziehen; `Accepted` ADRs nicht ändern |
| Läufe des Guard-Tests | Lesen von `tools/harness/run-schema-rollout-guard-test.sh` und `harness/targets/schema-rollout.md` §Belege | *(Implementer trägt ein)* | Zahl und Beschreibung nachziehen, falls sich die Läufe ändern |
| Beispielausgabe von `diagnose` im Handbuch (§4) und `diagnose`-Tests | `grep -rn 'diagnose' docs/user internal/bootstrap` | *(Implementer trägt ein)* | Beispiel und Tests an die neue Ausgabe |
| Zahl der Replication-Verbindungen je Container-Lauf (Handbuch §Grenzwerte: „zwei gleichzeitige Replication-Protokoll-Verbindungen … zählen gegen `max_wal_senders`") — während eines Runs kommt kurzzeitig der Walsender des temporären Slots hinzu | `grep -rn 'max_wal_senders\|Replication-Protokoll-Verbindungen' docs/user spec` | *(Implementer trägt ein)* | Aussage präzisieren („zwei; während der Slot-Anlage eines Backfills eine weitere"), sobald sie belegt ist |
| Rollen-Beschreibung im Handbuch (§2, Rollen-Tabelle) | Lesen | *(Implementer trägt ein)* | `cdc_capture` trägt zusätzlich die Betriebs-Vorbedingung `SELECT` auf Quelltabellen |
| Fortsetzungs-Idiome und `backfill`-Nennung im Handbuch | `grep -n 'letzte\|LIMIT\|backfill' docs/user/benutzerhandbuch.md` | *(Implementer trägt ein)* | Übergabe aus dem Spec-Nachzug, am Stand `89053d3b` nachgemessen: `commit_position > <letzte-gelesene-position>` samt `LIMIT 500` (SQL-Beispiel unter „Änderungen lesen", Z. 379/381) und `from = <letzte gelieferte commit_position> + 1` (unter „Changes lesen", Z. 642) tragen die Regel „Position und `limit`" aus [`SPEC-022`](../../../../spec/pflichtenheft.md) nicht; `grep -c 'backfill' docs/user/benutzerhandbuch.md` liefert 0 (das Handbuch nennt `backfill` nicht). Zeilennummern sind der Stand dieser Messung, am Start neu messen |
| E2E-Abdeckungs-Zeilennummern | `git diff --stat` auf `test/integration/**`, `tools/harness/run-integration-tests.sh` | *(Implementer trägt ein)* | dieser Slice berührt den Runner nicht; ein Treffer wäre ein Plan-Nachzug |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `run-usecase` und `run-store` in
`done/` liegen und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **erwartet als
  Möglichkeit** — dies ist der größte Slice der Welle. Sprengt der Diff die
  Größe eines Reviews, ist der abtrennbare Teil der dritte Liefer-Punkt
  (Status-View, `diagnose`, Handbuch) als eigener Slice
  `backfill-status-diagnose` mit Start nach diesem; die Welle nimmt ihn dann in
  §4 auf.
- `in-progress` → `open` (blockiert — Carveout?): falls die View im neutralen
  Modell nicht konvergiert (dann Nacharbeit-SQL und ein zweiter Guard-Eintrag)
  oder der Idempotenz-Guard die neue Funktion nicht erkennt (Exit 8) — beides
  Architect-Fragen nach [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test`, `make test-store` und
`make schema-rollout` (zweimal) real grün + der Alt-Tag-Lauf von
`tools/harness/run-schema-rollout-guard-test.sh` real grün + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Idempotenz-Guard erkennt die neue Funktion nicht** — ein zweiter
  Rollout bräche mit Exit 8. Die Signatur `backfill_table(in:text,in:text,in:text)`
  entspricht der Form von `enable_table` (gelesen in `guard.go`), die Erkennung
  der neuen ist ungemessen. *Erwartet, zu belegen durch:* zweiter
  `make schema-rollout` und `run-schema-rollout-guard-test.sh`. **Ausgang:**
  *(bei Closure)*
- **Die neue View oder die erweiterte CHECK-Menge konvergiert nicht über einen
  Alt-Bestand.** Eine neue View konvergiert über einen Alt-Bestand (gemessen im
  Architect-Verdikt `architect-verdict-schema-rollout-view-signatur`,
  Szenario 5: Exit 0, zweiter Lauf Exit 0); CHECK-Menge und Funktion laufen über
  `nacharbeit-administration.sql` und sind über einen Alt-Bestand ungemessen.
  *Erwartet, zu belegen durch:* der Alt-Tag-Lauf von
  `run-schema-rollout-guard-test.sh`. **Ausgang:** *(bei Closure)*
- **`applied` wird als „Bestand kopiert" gelesen.** Bei dieser Antragsart heißt
  `applied` „angenommen"; die Ausführung steht in `cdc.backfill_run`. Ein Leser
  von `cdc.administration_request`, `diagnose` oder Handbuch könnte es anders
  deuten. *Erwartet, zu belegen durch:* [`SPEC-019`](../../../../spec/pflichtenheft.md) und Handbuch sagen es
  ausdrücklich; Review liest beide. **Ausgang:** *(bei Closure)*
- **Die geschätzte Zeilenzahl wird zur Grenze** (`BEO-PGC/geschaetzter-wert-als-grenze`,
  offen, 1×): `pg_class.reltuples` ist eine Schätzung, für nie analysierte
  Tabellen als „unbekannt" erwartet (`-1`). Weder Handbuch noch Code dürfen sie
  als Grenze oder als `0` lesen; der Slice führt kein Ablehnen. *Erwartet, zu
  belegen durch:* Tests mit `-1` und die Wortwahl an jedem Träger. **Ausgang:**
  *(bei Closure)*
- **Ein `queued`-Run überlebt den Prozessstart nicht**, wenn der Worker nur einen
  In-Speicher-Kanal liest (`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger`,
  verkörpert): der Träger ist die Tabelle, das Signal weckt nur. *Erwartet, zu
  belegen durch:* der Neustart-Test und der Test „Signal während eines Runs" des
  zweiten Liefer-Punkts. **Ausgang:** *(bei Closure)*
- **Die Administrations-Goroutine blockiert auf einer Kopie.** *Erwartet, zu
  belegen durch:* ein Test mit einem blockierenden Fake-Run, während ein
  weiterer Antrag verarbeitet wird. **Ausgang:** *(bei Closure)*
- **Das Handbuch nennt eine Aussage ohne Beleg** (Startposition eines frisch
  registrierten Consumers; Richtgröße): beide bleiben bis `e2e` bzw.
  `bench-richtgroesse` **außerhalb** des Handbuchs ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B).
  *Erwartet, zu belegen durch:* Review des Abschnitts. **Ausgang:** *(bei
  Closure)*
- **Handbuch und Änderungshistorie** (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  und `BEO-PGC/handbuch-versionshistorie-uebersprungen`, verkörpert, je 3×): der
  Slice liefert mehrere Oberflächen (Funktion, View, `diagnose`); jede braucht
  ihre Stelle. *Erwartet, zu belegen durch:* der Suchlauf-Eintrag und die neue
  Historien-Zeile. **Ausgang:** *(bei Closure)*
- **Eine Instanz je Quelle** trägt den Start-Abgleich; laufen zwei Instanzen gegen
  dieselbe Quelle, setzte die zweite den Run der ersten auf `interrupted`.
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) legt „eine Instanz je Quelle" fest ([`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)-Muster für
  `pending`); die Annahme ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1) und die Aufnahme (Festlegung 2)
  stützen sich auf dieselbe Festlegung; ein Schutz ist nicht Teil. **Ausgang:**
  *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Bootstrap, Schema-Skripte und Handbuch sind keine
eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/geschaetzter-wert-als-grenze` (offen, 1×, einschlägig — Risiko §6),
`BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×) und
`BEO-PGC/d-migrate-nacharbeit` (verkörpert, 6×) — Guard-Eintrag und
Rückführung §4, `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger`
(verkörpert, Risiko §6), `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
und `BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert, je 3×,
Risiko §6), `BEO-PGC/rollen-test-abdeckungsluecken` (offen, 2×, einschlägig —
Rollen-Test der View), `BEO-PGC/limit-fortsetzung-innerhalb-einer-position`
(offen, 0×, einschlägig — das Handbuch trägt die Lese-Regel „Bestandsabzug ohne
`Limit`"; die Cursor-Form ist nicht Teil), `BEO-PGC/kein-admin-weg-schema-fehler-recovery`
(offen, 1×, gesichtet — ein fehlgeschlagener Run ist über Status und `diagnose`
sichtbar, ein Recovery-Weg für Schema-Fehler ist nicht Teil),
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (offen, 2×, Kommentare
zum Worker-Fehlerpfad sagen nur zu, was der Code trägt),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
