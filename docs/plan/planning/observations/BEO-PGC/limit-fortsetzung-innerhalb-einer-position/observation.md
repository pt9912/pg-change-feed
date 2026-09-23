# BEO-PGC/limit-fortsetzung-innerhalb-einer-position

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft den
Bereichs-Lesezugriff `GET /changes` und `cdc.changes`).

Die Beobachtung: Das Bereichslesen mit `Limit` kann innerhalb **einer**
Commit-Position nicht fortsetzen. Die Abfrage `SelectChanges`
(`internal/adapters/driven/postgresstorage/queries/queries.go`) schneidet den
Bereich über `t.commit_position >= $2 AND t.commit_position < $3`, sortiert
nach `(commit_position, transaction_id, sequence)` und begrenzt mit
`LIMIT $6`. Der einzige Fortsetzungswert ist eine Position; schneidet das
Limit mitten in einer Transaktion ab, gibt es keinen Cursor auf den Rest:
ein Folgelesen ab `commit_position + 1` überspringt ihn, ein Folgelesen ab
`commit_position` liefert dieselben Zeilen erneut.

Aus der Abfrage abgeleitet, **nicht** real reproduziert. Die Wirkung auf
einen Consumer (Rest verloren nach einer Bestätigung der Position, oder
Stillstand bei einer Transaktion größer als das Limit) ist damit erwartet,
nicht gemessen. Sie gilt für jede große WAL-Transaktion und wird durch einen
Backfill (ein Block ist eine synthetische Transaktion) zur Regel.

Bindung: `ADR-0081` (Lesen über die HTTP-API, Cursor-Trigger) und
`ADR-0111` Folgepflicht 6 (Cursor-Form; bis dahin die Lese-Regel
„Bestandsabzug ohne `Limit`").

## Benannt, nicht gezählt

Gefunden beim Schreiben von `ADR-0111` (Architect-Lauf, 2026-09-23); kein
abgeschlossener Vorgang, deshalb keine `evidence/`-Datei.
