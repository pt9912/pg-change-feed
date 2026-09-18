# Beleg: slice-087

Vorgang: `slice-087` — der E2E-Beleg für `GET /changes`.

Fund: Der Beleg prüft, dass **gelesen** wurde — aber nicht, dass der **Filter**
gewirkt hat. Die Phase liest mit `[from, to)` auf genau den `commit_position`
ihrer eigenen Sentinel-Zeile; der Bereich hat damit die **Breite 1**, und in
dieser Aufruf-Form ist „der Filter wirkt" nicht von „der Filter fehlt" zu
unterscheiden.

Gefunden hat es der Reviewer **nicht durch Lesen, sondern durch Mutieren der
Eingabeseite**: die Mutation „Filterachse entfällt" ließ den vollen
`make test-integration` **grün** (Review zu `slice-087`, F-1, Mutation C).

Damit ist dies die **zweite Gelegenheit desselben Mechanismus an einem anderen
Gegenstand**: `slice-086` betraf einen **Negativtest** (die Ablehnung hing an
keinem Eingabewert), `slice-087` einen **E2E-Beleg** (die Filterprüfung hing
nicht am Filter). Beide Male war der Beleg grün, egal was die Eingabe tat.

Quelle: Review zu `slice-087` (F-1 mit eigener Mutation C) ·
`tools/harness/run-integration-tests.sh` (die HTTP-Phase) ·
`tools/harness/httpclient/main.go` (der Lese-Aufruf).
