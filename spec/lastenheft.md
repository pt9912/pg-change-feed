# Lastenheft — PG Change Feed

**Projektname:** PG Change Feed
**Version:** 0.4.0 (`Major.Minor.Patch`); vor `Accepted` frei änderbar, ab
`Accepted` ist jede Änderung eine Vertragsänderung (siehe Historie).
**Status:** Draft
**Autor:** pt9912, **Datum:** 2026-09-12

---

## 1. Zweck und Geltungsbereich

PG Change Feed ist ein Open-Source-System zur Erfassung, Speicherung und
Bereitstellung von Änderungen aus PostgreSQL-Datenbanken (Change Data
Capture, CDC). Es richtet sich an Softwareentwickler, Datenbank- und
Systemadministratoren, Architekten, DevOps-/Plattform-Teams und Entwickler
von ETL-, Integrations- und Synchronisationslösungen, die Änderungen
konsumieren wollen, ohne selbst das Transaktionslog zu verarbeiten.

### Ausgangssituation

PostgreSQL stellt mit WAL, Logical Decoding, Logical Replication,
Publications und Replication Slots technische Mechanismen zur Erfassung und
Übertragung von Datenänderungen bereit. Diese Mechanismen bilden primär
technische Replikations- und Streaming-Primitiven. SQL Server CDC bietet
hingegen eine höherwertige Abstraktion, über die Änderungen aktiviert,
gespeichert und positions- bzw. bereichsbezogen abgefragt werden können.
Eine solche Abstraktion fehlt für PostgreSQL.

### Zielbild

PG Change Feed soll eine allgemeine CDC-Abstraktion für PostgreSQL
bereitstellen — PostgreSQL-native Fähigkeiten nutzend, ohne eine
vollständige technische Kopie von SQL Server CDC zu verlangen. Das
angestrebte Nutzungsmodell:

    Tabelle aktivieren
          |
          v
    Änderungen entstehen
          |
          v
    Änderungen werden persistent verfügbar
          |
          v
    Consumer liest ab bekannter Position
          |
          v
    Consumer bestätigt Verarbeitung
          |
          v
    Daten können gemäß Retention bereinigt werden

### Systemkontext

Quelle sind eine oder mehrere PostgreSQL-Datenbanken. Consumer können
unter anderem sein: ETL-/ELT-Prozesse, Data-Warehouse-Loader, Suchindizes,
Cache-Synchronisation, Replikations- und Integrationsdienste,
Microservices, kundenspezifische Anwendungen.

Zur Systemgrenze von PG Change Feed gehören:

- Erfassung relevanter Änderungen
- persistente Bereitstellung
- Konfiguration
- Consumer-Fortschritt
- Retention
- Betriebs- und Diagnoseinformationen

Die fachliche Verarbeitung der Änderungen durch Consumer liegt außerhalb
der Systemgrenze.

    +-----------------------------------------+
    |              PostgreSQL                 |
    |  Anwendungstabellen                     |
    |          |                              |
    |          v                              |
    |   Änderungsquelle                       |
    +----------+------------------------------+
               |
               v
    +-----------------------------------------+
    |               PG Change Feed            |
    |  Capture                                |
    |     |                                   |
    |     v                                   |
    |  Persistenz                             |
    |     |                                   |
    |     +--------------+                    |
    |     v              v                    |
    | Consumer State   Retention              |
    +----------+------------------------------+
               |
          +----+-----+
          v          v
      Consumer A  Consumer B

Die konkrete technische Änderungsquelle wird im Pflichtenheft bzw. in
ADRs festgelegt.

### MVP-Schnitt

Der MVP weist nach, dass eine robuste persistente CDC-Abstraktion mit
vertretbarem Betriebsaufwand realisierbar ist. Er umfasst die
Anforderungen mit der Kennzeichnung **MVP: ja**. Abnahmekriterien des MVP
— als Ende-zu-Ende-Prüfung über die referenzierten Anforderungen:

| MVP-Abnahmekriterium | Geführt als Akzeptanzkriterium von |
|---|---|
| INSERT, UPDATE und DELETE werden Ende-zu-Ende erfasst und gelesen | LH-FA-CAP-001 … 003, LH-FA-REA-002 |
| Die logische Reihenfolge ist korrekt | LH-FA-CAP-004, LH-FA-REA-004 |
| Rollback-Änderungen werden nicht ausgeliefert | LH-FA-CAP-007 |
| Ein Neustart verliert keine dauerhaft erfassten CDC-Daten | LH-FA-RET-001, LH-QA-REL-001 |
| Aufbewahrte Changes können erneut gelesen werden | LH-FA-REA-005 |
| Die Quellanwendung benötigt keine CDC-spezifischen SQL-Anpassungen | LH-FA-CFG-006 |
| Die vollständige Testumgebung ist automatisiert und reproduzierbar | LH-QA-POR-003 |

Der MVP-Integrationstest prüft mindestens den Ablauf: PostgreSQL starten →
CDC aktivieren → INSERT → UPDATE → DELETE → Changes lesen → Reihenfolge
und Inhalt prüfen.

### Erfolgskriterien

PG Change Feed gilt technisch als erfolgreich, wenn:

1. PostgreSQL-Änderungen zuverlässig erfasst werden.
2. Consumer die zugrunde liegende Änderungsquelle nicht kennen müssen.
3. Consumer Änderungen über stabile Positionen inkrementell verarbeiten
   können.
4. Neustarts ohne stille Datenverluste möglich sind.
5. die Auswirkung auf die Quelle akzeptabel bleibt.
6. Installation und Betrieb für PostgreSQL-Entwickler nachvollziehbar sind.
7. keine proprietäre Infrastruktur erforderlich ist.

## 2. Stakeholder

| Stakeholder | Rolle | Erwartung |
|---|---|---|
| Softwareentwickler | Consumer-Implementierer | Änderungen über stabile Positionen inkrementell lesen, ohne WAL-Kenntnisse |
| Datenbankadministratoren | Betreiber | SQL-basierte Administration, kontrollierte Retention, keine stillen Verluste |
| Software- und Systemarchitekten | Integrations-Entscheider | Allgemeine CDC-Abstraktion statt individueller Replikations-Primitiven |
| DevOps-/Plattform-Teams | Betreiber | Containerisierter Betrieb, Health Checks, maschinenlesbare Metriken |
| ETL-/Integrations-Entwickler | Consumer | Zuverlässiger, nachvollziehbar geordneter Änderungsstrom |
| Open-Source-Maintainer | Anbieter | Grundbetrieb ohne proprietäre Komponenten |

## 3. Funktionale Anforderungen

