# Slice transformationen-e2e-abhilfe: E2E Nichtanwendbarkeit und Abhilfe — der Prozess endet sichtbar, `cdc.remove_transformation` löst ihn

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Negative),
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (sichtbarer Fehlerzustand),
[`LH-FA-SCH-004`](../../../../spec/lastenheft.md) (Fehlerpfad der
Schemaänderung), [`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md) (kein
ACK, kein Datenverlust), [`LH-FA-REA-001`](../../../../spec/lastenheft.md)
(Lesen über `cdc.changes`),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 4 und Folgepflicht 5 (Akzeptanzkriterium der Abhilfe),
§Re-Evaluierungs-Trigger 4.

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Fehlerklasse `schema`), [`SPEC-019`](../../../../spec/pflichtenheft.md) —
gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-27.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Nichtanwendbarkeit einer Regel endet den Erfassungspfad sichtbar,
und die Abhilfe wirkt real: das **Abhilfe-Akzeptanzkriterium** von
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 5 ist am laufenden System belegt. Der Slice schließt erst, wenn es
belegt ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die deterministische Startreihenfolge** — `start-reihenfolge` liefert sie;
  dieser Slice belegt ihre Wirkung am System und ändert `Run` nicht.
- **Die Frist des Vorlaufs und der Beginn des Replikationsstroms** —
  `slice-start-vorlauf-grenze` liefert sie
  ([`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  Festlegung 1 bis 3: `START_REPLICATION` in `Stream.Run`, Vorlauf 30 s je
  Prozessstart, bei Ablauf startet der Stream und der Antrag bleibt `pending`);
  dieser Slice belegt die Abhilfe am Startpfad, den jener Slice ändert, und
  ändert weder Frist noch Aufrufstelle.
