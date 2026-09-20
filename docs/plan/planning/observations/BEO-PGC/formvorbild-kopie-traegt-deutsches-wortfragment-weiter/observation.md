# BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter

**Sub-Area:** `sdks/*/` (die SDK-Sprachpakete dieser Repo-Reihe —
`harness/conventions.md` §Modus-Deklaration, Default-Sub-Area `*`/`PGC`
Greenfield; ein neues SDK-Sprachpaket entsteht ausdrücklich mit einem
bereits bestehenden Sprachpaket als Formvorbild).

Die Beobachtung: Der englischsprachige Klassen-Doc-Kommentar von
`PgChangeFeedClientOptions` trägt in **zwei** der drei bislang existierenden
SDK-Sprachpakete dasselbe deutsche Wortfragment „unstrittige" mitten im
sonst durchgängig englischen Satz („this is the one unstrittige, shared
configuration denominator identified while writing this project
skeleton"). Der Ursprung ist real auffindbar: Der jeweilige Slice-Plan
formuliert an genau dieser Stelle auf Deutsch „ein echter, unstrittiger
gemeinsamer Nenner" (§1 Ziel/Abgrenzung) — das Wort wanderte von der
deutschen Plan-Prosa unübersetzt in den englischen Quellcode-Kommentar. Das
C#-SDK (`sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs:7`)
führte den Fehler zuerst ein; das Kotlin-SDK
(`sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/PgChangeFeedClientOptions.kt:9`)
übernahm denselben englischen Kommentartext wortgleich als Formvorbild und
kopierte damit auch das deutsche Wortfragment. Das Python-SDK trägt an der
entsprechenden Stelle eine andere Formulierung ohne dieses Wort — kein
Treffer dort (real geprüft, `grep -rn "unstrittige" sdks/python/`).

**Warum das zählt:** Wortgleiche Übernahme zwischen SDK-Sprachpaketen ist
die etablierte, gewünschte Praxis dieses Repos (Formvorbild-Prinzip,
`ADR-0106`/`ADR-0107`/`ADR-0109` je Festlegung 3) — sie trägt dabei aber
unbemerkt auch Fehler weiter, nicht nur korrekte Muster. Ein Review, das
nur den Diff des **neuen** Pakets liest, findet den übernommenen Fehler nur,
wenn es das kopierte Vorbild selbst noch einmal eigenständig auf
Sprachreinheit prüft, statt die wortgleiche Übereinstimmung mit dem Vorbild
als Bestätigung zu werten.

Deklaration: `slice-sdk-csharp-projektgeruest` (Erstauftreten, formal
befundet als F-1 in `review-slice-sdk-csharp-projektgeruest.md`, damals ohne
eigenen Beobachtungs-Register-Eintrag), `slice-sdk-kotlin-projektgeruest`
(zweites Auftreten, formal befundet als F-1 in
`review-slice-sdk-kotlin-projektgeruest.md`, Eintrag angelegt bei dieser
Planner-Closure, 2026-09-20).