Regeln dieser Sektion: ID-Schema `LH-FA-<BEREICH>-<NNN>` — das Präfix `LH`
ist im ganzen Repo dasselbe und taucht in Make-Target-Kommentaren, ADRs und
Commits wieder auf (Baseline-Regelwerk `grundlagen-source-precedence.md`
§ID-Schema als Klammer). Jede Anforderung trägt drei Pfade — Happy ·
Boundary · Negative — plus Out-of-Scope (Baseline-Regelwerk
`modul-03-spec.md` §Ziel-Form: Akzeptanzkriterium). Ein Pfad trägt `—`,
wenn er für die Anforderung keine eigene Aussage trägt und das Verhalten
bereits durch einen anderen Pfad oder eine andere Anforderung bindend
definiert ist; `—` ist dann eine Verlagerung der Aussage, kein Verzicht.
Eine Anforderung, deren Bedarf ersatzlos entfällt, wird **nicht gelöscht**,
sondern trägt den Vermerk *zurückgezogen* im Titel; die Nummer bleibt
vergeben.

Traceability-Muster (Beispiel):

    LH-FA-CAP-002
          |
          +-- [`LH-FA-CAP-002.a`](pflichtenheft.md) / SPEC-<NNN>   (Verfeinerung/Festlegung, Pflichtenheft)
          +-- ARC-<NNN>                       (Architektur-Sicht)
          +-- ADR-<NNNN>                      (Entscheidung)

Kennungen werden nach Veröffentlichung nicht wiederverwendet. Die exakten
technischen Namen von Objekten (Schema, Funktionen, Views) werden von
diesem Lastenheft nicht vorgeschrieben; sie sind Gegenstand des
Pflichtenhefts.

### LH-FA-CFG-001 — CDC-Aktivierung je Tabelle

**Beschreibung:** CDC muss für einzelne Tabellen aktiviert werden können.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine existierende Tabelle `t`, when CDC für `t`
  aktiviert wird, dann werden fortan Änderungen an `t` erfasst.
- **Boundary:** Given CDC ist für `t` bereits aktiviert, when die
  Aktivierung erneut ausgeführt wird, dann ist das Ergebnis definiert
  (idempotentes Verhalten, keine doppelte Erfassung).
- **Negative:** Given die Tabelle `t` existiert nicht, when die Aktivierung
  ausgeführt wird, dann folgt ein expliziter Fehlerpfad, keine stille
  Erfolgsmeldung.

**Out-of-Scope:** Keine Aktivierung auf Nicht-Tabellen-Objekten (Views,
Foreign Tables) gefordert.

### LH-FA-CFG-002 — CDC-Deaktivierung je Tabelle

**Beschreibung:** CDC muss für einzelne Tabellen deaktiviert werden können.

**Akzeptanzkriterien:**

- **Happy Path:** Given CDC ist für `t` aktiviert, when CDC für `t`
  deaktiviert wird, dann werden fortan keine Änderungen an `t` mehr
  erfasst.
- **Boundary:** Given CDC ist für `t` nicht aktiviert, when die
  Deaktivierung ausgeführt wird, dann ist das Ergebnis definiert
  (idempotentes Verhalten).
- **Negative:** Given die Tabelle `t` existiert nicht, when die
  Deaktivierung ausgeführt wird, dann folgt ein expliziter Fehlerpfad.

**Out-of-Scope:** Verhalten bereits persistierter Changes einer
deaktivierten Tabelle folgt der Retention (LH-FA-RET-001 ff.), nicht der
Deaktivierung.

### LH-FA-CFG-003 — CDC-Status einer Tabelle

**Beschreibung:** Der CDC-Status einer Tabelle muss abfragbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given CDC ist für `t` aktiviert, when der Status
  abgefragt wird, dann wird der Zustand „aktiviert" gemeldet.
- **Boundary:** Given `t` wurde nie aktiviert, when der Status abgefragt
  wird, dann wird der Zustand „nicht aktiviert" gemeldet.
- **Negative:** Given `t` existiert nicht, when der Status abgefragt wird,
  dann folgt ein expliziter Fehlerpfad.

**Out-of-Scope:** Keine Aussage über die Abfrageform (SQL, CLI) — sie folgt
aus LH-FA-ADM-001 und LH-FA-SST-003.

### LH-FA-CFG-004 — Liste aktivierter Tabellen

**Beschreibung:** Alle aktivierten Tabellen müssen aufgelistet werden
können.

**Akzeptanzkriterien:**

- **Happy Path:** Given `t1` und `t2` sind aktiviert, when die Liste
  abgefragt wird, dann enthält sie mindestens `t1` und `t2`.
- **Boundary:** Given keine Tabelle ist aktiviert, when die Liste abgefragt
  wird, dann wird eine leere Liste zurückgegeben.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-CFG-005 — Spaltenauswahl (perspektivisch)

**Beschreibung:** Perspektivisch sollen einzelne Spalten von der Erfassung
ausgeschlossen bzw. gezielt ausgewählt werden können.

**Akzeptanzkriterien:**

- **Happy Path:** Given CDC ist für `t` aktiviert, when Spalte `c` vom
  Ausschluss konfiguriert wird, dann sind künftige Changes von `t` ohne die
  für CDC relevanten Datenwerte von `c`.
- **Boundary:** Given ein Spaltenausschluss ist konfiguriert, when die
  Spalte gelöscht wird, dann gilt das Verhalten aus LH-FA-SCH-003.
- **Negative:** Given die Spalte `c` existiert nicht, when der Ausschluss
  konfiguriert wird, dann folgt ein expliziter Fehlerpfad.

**Out-of-Scope:** Kein Bestandteil des MVP; eine nachträgliche
Ergänzung ohne Neuanforderung ist nicht vorgesehen.

### LH-FA-CFG-006 — Keine Anwendungscode-Anpassung

**Beschreibung:** Das Aktivieren von CDC darf keine Änderungen am
Anwendungscode der Quellanwendung erfordern. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given eine bestehende Quellanwendung schreibt in `t`,
  when CDC für `t` aktiviert wird, dann funktioniert die Anwendung unverändert
  weiter und deren Schreibungen werden erfasst.
- **Boundary:** Given die Anwendung nutzt Prepared Statements bzw. geplante
  Schreibpfade unverändert fort, when Änderungen entstehen, dann werden sie
  erfasst.
- **Negative:** —

**Out-of-Scope:** Performance-Auswirkungen auf die Quelle sind Gegenstand
von LH-QA-PER-001, nicht dieser Anforderung.

### LH-FA-CAP-001 — Erfassung von INSERT

**Beschreibung:** INSERT-Operationen müssen erfasst werden. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given CDC ist für `t` aktiviert, when eine Zeile in `t`
  eingefügt wird, dann entsteht ein erfasster Change vom Typ INSERT.
- **Boundary:** Given mehrere Zeilen werden in einem Statement eingefügt,
  when die Transaktion committet, dann ist jede eingefügte Zeile als
  eigener Change erfasst.
- **Negative:** Given die Quelltransaktion wird zurückgerollt, when Changes
  gelesen werden, dann sind die INSERT-Changes nicht als committed
  enthalten (LH-FA-CAP-007).

