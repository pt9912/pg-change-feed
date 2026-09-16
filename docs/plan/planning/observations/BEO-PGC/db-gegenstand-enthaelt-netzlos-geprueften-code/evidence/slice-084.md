# Beleg: slice-084

Vorgang: `slice-084` — die Naht in `internal/adapters/driven/postgresack`.

Fund: Der Gegenstand der DB-Adapter-Coverage enthält nach diesem Zug Code, der
**ohne jede PostgreSQL-Verbindung** gedeckt wird. Gemessen (Reviewer **und**
Implementer, unabhängig): von den **32** Statements des Pakets sind **30**
netzlos gedeckt — `postgresack` steht bei **32/32**, und die Zahl ist damit
nicht mehr von einer echten Verbindung zu unterscheiden.

Die Bewegung in Zahlen (Null-Befund des Slice, beide Größen gemessen):

| | vor | nach |
|---|---|---|
| DB-Nenner gesamt | 650 | **659** (+9, der neue Paketcode) |
| davon `postgresack` | 23 | **32** |
| DB-Adapter-Quote | 73,38 % | **74,51 %** |
| Unit-Nenner (netzlos) | 1903 | 1903 (Δ 0) |

**Das Spiegelbild steht in `slice-081`:** dort hat ein **Transfer** den
DB-Nenner gedräniert (138 Statements wanderten aus dem Gegenstand, die Quote
fiel auf 73,38 %) — und das wurde **entschieden** (`ADR-0077`, teilweise
abgelöst durch `ADR-0078`). Hier wächst der Nenner und die Quote **steigt**,
ohne einen einzigen neuen DB-gestützten Beleg. Für die Abwärtsbewegung hält
`ADR-0078` einen Riegel bereit („fällt der Nenner ohne Ankunft, steht die
Schwelle"); für die **Aufwärtsbewegung** gibt es keinen.

`ADR-0080` hat die **Verdünnung** des DB-Nenners für dieses Paket ausdrücklich
als Trigger benannt; sie ist mit diesem Vorgang **eingetreten**.

Quelle: `docs/reviews/review-slice-084.md` (Schwerpunkt 4, mit eigener Messung
30/32) · Implementer-Bericht `slice-084` (unabhängige Messung im vollen
`-coverpkg`-Lauf) · `harness/sensors/db-adapter-coverage.md` §Zählbasis.
