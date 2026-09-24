# ADR-0115: Backfill — Spaltenwerte im Text-Ergebnisformat der Ausgabefunktion statt `col::text` (Supersedes ADR-0111, Teilfrage 1 und 2 teilweise)

**Status:** Proposed — Supersedes [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
in **Teilfrage 1 (Ablauf Nr. 3) und Teilfrage 2, jeweils teilweise**, sowie in
einer Zeile der §Fitness Function (genau drei Stellen, siehe §Entscheidung);
alles Übrige von `ADR-0111` bleibt in Kraft

**Datum:** 2026-09-24

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch den Befund F-1 des
Review-Reports zum Snapshot-Leser (`docs/reviews/`),
den diese ADR an einer Wegwerf-PostgreSQL selbst reproduziert und an einem breiten
Typ-Satz gegen mehrere Lesevarianten gemessen hat)

**Bezug:** [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Backfill des
Bestands), [`LH-FA-CAP-008`](../../../spec/lastenheft.md) (Row-Image-Form und
Abwesenheit), [`LH-QA-SEC-004`](../../../spec/lastenheft.md) (Spaltenausschluss,
unberührt), [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
(teilweise superseded — Haupt-Bezug), [ADR-0016](0016-jsonb-row-images-mvp.md)
(Row-Image-Form: Text-Stand der Quelle ohne Typ-Interpretation),
[ADR-0059](0059-spaltenauswahl-mechanismus.md) (Ausschluss, unberührt),
[ADR-0030](0030-testpyramide.md) (Testpyramide, DB-Tier)

**Schärft:** [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) (Absatz
„Markierung": „Das Row Image ist byte-gleich dem WAL-Image derselben Zeile …
Werte im Text-Stand der Quelle" — die Aussage bleibt wörtlich wahr; diese ADR
legt fest, woraus der Text-Stand entsteht und wie die Byte-Gleichheit belegt
wird), [`SPEC-012`](../../../spec/pflichtenheft.md) (PostgreSQL 17 und 18: die
Parität ist auf beiden gemessen)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) ist `Accepted` und
unberührbar (`AGENTS.md` §3.5). Seine Teilfrage 2 sagt, ein Backfill-Image
dürfe sich von einem WAL-Image derselben Zeile nicht unterscheiden; Teilfrage 1
(Ablauf Nr. 3) und Teilfrage 2 setzen die Mittel dafür: die Spaltenwerte
werden per `col::text` gelesen, „in derselben Sitzungs-GUC-Lage wie der
Walsender". Das Mittel trägt die Zusage nicht. Die Korrektur ändert den
Referenten, ist also keine Zitat-Korrektur (`AGENTS.md` §3.5) und braucht diese
Folge-ADR.

### Befund (reproduziert, gemessen)

Der Cast `col::text` ruft für Typen mit eigener Cast-Funktion in `pg_cast` diese
statt der Ausgabefunktion (`typoutput`) auf; `pgoutput` sendet im Text-Modus je
Spalte das Ergebnis der Ausgabefunktion. Beide Wege liefern deshalb nicht für
jeden Typ denselben Text (die Erklärung ist aus dem Review-Befund übernommen und
nicht am Katalog nachgelesen; belegt ist das Ergebnis durch die Messung unten).