**Out-of-Scope:** —

### LH-FA-CAP-002 — Erfassung von UPDATE

**Beschreibung:** UPDATE-Operationen müssen erfasst werden. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given CDC ist für `t` aktiviert, when eine Zeile in `t`
  geändert wird, dann entsteht ein erfasster Change vom Typ UPDATE.
- **Boundary:** Given CDC ist für `t` aktiviert, when ein UPDATE einen Zeilenwert zuweist, ohne ihn zu ändern (Zuweisung des identischen Wertes), then ist das Erfassungsverhalten definiert
  (dokumentiert, ob ein Change entsteht).
- **Negative:** Given eine Quelltransaktion erzeugte einen UPDATE-Change, when sie zurückgerollt wird, then ist der Change nicht als committed enthalten (LH-FA-CAP-007).

**Out-of-Scope:** —

### LH-FA-CAP-003 — Erfassung von DELETE

**Beschreibung:** DELETE-Operationen müssen erfasst werden. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given CDC ist für `t` aktiviert, when eine Zeile in `t`
  gelöscht wird, dann entsteht ein erfasster Change vom Typ DELETE.
- **Boundary:** Given mehrere Zeilen werden in einem Statement gelöscht,
  when die Transaktion committet, dann ist jede gelöschte Zeile als
  eigener Change erfasst.
- **Negative:** Given eine Quelltransaktion erzeugte einen DELETE-Change, when sie zurückgerollt wird, then ist der Change nicht als committed enthalten (LH-FA-CAP-007).

**Out-of-Scope:** TRUNCATE ist ausdrücklich nicht gefordert (siehe
Out-of-Scope-Punkte, Abschnitt 5).

### LH-FA-CAP-004 — Eindeutige logische Reihenfolge

**Beschreibung:** Es müssen genügend Informationen vorliegen, um die
logische Reihenfolge der Änderungen eindeutig zu bestimmen. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given zwei Änderungen an derselben Zeile in zwei
  aufeinanderfolgenden Transaktionen, when beide Changes gelesen werden,
  dann lässt sich ihre Ausführungsreihenfolge eindeutig rekonstruieren.
- **Boundary:** Given mehrere Änderungen derselben Quelltransaktion, when
  gelesen wird, dann entspricht ihre Reihenfolge der Ausführungsreihenfolge
  innerhalb der Transaktion.
- **Negative:** —

**Out-of-Scope:** Eine Totalordnung über unabhängige, parallel committete
Transaktionen hinaus ist nicht gefordert; gefordert ist eine konsistente,
eindeutige Ordnung (siehe LH-FA-REA-004).

### LH-FA-CAP-005 — Transaktionszusammengehörigkeit

**Beschreibung:** Änderungen derselben Quelltransaktion müssen als
zusammengehörig identifizierbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given drei Änderungen in einer Quelltransaktion, when die
  Changes gelesen werden, dann tragen alle dieselbe
  Transaktionskennung.
- **Boundary:** Given zwei Transaktionen schreiben zeitgleich, when
  gelesen wird, dann sind ihre Änderungen unterscheidbar zugeordnet.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-CAP-006 — Auslieferung erst nach Commit

**Beschreibung:** Reguläre Consumer dürfen Änderungen erst dann als
dauerhaft erfolgreich ansehen, wenn die zugehörige Quelltransaktion
erfolgreich abgeschlossen wurde.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine offene Quelltransaktion erzeugte Changes,
  when ein regulärer Consumer liest, dann sind diese Changes nicht
  enthalten, solange die Transaktion nicht committed hat.
- **Boundary:** Given ein Consumer liest, when die Quelltransaktion genau in diesem Moment committed, then ist das Verhalten definiert (der Commit ist spätestens beim
  nächsten Lesevorgang sichtbar).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-CAP-007 — Keine Auslieferung zurückgerollter Änderungen

**Beschreibung:** Zurückgerollte Änderungen dürfen nicht als erfolgreich
committed Changes ausgeliefert werden. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Quelltransaktion wird committet, when gelesen
  wird, dann sind ihre Changes als committed enthalten.
- **Boundary:** Given eine Quelltransaktion wird committet und danach eine
  zweite startet und rollt zurück, when gelesen wird, dann sind nur die
  Changes der ersten als committed enthalten.
- **Negative:** Given eine Quelltransaktion wird zurückgerollt, when
  gelesen wird, dann ist keiner ihrer Changes als committed enthalten.

**Out-of-Scope:** —

### LH-FA-CAP-008 — Vorherige und neue Werte

**Beschreibung:** Soweit die Quelle dies zuverlässig ermöglicht, sollen
vorherige und neue Werte bereitgestellt werden.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein UPDATE committet, when der Change gelesen
  wird, dann sind vorheriger und neuer Zeilenstand (soweit für die
  konfigurierte Erfassung relevant, LH-FA-DAT-005) abrufbar.
- **Boundary:** Given ein INSERT bzw. DELETE wird committed, when der Change gelesen wird, then ist der vorherige bzw. neue Wert entsprechend nicht vorhanden (die Abwesenheit ist definiert,
  nicht ein Fehler).
- **Negative:** Given die Quelle kann einen Wert nicht zuverlässig liefern (z. B. nach Typänderung ohne Anpassung, LH-FA-SCH-004), when der Change gelesen wird, then ist dies erkennbar, nicht still gefälscht.

**Out-of-Scope:** —

### LH-FA-DAT-001 — Eindeutige Identifizierung jeder Änderung

**Beschreibung:** Jede Änderung muss eindeutig identifizierbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given zwei erfasste Changes, when ihre Identifikatoren verglichen werden, then sind sie unterscheidbar.
- **Boundary:** Given derselbe Datensatz, when er zweimal nacheinander identisch geändert wird, then entstehen zwei unterscheidbare Changes.
- **Negative:** —

**Out-of-Scope:** Die konkrete Form des Identifikators (z. B. LSN-basiert)
wird im Pflichtenheft festgelegt.

### LH-FA-DAT-002 — Identifizierbarkeit der Quelltabelle

**Beschreibung:** Die Quelltabelle muss identifizierbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein erfasster Change, when er gelesen wird, then ist die Quelltabelle (Schema und Tabelle) eindeutig erkennbar.
- **Boundary:** Given zwei gleichnamige Tabellen in verschiedenen Schemata sind aktiviert, when ihre Changes gelesen werden, then sind sie unterscheidbar zugeordnet.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-DAT-003 — Operationstyp

**Beschreibung:** Der Operationstyp muss mindestens INSERT, UPDATE und
DELETE unterscheiden.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein jeweils committed erfasster Change, when er gelesen wird, then ist der Typ INSERT, UPDATE oder DELETE ablesbar.
- **Boundary:** Given die Quelle führt einen weiteren Operationstyp aus (z. B. TRUNCATE), when die Änderung die Erfassung erreicht, then ist das Verhalten definiert (dokumentiert, ob und wie sie erfasst
  werden).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-DAT-004 — Sortierbare Position

