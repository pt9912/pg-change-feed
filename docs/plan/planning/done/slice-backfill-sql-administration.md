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

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md) (Antrags-Datensatz) —
**geändert**: der Absatz „Grants“ der Antrags-Queue
([Architect-Verdikt](../../../reviews/architect-verdict-backfill-schema-klasse-rollen.md));
[`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) Absatz „Markierung“ — **geändert**: ein Satz zur
Schema-Version ([`ADR-0116`](../../adr/0116-backfill-schema-version-referenz-reichweite.md));
[`SPEC-029`](../../../../spec/pflichtenheft.md) (Run-Zustand), [`ARC-005`](../../../../spec/architecture.md),
[`ARC-007`](../../../../spec/architecture.md) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-24.

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

- [x] Antragsweg: `cdc.backfill_table` schreibt ausschließlich einen Antrag der
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
      einen Alt-Bestand; Lauf 5 prüft dafür im Skript die Rechte der drei Rollen,
      `EXECUTE` auf die Funktion und die CHECK-Menge.
- [x] Verarbeitung: ein `backfill`-Antrag gegen eine aktivierte Tabelle ruft
      `Request`, legt den Run `queued` mit Schätzung an und vermerkt ihn `applied`
      („angenommen") in einer Transaktion, weckt den Worker und blockiert die
      Administrations-Goroutine nicht; eine nicht aktivierte Tabelle und ein
      zweiter Antrag bei aktivem Run enden als `failed` mit Text; der Worker führt
      Runs nacheinander in Antragsreihenfolge `(requested_at, run_id)`; nach einem
      Prozessneustart ist ein vorheriger `running`-Run `interrupted` und ein
      `queued`-Run wird beim Start aufgenommen und ausgeführt; ein `interrupted`-Run
      startet nicht von selbst; ein Signal, das während eines Runs eintrifft, geht
      nicht verloren; fehlen Bindung oder Publication-Mitgliedschaft bei der
      Aufnahme, endet der Run `failed` (Klasse `configuration`) ohne Slot; ein
      `backfill`-Antrag einer anderen Quelle bleibt für diese Instanz `pending`. *Zu
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
- [x] Sichtbarkeit: `cdc.backfill_status` liefert je Tabelle den letzten Run mit
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
- [x] Spec-Zug (Architect-Verdikt): [`SPEC-019`](../../../../spec/pflichtenheft.md) trägt den Absatz „Grants“ der
      Antrags-Queue (`cdc_admin` `SELECT`, `UPDATE`; `cdc_capture` und `cdc_reader` kein Recht;
      niemand `INSERT` oder `DELETE`); [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) Absatz „Markierung“ trägt den
      Satz zur Schema-Version, das Handbuch denselben Inhalt im Abschnitt zum Backfill
      ([`ADR-0116`](../../adr/0116-backfill-schema-version-referenz-reichweite.md) Folgepflicht 2 und 3). *Zu belegen durch:* der Diff der
      Spec und des Handbuchs, `make docs-check`; die Rechte selbst der Rollen-Test
      (`make test-store`).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8);
      kein offenes HIGH oder MEDIUM nach der Fixrunde (F-1 bis F-6 behoben, F-7 bis F-9 behoben, F-10 als
      Quellvergleich umgesetzt, F-11 bis F-14 ohne Aktion).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: siehe dritter Liefer-Punkt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

**Umfangsentscheidung (vor dem Code, am Diff gemessen):** keine Rückführung nach `next/`. Der in §4 vorab benannte abtrennbare Teil (Status-View, `diagnose`, Handbuch) ist klein — die View ist eine YAML-Stelle samt einem Grant, `diagnose` eine Funktion mit Formatierung —, und drei Bindungen halten ihn im Slice: die Signatur der View soll von Anfang an stehen (Auflage aus dem Architect-Verdikt zur View-Signatur), die Handbuch-Pflicht für `cdc.backfill_table` verlangt den Abschnitt ohnehin, und der Login-Test der Verarbeitung liest den Status. Gemessen am Stand nach dem dritten Produktions-Commit (`git diff --shortstat 2d47d8a7 HEAD`): 33 Dateien, 2106 Zeilen dazu, davon 413 in acht Produktions-Go-Dateien und 1377 in 14 Testdateien; der Rest Spec, Handbuch, Schema-SQL und Harness-Doku. Vier Commits trennen Speicher, Verdrahtung, Spec und Handbuch. Die Spec-Zusätze aus dem Architect-Verdikt sind zwei kurze Absätze und tragen den Umfang nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/model/administrationrequest.go` (+ Test) | update | Antragsart `backfill`; der Doc-Kommentar zählt die Menge auf. |
| `internal/adapters/driven/postgresstorage/administrationrequest.go` (+ Test) | update | Abbildung der Antragsart. |
| `internal/adapters/driven/postgresstorage/backfilladmission.go`, `queries/queries.go` (+ Test) oder die Verarbeitung in `internal/bootstrap/wiring.go` | update (Übergabe aus `slice-backfill-run-store`) | Art- und Bezugs-Prüfung von Antrag und Run vor bzw. bei `Admit` (DoD „Bezug von Antrag und Run"). **Umgesetzt** in der zweiten Form des Plans („der Vermerk trägt beide Bedingungen in seiner `WHERE`-Klausel“): die neue Anweisung `UpdateAdministrationRequestAdmitted` trifft nur einen `pending`-Antrag der Art `backfill` mit Quelle, Schema und Tabelle des Runs; trifft sie keine Zeile, unterscheidet `rejectedRequest` (Antrag fehlt oder nicht offen → `ErrBackfillRequestNotPending`; offen, aber andere Art oder Adresse → `ErrBackfillRunInvalid`). Der Antragszweig baut den Run aus den Werten des Antrags. Port-Kommentar an `BackfillAdmissionPort.Admit` nachgezogen. |
| `internal/adapters/driven/postgresstorage/backfillhelpers_test.go`, `backfilladmission_test.go` | update | Der Helfer `pendingRequest` legt die Art `backfill` an (`pendingRequestOfKind` für abweichende Arten); neuer Test je Abweichung (Art, Quelle, Schema, Tabelle) als Unterfall. |
| `internal/domain/errors/errors.go`, `internal/application/port/outbound/administrationrequest.go` | update (nicht im Plan) | Doc-Kommentare zählen die fünf Antragsarten bzw. die fünf Funktionen (§3.13-Suchlauf). |
| `internal/bootstrap/backfill.go` (neu) | create (nicht im Plan als eigene Datei) | Worker (`runBackfillWorker`, `drainBackfillQueue`), Wecksignal (`newBackfillWake`, `signalBackfillWorker`), Start-Abgleich (`reconcileBackfillRuns`) und die Ausgabe von `diagnose` (`diagnoseBackfillStatus`); `wiring.go` bleibt bei der Verdrahtung. Der Worker liest nach einem Durchgang mit Fehler nach `backfillRetryInterval` erneut (ein Lesefehler bei der Aufnahme lässt eine `queued`-Zeile sonst bis zum nächsten Signal liegen; kein Poll im Normalbetrieb). |
| `internal/bootstrap/backfill_internal_test.go`, `backfill_endtoend_test.go` (neu), `administration_roles_internal_test.go`, `diagnose_test.go`, `roles_rollout_file_internal_test.go` | create / update | Whitebox-Tests mit Fakes (Reihenfolge, Aufnahme beim Start, Signal zwischen Lesung und Warten, erneutes Lesen nach jedem Run, nicht blockierender Sender, Ergebnis `failed` ohne weiteren Endzustand, Wiederholung nach Fehler, Antragszweig, Abgleich); der Abgleich gegen die reale Run-Tabelle und die Aufnahme einer `queued`-Zeile ohne Bindung (Worker, Use Case und Adapter real: der Run endet `failed`, Klasse `configuration`, vor dem Slot); der Login-Test führt die Art `backfill` (angenommen unter `cdc_admin`/`cdc_capture`, zweiter Antrag `failed`, nicht aktivierte Tabelle `failed`); `diagnose` unter einer `cdc_reader`-Identität; die neue View im Reader-Grant des Rollout-Texts. |
| `internal/adapters/driven/postgresstorage/administrationrequest_test.go`, `roles_test.go`, `sqlviews_test.go`, `backfillstatusview_test.go` (neu) | update / create | die Funktion `cdc.backfill_table` (schreibt nur den Antrag, `pg_notify`, `EXECUTE` nur `cdc_admin`), die CHECK-Menge mit genau fünf Werten, `cdc_reader` liest über die View und nicht die Basistabelle, die View trägt den letzten Run je Tabelle (Ordnung, Gleichstand, NULL/0, Warn-Spalten), und der Schlüsselvergleich der Fortsetzung in einer Position (Beleg der Handbuch-Aussage). |
| `spec/pflichtenheft.md` | update (Architect-Verdikt) | `SPEC-019`: Absatz „Grants“; `LH-FA-CAP-009.a` „Markierung“: Satz zur Schema-Version; Historie-Zeile. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | regeneriert | Ergebnis des Rollouts gegen eine leere Datenbank (die View `backfill_status` erhöht die Operationszahl von 15 auf 16); die `target`-Zeile bleibt die committete. |
| `tools/harness/run-schema-rollout-guard-test.sh` | update (Fixrunde zum Review, F-1/F-7) | Lauf 5 (Alt-Tag) prüft im Skript die Zusagen aus `harness/targets/schema-rollout.md` §Belege: 19 Tabellen-/View-Rechte der drei Rollen auf `cdc.administration_request`, `cdc.backfill_run`, `cdc.backfill_status`, `EXECUTE` auf `cdc.backfill_table` allein für `cdc_admin` (auch nicht `PUBLIC`) und die fünf Werte von `chk_administration_request_kind`; Vorbedingungen am Alt-Stand (Funktion fehlt, CHECK ohne `backfill`, `cdc_admin` ohne `UPDATE`). Die Zahl der Fremdobjekte in Kommentaren und der Ausgabe von Lauf 3 lautet „sieben“. Eingabeseiten-Mutationen gesehen (siehe Fixrunde). |
| `tools/schema/nacharbeit-administration.sql` | update | CHECK-Menge (fünf Werte), Funktion `cdc.backfill_table`, Kopfkommentar. |
| `tools/schema/rolloutguard/guard.go` (+ `guard_test.go`) | update | Eintrag der Funktion; der Kommentar „aktuell sechs Objekte" zählt neu. |
| `tools/schema/schema.yaml` | update | View `backfill_status` im neutralen Modell mit ihrer endgültigen Spaltenliste, die zwei Warn-Spalten eingeschlossen (Ausweichform: Nacharbeit-SQL, dann Guard-Eintrag). |
| `tools/schema/nacharbeit-roles.sql` (+ `roles_rollout_file_internal_test.go`) | update | `SELECT` auf die View für `cdc_reader`. |
| `internal/bootstrap/administration_roles_internal_test.go` | update (Übergabe aus `slice-backfill-run-store`) | der Login-Test zieht je Antragsart einen Antrag durch `processAdministrationRequests` unter einer `cdc_admin`- und einer `cdc_capture`-Login-Identität; die Art `backfill` kommt hinzu, damit der neue Zweig unter dem Rollenschnitt läuft, für den er gebaut ist. |
| `internal/bootstrap/wiring.go` (+ Tests) | update | Zweig `backfill` (ruft `Request`, sendet das Wecksignal), Worker-Goroutine mit Start-Aufnahme und Schleife „erst abarbeiten, dann warten", Pool, Start-Reihenfolge (Bindungsaufbau, Abgleich, Worker), `Diagnose`-Ausgabe; der Snapshot-Adapter (`postgressnapshot.New`, aus `slice-backfill-snapshot-reader`) entsteht hier mit `CDC_CAPTURE_DSN` und seinen Startwerten für Blockgröße und Zeitlimit (Setzung ohne Messung) — auch die Schätzung im Antrag liest er über diesen DSN, nicht über den Pool der Administrations-Goroutine. |
| `tools/harness/run-schema-rollout-guard-test.sh` | ausgeführt | trägt die Läufe des Guards, darunter den Alt-Tag-Lauf ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7); Änderung siehe die Zeile oben. |
| `docs/user/benutzerhandbuch.md` | update | neuer Abschnitt, §2 Rollen, §4 Diagnose, Glossar, die zwei Fortsetzungs-Idiome (§4 „Änderungen lesen", „Changes lesen"), `Version:`-Kopf und Änderungshistorie. |
| `harness/README.md` §Sensors | update | Zeile `make example-demo-up`, die die Fremdobjekt-Aufzählung wiederholt; die Zeile `make schema-rollout` trägt keine Aufzählung, sie verweist auf `harness/targets/schema-rollout.md`. |
| `harness/targets/schema-rollout.md` | update | trägt die Zahl und die Objektklassen der Fremdobjekte; der Makefile-Kommentar über `schema-rollout` trägt sie nicht; der Vertrag nennt fünf Views einschließlich `backfill_status` (Fixrunde F-6). |
| `internal/bootstrap/diagnose_test.go` | update (Fixrunde F-2) | `TestDiagnoseReportsTheLatestBackfillRunPerTable` legt zusätzlich einen Run einer fremden Quelle an; die Ausgabe darf seine Tabelle nicht tragen. |
| `internal/bootstrap/wiring.go`, `backfill_internal_test.go`, `administration_roles_internal_test.go` | update (Fixrunde F-10, nicht im Plan) | `administrationDeps.source` (`cfg.Source`); `processAdministrationRequests` übergeht einen `backfill`-Antrag einer anderen Quelle, er bleibt `pending` für die Instanz dieser Quelle ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1). Unit-Test mit zwei Quellen; der Login-Test legt einen Antrag einer fremden Quelle vor dem eigenen an (bleibt `pending`, keine Run-Zeile, ein Signal). Keine Port-Änderung: der Vergleich steht in der Verarbeitung, nicht in der Lesung. Grenze: die vier übrigen Antragsarten lesen weiter alle offenen Anträge ohne Quellbezug (Bestand aus [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)). |
| `internal/adapters/driven/postgresstorage/administrationrequest.go` (+ Test) | update (Fixrunde F-9) | `MarkApplied` meldet „Antrag erledigt“ nur, wenn eine `pending`-Zeile vermerkt wurde; ein `backfill`-Antrag, den die Annahme vermerkt hat, bleibt ohne Meldung. |
| `internal/adapters/driven/postgresstorage/queries/queries.go`, `administrationrequest_test.go`, `internal/bootstrap/administration_endtoend_test.go`, `tools/schema/schema.yaml`, `tools/schema/nacharbeit-administration.sql` | update (Fixrunde F-6) | Zähl- und Aufzähltexte: fünf Antragsarten, drei Tabellen-Antragsarten ohne Spalte, fünf Antrags-Funktionen; „Folge-Slice“-Hinweise für die gelieferte Administrations-Goroutine entfallen. |
| `internal/adapters/driven/postgresstorage/backfilladmission_test.go`, `internal/bootstrap/administration_roles_internal_test.go`, `diagnose_test.go`, `wiring.go` | update (Fixrunde F-8) | `gofmt` (gepinntes Toolchain-Image) für die vier Dateien, die der Review meldete. |
| `docs/user/benutzerhandbuch.md` | update (Fixrunde F-4/F-5, Version 1.49) | Antrag und Vermerk unter `cdc_admin`, Lesen von `cdc.backfill_status` unter `cdc_reader` (die Rolle `cdc_admin` trägt kein `SELECT` auf die View); Stichtag des Bestands ist der Start des Runs; ein Backfill-Antrag einer anderen Quelle bleibt `pending`. |