- **Ein allgemeiner Recovery-Weg für Schema-Fehler** —
  `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×): der Beleg deckt
  die Abhilfe für **diese** Ursache; die Frage, wie ein Betreiber nach einer
  inkompatiblen Typänderung oder einer entfernten Spalte weiterkommt, bleibt
  offen.
- **Die Abbildung der Nichtanwendbarkeit im Backfill-Pfad** — `backfill-pfad`.
- **Ein Beleg, der eine strikte Zeitordnung „vor der ersten Transaktion“
  misst** — der E2E-Lauf belegt die **Folge** der Ordnung (Antrag `applied`,
  kein zweiter `schema`-Fehler, Transaktion gelesen); die Ordnung selbst trägt
  der Test aus `start-reihenfolge`. Die Grenze steht im Bericht. **Was dieser
  Test trägt und was nicht** (Stand der Closure von
  `slice-transformationen-start-reihenfolge`): die Sequenz „Vorlauf, Goroutinen-Start,
  Stream“ steht in der Funktion `runStreamAfterAdministrationPass`
  (`internal/bootstrap/wiring.go`) und ist mit Fakes an ihre Eingabe gebunden
  (`TestRunStreamAfterAdministrationPassAppliesTheRuleRemovalBeforeTheStreamAssemblesTheFirstTransaction`,
  `…BindsAnEnabledTableBeforeTheStreamStarts`,
  `…KeepsTheStartOnAReadFailureAndReturnsTheStreamOutcome`,
  `…HoldsTheStreamUntilThePassEndsAndContextCancelEndsIt`, alle in
  `internal/bootstrap/administration_startorder_internal_test.go`); die
  **Aufrufstelle in `Run`** ist nur über zwei Quelltext-Tests gebunden
  (`TestRunSourceTextPassesStreamRunOnlyAsArgumentOfTheSequence`,
  `TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall`: Positionen
  der Aufrufe im Quelltext, kein Lauf). Vier Mutationen der Aufrufstelle bleiben im
  Unit-Lauf grün (Herkunft des Kontexts `streamCtx`, Rumpf von
  `startAdministration`, Inhalt der Wertegruppe, Aufruf in toter Verzweigung;
  Verifikation zu `slice-transformationen-start-reihenfolge` §4, V9b, V10b, V12,
  V14, **übernommen**). Bis zur Closure von `slice-transformationen-start-reihenfolge`
  legte kein Runner-Schritt einen beim Prozessstart `pending` stehenden Antrag
  (`git grep -n "'pending'" -- tools/harness/run-integration-tests.sh
  test/integration` traf keine Zeile, gemessen in jener Closure). Die erste
  Runner-Phase mit einem solchen Antrag (Antragsart `enable`, an einer Sperre
  wartend) liefert `slice-start-vorlauf-grenze`
  vor diesem Slice; der Lauf dieses Slice legt als erster einen
  `remove_transformation`-Antrag der Abhilfe über einen Prozessstart (der
  Schritt (b) unten).

## 2. Definition of Done

- [x] Nichtanwendbarkeit sichtbar (Kriterium (a)): nach einer aktiven
      `rename_column`-Regel `name` → `label` löst ein `ALTER TABLE … ADD COLUMN
      label` (kompatible Erweiterung, Zielname kollidiert) beim nächsten Change
      der Tabelle den Fehler aus: der Prozess endet, `diagnose` und der
      Heartbeat zeigen die Fehlerklasse `schema`, das Container-Log trägt den
      Text des neuen Sentinels; die Gegenprobe (`ADD COLUMN` mit anderem Namen)
      endet **nicht** mit diesem Fehler — der Beleg ist an die Kollision
      gebunden, nicht an irgendeine Spalten-Erweiterung; der
      spalten-entfernende Fall bleibt beim bestehenden Beleg (`relationOther`).
      *Zu belegen durch:* ein realer, grüner `make test-integration`-Lauf, Log-
      und `diagnose`-Ausgabe im Bericht.
- [x] Das Abhilfe-Akzeptanzkriterium aus
      [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
      Folgepflicht 5, **wörtlich**: „(a) eine Regel wird nichtanwendbar, der
      Prozess endet sichtbar mit Fehlerklasse `schema` (`diagnose`, Heartbeat);
      (b) `cdc.remove_transformation` wird beantragt, während der Prozess steht
      (Antrag bleibt `requested`); (c) nach dem Neustart ist der Antrag
      `applied`, **bevor** die erste Transaktion der Tabelle assembliert wird;
      (d) die zuvor nicht bestätigte Transaktion erscheint danach über
      `cdc.changes` — kein Datenverlust, keine zweite Neustart-Schleife.“ Der
      Zustand „`requested`“ des ADR-Wortlauts ist der Status `pending` von
      [`SPEC-019`](../../../../spec/pflichtenheft.md). (c) wird über seine
      Folge belegt (Antrag `applied` nach dem Neustart, Health ohne zweiten
      `schema`-Fehler, Transaktion in `cdc.changes`, Rohform, weil die Regel
      entfernt ist) und über die Ordnungs-Tests aus `start-reihenfolge` (Grenze:
      §1, vierter „NICHT“-Punkt — Fakes und Quelltext-Lesung der Aufrufstelle, kein
      Beleg am System außer diesem Lauf). *Zu
      belegen durch:* derselbe `make test-integration`-Lauf; die Phase läuft
      als eigener Aufruf nach der Container-Ende-Grenze des Runners (Muster von
      `TestE2ESchemaChangeDropColumn`: Neustart und Health-Poll davor und
      danach; am Start gelesen).
- [x] Die Abdeckung und die Klassen bleiben getragen: die Phase trägt die
      Kennung [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) im
      Deklarations-Anker und wird von einem `-run`-Muster erfasst; keine achte
      Fehlerklasse entsteht
      ([`ADR-0023`](../../adr/0023-fehlerklassifikation.md)); das
      Runner-Erzeugnis
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md) trägt die
      Zeile. *Zu belegen durch:* `make docs-check` und der `-run`-Abgleich im
      Bericht.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). Report:
      [`docs/reviews/review-slice-transformationen-e2e-abhilfe.md`](../../../reviews/review-slice-transformationen-e2e-abhilfe.md)
      (0 HIGH/MEDIUM/LOW, keine Fixrunde).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` §Sensors — die Zeile `make
      test-integration` nennt den Beleg; das Handbuch trägt die
      Abhilfe-Prozedur erst mit `slice-transformationen-betriebsdoku` (dessen
      §2 nennt sie als Gegenstand).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | update | `TestE2E…` für Nichtanwendbarkeit und Abhilfe ([`LH-FA-CFG-007`](../../../../spec/lastenheft.md) Negative). |
