# ADR-0139: Routing — Regelstand im Backfill-Run fail-closed, `target` außerhalb des Alphabets liefert leer (Schärft ADR-0137, ADR-0138)

**Status:** Accepted

**Datum:** 2026-10-01

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch die Architect-Fragen
A-1 und A-2 des Reviews `review-slice-routing-spec-nachzug`)

**Bezug:** [`LH-FA-CFG-008`](../../../spec/lastenheft.md) (Routing auf
Zustellziele), [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Backfill,
Fail-closed), [`LH-FA-SST-006`](../../../spec/lastenheft.md) (Gleichwertigkeit
der Zugriffswege), [`LH-FA-REA-006`](../../../spec/lastenheft.md) (Boundary:
kein Treffer ist kein Fehler),
[ADR-0137](0137-routing-zustellziele-persistiertes-ziel-label.md)
(Teilfragen 1, 5, 6; Entscheidung 6 — **nicht geändert**, nur ergänzt),
[ADR-0138](0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
(Festlegung 1 letzter Punkt, Festlegung 2 — **nicht geändert**, nur ergänzt),
[ADR-0111](0111-backfill-bestand-snapshot-bulk-copy.md) (Teilfrage 4),
[ADR-0113](0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[ADR-0117](0117-backfill-run-fehlerklasse-schema.md),
[ADR-0118](0118-backfill-umschreiben-im-snapshot-fenster.md),
[ADR-0081](0081-changes-lesen-ueber-die-http-api.md) (Teilfrage 4),
[ADR-0133](0133-tabellen-granulare-filterung-grpc-sse.md)

**Schärft:** [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) (Absätze
„Fail-closed vor dem Commit“ und „Ziel der Backfill-Changes“),
[`SPEC-020`](../../../spec/pflichtenheft.md#spec-020--grpc-live-change-stream-nachrichtenschema-rpc-name-stream-semantik),
[`SPEC-021`](../../../spec/pflichtenheft.md#spec-021--http-server-sent-events-für-den-live-change-stream),
[`SPEC-022`](../../../spec/pflichtenheft.md#spec-022--http-api-changes-lesen-get-changes),
[`SPEC-031`](../../../spec/pflichtenheft.md#spec-031--grpc-verwaltungs-api-dienst-rpcs-nachrichtenschema-fehlercodes)
(Zeile `ReadChanges`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0137` und `ADR-0138` sind `Accepted` und unberührbar (`AGENTS.md` §3.5). Das
Review `review-slice-routing-spec-nachzug` fand zwei Lücken, die keiner der beiden
Texte entscheidet; die Spec trägt sie als offen (F-1) bzw. als unbedingte Zusage
über die ADR hinaus (F-2).

- **A-1.** `ADR-0137` Teilfrage 6 gibt Backfill-Changes „das Label des Regelstands
  zum Run“. Der Absatz „Fail-closed vor dem Commit“ in `LH-FA-CAP-009.a` führt als
  zu prüfenden Stand nur den Ausschluss- und den Transformationsregelstand.
  Ungeklärt: gilt die Prüfung auch für den Routing-Regelstand, gegen welchen Stand
  wird verglichen, welche Fehlerklasse.
- **A-2.** `ADR-0138` Festlegung 1 sagt für `ReadChanges`: ein `target` außerhalb
  des Alphabets liefert eine leere Liste, kein Fehler — als *Erwartung*. Die Spec
  setzt den Satz unbedingt für alle Lesewege. Offen: gilt er überall, und was
  geschieht bei einem `target` mit U+0000?

### Befunde (Stand HEAD `b1b34ef3`)

Gelesen am Quelltext, **nicht gefahren** (keine Laufzeit-Erprobung in dieser ADR):

- Der Backfill-Dienst (`internal/application/usecase/backfill/service.go`,
  Kommentar zu `copyBlocks`) liest den Transformations-Regelstand einmal nach dem
  Öffnen des Snapshots („Stand des Runs“), prüft ihn auf Anwendbarkeit
  (`checkRulesApplicable`, Klasse `schema`), liest Ausschluss- und Regelstand je
  Block und unmittelbar vor dem Commit neu und vergleicht als Menge gegen den
  Start-Stand; eine Abweichung ist `ErrExclusionStateChanged`, abgebildet auf
  `configuration`. Eine Regel, die zwischen zwei Lesungen gesetzt und
  zurückgenommen wird, bleibt unsichtbar (kein Verlauf) — eine bekannte, dort
  benannte Grenze.
- Die Lese-Abfrage (`internal/adapters/driven/postgresstorage/queries/queries.go`,
  Zeile 70) bindet `schema` als Parameter `$4::text` und vergleicht gleich; es gibt
  weder im HTTP-Handler (`readchanges.go`) noch im Use Case
  (`readchanges/service.go`: nur `Source == ""` wird abgelehnt) eine Alphabet- oder
  Syntaxprüfung für `schema`/`table`. Unbekanntes `schema`/`table` liefert `200`
  mit leerer Liste (`ADR-0081` Teilfrage 4 Option B, `LH-FA-REA-006`); nur
  Parameterform (unbekannter Parameter, unlesbare Zahl) und Bereichsfehler enden
  `400`.
- *Hergeleitet, nicht gemessen:* PostgreSQL lehnt einen Text-Parameter mit U+0000
  ab (SQLSTATE 22021). Ein `schema`/`table` mit U+0000 endet deshalb heute auf den
  SQL-gestützten Wegen im Store-Fehler (Klasse `storage`, `500`/`Internal`), auf
  Live-Stream und SSE (Vergleich im Speicher) ergibt es leer. Das ist Bestand und
  wird hier nicht geändert.

## Entscheidung

Wir wählen **A-1: Option (a) — die Fail-closed-Prüfung gilt auch für den
Routing-Regelstand, Vergleichsstand ist der Stand der Lesung nach dem Öffnen des
Snapshots, Fehlerklasse `configuration`; A-2: ein `target` außerhalb des Alphabets
(einschließlich jedes `target` mit U+0000) liefert auf allen Lesewegen eine leere
Antwort ohne Fehler, gesichert durch eine einzige Prüfung im gemeinsamen Use Case
statt durch das Verhalten des Speichers.**

### Festlegung 1 (A-1) — Routing-Regelstand im Run fail-closed

- **Gegenstand.** Der Routing-Regelstand der Tabelle (die gefaltenen
  `applied`-Zeilen von `set_route`/`remove_route`, `ADR-0137` Teilfrage 6) wird
  einmal nach dem Öffnen des Snapshots gelesen, auf Anwendbarkeit geprüft
  (`ADR-0138` Festlegung 2, Klasse `schema`, vor der ersten Zeile) und ist von da an
  der **Stand des Runs** („Regelstand zum Run“). Alle Blöcke bestimmen
  `route_target` mit diesem Stand.
- **Vergleich.** Jeder weitere Block und der Zustand unmittelbar vor dem Commit
  lesen den Routing-Regelstand neu und vergleichen ihn als Menge gegen den
  Stand des Runs, derselbe Mechanismus wie für Ausschluss- und Transformationsstand.
  Eine Abweichung — auch ein nicht lesbarer Stand — rollt den Run zurück: `failed`,
  Klasse `configuration`, Grund im Fehlertext. Vergleichsstand ist weder der Stand
  bei Annahme des Antrags noch ein Stand je Block.
- **Reihenfolge.** Erst Anwendbarkeit (`schema`), dann Kopie mit Vergleich je
  Block (`configuration`). Die zwei Klassen kollidieren nicht: `schema` gilt der
  Nichtanwendbarkeit des Start-Stands (`ADR-0138`), `configuration` einer
  Änderung des Standes während des Runs.
- **Abhilfe.** Ein neuer Antrag; er nimmt den dann geltenden Stand.
- **Grenze (akzeptiertes Negativ).** Eine Routing-Regel, die zwischen zwei Lesungen
  gesetzt und wieder entfernt wird, bleibt unsichtbar, weil der Stand keinen
  Verlauf trägt — identisch zur bestehenden Grenze bei Ausschluss und
  Transformationen. Sie ändert kein Label, weil kein Block sie sieht.
- **Begründung.** Nicht der Leck-Schutz trägt die Prüfung: ein ausgeschlossener
  Wert gelangt über das Routing nicht in ein Bild, das deckt R3 (`ADR-0137`
  Teilfrage 6), das Label ist kein Teil des Row Images. Sie ist nötig, weil
  `ADR-0137` Teilfrage 6 **einen** Regelstand je Run zusagt und das Label nach
  Entscheidung 6 nicht rückwirkend neu bestimmt wird: ohne Prüfung trüge ein Abzug
  Labels aus zwei Ständen, und der Mischstand wäre nach dem Commit dauerhaft und
  unmarkiert; ein Consumer, der nach `target` filtert, sähe für dieselbe Tabelle
  einen halb nach altem, halb nach neuem Stand geteilten Bestand. Die Prüfung hat
  die Form, die der Dienst für die beiden anderen Stände schon trägt; es entsteht
  kein neuer Mechanismus, nur ein weiterer Stand im selben Vergleich.

### Festlegung 2 (A-2) — `target` außerhalb des Alphabets

- Ein `target`, das nicht dem Alphabet des Zielnamens
  (`[a-z0-9][a-z0-9_-]{0,62}`, `ADR-0137` Teilfrage 1) entspricht — auch jeder Wert
  mit U+0000, jeder Wert über 63 Zeichen, mit Großbuchstaben oder Trennzeichen —,
  liefert auf **allen** Lesewegen (`GET /changes`, gRPC-`ReadChanges`,
  gRPC-Live-Stream, SSE) eine leere Antwort: `200` mit `{"changes": []}` bzw. leere
  Liste bzw. ein Stream ohne Events, **kein** `400`/`InvalidArgument`. Auf dem
  SQL-Lesezugriff (`cdc.changes`, Spalte `route_target`) bestimmt der Aufrufer sein
  Prädikat selbst; die Zusage gilt dort nicht, hier trägt PostgreSQL die Semantik.
- **Mechanismus.** Der gemeinsame Use Case `ReadChangesUseCase` (bedient
  `GET /changes` und gRPC-`ReadChanges`, `ADR-0131`) prüft ein gesetztes `target`
  gegen das Alphabet mit derselben Prüffunktion der Domäne wie die Antragsprüfung
  (keine zweite Implementierung) und antwortet bei Verletzung mit der leeren Liste,
  **ohne** den Store anzufragen. Damit hängt „kein Fehler“ nicht am Verhalten von
  PostgreSQL gegenüber U+0000. Live-Stream und SSE vergleichen im Speicher; ein
  Wert außerhalb des Alphabets trifft nie ein vergebenes Label, die Antwort ist
  ohne weitere Prüfung leer.
- **Begründung.** Dieselbe Form wie bei `schema`/`table`: ein nicht treffender
  Filter ist `200` leer (`ADR-0081` Teilfrage 4 B, `LH-FA-REA-006`); ein Alphabet
  ist keine Syntaxprüfung der Parameterform wie „Zahl nicht lesbar“, sondern
  Namensraum eines vergebenen Labels. Ein Fehler für `target` allein würde die
  Filter-Dimensionen auf demselben Aufruf unterschiedlich behandeln
  (`schema=zzz` leer, `target=ZZZ` Fehler) und die Gleichwertigkeit der Wege
  (`LH-FA-SST-006`) nur durch eine Sonderform stützen. Der Preis — ein Tippfehler in
  `target` liefert still leer — ist derselbe wie bei `table` und akzeptiert.
- **Bestand `schema`/`table` mit U+0000.** Heute (hergeleitet) Store-Fehler auf den
  SQL-Wegen, leer auf den Streams. Diese ADR ändert das nicht; für `target` ist die
  Antwort einheitlich leer. Der Unterschied zwischen `target` und `schema`/`table`
  bei U+0000 ist ein einmaliger, harmloser Blindfleck: ein NUL im Filter ist
  keine legitime Eingabe, der Fehlerpfad (`500`) ist sichtbar statt still falsch.
  Eine Angleichung wäre eine eigene Änderung an `ADR-0081`/`ADR-0133` und wird nicht
  zu einem Vorgang gemacht.

### Offen, nicht entschieden — A-3

Ob die Lesart „Zustellung = Abruf/Abonnement eines Kanals“ (`ADR-0137`
Entscheidung 1) im Lastenheft geklärt werden soll, ist eine Frage an den
Auftraggeber, nicht an diese ADR. **Empfehlung:** ja, eine kurze Klarstellung in
`LH-FA-CFG-008` (Happy Path), in einem eigenen Commit vor dem umsetzenden Slice;
Begründung: im Konfliktfall gewinnt das Lastenheft, und eine einschränkende Lesart
sollte dort stehen, wo die Abnahme gegen den Text läuft. Keine Entscheidung hier.

## Verglichene Alternativen

**A-1 (Routing-Regelstand im Run):**

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (Spec bleibt offen) | keine Änderung | Plan von `slice-routing-backfill-pfad` setzt `configuration` bereits voraus, ohne dass eine Norm sie trägt; Auslegung des Implementers (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`) |
| B — keine Prüfung, Stand je Block, als Zusage benannt | kein Mehraufwand | Mischstand im Abzug, dauerhaft (Label nicht rückwirkend); „Regelstand zum Run“ wäre kein einzelner Stand |
| C — Stand bei Annahme des Antrags festgeschrieben | ein Run hätte den Stand der Entscheidung des Aufrufers | der Run kann lange `queued` warten; der Stand müsste am Run persistiert werden (Schema-Diff an `cdc.backfill_run`, `ADR-0113`-Rollenschnitt berührt) für einen Fall, den der Vergleich gegen den Start-Stand billiger löst |
| **D — Prüfung wie bei Ausschluss und Transformationen, Stand nach Snapshot-Öffnung, `configuration` (gewählt)** | vorhandener Mechanismus, ein Stand je Run, Run-lokal, kein Schema-Diff | je Lesung zusätzlich die `applied`-Zeilen der Routing-Antragsarten (Kosten wie bei den anderen Ständen); Set-und-Rücknahme zwischen zwei Lesungen unsichtbar |
| E — Klasse `schema` statt `configuration` | eine Klasse für alle Routing-Run-Fehler | `schema` meint Nichtanwendbarkeit gegen die Spalten (`ADR-0117`); ein Wechsel des Standes ist Konfigurationsänderung wie beim Ausschluss |

**A-2 (`target` außerhalb des Alphabets):**

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (ADR-0138-Wortlaut: nur `ReadChanges`, als Erwartung) | keine Änderung | andere Wege unentschieden, U+0000 offen |
| B — `400`/`InvalidArgument` bei Alphabet-Verletzung | Tippfehler sichtbar | `target` als einzige Filter-Dimension mit Wert-Prüfung; widerspricht `schema`/`table`-Verhalten (leer) und `LH-FA-REA-006`; Stream-Öffnung müsste bei SSE/gRPC-Stream einen Fehler vor dem Öffnen kennen |
| **C — leere Antwort überall, Prüfung im Use Case ohne Store-Aufruf (gewählt)** | einheitlich mit `schema`/`table`; U+0000 unabhängig vom Verhalten von PostgreSQL; eine Prüffunktion | Tippfehler still leer (wie `table`) |
| D — leer überall, Verlassen auf PostgreSQL für U+0000 | kein Code | U+0000 endet im Store-Fehler (hergeleitet), „kein Fehler“ wäre falsch |

## Konsequenzen

- Positiv: ein Regelstand je Run für alle drei Stände; `configuration` bedeutet im
  Run überall dasselbe; alle Lesewege antworten auf ein ungültiges `target`
  gleich.
- Negativ: eine zusätzliche Lesung des Routing-Regelstands je Block und vor dem
  Commit; ein Use-Case-Zweig mehr.
- Folgepflicht (Stand, nicht Chronik):
  - `spec/pflichtenheft.md`, `LH-FA-CAP-009.a`: Absatz „Fail-closed vor dem Commit“
    nennt den Routing-Regelstand; Absatz „Ziel der Backfill-Changes“ definiert
    „Regelstand zum Run“ als den nach dem Öffnen des Snapshots gelesenen Stand.
  - `spec/pflichtenheft.md`, `SPEC-020` (Zeile Request), `SPEC-021` (Zeile
    Query-Parameter), `SPEC-022` (Zeile Zustellziel), `SPEC-031` (Zeile
    `ReadChanges`): der Satz „leer, kein Fehler“ nennt U+0000 und den Mechanismus
    (Prüfung im Use Case); die SQL-Spalte bleibt ausgenommen.
  - Slice-Pläne: `slice-routing-backfill-pfad` (Fail-closed-Festlegung auf diese
    ADR stellen statt Annahme; Tests: `set_route` und `remove_route` zwischen zwei
    Blöcken enden `configuration`), `slice-routing-lesewege` (Use-Case-Prüfung,
    Tests für Wert außerhalb des Alphabets und U+0000 auf `GET /changes` und
    gRPC-`ReadChanges`, Streams), `slice-routing-spec-nachzug` (Fixrunde am
    Spec-Text; die DoD-Zeile „Review durchgeführt“ bleibt bis dahin offen),
    `welle-routing` (A-1/A-2 geschlossen, A-3 und V3 offen).

## Fitness Function (falls maschinell prüfbar)

Alle Zeilen sind **erwartet**: in dieser ADR wurde keine Mutation gefahren, keine
Aussage über eine Menge von Stellen ist erprobt. „Der Implementer fährt sie“ ist
eine Erwartung, keine Erprobung; Stellen und Instanz der Mutationen benennt der
Slice, der sie fährt.

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (erwartet) | Backfill-Dienst mit wechselndem Routing-Regelstand zwischen Block 1 und 2 und vor dem Commit endet `failed`, Klasse `configuration`; ohne Wechsel läuft er durch | `make test` |
| Store-Test (erwartet) | Run gegen reale PostgreSQL mit `set_route` zwischen zwei Blöcken endet `configuration`, keine Change sichtbar | `make test-store` |
| Go-Test (erwartet) | `ReadChangesService` mit `target` außerhalb des Alphabets (Großbuchstabe, 64 Zeichen, U+0000) antwortet leer und ruft den Store nicht auf | `make test` |
| Go-Test (erwartet) | `GET /changes` und gRPC-`ReadChanges` mit derselben Eingabe: leere Liste, kein Fehler | `make test` |

## Re-Evaluierungs-Trigger

Der Mischstand im Run wird anders erzeugt (zum Beispiel ein Routing, das
Row-Image-Werte ausgibt und damit nicht mehr nur ein Label ist), oder die
Gleichwertigkeitsprüfung über die Lesewege verlangt eine Fehlerform für
Filterwerte (dann auch für `schema`/`table`, `ADR-0081`).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-01 | Accepted (Kurz-ADR auf Auftrag des Hauptlaufs; Schärft ADR-0137, ADR-0138) | `review-slice-routing-spec-nachzug` F-1/F-2, A-1/A-2 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen entstehen als neue ADR mit `Supersedes ADR-0139`.
