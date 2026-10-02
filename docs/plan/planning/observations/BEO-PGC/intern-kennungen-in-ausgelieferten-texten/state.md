Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen. Behoben im
Vorgang: `git grep` der Kennungs- und Slice-/Welle-Muster über `sdks` (ohne
`grpc_gen`) druckt am Parent 499 Zeilen in 98 Dateien, am Ende des Slice 0; die
fünf ausgelieferten Artefakte (Wheel, sdist, `.nupkg` samt XML-Datei und `.dll`,
Jar, Sources-Jar) tragen 0 Treffer (Verifikation §5). Der Wächter greift an der
Eingabeseite (Mutationen G und H der Verifikation rot).

**Kandidat einer Regel (Planner-Entscheidung, 1× — kein `AGENTS.md`-Eintrag):** „Texte,
die das Repo verlassen (Paketbeschreibung, öffentliche API-Doku, Fehlertexte,
Kommentare erzeugter Stubs), tragen keine internen Kennungen.“ Das Register führt
Regeln erst ab 3×; die Klasse hat einen Wächter je Träger, der Regelwortlaut wartet
auf die zweite Ausprägung außerhalb der SDKs (ein weiteres ausgeliefertes Erzeugnis:
das Server-Image, die Docker-Hub-Beschreibung).

Zähler (abgeleitet): **1×** (evidence/slice-sdk-readme-nutzerdoku.md).

**Fangnetz für `docs/user/`:** `make handbuch-public-doc-check` (Gate in `make gates`) prüft
`benutzerhandbuch.md`, `benutzerhandbuch-standard.md` und `version.md` auf Kennungen und auf
Links nach `docs/plan/` und `docs/reviews/`; der Skill `.harness/skills/nutzerdoku-schreiben.md`
trägt die Schreib-Seite. Grenze: das Gate liest Kennungen, nicht Sinn (Chronik-Sprache ohne
Kennung bleibt grün). Die Texte, die das Repo verlassen, haben damit zwei Wächter
(`make sdk-public-doc-check` für `sdks/`, dieses Gate für `docs/user/`); der Zähler bleibt
**1×**, weil kein neuer Vorgang dieser Klasse dazukommt. Beide Wächter melden einen
Lesefehler als Exit 2; beim SDK-Wächter geschlossen am Commit `8f6430a0`
(`docs/plan/planning/done/slice-sdk-public-doc-check-lesefehler-fail-closed.md`).

**Verwandt, nicht gleich:** `BEO-PGC/slice-chronik-in-code-kommentar` (Chronik in
Code-Kommentaren, Leser ist der Entwickler im Repo) — hier ist der Leser ein
Anwender außerhalb des Repos, und der Träger ist ein Paket.