**Beschreibung:** Jede Änderung bzw. Transaktion muss einer eindeutig
sortierbaren Quellposition oder daraus abgeleiteten Position zugeordnet
werden können. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given zwei committed Changes, when ihre Positionen verglichen werden, then lässt sich die Ordnung bestimmen; sie entspricht der logischen Reihenfolge (LH-FA-CAP-004).
- **Boundary:** Given mehrere Changes derselben Transaktion, when ihre Positionen verglichen werden, then erhält sich ihre Reihenfolge innerhalb der Transaktion.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-DAT-005 — Relevante Datenwerte

**Beschreibung:** Die für die konfigurierte CDC-Erfassung relevanten
Datenwerte müssen bereitgestellt werden.

**Akzeptanzkriterien:**

- **Happy Path:** Given CDC ist ohne Spaltenausschluss aktiviert, when ein Change gelesen wird, then enthält er alle Zeilenwerte, die die Quelle zuverlässig liefert.
- **Boundary:** Given ein Spaltenausschluss ist konfiguriert (perspektivisch, LH-FA-CFG-005), when ein Change gelesen wird, then enthält er die Werte der nicht ausgeschlossenen Spalten.
- **Negative:** Given ein Wert ist nicht lieferbar, when der Change gelesen wird, then ist seine Abwesenheit erkennbar, nicht mit einem Platzhalter überdeckt.

**Out-of-Scope:** —

### LH-FA-DAT-006 — Metadaten-Erweiterbarkeit

**Beschreibung:** Das Datenmodell muss um zusätzliche Metadaten erweiterbar
sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine künftige Version ergänzt Metadaten, when bestehende gespeicherte Changes gelesen werden, then behalten sie ihre Lesbarkeit.
- **Boundary:** Given ein Change ohne die neuen Metadaten, when er gelesen wird, then liest er sich unverändert weiter (Erweiterung ist abwärtskompatibel).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-REA-001 — Lesebereich zwischen zwei Positionen

**Beschreibung:** Änderungen müssen zwischen zwei Positionen abgefragt
werden können.

**Akzeptanzkriterien:**

- **Happy Path:** Given committed Changes mit Positionen `p1 < p2 < p3`,
  when der Bereich `[p1, p2)` abgefragt wird, dann wird genau der Change an
  `p1` zurückgegeben, nicht der an `p2`.
- **Boundary:** Given ein leerer Bereich (keine Changes zwischen den Grenzen), when er abgefragt wird, then wird eine leere Menge zurückgegeben.
- **Negative:** Given die Endposition liegt vor der Startposition, when der Bereich abgefragt wird, then folgt ein definierter leerer Bereich oder ein expliziter Fehlerpfad
  (dokumentiert, nicht unbestimmt).

**Out-of-Scope:** —

### LH-FA-REA-002 — Lesen ab bekannter Position

**Beschreibung:** Änderungen müssen ab einer bekannten Position gelesen
werden können. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Consumer kennt die zuletzt verarbeitete
  Position `p`, when er ab `p` liest, dann erhält er ausschließlich Changes
  nach `p`.
- **Boundary:** Given `p` liegt nach dem letzten vorhandenen Change, when ab `p` gelesen wird, then wird eine leere Menge zurückgegeben.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-REA-003 — Lese-Limit

**Beschreibung:** Die Anzahl zurückgegebener Änderungen muss begrenzt
werden können.

**Akzeptanzkriterien:**

- **Happy Path:** Given mehr verfügbare Changes als das Limit `n`, when
  gelesen wird, dann werden höchstens `n` Changes zurückgegeben.
- **Boundary:** Given weniger verfügbare Changes als `n`, when mit dem Limit `n` gelesen wird, then werden alle verfügbaren zurückgegeben.
- **Negative:** Given ein nicht plausibles Limit (z. B. negativ), when mit diesem Limit gelesen wird, then folgt ein expliziter Fehlerpfad oder eine dokumentierte Normierung.

**Out-of-Scope:** —

### LH-FA-REA-004 — Deterministische Sortierung

**Beschreibung:** Die Sortierung muss bei identischer Eingabe deterministisch
sein. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given dieselbe Abfrage bei unverändertem Datenstand, when sie zweimal ausgeführt wird, then sind beide Ergebnisreihenfolgen identisch.
- **Boundary:** Given mehrere Changes innerhalb derselben Positionsgrenze, when der Bereich abgefragt wird, then ist auch die interne Reihenfolge stabil definiert.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-REA-005 — Erneutes Lesen innerhalb der Aufbewahrung

**Beschreibung:** Bereits gelesene Änderungen müssen innerhalb der
Aufbewahrungszeit erneut gelesen werden können. **MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Change wurde bereits gelesen und ist innerhalb
  der Aufbewahrungszeit, when er erneut ab derselben Position angefragt
  wird, dann wird er erneut zurückgegeben.
- **Boundary:** Given die Aufbewahrungszeit ist abgelaufen und der Change wurde bereinigt, when er erneut angefragt wird, then wird er nicht mehr zurückgegeben (kein
  Fehlverhalten des Systems).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-REA-006 — Filter nach Quelltabelle

**Beschreibung:** Änderungen müssen nach Quelltabelle filterbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given Changes aus `t1` und `t2`, when nach `t1` gefiltert
  wird, dann werden nur Changes aus `t1` zurückgegeben.
- **Boundary:** Given der Filter trifft keine erfassten Changes, when danach gefiltert wird, then wird eine leere Menge zurückgegeben.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-CON-001 — Registrierung benannter Consumer

**Beschreibung:** Benannte Consumer müssen registriert werden können.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein neuer Name `c`, when `c` registriert wird, dann
  kann `c` fortan lesen und Positionen bestätigen.
- **Boundary:** Given `c` ist bereits registriert, when die Registrierung
  erneut ausgeführt wird, dann ist das Ergebnis definiert (idempotent oder
  dokumentierter Fehler).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-CON-002 — Unabhängige Consumer

**Beschreibung:** Mehrere Consumer müssen dieselben Änderungen unabhängig
voneinander verarbeiten können.

**Akzeptanzkriterien:**

- **Happy Path:** Given Consumer `c1` und `c2`, when beide denselben Änderungsbereich lesen und bestätigen, then beflussen ihre Vorgänge einander nicht.
- **Boundary:** Given `c1` bestätigt Position `p`, when `c2` anschließend liest, then bleibt die gelesene Position von `c2` unverändert.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-CON-003 — Persistierung der Verarbeitungsposition

**Beschreibung:** Die zuletzt bestätigte Verarbeitungsposition muss pro
Consumer gespeichert werden können.

