# ADR-0126: Transformationen — Annahmemenge des Parameters `rule_spec` (Supersedes ADR-0125, teilweise)

**Status:** Accepted — Supersedes [`ADR-0125`](0125-transformationen-parametertyp-regelform-json.md)
in genau einer Stelle: der Satz in Festlegung 1, der die Annahmemenge des
Parameters `rule_spec` beschreibt („Syntaktisch ungültiges JSON scheitert beim
Aufruf, ohne Antrags-Zeile; alles, was gültiges JSON ist — auch `NULL`, JSON-`null`
und ein Wert ohne Objekt —, wird als Antrag angenommen und in Go abgelehnt“). Alles
Übrige von `ADR-0125` bleibt in Kraft, insbesondere der Parametertyp `json`, die
Spalte `jsonb`, die Aufrufformen, die Messtabelle (ihre Begrenzung auf sieben Formen
war richtig) und die Nennung von `NULL`, JSON-`null` und JSON ohne Objekt als
angenommene Formen.

**Datum:** 2026-09-26

**Autor:** Architect-Agent (Modul 8), Architect-Zug zu Finding F-1 (MEDIUM, Klasse
„ADR-Aussage breiter als ihre Messung“) des Reviews
`review-slice-transformationen-antragsweg-schema`; jede Tatsachenaussage trägt ihren
Beleg-Anker oder ist als hergeleitet gekennzeichnet (siehe §Kontext).