**Fixrunde zum Review — Eingabeseiten-Mutationen** (Zusage · mutierte Eingabe · gesehenes Rot; jede Datei danach zurückgenommen):

- `diagnose` zeigt nur die Tabellen der eigenen Quelle · `WHERE source_id = $1` → `WHERE $1::text IS NOT NULL` in `internal/bootstrap/backfill.go` · rot: `TestDiagnoseReportsTheLatestBackfillRunPerTable` (`make test-store`).
- Ein `backfill`-Antrag einer fremden Quelle bleibt `pending` · der Vergleich `request.Source != deps.source` in `processAdministrationRequests` entfernt · rot: `TestAdministrationBackfillBranchLeavesTheRequestOfAnotherSourcePending` (`make test` im Race-Image) und `TestAdministrationPathRunsUnderLeastPrivilegeLogins` (`make test-store`, beide im selben Lauf mit der Mutation von `diagnose`).
- `EXECUTE` auf `cdc.backfill_table` allein für `cdc_admin` (Lauf 5 des Guard-Test-Skripts) · `cdc.backfill_table(text, text, text)` aus der `GRANT`-Zeile von `nacharbeit-administration.sql` gestrichen · rot: „Lauf 5: cdc_admin trägt nach dem Upgrade über v0.1.2 auf cdc.backfill_table(text, text, text) das Recht EXECUTE nicht als t“; die `REVOKE`-Zeile ohne die Funktion · rot: „Lauf 5: cdc_capture trägt … das Recht EXECUTE nicht als f“.
- „Antrag erledigt“ steht nur bei vermerkter Zeile · `tag.RowsAffected() > 0` → `>= 0` in `MarkApplied` · rot: `TestAdministrationRequestAdapterMarkAppliedLogsOnlyARequestItMarked` („Meldungen … nach dem zweiten MarkApplied = 2, erwartet 1“, `make test-store`).

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