**Probe** (Wegwerf-Programm im Scratchpad, nicht committet; der Umsetzungs-Slice
wiederholt die Zeilen als Test; 2026-09-24, Wegwerf-Container, `wal_level=logical`; Container und Netz
danach entfernt): PostgreSQL **17.11** (`postgres:17-alpine`, Digest der
E2E-Matrix) und **18.6** (`postgres:18-alpine`, Digest von `PG_TEST_IMAGE`);
eine Tabelle mit **86** lesbaren Spalten (`id` und 85 Typ-Spalten; dazu eine
gelöschte, eine `STORED`-generierte und unter 18 eine `VIRTUAL`-generierte
Spalte, die im WAL und im Backfill fehlen) mit je einer Wert-Zeile, einer
Rand-Zeile (leere Zeichenketten, `infinity`, `NaN`, leere Arrays/Ranges,
Unicode) und einer `NULL`-Zeile. **WAL-Bild** = `pglogrepl`-Insert-Tuple aus
einem temporären Slot mit Publication (`proto_version '1'`, wie der
Replication-Pfad) unter einer Rolle mit `LOGIN REPLICATION`; **Backfill-Bild** =
Lesen derselben Tabelle mit derselben Rolle und demselben DSN in einer
`REPEATABLE READ`-Transaktion über einen `NO SCROLL`-Cursor mit `FETCH FORWARD`,
je Variante. Eine Spalte trifft, wenn alle drei Zeilen byte-gleich (auch in
NULL/leer) zum WAL-Text sind. Drei Rollen-GUC-Lagen je Version (Standard;
`timezone='Asia/Tokyo'` + `datestyle='German, DMY'` +
`intervalstyle='sql_standard'` + `bytea_output='escape'` +
`extra_float_digits=0`; `intervalstyle='postgres_verbose'` +
`datestyle='SQL, MDY'` + `extra_float_digits=-3` +
`timezone='America/New_York'`), dazu die Lage mit `lc_monetary='de_DE.utf8'` auf
`postgres:18` (Debian, 18.6, **ungepinnt**, nur für diese eine Lage). Der
Typ-Satz: Boolean; `int2/4/8`, `real`, `float8`, `numeric`; `money`;
`char(n)`, `varchar(n)`, `text`, `"char"`, `name`; `bytea`; `date`, `time`,
`timetz`, `timestamp`, `timestamptz`, `interval`; `bit`, `varbit`; `uuid`;
`inet`, `cidr`, `macaddr`, `macaddr8`; `json`, `jsonb`, `jsonpath`, `xml`;
`oid`, `regclass`, `regtype`, `regproc`, `xid`, `tid`, `pg_lsn`; `tsvector`,
`tsquery`; sieben geometrische Typen; `int4range`, `tstzrange`, `daterange`,
`int4multirange`, ein eigener Range-Typ; Enum; vier Domains (über `int`,
`text`, `int[]`, `bool`); zwei Composites (eines mit `bool`/`inet`/`char(4)`);
Extension-Typen `hstore`, `citext`, `ltree`; Arrays davon (`int4[]`, `text[]`,
`bool[]`, `inet[]`, `char(3)[]`, `bytea[]`, `timestamptz[]`, `float8[]`,
zweidimensional und mit Untergrenze, `money[]`, Enum-, Composite-, `interval`-,
`numeric`-, `uuid`-, `cidr`-, `bit`-, Range- und `hstore`-Array).

| Variante | Treffer (Spalten von 86) | Läufe |
|---|---|---|
| **B** — `col::text` (Wortlaut von `ADR-0111`) | **79** | 17.11 und 18.6 je drei GUC-Lagen, Debian 18.6 mit `de_DE.utf8`: in allen Lagen 79 |
| **A** — Cursor, Text-Ergebnisformat, **kein Cast** | **86** | dieselben sieben Lagen: in allen 86 |
| A2 — erweitertes Protokoll, Text-Format je Spalte angefordert, kein Cast | 86 | 17.11/18.6 je drei Lagen, Debian |
| A3 — `pgx.Query` mit `QueryResultFormats` Text und `RawValues` | 86 | dieselben |
| A4 — `pgx.Query` **ohne** Format-Vorgabe (Treiber-Default) und `RawValues` | 29 | dieselben — der Default liefert Binärbytes |
| C — Aufruf der Ausgabefunktion je Spalte (`pg_type.typoutput`) | 86 | dieselben |
| D — `COPY (SELECT …) TO STDOUT` (Text) mit Rück-Escaping | 86 | dieselben |
| F1 — `to_jsonb(col) #>> '{}'` | 55 | nur 18.6, Lage Tokyo |
| F2 — `CASE … format('%s', col)` mit NULL-Wache | 86 | nur 18.6, Lage Tokyo |