| `tools/harness/run-integration-tests.sh` | update | eigener Aufruf nach der Container-Ende-Grenze: Regel, `ADD COLUMN`, Log-/`diagnose`-Ausgabe, `cdc.remove_transformation` bei stehendem Prozess, `docker start`, Health-Poll, `cdc.changes`; `-run`-Muster; Deklarations-Anker. |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner. |
| `harness/README.md` §Sensors (`make test-integration`) | update | der Beleg in der Aufzählung. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „der Inhalt von
`make test-integration`“, „die Fälle, die den Erfassungspfad mit Klasse
`schema` beenden“; beide Stände gemessen).** Die Muster und Suchräume der Tabelle
stammen aus dem Schnitt der Welle und sind vor der Suchform von
[`AGENTS.md`](../../../../AGENTS.md) §3.13 geschnitten (ein Symbolname je Zeile,
Suchraum `harness docs`). Der Implementer trägt sie in Blöcke mit dem Etikett
`suchlauf` über und führt sie auf die Suchform: ganzer Baum, drei Arten des
Musters (Symbol, Zählwort, Beschreibung samt Hedge), jede Einschränkung mit Grund.
Die Aussage „nichts gefunden“ trägt nur, was Suchraum und Muster treffen
(`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`, Beleg
`slice-transformationen-start-reihenfolge`: die Spec führte eine Beschreibung der
Reihenfolge in anderen Worten als das Muster):

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Beschreibung der Belege von `make test-integration` | `grep -rn 'TestE2ESchemaChangeDropColumn' harness docs` | **Gefunden**: `harness/README.md` §Sensors nennt `TestE2ESchemaChangeDropColumn` als Anker eines bestehenden Belegs (Zeile 141, die sehr lange `make test-integration`-Zeile). | Dieselbe Zeile um den neuen Beleg ergänzt (`TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass`, Runner-Phase „Transformationen-Nichtanwendbarkeit und Abhilfe“). |
| Aufzählungen der Fälle mit Klasse `schema` | `grep -rn 'ErrIncompatibleSchemaChange\|ErrTransformationNotApplicable' docs harness spec internal` | **Gefunden**: beide Kennungen stehen bereits nebeneinander in `internal/bootstrap/wiring.go` (`classifyRunError`) und in `internal/adapters/driving/replication/mapper/mapper.go`; keine weitere Aufzählungsstelle in `docs`/`harness`/`spec` nennt beide Wege in einer Liste, die durch diesen Slice falsch würde. | Kein Nachzug nötig — dieser Slice fügt keinen dritten schema-Auslöser hinzu, er belegt nur den bereits existierenden `ErrTransformationNotApplicable`-Pfad am System. |
| Zeilenzahlen der Abdeckungstabelle | `git diff --stat` auf `tools/harness/run-integration-tests.sh` | 88 Zeilen berührt (84 Einfügungen, 4 Löschungen — die 4 Löschungen sind die korrigierte „letzter go-test-Aufruf“-Formulierung). | `docs/user/e2e-abdeckung.md` entsteht als Erzeugnis des Runner-Laufs neu (Anker-Zeilennummern verschieben sich automatisch). |
| `restart`-Verhalten des Feed-Containers | `grep -n 'restart' compose.yaml` | **Gefunden**: `restart: "no"` (Zeile 145) — unverändert seit dem Parent. | Die neue Phase beschreibt den Neustart korrekt als `docker start`-Betreiber-Handlung, kein automatischer Compose-Restart (siehe Kommentar in `tools/harness/run-integration-tests.sh`). |

Zusätzliche, auf die volle Suchform von [`AGENTS.md`](../../../../AGENTS.md)
§3.13 gezogene Blöcke (Parent `e6c5d087`, der Stand vor jeder Inhaltsänderung
dieses Slice; Suchraum: der ganze Baum ohne die drei Standard-Ausnahmen):

Symbolnamen der go-test-Aufrufe/Phasen, die den Erfassungspfad von
`make test-integration` mit Fehlerklasse `schema` dauerhaft beenden:

```suchlauf
e6c5d087 58 -n -E 'TestE2ESchemaChangeDropColumn|TestE2ESchemaChangeIncompatibleTypeChange|TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 69 -n -E 'TestE2ESchemaChangeDropColumn|TestE2ESchemaChangeIncompatibleTypeChange|TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Symbolnamen der beiden `schema`-Auslöser (Gefunden: beide bereits am Parent
nebeneinander in `wiring.go`/`mapper.go`; kein dritter Auslöser entsteht):

```suchlauf
e6c5d087 69 -n -E 'ErrIncompatibleSchemaChange|ErrTransformationNotApplicable' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 71 -n -E 'ErrIncompatibleSchemaChange|ErrTransformationNotApplicable' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Beschreibung „Container-Ende-Grenze“ (Zählwort/Beschreibung kombiniert — kein
separates Zählwort existiert für diese Eigenschaft, da keine Stelle im Baum
die Zahl der Container-Ende-Grenzen beziffert; geprüft und nicht gefunden):

