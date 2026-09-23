# Slice transformationen-e2e-abhilfe: E2E Nichtanwendbarkeit und Abhilfe — der Prozess endet sichtbar, `cdc.remove_transformation` löst ihn

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

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

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

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
- **Ein allgemeiner Recovery-Weg für Schema-Fehler** —
  `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×): der Beleg deckt
  die Abhilfe für **diese** Ursache; die Frage, wie ein Betreiber nach einer
  inkompatiblen Typänderung oder einer entfernten Spalte weiterkommt, bleibt
  offen.
- **Die Abbildung der Nichtanwendbarkeit im Backfill-Pfad** — `backfill-pfad`.
- **Ein Beleg, der eine strikte Zeitordnung „vor der ersten Transaktion“
  misst** — der E2E-Lauf belegt die **Folge** der Ordnung (Antrag `applied`,
  kein zweiter `schema`-Fehler, Transaktion gelesen); die Ordnung selbst trägt
  der Test aus `start-reihenfolge`. Die Grenze steht im Bericht.

## 2. Definition of Done

- [ ] Nichtanwendbarkeit sichtbar (Kriterium (a)): nach einer aktiven
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
- [ ] Das Abhilfe-Akzeptanzkriterium aus
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
      entfernt ist) und über den Ordnungs-Test aus `start-reihenfolge`. *Zu
      belegen durch:* derselbe `make test-integration`-Lauf; die Phase läuft
      als eigener Aufruf nach der Container-Ende-Grenze des Runners (Muster von
      `TestE2ESchemaChangeDropColumn`: Neustart und Health-Poll davor und
      danach; am Start gelesen).
- [ ] Die Abdeckung und die Klassen bleiben getragen: die Phase trägt die
      Kennung [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) im
      Deklarations-Anker und wird von einem `-run`-Muster erfasst; keine achte
      Fehlerklasse entsteht
      ([`ADR-0023`](../../adr/0023-fehlerklassifikation.md)); das
      Runner-Erzeugnis
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md) trägt die
      Zeile. *Zu belegen durch:* `make docs-check` und der `-run`-Abgleich im
      Bericht.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` §Sensors — die Zeile `make
      test-integration` nennt den Beleg; das Handbuch trägt die
      Abhilfe-Prozedur erst mit `slice-transformationen-betriebsdoku` (dessen
      §2 nennt sie als Gegenstand).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
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
`schema` beenden“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Beschreibung der Belege von `make test-integration` | `grep -rn 'TestE2ESchemaChangeDropColumn' harness docs` | *(Implementer trägt ein)* | Aufzählung in `harness/README.md` §Sensors trägt den Beleg |
| Aufzählungen der Fälle mit Klasse `schema` | `grep -rn 'ErrIncompatibleSchemaChange\|ErrTransformationNotApplicable' docs harness spec internal` | *(Implementer trägt ein)* | Aufzählungen nennen beide Wege; Handbuch-Stellen an `betriebsdoku` melden |
| Zeilenzahlen der Abdeckungstabelle | `git diff --stat` auf `tools/harness/run-integration-tests.sh` | *(Implementer trägt ein)* | Erzeugnis neu erzeugen |
| `restart`-Verhalten des Feed-Containers | `grep -n 'restart' compose.yaml` | *(Implementer trägt ein)* | der Neustart ist eine Betreiber-Handlung (`docker start`), kein automatischer Restart; der Beleg beschreibt ihn so |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn
`slice-transformationen-start-reihenfolge` und
`slice-transformationen-e2e-wirkung` in `done/` liegen und kein anderer Slice
in `in-progress/` liegt (WIP-Limit 1). Der Start-Trigger ist eine
**Vorab**-Bedingung: die Ordnung steht, bevor ihre Wirkung gemessen wird
(`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, offen, 2×).

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
  (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert, 6×).
  **Ausgang:** *(bei Closure)*
- **Kriterium (c) hält nicht** trotz `start-reihenfolge` (zweite
  Neustart-Schleife). *Erwartet, zu belegen durch:* der reale Lauf; tritt es
  ein, gilt die Rückführung §4 und
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  §Re-Evaluierungs-Trigger 4. **Ausgang:** *(bei Closure)*
- **Die Klassen-Unterscheidung ist im Heartbeat nicht sichtbar** — beide Wege
  enden mit `schema`; der Unterschied steht im Container-Log. *Erwartet, zu
  belegen durch:* die Log-Zeile im Beleg. **Ausgang:** *(bei Closure)*
- **Die Ordnung „vor der ersten Transaktion“ wird als gemessen ausgegeben**,
  obwohl der Lauf nur ihre Folge sieht (§3.12 Instanz B). *Erwartet, zu belegen
  durch:* die benannte Grenze in §1 und im Bericht. **Ausgang:** *(bei
  Closure)*
- **Laufzeit und Timing** von Neustart und Health-Poll
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert, 3×): Poll auf
  Zustand mit Frist. **Ausgang:** *(bei Closure)*
- **Eine neue Testfunktion fällt still aus dem Runner**
  (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). *Erwartet, zu belegen
  durch:* der `-run`-Abgleich. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Test-Runner und Testpaket sind keine eigenen Sub-Areas
— kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×, einschlägig —
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
