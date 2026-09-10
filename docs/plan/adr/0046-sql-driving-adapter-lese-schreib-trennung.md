# ADR-0046: SQL-Driving-Adapter — Lese-Views direkt, schreibende Funktionen über Inbound Ports

**Status:** Accepted — Supersedes [`ADR-0018`](0018-sql-driving-adapter.md)

**Datum:** 2026-09-10

**Autor:** pt9912

**Bezug:** [`LH-FA-SST-002`](../../../spec/lastenheft.md),
[`LH-FA-REA-002`](../../../spec/lastenheft.md)…006,
[`LH-FA-ADM-001`](../../../spec/lastenheft.md),
[ADR-0018](0018-sql-driving-adapter.md)

**Schärft:** [`ARC-005`](../../../spec/architecture.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0018`](0018-sql-driving-adapter.md) entschied: „SQL-Funktionen/Views
sind Driving Adapter und rufen Inbound Ports auf. Businesslogik wird nicht
in SQL dupliziert." Ein Review fand drei gelieferte Lese-Views
(`cdc.active_tables`, `cdc.consumer_status`, `cdc.changes`,
`tools/schema/nacharbeit-views.sql`) im Wortlaut-Verstoß: <!-- d-check:status-provenance -->
Sie selektieren direkt aus `cdc.source_table`, `cdc.consumer`,
`cdc.consumer_position`, `cdc.transaction`, `cdc.change` — denselben
Tabellen, die der `PostgresChangeStoreAdapter` (`ARC-006`, Driven) verwaltet
— ohne einen Inbound Port aufzurufen. Weder der Plankopf noch die
Implementer-Commits referenzierten `ADR-0018`, obwohl eine frühere
Closure-Notiz genau diese Lieferung dorthin verwiesen hatte
(Review-Report: [`docs/reviews/review-slice-010.md`](../../reviews/review-slice-010.md) <!-- d-check:status-provenance -->,
Finding F-1, HIGH).

**Physikalischer Befund, der `ADR-0018` beim Schreiben nicht auflöste:** Eine
PostgreSQL-`VIEW` ist eine gespeicherte `SELECT`-Anweisung. Sie hat kein
Sprachmittel, einen Go-Prozess (Inbound Port) synchron aufzurufen — dafür
bräuchte es eine Brücke (FDW, `dblink`, eine Erweiterung mit Netz-/
Socket-Zugriff), die in diesem Repo nirgends existiert, nirgends entschieden
ist und die für drei reine Lesezugriffe eine neue, schwergewichtige
RPC-Infrastruktur wäre. `ADR-0018`s Wortlaut „rufen Inbound Ports auf" nennt
kein Mittel für den Views-Teil von „SQL-Funktionen/**Views**" — diese Lücke
ist beim Schreiben der ADR nicht aufgelöst worden, nicht bewusst offen
gelassen.

**Was `ADR-0018` tatsächlich schützen wollte:** Die verworfene Option A
("Logik in PL/pgSQL-Funktionen") nennt den Grund explizit: „Duplikat der
Domänenlogik in einer zweiten Sprache; nicht domänentestbar; Wartung an
zwei Orten." Das ist eine Aussage über **Entscheidungslogik** — Code, der
eine Domänenregel auswertet oder einen Zustand ändert —, nicht über
Projektion. Die drei gelieferten Views enthalten keine solche Logik: Sie
joinen und projizieren bereits persistierte, bereits validierte Zeilen;
Bereich/Limit/Filter kommen vom aufrufenden SQL-Client (`WHERE`/`LIMIT` auf
der View), nicht aus der View selbst (`tools/schema/schema.yaml:16-27`
benennt das explizit — auch die bewusste Grenze, dass `active_tables` die
Publication-Mitgliedschaft nicht nachbildet, bleibt beim
`TableActivationPort`).

**Was unverändert bleibt:** `spec/architecture.md` §4 zeigt für
schreibende/aktionsauslösende SQL-Zugriffe — Consumer-ACK
(`LH-FA-CON-004`) und Tabelle aktivieren (`LH-FA-CFG-001.a`) — explizit den
Weg über `AcknowledgeConsumerUseCase`/`EnableTableUseCase` (Inbound Ports).
Beide tragen echte Domäneninvarianten (ACK „regulär nur vorwärts",
Replica-Identity-Prüfung) — genau die Business-Logik, deren
Zweitimplementierung in SQL `ADR-0018` verhindern wollte. Diese Pflicht
ändert diese ADR nicht.

**Was diese ADR nicht löst:** Für die noch nicht implementierten
schreibenden SQL-Funktionen (Enable/Disable/ACK — als
`ADR-0018`/`ADR-0019`-Rest-Arbeit benannt <!-- d-check:status-provenance -->)
besteht dasselbe physikalische Problem wie bei den Views: Eine reine
SQL-Funktion kann einen Go-Inbound-Port ebenfalls nicht synchron aufrufen.
Diese Frage bleibt **offen** und ist nicht Gegenstand dieser Entscheidung
(siehe Konsequenzen, Folgepflicht).

## Entscheidung

Wir wählen **C — Rollenspaltung nach Datenrichtung**: Reine Lese-SQL-Views
dürfen Driven-Adapter-Tabellen direkt lesen, sofern sie ausschließlich
Projektion/Join über bereits persistierte, bereits validierte Daten liefern
und keine Domänenregel auswerten (Bereich/Limit/Filter des aufrufenden
SQL-Clients zählt nicht als Geschäftslogik). Schreibende oder
aktionsauslösende SQL-Funktionen — jede SQL-Schnittstelle, die einen
Zustand ändert oder eine Autorisierungs-/Invarianten-Entscheidung trifft —
rufen weiterhin ausschließlich über Inbound Ports; diese Pflicht aus
`ADR-0018` gilt unverändert fort.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — `ADR-0018` wörtlich durchsetzen (auch Lese-Views rufen Inbound Ports auf) | keine Ausnahme, volle ADR-Treue | physisch nicht umsetzbar ohne neue FDW-/`dblink`-/Extension-Brücke, die nirgends entschieden oder gebaut ist; erzwingt schwergewichtige neue Infrastruktur für drei reine Lesezugriffe — dieser Pfad scheitert an der Physik, nicht an der Umsetzung |
| B — SQL-Lesezugriffe streichen (erneut Option B aus `ADR-0018`) | kein Konflikt, keine Ausnahme | verletzt [`LH-FA-SST-002`](../../../spec/lastenheft.md) und [`LH-FA-REA-002`](../../../spec/lastenheft.md)…006 — derselbe Grund, aus dem `ADR-0018` diese Option bereits verwarf |
| **C — Rollenspaltung: Lese-Views direkt, schreibende Funktionen über Ports (gewählt)** | physisch umsetzbar mit dem heutigen PostgreSQL-Only-Stack; hält `ADR-0018`s eigentlichen Schutz (keine Domänenlogik-Duplikation) vollständig aufrecht, weil reine Projektion keine Logik trägt; schreibende Pfade bleiben unverändert gebunden | zwei Kategorien statt einer Einheitsregel; künftige SQL-Objekte müssen bei Anlage eingeordnet werden (reine Projektion vs. Aktion), Grenzfälle bleiben Review-Prüfpflicht |

## Konsequenzen

- Positiv: Die drei bereits gelieferten Views (`active_tables`,
  `consumer_status`, `changes`) sind konform; keine Rückbau-Pflicht, kein
  Implementer-Fix-Zug nötig.
- Positiv: Der eigentliche Schutz von `ADR-0018` — keine
  Domänenlogik-Duplikation in SQL — bleibt für den Teil vollständig
  erhalten, der ihn tatsächlich betrifft: schreibende/aktionsauslösende
  SQL-Funktionen (Enable/Disable/ACK).
- Negativ: Zwei-Kategorien-Regel statt Einheitsregel. Eine künftige View
  mit einer berechneten Spalte, die einen Business-Zustand ausdrückt (z. B.
  eine Freigabe-Entscheidung), ist kein Grenzfall dieser ADR mehr, sondern
  fällt unter die schreibende/entscheidende Kategorie — Review-Prüfpflicht,
  kein Gate deckt das (wie schon `ADR-0018` feststellte; `.a-check.yml`
  kennt nur Go-Globs).
- Negativ / offen: Für die noch nicht implementierten schreibenden
  SQL-Funktionen bleibt ungeklärt, wie eine reine SQL-Funktion physisch
  einen Inbound Port aufrufen soll — dasselbe physikalische Problem wie bei
  den Views, nur dass hier keine Ausnahme greift. Diese Design-Lücke ist
  nicht Gegenstand dieser ADR und bleibt offen für den Slice, der das
  `ADR-0018`/`ADR-0019`-Rest umsetzt.
- Folgepflicht: `spec/architecture.md` §4, Sequenzdiagramm zu
  [`LH-FA-REA-002`](../../../spec/lastenheft.md) ("Consumer über CLI/SQL"
  → `ReadChangesUseCase` → `ChangeStorePort` → `PostgresChangeStoreAdapter`)
  zeigt für den SQL-Kanal einen Port-Umweg, den diese ADR für reine
  Lese-Views nicht mehr fordert — das Diagramm braucht eine
  Planner-/Architect-Korrektur (SQL-Kanal in der Lese-Sequenz als
  Direktzugriff kennzeichnen oder auf den CLI-Kanal beschränken), sonst
  zeigt die Sicht einen Ablauf, der für SQL nicht mehr gilt. Der
  betroffene Slice-Plankopf (`Bezug:`) referenziert bislang nur
  [`ADR-0009`](0009-change-store-outbound-port.md); Planner/Implementer
  tragen `ADR-0046` nach, sobald der Implementer-Zug fortgesetzt wird.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (Review-Prüfpflicht, wie `ADR-0018`) | SQL-Objekte unter `tools/schema/nacharbeit-*.sql` bzw. im `views:`-Knoten des Schemamodells: eine Lese-View enthält kein `INSERT`/`UPDATE`/`DELETE`, keine Funktionsdefinition mit Seiteneffekt und keine `WHERE`-Klausel, die eine Autorisierungs- oder Domänenentscheidung kodiert (statt sie an den aufrufenden Client zu übergeben) | kein Gate — `.a-check.yml` kennt nur Go-Globs, SQL-Layer-Edges sind ihm unsichtbar |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Wird eine technische Brücke eingeführt, die SQL-Objekten einen echten,
synchronen Port-Aufruf erlaubt (FDW, `dblink`, Extension mit Prozess-/
Socket-Zugriff) — dann Re-Evaluierung, ob die Rollenspaltung noch nötig ist
oder ob alle SQL-Driving-Adapter (Lese wie Schreib) wieder einheitlich über
Inbound Ports laufen sollen. Sonst permanent — die physikalische Grenze
reiner SQL-Views besteht unabhängig vom Datenbestand.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-10 | Accepted — Anlass: Review-Finding F-1 (Konflikt-Pfad, Modul 8); Verdikt 2 (Folge-ADR `supersedes`) <!-- d-check:status-provenance --> | [`docs/reviews/review-slice-010.md`](../../reviews/review-slice-010.md) <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0046` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