**Übergaben aus `slice-backfill-run-store`** (gemeldet, kein zusätzlicher Umfang; die Adapter
liegen in `internal/adapters/driven/postgresstorage/backfill{admission,run,writer}.go`; die
Rechte sind Messungen des Verifiers, die Handbuch-Version ist am Stand `c7045f81` gelesen):

- **Rechte, die die Verarbeitung voraussetzt.** `cdc_admin` trägt `SELECT`, `UPDATE` auf
  `cdc.administration_request` (kein `INSERT`, kein `DELETE`) und `SELECT`, `INSERT` auf
  `cdc.backfill_run`; `cdc_capture` trägt `SELECT`, `UPDATE` auf `cdc.backfill_run` und **kein**
  Recht auf die Antrags-Queue; `cdc_reader` trägt kein Recht auf eine der beiden Basistabellen
  (`has_table_privilege` an PostgreSQL 17.11 und 18.6, Verifikation `verifikation-slice-backfill-run-store`
  §3.2). Der Worker-Pool (`CDC_CAPTURE_DSN`) liest die Queue also nicht; was er vom Antrag
  braucht, steht in der Run-Zeile. Die View `backfill_status` bekommt hier ihr `SELECT` für
  `cdc_reader`.
- **Pools.** Jeder der drei Adapter baut seinen Verbindungspool im Konstruktor aus der DSN, die
  der Aufrufer übergibt (`NewBackfillAdmission(ctx, dsn, …)`, `NewBackfillRun(ctx, dsn, …)`,
  `NewBackfillWriter(ctx, dsn, …)`); die Annahme gehört an `cfg.AdminDSN`, Run-Zustand und
  Schreiber an `cfg.CaptureDSN`. Run-Zustand und Schreiber laufen auf getrennten Pools, damit der
  Fortschritt eines Runs nicht auf der Verbindung der offenen Schreibtransaktion liegt.
- **Sentinel-Fehler an den Ports.** `Admit` und der Schreiber melden Ablehnungen über drei
  Sentinels, unterscheidbar per `errors.Is`: `ErrBackfillRequestNotPending` (der Antrag ist
  nicht `pending`), `ErrBackfillRunInvalid` (der Run passt nicht zu `Admit`/`Begin`),
  `ErrBackfillBlockInvalid` (Block oder Commit verletzt den Vertrag des Schreibers); `classifyError`
  bildet sie auf die Klasse `internal` ab, Treiberfehler tragen `ErrBackfillStorage` (`storage`),
  ein aktiver Run ist `domainerrors.ErrBackfillRunActive`. Was der Zweig `backfill` mit einem
  Antrag tut, dessen `Admit` einen dieser Fehler meldet, legt dieser Plan nur für den aktiven Run
  und die nicht aktivierte Tabelle fest (DoD „Verarbeitung").
- **Test-Aufbau mit Logins (Muster).** Eine Rollen-Zusage wird unter einem Login belegt, nicht
  unter dem Superuser (`CREATE ROLE … LOGIN … IN ROLE <rolle>`, kein Eigentum an `cdc`-Objekten):
  `internal/bootstrap/administration_roles_internal_test.go` (der ganze Weg von
  `processAdministrationRequests`, alle vier Antragsarten, `applied` und `failed`),
  `internal/adapters/driven/postgresstorage/backfillroles_test.go` (Annahme unter `cdc_admin`,
  Run-Zustand und Schreiber unter `cdc_capture`, je unter der Rolle des anderen mit SQLSTATE `42501`),
  Rechte-Abfragen in Lauf 5 von `tools/harness/run-schema-rollout-guard-test.sh` nach dem Upgrade
  vom jüngsten `v*`-Tag. Die Grant-Mutation („Grant streichen") färbt alle drei rot.
- **Handbuch-Stand.** Am Stand `c7045f81` trägt das Benutzerhandbuch `Version: 1.47`; der Zug dieses
  Slice ist die Zeile 1.48 der Änderungshistorie. §2 (Rollen-Tabelle, Zeile `cdc_admin`), §4
  „Schema aktualisieren" (Absatz „Rechte der drei Rollen") und §5 (Zeile `CDC_ADMIN_DSN`) nennen den
  Rechteschnitt der Queue bereits; der Zug liest sie beim Nachziehen der Rollen-Beschreibung, statt
  sie zu überschreiben.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die geschlossene `request_kind`-Menge (vier → fünf)", „die Menge der Fremdobjekte außerhalb des neutralen Modells (sechs → sieben)", „die Ausgabe von `diagnose`", „die Rechte der Rollen und die Zahl der Lese-Views von `cdc_reader` (vier → fünf)"; beide Stände gemessen — Parent `2d47d8a7` mit `git grep -n … 2d47d8a7 -- <Pfade>`, Diff-Stand mit `grep -rn …` am Stand nach dem dritten Produktions-Commit `c188be43`; Zahlen sind gedruckte Zeilenzahlen dieser Läufe, keine Werte aus dem Plan):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Aufzählungen der Antragsarten und der Antragsfunktionen | `grep -rn -e 'exclude_column' internal tools docs spec harness` | Parent 103 Zeilen, Diff-Stand 105. Träger, die die Menge oder die Funktionen aufzählen: Fehlertext in `applyAdministrationRequest`, Doc-Kommentare in `model/administrationrequest.go`, `domain/errors/errors.go`, `outbound/administrationrequest.go`, Kopf und CHECK-Zeile in `nacharbeit-administration.sql`, Kommentare in `nacharbeit-roles.sql` und `schema.yaml` (Kopf, Beschreibung der Tabelle), Zeile 4 der Objektklassen in `harness/targets/schema-rollout.md`, Funktionszahl in `administrationrequest_test.go`, Meldung in `administration_internal_test.go`. Nicht gefunden: eine Aufzählung im Handbuch (es nennt die Antragsarten je Aufgabe, nicht als Menge). Schon fünfstellig: `SPEC-019`, Tabelle „Antragsart → Inbound Port“ in `spec/architecture.md`. | alle genannten Träger im Diff nachgezogen; historische Zeilen (Historie der Spec, Beobachtungs-Belege, Pläne fremder Slices) unberührt |