Die sieben Spalten, die B verfehlt (Lauf 18.6, Standard-GUC): `bool` (WAL `t`,
Cast `true`), `char(5)` (WAL `ab   `, Cast `ab`), `char(1)` (WAL ` `, Cast
leer), `inet` (WAL `1.2.3.4`, Cast `1.2.3.4/32`; ebenso `::1` und `2001:db8::1`
mit `/128`), `xml` (WAL `<b/>`, Cast `<?xml version="1.0"?><b/>` bei einem Dokument
mit Deklaration) und die Domain über `bool`. Das sind **fünf Typen**; der Review
hatte an 32 Typen drei gefunden. `cidr`, `macaddr`, `money`, Enum, Composite,
Ranges und Arrays stimmen mit dem Cast überein.

Die Rollen-GUC prägen das WAL-Bild real (gedruckt, 18.6 Lage Tokyo:
`23.09.2026 17:11:12.5 JST`, `+1-2 +3 +4:05:06.7`, `\336\255\276\357`; Lage `postgres_verbose`:
`09/23/2026 04:11:12.5 EDT`, `@ 1 year 2 mons 3 days 4 hours 5 mins 6.7 secs`;
`money` unter `de_DE.utf8`: `1.234.567,89 €`), und A folgt ihnen in jeder Lage.

**Kosten** (Lauf 2026-09-24, PostgreSQL 18.6 Alpine, Loopback im Docker-Netz,
Go-Client mit `pgconn`, `NO SCROLL`-Cursor in Blöcken von 1.000): 1.000.000 Zeilen
mit zehn Spalten (`int`, `timestamptz`, `float8`, `numeric`, `text`, `jsonb`,
`bool`, `inet`, `char(8)`, `text[]`), rund 191 MB Nutzdaten, je drei Läufe: A
1,72 / 1,70 / 1,79 s, B 2,44 / 2,48 / 2,49 s, D 1,63 / 1,52 / 1,40 s (D in einem
Puffer im Speicher, ohne Rück-Escaping). Gemessen ist der Lese-Lauf des Clients,
nicht der ganze Backfill; abgeleitet aus den Mittelwerten: A benötigt rund 29 %
weniger Zeit als B, D rund 13 % weniger als A.

Die Fixture der Paritätsprobe im Umsetzungs-Slice trägt acht Typen, keinen
davon aus der Fehlermenge; sie belegt die Zusage nur an ihrer Eingabeseite, nicht
am Typ-Satz.

### Konstraints

- `model.BuildRowImage(columns []string, values []*string, excluded []string)`
  ist die eine Konstruktionsstelle für Row Images (`ADR-0111` Teilfrage 2);
  die Lösung setzt an den **Werten** an und ändert die Domänenfunktion nicht.
- Ausschluss-Zusage (`LH-QA-SEC-004`) und Abwesenheits-Regel (`NULL` entfällt,
  `LH-FA-CAP-008`) bleiben.
- Der Adapter liest über `pgconn` im einfachen Query-Pfad (Text-Ergebnisformat
  ist dort Protokoll-Eigenschaft); die Cursor-Struktur (Blöcke, `NO SCROLL`,
  Snapshot-Import) bleibt.

## Entscheidung

Wir wählen **Option A: Der Backfill liest jede Spalte unverändert — ohne Cast und
ohne Funktion auf dem Spaltenwert — im Text-Ergebnisformat des Wire-Protokolls
und übernimmt die Rohbytes des Feldes als Wert** (`nil` für NULL, leeres Feld
für leeren Text). Das Ergebnis ist die Ausgabefunktion des Spaltentyps unter den
GUC der Sitzung — dieselbe Erzeugung wie der Text, den `pgoutput` je Spalte im
WAL-Pfad sendet. Vier Festlegungen.

