**Vorgang:** slice-code-kommentare-kennungen (Plan §6, Risiko „Der `DIFF`-Strom läuft über
Host-`git` und `bash`“)

**Fund:** Der Aufrufer `tools/harness/kommentar-kennungen.sh` fährt `git diff` auf dem Host und
läuft unter `bash`; der Plan führte die Host-Werkzeug-Frage in §6 als Risiko, mit der Erwartung,
der Vertrag nenne die Host-Werkzeuge und ein zweites Auftreten der Klasse sei die Architect-Frage
zu `AGENTS.md` §3.1. Der Vertrag `harness/sensors/kommentar-kennungen.md` nennt `bash`, `git` und
`docker` (Reviewer, Negativbefund: geprüft, ohne Befund). Das ist das zweite Auftreten der Klasse
nach `slice-harness-suchlauf-nachmessen`; ein drittes Plan-Risiko dieser Klasse führt
`slice-harness-fmt-check` §6 (Slice in `open/`, noch keine Belegdatei).

Quelle: Plan `slice-code-kommentare-kennungen` §6 · Review-Report
`review-slice-code-kommentare-kennungen` (Negativbefund `kommentar-kennungen.sh`).