| Zahl der Fremdobjekte („sechs") | `grep -rn -e 'sechs' harness Makefile tools docs/user` | Parent 19 Zeilen, Diff-Stand 11, Fixrunden-Stand (`4981766a`) 7. Verbleibende Treffer: Läufe des Guard-Tests („sechs Läufe“ in `harness/README.md`, `harness/targets/schema-rollout.md`, Kopf des Skripts — unverändert richtig), unabhängige Zahlen (`erfassung-feldliste.md`, `docs-check.md`, Handbuch „sechs Klassen“, Historie 1.18). Am Diff-Stand nicht nachgezogen und in der Fixrunde nachgezogen: `tools/harness/run-schema-rollout-guard-test.sh` Zeilen 9, 14, 44 und die Ausgabe von Lauf 3 (Zeile 171). | `guard.go` (Kommentar) und `guard_test.go`, `harness/targets/schema-rollout.md` (Zeilen 89, 151), die `harness/README.md`-Zeile `make example-demo-up` und in der Fixrunde das Skript nachgezogen |
| Läufe des Guard-Tests | Lesen von `tools/harness/run-schema-rollout-guard-test.sh` und `harness/targets/schema-rollout.md` §Belege | Die Zahl der Läufe (sechs) und ihre Beschreibung bleiben wahr; Lauf 5 prüfte am Diff-Stand die Rechte von `cdc_admin`, nicht die der Rollen `cdc_capture`/`cdc_reader`, nicht die neue Funktion und nicht die CHECK-Menge, obwohl `harness/targets/schema-rollout.md` §Belege (5) sie nennt | Fixrunde: Lauf 5 prüft alle genannten Zusagen im Skript (19 Rechte der drei Rollen, `EXECUTE` allein für `cdc_admin`, CHECK-Menge), gedruckte Zeile des Endlaufs: „Lauf 5 OK — Tag v0.1.2: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); … 19 Tabellen-/View-Rechte der drei Rollen …, request_kind-Menge backfill,disable,enable,exclude_column,include_column“ |
| Zählwörter der Antragsarten, Views und Funktionen | `grep -rn -e 'vier Antragsarten' -e 'die vier Views' -e 'vier Antrags' -e 'Folge-Slice' internal harness tools` (Parent: `git grep -n … 2d47d8a7 -- internal harness tools`) | Parent 18 Zeilen, Diff-Stand (`c188be43`) 16, Fixrunden-Stand (`4981766a`) 10. Am Diff-Stand nicht nachgezogen und in der Fixrunde behoben: `queries/queries.go` („Die vier Antragsarten … die beiden Tabellen-Antragsarten“), `harness/targets/schema-rollout.md` Zeile 12 („die vier Views“), `tools/schema/schema.yaml` (Tabellenbeschreibung: `backfill` lässt `column_name` NULL; „(Folge-Slice)“), `internal/bootstrap/administration_endtoend_test.go` („die vier Antrags-Funktionen“), `tools/schema/nacharbeit-administration.sql` („eines Folge-Slice“), `administrationrequest_test.go` („in einem Folge-Slice“). Verbleibende 10 Treffer gehören nicht zu den bewegten Eigenschaften: `Folge-Slice` als allgemeiner Begriff in `harness/README.md`/`harness/conventions.md` (Carveout, Modus-Tabelle), `internal/application/port/inbound/retention.go` (Retention), Kommentare der drei SDK-Integrationsläufer. Nicht gefunden: weitere Zählwörter „vier“ zu Antragsarten oder Views; `tools/schema/schema.yaml` Zeile 42 („Die drei Views“) benennt bewusst die drei Views des Kernlesezugriffs (`LH-FA-SST-002`), nicht die Menge aller Views. | die fünf genannten Träger nachgezogen; die übrigen unberührt |
| Stichtag des Bestands im Handbuch | `grep -n 'Zeitpunkt des Antrags' docs/user/benutzerhandbuch.md` | Fixrunden-Stand vorher 2 Zeilen (§4 Einleitung des Backfill-Abschnitts, §8 Glossar „Backfill“); der Snapshot entsteht mit dem Slot des Runs beim Beginn der Ausführung, nicht beim Antrag | beide Stellen nennen den Startzeitpunkt des Runs; nach dem Nachzug 0 Treffer |
| Beispielausgabe von `diagnose` im Handbuch (§4) und `diagnose`-Tests | `grep -rn -e 'diagnose' docs/user internal/bootstrap` | Parent 52 Zeilen, Diff-Stand 76; das Beispiel im Handbuch trägt die Zeilen der Abschnitte Betriebsstatus bis Speicherverbrauch, die Tests prüfen Teilzeichenketten (`strings.Contains`) | Beispiel und Erläuterung um „Backfill je Tabelle“ ergänzt; die bestehenden `diagnose`-Tests brauchen keine Änderung (additive Ausgabe), zwei neue Tests und ein Format-Test tragen den neuen Abschnitt |
| Zahl der Replication-Verbindungen je Container-Lauf (Handbuch §Grenzwerte) | `grep -rn -e 'max_wal_senders' -e 'Replication-Protokoll-Verbindungen' docs/user spec` | Parent 1 Zeile (Handbuch, „zwei gleichzeitige … zählen gegen `max_wal_senders`“), Diff-Stand 3 (neue Betriebs-Vorbedingung im Backfill-Abschnitt, die präzisierte Grenzwert-Zeile) | präzisiert („zwei; während der Slot-Anlage eines Backfills eine weitere, die nach dem Import endet“, abgeleitet aus dem Adapter, `snapshot.go` Paketkommentar); nicht während eines Container-Laufs mit laufendem Run gemessen |
| Rollen-Beschreibung im Handbuch (§2, Rollen-Tabelle) | Lesen | Die Zeilen `cdc_capture`, `cdc_admin` und `cdc_reader` nannten weder Backfill noch `cdc.backfill_status` | die drei Zeilen und ein Betriebs-Hinweis (`SELECT` auf Quelltabellen) nachgezogen; §4 „Schema aktualisieren“ (Absatz „Rechte der drei Rollen“) und §5 (zwei DSN-Zeilen) ebenso |
| Fortsetzungs-Idiome und `backfill`-Nennung im Handbuch | `grep -n -e 'letzte-gelesene' -e 'letzte gelieferte' -e 'LIMIT' -e 'backfill' docs/user/benutzerhandbuch.md` | Der Befehl druckt am Parent 10 Zeilen und am Diff-Stand (`c188be43`) 40 (`git grep -n … 2d47d8a7 -- docs/user/benutzerhandbuch.md` bzw. `… c188be43 …`); `-e 'backfill'` allein: Parent 7, Diff-Stand 30; am Fixrunden-Stand (`55085f76`, Handbuch 1.49) drucken die beiden Befehle 41 und 31. Beide Idiome (SQL-Beispiel unter „Änderungen lesen“, `from = <letzte gelieferte commit_position> + 1` unter „Changes lesen“) trugen die Regel „Position und `limit`“ nicht | beide tragen sie jetzt (Absatz „Fortsetzen und `LIMIT`“ mit dem Schlüsselvergleich; Satz im Absatz „Changes lesen“); der Schlüsselvergleich ist im Store-Tier gelaufen (`TestChangesViewKeysetContinuesInsideOnePosition`) |
| Nenner der Coverage-Messungen (`harness/sensors/coverage-gate.md` §Zählbasis, `harness/sensors/db-adapter-coverage.md` §Zählbasis) | Lesen; `make gates` und `make test-store` gefolgt von `make test-replication` drucken die Zahlen | Der Slice fügt Produktionscode in `internal/bootstrap` und `postgresstorage` hinzu: der gemergte DB-Nenner ist gedruckt **1025** (`DB-Adapter-Coverage: 81.76% (gedeckt 838 von 1025 Statements; Profile gemergt: store,replication)`, PostgreSQL 17 und 18), die Sensor-Datei nennt **1016**; die Unit-Quote druckt `Coverage 82.90%` im Endlauf von `make gates` (ein früherer Lauf desselben Slice: 82.80 %; Sensor-Datei: 84.90 % am Stand `c7045f81`); am Fixrunden-Stand (`4981766a` samt Plan-Nachzug) druckt `make gates` `coverage-gate: OK — Coverage 83.00% erfüllt Schwelle 80%`, `make test-store` nach frischer Messung des Replication-Teils `DB-Adapter-Coverage: 82.08% (gedeckt 843 von 1027 Statements; Profile gemergt: store,replication)` | in der Closure nachgezogen: die Zeile darunter |
| Nenner und gedeckte Zahl in den Sensor-Dokumenten (Suchlauf der Closure) | `git grep -c -w -e 1016 -e 2398 -e 830 -e 2033 -e 81.69 -e 84.90 <Stand> -- harness docs/user spec README.md AGENTS.md` | Stand `02b3059d` (vor dem Nachzug): `harness/sensors/coverage-gate.md` 4 Zeilen, `harness/sensors/db-adapter-coverage.md` 4 Zeilen, keine Zahl in `docs/user`, `spec`, `README.md`, `AGENTS.md`; Arbeitsbaum nach dem Nachzug: 0 Treffer. Neue Zahlen (`-e 1027 -e 2527 -e 843 -e 2096`): `coverage-gate.md` 2 Zeilen, `db-adapter-coverage.md` 3 Zeilen. Gemessen im Lauf der Closure: DB-Adapter `bash tools/harness/db-coverage.sh` über die abgelegten Profile druckt `DB-Adapter-Coverage: 82.08% (gedeckt 843 von 1027 Statements; Profile gemergt: store,replication)`, Anteile aus dem gemergten Profil abgeleitet: `postgresstorage` 686 (gedeckt 540), `postgresack` 32 (32), `postgressnapshot` 122 (118), `replication/receive` 187 (153); `make coverage-gate` (EXIT=0) druckt `total: (statements) 82.9%` und `coverage-gate: OK — Coverage 82.90% erfüllt Schwelle 80%`, das Profil des gebauten Images mit Awk über die Block-Position ausgezählt: gedeckt 2096 von 2527 (82,94 %), davon `postgressnapshot/snapshotlogic` 42 von 42. Übernommen: der Store-Lauf hinter dem Store-Profil (Verifikation §1, PostgreSQL 18); `83.00%` und `82.90%` desselben Stands `c092efa5` (Verifikation §1) | beide Sensor-Dokumente tragen die neuen Zahlen mit Lauf und Stand; die datierten Absätze zu `slice-097` und `slice-085` unverändert |
| Zählwort „beiden" vor Antragsarten (Verifikation V-1) | `git grep -n -e 'beiden Tabellen-Antragsarten' <Stand> -- internal tools harness spec` | Parent `2d47d8a7` 9 Zeilen, `895b1faf` 6, `02b3059d` 5. Der Treffer `queries.go` (Kommentar zu `SelectAppliedColumnRequests`: die Abfrage schließt drei Antragsarten aus, nicht zwei) ist in `02b3059d` behoben. Verbleibende fünf Treffer (`model/administrationrequest.go` Zeilen 16, 37, 56; `administrationrequest_test.go` Zeile 234; `wiring.go` Zeile 1345) bezeichnen `enable`/`disable` (Träger von Bindungs- und Publication-Menge), neben denen `backfill` namentlich genannt wird; Kontext gelesen | behoben in `02b3059d`; die fünf Treffer unverändert |
| E2E-Abdeckungs-Zeilennummern | `git diff --stat 2d47d8a7 -- test/integration tools/harness/run-integration-tests.sh` | leer (kein Treffer) | kein Nachzug; `make test-integration` lief unverändert grün (Regression der Verdrahtung, kein neuer Beleg) |


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
  **entfallen** — der Guard erkennt die Funktion: der Guard-Test-Lauf endet
  ungekürzt mit Exit 0 (Läufe 2 und 5 tragen den zweiten Rollout, Verifikation
  §1), der Unit-Test in `guard_test.go` ist grün (`make test`), und die Mutation
  „Eintrag `backfill_table(in:text,in:text,in:text)` aus `knownForeignObjects`
  entfernt" färbt `guard_test.go` rot (Review, Mutation K1).