**Was diese ADR von `ADR-0111` ersetzt** — genau drei Stellen: (a) in Teilfrage 1
im Ablauf Nr. 3 der Satzteil „mit `col::text` je Spalte" (gilt: Festlegung 1);
(b) in Teilfrage 2 im Absatz „Row-Image-Parität" der Satzteil „`col::text` in
derselben Sitzungs-GUC-Lage wie der Walsender" (gilt: Festlegung 1 und 2);
(c) in §Fitness Function die Typliste der Zeile „Bild-Parität" (`timestamptz`,
`float8`, `bytea`, `interval`, `numeric`, `jsonb`, Array, `date`, generierte
Spalte) (gilt: die Fitness Function dieser ADR). Alles Übrige — die
Spaltenliste, der `NO SCROLL`-Cursor in Blöcken, die Speichergrenze durch `B`,
die eine gemeinsame Bild-Funktion, generierte Spalten ausgenommen, `NULL`
entfällt, Ausschluss angewandt, `M1`–`M6` — bleibt gültig und wird hier nicht
wiederholt. `M5` (der Walsender trägt die Rollen-GUC) bleibt wahr; die
Sitzungs-GUC-Gleichheit einer regulären Verbindung, die `ADR-0111` als erwartet
führte, ist mit dieser Probe auf beiden Versionen belegt.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| B — `col::text` beibehalten (nichts tun) | einfachste Anweisung; `ADR-0111` wörtlich | 79 von 86 Spalten; verletzt die Kernzusage für `bool`, `char(n)`, `inet`, `xml` (und Domains darüber); Ausnahmen wären Typ-Wissen im Adapter (siehe B′); rund 40 % länger als A (2,44–2,49 s gegen 1,70–1,79 s) |
| B′ — `col::text` mit Ausnahmen je Typ | trifft an diesem Typ-Satz vollständig, wenn die Ausnahmen die fünf Typen nennen | Wartungsliste, die mit jedem Typ, jeder Extension, jeder Domain (der Cast hängt am Domain-Typ, nicht am Basistyp) und jeder PostgreSQL-Hauptversion falsch werden kann; die vollständige Form, die Ausnahmen aus `pg_cast` ableitet, läuft auf C hinaus |
| **A — Text-Ergebnisformat, kein Cast (gewählt)** | 86/86 in allen sieben gemessenen Lagen; **kein Typ-Wissen**: jeder heutige und künftige Typ, jede Domain, jedes Array und jede Extension liefert per Konstruktion die Ausgabefunktion; `SELECT "col"` ohne Aufruf ist die billigste Form (schneller als B); ändert nur die Cursor-Anweisung, `BuildRowImage` und die Wertbildung `[]*string` bleiben | die Zusage hängt an der Formatwahl des Lesepfads: ein Treiber-Default kann Binär liefern (A4: 29/86) — die Formatwahl wird Teil der Zusage und getestet (Festlegung 2) |
| C — Ausgabefunktion je Spalte im SQL aufrufen | 86/86 | Katalog-Nachschlagen je Spalte (`pg_type.typoutput`, Schema-Qualifizierung, polymorphe und Pseudo-Typ-Signaturen, `cstring` als Ergebnistyp, `EXECUTE`-Recht); gleiches Ergebnis wie A auf einem Umweg; Kosten nicht gemessen |
| D — `COPY (SELECT …) TO STDOUT` (Text) | 86/86; am schnellsten (rund 10 % vor A); gemessen mit demselben Ergebnis | eigener Zeilen-/Feld-Parser mit Rück-Escaping (`\N`, `\\`, Tab, Zeilenende) als zweiter Ort, der Text interpretiert; Strom statt `FETCH`-Blöcken (die Blockgrenze `B` und ihre Speichergrenze werden vom Adapter neu gebaut); der Gewinn rechtfertigt das nicht, solange die Ausbaustufe für Durchsatz (`ADR-0111` Re-Evaluierung) nicht ansteht |
| E — Bild über den Logical-Decoding-Pfad selbst ableiten | wäre der WAL-Weg selbst | nicht ausführbar: das WAL trägt die Bestandszeilen nicht; sie zu erzeugen hieße, in die Quelle zu schreiben (`ADR-0111` Teilfrage 1 Option D, verworfen) |
| F — JSON-/`format`-Funktionen (`to_jsonb(col) #>> '{}'`, `format('%s', col)`) | F2 86/86 | F1 nur 55/86 (JSON-Konventionen für Zeit, Zahl, Boolean); F2 braucht je Spalte eine NULL-Wache und ist ein Umweg auf die Ausgabefunktion ohne Vorteil gegenüber A (nur 18.6 Lage Tokyo gemessen) |

