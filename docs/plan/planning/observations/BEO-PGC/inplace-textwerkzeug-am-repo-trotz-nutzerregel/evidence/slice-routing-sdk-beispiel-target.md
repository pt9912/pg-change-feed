**Vorgang:** slice-routing-sdk-beispiel-target (Review-Bericht, Kopf „Ablage“)

**Fund:** Der Reviewer fuhr eine Mutationsreihe an den Go-Beispielen. Ein Mutationsaufruf lief wegen eines
fehlgeschlagenen `cd` im Arbeitsbaum statt in der Kopie im Scratchpad und änderte drei Repo-Dateien
(`examples/grpc-client/stream.go`, `examples/http-client/changes.go`, `examples/nats-stream-client/subject.go`).
Der Reviewer sah es an `git status`, nahm die Änderungen mit `git checkout` zurück, fuhr die Reihe erneut
in der Kopie und stellte `git status --short` leer fest; der Report nennt den Fehlgriff selbst. Welches
Werkzeug schrieb, nennt der Report nicht (nicht gemessen).

**Form (Ausprägung):** Mutationsprobe mit relativem Pfad nach einem fehlgeschlagenen Verzeichniswechsel,
mit Wirkung auf Repo-Dateien, vor dem Commit selbst bemerkt und zurückgenommen. Kein Host-Interpreter am
Kopf, also kein Beleg für das erste Neubewertungs-Kriterium der Kopf-Liste (`MR-003`/`MR-004`); die
Regel (`AGENTS.md` §3.1: eine Mutationsprobe arbeitet auf einer Kopie im Scratchpad) war formuliert, der
Weg liegt jenseits dessen, was der Guard liest. Die Rücknahme war wirksam, weil `git status` nach der
Reihe geprüft wurde.

Quelle: `docs/reviews/review-slice-routing-sdk-beispiel-target.md` (Kopf „Ablage“) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-sdk-beispiel-target.md` (Vorspann: `git status --short` im Repo nach allen eigenen Mutationen leer). <!-- d-check:status-provenance -->