- **Die neue View oder die erweiterte CHECK-Menge konvergiert nicht über einen
  Alt-Bestand.** Eine neue View konvergiert über einen Alt-Bestand (gemessen im
  Architect-Verdikt `architect-verdict-schema-rollout-view-signatur`,
  Szenario 5: Exit 0, zweiter Lauf Exit 0); CHECK-Menge und Funktion laufen über
  `nacharbeit-administration.sql` und sind über einen Alt-Bestand ungemessen.
  *Erwartet, zu belegen durch:* der Alt-Tag-Lauf von
  `run-schema-rollout-guard-test.sh`. **Ausgang:** **entfallen** — der
  Alt-Tag-Lauf gegen `v0.1.2` endet dreimal mit Exit 0 (Rollout des Tags,
  Arbeitsbaum mit Vorlauf, Arbeitsbaum zweiter Lauf) und prüft im Skript die
  fünf `request_kind`-Werte, `EXECUTE` allein für `cdc_admin` und 19 Rechte der
  drei Rollen; die Mutation „Funktion aus der `GRANT`-Zeile gestrichen" färbt
  Lauf 5 rot (Verifikation §1 und §4, M-F1).
- **`applied` wird als „Bestand kopiert" gelesen.** Bei dieser Antragsart heißt
  `applied` „angenommen"; die Ausführung steht in `cdc.backfill_run`. Ein Leser
  von `cdc.administration_request`, `diagnose` oder Handbuch könnte es anders
  deuten. *Erwartet, zu belegen durch:* [`SPEC-019`](../../../../spec/pflichtenheft.md) und Handbuch sagen es
  ausdrücklich; Review liest beide. **Ausgang:** **entfallen** —
  [`SPEC-019`](../../../../spec/pflichtenheft.md) („Für die Antragsart `backfill` heißt `applied`
  **angenommen**") und das Handbuch (§4: „`applied` heißt hier „angenommen"",
  Glossar „Angenommen (`applied` bei `backfill`)") sagen es; der Review las beide
  ohne Befund, der Verifier las beide am Stand `c092efa5`.
- **Die geschätzte Zeilenzahl wird zur Grenze** (`BEO-PGC/geschaetzter-wert-als-grenze`,
  offen, 1×): `pg_class.reltuples` ist eine Schätzung, für nie analysierte
  Tabellen als „unbekannt" erwartet (`-1`). Weder Handbuch noch Code dürfen sie
  als Grenze oder als `0` lesen; der Slice führt kein Ablehnen. *Erwartet, zu
  belegen durch:* Tests mit `-1` und die Wortwahl an jedem Träger. **Ausgang:**
  **entfallen** für diesen Slice — die Tests tragen `NULL` als „unbekannt", nie
  `0` (`TestFormatBackfillRunNamesTheEstimateAsEstimatedAndUnknownAsUnknown`,
  `TestBackfillStatusViewShowsTheLatestRunPerTable`; die Mutation
  `COALESCE(estimated_rows, 0)` färbt `TestDiagnoseReportsTheLatestBackfillRunPerTable`
  rot, Review D1), und der Review las „geschätzt" an jeder Stelle der Zeilenzahl;
  der Slice führt kein Ablehnen. Die Beobachtung `BEO-PGC/geschaetzter-wert-als-grenze`
  bleibt mit ihrem Zähler (1×) unberührt: dieser Slice liefert kein Vorkommen.
- **Ein `queued`-Run überlebt den Prozessstart nicht**, wenn der Worker nur einen
  In-Speicher-Kanal liest (`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger`,
  verkörpert): der Träger ist die Tabelle, das Signal weckt nur. *Erwartet, zu
  belegen durch:* der Neustart-Test und der Test „Signal während eines Runs" des
  zweiten Liefer-Punkts. **Ausgang:** **entfallen** — der Träger ist die Tabelle:
  `TestBackfillWorkerRunsQueuedRunsAtStartInOrderAndSkipsTheOthers` (Aufnahme beim
  Start, `interrupted` startet nicht), `TestReconcileBackfillRunsAgainstPostgreSQL`
  (Abgleich gegen die reale Run-Tabelle) und
  `TestBackfillWorkerKeepsASignalThatArrivesBetweenReadAndWait` (Signal zwischen
  Lesung und Warten); die Mutationen M1, M2, M4 und M5 färben sie rot (Review).
- **Die Administrations-Goroutine blockiert auf einer Kopie.** *Erwartet, zu
  belegen durch:* ein Test mit einem blockierenden Fake-Run, während ein
  weiterer Antrag verarbeitet wird. **Ausgang:** **entfallen** —
  `TestBackfillRunDoesNotBlockTheAdministrationGoroutine` belegt es mit einem
  blockierenden Fake-Run (Review und Verifikation §2 Zeile 2), `-race` grün.
- **Das Handbuch nennt eine Aussage ohne Beleg** (Startposition eines frisch
  registrierten Consumers; Richtgröße): beide bleiben bis `e2e` bzw.
  `bench-richtgroesse` **außerhalb** des Handbuchs ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B).
  *Erwartet, zu belegen durch:* Review des Abschnitts. **Ausgang:** **entfallen** —
  `grep -n -i 'Richtgr\|Startposition' docs/user/benutzerhandbuch.md` druckt keine
  Zeile (Verifikation §6); der Review nennt beide Aussagen ausdrücklich als nicht im
  Handbuch stehend.
