# ADR-0050: Schreibende SQL-Administration über eine Antrags-Queue mit LISTEN/NOTIFY, plus Live-Reload des Capture-Assemblers aus der DB

**Status:** Accepted — Supersedes: keine (ergänzt die offen gelassene Frage aus [`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md), ersetzt sie nicht)

**Datum:** 2026-09-13

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-FA-ADM-001`](../../../spec/lastenheft.md), [`LH-FA-CFG-001`](../../../spec/lastenheft.md), [`LH-FA-CFG-002`](../../../spec/lastenheft.md), [`LH-FA-CFG-003`](../../../spec/lastenheft.md), [`LH-FA-SST-002`](../../../spec/lastenheft.md), [ADR-0046](0046-sql-driving-adapter-lese-schreib-trennung.md), [ADR-0026](0026-composition-root.md), [ADR-0028](0028-inbound-use-cases.md), [ADR-0047](0047-rollenspezifische-dsn-verdrahtung.md)

**Schärft:** [`ARC-005`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`LH-FA-ADM-001`](../../../spec/lastenheft.md) (Lastenheft, Rang 1,
vertraglich abnahmebindend) verlangt, dass ein berechtigter Administrator
Aktivierung, Deaktivierung, Statusabfrage und Consumer-Verwaltung über SQL
ausführen kann (Boundary-Beispiele: `SELECT cdc.enable_table(...)`,
`SELECT cdc.disable_table(...)`). [`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)
(Accepted) hat entschieden, dass reine Lese-Views direkt gegen die
gespeicherten Tabellen lesen dürfen, schreibende/aktionsauslösende
SQL-Funktionen aber weiterhin ausschließlich über Inbound Ports
(`EnableTableUseCase`/`DisableTableUseCase`, `ADR-0028`) laufen — und hat
dabei ausdrücklich offen gelassen, **wie** eine reine SQL-Funktion einen
Go-Inbound-Port physisch erreichen soll: „Eine reine SQL-Funktion kann
einen Go-Inbound-Port ebenfalls nicht synchron aufrufen. Diese Frage
bleibt offen und ist nicht Gegenstand dieser Entscheidung."
Diese ADR löst genau diese Lücke; sie ändert nichts an `ADR-0046`s
Rollenspaltung (Lese-Views direkt, Schreiben über Ports) und ist deshalb
kein Supersedes, sondern eine Ergänzung.

**Zweiter, bislang undokumentierter Befund** (Fork-Recherche zur Eröffnung
der Feature-Welle „Verwaltungsfunktionen", registriert als
[`BEO-PGC/verwaltung-keine-sql-administration`](../planning/observations/BEO-PGC/verwaltung-keine-sql-administration/observation.md)):
`Assembler.tables` (`internal/adapters/driving/replication/mapper/mapper.go`)
trägt die Tabellen-Bindungen, gegen die der laufende Erfassungspfad
Änderungen filtert (`Assembler.Consume`/`Assembler.change`,
[`LH-FA-CFG-001`](../../../spec/lastenheft.md)). Diese Map wird in
`internal/bootstrap/wiring.go` **einmalig beim Prozessstart** aus der
Umgebungsvariable `CDC_TABLES` gebaut (`parseTables`) und an
`mapper.NewAssembler` übergeben; es existiert kein Mechanismus, der sie
zur Laufzeit erweitert. `Assembler` unterstützt bereits eine verwandte
Laufzeit-Mutation — `observeRelation` schreibt bei kompatibler
Schema-Erweiterung eine aktualisierte `TableBinding` in `a.tables[...]`
zurück (`ADR-0015`-Folgepflicht) —, aber ausschließlich für **bereits
aktivierte** Tabellen, nie für eine neu hinzukommende. Selbst wenn eine
schreibende SQL-Funktion die Bindungs-Zeile in `cdc.source_table` und die
Publication-Mitgliedschaft korrekt setzen könnte, würde der laufende
Erfassungspfad die neu aktivierte Tabelle **erst nach einem
Prozess-Neustart** sehen — und selbst dann nur, wenn `CDC_TABLES` (die
heutige alleinige Quelle für `Assembler.tables` beim Boot) die Tabelle
ebenfalls nennt, was ein rein SQL-seitig aktivierender Administrator
nicht steuert.

**Was der Code für den Schreibpfad bereits vorsieht:** `EnableTableService.Enable`
(`internal/application/usecase/enable/service.go`) prüft `TableExists`,
schreibt die Bindungs- und Schema-Version-Zeile idempotent
(`TableActivationAdapter.Register`) und trägt die Publication nach
(`TableActivationAdapter.Publish` — `CREATE PUBLICATION`/`ALTER
PUBLICATION ... ADD TABLE`, reine DDL). `spec/architecture.md`s
Sequenzdiagramm zu `LH-FA-CFG-001.a` benennt zusätzlich eine
`Replica Identity`-Prüfung als Schritt des Use Cases — eine echte
Domäneninvariante (dieselbe Kategorie, die `ADR-0046`s Kontext für
`EnableTableUseCase`/`AcknowledgeConsumerUseCase` als Grund nennt, warum
diese Pfade nicht in SQL dupliziert werden dürfen), auch wenn der
heutige Code sie noch nicht umsetzt (separater, hier nicht zu
schließender Befund). Jede Lösung, die eine SQL-Funktion die
Bindungs-Zeile/Publication **direkt** schreiben lässt, umgeht diese
Invariante strukturell oder zwingt zu ihrer Zweitimplementierung in SQL —
genau das Duplikations-Risiko, das `ADR-0018`/`ADR-0046` verhindern
wollen.

Das Repo hat für genau dieses Muster — eine Hintergrund-Aufgabe neben dem
Stream-Lauf, mit eigenem Verbindungspool derselben Rolle — bereits zwei
etablierte Beispiele in `internal/bootstrap/wiring.go`: `runHeartbeat`
und `runWALRetentionCheck`, beide als eigene Goroutine mit eigenem Pool
gegen `cfg.AdminDSN`/`cfg.CaptureDSN` (`ADR-0047`).

## Entscheidung

Wir wählen **D — Antrags-Queue mit LISTEN/NOTIFY-Weckung, verarbeitet vom
laufenden Capture-Prozess selbst**: Eine schreibende SQL-Funktion
(`cdc.enable_table(...)`, `cdc.disable_table(...)`) schreibt **ausschließlich**
einen Antrags-Datensatz in eine neue Tabelle im `cdc`-Schema (Name
Gegenstand des Pflichtenhefts, nicht dieser ADR — z. B.
`cdc.administration_request`: Quelle, Schema, Tabelle, Art
`enable`/`disable`, Zeitstempel, Status `pending`/`applied`/`failed`,
Fehlertext) und sendet `pg_notify` auf einem eigenen Administrations-Kanal.
Sie berührt `cdc.source_table`, `cdc.schema_version` und die Publication
**nicht direkt** — das bleibt exklusiv `TableActivationAdapter`
vorbehalten, aufgerufen über `EnableTableUseCase`/`DisableTableUseCase`.

Der bereits laufende Capture-Prozess bekommt eine neue Hintergrund-Goroutine
(Muster: `runHeartbeat`/`runWALRetentionCheck`, eigener Verbindungspool
gegen `cfg.AdminDSN`, Rolle `cdc_admin`), die den Kanal per `LISTEN` hält
und bei Wecksignal — sowie periodisch als Fallback für den Fall eines
verpassten `NOTIFY` (z. B. nach Verbindungsabbruch) — offene Anträge liest,
den passenden Inbound Port aufruft (einziger Schreibpfad bleibt der Port,
`ADR-0018`/`ADR-0046` unverändert), das Ergebnis im Antrags-Datensatz
vermerkt und **im selben Prozess** die laufende `Assembler`-Bindung
nachträgt bzw. entfernt (neue synchronisierte Methode am `Assembler`,
Muster bereits vorhanden in `observeRelation`, das `a.tables[...]` zur
Laufzeit schreibt — hier erstmals für eine neue statt eine bestehende
Bindung, und aus einer zweiten Goroutine, weshalb `a.tables` einen
Synchronisationsmechanismus braucht, den es heute nicht hat, siehe
Konsequenzen).

Zusätzlich baut `internal/bootstrap/wiring.go` `Assembler.tables` beim
Prozessstart künftig aus `cdc.source_table`/`cdc.schema_version`
(`TableActivationPort.List`, bereits vorhanden) statt ausschließlich aus
`CDC_TABLES` — sonst verliert ein Prozess-Neustart jede zwischenzeitlich
per SQL/CLI aktivierte Tabelle wieder, weil die Umgebungsvariable von der
laufenden Datenbank-Realität abweicht. `CDC_TABLES` bleibt als
Erstaktivierung/Bootstrap-Seed für eine leere Datenbank bestehen (dieselbe
idempotente `EnableTable`-Schleife wie heute), ist aber nicht mehr die
alleinige Quelle für den Laufzeit-Bindungsstand.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun (keine schreibenden SQL-Funktionen implementieren, `LH-FA-ADM-001`/`LH-FA-CFG-002`/`LH-FA-SST-003` bleiben beim heutigen Lese-only-Stand) | kein Aufwand, kein neues Risiko | verletzt dauerhaft eine vertraglich abnahmebindende Rang-1-Anforderung; genau der Zustand, dessentwegen die Feature-Welle vorgemerkt wurde — keine Option für ein Repo, das `LH-FA-ADM-001` ernst nimmt |
| B — Fork-Recherche-Vorschlag wörtlich: SQL-Funktion schreibt die Bindungs-Zeile und `ALTER PUBLICATION` direkt; der Go-Prozess pollt `cdc.source_table` periodisch und rekonziliert Differenzen gegen seine eigene `Assembler`-Bindung | physisch einfach, kein neuer Tabellentyp, löst das Reload-Problem beiläufig über den Poll | zwei Schreibpfade auf dieselben Zielzeilen (SQL-Funktion **und** `EnableTableService`/CLI/Tests) — genau die Domänenlogik-Duplikation, die `ADR-0018`/`ADR-0046` verhindern wollen, sobald die Replica-Identity-Prüfung aus `spec/architecture.md`s `LH-FA-CFG-001.a`-Diagramm real umgesetzt wird (die SQL-Funktion müsste sie nachbilden oder sie bliebe für den SQL-Pfad umgangen); Reconciliation muss zusätzlich Lösch-Fälle (Disable), Schema-Version-Zuordnung und Nebenläufigkeit zwischen SQL-Schreiber und Poll-Leser klären, ohne dass der Aufrufer ein Ergebnis-Feedback bekommt — mehr ungeklärter Zustand als die Queue-Variante, bei geringerem Nachvollziehbarkeitsgewinn |
| C — Dokumentierte Grenze „Aktivierung wirkt erst nach Neustart": SQL/CLI schreibt wie B direkt, Assembler liest beim Boot aus der DB statt aus `CDC_TABLES`, Betriebsdokumentation benennt die Neustart-Pflicht explizit | kein neuer Mechanismus (keine Queue, keine Goroutine), kleinster Implementierungs-Scope | verletzt den Zweck der Anforderung (administrative Aktivierung „im laufenden Betrieb", ohne dass ein Administrator den Prozess manuell neu startet) und löst das Duplikations-Problem aus B nicht — die SQL-Funktion schreibt weiterhin direkt an `TableActivationAdapter` vorbei |
| **D — Antrags-Queue mit LISTEN/NOTIFY-Weckung, verarbeitet vom laufenden Prozess (gewählt)** | ein einziger Schreibpfad für Bindungs-Zeile/Publication (voller Erhalt von `ADR-0018`/`ADR-0046`, auch sobald die Replica-Identity-Prüfung real umgesetzt wird); löst das Live-Reload-Problem strukturell, weil derselbe Prozess ausführt und seine eigene `Assembler`-Bindung nachträgt; niedrige Latenz durch `NOTIFY`, robust gegen verpasste Signale durch den Fallback-Poll; reines PostgreSQL-Bordmittel (kein FDW/`dblink`, passt zum PostgreSQL-Only-Stack); folgt einem im Repo bereits etablierten Muster (`runHeartbeat`, `runWALRetentionCheck`) | `SELECT cdc.enable_table(...)` wird asynchron — die Funktion bestätigt nur „beantragt", nicht „aktiv"; ein Aufrufer beobachtet den Abschluss über die Statusspalte bzw. eine Lese-View (`ADR-0046`-konform); ein neuer Tabellentyp und eine neue Hintergrund-Goroutine vergrößern die Betriebsfläche des Prozesses; `Assembler.tables` braucht eine Synchronisation, die es heute nicht hat, weil zwei Goroutinen darauf zugreifen; ein offener Antrag über einen Prozess-Neustart hinweg braucht eine definierte Fortsetzungsregel (bleibt `pending`, wird beim nächsten Boot/Poll erneut abgeholt) |

## Konsequenzen

- Positiv: `ADR-0018`s und `ADR-0046`s Grundsatz — schreibende/aktionsauslösende
  SQL-Zugriffe laufen ausschließlich über Inbound Ports, keine
  Domänenlogik-Duplikation in SQL — bleibt vollständig erhalten, auch für
  den bisher ungelösten Rest (`Enable`/`Disable`/`ACK` über SQL).
- Positiv: Der zweite Fund (`Assembler.tables` ohne Live-Reload) wird nicht
  separat gelöst, sondern fällt aus der gewählten Struktur ab — derselbe
  Prozess, der den Port aufruft, trägt die Bindung im selben Zug in seinen
  laufenden Zustand nach.
- Positiv: Der Boot-Wechsel von „`Assembler.tables` ausschließlich aus
  `CDC_TABLES`" auf „aus `cdc.source_table` via `TableActivationPort.List`,
  `CDC_TABLES` nur als Erstaktivierungs-Seed" schließt eine zweite,
  unabhängig von dieser ADR bestehende Inkonsistenz: Ein Prozess-Neustart
  ohne jede SQL-Aktivierung würde heute jede außerhalb von `CDC_TABLES`
  aktivierte Tabelle beim Hochfahren wieder verlieren.
- Negativ: `cdc.enable_table(...)`/`cdc.disable_table(...)` sind nach
  dieser Entscheidung **nicht** synchron im Sinn von „nach Rückkehr aktiv"
  — das weicht vom naheliegenden Lesart des Lastenheft-Beispiels
  (`SELECT cdc.enable_table(...)`) leicht ab; `LH-FA-CFG-001`s Happy Path
  („dann werden fortan Änderungen erfasst") verlangt aber keine
  Synchronität, nur Wirksamkeit, und ist damit erfüllt.
- Negativ: Zwei neue Komponenten (Antrags-Tabelle, Administrations-Goroutine)
  und ein Synchronisationsbedarf am `Assembler`, den es vorher nicht gab —
  Umsetzungs-Aufwand für den folgenden Implementer-Zug, kein Zusatzrisiko
  für bestehende Pfade (Store-, Heartbeat-, WAL-Retention-Goroutinen bleiben
  unverändert).
- Folgepflicht: `spec/architecture.md`s Sequenzdiagramm zu
  [`LH-FA-CFG-001.a`](../../../spec/architecture.md) zeigt heute
  „Administrator über SQL/CLI" undifferenziert als direkten Aufrufer von
  `EnableTableUseCase` — das gilt unverändert für CLI (eigener Go-Prozess,
  ruft den Port direkt auf, analog zu `register-consumer`), aber nicht mehr
  für SQL (Antrag → Queue → Hintergrund-Goroutine → Port, asynchron). Das
  Diagramm braucht eine Planner-/Architect-Korrektur, die beide Pfade
  trennt — derselbe offene Folgeauftrag, den bereits `ADR-0046` für den
  Lese-Pfad benannt hat, hier um den Schreib-Pfad erweitert.
- Folgepflicht: Die konkrete Form der Antrags-Tabelle (Spalten, Name,
  Retention der erledigten Anträge), das Antrags-Format für
  `cdc.disable_table`/eine künftige SQL-ACK-Funktion sowie die
  Synchronisationsstrategie am `Assembler` (Mutex vs. Kommando-Kanal in
  die bestehende `Consume`-Schleife) sind Gegenstand des Pflichtenhefts
  und der Slices, die diese ADR umsetzen — nicht dieser Entscheidung.
- Folgepflicht: `LH-FA-SST-003` (CLI-Diagnose) braucht keinen Teil dieser
  Entscheidung — ein Status-/Diagnose-Befehl liest nur (`cdc.active_tables`,
  `cdc.consumer_status`, Heartbeat), er schreibt nichts und braucht keinen
  Reload-Mechanismus. Er ist unabhängig von dieser ADR planbar.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (Review-Prüfpflicht, wie `ADR-0018`/`ADR-0046`) | Jede neue SQL-Funktion im `functions:`-Knoten des Schemamodells, die `cdc.source_table`, `cdc.schema_version` oder eine Publication direkt schreibt (`INSERT`/`UPDATE`/`DELETE`/`ALTER PUBLICATION` außerhalb der Antrags-Tabelle), verstößt gegen diese Entscheidung — sie darf ausschließlich einen Antrags-Datensatz schreiben und `pg_notify` senden | kein Gate — `.a-check.yml` kennt nur Go-Globs, SQL-Layer-Edges sind ihm unsichtbar |
| — (Code-Review) | `mapper.Assembler` greift von mehr als einer Goroutine auf `tables` zu, sobald die Administrations-Goroutine existiert — ohne Mutex/Kommando-Kanal ist das ein Data Race (`go test -race`, `ADR-0030`) | `make test` (bestehendes Ziel, Race-Detector als Umsetzungs-Detail des folgenden Slices) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Wird eine technische Brücke eingeführt, die SQL-Objekten einen echten,
synchronen Port-Aufruf erlaubt (FDW, `dblink`, Extension mit Prozess-/
Socket-Zugriff) — derselbe Trigger, den `ADR-0046` bereits nennt — dann
Re-Evaluierung, ob die Antrags-Queue noch nötig ist oder ob
`cdc.enable_table(...)`/`cdc.disable_table(...)` synchron über eine solche
Brücke laufen sollen. Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-13 | Accepted — Architect-Entscheidung vor dem Schneiden der Feature-Welle „Verwaltungsfunktionen — SQL-Administration & CLI-Diagnose"; löst die von `ADR-0046` bewusst offen gelassene Frage sowie den zweiten Fund der Fork-Recherche | [`docs/plan/planning/observations/BEO-PGC/verwaltung-keine-sql-administration/observation.md`](../planning/observations/BEO-PGC/verwaltung-keine-sql-administration/observation.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0050` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
