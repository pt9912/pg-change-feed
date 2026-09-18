# Beleg: slice-102

Vorgang: `slice-102` — gRPC-Client in C# (`ADR-0090`).

Fund: Der Slice-Plan selbst (§6 letzter Risiko-Punkt vor der Closure-Fassung
und §8 „Vorgelagert — offene Beobachtungen sichten") zitierte
`BEO-PGC/github-actions-unverifizierbar-lokal` mit „**5×**". Der Verifier hat
den abgeleiteten Zähler eigenständig nachgemessen
(Verifikationsbericht zu `slice-102`, #10, `ls evidence/`) und real **7×**
gefunden (`evidence/slice-039.md`, `-056.md`, `-064.md`, `-082.md`,
`-090.md`, `-098.md`, `-099.md`) — bereits zum Zeitpunkt der
Plan-Niederschrift (2026-09-17, nach `slice-098`/`-099`), also eine
**übernommene, nicht neu gemessene Zahl** statt einer Ableitung aus den
tatsächlich vorhandenen Belegdateien. Dieselbe Form wie
`evidence/slice-101.md` in diesem Register-Eintrag — beide Slices entstanden
praktisch zeitgleich, dieselbe veraltete Zahl wanderte in beide Pläne.

**Abgrenzung zu `BEO-PGC/arbeit-ueberholt-stehenden-traeger`:** Dort
überholt die *Arbeit dieses Slice* einen fremden, bei Niederschrift noch
korrekten Träger (`harness/README.md`, `slice-097`). Hier ist es umgekehrt:
Die im Plan selbst genannte Zahl war **schon bei ihrer eigenen
Niederschrift** falsch — kein Drift durch spätere Arbeit, sondern eine
Übernahme, die nie gegen die tatsächliche Beleg-Liste nachgezählt wurde.
Das ändert am Ausgang nichts (der Eintrag `github-actions-unverifizierbar-
lokal` ist unabhängig davon bereits verkörpert und bleibt bei real 7×),
gehört aber zur selben Klasse wie `slice-099`/`slice-101`.

**Behoben:** Beide Fundstellen im Slice-Plan (§6, §8) sind im Rahmen dieser
Closure auf „7×" korrigiert.

Quelle: Verifikationsbericht zu `slice-102`, §3 (Abschnitt „Abweichung bei
`AGENTS.md` §3.12") · `docs/plan/planning/done/slice-102-grpc-client-
csharp.md` §6/§8.