- **Handbuch und Änderungshistorie** (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  und `BEO-PGC/handbuch-versionshistorie-uebersprungen`, verkörpert, je 3×): der
  Slice liefert mehrere Oberflächen (Funktion, View, `diagnose`); jede braucht
  ihre Stelle. *Erwartet, zu belegen durch:* der Suchlauf-Eintrag und die neue
  Historien-Zeile. **Ausgang:** **entfallen** — `Version: 1.49` mit den
  Historienzeilen 1.48 und 1.49 (Verifikation §6), der Suchlauf-Eintrag steht in §3;
  Funktion, View und `diagnose` haben je ihre Stelle im Handbuch (§4, §2, Glossar).
- **Eine Instanz je Quelle** trägt den Start-Abgleich; laufen zwei Instanzen gegen
  dieselbe Quelle, setzte die zweite den Run der ersten auf `interrupted`.
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) legt „eine Instanz je Quelle" fest ([`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)-Muster für
  `pending`); die Annahme ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1) und die Aufnahme (Festlegung 2)
  stützen sich auf dieselbe Festlegung; ein Schutz gegen zwei Instanzen
  derselben Quelle ist nicht Teil. Ein `backfill`-Antrag einer anderen Quelle
  nimmt die Administrations-Goroutine nicht an (er bleibt `pending`, Test unter
  Login-Rollen mit zwei Quellen); gibt es für seine Quelle keine Instanz, bleibt
  er `pending`. Die vier übrigen Antragsarten lesen weiter alle offenen Anträge
  ohne Quellbezug (Bestand aus [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)); das ist eine benannte Grenze, kein
  Teil dieses Slice. **Ausgang:** **weiter offen** — der Quellvergleich für die
  Antragsart `backfill` ist umgesetzt und getestet; die Rest-Grenzen (die vier übrigen
  Antragsarten ohne Quellbezug, zwei Instanzen derselben Quelle ohne Schutz) stehen im
  Register `BEO-PGC/ein-instanz-annahme-ohne-erzwingung` (Zähler 2×).

## 7. Closure-Notiz

