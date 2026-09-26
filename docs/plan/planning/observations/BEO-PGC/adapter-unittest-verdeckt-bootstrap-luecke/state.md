Zustand: offen (**3×**) — die Schwelle ist erreicht, ein Ausgang ist nicht zugewiesen: der
Lese-Schritt der Closure von `welle-transformationen` liest den Eintrag. Ein möglicher
Träger ist ein schlanker Bootstrap-Smoke-Test (`ConfigFromEnv` + `Run()`-Teilaufruf gegen
einen minimalen ENV-Satz, der prüft, dass kein Adapter-Konstruktor einen
`nil`-Use-Case erhält) — das entscheidet erst ein Architect-Zug, kein
Slice hat das bislang vorgeschlagen. Kein Slice ist wegen des Eintrags fällig: der dritte
Beleg trifft nicht die Verdrahtung eines Adapters, sondern die Aufrufstelle einer
extrahierten Start-Sequenz in `Run`; sie ist über zwei Quelltext-Tests gebunden, und den
Goroutinen-Start deckt am System `make test-integration` (Mutation „startet nie“ rot,
Verifikation V10b, **übernommen**). Ein Smoke-Test der Aufrufstelle braucht die Datenbank wie
`Run` selbst (der Plan des Slice begründet die Quelltext-Tests damit, §3; hergeleitet, nicht
an einem Smoke-Test erprobt). Zähler
(abgeleitet): **3×** (evidence/slice-061.md, evidence/slice-072.md,
evidence/slice-transformationen-start-reihenfolge.md).
