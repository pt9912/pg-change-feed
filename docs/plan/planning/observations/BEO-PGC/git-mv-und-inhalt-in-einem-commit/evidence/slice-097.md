**Vorgang:** slice-097

**Fund:** Commit `1c1d937` verschiebt drei erzeugte Dateien
(`internal/adapters/driving/grpc/streamv1/*` → `gen/cdc/stream/v1/*`, von Git
korrekt als Umbenennung mit 98–100 % Similarity erkannt) und ändert im selben
Commit den Inhalt zweier dieser Dateien (die im Raw-Descriptor eingebettete
`go_package`-Zeichenkette in `changestream.pb.go`, der Import-Pfad in
`changestream_test.go`) sowie sechs weitere Dateien (`.a-check.yml`, vier
Importstellen, `proto/cdc/stream/v1/changestream.proto`) — statt Move-Commit
und Inhaltsänderung als zwei Commits zu führen (`AGENTS.md` §3.3, Regelfall).
Vom Reviewer als HIGH F-1 erkannt (`docs/reviews/review-slice-097.md`), vom
Verifier unabhängig bestätigt (`docs/reviews/verify-slice-097.md` §4) — beide
stufen den Fund als nicht-merge-blockierend und nicht-fix-pflichtig ein, weil
Git die Umbenennung trotzdem korrekt erkannte und `git log --follow` nicht
bricht.

Quelle: `docs/reviews/review-slice-097.md` (F-1) ·
`docs/reviews/verify-slice-097.md` (§4) · `AGENTS.md` §3.3.