**Bezug:** [`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Konfiguration der
Transformation; Haupt-Bezug),
[`ADR-0125`](0125-transformationen-parametertyp-regelform-json.md) (teilweise
superseded — Haupt-Bezug),
[`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(Zusage „Antrag wird angenommen, Prüfung in Go“, nur zur Abgrenzung),
[`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md) (keine Domänenlogik
in SQL)

**Schärft:** [`SPEC-019`](../../../spec/pflichtenheft.md) — der Absatz
„Transformations-Antragsarten“: welche Werte von `rule_spec` der Aufruf annimmt und
welche er ohne Antrags-Zeile ablehnt.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0125`](0125-transformationen-parametertyp-regelform-json.md) ist `Accepted` und
unberührbar (`AGENTS.md` §3.5). Sein Satz über die Annahmemenge nennt „alles, was
gültiges JSON ist“; seine Messung umfasst sieben Formen und sagt selbst, dass aus
ihnen keine Aussage über weitere Formen folgt. Der Satz reicht damit über die
Messung hinaus. Die Berichtigung ändert den Referenten (die Annahmemenge der
Funktion) und ist keine Zitat-Korrektur.

### Gemessen

**Die Funktion `cdc.set_transformation` mit dem Text aus
`tools/schema/nacharbeit-administration.sql` (der Text von `CREATE OR REPLACE FUNCTION`
bis `$$;`, mit `SECURITY DEFINER` und `SET search_path`), eine Tabelle
`cdc.administration_request` mit den Spalten der Funktion (`rule_spec jsonb`), je
eine frische PostgreSQL-Instanz: PostgreSQL 18.6 (`PG_TEST_IMAGE` aus dem
`Makefile`) und PostgreSQL 17.11 (`postgres:17-alpine`, Digest aus
`.github/workflows/e2e.yml`); der Arbeitsbaum des Repositories ist unberührt.**
Jede Form läuft als Literal vom Typ `unknown` (`SELECT cdc.set_transformation('s',
'public','t','r', <Form>)`, der Weg eines Clients mit Literal), unter dem Superuser.
Gedruckt, **identisch in beiden Instanzen** (33 Formen; ein Vergleich der
gedruckten Tabellen beider Läufe ergab keinen Unterschied):

| Form | Ergebnis |
|---|---|
| SQL-`NULL` | angenommen, Zeile mit `rule_spec` NULL |
| JSON-`null` | angenommen, gespeichert als `null` |
| `{}`, `[]`, `1`, `"x"`, `true` | angenommen |
| `{"":1}` (leerer Schlüssel) | angenommen |
| `{"kind":"a","kind":"b"}` (doppelter Schlüssel) | angenommen, gespeichert `{"kind": "b"}` (der letzte Wert) |
| `{"a":"\u00e4"}`, `{"a":"\ud83d\ude00"}` (Escape, Surrogat-Paar) | angenommen, gespeichert als `ä` bzw. `😀` |
| `{"a":"\\u0000"}` (maskierter Backslash, kein NUL-Escape) | angenommen |
| `{"a":1e400}` | angenommen |
| Schachtelung 1000 und 10000 Ebenen (`[[…]]`) | angenommen |
| Wert von 8 MiB, Objekt mit 200000 Schlüsseln | angenommen |
| `\u0000` als Escape im Wert, im Schlüssel, im Array-Element und als Skalar | „unsupported Unicode escape sequence“, keine Zeile |
| `{"a":1e200000}`, `{"a":1e-200000}`, Zahl mit 200000 Ziffern | „value overflows numeric format“, keine Zeile |
| Schachtelung 100000 Ebenen | „stack depth limit exceeded“, keine Zeile |
| einzelnes High- bzw. Low-Surrogat als Escape (`"\ud83d"`, `"\ude00"`) | „invalid input syntax for type json“, keine Zeile |
| `{oops`, leerer Text, nur Leerraum, `{"a":1}x`, `{"a":NaN}` | „invalid input syntax for type json“, keine Zeile |

Die Ablehnung des `\u0000`-Escapes und der Zahlen außerhalb des Bereichs entsteht
beim Cast `p_rule_spec::jsonb` in der Funktion, nicht beim Lesen des Parameters:
`SELECT '{"a":"\u0000"}'::json` und `SELECT '{"a":1e200000}'::json` liefern den Text,
`::jsonb` liefert den Fehler (PostgreSQL 18.6, gedruckt). Die Aussage über die
Schachtelung nennt die gemessenen Ebenen (1000 und 10000 angenommen, 100000
abgelehnt); die Grenze dazwischen ist nicht gemessen.

**Warum die Ablehnung des `\u0000`-Escapes keine Regel trifft, die ein Betreiber
setzen will** (hergeleitet aus zwei Messungen): ein Text mit dem Zeichen U+0000 ist
in PostgreSQL nicht darstellbar — `SELECT chr(0)` endet mit „null character not
permitted“, `SELECT E'a\x00b'` mit „invalid byte sequence for encoding "UTF8":
0x00“ (PostgreSQL 18.6, gedruckt). Ein Spaltenname der Quelle und ein Wert einer
Textspalte tragen es deshalb nie; eine Regel, die es nennt, hat keinen Gegenstand.
Dass kein Bezeichner-Pfad es doch trägt, ist nicht gemessen.

### Konstraints

- Keine Domänenlogik in SQL ([`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)):
  die Funktion schreibt den Antrag; die Prüfung der Regelform liegt in Go. Die Cast-Fehler
  oben sind Eigenschaft von PostgreSQL, keine Prüfung, die die Funktion selbst leistet.

## Entscheidung

Wir wählen **die Annahmemenge als das, was gemessen ist, in `SPEC-019` festzuschreiben,
ohne Änderung der Funktion**.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (Satz von `ADR-0125` und `SPEC-019` bleibt) | kein Aufwand | die Aussage ist an einer Randform falsch (`\u0000`); ein Leser leitet aus „gültiges JSON wird angenommen“ eine Zusage her, die die Funktion nicht hält |
| B — Funktion so ändern, dass `\u0000` und Zahlen außerhalb des Bereichs angenommen werden (Parametertyp `text`, Spalte `text`) | die Aussage „jedes gültige JSON“ würde wahr | Spalte, `SPEC-030`, Store und Guard-Schreibweise ändern sich für Formen, die keine Regel trägt (`\u0000` ist nicht darstellbar, s. o.); `ADR-0125` Option C lehnte den `text`-Weg aus demselben Grund ab |
| **C — Annahmemenge nach Messung festschreiben (gewählt)** | Funktion, Schema, Store, Guard unverändert; die Aussage ist an der Instanz gemessen; die abgelehnten Formen scheitern laut beim Aufruf, ohne Antrags-Zeile — sie erreichen die Queue nie | die Aussage nennt eine Ausnahmeliste, die vom PostgreSQL-Verhalten abhängt |

### Festlegung 1 — die Annahmemenge

`cdc.set_transformation` nimmt den Parameter `rule_spec` an, wenn er **`NULL`** ist
(SQL-`NULL`, kein JSON-Wert) **oder** ein Text, den PostgreSQL als `json` liest **und**
nach `jsonb` umwandelt. Angenommen sind damit auch JSON-`null`, ein Wert ohne Objekt
(Array, Zahl, Zeichenkette, `true`), das leere Objekt, ein leerer Schlüssel und ein
doppelter Schlüssel (in der Spalte gilt der letzte Wert). Der Aufruf scheitert ohne
Antrags-Zeile, wenn eine der Bedingungen nicht gilt:

1. **Syntax:** der Wert ist für PostgreSQL kein JSON — Syntaxfehler, leerer Text,
   `NaN`, einzelnes Surrogat-Escape.
2. **Umwandlung nach `jsonb`:** das Zeichen U+0000 als Escape `\u0000` in einem
   Schlüssel oder Wert; eine Zahl außerhalb des Zahlbereichs von `numeric` (gemessen:
   `1e200000`, `1e-200000` und eine Zahl mit 200000 Ziffern abgelehnt, `1e400`
   angenommen).
3. **Schachtelungstiefe:** jenseits der Stapeltiefe des Servers (gemessen: 10000
   Ebenen angenommen, 100000 abgelehnt).

Was angenommen ist, prüft Go (`SPEC-030`); die Annahme ist keine Zusage, dass der
Antrag verarbeitet wird. Diese Festlegung ersetzt den Satz über die Annahmemenge in
Festlegung 1 von `ADR-0125`.

## Konsequenzen

- Positiv: `SPEC-019` und diese Entscheidung tragen dieselbe, gemessene Aussage;
  kein Code-Stand ändert sich.
- Negativ (Grenze, benannt): die Aussage hängt am Verhalten von PostgreSQL 17.11 und
  18.6 (die beiden Digests der Spec); für andere Versionen ist sie nicht gemessen.
  Die Grenzen unter 2 (Zahlbereich) und 3 (Tiefe) sind Eigenschaften von
  PostgreSQL, ihre genauen Werte stehen hier nicht.
- Folgepflicht 1: **`SPEC-019` Nachzug** — der Absatz „Transformations-Antragsarten“
  trägt die Annahmemenge dieser Festlegung (in demselben Commit wie diese ADR).
- Folgepflicht 2: **Store-Test der Annahmemenge** — ein Test im Paket
  `postgresstorage` (`make test-store`) ruft die Funktion mit den Formen der
  Messtabelle auf und prüft je Form „Zeile vorhanden mit gelesenem Wert“ bzw.
  „Fehler, keine Zeile“; er trägt die Fitness-Function-Zeile unten. Träger:
  Implementer des Slice `slice-transformationen-antragsweg-schema`, Fixrunde.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test, reale PostgreSQL (Paket `postgresstorage`), **noch nicht vorhanden — Folgepflicht 2** | Tabellentest je Form der Messtabelle: `NULL`, JSON-`null`, `{}`, `[]`, `1`, `"x"`, doppelter Schlüssel angenommen (Zeile, gelesener Wert); `\u0000` im Wert, im Schlüssel, im Array-Element und als Skalar, `1e200000`, `{oops`, leerer Text, einzelnes Surrogat abgelehnt (Fehler, `count` der Zeilen des Aufrufs 0). Die Aussage dieser ADR ist bis zum Test **an der Instanz gemessen** (Messtabelle oben), nicht durch einen Test im Repository gebunden; die Mutation des Tests (Cast `::jsonb` aus der Funktion entfernt: die `\u0000`-Zeilen müssen rot färben) fährt der Implementer beim Schreiben | `make test-store` |

## Re-Evaluierungs-Trigger

- **Der Parametertyp oder der Cast der Funktion ändert sich** (Folge-ADR zu
  `ADR-0125`, z. B. nach dem Trigger „d-migrate rendert den Abbau mit `jsonb`“): die
  Messtabelle wiederholen.
- **Der gepinnte PostgreSQL-Digest wechselt die Hauptversion** (`SPEC-012`): die
  Messtabelle wiederholen.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-26 | Accepted — Architect-Zug (Vollmacht des Auftraggebers) zu Finding F-1 des Reviews `review-slice-transformationen-antragsweg-schema`; 33 Formen an PostgreSQL 18.6 und 17.11 gemessen, jede Tatsachenaussage an einer gedruckten Zeile oder als hergeleitet gekennzeichnet | [`LH-FA-CFG-007`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0126` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