**Akzeptanzkriterien:**

- **Happy Path:** Given Consumer `c` bestätigt Position `p`, when `c`
  danach seine gespeicherte Position abfragt, dann erhält er `p`.
- **Boundary:** Given `c` hat bereits eine Position bestätigt, when `c` erneut bestätigt, then wird die gespeicherte Position auf den neueren Wert fortgeschrieben.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-CON-004 — Bestätigung einer Position

**Beschreibung:** Ein Consumer muss eine erfolgreich verarbeitete Position
bestätigen können.

**Akzeptanzkriterien:**

- **Happy Path:** Given `c` hat Changes bis Position `p` erfolgreich
  verarbeitet, when `c` `p` bestätigt, dann gilt `p` als verarbeitete
  Position von `c` (LH-FA-CON-003).
- **Boundary:** Given `c` bestätigt Position `p`, when die Bestätigung für dieselbe Position wiederholt wird, then ist das Ergebnis definiert (idempotentes Verhalten).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-CON-005 — Fortsetzung nach Neustart

**Beschreibung:** Nach einem Neustart muss ein Consumer ab seiner zuletzt
bestätigten Position fortsetzen können.

**Akzeptanzkriterien:**

- **Happy Path:** Given `c` bestätigte `p` und startet neu, when `c` seine
  Position abfragt, dann erhält er `p` und setzt dort fort.
- **Boundary:** Given `c` hatte nie bestätigt, when `c` seine Position abfragt, then startet `c` an einer definierten Anfangsposition (dokumentiert).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-CON-006 — Administrative Entfernung von Consumern

**Beschreibung:** Consumer müssen administrativ entfernt werden können.

**Akzeptanzkriterien:**

- **Happy Path:** Given Consumer `c` existiert, when `c` administrativ
  entfernt wird, dann führt `c` keine Rolle mehr in der Retention
  (LH-FA-RET-004).
- **Boundary:** Given `c` hatte eine gespeicherte Position, when `c` entfernt wird, then ist deren Verhalten definiert (dokumentiert, ob sie entfernt wird).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-RET-001 — Persistente Speicherung

**Beschreibung:** Erfasste Changes müssen persistent gespeichert werden.
Ein Neustart darf bereits dauerhaft gespeicherte Änderungen nicht verlieren.
**MVP: ja.**

**Akzeptanzkriterien:**

- **Happy Path:** Given Changes wurden als dauerhaft gespeichert
  ausgewiesen, when das System kontrolliert neu startet, dann sind dieselben
  Changes weiterhin lesbar (LH-QA-REL-001, LH-QA-REL-002).
- **Boundary:** Given ein Change ist erfasst, aber die Quelltransaktion noch offen, when die Persistenz des Changes geprüft wird, then gilt die Aussage erst nach Commit (LH-FA-CAP-006).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-RET-002 — Konfigurierbare Aufbewahrung

**Beschreibung:** Die Aufbewahrung muss konfigurierbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine konfigurierte Aufbewahrungsregel, when die
  Bereinigung läuft, dann richtet sie sich nach dieser Regel.
- **Boundary:** Given eine geltende Aufbewahrungs-Konfiguration, when sie geändert wird, then gilt die neue Regel für nachfolgende Bereinigungsläufe.
- **Negative:** Given eine nicht gültige Konfiguration, when sie gesetzt wird, then folgt ein expliziter Fehlerpfad, keine stille Übernahme.

**Out-of-Scope:** —

### LH-FA-RET-003 — Zeitbasierte Retention

**Beschreibung:** Zeitbasierte Retention muss möglich sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Aufbewahrungszeit von `d`, when ein Change
  älter als `d` ist, dann kann er bereinigt werden.
- **Boundary:** Given ein Change ist exakt `d` alt, when die Bereinigung läuft, then ist definiert, ob er bereits bereinigt werden darf.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-RET-004 — Consumer-basierte Retention

**Beschreibung:** Consumer-basierte Retention soll ermöglichen, Änderungen
erst zu entfernen, wenn alle relevanten Consumer die erforderliche Position
bestätigt haben.

**Akzeptanzkriterien:**

- **Happy Path:** Given Consumer-basierte Retention ist aktiv und `c1`
  hat `p` noch nicht bestätigt, when die Bereinigung läuft, dann werden
  Changes ab `p` nicht entfernt.
- **Boundary:** Given alle relevanten Consumer haben `p` bestätigt, when die Bereinigung läuft, then dürfen Changes bis `p` gemäß Konfiguration entfernt werden.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-RET-005 — Erkennbarkeit blockierender Consumer

**Beschreibung:** Langsame Consumer, die Retention blockieren, müssen
erkennbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given `c1` blockiert die Bereinigung, when der Zustand beobachtet wird, then ist erkennbar, welcher Consumer blockiert und bis zu welcher Position.
- **Boundary:** Given mehrere blockierende Consumer, when der Zustand beobachtet wird, then sind alle einzeln erkennbar.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-RET-006 — Kontrolle des Datenwachstums

**Beschreibung:** Unbegrenztes Wachstum der CDC-Daten muss erkannt und
betrieblich kontrollierbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given wachsende CDC-Daten, when der Verbrauch beobachtet wird, then ist er über die Metriken ablesbar (LH-QA-OPS-003).
- **Boundary:** Given die Retention kann ein Wachstum nicht begrenzen (z. B. wegen blockierender Consumer), when dieser Zustand eintritt, then ist er erkennbar (LH-FA-RET-005).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-SCH-001 — Erkennung relevanter Schemaänderungen

**Beschreibung:** Für CDC relevante Schemaänderungen müssen erkannt
werden.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine aktivierte Tabelle, when sie eine neue Spalte erhält, then ist die Änderung für das System erkennbar.
- **Boundary:** Given eine nicht CDC-relevante Änderung (z. B. Index), when sie ausgeführt wird, then ist definiert, dass sie keine CDC-Reaktion auslöst.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-SCH-002 — Hinzugefügte Spalten

**Beschreibung:** Das Verhalten bei hinzugefügten Spalten muss definiert
und für Consumer nachvollziehbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Spalte wurde hinzugefügt, when die fortan erfassten Changes gelesen werden, then ist das Verhalten bezüglich der neuen Spalte definiert
  (z. B. Abwesenheit in älteren Changes ist erklärbar).
- **Boundary:** Given ältere Changes vor der Schemaänderung, when sie gelesen werden, then sind sie unverändert lesbar.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-SCH-003 — Entfernte Spalten

**Beschreibung:** Das Verhalten bei entfernten Spalten muss definiert sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine erfasste Spalte wurde entfernt, when künftige Changes erfasst werden, then ist ihr Verhalten definiert (dokumentiert).
- **Boundary:** Given ältere Changes, when sie gelesen werden, then ist definiert, ob und in welcher Form ihre historischen Werte lesbar bleiben.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-SCH-004 — Inkompatible Typänderungen