- **Was hat funktioniert:** die Rollen-Kette lief in getrennten Kontexten und fand, was kein Gate
  las. Der Review (3 HIGH · 3 MEDIUM · 3 LOW · 5 INFO, Summary des Reports) fand durch Nachfahren
  und Mutieren, was die Läufe des Implementers grün ließ: den Guard-Lauf, der die zugesagten Rechte
  nicht führte (F-1), die Quellfilter-Mutation von `diagnose`, die grün blieb (F-2, Mutation D2),
  die Suchlauf-Zeile, deren Zahlen zu einem anderen Befehl gehörten (F-3), und das Handbuch-SQL,
  das unter der genannten Rolle mit „permission denied" endete (F-4, ausgeführt gegen einen
  ausgerollten Arbeitsbaum). Die Fixrunde (`b019cc2e`, `cabc6d14`, `947d9960`, `c0c6e286`,
  `120139ec`, `d5e5ea99`, `55085f76`, `4981766a`, `c092efa5`) behob F-1 bis F-9 und setzte F-10 um; die
  Verifikation bestätigte mit eigenen Läufen (Verifikation §1: `make gates`, `make test`,
  `make a-check`, `make coverage-gate`, `make test-store` zweimal und der Guard-Test je EXIT=0),
  mit vier Eingabeseiten-Mutationen der Fixrunde, alle rot (§4, M-F1, M-F2, M-F9, M-F10), mit dem
  Alt-Tag-Lauf von `v0.1.2` (Lauf 5: Exit 0/0/0) und mit dem Suchlauf-Feld an beiden Ständen (§5).
  Der Alt-Bestand-Beleg der Rollen-, Funktions- und CHECK-Zusagen steht im Repo (Lauf 5 des
  Guard-Test-Skripts), nicht als Wegwerf-Messung im Scratchpad; die Vorbedingungen am Alt-Stand
  (Funktion fehlt, CHECK ohne `backfill`, `cdc_admin` ohne `UPDATE`) verhindern, dass der Lauf
  trivial grün ist.
