**Vorgang:** slice-104

**Fund:** Die neue Dockerfile-Stufe `proto-export` (mount-loser
`proto-generate`-Umbau, [`ADR-0060`](../../../../../adr/0060-grpc-streaming-mechanismus.md))
braucht die `.proto`-Quelle im Build-Kontext (`COPY proto/ proto/`).
`.dockerignore`s Allow-Listen-Default-Deny schloss das Verzeichnis `proto/`
zunächst aus — reale Gegenprobe (Reviewer): ohne eine `!proto/`-Ausnahmezeile
scheitert der Build real mit „`COPY failed: … "/proto": not found`". Von
Implementer, Reviewer (`docs/reviews/review-slice-104.md` F-2) und Verifier
(`docs/reviews/verify-slice-104.md` §6) unabhängig geprüft; beide Rollen
stimmen überein, dass dies **nicht** die dritte Instanz von
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` ist (jener Eintrag ist
auf ein Skript unter `tools/` verengt), sondern der Erstbeleg der hier
generalisierten Klasse — ein Quellverzeichnis für eine neue Erzeugungsstufe,
kein Werkzeug-Skript.

Quelle: `docs/reviews/review-slice-104.md` F-2 · `docs/reviews/verify-slice-104.md`
§6 · Planner-Closure-Entscheidung.