**Beschreibung:** Inkompatible Datentypänderungen müssen erkannt werden;
eine stille Fehlinterpretation ist nicht zulässig.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine kompatible Typänderung, when sie ausgeführt wird, then wird die Erfassung fortgesetzt.
- **Boundary:** Given eine Typänderung, deren Kompatibilität die Quelle zuverlässig nicht beurteilen kann, when sie eintritt, then ist das definierte Verhalten dokumentiert.
- **Negative:** Given eine inkompatible Typänderung, when sie eintritt, then ist sie erkennbar gemeldet; die Daten werden nicht still fehlinterpretiert.

**Out-of-Scope:** —

### LH-FA-SCH-005 — Schema-Version

**Beschreibung:** Änderungen sollen einer Schema-Version zugeordnet werden
können.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine erkannte Schemaänderung, when sie persistiert wird, then kann sie einer Version zugeordnet werden.
- **Boundary:** Given zwei Changes vor und nach einer Schemaänderung, when ihre Schema-Versionen gelesen werden, then lassen sie sich unterscheiden.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-ADM-001 — SQL-Administration

**Beschreibung:** Zentrale administrative Funktionen sollen über SQL
verfügbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein berechtigter Administrator, when er die administrativen Funktionen aufruft, then kann er Aktivierung, Deaktivierung, Statusabfrage und Consumer-Verwaltung über SQL ausführen.
- **Boundary:** Given die administrativen Objekte, when sie benannt werden, then ist die Namensgebung frei (die Anforderung schreibt Konzepte, keine Namen fest;
  Beispiele: `SELECT cdc.enable_table(...)`, `SELECT cdc.disable_table(...)`,
  `SELECT * FROM cdc.tables`, `SELECT * FROM cdc.consumers`).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-ADM-002 — Betriebsstatus

**Beschreibung:** Der Betriebsstatus muss abfragbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein betriebsbereites System, when der Status
  abgefragt wird, dann wird der aktuelle Betriebszustand gemeldet.
- **Boundary:** Given ein System in einem Fehlerzustand, when der Status abgefragt wird, then ist das auch über den Status erkennbar (LH-FA-ADM-003).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-ADM-003 — Sichtbare Fehlerzustände

**Beschreibung:** Fehlerzustände müssen sichtbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Fehlerzustand (z. B. Erfassung kann nicht fortsetzen), when er eintritt, then ist er erkennbar und unterscheidbar von normalem Betrieb (LH-QA-REL-003).
- **Boundary:** Given der Fehlerzustand endet, when der Zustand beobachtet wird, then ist das auch erkennbar.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-ADM-004 — Messbarer CDC-Abstand

**Beschreibung:** Der zeitliche bzw. positionsbezogene Abstand zwischen
Quelländerung und CDC-Verfügbarkeit muss messbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine committed Quelländerung, when ihr Abstand zur CDC-Verfügbarkeit gemessen wird, then lässt er sich bestimmen.
- **Boundary:** Given der Abstand wächst, when die Messung gelesen wird, then ist das beobachtbar (LH-QA-OPS-003).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-ADM-005 — Sichtbarer Verarbeitungsrückstand

**Beschreibung:** Ein Verarbeitungsrückstand muss sichtbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given unbestätigte Changes jenseits der bestätigten Position eines Consumers, when der Rückstand gemessen wird, then ist er erkennbar.
- **Boundary:** Given kein Rückstand, when die Messung gelesen wird, then zeigt sie das auch.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-SST-001 — Schnittstelle zur PostgreSQL-Quelle

**Beschreibung:** Es muss eine definierte Schnittstelle zur
PostgreSQL-Quelle bestehen.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine unterstützte PostgreSQL-Version, when die Erfassung betrieben wird, then erfolgt sie über eine dokumentierte, definierte Schnittstelle zur Quelle.
- **Boundary:** Given mehrere unterstützte PostgreSQL-Versionen (LH-QA-POR-001), when die Erfassung je Version betrieben wird, then ist die Schnittstelle je Version definiert.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-SST-002 — SQL-Zugriffe

**Beschreibung:** Konfiguration, Status und Kernlesezugriffe sollen über
SQL möglich sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein berechtigter Nutzer, when er Konfiguration, Status oder Kernlesezugriffe ausführt, then geschieht das über SQL.
- **Boundary:** Given der SQL-Zugriff, when die beteiligten Objekte gesucht werden, then sind sie dokumentiert (der Name bleibt frei, siehe LH-FA-ADM-001).
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-SST-003 — CLI

**Beschreibung:** Eine CLI soll Installation, Diagnose und Administration
unterstützen können.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Installation, when Installation, Diagnose oder Administration durchgeführt wird, then kann das über eine CLI erfolgen.
- **Boundary:** Given die CLI, when ihre Abdeckung geprüft wird, then deckt sie mindestens die Status-/Diagnoseabfragen ab, die LH-FA-ADM-002 … 005 nennen.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-SST-004 — Monitoring-Schnittstelle

**Beschreibung:** Eine Monitoring-kompatible Betriebsschnittstelle muss
möglich sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Monitoring-System, when es die Betriebsschnittstelle abfragt, then erhält es die Betriebsinformationen (Inhalt: LH-QA-OPS-003).
- **Boundary:** Given die Schnittstelle, when sie abgefragt wird, then ist ihr Format maschinenlesbar.
- **Negative:** —

**Out-of-Scope:** —

### LH-FA-SST-005 — Vorbereitung späterer HTTP-/gRPC-API

**Beschreibung:** Die Architektur soll eine spätere HTTP-/gRPC-API
ermöglichen, ohne das interne CDC-Modell grundlegend zu verändern.

**Akzeptanzkriterien:**

- **Happy Path:** Given das CDC-Modell, when eine andere Zugriffsschicht (z. B. HTTP-/gRPC-API) aufgesetzt wird, then lassen sich seine Lese- und Verwaltungsfähigkeiten anbieten, ohne das Modell umzubauen.
- **Boundary:** —
- **Negative:** —

**Out-of-Scope:** Protokoll-/Endpunkt-Details sind Sache der Spezifikation
und einer Architekturentscheidung (ADR), nicht dieses Lastenhefts — siehe
`LH-FA-SST-006`.

### LH-FA-SST-006 — Konkrete HTTP-/gRPC-API

**Beschreibung:** Eine konkrete HTTP-/gRPC-API muss die bestehenden
Lese- und Verwaltungsfähigkeiten (u. a. Consumer-Registrierung,
Positions-Bestätigung, Changes lesen, Status-/Diagnoseabfragen) zusätzlich
zu den bestehenden CLI-/SQL-Zugriffswegen über einen Netzwerkzugriffsweg
bereitstellen.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Netzwerk-Client, when er eine unterstützte
  Fähigkeit (z. B. Consumer-Registrierung) über die API aufruft, then wird
  sie ausgeführt, ohne dass CLI- oder SQL-Direktzugriff nötig sind.