### Festlegung 1 — Lese-Anweisung ohne Umwandlung

Die Cursor-Anweisung des Snapshot-Lesers wählt jede Spalte nur beim Namen
(`quote_ident`-gequotet, in der Reihenfolge der Spaltenliste) — kein `::`-Cast,
kein Funktionsaufruf, kein `CASE`, kein Alias mit Umwandlung. Der Adapter gibt
das Feld als Text weiter: `nil` ↔ NULL, sonst die Bytes des Feldes als Zeichenkette
(auch leer). Ein Typ-Katalog im Produktionscode entfällt: es gibt keine
Typ-Liste, die mit neuen Typen nachzuziehen ist.

### Festlegung 2 — Das Text-Ergebnisformat ist Teil der Zusage

Der Lesepfad fordert das Text-Ergebnisformat **ausdrücklich** oder nutzt einen
Protokollweg, in dem Text die Definition ist (der einfache Query-Pfad von
`pgconn`, wie ihn der Snapshot-Leser heute nutzt). Ein Wechsel auf einen
Treiberweg mit Binär-Vorgabe (`pgx.Query` liefert Binärbytes, A4) setzt das
Text-Format je Spalte (`QueryResultFormats`, A3: 86/86); der Paritätstest
(Fitness Function) fängt den Wechsel als Rot. Der Adapter setzt keine
Sitzungs-GUC (`SET`, `options`-Parameter): die Gleichheit mit dem Walsender
kommt aus demselben DSN und derselben Rolle (`ADR-0111` Teilfrage 1: alle
Verbindungen über `CDC_CAPTURE_DSN`).

### Festlegung 3 — Row Image, TOAST, Abwesenheit

`BuildRowImage(columns, values, excluded)` bleibt unverändert und bekommt als
`values` die nach Festlegung 1 gelesenen Zeichenketten. Der Backfill liest die
Spalten vollständig (die Zeile wird von der Abfrage entpackt); der
`UPDATE`-Abwesenheitsfall unveränderter TOAST-Spalten kommt im Backfill nicht
vor (`ADR-0111` Replay-Invariante, Ausnahme dort benannt). Ein `NULL` entfällt
im Bild, ein leerer Text steht als leerer String im Bild — beides gleich im
WAL-Pfad und gemessen (Rand-Zeile und `NULL`-Zeile).

### Festlegung 4 — Der Beleg der Byte-Gleichheit ist der Typ-Satz

Die Zusage „byte-gleich" gilt am **Typ-Satz**, nicht an einer Auswahl. Der
Paritätstest der Fitness Function trägt den Typ-Satz aus §Kontext als Daten
(Tabelle aus Spaltenname, Typ-DDL, zwei Wert-Ausdrücken, dazu eine
`NULL`-Zeile) und vergleicht den Roh-Text je Spalte zwischen WAL-Pfad und
Backfill-Pfad. Ein neuer Typ, der Anlass zu Zweifeln gibt, ist eine Zeile in
der Tabelle, keine Codeänderung im Adapter.

## Konsequenzen

- Positiv: die Kernzusage von `ADR-0111` Teilfrage 2 ist am Typ-Satz belegt
  (86 von 86 Spalten, PostgreSQL 17.11 und 18.6, drei Rollen-GUC-Lagen und
  `lc_monetary`), nicht mehr an acht Typen.
- Positiv: keine Typ-Wartung im Adapter — jeder künftige Typ, jede
  Domain, jede Extension und jede PostgreSQL-Hauptversion trägt die Parität
  per Konstruktion (hergeleitet; belegt an den 86 Spalten), weil derselbe
  Baustein (die Ausgabefunktion) beide Pfade speist.
- Positiv: die Lese-Anweisung wird einfacher und rund 29 % kürzer im Lese-Lauf als
  der Cast (abgeleitet aus der Messung an 1.000.000 Zeilen).
