**Vorgang:** slice-103

**Fund:** Review-F-1 (`docs/reviews/review-slice-103.md`, ausdrücklich als
„Fortsetzung von `review-slice-102.md` F-1" geführt): eigener Nachbau
bestätigt denselben Effekt auf `examples/kotlin/Dockerfile` —
`docker build examples/kotlin` ohne `--build-context proto=proto` bricht
bereits beim Auflösen der `docker.io/library/proto:latest`-Referenz ab,
unabhängig vom angeforderten `--target`, weil `COPY --from=proto …`
unbedingt in der von allen vier Modulen geteilten `build`-Stufe steht.
Zweite Instanz derselben strukturellen Ursache (eine gemeinsame
`build`-Stufe je Sprache). Anlass für die Anlage dieses Registereintrags
bei der `slice-103`-Closure (Reviewer-Empfehlung: „Empfehlung an die
Planner-Closure: neues Verzeichnis
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut/` … anlegen").

Quelle: `docs/reviews/review-slice-103.md` F-1.