- **Boundary:** Given dieselbe Fähigkeit über mehrere Zugriffswege (CLI,
  SQL, API), when sie über unterschiedliche Wege ausgeführt wird, then ist
  das Ergebnis fachlich gleichwertig (dieselbe Domänenlogik, kein
  Zweitpfad).
- **Negative:** Given ein nicht unterstützter oder nicht autorisierter
  Aufruf, when er über die API erfolgt, then wird er abgelehnt, nicht
  stillschweigend ignoriert.

**Out-of-Scope:** Wahl zwischen HTTP und gRPC, konkretes
Authentifizierungs-/Autorisierungsverfahren, Protokoll- und
Endpunkt-Details — das sind Architektur- (ADR) bzw. Spezifikations-Fragen
(`SPEC-*`), keine Lastenheft-Festlegung.

---

## 4. Nichtfunktionale Anforderungen

Regeln dieser Sektion: ID-Schema `LH-QA-<BEREICH>-<NN>`; Format: Anforderung
— messbare Formulierung — Messmethode.

Qualitätsziel-Prioritäten (Herkunft: Qualitätsziele des Projekts):

| Qualitätsziel | Priorität |
|---|---|
| Datenintegrität | sehr hoch |
| Betriebssicherheit | sehr hoch |
| Einfachheit (PostgreSQL-Entwickler ohne tiefgehende WAL-Kenntnisse) | hoch |
| Performance | hoch |
| Transparenz und Diagnosefähigkeit | hoch |
| Erweiterbarkeit | mittel bis hoch |

### LH-QA-REL-001 — Keine stillen Datenverluste

- **Anforderung:** Erkannte und bestätigte Quelländerungen dürfen nicht
  still verloren gehen. **MVP: ja.**
- **Messmethode:** Automatisierter Neustart- und Wiederholungstest der
  Erfassungs- und Lesepfade (Bestandteil des MVP-Integrations-tests,
  Abschnitt 1).

### LH-QA-REL-002 — Kontrollierter Neustart

- **Anforderung:** Das System muss kontrolliert neu gestartet und fortgesetzt
  werden können.
- **Messmethode:** Neustart-Test in der reproduzierbaren Testumgebung
  (LH-QA-POR-003).

### LH-QA-REL-003 — Sichtbarer Unzuverlässigkeitszustand

- **Anforderung:** Kann das System Änderungen nicht zuverlässig verarbeiten,
  muss dies deutlich erkennbar sein.
- **Messmethode:** Fehlerinjektionstest: Zustand herbeiführen und Erkennbarkeit
  über LH-FA-ADM-003 prüfen.

### LH-QA-REL-004 — Idempotente Consumer-Verarbeitung

- **Anforderung:** Modell und Schnittstellen sollen idempotente
  Consumer-Verarbeitung ermöglichen. Generisches Exactly-Once über externe
  Systeme hinweg wird nicht vorausgesetzt.
- **Messmethode:** Prüfung des Modells und der Schnittstellen (doppeltes
  Lesen desselben Bereichs liefert dieselben Changes, LH-FA-REA-005).

### LH-QA-PER-001 — Geringe Auswirkung auf die Quelle

- **Anforderung:** Die Auswirkung auf schreibende Quelltransaktionen soll
  minimiert werden.
- **Messmethode:** Benchmark: Schreibdurchsatz/-latenz der Quelle mit und
  ohne aktivierte CDC im Vergleich.

### LH-QA-PER-002 — Skalierbarkeit

- **Anforderung:** Das System soll von kleinen Datenbanken bis zu
  kontinuierlichen Änderungsvolumina skalieren können.
- **Messmethode:** Lasttests mit gestuften Änderungsvolumina; Volumina sind in
  `spec/pflichtenheft.md` §3 festgelegt ([`SPEC-014`](pflichtenheft.md)).

### LH-QA-PER-003 — Batch-Verarbeitung

- **Anforderung:** Effiziente Batch-Verarbeitung muss möglich sein.
- **Messmethode:** Benchmark: Lesevorgang über größere Change-Mengen im
  Batch gegen Einzelabruf.

### LH-QA-PER-004 — Latenz bis zur CDC-Verfügbarkeit

- **Anforderung:** Die Latenz zwischen Commit und CDC-Verfügbarkeit soll
  gering sein und über Benchmarks bewertet werden.
- **Messmethode:** Benchmark über LH-FA-ADM-004 (messbarer Abstand);
  Schwellen-Initialwerte sind in `spec/pflichtenheft.md` §3 festgelegt
  ([`SPEC-013`](pflichtenheft.md)).

### LH-QA-SEC-001 — Least-Privilege

- **Anforderung:** PostgreSQL-Berechtigungen sollen nach
  Least-Privilege-Prinzip vergeben werden können.
- **Messmethode:** Berechtigungsprüfung: dokumentierte Rollen reichen für
  die jeweiligen Aufgaben, nicht darüber hinaus.

### LH-QA-SEC-002 — Getrennte Berechtigbarkeit

- **Anforderung:** Administrative CDC-Funktionen und reguläre Lesezugriffe
  sollen getrennt berechtigbar sein.
- **Messmethode:** Berechtigungsprüfung: Lese-Rolle ohne Admin-Rolle kann
  lesen, nicht administrieren.

### LH-QA-SEC-003 — Beschränkbarkeit von CDC-Datenzugriffen

- **Anforderung:** CDC-Datenzugriffe müssen über PostgreSQL-Berechtigungen
  oder einen gleichwertigen Mechanismus beschränkbar sein.
- **Messmethode:** Berechtigungsprüfung: ohne die vergebene Berechtigung
  schlägt der Zugriff fehl.

### LH-QA-SEC-004 — Ausschluss sensibler Spalten

- **Anforderung:** Sensible Spalten sollen von der CDC-Erfassung
  ausgeschlossen werden können.
- **Messmethode:** Prüfung über LH-FA-CFG-005: ausgeschlossene Spaltenwerte
  erscheinen nicht in den Changes.

### LH-QA-OPS-001 — Containerisierter Betrieb

- **Anforderung:** Containerisierter Betrieb muss unterstützt werden.
- **Messmethode:** Deployment in der reproduzierbaren Docker-Umgebung
  (LH-QA-POR-003).

### LH-QA-OPS-002 — Automatisierte Health Checks

- **Anforderung:** Automatisierte Health Checks müssen möglich sein.
- **Messmethode:** Health-Check-Aufruf im Deployment, inklusive
  Fehlerzustandsfall (LH-FA-ADM-003).

### LH-QA-OPS-003 — Maschinenlesbare Metriken

- **Anforderung:** Maschinenlesbare Metriken sollen mindestens CDC-Lag,
  verarbeitete und ausstehende Changes, Speicherverbrauch, Alter des
  ältesten Changes, Consumer-Positionen und Fehler abbilden.
