# Beleg: slice-sdk-kotlin-projektgeruest

`review-slice-sdk-kotlin-projektgeruest.md` befundet real (F-2, LOW) das
deutsche Fachwort „vollinhalt" im README-Satz zu SSE-/NATS-Umfang
(`sdks/kotlin/pgchangefeed-kotlin/README.md:9`) — „SSE and NATS-vollinhalt
delivery remain out of scope for this package's planned first full release
…" — als dritte reale, aber erste formal befundete Instanz derselben
Formulierung in `sdks/{csharp,python,kotlin}/README.md`. Real
nachgemessen: `grep -n "vollinhalt"
sdks/kotlin/pgchangefeed-kotlin/README.md` → Zeile 9.

Der Reviewer verlangte keine Fixrunde (0 HIGH/MEDIUM, LOW ohne
semantische Auswirkung), empfahl aber ausdrücklich, ein viertes Auftreten
solle die §Pflege-Schärfung auslösen. Diese Planner-Closure legt den
Eintrag stattdessen bereits jetzt mit dem realen 3×-Stand an (siehe
`state.md`), statt auf ein hypothetisches viertes SDK zu warten.