- **Was ging anders als geplant:** der Plan wuchs um `internal/bootstrap/backfill.go` als eigene
  Datei, um die Wiederholung des Workers nach einem Fehlerdurchgang (`backfillRetryInterval`, 5 s;
  ein Lesefehler ließe eine `queued`-Zeile sonst bis zum nächsten Signal liegen), um
  `administrationDeps.source` (Quellvergleich für `backfill`, F-10) und um den Vermerk-Test von
  `MarkApplied` (F-9); die Umfangsentscheidung („keine Rückführung nach `next/`") hielt am Diff:
  33 Dateien, 2106 Zeilen am Stand `c188be43` (Verifikation §3.3, **übernommen**). Die Fixrunde lief
  ohne eigenen Review-Report; der Verifier las, fuhr und mutierte das neue Material selbst (V-2).
  Der Nachzug der Zahlen in den Sensor-Dokumenten fand in der Closure statt und steht mit
  Lauf und Stand in `harness/sensors/db-adapter-coverage.md` (1027 Statements, gedeckt 843,
  gedruckt `82.08%`) und `harness/sensors/coverage-gate.md` (2527 Statements, gedeckt 2096,
  gedruckt `82.90%`); der Suchlauf der Closure steht als Zeile im Feld in §3. Ein einmaliger
  Ausfall von `TestWALRetentionThresholdEndToEnd` im Lauf des Implementers ist mit dem
  Repo-Skript nicht reproduzierbar (Review F-11: 30 von 30 und 12 von 12 grün am Slice-Stand und
  am Parent) und ohne Fix im Register geführt.
- **Verifier-Beobachtungen (V-1 bis V-5):** *V-1* (LOW) gezogen: der Kommentar zu
  `SelectAppliedColumnRequests` nennt die drei Tabellen-Antragsarten (Commit `02b3059d`, reine
  Kommentar-Korrektur); das Suchmuster der Plan-Zeile „Zählwörter" enthielt „vier", nicht „beiden"
  (Lerneintrag unten). *V-2* benannt, kein Nachzug: die Fixrunde ohne eigenen Review-Report; die
  DoD-Zeile „Review durchgeführt" ist durch den Report erfüllt. *V-3* gezogen: die Sensor-Dokumente
  tragen die Zahlen dieses Laufs. *V-4* benannt, kein Nachzug: der Zweig „Publication-Mitgliedschaft
  fehlt" ist im Use Case netzlos getestet, im Worker-Test gegen die reale Run-Tabelle nur der
  Bindungs-Zweig; beide liegen in derselben Vorbedingungs-Prüfung des Use Cases. *V-5* im Register
  (`BEO-PGC/test-schreibt-in-committete-datei`, `BEO-PGC/tier-skripte-hinterlassen-anonyme-volumes`).
- **Steering-Loop-Eintrag (Lerneintrag):** *Neuer Sensor-Umfang:* Lauf 5 des Guard-Test-Skripts
  (`tools/harness/run-schema-rollout-guard-test.sh`) führt die Alt-Bestand-Zusagen im Repo — 19
  Tabellen-/View-Rechte der drei Rollen, `EXECUTE` auf `cdc.backfill_table` allein für `cdc_admin`
  (auch nicht `PUBLIC`) und die fünf `request_kind`-Werte —, gebunden an die Eingabeseite: die
  Mutation „Funktion aus der `GRANT`-Zeile gestrichen" färbt Lauf 5 rot (Exit 1, Verifikation §4,
  M-F1), die `REVOKE`-Zeile ohne die Funktion ebenfalls (§3, Fixrunden-Mutationen). Die Klasse ist nicht
  neu: `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` erreicht mit diesem Slice 12×, Ausgang
  bleibt **verkörpert** (der Reviewer fährt den Beleg; er fand hier einen Beleg-Befehl kleiner
  als seine Beschreibung); `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (9×, F-2) trägt die
  Ausprägung „Filter-Eingabe ohne Zeile außerhalb des Filters". *Geschärfte Regel
  (Ausgangs-Kandidat, nicht entschieden):* ein Suchlauf nach einer bewegten **Menge** (vier → fünf
  Antragsarten, sechs → sieben Fremdobjekte) sucht neben dem Symbolnamen auch die **Zählwörter** und
  Aufzählungen (`vier`, `beiden`, `zwei`, `sechs`); der Symbolnamen-Suchlauf fand den Fall
  in `queries.go` nicht, der Zählwort-Suchlauf des Verifiers fand den Rest, den die Fixrunde
  übersah (V-1). Das ist die Hälfte der Grenze von [`AGENTS.md`](../../../../AGENTS.md) §3.13
  („`grep` trifft zuverlässig Symbolnamen, nicht Zahlen"), die sich mit einem Zählwort-Muster
  mechanisch schließen lässt — anders als ein Zeilen-Lokator, dessen Verschiebung keine
  wiederholbare Spur trägt. Herkunft: `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
  (5×, F-6, F-7, V-1); Träger-Kandidaten sind [`AGENTS.md`](../../../../AGENTS.md) §3.13 und der
  Suchlauf-Schritt von `implement-slice`, die Entscheidung liegt beim Lese-Schritt der Closure von
  [welle-backfill-bestand](../welle-backfill-bestand.md) (Architect). *Benannte Lücke (Architect):*
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 („je
  Quelle nimmt eine Instanz Anträge an") ist für die Antragsart `backfill` am Code gebunden — der
  Vergleich `request.Source != deps.source` in `processAdministrationRequests`, Unit-Test mit zwei
  Quellen, Login-Test mit einem Antrag einer fremden Quelle, Mutation rot (Verifikation M-F10), ohne
  Port-Änderung. Die Rest-Grenze ist offen: die vier übrigen Antragsarten lesen weiter alle
  offenen Anträge ohne Quellbezug (Bestand aus
  [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)), und zwei Instanzen
  derselben Quelle sind ohne Schutz. Beide ADRs sind `Accepted` und unberührbar; ob eine Ergänzung
  nötig ist (Quellfilter der Queue-Lesung, neue ADR mit `Supersedes`), entscheidet der Architect
  (`BEO-PGC/ein-instanz-annahme-ohne-erzwingung`, 2×).
- **Beobachtungs-Register (`../observations/`):** je Vorkommen eine `evidence/`-Datei
  `slice-backfill-sql-administration.md`, Zähler = Zahl der Dateien (real ausgezählt).
  *Bestehende Klassen:* `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (F-1, F-3) **12×**,
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (F-2) **9×**,
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (F-6, F-7, V-1) **5×** und
  `BEO-PGC/test-schreibt-in-committete-datei` (V-5) **4×** stehen **über** der Schwelle 3×; ihr Ausgang
  gehört dem Lese-Schritt der Closure von [welle-backfill-bestand](../welle-backfill-bestand.md)
  (die ersten beiden sind verkörpert mit Ausgangs-Kandidat, die letzten beiden ohne zugewiesenen
  Ausgang). `BEO-PGC/ein-instanz-annahme-ohne-erzwingung` (F-10) **2×**, offen.
  *Neue Klassen (je 1×, offen):* `BEO-PGC/handbuch-beispiel-nicht-unter-genannter-rolle-lauffaehig` (F-4),
  `BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` (F-5), `BEO-PGC/formatierungs-drift-ohne-gate`
  (F-8), `BEO-PGC/nicht-reproduzierbarer-test-ausfall` (F-11) und
  `BEO-PGC/tier-skripte-hinterlassen-anonyme-volumes` (V-5). Kein Eintrag erreicht mit diesem Slice
  erstmals 3×. *Benannt, nicht gezählt:* F-9 (Log-Meldung folgt der Bedeutung des Zustands
  nicht, ein Vorkommen, in der Fixrunde behoben), F-12 (Zahl streut zwischen Läufen; die Sensor-Dateien
  trugen ihre Zahlen mit Lauf, §3.12 erfüllt), F-13 und F-14 (Hinweise), V-3 und V-4.
- **Träger-Übergaben (§3.13):** die vom Slice bewegten Eigenschaften, die offene Pläne beschreiben,
  geprüft an deren Text: `slice-transformationen-antragsweg-schema` (§3-Zeile zum Guard-Test-Skript:
  Lauf 5 vergleicht die `request_kind`-Menge exakt mit fünf Werten und prüft `EXECUTE` für die eine
  Funktion; die Zeile führt den Alt-Tag-Lauf als „ausgeführt und geändert" und benennt Menge, Meldung
  und `EXECUTE` je neue Funktion als Nachzug), `slice-backfill-bench-richtgroesse` (§3-Zeile
  zu `diagnose`: der Ort ist `internal/bootstrap/backfill.go`, `diagnoseBackfillStatus`, nicht
  `wiring.go`). Ohne Änderung geprüft: `slice-transformationen-start-reihenfolge` (die Start-Reihenfolge
  Bindungsaufbau → Abgleich → Worker → Administrations-Goroutine → `stream.Run` stimmt mit der Zeile
  „Startpfad des Backfill-Workers" überein), `slice-backfill-e2e` (nennt den Startzeitpunkt des
  Runs, den zweiten Antrag als `failed` und `cdc.backfill_status`),
  `slice-transformationen-antragsweg-usecase`, `slice-backfill-sdk-origin`,
  `welle-backfill-bestand` und die Roadmap (keine Zustandsaussage berührt; die Register-Zähler in den
  §8-Sichtungslisten der offenen Pläne sind Stände der jeweiligen Planungszeit, der Zähler steht im
  Register).
- **Validator (Modul 8):** entfällt ausdrücklich — der Slice belegt die Auslösung netzlos (Fakes) und im
  Store-Tier; der Container-Lauf, den ein Nutzer auslösen und lesen kann, gehört zu `slice-backfill-e2e`
  (§1 dieses Plans: „Der E2E-Beleg am Container — `e2e`"). Der Bedarf aus
  [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) („ein Backfill wird ausgelöst") wird mit dem
  Wellen-Beleg validierbar. Kein stilles Überspringen.
- **Closure-Notiz-Review (`.harness/skills/closure-note-reviewer.md`):** eine getrennte Rolle im frischen
  Kontext, kein Schritt der Planner-Closure; der Skill prüft Slices in `done/` und greift daher erst nach
  dem `git mv` — hier nicht ausgeführt.
- **Folge-Slices:** keine neuen — die Folge-Slices der Welle
  [welle-backfill-bestand](../welle-backfill-bestand.md) liegen als Dateien in `open/`
  (`slice-backfill-e2e`, `slice-backfill-bench-richtgroesse`, `slice-backfill-sdk-origin`); der
  Start-Trigger von `slice-transformationen-antragsweg-schema` und
  `slice-transformationen-start-reihenfolge` („nach `slice-backfill-sql-administration`") ist mit
  dem Übergang nach `done/` erfüllt. Die Frage an den Architect zu
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) und
  [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md) ist kein Slice, sie ist
  im Register adressiert.
- **Risiken aus §6:** je ein Ausgang am Ort — Risiko 1 (Guard erkennt die Funktion nicht)
  **entfallen**, Guard-Test-Lauf und Mutation K1; Risiko 2 (View/CHECK über einen Alt-Bestand)
  **entfallen**, Alt-Tag-Lauf `v0.1.2` Exit 0/0/0 mit den Rechte- und CHECK-Prüfungen im Skript;
  Risiko 3 (`applied` als „Bestand kopiert" gelesen) **entfallen**, Spec und Handbuch sagen es;
  Risiko 4 (geschätzte Zeilenzahl wird zur Grenze) **entfallen** für diesen Slice, Tests und Review;
  Risiko 5 (`queued`-Run überlebt den Prozessstart nicht) **entfallen**, Neustart-, Abgleich- und
  Signal-Tests; Risiko 6 (Administrations-Goroutine blockiert) **entfallen**, Test mit blockierendem
  Fake-Run; Risiko 7 (Handbuch-Aussage ohne Beleg) **entfallen**, Suchlauf ohne Treffer und Review;
  Risiko 8 (Handbuch und Änderungshistorie) **entfallen**, Version 1.49 mit Historienzeilen 1.48 und
  1.49; Risiko 9 (eine Instanz je Quelle) **weiter offen** im Register
  (`BEO-PGC/ein-instanz-annahme-ohne-erzwingung`, 2×), der Quellvergleich für `backfill` ist umgesetzt.
  Ein Risiko „DB-Adapter-Coverage bewegt Zähler und Nenner" führt dieser Plan nicht; der Nachzug der
  beiden Sensor-Dokumente ist im Slice erledigt (siehe „Was ging anders als geplant" und §3).
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure: (a) *Anker* — der Lerneintrag trägt kein Feld
  `liegt in <Zielort>`, der Ausgangs-Kandidat der Zählwort-Regel ist dem Lese-Schritt adressiert; (b)
  *Folge-Slice* — `slice-backfill-e2e`, `slice-backfill-bench-richtgroesse` und
  `slice-backfill-sdk-origin` existieren als Dateien in `open/`; (c) *Register* — die genannten
  Verzeichnisse `BEO-PGC/<slug>` tragen je ein nicht leeres `evidence/`.

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
