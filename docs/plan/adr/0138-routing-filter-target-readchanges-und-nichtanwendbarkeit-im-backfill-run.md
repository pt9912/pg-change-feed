# ADR-0138: Routing — Filter `target` am gRPC-`ReadChanges`, Nichtanwendbarkeit einer Routing-Regel im Backfill-Run (Schärft ADR-0137)

**Status:** Accepted

**Datum:** 2026-10-01

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch die
Vorab-Bedingungen V1 und V2 des Planners in `welle-routing`)

**Bezug:** [`LH-FA-CFG-008`](../../../spec/lastenheft.md) (Routing auf
Zustellziele), [`LH-FA-SST-006`](../../../spec/lastenheft.md) (Gleichwertigkeit
der Zugriffswege), [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Backfill),
[`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Negative: keine still
unveränderte Auslieferung),
[ADR-0137](0137-routing-zustellziele-persistiertes-ziel-label.md)
(Teilfrage 4, 5, 6 — Haupt-Bezug; **nicht geändert**, nur ergänzt),
[ADR-0131](0131-grpc-readchanges-zehnter-rpc.md) (zehnter RPC),
[ADR-0133](0133-tabellen-granulare-filterung-grpc-sse.md) (Filterform
`schema`/`table`; Muster Teilfrage 2),
[ADR-0117](0117-backfill-run-fehlerklasse-schema.md) (Klasse `schema` im Run),
[ADR-0112](0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(Folgepflicht 7)

**Schärft:** [`SPEC-031`](../../../spec/pflichtenheft.md#spec-031--grpc-verwaltungs-api-dienst-rpcs-nachrichtenschema-fehlercodes)
(Zeile `ReadChanges`: Request um `target` ergänzt),
[`SPEC-029`](../../../spec/pflichtenheft.md#spec-029--cdcbackfill_run-und-cdcbackfill_status-run-zustand-eines-backfills)
(Run-Fehlerklasse `schema` gilt auch für eine Routing-Regel)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0137`](0137-routing-zustellziele-persistiertes-ziel-label.md) ist
`Accepted` und unberührbar (`AGENTS.md` §3.5). Zwei Stellen seines Wortlauts
tragen die Umsetzung nicht; beide sind Lücken, keine Widersprüche, und werden
hier per Ergänzung geschlossen.

- **V1.** Teilfrage 5 listet je Weg einen Parameter `target`, aber nicht den
  zehnten RPC `ReadChanges` des `Administration`-Service (`ADR-0131`: derselbe
  Use Case wie `GET /changes`). Ohne `target` dort wäre ein Zugriffsweg, der
  `schema`/`table` filtern kann, für das Routing blind
  (`LH-FA-SST-006`).
- **V2.** `ADR-0137` Teilfrage 6 sagt, Backfill-Changes durchliefen dieselbe
  Auswertung; Teilfrage 4 bildet `ErrRoutingNotApplicable` im Erfassungspfad
  auf `schema` ab. Was ein **Run** bei einer nicht anwendbaren Routing-Regel
  tut, steht nirgends: `ADR-0117` gilt dem Wortlaut nach für
  Transformationsregeln.

### Befunde (erprobt am Quelltext, Stand HEAD `26392b65`)

- `ReadChangesRequest` in `proto/cdc/administration/v1/administration.proto`
  trägt die Felder 1 bis 6 (`source`, `schema`, `table`, `from`, `to`,
  `limit`); die Feldnummer 7 ist frei (gelesen, `grep` auf die Nachricht).
- `ADR-0131` Teilfrage 3 legt `ReadChangesRequest` als die sechs Felder fest;
  ein siebtes Feld ist eine additive Erweiterung der Nachricht, kein neuer RPC.