- **Messmethode:** Abzug der Metriken über LH-FA-SST-004 gegen einen
  kontrollierten Last-/Fehlerlauf.

### LH-QA-OPS-004 — Strukturiertes Logging

- **Anforderung:** Strukturiertes Logging soll unterstützt werden.
- **Messmethode:** Prüfung der Log-Ausgabe auf maschinenlesbare Struktur.

### LH-QA-OPS-005 — Upgrade-Sicherheit

- **Anforderung:** Upgrades dürfen persistierte CDC-Daten nicht verlieren.
- **Messmethode:** Upgrade-Test in der Testumgebung: Datenstand vor/nach
  Upgrade identisch lesbar (LH-QA-REL-001).

### LH-QA-POR-001 — PostgreSQL-Major-Versionen

- **Anforderung:** Mehrere aktiv unterstützte PostgreSQL-Major-Versionen
  sollen unterstützt werden.
- **Messmethode:** Testumgebung je unterstützter Version; konkrete Versionen
  sind in `spec/pflichtenheft.md` §3 festgelegt ([`SPEC-012`](pflichtenheft.md)).

### LH-QA-POR-002 — Primäre Zielplattform Linux

- **Anforderung:** Linux ist primäre Zielplattform.
- **Messmethode:** CI-/Testlauf unter Linux.

### LH-QA-POR-003 — Reproduzierbares Container-Deployment

- **Anforderung:** Ein offizielles bzw. reproduzierbares Docker-/OCI-Deployment
  soll bereitgestellt werden. **MVP: ja** (reproduzierbare lokale
  Docker-Umgebung).
- **Messmethode:** Aufbau der Umgebung aus dokumentierten Schritten/
  Artefakten auf leerer Maschine.

---

## 5. Globale Out-of-Scope-Punkte

Explizite Nicht-Anforderungen, die für das Gesamtsystem gelten:

- PG Change Feed ist **kein revisionssicheres Audit-System**.
- PG Change Feed ist **kein Message Broker** wie Kafka oder RabbitMQ.
- Technische Datenänderungen sind **nicht automatisch semantische Business
  Events**.
- PG Change Feed **ersetzt nicht** die native physische oder logische
  PostgreSQL-Replikation.
- Eine **vollständige Kompatibilität** mit SQL-Server-CDC-API und -Datenmodell
  wird nicht garantiert.
- TRUNCATE-Erfassung ist nicht gefordert.
- Kein generisches Exactly-Once über externe Systeme hinweg
  (LH-QA-REL-004).

Vorgesehene zukünftige Erweiterungen — die genannten Anforderungen bleiben
bindend; zurückgestellt ist jeweils ihre Produktionsreife bzw. Ausbaustufe
über die geforderte Fähigkeit hinaus:

- Mehrere unabhängige Consumer (LH-FA-CON-001 ff. — Fähigkeit gefordert;
  ihre Produktionsreife im Betrieb steht aus).
- Erweiterte Retention (LH-FA-RET-004 — Fähigkeit gefordert; Ausbau über
  die consumer- und zeitbasierte Basis hinaus).
- Schema Evolution (LH-FA-SCH-001 ff. — Erkennung und definiertes Verhalten
  sind gefordert; eine darüber hinausgehende Evolutions-Behandlung ist
  zukünftige Erweiterung).
- Produktionsreife Observability (LH-QA-OPS-001 ff. — Fähigkeiten
  gefordert; ihre Betriebshärtung in Produktionsumgebungen steht aus).
- High Availability (keine Anforderung dieses Lastenhefts).
- Exportadapter, beispielsweise Kafka, NATS, RabbitMQ, HTTP/Webhooks oder
  Object Storage (keine Anforderung dieses Lastenhefts; nicht Bestandteil
  des MVP).

## 6. Glossar

| Begriff | Bedeutung im Lastenheft |
|---|---|
| CDC (Change Data Capture) | Erfassung und Bereitstellung von Datenänderungen einer Quelle |
| Change (Änderungsdatensatz) | Ein erfasster, persistent bereitgestellter Einzelvorgang (INSERT/UPDATE/DELETE) |
| Position | Eindeutig sortierbare Kennung einer Änderung bzw. Transaktion, bezogen auf die Quelle oder daraus abgeleitet |
| Quelltransaktion | Die Transaktion in der PostgreSQL-Quelle, die die erfasste Änderung erzeugte |
| Quelltabelle | Die Tabelle in der PostgreSQL-Quelle, deren Änderung erfasst wurde |
| Consumer | Benannter Empfänger, der Changes ab Positionen liest und verarbeitete Positionen bestätigt |
| Retention | Konfigurierbare Regel, wie lange erfasste Changes aufbewahrt werden |
| Erfassung (Capture) | Der Vorgang, Quelländerungen in Changes zu überführen |
| Systemgrenze | Was PG Change Feed leistet; die fachliche Verarbeitung durch Consumer liegt außerhalb |
| MVP | Der erste lieferfähige Stand, der die Realisierbarkeit der CDC-Abstraktion nachweist |

## 7. Historie

Regeln dieser Sektion: Ab Status `Accepted` ist **jede** Änderung an diesem
Dokument eine Vertragsänderung — auch das **Hinzufügen** einer neuen
Anforderung. Sie entsteht **nur** aus einem Change Request, nie aus einem
ADR oder Slice. Vor `Accepted` frei änderbar ohne Change Request. Sind
Auftraggeber und Entwickler dieselbe Person, trägt der **Commit** die
Trennung: Die Änderung an dieser Datei liegt in einem eigenen Commit, **vor**
dem Slice, der sie umsetzt. Keine ADR-, Slice-, Carveout- oder Welle-Verweise
in dieser Tabelle (Decken-Regel).

| Version | Datum | Änderung | Verweis |
|---|---|---|---|
| 0.2.0 | 2026-09-09 | Initiale Fassung (Entwurf v0.2, vor Vorlagen-Überführung) | — |
| 0.3.0 | 2026-09-09 | Überführung in Lastenheft-Vorlagen-Struktur (Abschnitte 1–7); ID-Schema auf `LH-FA-<BEREICH>-<NNN>` / `LH-QA-<BEREICH>-<NNN>` normalisiert; Akzeptanzkriterien (Happy/Boundary/Negative) und Out-of-Scope je Anforderung ergänzt | — |
| 0.4.0 | 2026-09-12 | `LH-FA-SST-006` (konkrete HTTP-/gRPC-API) ergänzt; `LH-FA-SST-005`s Out-of-Scope-Klausel entsprechend angepasst — Auftraggeber und Entwickler sind dieselbe Person, Status ist `Draft` (frei änderbar ohne Change Request), diese Änderung liegt in einem eigenen Commit vor jedem umsetzenden Slice | — |