# Architect-Verdikt: `make schema-rollout` hebt kein Alt-Schema über eine View-Signaturänderung (`slice-backfill-change-origin`)

**Rolle:** Architect (Modul 8)
**Anlass:** Der Implementer von `slice-backfill-change-origin` hat gemessen,
dass ein Ziel mit Alt-Schema `make schema-rollout` nicht übersteht, sobald die
View `cdc.changes` `origin` als letzte Spalte trägt (Slice-Plan §3 „Messung
zu Risiko §6", DoD-Punkt 2 „Grenze (1)"), und den Zug angehalten, statt
`knownForeignObjects` oder den Rollout selbst zu ändern. Die Frage: gilt der
Blocker für die ganze Klasse, und welcher Weg öffnet das Upgrade?
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-24
**Bezug:** [`LH-QA-OPS-005`](../../spec/lastenheft.md) (Upgrade-Sicherheit),
[`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md),
[`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
Folge-ADR [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)
(Proposed), `tools/schema/rolloutguard/guard.go`, `tools/schema/apply-rollout.sh`,
`tools/harness/run-schema-rollout-guard-test.sh`, `examples/bootstrap.sh`

---

## Verdikt

**Die Entscheidung gilt nicht als Nachbesserung, sondern wird ergänzt: neue
ADR, [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md),
Status Proposed, „ergänzt `ADR-0043`".** Der Plan von
`slice-backfill-change-origin` hat nichts Falsches behauptet — er hat die
Grenze korrekt gemessen und die Frage an den Architect gegeben. Kein
`Accepted`-Text wird geändert.

Entscheidung in fünf Zeilen:

1. Der Blocker trifft eine Klasse, nicht den Einzelfall: **Signaturänderung
   einer im neutralen Modell deklarierten View** (anhängen, umordnen,
   umbenennen, Typ). Additive Änderungen sind verträglich (gemessen).
2. Gewählt ist ein **automatischer, eng begrenzter Vorlauf im Target
   `schema-rollout`**: nach dem Precheck `DROP VIEW cdc.<name>` (ohne
   `CASCADE`, gemeldet) für genau die Views, deren Blocker `ReplaceView` mit
   `VIEW_SIGNATURE_INCOMPATIBLE` ist — nur wenn kein anderer unbekannter
   Blocker im Report steht (alles oder nichts).
3. Kein Vorlauf ohne Anlass; additive Änderungen bekommen keinen.
4. Den Beleg trägt ein neuer Lauf im bestehenden
   `run-schema-rollout-guard-test.sh` plus ein generischer Lauf „Schema des
   letzten Release-Tags → Arbeitsbaum".
5. Die Umsetzung ist die **Fixrunde von `slice-backfill-change-origin`** —
   kein neuer Slice.

**Träger:** [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md), weil die Entscheidung `ADR-0043` Entscheidung 3
(Pflicht-Report, destruktive Operationen default blockiert) um eine berichtete
Ausnahme ergänzt und den Blocker-Filter von `rolloutguard` lockert
(`AGENTS.md` §3.6: Lockerung eines Gates nur per ADR). Ein Verdikt allein
trüge das nicht. Der ADR-Index ist Sache des Aufrufers.

## Reproduktion (Architect-Lauf, 2026-09-24)

Aufbau: Wegwerf-PostgreSQL 18 (`PG_TEST_IMAGE`, `wal_level=logical`), je
Szenario eigener Container und eigenes Docker-Netz, alle abgeräumt; die
Repo-Stände als `git archive`-Kopien von `f4ba82ab` (Parent), `HEAD` und
`v0.1.2` in einem Temp-Verzeichnis (der Arbeitsbaum blieb sauber);
`make schema-rollout` im jeweiligen Stand mit `SCHEMA_TARGET` und
`SCHEMA_ROLLOUT_NETWORK`, d-migrate im Makefile-Pin. Alle Zahlen sind im Lauf
gedruckt (gemessen).

| # | Szenario | Ergebnis |
|---|---|---|
| 1 | Parent-Schema `f4ba82ab` ausrollen | Exit 0 |
| 2 | dann `HEAD`-Schema (View mit `origin`), Makefile-Pin | **Exit 2**; Precheck-Report `exitCode` 8, Blocker `MANUAL_ACTION_REQUIRED` für `ReplaceView:VIEW … changes`, Diagnose `VIEW_SIGNATURE_INCOMPATIBLE`; `rolloutguard`: „unbekannte Blocker-Klasse"; danach fehlt `cdc.change.origin` (Zählung 0) — der Fehlschlag ändert das Ziel nicht |
| 3 | dasselbe mit d-migrate 1.7.1 (`ghcr.io/pt9912/d-migrate@sha256:af9d3eb3…`) | Exit 2, dieselbe Diagnose (auch im zweiten Lauf) — Pin-Hebung hilft nicht |
| 4 | Ausgang Schema-Stand `v0.1.2` statt Parent | Exit 0; danach `HEAD` ohne Vorlauf: Exit 2, dieselbe Diagnose; mit `DROP VIEW cdc.changes` davor: Exit 0. Zwischen `v0.1.2` und `f4ba82ab` hat sich `tools/schema/` nicht geändert (`git log` leer) |
| 5 | Additiv über den Parent-Stand: neue Tabelle mit Fremdschlüssel und Default, neue nullable Spalte an `administration_request`, neue View, keine Änderung an bestehenden Views | Exit 0; Tabelle, View und Spalte vorhanden; zweiter Lauf Exit 0 |
| 6 | Vorlauf-Prototyp: `DROP VIEW cdc.changes`, dann `HEAD`-Rollout | Exit 0 (rund 8,8 s); zweiter Lauf Exit 0; Zeile ohne `origin` liest als `wal` (`c1|wal`); `has_table_privilege('cdc_reader','cdc.changes','SELECT')` = `t` |
| 7 | Leser-Schleife als `cdc_reader` (alle rund 50 ms `SELECT count(*) FROM cdc.changes`) während Szenario 6 | 218 Abfragen, 126 scheiterten (davon 12 mit `permission denied for view changes`, also View schon neu, Recht noch nicht); erster bis letzter Fehler rund 7 s |
| 8 | Option F gemessen: `ALTER TABLE … ADD COLUMN origin`, natives `CREATE OR REPLACE VIEW` mit angehängter Spalte, dann Rollout | Exit 0, Recht bleibt (`t`) — trägt, aber nur für Anhängen und mit zweiter SQL-Fassung der View |

Nicht gemessen (abgeleitet): dass der Feed-Container von dem Fenster
unberührt ist — `queries.SelectChanges` liest `cdc.change`/`cdc.transaction`/
`cdc.source_table`, nicht die View (`queries.go`, Kopfkommentar
„dieselbe Projektion wie die View"); dass ein Alt-Binary gegen das Schema mit
der zusätzlichen nullable Spalte läuft (erwartet, additiv).

## Klasse, nicht Einzelfall

| Geplante Schema-Änderung (Träger) | Klasse | Beleg |
|---|---|---|
| `origin` in `cdc.changes` (`slice-backfill-change-origin`) | **View-Signatur** — braucht den Vorlauf | Szenarien 2–4 |
| Tabelle `cdc.backfill_run` samt Grants (`slice-backfill-run-store`) | additiv | Szenario 5 (Tabelle mit FK und Default) |
| View `cdc.backfill_status`, Antragsart `backfill` (CHECK), Funktion `cdc.backfill_table` (`slice-backfill-sql-administration`) | View additiv (Szenario 5); CHECK und Funktion laufen über `nacharbeit-administration.sql` (bestehender Weg, ein Upgrade-Beleg im Slice steht aus) | Szenario 5 für die View |
| zwei nullable Spalten `rule_name`/`rule_spec` an `administration_request`, zwei Funktionen (`slice-transformationen-antragsweg-schema`) | additiv; `jsonb` statt `text` ist nicht gemessen — erwartet, zu belegen durch den Upgrade-Lauf | Szenario 5 (`text`) |
| übrige Transformations- und Backfill-Slices | kein Schema | Suchlauf über `docs/plan/planning/open/` (Zeilen zu `schema.yaml`/`view`/`ALTER`) |

Eine Auflage an die Planung ergibt sich daraus: `cdc.backfill_status` liefert
seine Warn-Spalten schon in `slice-backfill-sql-administration` (dort als
`false`), damit `slice-backfill-bench-richtgroesse` die View nicht erneut in
der Signatur ändert. Träfe es doch, fängt der Vorlauf es ab — es wäre nur ein
weiteres Lesefenster.

## Optionen

Die Vergleichstabelle (A Handbuch-Anleitung, B Vorlauf, C View aus dem Modell,
D zweite View, E d-migrate anheben, F natives `CREATE OR REPLACE VIEW`,
G nichts tun) steht in [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md). Der Ausschlag: **A** verschiebt Handarbeit
auf jeden Betreiber und jedes Release, ohne das Fenster zu vermeiden; **C**
macht `slice-016` rückgängig und vergrößert die Fremdobjekt-Menge; **D**
verfehlt `SPEC-002`; **E** ist gemessen wirkungslos (Szenario 3); **F**
vermeidet das Fenster, führt aber eine zweite SQL-Fassung der View und deckt
nur „anhängen" — als Rückfall benannt ([`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Trigger b). **B** wird gewählt,
weil es für jede künftige Signaturänderung ohne neue Datei trägt und sich
selbst auflöst, sobald d-migrate die Klasse rendert.

`LH-QA-OPS-005` bleibt im Wortlaut erfüllt: verlangt ist, dass persistierte
CDC-Daten nicht verloren gehen und der Datenstand vor/nach identisch lesbar
ist — die View trägt keine Daten (Szenario 6: die Zeile ist danach lesbar);
das Lesefenster ist eine Verfügbarkeitsfrage, kein Datenverlust. `ADR-0064`
bleibt unberührt: `make test-integration` tauscht den Container, nicht das
Schema. Dessen Trigger 1 („erster Git-Tag gesetzt") ist mit `v0.1.0`
(2026-09-19) eingetreten; [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7 schließt ihn für die
Schema-Hälfte über den Alt-Tag-Lauf (der Image-Vergleich Alt/Neu bringt über
den Container-Tausch hinaus für das Schema nichts — das Image trägt es nicht).

## Folge-Arbeit

**`slice-backfill-change-origin` (Fixrunde, nach Accepted von [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)):**

1. `tools/schema/rolloutguard`: `decide` liefert zusätzlich die Liste der
   `ReplaceView`-Views der Klasse (der Report trägt `diagnostics[].operationId`
   und `code`; `report.go` dekodiert sie bisher nicht); Unit-Tests: Klasse +
   bekannte Fremdobjekte → erlaubt; Klasse + unbekannter Blocker → nichts;
   `MANUAL_ACTION_REQUIRED` ohne die Diagnose → nichts; kein Blocker → kein
   Vorlauf.
2. `Makefile`-Target `schema-rollout`: Vorlauf-Schritt zwischen Precheck und
   `--execute` (`psql` über `PG_TEST_IMAGE` wie die `nacharbeit`-Schritte,
   `DROP VIEW` ohne `CASCADE`, je View eine Meldung); der Kommentarblock über
   dem Target beschreibt die Klasse im Indikativ (`AGENTS.md` §3.7).
3. `run-schema-rollout-guard-test.sh`: neuer Lauf vor dem Lauf mit dem
   unbekannten Blocker (View auf abweichende Signatur setzen, Rollout Exit 0,
   Soll-Signatur, Recht, Zeile lesbar); zusätzlich der Alt-Tag-Lauf
   (`git archive` des jüngsten `v*`-Tags, ausrollen, Arbeitsbaum ausrollen,
   zweimal Exit 0, Datenstand lesbar). Wegwerf-Umgebung, in jedem Ausgang
   abgeräumt.
4. `docs/user/benutzerhandbuch.md` §Schema aktualisieren: Reihenfolge
   (Rollout vor Container-Tausch), Vorlauf-Meldung, Lesefenster, „ein
   fehlgeschlagener Rollout am Precheck ändert das Ziel nicht"; Version und
   Änderungshistorie. `harness/README.md` §Sensors, Zeile `make
   schema-rollout` (Lauf-Zahl, Vorlauf).
5. Slice-Plan: §3 „Messung zu Risiko §6" trägt den Ausgang, DoD-Punkt 2
   „Grenze (1)" verweist auf [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) statt auf eine offene Frage, §6-Risiko
   „d-migrate konvergiert nicht" bekommt seinen Ausgang; die Closure-Trigger
   verlangen den neuen Lauf. Das committete `plan.yaml`/`down.sql` bleibt das
   Ergebnis eines Laufs gegen eine leere Datenbank — kein Upgrade-Lauf
   schreibt sie neu.
6. §3.13-Suchlauf: bewegte Eigenschaft „`make schema-rollout` gegen ein
   Alt-Schema" — Träger sind Handbuch, `harness/README.md`, der
   Makefile-Kommentar, der Kopfkommentar von `tools/schema/schema.yaml` und
   der Kommentar in `examples/bootstrap.sh` (der Existenz-Check dort ist
   akzeptiertes Negativ, siehe [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md); sein Kommentar „nicht idempotent"
   ist seit der Wache überholt und wird nur dann angefasst, wenn der
   Suchlauf ihn als Träger trifft).

**Andere Slices (Planner, kein Architect-Zug nötig):**
`slice-backfill-run-store`, `slice-backfill-sql-administration` und
`slice-transformationen-antragsweg-schema` nehmen den Alt-Tag-Lauf in ihren
Rollout-DoD-Punkt („zweiter Rollout") auf; `slice-backfill-sql-administration`
kann sein Risiko „View konvergiert nicht" mit Szenario 5 belegen (neue View
konvergiert). Kein Slice in `open/` ändert die Signatur einer bestehenden View.

**Notiz, keine Pflicht:** d-migrate rendert `CREATE OR REPLACE VIEW` mit
angehängten Spalten nicht, obwohl PostgreSQL es erlaubt (1.3.1 und 1.7.1). Ein
Hinweis an das Repo des Werkzeugs liegt beim Auftraggeber; [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Trigger a
löst den Vorlauf ab, sobald es gelöst ist.

## Auftraggeber-Fragen

Keine. Das Lastenheft bleibt unverändert; `LH-QA-OPS-005` ist im Wortlaut
erfüllt. Zu bestätigen ist nur der Status von [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) (Proposed → Accepted),
damit der Implementer sie als Constraint liest.

## Grenzen dieser Messung

Eine Instanz, ein Lauf je Szenario (kein Wiederholungslauf der Fenster-
Messung; die 7 s sind ein Richtwert für die Größenordnung); Datenmenge
minimal (eine Zeile); nicht gemessen sind ein Rollout gegen ein Ziel mit
laufendem Feed-Container und eine Änderung an `nacharbeit-*.sql` über einen
Alt-Bestand.
