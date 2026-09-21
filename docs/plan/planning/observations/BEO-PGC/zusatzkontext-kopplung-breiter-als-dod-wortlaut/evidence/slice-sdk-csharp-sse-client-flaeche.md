**Vorgang:** slice-sdk-csharp-sse-client-flaeche

**Fund:** Verifier-Hinweis (Verifikationsbericht zu
`slice-sdk-csharp-sse-client-flaeche`, §Nicht blockierende Beobachtung):
`sdks/csharp/Dockerfile` trägt — wie `examples/csharp/Dockerfile` und
`examples/kotlin/Dockerfile` — eine einzige `build`-Stufe für den ganzen
Baum; `COPY --from=proto …` liegt darin unbedingt, deshalb ist
`--build-context proto=proto` für **jeden** `docker build`-Aufruf dieses
Dockerfiles zwingend, nicht nur für einen gRPC-spezifischen Bauweg. Anders
als bei den beiden bisherigen Belegen ist der betroffene Baum kein
Beispiel-Client (`examples/**`), sondern ein SDK-Package-Baum
(`sdks/csharp/**`) — **dritte**, strukturell andersartige Instanz derselben
Ursache (eine gemeinsame Bau-Stufe koppelt den Zusatzkontext breiter, als
der jeweils enger klingende DoD-Wortlaut vermuten lässt).

Der Planner hat den Ursprung real per `git log --diff-filter=A --oneline --
sdks/csharp/Dockerfile` und `git log --oneline -- sdks/csharp/Dockerfile`
nachgezogen: Die `COPY --from=proto …`-Zeile kam mit
`06e78545 feat(sdk): C#-gRPC-Client-Fläche für PgChangeFeed.Client` (Slice
`slice-sdk-csharp-grpc-client-flaeche`) in die einzige `build`-Stufe — zu
diesem Zeitpunkt trug weder Review noch Verifikation dieses Slice einen
Befund dazu (kein `F-*` zur Kopplung in
`docs/reviews/review-slice-sdk-csharp-grpc-client-flaeche.md` <!-- d-check:status-provenance -->/
`docs/reviews/verifikation-slice-sdk-csharp-grpc-client-flaeche.md` <!-- d-check:status-provenance -->) — die
Kopplung war zu diesem Zeitpunkt korrekt, weil noch keine dritte Fläche
(SSE) existierte, die sie hätte sichtbar machen können. Der Slice-Plan von
`slice-sdk-csharp-sse-client-flaeche` selbst (§3/§6) benannte die Kopplung
bereits transparent, bevor der Implementer-Zug begann, und die
Welle-Eröffnungs-Sichtung von `welle-sdk-csharp-vollabdeckung` (§6) sah sie
bereits als möglichen dritten Treffer voraus („ein dritter Treffer würde
die 3×-Schwelle erreichen"). Der Verifier bestätigte real, dass der Bau
ohne den Zusatzkontext an derselben `COPY --from=proto`-Zeile abbricht
(Verifikationsbericht §3.3, Exit 1, Dockerfile-Zeile 57) — vierte
unabhängige Bau-Bestätigung dieses konkreten Slice, aber die **erste**
außerhalb eines `examples/**`-Baums.

**Einordnung gegenüber den ersten beiden Belegen:** Bei
`examples/csharp`/`examples/kotlin` sind die vier bzw. fünf Programme einer
Sprache **eigenständige, unabhängig gedachte** Ziele (je ein `.csproj`),
die eine gemeinsame Docker-Stufe unnötig koppelt — eine vermeidbare
Bau-Entscheidung. Bei `sdks/csharp` ist HTTP+gRPC(+SSE) **ein einziges**
Package/`.csproj` (`ADR-0106` Festlegung 1, „Ein Package, nicht vier" —
bewusste Design-Entscheidung, keine Docker-Stufen-Nachlässigkeit); ein
`dotnet build`/`test` dieses einen Projekts kompiliert immer den ganzen
Baum, unabhängig vom Docker-Stufen-Design. Die **Ursache** ist damit nicht
identisch (vermeidbare Stufen-Kopplung vs. eine bereits vom Architect
gewollte Ein-Package-Bündelung), das **Symptom** — ein enger klingender
DoD-/Plan-Satz trifft die tatsächliche Zwangsbreite des Zusatzkontexts
nicht — ist es. Der Planner zählt diesen Beleg trotzdem zur bestehenden
Beobachtung, weil sowohl der Slice-Plan als auch die Welle-Eröffnung ihn
bereits selbst dieser Klasse zugeordnet hatten, bevor der Implementierungs-
Zug begann; die Ursachen-Nuance wird hier festgehalten, damit ein späterer
Lese-Schritt (Welle-Closure) sie nicht neu herleiten muss.

Quelle: Verifikationsbericht zu `slice-sdk-csharp-sse-client-flaeche`
(§Nicht blockierende Beobachtung), Planner-Nachweis per `git log` bei der
Slice-Closure.
