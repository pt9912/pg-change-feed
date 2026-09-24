# ADR-0114: Schema-Rollout — automatischer Vorlauf für View-Signatur-Änderungen (ergänzt ADR-0043)

**Status:** Accepted

**Datum:** 2026-09-24

**Autor:** pt9912 (Architect-Rolle, Modul 8; anderer Kontext als der
Implementer-Lauf von `slice-backfill-change-origin`, der den Blocker real
gemessen und den Slice-Plan angehalten hat, statt ihn selbst aufzulösen)

**Bezug:** [`LH-QA-OPS-005`](../../../spec/lastenheft.md) (Upgrade-Sicherheit),
[`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Backfill; `origin` als
Auslöser), [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) (ergänzte
ADR — Entscheidung 3 und Re-Evaluierungs-Trigger),
[`ADR-0064`](0064-lh-qa-ops-005-testansatz-korrektur.md) (Upgrade-Testansatz;
Re-Evaluierungs-Trigger 1 ist mit dem ersten Release-Tag eingetreten),
[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) (Feld `origin`,
letzte Spalte der View `cdc.changes`), Architect-Verdikt
`architect-verdict-schema-rollout-view-signatur`
(Reproduktion, Messungen, Klassenanalyse)

**Schärft:** — (Werkzeug-/Prozess-ADR; die Schema-Form aus
[`SPEC-001`](../../../spec/pflichtenheft.md)/[`SPEC-002`](../../../spec/pflichtenheft.md)
bleibt unverändert, ergänzt wird der Weg, auf dem ein bestehendes Ziel sie
erreicht)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`make schema-rollout` ist der dokumentierte Upgrade-Weg
(`docs/user/benutzerhandbuch.md` §Schema aktualisieren: „mit demselben
Befehl wie bei der Ersteinrichtung erneut ausrollen"). Releases existieren
(`v0.1.0`…`v0.1.2`), also gibt es Ziele mit Alt-Schema.

Gemessen am 2026-09-24 (Wegwerf-PostgreSQL 18, `git archive`-Kopien des
Parent-Stands `f4ba82ab` und von `HEAD`, d-migrate im Makefile-Pin):

- Parent-Schema ausgerollt (Exit 0), danach das Schema mit `origin` als
  vierzehnter Spalte der View `cdc.changes`: `make schema-rollout` endet mit
  **Exit 2** (Precheck Exit 8). Blocker `MANUAL_ACTION_REQUIRED` für
  `ReplaceView changes`, Diagnose `VIEW_SIGNATURE_INCOMPATIBLE` („only
  renderable when view columns keep the same count, order, names and visible
  types"). Dasselbe Ergebnis mit dem Schema-Stand von `v0.1.2` als Ausgang und
  mit d-migrate 1.7.1 (`ghcr.io/pt9912/d-migrate`, lokal vorhanden). Der
  fehlgeschlagene Lauf ändert nichts am Ziel (die Spalte `change.origin` ist
  danach nicht angelegt).
- Additive Änderungen rollen über einen Alt-Bestand: neue Tabelle mit
  Fremdschlüssel, neue nullable Spalte an einer bestehenden Tabelle, neue View
  — Exit 0, ein zweiter Lauf Exit 0.
- Ein `DROP VIEW cdc.changes` vor dem Rollout macht den Lauf grün (Exit 0,
  zweiter Lauf Exit 0, Zeile ohne `origin` liest über die View als `wal`,
  `cdc_reader` hat das `SELECT`-Recht wieder). Die Sicht war für SQL-Leser
  rund 7 s lang nicht oder ohne Recht lesbar.
- PostgreSQL erlaubt `CREATE OR REPLACE VIEW` mit angehängten Spalten
  nativ; nach so einem Vorab-Ersatz (und angelegter Spalte) konvergiert der
  Rollout ohne Fenster (Exit 0, Recht bleibt erhalten) — aber nur für
  Anhängen, und nur mit einer zweiten, handgeführten SQL-Fassung der View.

Die einzige Klasse, die an d-migrate scheitert, ist die **Signaturänderung
einer im neutralen Modell deklarierten View** (Spalte anhängen, umordnen,
umbenennen, Typ ändern). Von den geplanten Schema-Änderungen der
Backfill- und Transformationswelle trifft sie allein `origin` in
`cdc.changes`; alle übrigen sind additiv.

## Entscheidung

Wir wählen einen **automatischen, eng begrenzten Vorlauf im Rollout-Target**
(Option B unten):

1. **Klasse.** Ein Blocker gehört zur Klasse „View-Signatur", wenn seine
   Operation `ReplaceView` (Objekttyp `VIEW`) ist und der Report zu dieser
   Operation die Diagnose `VIEW_SIGNATURE_INCOMPATIBLE` trägt. Ein
   `ReplaceView` existiert nur für eine im neutralen Modell deklarierte View
   und ein bestehendes Ziel-Objekt.
2. **Alles oder nichts.** Der Vorlauf läuft nur, wenn **jeder** Blocker des
   Precheck-Reports entweder zur Klasse „View-Signatur" gehört oder eine der
   bekannten Fremdobjekt-Operationen mit dem Grund
   `DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION` ist
   (`knownForeignObjects`, unverändert). Jeder andere Blocker lässt beides
   ausfallen: kein Vorlauf, kein `--allow-destructive`, Exit 8 wie bisher.
3. **Vorlauf.** Zwischen Precheck und `--execute` entfernt das Target je
   betroffener View `DROP VIEW cdc.<name>` — ohne `CASCADE`, ein Statement je
   View, jedes auf stdout gemeldet. Hängt ein fremdes Objekt an der View,
   scheitert der Lauf laut; abhängige Objekte des Betreibers werden nie
   mitgelöscht. d-migrate legt die View anschließend selbst an (Operation
   `CreateView` im Pflicht-Report); die Rechte setzt `nacharbeit-roles.sql`
   im selben Lauf.
4. **Kein Vorlauf ohne Anlass.** Ein Lauf gegen ein Ziel, dessen Views die
   Soll-Signatur tragen, führt keinen Vorlauf aus; das Fenster ohne View
   entsteht nur in einem Lauf, der eine Signaturänderung ausliefert.
5. **Additive Änderungen** (neue Tabelle, neue nullable Spalte, neue View)
   brauchen keinen Vorlauf und bekommen keinen.
6. **Reihenfolge des Upgrades.** Schema-Rollout vor dem Container-Tausch.
   Der Handbuch-Abschnitt „Schema aktualisieren" nennt sie und das Fenster.
7. **Beleg.** `tools/harness/run-schema-rollout-guard-test.sh` erhält einen
   Lauf, der die View `cdc.changes` auf eine abweichende Signatur bringt und
   verlangt: Rollout Exit 0, Soll-Signatur, `SELECT`-Recht für `cdc_reader`,
   bestehende Zeilen lesbar; der Lauf mit dem unbekannten Blocker (Exit 8)
   bleibt unverändert der letzte. Ein zweiter, generischer Beleg trägt den
   echten Alt-Bestand: ein Lauf, der den Schema-Stand des letzten
   Release-Tags (`git archive`) ausrollt und danach den Arbeitsbaum
   (Exit 0 zweimal, Datenstand lesbar) — er schließt für das Schema den
   Trigger 1 von [`ADR-0064`](0064-lh-qa-ops-005-testansatz-korrektur.md).
   Kein Gate: er braucht DB-Zugang und Docker-Netz.

`ADR-0043` Entscheidung 3 bleibt unberührt (Pflicht-Report, destruktive
Operationen default blockiert); diese ADR ergänzt eine benannte, berichtete
Ausnahme für eine Objektklasse, die keine Daten trägt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Handbuch-Anleitung „`DROP VIEW cdc.changes` vor dem Rollout" | kein Code; Betreiber wählt den Zeitpunkt | jede künftige Signaturänderung ist eine Betreiber-Handarbeit je Release mit Release-Notiz; vergessen = Exit 2 ohne Zustandsänderung (sicher, aber der dokumentierte Upgrade-Befehl trägt nicht mehr); nicht testbar außer per Handschritt im Test; dasselbe Fenster wie B |
| **B — automatischer, eng begrenzter Vorlauf (gewählt)** | `make schema-rollout` hält sein Versprechen für jede künftige Signaturänderung, auch Umordnen/Typ (nicht nur Anhängen); ohne Anlass kein Fenster; per Lauf testbar; löst sich selbst auf, sobald d-migrate die Klasse rendert (dann meldet der Precheck den Blocker nicht mehr) | ein Schreibzugriff des Wrappers außerhalb des d-migrate-Reports (mitigiert: gemeldet, ohne `CASCADE`, nur diese Klasse, `plan.yaml` zeigt `CreateView`); nicht atomar — scheitert `--execute` nach dem `DROP`, fehlt die View bis zum Wiederholungslauf (heilt idempotent); ein Fenster von einigen Sekunden für SQL-Leser |
| C — View aus dem neutralen Modell nehmen (`nacharbeit-*.sql`, `knownForeignObjects`) | d-migrate sieht die View nie | macht die deklarative Rückführung der Views in `schema.yaml` rückgängig; vergrößert die Fremdobjekt-Menge (Idempotenz-Kosten, `BEO-PGC/schema-rollout-fremdobjekte`); die Signatur-Änderung wird eine handgeschriebene SQL-Migration je View |
| D — `origin` nicht in `cdc.changes`, sondern in einer zweiten View | keine Signaturänderung | verfehlt [`SPEC-002`](../../../spec/pflichtenheft.md) (`origin` ist die letzte Spalte von `cdc.changes`) und die Lesefläche der SDKs; zwei Lesewege für dieselbe Zusage |
| E — d-migrate anheben | löst die Ursache | gemessen: 1.7.1 meldet denselben Blocker; Pin-Hebung wäre ohne Wirkung; ein Hinweis an d-migrate bleibt Notiz |
| F — natives `CREATE OR REPLACE VIEW` als handgeführter Vorlauf je Änderung | kein Fenster, Rechte bleiben | zweite SQL-Fassung der View außerhalb des Modells (Doppelquelle); nur Anhängen, nicht Umordnen/Typ; je Änderung eine neue Datei |
| G — nichts tun | kein Aufwand | der dokumentierte Upgrade-Befehl scheitert für jedes Ziel mit Alt-Schema ab dem Release, der `origin` trägt |

## Konsequenzen

- Positiv: `LH-QA-OPS-005` bleibt im Wortlaut erfüllt — die View trägt keine
  Daten, der Datenstand ist nach dem Rollout über `cdc.changes` identisch
  lesbar; die Klasse der Signaturänderung kostet die übrigen Slices der
  Backfill- und Transformationswelle nichts.
- Negativ: Für SQL-Leser über `cdc_reader` fehlt die View in einem Lauf, der
  eine Signaturänderung ausliefert, für die Dauer des Rollouts (gemessen rund
  7 s auf einer Testinstanz); der Feed-Container liest den Store über
  Tabellen, nicht über diese View. `make test-integration` fährt weiterhin
  keinen Schema-Wandel (`ADR-0064` Negativ bleibt: der Container-Tausch prüft
  das Schema nicht) — den Schema-Wandel trägt der Beleg unter Entscheidung 7.
- Negativ (akzeptiert): `examples/bootstrap.sh` überspringt den Rollout gegen
  ein bereits migriertes Ziel; ein weiterlebendes Demo-Volume mit Alt-Schema
  bekommt das neue Schema nicht — `make example-demo-down` (mit
  `down -v`) baut es neu auf, die Demo trägt keine schützenswerten Daten.
- Folgepflicht: Umsetzung als Fixrunde von `slice-backfill-change-origin`
  (siehe Verdikt): `tools/schema/rolloutguard`, `Makefile`-Target
  `schema-rollout` samt Kommentar, `run-schema-rollout-guard-test.sh`,
  `docs/user/benutzerhandbuch.md` §Schema aktualisieren,
  `harness/README.md` §Sensors (Zeile `make schema-rollout`).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `tools/harness/run-schema-rollout-guard-test.sh` (neuer Lauf) und Unit-Tests von `tools/schema/rolloutguard` | ein Ziel mit abweichender View-Signatur wird von `make schema-rollout` ohne Handschritt auf die Soll-Signatur gehoben; jeder andere Blocker bricht weiter mit Exit 8 ab | Skript, kein Gate (DB-Zugang, wie `make schema-rollout`); Unit-Tests in `make test` |

## Re-Evaluierungs-Trigger

Zwei beobachtbare Trigger: (a) d-migrate rendert `ReplaceView` mit geänderter
Signatur, der Precheck meldet die Klasse nicht mehr — der Vorlauf wird
Folge-ADR-pflichtig entfernt; (b) ein Betreiber meldet das Lesefenster als
nicht tragbar — Option F (natives `CREATE OR REPLACE VIEW` je Änderung) wird
neu bewertet.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-24 | Proposed | `architect-verdict-schema-rollout-view-signatur` |
| 2026-09-24 | Accepted — Annahme durch den Auftraggeber samt dem gemessenen Lesefenster von rund 7 Sekunden während des Vorlaufs | `architect-verdict-schema-rollout-view-signatur` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