- `ADR-0117` Festlegung 2: die Nichtanwendbarkeit ist eine Eigenschaft von
  Regelmenge und Spaltenmenge, einmal je Run vor der ersten Zeile prüfbar.
  Für `when.column` (`ADR-0137` Teilfrage 4: „fehlt in der Relation“) gilt
  dieselbe Form: Prüfgegenstand ist die Existenz einer Spalte, kein
  Zeilenwert. Dass ein **abwesender Wert** keine Nichtanwendbarkeit ist
  (`ADR-0137` Teilfrage 4), bleibt unberührt.

### Konstraints

- Klassenmenge geschlossen (`ADR-0023`, `SPEC-008`): keine achte Klasse.
- Run-Fehler sind run-lokal (`ADR-0111` Teilfrage 5): weder Heartbeat-Zustand
  noch Halt des Erfassungspfads.
- `ChangeRecord` bleibt dreizehnfeldrig (`ADR-0131`); das Label ist nach
  `ADR-0137` Teilfrage 5 nicht Teil der Nachrichten der Wege.

## Entscheidung

Wir wählen **ein additives Feld `string target = 7` in `ReadChangesRequest`
(leer = kein Filter, Konjunktion mit `schema`/`table`) und die Klasse `schema`
für eine im Run nicht anwendbare Routing-Regel, run-lokal, vor der ersten
Kopie, wie `ADR-0117` für Transformationsregeln**. V3 bleibt offen (unten).

### Festlegung 1 (V1) — `ReadChangesRequest.target = 7`

- `ReadChangesRequest` erhält `string target = 7;` (Feldnummer 7: nächste freie
  nach 6). Leer = kein Filter; gesetzt = nur Changes mit `route_target = target`;
  Konjunktion mit `schema`/`table`. Die Filterform ist dieselbe wie in
  `ADR-0137` Teilfrage 5 für `GET /changes`; der Use Case `ReadChangesUseCase`
  (`ADR-0131`) erhält den Filter einmal und bedient beide Wege.
- Wire-kompatibel additiv: ein alter Client sendet das Feld nicht (leer = kein
  Filter, altes Verhalten); ein alter Server ignoriert ein unbekanntes Feld
  und lieferte ungefiltert (*hergeleitet* aus den Proto3-Regeln, nicht gefahren).
- `ChangeRecord` bleibt unverändert (`ADR-0137` Teilfrage 5: Label nicht in den
  Nachrichten); ein Filter bleibt ohne Anzeige des Labels lesbar über die
  SQL-Sicht.
- Kein Trennzeichen- oder Alphabet-Sonderfall: ein `target` außerhalb des
  Alphabets (`ADR-0137` Teilfrage 1) liefert eine leere Liste, kein Fehler
  (wie ein unbekanntes `schema`, `ADR-0131`, `LH-FA-REA-006`). *Erwartung*,
  vom Slice `slice-routing-lesewege` zu belegen.

### Festlegung 2 (V2) — Nichtanwendbarkeit einer Routing-Regel im Run

- Ist eine Regel des Routing-Regelstands der Tabelle auf die Spalten des
  Snapshots nicht anwendbar (`when.column` fehlt in `TableSnapshot.Columns()`),
  endet der Run `failed` mit Klasse `schema`: einmal je Run, vor der
  Schreibtransaktion und der ersten Zeile, `error_message` beginnt mit
  `schema: ` und nennt Regelname und Spalte; keine Change entsteht, der Snapshot
  ist geschlossen. Run-lokal (`ADR-0111`): kein Heartbeat-Zustand, kein Halt
  des Erfassungspfads.
- Dieselbe Prüffunktion der Domäne wie im Erfassungspfad (`ADR-0137`
  Teilfrage 6, Folgepflicht 6): keine zweite Implementierung.
- Abhilfe: Regel ändern/entfernen (`cdc.remove_route`), neuer Antrag.
- Nicht still roh kopieren und nicht Zeilen auslassen (`ADR-0117` Alternativen D/E
  gelten unverändert; hier nicht neu bewertet).

### Offen, nicht entschieden — V3

