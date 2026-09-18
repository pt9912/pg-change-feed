**Vorgang:** slice-103

**Fund:** Review-F-1 (Review zu `slice-103`, ausdrücklich als
„Fortsetzung von F-1 des Reviews zu `slice-102`" geführt): eigener Nachbau
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

Quelle: Review zu `slice-103`, F-1.