```suchlauf
e6c5d087 17 -n -F 'Container-Ende-Grenze' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 18 -n -F 'Container-Ende-Grenze' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Hedge „(als eigener,) letzter (Rundlauf|go-test-Aufruf)“ — die Stelle, die
diese Phase als letzten Rundlauf des Runners auszeichnet: am Parent stand die
Aussage fälschlich an `TestE2ESchemaChangeIncompatibleTypeChange` (korrigiert,
§3 oben), am Diff steht sie korrekt an der neuen Phase und im neuen
`harness/README.md`-Satz — kein Netto-Zuwachs an falschen Stellen, nur eine
Verschiebung an die jetzt zutreffende Stelle:

```suchlauf
e6c5d087 1 -n -i -E 'letzter (Rundlauf|go-test-Aufruf)|läuft.{0,4}letzter' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 2 -n -i -E 'letzter (Rundlauf|go-test-Aufruf)|läuft.{0,4}letzter' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

## 4. Trigger

**Start** (`next` → `in-progress`): wenn
`slice-transformationen-start-reihenfolge`,
`slice-transformationen-e2e-wirkung`, `slice-capture-leerlauf-quellbelege` und
`slice-start-vorlauf-grenze` in `done/` liegen
(`slice-capture-leerlauf-quellbelege` trägt den Belegaufbau „Fehlerschwelle
erreicht, Container endet“ im selben Runner und geht der zweiten
Container-Ende-Grenze dieses Slice voraus; `slice-start-vorlauf-grenze` ändert
den Startpfad, den dieser Slice belegt: Frist des Vorlaufs und Beginn des
Stroms in `Stream.Run`, [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md);
Kante aus dem Verdikt
[`architect-verdict-welle-transformationen-offene-fragen`](../../../reviews/architect-verdict-welle-transformationen-offene-fragen.md)
§9.1) und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1). Der
Start-Trigger ist eine
**Vorab**-Bedingung: die Ordnung steht, bevor ihre Wirkung gemessen wird
(`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, offen, 2×).

**Verifier-Hinweis** (Verdikt
[`architect-verdict-leerlauf-bestaetigung-intermittenz`](../../../reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md)
§4): ein Rot mit der Signatur „Fehlerklasse `replication` … WAL-Rückstand … über
Fehlerschwelle“ in der Phase „Leerlauf-Bestätigung“ ist, solange
`slice-leerlauf-phase-last-in-stuecken` nicht in `done/` liegt, weder Beleg noch Widerlegung
dieses Slice — der Verifier wiederholt den Lauf (`gh run rerun <Lauf> --failed`) und nennt Lauf,
Versuchsnummer und Job-Kennungen beider Versuche; nach jenem Slice (die Kette in
[welle-transformationen](welle-transformationen.md) §5 legt ihn vor diesen Slice) ist dasselbe
Rot ein Befund und ein Architect-Zug. Die Phasen dieses Slice laufen hinter der Leerlauf-Phase;
ein Rot dort lässt sie ungelaufen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls
  Nichtanwendbarkeit (a) und Abhilfe (b)–(d) nicht in einem Review tragen — der
  abtrennbare Teil ist der Nichtanwendbarkeits-Beleg (erster Liefer-Punkt) als
  eigener Slice, die Abhilfe folgt mit Start nach diesem.
- `in-progress` → `open` (blockiert): falls Kriterium (c) oder (d) real nicht
  hält (zweite Neustart-Schleife, Transaktion fehlt): das ist der
  Re-Evaluierungs-Trigger 4 von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  — Architect-Frage nach einer Folge-ADR (Eingriff in die Startreihenfolge über
  den Slice `start-reihenfolge` hinaus); der Slice geht dann zurück, statt das
  Kriterium abzuschwächen.

## 5. Closure-Trigger

DoD vollständig (das Abhilfe-Akzeptanzkriterium ist belegt) + `make gates` grün
+ ein realer, grüner `make test-integration`-Lauf + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Beleg trifft die falsche Ursache.** Eine Spalten-Entfernung endet schon
  an `relationOther` mit derselben Klasse; ein Test, der sie auslöst, belegt
  die neue Prüfung nicht. *Erwartet, zu belegen durch:* der Auslöser
  `ADD COLUMN` mit dem Zielnamen der Regel und die Gegenprobe
  (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert, 20×
  Dateien unter `evidence/`, real gezählt bei dieser Closure).
  **Ausgang: entfallen.** Der Go-Test (`test/integration/integration_test.go:1235–1298`)
  fährt die Gegenprobe (`ALTER TABLE … ADD COLUMN other text`) aktiv geprüft
  (Row Image von id=2 muss `label`/`other` tragen), bevor er die reale Kollision
  auslöst — Review und Verifikation lasen das wörtlich nach, beide unabhängig.
  Der Beleg ist an die Namenskollision gebunden, nicht an eine beliebige
  Spalten-Erweiterung; kein neuer Beleg für dieses Register nötig.
- **Kriterium (c) hält nicht** trotz `start-reihenfolge` (zweite
  Neustart-Schleife). *Erwartet, zu belegen durch:* der reale Lauf; tritt es
  ein, gilt die Rückführung §4 und
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  §Re-Evaluierungs-Trigger 4. **Ausgang: entfallen.** Der reale
  `make test-integration`-Lauf (Implementer und Reviewer je eigenständig
  gefahren) zeigt genau eine Neustart-Schleife: `cdc.remove_transformation` wird
  `applied`, kein zweiter `schema`-Fehler tritt auf, die zuvor nicht bestätigte
  Zeile erscheint über `cdc.changes`. Die Reihenfolge (a)/(b) ist real als
  Skript-Reihenfolge verankert (Verifikations-Report §4: `docker inspect`
  steht vor dem `remove_transformation`-SQL-Aufruf, sequenzielles Bash-Skript,
  keine Möglichkeit einer anderen Ausführungsreihenfolge). Der
  Re-Evaluierungs-Trigger von `ADR-0112` feuert nicht.
- **Die Klassen-Unterscheidung ist im Heartbeat nicht sichtbar** — beide Wege
  enden mit `schema`; der Unterschied steht im Container-Log. *Erwartet, zu
  belegen durch:* die Log-Zeile im Beleg. **Ausgang: entfallen (adressiert wie
  geplant).** Der reale Lauf zitiert den neuen Sentinel-Text im Container-Log,
  getrennt vom bestehenden `relationOther`-Sentinel der Spalten-Entfernung
  (Review-Report §„Eigene Messungen", Runner-Ausgabe „Kollision auf
  feed_e2e_transform_abhilfe beendete den Erfassungspfad real"); `diagnose` und
  der Heartbeat tragen beide weiterhin nur die Klasse `schema`, wie geplant —
  die Unterscheidung liegt im Log, nicht in einer neuen Klasse.
- **Die Ordnung „vor der ersten Transaktion“ wird als gemessen ausgegeben**,
  obwohl der Lauf nur ihre Folge sieht (§3.12 Instanz B). *Erwartet, zu belegen
  durch:* die benannte Grenze in §1 und im Bericht. **Ausgang: entfallen
  (vermieden).** §1 und §2 dieses Plans benannten die Grenze bereits beim
  Schreiben korrekt als Folge-Messung, nicht als Zeitordnungs-Messung — kein
  nachträglicher Fund nötig. Die Verifikation
  (`docs/reviews/verifikation-slice-transformationen-e2e-abhilfe.md` §9 Punkt 2)
  bestätigt dies unabhängig: „Plan §1 und §2 kennzeichnen die Grenze bereits
  korrekt beim Schreiben (kein nachträglicher Fund nötig)".
- **Laufzeit und Timing** von Neustart und Health-Poll
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert, 5× Dateien
  unter `evidence/`, real gezählt bei dieser Closure): Poll auf
  Zustand mit Frist. **Ausgang: entfallen.** Beide realen `make
  test-integration`-Läufe dieses Slices (Implementer, Reviewer) liefen ohne
  Timing-Flake grün — die Phase pollt auf `docker inspect`-Zustand und
  `bf_await_applied` mit Frist, kein neuer Beleg für dieses Register.
- **Eine neue Testfunktion fällt still aus dem Runner**
  (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). *Erwartet, zu belegen
  durch:* der `-run`-Abgleich. **Ausgang: entfallen.** Die neue Runner-Phase
  ruft `go test -run
  '^TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass$'` explizit
  auf, der reale Lauf zitiert `--- PASS:
  TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass` — kein
  stiller Ausschluss, `BEO-PGC/test-runner-stiller-ausschluss` bleibt bei 2×.
- **Die Frist des Vorlaufs ändert den Startpfad, den dieser Slice belegt.**
  Entschieden ist die Zeitgrenze
  ([`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md),
  Umsetzung in `slice-start-vorlauf-grenze`,
  das vor diesem Slice in `done/` liegt): der Vorlauf trägt 30 s je
  Prozessstart (Wert **hergeleitet**, nicht gemessen), bei Ablauf startet der
  Stream, und der unterbrochene Antrag und jeder dahinter bleiben `pending`;
  eine Anzeige des Wartens in `diagnose` und `--healthcheck` gibt es nicht
  (Warn-Eintrag im Log, Festlegung 4). Der Beleg (b) und (c) fährt einen
  kurzen Antrag (`cdc.remove_transformation`); er belegt die Ordnung und ihre
  Folge für einen Antrag, den der Vorlauf innerhalb der Frist erreicht, nicht
  das Verhalten bei einem Antrag davor, der länger als die Frist läuft: dann
  startet der Stream mit dem bisherigen Regelstand, und ist die Regel nicht
  anwendbar, endet der Prozess wieder mit `schema` (Folge der ADR,
  §Konsequenzen; hergeleitet, kein Lauf). Der Bericht nennt diese Grenze und
  den Stand von `slice-start-vorlauf-grenze` (Register:
  `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad`, Ausgang *geplant*).
  **Ausgang: weiter offen.** Der Beleg dieses Slice deckt einen Antrag
  innerhalb der Vorlauf-Frist; das Verhalten bei einem Antrag, der die Frist
  überschreitet (Stream startet mit dem bisherigen Regelstand, ein zweiter
  `schema`-Fehler ist dann Folge der ADR, nicht widerlegt und nicht gemessen),
  bleibt ungeprüft am System — wie im Plan von Anfang an benannt. Register
  `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad` bleibt bei 2× (kein neuer
  Beleg, keine neue Datei), Ausgang unverändert *geplant*: die Frage ist mit
  [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  und `slice-start-vorlauf-grenze` bereits entschieden, die Grenze bleibt eine
  benannte, akzeptierte Restlücke ohne eigenen Trigger für einen weiteren
  Slice.

## 7. Closure-Notiz

- **Was hat funktioniert:** (1) Der neue Go-Test bindet die Nichtanwendbarkeit
  aktiv an die Zielnamen-Kollision (Gegenprobe `ADD COLUMN other` gegen die
  reale Kollision `ADD COLUMN label` gehalten), statt an irgendeine
  Spalten-Erweiterung — Review und Verifikation lasen den Testkörper unabhängig
  wörtlich nach und kamen zum selben Ergebnis. (2) Die neue Runner-Phase bindet
  Kriterium (a)/(b) des Abhilfe-Akzeptanzkriteriums real als Skript-Reihenfolge
  (`docker inspect` vor dem SQL-Aufruf), nicht nur als Behauptung — der
  Verifier hat dies durch reine, eindeutige Textlektüre bestätigt, ohne einen
  eigenen Fehlerpfad-Mutationslauf an der Bash-Phase zu brauchen (Verifikation
  §4, §9 Punkt 3: „ein eigener Fehlerpfad-Mutationslauf … wäre eine schwere
  Compose-Last ohne neue Erkenntnis gegenüber der bereits eindeutigen
  Textlektüre"). (3) Kriterium (c) ist von Anfang an korrekt als Folge-Messung
  gekennzeichnet (`AGENTS.md` §3.12 Instanz B) — kein nachträglicher Fund
  nötig, die Verifikation bestätigt dies eigenständig (§9 Punkt 2). (4) Drei
  Rollen (Implementer, Reviewer, Verifier) fuhren je einen eigenen, realen
  `make test-integration`-Lauf bzw. übernahmen ihn begründet (die Verifikation
  verzichtete auf einen dritten Lauf, weil `git diff --name-only
  58e684fc..HEAD` seit dem Review-Commit leer ist) — alle drei Rollen kamen zum
  selben Ergebnis. (5) Der Escape-Mechanismus (`trap cleanup EXIT`) räumt
  unbedingt ab, auch bei einem Abbruch mitten in der letzten Runner-Phase —
  kein Zustand überlebt einen Fehlschlag.
- **Was ging anders als geplant:** (1) **Der Reviewer behauptete fälschlich,
  `make doc-immutable` existiere im Repo nicht** (INFO, kein Slice-Befund) —
  sein Suchraum (`grep -rn immutable Makefile harness/mk/*.mk
  harness/README.md`) deckte `d-check.mk` am Repo-Wurzelverzeichnis nicht ab,
  wohin das Target tatsächlich gehört (eingebunden über `harness/mk/doc-gate.mk`s
  `include d-check.mk`); die Verifikation fuhr das Target real
  (`make doc-immutable RANGE=e6c5d087..HEAD`, Exit 0) und widerlegte die
  Behauptung. Dies ist das **17. Vorkommen** von
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (Form **Befehl**, LOW,
  bekannter Trägertyp, Deckel bereits bei 14× erreicht) — nach `state.md` dieses
  Eintrags braucht ein solches Vorkommen keine eigene `evidence/`-Datei,
  sondern steht mit Finding-Kennung hier: **V-1**
  (`docs/reviews/verifikation-slice-transformationen-e2e-abhilfe.md` §2).
  (2) **Eine kleinere Mutations-Zuschreibungs-Ungenauigkeit im Review** (V-3,
  LOW): der Reviewer zitierte einen Mutations-Kommentar zugunsten zweier
  Domänen-Unit-Tests, das Zitat gehört tatsächlich zu einer dritten
  Testfunktion (`TestConsumeRuleColumnMissingInRelationIsNotApplicable`); die
  Verifikation fuhr drei eigene Mutationen und bestätigte, dass das
  Gesamtverhalten (Fehlerklasse `schema`) trotzdem real und mehrfach abgesichert
  ist — kein Blocker für diesen Slice, kein eigener Registereintrag (zu
  geringfügig), aber verwandt mit derselben Klasse wie V-1. (3) Ein
  Zahlendreher des Reviewers bei einer Commit-Range-Angabe (V-2, INFO: „4
  Commits" statt der real gemessenen 5 bei `RANGE=c674d637~1..HEAD`) —
  inhaltlich folgenlos. (4) Alle sieben §6-Risiken lösen sich als *entfallen*
  auf, bis auf eines (die Frist des Vorlaufs jenseits der Belegreichweite),
  das *weiter offen* bleibt — siehe unten.
- **Steering-Loop-Eintrag (Lerneintrag):** *(a) Neuer Sensor:* keiner gebaut —
  ein Sensor, der prüft, ob ein genannter Beleg-Befehl tatsächlich existiert
  bzw. ob ein Suchraum vollständig ist, bräuchte eine Semantik-Entscheidung
  darüber, welche Dateien „zum Thema" gehören (dieselbe strukturelle Grenze wie
  bei den bisherigen 16 Vorkommen dieses Registers). *(b) Geschärfte Regel:*
  keine neue Regel — das 17. Vorkommen von
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` bestätigt die bereits
  bestehende HIGH-Klasse des Reviewer-Skills („Beleg trägt seinen Satz nicht":
  wer einen Beleg nennt, fährt ihn) an einem neuen Träger (`include`-Kette
  zwischen `harness/mk/doc-gate.mk` und `d-check.mk` am Repo-Wurzelverzeichnis)
  — der Skill selbst ändert sich nicht, die Instanz zeigt nur, dass ein
  `grep`-Suchraum über mehrere Verzeichnisebenen (hier: eine `include`-Kette
  über zwei Ebenen) leicht unvollständig bleibt, auch wenn der gesuchte
  Dateiname korrekt ist. *(c) Benannte Spec-Lücke:* keine neue — die bereits
  bekannte Lücke `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (ein
  allgemeiner Recovery-Weg für Schema-Fehler jenseits der
  Transformations-Nichtanwendbarkeit) bleibt bewusst offen, wie in §1 benannt.
- **Beobachtungs-Register (`../observations/`):** kein neues Verzeichnis, keine
  neue `evidence/`-Datei. *Kein Anfall, mit Begründung:* die sechs im Plan
  benannten Register (`negativtest-ohne-bindung-an-seine-eingabe`,
  `test-integration-retention-timing-flake`, `test-runner-stiller-ausschluss`,
  `dod-begruendung-unzutreffende-tatsachenbehauptung`,
  `wartegrenze-ohne-zeitgrenze-im-startpfad`) tragen an diesem Slice keinen
  neuen Beleg — jedes zugehörige §6-Risiko löst sich auf, ohne dass die
  Beobachtung real eintrat (siehe „Risiken aus §6" unten). *Deckel-Fall ohne
  Datei, Finding-Kennung hier* (vor dem Merge vom Verifier gefunden, Schwere
  LOW, bekannter Trägertyp): **V-1**
  (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`, 17. Vorkommen, Deckel bei
  14× erreicht, Ausgang unverändert **verkörpert**, Zähler unter `evidence/`
  bleibt bei 16, real gezählt bei dieser Closure). *Zur Kenntnis, kein
  Registerfall:* V-2 (INFO, Zahlendreher ohne Fehlschlussfolge), V-3 (LOW,
  Mutations-Zuschreibung — verwandt mit V-1s Klasse, aber zu geringfügig für
  einen eigenen Kandidaten).
- **Folge-Slices:** keine angelegt. Nächster startbarer Slice der Kette: die
  Welle-Tabelle (§4 dort) nennt als letzten Slice
  `slice-transformationen-betriebsdoku`, dessen Start-Trigger
  `slice-transformationen-e2e-abhilfe` **und** `slice-sdk-regel-realserver-e2e`
  in `done/` verlangt. `slice-sdk-regel-realserver-e2e` liegt unter `open/`,
  sein eigener Start-Trigger (`slice-transformationen-e2e-wirkung` in `done/`)
  ist bereits erfüllt und unabhängig von diesem Slice — er ist der nächste
  startbare Slice der Welle, sobald kein anderer Slice in `in-progress/` liegt
  (mit dem Move dieses Slice nach `done/` erfüllt, gelesen mit `ls` in dieser
  Closure). Übergaben mit Adresse: keine — dieser Slice ändert weder
  `wiring.go` noch den Startpfad über den bereits von `slice-start-vorlauf-grenze`
  gelieferten Stand hinaus.
- **Risiken aus §6:** je ein Ausgang, mit Begründung an der jeweiligen Zeile in
  §6. *Entfallen (sechs von sieben):* der Beleg bindet an die Namenskollision,
  nicht an eine falsche Ursache; Kriterium (c) hält, keine zweite
  Neustart-Schleife, kein Re-Evaluierungs-Trigger; die Klassen-Unterscheidung
  steht real im Log; die Ordnung wurde nie fälschlich als gemessen ausgegeben
  (von Anfang an als Grenze benannt); kein Timing-Flake in beiden realen
  Läufen; kein stiller Testausschluss (`-run`-Pattern real getroffen). *Weiter
  offen (eines von sieben):* die Frist des Vorlaufs — der Beleg deckt einen
  Antrag innerhalb der Frist, das Verhalten bei einem Antrag davor, der die
  Frist überschreitet, bleibt ungeprüft am System (Register
  `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad`, unverändert bei 2×,
  Ausgang *geplant* — die Frage selbst ist mit `ADR-0128` und
  `slice-start-vorlauf-grenze` bereits entschieden, kein neuer Trigger).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure. Vorab gelesen: *Anker* — dieser Slice
  trägt kein Feld „liegt in"; *Folge-Slice* — `slice-sdk-regel-realserver-e2e`
  und `slice-transformationen-betriebsdoku` liegen als Dateien in `open/`
  (gelesen mit `ls` in dieser Closure); *Register* — die fünf oben genannten
  Kennungen `BEO-PGC/…` (ohne `beleg-befehl-traegt-seinen-satz-nicht`, das
  keinen neuen Beleg dieses Slice trägt) sind Verzeichnisse mit nicht leerem
  `evidence/` (gemessen mit `ls evidence | wc -l` in dieser Closure: 20, 5, 2,
  11 und 2).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Test-Runner und Testpaket sind keine eigenen Sub-Areas
— kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad` (2×, Ausgang *geplant* →
`slice-start-vorlauf-grenze`, einschlägig — Risiko §6 letzter Punkt),
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (4×, Ausgang *geplant* →
`slice-start-vorlauf-grenze`, einschlägig — die Aufrufstelle des Vorlaufs ist
nur über Quelltext gebunden, §1), `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×, einschlägig —
Abgrenzung §1; dieser Slice liefert den Beleg für **eine** Ursache),
`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×, einschlägig —
Start-Trigger), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
(verkörpert, 6×), `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×),
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 8×, jeder genannte
Beleg-Befehl wird gefahren),
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 6×,
Risiko §6 vierter Punkt), `BEO-PGC/spec008-replication-luecke` (geschlossen,
gesichtet — die Klassen-Abbildung ist der Bestand),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