Ob und wie `ErrRoutingNotApplicable` am laufenden System **erreichbar** ist
(naheliegende Ursache „Spalte entfernt“ endet vermutlich im
`ErrIncompatibleSchemaChange`-Pfad, *hergeleitet*, nicht gemessen), und ob der
Abhilfe-Weg „`cdc.remove_route` und Neustart“ greift, klärt die Messung in
`slice-routing-kern-label` (Unit-Ebene) und `slice-routing-e2e`
(Systemebene); ein unerreichbarer Fall verlangt ein Architect-Verdikt vor dem
Negative-Beleg. Diese ADR entscheidet V3 nicht; Festlegung 2 gilt
unabhängig davon, weil sie die Behandlung **für den Fall der Nichtanwendbarkeit**
festlegt, nicht ihre Erreichbarkeit.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (kein `target` am RPC, Run-Verhalten offen) | keine Änderung | `LH-FA-SST-006` für das Routing verletzt; die Run-Klasse bliebe Auslegung des Implementers (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`) |
| B — `target` am RPC über einen eigenen RPC `ReadRouted…` | kein Eingriff in die Nachricht | elfter RPC für einen Filter; weicht von der Form `schema`/`table` ab |
| **C — `target = 7` additiv; Run-Klasse `schema` (gewählt)** | Muster `ADR-0133`/`ADR-0117`; kein Schema-, Port- oder Modell-Diff; gleiche Ursache, gleiche Klasse in beiden Pfaden | `ReadChangesRequest` wächst auf sieben Felder |
| D — Run-Klasse `configuration` für Routing-Regeln | Regel ist Konfiguration | dieselbe Ursache hieße je Pfad anders; Gründe wie `ADR-0117` Option B |
| E — Run überspringt die Regel, `route_target = NULL` | Run läuft durch | stilles Fehl-Routing des Bestands, `LH-FA-CFG-008` Negative widersprochen; `ADR-0117` Option D |

## Konsequenzen

- Positiv: alle Lesewege filtern nach Ziel; kein Sonderpfad im Run.
- Negativ: `ReadChangesRequest` ändert sich (additiv); `make proto-generate`
  und `generated-sync` laufen im Slice.
- Folgepflicht: `SPEC-031`-Zeile `ReadChanges` und `SPEC-029`/`SPEC-008`
  um die Ergänzungen (`slice-routing-spec-nachzug` bzw. Fortschreibung dort);
  `tools/harness/grpcadminclient` um `target` erweitern;
  SDK-/Handbuch-Zeile (`ADR-0137` Folgepflicht 8) bleibt eigener Schritt.

## Fitness Function (falls maschinell prüfbar)

Alle Zeilen sind **erwartet** (nicht gefahren): bis zur Umsetzung wurde keine
Mutation erprobt; Stellen und Instanz der Mutationen benennt der Slice, der
sie fährt.

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (erwartet) | gRPC-`ReadChanges` mit `target` liefert nur Changes dieses Ziels; leeres `target` liefert alle; Konjunktion mit `table` | `make test` |
| Go-Test/Store-Test (erwartet) | Run mit nicht anwendbarer Routing-Regel endet `failed`, `error_message` beginnt mit `schema: `, keine Change | `make test`, `make test-store` |
| Tooling | Proto-Erzeugnis byte-gleich zum gepinnten Generator | `make generated-sync` |

## Re-Evaluierungs-Trigger

Die Messung zu V3 ergibt, dass die Nichtanwendbarkeit einer Routing-Regel am
System nie erreichbar ist (Festlegung 2 wäre dann nur Vorsorge), oder
`ReadChangesRequest` erhält weitere Filter, die eine gemeinsame Form
verlangen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-01 | Accepted (Kurz-ADR auf Auftrag des Hauptlaufs; Schärft ADR-0137) | `welle-routing`, Vorab-Bedingungen V1/V2 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen entstehen als neue ADR mit `Supersedes ADR-0138`.
