# Beleg: slice-085

Vorgang: `slice-085` — die Naht in `replication/receive`.

Fund: **Die zweite Naht derselben Form, und die Verdünnung springt von einem
Sonderfall zu einem Viertel des Gegenstands.** Gemessen (Unabhängig: Reviewer
und Verifier, je eigene Läufe):

| | `slice-084` (postgresack) | `slice-085` (receive) |
|---|---|---|
| netzlos gedeckt im Paket | 30/32 (**93,75 %**) | 112/187 (**59,89 %**) |
| netzlos gedeckt im ganzen Gegenstand | 41/491 (**8,35 %**) | 142/532 (**26,69 %**) |
| DB-Nenner | 650 → 659 | 659 → **691** |
| DB-Quote | 73,38 % → 74,51 % | 74,51 % → **76,99 %** |

**Von den 153 gedeckten Statements des größten Pakets brauchen nur 41 eine
PostgreSQL-Instanz.** Nach zwei Nähten ist also **mehr als ein Viertel** der
gedeckten Statements des DB-Gegenstands netzlos erreicht — und für die Quote
zählt jedes davon wie ein real geprüftes.

**Das ist die erste belastbare Zahl für den Trigger (a) aus [`ADR-0080`](../../../../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)**
(„materiell gewordene Verdünnung"): mit diesem Vorgang liest sich der
**DB-Hochschalt-Trigger** als fällig — 76,99 % ≥ 75 %, die nächste volle
5-%-Stufe über 70. **Nach dem Buchstaben ist er es; nach der Property nicht**,
denn der Anstieg ist Verdünnung, kein Ausbau. Diese Entscheidung gehört der
Wellen-Closure (`Trigger-Audit`, Modul 6 Schritt 2), nicht diesem Vorgang.

**Ein zweiter Slug wäre die Umformulierung, die die Register-Regel verbietet:**
es ist **eine** Beobachtung — der Mechanismus (netzlos geprüfter Code im
DB-Gegenstand) ist derselbe, nur die Größenordnung hat sich geändert. Der Zähler
zählt beide Vorgänge.

Quelle: Review zu `slice-085` (Schwerpunkt 5, eigenes Urteil) ·
Verifikationsbericht zu `slice-085` (eigene Messung des Anteils) ·
`harness/sensors/db-adapter-coverage.md` §Zählbasis.