- Negativ: die Parität hängt am Text-Format und an der Gleichheit der
  Sitzungs-GUC-Lage; beide sind Vertragsbestandteil (Festlegung 2) und durch den
  Paritätstest gebunden, nicht durch die Anweisung allein.
- Negativ: die Bindung an den Lesepfad — ein späterer Wechsel auf `COPY` (D,
  Durchsatz-Ausbaustufe) ist möglich (86/86 gemessen), bringt aber den
  Rück-Escaping-Parser als zweiten Text-Interpreten.
- Akzeptiertes Negativ (kurz begründet, keine Folgepflicht): eine Vollständigkeits-
  Prüfung des Typ-Satzes gegen `pg_type` (jeder Basistyp braucht eine Zeile)
  entfällt — die Zusage ist strukturell (kein Cast), nicht typ-weise; ein
  Vollständigkeits-Test bräuchte einen gültigen Wert je Basistyp der Version
  und würde eine Pflege erzeugen, die die Konstruktion gerade vermeidet.
- **Folgepflichten** (Umsetzung; Träger ist der Slice, der `ADR-0111`
  Teilfrage 1 umsetzt — die Fixrunde zum Review-Befund F-1 — und die Planung):
  1. **Adapter:** die Cursor-Anweisung des Snapshot-Lesers wählt Spalten ohne
     Cast (Festlegung 1); der Kommentar an der Anweisung und die Beschreibung
     der Zeilenwerte im Outbound-Port (dort steht „Text-Stand der Quelle")
     nennen die Ausgabefunktion als Erzeugung (`AGENTS.md` §3.7 Klassen
     Zusage/Kopplung).
  2. **Fixture und Paritätstest:** die Bild-Fixture wird die Typ-Tabelle aus
     Festlegung 4 — mindestens `bool`, `inet`, `cidr`, `macaddr`, `char(n)`,
     `varchar(n)`, `money`, `bit`/`varbit`, `uuid`, `xml`, `oid`/`regclass`,
     Enum, Domains (über `bool` und `int`), Composite, Range/Multirange,
     `tsvector`/`tsquery`, geometrische Typen, `time`/`timetz`/`timestamp`,
     `int`-Typen, `real`, `name`, `"char"` und Arrays davon zusätzlich zu den
     heutigen acht Typen, `NULL`, generierter Spalte; der Vergleich läuft je
     Spalte am Roh-Text (vor `BuildRowImage`) **und** am Bild, unter den drei
     GUC-Lagen, auf PostgreSQL 17 und 18 (`make test-replication`); die
     Fehlermeldung nennt Spalte und Typ.
  3. **Strukturtest netzlos:** ein Unit-Test der Cursor-Anweisung prüft, dass
     sie an einer Spaltenliste mit einfachen Namen nur gequotete Namen wählt
     (Form `SELECT "a", "b" FROM "s"."t"`, kein `::`, kein `(`) — Träger:
     `make test`.
  4. **Plan-Träger, die den Wortlaut tragen** (Suchlauf, beide Stände unten): der
     in Umsetzung befindliche Slice-Plan des Snapshot-Lesers (Ziel-Absatz mit
     `col::text`, die DoD-Zeile „Bild-Parität" mit ihrer Typliste und die
     Mess-/Risiko-Absätze zur GUC-Parität) zieht nach; der `done/`-Record des
     gemeinsamen Row-Image-Slice nennt den Wortlaut als Nicht-Umfang und bleibt als
     Record unberührt; die Berichte unter `docs/reviews/` sind Lauf-Belege.
  5. **Spec:** `LH-FA-CAP-009.a` und alle anderen Spec-Stellen tragen kein
     `::text` (Suchlauf unten); „Werte im Text-Stand der Quelle" bleibt wahr.
     Kein Spec-Nachzug nötig; eine Schärfung („Ausgabefunktion des Typs")
     kann ein späterer Spec-Nachzug mittragen, ohne Zwang.

**Suchlauf** (`grep -rn --include=*.md -e '::text' spec docs harness AGENTS.md
README.md`, Stand des Vorgängers `5e5d0b89`, gedruckt): 9 Zeilen in vier Dateien —
2 in `ADR-0111`, 1 im Slice-Plan des Snapshot-Lesers (`in-progress/`), 1 im
`done/`-Record des gemeinsamen Row-Image-Slice, 5 im Review-Report; kein Treffer in
`spec/`, `harness/`, `AGENTS.md`, `README.md`. Stand mit dieser ADR: dieselben 9
Zeilen zuzüglich der Erwähnungen in dieser Datei. Im Code trägt
`internal/adapters/driven/postgressnapshot/snapshot.go` die Stellen (Kommentar
Zeile 255, Anweisung Zeile 260); die übrigen `::text`-Treffer in `internal/`,
`tools/` betreffen andere Abfragen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test, reale PostgreSQL (Tier `test-replication`) | **Typ-Parität:** je Spalte des Typ-Satzes (Festlegung 4, mit `NULL`-Zeile und Rand-Werten) ist der Roh-Text des Backfill-Pfads byte-gleich dem des WAL-Pfads, und das über `BuildRowImage` gebaute Bild ebenso — unter drei Rollen-GUC-Lagen (Standard, `timezone`/`datestyle`/`intervalstyle`/`bytea_output`/`extra_float_digits`, `postgres_verbose`) auf PostgreSQL 17 und 18; ein Wechsel des Lesepfads auf ein Binär-Format (A4) fällt hier rot | `make test-replication` |
| Go-Test (netzlos) | Die Cursor-Anweisung wählt nur gequotete Spaltennamen (kein `::`, kein Funktionsaufruf) | `make test` |
| Review-Prüfpflicht | Kein Lesepfad für Quelltabellen-Werte wandelt sie per Cast oder Funktion um; die Row-Image-Konstruktion bleibt an einer Stelle | — (kein Gate) |

Die Zeile „Bild-Parität" von `ADR-0111` §Fitness Function wird durch die erste
Zeile ersetzt; die übrigen Zeilen von `ADR-0111` bleiben.

## Re-Evaluierungs-Trigger

- **Eine neue PostgreSQL-Hauptversion** (`SPEC-012`-Ausweitung, z. B. 19): der
  Paritätstest läuft dort mit; ein roter Typ heißt, die Ausgabefunktion oder
  das Protokoll hat sich geändert — neu bewerten (übernommen: kein Anzeichen
  dafür in 17/18).
- **Der WAL-Pfad aktiviert die `binary`-Option von `pgoutput`** (oder ein anderes
  Format, das nicht der Text-Ausgabe entspricht): die Erzeugung der Werte
  weicht dann vom Text-Ergebnisformat ab; Backfill-Format und WAL-Format neu
  paaren (übernommen aus der `pgoutput`-Optionsliste, nicht gemessen; der
  Replication-Pfad setzt sie nicht).
- **Der Lesepfad wechselt** (anderer Treiberweg oder `COPY`, D — z. B. mit der
  Durchsatz-Ausbaustufe von `ADR-0111`): Text-Format-Zusage (Festlegung 2) und
  Paritätstest prüfen; D ist gemessen paritätsgleich, braucht aber den
  Rück-Escaping-Parser.
- **Ein Betreiber-Typ, dessen Ausgabe sitzungs-abhängig ist** (Extension-Typ mit
  eigener GUC), dessen Text sich zwischen Walsender und regulärer Sitzung
  unterscheidet: Sitzungs-Angleichung an die Walsender-Werte prüfen (Fallback
  der Vorgänger-ADR).
- Sonst permanent — die Regel „Ausgabefunktion, kein Cast" gilt unabhängig von
  Typ-Menge und Version.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-24 | Proposed — Architect-Vorschlag zum Review-Befund F-1 (Bild-Parität `col::text`); Befund an Wegwerf-Containern auf PostgreSQL 17.11 und 18.6 reproduziert, neun Lesevarianten am 86-Spalten-Typ-Satz gegen den Walsender gemessen | Review-Report zum Snapshot-Leser, Befund F-1 (`docs/reviews/`) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0115` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
