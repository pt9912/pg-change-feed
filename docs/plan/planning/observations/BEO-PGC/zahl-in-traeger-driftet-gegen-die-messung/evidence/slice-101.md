# Beleg: slice-101

Vorgang: `slice-101` — NATS-Client in C# und Kotlin (`ADR-0090`).

Fund: Der Slice-Plan selbst (Kopf §Bezug-Nachbarschaft in §6 Risiko 3 und
§8 „Vorgelagert — offene Beobachtungen sichten") zitierte
`BEO-PGC/github-actions-unverifizierbar-lokal` mit „**5×**, verkörpert in
`AGENTS.md` §3.10". Der Verifier hat den abgeleiteten Zähler eigenständig
nachgemessen (`ls evidence/` in diesem Verzeichnis) und real **7×**
gefunden (`evidence/slice-039.md`, `-056.md`, `-064.md`, `-082.md`,
`-090.md`, `-098.md`, `-099.md`) — bereits zum Zeitpunkt der
Plan-Niederschrift (2026-09-17, nach `slice-098`/`-099`), also eine
**übernommene, nicht neu gemessene Zahl** statt einer Ableitung aus den
tatsächlich vorhandenen Belegdateien (Verifikationsbericht zu `slice-101`,
§5).

**Abgrenzung zu `BEO-PGC/arbeit-ueberholt-stehenden-traeger`:** Dort
überholt die *Arbeit dieses Slice* einen fremden, bei Niederschrift noch
korrekten Träger (`harness/README.md`). Hier ist es umgekehrt: Die im Plan
selbst genannte Zahl war **schon bei ihrer eigenen Niederschrift** falsch —
kein Drift durch spätere Arbeit, sondern eine Übernahme, die nie gegen die
tatsächliche Beleg-Liste nachgezählt wurde. Das ändert am Ausgang nichts
(der Eintrag `github-actions-unverifizierbar-lokal` ist unabhängig davon
bereits verkörpert), gehört aber zur selben Klasse wie `slice-099`
(dort ein ADR-Tabellenwert „ohne Ursprung", hier eine Plan-Prosa-Zeile
„ohne Neumessung").

**Behoben:** Beide Fundstellen im Slice-Plan (§6, §8) sind im Rahmen dieser
Closure auf „7×" korrigiert.

Quelle: Verifikationsbericht zu `slice-101`, §5 (Abschnitt „Abweichung bei
§3.12") · `docs/plan/planning/done/altbestand/slice-101-nats-client-csharp-kotlin.md`
§6/§8.
