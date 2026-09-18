# BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Dockerfiles der fremdsprachigen Beispiel-Client-Wurzeln, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: In den Dockerfiles der fremdsprachigen Beispiel-Client-
Wurzeln (`examples/csharp/Dockerfile`, `examples/kotlin/Dockerfile`) trägt
eine gemeinsame `build`-Stufe alle vier Programme einer Sprache. Der
benannte Zusatzkontext (`--build-context proto=proto`,
[`ADR-0090`](../../../../adr/0090-beispiel-clients-volle-matrix.md)
Festlegung 2), den `ADR-0090` §Fitness Function nur für den `grpc`-Bauweg
als zwingend beschreibt, ist durch diese gemeinsame Stufe strukturell für
**alle vier** `docker build`-Aufrufe der Sprache zwingend — nicht nur für
`grpc`. Der DoD-Wortlaut („für grpc zusätzlich") ist damit enger als das
tatsächliche Verhalten; die Abweichung ist in beiden Fällen transparent im
Dockerfile bzw. in `harness/mk/examples.mk` kommentiert, kein stiller
Defekt.

Zwei unabhängige Vorgänge, dieselbe strukturelle Ursache: `slice-102`
(Review-F-1, C#-Zelle) und `slice-103` (Review-F-1, Kotlin-Zelle,
ausdrücklich als „Fortsetzung von F-1 des Reviews zu `slice-102`" geführt).

Deklaration: Reviewer, Review zu `slice-103`
(F-1), Anlage durch den Planner bei der `slice-103`-Closure — `slice-102`s
eigene Instanz desselben Fundes wurde in dessen Closure-Notiz nur als Prosa
behandelt, nie als eigener Registereintrag angelegt.
