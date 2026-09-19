# Beleg: slice-sdk-csharp-projektgeruest

Vorgang: `slice-sdk-csharp-projektgeruest` (SDK-Projektgerüst
`sdks/csharp/PgChangeFeed.Client/`).

Fund: Der Implementer legte das erste Testprojekt dieser neuen Sub-Area
(`PgChangeFeed.Client.Tests`) mit `ImplicitUsings` an, ohne vorher zu
prüfen, ob derselbe Bestand das Problem „xUnit-Testdateien brauchen ein
explizites `using Xunit;`, `ImplicitUsings` deckt es nicht ab" bereits
gelöst hat — es tut es: alle 13 bestehenden Testdateien unter
`examples/csharp/*/*.Tests/*.cs` tragen bereits durchgängig ein
explizites `using Xunit;` (`grep -rl "using Xunit"
examples/csharp --include="*.cs"`, 13 Treffer, keiner ohne). Ein
`grep -rl "using Xunit" examples/csharp` vor dem ersten eigenen Testlauf
hätte das Muster gezeigt; stattdessen kostete es einen ersten roten
Docker-Build im Implementer-Lauf, bevor die (bereits im Bestand
etablierte) Lösung — explizites `using Xunit;` — nachgetragen wurde.
Anders als bei `slice-104`: Die *endgültige* Lösung deckt sich hier mit
dem etablierten Muster (kein abweichender Workaround blieb liegen); der
vermeidbare Aufwand ist die eine rote Bau-Iteration, nicht ein
dauerhaft suboptimaler Code-Zustand.

Aufgefallen ist es im Verifikations-/Closure-Schritt (Planner-Auftrag
für die Closure benannte den Fund explizit als Prüfpunkt für das
Beobachtungs-Register), nicht im Review — der Review-Report
`docs/reviews/review-slice-sdk-csharp-projektgeruest.md` <!-- d-check:status-provenance -->
bestätigt nur, dass die finale Fassung fehlerfrei ist, ohne die
vermeidbare erste rote Iteration als eigenen Fund zu benennen.

Quelle: Implementer-Bericht zu `slice-sdk-csharp-projektgeruest`
(erster roter Docker-Build wegen fehlendem `using Xunit;`);
`grep -rl "using Xunit" examples/csharp --include="*.cs"` (13 Treffer,
2026-09-19).
