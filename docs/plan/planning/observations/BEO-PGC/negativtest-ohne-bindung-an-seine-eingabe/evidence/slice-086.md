# Beleg: slice-086

Vorgang: `slice-086` — das Changes-Lesen über die HTTP-API.

Fund: Der Adapter-Test für `?limit=0 → 400`
(`internal/adapters/driving/http/readchanges_test.go`) stellte gegen einen Fake,
der `outbound.ErrNonPositiveLimit` **unabhängig von der Abfrage** lieferte — der
HTTP-Status hing damit an **keinem** Parameterwert.

Gefunden hat es der Reviewer **nicht durch Lesen, sondern durch Mutieren der
Eingabeseite**: `readChangesLimit("0") → nil` ließ die **ganze Suite grün**,
obwohl der Endpunkt danach bei `?limit=0` **still unbegrenzt** liest (Review
zu `slice-086`, F-2, eigene Mutation D).

Die Fixrunde hat die Kette geschlossen und ihre drei Glieder je an ihren Träger
gestellt: **Eingabe → unveränderte Weitergabe** (Adapter, jetzt getestet) →
`ErrNonPositiveLimit` (Port-Kontrakt `ChangeQuery.Validate()`) → **`400`**
(Adapter). Die Mutation des Reviewers wurde danach **rot** gesehen.

Quelle: Review zu `slice-086` (F-2 mit eigener Mutation) ·
Implementer-Bericht der Fixrunde (Commit `f20d2c9`) ·
`internal/adapters/driving/http/readchanges_test.go`.
