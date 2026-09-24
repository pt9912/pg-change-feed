# ADR-0116: Backfill — Bedeutung und Reichweite der Schema-Version-Referenz einer Backfill-Change (Supersedes ADR-0111, Konsequenzen teilweise)

**Status:** Accepted — Supersedes [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
in **§Konsequenzen, „Akzeptierte Negative", genau einem Aufzählungspunkt**
(„Schema-Version-Verweis"); alles Übrige von `ADR-0111` bleibt in Kraft

**Datum:** 2026-09-24

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch den Befund F-8 des
Review-Reports zum Run-Use-Case (`docs/reviews/`), den diese ADR an einer
Wegwerf-PostgreSQL mit den realen Adaptern reproduziert hat)

**Bezug:** [`LH-FA-SCH-005`](../../../spec/lastenheft.md) (Schema-Version;
Boundary: Versionen vor und nach einer Schemaänderung unterscheidbar),
[`LH-FA-SCH-002`](../../../spec/lastenheft.md) (hinzugefügte Spalten),
[`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Backfill des Bestands),
[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) (teilweise
superseded — Haupt-Bezug), [ADR-0015](0015-schema-evolution.md) (Schema
Evolution, dynamische Re-Versionierung),
[ADR-0029](0029-domain-invarianten.md) (Regel 7: jeder Change referenziert eine
Schema-Version), [ADR-0115](0115-backfill-spaltenwerte-text-ergebnisformat.md)
(nur zur Abgrenzung: Werte des Bildes)

**Schärft:** [`LH-FA-SCH-004.a`](../../../spec/pflichtenheft.md) (Absatz
„jeder Change referenziert eine Schema-Version" — die Aussage bleibt wörtlich
wahr; diese ADR legt fest, welche Version eine Backfill-Change referenziert und
was die Referenz nicht zusagt), [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md)
(Absatz „Markierung": ein Satz zur Version, siehe Folgepflicht 3),
[`SPEC-004`](../../../spec/pflichtenheft.md) (TableSchema/SchemaVersion)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) ist `Accepted` und
unberührbar (`AGENTS.md` §3.5). Sein akzeptiertes Negativ „Schema-Version-Verweis"
sagt: Backfill-Changes referenzieren die zum Run-Start aktuelle Version; eine
kompatible Spalten-Erweiterung **während des Runs** erscheine im Image ohne
Versions-Bump — harmlos, das Image sei ein selbstbeschreibendes JSON-Objekt. Die
Fälle **vor** dem Run nennt der Punkt nicht; die Ergänzung ändert den Referenten
und ist keine Zitat-Korrektur.

### Befunde am Bestand (Code-Anker als Symbolnamen, Stand `77fc6e43`)

- **Der Run liest die Version vor dem Snapshot.** `BackfillTableService.Execute`
  ruft `currentVersion` (`internal/application/usecase/backfill/service.go`) →
  `SchemaStorePort.CurrentVersion`, das die höchste Zeile aus
  `cdc.schema_version` der Tabelle liefert
  (`queries.SelectCurrentSchemaVersion`); jede Backfill-Change trägt sie, weil
  `model.NewChange` eine Schema-Version verlangt.
- **Die Version wechselt allein im Erfassungspfad.** Schreiber der Tabelle
  `cdc.schema_version` sind die Erstaktivierung (`TableActivationAdapter.Register`,
  nur die Versions-Zeile, **keine** Spaltenform) und
  `Assembler.observeRelation` (`internal/adapters/driving/replication/mapper/mapper.go`)
  über `SchemaStorePort.RegisterVersion`: bei der **ersten** Relation-Nachricht
  trägt es die Spaltenform der Version 1 nach, ohne zu wechseln; bei einer
  kompatiblen Erweiterung registriert es die nächste Version. Beides geschieht
  erst mit einer Relation-Nachricht, die pgoutput vor der ersten Change nach
  einer DDL-Änderung sendet — ein `ALTER TABLE` allein schreibt nichts nach
  `cdc.*`.
- **Der Snapshot trägt keine Typ-Kennungen.** `TableSnapshot.Columns()`
  (`internal/application/port/outbound/tablesnapshot.go`) liefert Spaltennamen
  ohne OIDs; `cdc.table_schema` führt `column_oid` je Spalte.
- **Nur der WAL-Pfad hält eine Versions-Kennung im Speicher.** Die
  `TableBinding` des `Assembler` trägt die Kennung der Version, mit der er
  Changes stempelt; `setSchemaVersion` wird ausschließlich in `observeRelation`
  aufgerufen, und nur im Zweig der kompatiblen Erweiterung (die Zweige „unverändert"
  und „nicht Erweiterung" kehren vorher zurück).
- **Rechte.** `cdc_capture` trägt `INSERT` auf `cdc.schema_version` und
  `cdc.table_schema` (`tools/schema/nacharbeit-roles.sql`), die Worker-Verbindung
  könnte also Versionen schreiben.

### Gemessen (Wegwerf-Läufe, nicht committet)

**Probe** (Wegwerf-Programm im Scratchpad, 2026-09-24; Wegwerf-Container mit
`wal_level=logical`, Container und Netz danach entfernt): PostgreSQL **17.11**
(`postgres:17-alpine`, Digest der E2E-Matrix) und **18.6** (`postgres:18-alpine`,
Digest von `PG_TEST_IMAGE`); Schema-Rollout des Arbeitsbaums über
`tools/schema/apply-rollout.sh`; die **realen Adapter** (Aktivierung, Schema
Store, Snapshot-Leser, Annahme, Run-Zustand, Schreiber) unter dem Use Case des
Runs, Superuser-DSN. Zwei Tabellen mit je zwei Zeilen (`id int, a text`); nach
der Aktivierung je `ALTER TABLE … ADD COLUMN b text DEFAULT 'neu'`, **keine**
Change über den WAL-Pfad, dann Antrag und Run. Der Nachtrag der Spaltenform der
Tabelle 2 ist mit `RegisterVersion` **nachgestellt**, wie `observeRelation` es
für die erste Relation-Nachricht ausführt (Code-Anker oben); der Replication-Strom
selbst lief in der Probe nicht.

| Fall | `CurrentVersion` vor dem Run | `TableSchema` dieser Version | Backfill-Changes (Lauf `completed`, 2 Zeilen) |
|---|---|---|---|
| 1 — Aktivierung, `ADD COLUMN`, keine WAL-Change | `sv-p1-1` (v1) | keine Zeilen in `cdc.table_schema` | `origin = backfill`, `schema_version = sv-p1-1`, Bild `{"a","b","id"}` |
| 2 — Spaltenform `[id, a]` nachgetragen, `ADD COLUMN`, keine WAL-Change | `sv-p2-1` (v1) | `[id, a]` | `origin = backfill`, `schema_version = sv-p2-1`, Bild `{"a","b","id"}` — `b` steht im Bild, nicht in der Spaltenform der Version |

Identisches Ergebnis auf beiden Versionen (17.11 und 18.6); nach dem Run steht in
beiden Fällen weiterhin nur Version 1 in `cdc.schema_version`. Der Wert aus F-8 ist
damit an den realen Adaptern belegt: die Referenz kann hinter den Spalten des
Bildes liegen (Fall 2) oder auf eine Version ohne Spaltenform zeigen (Fall 1).

### Konstraints

- [`LH-FA-SCH-005`](../../../spec/lastenheft.md) Boundary: zwei Changes vor und
  nach einer Schemaänderung sind an ihrer Version unterscheidbar — erfüllt für jede
  Änderung, die der Erfassungspfad beobachtet hat: im Fall 2 registriert die erste
  danach erfasste WAL-Change die nächste Version; im Fall 1 hat der Erfassungspfad
  vor der ersten Relation-Nachricht keine Änderung beobachtet, die Version 1 ist
  dort die Form ab dieser Nachricht.
- Kein Consumer-Weg liefert `cdc.table_schema`: `cdc_reader` trägt kein Recht darauf
  (`tools/schema/nacharbeit-roles.sql`), und keine Nachrichtenform der Zustellwege
  enthält eine Spaltenform; der einzige Leser der Spaltenform ist
  `observeRelation` als Vergleichsbasis der nächsten Relation-Nachricht
  (Suchlauf `grep -rn 'TableSchema(' internal cmd`, ohne Testdateien). Die
  Schema-Version ist für Consumer eine **Kennung**, kein Schlüssel zu einer
  Spaltenform.
- Der Capture-Pfad bleibt unberührt (`ADR-0111` Teilfrage 3/5: der Run
  berührt weder Position noch Heartbeat noch Bindung).

## Entscheidung

Wir wählen **die Version zum Run-Start als Marker — ohne Katalog-Abgleich, ohne
neue Version, ohne Abbruch**. Die Referenz ist eine Kennung der Reihenfolge, keine
Beschreibung der Bild-Spalten; ihre Grenze wird benannt statt beseitigt. Vier
Festlegungen.

**Was diese ADR von `ADR-0111` ersetzt** — genau einen Aufzählungspunkt: in
§Konsequenzen, „Akzeptierte Negative", den Punkt „Schema-Version-Verweis" (gilt
jetzt: Festlegung 1 und 2). Alles Übrige — Teilfragen 1–8, Festlegungen, die
übrigen akzeptierten Negative, Fitness Function — bleibt gültig.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (Wortlaut von `ADR-0111`) | keine Änderung | der Wortlaut nennt „während des Runs"; die zwei häufigeren Fälle **vor** dem Run (Fälle 1 und 2 der Probe) stehen nirgends, das Argument „selbstbeschreibendes Bild" trägt sie, ohne dass jemand sie so gelesen hätte; ein Leser der Grenze findet sie erst am Befund |
| B — beim Snapshot die Version aus dem Katalog bestimmen und bei Abweichung eine neue Version registrieren | die Referenz träfe die Bild-Spalten | (1) `Columns()` trägt Namen ohne OIDs, `cdc.table_schema` verlangt `column_oid`: der Vergleich sähe eine Typänderung nicht, eine neue Spaltenform ließe sich nicht vollständig schreiben; (2) der Worker schreibt eine Version, die der `Assembler` nicht kennt: seine `TableBinding` behielte die ältere Kennung, weil `setSchemaVersion` im Zweig „unverändert" nicht läuft — spätere WAL-Changes trügen eine **ältere** Version als die Backfill-Changes, der Versionsverlauf liefe rückwärts (aus dem Code hergeleitet, nicht gemessen); die Behebung ändert den Capture-Pfad; (3) zwei Schreiber derselben Versionsnummer (Worker-Pool, Erfassungspfad) ohne gemeinsame Sperre; (4) Fall 1 bräche die Zusage, dass die erste Relation-Nachricht die Spaltenform der Version 1 trägt |
| C — Run mit Klasse `configuration` abbrechen, wenn die Version nicht zu den Snapshot-Spalten passt | keine falsche Referenz | im Fall 1 (kein `TableSchema`) ist der Regelfall „frisch aktiviert, sofort Backfill" — ein Abbruch machte die Kernfähigkeit unbenutzbar oder verlangte eine Ausnahme, die den Fall aus der Prüfung nimmt; im Fall 2 ist die Abhilfe eine Schreibänderung an der Tabelle, damit der Erfassungspfad eine Version anlegt — bei einer ruhenden Tabelle eine künstliche Handlung, für eine Kennung, die kein Consumer nachschlägt |
| D — Version **nach** dem Snapshot lesen | kleineres Fenster gegen eine parallel eintreffende Relation-Nachricht | ändert an den Fällen 1 und 2 nichts (kein Wechsel ohne WAL-Change); ein Fenster gegen den Erfassungspfad bleibt in jeder Reihenfolge |
| **E — Marker: Version zum Run-Start, Reichweite benannt (gewählt)** | keine Änderung an Run, Capture-Pfad oder Schema; die Boundary von `LH-FA-SCH-005` bleibt erfüllt; die Grenze steht an den Stellen, an denen Leser sie suchen | die Referenz einer Backfill-Change ist kein Beleg für die Bild-Spalten; ein Leser, der die Spaltenform aus der Version ableitet, irrt in Fall 2 |

### Festlegung 1 — Bedeutung der Referenz

Die `schema_version` einer Backfill-Change ist die **höchste** in
`cdc.schema_version` registrierte Version der Tabelle **zum Zeitpunkt der
Vorbedingungsprüfung des Runs** (Aufruf von `CurrentVersion` in `Execute`, vor dem
Öffnen des Snapshots). Alle Changes des Runs tragen dieselbe Kennung. Sie ist
eine **Kennung der Reihenfolge**: sie unterscheidet die Backfill-Changes von
Changes einer später registrierten Version (`LH-FA-SCH-005` Boundary), sie sagt
**nicht**, welche Spalten das Bild trägt. Die Spaltenform der Backfill-Change
steht im Bild selbst (JSON-Schlüssel der Spalten des Snapshots); wer die Spalten
einer Change braucht, liest das Bild.

### Festlegung 2 — Reichweite (was das akzeptierte Negativ jetzt deckt)

Die Referenz kann von den Spalten des Bildes abweichen in drei Fällen; in keinem
ändert der Run die Version, bricht ab oder meldet:

1. **Erweiterung während des Runs** — wie in `ADR-0111`.
2. **Kompatible Erweiterung vor dem Run**, nach der letzten Relation-Nachricht der
   Tabelle und ohne WAL-Change der Tabelle dazwischen: das Bild trägt die neue
   Spalte, die Referenz die Version davor. Die erste danach erfasste WAL-Change
   registriert die nächste Version (Fall 2 der Probe).
3. **Version ohne Spaltenform** (Erstaktivierung ohne WAL-Change seither): die
   Referenz zeigt auf eine Version ohne Zeile in `cdc.table_schema`; die erste
   Relation-Nachricht trägt die Spaltenform nach und wechselt die Version nicht
   (Fall 1 der Probe).

**Nicht Gegenstand des Runs:** die Erkennung inkompatibler Änderungen. Eine
entfernte Spalte oder Typänderung endet im WAL-Pfad mit der Klasse `schema`
(`LH-FA-SCH-004`); ein Run vergleicht Version und Snapshot-Spalten nicht und
trägt den Stand der Tabelle zum Snapshot in das Bild.

### Festlegung 3 — Abhilfe (optional, erwartet)

Wer für einen Run eine Referenz will, deren Spaltenform die aktuellen Spalten
trägt, löst **vor dem Antrag** eine Änderung an einer Zeile der Tabelle aus: die
erste WAL-Change nach dem `ALTER TABLE` bringt die Relation-Nachricht, und
`observeRelation` registriert die Version. **Erwartet, nicht gemessen** (der
Replication-Strom lief in der Probe nicht); keine Zusage dieser ADR, keine
Vorbedingung des Backfills.

### Festlegung 4 — Wo die Grenze steht

Der Wortlaut „Referenz ist eine Kennung, kein Beleg für die Bild-Spalten" steht an
drei Stellen, an denen ein Leser die Version deutet: (a) als Grenze-Kommentar an
der Stelle im Use Case, die die Version liest (`currentVersion`); (b) im
Handbuch, im Abschnitt zum Backfill; (c) in `LH-FA-CAP-009.a` (Absatz
„Markierung") als ein Satz. Der E2E-Beleg des Backfills prüft an der Version
**genau das**: die Kennung der Backfill-Changes ist die zum Antrag aktuelle
Zeile aus `cdc.schema_version`, nicht ihre Übereinstimmung mit den Bild-Spalten.

## Konsequenzen

- Positiv: keine Änderung an Run, Capture-Pfad, Rollen oder Schema; die drei Fälle
  und ihre Reichweite sind an einer Stelle benannt und belegt (Fälle 1 und 2 real
  gemessen).
- Positiv: die Kernnutzung „aktivieren, dann sofort Backfill" bleibt ohne
  Vorbedingung an den Erfassungspfad.
- Negativ: ein Consumer, der aus der Version eine Spaltenform ableiten wollte,
  bekäme für Backfill-Changes im Fall 2 eine Form ohne die neue Spalte; der
  Weg dafür existiert für Consumer nicht (kein Zustellweg liefert
  `cdc.table_schema`, siehe Konstraints).
- **Akzeptierte Negative** (kurz begründet, keine Folgepflicht über die
  Folgepflichten unten hinaus):
  - *Backfill-Change trägt die Version „vor" einer Erweiterung, obwohl ihr Bild
    danach entstand:* die Reihenfolge der Kennungen bleibt monoton (die
    nächste WAL-Change trägt eine höhere Version), nur ihre Deutung „Bild passt
    zur Form" gilt für Backfill-Changes nicht.
  - *Version ohne Spaltenform im Fall 3:* der Run darf vor der ersten WAL-Change
    der Tabelle laufen (die Kernnutzung „aktivieren, dann Backfill"); die Referenz
    zeigt dann auf eine Version, deren Spaltenform erst die erste Relation-Nachricht
    nachträgt.

### Folgepflichten

Jede Pflicht hat einen Träger im Architect-Verdikt zu dieser ADR
(Verzeichnis `docs/reviews/`, Verdikt „backfill-schema-klasse-rollen“); **die
genannten Träger ändert diese ADR nicht** — der Planner zieht sie nach
(`AGENTS.md` §3.13).

1. **Grenze-Kommentar** an `currentVersion` im Use Case des Runs (kein
   Verhaltens-Diff).
2. **Handbuch:** ein Satz zur Referenz im Abschnitt zum Backfill.
3. **Spec:** ein Satz in `LH-FA-CAP-009.a`, Absatz „Markierung" („die
   Schema-Version einer Backfill-Change ist die zum Run-Start aktuelle Version der
   Tabelle; sie unterscheidet, sie beschreibt die Bild-Spalten nicht").
4. **E2E-Beleg:** die Prüfung der Version wie in Festlegung 4.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test, reale PostgreSQL (`make test-integration`) | die `schema_version` aller Backfill-Changes eines Runs ist die zum Antrag aktuelle Zeile aus `cdc.schema_version` der Tabelle; nach einer Erweiterung ohne WAL-Change bleibt sie unverändert und das Bild trägt die neue Spalte | `make test-integration` |
| Review-Prüfpflicht | kein Text im Repo sagt, die Referenz einer Backfill-Change beschreibe die Spalten des Bildes | — (kein Gate) |

## Re-Evaluierungs-Trigger

- **Ein Zustellweg oder ein Lese-Vertrag liefert die Spaltenform je Version**
  (`cdc.table_schema` oder ein Ableger wird Consumer-sichtbar): die Grenze
  wird zu einer Lücke — Folge-ADR (Option B mit den Nachteilen (2) und (3)
  aufgelöst, oder ein Kennzeichen an der Referenz).
- **Ein Betrieb meldet einen Consumer, der die Referenz als Spaltenform
  liest:** neu bewerten.
- **Der Erfassungspfad registriert Versionen anders** (z. B. bei der Aktivierung
  mit Spaltenform): Fall 3 entfällt, Fall 2 neu bewerten.
- Sonst permanent — die Bedeutung „Kennung der Reihenfolge" gilt, solange die
  Spaltenform Consumern nicht geliefert wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-24 | Proposed — Architect-Vorschlag zum Review-Befund F-8; die zwei Fälle vor dem Run an Wegwerf-Containern auf PostgreSQL 17.11 und 18.6 mit den realen Adaptern gemessen | Review-Report zum Run-Use-Case, Befund F-8 (`docs/reviews/`) |
| 2026-09-24 | Accepted — Annahme durch den Auftraggeber samt Reihenfolge-Marker-Bedeutung der Version und den drei benannten Reichweiten | [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0116` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
